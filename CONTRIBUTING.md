# Contributing to agentd

Thank you for contributing to agentd. We appreciate your time and help.
Here are some guidelines to get started.

Deep rules for agents and humans shipping code: [AGENTS.md](./AGENTS.md) · [CONVENTIONS.md](./CONVENTIONS.md) · architecture: [DESIGN.md](./DESIGN.md).

## Code of Conduct

Be kind and respectful to the members of the community. Take time to educate
others who are seeking help. Harassment of any kind will not be tolerated.

## Questions

1. Read the [user guide](./docs/en/) (or [русское руководство](./docs/ru/)) and [Troubleshooting](./docs/en/troubleshooting.md).
2. Check [DESIGN.md](./DESIGN.md) for behavior contracts and the roadmap (§13 / §14).
3. If it is still unclear, open a [GitHub Issue](https://github.com/macrox-pro/agentd/issues) with what you already tried (`agentd daemon status --json`, provider, OS).

**Contact:** [garin1221@yandex.ru](mailto:garin1221@yandex.ru)

## Filing a bug or feature

1. Before filing an issue, search [existing issues](https://github.com/macrox-pro/agentd/issues) for a similar report. If one exists, comment there.
2. **Bug reports** should include:
   - Steps to reproduce
   - `agentd` version (`agentd version` / release tag / `dev` build)
   - OS (linux / darwin / windows) and Go version if building from source
   - Coding agent / `--provider` (claude-code, cursor, codex, gemini, opencode, kimi-code)
   - Relevant config snippets (redact secrets)
   - Current vs expected behavior
   - For hook failures: whether the daemon was running; avoid pasting full tool outputs that contain secrets
3. **Feature / enhancement** requests should have a clear title, what you want, and why it helps users. Note whether it fits [DESIGN.md §11](./DESIGN.md#11-non-goals-v1) / planned work (§14 Trajectory). Large ideas are easier to land if they align with an open milestone in [PROGRESS.md](./PROGRESS.md).

## Submitting changes

1. **License:** Contributions are accepted under the project [MIT License](./LICENSE).
2. **Scope:** Change only what the PR needs. No drive-by refactors, unused “for later” APIs, or hand-edits under `gen/`.
3. **Style:** Follow [AGENTS.md](./AGENTS.md) / [CONVENTIONS.md](./CONVENTIONS.md) (errors, files, tests, protobuf, package boundaries).
4. **Tests:** Add or update tests for the change. Prefer `package xxx_test` + testify; table-driven when there are ≥2 similar cases.
5. **Checks before opening a PR:**

```bash
make lint
make test
make docs-check    # if you touched docs/en or docs/ru
make generate      # if you changed api/**/*.proto — then go build ./...
make e2e           # if you add/change scripts/e2e-m*.sh or milestone behavior
```

6. **Docs:** User-visible CLI, config, Status, install, or provider behavior → update [docs/en/](./docs/en/) and mirror [docs/ru/](./docs/ru/) ([maintaining](./docs/en/maintaining.md)). New CLI commands → [DESIGN.md §6](./DESIGN.md#6-cli-reference) as well.
7. **Roadmap / milestones:** Prefer work that matches [PROGRESS.md](./PROGRESS.md) / [DESIGN.md §13](./DESIGN.md#13-milestones). Updating PROGRESS (phase / next todo) is for maintainers and milestone owners — not required on every drive-by PR.

### Quick steps to contribute

1. Fork the project.
2. Clone your fork (`git clone https://github.com/YOUR_USERNAME/agentd.git && cd agentd`).
3. Create a feature branch (`git checkout -b my-new-feature`).
4. Make changes; run `make lint` and `make test` (and other checks above as needed).
5. Stage (`git add …` — only files for this change).
6. Commit with a concise message focused on **why**.
7. Push (`git push origin my-new-feature`).
8. Open a pull request against `main` with a short summary and test plan.

### Pull request tips

- Non-trivial PRs: fill the [PR template](./.github/pull_request_template.md) intent note and comprehension checklist.
- Title and body should explain the problem and approach, not only the diff.
- Link related issues.
- CI runs lint, intent-check, docs parity, unit/race tests, integration tags, and e2e — keep the PR green.
- Do not commit secrets, `.env`, or local binaries.

## Development pointers

Requires **Go 1.26+** (see [AGENTS.md](./AGENTS.md)).

| Need | Where |
|------|--------|
| Build / test / generate | [AGENTS.md § Commands](./AGENTS.md#commands) |
| Package layout & hot-path rules | [DESIGN.md](./DESIGN.md) |
| File naming, tests, protobuf | [CONVENTIONS.md](./CONVENTIONS.md) |
| Current milestone checklist | [PROGRESS.md](./PROGRESS.md) |
| Provider quirks | [docs/en/providers.md](./docs/en/providers.md) |

```bash
make build
make start    # local daemon
make stop
```

Never reimplement provider wire codecs outside `internal/hookedge` + [agenthooks](https://github.com/speakeasy-api/agenthooks).
