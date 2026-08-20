# agentd user guide

> **Language:** [English](./README.md) · [Русский](../ru/README.md)

How to install, configure, and operate agentd (v1). Facts below match the shipped CLI and config schema.

## Contents

| Page | Covers |
|------|--------|
| [Getting started](./getting-started.md) | Daemon, minimal config, install hooks, verify |
| [Installation](./installation.md) | `go install`, Releases, build from source |
| [Configuration](./configuration.md) | Layers, YAML keys, validate/show/patch, reload |
| [CLI](./cli.md) | Command and flag reference |
| [Guards](./guards.md) | secrets, shell, mcp, paths |
| [Dispatch](./dispatch.md) | Modes, targets, timeouts, async queue |
| [Approvals](./approvals.md) | Ask → record-decision, blocks, persist |
| [Providers](./providers.md) | Per-agent install and entrypoints |
| [Operations](./operations.md) | Status, stop, reload |
| [Troubleshooting](./troubleshooting.md) | Common failures |
| [Maintaining docs](./maintaining.md) | When/how to update EN+RU (contributors) |

Architecture depth: [DESIGN.md](../../DESIGN.md).
