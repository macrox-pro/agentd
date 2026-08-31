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
| `agent.stop` | `Stop`, `SessionEnd`, etc. | Agent session ending |
| `tool.post` | `PostToolUse`, etc. | After a tool finished |
| `notification` | Codex `notify`, observe-only frames | Fire-and-forget observation |
| `other` | Anything else | Default async-only handling |

## Modes

| Mode | Behavior |
|------|----------|
| `sync_only` | Sync chain → decision |
| `async_only` | Enqueue async; neutral wire decision |
| `parallel` | Sync + async start together; async never blocks the response |
| `after_sync` | Async after sync, with sync outcome |
| `sync_then_async` | Alias of `after_sync` |

Named `dispatch:` routes match top-down; first hit wins.

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
  tool.pre:         { mode: parallel }
  prompt.submitted: { mode: sync_only }
  agent.stop:       { mode: sync_then_async }
  tool.post:        { mode: parallel }
  notification:     { mode: async_only }
  other:            { mode: async_only }
```

## Sync timeout budget

Effective sync budget:

`min(provider_timeout − 10%, route.sync_timeout)` when `sync_timeout` is set; otherwise provider timeout minus 10%.

- If Invoke carries a deadline → that duration is the provider timeout.
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
