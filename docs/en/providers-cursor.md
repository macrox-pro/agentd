# Cursor

> **Language:** [English](./providers-cursor.md) · [Русский](../ru/providers-cursor.md)

Install and run agentd with Cursor (`--provider=cursor`).

`--provider=cursor`. Entrypoint: `agentd hook run` with **`--argv-payload`** (payload on argv, not only stdin).

## Install

```bash
# Project → .cursor/hooks.json
agentd install --provider=cursor --scope=project

# User → hooks.json under ~/.cursor
agentd install --provider=cursor --global
# same as: --scope=user

# Plugin → .cursor-plugin/plugin.json + hooks/hooks.json
agentd install --provider=cursor --scope=plugin --dir /path/to/plugin
```

## Runtime

1. `agentd daemon start`
2. Cursor runs the installed command (`agenthooks run --provider=cursor` + argv-payload as rendered).
3. Smoke:

```bash
agentd hook run --provider=cursor --argv-payload '{"session_id":"s","cwd":"/tmp","hook_event_name":"preToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}'
```

## Provider quirks

| Topic | Behavior |
|-------|----------|
| **Ask** | Enforced only on native `beforeShellExecution` / `beforeMCPExecution`. On generic `preToolUse`, Ask is ignored. `beforeReadFile` is allow/deny only |
| **Install + URLs** | If the hook **command string** contains a URL, install refuses to render the entry — Cursor would drop the entire `hooks.json`. Keep endpoints in config or environment variables, not in the command line |
| **Async vs sync** | Async/telemetry failure must **not** change the sync decision ([DESIGN.md §2](../../DESIGN.md#2-hook-dispatch-engine)) |
| **PromptSubmitted** | Capability surface is Deny-only at the kind level |
| **failClosed** | Install/runtime follow Cursor’s stricter fail-closed expectations via agenthooks |
| **Default install** | Includes `subagentStart` / `preCompact` and related observe hooks. Does **not** write `afterAgentThought` or `afterFileEdit`. Re-run `agentd install` after upgrading agentd |

See also: [Providers index](./providers.md), [Dispatch](./dispatch.md), [Troubleshooting](./troubleshooting.md).
