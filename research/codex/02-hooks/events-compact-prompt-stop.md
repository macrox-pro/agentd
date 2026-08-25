---
primary_sources:
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "PreCompact; PostCompact; UserPromptSubmit; Stop"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Compact, prompt, and stop hook events

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex Hooks — PreCompact

> ### PreCompact
>
> `PreCompact` runs before Codex compacts the chat. `matcher` is applied
> to `trigger`, whose values are `manual` and `auto`.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field     | Type     | Meaning                                        |
> | --------- | -------- | ---------------------------------------------- |
> | `turn_id` | `string` | Codex-specific extension. Active Codex turn id |
> | `trigger` | `string` | What triggered compaction: `manual` or `auto`  |
>
> Plain text on `stdout` is ignored.
>
> JSON on `stdout` supports [Common output fields](#common-output-fields). If a
> matching `PreCompact` hook returns `continue: false`, Codex stops before
> compacting.

### Source: Codex Hooks — PostCompact

> ### PostCompact
>
> `PostCompact` runs after Codex compacts the chat. `matcher` is applied
> to `trigger`, whose values are `manual` and `auto`.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field     | Type     | Meaning                                        |
> | --------- | -------- | ---------------------------------------------- |
> | `turn_id` | `string` | Codex-specific extension. Active Codex turn id |
> | `trigger` | `string` | What triggered compaction: `manual` or `auto`  |
>
> Plain text on `stdout` is ignored.
>
> JSON on `stdout` supports [Common output fields](#common-output-fields). If a
> matching `PostCompact` hook returns `continue: false`, Codex stops after
> compacting.

### Source: Codex Hooks — UserPromptSubmit

> ### UserPromptSubmit
>
> `matcher` isn't currently used for this event.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field     | Type     | Meaning                                        |
> | --------- | -------- | ---------------------------------------------- |
> | `turn_id` | `string` | Codex-specific extension. Active Codex turn id |
> | `prompt`  | `string` | User prompt that's about to be sent            |
>
> Plain text on `stdout` is added as extra developer context.
>
> JSON on `stdout` supports [Common output fields](#common-output-fields) and
> this hook-specific shape:
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "UserPromptSubmit",
>     "additionalContext": "Ask for a clearer reproduction before editing files."
>   }
> }
> ```
>
> That `additionalContext` text is added as extra developer context.
>
> To block the prompt, return:
>
> ```json
> {
>   "decision": "block",
>   "reason": "Ask for confirmation before doing that."
> }
> ```
>
> You can also use exit code `2` and write the blocking reason to `stderr`.

### Source: Codex Hooks — Stop

> ### Stop
>
> `matcher` isn't currently used for this event.
>
> Fields in addition to [Common input fields](#common-input-fields):
>
> | Field                    | Type             | Meaning                                           |
> | ------------------------ | ---------------- | ------------------------------------------------- |
> | `turn_id`                | `string`         | Codex-specific extension. Active Codex turn id    |
> | `stop_hook_active`       | `boolean`        | Whether this turn was already continued by `Stop` |
> | `last_assistant_message` | `string \| null` | Latest assistant message text, if available       |
>
> `Stop` expects JSON on `stdout` when it exits `0`. Plain text output is invalid
> for this event.
>
> JSON on `stdout` supports [Common output fields](#common-output-fields). To keep
> Codex going, return:
>
> ```json
> {
>   "decision": "block",
>   "reason": "Run one more pass over the failing tests."
> }
> ```
>
> You can also use exit code `2` and write the continuation reason to `stderr`.
>
> For this event, `decision: "block"` doesn't reject the turn. Instead, it tells
> Codex to continue and automatically creates a new continuation prompt that acts
> as a new user prompt, using your `reason` as that prompt text.
>
> If any matching `Stop` hook returns `continue: false`, that takes precedence
> over continuation decisions from other matching `Stop` hooks.
