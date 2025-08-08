# Root Makefile for Golang Service Skeleton
.PHONY: help build build-with-coverage test clean clean-coverage deep-clean run-service run-tests install-deps fmt vet lint coverage coverage-enforce

# Coverage enforcement thresholds
SERVICE_COVERAGE_THRESHOLD=25

# Default target
help:
	@echo "Golang Service Skeleton - Root Makefile"
	@echo "========================================"
	@echo ""
	@echo "Available targets:"
	@echo "  build           - Build all modules (service, testrunner, shared)"
	@echo "  build-with-coverage - Build with coverage enforcement (recommended for production)"
	@echo "  test            - Run all tests (unit and integration)"
	@echo "  coverage        - Generate combined service coverage (service + shared)"
	@echo "  coverage-enforce- Enforce service coverage threshold ($(SERVICE_COVERAGE_THRESHOLD)%)"
	@echo "  run-service     - Build with coverage enforcement and run the main service"
	@echo "  run-service-dev - Build without coverage enforcement and run (for development)"
	@echo "  run-tests       - Build and run integration tests"
	@echo "  clean           - Clean build artifacts from all projects"
	@echo "  clean-coverage  - Clean coverage files (service/shared/combined coverage.out)"
	@echo "  deep-clean      - Deep clean including vendor and module cache"
	@echo "  install-deps    - Download dependencies for all projects"
	@echo "  fmt             - Format code in all projects"
	@echo "  vet             - Vet code in all projects"
	@echo "  lint            - Lint code in all projects"
	@echo "  ci              - Full CI pipeline with coverage enforcement"
	@echo "  help            - Show this help"
	@echo ""
	@echo "Service Coverage Scope:"
	@echo "  ✅ service/*     - Main business logic (included)"
	@echo "  ✅ shared/*      - Common utilities (included)"
	@echo "  ❌ testrunner/*  - Testing infrastructure (excluded)"
	@echo ""
	@echo "Coverage Files (preserved for org tools):"
	@echo "  📁 service/coverage.out     - Service module coverage data"
	@echo "  📁 shared/coverage.out      - Shared module coverage data"
	@echo "  📁 combined_coverage.out    - Combined coverage data (root level)"
	@echo ""
	@echo "Individual project targets:"
	@echo "  service-*            - Run target for service (e.g., make service-build)"
	@echo "  service-build-with-coverage - Build service only with coverage enforcement"
	@echo "  testrunner-*         - Run target for testrunner (e.g., make testrunner-build)"
	@echo "  shared-*             - Run target for shared module (e.g., make shared-test)"

# Build all projects
build: service-build testrunner-build shared-build

# Build all projects with coverage enforcement (recommended for production)
build-with-coverage: coverage-enforce service-build testrunner-build shared-build
	@echo "✅ Build completed with coverage enforcement!"

# Run all tests
test: service-test testrunner-test shared-test

# Generate combined service coverage (service + shared modules)
coverage: service-coverage shared-coverage
	@echo ""
	@echo "=========================================="
	@echo "COMBINED SERVICE COVERAGE REPORT"
	@echo "=========================================="
	@echo "Scope: service/* + shared/* (excludes testrunner/*)"
	@echo ""
	@echo "📊 Service Module Coverage:"
	@cd service && go tool cover -func=coverage.out | tail -1 || echo "No service coverage file"
	@echo ""
	@echo "📊 Shared Module Coverage:"
	@cd shared && go tool cover -func=coverage.out | tail -1 || echo "No shared coverage file"
	@echo ""
	@echo "Calculating true combined coverage..."
	@$(MAKE) --no-print-directory coverage-calculate

# Calculate combined coverage from service and shared modules
coverage-calculate:
	@echo "🔢 True Combined Coverage Calculation:"
	@rm -f combined_coverage.out; \
	echo "mode: set" > combined_coverage.out; \
	if [ -f service/coverage.out ]; then \
		tail -n +2 service/coverage.out >> combined_coverage.out; \
	fi; \
	if [ -f shared/coverage.out ]; then \
		tail -n +2 shared/coverage.out >> combined_coverage.out; \
	fi; \
	if [ ! -s combined_coverage.out ] || [ $$(wc -l < combined_coverage.out) -eq 1 ]; then \
		echo "❌ No coverage data found"; \
		rm -f combined_coverage.out; \
		exit 1; \
	fi; \
	SERVICE_COV=$$(cd service && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	SHARED_COV=$$(cd shared && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	COMBINED_COV=$$(go tool cover -func=combined_coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	echo "Service Module: $$SERVICE_COV%"; \
	echo "Shared Module:  $$SHARED_COV%"; \
	echo ""; \
	echo "🎯 TRUE COMBINED COVERAGE: $$COMBINED_COV%"; \
	echo "🎯 THRESHOLD: $(SERVICE_COVERAGE_THRESHOLD)%"; \
	echo ""; \
	echo "📁 Coverage files preserved:"; \
	echo "   • service/coverage.out"; \
	echo "   • shared/coverage.out"; \
	echo "   • combined_coverage.out"; \
	echo ""

# Enforce combined service coverage threshold
coverage-enforce: service-coverage shared-coverage
	@echo "=========================================="
	@echo "SERVICE COVERAGE ENFORCEMENT"
	@echo "=========================================="
	@echo "Scope: service/* + shared/* (excludes testrunner/*)"
	@echo "Threshold: $(SERVICE_COVERAGE_THRESHOLD)%"
	@echo ""
	@rm -f combined_coverage.out; \
	echo "mode: set" > combined_coverage.out; \
	if [ -f service/coverage.out ]; then \
		tail -n +2 service/coverage.out >> combined_coverage.out; \
	fi; \
	if [ -f shared/coverage.out ]; then \
		tail -n +2 shared/coverage.out >> combined_coverage.out; \
	fi; \
	if [ ! -s combined_coverage.out ] || [ $$(wc -l < combined_coverage.out) -eq 1 ]; then \
		echo "❌ No coverage data found"; \
		rm -f combined_coverage.out; \
		exit 1; \
	fi; \
	SERVICE_COV=$$(cd service && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	SHARED_COV=$$(cd shared && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	COMBINED_COV=$$(go tool cover -func=combined_coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	echo "📊 Service Module: $$SERVICE_COV%"; \
	echo "📊 Shared Module:  $$SHARED_COV%"; \
	echo ""; \
	echo "🎯 TRUE COMBINED COVERAGE: $$COMBINED_COV%"; \
	echo "🎯 REQUIRED: $(SERVICE_COVERAGE_THRESHOLD)%"; \
	echo ""; \
	if [ $$(echo "$$COMBINED_COV >= $(SERVICE_COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "✅ SERVICE COVERAGE ENFORCEMENT PASSED!"; \
		echo "   True combined coverage $$COMBINED_COV% meets $(SERVICE_COVERAGE_THRESHOLD)% threshold"; \
	else \
		echo "❌ SERVICE COVERAGE ENFORCEMENT FAILED!"; \
		echo "   True combined coverage $$COMBINED_COV% below $(SERVICE_COVERAGE_THRESHOLD)% threshold"; \
		echo ""; \
		echo "💡 Coverage includes:"; \
		echo "   ✅ service/* (business logic)"; \
		echo "   ✅ shared/* (common utilities)"; \
		echo "   ❌ testrunner/* (excluded - testing infrastructure)"; \
		exit 1; \
	fi; \
	echo "📁 Coverage files preserved:"; \
	echo "   • service/coverage.out"; \
	echo "   • shared/coverage.out"; \
	echo "   • combined_coverage.out"

# Clean all projects
clean: service-clean testrunner-clean shared-clean

# Clean coverage files (when needed for fresh coverage run)
clean-coverage:
	@echo "Cleaning coverage files..."
	@rm -f service/coverage.out shared/coverage.out combined_coverage.out
	@echo "✅ Coverage files cleaned"

# Deep clean all projects
deep-clean: service-deep-clean testrunner-deep-clean shared-deep-clean

# Install dependencies for all projects
install-deps: service-deps testrunner-deps shared-deps

# Format code in all projects
fmt: service-fmt testrunner-fmt shared-fmt

# Vet code in all projects
vet: service-vet testrunner-vet shared-vet

# Lint code in all projects
lint: service-lint shared-lint

# Run the main service (with coverage enforcement)
run-service: service-build-with-coverage service-run

# Run the main service (without coverage enforcement - for development)
run-service-dev: service-build service-run

# Run integration tests
run-tests: testrunner-run

# Service targets
service-build:
	@echo "Building service..."
	@cd service && $(MAKE) build

service-build-with-coverage: service-coverage shared-coverage
	@echo "Building service with coverage enforcement..."
	@SERVICE_COV=$$(cd service && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	SHARED_COV=$$(cd shared && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	rm -f combined_coverage.out; \
	echo "mode: set" > combined_coverage.out; \
	if [ -f service/coverage.out ]; then \
		tail -n +2 service/coverage.out >> combined_coverage.out; \
	fi; \
	if [ -f shared/coverage.out ]; then \
		tail -n +2 shared/coverage.out >> combined_coverage.out; \
	fi; \
	COMBINED_COV=$$(go tool cover -func=combined_coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//' || echo "0"); \
	echo "🎯 Service Coverage: $$SERVICE_COV% | Shared: $$SHARED_COV% | Combined: $$COMBINED_COV%"; \
	if [ $$(echo "$$COMBINED_COV >= $(SERVICE_COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "✅ Coverage $$COMBINED_COV% meets $(SERVICE_COVERAGE_THRESHOLD)% threshold - proceeding with build"; \
		echo "📁 Coverage files preserved for org tools:"; \
		echo "   • service/coverage.out"; \
		echo "   • shared/coverage.out"; \
		echo "   • combined_coverage.out"; \
		cd service && $(MAKE) build; \
	else \
		echo "❌ Coverage $$COMBINED_COV% below $(SERVICE_COVERAGE_THRESHOLD)% threshold - build blocked"; \
		exit 1; \
	fi

service-coverage:
	@echo "Running service coverage..."
	@cd service && $(MAKE) coverage

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

# Shared module targets
shared-build:
	@echo "Building shared module..."
	@cd shared && $(MAKE) build

shared-test:
	@echo "Running shared module tests..."
	@cd shared && $(MAKE) test

shared-coverage:
	@echo "Running shared module coverage..."
	@cd shared && $(MAKE) coverage

shared-coverage-enforce:
	@echo "Enforcing shared module coverage..."
	@cd shared && $(MAKE) coverage-enforce

shared-clean:
	@echo "Cleaning shared module..."
	@cd shared && $(MAKE) clean

shared-deep-clean:
	@echo "Deep cleaning shared module..."
	@cd shared && $(MAKE) clean # shared module doesn't have vendor/

shared-deps:
	@echo "Installing shared module dependencies..."
	@cd shared && $(MAKE) deps

shared-fmt:
	@echo "Formatting shared module code..."
	@cd shared && $(MAKE) fmt

shared-vet:
	@echo "Vetting shared module code..."
	@cd shared && $(MAKE) vet

shared-lint:
	@echo "Linting shared module code..."
	@cd shared && $(MAKE) lint

shared-benchmark:
	@echo "Running shared module benchmarks..."
	@cd shared && $(MAKE) benchmark

shared-race:
	@echo "Running shared module race detection..."
	@cd shared && $(MAKE) race

# Quick development workflow
dev-setup: install-deps fmt vet build
	@echo "Development setup completed!"

# CI/CD workflow with coverage enforcement
ci: install-deps fmt vet test coverage-enforce build
	@echo "✅ CI pipeline completed successfully!"

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
