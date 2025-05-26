# Makefile for fair-project-go

# Variables
BINARY_SERVER = fair-project-server
BINARY_CLI = fair-project-cli
BUILD_DIR = bin
PACKAGE = github.com/sadeq/fair-project-go
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS = -ldflags="-X '$(PACKAGE)/pkg/version.BranchHash=$(COMMIT_HASH)'"
GO_FILES = $(shell find . -name "*.go" -not -path "./vendor/*")

# Default target
.PHONY: all
all: build

# Build both server and CLI binaries
.PHONY: build
build: build-server build-cli

# Build the server binary
.PHONY: build-server
build-server:
	@echo "Building server binary..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_SERVER) ./cmd/server

# Build the CLI binary
.PHONY: build-cli
build-cli:
	@echo "Building CLI binary..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CLI) ./cmd/cli

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Build Docker image
.PHONY: docker
docker:
	@echo "Building Docker image..."
	docker build -t fair-project-server .

# Run the server
.PHONY: run-server
run-server: build-server
	@echo "Running server..."
	./$(BUILD_DIR)/$(BINARY_SERVER)

# Run the CLI
.PHONY: run-cli
run-cli: build-cli
	@echo "Running CLI..."
	./$(BUILD_DIR)/$(BINARY_CLI) $(ARGS)

# Show version information
.PHONY: version
version: build-server
	@echo "Version information:"
	./$(BUILD_DIR)/$(BINARY_SERVER) version

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          : Build both server and CLI binaries (default)"
	@echo "  build        : Build both server and CLI binaries"
	@echo "  build-server : Build the server binary"
	@echo "  build-cli    : Build the CLI binary"
	@echo "  test         : Run tests"
	@echo "  clean        : Clean build artifacts"
	@echo "  docker       : Build Docker image"
	@echo "  run-server   : Build and run the server"
	@echo "  run-cli      : Build and run the CLI (use ARGS=\"arg1 arg2\" to pass arguments)"
	@echo "  version      : Show version information"
	@echo "  help         : Show this help message"
