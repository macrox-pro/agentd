# Trajectory (журнал сессий)

> **Language:** [English](../en/trajectory.md) · [Русский](./trajectory.md)

Опциональный append-only журнал вызовов хуков (`hook/invoked`, `hook/decided`, async meta). **По умолчанию выключен** — в payload могут быть секреты.

## Включение

```yaml
trajectory:
  enabled: true
  include_raw: false
  redact_secret_rules: true
  max_event_bytes: 262144
  queue_capacity: 1024
```

Хранение: `$XDG_STATE_HOME/agentd/sessions/<provider>/<session_id>.jsonl`.

Запись идёт по async-пути демона — синхронная задержка хука не меняется.

## CLI

| Команда | Назначение |
|---------|------------|
| `agentd session list [--provider ID] [--json]` | Список сессий (offline) |
| `agentd session show ID --provider ID [--json]` | События одной сессии |
| `agentd session export [--provider ID] [--session ID] [--out PATH]` | Экспорт JSONL |

## Покрытие (L0 и выше)

| Provider | Entrypoint | L0 live | L2 import | L3 thinking |
|----------|------------|---------|-----------|-------------|
| claude-code | `hook run` | обязательно | M10 | из session files, не из hooks |
| cursor | `hook run --argv-payload` | обязательно | M11 partial | часто redacted |
| codex | `run` + `hook notify` | обязательно | none | маловероятно в hooks |
| gemini | `hook run` | обязательно | none | unknown |
| opencode | `hook serve` | обязательно | none | unknown |
| kimi-code | `hook run` | обязательно | none | unknown |

Формулировка: все **поддерживаемые агенты** дают один поток hook-событий; глубина transcript/thinking различается ([DESIGN §14.6](../../DESIGN.md#146-provider-support-matrix-all-supported-agents)).

## Status

`trajectory_dropped_count` — счётчик overflow очереди trajectory ([Эксплуатация](./operations.md)).

См. также: [Configuration](./configuration.md), [CLI](./cli.md).
