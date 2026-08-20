# agentd — implementation progress

> Session handoff for agents. Roadmap to v1: [DESIGN.md §13](./DESIGN.md#13-milestones). Rules: [AGENTS.md](./AGENTS.md).

## Current phase

Phase: **m7** | Last: m6-f-checkpoint | Next: m7-a runtime approvals schema

## agents_md_ready

true

## Roadmap (summary)

| Milestone | Status | One-liner |
|-----------|--------|-----------|
| M0–M6 | **done** | Daemon through config layers, secrets+shell+mcp+paths guards |
| **M7** | planned | Approvals / RecordDecision, runtime persist, temporary blocks |
| M8 / v1 | planned | Overflow counters, conformance, docs freeze, release |

Full phases + acceptance: [DESIGN.md §13](./DESIGN.md#13-milestones).

## M6 checklist

### Phase A — Schema + compile

- [x] m6-a-schema — YAML `guards.shell` / `mcp` / `paths`
- [x] m6-a-types — compiled types + parse
- [x] m6-a-merge — field-level overlay
- [x] m6-a-defaults — enabled:false opt-in
- [x] m6-a-compile — allowlist + default routes attach enabled set
- [x] m6-a-test
- [x] m6-a-checkpoint

### Phase B — Shell

- [x] m6-b-shell — deny_patterns / ask_on substring on ToolShell
- [x] m6-b-attach — AttachShell (ToolPre + Permission)
- [x] m6-b-test
- [x] m6-b-checkpoint

### Phase C — MCP

- [x] m6-c-mcp — deny_servers glob on Tool.MCP.Server
- [x] m6-c-attach — AttachMCP
- [x] m6-c-test
- [x] m6-c-checkpoint

### Phase D — Paths

- [x] m6-d-paths — deny_read / deny_write with `**` globs
- [x] m6-d-attach — AttachPaths
- [x] m6-d-test
- [x] m6-d-checkpoint

### Phase E — Builtin attach

- [x] m6-e-builtin — route `guards: [...]` switch
- [x] m6-e-test — subset + multi-guard
- [x] m6-e-checkpoint

### Phase F — Close M6

- [x] m6-f-e2e — `scripts/e2e-m6.sh`
- [x] m6-f-makefile — wired into `make e2e`
- [x] m6-f-lint-test
- [x] m6-f-docs
- [x] m6-f-checkpoint

**M6 acceptance:** met.

## Later (do not start until M6 checkpoint)

### M7 — Approvals

- Runtime `approvals` + `blocks.temporary`
- `RecordDecision` + TTL (project 24h, session end)
- Debounced runtime.yaml flush
- `scripts/e2e-m7.sh`

### M8 / v1

- Overflow drop counter on Status
- Provider timeout margin polish
- agenthookstest / integration build tag
- Docs freeze + GitHub release binaries
- `scripts/e2e-v1.sh` + v1 exit criteria in DESIGN §13

## Session notes

- AGENTS.md / CONVENTIONS.md read: yes (2026-08-20)
- M6 complete: shell/mcp/paths schema+compile, `internal/guard/{shell,mcp,paths}.go`, builtin attach by route list, e2e-m6
- Files: `internal/config/{file,guards,merge,defaults,compile}.go`, `internal/guard/{attach,shell,mcp,paths}.go`, `internal/dispatch/targets/builtin.go`, `scripts/e2e-m6.sh`, `Makefile`
- **Convention:** `make e2e` runs all `scripts/e2e-mN.sh` — when adding `e2e-m7.sh` (etc.), append to Makefile `e2e` in the same PR

## Verify (last green)

```bash
go test ./internal/config/... ./internal/guard/... ./internal/dispatch/... ./internal/server/... ./cmd/... -race -count=1
make e2e   # includes scripts/e2e-m6.sh
```

## Blockers

(none)
