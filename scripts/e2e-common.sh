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
	"$BIN" daemon start --socket "$SOCK" "$@"
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

e2e_pass() {
	echo "${E2E_NAME}: ok"
}
