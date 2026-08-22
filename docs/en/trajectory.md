# Trajectory (session ledger)

> **Language:** [English](./trajectory.md) · [Русский](../ru/trajectory.md)

Opt-in append-only ledger of hook Invokes (`hook/invoked`, `hook/decided`, async meta). **Default off** — payloads may contain secrets.

## Enable

```yaml
trajectory:
  enabled: true
  include_raw: false          # default; Raw may hold secrets
  redact_secret_rules: true   # when include_raw: true
  max_event_bytes: 262144
  queue_capacity: 1024
```

Storage: `$XDG_STATE_HOME/agentd/sessions/<provider>/<session_id>.jsonl` (Windows: under `%LOCALAPPDATA%\agentd\sessions\`).

Recording happens on the daemon async path — sync hook latency is unchanged.

## CLI

| Command | Role |
|---------|------|
| `agentd session list [--provider ID] [--json]` | List sessions (offline) |
| `agentd session show ID --provider ID [--json]` | Print events |
| `agentd session export [--provider ID] [--session ID] [--out PATH]` | Export JSONL |

## Coverage (L0 vs richer tiers)

| Provider | Entrypoint | L0 live | L2 import | L3 thinking |
|----------|------------|---------|-----------|-------------|
| claude-code | `hook run` | required | M10 | from session files, not hooks |
| cursor | `hook run --argv-payload` | required | M11 partial | often redacted / absent |
| codex | `run` + `hook notify` | required | none until format known | unlikely via hooks |
| gemini | `hook run` | required | none | unknown |
| opencode | `hook serve` | required | none | unknown |
| kimi-code | `hook run` | required | none | unknown |

Honest claim: every **supported agent’s hooks** are traceable on one stream; transcript/thinking depth varies by provider ([DESIGN §14.6](../../DESIGN.md#146-provider-support-matrix-all-supported-agents)).

## Status

`trajectory_dropped_count` — monotonic overflow drops when the trajectory queue is full ([Operations](./operations.md)).

See also: [Configuration](./configuration.md), [CLI](./cli.md).
