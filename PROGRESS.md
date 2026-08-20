# agentd — implementation progress

> Session handoff for agents. Roadmap to v1: [DESIGN.md §13](./DESIGN.md#13-milestones). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **m8** | Last: m7-g-checkpoint | Next: m8-a async overflow Status counter

## agents_md_ready

true

## Roadmap (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M7 | **done** | Daemon through approvals / RecordDecision / runtime persist |
| **M8 / v1** | planned | Overflow counters, conformance, docs freeze, release |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## M7 checklist

### Phase A — Schema + compile

- [x] m7-a-schema — YAML `approvals` + `blocks.temporary`
- [x] m7-a-merge — upsert by fingerprint / tool+pattern
- [x] m7-a-compile — Snapshot + drop expired + ApprovalFingerprint
- [x] m7-a-test
- [x] m7-a-checkpoint

### Phase B — RecordDecision

- [x] m7-b-store — project 24h / session scope
- [x] m7-b-grpc — ConfigService + hookclient
- [x] m7-b-test
- [x] m7-b-checkpoint

### Phase C — Hot path skip Ask

- [x] m7-c-skip-ask — secrets + shell consult approvals
- [x] m7-c-emit-fp — `approval_fingerprint=` in Ask system_message
- [x] m7-c-test
- [x] m7-c-checkpoint

### Phase D — Temporary blocks

- [x] m7-d-blocks — Deny attach before guards
- [x] m7-d-test
- [x] m7-d-checkpoint

### Phase E — Persist

- [x] m7-e-persist — 500ms debounce atomic flush + IgnoreSelfWrite
- [x] m7-e-wire — PatchRuntime / RecordDecision / shutdown flush
- [x] m7-e-test
- [x] m7-e-checkpoint

### Phase F — Operator CLI

- [x] m7-f-cli — `agentd config record-decision` + DESIGN §6
- [x] m7-f-checkpoint

### Phase G — Close M7

- [x] m7-g-e2e — `scripts/e2e-m7.sh`
- [x] m7-g-lint-test
- [x] m7-g-docs
- [x] m7-g-checkpoint

**M7 acceptance:** met.

## Later (do not start until M7 checkpoint)

### M8 / v1

- Overflow drop counter on Status
- Provider timeout margin polish
- agenthookstest / integration build tag
- Docs freeze + GitHub release binaries
- `scripts/e2e-v1.sh` + v1 exit criteria in DESIGN §13

## Session notes

- AGENTS.md / CONVENTIONS.md read: yes (2026-08-20)
- M7 complete: approvals/blocks schema, RecordDecision, skip-Ask, temp blocks, runtime.yaml persist, `config record-decision`, e2e-m7
- Files: `internal/config/{approvals,blocks,persist,file,merge,compile,store}.go`, `internal/guard/{approve,block,attach,shell}.go`, `cmd/config_record_decision.go`, `scripts/e2e-m7.sh`
- Fingerprint format: `sha256:<kind>/<hex>` (kind embedded for RecordDecision routing)

## Verify (last green)

```bash
make lint
go test ./internal/config/... ./internal/guard/... ./internal/dispatch/... ./internal/server/... ./cmd/... -race -count=1
make e2e   # includes scripts/e2e-m7.sh
```

## Blockers

(none)
