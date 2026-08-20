# Конфигурация

> **Language:** [English](../en/configuration.md) · [Русский](./configuration.md)

Четыре слоя merge, YAML-поверхность и reload. Полные примеры: [DESIGN.md §7](../../DESIGN.md#7-configuration-schema).

## Слои (порядок merge)

| Порядок | Слой | Расположение |
|---------|------|--------------|
| 1 | defaults | в бинарнике |
| 2 | user | `--config` или `~/.agentd.yaml` |
| 3 | project | `.agentd.yaml` вверх от CWD / project root |
| 4 | runtime | overlay демона (approvals, temporary blocks) |

**Путь runtime**

- Unix: `$XDG_STATE_HOME/agentd/runtime.yaml`, иначе `~/.local/state/agentd/runtime.yaml`
- Windows: `%LOCALAPPDATA%\agentd\runtime.yaml`

Запись runtime с debounce **500ms**, режим `0600`, атомарный rename. Hot path — только `store.Current()`, без disk I/O на Invoke.

## Top-level ключи YAML

Из схемы файла: `version`, `policy`, `async`, `guards`, `approvals`, `blocks`, `dispatch_defaults`, `dispatch`.

В project обычно `guards` / `dispatch`. `approvals` и `blocks` чаще попадают в runtime через CLI/gRPC.

### policy

| Ключ | Значения | Default |
|------|----------|---------|
| `fail` | `fail_open` \| `fail_closed` | `fail_closed` |
| `unsupported` | `degrade` \| `strict` | `degrade` |
| `ask_fallback` | `deny` \| `no_decision` | `deny` |
| `offline` | `fail_open` \| `fail_closed` | `fail_closed` |

> **Note:** Hook CLI при недоступном демоне пишет в stderr `daemon not running` и выходит с кодом `1`; ветвления по `policy.offline` нет.

### async

| Ключ | Default |
|------|---------|
| `queue_capacity` | `1024` |
| `worker_limit` | `8` |
| `target_timeout` | `30s` |
| `on_overflow` | `drop` (`drop` \| `log`) |

Overflow всегда дропает job и увеличивает `async_dropped_count` в Status; режим `log` дополнительно пишет warn.

## CLI для конфига

| Команда | Роль |
|---------|------|
| `agentd config validate [--config] [--cwd]` | Offline parse + compile |
| `agentd config show [--merged] [--layer user\|project\|runtime] [--cwd]` | Просмотр слоёв |
| `agentd config patch --file DELTA.yaml` | Patch runtime (persist) |
| `agentd config record-decision …` | Upsert approval ([Approvals](./approvals.md)) |

## Reload

- Изменения user/project: fsnotify + debounce → re-merge
- `agentd daemon reload`: принудительный re-merge с диска
- Runtime patch / RecordDecision: in-memory + debounced flush

Status после успешного compile показывает `generation` и merged `fingerprint`.

См. также: [Guards](./guards.md), [Dispatch](./dispatch.md), [CLI](./cli.md).
