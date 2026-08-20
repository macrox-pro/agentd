#!/usr/bin/env bash
# M4 acceptance: grpc forward + install + OpenCode serve smoke
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m4

MAIN_DIR="$WORKDIR/main"
PEER_DIR="$WORKDIR/peer"
mkdir -p "$MAIN_DIR" "$PEER_DIR"
SOCK="$MAIN_DIR/agentd.sock"
PEER="$PEER_DIR/agentd.sock"
CFG="$MAIN_DIR/agentd.yaml"
PEER_CFG="$PEER_DIR/peer.yaml"
INST="$WORKDIR/install-root"
E2E_SOCKETS=("$SOCK" "$PEER")

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

e2e_build
e2e_daemon_start_at "$PEER" --config "$PEER_CFG"
e2e_daemon_start_at "$SOCK" --config "$CFG"

STATUS="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_matches "$STATUS" '"running"[[:space:]]*:[[:space:]]*true' status

ROUTES="$("$BIN" dispatch routes --config "$CFG")"
e2e_assert_contains "$ROUTES" 'grpc-async' routes
e2e_assert_contains "$ROUTES" 'grpc' routes

CLEAN_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"go test ./..."}}'
OUT="$(e2e_hook_run claude-code "$CLEAN_PAYLOAD")"
e2e_assert_eq "$OUT" '{}' clean-hook

# allow async grpc worker to finish
sleep 0.2

mkdir -p "$INST"
"$BIN" install --provider=claude-code --scope=project --dir "$INST"
test -f "$INST/.claude/settings.json"
e2e_assert_file_contains "$INST/.claude/settings.json" 'agenthooks'
e2e_assert_file_contains "$INST/.claude/settings.json" '--provider=claude-code'

"$BIN" install --provider=opencode --scope=project --dir "$INST"
test -f "$INST/.opencode/plugin/agenthooks.ts"
e2e_assert_file_contains "$INST/.opencode/plugin/agenthooks.ts" 'serve'

INIT_FRAME='{"seq":1,"hook":"initialize","input":{"serverUrl":"http://127.0.0.1:1","directory":"/work","worktree":""}}'
SERVE_OUT="$(printf '%s\n' "$INIT_FRAME" | "$BIN" hook serve --socket "$SOCK" --provider=opencode)"
e2e_assert_contains "$SERVE_OUT" '"seq":1' hook-serve

# sentinel path used by install-generated OpenCode shim
SERVE_OUT2="$(printf '%s\n' "$INIT_FRAME" | "$BIN" agenthooks serve --socket "$SOCK" --provider=opencode)"
e2e_assert_contains "$SERVE_OUT2" '"seq":1' agenthooks-serve

e2e_daemon_stop "$SOCK"
e2e_daemon_stop "$PEER"
e2e_pass
