# Configuration

> **Language:** [English](./configuration.md) · [Русский](../ru/configuration.md)

How agentd merges YAML layers into one effective config and where files live on disk.

agentd builds one effective config from **four layers**. Later layers override earlier ones. YAML keys and reload: below. Layer contract: [DESIGN.md §7](../../DESIGN.md#7-configuration-schema).

## State directory

User config is the file `~/.agentd.yaml` — not a `~/.agentd/` folder. Data the daemon writes (runtime overlay, log, session ledger) is **state**: regenerable, not for git. On Linux/macOS that lives under `$XDG_STATE_HOME` (unset → `~/.local/state`). Windows: `%LOCALAPPDATA%\agentd\`. The socket is **not** state — it lives under `$XDG_RUNTIME_DIR` and is removed on logout.

| Kind | Default location |
|------|------------------|
| User config | `~/.agentd.yaml` (file) |
| State directory | `$XDG_STATE_HOME/agentd/` else `~/.local/state/agentd/` (Windows: `%LOCALAPPDATA%\agentd\`) |
| → runtime overlay | `runtime.yaml` (daemon only) |
| → operational log | `agentd.log` |
| → session ledger | `sessions/<provider>/<session_id>.jsonl` (when enabled) |
| IPC socket (not state) | `$XDG_RUNTIME_DIR/agentd/agentd.sock` (macOS fallback `~/Library/Caches/agentd/`; Linux `~/.local/run/agentd/`; else temp) — [DESIGN.md §5](../../DESIGN.md#5-transport) |
| → PID file | `agentd.pid` (next to the socket in the runtime directory) |
| → lock file | `agentd.lock` (next to the socket; one daemon per user) |

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

| Order | Layer | Where it lives |
|-------|--------|----------------|
| 1 | defaults | Built into the binary |
| 2 | user | `--config` or `~/.agentd.yaml` |
| 3 | project | `.agentd.yaml` walking up from hook **payload cwd** (see below) or CLI `--cwd` |
| 4 | runtime | Daemon-managed file (approvals, temporary blocks) |

**Runtime path:** `runtime.yaml` in the [state directory](#state-directory).

Runtime writes are debounced (**500ms**), mode `0o600`, atomic rename. Each hook call reads the in-memory snapshot — no disk I/O per call.

### Hook cwd and project layer

When the **daemon is running**, `hook run`, `hook notify`, and `hook serve` send a **cwd** with each `Invoke`:

1. Top-level `cwd` in the event JSON (including Codex notify argv payload), when set
2. Else `workspace_roots[0]` when present (Cursor-style payloads)
3. Else the hook edge process working directory (used when the payload omits both fields above)

The daemon walks up from that path to find `.agentd.yaml` — the same rule as `config show --cwd`. This is usually the agent workspace from the wire payload, **not** necessarily the shell cwd of the parent process.

CLI `config validate|show|get` use the `--cwd` flag instead of hook JSON. See [CLI → config](./cli.md#config).

## Top-level YAML keys

From the file schema: `version`, `policy`, `async`, `logging`, `guards`, `approvals`, `blocks`, `dispatch_defaults`, `dispatch`, `trajectory`, `metrics`.

Project files typically carry `guards` / `dispatch`. `approvals` and `blocks` usually land in runtime via CLI/gRPC.

### policy

| Key | Values | Default | Meaning |
|-----|--------|---------|---------|
| `fail` | `fail_open` \| `fail_closed` | `fail_closed` | On **daemon** sync pipeline **errors** returned to the engine (for example sync budget deadline/cancel, or a sync target that surfaces an error). Normal guard **deny** is a decision, not remapped by `fail`. gRPC `on_error` handles many peer failures before `policy.fail` applies ([Dispatch](./dispatch.md)) |
| `ask_fallback` | `deny` \| `no_decision` | `deny` | When the agent cannot ask the user: **deny** blocks; **no_decision** returns a neutral allow |
| `offline` | `fail_open` \| `fail_closed` | `fail_open` | When the **hook edge** cannot reach the daemon (see below) |

`policy.fail` applies only on the **daemon** path. `policy.offline` applies only when the hook edge runs locally because dial/invoke failed.

When the daemon is unreachable, the hook edge loads local config (defaults merged with user merged with project(cwd) merged with runtime) and applies `policy.offline`. Default `fail_open` encodes a neutral decision (or exit 0 for notify) so agents keep working; `fail_closed` exits **1**. Stderr still prints `daemon not running` in both modes.

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

Session ledger ([Trajectory](./trajectory.md)). Default **on**.

| Key | Default |
|-----|---------|
| `enabled` | `true` |
| `statistics` | `true` (requires `enabled`; gates daemon rollup + `session stats`) |
| `include_raw` | `true` |
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

### metrics

Opt-in Prometheus scrape endpoint ([Operations → Metrics](./operations.md#prometheus-metrics)). Default **off**.

| Key | Default |
|-----|---------|
| `enabled` | `false` |
| `listen` | `127.0.0.1:2112` (`host:port`; required when enabled) |

When enabled, the daemon serves `/metrics` on loopback TCP at start time only. Changing `metrics.listen` or toggling `enabled` requires **`agentd daemon stop`** then **`agentd daemon start`** — `daemon reload` does not rebind the metrics HTTP listener. CLI `--metrics-listen` enables metrics for that process and overrides `listen`.

Binding `0.0.0.0` is allowed but exposes metrics on all interfaces — prefer loopback unless you understand the risk.

## Feature toggles

Use `agentd config enable|disable|get FEATURE` to flip curated booleans without hand-editing YAML. Full command reference: [CLI → config](./cli.md#config).

| Behavior | Detail |
|----------|--------|
| Layers written | User (`--config` / `~/.agentd.yaml`) or project (`.agentd.yaml` under `--cwd`) only |
| Runtime overlay | **Not** modified — use `config patch` for temporary runtime overrides |
| Missing user file | `config enable` creates the same bootstrap shape as `daemon start` |
| Reload | User/project file changes → config file watcher debounce; or `agentd daemon reload` |
| `config get` | Merges defaults with user with project; **excludes** runtime overlay |
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

- User/project file changes: file watcher + debounce → re-merge
- `agentd daemon reload`: force re-merge from disk
- Runtime patch / RecordDecision: in-memory update + debounced flush

Status exposes `generation` and merged `fingerprint` after each successful compile.

See also: [Guards](./guards.md), [Dispatch](./dispatch.md), [CLI](./cli.md).
