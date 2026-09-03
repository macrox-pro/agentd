# Glossary

> **Language:** [English](./glossary.md) · [Русский](../ru/glossary.md)

One-sentence definitions for terms used across the user guide. Each links to the page with full detail.

## Core concepts

| Term | Meaning | Detail |
|------|---------|--------|
| **Coding agent** | A tool that writes and edits code for you (Claude Code, Cursor, Codex, Gemini CLI, OpenCode, Kimi Code). | [Providers](./providers.md) |
| **Provider** | The agent id you pass as `--provider` on hook and install commands (for example `claude-code`, `cursor`). | [Providers](./providers.md) |
| **Hook** | A short program the agent runs at a lifecycle moment — for example before it executes a shell command. | [Why agentd](./why.md) |
| **Daemon** | One long-lived `agentd` process per user. Hooks send events here; policy is decided here. | [Getting started](./getting-started.md) |
| **Hook edge** | The thin `agentd hook run` / `notify` / `serve` process the agent spawns per event. It decodes the wire format, calls the daemon once, and encodes the reply — no policy logic. | [CLI → hook](./cli.md#hook) |

## Paths and policy

| Term | Meaning | Detail |
|------|---------|--------|
| **Sync path** | The blocking pipeline that decides allow, ask, or deny and shapes the reply the agent waits for. | [Why agentd](./why.md) · [Dispatch](./dispatch.md) |
| **Async path** | Side effects (logs, HTTP, exec) that run without delaying the sync reply. | [Why agentd](./why.md) · [Dispatch](./dispatch.md) |
| **Decision** | What the sync path returns: **allow** (continue), **ask** (prompt the user), or **deny** (block). | [Guards](./guards.md) |
| **Guard** | A declarative check (secrets, shell, MCP, paths) that can ask or deny before a tool runs. | [Guards](./guards.md) |
| **Dispatch** | Routing: which guards and targets run for each event, and in what order. | [Dispatch](./dispatch.md) |
| **Route** | One named rule in `dispatch:` — a `match` filter plus sync/async target lists. | [Dispatch](./dispatch.md) |
| **Target** | A destination for an event: `builtin` (guards), `log`, `http`, `exec`, `grpc`, or `file`. | [Dispatch](./dispatch.md) |
| **Mode** | How sync and async combine: `sync_only`, `async_only`, `parallel`, `after_sync` / `sync_then_async`. | [Dispatch](./dispatch.md) |
| **Approval** | A recorded "allow" after the user answered ask — the same action is not asked again within TTL. | [Approvals](./approvals.md) |
| **Temporary block** | A runtime rule that denies matching tools until an expiry time. | [Approvals](./approvals.md) |

## Config and storage

| Term | Meaning | Detail |
|------|---------|--------|
| **Config layers** | Four merged sources (low → high): defaults, user (`~/.agentd.yaml`), project (`.agentd.yaml`), runtime (`runtime.yaml`). | [Configuration](./configuration.md) |
| **Runtime layer** | Daemon-written overlay for approvals and temporary blocks (`runtime.yaml` in the state directory). | [Configuration](./configuration.md) |
| **Snapshot** | The in-memory merged config the daemon reads on each hook call — no disk I/O on the hot path. | [Configuration](./configuration.md) |
| **State directory** | Where the daemon writes runtime overlay, operational log, and session ledger files (not user config). | [Configuration → State directory](./configuration.md#state-directory) |
| **`policy.fail`** | When the **daemon** sync pipeline errors (timeout/cancel, sync target error surfaced to the engine): `fail_open` → neutral allow; `fail_closed` → deny/block when the event supports it. Not used for guard deny or for `policy.offline`. | [Configuration → policy](./configuration.md#policy) · [Dispatch](./dispatch.md) |
| **`policy.offline`** | What the hook edge does when the daemon is unreachable (`fail_open` → exit 0 + neutral wire; `fail_closed` → exit **1**). | [Troubleshooting → Daemon not running](./troubleshooting.md#daemon-not-running) |
| **`fail_open`** | Policy **mode**: on errors, allow the agent to continue with a neutral reply (used by `policy.fail` or `policy.offline` depending on context). | [Configuration → policy](./configuration.md#policy) |
| **`fail_closed`** | Policy **mode**: on errors, block or exit **1** (used by `policy.fail` on the daemon or `policy.offline` on the hook edge). | [Configuration → policy](./configuration.md#policy) |

## Session ledger (trajectory)

| Term | Meaning | Detail |
|------|---------|--------|
| **Session ledger** | A chronological JSONL log of hook calls and related events (`trajectory`, on by default). | [Trajectory](./trajectory.md) |
| **Transcript** | On-disk chat history from an agent; optional import adds `transcript/*` events to the ledger. | [Trajectory → Import](./trajectory.md#import) |

## Event kinds (`match: kind:`)

Wire names used in YAML `match:` and `dispatch_defaults`. Provider JSON may use different spellings; agentd normalizes them.

| Wire kind | Typical provider event | When it fires |
|-----------|------------------------|---------------|
| `tool.pre` | `PreToolUse` (Claude), `preToolUse` (Cursor), etc. | Before a tool runs |
| `prompt.submitted` | `PromptSubmitted`, `UserPromptSubmit`, etc. | User sent a prompt |
| `agent.stop` | `Stop` | Agent session ending (not `subagent.stop`) |
| `tool.post` | `PostToolUse`, etc. | After a tool finished |
| `tool.error` | `PostToolUseFailure`, etc. | Tool failed |
| `permission.request` | Permission / approval prompts | Agent asked to run a tool |
| `session.start` | `SessionStart` | New agent session |
| `session.end` | `SessionEnd` | Session closed |
| `subagent.start` | `SubagentStart` | Nested agent started |
| `subagent.stop` | `SubagentStop` | Nested agent finished |
| `compact.pre` | `PreCompact` | Context compaction about to run |
| `compact.post` | `PostCompact` | Compaction finished |
| `file.edited` | `afterFileEdit` (Cursor) | File changed (not in the default install set) |
| `model.response` | After-agent-thought style frames | Model output (not in the default install set) |
| `notification` | Codex `notify`, observe-only frames | Fire-and-forget observation |
| `other` | Unmapped native names | Catch-all when no exact default matches |

Full routing defaults: [Dispatch → Kind defaults](./dispatch.md#kind-defaults).

## Coverage tiers (trajectory)

| Tier | Meaning |
|------|---------|
| **L0 Live** | Every hook call → `hook/invoked` + `hook/decided` (required for all six providers). |
| **L1 Correlate** | Stable session and tool ids in events (quality varies; not a separate doc tier). |
| **L2 Import** | On-disk transcript → `transcript/*` events via `session import`. |
| **L3 Thinking** | Reasoning/thinking lines when the vendor persists them (provider-specific). |

Per-provider matrix: [Trajectory → Coverage tiers](./trajectory.md#coverage-tiers).

## Technical terms

| Term | Meaning |
|------|---------|
| **NDJSON** | Newline-delimited JSON — one JSON object per line (OpenCode `hook serve` uses this on stdin/stdout). |
| **TUI** | Text user interface — interactive terminal wizard (`agentd setup`, bare `agentd install`). |
| **IPC** | Inter-process communication — gRPC over a Unix socket (Linux/macOS) or named pipe (Windows). |
| **TTL** | Time to live — how long an approval stays valid before expiring. |
| **gRPC** | The RPC protocol the hook edge and management CLI use to talk to the daemon. |

See also: [Why agentd](./why.md) · [Getting started](./getting-started.md).
