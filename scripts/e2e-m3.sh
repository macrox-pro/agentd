#!/usr/bin/env bash
# M3 acceptance: dispatch YAML + file async + fsnotify reload
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

WORKDIR="$(mktemp -d)"
SOCK="$WORKDIR/agentd.sock"
BIN="$WORKDIR/agentd"
CFG="$WORKDIR/agentd.yaml"
AUDIT="$WORKDIR/audit.jsonl"

cleanup() {
  "$BIN" daemon stop --socket "$SOCK" --timeout 5s >/dev/null 2>&1 || true
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

cat >"$CFG" <<EOF
version: 1
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
      provider: ["*"]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: file
        path: $AUDIT
      - target: log
        level: info
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
GEN1="$(echo "$STATUS" | sed -n 's/.*"generation"[[:space:]]*:[[:space:]]*\([0-9][0-9]*\).*/\1/p')"

ROUTES="$("$BIN" dispatch routes --config "$CFG")"
echo "$ROUTES" | grep -q 'gate-and-audit'
echo "$ROUTES" | grep -q 'file'

CLEAN_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"go test ./..."}}'
OUT="$(echo "$CLEAN_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
test "$OUT" = '{}'

for _ in $(seq 1 50); do
  if [[ -s "$AUDIT" ]]; then
    break
  fi
  sleep 0.05
done
grep -q 'tool.pre' "$AUDIT"

cat >"$CFG" <<EOF
version: 1
policy:
  fail: fail_open
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: file
        path: $AUDIT
EOF

for _ in $(seq 1 80); do
  STATUS2="$("$BIN" daemon status --socket "$SOCK" --json)"
  GEN2="$(echo "$STATUS2" | sed -n 's/.*"generation"[[:space:]]*:[[:space:]]*\([0-9][0-9]*\).*/\1/p')"
  if [[ -n "$GEN2" && "$GEN2" -gt "$GEN1" ]]; then
    break
  fi
  sleep 0.05
done
STATUS2="$("$BIN" daemon status --socket "$SOCK" --json)"
GEN2="$(echo "$STATUS2" | sed -n 's/.*"generation"[[:space:]]*:[[:space:]]*\([0-9][0-9]*\).*/\1/p')"
test "$GEN2" -gt "$GEN1"

"$BIN" daemon stop --socket "$SOCK" --timeout 5s

echo "e2e-m3: ok"
