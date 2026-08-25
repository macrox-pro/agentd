---
primary_sources:
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "MCP tool hooks"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# MCP tool hooks

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex Hooks — MCP tool hooks

> ## MCP tool hooks
>
> An MCP tool hook lets a lifecycle event call a tool on an already-connected MCP
> server. It sends structured arguments directly to the tool and uses the same
> trust review and output contract as a command hook.
>
> ### Configure an MCP tool hook
>
> This hook asks the `scanner` MCP server to scan each patch after Codex writes or
> edits files:
>
> ```json
> {
>   "hooks": {
>     "PostToolUse": [
>       {
>         "matcher": "Write|Edit",
>         "hooks": [
>           {
>             "type": "mcp_tool",
>             "server": "scanner",
>             "tool": "scan_patch",
>             "input": { "patch": "${tool_input.command}" },
>             "timeout": 30,
>             "statusMessage": "Scanning edited files"
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> | Field           | Meaning                                                          |
> | --------------- | ---------------------------------------------------------------- |
> | `type`          | Must be `mcp_tool`.                                              |
> | `server`        | Required name of an already-connected MCP server.                |
> | `tool`          | Required name of a tool exposed by that server.                  |
> | `input`         | Optional JSON object of argument templates. Defaults to `{}`.    |
> | `timeout`       | Optional active execution timeout in seconds. Defaults to `600`. |
> | `statusMessage` | Optional message shown while the hook runs.                      |
>
> ### Expand arguments from the hook event
>
> Use `${field.nested}` to read a dotted field from the hook event. A placeholder
> that fills an entire value keeps its JSON type. A placeholder inside a larger
> string is rendered as text. Codex expands objects and arrays recursively.
>
> For an event containing `{"tool_input":{"file_path":"src/main.rs","count":3}}`,
> this argument template:
>
> ```json
> {
>   "path": "${tool_input.file_path}",
>   "count": "${tool_input.count}",
>   "message": "Scanning ${tool_input.file_path}"
> }
> ```
>
> becomes:
>
> ```json
> {
>   "path": "src/main.rs",
>   "count": 3,
>   "message": "Scanning src/main.rs"
> }
> ```
>
> ### Execution and lifecycle
>
> - Hooks use an existing MCP connection. They don't start or reconnect servers.
> - A hook can block an operation when the tool returns a blocking decision.
>   Errors, missing servers, and unavailable tools don't block the operation.
> - MCP tool hooks run synchronously. They don't request tool approval or trigger
>   other hooks.
> - The shorter hook or server timeout applies. Time spent waiting for an MCP
>   elicitation response doesn't count against the timeout.
> - `SessionStart` hooks can run before an MCP server is ready. If that happens,
>   they don't block the session.
> - `SessionEnd` doesn't support MCP tool hooks.
