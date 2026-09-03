# Dispatch

> **Language:** [English](./dispatch.md) · [Русский](../ru/dispatch.md)

How events are routed through the sync path (reply to the agent) and async path (side effects).

The **sync path** decides allow, ask, or deny and shapes the wire reply. The **async path** runs logs, HTTP, exec, and similar work without blocking that reply. Terms: [Glossary](./glossary.md).

Detail: [DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine).

## Event kinds (`kind`)

Wire names in YAML `match:` and `dispatch_defaults`. Provider JSON may use different spellings; agentd normalizes them.

| Wire kind | Provider examples | When it fires |
|-----------|-------------------|---------------|
| `tool.pre` | `PreToolUse` (Claude), `preToolUse` (Cursor), etc. | Before a tool runs |
| `prompt.submitted` | `PromptSubmitted`, `UserPromptSubmit`, etc. | User sent a prompt |
| `agent.stop` | `Stop` | Agent session ending (not `subagent.stop`) |
| `tool.post` | `PostToolUse`, etc. | After a tool finished |
| `tool.error` | `PostToolUseFailure`, etc. | Tool failed |
| `permission.request` | Permission / approval prompts | Agent asked to run a tool |
| `session.start` | `SessionStart` | New agent session |
| `session.end` | `SessionEnd` | Session closed |
| `subagent.start` | `SubagentStart` | Nested agent started |
| `subagent.stop` | `SubagentStop` | Nested agent finished |
| `compact.pre` | `PreCompact` | Context compaction about to run |
| `compact.post` | `PostCompact` | Compaction finished |
| `file.edited` | `afterFileEdit` (Cursor) | File changed (routed if subscribed; not in the default install set) |
| `model.response` | After-agent-thought style frames | Model output (routed if subscribed; not in the default install set) |
| `notification` | Codex `notify`, observe-only frames | Fire-and-forget observation |
| `other` | Unmapped native names | Catch-all when no exact default matches |

Unknown keys in `dispatch_defaults` or `match.kind` fail config compile (daemon start / `config validate`).

Named `dispatch:` routes match first. Then the exact-kind default. Then the default `other` catch-all so a new kind is observed instead of dropped. `async_only` defaults do not take the per-session lock.

## Modes

| Mode | Behavior |
|------|----------|
| `sync_only` | Sync chain → decision |
| `async_only` | Enqueue async; neutral wire decision |
| `parallel` | Sync + async start together; async never blocks the response |
| `after_sync` | Async after sync, with sync outcome |
| `sync_then_async` | Alias of `after_sync` |

Named `dispatch:` routes match first (top-down among user routes). Then the exact-kind default. Then default `other`.

## Targets

| Target | Sync | Async |
|--------|------|-------|
| `builtin` | guards / decision | observe |
| `grpc` | yes | yes |
| `http` | — | yes |
| `exec` | **no** | yes |
| `log` | — | yes |
| `file` | — | yes |

Sync `exec` with JSON decisions is not supported in the current release (DESIGN §11).

## Route fields reference

```yaml
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
      provider: ["*"]
    mode: parallel
    sync_timeout: 25s          # optional route cap (see below)
    sync:
      - target: builtin
        guards: [secrets, shell]
      - target: grpc
        endpoint: unix:///path/to/peer.sock
        timeout: 2s
        on_error: fail_closed  # or fail_open
        merge: first_conclusive
    async:
      - target: log
        level: info
      - target: exec
        command: ["notify", "--"]
        stdin: raw
```

| Field | Where | Meaning |
|-------|-------|---------|
| `sync_timeout` | route | Optional cap on the sync budget (see [Sync timeout budget](#sync-timeout-budget)) |
| `merge` | `grpc` sync target | `first_conclusive` — first non-neutral decision wins |
| `on_error` | `grpc` sync target | `fail_closed` (default) or `fail_open` when the peer fails |
| `endpoint` | `grpc` target | Peer address, for example `unix:///path/to/peer.sock` |
| `stdin` | `exec` async target | `raw` sends the event payload on the command’s stdin |

`match` accepts optional lists: `kind`, `provider`, `tools`.

## Kind defaults

Built-in defaults (override under `dispatch_defaults:` in config):

```yaml
dispatch_defaults:
  tool.pre:            { mode: parallel }
  prompt.submitted:    { mode: sync_only }
  agent.stop:          { mode: sync_then_async }
  tool.post:           { mode: parallel }
  notification:        { mode: async_only }
  other:               { mode: async_only }
  session.start:       { mode: async_only }
  session.end:         { mode: async_only }
  tool.error:          { mode: async_only }
  permission.request:  { mode: async_only }
  subagent.start:      { mode: async_only }
  subagent.stop:       { mode: async_only }
  compact.pre:         { mode: async_only }
  compact.post:        { mode: async_only }
  file.edited:         { mode: async_only }
  model.response:      { mode: async_only }
```

After upgrading, re-run `agentd install` so agents subscribe to the new kinds. `file.edited` and `model.response` are routed when a hook arrives but are not written by the default installer.

## Sync timeout budget

Effective sync budget:

`min(provider_timeout − 10%, route.sync_timeout)` when `sync_timeout` is set; otherwise provider timeout minus 10%.

- If the hook request carries a deadline → that duration is the provider timeout.
- Else kind defaults: `tool.pre` / `prompt.submitted` → **30s**; other kinds → **5s** (aligned with install hook timeouts).

Per-target gRPC `timeout` is clamped to remaining context budget.

## Async queue

Defaults: capacity `1024`, workers `8`, `target_timeout` `30s`, `on_overflow: drop`.

Full queue → drop job, bump Status `async_dropped_count` (`log` mode also warns). Async must not block the sync response.

Inspect compiled routes offline:

```bash
agentd dispatch routes --json
```

See also: [Configuration](./configuration.md), [Operations](./operations.md), [Glossary](./glossary.md).
