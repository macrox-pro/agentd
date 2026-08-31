# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**Prometheus metrics** — implemented (opt-in loopback `/metrics`).

- `internal/metrics` leaf: registry, runtime gauges, invoke/async histograms, reload counter
- Config `metrics.enabled` / `metrics.listen`; CLI `--metrics-listen`
- Daemon HTTP lifecycle; Status `metrics_listen`; docs EN+RU; DESIGN §1.5/§5/§7 + AGENTS row

```bash
make lint && make intent-check && make docs-check
go test ./internal/metrics/... ./internal/config/... ./internal/dispatch/... ./internal/daemon/... ./internal/server/... ./cmd/... -race -count=1
go fix ./internal/metrics/... ./internal/config/... ./internal/dispatch/... ./internal/daemon/... ./internal/server/... ./cmd/...
```

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-08-31 | metrics | Prometheus scrape HTTP + Observer histograms + reload hook |
| 2026-08-29 | v0.0.5 | Trajectory counters + session import `--out`; CHANGELOG + docs version bump |

## Verify (repo green)

```bash
make lint && make intent-check && make docs-check && make test
```

## Blockers

(none)
