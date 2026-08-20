# Providers

> **Language:** [English](../en/providers.md) · [Русский](./providers.md)

Цели install и entrypoint hooks по coding agent.

## Поддерживаемые `--provider`

| Флаг | Агент | Entrypoint |
|------|-------|------------|
| `claude-code` | Claude Code | `hook run` |
| `cursor` | Cursor | `hook run` (часто `--argv-payload`) |
| `codex` | OpenAI Codex | `hook run` / `hook notify` |
| `gemini` | Gemini CLI | `hook run` |
| `opencode` | OpenCode | `hook serve` |
| `kimicode` / `kimi-code` | Kimi Code | `hook run` |

Install принимает те же строки. Для encode у Kimi предпочтителен `kimi-code` (имя agenthooks); `kimicode` принимает parse в agentd.

## install

```bash
agentd install --provider=claude-code --scope=project --dir /path/to/repo
```

| `--scope` | Смысл |
|-----------|--------|
| `project` (default) | Project hook settings |
| `user` | User-level |
| `plugin` | Plugin target |

Пишет конфиги через agenthooks; `Command` — абсолютный путь к `agentd`. Сгенерированный argv использует hidden sentinel `agentd agenthooks …` (эквивалент `hook …`).

Таймауты HookSpec при install: ToolPre / PromptSubmitted **30s**; короткие hooks **5s** — defaults sync budget без deadline у Invoke ([Dispatch](./dispatch.md)).

## OpenCode

Долгоживущий NDJSON:

```bash
agentd hook serve --provider=opencode
```

Install может сгенерировать `agentd agenthooks serve --provider=opencode`.

## Codex notify

```bash
agentd hook notify --provider=codex '{"type":"agent-turn-complete"}'
```

Async-семантика; пустой stdout — валидный no-op для Codex.

Codecs: [agenthooks](https://github.com/speakeasy-api/agenthooks). Архитектура: [DESIGN.md §1](../../DESIGN.md).
