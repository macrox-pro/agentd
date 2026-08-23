#!/usr/bin/env bash
# Extract GitHub release notes for TAG from CHANGELOG.md (e.g. v0.0.2).
set -euo pipefail

tag="${1:?usage: release-notes.sh vX.Y.Z}"
ver="${tag#v}"

awk -v ver="$ver" '
  /^## \[/ {
    if (p && $0 !~ "\\[v" ver "\\]") exit
    if ($0 ~ "\\[v" ver "\\]") p = 1
  }
  p { print }
' CHANGELOG.md
