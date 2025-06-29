# gopenapi Makefile
# A comprehensive Makefile for the OpenAPI code generator

# Project configuration
PROJECT_NAME := gopenapi
BINARY_NAME := gopenapi
MODULE_NAME := github.com/shubhamku044/gopenapi
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build configuration
BUILD_DIR := bin
GENERATED_DIR := generated
SAMPLE_API := sample-api.yaml
INSTALL_PATH := /usr/local/bin

# Go configuration
GO := go
GOFLAGS := -v
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commitHash=$(COMMIT_HASH)"

# Cross-compilation targets
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

.PHONY: help build install uninstall clean test lint fmt vet deps update-deps \
        generate-sample run-sample docker-build release cross-compile \
        dev-setup check all run-generated build-generated release-check tag-release

# Default target
all: clean fmt vet test build

help: ## Show this help message
	@echo "$(CYAN)$(PROJECT_NAME) - OpenAPI Code Generator$(RESET)"
	@echo ""
	@echo "$(GREEN)Available targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Examples:$(RESET)"
	@echo "  make build                    # Build the binary"
	@echo "  make install                  # Build and install globally"
	@echo "  make demo-workflow            # Demonstrate complete user workflow"
	@echo "  make run-example             # Build and run the example server"
	@echo "  make dev-setup               # Set up development environment"
	@echo "  make release                 # Build for all platforms"
	@echo ""
	@echo "$(GREEN)Development targets:$(RESET)"
	@echo "  make test                    # Run all tests"
	@echo "  make test-coverage           # Run tests with coverage"
	@echo "  make lint                    # Run linter"
	@echo "  make fmt                     # Format code"
	@echo "  make clean                   # Clean build artifacts"
	@echo ""
	@echo "$(GREEN)Release targets:$(RESET)"
	@echo "  make release-check           # Run pre-release checks"
	@echo "  make release                 # Build release binaries for all platforms"
	@echo "  make tag-release             # Create and push a new git tag"

# Build targets
build: ## Build the gopenapi binary
	@echo "$(GREEN)Building $(BINARY_NAME)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(PROJECT_NAME)
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(RESET)"

build-debug: ## Build with debug symbols and race detection
	@echo "$(GREEN)Building $(BINARY_NAME) with debug symbols...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -race -gcflags "all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME)-debug ./cmd/$(PROJECT_NAME)
	@echo "$(GREEN)✓ Debug build complete: $(BUILD_DIR)/$(BINARY_NAME)-debug$(RESET)"

# Installation targets
install: build ## Build and install the binary globally
	@echo "$(GREEN)Installing $(BINARY_NAME) to $(INSTALL_PATH)...$(RESET)"
	@if [ -w "$(INSTALL_PATH)" ]; then \
		cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		echo "$(YELLOW)Requires sudo for installation to $(INSTALL_PATH)$(RESET)"; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME); \
	fi
	@echo "$(GREEN)✓ Installation complete. You can now use '$(BINARY_NAME)' from anywhere.$(RESET)"

uninstall: ## Remove the installed binary
	@echo "$(GREEN)Removing $(BINARY_NAME) from $(INSTALL_PATH)...$(RESET)"
	@if [ -w "$(INSTALL_PATH)" ]; then \
		rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	fi
	@echo "$(GREEN)✓ Uninstallation complete.$(RESET)"

# Development targets
dev-setup: ## Set up development environment
	@echo "$(GREEN)Setting up development environment...$(RESET)"
	$(GO) mod download
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "$(GREEN)✓ Development environment ready!$(RESET)"

deps: ## Download dependencies
	@echo "$(GREEN)Downloading dependencies...$(RESET)"
	$(GO) mod download
	$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies downloaded$(RESET)"

update-deps: ## Update dependencies
	@echo "$(GREEN)Updating dependencies...$(RESET)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(RESET)"

# Code quality targets
fmt: ## Format Go code
	@echo "$(GREEN)Formatting code...$(RESET)"
	$(GO) fmt ./...
	@which goimports >/dev/null 2>&1 && goimports -w . || true
	@echo "$(GREEN)✓ Code formatted$(RESET)"

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(RESET)"
	$(GO) vet ./...
	@echo "$(GREEN)✓ Vet check passed$(RESET)"

lint: ## Run linter
	@echo "$(GREEN)Running linter...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Run 'make dev-setup' first.$(RESET)"; \
	fi

staticcheck: ## Run staticcheck
	@echo "$(GREEN)Running staticcheck...$(RESET)"
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "$(YELLOW)staticcheck not installed. Run 'make dev-setup' first.$(RESET)"; \
	fi

check: fmt vet lint staticcheck ## Run all code quality checks

# Testing targets
test: ## Run tests
	@echo "$(GREEN)Running tests...$(RESET)"
	$(GO) test $(GOFLAGS) ./...
	@echo "$(GREEN)✓ Tests passed$(RESET)"

test-verbose: ## Run tests with verbose output
	@echo "$(GREEN)Running tests (verbose)...$(RESET)"
	$(GO) test -v ./...

test-race: ## Run tests with race detection
	@echo "$(GREEN)Running tests with race detection...$(RESET)"
	$(GO) test -race ./...

test-coverage: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(RESET)"
	$(GO) test -coverprofile=coverage.out ./... -coverpkg=./cmd/...,./internal/...,./pkg/...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(RESET)"
	@echo "$(CYAN)Coverage summary:$(RESET)"
	@$(GO) tool cover -func=coverage.out | tail -1

benchmark: ## Run benchmarks
	@echo "$(GREEN)Running benchmarks...$(RESET)"
	$(GO) test -bench=. -benchmem ./...

# Code generation targets
generate-example: build ## Generate code in the example project
	@echo "$(GREEN)Generating code in example project...$(RESET)"
	@cd example && ../$(BUILD_DIR)/$(BINARY_NAME) --spec=api.yaml
	@echo "$(GREEN)✓ Code generated in example/ directory$(RESET)"
	@echo "$(CYAN)Generated files:$(RESET)"
	@find example -type f -name "*.go" | grep -E "(generated|handlers|main)" | sort

# Cross-compilation targets
cross-compile: ## Build for all platforms
	@echo "$(GREEN)Cross-compiling for all platforms...$(RESET)"
	@mkdir -p $(BUILD_DIR)/releases
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		output="$(BUILD_DIR)/releases/$(BINARY_NAME)-$$os-$$arch"; \
		if [ "$$os" = "windows" ]; then output="$$output.exe"; fi; \
		echo "Building $$platform..."; \
		GOOS=$$os GOARCH=$$arch $(GO) build $(LDFLAGS) -o $$output ./cmd/$(PROJECT_NAME); \
	done
	@echo "$(GREEN)✓ Cross-compilation complete$(RESET)"

release: release-check
	@echo "🚀 Building release binaries..."
	@mkdir -p bin/releases
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then continue; fi; \
			echo "Building $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build -ldflags="-s -w" -o bin/releases/gopenapi-$$os-$$arch ./cmd/gopenapi; \
			if [ "$$os" = "windows" ]; then \
				mv bin/releases/gopenapi-$$os-$$arch bin/releases/gopenapi-$$os-$$arch.exe; \
			fi; \
		done; \
	done
	@echo "✅ Release binaries built in bin/releases/"

tag-release:
	@echo "📋 Current tags:"
	@git tag | tail -5
	@echo ""
	@read -p "Enter new version (e.g., v1.0.0): " version; \
	git tag $$version && \
	git push origin $$version && \
	echo "✅ Tagged and pushed $$version"

release-check:
	@echo "🔍 Pre-release checks..."
	@$(MAKE) test
	@$(MAKE) lint
	@echo "✅ All checks passed"

# Docker targets
docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(RESET)"
	docker build -t $(PROJECT_NAME):$(VERSION) -t $(PROJECT_NAME):latest .
	@echo "$(GREEN)✓ Docker image built$(RESET)"

# Utility targets
clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(RESET)"
	rm -rf $(BUILD_DIR)
	rm -rf demo-project example-basic example-custom
	rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Clean complete$(RESET)"

version: ## Show version information
	@echo "$(CYAN)$(PROJECT_NAME) version information:$(RESET)"
	@echo "  Version:     $(VERSION)"
	@echo "  Build time:  $(BUILD_TIME)"
	@echo "  Commit hash: $(COMMIT_HASH)"
	@echo "  Go version:  $$($(GO) version)"

show-config: ## Show build configuration
	@echo "$(CYAN)Build configuration:$(RESET)"
	@echo "  Project:     $(PROJECT_NAME)"
	@echo "  Module:      $(MODULE_NAME)"
	@echo "  Binary:      $(BINARY_NAME)"
	@echo "  Build dir:   $(BUILD_DIR)"
	@echo "  Install dir: $(INSTALL_PATH)"
	@echo "  Go version:  $$($(GO) version)"

# Example workflows
example-basic: build ## Run basic example in current directory
	@echo "$(GREEN)Running basic example...$(RESET)"
	@rm -rf example-basic && mkdir example-basic
	@cd example-basic && go mod init example-basic
	@cp example/api.yaml example-basic/
	@cd example-basic && ../$(BUILD_DIR)/$(BINARY_NAME) --spec=api.yaml
	@echo "$(GREEN)✓ Example complete. Check ./example-basic/$(RESET)"

example-custom: build ## Run example with custom package name
	@echo "$(GREEN)Running custom package example...$(RESET)"
	@rm -rf example-custom && mkdir example-custom
	@cd example-custom && go mod init github.com/myuser/custom-api
	@cp example/api.yaml example-custom/api.yaml
	@cd example-custom && ../$(BUILD_DIR)/$(BINARY_NAME) --spec=api.yaml --package=myapi
	@echo "$(GREEN)✓ Custom example complete. Check ./example-custom/$(RESET)"

# Watch for changes (requires entr or similar)
watch: ## Watch for changes and rebuild (requires 'entr')
	@if command -v entr >/dev/null 2>&1; then \
		echo "$(GREEN)Watching for changes... (Ctrl+C to stop)$(RESET)"; \
		find . -name "*.go" | entr -r make build; \
	else \
		echo "$(RED)Error: 'entr' not found. Install with: brew install entr$(RESET)"; \
	fi

# Help for specific file
inspect: ## Show information about generated files
	@if [ -d "example/generated" ]; then \
		echo "$(CYAN)Generated files in example/:$(RESET)"; \
		find example -type f -name "*.go" | while read file; do \
			echo "  $$file ($$(wc -l < "$$file") lines)"; \
		done; \
		if [ -f "example/api-server" ]; then \
			echo "$(CYAN)Compiled binary:$(RESET)"; \
			echo "  example/api-server ($$(ls -lh example/api-server | awk '{print $$5}'))"; \
		fi; \
	else \
		echo "$(YELLOW)No generated files found. Run 'make generate-example' first.$(RESET)"; \
	fi

# Install development tools
tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(RESET)"
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	$(GO) install github.com/goreleaser/goreleaser@latest
	@echo "$(GREEN)✓ Development tools installed$(RESET)" 