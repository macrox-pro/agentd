---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Hook events — session and setup"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook events — session and setup

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — Hook events

> ## Hook events
>
> Each event corresponds to a point in Claude Code's lifecycle where hooks can run. The sections below are ordered to match the lifecycle: from session setup through the agentic loop to session end. Each section describes when the event fires, what matchers it supports, the JSON input it receives, and how to control behavior through output.

### Source: Hooks reference — SessionStart

> ### SessionStart
>
> Runs when Claude Code starts a new session or resumes an existing session. Useful for loading development context like existing issues or recent changes to your codebase, or setting up environment variables. For static context that doesn't require a script, use [CLAUDE.md](/docs/en/memory) instead.
>
> SessionStart runs on every session, so keep these hooks fast. Only `type: "command"` and `type: "mcp_tool"` hooks are supported.
>
> The matcher value corresponds to how the session was initiated:
>
> | Matcher   | When it fires                                                                                                                          |
> | :-------- | :------------------------------------------------------------------------------------------------------------------------------------- |
> | `startup` | New session                                                                                                                            |
> | `resume`  | `--resume`, `--continue`, or `/resume`                                                                                                 |
> | `clear`   | `/clear`                                                                                                                               |
> | `compact` | Auto or manual compaction                                                                                                              |
> | `fork`    | A new session forked from an existing one: `--fork-session` with `--resume` or `--continue`, the `/fork` background copy, or `/branch` |
>
> Before v2.1.214, forked sessions reported source `"resume"`.
>
> #### SessionStart input
>
> In addition to the [common input fields](#common-input-fields), SessionStart hooks receive `source` and optionally `model`, `agent_type`, and `session_title`:
>
> | Field           | Description                                                                                                                                                                                                   |
> | :-------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `source`        | How the session started: `"startup"` for new sessions, `"resume"` for resumed sessions, `"clear"` after `/clear`, `"compact"` after compaction, or `"fork"` for a new session forked from an existing one     |
> | `model`         | The active model identifier. It can be omitted, for example after `/clear` or when a session is restored through conversation recovery, so check for the field before reading it                              |
> | `agent_type`    | The agent name, present when you start Claude Code with `claude --agent <name>`                                                                                                                               |
> | `session_title` | The current session title if one is already set, for example via `--name` or `/rename`. A hook that emits `sessionTitle` can check `session_title` first to avoid overwriting a title the user set explicitly |
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "hook_event_name": "SessionStart",
>   "source": "startup",
>   "model": "claude-sonnet-5"
> }
> ```
>
> #### SessionStart decision control
>
> Claude Code adds stdout it [treats as plain text](#exit-code-0) to Claude's context. In addition to the [JSON output fields](#json-output) available to all hooks, you can return these event-specific fields:
>
> | Field                | Description                                                                                                                                                                                                                                                                                                                          |
> | :------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `additionalContext`  | String added to Claude's context at the start of the conversation, before the first prompt. See [Add context for Claude](#add-context-for-claude) for how the text is delivered and what to put in it                                                                                                                                |
> | `initialUserMessage` | String used as the first user message of the session. Applies in [non-interactive mode](/docs/en/headless) with the `-p` flag, where it becomes the first turn even if no prompt is provided. If a prompt is provided, it follows as the next turn. Unlike `additionalContext`, which attaches to an existing turn, this creates the turn |
> | `sessionTitle`       | Sets the session title, with the same effect as `/rename`. Use to name sessions automatically from the launch folder, git branch, or worktree name. Applies when `source` is `"startup"`, `"resume"`, or `"fork"`; ignored on `"clear"` and `"compact"`                                                                              |
> | `watchPaths`         | Array of absolute paths to watch for [FileChanged](#filechanged) events during this session                                                                                                                                                                                                                                          |
> | `reloadSkills`       | Boolean. When `true`, Claude Code re-scans the [skill](/docs/en/skills) and command directories after the SessionStart hooks complete, so skills the hook installed are available in the same session, starting with the first prompt                                                                                                     |
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "SessionStart",
>     "additionalContext": "Current branch: feat/auth-refactor\nUncommitted changes: src/auth.ts, src/login.tsx\nActive issue: #4211 Migrate to OAuth2",
>     "sessionTitle": "auth-refactor"
>   }
> }
> ```
>
> Since plain stdout already reaches Claude for this event, a hook that only loads context can print to stdout directly without building JSON. Use the JSON form when you need to combine context with other fields such as `sessionTitle`.
>
> Use `reloadSkills` when a SessionStart hook installs or updates skills. Skill discovery normally runs before SessionStart hooks finish, so files the hook writes into `~/.claude/skills/` or `.claude/skills/` would otherwise only appear in the next session. This example syncs a shared skills repository and requests the re-scan:
>
> ```bash
> #!/bin/bash
>
> git -C ~/.claude/skills/team-skills pull --quiet 2>/dev/null || \
>   git clone --quiet https://git.example.com/your-org/team-skills.git ~/.claude/skills/team-skills
>
> echo '{"hookSpecificOutput": {"hookEventName": "SessionStart", "reloadSkills": true}}'
> ```
>
> The repository URL is a placeholder; replace it with your own skills repository. With the placeholder, the clone fails and prints a `fatal:` message to stderr. Stderr from a SessionStart hook that exits 0 is informational only, so the `reloadSkills` request still applies.
>
> #### Persist environment variables
>
> SessionStart hooks have access to the `CLAUDE_ENV_FILE` environment variable, which provides a file path where you can persist environment variables for subsequent Bash commands.
>
> To set individual environment variables, write `export` statements to `CLAUDE_ENV_FILE`. Use append (`>>`) to preserve variables set by other hooks:
>
> ```bash
> #!/bin/bash
>
> if [ -n "$CLAUDE_ENV_FILE" ]; then
>   echo 'export NODE_ENV=production' >> "$CLAUDE_ENV_FILE"
>   echo 'export DEBUG_LOG=true' >> "$CLAUDE_ENV_FILE"
>   echo 'export PATH="$PATH:./node_modules/.bin"' >> "$CLAUDE_ENV_FILE"
> fi
>
> exit 0
> ```
>
> To capture all environment changes from setup commands, compare the exported variables before and after:
>
> ```bash
> #!/bin/bash
>
> ENV_BEFORE=$(export -p | sort)
>
> # Run your setup commands that modify the environment
> source ~/.nvm/nvm.sh
> nvm use 20
>
> if [ -n "$CLAUDE_ENV_FILE" ]; then
>   ENV_AFTER=$(export -p | sort)
>   comm -13 <(echo "$ENV_BEFORE") <(echo "$ENV_AFTER") >> "$CLAUDE_ENV_FILE"
> fi
>
> exit 0
> ```
>
>
>   `CLAUDE_ENV_FILE` is available for SessionStart, [Setup](#setup), [CwdChanged](#cwdchanged), and [FileChanged](#filechanged) hooks. Other hook types don't have access to this variable.

### Source: Hooks reference — SessionEnd

> ### SessionEnd
>
> Runs when a Claude Code session ends. Useful for cleanup tasks, logging session
> statistics, or saving session state. Supports matchers to filter by exit reason.
>
> The `reason` field in the hook input indicates why the session ended:
>
> | Reason                        | Description                                                                               |
> | :---------------------------- | :---------------------------------------------------------------------------------------- |
> | `clear`                       | Session cleared with `/clear` command                                                     |
> | `resume`                      | Session switched via interactive `/resume`                                                |
> | `logout`                      | User logged out                                                                           |
> | `prompt_input_exit`           | User exited while prompt input was visible                                                |
> | `other`                       | Other exit reasons                                                                        |
> | `bypass_permissions_disabled` | Removed in v2.1.234; Claude Code doesn't send it. Drop it from your `SessionEnd` matchers |
>
> #### SessionEnd input
>
> In addition to the [common input fields](#common-input-fields), SessionEnd hooks receive a `reason` field indicating why the session ended. See the [reason table](#sessionend) above for all values.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "hook_event_name": "SessionEnd",
>   "reason": "other"
> }
> ```
>
> SessionEnd hooks have no decision control. They can't block session termination but can perform cleanup tasks. Claude Code discards their [JSON output fields](#json-output), such as `systemMessage`.
>
> SessionEnd hooks have a default timeout of 1.5 seconds. This applies to session exit, `/clear`, and switching sessions via interactive `/resume`. If a hook needs more time, set a per-hook `timeout` in the hook configuration. The overall budget is automatically raised to the highest per-hook timeout configured in settings files, up to 60 seconds. Timeouts set on plugin-provided hooks don't raise the budget. To override the budget explicitly, set the `CLAUDE_CODE_SESSIONEND_HOOKS_TIMEOUT_MS` environment variable in milliseconds.
>
> ```bash
> CLAUDE_CODE_SESSIONEND_HOOKS_TIMEOUT_MS=5000 claude
> ```

### Source: Hooks reference — Setup

> ### Setup
>
> Fires only when you launch Claude Code with `--init-only`, or with `--init` or `--maintenance` in [non-interactive mode](/docs/en/headless) with the `-p` flag. It doesn't fire on normal startup. Use it for one-time dependency installation or scheduled cleanup that you trigger explicitly from CI or scripts, separate from normal session startup. For per-session initialization, use [SessionStart](#sessionstart) instead.
>
> The matcher value corresponds to the CLI flag that triggered the hook:
>
> | Matcher       | When it fires                              |
> | :------------ | :----------------------------------------- |
> | `init`        | `claude --init-only` or `claude -p --init` |
> | `maintenance` | `claude -p --maintenance`                  |
>
> When you run `claude --init-only`, Claude Code runs Setup hooks and `SessionStart` hooks with the `startup` matcher, then exits without starting a conversation.
>
> When you start or continue a conversation with `-p`, you also need to supply a prompt, as an argument or piped on stdin. You can skip the prompt when a `SessionStart` hook supplies [`initialUserMessage`](#sessionstart-decision-control) or when you resume a session with a [deferred tool call](#defer-a-tool-call-for-later).
>
> On success, `--init-only` prints nothing to the terminal. To confirm the hooks ran, start with `claude --debug-file <path> --init-only`, replacing `<path>` with a log file location, and check the log for the Setup and SessionStart hook entries.
>
> Because Setup doesn't fire on every launch, a plugin that needs a dependency installed can't rely on Setup alone. The practical pattern is to check for the dependency on first use and install on miss, for example a hook or skill that tests for `${CLAUDE_PLUGIN_DATA}/node_modules` and runs `npm install` if absent. See the [persistent data directory](/docs/en/plugins-reference#persistent-data-directory) for where to store installed dependencies. If you distribute your plugin through a marketplace, you may not need this pattern: Claude Code [installs eligible Node.js package dependencies automatically](/docs/en/plugins-reference#node-js-package-dependencies) when it caches the plugin.
>
> #### Setup input
>
> In addition to the [common input fields](#common-input-fields), Setup hooks receive a `trigger` field set to either `"init"` or `"maintenance"`:
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "hook_event_name": "Setup",
>   "trigger": "init"
> }
> ```
>
> #### Setup decision control
>
> Setup hooks can't block; execution continues on any exit code. Exit code 2 surfaces stderr to the user as a `<hook name> hook error` notice whether or not you print JSON. On any other non-zero exit, Claude Code applies the [other exit codes](#other-exit-codes) rule: with a parsed object that passes schema validation it honors the fields and doesn't report the hook as an error, and with anything else on stdout it shows that notice. In [non-interactive mode](/docs/en/headless), hook output appears only when you launch with `--verbose`.
>
> To pass information into Claude's context, return `additionalContext` in JSON output; plain stdout is written to the debug log only. In addition to the [JSON output fields](#json-output) available to all hooks, you can return these event-specific fields:
>
> | Field               | Description                                                               |
> | :------------------ | :------------------------------------------------------------------------ |
> | `additionalContext` | String added to Claude's context. Multiple hooks' values are concatenated |
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "Setup",
>     "additionalContext": "Dependencies installed: node_modules, .venv"
>   }
> }
> ```
>
> Setup hooks have access to `CLAUDE_ENV_FILE`. Variables written to that file persist into subsequent Bash commands for the session, just as in [SessionStart hooks](#persist-environment-variables). Only `type: "command"` and `type: "mcp_tool"` hooks are supported.

### Source: Hooks reference — ConfigChange

> ### ConfigChange
>
> Runs when a configuration file changes during a session. Use this to audit settings changes, enforce security policies, or block unauthorized modifications to configuration files.
>
> Claude Code runs ConfigChange hooks when a settings file, a managed policy file, or a skill file changes. For managed policy, it runs them only when `managed-settings.json` or a file in `managed-settings.d/` changes. It applies [server-managed settings](/docs/en/server-managed-settings) and changes to macOS managed preferences or Windows registry policy without running them. On WSL with [`wslInheritsWindowsSettings`](/docs/en/settings#available-settings), it also applies a changed Windows-side managed settings file on its policy poll without running them.
>
> The matcher filters on the configuration source:
>
> | Matcher            | When it fires                                                      |
> | :----------------- | :----------------------------------------------------------------- |
> | `user_settings`    | `~/.claude/settings.json` changes                                  |
> | `project_settings` | `.claude/settings.json` changes                                    |
> | `local_settings`   | `.claude/settings.local.json` changes                              |
> | `policy_settings`  | `managed-settings.json` or a file in `managed-settings.d/` changes |
> | `skills`           | A skill file in `.claude/skills/` changes                          |
>
> This example logs all configuration changes for security auditing:
>
> ```json
> {
>   "hooks": {
>     "ConfigChange": [
>       {
>         "hooks": [
>           {
>             "type": "command",
>             "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/audit-config-change.sh",
>             "args": []
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> #### ConfigChange input
>
> In addition to the [common input fields](#common-input-fields), ConfigChange hooks receive `source` and optionally `file_path`. The `source` field indicates which configuration type changed, and `file_path` provides the path to the specific file that was modified.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "hook_event_name": "ConfigChange",
>   "source": "project_settings",
>   "file_path": "/Users/.../my-project/.claude/settings.json"
> }
> ```
>
> #### ConfigChange decision control
>
> ConfigChange hooks can block configuration changes from taking effect. Use exit code 2 or a JSON `decision` to prevent the change. When blocked, the new settings are not applied to the running session.
>
> | Field      | Description                                                                              |
> | :--------- | :--------------------------------------------------------------------------------------- |
> | `decision` | `"block"` prevents the configuration change from being applied. Omit to allow the change |
> | `reason`   | Accepted but never shown                                                                 |
>
> ```json
> {
>   "decision": "block",
>   "reason": "Configuration changes to project settings require admin approval"
> }
> ```
>
> `policy_settings` changes can't be blocked. Hooks still fire for `policy_settings` sources when a managed settings file on the machine changes, so you can use them to log those edits, but any blocking decision is ignored. This ensures enterprise-managed settings always take effect. Claude Code doesn't run `ConfigChange` hooks when [server-managed settings](/docs/en/server-managed-settings) arrive or refresh.
>
> Claude Code acts on the blocking decision from a ConfigChange hook's JSON output and discards `systemMessage` and `continue`. A blocked change surfaces no message to you or to Claude, whether you block with `reason` or with stderr on exit 2. Claude Code only writes a line to the debug log.

### Source: Hooks reference — CwdChanged

> ### CwdChanged
>
> Runs when the working directory changes during a session, for example when Claude executes a `cd` command. Use this to react to directory changes: reload environment variables, activate project-specific toolchains, or run setup scripts automatically. Pairs with [FileChanged](#filechanged) for tools like [direnv](https://direnv.net/) that manage per-directory environment.
>
> CwdChanged hooks have access to `CLAUDE_ENV_FILE`. Variables written to that file persist into subsequent Bash commands for the session, just as in [SessionStart hooks](#persist-environment-variables).
>
> CwdChanged doesn't support matchers and fires on every directory change.
>
> #### CwdChanged input
>
> In addition to the [common input fields](#common-input-fields), CwdChanged hooks receive `old_cwd` and `new_cwd`.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../transcript.jsonl",
>   "cwd": "/Users/my-project/src",
>   "hook_event_name": "CwdChanged",
>   "old_cwd": "/Users/my-project",
>   "new_cwd": "/Users/my-project/src"
> }
> ```
>
> #### CwdChanged output
>
> In addition to the [JSON output fields](#json-output) available to all hooks, CwdChanged hooks can return `watchPaths` to dynamically set which file paths [FileChanged](#filechanged) watches:
>
> | Field        | Description                                                                                                                                                                                                                    |
> | :----------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `watchPaths` | Array of absolute paths. Replaces the current dynamic watch list. Paths from your `matcher` configuration are always watched. Returning an empty array clears the dynamic list, which is typical when entering a new directory |
>
> CwdChanged hooks have no decision control. They can't block the directory change.
>
> Claude Code reads `watchPaths` and `systemMessage` from their JSON output and discards `continue`. In interactive sessions, it shows the `systemMessage` as a brief terminal notification. The message doesn't reach the SDK message stream.

### Source: Hooks reference — DirectoryAdded

> ### DirectoryAdded
>
> Runs after you add a working directory mid-session with the `/add-dir` command, or after an SDK client adds one with the `register_repo_root` control request. Use this to prepare a newly added repository, for example by installing its dependencies.
>
> Claude Code doesn't fire this event when:
>
> * You pass a directory with the `--add-dir` startup flag; [SessionStart](#sessionstart) covers those directories
> * You add a directory on the `/permissions` Workspace tab
> * You add a directory that is already a working directory; the add fails with an error
>
> Claude Code fires DirectoryAdded after refreshing sandbox and permission state, so sandboxed tools already see the new directory when your hook runs. Hook commands themselves run unsandboxed.
>
> Claude Code doesn't wait for the hook: the add completes immediately, and the hook runs in the background with the 600-second default timeout.
>
> The matcher filters on how the directory was added:
>
> | Matcher              | When it fires                                                                |
> | :------------------- | :--------------------------------------------------------------------------- |
> | `slash_command`      | You add a directory with `/add-dir`                                          |
> | `register_repo_root` | An SDK client adds a directory with the `register_repo_root` control request |
>
> #### DirectoryAdded input
>
> In addition to the [common input fields](#common-input-fields), DirectoryAdded hooks receive `directory` and `source`.
>
> | Field       | Description                                                                                                         |
> | :---------- | :------------------------------------------------------------------------------------------------------------------ |
> | `directory` | Absolute path of the directory that was added                                                                       |
> | `source`    | How the directory was added, `"slash_command"` for `/add-dir` or `"register_repo_root"` for the SDK control request |
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../transcript.jsonl",
>   "cwd": "/Users/my-project",
>   "hook_event_name": "DirectoryAdded",
>   "directory": "/Users/my-other-repo",
>   "source": "slash_command"
> }
> ```
>
> DirectoryAdded hooks have no decision control. They can't block the add, which has already completed when the hook runs. Claude Code discards the `continue` field from their JSON output and surfaces the rest differently per source:
>
> * `slash_command`: Claude Code delivers the hook's `systemMessage` to Claude as context on the next conversation turn, rather than showing it to you. A count of failed hooks appears in the transcript. Full failure output goes to the debug log
> * `register_repo_root`: Claude Code writes `systemMessage` output and failure output to the debug log only
