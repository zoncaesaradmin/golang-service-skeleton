# Service Makefile
.PHONY: build test clean deep-clean run deps tidy lint fmt vet coverage coverage-enforce coverage-html install-lint docker-build help

# Variables
BINARY_NAME=service
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/main.go
GO_FILES=$(shell find . -name "*.go" -type f)
COVERAGE_THRESHOLD=80
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

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
coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=$(COVERAGE_OUT) ./...
	go tool cover -func=$(COVERAGE_OUT)

# Generate HTML coverage report
coverage-html: coverage
	@echo "Generating HTML coverage report..."
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Enforce coverage threshold
coverage-enforce: coverage
	@echo "Enforcing coverage threshold ($(COVERAGE_THRESHOLD)%)..."
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_OUT) | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	echo "Current coverage: $$COVERAGE%"; \
	echo "Required coverage: $(COVERAGE_THRESHOLD)%"; \
	if [ $$(echo "$$COVERAGE >= $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "✅ Service coverage enforcement passed!"; \
	else \
		echo "❌ Service coverage enforcement failed! Current: $$COVERAGE%, Required: $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	fi

# Run tests with coverage (legacy alias)
test-coverage: coverage

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
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)
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
	@echo "Service Makefile - Available targets:"
	@echo "  build            - Build the binary"
	@echo "  run              - Build and run the service"
	@echo "  test             - Run tests"
	@echo "  coverage         - Run tests with coverage"
	@echo "  coverage-html    - Generate HTML coverage report"
	@echo "  coverage-enforce - Enforce coverage threshold ($(COVERAGE_THRESHOLD)%)"
	@echo "  deps             - Download dependencies"
	@echo "  tidy             - Tidy dependencies"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  lint             - Lint code"
	@echo "  clean            - Clean build artifacts and generated files"
	@echo "  deep-clean       - Deep clean including vendor and module cache"
	@echo "  install-lint     - Install golangci-lint"
	@echo "  docker-build     - Build Docker image"
	@echo "  help             - Show this help"
