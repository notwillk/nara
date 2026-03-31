#!/usr/bin/env bash
set -euo pipefail

V="${VERSION:-latest}"
if [ "$V" != "latest" ] && [ "$V" != "current" ]; then
  V="${V#v}"
fi

# Install nara via upstream installer, honoring VERSION.
# Clear the feature-provided VERSION env for latest/current so the upstream
# installer resolves the real latest release instead of treating "latest" as a tag.
if [ "$V" = "latest" ] || [ "$V" = "current" ]; then
  curl -fsSL https://raw.githubusercontent.com/notwillk/nara/main/install.sh | env -u VERSION bash
else
  curl -fsSL https://raw.githubusercontent.com/notwillk/nara/main/install.sh | VERSION="v${V}" bash
fi

nara --help >/dev/null
