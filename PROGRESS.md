# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**Full hook kind coverage** — config vocabulary, catch-all routing, stats billing gate, Tier A install.

### Next todo

1. Commit hook-kind coverage (tests + docs + `e2e-m21.sh` green).
2. Tag `v0.0.9` still pending from prior session if not already pushed.
3. Follow-up: `trajectory.DefaultSessionsDir()` Windows branch.

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-09-03 | Hook kinds | 16-kind defaults, catch-all, skip lock on async_only, billing only on `agent.stop`, Tier A install |
| 2026-09-02 | v0.0.9 prep | CHANGELOG/README/DESIGN/docs for stable v0.0.9 (Windows fixes + beta scope) |
| 2026-09-02 | Windows state paths | pipe endpoints keep pid/lock in the state dir; daemon tests listen on a pipe |
| 2026-09-02 | v0.0.9-beta | released; CHANGELOG/README/DESIGN release metadata |
| 2026-09-02 | E2E M15–M20 | `e2e-m15`/`m16`/`m18`/`m20` + `e2e_expect_exit`; DESIGN §13 M20 row |
| 2026-09-02 | Policy/reliability | policy.fail + ask_fallback; unsupported removed; Cwd; projectsMu; docs-check + platform CI |
| 2026-09-01 | v0.0.8-beta | Cursor stats per-generation sum |

## Blockers

(none)
