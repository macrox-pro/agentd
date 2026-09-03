# Trajectory (session ledger)

> **Language:** [English](./trajectory.md) · [Русский](../ru/trajectory.md)

Chronological log of hook calls and related events for audit and replay.

Terms: [Glossary](./glossary.md).

A chronological log of hook calls (`hook/invoked`, `hook/decided`, async metadata). **On by default.** Payloads can contain secrets; `redact_secret_rules` defaults to true.

## Enable

Trajectory is on by default (compile defaults). Verify or disable without hand-editing YAML ([Feature toggles](./configuration.md#feature-toggles)):

```bash
agentd config get trajectory          # trajectory: on (default)

# Disable when you do not want a session ledger
agentd config disable trajectory

# Raw is on by default; disable if you do not need session replay --policy
agentd config disable trajectory-raw
```

Or edit YAML:

```yaml
trajectory:
  enabled: true               # default on; set false to disable
  include_raw: true           # default on; set false to omit raw payloads
  redact_secret_rules: true   # default on when include_raw is true
  max_event_bytes: 262144
  queue_capacity: 1024
  import:
    claude-code:
      enabled: false            # optional daemon watches projects dir for new transcript lines
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

Offline commands read `sessions/` under the [state directory](./configuration.md#state-directory). `session subscribe` and `trajectory stats` need a running daemon. Full command and flag reference: [CLI → session](./cli.md#session).

## Daemon statistics

```bash
agentd config get trajectory
agentd config get trajectory-statistics
agentd trajectory stats [--provider ID] [--json]
```

Requires a **running daemon**. Counters reset when the daemon process restarts; `since` reflects daemon start time. Optional `--provider` filters the session rollup.

Daemon token totals are extracted from each hook event payload (not gated by `include_raw`):

- **Cursor** — billing tokens only on `agent.stop` (per generation, sum each stop). `subagent.stop` does not add billing. `context_tokens_last` still updates from `preCompact` / `compact.pre`.
- **Codex** — billing tokens from the rollout transcript tail on `agent.stop` when hook raw carries no usage (`transcript_path` in raw).
- Kinds without a proto enum (for example `subagent.*`, `compact.*`, `file.edited`) roll up as **OTHER** in `trajectory stats`.

Offline `session stats` token fields require `include_raw` in the JSONL ledger (Codex transcript fallback needs `transcript_path` in stored raw).

## Subscribe (live stream)

`session subscribe` tails the daemon in-memory ledger from dial time (gRPC `SessionService.Subscribe`). Filters: `--provider`, `--session`, `--source`. History is **not** replayed — use `session show` or `session export`.

- Requires a **running daemon** with `trajectory.enabled`.
- Offline `session import` / `fork` do not publish to Subscribe (no Hub); Claude daemon import watcher does.
- `schema_version: 1` on every event; `ignorable` marks forward-compat hints (readers may skip unknown **types**; Subscribe still delivers transcript events).
- `raw` on the stream follows the same redaction rules as JSONL (`include_raw`, `redact_secret_rules`).
- Global ledger mirror: `trajectory.enabled` records all Invokes — there is no separate webhook or dispatch target for the ledger.

## Trajectory contract

| Field | Meaning |
|-------|---------|
| `schema_version` | Frozen at `1` since the trajectory feature shipped |
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

For persisted ledger + Subscribe history, use default import (no `--out`). Details: [CLI §Import without writing the ledger](./cli.md#import-without-writing-the-ledger).

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

## Coverage tiers

| Tier | Meaning |
|------|---------|
| **L0 Live** | Every hook call → `hook/invoked` + `hook/decided` (required for all six providers). |
| **L1 Correlate** | Stable session and tool ids in events (quality varies by provider; not a separate documentation tier). |
| **L2 Import** | On-disk transcript → `transcript/*` events via `session import`. |
| **L3 Thinking** | Reasoning/thinking lines when the vendor persists them (provider-specific). |

| Provider | Entrypoint | L0 live | L2 import status | L3 thinking |
|----------|------------|---------|------------------|-------------|
| claude-code | `hook run` | required | **supported** | from session files, not hooks |
| cursor | `hook run --argv-payload` | required | **partial** (`--path`) | often redacted / absent |
| codex | `run` + `hook notify` | required | **supported** (`~/.codex/sessions` rollouts) | plaintext `agent_reasoning` only |
| gemini | `hook run` | required | none | unknown |
| opencode | `hook serve` | required | none | unknown |
| kimi-code | `hook run` | required | none | unknown |

Every supported agent’s hooks are traceable on one stream; transcript and thinking depth varies by provider ([DESIGN §14.3](../../DESIGN.md#143-provider-support-matrix)).

## Status

`trajectory_dropped_count` — monotonic overflow drops when the trajectory queue is full ([Operations](./operations.md)).

See also: [Configuration](./configuration.md), [CLI](./cli.md), [Glossary](./glossary.md).
