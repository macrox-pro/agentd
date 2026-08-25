---
primary_sources:
  - id: T1-MCP
    title: "Model Context Protocol"
    url: "https://learn.chatgpt.com/docs/extend/mcp.md"
    section: "Plugin-provided MCP servers; Examples of useful MCP servers"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Plugin-provided MCP servers

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Model Context Protocol — Plugin-provided MCP servers / Examples

> ### Plugin-provided MCP servers
>
> Installed plugins can bundle MCP servers in their plugin manifest. Those
> servers are launched from the plugin, so user config doesn't set their
> transport command. User config can still control on/off state and tool policy
> under `plugins.<plugin>.mcp_servers.<server>`.
>
> ```toml
> [plugins."sample@test".mcp_servers.sample]
> enabled = true
> default_tools_approval_mode = "prompt"
> enabled_tools = ["read", "search"]
>
> [plugins."sample@test".mcp_servers.sample.tools.search]
> approval_mode = "approve"
> ```
>
> ## Examples of useful MCP servers
>
> The list of MCP servers keeps growing. Here are a few common ones:
>
> - [OpenAI Docs MCP](https://developers.openai.com/learn/docs-mcp): Search and read OpenAI developer docs.
> - [Context7](https://github.com/upstash/context7): Connect to up-to-date developer documentation.
> - Figma [Local](https://developers.figma.com/docs/figma-mcp-server/local-server-installation/) and [Remote](https://developers.figma.com/docs/figma-mcp-server/remote-server-installation/): Access your Figma designs.
> - [Playwright](https://www.npmjs.com/package/@playwright/mcp): Control and inspect a browser using Playwright.
> - [Chrome Developer Tools](https://github.com/ChromeDevTools/chrome-devtools-mcp/): Control and inspect Chrome.
> - [Sentry](https://docs.sentry.io/product/sentry-mcp/#codex): Access Sentry logs.
> - [GitHub](https://github.com/github/github-mcp-server): Manage GitHub beyond what `git` supports (for example, pull requests and issues).
