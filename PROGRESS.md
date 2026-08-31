# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**M17/M18 setup wizard + doctor + trajectory defaults** — **done** (this session).

- PR1: trajectory compile defaults on (`enabled`, `include_raw`, `statistics`)
- PR2: `doctor`, `install --all-detected`, discovery/plan domain, `e2e-m17`
- PR3: `setup` TUI (`internal/install/tui`), interactive bare `install`, Charm deps
- Convention follow-up: plan tests (`hook_status_stale`, `run_all_partial_error`), `{file}_test.go` layout, `ErrNonInteractive` mapped in `cmd/` (no flag names in `internal/`)

**Next:** ship v0.0.7 tag when ready; Claude transcript token investigation (separate PR).

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-08-31 | M17/M18 | Doctor, `--all-detected`, setup TUI, trajectory default-on |
| 2026-08-31 | M16 Codex tokens | Transcript tail fallback on Stop; hook usage wins when present |
| 2026-08-31 | v0.0.6 | Prometheus metrics + trajectory token stats; CHANGELOG + docs version bump |
| 2026-08-31 | daemon tests | Isolate unit tests from production socket, agentd.log, runtime.yaml |
| 2026-08-31 | metrics | Prometheus scrape HTTP + Observer histograms + reload hook |
| 2026-08-29 | v0.0.5 | Trajectory counters + session import `--out`; CHANGELOG + docs version bump |

## Verify (repo green)

```bash
make lint && make intent-check && make docs-check && make test && make e2e
```

## Blockers

(none)
