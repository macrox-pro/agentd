# OpenAI Codex

> **Language:** [English](./providers-codex.md) · [Русский](../ru/providers-codex.md)

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

## Provider quirks

| Topic | Behavior |
|-------|----------|
| **No CapAsk** | ToolPre has Deny/Allow/… but **not Ask**. Guards with `action: ask` degrade via `policy.ask_fallback` (default deny) |
| **Neutral wire** | Explicit no-op = **empty stdout**, exit 0 — **not** `{}` |
| **Notify** | Always async semantics; never use notify for blocking gates |
| **Trust path** | Trust state keys embed the **absolute** path of `hooks.json`; moving CODEX_HOME requires reinstall |
| **Empty stdout meaning** | Unlike Claude, empty stdout means allow/no-op in Codex dialect (handled by agenthooks encode) |

See also: [Providers index](./providers.md), [CLI](./cli.md), [Troubleshooting](./troubleshooting.md).
