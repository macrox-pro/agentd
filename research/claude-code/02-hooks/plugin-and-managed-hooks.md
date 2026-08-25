---
primary_sources:
  - id: T1-PLUGINS-REF
    title: "Plugins reference"
    url: "https://code.claude.com/docs/en/plugins-reference.md"
    section: "Hooks"
  - id: T1-MANAGED
    title: "Managed settings"
    url: "https://code.claude.com/docs/en/managed-settings.md"
    section: "Deploy hooks via managed settings"
also_cited_in:
  - id: T1-SETTINGS-REF
    title: "Settings reference"
    url: "https://code.claude.com/docs/en/settings-reference.md"
    section: "allowManagedHooksOnly"
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Plugin and managed hooks

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Plugins reference — Hooks

> ### Hooks
>
> Plugins can provide event handlers that respond to Claude Code events automatically.
>
> **Location**: `hooks/hooks.json` in plugin root, or inline in plugin.json
>
> **Format**: JSON configuration with event matchers and actions
>
> **Hook configuration**:
>
> ```json
> {
>   "hooks": {
>     "PostToolUse": [
>       {
>         "matcher": "Write|Edit",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "\"${CLAUDE_PLUGIN_ROOT}\"/scripts/format-code.sh"
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> Plugin hooks respond to the same lifecycle events as [user-defined hooks](/docs/en/hooks):
>
> | Event                 | When it fires                                                                                                                                                                                                                                         |
> | :-------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `SessionStart`        | When a session begins or resumes                                                                                                                                                                                                                      |
> | `Setup`               | When you start Claude Code with `--init-only`, or with `--init` or `--maintenance` in `-p` mode. For one-time preparation in CI or scripts                                                                                                            |
> | `UserPromptSubmit`    | When you submit a prompt, before Claude processes it                                                                                                                                                                                                  |
> | `UserPromptExpansion` | When a user-typed command expands into a prompt, before it reaches Claude. Can block the expansion                                                                                                                                                    |
> | `PreToolUse`          | Before a tool call executes. Can block it                                                                                                                                                                                                             |
> | `PermissionRequest`   | When a tool call needs a permission decision                                                                                                                                                                                                          |
> | `PermissionDenied`    | When auto mode denies a tool call, including denials without a classifier verdict. Use JSON `hookSpecificOutput.retry: true` to tell the model it may retry the denied tool call. Claude Code ignores `retry` when the classifier produced no verdict |
> | `PostToolUse`         | After a tool call succeeds                                                                                                                                                                                                                            |
> | `PostToolUseFailure`  | After a tool call fails                                                                                                                                                                                                                               |
> | `PostToolBatch`       | After a full batch of parallel tool calls resolves, before the next model call                                                                                                                                                                        |
> | `Notification`        | When Claude Code sends a notification                                                                                                                                                                                                                 |
> | `MessageDisplay`      | While assistant message text is displayed                                                                                                                                                                                                             |
> | `SubagentStart`       | When a subagent is spawned                                                                                                                                                                                                                            |
> | `SubagentStop`        | When a subagent finishes                                                                                                                                                                                                                              |
> | `TaskCreated`         | When a task is being created via `TaskCreate`                                                                                                                                                                                                         |
> | `TaskCompleted`       | When a task is being marked as completed                                                                                                                                                                                                              |
> | `Stop`                | When Claude finishes responding                                                                                                                                                                                                                       |
> | `StopFailure`         | When the turn ends due to an API error                                                                                                                                                                                                                |
> | `TeammateIdle`        | When an [agent team](/docs/en/agent-teams) teammate is about to go idle                                                                                                                                                                                    |
> | `InstructionsLoaded`  | When a CLAUDE.md or `.claude/rules/*.md` file is loaded into context. Fires at session start and when files are lazily loaded during a session                                                                                                        |
> | `ConfigChange`        | When a configuration file changes during a session                                                                                                                                                                                                    |
> | `CwdChanged`          | When the working directory changes, for example when Claude executes a `cd` command. Useful for reactive environment management with tools like direnv                                                                                                |
> | `DirectoryAdded`      | When a working directory is added mid-session via `/add-dir` or the SDK `register_repo_root` control request                                                                                                                                          |
> | `FileChanged`         | When a watched file changes on disk. The `matcher` field specifies which filenames to watch                                                                                                                                                           |
> | `WorktreeCreate`      | When a worktree is being created via `--worktree`, `isolation: "worktree"`, or for a background session. Replaces default git behavior                                                                                                                |
> | `WorktreeRemove`      | When a worktree is being removed at session exit, when a subagent finishes, or when you delete a background session                                                                                                                                   |
> | `PreCompact`          | Before context compaction                                                                                                                                                                                                                             |
> | `PostCompact`         | After context compaction completes                                                                                                                                                                                                                    |
> | `Elicitation`         | When an MCP server requests user input during a tool call                                                                                                                                                                                             |
> | `ElicitationResult`   | After a user responds to an MCP elicitation, before the response is sent back to the server                                                                                                                                                           |
> | `SessionEnd`          | When a session terminates                                                                                                                                                                                                                             |
>
> **Hook types**:
>
> * `command`: execute shell commands or scripts
> * `http`: send the event JSON as a POST request to a URL
> * `mcp_tool`: call a tool on a configured [MCP server](/docs/en/mcp)
> * `prompt`: evaluate a prompt with an LLM (uses `$ARGUMENTS` placeholder for context)
> * `agent`: run an agentic verifier with tools for complex verification tasks
>
> Hooks that target the plugin's own [bundled MCP server](#mcp-servers) must use its scoped names. Tool matchers and `if` fields take the scoped tool name `mcp__plugin_<plugin-name>_<server-name>__<tool>`, and an `mcp_tool` hook's `server` field takes `plugin:<plugin-name>:<server-name>`. A matcher written against the bare server key never fires. See [Match MCP tools](/docs/en/hooks#match-mcp-tools) and [Plugin-provided MCP servers](/docs/en/mcp#plugin-provided-mcp-servers).

### Source: Managed settings — Deploy hooks via managed settings

> Hooks use the same format as in `settings.json`.
>
> This example runs an audit script after every file edit across the organization:
>
> ```json
> {
>   "hooks": {
>     "PostToolUse": [
>       {
>         "matcher": "Edit|Write",
>         "hooks": [
>           { "type": "command", "command": "/usr/local/bin/audit-edit.sh" }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> Because hooks execute shell commands, users in interactive sessions see a [security approval dialog](#security-approval-dialogs) before Claude Code applies them.

### Source: Settings reference — allowManagedHooksOnly

> ### `allowManagedHooksOnly`
>
> Restrict hook execution to hooks your organization deploys.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean
>   * `true`: only managed hooks run, plus Agent SDK hooks and hooks from plugins your managed settings force-enable. See [What runs under `allowManagedHooksOnly`](#what-runs-under-allowmanagedhooksonly)
>   * `false`: hooks from every settings file and plugin run
> * **Default**: unset, so hooks from every settings file and plugin run
>
> ```json managed-settings.json
> {
>   "allowManagedHooksOnly": true
> }
> ```
>
> #### What runs under `allowManagedHooksOnly`
>
> When you set it to `true`, Claude Code changes which hooks and hook-like commands load:
>
> * **Managed and SDK hooks run**: hooks from managed settings and hooks the [Agent SDK](/docs/en/agent-sdk/overview) registers in process
> * **Force-enabled plugin hooks run**: hooks from plugins your managed settings force-enable through [`enabledPlugins`](#enabledplugins). Claude Code matches on the full `plugin@marketplace` ID, so a plugin with the same name from a different marketplace stays blocked. This lets you distribute vetted hooks through an organization marketplace while blocking everything else
> * **Everything else is blocked**: user, project, and local hooks, hooks from other plugins, and hooks declared in agent frontmatter
> * **Command-sourced plugins are disabled**: Claude Code also disables plugins with a [`command` source](/docs/en/plugin-marketplaces#command-sources), including plugins force-enabled in managed `enabledPlugins`, unless you set [`disableCommandPluginSources`](#disablecommandpluginsources) to `false` explicitly
> * **Marketplace `headersHelper` commands are blocked**: Claude Code also blocks marketplace [`headersHelper` commands](/docs/en/plugin-marketplaces#authenticate-archive-downloads) unless [`disableCommandPluginSources`](#disablecommandpluginsources) is explicitly set to `false`, except for a marketplace that managed settings themselves declare. Requires Claude Code v2.1.238 or later
> * **Status line and file suggestion narrow to managed settings**: Claude Code reads [`statusLine`](/docs/en/statusline), [`fileSuggestion`](#filesuggestion), and [`subagentStatusLine`](/docs/en/statusline#subagent-status-lines) from managed settings only, following the [status line and file suggestion gates](#status-line-and-file-suggestion-gates)
>
> The [`/goal`](/docs/en/goal) command can't run while this key is set, because it depends on hooks.
