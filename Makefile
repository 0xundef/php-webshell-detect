# Makefile for webshell-detect

APP_NAME    := webshell-detect
CMD_DIR     := ./cmd/$(APP_NAME)
OUTPUT_DIR  := ./output
CONFIG_DIR  := ./config

VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME  := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION  := $(shell go version | awk '{print $$3}')

LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.GoVersion=$(GO_VERSION) -w -s"

# Handwritten packages. The vendored PHP frontend is excluded from vet/lint
# because it contains generated ragel output (see frontend/php/README.md).
GO_PACKAGES := ./cmd/... ./internal/... ./test/...

.PHONY: all
all: clean build

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@mkdir -p $(OUTPUT_DIR)

.PHONY: build
build:
	@echo "Building $(APP_NAME) for the current platform..."
	@mkdir -p $(OUTPUT_DIR)
	@go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME) $(CMD_DIR)

.PHONY: build-linux
build-linux:
	@echo "Building $(APP_NAME) for Linux (amd64)..."
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-linux $(CMD_DIR)

.PHONY: build-macos
build-macos:
	@echo "Building $(APP_NAME) for macOS (amd64)..."
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-macos $(CMD_DIR)

.PHONY: build-windows
build-windows:
	@echo "Building $(APP_NAME) for Windows (amd64)..."
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-windows.exe $(CMD_DIR)

.PHONY: build-all
build-all: build-linux build-macos build-windows

# package-<os> builds and writes a tarball/zip containing config/conf.yaml and
# config/mock.php next to the binary.
.PHONY: package-linux
package-linux: build-linux
	@python3 ./scripts/package.py linux $(VERSION)

.PHONY: package-macos
package-macos: build-macos
	@python3 ./scripts/package.py macos $(VERSION)

.PHONY: package-windows
package-windows: build-windows
	@python3 ./scripts/package.py windows $(VERSION)

.PHONY: package-all
package-all: package-linux package-macos package-windows

.PHONY: release
release: clean package-all
	@echo "Release $(VERSION) completed:"
	@ls -lh $(OUTPUT_DIR)/

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: vet
vet:
	@echo "Running go vet on handwritten packages..."
	@go vet $(GO_PACKAGES)

# Repo-wide vet. The generated scanner contains unreachable code by design, so
# the unreachable analyzer is disabled instead of editing generated files.
.PHONY: vet-all
vet-all:
	@echo "Running go vet on all packages..."
	@go vet -unreachable=false ./...

.PHONY: lint
lint:
	@echo "Running golangci-lint on handwritten packages..."
	@golangci-lint run $(GO_PACKAGES)

.PHONY: fmt
fmt:
	@echo "Formatting handwritten Go code..."
	@gofmt -w cmd internal test

.PHONY: fmt-check
fmt-check:
	@echo "Checking formatting..."
	@test -z "$$(gofmt -l cmd internal test)" || (gofmt -l cmd internal test; exit 1)

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: deps
deps:
	@go mod download

.PHONY: run
run: build
	@$(OUTPUT_DIR)/$(APP_NAME)

.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(GO_VERSION)"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all/all build  - Clean and build for the current platform"
	@echo "  build-*        - Build for linux/macos/windows (amd64)"
	@echo "  build-all      - Build for all supported platforms"
	@echo "  package-*      - Package one platform with config/"
	@echo "  package-all    - Package all builds"
	@echo "  release        - Clean, build and package everything"
	@echo "  test           - Run all Go tests"
	@echo "  test-coverage  - Run tests and write coverage.html"
	@echo "  vet/vet-all    - go vet handwritten / all packages"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt/fmt-check  - Format / verify formatting"
	@echo "  tidy/deps      - Tidy or download modules"
	@echo "  clean          - Remove build artifacts"
	@echo "  version        - Show version information"