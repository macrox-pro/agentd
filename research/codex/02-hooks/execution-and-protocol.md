---
primary_sources:
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "Matcher patterns; Common input fields; Common output fields; Large hook output"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Matcher patterns and hook protocol

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex Hooks — Matcher patterns / Common fields / Large hook output

> ## Matcher patterns
>
> The `matcher` field is a regex string that filters when hooks fire. Use `"*"`,
> `""`, or omit `matcher` entirely to match every occurrence of a supported
> event.
>
> Only some current Codex events honor `matcher`:
>
> | Event               | What `matcher` filters | Notes                                                        |
> | ------------------- | ---------------------- | ------------------------------------------------------------ |
> | `PermissionRequest` | tool name              | Support includes `Bash`, `apply_patch`\*, and MCP tool names |
> | `PostToolUse`       | tool name              | See [Tool coverage](#tool-coverage)                          |
> | `PostCompact`       | compaction trigger     | Values are `manual` or `auto`                                |
> | `PreCompact`        | compaction trigger     | Values are `manual` or `auto`                                |
> | `PreToolUse`        | tool name              | See [Tool coverage](#tool-coverage)                          |
> | `SessionEnd`        | end reason             | Currently only `other`                                       |
> | `SessionStart`      | start source           | Values are `startup`, `resume`, `clear`, and `compact`       |
> | `SubagentStart`     | subagent type          | Values depend on the subagent that starts                    |
> | `SubagentStop`      | subagent type          | Values depend on the subagent that stops                     |
> | `UserPromptSubmit`  | not supported          | Any configured `matcher` is ignored for this event           |
> | `Stop`              | not supported          | Any configured `matcher` is ignored for this event           |
>
> \*For `apply_patch`, `matcher` values can also use `Edit` or `Write`.
>
> Examples:
>
> - `Bash`
> - `^apply_patch$`
> - `Edit|Write`
> - `mcp__filesystem__read_file`
> - `mcp__filesystem__.*`
> - `startup|resume|clear|compact`
> - `manual|auto`
>
> ### Tool coverage
>
> `PreToolUse` and `PostToolUse` can observe more than shell and MCP calls. Most
> local function tools use the same hook path, so you can match their tool name,
> inspect their JSON arguments, and, for `PreToolUse`, block or rewrite the call.
>
> | Tool path                         | `PreToolUse` | `PostToolUse` | Notes                                                                                                                    |
> | --------------------------------- | ------------ | ------------- | ------------------------------------------------------------------------------------------------------------------------ |
> | Shell commands                    | Yes          | Yes           | Match as `Bash`.                                                                                                         |
> | Unified exec (`exec_command`)     | Yes          | Yes           | Match as `Bash`. A later `write_stdin` poll can deliver the original command's `PostToolUse` when that command finishes. |
> | `apply_patch`                     | Yes          | Yes           | Match as `apply_patch`, `Edit`, or `Write`.                                                                              |
> | MCP tools                         | Yes          | Yes           | Match the MCP tool name, such as `mcp__filesystem__read_file`.                                                           |
> | Other local function tools        | Yes          | Yes           | Match the function tool name, such as `update_plan`. `spawn_agent` also matches `Agent`.                                 |
> | Hosted tools, such as `WebSearch` | No           | No            | These don't use the local function-tool hook path.                                                                       |
>
> `write_stdin` is transport for an existing unified-exec session. It doesn't run
> `PreToolUse` again when it sends input or polls a command that already passed
> `PreToolUse`.
>
> Some specialized tool paths can opt out of the default hook path. Treat tool
> hooks as a useful guardrail, not a complete enforcement boundary.
>
> ## Common input fields
>
> Every command hook receives one JSON object on `stdin`.
>
> These are the shared fields you will usually use:
>
> | Field             | Type             | Meaning                                                             |
> | ----------------- | ---------------- | ------------------------------------------------------------------- |
> | `session_id`      | `string`         | Current Codex session id. Subagent hooks use the parent session id. |
> | `transcript_path` | `string \| null` | Path to the session transcript file, if any                         |
> | `cwd`             | `string`         | Working directory for the session                                   |
> | `hook_event_name` | `string`         | Current hook event name                                             |
> | `model`           | `string`         | Codex-specific extension. Active model slug                         |
>
> Turn-scoped hooks list `turn_id` as a Codex-specific extension in their
> event-specific tables.
>
> `SessionStart`, `PreToolUse`, `PermissionRequest`, `PostToolUse`,
> `UserPromptSubmit`, `SubagentStart`, `SubagentStop`, and `Stop` also include
> `permission_mode`, which describes the current permission mode as `default`,
> `acceptEdits`, `plan`, `dontAsk`, or `bypassPermissions`.
>
> `transcript_path` points to a chat transcript for convenience, but the
> transcript format isn't a stable interface for hooks and may change over time.
>
> If you need the full wire format, see [Schemas](#schemas).
>
> ## Common output fields
>
> `SessionStart`, `PreCompact`, `PostCompact`, `UserPromptSubmit`,
> `SubagentStop`, and `Stop` support these shared JSON fields. `SubagentStart`
> accepts the same shape for `systemMessage` and hook-specific context, but
> `continue: false` doesn't stop the subagent:
>
> ```json
> {
>   "continue": true,
>   "stopReason": "optional",
>   "systemMessage": "optional",
>   "suppressOutput": false
> }
> ```
>
> | Field            | Effect                                          |
> | ---------------- | ----------------------------------------------- |
> | `continue`       | If `false`, marks that hook run as stopped      |
> | `stopReason`     | Recorded as the reason for stopping             |
> | `systemMessage`  | Surfaced as a warning in the UI or event stream |
> | `suppressOutput` | Parsed today but not yet implemented            |
>
> Exit `0` with no output is treated as success and Codex continues.
>
> `PreToolUse` and `PermissionRequest` support `systemMessage`, but `continue`,
> `stopReason`, and `suppressOutput` aren't currently supported for those events.
> If a `PreToolUse` hook returns one of those unsupported fields, Codex marks
> that hook run as failed, reports the error, and continues the tool call.
>
> `PostToolUse` supports `systemMessage`, `continue: false`, and `stopReason`.
> `suppressOutput` is parsed but not currently supported for that event.
>
> ### Large hook output
>
> By default, Codex limits each model-visible hook-output message to roughly
> 2,500 tokens. If a hook returns more, Codex saves the full text under
> `<temp_dir>/hook_outputs/<session_id>/<uuid>.txt` and gives the model a
> head-and-tail preview with the saved-file path. This behavior is called
> **spilling**: Codex stores oversized output on disk and replaces it with a
> shorter, model-visible preview. If the file can't be written, the model still
> receives a truncated preview.
>
> Keep hook and plugin context concise. Context from multiple hooks and plugins
>   adds up and can degrade model performance. Raising `additionalContextLimit`
>   increases that risk. Avoid setting the limit to `0` unless the hook enforces a
>   strict output cap; otherwise, a single hook can consume the entire context
>   window.
>
> For any command hook that returns `additionalContext`, set
> `additionalContextLimit` on the handler to customize the approximate token
> threshold:
>
> ```json
> {
>   "type": "command",
>   "command": "python3 ~/.codex/hooks/session_start.py",
>   "additionalContextLimit": 5000
> }
> ```
>
> Omit `additionalContextLimit` to use the default `2500`-token threshold. Use a
> positive integer to select a different threshold, or `0` to pass the handler's
> complete additional context directly to the model. Codex evaluates each
> matching handler independently. For events that can't produce additional
> context, Codex ignores `additionalContextLimit` and reports a configuration
> warning.
>
> The setting applies only to `additionalContext`. Tool feedback and continuation
> prompts keep the default limit.
>
> Because oversized output can be written to disk, avoid returning secrets or
> other sensitive data in hook output.
