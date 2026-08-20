# Актуализация документации

> **Language:** [English](../en/maintaining.md) · [Русский](./maintaining.md)

Как держать `docs/en/` и `docs/ru/` в соответствии с кодом. Английская версия — канон; русская зеркалит те же страницы.

## Когда обновлять

| Изменение | Что править |
|-----------|-------------|
| Новая или изменённая команда/флаг CLI (`cmd/`) | [cli.md](./cli.md), [DESIGN.md §6](../../DESIGN.md#6-cli-reference), связанные инструкции при смене сценария |
| Ключ или перечисление YAML (`internal/config/file.go`, compile) | [configuration.md](./configuration.md); при необходимости [guards.md](./guards.md) / [dispatch.md](./dispatch.md) / [approvals.md](./approvals.md); DESIGN §7 при дрейфе примеров |
| Поведение охранника / Ask / Deny | [guards.md](./guards.md), [approvals.md](./approvals.md), [troubleshooting.md](./troubleshooting.md) |
| Режим маршрута, цель, таймаут, переполнение async | [dispatch.md](./dispatch.md), [operations.md](./operations.md) |
| Поля статуса / операции демона | [operations.md](./operations.md), [cli.md](./cli.md) |
| Провайдеры install / области / точки входа / особенности | [providers.md](./providers.md) + `providers-*.md`, [getting-started.md](./getting-started.md) |
| Установка / релизы / версия | [installation.md](./installation.md) |
| Режимы ошибок / offline / таймауты | [troubleshooting.md](./troubleshooting.md) |
| Формулировки позиционирования в README | [README.md](../../README.md) и [why.md](./why.md) |

Если код и DESIGN расходятся — **описывайте поведение кода** и по возможности поправьте DESIGN в том же изменении.

## Правила

- Сначала **английская** страница, затем **русская** (то же имя файла, тот же порядок разделов). Имена команд, ключей YAML и полей JSON в русской версии оставляйте на английском в `` `коде` ``; в прозе — русские пояснения.
- Без воды: команды, ключи, поля статуса.
- Не выдумывайте флаги и ключи вне `cmd/` / `file.go`.
- Ограничения v1 (DESIGN §11) не выдавайте за возможности продукта.
- После правок документации: `make docs-check` (одинаковый набор имён файлов EN/RU).

## Карта источников (для сверки)

| Тема | Источник в коде |
|------|-----------------|
| CLI | `cmd/*.go` |
| YAML | `internal/config/file.go` |
| Пути / сохранение runtime | `internal/config/store.go`, `persist.go` |
| Охранники | `internal/guard/` |
| Маршрутизация / таймауты | `internal/dispatch/`, `timeout.go` |
| Одобрения / блокировки | `internal/config/approvals.go`, `blocks.go` |
| JSON статуса | `internal/daemon/status_write.go`, `api/agentd/v1/daemon.proto` |
| Установка в агент | `internal/install/run.go` |
| Поведение при офлайне демона | `internal/hookedge/` |

## Чеклист для PR

- [ ] Изменение, видимое пользователю → обновлён `docs/en/`
- [ ] Зеркало в `docs/ru/` (или новая страница в обоих языках)
- [ ] DESIGN §6 / §7 при смене CLI или примеров схемы
- [ ] `make docs-check` проходит
