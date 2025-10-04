# Default target
.DEFAULT_GOAL := help

# Build configuration
BINARY_NAME := wild-cli
VERSION ?= 0.1.0
BUILD_DIR := build
DIST_DIR := dist

# Go build configuration
GO_VERSION := $(shell go version | cut -d' ' -f3)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)

.PHONY: help build clean test run install build-arm64 build-amd64 build-all check fmt vet deps-check version

help:
	@echo "💻 Wild CLI Build System"
	@echo ""
	@echo "📦 Build targets:"
	@echo "  build           - Build for current architecture"
	@echo "  build-arm64     - Build arm64 binary"
	@echo "  build-amd64     - Build amd64 binary"
	@echo "  build-all       - Build all architectures"
	@echo ""
	@echo "🔍 Quality assurance:"
	@echo "  check           - Run all checks (fmt + vet + test)"
	@echo "  fmt             - Format Go code"
	@echo "  vet             - Run go vet"
	@echo "  test            - Run tests"
	@echo ""
	@echo "🛠️  Development:"
	@echo "  run             - Run CLI locally"
	@echo "  clean           - Remove all build artifacts"
	@echo "  deps-check      - Verify and tidy dependencies"
	@echo "  version         - Show build information"
	@echo "  install         - Install to system"

build:
	@echo "Building $(BINARY_NAME) for current architecture..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-arm64:
	@echo "Building $(BINARY_NAME) for arm64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-arm64 .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)-arm64"

build-amd64:
	@echo "Building $(BINARY_NAME) for amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-amd64 .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)-amd64"

build-all: build-arm64 build-amd64

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@go clean
	@echo "✅ Clean complete"

test:
	@echo "🧪 Running tests..."
	@go test -v ./...

run:
	@echo "🚀 Running $(BINARY_NAME)..."
	@go run -ldflags="$(LDFLAGS)" .

# Code quality targets
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Format complete"

vet:
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "✅ Vet complete"

check: fmt vet test
	@echo "✅ All checks passed"

# Dependency management
deps-check:
	@echo "📦 Checking dependencies..."
	@go mod verify
	@go mod tidy
	@echo "✅ Dependencies verified"

# Version information
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ $(BINARY_NAME) installed to /usr/local/bin/"