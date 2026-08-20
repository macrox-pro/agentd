# Gemini CLI

> **Language:** [English](../en/providers-gemini.md) · [Русский](./providers-gemini.md)

`--provider=gemini`. Точка входа: `agentd hook run` (stdin).

## Установка

```bash
# Проект → .gemini/settings.json
agentd install --provider=gemini --scope=project

# Пользователь → settings.json в ~/.gemini
agentd install --provider=gemini --scope=user --dir ~/.gemini
```

## Работа

1. `agentd daemon start`
2. Запускайте Gemini CLI; хуки вызывают `agenthooks run --provider=gemini`.

## Особенности провайдера

| Тема | Поведение |
|------|-----------|
| **Единица таймаута** | В settings Gemini — **миллисекунды**; agenthooks переводит длительности при install |
| **Имена хуков** | Установщик задаёт отображаемые имена для UX `/hooks` |
| **stderr** | Hookedge **не** должен писать отладку в stderr (Gemini может иначе интерпретировать поток). Аудит — в async `file` / `log` ([DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine)) |
| **Ask** | На ToolPre Ask/Deny/Allow поддерживаются |
| **Коды выхода** | Семантика блокировки отличается от Claude; кодирование — зона agenthooks, не выдумывайте коды в agentd |

См. также: [Провайдеры](./providers.md), [Эксплуатация](./operations.md).
