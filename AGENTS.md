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

`package xxx_test` only. testify `assert` / `require` / `mock` — no suite, no ginkgo. Table when ≥2 similar cases:

```go
func TestParsePolicy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{name: "empty", in: ""},
		{name: "invalid", in: "nope", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ParsePolicy(tt.in)
			if tt.wantErr {
				require.Error(t, err, "ParsePolicy(%q)", tt.in)
				return
			}
			require.NoError(t, err, "ParsePolicy(%q)", tt.in)
		})
	}
}
```

- Subtest = `tt.name`; fail msgs include input; no `tt := tt` (Go 1.22+).
- `require` for preconditions; `assert` for independent checks; `t.Helper()` + `t.Cleanup` in helpers.
- Unit tests: bufconn / httptest / fakes — no real ports unless testing the network path.
- Integration: `//go:build integration`.

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

## Architecture

| Path | Owns | Must not |
|------|------|----------|
| `cmd/` | Cobra only — Short/Long/Example on every command | business logic |
| `internal/` | All business logic | — |
| `api/agentd/v1/` | Proto contracts | hand-written Go |
| `gen/` | Generated — never edit | manual edits |
| `internal/hookedge` | Provider codecs + wire I/O ([§1.5 invoke_sync](./DESIGN.md#15-hot-paths)) | policy, `Runner.Decide` |
| `internal/config` | Config merge/compile ([§1.5 config_reload](./DESIGN.md#15-hot-paths)) | dispatch, wire decode |
| `internal/dispatch` | Routing; new targets → `targets/` ([§1.5 invoke_sync, async_side](./DESIGN.md#15-hot-paths)) | wire decode, YAML compile |
| `internal/trajectory` | Session ledger append, persist, list/export ([§1.5 async_side](./DESIGN.md#15-hot-paths), §14) | wire decode, route match, config compile |
| `internal/guard` | secrets/shell/mcp/paths checks | routing, encode |
| `internal/daemon` | start/stop, lock, status, reload signal | dispatch, config compile |
| `internal/server` | Thin gRPC mapping | policy, dispatch logic |
| `internal/transport` | Unix socket / named pipe I/O | business logic |
| `internal/hookclient` | gRPC client to daemon | hook wire |
| `internal/install` | Provider hook install via agenthooks | daemon logic |

- Hook CLI: decode/encode only. `Runner.Decide` runs in the daemon.
- Never log to stdout on the hook path. Preserve `Event.Raw` verbatim.
- Hook entrypoint: public `agentd hook run|notify|serve`; `agenthooks/install` writes `agentd agenthooks …` (hidden alias, same `cmd/hook.go` path). Document `hook`, not `agenthooks`.
- ConfigStore hot path: `store.Current()` only — no disk I/O. Runtime overlay + one debounced reload goroutine.
- Async dispatch must not block the sync hook response.
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
- Never commit secrets / `.env`.
- Change only files required by the current todo.

## Session & PR

On stop / context limit: update [PROGRESS.md](./PROGRESS.md) (next todo + files touched).

PR: intent note + comprehension checklist ([template](./.github/pull_request_template.md)); `make lint` + `make intent-check` + `make test` on touched packages; `make e2e` when shipping a milestone e2e script (and wire it into Makefile `e2e`); `make generate` if `api/` changed; DESIGN.md CLI section **and** `docs/en`+`docs/ru` if commands or user-facing behavior changed (`make docs-check`).
