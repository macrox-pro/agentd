# Конфигурация

> **Language:** [English](../en/configuration.md) · [Русский](./configuration.md)

Четыре слоя, которые **сливаются** (merge) в один эффективный конфиг; поверхность YAML; как работает перезагрузка. Полные примеры: [DESIGN.md §7](../../DESIGN.md#7-configuration-schema).

## Слои (порядок слияния)

| Порядок | Слой | Где лежит |
|---------|------|-----------|
| 1 | встроенные значения (`defaults`) | внутри бинарника |
| 2 | пользовательский (`user`) | `--config` или `~/.agentd.yaml` |
| 3 | проектный (`project`) | `.agentd.yaml`, поиск вверх от текущей директории / корня проекта |
| 4 | runtime | наложение, которым управляет демон (одобрения, временные блокировки) |

**Путь runtime-файла**

- Unix: `$XDG_STATE_HOME/agentd/runtime.yaml`, иначе `~/.local/state/agentd/runtime.yaml`
- Windows: `%LOCALAPPDATA%\agentd\runtime.yaml`

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
| `offline` | `fail_open` \| `fail_closed` | `fail_closed` | Задумано для офлайна; см. примечание |

> **Важно:** клиент хука (`hook …`) при недоступном демоне пишет в stderr `daemon not running` и выходит с кодом `1`. Ветвления по `policy.offline` в текущем коде **нет**.

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
| `file` | `""` → `$XDG_STATE_HOME/agentd/agentd.log` (Windows: `%LOCALAPPDATA%\agentd\agentd.log`) |

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

При `import.claude-code.enabled: true` демон асинхронно следит за каталогом projects. CLI `session import` работает offline без этого флага.

Переполнение очереди увеличивает `trajectory_dropped_count` в Status.

## Команды для работы с конфигом

| Команда | Роль |
|---------|------|
| `agentd config validate [--config] [--cwd]` | Проверка и компиляция **без** демона |
| `agentd config show [--merged] [--layer user\|project\|runtime] [--cwd]` | Просмотр слоёв |
| `agentd config patch --file DELTA.yaml` | Изменить runtime (с сохранением на диск) |
| `agentd config record-decision …` | Записать одобрение ([Одобрения](./approvals.md)) |

## Перезагрузка

- Изменения user/project: наблюдатель файлов + отложенная пересборка слоёв
- `agentd daemon reload`: принудительно перечитать диск и слить слои
- `patch` / `record-decision`: обновление в памяти + отложенная запись runtime

После успешной компиляции в статусе видны `generation` (поколение) и `fingerprint` (отпечаток слитого конфига).

См. также: [Охранники](./guards.md), [Маршрутизация](./dispatch.md), [Справочник CLI](./cli.md).
