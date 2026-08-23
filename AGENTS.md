# agentd

Go **1.26.7** daemon that proxies coding-agent hooks via [agenthooks](https://github.com/speakeasy-api/agenthooks). Cobra CLI; gRPC over unix/npipe; Buf + protovalidate; testify.

Architecture: [DESIGN.md](./DESIGN.md) · Deep conventions: [CONVENTIONS.md](./CONVENTIONS.md) · Session: [PROGRESS.md](./PROGRESS.md)

## Commands

```bash
make lint                          # golangci-lint or go vet + buf lint
make test                          # go test ./... -race -count=1
make test-short                    # same with -short
make e2e                           # scripts/e2e-m*.sh (discovered; shared scripts/e2e-common.sh)
make build                         # go build -o agentd .
make start                         # build + agentd daemon start
make stop                          # agentd daemon stop
make generate                      # buf lint + buf generate
make intent-check                  # package comments present (internal/)
go test ./internal/daemon/... -race -count=1
go fix ./path/to/changed/...
```

When a milestone ships `scripts/e2e-mN.sh`, name it `e2e-mN.sh` under `scripts/` — `make e2e` discovers `scripts/e2e-m*.sh` automatically. Shared setup lives in `scripts/e2e-common.sh` (source only).
## Code style

No unused symbols. No “for later” APIs. No drive-by refactors. Comments only for non-obvious **why**.

- **Call directly** — use the owning `internal/` API from the call site. No passthrough helpers or duplicate `switch`/parse tables for the same concern.
- **One source of truth** — extend the package that owns the behavior; do not copy or re-export it elsewhere to “keep layers clean”.
- **Real abstractions welcome** — adapters, factories, and interfaces at package boundaries are fine when they encode a design decision (see [CONVENTIONS.md](./CONVENTIONS.md#do-not-duplicate-or-wrap)).

```go
// CORRECT — lowercase error; wrap with %w; sentinel for callers
return fmt.Errorf("acquire lock: %w", ErrAlreadyRunning)

// WRONG
return fmt.Errorf("Failed to Acquire Lock: %v", err)
```

Name domain strings and non-obvious numbers. Keep them next to the concern (never a kitchen-sink `consts.go`). `0` / `1` / `""` / `-1` are fine when the meaning is obvious.

```go
// CORRECT
const readyTimeout = 5 * time.Second
type FailMode string
const FailOpen FailMode = "fail_open"
if mode == FailOpen { ... }

// WRONG
time.After(5 * time.Second)
if mode == "fail_open" { ... }
```

Imports: std → third-party → `github.com/macrox-pro/agentd/...`.

### Files

One concern per file; `snake_case.go`; name the concern, not the package (`start.go`, not `daemon_start.go`). Never `util.go` / `helpers.go` / `types.go`.

```text
internal/daemon/start.go          # Start / StartOptions
internal/daemon/lock_unix.go      # //go:build unix
internal/daemon/lock_windows.go
internal/daemon/lock_other.go     # //go:build !unix && !windows
cmd/daemon.go                     # cobra: agentd daemon
store_test.go                     # tests store.go; package config_test
```

Platform: `{concern}.go` + `_unix.go` / `_windows.go` / `_other.go`. `_unix.go` always needs `//go:build unix`. Do not put portable code in `_{GOOS}.go`.

### Tests

- `package xxx_test` only — never `package foo` in `*_test.go`.
- testify `assert` / `require`; table when ≥2 similar cases; subtest = `tt.name`.
- Unit tests: fakes / bufconn / httptest — no real ports unless testing the network path.
- Details: [CONVENTIONS.md § Tests](./CONVENTIONS.md#tests).

## Intent before code

Before writing code, fill an **intent note** (PR body or session notes):

- **Problem** — what and why
- **Hot path** — `invoke_sync` | `config_reload` | `async_side` | `other` ([DESIGN §1.5](./DESIGN.md#15-hot-paths))
- **Invariants** — what must stay true
- **Corner cases** — scenarios to cover (test names, not prose in `_test.go`)
- **Out of scope** — what this change explicitly skips

Rules:

- One PR or agent session → one package or one hot path.
- Behavior or boundary change → update Tier-1 package comment and/or DESIGN §1.5.
- New corner case without a `TestXxx` / `tt.name` → do not merge.

### Plans (before implementation)

When preparing an implementation plan (Plan mode or session handoff), **segment TODOs** so nothing is missed:

- Split work into ordered segments (e.g. rules/constraints → preflight audit → code changes → tests → verify → PROGRESS handoff).
- **One actionable item per todo** — do not bundle unrelated steps (lint, code edit, docs update) into a single checkbox.
- Map every intent-note **corner case**, checklist row, and verify command to at least one todo.
- Call out **known exceptions** and **out-of-scope** items explicitly (separate todos or a dedicated segment) so they are not accidentally “fixed” mid-phase.
- Preserve execution order in the plan body (segment tables or numbered lists) matching the todo list.

A plan with fewer than ~5 coarse todos for a multi-file phase is usually under-segmented — expand until each step is independently checkable.

## Architecture

| Path | Owns | Must not |
|------|------|----------|
| `cmd/` | Cobra — flags, Short/Long/Example, **CLI input validation** in the subcommand file, `RunE` delegates to `internal/` | business logic, domain helpers, extra validation-only files |
| `internal/` | Domain logic, sentinel errors (`var Err…`) | import `cmd/`, Cobra/flag names in errors, CLI argument rules |
| `api/agentd/v1/` | Proto contracts | hand-written Go |
| `gen/` | Generated — never edit | manual edits |
| `internal/hookedge` | Provider codecs + wire I/O ([§1.5 invoke_sync](./DESIGN.md#15-hot-paths)) | policy, `Runner.Decide` |
| `internal/config` | Config merge/compile ([§1.5 config_reload](./DESIGN.md#15-hot-paths)) | dispatch, wire decode |
| `internal/dispatch` | Routing; Engine; queue; session lock ([§1.5 invoke_sync, async_side](./DESIGN.md#15-hot-paths)) | wire decode, YAML compile, Kind→impl switch |
| `internal/dispatch/targets` | Sync/Async target adapters; factories map CompiledTarget → invokers | route match, YAML compile, session lock |
| `internal/trajectory` | Session ledger append, persist, list/export ([§1.5 async_side](./DESIGN.md#15-hot-paths), §14) | wire decode, route match, config compile |
| `internal/guard` | secrets/shell/mcp/paths checks | routing, encode |
| `internal/daemon` | start/stop, lock, status, reload signal | dispatch, config compile |
| `internal/server` | Thin gRPC mapping | policy, dispatch logic |
| `internal/transport` | Unix socket / named pipe I/O | business logic |
| `internal/hookclient` | gRPC client to daemon | hook wire |
| `internal/install` | Provider hook install via agenthooks | daemon logic |
| `internal/provider` | Canonical coding-agent provider ids and enum mapping | wire I/O, ledger, importer logic |

- `cmd/`: flags + Cobra + **CLI input validation in the same file as the subcommand** (`session_import.go`, not a separate `session_import_validate.go`); domain logic in `internal/` ([CONVENTIONS.md § CLI](./CONVENTIONS.md#cli-cmd)).
- Hook CLI: decode/encode only. `Runner.Decide` runs in the daemon.
- Never log to stdout on the hook path. Preserve `Event.Raw` verbatim.
- Hook entrypoint: public `agentd hook run|notify|serve`; `agenthooks/install` writes `agentd agenthooks …` (hidden alias, same `cmd/hook.go` path). Document `hook`, not `agenthooks`.
- ConfigStore hot path: `store.Current()` only — no disk I/O. Runtime overlay + one debounced reload goroutine.
- Async dispatch must not block the sync hook response.
- **Target extensibility:** new sync/async target kinds are added under `internal/dispatch/targets` via factory; `Engine` must not grow a `switch` on target kind.
- **Guard extensibility:** built-in guards live in `internal/guard`; wiring into the agenthooks Runner stays in `targets` builtin (or a single registrar). Do not call `guard` from `Engine` directly.
- **Interfaces:** prefer narrow interfaces at package boundaries (`SyncInvoker`, `AsyncInvoker`, optional server-side ports). No kitchen-sink `interfaces.go`; define next to the consumer or implementor concern file.
- **Behavior-preserving refactors (R-phases):** default acceptance is table/golden parity on Decide outcomes for the same Snapshot + payload; document any intentional policy change in DESIGN.md first.
- New CLI command → update [DESIGN.md §6](./DESIGN.md#6-cli-reference) **and** [docs/en/cli.md](./docs/en/cli.md) + [docs/ru/cli.md](./docs/ru/cli.md).
- User-visible behavior / config / Status / install change → update matching pages under [docs/en/](./docs/en/) then mirror [docs/ru/](./docs/ru/) (see [docs/en/maintaining.md](./docs/en/maintaining.md)). Run `make docs-check`.

## Protobuf

- Buf v2 + protovalidate (`buf.validate`); plugins: `protoc-gen-go` + `protoc-gen-go-grpc` only (no grpc-gateway until needed).
- Package `agentd.v1`; Go: `github.com/macrox-pro/agentd/gen`. Additive in `v1`; breaking → `v2/`.
- RPCs: verb phrases; messages `{Rpc}Request`/`Response`; enums `*_UNSPECIFIED` zero.
- After proto edits: `make generate` then `go build ./...`.

## Boundaries

- Never edit `gen/` or hand-write `*.pb.go`.
- Never reimplement provider codecs outside `hookedge` + agenthooks.
- Never add a layer whose only job is forwarding to an existing function — call it directly (adapters/factories at real boundaries are fine; see [CONVENTIONS.md](./CONVENTIONS.md#do-not-duplicate-or-wrap)).
- **`internal/` must not import `cmd/`** or mention Cobra flags (`--provider`, `--path`, …) in errors or comments.
- **CLI argument validation** (required flags, flag combinations) lives in `cmd/` in the **same file** as the subcommand — not in `internal/`.
- **Sentinel errors** in `internal/`: `var ErrFoo = errors.New("…")` when the message is static; wrap with `fmt.Errorf("…: %w", err)` only when adding context. CLI maps sentinels to user-facing flag text in `RunE`.
- Never commit secrets / `.env`.
- Change only files required by the current todo.

## Session & PR

On stop / context limit: update [PROGRESS.md](./PROGRESS.md) (next todo + files touched).

PR: intent note + comprehension checklist ([template](./.github/pull_request_template.md)); `make lint` + `make intent-check` + `make test` on touched packages; `make e2e` when shipping a milestone e2e script (and wire it into Makefile `e2e`); `make generate` if `api/` changed; DESIGN.md CLI section **and** `docs/en`+`docs/ru` if commands or user-facing behavior changed (`make docs-check`).
