#!/bin/bash
# Rebuild samfile-audit, regenerate JSONL across the corpus,
# ingest into findings.db, mine reports. Idempotent.
#
# Requires ~/sam-corpus/.venv with pandas, scikit-learn, mlxtend.
# Bootstrap it with:
#   python3 -m venv ~/sam-corpus/.venv
#   ~/sam-corpus/.venv/bin/pip install pandas scikit-learn mlxtend

set -euo pipefail

REPO="${HOME}/git/samfile"
CORPUS="${HOME}/sam-corpus"
PY="${CORPUS}/.venv/bin/python"

if [ ! -x "$PY" ]; then
    echo "error: $PY not found. Bootstrap the venv:" >&2
    echo "  python3 -m venv ${CORPUS}/.venv && ${CORPUS}/.venv/bin/pip install pandas scikit-learn mlxtend" >&2
    exit 1
fi

cd "$REPO"
echo "==> building samfile-audit"
go build -o "${CORPUS}/samfile-audit" ./cmd/samfile

mkdir -p "${CORPUS}/outputs-jsonl" "${CORPUS}/analyses"
rm -f "${CORPUS}/outputs-jsonl"/*.jsonl

echo "==> running samfile-audit verify --format jsonl across the corpus"
count=0
for disk in "${CORPUS}/disks"/*.mgt; do
    [ -f "$disk" ] || continue
    name=$(basename "$disk" .mgt)
    "${CORPUS}/samfile-audit" verify --format jsonl -i "$disk" \
        > "${CORPUS}/outputs-jsonl/${name}.jsonl" 2>/dev/null || true
    count=$((count + 1))
done
echo "ran on $count disks"

echo "==> ingesting JSONL into ${CORPUS}/findings.db"
"$PY" "${REPO}/tools/audit/ingest.py"

echo "==> mining reports into ${CORPUS}/analyses/"
"$PY" "${REPO}/tools/audit/mine.py"

echo "==> done. reports in ${CORPUS}/analyses/"
ls -la "${CORPUS}/analyses/"
