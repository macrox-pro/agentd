#!/usr/bin/env bash
# M6 acceptance: shell / mcp / paths guards + route subset
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m6

cat >"$CFG" <<'EOF'
version: 1
guards:
  secrets:
    enabled: false
  shell:
    enabled: true
    deny_patterns: ["rm -rf /"]
    ask_on: [curl]
  mcp:
    enabled: true
    deny_servers: ["untrusted-*"]
  paths:
    enabled: true
    deny_read: ["/etc/shadow"]
    deny_write: ["**/.env"]
dispatch:
  - name: shell-only
    match:
      kind: [tool.pre]
      tools: [Bash]
      provider: ["codex"]
    mode: sync_only
    sync:
      - target: builtin
        guards: [shell]
  - name: all-guards
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: builtin
        guards: [shell, mcp, paths]
EOF

e2e_build

e2e_quiet "$BIN" config validate --config "$CFG"

SHOW="$("$BIN" config show --config "$CFG" --merged)"
e2e_assert_contains "$SHOW" 'deny_patterns' show
e2e_assert_contains "$SHOW" 'deny_servers' show
e2e_assert_contains "$SHOW" 'deny_read' show

ROUTES_JSON="$("$BIN" dispatch routes --config "$CFG" --json)"
e2e_assert_contains "$ROUTES_JSON" 'all-guards' routes
e2e_assert_contains "$ROUTES_JSON" 'shell-only' routes
e2e_assert_contains "$ROUTES_JSON" '"shell"' routes
e2e_assert_contains "$ROUTES_JSON" '"mcp"' routes
e2e_assert_contains "$ROUTES_JSON" '"paths"' routes

e2e_daemon_start --config "$CFG"

DENY_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"sudo rm -rf /"}}'
DOUT="$(e2e_hook_run claude-code "$DENY_PAYLOAD")"
e2e_assert_contains "$DOUT" '"permissionDecision":"deny"' shell-deny

ASK_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t2","tool_input":{"command":"curl https://example.com"}}'
AOUT="$(e2e_hook_run claude-code "$ASK_PAYLOAD")"
e2e_assert_contains "$AOUT" '"permissionDecision":"ask"' shell-ask

CLEAN_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t3","tool_input":{"command":"go test ./..."}}'
COUT="$(e2e_hook_run claude-code "$CLEAN_PAYLOAD")"
e2e_assert_eq "$COUT" '{}' shell-clean

MCP_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"mcp__untrusted-foo__bar","tool_use_id":"t4","tool_input":{}}'
MOUT="$(e2e_hook_run claude-code "$MCP_PAYLOAD")"
e2e_assert_contains "$MOUT" '"permissionDecision":"deny"' mcp-deny

PATH_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Read","tool_use_id":"t5","tool_input":{"file_path":"/etc/shadow"}}'
POUT="$(e2e_hook_run claude-code "$PATH_PAYLOAD")"
e2e_assert_contains "$POUT" '"permissionDecision":"deny"' path-deny-read

ENV_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Write","tool_use_id":"t6","tool_input":{"file_path":"repo/.env"}}'
EOUT="$(e2e_hook_run claude-code "$ENV_PAYLOAD")"
e2e_assert_contains "$EOUT" '"permissionDecision":"deny"' path-deny-write

# subset route: shell-only for codex Bash — ask falls back to deny (no CapAsk)
SUB_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t7","tool_input":{"command":"curl https://example.com"}}'
SOUT="$(e2e_hook_run codex "$SUB_PAYLOAD")"
e2e_assert_matches "$SOUT" 'deny|blocked|ask_on' shell-subset

e2e_daemon_stop
e2e_pass
