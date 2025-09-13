#!/usr/bin/env bash
set -euo pipefail

if ! command -v gitleaks >/dev/null 2>&1; then
  echo "Installing gitleaks..."
  curl -sSL https://raw.githubusercontent.com/gitleaks/gitleaks/master/install.sh | bash
  sudo mv gitleaks /usr/local/bin/ || true
fi

echo "==> Scanning full git history with Gitleaks"
gitleaks detect --redact -c .gitleaks.toml --log-opts="--all" -v
