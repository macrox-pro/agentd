# Агенты

> **Language:** [English](../en/providers.md) · [Русский](./providers.md)

Как подключить каждого поддерживаемого ИИ-агента, какой командой `agentd hook` он вызывает службу и чего этот агент не умеет.

**Перед установкой хуков:** запустите демон ([Быстрый старт](./getting-started.md)). Если он не запущен, установленные хуки следуют `policy.offline` (по умолчанию `fail_open`) — агент не останавливается.

## Руководства по агентам

| Агент | `--provider` | Как агент вызывает agentd | Страница |
|-------|--------------|---------------------------|----------|
| Claude Code | `claude-code` | `hook run` (JSON во входной поток) | [providers-claude-code.md](./providers-claude-code.md) |
| Cursor | `cursor` | `hook run` (`--argv-payload`) | [providers-cursor.md](./providers-cursor.md) |
| OpenAI Codex | `codex` | `hook run` + `hook notify` | [providers-codex.md](./providers-codex.md) |
| Gemini CLI | `gemini` | `hook run` (JSON во входной поток) | [providers-gemini.md](./providers-gemini.md) |
| OpenCode | `opencode` | `hook serve` (поток NDJSON) | [providers-opencode.md](./providers-opencode.md) |
| Kimi Code | `kimi-code` / `kimicode` | `hook run` | [providers-kimi.md](./providers-kimi.md) |

## Что умеет каждый агент

| Агент | Может спросить вас перед инструментом | Ответ «ничего не менять» | Ограничения |
|-------|----------------------------------------|--------------------------|-------------|
| Claude | да | `{}` и код выхода 0 | Полный набор: спросить / запретить / разрешить на PreToolUse |
| Cursor | только встроенные оболочка и MCP | свой JSON Cursor | Адрес URL в команде хука ломает установку; фон не должен менять синхронный ответ |
| Codex | **нет** | **пустой** stdout и код 0 | Ключи доверия в `config.toml`; `notify` только наблюдение |
| Gemini | да | JSON Gemini; stderr без отладки | Таймауты в **миллисекундах**; с хука не писать отладку в stderr |
| OpenCode | **нет** на tool.pre | кадры `serve` | Долгий `serve`; в демоне одна сессия за раз |
| Kimi | **нет** | **пустой** stdout и код 0 | Только область `user`; в JSON только запрет/разрешение |

Возможности агента задаёт agenthooks. agentd накладывает ваши проверки на эту поверхность (`policy.ask_fallback`, если агент не умеет спрашивать).

## Общая установка хуков

```bash
agentd install --provider=PROVIDER --scope=SCOPE [--dir PATH]
```

| `--scope` | Куда пишутся файлы |
|-----------|-------------------|
| `project` (по умолчанию) | Этот репозиторий (Codex: `./.codex`) |
| `user` | Домашний каталог агента (например `~/.cursor`, `~/.claude`; Codex: `$CODEX_HOME` или `~/.codex`). То же, что `--global` |
| `plugin` | Корень расширения — нужен `--dir` (Claude, Cursor) |

При успехе `agentd install` печатает каждый созданный, обновлённый или неизменённый файл с абсолютным путём.

В настройках агента вызывается `agentd agenthooks …` (то же, что `hook …`). Таймауты: «перед инструментом» и отправка запроса **30 с**; более короткие события **5 с**.

Устройство: [DESIGN.md §1–§2](../../DESIGN.md) · форматы: [agenthooks](https://github.com/speakeasy-api/agenthooks).

## Автообнаружение

`agentd doctor` и `agentd install --all-detected` смотрят рабочий каталог, домашний каталог и `PATH`.

| Уверенность | Признак | `--all-detected --yes` | `doctor` |
|-------------|---------|------------------------|----------|
| Высокая | Есть каталог настроек агента | Записывает хуки | Показывает |
| Средняя | Есть только программа в `PATH` | **Пропускает** | Показывает с пометкой |

| Агент | Каталог в проекте | Каталог пользователя | Автоустановка |
|-------|-------------------|----------------------|---------------|
| `claude-code` | `.claude/` | `~/.claude/` | проект и пользователь, если есть оба |
| `cursor` | `.cursor/` | `~/.cursor/` | проект и пользователь, если есть оба |
| `codex` | `.codex/` | `$CODEX_HOME` или `~/.codex` | проект и пользователь, если есть оба |
| `gemini` | `.gemini/` | `~/.gemini/` | проект и пользователь, если есть оба |
| `opencode` | `.opencode/` | — | только проект (`user` нужен `--dir`) |
| `kimi-code` | — | `$KIMI_CODE_HOME` или `~/.kimi-code` | только пользователь |

Область `plugin` сама не выбирается. Явно: `agentd install --provider=… --scope=plugin --dir=…`.
