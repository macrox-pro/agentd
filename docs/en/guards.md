# Guards

> **Language:** [English](./guards.md) · [Русский](../ru/guards.md)

Declarative checks that run **before** the agent’s reply is sent (typically on `tool.pre`). Configured under `guards:`. Terms: [Glossary](./glossary.md).

## Names

| Guard | Default `enabled` | Main fields |
|-------|-------------------|-------------|
| `secrets` | `true` | `action`: `ask` \| `deny` (default `ask`); optional `rules` |
| `shell` | `false` | `deny_patterns`, `ask_on` |
| `mcp` | `false` | `deny_servers` |
| `paths` | `false` | `deny_read`, `deny_write` |

On a `builtin` sync target, `guards: [secrets, shell]` runs only those guards for that route. Omit the list to run every enabled guard.

`ask_on` (shell guard): tool names that trigger Ask instead of silent allow — for example `curl`, `wget`, `ssh`.

## Enable via CLI

Shell, MCP, and paths guards default to **project** scope — from the repo root:

```bash
cd /path/to/repo
agentd config enable guard-shell
agentd config get guard-shell          # guard-shell: on (project)
```

Features: `guard-shell`, `guard-mcp`, `guard-paths` ([CLI](./cli.md#config)). **`secrets` is not a curated toggle** — it stays on by default in bootstrap YAML; change `action` / `rules` in user or project config.

## Ask vs Deny

- **Deny** — hard stop; provider encodes deny.
- **Ask** — provider Ask / permission prompt; message may include `approval_fingerprint=sha256:…` for later [Approvals](./approvals.md).

Provider capability limits still apply (some agents cannot ask). When Ask is unsupported, `policy.ask_fallback` controls the outcome: `deny` (default) or `no_decision` ([Configuration](./configuration.md#policy)).

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

See also: [Configuration](./configuration.md), [Dispatch](./dispatch.md), [Glossary](./glossary.md).
