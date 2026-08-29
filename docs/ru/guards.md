# Охранники (guards)

> **Language:** [English](../en/guards.md) · [Русский](./guards.md)

Декларативные **синхронные** проверки на событии `tool.pre` (и смежных видах через встроенную цель `builtin`). Настройки — в секции `guards:`.

## Имена

| Охранник | `enabled` по умолчанию | Основные поля |
|----------|------------------------|---------------|
| `secrets` | `true` | `action`: `ask` (спросить) \| `deny` (запретить); опционально `rules` |
| `shell` | `false` | `deny_patterns`, `ask_on` |
| `mcp` | `false` | `deny_servers` |
| `paths` | `false` | `deny_read`, `deny_write` |

В маршруте у sync-цели `builtin` можно сузить список: `guards: [secrets, shell]`. Пустой список — встроенный набор по умолчанию для этой цели.

## Включение через CLI

Охранники shell, MCP и paths по умолчанию пишутся в **project** scope — из корня репозитория:

```bash
cd /path/to/repo
agentd config enable guard-shell
agentd config get guard-shell          # guard-shell: on (project)
```

Features: `guard-shell`, `guard-mcp`, `guard-paths` ([CLI](./cli.md#config-конфиг)). **`secrets` не curated toggle** — по умолчанию включён в bootstrap YAML; меняйте `action` / `rules` в user или project конфиге.

## Спросить (Ask) и запретить (Deny)

- **Deny** — жёсткий отказ; провайдер получает запрет в своём формате.
- **Ask** — запрос разрешения у пользователя агента; в тексте может быть `approval_fingerprint=sha256:…` для последующей записи одобрения ([Одобрения](./approvals.md)).

Ограничения протокола конкретного агента сохраняются: не все умеют режим «спросить».

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

Временные блокировки из runtime проверяются **до** охранников — см. [Одобрения](./approvals.md).

Архитектура: [DESIGN.md](../../DESIGN.md). Как это стыкуется с маршрутами: [Маршрутизация](./dispatch.md).
