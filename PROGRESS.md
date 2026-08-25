# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **policy.offline hook edge done** | Last: DialReady + Serve offline cache | Next: **none**

### policy.offline hook edge (2026-08-25)

**Problem:** Installed hooks block the agent when agentd is down because hook CLI always exits 1; policy.offline is parsed but unused.

**Hot path:** invoke_sync (disk I/O only on DialReady/Invoke failure; live Invoke stays Store.Current-only in daemon).

**Invariants:**
- Never write debug to stdout on the hook path.
- No Runner.Decide / guards / route compile on hookedge.
- Default Offline is fail_open (defaults.go).
- Offline mode = defaults ⊕ user ⊕ project(cwd) ⊕ runtime via OfflineFor.
- Hook entrypoints use `hookclient.DialReady` (lazy gRPC + Health).
- fail_open -> exit 0 + encode Neutral (Run/Serve) or exit 0 (Notify); stderr still prints daemon not running.
- fail_closed OR OfflineFor error -> stderr daemon not running, exit 1.
- Serve caches offline mode after the first mid-stream Invoke failure (stderr once).

**Corner cases (`tt.name`):**
- config OfflineFor: `default fail_open` | `user fail_closed` | `project overrides` | `runtime overrides` | `missing user file` | `invalid yaml`
- run: `daemon down default fail_open` | `daemon down fail_closed` | `daemon down project fail_closed` | `daemon down invalid config`
- notify: `daemon down default fail_open` | `daemon down fail_closed`
- serve: `daemon down default fail_open` | `daemon down fail_closed` | `invoke fail caches offline once`
- hookclient Dial: `dial ready missing socket` | `dial ready healthy`

**Out of scope:** offline state-dir cache; auto-start daemon; change install FailClosed; proto/gen; change policy.fail default; new e2e script; OfflineFor(ctx) (LoadWith ignores ctx).

**Done:**
- `internal/config`: `OfflineFor`, default `Offline: FailOpen`, `offline_test.go`
- `internal/hookclient`: `DialReady` (+ tests)
- `internal/hookedge`: `resolveOffline`, `runOffline` / `notifyOffline` / `serveOffline`; Serve mid-Invoke offline cache
- `cmd/hook.go` ConfigPath + Long texts
- AGENTS.md / DESIGN §1.5 + §6; docs EN/RU + README

**Files:** `internal/config/{defaults,offline,offline_test,config,compile_test}.go`, `internal/hookclient/{client,client_test}.go`, `internal/hookedge/{hookedge,encode,options,offline,run,notify,serve,*_test}.go`, `cmd/hook.go`, `cmd/hook_run.go`, `AGENTS.md`, `DESIGN.md`, `docs/en|ru/{configuration,troubleshooting,cli,getting-started,providers}.md`, `README.md`, `PROGRESS.md`

---

Phase: **CLI version** | Last: `agentd version` + daemon status running-process version | Next: **none**

### CLI version (2026-08-23)

**Problem:** Users need this binary's version without a daemon (`agentd version`). `daemon status` must keep reporting the **running process** version from Status RPC (can differ from this CLI after an upgrade).

**Hot path:** other

**Invariants:**
- `internal/version.Version` is the only build string for this binary
- `agentd version` does not dial the daemon
- `daemon status` `version` comes from the Status RPC, not `version.Version`
- JSON key stays `version` (proto field)

**Corner cases (`tt.name`):**
- version: `default dev`, `release semver`
- daemon status: `running human daemon version`, `running json daemon version`, `running uses daemon version not cli`

**Out of scope:** `--json` on `version`; root `--version`; proto changes; git commit

**Done:** `cmd/version.go`; CLI tables in `cmd/version_test.go` + `cmd/daemon_status_test.go`; DESIGN §6; docs EN/RU cli + operations + installation; CONTRIBUTING

---

Phase: **install UX done** | Last: install defaults + stdout report | Next: **none**

> Milestones M0–M12: **done**. Post-release **R-series** refactor: **R1–R11 done**. No R12.
> Post-R follow-through: **F1–F4 done**. Coverage goal ≥70% still not met — exception (no filler).

### Install UX (2026-08-23)

**Problem:** `agentd install --scope=user` without `--dir` wrote into CWD; success was silent.

**Done:**
- `internal/install`: `ResolveDir`, `WriteReport`, `Result`; user scope defaults to agent home; codex project → `./.codex`; plugin/opencode-user require `--dir`
- `cmd/install`: report on stdout; `SilenceErrors` removed; sentinel → CLI errors
- Docs: DESIGN §6, `docs/en`+`docs/ru` cli + providers
- Tests: `dir_test`, `run_test`, `cmd/install_test`

**Files:** `internal/install/{errors,dir,report,run}.go`, `cmd/install.go`, docs EN/RU

---

## F1 intent note — Server black-box tests + comment hygiene

**Problem:** `invoke_mapping_test.go` / `invoke_mode_test.go` use `package server` (CONVENTIONS requires `package xxx_test`); server package comment missing Invariants; stale `v1.1` in trajectory comments (product is v0.0.2).

**Hot path:** `other`.

**Invariants:**
- black-box `*_test.go` only (`package server_test`)
- `hookService` stays unexported; seam is `NewHookService`
- do not export `normalizeInvocationMode`
- `Options.Store` / `Options.Engine` stay concrete types
- `SchemaVersion` value unchanged (`1`)

**Corner cases (`tt.name`):**
- mapping: `snap_nil_neutral`, `engine_nil_neutral`, `engine_error_neutral`, `engine_success_maps_decision`, `snap_generation_in_response`, `async_count_forwarded`
- mode: `cursor unspecified to argv`, `cursor argv preserved`, `claude unspecified stays unspecified`, `codex unspecified to stdin`

**Out of scope:** F2 decision codec; F3 Subscribe table; F4 coverage; proto/`gen/`; Options interface change.

### F1 checklist

- [x] f1-intent
- [x] f1-factory (`NewHookService` + wire `New`)
- [x] f1-mapping-blackbox
- [x] f1-mode-table (+ delete white-box)
- [x] f1-invariants + v0.0.2 comments
- [x] f1-verify

**F1 done.**

---

## F2 intent note — internal/decision codec

**Problem:** Duplicate proto↔agenthooks Decision Kind switches in `dispatch.DecisionToProto`, `hookedge.fromProto`, and `targets/grpc.decisionFromProto`. `targets` cannot import `dispatch` (cycle).

**Hot path:** `invoke_sync`.

**Invariants:**
- one Kind table in `internal/decision` (ToProto / FromProto / Neutral)
- no `targets` → `dispatch` import
- no re-export aliases in `dispatch`
- keep `promptFromProto` Kind switch in hookedge (wire-only)

**Corner cases (`tt.name`):** FromProto `nil`, `no decision`, `deny+extras`, `ask+extras`, `allow+extras`, `block_prompt`; round-trip `deny`, `ask`, `allow`, `block_prompt`, `no decision`.

**Out of scope:** SyncInvoker contract change; process globals; F3/F4.

### F2 checklist

- [x] f2-intent
- [x] f2-pkg + move ToProto/Neutral/FromProto
- [x] f2-callers (engine, hookedge, grpc, server)
- [x] f2-tests + DESIGN/AGENTS
- [x] f2-verify

**F2 done.**

---

## F3 intent note — Subscribe filter table

**Problem:** Four structurally identical Subscribe filter tests should be one table (R10 leftover).

**Hot path:** `other` (test-only).

**Invariants:** per-row hub via `newSessionTestServer`; keep Cancel / IdleThenEvent / NilHub separate.

**Corner cases (`tt.name`):** `no filter`, `provider`, `session`, `source`.

**Out of scope:** rewriting Cancel/IdleThenEvent/NilHub; F4 coverage.

### F3 checklist

- [x] f3-intent
- [x] f3-table
- [x] f3-verify

**F3 done.**

---

## F4 intent note — Coverage toward 70%

**Problem:** Repo cover ~61.5% vs ≥70% goal. Fill real behavior holes only.

**Hot path:** `other` (test-only).

**Invariants:** no filler; stop at init/Cobra/`_other.go`; measure before writing tests.

**Out of scope:** padding to 100%; new CLI commands; drive-by refactors.

### F4 checklist

- [x] f4-intent
- [x] f4-measure + gap list
- [x] f4-cmd tables
- [x] f4-targets / trajectory redact (daemon skipped — detach/import_watch/signals)
- [x] f4-stop + verify + handoff

### F4 coverage gap list (2026-08-23 measure; start total **61.8%**)

| Package | Cover | Gaps to fill (behavior) | Skip (boilerplate / hard) |
|---------|-------|---------------------------|---------------------------|
| `cmd` | 66.9% → **76.1%** | `session list` empty/json/provider/unknown; `config patch` missing file / unavailable; `record-decision` invalid scope / unavailable; `install` unknown provider; `hook` notify/run argv errors | `Execute`; hook success → `os.Exit` |
| `dispatch/targets` | 61.8% → **71.7%** | `NewAsyncInvoker` kinds; `EventKindOf`; `ProjectRootOf`; log levels | `dialAndInvoke` real dial |
| `daemon` | 65.0% | (none this phase) | `detach`, `import_watch`, unix signals |
| `trajectory` | 68.2% → **72.0%** | `PrepareRaw` / `PrepareTranscriptText` | recorder accessors via other pkgs |
| `decision` | — → **100%** | new package | — |

**F4 stop:** remaining holes are `gen/`, `version`, daemon detach/import_watch, hook `os.Exit` success paths, `dialAndInvoke`. **No filler.** Repo total **62.8%** (vs ≥70% — exception accepted, same rule as R11).

**F4 done.**
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

Full phases + acceptance: below (§ Milestones archive).

## Milestones archive

Historical acceptance criteria (M0–M12 **done**). For architecture, see [DESIGN.md](./DESIGN.md).

### M0–M7 acceptance

| Milestone | Acceptance |
|-----------|------------|
| M1 | `daemon start\|status\|reload\|stop` + `hook run --provider=…` round-trip |
| M2 | Dispatch (parallel/after_sync), async queue, secrets Ask/Deny on tool.pre, `dispatch routes`, `e2e-m2` |
| M3 | Declarative `dispatch:` + async exec/http/log/file; fsnotify reload; `e2e-m3` |
| M4 | `target: grpc`; `hook serve`/`notify` + agenthooks sentinel; `install`; Windows SID pipe; `e2e-m4` |
| M5 | Four-layer merge; ConfigService Get/Patch; `config validate\|show\|patch`; fingerprint; `e2e-m5` |
| M6 | `guards.shell`/`mcp`/`paths`; route `guards: [...]` subset; `e2e-m6` |
| M7 | Approve once → allow within TTL; restart reloads approvals; expired gone; `e2e-m7` |

### M8 — v1 gate (done)

v1 exit: four-layer config + RecordDecision; all guards; sync+async dispatch (exec async-only); install + hook run/notify/serve; unix + Windows IPC; no `not implemented` on documented commands; README/DESIGN match; lint + race + `e2e-m8` green; release published.

### M9 — Trajectory P0 (done)

L0 live ledger for all six providers; contiguous `seq` on hook run/notify/serve fixtures; no disk I/O in Decide; `session export`; default off + §14.3 matrix in docs; `e2e-m9` green.

### M10 — Trajectory P1 (done)

`session search`; Claude import with `source=transcript`; non-Claude importer status explicit; append-only merge; `e2e-m10` green.

### M11 — Trajectory P2 (done)

Every provider importer row implemented or `none`/`partial`; policy replay all six wire dialects; fork with `parent_session`; `e2e-m11` green.

### M12 / v0.0.2 — Trajectory P3 (done)

`SessionService.Subscribe`; schema versioned + `ignorable`; honest product copy; M9–M11 met; `e2e-m9…m12` green.

## Session notes

- M0–M12 / v0.0.2 shipped (see Milestones archive above)
- **R1–R11 done** — R-series complete; no R12
- **F1–F4 done** — CONVENTIONS follow-through; coverage **62.8%** (exception vs ≥70%; no pad)

## Verify (last green)

```bash
make lint
make intent-check
make docs-check
make test
go test ./... -coverprofile=/tmp/agentd_f.out -count=1
go tool cover -func=/tmp/agentd_f.out | tail -1  # total 62.8%
```

**F-series files touched:** `internal/server/{invoke,server,daemon_test,invoke_mapping_test,session_test}.go` (deleted `invoke_mode_test.go`), `internal/trajectory/{trajectory,event}.go` + `redact_test.go`, `internal/decision/{decision,decision_test}.go`, `internal/dispatch/engine.go` (deleted `decision.go`/`decision_test.go`), `internal/hookedge/decision.go`, `internal/dispatch/targets/{grpc,targets_test}.go`, `cmd/{session_list,config_patch,config_record_decision,install,hook}_test.go`, `AGENTS.md`, `DESIGN.md` §1.5, `PROGRESS.md`

## Blockers

(none)
