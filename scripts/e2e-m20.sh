#!/usr/bin/env bash
# M20 acceptance: policy.fail on daemon path, ask_fallback, cwd on notify/serve
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=e2e-common.sh
source "${SCRIPT_DIR}/e2e-common.sh"

e2e_setup e2e-m20

export XDG_STATE_HOME="$WORKDIR/state"
mkdir -p "$XDG_STATE_HOME"
PROJ="$WORKDIR/proj"
OTHER="$WORKDIR/other"
mkdir -p "$PROJ" "$OTHER"
DEAD_SOCK="$WORKDIR/dead.sock"
SESSIONS="$XDG_STATE_HOME/agentd/sessions"

write_sync_fail_config() {
	local fail_mode="$1"
	cat >"$CFG" <<EOF
version: 1
trajectory:
  enabled: true
guards:
  secrets:
    enabled: false
policy:
  fail: ${fail_mode}
dispatch:
  - name: grpc-sync-fail-tool
    match:
      kind: [tool.pre]
      provider: [claude-code]
    mode: sync_only
    sync_timeout: 1ns
    sync:
      - target: grpc
        endpoint: unix://${DEAD_SOCK}
        timeout: 2s
  - name: grpc-sync-fail-prompt
    match:
      kind: [prompt.submitted]
    mode: sync_only
    sync_timeout: 1ns
    sync:
      - target: grpc
        endpoint: unix://${DEAD_SOCK}
        timeout: 2s
EOF
}

write_ask_fallback_config() {
	local fb="${1:-}"
	cat >"$CFG" <<EOF
version: 1
trajectory:
  enabled: true
guards:
  secrets:
    enabled: false
  shell:
    enabled: true
    ask_on: [curl]
policy:
  fail: fail_open
$(if [[ -n "$fb" ]]; then echo "  ask_fallback: ${fb}"; fi)
dispatch:
  - name: shell-only
    match:
      kind: [tool.pre]
      provider: [codex]
    mode: parallel
    sync:
      - target: builtin
        guards: [shell]
    async:
      - target: builtin
        observe: true
EOF
}

write_cwd_config() {
	cat >"$CFG" <<EOF
version: 1
trajectory:
  enabled: true
guards:
  secrets:
    enabled: false
policy:
  fail: fail_open
dispatch:
  - name: grpc-codex-notify-fail
    match:
      kind: [notification]
      provider: [codex]
    mode: sync_only
    sync_timeout: 1ns
    sync:
      - target: grpc
        endpoint: unix://${DEAD_SOCK}
        timeout: 2s
  - name: grpc-opencode-fail
    match:
      kind: [tool.pre]
      provider: [opencode]
    mode: sync_only
    sync_timeout: 1ns
    sync:
      - target: grpc
        endpoint: unix://${DEAD_SOCK}
        timeout: 2s
EOF
	cat >"$PROJ/.agentd.yaml" <<'EOF'
version: 1
policy:
  fail: fail_closed
EOF
}

e2e_build

CLAUDE_TOOL='{"session_id":"m20","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}'
CODEX_CURL='{"session_id":"m20-codex","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"shell","tool_use_id":"c1","tool_input":{"command":"curl https://example.com"}}'
CODEX_CURL_NOD='{"session_id":"m20-codex-nod","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"shell","tool_use_id":"c2","tool_input":{"command":"curl https://example.com"}}'
CODEX_CURL_DENY='{"session_id":"m20-codex-deny","cwd":"'"$PROJ"'","hook_event_name":"PreToolUse","tool_name":"shell","tool_use_id":"c3","tool_input":{"command":"curl https://example.com"}}'

# sync_failure_fail_closed_denies
write_sync_fail_config fail_closed
e2e_daemon_start --config "$CFG"
OUT="$(e2e_hook_run claude-code "$CLAUDE_TOOL")"
e2e_assert_contains "$OUT" '"permissionDecision":"deny"' sync_failure_fail_closed_denies
e2e_daemon_stop

# sync_failure_fail_open_neutral
write_sync_fail_config fail_open
e2e_daemon_start --config "$CFG"
OUT="$(e2e_hook_run claude-code "$CLAUDE_TOOL")"
e2e_assert_eq "$OUT" '{}' sync_failure_fail_open_neutral
e2e_daemon_stop

# sync_failure_prompt_blocked
write_sync_fail_config fail_closed
e2e_daemon_start --config "$CFG"
PROMPT='{"session_id":"m20-prompt","cwd":"'"$PROJ"'","hook_event_name":"UserPromptSubmit","prompt":"hello"}'
OUT="$(printf '%s' "$PROMPT" | "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=claude-code)"
e2e_assert_not_eq "$OUT" '{}' sync_failure_prompt_blocked
e2e_assert_matches "$OUT" 'sync pipeline failed|block|deny' sync_failure_prompt_blocked
e2e_daemon_stop

# ask_fallback_no_decision
rm -rf "$SESSIONS/codex" 2>/dev/null || true
write_ask_fallback_config no_decision
e2e_daemon_start --config "$CFG"
e2e_run_hook "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=codex <<<"$CODEX_CURL_NOD"
e2e_trajectory_settle
e2e_wait_provider_sessions codex
CODEX_LOG="$SESSIONS/codex/m20-codex-nod.jsonl"
python3 - "$CODEX_LOG" "m20-codex-nod" "DECISION_KIND_NO_DECISION" <<'PY'
import json, sys
path, session, want = sys.argv[1], sys.argv[2], sys.argv[3]
found = False
with open(path) as f:
    for line in f:
        ev = json.loads(line)
        if ev.get("type") != "hook/decided" or ev.get("session_id") != session:
            continue
        data = ev.get("data", {})
        if isinstance(data, str):
            data = json.loads(data)
        if data.get("decision") == want:
            found = True
            break
if not found:
    raise SystemExit(f"ask_fallback_no_decision: expected {want} for session {session}")
PY
e2e_daemon_stop

# ask_fallback_deny_default
write_ask_fallback_config ""
e2e_daemon_start --config "$CFG"
e2e_run_hook "$BIN" hook run --socket "$SOCK" --config "$CFG" --provider=codex <<<"$CODEX_CURL_DENY"
e2e_trajectory_settle
e2e_wait_provider_sessions codex
CODEX_LOG="$SESSIONS/codex/m20-codex-deny.jsonl"
python3 - "$CODEX_LOG" "m20-codex-deny" "DECISION_KIND_DENY" <<'PY'
import json, sys
path, session, want = sys.argv[1], sys.argv[2], sys.argv[3]
found = False
with open(path) as f:
    for line in f:
        ev = json.loads(line)
        if ev.get("type") != "hook/decided" or ev.get("session_id") != session:
            continue
        data = ev.get("data", {})
        if isinstance(data, str):
            data = json.loads(data)
        if data.get("decision") == want:
            found = True
            break
if not found:
    raise SystemExit(f"ask_fallback_deny_default: expected {want} for session {session}")
PY
e2e_daemon_stop

# notify_uses_payload_cwd + serve_resolves_cwd_per_frame
write_cwd_config
rm -rf "$SESSIONS/codex" 2>/dev/null || true
e2e_daemon_start --config "$CFG"

NOTIFY_OTHER='{"type":"agent-turn-complete","thread_id":"m20-notify-other","cwd":"'"$OTHER"'"}'
NOTIFY_PROJ='{"type":"agent-turn-complete","thread_id":"m20-notify-proj","cwd":"'"$PROJ"'"}'
e2e_run_hook "$BIN" hook notify --socket "$SOCK" --config "$CFG" --provider=codex "$NOTIFY_OTHER"
e2e_trajectory_settle
e2e_wait_provider_sessions codex
FP_OTHER="$(python3 - "$SESSIONS/codex" "$OTHER" <<'PY'
import json, os, sys
root, cwd = sys.argv[1], sys.argv[2]
fp = ""
for dirpath, _, names in os.walk(root):
    for name in names:
        if not name.endswith(".jsonl"):
            continue
        with open(os.path.join(dirpath, name)) as f:
            for line in f:
                ev = json.loads(line)
                if ev.get("type") != "hook/decided" or ev.get("cwd") != cwd:
                    continue
                data = ev.get("data", {})
                if isinstance(data, str):
                    data = json.loads(data)
                if data.get("decision") != "DECISION_KIND_NO_DECISION":
                    continue
                fp = data.get("config_fingerprint", "")
                break
        if fp:
            break
    if fp:
        break
if not fp:
    raise SystemExit(f"notify_uses_payload_cwd: no hook/decided for cwd {cwd}")
print(fp)
PY
)"
e2e_run_hook "$BIN" hook notify --socket "$SOCK" --config "$CFG" --provider=codex "$NOTIFY_PROJ"
e2e_trajectory_settle
python3 - "$SESSIONS/codex" "$PROJ" "$FP_OTHER" <<'PY'
import json, os, sys
root, cwd, fp_other = sys.argv[1], sys.argv[2], sys.argv[3]
fp = ""
for dirpath, _, names in os.walk(root):
    for name in names:
        if not name.endswith(".jsonl"):
            continue
        with open(os.path.join(dirpath, name)) as f:
            for line in f:
                ev = json.loads(line)
                if ev.get("type") != "hook/decided" or ev.get("cwd") != cwd:
                    continue
                data = ev.get("data", {})
                if isinstance(data, str):
                    data = json.loads(data)
                fp = data.get("config_fingerprint", "")
                break
        if fp:
            break
    if fp:
        break
if not fp:
    raise SystemExit(f"notify_uses_payload_cwd: no hook/decided for cwd {cwd}")
if fp == fp_other:
    raise SystemExit(f"notify_uses_payload_cwd: project cwd should change config fingerprint (got {fp})")
PY

OC_OTHER='{"seq":1,"cwd":"'"$OTHER"'","hook":"tool.execute.before","input":{"sessionID":"m20-oc-other","callID":"oc1","tool":"bash"},"output":{"args":{"command":"echo ok"}}}'
OC_PROJ='{"seq":2,"cwd":"'"$PROJ"'","hook":"tool.execute.before","input":{"sessionID":"m20-oc-proj","callID":"oc2","tool":"bash"},"output":{"args":{"command":"echo ok"}}}'
SERVE_OUT="$(printf '%s\n' "$OC_OTHER" "$OC_PROJ" | "$BIN" hook serve --socket "$SOCK" --config "$CFG" --provider=opencode)"
e2e_assert_contains "$SERVE_OUT" '"seq":1' serve_resolves_cwd_per_frame
e2e_assert_contains "$SERVE_OUT" '"seq":2' serve_resolves_cwd_per_frame
# project cwd frame should carry sync-failure deny on wire (non-empty error path)
e2e_assert_matches "$SERVE_OUT" 'deny|block|error|sync pipeline failed' serve_resolves_cwd_per_frame

e2e_daemon_stop
e2e_pass
