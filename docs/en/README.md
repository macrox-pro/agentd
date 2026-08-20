# agentd user guide

> **Language:** [English](./README.md) · [Русский](../ru/README.md)

How to install, configure, and operate agentd (v1). Facts below match the shipped CLI and config schema.

Start with [Why agentd](./why.md) if you need the problem statement; otherwise [Getting started](./getting-started.md).

## Contents

| Page | Covers |
|------|--------|
| [Why agentd](./why.md) | Problems, pains, what v1 is / is not |
| [Getting started](./getting-started.md) | Daemon, minimal config, install hooks, verify |
| [Installation](./installation.md) | `go install`, Releases, build from source |
| [Configuration](./configuration.md) | Layers, YAML keys, validate/show/patch, reload |
| [CLI](./cli.md) | Command and flag reference |
| [Guards](./guards.md) | secrets, shell, mcp, paths |
| [Dispatch](./dispatch.md) | Modes, targets, timeouts, async queue |
| [Approvals](./approvals.md) | Ask → record-decision, blocks, persist |
| [Providers](./providers.md) | Per-agent install, entrypoints, **quirks** — [Claude](./providers-claude-code.md) · [Cursor](./providers-cursor.md) · [Codex](./providers-codex.md) · [Gemini](./providers-gemini.md) · [OpenCode](./providers-opencode.md) · [Kimi](./providers-kimi.md) |
| [Operations](./operations.md) | Status, stop, reload |
| [Troubleshooting](./troubleshooting.md) | Common failures |
| [Maintaining docs](./maintaining.md) | When/how to update EN+RU (contributors) |

Architecture depth: [DESIGN.md](../../DESIGN.md).
