# Актуализация user docs

> **Language:** [English](../en/maintaining.md) · [Русский](./maintaining.md)

Как держать `docs/en/` и `docs/ru/` в синхроне с кодом. EN — канон; RU зеркалит те же страницы.

## Когда обновлять

| Изменение | Обновить |
|-----------|----------|
| Новая/изменённая CLI-команда или флаг (`cmd/`) | [cli.md](./cli.md), [DESIGN.md §6](../../DESIGN.md#6-cli-reference), связанные how-to при смене UX |
| YAML-ключ / enum (`internal/config/file.go`, compile) | [configuration.md](./configuration.md), при необходимости [guards.md](./guards.md) / [dispatch.md](./dispatch.md) / [approvals.md](./approvals.md); DESIGN §7 при дрейфе примеров |
| Поведение guard / Ask / Deny | [guards.md](./guards.md), [approvals.md](./approvals.md), [troubleshooting.md](./troubleshooting.md) |
| Dispatch mode, target, timeout, async overflow | [dispatch.md](./dispatch.md), [operations.md](./operations.md) |
| Поля Status / ops демона | [operations.md](./operations.md), [cli.md](./cli.md) |
| Install providers / scopes / entrypoints | [providers.md](./providers.md), [getting-started.md](./getting-started.md) |
| Install / Releases / version | [installation.md](./installation.md) |
| Failure modes / offline / timeouts | [troubleshooting.md](./troubleshooting.md) |
| User-visible заявления в README | [README.md](../../README.md) + соответствующая страница docs |

Если код и DESIGN расходятся — **документируйте код** и по возможности поправьте DESIGN в том же изменении.

## Правила

- Сначала **EN**, затем **RU** (то же имя файла, тот же порядок секций). Идентификаторы в RU остаются на английском.
- Без воды; указывайте команды, ключи и поля Status.
- Не выдумывайте флаги и YAML-ключи вне `cmd/` / `file.go`.
- Non-goals (DESIGN §11) не рекламировать как фичи.
- После правок docs: `make docs-check` (parity имён EN/RU).

## Карта источников (сверка)

| Тема | Источник |
|------|----------|
| CLI | `cmd/*.go` |
| YAML | `internal/config/file.go` |
| Paths / persist | `internal/config/store.go`, `persist.go` |
| Guards | `internal/guard/` |
| Dispatch / timeout | `internal/dispatch/`, `timeout.go` |
| Approvals / blocks | `internal/config/approvals.go`, `blocks.go` |
| Status JSON | `internal/daemon/status_write.go`, `api/agentd/v1/daemon.proto` |
| Install | `internal/install/run.go` |
| Hook offline | `internal/hookedge/` |

## PR checklist (docs)

- [ ] User-visible изменение → обновлён `docs/en/`
- [ ] Зеркало в `docs/ru/` (или новая страница в обоих locales)
- [ ] DESIGN §6 / §7 при смене CLI или примеров schema
- [ ] `make docs-check` проходит
