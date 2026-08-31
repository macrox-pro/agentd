<div align="center">
<img width="1280" height="640" alt="agentd-about" src="https://github.com/user-attachments/assets/399f29b3-958b-43ce-8166-a548bc7dde19" />
</div>

<hr />

<p align="center">
  <h1 align="center"><b>agentd</b></h1>
  <p align="center">A local daemon that proxies, guards, and observes coding-agent hooks — once, for every agent.</p>
</p>

<hr />

agentd sits between your AI coding agents (Claude Code, Cursor, Codex, Gemini CLI, OpenCode, Kimi Code) and your hook logic. Agents invoke a thin CLI entrypoint; a user-level daemon applies policies, dispatches sync and async pipelines, and returns provider-correct responses. Built on [agenthooks](https://github.com/speakeasy-api/agenthooks) for wire compatibility.

[![Go Reference](https://pkg.go.dev/badge/github.com/macrox-pro/agentd.svg)](https://pkg.go.dev/github.com/macrox-pro/agentd)
![Go Version](https://img.shields.io/badge/go-1.26+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

> **Status:** [v0.0.7](./CHANGELOG.md#v007--2026-08-31) (setup wizard, doctor, trajectory defaults on). Roadmap history: [PROGRESS.md](./PROGRESS.md).

## Documentation

- [User guide (English)](./docs/en/) — default ([why agentd](./docs/en/why.md) · [getting started](./docs/en/getting-started.md))
- [Руководство (русский)](./docs/ru/) ([зачем нужен](./docs/ru/why.md) · [быстрый старт](./docs/ru/getting-started.md))
- Keeping docs current: [docs/en/maintaining.md](./docs/en/maintaining.md)

Contributor design and conventions: [DESIGN.md](./DESIGN.md), [AGENTS.md](./AGENTS.md), [CONVENTIONS.md](./CONVENTIONS.md). How to contribute: [CONTRIBUTING.md](./CONTRIBUTING.md).

## Research

Structured verbatim excerpts from primary sources — used when designing provider support, hooks, and conventions. Each tree has its own index and `SOURCES.md`.

| Tree | Focus |
|------|--------|
| [research/best-practice](./research/best-practice/) | Go best practices (Go ≥ 1.26.7) |
| [research/claude-code](./research/claude-code/) | Claude Code docs — agent loop, hooks, MCP, settings, skills, plugins, cloud, enterprise |
| [research/codex](./research/codex/) | Codex / ChatGPT Learn docs — sandbox, hooks, MCP, cloud, enterprise |
| [research/cursor](./research/cursor/) | Cursor docs — agent loop, hooks, MCP, settings, enterprise |
| [research/gemini](./research/gemini/) | Gemini CLI docs (snapshot 2026-08-29) — agent loop, hooks, MCP, settings, skills, extensions, Managed Agents API, Antigravity migration delta |

Stub dirs for other agents (`opencode`, `kimi-code`) live under [`research/`](./research/) and will fill in the same shape.

## Why agentd?

Coding-agent hooks are powerful but painful to operationalize:

- **Duplicated glue** — each provider speaks a slightly different JSON dialect, timeout unit, and failure mode.
- **Heavy cold starts** — spawning full hook logic on every tool call adds latency.
- **Mixed concerns** — blocking guards, audit webhooks, and metrics want different lifecycles but share one process.

agentd centralizes hook logic in a **long-lived daemon** while keeping the **agent-facing contract** compatible with agenthooks. You configure declarative guards and dispatch routes; the daemon hot-reloads config without re-reading disk on every event.

## Features

- **Universal hook proxy** — one CLI surface (`agentd hook run` / `notify` / `serve`) for all supported agents
- **Sync + async + hybrid dispatch** — blocking decisions for the agent, fire-and-forget observability in parallel or after sync
- **Declarative guards** — secrets, shell, MCP, path policies via YAML
- **Approvals & temporary blocks** — Ask once / approve with TTL; runtime overlay persisted across restarts
- **Session ledger** — hook calls logged by default (`agentd config get trajectory`); disable with `config disable trajectory`
- **Detect and install** — `agentd doctor` (read-only); `install --all-detected` (plan unless `--yes`); `setup` in a terminal
- **Efficient config reload** — in-memory snapshots, fsnotify with debounce; zero config I/O on the hot path
- **Cross-platform IPC** — gRPC over Unix domain sockets (Linux/macOS) or named pipes (Windows)
- **Provider-faithful I/O** — stdout/stderr discipline and exit codes handled per agenthooks codecs
- **Ops Status** — queue depth, async drops, optional Prometheus listen address on `daemon status --json`

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

**1. Start the daemon** (one instance per user). If `~/.agentd.yaml` is missing, start creates a [minimal bootstrap](./docs/en/configuration.md#user-config-bootstrap) automatically:

```bash
agentd daemon start
agentd daemon status
```

**2. Customize user config** (optional — edit `~/.agentd.yaml` after start, or create it yourself first):

```yaml
version: 1
policy:
  fail: fail_closed
  # offline defaults to fail_open — agents keep working if the daemon is down
guards:
  secrets:
    enabled: true
    action: ask
```

**3. Connect an agent.** See what agentd finds (read-only), then install:

```bash
agentd doctor
agentd install --all-detected          # plan only
agentd install --all-detected --yes    # write hook files
```

Or one agent, this repository:

```bash
cd your-repo
agentd install --provider=claude-code --scope=project
```

In a terminal, `agentd setup` walks the same flow. CI: `AGENTD_NO_TUI=1`.

**4. Verify** — trigger a tool call in your agent; check daemon status:

```bash
agentd daemon status --json
agentd config get trajectory    # trajectory: on (default)
```

For OpenCode, generated plugin config uses `agentd hook serve --provider=opencode` (see [DESIGN.md §1](./DESIGN.md#1-architecture)).

Full walkthrough: [docs/en/getting-started.md](./docs/en/getting-started.md).

## Configuration

Configuration merges four layers: defaults → `~/.agentd.yaml` → `.agentd.yaml` (project) → runtime overlay (daemon-written). State (log, runtime overlay, sessions) lives under the [state directory](./docs/en/configuration.md#state-directory), not `~/.agentd/`.

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

Full schema: [docs/en/configuration.md](./docs/en/configuration.md) · layer/runtime overlay: [DESIGN.md §7](./DESIGN.md#7-configuration-schema)

## CLI overview

| Command | Purpose |
|---------|---------|
| `agentd daemon start` | Start the user-level daemon |
| `agentd daemon enable` | Register login autostart (see [Operations](./docs/en/operations.md#autostart-at-login)) |
| `agentd daemon disable` | Remove login autostart |
| `agentd daemon stop` | Graceful shutdown |
| `agentd daemon status` | Health, config generation, queue depth, async drops |
| `agentd hook run` | **Agent entrypoint** — blocking hooks |
| `agentd hook notify` | Codex notify path (async) |
| `agentd hook serve` | OpenCode NDJSON bridge |
| `agentd doctor` | Detect agents and hook install status (read-only) |
| `agentd install` | Write agent hook configs (`--provider` or `--all-detected`) |
| `agentd setup` | Interactive install (terminal) |
| `agentd config validate` | Validate YAML offline (CI-friendly) |
| `agentd config enable FEATURE` | Curated toggles (trajectory, guards) — user/project YAML |
| `agentd config disable FEATURE` | Turn off a curated toggle |
| `agentd config get FEATURE` | Effective on/off + winning layer (no runtime) |
| `agentd config show` | Inspect merged config |
| `agentd config patch` | Patch runtime overlay (persisted) |
| `agentd config record-decision` | Record approval after Ask |
| `agentd dispatch routes` | Show compiled dispatch routes |
| `agentd session list` / `show` / `export` | Inspect the session ledger (offline) |
| `agentd session subscribe` | Live trajectory stream (daemon required) |
| `agentd trajectory stats` | Daemon-lifetime counters (daemon required) |

**Trajectory (default on):** every supported agent’s hooks are traceable on one stream; transcript/thinking depth varies by provider — not “everything the model sees everywhere.”

| Provider | L2 import |
|----------|-----------|
| claude-code, codex | **supported** |
| cursor | **partial** (`--path`) |
| gemini, opencode, kimi-code | none |

Details: [docs/en/trajectory.md](./docs/en/trajectory.md) · [DESIGN §14.3](./DESIGN.md#143-provider-support-matrix)

Rationale for each command: [docs/en/cli.md](./docs/en/cli.md) · [DESIGN §6](./DESIGN.md#6-cli-reference) (architecture notes)

## Development

```bash
git clone https://github.com/macrox-pro/agentd.git
cd agentd
make generate   # protobuf (requires buf)
make test       # go test -race
make lint       # golangci-lint + buf lint
make e2e
go test -tags=integration ./...   # optional daemon↔hook integration
```

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [AGENTS.md](./AGENTS.md).

## Contributing

Issues and pull requests are welcome. Please read [CONTRIBUTING.md](./CONTRIBUTING.md) (and [AGENTS.md](./AGENTS.md) before submitting code).

## License

MIT — see [LICENSE](./LICENSE).

## Acknowledgements

Hook wire formats and provider compatibility powered by [speakeasy-api/agenthooks](https://github.com/speakeasy-api/agenthooks).
