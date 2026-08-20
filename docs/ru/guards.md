# Guards

> **Language:** [English](../en/guards.md) · [Русский](./guards.md)

Декларативные sync-проверки на `tool.pre` (и смежных kind через builtin). Конфиг в `guards:`.

## Имена

| Guard | Default `enabled` | Основные поля |
|-------|-------------------|---------------|
| `secrets` | `true` | `action`: `ask` \| `deny` (default `ask`); опционально `rules` |
| `shell` | `false` | `deny_patterns`, `ask_on` |
| `mcp` | `false` | `deny_servers` |
| `paths` | `false` | `deny_read`, `deny_write` |

Sync builtin на route может сузить набор: `guards: [secrets, shell]`. Пустой список — compiled default для этого target.

## Ask vs Deny

- **Deny** — жёсткий стоп; provider кодирует deny.
- **Ask** — Ask / permission у агента; в сообщении может быть `approval_fingerprint=sha256:…` для [Approvals](./approvals.md).

Ограничения capability провайдера сохраняются (не все агенты умеют Ask).

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

Temporary blocks (runtime) выполняются до guards — см. [Approvals](./approvals.md).

Архитектура: [DESIGN.md](../../DESIGN.md). Поведение на wire: [Dispatch](./dispatch.md).
