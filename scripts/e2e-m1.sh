#!/usr/bin/env bash
# M1 acceptance: detached start → status → hook run → reload → already-running → stop
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m1
e2e_build

e2e_daemon_start

STATUS="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_matches "$STATUS" '"running"[[:space:]]*:[[:space:]]*true' status

OUT="$(e2e_hook_run claude-code '{}')"
e2e_assert_eq "$OUT" '{}' hook-run

RELOAD="$("$BIN" daemon reload --socket "$SOCK")"
e2e_assert_contains "$RELOAD" 'generation=' reload

if "$BIN" daemon start --socket "$SOCK" 2>/dev/null; then
	e2e_fail 'expected already-running error'
fi

e2e_daemon_stop
e2e_pass
