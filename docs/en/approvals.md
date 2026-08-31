# Approvals and temporary blocks

> **Language:** [English](./approvals.md) · [Русский](../ru/approvals.md)

Record “allow” after an ask so the same action is not asked again, and set temporary tool blocks in the runtime layer.

## Ask → fingerprint → record

1. Guard returns Ask; encoded `system_message` includes `approval_fingerprint=sha256:<kind>/<hex>` where `kind` is `secrets` or `shell`.
2. Operator records Allow:

```bash
agentd config record-decision \
  --fingerprint 'sha256:shell/…' \
  --scope session \
  --session-id s1
```

Project scope with wall-clock expiry:

```bash
agentd config record-decision \
  --fingerprint 'sha256:secrets/…' \
  --scope project \
  --expires-at 2026-12-31T23:59:59Z
```

| `--scope` | Expiry |
|-----------|--------|
| `project` (default) | **24h** unless `--expires-at` (RFC3339, for example `2026-12-31T23:59:59Z`) |
| `session` | until cleared; matches `--session-id` (no wall clock by default) |

`--project-root` scopes project approvals. Default `granted_by` is `ask_user`.

Matching later `tool.pre` skips Ask within TTL.

## Temporary blocks

Runtime YAML (via `config patch`):

```yaml
version: 1
blocks:
  temporary:
    - tool: Bash
      pattern: "blocked-cmd"
      reason: operator
      until: "2026-12-31T23:59:59Z"
```

Evaluated before guards; match → Deny.

## Persist

Runtime file flush is debounced **500ms**, atomic write. Restart reloads approvals/blocks from `runtime.yaml` in the [state directory](./configuration.md#state-directory). Expired entries are dropped on compile.

See also: [Configuration](./configuration.md), [Guards](./guards.md), [CLI](./cli.md).
