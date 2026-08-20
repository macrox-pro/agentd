# Эксплуатация

> **Language:** [English](../en/operations.md) · [Русский](./operations.md)

Day-2 управление user-демоном.

## Один демон на пользователя

Lock + socket не дают второму start. IPC переопределяется `--socket`.

## Status

```bash
agentd daemon status
agentd daemon status --json
```

JSON при `running`:

| Поле | Смысл |
|------|--------|
| `running` | bool |
| `socket` | path / pipe |
| `version` | версия сборки (`dev`, если без ldflags/tag) |
| `started_at` | RFC3339 UTC |
| `generation` | generation конфига |
| `fingerprint` | fingerprint merged-конфига |
| `async_queue_depth` | jobs в async queue |
| `async_dropped_count` | overflow drops (монотонно) |
| `compiled_route_count` | число routes в snapshot |

Human: `agentd: running (version …, generation …)` — без depth/drops (нужен `--json`).

## Reload / stop

```bash
agentd daemon reload
agentd daemon stop --timeout 10s
```

Stop: drain sync, затем async (до timeout), снятие socket/PID.

## Логирование

Hook path: без debug в stdout. Демон / async `log` targets — structured logs (stderr / sinks).

См. также: [Dispatch](./dispatch.md) (overflow), [Конфигурация](./configuration.md) (reload).
