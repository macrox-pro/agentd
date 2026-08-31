# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**Documentation refactoring** — complete.

- Added `docs/en/glossary.md` + `docs/ru/glossary.md`; style rules in [CONVENTIONS.md § Documentation style](./CONVENTIONS.md#documentation-style)
- Factual fixes: `daemon restart` → stop+start; six CLI toggles; anchor/link repairs; RU translation cleanup
- EN+RU page rewrites: dispatch kind table, trajectory readability, provider uniformity, de-duplicated tables
- `scripts/check-docs-terms.sh` + `scripts/check-docs-links.sh` wired into `make docs-check`
- [README.md](./README.md) and [DESIGN.md](./DESIGN.md) updated (glossary link, IPC/TUI/NDJSON/TTL, hook edge naming)

```bash
make docs-check && make lint && make intent-check
```

**Next:** tag v0.0.7 if not yet pushed; Claude transcript token investigation (separate PR).

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-09-01 | Docs refactor | Glossary, EN/RU rewrite, docs-check linters, README/DESIGN/CONVENTIONS style |
| 2026-08-31 | v0.0.7 | Setup wizard, doctor, `--all-detected`, trajectory default-on; CHANGELOG + docs version bump |

## Verify (repo green)

```bash
make lint && make intent-check && make docs-check && make test && make e2e
```

## Blockers

(none)
