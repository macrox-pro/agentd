# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**v0.0.8-beta release** — ready for tag.

- Cursor trajectory stats: sum billing tokens per `stop` (per generation); drop session-delta aggregation
- CHANGELOG + docs version bump to v0.0.8-beta

```bash
make lint && make intent-check && make docs-check && make test
# tag (user): git tag -a v0.0.8-beta -m "v0.0.8-beta" && git push origin v0.0.8-beta
# release binaries: goreleaser (GitHub Actions on tag push; prerelease auto for -beta)
```

**Next:** Claude transcript token investigation (separate PR); stable v0.0.8 when beta validates.

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-09-02 | OpenCode docs research | 52 topic files under `research/opencode/`; README + SOURCES/MANIFEST/CHECKLIST |
| 2026-09-01 | v0.0.8-beta | Cursor stats per-generation sum; CHANGELOG + docs version bump |
| 2026-09-01 | Cursor stats tokens | Sum each Cursor stop (per generation); drop session delta; DESIGN §14.6 + docs |
| 2026-09-01 | Docs refactor | Glossary, EN/RU rewrite, docs-check linters, README/DESIGN/CONVENTIONS style |
| 2026-08-31 | v0.0.7 | Setup wizard, doctor, `--all-detected`, trajectory default-on; CHANGELOG + docs version bump |

## Verify (repo green)

```bash
make lint && make intent-check && make docs-check && make test && make e2e
```

## Blockers

(none)
