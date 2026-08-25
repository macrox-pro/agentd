---
primary_sources:
  - id: T1-MCP
    title: "Model Context Protocol"
    url: "https://learn.chatgpt.com/docs/extend/mcp.md"
    section: "Connect Codex to an MCP server; Configure (app/CLI/IDE/config.toml); config.toml examples"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# MCP configuration

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Model Context Protocol — Connect / Configure / Examples

> # Model Context Protocol
>
> > For the complete documentation index, see [llms.txt](https://learn.chatgpt.com/llms.txt). Markdown versions of documentation pages are available by appending `.md` to the page URL.
>
> Model Context Protocol (MCP) connects models to tools and context. Use it to
> give ChatGPT or Codex access to third-party documentation, or to let it
> interact with developer tools like your browser or Figma.
>
> ChatGPT web can use remote MCP-backed tools supplied by plugins. Local Codex
> clients can also connect directly to MCP servers and share their configuration.
>
> ## Connect Codex to an MCP server
>
> Codex stores MCP configuration in `config.toml` alongside other Codex configuration settings. By default this is `~/.codex/config.toml`, but you can also scope MCP servers to a project with `.codex/config.toml` (trusted projects only).
>
> The ChatGPT desktop app, Codex CLI, and IDE extension share this configuration.
> Once you configure your MCP servers, you can switch among those clients without
> redoing setup.
>
> ### Configure in the ChatGPT desktop app
>
> 1. Open **Settings**, then select **MCP servers**.
> 2. Select **Add server**.
> 3. Enter a name, choose **STDIO** or **Streamable HTTP**, and provide the
>    server's command or URL.
> 4. Save the server, then select **Restart**.
>
> The server list shows which servers are enabled and which require OAuth. Select
> **Authenticate** when an OAuth server requires sign-in. In the composer, type `/mcp`
> to view connected servers.
>
> ## Use MCP-backed tools in ChatGPT web
>
> In a hosted ChatGPT Work chat, install a [plugin](https://learn.chatgpt.com/docs/plugins) to use its
> bundled connectors and remote MCP tools. After installation, Chat and Work can
> use those tools. Workspace administrators can control which plugins and tools
> are available.
>
> ChatGPT web doesn't read local Codex configuration files or expose the local
> Codex command menu. Open the **Plugins** tab to browse and manage available
> tools.
>
> ### Configure with the CLI
>
> #### Add an MCP server
>
> ```bash
> codex mcp add <server-name> --env VAR1=VALUE1 --env VAR2=VALUE2 -- <stdio server-command>
> ```
>
> For example, to add Context7 (a free MCP server for developer documentation), you can run the following command:
>
> ```bash
> codex mcp add context7 -- npx -y @upstash/context7-mcp
> ```
>
> #### Other CLI commands
>
> Run `codex mcp list` to see configured servers. To see all available MCP
> commands, run `codex mcp --help`. For a server that supports OAuth, run
> `codex mcp login <server-name>`.
>
> #### Terminal UI (TUI)
>
> In the `codex` TUI, use `/mcp` to see your active MCP servers.
>
> ### Configure in the IDE extension
>
> 1. Open the gear menu, then select **MCP servers**.
> 2. Select **Add server**.
> 3. Enter a name, choose **STDIO** or **Streamable HTTP**, and provide the
>    server's command or URL.
> 4. Save the server, then select **Restart extension**.
>
> The MCP server list shows which servers are enabled and which require OAuth.
> Select **Authenticate** when an OAuth server requires sign-in.
>
> ### Configure with config.toml
>
> For more fine-grained control, edit `~/.codex/config.toml` or a project-scoped
> `.codex/config.toml`. See the [configuration reference](https://learn.chatgpt.com/docs/config-file/config-reference)
> for a searchable list of every supported MCP option.
>
> Configure each MCP server with a `[mcp_servers.<server-name>]` table in the configuration file.
>
> #### config.toml examples
>
> ```toml
> [mcp_servers.context7]
> command = "npx"
> args = ["-y", "@upstash/context7-mcp"]
> env_vars = ["LOCAL_TOKEN"]
>
> [mcp_servers.context7.env]
> MY_ENV_VAR = "MY_ENV_VALUE"
> ```
>
> ```toml
> # Optional MCP OAuth callback overrides (used by `codex mcp login`)
> mcp_oauth_callback_port = 5555
> mcp_oauth_callback_url = "https://devbox.example.internal/callback"
> ```
>
> ```toml
> [mcp_servers.figma]
> url = "https://mcp.figma.com/mcp"
> bearer_token_env_var = "FIGMA_OAUTH_TOKEN"
> http_headers = { "X-Figma-Region" = "us-east-1" }
> ```
>
> ```toml
> [mcp_servers.chrome_devtools]
> url = "http://localhost:3000/mcp"
> enabled_tools = ["open", "screenshot"]
> disabled_tools = ["screenshot"] # applied after enabled_tools
> default_tools_approval_mode = "prompt"
> startup_timeout_sec = 20
> tool_timeout_sec = 45
> enabled = true
>
> [mcp_servers.chrome_devtools.tools.open]
> approval_mode = "approve"
> ```
