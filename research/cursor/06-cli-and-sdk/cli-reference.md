---
primary_sources:
  - id: T1-CLI-PARAMS
    title: "Reference"
    url: "https://cursor.com/docs/cli/reference/parameters.md"
    section: "Reference"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# CLI reference

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Parameters

> # Parameters
>
> ## Global options
>
> Global options can be used with any command:
>
> | Option                     | Description                                                                                                          |
> | -------------------------- | -------------------------------------------------------------------------------------------------------------------- |
> | `-v, --version`            | Output the version number                                                                                            |
> | `--api-key <key>`          | API key for authentication (can also use `CURSOR_API_KEY` env var)                                                   |
> | `-H, --header <header>`    | Add custom header to agent requests (format: `Name: Value`, can be used multiple times)                              |
> | `-p, --print`              | Print responses to console (for scripts or non-interactive use). Has access to all tools, including write and shell. |
> | `--output-format <format>` | Output format (only works with `--print`): `text`, `json`, or `stream-json` (default: `text`)                        |
> | `--stream-partial-output`  | Stream partial output as individual text deltas (only works with `--print` and `stream-json` format)                 |
> | `--resume [chatId]`        | Resume a chat session                                                                                                |
> | `--continue`               | Continue the previous session (alias for `--resume=-1`)                                                              |
> | `--model <model>`          | Model to use                                                                                                         |
> | `--mode <mode>`            | Set agent mode: `plan` or `ask` (agent is the default when no mode is specified)                                     |
> | `--plan`                   | Start in plan mode (shorthand for `--mode=plan`)                                                                     |
> | `--list-models`            | List all available models                                                                                            |
> | `-f, --force`              | Force allow commands unless explicitly denied                                                                        |
> | `--yolo`                   | Alias for `--force`                                                                                                  |
> | `--sandbox <mode>`         | Set sandbox mode: `enabled` or `disabled`                                                                            |
> | `--approve-mcps`           | Automatically approve all MCP servers                                                                                |
> | `--trust`                  | Trust the workspace without prompting (headless mode only)                                                           |
> | `--workspace <path>`       | Workspace directory to use                                                                                           |
> | `--plugin-dir <path>`      | Load a local plugin directory (can be specified multiple times)                                                      |
> | `-w, --worktree [name]`    | Run in a new Git worktree under `~/.cursor/worktrees/<reponame>/<name>`. If omitted, a name is generated.            |
> | `--worktree-base <branch>` | Branch or ref to base the new worktree on (default: current HEAD)                                                    |
> | `--skip-worktree-setup`    | Skip running worktree setup scripts from `.cursor/worktrees.json`                                                    |
> | `-h, --help`               | Display help for command                                                                                             |
>
> ## Commands
>
> | Command                       | Description                                                       | Usage                               |
> | ----------------------------- | ----------------------------------------------------------------- | ----------------------------------- |
> | `agent [prompt...]`           | Start in agent mode (the default)                                 | `agent agent "fix the tests"`       |
> | `login`                       | Authenticate with Cursor                                          | `agent login`                       |
> | `logout`                      | Sign out and clear stored authentication                          | `agent logout`                      |
> | `status` \| `whoami`          | View authentication status                                        | `agent status`                      |
> | `about`                       | Display version, system, and account info                         | `agent about`                       |
> | `models`                      | List available models for this account                            | `agent models`                      |
> | `mcp`                         | Manage MCP servers                                                | `agent mcp`                         |
> | `sandbox`                     | Configure sandbox mode or run one command in a sandbox (hidden)   | `agent sandbox enable`              |
> | `worker`                      | Start a private cloud worker that runs agents in your environment | `agent worker start`                |
> | `acp`                         | Start ACP server mode (advanced, hidden command)                  | `agent acp`                         |
> | `update`                      | Update Cursor Agent to the latest version                         | `agent update`                      |
> | `ls`                          | Resume a chat session                                             | `agent ls`                          |
> | `resume`                      | Resume the latest chat session                                    | `agent resume`                      |
> | `create-chat`                 | Create a new empty chat and return its ID                         | `agent create-chat`                 |
> | `generate-rule` \| `rule`     | Generate a new Cursor rule with interactive prompts               | `agent generate-rule`               |
> | `install-shell-integration`   | Install shell integration to `~/.zshrc`                           | `agent install-shell-integration`   |
> | `uninstall-shell-integration` | Remove shell integration from `~/.zshrc`                          | `agent uninstall-shell-integration` |
> | `help [command]`              | Display help for command                                          | `agent help [command]`              |
>
> `agent acp` is intended for custom ACP clients and advanced integrations. It is
> hidden from default command help output.
>
> When no command is specified, Cursor Agent starts in interactive agent mode by
> default.
>
> ## MCP
>
> Manage MCP servers configured for Cursor Agent.
>
> | Subcommand                | Description                                                                              | Usage                               |
> | ------------------------- | ---------------------------------------------------------------------------------------- | ----------------------------------- |
> | `login <identifier>`      | Authenticate with an MCP server configured in `.cursor/mcp.json` or `~/.cursor/mcp.json` | `agent mcp login <identifier>`      |
> | `list`                    | List configured MCP servers and their status                                             | `agent mcp list`                    |
> | `list-tools <identifier>` | List available tools and their argument names for a specific MCP                         | `agent mcp list-tools <identifier>` |
> | `enable <identifier>`     | Add an MCP server to the local approved list                                             | `agent mcp enable <identifier>`     |
> | `disable <identifier>`    | Disable an MCP server so it won't load or prompt for approval                            | `agent mcp disable <identifier>`    |
>
> All MCP commands support `-h, --help` for command-specific help.
>
> ## Sandbox
>
> Configure sandbox mode or run one command in a sandbox.
>
> | Subcommand            | Description                                          | Usage                   |
> | --------------------- | ---------------------------------------------------- | ----------------------- |
> | `enable`              | Enable sandbox mode for command execution            | `agent sandbox enable`  |
> | `disable`             | Disable sandbox mode and use allowlist mode          | `agent sandbox disable` |
> | `reset`               | Reset sandbox configuration to defaults              | `agent sandbox reset`   |
> | `run <cmd> [args...]` | Run a command in a sandbox with workspace read/write | `agent sandbox run ls`  |
> | `help [command]`      | Display help for command                             | `agent sandbox help`    |
>
> | Command       | Option                          | Description                                                        |
> | ------------- | ------------------------------- | ------------------------------------------------------------------ |
> | `sandbox run` | `--allow-paths <paths>`         | Comma-separated list of extra read/write paths                     |
> | `sandbox run` | `--readonly-paths <paths>`      | Comma-separated list of extra read-only paths                      |
> | `sandbox run` | `--blocked-patterns <patterns>` | Comma-separated list of gitignore-style patterns to block          |
> | `sandbox run` | `--sandbox`                     | Run with the workspace read/write sandbox policy (default: `true`) |
> | `sandbox run` | `--network`                     | Enable network access in the sandbox (default: `false`)            |
> | `sandbox run` | `--sb-debug`                    | Write sandbox debug logs to a temp folder and print the path       |
>
> All sandbox commands support `-h, --help` for command-specific help.
>
> ## Worker
>
> Start a private cloud worker that connects to Cursor and runs agents in your environment.
>
> | Subcommand       | Description                                                             | Usage                |
> | ---------------- | ----------------------------------------------------------------------- | -------------------- |
> | `start`          | Start the worker and connect to Cursor                                  | `agent worker start` |
> | `debug`          | Run private worker preflight diagnostics for auth, privacy, and routing | `agent worker debug` |
> | `help [command]` | Display help for command                                                | `agent worker help`  |
>
> | Command        | Option                             | Description                                                                                         |
> | -------------- | ---------------------------------- | --------------------------------------------------------------------------------------------------- |
> | `worker`       | `--auth-token-file <path>`         | Path to a file containing the worker auth token                                                     |
> | `worker`       | `--worker-dir <path>`              | Workspace root to expose to agents. Repeatable. The first value is the assignment identity.         |
> | `worker`       | `--management-addr <address>`      | Listen address for `/healthz`, `/readyz`, and `/metrics`                                            |
> | `worker`       | `--label <key=value>`              | Add a worker label. Can be used multiple times. Can't be used with `--labels-file`.                 |
> | `worker`       | `--labels-file <path>`             | Path to a JSON or TOML labels file. Can also use `CURSOR_WORKER_LABELS_FILE`.                       |
> | `worker`       | `--idle-release-timeout <seconds>` | Seconds the worker may stay connected after becoming idle. Default `0` disables idle-based release. |
> | `worker`       | `--pool`                           | Register for pool assignment. One cloud agent claims the worker at a time.                          |
> | `worker`       | `--single-use`                     | Legacy alias for `--pool`                                                                           |
> | `worker`       | `--pool-name <name>`               | Pool label for pool workers. Requires `--pool` or `--single-use`. Defaults to `default`.            |
> | `worker`       | `--name <name>`                    | Custom display name. Defaults to the machine hostname.                                              |
> | `worker`       | `--data-dir <path>`                | Base directory for logs, artifacts, and recording data                                              |
> | `worker`       | `--debug`                          | Print worker debug diagnostics before starting bridge mode                                          |
> | `worker start` | `--verbose`                        | Enable verbose startup logs                                                                         |
> | `worker debug` | `--json`                           | Output the debug report as JSON                                                                     |
>
> ## Command-specific options
>
> | Command            | Option              | Description                                       |
> | ------------------ | ------------------- | ------------------------------------------------- |
> | `status`, `whoami` | `--format <format>` | Output format: `text` or `json` (default: `text`) |
> | `about`            | `--format <format>` | Output format: `text` or `json` (default: `text`) |
>
> ## Arguments
>
> When starting in chat mode (default behavior), you can provide an initial prompt:
>
> **Arguments:**
>
> - `prompt` — Initial prompt for the agent
>
> ## Getting help
>
> All commands support the global `-h, --help` option to display command-specific help.
>
>

### Source: Authentication

> # Authentication
>
> Cursor CLI supports two authentication methods: browser-based login (recommended) and API keys.
>
> ## Browser authentication (recommended)
>
> Use the browser flow for the easiest authentication experience:
>
> ```bash
> # Log in using browser flow
> agent login
>
> # Check authentication status
> agent status
>
> # Log out and clear stored authentication
> agent logout
> ```
>
> The login command opens your default browser and prompts you to authenticate with your Cursor account. Set `NO_OPEN_BROWSER=1` to print the login URL without opening a browser. Once complete, your credentials are securely stored locally.
>
> ## API key authentication
>
> For automation, scripts, or CI environments, use API key authentication:
>
> ### Step 1: Generate an API key
>
> Generate a user API key from [Cursor Dashboard → API Keys](https://cursor.com/dashboard/api).
>
> ### Step 2: Set the API key
>
> You can provide the API key in two ways:
>
> **Option 1: Environment variable (recommended)**
>
> ```bash
> export CURSOR_API_KEY=your_api_key_here
> agent "implement user authentication"
> ```
>
> **Option 2: Command line flag**
>
> ```bash
> agent --api-key your_api_key_here "implement user authentication"
> ```
>
> ## Authentication status
>
> Check your current authentication status:
>
> ```bash
> agent status
> ```
>
> This command will display:
>
> - Whether you're authenticated
> - Your account information
> - Current endpoint configuration
>
> ## Troubleshooting
>
> - **"Not authenticated" errors:** Run `agent login` or ensure your API key is correctly set
> - **Browser doesn't open:** Run `NO_OPEN_BROWSER=1 agent login` and open the printed URL manually
>
>

### Source: Slash commands

> # Slash commands
>
> | Command                                | Description                                                                                                            |
> | -------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
> | `/model [filter]`                      | Select a model. Press `Tab` to edit.                                                                                   |
> | `/run-everything [on\|off\|status]`    | Toggle Run Everything or show its status. `/auto-run` is an alias.                                                     |
> | `/plan [prompt]`                       | Switch to Plan mode, show the current plan, or submit a prompt in Plan mode                                            |
> | `/ask`                                 | Toggle Ask mode for read-only questions                                                                                |
> | `/debug [prompt]`                      | Toggle Debug mode or submit a prompt in Debug mode                                                                     |
> | `/goal [objective]`                    | Give the agent a long-lived objective to work towards until it's fully complete. Rolling out.                          |
> | `/logs`                                | Show the debug log path and copy it to the clipboard                                                                   |
> | `/update`                              | Update Cursor Agent to the latest version                                                                              |
> | `/max-mode`                            | Toggle Max Mode on legacy request-based plans                                                                          |
> | `/rename <name>`                       | Rename the current chat session                                                                                        |
> | `/clear`                               | Start a new chat session. `/new`, `/new-chat`, and `/newchat` are aliases.                                             |
> | `/resume`                              | Open recent chats and resume one                                                                                       |
> | `/fork`                                | Fork the current chat into a new session                                                                               |
> | `/summarize`                           | Summarize the conversation to reduce context. `/compress` is an alias.                                                 |
> | `/rewind`                              | Jump back to a previous message                                                                                        |
> | `/vim`                                 | Toggle Vim keys                                                                                                        |
> | `/line-numbers`                        | Toggle line numbers in code blocks                                                                                     |
> | `/show-thinking`                       | Toggle thinking block display                                                                                          |
> | `/status-indicators`                   | Toggle terminal title status indicators                                                                                |
> | `/shell [command]`                     | Enter Shell Mode. `/sh` and `/run` are aliases.                                                                        |
> | `/about`                               | Show CLI version, system, and account info. Also copies it to the clipboard.                                           |
> | `/setup-terminal`                      | Configure terminal newline keybindings. See [Terminal setup](https://cursor.com/docs/cli/reference/terminal-setup.md). |
> | `/help [command]`                      | Show help. Use `/help <command>` for command details.                                                                  |
> | `/feedback <message>`                  | Share feedback with the team                                                                                           |
> | `/open`                                | Open the repository's Git root in Cursor. `/cursor` is an alias.                                                       |
> | `/copy-request-id`                     | Copy the last request ID to the clipboard                                                                              |
> | `/copy-conversation-id`                | Copy the current conversation ID to the clipboard                                                                      |
> | `/logout`                              | Sign out from Cursor                                                                                                   |
> | `/quit`                                | Exit                                                                                                                   |
> | `/exit`                                | Exit                                                                                                                   |
> | `/mcp [list\|list-tools] [identifier]` | Manage MCP servers and list tools for a server                                                                         |
> | `/plugin [subcommand]`                 | Manage plugins and marketplaces                                                                                        |
> | `/config`                              | Configure CLI settings interactively                                                                                   |
> | `/copy`                                | Copy a previous user message to the clipboard                                                                          |
> | `/sandbox`                             | Configure sandbox mode and network access settings                                                                     |
> | `/bedrock [subcommand]`                | Configure Bedrock when the Bedrock feature is enabled                                                                  |
>
>

### Source: Permissions

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

### Source: Configuration

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

### Source: Output format

> # Output format
>
> The Cursor Agent CLI provides multiple output formats with the `--output-format` option when combined with `--print`. These formats include structured formats for programmatic use (`json`, `stream-json`) and a simplified text format for human-readable output (`text`).
>
> The default `--output-format` is `text`. This option is only valid when
> printing (`--print`) or when print mode is inferred (non-TTY stdout or piped
> stdin).
>
> ## JSON format
>
> The `json` output format emits a single JSON object (followed by a newline) when the run completes successfully. Deltas and tool events are not emitted; text is aggregated into the final result.
>
> On failure, the process exits with a non-zero code and writes an error message to stderr. No well-formed JSON object is emitted in failure cases.
>
> ### Success response
>
> When successful, the CLI outputs a JSON object with the following structure:
>
> ```json
> {
>   "type": "result",
>   "subtype": "success",
>   "is_error": false,
>   "duration_ms": 1234,
>   "duration_api_ms": 1234,
>   "result": "<full assistant text>",
>   "session_id": "<uuid>",
>   "request_id": "<optional request id>"
> }
> ```
>
> | Field             | Description                                                         |
> | ----------------- | ------------------------------------------------------------------- |
> | `type`            | Always `"result"` for terminal results                              |
> | `subtype`         | Always `"success"` for successful completions                       |
> | `is_error`        | Always `false` for successful responses                             |
> | `duration_ms`     | Total execution time in milliseconds                                |
> | `duration_api_ms` | API request time in milliseconds (currently equal to `duration_ms`) |
> | `result`          | Complete assistant response text (concatenation of all text deltas) |
> | `session_id`      | Unique session identifier                                           |
> | `request_id`      | Optional request identifier (may be omitted)                        |
>
> ## Stream JSON format
>
> The `stream-json` output format emits newline-delimited JSON (NDJSON). Each line contains a single JSON object representing an event during execution. This format aggregates text deltas and outputs **one line per assistant message** (the complete message between tool calls).
>
> The stream ends with a terminal `result` event on success. On failure, the process exits with a non-zero code and the stream may end early without a terminal event; an error message is written to stderr.
>
> **Streaming partial output:** For real-time character-level streaming, use `--stream-partial-output` with `--output-format stream-json`. This emits text as it's generated in small chunks, with multiple `assistant` events per message.
>
> With `--stream-partial-output`, the CLI emits three kinds of `assistant` events. Only the first kind contains new text:
>
> | `timestamp_ms` | `model_call_id` | What it is                                    | Action                                    |
> | -------------- | --------------- | --------------------------------------------- | ----------------------------------------- |
> | Present        | Absent          | Streaming delta with new text                 | **Use** — append `message.content[].text` |
> | Present        | Present         | Buffered flush before a tool call (duplicate) | **Skip**                                  |
> | Absent         | Absent          | Final flush at end of turn (duplicate)        | **Skip**                                  |
>
> If you don't need real-time streaming and only want the finished answer, skip all `assistant` events and read the `result` field from the terminal `result` event.
>
> ### Event types
>
> #### System initialization
>
> Emitted once at the beginning of each session:
>
> ```json
> {
>   "type": "system",
>   "subtype": "init",
>   "apiKeySource": "env|flag|login",
>   "cwd": "/absolute/path",
>   "session_id": "<uuid>",
>   "model": "<model display name>",
>   "permissionMode": "default"
> }
> ```
>
> Future fields like `tools` and `mcp_servers` may be added to this event.
>
> #### User message
>
> Contains the user's input prompt:
>
> ```json
> {
>   "type": "user",
>   "message": {
>     "role": "user",
>     "content": [{ "type": "text", "text": "<prompt>" }]
>   },
>   "session_id": "<uuid>"
> }
> ```
>
> #### Assistant message
>
> Emitted once per complete assistant message (between tool calls). Each event contains the full text of that message segment:
>
> ```json
> {
>   "type": "assistant",
>   "message": {
>     "role": "assistant",
>     "content": [{ "type": "text", "text": "<complete message text>" }]
>   },
>   "session_id": "<uuid>"
> }
> ```
>
> When `--stream-partial-output` is enabled, assistant events may include two additional fields:
>
> | Field           | Description                                                                                                  |
> | --------------- | ------------------------------------------------------------------------------------------------------------ |
> | `timestamp_ms`  | Present on streaming deltas and pre-tool-call flushes. Absent on the final flush at the end of a turn.       |
> | `model_call_id` | Present only on the buffered flush emitted before a tool call. Use this to identify and skip duplicate text. |
>
> See the [streaming partial output note](https://cursor.com/docs/cli/reference/output-format.md#stream-json-format) above for how to filter these events.
>
> #### Tool call events
>
> Tool calls are tracked with start and completion events:
>
> **Tool call started:**
>
> ```json
> {
>   "type": "tool_call",
>   "subtype": "started",
>   "call_id": "<string id>",
>   "tool_call": {
>     "readToolCall": {
>       "args": { "path": "file.txt" }
>     }
>   },
>   "session_id": "<uuid>"
> }
> ```
>
> **Tool call completed:**
>
> ```json
> {
>   "type": "tool_call",
>   "subtype": "completed",
>   "call_id": "<string id>",
>   "tool_call": {
>     "readToolCall": {
>       "args": { "path": "file.txt" },
>       "result": {
>         "success": {
>           "content": "file contents...",
>           "isEmpty": false,
>           "exceededLimit": false,
>           "totalLines": 54,
>           "totalChars": 1254
>         }
>       }
>     }
>   },
>   "session_id": "<uuid>"
> }
> ```
>
> #### Tool call types
>
> **Read file tool:**
>
> - **Started**: `tool_call.readToolCall.args` contains `{ "path": "file.txt" }`
> - **Completed**: `tool_call.readToolCall.result.success` contains file metadata and content
>
> **Write file tool:**
>
> - **Started**: `tool_call.writeToolCall.args` contains `{ "path": "file.txt", "fileText": "content...", "toolCallId": "id" }`
> - **Completed**: `tool_call.writeToolCall.result.success` contains `{ "path": "/absolute/path", "linesCreated": 19, "fileSize": 942 }`
>
> **Other tools:**
>
> - May use `tool_call.function` structure with `{ "name": "tool_name", "arguments": "..." }`
>
> #### Terminal result
>
> The final event emitted on successful completion:
>
> ```json
> {
>   "type": "result",
>   "subtype": "success",
>   "duration_ms": 1234,
>   "duration_api_ms": 1234,
>   "is_error": false,
>   "result": "<full assistant text>",
>   "session_id": "<uuid>",
>   "request_id": "<optional request id>"
> }
> ```
>
> ### Example sequence
>
> Here's a representative NDJSON sequence showing the typical flow of events:
>
> ```json
> {"type":"system","subtype":"init","apiKeySource":"login","cwd":"/Users/user/project","session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff","model":"Claude 4 Sonnet","permissionMode":"default"}
> {"type":"user","message":{"role":"user","content":[{"type":"text","text":"Read README.md and create a summary"}]},"session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff"}
> {"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"I'll read the README.md file"}]},"session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff"}
> {"type":"tool_call","subtype":"started","call_id":"toolu_vrtx_01NnjaR886UcE8whekg2MGJd","tool_call":{"readToolCall":{"args":{"path":"README.md"}}},"session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff"}
> {"type":"tool_call","subtype":"completed","call_id":"toolu_vrtx_01NnjaR886UcE8whekg2MGJd","tool_call":{"readToolCall":{"args":{"path":"README.md"},"result":{"success":{"content":"# Project\n\nThis is a sample project...","isEmpty":false,"exceededLimit":false,"totalLines":54,"totalChars":1254}}}},"session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff"}
> {"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Based on the README, I'll create a summary"}]},"session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff"}
> {"type":"tool_call","subtype":"started","call_id":"toolu_vrtx_01Q3VHVnWFSKygaRPT7WDxrv","tool_call":{"writeToolCall":{"args":{"path":"summary.txt","fileText":"# README Summary\n\nThis project contains...","toolCallId":"toolu_vrtx_01Q3VHVnWFSKygaRPT7WDxrv"}}},"session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff"}
> {"type":"tool_call","subtype":"completed","call_id":"toolu_vrtx_01Q3VHVnWFSKygaRPT7WDxrv","tool_call":{"writeToolCall":{"args":{"path":"summary.txt","fileText":"# README Summary\n\nThis project contains...","toolCallId":"toolu_vrtx_01Q3VHVnWFSKygaRPT7WDxrv"},"result":{"success":{"path":"/Users/user/project/summary.txt","linesCreated":19,"fileSize":942}}}},"session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff"}
> {"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Done! I've created the summary in summary.txt"}]},"session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff"}
> {"type":"result","subtype":"success","duration_ms":5234,"duration_api_ms":5234,"is_error":false,"result":"I'll read the README.md fileBased on the README, I'll create a summaryDone! I've created the summary in summary.txt","session_id":"c6b62c6f-7ead-4fd6-9922-e952131177ff","request_id":"10e11780-df2f-45dc-a1ff-4540af32e9c0"}
> ```
>
> ## Text format
>
> The `text` output format provides only the final assistant message without any intermediate progress updates or tool call summaries. This is the cleanest output format for scripts that only need the agent's final response.
>
> This format is ideal when you want just the answer or final message from the agent, without any progress indicators or tool execution details.
>
> ### Example output
>
> ```
> The command to move this branch onto main is `git rebase --onto main HEAD~3`.
> ```
>
> Only the final assistant message (after the last tool call) is output, with no tool call summaries or intermediate text.
>
> ## Notes
>
> - Each event is emitted as a single line terminated by `\n`
> - `thinking` events are suppressed in print mode and will not appear in any output format
> - Field additions may occur over time in a backward-compatible way (consumers should ignore unknown fields)
> - The `json` format waits for completion before outputting results
> - The `stream-json` format outputs complete agent messages
> - The `--stream-partial-output` flag provides real-time text deltas for character-level streaming (only works with `stream-json` format)
> - Tool call IDs can be used to correlate start/completion events
> - Session IDs remain consistent throughout a single agent execution
>
>

### Source: Terminal setup

> # Terminal setup
>
> Configure your terminal for the best Cursor CLI experience. This guide covers keybindings for multi-line input, Vim mode, and theme synchronization.
>
> ## Quick start
>
> If Shift+Enter doesn't work for newlines in your terminal, run `/setup-terminal` for guidance on configuring alternatives:
>
> ```bash
> /setup-terminal
> ```
>
> This command detects your terminal and provides instructions for configuring Option+Enter as an alternative way to insert newlines.
>
> ## Universal options
>
> These methods work in **all terminals**, including tmux, screen, and SSH sessions:
>
> | Method | Description                                              |
> | :----- | :------------------------------------------------------- |
> | +Enter | Type a backslash, then press Enter to insert a newline   |
> | Ctrl+J | Standard control character for newline (ASCII line feed) |
>
> If you're in tmux or having trouble with other keybindings, Ctrl+J is the most reliable option.
>
> ## Terminal support
>
> ### Native Shift+Enter support
>
> These terminals support Shift+Enter for newlines out of the box:
>
> - **iTerm2** (macOS)
> - **Ghostty**
> - **Kitty**
> - **Warp**
> - **Zed** (integrated terminal)
>
> ### Requires `/setup-terminal`
>
> These terminals need `/setup-terminal` to configure Option+Enter for newlines:
>
> - **Apple Terminal** (macOS)
> - **Alacritty**
> - **VS Code** (integrated terminal)
>
> ### Terminal multiplexers
>
> **tmux** and **screen** intercept Shift+Enter before it reaches applications. Use the universal options instead:
>
> - Ctrl+J — Works reliably in all multiplexer sessions
> - +Enter — Also works universally
>
> You can configure your outer terminal (e.g., iTerm2) for Shift+Enter, but the keybinding won't pass through tmux. Use the universal options for the most consistent experience.
>
> ## Vim mode
>
> Enable Vim keybindings for navigation and editing in the CLI input area.
>
> ### Toggle with slash command
>
> ```bash
> /vim
> ```
>
> This toggles Vim mode on or off for the current session and saves the preference.
>
> ### Configure in settings
>
> Add to your `~/.cursor/cli-config.json`:
>
> ```json
> {
>   "version": 1,
>   "editor": { "vimMode": true },
>   "permissions": { "allow": [], "deny": [] }
> }
> ```
>
> ### Modes
>
> Vim mode uses modal editing:
>
> - **Normal mode** — Navigate and execute commands (default when Vim mode is enabled)
> - **Insert mode** — Type text normally
>
> Press Esc to return to normal mode from insert mode.
>
> ### Navigation
>
> | Key     | Description                                           |
> | :------ | :---------------------------------------------------- |
> | h, l    | Move left / right                                     |
> | j, k    | Move down / up                                        |
> | w, b    | Next / previous word                                  |
> | e       | End of word                                           |
> | W, B, E | Same as above, but for WORD (non-whitespace sequence) |
> | 0, $    | Start / end of line                                   |
>
> ### Editing
>
> | Key        | Description                                 |
> | :--------- | :------------------------------------------ |
> | x          | Delete character under cursor               |
> | X          | Delete character before cursor              |
> | d + motion | Delete range (e.g., `dw` deletes word)      |
> | dd         | Delete entire line                          |
> | D          | Delete to end of line                       |
> | s          | Substitute character (delete + insert mode) |
> | S, cc      | Change entire line                          |
> | C          | Change to end of line                       |
>
> ### Entering insert mode
>
> | Key | Description             |
> | :-- | :---------------------- |
> | i   | Insert at cursor        |
> | a   | Insert after cursor     |
> | I   | Insert at start of line |
> | A   | Insert at end of line   |
> | o   | Open new line below     |
> | O   | Open new line above     |
>
> ### Counts
>
> Prefix commands with a number to repeat them. For example, `3w` moves forward 3 words, `2dd` deletes 2 lines.
>
> Vim mode affects the input area only. Navigation through chat history and other UI elements uses standard keybindings.
>
> ## Terminal theme
>
> Cursor CLI automatically detects your terminal's color scheme and adapts its appearance.
>
> ### Automatic detection
>
> The CLI queries your terminal for its background color using standard escape sequences. Most modern terminals support this:
>
> - **Dark terminals** → CLI uses dark theme
> - **Light terminals** → CLI uses light theme
>
> ### Terminals with automatic detection
>
> These terminals report their color scheme correctly:
>
> - iTerm2
> - Ghostty
> - Kitty
> - Alacritty
> - Apple Terminal
> - Windows Terminal
> - VS Code integrated terminal
>
> ### Forcing a theme
>
> If automatic detection doesn't work, you can override it with an environment variable:
>
> ```bash
> # Force dark theme
> export COLORFGBG="15;0"
>
> # Force light theme
> export COLORFGBG="0;15"
> ```
>
> Add this to your shell profile (`.bashrc`, `.zshrc`, etc.) to make it permanent.
>
> ### Troubleshooting theme issues
>
> **Colors look wrong:**
>
> - Ensure your terminal supports 256 colors or true color
> - Check that `TERM` is set correctly (e.g., `xterm-256color`)
> - Try setting `COLORFGBG` explicitly
>
> **tmux users:**
>
> - Add to your `.tmux.conf` to pass through color detection:
>   ```
>   set -g default-terminal "tmux-256color"
>   set -ag terminal-overrides ",xterm-256color:RGB"
>   ```
> - Restart tmux after making changes
>
> ## Manual configuration
>
> If `/setup-terminal` doesn't work for your terminal, you can manually configure keybindings.
>
> ### Option+Enter for newlines
>
> Option+Enter sends a special escape sequence that Cursor CLI recognizes as a newline. Configure your terminal to send `\x1b\r` (Escape followed by carriage return) when Option+Enter is pressed.
>
> **iTerm2:**
>
> 1. Open **Preferences** → **Profiles** → **Keys** → **Key Mappings**
> 2. Click **+** to add a new mapping
> 3. Set **Keyboard Shortcut** to Option+Enter
> 4. Set **Action** to "Send Escape Sequence"
> 5. Enter `\r` as the escape sequence
>
> **Alacritty:**
>
> Add to your `alacritty.toml`:
>
> ```toml
> [keyboard]
> bindings = [
>   { key = "Return", mods = "Alt", chars = "\u001b\r" }
> ]
> ```
>
> **Kitty:**
>
> Add to your `kitty.conf`:
>
> ```
> map alt+enter send_text all \x1b\r
> ```
>
> ### Shift+Enter
>
> Shift+Enter support depends on your terminal correctly reporting the key modifier. Most modern terminals handle this automatically, but some may need configuration.
>
> **VS Code terminal:**
>
> VS Code's integrated terminal may not pass Shift+Enter correctly. Add to your `keybindings.json`:
>
> ```json
> {
>   "key": "shift+enter",
>   "command": "workbench.action.terminal.sendSequence",
>   "args": { "text": "\u001b[13;2u" },
>   "when": "terminalFocus"
> }
> ```
>
> ## Troubleshooting
>
> **Keybindings not working:**
>
> - Verify your terminal is detecting the keys correctly using `cat` or `showkey`
> - Check if a terminal multiplexer (tmux/screen) is intercepting the keys
> - Use Ctrl+J as a reliable fallback
>
> **tmux users:**
>
> - Shift+Enter and Option+Enter won't work through tmux
> - Use Ctrl+J or +Enter instead
> - These universal options work everywhere, including nested tmux sessions
>
> **SSH sessions:**
>
> - Remote terminal capabilities depend on your local terminal emulator
> - Ctrl+J works reliably over SSH
> - +Enter is another reliable option
>
> ## Summary
>
> | Keybinding   | Works in                          | Notes                                                      |
> | :----------- | :-------------------------------- | :--------------------------------------------------------- |
> | Ctrl+J       | All terminals                     | Most reliable, works everywhere                            |
> | +Enter       | All terminals                     | Universal alternative                                      |
> | Shift+Enter  | iTerm2, Ghostty, Kitty, Warp, Zed | Native support, no config needed                           |
> | Option+Enter | After `/setup-terminal`           | Newline alternative for Apple Terminal, Alacritty, VS Code |
>
>
