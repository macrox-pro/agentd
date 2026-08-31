# Gemini CLI

> **Language:** [English](./providers-gemini.md) · [Русский](../ru/providers-gemini.md)

Install and run agentd with Gemini CLI (`--provider=gemini`).

`--provider=gemini`. Entrypoint: `agentd hook run` (stdin).

## Install

```bash
# Project → .gemini/settings.json
agentd install --provider=gemini --scope=project

# User → settings.json under ~/.gemini
agentd install --provider=gemini --scope=user
```

## Runtime

1. `agentd daemon start`
2. Run Gemini CLI; hooks call `agenthooks run --provider=gemini`.
3. Smoke:

```bash
echo '{"session_id":"s","cwd":"/tmp","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}' \
  | agentd hook run --provider=gemini
```

## Provider quirks

| Topic | Behavior |
|-------|----------|
| **Timeout unit** | Gemini settings use **milliseconds**; agenthooks converts Go durations on install |
| **Hook names** | Installer sets display names required by Gemini `/hooks` UX |
| **stderr** | The hook edge must **never** write debug to stderr on this path (Gemini may treat stderr specially). Put audit on async `file` / `log` ([DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine)) |
| **Ask** | `tool.pre` supports Ask/Deny/Allow |
| **Exit codes** | Gemini blocking/prompt semantics differ from Claude; agenthooks owns encoding — do not invent exit codes in agentd |

See also: [Providers index](./providers.md), [Operations](./operations.md).
