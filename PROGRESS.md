# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub-post-v1). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **R2 done** | Last: r2-trajectory-domain | Next: **R3 importer registry**

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
| `cmd` | 35.8% | **P0** — session/hook RunE mostly untested |
| `internal/daemon` | 47.1% | **P1** — start/stop/reload lifecycle |
| `internal/hookclient` | 52.4% | **P1** — gRPC client helpers |
| `internal/dispatch/targets` | 60.1% | P2 |
| `internal/trajectory/importer` | 63.8% | P2 — registry phase |
| `internal/trajectory` | 68.3% | P2 — was 66.1% |
| `internal/dispatch` | 69.5% | P2 |
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
- **`cmd` table tests** — reset Cobra flags between rows (`resetCommandFlags`); stale `Changed` bits cause flaky cross-row pollution.
- **`internal/` errors** — no `--provider` / flag names in sentinels; map in `cmd/RunE` only.
- **Importer registry** — must not pull `dispatch.InvokeInput` (PolicyInvoker still out of scope); registry returns `(ImportResult, error)` only.
- **No drive-by refactors** — one R-phase touches one boundary; unrelated cleanups → separate PR.

---

## R-series roadmap (summary)

| Phase | Status | One-liner |
|-------|--------|-----------|
| **R1** | **done** | `internal/provider` — single source for ids, proto/agenthooks mapping |
| **R2** | **done** | Trajectory domain — typed ids at boundaries, errors, grpc event mapping |
| **R3** | **next** | Importer registry — kill `import.go` provider switch + string literals |
| **R4** | pending | `cmd/` coverage — table-driven session + config CLI |
| **R5** | pending | `daemon` + `hookclient` lifecycle tests |
| **R6** | pending | `dispatch` decode/targets — provider mapping consolidation |
| **R7** | pending | Test hygiene — table-driven migration where structure matches |
| **R8** | pending | Tier-1 package comments audit + docs sync |

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

**Out of scope:** Changing proto enum values; Kimi alias policy beyond existing `kimicode`; importer registry (R3 — `importer/import.go` switches remain until then).

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

**R1 acceptance:** zero `switch` on raw provider strings outside `internal/provider`, `trajectory/importer_status.go`, and `trajectory/importer/` (R3 registry).

---

## R2 intent note — Trajectory domain boundaries

**Problem:** Trajectory still carries provider as bare `string`; errors scattered; `ReplayPolicyFromConfig` and grpc event mapping live beside unrelated concerns.

**Hot path:** `async_side` (append, hub publish); offline `other` (replay, fork, export).

**Invariants:** append-only; `schema_version` on append; Subscribe never blocks Invoke; grpc mapping lossless for `Data`/`Raw`/`Ts`.

**Corner cases (test names):**
- `TestEventSessionEventRoundTrip` — nil, empty, zero TS, large Raw (`session_event_grpc_test.go`)
- `TestReplayPolicyFromConfig` — bad config path, missing session, seq bounds
- `TestResolveSessionKeyUsesCanonicalProvider` — alias + weak id stability
- `TestProviderImporterStatusTable` — all six providers × status enum

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

- [ ] r3-intent
- [ ] r3-registry — `var importers = map[provider.ID]importerFunc{…}` or small registry type in `importer/`
- [ ] r3-import-facade — `ImportSession(provider.ID, ImportOptions)`; CLI calls facade
- [ ] r3-status-colocate — move `ProviderImporterStatus` next to registry (single table)
- [ ] r3-delete-switch — `import.go` dispatches via registry only
- [ ] r3-verify — `go test ./internal/trajectory/importer/... -race` + importer ≥ 70%

---

## R4 intent note — cmd coverage & CLI tables

**Problem:** `cmd` at **35.8%** — most `RunE` paths (hook, session fork/replay/export, config) untested; validation logic hard to regress.

**Hot path:** `other` (all mgmt CLI).

**Invariants:** validation stays in subcommand file; tests use `package cmd_test` + `RootCommand()`; no business logic added to `cmd/`.

**Corner cases (test names):**
- `TestSessionProviderCLI` — extend rows: show/export/search/replay/fork happy + error
- `TestSessionReplayPolicy` — `--policy` requires config; missing `--session`
- `TestSessionSubscribeFilter` — provider/session/source filter flag combos
- `TestConfigValidate` / `TestDispatchRoutes` — table-driven stdout snapshots

**Out of scope:** E2E replacement; testing `main()`; hook integration (hookedge owns wire).

### R4 checklist

- [ ] r4-intent
- [ ] r4-helper — shared `resetCommandFlags` + temp state dir helper in `cmd/*_test.go`
- [ ] r4-session-table — one table per subcommand family (errors + one happy path)
- [ ] r4-config-dispatch — table tests for validate/show/routes
- [ ] r4-verify — `go test ./cmd/... -race` + cmd ≥ 55%

---

## R5 intent note — daemon & hookclient lifecycle

**Problem:** `daemon` **47.1%**, `hookclient` **52.4%** — lock/PID/reload and client dial failures under-tested.

**Hot path:** `config_reload` (SIGHUP debounce); mgmt CLI → gRPC.

**Invariants:** one daemon per user; lock released on stop; client respects `AGENTD_SOCKET` / default path.

**Corner cases (test names):**
- `TestStartStopTable` — already running, stale PID, lock held
- `TestReloadDebounced` — rapid SIGHUP coalesced
- `TestHookclientDial` — missing socket, canceled ctx, bufconn success
- `TestSubscribeClient` — cancel propagation (extend `hookclient` tests)

**Out of scope:** Real system systemd; Windows service integration beyond existing `_windows.go` splits.

### R5 checklist

- [ ] r5-intent
- [ ] r5-daemon-table — start/stop/status error rows
- [ ] r5-hookclient — Invoke + Subscribe helpers via bufconn
- [ ] r5-verify — daemon ≥ 65%, hookclient ≥ 70%

---

## R6 intent note — dispatch boundaries

**Problem:** Residual provider/string routing in dispatch decode and grpc forward target; engine hybrid mode branches under-covered.

**Hot path:** `invoke_sync`, `async_side`.

**Invariants:** sync path never blocks on trajectory; async queue bounded; decode errors map to gRPC status in server only.

**Corner cases (test names):**
- `TestDecodeProvider` — table: six proto providers → canonical id
- `TestGRPCTargetForward` — provider header mapping
- `TestEngineHybrid` — sync deny + async fan-out row

**Out of scope:** New target types; route YAML schema changes.

### R6 checklist

- [ ] r6-intent
- [ ] r6-decode-table — provider mapping tests in `dispatch/decode_test.go`
- [ ] r6-targets-grpc — consolidate with `provider.FromProto`
- [ ] r6-engine-hybrid — table in `engine_test.go`
- [ ] r6-verify — dispatch ≥ 75%

---

## R7 intent note — Test hygiene (table-driven migration)

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

**Out of scope:** testify/suite; rewriting e2e scripts; 100% coverage chase.

### R7 checklist

- [ ] r7-intent
- [ ] r7-hub-table — hub deliver/drop/unregister merged where safe
- [ ] r7-fork-table
- [ ] r7-invoke-table
- [ ] r7-audit — grep for `if tt.` blocks inside single `Test` without `t.Run`
- [ ] r7-verify — full `make test`; total coverage ≥ 70%

---

## R8 intent note — Package comments & docs sync

**Problem:** New packages (`provider`, trajectory splits) need Tier-1 comments; user-facing behavior unchanged but boundaries shifted.

**Hot path:** n/a.

**Invariants:** package comment template per [CONVENTIONS](./CONVENTIONS.md#comments); DESIGN §1.5 hot path tags accurate.

**Corner cases:** `make intent-check` passes; no stale `hookedge/provider` references in docs.

**Out of scope:** New CLI commands; proto changes; README feature marketing.

### R8 checklist

- [ ] r8-intent
- [ ] r8-provider-comment — verify `internal/provider/provider.go` header
- [ ] r8-trajectory-comments — `hub.go`, `replay_config.go`, `session_event_grpc.go`
- [ ] r8-design-15 — DESIGN §1.5 row for `internal/provider` if missing
- [ ] r8-docs-check — `make docs-check` (only if user-visible text changed)
- [ ] r8-verify — `make intent-check` + full verify block

**R-series done when:** R1–R8 checklists complete, total coverage ≥ 70%, no duplicate provider switches outside owner tables, user tags v1.1.0.

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
- **R2 done:** sentinels in `errors.go`; `ResolveSessionKeyID`; importer status `provider.ID` switch; grpc/replay_config table tests; cmd error mapping
- Next agent session: open **R3** — importer registry (`import.go` dispatch switch)

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
go test ./internal/trajectory/... -race -count=1
go test ./internal/trajectory/... -coverprofile=/tmp/r2.cover
go tool cover -func=/tmp/r2.cover | tail -1   # trajectory 68.3%, importer 63.8%, module 66.7%
```

**R2 files touched:** `internal/trajectory/{errors,replay,fork,list,search,session_key,importer_status,trajectory}.go`, `*_test.go` (session_key, import_checkpoint, fork, session_event_grpc, replay_config), `cmd/session_{replay,show,fork}.go`, `PROGRESS.md`

## Blockers

(none)
