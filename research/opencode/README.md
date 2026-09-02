# OpenCode Documentation Research

Structured verbatim excerpts from [opencode.ai/docs](https://opencode.ai/docs/).

**Snapshot date: 2026-09-02** (see `opencode_docs_snapshot` in topic front matter).

OpenCode extends the agent loop via **plugins** (event hooks such as `tool.execute.before`, `permission.asked`, `session.idle`) — not a separate `hooks.json` product like Cursor or Codex.

## How to read

- Each topic file under `01-agent-loop/` … `08-enterprise/` starts with YAML front matter (`primary_sources`, `studied_at`, `opencode_docs_snapshot`, `applicability`).
- Body text is verbatim from primary sources under `### Source:` headings.
- Raw fetches: [`raw/fetched/`](./raw/fetched/). Machine-readable schemas: [`schemas/`](./schemas/).
- Coverage ledgers: [`MANIFEST.md`](./MANIFEST.md), [`CHECKLIST.md`](./CHECKLIST.md), [`SOURCES.md`](./SOURCES.md).

## Topic index

| Topic | File | Primary sources |
|-------|------|-----------------|
| Overview and getting started | [01-agent-loop/overview-and-getting-started.md](./01-agent-loop/overview-and-getting-started.md) | T1-INTRO |
| Built-in agents and modes | [01-agent-loop/agents-builtin-and-modes.md](./01-agent-loop/agents-builtin-and-modes.md) | T1-AGENTS |
| Agent configuration | [01-agent-loop/agent-configuration.md](./01-agent-loop/agent-configuration.md) | T1-AGENTS |
| Models and variants | [01-agent-loop/models-and-variants.md](./01-agent-loop/models-and-variants.md) | T1-MODELS |
| Compaction and context | [01-agent-loop/compaction-and-context.md](./01-agent-loop/compaction-and-context.md) | T1-PLUGINS, T1-CONFIG |
| Built-in slash commands | [01-agent-loop/built-in-slash-commands.md](./01-agent-loop/built-in-slash-commands.md) | T1-INTRO, T1-COMMANDS, T1-CLI |
| Plugins overview | [02-plugins/plugins-overview.md](./02-plugins/plugins-overview.md) | T1-PLUGINS |
| Plugin events — tool | [02-plugins/plugin-events-tool.md](./02-plugins/plugin-events-tool.md) | T1-PLUGINS, T1-TOOLS |
| Plugin events — permission | [02-plugins/plugin-events-permission.md](./02-plugins/plugin-events-permission.md) | T1-PLUGINS |
| Plugin events — session | [02-plugins/plugin-events-session.md](./02-plugins/plugin-events-session.md) | T1-PLUGINS |
| Plugin events — file, message, shell | [02-plugins/plugin-events-file-message-shell.md](./02-plugins/plugin-events-file-message-shell.md) | T1-PLUGINS |
| Plugin events — TUI | [02-plugins/plugin-events-tui.md](./02-plugins/plugin-events-tui.md) | T1-PLUGINS |
| Plugin events — generic handler | [02-plugins/plugin-events-generic-handler.md](./02-plugins/plugin-events-generic-handler.md) | T1-PLUGINS |
| Plugin examples and custom tools | [02-plugins/plugin-examples-and-custom-tools.md](./02-plugins/plugin-examples-and-custom-tools.md) | T1-PLUGINS |
| MCP configuration | [03-mcp/configuration.md](./03-mcp/configuration.md) | T1-MCP |
| MCP OAuth and management | [03-mcp/oauth-and-management.md](./03-mcp/oauth-and-management.md) | T1-MCP |
| MCP tools integration | [03-mcp/tools-integration.md](./03-mcp/tools-integration.md) | T1-MCP |
| Rules and AGENTS.md | [04-customization/rules-and-agents-md.md](./04-customization/rules-and-agents-md.md) | T1-RULES |
| Agent Skills | [04-customization/skills.md](./04-customization/skills.md) | T1-SKILLS |
| Commands | [04-customization/commands.md](./04-customization/commands.md) | T1-COMMANDS |
| Custom tools | [04-customization/custom-tools.md](./04-customization/custom-tools.md) | T1-CUSTOM-TOOLS |
| References | [04-customization/references.md](./04-customization/references.md) | T1-REFERENCES |
| Formatters | [04-customization/formatters.md](./04-customization/formatters.md) | T1-FORMATTERS |
| Themes and keybinds | [04-customization/themes-and-keybinds.md](./04-customization/themes-and-keybinds.md) | T1-THEMES, T1-KEYBINDS |
| Built-in tools | [04-customization/tools-built-in.md](./04-customization/tools-built-in.md) | T1-TOOLS |
| LSP servers | [04-customization/lsp-servers.md](./04-customization/lsp-servers.md) | T1-LSP |
| Config locations and precedence | [05-settings-and-paths/config-locations-and-precedence.md](./05-settings-and-paths/config-locations-and-precedence.md) | T1-CONFIG |
| Config reference — runtime | [05-settings-and-paths/config-reference-runtime.md](./05-settings-and-paths/config-reference-runtime.md) | T1-CONFIG |
| Config variables substitution | [05-settings-and-paths/config-variables-substitution.md](./05-settings-and-paths/config-variables-substitution.md) | T1-CONFIG |
| Config reference — experimental | [05-settings-and-paths/config-reference-experimental.md](./05-settings-and-paths/config-reference-experimental.md) | T1-CONFIG |
| Auth and credentials | [05-settings-and-paths/auth-and-credentials.md](./05-settings-and-paths/auth-and-credentials.md) | T1-CLI, T1-PROVIDERS, T1-GITLAB |
| Project vs global matrix | [05-settings-and-paths/project-vs-global-matrix.md](./05-settings-and-paths/project-vs-global-matrix.md) | T1-CONFIG, T1-RULES, T1-PLUGINS, T1-SKILLS, T1-TROUBLESHOOTING |
| Environment variables | [05-settings-and-paths/environment-variables.md](./05-settings-and-paths/environment-variables.md) | T1-CLI, T1-CONFIG |
| Storage and runtime paths | [05-settings-and-paths/storage-and-runtime-paths.md](./05-settings-and-paths/storage-and-runtime-paths.md) | T1-TROUBLESHOOTING |
| Permissions | [05-settings-and-paths/permissions.md](./05-settings-and-paths/permissions.md) | T1-PERMISSIONS |
| Policies | [05-settings-and-paths/policies.md](./05-settings-and-paths/policies.md) | T1-POLICIES |
| Providers overview | [05-settings-and-paths/providers-overview.md](./05-settings-and-paths/providers-overview.md) | T1-PROVIDERS |
| CLI reference | [06-cli-and-sdk/cli-reference.md](./06-cli-and-sdk/cli-reference.md) | T1-CLI |
| Headless run and attach | [06-cli-and-sdk/headless-run-and-attach.md](./06-cli-and-sdk/headless-run-and-attach.md) | T1-CLI |
| Session export and import | [06-cli-and-sdk/session-export-import.md](./06-cli-and-sdk/session-export-import.md) | T1-CLI |
| Server and HTTP API | [06-cli-and-sdk/server-and-http-api.md](./06-cli-and-sdk/server-and-http-api.md) | T1-SERVER |
| SDK TypeScript | [06-cli-and-sdk/sdk-typescript.md](./06-cli-and-sdk/sdk-typescript.md) | T1-SDK |
| ACP support | [06-cli-and-sdk/acp.md](./06-cli-and-sdk/acp.md) | T1-ACP |
| GitHub | [07-integrations/github.md](./07-integrations/github.md) | T1-GITHUB |
| GitLab | [07-integrations/gitlab.md](./07-integrations/gitlab.md) | T1-GITLAB |
| Zen provider | [07-integrations/zen-provider.md](./07-integrations/zen-provider.md) | T1-ZEN |
| Share | [07-integrations/share.md](./07-integrations/share.md) | T1-SHARE |
| Network and proxy | [07-integrations/network-and-proxy.md](./07-integrations/network-and-proxy.md) | T1-NETWORK |
| Ecosystem index | [07-integrations/ecosystem-index.md](./07-integrations/ecosystem-index.md) | T1-ECOSYSTEM |
| Troubleshooting | [07-integrations/troubleshooting.md](./07-integrations/troubleshooting.md) | T1-TROUBLESHOOTING |
| Enterprise | [08-enterprise/enterprise.md](./08-enterprise/enterprise.md) | T1-ENTERPRISE |
| Managed config | [08-enterprise/managed-config.md](./08-enterprise/managed-config.md) | T1-CONFIG |

## Meta

- Source registry: [SOURCES.md](./SOURCES.md)
- File checklist: [MANIFEST.md](./MANIFEST.md)
- Coverage checklists: [CHECKLIST.md](./CHECKLIST.md)
- Topic file pattern: [TEMPLATE.md](./TEMPLATE.md)
- Schemas: [schemas/config.schema.json](./schemas/config.schema.json), [schemas/tui.schema.json](./schemas/tui.schema.json)

## Invariants

1. YAML front matter on every topic `.md` file.
2. Verbatim excerpts only (minimal glue).
3. No project-specific mapping sections in topic files.
4. English only; canonical URLs from opencode.ai/docs.
5. Update [SOURCES.md](./SOURCES.md) when re-fetching after doc changes.

## agentd relevance

Reference only — not duplicated inside topic files.

| OpenCode surface | agentd touchpoint |
|------------------|-------------------|
| Plugin events (`tool.execute.before`, `permission.*`, `session.idle`) | [`hook serve`](../../internal/hookedge/serve.go) NDJSON via [agenthooks](https://github.com/speakeasy-api/agenthooks); not `hooks.json` |
| Install path `.opencode/plugin/agenthooks.ts` | [`agentd install --provider=opencode`](../../docs/en/providers-opencode.md); **project-only** |
| `permission` / tool gates | [`internal/guard`](../../internal/guard/); OpenCode `tool.pre` = deny/update-input only (no Ask) |
| `AGENTS.md` + `instructions` | Prompt context; parallels other providers' rules |
| Skills (`skill` tool) | On-demand SKILL.md; compare [research/cursor/04-customization/skills.md](../cursor/04-customization/skills.md) |
| MCP `mcp` config | Guard/MCP checks; OAuth at `~/.local/share/opencode/mcp-auth.json` |
| Config merge + precedence | [`internal/config`](../../internal/config/) overlay model |
| Server `/session/*/message`, `/event` SSE | Trajectory/export research; importer **none** for opencode today |
| Managed config / enterprise | Org deployment constraints |

## Gaps vs agentd assumptions

1. **Hooks naming**: OpenCode docs say “plugins”; agentd/agenthooks say “hook serve”.
2. **Plugin directory**: official `.opencode/plugins/` (plural) vs agentd install `.opencode/plugin/` (singular).
3. **No user-scope install** for OpenCode ([`internal/install/target.go`](../../internal/install/target.go)).
4. **Permission model**: OpenCode native allow/ask/deny UI vs agentd Decide on serve path.
5. **No `llms.txt`** — manual source catalog.
6. **Trajectory import**: `ImporterNone` for opencode.
7. **Plugin event I/O schemas**: not published by OpenCode; wire shapes in agenthooks.
8. **Legacy `tools` config**: deprecated → `permission` (v1.1.1+).
9. **`modes/` directory**: in config; included in matrix.
10. **Testing isolation**: `--pure` / `OPENCODE_DISABLE_DEFAULT_PLUGINS` when debugging agenthooks alongside other plugins.
11. **MCP doc slug**: `/docs/mcp-servers`, not `/docs/mcp`.

## Out of scope (this corpus)

- Full `providers.md` per-provider catalog (essentials only)
- Full Zen model endpoint table (gateway + `/connect` only)
- agenthooks NDJSON protocol (cross-ref only)
- P2 usage pages TUI/Web/IDE/Go (raw cached; skipped-overlap in SOURCES)
