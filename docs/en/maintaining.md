# Maintaining user docs

> **Language:** [English](./maintaining.md) · [Русский](../ru/maintaining.md)

How to keep `docs/en/` and `docs/ru/` accurate when code changes. EN is canonical; RU must mirror the same pages.

## When to update

| Change | Update |
|--------|--------|
| New/changed CLI command or flag (`cmd/`) | [cli.md](./cli.md), [DESIGN.md §6](../../DESIGN.md#6-cli-reference), related how-to if UX changes |
| YAML key / enum (`internal/config/file.go`, compile) | [configuration.md](./configuration.md), plus [guards.md](./guards.md) / [dispatch.md](./dispatch.md) / [approvals.md](./approvals.md) as needed; DESIGN §7 if schema example drifts |
| Guard / Ask / Deny behavior | [guards.md](./guards.md), [approvals.md](./approvals.md), [troubleshooting.md](./troubleshooting.md) |
| Dispatch mode, target, timeout, async overflow | [dispatch.md](./dispatch.md), [operations.md](./operations.md) |
| Status / daemon ops fields | [operations.md](./operations.md), [cli.md](./cli.md) |
| Install providers / scopes / entrypoints | [providers.md](./providers.md), [getting-started.md](./getting-started.md) |
| Install / Releases / version wiring | [installation.md](./installation.md) |
| Failure modes / offline / timeouts | [troubleshooting.md](./troubleshooting.md) |
| User-visible README claims | [README.md](../../README.md) + matching docs page ([why.md](./why.md) if positioning changes) |

If code and DESIGN disagree, **document code** and fix DESIGN in the same change when practical.

## Rules

- Edit **EN first**, then **RU** (same filename, same section order). Identifiers stay English in RU.
- No fluff; state commands, keys, and Status fields.
- Do not invent flags or YAML keys not present in `cmd/` / `file.go`.
- Non-goals (DESIGN §11) stay non-goals — do not oversell.
- After doc edits: `make docs-check` (EN/RU filename parity).

## Source map (verify against)

| Topic | Primary source |
|-------|----------------|
| CLI | `cmd/*.go` |
| YAML | `internal/config/file.go` |
| Paths / persist | `internal/config/store.go`, `persist.go` |
| Guards | `internal/guard/` |
| Dispatch / timeout | `internal/dispatch/`, `timeout.go` |
| Approvals / blocks | `internal/config/approvals.go`, `blocks.go` |
| Status JSON | `internal/daemon/status_write.go`, `api/agentd/v1/daemon.proto` |
| Install | `internal/install/run.go` |
| Hook offline | `internal/hookedge/` |

## PR checklist (docs)

- [ ] User-visible behavior change → `docs/en/` updated
- [ ] Matching `docs/ru/` page(s) updated (or new page added in both)
- [ ] DESIGN §6 / §7 if CLI or schema examples changed
- [ ] `make docs-check` passes
