# Dispatch

> **Language:** [English](../en/dispatch.md) · [Русский](./dispatch.md)

Match маршрутов, sync/async pipelines, targets и таймауты. Детали: [DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine).

## Режимы

| Mode | Поведение |
|------|-----------|
| `sync_only` | Sync chain → decision |
| `async_only` | Enqueue async; нейтральный wire decision |
| `parallel` | Sync + async стартуют вместе; async не блокирует ответ |
| `after_sync` | Async после sync, с outcome |
| `sync_then_async` | Алиас `after_sync` |

Defaults по kind — в `dispatch_defaults:` (можно переопределить). Именованные `dispatch:` routes сверху вниз; первое совпадение побеждает.

## Targets

| Target | Sync | Async |
|--------|------|-------|
| `builtin` | Runner.Decide / guards | observe |
| `grpc` | да | да |
| `http` | — | да |
| `exec` | **нет (v1)** | да |
| `log` | — | да |
| `file` | — | да |

v1 **не** поддерживает sync `exec` JSON decisions (DESIGN §11 / §12).

## Поля route

```yaml
dispatch:
  - name: gate-and-audit
    match:
      kind: [tool.pre]
      provider: ["*"]
    mode: parallel
    sync_timeout: 25s          # optional cap
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

Match: `kind`, `provider`, `tools` (опциональные списки).

## Sync timeout budget

Эффективный budget:

`min(provider_timeout − 10%, route.sync_timeout)` если задан `sync_timeout`; иначе provider timeout минус 10%.

- Если у Invoke есть deadline → это provider timeout.
- Иначе defaults по kind: `tool.pre` / `prompt.submitted` → **30s**; остальные → **5s** (как install HookSpec).

Per-target gRPC `timeout` clamp’ится оставшимся budget контекста.

## Async queue

Defaults: capacity `1024`, workers `8`, `target_timeout` `30s`, `on_overflow: drop`.

Полная очередь → drop job, инкремент Status `async_dropped_count` (режим `log` ещё и warn). Async не должен блокировать sync response.

Compiled routes offline:

```bash
agentd dispatch routes --json
```

См. также: [Конфигурация](./configuration.md), [Эксплуатация](./operations.md).
