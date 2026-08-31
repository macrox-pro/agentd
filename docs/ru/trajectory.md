# Trajectory (журнал сессий)

> **Language:** [English](../en/trajectory.md) · [Русский](./trajectory.md)

Опциональный append-only журнал вызовов хуков (`hook/invoked`, `hook/decided`, async meta). **По умолчанию выключен** — в payload могут быть секреты.

## Включение

Без правки YAML ([Переключатели features](./configuration.md#переключатели-features)):

```bash
agentd config enable trajectory
agentd config get trajectory          # trajectory: on (user)

# Raw для session replay --policy (можно до или после trajectory)
agentd config enable trajectory-raw
```

Или через YAML:

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

Хранение: `sessions/<provider>/<session_id>.jsonl` в [state directory](./configuration.md#state-directory).

Запись идёт по async-пути демона — синхронная задержка хука не меняется.

## CLI

| Команда | Назначение |
|---------|------------|
| `agentd session list [--provider ID] [--json]` | Список сессий (`importer_status` в `--json`) |
| `agentd session show ID --provider ID [--json]` | События одной сессии |
| `agentd session export …` | Экспорт JSONL |
| `agentd session search …` | Поиск (O(n) скан JSONL) |
| `agentd session import --provider ID …` | Импорт transcript (`source=transcript`) или `--out` — parse-only JSONL |
| `agentd session replay --policy --provider ID --session ID` | Dry-run policy по сохранённому Raw |
| `agentd session fork --provider ID --session SRC --new-session DST` | Копия префикса ledger (аудит) |
| `agentd session stats ID --provider ID [--json]` | Offline-статистика сессии (нужен `trajectory.statistics`) |
| `agentd session subscribe [--json]` | **Live** поток от демона (нужен запущенный daemon + trajectory.enabled) |

## Статистика демона

```bash
agentd config enable trajectory
agentd config enable trajectory-statistics
agentd trajectory stats [--provider ID] [--json]
```

Нужен **запущенный daemon**. Счётчики сбрасываются при перезапуске; `since` — время старта демона. Опциональный `--provider` фильтрует rollup. Токены daemon rollup извлекаются из `RawPayload` каждого Invoke (не зависят от `include_raw`). Billing-токены Cursor — на хуке `stop` (кумулятивные по сессии, дельта-агрегация); `context_tokens_last` — на `preCompact`. Для offline `session stats` токены в JSONL нужен `include_raw`.

## Subscribe (live-поток)

`session subscribe` читает in-memory ledger демона с момента подключения (gRPC `SessionService.Subscribe`). Фильтры: `--provider`, `--session`, `--source`. История **не** воспроизводится — используйте `session show` или `session export`.

- Нужен **запущенный демон** с `trajectory.enabled`.
- Offline `session import` / `fork` не публикуют в Subscribe; watcher импорта Claude — публикует.
- `schema_version: 1` на каждом событии; `ignorable` — forward-compat (пропуск неизвестных **type**; transcript всё равно доставляется).
- `raw` на потоке следует тем же правилам redaction, что и JSONL.
- Зеркало ledger: `trajectory.enabled` — отдельный webhook / `target: trajectory` в M12 нет.

## Контракт trajectory

| Поле | Смысл |
|------|-------|
| `schema_version` | Зафиксировано `1` для v0.0.2 |
| `seq` | Монотонный в рамках сессии |
| `type` | Каталог событий |
| `source` | `hook`, `decision`, `transcript`, `system` |
| `ignorable` | Forward-compat для неизвестных **type** |

Каталог событий: [DESIGN §14.2](../../DESIGN.md#142-event-catalog).

## Поиск

`session search` обходит JSONL построчно без индекса. Фильтры: `--provider`, `--session`, `--kind`, `--source`, `--query`, `--limit`.

## Импорт

```bash
agentd session import --provider claude-code --session SESSION_ID
agentd session import --provider cursor --path /path/to/transcript.jsonl
agentd session import --provider codex --session SESSION_ID
agentd session import --provider codex --path /path/to/rollout-…-SESSION_ID.jsonl
```

Transcript-события дописываются после hook-событий (монотонный `seq`). Повторный импорт пропускает строки из sidecar `<session_id>.import.json`. Cursor — **partial** (лучше `--path`; thinking/tool-output не выдумываются). Codex — **supported** через `~/.codex/sessions/**/rollout-*-{session_id}.jsonl` (thinking только из plaintext `agent_reasoning`).

### Импорт в stdout или файл (`--out`)

Preview/transcode без изменения ledger:

```bash
agentd session import --provider claude-code --path /path/to/session.jsonl --out -
agentd session import --provider codex --session SESSION_ID --out /tmp/events.jsonl
agentd session import --provider claude-code --session s1 --out - 2>/dev/null | wc -l
```

- Не дописывает в `sessions/…jsonl` и не обновляет `<session>.import.json`.
- Учитывает инкрементальный импорт (читает checkpoint для `startIndex`, если sidecar есть).
- `seq` совпадает с обычным import (включая продолжение после hook-событий).
- При `--out` summary на stderr; `--json` — машиночитаемый summary там же.

Для записи в ledger и Subscribe — обычный import (без `--out`). Подробнее: [CLI §session import --out](./cli.md#session-import-out).

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
| codex | `run` + `hook notify` | обязательно | **supported** (`~/.codex/sessions` rollouts) | только plaintext `agent_reasoning` |
| gemini | `hook run` | обязательно | none | unknown |
| opencode | `hook serve` | обязательно | none | unknown |
| kimi-code | `hook run` | обязательно | none | unknown |

Формулировка: все **поддерживаемые агенты** дают один поток hook-событий; глубина transcript/thinking различается ([DESIGN §14.3](../../DESIGN.md#143-provider-support-matrix)).

## Status

`trajectory_dropped_count` — счётчик overflow очереди trajectory ([Эксплуатация](./operations.md)).

См. также: [Configuration](./configuration.md), [CLI](./cli.md).
