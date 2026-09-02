# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**v0.0.9 release prep** — changelog + metadata updated locally; tag `v0.0.9` not pushed yet.

### Next todo

1. Commit release metadata, tag `v0.0.9`, push — `release.yml` runs goreleaser + GitHub Release.
2. Watch CI after push: `platform-test (windows-latest)` validates the pipe/state-dir fixes (no local Windows runner).
3. Follow-up (separate session): `trajectory.DefaultSessionsDir()` has no Windows branch — the session ledger lands in `%USERPROFILE%\.local\state\agentd\sessions` instead of `%LOCALAPPDATA%\agentd`, and daemon tests do not isolate it.

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-09-02 | v0.0.9 prep | CHANGELOG/README/DESIGN/docs for stable v0.0.9 (Windows fixes + beta scope) |
| 2026-09-02 | Windows state paths | pipe endpoints keep pid/lock in the state dir; daemon tests listen on a pipe |
| 2026-09-02 | v0.0.9-beta | released; CHANGELOG/README/DESIGN release metadata |
| 2026-09-02 | E2E M15–M20 | `e2e-m15`/`m16`/`m18`/`m20` + `e2e_expect_exit`; DESIGN §13 M20 row |
| 2026-09-02 | Policy/reliability | policy.fail + ask_fallback; unsupported removed; Cwd; projectsMu; docs-check + platform CI |
| 2026-09-01 | v0.0.8-beta | Cursor stats per-generation sum |

## Blockers

(none)
