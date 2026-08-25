# Providers

> **Language:** [English](./providers.md) · [Русский](../ru/providers.md)

How to install agentd hooks into each supported coding agent, which CLI entrypoint it uses, and **provider-specific quirks** (capabilities, wire, install paths).

**Prerequisite:** Prefer `agentd daemon start` ([Getting started](./getting-started.md)). When the daemon is down, installed hooks apply `policy.offline` (default `fail_open`) so agents are not blocked.

## Agent guides

| Agent | `--provider` | Entrypoint | Guide |
|-------|--------------|------------|--------|
| Claude Code | `claude-code` | `hook run` (stdin) | [providers-claude-code.md](./providers-claude-code.md) |
| Cursor | `cursor` | `hook run` (`--argv-payload`) | [providers-cursor.md](./providers-cursor.md) |
| OpenAI Codex | `codex` | `hook run` + `hook notify` | [providers-codex.md](./providers-codex.md) |
| Gemini CLI | `gemini` | `hook run` (stdin) | [providers-gemini.md](./providers-gemini.md) |
| OpenCode | `opencode` | `hook serve` (NDJSON) | [providers-opencode.md](./providers-opencode.md) |
| Kimi Code | `kimi-code` / `kimicode` | `hook run` | [providers-kimi.md](./providers-kimi.md) |

## Quirks at a glance

| Provider | Ask on tool.pre | Neutral / no-op wire | Notable constraints |
|----------|-----------------|----------------------|---------------------|
| Claude | yes | `{}` + exit 0 | Full Ask/Deny/Allow surface on PreToolUse |
| Cursor | only shell/MCP natives | provider dialect | URL in command breaks install; async must not flip sync |
| Codex | **no** CapAsk | **empty** stdout + exit 0 | Trust keys in `config.toml`; notify is async-only |
| Gemini | yes | dialect + stderr discipline | Timeouts in **ms**; never debug on stderr from hookedge |
| OpenCode | **no** CapAsk on tool.pre | serve frames | Long-lived `serve`; per-session mutex in daemon |
| Kimi | **no** CapAsk | **empty** stdout + exit 0 | User-scope install only; deny/allow JSON only |

Capabilities come from agenthooks’ matrix; agentd guards honor what the provider can express (`policy.ask_fallback` when Ask is unsupported).

## Shared install behavior

```bash
agentd install --provider=PROVIDER --scope=SCOPE [--dir PATH]
```

| `--scope` | Meaning |
|-----------|---------|
| `project` (default) | Project-local settings under CWD (codex: `./.codex`) |
| `user` | Agent home directory (e.g. `~/.cursor`, `~/.claude`; codex: `$CODEX_HOME` or `~/.codex`). Alias: `--global` |
| `plugin` | Plugin root — `--dir` required (Claude, Cursor) |

On success, `agentd install` prints a summary with absolute paths for each created, updated, or unchanged file.

Generated argv uses `agentd agenthooks …` (same as `hook …`). HookSpec timeouts: ToolPre / PromptSubmitted **30s**; shorter kinds **5s**.

Deep design: [DESIGN.md §1–§2](../../DESIGN.md) · codecs: [agenthooks](https://github.com/speakeasy-api/agenthooks).
