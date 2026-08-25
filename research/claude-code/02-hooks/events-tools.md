---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Hook events — tools and permissions"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook events — tools and permissions

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — Hook events

> ## Hook events
>
> Each event corresponds to a point in Claude Code's lifecycle where hooks can run. The sections below are ordered to match the lifecycle: from session setup through the agentic loop to session end. Each section describes when the event fires, what matchers it supports, the JSON input it receives, and how to control behavior through output.

### Source: Hooks reference — PreToolUse

> ### PreToolUse
>
> Runs after Claude creates tool parameters and before processing the tool call. Matches on tool name: `Bash`, `PowerShell`, `Edit`, `Write`, `Read`, `Glob`, `Grep`, `Agent`, `WebFetch`, `WebSearch`, `AskUserQuestion`, `ExitPlanMode`, and any [MCP tool names](#match-mcp-tools).
>
> To run a hook when a specific file changes on disk, whatever wrote it, use [FileChanged](#filechanged) instead of matching file-editing tools by name. Unlike PreToolUse, Claude Code runs FileChanged hooks after the change, and they have no decision control, so they can't block the write.
>
>
>   PreToolUse runs only when Claude calls a tool. Files you [reference with `@` in your prompt](/docs/en/common-workflows#reference-files-and-directories) are added without any tool call: Claude Code inserts their contents while building the prompt, so no PreToolUse hook fires for them, including hooks matching `Read`. To block specific paths from `@` references, use a [`Read` deny rule](/docs/en/permissions#read-and-edit) instead.
>
>   PreToolUse also doesn't fire for [`EndConversation`](/docs/en/tools-reference#endconversation-tool-behavior).
>
>
> Use [PreToolUse decision control](#pretooluse-decision-control) to allow, deny, ask, or defer the tool call.
>
> An [Agent SDK callback hook](/docs/en/agent-sdk/hooks) on `PreToolUse` that exceeds its timeout blocks the tool call, and Claude receives an error result naming the timeout. An explicit deny returned by another hook still takes precedence.
>
> #### PreToolUse input
>
> In addition to the [common input fields](#common-input-fields), PreToolUse hooks receive `tool_name`, `tool_input`, and `tool_use_id`.
>
> For the file tools `Write`, `Edit`, and `Read`, `tool_input.file_path` is always absolute:
>
> * Claude Code expands `~` and relative paths before hooks run, so a hook that matches on paths can't be bypassed via `~` or a relative spelling of the same path
> * On Windows, the path arrives with backslash separators, even when your hook runs under Git Bash where `$PWD` looks like `/c/project`
> * A comparison written with forward slashes, such as a `/src/` check, never matches a backslash path, and the tool call proceeds as if the hook had nothing to block
> * Normalize separators before comparing: `FILE_PATH="${FILE_PATH//\\//}"` in Bash, or `file_path.replace("\\", "/")` in Python, then match a path segment such as `/src/` rather than anchoring with `^`, since the path is absolute
>
> A `Write` call on Windows delivers:
>
> ```json
> {
>   "hook_event_name": "PreToolUse",
>   "tool_name": "Write",
>   "tool_input": {
>     "file_path": "C:\\project\\src\\index.ts",
>     "content": "..."
>   },
>   ...
> }
> ```
>
> The `tool_input` fields depend on the tool:
>
> ##### Bash
>
> Executes shell commands.
>
> | Field               | Type    | Example            | Description                                                                                                                                          |
> | :------------------ | :------ | :----------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `command`           | string  | `"npm test"`       | The shell command to execute                                                                                                                         |
> | `description`       | string  | `"Run test suite"` | Optional description of what the command does                                                                                                        |
> | `timeout`           | number  | `120000`           | Optional timeout in milliseconds. Values above the [maximum](/docs/en/tools-reference#bash-tool-behavior) are reduced to the maximum rather than rejected |
> | `run_in_background` | boolean | `false`            | Whether to run the command in background                                                                                                             |
>
>
> ##### PowerShell
>
> Executes PowerShell commands. See the [PowerShell tool](/docs/en/tools-reference#powershell-tool) for availability by platform.
>
> The fields match the Bash tool, with the command string in `command`:
>
> | Field               | Type    | Example                    | Description                                   |
> | :------------------ | :------ | :------------------------- | :-------------------------------------------- |
> | `command`           | string  | `"Get-ChildItem -Recurse"` | The PowerShell command to execute             |
> | `description`       | string  | `"List files recursively"` | Optional description of what the command does |
> | `timeout`           | number  | `120000`                   | Optional timeout in milliseconds              |
> | `run_in_background` | boolean | `false`                    | Whether to run the command in background      |
>
> Match `Bash|PowerShell` in hooks that inspect shell commands, so they cover both tools:
>
> * On Windows, wherever the PowerShell tool is enabled, Claude treats PowerShell as the primary shell and routes shell commands through it.
> * On Windows without Git Bash, the tool is enabled automatically and Claude Code doesn't register the Bash tool at all.
> * A hook that matches only `Bash` never fires there.
>
> ##### Write
>
> Creates or overwrites a file.
>
> | Field       | Type   | Example               | Description                        |
> | :---------- | :----- | :-------------------- | :--------------------------------- |
> | `file_path` | string | `"/path/to/file.txt"` | Absolute path to the file to write |
> | `content`   | string | `"file content"`      | Content to write to the file       |
>
> ##### Edit
>
> Replaces a string in an existing file.
>
> | Field         | Type    | Example               | Description                        |
> | :------------ | :------ | :-------------------- | :--------------------------------- |
> | `file_path`   | string  | `"/path/to/file.txt"` | Absolute path to the file to edit  |
> | `old_string`  | string  | `"original text"`     | Text to find and replace           |
> | `new_string`  | string  | `"replacement text"`  | Replacement text                   |
> | `replace_all` | boolean | `false`               | Whether to replace all occurrences |
>
> ##### Read
>
> Reads file contents.
>
> | Field       | Type   | Example               | Description                                |
> | :---------- | :----- | :-------------------- | :----------------------------------------- |
> | `file_path` | string | `"/path/to/file.txt"` | Absolute path to the file to read          |
> | `offset`    | number | `10`                  | Optional line number to start reading from |
> | `limit`     | number | `50`                  | Optional number of lines to read           |
>
> ##### Glob
>
> Finds files matching a glob pattern.
>
> | Field     | Type   | Example          | Description                                                            |
> | :-------- | :----- | :--------------- | :--------------------------------------------------------------------- |
> | `pattern` | string | `"**/*.ts"`      | Glob pattern to match files against                                    |
> | `path`    | string | `"/path/to/dir"` | Optional directory to search in. Defaults to current working directory |
>
> ##### Grep
>
> Searches file contents with regular expressions.
>
> | Field         | Type    | Example          | Description                                                                           |
> | :------------ | :------ | :--------------- | :------------------------------------------------------------------------------------ |
> | `pattern`     | string  | `"TODO.*fix"`    | Regular expression pattern to search for                                              |
> | `path`        | string  | `"/path/to/dir"` | Optional file or directory to search in                                               |
> | `glob`        | string  | `"*.ts"`         | Optional glob pattern to filter files                                                 |
> | `output_mode` | string  | `"content"`      | `"content"`, `"files_with_matches"`, or `"count"`. Defaults to `"files_with_matches"` |
> | `-i`          | boolean | `true`           | Case insensitive search                                                               |
> | `multiline`   | boolean | `false`          | Enable multiline matching                                                             |
>
> ##### WebFetch
>
> Fetches and processes web content.
>
> | Field    | Type   | Example                       | Description                          |
> | :------- | :----- | :---------------------------- | :----------------------------------- |
> | `url`    | string | `"https://example.com/api"`   | URL to fetch content from            |
> | `prompt` | string | `"Extract the API endpoints"` | Prompt to run on the fetched content |
>
> ##### WebSearch
>
> Searches the web.
>
> | Field             | Type   | Example                        | Description                                       |
> | :---------------- | :----- | :----------------------------- | :------------------------------------------------ |
> | `query`           | string | `"react hooks best practices"` | Search query                                      |
> | `allowed_domains` | array  | `["docs.example.com"]`         | Optional: only include results from these domains |
> | `blocked_domains` | array  | `["spam.example.com"]`         | Optional: exclude results from these domains      |
>
> ##### Agent
>
> Spawns a [subagent](/docs/en/sub-agents).
>
> | Field           | Type   | Example                    | Description                                  |
> | :-------------- | :----- | :------------------------- | :------------------------------------------- |
> | `prompt`        | string | `"Find all API endpoints"` | The task for the agent to perform            |
> | `description`   | string | `"Find API endpoints"`     | Short description of the task                |
> | `subagent_type` | string | `"Explore"`                | Type of specialized agent to use             |
> | `model`         | string | `"sonnet"`                 | Optional model alias to override the default |
>
> When a foreground Agent call completes, your [PostToolUse hook](#posttooluse) receives the subagent's final text and run telemetry in `tool_response`. Read these fields to inspect the run; for token and cost rollups across subagents, use the [token and cost counters](/docs/en/monitoring-usage#token-counter) filtered to `query_source` `"subagent"`, since `totalTokens` and `usage` cover the final request only:
>
> | Field               | Type   | Example                                               | Description                                                                                                                                                                                                         |
> | :------------------ | :----- | :---------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `status`            | string | `"completed"`                                         | `"completed"` for foreground subagents, `"async_launched"` for background subagents. As of v2.1.198, subagents run in the background by default, so an omitted `run_in_background` also produces `"async_launched"` |
> | `agentId`           | string | `"a4d2c8f1e0b3a297"`                                  | Identifier for the subagent run                                                                                                                                                                                     |
> | `content`           | array  | `[{"type": "text", "text": "Found 12 endpoints..."}]` | The subagent's final text blocks                                                                                                                                                                                    |
> | `resolvedModel`     | string | `"claude-sonnet-4-5"`                                 | Model the subagent started on, which may differ from the requested model. Requires Claude Code v2.1.174 or later                                                                                                    |
> | `modelsUsed`        | array  | `["claude-sonnet-4-5", "claude-haiku-4-5"]`           | Models used in order, with consecutive repeats collapsed; set only when the model was swapped mid-run. Requires Claude Code v2.1.212 or later                                                                       |
> | `totalTokens`       | number | `12450`                                               | Token count from the subagent's final API request: input, output, and cache tokens combined. This isn't a total across the whole run                                                                                |
> | `totalDurationMs`   | number | `48211`                                               | Wall-clock duration of the subagent run                                                                                                                                                                             |
> | `totalToolUseCount` | number | `7`                                                   | Count of tool calls the subagent made                                                                                                                                                                               |
> | `usage`             | object | `{"input_tokens": 8320, ...}`                         | Per-type token breakdown of the final API request: `input_tokens`, `output_tokens`, `cache_creation_input_tokens`, `cache_read_input_tokens`                                                                        |
>
> For background subagents, the tool returns when the task moves to the background, so `tool_response` carries no usage fields: a background launch returns immediately, and a foreground task that Claude Code backgrounds mid-run returns at that transition. It has `status: "async_launched"`, `agentId`, `description`, `prompt`, `outputFile`, and `resolvedModel`.
>
> On a `completed` response, `resolvedModel` names the model the subagent started on, which can differ from the `model` value in `tool_input`, such as when `availableModels` or another override applies. It requires Claude Code v2.1.174 or later. On an `async_launched` response, `resolvedModel` names the model in use when the agent moved to the background, so a swap that happened before backgrounding is reflected there. `modelsUsed` and the backgrounding-time `resolvedModel` behavior require Claude Code v2.1.212 or later.
>
>
> ##### AskUserQuestion
>
> Asks the user one to four multiple-choice questions.
>
> | Field       | Type   | Example                                                                                                            | Description                                                                                                                                                                                     |
> | :---------- | :----- | :----------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `questions` | array  | `[{"question": "Which framework?", "header": "Framework", "options": [{"label": "React"}], "multiSelect": false}]` | Questions to present, each with a `question` string, short `header`, `options` array, and optional `multiSelect` flag                                                                           |
> | `answers`   | object | `{"Which framework?": "React"}`                                                                                    | Optional. Maps question text to the selected option label. Multi-select answers join labels with commas. Claude doesn't set this field; supply it via `updatedInput` to answer programmatically |
>
> ##### ExitPlanMode
>
> Presents a plan and asks the user to approve it before Claude leaves [plan mode](/docs/en/permission-modes#analyze-before-you-edit-with-plan-mode). Claude writes the plan to a file on disk before calling the tool, so the literal `tool_input` from the model is typically empty. Claude Code injects the plan content and file path before passing the input to hooks.
>
> | Field            | Type   | Example                                     | Description                                                                                                                                           |
> | :--------------- | :----- | :------------------------------------------ | :---------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `plan`           | string | `"## Refactor auth\n1. Extract..."`         | Plan content in Markdown. Injected from the plan file on disk                                                                                         |
> | `planFilePath`   | string | `"/Users/.../plans/refactor-auth.md"`       | Path to the plan file. Injected                                                                                                                       |
> | `allowedPrompts` | array  | `[{"tool": "Bash", "prompt": "run tests"}]` | Deprecated. Claude Code accepts the field but ignores it. Before v2.1.205, it carried prompt-based permissions Claude requested to implement the plan |
>
> In `PostToolUse`, `tool_response` is an object with `plan` and `filePath` fields holding the approved plan, plus internal status flags. Read `tool_response.plan` for the plan content rather than re-reading the file from disk.
>
> #### PreToolUse decision control
>
> `PreToolUse` hooks can control whether a tool call proceeds. Unlike other hooks that use a top-level `decision` field, PreToolUse returns its decision inside a `hookSpecificOutput` object. This gives it richer control: four outcomes (allow, deny, ask, or defer) plus the ability to modify tool input before execution.
>
> | Field                      | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
> | :------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `permissionDecision`       | `"allow"` skips the permission prompt, except for the [actions no mode auto-approves](/docs/en/permission-modes#actions-no-mode-auto-approves) and for `AskUserQuestion` and `ExitPlanMode`, which need [`updatedInput` paired with it](#allow-with-updatedinput). `"deny"` prevents the tool call. `"ask"` prompts the user to confirm. `"defer"` exits gracefully so the tool can be resumed later. [Deny and ask rules](/docs/en/permissions#manage-permissions) are still evaluated regardless of what the hook returns |
> | `permissionDecisionReason` | For `"allow"` and `"ask"`, shown to the user but not Claude. For `"deny"`, shown to Claude. For `"defer"`, ignored                                                                                                                                                                                                                                                                                                                                                                                                |
> | `updatedInput`             | Modifies the tool's input parameters before execution. Replaces the entire input object, so include unchanged fields alongside modified ones. Combine with `"allow"` to auto-approve, or `"ask"` to show the modified input to the user. For `"defer"`, ignored                                                                                                                                                                                                                                                   |
> | `additionalContext`        | String added to Claude's context alongside the tool result. Ignored when `permissionDecision` is `"defer"`. See [Add context for Claude](#add-context-for-claude)                                                                                                                                                                                                                                                                                                                                                 |
>
> When multiple PreToolUse hooks return different decisions, precedence is `deny` > `defer` > `ask` > `allow`.
>
> A hook that blocks by exiting 2 routes the same way as `"deny"`: Claude sees the stderr message as the denial reason.
>
> When a hook returns `"ask"`, the permission prompt displayed to the user includes a label identifying where the hook came from: for example, `[User]`, `[Project]`, `[Plugin]`, or `[Local]`. This helps users understand which configuration source is requesting confirmation.
>
> A hook's `"ask"` also forces a permission prompt in [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode): the classifier can still deny the tool call, but it can't approve the call silently. Before v2.1.211, the classifier could approve a Bash command running outside the [sandbox](/docs/en/sandboxing) without showing the prompt the hook requested; the classifier still applied its own safety rules to that command, and a hook `"deny"` was always honored.
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PreToolUse",
>     "permissionDecision": "allow",
>     "permissionDecisionReason": "My reason here",
>     "updatedInput": {
>       "field_to_modify": "new value"
>     },
>     "additionalContext": "Current environment: production. Proceed with caution."
>   }
> }
> ```
>
>
> `AskUserQuestion` and `ExitPlanMode` require user interaction and normally block in [non-interactive mode](/docs/en/headless) with the `-p` flag. Returning `permissionDecision: "allow"` together with `updatedInput` satisfies that requirement: the hook reads the tool's input from stdin, collects the answer through your own UI, and returns it in `updatedInput` so the tool runs without prompting. Returning `"allow"` alone is not sufficient for these tools. For `AskUserQuestion`, echo back the original `questions` array and add an [`answers`](#askuserquestion) object mapping each question's text to the chosen answer.
>
> As of v2.1.199, an MCP tool whose server marks it with [`_meta["anthropic/requiresUserInteraction"]`](/docs/en/mcp#require-approval-for-a-specific-tool) is stricter: a hook can't skip its approval prompt with `"allow"`, with or without `updatedInput`, because Claude Code can't confirm the hook collected the interaction the tool needs.
>
>
>   PreToolUse previously used top-level `decision` and `reason` fields, but these are deprecated for this event. Use `hookSpecificOutput.permissionDecision` and `hookSpecificOutput.permissionDecisionReason` instead. The deprecated values `"approve"` and `"block"` map to `"allow"` and `"deny"` respectively. Other events like PostToolUse and Stop continue to use top-level `decision` and `reason` as their current format.
>
>
> #### Defer a tool call for later
>
> `"defer"` is for integrations that run `claude -p` as a subprocess and read its JSON output, such as an Agent SDK app or a custom UI built on top of Claude Code. It lets that calling process pause Claude at a tool call, collect input through its own interface, and resume where it left off. Claude Code honors this value only in [non-interactive mode](/docs/en/headless) with the `-p` flag. In interactive sessions it logs a warning and ignores the hook result.
>
> The `AskUserQuestion` tool is the typical case: Claude wants to ask the user something, but there is no terminal to answer in. The round trip works like this:
>
> 1. Claude calls `AskUserQuestion`. The `PreToolUse` hook fires.
> 2. The hook returns `permissionDecision: "defer"`. The tool doesn't execute. The process exits with `stop_reason: "tool_deferred"` and the pending tool call preserved in the transcript.
> 3. The calling process reads `deferred_tool_use` from the SDK result, surfaces the question in its own UI, and waits for an answer.
> 4. The calling process runs `claude -p --resume <session-id>`. The same tool call fires `PreToolUse` again.
> 5. The hook returns `permissionDecision: "allow"` with the answer in `updatedInput`. The tool executes and Claude continues.
>
> The `deferred_tool_use` field carries the tool's `id`, `name`, and `input`. The `input` is the parameters Claude generated for the tool call, captured before execution:
>
> ```json
> {
>   "type": "result",
>   "subtype": "success",
>   "stop_reason": "tool_deferred",
>   "session_id": "abc123",
>   "deferred_tool_use": {
>     "id": "toolu_01abc",
>     "name": "AskUserQuestion",
>     "input": { "questions": [{ "question": "Which framework?", "header": "Framework", "options": [{"label": "React"}, {"label": "Vue"}], "multiSelect": false }] }
>   }
> }
> ```
>
> There is no timeout or retry limit. The session remains on disk until you resume it, subject to the [`cleanupPeriodDays`](/docs/en/settings-reference#cleanupperioddays) retention sweep, which deletes session files after 30 days by default, following the [retention sweep rules](/docs/en/claude-directory#cleaned-up-automatically). If the answer is not ready when you resume, the hook can return `"defer"` again and the process exits the same way. The calling process controls when to break the loop by eventually returning `"allow"` or `"deny"` from the hook.
>
> `"defer"` only works when Claude makes a single tool call in the turn. If Claude makes several tool calls at once, `"defer"` is ignored with a warning and the tool proceeds through the normal permission flow. The constraint exists because resume can only re-run one tool: there is no way to defer one call from a batch without leaving the others unresolved.
>
> If the deferred tool is no longer available when you resume, the process exits with `stop_reason: "tool_deferred_unavailable"` and `is_error: true` before the hook fires. This happens when an MCP server that provided the tool is not connected for the resumed session. The `deferred_tool_use` payload is still included so you can identify which tool went missing.
>
>
>   `--resume` restores the permission mode that was active when the tool was deferred, so you don't need to pass `--permission-mode` again. The exceptions are `plan` and `bypassPermissions`, which are never carried over, and `auto`, which is restored only when your account still meets the [auto mode requirements](/docs/en/permission-modes#eliminate-prompts-with-auto-mode). Passing `--permission-mode` explicitly on resume overrides the restored value.

### Source: Hooks reference — PermissionRequest

> ### PermissionRequest
>
> Runs when Claude Code is about to ask you for permission. In sessions that can't show a prompt, such as background subagents in [non-interactive mode](/docs/en/headless), Claude Code still runs these hooks, and if no hook returns a decision, it denies the tool call.
> Use [PermissionRequest decision control](#permissionrequest-decision-control) to allow or deny on behalf of the user.
>
> Use this event when you need a signal the moment Claude asks for permission. The [Notification](#notification) event's `permission_prompt` type reaches you only after the prompt has waited about six seconds.
>
> Matches on tool name, same values as PreToolUse.
>
> #### PermissionRequest input
>
> PermissionRequest hooks receive `tool_name` and `tool_input` fields like PreToolUse hooks, but without `tool_use_id`. An optional `permission_suggestions` array contains the "always allow" options the user would normally see in the permission dialog.
>
> PreToolUse hooks run before every tool call, whether or not it needs permission. PermissionRequest hooks run only when Claude Code is about to ask you for permission, or when it would otherwise auto-deny a call that can't prompt. Neither event fires for [`EndConversation`](/docs/en/tools-reference#endconversation-tool-behavior).
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "PermissionRequest",
>   "tool_name": "Bash",
>   "tool_input": {
>     "command": "rm -rf node_modules",
>     "description": "Remove node_modules directory"
>   },
>   "permission_suggestions": [
>     {
>       "type": "addRules",
>       "rules": [{ "toolName": "Bash", "ruleContent": "rm -rf node_modules" }],
>       "behavior": "allow",
>       "destination": "localSettings"
>     }
>   ]
> }
> ```
>
> #### PermissionRequest decision control
>
> `PermissionRequest` hooks can allow or deny permission requests. In addition to the [JSON output fields](#json-output) available to all hooks, your hook script can return a `decision` object with these event-specific fields:
>
> | Field                | Description                                                                                                                                                                                                                     |
> | :------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `behavior`           | `"allow"` grants the permission, `"deny"` denies it. [Deny and ask rules](/docs/en/permissions#manage-permissions) are still evaluated, so a hook returning `"allow"` doesn't override a matching deny rule                          |
> | `updatedInput`       | For `"allow"` only: modifies the tool's input parameters before execution. Replaces the entire input object, so include unchanged fields alongside modified ones. The modified input is re-evaluated against deny and ask rules |
> | `updatedPermissions` | For `"allow"` only: array of [permission update entries](#permission-update-entries) to apply, such as adding an allow rule or changing the session permission mode                                                             |
> | `message`            | For `"deny"` only: tells Claude why the permission was denied                                                                                                                                                                   |
> | `interrupt`          | For `"deny"` only: if `true`, stops Claude                                                                                                                                                                                      |
>
> A hook that exits 2 without a `decision` object leaves the permission flow unchanged, and its stderr is discarded. Only the `decision` object can grant or deny the request.
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PermissionRequest",
>     "decision": {
>       "behavior": "allow",
>       "updatedInput": {
>         "command": "npm run lint"
>       }
>     }
>   }
> }
> ```
>
> #### Permission update entries
>
> The `updatedPermissions` output field and the [`permission_suggestions` input field](#permissionrequest-input) both use the same array of entry objects. Each entry has a `type` that determines its other fields, and a `destination` that controls where the change is written.
>
> | `type`              | Fields                             | Effect                                                                                                                                                                                                                   |
> | :------------------ | :--------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `addRules`          | `rules`, `behavior`, `destination` | Adds permission rules. `rules` is an array of `{toolName, ruleContent?}` objects. Omit `ruleContent` to match the whole tool. `behavior` is `"allow"`, `"deny"`, or `"ask"`                                              |
> | `replaceRules`      | `rules`, `behavior`, `destination` | Replaces all rules of the given `behavior` at the `destination` with the provided `rules`                                                                                                                                |
> | `removeRules`       | `rules`, `behavior`, `destination` | Removes matching rules of the given `behavior`                                                                                                                                                                           |
> | `setMode`           | `mode`, `destination`              | Changes the permission mode. Valid modes are `default`, `auto`, `acceptEdits`, `dontAsk`, `bypassPermissions`, `plan`, and `manual` as an alias for `default`. The `manual` alias requires Claude Code v2.1.200 or later |
> | `addDirectories`    | `directories`, `destination`       | Adds working directories. `directories` is an array of path strings                                                                                                                                                      |
> | `removeDirectories` | `directories`, `destination`       | Removes working directories                                                                                                                                                                                              |
>
>
>   `setMode` with `bypassPermissions` only takes effect if the session was launched with bypass mode already available: `--dangerously-skip-permissions`, `--permission-mode bypassPermissions`, `--allow-dangerously-skip-permissions`, or `permissions.defaultMode: "bypassPermissions"` in settings, and the mode is not disabled by [`permissions.disableBypassPermissionsMode`](/docs/en/permissions#managed-settings). Otherwise the update is a no-op. `bypassPermissions` is never persisted as `defaultMode` regardless of `destination`.
>
>
> The `destination` field on every entry determines whether the change stays in memory or persists to a settings file.
>
> | `destination`     | Writes to                                       |
> | :---------------- | :---------------------------------------------- |
> | `session`         | in-memory only, discarded when the session ends |
> | `localSettings`   | `.claude/settings.local.json`                   |
> | `projectSettings` | `.claude/settings.json`                         |
> | `userSettings`    | `~/.claude/settings.json`                       |
>
> A hook can echo one of the `permission_suggestions` it received as its own `updatedPermissions` output, which is equivalent to the user selecting that "always allow" option in the dialog.

### Source: Hooks reference — PermissionDenied

> ### PermissionDenied
>
> Runs when [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) denies a tool call, including when it denies without a classifier verdict because [a safety check separate from auto mode refused the classifier's own request](/docs/en/errors#auto-mode-cannot-determine-the-safety-of-an-action) or its response didn't parse. This hook only fires in auto mode: it doesn't run when you manually deny a permission dialog, when a `PreToolUse` hook blocks a call, or when a `deny` rule matches. Use it to log denials, adjust configuration, or tell the model it may retry the tool call.
>
> Matches on tool name, same values as PreToolUse.
>
> #### PermissionDenied input
>
> In addition to the [common input fields](#common-input-fields), PermissionDenied hooks receive `tool_name`, `tool_input`, `tool_use_id`, and `reason`.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "auto",
>   "hook_event_name": "PermissionDenied",
>   "tool_name": "Bash",
>   "tool_input": {
>     "command": "rm -rf /tmp/build",
>     "description": "Clean build directory"
>   },
>   "tool_use_id": "toolu_01ABC123...",
>   "reason": "Blocked by classifier"
> }
> ```
>
> | Field    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
> | :------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `reason` | The denial reason: the fixed text `Blocked by classifier` in most sessions, or the classifier's written explanation when the session's classifier model provides one. For a denial where a safety check separate from auto mode refused the classifier's request or its response didn't parse, the reason starts with `Auto mode could not evaluate this action and is blocking it for safety`. For a denial because the classifier model was unavailable, the reason is the fixed text `Classifier unavailable`. See [Review denials](/docs/en/auto-mode-config#review-denials) |
>
> #### PermissionDenied decision control
>
> PermissionDenied hooks can tell the model it may retry the denied tool call. Return a JSON object with `hookSpecificOutput.retry` set to `true`:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PermissionDenied",
>     "retry": true
>   }
> }
> ```
>
> When `retry` is `true`, Claude Code adds a message to the conversation telling the model it may retry the tool call. Claude Code doesn't reverse the denial itself. If your hook doesn't return JSON, or returns `retry: false`, the denial stands and the model receives the original rejection message.
>
> Claude Code ignores `retry: true` when the classifier produced [no verdict on the action](/docs/en/errors#auto-mode-cannot-determine-the-safety-of-an-action): its response didn't parse, or a safety check separate from auto mode refused the classifier's own request. For those denials, Claude Code already tells the model in the rejection message whether to retry later or move on.

### Source: Hooks reference — PostToolUse

> ### PostToolUse
>
> Runs immediately after a tool completes successfully.
>
> Matches on tool name, same values as PreToolUse.
>
> Match more broadly when the tool name isn't the right filter:
>
> * To run a hook after any tool completes successfully, omit the `matcher` or set it to `"*"`. Your hook can then discover what changed itself, for example by running `git status --porcelain`, which also lists untracked files that `git diff` misses. For tool calls that fail, add the same hook under [PostToolUseFailure](#posttoolusefailure).
> * To run a hook when a specific file changes on disk, whatever wrote it, use [FileChanged](#filechanged). Claude Code doesn't run a `PostToolUse` hook matching `Edit|Write` when a `Bash` command or a process outside Claude Code rewrites the same file.
>
> #### PostToolUse input
>
> `PostToolUse` hooks fire after a tool has already executed successfully. The input includes both `tool_input`, the arguments sent to the tool, and `tool_response`, the result it returned. The exact schema for both depends on the tool. File-tool `tool_input` paths arrive in the same format as for [PreToolUse](#pretooluse-input): always absolute, with the platform's native separators, so backslashes on Windows.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "PostToolUse",
>   "tool_name": "Write",
>   "tool_input": {
>     "file_path": "/path/to/file.txt",
>     "content": "file content"
>   },
>   "tool_response": {
>     "filePath": "/path/to/file.txt",
>     "success": true
>   },
>   "tool_use_id": "toolu_01ABC123...",
>   "duration_ms": 12
> }
> ```
>
> | Field         | Description                                                                                                   |
> | :------------ | :------------------------------------------------------------------------------------------------------------ |
> | `duration_ms` | Optional. Tool execution time in milliseconds. Excludes time spent in permission prompts and PreToolUse hooks |
>
> #### PostToolUse decision control
>
> `PostToolUse` hooks can provide feedback to Claude after tool execution. In addition to the [JSON output fields](#json-output) available to all hooks, your hook script can return these event-specific fields:
>
> | Field                  | Description                                                                                                                                                                                                                                                                                     |
> | :--------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `decision`             | `"block"` adds the `reason` next to the tool result. Claude still sees the original output; to replace it, use `updatedToolOutput`                                                                                                                                                              |
> | `reason`               | Explanation shown to Claude when `decision` is `"block"`                                                                                                                                                                                                                                        |
> | `additionalContext`    | String added to Claude's context alongside the tool result. See [Add context for Claude](#add-context-for-claude)                                                                                                                                                                               |
> | `classifierContext`    | Short note about this call's result for the [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) classifier rather than for Claude. See [Annotate a result for the auto mode classifier](#annotate-a-result-for-the-auto-mode-classifier). Requires Claude Code v2.1.236 or later |
> | `updatedToolOutput`    | Replaces the tool's output with the provided value before it is sent to Claude. The value must match the tool's output shape                                                                                                                                                                    |
> | `updatedMCPToolOutput` | Replaces the output for [MCP tools](#match-mcp-tools) only. Prefer `updatedToolOutput`, which works for all tools                                                                                                                                                                               |
>
> The example below replaces the output of a `Bash` call. The replacement value matches the `Bash` tool's output shape:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PostToolUse",
>     "additionalContext": "Additional information for Claude",
>     "updatedToolOutput": {
>       "stdout": "[redacted]",
>       "stderr": "",
>       "interrupted": false,
>       "isImage": false
>     }
>   }
> }
> ```
>
>
>   `updatedToolOutput` only changes what Claude sees. The tool has already run by the time the hook fires, so any files written, commands executed, or network requests sent have already taken effect. Telemetry such as OpenTelemetry tool spans and analytics events also captures the original output before the hook runs. To prevent or modify a tool call before it runs, use a [PreToolUse](#pretooluse) hook instead.
>
>   The replacement value must match the tool's output shape. Built-in tools return structured objects rather than plain strings. For example, `Bash` returns an object with `stdout`, `stderr`, `interrupted`, and `isImage` fields. For built-in tools, a value that doesn't match the tool's output schema is ignored and the original output is used. MCP tool output is passed through without schema validation. Stripping error details that Claude needs can cause it to proceed on a false assumption.
>
>
> #### Annotate a result for the auto mode classifier
>
> Return `classifierContext` to send a short note about the tool call's result to the [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) classifier rather than to Claude. The classifier [never receives tool results themselves](/docs/en/permission-modes#how-the-classifier-evaluates-actions), so this field is the supported way to tell it something about what a call returned before it reviews later actions. The field requires Claude Code v2.1.236 or later.
>
> The example below tells the classifier where a query's output came from:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PostToolUse",
>     "classifierContext": "This query ran against the staging database, not production."
>   }
> }
> ```
>
> How much weight the classifier gives the note depends on where you configured the hook:
>
> * **Hooks configured in Claude Code**: for hooks from settings files, plugins, skills, and agent frontmatter, the classifier treats the note as unverified, application-provided context. The note never establishes user intent, and if it claims you approved or requested something, the classifier checks that claim against your own messages in the conversation
> * **In-process Agent SDK callbacks**: when an application embedding Claude Code registers the hook as a [TypeScript SDK callback](/docs/en/agent-sdk/hooks) and returns the note during the live session, the classifier may weigh a user statement relayed in the note as user intent. Such a statement can satisfy a consent requirement the classifier would accept from a message you send, but it never lifts a block that your own message couldn't lift either. After a session resumes, Claude Code treats restored notes as unverified context. When hooks from both groups annotate the same call, the classifier treats the combined note as unverified
>
> Claude Code applies these limits when delivering the note:
>
> * **Length**: Claude Code caps the notes for one tool call at 2,000 characters and truncates the rest. The cap is shared across every hook that responds to that call
> * **Synchronous responses only**: Claude Code ignores the field in the response of a hook that [runs in the background](#run-hooks-in-the-background), because that response arrives after Claude Code records the tool result
> * **Calls that the classifier doesn't record**: the classifier's transcript omits read-only lookups such as file reads and searches. Claude Code discards a note attached to one of those calls
> * **Interaction with rewrites**: when the note describes output you're replacing with `updatedToolOutput`, return both fields in the same hook response. Claude Code drops the note if that rewrite is rejected or another hook's rewrite replaces it. Claude Code delivers a note you return without a rewrite even when another hook rewrites the output
>
>
>   The classifier reads content you place in `classifierContext` as information from the application hosting the session, so don't copy untrusted tool output or third-party text into it. Keep the note to a short assertion about this one call, such as a fact about its origin or a user statement about it; don't use the field to deliver unrelated messages or a stream of events.

### Source: Hooks reference — PostToolUseFailure

> ### PostToolUseFailure
>
> Runs when a tool that started executing fails: the tool threw an error, or an MCP tool returned an error result. Use this to log failures, send alerts, or provide corrective feedback to Claude.
>
> Matches on tool name, same values as PreToolUse.
>
>
>   This event doesn't fire for tool calls rejected before execution: an unknown tool name, input that fails schema or tool-specific validation, or a permission denial. Validation rejections are returned as `tool_use_error` results and happen before hooks run, so they fire neither `PreToolUse` nor `PostToolUseFailure`. Permission denials fire `PreToolUse` but not this event; see [PermissionDenied](#permissiondenied).
>
>
> #### PostToolUseFailure input
>
> PostToolUseFailure hooks receive the same `tool_name` and `tool_input` fields as PostToolUse, along with error information as top-level fields. For example, a failed `npm test` command might deliver:
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "PostToolUseFailure",
>   "tool_name": "Bash",
>   "tool_input": {
>     "command": "npm test",
>     "description": "Run test suite"
>   },
>   "tool_use_id": "toolu_01ABC123...",
>   "error": "Exit code 1\nError: Cannot find module 'express'",
>   "is_interrupt": false,
>   "duration_ms": 4187
> }
> ```
>
> | Field          | Description                                                                                                                                                                                                                    |
> | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `error`        | String describing what went wrong. The format depends on the tool that failed                                                                                                                                                  |
> | `is_interrupt` | Optional boolean. True when the failure reached Claude Code as an abort rather than as an error the tool reported. Cancelling a running tool does not fire this hook; the tool result carries the interruption message instead |
> | `duration_ms`  | Optional. Tool execution time in milliseconds. Excludes time spent in permission prompts and PreToolUse hooks                                                                                                                  |
>
> The `error` string is generally the same text Claude receives as the failed tool's result. Its format varies by tool and failure. Key your hook on `tool_name`, `is_interrupt`, and the `Exit code N` first line; treat the rest of the string as display text, not a stable format.
>
> * For Bash and PowerShell, a command that ran and exited produces a first line `Exit code N`, then any output the command produced as one block with stdout and stderr interleaved
> * A payload may also carry a bare failure message with no exit-code line, when Claude Code could not start the shell process itself
> * Claude Code middle-truncates strings longer than 10,000 characters around a `... [N characters truncated] ...` marker, and can insert lines of its own, such as `Command timed out after 2m 0s`
>
> #### PostToolUseFailure decision control
>
> `PostToolUseFailure` hooks can provide context to Claude after a tool failure. In addition to the [JSON output fields](#json-output) available to all hooks, your hook script can return these event-specific fields:
>
> | Field               | Description                                                                                                 |
> | :------------------ | :---------------------------------------------------------------------------------------------------------- |
> | `additionalContext` | String added to Claude's context alongside the error. See [Add context for Claude](#add-context-for-claude) |
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PostToolUseFailure",
>     "additionalContext": "Additional information about the failure for Claude"
>   }
> }
> ```

### Source: Hooks reference — PostToolBatch

> ### PostToolBatch
>
> Runs once after every tool call in a batch has resolved, before Claude Code sends the next request to the model. `PostToolUse` fires once per tool, which means it fires concurrently when Claude makes parallel tool calls. `PostToolBatch` fires exactly once with the full batch, so it is the right place to inject context that depends on the set of tools that ran rather than on any single tool. There is no matcher for this event.
>
> #### PostToolBatch input
>
> In addition to the [common input fields](#common-input-fields), PostToolBatch hooks receive `tool_calls`, an array describing every tool call in the batch:
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "PostToolBatch",
>   "tool_calls": [
>     {
>       "tool_name": "Read",
>       "tool_input": {"file_path": "/.../ledger/accounts.py"},
>       "tool_use_id": "toolu_01...",
>       "tool_response": "     1\tfrom __future__ import annotations\n     2\t..."
>     },
>     {
>       "tool_name": "Read",
>       "tool_input": {"file_path": "/.../ledger/transactions.py"},
>       "tool_use_id": "toolu_02...",
>       "tool_response": "     1\tfrom __future__ import annotations\n     2\t..."
>     }
>   ]
> }
> ```
>
> `tool_response` contains the same content the model receives in the corresponding `tool_result` block. The value is a serialized string or content-block array, exactly as the tool emitted it. For `Read`, that means line-number-prefixed text rather than raw file contents. Responses can be large, so parse only the fields you need.
>
>
>   The `tool_response` shape differs from `PostToolUse`'s. `PostToolUse` passes the tool's structured `Output` object, such as `{filePath: "...", success: true}` for `Write`; `PostToolBatch` passes the serialized `tool_result` content the model sees.
>
>
> #### PostToolBatch decision control
>
> `PostToolBatch` hooks can inject context for Claude. In addition to the [JSON output fields](#json-output) available to all hooks, your hook script can return these event-specific fields:
>
> | Field               | Description                                                                                                                                                                                         |
> | :------------------ | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `additionalContext` | Context string injected once before the next model call. See [Add context for Claude](#add-context-for-claude) for delivery details, what to put in it, and how resumed sessions handle past values |
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PostToolBatch",
>     "additionalContext": "These files are part of the ledger module. Run pytest before marking the task complete."
>   }
> }
> ```
>
> Returning `decision: "block"` or `continue: false` stops the agentic loop before the next model call. The blocking message comes from the JSON `reason` or `stopReason`, or from stderr on exit 2. You see it as a warning in the transcript, and it stays in the conversation, so Claude sees it when the conversation continues.
