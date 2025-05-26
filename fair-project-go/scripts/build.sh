#!/bin/bash

# Build script for fair-project-go

# Set variables
BINARY_SERVER="fair-project-server"
BINARY_CLI="fair-ctl"
BUILD_DIR="bin"
PACKAGE="github.com/sadeq/fair-project-go"
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS="-ldflags=\"-X '${PACKAGE}/pkg/version.BranchHash=${COMMIT_HASH}'\""

# Create build directory if it doesn't exist
mkdir -p ${BUILD_DIR}

# Function to build the server
build_server() {
    echo "Building server binary..."
    eval "go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_SERVER} ./cmd/server"
}

# Function to build the CLI
build_cli() {
    echo "Building CLI binary..."
    eval "go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_CLI} ./cmd/cli"
}

# Function to run tests
run_tests() {
    echo "Running tests..."
    go test ./...
}

# Function to clean build artifacts
clean() {
    echo "Cleaning build artifacts..."
    rm -rf ${BUILD_DIR}
}

# Function to build Docker image
build_docker() {
    echo "Building Docker image..."
    docker build -t fair-project-server .
}

# Function to show help
show_help() {
    echo "Usage: $0 [command]"
    echo "Commands:"
    echo "  server    : Build the server binary"
    echo "  cli       : Build the CLI binary"
    echo "  all       : Build both server and CLI binaries (default)"
    echo "  test      : Run tests"
    echo "  clean     : Clean build artifacts"
    echo "  docker    : Build Docker image"
    echo "  help      : Show this help message"
}

# Parse command line arguments
case "$1" in
    server)
        build_server
        ;;
    cli)
        build_cli
        ;;
    all|"")
        build_server
        build_cli
        ;;
    test)
        run_tests
        ;;
    clean)
        clean
        ;;
    docker)
        build_docker
        ;;
    help)
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        show_help
        exit 1
        ;;
esac

exit 0
