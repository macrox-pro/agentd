# CLI reference

> **Language:** [English](./cli.md) · [Русский](../ru/cli.md)

Commands and flags as implemented under `cmd/`. Architecture notes: [DESIGN.md §6](../../DESIGN.md#6-cli-reference).

## Persistent flags

| Flag | Default | Meaning |
|------|---------|---------|
| `--config` | `~/.agentd.yaml` | User config path |
| `--socket` | OS default | Daemon IPC endpoint |
| `-v` / `--verbose` | off | Extra stderr (never hook stdout) |

## version

Print the CLI version. Prefer goreleaser ldflags; otherwise `go install` module version from BuildInfo; local devel may be `dev` or `dev+<shortrev>`. Does not contact the daemon.

Daemon process version: field `version` on `agentd daemon status`.

| Command | Flags | Notes |
|---------|-------|-------|
| `version` | — | CLI version on stdout |

## daemon

| Command | Flags | Notes |
|---------|-------|-------|
| `daemon start` | `--foreground`, `--log-level`, `--log-file` | Detach by default; waits until Health succeeds; logs to [state directory](./configuration.md#state-directory) `agentd.log` |
| `daemon stop` | `--timeout` (`10s`) | gRPC Shutdown, then SIGTERM fallback |
| `daemon status` | `--json` | Daemon health + `autostart` block ([Operations](./operations.md#autostart-at-login)) |
| `daemon reload` | — | Force config re-merge |
| `daemon enable` | — | Registers login autostart and starts the daemon if down. May exit with an error while autostart is already enabled — [Operations → Autostart at login](./operations.md#autostart-at-login) |
| `daemon disable` | — | Removes login autostart only; does **not** stop a running daemon |

## hook

Thin edge: decode → gRPC Invoke → encode. Full Decide/guards stay in the daemon. When the daemon is unreachable, the edge applies `policy.offline` from local config.

| Command | Flags | Notes |
|---------|-------|-------|
| `hook run` | `--provider` (required), `--argv-payload`, `--timeout` (`0` = unset) | Stdin (or argv) hooks |
| `hook notify` | `--provider`, `--timeout` | Codex notify (argv JSON) |
| `hook serve` | `--provider`, `--timeout` | OpenCode NDJSON; provider must be `opencode` |

If dial/Invoke fails: stderr `daemon not running`, then `policy.offline` (`fail_open` default → exit 0 / neutral wire; `fail_closed` → exit **1**). Do not write debug to stdout on the hook path.

### agenthooks (hidden)

**Why:** `agentd install` uses [agenthooks/install](https://github.com/speakeasy-api/agenthooks), which writes provider hook configs with `agentd agenthooks …`, not `hook`. agentd registers matching hidden subcommands so those configs work unchanged. Prefer `hook run` / `hook serve` / `hook notify` in docs and manual hook settings — same flags, same wire path (`cmd/hook.go`).

Install writes `agentd agenthooks run|notify|serve --provider=…`. Same behavior as `hook …`. `agenthooks serve` defaults `--provider` to `opencode`.

## config

Curated feature toggles (`enable` / `disable` / `get`) write **user or project** YAML only — not the runtime overlay. Use `config patch` for temporary runtime overrides. **`daemon enable`** is login autostart, not a config toggle.

| Command | Flags | Notes |
|---------|-------|-------|
| `config validate` | `--cwd` | Offline parse + compile |
| `config show` | `--merged`, `--layer user\|project\|runtime`, `--cwd` | Inspect layers |
| `config enable FEATURE` | `--scope user\|project`, `--cwd` | Persist `enabled: true` (creates user bootstrap if missing) |
| `config disable FEATURE` | `--scope user\|project`, `--cwd` | Persist `enabled: false` |
| `config get FEATURE` | `--cwd` | Effective on/off + winning layer (`default` \| `user` \| `project`); runtime excluded |
| `config patch` | `--file` (required) | Patch runtime overlay (daemon required) |
| `config record-decision` | `--fingerprint` (required), `--scope` (`project` default), `--project-root`, `--session-id`, `--expires-at` (RFC3339) | Upsert approval |

**Features:** `trajectory`, `trajectory-raw`, `guard-shell`, `guard-mcp`, `guard-paths`. Default `--scope`: user for trajectory toggles; project for guard toggles. Offline; running daemon reloads via fsnotify.

| Feature | YAML path | Default scope |
|---------|-----------|---------------|
| `trajectory` | `trajectory.enabled` | `user` |
| `trajectory-raw` | `trajectory.include_raw` | `user` |
| `guard-shell` | `guards.shell.enabled` | `project` |
| `guard-mcp` | `guards.mcp.enabled` | `project` |
| `guard-paths` | `guards.paths.enabled` | `project` |

**Output:** `get` prints `FEATURE: on|off (SOURCE)` where `SOURCE` is `default`, `user`, or `project`. `enable` / `disable` print `FEATURE: enabled|disabled (SCOPE PATH)`; idempotent re-run prints `already enabled|disabled` and exits 0.

```bash
# trajectory (default scope user → ~/.agentd.yaml)
agentd config enable trajectory
# trajectory: enabled (user /home/you/.agentd.yaml)

agentd config get trajectory
# trajectory: on (user)

# guards (default scope project → .agentd.yaml under --cwd)
cd /path/to/repo
agentd config enable guard-shell
# guard-shell: enabled (project /path/to/repo/.agentd.yaml)
```

**Project path asymmetry:** `config get --cwd DIR` walks up from `DIR` to find `.agentd.yaml` (same as hook merge). `config enable|disable` with project scope writes **only** `DIR/.agentd.yaml` — it does not update a parent repo config. Run from the repo root (or pass `--cwd` there) when you intend project-level guards.

**Not toggleable via CLI:** `guards.secrets` (and other non-boolean knobs) — edit YAML or use `config show` / `config patch`. Only the five features above are curated toggles.

## install

| Flag | Default |
|------|---------|
| `--provider` | required |
| `--scope` | `project` (`user`, `plugin`) |
| `--global` | false — same as `--scope=user` |
| `--dir` | `scope=project`: CWD (codex: `./.codex`); `scope=user`: agent home (e.g. `~/.cursor`); `scope=plugin`: required |

`--global` conflicts with an explicit `--scope` other than `user`. On success, prints provider, scope, install root, and per-file `create` / `update` / `unchanged` with absolute paths.

## dispatch

| Command | Flags |
|---------|-------|
| `dispatch routes` | `--json`, `--cwd` |

Offline compile of defaults ⊕ user ⊕ optional project (no daemon required).

## session

Trajectory ledger inspect/export ([Trajectory](./trajectory.md)). Offline — reads `sessions/` under the [state directory](./configuration.md#state-directory). **Exception:** `session subscribe` requires a running daemon.

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
