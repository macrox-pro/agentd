---
primary_sources:
  - id: T1-DEBUG
    title: "Debug your configuration"
    url: "https://code.claude.com/docs/en/debug-your-config.md"
    section: ""
  - id: T1-HOOKS-GUIDE
    title: "Hooks guide"
    url: "https://code.claude.com/docs/en/hooks-guide.md"
    section: "Troubleshooting"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook debug and troubleshooting

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

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

### Source: Hooks reference — Debug hooks

> ## Debug hooks
>
> Hook execution details, including which hooks matched, their exit codes, and full stdout and stderr, are written to the debug log file. Start Claude Code with `claude --debug-file <path>` to write the log to a known location, or run `claude --debug` and read the log at `~/.claude/debug/<session-id>.txt`. The `--debug` flag doesn't print to the terminal.
>
> For example, a `PostToolUse` hook on `Write` whose command prints `hook-ran` produces entries like:
>
> ```text
> 2026-07-19T02:03:24.382Z [DEBUG] Hook output does not start with {, treating as plain text
> 2026-07-19T02:03:24.382Z [DEBUG] Hook PostToolUse:Write (PostToolUse) success:
> hook-ran
> ```
>
> For more granular hook matching details, set `CLAUDE_CODE_DEBUG_LOG_LEVEL=verbose` to see additional log lines such as hook matcher counts and query matching.
>
> For troubleshooting common issues like hooks not firing, Stop hooks that keep blocking, or configuration errors, see [Limitations and troubleshooting](/docs/en/hooks-guide#limitations-and-troubleshooting) in the guide. For a broader diagnostic walkthrough covering `/context`, `/doctor`, and settings precedence, see [Debug your config](/docs/en/debug-your-config).

### Source: Hooks guide — Limitations and troubleshooting

> ## Limitations and troubleshooting
>
> ### Limitations
>
> Keep these constraints in mind when designing hooks:
>
> * Command hooks communicate through stdout, stderr, and exit codes only. They can't trigger `/` commands or tool calls. Text returned via `additionalContext` is injected as a system reminder that Claude reads as plain text. HTTP hooks communicate through the response body instead.
> * Hook timeouts vary by type. Override per hook with the `timeout` field in seconds.
>   * `command`, `http`, `mcp_tool`: 10 minutes. `UserPromptSubmit` lowers these to 30 seconds, and `MessageDisplay` lowers them to 10 seconds.
>   * `prompt`: 30 seconds.
>   * `agent`: 60 seconds.
>   * [`SessionEnd`](/docs/en/hooks#sessionend) hooks of any type share a 1.5-second budget. If your settings set a longer per-hook `timeout`, Claude Code raises the budget to match, up to 60 seconds.
> * `PostToolUse` hooks can't undo actions since the tool has already executed.
> * `PermissionRequest` hooks fire when Claude Code is about to ask you for permission.
>   * In [non-interactive mode](/docs/en/headless) with the `-p` flag, that prompt only exists when the Agent SDK's [`canUseTool` callback](/docs/en/agent-sdk/permissions) supplies it. In plain `-p` runs or with `--permission-prompt-tool`, use `PreToolUse` hooks for automated permission decisions instead.
>   * Background subagents can't show a prompt in non-interactive mode. Claude Code still runs the hooks for their tool calls, and if no hook returns a decision, it denies the call. In an interactive session, background subagent prompts surface in your main session and the hooks fire as usual.
> * `Stop` hooks fire whenever Claude finishes responding, not only at task completion. They don't fire on user interrupts. API errors fire [StopFailure](/docs/en/hooks#stopfailure) instead.
> * When multiple `PreToolUse` hooks return [`updatedInput`](/docs/en/hooks#pretooluse) to rewrite a tool's arguments, the last one to finish takes effect. Since hooks run in parallel, the order is non-deterministic. Avoid having more than one hook modify the same tool's input.
>
> ### Hooks and permission modes
>
> `PreToolUse` hooks fire before any permission-mode check, in every [permission mode](/docs/en/permission-modes), including `dontAsk`. A hook that returns `permissionDecision: "deny"` blocks the tool even in `bypassPermissions` mode or with `--dangerously-skip-permissions`. This lets you enforce policy that users can't bypass by changing their permission mode.
>
> The reverse is not true: a hook returning `"allow"` doesn't bypass deny rules from settings, and it can't suppress the prompt for connector tools [your organization set to `ask`](/docs/en/mcp#organization-controls-on-connector-tools) or MCP tools marked [`requiresUserInteraction`](/docs/en/mcp#require-approval-for-a-specific-tool). Hooks can tighten restrictions but not loosen them past what permission rules allow.
>
> ### Hook not firing
>
> The hook is configured but never executes.
>
> * Run `/hooks` and confirm the hook appears under the correct event
> * Check that the matcher pattern matches the tool name exactly. Matchers are case-sensitive
> * Verify you're triggering the right event type: `PreToolUse` fires before tool execution, `PostToolUse` fires after. A `PermissionRequest` hook fires when Claude Code is about to ask you for permission; see the [limitations](#limitations) for the non-interactive cases
>
> ### Hook error in output
>
> You see a message like "PreToolUse hook error: ..." in the transcript.
>
> * Your script exited with a non-zero code unexpectedly. Test it manually by piping sample JSON:
>   ```bash
>   echo '{"tool_name":"Bash","tool_input":{"command":"ls"}}' | ./my-hook.sh
>   echo $?  # Check the exit code
>   ```
> * If you see "command not found", use absolute paths or `${CLAUDE_PROJECT_DIR}` to reference scripts. To avoid shell quoting entirely, add `"args": []` to switch to [exec form](/docs/en/hooks#exec-form-and-shell-form), which spawns the script directly without a shell
> * If you see "jq: command not found", install `jq` or use Python/Node.js for JSON parsing
> * If the notice shows a JSON validation message, your hook's stdout parsed as JSON but failed schema validation. This happens even on exit 0. The reference's [Exit code output](/docs/en/hooks#exit-code-0) section covers the exit-code and JSON combinations
> * If the script isn't running at all, make it executable: `chmod +x ./my-hook.sh`
>
> ### `/hooks` shows no hooks configured
>
> You edited a settings file but the hooks don't appear in the menu.
>
> * File edits are normally picked up automatically. If they haven't appeared after a few seconds, the file watcher may have missed the change: restart your session to force a reload.
> * Verify your JSON is valid: trailing commas and comments aren't allowed
> * Confirm the settings file is in the correct location: `.claude/settings.json` for project hooks, `~/.claude/settings.json` for global hooks
>
> ### Stop hook hits the block cap
>
> Claude keeps working instead of stopping, then ends the turn with a warning that the Stop hook blocked too many consecutive times.
>
> Claude Code overrides a Stop hook after it blocks eight times in a row without progress. Your hook script needs to check whether it already triggered a continuation. Parse the `stop_hook_active` field from the JSON input and exit early if it's `true`:
>
> ```bash
> #!/bin/bash
> INPUT=$(cat)
> if [ "$(echo "$INPUT" | jq -r '.stop_hook_active')" = "true" ]; then
>   exit 0  # Allow Claude to stop
> fi
> # ... rest of your hook logic
> ```
>
> If your hook legitimately needs more than eight iterations to converge, raise the cap with [`CLAUDE_CODE_STOP_HOOK_BLOCK_CAP`](/docs/en/env-vars).
>
> ### Hook JSON has no effect
>
> Your hook prints valid JSON, but the decision doesn't take effect and no error appears in the transcript.
>
> When Claude Code runs a shell-form command hook, one without `args`, it spawns `sh -c` on macOS and Linux, Git Bash on Windows, or PowerShell when Git Bash isn't installed by default. This shell is non-interactive, but Git Bash and some configurations, such as `BASH_ENV` pointing at `~/.bashrc`, still source your profile. If that profile contains unconditional `echo` statements, the output gets prepended to your hook's JSON:
>
> ```text
> Shell ready on arm64
> {"decision": "block", "reason": "Not allowed"}
> ```
>
> The combined output no longer starts with `{`, so Claude Code treats all of stdout as plain text and ignores the JSON. On exit 0 nothing is reported in the transcript; the parse attempt is recorded only in the [debug log](/docs/en/hooks#debug-hooks). To fix this, wrap echo statements in your shell profile so they only run in interactive shells:
>
> ```bash
> # In ~/.zshrc or ~/.bashrc
> if [[ $- == *i* ]]; then
>   echo "Shell ready"
> fi
> ```
>
> The `$-` variable contains shell flags, and `i` means interactive. Hooks run in non-interactive shells, so the echo is skipped.
>
> ### Debug techniques
>
> Press `Ctrl+O` to open the transcript view to check the outcome of a hook run:
>
> * **Successful run**: you see nothing, unless the hook's JSON surfaces something, such as `systemMessage` or Stop hook feedback.
>   * To confirm a hook ran, check for its effect, like a reformatted file, or turn on debug logging as described below and trigger the hook again
> * **Blocking error**: on most events you see the hook's feedback. When the hook's JSON made a blocking decision, the feedback is the reason from that decision; otherwise it is the hook's stderr. On a few events, such as `ConfigChange` and `Elicitation`, a block surfaces no message.
> * **Non-blocking error**: the action proceeded, and you see a `<hook name> hook error` notice with a short explanation, such as the first line of stderr prefixed with `Failed with non-blocking status code:` or a JSON validation message.
>
> Which exit-code and JSON combinations produce each outcome, including the per-event exceptions, is defined in the reference's [Exit code output](/docs/en/hooks#exit-code-output) section.
>
> For full execution details including which hooks matched, their exit codes, stdout, and stderr, read the debug log. Start Claude Code with `claude --debug-file /tmp/claude.log` to write to a known path, then `tail -f /tmp/claude.log` in another terminal. If you started without that flag, run `/debug` mid-session to enable logging and find the log path.
