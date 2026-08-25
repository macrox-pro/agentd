---
primary_sources:
  - id: T1-SDK-HOOKS
    title: "SDK hooks"
    url: "https://code.claude.com/docs/en/agent-sdk/hooks.md"
    section: ""
  - id: T1-SDK-PERMS
    title: "SDK permissions"
    url: "https://code.claude.com/docs/en/agent-sdk/permissions.md"
    section: ""
  - id: T1-SDK-FEATURES
    title: "SDK features"
    url: "https://code.claude.com/docs/en/agent-sdk/claude-code-features.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Agent SDK hooks and permissions

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Agent SDK — Hooks
>
> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Intercept and control agent behavior with hooks
>
> > Intercept and customize agent behavior at key execution points with hooks
>
> Hooks are callback functions that run your code in response to agent events, like a tool being called, a session starting, or execution stopping. With hooks, you can:
>
> * **Block dangerous operations** before they execute, like destructive shell commands or unauthorized file access
> * **Log and audit** every tool call for compliance, debugging, or analytics
> * **Transform inputs and outputs** to sanitize data, inject credentials, or redirect file paths
> * **Require human approval** for sensitive actions like database writes or API calls
> * **Track session lifecycle** to manage state, clean up resources, or send notifications
>
> ## How hooks work
>
>   **An event fires**: Something happens during agent execution and the SDK fires an event: a tool is about to be called (`PreToolUse`), a tool returned a result (`PostToolUse`), a subagent started or stopped, the agent is idle, or execution finished. See the [full list of events](#available-hooks).
>
>   **The SDK collects registered hooks**: The SDK checks for hooks registered for that event type. This includes callback hooks you pass in `options.hooks` and shell command hooks from settings files when the corresponding [`settingSources`](/docs/en/agent-sdk/typescript#settingsource) or [`setting_sources`](/docs/en/agent-sdk/python#settingsource) entry is enabled, which it is for default `query()` options.
>
>   **Matchers filter which hooks run**: If a hook has a [`matcher`](#matchers) pattern (like `"Write|Edit"`), the SDK tests it against the event's target (for example, the tool name). Hooks without a matcher run for every event of that type.
>
>   **Callback functions execute**: Each matching hook's [callback function](#callback-functions) receives input about what's happening: the tool name, its arguments, the session ID, and other event-specific details.
>
>   **Your callback returns a decision**: After performing any operations (logging, API calls, validation), your callback returns an [output object](#outputs) that tells the agent what to do: allow the operation, block it, modify the input, or inject context into the conversation.
>
> The following example puts these steps together. It registers a `PreToolUse` hook (step 1) with a `"Write|Edit"` matcher (step 3) so the callback only fires for file-writing tools. When triggered, the callback receives the tool's input (step 4), checks if the file path targets a `.env` file, and returns `permissionDecision: "deny"` to block the operation (step 5):
>
>   ```python Python
>   import asyncio
>   from claude_agent_sdk import (
>       AssistantMessage,
>       ClaudeSDKClient,
>       ClaudeAgentOptions,
>       HookMatcher,
>       ResultMessage,
>   )
>
>   # Define a hook callback that receives tool call details
>   async def protect_env_files(input_data, tool_use_id, context):
>       # Extract the file path from the tool's input arguments
>       file_path = input_data["tool_input"].get("file_path", "")
>       file_name = file_path.split("/")[-1]
>
>       # Block the operation if targeting a .env file
>       if file_name == ".env":
>           return {
>               "hookSpecificOutput": {
>                   "hookEventName": input_data["hook_event_name"],
>                   "permissionDecision": "deny",
>                   "permissionDecisionReason": "Cannot modify .env files",
>               }
>           }
>
>       # Return empty object to allow the operation
>       return {}
>
>   async def main():
>       options = ClaudeAgentOptions(
>           hooks={
>               # Register the hook for PreToolUse events
>               # The matcher filters to only Write and Edit tool calls
>               "PreToolUse": [HookMatcher(matcher="Write|Edit", hooks=[protect_env_files])]
>           }
>       )
>
>       async with ClaudeSDKClient(options=options) as client:
>           await client.query("Create a .env file with the standard local development database configuration")
>           async for message in client.receive_response():
>               # Filter for assistant and result messages
>               if isinstance(message, (AssistantMessage, ResultMessage)):
>                   print(message)
>
>   asyncio.run(main())
>   ```
>
>   ```typescript TypeScript
>   import { query, HookCallback, PreToolUseHookInput } from "@anthropic-ai/claude-agent-sdk";
>
>   // Define a hook callback with the HookCallback type
>   const protectEnvFiles: HookCallback = async (input, toolUseID, { signal }) => {
>     // Cast input to the specific hook type for type safety
>     const preInput = input as PreToolUseHookInput;
>
>     // Cast tool_input to access its properties (typed as unknown in the SDK)
>     const toolInput = preInput.tool_input as Record<string, unknown>;
>     const filePath = toolInput?.file_path as string;
>     const fileName = filePath?.split("/").pop();
>
>     // Block the operation if targeting a .env file
>     if (fileName === ".env") {
>       return {
>         hookSpecificOutput: {
>           hookEventName: preInput.hook_event_name,
>           permissionDecision: "deny",
>           permissionDecisionReason: "Cannot modify .env files"
>         }
>       };
>     }
>
>     // Return empty object to allow the operation
>     return {};
>   };
>
>   for await (const message of query({
>     prompt: "Create a .env file with the standard local development database configuration",
>     options: {
>       hooks: {
>         // Register the hook for PreToolUse events
>         // The matcher filters to only Write and Edit tool calls
>         PreToolUse: [{ matcher: "Write|Edit", hooks: [protectEnvFiles] }]
>       }
>     }
>   })) {
>     // Filter for assistant and result messages
>     if (message.type === "assistant" || message.type === "result") {
>       console.log(message);
>     }
>   }
>   ```
>
> When you run either script, Claude attempts to create the `.env` file, the hook denies the tool call, and Claude's final response explains that it can't create `.env` files.
>
> ## Available hooks
>
> The SDK provides hooks for different stages of agent execution. Some hooks are available in both SDKs, while others are TypeScript-only.
>
> | Hook Event                                             | Python SDK | TypeScript SDK | What triggers it                                                                                                                        | Example use case                                                                                                                                          |
> | ------------------------------------------------------ | ---------- | -------------- | --------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `PreToolUse`                                           | Yes        | Yes            | Tool call request (can block or modify)                                                                                                 | Block dangerous shell commands                                                                                                                            |
> | `PostToolUse`                                          | Yes        | Yes            | Tool execution result                                                                                                                   | Log all file changes to audit trail                                                                                                                       |
> | `PostToolUseFailure`                                   | Yes        | Yes            | Tool execution failure                                                                                                                  | Handle or log tool errors                                                                                                                                 |
> | `PostToolBatch`                                        | No         | Yes            | A full batch of tool calls resolves, once per batch before the next model call                                                          | Inject conventions once for the whole batch                                                                                                               |
> | `UserPromptSubmit`                                     | Yes        | Yes            | User prompt submission                                                                                                                  | Inject additional context into prompts                                                                                                                    |
> | [`UserPromptExpansion`](/docs/en/hooks#userpromptexpansion) | No         | Yes            | A user-typed command, or an MCP prompt, expands into a prompt before it reaches Claude. Doesn't fire when Claude invokes a skill itself | Block a command from direct invocation or add context when a skill is typed                                                                               |
> | `MessageDisplay`                                       | No         | Yes            | An assistant message with text completes, once per message with the full message text                                                   | Redact or reformat the displayed text without changing the transcript                                                                                     |
> | `Stop`                                                 | Yes        | Yes            | Agent execution stop                                                                                                                    | Save session state before exit                                                                                                                            |
> | `StopFailure`                                          | No         | Yes            | The turn ends with an API error instead of a normal stop                                                                                | Log failures or send alerts                                                                                                                               |
> | `SubagentStart`                                        | Yes        | Yes            | Subagent initialization                                                                                                                 | Track parallel task spawning                                                                                                                              |
> | `SubagentStop`                                         | Yes        | Yes            | Subagent completion                                                                                                                     | Aggregate results from parallel tasks                                                                                                                     |
> | `PreCompact`                                           | Yes        | Yes            | Conversation compaction request                                                                                                         | Archive full transcript before summarizing                                                                                                                |
> | `PostCompact`                                          | No         | Yes            | Conversation compaction completes                                                                                                       | Log the generated summary                                                                                                                                 |
> | `PermissionRequest`                                    | Yes        | Yes            | A tool call needs a permission decision                                                                                                 | Custom permission handling                                                                                                                                |
> | `PermissionDenied`                                     | No         | Yes            | Auto mode denies a tool call, including denials without a classifier verdict                                                            | Log denials, or tell the model it may retry; Claude Code ignores `retry: true` for no-verdict denials. See [PermissionDenied](/docs/en/hooks#permissiondenied) |
> | `SessionStart`                                         | No         | Yes            | Session initialization                                                                                                                  | Initialize logging and telemetry                                                                                                                          |
> | `SessionEnd`                                           | No         | Yes            | Session termination                                                                                                                     | Clean up temporary resources                                                                                                                              |
> | `Notification`                                         | Yes        | Yes            | Agent status messages                                                                                                                   | Send agent status updates to Slack or PagerDuty                                                                                                           |
> | `Setup`                                                | No         | Yes            | Session setup/maintenance                                                                                                               | Run initialization tasks                                                                                                                                  |
> | `TeammateIdle`                                         | No         | Yes            | Teammate becomes idle                                                                                                                   | Reassign work or notify                                                                                                                                   |
> | `TaskCreated`                                          | No         | Yes            | A task is created via the `TaskCreate` tool                                                                                             | Enforce task naming conventions                                                                                                                           |
> | [`TaskCompleted`](/docs/en/hooks#taskcompleted)             | No         | Yes            | A task is marked completed                                                                                                              | Require passing tests before a task closes                                                                                                                |
> | `Elicitation`                                          | No         | Yes            | An MCP server requests user input mid-task                                                                                              | Respond to MCP input requests programmatically                                                                                                            |
> | `ElicitationResult`                                    | No         | Yes            | A user responds to an MCP elicitation                                                                                                   | Modify or block the response before it returns to the server                                                                                              |
> | `ConfigChange`                                         | No         | Yes            | Configuration file changes                                                                                                              | Reload settings dynamically                                                                                                                               |
> | `InstructionsLoaded`                                   | No         | Yes            | A `CLAUDE.md` or rules file is loaded into context                                                                                      | Audit which instruction files load                                                                                                                        |
> | `WorktreeCreate`                                       | No         | Yes            | Git worktree created                                                                                                                    | Track isolated workspaces                                                                                                                                 |
> | `WorktreeRemove`                                       | No         | Yes            | Git worktree removed                                                                                                                    | Clean up workspace resources                                                                                                                              |
> | `CwdChanged`                                           | No         | Yes            | The working directory changes during a session                                                                                          | Reload environment variables per directory                                                                                                                |
> | `FileChanged`                                          | No         | Yes            | A watched file is modified, created, or deleted                                                                                         | Reload configuration when project files change                                                                                                            |
> | `DirectoryAdded`                                       | No         | Yes            | A working directory is added during a session                                                                                           | Install dependencies for a repository added mid-session                                                                                                   |
>
> ## Configure hooks
>
> To configure a hook, pass it in the `hooks` field of your agent options (`ClaudeAgentOptions` in Python, the `options` object in TypeScript). This snippet assumes you have already defined a hook callback, like `protect_env_files` in Python or `protectEnvFiles` in TypeScript from the example above:
>
>   ```python Python
>   options = ClaudeAgentOptions(
>       hooks={"PreToolUse": [HookMatcher(matcher="Bash", hooks=[my_callback])]}
>   )
>
>   async with ClaudeSDKClient(options=options) as client:
>       await client.query("Your prompt")
>       async for message in client.receive_response():
>           print(message)
>   ```
>
>   ```typescript TypeScript
>   for await (const message of query({
>     prompt: "Your prompt",
>     options: {
>       hooks: {
>         PreToolUse: [{ matcher: "Bash", hooks: [myCallback] }]
>       }
>     }
>   })) {
>     console.log(message);
>   }
>   ```
>
> The `hooks` option is a dictionary in Python or an object in TypeScript, where:
>
> * **Keys**: [hook event names](#available-hooks) such as `'PreToolUse'`, `'PostToolUse'`, and `'Stop'`
> * **Values**: arrays of [matchers](#matchers), each containing an optional filter pattern and your [callback functions](#callback-functions)
>
> ### Matchers
>
> Use matchers to filter when your callbacks fire. The `matcher` field matches against a different value depending on the hook event type. For example, tool-based hooks match against the tool name, while `Notification` hooks match against the notification type.
>
> SDK matchers follow the same rules as [matchers in settings files](/docs/en/hooks#matcher-patterns). That section documents the exact-string and regular-expression evaluation paths, their version requirements, and the matcher values for each event type.
>
> | Option    | Type             | Default     | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
> | --------- | ---------------- | ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `matcher` | `string`         | `undefined` | Pattern matched against the event's filter field, following the [rules for matchers in settings files](/docs/en/hooks#matcher-patterns). For tool hooks, this is the tool name. Built-in tools include `Bash`, `Read`, `Write`, `Edit`, `Glob`, `Grep`, `WebFetch`, `Agent`, and others (see [Tool Input Types](/docs/en/agent-sdk/typescript#tool-input-types) for the full list). MCP tools use the pattern `mcp__<server>__<action>`, where `<server>` is the key you use in the `mcpServers` configuration. |
> | `hooks`   | `HookCallback[]` | -           | Required. Array of callback functions to execute when the pattern matches                                                                                                                                                                                                                                                                                                                                                                                                                             |
> | `timeout` | `number`         | `undefined` | Timeout in seconds. When omitted, Claude Code applies the [event's default timeout](#hook-timeout). Your SDK callbacks follow the `command` hook defaults                                                                                                                                                                                                                                                                                                                                             |
>
> Use the `matcher` pattern to target specific tools whenever possible. A matcher with `'Bash'` only runs for Bash commands, while omitting the pattern runs your callbacks for every occurrence of the event. Omit it on purpose to log every tool call your session makes.
>
> ### Callback functions
>
> #### Inputs
>
> Every hook callback receives three arguments:
>
> * **Input data:** a typed object containing event details. Each hook type has its own input shape. For example, `PreToolUseHookInput` includes `tool_name` and `tool_input`, while `NotificationHookInput` includes `message`. See the full type definitions in the [TypeScript](/docs/en/agent-sdk/typescript#hookinput) and [Python](/docs/en/agent-sdk/python#hookinput) SDK references.
>   * All hook inputs share `session_id`, `cwd`, and `hook_event_name`.
>   * `agent_id` and `agent_type` are populated when the hook fires inside a subagent. In TypeScript, these are on the base hook input and available to all hook types. In Python, they are optional fields on `PreToolUse`, `PostToolUse`, `PostToolUseFailure`, and `PermissionRequest`, and required fields on `SubagentStart` and `SubagentStop`.
> * **Tool use ID** (`str | None` / `string | undefined`): correlates `PreToolUse` and `PostToolUse` events for the same tool call.
> * **Context:** in TypeScript, contains a `signal` property (`AbortSignal`) for cancellation. In Python, this argument is reserved for future use.
>
> #### Outputs
>
> Your callback returns an object with two categories of fields:
>
> * **Top-level fields** are accepted on every event: `systemMessage` shows a message to the user, and `continue` (`continue_` in Python) determines whether the agent keeps running after this hook. Some events discard them or deliver them elsewhere. Each [event's section](/docs/en/hooks#hook-events) on the hooks page says where they land.
> * **`hookSpecificOutput`** controls the current operation. The fields you set inside depend on the hook event type:
>   * For `PreToolUse` hooks, this is where you set `permissionDecision` (`"allow"`, `"deny"`, `"ask"`, or `"defer"`), `permissionDecisionReason`, and `updatedInput`. If you return `"defer"`, the query ends so you can [resume it later](/docs/en/hooks#defer-a-tool-call-for-later).
>   * For `PostToolUse` hooks, you can set `additionalContext` to append information to the tool result. To replace the tool's output before Claude sees it, set `updatedToolOutput`, which works for any tool in both SDKs. The older `updatedMCPToolOutput` field replaces MCP tool output only and is deprecated.
>   * In the TypeScript SDK, a `PostToolUse` callback can also return `classifierContext`, a short note about the tool call's result for the [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) permission classifier. Because your callback runs in your application's own process, the classifier may weigh a user statement you relay in the note as user intent. The field requires TypeScript Agent SDK v0.3.236 or later. [Annotate a result for the auto mode classifier](/docs/en/hooks#annotate-a-result-for-the-auto-mode-classifier) covers the length cap, the synchronous-only rule, and what not to put in the note.
>
> Return `{}` to allow the operation without changes. SDK callback hooks use the same JSON output format as [Claude Code shell command hooks](/docs/en/hooks#json-output), which documents every field and event-specific option. For the SDK type definitions, see the [TypeScript](/docs/en/agent-sdk/typescript#synchookjsonoutput) and [Python](/docs/en/agent-sdk/python#synchookjsonoutput) SDK references.
>
>   When multiple hooks or permission rules apply, `deny` takes priority over `defer`, which takes priority over `ask`, which takes priority over `allow`. If any hook returns `deny`, the operation is blocked regardless of other hooks.
>
> #### Asynchronous output
>
> By default, the agent waits for your hook to return before proceeding. If your hook performs a side effect, such as logging or sending a webhook, and doesn't need to influence the agent's behavior, you can return an async output instead. This tells the agent to continue immediately without waiting for the hook to finish. In this snippet, `send_to_logging_service` in Python and `sendToLoggingService` in TypeScript stand in for any logging function you define:
>
>   ```python Python
>   async def async_hook(input_data, tool_use_id, context):
>       # Start a background task, then return immediately
>       asyncio.create_task(send_to_logging_service(input_data))
>       return {"async_": True, "asyncTimeout": 30000}
>   ```
>
>   ```typescript TypeScript
>   const asyncHook: HookCallback = async (input, toolUseID, { signal }) => {
>     // Start a background task, then return immediately
>     sendToLoggingService(input).catch(console.error);
>     return { async: true, asyncTimeout: 30000 };
>   };
>   ```
>
> | Field          | Type     | Description                                                                                                    |
> | -------------- | -------- | -------------------------------------------------------------------------------------------------------------- |
> | `async`        | `true`   | Signals async mode. The agent proceeds without waiting. In Python, use `async_` to avoid the reserved keyword. |
> | `asyncTimeout` | `number` | Optional timeout in milliseconds for the background operation                                                  |
>
>   Async outputs can't block, modify, or inject context into the operation since the agent has already moved on. Use them only for side effects like logging, metrics, or notifications.
>
> ## Examples
>
> Several examples in this section show only the callback function. To run one, register the callback under the matching event in the `hooks` field of your options, as shown in [Configure hooks](#configure-hooks).
>
> ### Modify tool input
>
> This example intercepts Write tool calls and rewrites the `file_path` argument to prepend `/sandbox`, redirecting all file writes to a sandboxed directory. The callback returns `updatedInput` with the modified path and `permissionDecision: 'allow'` to auto-approve the rewritten operation:
>
>   ```python Python
>   async def redirect_to_sandbox(input_data, tool_use_id, context):
>       if input_data["hook_event_name"] != "PreToolUse":
>           return {}
>
>       if input_data["tool_name"] == "Write":
>           original_path = input_data["tool_input"].get("file_path", "")
>           return {
>               "hookSpecificOutput": {
>                   "hookEventName": input_data["hook_event_name"],
>                   "permissionDecision": "allow",
>                   "updatedInput": {
>                       **input_data["tool_input"],
>                       "file_path": f"/sandbox{original_path}",
>                   },
>               }
>           }
>       return {}
>   ```
>
>   ```typescript TypeScript
>   const redirectToSandbox: HookCallback = async (input, toolUseID, { signal }) => {
>     if (input.hook_event_name !== "PreToolUse") return {};
>
>     const preInput = input as PreToolUseHookInput;
>     const toolInput = preInput.tool_input as Record<string, unknown>;
>     if (preInput.tool_name === "Write") {
>       const originalPath = toolInput.file_path as string;
>       return {
>         hookSpecificOutput: {
>           hookEventName: preInput.hook_event_name,
>           permissionDecision: "allow",
>           updatedInput: {
>             ...toolInput,
>             file_path: `/sandbox${originalPath}`
>           }
>         }
>       };
>     }
>     return {};
>   };
>   ```
>
>   Pair `updatedInput` with `permissionDecision: 'allow'` to auto-approve the modified input, or `permissionDecision: 'ask'` to show it to the user. If you omit `permissionDecision`, the modified input still applies and flows through the normal permission evaluation. With `'defer'`, `updatedInput` is ignored. Always return a new object rather than mutating the original `tool_input`.
>
> To confirm the redirect, set the prefix to a path you can write to, such as `./sandbox` or `/tmp/sandbox` (macOS doesn't allow creating a root-level `/sandbox` directory), then ask the agent to write a file: the Write tool's result in the message stream names the path with your sandbox prefix rather than the one Claude requested.
>
> ### Add context and block a tool
>
> This example blocks writes to the `/etc` directory and explains why to both the model and the user:
>
> * `permissionDecision: 'deny'` stops the tool call.
> * `permissionDecisionReason` tells the model why, so it avoids retrying.
> * `systemMessage` shows the user what happened.
>
>   ```python Python
>   async def block_etc_writes(input_data, tool_use_id, context):
>       file_path = input_data["tool_input"].get("file_path", "")
>
>       if file_path.startswith("/etc"):
>           return {
>               # Top-level field: message shown to the user
>               "systemMessage": "Remember: system directories like /etc are protected.",
>               # hookSpecificOutput: block the operation
>               "hookSpecificOutput": {
>                   "hookEventName": input_data["hook_event_name"],
>                   "permissionDecision": "deny",
>                   "permissionDecisionReason": "Writing to /etc is not allowed",
>               },
>           }
>       return {}
>   ```
>
>   ```typescript TypeScript
>   const blockEtcWrites: HookCallback = async (input, toolUseID, { signal }) => {
>     const preInput = input as PreToolUseHookInput;
>     const toolInput = preInput.tool_input as Record<string, unknown>;
>     const filePath = toolInput?.file_path as string;
>
>     if (filePath?.startsWith("/etc")) {
>       return {
>         // Top-level field: message shown to the user
>         systemMessage: "Remember: system directories like /etc are protected.",
>         // hookSpecificOutput: block the operation
>         hookSpecificOutput: {
>           hookEventName: preInput.hook_event_name,
>           permissionDecision: "deny",
>           permissionDecisionReason: "Writing to /etc is not allowed"
>         }
>       };
>     }
>     return {};
>   };
>   ```
>
> ### Auto-approve specific tools
>
> By default, the agent may prompt for permission before using certain tools. This example auto-approves read-only filesystem tools (Read, Glob, Grep) by returning `permissionDecision: 'allow'`, letting them run without user confirmation while leaving all other tools subject to normal permission checks:
>
>   ```python Python
>   async def auto_approve_read_only(input_data, tool_use_id, context):
>       if input_data["hook_event_name"] != "PreToolUse":
>           return {}
>
>       read_only_tools = ["Read", "Glob", "Grep"]
>       if input_data["tool_name"] in read_only_tools:
>           return {
>               "hookSpecificOutput": {
>                   "hookEventName": input_data["hook_event_name"],
>                   "permissionDecision": "allow",
>                   "permissionDecisionReason": "Read-only tool auto-approved",
>               }
>           }
>       return {}
>   ```
>
>   ```typescript TypeScript
>   const autoApproveReadOnly: HookCallback = async (input, toolUseID, { signal }) => {
>     if (input.hook_event_name !== "PreToolUse") return {};
>
>     const preInput = input as PreToolUseHookInput;
>     const readOnlyTools = ["Read", "Glob", "Grep"];
>     if (readOnlyTools.includes(preInput.tool_name)) {
>       return {
>         hookSpecificOutput: {
>           hookEventName: preInput.hook_event_name,
>           permissionDecision: "allow",
>           permissionDecisionReason: "Read-only tool auto-approved"
>         }
>       };
>     }
>     return {};
>   };
>   ```
>
> ### Register multiple hooks
>
> When an event fires, all matching hooks run in parallel. For permission decisions, the most restrictive result applies: a single `deny` blocks the tool call regardless of what the other hooks return. Because completion order is non-deterministic, write each hook to act independently rather than relying on another hook having run first.
>
> The example below registers three independent checks for every tool call:
>
>   ```python Python
>   options = ClaudeAgentOptions(
>       hooks={
>           "PreToolUse": [
>               HookMatcher(hooks=[authorization_check]),
>               HookMatcher(hooks=[input_validator]),
>               HookMatcher(hooks=[audit_logger]),
>           ]
>       }
>   )
>   ```
>
>   ```typescript TypeScript
>   const options = {
>     hooks: {
>       PreToolUse: [
>         { hooks: [authorizationCheck] },
>         { hooks: [inputValidator] },
>         { hooks: [auditLogger] }
>       ]
>     }
>   };
>   ```
>
> ### Filter with multi-tool matchers
>
> Use multi-tool matchers to share one callback across related tools. This example registers three matchers with different scopes:
>
> * A pipe-separated exact list (`Write|Edit|NotebookEdit`) triggers `file_security_hook` only for file modification tools.
> * A regex (`^mcp__`) triggers `mcp_audit_hook` for any MCP tool whose name starts with `mcp__`.
> * An omitted matcher triggers `global_logger` for every tool call regardless of name.
>
>   ```python Python
>   options = ClaudeAgentOptions(
>       hooks={
>           "PreToolUse": [
>               # Match file modification tools
>               HookMatcher(matcher="Write|Edit|NotebookEdit", hooks=[file_security_hook]),
>               # Match all MCP tools
>               HookMatcher(matcher="^mcp__", hooks=[mcp_audit_hook]),
>               # Match everything (no matcher)
>               HookMatcher(hooks=[global_logger]),
>           ]
>       }
>   )
>   ```
>
>   ```typescript TypeScript
>   const options = {
>     hooks: {
>       PreToolUse: [
>         // Match file modification tools
>         { matcher: "Write|Edit|NotebookEdit", hooks: [fileSecurityHook] },
>
>         // Match all MCP tools
>         { matcher: "^mcp__", hooks: [mcpAuditHook] },
>
>         // Match everything (no matcher)
>         { hooks: [globalLogger] }
>       ]
>     }
>   };
>   ```
>
> ### Track subagent activity
>
> Use `SubagentStop` hooks to monitor when subagents finish their work. See the full input type in the [TypeScript](/docs/en/agent-sdk/typescript#hookinput) and [Python](/docs/en/agent-sdk/python#hookinput) SDK references. This example logs a summary each time a subagent completes:
>
>   ```python Python
>   async def subagent_tracker(input_data, tool_use_id, context):
>       # Log subagent details when it finishes
>       print(f"[SUBAGENT] Completed: {input_data['agent_id']}")
>       print(f"  Transcript: {input_data['agent_transcript_path']}")
>       print(f"  Tool use ID: {tool_use_id}")
>       print(f"  Stop hook active: {input_data.get('stop_hook_active')}")
>       return {}
>
>   options = ClaudeAgentOptions(
>       hooks={"SubagentStop": [HookMatcher(hooks=[subagent_tracker])]}
>   )
>   ```
>
>   ```typescript TypeScript
>   import { HookCallback, SubagentStopHookInput } from "@anthropic-ai/claude-agent-sdk";
>
>   const subagentTracker: HookCallback = async (input, toolUseID, { signal }) => {
>     // Cast to SubagentStopHookInput to access subagent-specific fields
>     const subInput = input as SubagentStopHookInput;
>
>     // Log subagent details when it finishes
>     console.log(`[SUBAGENT] Completed: ${subInput.agent_id}`);
>     console.log(`  Transcript: ${subInput.agent_transcript_path}`);
>     console.log(`  Tool use ID: ${toolUseID}`);
>     console.log(`  Stop hook active: ${subInput.stop_hook_active}`);
>     return {};
>   };
>
>   const options = {
>     hooks: {
>       SubagentStop: [{ hooks: [subagentTracker] }]
>     }
>   };
>   ```
>
> ### Make HTTP requests from hooks
>
> Hooks can perform asynchronous operations like HTTP requests. Catch errors inside your hook instead of letting them propagate, since an unhandled exception can interrupt the agent.
>
> This example sends a webhook after each tool completes, logging which tool ran and when. The hook catches errors so a failed webhook doesn't interrupt the agent:
>
>   ```python Python
>   import asyncio
>   import json
>   import urllib.request
>   from datetime import datetime
>
>   def _send_webhook(tool_name):
>       """Synchronous helper that POSTs tool usage data to an external webhook."""
>       data = json.dumps(
>           {
>               "tool": tool_name,
>               "timestamp": datetime.now().isoformat(),
>           }
>       ).encode()
>       req = urllib.request.Request(
>           "https://api.example.com/webhook",
>           data=data,
>           headers={"Content-Type": "application/json"},
>           method="POST",
>       )
>       urllib.request.urlopen(req)
>
>   async def webhook_notifier(input_data, tool_use_id, context):
>       # Only fire after a tool completes (PostToolUse), not before
>       if input_data["hook_event_name"] != "PostToolUse":
>           return {}
>
>       try:### Source: Agent SDK — Permissions
>
> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Configure permissions
>
> > Control how your agent uses tools with permission modes, hooks, and declarative allow/deny rules.
>
> The Claude Agent SDK provides permission controls to manage how Claude uses tools. Use permission modes and rules to define what's allowed automatically, and the [`canUseTool` callback](/docs/en/agent-sdk/user-input) to handle everything else at runtime.
>
> ## How permissions are evaluated
>
> When Claude requests a tool, the SDK checks permissions in this order:
>
>   **Hooks**: Run [hooks](/docs/en/agent-sdk/hooks) first. A hook can deny the call outright or pass it on. A hook that returns `allow` does not skip the deny and ask rules below; those are evaluated regardless of the hook result. A `PreToolUse` hook allow also can't approve an `rm` or `rmdir` removal targeting a [critical path](/docs/en/permission-modes#critical-paths).
>
>   **Deny rules**: Check `deny` rules (from `disallowed_tools` and [settings.json](/docs/en/settings-reference#permission-settings)). If a deny rule matches, the tool is blocked, even in `bypassPermissions` mode. Bare-name deny rules like `Bash` remove the tool from Claude's context before this evaluation begins, so only scoped rules like `Bash(rm *)` are checked at this step.
>
>   **Ask rules**: Check `ask` rules from [settings.json](/docs/en/settings-reference#permission-settings). If an ask rule matches, the call falls through to your [`canUseTool` callback](/docs/en/agent-sdk/user-input) for confirmation, even in `bypassPermissions` mode.
>
>     Tools that require user interaction behave the same way: `AskUserQuestion` and MCP tools whose server sets [`_meta["anthropic/requiresUserInteraction"]`](/docs/en/mcp#require-approval-for-a-specific-tool) always fall through to the callback, even when an allow rule matches. In `dontAsk` mode both cases are denied instead, because that mode never prompts. The MCP annotation requires Claude Code v2.1.199 or later.
>
>     [claude.ai connector](/docs/en/mcp#organization-controls-on-connector-tools) tools your organization has set to `ask` also leave the flow at this step. Every call falls through to the callback, even in `bypassPermissions` mode and even when an allow rule matches. The callback receives the reason `Your organization requires approval for this tool`. In `dontAsk` mode the call is denied instead, because that mode never prompts.
>
>   **Permission mode**: Apply the active [permission mode](#permission-modes). `bypassPermissions` approves everything that reaches this step except `rm` and `rmdir` removals targeting a [critical path](/docs/en/permission-modes#critical-paths), which fall through instead. `acceptEdits` approves the file operations listed under [Accept edits mode](#accept-edits-mode-acceptedits). `plan` routes file-edit and shell-write tools to your `canUseTool` callback regardless of allow rules, so write operations cannot be auto-approved while planning. Other modes fall through.
>
>   **Allow rules**: Check `allow` rules (from `allowed_tools` and settings.json). If a rule matches, the tool is approved. `rm` and `rmdir` removals targeting a [critical path](/docs/en/permission-modes#critical-paths) are never approved by an allow rule: they reach your callback in the modes that prompt, go to the [classifier](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) in `auto` mode on Claude Code v2.1.218 or later, and are denied in `dontAsk` mode.
>
>   **canUseTool callback**: If not resolved by any of the above, call your [`canUseTool` callback](/docs/en/agent-sdk/user-input) for a decision. In `dontAsk` mode, this step is skipped and the tool is denied.
>
> If you pass a `canUseTool` callback in a configuration where the TypeScript SDK expects the evaluation order to auto-approve calls before the callback is consulted, the SDK emits a Node.js process warning once when the query is constructed. The warning's code is `CLAUDE_SDK_CAN_USE_TOOL_SHADOWED`. Two configurations trigger it:
>
> * `permissionMode: 'bypassPermissions'`, which auto-approves every call that reaches the permission mode step apart from the [actions no mode auto-approves](/docs/en/permission-modes#actions-no-mode-auto-approves)
> * Each bare `allowedTools` entry such as `"Read"`, which auto-approves that whole tool before the callback is consulted, apart from the [actions no mode auto-approves](/docs/en/permission-modes#actions-no-mode-auto-approves)
>
> Entries with a specifier such as `Bash(ls *)` and the `acceptEdits` mode don't trigger it, and allow rules coming from settings files aren't visible to the check.
>
> Listen with `process.on('warning', ...)` and match the code to log or suppress it. To gate every tool call regardless of mode and rules, use a [`PreToolUse` hook](/docs/en/agent-sdk/hooks) instead.
>
> This page focuses on **allow and deny rules** and **permission modes**. For the other steps:
>
> * **Hooks:** run custom code to allow, deny, or modify tool requests. See [Control execution with hooks](/docs/en/agent-sdk/hooks).
> * **canUseTool callback:** prompt users for approval at runtime, when no earlier step resolves the call. See [Handle approvals and user input](/docs/en/agent-sdk/user-input).
>
> ## Allow and deny rules
>
> `allowed_tools` and `disallowed_tools` (TypeScript: `allowedTools` / `disallowedTools`) add entries to the allow and deny rule lists in the evaluation flow above. If you name one of the [task-tracking tools](/docs/en/agent-sdk/todo-tracking#model-availability) in `allowed_tools`, Claude Code also opts the session in. Any other tool not listed in `allowed_tools` is still available to Claude and falls through to the permission mode. Deny rules behave differently depending on whether they name a tool or scope a pattern within one.
>
> | Option                            | Effect                                                                                                                                                                             |
> | :-------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `allowed_tools=["Read", "Grep"]`  | `Read` and `Grep` are auto-approved. Other tools not listed here still exist and fall through to the permission mode and `canUseTool`.                                             |
> | `disallowed_tools=["Bash"]`       | The `Bash` tool definition is removed from the request. Claude does not see the tool and cannot attempt it.                                                                        |
> | `disallowed_tools=["Bash(rm *)"]` | `Bash` stays available. Calls matching `rm *` are denied in every permission mode, including `bypassPermissions`. Other `Bash` calls fall through to the permission mode.          |
> | `disallowed_tools=["*"]`          | Every tool definition is removed from the request. Tool-name globs are supported in deny rules: `"*"` matches every tool and `"mcp__*"` matches every MCP tool across all servers. |
>
> Allow rules accept tool-name globs only after a literal `mcp__<server>__` prefix. The server segment must be glob-free so the rule names a specific server you configured: `mcp__puppeteer__*` matches every tool from the `puppeteer` server, and `mcp__github__get_*` matches its `get_` tools. An unanchored entry like `allowed_tools=["*"]` or `allowed_tools=["mcp__*"]` is ignored with a startup warning and does not auto-approve anything.
>
> Scoped rules for `Read` and `Edit` take a path pattern. `Edit(path)` rules govern all built-in tools that write files, including `Write` and `NotebookEdit`; a `Write(path)` rule is never matched by the file permission checks.
>
> Use `//path` for an absolute filesystem path: a deny rule of `Edit(//secrets/**)` blocks writes anywhere under `/secrets` on disk. With a single leading slash, `Edit(/secrets/**)` anchors at the rule's source instead. For rules passed through `allowed_tools` or `disallowed_tools`, that means the session's working directory, so the rule doesn't block `/secrets` on disk. See [Read and Edit rules](/docs/en/permissions#read-and-edit) for the four anchor forms and how rules from settings files resolve.
>
>   **Auto-approved tools never reach `canUseTool`.** A tool call approved at any earlier step, by `acceptEdits` or `bypassPermissions`, or by an allow rule, skips your `canUseTool` callback, so permission checks you put there are silently bypassed for that tool. `AskUserQuestion`, MCP tools marked [`_meta["anthropic/requiresUserInteraction"]`](/docs/en/mcp#require-approval-for-a-specific-tool), connector tools [your organization set to `ask`](/docs/en/mcp#organization-controls-on-connector-tools), and `rm` and `rmdir` removals targeting a [critical path](/docs/en/permission-modes#critical-paths) still reach the callback, even when an allow rule matches. In `auto` mode, critical-path removals go to the [classifier](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) instead of the callback, while the other calls listed here still reach it; the classifier routing requires Claude Code v2.1.218 or later. In `dontAsk` mode these calls are denied instead, without invoking the callback.
>
>   Coverage depends on the entry's form: a bare name like `Read` or `mcp__github__get_issue` auto-approves every call to that tool apart from the exceptions above, while a scoped rule like `Bash(ls *)` auto-approves only matching calls and other `Bash` calls still fall through to the callback. For checks that must run on every tool call, use a [`PreToolUse` hook](/docs/en/agent-sdk/hooks): hooks run before every other step, and a hook deny applies even in `bypassPermissions` mode.
>
> For a locked-down agent, pair `allowedTools` with `permissionMode: "dontAsk"`. Listed tools are approved, apart from the always-prompt tools in the Warning above; anything else is denied outright instead of prompting:
>
> ```typescript
> const options = {
>   allowedTools: ["Read", "Glob", "Grep"],
>   permissionMode: "dontAsk"
> };
> ```
>
>   **`allowed_tools` does not constrain `bypassPermissions`.** `allowed_tools` pre-approves the tools you list. Other unlisted tools are not matched by any allow rule and fall through to the permission mode, where `bypassPermissions` approves them. Setting `allowed_tools=["Read"]` alongside `permission_mode="bypassPermissions"` still approves every tool, including `Bash`, `Write`, and `Edit`. If you need `bypassPermissions` but want specific tools blocked, use `disallowed_tools`.
>
> You can also configure allow, deny, and ask rules declaratively in `.claude/settings.json`. These rules are read when the `project` setting source is enabled, which it is for default `query()` options. If you set `setting_sources` (TypeScript: `settingSources`) explicitly, include `"project"` for them to apply. See [Permission settings](/docs/en/settings-reference#permission-settings) for the rule syntax.
>
> ## Permission modes
>
> Permission modes provide global control over how Claude uses tools. You can set the permission mode when calling `query()` or change it dynamically during streaming sessions.
>
> ### Available modes
>
> The SDK supports these permission modes:
>
> | Mode                | Description                  | Tool behavior                                                                                                                                                                                                                                                                                                                                                                            |
> | :------------------ | :--------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `default`           | Standard permission behavior | No auto-approvals; unmatched tools trigger your `canUseTool` callback                                                                                                                                                                                                                                                                                                                    |
> | `dontAsk`           | Deny instead of prompting    | Anything not pre-approved by `allowed_tools` or rules is denied; connector tools [your organization set to `ask`](/docs/en/mcp#organization-controls-on-connector-tools) and tools that require user interaction are denied even if you've pre-approved them, as are `rm` and `rmdir` removals targeting a [critical path](/docs/en/permission-modes#critical-paths). `canUseTool` is never called |
> | `acceptEdits`       | Auto-accept file edits       | File edits and [filesystem operations](#accept-edits-mode-acceptedits) (`mkdir`, `rm`, `mv`, etc.) are automatically approved                                                                                                                                                                                                                                                            |
> | `bypassPermissions` | Bypass permission checks     | Tools run without permission prompts, except for the [actions no mode auto-approves](/docs/en/permission-modes#actions-no-mode-auto-approves). Use with caution                                                                                                                                                                                                                               |
> | `plan`              | Planning mode                | Claude explores and plans without editing your source files; file edits are never auto-approved and prompt through your `canUseTool` callback                                                                                                                                                                                                                                            |
> | `auto`              | Model-classified approvals   | A model classifier approves or denies permission prompts. See [Auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) for availability                                                                                                                                                                                                                                        |
>
>   **Subagent inheritance:** Subagents inherit the parent session's permission mode. An [`AgentDefinition`'s `permissionMode`](/docs/en/agent-sdk/typescript#agentdefinition) can override it, except when the parent uses `bypassPermissions`, `acceptEdits`, or `auto`: those modes apply to every subagent and can't be overridden per subagent. Claude Code also ignores a definition's `permissionMode: "bypassPermissions"` when bypass mode is disabled by [`permissions.disableBypassPermissionsMode`](/docs/en/permissions#managed-settings), so that subagent runs with the parent session's mode.
>
>   Subagents may have different system prompts and less constrained behavior than your main agent, so inheriting `bypassPermissions` grants them full, autonomous system access. The [actions no mode auto-approves](/docs/en/permission-modes#actions-no-mode-auto-approves) still apply.
>
> ### Set permission mode
>
> You can set the permission mode once when starting a query, or change it dynamically while the session is active.
>
>   #### At query time
>
> Pass `permission_mode` (Python) or `permissionMode` (TypeScript) when creating a query. This mode applies for the entire session unless changed dynamically.
>
>     <CodeGroup>
>       ```python Python
>       import asyncio
>       from claude_agent_sdk import query, ClaudeAgentOptions
>
>       async def main():
>           async for message in query(
>               prompt="Help me refactor this code",
>               options=ClaudeAgentOptions(
>                   permission_mode="default",  # Set the mode here
>               ),
>           ):
>               if hasattr(message, "result"):
>                   print(message.result)
>
>       asyncio.run(main())
>       ```
>
>       ```typescript TypeScript
>       import { query } from "@anthropic-ai/claude-agent-sdk";
>
>       async function main() {
>         for await (const message of query({
>           prompt: "Help me refactor this code",
>           options: {
>             permissionMode: "default" // Set the mode here
>           }
>         })) {
>           if ("result" in message) {
>             console.log(message.result);
>           }
>         }
>       }
>
>       main();
>       ```
>     </CodeGroup>
>
>   #### During streaming
>
> Call `set_permission_mode()` (Python) or `setPermissionMode()` (TypeScript) to change the mode mid-session. The new mode takes effect immediately for all subsequent tool requests. This lets you start restrictive and loosen permissions as trust builds, for example switching to `acceptEdits` after reviewing Claude's initial approach.
>
>     <CodeGroup>
>       ```python Python
>       import asyncio
>       from claude_agent_sdk import ClaudeSDKClient, ClaudeAgentOptions
>
>       async def main():
>           async with ClaudeSDKClient(
>               options=ClaudeAgentOptions(
>                   permission_mode="default",  # Start in default mode
>               )
>           ) as client:
>               await client.query("Help me refactor this code")
>
>               # Change mode dynamically mid-session
>               await client.set_permission_mode("acceptEdits")
>
>               # Process messages with the new permission mode
>               async for message in client.receive_response():
>                   if hasattr(message, "result"):
>                       print(message.result)
>
>       asyncio.run(main())
>       ```
>
>       ```typescript TypeScript
>       import { query } from "@anthropic-ai/claude-agent-sdk";
>
>       async function main() {
>         const q = query({
>           prompt: "Help me refactor this code",
>           options: {
>             permissionMode: "default" // Start in default mode
>           }
>         });
>
>         // Change mode dynamically mid-session
>         await q.setPermissionMode("acceptEdits");
>
>         // Process messages with the new permission mode
>         for await (const message of q) {
>           if ("result" in message) {
>             console.log(message.result);
>           }
>         }
>       }
>
>       main();
>       ```
>     </CodeGroup>
>
> ### Mode details
>
> #### Accept edits mode (`acceptEdits`)
>
> Auto-approves file operations so Claude can edit code without prompting. Other tools (like Bash commands that aren't filesystem operations) still require normal permissions.
>
> **Auto-approved operations:**
>
> * File edits (Edit, Write tools)
> * Filesystem commands: `mkdir`, `touch`, `rm`, `rmdir`, `mv`, `cp`, `sed`
>
> Both apply only to paths inside the working directory or `additionalDirectories`. Paths outside that scope, writes to protected paths, and `rm` and `rmdir` removals targeting a [critical path](/docs/en/permission-modes#critical-paths) still prompt.
>
> **Use when:** you trust Claude's edits and want faster iteration, such as during prototyping or when working in an isolated directory.
>
> #### Don't ask mode (`dontAsk`)
>
> Converts any permission prompt into a denial. Tools pre-approved by `allowed_tools`, `settings.json` allow rules, or a hook run as normal. Connector tools [your organization set to `ask`](/docs/en/mcp#orga### Source: Agent SDK — Claude Code features
>
> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Use Claude Code features in the SDK
>
> > Load project instructions, skills, hooks, and other Claude Code features into your SDK agents.
>
> The Agent SDK is built on the same foundation as Claude Code, which means your SDK agents have access to the same filesystem-based features: project instructions (`CLAUDE.md` and rules), skills, hooks, and more.
>
> When you omit `settingSources`, `query()` reads the same filesystem settings as the Claude Code CLI: user, project, and local settings, CLAUDE.md files, and `.claude/` skills, agents, and commands. To run without these, pass `settingSources: []`, which limits the agent to what you configure programmatically. Managed policy settings and the global `~/.claude.json` config are read regardless of this option. See [What settingSources does not control](#what-settingsources-does-not-control).
>
> For a conceptual overview of what each feature does and when to use it, see [Extend Claude Code](/docs/en/features-overview).
>
> ## Control filesystem settings with settingSources
>
> The setting sources option ([`setting_sources`](/docs/en/agent-sdk/python#claudeagentoptions) in Python, [`settingSources`](/docs/en/agent-sdk/typescript#settingsource) in TypeScript) controls which filesystem-based settings the SDK loads. Pass an explicit list to opt in to specific sources, or pass an empty array to disable user, project, and local settings.
>
> This example loads both user-level and project-level settings by setting `settingSources` to `["user", "project"]`:
>
>   ```python Python
>   from claude_agent_sdk import query, ClaudeAgentOptions, AssistantMessage, ResultMessage
>   import asyncio
>
>   async def main():
>       async for message in query(
>           prompt="Help me refactor the auth module",
>           options=ClaudeAgentOptions(
>               # "user" loads from ~/.claude/, "project" loads from ./.claude/ in cwd.
>               # Together they give the agent access to CLAUDE.md, skills, hooks, and
>               # permissions from both locations.
>               setting_sources=["user", "project"],
>               allowed_tools=["Read", "Edit", "Bash"],
>           ),
>       ):
>           if isinstance(message, AssistantMessage):
>               for block in message.content:
>                   if hasattr(block, "text"):
>                       print(block.text)
>           if isinstance(message, ResultMessage) and message.subtype == "success":
>               print(f"\nResult: {message.result}")
>
>   asyncio.run(main())
>   ```
>
>   ```typescript TypeScript
>   import { query } from "@anthropic-ai/claude-agent-sdk";
>
>   for await (const message of query({
>     prompt: "Help me refactor the auth module",
>     options: {
>       // "user" loads from ~/.claude/, "project" loads from ./.claude/ in cwd.
>       // Together they give the agent access to CLAUDE.md, skills, hooks, and
>       // permissions from both locations.
>       settingSources: ["user", "project"],
>       allowedTools: ["Read", "Edit", "Bash"]
>     }
>   })) {
>     if (message.type === "assistant") {
>       for (const block of message.message.content) {
>         if (block.type === "text") console.log(block.text);
>       }
>     }
>     if (message.type === "result" && message.subtype === "success") {
>       console.log(`\nResult: ${message.result}`);
>     }
>   }
>   ```
>
> When this runs, the assistant's response prints to stdout, followed by a final result line once the run completes.
>
> Each source loads settings from a specific location, where `<cwd>` is the working directory you pass via the `cwd` option, or the process's current directory if unset. For the full type definition, see [`SettingSource`](/docs/en/agent-sdk/typescript#settingsource) (TypeScript) or [`SettingSource`](/docs/en/agent-sdk/python#settingsource) (Python).
>
> | Source      | What it loads                                                                                                          | Location                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
> | :---------- | :--------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `"project"` | Project `settings.json` and hooks; project CLAUDE.md and `.claude/rules/*.md`; project skills, commands, and subagents | `<cwd>/.claude/` for `settings.json` and hooks; `<cwd>` and every parent directory for CLAUDE.md and rules; `<cwd>` and every parent directory up to the repository root for skills, commands, and subagents, plus the `.claude/skills/`, `.claude/commands/`, and `.claude/agents/` folders of each directory you pass through the `additionalDirectories` or `add_dirs` option, which the SDK passes to Claude Code as [`--add-dir`](/docs/en/permissions#additional-directories-grant-file-access-not-configuration) |
> | `"user"`    | User `settings.json`; user CLAUDE.md and `~/.claude/rules/*.md`; user skills, commands, and subagents                  | `~/.claude/` for `settings.json`, CLAUDE.md, and rules; `~/.claude/skills/`, `~/.claude/commands/`, and `~/.claude/agents/` for skills, commands, and subagents                                                                                                                                                                                                                                                                                                                                                    |
> | `"local"`   | CLAUDE.local.md, `.claude/settings.local.json`                                                                         | `<cwd>/.claude/` for `settings.local.json`; `<cwd>` and every parent directory for CLAUDE.local.md                                                                                                                                                                                                                                                                                                                                                                                                                 |
>
> Omitting `settingSources` is equivalent to `["user", "project", "local"]`.
>
> The `cwd` option determines where the SDK looks for project-level inputs. Project `settings.json` and hooks load only from `<cwd>/.claude/` with no parent-directory fallback.
>
> ### What settingSources does not control
>
> `settingSources` covers user, project, and local settings. A few inputs are read regardless of its value:
>
> | Input                                                              | Behavior                                                                                                                                                                                                                                                                                                                                                             | To disable                                                                                                                                                                         |
> | :----------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | Managed policy settings                                            | Endpoint-managed policy, such as an MDM plist, registry policy, or managed settings file, loads from the host. [Server-managed settings](/docs/en/server-managed-settings) are fetched on an [eligible configuration](/docs/en/server-managed-settings#platform-availability) when the session authenticates with an organization OAuth login or a directly configured API key | Endpoint policy: remove the managed settings file, plist, or registry policy from the host. Server-managed settings: controlled by your org admin; cannot be disabled from the SDK |
> | `~/.claude.json` global config                                     | Always read                                                                                                                                                                                                                                                                                                                                                          | Relocate with `CLAUDE_CONFIG_DIR` in `env`                                                                                                                                         |
> | Auto memory at `~/.claude/projects/<project>/memory/`              | Loaded into the system prompt at session start. The agent writes new memories there with the standard `Write` and `Edit` tools rather than a dedicated memory tool, so those tools must be enabled for the agent to save memories                                                                                                                                    | Set `autoMemoryEnabled: false` in settings, or `CLAUDE_CODE_DISABLE_AUTO_MEMORY=1` in `env`                                                                                        |
> | [claude.ai MCP connectors](/docs/en/mcp#use-mcp-servers-from-claude-ai) | Loaded when the session authenticates with your claude.ai login. Not loaded when `CLAUDE_CODE_OAUTH_TOKEN` holds a token from [`claude setup-token`](/docs/en/authentication#generate-a-long-lived-token), which can only make model requests. Passing `mcpServers: {}` does not suppress the connectors                                                                  | Set `strictMcpConfig: true`, [`disableClaudeAiConnectors: true`](/docs/en/mcp#disable-claude-ai-connectors) in settings, or `ENABLE_CLAUDEAI_MCP_SERVERS=false` in `env`                |
>
>   Do not rely on default `query()` options for multi-tenant isolation. Because the inputs above are read regardless of `settingSources`, an SDK process can pick up host-level configuration and per-directory memory. For multi-tenant deployments, run each tenant in its own filesystem and set `settingSources: []` plus `CLAUDE_CODE_DISABLE_AUTO_MEMORY=1` in `env`. [Server-managed settings](/docs/en/server-managed-settings) are fetched when the process authenticates with an organization credential; filesystem isolation does not remove them. See [Secure deployment](/docs/en/agent-sdk/secure-deployment).
>
> ## Project instructions (CLAUDE.md and rules)
>
> `CLAUDE.md` files and `.claude/rules/*.md` files give your agent persistent context about your project: coding conventions, build commands, architecture decisions, and instructions. When `settingSources` includes `"project"` (as in the example above), the SDK loads these files into context at session start. The agent then follows your project conventions without you repeating them in every prompt.
>
> ### CLAUDE.md load locations
>
> | Level                 | Location                                                                      | When loaded                                                                                         |
> | :-------------------- | :---------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------- |
> | Project (root)        | `<cwd>/CLAUDE.md` or `<cwd>/.claude/CLAUDE.md`                                | `settingSources` includes `"project"`                                                               |
> | Project rules         | `<cwd>/.claude/rules/*.md` and `.claude/rules/*.md` in every parent directory | `settingSources` includes `"project"`                                                               |
> | Project (parent dirs) | `CLAUDE.md` files in directories above `cwd`                                  | `settingSources` includes `"project"`, loaded at session start                                      |
> | Project (child dirs)  | `CLAUDE.md` files in subdirectories of `cwd`                                  | `settingSources` includes `"project"`, loaded on demand when the agent reads a file in that subtree |
> | Local                 | `<cwd>/CLAUDE.local.md` and `CLAUDE.local.md` in every parent directory       | `settingSources` includes `"local"`                                                                 |
> | User                  | `~/.claude/CLAUDE.md`                                                         | `settingSources` includes `"user"`                                                                  |
> | User rules            | `~/.claude/rules/*.md`                                                        | `settingSources` includes `"user"`                                                                  |
>
> All levels are additive: if both project and user CLAUDE.md files exist, the agent sees both. There is no hard precedence rule between levels; if instructions conflict, the outcome depends on how Claude interprets them. Write non-conflicting rules, or state precedence explicitly in the more specific file ("These project instructions override any conflicting user-level defaults").
>
>   You can also inject context directly via `systemPrompt` without using CLAUDE.md files. See [Modify system prompts](/docs/en/agent-sdk/modifying-system-prompts). Use CLAUDE.md when you want the same context shared between interactive Claude Code sessions and your SDK agents.
>
> For how to structure and organize CLAUDE.md content, see [Manage Claude's memory](/docs/en/memory).
>
> ## Skills
>
> Skills are markdown files that give your agent specialized knowledge and invocable workflows. Unlike `CLAUDE.md` (which loads every session), skills load on demand. The agent receives skill descriptions at startup and loads the full content when relevant.
>
> Skills are discovered from the filesystem through `settingSources`. When the `skills` option on `query()` is omitted, discovered user and project skills are enabled and the Skill tool is available, matching CLI behavior. To control which skills are enabled, pass `skills` as `"all"`, a list of skill names, or `[]` to disable all. When `skills` is set, the SDK adds the Skill tool to `allowedTools` automatically. If you also pass an explicit `tools` list, include `"Skill"` in that list so Claude can invoke skills.
>
>   ```python Python
>   from claude_agent_sdk import query, ClaudeAgentOptions, ResultMessage
>   import asyncio
>
>   # Skills in .claude/skills/ are discovered automatically
>   # when settingSources includes "project"
>   async def main():
>       async for message in query(
>           prompt="Review this PR using our code review checklist",
>           options=ClaudeAgentOptions(
>               setting_sources=["user", "project"],
>               skills="all",
>               allowed_tools=["Read", "Grep", "Glob"],
>           ),
>       ):
>           if isinstance(message, ResultMessage) and message.subtype == "success":
>               print(message.result)
>
>   asyncio.run(main())
>   ```
>
>   ```typescript TypeScript
>   import { query } from "@anthropic-ai/claude-agent-sdk";
>
>   // Skills in .claude/skills/ are discovered automatically
>   // when settingSources includes "project"
>   for await (const message of query({
>     prompt: "Review this PR using our code review checklist",
>     options: {
>       settingSources: ["user", "project"],
>       skills: "all",
>       allowedTools: ["Read", "Grep", "Glob"]
>     }
>   })) {
>     if (message.type === "result" && message.subtype === "success") {
>       console.log(message.result);
>     }
>   }
>   ```
>
>   Skills must be created as filesystem artifacts (`.claude/skills/<name>/SKILL.md`). The SDK does not have a programmatic API for registering skills. See [Agent Skills in the SDK](/docs/en/agent-sdk/skills) for full details.
>
> ## Hooks
>
> The SDK supports two ways to define hooks, and they run side by side:
>
> * **Filesystem hooks:** shell commands defined in `settings.json`, loaded when `settingSources` includes the relevant source. These are the same hooks you'd configure for [interactive Claude Code sessions](/docs/en/hooks-guide).
> * **Programmatic hooks:** callback functions passed directly to `query()`. These run in your application process and can return structured decisions. See [Control execution with hooks](/docs/en/agent-sdk/hooks).
>
> Both types execute during the same hook lifecycle. If you already have hooks in your project's `.claude/settings.json` and you set `settingSources: ["project"]`, those hooks run automatically in the SDK with no extra configuration.
>
> Hook callbacks receive the tool input and return a decision dict. Returning `{}` means allow the tool to proceed. To block execution, return a `hookSpecificOutput` object with `permissionDecision: "deny"` and a `permissionDecisionReason`. The reason is sent to Claude as the tool result. See the [hooks guide](/docs/en/agent-sdk/hooks) for the full callback signature and return types.
>
>   ```python Python
>   from claude_agent_sdk import query, ClaudeAgentOptions, HookMatcher, ResultMessage
>   import asyncio
>
>   # PreToolUse hook callback. Positional args:
>   #   input_data: HookInput dict with tool_name, tool_input, hook_event_name
>   #   tool_use_id: str | None, the ID of the tool call being intercepted
>   #   context: HookContext, reserved for future abort-signal support
>   async def audit_bash(input_data, tool_use_id, context):
>       command = input_data.get("tool_input", {}).get("command", "")
>       if "rm -rf" in command:
>           return {
>               "hookSpecificOutput": {
>                   "hookEventName": "PreToolUse",
>                   "permissionDecision": "deny",
>                   "permissionDecisionReason": "Destructive command blocked",
>               }
>           }
>       return {}  # Empty dict: allow the tool to proceed
>
>   # Filesystem hooks from .claude/settings.json run automatically
>   # when settingSources loads them. You can also add programmatic hooks:
>   async def main():
>       async for message in query(
>           prompt="Refactor the auth module",
>           options=ClaudeAgentOptions(
>               setting_sources=["project"],  # Loads hooks from .claude/settings.json
>               hooks={
>                   "PreToolUse": [
>                       HookMatcher(matcher="Bash", hooks=[audit_bash]),
>                   ]
>               },
>           ),
>       ):
>           if isinstance(message, ResultMessage) and message.subtype == "success":
>               print(message.result)
>
>   asyncio.run(main())
>   ```
>
>   ```typescript TypeScript
>   import { query, type HookInput, type HookJSONOutput } from "@anthropic-ai/claude-agent-sdk";
>
>   // PreT