# Approvals и temporary blocks

> **Language:** [English](../en/approvals.md) · [Русский](./approvals.md)

Зафиксировать Allow после Ask и временно запретить tool по pattern. Хранится в runtime-слое.

## Ask → fingerprint → record

1. Guard возвращает Ask; в encoded `system_message` есть `approval_fingerprint=sha256:<kind>/<hex>`.
2. `kind` — `secrets` или `shell`.
3. Оператор записывает Allow:

```bash
agentd config record-decision \
  --fingerprint 'sha256:shell/…' \
  --scope session \
  --session-id s1
```

| `--scope` | Срок |
|-----------|------|
| `project` (default) | **24h**, если нет `--expires-at` (RFC3339) |
| `session` | до очистки; матч по `--session-id` (без wall-clock по умолчанию) |

`--project-root` ограничивает project approvals. Default `granted_by` — `ask_user`.

Повторный matching `tool.pre` в пределах TTL пропускает Ask (e2e-m7).

## Temporary blocks

Runtime YAML через `config patch`:

```yaml
version: 1
blocks:
  temporary:
    - tool: Bash
      pattern: "blocked-cmd"
      reason: operator
      until: "2026-12-31T23:59:59Z"
```

Проверка до guards; match → Deny.

## Persist

Flush runtime с debounce **500ms**, атомарная запись. После restart approvals/blocks читаются из runtime.yaml. Истёкшие записи отбрасываются при compile.

См. также: [Конфигурация](./configuration.md), [Guards](./guards.md), [CLI](./cli.md).
