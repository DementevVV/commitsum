# Makefile for commitsum
# Inspired by modern Go project conventions

.PHONY: help build run clean deps fmt lint install uninstall build-all release-tag

# Default target
.DEFAULT_GOAL := help

# Application name
APP_NAME := commitsum
MAIN_PKG := ./cmd/commitsum
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Build targets for multiple platforms
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

help: ## Show this help message
	@echo '📊 GitHub Commit Summarizer - Available commands:'
	@echo ''
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ''

build: ## Build the application
	@echo "🔨 Building $(APP_NAME)..."
	@go build $(LDFLAGS) -o $(APP_NAME) $(MAIN_PKG)
	@echo "✅ Build complete: ./$(APP_NAME)"

build-all: ## Build for all platforms (Linux, macOS, Windows)
	@echo "🔨 Building for all platforms..."
	@mkdir -p bin
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		out="bin/$(APP_NAME)-$${platform%/*}-$${platform#*/}"; \
		if [ "$$GOOS" = "windows" ]; then out="$$out.exe"; fi; \
		go build $(LDFLAGS) -o $$out $(MAIN_PKG) ; \
		echo "✅ Built for $$platform"; \
	done
	@echo "✨ All builds complete in ./bin/"

run: build ## Build and run the application
	@echo "🚀 Running $(APP_NAME)..."
	@./$(APP_NAME)

deps: ## Install/update dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

fmt: ## Format code with gofmt
	@echo "🎨 Formatting code..."
	@gofmt -s -w .
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@if [ ! -f ./bin/golangci-lint ]; then \
		echo "📦 Installing golangci-lint to ./bin..."; \
		curl -sSfL https://golangci-lint.run/install.sh | sh -s v2.8.0; \
	fi
	./bin/golangci-lint run

clean: ## Remove build artifacts
	@echo "🧹 Cleaning..."
	@rm -f $(APP_NAME)
	@rm -rf bin/
	@echo "✅ Clean complete"

install: build ## Install to /usr/local/bin
	@echo "📥 Installing $(APP_NAME) to /usr/local/bin..."
	@sudo cp $(APP_NAME) /usr/local/bin/
	@echo "✅ Installed successfully"

uninstall: ## Remove from /usr/local/bin
	@echo "📤 Uninstalling $(APP_NAME)..."
	@sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "✅ Uninstalled successfully"

dev: ## Run with auto-reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

release: ## Create and push a release tag (triggers GitHub Actions)
	@read -p "Release tag (vX.Y.Z): " tag; \
	if ! echo $$tag | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' > /dev/null; then \
		echo "❌ Invalid tag format. Use vX.Y.Z (e.g., v1.2.3)"; \
		exit 1; \
	fi; \
	if [ ! -f CHANGELOG.md ]; then \
		echo "❌ CHANGELOG.md not found. Create it and add a section for $$tag."; \
		exit 1; \
	fi; \
	if ! grep -Eq "^(##[[:space:]]*)?$$tag$$" CHANGELOG.md; then \
		echo "❌ CHANGELOG.md missing section for $$tag. Add a heading like '## $$tag'."; \
		exit 1; \
	fi; \
	git tag -a $$tag -m "Release $$tag"; \
	git push origin $$tag; \
	echo "✅ Pushed tag $$tag. GitHub Actions will build and publish the release."
