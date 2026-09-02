---
primary_sources:
  - id: T1-MCP
    title: "MCP servers"
    url: "https://opencode.ai/docs/mcp-servers.md"
    section: "Manage"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# MCP tools integration

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode MCP — Manage tools and Examples

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
>
> ---
>
> ## Examples
>
> Below are examples of some common MCP servers. You can submit a PR if you want to document other servers.
>
> ---
>
> ### Sentry
>
> Add the [Sentry MCP server](https://mcp.sentry.dev) to interact with your Sentry projects and issues.
>
> ```json title="opencode.json" {4-8}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "sentry": {
>       "type": "remote",
>       "url": "https://mcp.sentry.dev/mcp",
>       "oauth": {}
>     }
>   }
> }
> ```
>
> After adding the configuration, authenticate with Sentry:
>
> ```bash
> opencode mcp auth sentry
> ```
>
> This will open a browser window to complete the OAuth flow and connect OpenCode to your Sentry account.
>
> Once authenticated, you can use Sentry tools in your prompts to query issues, projects, and error data.
>
> ```txt "use sentry"
> Show me the latest unresolved issues in my project. use sentry
> ```
>
> ---
>
> ### Context7
>
> Add the [Context7 MCP server](https://github.com/upstash/context7) to search through docs.
>
> ```json title="opencode.json" {4-7}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "context7": {
>       "type": "remote",
>       "url": "https://mcp.context7.com/mcp"
>     }
>   }
> }
> ```
>
> If you have signed up for a free account, you can use your API key and get higher rate-limits.
>
> ```json title="opencode.json" {7-9}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "context7": {
>       "type": "remote",
>       "url": "https://mcp.context7.com/mcp",
>       "headers": {
>         "CONTEXT7_API_KEY": "{env:CONTEXT7_API_KEY}"
>       }
>     }
>   }
> }
> ```
>
> Here we are assuming that you have the `CONTEXT7_API_KEY` environment variable set.
>
> Add `use context7` to your prompts to use Context7 MCP server.
>
> ```txt "use context7"
> Configure a Cloudflare Worker script to cache JSON API responses for five minutes. use context7
> ```
>
> Alternatively, you can add something like this to your [AGENTS.md](/docs/rules/).
>
> ```md title="AGENTS.md"
> When you need to search docs, use `context7` tools.
> ```
>
> ---
>
> ### Grep by Vercel
>
> Add the [Grep by Vercel](https://grep.app) MCP server to search through code snippets on GitHub.
>
> ```json title="opencode.json" {4-7}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {
>     "gh_grep": {
>       "type": "remote",
>       "url": "https://mcp.grep.app"
>     }
>   }
> }
> ```
>
> Since we named our MCP server `gh_grep`, you can add `use the gh_grep tool` to your prompts to get the agent to use it.
>
> ```txt "use the gh_grep tool"
> What's the right way to set a custom domain in an SST Astro component? use the gh_grep tool
> ```
>
> Alternatively, you can add something like this to your [AGENTS.md](/docs/rules/).
>
> ```md title="AGENTS.md"
> If you are unsure how to do something, use `gh_grep` to search code examples from GitHub.
> ```
