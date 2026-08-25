---
primary_sources:
  - id: T1-HOOKS
    title: "Cursor Hooks"
    url: "https://cursor.com/docs/hooks.md"
    section: "Tab and lifecycle hook events"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Tab and lifecycle hook events

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Cursor Hooks — Reference (Tab and lifecycle)

> #### beforeTabFileRead
>
> Called before Tab (inline completions) reads a file. Enable redaction or access control before Tab accesses file contents.
> **Key differences from `beforeReadFile`:**
> - Only triggered by Tab, not Agent
> - Does not include `attachments` field (Tab doesn't use prompt attachments)
> - Useful for applying different policies to autonomous Tab operations
> ```json
> // Input
> {
>   "file_path": "<absolute path>",
>   "content": "<file contents>"
> }
> // Output
>   "permission": "allow" | "deny"
> ```
> #### afterTabFileEdit
> Called after Tab (inline completions) edits a file. Useful for formatters or auditing of Tab-written code.
> **Key differences from `afterFileEdit`:**
> - Includes detailed edit information: `range`, `old_line`, and `new_line` for precise edit tracking
> - Useful for fine-grained formatting or analysis of Tab edits
>   "edits": [
>     {
>       "old_string": "<search>",
>       "new_string": "<replace>",
>       "range": {
>         "start_line_number": 10,
>         "start_column": 5,
>         "end_line_number": 10,
>         "end_column": 20
>       },
>       "old_line": "<line before edit>",
>       "new_line": "<line after edit>"
>     }
>   ]
>   // No output fields currently supported
> #### beforeSubmitPrompt
> Called right after user hits send but before backend request. Can prevent submission.
>   "prompt": "<user prompt text>",
>   "attachments": [
>       "type": "file" | "rule",
>       "file_path": "<absolute path>"
>   "continue": true | false,
>   "user_message": "<message shown to user when blocked>"
> | Output Field   | Type              | Description                                          |
> | -------------- | ----------------- | ---------------------------------------------------- |
> | `continue`     | boolean           | Whether to allow the prompt submission to proceed    |
> | `user_message` | string (optional) | Message shown to the user when the prompt is blocked |
> #### afterAgentResponse
> Called after the agent has completed an assistant message.
>   "text": "<assistant final text>"
> #### afterAgentThought
> Called after the agent completes a thinking block. Useful for observing the agent's reasoning process.
>   "text": "<fully aggregated thinking text>",
>   "duration_ms": 5000
> | Field         | Type              | Description                                            |
> | ------------- | ----------------- | ------------------------------------------------------ |
> | `text`        | string            | Fully aggregated thinking text for the completed block |
> | `duration_ms` | number (optional) | Duration in milliseconds for the thinking block        |
> #### stop
> Called when the agent loop ends. Can optionally auto-submit a follow-up user message to keep iterating.
>   "status": "completed" | "aborted" | "error",
>   "loop_count": 0
>   "followup_message": "<message text>"
> - The optional `followup_message` is a string. When provided and non-empty, Cursor will automatically submit it as the next user message. This enables loop-style flows (e.g., iterate until a goal is met).
> - The `loop_count` field indicates how many times the stop hook has already triggered an automatic follow-up for this conversation (starts at 0). The default limit is 5 auto follow-ups per script, configurable via the `loop_limit` option. Set `loop_limit` to `null` to remove the cap. The same limit applies to `subagentStop` follow-ups.
> #### sessionStart
> Called when a new composer conversation is created. This hook runs as fire-and-forget; the agent loop does not wait for or enforce a blocking response. Use it to set up session-specific environment variables or inject additional context.
>   "session_id": "<unique session identifier>",
>   "is_background_agent": true | false,
>   "composer_mode": "agent" | "ask" | "edit"
>   "env": { "<key>": "<value>" },
>   "additional_context": "<context to add to conversation>"
> | Input Field           | Type              | Description                                                         |
> | --------------------- | ----------------- | ------------------------------------------------------------------- |
> | `session_id`          | string            | Unique identifier for this session (same as `conversation_id`)      |
> | `is_background_agent` | boolean           | Whether this is a background agent session vs interactive session   |
> | `composer_mode`       | string (optional) | The mode the composer is starting in (e.g., "agent", "ask", "edit") |
> | Output Field         | Type              | Description                                                                                |
> | -------------------- | ----------------- | ------------------------------------------------------------------------------------------ |
> | `env`                | object (optional) | Environment variables to set for this session. Available to all subsequent hook executions |
> | `additional_context` | string (optional) | Additional context to add to the conversation's initial system context                     |
> The schema also accepts `continue` and `user_message` fields, but current callers do not enforce them. Session creation is not blocked even when `continue` is `false`.
> #### sessionEnd
> Called when a composer conversation ends. This is a fire-and-forget hook useful for logging, analytics, or cleanup tasks. The response is logged but not used.
>   "reason": "completed" | "aborted" | "error" | "window_close" | "user_close",
>   "duration_ms": 45000,
>   "final_status": "<status string>",
>   "error_message": "<error details if reason is 'error'>"
>   // No output fields - fire and forget
> | Input Field           | Type              | Description                                                                               |
> | --------------------- | ----------------- | ----------------------------------------------------------------------------------------- |
> | `session_id`          | string            | Unique identifier for the session that is ending                                          |
> | `reason`              | string            | How the session ended: "completed", "aborted", "error", "window\_close", or "user\_close" |
> | `duration_ms`         | number            | Total duration of the session in milliseconds                                             |
> | `is_background_agent` | boolean           | Whether this was a background agent session                                               |
> | `final_status`        | string            | Final status of the session                                                               |
> | `error_message`       | string (optional) | Error message if reason is "error"                                                        |
> #### preCompact
> Called before context window compaction/summarization occurs. This is an observational hook that cannot block or modify the compaction behavior. Useful for logging when compaction happens or notifying users.
>   "trigger": "auto" | "manual",
>   "context_usage_percent": 85,
>   "context_tokens": 120000,
>   "context_window_size": 128000,
>   "message_count": 45,
>   "messages_to_compact": 30,
>   "is_first_compaction": true | false
>   "user_message": "<message to show when compaction occurs>"
> | Input Field             | Type    | Description                                                |
> | ----------------------- | ------- | ---------------------------------------------------------- |
> | `trigger`               | string  | What triggered the compaction: "auto" or "manual"          |
> | `context_usage_percent` | number  | Current context window usage as a percentage (0-100)       |
> | `context_tokens`        | number  | Current context window token count                         |
> | `context_window_size`   | number  | Maximum context window size in tokens                      |
> | `message_count`         | number  | Number of messages in the conversation                     |
> | `messages_to_compact`   | number  | Number of messages that will be summarized                 |
> | `is_first_compaction`   | boolean | Whether this is the first compaction for this conversation |
> | Output Field   | Type              | Description                                        |
> | -------------- | ----------------- | -------------------------------------------------- |
> | `user_message` | string (optional) | Message to show to the user when compaction occurs |
> #### workspaceOpen
> Fires once when Cursor opens a workspace and again on every workspace folder change. Skipped when the window has zero workspace folders. Runs in the Cursor desktop app and CLI.
>   "hook_event_name": "workspaceOpen",
>   "cursor_version": "string",
>   "workspace_roots": ["<absolute path>"],
>   "user_email": "string | null"
>   "pluginPaths": ["<absolute path>", "..."]
> | Output Field  | Type                 | Description                                                             |
> | ------------- | -------------------- | ----------------------------------------------------------------------- |
> | `pluginPaths` | string\[] (optional) | Absolute paths to plugin directories to load for the current workspace. |
> #### beforeTabFileRead
> #### afterTabFileEdit
> #### workspaceOpen
