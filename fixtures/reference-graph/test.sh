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

nara --config "$FIXTURE_DIR/nara.yaml" lint
nara --config "$FIXTURE_DIR/nara.yaml" compile "$FIXTURE_DIR"/entries/*.person.yaml --format json --out "$TMP_DIR/out.json"

python3 - "$TMP_DIR/out.json" <<'PY'
import json
import sys

with open(sys.argv[1], 'r', encoding='utf-8') as f:
    data = json.load(f)

assert len(data) == 4, data
team = next(item for item in data if item["name"] == "Platform Team")
assert [report["name"] for report in team["reports"]] == ["Riley Lead", "Quinn Dev"], team
lead = next(item for item in data if item["name"] == "Riley Lead")
assert lead["manager"]["name"] == "Casey Chief", lead
PY
