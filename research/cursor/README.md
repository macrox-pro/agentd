# Cursor Documentation Research

Structured verbatim excerpts from [cursor.com/docs](https://cursor.com/docs) and selected help pages.

**Snapshot date: 2026-08-25** (see `cursor_docs_snapshot` in topic front matter).

## How to read

- Each topic file under `01-agent-loop/` … `09-enterprise/` starts with YAML front matter (`primary_sources`, `studied_at`, `cursor_docs_snapshot`, `applicability`).
- Body text is verbatim from primary sources; see `### Source:` headings.
- Machine-readable schemas live in [`schemas/`](./schemas/).

## Topic index

| Topic | File | Primary sources |
|-------|------|-----------------|
| Agent overview and modes | [01-agent-loop/overview-and-modes.md](./01-agent-loop/overview-and-modes.md) | T1-AGENT-OVERVIEW, T1-AGENT-PLAN, T1-AGENT-DEBUG, T1-AGENT-DESIGN, T1-AGENT-PROMPT |
| Agent tools | [01-agent-loop/tools-terminal-browser.md](./01-agent-loop/tools-terminal-browser.md) | T1-AGENT-TERMINAL, T1-AGENT-BROWSER, T1-AGENT-SEARCH, T1-AGENT-CANVAS, T1-AGENTS-WINDOW, T1-AGENT-REVIEW |
| Agent security and run modes | [01-agent-loop/security-and-run-modes.md](./01-agent-loop/security-and-run-modes.md) | T1-AGENT-SECURITY, T1-RUN-MODES |
| Subagents | [01-agent-loop/subagents.md](./01-agent-loop/subagents.md) | T1-SUBAGENTS |
| Hooks overview | [02-hooks/hooks-overview.md](./02-hooks/hooks-overview.md) | T1-HOOKS |
| Agent hook events | [02-hooks/hook-events-agent.md](./02-hooks/hook-events-agent.md) | T1-HOOKS |
| Tab and lifecycle hooks | [02-hooks/hook-events-tab-and-lifecycle.md](./02-hooks/hook-events-tab-and-lifecycle.md) | T1-HOOKS |
| Hook execution and protocol | [02-hooks/execution-and-protocol.md](./02-hooks/execution-and-protocol.md) | T1-HOOKS |
| Hook config and precedence | [02-hooks/config-and-precedence.md](./02-hooks/config-and-precedence.md) | T1-HOOKS |
| Cloud hook support matrix | [02-hooks/cloud-support-matrix.md](./02-hooks/cloud-support-matrix.md) | T1-HOOKS |
| Hook examples | [02-hooks/hooks-examples.md](./02-hooks/hooks-examples.md) | T1-HOOKS |
| Partner integrations | [02-hooks/partner-integrations.md](./02-hooks/partner-integrations.md) | T1-HOOKS |
| Third-party hooks | [02-hooks/third-party-hooks.md](./02-hooks/third-party-hooks.md) | T1-THIRD-PARTY-HOOKS |
| MCP configuration | [03-mcp/configuration.md](./03-mcp/configuration.md) | T1-MCP |
| MCP transports and features | [03-mcp/transports-auth-features.md](./03-mcp/transports-auth-features.md) | T1-MCP |
| MCP install links | [03-mcp/install-links.md](./03-mcp/install-links.md) | T1-MCP-INSTALL |
| MCP enterprise policy | [03-mcp/enterprise-policy.md](./03-mcp/enterprise-policy.md) | T1-ENT-MCP, T1-PERMISSIONS |
| Rules and AGENTS.md | [04-customization/rules-and-agents-md.md](./04-customization/rules-and-agents-md.md) | T1-RULES |
| Skills | [04-customization/skills.md](./04-customization/skills.md) | T1-SKILLS |
| Plugins and marketplace | [04-customization/plugins-and-marketplace.md](./04-customization/plugins-and-marketplace.md) | T1-PLUGINS, T1-PLUGINS-REF |
| Customize Cursor | [04-customization/customize-cursor.md](./04-customization/customize-cursor.md) | T1-CUSTOMIZE |
| Deeplinks | [04-customization/deeplinks.md](./04-customization/deeplinks.md) | T1-DEEPLINKS |
| Project vs global matrix | [05-settings-and-paths/project-vs-global-matrix.md](./05-settings-and-paths/project-vs-global-matrix.md) | T1-HOOKS, T1-RULES |
| CLI config and permissions | [05-settings-and-paths/cli-config-and-permissions.md](./05-settings-and-paths/cli-config-and-permissions.md) | T1-CLI-CONFIG, T1-CLI-PERMS |
| Ignore files and worktrees | [05-settings-and-paths/ignore-and-worktrees.md](./05-settings-and-paths/ignore-and-worktrees.md) | T1-IGNORE, T1-WORKTREES |
| Cloud environment | [05-settings-and-paths/cloud-environment.md](./05-settings-and-paths/cloud-environment.md) | T1-CLOUD-SETUP |
| IDE settings gaps | [05-settings-and-paths/ide-settings-gaps.md](./05-settings-and-paths/ide-settings-gaps.md) | T1-IGNORE, T1-RULES |
| Build an AI coding agent | [06-cli-and-sdk/build-ai-coding-agent.md](./06-cli-and-sdk/build-ai-coding-agent.md) | T2-HELP-BUILD |
| CLI overview and usage | [06-cli-and-sdk/cli-overview-and-using.md](./06-cli-and-sdk/cli-overview-and-using.md) | T1-CLI-INSTALL, T1-CLI-OVERVIEW, T1-CLI-USING |
| Headless and shell mode | [06-cli-and-sdk/headless-and-shell-mode.md](./06-cli-and-sdk/headless-and-shell-mode.md) | T1-CLI-HEADLESS, T1-CLI-SHELL |
| Agent Client Protocol | [06-cli-and-sdk/acp.md](./06-cli-and-sdk/acp.md) | T1-CLI-ACP |
| CLI GitHub Actions | [06-cli-and-sdk/github-actions.md](./06-cli-and-sdk/github-actions.md) | T1-CLI-GHA |
| CLI reference | [06-cli-and-sdk/cli-reference.md](./06-cli-and-sdk/cli-reference.md) | T1-CLI-PARAMS … T1-CLI-TERM |
| Cursor SDK | [06-cli-and-sdk/sdk-typescript-python-bridge.md](./06-cli-and-sdk/sdk-typescript-python-bridge.md) | T1-SDK-TS, T1-SDK-PY, T1-SDK-BRIDGE, T1-SDK-CHANGELOG |
| Cloud overview | [07-cloud-agents/overview-setup-builds.md](./07-cloud-agents/overview-setup-builds.md) | T1-CLOUD, T1-CLOUD-SETUP, T1-CLOUD-BUILDS, T1-CLOUD-BP, T1-CLOUD-SETTINGS |
| Cloud capabilities | [07-cloud-agents/capabilities-metadata.md](./07-cloud-agents/capabilities-metadata.md) | T1-CLOUD-CAP, T1-CLOUD-META |
| Cloud automations | [07-cloud-agents/automations-and-webhooks.md](./07-cloud-agents/automations-and-webhooks.md) | T1-CLOUD-AUTO, T1-CLOUD-WH |
| Cloud Agents API | [07-cloud-agents/api-endpoints.md](./07-cloud-agents/api-endpoints.md) | T1-CLOUD-API, T1-API |
| Cloud security | [07-cloud-agents/security-network-identity.md](./07-cloud-agents/security-network-identity.md) | T1-CLOUD-SEC, T1-CLOUD-NET, T1-CLOUD-ID, T1-CLOUD-PRIV |
| Self-hosted agents | [07-cloud-agents/self-hosted-agents.md](./07-cloud-agents/self-hosted-agents.md) | T1-CLOUD-SELF |
| Bugbot | [08-review-and-automation/bugbot.md](./08-review-and-automation/bugbot.md) | T1-BUGBOT |
| Security agents | [08-review-and-automation/security-agents.md](./08-review-and-automation/security-agents.md) | T1-SEC-AGENTS |
| Approval agents | [08-review-and-automation/approval-agents.md](./08-review-and-automation/approval-agents.md) | T1-APPROVAL |
| GitHub integration | [08-review-and-automation/integrations-github.md](./08-review-and-automation/integrations-github.md) | T1-GITHUB |
| Enterprise overview | [09-enterprise/overview-and-llm-controls.md](./09-enterprise/overview-and-llm-controls.md) | T1-ENTERPRISE, T1-LLM-SAFETY |
| Endpoint security | [09-enterprise/endpoint-security-and-network.md](./09-enterprise/endpoint-security-and-network.md) | T1-EP-SEC, T1-NET-CFG |
| Analytics and Admin API | [09-enterprise/analytics-admin-api.md](./09-enterprise/analytics-admin-api.md) | T1-ANALYTICS, T1-ADMIN |

## Meta

- Source registry: [SOURCES.md](./SOURCES.md)
- Topic file pattern: [TEMPLATE.md](./TEMPLATE.md)
- Schemas: [schemas/cloud-agents-openapi.yaml](./schemas/cloud-agents-openapi.yaml), [schemas/environment.schema.json](./schemas/environment.schema.json)

## Invariants

1. YAML front matter on every topic `.md` file.
2. Verbatim excerpts only (minimal glue).
3. No project-specific mapping sections in topic files.
4. English only; canonical URLs from [cursor.com/llms.txt](https://cursor.com/llms.txt).
5. Update [SOURCES.md](./SOURCES.md) when re-fetching after doc changes.

## agentd relevance

Reference only — not duplicated inside topic files.

| Cursor surface | agentd touchpoint |
|----------------|-------------------|
| Native hooks (`.cursor/hooks.json`, stdio JSON) | Analog to hook lifecycle; different wire than `agentd hook run --provider=cursor --argv-payload` |
| Install root | Cursor user scope → `~/.cursor` (`internal/install/dir.go`) |
| Policy layers | Hooks/MCP intercept at runtime; rules/skills inject prompt context |
| Cloud agents | User hooks unavailable in cloud VMs; project `.cursor/hooks.json` only |
| Trajectory | Partial Cursor transcript import: `agentd session import --provider cursor --path …` |
| Guards | `beforeMCPExecution`, `beforeShellExecution` mirror agentd guard concerns |

## Help center supplement

Help pages under `cursor.com/help/` were fetched (see SOURCES `T2-HELP-*`). Content largely overlaps `/docs/`; unique narrative retained in [build-ai-coding-agent.md](./06-cli-and-sdk/build-ai-coding-agent.md). Other help pages marked `done` in SOURCES without separate topic files (no unique facts beyond `/docs/`).

## Out of scope

Billing/pricing, per-model pages, i18n locales, Grok Bot help, Origin docs.
