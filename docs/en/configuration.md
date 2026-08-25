# Configuration

> **Language:** [English](./configuration.md) · [Русский](../ru/configuration.md)

Four-layer merge, YAML surface, and how reloads work. Layer/runtime overlay contract: [DESIGN.md §7](../../DESIGN.md#7-configuration-schema).

## Layers (merge order)

| Order | Layer | Location |
|-------|--------|----------|
| 1 | defaults | compiled into the binary |
| 2 | user | `--config` or `~/.agentd.yaml` |
| 3 | project | `.agentd.yaml` walking up from CWD / project root |
| 4 | runtime | daemon-managed overlay (approvals, temporary blocks) |

**Runtime path**

- Unix: `$XDG_STATE_HOME/agentd/runtime.yaml`, else `~/.local/state/agentd/runtime.yaml`
- Windows: `%LOCALAPPDATA%\agentd\runtime.yaml`

Runtime writes are debounced (**500ms**), mode `0600`, atomic rename. Hot path uses `store.Current()` only — no disk I/O per Invoke.

## Top-level YAML keys

From the file schema: `version`, `policy`, `async`, `logging`, `guards`, `approvals`, `blocks`, `dispatch_defaults`, `dispatch`, `trajectory`.

Project files typically carry `guards` / `dispatch`. `approvals` and `blocks` usually land in runtime via CLI/gRPC.

### policy

| Key | Values | Default |
|-----|--------|---------|
| `fail` | `fail_open` \| `fail_closed` | `fail_closed` |
| `unsupported` | `degrade` \| `strict` | `degrade` |
| `ask_fallback` | `deny` \| `no_decision` | `deny` |
| `offline` | `fail_open` \| `fail_closed` | `fail_open` |

When the daemon is unreachable, the hook edge loads local config (defaults ⊕ user ⊕ project(cwd) ⊕ runtime) and applies `policy.offline`. Default `fail_open` encodes a neutral decision (or exit 0 for notify) so agents keep working; `fail_closed` exits **1**. Stderr still prints `daemon not running` in both modes.

### async

| Key | Default |
|-----|---------|
| `queue_capacity` | `1024` |
| `worker_limit` | `8` |
| `target_timeout` | `30s` |
| `on_overflow` | `drop` (`drop` \| `log`) |

Overflow always drops the job and increments `async_dropped_count` on Status; `log` also emits a warn log.

### logging

Daemon operational logging (not the async dispatch `target: log`).

| Key | Default |
|-----|---------|
| `level` | `info` (`debug` \| `info` \| `warn` \| `error`) |
| `file` | `""` → `$XDG_STATE_HOME/agentd/agentd.log` (Windows: `%LOCALAPPDATA%\agentd\agentd.log`) |

`agentd daemon start --foreground` mirrors logs to stderr as well as the file. CLI `--log-level` and `--log-file` override YAML for that process only.

### trajectory

Opt-in session ledger ([Trajectory](./trajectory.md)). Default **off**.

| Key | Default |
|-----|---------|
| `enabled` | `false` |
| `include_raw` | `false` |
| `redact_secret_rules` | `true` |
| `max_event_bytes` | `262144` |
| `queue_capacity` | `1024` |
| `import.claude-code.enabled` | `false` |
| `import.claude-code.path` | `""` (default `~/.claude/projects`) |
| `import.cursor.enabled` | `false` |
| `import.cursor.path` | `""` (prefer CLI `--path`) |
| `import.codex.enabled` | `false` |
| `import.codex.path` | `""` (default `$CODEX_HOME/sessions` or `~/.codex/sessions`) |

When `import.claude-code.enabled` is true, the daemon watches the projects directory and appends new transcript lines asynchronously. CLI `session import` works offline without this flag. Set `include_raw: true` if you need `session replay --policy`.

Overflow increments `trajectory_dropped_count` on Status.

## CLI against config

| Command | Role |
|---------|------|
| `agentd config validate [--config] [--cwd]` | Offline parse + compile |
| `agentd config show [--merged] [--layer user\|project\|runtime] [--cwd]` | Inspect layers |
| `agentd config patch --file DELTA.yaml` | Patch runtime (persisted) |
| `agentd config record-decision …` | Upsert approval ([Approvals](./approvals.md)) |

## Reload

- User/project file changes: fsnotify + debounce → re-merge
- `agentd daemon reload`: force re-merge from disk
- Runtime patch / RecordDecision: in-memory update + debounced flush

Status exposes `generation` and merged `fingerprint` after each successful compile.

See also: [Guards](./guards.md), [Dispatch](./dispatch.md), [CLI](./cli.md).
