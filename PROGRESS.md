# agentd — implementation progress

> Source of truth for multi-session / context-compressed agents.
> Full design: [DESIGN.md](./DESIGN.md) | Contributor rules: [AGENTS.md](./AGENTS.md)

## Current phase

Phase: **complete (doc + scaffold)** | Last todo: p4-checkpoint | Next: **M1 implementation**

## agents_md_ready

true

## Todo checklist

### Phase 0
- [x] p0-progress-file

### Phase 1 — Documentation
- [x] p1-readme
- [x] p1-design-skeleton
- [x] p1-design-dispatch
- [x] p1-design-configstore
- [x] p1-design-grpc
- [x] p1-design-transport
- [x] p1-design-cli
- [x] p1-design-config-schema
- [x] p1-design-ops
- [x] p1-agents-md
- [x] p1-agents-gate
- [x] p1-doc-crosslinks
- [x] p1-checkpoint

### Phase 2 — Protobuf
- [x] p2-buf-config
- [x] p2-proto-common
- [x] p2-proto-hook
- [x] p2-proto-daemon
- [x] p2-proto-config
- [x] p2-buf-lint-breaking (lint clean; breaking skipped — no baseline on main)
- [x] p2-buf-generate
- [x] p2-design-proto-sync
- [x] p2-checkpoint

### Phase 3 — Skeleton
- [x] p3-deps
- [x] p3-internal-layout
- [x] p3-transport-stub
- [x] p3-config-stub
- [x] p3-dispatch-stub
- [x] p3-cli-root
- [x] p3-cli-daemon
- [x] p3-cli-hook
- [x] p3-cli-config
- [x] p3-cli-install
- [x] p3-cli-dispatch
- [x] p3-checkpoint

### Phase 4 — Verify
- [x] p4-makefile
- [x] p4-final-review
- [x] p4-checkpoint

## Files touched (scaffold phase)

- PROGRESS.md, README.md, DESIGN.md, AGENTS.md
- api/agentd/v1/*.proto, gen/agentd/v1/*.pb.go
- buf.gen.yaml, buf.lock, go.mod, go.sum, Makefile, .gitignore
- cmd/*.go, internal/**

## Verify commands (last green)

```bash
go build ./...
./agentd --help
buf lint api/
make lint
make test
```

## Blockers / notes

(none)

## Deferred (M1–M4)

See [DESIGN.md § Milestones](./DESIGN.md#13-milestones):

- **M1** — daemon start/stop/status; ConfigStore wiring; `Invoke` stub
- **M2** — Dispatch Engine; hybrid modes; async queue; secrets guard
- **M3** — forward targets; full dispatch YAML; fsnotify reload
- **M4** — gRPC forward; OpenCode serve bridge; install wrapper; Windows npipe hardening
