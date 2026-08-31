# Providers

> **Language:** [English](./providers.md) · [Русский](../ru/providers.md)

How to connect each supported coding agent, which `agentd hook` command it uses, and what that agent cannot do.

**Before you install:** start the daemon ([Getting started](./getting-started.md)). If it is down, installed hooks follow `policy.offline` (default `fail_open`) so the agent is not blocked.

## Agent guides

| Agent | `--provider` | How the agent calls agentd | Guide |
|-------|--------------|----------------------------|--------|
| Claude Code | `claude-code` | `hook run` (JSON on stdin) | [providers-claude-code.md](./providers-claude-code.md) |
| Cursor | `cursor` | `hook run` (`--argv-payload`) | [providers-cursor.md](./providers-cursor.md) |
| OpenAI Codex | `codex` | `hook run` + `hook notify` | [providers-codex.md](./providers-codex.md) |
| Gemini CLI | `gemini` | `hook run` (JSON on stdin) | [providers-gemini.md](./providers-gemini.md) |
| OpenCode | `opencode` | `hook serve` (NDJSON stream) | [providers-opencode.md](./providers-opencode.md) |
| Kimi Code | `kimi-code` / `kimicode` | `hook run` | [providers-kimi.md](./providers-kimi.md) |

## What each agent can express

| Agent | Can ask the user before a tool | “Do nothing” reply | Limits |
|-------|--------------------------------|--------------------|--------|
| Claude | yes | `{}` and exit 0 | Full ask / deny / allow on PreToolUse |
| Cursor | only native shell and MCP | Cursor’s own JSON | A URL in the hook command breaks install; async must not change the sync reply |
| Codex | **no** | **empty** stdout and exit 0 | Trust keys in `config.toml`; `notify` is observation only |
| Gemini | yes | Gemini JSON; keep stderr clean | Timeouts in **milliseconds**; never print debug on stderr from the hook |
| OpenCode | **no** on tool.pre | serve frames | Long-lived `serve`; one session at a time in the daemon |
| Kimi | **no** | **empty** stdout and exit 0 | User-scope install only; JSON is deny/allow only |

What an agent can say comes from agenthooks. agentd maps your guards onto that surface (`policy.ask_fallback` when the agent cannot ask).

## Shared install

```bash
agentd install --provider=PROVIDER --scope=SCOPE [--dir PATH]
```

| `--scope` | Where files are written |
|-----------|-------------------------|
| `project` (default) | This repository (Codex: `./.codex`) |
| `user` | Agent home (for example `~/.cursor`, `~/.claude`; Codex: `$CODEX_HOME` or `~/.codex`). Same as `--global` |
| `plugin` | Plugin root — `--dir` required (Claude, Cursor) |

On success, `agentd install` prints each created, updated, or unchanged file with an absolute path.

Generated commands use `agentd agenthooks …` (same as `hook …`). Timeouts: before-tool and prompt-submit **30s**; shorter events **5s**.

Design: [DESIGN.md §1–§2](../../DESIGN.md) · codecs: [agenthooks](https://github.com/speakeasy-api/agenthooks).

## Auto-detection

`agentd doctor` and `agentd install --all-detected` look at the working directory, your home directory, and `PATH`.

| Confidence | Signal | `--all-detected --yes` | `doctor` |
|------------|--------|------------------------|----------|
| High | Agent config folder exists | Writes hooks | Shown |
| Medium | Binary on `PATH` only | **Skipped** | Shown with a note |

| Agent | Folder in the project | Folder for this user | Auto-install |
|-------|----------------------|----------------------|--------------|
| `claude-code` | `.claude/` | `~/.claude/` | project and user if both exist |
| `cursor` | `.cursor/` | `~/.cursor/` | project and user if both exist |
| `codex` | `.codex/` | `$CODEX_HOME` or `~/.codex` | project and user if both exist |
| `gemini` | `.gemini/` | `~/.gemini/` | project and user if both exist |
| `opencode` | `.opencode/` | — | project only (`user` needs `--dir`) |
| `kimi-code` | — | `$KIMI_CODE_HOME` or `~/.kimi-code` | user only |

Plugin scope is never chosen automatically. Use `agentd install --provider=… --scope=plugin --dir=…`.
