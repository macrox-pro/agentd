# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub-post-v1). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **m10** | Last: m9-trajectory-p0 | Next: m10-a search + Claude import

## agents_md_ready

true

## Roadmap (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M7 | **done** | Daemon through approvals / RecordDecision / runtime persist |
| M8 / v1 | **done** | Overflow counters, conformance, docs freeze, release |
| **M9** | **done** | Trajectory P0 — L0 live ledger for all six providers + export |
| **M10** | planned | Trajectory P1 — search + Claude import; others stay L0 |
| **M11** | planned | Trajectory P2 — importers if format exists; policy replay all dialects |
| **M12 / v1.1** | planned | Trajectory P3 — Subscribe; contract + §14.6; **v1.1 release** |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## M9 checklist — Trajectory hub P0

**Constraint:** L0 live ledger for **all six** providers (DESIGN §14.6). Entrypoints: run / argv-payload / notify / serve.

### Phase A — Store + catalog

- [x] m9-a-catalog — draft event types in code (`hook/invoked`, `hook/decided`, …) aligned with DESIGN §14.3; include `provider` + `invocation_mode`
- [x] m9-a-store — `internal/trajectory` append-only in-memory seq + JSONL under `sessions/<provider>/…`
- [x] m9-a-test — unit tests: contig seq, immutable after append, truncate `max_event_bytes`
- [x] m9-a-checkpoint

### Phase B — Engine wiring (all providers)

- [x] m9-b-wire — enqueue trajectory record from Invoke path (after sync / async_only); never block wire
- [x] m9-b-providers — fixtures or e2e coverage: claude-code, cursor (argv), codex (run+notify), gemini, opencode (serve), kimi-code
- [x] m9-b-overflow — drop + Status field `trajectory_dropped_count`
- [x] m9-b-test — server/dispatch tests with trajectory enabled
- [x] m9-b-checkpoint

### Phase C — Config

- [x] m9-c-schema — `trajectory:` in config (enabled, include_raw, redact, max_event_bytes); default **off**
- [x] m9-c-compile — merge/validate; fingerprint includes trajectory knobs
- [x] m9-c-docs — DESIGN §7 snippet
- [x] m9-c-checkpoint

### Phase D — CLI + docs

- [x] m9-d-cli — `agentd session list|show|export` (`--provider` filter)
- [x] m9-d-matrix — docs en/ru: trajectory page linking §14.6 limits per agent
- [x] m9-d-design-cli — DESIGN §6 + `make docs-check`
- [x] m9-d-checkpoint

### Phase E — Close M9

- [x] m9-e-e2e — `scripts/e2e-m9.sh` (multi-provider smoke)
- [x] m9-e-verify — lint + test + e2e
- [x] m9-e-checkpoint

**M9 acceptance:** see [DESIGN.md §13 M9](./DESIGN.md#m9--trajectory-hub-p0-live-ledger).

## M10–M12 / v1.1 (outline only)

Checklists expand when M10 starts. Summary:

| Milestone | Focus |
|-----------|--------|
| M10 | `session search`; Claude JSONL import; `source=transcript`; other providers stay L0 + importer status `none`/`partial` |
| M11 | Importers only where format exists; `session replay --policy` **all six** wire dialects; `session fork` |
| M12 / v1.1 | Subscribe; versioned contract; README §14.6; tag **v1.1.0** |

Do **not**: invent thinking/tool results; Claude-only L0; agent-loop resume (DESIGN §12.8 / §14.6–§14.7).

## Session notes

- M9 shipped: `internal/trajectory/`, `trajectory:` config, `session` CLI, `trajectory_dropped_count`, `scripts/e2e-m9.sh`
- Next: M10 search + Claude import

## Verify (last green)

```bash
make lint
make intent-check
make docs-check
make test
go test -tags=integration ./internal/hookedge/ -race -count=1
make e2e   # includes scripts/e2e-m9.sh
```

## Blockers

(none)
