# agentd user guide

> **Language:** [English](./README.md) · [Русский](../ru/README.md)

Install, configure, and run agentd (v1). Pages match the shipped commands and YAML schema.

> **Release:** [v0.0.7](../../CHANGELOG.md#v007--2026-08-31) — setup wizard, doctor, session ledger on by default.

New here? [Why agentd](./why.md), then [Getting started](./getting-started.md).

## Contents

| Page | What you will find |
|------|--------------------|
| [Why agentd](./why.md) | Problem, limits of v1 |
| [Getting started](./getting-started.md) | Daemon, first config, connect an agent |
| [Installation](./installation.md) | `go install`, GitHub Releases, build from source |
| [Configuration](./configuration.md) | Layers, on-disk paths, YAML keys, reload |
| [CLI](./cli.md) | Commands and flags |
| [Guards](./guards.md) | Secrets, shell, MCP, paths |
| [Dispatch](./dispatch.md) | How events are routed, timeouts, async queue |
| [Approvals](./approvals.md) | Ask once, then allow; temporary blocks |
| [Providers](./providers.md) | Per-agent install and limits — [Claude](./providers-claude-code.md) · [Cursor](./providers-cursor.md) · [Codex](./providers-codex.md) · [Gemini](./providers-gemini.md) · [OpenCode](./providers-opencode.md) · [Kimi](./providers-kimi.md) |
| [Operations](./operations.md) | Status, stop, reload, session stats |
| [Troubleshooting](./troubleshooting.md) | Common failures |
| [Maintaining docs](./maintaining.md) | When to update EN+RU (contributors) |

Architecture: [DESIGN.md](../../DESIGN.md).
