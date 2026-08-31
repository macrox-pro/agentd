# Журнал сессий (trajectory)

> **Language:** [English](../en/trajectory.md) · [Русский](./trajectory.md)

Хронологическая запись вызовов хуков и связанных событий. **Включена по умолчанию.** В данных могут быть секреты; маскирование по правилам секретов (`redact_secret_rules`) включено по умолчанию. Термины: [Глоссарий](./glossary.md).

## Включение

Журнал включён по умолчанию. Проверка или отключение без правки YAML ([Переключатели](./configuration.md#переключатели)):

```bash
agentd config get trajectory          # trajectory: on (default)

agentd config disable trajectory

agentd config disable trajectory-raw   # если не нужен session replay --policy
```

Или через YAML:

```yaml
trajectory:
  enabled: true
  include_raw: true
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

Хранение: `sessions/<provider>/<session_id>.jsonl` в [каталоге состояния](./configuration.md#каталог-состояния). Запись идёт асинхронно — синхронная задержка хука не меняется.

## CLI

Офлайн-команды читают `sessions/` в [каталоге состояния](./configuration.md#каталог-состояния). `session subscribe` и `trajectory stats` нужен запущенный демон. Полный справочник: [CLI → session](./cli.md#session).

## Статистика демона

```bash
agentd config get trajectory
agentd config get trajectory-statistics
agentd trajectory stats [--provider ID] [--json]
```

Нужен **запущенный демон**. Счётчики сбрасываются при перезапуске процесса; `since` — время старта демона. `--provider` фильтрует сводку.

Счётчики токенов демона берутся из сырой полезной нагрузки каждого вызова (не зависят от `include_raw`):

- **Cursor** — токены биллинга на хуке `stop` (кумулятивно по сессии, с дельтой); `context_tokens_last` на `preCompact`.
- **Codex** — токены из хвоста rollout-транскрипта на `Stop`, если в raw хука нет usage (`transcript_path` в raw).

Для офлайн `session stats` поля токенов в JSONL требуют `include_raw` (запасной путь Codex — `transcript_path` в сохранённом raw).

## Subscribe (поток в реальном времени)

`session subscribe` читает журнал демона в памяти с момента подключения. Фильтры: `--provider`, `--session`, `--source`. История **не** воспроизводится — `session show` или `session export`.

- Нужен **запущенный демон** с `trajectory.enabled`.
- Офлайн `session import` / `fork` не публикуют в поток; наблюдатель импорта Claude в демоне — публикует.
- `schema_version: 1` на каждом событии; `ignorable` — подсказки совместимости.
- `raw` на потоке следует тем же правилам маскирования, что и JSONL.
- При `trajectory.enabled` записывается каждый вызов хука — отдельного webhook или цели dispatch для журнала нет.

## Контракт журнала

| Поле | Смысл |
|------|-------|
| `schema_version` | Зафиксировано `1` с момента появления журнала сессий |
| `seq` | Монотонный номер в рамках сессии |
| `type` | Каталог событий (`hook/invoked`, `transcript/message`, …) |
| `source` | `hook`, `decision`, `transcript`, `system` |
| `ignorable` | Старые читатели могут пропускать неизвестные **type** |

Каталог: [DESIGN §14.2](../../DESIGN.md#142-event-catalog).

## Поиск

`session search` обходит JSONL построчно. Индекса нет — фильтры: `--provider`, `--session`, `--kind`, `--source`, `--query`, `--limit`.

## Импорт

```bash
agentd session import --provider claude-code --session SESSION_ID
agentd session import --provider cursor --path /path/to/transcript.jsonl
agentd session import --provider codex --session SESSION_ID
```

События транскрипта дописываются после событий хуков. Повторный импорт использует сопутствующий файл `<session_id>.import.json`. Cursor — **частично** (лучше `--path`). Codex — **поддерживается**.

### Импорт в stdout или файл (`--out`)

```bash
agentd session import --provider claude-code --path /path/to/session.jsonl --out -
```

Не дописывает в журнал на диске. Подробнее: [CLI → Импорт без записи в журнал](./cli.md#импорт-без-записи-в-журнал).

## Policy replay

```bash
agentd session replay --policy --provider claude-code --session s1 --json
```

Пробный прогон сохранённых событий через маршрутизацию офлайн. Нужен `include_raw` при записи. **Не** обращается к живому агенту.

## Fork

```bash
agentd session fork --provider claude-code --session s1 --new-session s1-fork --at-seq 4
```

Копия префикса в новый id сессии. Только аудит происхождения — не возобновление агента.

## Уровни покрытия

| Уровень | Смысл |
|---------|--------|
| **L0 Live** | Каждый вызов хука → `hook/invoked` + `hook/decided` (обязательно для всех шести агентов). |
| **L1 Correlate** | Стабильные id сессии и инструмента (качество разное; отдельно в доке не выделяем). |
| **L2 Import** | Транскрипт с диска → события `transcript/*` через `session import`. |
| **L3 Thinking** | Строки рассуждений, если поставщик их сохраняет. |

| Поставщик | Точка входа | L0 | L2 импорт | L3 |
|-----------|-------------|----|-----------|-----|
| claude-code | `hook run` | да | **поддерживается** | из файлов сессии |
| cursor | `hook run --argv-payload` | да | **частично** (`--path`) | часто нет |
| codex | `run` + `hook notify` | да | **поддерживается** | plaintext `agent_reasoning` |
| gemini | `hook run` | да | нет | неизвестно |
| opencode | `hook serve` | да | нет | неизвестно |
| kimi-code | `hook run` | да | нет | неизвестно |

Все поддерживаемые агенты дают один поток событий хуков; глубина транскрипта и рассуждений различается ([DESIGN §14.3](../../DESIGN.md#143-provider-support-matrix)).

## Статус

`trajectory_dropped_count` — счётчик отброшенных записей при переполнении очереди ([Эксплуатация](./operations.md)).

См. также: [Конфигурация](./configuration.md), [Справочник команд](./cli.md), [Глоссарий](./glossary.md).
