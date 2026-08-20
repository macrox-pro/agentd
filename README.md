# agentd

**A local daemon that proxies, guards, and observes coding-agent hooks — once, for every agent.**

agentd sits between your AI coding agents (Claude Code, Cursor, Codex, Gemini CLI, OpenCode, Kimi Code) and your hook logic. Agents invoke a thin CLI entrypoint; a user-level daemon applies policies, dispatches sync and async pipelines, and returns provider-correct responses. Built on [agenthooks](https://github.com/speakeasy-api/agenthooks) for wire compatibility.

[![Go Reference](https://pkg.go.dev/badge/github.com/macrox-pro/agentd.svg)](https://pkg.go.dev/github.com/macrox-pro/agentd)
![Go Version](https://img.shields.io/badge/go-1.26+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

> **Status:** v1 (M0–M8). Roadmap history: [DESIGN.md §13](./DESIGN.md#13-milestones).

## Documentation

- [User guide (English)](./docs/en/) — default ([why agentd](./docs/en/why.md))
- [Руководство (русский)](./docs/ru/) ([зачем нужен](./docs/ru/why.md))
- Keeping docs current: [docs/en/maintaining.md](./docs/en/maintaining.md)

Contributor design and conventions: [DESIGN.md](./DESIGN.md), [AGENTS.md](./AGENTS.md), [CONVENTIONS.md](./CONVENTIONS.md). How to contribute: [CONTRIBUTING.md](./CONTRIBUTING.md).

## Why agentd?

Coding-agent hooks are powerful but painful to operationalize:

- **Duplicated glue** — each provider speaks a slightly different JSON dialect, timeout unit, and failure mode.
- **Heavy cold starts** — spawning full hook logic on every tool call adds latency.
- **Mixed concerns** — blocking guards, audit webhooks, and metrics want different lifecycles but share one process.

agentd centralizes hook logic in a **long-lived daemon** while keeping the **agent-facing contract** compatible with agenthooks. You configure declarative guards and dispatch routes; the daemon hot-reloads config without re-reading disk on every event.

## Features

- **Universal hook proxy** — one CLI surface (`agentd hook run`) for all supported agents
- **Sync + async + hybrid dispatch** — blocking decisions for the agent, fire-and-forget observability in parallel or after sync
- **Declarative guards** — secrets, shell, MCP, path policies via YAML
- **Approvals & temporary blocks** — Ask once / approve with TTL; runtime overlay persisted across restarts
- **Efficient config reload** — in-memory snapshots, fsnotify with debounce; zero config I/O on the hot path
- **Cross-platform IPC** — gRPC over Unix domain sockets (Linux/macOS) or named pipes (Windows)
- **Provider-faithful I/O** — stdout/stderr discipline and exit codes handled per agenthooks codecs
- **Ops Status** — queue depth and async overflow drop counter on `daemon status`

## Supported agents

| Agent | Hook install target | Entry command | Guide |
|-------|---------------------|---------------|-------|
| Claude Code | `.claude/settings.json`, plugins | `agentd hook run --provider=claude-code` | [docs](./docs/en/providers-claude-code.md) |
| Cursor | `.cursor/hooks.json` | `agentd hook run --provider=cursor` | [docs](./docs/en/providers-cursor.md) |
| OpenAI Codex | `hooks.json` / `config.toml` | `agentd hook run --provider=codex` | [docs](./docs/en/providers-codex.md) |
| Gemini CLI | `.gemini/settings.json` | `agentd hook run --provider=gemini` | [docs](./docs/en/providers-gemini.md) |
| OpenCode | `.opencode/plugin` shim | `agentd hook serve --provider=opencode` | [docs](./docs/en/providers-opencode.md) |
| Kimi Code | user `~/.kimi-code/config.toml` only | `agentd hook run --provider=kimi-code` | [docs](./docs/en/providers-kimi.md) |

Provider quirks (Ask support, empty stdout, timeouts, install scope): [docs/en/providers.md](./docs/en/providers.md).

## Architecture

```
 Agent (Claude/Cursor/…)          agentd CLI              agentd daemon
        │                    (hook edge)              (gRPC + dispatch)
        │  spawn per event          │                         │
        └──── hook run ──────────►│ decode ── Invoke ──────►│ sync pipeline  ──► decision
                                  │                         │ async pipeline ──► queue → sinks
                                  │◄── encode stdout ───────┘
        ◄── JSON + exit code ─────┘
```

- **Hook CLI** — decode/encode only; no business logic
- **Daemon** — routing, guards, forward targets (HTTP, exec, gRPC, logs)
- **Config** — layered YAML with atomic in-memory snapshots

Details: [DESIGN.md](./DESIGN.md)

## Requirements

- Go 1.26+ (to build from source)
- A supported coding agent (see table above)
- Linux, macOS, or Windows

## Installation

```bash
go install github.com/macrox-pro/agentd@latest
```

Pre-built binaries for linux/darwin/windows are published on [GitHub Releases](https://github.com/macrox-pro/agentd/releases) (goreleaser).

Details: [docs/en/installation.md](./docs/en/installation.md).

## Quick start

**1. Start the daemon** (one instance per user):

```bash
agentd daemon start
agentd daemon status
```

**2. Add a minimal config** (`~/.agentd.yaml`):

```yaml
version: 1
policy:
  fail: fail_closed
guards:
  secrets:
    enabled: true
    action: ask
```

**3. Install hooks for your agent** (example: Claude Code, project scope):

```bash
cd your-repo
agentd install --provider=claude-code --scope=project
```

**4. Verify** — trigger a tool call in your agent; check daemon status:

```bash
agentd daemon status --json
```

For OpenCode, use `agentd hook serve --provider=opencode` in generated plugin config (see [DESIGN.md](./DESIGN.md) — OpenCode integration).

Full walkthrough: [docs/en/getting-started.md](./docs/en/getting-started.md).

## Configuration

Configuration merges four layers: defaults → `~/.agentd.yaml` → `.agentd.yaml` (project) → runtime overlay (daemon-written).

Minimal dispatch example (sync guard + async audit):

```yaml
dispatch:
  - name: gate-and-audit
    match: { kind: [tool.pre] }
    mode: parallel
    sync:
      - target: builtin
        guards: [secrets]
    async:
      - target: log
        level: info
```

Full schema, merge rules, and reload behavior: [DESIGN.md § Configuration schema](./DESIGN.md#7-configuration-schema) · [docs/en/configuration.md](./docs/en/configuration.md)

## CLI overview

| Command | Purpose |
|---------|---------|
| `agentd daemon start` | Start the user-level daemon |
| `agentd daemon stop` | Graceful shutdown |
| `agentd daemon status` | Health, config generation, queue depth, async drops |
| `agentd hook run` | **Agent entrypoint** — blocking hooks |
| `agentd hook notify` | Codex notify path (async) |
| `agentd hook serve` | OpenCode NDJSON bridge |
| `agentd install` | Write agent hook configs (via agenthooks) |
| `agentd config validate` | Validate YAML offline (CI-friendly) |
| `agentd config show` | Inspect merged config |
| `agentd config patch` | Patch runtime overlay (persisted) |
| `agentd config record-decision` | Record approval after Ask |
| `agentd dispatch routes` | Show compiled dispatch routes |

Rationale for each command: [DESIGN.md § CLI Reference](./DESIGN.md#6-cli-reference) · [docs/en/cli.md](./docs/en/cli.md)

## Development

```bash
git clone https://github.com/macrox-pro/agentd.git
cd agentd
make generate   # protobuf (requires buf)
make test       # go test -race
make lint       # golangci-lint + buf lint
go test -tags=integration ./...   # optional daemon↔hook integration
```

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [AGENTS.md](./AGENTS.md).

## Contributing

Issues and pull requests are welcome. Please read [CONTRIBUTING.md](./CONTRIBUTING.md) (and [AGENTS.md](./AGENTS.md) before submitting code).

## License

MIT — see [LICENSE](./LICENSE).

## Acknowledgements

Hook wire formats and provider compatibility powered by [speakeasy-api/agenthooks](https://github.com/speakeasy-api/agenthooks).
