# Dispatch

> **Language:** [English](./dispatch.md) · [Русский](../ru/dispatch.md)

Which checks and side effects run for an event, and in what order. Detail: [DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine).

## Modes

| Mode | Behavior |
|------|----------|
| `sync_only` | Sync chain → decision |
| `async_only` | Enqueue async; neutral wire decision |
| `parallel` | Sync + async start together; async never blocks the response |
| `after_sync` | Async after sync, with sync outcome |
| `sync_then_async` | Alias of `after_sync` |

Kind defaults live under `dispatch_defaults:` (overridable). Named `dispatch:` routes match top-down; first hit wins.

## Targets

| Target | Sync | Async |
|--------|------|-------|
| `builtin` | Runner.Decide / guards | observe |
| `grpc` | yes | yes |
| `http` | — | yes |
| `exec` | **no (v1)** | yes |
| `log` | — | yes |
| `file` | — | yes |

v1 does **not** support sync `exec` JSON decisions (DESIGN §11).

## Route fields

```yaml
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
      provider: ["*"]
    mode: parallel
    sync_timeout: 25s          # optional route cap
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

Match: `kind`, `provider`, `tools` (optional lists).

## Sync timeout budget

Effective sync budget:

`min(provider_timeout − 10%, route.sync_timeout)` when `sync_timeout` is set; otherwise provider timeout minus 10%.

- If Invoke carries a deadline → that duration is the provider timeout.
- Else kind defaults: `tool.pre` / `prompt.submitted` → **30s**; other kinds → **5s** (aligned with install HookSpec).

Per-target gRPC `timeout` is clamped to remaining context budget.

## Async queue

Defaults: capacity `1024`, workers `8`, `target_timeout` `30s`, `on_overflow: drop`.

Full queue → drop job, bump Status `async_dropped_count` (`log` mode also warns). Async must not block the sync response.

Inspect compiled routes offline:

```bash
agentd dispatch routes --json
```

See also: [Configuration](./configuration.md), [Operations](./operations.md).
