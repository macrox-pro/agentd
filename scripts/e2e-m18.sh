#!/usr/bin/env bash
# M18 acceptance: non-interactive TUI gate (setup / install)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m18
e2e_build

STDERR="$WORKDIR/hook.stderr"

# setup_without_tty
E2E_LAST_ERR_FILE="$STDERR"
e2e_expect_exit 1 -- env AGENTD_NO_TUI=1 "$BIN" setup
e2e_assert_contains "$(cat "$STDERR")" 'non-interactive environment' setup_without_tty

# setup_ci_env
: >"$STDERR"
E2E_LAST_ERR_FILE="$STDERR"
e2e_expect_exit 1 -- env CI=true "$BIN" setup
e2e_assert_contains "$(cat "$STDERR")" 'non-interactive environment' setup_ci_env

# install_bare_without_tty
: >"$STDERR"
E2E_LAST_ERR_FILE="$STDERR"
e2e_expect_exit 1 -- env AGENTD_NO_TUI=1 "$BIN" install
e2e_assert_contains "$(cat "$STDERR")" '--provider or --all-detected is required' install_bare_without_tty

# setup_yes_dryrun_conflict (install flags; setup requires TTY before flag validation)
: >"$STDERR"
E2E_LAST_ERR_FILE="$STDERR"
e2e_expect_exit 1 -- env AGENTD_NO_TUI=1 "$BIN" install --all-detected --yes --dry-run
e2e_assert_contains "$(cat "$STDERR")" 'mutually exclusive' setup_yes_dryrun_conflict

e2e_pass
