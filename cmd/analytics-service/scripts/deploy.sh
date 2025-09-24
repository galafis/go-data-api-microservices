#!/usr/bin/env bash
# deploy.sh - Automated deployment script for Analytics Service
# Author: Comet Assistant
# Description: Performs deployments to dev, staging, and prod with confirmations,
#              dry-run support for production, rollback capability, and
#              comprehensive logging and safety guards.
# Usage:
#   ./deploy.sh <dev|staging|prod> [--confirm] [--dry-run] [--rollback]
# Options:
#   --confirm     Require manual confirmation (mandatory for prod unless --dry-run)
#   --dry-run     Simulate actions without making changes (allowed/encouraged for prod)
#   --rollback    Roll back to previous version/tag (requires DEPLOY_ROLLBACK_REF)
#   --help        Show help
# Notes:
#   - Follows Bash/DevOps best practices (set -euo pipefail, traps, logging)
#   - Validates required environment variables before proceeding
#   - Integrates with Docker and Kubernetes if configured

set -euo pipefail
set -o errtrace

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
PROJECT_ROOT="$(cd "${SERVICE_ROOT}/../.." && pwd)"

# Defaults and flags
ENVIRONMENT=""
REQUIRE_CONFIRM=false
DRY_RUN=false
DO_ROLLBACK=false

# Colors
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'

# Logging helpers
log_info()    { echo -e "${BLUE}INFO:${NC} $*"; }
log_warn()    { echo -e "${YELLOW}WARN:${NC} $*"; }
log_error()   { echo -e "${RED}ERROR:${NC} $*" 1>&2; }
log_success() { echo -e "${GREEN}SUCCESS:${NC} $*"; }

cleanup() {
  log_info "Cleaning up temporary resources (if any)..."
}

error_handler() {
  local exit_code=$1 line=$2
  log_error "Deployment failed with exit code ${exit_code} at line ${line}"
  cleanup
  exit "${exit_code}"
}
trap 'error_handler $? $LINENO' ERR
trap cleanup EXIT

show_help() {
  cat <<EOF
Usage: $0 <dev|staging|prod> [--confirm] [--dry-run] [--rollback]

Environments:
  dev       Deploy to development
  staging   Deploy to staging
  prod      Deploy to production (requires --confirm or --dry-run)

Options:
  --confirm   Ask for explicit confirmation before executing changes
  --dry-run   Print actions without executing them
  --rollback  Roll back to previous version (requires DEPLOY_ROLLBACK_REF)
  --help      Show this help

Environment variables (examples / expectations):
  # Common
  SERVICE_NAME=analytics-service
  VERSION?=auto-derived-from-git
  DOCKER_REGISTRY=registry.example.com/namespace
  DOCKER_IMAGE="${DOCKER_REGISTRY:-your-registry.com}/analytics-service"

  # Kubernetes
  KUBECONFIG=~/.kube/config
  K8S_NAMESPACE=analytics
  K8S_DEPLOYMENT=analytics-service

  # Rollback
  DEPLOY_ROLLBACK_REF=<tag|sha|image-tag>
EOF
}

confirm() {
  local prompt_msg=${1:-"Proceed? (y/N): "}
  read -r -p "${prompt_msg}" ans || true
  case "${ans:-}" in
    y|Y|yes|YES) return 0 ;;
    *) log_warn "User cancelled."; return 1 ;;
  esac
}

parse_args() {
  if [[ $# -lt 1 ]]; then
    show_help; exit 1
  fi
  case "$1" in
    dev|development) ENVIRONMENT=dev ;;
    staging)         ENVIRONMENT=staging ;;
    prod|production) ENVIRONMENT=prod ;;
    --help|-h)       show_help; exit 0 ;;
    *) log_error "Unknown environment: $1"; show_help; exit 1 ;;
  esac
  shift || true
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --confirm)  REQUIRE_CONFIRM=true ;;
      --dry-run)  DRY_RUN=true ;;
      --rollback) DO_ROLLBACK=true ;;
      --help|-h)  show_help; exit 0 ;;
      *) log_error "Unknown option: $1"; show_help; exit 1 ;;
    esac
    shift || true
  done
}

require_env() {
  local name=$1
  if [[ -z "${!name:-}" ]]; then
    log_error "Missing required environment variable: ${name}"
    exit 1
  fi
}

check_prerequisites() {
  log_info "Checking prerequisites..."
  command -v git >/dev/null 2>&1 || { log_error "git not found"; exit 1; }
  command -v docker >/dev/null 2>&1 || log_warn "docker not found - skipping docker steps"
  command -v kubectl >/dev/null 2>&1 || log_warn "kubectl not found - skipping Kubernetes steps"

  # Required for all envs
  export SERVICE_NAME=${SERVICE_NAME:-analytics-service}

  # Optional defaults
  export DOCKER_REGISTRY=${DOCKER_REGISTRY:-your-registry.com}
  export DOCKER_IMAGE=${DOCKER_IMAGE:-${DOCKER_REGISTRY}/${SERVICE_NAME}}
  export K8S_NAMESPACE=${K8S_NAMESPACE:-analytics}
  export K8S_DEPLOYMENT=${K8S_DEPLOYMENT:-analytics-service}

  # Versioning
  if [[ -z "${VERSION:-}" ]]; then
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
      VERSION=$(git describe --tags --always --dirty 2>/dev/null || git rev-parse --short HEAD)
    else
      VERSION="unknown"
    fi
  fi
  export VERSION

  log_info "Environment: ${ENVIRONMENT} | Version: ${VERSION}"
  log_info "Docker image: ${DOCKER_IMAGE}:${VERSION}"
}

ensure_prod_safety() {
  if [[ "${ENVIRONMENT}" == "prod" ]]; then
    if [[ "${DRY_RUN}" != true && "${REQUIRE_CONFIRM}" != true ]]; then
      log_error "Production deploy requires --confirm or --dry-run"
      exit 1
    fi
  fi
}

run_cmd() {
  if [[ "${DRY_RUN}" == true ]]; then
    echo "+ $*"
  else
    eval "$*"
  fi
}

build_and_push_image() {
  if ! command -v docker >/dev/null 2>&1; then
    log_warn "Docker not available, skipping image build/push"
    return
  fi
  log_info "Building Docker image ${DOCKER_IMAGE}:${VERSION}"
  run_cmd "docker build -t ${DOCKER_IMAGE}:${VERSION} -f ${PROJECT_ROOT}/Dockerfile ${PROJECT_ROOT}"
  log_info "Pushing Docker image ${DOCKER_IMAGE}:${VERSION}"
  run_cmd "docker push ${DOCKER_IMAGE}:${VERSION}"
}

k8s_deploy() {
  if ! command -v kubectl >/dev/null 2>&1; then
    log_warn "kubectl not available, skipping Kubernetes deploy"
    return
  fi
  require_env K8S_NAMESPACE
  require_env K8S_DEPLOYMENT

  local image="${DOCKER_IMAGE}:${VERSION}"
  log_info "Setting image on deployment/${K8S_DEPLOYMENT} to ${image} in namespace ${K8S_NAMESPACE}"
  run_cmd "kubectl -n ${K8S_NAMESPACE} set image deployment/${K8S_DEPLOYMENT} ${SERVICE_NAME}=${image}"

  log_info "Waiting for rollout to complete"
  run_cmd "kubectl -n ${K8S_NAMESPACE} rollout status deployment/${K8S_DEPLOYMENT} --timeout=300s"
}

k8s_rollback() {
  if ! command -v kubectl >/dev/null 2>&1; then
    log_warn "kubectl not available, cannot rollback"
    return
  fi
  require_env K8S_NAMESPACE
  require_env K8S_DEPLOYMENT

  if [[ -n "${DEPLOY_ROLLBACK_REF:-}" ]]; then
    local image="${DOCKER_IMAGE}:${DEPLOY_ROLLBACK_REF}"
    log_warn "Rolling back by setting image to ${image}"
    run_cmd "kubectl -n ${K8S_NAMESPACE} set image deployment/${K8S_DEPLOYMENT} ${SERVICE_NAME}=${image}"
  else
    log_warn "DEPLOY_ROLLBACK_REF not set, using kubectl rollout undo"
    run_cmd "kubectl -n ${K8S_NAMESPACE} rollout undo deployment/${K8S_DEPLOYMENT}"
  fi

  log_info "Waiting for rollback to complete"
  run_cmd "kubectl -n ${K8S_NAMESPACE} rollout status deployment/${K8S_DEPLOYMENT} --timeout=300s"
}

pre_deploy_checks() {
  log_info "Running pre-deploy checks..."
  # Example: ensure go.mod exists at project root
  if [[ ! -f "${PROJECT_ROOT}/go.mod" ]]; then
    log_error "go.mod not found at project root (${PROJECT_ROOT})"
    exit 1
  fi
}

post_deploy_validation() {
  log_info "Running post-deploy validation..."
  if command -v kubectl >/dev/null 2>&1; then
    run_cmd "kubectl -n ${K8S_NAMESPACE} get pods -l app=${SERVICE_NAME}"
  fi
}

deploy_flow() {
  pre_deploy_checks

  case "${ENVIRONMENT}" in
    dev)
      log_info "Deploying to development"
      build_and_push_image
      k8s_deploy
      ;;
    staging)
      log_info "Deploying to staging"
      build_and_push_image
      k8s_deploy
      ;;
    prod)
      log_warn "Preparing production deployment"
      ensure_prod_safety
      if [[ "${REQUIRE_CONFIRM}" == true && "${DRY_RUN}" != true ]]; then
        confirm "Proceed with PRODUCTION deployment of ${VERSION}? (y/N): " || exit 1
      fi
      build_and_push_image
      k8s_deploy
      ;;
    *)
      log_error "Unknown environment: ${ENVIRONMENT}"; exit 1 ;;
  esac

  post_deploy_validation
  log_success "Deployment flow completed for ${ENVIRONMENT}"
}

main() {
  echo "======================================"
  echo "  Analytics Service Deployment Script"
  echo "======================================"

  parse_args "$@"
  check_prerequisites

  if [[ "${DO_ROLLBACK}" == true ]]; then
    log_warn "Rollback requested"
    ensure_prod_safety
    if [[ "${REQUIRE_CONFIRM}" == true && "${DRY_RUN}" != true ]]; then
      confirm "Proceed with ROLLBACK on ${ENVIRONMENT}? (y/N): " || exit 1
    fi
    k8s_rollback
    log_success "Rollback completed"
    return 0
  fi

  deploy_flow

  if [[ "${DRY_RUN}" == true ]]; then
    log_warn "Dry run completed. No changes were made."
  fi
}

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  show_help; exit 0
fi

main "$@"
