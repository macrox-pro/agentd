#!/usr/bin/env bash
# M4 acceptance: grpc forward + install + OpenCode serve smoke
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

WORKDIR="$(mktemp -d)"
MAIN_DIR="$WORKDIR/main"
PEER_DIR="$WORKDIR/peer"
mkdir -p "$MAIN_DIR" "$PEER_DIR"
SOCK="$MAIN_DIR/agentd.sock"
PEER="$PEER_DIR/agentd.sock"
BIN="$WORKDIR/agentd"
CFG="$MAIN_DIR/agentd.yaml"
PEER_CFG="$PEER_DIR/peer.yaml"
INST="$WORKDIR/install-root"

cleanup() {
  "$BIN" daemon stop --socket "$SOCK" --timeout 5s >/dev/null 2>&1 || true
  "$BIN" daemon stop --socket "$PEER" --timeout 5s >/dev/null 2>&1 || true
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

cat >"$PEER_CFG" <<EOF
version: 1
EOF

cat >"$CFG" <<EOF
version: 1
dispatch:
  - name: grpc-async
    match:
      kind: [tool.pre]
      provider: ["*"]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: grpc
        endpoint: unix://$PEER
        timeout: 2s
  - name: grpc-sync
    match:
      kind: [prompt.submitted]
    mode: sync_only
    sync:
      - target: grpc
        endpoint: unix://$PEER
        timeout: 2s
        on_error: fail_open
EOF

go build -o "$BIN" .

"$BIN" daemon start --socket "$PEER" --config "$PEER_CFG"
"$BIN" daemon start --socket "$SOCK" --config "$CFG"

for _ in $(seq 1 50); do
  if "$BIN" daemon status --socket "$SOCK" --json 2>/dev/null | grep -qE '"running"[[:space:]]*:[[:space:]]*true'; then
    break
  fi
  sleep 0.1
done

STATUS="$("$BIN" daemon status --socket "$SOCK" --json)"
echo "$STATUS" | grep -qE '"running"[[:space:]]*:[[:space:]]*true'

ROUTES="$("$BIN" dispatch routes --config "$CFG")"
echo "$ROUTES" | grep -q 'grpc-async'
echo "$ROUTES" | grep -q 'grpc'

CLEAN_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"go test ./..."}}'
OUT="$(echo "$CLEAN_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
test "$OUT" = '{}'

# allow async grpc worker to finish
sleep 0.2

mkdir -p "$INST"
"$BIN" install --provider=claude-code --scope=project --dir "$INST"
test -f "$INST/.claude/settings.json"
grep -q 'agenthooks' "$INST/.claude/settings.json"
grep -q -- '--provider=claude-code' "$INST/.claude/settings.json"

"$BIN" install --provider=opencode --scope=project --dir "$INST"
test -f "$INST/.opencode/plugin/agenthooks.ts"
grep -q 'serve' "$INST/.opencode/plugin/agenthooks.ts"

INIT_FRAME='{"seq":1,"hook":"initialize","input":{"serverUrl":"http://127.0.0.1:1","directory":"/work","worktree":""}}'
SERVE_OUT="$(printf '%s\n' "$INIT_FRAME" | "$BIN" hook serve --socket "$SOCK" --provider=opencode)"
echo "$SERVE_OUT" | grep -q '"seq":1'

# sentinel path used by install-generated OpenCode shim
SERVE_OUT2="$(printf '%s\n' "$INIT_FRAME" | "$BIN" agenthooks serve --socket "$SOCK" --provider=opencode)"
echo "$SERVE_OUT2" | grep -q '"seq":1'

"$BIN" daemon stop --socket "$SOCK" --timeout 5s
"$BIN" daemon stop --socket "$PEER" --timeout 5s

echo "e2e-m4: ok"
