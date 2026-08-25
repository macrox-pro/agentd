---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Matcher patterns"
  - id: T1-DEBUG
    title: "Debug your configuration"
    url: "https://code.claude.com/docs/en/debug-your-config.md"
    section: "Check hooks"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Matcher patterns

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — Matcher patterns

> ### Matcher patterns
>
> The `matcher` field filters when hooks fire. How a matcher is evaluated depends on the characters it contains:
>
> | Matcher value                                         | Evaluated as                                                                                         | Example                                                                                                                                         |
> | :---------------------------------------------------- | :--------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------- |
> | `"*"`, `""`, or omitted                               | Match all                                                                                            | fires on every occurrence of the event                                                                                                          |
> | Only letters, digits, `_`, `-`, spaces, `,`, and `\|` | Exact string, or list of exact strings separated by `\|` or `,` with optional surrounding whitespace | `Bash` matches only the Bash tool; `Edit\|Write` and `Edit, Write` each match either tool exactly; `code-reviewer` matches only that agent type |
> | Contains any other character                          | JavaScript regular expression, unanchored                                                            | `^Notebook` matches any tool whose name starts with `Notebook`; `mcp__memory__.*` matches every tool from the `memory` server                   |
>
> A matcher on the regular-expression path is tested with JavaScript's `RegExp.prototype.test`, which succeeds on a match anywhere in the value. `Edit.*` matches both `Edit` and `NotebookEdit`; wrap the pattern in `^` and `$`, as in `^Edit$`, when you need a whole-string match.
>
> Comma separators and the surrounding whitespace tolerance require Claude Code v2.1.191 or later.
>
> Hyphens in the exact-match set require Claude Code v2.1.195 or later. On earlier versions a hyphenated name like `code-reviewer` is evaluated as an unanchored regular expression, so it also fires for `senior-code-reviewer`; anchor it as `^code-reviewer$` on those versions to match only that name.
>
> `FileChanged` and `StopFailure` use a narrower exact-match set of letters, digits, `_`, and `|` only. A hyphen, space, or comma in a matcher for those two events keeps it on the regular-expression path, and only `|` separates alternatives. Every other event with matcher support in the table that follows accepts `|` or `,`.
>
> The `FileChanged` event doesn't follow these rules when building its watch list. See [FileChanged](#filechanged).
>
> Each event type matches on a different field:
>
> | Event                                                                                                                                             | What the matcher filters                                     | Example matcher values                                                                                                                                                                                                                                                         |
> | :------------------------------------------------------------------------------------------------------------------------------------------------ | :----------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `PreToolUse`, `PostToolUse`, `PostToolUseFailure`, `PermissionRequest`, `PermissionDenied`                                                        | tool name                                                    | `Bash`, `Edit\|Write`, `mcp__.*`                                                                                                                                                                                                                                               |
> | `SessionStart`                                                                                                                                    | how the session started                                      | `startup`, `resume`, `clear`, `compact`, `fork`                                                                                                                                                                                                                                |
> | `Setup`                                                                                                                                           | which CLI flag triggered setup                               | `init`, `maintenance`                                                                                                                                                                                                                                                          |
> | `SessionEnd`                                                                                                                                      | why the session ended                                        | `clear`, `resume`, `logout`, `prompt_input_exit`, `other`                                                                                                                                                                                                                      |
> | `Notification`                                                                                                                                    | notification type                                            | `permission_prompt`, `idle_prompt`, `auth_success`, `elicitation_dialog`, `elicitation_url_dialog`, `elicitation_complete`, `elicitation_response`, `agent_needs_input`, `agent_completed`, `quota_auto_resume_fired`, `quota_auto_resume_stale`, `quota_auto_resume_disabled` |
> | `SubagentStart`                                                                                                                                   | agent type                                                   | `general-purpose`, `Explore`, `Plan`, custom agent names, or plugin-scoped names like `^my-plugin:reviewer$`                                                                                                                                                                   |
> | `PreCompact`, `PostCompact`                                                                                                                       | what triggered compaction                                    | `manual`, `auto`                                                                                                                                                                                                                                                               |
> | `SubagentStop`                                                                                                                                    | agent type                                                   | same values as `SubagentStart`                                                                                                                                                                                                                                                 |
> | `ConfigChange`                                                                                                                                    | configuration source                                         | `user_settings`, `project_settings`, `local_settings`, `policy_settings`, `skills`                                                                                                                                                                                             |
> | `CwdChanged`                                                                                                                                      | no matcher support                                           | always fires on every directory change                                                                                                                                                                                                                                         |
> | `DirectoryAdded`                                                                                                                                  | how the directory was added                                  | `slash_command`, `register_repo_root`                                                                                                                                                                                                                                          |
> | `FileChanged`                                                                                                                                     | literal filenames to watch (see [FileChanged](#filechanged)) | `.envrc\|.env`                                                                                                                                                                                                                                                                 |
> | `StopFailure`                                                                                                                                     | error type                                                   | `rate_limit`, `overloaded`, `authentication_failed`, `oauth_org_not_allowed`, `billing_error`, `invalid_request`, `model_not_found`, `server_error`, `max_output_tokens`, `unknown`                                                                                            |
> | `InstructionsLoaded`                                                                                                                              | load reason                                                  | `session_start`, `nested_traversal`, `path_glob_match`, `include`, `compact`                                                                                                                                                                                                   |
> | `UserPromptExpansion`                                                                                                                             | command name                                                 | your skill or command names                                                                                                                                                                                                                                                    |
> | `Elicitation`                                                                                                                                     | MCP server name                                              | your configured MCP server names                                                                                                                                                                                                                                               |
> | `ElicitationResult`                                                                                                                               | MCP server name                                              | same values as `Elicitation`                                                                                                                                                                                                                                                   |
> | `UserPromptSubmit`, `PostToolBatch`, `Stop`, `TeammateIdle`, `TaskCreated`, `TaskCompleted`, `WorktreeCreate`, `WorktreeRemove`, `MessageDisplay` | no matcher support                                           | always fires on every occurrence                                                                                                                                                                                                                                               |
>
> The matcher runs against a field from the [JSON input](#hook-input-and-output) that Claude Code sends to your hook on stdin. For tool events, that field is `tool_name`. Each [hook event](#hook-events) section lists the full set of matcher values and the input schema for that event.
>
> This example runs a linting script only when Claude writes or edits a file:
>
> ```json
> {
>   "hooks": {
>     "PostToolUse": [
>       {
>         "matcher": "Edit|Write",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/path/to/lint-check.sh"
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> If you add a `matcher` field to an event without matcher support, it is silently ignored.
>
> For tool events, you can filter more narrowly by setting the [`if` field](#common-fields) on individual hook handlers. `if` uses [permission rule syntax](/docs/en/permissions) to match against the tool name and arguments together, so `"Bash(git *)"` runs when any subcommand of the Bash input matches `git *` and `"Edit(*.ts)"` runs only for TypeScript files.
>
> #### Match MCP tools
>
> [MCP](/docs/en/mcp) server tools appear as regular tools in tool events (`PreToolUse`, `PostToolUse`, `PostToolUseFailure`, `PermissionRequest`, `PermissionDenied`), so you can match them the same way you match any other tool name.
>
> MCP tools follow the naming pattern `mcp__<server>__<tool>`, for example:
>
> * `mcp__memory__create_entities`: Memory server's create entities tool
> * `mcp__filesystem__read_file`: Filesystem server's read file tool
> * `mcp__github__search_repositories`: GitHub server's search tool
>
> To match every tool from a server, append `.*` to the server prefix. The `.*` is required: a matcher like `mcp__memory` or `mcp__brave-search` contains only exact-match characters, so it is compared as an exact string and matches no tool.
>
> * `mcp__memory__.*` matches all tools from the `memory` server
> * `mcp__brave-search__.*` matches all tools from a server whose name contains a hyphen
> * `mcp__.*__write.*` matches any tool whose name starts with `write` from any server
>
> Hyphens in the exact-match set require Claude Code v2.1.195 or later. On earlier versions a bare hyphenated prefix like `mcp__brave-search` is evaluated as an unanchored regular expression and matches every tool from that server. The `mcp__brave-search__.*` form works on every version.
>
> Tools from a [plugin-bundled MCP server](/docs/en/mcp#plugin-provided-mcp-servers) use a scoped server segment that includes the plugin name: `mcp__plugin_<plugin-name>_<server-name>__<tool>`. A matcher written against the bare server key never fires for these tools. For a plugin named `my-plugin` that bundles a server under the key `db`, a `query` tool appears as `mcp__plugin_my-plugin_db__query`, so the matcher for every tool from that server is `mcp__plugin_my-plugin_db__.*`. Use the same scoped tool name in a handler's [`if` field](#common-fields). See [Plugin-provided MCP servers](/docs/en/mcp#plugin-provided-mcp-servers) for how the scoped name is built.
>
> This example logs all memory server operations and validates write operations from any MCP server:
>
> ```json
> {
>   "hooks": {
>     "PreToolUse": [
>       {
>         "matcher": "mcp__memory__.*",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "echo 'Memory operation initiated' >> ~/mcp-operations.log"
>           }
>         ]
>       },
>       {
>         "matcher": "mcp__.*__write.*",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/home/user/scripts/validate-mcp-write.py"
>           }
>         ]
>       }
>     ]
>   }
> }
> ```

### Source: Debug your configuration — Check hooks

> ## Check hooks
>
> Run `/hooks` to list every hook registered for the current session, grouped by event. If a hook you defined doesn't appear, it isn't being read: hooks go under the `"hooks"` key in a settings file, not in a standalone file.
>
> If the hook appears but doesn't fire, the matcher is the usual cause. Check it for these mistakes:
>
> * The `matcher` field is a single string that uses `|` to match multiple tool names, for example `"Edit|Write"`. A `,` separator is equivalent, so `"Edit,Write"` matches the same tools. Before v2.1.191, a comma fell through to regex evaluation and the matcher never matched, so use `|` if you aren't on v2.1.191 yet.
> * A misspelled tool name produces a matcher that matches nothing, so the hook fails silently.
> * An array value is a schema error: Claude Code shows a settings error notice and rejects the whole user, project, or local settings file, `claude doctor` reports the validation failure, and no hook from that file appears in `/hooks`. In [managed settings](/docs/en/managed-settings), Claude Code drops the whole `hooks` key from the file that contains the array, so none of that file's hooks apply. The file's other settings still apply, and `claude doctor` lists the dropped key.
>
> Edits to `settings.json` take effect in the running session after a brief file-stability delay. You don't need to restart. If `/hooks` still shows the old definition a few seconds after saving, run `/hooks` again to refresh the view.
>
> If `/hooks` shows the hook but it still does not fire, the next step is to watch hook evaluation live. Start a session with `claude --debug` and trigger the tool call. The debug log records each event, which matchers were checked, and the hook's exit code and output. See [Debug hooks](/docs/en/hooks#debug-hooks) for the log format and [hooks troubleshooting](/docs/en/hooks-guide#limitations-and-troubleshooting) for common failure patterns.
