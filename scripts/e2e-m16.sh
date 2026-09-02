#!/usr/bin/env bash
# M16 acceptance: Prometheus /metrics HTTP
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m16

PROJ="$WORKDIR/repo"
mkdir -p "$PROJ"
METRICS_PORT="$(e2e_free_tcp_port)"
METRICS_ADDR="127.0.0.1:${METRICS_PORT}"
FLAG_PORT="$(e2e_free_tcp_port)"
FLAG_ADDR="127.0.0.1:${FLAG_PORT}"

cat >"$CFG" <<'EOF'
version: 1
guards:
  secrets:
    enabled: false
dispatch:
  - name: observe
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: builtin
        guards: [secrets]
EOF

e2e_build

# metrics_off_by_default — no metrics listener
e2e_daemon_start --config "$CFG"
if curl -sf --max-time 1 "http://127.0.0.1:2112/metrics" >/dev/null 2>&1; then
	e2e_fail "metrics_off_by_default: unexpected listener on 127.0.0.1:2112"
fi
e2e_daemon_stop

# metrics_enabled_serves
cat >"$CFG" <<EOF
version: 1
metrics:
  enabled: true
  listen: ${METRICS_ADDR}
guards:
  secrets:
    enabled: false
dispatch:
  - name: observe
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: builtin
        guards: [secrets]
EOF

e2e_daemon_start --config "$CFG"
METRICS_BODY="$(curl -sf --max-time 2 "http://${METRICS_ADDR}/metrics")"
e2e_assert_contains "$METRICS_BODY" 'agentd_' metrics_enabled_serves
e2e_assert_contains "$METRICS_BODY" 'agentd_async_queue_depth' metrics_runtime_gauges
e2e_assert_contains "$METRICS_BODY" 'agentd_async_queue_dropped_total' metrics_runtime_gauges

# metrics_invoke_histogram
CLAUDE='{"session_id":"m16","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}'
printf '%s' "$CLAUDE" | e2e_run_hook "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=claude-code
sleep 0.2
METRICS_HOOK="$(curl -sf --max-time 2 "http://${METRICS_ADDR}/metrics")"
e2e_assert_contains "$METRICS_HOOK" 'agentd_invoke_duration_seconds' metrics_invoke_histogram
python3 -c '
import re, sys
body = sys.stdin.read()
if not re.search(r"^agentd_invoke_duration_seconds_count\{[^}]*\}\s+[1-9][0-9]*", body, re.M):
    raise SystemExit("metrics_invoke_histogram: invoke histogram count < 1")
' <<<"$METRICS_HOOK"

e2e_daemon_stop

# metrics_listener_released
if curl -sf --max-time 1 "http://${METRICS_ADDR}/metrics" >/dev/null 2>&1; then
	e2e_fail "metrics_listener_released: metrics still reachable after stop"
fi

# metrics_listen_flag_override
cat >"$CFG" <<'EOF'
version: 1
metrics:
  enabled: false
guards:
  secrets:
    enabled: false
dispatch:
  - name: observe
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: builtin
        guards: [secrets]
EOF

e2e_daemon_start --config "$CFG" --metrics-listen "$FLAG_ADDR"
FLAG_BODY="$(curl -sf --max-time 2 "http://${FLAG_ADDR}/metrics")"
e2e_assert_contains "$FLAG_BODY" 'agentd_' metrics_listen_flag_override

e2e_daemon_stop
if curl -sf --max-time 1 "http://${FLAG_ADDR}/metrics" >/dev/null 2>&1; then
	e2e_fail "metrics_listener_released: flag override listener still up after stop"
fi

e2e_pass
