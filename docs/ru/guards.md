# Проверки

> **Language:** [English](../en/guards.md) · [Русский](./guards.md)

Декларативные проверки **до** ответа агенту (обычно `tool.pre`). Настройки — в `guards:`. Термины: [Глоссарий](./glossary.md).

## Имена

| Проверка | `enabled` по умолчанию | Основные поля |
|----------|------------------------|---------------|
| `secrets` | `true` | `action`: `ask` \| `deny` (по умолчанию `ask`); опционально `rules` |
| `shell` | `false` | `deny_patterns`, `ask_on` |
| `mcp` | `false` | `deny_servers` |
| `paths` | `false` | `deny_read`, `deny_write` |

У синхронной цели `builtin` список `guards: [secrets, shell]` запускает только эти проверки. Если список **не указан** — все включённые проверки.

`ask_on` (проверка оболочки): имена инструментов, при которых срабатывает Ask, а не молчаливое разрешение — например `curl`, `wget`, `ssh`.

## Включение через CLI

Проверки shell, MCP и paths по умолчанию в области **project** — из корня репозитория:

```bash
cd /path/to/repo
agentd config enable guard-shell
agentd config get guard-shell
```

Команды: `guard-shell`, `guard-mcp`, `guard-paths` ([CLI](./cli.md#config)). **`secrets` не переключается через CLI** — включён в автосозданном YAML; меняйте `action` / `rules` в конфиге.

## Ask и Deny

- **Deny** — жёсткий отказ.
- **Ask** — запрос разрешения; в тексте может быть `approval_fingerprint=sha256:…` для [одобрения](./approvals.md).

Не все агенты умеют Ask. Если Ask недоступен, действует `policy.ask_fallback`: `deny` (по умолчанию) или `no_decision` ([Конфигурация](./configuration.md#policy)).

## Пример

```yaml
guards:
  secrets:
    enabled: true
    action: ask
  shell:
    enabled: true
    deny_patterns: ["rm -rf /"]
    ask_on: [curl, wget, ssh]
  mcp:
    enabled: true
    deny_servers: ["untrusted-*"]
  paths:
    enabled: true
    deny_read: ["/etc/shadow"]
    deny_write: ["**/.env"]
```

Временные блокировки проверяются **до** проверок — [Одобрения](./approvals.md).

См. также: [Конфигурация](./configuration.md), [Маршрутизация](./dispatch.md), [Глоссарий](./glossary.md).
