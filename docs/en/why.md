# Why agentd

> **Language:** [English](./why.md) · [Русский](../ru/why.md)

agentd is a **user-level daemon** that sits between coding agents and your hook logic: one policy surface, provider-correct wire I/O, sync decisions and async side effects without re-implementing each agent’s hook dialect.

Built on [agenthooks](https://github.com/speakeasy-api/agenthooks) for codecs and install targets.

## Problems it addresses

| Pain | Without agentd | With agentd |
|------|----------------|-------------|
| **N providers, N dialects** | Separate scripts/timeouts/exit codes for Claude, Cursor, Codex, Gemini, OpenCode, Kimi | Thin `hook run` / `serve` / `notify`; daemon owns policy |
| **Cold start on every tool call** | Heavy logic spawned per hook process → latency and flaky timeouts | Long-lived daemon; hot path uses in-memory config snapshot (no disk I/O per Invoke) |
| **Guards mixed with audit** | One process must both block the agent and fire webhooks/metrics | Sync pipeline (Ask/Deny) vs async queue (log/http/exec/…) — async never blocks the wire response |
| **Policy drift across repos** | Copy-paste hook configs; hard to approve once / block temporarily | Layered YAML (user ⊕ project ⊕ runtime); approvals + temporary blocks with persist |
| **Ops blind spot** | No single place for queue pressure or “is the gate up?” | `daemon status --json`: generation, fingerprint, `async_queue_depth`, `async_dropped_count` |

## What it is not (v1)

Not an agent auth product, transcript pipeline, plugin runtime, or general hooks DSL. Targets are declarative YAML; exec stays **async-only**. See [DESIGN.md §11](../../DESIGN.md#11-non-goals-v1).

## Who it is for

Engineers who run coding-agent hooks in real workflows and need **one** place for secrets/shell/MCP/path guards, optional forward targets, and provider-faithful responses.

Next: [Getting started](./getting-started.md) · architecture: [DESIGN.md](../../DESIGN.md).
