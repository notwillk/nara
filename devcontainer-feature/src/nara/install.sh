#!/usr/bin/env bash
set -euo pipefail

V="${VERSION:-latest}"
if [ "$V" != "latest" ] && [ "$V" != "current" ]; then
  V="${V#v}"
fi

# Install nara via upstream installer, honoring VERSION
if [ "$V" = "latest" ] || [ "$V" = "current" ]; then
  curl -fsSL https://raw.githubusercontent.com/notwillk/nara/main/install.sh | bash
else
  curl -fsSL https://raw.githubusercontent.com/notwillk/nara/main/install.sh | VERSION="v${V}" bash
fi

nara --help >/dev/null || true
