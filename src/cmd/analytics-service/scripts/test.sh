#!/bin/bash

# test.sh - Comprehensive Go Testing Script
# Author: Generated for go-data-api-microservices
# Description: Executes unit tests, integration tests, coverage analysis,
#              race condition detection with configurable options
# Usage: ./test.sh [OPTIONS]
# Options:
#   --verbose         Enable verbose test output
#   --html-coverage   Generate HTML coverage report
#   --package=PKG     Run tests for specific package only
#   --integration     Run integration tests only
#   --unit           Run unit tests only
#   --race           Enable race condition detection
#   --help           Show this help message

set -euo pipefail  # Exit on error, undefined vars, pipe failures
set -o errtrace    # Inherit ERR trap in functions

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
COVERAGE_DIR="$PROJECT_ROOT/coverage"
COVERAGE_FILE="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="$COVERAGE_DIR/coverage.html"

# Default options
VERBOSE=false
HTML_COVERAGE=false
SPECIFIC_PACKAGE=""
RUN_INTEGRATION=true
RUN_UNIT=true
ENABLE_RACE=false

# Colors for output
READONLY RED='\033[0;31m'
READONLY GREEN='\033[0;32m'
READONLY YELLOW='\033[1;33m'
READONLY BLUE='\033[0;34m'
READONLY NC='\033[0m' # No Color

# Error handling
trap 'error_handler $? $LINENO' ERR

error_handler() {
    local exit_code=$1
    local line_number=$2
    echo -e "${RED}ERROR: Script failed with exit code $exit_code at line $line_number${NC}" >&2
    cleanup
    exit $exit_code
}

cleanup() {
    echo -e "${YELLOW}Cleaning up temporary files...${NC}"
    # Add cleanup logic here if needed
}

log_info() {
    echo -e "${BLUE}INFO: $1${NC}"
}

log_success() {
    echo -e "${GREEN}SUCCESS: $1${NC}"
}

log_error() {
    echo -e "${RED}ERROR: $1${NC}" >&2
}

log_warning() {
    echo -e "${YELLOW}WARNING: $1${NC}"
}

show_help() {
    cat << EOF
Usage: $0 [OPTIONS]

Comprehensive Go testing script with coverage analysis and race detection.

Options:
  --verbose         Enable verbose test output (-v flag)
  --html-coverage   Generate HTML coverage report
  --package=PKG     Run tests for specific package only (e.g., --package=./cmd/analytics-service)
  --integration     Run integration tests only (files ending with _integration_test.go)
  --unit           Run unit tests only (exclude integration tests)
  --race           Enable race condition detection (-race flag)
  --help           Show this help message

Examples:
  $0                                    # Run all tests
  $0 --verbose --html-coverage         # Verbose output with HTML coverage
  $0 --package=./cmd/analytics-service # Test specific package
  $0 --unit --race                     # Unit tests with race detection
  $0 --integration --verbose           # Integration tests with verbose output

EOF
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --verbose)
                VERBOSE=true
                shift
                ;;
            --html-coverage)
                HTML_COVERAGE=true
                shift
                ;;
            --package=*)
                SPECIFIC_PACKAGE="${1#*=}"
                shift
                ;;
            --integration)
                RUN_INTEGRATION=true
                RUN_UNIT=false
                shift
                ;;
            --unit)
                RUN_UNIT=true
                RUN_INTEGRATION=false
                shift
                ;;
            --race)
                ENABLE_RACE=true
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    local go_version
    go_version=$(go version | cut -d' ' -f3)
    log_info "Using Go version: $go_version"
    
    if [[ ! -f "$PROJECT_ROOT/go.mod" ]]; then
        log_error "go.mod not found in project root: $PROJECT_ROOT"
        exit 1
    fi
    
    # Create coverage directory
    mkdir -p "$COVERAGE_DIR"
}

build_test_command() {
    local test_cmd="go test"
    local test_args=()
    
    # Add verbose flag
    if [[ "$VERBOSE" == "true" ]]; then
        test_args+=("-v")
    fi
    
    # Add race detection
    if [[ "$ENABLE_RACE" == "true" ]]; then
        test_args+=("-race")
    fi
    
    # Add coverage
    test_args+=("-coverprofile=$COVERAGE_FILE")
    test_args+=("-covermode=atomic")
    
    # Add timeout
    test_args+=("-timeout=300s")
    
    echo "$test_cmd ${test_args[*]}"
}

get_test_packages() {
    local packages
    
    if [[ -n "$SPECIFIC_PACKAGE" ]]; then
        packages="$SPECIFIC_PACKAGE"
    else
        packages="./..."
    fi
    
    echo "$packages"
}

build_test_pattern() {
    local pattern=""
    
    if [[ "$RUN_UNIT" == "true" && "$RUN_INTEGRATION" == "false" ]]; then
        # Run only unit tests (exclude integration tests)
        pattern="-run='^Test.*' -skip='.*Integration.*'"
    elif [[ "$RUN_INTEGRATION" == "true" && "$RUN_UNIT" == "false" ]]; then
        # Run only integration tests
        pattern="-run='.*Integration.*'"
    fi
    # If both are true or both are false, run all tests (no pattern)
    
    echo "$pattern"
}

run_tests() {
    local test_cmd
    local packages
    local pattern
    
    test_cmd=$(build_test_command)
    packages=$(get_test_packages)
    pattern=$(build_test_pattern)
    
    log_info "Running tests..."
    log_info "Command: $test_cmd $pattern $packages"
    
    cd "$PROJECT_ROOT"
    
    # Execute the test command
    if [[ -n "$pattern" ]]; then
        eval "$test_cmd $pattern $packages"
    else
        eval "$test_cmd $packages"
    fi
    
    log_success "Tests completed successfully"
}

generate_coverage_report() {
    if [[ ! -f "$COVERAGE_FILE" ]]; then
        log_warning "Coverage file not found: $COVERAGE_FILE"
        return
    fi
    
    log_info "Generating coverage report..."
    
    # Generate text coverage summary
    local coverage_percent
    coverage_percent=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}')
    log_info "Total coverage: $coverage_percent"
    
    # Generate HTML coverage report if requested
    if [[ "$HTML_COVERAGE" == "true" ]]; then
        log_info "Generating HTML coverage report: $COVERAGE_HTML"
        go tool cover -html="$COVERAGE_FILE" -o="$COVERAGE_HTML"
        log_success "HTML coverage report generated: $COVERAGE_HTML"
    fi
}

run_additional_checks() {
    log_info "Running additional code quality checks..."
    
    cd "$PROJECT_ROOT"
    
    # Check for Go modules issues
    log_info "Verifying Go modules..."
    go mod verify
    
    # Check for unused dependencies
    log_info "Tidying Go modules..."
    go mod tidy
    
    # Run go vet for static analysis
    log_info "Running go vet..."
    if [[ -n "$SPECIFIC_PACKAGE" ]]; then
        go vet "$SPECIFIC_PACKAGE"
    else
        go vet ./...
    fi
    
    log_success "Additional checks completed"
}

main() {
    echo "======================================"
    echo "  Go Testing Script Started"
    echo "======================================"
    
    parse_arguments "$@"
    check_prerequisites
    run_tests
    generate_coverage_report
    run_additional_checks
    
    echo "======================================"
    log_success "All tests and checks completed successfully!"
    echo "======================================"
    
    if [[ "$HTML_COVERAGE" == "true" && -f "$COVERAGE_HTML" ]]; then
        log_info "Open coverage report: file://$COVERAGE_HTML"
    fi
}

# Execute main function with all arguments
main "$@"
