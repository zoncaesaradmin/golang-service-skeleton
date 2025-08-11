# Root Makefile for Golang Component Skeleton
.PHONY: help build build-with-coverage build-local test test-local clean clean-coverage clean-all-coverage deep-clean run-component run-component-local run-tests run-tests-local install-deps fmt vet lint coverage coverage-enforce start start-local component-build component-build-local component-build-with-coverage component-build-with-coverage-local component-coverage component-coverage-local component-test component-test-local component-run component-run-local testrunner-build testrunner-build-local testrunner-test testrunner-test-local testrunner-run testrunner-run-local shared-build shared-build-local shared-test shared-test-local shared-coverage shared-coverage-local

# Coverage enforcement thresholds
SERVICE_COVERAGE_THRESHOLD = 60

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
	@echo "Golang Component Skeleton - Root Makefile"
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
	@echo "  coverage        - Generate combined component coverage (component + shared)"
	@echo "  coverage-enforce- Enforce component coverage threshold ($(SERVICE_COVERAGE_THRESHOLD)%)"
	@echo "  run-component     - Build with coverage enforcement and run the main component (production)"
	@echo "  run-component-local - Build with coverage enforcement and run the main component (local dev)"
	@echo "  run-component-dev - Build without coverage enforcement and run (for development)"
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
	@echo "  ✅ component/*     - Main business logic (included)"
	@echo "  ✅ shared/*      - Common utilities (included)"
	@echo "  ❌ testrunner/*  - Testing infrastructure (excluded)"
	@echo ""
	@echo "Coverage Files (preserved for org tools):"
	@echo "  📁 coverage.out            - Combined coverage data for entire repository"
	@echo ""
	@echo "Individual project targets:"
	@echo "  component-*            - Run target for component (e.g., make component-build)"
	@echo "  component-build-with-coverage - Build component only with coverage enforcement"
	@echo "  testrunner-*         - Run target for testrunner (e.g., make testrunner-build)"
	@echo "  shared-*             - Run target for shared module (e.g., make shared-test)"

# Start target - Build with enforced combined coverage and run the component
start: coverage-enforce component-build

start-local: coverage-enforce component-build-local
	@echo "✅ Service built with enforced coverage!"

# Build all projects
build: component-build testrunner-build shared-build

build-local: component-build-local testrunner-build-local shared-build-local

# Build all projects with coverage enforcement (recommended for production)
build-with-coverage: coverage-enforce component-build testrunner-build shared-build
	@echo "✅ Build completed with coverage enforcement!"

# Run all tests
test: component-test testrunner-test shared-test

test-local: component-test-local testrunner-test-local shared-test-local

# Generate combined component coverage (component + shared modules)
coverage: component-coverage shared-coverage
	@echo ""
	@echo "=========================================="
	@echo "COMBINED SERVICE COVERAGE REPORT"
	@echo "=========================================="
	@echo "Scope: component/* + shared/* (excludes testrunner/*)"
	@echo ""
	@echo "📊 Service Module Coverage:"
	@cd component && go tool cover -func=component_coverage.tmp | tail -1 || echo "No component coverage file"
	@echo ""
	@echo "📊 Shared Module Coverage:"
	@cd shared && go tool cover -func=shared_coverage.tmp | tail -1 || echo "No shared coverage file"
	@echo ""
	@echo "Calculating true combined coverage..."
	@$(MAKE) --no-print-directory coverage-calculate

# Calculate combined coverage from component and shared modules
coverage-calculate:
	@echo "🔢 True Combined Coverage Calculation:"
	@rm -f coverage.out; \
	echo "mode: set" > coverage.out; \
	if [ -f component/component_coverage.tmp ]; then \
		tail -n +2 component/component_coverage.tmp >> coverage.out; \
	fi; \
	if [ -f shared/shared_coverage.tmp ]; then \
		tail -n +2 shared/shared_coverage.tmp >> coverage.out; \
	fi; \
	if [ ! -s coverage.out ] || [ $$(wc -l < coverage.out) -eq 1 ]; then \
		echo "❌ No coverage data found"; \
		rm -f coverage.out; \
		exit 1; \
	fi; \
	SERVICE_COV=$$(cd component && go tool cover -func=component_coverage.tmp 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SERVICE_COV" ]; then SERVICE_COV="0"; fi; \
	SHARED_COV=$$(cd shared && go tool cover -func=shared_coverage.tmp 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
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
	echo "   • component/component_coverage.tmp"; \
	echo "   • shared/shared_coverage.tmp"; \
	echo "   • coverage.out"; \
	echo ""

# Enforce combined component coverage threshold
coverage-enforce: component-coverage shared-coverage
	@echo "=========================================="
	@echo "SERVICE COVERAGE ENFORCEMENT"
	@echo "=========================================="
	@echo "Scope: component/* + shared/* (excludes testrunner/*)"
	@echo "Threshold: $(SERVICE_COVERAGE_THRESHOLD)%"
	@echo ""
	@rm -f coverage.out; \
	echo "mode: set" > coverage.out; \
	if [ -f component/component_coverage.tmp ]; then \
		tail -n +2 component/component_coverage.tmp >> coverage.out; \
	fi; \
	if [ -f shared/shared_coverage.tmp ]; then \
		tail -n +2 shared/shared_coverage.tmp >> coverage.out; \
	fi; \
	if [ ! -s coverage.out ] || [ $$(wc -l < coverage.out) -eq 1 ]; then \
		echo "❌ No coverage data found"; \
		rm -f coverage.out; \
		exit 1; \
	fi; \
	SERVICE_COV=$$(cd component && go tool cover -func=component_coverage.tmp 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SERVICE_COV" ]; then SERVICE_COV="0"; fi; \
	SHARED_COV=$$(cd shared && go tool cover -func=shared_coverage.tmp 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
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
		echo "   ✅ component/* (business logic)"; \
		echo "   ✅ shared/* (common utilities)"; \
		echo "   ❌ testrunner/* (excluded - testing infrastructure)"; \
		exit 1; \
	fi; \
	echo "📁 Coverage files preserved:"; \
	echo "   • component/coverage.out"; \
	echo "   • shared/coverage.out"; \
	echo "   • coverage.out"

# Clean all projects
clean: component-clean testrunner-clean shared-clean
	@echo "Cleaning individual module coverage files while preserving main coverage.out..."
	@rm -f component/component_coverage.tmp shared/shared_coverage.tmp 2>/dev/null || true
	@echo "✅ All projects cleaned (main coverage.out preserved)"

# Clean coverage files (when needed for fresh coverage run)
clean-coverage:
	@echo "Cleaning individual module coverage files..."
	@rm -f component/component_coverage.tmp shared/shared_coverage.tmp
	@echo "✅ Individual coverage files cleaned (main coverage.out preserved)"

# Clean all coverage files including main combined coverage
clean-all-coverage:
	@echo "Cleaning all coverage files including main combined coverage..."
	@rm -f component/component_coverage.tmp shared/shared_coverage.tmp coverage.out
	@echo "✅ All coverage files cleaned"

# Deep clean all projects
deep-clean: component-deep-clean testrunner-deep-clean shared-deep-clean

# Install dependencies for all projects
install-deps: component-deps testrunner-deps shared-deps

# Format code in all projects
fmt: component-fmt testrunner-fmt shared-fmt

# Vet code in all projects
vet: component-vet testrunner-vet shared-vet

# Lint code in all projects
lint: component-lint shared-lint

# Run the main component (with coverage enforcement)
run-component: component-build-with-coverage component-run

run-component-local: component-build-with-coverage-local component-run-local

# Run the main component (without coverage enforcement - for development)
run-component-dev: component-build component-run

# Run integration tests
run-tests: testrunner-run

run-tests-local: testrunner-run-local

# Service targets
component-build:
	@echo "Building component..."
	@cd component && $(MAKE) build BUILD_TAGS="$(BUILD_TAGS)"

component-build-local:
	@echo "Building component (local development mode)..."
	@cd component && $(MAKE) build BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

component-build-with-coverage: component-coverage shared-coverage
	@echo "Building component with coverage enforcement..."
	@SERVICE_COV=$$(cd component && go tool cover -func=component_coverage.tmp 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SERVICE_COV" ]; then SERVICE_COV="0"; fi; \
	SHARED_COV=$$(cd shared && go tool cover -func=shared_coverage.tmp 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SHARED_COV" ]; then SHARED_COV="0"; fi; \
	rm -f coverage.out; \
	echo "mode: set" > coverage.out; \
	if [ -f component/component_coverage.tmp ]; then \
		tail -n +2 component/component_coverage.tmp >> coverage.out; \
	fi; \
	if [ -f shared/shared_coverage.tmp ]; then \
		tail -n +2 shared/shared_coverage.tmp >> coverage.out; \
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
		echo "   • component/component_coverage.tmp"; \
		echo "   • shared/shared_coverage.tmp"; \
		echo "   • coverage.out"; \
		cd component && $(MAKE) build BUILD_TAGS="$(BUILD_TAGS)"; \
	else \
		echo "❌ Coverage $$COMBINED_COV% below $(SERVICE_COVERAGE_THRESHOLD)% threshold - build blocked"; \
		exit 1; \
	fi

component-build-with-coverage-local: component-coverage-local shared-coverage-local
	@echo "Building component with coverage enforcement (local development mode)..."
	@SERVICE_COV=$$(cd component && go tool cover -func=component_coverage.tmp 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SERVICE_COV" ]; then SERVICE_COV="0"; fi; \
	SHARED_COV=$$(cd shared && go tool cover -func=shared_coverage.tmp 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$SHARED_COV" ]; then SHARED_COV="0"; fi; \
	rm -f coverage.out; \
	echo "mode: set" > coverage.out; \
	if [ -f component/component_coverage.tmp ]; then \
		tail -n +2 component/component_coverage.tmp >> coverage.out; \
	fi; \
	if [ -f shared/shared_coverage.tmp ]; then \
		tail -n +2 shared/shared_coverage.tmp >> coverage.out; \
	fi; \
	COMBINED_COV=$$(go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$COMBINED_COV" ]; then \
		echo "⚠️  Warning: Failed to calculate combined coverage, trying alternative method..."; \
		COMBINED_COV=$$(go tool cover -func=coverage.out 2>&1 | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
		if [ -z "$$COMBINED_COV" ]; then COMBINED_COV="0"; fi; \
	fi; \
	echo "🎯 Component Coverage (local): $$SERVICE_COV% | Shared: $$SHARED_COV% | Combined: $$COMBINED_COV%"; \
	if [ $$(echo "$$COMBINED_COV >= $(SERVICE_COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "✅ Coverage $$COMBINED_COV% meets $(SERVICE_COVERAGE_THRESHOLD)% threshold - proceeding with build (local)"; \
		echo "📁 Coverage files preserved for org tools:"; \
		echo "   • component/component_coverage.tmp"; \
		echo "   • shared/shared_coverage.tmp"; \
		echo "   • coverage.out"; \
		cd component && $(MAKE) build BUILD_TAGS="$(LOCAL_BUILD_TAGS)"; \
	else \
		echo "❌ Coverage $$COMBINED_COV% below $(SERVICE_COVERAGE_THRESHOLD)% threshold - build blocked"; \
		exit 1; \
	fi

component-coverage:
	@echo "Running component coverage..."
	@cd component && $(MAKE) coverage BUILD_TAGS="$(BUILD_TAGS)"

component-coverage-local:
	@echo "Running component coverage (local development mode)..."
	@cd component && $(MAKE) coverage BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

component-test:
	@echo "Running component tests..."
	@cd component && $(MAKE) test BUILD_TAGS="$(BUILD_TAGS)"

component-test-local:
	@echo "Running component tests (local development mode)..."
	@cd component && $(MAKE) test BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

component-clean:
	@echo "Cleaning component..."
	@cd component && $(MAKE) clean

component-deep-clean:
	@echo "Deep cleaning component..."
	@cd component && $(MAKE) deep-clean

component-deps:
	@echo "Installing component dependencies..."
	@cd component && $(MAKE) deps

component-fmt:
	@echo "Formatting component code..."
	@cd component && $(MAKE) fmt

component-vet:
	@echo "Vetting component code..."
	@cd component && $(MAKE) vet

component-lint:
	@echo "Linting component code..."
	@cd component && $(MAKE) lint

component-run:
	@echo "Running component..."
	@cd component && $(MAKE) run BUILD_TAGS="$(BUILD_TAGS)"

component-run-local:
	@echo "Running component (local development mode)..."
	@cd component && $(MAKE) run BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

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

# Full end-to-end test (requires component to be running)
e2e-test: component-build testrunner-build
	@echo "Starting end-to-end test..."
	@echo "Starting component in background..."
	@cd component && ./bin/component & SERVICE_PID=$$!; \
	sleep 3; \
	echo "Running integration tests..."; \
	cd testrunner && ./bin/testrunner; \
	TEST_RESULT=$$?; \
	echo "Stopping component..."; \
	kill $$SERVICE_PID 2>/dev/null || true; \
	exit $$TEST_RESULT

# ==============================================================================
# LOCAL DEVELOPMENT TARGETS (using run_tests_local.sh)
# ==============================================================================

.PHONY: local-test local-test-run local-build local-clean local-dev local-watch local-setup local-lint local-unit-test local-coverage local-help

# Build components with coverage and run tests (main workflow)
local-test:
	@echo "🚀 Running full test suite with build..."
	./run_tests_local.sh build

# Run tests with existing binaries (quick workflow)
local-test-run:
	@echo "⚡ Running tests with existing binaries..."
	./run_tests_local.sh run

# Build only (without running tests)
local-build:
	@echo "🔨 Building components for local development..."
	@cd component && go build -tags local -cover -o bin/component ./cmd
	@cd testrunner && go build -tags local -o bin/testrunner ./cmd
	@echo "✅ Local build completed"

# Generate coverage report only
local-coverage:
	@echo "📊 Generating coverage report..."
	@cd component && \
	if [ -d "coverage" ] && [ -n "$$(ls -A coverage 2>/dev/null)" ]; then \
		go tool covdata textfmt -i=coverage -o=coverage/coverage.out && \
		go tool cover -html=coverage/coverage.out -o=coverage/coverage.html && \
		echo "Coverage report: component/coverage/coverage.html"; \
	else \
		echo "No coverage data found. Run 'make local-test' first."; \
	fi

# Clean all build artifacts
local-clean:
	@echo "🧹 Cleaning local build artifacts..."
	rm -rf component/bin/ component/coverage/
	rm -rf testrunner/bin/ testrunner/results/
	rm -f testrunner/config.yaml
	rm -rf testrunner/testdata/
	@echo "✅ Local clean completed"

# Quick development cycle (build + run tests + show coverage)
local-dev: local-test local-coverage
	@echo "�� Local development cycle completed!"

# Watch mode (requires fswatch - install with: brew install fswatch)
local-watch:
	@echo "👀 Watching for changes (requires fswatch)..."
	@which fswatch > /dev/null || (echo "❌ fswatch not found. Install with: brew install fswatch" && exit 1)
	@fswatch -o component/ testrunner/ shared/ | while read num; do \
		echo "🔄 Changes detected, running tests..."; \
		make local-test-run; \
		echo "⏳ Waiting for changes..."; \
	done

# Setup development environment
local-setup:
	@echo "🛠️  Setting up local development environment..."
	@go mod download
	@cd component && go mod download
	@cd testrunner && go mod download  
	@cd shared && go mod download
	@echo "✅ Local setup completed"

# Run linting and formatting
local-lint:
	@echo "🔍 Running linting and formatting..."
	@cd component && go fmt ./... && go vet ./...
	@cd testrunner && go fmt ./... && go vet ./...
	@cd shared && go fmt ./... && go vet ./...
	@echo "✅ Local linting completed"

# Run unit tests only (not integration tests)
local-unit-test:
	@echo "🧪 Running unit tests..."
	@cd component && go test -tags local ./...
	@cd testrunner && go test -tags local ./...
	@cd shared && go test -tags local ./...

# Help for local development targets
local-help:
	@echo "📚 Katharos Local Development Commands:"
	@echo ""
	@echo "🎯 Quick Start:"
	@echo "  make local-test        - Build components with coverage and run integration tests"
	@echo "  make local-test-run    - Run integration tests with existing binaries (quick)"
	@echo "  make local-dev         - Full development cycle (test + coverage report)"
	@echo ""
	@echo "🔨 Building:"
	@echo "  make local-build       - Build components in local mode"
	@echo "  make local-clean       - Clean all build artifacts"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make local-unit-test   - Run unit tests only"
	@echo "  make local-coverage    - Generate coverage report"
	@echo ""
	@echo "⚡ Development:"
	@echo "  make local-watch       - Watch for changes and run tests (requires fswatch)"
	@echo "  make local-lint        - Run linting and formatting"
	@echo "  make local-setup       - Setup development environment"
	@echo ""
	@echo "🚀 Script Interface:"
	@echo "  ./run_tests_local.sh build  - Full workflow with script"
	@echo "  ./run_tests_local.sh run    - Quick run with script"
	@echo ""
	@echo "ℹ️  Use 'make help' for original makefile commands"
