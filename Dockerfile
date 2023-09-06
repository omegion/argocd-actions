# Use the official Golang image from the DockerHub
FROM golang:1.17 as builder

# Install git, required for fetching Go dependencies
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Set the Current Working Directory inside the container
WORKDIR /app
COPY . .

# Download all Go dependencies
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o argocd-sync

FROM alpine:3.14
COPY --from=builder /app/argocd-sync /argocd-sync

ENTRYPOINT ["/argocd-sync"]
