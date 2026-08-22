#!/usr/bin/env bash
# M11 acceptance: multi-import (cursor partial), policy replay all six, session fork
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m11

export XDG_STATE_HOME="$WORKDIR/state"
mkdir -p "$XDG_STATE_HOME"
PROJ="$WORKDIR/repo"
mkdir -p "$PROJ"
SESSIONS="$XDG_STATE_HOME/agentd/sessions"
CURSOR_TX="$WORKDIR/cursor_transcript.jsonl"

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

# --- L0 with raw for all six providers ---
CLAUDE='{"session_id":"m11-claude","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}'
printf '%s' "$CLAUDE" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=claude-code

CURSOR='{"conversation_id":"m11-cursor","cwd":"'"$PROJ"'","generation_id":"g1","hook_event_name":"preToolUse","workspace_roots":["'"$PROJ"'"],"tool_name":"Read","tool_input":{"path":"'"$PROJ"'/main.go"}}'
e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=cursor --argv-payload "$CURSOR"

CODEX='{"session_id":"m11-codex","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"shell","tool_use_id":"c1","tool_input":{"command":["bash","-lc","echo ok"]}}'
printf '%s' "$CODEX" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=codex

GEMINI='{"session_id":"m11-gemini","cwd":"'"$PROJ"'","hook_event_name":"BeforeTool","tool_name":"run_shell_command","tool_input":{"command":"echo ok"},"tool_call_id":"g1"}'
printf '%s' "$GEMINI" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=gemini

KIMI='{"session_id":"m11-kimi","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"echo ok"},"tool_call_id":"k1"}'
printf '%s' "$KIMI" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --provider=kimi-code

OC_INIT='{"seq":1,"hook":"initialize","input":{"serverUrl":"http://127.0.0.1:1","directory":"'"$PROJ"'","worktree":""}}'
OC_TOOL='{"seq":2,"hook":"tool.execute.before","input":{"sessionID":"m11-opencode","callID":"oc1","tool":"bash"},"output":{"args":{"command":"echo ok","timeout":30}}}'
printf '%s\n%s\n' "$OC_INIT" "$OC_TOOL" | e2e_run_hook "$BIN" hook serve --socket "$SOCK" --provider=opencode

e2e_trajectory_settle
for prov in claude-code cursor codex gemini kimi-code opencode; do
	e2e_wait_provider_sessions "$prov"
done

# --- Policy replay for each provider ---
for prov_sid in \
	"claude-code:m11-claude" \
	"cursor:m11-cursor" \
	"codex:m11-codex" \
	"gemini:m11-gemini" \
	"kimi-code:m11-kimi" \
	"opencode:m11-opencode"; do
	prov="${prov_sid%%:*}"
	sid="${prov_sid##*:}"
	REPLAY="$("$BIN" session replay --policy --provider "$prov" --session "$sid" --config "$CFG" --json)"
	e2e_assert_contains "$REPLAY" '"hits"' "replay-$prov"
	e2e_assert_contains "$REPLAY" 'replay_decision' "replay-decision-$prov"
done

# --- Cursor partial import ---
cat >"$CURSOR_TX" <<'JSONL'
{"role":"user","message":{"content":[{"type":"text","text":"<user_query>\nhello m11\n</user_query>"}]}}
{"role":"assistant","message":{"content":[{"type":"text","text":"imported cursor text"}]}}
JSONL
IMPORT_OUT="$("$BIN" session import --provider cursor --session m11-cursor --path "$CURSOR_TX" --json)"
e2e_assert_contains "$IMPORT_OUT" '"importer_status": "partial"' import-cursor-status
e2e_assert_contains "$IMPORT_OUT" '"imported"' import-cursor-count

SEARCH="$("$BIN" session search --provider cursor --session m11-cursor --source transcript --query "imported cursor" --json)"
e2e_assert_contains "$SEARCH" 'transcript/message' search-cursor-import

LIST_JSON="$("$BIN" session list --provider cursor --json)"
e2e_assert_contains "$LIST_JSON" '"importer_status": "partial"' list-cursor-partial

# --- Codex supported import (temp sessions tree; never real ~/.codex) ---
CODEX_SESSIONS="$WORKDIR/codex/sessions"
CODEX_TX="$CODEX_SESSIONS/2026/07/26/rollout-2026-07-26T17-33-55-m11-codex.jsonl"
mkdir -p "$(dirname "$CODEX_TX")"
cat >"$CODEX_TX" <<'JSONL'
{"timestamp":"2026-07-26T14:33:55.115Z","type":"session_meta","payload":{"id":"m11-codex","session_id":"m11-codex","cwd":"/tmp"}}
{"timestamp":"2026-07-26T14:33:57.000Z","type":"event_msg","payload":{"type":"user_message","message":"hello codex m11"}}
{"timestamp":"2026-07-26T14:33:58.000Z","type":"event_msg","payload":{"type":"agent_reasoning","text":"**Planning m11**"}}
{"timestamp":"2026-07-26T14:33:59.000Z","type":"event_msg","payload":{"type":"agent_message","message":"imported codex text"}}
JSONL
# Point import.codex.path via a small overlay config for offline CLI
CODEX_CFG="$WORKDIR/codex-import.yaml"
cat >"$CODEX_CFG" <<EOF
version: 1
trajectory:
  enabled: true
  include_raw: true
  import:
    codex:
      enabled: false
      path: "$CODEX_SESSIONS"
EOF
CODEX_IMPORT="$("$BIN" session import --provider codex --session m11-codex --config "$CODEX_CFG" --json)"
e2e_assert_contains "$CODEX_IMPORT" '"importer_status": "supported"' import-codex-status
e2e_assert_contains "$CODEX_IMPORT" '"imported"' import-codex-count

CODEX_SEARCH="$("$BIN" session search --provider codex --session m11-codex --source transcript --query "imported codex" --json)"
e2e_assert_contains "$CODEX_SEARCH" 'transcript/message' search-codex-import
CODEX_THINK="$("$BIN" session search --provider codex --session m11-codex --source transcript --query "Planning m11" --json)"
e2e_assert_contains "$CODEX_THINK" 'transcript/thinking' search-codex-thinking

LIST_CODEX="$("$BIN" session list --provider codex --json)"
e2e_assert_contains "$LIST_CODEX" '"importer_status": "supported"' list-codex-supported

LIST_GEMINI="$("$BIN" session list --provider gemini --json)"
e2e_assert_contains "$LIST_GEMINI" '"importer_status": "none"' list-gemini-none

# --- Fork ---
FORK_OUT="$("$BIN" session fork --provider claude-code --session m11-claude --new-session m11-claude-fork --at-seq 2 --json)"
e2e_assert_contains "$FORK_OUT" '"new_session_id": "m11-claude-fork"' fork-new
e2e_assert_contains "$FORK_OUT" '"parent_session": "m11-claude"' fork-parent

SHOW_FORK="$("$BIN" session show m11-claude-fork --provider claude-code --json)"
e2e_assert_contains "$SHOW_FORK" 'session/fork' show-fork-event
e2e_assert_contains "$SHOW_FORK" 'parent_session' show-fork-meta

LIST_FORK="$("$BIN" session list --provider claude-code)"
e2e_assert_contains "$LIST_FORK" 'm11-claude-fork' list-fork

# Source still present
test -f "$SESSIONS/claude-code/m11-claude.jsonl" || e2e_fail "source ledger missing after fork"

e2e_daemon_stop
e2e_pass
