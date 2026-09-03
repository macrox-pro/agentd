#!/usr/bin/env bash
# M21 acceptance: full hook kind coverage (Tier A install + default routes)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m21
cat >"$CFG" <<'EOF'
version: 1
EOF

e2e_build
e2e_isolated_home

cd "$WORKDIR"

mkdir -p "$HOME/.cursor" "$HOME/.gemini"

"$BIN" install --provider=cursor --scope=user >/dev/null
CURSOR_HOOKS="$(cat "$HOME/.cursor/hooks.json")"
e2e_assert_contains "$CURSOR_HOOKS" 'subagentStart' cursor-tier-a-subagent
e2e_assert_contains "$CURSOR_HOOKS" 'preCompact' cursor-tier-a-compact
e2e_assert_not_contains "$CURSOR_HOOKS" 'afterAgentThought' cursor-no-tier-b-thought
e2e_assert_not_contains "$CURSOR_HOOKS" 'afterFileEdit' cursor-no-tier-b-file

"$BIN" install --provider=gemini --scope=user >/dev/null
GEMINI_HOOKS="$(cat "$HOME/.gemini/settings.json")"
e2e_assert_not_contains "$GEMINI_HOOKS" 'subagentStart' gemini-skip-subagent

ROUTES="$("$BIN" dispatch routes --config "$CFG")"
e2e_assert_contains "$ROUTES" 'default-subagent.start' routes-subagent
e2e_assert_contains "$ROUTES" 'default-compact.pre' routes-compact

e2e_pass
