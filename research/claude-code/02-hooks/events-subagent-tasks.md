---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Hook events — subagents and tasks"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook events — subagents and tasks

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — Hook events

> ## Hook events
>
> Each event corresponds to a point in Claude Code's lifecycle where hooks can run. The sections below are ordered to match the lifecycle: from session setup through the agentic loop to session end. Each section describes when the event fires, what matchers it supports, the JSON input it receives, and how to control behavior through output.

### Source: Hooks reference — SubagentStart

> ### SubagentStart
>
> Runs when a Claude Code subagent is spawned via the Agent tool. Supports matchers to filter by agent type name. For built-in agents, this is the agent name like `general-purpose`, `Explore`, or `Plan`. For [custom subagents](/docs/en/sub-agents), this is the `name` field from the agent's frontmatter, not the filename.
>
> For subagents shipped by a [plugin](/docs/en/plugins), the agent type is the plugin-scoped identifier such as `my-plugin:reviewer`, not the bare frontmatter name. The colon places a plugin-scoped name on the regular-expression path, so anchor the matcher with `^` and `$` for an exact match: `^my-plugin:reviewer$`.
>
> #### SubagentStart input
>
> In addition to the [common input fields](#common-input-fields), SubagentStart hooks receive `agent_id` with the unique identifier for the subagent and `agent_type` with the agent name that the matcher filters on.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "hook_event_name": "SubagentStart",
>   "agent_id": "agent-abc123",
>   "agent_type": "Explore"
> }
> ```
>
> SubagentStart hooks can't block subagent creation, but they can inject context into the subagent. In addition to the [JSON output fields](#json-output) available to all hooks, you can return:
>
> | Field               | Description                                                                                                                                             |
> | :------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `additionalContext` | String added to the subagent's context at the start of its conversation, before its first prompt. See [Add context for Claude](#add-context-for-claude) |
>
> ```json
> {
>   "hookSpecificOutput": {
>     "hookEventName": "SubagentStart",
>     "additionalContext": "Follow security guidelines for this task"
>   }
> }
> ```

### Source: Hooks reference — SubagentStop

> ### SubagentStop
>
> Runs when a Claude Code subagent has finished responding. Matches on agent type, same values as SubagentStart.
>
> #### SubagentStop input
>
> In addition to the [common input fields](#common-input-fields), SubagentStop hooks receive `stop_hook_active`, `agent_id`, `agent_type`, `agent_transcript_path`, and `last_assistant_message`. The `agent_type` field is the value used for matcher filtering. The `transcript_path` is the main session's transcript, while `agent_transcript_path` is the subagent's own transcript stored in a nested `subagents/` folder. The `last_assistant_message` field contains the text content of the subagent's final response, so hooks can access it without parsing the transcript file.
>
> SubagentStop hooks also receive the `background_tasks` and `session_crons` arrays described under [Stop input](#stop-input). Both arrays are scoped to the parent session, not the subagent.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "~/.claude/projects/.../abc123.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "SubagentStop",
>   "stop_hook_active": false,
>   "agent_id": "def456",
>   "agent_type": "Explore",
>   "agent_transcript_path": "~/.claude/projects/.../abc123/subagents/agent-def456.jsonl",
>   "last_assistant_message": "Analysis complete. Found 3 potential issues...",
>   "background_tasks": [],
>   "session_crons": []
> }
> ```
>
> SubagentStop hooks use the same decision control format as [Stop hooks](#stop-decision-control), including `hookSpecificOutput.additionalContext` with `hookEventName` set to `"SubagentStop"`, for non-error feedback that keeps the subagent running. Returning `decision: "block"` with a `reason` keeps the subagent running and delivers `reason` to the subagent as its next instruction. A hook that blocks by exiting 2 delivers its stderr message the same way. To inject context into the parent session after a subagent returns, use a [`PostToolUse`](#posttooluse) hook on the `Agent` tool instead.

### Source: Hooks reference — TaskCreated

> ### TaskCreated
>
> Runs when a task is being created via the `TaskCreate` tool. Use this to enforce naming conventions, require task descriptions, or prevent certain tasks from being created. In a [session without the Task tools](/docs/en/tools-reference#task-tool-availability), this event doesn't fire.
>
> TaskCreated hooks don't support matchers and fire on every occurrence.
>
> #### TaskCreated input
>
> In addition to the [common input fields](#common-input-fields), TaskCreated hooks receive `task_id`, `task_subject`, and optionally `task_description`, `teammate_name`, and `team_name`.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "TaskCreated",
>   "task_id": "task-001",
>   "task_subject": "Implement user authentication",
>   "task_description": "Add login and signup endpoints",
>   "teammate_name": "implementer",
>   "team_name": "session-a1b2c3d4"
> }
> ```
>
> | Field              | Description                                                                |
> | :----------------- | :------------------------------------------------------------------------- |
> | `task_id`          | Identifier of the task being created                                       |
> | `task_subject`     | Title of the task                                                          |
> | `task_description` | Detailed description of the task. May be absent                            |
> | `teammate_name`    | Name of the teammate creating the task. May be absent                      |
> | `team_name`        | Deprecated. Session-derived team name; will be removed in a future release |
>
> #### TaskCreated decision control
>
> A TaskCreated hook can block the creation in two ways. Either way, Claude Code deletes the task and returns your message to Claude as the tool's error. Claude Code ignores `continue: false` from this event and Claude keeps working.
>
> * **Exit code 2**: Claude Code returns the stderr text as the message.
> * **JSON `{"decision": "block", "reason": "..."}`**: Claude Code returns `reason` as the message.
>
> This example blocks tasks whose subjects don't follow the required format:
>
> ```bash
> #!/bin/bash
> INPUT=$(cat)
> TASK_SUBJECT=$(echo "$INPUT" | jq -r '.task_subject')
>
> if [[ ! "$TASK_SUBJECT" =~ ^\[TICKET-[0-9]+\] ]]; then
>   echo "Task subject must start with a ticket number, e.g. '[TICKET-123] Add feature'" >&2
>   exit 2
> fi
>
> exit 0
> ```

### Source: Hooks reference — TaskCompleted

> ### TaskCompleted
>
> Runs when a task is being marked as completed. This fires in two situations: when any agent explicitly marks a task as completed through the TaskUpdate tool, or when an [agent team](/docs/en/agent-teams) teammate finishes its turn with in-progress tasks. Use this to enforce completion criteria like passing tests or lint checks before a task can close.
>
> TaskCompleted hooks don't support matchers and fire on every occurrence.
>
> #### TaskCompleted input
>
> In addition to the [common input fields](#common-input-fields), TaskCompleted hooks receive `task_id`, `task_subject`, and optionally `task_description`, `teammate_name`, and `team_name`.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "TaskCompleted",
>   "task_id": "task-001",
>   "task_subject": "Implement user authentication",
>   "task_description": "Add login and signup endpoints",
>   "teammate_name": "implementer",
>   "team_name": "session-a1b2c3d4"
> }
> ```
>
> | Field              | Description                                                                |
> | :----------------- | :------------------------------------------------------------------------- |
> | `task_id`          | Identifier of the task being completed                                     |
> | `task_subject`     | Title of the task                                                          |
> | `task_description` | Detailed description of the task. May be absent                            |
> | `teammate_name`    | Name of the teammate completing the task. May be absent                    |
> | `team_name`        | Deprecated. Session-derived team name; will be removed in a future release |
>
> #### TaskCompleted decision control
>
> TaskCompleted hooks support two ways to control task completion:
>
> * **Exit code 2**: the task is not marked as completed and the stderr message is fed back to the model as feedback.
> * **JSON `{"continue": false, "stopReason": "..."}`**: stops the teammate entirely, matching `Stop` hook behavior. The `stopReason` is shown to the user.
>
> This example runs tests and blocks task completion if they fail:
>
> ```bash
> #!/bin/bash
> INPUT=$(cat)
> TASK_SUBJECT=$(echo "$INPUT" | jq -r '.task_subject')
>
> # Run the test suite
> if ! npm test 2>&1; then
>   echo "Tests not passing. Fix failing tests before completing: $TASK_SUBJECT" >&2
>   exit 2
> fi
>
> exit 0
> ```

### Source: Hooks reference — TeammateIdle

> ### TeammateIdle
>
> Runs when an [agent team](/docs/en/agent-teams) teammate is about to go idle after finishing its turn. Use this to enforce quality gates before a teammate stops working, such as requiring passing lint checks or verifying that output files exist.
>
> TeammateIdle hooks don't support matchers and fire on every occurrence.
>
> #### TeammateIdle input
>
> In addition to the [common input fields](#common-input-fields), TeammateIdle hooks receive `teammate_name` and `team_name`.
>
> ```json
> {
>   "session_id": "abc123",
>   "transcript_path": "/Users/.../.claude/projects/.../00893aaf-19fa-41d2-8238-13269b9b3ca0.jsonl",
>   "cwd": "/Users/...",
>   "permission_mode": "default",
>   "hook_event_name": "TeammateIdle",
>   "teammate_name": "researcher",
>   "team_name": "session-a1b2c3d4"
> }
> ```
>
> | Field           | Description                                                                |
> | :-------------- | :------------------------------------------------------------------------- |
> | `teammate_name` | Name of the teammate that is about to go idle                              |
> | `team_name`     | Deprecated. Session-derived team name; will be removed in a future release |
>
> #### TeammateIdle decision control
>
> TeammateIdle hooks support two ways to control teammate behavior:
>
> * **Exit code 2**: the teammate receives the stderr message as feedback and continues working instead of going idle.
> * **JSON `{"continue": false, "stopReason": "..."}`**: stops the teammate entirely, matching `Stop` hook behavior. The `stopReason` is shown to the user.
>
> This example checks that a build artifact exists before allowing a teammate to go idle:
>
> ```bash
> #!/bin/bash
>
> if [ ! -f "./dist/output.js" ]; then
>   echo "Build artifact missing. Run the build before stopping." >&2
>   exit 2
> fi
>
> exit 0
> ```
