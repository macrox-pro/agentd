#!/usr/bin/env bash
# M7 acceptance: approvals RecordDecision + runtime persist + temporary blocks
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m7

export XDG_STATE_HOME="$WORKDIR/state"
mkdir -p "$XDG_STATE_HOME"
PROJ="$WORKDIR/repo"
mkdir -p "$PROJ"

# Known fingerprint for shell ask_on hit "curl" on tool Bash (config.ApprovalFingerprint).
SHELL_FP='sha256:shell/cf9a28b78b54bc78c4242098c7860d460801c515dac2fc69e2b0f53f4546d513'

cat >"$CFG" <<'EOF'
version: 1
guards:
  secrets:
    enabled: false
  shell:
    enabled: true
    ask_on: [curl]
dispatch:
  - name: shell-gate
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: builtin
        guards: [shell]
EOF

cat >"$PROJ/.agentd.yaml" <<'EOF'
version: 1
EOF

e2e_build

e2e_daemon_start --config "$CFG"

ASK_PAYLOAD="{\"session_id\":\"s7\",\"cwd\":\"$PROJ\",\"hook_event_name\":\"PreToolUse\",\"tool_name\":\"Bash\",\"tool_use_id\":\"t1\",\"tool_input\":{\"command\":\"curl https://example.com\"}}"
AOUT="$(e2e_hook_run claude-code "$ASK_PAYLOAD")"
e2e_assert_contains "$AOUT" '"permissionDecision":"ask"' shell-ask
e2e_assert_contains "$AOUT" 'approval_fingerprint=' shell-ask-fp

# Prefer fingerprint from wire when present; fall back to known value.
FP="$(sed -n 's/.*approval_fingerprint=\(sha256:[^"\\ ]*\).*/\1/p' <<<"$AOUT" | head -n1)"
if [[ -z "$FP" ]]; then
	FP="$SHELL_FP"
fi
e2e_assert_eq "$FP" "$SHELL_FP" fingerprint-match

"$BIN" config record-decision --socket "$SOCK" \
	--fingerprint "$FP" --scope session --session-id s7 >/dev/null

ALLOW_OUT="$(e2e_hook_run claude-code "$ASK_PAYLOAD")"
e2e_assert_eq "$ALLOW_OUT" '{}' shell-approved

# Wrong session still asks.
WRONG_PAYLOAD="{\"session_id\":\"other\",\"cwd\":\"$PROJ\",\"hook_event_name\":\"PreToolUse\",\"tool_name\":\"Bash\",\"tool_use_id\":\"t2\",\"tool_input\":{\"command\":\"curl https://example.com\"}}"
WOUT="$(e2e_hook_run claude-code "$WRONG_PAYLOAD")"
e2e_assert_contains "$WOUT" '"permissionDecision":"ask"' wrong-session

# Persist across restart.
e2e_daemon_stop
RUNTIME_FILE="$XDG_STATE_HOME/agentd/runtime.yaml"
# Flush on stop; wait briefly if debounce race.
for _ in $(seq 1 20); do
	if [[ -f "$RUNTIME_FILE" ]]; then
		break
	fi
	sleep 0.1
done
e2e_assert_file_contains "$RUNTIME_FILE" "$FP"

e2e_daemon_start --config "$CFG"
ALLOW2="$(e2e_hook_run claude-code "$ASK_PAYLOAD")"
e2e_assert_eq "$ALLOW2" '{}' after-restart

# Temporary block denies.
UNTIL="$(date -u -v+1H +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -d '+1 hour' +%Y-%m-%dT%H:%M:%SZ)"
cat >"$WORKDIR/block.yaml" <<EOF
version: 1
blocks:
  temporary:
    - tool: Bash
      pattern: "blocked-cmd"
      reason: e2e-block
      until: "$UNTIL"
EOF
"$BIN" config patch --socket "$SOCK" --file "$WORKDIR/block.yaml" >/dev/null
BLOCK_PAYLOAD="{\"session_id\":\"s7\",\"cwd\":\"$PROJ\",\"hook_event_name\":\"PreToolUse\",\"tool_name\":\"Bash\",\"tool_use_id\":\"t3\",\"tool_input\":{\"command\":\"run blocked-cmd now\"}}"
BOUT="$(e2e_hook_run claude-code "$BLOCK_PAYLOAD")"
e2e_assert_contains "$BOUT" '"permissionDecision":"deny"' temp-block

e2e_daemon_stop
e2e_pass
