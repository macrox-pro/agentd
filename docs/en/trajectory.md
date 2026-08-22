# Trajectory (session ledger)

> **Language:** [English](./trajectory.md) · [Русский](../ru/trajectory.md)

Opt-in append-only ledger of hook Invokes (`hook/invoked`, `hook/decided`, async meta). **Default off** — payloads may contain secrets.

## Enable

```yaml
trajectory:
  enabled: true
  include_raw: false          # default; set true to enable session replay --policy
  redact_secret_rules: true   # when include_raw: true
  max_event_bytes: 262144
  queue_capacity: 1024
  import:
    claude-code:
      enabled: false            # optional daemon fsnotify import
      path: ""                  # default ~/.claude/projects
    cursor:
      enabled: false
      path: ""                  # optional scan root; prefer CLI --path
    codex:
      enabled: false
      path: ""                  # default $CODEX_HOME/sessions or ~/.codex/sessions
```

Storage: `$XDG_STATE_HOME/agentd/sessions/<provider>/<session_id>.jsonl` (Windows: under `%LOCALAPPDATA%\agentd\sessions\`).

Recording happens on the daemon async path — sync hook latency is unchanged.

## CLI

| Command | Role |
|---------|------|
| `agentd session list [--provider ID] [--json]` | List sessions (offline; `--json` includes `importer_status`) |
| `agentd session show ID --provider ID [--json]` | Print events |
| `agentd session export [--provider ID] [--session ID] [--out PATH]` | Export JSONL |
| `agentd session search [--provider ID] [--query TEXT] …` | Filter ledger (O(n) JSONL scan) |
| `agentd session import --provider ID …` | Append transcript events (`source=transcript`) |
| `agentd session replay --policy --provider ID --session ID` | Dry-run stored Raw through Dispatch Engine |
| `agentd session fork --provider ID --session SRC --new-session DST` | Copy ledger prefix (audit lineage) |
| `agentd session subscribe [--json]` | **Live** stream from daemon (requires running daemon + trajectory.enabled) |

## Subscribe (live stream)

`session subscribe` tails the daemon in-memory ledger from dial time (gRPC `SessionService.Subscribe`). Filters: `--provider`, `--session`, `--source`. History is **not** replayed — use `session show` or `session export`.

- Requires a **running daemon** with `trajectory.enabled`.
- Offline `session import` / `fork` do not publish to Subscribe (no Hub); Claude daemon import watcher does.
- `schema_version: 1` on every event; `ignorable` marks forward-compat hints (readers may skip unknown **types**; Subscribe still delivers transcript events).
- `raw` on the stream follows the same redaction rules as JSONL (`include_raw`, `redact_secret_rules`).
- Global ledger mirror: `trajectory.enabled` records all Invokes — no separate webhook or `target: trajectory` in M12.

## Trajectory contract (v1.1)

| Field | Meaning |
|-------|---------|
| `schema_version` | Frozen at `1` for v1.1 |
| `seq` | Contiguous per session |
| `type` | Event catalog (`hook/invoked`, `transcript/message`, …) |
| `source` | `hook`, `decision`, `transcript`, `system` |
| `ignorable` | Forward-compat: old readers may skip unknown **types** |

Event catalog matches [DESIGN §14.3](../../DESIGN.md#143-event-model-draft-catalog).

## Search

`session search` walks matching session JSONL files line-by-line. There is no search index — acceptable for moderate logs; filter with `--provider`, `--session`, `--kind`, `--source`, `--query`, `--limit`.

## Import

```bash
agentd session import --provider claude-code --session SESSION_ID
agentd session import --provider cursor --path /path/to/transcript.jsonl
agentd session import --provider codex --session SESSION_ID
agentd session import --provider codex --path /path/to/rollout-…-SESSION_ID.jsonl
```

Transcript events append after existing hook events (monotonic `seq`). Re-import skips lines recorded in `<session_id>.import.json` sidecar. Correlation uses `session_id` and `tool_use_id` / `call_id` when present — never merges unrelated runs. Cursor is **partial** (prefer `--path`; never invent thinking or tool outputs). Codex is **supported** via `~/.codex/sessions/**/rollout-*-{session_id}.jsonl` (thinking only from plaintext `agent_reasoning`).

## Policy replay

```bash
agentd session replay --policy --provider claude-code --session s1 --json
```

Re-invokes stored `Raw` through the Dispatch Engine offline. Requires `trajectory.include_raw: true` when events were recorded. Does **not** talk to a live agent or resume an agent loop.

## Fork

```bash
agentd session fork --provider claude-code --session s1 --new-session s1-fork --at-seq 4
```

Copies a prefix into a new session id and appends `session/fork` + `session/end-seed` metadata. Source ledger stays immutable. Audit lineage only — not agent resume.

## Coverage (L0 vs richer tiers)

| Provider | Entrypoint | L0 live | L2 import status | L3 thinking |
|----------|------------|---------|------------------|-------------|
| claude-code | `hook run` | required | **supported** | from session files, not hooks |
| cursor | `hook run --argv-payload` | required | **partial** (`--path`) | often redacted / absent |
| codex | `run` + `hook notify` | required | **supported** (`~/.codex/sessions` rollouts) | plaintext `agent_reasoning` only |
| gemini | `hook run` | required | none | unknown |
| opencode | `hook serve` | required | none | unknown |
| kimi-code | `hook run` | required | none | unknown |

Honest claim: every **supported agent’s hooks** are traceable on one stream; transcript/thinking depth varies by provider ([DESIGN §14.6](../../DESIGN.md#146-provider-support-matrix-all-supported-agents)).

## Status

`trajectory_dropped_count` — monotonic overflow drops when the trajectory queue is full ([Operations](./operations.md)).

See also: [Configuration](./configuration.md), [CLI](./cli.md).
