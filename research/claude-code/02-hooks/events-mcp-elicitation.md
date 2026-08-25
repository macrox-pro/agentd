---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Hook events — MCP elicitation"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook events — MCP elicitation

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — Hook events

> ## Hook events
>
> Each event corresponds to a point in Claude Code's lifecycle where hooks can run. The sections below are ordered to match the lifecycle: from session setup through the agentic loop to session end. Each section describes when the event fires, what matchers it supports, the JSON input it receives, and how to control behavior through output.

### Source: Hooks reference — Elicitation

> ### Elicitation
>
> Runs when an MCP server requests user input mid-task. By default, Claude Code shows an interactive dialog for the user to respond. Hooks can intercept this request and respond programmatically, skipping the dialog entirely.
>
> The matcher field matches against the MCP server name.
>
> #### Elicitation input
>
> In addition to the [common input fields](#common-input-fields), Elicitation hooks receive `mcp_server_name`, `message`, and optional `mode`, `url`, `elicitation_id`, and `requested_schema` fields.
>
> For form-mode elicitation, the most common case:
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "Elicitation",
>   "mcp_server_name": "my-mcp-server",
>   "message": "Please provide your credentials",
>   "mode": "form",
>   "requested_schema": {
>     "type": "object",
>     "properties": {
>       "username": { "type": "string", "title": "Username" }
>     }
>   }
> }
> ```
>
> For URL-mode elicitation, used for browser-based authentication:
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "Elicitation",
>   "mcp_server_name": "my-mcp-server",
>   "message": "Please authenticate",
>   "mode": "url",
>   "url": "https://auth.example.com/login"
> }
> ```
>
> #### Elicitation output
>
> To respond programmatically without showing the dialog, return a JSON object with `hookSpecificOutput`:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "Elicitation",
>     "action": "accept",
>     "content": {
>       "username": "alice"
>     }
>   }
> }
> ```
>
> | Field     | Values                        | Description                                                      |
> | :-------- | :---------------------------- | :--------------------------------------------------------------- |
> | `action`  | `accept`, `decline`, `cancel` | Whether to accept, decline, or cancel the request                |
> | `content` | object                        | Form field values to submit. Only used when `action` is `accept` |
>
> Exit code 2 denies the elicitation. Claude Code doesn't show your stderr message anywhere.
>
> Claude Code acts on `hookSpecificOutput` from an Elicitation hook's JSON output and discards `systemMessage` and `continue`.

### Source: Hooks reference — ElicitationResult

> ### ElicitationResult
>
> Runs after a user responds to an MCP elicitation. Hooks can observe, modify, or block the response before it is sent back to the MCP server.
>
> The matcher field matches against the MCP server name.
>
> #### ElicitationResult input
>
> In addition to the [common input fields](#common-input-fields), ElicitationResult hooks receive `mcp_server_name`, `action`, and optional `mode`, `elicitation_id`, and `content` fields.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "ElicitationResult",
>   "mcp_server_name": "my-mcp-server",
>   "action": "accept",
>   "content": { "username": "alice" },
>   "mode": "form",
>   "elicitation_id": "elicit-123"
> }
> ```
>
> #### ElicitationResult output
>
> To override the user's response, return a JSON object with `hookSpecificOutput`:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "ElicitationResult",
>     "action": "decline",
>     "content": {}
>   }
> }
> ```
>
> | Field     | Values                        | Description                                                            |
> | :-------- | :---------------------------- | :--------------------------------------------------------------------- |
> | `action`  | `accept`, `decline`, `cancel` | Overrides the user's action                                            |
> | `content` | object                        | Overrides form field values. Only meaningful when `action` is `accept` |
>
> Exit code 2 blocks the response, changing the effective action to `decline`. Claude Code doesn't show your stderr message anywhere.
>
> Claude Code acts on `hookSpecificOutput` from an ElicitationResult hook's JSON output and discards `systemMessage` and `continue`.
>
> ## Prompt-based hooks
>
> In addition to command, HTTP, and MCP tool hooks, Claude Code supports prompt-based hooks (`type: "prompt"`) that use an LLM to evaluate whether to allow or block an action, and agent hooks (`type: "agent"`) that spawn an agentic verifier with tool access. Not all events support every hook type.
>
> Events that support all five hook types (`command`, `http`, `mcp_tool`, `prompt`, and `agent`):
>
> * `PermissionDenied`
> * `PermissionRequest`
> * `PostToolBatch`
> * `PostToolUse`
> * `PostToolUseFailure`
> * `PreToolUse`
> * `Stop`
> * `SubagentStop`
> * `TaskCompleted`
> * `TaskCreated`
> * `TeammateIdle`
> * `UserPromptExpansion`
> * `UserPromptSubmit`
>
> Events that support `command`, `http`, and `mcp_tool` hooks but not `prompt` or `agent`:
>
> * `ConfigChange`
> * `CwdChanged`
> * `DirectoryAdded`
> * `Elicitation`
> * `ElicitationResult`
> * `FileChanged`
> * `InstructionsLoaded`
> * `MessageDisplay`
> * `Notification`
> * `PostCompact`
> * `PreCompact`
> * `SessionEnd`
> * `StopFailure`
> * `SubagentStart`
> * `WorktreeCreate`
> * `WorktreeRemove`
>
> `SessionStart` and `Setup` support `command` and `mcp_tool` hooks. They don't support `http`, `prompt`, or `agent` hooks.
