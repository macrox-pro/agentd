#!/usr/bin/env bash
# M9 acceptance: trajectory L0 ledger for all six providers + session export
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m9

export XDG_STATE_HOME="$WORKDIR/state"
mkdir -p "$XDG_STATE_HOME"
PROJ="$WORKDIR/repo"
mkdir -p "$PROJ"
SESSIONS="$XDG_STATE_HOME/agentd/sessions"

cat >"$CFG" <<'EOF'
version: 1
trajectory:
  enabled: true
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

e2e_run_hook() {
	if ! "$@" >/dev/null; then
		e2e_fail "hook failed: $*"
	fi
}

e2e_wait_provider_sessions() {
	local prov="$1" n=0
	while (( n < 50 )); do
		if find "$SESSIONS/$prov" -name '*.jsonl' -print -quit 2>/dev/null | grep -q .; then
			return 0
		fi
		sleep 0.1
		n=$((n + 1))
	done
	e2e_fail "timeout waiting for sessions under $SESSIONS/$prov"
}

STATUS="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_contains "$STATUS" '"running":true' status-running
e2e_assert_contains "$STATUS" 'trajectory_dropped_count' status-trajectory-field

# claude-code — hook run (stdin)
CLAUDE='{"session_id":"m9-claude","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}'
printf '%s' "$CLAUDE" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=claude-code

# cursor — argv-payload
CURSOR='{"conversation_id":"m9-cursor","generation_id":"g1","hook_event_name":"preToolUse","workspace_roots":["'"$PROJ"'"],"tool_name":"Read","tool_input":{"path":"'"$PROJ"'/main.go"}}'
e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=cursor --argv-payload "$CURSOR"

# codex — run + notify
CODEX='{"session_id":"m9-codex","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"shell","tool_use_id":"c1","tool_input":{"command":["bash","-lc","echo ok"]}}'
printf '%s' "$CODEX" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=codex
CODEX_NOTIFY='{"session_id":"m9-codex-notify","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"shell","tool_use_id":"n1","tool_input":{"command":["bash","-lc","echo notify"]}}'
e2e_run_hook "$BIN" hook notify --socket "$SOCK" --provider=codex "$CODEX_NOTIFY"

# gemini — stdin
GEMINI='{"session_id":"m9-gemini","cwd":"'"$PROJ"'","hook_event_name":"BeforeTool","tool_name":"run_shell_command","tool_input":{"command":"echo ok"},"tool_call_id":"g1"}'
printf '%s' "$GEMINI" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=gemini

# kimi-code — stdin
KIMI='{"session_id":"m9-kimi","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"echo ok"},"tool_call_id":"k1"}'
printf '%s' "$KIMI" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=kimi-code

# opencode — serve NDJSON frame
OC_INIT='{"seq":1,"hook":"initialize","input":{"serverUrl":"http://127.0.0.1:1","directory":"'"$PROJ"'","worktree":""}}'
OC_TOOL='{"seq":2,"hook":"tool.execute.before","input":{"sessionID":"m9-opencode","callID":"oc1","tool":"bash"},"output":{"args":{"command":"echo ok","timeout":30}}}'
printf '%s\n%s\n' "$OC_INIT" "$OC_TOOL" | e2e_run_hook "$BIN" hook serve --socket "$SOCK" --provider=opencode

for prov in claude-code cursor codex gemini kimi-code opencode; do
	e2e_wait_provider_sessions "$prov"
done

# Contiguous seq in claude ledger
CLAUDE_LOG="$(find "$SESSIONS/claude-code" -name '*.jsonl' | head -n1)"
python3 - "$CLAUDE_LOG" <<'PY'
import json, sys
seqs = []
with open(sys.argv[1]) as f:
    for line in f:
        line = line.strip()
        if not line:
            continue
        seqs.append(json.loads(line)["seq"])
want = list(range(1, len(seqs) + 1))
if seqs != want:
    raise SystemExit(f"non-contiguous seq: {seqs}")
PY

EXPORT="$WORKDIR/export.jsonl"
e2e_quiet "$BIN" session export --provider=claude-code --out "$EXPORT"
test -s "$EXPORT"
e2e_assert_contains "$(cat "$EXPORT")" 'hook/invoked' export-invoked

LIST="$("$BIN" session list --provider=claude-code)"
e2e_assert_contains "$LIST" 'claude-code' session-list

e2e_daemon_stop
e2e_pass
