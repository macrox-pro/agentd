# Operations

> **Language:** [English](./operations.md) · [Русский](../ru/operations.md)

Day-2 controls for the user daemon.

## One daemon per user

Lock + socket prevent a second start. Override IPC with `--socket`.

## Status

```bash
agentd daemon status
agentd daemon status --json
```

JSON when running:

| Field | Meaning |
|-------|---------|
| `running` | bool |
| `socket` | path / pipe |
| `version` | build version (`dev` unless ldflags/tag) |
| `started_at` | RFC3339 UTC |
| `generation` | config generation |
| `fingerprint` | merged config fingerprint |
| `async_queue_depth` | queued async jobs |
| `async_dropped_count` | overflow drops (monotonic) |
| `compiled_route_count` | routes in snapshot |

Human line: `agentd: running (version …, generation …)` — does not print depth/drops (use `--json`).

## Reload / stop

```bash
agentd daemon reload
agentd daemon stop --timeout 10s
```

Stop drains sync then async (up to timeout), then removes socket/PID.

## Logging

Hook path: never debug on stdout. Daemon / async `log` targets use structured logs (stderr / configured sinks).

See also: [Dispatch](./dispatch.md) async overflow, [Configuration](./configuration.md) reload.
