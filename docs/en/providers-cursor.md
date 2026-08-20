# Cursor

> **Language:** [English](./providers-cursor.md) · [Русский](../ru/providers-cursor.md)

`--provider=cursor`. Entrypoint: `agentd hook run` with **`--argv-payload`** (payload on argv, not only stdin).

## Install

```bash
# Project → .cursor/hooks.json
agentd install --provider=cursor --scope=project

# User → hooks.json under ~/.cursor
agentd install --provider=cursor --scope=user --dir ~/.cursor

# Plugin → .cursor-plugin/plugin.json + hooks/hooks.json
agentd install --provider=cursor --scope=plugin --dir /path/to/plugin
```

## Runtime

1. `agentd daemon start`
2. Cursor runs the installed command (`agenthooks run --provider=cursor` + argv-payload as rendered).

```bash
agentd hook run --provider=cursor --argv-payload '<json>'
```

## Provider quirks

| Topic | Behavior |
|-------|----------|
| **Ask** | Enforced only on native `beforeShellExecution` / `beforeMCPExecution`. On generic `preToolUse`, Ask is ignored. `beforeReadFile` is allow/deny only |
| **Install + URLs** | If a hook **command string contains a URL**, agenthooks **refuses to render** — Cursor would drop the entire `hooks.json`. Keep endpoints in config/env, not in argv |
| **Async vs sync** | Async/telemetry failure must **not** change the sync decision ([DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine)) |
| **PromptSubmitted** | Capability surface is Deny-only at the kind level |
| **failClosed** | Install/runtime follow Cursor’s stricter fail-closed expectations via agenthooks |

See also: [Providers index](./providers.md), [Dispatch](./dispatch.md), [Troubleshooting](./troubleshooting.md).
