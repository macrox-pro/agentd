# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**Daemon test isolation** — done.

- `internal/daemon/start_test.go`: `testEnv` (XDG_STATE_HOME / XDG_RUNTIME_DIR / LOCALAPPDATA), `skipUnlessDefaultSocketIsolated`, `TestTestEnv` table test
- `testSocket` calls `testEnv` for state/log/runtime isolation; socket dir stays short `MkdirTemp` (macOS unix path length)
- `default socket not running` subtests in stop/reload/status use isolated env + skip off darwin/linux
- `listen error` subtest uses `testEnv` for runtime overlay isolation

```bash
make lint && make test
go test ./internal/daemon/... -race -count=1
# acceptance: with `agentd daemon start` running, tests must not stop it or append version=test to ~/.local/state/agentd/agentd.log
```

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-08-31 | daemon tests | Isolate unit tests from production socket, agentd.log, runtime.yaml |
| 2026-08-31 | metrics | Prometheus scrape HTTP + Observer histograms + reload hook |
| 2026-08-29 | v0.0.5 | Trajectory counters + session import `--out`; CHANGELOG + docs version bump |

## Verify (repo green)

```bash
make lint && make intent-check && make docs-check && make test
```

## Blockers

(none)
