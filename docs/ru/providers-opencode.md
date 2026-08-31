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
2. OpenCode загружает плагин; shim держит дочерний `serve`.
3. Каждый кадр NDJSON → gRPC `Invoke`.

```bash
agentd hook serve --provider=opencode
```

`hook serve` принимает только `--provider=opencode`.

## Особенности агента

| Тема | Поведение |
|------|-----------|
| **Модель процесса** | Один долгий процесс serve; события — кадры NDJSON |
| **Порядок сессий** | В демоне — **mutex на сессию** для sync ([DESIGN.md §1](../../DESIGN.md)) |
| **Нет Ask на tool.pre** | Возможности ToolPre: Deny + update-input, без Ask. Для «спросить» опирайтесь на Deny / `ask_fallback` |
| **Stop / session.idle** | Не на каждом stop есть Continue — часть событий observe-only |
| **Permission** | Идут через канал permission OpenCode (allow/deny), не через Claude-style Ask JSON |

См. также: [Агенты](./providers.md), [Справочник команд](./cli.md).
