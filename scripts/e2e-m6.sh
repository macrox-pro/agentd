#!/usr/bin/env bash
# M6 acceptance: shell / mcp / paths guards + route subset
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

WORKDIR="$(mktemp -d)"
SOCK="$WORKDIR/agentd.sock"
BIN="$WORKDIR/agentd"
CFG="$WORKDIR/agentd.yaml"

cleanup() {
  "$BIN" daemon stop --socket "$SOCK" --timeout 5s >/dev/null 2>&1 || true
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

cat >"$CFG" <<'EOF'
version: 1
guards:
  secrets:
    enabled: false
  shell:
    enabled: true
    deny_patterns: ["rm -rf /"]
    ask_on: [curl]
  mcp:
    enabled: true
    deny_servers: ["untrusted-*"]
  paths:
    enabled: true
    deny_read: ["/etc/shadow"]
    deny_write: ["**/.env"]
dispatch:
  - name: shell-only
    match:
      kind: [tool.pre]
      tools: [Bash]
      provider: ["codex"]
    mode: sync_only
    sync:
      - target: builtin
        guards: [shell]
  - name: all-guards
    match:
      kind: [tool.pre]
    mode: sync_only
    sync:
      - target: builtin
        guards: [shell, mcp, paths]
EOF

go build -o "$BIN" .

"$BIN" config validate --config "$CFG"

SHOW="$("$BIN" config show --config "$CFG" --merged)"
echo "$SHOW" | grep -q 'deny_patterns'
echo "$SHOW" | grep -q 'deny_servers'
echo "$SHOW" | grep -q 'deny_read'

ROUTES_JSON="$("$BIN" dispatch routes --config "$CFG" --json)"
echo "$ROUTES_JSON" | grep -q 'all-guards'
echo "$ROUTES_JSON" | grep -q 'shell-only'
echo "$ROUTES_JSON" | grep -q '"shell"'
echo "$ROUTES_JSON" | grep -q '"mcp"'
echo "$ROUTES_JSON" | grep -q '"paths"'

"$BIN" daemon start --socket "$SOCK" --config "$CFG"

for _ in $(seq 1 50); do
  if "$BIN" daemon status --socket "$SOCK" --json 2>/dev/null | grep -qE '"running"[[:space:]]*:[[:space:]]*true'; then
    break
  fi
  sleep 0.1
done

# shell deny
DENY_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"sudo rm -rf /"}}'
DOUT="$(echo "$DENY_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
echo "$DOUT" | grep -q '"permissionDecision":"deny"'

# shell ask
ASK_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t2","tool_input":{"command":"curl https://example.com"}}'
AOUT="$(echo "$ASK_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
echo "$AOUT" | grep -q '"permissionDecision":"ask"'

# shell clean
CLEAN_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t3","tool_input":{"command":"go test ./..."}}'
COUT="$(echo "$CLEAN_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
test "$COUT" = '{}'

# mcp deny
MCP_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"mcp__untrusted-foo__bar","tool_use_id":"t4","tool_input":{}}'
MOUT="$(echo "$MCP_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
echo "$MOUT" | grep -q '"permissionDecision":"deny"'

# path deny read
PATH_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Read","tool_use_id":"t5","tool_input":{"file_path":"/etc/shadow"}}'
POUT="$(echo "$PATH_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
echo "$POUT" | grep -q '"permissionDecision":"deny"'

# path deny write **/.env
ENV_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Write","tool_use_id":"t6","tool_input":{"file_path":"repo/.env"}}'
EOUT="$(echo "$ENV_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=claude-code)"
echo "$EOUT" | grep -q '"permissionDecision":"deny"'

# subset route: shell-only for codex Bash — ask falls back to deny (no CapAsk)
SUB_PAYLOAD='{"session_id":"s","cwd":"/w","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t7","tool_input":{"command":"curl https://example.com"}}'
SOUT="$(echo "$SUB_PAYLOAD" | "$BIN" hook run --socket "$SOCK" --provider=codex)"
echo "$SOUT" | grep -qiE 'deny|blocked|ask_on'

"$BIN" daemon stop --socket "$SOCK" --timeout 5s

echo "e2e-m6: ok"
