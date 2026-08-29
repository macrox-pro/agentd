# Getting started

> **Language:** [English](./getting-started.md) · [Русский](../ru/getting-started.md)

Bring up a user-level daemon, write a minimal policy, install agent hooks, and confirm Status.

Motivation and pain points: [Why agentd](./why.md).

## 1. Install binary

See [Installation](./installation.md). Quick path:

```bash
go install github.com/macrox-pro/agentd@latest
```

## 2. Start the daemon

One instance per user. Default socket is OS-specific (`--socket` overrides).

```bash
agentd daemon start
agentd daemon status
```

`--foreground` keeps the process attached (dev / process managers). Operational logs go to `agentd.log` in the [state directory](./configuration.md#state-directory) (`--foreground` also mirrors stderr).

If `~/.agentd.yaml` (or `--config`) is missing, **`daemon start` creates a minimal user config** automatically ([details](./configuration.md#user-config-bootstrap)). You can edit it after the daemon is up.

## 3. User config (optional at first start)

Default path: `~/.agentd.yaml` (or `--config`). Skip creating it manually if you already ran `daemon start` — bootstrap matches:

```yaml
version: 1
policy:
  fail: fail_closed
  # offline defaults to fail_open — agents keep working if the daemon is down
  # offline: fail_closed  # lockdown when daemon unreachable
guards:
  secrets:
    enabled: true
    action: ask
```

After edit, fsnotify reloads the user file; `agentd daemon reload` forces a re-merge.

## 4. Install hooks for an agent

From a project directory (example: Claude Code, project scope):

```bash
agentd install --provider=claude-code --scope=project
```

Generated configs call `agentd agenthooks …` (hidden alias of `agentd hook …` — [why](./cli.md#agenthooks-hidden)). Prefer documenting `hook run` / `hook serve` / `hook notify`.

Providers, scopes, and quirks: [Providers](./providers.md).

## 5. Verify

```bash
agentd daemon status --json
```

Expect `"running": true` plus `generation`, `fingerprint`, `async_queue_depth`, `async_dropped_count`.

Trigger a tool call in the agent, or pipe a fixture through the edge:

```bash
echo '{"session_id":"s","cwd":"/tmp","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}' \
  | agentd hook run --provider=claude-code
```

Clean tool.pre with defaults usually encodes a provider no-op (Claude: `{}`).

## Next

- [Configuration](./configuration.md) — layers and schema
- [Guards](./guards.md) / [Dispatch](./dispatch.md) — policy and routing
- [Approvals](./approvals.md) — Ask once, then allow
- [CLI](./cli.md) — full flag list
