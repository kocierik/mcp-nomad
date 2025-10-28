# Makefile for mcp-nomad

# Variables
BINARY_NAME=mcp-nomad
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(shell go version | awk '{print $$3}')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.goVersion=$(GO_VERSION)"

# Default target
.DEFAULT_GOAL := help

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
PURPLE=\033[0;35m
CYAN=\033[0;36m
NC=\033[0m # No Color

# Enable color output
export TERM=xterm-256color

.PHONY: help build run clean test test-unit test-integration test-coverage test-race test-benchmark test-all clean-test test-deps lint format deps install uninstall release docker dev start-nomad stop-nomad status

# =============================================================================
# HELP
# =============================================================================

help: ## Show this help message
	@printf "$(CYAN)mcp-nomad - Nomad MCP Server$(NC)\n"
	@printf "$(YELLOW)Version: $(VERSION)$(NC)\n"
	@printf "\n"
	@printf "$(GREEN)Available targets:$(NC)\n"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(CYAN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# =============================================================================
# BUILD & RUN
# =============================================================================

build: ## Build the binary
	@printf "$(BLUE)Building $(BINARY_NAME)...$(NC)\n"
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@printf "$(GREEN)✅ Build completed: $(BINARY_NAME)$(NC)\n"

build-linux: ## Build for Linux
	@echo "$(BLUE)Building for Linux...$(NC)"
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux .
	@echo "$(GREEN)✅ Linux build completed: $(BINARY_NAME)-linux$(NC)"

build-darwin: ## Build for macOS
	@echo "$(BLUE)Building for macOS...$(NC)"
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin .
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .
	@echo "$(GREEN)✅ macOS builds completed$(NC)"

build-windows: ## Build for Windows
	@echo "$(BLUE)Building for Windows...$(NC)"
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows.exe .
	@GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-arm64.exe .
	@echo "$(GREEN)✅ Windows builds completed$(NC)"

build-all: build-linux build-darwin build-windows ## Build for all platforms

run: build ## Build and run the server
	@echo "$(BLUE)Starting $(BINARY_NAME) server...$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"
	@./$(BINARY_NAME)

run-stdio: build ## Run with stdio transport
	@echo "$(BLUE)Starting $(BINARY_NAME) with stdio transport...$(NC)"
	@./$(BINARY_NAME) -transport=stdio

run-sse: build ## Run with SSE transport
	@echo "$(BLUE)Starting $(BINARY_NAME) with SSE transport...$(NC)"
	@./$(BINARY_NAME) -transport=sse -port=8080

run-http: build ## Run with HTTP transport
	@echo "$(BLUE)Starting $(BINARY_NAME) with HTTP transport...$(NC)"
	@./$(BINARY_NAME) -transport=streamable-http -port=8080

dev: ## Run in development mode with hot reload
	@echo "$(BLUE)Starting development server...$(NC)"
	@echo "$(YELLOW)Installing air for hot reload...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@air

# =============================================================================
# INSTALLATION
# =============================================================================

install: build ## Install the binary to /usr/local/bin
	@echo "$(BLUE)Installing $(BINARY_NAME) to /usr/local/bin...$(NC)"
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)✅ $(BINARY_NAME) installed successfully$(NC)"

uninstall: ## Remove the binary from /usr/local/bin
	@echo "$(BLUE)Removing $(BINARY_NAME) from /usr/local/bin...$(NC)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ $(BINARY_NAME) uninstalled successfully$(NC)"

# =============================================================================
# DEPENDENCIES
# =============================================================================

deps: ## Download and tidy dependencies
	@echo -e "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo -e "$(GREEN)✅ Dependencies updated$(NC)"

deps-update: ## Update dependencies to latest versions
	@echo "$(BLUE)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies updated to latest versions$(NC)"

# =============================================================================
# TESTING
# =============================================================================

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	@echo -e "$(BLUE)Running unit tests...$(NC)"
	@go test -v ./test/unit/...

test-integration: ## Run integration tests
	@echo -e "$(BLUE)Running integration tests...$(NC)"
	@go test -v ./test/integration/...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

test-race: ## Run tests with race detection
	@echo "$(BLUE)Running tests with race detection...$(NC)"
	@go test -v -race ./test/...

test-benchmark: ## Run benchmark tests
	@echo "$(BLUE)Running benchmark tests...$(NC)"
	@go test -v -bench=. ./test/...

test-all: test-unit test-integration test-coverage test-race ## Run all test types

test-deps: ## Install test dependencies
	@echo -e "$(BLUE)Installing test dependencies...$(NC)"
	@go mod tidy
	@echo -e "$(GREEN)✅ Test dependencies installed$(NC)"

test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)Running tests with verbose output...$(NC)"
	@go test -v -count=1 ./test/...

test-parallel: ## Run tests in parallel
	@echo "$(BLUE)Running tests in parallel...$(NC)"
	@go test -v -parallel 4 ./test/...

test-timeout: ## Run tests with timeout
	@echo "$(BLUE)Running tests with timeout...$(NC)"
	@go test -v -timeout 30s ./test/...

test-badge: ## Generate test coverage badge
	@echo "$(BLUE)Generating test coverage badge...$(NC)"
	@go test -coverprofile=coverage.out ./test/...
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Coverage: $$COVERAGE%"; \
	if [ $$COVERAGE -ge 80 ]; then \
		echo "$(GREEN)✅ Coverage is good ($$COVERAGE%)$(NC)"; \
	else \
		echo "$(RED)❌ Coverage is low ($$COVERAGE%)$(NC)"; \
	fi

# =============================================================================
# CODE QUALITY
# =============================================================================

lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(GREEN)✅ Linting completed$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not found, installing...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		$(GOPATH)/bin/golangci-lint run || $(HOME)/go/bin/golangci-lint run; \
		echo "$(GREEN)✅ Linting completed$(NC)"; \
	fi

lint-fix: ## Run linter with auto-fix
	@echo "$(BLUE)Running linter with auto-fix...$(NC)"
	@golangci-lint run --fix
	@echo "$(GREEN)✅ Linting completed with fixes$(NC)"

format: ## Format code
	@echo -e "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .
	@echo -e "$(GREEN)✅ Code formatted$(NC)"

vet: ## Run go vet
	@echo -e "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo -e "$(GREEN)✅ Go vet completed$(NC)"

security: ## Run security scan
	@echo -e "$(BLUE)Running security scan...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
		echo -e "$(GREEN)✅ Security scan completed$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️  gosec not found, skipping security scan$(NC)"; \
		echo -e "$(GREEN)✅ Security scan skipped$(NC)"; \
	fi

# =============================================================================
# CLEANUP
# =============================================================================

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-*
	@rm -f *.exe
	@go clean
	@echo "$(GREEN)✅ Build artifacts cleaned$(NC)"

clean-test: ## Clean test artifacts
	@echo "$(BLUE)Cleaning test artifacts...$(NC)"
	@rm -f coverage.out coverage.html
	@go clean -testcache
	@echo "$(GREEN)✅ Test artifacts cleaned$(NC)"

clean-all: clean clean-test ## Clean all artifacts

# =============================================================================
# DOCKER
# =============================================================================

docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	@docker build -t mcp-nomad:$(VERSION) .
	@docker tag mcp-nomad:$(VERSION) mcp-nomad:latest
	@echo "$(GREEN)✅ Docker image built: mcp-nomad:$(VERSION)$(NC)"

docker-run: docker-build ## Build and run Docker container
	@echo "$(BLUE)Running Docker container...$(NC)"
	@docker run --rm -p 8080:8080 mcp-nomad:latest

docker-push: docker-build ## Push Docker image to registry
	@echo "$(BLUE)Pushing Docker image...$(NC)"
	@docker push mcp-nomad:$(VERSION)
	@docker push mcp-nomad:latest
	@echo "$(GREEN)✅ Docker image pushed$(NC)"

# =============================================================================
# NOMAD INTEGRATION
# =============================================================================

start-nomad: ## Start Nomad server (requires Docker)
	@echo "$(BLUE)Starting Nomad server...$(NC)"
	@docker run -d --name nomad-server \
		-p 4646:4646 \
		-p 4647:4647 \
		-p 4648:4648 \
		hashicorp/nomad:latest agent -dev
	@echo "$(GREEN)✅ Nomad server started on http://localhost:4646$(NC)"

stop-nomad: ## Stop Nomad server
	@echo "$(BLUE)Stopping Nomad server...$(NC)"
	@docker stop nomad-server || true
	@docker rm nomad-server || true
	@echo "$(GREEN)✅ Nomad server stopped$(NC)"

nomad-status: ## Check Nomad server status
	@echo "$(BLUE)Checking Nomad server status...$(NC)"
	@curl -s http://localhost:4646/v1/status/leader || echo "$(RED)❌ Nomad server not running$(NC)"

# =============================================================================
# RELEASE
# =============================================================================

# release: clean test build-all ## Create release builds
# 	@echo "$(BLUE)Creating release...$(NC)"
# 	@mkdir -p dist
# 	@cp $(BINARY_NAME)-linux dist/
# 	@cp $(BINARY_NAME)-darwin dist/
# 	@cp $(BINARY_NAME)-darwin-arm64 dist/
# 	@cp $(BINARY_NAME)-windows.exe dist/
# 	@cp $(BINARY_NAME)-windows-arm64.exe dist/
# 	@echo "$(GREEN)✅ Release builds created in dist/$(NC)"
#
# release-tar: release ## Create tar archives for release
# 	@echo "$(BLUE)Creating tar archives...$(NC)"
# 	@cd dist && tar -czf mcp-nomad-linux-amd64.tar.gz mcp-nomad-linux
# 	@cd dist && tar -czf mcp-nomad-darwin-amd64.tar.gz mcp-nomad-darwin
# 	@cd dist && tar -czf mcp-nomad-darwin-arm64.tar.gz mcp-nomad-darwin-arm64
# 	@echo "$(GREEN)✅ Tar archives created$(NC)"
#
# release-zip: release ## Create zip archives for release
# 	@echo "$(BLUE)Creating zip archives...$(NC)"
# 	@cd dist && zip mcp-nomad-windows-amd64.zip mcp-nomad-windows.exe
# 	@cd dist && zip mcp-nomad-windows-arm64.zip mcp-nomad-windows-arm64.exe
# 	@echo "$(GREEN)✅ Zip archives created$(NC)"

# =============================================================================
# DEVELOPMENT TOOLS
# =============================================================================

install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@echo "$(YELLOW)Installing air for hot reload...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@echo "$(YELLOW)Installing golangci-lint...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(YELLOW)Installing gosec for security scanning...$(NC)"
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "$(YELLOW)Installing goimports...$(NC)"
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)✅ Development tools installed$(NC)"

# =============================================================================
# STATUS & INFO
# =============================================================================

status: ## Show project status
	@echo -e "$(CYAN)mcp-nomad Status$(NC)"
	@echo -e "$(YELLOW)Version: $(VERSION)$(NC)"
	@echo -e "$(YELLOW)Go Version: $(GO_VERSION)$(NC)"
	@echo -e "$(YELLOW)Build Time: $(BUILD_TIME)$(NC)"
	@echo ""
	@echo -e "$(GREEN)Project Structure:$(NC)"
	@find . -name "*.go" -not -path "./test/*" | wc -l | xargs echo "Go files:"
	@find . -name "*_test.go" | wc -l | xargs echo "Test files:"
	@echo ""
	@echo -e "$(GREEN)Dependencies:$(NC)"
	@go list -m all | wc -l | xargs echo "Total dependencies:"

version: ## Show version information
	@echo -e "$(CYAN)mcp-nomad $(VERSION)$(NC)"
	@echo "Build: $(BUILD_TIME)"
	@echo "Go: $(GO_VERSION)"

# =============================================================================
# QUICK COMMANDS
# =============================================================================

quick-test: ## Quick test run
	@echo "$(BLUE)Running quick tests...$(NC)"
	@go test -short ./test/unit/...

quick-build: ## Quick build without optimizations
	@echo "$(BLUE)Quick building...$(NC)"
	@go build -o $(BINARY_NAME) .
	@echo "$(GREEN)✅ Quick build completed$(NC)"

quick-run: quick-build ## Quick build and run
	@echo "$(BLUE)Quick running...$(NC)"
	@./$(BINARY_NAME)

# =============================================================================
# CI/CD HELPERS
# =============================================================================

ci-test: deps test-deps test ## CI test pipeline
	@echo -e "$(GREEN)✅ CI tests passed$(NC)"

ci-build: deps build ## CI build pipeline
	@echo -e "$(GREEN)✅ CI build completed$(NC)"

ci-lint: deps ## CI lint pipeline
	@echo -e "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo -e "$(GREEN)✅ CI linting passed$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠️  golangci-lint not found, skipping linting$(NC)"; \
		echo -e "$(GREEN)✅ CI linting skipped$(NC)"; \
	fi

ci-all: ci-test ci-build ci-lint vet security ## Run all CI checks
	@echo -e "$(GREEN)✅ All CI checks passed$(NC)"
