# Use the official Go image as the base
FROM golang:1.25.1-bookworm

# Set the Go working directory (standard practice)
WORKDIR /go/src/tools

RUN go mod init tools || true

# Set the environment variable for Go install location
ENV GO_BIN_PATH /usr/local/bin

# 1. Install Protocol Buffers compiler (protoc)
# We use apt to install protoc and curl, unzip (needed for protoc installation in case of a different base image)
# Note: golang:1.25.1-bookworm is based on Debian, so apt is used.
RUN apt-get update && apt-get install -y --no-install-recommends \
    protobuf-compiler \
    git \
    make \
    vim \
    && rm -rf /var/lib/apt/lists/*

# 2. Install Go Protobuf plugins
# This installs:
# - google.golang.org/protobuf/cmd/protoc-gen-go (protoc-gen-go)
# - google.golang.org/grpc/cmd/protoc-gen-go-grpc (protoc-gen-go-grpc)
# - github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway (protoc-gen-grpc-gateway)
# - github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 (protoc-gen-openapiv2)
#
# 'go install' places the binaries in $GOPATH/bin, which should be in the PATH
# for the default Go image.
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Set the working directory inside the container
WORKDIR /app
