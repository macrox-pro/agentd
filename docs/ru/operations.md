# Эксплуатация

> **Language:** [English](../en/operations.md) · [Русский](./operations.md)

Повседневное управление демоном после установки.

## Один демон на пользователя

Файловая блокировка и сокет не дают запустить второй экземпляр. Точку IPC можно сменить флагом `--socket`.

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
| `generation` | поколение конфига после слияния |
| `fingerprint` | отпечаток слитого конфига |
| `async_queue_depth` | сколько задач ждёт в асинхронной очереди |
| `async_dropped_count` | сколько задач отброшено из‑за переполнения (растёт только вверх) |
| `trajectory_dropped_count` | overflow очереди trajectory (если ledger включён) |
| `compiled_route_count` | число скомпилированных маршрутов |

В `--json` всегда есть объект `autostart` (даже при `running: false`):

| Поле | Смысл |
|------|--------|
| `autostart.enabled` | Автозапуск при входе зарегистрирован в ОС |
| `autostart.backend` | `systemd`, `launchd` или `schtasks` |
| `autostart.manifest_path` | Путь к unit/plist, если применимо |
| `autostart.registered_exe` | Путь к бинарнику для старта при входе |
| `autostart.stale` | `true`, если `registered_exe` не совпадает с этим CLI (после обновления снова выполните `daemon enable`) |

Человекочитаемая строка: `agentd: running (version …, generation …)` — без autostart (нужен `--json`).

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

После установки нового бинарника (например `go install …@latest`) снова выполните `agentd daemon enable`. Смотрите `autostart.stale` в `daemon status --json`.

### `disable` и `stop`

`daemon disable` снимает только автозапуск при входе. Работающий процесс **не** останавливается. Для остановки используйте `daemon stop`.

## Перезагрузка и остановка

```bash
agentd daemon reload
agentd daemon stop --timeout 10s
```

При остановке: сначала дождаться синхронных запросов, затем асинхронной очереди (не дольше `--timeout`), потом снять сокет и PID-файл.

## Логирование

На пути хука не писать отладку в stdout. Демон дописывает операционные логи в `agentd.log` в [state directory](./configuration.md#state-directory); `agentd daemon start --foreground` также дублирует в stderr. Асинхронная цель dispatch `target: log` использует тот же slog-логгер.

См. также: [Маршрутизация](./dispatch.md) (переполнение очереди), [Конфигурация](./configuration.md) (перезагрузка).
