# Providers

> **Language:** [English](./providers.md) · [Русский](../ru/providers.md)

Install targets and hook entrypoints per coding agent.

## Supported `--provider` values

| Provider flag | Typical agent | Entrypoint |
|---------------|---------------|------------|
| `claude-code` | Claude Code | `hook run` |
| `cursor` | Cursor | `hook run` (often `--argv-payload`) |
| `codex` | OpenAI Codex | `hook run` / `hook notify` |
| `gemini` | Gemini CLI | `hook run` |
| `opencode` | OpenCode | `hook serve` |
| `kimicode` / `kimi-code` | Kimi Code | `hook run` |

Install accepts the same provider strings. Encode path for Kimi should use `kimi-code` (agenthooks name); `kimicode` is accepted by agentd parse.

## install

```bash
agentd install --provider=claude-code --scope=project --dir /path/to/repo
```

| `--scope` | Meaning |
|-----------|---------|
| `project` (default) | Project hook settings |
| `user` | User-level settings |
| `plugin` | Plugin install target |

Writes provider configs via agenthooks; `Command` is the absolute `agentd` binary. Generated argv uses the hidden `agentd agenthooks …` sentinel (same as `hook …`).

Install HookSpec timeouts: ToolPre / PromptSubmitted **30s**; shorter hooks **5s** — used as sync budget defaults when Invoke has no deadline ([Dispatch](./dispatch.md)).

## OpenCode

Long-lived NDJSON:

```bash
agentd hook serve --provider=opencode
```

Install may generate `agentd agenthooks serve --provider=opencode`.

## Codex notify

```bash
agentd hook notify --provider=codex '{"type":"agent-turn-complete"}'
```

Async semantics; empty stdout is a valid no-op for Codex.

Wire codecs: [agenthooks](https://github.com/speakeasy-api/agenthooks). Architecture: [DESIGN.md §1](../../DESIGN.md).
