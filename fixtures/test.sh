#!/usr/bin/env bash
set -euo pipefail

FIXTURES_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$FIXTURES_DIR/.." && pwd)"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

GOOS="$(go env GOOS)"
GOARCH="$(go env GOARCH)"
go build -o "$TMP_DIR/nara" "$REPO_ROOT"

for fixture_test in "$FIXTURES_DIR"/*/test.sh; do
  echo "running $(basename "$(dirname "$fixture_test")")"
  NARA_BIN="$TMP_DIR/nara" bash "$fixture_test"
done
