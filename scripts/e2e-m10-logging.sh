#!/usr/bin/env bash
# Milestone: daemon operational logging — detached start writes agentd.log
set -euo pipefail
source "$(dirname "$0")/e2e-common.sh"

e2e_setup "e2e-m10-logging"
e2e_build

STATE="$WORKDIR/state"
export XDG_STATE_HOME="$STATE"
LOG_FILE="$STATE/agentd/agentd.log"

e2e_daemon_start_at "$SOCK"

if [[ ! -f "$LOG_FILE" ]]; then
	echo "e2e-m10-logging: missing log file $LOG_FILE" >&2
	exit 1
fi
if ! grep -q "daemon ready" "$LOG_FILE"; then
	echo "e2e-m10-logging: log missing daemon ready line" >&2
	cat "$LOG_FILE" >&2 || true
	exit 1
fi

e2e_quiet "$BIN" daemon stop --socket "$SOCK" --timeout 5s
if ! grep -q "daemon shutdown" "$LOG_FILE"; then
	echo "e2e-m10-logging: log missing shutdown line" >&2
	cat "$LOG_FILE" >&2 || true
	exit 1
fi

echo "e2e-m10-logging: ok"
