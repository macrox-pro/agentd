# Gemini Documentation Research

Structured verbatim excerpts from [geminicli.com/docs](https://geminicli.com/docs), [ai.google.dev/gemini-api/docs](https://ai.google.dev/gemini-api/docs) (Managed Agents), and selected migration sources.

**Snapshot date: 2026-08-29** (see `gemini_docs_snapshot` in topic front matter).

> **Deprecation note:** Gemini CLI stopped serving individual-account requests on **2026-06-18**; successor terminal product is **Antigravity CLI**. Enterprise/API-key Gemini CLI remains supported. See [09-migration-and-antigravity/](./09-migration-and-antigravity/).

## How to read

- Each topic file under `01-agent-loop/` … `09-migration-and-antigravity/` starts with YAML front matter.
- Body text is verbatim from primary sources under `### Source:` headings.
- Machine-readable schemas: [`schemas/`](./schemas/).

## Topic index

| Topic | File |
|-------|------|
| Checkpointing and rewind | [01-agent-loop/checkpointing-and-rewind.md](./01-agent-loop/checkpointing-and-rewind.md) |
| Core turn loop | [01-agent-loop/core-turn-loop.md](./01-agent-loop/core-turn-loop.md) |
| Git worktrees | [01-agent-loop/git-worktrees.md](./01-agent-loop/git-worktrees.md) |
| Overview and modes | [01-agent-loop/overview-and-modes.md](./01-agent-loop/overview-and-modes.md) |
| Policy engine | [01-agent-loop/policy-engine.md](./01-agent-loop/policy-engine.md) |
| Sandbox and trusted folders | [01-agent-loop/sandbox-and-trust.md](./01-agent-loop/sandbox-and-trust.md) |
| Sessions and history | [01-agent-loop/sessions-and-history.md](./01-agent-loop/sessions-and-history.md) |
| Subagents | [01-agent-loop/subagents.md](./01-agent-loop/subagents.md) |
| Tools reference | [01-agent-loop/tools-reference.md](./01-agent-loop/tools-reference.md) |
| Hooks Best Practices | [02-hooks/best-practices.md](./02-hooks/best-practices.md) |
| Cloud hooks matrix | [02-hooks/cloud-hooks-matrix.md](./02-hooks/cloud-hooks-matrix.md) |
| Hook events — agent | [02-hooks/events-reference-agent.md](./02-hooks/events-reference-agent.md) |
| Hook events — lifecycle | [02-hooks/events-reference-lifecycle.md](./02-hooks/events-reference-lifecycle.md) |
| Hook events — model | [02-hooks/events-reference-model.md](./02-hooks/events-reference-model.md) |
| Hook events — tools | [02-hooks/events-reference-tools.md](./02-hooks/events-reference-tools.md) |
| Hook execution and protocol | [02-hooks/execution-and-protocol.md](./02-hooks/execution-and-protocol.md) |
| Hooks overview | [02-hooks/hooks-overview.md](./02-hooks/hooks-overview.md) |
| Notifications | [02-hooks/notifications.md](./02-hooks/notifications.md) |
| Writing hooks for Gemini CLI | [02-hooks/writing-hooks.md](./02-hooks/writing-hooks.md) |
| MCP servers with Gemini CLI | [03-mcp/mcp-servers.md](./03-mcp/mcp-servers.md) |
| Set up an MCP server | [03-mcp/mcp-setup.md](./03-mcp/mcp-setup.md) |
| MCP resource tools | [03-mcp/resource-tools.md](./03-mcp/resource-tools.md) |
| Custom commands | [04-customization/custom-commands.md](./04-customization/custom-commands.md) |
| Extensions — developer guide | [04-customization/extensions-developer.md](./04-customization/extensions-developer.md) |
| Gemini CLI extensions | [04-customization/extensions-overview-install.md](./04-customization/extensions-overview-install.md) |
| Provide context with GEMINI.md files | [04-customization/gemini-md.md](./04-customization/gemini-md.md) |
| Memory and context | [04-customization/memory.md](./04-customization/memory.md) |
| Agent Skills | [04-customization/skills.md](./04-customization/skills.md) |
| System Prompt Override (GEMINI_SYSTEM_MD) | [04-customization/system-prompt-override.md](./04-customization/system-prompt-override.md) |
| Configuration reference — advanced | [05-settings-and-paths/config-reference-advanced.md](./05-settings-and-paths/config-reference-advanced.md) |
| Configuration reference — general | [05-settings-and-paths/config-reference-general.md](./05-settings-and-paths/config-reference-general.md) |
| Configuration reference — MCP and model | [05-settings-and-paths/config-reference-mcp-model.md](./05-settings-and-paths/config-reference-mcp-model.md) |
| Environment variables | [05-settings-and-paths/environment-variables.md](./05-settings-and-paths/environment-variables.md) |
| Gemini directory layout | [05-settings-and-paths/gemini-directory.md](./05-settings-and-paths/gemini-directory.md) |
| Ignore files (.geminiignore) | [05-settings-and-paths/gemini-ignore.md](./05-settings-and-paths/gemini-ignore.md) |
| Project vs global path matrix | [05-settings-and-paths/project-vs-global-matrix.md](./05-settings-and-paths/project-vs-global-matrix.md) |
| Settings and precedence | [05-settings-and-paths/settings-and-precedence.md](./05-settings-and-paths/settings-and-precedence.md) |
| CLI flags and cheatsheet | [06-cli-and-sdk/cli-flags.md](./06-cli-and-sdk/cli-flags.md) |
| Command reference | [06-cli-and-sdk/commands-reference.md](./06-cli-and-sdk/commands-reference.md) |
| Headless mode | [06-cli-and-sdk/headless-mode.md](./06-cli-and-sdk/headless-mode.md) |
| IDE integration and ACP | [06-cli-and-sdk/ide-integration-acp.md](./06-cli-and-sdk/ide-integration-acp.md) |
| Telemetry | [06-cli-and-sdk/telemetry.md](./06-cli-and-sdk/telemetry.md) |
| Agent environments | [07-cloud-and-api-agents/agent-environment.md](./07-cloud-and-api-agents/agent-environment.md) |
| Agents overview | [07-cloud-and-api-agents/agents-overview.md](./07-cloud-and-api-agents/agents-overview.md) |
| Managed agents API hooks | [07-cloud-and-api-agents/api-hooks.md](./07-cloud-and-api-agents/api-hooks.md) |
| Antigravity agent | [07-cloud-and-api-agents/building-managed-agents.md](./07-cloud-and-api-agents/building-managed-agents.md) |
| Coding agent setup | [07-cloud-and-api-agents/coding-agent-setup.md](./07-cloud-and-api-agents/coding-agent-setup.md) |
| Enterprise managed agents | [07-cloud-and-api-agents/enterprise-managed-agents.md](./07-cloud-and-api-agents/enterprise-managed-agents.md) |
| Managed agents hooks announcement | [07-cloud-and-api-agents/hooks-announcement.md](./07-cloud-and-api-agents/hooks-announcement.md) |
| Managed agents quickstart | [07-cloud-and-api-agents/managed-agents-quickstart.md](./07-cloud-and-api-agents/managed-agents-quickstart.md) |
| Enterprise configuration | [08-enterprise/enterprise-config.md](./08-enterprise/enterprise-config.md) |
| Enterprise admin controls | [08-enterprise/enterprise-controls.md](./08-enterprise/enterprise-controls.md) |
| Gemini CLI to Antigravity migration | [09-migration-and-antigravity/gcli-migration.md](./09-migration-and-antigravity/gcli-migration.md) |
| Gemini CLI transition timeline | [09-migration-and-antigravity/transition-timeline.md](./09-migration-and-antigravity/transition-timeline.md) |

## Meta

- Source registry: [SOURCES.md](./SOURCES.md)
- Topic file pattern: [TEMPLATE.md](./TEMPLATE.md)
- Schemas: [schemas/settings.schema.json](./schemas/settings.schema.json), [schemas/hooks/](./schemas/hooks/)

## Invariants

1. YAML front matter on every topic `.md` file.
2. Verbatim excerpts only (minimal glue / research notes labeled explicitly).
3. No project-specific mapping sections in topic files.
4. English only; canonical URLs from geminicli.com llms.txt and ai.google.dev.
5. Update SOURCES.md when re-fetching after doc changes.

## Three hook surfaces

| Surface | Config | Events | agentd today? |
|---------|--------|--------|---------------|
| **Gemini CLI** | `.gemini/settings.json` → `hooks` | `BeforeTool`, `AfterAgent`, … | **Yes** — `--provider=gemini` |
| **Managed Agents API** | `.agents/hooks.json` in sandbox | `pre_tool_execution`, `post_tool_execution` | No |
| **Antigravity CLI** | `.agents/hooks.json` / plugin hooks | `PreToolUse`, `PostToolUse`, … | No — successor wire |

## agentd relevance

| Gemini surface | agentd touchpoint |
|----------------|-------------------|
| `settings.json` hooks | `agentd install --provider=gemini` → agenthooks `render_gemini.go` |
| stdin JSON → stdout `decision` | `agentd hook run --provider=gemini` |
| Policy engine / sandbox | `internal/guard` analog |
| Sessions / `transcript_path` | Trajectory import (layout partially undocumented) |
| Extension `gemini-extension.json` MCP | MCP guard; quirk #27 |
| Managed API `.agents/hooks.json` | Different wire — no current touchpoint |

## Gaps vs agentd / agenthooks assumptions

1. **`SessionStart` may not fire** — agenthooks documents open upstream regression; see hooks lifecycle sources.
2. **`decision: "ask"` on BeforeTool** — implemented in Gemini CLI (GitHub #28046) but omitted from hooks reference allow/deny table at snapshot.
3. **Exit codes ≠ docs (quirk #11)** — any non-zero except 1 can block; stderr parsed when stdout empty; agenthooks always writes explicit JSON stdout.
4. **Shallow-merge `tool_input` (quirk #12)** — key removal impossible; agenthooks surfaces `ErrLossyUpdate`.
5. **`additionalContext` HTML-escaping (quirk #13)** — documented loss in agenthooks DESIGN.md.
6. **Timeout milliseconds (quirk #14)** — Gemini uses ms; agenthooks converts from `time.Duration`.
7. **MCP naming `mcp_<server>_<tool>` (quirk #15)** — matcher compiler in agenthooks.
8. **Extension MCP in `gemini-extension.json` (quirk #27)** — not only `settings.json` mcpServers.
9. **Session file layout for trajectory** — official docs cover retention and `transcript_path`; nested rollout JSONL paths not found in reviewed sources.
10. **`T2-CODING-SETUP`** — ai.google.dev path 404 at snapshot; use agents overview + quickstart instead.

## Out of scope

Billing/pricing, themes, keyboard shortcuts, FAQ, troubleshooting, changelogs, pure UX tutorials, full Antigravity product dump (~60 pages) — migration delta only in `09-migration-and-antigravity/`.
