---
primary_sources:
  - id: T1-CLI-CONFIG
    title: "CLI config"
    url: "https://cursor.com/docs/cli/reference/configuration.md"
    section: "CLI config"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# CLI config and permissions

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: CLI configuration

> # Configuration
>
> Configure the Agent CLI using the `cli-config.json` file.
>
> ## File location
>
> | Type    | Platform    | Path                                       |
> | :------ | :---------- | :----------------------------------------- |
> | Global  | macOS/Linux | `~/.cursor/cli-config.json`                |
> | Global  | Windows     | `$env:USERPROFILE\.cursor\cli-config.json` |
> | Project | All         | `<project>/.cursor/cli.json`               |
>
> Only permissions can be configured at the project level. All other CLI
> settings must be set globally.
>
> Override with environment variables:
>
> - **`CURSOR_CONFIG_DIR`**: custom directory path
> - **`XDG_CONFIG_HOME`** (Linux/BSD): uses `$XDG_CONFIG_HOME/cursor/cli-config.json`
>
> ## Schema
>
> ### Required fields
>
> | Field               | Type      | Description                                                                                    |
> | :------------------ | :-------- | :--------------------------------------------------------------------------------------------- |
> | `version`           | number    | Config schema version (current: `1`)                                                           |
> | `editor.vimMode`    | boolean   | Enable Vim keybindings (default: `false`)                                                      |
> | `permissions.allow` | string\[] | Permitted operations (see [Permissions](https://cursor.com/docs/cli/reference/permissions.md)) |
> | `permissions.deny`  | string\[] | Forbidden operations (see [Permissions](https://cursor.com/docs/cli/reference/permissions.md)) |
>
> ### Optional fields
>
> | Field                                 | Type    | Description                                                             |
> | :------------------------------------ | :------ | :---------------------------------------------------------------------- |
> | `channel`                             | string  | Release channel used for CLI updates                                    |
> | `model`                               | object  | Selected model configuration                                            |
> | `maxMode`                             | boolean | Persisted preference for max mode in the model picker                   |
> | `hasChangedDefaultModel`              | boolean | CLI-managed model override flag                                         |
> | `notifications`                       | boolean | Send a terminal notification when the agent finishes or needs input     |
> | `hints`                               | boolean | Show CLI hints while the agent is working                               |
> | `rewind`                              | boolean | Enable `/rewind` to restore an earlier message in the session           |
> | `suggestNextPrompt`                   | boolean | Suggest a follow-up prompt at the end of each turn                      |
> | `display.showLineNumbers`             | boolean | Show line numbers in rendered code blocks                               |
> | `display.showThinkingBlocks`          | boolean | Render model thinking blocks when available                             |
> | `display.showStatusIndicators`        | boolean | Enable terminal title status indicators                                 |
> | `display.showStatusLineRunningTime`   | boolean | Show elapsed running time in the status line                            |
> | `approvalMode`                        | string  | Approval mode: `allowlist`, `auto-review`, or `unrestricted`            |
> | `sandbox.mode`                        | string  | Sandbox mode override                                                   |
> | `sandbox.networkAccess`               | string  | Network access setting for sandbox mode                                 |
> | `network.useHttp1ForAgent`            | boolean | Use HTTP/1.1 instead of HTTP/2 for agent connections (default: `false`) |
> | `attribution.attributeCommitsToAgent` | boolean | Add "Made with Cursor" trailer to Agent commits (default: `true`)       |
> | `attribution.attributePRsToAgent`     | boolean | Add "Made with Cursor" footer to Agent PRs (default: `true`)            |
>
> ## Examples
>
> ### Minimal config
>
> ```json
> {
>   "version": 1,
>   "editor": { "vimMode": false },
>   "permissions": { "allow": ["Shell(ls)"], "deny": [] }
> }
> ```
>
> ### Enable Vim mode
>
> ```json
> {
>   "version": 1,
>   "editor": { "vimMode": true },
>   "permissions": { "allow": ["Shell(ls)"], "deny": [] }
> }
> ```
>
> ### Configure permissions
>
> ```json
> {
>   "version": 1,
>   "editor": { "vimMode": false },
>   "permissions": {
>     "allow": ["Shell(ls)", "Shell(echo)"],
>     "deny": ["Shell(rm)"]
>   }
> }
> ```
>
> See [Permissions](https://cursor.com/docs/cli/reference/permissions.md) for available permission types and examples.
>
> ## Troubleshooting
>
> **Config errors**: Move the file aside and restart:
>
> ```bash
> mv ~/.cursor/cli-config.json ~/.cursor/cli-config.json.bad
> ```
>
> **Changes don't persist**: Ensure valid JSON and write permissions. Some fields are CLI-managed and may be overwritten.
>
> ## Notes
>
> - Pure JSON format (no comments)
> - CLI performs self-repair for missing fields
> - Corrupted files are backed up as `.bad` and recreated
> - Permission entries are exact strings (see [Permissions](https://cursor.com/docs/cli/reference/permissions.md) for details)
>
> ## Models
>
> You can select a model for the CLI using the `/model` slash command.
>
> ```bash
> /model auto
> /model gpt-5
> /model sonnet-4-thinking
> ```
>
> See the [Slash commands](https://cursor.com/docs/cli/reference/slash-commands.md) docs for other commands.
>
> ## Proxy configuration
>
> If your network routes traffic through a proxy server, configure the CLI using environment variables and the config file.
>
> ### Environment variables
>
> Set these environment variables before running the CLI:
>
> ```bash
> export HTTP_PROXY=http://your-proxy:port
> export HTTPS_PROXY=http://your-proxy:port
> export NODE_USE_ENV_PROXY=1
> ```
>
> If your proxy performs SSL inspection (man-in-the-middle), also trust your organization's CA certificate:
>
> ```bash
> export NODE_EXTRA_CA_CERTS=/path/to/corporate-ca-cert.pem
> ```
>
> ### HTTP/1.1 fallback
>
> Some enterprise proxies (like Zscaler) don't support HTTP/2 bidirectional streaming. Enable HTTP/1.1 mode in your config:
>
> ```json
> {
>   "version": 1,
>   "editor": { "vimMode": false },
>   "permissions": { "allow": [], "deny": [] },
>   "network": {
>     "useHttp1ForAgent": true
>   }
> }
> ```
>
> This switches agent connections to HTTP/1.1 with Server-Sent Events (SSE), which works with most corporate proxies.
>
> See [Network Configuration](https://cursor.com/docs/enterprise/network-configuration.md) for proxy testing commands and troubleshooting.
>
>

### Source: CLI permissions

> # Permissions
>
> Configure what the agent is allowed to do using permission tokens in your CLI configuration. Permissions are set in `~/.cursor/cli-config.json` (global) or `<project>/.cursor/cli.json` (project-specific).
>
> ## Permission types
>
> ### Shell commands
>
> **Format:** `Shell(commandBase)`
>
> Controls access to shell commands. The `commandBase` is the first token in the command line. Supports glob patterns and an optional `command:args` syntax for finer control.
>
> | Example         | Description                                        |
> | --------------- | -------------------------------------------------- |
> | `Shell(ls)`     | Allow running `ls` commands                        |
> | `Shell(git)`    | Allow any `git` subcommand                         |
> | `Shell(npm)`    | Allow npm package manager commands                 |
> | `Shell(curl:*)` | Allow `curl` with any arguments                    |
> | `Shell(rm)`     | Deny destructive file removal (commonly in `deny`) |
>
> ### File reads
>
> **Format:** `Read(pathOrGlob)`
>
> Controls read access to files and directories. Supports glob patterns.
>
> | Example             | Description                             |
> | ------------------- | --------------------------------------- |
> | `Read(src/**/*.ts)` | Allow reading TypeScript files in `src` |
> | `Read(**/*.md)`     | Allow reading markdown files anywhere   |
> | `Read(.env*)`       | Deny reading environment files          |
> | `Read(/etc/passwd)` | Deny reading system files               |
>
> ### File writes
>
> **Format:** `Write(pathOrGlob)`
>
> Controls write access to files and directories. Supports glob patterns. Print mode can use write and shell tools. Use `permissions.allow`, `permissions.deny`, and `--force` to control what runs without prompts.
>
> | Example               | Description                           |
> | --------------------- | ------------------------------------- |
> | `Write(src/**)`       | Allow writing to any file under `src` |
> | `Write(package.json)` | Allow modifying package.json          |
> | `Write(**/*.key)`     | Deny writing private key files        |
> | `Write(**/.env*)`     | Deny writing environment files        |
>
> ### Web fetch
>
> **Format:** `WebFetch(domainOrPattern)`
>
> Controls which domains the agent can fetch when using the web fetch tool (e.g., to retrieve documentation or web pages). Without an allowlist entry, each fetch prompts for approval. Add domains to `allow` to auto-approve fetches from trusted sources.
>
> | Example                     | Description                                       |
> | --------------------------- | ------------------------------------------------- |
> | `WebFetch(docs.github.com)` | Allow fetches from `docs.github.com`              |
> | `WebFetch(*.example.com)`   | Allow fetches from any subdomain of `example.com` |
> | `WebFetch(*)`               | Allow fetches from any domain (use with caution)  |
>
> **Domain pattern matching:**
>
> - `*` matches all domains
> - `*.example.com` matches subdomains (e.g., `docs.example.com`, `api.example.com`)
> - `example.com` matches that exact domain only
>
> ### MCP tools
>
> **Format:** `Mcp(server:tool)`
>
> Controls which MCP (Model Context Protocol) tools the agent can run. Use `server` (from `mcp.json`) and `tool` name, with `*` for wildcards.
>
> | Example          | Description                                 |
> | ---------------- | ------------------------------------------- |
> | `Mcp(datadog:*)` | Allow all tools from the Datadog MCP server |
> | `Mcp(*:search)`  | Allow any server's `search` tool            |
> | `Mcp(*:*)`       | Allow all MCP tools (use with caution)      |
>
> ## Configuration
>
> Add permissions to the `permissions` object in your CLI configuration file:
>
> ```json
> {
>   "permissions": {
>     "allow": [
>       "Shell(ls)",
>       "Shell(git)",
>       "Read(src/**/*.ts)",
>       "Write(package.json)",
>       "WebFetch(docs.github.com)",
>       "WebFetch(*.github.com)",
>       "Mcp(datadog:*)"
>     ],
>     "deny": [
>       "Shell(rm)",
>       "Read(.env*)",
>       "Write(**/*.key)",
>       "WebFetch(malicious-site.com)"
>     ]
>   }
> }
> ```
>
> ## Pattern matching
>
> - Glob patterns use `**`, `*`, and `?` wildcards
> - Relative paths are scoped to the current workspace
> - Absolute paths can target files outside the project
> - Deny rules take precedence over allow rules
> - Use `command:args` (e.g., `curl:*`) to match both command and arguments with globs
>
>
