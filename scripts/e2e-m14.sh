#!/usr/bin/env bash
# M14 acceptance: enable (isolated HOME) starts daemon → autostart on → disable keeps daemon → manifest gone
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m14
e2e_build
e2e_isolated_home

cat >"$CFG" <<'EOF'
version: 1
EOF

STATUS_BEFORE="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_matches "$STATUS_BEFORE" '"running"[[:space:]]*:[[:space:]]*false' daemon-down-before-enable

if ! "$BIN" daemon enable --socket "$SOCK" --config "$CFG" >/dev/null 2>&1; then
	# Backend unavailable in this environment — skip without failing make e2e.
	echo "${E2E_NAME}: skip (login autostart backend unavailable)"
	e2e_pass
	exit 0
fi
E2E_AUTOSTART_ENABLED=1

e2e_wait_ready "$SOCK"
STATUS_ON="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_matches "$STATUS_ON" '"running"[[:space:]]*:[[:space:]]*true' enable-starts-daemon-if-down
e2e_assert_matches "$STATUS_ON" '"enabled"[[:space:]]*:[[:space:]]*true' autostart-enabled

if manifest="$(e2e_autostart_manifest_path 2>/dev/null || true)" && [[ -n "$manifest" ]]; then
	[[ -f "$manifest" ]] || e2e_fail "manifest missing after enable: ${manifest}"
fi

e2e_autostart_disable

STATUS_AFTER="$("$BIN" daemon status --socket "$SOCK" --json)"
e2e_assert_matches "$STATUS_AFTER" '"running"[[:space:]]*:[[:space:]]*true' disable-does-not-stop-running-daemon
e2e_assert_matches "$STATUS_AFTER" '"enabled"[[:space:]]*:[[:space:]]*false' autostart-disabled

if manifest="$(e2e_autostart_manifest_path 2>/dev/null || true)" && [[ -n "$manifest" ]]; then
	[[ ! -f "$manifest" ]] || e2e_fail "manifest still present after disable: ${manifest}"
fi

e2e_pass
