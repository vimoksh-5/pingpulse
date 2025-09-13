#!/usr/bin/env bash
set -euo pipefail

# Local sweep: secrets, deps, fs misconfig, email/phone
echo "==> Running Gitleaks (staged + working tree)"
if ! command -v gitleaks >/dev/null 2>&1; then
  echo "Installing gitleaks..."
  curl -sSL https://raw.githubusercontent.com/gitleaks/gitleaks/master/install.sh | bash
  sudo mv gitleaks /usr/local/bin/ || true
fi
gitleaks detect --redact -c .gitleaks.toml

echo "==> OSV-Scanner (deps)"
if ! command -v osv-scanner >/dev/null 2>&1; then
  go install github.com/google/osv-scanner/cmd/osv-scanner@latest || true
fi
osv-scanner --recursive . || true

echo "==> Trivy filesystem (vuln+secret+config)"
if ! command -v trivy >/dev/null 2>&1; then
  curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -
  sudo mv ./bin/trivy /usr/local/bin/ || true
fi
trivy fs --scanners vuln,secret,config --ignore-unfixed --exit-code 0 .

echo "==> Quick PII grep (emails/phones)"
bash scripts/find_emails.sh
echo "Done."
