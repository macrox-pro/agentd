# Проверки (guards)

> **Language:** [English](../en/guards.md) · [Русский](./guards.md)

Проверки, которые выполняются **до** ответа агенту (обычно событие «перед инструментом»). Настройки — в секции `guards:`.

## Имена

| Проверка | `enabled` по умолчанию | Основные поля |
|----------|------------------------|---------------|
| `secrets` | `true` | `action`: `ask` (спросить) \| `deny` (запретить); опционально `rules` |
| `shell` | `false` | `deny_patterns`, `ask_on` |
| `mcp` | `false` | `deny_servers` |
| `paths` | `false` | `deny_read`, `deny_write` |

В маршруте у sync-цели `builtin` можно сузить список: `guards: [secrets, shell]`. Пустой список — встроенный набор по умолчанию для этой цели.

## Включение через CLI

Проверки оболочки, MCP и путей по умолчанию пишутся в область **проекта** — из корня репозитория:

```bash
cd /path/to/repo
agentd config enable guard-shell
agentd config get guard-shell          # guard-shell: on (project)
```

Команды: `guard-shell`, `guard-mcp`, `guard-paths` ([справочник](./cli.md#config-конфиг)). **`secrets` через CLI не переключается** — по умолчанию включён в автосозданном YAML; меняйте `action` / `rules` в конфиге пользователя или проекта.

## Спросить (Ask) и запретить (Deny)

- **Deny** — жёсткий отказ; агент получает запрет в своём формате.
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
