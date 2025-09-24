# 📦 Analytics Service Scripts — Guia Rápido / Quick Start
Autor/Author: Gabriel Demetrios Lafis
Este README foi projetado para onboarding rápido, com exemplos práticos, tabelas resumidas, dicas de CI/CD, troubleshooting e boas práticas. Versão bilíngue: Português e Inglês, lado a lado quando aplicável.
This README is designed for rapid onboarding, with practical examples, summary tables, CI/CD tips, troubleshooting, and best practices. Bilingual: Portuguese and English, side-by-side when applicable.
---
## 🔗 Índice / Table of Contents
- Visão Geral / Overview
- Tabela de Scripts / Scripts Matrix
- Como usar por ambiente / Environment-based usage
- Exemplos práticos / Practical examples
- Variáveis de ambiente / Environment variables
- Integração CI/CD / CI/CD integration
- Boas práticas / Best practices
- Segurança de Scripts / Script Security
- Troubleshooting / Troubleshooting
- Manutenção / Maintenance
- Créditos / Credits
---
## 🧭 Visão Geral / Overview
- PT: Scripts de automação para desenvolvimento, testes, build, deploy, banco de dados e manutenção do Analytics Service.
- EN: Automation scripts for development, testing, build, deployment, database, and maintenance for the Analytics Service.
---
## 📚 Tabela de Scripts / Scripts Matrix
| Categoria / Category | Script | Descrição (PT) | Description (EN) |
|---|---|---|---|
| Build & Dev | build.sh | Compila o serviço com otimizações | Builds service with optimizations |
| Build & Dev | dev.sh | Sobe servidor com hot reload | Starts dev server with hot reload |
| Build & Dev | clean.sh | Limpa artefatos de build | Cleans build artifacts |
| Build & Dev | deps.sh | Instala/atualiza dependências | Installs/updates dependencies |
| Test & Quality | test.sh | Executa suíte completa com cobertura | Runs full test suite with coverage |
| Test & Quality | test-unit.sh | Executa testes unitários | Runs unit tests |
| Test & Quality | test-integration.sh | Executa testes de integração | Runs integration tests |
| Test & Quality | lint.sh | Lint e formatação | Linting and formatting |
| Test & Quality | security.sh | Scan de vulnerabilidades | Security vulnerability scanning |
| Deployment | deploy.sh | Deploy automatizado por ambiente | Automated environment deployment |
| Deployment | docker-build.sh | Build de imagem Docker | Docker image build |
| Deployment | k8s-deploy.sh | Deploy no Kubernetes | Kubernetes deployment |
| Deployment | rollback.sh | Rollback de versão | Version rollback |
| Database & Data | migrate.sh | Migrações de schema | Schema migrations |
| Database & Data | seed.sh | Seed de dados de teste | Test data seeding |
| Database & Data | backup.sh | Backup do banco | Database backup |
| Database & Data | restore.sh | Restore do banco | Database restore |
| Monitoring & Ops | health-check.sh | Verifica saúde do serviço | Service health check |
| Monitoring & Ops | logs.sh | Coleta/análise de logs | Log collection/analysis |
| Monitoring & Ops | metrics.sh | Coleta de métricas | Metrics collection |
| Monitoring & Ops | cleanup.sh | Limpeza de manutenção | Maintenance cleanup |
> Dica/Tip: Todos os scripts aceitam --help quando disponível. Many scripts support --help.
---
## 🏗️ Como usar por ambiente / Environment-based usage
- Desenvolvimento / Development:
  - PT: Rodar deps, testes e servidor de desenvolvimento.
  - EN: Run deps, tests, and development server.
  - bash:
    - ./scripts/deps.sh
    - ./scripts/test.sh --verbose
    - ./scripts/dev.sh
- Staging:
  - PT: Build, imagem Docker e deploy em staging.
  - EN: Build, Docker image, and deploy to staging.
  - bash:
    - ./scripts/build.sh --version=$(git rev-parse --short HEAD)
    - ./scripts/docker-build.sh
    - ./scripts/deploy.sh staging
- Produção / Production:
  - PT: Confirmar e registrar rollout.
  - EN: Confirm and record rollout.
  - bash:
    - ./scripts/test.sh
    - ./scripts/build.sh --version=$TAG
    - ./scripts/deploy.sh prod --confirm
Rollback:
- PT: Reverter rapidamente em caso de falha.
- EN: Quick revert on failure.
- bash:
  - ./scripts/rollback.sh --to=$PREV_TAG
---
## 🧪 Exemplos práticos / Practical examples
- Build multiplataforma / Cross-platform build:
  - ./scripts/build.sh --os=linux --arch=amd64 --version=v1.2.3
- Testes com cobertura HTML / HTML coverage:
  - ./scripts/test.sh --html-coverage
- Migrações / Migrations:
  - ./scripts/migrate.sh up
  - ./scripts/migrate.sh create add_user_preferences
- Health check pré-deploy / Pre-deploy health check:
  - ./scripts/health-check.sh --prereq
- Debug detalhado / Verbose debug:
  - DEBUG=1 ./scripts/build.sh
  - ./scripts/build.sh --debug
---
## ⚙️ Variáveis de ambiente / Environment variables
Exemplo / Example:
```
# Service
export ANALYTICS_ENV="development"
export ANALYTICS_PORT="8083"
export LOG_LEVEL="info"

# Database
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_NAME="analytics_dev"
export DB_USER="analytics_user"

# Docker
export DOCKER_REGISTRY="your-registry.com"
export DOCKER_TAG="latest"

# Kubernetes
export KUBECONFIG="~/.kube/config"
export K8S_NAMESPACE="analytics"
```
Arquivos de config / Config files:
- .env.scripts, config/build.yaml, config/deploy.yaml, config/test.yaml
---
## 🤖 Integração CI/CD / CI/CD integration
GitHub Actions (trecho) / snippet:
```
name: Analytics Service CI/CD
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Tests
        run: ./cmd/analytics-service/scripts/test.sh
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Service
        run: ./cmd/analytics-service/scripts/build.sh
      - name: Build Docker Image
        run: ./cmd/analytics-service/scripts/docker-build.sh
```
Jenkins (trecho) / snippet:
```
pipeline {
  agent any
  stages {
    stage('Test') { steps { sh './cmd/analytics-service/scripts/test.sh' } }
    stage('Build') { steps { sh './cmd/analytics-service/scripts/build.sh'; sh './cmd/analytics-service/scripts/docker-build.sh' } }
    stage('Deploy') { when { branch 'main' } steps { sh './cmd/analytics-service/scripts/deploy.sh prod' } }
  }
}
```
Dicas / Tips:
- PT: Use matrizes (matrix) para múltiplas plataformas; armazene DOCKER_REGISTRY/DOCKER_TOKEN como secrets; gere SBOM (syft) e varredura (grype).
- EN: Use matrix builds; store DOCKER_REGISTRY/DOCKER_TOKEN as secrets; generate SBOM (syft) and scan (grype).
---
## ✅ Boas práticas / Best practices
- set -e, set -u, set -o pipefail
- Flags --help e validação de inputs
- Sem segredos hardcoded; use variáveis de ambiente/secret manager
- Logs claros e timestamps; níveis de log
- Idempotência: reentrância segura nos scripts
- Checks de pré-requisito (docker, kubectl, go)

---
## 🔒 Segurança de Scripts / Script Security
PT: Diretrizes práticas para proteger scripts de automação. EN: Practical guidelines to secure automation scripts.

1) Sanitização de entradas / Input sanitization
- PT: Nunca executar entradas sem validação; use whitelist de flags/valores e regex. 
- EN: Never execute unvalidated input; use whitelist of flags/values and regex.
Exemplo (bash seguro / safe bash):
```bash
set -euo pipefail

ALLOWED_ENVS="dev|staging|prod"
ENVIRONMENT="${1:-dev}"
if ! [[ "$ENVIRONMENT" =~ ^(${ALLOWED_ENVS})$ ]]; then
  echo "[ERR] invalid env: $ENVIRONMENT" >&2; exit 2
fi

# sanitize numeric
CONCURRENCY="${CONCURRENCY:-4}"
if ! [[ "$CONCURRENCY" =~ ^[0-9]{1,2}$ ]]; then
  echo "[ERR] invalid concurrency" >&2; exit 2
fi
```

2) Variáveis de ambiente seguras / Secure environment variables
- PT: Carregue de .env.scripts com set -a; nunca echo segredos; redija logs. 
- EN: Load from .env.scripts using set -a; never echo secrets; redact logs.
Exemplo:
```bash
set -euo pipefail
# load env safely
if [[ -f .env.scripts ]]; then set -a; source .env.scripts; set +a; fi

# redact secrets in logs
log_redact() { sed -E "s/(${DB_PASS}|${API_TOKEN})/[REDACTED]/g"; }

# usage example
some_command --token "$API_TOKEN" 2>&1 | log_redact
```

3) Proteção contra comandos perigosos / Guard against dangerous commands
- PT: Use set -o noclobber; proteja rm, sed, find com checagens de path. 
- EN: Use set -o noclobber; guard rm, sed, find with path checks.
Exemplo safe-rm:
```bash
safe_rm() {
  local target="$1"
  [[ -z "$target" ]] && { echo "empty target" >&2; return 2; }
  [[ "$target" == "/" || "$target" == "/*" ]] && { echo "refusing to delete root" >&2; return 3; }
  [[ ! -e "$target" ]] && { echo "not found: $target" >&2; return 4; }
  rm -rf -- "$target"
}
```
Guard for kubectl context/namespace:
```bash
require_k8s() {
  local ns="${K8S_NAMESPACE:-}"
  local ctx
  ctx=$(kubectl config current-context 2>/dev/null || true)
  [[ -z "$ctx" || -z "$ns" ]] && { echo "k8s context/namespace missing" >&2; exit 2; }
  echo "Using k8s: context=$ctx ns=$ns"
}
```

4) Modo dry-run e confirmação / Dry-run and confirmation
- PT: Forneça --dry-run e --confirm antes de ações destrutivas. 
- EN: Provide --dry-run and --confirm before destructive actions.
Exemplo:
```bash
DRY_RUN=0; CONFIRM=0
for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=1;;
    --confirm) CONFIRM=1;;
  esac
done

run_or_echo() { [[ $DRY_RUN -eq 1 ]] && echo "+ $*" || eval "$*"; }
confirm_or_exit() { [[ $CONFIRM -eq 1 ]] || { echo "--confirm required"; exit 3; }; }

confirm_or_exit
run_or_echo kubectl apply -f deploy.yaml
```

5) Auditoria de logs / Log auditing
- PT: Envie logs para arquivo e syslog, retenção rotacionada, checksum (sha256). 
- EN: Send logs to file and syslog, rotated retention, checksum (sha256).
Exemplo:
```bash
LOG_DIR="./logs"; mkdir -p "$LOG_DIR"
LOG_FILE="$LOG_DIR/$(date -u +%F)_deploy.log"
{
  date -u +"%FT%TZ"; set -x
  ./scripts/deploy.sh staging
} &>"$LOG_FILE"
sha256sum "$LOG_FILE" >"$LOG_FILE.sha256"
logger -t analytics-deploy -- "deploy finished: $(tail -n1 "$LOG_FILE")"
```

6) Checklist de compliance / Compliance checklist
- PT:
  - [ ] Segredos via secret manager (GH Secrets, SOPS, Vault)
  - [ ] SBOM gerado e verificado (syft/grype)
  - [ ] Imagens assinadas (cosign) e policy (OPA/Gatekeeper)
  - [ ] Princípio do menor privilégio (IAM/K8s RBAC)
  - [ ] Artefatos com checksum/assinatura (sha256/cosign)
  - [ ] Backups e testes de restore regulares
  - [ ] Logs com request_id e redacted
- EN:
  - [ ] Secrets via secret manager (GH Secrets, SOPS, Vault)
  - [ ] SBOM generated and scanned (syft/grype)
  - [ ] Images signed (cosign) and policy (OPA/Gatekeeper)
  - [ ] Least privilege (IAM/K8s RBAC)
  - [ ] Artifacts checksummed/signed (sha256/cosign)
  - [ ] Backups and periodic restore tests
  - [ ] Logs with request_id and redacted

7) Referências / References
- OWASP Cheatsheets: https://cheatsheetseries.owasp.org/
- OWASP Top 10: https://owasp.org/www-project-top-ten/
- CIS Benchmarks: https://www.cisecurity.org/cis-benchmarks
- Shell Script Security: https://github.com/koalaman/shellcheck/wiki/ShellCheck

---
## 🐛 Troubleshooting
- Permissão negada / Permission denied:
  - chmod +x scripts/*.sh
  - find scripts/ -name "*.sh" -exec chmod +x {} \;
- Dependências ausentes / Missing dependencies:
  - ./scripts/deps.sh
  - ./scripts/health-check.sh --prereq
- Falhas de build / Build failures:
  - ./scripts/clean.sh
  - go mod download && go mod tidy
  - ./scripts/build.sh --debug
- Deploy falhou / Deploy failed:
  - ./scripts/logs.sh --since=1h
  - ./scripts/rollback.sh --to=$PREV_TAG
- Banco de dados / Database:
  - ./scripts/backup.sh --output backup_$(date +%F).sql
  - ./scripts/restore.sh --input backup.sql
---
## 🛠️ Manutenção / Maintenance
- Versionar scripts (semver) e manter CHANGELOG
- Testes em CI para cada alteração de script
- Atualizar dependências periodicamente
- Monitorar duração dos jobs e otimizar cache
---
## 👤 Créditos / Credits
- Autor/Author: Gabriel Demetrios Lafis
- Parte do ecossistema Go Data API Microservices / Part of the Go Data API Microservices ecosystem
- Feedback e melhorias são bem-vindos! / Feedback and improvements are welcome!
