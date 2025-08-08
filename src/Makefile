# Root Makefile for Golang Service Skeleton
.PHONY: help build build-with-coverage build-local test test-local clean clean-coverage clean-all-coverage deep-clean run-service run-service-local run-tests run-tests-local install-deps fmt vet lint coverage coverage-enforce start start-local service-build service-build-local service-build-with-coverage service-build-with-coverage-local service-coverage service-coverage-local service-test service-test-local service-run service-run-local testrunner-build testrunner-build-local testrunner-test testrunner-test-local testrunner-run testrunner-run-local shared-build shared-build-local shared-test shared-test-local shared-coverage shared-coverage-local

# Coverage enforcement thresholds
SERVICE_COVERAGE_THRESHOLD = 35

# Build tags (empty by default, but can be set from environment/command line/parent)
BUILD_TAGS ?=

# Helper function to safely append a tag if not already present
define append_tag
$(if $(filter $(1),$(BUILD_TAGS)),$(BUILD_TAGS),$(if $(BUILD_TAGS),$(BUILD_TAGS) $(1),$(1)))
endef

# Local build tags - safely append 'local' to existing BUILD_TAGS
LOCAL_BUILD_TAGS = $(call append_tag,local)

# Default target
help:
	@echo "Golang Service Skeleton - Root Makefile"
	@echo "========================================"
	@echo ""
	@echo "Available targets:"
	@echo "  start           - Build with enforced combined coverage (production mode, default)"
	@echo "  start-local     - Build with enforced combined coverage (local development mode)"
	@echo "  build           - Build all modules (production mode)"
	@echo "  build-local     - Build all modules (local development mode)"
	@echo "  build-with-coverage - Build with coverage enforcement (production mode)"
	@echo "  test            - Run all tests (production mode)"
	@echo "  test-local      - Run all tests (local development mode)"
	@echo "  coverage        - Generate combined service coverage (service + shared)"
	@echo "  coverage-enforce- Enforce service coverage threshold ($(SERVICE_COVERAGE_THRESHOLD)%)"
	@echo "  run-service     - Build with coverage enforcement and run the main service (production)"
	@echo "  run-service-local - Build with coverage enforcement and run the main service (local dev)"
	@echo "  run-service-dev - Build without coverage enforcement and run (for development)"
	@echo "  run-tests       - Build and run integration tests (production)"
	@echo "  run-tests-local - Build and run integration tests (local dev)"
	@echo "  clean           - Clean build artifacts from all projects"
	@echo "  clean-coverage  - Clean individual module coverage files (preserves main coverage.out)"
	@echo "  clean-all-coverage - Clean all coverage files including main combined coverage"
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
	@echo "  📁 coverage.out            - Combined coverage data for entire repository"
	@echo ""
	@echo "Individual project targets:"
	@echo "  service-*            - Run target for service (e.g., make service-build)"
	@echo "  service-build-with-coverage - Build service only with coverage enforcement"
	@echo "  testrunner-*         - Run target for testrunner (e.g., make testrunner-build)"
	@echo "  shared-*             - Run target for shared module (e.g., make shared-test)"

# Start target - Build with enforced combined coverage and run the service
start: coverage-enforce service-build

start-local: coverage-enforce service-build-local
	@echo "✅ Service built with enforced coverage!"

# Build all projects
build: service-build testrunner-build shared-build

build-local: service-build-local testrunner-build-local shared-build-local

# Build all projects with coverage enforcement (recommended for production)
build-with-coverage: coverage-enforce service-build testrunner-build shared-build
	@echo "✅ Build completed with coverage enforcement!"

# Run all tests
test: service-test testrunner-test shared-test

test-local: service-test-local testrunner-test-local shared-test-local

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
	@rm -f coverage.out; \
	echo "mode: set" > coverage.out; \
	if [ -f service/coverage.out ]; then \
		tail -n +2 service/coverage.out >> coverage.out; \
	fi; \
	if [ -f shared/coverage.out ]; then \
		tail -n +2 shared/coverage.out >> coverage.out; \
	fi; \
	if [ ! -s coverage.out ] || [ $$(wc -l < coverage.out) -eq 1 ]; then \
		echo "❌ No coverage data found"; \
		rm -f coverage.out; \
		exit 1; \
	fi; \
	SERVICE_COV=$$(cd service && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SERVICE_COV" ]; then SERVICE_COV="0"; fi; \
	SHARED_COV=$$(cd shared && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SHARED_COV" ]; then SHARED_COV="0"; fi; \
	COMBINED_COV=$$(go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$COMBINED_COV" ]; then \
		echo "⚠️  Warning: Failed to calculate combined coverage, trying alternative method..."; \
		echo "📋 Debug: Coverage file size: $$(wc -l < coverage.out) lines"; \
		COMBINED_COV=$$(go tool cover -func=coverage.out 2>&1 | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
		if [ -z "$$COMBINED_COV" ]; then \
			echo "📋 Debug: Alternative method also failed, setting to 0"; \
			COMBINED_COV="0"; \
		fi; \
	fi; \
	echo "Service Module: $$SERVICE_COV%"; \
	echo "Shared Module:  $$SHARED_COV%"; \
	echo ""; \
	echo "🎯 TRUE COMBINED COVERAGE: $$COMBINED_COV%"; \
	echo "🎯 THRESHOLD: $(SERVICE_COVERAGE_THRESHOLD)%"; \
	echo ""; \
	echo "📁 Coverage files preserved:"; \
	echo "   • service/coverage.out"; \
	echo "   • shared/coverage.out"; \
	echo "   • coverage.out"; \
	echo ""

# Enforce combined service coverage threshold
coverage-enforce: service-coverage shared-coverage
	@echo "=========================================="
	@echo "SERVICE COVERAGE ENFORCEMENT"
	@echo "=========================================="
	@echo "Scope: service/* + shared/* (excludes testrunner/*)"
	@echo "Threshold: $(SERVICE_COVERAGE_THRESHOLD)%"
	@echo ""
	@rm -f coverage.out; \
	echo "mode: set" > coverage.out; \
	if [ -f service/coverage.out ]; then \
		tail -n +2 service/coverage.out >> coverage.out; \
	fi; \
	if [ -f shared/coverage.out ]; then \
		tail -n +2 shared/coverage.out >> coverage.out; \
	fi; \
	if [ ! -s coverage.out ] || [ $$(wc -l < coverage.out) -eq 1 ]; then \
		echo "❌ No coverage data found"; \
		rm -f coverage.out; \
		exit 1; \
	fi; \
	SERVICE_COV=$$(cd service && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SERVICE_COV" ]; then SERVICE_COV="0"; fi; \
	SHARED_COV=$$(cd shared && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SHARED_COV" ]; then SHARED_COV="0"; fi; \
	COMBINED_COV=$$(go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$COMBINED_COV" ]; then \
		echo "⚠️  Warning: Failed to calculate combined coverage, trying alternative method..."; \
		COMBINED_COV=$$(go tool cover -func=coverage.out 2>&1 | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
		if [ -z "$$COMBINED_COV" ]; then COMBINED_COV="0"; fi; \
	fi; \
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
	echo "   • coverage.out"

# Clean all projects
clean: service-clean testrunner-clean shared-clean
	@echo "Cleaning individual module coverage files while preserving main coverage.out..."
	@rm -f service/coverage.out shared/coverage.out 2>/dev/null || true
	@echo "✅ All projects cleaned (main coverage.out preserved)"

# Clean coverage files (when needed for fresh coverage run)
clean-coverage:
	@echo "Cleaning individual module coverage files..."
	@rm -f service/coverage.out shared/coverage.out
	@echo "✅ Individual coverage files cleaned (main coverage.out preserved)"

# Clean all coverage files including main combined coverage
clean-all-coverage:
	@echo "Cleaning all coverage files including main combined coverage..."
	@rm -f service/coverage.out shared/coverage.out coverage.out
	@echo "✅ All coverage files cleaned"

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

run-service-local: service-build-with-coverage-local service-run-local

# Run the main service (without coverage enforcement - for development)
run-service-dev: service-build service-run

# Run integration tests
run-tests: testrunner-run

run-tests-local: testrunner-run-local

# Service targets
service-build:
	@echo "Building service..."
	@cd service && $(MAKE) build BUILD_TAGS="$(BUILD_TAGS)"

service-build-local:
	@echo "Building service (local development mode)..."
	@cd service && $(MAKE) build BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

service-build-with-coverage: service-coverage shared-coverage
	@echo "Building service with coverage enforcement..."
	@SERVICE_COV=$$(cd service && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SERVICE_COV" ]; then SERVICE_COV="0"; fi; \
	SHARED_COV=$$(cd shared && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SHARED_COV" ]; then SHARED_COV="0"; fi; \
	rm -f coverage.out; \
	echo "mode: set" > coverage.out; \
	if [ -f service/coverage.out ]; then \
		tail -n +2 service/coverage.out >> coverage.out; \
	fi; \
	if [ -f shared/coverage.out ]; then \
		tail -n +2 shared/coverage.out >> coverage.out; \
	fi; \
	COMBINED_COV=$$(go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$COMBINED_COV" ]; then \
		echo "⚠️  Warning: Failed to calculate combined coverage, trying alternative method..."; \
		COMBINED_COV=$$(go tool cover -func=coverage.out 2>&1 | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
		if [ -z "$$COMBINED_COV" ]; then COMBINED_COV="0"; fi; \
	fi; \
	echo "🎯 Service Coverage: $$SERVICE_COV% | Shared: $$SHARED_COV% | Combined: $$COMBINED_COV%"; \
	if [ $$(echo "$$COMBINED_COV >= $(SERVICE_COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "✅ Coverage $$COMBINED_COV% meets $(SERVICE_COVERAGE_THRESHOLD)% threshold - proceeding with build"; \
		echo "📁 Coverage files preserved for org tools:"; \
		echo "   • service/coverage.out"; \
		echo "   • shared/coverage.out"; \
		echo "   • coverage.out"; \
		cd service && $(MAKE) build BUILD_TAGS="$(BUILD_TAGS)"; \
	else \
		echo "❌ Coverage $$COMBINED_COV% below $(SERVICE_COVERAGE_THRESHOLD)% threshold - build blocked"; \
		exit 1; \
	fi

service-build-with-coverage-local: service-coverage-local shared-coverage-local
	@echo "Building service with coverage enforcement (local development mode)..."
	@SERVICE_COV=$$(cd service && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SERVICE_COV" ]; then SERVICE_COV="0"; fi; \
	SHARED_COV=$$(cd shared && go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SHARED_COV" ]; then SHARED_COV="0"; fi; \
	rm -f coverage.out; \
	echo "mode: set" > coverage.out; \
	if [ -f service/coverage.out ]; then \
		tail -n +2 service/coverage.out >> coverage.out; \
	fi; \
	if [ -f shared/coverage.out ]; then \
		tail -n +2 shared/coverage.out >> coverage.out; \
	fi; \
	COMBINED_COV=$$(go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$COMBINED_COV" ]; then \
		echo "⚠️  Warning: Failed to calculate combined coverage, trying alternative method..."; \
		COMBINED_COV=$$(go tool cover -func=coverage.out 2>&1 | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
		if [ -z "$$COMBINED_COV" ]; then COMBINED_COV="0"; fi; \
	fi; \
	echo "🎯 Service Coverage (local): $$SERVICE_COV% | Shared: $$SHARED_COV% | Combined: $$COMBINED_COV%"; \
	if [ $$(echo "$$COMBINED_COV >= $(SERVICE_COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "✅ Coverage $$COMBINED_COV% meets $(SERVICE_COVERAGE_THRESHOLD)% threshold - proceeding with build (local)"; \
		echo "📁 Coverage files preserved for org tools:"; \
		echo "   • service/coverage.out"; \
		echo "   • shared/coverage.out"; \
		echo "   • coverage.out"; \
		cd service && $(MAKE) build BUILD_TAGS="$(LOCAL_BUILD_TAGS)"; \
	else \
		echo "❌ Coverage $$COMBINED_COV% below $(SERVICE_COVERAGE_THRESHOLD)% threshold - build blocked"; \
		exit 1; \
	fi

service-coverage:
	@echo "Running service coverage..."
	@cd service && $(MAKE) coverage BUILD_TAGS="$(BUILD_TAGS)"

service-coverage-local:
	@echo "Running service coverage (local development mode)..."
	@cd service && $(MAKE) coverage BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

service-test:
	@echo "Running service tests..."
	@cd service && $(MAKE) test BUILD_TAGS="$(BUILD_TAGS)"

service-test-local:
	@echo "Running service tests (local development mode)..."
	@cd service && $(MAKE) test BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

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
	@cd service && $(MAKE) run BUILD_TAGS="$(BUILD_TAGS)"

service-run-local:
	@echo "Running service (local development mode)..."
	@cd service && $(MAKE) run BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

# Testrunner targets
testrunner-build:
	@echo "Building testrunner..."
	@cd testrunner && $(MAKE) build BUILD_TAGS="$(BUILD_TAGS)"

testrunner-build-local:
	@echo "Building testrunner (local development mode)..."
	@cd testrunner && $(MAKE) build BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

testrunner-test:
	@echo "Running testrunner unit tests..."
	@cd testrunner && $(MAKE) test BUILD_TAGS="$(BUILD_TAGS)"

testrunner-test-local:
	@echo "Running testrunner unit tests (local development mode)..."
	@cd testrunner && $(MAKE) test BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

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
	@cd testrunner && $(MAKE) run BUILD_TAGS="$(BUILD_TAGS)"

testrunner-run-local:
	@echo "Running testrunner (local development mode)..."
	@cd testrunner && $(MAKE) run BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

testrunner-integration:
	@echo "Running integration tests..."
	@cd testrunner && $(MAKE) run-integration-tests

testrunner-performance:
	@echo "Running performance tests..."
	@cd testrunner && $(MAKE) run-performance-tests

# Shared module targets
shared-build:
	@echo "Building shared module..."
	@cd shared && $(MAKE) build BUILD_TAGS="$(BUILD_TAGS)"

shared-build-local:
	@echo "Building shared module (local development mode)..."
	@cd shared && $(MAKE) build BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

shared-test:
	@echo "Running shared module tests..."
	@cd shared && $(MAKE) test BUILD_TAGS="$(BUILD_TAGS)"

shared-test-local:
	@echo "Running shared module tests (local development mode)..."
	@cd shared && $(MAKE) test BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

shared-coverage:
	@echo "Running shared module coverage..."
	@cd shared && $(MAKE) coverage BUILD_TAGS="$(BUILD_TAGS)"

shared-coverage-local:
	@echo "Running shared module coverage (local development mode)..."
	@cd shared && $(MAKE) coverage BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

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
