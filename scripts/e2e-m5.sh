#!/usr/bin/env bash
# M5 acceptance: config layers + ConfigService + config CLI + merged fingerprint
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m5
CFG="$WORKDIR/user.yaml"
PROJ="$WORKDIR/repo"
PATCH="$WORKDIR/patch.yaml"
export XDG_STATE_HOME="$WORKDIR/state"
mkdir -p "$PROJ" "$XDG_STATE_HOME"

cat >"$CFG" <<EOF
version: 1
policy:
  fail: fail_closed
dispatch:
  - name: user-route
    match:
      kind: [tool.pre]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: log
        level: info
EOF

cat >"$PROJ/.agentd.yaml" <<EOF
version: 1
policy:
  fail: fail_open
dispatch:
  - name: project-route
    match:
      kind: [tool.pre]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: log
        level: info
EOF

cat >"$PATCH" <<EOF
version: 1
async:
  queue_capacity: 42
EOF

e2e_build

# validate prints "ok" on success — assert via exit status only
e2e_quiet "$BIN" config validate --config "$CFG"
e2e_quiet "$BIN" config validate --config "$CFG" --cwd "$PROJ"

SHOW_USER="$("$BIN" config show --config "$CFG" --layer user)"
e2e_assert_contains "$SHOW_USER" 'fail_closed' show-user

SHOW_MERGED="$("$BIN" config show --config "$CFG" --merged --cwd "$PROJ")"
e2e_assert_contains "$SHOW_MERGED" 'fail_open' show-merged

ROUTES_USER="$("$BIN" dispatch routes --config "$CFG")"
e2e_assert_contains "$ROUTES_USER" 'user-route' routes-user

ROUTES_PROJ="$("$BIN" dispatch routes --config "$CFG" --cwd "$PROJ")"
e2e_assert_contains "$ROUTES_PROJ" 'project-route' routes-project

e2e_daemon_start --config "$CFG"

STATUS="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_matches "$STATUS" '"running"[[:space:]]*:[[:space:]]*true' status
GEN1="$(e2e_json_field "$STATUS" generation)"
FP1="$(e2e_json_field "$STATUS" fingerprint)"
test -n "$GEN1"
test -n "$FP1"

PATCH_OUT="$("$BIN" config patch --socket "$SOCK" --file "$PATCH")"
e2e_assert_contains "$PATCH_OUT" 'generation=' patch
GEN2="$(e2e_json_field "$PATCH_OUT" generation)"
# patch prints generation=N fingerprint=... (not JSON) — parse generation=
GEN2="${PATCH_OUT#*generation=}"
GEN2="${GEN2%% *}"
test "$GEN2" -gt "$GEN1"

STATUS2="$("$BIN" daemon status --socket "$SOCK" --json)"
FP2="$(e2e_json_field "$STATUS2" fingerprint)"
test "$FP2" != "$FP1"

PAYLOAD=$(printf '%s' '{"session_id":"s","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo hi"}}')
OUT="$(e2e_hook_run claude-code "$PAYLOAD")"
e2e_assert_eq "$OUT" '{}' project-hook

e2e_pass
