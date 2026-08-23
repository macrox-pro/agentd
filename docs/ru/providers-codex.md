# OpenAI Codex

> **Language:** [English](../en/providers-codex.md) · [Русский](./providers-codex.md)

`--provider=codex`. Точки входа: `hook run` и `hook notify` (JSON в argv, только асинхронная семантика).

## Установка

```bash
# Пользователь → $CODEX_HOME или ~/.codex/hooks.json
agentd install --provider=codex --scope=user

# Проект → .codex/hooks.json в cwd
agentd install --provider=codex --scope=project
```

Пишутся `hooks.json` и управляемый фрагмент `config.toml` (ключи доверия, без интерактивного trust-цикла).

Неблокирующие хуки могут получить `--async` в командной строке.

## Работа

```bash
agentd daemon start

agentd hook notify --provider=codex '{"type":"agent-turn-complete"}'
```

## Особенности провайдера

| Тема | Поведение |
|------|-----------|
| **Нет Ask** | На ToolPre есть Deny/Allow/…, но **нет CapAsk**. Охранник с `action: ask` уходит в `policy.ask_fallback` (по умолчанию deny) |
| **Пустой ответ** | No-op = **пустой stdout**, код 0 — **не** `{}` |
| **Notify** | Только async; не использовать для блокирующих гейтов |
| **Путь trust** | Ключи доверия содержат **абсолютный** путь к `hooks.json`; смена CODEX_HOME → переустановка |
| **Смысл пустого stdout** | В отличие от Claude, пустой stdout в диалекте Codex — allow/no-op (кодирует agenthooks) |

См. также: [Провайдеры](./providers.md), [Справочник CLI](./cli.md), [Диагностика](./troubleshooting.md).
