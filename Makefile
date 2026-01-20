# DevSpace MCP Server Makefile

# Binary name
BINARY_NAME := devspace-mcp

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOCLEAN := $(GOCMD) clean
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet

# Build flags
LDFLAGS := -s -w
BUILD_FLAGS := -ldflags "$(LDFLAGS)"

# Directories
COVERAGE_DIR := coverage

.PHONY: all build test test-verbose test-coverage clean deps fmt vet lint run install help

# Default target
.DEFAULT_GOAL := help

all: fmt vet test build ## Run fmt, vet, test, and build

## Build Commands

build: ## Build the binary
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) .

build-debug: ## Build with debug symbols
	$(GOBUILD) -o $(BINARY_NAME) .

## Test Commands

test: ## Run tests
	$(GOTEST) -race ./...

test-verbose: ## Run tests with verbose output
	$(GOTEST) -race -v ./...

test-coverage: ## Run tests with coverage report
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

test-coverage-func: ## Show coverage by function
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out

## Code Quality

fmt: ## Format code
	$(GOFMT) -s -w .

fmt-check: ## Check code formatting (fails if not formatted)
	@test -z "$$($(GOFMT) -s -l . | tee /dev/stderr)" || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

vet: ## Run go vet
	$(GOVET) ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	golangci-lint run ./...

## Dependency Management

deps: ## Download dependencies
	$(GOMOD) download

deps-tidy: ## Tidy dependencies
	$(GOMOD) tidy

deps-verify: ## Verify dependencies
	$(GOMOD) verify

deps-update: ## Update all dependencies
	$(GOGET) -u ./...
	$(GOMOD) tidy

## Run & Install

run: build ## Build and run the server
	./$(BINARY_NAME)

install: build ## Install binary to GOPATH/bin
	$(GOCMD) install .

## Cleanup

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(COVERAGE_DIR)

## Development Helpers

dev-setup: deps ## Setup development environment
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development environment ready!"

check: fmt-check vet test ## Run all checks (CI)

## MCP Testing

mcp-test: build ## Test MCP server with initialize request
	@echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | ./$(BINARY_NAME)

mcp-list-tools: build ## List available MCP tools
	@echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' > /tmp/mcp-init.json
	@echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' >> /tmp/mcp-init.json
	@cat /tmp/mcp-init.json | ./$(BINARY_NAME)
	@rm -f /tmp/mcp-init.json

mcp-inspector: build ## Run MCP inspector (requires npx)
	@command -v npx >/dev/null 2>&1 || { echo "npx not installed. Install Node.js first."; exit 1; }
	npx @anthropic-ai/mcp-inspector ./$(BINARY_NAME)

## Help

help: ## Show this help
	@echo "DevSpace MCP Server - Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
