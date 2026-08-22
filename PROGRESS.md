# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub-post-v1). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **m10 done** | Last: m10-trajectory-p1 | Next: m11-a importers + policy replay

## agents_md_ready

true

## Roadmap (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M7 | **done** | Daemon through approvals / RecordDecision / runtime persist |
| M8 / v1 | **done** | Overflow counters, conformance, docs freeze, release |
| **M9** | **done** | Trajectory P0 — L0 live ledger for all six providers + export |
| **M10** | **done** | Trajectory P1 — search + Claude import; others stay L0 |
| **M11** | planned | Trajectory P2 — importers if format exists; policy replay all dialects |
| **M12 / v1.1** | planned | Trajectory P3 — Subscribe; contract + §14.6; **v1.1 release** |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## M10 checklist — Trajectory hub P1

**Constraint:** L0 live path unchanged for all six providers; Claude L2 import only.

### Phase A — Search

- [x] m10-a-events — transcript/message, transcript/thinking, source=transcript
- [x] m10-a-search — internal/trajectory/search.go (O(n) JSONL scan)
- [x] m10-a-cli-search — `agentd session search`
- [x] m10-a-checkpoint

### Phase B — Config

- [x] m10-b-config-schema — trajectory.import.claude-code
- [x] m10-b-config-compile — fingerprint includes import knobs
- [x] m10-b-config-test
- [x] m10-b-config-docs
- [x] m10-b-checkpoint

### Phase C — Claude import + merge

- [x] m10-c-import-pkg — internal/trajectory/importer/
- [x] m10-c-claude-map — agenthooks/transcript + thinking blocks
- [x] m10-c-checkpoint-dedup — import sidecar .import.json
- [x] m10-c-import-append — append-only seq merge
- [x] m10-c-cli-import — reject non-Claude with explicit none
- [x] m10-c-checkpoint

### Phase D — Daemon watcher

- [x] m10-d-watcher — fsnotify debounced import (async_side)
- [x] m10-d-watcher-wire — start.go + config reload
- [x] m10-d-checkpoint

### Phase E — L0 regression + list status

- [x] m10-e-list-status — importer_status in session list --json
- [x] m10-e-l0-regression — e2e-m9 unchanged
- [x] m10-e-checkpoint

### Phase F — Close M10

- [x] m10-f-design-cli — DESIGN §6/§7/§13
- [x] m10-f-docs — docs en/ru cli, trajectory, configuration
- [x] m10-f-e2e — scripts/e2e-m10.sh
- [x] m10-f-verify — lint + test + e2e
- [x] m10-f-checkpoint

**M10 acceptance:** see [DESIGN.md §13 M10](./DESIGN.md#m10--trajectory-p1-search--claude-import).

## M11–M12 / v1.1 (outline only)

| Milestone | Focus |
|-----------|--------|
| M11 | Importers where format exists; `session replay --policy` all six wire dialects; `session fork` |
| M12 / v1.1 | Subscribe; versioned contract; README §14.6; tag **v1.1.0** |

## Session notes

- M10 shipped: `session search|import`, Claude transcript import, importer_status, `trajectory.import`, `scripts/e2e-m10.sh`
- Daemon operational logging: `logging` config, `SetupLog`, default `$XDG_STATE_HOME/agentd/agentd.log`, `--log-level`/`--log-file`, `scripts/e2e-m10-logging.sh`
- Next: M11 multi-import + policy replay

## Verify (last green)

```bash
make lint
make intent-check
make docs-check
make test
go test -tags=integration ./internal/hookedge/ -race -count=1
make e2e   # includes scripts/e2e-m9.sh, e2e-m10.sh
```

## Blockers

(none)
