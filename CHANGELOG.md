# Changelog

## Unreleased

## [v0.0.9-beta] — 2026-09-02

Policy reliability on the daemon path and milestone e2e wire coverage.

### Fixed

- **policy.fail on daemon path** — sync pipeline errors (grpc timeout/cancel, guard failures) now map through `policy.fail` in `dispatch.Engine` instead of returning neutral allow from the server.
- **policy.ask_fallback** — when a guard would Ask but the event lacks `CapAsk`, `ask_fallback` (`deny` vs `no_decision`) is honored instead of always denying.
- **hook notify/serve Cwd** — `InvokeRequest.Cwd` is forwarded so project `.agentd.yaml` applies on notify and serve paths.
- **ConfigStore project cache** — project snapshot cache hits use `projectsMu` RLock so concurrent sessions do not block on `reloadMu`.
- **Makefile docs-check** — terms/links checker failures now fail the target (was masked by `;` before final echo).
- **CI platform tests** — `platform-test` job on `macos-latest` and `windows-latest` for config/daemon/transport packages.

### Removed

- **policy.unsupported** — YAML key is ignored for compatibility; was parsed but never applied on the daemon path.

### Testing

- **E2E milestones** — `e2e-m15` (`trajectory stats` / `session stats`, gates, **M19** Cursor per-generation token sum via `cursor_two_stops_sum_tokens`); `e2e-m16` (Prometheus `/metrics`); `e2e-m18` (non-interactive TUI gate); `e2e-m20` (`policy.fail`, `ask_fallback`, notify/serve `Cwd` on wire).

## [v0.0.8-beta] — 2026-09-01

Cursor trajectory stats: per-generation token sum.

### Fixed

- **Cursor trajectory stats** — `trajectory stats` and offline `session stats` now sum billing tokens from each Cursor `stop` hook (per generation). Multi-stop sessions report higher totals than the previous per-session delta aggregation.

## [v0.0.7] — 2026-08-31

Setup wizard, doctor, auto-install, and trajectory defaults on.

### Highlights

- **Setup wizard** — `agentd setup` and interactive bare `agentd install` on TTY; `AGENTD_NO_TUI` / `CI` bypass
- **Doctor + auto-install** — `agentd doctor` (read-only); `agentd install --all-detected` (plan-only default, `--yes` to apply)
- **Codex trajectory stats** — billing token extraction via transcript tail-read on `Stop` when hook raw carries no usage
- **Trajectory defaults on** — compile defaults set `trajectory.enabled`, `include_raw`, and `trajectory.statistics` to true; `redact_secret_rules` stays true. Disable with `agentd config disable trajectory` or YAML `trajectory.enabled: false`.

## [v0.0.6] — 2026-08-31

Prometheus observability and richer trajectory token stats.

### Highlights

- **Prometheus metrics** — opt-in loopback `/metrics` (default off; `127.0.0.1:2112`); runtime gauges, invoke/async histograms, config reload counter; `daemon status --json` field `metrics_listen`; CLI `--metrics-listen`
- **`trajectory stats`** — Cursor `stop` billing token extract + per-session delta aggregation; CLI `--json` uses enum names and numeric counters (not proto int map keys)

## [v0.0.5] — 2026-08-29

Daemon-lifetime trajectory counters and offline session stats.

### Highlights

- **`trajectory stats`** — in-memory daemon rollup via `TrajectoryService.Statistics`; `since` is daemon start; counters reset on restart; optional `--provider` filter
- **`session stats`** — offline JSONL scan for one session (`--provider` required); no daemon needed
- **`trajectory.statistics`** — both surfaces require `trajectory.enabled` and this key (default off); toggle `trajectory-statistics`
- **`session import --out`** — emit parsed transcript events as JSONL to stdout (`-`) or a file without writing the session ledger or import checkpoint; summary on stderr when `--out` is set

## [v0.0.4] — 2026-08-29

Login autostart and curated config toggles for trajectory and guards.

### Highlights

- **`config enable` / `disable` / `get`** — curated offline toggles for trajectory and guards (user/project YAML); distinct from `daemon enable` (login autostart)
- **`daemon enable` / `disable`** — user-level login autostart (systemd user unit, macOS LaunchAgent, Windows Task Scheduler); `daemon status --json` always includes `autostart`
- **`e2e-m14.sh`** — isolated HOME: enable starts daemon + autostart on → disable keeps daemon running → manifest removed

## [v0.0.3] — 2026-08-26

Hook edge honors `policy.offline` when the daemon is down; research corpus for provider design.

### Highlights

- **`policy.offline` on hook edge** — when the daemon is unreachable, `hook run|notify|serve` loads local config (defaults ⊕ user ⊕ project(cwd) ⊕ runtime) via `OfflineFor` and applies `policy.offline`; default `fail_open` → exit 0 + neutral wire (agents keep working); `fail_closed` → exit **1**; stderr still prints `daemon not running`
- **`hookclient.DialReady`** — lazy gRPC dial with Health check before Invoke
- **OpenCode serve offline cache** — mid-stream Invoke failure caches offline mode for the rest of the NDJSON session (stderr once)
- User docs EN/RU updated: configuration, troubleshooting, cli, getting-started, providers

### Contributor / research

- Structured research trees under `research/` — Claude Code, Codex, Cursor, Go best practices (verbatim excerpts + indexes)
- DESIGN.md streamlining; documentation cross-reference cleanup

### Explicitly not in v0.0.3

Offline state-dir cache; auto-start daemon; change install FailClosed; new e2e script.

## [v0.0.2] — 2026-08-23

Trajectory hub (M9–M12): live ledger, import, replay/fork, and live Subscribe.

### Highlights

- **Live trajectory ledger** — opt-in append-only session log for all six providers; `session list`, `show`, `export`
- **Search & import** — `session search`; Claude transcript import; Codex rollout JSONL import (**supported**); Cursor import with explicit per-provider status matrix (§14.3)
- **Replay & fork** — `session replay --policy` (dry-run policy re-check from stored raw); log fork for audit lineage
- **Live Subscribe** — gRPC `SessionService.Subscribe` and `agentd session subscribe` firehose from the daemon
- **`schema_version: 1`** frozen on all ledger events (JSONL + stream)
- Trajectory contract docs; honest coverage matrix in user guide

### Architecture (no user-facing behavior change)

- R-series refactor (R1–R11): dispatch/target boundaries, guard registry, provider IDs, test coverage uplift
- `internal/decision` owns proto↔agenthooks Decision mapping; no rename-only passthrough helpers

### Explicitly not in v0.0.2

Historical catch-up on Subscribe (`after_seq`), HTTP webhook mirror, agent-loop resume, git tag automation in-repo.

## [v0.0.1] — 2026-08-20

First version release (M0–M8).

### Highlights

- Four-layer config merge (defaults ⊕ user ⊕ project ⊕ runtime) with ConfigService Get/Patch/RecordDecision
- Guards: secrets, shell, MCP, paths; approvals + temporary blocks with runtime persist
- Sync + async dispatch (builtin, exec async-only, http, grpc, log, file)
- Install + `hook run` / `notify` / `serve` for supported providers
- Cross-platform IPC (unix sockets + Windows SID named pipes)
- Status exposes async queue depth and overflow drop counter
- Provider-aware sync timeout margin (`min(provider_timeout - margin, route.sync_timeout)`)
- Conformance fixtures via agenthookstest; optional `//go:build integration` round-trip
- GitHub Releases binaries via goreleaser (linux/darwin/windows)

### Explicitly not in v0.0.1

Agent auth, transcripts, plugins, hooks DSL, async retry storms, exec sync JSON decisions.
