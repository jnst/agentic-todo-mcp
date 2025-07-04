# Makefile for agentic-todo-mcp

.PHONY: all build test lint fmt fix clean deps generate

# Default target
all: fmt lint test build

# Build the project
build:
	@echo "Building agentic-todo-mcp..."
	go build -o bin/agentic-todo-mcp cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Auto-fix common issues
fix:
	@echo "Auto-fixing common issues..."
	go fmt ./...
	goimports -w .
	golangci-lint run --fix

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Generate mocks
generate:
	@echo "Generating mocks..."
	go generate ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install go.uber.org/mock/mockgen@latest

# Run development server
run:
	@echo "Running development server..."
	go run cmd/server/main.go

# Development workflow (format, lint, test)
dev: fmt lint test

# CI workflow
ci: deps fmt lint test-race

# Help
help:
	@echo "Available targets:"
	@echo "  all             - Format, lint, test, and build"
	@echo "  build           - Build the project"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-race       - Run tests with race detection"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  fix             - Auto-fix common issues"
	@echo "  clean           - Clean build artifacts"
	@echo "  deps            - Download dependencies"
	@echo "  generate        - Generate mocks"
	@echo "  install-tools   - Install development tools"
	@echo "  run             - Run development server"
	@echo "  dev             - Development workflow (format, lint, test)"
	@echo "  ci              - CI workflow"
	@echo "  help            - Show this help"