# Guards

> **Language:** [English](./guards.md) · [Русский](../ru/guards.md)

Declarative sync checks on `tool.pre` (and related kinds via builtin). Config under `guards:`.

## Names

| Guard | Default `enabled` | Main fields |
|-------|-------------------|-------------|
| `secrets` | `true` | `action`: `ask` \| `deny` (default `ask`); optional `rules` |
| `shell` | `false` | `deny_patterns`, `ask_on` |
| `mcp` | `false` | `deny_servers` |
| `paths` | `false` | `deny_read`, `deny_write` |

Route sync builtin can subset: `guards: [secrets, shell]`. Omitting the list uses the compiled default set for that target.

## Ask vs Deny

- **Deny** — hard stop; provider encodes deny.
- **Ask** — provider Ask / permission prompt; message may include `approval_fingerprint=sha256:…` for later [Approvals](./approvals.md).

Provider capability limits still apply (some agents cannot Ask).

## Example

```yaml
guards:
  secrets:
    enabled: true
    action: ask
  shell:
    enabled: true
    deny_patterns: ["rm -rf /"]
    ask_on: [curl, wget, ssh]
  mcp:
    enabled: true
    deny_servers: ["untrusted-*"]
  paths:
    enabled: true
    deny_read: ["/etc/shadow"]
    deny_write: ["**/.env"]
```

Temporary tool blocks (runtime) run before guards — see [Approvals](./approvals.md).

Architecture: [DESIGN.md](../../DESIGN.md). Wire behavior: [Dispatch](./dispatch.md).
