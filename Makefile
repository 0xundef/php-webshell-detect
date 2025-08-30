# Makefile for webshell-detect

# Variables
APP_NAME := webshell-detect
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION := $(shell go version | awk '{print $$3}')

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.GoVersion=$(GO_VERSION) -w -s"

# Directories
CMD_DIR := ./cmd/$(APP_NAME)
OUTPUT_DIR := ./output
CONFIG_DIR := ./config

# Default target
.PHONY: all
all: clean build

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@mkdir -p $(OUTPUT_DIR)

# Build for current platform
.PHONY: build
build:
	@echo "Building $(APP_NAME) for current platform..."
	@go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME) $(CMD_DIR)

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building $(APP_NAME) for Linux..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-linux $(CMD_DIR)

# Build for macOS
.PHONY: build-macos
build-macos:
	@echo "Building $(APP_NAME) for macOS..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-macos $(CMD_DIR)

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building $(APP_NAME) for Windows..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-windows.exe $(CMD_DIR)

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-macos build-windows

# Package Linux build
.PHONY: package-linux
package-linux: build-linux
	@echo "Packaging Linux build..."
	@tar -czf $(OUTPUT_DIR)/$(APP_NAME)-linux-default-x86_64-$(VERSION).tar.gz -C $(OUTPUT_DIR) $(APP_NAME)-linux -C ../$(CONFIG_DIR) conf.yaml mock.php

# Package macOS build
.PHONY: package-macos
package-macos: build-macos
	@echo "Packaging macOS build..."
	@tar -czf $(OUTPUT_DIR)/$(APP_NAME)-macos-default-x86_64-$(VERSION).tar.gz -C $(OUTPUT_DIR) $(APP_NAME)-macos -C ../$(CONFIG_DIR) conf.yaml mock.php

# Package Windows build
.PHONY: package-windows
package-windows: build-windows
	@echo "Packaging Windows build..."
	@cd $(OUTPUT_DIR) && zip $(APP_NAME)-windows-default-x86_64-$(VERSION).zip $(APP_NAME)-windows.exe ../$(CONFIG_DIR)/conf.yaml ../$(CONFIG_DIR)/mock.php

# Package all builds
.PHONY: package-all
package-all: package-linux package-macos package-windows

# Release (build and package all)
.PHONY: release
release: clean build-all package-all
	@echo "Release $(VERSION) completed!"
	@ls -la $(OUTPUT_DIR)/

# Test
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Lint
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download

# Development build (with race detection)
.PHONY: dev
dev:
	@echo "Building development version..."
	@go build -race $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-dev $(CMD_DIR)

# Run the application
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	@$(OUTPUT_DIR)/$(APP_NAME)

# Show version info
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(GO_VERSION)"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build for current platform"
	@echo "  build        - Build for current platform"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-macos  - Build for macOS"
	@echo "  build-windows- Build for Windows"
	@echo "  build-all    - Build for all platforms"
	@echo "  package-*    - Package builds for specific platforms"
	@echo "  package-all  - Package all builds"
	@echo "  release      - Build and package everything"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  tidy         - Tidy dependencies"
	@echo "  deps         - Install dependencies"
	@echo "  dev          - Development build with race detection"
	@echo "  run          - Build and run the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  version      - Show version information"
	@echo "  help         - Show this help"