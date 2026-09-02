# Маршрутизация

> **Language:** [English](../en/dispatch.md) · [Русский](./dispatch.md)

Как события проходят через синхронный путь (ответ агенту) и асинхронный путь (побочные эффекты).

**Синхронный путь** решает: разрешить, спросить или запретить — и формирует ответ на проводе. **Асинхронный путь** выполняет журнал, HTTP, exec и подобное, не задерживая этот ответ. Термины: [Глоссарий](./glossary.md).

Подробности: [DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine).

## Виды событий (`kind`)

Имена в YAML `match:` и `dispatch_defaults`. В JSON агента могут быть другие строки — agentd приводит их к этим.

| Вид (wire) | Примеры у агента | Когда |
|------------|------------------|--------|
| `tool.pre` | `PreToolUse` (Claude), `preToolUse` (Cursor) и т. д. | Перед инструментом |
| `prompt.submitted` | `PromptSubmitted`, `UserPromptSubmit` и т. д. | Пользователь отправил запрос |
| `agent.stop` | `Stop`, `SessionEnd` и т. д. | Завершение сессии |
| `tool.post` | `PostToolUse` и т. д. | После инструмента |
| `notification` | Codex `notify`, observe-only кадры | Только наблюдение |
| `other` | Прочее | По умолчанию только async |

## Режимы

| Режим | Поведение |
|-------|-----------|
| `sync_only` | Только синхронная цепочка → решение |
| `async_only` | Только async; нейтральный ответ агенту |
| `parallel` | Sync и async вместе; async не задерживает ответ |
| `after_sync` | Async после sync, с результатом sync |
| `sync_then_async` | То же, что `after_sync` |

Маршруты в `dispatch:` проверяются сверху вниз; побеждает первое совпадение.

## Цели

| Цель | Sync | Async |
|------|------|-------|
| `builtin` | проверки / решение | observe |
| `grpc` | да | да |
| `http` | — | да |
| `exec` | **нет** | да |
| `log` | — | да |
| `file` | — | да |

Синхронный `exec` с JSON-решением в текущем релизе не поддерживается (DESIGN §11).

## Справочник полей маршрута

```yaml
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
      provider: ["*"]
    mode: parallel
    sync_timeout: 25s          # необязательный потолок (см. ниже)
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

| Поле | Где | Смысл |
|------|-----|--------|
| `sync_timeout` | маршрут | Необязательный потолок бюджета sync (см. [Бюджет sync](#бюджет-времени-на-sync)) |
| `merge` | sync `grpc` | `first_conclusive` — побеждает первое ненейтральное решение |
| `on_error` | sync `grpc` | `fail_closed` (по умолчанию) или `fail_open` при сбое узла gRPC |
| `endpoint` | `grpc` | Адрес peer, например `unix:///path/to/peer.sock` |
| `stdin` | async `exec` | `raw` — полезная нагрузка события в stdin команды |

В `match` необязательные списки: `kind`, `provider`, `tools`.

## Значения по умолчанию по виду события

Встроенные (переопределяются в `dispatch_defaults:`):

```yaml
dispatch_defaults:
  tool.pre:         { mode: parallel }
  prompt.submitted: { mode: sync_only }
  agent.stop:       { mode: sync_then_async }
  tool.post:        { mode: parallel }
  notification:     { mode: async_only }
  other:            { mode: async_only }
```

## Бюджет времени на sync

Эффективный лимит:

`min(таймаут_агента − 10%, route.sync_timeout)`, если задан `sync_timeout`; иначе таймаут агента минус 10%.

- Если в запросе уже есть срок — от него считается таймаут агента.
- Иначе по виду: `tool.pre` / `prompt.submitted` → **30 с**; остальные → **5 с** (как при установке хуков).

Таймаут gRPC-цели ограничивается оставшимся временем контекста.

## Асинхронная очередь

По умолчанию: ёмкость `1024`, воркеров `8`, `target_timeout` `30s`, `on_overflow: drop`.

Очередь полна → задача отбрасывается, растёт `async_dropped_count` (при `log` ещё предупреждение). Async не должен блокировать синхронный ответ.

Без демона:

```bash
agentd dispatch routes --json
```

См. также: [Конфигурация](./configuration.md), [Эксплуатация](./operations.md), [Глоссарий](./glossary.md).
