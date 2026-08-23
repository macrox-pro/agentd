# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub-post-v1). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **R6 done** | Last: r6-syncinvoker | Next: **R7 aggregation policy**

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

### Coverage baseline (2026-08-23)

| Package | Statements | Priority |
|---------|------------|----------|
| `cmd` | 60.7% | maintain — R4 table tests done |
| `internal/daemon` | 65.0% | maintain — R5 lifecycle tests done |
| `internal/hookclient` | 81.0% | maintain — R5 client tests done |
| `internal/dispatch/targets` | 60.7% | maintain — R6 SyncInvoker done |
| `internal/trajectory/importer` | 71.8% | maintain — registry done |
| `internal/trajectory` | 68.3% | P2 — was 66.1% |
| `internal/dispatch` | 73.7% | maintain — R6 SyncInvoker done |
| `internal/hookedge` | 70.4% | maintain |
| `internal/config` | 74.8% | maintain |
| `internal/transport` | 77.5% | maintain |
| `internal/server` | 78.3% | maintain |
| `internal/guard` | 85.9% | maintain |
| `internal/provider` | 93.0% | R1 target package |

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
| **R7** | **next** | `dispatch` — extract `first_conclusive` aggregation |
| **R8** | pending | `guard` + `targets/builtin` — Checker registry (incremental) |
| **R9** | pending | `server` — narrow ports only if tests require |
| **R10** | pending | Test hygiene — table-driven migration where structure matches |
| **R11** | pending | Tier-1 package comments audit + docs sync |

---

## R1 intent note — Provider single source of truth

**Problem:** Provider id parsing/mapping duplicated across `hookedge`, `install`, `cmd/session_*`, `dispatch/decode`, `trajectory/session_key` — drift risk and repeated switches.

**Hot path:** `other` (CLI, install, offline import); read-only `Lookup` on `invoke_sync` where enum already validated.

**Invariants:**
- Canonical ids remain lowercase hyphenated (`claude-code`, …)
- `internal/provider` owns Parse/Lookup/Proto/Agenthooks/FromProto — no wire I/O
- Callers import `provider` directly — no passthrough wrappers in `cmd/`

**Corner cases (test names):**
- `TestParse` / `TestParseFilter` / `TestLookup` / `TestProtoRoundTrip` / `TestAgenthooksRoundTrip` / `TestFromProto` (`provider_test.go`)
- `TestSessionProviderCLI` (`cmd/session_provider_test.go`) — list/import validation rows
- `TestCanonicalProvider` — still passes via `trajectory` → `provider.Lookup`

- `TestCompileTrajectoryImportProviderAlias` — config YAML alias → canonical import key (`trajectory_test.go`)
- `TestConformanceFixtures` — all six providers via `provider.Parse`→`Agenthooks()` (`conformance_test.go`)

**Out of scope:** Changing proto enum values; Kimi alias policy beyond existing `kimicode`.

### R1 checklist

- [x] r1-intent — intent note (this section)
- [x] r1-package — `internal/provider/{provider,errors}.go` + table tests
- [x] r1-hookedge — `hookedge/run|notify|serve` use `provider.Parse`; delete `hookedge/provider.go`
- [x] r1-install — `install/run.go` uses `provider.Agenthooks()`
- [x] r1-cmd-session — `session_{list,show,export,search,replay,fork,import,subscribe}.go` use `provider.Parse` / `ParseFilter`
- [x] r1-dispatch-decode — `dispatch/decode.go` uses `provider.FromProto`
- [x] r1-trajectory-key — `CanonicalProvider` delegates to `provider.Lookup`
- [x] r1-dispatch-grpc — `dispatch/targets/grpc.go` uses `provider.Parse`
- [x] r1-config-trajectory — `config/trajectory.go` uses `provider.Lookup` + `provider.*` constants; `TestCompileTrajectoryImportProviderAlias`
- [x] r1-conformance — hookedge conformance table all six providers; `Parse`→`Agenthooks()` in loop
- [x] r1-verify — `make lint` + `make intent-check` + race tests on provider/cmd/hookedge/config

**R1 acceptance:** zero `switch` on raw provider strings outside `internal/provider` and the importer registry table in `internal/trajectory/importer/registry.go`.

---

## R2 intent note — Trajectory domain boundaries

**Problem:** Trajectory still carries provider as bare `string`; errors scattered; `ReplayPolicyFromConfig` and grpc event mapping live beside unrelated concerns.

**Hot path:** `async_side` (append, hub publish); offline `other` (replay, fork, export).

**Invariants:** append-only; `schema_version` on append; Subscribe never blocks Invoke; grpc mapping lossless for `Data`/`Raw`/`Ts`.

**Corner cases (test names):**
- `TestEventSessionEventRoundTrip` — nil, empty, zero TS, large Raw (`session_event_grpc_test.go`)
- `TestReplayPolicyFromConfig` — bad config path, missing session, seq bounds
- `TestResolveSessionKeyUsesCanonicalProvider` — alias + weak id stability
- `TestProviderImporterStatusTable` — all six providers × status enum (moved to `importer/status_test.go` in R3)

**Out of scope:** Importer file layout; Persister.Close; changing on-disk JSONL layout.

### R2 checklist

- [x] r2-intent — intent note finalized in PR
- [x] r2-errors — extend `trajectory/errors.go`; replace ad-hoc `fmt.Errorf` where message is static
- [x] r2-session-key — `ResolveSessionKeyID`; `TestResolveSessionKeyUsesCanonicalProvider` alias + weak-id rows
- [x] r2-importer-status — `TestProviderImporterStatusTable`; `provider.ID` switch in `importer_status.go`
- [x] r2-replay-config — `replay_config_test.go`; `cmd/session_replay.go` sole caller; sentinel mapping in cmd
- [x] r2-grpc-map — `TestEventSessionEventRoundTrip` table in `session_event_grpc_test.go`
- [x] r2-verify — `go test ./internal/trajectory/... -race` + root coverage **68.3%** (≥ 68%)

---

## R3 intent note — Importer registry

**Problem:** `importer/import.go` duplicates provider switch + `"claude-code"` string literals for `ProjectsRoot`; adding a provider touches multiple files.

**Hot path:** `async_side` (daemon watcher); offline `session import`.

**Invariants:** no invented L2 fields; `LastLineIndex` = physical line index; status enum honest per §14.6.

**Corner cases (test names):**
- `TestImportDispatchesByProvider` — supported / partial / none
- `TestImportSetsProjectsRootFromConfig` — per-provider cfg path
- `TestImportSessionFacade` — `import_session.go` single entry for CLI + watcher
- Existing codex/cursor/claude tests must stay green unchanged

**Out of scope:** New provider formats; Codex/Cursor fsnotify watchers; changing testdata fixtures.

### R3 checklist

- [x] r3-intent
- [x] r3-registry — `map[provider.ID]importerEntry` in `importer/registry.go`
- [x] r3-import-facade — `ImportSession(provider.ID, …)`; CLI + daemon watcher call facade
- [x] r3-status-colocate — `importer/status.go` reads registry; deleted `trajectory/importer_status.go`
- [x] r3-delete-switch — dispatch + projectsRoot via registry only
- [x] r3-verify — `go test ./internal/trajectory/importer/... -race` + importer **71.8%** (≥ 70%)

---

## R4 intent note — cmd coverage & CLI tables

**Problem:** `cmd` at **35.8%** — most `RunE` paths (hook, session fork/replay/export, config) untested; validation logic hard to regress.

**Hot path:** `other` (all mgmt CLI).

**Invariants:** validation stays in subcommand file; tests use `package cmd_test` + `RootCommand()`; no business logic added to `cmd/`.

**Corner cases (test names):**
- `TestSessionProviderCLI` — list/import validation + `list json importer_status`
- `TestSessionShowCLI` / `TestSessionExportCLI` / `TestSessionSearchCLI` / `TestSessionForkCLI` / `TestSessionImportCLI` — one table per subcommand family
- `TestSessionReplayPolicy` — `--policy` required; missing `--session`; no raw; happy hits
- `TestSessionSubscribeFilter` — provider validation + dial failure rows only
- `TestConfigValidate` / `TestConfigShow` / `TestDispatchRoutes` — table-driven stdout/JSON

**Out of scope:** E2E replacement; testing `main()`; hook integration (hookedge owns wire); `session subscribe` streaming happy path (R5).

### R4 checklist

- [x] r4-intent
- [x] r4-helper — shared harness in `cmd/root_test.go` (`resetCommandFlags`, `executeRoot`, ledger fixtures)
- [x] r4-session-table — one table per subcommand family (errors + one happy path)
- [x] r4-config-dispatch — table tests for validate/show/routes
- [x] r4-verify — `go test ./cmd/... -race` + cmd **60.7%** (≥ 55%)

---

## R5 intent note — daemon & hookclient lifecycle

**Problem:** `daemon` **47.1%**, `hookclient` **52.4%** — lock/PID/reload and client dial failures under-tested.

**Hot path:** `config_reload` (SIGHUP debounce); mgmt CLI → gRPC.

**Invariants:** one daemon per user; lock released on stop; client uses `--socket` / default path from `transport.DefaultSocketPath()`.

**Corner cases (test names):**
- `TestStartStopTable` — foreground start/stop, already running (stub health), stale PID, lock held, listen error
- `TestStopTable` — not running, clean shutdown, stale dead PID, own-PID timeout, shutdown RPC
- `TestStatusTable` / `TestReloadRPC` — stub server rows + default socket
- `TestReloadSIGHUP` — generation bump after SIGHUP (unix)
- `TestHookclientDial` / `TestHookclientDaemonRPC` / `TestSubscribeClient` — dial, mgmt RPC, subscribe cancel
- `TestDaemonStatusCLI` / `TestDaemonStopCLI` / `TestDaemonReloadCLI` — cmd daemon subcommands

**Out of scope:** Real system systemd; Windows service integration beyond existing `_windows.go` splits; detached `Start` (`Foreground: false`).

### R5 checklist

- [x] r5-intent
- [x] r5-daemon-table — start/stop/status/reload error rows + SIGHUP
- [x] r5-hookclient — dial, Status/Reload/Shutdown, Subscribe cancel (`client_test.go` only)
- [x] r5-cmd-daemon — `cmd/daemon_{status,stop,reload}_test.go` + subscribe streaming happy row
- [x] r5-verify — daemon **65.0%**, hookclient **81.0%**; `-race -timeout=20s` green

---

## R6 intent note — SyncInvoker boundary

**Problem:** Sync target Kind switch lives in `Engine.runSync` while async already uses `NewAsyncInvoker`; adding sync kinds forces Engine edits and duplicates the Kind→impl map.

**Hot path:** `invoke_sync` (call-site in `dispatch`; factory in `targets`).

**Invariants:**
- Identical Decide outcomes for same Snapshot + payload (builtin / grpc / fail_open / fail_closed / unknown kind skip)
- Sync response never waits on async queue
- No disk I/O on Invoke; `Event.Raw` unchanged
- Kind→impl mapping has **one** home: `targets` factory (not Engine)
- `OnError` applies to grpc sync only; `GRPC.InvokeSync` returns raw errors (FailMode is a factory wrapper)

**Corner cases (test names):**
- `TestNewSyncInvoker` — tt: `builtin`, `grpc`, `log_not_sync`, `unknown_kind`
- `TestEngine_RunSyncParity` — tt: `builtin_deny_first`, `grpc_fail_open`, `grpc_fail_closed`, `skip_non_sync_kind`, `first_conclusive_stops`, `empty_sync_list_neutral`
- `TestEngineHybrid` — sync deny + async fan-out; async must not change sync decision
- `TestGRPCInvokeSync` / `peer error` — raw error from `GRPC.InvokeSync` (FailMode not at that layer)

**Out of scope:** aggregation extract (R7); guard registry (R8); YAML merge-policy schema; new target kinds; fail_mode semantic changes; `config.Store` rewrite; trajectory; cmd/.

**Old R6 disposition (parked):**
- `r6-decode-table` / `TestDecodeProvider` — done-by-R1 (`provider.FromProto`); residual → R10
- `r6-targets-grpc` / `TestGRPCTargetForward` — done-by-R1 (`r1-dispatch-grpc`); residual → R10
- dispatch ≥ 75% hard gate — dropped; no-regress vs baseline

### R6 checklist

- [x] r6-intent — SyncInvoker intent + R6–R11 roadmap remap
- [x] r6-agents-targets-row / r6-conventions-kind-home
- [x] r6-sync-iface + builtin adapter + grpc FailMode wrapper + `NewSyncInvoker`
- [x] r6-engine-callsite — `runSync` uses factory; no Kind switch in Engine
- [x] r6-test-factory / r6-test-parity / r6-test-hybrid
- [x] r6-verify — race + lint + intent-check; dispatch **73.7%**, targets **60.7%** (no regress)
- [x] r6-design / package comments / PROGRESS handoff → R7

**R6 acceptance:** zero `switch` on target Kind in `Engine`; Kind→impl only in `targets.NewSyncInvoker` / `NewAsyncInvoker`; Decide parity for SyncInvoker corner cases.

**Files touched:** `internal/dispatch/{dispatch,engine}.go`, `internal/dispatch/engine_test.go`, `internal/dispatch/targets/{sync,factory,grpc,targets}.go`, `internal/dispatch/targets/targets_test.go`, `AGENTS.md`, `CONVENTIONS.md`, `DESIGN.md` §2, `PROGRESS.md`.

---

## R7 intent note — Aggregation policy (roadmap stub)

**Problem:** `first_conclusive` is inline in `runSync`; DESIGN names other merge policies.

**Hot path:** `invoke_sync`.

**Out of scope until phase starts:** implementing `all_restrictive` / `sequential_neutral_merge`; YAML schema unless compile already has a field.

### R7 checklist

- [ ] r7-intent (full note when phase starts)
- [ ] r7-extract — named aggregation type/func
- [ ] r7-tests — stop-on-conclusive + fail modes
- [ ] r7-verify

---

## R8 intent note — Guard Checker registry (roadmap stub)

**Problem:** Guard wiring in `targets/builtin` is an ad-hoc switch; wants incremental named registry.

**Hot path:** `invoke_sync`.

**Out of scope until phase starts:** YAML schema change; breaking agenthooks `Runner` wiring in one jump.

### R8 checklist

- [ ] r8-intent (full note when phase starts)
- [ ] r8-registry — incremental Checker registry
- [ ] r8-verify

---

## R9 intent note — Server narrow ports (roadmap stub)

**Problem:** Optional test ports at `server` ↔ `dispatch` / `config`.

**Out of scope:** god-interface over `Store`; add only if a concrete unit test requires it.

### R9 checklist

- [ ] r9-intent (full note when phase starts)
- [ ] r9-ports — only if needed
- [ ] r9-verify

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

**Problem:** New packages (`provider`, trajectory splits, SyncInvoker boundary) need Tier-1 comments; user-facing behavior unchanged but boundaries shifted.

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

**R-series done when:** R1–R11 checklists complete, total coverage ≥ 70%, no duplicate provider switches outside owner tables, user tags v1.1.0.

---

## Codex rollout importer intent note

**Problem:** `ImportCodex` treated Codex like Claude JSONL and resolved `{session_id}.jsonl`; real sessions live at `~/.codex/sessions/YYYY/MM/DD/rollout-*-{session_id}.jsonl` with `{timestamp,type,payload}` → 0 events + misleading checkpoint.

**Hot path:** `async_side` (offline CLI import; no Codex fsnotify watcher).

**Invariants:**
- Append-only; never invent thinking/tool-output
- Thinking only from plaintext `event_msg.agent_reasoning`
- `source=transcript`; `LastLineIndex` = file line index (including skipped lines)

**Corner cases (test names):**
- `TestResolveCodexBySessionID` (`nested rollout newest`, `not found`, `explicit path`, `exact sid.jsonl absent`)
- `TestImportCodexSessionIDFromMeta`, `TestImportCodexSessionIDFromFilenameFallback`
- `TestImportCodexRolloutMapsMessagesAndTools`, `TestImportCodexSkipsEncryptedReasoning`
- `TestMapCodexRolloutLineTable` (`custom_tool_call`, `agent_reasoning empty`/`present`, `skip meta telemetry`, `malformed line`, …)
- `TestImportCodexStartIndexResume`, `TestImportCodexEmptyFile`
- `TestProviderImporterStatus` / `codex` → `supported`

**Out of scope:** Codex daemon watcher; Cursor/Claude importer fixes; agenthooks changes.

**Status:** **done** — verify: `go test ./internal/trajectory/... -race`, `make intent-check`, `make docs-check`, `make lint`, `scripts/e2e-m11.sh`.

**Files touched:**
- `internal/trajectory/importer/codex.go`, `codex_map.go`, `codex_test.go`, `codex_map_test.go`, `importer.go`, `event.go`, `claude_code.go`, `cursor.go`, `import.go`, `testdata/codex_transcript.jsonl`
- removed: `map.go`, `resolve.go` (merged into provider files)
- `internal/trajectory/importer_status.go`, `import_checkpoint_test.go`
- `cmd/session_import.go`, `scripts/e2e-m11.sh`
- `DESIGN.md`, `docs/en|ru/{trajectory,configuration,cli}.md`, `PROGRESS.md`

## M11 intent note

**Problem:** Trajectory P2 — Cursor (and Codex when path known) L2 import; offline policy replay for all six wire dialects; audit-only session fork.

**Hot path:** importers/watcher → `async_side`; replay/fork CLI → `other` (offline).

**Invariants:**
- Append-only ledger; original immutable on fork
- No invented thinking/tool-output from transcripts
- Policy replay never talks to a live agent; requires stored `Raw` (`include_raw`)
- No disk I/O on sync Invoke

**Corner cases (test names):**
- `TestImportCursor`, `TestImportCodex`, `TestProviderImporterStatus`
- `TestReplayPolicy`, `TestReplayMissingRaw`
- `TestForkSession`, `TestForkDuplicateRejected`

**Out of scope:** agent-loop resume; inventing L3 events; M12 Subscribe.

## Roadmap (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M7 | **done** | Daemon through approvals / RecordDecision / runtime persist |
| M8 / v1 | **done** | Overflow counters, conformance, docs freeze, release |
| **M9** | **done** | Trajectory P0 — L0 live ledger for all six providers + export |
| **M10** | **done** | Trajectory P1 — search + Claude import; others stay L0 |
| **M11** | **done** | Trajectory P2 — importers if format exists; policy replay all dialects |
| **M12 / v1.1** | **done** | Trajectory P3 — Subscribe; contract freeze; depth = §14.6 matrix |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## M11 checklist — Trajectory hub P2

### Phase A — Spike + fixtures

- [x] m11-a-intent — intent note
- [x] m11-a-codex-spike — agenthooks parses codex → **partial**
- [x] m11-a-cursor-fixture — testdata/cursor_transcript.jsonl

### Phase B — Config

- [x] m11-b-config — trajectory.import.cursor + codex
- [x] m11-b-config-test

### Phase C — Importers

- [x] m11-c-map — importer/map.go shared mapping
- [x] m11-c-cursor — importer/cursor.go
- [x] m11-c-codex — importer/codex.go
- [x] m11-c-import-cli — provider dispatch in session import

### Phase D — Status

- [x] m11-d-status — supported/partial/none for all six
- [x] m11-d-watcher — StartIndex resume fix only (Claude watch root unchanged)

### Phase E — Policy replay

- [x] m11-e-replay-core — internal/trajectory/replay.go
- [x] m11-e-replay-cli — session replay --policy
- [x] m11-e-replay-tests — six providers + missing Raw

### Phase F — Fork

- [x] m11-f-events — session/fork, session/end-seed
- [x] m11-f-fork-core — internal/trajectory/fork.go
- [x] m11-f-fork-cli — session fork
- [x] m11-f-fork-tests

### Phase G–H — Docs + e2e

- [x] m11-g-docs — DESIGN §6/§13/§14.6; docs en/ru
- [x] m11-h-e2e — scripts/e2e-m11.sh
- [x] m11-h-verify

**M11 acceptance:** see [DESIGN.md §13 M11](./DESIGN.md#m11--trajectory-p2-multi-import--policy-replay).

## M12 intent note

**Problem:** External Trajectory UIs cannot live-tail the ledger; schema/docs not frozen for v1.1.

**Hot path:** Fan-out after in-memory append → `async_side`; Subscribe RPC / CLI → `other`.

**Invariants:**
- Subscribe never blocks Invoke; append-only; `ignorable` = forward-compat hint (not Subscribe filter)
- `schema_version: 1` on all appended/streamed events
- Honest §14.6 claim; no agent-loop

**Corner cases (test names):**
- `TestHubDeliverAfterAppend`, `TestHubSlowConsumerDrop`, `TestHubUnregister`, `TestHubMultipleSubscribers`, `TestHubPublishNoSubscribers`, `TestHubCloseEndsSubscribers`
- `TestSchemaVersionOnAppend`, `TestIgnorablePreserved`
- `TestSubscribeFilterProvider`, `TestSubscribeFilterSession`, `TestSubscribeFilterSource`, `TestSubscribeNoFilter`, `TestSubscribeCancel`, `TestSubscribeIdleThenEvent`, `TestSubscribeNilHub`

**Out of scope:** grpc-gateway; webhook; `after_seq` catch-up; Codex/Cursor fsnotify watchers; importer refactors; git tag (user).

## M12 checklist — Trajectory P3 / v1.1

### Phase 0 — Process

- [x] m12-0-agents — AGENTS.md + CONVENTIONS
- [x] m12-0-intent — intent note (this section)
- [x] m12-0-checklist — segmented checklist

### Phase A1 — Proto

- [x] m12-a1-proto — `api/agentd/v1/session.proto`
- [x] m12-a1-generate — `make generate` + build
- [x] m12-a1-design4 — DESIGN §4 SessionService

### Phase A2 — Hub + schema

- [x] m12-a2-schema — SchemaVersion stamp in Store.Append + AppendEvents
- [x] m12-a2-schema-test — TestSchemaVersionOnAppend, TestIgnorablePreserved
- [x] m12-a2-hub — `internal/trajectory/hub.go`
- [x] m12-a2-hub-tests — hub corner cases
- [x] m12-a2-queue — Publish(appended) after Append
- [x] m12-a2-recorder — Hub + Close closes Hub

### Phase A3 — gRPC server

- [x] m12-a3-server — `internal/server/session.go` + register
- [x] m12-a3-server-opts — daemon wires Hub
- [x] m12-a3-server-tests — bufconn Subscribe tests

### Phase A4 — Client + CLI

- [x] m12-a4-hookclient — Subscribe helper
- [x] m12-a4-cli — `cmd/session_subscribe.go`
- [x] m12-a4-cli-wire — session.go AddCommand
- [x] m12-a4-design6 — DESIGN §6
- [x] m12-a4-docs-cli — docs en/ru cli + trajectory

### Phase A5 — Import live fan-out

- [x] m12-a5-import-hub — Hub.Publish after AppendImported (Claude watcher)
- [x] m12-a5-import-wire — ImportWatcher gets Hub

### Phase B — Mirror docs

- [x] m12-b-docs-mirror — store mirror documented; no webhook
- [x] m12-b-progress — Phase B checkbox

### Phase C — Contract freeze

- [x] m12-c-contract-en — docs/en/trajectory.md
- [x] m12-c-contract-ru — docs/ru/trajectory.md
- [x] m12-c-readme — README v1.1 + matrix
- [x] m12-c-design-accept — DESIGN §13/§14.3/§14.5/§14.8
- [x] m12-c-docs-check — `make docs-check`

### Phase D — e2e + CHANGELOG

- [x] m12-d-e2e — `scripts/e2e-m12.sh`
- [x] m12-d-changelog — CHANGELOG `[v1.1.0]`
- [x] m12-d-progress-done — m12 done + verify block
- [x] m12-d-verify — lint + intent-check + test + e2e

**M12 acceptance:** see [DESIGN.md §13 M12](./DESIGN.md#m12--v11--trajectory-p3-stream-out).

## Session notes

- M11 shipped: Cursor/Codex partial import, `session replay --policy`, `session fork`, `scripts/e2e-m11.sh`
- M12 shipped: Subscribe RPC/CLI, hub fan-out, schema_version contract, `scripts/e2e-m12.sh`
- Trajectory package refactor (pre-M12): AppendEvents sync write (no Persister leak); CanonicalProvider lowercase; importer `Import` facade + file split; thin `cmd/session_import`; DESIGN §14.8 `importer/`
- **R2 done:** sentinels in `errors.go`; `ResolveSessionKeyID`; grpc/replay_config table tests; cmd error mapping
- **R3 done:** importer registry + status colocated; `ImportSession` facade for CLI + watcher; `session list --json` enriches status in `cmd/`
- **R4 done:** shared `cmd/root_test.go` harness; table-driven session/config/dispatch CLI tests; cmd **60.7%**
- **R5 done:** daemon lifecycle tables + hookclient client tests + cmd daemon CLI; daemon **65.0%**, hookclient **81.0%**
- **R6 done:** `SyncInvoker` + `NewSyncInvoker`; `GRPCSync` FailMode wrapper; Engine Kind switch removed; Decide parity + hybrid tests; dispatch **73.7%**, targets **60.7%**
- Next agent session: open **R7** — extract `first_conclusive` aggregation from `runSync`

### Refactor intent note (pre-R-series, archived)

**Problem:** Post-M11 drift — offline AppendEvents leaked Persister goroutine; shared importer API in claude_code.go; provider switch in cmd/; stale Entry/DESIGN paths; CanonicalProvider casing.

**Hot path:** `other` (offline CLI/import/fork/replay); live `async_side` unchanged.

**Invariants:** append-only; fork immutable source; no invented L2/L3; no unused APIs; cmd Cobra-only; no kitchen-sink filenames.

**Corner cases:** `TestAppendEventsNoLeakedGoroutine`, `TestCanonicalProvider`, `TestImportDispatchesByProvider`.

**Out of scope:** M12 Subscribe; Persister.Close; RunOfflineImport; PolicyInvoker without net win; coverage marathon.

## Verify (last green)

```bash
make lint
make intent-check
go test ./internal/dispatch/... ./internal/dispatch/targets/... -race -count=1
go test ./internal/dispatch -coverprofile=/tmp/dispatch_r6.out -count=1
go tool cover -func=/tmp/dispatch_r6.out | tail -1   # dispatch 73.7%
go test ./internal/dispatch/targets -coverprofile=/tmp/targets_r6.out -count=1
go tool cover -func=/tmp/targets_r6.out | tail -1   # targets 60.7%
```

**R6 files touched:** `internal/dispatch/{dispatch,engine}.go`, `internal/dispatch/engine_test.go`, `internal/dispatch/targets/{sync,factory,grpc,targets}.go`, `internal/dispatch/targets/targets_test.go`, `AGENTS.md`, `CONVENTIONS.md`, `DESIGN.md`, `PROGRESS.md`

**R5 files touched:** `internal/daemon/{start,stop,status,reload,reload_unix,paths}_test.go`, `internal/hookclient/client_test.go`, `cmd/{root,daemon_status,daemon_stop,daemon_reload,session_subscribe}_test.go`, `PROGRESS.md`

**R4 files touched:** `cmd/root_test.go`, `cmd/session_{provider,show,export,search,fork,import,replay,subscribe}_test.go`, `cmd/config_{validate,show}_test.go`, `cmd/dispatch_routes_test.go`, `PROGRESS.md`

**R3 files touched:** `internal/trajectory/importer/{registry,status,import,import_session,importer}.go`, `*_test.go` (status, import_session, import), deleted `internal/trajectory/importer_status.go`, `internal/trajectory/{list,trajectory}.go`, `internal/daemon/import_watch.go`, `cmd/session_{import,list}.go`, `PROGRESS.md`

**R2 files touched:** `internal/trajectory/{errors,replay,fork,list,search,session_key,importer_status,trajectory}.go`, `*_test.go` (session_key, import_checkpoint, fork, session_event_grpc, replay_config), `cmd/session_{replay,show,fork}.go`, `PROGRESS.md`

## Blockers

(none)
