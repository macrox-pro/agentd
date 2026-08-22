# CLI reference

> **Language:** [English](./cli.md) · [Русский](../ru/cli.md)

Commands and flags as implemented under `cmd/`. Narrative rationale: [DESIGN.md §6](../../DESIGN.md#6-cli-reference).

## Persistent flags

| Flag | Default | Meaning |
|------|---------|---------|
| `--config` | `~/.agentd.yaml` | User config path |
| `--socket` | OS default | Daemon IPC endpoint |
| `-v` / `--verbose` | off | Extra stderr (never hook stdout) |

## daemon

| Command | Flags | Notes |
|---------|-------|-------|
| `daemon start` | `--foreground`, `--log-level`, `--log-file` | Detach by default; waits until Health succeeds; logs to state-dir file |
| `daemon stop` | `--timeout` (`10s`) | gRPC Shutdown, then SIGTERM fallback |
| `daemon status` | `--json` | Runtime snapshot ([Operations](./operations.md)) |
| `daemon reload` | — | Force config re-merge |

## hook

Thin edge: decode → gRPC Invoke → encode. No policy in the CLI.

| Command | Flags | Notes |
|---------|-------|-------|
| `hook run` | `--provider` (required), `--argv-payload`, `--timeout` (`0` = unset) | Stdin (or argv) hooks |
| `hook notify` | `--provider`, `--timeout` | Codex notify (argv JSON) |
| `hook serve` | `--provider`, `--timeout` | OpenCode NDJSON; provider must be `opencode` |

If dial/Invoke fails: stderr `daemon not running`, exit **1**. Do not write debug to stdout on the hook path.

### agenthooks (hidden)

**Why:** `agentd install` uses [agenthooks/install](https://github.com/speakeasy-api/agenthooks), which writes provider hook configs with `agentd agenthooks …`, not `hook`. agentd registers matching hidden subcommands so those configs work unchanged. Prefer `hook run` / `hook serve` / `hook notify` in docs and manual hook settings — same flags, same wire path (`cmd/hook.go`).

Install writes `agentd agenthooks run|notify|serve --provider=…`. Same behavior as `hook …`. `agenthooks serve` defaults `--provider` to `opencode`.

## config

| Command | Flags |
|---------|-------|
| `config validate` | `--cwd` |
| `config show` | `--merged`, `--layer user\|project\|runtime`, `--cwd` |
| `config patch` | `--file` (required) |
| `config record-decision` | `--fingerprint` (required), `--scope` (`project` default), `--project-root`, `--session-id`, `--expires-at` (RFC3339) |

## install

| Flag | Default |
|------|---------|
| `--provider` | required |
| `--scope` | `project` (`user`, `plugin`) |
| `--dir` | CWD |

## dispatch

| Command | Flags |
|---------|-------|
| `dispatch routes` | `--json`, `--cwd` |

Offline compile of defaults ⊕ user ⊕ optional project (no daemon required).

## session

Trajectory ledger inspect/export ([Trajectory](./trajectory.md)). Offline — reads `$XDG_STATE_HOME/agentd/sessions/`. **Exception:** `session subscribe` requires a running daemon.

| Command | Flags |
|---------|-------|
| `session list` | `--provider`, `--json` (includes `importer_status`) |
| `session show SESSION_ID` | `--provider` (required), `--json` |
| `session export` | `--provider`, `--session`, `--out` |
| `session search` | `--provider`, `--session`, `--kind` (repeatable), `--source`, `--query`, `--limit`, `--json` |
| `session import` | `--provider` (required), `--session`, `--path`, `--dry-run`, `--json` |
| `session replay` | `--policy` (required), `--provider`, `--session`, `--seq`, `--json` |
| `session fork` | `--provider`, `--session`, `--new-session`, `--at-seq`, `--json` |
| `session subscribe` | `--provider`, `--session`, `--source`, `--json` (live firehose; daemon required) |

`session search` scans JSONL line-by-line (O(total bytes); no index). `session import`: Claude Code and Codex `supported`; Cursor `partial` (prefer `--path`); others explicit `none`. `session replay --policy` needs `include_raw` at record time. `session fork` is audit lineage only (source immutable). `session subscribe` is live-only from dial time — use show/export for history.

See also: [Getting started](./getting-started.md), [Providers](./providers.md).
