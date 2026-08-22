# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub-post-v1). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **m9** | Last: plan-trajectory-hub | Next: m9-a event catalog + `internal/trajectory` store

## agents_md_ready

true

## Roadmap (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M7 | **done** | Daemon through approvals / RecordDecision / runtime persist |
| **M8 / v1** | **done** | Overflow counters, conformance, docs freeze, release |
| **M9** | **planned** | Trajectory P0 — L0 for **all six** agents + export ([§14.6](./DESIGN.md#146-provider-support-matrix-all-supported-agents)) |
| **M10** | planned | Trajectory P1 — search + Claude import; others stay L0 |
| **M11** | planned | Trajectory P2 — importers if format exists; policy replay all dialects |
| **M12 / v1.1** | planned | Trajectory P3 — Subscribe; contract + §14.6; **v1.1 release** |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## M9 checklist — Trajectory hub P0

**Constraint:** L0 live ledger for **all six** providers (DESIGN §14.6). Entrypoints: run / argv-payload / notify / serve.

### Phase A — Store + catalog

- [ ] m9-a-catalog — draft event types in code (`hook/invoked`, `hook/decided`, …) aligned with DESIGN §14.3; include `provider` + `invocation_mode`
- [ ] m9-a-store — `internal/trajectory` append-only in-memory seq + JSONL under `sessions/<provider>/…`
- [ ] m9-a-test — unit tests: contig seq, immutable after append, truncate `max_event_bytes`
- [ ] m9-a-checkpoint

### Phase B — Engine wiring (all providers)

- [ ] m9-b-wire — enqueue trajectory record from Invoke path (after sync / async_only); never block wire
- [ ] m9-b-providers — fixtures or e2e coverage: claude-code, cursor (argv), codex (run+notify), gemini, opencode (serve), kimi-code
- [ ] m9-b-overflow — drop + Status field (reuse or `trajectory_dropped_count`)
- [ ] m9-b-test — server/dispatch tests with trajectory enabled
- [ ] m9-b-checkpoint

### Phase C — Config

- [ ] m9-c-schema — `trajectory:` in config (enabled, include_raw, redact, max_event_bytes); default **off**
- [ ] m9-c-compile — merge/validate; fingerprint includes trajectory knobs
- [ ] m9-c-docs — DESIGN §7 snippet if needed; user docs later with CLI
- [ ] m9-c-checkpoint

### Phase D — CLI + docs

- [ ] m9-d-cli — `agentd session list|show|export` (`--provider` filter)
- [ ] m9-d-matrix — docs en/ru: trajectory page or section linking §14.6 limits per agent
- [ ] m9-d-design-cli — DESIGN §6 + `make docs-check`
- [ ] m9-d-checkpoint

### Phase E — Close M9

- [ ] m9-e-e2e — `scripts/e2e-m9.sh` (multi-provider smoke)
- [ ] m9-e-verify — lint + test + e2e
- [ ] m9-e-checkpoint

**M9 acceptance:** see [DESIGN.md §13 M9](./DESIGN.md#m9--trajectory-hub-p0-live-ledger).

## M10–M12 / v1.1 (outline only)

Checklists expand when M9 ships. Summary:

| Milestone | Focus |
|-----------|--------|
| M10 | `session search`; Claude JSONL import; `source=transcript`; other providers stay L0 + importer status `none`/`partial` |
| M11 | Importers only where format exists; `session replay --policy` **all six** wire dialects; `session fork` |
| M12 / v1.1 | Subscribe; versioned contract; README §14.6; tag **v1.1.0** |

Do **not**: invent thinking/tool results; Claude-only L0; agent-loop resume (DESIGN §12.8 / §14.6–§14.7).

## Session notes

- Mental model: intent-first comprehension (Tier-1 package comments, DESIGN §1.5 hot paths, PR gate)
- AGENTS.md / CONVENTIONS.md read: yes (2026-08-21)
- Trajectory hub planned: universal **L0** for all supported agents; L2/L3 per-provider matrix in DESIGN §14.6
- DESIGN: §11 transcripts → §14; §12 Q6–Q8; §13 M9–M12; §14 (+ full provider limits)
- Next implementation todo: **m9-a-catalog** / **m9-a-store**
- Key files (plan): `DESIGN.md`, `PROGRESS.md`; future: `internal/trajectory/`, `cmd/session_*.go`, `api/agentd/v1/session.proto` (M12 / v1.1)

## Verify (last green)

```bash
make lint
make intent-check
make docs-check
make test
go test -tags=integration ./internal/hookedge/ -race -count=1
make e2e   # includes scripts/e2e-m8.sh; e2e-m9 when added
```

## Blockers

(none)
