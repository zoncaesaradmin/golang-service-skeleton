# Root Makefile for Katharos Project
.PHONY: help package package-local pre_build post_build clean

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
	@echo "Parameters:"
	@echo "  ENTRYPOINT_MODE - Mode to pass to src/Makefile (default: start)"
	@echo "  BUILD_TAGS      - Additional build tags to use (empty by default)"
	@echo ""
	@echo "Examples:"
	@echo "  make package                       # Uses default ENTRYPOINT_MODE=start (production)"
	@echo "  make package-local                 # Uses default ENTRYPOINT_MODE=start (local dev)"
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
