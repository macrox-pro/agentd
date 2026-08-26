# Одобрения и временные блокировки

> **Language:** [English](../en/approvals.md) · [Русский](./approvals.md)

Зафиксировать разрешение после режима «спросить» (Ask) и временно запретить инструмент по шаблону. Данные лежат в runtime-слое.

## Спросить → отпечаток → запись

1. Охранник возвращает Ask; в закодированном `system_message` есть `approval_fingerprint=sha256:<kind>/<hex>`.
2. `kind` — `secrets` или `shell`.
3. Оператор записывает разрешение (Allow):

```bash
agentd config record-decision \
  --fingerprint 'sha256:shell/…' \
  --scope session \
  --session-id s1
```

| `--scope` | Срок действия |
|-----------|----------------|
| `project` (по умолчанию) | **24 часа**, если нет `--expires-at` (RFC3339) |
| `session` | до очистки; сопоставление по `--session-id` (без календарного срока по умолчанию) |

`--project-root` ограничивает одобрения уровнем проекта. Поле `granted_by` по умолчанию — `ask_user`.

Повторный подходящий `tool.pre` в пределах срока жизни **не** спрашивает снова (сценарий e2e-m7).

## Временные блокировки

YAML runtime через `config patch`:

```yaml
version: 1
blocks:
  temporary:
    - tool: Bash
      pattern: "blocked-cmd"
      reason: operator
      until: "2026-12-31T23:59:59Z"
```

Проверка идёт **до** охранников; совпадение → Deny (запрет).

## Сохранение на диск

Сброс runtime откладывается на **500 ms**, запись атомарная. После перезапуска демона одобрения и блокировки читаются из `runtime.yaml` в [state directory](./configuration.md#state-directory). Просроченные записи отбрасываются при компиляции снимка.

См. также: [Конфигурация](./configuration.md), [Охранники](./guards.md), [Справочник CLI](./cli.md).
