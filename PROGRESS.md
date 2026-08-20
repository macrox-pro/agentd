# agentd — implementation progress

> Session handoff for agents. Roadmap to v1: [DESIGN.md §13](./DESIGN.md#13-milestones). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **m6** | Last: m5-g-checkpoint | Next: m6-a guards schema

## agents_md_ready

true

## Roadmap (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M5 | **done** | Daemon through config layers, ConfigService, config CLI, merged fingerprint |
| **M6** | planned | Guards: shell, mcp, paths |
| M7 | planned | Approvals / RecordDecision, runtime persist, temporary blocks |
| M8 / v1 | planned | Overflow counters, conformance, docs freeze, release |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## M5 checklist

### Phase A — Project layer

- [x] m5-a-resolve — nearest `.agentd.yaml` from cwd / project_root
- [x] m5-a-merge — `defaults ⊕ user ⊕ project`
- [x] m5-a-watch — lazy fsnotify on first project sighting
- [x] m5-a-test
- [x] m5-a-checkpoint

### Phase B — Runtime overlay

- [x] m5-b-path — `$XDG_STATE_HOME/agentd/runtime.yaml` (+ defaults)
- [x] m5-b-load — merge as highest layer on startup / reload
- [x] m5-b-watch — ignore self-writes (atomic rename)
- [x] m5-b-test
- [x] m5-b-checkpoint

### Phase C — Fingerprint

- [x] m5-c-canonical — `sha256(canonical_json(merged_config))`
- [x] m5-c-status — Status exposes new fingerprint + generation
- [x] m5-c-test
- [x] m5-c-checkpoint

### Phase D — ConfigService (Get / Patch)

- [x] m5-d-get — `GetConfig` (merged + per-layer)
- [x] m5-d-patch — `PatchConfig` → in-memory merge + snapshot swap
- [x] m5-d-register — register on daemon gRPC server
- [x] m5-d-test
- [x] m5-d-checkpoint

### Phase E — Config CLI

- [x] m5-e-validate — `agentd config validate` (offline compile)
- [x] m5-e-show — `agentd config show` (merged / layer)
- [x] m5-e-patch — `agentd config patch` via gRPC
- [x] m5-e-test
- [x] m5-e-checkpoint

### Phase F — Hot path cwd

- [x] m5-f-invoke — Invoke / routes honor project-aware snapshot when cwd set
- [x] m5-f-test
- [x] m5-f-checkpoint

### Phase G — Close M5

- [x] m5-g-e2e — `scripts/e2e-m5.sh`
- [x] m5-g-lint-test — `make lint` + `make test` on touched packages
- [x] m5-g-docs — DESIGN/README/PROGRESS sync
- [x] m5-g-checkpoint — mark M5 done in DESIGN §13

**M5 acceptance:** met. `RecordDecision` remains Unimplemented until M7.

## Later (do not start until M5 checkpoint)

### M6 — Guards

- Schema + compile for shell / mcp / paths
- `internal/guard` handlers + builtin attach by route list
- `scripts/e2e-m6.sh`

### M7 — Approvals

- Runtime `approvals` + `blocks.temporary`
- `RecordDecision` + TTL (project 24h, session end)
- Debounced runtime.yaml flush
- `scripts/e2e-m7.sh`

### M8 / v1

- Overflow drop counter on Status
- Provider timeout margin polish
- agenthookstest / integration build tag
- Docs freeze + GitHub release binaries
- `scripts/e2e-v1.sh` + v1 exit criteria in DESIGN §13

## Session notes

- AGENTS.md / CONVENTIONS.md read: yes (2026-08-20)
- M5 complete: project + runtime layers, Fingerprint(canonical JSON), ConfigService Get/Patch, config CLI, SnapshotFor on Invoke, e2e-m5
- Files: `internal/config/{project,fingerprint,paths_*,store,watch,compile,merge}.go`, `internal/server/config.go`, `cmd/config_*.go`, `scripts/e2e-m5.sh`
- **Convention:** `make e2e` runs all `scripts/e2e-mN.sh` — when adding `e2e-m6.sh` (etc.), append to Makefile `e2e` in the same PR (see AGENTS.md Commands)

## Verify (last green)

```bash
go test ./internal/config/... ./internal/server/... ./internal/hookedge/... ./internal/daemon/... ./cmd/... -race -count=1
make e2e   # includes scripts/e2e-m5.sh
```

## Blockers

(none)
