# Справочник CLI

> **Language:** [English](../en/cli.md) · [Русский](./cli.md)

Команды и флаги, как в пакете `cmd/`. Развёрнутое описание в DESIGN: [§6](../../DESIGN.md#6-cli-reference).

## Общие флаги (для всех подкоманд)

| Флаг | По умолчанию | Смысл |
|------|--------------|--------|
| `--config` | `~/.agentd.yaml` | Путь пользовательского конфига |
| `--socket` | зависит от ОС | Точка IPC (сокет или именованный канал) |
| `-v` / `--verbose` | выкл. | Доп. сообщения в stderr (**не** в stdout хука) |

## version

Выводит версию CLI (`dev` без ldflags/тега релиза). Без обращения к демону.

Версия процесса демона — поле `version` в `agentd daemon status`.

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `version` | — | версия CLI в stdout |

## daemon (демон)

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `daemon start` | `--foreground`, `--log-level`, `--log-file` | По умолчанию уходит в фон; ждёт Health; логи в файл state-dir |
| `daemon stop` | `--timeout` (`10s`) | Корректное завершение по gRPC, иначе SIGTERM |
| `daemon status` | `--json` | Статус работающего демона, включая версию процесса (`version`; [Эксплуатация](./operations.md)) |
| `daemon reload` | — | Принудительно пересобрать конфиг с диска |

## hook (клиент хука)

Тонкий край: разобрать вход → вызвать демону `Invoke` → закодировать ответ. Политик в CLI нет.

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `hook run` | `--provider` (обязателен), `--argv-payload`, `--timeout` (`0` = не задан) | Хук из stdin (или argv) |
| `hook notify` | `--provider`, `--timeout` | Путь notify у Codex (JSON в argv) |
| `hook serve` | `--provider`, `--timeout` | Долгий мост OpenCode (NDJSON); только `opencode` |

Если не удалось связаться с демоном или вызвать `Invoke`: в stderr — `daemon not running`, код выхода **1**. На пути хука не писать отладочный вывод в stdout.

### agenthooks (скрытая команда)

**Зачем:** `agentd install` вызывает [agenthooks/install](https://github.com/speakeasy-api/agenthooks), который прописывает в настройках агента `agentd agenthooks …`, а не `hook`. agentd регистрирует скрытые подкоманды с тем же поведением, чтобы сгенерированные конфиги работали без правок. В документации и при ручной настройке удобнее `hook run` / `hook serve` / `hook notify` — те же флаги и тот же путь через `hookedge` (`cmd/hook.go`).

`install` прописывает в конфиг агента `agentd agenthooks run|notify|serve --provider=…`. Поведение совпадает с `hook …`. У `agenthooks serve` значение `--provider` по умолчанию — `opencode`.

## config (конфиг)

| Команда | Флаги |
|---------|-------|
| `config validate` | `--cwd` |
| `config show` | `--merged`, `--layer user\|project\|runtime`, `--cwd` |
| `config patch` | `--file` (обязателен) |
| `config record-decision` | `--fingerprint` (обязателен), `--scope` (по умолчанию `project`), `--project-root`, `--session-id`, `--expires-at` (RFC3339) |

## install (установка хуков в агент)

| Флаг | По умолчанию |
|------|--------------|
| `--provider` | обязателен |
| `--scope` | `project` (также `user`, `plugin`) |
| `--global` | false — то же, что `--scope=user` |
| `--dir` | `scope=project`: cwd (codex: `./.codex`); `scope=user`: home агента (напр. `~/.cursor`); `scope=plugin`: обязателен |

`--global` конфликтует с явным `--scope`, отличным от `user`. При успехе печатает provider, scope, корень установки и по каждому файлу `create` / `update` / `unchanged` с абсолютными путями.

## dispatch (маршруты)

| Команда | Флаги |
|---------|-------|
| `dispatch routes` | `--json`, `--cwd` |

Компиляция маршрутов из defaults ⊕ user ⊕ (опционально) project **без** работающего демона.

## session (журнал trajectory)

Просмотр и экспорт JSONL ([Trajectory](./trajectory.md)). Offline — читает `$XDG_STATE_HOME/agentd/sessions/`. **Исключение:** `session subscribe` требует запущенный daemon.

| Команда | Флаги |
|---------|-------|
| `session list` | `--provider`, `--json` (в `--json` есть `importer_status`) |
| `session show SESSION_ID` | `--provider` (обязателен), `--json` |
| `session export` | `--provider`, `--session`, `--out` |
| `session search` | `--provider`, `--session`, `--kind` (повтор), `--source`, `--query`, `--limit`, `--json` |
| `session import` | `--provider` (обязателен), `--session`, `--path`, `--dry-run`, `--json` |
| `session replay` | `--policy` (обязателен), `--provider`, `--session`, `--seq`, `--json` |
| `session fork` | `--provider`, `--session`, `--new-session`, `--at-seq`, `--json` |
| `session subscribe` | `--provider`, `--session`, `--source`, `--json` (live; нужен daemon) |

`session search` сканирует JSONL построчно (O(объём); без индекса). `session import`: Claude Code и Codex — `supported`; Cursor — `partial` (лучше `--path`); остальные — явный `none`. `session replay --policy` требует `include_raw` при записи. `session fork` — только аудит (исходник неизменяем). `session subscribe` — live с момента подключения; история через show/export.

См. также: [Быстрый старт](./getting-started.md), [Провайдеры](./providers.md).
