# Root Makefile for Golang Service Skeleton
.PHONY: help build test clean deep-clean run-service run-tests install-deps fmt vet lint

# Default target
help:
	@echo "Golang Service Skeleton - Root Makefile"
	@echo "========================================"
	@echo ""
	@echo "Available targets:"
	@echo "  build         - Build both service and testrunner binaries"
	@echo "  test          - Run all tests (unit and integration)"
	@echo "  run-service   - Build and run the main service"
	@echo "  run-tests     - Build and run integration tests"
	@echo "  clean         - Clean build artifacts from both projects"
	@echo "  deep-clean    - Deep clean including vendor and module cache"
	@echo "  install-deps  - Download dependencies for both projects"
	@echo "  fmt           - Format code in both projects"
	@echo "  vet           - Vet code in both projects"
	@echo "  lint          - Lint code in both projects"
	@echo "  help          - Show this help"
	@echo ""
	@echo "Individual project targets:"
	@echo "  service-*     - Run target for service (e.g., make service-build)"
	@echo "  testrunner-*  - Run target for testrunner (e.g., make testrunner-build)"

# Build both projects
build: service-build testrunner-build

# Run all tests
test: service-test testrunner-test

# Clean both projects
clean: service-clean testrunner-clean

# Deep clean both projects
deep-clean: service-deep-clean testrunner-deep-clean

# Install dependencies for both projects
install-deps: service-deps testrunner-deps

# Format code in both projects
fmt: service-fmt testrunner-fmt

# Vet code in both projects
vet: service-vet testrunner-vet

# Lint code in both projects
lint: service-lint

# Run the main service
run-service: service-run

# Run integration tests
run-tests: testrunner-run

# Service targets
service-build:
	@echo "Building service..."
	@cd service && $(MAKE) build

service-test:
	@echo "Running service tests..."
	@cd service && $(MAKE) test

service-clean:
	@echo "Cleaning service..."
	@cd service && $(MAKE) clean

service-deep-clean:
	@echo "Deep cleaning service..."
	@cd service && $(MAKE) deep-clean

service-deps:
	@echo "Installing service dependencies..."
	@cd service && $(MAKE) deps

service-fmt:
	@echo "Formatting service code..."
	@cd service && $(MAKE) fmt

service-vet:
	@echo "Vetting service code..."
	@cd service && $(MAKE) vet

service-lint:
	@echo "Linting service code..."
	@cd service && $(MAKE) lint

service-run:
	@echo "Running service..."
	@cd service && $(MAKE) run

# Testrunner targets
testrunner-build:
	@echo "Building testrunner..."
	@cd testrunner && $(MAKE) build

testrunner-test:
	@echo "Running testrunner unit tests..."
	@cd testrunner && $(MAKE) test

testrunner-clean:
	@echo "Cleaning testrunner..."
	@cd testrunner && $(MAKE) clean

testrunner-deep-clean:
	@echo "Deep cleaning testrunner..."
	@cd testrunner && $(MAKE) deep-clean

testrunner-deps:
	@echo "Installing testrunner dependencies..."
	@cd testrunner && $(MAKE) deps

testrunner-fmt:
	@echo "Formatting testrunner code..."
	@cd testrunner && $(MAKE) fmt

testrunner-vet:
	@echo "Vetting testrunner code..."
	@cd testrunner && $(MAKE) vet

testrunner-run:
	@echo "Running testrunner..."
	@cd testrunner && $(MAKE) run

testrunner-integration:
	@echo "Running integration tests..."
	@cd testrunner && $(MAKE) run-integration-tests

testrunner-performance:
	@echo "Running performance tests..."
	@cd testrunner && $(MAKE) run-performance-tests

# Quick development workflow
dev-setup: install-deps fmt vet build
	@echo "Development setup completed!"

# CI/CD workflow
ci: install-deps fmt vet lint test build
	@echo "CI pipeline completed!"

# Full end-to-end test (requires service to be running)
e2e-test: service-build testrunner-build
	@echo "Starting end-to-end test..."
	@echo "Starting service in background..."
	@cd service && ./bin/service & SERVICE_PID=$$!; \
	sleep 3; \
	echo "Running integration tests..."; \
	cd testrunner && ./bin/testrunner; \
	TEST_RESULT=$$?; \
	echo "Stopping service..."; \
	kill $$SERVICE_PID 2>/dev/null || true; \
	exit $$TEST_RESULT
