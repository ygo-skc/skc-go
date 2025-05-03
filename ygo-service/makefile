# vars
GENERATED_DIR := ./pb

# default
all: clean build

build:
	go mod tidy
	env GOOS=linux GOARCH=arm64 go build .

clean:
	rm -f ygo-service

help:
	@echo "Makefile Commands:"
	@echo "  all            - Default target. Installs deps, cleans, and builds the binary."
	@echo "  clean          - Remove all generated files and the binary."
	@echo "  build          - Compile the Go project for Linux ARM64."
	@echo "  help           - Show this help message."