# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**M16 follow-up: Codex transcript token fallback** — in progress.

- Codex `Stop` billing tokens via transcript tail-scan when hook raw has no usage
- Files: `internal/trajectory/statistics/extract/{extract,codex}.go`, `collector.go`, `from_events.go`, `statistics.go`, tests, docs

**Next:** Claude transcript token investigation (separate PR).

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-08-31 | M16 Codex tokens | Transcript tail fallback on Stop; hook usage wins when present |
| 2026-08-31 | v0.0.6 | Prometheus metrics + trajectory token stats; CHANGELOG + docs version bump |
| 2026-08-31 | daemon tests | Isolate unit tests from production socket, agentd.log, runtime.yaml |
| 2026-08-31 | metrics | Prometheus scrape HTTP + Observer histograms + reload hook |
| 2026-08-29 | v0.0.5 | Trajectory counters + session import `--out`; CHANGELOG + docs version bump |

## Verify (repo green)

```bash
make lint && make intent-check && make docs-check && make test
```

## Blockers

(none)
