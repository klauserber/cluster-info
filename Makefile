# Makefile for Cluster Info

# Define variables
GOPATH := $(shell go env GOPATH)
REGISTRY_NAME ?= isi006
IMAGE_NAME ?= cluster-info
VERSION ?= 1.0.0

# Default target
all: build run

# Build the application
build:
	@echo "Building the application..."
	go build -o bin/cluster-info main.go

run: build
	@echo "Running the application..."
	./bin/cluster-info

# Run the application locally in docker
docker-run:
	@echo "Running the application locally..."
	# docker run -p 8080:8080 --rm -it --entrypoint sh ${REGISTRY_NAME}/${IMAGE_NAME}:${VERSION}
	docker run -p 8080:8080 --rm ${REGISTRY_NAME}/${IMAGE_NAME}:${VERSION}

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build \
		. \
  	-t ${REGISTRY_NAME}/${IMAGE_NAME}:${VERSION} \
  	-t ${REGISTRY_NAME}/${IMAGE_NAME}:latest \
  	--platform=linux/arm64,linux/amd64

# Push Docker image to registry
docker-push: docker-build
	@echo "Pushing Docker image to registry..."
	docker push ${REGISTRY_NAME}/${IMAGE_NAME}:${VERSION}
	docker push ${REGISTRY_NAME}/${IMAGE_NAME}:latest

k8s-deploy: docker-push
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s.yaml

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	go mod tidy

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up build artifacts..."
	rm -f bin/cluster-info
	docker rmi -f ${REGISTRY_NAME}/${IMAGE_NAME}:${VERSION} || true
	docker rmi -f ${REGISTRY_NAME}/${IMAGE_NAME}:latest || true

.PHONY: all build run docker-run docker-build docker-push update-deps test clean