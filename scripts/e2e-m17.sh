#!/usr/bin/env bash
# M17 acceptance: doctor + install --all-detected (plan/apply) + trajectory default-on
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m17
e2e_build
e2e_isolated_home

cd "$WORKDIR"

mkdir -p "$HOME/.cursor" "$HOME/.claude"

DOCTOR_JSON="$("$BIN" doctor --json)"
e2e_assert_contains "$DOCTOR_JSON" '"cursor"' doctor-json-cursor
e2e_assert_contains "$DOCTOR_JSON" '"claude-code"' doctor-json-claude

PLAN_OUT="$("$BIN" install --all-detected)"
e2e_assert_contains "$PLAN_OUT" 'provider=cursor' all-detected-plan-cursor
e2e_assert_contains "$PLAN_OUT" 'hooks=missing' all-detected-plan-status
[[ ! -f "$HOME/.cursor/hooks.json" ]] || e2e_fail "hooks.json written before --yes"

APPLY_OUT="$("$BIN" install --all-detected --yes)"
e2e_assert_contains "$APPLY_OUT" 'hooks.json' all-detected-yes-hooks
e2e_assert_contains "$(cat "$HOME/.cursor/hooks.json")" 'agenthooks' hooks-json-agenthooks

TRAJ="$("$BIN" config get trajectory --config "$CFG")"
e2e_assert_contains "$TRAJ" 'on' trajectory-default-on
e2e_assert_contains "$TRAJ" '(default)' trajectory-default-source

e2e_daemon_start --config "$CFG"
DOCTOR_UP="$("$BIN" doctor --json --socket "$SOCK")"
e2e_assert_contains "$DOCTOR_UP" '"DaemonReachable": true' doctor-daemon-reachable
e2e_daemon_stop

e2e_pass
