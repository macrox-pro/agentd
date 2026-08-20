# Зачем нужен agentd

> **Language:** [English](../en/why.md) · [Русский](./why.md)

agentd — **user-level демон** между coding agents и вашей hook-логикой: одна поверхность политик, корректный wire I/O под провайдера, sync-решения и async side effects без переписывания диалекта hooks каждого агента.

Codecs и install — через [agenthooks](https://github.com/speakeasy-api/agenthooks).

## Какие боли закрывает

| Боль | Без agentd | С agentd |
|------|------------|----------|
| **N провайдеров, N диалектов** | Отдельные скрипты/таймауты/exit codes под Claude, Cursor, Codex, Gemini, OpenCode, Kimi | Тонкий `hook run` / `serve` / `notify`; политика в демоне |
| **Cold start на каждый tool call** | Тяжёлая логика в каждом hook-процессе → latency и срывы по timeout | Долгоживущий демон; hot path — in-memory snapshot (без disk I/O на Invoke) |
| **Guards смешаны с audit** | Один процесс и блокирует агента, и шлёт webhooks/metrics | Sync (Ask/Deny) vs async queue (log/http/exec/…) — async не блокирует wire response |
| **Дрейф политик по репозиториям** | Copy-paste hook-конфигов; сложно «одобрить раз» / временно заблокировать | Слои YAML (user ⊕ project ⊕ runtime); approvals + temporary blocks с persist |
| **Слепая зона ops** | Нет единой точки «жив ли gate» и давления на очередь | `daemon status --json`: generation, fingerprint, `async_queue_depth`, `async_dropped_count` |

## Чем не является (v1)

Не продукт agent auth, не transcript pipeline, не plugin runtime и не общий hooks DSL. Targets — декларативный YAML; exec остаётся **async-only**. См. [DESIGN.md §11](../../DESIGN.md#11-non-goals-v1).

## Для кого

Инженеры, у которых coding-agent hooks в реальных workflow и нужна **одна** точка для secrets/shell/MCP/path guards, опциональных forward targets и ответов в формате провайдера.

Дальше: [Быстрый старт](./getting-started.md) · архитектура: [DESIGN.md](../../DESIGN.md).
