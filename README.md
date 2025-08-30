# cccluster: PoC for k3d, MinIO, and Go Services

This project is a **proof of concept (PoC)** for experimenting with Kubernetes (using k3d), MinIO object storage, and Go-based microservices. It is designed for local development and learning, with a focus on practical experimentation and rapid iteration.

## Overview

- **Kubernetes with k3d:** Lightweight local Kubernetes cluster for development and testing.
- **MinIO:** S3-compatible object storage deployed in the cluster for experimentation.
- **Go Services:** Example Go microservices (such as `storage-sync`) for learning and prototyping.
- **Podman:** Used for building and running container images locally (as a Docker alternative). However for simplicity the images are exported to docker (to simplify work with k3d registry).

## Structure

```
input/           # Input files and deployment manifests
minio/           # MinIO deployment and secret manifests
output/          # Output files and manifests
services/
   storage-sync/  # Example Go service with Dockerfile, k8s manifests, and source code
```

## Key Components

### MinIO
- Deployed as a Pod and Service in its own namespace (`minio`).
- Credentials are managed via Kubernetes Secrets (see `minio/minio-secret.example.yaml`).
- Exposes S3 API and web UI via NodePort for local access.

### Storage-Sync Service
- Written in Go, containerized with Podman.
- Watches a shared volume for file changes and syncs with a MinIO bucket.
- Includes Kubernetes manifests for running as a CronJob and for volume sharing with other pods.

## Getting Started

1. **Create k3d cluster with registry:**
   ```sh
   k3d cluster create my-cluster --agents 2 -v ~/k3d/minio-data/:/data --registry-create myregistry.localhost:5000 -p 8090:80@loadbalancer
   ```

2. **Build and push images:**
   ```sh
   cd services/storage-sync
   make pm-build  # Build with podman
   make pm-push   # Export to docker and push to registry
   ```

3. **Create namespaces:**
   ```sh
   kubectl apply -f minio/k8s/namespaces.yaml
   ```
   Creates:
   - `minio` namespace for minio application
   - `test` namespace for other application, custom Go apps, etc. 

4. **Create MinIO secrets:**
   - Copy `minio/k8s/minio-secret.example.yaml` to `minio/k8s/minio-secret.yaml` and fill in your base64-encoded credentials.
   - Apply the secret:
     ```sh
     kubectl apply -f minio/k8s/minio-secret.yaml
     ```

5. **Deploy MinIO:**
   ```sh
   kubectl apply -f minio/k8s/minio-dev.yaml
   ```

6. **Deploy storage-sync and sample pods:**
   ```sh
   kubectl apply -f services/storage-sync/k8s/sync-sample.yaml
   ```
   
   This creates:
   - A PersistentVolume and PersistentVolumeClaim for shared storage
   - A writer pod that continuously writes to the shared volume
   - A CronJob that runs every 5 minutes to sync data to MinIO

7. **Access MinIO UI:**
   - Open [http://localhost:30090](http://localhost:30090) in your browser (or http://localhost:8090 via loadbalancer).
   - Login with the credentials you set in the secret.


## Notes
- This project is for local development and experimentation only.
- Secrets are not committed to version control; use the provided example file.
- Podman is used for local image builds; images are loaded into k3d via the `localhost/` prefix.
- The storage-sync service is a starting point for experimenting with Go and Kubernetes patterns.

## License
MIT
