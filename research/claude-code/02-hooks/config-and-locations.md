---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Configuration; Hook locations"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook configuration and locations

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — Configuration

> ## Configuration
>
> Hooks are defined in JSON settings files. The configuration has three levels of nesting:
>
> 1. Choose a [hook event](#hook-events) to respond to, like `PreToolUse` or `Stop`
> 2. Add a [matcher group](#matcher-patterns) to filter when it fires, like "only for the Bash tool"
> 3. Define one or more [hook handlers](#hook-handler-fields) to run when matched
>
> See [How a hook resolves](#how-a-hook-resolves) above for a complete walkthrough with an annotated example.
>
>
>   This page uses specific terms for each level: **hook event** for the lifecycle point, **matcher group** for the filter, and **hook handler** for the shell command, HTTP endpoint, MCP tool, prompt, or agent that runs. "Hook" on its own refers to the general feature.
>
>
> ### Hook locations
>
> Where you define a hook determines its scope:
>
> | Location                                 | Scope                                                                                                            | Shareable                                             |
> | :--------------------------------------- | :--------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------- |
> | `~/.claude/settings.json`                | All your projects                                                                                                | No, local to your machine                             |
> | `.claude/settings.json`                  | Single project                                                                                                   | Yes, can be committed to the repo                     |
> | `.claude/settings.local.json`            | Single project                                                                                                   | No, gitignored when Claude Code saves a setting to it |
> | Managed policy settings                  | Organization-wide                                                                                                | Yes, admin-controlled                                 |
> | [Plugin](/docs/en/plugins) `hooks/hooks.json` | When plugin is enabled                                                                                           | Yes, bundled with the plugin                          |
> | [Skill](/docs/en/skills) frontmatter          | The rest of the session once the skill is invoked. See [Hooks in skills and agents](#hooks-in-skills-and-agents) | Yes, defined in the skill file                        |
> | [Subagent](/docs/en/sub-agents) frontmatter   | While that subagent is running                                                                                   | Yes, defined in the subagent file                     |
>
> Cloud sessions on [Claude Code on the web](/docs/en/claude-code-on-the-web) don't read your local `~/.claude/settings.json`; hooks there come from the repo and from your organization's server-managed settings. In a [self-hosted environment](/docs/en/self-hosted-environments-configuration#permissions-and-tool-approval), Claude Code also runs the hooks the operator seeded from the runner host's `~/.claude/`, and it runs the hooks in the runner image's managed settings file when neither [server-managed settings nor an MDM-delivered Claude Code policy](/docs/en/settings#precedence-within-the-managed-tier) supplies the managed tier. See [what carries over from your setup](/docs/en/cloud-environments#what-carries-over-from-your-setup) for which files reach a cloud session.
>
> For details on settings file resolution, see [settings](/docs/en/settings).
>
> Hooks from settings files, managed policy settings, and plugins also run inside [subagents](/docs/en/sub-agents). When a subagent calls a tool, tool events such as `PreToolUse` and `PostToolUse` fire the same configured hooks as in the main conversation, and the input carries the `agent_id` and `agent_type` [common input fields](#common-input-fields) that identify the subagent.
>
> Enterprise administrators can use `allowManagedHooksOnly` to restrict which hooks run:
>
> * Your user, project, local, and plugin hooks are blocked. Hooks from plugins force-enabled in managed settings `enabledPlugins` are exempt
> * Claude Code also narrows your [`statusLine`](/docs/en/statusline), [`fileSuggestion`](/docs/en/settings-reference#filesuggestion), and [`subagentStatusLine`](/docs/en/statusline#subagent-status-lines) settings to managed settings
> * Claude Code also disables plugins with a [`command` source](/docs/en/plugin-marketplaces#command-sources), including plugins force-enabled in managed settings `enabledPlugins`, unless [`disableCommandPluginSources`](/docs/en/settings-reference#disablecommandpluginsources) is explicitly set to `false`
> * Claude Code also blocks marketplace [`headersHelper` commands](/docs/en/plugin-marketplaces#authenticate-archive-downloads) unless [`disableCommandPluginSources`](/docs/en/settings-reference#disablecommandpluginsources) is explicitly set to `false`, except for a marketplace that managed settings themselves declare
>
> See [what runs under `allowManagedHooksOnly`](/docs/en/settings-reference#what-runs-under-allowmanagedhooksonly).
>
> Hook entries merge across settings levels rather than replacing each other: user, project, and local settings add their own hooks without removing managed ones, and the [`disableAllHooks`](#disable-or-remove-hooks) setting can't disable managed hooks from outside managed settings.
>
> The [HTTP hook allowlists](/docs/en/settings-reference#hook-and-skill-settings) apply to hooks from every source, including managed policy settings:
>
> * `allowedHttpHookUrls`: when defined at any settings level, Claude Code runs an HTTP hook handler only if its URL matches the merged allowlist
> * `httpHookAllowedEnvVars`: when defined, Claude Code interpolates only the environment variables on that list into hook headers

### Source: Hooks reference — Hook locations

> ### Hook locations
>
> Where you define a hook determines its scope:
>
> | Location                                 | Scope                                                                                                            | Shareable                                             |
> | :--------------------------------------- | :--------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------- |
> | `~/.claude/settings.json`                | All your projects                                                                                                | No, local to your machine                             |
> | `.claude/settings.json`                  | Single project                                                                                                   | Yes, can be committed to the repo                     |
> | `.claude/settings.local.json`            | Single project                                                                                                   | No, gitignored when Claude Code saves a setting to it |
> | Managed policy settings                  | Organization-wide                                                                                                | Yes, admin-controlled                                 |
> | [Plugin](/docs/en/plugins) `hooks/hooks.json` | When plugin is enabled                                                                                           | Yes, bundled with the plugin                          |
> | [Skill](/docs/en/skills) frontmatter          | The rest of the session once the skill is invoked. See [Hooks in skills and agents](#hooks-in-skills-and-agents) | Yes, defined in the skill file                        |
> | [Subagent](/docs/en/sub-agents) frontmatter   | While that subagent is running                                                                                   | Yes, defined in the subagent file                     |
>
> Cloud sessions on [Claude Code on the web](/docs/en/claude-code-on-the-web) don't read your local `~/.claude/settings.json`; hooks there come from the repo and from your organization's server-managed settings. In a [self-hosted environment](/docs/en/self-hosted-environments-configuration#permissions-and-tool-approval), Claude Code also runs the hooks the operator seeded from the runner host's `~/.claude/`, and it runs the hooks in the runner image's managed settings file when neither [server-managed settings nor an MDM-delivered Claude Code policy](/docs/en/settings#precedence-within-the-managed-tier) supplies the managed tier. See [what carries over from your setup](/docs/en/cloud-environments#what-carries-over-from-your-setup) for which files reach a cloud session.
>
> For details on settings file resolution, see [settings](/docs/en/settings).
>
> Hooks from settings files, managed policy settings, and plugins also run inside [subagents](/docs/en/sub-agents). When a subagent calls a tool, tool events such as `PreToolUse` and `PostToolUse` fire the same configured hooks as in the main conversation, and the input carries the `agent_id` and `agent_type` [common input fields](#common-input-fields) that identify the subagent.
>
> Enterprise administrators can use `allowManagedHooksOnly` to restrict which hooks run:
>
> * Your user, project, local, and plugin hooks are blocked. Hooks from plugins force-enabled in managed settings `enabledPlugins` are exempt
> * Claude Code also narrows your [`statusLine`](/docs/en/statusline), [`fileSuggestion`](/docs/en/settings-reference#filesuggestion), and [`subagentStatusLine`](/docs/en/statusline#subagent-status-lines) settings to managed settings
> * Claude Code also disables plugins with a [`command` source](/docs/en/plugin-marketplaces#command-sources), including plugins force-enabled in managed settings `enabledPlugins`, unless [`disableCommandPluginSources`](/docs/en/settings-reference#disablecommandpluginsources) is explicitly set to `false`
> * Claude Code also blocks marketplace [`headersHelper` commands](/docs/en/plugin-marketplaces#authenticate-archive-downloads) unless [`disableCommandPluginSources`](/docs/en/settings-reference#disablecommandpluginsources) is explicitly set to `false`, except for a marketplace that managed settings themselves declare
>
> See [what runs under `allowManagedHooksOnly`](/docs/en/settings-reference#what-runs-under-allowmanagedhooksonly).
>
> Hook entries merge across settings levels rather than replacing each other: user, project, and local settings add their own hooks without removing managed ones, and the [`disableAllHooks`](#disable-or-remove-hooks) setting can't disable managed hooks from outside managed settings.
>
> The [HTTP hook allowlists](/docs/en/settings-reference#hook-and-skill-settings) apply to hooks from every source, including managed policy settings:
>
> * `allowedHttpHookUrls`: when defined at any settings level, Claude Code runs an HTTP hook handler only if its URL matches the merged allowlist
> * `httpHookAllowedEnvVars`: when defined, Claude Code interpolates only the environment variables on that list into hook headers
