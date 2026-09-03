# Эксплуатация

> **Language:** [English](../en/operations.md) · [Русский](./operations.md)

Повседневное управление демоном после установки.

## Один демон на пользователя

Файловая блокировка и сокет не дают запустить второй экземпляр. Сокет можно сменить флагом `--socket`. PID пишется в `agentd.pid` рядом с сокетом; в Windows сокет — именованный канал без своего каталога, поэтому оба файла лежат в [каталоге состояния](./configuration.md#каталог-состояния).

## Статус

```bash
agentd daemon status
agentd daemon status --json
```

Поля JSON, когда демон запущен:

| Поле | Смысл |
|------|--------|
| `running` | работает ли демон |
| `socket` | путь к сокету или имя pipe |
| `version` | версия демона (`dev` / `dev+rev` / модульная версия — как собран процесс) |
| `started_at` | время старта (RFC3339, UTC) |
| `generation` | Сколько раз конфиг загружался с момента старта (растёт при перезагрузке) |
| `fingerprint` | Контрольная сумма собранного конфига — меняется, когда меняется политика |
| `async_queue_depth` | Сколько фоновых задач ждёт очереди |
| `async_dropped_count` | Сколько задач отброшено из‑за переполнения (до перезапуска только растёт) |
| `trajectory_dropped_count` | Сколько записей журнала сессий отброшено (если журнал включён) |
| `compiled_route_count` | Сколько маршрутов сейчас действует |
| `metrics_listen` | Адрес метрик Prometheus, если они включены; иначе пусто |

В `--json` всегда есть объект `autostart` (даже при `running: false`):

| Поле | Смысл |
|------|--------|
| `autostart.enabled` | Автозапуск при входе зарегистрирован в ОС |
| `autostart.backend` | `systemd`, `launchd` или `schtasks` |
| `autostart.manifest_path` | Путь к unit/plist, если применимо |
| `autostart.registered_exe` | Путь к бинарнику для старта при входе |
| `autostart.stale` | `true`, если `registered_exe` не совпадает с этим CLI (после обновления снова выполните `daemon enable`) |

Человекочитаемая строка: `agentd: running (version …, generation …)` — без autostart (нужен `--json`).

## Метрики Prometheus

Включение в user config или при старте:

```yaml
metrics:
  enabled: true
  listen: "127.0.0.1:2112"
```

Или разовый override:

```bash
agentd daemon start --metrics-listen 127.0.0.1:2112
```

При работе демона `agentd daemon status --json` содержит `metrics_listen` (пусто, если выключено). Адрес сбора: `http://<metrics_listen>/metrics`.

Пример job Prometheus:

```yaml
scrape_configs:
  - job_name: agentd
    static_configs:
      - targets: ["127.0.0.1:2112"]
```

Адрес metrics фиксируется при **старте демона**. После изменения `metrics` в YAML выполните **`agentd daemon stop`**, затем **`agentd daemon start`** — `daemon reload` обновляет снимок конфига, но не перепривязывает HTTP listener.

По умолчанию bind — loopback (`127.0.0.1`). Публикация на `0.0.0.0` возможна, но не рекомендуется на общих машинах.

## Автозапуск при входе

Выполните `agentd daemon enable` один раз, если agentd должен стартовать автоматически при входе в систему. `agentd daemon disable` снимает регистрацию. Отключение автозапуска **не** останавливает уже работающий демон.

```bash
agentd daemon enable
agentd daemon disable
```

### Что делает `enable`

`daemon enable` регистрирует agentd в ОС (systemd user unit, launchd LaunchAgent или планировщик Windows), затем пытается сразу запустить демон, если он ещё не работает.

### Команда завершилась с ошибкой, но автозапуск мог уже быть включён

Если `daemon enable` завершился с ошибкой, **автозапуск мог уже быть включён**, даже если демон сейчас не запущен. Обычно это значит: при входе в систему agentd будет стартовать сам, но немедленный запуск не удался — часто из‑за ошибки в `~/.agentd.yaml`.

**Что делать:**

1. Проверьте `agentd daemon status --json` (`autostart.enabled`).
2. Исправьте проблему из текста ошибки (`agentd config validate --config ~/.agentd.yaml`).
3. Выполните `agentd daemon start` или перелогиньтесь — повторный `enable` не нужен.
4. Чтобы отключить только автозапуск: `agentd daemon disable`.

### После обновления agentd

После установки нового бинарника (например `go install …@latest`) снова выполните `agentd daemon enable`. Смотрите `autostart.stale` в `daemon status --json`. Снова выполните `agentd install`, чтобы агенты подписались на новые виды хуков.

### `disable` и `stop`

`daemon disable` снимает только автозапуск при входе. Работающий процесс **не** останавливается. Для остановки используйте `daemon stop`.

## Перезагрузка и остановка

```bash
agentd daemon reload
agentd daemon stop --timeout 10s
```

При остановке: сначала дождаться синхронных запросов, затем асинхронной очереди (не дольше `--timeout`), потом снять сокет, PID-файл (`agentd.pid`) и блокировку. Удаление устаревшего сокета при следующем старте выполняется только под `agentd.lock`.

## Статистика журнала сессий

Журнал сессий и статистика включены по умолчанию. Явное включение (запись в пользовательский слой) опционально:

```bash
agentd config get trajectory
agentd config get trajectory-statistics
agentd trajectory stats [--provider ID] [--json]
agentd session stats SESSION_ID --provider ID [--json]
```

`trajectory stats` читает счётчики демона в памяти (`TrajectoryService.Statistics`) и требует запущенный демон; счётчики сбрасываются при перезапуске. `session stats` сканирует локальный JSONL-журнал и демон не нужен. Оба требуют `trajectory.enabled` и `trajectory.statistics`. См. [Журнал сессий → Статистика демона](./trajectory.md#статистика-демона).

## Логирование

На пути хука не писать отладку в stdout. Демон дописывает операционные логи в `agentd.log` в [каталоге состояния](./configuration.md#каталог-состояния); `agentd daemon start --foreground` также дублирует в stderr. Асинхронная цель dispatch `target: log` пишет в тот же журнал работы демона.

См. также: [Маршрутизация](./dispatch.md) (переполнение очереди), [Конфигурация](./configuration.md) (перезагрузка).
