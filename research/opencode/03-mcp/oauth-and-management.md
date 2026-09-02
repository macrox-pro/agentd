---
primary_sources:
  - id: T1-MCP
    title: "MCP servers"
    url: "https://opencode.ai/docs/mcp-servers.md"
    section: "OAuth"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# MCP OAuth and management

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode MCP — OAuth and CLI management

> ## OAuth
>
> OpenCode automatically handles OAuth authentication for remote MCP servers. When a server requires authentication, OpenCode will:
>
> 1. Detect the 401 response and initiate the OAuth flow
> 2. Use **Dynamic Client Registration (RFC 7591)** if supported by the server
> 3. Store tokens securely for future requests
>
> ---
>
> ### Automatic
>
> For most OAuth-enabled MCP servers, no special configuration is needed. Just configure the remote server:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "my-oauth-server": {
>       "type": "remote",
>       "url": "https://mcp.example.com/mcp"
>     }
>   }
> }
> ```
>
> If the server requires authentication, OpenCode will prompt you to authenticate when you first try to use it. If not, you can [manually trigger the flow](#authenticating) with `opencode mcp auth <server-name>`.
>
> ---
>
> ### Pre-registered
>
> If you have client credentials from the MCP server provider, you can configure them:
>
> ```json title="opencode.json" {7-11}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "my-oauth-server": {
>       "type": "remote",
>       "url": "https://mcp.example.com/mcp",
>       "oauth": {
>         "clientId": "{env:MY_MCP_CLIENT_ID}",
>         "clientSecret": "{env:MY_MCP_CLIENT_SECRET}",
>         "scope": "tools:read tools:execute"
>       }
>     }
>   }
> }
> ```
>
> ---
>
> ### Authenticating
>
> You can manually trigger authentication or manage credentials.
>
> Authenticate with a specific MCP server:
>
> ```bash
> opencode mcp auth my-oauth-server
> ```
>
> List all MCP servers and their auth status:
>
> ```bash
> opencode mcp list
> ```
>
> Remove stored credentials:
>
> ```bash
> opencode mcp logout my-oauth-server
> ```
>
> The `mcp auth` command will open your browser for authorization. After you authorize, OpenCode will store the tokens securely in `~/.local/share/opencode/mcp-auth.json`.
>
> ---
>
> #### Disabling OAuth
>
> If you want to disable automatic OAuth for a server (e.g., for servers that use API keys instead), set `oauth` to `false`:
>
> ```json title="opencode.json" {7}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "my-api-key-server": {
>       "type": "remote",
>       "url": "https://mcp.example.com/mcp",
>       "oauth": false,
>       "headers": {
>         "Authorization": "Bearer {env:MY_API_KEY}"
>       }
>     }
>   }
> }
> ```
>
> ---
>
> #### OAuth Options
>
> | Option         | Type            | Description                                                                      |
> | -------------- | --------------- | -------------------------------------------------------------------------------- |
> | `oauth`        | Object \| false | OAuth config object, or `false` to disable OAuth auto-detection.                 |
> | `clientId`     | String          | OAuth client ID. If not provided, dynamic client registration will be attempted. |
> | `clientSecret` | String          | OAuth client secret, if required by the authorization server.                    |
> | `scope`        | String          | OAuth scopes to request during authorization.                                    |
>
> #### Debugging
>
> If a remote MCP server is failing to authenticate, you can diagnose issues with:
>
> ```bash
> # View auth status for all OAuth-capable servers
> opencode mcp auth list
>
> # Debug connection and OAuth flow for a specific server
> opencode mcp debug my-oauth-server
> ```
>
> The `mcp debug` command shows the current auth status, tests HTTP connectivity, and attempts the OAuth discovery flow.
>
> ---
>
> ---
>
> ## Manage
>
> Your MCPs are available as tools in OpenCode, alongside built-in tools. So you can manage them through the OpenCode config like any other tool.
>
> ---
>
> ### Global
>
> This means that you can enable or disable them globally.
>
> ```json title="opencode.json" {14}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "my-mcp-foo": {
>       "type": "local",
>       "command": ["bun", "x", "my-mcp-command-foo"]
>     },
>     "my-mcp-bar": {
>       "type": "local",
>       "command": ["bun", "x", "my-mcp-command-bar"]
>     }
>   },
>   "tools": {
>     "my-mcp-foo": false
>   }
> }
> ```
>
> We can also use a glob pattern to disable all matching MCPs.
>
> ```json title="opencode.json" {14}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "my-mcp-foo": {
>       "type": "local",
>       "command": ["bun", "x", "my-mcp-command-foo"]
>     },
>     "my-mcp-bar": {
>       "type": "local",
>       "command": ["bun", "x", "my-mcp-command-bar"]
>     }
>   },
>   "tools": {
>     "my-mcp*": false
>   }
> }
> ```
>
> Here we are using the glob pattern `my-mcp*` to disable all MCPs.
>
> ---
>
> ### Per agent
>
> If you have a large number of MCP servers you may want to only enable them per agent and disable them globally. To do this:
>
> 1. Disable it as a tool globally.
> 2. In your [agent config](/docs/agents#tools), enable the MCP server as a tool.
>
> ```json title="opencode.json" {11, 14-18}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "my-mcp": {
>       "type": "local",
>       "command": ["bun", "x", "my-mcp-command"],
>       "enabled": true
>     }
>   },
>   "tools": {
>     "my-mcp*": false
>   },
>   "agent": {
>     "my-agent": {
>       "tools": {
>         "my-mcp*": true
>       }
>     }
>   }
> }
> ```
>
> ---
>
> #### Glob patterns
>
> The glob pattern uses simple regex globbing patterns:
>
> - `*` matches zero or more of any character (e.g., `"my-mcp*"` matches `my-mcp_search`, `my-mcp_list`, etc.)
> - `?` matches exactly one character
> - All other characters match literally
>
> :::note
> MCP server tools are registered with server name as prefix, so to disable all tools for a server simply use:
>
> ```
> "mymcpservername_*": false
> ```
>
> :::
>
> ---
