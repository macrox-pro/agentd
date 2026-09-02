# OpenCode

> **Language:** [English](../en/providers-opencode.md) · [Русский](./providers-opencode.md)

`--provider=opencode`. Точка входа: долгоживущий **`agentd hook serve`** (NDJSON), не `hook run` на каждое событие.

## Установка

```bash
agentd install --provider=opencode --scope=project --dir /path/to/repo
```

Пишет `.opencode/plugin/agenthooks.ts`, порождающий:

`agentd agenthooks serve --provider=opencode`

(= `agentd hook serve --provider=opencode`).

## Работа

1. `agentd daemon start`
2. OpenCode загружает плагин; обёртка плагина держит дочерний процесс `serve`.
3. Каждая строка NDJSON на stdin/stdout — один обмен событием с демоном.

```bash
agentd hook serve --provider=opencode
```

`hook serve` принимает только `--provider=opencode`.

При работающем демоне рабочий каталог **каждого** NDJSON-кадра определяет, какой проектный `.agentd.yaml` подмешать ([Конфигурация → Рабочий каталог события](./configuration.md#рабочий-каталог-события-и-проектный-конфиг)). Если в кадрах нет `cwd` и `workspace_roots`, демон может использовать каталог процесса `serve` — задайте `cwd` в JSON кадра, чтобы политика соответствовала репозиторию.

## Особенности агента

| Тема | Поведение |
|------|-----------|
| **Модель процесса** | Один долгий процесс serve; события — кадры NDJSON |
| **Проектный конфиг по кадру** | Каждая строка NDJSON — отдельный вызов демона. Проектный `.agentd.yaml` выбирается по `cwd` или `workspace_roots[0]` **в этом кадре**. Поле `input.directory` OpenCode **не** используется — обёртка из установки должна проставлять top-level `cwd` на каждый кадр |
| **Порядок сессий** | В демоне — **mutex на сессию** для sync ([DESIGN.md §1](../../DESIGN.md)) |
| **Нет Ask на tool.pre** | Возможности ToolPre: Deny + update-input, без Ask. Для «спросить» опирайтесь на Deny / `ask_fallback` |
| **Stop / session.idle** | Не на каждом stop есть Continue — часть событий observe-only |
| **Permission** | Идут через канал permission OpenCode (allow/deny), не через Claude-style Ask JSON |

См. также: [Агенты](./providers.md), [Справочник команд](./cli.md).
