# Root Makefile for Katharos Project
.PHONY: help package package-local pre_build post_build clean test test-local test-run test-run-local test-build test-clean test-coverage test-setup test-help

# Default ENTRYPOINT_MODE if not set
ENTRYPOINT_MODE ?= start

# Build tags (empty by default, but can be set from environment/command line)
BUILD_TAGS ?=

# Helper function to safely append a tag if not already present
define append_tag
$(if $(filter $(1),$(BUILD_TAGS)),$(BUILD_TAGS),$(if $(BUILD_TAGS),$(BUILD_TAGS) $(1),$(1)))
endef

# Local build tags - safely append 'local' to existing BUILD_TAGS
LOCAL_BUILD_TAGS = $(call append_tag,local)

# Default target
help:
	@echo "Katharos Project - Root Makefile"
	@echo "================================"
	@echo ""
	@echo "Available targets:"
	@echo "  package      - Run specified ENTRYPOINT_MODE target in src/ (production mode, default)"
	@echo "  package-local - Run specified ENTRYPOINT_MODE target in src/ (local development mode)"
	@echo "  pre_build    - Pre-build target (depends on package)"
	@echo "  post_build   - Post-build target (depends on clean)"
	@echo "  clean        - Clean target"
	@echo "  help         - Show this help"
	@echo ""
	@echo " Test Infrastructure (in test/ directory):"
	@echo "  test         - Build components and run integration tests (primary workflow)"
	@echo "  test-local   - Same as test (local dev mode - default behavior)"
	@echo "  test-run     - Run tests with existing binaries (quick workflow)"
	@echo "  test-build   - Build components only (no tests)"
	@echo "  test-clean   - Clean test artifacts and logs"
	@echo "  test-coverage- Generate coverage report from latest test run"
	@echo "  test-setup   - Setup test environment"
	@echo "  test-help    - Show detailed test help"
	@echo ""
	@echo "Parameters:"
	@echo "  ENTRYPOINT_MODE - Mode to pass to src/Makefile (default: start)"
	@echo "  BUILD_TAGS      - Additional build tags to use (empty by default)"
	@echo ""
	@echo "Examples:"
	@echo "  make package                       # Uses default ENTRYPOINT_MODE=start (production)"
	@echo "  make package-local                 # Uses default ENTRYPOINT_MODE=start (local dev)"
	@echo "  make test                          # Build and run integration tests"
	@echo "  make test-run                      # Quick test run with existing binaries"
	@echo "  make package ENTRYPOINT_MODE=build    # Uses ENTRYPOINT_MODE=build (production)"
	@echo "  make package-local ENTRYPOINT_MODE=test # Uses ENTRYPOINT_MODE=test (local dev)"
	@echo "  make package-local BUILD_TAGS=\"debug test\" # Appends 'local' to existing tags"

# Package target - runs the specified ENTRYPOINT_MODE in src/ (production mode)
package:
	@echo "Running package target with ENTRYPOINT_MODE=$(ENTRYPOINT_MODE) (production mode)"
	@$(MAKE) -C src $(ENTRYPOINT_MODE) BUILD_TAGS="$(BUILD_TAGS)"

# Package target with local development build tags
package-local:
	@echo "Running package target with ENTRYPOINT_MODE=$(ENTRYPOINT_MODE) (local development mode)"
	@$(MAKE) -C src $(ENTRYPOINT_MODE) BUILD_TAGS="$(LOCAL_BUILD_TAGS)"

# Pre-build target that depends on package
pre_build: package
	@echo "Pre-build completed"

# Clean target
clean:
	@echo "Cleaning project..."
	@$(MAKE) -C src clean
	@echo "✅ Clean completed"

# Post-build target that depends on clean
post_build: clean
	@echo "Post-build completed"

# ==============================================================================
# TEST INFRASTRUCTURE (test/ directory)
# ==============================================================================

# Primary test workflow - build and run integration tests
test:
	@echo "🚀 Running integration tests (build + test)..."
	@cd test && ./run_tests_local.sh build

# Alias for test (local development is the default mode)
test-local: test

# Quick test run with existing binaries
test-run:
	@echo "⚡ Running tests with existing binaries..."
	@cd test && ./run_tests_local.sh run

# Alias for test-run in local mode
test-run-local: test-run

# Build components only (no tests)
test-build:
	@echo "🔨 Building components for testing..."
	@cd test && ./run_tests_local.sh build > /dev/null 2>&1 || true
	@echo "✅ Test build completed"

# Clean test artifacts
test-clean:
	@echo "🧹 Cleaning test artifacts..."
	@cd test && rm -rf results/ coverage/ *.pid 2>/dev/null || true
	@echo "✅ Test artifacts cleaned"

# Generate coverage report from latest test run
test-coverage:
	@echo "📊 Generating coverage report..."
	@cd test && \
	if [ -d "coverage" ] && [ -n "$$(ls -A coverage 2>/dev/null)" ]; then \
		mkdir -p results && \
		cd ../src/component && \
		go tool covdata textfmt -i=../../test/coverage -o=../../test/results/coverage.out && \
		go tool cover -html=../../test/results/coverage.out -o=../../test/results/coverage.html && \
		cd ../../test && \
		echo "✅ Coverage report: test/results/coverage.html"; \
		go tool cover -func=results/coverage.out | tail -1; \
	else \
		echo "❌ No coverage data found. Run 'make test' first."; \
	fi

# Setup test environment
test-setup:
	@echo "🛠️  Setting up test environment..."
	@cd src && go mod download
	@cd src/component && go mod download
	@cd src/testrunner && go mod download  
	@cd src/shared && go mod download
	@echo "✅ Test environment setup completed"

# Help for test targets
test-help:
	@echo "📚 Katharos Test Infrastructure (test/ directory):"
	@echo ""
	@echo "🎯 Primary Commands:"
	@echo "  make test          - Build components and run integration tests (main workflow)"
	@echo "  make test-run      - Run tests with existing binaries (quick iteration)"
	@echo ""
	@echo "🔨 Building:"
	@echo "  make test-build    - Build components only (no tests)"
	@echo "  make test-clean    - Clean test artifacts and logs"
	@echo ""
	@echo "📊 Reports:"
	@echo "  make test-coverage - Generate coverage report from latest test run"
	@echo ""
	@echo "⚙️  Setup:"
	@echo "  make test-setup    - Setup test environment and dependencies"
	@echo ""
	@echo "🚀 Script Interface:"
	@echo "  cd test && ./run_tests_local.sh build  - Full workflow"
	@echo "  cd test && ./run_tests_local.sh run    - Quick run"
	@echo ""
	@echo "📁 Output Directories (in test/):"
	@echo "  • results/         - Test reports and logs"
	@echo "  • results/logs/    - Component and testrunner logs"
	@echo "  • coverage/        - Coverage data"
	@echo ""
	@echo "ℹ️  The test infrastructure is now in test/ parallel to src/"
	@echo "ℹ️  Use 'make help' for original makefile commands"
