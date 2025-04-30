# Variables
GENERATED_DIR := ./pb

# Default target
all: clean generate

# Target to generate Go code from .proto files
generate:
	protoc --go_out=./ --go-grpc_out=./ ygo-service.proto

# Clean up generated files
clean:
	rm -rf $(GENERATED_DIR)/*.pb.go

# Install dependencies
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Help message
help:
	@echo "Makefile Commands:"
	@echo "  all         - Default target, generates Go code from Protobuf files."
	@echo "  generate    - Generate Go and gRPC code from .proto files."
	@echo "  clean       - Remove all generated .pb.go files."
	@echo "  deps        - Install Protobuf plugins for Go."
	@echo "  help        - Show this help message."