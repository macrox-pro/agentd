---
primary_sources:
  - id: T1-MCP
    title: "MCP"
    url: "https://code.claude.com/docs/en/mcp.md"
    section: ""
  - id: T1-DEBUG
    title: "Debug"
    url: "https://code.claude.com/docs/en/debug-your-config.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# MCP debug and approval

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: MCP — Approval

> approves that file.
>
> Each `claude mcp add` and `claude mcp add-json` command prints an `Added ...` line. To check that Claude Code connected, run `claude mcp get <name>`; [Server status](#server-status) covers the statuses it shows and the approval step for `.mcp.json` servers.
>
> ### Managing your servers
>
> Once configured, you can manage your MCP servers with these commands:
>
> ```bash
> # List all configured servers
> claude mcp list
>
> # Get details for a specific server
> claude mcp get notion
>
> # Remove a server
> claude mcp remove notion
>
> # (within Claude Code) Check server status
> /mcp
> ```
>
> #### Server status
>
> `claude mcp add` confirms a successful add by printing an `Added ...` line, which means the configuration was written. `claude mcp list` then shows a health status next to each server it lists, such as `✔ Connected`, `! Needs authentication`, or `✘ Failed to connect`. A failure status means Claude Code couldn't connect to that server, not that the list command failed.
>
> Project-scoped servers from `.mcp.json` that are awaiting your approval appear in `claude mcp list` and `claude mcp get <name>` as ``⏸ Pending approval (run `claude` to approve)``. Run `claude` interactively to review and approve them. `claude mcp get <name>` shows rejected servers as `✘ Rejected (see disabledMcpjsonServers in settings)`.
>
> WebSocket servers don't appear in `claude mcp list` output. Use `claude mcp get <name>` or the `/mcp` panel to check them.
>
> #### Project server approvals and workspace trust
>
> As of v2.1.196, `claude mcp list` and `claude mcp get` read `.mcp.json` approvals only from settings files that aren't checked into the repository until you trust the workspace by running `claude` in it and accepting the workspace trust dialog. A cloned repository can't approve its own servers: [`enableAllProjectMcpServers`](/docs/en/settings-reference#enableallprojectmcpservers) or [`enabledMcpjsonServers`](/docs/en/settings-reference#enabledmcpjsonservers) committed to the project's `.claude/settings.json` is ignored in an untrusted folder, and the server stays at `⏸ Pending approval` instead of being connected and health-checked.
>
> Approvals from these sources still apply in an untrusted folder:
>
> * your user `~/.claude/settings.json`
> * managed settings
> * settings passed with `--settings`
>
> Claude Code also applies approvals from an untracked `.claude/settings.local.json`, but it runs git to check whether the file is tracked, and it runs that check only in a [trusted folder](/docs/en/permissions#project-allow-rules-and-workspace-trust). In a folder you've never trusted, Claude Code waits for the trust dialog before applying the file's approvals, unless the folder is your own configuration home: your home directory, or a directory whose `.claude` you've set as [`CLAUDE_CONFIG_DIR`](/docs/en/env-vars). Before v2.1.207, Claude Code applied approvals from an untracked `.claude/settings.local.json` even in a folder you'd never trusted.
>
> A `disabledMcpjsonServers` entry in any settings file still rejects the server.
>
> #### Server status detail
>
> In `/mcp`, including a server's menu there, and in the [`/plugin`](/docs/en/plugins) manager, a remote HTTP or SSE server you've used before can show a `cached` status such as `cached 2h ago · connects on first use · 5 tools`. Claude Code loaded the server's tool list from its discovery cache, saved in a previous session, instead of connecting at startup, and Claude Code connects the server the first time Claude calls one of the server's tools. The tools are available from your first message, so you don't need to do anything. The discovery cache and its `cached` status require Claude Code v2.1.221 or later.
>
> The discovery cache is off by default unless a gradual rollout has enabled it for your account. Set [`MCP_DISCOVERY_CACHE=1`](/docs/en/env-vars) to turn it on, or `0` to keep it off even when the rollout has enabled it. Before v2.1.238, the cache was on by default.
>
> Two actions in a server's menu in `/mcp` also affect that server's cache entry:
>
> * **Reconnect**: on a `cached` server, Claude Code connects it now rather than on its first tool call and keeps the entry. On a connected or failed server, Claude Code reconnects it and also discards the entry.
> * **Clear authentication**: Claude Code revokes the server's authentication and also discards the entry.
>
> After discarding the entry, Claude Code fetches the server's tool list from the server instead of from the cache.
>
> When a server's status is `✘ Failed to connect`, `claude mcp list` appends the failure detail to that status line, and `claude mcp get <name>` shows it on an `Issue:` line: the HTTP status or error code, plus any error text the server returned. The server's detail view in `/mcp` includes the same server-reported text in its `Issue:` row. Claude Code redacts credential-like text from this detail and never includes the expanded server URL, which can carry secrets. Claude Code appends no detail to a `✘ Connection error` status, because the exception text it would print there can embed that URL. Before v2.1.219, both commands showed only the bare failure status, without the status code or the server's error text.
>
> A remote server whose configuration has an empty `url` shows as `not configured` in `/mcp`, in `claude mcp list`, and in the [`/plugin`](/docs/en/plugins) manager, and Claude Code doesn't attempt to connect to it. A plugin can include a placeholder entry like this for a connector you configure later, so Claude Code doesn't report it as an error or a setup issue. The server's detail view in `/mcp` reads `No URL configured for this server`; set the entry's `url` to connect it. Before v2.1.208, Claude Code reported an empty `url` as a configuration issue with a prompt to reconnect.
>
> #### Configuration warnings
>
> Claude Code also warns when an MCP config value carries hidden leading or trailing whitespace, which often comes from pasting a token with a trailing newline. Claude Code checks `command`, `url`, each `args` entry, and the values and key names under `env` and `headers`. Claude Code shows the warning in `claude mcp list` output and in `/mcp`, naming the affected fields without echoing their values, for example `Leading or trailing whitespace in: headers.Authorization`. Claude Code doesn't trim the whitespace and uses the values exactly as written, so edit the configuration to remove it.
>
> Some server names are reserved for Claude Code's built-in servers: `workspace`, `claude-in-chrome`, `computer-use`, `Claude Preview`, and `Claude Browser`. If your configuration defines a server with a reserved name, Claude Code skips it at load time and shows a warning asking you to rename it. `claude mcp add` rejects a reserved name with an error.
>
> `Claude Preview` and `Claude Browser` both name the built-in server that the [Claude Code desktop app's preview pane](/docs/en/desktop#preview-your-app) uses. Before v2.1.205, `Claude Browser` wasn't reserved, so a user-configured server could register under that name.
>
> #### Tool availability
>
> The `/mcp` panel shows the tool count next to each connected server and flags servers that advertise the tools capability but expose no tools.
>
> If your request needs tools from a server that is still connecting in the background, Claude waits for that server before continuing. How the wait happens depends on your configuration:
>
> * **With [tool search](#scale-with-mcp-tool-search), the default**: the wait happens inside the `ToolSearch` call.
> * **Without tool search**: Claude uses the `WaitForMcpServers` tool instead. Configurations without tool search include a custom `ANTHROPIC_BASE_URL`, `ENABLE_TOOL_SEARCH=false`, and a model earlier than the Claude 4.5 generation on Google Cloud's Agent Platform.
> * **On a Microsoft Foundry [deployment hosted on Azure](https://platform.claude.com/docs/en/build-with-claude/claude-in-microsoft-foundry#hosting-options)**: Claude starts on the tool-search path rather than with `WaitForMcpServers`, since Claude Code di

### Source: Debug — Check MCP

> Check MCP servers
>
> Run `/mcp` to see every configured server, its connection status, and whether you have approved it for the current project. A server can be defined correctly but still not provide tools for a few common reasons:
>
> * Project-scoped servers in `.mcp.json` require a one-time approval. If the prompt was dismissed, the server stays disabled until you approve it from `/mcp`.
> * A server that fails to start shows as failed in `/mcp`. Relative file paths in `command` or `args` are a frequent cause, since they resolve against the directory you launched Claude Code from rather than the location of `.mcp.json`.
> * A server that shows as connected but lists zero tools has started successfully but isn't returning a tool list. Select **Reconnect** from `/mcp`. If the count stays at zero, run `claude --debug=mcp` and read the server's stderr in the debug log at `~/.claude/debug/<session-id>.txt`.
>
> For configuration locations and scope rules, see [MCP](/docs/en/mcp).
