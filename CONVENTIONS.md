# Coding conventions (encyclopedia)

Deep reference for contributors. Agents: prefer [AGENTS.md](./AGENTS.md); open this file when naming files, writing tests, or changing proto. Architecture: [DESIGN.md](./DESIGN.md).

## Code quality

### No unused code

- Do not add functions, types, flags, or proto fields "for later" — every symbol must have a caller in the same change (skeleton stubs may return `errNotImplemented` until wired).
- Keep `golangci-lint` / `staticcheck` / `unused` clean in packages you touch.
- Delete dead code; never leave commented-out blocks "just in case".

### Comments

- Prefer clear names over comments. Comment only non-obvious invariants, provider quirks, or **why** (not **what**).
- Do not restate the function name. No narration (`// loop over items`).
- Unexported symbols: no doc comment unless there is a quirk worth documenting.
- Exported API: one-line godoc when behavior is not obvious from the signature.

### Do not duplicate logic

- Hook wire I/O: `internal/hookedge` + agenthooks — do not reimplement provider codecs.
- Config merge/compile: `internal/config` only.
- Dispatch routing: `internal/dispatch` only.
- Extract shared test helpers once; use `t.Helper()`. Before adding a helper, search `internal/`.

### Minimal scope

Change only files required by the current todo. Do not drive-by refactor unrelated code.

### No magic literals

Name domain strings and non-obvious numbers/durations. Put `const` next to the concern that owns them — not in a catch-all `consts.go` (see [Files](#files)).

```go
// CORRECT — typed domain string + named timeout beside Start
const readyTimeout = 5 * time.Second
type FailMode string
const FailOpen FailMode = "fail_open"

// WRONG — bare wire/config values and unexplained numbers
if mode == "fail_open" { ... }
time.After(5 * time.Second)
```

Leave obvious literals alone: `0`, `1`, `-1`, `""`, `nil`, loop bounds, and test table inputs that appear once. Prefer a typed string (`FailMode`) when the value is part of config/wire/protocol vocabulary.

## Naming

### Identifiers

- MixedCaps; acronyms `ID`, `URL`, `HTTP`, `API`.
- `context.Context` as first parameter when needed; never store context in structs.

### Files

A filename is a locator: from the path alone, know what lives in the file. Do not dump unrelated types into a catch-all file.

Form follows Go toolchain rules ([build constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)) and source-file convention (lowercase, underscores — [golang/go#36060](https://github.com/golang/go/issues/36060)). Identifiers stay MixedCaps; **filenames are not identifiers** ([Google Go style](https://google.github.io/styleguide/go/decisions.html#underscores)). Package directories: short, all lowercase, no underscores (`hookclient`, not `hook_client`) — [package names](https://go.dev/blog/package-names).

**Form**

- ASCII lowercase letters, digits, and underscores only: `store.go`, `lock_unix.go`
- Multi-word stems: `snake_case.go`. Do not use `camelCase.go`, `PascalCase.go`, `kebab-case.go`, or concatenated blobs (`lockunix.go`)
- Prefer one short English word when it uniquely names the concern (`start.go`, `paths.go`, `client.go`)
- Never start a compiled source filename with `.` or `_` — `go build` ignores those files
- No extra dots in the stem (`foo.bar.go`); the only dot is the `.go` / `.proto` extension

**Name the concern, not the package or the layer**

| Path | Contains |
|------|----------|
| `internal/config/store.go` | `Store` and merge/load |
| `internal/daemon/start.go` | `Start` / `StartOptions` |
| `internal/daemon/lock_unix.go` | unix `AcquireLock` / `ReleaseLock` |
| `cmd/daemon.go` | Cobra `agentd daemon` |

- Do not prefix with the package name (stutter): `start.go`, not `daemon_start.go`
- One primary concern per file; unexported helpers for that concern stay in the same file
- Do not use kitchen-sink names: `util.go`, `helpers.go`, `common.go`, `misc.go`, `types.go`, `interfaces.go`, `consts.go`
- `cmd/`: one file per subcommand, stem = command name (`hook.go` → `agentd hook`)
- `{package}.go`: holds the package comment (immediately before `package`); may also contain package-level code for that package (e.g. `transport.go` in package `transport`). Do not put the package comment in a random concern file.
- `errors.go`: package sentinel errors only (`ErrAlreadyRunning`, …)
- `main.go`: process entry only — no business logic

**Reserved suffixes** — use only when the special behavior is intended:

| Suffix | Meaning |
|--------|---------|
| `_test.go` | compiled only by `go test` |
| `_{GOOS}.go` / `_{GOARCH}.go` | automatic OS/arch constraint (`_windows.go`, `_amd64.go`) |
| `_unix.go` | **not** a GOOS; always add `//go:build unix` |
| `_other.go` | unsupported-platform fallback; always add a matching `//go:build` (e.g. `!unix && !windows`) |

Platform split: `{concern}.go` (shared API/types) + `{concern}_unix.go` + `{concern}_windows.go` + `{concern}_other.go`. Do not put portable code in a `_{GOOS}.go` file — the suffix constrains the build even without a `//go:build` line.

**Tests**

- `{file}_test.go` beside `{file}.go` (`store_test.go` tests `store.go`)
- `{package}_test.go` only when the test spans several files in the package
- Platform tests: `{file}_unix_test.go` (or `{package}_unix_test.go`) plus `//go:build unix`
- Do not put tests for `foo.go` in `bar_test.go`

**Do not**

| Wrong | Right |
|-------|--------|
| `ConfigStore.go` / `configStore.go` | `store.go` |
| `http-client.go` | `client.go` in package `hookclient` |
| `internal/daemon/daemon_start.go` | `start.go` |
| `utils.go` | `paths.go`, `lock_unix.go` |
| `foo_linux.go` for portable code | `foo.go` |
| `listener.go` that also dials and encodes | split files or name the real concern |

Non-Go: scripts may use kebab-case (`scripts/e2e-m1.sh`). Root docs keep existing names (`AGENTS.md`, `DESIGN.md`). Generated protobuf stays in `gen/` (`*.pb.go`, `*_grpc.pb.go`) — never hand-name or edit those files.

## Error handling

- Error strings lowercase
- Wrap with `%w` for context
- Sentinels: `ErrDaemonNotRunning`, etc.

## Testing

Reference: [TableDrivenTests](https://go.dev/wiki/TableDrivenTests). Assertions and mocks: [testify](https://github.com/stretchr/testify) (`assert`, `require`, `mock`). No ginkgo or other frameworks. Enable [testifylint](https://github.com/Antonboom/testifylint) via `golangci-lint` when configured.

### Package & layout

- Tests in **`package xxx_test` only** — never `package xxx` for unit tests (black-box; import `github.com/macrox-pro/agentd/internal/xxx`).
- File: `{file}_test.go` next to `{file}.go`; helpers in the same `_test` package or a shared `internal/.../testutil` if reused across packages.
- Integration / slow tests: `//go:build integration` at top of file; run with `go test -tags=integration ./...`.

### Table-driven tests (required pattern)

Use a table when the same assertion logic applies to **two or more** input/expected pairs. Do not copy-paste nearly identical test functions.

**Slice of structs** (default — deterministic order):

```go
import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    Policy
		wantErr bool
	}{
		{name: "empty", in: "", want: Policy{Fail: FailOpen}},
		{name: "fail closed", in: "fail_closed", want: Policy{Fail: FailClosed}},
		{name: "invalid", in: "nope", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParsePolicy(tt.in)
			if tt.wantErr {
				require.Error(t, err, "ParsePolicy(%q)", tt.in)
				return
			}
			require.NoError(t, err, "ParsePolicy(%q)", tt.in)
			assert.Equal(t, tt.want, got, "ParsePolicy(%q)", tt.in)
		})
	}
}
```

Rules:

- Every row **must** have a `name string` field; subtest name is `tt.name`, not the loop index.
- Loop body **must** be `for _, tt := range tests { t.Run(tt.name, func(t *testing.T) { ... }) }` — never test all rows inside one `t.Run` or one function without subtests.
- Fail messages **must** name the input: `require.NoError(t, err, "ParsePolicy(%q)", tt.in)` — not bare `require.NoError(t, err)`.
- **`require.*`** when later assertions depend on success — it calls `t.FailNow()`.
- **`assert.*`** for independent checks in one subtest.
- Go **1.22+** (this repo): do **not** add `tt := tt` before `t.Run`.

**Map of structs** — only when iteration order must not affect results. Map key is the subtest name:

```go
tests := map[string]struct {
	in   string
	want string
}{
	"empty": {in: "", want: ""},
	"hello": {in: "hello", want: "HELLO"},
}
for name, tt := range tests {
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, tt.want, Upper(tt.in), "Upper(%q)", tt.in)
	})
}
```

### When *not* to use a table

- Single scenario with multi-step setup (daemon start/stop, lock + PID lifecycle) — one `TestXxx` is fine.
- Subtests differ in **structure** (different APIs, different cleanup), not just inputs — split into separate top-level tests or helpers, not one wide table with `setup func(t *testing.T)`.

### Parallelism

- Call `t.Parallel()` at the start of `TestXxx` and inside each `t.Run` when the case has no shared mutable state.
- Do **not** call `t.Parallel()` on tests that bind ports, write global temp dirs, or race on the same file without isolation.
- Subtests that share one `t.TempDir()` from the parent: parallelize only if each row uses disjoint paths under that dir.

### Helpers & teardown

- Shared setup: extract `func helper(t *testing.T) ...` and call **`t.Helper()`** as the first line.
- Register cleanup with **`t.Cleanup(func() { ... })`** — not bare `defer` in helpers that return before the test finishes.
- Prefer **`t.TempDir()`** over manual `os.MkdirTemp` + defer remove.

### Testify

| Package | Use |
|---------|-----|
| `github.com/stretchr/testify/assert` | Value checks; multiple per subtest |
| `github.com/stretchr/testify/require` | Preconditions; stops subtest on failure |
| `github.com/stretchr/testify/mock` | `mock.Mock` embed + `On`/`Called`/`AssertExpectations` |

- Always pass **`t`** as the first argument.
- Sentinels: `require.ErrorIs(t, err, ErrFoo)`, `require.ErrorAs(t, err, &target)`.
- Proto / deep structs: `assert.Equal` or `proto.Equal`; `assert.Empty`, `assert.NotNil`, `assert.True`/`False` as needed.
- Many asserts in one subtest: optional `a := assert.New(t)` / `r := require.New(t)`.
- **`mock`**: embed `mock.Mock`, record with `m.On(...).Return(...)`, finish with `m.AssertExpectations(t)`. Use `mock.Anything` only when input is dynamic.
- **Do not use** `testify/suite` — no parallel subtests ([testify#934](https://github.com/stretchr/testify/issues/934)).
- **Do not use** deprecated `testify/http`.

### Assertions & I/O

- Prefer testify; raw `t.Errorf` / `t.Fatalf` only when testify has no clearer helper.
- Unit tests: **`bufconn`**, **`httptest`**, in-memory fakes — no real network or ephemeral TCP ports unless the code under test is the network path.
- Time: inject clocks or bounded `context.WithTimeout`; avoid `time.Sleep` except when polling an external process (e.g. daemon PID file).

### Anti-patterns

| Wrong | Right |
|-------|-------|
| `package config` in `*_test.go` | `package config_test` |
| One `TestFoo` with `if a { ... }; if b { ... }` blocks | Table + `t.Run(tt.name, ...)` |
| `for i, tt := range tests { ... }` without `t.Run` | Subtest per row |
| `require.NoError(t, err)` with no context | `require.NoError(t, err, "Op(%q)", tt.in)` |
| `assert` where next line dereferences result | `require` for precondition, then `assert` |
| `testify/suite` + `t.Parallel()` | Plain `TestXxx` + table + `t.Run` |
| Table with one row | Single test function or add more cases |
| `tt := tt` inside loop (Go 1.22+) | Remove — shadows incorrectly on 1.26 |

After changes: `go test ./path/to/pkg/... -race -count=1`.

## Protobuf & Buf

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

Or: `make generate`.

### gRPC in Go

- Implement services in `internal/server`; thin mapping only
- Register on unix/npipe listener (v1)

### Testing proto

- Round-trip in `*_test` with `proto.Marshal` / `proto.Equal`
- Table-driven validation failure cases

DESIGN.md §4 lists RPC catalog only — do not duplicate buf rules there.

## Cross-platform

Platform I/O: `{concern}_unix.go` / `{concern}_windows.go` / `{concern}_other.go` (see [Files](#files)).

## Domain invariants

### agenthooks

- Never log to stdout on hook path
- Preserve `Event.Raw` verbatim
- `Runner.Decide` in daemon; codecs only in `hookedge`
- Install command must use `agentd hook run --provider=...`

### ConfigStore

- Hot path: `store.Current()` only — no disk I/O
- Daemon writes runtime overlay only
- Reload in one goroutine; debounced fsnotify

### Dispatch

- Hook CLI: decode/encode only
- Async never blocks sync response
- New targets: `internal/dispatch/targets/` + table-driven tests in `targets_test`

### CLI documentation

Every new subcommand:

1. Cobra Short, Long, Example (English)
2. Matching section in [DESIGN.md § CLI Reference](./DESIGN.md#6-cli-reference)

## Session handoff & PR

On stop or context limit:

1. Update [PROGRESS.md](./PROGRESS.md) with next todo and files touched
2. Confirm [AGENTS.md](./AGENTS.md) was read this session

PR checklist:

- [ ] AGENTS.md / this file compliance (code quality, no duplication)
- [ ] `golangci-lint run` on changed packages (includes unused)
- [ ] `buf lint` if `api/` changed
- [ ] `go test -race` on changed packages
- [ ] `make generate` if proto changed
- [ ] DESIGN.md CLI section if commands changed
- [ ] PROGRESS.md updated
