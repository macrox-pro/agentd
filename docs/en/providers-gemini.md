# Gemini CLI

> **Language:** [English](./providers-gemini.md) · [Русский](../ru/providers-gemini.md)

`--provider=gemini`. Entrypoint: `agentd hook run` (stdin).

## Install

```bash
# Project → .gemini/settings.json
agentd install --provider=gemini --scope=project

# User → settings.json under ~/.gemini
agentd install --provider=gemini --scope=user --dir ~/.gemini
```

## Runtime

1. `agentd daemon start`
2. Run Gemini CLI; hooks call `agenthooks run --provider=gemini`.

## Provider quirks

| Topic | Behavior |
|-------|----------|
| **Timeout unit** | Gemini settings use **milliseconds**; agenthooks converts Go durations on install |
| **Hook names** | Installer sets display names required by Gemini `/hooks` UX |
| **stderr** | Hookedge must **never** write debug to stderr on this path (Gemini may treat stderr specially). Put audit on async `file` / `log` ([DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine)) |
| **Ask** | ToolPre supports Ask/Deny/Allow (capability matrix) |
| **Exit codes** | Gemini blocking/prompt semantics differ from Claude; agenthooks owns encoding — do not invent exit codes in agentd |

See also: [Providers index](./providers.md), [Operations](./operations.md).
