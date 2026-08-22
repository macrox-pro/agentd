#!/usr/bin/env bash
# M10 acceptance: Claude import + search; L0 unchanged; importer status; daemon log file
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m10

export XDG_STATE_HOME="$WORKDIR/state"
mkdir -p "$XDG_STATE_HOME"
LOG_FILE="$XDG_STATE_HOME/agentd/agentd.log"
PROJ="$WORKDIR/repo"
mkdir -p "$PROJ"
SESSIONS="$XDG_STATE_HOME/agentd/sessions"
TRANSCRIPT_DIR="$WORKDIR/claude-projects"
mkdir -p "$TRANSCRIPT_DIR"

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

[[ -f "$LOG_FILE" ]] || e2e_fail "missing log file ${LOG_FILE}"
e2e_assert_file_contains "$LOG_FILE" "daemon ready" log-ready

CLAUDE='{"session_id":"m10-claude","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"toolu_01","tool_input":{"command":"echo ok"}}'
printf '%s' "$CLAUDE" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=claude-code
e2e_wait_provider_sessions claude-code

cat >"$TRANSCRIPT_DIR/m10-claude.jsonl" <<'JSONL'
{"type":"user","message":{"role":"user","content":"hello agentd"}}
{"type":"assistant","message":{"role":"assistant","content":[{"type":"thinking","thinking":"plan the response"},{"type":"text","text":"hi there"},{"type":"tool_use","id":"toolu_01","name":"Bash","input":{"command":"echo ok"}}]}}
JSONL

IMPORT_OUT="$("$BIN" session import --provider=claude-code --session=m10-claude --path="$TRANSCRIPT_DIR/m10-claude.jsonl" --json)"
e2e_assert_contains "$IMPORT_OUT" 'importer_status' import-status-field
e2e_assert_contains "$IMPORT_OUT" 'supported' import-status-value
e2e_assert_contains "$IMPORT_OUT" '"imported"' import-count

SHOW="$("$BIN" session show m10-claude --provider=claude-code --json)"
e2e_assert_contains "$SHOW" 'transcript/thinking' show-thinking
e2e_assert_contains "$SHOW" 'transcript' show-transcript-source
e2e_assert_contains "$SHOW" 'hook/invoked' show-hook

SEARCH="$("$BIN" session search --provider=claude-code --query=thinking --json)"
e2e_assert_contains "$SEARCH" 'transcript/thinking' search-hit

LIST_JSON="$("$BIN" session list --provider=claude-code --json)"
e2e_assert_contains "$LIST_JSON" 'importer_status' list-importer-claude
e2e_assert_contains "$LIST_JSON" 'supported' list-importer-claude-val

if "$BIN" session import --provider=cursor --session=x 2>/dev/null; then
	e2e_fail "cursor import should fail"
fi

# cursor L0 smoke
CURSOR='{"conversation_id":"m10-cursor","cwd":"'"$PROJ"'","generation_id":"g1","hook_event_name":"preToolUse","workspace_roots":["'"$PROJ"'"],"tool_name":"Read","tool_input":{"path":"'"$PROJ"'/main.go"}}'
e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=cursor --argv-payload "$CURSOR"
e2e_trajectory_settle
e2e_assert_trajectory_no_drops
e2e_wait_provider_sessions cursor
LIST_CURSOR="$("$BIN" session list --provider=cursor --json)"
e2e_assert_contains "$LIST_CURSOR" 'importer_status' list-importer-cursor
e2e_assert_contains "$LIST_CURSOR" 'none' list-importer-cursor-val

e2e_daemon_stop
e2e_assert_file_contains "$LOG_FILE" "daemon shutdown" log-shutdown
e2e_pass
