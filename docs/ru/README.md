# Руководство пользователя agentd

> **Language:** [English](../en/README.md) · [Русский](./README.md)

Установка, конфигурация и эксплуатация agentd (v1). Текст сверен с CLI и YAML-схемой.

## Содержание

| Страница | Содержание |
|----------|------------|
| [Быстрый старт](./getting-started.md) | Демон, минимальный конфиг, install hooks, проверка |
| [Установка](./installation.md) | `go install`, Releases, сборка из исходников |
| [Конфигурация](./configuration.md) | Слои, ключи YAML, validate/show/patch, reload |
| [CLI](./cli.md) | Команды и флаги |
| [Guards](./guards.md) | secrets, shell, mcp, paths |
| [Dispatch](./dispatch.md) | Режимы, targets, таймауты, async queue |
| [Approvals](./approvals.md) | Ask → record-decision, blocks, persist |
| [Providers](./providers.md) | Install и entrypoint по агентам |
| [Эксплуатация](./operations.md) | Status, stop, reload |
| [Troubleshooting](./troubleshooting.md) | Типичные сбои |
| [Актуализация docs](./maintaining.md) | Когда/как обновлять EN+RU (для контрибьюторов) |

Архитектура: [DESIGN.md](../../DESIGN.md).
