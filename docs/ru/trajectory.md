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
  import:
    claude-code:
      enabled: false
      path: ""
```

Хранение: `$XDG_STATE_HOME/agentd/sessions/<provider>/<session_id>.jsonl`.

Запись идёт по async-пути демона — синхронная задержка хука не меняется.

## CLI

| Команда | Назначение |
|---------|------------|
| `agentd session list [--provider ID] [--json]` | Список сессий (`importer_status` в `--json`) |
| `agentd session show ID --provider ID [--json]` | События одной сессии |
| `agentd session export …` | Экспорт JSONL |
| `agentd session search …` | Поиск (O(n) скан JSONL) |
| `agentd session import --provider claude-code …` | Импорт transcript (`source=transcript`) |

## Поиск

`session search` обходит JSONL построчно без индекса. Фильтры: `--provider`, `--session`, `--kind`, `--source`, `--query`, `--limit`.

## Импорт (Claude Code)

```bash
agentd session import --provider claude-code --session SESSION_ID
agentd session import --provider claude-code --path /path/to/session.jsonl
```

Transcript-события дописываются после hook-событий (монотонный `seq`). Повторный импорт пропускает строки из sidecar `<session_id>.import.json`. Корреляция по `session_id` и `tool_use_id`.

## Покрытие (L0 и выше)

| Provider | Entrypoint | L0 live | L2 import status | L3 thinking |
|----------|------------|---------|------------------|-------------|
| claude-code | `hook run` | обязательно | **supported** | из session files |
| cursor | `hook run --argv-payload` | обязательно | none (M11 partial) | часто redacted |
| codex | `run` + `hook notify` | обязательно | none | маловероятно |
| gemini | `hook run` | обязательно | none | unknown |
| opencode | `hook serve` | обязательно | none | unknown |
| kimi-code | `hook run` | обязательно | none | unknown |

Формулировка: все **поддерживаемые агенты** дают один поток hook-событий; глубина transcript/thinking различается ([DESIGN §14.6](../../DESIGN.md#146-provider-support-matrix-all-supported-agents)).

## Status

`trajectory_dropped_count` — счётчик overflow очереди trajectory ([Эксплуатация](./operations.md)).

См. также: [Configuration](./configuration.md), [CLI](./cli.md).
