#!/usr/bin/env bash
# migrate.sh — Database migration management for Analytics Service
#
# Description:
#   Wrapper script to manage database migrations using common Go-friendly tools
#   (prefers golang-migrate if available, then goose, then sql-migrate). Provides
#   commands: up, down, status, create. Includes logging, environment validation,
#   safe defaults, and helpful usage examples.
#
# Best practices:
#   - set -euo pipefail for safer bash
#   - Explicit env var validation and .env loading support
#   - Consistent, structured logging with timestamps
#   - Dry-run and confirmation options for destructive operations
#
# Usage:
#   ./scripts/migrate.sh up [--steps N] [--dry-run]
#   ./scripts/migrate.sh down [--steps N] [--confirm]
#   ./scripts/migrate.sh status
#   ./scripts/migrate.sh create <name> [--type sql|go]
#   ./scripts/migrate.sh --help
#
# Examples:
#   DB_NAME=analytics_dev ./scripts/migrate.sh up
#   ./scripts/migrate.sh down --steps 1 --confirm
#   ./scripts/migrate.sh create add_user_preferences --type sql
#
# Environment (can be set or provided via .env files):
#   DB_HOST (default: localhost)
#   DB_PORT (default: 5432)
#   DB_NAME (required)
#   DB_USER (default: $USER)
#   DB_PASSWORD (optional; use secret store when possible)
#   DB_SSLMODE (default: disable)
#   MIGRATIONS_DIR (default: cmd/analytics-service/migrations)
#   DB_SCHEMA (optional, defaults to public)
#
# Notes:
#   - Ensure the migration tool of choice is installed and on PATH.
#   - For goose/sql-migrate, this script expects conventional folder layouts.
#   - For golang-migrate, this uses file:// scheme for MIGRATIONS_DIR.

set -euo pipefail
IFS=$'\n\t'

# --- Globals & defaults ---
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
ROOT_DIR="$(cd -- "${SCRIPT_DIR}/../../.." &>/dev/null && pwd)"
SERVICE_DIR="${ROOT_DIR}/cmd/analytics-service"
DEFAULT_MIGRATIONS_DIR="${SERVICE_DIR}/migrations"
LOG_TS_FORMAT="+%Y-%m-%dT%H:%M:%S%z"

# Allow overriding via env
MIGRATIONS_DIR="${MIGRATIONS_DIR:-${DEFAULT_MIGRATIONS_DIR}}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-}"
DB_USER="${DB_USER:-${USER:-analytics}}"
DB_PASSWORD="${DB_PASSWORD:-}"
DB_SSLMODE="${DB_SSLMODE:-disable}"
DB_SCHEMA="${DB_SCHEMA:-public}"

DEBUG_MODE="${DEBUG:-0}"
DRY_RUN=0
STEPS=""
CONFIRM=0
CREATE_TYPE="sql"

# --- Logging helpers ---
log()       { printf '%s [INFO ] %s\n'  "$(date "${LOG_TS_FORMAT}")" "$*"; }
warn()      { printf '%s [WARN ] %s\n'  "$(date "${LOG_TS_FORMAT}")" "$*" 1>&2; }
error()     { printf '%s [ERROR] %s\n' "$(date "${LOG_TS_FORMAT}")" "$*" 1>&2; }
debug()     { [[ "${DEBUG_MODE}" == "1" ]] && printf '%s [DEBUG] %s\n' "$(date "${LOG_TS_FORMAT}")" "$*"; }

# --- Utility: command existence ---
has() { command -v "$1" >/dev/null 2>&1; }

# --- Load .env files (non-fatal) ---
load_env() {
  local files=("${ROOT_DIR}/.env" "${SERVICE_DIR}/.env" "${SERVICE_DIR}/.env.scripts")
  for f in "${files[@]}"; do
    if [[ -f "$f" ]]; then
      debug "Loading env file: $f"
      # shellcheck disable=SC2046
      set -a; source "$f"; set +a
    fi
  done
}

# --- Print usage ---
usage() {
  cat <<'USAGE'
Usage:
  migrate.sh up [--steps N] [--dry-run]
  migrate.sh down [--steps N] [--confirm]
  migrate.sh status
  migrate.sh create <name> [--type sql|go]
  migrate.sh --help | -h

Options:
  --steps N    Number of migration steps to apply/rollback (default: all for up, 1 for down)
  --dry-run    Show what would run without executing (wrapper-level; tools may vary)
  --confirm    Required for destructive operations (down)
  --type       Migration file type (for create). Default: sql

Environment:
  DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD, DB_SSLMODE, DB_SCHEMA, MIGRATIONS_DIR

Examples:
  DB_NAME=analytics_dev ./scripts/migrate.sh up
  ./scripts/migrate.sh down --steps 1 --confirm
  ./scripts/migrate.sh status
  ./scripts/migrate.sh create add_user_preferences --type sql
USAGE
}

# --- Ensure prerequisites and env ---
validate_env() {
  [[ -z "${DB_NAME}" ]] && { error "DB_NAME is required"; exit 1; }
  if [[ ! -d "${MIGRATIONS_DIR}" ]]; then
    warn "MIGRATIONS_DIR not found: ${MIGRATIONS_DIR}. Creating it."
    mkdir -p "${MIGRATIONS_DIR}"
  fi
}

# Build a Postgres connection string for tools that accept DSN/URL
pg_url() {
  local auth
  if [[ -n "${DB_PASSWORD}" ]]; then
    auth="${DB_USER}:${DB_PASSWORD}"
  else
    auth="${DB_USER}"
  fi
  local qp="sslmode=${DB_SSLMODE}"
  if [[ -n "${DB_SCHEMA}" ]]; then
    qp+="&search_path=${DB_SCHEMA}"
  fi
  printf 'postgres://%s@%s:%s/%s?%s' "${auth}" "${DB_HOST}" "${DB_PORT}" "${DB_NAME}" "${qp}"
}

# Determine backend tool priority
select_tool() {
  if has migrate; then echo "migrate"; return; fi
  if has goose; then echo "goose"; return; fi
  if has sql-migrate; then echo "sql-migrate"; return; fi
  echo ""; return
}

confirm_destructive() {
  if [[ ${CONFIRM} -ne 1 ]]; then
    error "This operation is destructive. Re-run with --confirm to proceed."
    exit 1
  fi
}

# --- Parsers ---
parse_global_flags() {
  local arg
  while [[ $# -gt 0 ]]; do
    arg="$1"; shift || true
    case "$arg" in
      --help|-h) usage; exit 0;;
      --dry-run) DRY_RUN=1;;
      --steps)   STEPS="${1:-}"; shift || true;;
      --confirm) CONFIRM=1;;
      --type)    CREATE_TYPE="${1:-sql}"; shift || true;;
      *)         set -- "$arg" "$@"; break;;
    esac
  done
  echo "$@"
}

# --- Actions ---
action_up() {
  local tool="$1"
  local steps_opt=""
  [[ -n "${STEPS}" ]] && steps_opt="${STEPS}"
  log "Applying migrations (tool=${tool}, steps=${steps_opt:-all})"
  if [[ ${DRY_RUN} -eq 1 ]]; then warn "Dry-run requested. No changes will be applied."; fi
  case "${tool}" in
    migrate)
      local db="$(pg_url)"
      local cmd=(migrate -path "${MIGRATIONS_DIR}" -database "${db}" up)
      [[ -n "${STEPS}" ]] && cmd+=("${STEPS}")
      [[ ${DRY_RUN} -eq 1 ]] && { printf 'DRY-RUN: %q ' "${cmd[@]}"; echo; return; }
      "${cmd[@]}"
      ;;
    goose)
      local conn="user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} sslmode=${DB_SSLMODE} search_path=${DB_SCHEMA}"
      local cmd=(goose -dir "${MIGRATIONS_DIR}" postgres "${conn}" up)
      if [[ -n "${STEPS}" ]]; then
        # goose lacks direct up N; run up-by-one N times
        for (( i=0; i<STEPS; i++ )); do
          [[ ${DRY_RUN} -eq 1 ]] && { printf 'DRY-RUN[%d]: %q ' "$i" "${cmd[@]}"; echo; continue; }
          "${cmd[@]}"
        done
        return
      fi
      [[ ${DRY_RUN} -eq 1 ]] && { printf 'DRY-RUN: %q ' "${cmd[@]}"; echo; return; }
      "${cmd[@]}"
      ;;
    sql-migrate)
      local cmd=(sql-migrate up)
      [[ -n "${STEPS}" ]] && cmd+=("-limit" "${STEPS}")
      [[ ${DRY_RUN} -eq 1 ]] && cmd+=("-dryrun")
      "${cmd[@]}"
      ;;
    *) error "No supported migration tool found (install migrate, goose, or sql-migrate)"; exit 1;;
  esac
}

action_down() {
  local tool="$1"
  confirm_destructive
  local steps_opt="${STEPS:-1}"
  log "Rolling back migrations (tool=${tool}, steps=${steps_opt})"
  case "${tool}" in
    migrate)
      local db="$(pg_url)"
      local cmd=(migrate -path "${MIGRATIONS_DIR}" -database "${db}" down "${steps_opt}")
      [[ ${DRY_RUN} -eq 1 ]] && { printf 'DRY-RUN: %q ' "${cmd[@]}"; echo; return; }
      "${cmd[@]}"
      ;;
    goose)
      local conn="user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} sslmode=${DB_SSLMODE} search_path=${DB_SCHEMA}"
      for (( i=0; i<steps_opt; i++ )); do
        local cmd=(goose -dir "${MIGRATIONS_DIR}" postgres "${conn}" down)
        [[ ${DRY_RUN} -eq 1 ]] && { printf 'DRY-RUN[%d]: %q ' "$i" "${cmd[@]}"; echo; continue; }
        "${cmd[@]}"
      done
      ;;
    sql-migrate)
      local cmd=(sql-migrate down -limit "${steps_opt}")
      [[ ${DRY_RUN} -eq 1 ]] && cmd+=("-dryrun")
      "${cmd[@]}"
      ;;
    *) error "No supported migration tool found (install migrate, goose, or sql-migrate)"; exit 1;;
  esac
}

action_status() {
  local tool="$1"
  log "Showing migration status (tool=${tool})"
  case "${tool}" in
    migrate)
      local db="$(pg_url)"
      migrate -path "${MIGRATIONS_DIR}" -database "${db}" version || true
      ;;
    goose)
      local conn="user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} sslmode=${DB_SSLMODE} search_path=${DB_SCHEMA}"
      goose -dir "${MIGRATIONS_DIR}" postgres "${conn}" status || true
      ;;
    sql-migrate)
      sql-migrate status || true
      ;;
    *) error "No supported migration tool found"; exit 1;;
  esac
}

action_create() {
  local tool="$1"; shift || true
  local name="${1:-}"
  [[ -z "${name}" ]] && { error "Missing migration name. Usage: migrate.sh create <name> [--type sql|go]"; exit 1; }
  log "Creating new migration (tool=${tool}, type=${CREATE_TYPE}, name=${name})"
  mkdir -p "${MIGRATIONS_DIR}"
  case "${tool}" in
    migrate)
      local ts
      ts="$(date +%Y%m%d%H%M%S)"
      case "${CREATE_TYPE}" in
        sql)
          local up="${MIGRATIONS_DIR}/${ts}_${name}.up.sql"
          local down="${MIGRATIONS_DIR}/${ts}_${name}.down.sql"
          printf '-- +migrate Up\n-- SQL statements for up migration\n' >"${up}"
          printf '-- +migrate Down\n-- SQL statements for down migration\n' >"${down}"
          log "Created ${up} and ${down}"
          ;;
        go)
          local up="${MIGRATIONS_DIR}/${ts}_${name}.go"
          cat >"${up}" <<EOF
package migrations

// TODO: implement migration ${name}
EOF
          log "Created ${up} (skeleton)"
          ;;
        *) error "Unsupported --type: ${CREATE_TYPE}"; exit 1;;
      esac
      ;;
    goose)
      local cmd=(goose -dir "${MIGRATIONS_DIR}" create "${name}" "${CREATE_TYPE}")
      [[ ${DRY_RUN} -eq 1 ]] && { printf 'DRY-RUN: %q ' "${cmd[@]}"; echo; return; }
      "${cmd[@]}"
      ;;
    sql-migrate)
      local cwd
      cwd="$(pwd)"
      cd "${MIGRATIONS_DIR}"
      local cmd=(sql-migrate new "${name}")
      [[ ${DRY_RUN} -eq 1 ]] && { printf 'DRY-RUN: (in %s) %q ' "${MIGRATIONS_DIR}" "${cmd[@]}"; echo; cd "${cwd}"; return; }
      "${cmd[@]}"
      cd "${cwd}"
      ;;
    *) error "No supported migration tool found"; exit 1;;
  esac
}

main() {
  load_env || true
  local cmd="${1:-}"; shift || true
  if [[ -z "${cmd}" ]]; then usage; exit 1; fi
  # Parse flags following the command
  # shellcheck disable=SC2206
  local rest=( $(parse_global_flags "$@") )
  : "${rest[@]:-}"
  validate_env
  local tool
  tool="$(select_tool)"
  if [[ -z "${tool}" ]]; then
    error "No migration tool found. Install one of: golang-migrate (migrate), goose, sql-migrate"
    exit 1
  fi
  log "Using migration tool: ${tool}"
  log "Migrations dir: ${MIGRATIONS_DIR} | DB: ${DB_HOST}:${DB_PORT}/${DB_NAME} schema=${DB_SCHEMA}"
  case "${cmd}" in
    up)     action_up     "${tool}" ;;
    down)   action_down   "${tool}" ;;
    status) action_status "${tool}" ;;
    create) action_create "${tool}" "$@" ;;
    --help|-h) usage ;;
    *) error "Unknown command: ${cmd}"; usage; exit 1;;
  esac
}

main "$@"
