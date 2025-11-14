## Simple Makefile for switcher

SHELL := /bin/bash

GO ?= go
BINARY ?= switcher
INSTALL_DIR ?= /usr/bin

# Version information (can be overridden)
# 默认使用最新的Git tag，如果没有tag则使用git describe格式，如果连git都没有则使用dev
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILT_BY ?= make

# Build flags with version information
LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(DATE) \
	-X main.builtBy=$(BUILT_BY)

.PHONY: build install clean build-all build-linux build-darwin build-windows \
	release-test release-snapshot tag help

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY) .

build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY)-linux-arm64 .

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY)-darwin-arm64 .

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY)-windows-arm64.exe .

install: build
	@echo "Installing $(BINARY) to $(INSTALL_DIR)"
	install -d -m 0755 $(INSTALL_DIR)
	install -m 0755 $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Done. Run: $(BINARY)"

clean:
	rm -f $(BINARY) $(BINARY)-*
	rm -rf dist/
	@echo "Cleaned build artifacts."

# Test GoReleaser configuration without publishing
release-test:
	@if ! command -v goreleaser &> /dev/null; then \
		echo "GoReleaser not found. Install it from: https://goreleaser.com/install/"; \
		exit 1; \
	fi
	goreleaser release --snapshot --clean --skip=publish

# Create a snapshot release (for testing)
release-snapshot:
	@if ! command -v goreleaser &> /dev/null; then \
		echo "GoReleaser not found. Install it from: https://goreleaser.com/install/"; \
		exit 1; \
	fi
	goreleaser release --snapshot --clean

# Helper to create and push a new version tag
tag:
	@echo "Current tags:"
	@git tag -l | tail -5
	@echo ""
	@read -p "Enter new version tag (e.g., v1.0.0): " TAG; \
	if [ -z "$$TAG" ]; then \
		echo "Tag cannot be empty"; \
		exit 1; \
	fi; \
	echo "Creating tag $$TAG..."; \
	git tag -a $$TAG -m "Release $$TAG"; \
	echo "Pushing tag $$TAG..."; \
	git push origin $$TAG; \
	echo "Done! GitHub Actions will now build and release $$TAG"

help:
	@echo "Available targets:"
	@echo "  build            - Build binary for current platform"
	@echo "  build-all        - Build binaries for all platforms"
	@echo "  install          - Install binary to $(INSTALL_DIR)"
	@echo "  clean            - Remove build artifacts"
	@echo "  release-test     - Test GoReleaser config locally (no publish)"
	@echo "  release-snapshot - Create snapshot release locally"
	@echo "  tag              - Create and push a new version tag"
	@echo "  help             - Show this help message"
