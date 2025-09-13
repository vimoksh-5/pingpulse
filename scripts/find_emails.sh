#!/usr/bin/env bash
set -euo pipefail

# 🔧 Add your real emails to highlight
EMAILS=("YOUR_EMAIL_1" "YOUR_EMAIL_2" "YOUR_WORK_EMAIL")

echo "==> Searching for any emails in tracked text files"
grep -RIn --exclude-dir={.git,node_modules,dist,build,.next,.venv,venv} \
  --include='*.{go,js,ts,tsx,jsx,py,md,env,yml,yaml,ini,txt,json,tf,sh}' \
  -E "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}" . || true

for e in "${EMAILS[@]}"; do
  if [[ "$e" != "YOUR_EMAIL_1" ]]; then
    echo "==> Highlighting occurrences of $e"
    git grep -n -- "$e" || true
  fi
done

echo "==> Searching for Indian phone formats"
grep -RIn --exclude-dir={.git,node_modules,dist,build,.next,.venv,venv} \
  -E "(\\+?91[\\s-]?)?[6-9][0-9]{9}" . || true
