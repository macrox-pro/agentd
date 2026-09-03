# Конфигурация

> **Language:** [English](../en/configuration.md) · [Русский](./configuration.md)

agentd собирает один рабочий конфиг из **четырёх слоёв**. Поздний слой перекрывает ранний. Ключи YAML и перезагрузка — ниже. Договорённость о слоях: [DESIGN.md §7](../../DESIGN.md#7-configuration-schema).

## Каталог состояния

Конфиг пользователя — файл `~/.agentd.yaml`, не папка `~/.agentd/`. То, что пишет демон (временный слой, журнал работы, журнал сессий) — **состояние**: его можно заново создать, в git его не кладут. В Linux и macOS это `$XDG_STATE_HOME` (если не задан → `~/.local/state`). Windows: `%LOCALAPPDATA%\agentd\`. Сокет — **не** состояние: он живёт в `$XDG_RUNTIME_DIR` и исчезает при выходе из сеанса.

| Что | Путь по умолчанию |
|-----|-------------------|
| Пользовательский конфиг | `~/.agentd.yaml` (файл) |
| Каталог состояния | `$XDG_STATE_HOME/agentd/` иначе `~/.local/state/agentd/` (Windows: `%LOCALAPPDATA%\agentd\`) |
| → временный слой | `runtime.yaml` (только демон) |
| → журнал работы | `agentd.log` |
| → журнал сессий | `sessions/<provider>/<session_id>.jsonl` (если включён) |
| Сокет (не состояние) | `$XDG_RUNTIME_DIR/agentd/agentd.sock` (запасной путь на Darwin `~/Library/Caches/agentd/`; Linux `~/.local/run/agentd/`; иначе временный каталог) — [DESIGN.md §5](../../DESIGN.md#5-transport) |
| → PID-файл | `agentd.pid` (рядом с сокетом; в Windows — в каталоге состояния) |
| → файл блокировки | `agentd.lock` (рядом с сокетом; в Windows — в каталоге состояния; один демон на пользователя) |

## Автосоздание пользовательского конфига

Только при **`agentd daemon start`**: если файла нет, демон записывает минимальный конфиг (те же ключи, что в [быстром старте](./getting-started.md)) и продолжает запуск. Команды только для чтения (`config show`, `config validate`, хуки) **не** создают файл.

| Ситуация | Поведение |
|----------|-----------|
| Файла нет | Файл создаётся без сообщений; старт продолжается |
| Файл верный | Без изменений; старт продолжается |
| Неверный YAML или ошибка сборки | Stderr: `agentd: invalid user config <path>: …`; старт не выполняется; файл не меняется |
| Путь недоступен / ошибка ввода-вывода | Старт не выполняется; без сообщения «invalid user config» |

Проверка и правка офлайн, затем перезапуск:

```bash
agentd config validate --config ~/.agentd.yaml
```

`config show` нормализует YAML (пустые поля опускаются; ключей `null` нет).

## Слои (порядок слияния)

| Порядок | Слой | Где лежит |
|---------|------|-----------|
| 1 | встроенные значения (`defaults`) | внутри бинарника |
| 2 | пользовательский (`user`) | `--config` или `~/.agentd.yaml` |
| 3 | проектный (`project`) | `.agentd.yaml`, поиск вверх от **рабочего каталога в JSON события** (см. ниже) или от `--cwd` в CLI |
| 4 | временный слой (`runtime`) | файл демона (одобрения, временные блокировки) |

**Путь временного слоя:** `runtime.yaml` в [каталоге состояния](#каталог-состояния).

Запись на диск откладывается на **500 мс**, права файла `0600`, запись атомарная. Каждый вызов хука читает снимок в памяти — без чтения диска на вызов.

### Рабочий каталог события и проектный конфиг

При **работающем демоне** `hook run`, `hook notify` и `hook serve` передают демону рабочий каталог вместе с каждым вызовом. Порядок выбора:

1. Поле `cwd` в JSON события (для Codex notify — в argv), если указано
2. Иначе первый элемент `workspace_roots[]`, если он есть (формат Cursor)
3. Иначе текущая директория процесса внешней части хука — когда в JSON нет ни `cwd`, ни `workspace_roots`

Демон поднимается вверх от этого пути и подмешивает `.agentd.yaml` — по тому же правилу, что `config show --cwd`. Обычно это каталог проекта из JSON агента, а не каталог shell, из которого агент запустил процесс.

Команды `config validate`, `config show` и `config get` берут каталог из флага `--cwd`, а не из JSON хука. См. [CLI → config](./cli.md#config).

## Ключи верхнего уровня YAML

Из схемы файла: `version`, `policy`, `async`, `logging`, `guards`, `approvals`, `blocks`, `dispatch_defaults`, `dispatch`, `trajectory`, `metrics`.

Ключи `dispatch_defaults` и значения `dispatch[].match.kind` — только известные wire-виды ([Маршрутизация](./dispatch.md#виды-событий-kind)). Опечатка валит `config validate` и `daemon start`. Bootstrap `~/.agentd.yaml` остаётся коротким — defaults видов живут в бинарнике.

В проектном файле обычно `guards` / `dispatch`. Блоки `approvals` и `blocks` чаще попадают в runtime через CLI или gRPC.

### policy

| Ключ | Значения | По умолчанию | Смысл |
|------|----------|--------------|-------|
| `fail` | `fail_open` \| `fail_closed` | `fail_closed` | Если синхронный конвейер **в демоне** завершился **ошибкой** (истёк бюджет, отмена контекста, ошибка sync-цели). Обычный **deny** от проверки — уже готовое решение; `fail` его не перекрывает. У sync-цели `grpc` часть сбоев узла закрывается полем `on_error` раньше, чем применится `policy.fail` ([Маршрутизация](./dispatch.md)) |
| `ask_fallback` | `deny` \| `no_decision` | `deny` | Когда агент не умеет спрашивать пользователя: **deny** блокирует; **no_decision** — нейтральное разрешение |
| `offline` | `fail_open` \| `fail_closed` | `fail_open` | Когда **внешняя часть хука** не может вызвать демон (см. ниже) |

`policy.fail` действует только когда запрос обрабатывается **демоном**. `policy.offline` — только когда внешняя часть хука работает локально после неудачного подключения или вызова демона.

Когда демон недоступен, внешняя часть хука читает локальный конфиг (встроенные defaults, затем user, проектный слой по рабочему каталогу и runtime, слой за слоем) и применяет `policy.offline`. По умолчанию `fail_open` — нейтральное решение (или код 0 для notify), агент продолжает работу; `fail_closed` — выход с кодом **1**. В обоих режимах в stderr пишется `daemon not running`.

### async

| Ключ | По умолчанию |
|------|--------------|
| `queue_capacity` | `1024` |
| `worker_limit` | `8` |
| `target_timeout` | `30s` |
| `on_overflow` | `drop` (`drop` \| `log`) |

При переполнении задача **всегда** отбрасывается, счётчик `async_dropped_count` в статусе растёт.

### logging

Операционные логи демона (не асинхронная цель dispatch `target: log`).

| Ключ | По умолчанию |
|------|--------------|
| `level` | `info` (`debug` \| `info` \| `warn` \| `error`) |
| `file` | `""` → `agentd.log` в [каталоге состояния](#каталог-состояния) |

`agentd daemon start --foreground` дублирует логи в stderr и в файл. Флаги CLI `--log-level` и `--log-file` переопределяют YAML только для этого процесса.

### trajectory

Журнал сессий ([Журнал сессий](./trajectory.md)). По умолчанию **включён**.

| Ключ | По умолчанию |
|------|--------------|
| `enabled` | `true` |
| `statistics` | `true` (нужен `enabled`; для сводной статистики демона и `session stats`) |
| `include_raw` | `true` |
| `redact_secret_rules` | `true` |
| `max_event_bytes` | `262144` |
| `queue_capacity` | `1024` |
| `import.claude-code.enabled` | `false` |
| `import.claude-code.path` | `""` (по умолчанию `~/.claude/projects`) |
| `import.cursor.enabled` | `false` |
| `import.cursor.path` | `""` (лучше CLI `--path`) |
| `import.codex.enabled` | `false` |
| `import.codex.path` | `""` (по умолчанию `$CODEX_HOME/sessions` или `~/.codex/sessions`) |

При `import.claude-code.enabled: true` демон следит за каталогом projects и асинхронно дописывает новые строки транскрипта. CLI `session import` работает офлайн без этого флага. Задайте `include_raw: true`, если нужен `session replay --policy`.

Переполнение очереди увеличивает `trajectory_dropped_count` в Status.

### metrics

Опциональный endpoint Prometheus ([Эксплуатация → Метрики Prometheus](./operations.md#метрики-prometheus)). По умолчанию **выключен**.

| Ключ | По умолчанию |
|------|--------------|
| `enabled` | `false` |
| `listen` | `127.0.0.1:2112` (`host:port`; обязателен при `enabled: true`) |

При включении демон отдаёт `/metrics` по loopback TCP только на старте. Смена `metrics.listen` или `enabled` требует **`agentd daemon stop`**, затем **`agentd daemon start`** — `daemon reload` не перепривязывает HTTP listener метрик. CLI `--metrics-listen` включает метрики для этого процесса и переопределяет `listen`.

Привязка к `0.0.0.0` допустима, но открывает метрики на всех интерфейсах — предпочитайте loopback, если не понимаете риск.

## Переключатели

`agentd config enable|disable|get FEATURE` включает и выключает готовые флаги без ручного YAML ([команды](./cli.md#config)).

| Поведение | Детали |
|-----------|--------|
| Запись | Пользователь (`--config` / `~/.agentd.yaml`) или проект (`.agentd.yaml` под `--cwd`) |
| Временный слой | **Не** меняется — временные правки через `config patch` |
| Нет файла пользователя | `config enable` создаёт тот же файл, что и `daemon start` |
| Перезагрузка | Слежение за файлами user/project; или `agentd daemon reload` |
| `config get` | встроенные, user и project, слой за слоем; **без** временного слоя |
| Повтор | Повторный enable/disable при том же значении → код 0, без записи |
| YAML | Сохранение может убрать ручные комментарии в изменённых файлах |
| Проверка secrets | Не переключается этой командой — правьте `guards.secrets` в YAML |
| Чтение vs запись проекта | `config get` ищет `.agentd.yaml` вверх; `enable`/`disable` (project) пишет только под `--cwd` |

**Примеры stdout:**

```text
agentd config get trajectory
trajectory: on (user)

agentd config enable guard-shell    # из корня репо; project scope по умолчанию
guard-shell: enabled (project /path/to/repo/.agentd.yaml)

agentd config enable trajectory     # уже включено
trajectory: already enabled (user /home/you/.agentd.yaml)
```

См. также: [Журнал сессий](./trajectory.md#включение), [Проверки](./guards.md#включение-через-cli).

## Команды для работы с конфигом

| Команда | Роль |
|---------|------|
| `agentd config validate [--config] [--cwd]` | Офлайн-разбор и компиляция |
| `agentd config show [--merged] [--layer user\|project\|runtime] [--cwd]` | Просмотр слоёв |
| `agentd config enable\|disable\|get FEATURE` | Готовые постоянные переключатели ([Переключатели](#переключатели)) |
| `agentd config patch --file DELTA.yaml` | Изменить runtime (с сохранением на диск) |
| `agentd config record-decision …` | Записать одобрение ([Одобрения](./approvals.md)) |

## Перезагрузка

- Изменения user/project: наблюдатель файлов + отложенная пересборка слоёв
- `agentd daemon reload`: принудительно перечитать диск и слить слои
- `patch` / `record-decision`: обновление в памяти + отложенная запись runtime

После успешной компиляции в статусе видны `generation` (поколение) и `fingerprint` (отпечаток слитого конфига).

См. также: [Проверки](./guards.md), [Маршрутизация](./dispatch.md), [Справочник команд](./cli.md).
