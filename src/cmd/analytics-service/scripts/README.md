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
---
## 🛡️ Auditoria Periódica / Periodic Auditing
PT: Diretrizes e exemplos para auditorias mensais/trimestrais de scripts e infraestrutura.
EN: Guidelines and examples for monthly/quarterly audits of scripts and infrastructure.

1) Checklist de auditoria mensal/trimestral (linhas gerais) / Monthly/Quarterly audit checklist (high-level)
- PT:
  - [ ] Revisar mudanças em scripts (git log, diffs, owners) e cobertura de testes
  - [ ] Rodar linters e análise estática (ShellCheck) e corrigir findings críticos
  - [ ] Revalidar variáveis sensíveis e escopos (mínimo privilégio) em CI/CD e cloud
  - [ ] Regerar SBOM e re-escanear imagens/dep. (syft/grype/trivy/snyk)
  - [ ] Verificar política de assinatura/verificação de imagens (cosign/OPA/Gatekeeper)
  - [ ] Validar configurações de Docker/K8s (rootless, readOnlyRootFilesystem, seccomp)
  - [ ] Executar benchmarks e hardening (kube-bench, kube-hunter, CIS Docker)
  - [ ] Revisar logs/alertas, SLAs de backup/restore e testes de recuperação
  - [ ] Validar rotação de chaves/tokens e expiração de credenciais
  - ■ Evidências anexadas no relatório e tickets abertos para pendências
- EN:
  - [ ] Review script changes (git log, diffs, owners) and test coverage
  - [ ] Run linters and static analysis (ShellCheck) and fix critical findings
  - [ ] Revalidate sensitive vars and scopes (least privilege) in CI/CD and cloud
  - [ ] Regenerate SBOM and rescan images/deps (syft/grype/trivy/snyk)
  - [ ] Verify image signing/enforcement (cosign/OPA/Gatekeeper)
  - [ ] Validate Docker/K8s configs (rootless, readOnlyRootFilesystem, seccomp)
  - [ ] Run benchmarks/hardening (kube-bench, kube-hunter, CIS Docker)
  - [ ] Review logs/alerts, backup/restore SLAs and recovery drills
  - [ ] Validate key/token rotation and credential expiration
  - ■ Attach evidence in report and open follow-up tickets

2) Procedimentos recomendados / Recommended procedures
- Docker:
  - PT: Imagens mínimas (distroless/alpine), rootless, USER não-root, drop de capabilities, scan com Trivy/Snyk, política de resource limits, não usar latest.
  - EN: Minimal images (distroless/alpine), rootless, non-root USER, drop capabilities, scan with Trivy/Snyk, resource limits policy, avoid latest tags.
- Kubernetes:
  - PT: PodSecurity/PSA enforced, NetworkPolicies, readOnlyRootFilesystem, seccomp/apparmor, requests/limits, liveness/readiness, secrets via K8s/Vault, RBAC least privilege.
  - EN: Enforce PodSecurity/PSA, NetworkPolicies, readOnlyRootFilesystem, seccomp/apparmor, requests/limits, probes, secrets via K8s/Vault, least-privilege RBAC.
- Cloud:
  - PT: Contas separadas por ambiente, IAM com least privilege, rotação de chaves, criptografia at-rest/in-transit, CloudTrail/Audit Logs, S3/Object lock e versionamento.
  - EN: Separate accounts by env, least-privilege IAM, key rotation, at-rest/in-transit encryption, CloudTrail/Audit logs, object lock/versioning.

3) Ferramentas automatizadas / Automated tools
- Scripts: ShellCheck, shfmt, Bandit (para Python utilitário), Semgrep rules.
- Containers: Trivy, Grype, Syft (SBOM), Snyk.
- Kubernetes: kube-bench, kube-hunter, Polaris, Kubesec.
- Cloud/Configs: tfsec, Checkov, Terrascan, OpenSCAP, AWS Config/Config Rules.
- Supply-chain: cosign (assinatura/verificação), Sigstore policy-controller, osv-scanner.

4) Modelo de relatório resumido / Summary report template
```
Título/Title: Auditoria Scripts & Infra - <YYYY-MM> (Mensal/Trimestral)
Escopo/Scope: Repositórios, imagens, clusters, contas cloud
Resumo Executivo/Executive Summary: <3-5 bullets com principais riscos e ações>
Metodologia/Methodology: Ferramentas e checks executados
Achados/Findings:
  - Criticidade/Severity: High | Medium | Low
  - Descrição/Description: <texto>
  - Evidência/Evidence: <arquivo/link>
  - Ação/Action: <fix/owner/data>
Métricas/Metrics: #CVEs High, %pods com R/O FS, %jobs com sucesso, tempo MTTR
Riscos Abertos/Open Risks: <lista>
Plano de Ação/Action Plan: <tarefas, prazos, responsáveis>
Aprovação/Approval: <assinado por>
```

5) Agendamento via CI/CD / Scheduling via CI/CD
- GitHub Actions (cron mensal):
```
name: Monthly Audit
on:
  schedule:
    - cron: '0 3 1 * *'   # 03:00 UTC todo dia 1
jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Security Tools
        run: |
          ./scripts/ci/audit.sh
      - name: Upload report artifact
        uses: actions/upload-artifact@v4
        with:
          name: audit-report
          path: reports/audit-*.md
```
- Jenkins (pipeline nightly/weekly):
```
pipeline {
  triggers { cron('H H(2-4) 1 * *') } // mensal, janela 02:00-04:00
  agent any
  stages {
    stage('Run Audit') { steps { sh './scripts/ci/audit.sh' } }
    stage('Archive Report') { steps { archiveArtifacts artifacts: 'reports/audit-*.md', fingerprint: true } }
  }
}
```

6) Exemplo visual de agendamento + Slack/Email / Visual scheduling + Slack/Email integration
- Script de auditoria (exemplo):
```
#!/usr/bin/env bash
set -euo pipefail
mkdir -p reports
REPORT="reports/audit-$(date -u +%F).md"
{
  echo "# Audit $(date -u +%F)";
  echo "## Tools";
  shellcheck -V || true
  echo "## Results";
  trivy fs --quiet --exit-code 0 . || true
} > "$REPORT"
# Slack (via webhook)
[ -n "${SLACK_WEBHOOK_URL:-}" ] && curl -s -X POST -H 'Content-type: application/json' \
  --data "{\"text\":\"Audit finished: $REPORT\"}" "$SLACK_WEBHOOK_URL" || true
```
- GitHub Actions passo Slack/Email:
```
- name: Notify Slack
  if: always()
  uses: rtCamp/action-slack-notify@v2
  env:
    SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK_URL }}
    SLACK_MESSAGE: "Audit finished. Artifacts attached."

- name: Send Email (SMTP)
  if: always()
  uses: dawidd6/action-send-mail@v3
  with:
    server_address: smtp.example.com
    server_port: 587
    username: ${{ secrets.SMTP_USER }}
    password: ${{ secrets.SMTP_PASS }}
    subject: Audit finished
    to: security@example.com
    from: ci@example.com
    attachments: reports/audit-*.md
```

Notas / Notes:
- Execute auditorias em ambientes efêmeros para evitar impacto.
- Centralize relatórios (S3/Artifacts) e aplique retenção/ACLs.
- Defina SLOs para correção de findings com ownership claro.
