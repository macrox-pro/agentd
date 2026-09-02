# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**Windows daemon state paths** — done locally, unverified on Windows.

`platform-test (windows-latest)` was red since the job landed in `2ab9dfa`: every `transport.Listen` in `internal/daemon` tests got `Incorrect function`, because the tests pass a file path where Windows needs a named pipe. Under it sat a real bug — `NewPaths` derived pid/lock from `filepath.Dir(socket)`, which is `\\.\pipe` for the default Windows endpoint.

- `internal/transport/path_windows.go` — `pipeNamespace` const + `IsPipePath`
- `internal/daemon/paths.go` + `paths_unix.go` / `paths_windows.go` / `paths_other.go` — `stateDir` seam; pipe endpoints use `config.DefaultStateDir()`
- `internal/daemon/start_test.go` — `testSocket` returns a unique pipe on Windows
- Tests: `TestIsPipePath`, `TestPathsStateDirWindows` (both Windows-only)
- Docs: `docs/{en,ru}/configuration.md`, `docs/{en,ru}/operations.md`, `docs/en/troubleshooting.md`, DESIGN §5, CHANGELOG Unreleased

### Next todo

1. Push and watch `platform-test (windows-latest)` — `make lint` / `make test` cover darwin only; the Windows path has no local runner.
2. If daemon tests still fail on Windows, read the full job log: `internal/config` results were cut off in the reported tail.
3. Follow-up (separate session): `trajectory.DefaultSessionsDir()` has no Windows branch — the session ledger lands in `%USERPROFILE%\.local\state\agentd\sessions` instead of `%LOCALAPPDATA%\agentd`, and daemon tests do not isolate it.

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-09-02 | Windows state paths | pipe endpoints keep pid/lock in the state dir; daemon tests listen on a pipe |
| 2026-09-02 | v0.0.9-beta | released; CHANGELOG/README/DESIGN release metadata |
| 2026-09-02 | E2E M15–M20 | `e2e-m15`/`m16`/`m18`/`m20` + `e2e_expect_exit`; DESIGN §13 M20 row |
| 2026-09-02 | Policy/reliability | policy.fail + ask_fallback; unsupported removed; Cwd; projectsMu; docs-check + platform CI |
| 2026-09-01 | v0.0.8-beta | Cursor stats per-generation sum |

## Blockers

(none)
