# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub-post-v1). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **R9 done** | Last: r9-server-ports | Next: **R10 test hygiene**

> Milestones M0–M12: **done**. Post-release work is the **R-series** refactor below — one phase = one PR or agent session, one package or one hot path ([AGENTS.md](./AGENTS.md) intent rules).

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

### Coverage baseline (2026-08-23; R9 update)

| Package | Statements | Priority |
|---------|------------|----------|
| `cmd` | 60.7% | maintain — R4 |
| `internal/daemon` | 65.0% | maintain — R5 |
| `internal/hookclient` | 81.0% | maintain — R5 |
| `internal/dispatch/targets` | 59.7% | maintain — R6/R8 |
| `internal/trajectory/importer` | 71.8% | maintain — R3 |
| `internal/trajectory` | 68.3% | P2 |
| `internal/dispatch` | 73.8% | maintain — R7 |
| `internal/hookedge` | 70.4% | maintain |
| `internal/config` | 74.8% | maintain |
| `internal/transport` | 77.5% | maintain |
| `internal/server` | 81.7% | maintain — R9 |
| `internal/guard` | 85.8% | maintain — R8 |
| `internal/provider` | 93.0% | maintain — R1 |

Verify after each phase: `make lint` · `make intent-check` · `make test` · touched-package `-coverprofile`.

### Pitfalls (read before every R-phase)

- **`provider.Parse` on hot paths** — Invoke already carries proto enum; use `provider.FromProto` / `Lookup` where id was validated upstream; reserve `Parse` for CLI and strict entrypoints.
- **`SessionKey.Provider` is still `string`** — migrate call sites to `provider.ID` gradually; do not change `weakSessionID` hash inputs without a migration test.
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
| **R10** | pending | Test hygiene — table-driven migration where structure matches |
| **R11** | pending | Tier-1 package comments audit + docs sync |

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

- [ ] r10-intent
- [ ] r10-hub-table — hub deliver/drop/unregister merged where safe
- [ ] r10-fork-table
- [ ] r10-invoke-table
- [ ] r10-audit — grep for `if tt.` blocks inside single `Test` without `t.Run`
- [ ] r10-verify — full `make test`; total coverage ≥ 70%

---

## R11 intent note — Package comments & docs sync

**Problem:** New packages (`provider`, trajectory splits, SyncInvoker / guard registry / server ports) need Tier-1 comments; user-facing behavior unchanged but boundaries shifted.

**Hot path:** n/a.

**Invariants:** package comment template per [CONVENTIONS](./CONVENTIONS.md#comments); DESIGN §1.5 hot path tags accurate.

**Corner cases:** `make intent-check` passes; no stale `hookedge/provider` references in docs.

**Out of scope:** New CLI commands; proto changes; README feature marketing.

### R11 checklist

- [ ] r11-intent
- [ ] r11-provider-comment — verify `internal/provider/provider.go` header
- [ ] r11-trajectory-comments — `hub.go`, `replay_config.go`, `session_event_grpc.go`
- [ ] r11-design-15 — DESIGN §1.5 row for `internal/provider` if missing
- [ ] r11-docs-check — `make docs-check` (only if user-visible text changed)
- [ ] r11-verify — `make intent-check` + full verify block

**R-series done when:** R1–R11 checklists complete, total coverage ≥ 70%, no duplicate provider switches outside owner tables.

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
- **R1–R9 done** — see summary table above
- Next agent session: open **R10** — test hygiene (table-driven migration)

## Verify (last green)

```bash
make lint
make intent-check
go test ./internal/server/... -race -count=1
go test ./internal/server/... -coverprofile=/tmp/server_r9.out -count=1
go tool cover -func=/tmp/server_r9.out | tail -1   # server 81.7%
make test
```

**R9 files touched:** `internal/server/{invoke,invoke_mapping_test,server}.go`, `CONVENTIONS.md`, `DESIGN.md` §2, `PROGRESS.md`

## Blockers

(none)
