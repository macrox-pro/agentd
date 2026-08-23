# Cursor

> **Language:** [English](../en/providers-cursor.md) · [Русский](./providers-cursor.md)

`--provider=cursor`. Точка входа: `agentd hook run` с **`--argv-payload`** (полезная нагрузка в argv).

## Установка

```bash
# Проект → .cursor/hooks.json
agentd install --provider=cursor --scope=project

# Пользователь → hooks.json в ~/.cursor
agentd install --provider=cursor --global
# то же, что: --scope=user

# Плагин → .cursor-plugin/plugin.json + hooks/hooks.json
agentd install --provider=cursor --scope=plugin --dir /path/to/plugin
```

## Работа

1. `agentd daemon start`
2. Cursor вызывает установленную команду.

```bash
agentd hook run --provider=cursor --argv-payload '<json>'
```

## Особенности провайдера

| Тема | Поведение |
|------|-----------|
| **Ask** | Реально срабатывает только на нативных `beforeShellExecution` / `beforeMCPExecution`. На обычном `preToolUse` Ask игнорируется. `beforeReadFile` — только allow/deny |
| **Install и URL** | Если в **строке команды** хука есть URL, agenthooks **не сгенерирует** конфиг — Cursor иначе выбросил бы весь `hooks.json`. Endpoint держите в конфиге/env, не в argv |
| **Async и sync** | Сбой асинхронной телеметрии **не** должен менять sync-решение ([DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine)) |
| **PromptSubmitted** | На уровне kind — в основном Deny |
| **fail-closed** | Ожидания Cursor по fail-closed учитывает agenthooks при install/runtime |

См. также: [Провайдеры](./providers.md), [Маршрутизация](./dispatch.md), [Диагностика](./troubleshooting.md).
