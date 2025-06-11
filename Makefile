# Default target
.DEFAULT_GOAL := help

# Build configuration
BINARY_NAME := wild-cloud-central
VERSION ?= 0.1.1
BUILD_DIR := build
DIST_DIR := dist
DEB_DIR := debian-package

# Go build configuration
GO_VERSION := $(shell go version | cut -d' ' -f3)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)

.PHONY: help build clean test run install package repo deploy-repo build-arm64 build-amd64 package-arm64 package-amd64 package-all deb deb-arm64 deb-amd64 deb-all check fmt vet lint deps-check version

# Function to create debian package for specific architecture
# Usage: $(call package_deb,architecture,binary_name)
help:
	@echo "🏗️  Wild Cloud Central Build System"
	@echo ""
	@echo "📦 Build targets (compile binaries):"
	@echo "  build           - Build for current architecture"
	@echo "  build-arm64     - Build arm64 binary"
	@echo "  build-amd64     - Build amd64 binary"
	@echo "  build-all       - Build all architectures"
	@echo ""
	@echo "📋 Package targets (create .deb packages):"
	@echo "  package         - Create .deb package for current arch"
	@echo "  package-arm64   - Create arm64 .deb package"
	@echo "  package-amd64   - Create amd64 .deb package"
	@echo "  package-all     - Create all .deb packages"
	@echo ""
	@echo "🚀 Repository targets:"
	@echo "  repo            - Build APT repository from packages"
	@echo "  deploy-repo     - Deploy repository to server"
	@echo ""
	@echo "🔍 Quality assurance:"
	@echo "  check           - Run all checks (fmt + vet + test)"
	@echo "  fmt             - Format Go code"
	@echo "  vet             - Run go vet"
	@echo "  test            - Run tests"
	@echo ""
	@echo "🛠️  Development:"
	@echo "  run             - Run application locally"
	@echo "  clean           - Remove all build artifacts"
	@echo "  deps-check      - Verify and tidy dependencies"
	@echo "  version         - Show build information"
	@echo "  install         - Install to system"
	@echo ""
	@echo "📜 Legacy aliases (deprecated):"
	@echo "  deb, deb-arm64, deb-amd64, deb-all"
	@echo ""
	@echo "📁 Directory structure:"
	@echo "  build/          - Intermediate build artifacts"
	@echo "  dist/bin/       - Final binaries for distribution"
	@echo "  dist/packages/  - OS packages (.deb files)"
	@echo "  dist/repositories/ - APT repository for deployment"
	@echo ""
	@echo "💡 Example workflows:"
	@echo "  make check && make build     - Safe development build"
	@echo "  make clean && make repo      - Full release build"

define package_deb
	@echo "📦 Creating .deb package for $(1)..."
	@mkdir -p $(DIST_DIR)/bin $(DIST_DIR)/packages
	
	@# Copy debian package structure
	@cp -r debian/ $(BUILD_DIR)/$(DEB_DIR)-$(1)/
	
	@# Copy binary to correct location
	@cp $(BUILD_DIR)/$(2) $(BUILD_DIR)/$(DEB_DIR)-$(1)/usr/bin/$(BINARY_NAME)
	
	@# Copy static web files
	@cp -r static/* $(BUILD_DIR)/$(DEB_DIR)-$(1)/var/www/html/wild-central/
	@chmod 755 $(BUILD_DIR)/$(DEB_DIR)-$(1)/DEBIAN/postinst
	@chmod 755 $(BUILD_DIR)/$(DEB_DIR)-$(1)/DEBIAN/prerm
	@chmod 755 $(BUILD_DIR)/$(DEB_DIR)-$(1)/DEBIAN/postrm
	
	@# Substitute placeholders in control file
	@sed -i 's/VERSION_PLACEHOLDER/$(VERSION)/g' $(BUILD_DIR)/$(DEB_DIR)-$(1)/DEBIAN/control
	@sed -i 's/ARCH_PLACEHOLDER/$(1)/g' $(BUILD_DIR)/$(DEB_DIR)-$(1)/DEBIAN/control
	
	
	@# Build package and copy to dist directories
	dpkg-deb --build $(BUILD_DIR)/$(DEB_DIR)-$(1) $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_$(1).deb
	@cp $(BUILD_DIR)/$(2) $(DIST_DIR)/bin/$(BINARY_NAME)-$(1)
	@cp $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_$(1).deb $(DIST_DIR)/packages/
	@echo "✅ Package created: $(DIST_DIR)/packages/$(BINARY_NAME)_$(VERSION)_$(1).deb"
	@echo "✅ Binary copied: $(DIST_DIR)/bin/$(BINARY_NAME)-$(1)"
endef

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

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(DEB_DIR)-* $(DEB_DIR)
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
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/bin/
	sudo cp wild-cloud-central.service /etc/systemd/system/
	sudo mkdir -p /etc/wild-cloud-central
	sudo cp config.yaml.example /etc/wild-cloud-central/config.yaml.example
	sudo systemctl daemon-reload

# Build targets - compile binaries only
build: build-amd64

build-all: build-arm64 build-amd64

# Package targets - create .deb packages (requires binaries)
package-arm64: build-arm64
	$(call package_deb,arm64,$(BINARY_NAME)-arm64)

package-amd64: build-amd64
	$(call package_deb,amd64,$(BINARY_NAME)-amd64)

package-all: package-arm64 package-amd64

package: package-amd64

repo: package-all
	./scripts/build-apt-repository.sh

deploy-repo: repo
	./scripts/deploy-apt-repository.sh

dev:
	go run . &
	echo "Server started on http://localhost:5055"