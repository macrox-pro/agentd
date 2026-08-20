# agentd — implementation progress

> Source of truth for multi-session / context-compressed agents.
> Full design: [DESIGN.md](./DESIGN.md) | Contributor rules: [AGENTS.md](./AGENTS.md)

## Current phase

Phase: m3 | Last: m2-complete | Next: forward-targets (M3)

## agents_md_ready

true

## Session notes

- AGENTS.md read: yes (2026-08-20)
- CONVENTIONS.md read: yes (2026-08-20)
- M2 complete: Dispatch Engine, async queue, secrets guard, hookedge decision encode, e2e-m2

## M2 checklist

### Phase 0
- [x] m2-p0-agents

### Phase A — Config compile
- [x] m2-a-config-schema
- [x] m2-a-config-defaults
- [x] m2-a-config-compile
- [x] m2-a-config-test
- [x] m2-a-checkpoint

### Phase B — Secrets guard
- [x] m2-b-secrets-scan
- [x] m2-b-secrets-attach
- [x] m2-b-secrets-test
- [x] m2-b-checkpoint

### Phase C — Builtin target
- [x] m2-c-decision-map
- [x] m2-c-builtin-sync
- [x] m2-c-builtin-async
- [x] m2-c-builtin-test
- [x] m2-c-checkpoint

### Phase D — Async queue
- [x] m2-d-queue
- [x] m2-d-queue-overflow
- [x] m2-d-queue-test
- [x] m2-d-checkpoint

### Phase E — Dispatch Engine
- [x] m2-e-decode
- [x] m2-e-route-match
- [x] m2-e-modes
- [x] m2-e-engine-api
- [x] m2-e-engine-test
- [x] m2-e-checkpoint

### Phase F — Server + daemon
- [x] m2-f-server-invoke
- [x] m2-f-status-metrics
- [x] m2-f-daemon-wire
- [x] m2-f-server-test
- [x] m2-f-checkpoint

### Phase G — hookedge encode
- [x] m2-g-hookedge-encode
- [x] m2-g-hookedge-test
- [x] m2-g-checkpoint

### Phase H — dispatch routes CLI
- [x] m2-h-dispatch-routes

### Phase I — Close
- [x] m2-i-e2e
- [x] m2-i-lint-test
- [x] m2-i-docs
- [x] m2-i-checkpoint

## Files touched (this session)

- internal/config/{schema,defaults,merge,compile,store}.go + store_test.go
- internal/guard/{secrets,attach}.go + secrets_test.go
- internal/dispatch/{engine,queue,decode,route,decision}.go + tests
- internal/dispatch/targets/builtin.go + builtin_test.go
- internal/server/server.go + server_test.go
- internal/daemon/{start,status}.go
- internal/hookedge/run.go + run_test.go
- cmd/{dispatch,daemon}.go
- scripts/e2e-m2.sh
- DESIGN.md §13
- PROGRESS.md

## Verify (last green)

```bash
bash scripts/e2e-m2.sh
bash scripts/e2e-m1.sh
go test ./internal/config/... ./internal/guard/... ./internal/dispatch/... ./internal/server/... ./internal/hookedge/... ./internal/daemon/... -race -count=1
golangci-lint run ./internal/config/... ./internal/guard/... ./internal/dispatch/... ./internal/server/... ./internal/hookedge/... ./internal/daemon/... ./cmd/...
```

## Blockers

(none)

## Deferred (M3+)

Full dispatch YAML, fsnotify, forward targets (exec/http/log/file), gRPC forward, OpenCode serve, install, ConfigService
