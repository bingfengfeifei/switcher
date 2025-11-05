## Simple Makefile for switcher

SHELL := /bin/bash

GO ?= go
BINARY ?= switcher
INSTALL_DIR ?= /usr/bin

.PHONY: build install clean

build:
	$(GO) build -o $(BINARY) .

install: build
	@echo "Installing $(BINARY) to $(INSTALL_DIR)"
	install -d -m 0755 $(INSTALL_DIR)
	install -m 0755 $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Done. Run: $(BINARY)"

clean:
	rm -f $(BINARY)
	@echo "Cleaned build artifacts."
