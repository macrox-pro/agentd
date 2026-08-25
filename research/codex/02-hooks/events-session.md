---
primary_sources:
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "SessionStart; SessionEnd; SubagentStart; SubagentStop"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Session and subagent hook events

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex Hooks — SessionStart

> ### SessionStart
>
> `matcher` is applied to `source` for this event.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field    | Type     | Meaning                                                             |
> | -------- | -------- | ------------------------------------------------------------------- |
> | `source` | `string` | How the session started: `startup`, `resume`, `clear`, or `compact` |
>
> Plain text on `stdout` is added as extra developer context.
>
> JSON on `stdout` supports [Common output fields](#common-output-fields) and this
> hook-specific shape:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "SessionStart",
>     "additionalContext": "Load the workspace conventions before editing."
>   }
> }
> ```
>
> That `additionalContext` text is added as extra developer context.
>
> After Codex compacts a root session, `SessionStart` hooks that match
> `source: "compact"` run before the next model request. This also applies when
> automatic compaction happens in the middle of a turn: Codex delivers the hook's
> additional context to the immediate continuation instead of waiting for a
> later user turn. If the hook returns `continue: false`, Codex ends the turn
> without sending another model request.

### Source: Codex Hooks — SessionEnd

> ### SessionEnd
>
> `SessionEnd` lets you run a command when a session ends, such as saving final
> notes or cleaning up files. It runs for the main thread when you archive or
> delete a conversation that's still open, when Codex closes normally, or after a
> conversation has been idle and isn't open in any connected client for 30
> minutes. It won't run for subagents.
>
> Switching away from a conversation or calling `thread/unsubscribe` doesn't end
> the session right away, so it won't immediately run `SessionEnd`. Your hook can
> still read the session transcript while it runs.
>
> `matcher` filters `reason` for this event. For now, `reason` is always `other`.
> You can omit `matcher` or use `other` to run on every `SessionEnd` event.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field    | Type     | Meaning                        |
> | -------- | -------- | ------------------------------ |
> | `reason` | `string` | Why the session ended: `other` |
>
> For example, a `SessionEnd` command receives:
>
> ```json
> {
>   "session_id": "thr_123",
>   "transcript_path": "/workspace/.codex/rollout.jsonl",
>   "cwd": "/workspace",
>   "hook_event_name": "SessionEnd",
>   "reason": "other"
> }
> ```
>
> `SessionEnd` hooks always run synchronously, even when `async` is `true`. They
> are advisory, so their output won't steer Codex or keep the thread open. If a
> command times out or exits with an error, Codex reports it as a hook failure.

### Source: Codex Hooks — SubagentStart

> ### SubagentStart
>
> `matcher` is applied to `agent_type` for this event.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field             | Type     | Meaning                                        |
> | ----------------- | -------- | ---------------------------------------------- |
> | `turn_id`         | `string` | Codex-specific extension. Active Codex turn id |
> | `agent_id`        | `string` | Identifier for the subagent                    |
> | `agent_type`      | `string` | Subagent type or profile                       |
> | `permission_mode` | `string` | Current permission mode                        |
>
> Plain text on `stdout` is added as extra developer context for the subagent.
>
> JSON on `stdout` supports `systemMessage` and this hook-specific shape:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "SubagentStart",
>     "additionalContext": "Review the repository test conventions first."
>   }
> }
> ```
>
> That `additionalContext` text is added as extra developer context for the
> subagent. `continue: false` is parsed for compatibility, but it doesn't stop the
> subagent from starting.

### Source: Codex Hooks — SubagentStop

> ### SubagentStop
>
> `matcher` is applied to `agent_type` for this event.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field                    | Type             | Meaning                                         |
> | ------------------------ | ---------------- | ----------------------------------------------- |
> | `turn_id`                | `string`         | Codex-specific extension. Active Codex turn id  |
> | `agent_id`               | `string`         | Identifier for the subagent                     |
> | `agent_type`             | `string`         | Subagent type or profile                        |
> | `agent_transcript_path`  | `string \| null` | Path to the subagent transcript file, if any    |
> | `stop_hook_active`       | `boolean`        | Whether this subagent was already continued     |
> | `last_assistant_message` | `string \| null` | Latest subagent assistant message, if available |
>
> `SubagentStop` expects JSON on `stdout` when it exits `0`. Plain text output is
> invalid for this event.
>
> JSON on `stdout` supports [Common output fields](#common-output-fields). To ask
> Codex to continue the subagent flow, return:
>
> ```json
> {
>   "decision": "block",
>   "reason": "Run one more focused pass inside the subagent."
> }
> ```
>
> You can also use exit code `2` and write the continuation reason to `stderr`.
>
> If any matching `SubagentStop` hook returns `continue: false`, that takes
> precedence over continuation decisions from other matching `SubagentStop`
> hooks.
