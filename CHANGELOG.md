# Changelog

## [v1.0.0] — 2026-08-20

First v1 release (M0–M8).

### Highlights

- Four-layer config merge (defaults ⊕ user ⊕ project ⊕ runtime) with ConfigService Get/Patch/RecordDecision
- Guards: secrets, shell, MCP, paths; approvals + temporary blocks with runtime persist
- Sync + async dispatch (builtin, exec async-only, http, grpc, log, file)
- Install + `hook run` / `notify` / `serve` for supported providers
- Cross-platform IPC (unix sockets + Windows SID named pipes)
- Status exposes async queue depth and overflow drop counter
- Provider-aware sync timeout margin (`min(provider_timeout - margin, route.sync_timeout)`)
- Conformance fixtures via agenthookstest; optional `//go:build integration` round-trip
- GitHub Releases binaries via goreleaser (linux/darwin/windows)

### Explicitly not in v1

Agent auth, transcripts, plugins, hooks DSL, async retry storms, exec sync JSON decisions.
