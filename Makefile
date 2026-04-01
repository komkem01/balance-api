APP_NAME := balance-api
BIN_DIR := bin
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.0.0)
LDFLAGS := -X balance/internal/config.version=$(VERSION)

.PHONY: help test run build build-release tag-release

help:
	@echo "Targets:"
	@echo "  make test                       - Run go tests"
	@echo "  make run                        - Run HTTP server (default version)"
	@echo "  make build                      - Build binary (default version)"
	@echo "  make build-release VERSION=vX.Y.Z - Build binary with injected version"
	@echo "  make tag-release VERSION=vX.Y.Z - Create and push annotated git tag"

test:
	go test ./...

run:
	go run . http

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) .

build-release:
	mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-$(VERSION) .
	@echo "Built $(BIN_DIR)/$(APP_NAME)-$(VERSION)"

tag-release:
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required"; exit 1; fi
	@echo "Creating tag $(VERSION)"
	git tag -a $(VERSION) -m "release $(VERSION)"
	git push origin $(VERSION)
