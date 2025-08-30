
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o storage-sync-cli ./cmd/cli.go

# Final image
FROM alpine:latest

# Copy only the compiled binary from the build stage
COPY --from=builder /app/storage-sync-cli /usr/local/bin/storage-sync-cli

# No entrypoint for now, will specify command in the pod/job template
