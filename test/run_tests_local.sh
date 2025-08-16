#!/bin/bash

# Local Development Test Runner
# Usage: ./run_tests_local.sh [build|run]
# Default mode: build (builds and runs)

set -e
# Note: We handle test failures gracefully to ensure log collection

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration - Updated paths to work from both root and test directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ "$(basename "$SCRIPT_DIR")" == "test" ]]; then
    # Running from test directory or script in test directory
    ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
    COMPONENT_DIR="$ROOT_DIR/src/component"
    TESTRUNNER_DIR="$ROOT_DIR/src/testrunner"
    RESULTS_DIR="$SCRIPT_DIR/results"
    LOGS_DIR="$SCRIPT_DIR/results/logs"
    COVERAGE_DIR="$SCRIPT_DIR/coverage"
else
    # Fallback: assume we're in root and test is a subdirectory
    ROOT_DIR="$(pwd)"
    COMPONENT_DIR="$ROOT_DIR/src/component"
    TESTRUNNER_DIR="$ROOT_DIR/src/testrunner"
    RESULTS_DIR="$ROOT_DIR/test/results"
    LOGS_DIR="$ROOT_DIR/test/results/logs"
    COVERAGE_DIR="$ROOT_DIR/test/coverage"
fi
BUILD_MODE="${1:-build}"

# Function to print colored output
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to cleanup processes and directories
cleanup() {
    log_info "Cleaning up processes and temporary files..."
    
    # Kill any running component or testrunner processes
    pkill -f "component.bin" 2>/dev/null || true
    pkill -f "testrunner.bin" 2>/dev/null || true
    
    # Clean up coverage data
    if [ -d "$COVERAGE_DIR" ]; then
        rm -rf "$COVERAGE_DIR"
    fi
    
    log_success "Cleanup completed"
}

# Function to setup directories
setup_directories() {
    log_info "Setting up directories..."
    mkdir -p "$RESULTS_DIR"
    mkdir -p "$LOGS_DIR"
    mkdir -p "$COVERAGE_DIR"
    log_success "Directories created"
}

# Function to build component with coverage
build_component() {
    log_info "Component will be built by make run-local-coverage with coverage instrumentation and local tags..."
    
    # Verify component directory exists and has Makefile
    if [ ! -f "$COMPONENT_DIR/Makefile" ]; then
        log_error "Component Makefile not found at $COMPONENT_DIR/Makefile"
        exit 1
    fi
    
    log_success "Component build preparation completed (will use make run-local-coverage)"
}

# Function to build testrunner
build_testrunner() {
    log_info "Building testrunner..."
    cd "$TESTRUNNER_DIR"
    
    # Build testrunner
    go build -tags local -o bin/testrunner.bin cmd/main.go
    
    if [ $? -eq 0 ]; then
        log_success "Testrunner built successfully"
    else
        log_error "Testrunner build failed"
        exit 1
    fi
    
    cd - > /dev/null
}

# Function to run component
run_component() {
    log_info "Starting component using make run-local-coverage..."
    cd "$COMPONENT_DIR"
    
    # Set coverage directory and log file path (using absolute paths)
    export GOCOVERDIR="$COVERAGE_DIR"
    export LOG_FILE_PATH="$LOGS_DIR/component.log"
    
    # Set test-specific topics for integration testing
    export PROCESSING_INPUT_TOPICS="test_input"
    export PROCESSING_OUTPUT_TOPIC="test_output"
    
    # Use make run-local-coverage which automatically sets HOME_DIR and builds with local tags + coverage
    # Run in background and capture stdout/stderr
    make run-local-coverage > "$LOGS_DIR/component_stdout.log" 2> "$LOGS_DIR/component_stderr.log" &
    COMPONENT_PID=$!
    echo $COMPONENT_PID > "$ROOT_DIR/test/component.pid"
    
    log_success "Component started with PID $COMPONENT_PID using make run-local-coverage (HOME_DIR=$ROOT_DIR)"
    log_info "Component logs: $LOGS_DIR/component.log, $LOGS_DIR/component_stdout.log, $LOGS_DIR/component_stderr.log"
    cd - > /dev/null
    
    # Wait a moment for component to start
    sleep 2
}

# Function to run testrunner
run_testrunner() {
    log_info "Running testrunner..."
    cd "$TESTRUNNER_DIR"
    
    # Set HOME_DIR and testrunner log file path (using absolute paths)
    export HOME_DIR="$ROOT_DIR"
    export LOG_FILE_PATH="$LOGS_DIR/testrunner.log"
    
    # Temporarily disable strict error handling for testrunner execution
    set +e
    # Run testrunner and capture output
    ./bin/testrunner.bin > "$LOGS_DIR/testrunner_stdout.log" 2> "$LOGS_DIR/testrunner_stderr.log"
    TEST_RESULT=$?
    set -e
    
    # Also capture combined output for backward compatibility
    cat "$LOGS_DIR/testrunner_stdout.log" "$LOGS_DIR/testrunner_stderr.log" > "$RESULTS_DIR/testrunner_output.log"
    
    cd - > /dev/null
    
    if [ $TEST_RESULT -eq 0 ]; then
        log_success "Tests completed successfully"
    else
        log_warning "Tests completed with issues (exit code: $TEST_RESULT)"
    fi
    
    return $TEST_RESULT
}

# Function to stop component
stop_component() {
    PIDFILE="$ROOT_DIR/test/component.pid"
    if [ -f "$PIDFILE" ]; then
        COMPONENT_PID=$(cat "$PIDFILE")
        log_info "Stopping component (PID: $COMPONENT_PID)..."
        kill $COMPONENT_PID 2>/dev/null || true
        rm -f "$PIDFILE"
        log_success "Component stopped"
        
        # Wait a moment for clean shutdown and final logs
        sleep 1
    fi
}

# Function to collect and organize logs
collect_logs() {
    log_info "Collecting and organizing logs..."
    
    # Create a consolidated log file with timestamps
    CONSOLIDATED_LOG="$RESULTS_DIR/all_logs.txt"
    
    {
        echo "=================================================="
        echo "Consolidated Log Report"
        echo "Generated: $(date)"
        echo "=================================================="
        echo
        
        echo "=== COMPONENT APPLICATION LOGS ==="
        if [ -f "$LOGS_DIR/component.log" ]; then
            cat "$LOGS_DIR/component.log"
        else
            echo "No component application logs found"
        fi
        echo
        
        echo "=== COMPONENT STDOUT ==="
        if [ -f "$LOGS_DIR/component_stdout.log" ]; then
            cat "$LOGS_DIR/component_stdout.log"
        else
            echo "No component stdout logs found"
        fi
        echo
        
        echo "=== COMPONENT STDERR ==="
        if [ -f "$LOGS_DIR/component_stderr.log" ]; then
            cat "$LOGS_DIR/component_stderr.log"
        else
            echo "No component stderr logs found"
        fi
        echo
        
        echo "=== TESTRUNNER APPLICATION LOGS ==="
        if [ -f "$LOGS_DIR/testrunner.log" ]; then
            cat "$LOGS_DIR/testrunner.log"
        else
            echo "No testrunner application logs found"
        fi
        echo
        
        echo "=== TESTRUNNER STDOUT ==="
        if [ -f "$LOGS_DIR/testrunner_stdout.log" ]; then
            cat "$LOGS_DIR/testrunner_stdout.log"
        else
            echo "No testrunner stdout logs found"
        fi
        echo
        
        echo "=== TESTRUNNER STDERR ==="
        if [ -f "$LOGS_DIR/testrunner_stderr.log" ]; then
            cat "$LOGS_DIR/testrunner_stderr.log"
        else
            echo "No testrunner stderr logs found"
        fi
        echo
        
    } > "$CONSOLIDATED_LOG"
    
    log_success "Consolidated logs created at $CONSOLIDATED_LOG"
    
    # Display log file summary
    echo
    log_info "Log Files Summary:"
    if [ -d "$LOGS_DIR" ]; then
        ls -la "$LOGS_DIR" | while read line; do
            echo "  $line"
        done
    fi
    echo
}

# Function to generate coverage report
generate_coverage_report() {
    log_info "Generating coverage report..."
    
    if [ -d "$COVERAGE_DIR" ] && [ "$(ls -A "$COVERAGE_DIR")" ]; then
        # Convert binary coverage data to text format
        # Use the component's working directory for proper module resolution
        cd "$COMPONENT_DIR"
        go tool covdata textfmt -i="$COVERAGE_DIR" -o="$RESULTS_DIR/coverage.out"
        
        # Generate HTML coverage report
        go tool cover -html="$RESULTS_DIR/coverage.out" -o="$RESULTS_DIR/coverage.html"
        
        # Generate coverage summary
        go tool cover -func="$RESULTS_DIR/coverage.out" > "$RESULTS_DIR/coverage_summary.txt"
        
        cd - > /dev/null
        
        log_success "Coverage report generated at $RESULTS_DIR/coverage.html"
        
        # Show coverage summary
        echo
        log_info "Coverage Summary:"
        cat "$RESULTS_DIR/coverage_summary.txt" | tail -1
        echo
    else
        log_warning "No coverage data found"
    fi
}

# Function to generate final report
generate_report() {
    log_info "Generating test report..."
    
    REPORT_FILE="$RESULTS_DIR/test_report.txt"
    
    {
        echo "=================================================="
        echo "Local Development Test Report"
        echo "Generated: $(date)"
        echo "=================================================="
        echo
        
        echo "Build Results:"
        echo "- Component: Built successfully"
        echo "- Testrunner: Built successfully"
        echo
        
        echo "Test Execution:"
        if [ $TEST_RESULT -eq 0 ]; then
            echo "- Status: PASSED"
        else
            echo "- Status: FAILED (exit code: $TEST_RESULT)"
        fi
        echo
        
        echo "Coverage Report:"
        if [ -f "$RESULTS_DIR/coverage_summary.txt" ]; then
            cat "$RESULTS_DIR/coverage_summary.txt" | tail -1
        else
            echo "- No coverage data available"
        fi
        echo
        
        echo "Log Files:"
        echo "- Component application logs: $LOGS_DIR/component.log"
        echo "- Component stdout: $LOGS_DIR/component_stdout.log"
        echo "- Component stderr: $LOGS_DIR/component_stderr.log"
        echo "- Testrunner application logs: $LOGS_DIR/testrunner.log"
        echo "- Testrunner stdout: $LOGS_DIR/testrunner_stdout.log"
        echo "- Testrunner stderr: $LOGS_DIR/testrunner_stderr.log"
        echo "- Consolidated logs: $RESULTS_DIR/all_logs.txt"
        echo
        
        echo "Output Files:"
        echo "- Test logs (legacy): $RESULTS_DIR/testrunner_output.log"
        echo "- Coverage report: $RESULTS_DIR/coverage.html"
        echo "- Coverage summary: $RESULTS_DIR/coverage_summary.txt"
        echo "- This report: $REPORT_FILE"
        echo
        
    } > "$REPORT_FILE"
    
    log_success "Test report generated at $REPORT_FILE"
    echo
    log_info "=== QUICK SUMMARY ==="
    cat "$REPORT_FILE" | grep -A 10 "Test Execution:"
}

# Main execution
main() {
    echo
    log_info "Starting Local Development Test Runner (mode: $BUILD_MODE)"
    log_info "Script directory: $SCRIPT_DIR"
    log_info "Repository root: $ROOT_DIR"
    log_info "Component directory: $COMPONENT_DIR"
    log_info "Results directory: $RESULTS_DIR"
    echo
    
    # Setup trap for cleanup on exit - but not for normal script completion
    trap 'cleanup' INT TERM
    
    # Setup directories
    setup_directories
    
    if [ "$BUILD_MODE" = "build" ] || [ "$BUILD_MODE" = "all" ]; then
        # Build phase
        build_component
        build_testrunner
    fi
    
    if [ "$BUILD_MODE" = "run" ] || [ "$BUILD_MODE" = "build" ] || [ "$BUILD_MODE" = "all" ]; then
        # Run phase
        run_component
        
        # Run tests (handle failures gracefully)
        set +e
        run_testrunner
        TEST_RESULT=$?
        set -e
        
        echo "DEBUG: About to stop component"        # Stop component
        stop_component
        
        echo "DEBUG: About to collect logs"        # Collect and organize logs
        collect_logs
        
        # Generate reports
        generate_coverage_report
        generate_report
        
        # Manual cleanup at the end
        cleanup
        
        # Show final status
        echo
        if [ $TEST_RESULT -eq 0 ]; then
            log_success "All tests passed! Check $RESULTS_DIR/ for detailed reports."
            log_info "All logs are organized in $LOGS_DIR/"
        else
            log_error "Some tests failed. Check $RESULTS_DIR/ and $LOGS_DIR/ for detailed logs."
        fi
        echo
    fi
}

# Check if we're being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
