# agentd — implementation progress

> Session handoff for agents. Roadmap to v1: [DESIGN.md §13](./DESIGN.md#13-milestones). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **m8** | Last: m8-g-checkpoint | Next: (tag v1.0.0 release when ready)

## agents_md_ready

true

## Roadmap (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M7 | **done** | Daemon through approvals / RecordDecision / runtime persist |
| **M8 / v1** | **done** | Overflow counters, conformance, docs freeze, release |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## M8 checklist

### Phase A — Async overflow Status

- [x] m8-a-proto — `async_dropped_count` on StatusResponse
- [x] m8-a-wire — Queue.Dropped → Status + CLI JSON
- [x] m8-a-test-docs — server/daemon tests + DESIGN §5/§6
- [x] m8-a-checkpoint

### Phase B — Timeout margin

- [x] m8-b-schema — route `sync_timeout`
- [x] m8-b-budget — SyncBudget (10% margin)
- [x] m8-b-engine — Engine applies deadline; grpc clamp
- [x] m8-b-docs
- [x] m8-b-checkpoint

### Phase C — Conformance

- [x] m8-c-conformance — agenthookstest fixtures (one per provider)
- [x] m8-c-checkpoint

### Phase D — Integration

- [x] m8-d-integration — `roundtrip_integration_test.go`
- [x] m8-d-checkpoint

### Phase E — Docs freeze

- [x] m8-e-readme
- [x] m8-e-design
- [x] m8-e-checkpoint

### Phase F — Release

- [x] m8-f-version — `internal/version`
- [x] m8-f-release — goreleaser + CI/release workflows + CHANGELOG
- [x] m8-f-checkpoint

### Phase G — Close M8

- [x] m8-g-e2e — `scripts/e2e-m8.sh`
- [x] m8-g-verify — lint + test + e2e
- [x] m8-g-checkpoint

**M8 acceptance:** met.

## Session notes

- AGENTS.md / CONVENTIONS.md read: yes (2026-08-20)
- M8 complete: Status drop counter; SyncBudget + route.sync_timeout; conformance; integration tag; docs freeze; goreleaser/CI; e2e-m8
- Key files: `api/agentd/v1/daemon.proto`, `internal/dispatch/timeout.go`, `internal/hookedge/{conformance_test,roundtrip_integration_test}.go`, `internal/version/version.go`, `scripts/e2e-m8.sh`, `.goreleaser.yaml`, `.github/workflows/`

## Verify (last green)

```bash
make lint
make test
go test -tags=integration ./internal/hookedge/ -race -count=1
make e2e   # includes scripts/e2e-m8.sh
```

## Blockers

(none)
