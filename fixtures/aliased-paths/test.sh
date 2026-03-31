#!/usr/bin/env bash
set -euo pipefail

FIXTURE_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$FIXTURE_DIR/../.." && pwd)"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

nara() {
  if [ -n "${NARA_BIN:-}" ]; then
    "$NARA_BIN" "$@"
  else
    (cd "$REPO_ROOT" && go run . "$@")
  fi
}

LIST_OUTPUT="$(nara --config "$FIXTURE_DIR/nara.yaml" list entries)"
printf '%s\n' "$LIST_OUTPUT" | grep -F $'acme\tsupplier\tshared/acme.supplier.yaml' >/dev/null
printf '%s\n' "$LIST_OUTPUT" | grep -F $'widget\tproduct\tentries/widget.product.yaml' >/dev/null

nara --config "$FIXTURE_DIR/nara.yaml" validate "$FIXTURE_DIR"/entries/*.product.yaml
nara --config "$FIXTURE_DIR/nara.yaml" compile "$FIXTURE_DIR"/entries/*.product.yaml --format sqlite --out "$TMP_DIR/out.db"

python3 - "$TMP_DIR/out.db" <<'PY'
import sqlite3
import sys

conn = sqlite3.connect(sys.argv[1])
cur = conn.cursor()
cur.execute('select id, schema from entities order by id')
rows = cur.fetchall()
assert rows == [('acme', 'supplier'), ('widget', 'product')], rows
cur.execute('select from_id, field, to_id from edges')
edge_rows = cur.fetchall()
assert edge_rows == [('widget', '$.supplier', 'acme')], edge_rows
conn.close()
PY
