#!/usr/bin/env bash
# M1 acceptance: detached start → status → hook run → reload → already-running → stop
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

WORKDIR="$(mktemp -d)"
SOCK="$WORKDIR/agentd.sock"
BIN="$WORKDIR/agentd"

cleanup() {
  "$BIN" daemon stop --socket "$SOCK" --timeout 5s >/dev/null 2>&1 || true
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

go build -o "$BIN" .

"$BIN" daemon start --socket "$SOCK"

for _ in $(seq 1 50); do
  if "$BIN" daemon status --socket "$SOCK" --json 2>/dev/null | grep -qE '"running"[[:space:]]*:[[:space:]]*true'; then
    break
  fi
  sleep 0.1
done

"$BIN" daemon status --socket "$SOCK" --json | grep -qE '"running"[[:space:]]*:[[:space:]]*true'

OUT="$(echo '{}' | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
test "$OUT" = '{}'

"$BIN" daemon reload --socket "$SOCK" | grep -q generation=

if "$BIN" daemon start --socket "$SOCK" 2>/dev/null; then
  echo "expected already-running error" >&2
  exit 1
fi

"$BIN" daemon stop --socket "$SOCK" --timeout 5s

echo "e2e-m1: ok"
