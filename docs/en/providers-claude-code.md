# Claude Code

> **Language:** [English](./providers-claude-code.md) · [Русский](../ru/providers-claude-code.md)

`--provider=claude-code`. Entrypoint: `agentd hook run` (stdin JSON).

## Install

```bash
# Project → .claude/settings.json under --dir
agentd install --provider=claude-code --scope=project

# User (Dir = ~/.claude → settings.json)
agentd install --provider=claude-code --scope=user --dir ~/.claude

# Plugin → .claude-plugin/plugin.json + hooks/hooks.json
agentd install --provider=claude-code --scope=plugin --dir /path/to/plugin
```

Generated command shape: `…/agentd agenthooks run --provider=claude-code`.

## Runtime

1. `agentd daemon start`
2. Use Claude Code; installed hooks (PreToolUse, …) spawn the CLI.
3. Smoke:

```bash
echo '{"session_id":"s","cwd":"/tmp","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}' \
  | agentd hook run --provider=claude-code
```

## Provider quirks

| Topic | Behavior |
|-------|----------|
| **Wire no-op** | Neutral decision → stdout `{}`, exit 0 |
| **Ask / Deny** | ToolPre supports Deny, Ask, Allow, update-input, system message, stop-agent (agenthooks capability matrix) |
| **PromptSubmitted** | Can Deny / add context / system message — **no Ask** |
| **Timeouts** | HookSpec seconds; sync budget uses install defaults when Invoke has no deadline |
| **Blocking** | ToolPre / PromptSubmitted / Stop are blocking in the default install set |

See also: [Providers index](./providers.md), [Approvals](./approvals.md), [Guards](./guards.md).
