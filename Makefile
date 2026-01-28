# Makefile for openCMD project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
VERSION=$(shell git rev-parse --short HEAD)-$(shell git branch --show-current)

# Binary name
BINARY_NAME=openCMD

# Build directory
BUILD_DIR=./build

# Default target
all: deps test

# Build the binary (if main.go exists)
build:
	@if [ -f "./main.go" ]; then \
		echo "Building $(BINARY_NAME)..."; \
		mkdir -p $(BUILD_DIR); \
		$(GOBUILD) -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./main.go; \
		echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"; \
		echo "Version: $(VERSION)"; \
	else \
		echo "No main.go file found, skipping build..."; \
	fi

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) ./...
	@echo "Tests completed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean completed"

# Update dependencies
deps:
	@echo "Updating dependencies..."
	@$(GOMOD) tidy
	@echo "Dependencies updated"

# Run the application (if main.go exists)
run:
	@if [ -f "./main.go" ]; then \
		echo "Running $(BINARY_NAME)..."; \
		mkdir -p $(BUILD_DIR); \
		$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./main.go; \
		$(BUILD_DIR)/$(BINARY_NAME); \
	else \
		echo "No main.go file found, cannot run..."; \
	fi

# Check code formatting
fmt:
	@echo "Checking code formatting..."
	@$(GOCMD) fmt ./...
	@echo "Code formatting checked"

# Run linting (requires golangci-lint)
lint:
	@echo "Running linting..."
	@golangci-lint run
	@echo "Linting completed"

# Echo version information
echo-version:
	@echo "Version: $(VERSION)"

.PHONY: all build test clean deps run fmt lint echo-version