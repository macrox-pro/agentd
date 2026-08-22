# Trajectory (журнал сессий)

> **Language:** [English](../en/trajectory.md) · [Русский](./trajectory.md)

Опциональный append-only журнал вызовов хуков (`hook/invoked`, `hook/decided`, async meta). **По умолчанию выключен** — в payload могут быть секреты.

## Включение

```yaml
trajectory:
  enabled: true
  include_raw: false          # true нужен для session replay --policy
  redact_secret_rules: true
  max_event_bytes: 262144
  queue_capacity: 1024
  import:
    claude-code:
      enabled: false
      path: ""
    cursor:
      enabled: false
      path: ""
    codex:
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
| `agentd session import --provider ID …` | Импорт transcript (`source=transcript`) |
| `agentd session replay --policy --provider ID --session ID` | Dry-run policy по сохранённому Raw |
| `agentd session fork --provider ID --session SRC --new-session DST` | Копия префикса ledger (аудит) |

## Поиск

`session search` обходит JSONL построчно без индекса. Фильтры: `--provider`, `--session`, `--kind`, `--source`, `--query`, `--limit`.

## Импорт

```bash
agentd session import --provider claude-code --session SESSION_ID
agentd session import --provider cursor --path /path/to/transcript.jsonl
agentd session import --provider codex --path /path/to/transcript.jsonl
```

Transcript-события дописываются после hook-событий (монотонный `seq`). Повторный импорт пропускает строки из sidecar `<session_id>.import.json`. Cursor/Codex — **partial**: предпочитайте `--path`; thinking/tool-output не выдумываются.

## Policy replay

```bash
agentd session replay --policy --provider claude-code --session s1 --json
```

Повторный Invoke сохранённого `Raw` через Dispatch Engine offline. Нужен `trajectory.include_raw: true` при записи. Не обращается к живому агенту и не делает resume цикла.

## Fork

```bash
agentd session fork --provider claude-code --session s1 --new-session s1-fork --at-seq 4
```

Копирует префикс в новый session id и добавляет `session/fork` + `session/end-seed`. Исходный ledger неизменяем. Только аудит — не resume агента.

## Покрытие (L0 и выше)

| Provider | Entrypoint | L0 live | L2 import status | L3 thinking |
|----------|------------|---------|------------------|-------------|
| claude-code | `hook run` | обязательно | **supported** | из session files |
| cursor | `hook run --argv-payload` | обязательно | **partial** (`--path`) | часто redacted |
| codex | `run` + `hook notify` | обязательно | **partial** (`--path`) | маловероятно |
| gemini | `hook run` | обязательно | none | unknown |
| opencode | `hook serve` | обязательно | none | unknown |
| kimi-code | `hook run` | обязательно | none | unknown |

Формулировка: все **поддерживаемые агенты** дают один поток hook-событий; глубина transcript/thinking различается ([DESIGN §14.6](../../DESIGN.md#146-provider-support-matrix-all-supported-agents)).

## Status

`trajectory_dropped_count` — счётчик overflow очереди trajectory ([Эксплуатация](./operations.md)).

См. также: [Configuration](./configuration.md), [CLI](./cli.md).
