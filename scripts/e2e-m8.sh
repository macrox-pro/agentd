#!/usr/bin/env bash
# M8 / v1 gate: Status drop counter + smoke across config/guards/approvals paths
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m8

export XDG_STATE_HOME="$WORKDIR/state"
mkdir -p "$XDG_STATE_HOME"
PROJ="$WORKDIR/repo"
mkdir -p "$PROJ"

cat >"$CFG" <<'EOF'
version: 1
async:
  queue_capacity: 8
  worker_limit: 1
  on_overflow: drop
guards:
  secrets:
    enabled: true
    action: ask
  shell:
    enabled: true
    ask_on: [curl]
dispatch:
  - name: gate
    match:
      kind: [tool.pre]
    mode: parallel
    sync_timeout: 20s
    sync:
      - target: builtin
        guards: [secrets, shell]
    async:
      - target: log
        level: info
EOF

cat >"$PROJ/.agentd.yaml" <<'EOF'
version: 1
EOF

e2e_build

e2e_daemon_start --config "$CFG"

# Status JSON includes async drop counter (zero before overflow).
STATUS="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_contains "$STATUS" '"running":true' status-running
e2e_assert_contains "$STATUS" 'async_dropped_count' status-dropped-field
e2e_assert_contains "$STATUS" 'async_queue_depth' status-depth-field

# Clean tool.pre → no-op.
CLEAN="{\"session_id\":\"s8\",\"cwd\":\"$PROJ\",\"hook_event_name\":\"PreToolUse\",\"tool_name\":\"Bash\",\"tool_use_id\":\"t1\",\"tool_input\":{\"command\":\"echo ok\"}}"
COUT="$(e2e_hook_run claude-code "$CLEAN")"
e2e_assert_eq "$COUT" '{}' clean-noop

# Shell ask_on still works (m5–m7 smoke).
ASK="{\"session_id\":\"s8\",\"cwd\":\"$PROJ\",\"hook_event_name\":\"PreToolUse\",\"tool_name\":\"Bash\",\"tool_use_id\":\"t2\",\"tool_input\":{\"command\":\"curl https://example.com\"}}"
AOUT="$(e2e_hook_run claude-code "$ASK")"
e2e_assert_contains "$AOUT" '"permissionDecision":"ask"' shell-ask
e2e_assert_contains "$AOUT" 'approval_fingerprint=' shell-ask-fp

# Config show / validate offline.
e2e_quiet "$BIN" config validate --config "$CFG" --cwd "$PROJ"
e2e_quiet "$BIN" config show --config "$CFG" --cwd "$PROJ" --merged

e2e_daemon_stop
e2e_pass
