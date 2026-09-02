#!/usr/bin/env bash
# M15 acceptance: trajectory stats + session stats; M19 cursor per-stop token sum
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m15

export XDG_STATE_HOME="$WORKDIR/state"
mkdir -p "$XDG_STATE_HOME"
PROJ="$WORKDIR/repo"
mkdir -p "$PROJ"
SESSIONS="$XDG_STATE_HOME/agentd/sessions"
STDERR="$WORKDIR/hook.stderr"

cat >"$CFG" <<'EOF'
version: 1
trajectory:
  enabled: true
  statistics: true
guards:
  secrets:
    enabled: false
dispatch:
  - name: observe-all
    match:
      kind: ["*"]
      provider: ["*"]
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: builtin
        observe: true
EOF

e2e_build

# stats_requires_daemon — no listener
E2E_LAST_ERR_FILE="$STDERR"
e2e_expect_exit 1 -- "$BIN" trajectory stats --config "$CFG" --socket "$WORKDIR/missing.sock"
e2e_assert_contains "$(cat "$STDERR")" 'daemon not running' stats_requires_daemon

e2e_daemon_start --config "$CFG"

CLAUDE_SESSION="m15-claude"
CLAUDE='{"session_id":"'"$CLAUDE_SESSION"'","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}'
printf '%s' "$CLAUDE" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=claude-code
e2e_trajectory_settle
e2e_wait_provider_sessions claude-code

STATS_JSON="$("$BIN" trajectory stats --config "$CFG" --socket "$SOCK" --json)"
python3 -c '
import json, sys
d = json.load(sys.stdin)
total = sum(d.get("rollup", {}).get("hooks_by_kind", {}).values())
if total < 1:
    raise SystemExit(f"stats_rollup_after_hooks: hooks total={total}")
' <<<"$STATS_JSON"

STATS_CLAUDE="$("$BIN" trajectory stats --config "$CFG" --socket "$SOCK" --provider claude-code --json)"
e2e_assert_contains "$STATS_CLAUDE" '"since"' stats_provider_filter

# stats_gate_statistics_off
e2e_quiet "$BIN" config disable trajectory-statistics --config "$CFG"
E2E_LAST_ERR_FILE="$STDERR"
e2e_expect_exit 1 -- "$BIN" trajectory stats --config "$CFG" --socket "$SOCK"
e2e_assert_contains "$(cat "$STDERR")" 'trajectory statistics is disabled' stats_gate_statistics_off
e2e_quiet "$BIN" config enable trajectory-statistics --config "$CFG"

# stats_gate_trajectory_off
e2e_quiet "$BIN" config disable trajectory --config "$CFG"
E2E_LAST_ERR_FILE="$STDERR"
e2e_expect_exit 1 -- "$BIN" trajectory stats --config "$CFG" --socket "$SOCK"
e2e_assert_contains "$(cat "$STDERR")" 'trajectory is disabled' stats_gate_trajectory_off
e2e_quiet "$BIN" config enable trajectory --config "$CFG"

e2e_daemon_stop

# session_stats_offline — daemon stopped
SESSION_STATS="$("$BIN" session stats "$CLAUDE_SESSION" --provider claude-code --config "$CFG" --json)"
e2e_assert_contains "$SESSION_STATS" '"session_id"' session_stats_offline
e2e_assert_contains "$SESSION_STATS" '"event_count"' session_stats_offline

# session_stats_not_found
E2E_LAST_ERR_FILE="$STDERR"
e2e_expect_exit 1 -- "$BIN" session stats missing-session --provider claude-code --config "$CFG"
e2e_assert_contains "$(cat "$STDERR")" 'not found' session_stats_not_found

# cursor_two_stops_sum_tokens (M19)
e2e_daemon_start --config "$CFG"
CURSOR1='{"conversation_id":"m15-cur","cwd":"'"$PROJ"'","generation_id":"g1","hook_event_name":"stop","input_tokens":100,"output_tokens":10,"cache_read_tokens":0}'
CURSOR2='{"conversation_id":"m15-cur","cwd":"'"$PROJ"'","generation_id":"g2","hook_event_name":"stop","input_tokens":200,"output_tokens":20,"cache_read_tokens":0}'
e2e_run_hook "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=cursor --argv-payload "$CURSOR1"
e2e_run_hook "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=cursor --argv-payload "$CURSOR2"
e2e_trajectory_settle

CURSOR_STATS="$("$BIN" trajectory stats --config "$CFG" --socket "$SOCK" --provider cursor --json)"
python3 -c '
import json, sys
d = json.load(sys.stdin)
r = d.get("rollup", {})
inp = r.get("input_tokens_total", 0)
out = r.get("output_tokens_total", 0)
if inp != 300 or out != 30:
    raise SystemExit(f"cursor_two_stops_sum_tokens: input={inp} output={out} want 300/30")
' <<<"$CURSOR_STATS"

e2e_daemon_stop
e2e_pass
