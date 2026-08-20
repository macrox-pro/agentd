# AGENTS.md — agentd contributor guide

Rules for human and AI contributors. **After Phase 1, this file is the source of truth** for coding conventions (over plan or chat history).

- Architecture: [DESIGN.md](./DESIGN.md)
- Progress handoff: [PROGRESS.md](./PROGRESS.md)
- Go idioms: [go-guidelines](https://github.com/mhmtszr/go-guidelines) skill (Go **1.26.7** detected in `go.mod`)

## Project overview

**agentd** is a user-level daemon that proxies, guards, and observes coding-agent hooks via [agenthooks](https://github.com/speakeasy-api/agenthooks). Hook CLI is thin; daemon owns dispatch, config, and gRPC.

## Go version

Go **1.26.7**. After changes:

```bash
golangci-lint run ./path/to/changed/...   # or go vet
go test ./path/to/changed/... -race
go fix ./path/to/changed/...              # Go 1.26+
```

## Package conventions

- `cmd/` — Cobra wiring only; full Short/Long/Example on every command
- `internal/` — all business logic
- `gen/` — generated protobuf; never edit
- Import order: std → third-party → `github.com/macrox-pro/agentd/...`

## Code quality

### No unused code

- Do not add functions, types, flags, or proto fields "for later" — every symbol must have a caller in the same change (skeleton phase p3 excepted: stubs may return `errNotImplemented` until wired).
- Keep `golangci-lint` / `staticcheck` / `unused` clean in packages you touch.
- Delete dead code; never leave commented-out blocks "just in case".

### Comments sparingly

- Prefer clear names over comments. Comment only non-obvious invariants, provider quirks, or **why** (not **what**).
- Do not restate the function name in a comment. No narration (`// loop over items`).
- Unexported symbols: no doc comment unless there is a quirk worth documenting.
- Exported API: one-line godoc when behavior is not obvious from the signature.

### Do not duplicate logic

- Hook wire I/O lives in `internal/hookedge` + agenthooks — do not reimplement provider codecs.
- Config merge/compile lives in `internal/config` only.
- Dispatch routing lives in `internal/dispatch` only.
- Extract shared test helpers once; use `t.Helper()`.
- Before adding a helper, search `internal/` for an existing one.

### Minimal scope

- Change only files required by the current todo. Do not drive-by refactor unrelated code.

## Naming

- MixedCaps; acronyms `ID`, `URL`, `HTTP`, `API`
- `context.Context` as first parameter when needed; never store context in structs

## Error handling

- Error strings lowercase
- Wrap with `%w` for context
- Sentinels: `ErrDaemonNotRunning`, etc.

## Testing

- Tests in **`package xxx_test` only** — never `package xxx` for unit tests
- Table-driven tests with `t.Run(tt.name, ...)` ([TableDrivenTests](https://go.dev/wiki/TableDrivenTests))
- `t.Helper()` in helpers; `t.Cleanup()` for teardown
- `go test ./... -race`; integration behind `//go:build integration`
- Prefer bufconn/httptest over real network in unit tests

## Protobuf & Buf (2026)

Canonical API contracts live in `api/`. Generated code in `gen/` — never edit by hand.

### Toolchain

- **Buf v2** (`buf.yaml`, `buf.gen.yaml`, `buf.lock` committed)
- **protovalidate** (`buf.build/bufbuild/protovalidate`) — `(buf.validate.field)` / `(buf.validate.message)` CEL rules; do **not** add legacy `protoc-gen-validate`
- **Plugins (v1 IPC only):** `protoc-gen-go` + `protoc-gen-go-grpc`
- **No grpc-gateway** until HTTP/JSON is required (YAGNI)

### Layout & versioning

- Directory: `api/agentd/v1/*.proto`
- Proto package: `agentd.v1`
- Go import via buf managed `go_package_prefix`: `github.com/macrox-pro/agentd/gen`
- Additive changes only within `v1`; breaking changes → `api/agentd/v2/`
- One service per file + `common.proto` for shared types

### Naming (buf STANDARD lint)

- Services: `PascalCase` + `Service` suffix
- RPCs: verb phrases (`Invoke`, `Health`, `PatchConfig`)
- Messages: `{RpcName}Request` / `{RpcName}Response`
- Enums: `UPPER_SNAKE`; zero = `*_UNSPECIFIED`
- Fields: `snake_case`
- Use `google.protobuf.Timestamp` and `google.protobuf.Duration`

### Validation & fields

- Required semantics via protovalidate, not proto2 `required`
- Reserve removed fields: `reserved 2, 3; reserved "old_field";`
- Prefer enums for `Provider`, `EventKind`, `DecisionKind`

### Buf workflow (every proto change)

```bash
buf lint api/
buf generate
buf breaking --against .git#branch=main   # when baseline exists
go build ./...
```

### gRPC in Go

- Implement services in `internal/server`; thin mapping only
- Register on unix/npipe listener (v1)

### Testing proto

- Round-trip in `*_test` with `proto.Marshal` / `proto.Equal`
- Table-driven validation failure cases

DESIGN.md §4 lists RPC catalog only — do not duplicate buf rules there.

## Cross-platform

- Platform I/O in `internal/transport` with build tags (`unix.go`, `windows.go`)

## agenthooks integration

- Never log to stdout on hook path
- Preserve `Event.Raw` verbatim
- `Runner.Decide` in daemon; codecs only in `hookedge`
- Install command must use `agentd hook run --provider=...`

## ConfigStore

- Hot path: `store.Current()` only — no disk I/O
- Daemon writes runtime overlay only
- Reload in one goroutine; debounced fsnotify

## Dispatch

- Hook CLI: decode/encode only
- Async never blocks sync response
- New targets: `internal/dispatch/targets/` + table-driven tests in `targets_test`

## CLI documentation

Every new subcommand:

1. Cobra Short, Long, Example (English)
2. Matching section in [DESIGN.md § CLI Reference](./DESIGN.md#6-cli-reference)

## Session handoff

On stop or context limit:

1. Update [PROGRESS.md](./PROGRESS.md) with next todo and files touched
2. If p2+: confirm AGENTS.md was read this session

## PR checklist

- [ ] AGENTS.md compliance (code quality, no duplication)
- [ ] `golangci-lint run` on changed packages (includes unused)
- [ ] `buf lint` if `api/` changed
- [ ] `go test -race` on changed packages
- [ ] `make generate` if proto changed
- [ ] DESIGN.md CLI section if commands changed
- [ ] PROGRESS.md updated
