# Service Makefile
.PHONY: build test clean deep-clean run deps tidy lint fmt vet install-lint docker-build help

# Variables
BINARY_NAME=service
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/main.go
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
all: clean deps fmt vet lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build completed: $(BINARY_PATH)"

# Run the service
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BINARY_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

# Clean build artifacts and generated files
clean:
	@echo "Cleaning build artifacts and generated files..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f *.log
	@rm -rf tmp/ temp/
	@find . -name "*.test" -type f -delete
	@find . -name "*.out" -type f -delete
	@echo "Clean completed"

# Deep clean - removes all generated files including vendor and module cache
deep-clean: clean
	@echo "Deep cleaning..."
	@rm -rf vendor/
	@go clean -modcache
	@echo "Deep clean completed"

# Install golangci-lint (optional)
install-lint:
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker build (optional)
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  run            - Build and run the service"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy dependencies"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  lint           - Lint code"
	@echo "  clean          - Clean build artifacts and generated files"
	@echo "  deep-clean     - Deep clean including vendor and module cache"
	@echo "  install-lint   - Install golangci-lint"
	@echo "  docker-build   - Build Docker image"
	@echo "  help           - Show this help"
