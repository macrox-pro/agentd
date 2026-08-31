# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**v0.0.7 release** — ready for tag.

- Setup wizard (`agentd setup`, interactive `install` on TTY; `AGENTD_NO_TUI` / `CI` bypass)
- Doctor + `install --all-detected` (plan-only default, `--yes` to apply); `e2e-m17`
- Trajectory compile defaults on (`enabled`, `include_raw`, `statistics`)
- Codex transcript tail token fallback on `Stop` when hook raw has no usage
- CHANGELOG + docs version bump to v0.0.7

```bash
make lint && make intent-check && make docs-check && make test && make e2e
# tag (user): git tag -a v0.0.7 -m "v0.0.7" && git push origin v0.0.7
# release binaries: goreleaser (GitHub Actions on tag push)
```

**Next:** Claude transcript token investigation (separate PR).

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-08-31 | v0.0.7 | Setup wizard, doctor, `--all-detected`, trajectory default-on; CHANGELOG + docs version bump |
| 2026-08-31 | M17/M18 | Doctor, `--all-detected`, setup TUI, trajectory default-on |
| 2026-08-31 | M16 Codex tokens | Transcript tail fallback on Stop; hook usage wins when present |
| 2026-08-31 | v0.0.6 | Prometheus metrics + trajectory token stats; CHANGELOG + docs version bump |
| 2026-08-31 | daemon tests | Isolate unit tests from production socket, agentd.log, runtime.yaml |
| 2026-08-31 | metrics | Prometheus scrape HTTP + Observer histograms + reload hook |
| 2026-08-29 | v0.0.5 | Trajectory counters + session import `--out`; CHANGELOG + docs version bump |

## Verify (repo green)

```bash
make lint && make intent-check && make docs-check && make test && make e2e
```

## Blockers

(none)
