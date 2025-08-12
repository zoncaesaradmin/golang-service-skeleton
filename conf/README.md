# Configuration Documentation

This document explains the configuration system for the Katharos component and test infrastructure.

## Configuration Files

### Primary Configuration: `conf/config.yaml`
- **Location**: `/conf/config.yaml` (root of project)
- **Purpose**: Main configuration file for the component (production and development)
- **Tracked**: Yes (committed to git)
- **Usage**: Contains settings for server, logging, etc.

### Test Configuration: `conf/testconfig.yaml`
- **Location**: `/conf/testconfig.yaml` (root of project)
- **Purpose**: Configuration for integration test runner
- **Tracked**: Yes (committed to git)
- **Usage**: Contains testrunner-specific settings like timeouts, paths, etc.

## Configuration Loading

### Component Configuration
The component loads configuration using the `HOME_DIR` environment variable:
- **Path**: `$HOME_DIR/conf/config.yaml`
- **Required**: `HOME_DIR` must be set to repository root
- **Method**: Absolute path loading for location independence

### Testrunner Configuration
The testrunner loads configuration using the `HOME_DIR` environment variable:
- **Path**: `$HOME_DIR/conf/testconfig.yaml`
- **Required**: `HOME_DIR` must be set to repository root
- **Method**: Automatic loading without command line arguments

## Configuration Structure

### Component Configuration (`config.yaml`)
```yaml
# Server configuration
server:
  host: "localhost"              # SERVER_HOST
  port: 8080                     # SERVER_PORT
  read_timeout: 10               # SERVER_READ_TIMEOUT
  write_timeout: 10              # SERVER_WRITE_TIMEOUT

# Logging configuration
logging:
  level: "info"                  # LOG_LEVEL (debug, info, warn, error)
  format: "json"                 # LOG_FORMAT (json, text)
  file_path: "/tmp/katharos-component.log"  # LOG_FILE_PATH
```

### Test Configuration (`testconfig.yaml`)
```yaml
# Component settings for testing
component:
  binary_path: "../component/bin/component"
  port: 8080
  timeout: 30s

# Message bus configuration
messagebus:
  type: local

# Test data paths
testdata:
  scenarios_path: testdata/scenarios
  fixtures_path: testdata/fixtures

# Validation settings
validation:
  timeout: 1m0s
  max_retries: 3
  retry_delay: 1s
```

## Environment Variables

Component configuration values can be overridden using environment variables. The environment variable names are shown in comments in the YAML structure above.

### Example:
```bash
export HOME_DIR="/path/to/katharos"
export SERVER_PORT=9090
export LOG_LEVEL=debug
./bin/component
```

## Development Workflow

### For Local Development:
1. Set `HOME_DIR` environment variable to repository root
2. Use `make run-local` to automatically set HOME_DIR and run component
3. Modify `conf/config.yaml` as needed for development
4. All changes are tracked in git

### For Production:
1. Set `HOME_DIR` to deployment directory
2. Override sensitive/environment-specific values using environment variables
3. Never commit passwords or secrets to git

### For Testing:
1. Test infrastructure automatically sets `HOME_DIR`
2. Use `./test/run_tests_local.sh` to run integration tests
3. Testrunner uses `conf/testconfig.yaml` automatically
4. Log files are directed to `test/results/logs/` during test runs

## Make Targets

- `make run-local`: Build and run component with HOME_DIR automatically set
- `make run-local-coverage`: Build and run with coverage instrumentation for testing

## File Locations When Running

- **Component**: Always loads from `$HOME_DIR/conf/config.yaml`
- **Testrunner**: Always loads from `$HOME_DIR/conf/testconfig.yaml`
- **Location Independent**: Both work from any directory with proper HOME_DIR

## Best Practices

1. **Set HOME_DIR**: Always ensure HOME_DIR points to repository root
2. **Use make targets**: Prefer `make run-local` over direct binary execution
3. **Keep secrets out of config files**: Use environment variables for passwords, API keys, etc.
4. **Document changes**: Update this file when adding new configuration options
5. **Test with defaults**: Ensure components work with the provided configurations
6. **Validate environment variables**: Check that environment overrides work correctly
