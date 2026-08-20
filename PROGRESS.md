# agentd — implementation progress

> Source of truth for multi-session / context-compressed agents.
> Full design: [DESIGN.md](./DESIGN.md) | Contributor rules: [AGENTS.md](./AGENTS.md)

## Current phase

Phase: m5 | Last: m4-complete | Next: ConfigService (Deferred)

## agents_md_ready

true

## Session notes

- AGENTS.md read: yes (2026-08-20)
- CONVENTIONS.md read: yes (2026-08-20) — root `CONVENTIONS.md`
- M4 complete: gRPC forward (sync+async), OpenCode serve + notify, agenthooks sentinel, install wrapper, Windows SID pipe path, e2e-m4

## M4 checklist

### Phase 0
- [x] m4-p0-agents

### Phase A — Config grpc schema
- [x] m4-a-schema
- [x] m4-a-compile
- [x] m4-a-test
- [x] m4-a-checkpoint

### Phase B — Transport endpoint
- [x] m4-b-endpoint
- [x] m4-b-test
- [x] m4-b-checkpoint

### Phase C — gRPC forward target
- [x] m4-c-grpc-target
- [x] m4-c-factory
- [x] m4-c-test
- [x] m4-c-checkpoint

### Phase D — Session mutex
- [x] m4-d-session
- [x] m4-d-test
- [x] m4-d-checkpoint

### Phase E — OpenCode serve / notify / sentinel
- [x] m4-e-serve
- [x] m4-e-notify
- [x] m4-e-cli
- [x] m4-e-test
- [x] m4-e-checkpoint

### Phase F — Install
- [x] m4-f-install-pkg
- [x] m4-f-cli
- [x] m4-f-test
- [x] m4-f-checkpoint

### Phase G — Windows npipe
- [x] m4-g-sid-path
- [x] m4-g-test
- [x] m4-g-docs
- [x] m4-g-checkpoint

### Phase H — CLI docs
- [x] m4-h-cli-docs

### Phase I — Close
- [x] m4-i-e2e
- [x] m4-i-lint-test
- [x] m4-i-progress
- [x] m4-i-checkpoint

## Files touched (this session)

- internal/config/{schema,compile}.go + store_test.go
- internal/transport/{endpoint,path_windows,listen_windows}.go + tests
- internal/dispatch/{engine,session}.go + targets/{grpc,factory}.go + tests
- internal/server/server.go
- internal/hookedge/{serve,notify}.go + tests
- internal/install/install.go + tests
- cmd/{hook,agenthooks,install}.go
- scripts/e2e-m4.sh
- DESIGN.md §1/§5/§6/§9/§13
- PROGRESS.md

## Verify (last green)

```bash
bash scripts/e2e-m4.sh
bash scripts/e2e-m3.sh
bash scripts/e2e-m2.sh
bash scripts/e2e-m1.sh
go test ./internal/config/... ./internal/transport/... ./internal/dispatch/... ./internal/hookedge/... ./internal/install/... ./internal/server/... ./internal/daemon/... -race -count=1
golangci-lint run ./internal/config/... ./internal/transport/... ./internal/dispatch/... ./internal/hookedge/... ./internal/install/... ./internal/server/... ./internal/daemon/... ./cmd/...
```

## Blockers

(none)

## Deferred (M5+)

ConfigService, shell/mcp/paths guards, project layer / runtime overlay
