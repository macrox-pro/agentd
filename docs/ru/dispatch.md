# Маршрутизация (dispatch)

> **Language:** [English](../en/dispatch.md) · [Русский](./dispatch.md)

Как события попадают на маршруты, синхронный и асинхронный конвейеры, типы целей и таймауты. Подробности: [DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine).

## Режимы

| Режим | Поведение |
|-------|-----------|
| `sync_only` | Только синхронная цепочка → решение для агента |
| `async_only` | Только постановка в асинхронную очередь; на wire — нейтральный ответ |
| `parallel` | Sync и async стартуют вместе; async **не** держит ответ агенту |
| `after_sync` | Async после sync, с результатом sync |
| `sync_then_async` | То же, что `after_sync` |

Значения по умолчанию по виду события (`kind`) задаются в `dispatch_defaults:` (можно переопределить). Именованные маршруты в `dispatch:` проверяются сверху вниз; побеждает первое совпадение.

## Типы целей (targets)

| Цель | Sync | Async | Назначение |
|------|------|-------|------------|
| `builtin` | да | observe | Встроенные охранники / наблюдение |
| `grpc` | да | да | Пересылка другому демону/сервису |
| `http` | — | да | HTTP-уведомление |
| `exec` | **нет в v1** | да | Запуск внешней команды |
| `log` | — | да | Структурированный лог |
| `file` | — | да | Дозапись JSONL |

В v1 **нет** синхронного `exec` с JSON-решением (DESIGN §11).

## Поля маршрута

```yaml
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
      provider: ["*"]
    mode: parallel
    sync_timeout: 25s          # необязательный потолок для sync
    sync:
      - target: builtin
        guards: [secrets, shell]
      - target: grpc
        endpoint: unix:///path/to/peer.sock
        timeout: 2s
        on_error: fail_closed  # или fail_open
        merge: first_conclusive
    async:
      - target: log
        level: info
      - target: exec
        command: ["notify", "--"]
        stdin: raw
```

Условие `match`: списки `kind`, `provider`, `tools` (все необязательны).

## Бюджет времени на sync

Эффективный лимит:

`min(таймаут_провайдера − 10%, route.sync_timeout)`, если задан `sync_timeout`; иначе таймаут провайдера минус 10%.

- Если в запросе `Invoke` уже есть deadline — от него считается таймаут провайдера.
- Иначе по виду события: `tool.pre` / `prompt.submitted` → **30 s**; остальные → **5 s** (как в спецификации хуков при `install`).

Таймаут отдельной gRPC-цели ограничивается оставшимся временем контекста.

## Асинхронная очередь

По умолчанию: ёмкость `1024`, воркеров `8`, `target_timeout` `30s`, `on_overflow: drop`.

Очередь полна → задача отбрасывается, в статусе растёт `async_dropped_count` (при `log` ещё предупреждение в лог). Асинхронная обработка не должна блокировать синхронный ответ агенту.

Посмотреть скомпилированные маршруты без демона:

```bash
agentd dispatch routes --json
```

См. также: [Конфигурация](./configuration.md), [Эксплуатация](./operations.md).
