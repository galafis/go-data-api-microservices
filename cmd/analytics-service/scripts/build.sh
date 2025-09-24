#!/bin/bash

# Analytics Service Build Script
# Builds and compiles the Analytics Service with optimization flags
#
# Usage:
#   ./build.sh [OPTIONS]
#
# Options:
#   --version=VERSION    Set build version (default: auto-generated)
#   --os=OS              Target operating system (default: current)
#   --arch=ARCH          Target architecture (default: current) 
#   --debug              Build with debug symbols
#   --help               Show this help message
#
# Environment Variables:
#   BUILD_VERSION        Override version
#   BUILD_OS             Override target OS
#   BUILD_ARCH           Override target architecture
#   DEBUG                Enable debug mode (1/true)

set -euo pipefail

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
SERVICE_DIR="$PROJECT_ROOT/cmd/analytics-service"
BIN_DIR="$PROJECT_ROOT/bin"
BUILD_DIR="$PROJECT_ROOT/build"
LOG_FILE="$BUILD_DIR/build.log"

# Default values
VERSION=${BUILD_VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev-$(date +%Y%m%d%H%M%S)")}
TARGET_OS=${BUILD_OS:-$(go env GOOS)}
TARGET_ARCH=${BUILD_ARCH:-$(go env GOARCH)}
DEBUG_MODE=${DEBUG:-false}
VERBOSE=false
BINARY_NAME="analytics-service"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${BLUE}[INFO]${NC} $*" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*" | tee -a "$LOG_FILE"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*" | tee -a "$LOG_FILE"
}

# Usage function
show_help() {
    cat << EOF
Analytics Service Build Script

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --version=VERSION    Set build version (default: $VERSION)
    --os=OS              Target operating system (default: $TARGET_OS)
    --arch=ARCH          Target architecture (default: $TARGET_ARCH)
    --debug              Build with debug symbols
    --verbose            Enable verbose output
    --help               Show this help message

ENVIRONMENT VARIABLES:
    BUILD_VERSION        Override version
    BUILD_OS             Override target OS
    BUILD_ARCH           Override target architecture
    DEBUG                Enable debug mode (1/true)

EXAMPLES:
    $0                                    # Basic build
    $0 --version=v1.2.3                  # Build with specific version
    $0 --os=linux --arch=amd64           # Cross-compile for Linux/AMD64
    $0 --debug                           # Build with debug symbols
    DEBUG=1 $0                           # Enable debug mode via env var

EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version=*)
                VERSION="${1#*=}"
                shift
                ;;
            --os=*)
                TARGET_OS="${1#*=}"
                shift
                ;;
            --arch=*)
                TARGET_ARCH="${1#*=}"
                shift
                ;;
            --debug)
                DEBUG_MODE=true
                shift
                ;;
            --verbose)
                VERBOSE=true
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

# Validate prerequisites
validate_prerequisites() {
    log "Validating prerequisites..."
    
    # Check if Go is installed
    if ! command -v go >/dev/null 2>&1; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    local go_version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    log "Go version: $go_version"
    
    # Check if we're in a Go module
    if [[ ! -f "$PROJECT_ROOT/go.mod" ]]; then
        log_error "go.mod not found. Please run from project root or ensure Go module is initialized."
        exit 1
    fi
    
    # Validate target OS/ARCH combination
    if ! go tool dist list | grep -q "^${TARGET_OS}/${TARGET_ARCH}$"; then
        log_error "Unsupported OS/ARCH combination: ${TARGET_OS}/${TARGET_ARCH}"
        log "Supported combinations:"
        go tool dist list | head -10
        log "... and more. Run 'go tool dist list' for complete list."
        exit 1
    fi
    
    log_success "Prerequisites validated"
}

# Setup build environment
setup_build_env() {
    log "Setting up build environment..."
    
    # Create necessary directories
    mkdir -p "$BIN_DIR" "$BUILD_DIR"
    
    # Clean up old build artifacts if they exist
    local binary_path="$BIN_DIR/${BINARY_NAME}"
    if [[ "$TARGET_OS" == "windows" ]]; then
        binary_path="${binary_path}.exe"
    fi
    
    if [[ -f "$binary_path" ]]; then
        log "Removing existing binary: $binary_path"
        rm -f "$binary_path"
    fi
    
    log_success "Build environment ready"
}

# Build the service
build_service() {
    log "Building Analytics Service..."
    log "Version: $VERSION"
    log "Target OS: $TARGET_OS"
    log "Target Architecture: $TARGET_ARCH"
    log "Debug Mode: $DEBUG_MODE"
    
    cd "$PROJECT_ROOT"
    
    # Set build variables
    local ldflags="-s -w"
    local gcflags=""
    
    # Add version information
    ldflags="$ldflags -X main.version=$VERSION"
    ldflags="$ldflags -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    ldflags="$ldflags -X main.gitCommit=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
    
    # Debug mode adjustments
    if [[ "$DEBUG_MODE" == "true" ]]; then
        ldflags="-X main.version=$VERSION -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
        gcflags="-N -l"
        log "Debug symbols enabled"
    fi
    
    # Verbose mode adjustments
    local build_flags=()
    if [[ "$VERBOSE" == "true" ]]; then
        build_flags+=("-v")
    fi
    
    # Set target environment
    export GOOS="$TARGET_OS"
    export GOARCH="$TARGET_ARCH"
    export CGO_ENABLED=0
    
    # Build the binary
    local binary_path="$BIN_DIR/${BINARY_NAME}"
    if [[ "$TARGET_OS" == "windows" ]]; then
        binary_path="${binary_path}.exe"
    fi
    
    local build_cmd=(
        go build
        "${build_flags[@]}"
        -ldflags "$ldflags"
    )
    
    if [[ -n "$gcflags" ]]; then
        build_cmd+=(-gcflags "$gcflags")
    fi
    
    build_cmd+=(
        -o "$binary_path"
        "./cmd/analytics-service"
    )
    
    log "Build command: ${build_cmd[*]}"
    
    if "${build_cmd[@]}" 2>&1 | tee -a "$LOG_FILE"; then
        log_success "Build completed successfully"
        
        # Display binary information
        if [[ -f "$binary_path" ]]; then
            local file_size
            file_size=$(du -h "$binary_path" | cut -f1)
            log "Binary location: $binary_path"
            log "Binary size: $file_size"
            
            # Make binary executable on Unix systems
            if [[ "$TARGET_OS" != "windows" ]]; then
                chmod +x "$binary_path"
                log "Binary made executable"
            fi
        fi
    else
        log_error "Build failed. Check $LOG_FILE for details."
        exit 1
    fi
}

# Verify build
verify_build() {
    log "Verifying build..."
    
    local binary_path="$BIN_DIR/${BINARY_NAME}"
    if [[ "$TARGET_OS" == "windows" ]]; then
        binary_path="${binary_path}.exe"
    fi
    
    if [[ ! -f "$binary_path" ]]; then
        log_error "Binary not found: $binary_path"
        exit 1
    fi
    
    # Test if binary can show version (cross-platform check)
    if [[ "$TARGET_OS" == "$(go env GOOS)" && "$TARGET_ARCH" == "$(go env GOARCH)" ]]; then
        log "Testing binary execution..."
        if "$binary_path" --version 2>/dev/null || "$binary_path" -version 2>/dev/null; then
            log_success "Binary verification completed"
        else
            log_warn "Binary exists but version check failed (this might be normal)"
        fi
    else
        log "Cross-compiled binary - skipping execution test"
    fi
    
    log_success "Build verification completed"
}

# Cleanup function
cleanup() {
    log "Cleaning up temporary files..."
    # Add any necessary cleanup here
    log_success "Cleanup completed"
}

# Main execution
main() {
    log "Analytics Service Build Script Started"
    log "Timestamp: $(date)"
    
    # Initialize log file
    echo "=== Analytics Service Build Log ===" > "$LOG_FILE"
    echo "Started: $(date)" >> "$LOG_FILE"
    
    # Parse arguments
    parse_args "$@"
    
    # Setup trap for cleanup
    trap cleanup EXIT
    
    # Execute build steps
    validate_prerequisites
    setup_build_env
    build_service
    verify_build
    
    log_success "Analytics Service build completed successfully!"
    log "Build artifacts:"
    log "  Binary: $BIN_DIR/${BINARY_NAME}$([ "$TARGET_OS" == "windows" ] && echo ".exe")"
    log "  Log: $LOG_FILE"
    
    echo
    log_success "Build Summary:"
    log "  Version: $VERSION"
    log "  Target: ${TARGET_OS}/${TARGET_ARCH}"
    log "  Debug: $DEBUG_MODE"
    log "  Location: $BIN_DIR"
}

# Execute main function if script is run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
