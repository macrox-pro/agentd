# Руководство пользователя agentd

> **Language:** [English](../en/README.md) · [Русский](./README.md)

Установка, настройка и эксплуатация agentd (v1). Содержание сверено с CLI и схемой YAML.

> **Релиз:** [v0.0.3](../../CHANGELOG.md#v003--2026-08-26) — `policy.offline` на краю хука при недоступном демоне.

Начните с [Зачем нужен agentd](./why.md), если нужна постановка задачи; иначе — с [Быстрого старта](./getting-started.md).

## Содержание

| Страница | О чём |
|----------|--------|
| [Зачем нужен agentd](./why.md) | Задачи продукта, боли, границы v1 |
| [Быстрый старт](./getting-started.md) | Демон, минимальный конфиг, установка хуков, проверка |
| [Установка](./installation.md) | `go install`, релизы, сборка из исходников |
| [Конфигурация](./configuration.md) | Слои, пути на диске (state directory), ключи YAML, validate/show/patch, перезагрузка |
| [Справочник CLI](./cli.md) | Команды и флаги |
| [Охранники (guards)](./guards.md) | `secrets`, `shell`, `mcp`, `paths` |
| [Маршрутизация (dispatch)](./dispatch.md) | Режимы, цели, таймауты, асинхронная очередь |
| [Одобрения (approvals)](./approvals.md) | Ask → `record-decision`, временные блокировки, сохранение |
| [Провайдеры](./providers.md) | Install, точки входа и **особенности** — [Claude](./providers-claude-code.md) · [Cursor](./providers-cursor.md) · [Codex](./providers-codex.md) · [Gemini](./providers-gemini.md) · [OpenCode](./providers-opencode.md) · [Kimi](./providers-kimi.md) |
| [Эксплуатация](./operations.md) | Статус, остановка, перезагрузка конфига |
| [Диагностика](./troubleshooting.md) | Типичные сбои |
| [Актуализация документации](./maintaining.md) | Когда и как обновлять EN+RU (для разработчиков) |

Архитектура: [DESIGN.md](../../DESIGN.md).
