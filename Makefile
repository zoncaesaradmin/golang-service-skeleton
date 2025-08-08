# Root Makefile for Katharos Project
.PHONY: help package pre_build post_build clean

# Default ENTRYPOINT_MODE if not set
ENTRYPOINT_MODE ?= start

# Default target
help:
	@echo "Katharos Project - Root Makefile"
	@echo "================================"
	@echo ""
	@echo "Available targets:"
	@echo "  package      - Run specified ENTRYPOINT_MODE target in src/ (default: start)"
	@echo "  pre_build    - Pre-build target (depends on package)"
	@echo "  post_build   - Post-build target (depends on clean)"
	@echo "  clean        - Clean target"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Parameters:"
	@echo "  ENTRYPOINT_MODE - Mode to pass to src/Makefile (default: start)"
	@echo ""
	@echo "Examples:"
	@echo "  make package                       # Uses default ENTRYPOINT_MODE=start"
	@echo "  make package ENTRYPOINT_MODE=build    # Uses ENTRYPOINT_MODE=build"
	@echo "  make package ENTRYPOINT_MODE=test     # Uses ENTRYPOINT_MODE=test"

# Package target - runs the specified ENTRYPOINT_MODE in src/
package:
	@echo "Running package target with ENTRYPOINT_MODE=$(ENTRYPOINT_MODE)"
	@$(MAKE) -C src $(ENTRYPOINT_MODE)

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
