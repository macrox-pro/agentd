#!/usr/bin/env bash
# M13 acceptance: policy.offline on hook edge when daemon is down (fail_open / fail_closed, project merge)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m13
e2e_build

PROJ="$WORKDIR/proj"
mkdir -p "$PROJ"
STDERR="$WORKDIR/hook.stderr"

# Daemon intentionally not started — SOCK has no listener.

# Default fail_open (no user config): neutral wire, exit 0, stderr notice.
: >"$STDERR"
OUT="$(printf '{}' | "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=claude-code 2>"$STDERR")"
e2e_assert_eq "$OUT" '{}' offline-fail-open-out
e2e_assert_contains "$(cat "$STDERR")" 'daemon not running' offline-fail-open-stderr

# User policy.offline fail_closed → exit 1.
cat >"$CFG" <<'EOF'
version: 1
policy:
  offline: fail_closed
EOF
: >"$STDERR"
set +e
printf '{}' | "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=claude-code 2>"$STDERR"
CODE=$?
set -e
e2e_assert_eq "$CODE" 1 offline-fail-closed-exit
e2e_assert_contains "$(cat "$STDERR")" 'daemon not running' offline-fail-closed-stderr

# Project layer overrides user fail_open when payload cwd points at project config.
cat >"$CFG" <<'EOF'
version: 1
policy:
  offline: fail_open
EOF
cat >"$PROJ/.agentd.yaml" <<'EOF'
version: 1
policy:
  offline: fail_closed
EOF
PAYLOAD="{\"session_id\":\"s13\",\"cwd\":\"$PROJ\",\"hook_event_name\":\"PreToolUse\",\"tool_name\":\"Bash\",\"tool_use_id\":\"t1\",\"tool_input\":{\"command\":\"echo ok\"}}"
: >"$STDERR"
set +e
printf '%s' "$PAYLOAD" | "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=claude-code 2>"$STDERR"
PCODE=$?
set -e
e2e_assert_eq "$PCODE" 1 offline-project-fail-closed-exit
e2e_assert_contains "$(cat "$STDERR")" 'daemon not running' offline-project-fail-closed-stderr

# Codex notify offline: default fail_open → exit 0.
rm -f "$CFG"
: >"$STDERR"
set +e
"$BIN" hook notify --socket "$SOCK" --config "$CFG" --provider=codex '{"type":"agent-turn-complete","thread_id":"t1"}' 2>"$STDERR"
NCODE=$?
set -e
e2e_assert_eq "$NCODE" 0 offline-notify-fail-open
e2e_assert_contains "$(cat "$STDERR")" 'daemon not running' offline-notify-stderr

# Online baseline: daemon up → normal hook (not offline path).
cat >"$CFG" <<'EOF'
version: 1
guards:
  secrets:
    enabled: false
EOF
rm -f "$PROJ/.agentd.yaml"
e2e_daemon_start --config "$CFG"
ONLINE="$(e2e_hook_run claude-code "$PAYLOAD")"
e2e_assert_eq "$ONLINE" '{}' online-hook-out
e2e_daemon_stop

e2e_pass
