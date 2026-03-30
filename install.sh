#!/usr/bin/env bash
# Install nara from GitHub Releases (PRD §13).
# Usage:
#   curl -sL https://raw.githubusercontent.com/notwillk/nara/main/install.sh | bash
# Optional:
#   VERSION=v1.2.3 bash install.sh
#   GITHUB_REPO=owner/nara bash install.sh

set -euo pipefail

REPO="${GITHUB_REPO:-notwillk/nara}"
BINARY="${BINARY_NAME:-nara}"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "install.sh: required command not found: $1" >&2
    exit 1
  }
}

fetch() {
  local url="$1"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "$url"
  else
    echo "install.sh: need curl or wget" >&2
    exit 1
  fi
}

# OS / arch for Go release asset names (match .goreleaser.yaml archives).
uname_s="$(uname -s 2>/dev/null || echo unknown)"
uname_m="$(uname -m 2>/dev/null || echo unknown)"
uname_s_lc="$(printf '%s' "$uname_s" | tr '[:upper:]' '[:lower:]')"

case "$uname_m" in
  x86_64|amd64) GOARCH="amd64" ;;
  aarch64|arm64) GOARCH="arm64" ;;
  *)
    echo "install.sh: unsupported architecture: $uname_m" >&2
    exit 1
    ;;
esac

case "$uname_s_lc" in
  linux*) GOOS="linux" ;;
  darwin*) GOOS="darwin" ;;
  mingw*|msys*|cygwin*) GOOS="windows" ;;
  *)
    echo "install.sh: unsupported OS: $uname_s" >&2
    exit 1
    ;;
esac

if [ "$GOOS" = "windows" ] && [ "$GOARCH" != "amd64" ]; then
  echo "install.sh: only windows/amd64 is published" >&2
  exit 1
fi

if [ -z "${VERSION:-}" ]; then
  need_cmd sed
  TAG="$(fetch "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n1)"
  if [ -z "$TAG" ]; then
    echo "install.sh: could not determine latest release tag" >&2
    exit 1
  fi
else
  TAG="$VERSION"
fi

case "$TAG" in
  v*) VERSION_NUM="${TAG#v}" ;;
  *)
    VERSION_NUM="$TAG"
    TAG="v${TAG}"
    ;;
esac

# Asset names from GoReleaser: nara_<semver>_linux_amd64.tar.gz (windows zip).
if [ "$GOOS" = "windows" ]; then
  ASSET="${BINARY}_${VERSION_NUM}_${GOOS}_${GOARCH}.zip"
else
  ASSET="${BINARY}_${VERSION_NUM}_${GOOS}_${GOARCH}.tar.gz"
fi

URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "install.sh: downloading ${URL}"
fetch "$URL" >"${TMP}/asset"

if [ "$GOOS" = "windows" ]; then
  need_cmd unzip
  unzip -q -o "${TMP}/asset" -d "$TMP"
  INSTALL_DIR="${USERPROFILE:-$HOME}/bin"
  mkdir -p "$INSTALL_DIR"
  cp -f "${TMP}/${BINARY}.exe" "${INSTALL_DIR}/${BINARY}.exe"
  echo "install.sh: installed ${INSTALL_DIR}/${BINARY}.exe"
  echo "install.sh: add ${INSTALL_DIR} to PATH if needed"
else
  need_cmd tar
  tar -xzf "${TMP}/asset" -C "$TMP"
  INSTALL_DIR="/usr/local/bin"
  if [ ! -w "$INSTALL_DIR" ] 2>/dev/null; then
    echo "install.sh: ${INSTALL_DIR} not writable; using sudo"
    need_cmd sudo
    sudo install -m 0755 "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  else
    install -m 0755 "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  fi
  echo "install.sh: installed ${INSTALL_DIR}/${BINARY}"
fi
