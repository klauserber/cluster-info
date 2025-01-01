# Cluster Info

## Overview

Cluster Info is a Go-based application designed to provide information about a Kubernetes cluster. It can be built, containerized, and deployed using Docker and Kubernetes.

Features:
* Show all ingresses as links

## Prerequisites

- **Go**: Ensure you have Go installed (version 1.23 or later).
- **Docker**: Ensure you have Docker installed.
- **Kubernetes Cluster**: Access to a Kubernetes cluster is required for deployment.
- **kubectl**: The `kubectl` command-line tool should be configured to interact with your Kubernetes cluster.

## Structure

The project consists of the following files and directories:

- `Makefile`: Makefile to build and run the application, Docker image.
- `go.mod`: Go module file that lists all dependencies.
- `main.go`: The main entry point of the application.
- `k8s.yaml`: Kubernetes deployment configuration for deploying the application.

## Building

Build locally, the binary is `bin/cluster-info`.

```bash
make build
```

Build the Docker image:

```bash
make docker-build
```

Push the Docker image to a registry:

```bash
make docker-push
```


## Running Locally

To run the application locally, you can use Docker:

```bash
docker run -p 8080:8080 ${REGISTRY_NAME}/${IMAGE_NAME}:${VERSION}
```

Visit `http://localhost:8080` in your web browser to access the Cluster Info application.

## Deployment

To deploy the application to a Kubernetes cluster, use the provided `k8s.yaml` file:

```bash
kubectl apply -f k8s.yaml
```

This will create an ingress resource and a service named `cluster-info-service`. The application will be accessible at `http://info-apps.bm.dc1.isium.de`.

## Development

- **Dependencies**: Update dependencies using Go modules with:
  ```bash
  go mod tidy
  ```

- **Running Tests**: If you have tests, you can run them using:
  ```bash
  go test ./...
  ```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.