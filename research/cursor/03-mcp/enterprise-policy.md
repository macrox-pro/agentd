---
primary_sources:
  - id: T1-ENT-MCP
    title: "Enterprise MCP"
    url: "https://cursor.com/docs/enterprise/model-and-integration-management.md"
    section: "Enterprise MCP"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# MCP enterprise policy

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Enterprise — Model and integration management

> # Model and Integration Management
>
> Your team can access multiple AI models and integrate Cursor with various services. This documentation covers how to control which models are available, manage MCP server trust, and set up integrations with tools like Slack, GitHub, and Linear.
>
> ## Model access control
>
> Enterprise teams can control which AI models team members can use, [contact sales](https://cursor.com/contact-sales?source=docs-model-controls) to get access. This helps manage costs, ensure appropriate usage, and comply with organizational policies.
>
> Configure model access in two places:
>
> 1. **Team Settings → Models** in the [team dashboard](https://cursor.com/docs/account/teams/dashboard.md) (Enterprise only). From Team Settings, open the **Model Providers** section to manage providers, models, defaults, and personal API key (BYOK) controls. This is the team baseline.
> 2. **Organization → Groups → \[group] → Models**, when you use [Organization Groups](https://cursor.com/docs/enterprise/organization-groups.md#model-access). Use this to widen access for specific cohorts.
>
> You can also manage the team baseline programmatically with the [Admin API model access](https://cursor.com/docs/account/teams/admin-api.md#model-access) routes, or across linked teams with the [Organization API](https://cursor.com/docs/account/organizations/organization-admin-api.md#model-access).
>
> ### How team and group model access combine
>
> Cursor reconciles team and Organization Group model settings with a **most-permissive (union)** model. Neither layer fully overrides the other:
>
> - A model is allowed if the **team** or **any of the user's Organization Groups** allows it.
> - A group setting cannot make a model more restrictive than what another allowing source already grants. Groups are for widening access, not tightening it below an allow from the team or another group.
> - Put your strictest defaults on the **team**. Use Organization Groups only to grant additional models to selected cohorts.
>
> Personal API key (BYOK) controls remain on **Team Settings → Models** only. Organization Group Models settings do not configure BYOK.
>
> See [Organization Groups](https://cursor.com/docs/enterprise/organization-groups.md#how-group-and-team-settings-combine) and [How limits and permissions combine](https://cursor.com/docs/enterprise/organizations.md#how-limits-and-permissions-combine) for the full merge rules across settings.
>
> ### How enterprise model rollout works
>
> When new models become available, Cursor doesn't immediately enable them for all enterprise teams.
>
> Instead, Enterprise teams can opt in to new models for their organization.
>
> See [Models](https://cursor.com/docs/models-and-pricing.md) for the current list of available models.
>
> ### Auto-review and model access
>
> [Auto-review](https://cursor.com/docs/agent/security/run-modes.md#run-mode) uses a background classifier that runs on [Claude 4.5 Haiku](https://cursor.com/docs/models/claude-4-5-haiku.md) or [GPT-5.4 Mini](https://cursor.com/docs/models/gpt-5-4-mini.md). Blocking all of them disables Auto-review in the IDE, even when team Run Modes includes it. See [Auto-review classifier requirements](https://cursor.com/docs/agent/security/run-modes.md#auto-review-model-requirements).
>
> ## Restrict personal API keys (BYOK controls)
>
> Enterprise teams can prevent team members from using their own API keys with third-party providers (OpenAI, Anthropic, Azure, AWS Bedrock) in Cursor. All usage goes through Cursor's included models and usage pool.
>
> Configure this in the [team dashboard](https://cursor.com/docs/account/teams/dashboard.md) under **Team Settings → Models** (Enterprise only).
>
> ## MCP server trust management
>
> The Model Context Protocol (MCP) lets you connect external tools and data sources to Cursor. MCP servers can:
>
> - Read files from external systems
> - Execute operations on your behalf
> - Access databases and APIs
> - Integrate with third-party services
>
> MCP servers are designed and implemented by external vendors, not Cursor. We work with partners to provide a [vetted marketplace](/marketplace) of trusted servers, but you should review each server's capabilities and permissions before enabling it for your team.
>
> Because MCP servers have significant capabilities, you need to manage which servers your team can use.
>
> ### MCP Allowlist
>
> Enterprise teams can control which MCP servers team members are allowed to use. Configure this in the [team dashboard](https://cursor.com/docs/account/teams/dashboard.md) under "MCP Configuration" (Enterprise only).
>
> Add each approved server as a command or URL entry, then configure its tool controls and network policy. Approving a trusted set of servers and domains is usually enough; apply stricter tool and network controls per server when you need them.
>
> You can also distribute `~/.cursor/permissions.json` through MDM to set the per-user MCP auto-run allowlist from a managed file.
>
> In that file, `mcpAllowlist` must be a JSON array of strings using `server:tool` syntax:
>
> | Entry         | Meaning                                      |
> | :------------ | :------------------------------------------- |
> | `server:tool` | One specific tool on one specific MCP server |
> | `server:*`    | All tools from one MCP server                |
> | `*:tool`      | One tool name from any MCP server            |
> | `*:*`         | All MCP tools                                |
>
> Cursor resolves the effective MCP allowlist in this order:
>
> 1. Team dashboard or other admin-controlled settings
> 2. `~/.cursor/permissions.json`
> 3. The MCP allowlist in editor settings and inline **Add to allowlist**
>
> Higher-priority sources replace lower-priority ones. They do not merge.
>
> When an allowlist is active, only servers matching an allowlist entry can run. Servers that don't match are blocked.
>
> Adding a server to the allowlist does not push it to users' machines. Team members still need to configure the server in their own [Cursor settings](https://cursor.com/docs/mcp.md).
>
> To distribute an approved server, add it to a [team marketplace](https://cursor.com/docs/plugins.md#team-marketplaces). Admins can link existing standalone Team MCP servers to the Default marketplace so teammates can install and configure them in the Agent Window, IDE, and CLI.
>
> All allowlist entries support wildcards using `*` to match any sequence of characters.
>
> #### Command-based servers (stdio)
>
> For local MCP servers configured with `command` and `args`, the allowlist matches against the **full command string**: the `command` value and all `args` values joined with spaces.
>
> Given this `mcp.json` config:
>
> ```json
> {
>   "mcpServers": {
>     "my-tool": {
>       "command": "npx",
>       "args": ["-y", "@acme/mcp-tool@latest"]
>     }
>   }
> }
> ```
>
> The full command string is `npx -y @acme/mcp-tool@latest`. On most systems, the shell resolves `npx` to a full path like `/usr/local/bin/npx` or `/opt/homebrew/bin/npx`, so the actual string becomes `/usr/local/bin/npx -y @acme/mcp-tool@latest`.
>
> Use a leading `*` wildcard to match regardless of the install path:
>
> | Allowlist entry                               | Matches                                                           |
> | :-------------------------------------------- | :---------------------------------------------------------------- |
> | `*npx -y @acme/mcp-tool@latest`               | `npx` at any path, with these exact arguments                     |
> | `/usr/local/bin/npx -y @acme/mcp-tool@latest` | Only this exact path                                              |
> | `*npx -y @acme/*`                             | Any `@acme`-scoped MCP package                                    |
> | `*python */scripts/mcp-server.py*`            | A Python server at any matching path, with any trailing arguments |
>
> #### URL-based servers (HTTP/SSE)
>
> For remote MCP servers configured with `url`, the allowlist matches against the URL.
>
> Given this `mcp.json` config:
>
> ```json
> {
>   "mcpServers": {
>     "acme-tools": {
>       "url": "https://mcp.acme.com/sse"
>     }
>   }
> }
> ```
>
> The allowlist entry matches against the full URL `https://mcp.acme.com/sse`:
>
> | Allowlist entry            | Matches                                 |
> | :------------------------- | :-------------------------------------- |
> | `https://mcp.acme.com/sse` | This exact URL                          |
> | `https://*.acme.com/*`     | Any subdomain and path under `acme.com` |
> | `https://mcp.acme.com/*`   | Any path on this host                   |
>
> ### Per-server tool controls
>
> Tool controls live in the MCP Configuration section and are set per server, not in a separate auto-run list. For each approved server, restrict which tools can run by listing them in that server's Tools field. Leave the field empty to allow all tools from that server.
>
> ### Per-server network controls
>
> Each approved server has its own network policy, so you control what it can reach.
>
> Remote (URL) MCP servers are restricted to the configured URL entry pattern.
>
> Local command-based (`stdio`) servers run in a sandbox with one of these network modes:
>
> | Network mode   | Behavior                                                |
> | :------------- | :------------------------------------------------------ |
> | **Allow all**  | No egress restrictions.                                 |
> | **Allowlist**  | Only listed destinations are reachable.                 |
> | **Deny all**   | Run the server locally with no outbound network access. |
> | **No sandbox** | Run without command or network sandboxing.              |
>
> ## Git repository blocklist
>
> You can prevent Cursor from accessing specific repositories.
>
> Add repository URLs or patterns in the [team dashboard](https://cursor.com/docs/account/teams/dashboard.md) under "Repository Blocklist" (Enterprise only). Cursor will refuse to index or work with blocked repositories.
>
> ## Protected Git Scopes
>
> Lock a Git organization, group, or namespace to your Cursor organization so only your teams can use its repositories with [Cloud Agents](https://cursor.com/docs/cloud-agent.md), [automations](https://cursor.com/docs/cloud-agent/automations.md), and [Bugbot](https://cursor.com/docs/bugbot.md). Cursor always verifies that a user can access a repository's connected source before it runs an agent or Bugbot check. Protected Git Scopes adds an organization-level guarantee on top of that per-user check, so enterprises can be confident their code can't be reached through unsanctioned ("shadow IT") Cursor accounts or outside teams, even ones that already have legitimate Git access.
>
> Protect or remove a scope from the [Integrations & MCP](https://cursor.com/dashboard/integrations) tab of your dashboard (Teams and Enterprise). Claiming a scope requires a Cursor team admin who is also a Git provider admin. Works with cloud and self-hosted GitHub and GitLab.
>
> ## Integration: Slack
>
> The Slack integration enables Cloud Agents to run directly from Slack. Team members can mention `@cursor` with a prompt and get automated code changes delivered as pull requests.
>
> Cursor requires permissions to read messages, post responses, and access channel metadata. See the [Slack integration documentation](https://cursor.com/docs/integrations/slack.md#permissions) for the full list.
>
> See [Slack integration](https://cursor.com/docs/integrations/slack.md) for detailed setup and usage instructions.
>
> ## Integration: GitHub, GHES, and GitLab
>
> Connect Cursor to your version control system to work with Cloud Agents.
>
> Cursor requires read access to repositories and write access to create PRs. You control which repositories the Cursor app can access.
>
> See [GitHub integration](https://cursor.com/docs/integrations/github.md) for setup.
>
> ## Integration: Linear
>
> Connect Linear to start Cloud Agents from issues.
>
> Cursor requires read access to issues and write access to update issue status.
>
> See [Linear integration](https://cursor.com/docs/integrations/linear.md) for details.
>
> ### Model controls are available on the Enterprise plan
>
> Contact our team to learn about model restrictions and MCP management.
>
>

### Source: Permissions reference

> # permissions.json reference
>
> Use `permissions.json` to configure MCP tool and terminal command allowlists and to steer the [Auto-review mode](https://cursor.com/docs/agent/security/run-modes.md#run-mode) classifier so tools run without approval.
>
> When `permissions.json` defines an allowlist, it **overrides** the corresponding in-app allowlist in Cursor Settings. The in-app allowlist editor becomes read-only for that allowlist type.
>
> ## File location
>
> Cursor reads `permissions.json` from two locations:
>
> ```text
> ~/.cursor/permissions.json              # per-user (applies everywhere)
> <workspace>/.cursor/permissions.json    # per-repo (applies in this workspace)
> ```
>
> Both files are optional. When both exist, Cursor **concatenates** the arrays inside every field. Per-user and per-repo entries combine; one does not replace the other. Commit the per-repo file so teammates inherit the same rules.
>
> The files are read on startup and re-read automatically whenever they change. JSONC (JSON with comments) is supported.
>
> ## Top-level fields
>
> All fields are optional. Unknown keys are ignored.
>
> | Field               | Type       | Default | Description                                                                                                                                                               |
> | :------------------ | :--------- | :------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `mcpAllowlist`      | `string[]` | not set | MCP tools that can run without approval. When set, overrides the in-app MCP allowlist.                                                                                    |
> | `terminalAllowlist` | `string[]` | not set | Terminal commands that can run without approval. When set, overrides the in-app terminal allowlist.                                                                       |
> | `autoRun`           | `object`   | not set | Natural-language guidance for the **Auto-review** mode classifier. See [`autoRun` configuration](https://cursor.com/docs/reference/permissions.md#autorun-configuration). |
>
> Non-string entries inside either array are silently dropped.
>
> ## Precedence
>
> Allowlists come from three sources, evaluated in strict priority order:
>
> ```text
> team admin (dashboard)  >  permissions.json (per-user ∪ per-repo)  >  IDE settings UI
>        (highest)                                                          (lowest)
> ```
>
> - **Team admin controls.** If your team admin has configured Run Mode controls through the dashboard, those settings take effect. Neither `permissions.json` nor the IDE allowlist can add extra entries.
> - **permissions.json.** When Run Mode is not admin-controlled and `permissions.json` defines a key, that key's value **replaces** the corresponding IDE allowlist entirely. Arrays from `~/.cursor/permissions.json` and `<workspace>/.cursor/permissions.json` are concatenated before being applied. The in-app editor for that allowlist becomes read-only and the "Add to allowlist" button is hidden.
> - **IDE settings.** When Run Mode is not admin-controlled and neither permissions file defines a given key, the IDE allowlist from Cursor Settings is used.
>
> MCP, terminal, and `autoRun` are independent. You can define one in `permissions.json` and manage the others in the IDE. Defining only `mcpAllowlist` in the file overrides the MCP allowlist but leaves the terminal allowlist under IDE control.
>
> If neither file is present, both are unparseable, or no file contains a given key, Cursor falls back to the IDE allowlist for that key. If a key is present in either file but evaluates to an empty array (after concatenation), the effective allowlist for that type is empty. Cursor does not fall back to the IDE allowlist in that case.
>
> ## How it appears in Cursor Settings
>
> When `permissions.json` defines an allowlist, Cursor Settings notes that the allowlist is configured via `permissions.json`.
>
> - If the allowlist is controlled by `permissions.json`, the editor becomes read-only and shows the file-defined entries. The "Add to allowlist" option is not available for that allowlist type.
> - If the allowlist is admin-controlled, the editor becomes read-only and shows the admin-defined entries.
>
> ## MCP allowlist format
>
> Each entry is a `server:tool` string. Both parts are matched case-insensitively. The `*` wildcard matches any value for that part.
>
> | Pattern             | Matches                                                      |
> | :------------------ | :----------------------------------------------------------- |
> | `my-server:my_tool` | Exactly the tool `my_tool` from the server named `my-server` |
> | `my-server:*`       | All tools from `my-server`                                   |
> | `*:my_tool`         | The tool `my_tool` from any server                           |
> | `*:*`               | All tools from all servers                                   |
>
> The server name is the key you used in `mcp.json` (e.g. `"github"`, `"linear"`). Glob-style `*` patterns also work inside names (e.g. `my-server:list_*` matches `list_issues`, `list_users`, etc.).
>
> Entries that do not contain a `:` are ignored.
>
> ## `autoRun` configuration
>
> The `autoRun` object steers the LLM classifier that gates shell, MCP, and Fetch tool calls when [Auto-review mode](https://cursor.com/docs/agent/security/run-modes.md#run-mode) is active. It has no effect in **Allowlist** or **Run Everything**.
>
> | Field                | Type       | Description                                                                                                                     |
> | :------------------- | :--------- | :------------------------------------------------------------------------------------------------------------------------------ |
> | `allow_instructions` | `string[]` | Natural-language hints describing call shapes the classifier should lean toward allowing.                                       |
> | `block_instructions` | `string[]` | Natural-language hints describing call shapes the classifier should lean toward blocking, surfacing an approval prompt instead. |
>
> Each entry is a free-form sentence. Write the instruction the way you would tell a teammate what to watch for. Calls that match an `allow_instructions` entry still go through the safety check; calls that match a `block_instructions` entry can still be approved when Cursor insists. Treat both as steering, not enforcement.
>
> Per-user and per-repo entries are concatenated, so a workspace can layer repo-specific guardrails on top of your personal defaults.
>
> ## Terminal allowlist format
>
> Each entry is a command or command prefix string.
>
> | Pattern        | Matches                                                                                          |
> | :------------- | :----------------------------------------------------------------------------------------------- |
> | `git`          | Any command starting with `git` (e.g. `git status`, `git diff`)                                  |
> | `git status`   | Only `git status` (and anything starting with `git status `)                                     |
> | `npm:install*` | `npm install`, `npm install express`, etc. The `:` separates the base command from an args glob. |
>
> Matching is case-sensitive and uses prefix semantics: `git` matches `git status` but not `gitk`.
>
> ## Examples
>
> ### Set MCP allowlist globally
>
> ```jsonc
> {
>   // Overrides the in-app MCP allowlist entirely.
>   "mcpAllowlist": [
>     "github:*",
>     "linear:list_issues"
>   ]
> }
> ```
>
> ### Set terminal allowlist globally
>
> ```jsonc
> {
>   "terminalAllowlist": [
>     "git",
>     "npm",
>     "yarn",
>     "pnpm",
>     "cargo",
>     "make"
>   ]
> }
> ```
>
> ### Override only one allowlist type
>
> If `permissions.json` only defines `mcpAllowlist`, the MCP allowlist is taken from the file while the terminal allowlist remains under IDE control:
>
> ```jsonc
> {
>   "mcpAllowlist": [
>     "github:*",
>     "linear:*"
>   ]
> }
> ```
>
> Any MCP entries previously set in Cursor Settings are ignored while this file is present. Terminal allowlist entries in Cursor Settings still apply.
>
> ### Combined setup
>
> ```jsonc
> {
>   "mcpAllowlist": [
>     "github:*",
>     "linear:*",
>     "notion:search"
>   ],
>   "terminalAllowlist": [
>     "git",
>     "npm",
>     "cargo build",
>     "cargo test"
>   ]
> }
> ```
>
> ### Steer the Auto-review classifier
>
> ```jsonc
> {
>   "autoRun": {
>     "allow_instructions": [
>       "Read-only inspections of build artifacts under ./dist are fine."
>     ],
>     "block_instructions": [
>       "Especially for delete operations, I like for the classifier to reject so I can have a chance to review the operation."
>     ]
>   }
> }
> ```
>
> ### Combine per-user and per-repo files
>
> `~/.cursor/permissions.json`:
>
> ```jsonc
> {
>   "terminalAllowlist": ["git", "npm", "pnpm"],
>   "autoRun": {
>     "block_instructions": [
>       "Anything that touches my SSH config or shell rc files."
>     ]
>   }
> }
> ```
>
> `<workspace>/.cursor/permissions.json`:
>
> ```jsonc
> {
>   "terminalAllowlist": ["cargo build", "cargo test"],
>   "autoRun": {
>     "block_instructions": [
>       "Never run database migrations against the production schema in this repo."
>     ]
>   }
> }
> ```
>
> The effective config is the concatenation of both files:
>
> ```jsonc
> {
>   "terminalAllowlist": ["git", "npm", "pnpm", "cargo build", "cargo test"],
>   "autoRun": {
>     "block_instructions": [
>       "Anything that touches my SSH config or shell rc files.",
>       "Never run database migrations against the production schema in this repo."
>     ]
>   }
> }
> ```
>
> ## Notes
>
> - **Run Mode required.** `permissions.json` only takes effect when Run Mode is enabled in Cursor Settings (**Auto-review**, **Allowlist**, or **Run Everything**). `autoRun` instructions are only consulted in **Auto-review** mode. Before Cursor 3.5, allowlists were not consulted in the deprecated **Ask Every Time** mode.
> - **Not a security boundary.** Allowlists and `autoRun` instructions are best-effort convenience. They are not a security guarantee. See [agent security](https://cursor.com/docs/agent/security.md) for details.
> - **Override IDE, merge files.** When `permissions.json` defines a key, it fully replaces the in-app allowlist for that type. Entries from per-user and per-repo files are concatenated; IDE entries are not merged in.
> - **IDE display.** When `permissions.json` controls an allowlist, the corresponding settings section becomes read-only and shows the file-defined entries. The "Add to allowlist" option is hidden.
> - **CLI permissions are separate.** The Cursor CLI has its own permissions system. See [CLI Permissions](https://cursor.com/docs/cli/reference/permissions.md) for that reference.
>
>
