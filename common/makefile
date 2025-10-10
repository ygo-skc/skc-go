.PHONY: deps clean help

# vars
YGO_GRPC_DIR := ./ygo
API_HEALTH_GRPC_DIR := ./health

# default
all: deps clean generate-grpc

generate-grpc:
	protoc --go_out=./ --go-grpc_out=./ common.proto ygo_service.proto health.proto
	go vet ./...

clean:
	rm -rf $(YGO_GRPC_DIR)/*.pb.go $(API_HEALTH_GRPC_DIR)/*.pb.go

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

help:
	@echo "Makefile Commands:"
	@echo "  all            - Default target. Installs deps, cleans, generates code, etc."
	@echo "  deps           - Install Protobuf plugins for Go."
	@echo "  clean          - Remove all generated files and the binary."
	@echo "  generate-grpc  - Generate Go and gRPC code from .proto files."
	@echo "  help           - Show this help message."