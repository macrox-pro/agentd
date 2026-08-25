---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Hook input and output"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook execution and protocol

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — Hook input and output
>
> ## Hook input and output
>
> Command hooks receive JSON data via stdin and communicate results through exit codes, stdout, and stderr. HTTP hooks receive the same JSON as the POST request body and communicate results through the HTTP response body. This section covers fields and behavior common to all events. Each event's section under [Hook events](#hook-events) includes its specific input schema and decision control options.
>
> On macOS and Linux, command hooks run in their own session without a controlling terminal. The hook process and any child processes can't open `/dev/tty` or send escape sequences directly to the Claude Code interface. Windows has no `/dev/tty`.
>
> To surface a message to the user on any platform, return [`systemMessage`](#json-output) in JSON output. Some events discard it or deliver it elsewhere, and each [event's section](#hook-events) says so. To trigger a desktop notification, set a window title, or ring the bell, return [`terminalSequence`](#emit-terminal-notifications) instead.
>
> ### Common input fields
>
> Hook events receive these fields as JSON, in addition to event-specific fields documented in each [hook event](#hook-events) section. For command hooks, this JSON arrives via stdin. For HTTP hooks, it arrives as the POST request body.
>
> | Field             | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
> | :---------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `session_id`      | Current session identifier                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
> | `prompt_id`       | UUID identifying the user prompt currently being processed. Matches the [`prompt.id` attribute on OpenTelemetry events](/docs/en/monitoring-usage#event-correlation-attributes), so you can correlate hook output with telemetry for a single prompt. Absent until the first user input. Requires Claude Code v2.1.196 or later                                                                                                                                                                                                                                                                                                                                                                                                                       |
> | `transcript_path` | Path to conversation JSON. The transcript file is written asynchronously and may lag the in-memory conversation, so it may not yet include the current turn's most recent messages when a hook fires. Hooks that need the final assistant text of the current turn should use `last_assistant_message` on [Stop](#stop) and [SubagentStop](#subagentstop) instead of reading the transcript                                                                                                                                                                                                                                                                                                                                                      |
> | `cwd`             | Current working directory when the hook is invoked                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
> | `permission_mode` | Current [permission mode](/docs/en/permissions#permission-modes): `"default"`, `"plan"`, `"acceptEdits"`, `"auto"`, `"dontAsk"`, or `"bypassPermissions"`. The mode labeled **Manual** arrives as `"default"`, never as `"manual"`, so scripts that match `"default"` keep working. Not all events receive this field. Check the JSON example in each [hook event](#hook-events) section                                                                                                                                                                                                                                                                                                                                                              |
> | `effort`          | Object with a `level` field holding the active [effort level](/docs/en/model-config#adjust-effort-level) for the turn: `"low"`, `"medium"`, `"high"`, `"xhigh"`, or `"max"`. If the requested model effort exceeds what the current model supports, this is the downgraded level the model actually used. Ultracode is not a distinct level and reports as `"xhigh"`. The object matches the [status line](/docs/en/statusline#available-data) `effort` field. Present for events that fire within a tool-use context, such as `PreToolUse`, `PostToolUse`, `Stop`, and `SubagentStop`, when the current model supports the effort parameter. The level is also available to hook commands and the Bash tool as the `$CLAUDE_EFFORT` environment variable. |
> | `hook_event_name` | Name of the event that fired                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
>
> When running with `--agent` or inside a subagent, two additional fields are included:
>
> | Field        | Description                                                                                                                                                                                                                                                                                                                                                                         |
> | :----------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `agent_id`   | Unique identifier for the subagent. Present only when the hook fires inside a subagent call. Use this to distinguish subagent hook calls from main-thread calls.                                                                                                                                                                                                                    |
> | `agent_type` | Agent name (for example, `"Explore"` or `"security-reviewer"`). Present when the session uses `--agent` or the hook fires inside a subagent. For subagents, the subagent's type takes precedence over the session's `--agent` value. See [SubagentStart](#subagentstart) for the values custom and plugin subagents report and how to write a matcher against a plugin-scoped name. |
>
> Only [`SessionStart`](#sessionstart) hooks can receive a `model` field, and it is not guaranteed to be present. There is no `$CLAUDE_MODEL` environment variable. A hook process inherits the parent environment, so it can read `$ANTHROPIC_MODEL` if you set it in your shell, but that value doesn't change when you switch models with `/model` during a session. One set of variables is not inherited: Claude Code [removes `OTEL_*` exporter variables from every subprocess it spawns](/docs/en/monitoring-usage#administrator-configuration), including hooks.
>
> For example, a `PreToolUse` hook for a Bash command receives this on stdin:
>
> ```json
> {
>   "session_id": "abc123",
>   "prompt_id": "550e8400-e29b-41d4-a716-446655440000",
>   "transcript_path": "/home/user/.claude/projects/.../transcript.jsonl",
>   "cwd": "/home/user/my-project",
>   "permission_mode": "default",
>   "hook_event_name": "PreToolUse",
>   "tool_name": "Bash",
>   "tool_input": {
>     "command": "npm test",
>     "description": "Run test suite",
>     "timeout": 120000,
>     "run_in_background": false
>   },
>   "tool_use_id": "toolu_01ABC123..."
> }
> ```
>
> The `tool_name`, `tool_input`, and `tool_use_id` fields are event-specific. Each [hook event](#hook-events) section documents the additional fields for that event.
>
> ### Exit code output
>
> The exit code from your hook command tells Claude Code whether the action should proceed, be blocked, or be ignored. The exit code doesn't act alone. Claude Code reads [JSON output fields](#json-output) from stdout on every exit code, not just 0, and for events that use the standard decision model, a parsed object that passes schema validation takes effect alongside the code. Exit 2's block is the one outcome JSON can't override.
>
> Two tables own the per-event exceptions: [Exit code 2 behavior per event](#exit-code-2-behavior-per-event) says what exit codes do for each event, and [Decision control](#decision-control) says which decision fields each event honors. Universal fields such as `systemMessage` work across most events and are listed in the [JSON output](#json-output) table.
>
> #### Exit code 0
>
> Exit 0 means success, and is the intended exit code when you print JSON for structured control. For most events, stdout is written to the debug log but not shown in the transcript. The exceptions are `UserPromptSubmit`, `UserPromptExpansion`, and `SessionStart`, where Claude Code adds plain-text stdout as context that Claude can see and act on.
>
> Whether Claude Code reads your stdout as [JSON output](#json-output) or as plain text depends on its first character, ignoring leading whitespace:
>
> * **Starts with `{`**: Claude Code parses it as JSON. If it isn't valid JSON, Claude Code treats it as plain text.
> * **Starts with anything else**: Claude Code treats it as plain text, a JSON array or a quoted JSON string included.
>
> For events that use the standard decision model, exit 0 with a parsed object that fails schema validation is a non-blocking error: the action proceeds, and the transcript shows a `<hook name> hook error` notice with the validation message. The same happens on any exit code other than 2, while [exit 2 still blocks](#exit-code-2).
>
> Stderr from a hook that exits 0 goes to the debug log only, never the transcript, and Claude never sees it. To read it yourself, enable [debug logging](#debug-hooks). To surface a warning to Claude from a `PostToolUse` or `PostToolUseFailure` hook, exit 2 instead so [Claude sees the stderr](#exit-code-2-behavior-per-event) even though the tool already ran.
>
> #### Exit code 2
>
> Exit 2 means a blocking error. On [events that can block](#exit-code-2-behavior-per-event), exit 2 blocks whether or not you print JSON: even a JSON `permissionDecision` of `"allow"` can't override it. Claude Code still reads any valid [JSON output](#json-output) on stdout. On `Elicitation` and `ElicitationResult`, an exit-2 hook's `hookSpecificOutput` is ignored.
>
> The blocking message is the reason from your JSON's blocking decision when it makes one, and your stderr text otherwise. What the block does varies by event: `PreToolUse` blocks the tool call, `UserPromptSubmit` rejects the prompt, and so on. [Exit code 2 behavior per event](#exit-code-2-behavior-per-event) lists the effect for every event, and each event's section says where the message goes.
>
> A hook that exits 2 while printing JSON that fails [JSON output](#json-output) schema validation still blocks: Claude Code uses stderr as the blocking reason and records the validation failure in the debug log. Before v2.1.214, Claude Code treated that combination as a non-blocking error and the action proceeded.
>
> This script blocks `rm` commands by exiting 2 and leaves every other command to the normal permission flow:
>
> ```bash
> #!/bin/bash
> # Reads JSON input from stdin, checks the command
> input=$(cat)
> command=$(jq -r '.tool_input.command' <<<"$input")
>
> if [[ "$command" == rm* ]]; then
>   echo "Blocked: rm commands are not allowed" >&2
>   exit 2  # Blocking error: tool call is prevented
> fi
>
> exit 0  # No decision: the normal permission flow applies
> ```
>
> #### Other exit codes
>
> Any other exit code doesn't block on its own for most hook events. What happens depends on your stdout:
>
> * With a parsed object that passes schema validation, for events that use the standard decision model, Claude Code ignores the exit code and the JSON alone decides the outcome:
>   * Each field the event supports is honored, including `permissionDecision`, `additionalContext`, `updatedInput`, and `systemMessage`, and the hook isn't reported as an error.
>   * [Decision control](#decision-control) lists the decision fields per event; universal fields like `systemMessage` follow the [JSON output](#json-output) table.
> * With a parsed object that fails schema validation, for events that use the standard decision model, it's the same non-blocking error as [on exit 0](#exit-code-0): the action proceeds, and the `<hook name> hook error` notice carries the validation message.
> * With stdout that Claude Code [treats as plain text](#exit-code-0), or with empty stdout, it's a non-blocking error for most hook events: the action proceeds, and the transcript shows a `<hook name> hook error` notice followed by the first line of stderr, prefixed with `Failed with non-blocking status code:`. To capture the full stderr, enable [debug logging](#debug-hooks).
>
> Events outside the standard decision model keep their own rows in the [per-event table](#exit-code-2-behavior-per-event): `WorktreeCreate` fails creation on any nonzero exit no matter what your JSON says, and events that discard hook output entirely, like `StopFailure`, ignore your JSON on every exit code, apart from side-effect fields like `terminalSequence`, which still fire.
>
> A hook that can't start lands in the same non-blocking bucket. When the script path doesn't exist or isn't executable, the shell exits with a code like 127 and you see the same notice with the interpreter's message, for example `Failed with non-blocking status code: /bin/sh: /path/to/hook.sh: No such file or directory`. For most hook events, the action proceeds. When you set up a policy hook, watch for this notice on its first run: a mistyped path in `settings.json` leaves the gate silently disabled.
>
>   For most hook events, exit code 2 is the only exit code that blocks through the code alone. Without valid JSON on stdout, Claude Code treats exit code 1 as a non-blocking error and proceeds with the action, even though 1 is the conventional Unix failure code. If your hook is meant to enforce a policy, use `exit 2`. The exception is `WorktreeCreate`, where any non-zero exit code aborts worktree creation.
>
> #### Timeouts
>
> A `command`, `http`, or `mcp_tool` hook that reaches its [`timeout`](#common-fields) is canceled: Claude Code discards the hook's output, and the hook renders no decision. On `PreToolUse`, the two hook families differ:
>
> * A timed-out `command`, `http`, or `mcp_tool` hook doesn't block the tool call. The call continues through the normal [permission flow](/docs/en/permissions), so don't count on a stalled hook to act as a gate.
> * An [Agent SDK callback hook](/docs/en/agent-sdk/hooks) that exceeds its timeout [blocks the tool call](#pretooluse).
>
> #### Exit code 2 behavior per event
>
> Exit code 2 is the way a hook signals "stop, don't do this." The effect depends on the event, because some events represent actions that can be blocked (like a tool call that hasn't happened yet) and others represent things that already happened or can't be prevented.
>
> | Hook event            | Can block? | What happens on exit 2                                                                                                                                                                                                                         |
> | :-------------------- | :--------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `PreToolUse`          | Yes        | Blocks the tool call                                                                                                                                                                                                                           |
> | `PermissionRequest`   | No         | Exit code 2 isn't honored for this event and the permission flow proceeds unchanged. Deny through the [`decision` object](#permissionrequest-decision-control) instead                                                                         |
> | `UserPromptSubmit`    | Yes        | Blocks prompt processing and erases the prompt                                                                                                                                                                                                 |
> | `UserPromptExpansion` | Yes        | Blocks the expansion                                                                                                                                                                                                                           |
> | `Stop`                | Yes        | Prevents Claude from stopping, continues the conversation                                                                                                                                                                                      |
> | `SubagentStop`        | Yes        | Prevents the subagent from stopping                                                                                                                                                                                                            |
> | `TeammateIdle`        | Yes        | Prevents the teammate from going idle, so it continues working                                                                                                                                                                                 |
> | `TaskCreated`         | Yes        | Rolls back the task creation                                                                                                                                                                                                                   |
> | `TaskCompleted`       | Yes        | Prevents the task from being marked as completed                                                                                                                                                                                               |
> | `ConfigChange`        | Yes        | Blocks the configuration change from taking effect (except `policy_settings`)                                                                                                                                                                  |
> | `StopFailure`         | No         | Output and exit code are ignored, except `terminalSequence`                                                                                                                                                                                    |
> | `PostToolUse`         | No         | Shows stderr to Claude; the tool already ran                                                                                                                                                                                                   |
> | `PostToolUseFailure`  | No         | Shows stderr to Claude; the tool already failed                                                                                                                                                                                                |
> | `PostToolBatch`       | Yes        | Stops the agentic loop before the next model call                                                                                                                                                                                              |
> | `PermissionDenied`    | No         | Exit code and stderr are ignored because the denial already occurred. Use JSON `hookSpecificOutput.retry: true` to tell the model it may retry; Claude Code ignores `retry: true` for [no-verdict denials](#permissiondenied-decision-control) |
> | `Notification`        | No         | Exit code and stderr are ignored                                                                                                                                                                                                               |
> | `SubagentStart`       | No         | Shows stderr to user only                                                                                                                                                                                                                      |
> | `SessionStart`        | No         | Shows stderr to user only                                                                                                                                                                                                                      |
> | `Setup`               | No         | Shows stderr to user only                                                                                                                                                                                                                      |
> | `SessionEnd`          | No         | Shows stderr to user only                                                                                                                                                                                                                      |
> | `CwdChanged`          | No         | Shows stderr to user only                                                                                                                                                                                                                      |
> | `DirectoryAdded`      | No         | Stderr goes to the debug log; the directory is already added                                                                                                                                                                                   |
> | `FileChanged`         | No         | Shows stderr to user only                                                                                                                                                                                                                      |
> | `PreCompact`          | Yes        | Blocks compaction                                                                                                                                                                                                                              |
> | `PostCompact`         | No         | Shows stderr to user only                                                                                                                                                                                                                      |
> | `Elicitation`         | Yes        | Denies the elicitation                                                                                                                                                                                                                         |
> | `ElicitationResult`   | Yes        | Blocks the response (action becomes decline)                                                                                                                                                                                                   |
> | `WorktreeCreate`      | Yes        | Any non-zero exit code causes worktree creation to fail                                                                                                                                                                                        |
> | `WorktreeRemove`      | No         | Failures are logged in debug mode only                                                                                                                                                                                                         |
> | `InstructionsLoaded`  | No         | Exit code is ignored                                                                                                                                                                                                                           |
> | `MessageDisplay`      | No         | The original text is displayed                                                                                                                                                                                                                 |
>
> For `SessionStart`, `Setup`, and `SubagentStart`, the exit code 2 stderr renders in the transcript as a `<hook name> hook error` notice, the same way a [non-blocking error](#exit-code-output) does. Claude doesn't see it, and the session or subagent proceeds. For `SubagentStart`, the notice appears in the subagent's own transcript, not in the parent conversation.
>
> ### HTTP response handling
>
> HTTP hooks use HTTP status codes and response bodies instead of exit codes and stdout. The outcomes below apply to most events; an event with its own failure contract in the [per-event table](#exit-code-2-behavior-per-event), such as `WorktreeCreate`, applies that contract to a failed HTTP hook too:
>
> * **2xx with an empty body**: success, equivalent to exit code 0 with no output
> * **2xx with a JSON object body**: parsed using the same [JSON output](#json-output) schema as command hooks. A body that fails schema validation is a non-blocking error
> * **2xx with any other body, such as plain text**: non-blocking error, handled the same as a non-2xx status. Claude Code doesn't add the text to Claude's context
> * **Non-2xx status**: non-blocking error, execution continues
> * **Connection failure**: non-blocking error, execution continues
> * **Timeout**: the hook is canceled and renders no decision, and execution continues
>
> Unlike command hooks, HTTP hooks can't signal a blocking error through status codes alone. To block a tool call or deny a permission, return a 2xx response with a JSON body containing the appropriate decision fields.
>
> ### JSON output
>
> Exit codes only let you block or stay silent, but JSON output gives you finer-grained control. Instead of exiting with code 2 to block, exit 0 and print a JSON object to stdout. Claude Code reads specific fields from that JSON to control behavior, including [decision control](#decision-control) for blocking, allowing, or escalating to the user.
>
>   Choose one approach per hook: either use exit codes alone for signaling, or exit 0 and print JSON for structured control. If you mix them, exit 2 keeps its [blocking effect](#exit-code-2-behavior-per-event), and Claude Code still reads the JSON fields, with the one elicitation exception noted under [Exit code 2](#exit-code-2).
>
> Your hook's stdout must contain only the JSON object. If your shell profile prints text on startup, it can interfere with JSON parsing. See [Hook JSON has no effect](/docs/en/hooks-guide#hook-json-has-no-effect) in the troubleshooting guide.
>
> Hook output strings, including `additionalContext`, `systemMessage`, and plain stdout, are capped at 10,000 characters. Output that exceeds this limit is saved to a file and replaced with a preview and file path, the same way a large valid Bash result is handled under [Output limits](/docs/en/tools-reference#output-limits).
>
> The JSON object supports three kinds of fields:
>
> * **Universal fields** like `continue` are listed in the table below. Every event accepts them, but some events discard them or deliver `systemMessage` somewhere other than the transcript. Each event's section says so. `terminalSequence` works on those events too, with the exceptions listed under [Emit terminal notifications](#emit-terminal-notifications).
> * **Top-level `decision` and `reason`** are used by some events to block or provide feedback.
> * **`hookSpecificOutput`** is a nested object for events that need richer control. It requires a `hookEventName` field set to the event name.
>
> | Field              | Default | Description                                                                                                                                                                                                                                                                                                                          |
> | :----------------- | :------ | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `continue`         | `true`  | If `false`, Claude stops processing entirely after the hook runs. Takes precedence over any event-specific decision fields                                                                                                                                                                                                           |
> | `stopReason`       | none    | Message shown to the user when `continue` is `false`. Not shown to Claude                                                                                                                                                                                                                                                            |
> | `suppressOutput`   | `false` | Has no effect: Claude Code accepts the field but doesn't act on it. A successful hook's stdout is never shown in the transcript and is recorded in the debug log                                                                                                                                                                     |
> | `systemMessage`    | none    | Warning message shown to the user. In [Agent SDK](/docs/en/agent-sdk/overview) and [`--output-format stream-json`](/docs/en/headless) output, it can arrive as an [`SDKInformationalMessage`](/docs/en/agent-sdk/typescript#sdkinformationalmessage)                                                                                                |
> | `terminalSequence` | none    | A terminal escape sequence for Claude Code to emit on your behalf, such as a desktop notification, window title, or bell. Restricted to OSC `0`/`1`/`2`/`9`/`99`/`777` and BEL. If the value contains anything outside the allowlist, the field is ignored. Use this instead of writing to `/dev/tty`, which is unavailable to hooks |
>
> To stop Claude entirely:
>
> ```json
> { "continue": false, "stopReason": "Build failed, fix errors before continuing" }
> ```
>
> For `PreToolUse` and `PostToolUse` hooks, the stop applies even when the tool call fails or completes while Claude is still streaming a response.
>
> #### Emit terminal notifications
>
> Hooks run without a controlling terminal, so writing escape sequences directly to `/dev/tty` fails. Instead, return the escape sequence in the `terminalSequence` field and Claude Code emits it for you through its own terminal write path. This is race-free, works inside tmux and GNU screen, and works on Windows where there is no `/dev/tty`.
>
> The field accepts a string of one or more allowlisted escape sequences:
>
> * OSC `0`, `1`, `2`: window and icon titles
> * OSC `9`: iTerm2, ConEmu, Windows Terminal, and WezTerm notifications, including `9;4` taskbar progress
> * OSC `99`: Kitty notifications
> * OSC `777`: urxvt, Ghostty, and Warp notifications
> * Bare BEL
>
> Sequences may be terminated with BEL or with ST. Anything outside the allowlist, including CSI cursor and color sequences, OSC palette sequences, OSC 8 hyperlinks, OSC 52 clipboard writes, and OSC 1337, is rejected and the field is ignored.
>
> Claude Code writes the sequence itself when it processes your hook's output, so the field works on events that discard `systemMessage` and `continue`, such as `Notification` and `StopFailure`. It has two limits:
>
> * Claude Code writes the sequence only in an interactive session, and only while its interface is on screen. In non-interactive mode with the `-p` flag and in the Agent SDK, it ignores the field.
> * A `WorktreeCreate` command hook can't return JSON, because Claude Code reads its stdout as the worktree path. An HTTP `WorktreeCreate` hook returns JSON and can include the field.
>
> The example below fires a desktop notification from a `Notification` hook. The escape sequence is built with `printf` octal escapes so the control bytes never appear on the shell command line, and `jq -n --arg` builds the JSON output so quotes, backslashes, and newlines in the notification message are escaped correctly:
>
> ```bash
> #!/bin/bash
> # Notification hook: ping the desktop when Claude Code needs attention.
> input=$(cat)
> title="Claude Code"
> body=$(jq -r '.message // "Needs your attention"' <<<"$input")
> seq=$(printf '\033]777;notify;%s;%s\007' "$title" "$body")
> jq -nc --arg seq "$seq" '{terminalSequence: $seq}'
> ```
>
> The `{ "terminalSequence": "..." }` shape is the same from any shell or language.
>
> #### Add context for Claude
>
> The `additionalContext` field passes a string from your hook into Claude's context window. Claude Code wraps the string in a system reminder and inserts it into the conversation at the point where the hook fired. Claude reads the reminder on the next model request, but it doesn't appear as a chat message in the interface.
>
> Return `additionalContext` inside `hookSpecificOutput` alongside the event name:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "PostToolUse",
>     "additionalContext": "This file is generated. Edit src/schema.ts and run `bun generate` instead."
>   }
> }
> ```
>
> Where the reminder appears depends on the event:
>
> * [SessionStart](#sessionstart), [Setup](#setup), and [SubagentStart](#subagentstart): at the start of the conversation, before the first prompt
> * [UserPromptSubmit](#userpromptsubmit) and [UserPromptExpansion](#userpromptexpansion): alongside the submitted prompt
> * [PreToolUse](#pretooluse), [PostToolUse](#posttooluse), [PostToolUseFailure](#posttoolusefailure), and [PostToolBatch](#posttoolbatch): next to the tool result
> * [Stop](#stop) and [SubagentStop](#subagentstop): at the end of the turn. The conversation continues so Claude can act on the feedback. See [Stop decision control](#stop-decision-control)
>
> When several hooks return `additionalContext` for the same event, Claude receives all of the values. If a value exceeds 10,000 characters, Claude Code writes the full text to a file in the session directory and passes Claude the file path with a short preview instead.
>
> Use `additionalContext` for information Claude should know about the current state of your environment or the operation that just ran:
>
> * **Environment state**: the current branch, deployment target, or active feature flags
> * **Conditional project rules**: which test command applies to the file just edited, which directories are read-only in this worktree
> * **External data**: open issues assigned to you, recent CI results, content fetched from an internal service
>
> For instructions that never change, prefer [CLAUDE.md](/docs/en/memory). It loads without running a script and is the standard place for static project conventions.
>
> Write the text as factual statements rather than imperative system instructions. Phrasing such as "The deployment target is production" or "This repo uses `bun test`" reads as project information. Text framed as out-of-band system commands can trigger Claude's prompt-injection defenses, which causes Claude to surface the text to you instead of treating it as context.
>
> Claude Code saves the injected text in the session transcript. For mid-session events like `PostToolUse` or `UserPromptSubmit`, when you resume with `--continue` or `--resume`, Claude Code replays the saved text rather than re-running the hook for past turns, so values like timestamps or commit SHAs become stale. `SessionStart` hooks run again on resume with `source` set to `"resume"`, or `"fork"` if you added `--fork-session`, so they can refresh their context.
>
> #### Decision control
>
> Not every event supports blocking or controlling behavior through JSON. The events that do each use a different set of fields to express that decision. Use this table as a quick reference before writing a hook:
>
> | Events                                                                                                                              | Decision pattern                  | Key fields                                                                                                                                                                                                                          |
> | :---------------------------------------------------------------------------------------------------------------------------------- | :-------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | UserPromptSubmit, UserPromptExpansion, PostToolUse, PostToolUseFailure, PostToolBatch, Stop, SubagentStop, ConfigChange, PreCompact | Top-level `decision`              | `decision: "block"`, `reason`. Stop and SubagentStop also accept `hookSpecificOutput.additionalContext` for [non-error feedback that continues the conversation](#stop-decision-control)                                            |
> | TeammateIdle, TaskCompleted                                                                                                         | Exit code or `continue: false`    | Exit code 2 blocks the action with stderr feedback. JSON `{"continue": false, "stopReason": "..."}` also stops the teammate entirely, matching `Stop` hook behavior                                                                 |
> | TaskCreated                                                                                                                         | Exit code or top-level `decision` | Exit code 2 or `decision: "block"` [cancels the task](#taskcreated-decision-control) and returns the message to Claude. `continue: false` is ignored                                                                                |
> | PreToolUse                                                                                                                          | `hookSpecificOutput`              | `permissionDecision` (allow/deny/ask/defer), `permissionDecisionReason`                                                                                                                                                             |
> | PermissionRequest                                                                                                                   | `hookSpecificOutput`              | `decision.behavior` (allow/deny)                                                                                                                                                                                                    |
> | PermissionDenied                                                                                                                    | `hookSpecificOutput`              | `retry: true` tells the model it may retry the denied tool call; Claude Code ignores it for [no-verdict denials](#permissiondenied-decision-control)                                                                                |
> | WorktreeCreate                                                                                                                      | path return                       | Command hook prints path on stdout; HTTP hook returns `hookSpecificOutput.worktreePath`. Hook failure or missing path fails creation                                                                                                |
> | Elicitation                                                                                                                         | `hookSpecificOutput`              | `action` (accept/decline/cancel), `content` (form field values for accept)                                                                                                                                                          |
> | ElicitationResult                                                                                                                   | `hookSpecificOutput`              | `action` (accept/decline/cancel), `content` (form field values override)                                                                                                                                                            |
> | MessageDisplay                                                                                                                      | `hookSpecificOutput`              | `displayContent` replaces the displayed text on screen. Display-only: the transcript and what Claude sees keep the original                                                                                                         |
> | SessionStart, Setup, SubagentStart                                                                                                  | Context only                      | `hookSpecificOutput.additionalContext` adds context for Claude. SessionStart also accepts [`initialUserMessage`, `watchPaths`, `sessionTitle`, and `reloadSkills`](#sessionstart-decision-control). No blocking or decision control |
> | WorktreeRemove, Notification, SessionEnd, PostCompact, InstructionsLoaded, StopFailure, CwdChanged, DirectoryAdded, FileChanged     | None                              | No decision control. Used for side effects like logging or cleanup                                                                                                                                                                  |
>
> A few events can also rewrite content rather than only allow or block it:
>
> * `PreToolUse`: `updatedInput` directly under `hookSpecificOutput` replaces a tool's arguments before it runs. See [PreToolUse decision control](#pretooluse-decision-control)
> * `PermissionRequest`: `updatedInput` inside the `decision` object. See [PermissionRequest decision control](#permissionrequest-decision-control)
> * `PostToolUse`: `updatedToolOutput` replaces the tool's result. See [PostToolUse decision control](#posttooluse-decision-control)
> * `UserPromptSubmit`: can't replace the prompt; it only injects `additionalContext` alongside it
>
> For redaction or transformation use cases, intercept at `PreToolUse` for outbound tool inputs and `PostToolUse` for inbound tool results.
>
> Here are examples of each pattern in action:
>
>   #### Top-level decision
>
> The only value for `decision` is `"block"`. To allow the action to proceed, omit `decision` from your JSON, or exit 0 without any JSON at all:
>
>     ```json
>     {
>       "decision": "block",
>       "reason": "Test suite must pass before proceeding"
>     }
>     ```
>
>   #### PreToolUse
>
> Uses `hookSpecificOutput` for richer control: allow, deny, or escalate to the user. You can also modify tool input before it runs or inject additional context for Claude. See [PreToolUse decision control](#pretooluse-decision-control) for the full set of options.
>
>     ```json
>     {
>       "hookSpecificOutput": {
>         "hookEventName": "PreToolUse",
>         "permissionDecision": "deny",
>         "permissionDecisionReason": "Database writes are not allowed"
>       }
>     }
>     ```
>
>   #### PermissionRequest
>
> Uses `hookSpecificOutput` to allow or deny a permission request on behalf of the user. When allowing, you can also modify the tool's input or apply permission rules so the user isn't prompted again. See [PermissionRequest decision control](#permissionrequest-decision-control) for the full set of options.
>
>     ```json
>     {
>       "hookSpecificOutput": {
>         "hookEventName": "PermissionRequest",
>         "decision": {
>           "behavior": "allow",
>           "updatedInput": {
>             "command": "npm run lint"
>           }
>         }
>       }
>     }
>     ```
>
> For extended examples including Bash command validation, prompt filtering, and auto-approval scripts, see [What you can automate](/docs/en/hooks-guide#what-you-can-automate) in the guide and the [Bash command validator reference implementation](https://github.com/anthropics/claude-code/blob/main/examples/hooks/bash_command_validator_example.py).### Source: Hooks reference — Prompt-based hooks
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
>
> ### How prompt-based hooks work
>
> Instead of executing a Bash command, prompt-based hooks:
>
> 1. Send the hook input and your prompt to a Claude model, Haiku by default
> 2. The LLM responds with structured JSON containing a decision
> 3. Claude Code processes the decision automatically
>
> ### Prompt hook configuration
>
> Set `type` to `"prompt"` and provide a `prompt` string instead of a `command`. Use the `$ARGUMENTS` placeholder to inject the hook's JSON input data into your prompt text.
>
> This `Stop` hook asks the LLM to evaluate whether all tasks are complete before allowing Claude to finish:
>
> ```json
> {
>   "hooks": {
>     "Stop": [
>       {
>         "hooks": [
>           {
>             "type": "prompt",
>             "prompt": "Evaluate if Claude should stop: $ARGUMENTS. Check if all tasks are complete."
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> | Field             | Required | Description                                                                                                                                                                                               |
> | :---------------- | :------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `type`            | yes      | Must be `"prompt"`                                                                                                                                                                                        |
> | `prompt`          | yes      | The prompt text to send to the LLM. Use `$ARGUMENTS` as a placeholder for the hook input JSON. If `$ARGUMENTS` is not present, input JSON is appended to the prompt                                       |
> | `model`           | no       | Model to use for evaluation. Defaults to a fast model                                                                                                                                                     |
> | `timeout`         | no       | Timeout in seconds. Default: 30                                                                                                                                                                           |
> | `continueOnBlock` | no       | On the events it applies to, `true` feeds an `ok: false` reason back to Claude and continues instead of ending the turn. Default: `false`. See [Response schema](#response-schema) for per-event behavior |
>
> ### Response schema
>
> The LLM must respond with JSON containing:
>
> ```json
> {
>   "ok": true | false,
>   "reason": "Explanation for the decision",
>   "impossible": true | false
> }
> ```
>
> | Field        | Description                                                                                                                                                                                                                                      |
> | :----------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `ok`         | `true` to allow. For `false`, see the per-event behavior below                                                                                                                                                                                   |
> | `reason`     | Required when `ok` is `false`                                                                                                                                                                                                                    |
> | `impossible` | Optional. The model returns it with `ok: false` when it judges the condition can never be satisfied. On `Stop` and `SubagentStop`, Claude Code then lets the turn end instead of feeding the reason back. Agent hooks and other events ignore it |
>
> What happens on `ok: false` depends on the event:
>
> * `Stop` and `SubagentStop`: the reason is fed back to Claude as its next instruction and the turn continues, unless the response also sets `impossible: true`, in which case Claude Code allows the stop and the turn ends
> * `PreToolUse`: the tool call is denied; by default the turn ends and the deny reason appears in the chat as a warning line. Set `continueOnBlock: true` to instead return the reason to Claude as the tool error so it can adjust and continue, equivalent to a command hook's `permissionDecision: "deny"`. Before v2.1.210, the deny reason was returned to Claude as the tool error and the turn continued
> * `PostToolUse`: by default the turn ends and the reason appears in the chat as a warning line. Set `continueOnBlock: true` to feed the reason back to Claude and continue the turn instead
> * `PostToolBatch`, `UserPromptSubmit`, and `UserPromptExpansion`: the turn ends and the reason appears as a warning line. These events end the turn on `decision: "block"` regardless of `continue`
> * `PostToolUseFailure` and `TaskCreated`: the reason is returned to Claude as a tool error and the turn continues, regardless of `continueOnBlock`
> * `TaskCompleted`: when it fires because a task is marked completed during a turn, the reason is returned to Claude as a tool error and the turn continues, regardless of `continueOnBlock`. When it fires because a teammate stops, it behaves like `TeammateIdle` and halts the teammate by default
> * `TeammateIdle`: by default the teammate stops and the reason appears as a warning line. Set `continueOnBlock: true` to feed the reason back to the teammate and keep it working instead
> * `PermissionRequest`: `ok: false` has no effect. To deny an approval from a hook, use a [command hook](#command-hook-fields) returning `hookSpecificOutput.decision.behavior: "deny"`
> * `PermissionDenied`: `ok: false` has no effect because the denial already happened. The only output this event reads is `hookSpecificOutput.retry`, which prompt and agent hooks can't set. They run on this event, but their output is discarded. Use a [command hook](#command-hook-fields) to return `retry`
>
> If you need finer control on any event, use a [command hook](#command-hook-fields) with the per-event fields described in [Decision control](#decision-control).
>
> ### Check multiple conditions before stopping
>
> This `Stop` hook uses a detailed prompt to check three conditions before allowing Claude to stop. `SubagentStop` hooks use the same format to evaluate whether a [subagent](/docs/en/sub-agents) should stop. If the model returns `"ok": false` because the condition isn't met yet, Claude continues working with the provided reason as its next instruction:
>
> ```json
> {
>   "hooks": {
>     "Stop": [
>       {
>         "hooks": [
>           {
>             "type": "prompt",
>             "prompt": "You are evaluating whether Claude should stop working. Context: $ARGUMENTS\n\nAnalyze the conversation and determine if:\n1. All user-requested tasks are complete\n2. Any errors need to be addressed\n3. Follow-up work is needed\n\nRespond with JSON: {\"ok\": true} to allow stopping, or {\"ok\": false, \"reason\": \"your explanation\"} to continue working.",
>             "timeout": 30
>           }
>         ]
>       }
>     ]
>   }
> }
> ```### Source: Hooks reference — Agent-based hooks
>
> ## Agent-based hooks
>
>   Agent hooks are experimental. Behavior and configuration may change in future releases. For production workflows, prefer [command hooks](#command-hook-fields).
>
> Agent-based hooks (`type: "agent"`) are like prompt-based hooks but with multi-turn tool access. Instead of a single LLM call, an agent hook spawns a subagent that can read files, search code, and inspect the codebase to verify conditions. Agent hooks support the same events as prompt-based hooks.
>
> ### How agent hooks work
>
> When an agent hook fires:
>
> 1. Claude Code spawns a subagent with your prompt and the hook's JSON input
> 2. The subagent can use tools like Read, Grep, and Glob to investigate
> 3. After up to 50 turns, the subagent returns a structured `{ "ok": true/false }` decision
> 4. Claude Code allows the action if `ok` is `true`. If `ok` is `false`, Claude Code handles the block the same way as a prompt hook with `continueOnBlock: true` on that event, as listed under [Response schema](#response-schema)
>
> Agent hooks are useful when verification requires inspecting actual files or test output, not just evaluating the hook input data alone.
>
> ### Agent hook configuration
>
> Set `type` to `"agent"` and provide a `prompt` string. The configuration fields are the same as [prompt hooks](#prompt-hook-configuration), except that agent hooks have a longer default timeout and no `continueOnBlock` field:
>
> | Field     | Required | Description                                                                                 |
> | :-------- | :------- | :------------------------------------------------------------------------------------------ |
> | `type`    | yes      | Must be `"agent"`                                                                           |
> | `prompt`  | yes      | Prompt describing what to verify. Use `$ARGUMENTS` as a placeholder for the hook input JSON |
> | `model`   | no       | Model to use. Defaults to a fast model                                                      |
> | `timeout` | no       | Timeout in seconds. Default: 60                                                             |
>
> The response schema is `{ "ok": true }` to allow or `{ "ok": false, "reason": "..." }` to block. On `ok: false`, Claude Code handles an agent hook the way it handles a [prompt hook with `continueOnBlock: true`](#response-schema) on the same event; agent hooks have no `continueOnBlock` field, and don't support the prompt-hook `impossible` field.
>
> This `Stop` hook verifies that all unit tests pass before allowing Claude to finish:
>
> ```json
> {
>   "hooks": {
>     "Stop": [
>       {
>         "hooks": [
>           {
>             "type": "agent",
>             "prompt": "Verify that all unit tests pass. Run the test suite and check the results. $ARGUMENTS",
>             "timeout": 120
>           }
>         ]
>       }
>     ]
>   }
> }
> ```### Source: Hooks reference — Run hooks in the background
>
> ## Run hooks in the background
>
> By default, hooks block Claude's execution until they complete. For long-running tasks like deployments, test suites, or external API calls, set `"async": true` to run the hook in the background while Claude continues working. Async hooks can't block or control Claude's behavior: response fields like `decision`, `permissionDecision`, and `continue` have no effect, because the action they would have controlled has already completed.
>
> ### Configure an async hook
>
> Add `"async": true` to a command hook's configuration to run it in the background without blocking Claude. This field is only available on `type: "command"` hooks.
>
> This hook runs a test script after every `Write` tool call. Claude continues working immediately while `run-tests.sh` executes for up to 120 seconds. When the script finishes, its output is delivered on the next conversation turn:
>
> ```json
> {
>   "hooks": {
>     "PostToolUse": [
>       {
>         "matcher": "Write",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/path/to/run-tests.sh",
>             "async": true,
>             "timeout": 120
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> The `timeout` field sets the maximum time in seconds for the background process. If not specified, async hooks use the same 10-minute default as sync hooks.
>
> Claude Code delivers an async hook's results only while the session runs:
>
> * In [non-interactive mode](/docs/en/headless) with the `-p` flag, Claude Code kills any async hook still running at teardown and finalizes it with outcome `cancelled`
> * If your hook's work must outlive a `claude -p` session, start a fully detached process from it
>
> ### How async hooks execute
>
> When an async hook fires, Claude Code starts the hook process and immediately continues without waiting for it to finish. The hook receives the same JSON input via stdin as a synchronous hook.
>
> After the background process exits, Claude Code delivers the `additionalContext` and `systemMessage` fields from the hook's JSON response to Claude on the next conversation turn. Unlike a synchronous hook's `systemMessage`, neither field is shown to you.
>
> Claude Code validates that JSON response against the same [output schema](#json-output) as synchronous hooks, and drops any field whose value has the wrong type, such as a `systemMessage` that isn't a string, instead of delivering it. Run with `--debug` to see a warning naming each dropped field. Before v2.1.202, malformed JSON output from an async hook could crash the session, and the crash recurred each time the session was resumed.
>
> Async hook completion notifications are suppressed by default. To see them, enable verbose mode with `Ctrl+O` or start Claude Code with `--verbose`.
>
> ### Run tests after file changes
>
> This hook starts a test suite in the background whenever Claude writes a file, then reports the results back to Claude when the tests finish. Save this script to `.claude/hooks/run-tests-async.sh` in your project and make it executable with `chmod +x`:
>
> ```bash
> #!/bin/bash
> # run-tests-async.sh
>
> # Read hook input from stdin
> INPUT=$(cat)
> FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')
>
> # Only run tests for source files
> if [[ "$FILE_PATH" != *.ts && "$FILE_PATH" != *.js ]]; then
>   exit 0
> fi
>
> # Run tests and report results to Claude via additionalContext
> RESULT=$(npm test 2>&1)
> EXIT_CODE=$?
>
> if [ $EXIT_CODE -eq 0 ]; then
>   MSG="Tests passed after editing $FILE_PATH"
> else
>   MSG="Tests failed after editing $FILE_PATH: $RESULT"
> fi
> jq -nc --arg msg "$MSG" '{hookSpecificOutput: {hookEventName: "PostToolUse", additionalContext: $msg}}'
> ```
>
> Then add this configuration to `.claude/settings.json` in your project root. The `async: true` flag lets Claude keep working while tests run:
>
> ```json
> {
>   "hooks": {
>     "PostToolUse": [
>       {
>         "matcher": "Write|Edit",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/run-tests-async.sh",
>             "args": [],
>             "async": true,
>             "timeout": 300
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> ### Limitations
>
> Async hooks have additional constraints compared to synchronous hooks:
>
> * Hook output is delivered on the next conversation turn. If the session is idle, the response waits until the next user interaction. Exception: an `asyncRewake` hook that exits with code 2 wakes Claude immediately even when the session is idle.
> * Each execution creates a separate background process. There is no deduplication across multiple firings of the same async hook.