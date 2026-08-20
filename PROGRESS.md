# agentd — implementation progress

> Source of truth for multi-session / context-compressed agents.
> Full design: [DESIGN.md](./DESIGN.md) | Contributor rules: [AGENTS.md](./AGENTS.md)

## Current phase

Phase: m5 | Last: packages-conventions-refactor | Next: ConfigService (Deferred)

## agents_md_ready

true

## Session notes

- AGENTS.md read: yes (2026-08-20)
- CONVENTIONS.md read: yes (2026-08-20) — root `CONVENTIONS.md`
- packages conventions refactor: install/hookedge/hookclient/server/guard — concern files, named consts, dead API removed, test layout
- M4 complete: gRPC forward (sync+async), OpenCode serve + notify, agenthooks sentinel, install wrapper, Windows SID pipe path, e2e-m4
- cmd/ refactor (AGENTS/CONVENTIONS): one file per subcommand; Cobra-thin; WriteStatus/Reload/DefaultUserPath in internal/; hook↔agenthooks builders shared; notimpl.go removed
- config/ refactor (AGENTS/CONVENTIONS): split schema/merge by concern; FormatRoutes folded into cmd/dispatch_routes.go; Store.reloadMu; removed KindDefault.Blocking + Snapshot.RawYAML; fingerprint remains sha256(raw user YAML) until M5

## Files touched (packages conventions refactor)

- internal/install: install.go→run.go (+ run_test.go); identity/timeout consts; provider normalize; table TestRun
- internal/hookedge: split options/provider/payload/decision/encode; hookedge.go package doc; notify_test.go; unified fromProto
- internal/hookclient: unexport daemon/hook; grpcDialTarget; client_test.go
- internal/server: split server.go / daemon.go / invoke.go; drop Options.Log; daemon_test.go + invoke_test.go
- internal/daemon/start.go: remove Log from server.Options
- internal/guard: rule-id consts; RuleIDs(); attach_test.go; align test vs config.DefaultSecretsRules
- PROGRESS.md

## Files touched (config conventions refactor)

- deleted internal/config/{schema,format_routes}.go (+ format_routes_test)
- added internal/config/{config,file,policy,async,guards,mode,route}.go + compile_test/mode_test
- updated internal/config/{merge,compile,store,defaults}.go + store_test
- cmd/dispatch_routes.go (+ dispatch_routes_test.go); FormatRoutes unexported helper in cmd
- DESIGN.md §3 fingerprint note; PROGRESS.md

## Files touched (cmd conventions refactor)

- cmd/{root,daemon,daemon_*,hook,hook_*,agenthooks,agenthooks_*,install,config,config_*,dispatch,dispatch_routes}.go; deleted notimpl.go
- internal/config/{paths}.go + tests
- internal/daemon/{status_write,reload}.go + status_write_test.go
- DESIGN.md §6 CLI tree (agenthooks/)
- PROGRESS.md

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
go test ./internal/install/... ./internal/hookedge/... ./internal/hookclient/... ./internal/server/... ./internal/guard/... ./internal/daemon/... ./internal/dispatch/... -race -count=1
golangci-lint run ./internal/install/... ./internal/hookedge/... ./internal/hookclient/... ./internal/server/... ./internal/guard/... ./internal/daemon/... ./cmd/...
```

## Blockers

(none)

## Deferred (M5+)

ConfigService, shell/mcp/paths guards, project layer / runtime overlay
