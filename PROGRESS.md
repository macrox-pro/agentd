# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**E2E coverage (M15/M16/M18/M20)** — plan complete.

### Segment checkpoints

| Segment | Status | Script / verify |
|---------|--------|-----------------|
| S0 — `e2e_expect_exit` helper | done | `scripts/e2e-common.sh` |
| S1 — M15 + M19 | done | `scripts/e2e-m15.sh` |
| S2 — M16 metrics | done | `scripts/e2e-m16.sh` |
| S3 — M18 TUI gate | done | `scripts/e2e-m18.sh` |
| S4 — M20 policy wire | done | `scripts/e2e-m20.sh` |
| S5 — docs handoff | done | DESIGN §13, CHANGELOG Unreleased, PROGRESS |
| S6 — full gate | done | lint, intent-check, docs-check, test, `make e2e`, `E2E_RETRIES=3 make e2e` |

### E2E rows covered

- **m15** — `stats_requires_daemon`, rollup/provider filter, statistics/trajectory gates, offline `session stats`, `cursor_two_stops_sum_tokens`
- **m16** — metrics off/enabled, runtime gauges, invoke histogram, `--metrics-listen` override, listener released
- **m18** — `AGENTD_NO_TUI` / `CI=true` setup gate, bare install validation, `--yes`/`--dry-run` conflict
- **m20** — sync failure `fail_closed`/`fail_open`, prompt block, `ask_fallback`, notify cwd → project fingerprint, serve cwd per frame

### Out of scope (this session)

- Product code changes in `cmd/` / `internal/` (bugs → separate PR)
- `docs/en` / `docs/ru` (no user-visible CLI changes)
- `Makefile` (e2e glob discovery already works)
- TUI wizard PTY e2e (gate only)

**Next:** open PR with intent note + comprehension checklist, or tag v0.0.9-beta.

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-09-02 | E2E M15–M20 | `e2e-m15`/`m16`/`m18`/`m20` + `e2e_expect_exit`; DESIGN §13 M20 row |
| 2026-09-02 | Policy/reliability | policy.fail + ask_fallback; unsupported removed; Cwd; projectsMu; docs-check + platform CI |
| 2026-09-02 | OpenCode docs research | 52 topic files under `research/opencode/` |
| 2026-09-01 | v0.0.8-beta | Cursor stats per-generation sum |

## Blockers

(none)
