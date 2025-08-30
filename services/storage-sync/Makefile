# Makefile commands for storage-sync

APP_NAME=storage-sync
IMAGE_NAME=storage-sync
BINARY_NAME=storage-sync-cli
REGISTRY=localhost:5000
K3D_REGISTRY=k3d-myregistry.localhost:5000

.PHONY: all build test clean pm-build pm-push pm-build-push run

all: build

build:
	go build -o $(BINARY_NAME) ./cmd/cli.go

test:
	go test ./storage-sync/...

clean:
	rm -f $(BINARY_NAME)


# Build Docker image using podman
pm-build:
	podman build -t $(IMAGE_NAME) .

# Tag and push the image to the local k3d registry using Docker (via podman export)
pm-push:
	podman save localhost/$(IMAGE_NAME):latest | docker load
	docker tag localhost/$(IMAGE_NAME):latest $(K3D_REGISTRY)/$(IMAGE_NAME):latest
	docker push $(K3D_REGISTRY)/$(IMAGE_NAME):latest

# Build and push in one step  
pm-build-push: pm-build pm-push

# Run the CLI directly (example)
run:
	./$(BINARY_NAME) --help
