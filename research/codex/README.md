# Codex Documentation Research

Structured verbatim excerpts from [learn.chatgpt.com/docs](https://learn.chatgpt.com/docs) (Codex / ChatGPT Learn).

**Snapshot date: 2026-08-25** (see `codex_docs_snapshot` in topic front matter).

## How to read

- Each topic file under `01-agent-loop/` … `08-enterprise/` starts with YAML front matter (`primary_sources`, `studied_at`, `codex_docs_snapshot`, `applicability`).
- Body text is verbatim from primary sources; see `### Source:` headings.
- Machine-readable schemas live in [`schemas/`](./schemas/).

## Topic index

| Topic | File | Primary sources |
|-------|------|-----------------|
| Sandbox and approvals | [01-agent-loop/sandbox-and-approvals.md](./01-agent-loop/sandbox-and-approvals.md) | T1-SANDBOX, T1-APPROVALS, T1-WIN-SANDBOX, T1-WSL |
| Permissions and profiles | [01-agent-loop/permissions-and-profiles.md](./01-agent-loop/permissions-and-profiles.md) | T1-PERMISSIONS, T1-PERM-MODES |
| Auto-review | [01-agent-loop/auto-review.md](./01-agent-loop/auto-review.md) | T1-AUTO-REVIEW |
| Hooks overview | [02-hooks/hooks-overview.md](./02-hooks/hooks-overview.md) | T1-HOOKS |
| Review, trust, and config shape | [02-hooks/config-and-trust.md](./02-hooks/config-and-trust.md) | T1-HOOKS |
| MCP tool hooks | [02-hooks/mcp-tool-hooks.md](./02-hooks/mcp-tool-hooks.md) | T1-HOOKS |
| Matcher patterns and hook protocol | [02-hooks/execution-and-protocol.md](./02-hooks/execution-and-protocol.md) | T1-HOOKS |
| Background hooks | [02-hooks/background-hooks.md](./02-hooks/background-hooks.md) | T1-HOOKS |
| Notifications (`notify` argv) | [02-hooks/notify.md](./02-hooks/notify.md) | T1-CFG-ADV |
| Notifications (UI surfaces) | [02-hooks/notifications-ui.md](./02-hooks/notifications-ui.md) | T2-NOTIFICATIONS |
| Session and subagent hook events | [02-hooks/events-session.md](./02-hooks/events-session.md) | T1-HOOKS |
| Tool hook events | [02-hooks/events-tools.md](./02-hooks/events-tools.md) | T1-HOOKS |
| Compact, prompt, and stop hook events | [02-hooks/events-compact-prompt-stop.md](./02-hooks/events-compact-prompt-stop.md) | T1-HOOKS |
| Managed and plugin-bundled hooks | [02-hooks/plugin-and-managed.md](./02-hooks/plugin-and-managed.md) | T1-HOOKS |
| Hook schemas pointer | [02-hooks/schemas.md](./02-hooks/schemas.md) | T1-HOOKS, T1-OPEN-SOURCE |
| Cloud hooks matrix | [02-hooks/cloud-hooks-matrix.md](./02-hooks/cloud-hooks-matrix.md) | T1-CLOUD, T1-ENV-MODES, T1-MANAGED, T2-REMOTE |
| MCP configuration | [03-mcp/configuration.md](./03-mcp/configuration.md) | T1-MCP |
| MCP transports and auth | [03-mcp/transports-auth.md](./03-mcp/transports-auth.md) | T1-MCP |
| Plugin-provided MCP servers | [03-mcp/plugin-mcp.md](./03-mcp/plugin-mcp.md) | T1-MCP |
| Customization overview | [04-customization/overview.md](./04-customization/overview.md) | T1-CUSTOM |
| Custom instructions with AGENTS.md | [04-customization/agents-md.md](./04-customization/agents-md.md) | T1-AGENTS-MD |
| Rules | [04-customization/rules.md](./04-customization/rules.md) | T1-RULES |
| Build skills | [04-customization/skills.md](./04-customization/skills.md) | T1-SKILLS |
| Skills & Plugins | [04-customization/skills-and-plugins.md](./04-customization/skills-and-plugins.md) | T1-SKILLS-PLUGINS |
| Plugins | [04-customization/plugins.md](./04-customization/plugins.md) | T1-PLUGINS, T1-BUILD-PLUGINS |
| Subagents | [04-customization/subagents.md](./04-customization/subagents.md) | T1-SUBAGENTS |
| Memories | [04-customization/memories.md](./04-customization/memories.md) | T2-MEMORIES |
| Record & Replay | [04-customization/record-and-replay.md](./04-customization/record-and-replay.md) | T2-RECORD |
| Import from another agent | [04-customization/import-from-another-agent.md](./04-customization/import-from-another-agent.md) | T2-IMPORT |
| Custom prompts (deprecated) | [04-customization/custom-prompts.md](./04-customization/custom-prompts.md) | T2-CUSTOM-PROMPTS |
| GitHub code review rules | [04-customization/github-code-review-rules.md](./04-customization/github-code-review-rules.md) | T2-GH-REVIEW |
| Project vs global path matrix | [05-settings-and-paths/project-vs-global-matrix.md](./05-settings-and-paths/project-vs-global-matrix.md) | T1-CFG-BASIC, T1-CFG-ADV, T1-HOOKS, … |
| Config basics and precedence | [05-settings-and-paths/config-basics-and-precedence.md](./05-settings-and-paths/config-basics-and-precedence.md) | T1-CFG-BASIC |
| Advanced configuration | [05-settings-and-paths/config-advanced.md](./05-settings-and-paths/config-advanced.md) | T1-CFG-ADV |
| Config reference — layers | [05-settings-and-paths/config-reference-layers.md](./05-settings-and-paths/config-reference-layers.md) | T1-CFG-REF |
| Config reference — hooks and MCP | [05-settings-and-paths/config-reference-hooks-mcp.md](./05-settings-and-paths/config-reference-hooks-mcp.md) | T1-CFG-REF |
| Config reference — sandbox | [05-settings-and-paths/config-reference-sandbox.md](./05-settings-and-paths/config-reference-sandbox.md) | T1-CFG-REF |
| Config reference — complete page | [05-settings-and-paths/config-reference-other.md](./05-settings-and-paths/config-reference-other.md) | T1-CFG-REF |
| Environment variables | [05-settings-and-paths/environment-variables.md](./05-settings-and-paths/environment-variables.md) | T1-CFG-ENV, T1-AUTH |
| Sample configuration | [05-settings-and-paths/config-sample.md](./05-settings-and-paths/config-sample.md) | T1-CFG-SAMPLE |
| Sessions and history | [05-settings-and-paths/sessions-and-history.md](./05-settings-and-paths/sessions-and-history.md) | T1-CFG-ADV |
| IDE and developer settings | [05-settings-and-paths/ide-settings.md](./05-settings-and-paths/ide-settings.md) | T1-DEV-SETTINGS |
| Feature maturity | [05-settings-and-paths/feature-maturity.md](./05-settings-and-paths/feature-maturity.md) | T2-FEATURE-MATURITY |
| Speed | [05-settings-and-paths/speed.md](./05-settings-and-paths/speed.md) | T2-SPEED |
| Codex CLI | [06-cli-and-sdk/cli-overview.md](./06-cli-and-sdk/cli-overview.md) | T1-CLI |
| Developer commands | [06-cli-and-sdk/developer-commands.md](./06-cli-and-sdk/developer-commands.md) | T1-DEV-CMDS |
| Non-interactive mode | [06-cli-and-sdk/non-interactive.md](./06-cli-and-sdk/non-interactive.md) | T1-NONINTERACTIVE |
| Codex SDK | [06-cli-and-sdk/codex-sdk.md](./06-cli-and-sdk/codex-sdk.md) | T2-SDK |
| Codex GitHub Action | [06-cli-and-sdk/github-action.md](./06-cli-and-sdk/github-action.md) | T2-GHA |
| Codex cloud | [07-cloud-and-environments/cloud.md](./07-cloud-and-environments/cloud.md) | T1-CLOUD, T1-ENV-CLOUD |
| Agent internet access | [07-cloud-and-environments/internet-access.md](./07-cloud-and-environments/internet-access.md) | T1-CLOUD-NET |
| Codex environments | [07-cloud-and-environments/environment-modes.md](./07-cloud-and-environments/environment-modes.md) | T1-ENV-MODES |
| Local environments | [07-cloud-and-environments/local-environment.md](./07-cloud-and-environments/local-environment.md) | T1-ENV-LOCAL |
| Worktrees | [07-cloud-and-environments/worktrees.md](./07-cloud-and-environments/worktrees.md) | T1-WORKTREES |
| Managed configuration | [08-enterprise/managed-configuration.md](./08-enterprise/managed-configuration.md) | T1-MANAGED |
| Admin rollout / team config | [08-enterprise/admin-setup-team-config.md](./08-enterprise/admin-setup-team-config.md) | T1-ADMIN-SETUP |
| Skill controls | [08-enterprise/skill-controls.md](./08-enterprise/skill-controls.md) | T1-ENT-SKILLS |
| Plugin controls | [08-enterprise/plugin-controls.md](./08-enterprise/plugin-controls.md) | T1-ENT-PLUGINS |
| Security administration | [08-enterprise/security-administration.md](./08-enterprise/security-administration.md) | T1-SEC-ADMIN |

## Meta

- Source registry: [SOURCES.md](./SOURCES.md)
- Topic file pattern: [TEMPLATE.md](./TEMPLATE.md)
- Schemas: [schemas/config-schema.json](./schemas/config-schema.json), [schemas/hooks/](./schemas/hooks/) (23 generated JSON files)

## Invariants

1. YAML front matter on every topic `.md` file.
2. Verbatim excerpts only (minimal glue).
3. No project-specific mapping sections in topic files.
4. English only; canonical URLs from [learn.chatgpt.com/llms.txt](https://learn.chatgpt.com/llms.txt).
5. Update [SOURCES.md](./SOURCES.md) when re-fetching after doc changes.

## agentd relevance

Reference only — not duplicated inside topic files.

| Codex surface | agentd touchpoint |
|---------------|-------------------|
| `~/.codex/hooks.json` + inline `[hooks]` in `config.toml` | `agentd install --provider=codex --scope=user` → `$CODEX_HOME` or `~/.codex` |
| `<repo>/.codex/hooks.json` | `--scope=project` → `cwd/.codex` |
| Lifecycle command hooks (stdin JSON → stdout JSON; empty stdout + exit 0 = continue) | `agentd hook run --provider=codex` |
| `notify = ["…"]` argv JSON (`type: agent-turn-complete`, `thread-id`, …) | `agentd hook notify` — **separate wire**; async-only, never a gate |
| Hook trust review of exact definition (hash); `/hooks`; managed hooks | Install writes trust state; moving `CODEX_HOME` requires reinstall |
| Project trust via `projects.<path>.trust_level` | Untrusted skips project hooks/config/rules |
| `allow_managed_hooks_only` + `[hooks]` in `requirements.toml` | Managed hooks bypass user trust review |
| `PreToolUse` / `PermissionRequest` decision fields | No CapAsk; `ask` degrades via `policy.ask_fallback` |
| MCP in `[mcp_servers.<name>]` + MCP **tool hooks** | Guard `mcp`; second hook handler type |
| `CODEX_HOME/sessions/**/rollout-*-{id}.jsonl` | Trajectory L2 import (see gaps) |
| `features.hooks` | Codex can disable lifecycle hooks entirely |
| Plugin-bundled hooks | Future install/scope concerns |

## Gaps vs agentd assumptions

Facts recorded when official docs omit or conflict with known layouts:

1. **Rollout session layout** — Official docs emphasize `history.jsonl` / session logging under `CODEX_HOME`. Nested `sessions/YYYY/MM/DD/rollout-*-{session_id}.jsonl` (used by trajectory import) was **not** found in the config-advanced / history sources reviewed. See [sessions-and-history.md](./05-settings-and-paths/sessions-and-history.md) research note.
2. **Hook trust storage file/key** — Docs say Codex records trust against the hook's current hash and that trust is persisted, but they do **not** name the on-disk file or `config.toml` key. agentd install assumes trust keys in `config.toml`; treat as **undocumented**.
3. **Cloud vs local hooks** — Docs state Local/Worktree run on the user's computer; they do **not** state whether Codex cloud loads user-level `~/.codex/hooks.json`. See [cloud-hooks-matrix.md](./02-hooks/cloud-hooks-matrix.md).
4. **Skills locations** — Skills are documented under `.agents/skills` / `~/.agents/skills` / `/etc/codex/skills`, while `CODEX_HOME` env docs also say the home root includes “skills”. Capture both claims; do not assume `~/.codex/skills` is the primary user skills path.
5. **`notify` vs lifecycle hooks** — `notify` is a user-config argv program (currently `agent-turn-complete` only), not an event in `hooks.json`. Do not conflate with `Stop` / `SessionEnd` hooks.

## Out of scope

Billing/pricing, pets, visualizations, voice, images, ads, commerce, Codex Security plugin deep-dive walkthroughs, ChatGPT web-only marketing UI, dumping `llms-full.txt` as one blob.
