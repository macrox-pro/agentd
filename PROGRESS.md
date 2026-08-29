# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**v0.0.5** — changelog + docs version bump (ready to tag).

- `agentd trajectory stats` (daemon rollup) + `agentd session stats` (offline JSONL)
- `trajectory.statistics` / toggle `trajectory-statistics`
- `session import --out`

```bash
make lint && make intent-check && make docs-check && make test
```

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-08-29 | v0.0.5 | Trajectory counters + session import `--out`; CHANGELOG + docs version bump |
| 2026-08-29 | v0.0.4 | M14 autostart + config toggles; CHANGELOG + docs version bump |

## Verify (repo green)

```bash
make lint && make intent-check && make docs-check && make test
```

## Blockers

(none)
