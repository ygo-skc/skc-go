# vars
GENERATED_DIR := ./pb

# default
all: deps clean generate-grpc build

generate-grpc:
	protoc --go_out=./ --go-grpc_out=./ ygo-service.proto

build:
	env GOOS=linux GOARCH=arm64 go build .

clean:
	rm -rf $(GENERATED_DIR)/*.pb.go
	rm -f ygo-service

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go mod tidy

help:
	@echo "Makefile Commands:"
	@echo "  all            - Default target. Installs deps, cleans, generates code, and builds the binary."
	@echo "  deps           - Install Protobuf plugins for Go."
	@echo "  clean          - Remove all generated files and the binary."
	@echo "  generate-grpc  - Generate Go and gRPC code from .proto files."
	@echo "  build          - Compile the Go project for Linux ARM64."
	@echo "  help           - Show this help message."