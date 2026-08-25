---
primary_sources:
  - id: T1-SETTINGS-REF
    title: "Settings reference"
    url: "https://code.claude.com/docs/en/settings-reference.md"
    section: "Hooks and automation; Permission settings; env"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Settings reference — hooks and permissions

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Settings reference — hooks, permissions, env

> ## Hooks and automation
>
> Register hooks, restrict which hooks run, and control workflows. For hook events and payloads, see the [hooks reference](/docs/en/hooks).
>
> ### `allowedHttpHookUrls`
>
> Limit which URLs [HTTP hooks](/docs/en/hooks#http-hook-fields) can target. When you define this key, Claude Code runs an HTTP hook only if its URL matches one of the patterns and blocks the rest without running them; an empty array blocks every HTTP hook.
>
> * **Scope**: [`Any file`](#scopes). Arrays merge across settings files.
> * **Type**: array of URL patterns, with `*` as a wildcard
> * **Default**: unset, so any URL is allowed
>
> This example allows any URL under `https://hooks.example.com/` and any `http://localhost` URL:
>
> ```json settings.json
> {
>   "allowedHttpHookUrls": ["https://hooks.example.com/*", "http://localhost:*"]
> }
> ```
>
> Hostname matching is case-insensitive and treats `hooks.example.com.`, with the trailing dot that marks a fully qualified domain name, the same as `hooks.example.com`, which is how DNS treats them. The allowlist applies to hooks from every source, including managed settings.
>
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
>
> ### `disableAllHooks`
>
> Turn off [hooks](/docs/en/hooks#disable-or-remove-hooks), any custom [status line](/docs/en/statusline), and any custom [file suggestion](#filesuggestion) command. Use it to turn all of these off temporarily without deleting them from your settings.
>
> * **Scope**: [`Any file`](#scopes). Only managed settings can disable managed hooks.
> * **Type**: Boolean
>   * `true`: Claude Code turns off hooks, any custom status line, and any custom file suggestion command
>   * `false`: hooks, the status line, and the file suggestion command run
> * **Default**: unset, so hooks run
>
> ```json settings.json
> {
>   "disableAllHooks": true
> }
> ```
>
> The reach depends on which file carries the key:
>
> * **In managed settings**: Claude Code disables every configured hook, including managed ones, and keeps running the hooks the [Agent SDK](/docs/en/agent-sdk/overview) registers in process
> * **In any other settings file**: Claude Code disables user, project, local, and plugin hooks; managed hooks, Agent SDK hooks, and hooks from plugins force-enabled in managed [`enabledPlugins`](#enabledplugins) keep running
>
> Keeping Agent SDK hooks running when managed settings set this key requires Claude Code v2.1.242 or later.
>
> The [`/goal`](/docs/en/goal) command can't run while hooks are disabled, and the `/hooks` menu shows a notice instead of your hooks.
>
> #### Status line and file suggestion gates
>
> Claude Code makes two decisions for `statusLine`, `fileSuggestion`, and `subagentStatusLine`, in this order:
>
> * **Off entirely**: when managed settings set `disableAllHooks`, or when the folder isn't trusted under the same [workspace trust rule as hooks in settings files](/docs/en/permissions#what-runs-before-you-trust-a-folder)
> * **Narrowed to managed settings**: when [`allowManagedHooksOnly`](#allowmanagedhooksonly) is set, when `disableAllHooks` is `true` outside managed settings after [settings precedence](/docs/en/hooks#disable-or-remove-hooks) applies, or when you start Claude Code with `--safe-mode`
>
> Under narrowing, Claude Code runs a managed value if one is deployed. Otherwise it skips your value without warning: the status line is disabled, and `@` autocomplete falls back to the built-in file suggestion.
>
> ### `disableWorkflows`
>
> Turn off [dynamic workflows](/docs/en/workflows#turn-workflows-off) and the bundled workflow commands for everyone your settings reach, such as an organization through managed settings. To turn workflows on or off just for yourself, use [`enableWorkflows`](#enableworkflows) instead, which the **Dynamic workflows** toggle in `/config` writes to your user settings.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code turns off dynamic workflows and the bundled workflow commands for everyone your settings reach
>   * `false`: the same as unset; whether workflows are on then follows [`enableWorkflows`](#enableworkflows) and your plan's default
> * **Default**: `false`
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_WORKFLOWS`](/docs/en/env-vars) turns workflows off for one session; whichever of the two turns them off, the other can't turn them back on
>
> ```json settings.json
> {
>   "disableWorkflows": true
> }
> ```
>
> ### `enableWorkflows`
>
> Turn [dynamic workflows](/docs/en/workflows) on or off for yourself when your plan's default isn't what you want. Appears in `/config` as **Dynamic workflows**, which writes this key to your user settings and removes it again when you toggle back to your plan's default. To turn workflows off for everyone from managed settings, use [`disableWorkflows`](#disableworkflows) instead.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code turns dynamic workflows on for you
>   * `false`: Claude Code turns dynamic workflows off for you
> * **Default**: unset, so workflows are on unless you're on the Pro plan, where they're off
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_WORKFLOWS`](/docs/en/env-vars) turns workflows off for one session, and `true` here can't turn them back on while it's set
>
> ```json settings.json
> {
>   "enableWorkflows": true
> }
> ```
>
> [`disableWorkflows`](#disableworkflows) and your organization's workflows policy also take precedence: `enableWorkflows: true` can't turn workflows back on while any source turns workflows off. Claude Code hides the `/config` row while a source other than your user settings sets `enableWorkflows`, or sets `disableWorkflows` to `true`.
>
> ### `hooks`
>
> Run your own commands, prompts, agents, HTTP requests, or MCP tools as [hooks](/docs/en/hooks) at points in Claude Code's lifecycle, such as before a tool call or when a session starts; the [hooks reference](/docs/en/hooks#hook-events) lists every event, its payload, and its exit codes. Each event maps to a list of matcher groups, and each group lists the handlers to run when the matcher applies.
>
> * **Scope**: [`Any file`](#scopes). Hooks merge across files rather than replacing each other, and hooks from managed settings can't be removed from other files.
> * **Type**: object keyed by [hook event](/docs/en/hooks#hook-events); each value is an array of `{ "matcher", "hooks" }` groups whose `hooks` entries have a `type` of `"command"`, `"prompt"`, `"agent"`, `"http"`, or `"mcp_tool"`
> * **Default**: unset, so no hooks run
>
> This example runs a script before every Bash tool call:
>
> ```json settings.json
> {
>   "hooks": {
>     "PreToolUse": [
>       {
>         "matcher": "Bash",
>         "hooks": [
>           { "type": "command", "command": "~/.claude/hooks/check-bash.sh" }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> For every event, matcher pattern, and handler field, see the [hooks reference](/docs/en/hooks#configuration). To turn hooks off, see [`disableAllHooks`](#disableallhooks); to limit hooks to the ones your organization deploys, see [`allowManagedHooksOnly`](#allowmanagedhooksonly).
>
> ### `httpHookAllowedEnvVars`
>
> An [HTTP hook](/docs/en/hooks#http-hook-fields) can put the value of an environment variable into a request header, for example an `Authorization: Bearer $HOOK_TOKEN` header, but only for variables the hook lists in its own `allowedEnvVars`. This key sets an outer limit on that list for every HTTP hook: a hook can use a variable only if both its own `allowedEnvVars` and this key name it. Use it to stop a hook from reading a secret it shouldn't, even when the hook's definition asks for it.
>
> * **Scope**: [`Any file`](#scopes). Arrays merge across settings files.
> * **Type**: array of environment variable names
> * **Default**: unset, so each hook's own `allowedEnvVars` list applies
>
> This example limits header interpolation to `MY_TOKEN` and `HOOK_SECRET`:
>
> ```json settings.json
> {
>   "httpHookAllowedEnvVars": ["MY_TOKEN", "HOOK_SECRET"]
> }
> ```
>
> The allowlist applies to hooks from every source, including managed settings.
>
> ### `workflowKeywordTriggerEnabled`
>
> Choose whether typing the keyword `ultracode` in a prompt triggers a [dynamic workflow](/docs/en/workflows#ask-for-a-workflow-in-your-prompt). Set it to `false` to type the word without triggering one. Requires Claude Code v2.1.157 or later.
>
> * **Scope**: [`Any file`](#scopes). Appears in `/config` as **Ultracode keyword trigger**.
> * **Type**: Boolean
>   * `true`: typing `ultracode` in a prompt triggers a dynamic workflow
>   * `false`: you can type the word without triggering one
> * **Default**: `true`
>
> ```json settings.json
> {
>   "workflowKeywordTriggerEnabled": false
> }
> ```
>
> The `ultracode` effort setting, `/workflows`, and saved workflow commands are unaffected. Requires Claude Code v2.1.157 or later. Before v2.1.160, the trigger keyword was `workflow`.
>
> ### `workflowSizeGuideline`
>
> Set the [agent count Claude aims for](/docs/en/workflows#set-a-size-guideline) in the dynamic workflows it writes. Claude Code sends the value to Claude as advice, not an enforced cap: `"small"` asks for fewer than 5 agents, `"medium"` fewer than 15, and `"large"` fewer than 50. Choose `"small"` when you want to bound what a workflow spends. Requires Claude Code v2.1.219 or later.
>
> * **Scope**: [`Any file`](#scopes). A value there takes precedence over the **Dynamic workflow size** choice in `/config`, which Claude Code stores in `~/.claude.json`, and Claude Code hides that row while a settings file sets the key.
> * **Type**: string, one of:
>   * `"unrestricted"`: no guideline, so Claude sizes the workflow to the task
>   * `"small"`: Claude aims for fewer than 5 agents
>   * `"medium"`: Claude aims for fewer than 15 agents
>   * `"large"`: Claude aims for fewer than 50 agents
> * **Default**: `"medium"`
>
> ```json settings.json
> {
>   "workflowSizeGuideline": "small"
> }
> ```
>
> Requires Claude Code v2.1.219 or later; on v2.1.202 through v2.1.218, set the guideline in `/config` instead.
>
>
>
>
> ## Permission settings
>
> Decide what Claude can do without asking, which permission mode a session starts in, and what auto mode's classifier allows. For rule syntax and the permission model, see [Configure permissions](/docs/en/permissions).
>
> ### `allowManagedPermissionRulesOnly`
>
> Make managed settings the only source of `allow`, `ask`, and `deny` permission rules. Claude Code then ignores rules in user, project, local, and `--settings` files, ignores `--allowedTools`, hides the always-allow choices in permission prompts, and stops saving new rules.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean
>   * `true`: managed settings are the only source of `allow`, `ask`, and `deny` rules; Claude Code ignores rules from other files and `--allowedTools`, hides always-allow choices, and stops saving new rules
>   * `false`: Claude Code applies permission rules from user, project, local, and `--settings` files in addition to the managed ones
> * **Default**: unset, so Claude Code applies permission rules from user, project, and local settings and from `--settings`, in addition to the managed ones
>
> ```json managed-settings.json
> {
>   "allowManagedPermissionRulesOnly": true
> }
> ```
>
> This key doesn't lock down the MCP server allowlist; for that, set [`allowManagedMcpServersOnly`](#allowmanagedmcpserversonly). See [Managed-only settings](/docs/en/managed-settings#managed-only-settings).
>
> ### `autoMode`
>
> Add your own rules to what the [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) classifier blocks and allows. Use it to tell the classifier which repos, buckets, and domains your organization trusts, so it stops blocking routine internal operations. The classifier ships with [built-in allow and deny rules](/docs/en/auto-mode-config#inspect-the-defaults-and-your-effective-config). Include the literal string `"$defaults"` in an array to keep those built-in rules at that position and add yours around them; leave it out to replace them with yours.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: object with `environment`, `allow`, `soft_deny`, and `hard_deny` arrays of prose rules, plus the [`classifyAllShell`](#automode-classifyallshell) Boolean
> * **Default**: unset, so the classifier uses only its [built-in rules](/docs/en/auto-mode-config#inspect-the-defaults-and-your-effective-config)
>
> This example keeps the built-in `soft_deny` rules, through `"$defaults"`, and adds one more that blocks `terraform apply`:
>
> ```json settings.json
> {
>   "autoMode": {
>     "soft_deny": ["$defaults", "Never run terraform apply"]
>   }
> }
> ```
>
> When more than one of those files sets the same array, Claude Code concatenates the entries. For the rule format and how each array is applied, see [Configure auto mode](/docs/en/auto-mode-config).
>
> ### `autoMode.classifyAllShell`
>
> Send every Bash and PowerShell command through the auto mode classifier while auto mode is active. By default, auto mode suspends only allow rules that could run arbitrary code: tool-wide and wildcard rules such as `Bash(*)`, and interpreter or shell-wrapper prefixes such as `Bash(python *)`. A command that any other allow rule matches, such as `Bash(npm test)`, skips the classifier, and a destructive argument the rule's prefix didn't anticipate can get through unseen. Setting this key suspends every shell allow rule for the session so the classifier sees every command. Requires Claude Code v2.1.193 or later.
>
> * **Scope**: [`User or managed`](#scopes). Read wherever [`autoMode`](#automode) is read.
> * **Type**: Boolean
>   * `true`: while auto mode is active, Claude Code sends every Bash and PowerShell command through the classifier and suspends your shell allow rules; outside auto mode the rules still apply
>   * `false`: auto mode suspends only allow rules that could run arbitrary code, such as `Bash(*)` and `Bash(python *)`; a command that any other allow rule matches skips the classifier, and every other shell command goes through it
> * **Default**: `false`
>
> ```json settings.json
> {
>   "autoMode": {
>     "classifyAllShell": true
>   }
> }
> ```
>
> See [Route all shell commands through the classifier](/docs/en/auto-mode-config#route-all-shell-commands-through-the-classifier). Requires Claude Code v2.1.193 or later.
>
> ### `disableAutoMode`
>
> Remove [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) from the `Shift+Tab` cycle. Any session that would otherwise [start in auto mode](/docs/en/permission-modes#which-mode-a-session-starts-in), whether from `--permission-mode auto`, a settings file, or the built-in default, starts in `default` instead. Administrators set it in managed settings to prevent developers in their organization from using auto mode.
>
> * **Scope**: [`Any file`](#scopes). Most useful in [managed settings](/docs/en/managed-settings), where users can't override it. Also accepted under `permissions` as `permissions.disableAutoMode`.
> * **Type**: the string `"disable"`
> * **Default**: unset
>
> ```json settings.json
> {
>   "disableAutoMode": "disable"
> }
> ```
>
> ### `permissions`
>
> Control which tools Claude can use without asking, which ones always prompt, and which ones are blocked, and set the [permission mode](/docs/en/permission-modes) a session starts in. Every `permissions.*` key below nests under this object.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object with `allow`, `ask`, `deny`, `additionalDirectories`, `defaultMode`, `disableBypassPermissionsMode`, and `disableAutoMode`
> * **Default**: unset
>
> This example approves `npm run` commands without asking, prompts before `git push`, blocks reads of `.env`, and starts sessions in `acceptEdits`:
>
> ```json settings.json
> {
>   "permissions": {
>     "allow": ["Bash(npm run *)"],
>     "ask": ["Bash(git push *)"],
>     "deny": ["Read(./.env)"],
>     "defaultMode": "acceptEdits"
>   }
> }
> ```
>
> The three rule arrays share one syntax; see [Permission rule syntax](#permission-rule-syntax) under `permissions.allow`. For how permission rules from different files combine, see [how permission rules merge across scopes](/docs/en/permissions#settings-precedence); for how settings keys in general combine, see [Settings precedence](/docs/en/settings#settings-precedence) on the settings guide.
>
> ### `useAutoModeDuringPlan`
>
> Choose whether Claude Code uses the auto mode classifier to review shell commands in plan mode. With the default `true`, the classifier reviews each command during planning when auto mode is available and you see no prompt. Set `false` to get a permission prompt for every command outside the built-in read-only set. Appears in `/config` as **Use auto mode during plan**.
>
> * **Scope**: [`User, local, or managed`](#scopes). A repository can't turn it off for you.
> * **Type**: Boolean
>   * `true`: the same as unset; when auto mode is available, the classifier reviews each shell command during planning instead of prompting you for it. A `false` in any of these files still turns it off
>   * `false`: you get a permission prompt for every command outside the built-in read-only set
> * **Default**: `true`
>
> ```json settings.json
> {
>   "useAutoModeDuringPlan": false
> }
> ```
>
> ### `permissions.allow`
>
> List the tool uses Claude Code approves without asking you. In an MCP rule, `*` can appear only in the tool name after the `mcp__<server>__` prefix, such as `mcp__github__get_*`; it can't appear in the server name.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of permission rule strings
> * **Default**: unset
> * **Per-session overrides**: `--allowedTools` adds allow rules for one session, and a deny rule from any settings file still blocks a tool it names
>
> This example approves `git diff` and lets Claude Code read your `.zshrc` without asking:
>
> ```json settings.json
> {
>   "permissions": {
>     "allow": ["Bash(git diff *)", "Read(~/.zshrc)"]
>   }
> }
> ```
>
> Claude Code applies `allow` rules from a project's `.claude/settings.json` only after you accept the [workspace trust dialog](/docs/en/permissions#project-allow-rules-and-workspace-trust) for that folder.
>
> #### Permission rule syntax
>
> Permission rules follow the format `Tool` or `Tool(specifier)`. Claude Code evaluates `deny` rules first, then `ask`, then `allow`, and the first match decides regardless of how specific each rule is; see the [permission rule evaluation order](/docs/en/permissions#manage-permissions).
>
> Each row shows one rule shape and what it matches.
>
> | Rule                           | What it matches                  |
> | :----------------------------- | :------------------------------- |
> | `Bash`                         | Every Bash command               |
> | `Bash(npm run *)`              | Commands starting with `npm run` |
> | `Read(./.env)`                 | Reads of the `.env` file         |
> | `WebFetch(domain:example.com)` | Fetch requests to example.com    |
>
> For the complete rule syntax, including wildcard behavior, tool-specific patterns for Read, Edit, WebFetch, MCP, and Agent rules, and the security limitations of Bash patterns, see [Permission rule syntax](/docs/en/permissions#permission-rule-syntax).
>
> ### `permissions.ask`
>
> List the tool uses that prompt you for confirmation even in a permission mode that would otherwise approve them, such as `acceptEdits` or `bypassPermissions`. In `dontAsk` mode Claude Code denies a matching tool use instead of prompting.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of permission rule strings
> * **Default**: unset
>
> ```json settings.json
> {
>   "permissions": {
>     "ask": ["Bash(git push *)"]
>   }
> }
> ```
>
>
> ### `permissions.deny`
>
> List the tool uses Claude Code blocks. Use it for files that hold API keys, secrets, or environment values: Claude Code excludes matching files from file discovery and search results, denies reads of them, and blocks the [Edit and Write tools](/docs/en/permissions#read-and-edit) on the matching paths. Read and Edit deny rules apply to Claude's built-in file tools and to file commands Claude Code recognizes in Bash, such as `cat`, `head`, `tail`, and `sed`; they don't apply to arbitrary subprocesses, so for OS-level enforcement [enable the sandbox](/docs/en/sandboxing).
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of permission rule strings
> * **Default**: unset
> * **Per-session overrides**: `--disallowedTools` adds deny rules for one session alongside this key
>
> This example denies reads of `.env` files, the `secrets` directory, and a credentials file, and blocks `curl` commands:
>
> ```json settings.json
> {
>   "permissions": {
>     "deny": [
>       "Read(./.env)",
>       "Read(./.env.*)",
>       "Read(./secrets/**)",
>       "Read(./config/credentials.json)",
>       "Bash(curl *)"
>     ]
>   }
> }
> ```
>
> Tool names accept glob patterns, so `"*"` denies every tool and `"mcp__*"` denies every MCP tool. Claude Code ignores a deny rule for the [`EndConversation`](/docs/en/tools-reference#endconversation-tool-behavior) tool as long as any other tool is still available to Claude. For what a `Bash` deny rule can and can't catch, see [Bash permission limitations](/docs/en/permissions#tool-specific-permission-rules). This key replaces the deprecated `ignorePatterns` configuration.
>
> ### `permissions.additionalDirectories`
>
> Give Claude file access to directories outside the one you started in, as additional [working directories](/docs/en/permissions#working-directories). Most `.claude/` configuration is [not discovered](/docs/en/permissions#additional-directories-grant-file-access-not-configuration) from these directories.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of directory paths
> * **Default**: unset
> * **Per-session overrides**: `--add-dir` and `/add-dir` add directories for one session alongside this key
>
> ```json settings.json
> {
>   "permissions": {
>     "additionalDirectories": ["../docs/"]
>   }
> }
> ```
>
> Like `allow` rules, entries in a project's `.claude/settings.json` take effect only after you accept the [workspace trust dialog](/docs/en/permissions#project-allow-rules-and-workspace-trust) for that folder.
>
> ### `permissions.defaultMode`
>
> Set the [permission mode](/docs/en/permission-modes) new sessions start in. When you leave it unset, sessions start in the [built-in default](/docs/en/permission-modes#which-mode-a-session-starts-in) for your plan and surface.
>
> * **Scope**: [`Any file`](#scopes). `auto` doesn't take effect from project or local settings, so set it in `~/.claude/settings.json` instead. Conversations the VS Code extension starts read only user, managed, and `--settings` values.
> * **Type**: string, one of:
>   * `"default"`: Claude Code runs only reads without asking
>   * `"acceptEdits"`: Claude Code also runs file edits and common filesystem commands such as `mkdir` and `mv` without asking
>   * `"plan"`: Claude Code reads and plans but blocks edits until you approve a plan
>   * `"auto"`: Claude Code runs everything, with background safety checks
>   * `"dontAsk"`: Claude Code runs only pre-approved tools and auto-denies every call that would otherwise prompt
>   * `"bypassPermissions"`: Claude Code runs everything without asking
>   * `"manual"`: an alias for `"default"`, in Claude Code v2.1.200 or later
> * **Default**: unset
> * **Per-session overrides**: `--permission-mode`, and its equivalent `--dangerously-skip-permissions` for `bypassPermissions`, take precedence over this key for one session
>
> ```json settings.json
> {
>   "permissions": {
>     "defaultMode": "acceptEdits"
>   }
> }
> ```
>
> Permission rules layer on top of every mode: `deny` rules block in every mode, including `bypassPermissions`. See [Permission modes](/docs/en/permission-modes). `manual` names the permission mode labeled Manual in the CLI and the VS Code extension; the alias requires Claude Code v2.1.200 or later. In Claude Code on the web, Claude Code honors only `acceptEdits`, `plan`, `default`, and `auto` from this key. For conversations the VS Code extension starts, see [which setting the extension reads for the starting permission mode](/docs/en/permission-modes#switch-permission-modes).
>
> ### `permissions.disableBypassPermissionsMode`
>
> Prevent anyone from entering `bypassPermissions` mode. Claude Code then rejects the `--dangerously-skip-permissions` flag, and ignores an [agent definition's](/docs/en/sub-agents#permission-modes) `permissionMode: bypassPermissions`, so the subagent runs with the parent session's permission mode.
>
> * **Scope**: [`Any file`](#scopes). Typically set in [managed settings](/docs/en/managed-settings) to enforce organizational policy.
> * **Type**: the string `"disable"`
> * **Default**: unset
> * **Per-session overrides**: this key takes precedence over `--dangerously-skip-permissions`, which Claude Code rejects while the key is set
>
> ```json settings.json
> {
>   "permissions": {
>     "disableBypassPermissionsMode": "disable"
>   }
> }
> ```
>
> Before v2.1.223, Claude Code applied the frontmatter permission mode even with bypass disabled.
>
> ### `skipAutoPermissionPrompt`
>
> Skip the one-time notice describing [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) that Claude Code shows when you first enter auto mode yourself, for example through your own settings or the mode selector, rather than when the built-in default starts a session in it. Claude Code shows that notice once and then records that it was shown, so this key only matters where the notice hasn't appeared yet.
>
> * **Scope**: [`User or managed`](#scopes). A repository can't set it for you.
> * **Type**: Boolean
>   * `true`: Claude Code skips the notice
>   * `false`: the same as unset; the notice appears once unless another of these files sets `true`
> * **Default**: unset, so the notice appears once
>
> ```json settings.json
> {
>   "skipAutoPermissionPrompt": true
> }
> ```
>
> ### `skipDangerousModePermissionPrompt`
>
> Skip the confirmation dialog Claude Code shows before a session enters `bypassPermissions` mode, whether from `--dangerously-skip-permissions` or from `defaultMode: "bypassPermissions"`. Claude Code writes `true` here in your user settings when you accept that dialog once.
>
> * **Scope**: [`User, local, or managed`](#scopes). An untrusted repository can't skip the dialog for you.
> * **Type**: Boolean
>   * `true`: Claude Code skips the confirmation dialog before a session enters `bypassPermissions` mode
>   * `false`: the same as unset; the dialog appears unless another of these files sets `true`
> * **Default**: unset, so the dialog appears
>
> ```json settings.json
> {
>   "skipDangerousModePermissionPrompt": true
> }
> ```
>
> ## Memory and context
>
> ### `env`
>
> Set environment variables for every session and for the subprocesses Claude Code starts from it. Any variable in the [environment variables reference](/docs/en/env-vars) can go here, which is how you apply one to every session or roll it out to your team.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object mapping variable names to string values
> * **Default**: unset
>
> This example turns off automatic compaction and routes API requests through a proxy:
>
> ```json settings.json
> {
>   "env": {
>     "DISABLE_AUTO_COMPACT": "1",
>     "ANTHROPIC_BASE_URL": "https://proxy.example.com"
>   }
> }
> ```
>
> #### How `env` values interact with your shell
>
> * A value here overwrites the same variable exported in your shell, and when more than one settings file sets a variable, the [highest-precedence](/docs/en/settings#settings-precedence) one applies.
> * To cancel a shell export, set the variable to `""`. Claude Code treats an empty value as unset for provider selection, and subprocesses inherit the empty value.
> * `NO_COLOR` and `FORCE_COLOR` set here reach only subprocesses. To change Claude Code's own interface colors, set them in your shell before launching `claude`.
> * Values here are plain text in the settings file and reach every subprocess Claude Code starts. For an OTLP bearer token that rotates, use [`otelHeadersHelper`](#otelheadershelper); for API credentials, use [`apiKeyHelper`](#apikeyhelper).
>
> #### When Claude Code applies `env` values
>
> * From user settings, `--settings`, and managed settings: at startup, and again in the running session when a saved change alters the merged `env`.
> * From project and local settings: after you trust the workspace, or at startup in `-p` mode, which never shows the trust dialog, and again when a saved change alters the merged `env`.
> * Variables Claude Code classifies as safe, such as model selection, timeouts and limits, feature toggles, and telemetry settings: at startup from every settings file.
>
> #### Variables Claude Code ignores in `env`
>
> * Project and local settings can't set a few variables, such as `CLAUDE_CODE_PROCESS_WRAPPER`, `CLAUDE_CODE_SYNC_SKILLS`, `CLAUDE_CODE_SYNC_PLUGINS`, `CLAUDE_CODE_PLUGIN_CACHE_DIR`, and `CLAUDE_CODE_PLUGIN_SEED_DIR`; set those in user or managed settings.
> * Identity variables that Claude Code's hosting environments own, such as `CLAUDE_CODE_REMOTE` and `CLAUDE_CODE_ACCOUNT_UUID`, are ignored from every file.
> * [`CLAUDE_CODE_MESSAGING_SOCKET` and `CLAUDE_CODE_MESSAGING_TOKEN`](/docs/en/env-vars#variables), which Claude Code exports itself, are ignored from every file. Ignoring the socket variable requires Claude Code v2.1.224 or later, and ignoring the token requires v2.1.228 or later.
> * [`CLAUDE_CODE_PROJECT_DIR_NAME`](/docs/en/sessions#name-the-project-directory-yourself), which Claude Code reads from the launch environment only, is ignored from every file; requires v2.1.234 or later.
