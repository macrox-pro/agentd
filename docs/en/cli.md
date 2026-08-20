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
| `daemon start` | `--foreground` | Detach by default; waits until Health succeeds |
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

See also: [Getting started](./getting-started.md), [Providers](./providers.md).
