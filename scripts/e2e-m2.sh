#!/usr/bin/env bash
# M2 acceptance: daemon + secrets guard + dispatch routes + status metrics
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m2

cat >"$CFG" <<'EOF'
version: 1
EOF

e2e_build
e2e_daemon_start --config "$CFG"

STATUS="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_matches "$STATUS" '"running"[[:space:]]*:[[:space:]]*true' status
e2e_assert_matches "$STATUS" '"compiled_route_count"[[:space:]]*:[[:space:]]*[1-9]' status

CLEAN_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"go test ./..."}}'
OUT="$(e2e_hook_run claude-code "$CLEAN_PAYLOAD")"
e2e_assert_eq "$OUT" '{}' clean-hook

SECRET_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE"}}'
SOUT="$(e2e_hook_run claude-code "$SECRET_PAYLOAD")"
e2e_assert_contains "$SOUT" '"permissionDecision":"ask"' secret-hook
e2e_assert_not_contains "$SOUT" 'AKIAIOSFODNN7EXAMPLE' secret-hook

ROUTES="$("$BIN" dispatch routes --config "$CFG")"
e2e_assert_contains "$ROUTES" 'tool.pre' routes

e2e_daemon_stop
e2e_pass
