---
primary_sources:
  - id: T1-HOOKS
    title: "Hook Types and Configuration"
    url: "https://cursor.com/docs/hooks.md"
    section: "Hook Types and Configuration"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook execution and protocol

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Cursor Hooks — Hook Types

> ## Hook Types
>
> Hooks support two execution types: command-based (default) and prompt-based (LLM-evaluated).
>
> ### Command-Based Hooks
>
> Command hooks execute shell scripts that receive JSON input via stdin and return JSON output via stdout.
>
> ```json
> {
>   "hooks": {
>     "beforeShellExecution": [
>       {
>         "command": "./scripts/approve-network.sh",
>         "timeout": 30,
>         "matcher": "curl|wget|nc"
>       }
>     ]
>   }
> }
> ```
>
> **Exit code behavior:**
>
> - Exit code `0` - Hook succeeded, use the JSON output
> - Exit code `2` - Block the action (equivalent to returning `permission: "deny"`)
> - Other exit codes - Hook failed, action proceeds (fail-open by default)
>
> ### Prompt-Based Hooks
>
> Prompt hooks use an LLM to evaluate a natural language condition. They're useful for policy enforcement without writing custom scripts.
>
> ```json
> {
>   "hooks": {
>     "beforeShellExecution": [
>       {
>         "type": "prompt",
>         "prompt": "Does this command look safe to execute? Only allow read-only operations.",
>         "timeout": 10
>       }
>     ]
>   }
> }
> ```
>
> **Features:**
>
> - Returns structured `{ ok: boolean, reason?: string }` response
> - Uses a fast model for quick evaluation
> - `$ARGUMENTS` placeholder is auto-replaced with hook input JSON
> - If `$ARGUMENTS` is absent, hook input is auto-appended
> - Optional `model` field to override the default LLM model

### Source: Cursor Hooks — Configuration — Global and Per-Script Options

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
