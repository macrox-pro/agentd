# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub-post-v1). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **R11 done** | Last: r11-comments-docs + audit follow-through | Next: **none (R-series complete; leftover = coverage only)**

> Milestones M0–M12: **done**. Post-release **R-series** refactor: **R1–R11 done**. No R12.

## agents_md_ready

true

## Refactoring goals (R-series)

| Goal | How we measure done |
|------|---------------------|
| **Code quality** | No duplicate provider/import switches; sentinels not stringly-typed errors; `make lint` + `staticcheck` clean on touched packages |
| **Architecture** | One owner per concern ([DESIGN §1](./DESIGN.md#1-architecture), [CONVENTIONS § Do not duplicate](./CONVENTIONS.md#do-not-duplicate-or-wrap)); purposeful adapters at boundaries only |
| **Test coverage** | Repo statements **≥ 70%** (baseline **57.3%**); no touched package regresses; gaps closed in priority order (see baseline) |
| **Table-driven tests** | ≥2 similar cases → table + `tt.name` ([CONVENTIONS § Table-driven](./CONVENTIONS.md#table-driven-tests-required-pattern)); no scenario prose in `_test.go` |
| **Corner cases** | Every new behavior row in intent note → named `TestXxx` / `tt.name`; no merge without test |
| **Duplication** | Litmus: removing helper leaves call sites importing owner package directly → delete helper |

### Coverage baseline (2026-08-23; R11 update)

| Package | Statements | Priority |
|---------|------------|----------|
| `cmd` | 60.7% | maintain — R4 |
| `internal/daemon` | 65.0% | maintain — R5 |
| `internal/hookclient` | 81.0% | maintain — R5 |
| `internal/dispatch/targets` | 59.7% | maintain — R6/R8 |
| `internal/trajectory/importer` | 71.8% | maintain — R3 |
| `internal/trajectory` | 68.2% | maintain — R10 |
| `internal/dispatch` | 74.0% | maintain — R7/R10 |
| `internal/hookedge` | 70.4% | maintain |
| `internal/config` | 75.0% | maintain |
| `internal/transport` | 77.5% | maintain |
| `internal/server` | 82.2% | maintain — R9/R11 |
| `internal/guard` | 86.0% | maintain — R8/R10 |
| `internal/provider` | 93.0% | maintain — R1 |

**Repo total (R11):** 61.5% (baseline 57.3%; R-series ≥70% goal **not** met — exception accepted; no filler tests).

Verify after each phase: `make lint` · `make intent-check` · `make test` · touched-package `-coverprofile`.

### Pitfalls (read before every R-phase)

- **`provider.Parse` on hot paths** — Invoke already carries proto enum; use `provider.FromProto` / `Lookup` where id was validated upstream; reserve `Parse` for CLI and strict entrypoints.
- **`SessionKey.Provider` is `provider.ID`** — `CanonicalProvider` stays `string` for paths/filters; JSONL `Event.Provider` stays `string`. Do not change `weakSessionID` hash inputs without a migration test.
- **Table-driven hub/subscribe tests** — shared `Hub` + channel state: parallelize subtests only when each row uses disjoint hubs ([CONVENTIONS § Parallelism](./CONVENTIONS.md#parallelism)).
- **`cmd` table tests** — reset Cobra flags between rows (`resetCommandFlags`); stale `Changed` bits cause flaky cross-row pollution. For `stringSlice`/`stringArray`, reset with `Set("")` not `Set(DefValue)` — nil defaults use DefValue `"[]"`, which parses as one element `"[]"`.
- **`daemon` lifecycle tests** — short unix socket paths (`os.MkdirTemp("", "agentd-daemon-")` + `s.sock`); never `net.Listen` before `Start` on the same path — `transport.Listen` removes the file and the daemon may start instead of failing. Prefer gRPC stubs for Status/Reload rows; reserve full `Start` for lock/PID/SIGHUP paths.
- **`hookclient` tests** — one prod file → one `client_test.go`; real unix socket for dial (not bufconn) when exercising `transport.Dial`.
- **`internal/` errors** — no `--provider` / flag names in sentinels; map in `cmd/RunE` only.
- **Importer registry** — must not pull `dispatch.InvokeInput` (PolicyInvoker still out of scope); registry returns `(ImportResult, error)` only.
- **Server HookService ports** — `Invoker` / `SnapshotSource` in `server/invoke.go`; typed-nil-safe `New` wiring (`opts.Engine` / `opts.Store` assigned only when non-nil).
- **No drive-by refactors** — one R-phase touches one boundary; unrelated cleanups → separate PR.

---

## R-series roadmap (summary)

| Phase | Status | One-liner |
|-------|--------|-----------|
| **R1** | **done** | `internal/provider` — single source for ids, proto/agenthooks mapping |
| **R2** | **done** | Trajectory domain — typed ids at boundaries, errors, grpc event mapping |
| **R3** | **done** | Importer registry — `registry.go` + `ImportSession` facade |
| **R4** | **done** | `cmd/` coverage — table-driven session + config CLI |
| **R5** | **done** | `daemon` + `hookclient` lifecycle tests |
| **R6** | **done** | `dispatch/targets` — SyncInvoker + factory; Engine Kind switch removed |
| **R7** | **done** | `dispatch` — extract `first_conclusive` aggregation |
| **R8** | **done** | `guard` + `targets/builtin` — Checker registry (`AttachCheckers`) |
| **R9** | **done** | `server` — narrow HookService ports (`Invoker`, `SnapshotSource`) |
| **R10** | **done** | Test hygiene — table-driven migration where structure matches |
| **R11** | **done** | Tier-1 package comments + DESIGN §1.5; SessionKey/Invoker/wrapper audit follow-through |

Done phases: acceptance and files live in git history / PRs. Do not re-expand intent notes here.

---

## R10 intent note — Test hygiene (table-driven migration)

**Problem:** Several packages use copy-pasted `TestXxx` functions where assertion skeleton is identical ([CONVENTIONS](./CONVENTIONS.md#when-not-to-use-a-table) allows exceptions — apply selectively).

**Hot path:** n/a (test-only).

**Invariants:** do **not** table-ify multi-step daemon/e2e flows; keep integration tag on slow tests.

**Merge candidates (keep separate tests when setup diverges):**

| File | Current | Target |
|------|---------|--------|
| `trajectory/hub_test.go` | 10 top-level tests | `TestHubDeliverTable` (`deliver`, `slow drop`, `unregister`, `no subscribers`) |
| `trajectory/fork_test.go` | 2 tests | `TestForkSessionTable` (`ok`, `duplicate rejected`) |
| `server/invoke_test.go` | 3 trajectory variants | `TestInvokeTrajectoryTable` (claude, cursor argv, …) |
| `guard/approve_test.go` | 3 separate | extend existing `TestTemporaryBlockDenies` pattern |

**Also parked from old R6:** residual `TestDecodeProvider` / `TestGRPCTargetForward` if gaps remain after R1.

**Out of scope:** testify/suite; rewriting e2e scripts; 100% coverage chase.

### R10 checklist

- [x] r10-intent
- [x] r10-hub-table — `TestHubDeliverTable` (`deliver`, `ignorable preserved`); kept separate: slow drop, unregister, no subscribers, multiple/close/concurrent/enqueue/schema
- [x] r10-fork-table — `TestForkSessionTable` (`ok`, `duplicate rejected`)
- [x] r10-invoke-table — `TestInvokeTrajectoryTable` (`claude`, `cursor argv`)
- [x] r10-audit — greps clean on touched packages; subscribe merge skipped (optional)
- [x] r10-verify — `make lint` · `make intent-check` · `make test`; total 61.5% (no padding)

**R10 hub keep-separate (CONVENTIONS skeleton gate):** `TestHubSlowConsumerDrop`, `TestHubUnregister`, `TestHubPublishNoSubscribers`, `TestHubMultipleSubscribers`, `TestHubCloseEndsSubscribers`, `TestHubConcurrentRegisterPublish`, `TestSubscribeDoesNotBlockEnqueue`, `TestSchemaVersionOnAppend`.

**R10 parked R6:** `TestDecodeTypedTable` merged; GRPC already tabled in `targets_test.go` — no change.

**R10 optional skipped:** `session_test.go` subscribe filter merge (non-blocking).

---

## R11 intent note — Package comments & docs sync

**Problem:** New packages (`provider`, trajectory splits, SyncInvoker / guard registry / server ports) need Tier-1 comments; user-facing behavior unchanged but boundaries shifted.

**Hot path:** n/a.

**Invariants:** package comment template per [CONVENTIONS](./CONVENTIONS.md#comments); DESIGN §1.5 hot path tags accurate.

**Corner cases:** `make intent-check` passes; no stale `hookedge/provider` references in docs.

**Out of scope:** New CLI commands; proto changes; README feature marketing.

### R11 checklist

- [x] r11-intent
- [x] r11-provider-comment — verify `internal/provider/provider.go` header
- [x] r11-trajectory-comments — Entry + godoc verified; no Package on concern files
- [x] r11-design-15 — DESIGN §1.5 Package tags table with `internal/provider` as `other`
- [x] r11-docs-check — no user-doc edits (no stale paths); `make docs-check` ok
- [x] r11-verify — `make lint` · `make intent-check` · `make test`; total 61.5% (no padding)

**Same-session follow-through (no R12) — done:**

- [x] `SessionKey.Provider` → `provider.ID`; `weak id hash unchanged`
- [x] `dispatch.Invoker` for replay; `server.Invoker` alias; importer clean
- [x] delete `providerFromProto` / `canonicalImportProvider`; `provider.Cursor` compares
- [x] export `config.GuardSecrets`/`Shell`/`MCP`/`Paths`; guard registry keys use them
- [x] nil `SnapshotSource` → neutral; `snap_nil_neutral`

**R-series complete.** Leftover: coverage gap only (61.5% vs ≥70% goal — exception; no pad).

---

## Milestones (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M7 | **done** | Daemon through approvals / RecordDecision / runtime persist |
| M8 / v0.0.1 | **done** | Overflow counters, conformance, docs freeze, release |
| M9–M12 / v0.0.2 | **done** | Trajectory P0–P3 — ledger, import, replay/fork, Subscribe |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## Session notes

- M0–M12 / v0.0.2 shipped (see DESIGN §13)
- **R1–R11 done** — R-series complete; leftover = coverage only (61.5%)
- No R12

## Verify (last green)

```bash
make lint
make intent-check
make docs-check
go test ./internal/trajectory/... ./internal/dispatch/... ./internal/server/... ./internal/guard/... ./internal/config/... -race -count=1
make test
go test ./... -coverprofile=/tmp/agentd_r11.out -count=1
go tool cover -func=/tmp/agentd_r11.out | tail -1  # total 61.5%
```

**R11 files touched:** `internal/provider/provider.go`, `DESIGN.md` §1.5, `PROGRESS.md`, `internal/trajectory/{session_key,paths,store,fork,recorder,import_append,replay,trajectory}.go` + tests, `internal/dispatch/{engine,decode}.go`, `internal/server/invoke.go` + `invoke_mapping_test.go`, `internal/config/{guards,trajectory}.go`, `internal/guard/registry.go`

## Blockers

(none)
