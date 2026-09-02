---
primary_sources:
  - id: T1-MCP
    title: "MCP servers"
    url: "https://opencode.ai/docs/mcp-servers.md"
    section: "Enable through Remote"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# MCP configuration

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode MCP — Enable, Local, Remote

> ## Enable
>
> You can define MCP servers in your [OpenCode Config](https://opencode.ai/docs/config/) under `mcp`. Add each MCP with a unique name. You can refer to that MCP by name when prompting the LLM.
>
> ```jsonc title="opencode.jsonc" {6}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "name-of-mcp-server": {
>       // ...
>       "enabled": true,
>     },
>     "name-of-other-mcp-server": {
>       // ...
>     },
>   },
> }
> ```
>
> You can also disable a server by setting `enabled` to `false`. This is useful if you want to temporarily disable a server without removing it from your config.
>
> ---
>
> ### Overriding remote defaults
>
> Organizations can provide default MCP servers via their `.well-known/opencode` endpoint. These servers may be disabled by default, allowing users to opt-in to the ones they need.
>
> To enable a specific server from your organization's remote config, add it to your local config with `enabled: true`:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "jira": {
>       "type": "remote",
>       "url": "https://jira.example.com/mcp",
>       "enabled": true
>     }
>   }
> }
> ```
>
> Your local config values override the remote defaults. See [config precedence](/docs/config#precedence-order) for more details.
>
> ---
>
> ---
>
> ## Local
>
> Add local MCP servers using `type` to `"local"` within the MCP object.
>
> ```jsonc title="opencode.jsonc" {15}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "my-local-mcp-server": {
>       "type": "local",
>       // Or ["bun", "x", "my-mcp-command"]
>       "command": ["npx", "-y", "my-mcp-command"],
>       "enabled": true,
>       "environment": {
>         "MY_ENV_VAR": "my_env_var_value",
>       },
>     },
>   },
> }
> ```
>
> The command is how the local MCP server is started. You can also pass in a list of environment variables as well.
>
> For example, here's how you can add the test [`@modelcontextprotocol/server-everything`](https://www.npmjs.com/package/@modelcontextprotocol/server-everything) MCP server.
>
> ```jsonc title="opencode.jsonc"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "mcp_everything": {
>       "type": "local",
>       "command": ["npx", "-y", "@modelcontextprotocol/server-everything"],
>     },
>   },
> }
> ```
>
> And to use it I can add `use the mcp_everything tool` to my prompts.
>
> ```txt "mcp_everything"
> use the mcp_everything tool to add the number 3 and 4
> ```
>
> ---
>
> #### Options
>
> Here are all the options for configuring a local MCP server.
>
> | Option        | Type    | Required | Description                                                                              |
> | ------------- | ------- | -------- | ---------------------------------------------------------------------------------------- |
> | `type`        | String  | Y        | Type of MCP server connection, must be `"local"`.                                        |
> | `command`     | Array   | Y        | Command and arguments to run the MCP server.                                             |
> | `cwd`         | String  |          | Working directory for the MCP server process. Relative paths resolve from the workspace. |
> | `environment` | Object  |          | Environment variables to set when running the server.                                    |
> | `enabled`     | Boolean |          | Enable or disable the MCP server on startup.                                             |
> | `timeout`     | Number  |          | Timeout in ms for fetching tools from the MCP server. Defaults to 5000 (5 seconds).      |
>
> ---
>
> ---
>
> ## Remote
>
> Add remote MCP servers by setting `type` to `"remote"`.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "my-remote-mcp": {
>       "type": "remote",
>       "url": "https://my-mcp-server.com",
>       "enabled": true,
>       "headers": {
>         "Authorization": "Bearer MY_API_KEY"
>       }
>     }
>   }
> }
> ```
>
> The `url` is the URL of the remote MCP server and with the `headers` option you can pass in a list of headers.
>
> ---
>
> #### Options
>
> | Option    | Type    | Required | Description                                                                         |
> | --------- | ------- | -------- | ----------------------------------------------------------------------------------- |
> | `type`    | String  | Y        | Type of MCP server connection, must be `"remote"`.                                  |
> | `url`     | String  | Y        | URL of the remote MCP server.                                                       |
> | `enabled` | Boolean |          | Enable or disable the MCP server on startup.                                        |
> | `headers` | Object  |          | Headers to send with the request.                                                   |
> | `oauth`   | Object  |          | OAuth authentication configuration. See [OAuth](#oauth) section below.              |
> | `timeout` | Number  |          | Timeout in ms for fetching tools from the MCP server. Defaults to 5000 (5 seconds). |
>
> ---
