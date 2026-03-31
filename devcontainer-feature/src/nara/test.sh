#!/usr/bin/env bash
set -euo pipefail

FEATURE_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TMP_PROJECT="$(mktemp -d)"
trap 'rm -rf "$TMP_PROJECT"' EXIT

cp -R "$FEATURE_ROOT/src" "$TMP_PROJECT/src"
mkdir -p "$TMP_PROJECT/test/nara"

cat >"$TMP_PROJECT/test/nara/test.sh" <<'EOS'
#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=/dev/null
source dev-container-features-test-lib

check "nara available" nara --help

reportResults
EOS
chmod +x "$TMP_PROJECT/test/nara/test.sh"

devcontainer features test   --project-folder "$TMP_PROJECT"   --features nara   --base-image mcr.microsoft.com/devcontainers/base:ubuntu
