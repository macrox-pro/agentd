#!/usr/bin/env bash
# M2 acceptance: daemon + secrets guard + dispatch routes + status metrics
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

WORKDIR="$(mktemp -d)"
SOCK="$WORKDIR/agentd.sock"
BIN="$WORKDIR/agentd"
CFG="$WORKDIR/agentd.yaml"

cleanup() {
  "$BIN" daemon stop --socket "$SOCK" --timeout 5s >/dev/null 2>&1 || true
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

cat >"$CFG" <<'EOF'
version: 1
EOF

go build -o "$BIN" .

"$BIN" daemon start --socket "$SOCK" --config "$CFG"

for _ in $(seq 1 50); do
  if "$BIN" daemon status --socket "$SOCK" --json 2>/dev/null | grep -qE '"running"[[:space:]]*:[[:space:]]*true'; then
    break
  fi
  sleep 0.1
done

STATUS="$("$BIN" daemon status --socket "$SOCK" --json)"
echo "$STATUS" | grep -qE '"running"[[:space:]]*:[[:space:]]*true'
echo "$STATUS" | grep -qE '"compiled_route_count"[[:space:]]*:[[:space:]]*[1-9]'

CLEAN_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"go test ./..."}}'
OUT="$(echo "$CLEAN_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
test "$OUT" = '{}'

SECRET_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"}}'
SOUT="$(echo "$SECRET_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
echo "$SOUT" | grep -q '"permissionDecision":"ask"'
echo "$SOUT" | grep -vq 'AKIAIOSFODNN7EXAMPLE'

ROUTES="$("$BIN" dispatch routes --config "$CFG")"
echo "$ROUTES" | grep -q 'tool.pre'

"$BIN" daemon stop --socket "$SOCK" --timeout 5s

echo "e2e-m2: ok"
