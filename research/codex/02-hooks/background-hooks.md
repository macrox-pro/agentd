---
primary_sources:
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "Run hooks in the background"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Background hooks

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex Hooks — Run hooks in the background

> ## Run hooks in the background
>
> By default, Codex waits for a command hook to finish before continuing the
> operation that triggered it. Set `async` to `true` to run a command hook in the
> background while Codex continues.
>
> ### Configure a background hook
>
> Add `"async": true` to a command handler in `hooks.json`:
>
> ```json
> {
>   "hooks": {
>     "PostToolUse": [
>       {
>         "matcher": "Bash",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "python3 ~/.codex/hooks/post_tool_use.py",
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
> For an inline hook in `config.toml`, set `async = true`:
>
> ```toml
> [[hooks.PostToolUse]]
> matcher = "Bash"
>
> [[hooks.PostToolUse.hooks]]
> type = "command"
> command = "python3 ~/.codex/hooks/post_tool_use.py"
> async = true
> timeout = 120
> ```
>
> Background hooks use the same input, matcher, trust review, timeout, and
> [large-output handling](#large-hook-output) as synchronous command hooks. As
> with other command hooks, `timeout` is measured in seconds and defaults to
> `600`.
>
> ### How background hooks run
>
> When a background hook finishes, Codex delivers supported informational output
> at the next safe point in the conversation:
>
> - If a turn is active, Codex waits for the current model request and tool calls
>   to finish, then makes the output available to the next model request in that
>   turn.
> - If no turn is active, Codex waits until the next user turn. Finishing a
>   background hook doesn't start a new turn.
>
> Use the same event-specific JSON output as a synchronous hook. Codex adds
> `additionalContext` to the model's context and surfaces `systemMessage` as a
> warning.
>
> Background hooks can't block, approve, rewrite, or otherwise control the
>   operation that triggered them. Use synchronous hooks for tool policies,
>   permission decisions, prompt rejection, or turn continuation.
>
> ### Limitations
>
> - Codex runs up to eight background hooks concurrently per session. Additional
>   hooks wait until a running hook finishes.
> - Each matching invocation runs independently, and background hooks can finish
>   in a different order than they started.
> - When the session ends, Codex cancels unfinished background hooks and discards
>   output that hasn't been delivered.
> - `SessionEnd` hooks always run synchronously.
