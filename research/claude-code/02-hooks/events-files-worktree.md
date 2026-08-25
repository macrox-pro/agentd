---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Hook events — files, worktrees, instructions"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook events — files, worktrees, instructions

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — Hook events

> ## Hook events
>
> Each event corresponds to a point in Claude Code's lifecycle where hooks can run. The sections below are ordered to match the lifecycle: from session setup through the agentic loop to session end. Each section describes when the event fires, what matchers it supports, the JSON input it receives, and how to control behavior through output.

### Source: Hooks reference — FileChanged

> ### FileChanged
>
> Runs when a watched file changes on disk. Claude Code detects changes with a filesystem watcher, not by inspecting tool calls, so it runs the hook no matter what changed the file: an `Edit` or `Write` tool call, a script Claude runs with `Bash`, or a process outside Claude Code entirely. A common use is reloading environment variables when project configuration files change.
>
> The `matcher` for this event serves two roles:
>
> * **Build the watch list**: the value is split on `|` and each segment is registered as a literal filename in the working directory, so `".envrc|.env"` watches exactly those two files. Regex patterns are not useful here: a value like `^\.env` would watch a file literally named `^\.env`.
> * **Filter which hooks run**: when a watched file changes, the same value filters which hook groups run using the standard [matcher rules](#matcher-patterns) against the changed file's basename.
>
> This example normalizes line endings in `data.csv` after any change, including a `Bash` command or an external script rewriting the file:
>
> ```json
> {
>   "hooks": {
>     "FileChanged": [
>       {
>         "matcher": "data.csv",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "/path/to/normalize-line-endings.sh"
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> The hook reads the changed file's absolute path from the `file_path` field of the [JSON input](#filechanged-input) on stdin. Its `grep` guard tests for the same thing `perl` removes, a CR at the end of a line, so the run after a normalization exits without touching the file. A looser guard loops forever, because `perl -i` rewrites the file even when it substitutes nothing and Claude Code runs the hook again after every rewrite. Save this script at `/path/to/normalize-line-endings.sh` and make it executable:
>
> ```bash
> #!/bin/bash
> FILE=$(jq -r .file_path)
> if grep -q $'\r$' "$FILE"; then
>   perl -pi -e 's/\r$//' "$FILE"
> fi
> ```
>
> To confirm the hook works, ask Claude to append a CRLF line to `data.csv` with a `Bash` command. Claude Code runs the hook and the file ends up with LF endings.
>
> To watch files you can't name up front, return [`watchPaths`](#filechanged-output) from a hook to update the watch list dynamically. Claude Code starts the watcher only when something names a file to watch, so seed the list with a FileChanged group whose matcher names at least one file, or with a [SessionStart](#sessionstart-decision-control) or [CwdChanged](#cwdchanged) hook that returns `watchPaths`. The matcher still filters which hook groups run when a watched file changes, so give the group that handles dynamic paths an omitted matcher, which matches every watched file and adds nothing to the watch list. A `"*"` matcher also matches every file, but Claude Code registers it in the watch list like any other value, as a literal file named `*`.
>
> FileChanged hooks have access to `CLAUDE_ENV_FILE`. Variables written to that file persist into subsequent Bash commands for the session, just as in [SessionStart hooks](#persist-environment-variables).
>
> #### FileChanged input
>
> In addition to the [common input fields](#common-input-fields), FileChanged hooks receive `file_path` and `event`.
>
> | Field       | Description                                                                                                 |
> | :---------- | :---------------------------------------------------------------------------------------------------------- |
> | `file_path` | Absolute path to the file that changed                                                                      |
> | `event`     | What happened: `"change"` for a modified file, `"add"` for a created file, or `"unlink"` for a deleted file |
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../transcript.jsonl",
>   "cwd": "/Users/my-project",
>   "hook_event_name": "FileChanged",
>   "file_path": "/Users/my-project/.envrc",
>   "event": "change"
> }
> ```
>
> #### FileChanged output
>
> In addition to the [JSON output fields](#json-output) available to all hooks, FileChanged hooks can return `watchPaths` to dynamically update which file paths are watched:
>
> | Field        | Description                                                                                                                                                                                                                |
> | :----------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `watchPaths` | Array of absolute paths. Replaces the current dynamic watch list. Paths from your `matcher` configuration are always watched. Use this when your hook script discovers additional files to watch based on the changed file |
>
> FileChanged hooks have no decision control. They can't block the file change from occurring.
>
> Claude Code reads `watchPaths` and `systemMessage` from their JSON output and discards `continue`. In interactive sessions, it shows the `systemMessage` as a brief terminal notification. The message doesn't reach the SDK message stream.

### Source: Hooks reference — WorktreeCreate

> ### WorktreeCreate
>
> Runs when a worktree is being created, whether from `claude --worktree`, from a [subagent using `isolation: "worktree"`](/docs/en/sub-agents#choose-the-subagent-scope), or for a [background session](/docs/en/agent-view#how-file-edits-are-isolated) that Claude Code isolates in its own worktree. By default Claude Code creates the isolated working copy with `git worktree`. Configuring a WorktreeCreate hook replaces that default git behavior, letting you use a different version control system like SVN, Perforce, or Mercurial.
>
> Because the hook replaces the default behavior entirely, [`.worktreeinclude`](/docs/en/worktrees#copy-gitignored-files-into-worktrees) is not processed. If you need to copy local configuration files like `.env` into the new worktree, do it inside your hook script.
>
> The hook must return the path to the created worktree directory. Claude Code uses this path as the working directory for the isolated session. See [WorktreeCreate output](#worktreecreate-output) for how each hook type returns the path.
>
> Claude Code acts on the hook's success and the returned path, and discards `systemMessage` and `continue`.
>
> This example creates an SVN working copy and prints the path for Claude Code to use. Replace the repository URL with your own:
>
> ```json
> {
>   "hooks": {
>     "WorktreeCreate": [
>       {
>         "hooks": [
>           {
>             "type": "command",
>             "command": "bash -c 'NAME=$(jq -r .name); DIR=\"$HOME/.claude/worktrees/$NAME\"; svn checkout https://svn.example.com/repo/trunk \"$DIR\" >&2 && echo \"$DIR\"'"
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> The hook reads the worktree `name` from the JSON input on stdin, checks out a fresh copy into a new directory, and prints the directory path. The `echo` on the last line is what Claude Code reads as the worktree path. Redirect any other output to stderr so it doesn't interfere with the path.
>
> #### WorktreeCreate input
>
> In addition to the [common input fields](#common-input-fields), WorktreeCreate hooks receive the `name` field. This is a slug identifier for the new worktree, either specified by the user or auto-generated, for example `bold-oak-a3f2`.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "hook_event_name": "WorktreeCreate",
>   "name": "feature-auth"
> }
> ```
>
> #### WorktreeCreate output
>
> WorktreeCreate hooks don't use the standard allow/block decision model. Instead, the hook's success or failure determines the outcome. The hook must return the path to the created worktree directory:
>
> * **Command hooks** (`type: "command"`): print the path as the last non-empty line of stdout. Claude Code strips ANSI escape codes before reading that line, so shell startup banners printed before your `echo` are ignored. Redirect any other hook output to stderr.
> * **HTTP hooks** (`type: "http"`): return `{ "hookSpecificOutput": { "hookEventName": "WorktreeCreate", "worktreePath": "/absolute/path" } }` in the response body.
>
> If the hook fails or produces no path, worktree creation fails with an error.
>
> Claude Code resolves a relative path against the directory the hook ran in, collapsing any `.` or `..` segments in it. If the resulting path isn't a directory Claude Code can enter, the session prints an error naming the path and exits with code 1.
>
> Claude Code refuses an absolute path that contains `.` or `..` segments, and any path that passes through a symlink below the repository root, because a symlink committed to the repository could redirect the worktree outside it. The error names the rejected component. Return a normalized path that doesn't pass through a symlink inside the repository. Before v2.1.216, worktree creation followed the hook's path without this screening.

### Source: Hooks reference — WorktreeRemove

> ### WorktreeRemove
>
> Runs when a worktree is being removed. This is the cleanup counterpart to [WorktreeCreate](#worktreecreate). The event fires when:
>
> * you exit a `--worktree` session and choose to remove it
> * a subagent with `isolation: "worktree"` finishes
> * you delete a [background session](/docs/en/agent-view#what-deleting-a-session-removes) whose worktree the hook created
>
> For git-based worktrees, Claude Code handles cleanup automatically with `git worktree remove`. If you configured a WorktreeCreate hook for a non-git version control system, pair it with a WorktreeRemove hook to handle cleanup. Without one, the worktree directory is left on disk.
>
> Claude Code discards a WorktreeRemove hook's [JSON output fields](#json-output), such as `systemMessage` and `continue`.
>
> For a background-session delete, Claude Code verifies the stored worktree path before running the hook and refuses a path that is a symlink or passes through one below the repository root. The hook runs for a worktree that still contains files only when you confirm the delete in [agent view](/docs/en/agent-view#what-deleting-a-session-removes); for such a worktree, [`claude rm`](/docs/en/agent-view#manage-sessions-from-the-shell) keeps the session and worktree instead. Before v2.1.216, the hook ran on the stored path without these checks.
>
> Claude Code passes the path returned by WorktreeCreate as `worktree_path` in the hook input. This example reads that path and removes the directory:
>
> ```json
> {
>   "hooks": {
>     "WorktreeRemove": [
>       {
>         "hooks": [
>           {
>             "type": "command",
>             "command": "bash -c 'jq -r .worktree_path | xargs rm -rf'"
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> #### WorktreeRemove input
>
> In addition to the [common input fields](#common-input-fields), WorktreeRemove hooks receive the `worktree_path` field, which is the absolute path to the worktree being removed.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "hook_event_name": "WorktreeRemove",
>   "worktree_path": "/Users/.../my-project/.claude/worktrees/feature-auth"
> }
> ```
>
> WorktreeRemove hooks have no decision control. They can't block worktree removal but can perform cleanup tasks like removing version control state or archiving changes. Hook failures are logged in debug mode only.

### Source: Hooks reference — InstructionsLoaded

> ### InstructionsLoaded
>
> Fires when a `CLAUDE.md` or `.claude/rules/*.md` file is loaded into context. This event fires at session start for eagerly-loaded files and again later when files are lazily loaded, for example when Claude accesses a subdirectory that contains a nested `CLAUDE.md` or when conditional rules with `paths:` frontmatter match. The hook doesn't support blocking or decision control. It runs asynchronously for observability purposes.
>
> The matcher runs against `load_reason`. For example, use `"matcher": "session_start"` to fire only for files loaded at session start, or `"matcher": "path_glob_match|nested_traversal"` to fire only for lazy loads.
>
> #### InstructionsLoaded input
>
> In addition to the [common input fields](#common-input-fields), InstructionsLoaded hooks receive these fields:
>
> | Field               | Description                                                                                                                                                                                                   |
> | :------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `file_path`         | Absolute path to the instruction file that was loaded                                                                                                                                                         |
> | `memory_type`       | Scope of the file: `"User"`, `"Project"`, `"Local"`, or `"Managed"`                                                                                                                                           |
> | `load_reason`       | Why the file was loaded: `"session_start"`, `"nested_traversal"`, `"path_glob_match"`, `"include"`, or `"compact"`. The `"compact"` value fires when instruction files are re-loaded after a compaction event |
> | `globs`             | Path glob patterns from the file's `paths:` frontmatter, if any. Present only for `path_glob_match` loads                                                                                                     |
> | `trigger_file_path` | Path to the file whose access triggered this load, for lazy loads                                                                                                                                             |
> | `parent_file_path`  | Path to the parent instruction file that included this one, for `include` loads                                                                                                                               |
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../transcript.jsonl",
>   "cwd": "/Users/my-project",
>   "hook_event_name": "InstructionsLoaded",
>   "file_path": "/Users/my-project/CLAUDE.md",
>   "memory_type": "Project",
>   "load_reason": "session_start"
> }
> ```
>
> #### InstructionsLoaded decision control
>
> InstructionsLoaded hooks have no decision control. They can't block or modify instruction loading. Claude Code discards their [JSON output fields](#json-output), such as `systemMessage` and `continue`. Use this event for audit logging, compliance tracking, or observability.
