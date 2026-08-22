#!/usr/bin/env bash
# Shared setup and assertions for scripts/e2e-m*.sh. Source only — do not execute.
#
# Conventions:
# - Intermediate CLI success chatter is silenced (exit status is the signal).
# - Failures print context to stderr and exit non-zero.
# - Each script prints exactly one success line: "<name>: ok"

# e2e_setup NAME — create temp workdir, binary path, default socket/config, EXIT trap.
e2e_setup() {
	E2E_NAME="${1:?e2e_setup requires a name}"
	E2E_ROOT="$(cd "$(dirname "${BASH_SOURCE[1]}")/.." && pwd)"
	cd "$E2E_ROOT"

	WORKDIR="$(mktemp -d "${TMPDIR:-/tmp}/agentd-${E2E_NAME}.XXXXXX")"
	BIN="$WORKDIR/agentd"
	SOCK="$WORKDIR/agentd.sock"
	CFG="$WORKDIR/agentd.yaml"
	E2E_SOCKETS=("$SOCK")

	trap e2e_cleanup EXIT
}

# e2e_track_socket PATH — also stop this socket on cleanup (multi-daemon tests).
e2e_track_socket() {
	E2E_SOCKETS+=("${1:?}")
}

e2e_cleanup() {
	local s
	if [[ -n "${BIN:-}" && -x "${BIN}" ]]; then
		for s in "${E2E_SOCKETS[@]:-}"; do
			[[ -n "$s" ]] || continue
			"$BIN" daemon stop --socket "$s" --timeout 5s >/dev/null 2>&1 || true
		done
	fi
	if [[ -n "${E2E_DAEMON_PID:-}" ]]; then
		wait "$E2E_DAEMON_PID" 2>/dev/null || true
	fi
	if [[ -n "${WORKDIR:-}" && -d "${WORKDIR}" ]]; then
		rm -rf "$WORKDIR"
	fi
}

e2e_build() {
	go build -o "$BIN" .
}

# e2e_quiet CMD... — run command, keep stderr, discard stdout (for "ok"-style CLIs).
e2e_quiet() {
	"$@" >/dev/null
}

	e2e_daemon_start() {
		# Foreground avoids detach re-exec races; run in background with tracked PID.
		XDG_STATE_HOME="${XDG_STATE_HOME:-}" "$BIN" daemon start --foreground --socket "$SOCK" "$@" &
		E2E_DAEMON_PID=$!
		e2e_wait_ready "$SOCK"
	}

e2e_daemon_start_at() {
	local sock="$1"
	shift
	"$BIN" daemon start --socket "$sock" "$@"
	e2e_wait_ready "$sock"
}

e2e_wait_ready() {
	local sock="${1:-$SOCK}"
	local _
	for _ in $(seq 1 50); do
		if "$BIN" daemon status --socket "$sock" --json 2>/dev/null | grep -qE '"running"[[:space:]]*:[[:space:]]*true'; then
			return 0
		fi
		sleep 0.1
	done
	echo "${E2E_NAME}: daemon not ready on ${sock}" >&2
	"$BIN" daemon status --socket "$sock" --json >&2 || true
	exit 1
}

e2e_daemon_stop() {
	local sock="${1:-$SOCK}"
	"$BIN" daemon stop --socket "$sock" --timeout 5s
}

e2e_fail() {
	echo "${E2E_NAME}: $*" >&2
	exit 1
}

e2e_assert_eq() {
	local got="$1" want="$2" label="${3:-value}"
	if [[ "$got" != "$want" ]]; then
		e2e_fail "${label}: got $(printf %q "$got"), want $(printf %q "$want")"
	fi
}

e2e_assert_contains() {
	local haystack="$1" needle="$2" label="${3:-output}"
	if ! grep -Fq -- "$needle" <<<"$haystack"; then
		e2e_fail "${label}: expected to contain $(printf %q "$needle")"
	fi
}

e2e_assert_matches() {
	local haystack="$1" pattern="$2" label="${3:-output}"
	if ! grep -Eq -- "$pattern" <<<"$haystack"; then
		e2e_fail "${label}: expected to match /${pattern}/"
	fi
}

e2e_assert_not_contains() {
	local haystack="$1" needle="$2" label="${3:-output}"
	if grep -Fq -- "$needle" <<<"$haystack"; then
		e2e_fail "${label}: must not contain $(printf %q "$needle")"
	fi
}

# e2e_assert_trajectory_no_drops — trajectory queue overflow must stay zero.
e2e_assert_trajectory_no_drops() {
	local status
	status="$("$BIN" daemon status --socket "$SOCK" --json)"
	if grep -Eq '"trajectory_dropped_count"[[:space:]]*:[[:space:]]*[1-9][0-9]*' <<<"$status"; then
		e2e_fail "trajectory queue overflow: ${status}"
	fi
}

e2e_assert_file_contains() {
	local path="$1" needle="$2"
	if ! grep -Fq -- "$needle" "$path"; then
		e2e_fail "file ${path}: expected to contain $(printf %q "$needle")"
	fi
}

e2e_json_field() {
	# Extract first JSON string or number field value (best-effort for status blobs).
	local json="$1" field="$2"
	sed -n "s/.*\"${field}\"[[:space:]]*:[[:space:]]*\"\\([^\"]*\\)\".*/\\1/p;s/.*\"${field}\"[[:space:]]*:[[:space:]]*\\([0-9][0-9]*\\).*/\\1/p" <<<"$json" | head -n1
}

e2e_hook_run() {
	local provider="$1" payload="$2"
	printf '%s' "$payload" | "$BIN" hook run --socket "$SOCK" --provider="$provider"
}

# e2e_run_hook CMD... — run hook CLI; fail on non-zero exit.
e2e_run_hook() {
	if ! "$@" >/dev/null; then
		e2e_fail "hook failed: $*"
	fi
}

# e2e_trajectory_settle — brief pause for async trajectory queue + JSONL debounce.
e2e_trajectory_settle() {
	local secs="${E2E_TRAJECTORY_SETTLE_SECS:-0.2}"
	# bash sleep accepts fractions on modern macOS/Linux
	sleep "$secs"
}

# e2e_sessions_dir — ledger root (requires SESSIONS or XDG_STATE_HOME).
e2e_sessions_dir() {
	if [[ -n "${SESSIONS:-}" ]]; then
		printf '%s' "$SESSIONS"
		return 0
	fi
	if [[ -n "${XDG_STATE_HOME:-}" ]]; then
		printf '%s/agentd/sessions' "$XDG_STATE_HOME"
		return 0
	fi
	return 1
}

# e2e_poll_provider_sessions PROV [ROOT] — single check for *.jsonl under provider dir.
e2e_poll_provider_sessions() {
	local prov="$1" root="${2:-}"
	if [[ -z "$root" ]]; then
		root="$(e2e_sessions_dir)" || return 1
	fi
	find "$root/$prov" -name '*.jsonl' -print -quit 2>/dev/null | grep -q .
}

# _e2e_poll_wait_cycle PROV [ROOT] — one settle + poll window; returns 0 when ledger exists.
_e2e_poll_wait_cycle() {
	local prov="$1" root="${2:-}"
	if [[ -z "$root" ]]; then
		root="$(e2e_sessions_dir)" || return 1
	fi
	local max="${E2E_SESSION_WAIT_MAX:-100}"
	local poll="${E2E_SESSION_POLL_SECS:-0.1}"
	e2e_trajectory_settle
	local n=0
	while (( n < max )); do
		if e2e_poll_provider_sessions "$prov" "$root"; then
			return 0
		fi
		sleep "$poll"
		n=$((n + 1))
	done
	return 1
}

# e2e_wait_provider_sessions PROV [ROOT] — poll until ledger file exists or timeout.
# Tunables: E2E_SESSION_WAIT_MAX (iterations, default 100), E2E_SESSION_POLL_SECS (default 0.1),
# E2E_WAIT_RETRIES (full poll cycles, default 3), E2E_TRAJECTORY_SETTLE_SECS (default 0.2).
e2e_wait_provider_sessions() {
	local prov="$1" root="${2:-}"
	if [[ -z "$root" ]]; then
		root="$(e2e_sessions_dir)" || e2e_fail "sessions dir unavailable"
	fi
	local attempts="${E2E_WAIT_RETRIES:-3}"
	local cycle=1
	while (( cycle <= attempts )); do
		if _e2e_poll_wait_cycle "$prov" "$root"; then
			return 0
		fi
		if (( cycle < attempts )); then
			echo "${E2E_NAME}: timeout waiting for sessions under ${root}/${prov} (attempt ${cycle}/${attempts}), retrying..." >&2
		fi
		cycle=$((cycle + 1))
	done
	_e2e_sessions_timeout_diag "$prov" "$root"
	e2e_fail "timeout waiting for sessions under ${root}/${prov}"
}

# e2e_wait_provider_sessions_with_hook PROV -- HOOK_CMD...
# Re-runs HOOK_CMD before each retry cycle after the first (hook should have run once before calling).
e2e_wait_provider_sessions_with_hook() {
	local prov="$1"
	shift
	[[ "${1:-}" == "--" ]] || e2e_fail "e2e_wait_provider_sessions_with_hook: expected -- before hook command"
	shift
	local hook_cmd=("$@")
	local root
	root="$(e2e_sessions_dir)" || e2e_fail "sessions dir unavailable"
	local attempts="${E2E_WAIT_RETRIES:-3}"
	local cycle=1
	while (( cycle <= attempts )); do
		if (( cycle > 1 )); then
			e2e_run_hook "${hook_cmd[@]}"
		fi
		if _e2e_poll_wait_cycle "$prov" "$root"; then
			return 0
		fi
		if (( cycle < attempts )); then
			echo "${E2E_NAME}: timeout waiting for sessions under ${root}/${prov} (attempt ${cycle}/${attempts}), re-hook and retry..." >&2
		fi
		cycle=$((cycle + 1))
	done
	_e2e_sessions_timeout_diag "$prov" "$root"
	e2e_fail "timeout waiting for sessions under ${root}/${prov}"
}

_e2e_sessions_timeout_diag() {
	local prov="$1" root="$2"
	echo "${E2E_NAME}: timeout waiting for sessions under ${root}/${prov} after ${E2E_WAIT_RETRIES:-3} attempts" >&2
	if [[ -n "${BIN:-}" && -x "${BIN}" && -n "${SOCK:-}" ]]; then
		"$BIN" daemon status --socket "$SOCK" --json 2>/dev/null >&2 || true
	fi
	if [[ -d "$root" ]]; then
		find "$root" -maxdepth 2 -name '*.jsonl' 2>/dev/null >&2 || true
	fi
}

# e2e_retry ATTEMPTS LABEL CMD... — retry a command with delay between attempts.
e2e_retry() {
	local attempts="$1" label="$2"
	shift 2
	local n=1 delay="${E2E_RETRY_DELAY_SECS:-1}"
	while (( n <= attempts )); do
		if "$@"; then
			return 0
		fi
		if (( n < attempts )); then
			echo "${E2E_NAME}: ${label}: attempt ${n}/${attempts} failed, retrying in ${delay}s..." >&2
			sleep "$delay"
		fi
		n=$((n + 1))
	done
	e2e_fail "${label}: failed after ${attempts} attempts"
}

e2e_pass() {
	echo "${E2E_NAME}: ok"
}
