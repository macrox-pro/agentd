# Справочник команд

> **Language:** [English](../en/cli.md) · [Русский](./cli.md)

Команды и флаги, как в пакете `cmd/`. Заметки по устройству CLI: [§6](../../DESIGN.md#6-cli-reference).

## Общие флаги (для всех подкоманд)

Флаги есть у каждой команды.

| Флаг | По умолчанию | Смысл |
|------|--------------|--------|
| `--config` | `~/.agentd.yaml` | Файл конфига пользователя |
| `--socket` | зависит от ОС | Как эта команда связывается с демоном |
| `-v` / `--verbose` | выкл. | Дополнительные сообщения в stderr (**не** в stdout хука) |

## version

Выводит версию CLI. Сначала ldflags goreleaser; иначе модульная версия из BuildInfo (`go install`); локальный devel — `dev` или `dev+<shortrev>`. Без обращения к демону.

Версия процесса демона — поле `version` в `agentd daemon status`.

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `version` | — | версия CLI в stdout |

## daemon (демон)

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `daemon start` | `--foreground`, `--log-level`, `--log-file`, `--metrics-listen` | По умолчанию уходит в фон; ждёт готовности; журнал в `agentd.log` в [каталоге состояния](./configuration.md#каталог-состояния); `--metrics-listen host:port` включает сбор метрик Prometheus для этого процесса |
| `daemon stop` | `--timeout` (`10s`) | Корректное завершение по gRPC, иначе SIGTERM |
| `daemon status` | `--json` | Состояние демона + блок `autostart` ([Эксплуатация](./operations.md#автозапуск-при-входе)) |
| `daemon reload` | — | Принудительно пересобрать конфиг с диска |
| `daemon enable` | — | Включает автозапуск при входе и стартует демон, если он не запущен. Может завершиться с ошибкой, хотя автозапуск уже включён — [Эксплуатация → Автозапуск](./operations.md#автозапуск-при-входе) |
| `daemon disable` | — | Отключает только автозапуск; **не** останавливает работающий демон |

## hook (хук)

Процесс, который запускает ИИ-агент. Читает событие, вызывает демон, пишет ответ. Разрешить / спросить / запретить решает **демон**, не эта команда. Если демон недоступен, применяется `policy.offline` из локального конфига.

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `hook run` | `--provider` (обязателен), `--argv-payload`, `--timeout` (`0` = не задан) | Событие во входном потоке (или в argv) |
| `hook notify` | `--provider`, `--timeout` | Уведомление Codex (JSON в argv) |
| `hook serve` | `--provider`, `--timeout` | Поток OpenCode (NDJSON); `--provider` только `opencode` |

Если демон недоступен: в stderr — `daemon not running`, затем `policy.offline`. По умолчанию `fail_open` → код 0 и ответ «ничего не менять», агент продолжает работу. `fail_closed` → код **1**. На этом пути не писать отладку в stdout.

### agenthooks (скрытая команда)

`agentd install` прописывает в настройках агента `agentd agenthooks …`, а не `hook`. Скрытые подкоманды делают то же самое, чтобы эти файлы работали. В документации и при ручной настройке удобнее `hook run` / `hook serve` / `hook notify`.

Установка пишет `agentd agenthooks run|notify|serve --provider=…`. Флаги те же, что у `hook …`. У `agenthooks serve` значение `--provider` по умолчанию — `opencode`.

## config (конфиг)

Команды `enable` / `disable` / `get` пишут только в YAML **пользователя или проекта**, не во временный слой. Временные правки — `config patch`. **`daemon enable`** — автозапуск при входе, не переключатель конфига.

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `config validate` | `--cwd` | Offline parse + compile |
| `config show` | `--merged`, `--layer user\|project\|runtime`, `--cwd` | Просмотр слоёв |
| `config enable FEATURE` | `--scope user\|project`, `--cwd` | Записать `enabled: true` (создаёт файл пользователя, если его нет) |
| `config disable FEATURE` | `--scope user\|project`, `--cwd` | Записать `enabled: false` |
| `config get FEATURE` | `--cwd` | Фактическое on/off и слой-победитель (`default` \| `user` \| `project`); временный слой не учитывается |
| `config patch` | `--file` (обязателен) | Правка временного слоя (нужен демон) |
| `config record-decision` | `--fingerprint` (обязателен), `--scope` (по умолчанию `project`), `--project-root`, `--session-id`, `--expires-at` (RFC3339) | Upsert approval |

**Features:** `trajectory`, `trajectory-raw`, `guard-shell`, `guard-mcp`, `guard-paths`. Scope по умолчанию: `user` для trajectory; `project` для guards. Offline; работающий демон подхватывает через fsnotify.

| Feature | YAML path | Default scope |
|---------|-----------|---------------|
| `trajectory` | `trajectory.enabled` | `user` |
| `trajectory-raw` | `trajectory.include_raw` | `user` |
| `trajectory-statistics` | `trajectory.statistics` | `user` |
| `guard-shell` | `guards.shell.enabled` | `project` |
| `guard-mcp` | `guards.mcp.enabled` | `project` |
| `guard-paths` | `guards.paths.enabled` | `project` |

**Вывод:** `get` печатает `FEATURE: on|off (SOURCE)`, где `SOURCE` — `default`, `user` или `project`. `enable` / `disable` — `FEATURE: enabled|disabled (SCOPE PATH)`; повтор при том же состоянии — `already enabled|disabled`, exit 0.

```bash
# trajectory (scope по умолчанию user → ~/.agentd.yaml)
agentd config enable trajectory
# trajectory: enabled (user /home/you/.agentd.yaml)

agentd config get trajectory
# trajectory: on (user)

# guards (scope по умолчанию project → .agentd.yaml под --cwd)
cd /path/to/repo
agentd config enable guard-shell
# guard-shell: enabled (project /path/to/repo/.agentd.yaml)
```

**Асимметрия project-пути:** `config get --cwd DIR` ищет `.agentd.yaml` вверх от `DIR` (как при merge хуков). `config enable|disable` с project scope пишет **только** в `DIR/.agentd.yaml` — родительский конфиг репозитория не обновляется. Для project-guards запускайте из корня репо (или `--cwd` на корень).

**Не через CLI:** `guards.secrets` и прочие не-boolean поля — правка YAML или `config show` / `config patch`. Curated toggles — только пять features выше.

## doctor (диагностика)

Список ИИ-агентов на этой машине и совпадают ли их файлы хуков с agentd. **Только чтение** — настройки агента не меняет.

| Флаг | По умолчанию | Смысл |
|------|--------------|--------|
| `--json` | выкл. | Отчёт в JSON |
| `--cwd` | текущий каталог | Корень проекта для поиска |

Если сокет демона отвечает (`--socket`), в выводе `daemon: reachable` или `"DaemonReachable": true`. Недоступный демон — не ошибка, код выхода 0.

```bash
agentd doctor
agentd doctor --json
agentd doctor --cwd /path/to/repo
```

## install (установка хуков)

Записывает настройки агента, чтобы он вызывал agentd. Демон при этом не запускается.

| Флаг | По умолчанию |
|------|--------------|
| `--provider` | обязателен, если нет `--all-detected` |
| `--scope` | `project` (также `user`, `plugin`) |
| `--global` | выкл. — то же, что `--scope=user` |
| `--dir` | `project`: текущий каталог (Codex: `./.codex`); `user`: домашний каталог агента (например `~/.cursor`); `plugin`: обязателен |
| `--all-detected` | выкл. — агенты с каталогом настроек; без `--yes` только план |
| `--yes` | выкл. — действительно записать файлы для `--all-detected` |
| `--dry-run` | выкл. — показать план, ничего не писать (один `--provider` или вместе с `--all-detected --yes`) |

`--global` нельзя сочетать с `--scope`, отличным от `user`. При успехе печатает агента, область, корень установки и по каждому файлу `create` / `update` / `unchanged` с абсолютными путями.

```bash
agentd install --provider=claude-code --scope=project
agentd install --all-detected              # только план
agentd install --all-detected --yes        # записать файлы
agentd install --provider=cursor --dry-run
```

В интерактивном терминале `agentd install` без флагов открывает короткий мастер ([`setup`](#setup-мастер-настройки)). В скриптах и CI нужны `--provider` или `--all-detected`.

## setup (мастер настройки)

Сценарий в настоящем терминале: найти агентов, выбрать куда ставить, показать план, по желанию записать файлы.

| Флаг | По умолчанию | Смысл |
|------|--------------|--------|
| `--yes` | выкл. | Записать без повторного подтверждения |
| `--dry-run` | выкл. | Показать план, ничего не писать |
| `--cwd` | текущий каталог | Корень проекта для поиска |

`AGENTD_NO_TUI=1` или `CI=true` отключают мастер (команда подскажет, как поставить хуки без терминала).

```bash
agentd setup
agentd setup --yes
```

## dispatch (маршруты)

| Команда | Флаги |
|---------|-------|
| `dispatch routes` | `--json`, `--cwd` |

Компиляция маршрутов из defaults ⊕ user ⊕ (опционально) project **без** работающего демона.

## session (журнал trajectory)

Просмотр и экспорт JSONL ([журнале сессий](./trajectory.md)). Offline — читает `sessions/` в [каталоге состояния](./configuration.md#каталог-состояния). **Исключения:** `session subscribe` и `trajectory stats` требуют запущенный daemon.

## trajectory stats

```bash
agentd trajectory stats [--provider ID] [--json]
```

Счётчики за время жизни демона с момента старта процесса (`since` в выводе). Нужны `trajectory.enabled` и `trajectory.statistics`. См. [Trajectory § Статистика демона](./trajectory.md#статистика-демона).

| Команда | Флаги |
|---------|-------|
| `session list` | `--provider`, `--json` (в `--json` есть `importer_status`) |
| `session show SESSION_ID` | `--provider` (обязателен), `--json` |
| `session export` | `--provider`, `--session`, `--out` |
| `session search` | `--provider`, `--session`, `--kind` (повтор), `--source`, `--query`, `--limit`, `--json` |
| `session import` | `--provider` (обязателен), `--session`, `--path`, `--out`, `--dry-run`, `--json` |
| `session replay` | `--policy` (обязателен), `--provider`, `--session`, `--seq`, `--json` |
| `session fork` | `--provider`, `--session`, `--new-session`, `--at-seq`, `--json` |
| `session stats` | `SESSION_ID`, `--provider` (обязателен), `--json` |
| `session subscribe` | `--provider`, `--session`, `--source`, `--json` (live; нужен daemon) |

`session search` сканирует JSONL построчно (O(объём); без индекса). `session import`: Claude Code и Codex — `supported`; Cursor — `partial` (лучше `--path`); остальные — явный `none`. `session replay --policy` требует `include_raw` при записи. `session fork` — только аудит (исходник неизменяем). `session subscribe` — live с момента подключения; история через show/export.

### session import `--out`

**По умолчанию:** дописывает распарсенные transcript-события в ledger в [каталоге состояния](./configuration.md#каталог-состояния); summary на **stdout**.

| | `session export` | `session import` |
|--|------------------|------------------|
| **Вывод по умолчанию** | stdout (ledger JSONL) | ledger на диске (без потока событий) |
| **Без `--out`** | данные → stdout | данные → `sessions/<provider>/<id>.jsonl` |
| **`--out PATH`** | файл вместо stdout | файл вместо ledger (только parse) |
| **`--out -`** | N/A (stdout уже по умолчанию) | stdout вместо ledger (только parse) |
| **Summary при потоке данных** | N/A (только данные) | stderr при `--out` |

- **`--out -`:** JSONL на **stdout**; ledger и import checkpoint **не** обновляются; summary на **stderr** (`--json` — JSON summary на stderr).
- **`--out PATH`:** как `-`, но в файл (truncate/create; файл `0o600`, каталог `0o700`).
- **`--dry-run`:** только summary, без записи в ledger (как раньше). **`--out` ≠ `--dry-run`** — `dry_run` в JSON summary true только при флаге `--dry-run`. Комбинация `--dry-run --out -` отдаёт события и помечает summary как dry-run.
- **Формат:** один JSON-объект на строку — та же схема, что у `session export` / ledger JSONL.
- **Инкрементальный импорт:** читает `<session_id>.import.json` для `startIndex`; при `--out` sidecar не обновляется.
- **Ошибка записи:** ненулевой exit; выходной файл может быть усечён.

Пример pipe:

```bash
agentd session import --provider claude-code --path /path/to/session.jsonl --out - | jq -c 'select(.type=="transcript/message")'
```

См. также: [Быстрый старт](./getting-started.md), [Агенты](./providers.md).
