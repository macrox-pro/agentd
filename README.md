# agentd

**A local daemon that proxies, guards, and observes coding-agent hooks — once, for every agent.**

agentd sits between your AI coding agents (Claude Code, Cursor, Codex, Gemini CLI, OpenCode, Kimi Code) and your hook logic. Agents invoke a thin CLI entrypoint; a user-level daemon applies policies, dispatches sync and async pipelines, and returns provider-correct responses. Built on [agenthooks](https://github.com/speakeasy-api/agenthooks) for wire compatibility.

[![Go Reference](https://pkg.go.dev/badge/github.com/macrox-pro/agentd.svg)](https://pkg.go.dev/github.com/macrox-pro/agentd)
![Go Version](https://img.shields.io/badge/go-1.26+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

> **Status:** Early development (M0–M4 done; M5→v1 planned). API and config schema may change before v1. Roadmap: [DESIGN.md §13](./DESIGN.md#13-milestones).

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
- **Runtime config overlay** — daemon-managed approvals and temporary blocks (without editing your repo config)
- **Efficient config reload** — in-memory snapshots, fsnotify with debounce; zero config I/O on the hot path
- **Cross-platform IPC** — gRPC over Unix domain sockets (Linux/macOS) or named pipes (Windows)
- **Provider-faithful I/O** — stdout/stderr discipline and exit codes handled per agenthooks codecs

## Supported agents

| Agent | Hook install target | Entry command |
|-------|---------------------|---------------|
| Claude Code | `.claude/settings.json`, plugins | `agentd hook run --provider=claude-code` |
| Cursor | `hooks.json` | `agentd hook run --provider=cursor` |
| OpenAI Codex | `hooks.json` / `config.toml` | `agentd hook run --provider=codex` |
| Gemini CLI | `settings.json` | `agentd hook run --provider=gemini` |
| OpenCode | plugin shim (stdio) | `agentd hook serve --provider=opencode` |
| Kimi Code | project settings | `agentd hook run --provider=kimicode` |

Provider quirks (timeouts, fail-open vs fail-closed, MCP naming) are encoded once — via agenthooks — not in your hook scripts.

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

Pre-built binaries will be published with GitHub Releases (planned).

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

Full schema, merge rules, and reload behavior: [DESIGN.md § Configuration schema](./DESIGN.md#7-configuration-schema)

## CLI overview

| Command | Purpose |
|---------|---------|
| `agentd daemon start` | Start the user-level daemon |
| `agentd daemon stop` | Graceful shutdown |
| `agentd daemon status` | Health, config generation, queue depth |
| `agentd hook run` | **Agent entrypoint** — blocking hooks |
| `agentd hook notify` | Codex notify path (async) |
| `agentd hook serve` | OpenCode NDJSON bridge |
| `agentd install` | Write agent hook configs (via agenthooks) |
| `agentd config validate` | Validate YAML offline (CI-friendly) |
| `agentd config show` | Inspect merged config |
| `agentd dispatch routes` | Show compiled dispatch routes |

Rationale for each command: [DESIGN.md § CLI Reference](./DESIGN.md#6-cli-reference)

## Development

```bash
git clone https://github.com/macrox-pro/agentd.git
cd agentd
make generate   # protobuf (requires buf)
make test       # go test -race
make lint       # golangci-lint + buf lint
```

Contributor conventions: [AGENTS.md](./AGENTS.md)

## Contributing

Issues and pull requests are welcome. Please read [AGENTS.md](./AGENTS.md) before submitting code changes.

## License

MIT — see [LICENSE](./LICENSE).

## Acknowledgements

Hook wire formats and provider compatibility powered by [speakeasy-api/agenthooks](https://github.com/speakeasy-api/agenthooks).
