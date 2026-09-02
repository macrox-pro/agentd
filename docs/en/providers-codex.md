# OpenAI Codex

> **Language:** [English](./providers-codex.md) · [Русский](../ru/providers-codex.md)

Install and run agentd with OpenAI Codex (`--provider=codex`).

`--provider=codex`. Entrypoints: `hook run` (blocking/stdin path) and `hook notify` (argv JSON, async-only).

## Install

```bash
# User → $CODEX_HOME or ~/.codex/hooks.json
agentd install --provider=codex --scope=user

# Project → .codex/hooks.json under CWD
agentd install --provider=codex --scope=project
```

Writes `hooks.json` plus a managed region in `config.toml` (trust keys so Codex does not loop on interactive trust).

Non-blocking hooks may be rendered with `--async` in the command line (Codex backgrounder).

## Runtime

```bash
agentd daemon start

agentd hook notify --provider=codex '{"type":"agent-turn-complete"}'
```

Smoke (`tool.pre`):

```bash
echo '{"session_id":"s","cwd":"/tmp","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}' \
  | agentd hook run --provider=codex
```

## Provider quirks

| Topic | Behavior |
|-------|----------|
| **Cannot ask** | `tool.pre` has Deny/Allow/… but **not Ask**. Guards with `action: ask` degrade via `policy.ask_fallback` (default deny) |
| **Neutral wire** | Explicit no-op = **empty stdout**, exit 0 — **not** `{}` |
| **Notify** | Always async semantics; never use notify for blocking gates |
| **Payload cwd** | Optional top-level `cwd` in run stdin or notify argv JSON selects project `.agentd.yaml` in the daemon ([Configuration → Hook cwd](./configuration.md#hook-cwd-and-project-layer)) |
| **Trust path** | Trust state keys embed the **absolute** path of `hooks.json`; moving CODEX_HOME requires reinstall |
| **Empty stdout meaning** | Unlike Claude, empty stdout means allow/no-op in Codex dialect (handled by agenthooks encode) |

See also: [Providers index](./providers.md), [CLI](./cli.md), [Troubleshooting](./troubleshooting.md).
