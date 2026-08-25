# Claude Code Documentation Research

Structured verbatim excerpts from [code.claude.com/docs](https://code.claude.com/docs).

**Snapshot date: 2026-08-25** (see `claude_code_docs_snapshot` in topic front matter).

## How to read

- Each topic file under `01-agent-loop/` … `08-enterprise/` starts with YAML front matter (`primary_sources`, `studied_at`, `claude_code_docs_snapshot`, `applicability`).
- Body text is verbatim from primary sources; see `### Source:` headings.
- Machine-readable schemas live in [`schemas/`](./schemas/).

## Topic index

| Topic | File | Primary sources |
|-------|------|-----------------|
| Context window and compaction | [01-agent-loop/context-window-and-compaction.md](./01-agent-loop/context-window-and-compaction.md) | T1-CONTEXT, T1-PROMPT-CACHE |
| Overview and agent loop | [01-agent-loop/overview-and-loop.md](./01-agent-loop/overview-and-loop.md) | T1-AGENT-LOOP |
| Permissions and modes | [01-agent-loop/permissions-and-modes.md](./01-agent-loop/permissions-and-modes.md) | T1-PERMISSIONS, T1-PERM-MODES |
| Sandbox and checkpoints | [01-agent-loop/sandbox-and-checkpoints.md](./01-agent-loop/sandbox-and-checkpoints.md) | T1-SANDBOX, T1-SANDBOX-ENV, T1-CHECKPOINT |
| Subagents | [01-agent-loop/subagents.md](./01-agent-loop/subagents.md) | T1-SUBAGENTS |
| Tools reference | [01-agent-loop/tools-reference-pointer.md](./01-agent-loop/tools-reference-pointer.md) | T1-TOOLS-REF |
| Cloud hooks matrix | [02-hooks/cloud-hooks-matrix.md](./02-hooks/cloud-hooks-matrix.md) | T1-HOOKS, T1-CLOUD-ENV |
| Hook configuration and locations | [02-hooks/config-and-locations.md](./02-hooks/config-and-locations.md) | T1-HOOKS |
| Hook debug and troubleshooting | [02-hooks/debug-and-troubleshooting.md](./02-hooks/debug-and-troubleshooting.md) | T1-DEBUG, T1-HOOKS-GUIDE |
| Hook events — compaction | [02-hooks/events-compact.md](./02-hooks/events-compact.md) | T1-HOOKS |
| Hook events — files, worktrees, instructions | [02-hooks/events-files-worktree.md](./02-hooks/events-files-worktree.md) | T1-HOOKS |
| Hook events — MCP elicitation | [02-hooks/events-mcp-elicitation.md](./02-hooks/events-mcp-elicitation.md) | T1-HOOKS |
| Hook events — prompt and stop | [02-hooks/events-prompt-stop.md](./02-hooks/events-prompt-stop.md) | T1-HOOKS |
| Hook events — session and setup | [02-hooks/events-session.md](./02-hooks/events-session.md) | T1-HOOKS |
| Hook events — subagents and tasks | [02-hooks/events-subagent-tasks.md](./02-hooks/events-subagent-tasks.md) | T1-HOOKS |
| Hook events — tools and permissions | [02-hooks/events-tools.md](./02-hooks/events-tools.md) | T1-HOOKS |
| Hook execution and protocol | [02-hooks/execution-and-protocol.md](./02-hooks/execution-and-protocol.md) | T1-HOOKS |
| Hooks overview | [02-hooks/hooks-overview.md](./02-hooks/hooks-overview.md) | T1-HOOKS-GUIDE |
| Matcher patterns | [02-hooks/matcher-patterns.md](./02-hooks/matcher-patterns.md) | T1-HOOKS, T1-DEBUG |
| Plugin and managed hooks | [02-hooks/plugin-and-managed-hooks.md](./02-hooks/plugin-and-managed-hooks.md) | T1-PLUGINS-REF, T1-MANAGED |
| Hook schemas | [02-hooks/schemas.md](./02-hooks/schemas.md) | T1-HOOKS |
| Channels | [03-mcp/channels.md](./03-mcp/channels.md) | T1-CHANNELS, T1-CHANNELS-REF |
| MCP configuration | [03-mcp/configuration.md](./03-mcp/configuration.md) | T1-MCP |
| MCP debug and approval | [03-mcp/debug-and-approval.md](./03-mcp/debug-and-approval.md) | T1-MCP, T1-DEBUG |
| Plugin MCP | [03-mcp/plugin-mcp.md](./03-mcp/plugin-mcp.md) | T1-PLUGINS-REF |
| MCP scopes and paths | [03-mcp/scopes-and-paths.md](./03-mcp/scopes-and-paths.md) | T1-MCP-QS |
| MCP transports, auth, and tool search | [03-mcp/transports-auth-tool-search.md](./03-mcp/transports-auth-tool-search.md) | T1-MCP |
| Features overview | [04-customization/features-overview.md](./04-customization/features-overview.md) | T1-FEATURES |
| Large codebases | [04-customization/large-codebases.md](./04-customization/large-codebases.md) | T1-LARGE-CODEBASES |
| Memory and CLAUDE.md | [04-customization/memory-and-claude-md.md](./04-customization/memory-and-claude-md.md) | T1-MEMORY |
| Discover and install plugins | [04-customization/discover-plugins.md](./04-customization/discover-plugins.md) | T1-DISCOVER-PLUGINS |
| Plugin marketplaces | [04-customization/plugin-marketplaces.md](./04-customization/plugin-marketplaces.md) | T2-PLUGIN-MKT |
| Plugins | [04-customization/plugins.md](./04-customization/plugins.md) | T1-PLUGINS |
| Rules | [04-customization/rules.md](./04-customization/rules.md) | T1-MEMORY |
| Skills | [04-customization/skills.md](./04-customization/skills.md) | T1-SKILLS |
| Workflows | [04-customization/workflows.md](./04-customization/workflows.md) | T1-WORKFLOWS |
| `.claude` directory reference | [05-settings-and-paths/claude-directory-reference.md](./05-settings-and-paths/claude-directory-reference.md) | T1-CLAUDE-DIR |
| claude.json vs settings.json | [05-settings-and-paths/claude-json-vs-settings-json.md](./05-settings-and-paths/claude-json-vs-settings-json.md) | T1-DEBUG, T1-CLAUDE-DIR |
| Debug configuration commands | [05-settings-and-paths/debug-config-commands.md](./05-settings-and-paths/debug-config-commands.md) | T1-DEBUG |
| Environment variables | [05-settings-and-paths/environment-variables.md](./05-settings-and-paths/environment-variables.md) | T1-ENV-VARS |
| Project vs global path matrix | [05-settings-and-paths/project-vs-global-matrix.md](./05-settings-and-paths/project-vs-global-matrix.md) | T1-SETTINGS, T1-HOOKS, T1-MCP, T1-MEMORY, T1-SKILLS, T1-SUBAGENTS |
| Sessions and transcripts | [05-settings-and-paths/sessions-and-transcripts.md](./05-settings-and-paths/sessions-and-transcripts.md) | T1-SESSIONS, T1-AGENT-LOOP |
| Settings basics and precedence | [05-settings-and-paths/settings-basics-and-precedence.md](./05-settings-and-paths/settings-basics-and-precedence.md) | T1-SETTINGS |
| Settings reference — hooks and permissions | [05-settings-and-paths/settings-reference-hooks-permissions.md](./05-settings-and-paths/settings-reference-hooks-permissions.md) | T1-SETTINGS-REF |
| Settings reference — other keys | [05-settings-and-paths/settings-reference-other.md](./05-settings-and-paths/settings-reference-other.md) | T1-SETTINGS-REF, T1-SETTINGS-EXAMPLE |
| Agent SDK hooks and permissions | [06-cli-and-sdk/agent-sdk-hooks-permissions.md](./06-cli-and-sdk/agent-sdk-hooks-permissions.md) | T1-SDK-HOOKS, T1-SDK-PERMS, T1-SDK-FEATURES |
| CLI reference — hooks flags | [06-cli-and-sdk/cli-reference-hooks-flags.md](./06-cli-and-sdk/cli-reference-hooks-flags.md) | T1-CLI |
| Commands reference | [06-cli-and-sdk/commands-reference.md](./06-cli-and-sdk/commands-reference.md) | T2-COMMANDS |
| Deep links | [06-cli-and-sdk/deep-links.md](./06-cli-and-sdk/deep-links.md) | T2-DEEP-LINKS |
| Errors reference | [06-cli-and-sdk/errors-reference.md](./06-cli-and-sdk/errors-reference.md) | T2-ERRORS |
| GitHub Actions | [06-cli-and-sdk/github-actions.md](./06-cli-and-sdk/github-actions.md) | T2-GHA |
| Headless mode | [06-cli-and-sdk/headless-and-init.md](./06-cli-and-sdk/headless-and-init.md) | T1-HEADLESS |
| Interactive mode | [06-cli-and-sdk/interactive-mode.md](./06-cli-and-sdk/interactive-mode.md) | T2-INTERACTIVE |
| Output styles | [06-cli-and-sdk/output-styles.md](./06-cli-and-sdk/output-styles.md) | T2-OUTPUT-STYLES |
| Troubleshooting | [06-cli-and-sdk/troubleshooting.md](./06-cli-and-sdk/troubleshooting.md) | T2-TROUBLESHOOTING |
| Claude Code on the web | [07-cloud-and-environments/claude-code-on-web.md](./07-cloud-and-environments/claude-code-on-web.md) | T1-CLOUD-WEB |
| Cloud environments | [07-cloud-and-environments/cloud-environments.md](./07-cloud-and-environments/cloud-environments.md) | T1-CLOUD-ENV |
| Self-hosted environments | [07-cloud-and-environments/self-hosted-environments.md](./07-cloud-and-environments/self-hosted-environments.md) | T1-SELF-HOSTED |
| Worktrees | [07-cloud-and-environments/worktrees.md](./07-cloud-and-environments/worktrees.md) | T1-WORKTREES |
| Admin setup | [08-enterprise/admin-setup.md](./08-enterprise/admin-setup.md) | T1-ADMIN |
| Auto mode config | [08-enterprise/auto-mode-config.md](./08-enterprise/auto-mode-config.md) | T1-AUTO-MODE |
| Managed MCP | [08-enterprise/managed-mcp.md](./08-enterprise/managed-mcp.md) | T1-MANAGED-MCP |
| Managed settings | [08-enterprise/managed-settings.md](./08-enterprise/managed-settings.md) | T1-MANAGED, T1-SERVER-MANAGED |
| Monitoring and OTEL | [08-enterprise/monitoring-usage.md](./08-enterprise/monitoring-usage.md) | T2-MONITORING |
| Security | [08-enterprise/security.md](./08-enterprise/security.md) | T1-SECURITY |

## Meta

- Source registry: [SOURCES.md](./SOURCES.md)
- Topic file pattern: [TEMPLATE.md](./TEMPLATE.md)
- Schemas: [schemas/hooks/](./schemas/hooks/) (31 input example JSON files)

## Invariants

1. YAML front matter on every topic `.md` file.
2. Verbatim excerpts only (minimal glue).
3. No project-specific mapping sections in topic files.
4. English only; canonical URLs from [code.claude.com/docs/llms.txt](https://code.claude.com/docs/llms.txt).
5. Update [SOURCES.md](./SOURCES.md) when re-fetching after doc changes.

## agentd relevance

Reference only — not duplicated inside topic files.

| Claude surface | agentd touchpoint |
|----------------|-------------------|
| `.claude/settings.json` hooks (`"hooks"` key) | `agentd install --provider=claude-code --scope=project` → `.claude/settings.json` |
| `~/.claude/settings.json` | `--scope=user` |
| Plugin `hooks/hooks.json` | `--scope=plugin` → `.claude-plugin/plugin.json` |
| Hook wire: stdin JSON → stdout JSON; `{}` + exit 0 | `agentd hook run --provider=claude-code` |
| Events: `PreToolUse`, `UserPromptSubmit`, `Stop`, … | agenthooks capability matrix; blocking vs async |
| `.mcp.json` / `~/.claude.json` MCP | guard `mcp`; separate from settings hooks |
| `~/.claude/projects/**/transcript.jsonl` | trajectory import (`internal/trajectory/importer`) |
| Channels (MCP → session) | async side design reference |
| Security safeguards | `internal/guard` cross-check |
| Large codebase / monorepo | nested CLAUDE.md, per-package skills — install scope |
| `/hooks`, `/mcp`, `/doctor` | debug installed hooks and MCP |
| Cloud: no local `~/.claude/settings.json` | install scope implications |

## Gaps vs agentd assumptions

Facts recorded when official docs omit or conflict:

1. **`~/.claude.json` vs `~/.claude/settings.json`** — `permissions`, `hooks`, and `env` belong in `settings.json`; `~/.claude.json` holds app state and user MCP UI config.
2. **No standalone project `hooks.json`** (except plugins) — hooks live under `"hooks"` in settings files; plugins use `hooks/hooks.json`.
3. **Matcher comma vs pipe** — before v2.1.191, comma in matcher is literal regex, not a list separator; use `|`.
4. **Cloud hook sources** — web/cloud sessions use repo hooks + server-managed settings; local user settings not loaded.
5. **Open-source hook JSON Schema** — no published generated schema repo found (unlike Codex); docs + `schemas/hooks/*.example.json` are the reference.
6. **Event name parity** — Claude `PreToolUse` maps to agenthooks `KindToolPre`; wire field names differ slightly by provider.

## Out of scope

See [SOURCES.md](./SOURCES.md) skip table (whats-new, desktop UI, gateways, full SDK API, billing/legal, glossary, i18n).
