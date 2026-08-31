# Why agentd

> **Language:** [English](./why.md) · [Русский](../ru/why.md)

What agentd is for and what it is not in the current release.

Terms: [Glossary](./glossary.md).

**agentd** is a background service on your machine. It sits between coding agents and your rules: one policy for every agent, a reply in the format that agent expects, and side effects (logs, HTTP) that never delay that reply.

Hook formats and writing agent settings use [agenthooks](https://github.com/speakeasy-api/agenthooks).

| Term | Meaning |
|------|---------|
| **Hook** | Callback the agent runs at an event (for example: before a tool). |
| **Coding agent** | Product such as Claude Code or Cursor (`--provider` in commands). |
| **Daemon** | One long-lived `agentd` process per user. Hooks send events here; policy is decided here. |
| **Sync path** | Decides allow / ask / deny and shapes the reply to the agent. |
| **Async path** | Audit and notifications. Must not block the reply. |

## Problems it addresses

| Pain | Without agentd | With agentd |
|------|----------------|-------------|
| **Many agents, many formats** | A separate script, timeout, and exit code per product | Thin `hook run` / `serve` / `notify`; policy lives in the daemon |
| **Cold start on every tool call** | Heavy logic in a new process → latency and flaky timeouts | One long-lived process; each call uses an in-memory config snapshot |
| **Guards mixed with audit** | The same script must both block the agent and fire webhooks | Sync pipeline (ask / deny) separate from an async queue (`log` / `http` / `exec`) |
| **Policy drift across repos** | Copied hook files; hard to approve once or block for a while | Layered YAML (user + project + runtime); approvals and temporary blocks |
| **No operations picture** | Unclear whether the gate is up or the queue is full | `daemon status --json`: config generation, fingerprint, queue depth, drops |

## What it is not (current release)

Not a login product for agent accounts, not a full transcript pipeline, not a plugin loader, and not a separate rules language. Routes are YAML. Running an external program (`exec`) is **async only**. See [DESIGN.md §11](../../DESIGN.md#11-non-goals-v1).

## Who it is for

People who already use coding agents at work and want **one** place to:

- **Guard** — secrets, shell, MCP, paths; ask or deny; temporary blocks
- **Observe** — write events to `log` / `http` / `exec` without delaying the agent
- **See what happened** — one hook stream across agents, plus `daemon status`
- **Reply in the agent’s native format** — without re-implementing each dialect in your scripts

Typical uses: audit of “what the agent tried,” alerts on hooks, one policy for Claude + Cursor + others.

Next: [Getting started](./getting-started.md) · architecture: [DESIGN.md](../../DESIGN.md).
