# Trajectory (session ledger)

> **Language:** [English](./trajectory.md) · [Русский](../ru/trajectory.md)

Opt-in append-only ledger of hook Invokes (`hook/invoked`, `hook/decided`, async meta). **Default off** — payloads may contain secrets.

## Enable

Without hand-editing YAML ([Feature toggles](./configuration.md#feature-toggles)):

```bash
agentd config enable trajectory
agentd config get trajectory          # trajectory: on (user)

# Raw payloads for session replay --policy (can enable before or after trajectory)
agentd config enable trajectory-raw
```

Or edit YAML:

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

Storage: `sessions/<provider>/<session_id>.jsonl` under the [state directory](./configuration.md#state-directory).

Recording happens on the daemon async path — sync hook latency is unchanged.

## CLI

| Command | Role |
|---------|------|
| `agentd session list [--provider ID] [--json]` | List sessions (offline; `--json` includes `importer_status`) |
| `agentd session show ID --provider ID [--json]` | Print events |
| `agentd session export [--provider ID] [--session ID] [--out PATH]` | Export JSONL |
| `agentd session search [--provider ID] [--query TEXT] …` | Filter ledger (O(n) JSONL scan) |
| `agentd session import --provider ID …` | Append transcript events (`source=transcript`) or `--out` parse-only JSONL emit |
| `agentd session replay --policy --provider ID --session ID` | Dry-run stored Raw through Dispatch Engine |
| `agentd session fork --provider ID --session SRC --new-session DST` | Copy ledger prefix (audit lineage) |
| `agentd session stats ID --provider ID [--json]` | Offline session ledger statistics (requires `trajectory.statistics`) |
| `agentd session subscribe [--json]` | **Live** stream from daemon (requires running daemon + trajectory.enabled) |

## Daemon statistics

```bash
agentd config enable trajectory
agentd config enable trajectory-statistics
agentd trajectory stats [--provider ID] [--json]
```

Requires a **running daemon**. Counters reset on daemon restart; `since` reflects daemon start time. Optional `--provider` filters the rollup. Token totals appear only when `trajectory.include_raw` was true at record time.

## Subscribe (live stream)

`session subscribe` tails the daemon in-memory ledger from dial time (gRPC `SessionService.Subscribe`). Filters: `--provider`, `--session`, `--source`. History is **not** replayed — use `session show` or `session export`.

- Requires a **running daemon** with `trajectory.enabled`.
- Offline `session import` / `fork` do not publish to Subscribe (no Hub); Claude daemon import watcher does.
- `schema_version: 1` on every event; `ignorable` marks forward-compat hints (readers may skip unknown **types**; Subscribe still delivers transcript events).
- `raw` on the stream follows the same redaction rules as JSONL (`include_raw`, `redact_secret_rules`).
- Global ledger mirror: `trajectory.enabled` records all Invokes — no separate webhook or `target: trajectory` in M12.

## Trajectory contract

| Field | Meaning |
|-------|---------|
| `schema_version` | Frozen at `1` for v0.0.2 |
| `seq` | Contiguous per session |
| `type` | Event catalog (`hook/invoked`, `transcript/message`, …) |
| `source` | `hook`, `decision`, `transcript`, `system` |
| `ignorable` | Forward-compat: old readers may skip unknown **types** |

Event catalog matches [DESIGN §14.2](../../DESIGN.md#142-event-catalog).

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

### Import to stdout or file (`--out`)

Preview or transcode pipelines without mutating the ledger:

```bash
agentd session import --provider claude-code --path /path/to/session.jsonl --out -
agentd session import --provider codex --session SESSION_ID --out /tmp/events.jsonl
agentd session import --provider claude-code --session s1 --out - 2>/dev/null | wc -l
```

- Does not append to `sessions/…jsonl` or update `<session>.import.json`.
- Still respects incremental import (reads checkpoint for `startIndex` if sidecar exists).
- Event `seq` values match what a normal import would assign (including continuation after existing hook events).
- With `--out`, human summary moves to stderr; use `--json` for machine-readable summary there.

For persisted ledger + Subscribe history, use default import (no `--out`). Details: [CLI §session import --out](./cli.md#session-import-out).

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

Honest claim: every **supported agent’s hooks** are traceable on one stream; transcript/thinking depth varies by provider ([DESIGN §14.3](../../DESIGN.md#143-provider-support-matrix)).

## Status

`trajectory_dropped_count` — monotonic overflow drops when the trajectory queue is full ([Operations](./operations.md)).

See also: [Configuration](./configuration.md), [CLI](./cli.md).
