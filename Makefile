## Simple Makefile for switcher

SHELL := /bin/bash

GO ?= go
BINARY ?= switcher
INSTALL_DIR ?= /usr/bin

.PHONY: build install clean build-all build-linux build-darwin build-windows

build:
	$(GO) build -o $(BINARY) .

build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY)-linux-amd64 .

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINARY)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BINARY)-darwin-arm64 .

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY)-windows-amd64.exe .

install: build
	@echo "Installing $(BINARY) to $(INSTALL_DIR)"
	install -d -m 0755 $(INSTALL_DIR)
	install -m 0755 $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Done. Run: $(BINARY)"

clean:
	rm -f $(BINARY) $(BINARY)-*
	@echo "Cleaned build artifacts."
