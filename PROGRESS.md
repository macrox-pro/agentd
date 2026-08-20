# agentd — implementation progress

> Source of truth for multi-session / context-compressed agents.
> Full design: [DESIGN.md](./DESIGN.md) | Contributor rules: [AGENTS.md](./AGENTS.md)

## Current phase

Phase: m4 | Last: m3-complete | Next: grpc-forward (M4)

## agents_md_ready

true

## Session notes

- AGENTS.md read: yes (2026-08-20)
- CONVENTIONS.md read: yes (2026-08-20)
- M3 complete: full dispatch YAML, forward targets (log/file/http/exec), fsnotify reload, e2e-m3
- Queue capacity/workers still fixed at daemon start; timeout comes from snap at enqueue; capacity change needs restart

## M3 checklist

### Phase 0
- [x] m3-p0-agents

### Phase A — Config dispatch YAML
- [x] m3-a-schema
- [x] m3-a-target-kinds
- [x] m3-a-compile
- [x] m3-a-test
- [x] m3-a-checkpoint

### Phase B — Route match
- [x] m3-b-match
- [x] m3-b-test
- [x] m3-b-checkpoint

### Phase C — AsyncInvoker
- [x] m3-c-iface
- [x] m3-c-engine-wire
- [x] m3-c-sync-builtin-only
- [x] m3-c-checkpoint

### Phase D — Forward targets
- [x] m3-d-log
- [x] m3-d-file
- [x] m3-d-http
- [x] m3-d-exec
- [x] m3-d-test
- [x] m3-d-checkpoint

### Phase E — Engine
- [x] m3-e-engine
- [x] m3-e-test
- [x] m3-e-checkpoint

### Phase F — fsnotify
- [x] m3-f-watch
- [x] m3-f-daemon
- [x] m3-f-test
- [x] m3-f-checkpoint

### Phase G — CLI
- [x] m3-g-cli

### Phase H — Close
- [x] m3-h-e2e
- [x] m3-h-lint-test
- [x] m3-h-docs
- [x] m3-h-checkpoint

## Files touched (this session)

- internal/config/{schema,compile,merge,watch}.go + tests
- internal/dispatch/{engine,route}.go + tests
- internal/dispatch/targets/{async,log,file,http,exec,observe,factory}.go + tests
- internal/daemon/{start,reload_unix,reload_other}.go
- cmd/dispatch.go
- scripts/e2e-m3.sh
- DESIGN.md §6/§13
- PROGRESS.md
- go.mod (fsnotify direct)

## Verify (last green)

```bash
bash scripts/e2e-m3.sh
bash scripts/e2e-m2.sh
bash scripts/e2e-m1.sh
go test ./internal/config/... ./internal/dispatch/... ./internal/daemon/... -race -count=1
golangci-lint run ./internal/config/... ./internal/dispatch/... ./internal/daemon/... ./cmd/...
```

## Blockers

(none)

## Deferred (M4+)

gRPC forward, OpenCode serve, install, ConfigService, shell/mcp/paths guards, project layer / runtime overlay
