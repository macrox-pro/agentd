---
primary_sources:
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "Review and trust hooks; Config shape"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Review, trust, and config shape

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex Hooks — Review and trust hooks / Config shape

> ## Review and trust hooks
>
> Codex lists configured hooks before deciding which ones can run. Before a
> non-managed hook can run, Codex requires you to review and trust the exact hook
> definition. Codex records trust against the hook's current hash, so new or
> changed hooks are marked for review and skipped until trusted.
>
> Use `/hooks` in the CLI to inspect hook sources, review new or changed hooks,
> trust hooks, or disable individual non-managed hooks. If hooks need review at
> startup, Codex prints a warning that tells you to open `/hooks`.
>
> Managed hooks from system, MDM, cloud, or `requirements.toml` sources are marked
> as managed, trusted by policy, and can't be disabled from the user hook browser.
>
> For one-off automation that already vets hook sources outside Codex, pass
> `--dangerously-bypass-hook-trust` to run enabled hooks without requiring
> persisted hook trust for that invocation.
>
> ## Config shape
>
> Hooks are organized in three levels:
>
> - A hook event such as `PreToolUse`, `PostToolUse`, `PreCompact`,
>   `SubagentStart`, or `Stop`
> - A matcher group that decides when that event matches
> - One or more hook handlers that run when the matcher group matches
>
> ```json
> {
>   "description": "Optional lifecycle hooks for this workspace.",
>   "hooks": {
>     "SessionStart": [
>       {
>         "matcher": "startup|resume",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "python3 ~/.codex/hooks/session_start.py",
>             "statusMessage": "Loading session notes",
>             "additionalContextLimit": 5000
>           }
>         ]
>       }
>     ],
>     "SessionEnd": [
>       {
>         "hooks": [
>           {
>             "type": "command",
>             "command": "python3 ~/.codex/hooks/session_end.py",
>             "timeout": 3
>           }
>         ]
>       }
>     ],
>     "PreToolUse": [
>       {
>         "matcher": "Bash",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/usr/bin/python3 \"$(git rev-parse --show-toplevel)/.codex/hooks/pre_tool_use_policy.py\"",
>             "statusMessage": "Checking Bash command"
>           }
>         ]
>       }
>     ],
>     "PermissionRequest": [
>       {
>         "matcher": "Bash",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/usr/bin/python3 \"$(git rev-parse --show-toplevel)/.codex/hooks/permission_request.py\"",
>             "statusMessage": "Checking approval request"
>           }
>         ]
>       }
>     ],
>     "PostToolUse": [
>       {
>         "matcher": "Bash",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/usr/bin/python3 \"$(git rev-parse --show-toplevel)/.codex/hooks/post_tool_use_review.py\"",
>             "statusMessage": "Reviewing Bash output"
>           }
>         ]
>       }
>     ],
>     "UserPromptSubmit": [
>       {
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/usr/bin/python3 \"$(git rev-parse --show-toplevel)/.codex/hooks/user_prompt_submit_data_flywheel.py\""
>           }
>         ]
>       }
>     ],
>     "Stop": [
>       {
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/usr/bin/python3 \"$(git rev-parse --show-toplevel)/.codex/hooks/stop_continue.py\"",
>             "timeout": 30
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> Notes:
>
> - `description` is optional top-level metadata for a `hooks.json` file. It
>   doesn't change which hooks run.
> - `timeout` is in seconds.
> - If `timeout` is omitted, Codex uses `600` seconds for most hooks.
>   - `SessionEnd` uses `1` second by default and supports up to `3` seconds.
> - `statusMessage` is optional.
> - `additionalContextLimit` sets how much `additionalContext` a command hook can
>   send to the model before Codex saves the full text to disk and sends a shorter
>   preview instead. See [Large hook output](#large-hook-output).
> - `commandWindows` is an optional Windows-only command override. In TOML, use
>   `command_windows` or `commandWindows`.
> - Set `async` to `true` to [run a command hook in the
>   background](#run-hooks-in-the-background).
> - `command` and `mcp_tool` handlers are supported. `prompt` and `agent`
>   handlers are parsed but skipped.
> - Commands run with the session `cwd` as their working directory.
> - For repo-local hooks, prefer resolving from the git root instead of using a
>   relative path such as `.codex/hooks/...`. Codex may be started from a
>   subdirectory, and a git-root-based path keeps the hook location stable.
>
> Equivalent inline TOML in `config.toml`:
>
> ```toml
> [[hooks.SessionStart]]
> matcher = "^compact$"
>
> [[hooks.SessionStart.hooks]]
> type = "command"
> command = '/usr/bin/python3 "$(git rev-parse --show-toplevel)/.codex/hooks/session_start.py"'
> additionalContextLimit = 5000
>
> [[hooks.PreToolUse]]
> matcher = "^Bash$"
>
> [[hooks.PreToolUse.hooks]]
> type = "command"
> command = '/usr/bin/python3 "$(git rev-parse --show-toplevel)/.codex/hooks/pre_tool_use_policy.py"'
> timeout = 30
> statusMessage = "Checking Bash command"
>
> [[hooks.PostToolUse]]
> matcher = "^Bash$"
>
> [[hooks.PostToolUse.hooks]]
> type = "command"
> command = '/usr/bin/python3 "$(git rev-parse --show-toplevel)/.codex/hooks/post_tool_use_review.py"'
> timeout = 30
> statusMessage = "Reviewing Bash output"
> ```
