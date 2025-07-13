# Article Processing Microservice
# ================================================

# Variables
BINARY_NAME=article-server
BUILD_DIR=build

# Go commands
GO=go
GOBUILD=$(GO) build
GORUN=$(GO) run
GOTEST=$(GO) test
GOCLEAN=$(GO) clean
GOGET=$(GO) get

# Protobuf compiler
PROTOC=protoc


# Default target
.DEFAULT_GOAL := help

# Help - shows all available commands
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make proto     - Generate protobuf files"
	@echo "  make clean     - Clean build files"
	@echo "  make deps      - Install dependencies"


# Build the application
.PHONY: build
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
.PHONY: run
run:
	@echo "Starting article processing microservice..."
	$(GORUN) main.go

# Generate protobuf files
.PHONY: proto
proto:
	@echo "Generating protobuf files..."
	cd proto && $(PROTOC) --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative article.proto
	@echo "Protobuf files generated successfully"

# Clean build files
.PHONY: clean
clean:
	@echo "Cleaning build files..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete" 
	
# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOGET) -v ./...
	$(GO) mod tidy
	@echo "Dependencies installed"

