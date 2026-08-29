# Changelog

## Unreleased

### Highlights

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
