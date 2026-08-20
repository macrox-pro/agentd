#!/usr/bin/env bash
# M3 acceptance: dispatch YAML + file async + fsnotify reload
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m3
AUDIT="$WORKDIR/audit.jsonl"

cat >"$CFG" <<EOF
version: 1
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
      provider: ["*"]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: file
        path: $AUDIT
      - target: log
        level: info
EOF

e2e_build
e2e_daemon_start --config "$CFG"

STATUS="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_matches "$STATUS" '"running"[[:space:]]*:[[:space:]]*true' status
GEN1="$(e2e_json_field "$STATUS" generation)"
test -n "$GEN1"

ROUTES="$("$BIN" dispatch routes --config "$CFG")"
e2e_assert_contains "$ROUTES" 'gate-and-audit' routes
e2e_assert_contains "$ROUTES" 'file' routes

CLEAN_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"go test ./..."}}'
OUT="$(e2e_hook_run claude-code "$CLEAN_PAYLOAD")"
e2e_assert_eq "$OUT" '{}' clean-hook

for _ in $(seq 1 50); do
	if [[ -s "$AUDIT" ]]; then
		break
	fi
	sleep 0.05
done
e2e_assert_file_contains "$AUDIT" 'tool.pre'

cat >"$CFG" <<EOF
version: 1
policy:
  fail: fail_open
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: file
        path: $AUDIT
EOF

GEN2=""
for _ in $(seq 1 80); do
	STATUS2="$("$BIN" daemon status --socket "$SOCK" --json)"
	GEN2="$(e2e_json_field "$STATUS2" generation)"
	if [[ -n "$GEN2" && "$GEN2" -gt "$GEN1" ]]; then
		break
	fi
	sleep 0.05
done
test "$GEN2" -gt "$GEN1"

e2e_daemon_stop
e2e_pass
