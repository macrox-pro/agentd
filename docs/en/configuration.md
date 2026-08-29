# Configuration

> **Language:** [English](./configuration.md) · [Русский](../ru/configuration.md)

Four-layer merge, YAML surface, and how reloads work. Layer/runtime overlay contract: [DESIGN.md §7](../../DESIGN.md#7-configuration-schema).

## State directory

User config is the file `~/.agentd.yaml` — not a `~/.agentd/` tree. Mutable daemon data (runtime overlay, operational log, session ledger) is **state**: daemon-written, regenerable, and not meant to travel with config backups or git. [XDG Base Directory](https://specifications.freedesktop.org/basedir-spec/latest/) puts that class under `$XDG_STATE_HOME` (unset → `~/.local/state`). Windows uses `%LOCALAPPDATA%\agentd\`. IPC (socket / lock) is **runtime**, not state — `$XDG_RUNTIME_DIR` is session-scoped and cleaned on logout.

| Kind | Default location |
|------|------------------|
| User config | `~/.agentd.yaml` (file) |
| State directory | `$XDG_STATE_HOME/agentd/` else `~/.local/state/agentd/` (Windows: `%LOCALAPPDATA%\agentd\`) |
| → runtime overlay | `runtime.yaml` (daemon only) |
| → operational log | `agentd.log` |
| → trajectory ledger | `sessions/<provider>/<session_id>.jsonl` (when enabled) |
| IPC socket (not state) | `$XDG_RUNTIME_DIR/agentd/agentd.sock` (Darwin fallback `~/Library/Caches/agentd/`; Linux `~/.local/run/agentd/`; else temp) — [DESIGN.md §5](../../DESIGN.md#5-transport) |

## User config bootstrap

On **`agentd daemon start` only**, if the user config file is missing, the daemon writes a minimal bootstrap file (same keys as the [getting-started](./getting-started.md) example) and continues. Read-only commands (`config show`, `config validate`, hooks) **never** create the file.

| Situation | Behavior |
|-----------|----------|
| File missing | Bootstrap written silently; start continues |
| File valid | No change; start continues |
| Invalid YAML or compile error | Stderr: `agentd: invalid user config <path>: …`; start fails; file unchanged |
| Unreadable path / I/O error | Start fails; no invalid-config message |

Fix invalid config offline, then restart:

```bash
agentd config validate --config ~/.agentd.yaml
```

`config show` normalizes YAML output (empty fields omitted; no `null` keys).

## Layers (merge order)

| Order | Layer | Location |
|-------|--------|----------|
| 1 | defaults | compiled into the binary |
| 2 | user | `--config` or `~/.agentd.yaml` |
| 3 | project | `.agentd.yaml` walking up from CWD / project root |
| 4 | runtime | daemon-managed overlay (approvals, temporary blocks) |

**Runtime path:** `runtime.yaml` in the [state directory](#state-directory).

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
| `file` | `""` → `agentd.log` in the [state directory](#state-directory) |

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

## Feature toggles

Use `agentd config enable|disable|get FEATURE` to flip curated booleans without hand-editing YAML ([CLI](./cli.md#config)).

| Behavior | Detail |
|----------|--------|
| Layers written | User (`--config` / `~/.agentd.yaml`) or project (`.agentd.yaml` under `--cwd`) only |
| Runtime overlay | **Not** modified — use `config patch` for temporary runtime overrides |
| Missing user file | `config enable` creates the same bootstrap shape as `daemon start` |
| Reload | User/project file changes → fsnotify debounce; or `agentd daemon reload` |
| `config get` | Merges defaults ⊕ user ⊕ project; **excludes** runtime overlay |
| Idempotent | Re-enable / disable when already at effective value → exit 0, no write |
| YAML round-trip | Marshal may drop hand-written comments in touched files |
| Secrets guard | Not a curated toggle — edit `guards.secrets` in YAML |
| Project read vs write | `config get` walks up for `.agentd.yaml`; `enable`/`disable` (project scope) writes only under `--cwd` |

**Stdout examples:**

```text
agentd config get trajectory
trajectory: on (user)

agentd config enable guard-shell    # from repo root; default project scope
guard-shell: enabled (project /path/to/repo/.agentd.yaml)

agentd config enable trajectory     # already on
trajectory: already enabled (user /home/you/.agentd.yaml)
```

See also: [Trajectory](./trajectory.md#enable), [Guards](./guards.md#enable-via-cli).

## CLI against config

| Command | Role |
|---------|------|
| `agentd config validate [--config] [--cwd]` | Offline parse + compile |
| `agentd config show [--merged] [--layer user\|project\|runtime] [--cwd]` | Inspect layers |
| `agentd config enable\|disable\|get FEATURE` | Curated persistent toggles ([Feature toggles](#feature-toggles)) |
| `agentd config patch --file DELTA.yaml` | Patch runtime (persisted) |
| `agentd config record-decision …` | Upsert approval ([Approvals](./approvals.md)) |

## Reload

- User/project file changes: fsnotify + debounce → re-merge
- `agentd daemon reload`: force re-merge from disk
- Runtime patch / RecordDecision: in-memory update + debounced flush

Status exposes `generation` and merged `fingerprint` after each successful compile.

See also: [Guards](./guards.md), [Dispatch](./dispatch.md), [CLI](./cli.md).
