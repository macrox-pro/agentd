---
primary_sources:
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "PreToolUse; PermissionRequest; PostToolUse"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Tool hook events

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex Hooks — PreToolUse

> ### PreToolUse
>
> `PreToolUse` can intercept Bash, file edits performed through `apply_patch`,
> MCP tool calls, and other local function tools. See [Tool
> coverage](#tool-coverage) for the supported paths and exceptions.
>
> `matcher` is applied to `tool_name` and matcher aliases. For file edits through
> `apply_patch`, `matcher` values can use `apply_patch`, `Edit`, or `Write`; hook input
> still reports `tool_name: "apply_patch"`.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field         | Type         | Meaning                                                                                                                          |
> | ------------- | ------------ | -------------------------------------------------------------------------------------------------------------------------------- |
> | `turn_id`     | `string`     | Codex-specific extension. Active Codex turn id                                                                                   |
> | `tool_name`   | `string`     | Canonical hook tool name, such as `Bash`, `apply_patch`, or an MCP name like `mcp__fs__read`                                     |
> | `tool_use_id` | `string`     | Tool-call id for this invocation                                                                                                 |
> | `tool_input`  | `JSON value` | Tool-specific input. `Bash` and `apply_patch` use `tool_input.command`. MCP and other local function tools send their arguments. |
>
> Plain text on `stdout` is ignored.
>
> JSON on `stdout` can use `systemMessage`. To deny a supported tool call, return
> this hook-specific shape:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PreToolUse",
>     "permissionDecision": "deny",
>     "permissionDecisionReason": "Destructive command blocked by hook."
>   }
> }
> ```
>
> Codex also accepts this older block shape:
>
> ```json
> {
>   "decision": "block",
>   "reason": "Destructive command blocked by hook."
> }
> ```
>
> You can also use exit code `2` and write the blocking reason to `stderr`.
>
> To add model-visible context without blocking, return
> `hookSpecificOutput.additionalContext`:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PreToolUse",
>     "additionalContext": "The pending command touches generated files."
>   }
> }
> ```
>
> To rewrite a supported tool call without blocking, return
> `permissionDecision: "allow"` with `updatedInput`:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PreToolUse",
>     "permissionDecision": "allow",
>     "updatedInput": {
>       "command": "echo rewritten"
>     }
>   }
> }
> ```
>
> For Bash commands and `apply_patch`, `updatedInput` must include a string
> `command` field. For MCP and other local function tools, `updatedInput` is the
> replacement arguments object. Return `updatedInput` only with
> `permissionDecision: "allow"`; other `updatedInput` shapes are reported as
> errors.
>
> `permissionDecision: "ask"`, legacy `decision: "approve"`, `continue: false`,
> `stopReason`, and `suppressOutput` are parsed but not supported yet. Codex marks
> the hook run as failed, reports the error, and continues the tool call.

### Source: Codex Hooks — PermissionRequest

> ### PermissionRequest
>
> `PermissionRequest` runs when Codex is about to ask for approval, such as a
> shell escalation or managed-network approval. It can allow the request, deny
> the request, or decline to decide and let the normal approval prompt continue.
> It doesn't run for commands that don't need approval.
>
> `matcher` is applied to `tool_name` and matcher aliases. Current canonical
> values include `Bash`, `apply_patch`, and MCP tool names such as
> `mcp__server__tool`; `apply_patch` also matches `Edit` and `Write`.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field                    | Type             | Meaning                                                                                                        |
> | ------------------------ | ---------------- | -------------------------------------------------------------------------------------------------------------- |
> | `turn_id`                | `string`         | Codex-specific extension. Active Codex turn id                                                                 |
> | `tool_name`              | `string`         | Canonical hook tool name, such as `Bash`, `apply_patch`, or an MCP name like `mcp__fs__read`                   |
> | `tool_input`             | `JSON value`     | Tool-specific input. `Bash` and `apply_patch` use `tool_input.command` while MCP tools send all the arguments. |
> | `tool_input.description` | `string \| null` | Human-readable approval reason, when Codex has one                                                             |
>
> Plain text on `stdout` is ignored.
>
> Some tool inputs may include a human-readable description, but don't rely on a
> `tool_input.description` field for every tool.
>
> To approve the request, return:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PermissionRequest",
>     "decision": {
>       "behavior": "allow"
>     }
>   }
> }
> ```
>
> To deny the request, return:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PermissionRequest",
>     "decision": {
>       "behavior": "deny",
>       "message": "Blocked by repository policy."
>     }
>   }
> }
> ```
>
> If multiple matching hooks return decisions, any `deny` wins. Otherwise, an
> `allow` lets the request proceed without surfacing the approval prompt. If no
> matching hook decides, Codex uses the normal approval flow.
>
> Don't return `updatedInput`, `updatedPermissions`, or `interrupt` for
> `PermissionRequest`; those fields are reserved for future behavior and fail
> closed today.

### Source: Codex Hooks — PostToolUse

> ### PostToolUse
>
> `PostToolUse` runs after supported tools produce output, including Bash,
> `apply_patch`, MCP tool calls, and other local function tools. For Bash, it
> also runs after commands that exit with a non-zero status. It can't undo side
> effects from a tool that already ran. See [Tool coverage](#tool-coverage) for
> the supported paths and exceptions.
>
> `matcher` is applied to `tool_name` and matcher aliases. For file edits through
> `apply_patch`, `matcher` values can use `apply_patch`, `Edit`, or `Write`; hook input
> still reports `tool_name: "apply_patch"`.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field           | Type         | Meaning                                                                                                                          |
> | --------------- | ------------ | -------------------------------------------------------------------------------------------------------------------------------- |
> | `turn_id`       | `string`     | Codex-specific extension. Active Codex turn id                                                                                   |
> | `tool_name`     | `string`     | Canonical hook tool name, such as `Bash`, `apply_patch`, or an MCP name like `mcp__fs__read`                                     |
> | `tool_use_id`   | `string`     | Tool-call id for this invocation                                                                                                 |
> | `tool_input`    | `JSON value` | Tool-specific input. `Bash` and `apply_patch` use `tool_input.command`. MCP and other local function tools send their arguments. |
> | `tool_response` | `JSON value` | Tool-specific output. MCP tools send the MCP call result. Other local function tools normally send their model-facing output.    |
>
> Plain text on `stdout` is ignored.
>
> JSON on `stdout` can use `systemMessage` and this hook-specific shape:
>
> ```json
> {
>   "decision": "block",
>   "reason": "The Bash output needs review before continuing.",
>   "hookSpecificOutput": {
>     "hookEventName": "PostToolUse",
>     "additionalContext": "The command updated generated files."
>   }
> }
> ```
>
> That `additionalContext` text is added as extra developer context.
>
> For this event, `decision: "block"` doesn't undo the completed Bash command.
> Instead, Codex records the feedback, replaces the tool result with that
> feedback, and continues the model from the hook-provided message.
>
> You can also use exit code `2` and write the feedback reason to `stderr`.
>
> To stop normal processing of the original tool result after the command has
> already run, return `continue: false`. Codex will replace the tool result with
> your feedback or stop text and continue from there.
>
> `updatedMCPToolOutput` and `suppressOutput` are parsed but not supported yet.
> Codex marks the hook run as failed, reports the error, and continues normal
> processing of the tool result.
>
> #### Tool calls from code mode
>
> When a model uses code mode to call a tool from JavaScript, hook decisions apply
> to that nested call. `PreToolUse` can stop the tool before it runs or rewrite
> its input. A blocking `PostToolUse` can't undo the tool's side effects, but it
> can keep the original result from reaching the running script.
>
> | Hook result                                                      | What code mode sees                                                                                    |
> | ---------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
> | `PreToolUse` blocks                                              | The tool promise rejects before the tool runs.                                                         |
> | `PreToolUse` returns `updatedInput`                              | The tool runs with the rewritten input and the promise resolves with that result.                      |
> | `PostToolUse` returns `decision: "block"` or exits with code `2` | The tool runs, then the promise rejects with the hook reason.                                          |
> | `PostToolUse` returns `continue: false`                          | Codex uses the hook feedback for the model-visible result, but doesn't reject the nested tool promise. |
