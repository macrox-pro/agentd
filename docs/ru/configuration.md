# Конфигурация

> **Language:** [English](../en/configuration.md) · [Русский](./configuration.md)

Четыре слоя, которые **сливаются** (merge) в один эффективный конфиг; поверхность YAML; как работает перезагрузка. Контракт слоёв и runtime overlay: [DESIGN.md §7](../../DESIGN.md#7-configuration-schema).

## State directory (каталог состояния)

Пользовательский конфиг — файл `~/.agentd.yaml`, не дерево `~/.agentd/`. Изменяемые данные демона (runtime overlay, операционный лог, ledger сессий) — это **state**: пишет демон, данные восстановимы и не должны ехать вместе с бэкапами конфига или в git. [XDG Base Directory](https://specifications.freedesktop.org/basedir-spec/latest/) относит этот класс к `$XDG_STATE_HOME` (если не задан → `~/.local/state`). На Windows — `%LOCALAPPDATA%\agentd\`. IPC (сокет / блокировка) — **runtime**, не state: `$XDG_RUNTIME_DIR` живёт в сессии ОС и очищается при logout.

| Что | Путь по умолчанию |
|-----|-------------------|
| Пользовательский конфиг | `~/.agentd.yaml` (файл) |
| State directory | `$XDG_STATE_HOME/agentd/` иначе `~/.local/state/agentd/` (Windows: `%LOCALAPPDATA%\agentd\`) |
| → runtime overlay | `runtime.yaml` (только демон) |
| → операционный лог | `agentd.log` |
| → ledger trajectory | `sessions/<provider>/<session_id>.jsonl` (если включён) |
| IPC-сокет (не state) | `$XDG_RUNTIME_DIR/agentd/agentd.sock` (fallback Darwin: `~/Library/Caches/agentd/`; Linux: `~/.local/run/agentd/`; иначе temp) — [DESIGN.md §5](../../DESIGN.md#5-transport) |

## Bootstrap пользовательского конфига

Только при **`agentd daemon start`**: если файла пользовательского конфига нет, демон записывает минимальный bootstrap (те же ключи, что в [быстром старте](./getting-started.md)) и продолжает запуск. Команды только для чтения (`config show`, `config validate`, хуки) **не** создают файл.

| Ситуация | Поведение |
|----------|-----------|
| Файла нет | Bootstrap пишется без сообщений; старт продолжается |
| Файл валиден | Без изменений; старт продолжается |
| Невалидный YAML или ошибка compile | Stderr: `agentd: invalid user config <path>: …`; старт не выполняется; файл не меняется |
| Путь недоступен / ошибка I/O | Старт не выполняется; без сообщения «invalid user config» |

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
| 3 | проектный (`project`) | `.agentd.yaml`, поиск вверх от текущей директории / корня проекта |
| 4 | runtime | наложение, которым управляет демон (одобрения, временные блокировки) |

**Путь runtime-файла:** `runtime.yaml` в [state directory](#state-directory).

Запись на диск откладывается на **500 ms** (debounce), права файла `0600`, запись атомарная (временный файл + rename). На горячем пути обработки запроса используется только снимок в памяти (`store.Current()`), без чтения диска на каждый `Invoke`.

## Ключи верхнего уровня YAML

Из схемы файла: `version`, `policy`, `async`, `logging`, `guards`, `approvals`, `blocks`, `dispatch_defaults`, `dispatch`, `trajectory`.

В проектном файле обычно `guards` / `dispatch`. Блоки `approvals` и `blocks` чаще попадают в runtime через CLI или gRPC.

### policy (политика ошибок)

| Ключ | Значения | По умолчанию | Смысл |
|------|----------|--------------|--------|
| `fail` | `fail_open` \| `fail_closed` | `fail_closed` | При сбое: пропустить или закрыть |
| `unsupported` | `degrade` \| `strict` | `degrade` | Неподдерживаемое: смягчить или строго |
| `ask_fallback` | `deny` \| `no_decision` | `deny` | Если «спросить» недоступно |
| `offline` | `fail_open` \| `fail_closed` | `fail_open` | Демон недоступен: пропустить или закрыть |

Когда демон недоступен, край хука читает локальный конфиг (defaults ⊕ user ⊕ project(cwd) ⊕ runtime) и применяет `policy.offline`. По умолчанию `fail_open` — нейтральное решение (или код 0 для notify), агент продолжает работу; `fail_closed` — выход с кодом **1**. В обоих режимах в stderr пишется `daemon not running`.

### async (асинхронная очередь)

| Ключ | По умолчанию | Смысл |
|------|--------------|--------|
| `queue_capacity` | `1024` | Ёмкость очереди |
| `worker_limit` | `8` | Число воркеров |
| `target_timeout` | `30s` | Таймаут одной асинхронной цели |
| `on_overflow` | `drop` (`drop` \| `log`) | При переполнении: отбросить; `log` ещё пишет предупреждение |

При переполнении задача **всегда** отбрасывается, счётчик `async_dropped_count` в статусе растёт.

### logging (операционные логи демона)

Не путать с асинхронной целью dispatch `target: log`.

| Ключ | По умолчанию |
|------|--------------|
| `level` | `info` (`debug` \| `info` \| `warn` \| `error`) |
| `file` | `""` → `agentd.log` в [state directory](#state-directory) |

`agentd daemon start --foreground` дублирует логи в stderr и в файл. Флаги CLI `--log-level` и `--log-file` переопределяют YAML только для этого процесса.

### trajectory (журнал сессий)

Опциональный ledger ([Trajectory](./trajectory.md)). По умолчанию **выключен**.

| Ключ | По умолчанию |
|------|--------------|
| `enabled` | `false` |
| `include_raw` | `false` |
| `redact_secret_rules` | `true` |
| `max_event_bytes` | `262144` |
| `queue_capacity` | `1024` |
| `import.claude-code.enabled` | `false` |
| `import.claude-code.path` | `""` (по умолчанию `~/.claude/projects`) |
| `import.cursor.enabled` | `false` |
| `import.cursor.path` | `""` (лучше CLI `--path`) |
| `import.codex.enabled` | `false` |
| `import.codex.path` | `""` (по умолчанию `$CODEX_HOME/sessions` или `~/.codex/sessions`) |

При `import.claude-code.enabled: true` демон асинхронно следит за каталогом projects. CLI `session import` работает offline без этого флага. Для `session replay --policy` нужен `include_raw: true` при записи.

Переполнение очереди увеличивает `trajectory_dropped_count` в Status.

## Переключатели features

`agentd config enable|disable|get FEATURE` переключает curated-boolean без ручного YAML ([CLI](./cli.md#config-конфиг)).

| Поведение | Детали |
|-----------|--------|
| Запись | User (`--config` / `~/.agentd.yaml`) или project (`.agentd.yaml` под `--cwd`) |
| Runtime overlay | **Не** изменяется — временные override через `config patch` |
| Нет user-файла | `config enable` создаёт тот же bootstrap, что и `daemon start` |
| Reload | fsnotify на user/project; или `agentd daemon reload` |
| `config get` | defaults ⊕ user ⊕ project; **без** runtime |
| Idempotent | Повторный enable/disable при том же effective → exit 0, без записи |
| YAML round-trip | Marshal может удалить ручные комментарии в изменённых файлах |
| Охранник secrets | Не curated toggle — правьте `guards.secrets` в YAML |
| Чтение vs запись project | `config get` ищет `.agentd.yaml` вверх; `enable`/`disable` (project) пишет только под `--cwd` |

**Примеры stdout:**

```text
agentd config get trajectory
trajectory: on (user)

agentd config enable guard-shell    # из корня репо; project scope по умолчанию
guard-shell: enabled (project /path/to/repo/.agentd.yaml)

agentd config enable trajectory     # уже включено
trajectory: already enabled (user /home/you/.agentd.yaml)
```

См. также: [Trajectory](./trajectory.md#включение), [Guards](./guards.md#включение-через-cli).

## Команды для работы с конфигом

| Команда | Роль |
|---------|------|
| `agentd config validate [--config] [--cwd]` | Проверка и компиляция **без** демона |
| `agentd config show [--merged] [--layer user\|project\|runtime] [--cwd]` | Просмотр слоёв |
| `agentd config enable\|disable\|get FEATURE` | Curated persistent toggles ([Переключатели features](#переключатели-features)) |
| `agentd config patch --file DELTA.yaml` | Изменить runtime (с сохранением на диск) |
| `agentd config record-decision …` | Записать одобрение ([Одобрения](./approvals.md)) |

## Перезагрузка

- Изменения user/project: наблюдатель файлов + отложенная пересборка слоёв
- `agentd daemon reload`: принудительно перечитать диск и слить слои
- `patch` / `record-decision`: обновление в памяти + отложенная запись runtime

После успешной компиляции в статусе видны `generation` (поколение) и `fingerprint` (отпечаток слитого конфига).

См. также: [Охранники](./guards.md), [Маршрутизация](./dispatch.md), [Справочник CLI](./cli.md).
