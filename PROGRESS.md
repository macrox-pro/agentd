# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**v0.0.6 release** — ready for tag.

- Prometheus metrics (`internal/metrics`, daemon HTTP, Observer histograms, reload counter)
- Trajectory stats: Cursor `stop` billing tokens, per-session delta aggregation, `--json` enum names
- Daemon unit tests isolated from production socket/state/log (`testEnv` in `start_test.go`)
- CHANGELOG + docs version bump to v0.0.6

```bash
make lint && make intent-check && make docs-check && make test
# tag (user): git tag -a v0.0.6 -m "v0.0.6" && git push origin v0.0.6
# release binaries: goreleaser (GitHub Actions on tag push)
```

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
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
