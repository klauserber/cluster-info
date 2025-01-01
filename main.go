package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type IngressInfo struct {
	Host      string
	Namespace string
}

type CategorizedIngresses struct {
	Categories map[string][]IngressInfo
}

func initKubernetesClient() (*kubernetes.Clientset, error) {
	// Determine if running inside or outside of the cluster
	inCluster := os.Getenv("KUBERNETES_SERVICE_HOST") != ""

	var config *rest.Config
	var err error

	if inCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		kubeconfigPath := os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}

	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig: %v", err)
	}

	// Create the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func main() {
	http.HandleFunc("/", handleIndex)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	clientset, err := initKubernetesClient()
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	err, categorizedIngresses := getIngresses(clientset)
	if err != nil {
		log.Printf("Error loading ingresses: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmplptr := template.Must(template.ParseFiles("templates/index.html"))
	err = tmplptr.Execute(w, categorizedIngresses)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func getIngresses(clientset *kubernetes.Clientset) (error, CategorizedIngresses) {
	ingresses, err := clientset.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err, CategorizedIngresses{}
	}

	categorizedIngresses := CategorizedIngresses{Categories: make(map[string][]IngressInfo)}
	for _, ingress := range ingresses.Items {
		skip, skipFound := ingress.Annotations["cluster-info.isium.de/skip"]
		if skipFound && skip == "true" {
			continue
		}

		category, catFound := ingress.Annotations["cluster-info.isium.de/category"]
		if !catFound || category == "" {
			category = "unsorted"
		}

		for _, rule := range ingress.Spec.Rules {
			categorizedIngresses.Categories[category] = append(categorizedIngresses.Categories[category], IngressInfo{Host: rule.Host, Namespace: ingress.Namespace})
		}
	}
	return nil, categorizedIngresses
}
