#!/usr/bin/env bash
# M12 acceptance: live Subscribe, schema_version, multi-provider, no Invoke hang
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m12

export XDG_STATE_HOME="$WORKDIR/state"
mkdir -p "$XDG_STATE_HOME"
PROJ="$WORKDIR/repo"
mkdir -p "$PROJ"
SESSIONS="$XDG_STATE_HOME/agentd/sessions"
SUB_OUT="$WORKDIR/subscribe.jsonl"

cat >"$CFG" <<'EOF'
version: 1
trajectory:
  enabled: true
  include_raw: true
dispatch:
  - name: observe-all
    match:
      kind: ["*"]
      provider: ["*"]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: builtin
        observe: true
EOF

e2e_build
e2e_daemon_start --config "$CFG"

# Background subscribe before hooks (live firehose)
"$BIN" session subscribe --socket "$SOCK" --json >"$SUB_OUT" 2>/dev/null &
SUB_PID=$!
trap 'kill "$SUB_PID" 2>/dev/null || true' EXIT
sleep 0.3

CLAUDE='{"session_id":"m12-claude","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}'
printf '%s' "$CLAUDE" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=claude-code

CURSOR='{"conversation_id":"m12-cursor","cwd":"'"$PROJ"'","generation_id":"g1","hook_event_name":"preToolUse","workspace_roots":["'"$PROJ"'"],"tool_name":"Read","tool_input":{"path":"'"$PROJ"'/main.go"}}'
e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=cursor --argv-payload "$CURSOR"

e2e_trajectory_settle
e2e_wait_provider_sessions claude-code
e2e_wait_provider_sessions cursor

CAPTURE="$(cat "$SUB_OUT" || true)"
e2e_assert_contains "$CAPTURE" 'hook/invoked' subscribe-invoked
e2e_assert_contains "$CAPTURE" 'hook/decided' subscribe-decided
e2e_assert_contains "$CAPTURE" '"schema_version":1' subscribe-schema
e2e_assert_contains "$CAPTURE" 'm12-claude' subscribe-claude
e2e_assert_contains "$CAPTURE" 'm12-cursor' subscribe-cursor

# Filtered subscribe
SUB_FILT="$WORKDIR/subscribe-filter.jsonl"
"$BIN" session subscribe --socket "$SOCK" --provider claude-code --json >"$SUB_FILT" 2>/dev/null &
FILT_PID=$!
sleep 0.2
printf '%s' "$CLAUDE" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=claude-code
e2e_trajectory_settle
sleep 0.3
kill "$FILT_PID" 2>/dev/null || true
FILT_CAPTURE="$(cat "$SUB_FILT" || true)"
e2e_assert_contains "$FILT_CAPTURE" 'm12-claude' subscribe-filter-claude
e2e_assert_not_contains "$FILT_CAPTURE" 'm12-cursor' subscribe-filter-no-cursor

kill "$SUB_PID" 2>/dev/null || true
trap - EXIT

e2e_daemon_stop
test -f "$SESSIONS/claude-code/m12-claude.jsonl" || e2e_fail "claude ledger missing"
e2e_pass
