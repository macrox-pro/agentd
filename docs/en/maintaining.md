# Maintaining user docs

> **Language:** [English](./maintaining.md) · [Русский](../ru/maintaining.md)

How to keep `docs/en/` and `docs/ru/` accurate when code changes. EN is canonical; RU must mirror the same pages.

## When to update

| Change | Update |
|--------|--------|
| New/changed CLI command or flag (`cmd/`) | [cli.md](./cli.md) (canonical), [DESIGN.md §6](../../DESIGN.md#6-cli-reference) if architecture notes change, related how-to if UX changes |
| YAML key / enum (`internal/config/file.go`, compile) | [configuration.md](./configuration.md) (canonical), plus [guards.md](./guards.md) / [dispatch.md](./dispatch.md) / [approvals.md](./approvals.md) as needed; DESIGN §7 if layer/runtime overlay contract changes |
| Default on-disk paths (state dir, log, sessions, socket) | [configuration.md](./configuration.md#state-directory) (canonical); DESIGN §3 / §5 / §14 pointers only |
| Guard / Ask / Deny behavior | [guards.md](./guards.md), [approvals.md](./approvals.md), [troubleshooting.md](./troubleshooting.md) |
| Dispatch mode, target, timeout, async overflow | [dispatch.md](./dispatch.md), [operations.md](./operations.md) |
| Status / daemon ops fields | [operations.md](./operations.md), [cli.md](./cli.md) |
| Install providers / scopes / entrypoints / quirks | [providers.md](./providers.md) + `providers-*.md`, [getting-started.md](./getting-started.md) |
| Install / Releases / version wiring | [installation.md](./installation.md) |
| Failure modes / offline / timeouts | [troubleshooting.md](./troubleshooting.md) |
| Trajectory / session ledger / import / subscribe | [trajectory.md](./trajectory.md), [cli.md](./cli.md) (session commands) |
| User-visible README claims | [README.md](../../README.md) + matching docs page ([why.md](./why.md) if positioning changes) |

If code and DESIGN disagree, **document code** and fix DESIGN in the same change when practical.

## Rules

- Edit **EN first**, then **RU** (same filename, same section order). Identifiers stay English in RU.
- **Writing style:** [CONVENTIONS.md § Documentation style](../../CONVENTIONS.md#documentation-style) — page types, beginner terms, RU terminology table, no internal names in user docs.
- No fluff; state commands, keys, and Status fields.
- Do not invent flags or YAML keys not present in `cmd/` / `file.go`.
- Non-goals (DESIGN §11) stay non-goals — do not oversell.
- After doc edits: `make docs-check` (EN/RU filename parity).

## Source map (verify against)

| Topic | Primary source |
|-------|----------------|
| CLI | `cmd/*.go` |
| YAML | `internal/config/file.go` |
| Paths / persist | `internal/config/paths.go`, `paths_unix.go`, `paths_windows.go`, `store.go`, `persist.go` |
| Socket defaults | `internal/transport/path_*.go` |
| Guards | `internal/guard/` |
| Dispatch / timeout | `internal/dispatch/`, `timeout.go` |
| Approvals / blocks | `internal/config/approvals.go`, `blocks.go` |
| Status JSON | `internal/daemon/status_write.go`, `api/agentd/v1/daemon.proto` |
| Install | `internal/install/run.go` |
| Hook offline | `internal/hookedge/` |

## PR checklist (docs)

- [ ] User-visible behavior change → `docs/en/` updated
- [ ] Matching `docs/ru/` page(s) updated (or new page added in both)
- [ ] DESIGN §6 / §7 if CLI architecture notes or schema layer docs change
- [ ] `make docs-check` passes
