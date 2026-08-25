---
primary_sources:
  - id: T1-HOOKS
    title: "Configuration and Team Distribution"
    url: "https://cursor.com/docs/hooks.md"
    section: "Configuration and Team Distribution"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook configuration and precedence

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Cursor Hooks — Configuration

> ## Configuration
>
> Define hooks in a `hooks.json` file. Configuration can exist at multiple levels. All matching hooks from every source run; when responses conflict, higher-priority sources take precedence during merge:
>
> ```sh
> ~/.cursor/
> ├── hooks.json
> └── hooks/
>     ├── audit.sh
>     └── block-git.sh
> ```
>
> - **Enterprise** (MDM-managed, system-wide):
>   - macOS: `/Library/Application Support/Cursor/hooks.json`
>   - Linux/WSL: `/etc/cursor/hooks.json`
>   - Windows: `C:\\ProgramData\\Cursor\\hooks.json`
> - **Team** (Cloud-distributed, enterprise only):
>   - Configured in the [web dashboard](https://cursor.com/dashboard/team-content?section=hooks) and synced to all team members automatically
> - **Project** (Project-specific):
>   - `<project-root>/.cursor/hooks.json`
>   - Project hooks run in any trusted workspace and are checked into version control with your project
> - **User** (User-specific):
>   - `~/.cursor/hooks.json`
>
> Priority order (highest to lowest): Enterprise → Team → Project → User
>
> The `hooks` object maps hook names to arrays of hook definitions. Each definition currently supports a `command` property that can be a shell string, an absolute path, or a relative path. The working directory depends on the hook source:
>
> - **Project hooks** (`.cursor/hooks.json` in a repository): Run from the **project root**
> - **User hooks** (`~/.cursor/hooks.json`): Run from `~/.cursor/`
> - **Enterprise hooks** (system-wide config): Run from the enterprise config directory
> - **Team hooks** (cloud-distributed): Run from the managed hooks directory
>
> For project hooks, use paths like `.cursor/hooks/script.sh` (relative to project root), not `./hooks/script.sh` (which would look for `<project>/hooks/script.sh`).
>
> ### Configuration file
>
> This example shows a user-level hooks file (`~/.cursor/hooks.json`). For project-level hooks, change paths like `./hooks/script.sh` to `.cursor/hooks/script.sh`:
>
> ```json
> {
>   "version": 1,
>   "hooks": {
>     "sessionStart": [{ "command": "./session-init.sh" }],
>     "sessionEnd": [{ "command": "./audit.sh" }],
>     "preToolUse": [
>       {
>         "command": "./hooks/validate-tool.sh",
>         "matcher": "Shell|Read|Write"
>       }
>     ],
>     "postToolUse": [{ "command": "./hooks/audit-tool.sh" }],
>     "subagentStart": [{ "command": "./hooks/validate-subagent.sh" }],
>     "subagentStop": [{ "command": "./hooks/audit-subagent.sh" }],
>     "beforeShellExecution": [{ "command": "./script.sh" }],
>     "afterShellExecution": [{ "command": "./script.sh" }],
>     "afterMCPExecution": [{ "command": "./script.sh" }],
>     "afterFileEdit": [{ "command": "./format.sh" }],
>     "preCompact": [{ "command": "./audit.sh" }],
>     "stop": [{ "command": "./audit.sh", "loop_limit": 10 }],
>     "beforeTabFileRead": [{ "command": "./redact-secrets-tab.sh" }],
>     "afterTabFileEdit": [{ "command": "./format-tab.sh" }],
>     "workspaceOpen": [{ "command": "./register-workspace-plugins.sh" }]
>   }
> }
> ```
>
> The Agent hooks (`sessionStart`, `sessionEnd`, `preToolUse`, `postToolUse`, `postToolUseFailure`, `subagentStart`, `subagentStop`, `beforeShellExecution`, `afterShellExecution`, `beforeMCPExecution`, `afterMCPExecution`, `beforeReadFile`, `afterFileEdit`, `beforeSubmitPrompt`, `preCompact`, `stop`, `afterAgentResponse`, `afterAgentThought`) apply to Cmd+K and Agent Chat operations. The Tab hooks (`beforeTabFileRead`, `afterTabFileEdit`) apply specifically to inline Tab completions. The app lifecycle hook (`workspaceOpen`) fires when a workspace opens and on workspace folder changes, independent of any agent session.
>
> ### Global Configuration Options
>
> | Option    | Type   | Default | Description           |
> | --------- | ------ | ------- | --------------------- |
> | `version` | number | `1`     | Config schema version |
>
> ### Per-Script Configuration Options
>
> | Option       | Type                      | Default          | Description                                                                                                                                    |
> | ------------ | ------------------------- | ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
> | `command`    | string                    | required         | Script path or command                                                                                                                         |
> | `type`       | `"command"` \| `"prompt"` | `"command"`      | Hook execution type                                                                                                                            |
> | `timeout`    | number                    | platform default | Execution timeout in seconds                                                                                                                   |
> | `loop_limit` | number \| null            | `5`              | Per-script loop limit for stop/subagentStop hooks. `null` means no limit. Default is `5` for Cursor hooks, `null` for Claude Code hooks.       |
> | `failClosed` | boolean                   | `false`          | When `true`, hook failures (crash, timeout, invalid JSON) block the action instead of allowing it through. Useful for security-critical hooks. |
> | `matcher`    | object                    | -                | Filter criteria for when hook runs                                                                                                             |
>
> ### Matcher Configuration
>
> Matchers let you filter when a hook runs. Which field the matcher applies to depends on the hook:
>
> ```json
> {
>   "hooks": {
>     "preToolUse": [
>       {
>         "command": "./validate-shell.sh",
>         "matcher": "Shell"
>       }
>     ],
>     "subagentStart": [
>       {
>         "command": "./validate-explore.sh",
>         "matcher": "explore|shell"
>       }
>     ],
>     "beforeShellExecution": [
>       {
>         "command": "./approve-network.sh",
>         "matcher": "curl|wget|nc "
>       }
>     ]
>   }
> }
> ```
>
> - **subagentStart**: The matcher runs against the **subagent type** (e.g. `explore`, `shell`, `generalPurpose`). Use it to run hooks only when a specific kind of subagent is started. The example above runs `validate-explore.sh` only for explore or shell subagents.
> - **beforeShellExecution**: The matcher runs against the **shell command** string. Use it to run hooks only when the command matches a pattern (e.g. network calls, file deletions). The example above runs `approve-network.sh` only when the command contains `curl`, `wget`, or `nc `.
>
> **Available matchers by hook:**
>
> - **preToolUse / postToolUse / postToolUseFailure**: Filter by tool type. Values include `Shell`, `Read`, `Write`, `Grep`, `Delete`, `Task`, and MCP tools using the `MCP:<tool_name>` format.
> - **subagentStart / subagentStop**: Filter by subagent type (`generalPurpose`, `explore`, `shell`, etc.).
> - **beforeShellExecution / afterShellExecution**: Filter by the shell command text; the matcher is matched against the full command string.
> - **beforeReadFile**: Filter by tool type (`TabRead`, `Read`, etc.).
> - **afterFileEdit**: Filter by tool type (`TabWrite`, `Write`, etc.).
> - **beforeSubmitPrompt**: Matched against the value `UserPromptSubmit`.
> - **stop**: Matched against the value `Stop`.
> - **afterAgentResponse**: Matched against the value `AgentResponse`.
> - **afterAgentThought**: Matched against the value `AgentThought`.

### Source: Cursor Hooks — Team Distribution

> ## Team Distribution
>
> Hooks can be distributed to team members using project hooks (via version control), MDM tools, or Cursor's cloud distribution system.
>
> ### Project Hooks (Version Control)
>
> Project hooks are the simplest way to share hooks with your team. Place a `hooks.json` file at `<project-root>/.cursor/hooks.json` and commit it to your repository. When team members open the project in a trusted workspace, Cursor automatically loads and runs the project hooks.
>
> Cloud agents also load these project hooks when they work on your repository in
> the cloud.
>
> Project hooks:
>
> - Are stored in version control alongside your code
> - Automatically load for all team members in trusted workspaces
> - Can be project-specific (e.g., enforce formatting standards for a particular codebase)
> - Require the workspace to be trusted to run (for security)
>
> ### MDM Distribution
>
> Distribute hooks across your organization using Mobile Device Management (MDM) tools. Place the `hooks.json` file and hook scripts in the target directories on each machine.
>
> **User home directory** (per-user distribution):
>
> - `~/.cursor/hooks.json`
> - `~/.cursor/hooks/` (for hook scripts)
>
> **Global directories** (system-wide distribution):
>
> - macOS: `/Library/Application Support/Cursor/hooks.json`
> - Linux/WSL: `/etc/cursor/hooks.json`
> - Windows: `C:\\ProgramData\\Cursor\\hooks.json`
>
> Note: MDM-based distribution is fully managed by your organization. Cursor does not deploy or manage files through your MDM solution. Ensure your internal IT or security team handles configuration, deployment, and updates in accordance with your organization's policies.
>
> ### Cloud Distribution (Enterprise Only)
>
> Enterprise teams can use Cursor's native cloud distribution to automatically sync hooks to all team members. Configure hooks in the [web dashboard](https://cursor.com/dashboard/team-content?section=hooks). Cursor automatically delivers configured hooks to all client machines when team members log in.
>
> Cloud distribution provides:
>
> - Automatic synchronization to all team members (every thirty minutes)
> - Operating system targeting for platform-specific hooks
> - Centralized management through the dashboard
>
> Enterprise administrators can create, edit, and manage team hooks from the dashboard without requiring access to individual machines.
>
> [Contact sales](https://cursor.com/contact-sales?source=docs-hooks-cloud) to get Enterprise cloud hook distribution.

### Source: Cursor Hooks — Troubleshooting

> ## Troubleshooting
>
> **How to confirm hooks are active**
>
> There is a Hooks tab in **Customize** and a Hooks output channel to debug configured and executed hooks and see errors.
>
> **If hooks are not working**
>
> - Cursor watches `hooks.json` files and reloads them on save. If hooks still do not load, restart Cursor.
> - Check that relative paths are correct for your hook source:
>   - For **project hooks**, paths are relative to the **project root** (e.g., `.cursor/hooks/script.sh`)
>   - For **user hooks**, paths are relative to `~/.cursor/` (e.g., `./hooks/script.sh` or `hooks/script.sh`)
>
> **Exit code blocking**
>
> Exit code `2` from command hooks blocks the action (equivalent to returning `permission: "deny"`). This matches Claude Code behavior for compatibility.
>
> ### Enterprise hooks and distribution
>
> Cloud distribution and team-wide hook management are available on Enterprise.
>
>
