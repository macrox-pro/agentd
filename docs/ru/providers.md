# Провайдеры

> **Language:** [English](../en/providers.md) · [Русский](./providers.md)

Как установить хуки agentd в каждого поддерживаемого агента, какой командой он вызывает agentd и **какие особенности** (возможности протокола, формат ответа, пути install).

**Условие:** предпочтительно `agentd daemon start` ([Быстрый старт](./getting-started.md)). Если демон не запущен, установленные хуки применяют `policy.offline` (по умолчанию `fail_open`), чтобы агенты не блокировались.

## Руководства по агентам

| Агент | `--provider` | Точка входа | Страница |
|-------|--------------|-------------|----------|
| Claude Code | `claude-code` | `hook run` (stdin) | [providers-claude-code.md](./providers-claude-code.md) |
| Cursor | `cursor` | `hook run` (`--argv-payload`) | [providers-cursor.md](./providers-cursor.md) |
| OpenAI Codex | `codex` | `hook run` + `hook notify` | [providers-codex.md](./providers-codex.md) |
| Gemini CLI | `gemini` | `hook run` (stdin) | [providers-gemini.md](./providers-gemini.md) |
| OpenCode | `opencode` | `hook serve` (NDJSON) | [providers-opencode.md](./providers-opencode.md) |
| Kimi Code | `kimi-code` / `kimicode` | `hook run` | [providers-kimi.md](./providers-kimi.md) |

## Особенности одним взглядом

| Провайдер | Ask на tool.pre | «Пустой» ответ | Важные ограничения |
|-----------|-----------------|----------------|--------------------|
| Claude | да | `{}` + код 0 | Полный Ask/Deny/Allow на PreToolUse |
| Cursor | только shell/MCP | диалект Cursor | URL в команде ломает install; async не меняет sync |
| Codex | **нет** Ask | **пустой** stdout + код 0 | Trust в `config.toml`; notify только async |
| Gemini | да | свой диалект + дисциплина stderr | Таймауты в **мс**; не писать debug в stderr с hookedge |
| OpenCode | **нет** Ask на tool.pre | кадры serve | Долгий `serve`; mutex сессии в демоне |
| Kimi | **нет** Ask | **пустой** stdout + код 0 | Только user install; в JSON только deny/allow |

Матрица возможностей — из agenthooks; охранники agentd опираются на то, что провайдер умеет выразить (`policy.ask_fallback`, если Ask недоступен).

## Общее поведение install

```bash
agentd install --provider=PROVIDER --scope=SCOPE [--dir PATH]
```

| `--scope` | Смысл |
|-----------|--------|
| `project` (по умолчанию) | Настройки проекта в cwd (codex: `./.codex`) |
| `user` | Home агента (напр. `~/.cursor`, `~/.claude`; codex: `$CODEX_HOME` или `~/.codex`). Алиас: `--global` |
| `plugin` | Корень плагина — `--dir` обязателен (Claude, Cursor) |

При успехе `agentd install` печатает сводку с абсолютными путями для каждого созданного, обновлённого или неизменённого файла.

В argv агента — `agentd agenthooks …`. Таймауты HookSpec: ToolPre / PromptSubmitted **30 s**; короткие виды **5 s**.

Дизайн: [DESIGN.md §1–§2](../../DESIGN.md) · кодеки: [agenthooks](https://github.com/speakeasy-api/agenthooks).
