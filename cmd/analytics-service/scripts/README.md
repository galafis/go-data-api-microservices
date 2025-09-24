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
- Troubleshooting / Troubleshooting
- Manutenção / Maintenance
- Créditos / Credits
- ETL: Exemplos Práticos e Visuais / ETL: Practical and Visual Examples
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

## ⚙️ Variáveis de ambiente / Environment variables
Exemplo / Example:

```bash
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

## 🤖 Integração CI/CD / CI/CD integration
GitHub Actions (trecho) / snippet:

```yaml
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

```groovy
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

## ✅ Boas práticas / Best practices
- set -e, set -u, set -o pipefail
- Flags --help e validação de inputs
- Sem segredos hardcoded; use variáveis de ambiente/secret manager
- Logs claros e timestamps; níveis de log
- Idempotência: reentrância segura nos scripts
- Checks de pré-requisito (docker, kubectl, go)

Segurança / Security:
- Sanitizar entradas, princípio do menor privilégio, não expor tokens em logs

Observabilidade / Observability:
- Coletar métricas, logs estruturados, códigos de retorno consistentes

## 🐛 Troubleshooting
Permissão negada / Permission denied:

```bash
chmod +x scripts/*.sh
find scripts/ -name "*.sh" -exec chmod +x {} \;
```

Dependências ausentes / Missing dependencies:

```bash
./scripts/deps.sh
./scripts/health-check.sh --prereq
```

Falhas de build / Build failures:

```bash
./scripts/clean.sh
go mod download && go mod tidy
./scripts/build.sh --debug
```

Deploy falhou / Deploy failed:

```bash
./scripts/logs.sh --since=1h
./scripts/rollback.sh --to=$PREV_TAG
```

Banco de dados / Database:

```bash
./scripts/backup.sh --output backup_$(date +%F).sql
./scripts/restore.sh --input backup.sql
```

## 🛠️ Manutenção / Maintenance
- Versionar scripts (semver) e manter CHANGELOG
- Testes em CI para cada alteração de script
- Atualizar dependências periodicamente
- Monitorar duração dos jobs e otimizar cache

## 👤 Créditos / Credits
- Autor/Author: Gabriel Demetrios Lafis
- Parte do ecossistema Go Data API Microservices / Part of the Go Data API Microservices ecosystem
- Feedback e melhorias são bem-vindos! / Feedback and improvements are welcome!

---
## 🧱 ETL: Exemplos Práticos e Visuais / ETL: Practical and Visual Examples

Este bloco acrescenta exemplos didáticos de pipeline ETL usando o script test_analytics_etl.sh, incluindo: automação com parâmetros, mock de variáveis de ambiente, outputs simulados, e instruções de logs/monitoramento. Seção separada em dois blocos: exemplos de execução e exemplos visuais para copy-paste.

This section adds educational ETL pipeline examples using test_analytics_etl.sh, including: parameterized automation, environment variable mocks, simulated outputs, and logging/monitoring guidance. Split into two blocks: runnable examples and visual copy-paste snippets.

### 1) Execução guiada / Guided execution

```bash
# PT: Executa ETL completo com parâmetros de data e modo verbose
# EN: Run full ETL with date parameters and verbose mode
./cmd/analytics-service/scripts/test_analytics_etl.sh \
  --source=s3://raw-bucket/daily/ \
  --target=postgresql://$DB_USER@${DB_HOST}:${DB_PORT}/${DB_NAME} \
  --date=$(date -u +%F) \
  --stages=extract,transform,load \
  --concurrency=4 \
  --verbose
```

```bash
# PT: Executa apenas transformação e carga, filtrando partição
# EN: Run only transform and load, filtering a partition
./cmd/analytics-service/scripts/test_analytics_etl.sh \
  --stages=transform,load \
  --partition=dt=2025-09-24/region=BR \
  --fail-fast
```

```bash
# PT: Execução dry-run para validar dependências e plano de execução
# EN: Dry-run to validate dependencies and execution plan
./cmd/analytics-service/scripts/test_analytics_etl.sh --dry-run --plan --verbose
```

### 2) Mock de variáveis de ambiente / Environment variable mocks

```bash
# PT: Mock controlado para rodar localmente sem tocar em serviços reais
# EN: Controlled mock to run locally without touching real services
export ETL_MOCK_MODE=true
export ANALYTICS_ENV="test"
export DB_HOST="127.0.0.1"
export DB_PORT="5432"
export DB_NAME="analytics_test"
export DB_USER="analytics_user"
export S3_ENDPOINT="http://localhost:4566"   # localstack
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"

./cmd/analytics-service/scripts/test_analytics_etl.sh \
  --source=s3://mock-raw/daily/ \
  --target=postgresql://$DB_USER@${DB_HOST}:${DB_PORT}/${DB_NAME} \
  --stages=extract,transform,load \
  --verbose
```

Dicas / Tips:
- PT: Combine ETL_MOCK_MODE com fixtures em ./testdata/ para entradas determinísticas.
- EN: Combine ETL_MOCK_MODE with ./testdata/ fixtures for deterministic inputs.

### 3) Outputs simulados / Simulated outputs

```text
[2025-09-24T10:00:00Z] [INFO] ETL start {date="2025-09-24", stages=[extract,transform,load]}
[2025-09-24T10:00:02Z] [INFO] extract: 24 files discovered, 24 queued
[2025-09-24T10:00:07Z] [INFO] transform: 24 -> 24 records normalized (schema v3)
[2025-09-24T10:00:09Z] [INFO] load: batch=8, inserted=24, updated=0, upsert_key="event_id"
[2025-09-24T10:00:10Z] [METRIC] etl_duration_seconds=10.2 stage="all"
[2025-09-24T10:00:10Z] [METRIC] etl_records_total=24 labels={stage="load"}
[2025-09-24T10:00:10Z] [INFO] ETL success ✔
```

Falhas comuns / Common failures:

```text
[ERROR] extract: S3 403 AccessDenied (verifique credenciais / check credentials)
[WARN ] transform: 3 records dropped by validation (schema mismatch)
[ERROR] load: failed pq: relation "events" does not exist (rodar migrate.sh)
```

### 4) Logs e monitoramento / Logging and monitoring

- PT: Os scripts emitem logs estruturados (JSON opcional via LOG_FORMAT=json). Use logs.sh para coleta e filtros por estágio.
- EN: Scripts emit structured logs (JSON optional via LOG_FORMAT=json). Use logs.sh for collection and stage filters.

Exemplos / Examples:

```bash
# PT: Seguir logs da última execução ETL
# EN: Tail logs of the last ETL run
./cmd/analytics-service/scripts/logs.sh --component=etl --since=2h --follow
```

```bash
# PT: Exportar métricas para Prometheus via textfile collector
# EN: Export metrics to Prometheus via textfile collector
./cmd/analytics-service/scripts/metrics.sh --component=etl --output=/var/lib/node_exporter/textfile/etl.prom
```

Boas práticas / Best practices:
- PT: Inclua request_id/trace_id em cada estágio; log de progresso a cada N registros; nunca logar dados sensíveis.
- EN: Include request_id/trace_id per stage; progress logging every N records; never log sensitive data.

Alertas e SLAs / Alerts and SLAs:
- PT: Configure alertas para etl_duration_seconds alto e taxa de erro por estágio (>1%).
- EN: Configure alerts for high etl_duration_seconds and per-stage error rate (>1%).

---
## 🧩 Exemplos visuais (copy-paste) / Visual examples (copy-paste)

Copie e cole conforme necessário para seu caso de uso.
Copy and paste as needed for your use case.

```bash
# FULL RUN
./cmd/analytics-service/scripts/test_analytics_etl.sh \
  --source=s3://raw-bucket/daily/ \
  --target=postgresql://$DB_USER@${DB_HOST}:${DB_PORT}/${DB_NAME} \
  --date=$(date -u +%F) \
  --stages=extract,transform,load \
  --concurrency=4 \
  --verbose
```

```bash
# PARTIAL RUN (TRANSFORM+LOAD)
./cmd/analytics-service/scripts/test_analytics_etl.sh \
  --stages=transform,load \
  --partition=dt=2025-09-24/region=BR \
  --fail-fast
```

```bash
# DRY RUN + PLAN
./cmd/analytics-service/scripts/test_analytics_etl.sh --dry-run --plan --verbose
```

```bash
# MOCKED ENV
export ETL_MOCK_MODE=true
export ANALYTICS_ENV="test"
export DB_HOST="127.0.0.1"; export DB_PORT="5432"; export DB_NAME="analytics_test"; export DB_USER="analytics_user"
export S3_ENDPOINT="http://localhost:4566"; export AWS_ACCESS_KEY_ID="test"; export AWS_SECRET_ACCESS_KEY="test"
./cmd/analytics-service/scripts/test_analytics_etl.sh --stages=extract,transform,load --verbose \
  --source=s3://mock-raw/daily/ \
  --target=postgresql://$DB_USER@${DB_HOST}:${DB_PORT}/${DB_NAME}
```

```bash
# LOGS + METRICS
./cmd/analytics-service/scripts/logs.sh --component=etl --since=2
