# Справочник CLI

> **Language:** [English](../en/cli.md) · [Русский](./cli.md)

Команды и флаги, как в пакете `cmd/`. Архитектурные заметки: [§6](../../DESIGN.md#6-cli-reference).

## Общие флаги (для всех подкоманд)

| Флаг | По умолчанию | Смысл |
|------|--------------|--------|
| `--config` | `~/.agentd.yaml` | Путь пользовательского конфига |
| `--socket` | зависит от ОС | Точка IPC (сокет или именованный канал) |
| `-v` / `--verbose` | выкл. | Доп. сообщения в stderr (**не** в stdout хука) |

## version

Выводит версию CLI. Сначала ldflags goreleaser; иначе модульная версия из BuildInfo (`go install`); локальный devel — `dev` или `dev+<shortrev>`. Без обращения к демону.

Версия процесса демона — поле `version` в `agentd daemon status`.

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `version` | — | версия CLI в stdout |

## daemon (демон)

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `daemon start` | `--foreground`, `--log-level`, `--log-file`, `--metrics-listen` | По умолчанию уходит в фон; ждёт Health; логи в `agentd.log` в [state directory](./configuration.md#state-directory); `--metrics-listen host:port` включает Prometheus scrape для этого процесса |
| `daemon stop` | `--timeout` (`10s`) | Корректное завершение по gRPC, иначе SIGTERM |
| `daemon status` | `--json` | Состояние демона + блок `autostart` ([Эксплуатация](./operations.md#автозапуск-при-входе)) |
| `daemon reload` | — | Принудительно пересобрать конфиг с диска |
| `daemon enable` | — | Включает автозапуск при входе и стартует демон, если он не запущен. Может завершиться с ошибкой, хотя автозапуск уже включён — [Эксплуатация → Автозапуск](./operations.md#автозапуск-при-входе) |
| `daemon disable` | — | Отключает только автозапуск; **не** останавливает работающий демон |

## hook (клиент хука)

Тонкий край: разобрать вход → вызвать демону `Invoke` → закодировать ответ. Полный Decide/guards остаются в демоне. Если демон недоступен, край применяет `policy.offline` из локального конфига.

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `hook run` | `--provider` (обязателен), `--argv-payload`, `--timeout` (`0` = не задан) | Хук из stdin (или argv) |
| `hook notify` | `--provider`, `--timeout` | Путь notify у Codex (JSON в argv) |
| `hook serve` | `--provider`, `--timeout` | Долгий мост OpenCode (NDJSON); только `opencode` |

Если не удалось связаться с демоном или вызвать `Invoke`: в stderr — `daemon not running`, затем `policy.offline` (по умолчанию `fail_open` → код 0 / нейтральный wire; `fail_closed` → код **1**). На пути хука не писать отладочный вывод в stdout.

### agenthooks (скрытая команда)

**Зачем:** `agentd install` вызывает [agenthooks/install](https://github.com/speakeasy-api/agenthooks), который прописывает в настройках агента `agentd agenthooks …`, а не `hook`. agentd регистрирует скрытые подкоманды с тем же поведением, чтобы сгенерированные конфиги работали без правок. В документации и при ручной настройке удобнее `hook run` / `hook serve` / `hook notify` — те же флаги и тот же путь через `hookedge` (`cmd/hook.go`).

`install` прописывает в конфиг агента `agentd agenthooks run|notify|serve --provider=…`. Поведение совпадает с `hook …`. У `agenthooks serve` значение `--provider` по умолчанию — `opencode`.

## config (конфиг)

Curated-переключатели (`enable` / `disable` / `get`) пишут только в **user или project** YAML — не в runtime overlay. Для временных override — `config patch`. **`daemon enable`** — автозапуск при входе, не toggle конфига.

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `config validate` | `--cwd` | Offline parse + compile |
| `config show` | `--merged`, `--layer user\|project\|runtime`, `--cwd` | Просмотр слоёв |
| `config enable FEATURE` | `--scope user\|project`, `--cwd` | Записать `enabled: true` (создаёт user bootstrap при отсутствии файла) |
| `config disable FEATURE` | `--scope user\|project`, `--cwd` | Записать `enabled: false` |
| `config get FEATURE` | `--cwd` | Эффективное on/off + слой-победитель (`default` \| `user` \| `project`); runtime не учитывается |
| `config patch` | `--file` (обязателен) | Patch runtime overlay (нужен демон) |
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

Только чтение: обнаруженные агенты и статус установки хуков. Конфиг агента не меняется.

| Флаг | По умолчанию | Смысл |
|------|--------------|--------|
| `--json` | выкл. | JSON-отчёт (`DoctorReport`) |
| `--cwd` | текущий каталог | Рабочий каталог для discovery |

Если сокет демона доступен (`--socket`), в отчёте `daemon: reachable` (текст) или `"DaemonReachable": true` (JSON). Недоступный демон — не ошибка, exit 0.

```bash
agentd doctor
agentd doctor --json
agentd doctor --cwd /path/to/repo
```

## install (установка хуков в агент)

| Флаг | По умолчанию |
|------|--------------|
| `--provider` | обязателен, если нет `--all-detected` |
| `--scope` | `project` (также `user`, `plugin`) |
| `--global` | false — то же, что `--scope=user` |
| `--dir` | `scope=project`: cwd (codex: `./.codex`); `scope=user`: home агента (напр. `~/.cursor`); `scope=plugin`: обязателен |
| `--all-detected` | выкл. — high-confidence агенты; только план без `--yes` |
| `--yes` | выкл. — применить `--all-detected` (обязателен для записи) |
| `--dry-run` | выкл. — план без записи (один `--provider` или с `--all-detected --yes`) |

`--global` конфликтует с явным `--scope`, отличным от `user`. При успехе печатает provider, scope, корень установки и по каждому файлу `create` / `update` / `unchanged` с абсолютными путями.

```bash
agentd install --provider=claude-code --scope=project
agentd install --all-detected              # только план
agentd install --all-detected --yes        # установка high-confidence
agentd install --provider=cursor --dry-run
```

На TTY голый `agentd install` (без флагов) открывает интерактивный мастер ([`setup`](#setup-мастер-настройки)). В non-TTY нужны `--provider` или `--all-detected`.

## setup (мастер настройки)

Интерактивный мастер (TTY): discovery, выбор целей, превью плана, установка.

| Флаг | По умолчанию | Смысл |
|------|--------------|--------|
| `--yes` | выкл. | Без подтверждения, сразу применить |
| `--dry-run` | выкл. | Только превью |
| `--cwd` | текущий каталог | Каталог для discovery |

`AGENTD_NO_TUI=1` или `CI=true` отключают TUI (подсказка для non-interactive).

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

Просмотр и экспорт JSONL ([Trajectory](./trajectory.md)). Offline — читает `sessions/` в [state directory](./configuration.md#state-directory). **Исключения:** `session subscribe` и `trajectory stats` требуют запущенный daemon.

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

**По умолчанию:** дописывает распарсенные transcript-события в ledger в [state directory](./configuration.md#state-directory); summary на **stdout**.

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

См. также: [Быстрый старт](./getting-started.md), [Провайдеры](./providers.md).
