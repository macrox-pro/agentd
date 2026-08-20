# Troubleshooting

> **Language:** [English](../en/troubleshooting.md) · [Русский](./troubleshooting.md)

Типичные сбои и фактическое поведение agentd.

## Демон не запущен

`hook run|notify|serve` пишет в stderr `daemon not running` и выходит с кодом **1**.

```bash
agentd daemon start
agentd daemon status --json
```

`--socket` у edge и демона должен совпадать. Stale socket чистится только под start lock.

## Таймауты

Sync budget ≈ 90% provider timeout, опционально capped `route.sync_timeout` ([Dispatch](./dispatch.md)). CLI `--timeout 0` не задаёт deadline; демон берёт kind defaults (30s / 5s).

Превышение budget → cancel context; более широкий fail mode на стороне демона задаёт `policy.fail`.

## Циклы Ask

Тот же tool снова Ask → нет matching approval. Возьмите `approval_fingerprint=` из Ask и выполните `config record-decision` ([Approvals](./approvals.md)). Неверный `--session-id` не матчит session-scope.

## Async drops

Растёт `async_dropped_count` → очередь полна (`queue_capacity`). Увеличьте capacity/workers или уменьшите fan-out; `on_overflow: log` добавляет warn, но всё равно drop.

## Windows named pipe

Default pipe привязан к SID (`\\.\pipe\agentd-<SID>`). Другой user/session → dial failure, похожий на «daemon not running».

## Codex empty no-op

Для Codex/Kimi no-op — **пустой stdout**, exit 0, не `{}`. Пустой stdout не считать ошибкой.

## Поле offline

`policy.offline` парсится и хранится; hook edge **сейчас его не читает**. Недоступный демон → exit 1, как выше.

## Не входит в v1

Нет agent auth, transcript pipelines, Go plugins, hooks DSL, async retry storms и sync `exec` decisions ([DESIGN.md §11](../../DESIGN.md#11-non-goals-v1)).

См. также: [CLI](./cli.md), [Эксплуатация](./operations.md).
