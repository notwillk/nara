#!/usr/bin/env bash
# Bump from latest GitHub release (semver); requires clean workspace (git status), main == origin/main, bump2version, jq, curl.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

PART="${1:-patch}"
case "$PART" in
    major|minor|patch) ;;
    *) echo "Usage: $0 [major|minor|patch] (default: patch)" >&2; exit 1 ;;
esac

if ! command -v bump2version >/dev/null 2>&1; then
    echo "release: bump2version not found (pip install bump2version)" >&2
    exit 1
fi
if ! command -v jq >/dev/null 2>&1; then
    echo "release: jq not found" >&2
    exit 1
fi

BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [ "$BRANCH" != "main" ]; then
    echo "release: must be on branch main (current: $BRANCH)" >&2
    exit 1
fi
if [ -n "$(git status --porcelain)" ]; then
    echo "release: workspace is not clean" >&2
    git status --short >&2
    exit 1
fi
git fetch origin main
UPSTREAM="$(git rev-parse origin/main)"
HEAD_SHA="$(git rev-parse HEAD)"
if [ "$HEAD_SHA" != "$UPSTREAM" ]; then
    echo "release: HEAD is not origin/main (local: $HEAD_SHA, origin/main: $UPSTREAM)" >&2
    exit 1
fi

REPO_API="https://api.github.com/repos/notwillk/nara/releases/latest"
RESP="$(curl -sSL -H "Accept: application/vnd.github+json" "$REPO_API")"
if echo "$RESP" | jq -e '.message == "Not Found"' >/dev/null 2>&1; then
    CURRENT_SEMVER="0.0.0"
else
    TAG="$(echo "$RESP" | jq -r '.tag_name // empty')"
    if [ -z "$TAG" ]; then
        echo "release: unexpected GitHub API response (no tag_name)" >&2
        echo "$RESP" | head -c 500 >&2
        exit 1
    fi
    TAG="${TAG#v}"
    if ! [[ "$TAG" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "release: latest release tag must be vMAJOR.MINOR.PATCH (got: $TAG)" >&2
        exit 1
    fi
    CURRENT_SEMVER="$TAG"
fi

set +e
LIST="$(bump2version --dry-run --list --no-configured-files --current-version "$CURRENT_SEMVER" "$PART" 2>&1)"
BV_EXIT=$?
set -e
if [ "$BV_EXIT" -ne 0 ]; then
    echo "release: bump2version failed" >&2
    echo "$LIST" >&2
    exit 1
fi
NEW_VERSION="$(echo "$LIST" | grep -E '^new_version=' | tail -1 | cut -d= -f2-)"
if [ -z "$NEW_VERSION" ]; then
    echo "release: bump2version did not produce new_version" >&2
    echo "$LIST" >&2
    exit 1
fi
if ! [[ "$NEW_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "release: unexpected new_version from bump2version: $NEW_VERSION" >&2
    exit 1
fi

TAG_NAME="v${NEW_VERSION}"
if git rev-parse "$TAG_NAME" >/dev/null 2>&1; then
    echo "release: tag $TAG_NAME already exists locally" >&2
    exit 1
fi
if [ -n "$(git ls-remote origin "refs/tags/$TAG_NAME" 2>/dev/null)" ]; then
    echo "release: tag $TAG_NAME already exists on origin" >&2
    exit 1
fi

echo "release: latest GitHub release baseline $CURRENT_SEMVER -> $TAG_NAME (bump $PART)"
git tag -a "$TAG_NAME" -m "Release $TAG_NAME"
git push origin "$TAG_NAME"
