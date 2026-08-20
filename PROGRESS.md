# agentd — implementation progress

> Session handoff for agents. Roadmap to v1: [DESIGN.md §13](./DESIGN.md#13-milestones). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **m5** | Last: roadmap-to-v1 | Next: m5-a project layer

## agents_md_ready

true

## Roadmap (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M4 | **done** | Daemon, dispatch, targets, secrets, install, OpenCode, Windows npipe |
| **M5** | **in progress** | Project + runtime layers, ConfigService, config CLI, merged fingerprint |
| M6 | planned | Guards: shell, mcp, paths |
| M7 | planned | Approvals / RecordDecision, runtime persist, temporary blocks |
| M8 / v1 | planned | Overflow counters, conformance, docs freeze, release |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## M5 checklist

### Phase A — Project layer

- [ ] m5-a-resolve — nearest `.agentd.yaml` from cwd / project_root
- [ ] m5-a-merge — `defaults ⊕ user ⊕ project`
- [ ] m5-a-watch — lazy fsnotify on first project sighting
- [ ] m5-a-test
- [ ] m5-a-checkpoint

### Phase B — Runtime overlay

- [ ] m5-b-path — `$XDG_STATE_HOME/agentd/runtime.yaml` (+ defaults)
- [ ] m5-b-load — merge as highest layer on startup / reload
- [ ] m5-b-watch — ignore self-writes (atomic rename)
- [ ] m5-b-test
- [ ] m5-b-checkpoint

### Phase C — Fingerprint

- [ ] m5-c-canonical — `sha256(canonical_json(merged_config))`
- [ ] m5-c-status — Status exposes new fingerprint + generation
- [ ] m5-c-test
- [ ] m5-c-checkpoint

### Phase D — ConfigService (Get / Patch)

- [ ] m5-d-get — `GetConfig` (merged + per-layer)
- [ ] m5-d-patch — `PatchConfig` → in-memory merge + snapshot swap
- [ ] m5-d-register — register on daemon gRPC server
- [ ] m5-d-test
- [ ] m5-d-checkpoint

### Phase E — Config CLI

- [ ] m5-e-validate — `agentd config validate` (offline compile)
- [ ] m5-e-show — `agentd config show` (merged / layer)
- [ ] m5-e-patch — `agentd config patch` via gRPC
- [ ] m5-e-test
- [ ] m5-e-checkpoint

### Phase F — Hot path cwd

- [ ] m5-f-invoke — Invoke / routes honor project-aware snapshot when cwd set
- [ ] m5-f-test
- [ ] m5-f-checkpoint

### Phase G — Close M5

- [ ] m5-g-e2e — `scripts/e2e-m5.sh`
- [ ] m5-g-lint-test — `make lint` + `make test` on touched packages
- [ ] m5-g-docs — DESIGN/README/PROGRESS sync
- [ ] m5-g-checkpoint — mark M5 done in DESIGN §13

**M5 acceptance:** see DESIGN §13 M5. `RecordDecision` may remain stub until M7.

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
- M4 complete (grpc forward, OpenCode serve/notify, install, Windows SID pipe, e2e-m4)
- Post-M4: cmd/ + config/ + packages conventions refactors (concern files; config CLI still stubs)
- ConfigStore today: defaults ⊕ user only; fingerprint = sha256(raw user YAML) until M5-C
- Roadmap M5–M8 / v1 written into DESIGN §13 (2026-08-20)

## Verify (last green)

```bash
go test ./internal/install/... ./internal/hookedge/... ./internal/hookclient/... ./internal/server/... ./internal/guard/... ./internal/daemon/... ./internal/dispatch/... -race -count=1
golangci-lint run ./internal/install/... ./internal/hookedge/... ./internal/hookclient/... ./internal/server/... ./internal/guard/... ./internal/daemon/... ./cmd/...
```

## Blockers

(none)
