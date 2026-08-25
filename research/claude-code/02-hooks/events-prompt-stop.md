---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Hook events — prompt and stop"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook events — prompt and stop

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — Hook events
>
> ## Hook events
>
> Each event corresponds to a point in Claude Code's lifecycle where hooks can run. The sections below are ordered to match the lifecycle: from session setup through the agentic loop to session end. Each section describes when the event fires, what matchers it supports, the JSON input it receives, and how to control behavior through output.### Source: Hooks reference — UserPromptSubmit
>
> ### UserPromptSubmit
>
> Runs when the user submits a prompt, before Claude processes it. This allows you
> to add additional context based on the prompt/conversation, validate prompts, or
> block certain types of prompts.
>
> `UserPromptSubmit` hooks have a default timeout of 30 seconds for `command`, `http`, and `mcp_tool` types, shorter than the 600-second default for those types on most other events. Because this hook runs before every prompt and blocks model processing until it completes, a stuck hook stalls the session. If your hook needs more time, set the `timeout` field in the hook entry.
>
> A `UserPromptSubmit` command, HTTP, or MCP tool hook that reaches its timeout is canceled and its output, including any `additionalContext`, is discarded. The prompt still reaches Claude without that context. The transcript shows a notice naming the hook, the timeout that fired, and that the output was discarded.
>
> An [Agent SDK callback hook](/docs/en/agent-sdk/hooks) on `UserPromptSubmit` that reaches its timeout blocks the prompt with a message naming the hook and the timeout, because a callback there can be acting as a policy gate that must not fail open. The session continues. Before v2.1.208, a callback timeout on that event ended the turn with an execution error.
>
> #### UserPromptSubmit input
>
> In addition to the [common input fields](#common-input-fields), UserPromptSubmit hooks receive the `prompt` field containing the text the user submitted.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "UserPromptSubmit",
>   "prompt": "Write a function to calculate the factorial of a number"
> }
> ```
>
> #### UserPromptSubmit decision control
>
> `UserPromptSubmit` hooks can control whether a user prompt is processed and add context. All [JSON output fields](#json-output) are available.
>
> There are two ways to add context to the conversation on exit code 0:
>
> * **Plain text stdout**: Claude Code adds stdout it [treats as plain text](#exit-code-0) to Claude's context
> * **JSON with `additionalContext`**: use the JSON format below for more control. The `additionalContext` field is added as context
>
> Neither channel produces a visible transcript entry. Plain stdout and the `additionalContext` value are each injected as a system reminder that starts with the hook's name; Claude reads both. To confirm delivery, check the [debug log](#debug-hooks).
>
> To block a prompt, return a JSON object with `decision` set to `"block"`:
>
> | Field                    | Description                                                                                                            |
> | :----------------------- | :--------------------------------------------------------------------------------------------------------------------- |
> | `decision`               | `"block"` prevents the prompt from being processed and erases it from context. Omit to allow the prompt to proceed     |
> | `reason`                 | Shown to the user when `decision` is `"block"`. Not added to context                                                   |
> | `additionalContext`      | String added to Claude's context alongside the submitted prompt. See [Add context for Claude](#add-context-for-claude) |
> | `sessionTitle`           | Sets the session title. Use to name sessions automatically based on the prompt content                                 |
> | `suppressOriginalPrompt` | If `true` when `decision` is `"block"`, omits the original prompt text from the block message shown to the user        |
>
> A hook that blocks by exiting 2 routes the same way as `reason`: the block message shows the stderr text to the user, and it isn't added to context.
>
> ```json
> {
>   "decision": "block",
>   "reason": "Explanation for decision",
>   "hookSpecificOutput": {
>     "hookEventName": "UserPromptSubmit",
>     "additionalContext": "My additional context here",
>     "sessionTitle": "My session title"
>   }
> }
> ```### Source: Hooks reference — UserPromptExpansion
>
> ### UserPromptExpansion
>
> Runs when a user-typed command expands into a prompt before reaching Claude. Use this to block specific commands from direct invocation, inject context for a particular skill, or log which commands users invoke. For example, a hook matching `deploy` can block `/deploy` unless an approval file is present, or a hook matching a review skill can append the team's review checklist as `additionalContext`.
>
> This event covers the path `PreToolUse` doesn't: a `PreToolUse` hook matching the `Skill` tool fires only when Claude calls the tool, but typing `/skillname` directly bypasses `PreToolUse`. `UserPromptExpansion` fires on that direct path.
>
> Matches on `command_name`. Leave the matcher empty to fire on every prompt-type command.
>
> #### UserPromptExpansion input
>
> In addition to the [common input fields](#common-input-fields), UserPromptExpansion hooks receive `expansion_type`, `command_name`, `command_args`, `command_source`, and the original `prompt` string. The `expansion_type` field is `slash_command` for skill and custom commands, or `mcp_prompt` for MCP server prompts.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../00893aaf.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "UserPromptExpansion",
>   "expansion_type": "slash_command",
>   "command_name": "example-skill",
>   "command_args": "arg1 arg2",
>   "command_source": "plugin",
>   "prompt": "/example-skill arg1 arg2"
> }
> ```
>
> #### UserPromptExpansion decision control
>
> `UserPromptExpansion` hooks can block the expansion or add context. All [JSON output fields](#json-output) are available.
>
> | Field               | Description                                                                                                           |
> | :------------------ | :-------------------------------------------------------------------------------------------------------------------- |
> | `decision`          | `"block"` prevents the command from expanding. Omit to allow it to proceed                                            |
> | `reason`            | Shown to the user when `decision` is `"block"`                                                                        |
> | `additionalContext` | String added to Claude's context alongside the expanded prompt. See [Add context for Claude](#add-context-for-claude) |
>
> A hook that blocks by exiting 2 routes the same way as `reason`: the block message shows the stderr text to the user.
>
> ```json
> {
>   "decision": "block",
>   "reason": "This slash command is not available",
>   "hookSpecificOutput": {
>     "hookEventName": "UserPromptExpansion",
>     "additionalContext": "Additional context for this expansion"
>   }
> }
> ```### Source: Hooks reference — Stop
>
> ### Stop
>
> Runs when the main Claude Code agent has finished responding. Does not run if
> the stoppage occurred due to a user interrupt. API errors fire
> [StopFailure](#stopfailure) instead.
>
>   The [`/goal`](/docs/en/goal) command is a built-in shortcut for a session-scoped prompt-based Stop hook. Use it when you want Claude to keep working toward a condition without writing hook configuration.
>
> #### Stop input
>
> In addition to the [common input fields](#common-input-fields), Stop hooks receive `stop_hook_active`, `last_assistant_message`, `background_tasks`, and `session_crons`. The `stop_hook_active` field is `true` when Claude Code is already continuing as a result of a stop hook. Check this value or process the transcript to avoid blocking on a condition that will never resolve. Claude Code overrides the hook and ends the turn after 8 consecutive blocks.
>
> The `last_assistant_message` field contains the text content of Claude's final response, so hooks can access it without parsing the transcript file. For hooks that act on the just-completed turn, such as read-aloud or notification hooks, use this field rather than reading `transcript_path`: the transcript file isn't guaranteed to include the final message at Stop time on all versions.
>
> The `background_tasks` and `session_crons` arrays let hooks distinguish "session is done" from "session is paused waiting for background work to wake it back up". Both arrays are present when the task registry is reachable and are empty when nothing is in flight or scheduled.
>
> Each entry in `background_tasks` describes one in-flight task and uses these fields:
>
> | Field         | Description                                                                                                                                                                                                                                          |
> | :------------ | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `id`          | Task identifier                                                                                                                                                                                                                                      |
> | `type`        | Friendly task-type label such as `shell`, `subagent`, `monitor`, `workflow`, `teammate`, `cloud session`, or `MCP task`. Each label identifies which Claude Code feature created the task. Falls back to the raw discriminant for unrecognized types |
> | `status`      | Current task status                                                                                                                                                                                                                                  |
> | `description` | Free-text description, capped at 1000 characters with an in-string `… [+N chars]` marker when clipped                                                                                                                                                |
> | `command`     | Shell command line, capped at 1000 characters. Present only for `shell` tasks                                                                                                                                                                        |
> | `agent_type`  | Subagent type name. Present only for `subagent` tasks                                                                                                                                                                                                |
> | `server`      | MCP server name. Present only for `monitor` and `MCP task` tasks                                                                                                                                                                                     |
> | `tool`        | MCP tool name. Present only for `monitor` and `MCP task` tasks                                                                                                                                                                                       |
> | `name`        | Workflow name. Present only for `workflow` tasks                                                                                                                                                                                                     |
>
> Each entry in `session_crons` describes one session-scoped scheduled wakeup, sourced from `CronCreate`, `ScheduleWakeup`, and `/loop`:
>
> | Field       | Description                                                                                                          |
> | :---------- | :------------------------------------------------------------------------------------------------------------------- |
> | `id`        | Cron task identifier                                                                                                 |
> | `schedule`  | Cron expression, for example `0 9 * * 1-5`                                                                           |
> | `recurring` | `false` for one-shot wakeups whose schedule encodes a single fire time, `true` for tasks that re-fire on every match |
> | `prompt`    | Prompt submitted when the cron fires, capped at 1000 characters with the same `… [+N chars]` marker                  |
>
> This example shows a Stop input with one in-flight shell task and one recurring cron:
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "~/.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "Stop",
>   "stop_hook_active": true,
>   "last_assistant_message": "I've completed the refactoring. Here's a summary...",
>   "background_tasks": [
>     {
>       "id": "task-001",
>       "type": "shell",
>       "status": "running",
>       "description": "tail logs",
>       "command": "tail -f /var/log/syslog"
>     }
>   ],
>   "session_crons": [
>     {
>       "id": "cron-001",
>       "schedule": "0 9 * * 1-5",
>       "recurring": true,
>       "prompt": "check the build"
>     }
>   ]
> }
> ```
>
> #### Stop decision control
>
> `Stop` and `SubagentStop` hooks can control whether Claude continues. In addition to the [JSON output fields](#json-output) available to all hooks, your hook script can return these event-specific fields:
>
> | Field                                  | Description                                                                                                                                                                               |
> | :------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `decision`                             | `"block"` prevents Claude from stopping. Omit to allow Claude to stop                                                                                                                     |
> | `reason`                               | Required when `decision` is `"block"`. Tells Claude why it should continue                                                                                                                |
> | `hookSpecificOutput.additionalContext` | Non-error feedback for Claude. The conversation continues so Claude can act on it, but unlike `decision: "block"` it is shown in the transcript as hook feedback rather than a hook error |
>
> A hook that blocks by exiting 2 routes the same way as `reason`: Claude receives the stderr message as the explanation for why it should continue.
>
> ```json
> {
>   "decision": "block",
>   "reason": "Must be provided when Claude is blocked from stopping"
> }
> ```
>
> Use `additionalContext` when the hook is working as designed and giving Claude guidance, such as "run the test suite before finishing". It keeps the conversation going through the same loop protections as `decision: "block"`, namely the `stop_hook_active` input and the 8-consecutive-continuation cap, but the transcript labels it `Stop hook feedback` and no hook error notification is shown:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "Stop",
>     "additionalContext": "Please run the test suite before finishing"
>   }
> }
> ```### Source: Hooks reference — StopFailure
>
> ### StopFailure
>
> Runs instead of [Stop](#stop) when the turn ends due to an API error. Claude Code ignores the hook's output and exit code, apart from [`terminalSequence`](#emit-terminal-notifications). Use this to log failures, send alerts, or take recovery actions when Claude can't complete a response due to rate limits, authentication problems, or other API errors.
>
> #### StopFailure input
>
> In addition to the [common input fields](#common-input-fields), StopFailure hooks receive `error`, optional `error_details`, and optional `last_assistant_message`. The `error` field identifies the error type and is used for matcher filtering.
>
> | Field                    | Description                                                                                                                                                                                                                                      |
> | :----------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `error`                  | Error type: `rate_limit`, `overloaded`, `authentication_failed`, `oauth_org_not_allowed`, `billing_error`, `invalid_request`, `model_not_found`, `server_error`, `max_output_tokens`, or `unknown`                                               |
> | `error_details`          | Additional details about the error, when available                                                                                                                                                                                               |
> | `last_assistant_message` | The rendered error text shown in the conversation. Unlike `Stop` and `SubagentStop`, where this field holds Claude's conversational output, for `StopFailure` it contains the API error string itself, such as `"API Error: Rate limit reached"` |
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "hook_event_name": "StopFailure",
>   "error": "rate_limit",
>   "error_details": "429 Too Many Requests",
>   "last_assistant_message": "API Error: Rate limit reached"
> }
> ```
>
> StopFailure hooks have no decision control. They run for notification and logging purposes only.### Source: Hooks reference — Notification
>
> ### Notification
>
> Runs when Claude Code sends notifications. Matches on notification type. Omit the matcher to run hooks for all notification types.
>
> You receive these hook events even with desktop notifications turned off: the `preferredNotifChannel` setting, including `notifications_disabled`, changes only how you're alerted, not whether your hook runs.
>
> | Matcher                      | When it fires                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
> | :--------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `permission_prompt`          | Claude needs you to approve a tool use and the prompt has waited about six seconds                                                                                                                                                                                                                                                                                                                                                                                    |
> | `idle_prompt`                | Claude finished responding about 60 seconds ago and you haven't typed since                                                                                                                                                                                                                                                                                                                                                                                           |
> | `auth_success`               | Authentication completes                                                                                                                                                                                                                                                                                                                                                                                                                                              |
> | `elicitation_dialog`         | An MCP server opens an elicitation form and you haven't typed for about six seconds                                                                                                                                                                                                                                                                                                                                                                                   |
> | `elicitation_url_dialog`     | An MCP server asks you to open a browser URL and you haven't typed for about six seconds                                                                                                                                                                                                                                                                                                                                                                              |
> | `elicitation_complete`       | An MCP elicitation form is submitted or dismissed                                                                                                                                                                                                                                                                                                                                                                                                                     |
> | `elicitation_response`       | An MCP elicitation response is sent back to the server                                                                                                                                                                                                                                                                                                                                                                                                                |
> | `agent_needs_input`          | A background session starts waiting on your input. Fires only while [agent view](/docs/en/agent-view) is open in a terminal                                                                                                                                                                                                                                                                                                                                                |
> | `agent_completed`            | A background session finishes or fails. Fires only while [agent view](/docs/en/agent-view) is open in a terminal                                                                                                                                                                                                                                                                                                                                                           |
> | `quota_auto_resume_fired`    | Claude Code continues your task after a claude.ai usage limit paused it: at the reset, or sooner when something you do in Claude Code during the wait, such as adding usage credits, upgrading your plan, or switching models, makes usage available again, with the [model-setting exception](/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset)                                                                                                                  |
> | `quota_auto_resume_stale`    | A claude.ai usage limit reset while your computer slept for more than about 30 minutes. Claude Code waits for you to press `Enter` instead of continuing. After a shorter sleep it continues and fires `quota_auto_resume_fired` instead                                                                                                                                                                                                                              |
> | `quota_auto_resume_disabled` | Claude Code ends its wait for a claude.ai usage limit without continuing your task: [`autoContinueAtUsageLimit`](/docs/en/settings-reference#autocontinueatusagelimit) turned off or the reset moved more than 24 hours away during a wait Claude Code started on its own, the continued task kept hitting the limit, or the continuation was blocked before it reached the model. Doesn't fire when you press `Esc` or `Ctrl+C`, or pick **Don't continue automatically** |
>
> The `agent_needs_input` and `agent_completed` types require Claude Code v2.1.198 or later.
>
> The `quota_auto_resume_fired`, `quota_auto_resume_stale`, and `quota_auto_resume_disabled` types require Claude Code v2.1.234 or later.
>
>   The `permission_prompt`, `idle_prompt`, `elicitation_dialog`, and `elicitation_url_dialog` types share their timing with desktop notifications, so in terminal sessions you only see them when you appear to be away from the terminal:
>
>   * Expect `permission_prompt` once you haven't typed for about six seconds. The timer starts when the permission prompt appears, and each keystroke defers it. To run a hook immediately on every permission ask, use [PermissionRequest](#permissionrequest) instead.
>   * Expect `idle_prompt` about 60 seconds after Claude finishes responding, and only if you haven't typed since. Claude Code doesn't send `idle_prompt` while it waits for a claude.ai usage limit to reset. When the wait ends on its own, one of the `quota_auto_resume_*` types fires instead.
>   * Expect `elicitation_dialog` for an elicitation form, or `elicitation_url_dialog` for a browser URL request, once you haven't typed for about six seconds. Both share the same six-second gate as `permission_prompt`: the timer starts when the dialog appears, and each keystroke defers it.
>
> Claude Code times `permission_prompt` differently in sessions where it sends permission requests to the Agent SDK's [`canUseTool` callback](/docs/en/agent-sdk/user-input), which is how Claude Desktop and the VS Code extension host Claude Code:
>
> * Expect `permission_prompt` about six seconds after Claude asks for permission. Claude Code doesn't defer it while you type.
> * If you or a [PermissionRequest](#permissionrequest) hook answer sooner, Claude Code doesn't run `permission_prompt`.
> * Set [`CLAUDE_CODE_DISABLE_PERMISSION_PROMPT_NOTIFY_HOOKS`](/docs/en/env-vars) to `1` to turn `permission_prompt` off in these sessions.
>
> Before v2.1.233, `permission_prompt` didn't fire in these sessions.
>
> Use separate matchers to run different handlers depending on the notification type. This configuration triggers a permission-specific alert script when Claude needs permission approval and a different notification when Claude has been idle:
>
> ```json
> {
>   "hooks": {
>     "Notification": [
>       {
>         "matcher": "permission_prompt",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/path/to/permission-alert.sh"
>           }
>         ]
>       },
>       {
>         "matcher": "idle_prompt",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/path/to/idle-notification.sh"
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> #### Notification input
>
> In addition to the [common input fields](#common-input-fields), Notification hooks receive `message` with the notification text, an optional `title`, and `notification_type` indicating which type fired.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "hook_event_name": "Notification",
>   "message": "Claude needs your permission",
>   "title": "Permission needed",
>   "notification_type": "permission_prompt"
> }
> ```
>
> Notification hooks can't block or modify notifications. Claude Code discards their `systemMessage` and `continue` fields but still emits [`terminalSequence`](#emit-terminal-notifications), which is what the desktop notification example relies on. Notification hooks are intended for side effects such as forwarding the notification to an external service.### Source: Hooks reference — MessageDisplay
>
> ### MessageDisplay
>
> Runs while an assistant message streams to the screen. Claude Code displays the message in increments: each time a batch of newly completed lines is ready to render, the hook runs once with those lines and Claude Code renders the hook's replacement text in their place. A long message produces several calls; a short message may produce only one.
>
> Use MessageDisplay to:
>
> * strip markdown for a minimal display
> * transform the text an Agent SDK application shows its users
> * redact API keys or internal hostnames from Claude's responses
>
> Claude Code holds each batch until your hook returns, so keep the hook fast. If the hook fails or times out, Claude Code displays the original text. The default timeout for this event is 10 seconds; if your hook needs more time, set the `timeout` field in the hook entry.
>
> MessageDisplay is display-only: the replacement text changes only what is rendered on screen. The transcript and what Claude sees keep the original text, so Claude never sees the replacement, and verbose mode shows the original. The hook receives assistant message text only, so tool results and the text you type render unchanged.
>
> MessageDisplay doesn't support matchers and fires for every assistant message that streams text; messages with no text, such as tool-call-only responses, don't trigger it.
>
> In non-interactive runs, including Agent SDK queries and `claude -p`, MessageDisplay runs once per assistant message instead of once per batch of lines. The single call arrives after the message completes and carries the full message text: `index` is `0`, `final` is `true`, and `delta` holds the entire message. A hook that collects the `delta` text for each message receives the same total text in both modes.
>
> #### MessageDisplay input
>
> In addition to the [common input fields](#common-input-fields), MessageDisplay hooks receive identifiers for the turn and message, the position of this call within the message, and the new text in `delta`. Batch boundaries depend on how the text streams, so use `index` and `final` to track progress through a message rather than expecting lines to be grouped a particular way.
>
> | Field        | Description                                                                                                                                                                                                                                                                                                                                                                                       |
> | :----------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `turn_id`    | UUID of the current turn                                                                                                                                                                                                                                                                                                                                                                          |
> | `message_id` | UUID of the assistant message being displayed. Stable across every batch of the same message. This is not the API `msg_…` id, so it can't be correlated with transcript message ids                                                                                                                                                                                                               |
> | `index`      | Zero-based index of this batch within the message                                                                                                                                                                                                                                                                                                                                                 |
> | `final`      | `true` on the message's last batch. Each message has exactly one final batch                                                                                                                                                                                                                                                                                                                      |
> | `delta`      | The newly completed lines since the prior batch, terminating newlines included. Always whole lines, except the final batch which may end mid-line. In interactive runs, the final batch's delta is empty when the message ends on a newline, so treat `final`, not a non-empty delta, as the end-of-message signal. In Agent SDK and `claude -p` runs, the single call carries the entire message |
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../transcript.jsonl",
>   "cwd": "/Users/my-project",
>   "hook_event_name": "MessageDisplay",
>   "turn_id": "0c9e6a2f-7d41-4f4e-9a15-3f4f7c2b8d10",
>   "message_id": "5b2a9c8e-1f63-4d8a-b7c4-9e0d2a6f1c3b",
>   "index": 0,
>   "final": false,
>   "delta": "Here is the plan:\n"
> }
> ```
>
> #### MessageDisplay output
>
> In addition to the [JSON output fields](#json-output) available to all hooks, MessageDisplay hooks can return `displayContent` to replace the delta on screen:
>
> | Field            | Description                                                           |
> | :--------------- | :-------------------------------------------------------------------- |
> | `displayContent` | Text displayed in place of the delta. Omit it to display the original |
>
> MessageDisplay hooks have no decision control. They can't block the message or change what is stored in the transcript or sent to Claude. Claude Code acts on `displayContent` from their JSON output and discards `systemMessage` and `continue`.
>
> This example strips markdown formatting from Claude's responses for a plain-text display. The script reads each batch from stdin, removes bold markers and inline code backticks from `delta`, and returns the result as `displayContent`.
>
>   #### macOS/Linux
>
> Register a command hook for the event in your settings file:
>
>     ```json
>     {
>       "hooks": {
>         "MessageDisplay": [
>           {
>             "hooks": [
>               {
>                 "type": "command",
>                 "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/plain-display.sh",
>                 "args": []
>               }
>             ]
>           }
>         ]
>       }
>     }
>     ```
>
>     Save this script to `.claude/hooks/plain-display.sh` in your project and make it executable with `chmod +x`:
>
>     ```bash
>     #!/bin/bash
>     jq '{hookSpecificOutput: {hookEventName: "MessageDisplay", displayContent: (.delta | gsub("\\*\\*"; "") | gsub("`"; ""))}}'
>     ```
>
>   #### Windows (PowerShell)
>
> Register a command hook that runs the script through PowerShell:
>
>     ```json
>     {
>       "hooks": {
>         "MessageDisplay": [
>           {
>             "hooks": [
>               {
>                 "type": "command",
>                 "command": "powershell.exe",
>                 "args": [
>                   "-NoProfile",
>                   "-ExecutionPolicy",
>                   "Bypass",
>                   "-File",
>                   "${CLAUDE_PROJECT_DIR}/.claude/hooks/plain-display.ps1"
>                 ]
>               }
>             ]
>           }
>         ]
>       }
>     }
>     ```
>
>     The `-NoProfile` flag skips loading your PowerShell profile so the hook starts fast, and `-ExecutionPolicy Bypass` lets PowerShell run the local script file.
>
>     Save this script to `.claude/hooks/plain-display.ps1` in your project:
>
>     ```powershell
>     $batch = [Console]::In.ReadToEnd() | ConvertFrom-Json
>     $text = $batch.delta -replace '\*\*', '' -replace '`', ''
>     @{
>       hookSpecificOutput = @{
>         hookEventName = "MessageDisplay"
>         displayContent = $text
>       }
>     } | ConvertTo-Json
>     ```
>
> Batches with no markdown pass through unchanged. If the script fails, for example because `jq` is missing, Claude Code displays the original text and notes the failure only in [debug output](#debug-hooks), not in the session.