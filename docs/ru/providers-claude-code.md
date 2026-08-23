# Claude Code

> **Language:** [English](../en/providers-claude-code.md) · [Русский](./providers-claude-code.md)

`--provider=claude-code`. Точка входа: `agentd hook run` (JSON на stdin).

## Установка

```bash
# Проект → .claude/settings.json относительно --dir
agentd install --provider=claude-code --scope=project

# Пользователь (→ ~/.claude/settings.json)
agentd install --provider=claude-code --scope=user

# Плагин → .claude-plugin/plugin.json + hooks/hooks.json
agentd install --provider=claude-code --scope=plugin --dir /path/to/plugin
```

Вид команды: `…/agentd agenthooks run --provider=claude-code`.

## Работа

1. `agentd daemon start`
2. Пользуйтесь Claude Code; установленные хуки порождают CLI.
3. Проверка:

```bash
echo '{"session_id":"s","cwd":"/tmp","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}' \
  | agentd hook run --provider=claude-code
```

## Особенности провайдера

| Тема | Поведение |
|------|-----------|
| **Пустой ответ** | Нейтральное решение → stdout `{}`, код 0 |
| **Ask / Deny** | На ToolPre доступны Deny, Ask, Allow, правка входа, system message, stop-agent (матрица agenthooks) |
| **PromptSubmitted** | Deny / контекст / system message — **без Ask** |
| **Таймауты** | В HookSpec — секунды; бюджет sync берёт defaults install, если у Invoke нет deadline |
| **Блокирующие хуки** | ToolPre / PromptSubmitted / Stop в наборе install по умолчанию — blocking |

См. также: [Провайдеры](./providers.md), [Одобрения](./approvals.md), [Охранники](./guards.md).
