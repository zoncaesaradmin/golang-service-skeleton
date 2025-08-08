# Shared Module

This module contains common utilities and packages shared between the service and testrunner modules.

## Packages

- `utils/` - Common utility functions
- `types/` - Shared data types and structures
- `validator/` - Common validation logic

## Coverage Target

- **Target**: 85% coverage
- **Enforcement**: Enabled via Makefile targets
- **Reporting**: HTML and terminal output available

## Usage

Import from other modules in the workspace:
```go
import "sharedmodule/utils"
import "sharedmodule/types"
import "sharedmodule/validator"
```
