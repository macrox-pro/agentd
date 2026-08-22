# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub-post-v1). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **m12 done** | Last: m12-subscribe-v1.1 | Next: user tag v1.1.0

## agents_md_ready

true

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
- Trajectory package refactor (pre-M12): AppendEvents sync write (no Persister leak); CanonicalProvider lowercase; importer `Import` facade + file split (`import.go`, `resolve.go`, `mapEntriesFrom` in `map.go`); thin `cmd/session_import`; DESIGN §14.8 `importer/`; skipped PolicyInvoker (would still import `dispatch.InvokeInput`)
- Files touched: `internal/trajectory/*` (persist, session_key, trajectory, list, replay, import_append, queue, replay_test), `internal/trajectory/importer/*`, `cmd/session_import.go`, `DESIGN.md`
- Next: M12 Subscribe + trajectory contract freeze + v1.1 release

### Refactor intent note

**Problem:** Post-M11 drift — offline AppendEvents leaked Persister goroutine; shared importer API in claude_code.go; provider switch in cmd/; stale Entry/DESIGN paths; CanonicalProvider casing.

**Hot path:** `other` (offline CLI/import/fork/replay); live `async_side` unchanged.

**Invariants:** append-only; fork immutable source; no invented L2/L3; no unused APIs; cmd Cobra-only; no kitchen-sink filenames.

**Corner cases:** `TestAppendEventsNoLeakedGoroutine`, `TestCanonicalProvider`, `TestImportDispatchesByProvider`.

**Out of scope:** M12 Subscribe; Persister.Close; RunOfflineImport; PolicyInvoker without net win; coverage marathon.

## Verify (last green)

```bash
make lint
make intent-check
make docs-check
make test
go test -tags=integration ./internal/hookedge/ -race -count=1
make e2e   # includes scripts/e2e-m9.sh … e2e-m12.sh
```

## Blockers

(none)
