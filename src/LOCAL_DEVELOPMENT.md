# Local Development Guide

This guide covers the local development setup for the Katharos project, including building, testing, and generating coverage reports with comprehensive log management.

## Quick Start

### One-Step Build and Test
```bash
# Build both components and run full test suite
./run_tests_local.sh

# Or using make
make local-test
```

### Available Commands

#### Shell Script Commands
```bash
# Build and run (default)
./run_tests_local.sh build

# Run only (assumes binaries exist)
./run_tests_local.sh run
```

#### Make Targets
```bash
# Main workflows
make local-test          # Build and run full test suite
make local-test-run       # Run tests with existing binaries
make local-build          # Build both component and testrunner
make local-coverage       # Generate coverage reports
make local-clean          # Clean up build artifacts

# Development utilities
make local-dev            # Start development environment
make local-watch          # Watch mode (requires fswatch)
make local-lint           # Run linting
make local-unit-test      # Run unit tests only

# Setup and help
make local-setup          # One-time setup
make local-help           # Show help
```

## Output Files and Logs

All test results, logs, and reports are generated in the `results/` directory:

### Main Reports
- `test_report.txt` - Comprehensive test report with all details
- `coverage.html` - Interactive HTML coverage report
- `coverage.out` - Raw coverage data
- `coverage_summary.txt` - Coverage statistics

### Log Files (organized in `results/logs/`)
- `component.log` - Component application logs (structured JSON)
- `component_stdout.log` - Component standard output
- `component_stderr.log` - Component standard error output
- `testrunner.log` - Testrunner application logs (if using shared logging)
- `testrunner_stdout.log` - Testrunner standard output
- `testrunner_stderr.log` - Testrunner standard error output

### Consolidated Logs
- `all_logs.txt` - All logs consolidated in one file with sections
- `testrunner_output.log` - Combined testrunner output (legacy format)

## Logging Configuration

The local development setup automatically configures logging for both components:

### Component Logging
- **Structured Logging**: Uses JSON format via zerolog library
- **Configurable Path**: Set via `LOG_FILE_PATH` environment variable
- **Default Path**: `/tmp/katharos-component.log` (production) or `results/logs/component.log` (local dev)
- **Level Control**: Configurable via `LOG_LEVEL` environment variable

### Testrunner Logging
- **Standard Logging**: Uses Go's standard log package
- **Output Capture**: Both stdout and stderr are captured separately
- **Organized Storage**: All logs stored in `results/logs/` directory

### Environment Variables for Log Control
```bash
# Component logging
export LOG_FILE_PATH="./results/logs/component.log"
export LOG_LEVEL="info"  # debug, info, warn, error

# These are automatically set by run_tests_local.sh
```

## Coverage Reports

The system automatically instruments binaries with coverage tracking and generates:

1. **HTML Report**: Open `results/coverage.html` in your browser for interactive coverage exploration
2. **Summary Report**: View `results/coverage_summary.txt` for quick coverage statistics
3. **Console Output**: Coverage percentage is displayed after test completion

## Build Tags

The local development setup uses `-tags local` to:
- Avoid external dependencies (like Kafka)
- Use local message bus implementation
- Enable local-specific configurations

## Process Management

The test runner automatically:
- Starts the component with coverage instrumentation and log configuration
- Runs the test suite with organized log capture
- Stops all processes on completion
- Organizes and consolidates all logs
- Cleans up temporary files

## Requirements

### Required
- Go 1.22.5 or later
- Make utility

### Optional
- `fswatch` (for watch mode)
- `jq` (for JSON processing)

### Installation on macOS
```bash
# Optional dependencies
brew install fswatch jq
```

## Directory Structure

```
src/
├── component/           # Main component source
│   ├── bin/            # Built binaries
│   └── cmd/main.go     # Component entry point
├── testrunner/         # Test runner source
│   ├── bin/            # Built binaries
│   └── cmd/main.go     # Testrunner entry point
├── results/            # Generated test results and reports
│   ├── logs/           # All log files organized here
│   ├── coverage.html   # Coverage report
│   ├── test_report.txt # Comprehensive test report
│   └── all_logs.txt    # Consolidated log file
├── coverage/           # Coverage data files (temporary)
└── run_tests_local.sh  # Main development script
```

## Log Analysis

### Quick Log Review
```bash
# View all logs consolidated
cat results/all_logs.txt

# View specific component logs
cat results/logs/component.log

# View test execution output
cat results/logs/testrunner_stdout.log

# Check for errors across all logs
grep -i error results/logs/*.log
```

### Log File Descriptions

1. **`component.log`**: Structured application logs from the main component
   - Server startup/shutdown events
   - Request processing logs
   - Error conditions and warnings
   - Configuration information

2. **`component_stdout.log`** / **`component_stderr.log`**: Raw output streams
   - System-level output
   - Panic traces (if any)
   - Operating system messages

3. **`testrunner_stdout.log`** / **`testrunner_stderr.log`**: Test execution logs
   - Test scenario execution details
   - Success/failure messages
   - Performance metrics
   - Validation results

4. **`all_logs.txt`**: Complete consolidated view
   - All logs in chronological order by component
   - Useful for timeline analysis
   - Single file for external log analysis tools

## Troubleshooting

### Build Issues
1. Ensure Go modules are properly initialized
2. Check for import cycle issues
3. Verify all dependencies are available

### Test Failures
1. Check `results/logs/testrunner_stdout.log` for detailed error messages
2. Review `results/logs/component.log` for component-side issues
3. Verify component starts successfully in `results/logs/component_stdout.log`
4. Check port conflicts or permission issues

### Coverage Issues
1. Ensure `GOCOVERDIR` environment is set correctly
2. Check that binaries are built with `-cover` flag
3. Verify coverage directory permissions

### Logging Issues
1. Check `LOG_FILE_PATH` environment variable settings
2. Verify directory permissions for log file creation
3. Review file system space availability
4. Check for file handle limits

### Watch Mode Issues
1. Install fswatch: `brew install fswatch`
2. Check file permissions in watched directories
3. Verify PATH includes fswatch location

## Daily Workflow

1. **Start Development**:
   ```bash
   make local-test
   ```

2. **Check Results**:
   - Open `results/coverage.html` for coverage analysis
   - Review `results/test_report.txt` for summary
   - Check `results/all_logs.txt` for comprehensive log review

3. **Debug Issues**:
   ```bash
   # Quick error check
   grep -i error results/logs/*.log
   
   # Component-specific logs
   cat results/logs/component.log
   
   # Test execution details
   cat results/logs/testrunner_stdout.log
   ```

4. **Development Cycle**:
   ```bash
   # After making changes
   make local-test
   
   # Or for continuous development
   make local-watch
   ```

5. **Clean Up**:
   ```bash
   make local-clean
   ```

## Integration with IDE

The local development setup integrates with VS Code through:
- Build tasks for compilation
- Coverage reports viewable in browser
- Terminal integration for script execution
- Problem detection through build output
- Log files can be opened directly in editor

## Performance Notes

- Initial builds may take longer due to coverage instrumentation
- Subsequent runs are faster with cached builds
- Watch mode provides near-instant feedback for changes
- Coverage collection adds minimal runtime overhead
- Log file I/O is optimized for development use

## Best Practices

1. Always run full test suite before committing changes
2. Monitor coverage reports to maintain test quality
3. Use watch mode during active development
4. Clean build artifacts regularly
5. Review test logs for performance insights
6. Check consolidated logs for issue patterns
7. Use structured logging data for debugging

## Log Management Tips

1. **Regular Cleanup**: Use `make local-clean` to remove old log files
2. **Log Analysis**: Use grep, awk, or jq for structured log analysis
3. **Error Tracking**: Monitor error patterns across test runs
4. **Performance Monitoring**: Track component startup and response times in logs
5. **Debug Information**: Enable verbose logging for detailed troubleshooting
