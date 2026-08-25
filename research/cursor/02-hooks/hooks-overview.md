---
primary_sources:
  - id: T1-HOOKS
    title: "Intro and categories"
    url: "https://cursor.com/docs/hooks.md"
    section: "Intro and categories"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hooks overview

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Cursor Hooks — Introduction

> # Hooks
>
> Hooks let you observe, control, and extend the agent loop using custom scripts. Define hooks in `hooks.json` files at the project or user level, or install them through plugins from **Customize**. Hooks are spawned processes that communicate over stdio using JSON in both directions. They run before or after defined stages of the agent loop and can observe, block, or modify behavior.
>
> [Media](/docs-static/images/agent/hooks.mp4)
>
> With hooks, you can:
>
> - Run formatters after edits
> - Add analytics for events
> - Scan for PII or secrets
> - Gate risky operations (e.g., SQL writes)
> - Control subagent (Task tool) execution
> - Inject context at session start
>
> Looking for ready-to-use integrations? See [Partner Integrations](https://cursor.com/docs/hooks.md#partner-integrations) for security, governance, and secrets management solutions from our ecosystem partners.
>
> Cursor supports loading hooks from third-party tools like Claude Code. See [Third Party Hooks](https://cursor.com/docs/reference/third-party-hooks.md) for details on compatibility and configuration.

### Source: Cursor Hooks — Hook categories

> ## Hook categories
>
> Hooks fall into three categories based on what triggers them:
>
> **Agent hooks (Cmd+K/Agent Chat)** fire during an agent session:
>
> - `sessionStart` / `sessionEnd` - Session lifecycle management
> - `preToolUse` / `postToolUse` / `postToolUseFailure` - Generic tool use hooks (fires for all tools)
> - `subagentStart` / `subagentStop` - Subagent (Task tool) lifecycle
> - `beforeShellExecution` / `afterShellExecution` - Control shell commands
> - `beforeMCPExecution` / `afterMCPExecution` - Control MCP tool usage
> - `beforeReadFile` / `afterFileEdit` - Control file access and edits
> - `beforeSubmitPrompt` - Validate prompts before submission
> - `preCompact` - Observe context window compaction
> - `stop` - Handle agent completion
> - `afterAgentResponse` / `afterAgentThought` - Track agent responses
>
> **Tab hooks (inline completions)** fire for autonomous Tab operations:
>
> - `beforeTabFileRead` - Control file access for Tab completions
> - `afterTabFileEdit` - Post-process Tab edits
>
> **App lifecycle hooks** fire outside any agent session:
>
> - `workspaceOpen` - Fires when Cursor opens a workspace and on every workspace folder change. Can return additional plugin paths to load for the current workspace.
>
> These separate hook surfaces let you apply different policies to autonomous Tab operations, user-directed Agent operations, and workspace startup.

### Source: Cursor Hooks — Quickstart

> ## Quickstart
>
> Create a `hooks.json` file. You can create it at the project level (`<project>/.cursor/hooks.json`) or in your home directory (`~/.cursor/hooks.json`). Project-level hooks apply only to that specific project, while home directory hooks apply globally.
>
> ### User hooks (\~/.cursor/)
>
> For user-level hooks that apply globally, create `~/.cursor/hooks.json`:
>
> ```json
> {
>   "version": 1,
>   "hooks": {
>     "afterFileEdit": [{ "command": "./hooks/format.sh" }]
>   }
> }
> ```
>
> Create your hook script at `~/.cursor/hooks/format.sh`:
>
> ```bash
> #!/bin/bash
> # Read input, do something, exit 0
> cat > /dev/null
> exit 0
> ```
>
> Make it executable:
>
> ```bash
> chmod +x ~/.cursor/hooks/format.sh
> ```
>
> ### Project hooks (.cursor/)
>
> For project-level hooks that apply to a specific repository, create `<project>/.cursor/hooks.json`:
>
> ```json
> {
>   "version": 1,
>   "hooks": {
>     "afterFileEdit": [{ "command": ".cursor/hooks/format.sh" }]
>   }
> }
> ```
>
> Note: Project hooks run from the **project root**, so use `.cursor/hooks/format.sh` (not `./hooks/format.sh`).
>
> Create your hook script at `<project>/.cursor/hooks/format.sh`:
>
> ```bash
> #!/bin/bash
> # Read input, do something, exit 0
> cat > /dev/null
> exit 0
> ```
>
> Make it executable:
>
> ```bash
> chmod +x .cursor/hooks/format.sh
> ```
>
> Cursor watches hooks config files and reloads them automatically. Your hook runs after every file edit.

### Source: Cursor Hooks — Environment Variables

> ## Environment Variables
>
> Hook scripts receive environment variables when executed:
>
> | Variable                 | Description                                                   | Always Present         |
> | ------------------------ | ------------------------------------------------------------- | ---------------------- |
> | `CURSOR_PROJECT_DIR`     | Workspace root directory                                      | Yes                    |
> | `CURSOR_VERSION`         | Cursor version string                                         | Yes                    |
> | `CURSOR_USER_EMAIL`      | Authenticated user email                                      | If logged in           |
> | `CURSOR_TRANSCRIPT_PATH` | Path to the conversation transcript file                      | If transcripts enabled |
> | `CURSOR_CODE_REMOTE`     | Set to the string `"true"` when running in a remote workspace | For remote workspaces  |
> | `CLAUDE_PROJECT_DIR`     | Alias for project dir (Claude compatibility)                  | Yes                    |
>
> Session-scoped environment variables from `sessionStart` hooks are passed to all subsequent hook executions within that session.
