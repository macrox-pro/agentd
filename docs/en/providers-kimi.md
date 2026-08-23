# Kimi Code

> **Language:** [English](./providers-kimi.md) · [Русский](../ru/providers-kimi.md)

`--provider=kimi-code` (preferred for encode) or `kimicode` (accepted by agentd parse). Entrypoint: `agentd hook run`.

## Install

Hooks live only in **user-level** `config.toml` (`$KIMI_CODE_HOME`, default `~/.kimi-code`). **`--scope=project` fails**.

```bash
agentd install --provider=kimi-code --scope=user
```

Updates a managed `[[hooks]]` region in `config.toml`.

## Runtime

1. `agentd daemon start`
2. Use Kimi; hooks invoke `agenthooks run --provider=kimi-code`.

Prefer **`kimi-code`** on the CLI so agenthooks recognizes the provider.

## Provider quirks

| Topic | Behavior |
|-------|----------|
| **No CapAsk** | ToolPre is Deny/Allow only — no Ask, no update-input, no additionalContext in JSON decisions |
| **Neutral wire** | No-op = **empty stdout**, exit 0 (same family as Codex) |
| **Prompt / stop blocking** | Block-style outcomes use exit **2** + stderr (agenthooks quirk); do not parse as Claude `{}` |
| **Observation-only events** | Many kinds (e.g. PostToolUse, PermissionRequest) are observe-only — empty capability set |
| **Install scope** | User only; project-level hooks config does not exist for Kimi |

See also: [Providers index](./providers.md), [Troubleshooting](./troubleshooting.md), [Guards](./guards.md).
