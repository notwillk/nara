#!/usr/bin/env bash
set -euo pipefail

FEATURE_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
TMP_PROJECT="$(mktemp -d)"
trap 'rm -rf "$TMP_PROJECT"' EXIT

case "$(uname -m)" in
  x86_64|amd64) GOARCH="amd64" ;;
  aarch64|arm64) GOARCH="arm64" ;;
  *)
    echo "unsupported architecture for feature test: $(uname -m)" >&2
    exit 1
    ;;
esac

cp -R "$FEATURE_ROOT/src" "$TMP_PROJECT/src"
mkdir -p "$TMP_PROJECT/test/nara"

GOOS=linux GOARCH="$GOARCH" CGO_ENABLED=0 go build -o "$TMP_PROJECT/src/nara/nara" "$REPO_ROOT"

cat >"$TMP_PROJECT/src/nara/install.sh" <<'EOS'
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
install -m 0755 "$SCRIPT_DIR/nara" /usr/local/bin/nara
nara --help >/dev/null
EOS
chmod +x "$TMP_PROJECT/src/nara/install.sh"

cat >"$TMP_PROJECT/test/nara/test.sh" <<'EOS'
#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=/dev/null
source dev-container-features-test-lib

check "nara available" nara --help

reportResults
EOS
chmod +x "$TMP_PROJECT/test/nara/test.sh"

devcontainer features test \
  --project-folder "$TMP_PROJECT" \
  --features nara \
  --base-image mcr.microsoft.com/devcontainers/base:ubuntu
