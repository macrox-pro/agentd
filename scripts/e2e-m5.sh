#!/usr/bin/env bash
# M5 acceptance: config layers + ConfigService + config CLI + merged fingerprint
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

WORKDIR="$(mktemp -d)"
SOCK="$WORKDIR/agentd.sock"
BIN="$WORKDIR/agentd"
CFG="$WORKDIR/user.yaml"
PROJ="$WORKDIR/repo"
PATCH="$WORKDIR/patch.yaml"
export XDG_STATE_HOME="$WORKDIR/state"

cleanup() {
  "$BIN" daemon stop --socket "$SOCK" --timeout 5s >/dev/null 2>&1 || true
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

mkdir -p "$PROJ" "$XDG_STATE_HOME"

cat >"$CFG" <<EOF
version: 1
policy:
  fail: fail_closed
dispatch:
  - name: user-route
    match:
      kind: [tool.pre]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: log
        level: info
EOF

cat >"$PROJ/.agentd.yaml" <<EOF
version: 1
policy:
  fail: fail_open
dispatch:
  - name: project-route
    match:
      kind: [tool.pre]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: log
        level: info
EOF

cat >"$PATCH" <<EOF
version: 1
async:
  queue_capacity: 42
EOF

go build -o "$BIN" .

"$BIN" config validate --config "$CFG"
"$BIN" config validate --config "$CFG" --cwd "$PROJ"

SHOW_USER="$("$BIN" config show --config "$CFG" --layer user)"
echo "$SHOW_USER" | grep -q 'fail_closed'

SHOW_MERGED="$("$BIN" config show --config "$CFG" --merged --cwd "$PROJ")"
echo "$SHOW_MERGED" | grep -q 'fail_open'

ROUTES_USER="$("$BIN" dispatch routes --config "$CFG")"
echo "$ROUTES_USER" | grep -q 'user-route'

ROUTES_PROJ="$("$BIN" dispatch routes --config "$CFG" --cwd "$PROJ")"
echo "$ROUTES_PROJ" | grep -q 'project-route'

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
FP1="$(echo "$STATUS" | sed -n 's/.*"fingerprint"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')"
test -n "$GEN1"
test -n "$FP1"

PATCH_OUT="$("$BIN" config patch --socket "$SOCK" --file "$PATCH")"
echo "$PATCH_OUT" | grep -q 'generation='
GEN2="$(echo "$PATCH_OUT" | sed -n 's/.*generation=\([0-9][0-9]*\).*/\1/p')"
test "$GEN2" -gt "$GEN1"

STATUS2="$("$BIN" daemon status --socket "$SOCK" --json)"
FP2="$(echo "$STATUS2" | sed -n 's/.*"fingerprint"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')"
test "$FP2" != "$FP1"

# shellcheck disable=SC2016
PAYLOAD=$(printf '%s' '{"session_id":"s","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo hi"}}')
OUT="$(echo "$PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
test "$OUT" = '{}'

echo "e2e-m5: ok"
