# OpenCode

> **Language:** [English](./providers-opencode.md) · [Русский](../ru/providers-opencode.md)

Install and run agentd with OpenCode (`--provider=opencode`).

`--provider=opencode`. Entrypoint: long-lived **`agentd hook serve`** (NDJSON stdio), not per-event `hook run`.

## Install

```bash
agentd install --provider=opencode --scope=project --dir /path/to/repo
```

Writes `.opencode/plugin/agenthooks.ts`, spawning:

`agentd agenthooks serve --provider=opencode`

(= `agentd hook serve --provider=opencode`).

## Runtime

1. `agentd daemon start`
2. OpenCode loads the plugin; a small TypeScript wrapper keeps a child process on `hook serve`.
3. Each NDJSON line on stdin/stdout is one event round-trip.

```bash
agentd hook serve --provider=opencode
```

`hook serve` rejects any `--provider` other than `opencode`.

When the daemon is up, each frame's resolved cwd selects the project config layer ([Configuration → Hook cwd](./configuration.md#hook-cwd-and-project-layer)). If frames omit `cwd` and `workspace_roots`, the daemon may fall back to the serve process working directory — set `cwd` in the wire JSON for per-repo policy.

## Provider quirks

| Topic | Behavior |
|-------|----------|
| **Process model** | One long-lived serve process; events multiplexed over stdio |
| **Project cwd per frame** | Each NDJSON line is one Invoke. Project `.agentd.yaml` follows cwd from **that frame's JSON**: top-level `cwd`, else `workspace_roots[0]`. agentd does **not** read OpenCode's `input.directory` field — the install shim should set top-level `cwd` on each frame (or ensure tool frames carry `cwd`) |
| **Session ordering** | The daemon serializes sync work per session (and coordinates async) so replies stay in order ([DESIGN.md §1](../../DESIGN.md)) |
| **Cannot ask on tool.pre** | `tool.pre` capabilities are Deny + update-input only — no Ask. Use Deny or `ask_fallback` policy |
| **Stop / session.idle** | Stop-like events may be non-blocking / empty capability set — do not expect Continue on every stop shape |
| **Permission** | Permission events go through OpenCode’s permission channel (allow/deny), not Claude-style Ask JSON |

See also: [Providers index](./providers.md), [CLI](./cli.md), [Glossary](./glossary.md).
