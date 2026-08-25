---
primary_sources:
  - id: T1-CFG-REF
    title: "Configuration Reference"
    url: "https://learn.chatgpt.com/docs/config-file/config-reference.md"
    section: "notify; features.hooks; hooks.*; mcp_servers.*"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Config reference — hooks and MCP

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Configuration Reference — notify, features.hooks, hooks.*, mcp_servers.*

> <ConfigTable
>   options={[
>     {
>       key: "notify",
>       type: "array<string>",
>       description:
>         "Command invoked for notifications; receives a JSON payload from Codex.",
>     },
>     {
>       key: "features.hooks",
>       type: "boolean",
>       description:
>         "Enable lifecycle hooks loaded from `hooks.json` or inline `[hooks]` config. `features.codex_hooks` is a deprecated alias.",
>     },
>     {
>       key: "hooks",
>       type: "table",
>       description:
>         "Lifecycle hooks configured inline in `config.toml`. Uses the same event schema as `hooks.json`; see the Hooks guide for examples and supported events.",
>     },
>     {
>       key: "hooks.<Event>",
>       type: "array<table>",
>       description:
>         "Matcher groups for hook events such as `PreToolUse`, `PermissionRequest`, `PostToolUse`, `PreCompact`, `PostCompact`, `SessionStart`, `SessionEnd`, `SubagentStart`, `SubagentStop`, `UserPromptSubmit`, or `Stop`.",
>     },
>     {
>       key: "hooks.<Event>[].hooks",
>       type: "array<table>",
>       description:
>         "Hook handlers for a matcher group. Command and MCP tool hooks are supported while prompt and agent hook handlers are parsed but skipped.",
>     },
>     {
>       key: "hooks.<Event>[].hooks[].async",
>       type: "boolean",
>       description:
>         "Run a command hook in the background without delaying the triggering operation. Defaults to `false`; `SessionEnd` always runs synchronously. See [Run hooks in the background](https://learn.chatgpt.com/docs/hooks#run-hooks-in-the-background).",
>     },
>     {
>       key: "hooks.<Event>[].hooks[].additionalContextLimit",
>       type: "integer",
>       description:
>         "Approximate per-handler token threshold for saving oversized `additionalContext` to disk and showing the model a shorter preview. Defaults to `2500`; `0` passes the full context directly to the model. See [Large hook output](https://learn.chatgpt.com/docs/hooks#large-hook-output).",
>     },
>     {
>       key: "hooks.<Event>[].hooks[].commandWindows",
>       type: "string",
>       description:
>         "Windows-only command override for command hooks. The TOML alias `command_windows` is also accepted.",
>     },
>     {
>       key: "mcp_servers.<id>.command",
>       type: "string",
>       description: "Launcher command for an MCP stdio server.",
>     },
>     {
>       key: "mcp_servers.<id>.args",
>       type: "array<string>",
>       description: "Arguments passed to the MCP stdio server command.",
>     },
>     {
>       key: "mcp_servers.<id>.env",
>       type: "map<string,string>",
>       description: "Environment variables forwarded to the MCP stdio server.",
>     },
>     {
>       key: "mcp_servers.<id>.env_vars",
>       type: 'array<string | { name = string, source = "local" | "remote" }>',
>       description:
>         'Additional environment variables to whitelist for an MCP stdio server. String entries default to `source = "local"`; use `source = "remote"` only with executor-backed remote stdio.',
>     },
>     {
>       key: "mcp_servers.<id>.cwd",
>       type: "string",
>       description: "Working directory for the MCP stdio server process.",
>     },
>     {
>       key: "mcp_servers.<id>.url",
>       type: "string",
>       description: "Endpoint for an MCP streamable HTTP server.",
>     },
>     {
>       key: "mcp_servers.<id>.auth",
>       type: "oauth | chatgpt",
>       description:
>         "Authentication fallback for an MCP HTTP server after configured bearer tokens and authorization headers. `oauth` (default) uses stored MCP OAuth credentials when available. `chatgpt` uses the current ChatGPT session for the trusted first-party ChatGPT origin, then falls back to stored OAuth. Both modes can connect without authentication if no credential source resolves.",
>     },
>     {
>       key: "mcp_servers.<id>.bearer_token_env_var",
>       type: "string",
>       description:
>         "Environment variable sourcing the bearer token for an MCP HTTP server.",
>     },
>     {
>       key: "mcp_servers.<id>.http_headers",
>       type: "map<string,string>",
>       description: "Static HTTP headers included with each MCP HTTP request.",
>     },
>     {
>       key: "mcp_servers.<id>.env_http_headers",
>       type: "map<string,string>",
>       description:
>         "HTTP headers populated from environment variables for an MCP HTTP server.",
>     },
>     {
>       key: "mcp_servers.<id>.enabled",
>       type: "boolean",
>       description: "Disable an MCP server without removing its configuration.",
>     },
>     {
>       key: "mcp_servers.<id>.required",
>       type: "boolean",
>       description:
>         "When true, fail startup/resume if this enabled MCP server cannot initialize.",
>     },
>     {
>       key: "mcp_servers.<id>.startup_timeout_sec",
>       type: "number",
>       description:
>         "Override the default 10s startup timeout for an MCP server.",
>     },
>     {
>       key: "mcp_servers.<id>.startup_timeout_ms",
>       type: "number",
>       description: "Alias for `startup_timeout_sec` in milliseconds.",
>     },
>     {
>       key: "mcp_servers.<id>.tool_timeout_sec",
>       type: "number",
>       description:
>         "Override the default 60s per-tool timeout for an MCP server.",
>     },
>     {
>       key: "mcp_servers.<id>.enabled_tools",
>       type: "array<string>",
>       description: "Allow list of tool names exposed by the MCP server.",
>     },
>     {
>       key: "mcp_servers.<id>.disabled_tools",
>       type: "array<string>",
>       description:
>         "Deny list applied after `enabled_tools` for the MCP server.",
>     },
>     {
>       key: "mcp_servers.<id>.default_tools_approval_mode",
>       type: "auto | prompt | writes | approve",
>       description:
>         "Default approval behavior for MCP tools on this server unless a per-tool override exists.",
>     },
>     {
>       key: "mcp_servers.<id>.tools.<tool>.approval_mode",
>       type: "auto | prompt | writes | approve",
>       description:
>         "Per-tool approval behavior override for one MCP tool on this server.",
>     },
>     {
>       key: "mcp_servers.<id>.scopes",
>       type: "array<string>",
>       description:
>         "OAuth scopes to request when authenticating to that MCP server.",
>     },
>     {
>       key: "mcp_servers.<id>.oauth_resource",
>       type: "string",
>       description:
>         "Optional RFC 8707 OAuth resource parameter to include during MCP login.",
>     },
>     {
>       key: "mcp_servers.<id>.experimental_environment",
>       type: "local | remote",
>       description:
>         "Experimental placement for an MCP server. `remote` starts stdio servers through a remote executor environment; streamable HTTP remote placement is not implemented.",
>     }
>   ]}
>   client:load
> />
