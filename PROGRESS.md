# agentd — implementation progress

> Session handoff for agents. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones) · Trajectory: [DESIGN.md §14](./DESIGN.md#14-trajectory-hub) · Rules: [AGENTS.md](./AGENTS.md).

## Current

**v0.0.9-beta** — release prep ready; tag after commit.

### Release checklist (commit locally)

1. **Review** unstaged/staged diff for release metadata:
   - `CHANGELOG.md` — `## [v0.0.9-beta] — 2026-09-02`
   - `README.md`, `docs/README.md` — current release pointer
   - `DESIGN.md` — `Shipped: v0.0.9-beta`, M20 milestone tag
2. **Verify** (already green on 2026-09-02):

```bash
make lint && make intent-check && make docs-check && make test && make e2e
```

3. **Commit** (example message below).
4. **Tag** and push:

```bash
git tag -a v0.0.9-beta -m "v0.0.9-beta"
git push origin main
git push origin v0.0.9-beta
```

5. **Release workflow** — push tag triggers `.github/workflows/release.yml` (goreleaser + notes from `scripts/release-notes.sh`).

### Suggested commit message

```
release: v0.0.9-beta

Policy reliability on the daemon path (policy.fail, ask_fallback, notify/serve
Cwd, project cache) plus e2e wire coverage for M15/M16/M18/M20 (M19 in e2e-m15).
```

### v0.0.9-beta scope

- **Fixed** — `policy.fail`, `ask_fallback`, hook `Cwd`, `projectsMu` cache, docs-check, CI platform tests
- **Removed** — dead `policy.unsupported` config key
- **Testing** — `e2e-m15` (M15 + M19), `e2e-m16`, `e2e-m18`, `e2e-m20`

## Recent (done)

| When | Phase | One-liner |
|------|-------|-----------|
| 2026-09-02 | v0.0.9-beta prep | CHANGELOG/README/DESIGN release metadata |
| 2026-09-02 | E2E M15–M20 | `e2e-m15`/`m16`/`m18`/`m20` + `e2e_expect_exit`; DESIGN §13 M20 row |
| 2026-09-02 | Policy/reliability | policy.fail + ask_fallback; unsupported removed; Cwd; projectsMu; docs-check + platform CI |
| 2026-09-01 | v0.0.8-beta | Cursor stats per-generation sum |

## Blockers

(none)
