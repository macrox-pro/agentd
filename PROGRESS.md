# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**v0.0.4** (2026-08-29) · shipped · next: **none**

Login autostart (`daemon enable`/`disable`) + curated config toggles (`config enable`/`disable`/`get`). Tag `v0.0.4` — push to trigger goreleaser.

```bash
make lint && make intent-check && make docs-check && make test
./scripts/release-notes.sh v0.0.4
```

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-08-29 | v0.0.4 | M14 autostart + config toggles; CHANGELOG + docs version bump |
| 2026-08-29 | Gemini research | `research/gemini/` hooks/MCP/settings/skills/GEMINI.md + T2 managed agents + migration delta |
| 2026-08-29 | M14 autostart | `daemon enable/disable`, status autostart JSON, `e2e-m14` |
| 2026-08-29 | user config bootstrap | `PrepareUserConfig` on daemon start; omitempty YAML; `LayerYAML` normalize |
| 2026-08-26 | state-directory docs | XDG state dir + cross-links EN/RU |
| 2026-08-26 | version BuildInfo | `version.Resolve` / `String` from `debug.BuildInfo` |
| 2026-08-25 | policy.offline hook edge | `OfflineFor`, `DialReady`, serve offline cache |
| 2026-08-23 | CLI version | `agentd version`; status RPC version unchanged |
| 2026-08-23 | install UX | scope defaults + stdout report |

## Refactoring (done)

**R1–R11** — provider ids; trajectory/importer registry; `dispatch/targets` factory; `first_conclusive`; guard registry; server ports; table-driven tests; tier-1 package comments. **F1–F4** — `server_test` black-box; `internal/decision` codec; subscribe filter table; coverage **62.8%** (≥70% goal exception, no filler).

| R | Done |
|---|------|
| R1 | `internal/provider` |
| R2 | trajectory typed ids |
| R3 | importer registry |
| R4–R5 | cmd + daemon/hookclient tests |
| R6–R8 | targets factory, sync merge, guard registry |
| R9–R11 | server ports, test hygiene, package comments |

Acceptance detail: git history / PRs — do not re-expand intent notes here.

## Milestones archive

Historical acceptance (M0–M14 shipped). Architecture: [DESIGN.md](./DESIGN.md).

| Milestone | Acceptance |
|-----------|------------|
| M1 | `daemon start\|status\|reload\|stop` + hook round-trip |
| M2 | Dispatch parallel/after_sync, async queue, secrets, `e2e-m2` |
| M3 | Declarative `dispatch:`, fsnotify reload, `e2e-m3` |
| M4 | gRPC target, hook serve/notify, install, Windows pipe, `e2e-m4` |
| M5 | Four-layer merge, ConfigService, `config` CLI, `e2e-m5` |
| M6 | guards shell/mcp/paths, route `guards:`, `e2e-m6` |
| M7 | Approvals TTL + persist, `e2e-m7` |
| M8 / v0.0.1 | v1 gate: full guards + dispatch + IPC + `e2e-m8` release |
| M9 | Trajectory P0 ledger, `session export`, `e2e-m9` |
| M10 | `session search`, Claude import, `e2e-m10` |
| M11 | All importer rows, replay/fork, `e2e-m11` |
| M12 / v0.0.2 | `SessionService.Subscribe`, `e2e-m9…m12` |
| M13 / v0.0.3 | `policy.offline` hook edge, docs EN/RU |
| M14 / v0.0.4 | `daemon enable\|disable`, config toggles, `e2e-m14` |

## Verify (repo green)

```bash
make lint && make intent-check && make docs-check && make test
```

## Blockers

(none)
