.PHONY: all deps clean build generate-grpc help

# vars
API_HEALTH_GRPC_DIR := ./health

# default
all: deps clean generate-grpc build

generate-grpc:
	protoc --go_out=./ --go-grpc_out=./ health.proto

build:
	go mod tidy
	env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -trimpath .

clean:
	rm -rf "${API_HEALTH_GRPC_DIR}"

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

help:
	@echo "Makefile Commands:"
	@echo "  all            - Default target. Installs deps, cleans, and builds the binary."
	@echo "  deps           - Install Protobuf plugins for Go."
	@echo "  clean          - Remove all generated files and the binary."
	@echo "  build          - Compile the Go project for Linux ARM64."
	@echo "  generate-grpc  - Generate Go and gRPC code from .proto files."
	@echo "  help           - Show this help message."