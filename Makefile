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
        dev-setup check all run-generated build-generated

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
	@echo "  make generate-sample          # Generate code from sample API"
	@echo "  make run-generated           # Build and run the generated server"
	@echo "  make dev-setup               # Set up development environment"
	@echo "  make release                 # Build for all platforms"

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
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(RESET)"

benchmark: ## Run benchmarks
	@echo "$(GREEN)Running benchmarks...$(RESET)"
	$(GO) test -bench=. -benchmem ./...

# Code generation targets
generate-sample: build ## Generate code from sample API
	@echo "$(GREEN)Generating code from $(SAMPLE_API)...$(RESET)"
	@rm -rf $(GENERATED_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) --spec=$(SAMPLE_API) --output=$(GENERATED_DIR) --package=sample
	@echo "$(GREEN)✓ Code generated in $(GENERATED_DIR)/$(RESET)"
	@echo "$(CYAN)Generated files:$(RESET)"
	@find $(GENERATED_DIR) -type f | sort

run-sample: generate-sample ## Build and run sample generation
	@echo "$(GREEN)Sample generation completed!$(RESET)"
	@echo "$(CYAN)Check the $(GENERATED_DIR)/ directory for generated code.$(RESET)"

# Generated code targets
build-generated: ## Build the generated API server
	@if [ ! -d "$(GENERATED_DIR)" ]; then \
		echo "$(RED)Error: No generated code found. Run 'make generate-sample' first.$(RESET)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Building generated API server...$(RESET)"
	@cd $(GENERATED_DIR) && go mod tidy && go build -o api-server .
	@echo "$(GREEN)✓ Generated server built: $(GENERATED_DIR)/api-server$(RESET)"

run-generated: build-generated ## Build and run the generated API server
	@echo "$(GREEN)Starting generated API server...$(RESET)"
	@echo "$(CYAN)Server will start on http://localhost:8080$(RESET)"
	@echo "$(CYAN)Available endpoints:$(RESET)"
	@echo "  GET  /health      - Health check"
	@echo "  GET  /users       - List users"
	@echo "  POST /users       - Create user"
	@echo "  GET  /users/:id   - Get user by ID"
	@echo ""
	@echo "$(YELLOW)Press Ctrl+C to stop the server$(RESET)"
	@cd $(GENERATED_DIR) && ./api-server

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

release: clean check cross-compile ## Build release artifacts for all platforms
	@echo "$(GREEN)Creating release artifacts...$(RESET)"
	@cd $(BUILD_DIR)/releases && \
	for binary in *; do \
		if [ "$$binary" != "*" ]; then \
			echo "Creating archive for $$binary..."; \
			tar -czf "$$binary.tar.gz" "$$binary"; \
		fi \
	done
	@echo "$(GREEN)✓ Release artifacts created in $(BUILD_DIR)/releases/$(RESET)"
	@ls -la $(BUILD_DIR)/releases/

# Docker targets
docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(RESET)"
	docker build -t $(PROJECT_NAME):$(VERSION) -t $(PROJECT_NAME):latest .
	@echo "$(GREEN)✓ Docker image built$(RESET)"

# Utility targets
clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(RESET)"
	rm -rf $(BUILD_DIR)
	rm -rf $(GENERATED_DIR)
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
example-basic: build ## Run basic example
	@echo "$(GREEN)Running basic example...$(RESET)"
	./$(BUILD_DIR)/$(BINARY_NAME) --spec=$(SAMPLE_API) --output=./example-output
	@echo "$(GREEN)✓ Example complete. Check ./example-output/$(RESET)"

example-custom: build ## Run example with custom package name
	@echo "$(GREEN)Running custom package example...$(RESET)"
	./$(BUILD_DIR)/$(BINARY_NAME) --spec=$(SAMPLE_API) --output=./custom-output --package=myapi
	@echo "$(GREEN)✓ Custom example complete. Check ./custom-output/$(RESET)"

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
	@if [ -d "$(GENERATED_DIR)" ]; then \
		echo "$(CYAN)Generated files:$(RESET)"; \
		find $(GENERATED_DIR) -type f -name "*.go" | while read file; do \
			echo "  $$file ($$(wc -l < "$$file") lines)"; \
		done; \
		if [ -f "$(GENERATED_DIR)/api-server" ]; then \
			echo "$(CYAN)Compiled binary:$(RESET)"; \
			echo "  $(GENERATED_DIR)/api-server ($$(ls -lh $(GENERATED_DIR)/api-server | awk '{print $$5}'))"; \
		fi; \
	else \
		echo "$(YELLOW)No generated files found. Run 'make generate-sample' first.$(RESET)"; \
	fi

# Install development tools
tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(RESET)"
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	$(GO) install github.com/goreleaser/goreleaser@latest
	@echo "$(GREEN)✓ Development tools installed$(RESET)" 