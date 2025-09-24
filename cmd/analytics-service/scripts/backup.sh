#!/usr/bin/env bash
# ==============================================================================
# backup.sh — Backup seguro do banco de dados do Analytics Service
# ==============================================================================
# Descrição:
#   Executa backup do banco de dados relacionado ao analytics-service com opções
#   de destino (local ou remoto), compressão, logs, timestamp e limpeza de 
#   backups antigos. Segue boas práticas de Bash e segurança.
#
# Uso:
#   ./scripts/backup.sh [opções]
#
# Opções:
#   -d, --dest <local|s3|scp>     Destino do backup (padrão: local)
#   -o, --output-dir <dir>        Diretório de saída para backups locais (padrão: ./backups)
#   -b, --bucket <s3://bucket/prefix>  Bucket/prefixo S3 (quando dest=s3)
#   -r, --remote <user@host:/path>     Destino remoto SCP (quando dest=scp)
#   -c, --compression <gz|zstd|none>    Tipo de compressão (padrão: gz)
#   -k, --keep-days <N>           Dias para manter backups (padrão: 7)
#   -n, --name <basename>         Nome base do arquivo (padrão: analytics)
#   --no-clean                    Não remover backups antigos
#   --dry-run                     Simula execução sem efetuar upload/movimento
#   -h, --help                    Exibe ajuda
#
# Variáveis de ambiente esperadas (podem vir de .env.scripts):
#   DB_HOST, DB_PORT, DB_NAME, DB_USER
#   DB_PASSWORD (recomendado via variável/secret; não hardcode)
#   AWS_PROFILE / AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY (para S3)
#
# Exemplos:
#   ./scripts/backup.sh -d local -o ./backups -c gz -k 14
#   ./scripts/backup.sh -d s3 -b s3://my-bucket/analytics/backups -c zstd
#   ./scripts/backup.sh -d scp -r user@server:/data/backups -c gz
# ==============================================================================

set -euo pipefail
IFS=$'\n\t'

# ------------------------------- Cores/Logs ----------------------------------
RED="\033[31m"; GREEN="\033[32m"; YELLOW="\033[33m"; BLUE="\033[34m"; BOLD="\033[1m"; RESET="\033[0m"
log_info()  { echo -e "${BLUE}[INFO]${RESET} $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${RESET} $*"; }
log_error() { echo -e "${RED}[ERROR]${RESET} $*" 1>&2; }
log_ok()    { echo -e "${GREEN}[OK]${RESET} $*"; }

print_help() {
  sed -n '1,80p' "$0"
}

# --------------------------- Defaults e argumentos ---------------------------
DESTINATION="local"          # local|s3|scp
OUTPUT_DIR="./backups"
S3_BUCKET=""
SCP_TARGET=""
COMPRESSION="gz"            # gz|zstd|none
KEEP_DAYS=7
BASENAME="analytics"
DO_CLEAN=1
DRY_RUN=0

# ---------------------------- Validação de deps ------------------------------
require_cmd() {
  command -v "$1" >/dev/null 2>&1 || { log_error "Comando requerido não encontrado: $1"; exit 1; }
}

# ------------------------------ Parse args -----------------------------------
while [[ $# -gt 0 ]]; do
  case "$1" in
    -d|--dest) DESTINATION="${2:-}"; shift 2 ;;
    -o|--output-dir) OUTPUT_DIR="${2:-}"; shift 2 ;;
    -b|--bucket) S3_BUCKET="${2:-}"; shift 2 ;;
    -r|--remote) SCP_TARGET="${2:-}"; shift 2 ;;
    -c|--compression) COMPRESSION="${2:-}"; shift 2 ;;
    -k|--keep-days) KEEP_DAYS="${2:-}"; shift 2 ;;
    -n|--name) BASENAME="${2:-}"; shift 2 ;;
    --no-clean) DO_CLEAN=0; shift ;;
    --dry-run) DRY_RUN=1; shift ;;
    -h|--help) print_help; exit 0 ;;
    *) log_error "Argumento desconhecido: $1"; print_help; exit 1 ;;
  esac
done

# --------------------------- Validação de ambiente ---------------------------
: "${DB_HOST:?Defina DB_HOST}"
: "${DB_PORT:?Defina DB_PORT}"
: "${DB_NAME:?Defina DB_NAME}"
: "${DB_USER:?Defina DB_USER}"
# DB_PASSWORD é opcional se .pgpass for utilizado

# Ferramentas necessárias
require_cmd date
require_cmd mktemp
require_cmd psql
require_cmd pg_dump
case "$COMPRESSION" in
  gz)   require_cmd gzip ;;
  zstd) require_cmd zstd || { log_error "zstd requerido para compressão zstd"; exit 1; } ;;
  none) : ;;
  *) log_error "Compressão inválida: $COMPRESSION"; exit 1 ;;
end

# Para destinos
case "$DESTINATION" in
  local) : ;;
  s3)    require_cmd aws ;;
  scp)   require_cmd scp ;;
  *) log_error "Destino inválido: $DESTINATION"; exit 1 ;;
end

# ----------------------------- Preparação backup -----------------------------
TS="$(date -u +%Y%m%dT%H%M%SZ)"
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT
FILE_BASE="${BASENAME}_${DB_NAME}_${TS}.sql"
DUMP_PATH="$TMPDIR/$FILE_BASE"

# String de conexão segura para pg_dump/psql
PGPASSFILE_TMP=""
if [[ -n "${DB_PASSWORD:-}" ]]; then
  PGPASSFILE_TMP="$TMPDIR/.pgpass"
  chmod 600 "$PGPASSFILE_TMP"
  echo "${DB_HOST}:${DB_PORT}:${DB_NAME}:${DB_USER}:${DB_PASSWORD}" > "$PGPASSFILE_TMP"
  export PGPASSFILE="$PGPASSFILE_TMP"
fi

export PGHOST="$DB_HOST"
export PGPORT="$DB_PORT"
export PGDATABASE="$DB_NAME"
export PGUSER="$DB_USER"

log_info "Iniciando backup do banco ${BOLD}$DB_NAME${RESET} em $TS (dest=$DESTINATION, comp=$COMPRESSION)"

# Check conexão rápida
if ! psql -c 'SELECT 1;' >/dev/null 2>&1; then
  log_error "Não foi possível conectar ao banco. Verifique credenciais e rede."
  exit 1
fi
log_ok "Conexão com banco validada."

# ------------------------------- Dump do banco -------------------------------
log_info "Gerando dump lógico com pg_dump..."
pg_dump --no-owner --format=plain --encoding=UTF8 \
  --quote-all-identifiers \
  --file "$DUMP_PATH"
log_ok "Dump gerado em $DUMP_PATH"

# ------------------------------- Compressão ---------------------------------
OUT_PATH="$DUMP_PATH"
case "$COMPRESSION" in
  gz)
    log_info "Comprimindo com gzip..."
    gzip -9 "$DUMP_PATH"
    OUT_PATH="${DUMP_PATH}.gz"
    ;;
  zstd)
    log_info "Comprimindo com zstd..."
    zstd -19 --quiet "$DUMP_PATH"
    OUT_PATH="${DUMP_PATH}.zst"
    ;;
  none)
    log_warn "Sem compressão selecionada."
    ;;
esac
log_ok "Arquivo final: $OUT_PATH"

# ------------------------------- Entrega -------------------------------------
case "$DESTINATION" in
  local)
    mkdir -p "$OUTPUT_DIR"
    TARGET_PATH="$OUTPUT_DIR/$(basename "$OUT_PATH")"
    if [[ "$DRY_RUN" -eq 1 ]]; then
      log_info "[DRY-RUN] Copiaria $OUT_PATH -> $TARGET_PATH"
    else
      cp "$OUT_PATH" "$TARGET_PATH"
      log_ok "Backup salvo em $TARGET_PATH"
    fi
    ;;
  s3)
    if [[ -z "$S3_BUCKET" ]]; then
      log_error "Informe --bucket s3://bucket/prefix ao usar dest=s3"
      exit 1
    fi
    TARGET_PATH="$S3_BUCKET/$(basename "$OUT_PATH")"
    if [[ "$DRY_RUN" -eq 1 ]]; then
      log_info "[DRY-RUN] Enviaria $OUT_PATH -> $TARGET_PATH"
    else
      aws s3 cp "$OUT_PATH" "$TARGET_PATH" --only-show-errors
      log_ok "Backup enviado para $TARGET_PATH"
    fi
    ;;
  scp)
    if [[ -z "$SCP_TARGET" ]]; then
      log_error "Informe --remote user@host:/path ao usar dest=scp"
      exit 1
    fi
    if [[ "$DRY_RUN" -eq 1 ]]; then
      log_info "[DRY-RUN] Enviaria $OUT_PATH -> $SCP_TARGET"
    else
      scp "$OUT_PATH" "$SCP_TARGET"
      log_ok "Backup transferido para $SCP_TARGET"
    fi
    ;;
esac

# -------------------------- Limpeza de backups antigos -----------------------
if [[ "$DO_CLEAN" -eq 1 && "$DESTINATION" == "local" ]]; then
  if [[ -d "$OUTPUT_DIR" ]]; then
    log_info "Limpando backups locais com mais de $KEEP_DAYS dias em $OUTPUT_DIR..."
    if [[ "$DRY_RUN" -eq 1 ]]; then
      log_info "[DRY-RUN] Removeria arquivos antigos correspondentes a ${BASENAME}_*."
    else
      find "$OUTPUT_DIR" -type f -name "${BASENAME}_*.sql*" -mtime +"$KEEP_DAYS" -print -delete || true
      log_ok "Limpeza concluída."
    fi
  fi
elif [[ "$DO_CLEAN" -eq 1 && "$DESTINATION" == "s3" ]]; then
  if [[ -n "$S3_BUCKET" ]]; then
    log_info "Limpando backups S3 com mais de $KEEP_DAYS dias em $S3_BUCKET..."
    if [[ "$DRY_RUN" -eq 1 ]]; then
      log_info "[DRY-RUN] Listaria e deletaria objetos antigos via aws s3 rm."
    else
      # Remove objetos antigos por data de modificação
      aws s3 ls "$S3_BUCKET/" | awk -v days="$KEEP_DAYS" -v base="$BASENAME" '
        BEGIN{ cmd="date -u +%s"; cmd | getline now; close(cmd); cutoff=now-(days*24*3600) }
        { date=$1 " " $2; file=$4; if (file ~ (base "_")) {
            cmd = "date -u -d \"" date "\" +%s"; cmd | getline ts; close(cmd);
            if (ts < cutoff) print file;
          }
        }' | while read -r key; do
          aws s3 rm "$S3_BUCKET/$key" --only-show-errors || true
        done
      log_ok "Limpeza em S3 concluída."
    fi
  fi
else
  log_info "Limpeza automática desativada ou não suportada para este destino."
fi

log_ok "Backup finalizado com sucesso."
