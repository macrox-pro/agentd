# Kimi Code

> **Language:** [English](../en/providers-kimi.md) · [Русский](./providers-kimi.md)

`--provider=kimi-code` (лучше для encode) или `kimicode` (принимает parse в agentd). Точка входа: `agentd hook run`.

## Установка

Хуки только в **пользовательском** `config.toml` (`$KIMI_CODE_HOME`, обычно `~/.kimi-code`). **`--scope=project` завершится ошибкой**.

```bash
agentd install --provider=kimi-code --scope=user --dir "${KIMI_CODE_HOME:-$HOME/.kimi-code}"
```

Обновляет управляемый регион `[[hooks]]` в `config.toml`.

## Работа

1. `agentd daemon start`
2. Работайте в Kimi; хуки вызывают `agenthooks run --provider=kimi-code`.

На CLI предпочитайте **`kimi-code`**.

## Особенности провайдера

| Тема | Поведение |
|------|-----------|
| **Нет Ask** | ToolPre — только Deny/Allow; нет Ask, update-input и additionalContext в JSON |
| **Пустой ответ** | No-op = **пустой stdout**, код 0 (как у Codex) |
| **Блокировка prompt/stop** | Блокирующий исход — код **2** + stderr (quirk agenthooks); не ждите Claude-`{}` |
| **Только наблюдение** | Многие виды событий (PostToolUse, PermissionRequest, …) — observe-only |
| **Область install** | Только user; project-level hooks у Kimi нет |

См. также: [Провайдеры](./providers.md), [Диагностика](./troubleshooting.md), [Охранники](./guards.md).
