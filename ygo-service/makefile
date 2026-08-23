.PHONY: help

# default
all: build

build:
	go mod tidy
	go vet ./...
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM64=v8.2 go build -ldflags="-s -w" -trimpath .

help:
	@echo "Makefile Commands:"
	@echo "  all            - Default target. Installs deps, cleans, and builds the binary."
	@echo "  build          - Compile the Go project for Linux ARM64."
	@echo "  help           - Show this help message."
