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

nara --config "$FIXTURE_DIR/nara.yaml" list schemas >/dev/null
nara --config "$FIXTURE_DIR/nara.yaml" validate "$FIXTURE_DIR"/entries/*.note.yaml
nara --config "$FIXTURE_DIR/nara.yaml" compile "$FIXTURE_DIR"/entries/*.note.yaml --format json --out "$TMP_DIR/out.json"

python3 - "$TMP_DIR/out.json" <<'PY'
import json
import sys

with open(sys.argv[1], 'r', encoding='utf-8') as f:
    data = json.load(f)

assert len(data) == 1, data
assert data[0]["title"] == "Hello", data
assert data[0]["body"] == "A minimal single-file note fixture.", data
assert data[0]["tags"] == ["intro", "yaml"], data
PY
