---
primary_sources:
  - id: T1-CLI
    title: "CLI"
    url: "https://opencode.ai/docs/cli.md"
    section: "Environment variables"
  - id: T1-CONFIG
    title: "Config"
    url: "https://opencode.ai/docs/config.md"
    section: "Experimental"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Environment variables

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode CLI — Environment variables

> | <nobr><code>{"--password"}</code></nobr> | `-p`  | Basic auth password (defaults to `OPENCODE_SERVER_PASSWORD`)               |
> | <nobr><code>{"--username"}</code></nobr> | `-u`  | Basic auth username (defaults to `OPENCODE_SERVER_USERNAME` or `opencode`) |
>
> ---
>
> ### auth
>
> Command to manage credentials and login for providers.
>
> ```bash
> opencode auth [command]
> ```
>
> ---
>
> #### login
>
> OpenCode is powered by the provider list at [Models.dev](https://models.dev), so you can use `opencode auth login` to configure API keys for any provider you'd like to use. This is stored in `~/.local/share/opencode/auth.json`.
>
> ```bash
> opencode auth login
> ```
>
> When OpenCode starts up it loads the providers from the credentials file. And if there are any keys defined in your environments or a `.env` file in your project.
>
> ##### Flags
>
> | Flag                                     | Short | Description                                          |
> | ---------------------------------------- | ----- | ---------------------------------------------------- |
> | <nobr><code>{"--provider"}</code></nobr> | `-p`  | Provider ID or name to log in to                     |
> | <nobr><code>{"--method"}</code></nobr>   | `-m`  | Login method label to use, skipping method selection |
>
> ---
>
> #### list
>
> Lists all the authenticated providers as stored in the credentials file.
>
> ```bash
> opencode auth list
> ```
>
> Or the short version.
>
> ```bash
> opencode auth ls
> ```
>
> ---
>
> #### logout
>
> Logs you out of a provider by clearing it from the credentials file.
>
> ```bash
> opencode auth logout
> ```
>
> ---
>
> ### github
>
> Manage the GitHub agent for repository automation.
>
> ```bash
> opencode github [command]
> ```
>
> ---
>
> #### install
>
> Install the GitHub agent in your repository.
>
> ```bash
> opencode github install
> ```
>
> This sets up the necessary GitHub Actions workflow and guides you through the configuration process. [Learn more](/docs/github).
>
> ---
>
> #### run
>
> Run the GitHub agent. This is typically used in GitHub Actions.
>
> ```bash
> opencode github run
> ```
>
> ##### Flags
>
> | Flag                                  | Description                            |
> | ------------------------------------- | -------------------------------------- |
> | <nobr><code>{"--event"}</code></nobr> | GitHub mock event to run the agent for |
> | <nobr><code>{"--token"}</code></nobr> | GitHub personal access token           |
>
> ---
>
> ### mcp
>
> Manage Model Context Protocol servers.
>
> ```bash
> opencode mcp [command]
> ```
>
> ---
>
> #### add
>
> Add an MCP server to your configuration.
>
> ```bash
> opencode mcp add
> ```
>
> This command will guide you through adding either a local or remote MCP server.
>
> ---
>
> #### list
>
> List all configured MCP servers and their connection status.
>
> ```bash
> opencode mcp list
> ```
>
> Or use the short version.
>
> ```bash
> opencode mcp ls
> ```
>
> ---
>
> #### auth
>
> Authenticate with an OAuth-enabled MCP server.
>
> ```bash
> opencode mcp auth [name]
> ```
>
> If you don't provide a server name, you'll be prompted to select from available OAuth-capable servers.
>
> You can also list OAuth-capable servers and their authentication status.
>
> ```bash
> opencode mcp auth list
> ```
>
> Or use the short version.
>
> ```bash
> opencode mcp auth ls
> ```
>
> ---
>
> #### logout
>
> Remove OAuth credentials for an MCP server.
>
> ```bash
> opencode mcp logout [name]
> ```
>
> ---
>
> #### debug
>
> Debug OAuth connection issues for an MCP server.
>
> ```bash
> opencode mcp debug <name>
> ```
>
> ---
>
> ### models
>
> List all available models from configured providers.
>
> ```bash
> opencode models [provider]
> ```
>
> This command displays all models available across your configured providers in the format `provider/model`.
>
> This is useful for figuring out the exact model name to use in [your config](/docs/config/).
>
> You can optionally pass a provider ID to filter models by that provider.
>
> ```bash
> opencode models anthropic
> ```
>
> #### Flags
>
> | Flag                                    | Description                                                  |
> | --------------------------------------- | ------------------------------------------------------------ |
> | <nobr><code>{"--refresh"}</code></nobr> | Refresh the models cache from models.dev                     |
> | <nobr><code>{"--verbose"}</code></nobr> | Use more verbose model output (includes metadata like costs) |
>
> Use the `--refresh` flag to update the cached model list. This is useful when new models have been added to a provider and you want to see them in OpenCode.
>
> ```bash
> opencode models --refresh
> ```
>
> ---
>
> ### run
>
> Run opencode in non-interactive mode by passing a prompt directly.
>
> ```bash
> opencode run [message..]
> ```
>
> This is useful for scripting, automation, or when you want a quick answer without launching the full TUI. For example.
>
> ```bash "opencode run"
> opencode run Explain the use of context in Go
> ```
>
> You can also attach to a running `opencode serve` instance to avoid MCP server cold boot times on every run:
>
> ```bash
> # Start a headless server in one terminal
> opencode serve
>
> # In another terminal, run commands that attach to it
> opencode run --attach http://localhost:4096 "Explain async/await in JavaScript"
> ```
>
> #### Flags
>
> | Flag                                     | Short | Description                                                                |
> | ---------------------------------------- | ----- | -------------------------------------------------------------------------- |
> | <nobr><code>{"--command"}</code></nobr>  |       | The command to run, use message for args                                   |
> | <nobr><code>{"--continue"}</code></nobr> | `-c`  | Continue the last session                                                  |
> | <nobr><code>{"--session"}</code></nobr>  | `-s`  | Session ID to continue                                                     |
> | <nobr><code>{"--fork"}</code></nobr>     |       | Fork the session when continuing (use with `--continue` or `--session`)    |
> | <nobr><code>{"--share"}</code></nobr>    |       | Share the session                                                          |
> | <nobr><code>{"--model"}</code></nobr>    | `-m`  | Model to use in the form of provider/model                                 |
> | <nobr><code>{"--agent"}</code></nobr>    |       | Agent to use                                                               |
> | <nobr><code>{"--file"}</code></nobr>     | `-f`  | File(s) to attach to message                                               |
> | <nobr><code>{"--format"}</code></nobr>   |       | Format: default (formatted) or json (raw JSON events)                      |
> | <nobr><code>{"--title"}</code></nobr>    |       | Title for the session (uses truncated prompt if no value provided)         |
> | <nobr><code>{"--attach"}</code></nobr>   |       | Attach to a running opencode server (e.g., http://localhost:4096)          |
> | <nobr><code>{"--password"}</code></nobr> | `-p`  | Basic auth password (defaults to `OPENCODE_SERVER_PASSWORD`)               |
> | <nobr><code>{"--username"}</code></nobr> | `-u`  | Basic auth username (defaults to `OPENCODE_SERVER_USERNAME` or `opencode`) |
> | <nobr><code>{"--dir"}</code></nobr>      |       | Directory to run in, or path on the remote server when attaching           |
> | <nobr><code>{"--port"}</code></nobr>     |       | Port for the local server (defaults to random port)                        |
> | <nobr><code>{"--variant"}</code></nobr>  |       | Model variant (provider-specific reasoning effort)                         |
> | <nobr><code>{"--thinking"}</code></nobr> |       | Show thinking blocks                                                       |
> | <nobr><code>{"--auto"}</code></nobr>     |       | Auto-approve permissions that are not explicitly denied                    |
>
> ---
>
> ### serve
>
> Start a headless OpenCode server for API access. Check out the [server docs](/docs/server) for the full HTTP interface.
>
> ```bash
> opencode serve
> ```
>
> This starts an HTTP server that provides API access to opencode functionality without the TUI interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### session
>
> Manage OpenCode sessions.
>
> ```bash
> opencode session [command]
> ```
>
> ---
>
> #### list
>
> List all OpenCode sessions.
>
> ```bash
> opencode session list
> ```
>
> ##### Flags
>
> | Flag                                      | Short | Description                          |
> | ----------------------------------------- | ----- | ------------------------------------ |
> | <nobr><code>{"--max-count"}</code></nobr> | `-n`  | Limit to N most recent sessions      |
> | <nobr><code>{"--format"}</code></nobr>    |       | Output format: table or json (table) |
>
> ---
>
> #### delete
>
> Delete an OpenCode session.
>
> ```bash
> opencode session delete <sessionID>
> ```
>
> ---
>
> ### stats
>
> Show token usage and cost statistics for your OpenCode sessions.
>
> ```bash
> opencode stats
> ```
>
> #### Flags
>
> | Flag                                    | Description                                                                 |
> | --------------------------------------- | --------------------------------------------------------------------------- |
> | <nobr><code>{"--days"}</code></nobr>    | Show stats for the last N days (all time)                                   |
> | <nobr><code>{"--tools"}</code></nobr>   | Number of tools to show (all)                                               |
> | <nobr><code>{"--models"}</code></nobr>  | Show model usage breakdown (hidden by default). Pass a number to show top N |
> | <nobr><code>{"--project"}</code></nobr> | Filter by project (all projects, empty string: current project)             |
>
> ---
>
> ### export
>
> Export session data as JSON.
>
> ```bash
> opencode export [sessionID]
> ```
>
> If you don't provide a session ID, you'll be prompted to select from available sessions.
>
> #### Flags
>
> | Flag                                     | Description                           |
> | ---------------------------------------- | ------------------------------------- |
> | <nobr><code>{"--sanitize"}</code></nobr> | Redact sensitive transcript/file data |
>
> ---
>
> ### import
>
> Import session data from a JSON file or OpenCode share URL.
>
> ```bash
> opencode import <file>
> ```
>
> You can import from a local file or an OpenCode share URL.
>
> ```bash
> opencode import session.json
> opencode import https://opncd.ai/s/abc123
> ```
>
> ---
>
> ### web
>
> Start a headless OpenCode server with a web interface.
>
> ```bash
> opencode web
> ```
>
> This starts an HTTP server and opens a web browser to access OpenCode through a web interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### acp
>
> Start an ACP (Agent Client Protocol) server.
>
> ```bash
> opencode acp
> ```
>
> This command starts an ACP server that communicates via stdin/stdout using nd-JSON.
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--cwd"}</code></nobr>         | Working directory                          |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### plugin
>
> Install a plugin and update your config.
>
> ```bash
> opencode plugin <module>
> ```
>
> Or use the alias.
>
> ```bash
> opencode plug <module>
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                     |
> | -------------------------------------- | ----- | ------------------------------- |
> | <nobr><code>{"--global"}</code></nobr> | `-g`  | Install in global config        |
> | <nobr><code>{"--force"}</code></nobr>  | `-f`  | Replace existing plugin version |
>
> ---
>
> ### pr
>
> Fetch and checkout a GitHub PR branch, then run OpenCode.
>
> ```bash
> opencode pr <number>
> ```
>
> ---
>
> ### db
>
> Database tools.
>
> ```bash
> opencode db [query]
> ```
>
> #### Flags
>
> | Flag                                   | Description                    |
> | -------------------------------------- | ------------------------------ |
> | <nobr><code>{"--format"}</code></nobr> | Output format: `json` or `tsv` |
>
> ---
>
> #### path
>
> Print the database path.
>
> ```bash
> opencode db path
> ```
>
> ---
>
> ### debug
>
> Debugging and troubleshooting tools.
>
> ```bash
> opencode debug [command]
> ```
>
> ---
>
> ### uninstall
>
> Uninstall OpenCode and remove all related files.
>
> ```bash
> opencode uninstall
> ```
>
> #### Flags
>
> | Flag                                        | Short | Description                                 |
> | ------------------------------------------- | ----- | ------------------------------------------- |
> | <nobr><code>{"--keep-config"}</code></nobr> | `-c`  | Keep configuration files                    |
> | <nobr><code>{"--keep-data"}</code></nobr>   | `-d`  | Keep session data and snapshots             |
> | <nobr><code>{"--dry-run"}</code></nobr>     |       | Show what would be removed without removing |
> | <nobr><code>{"--force"}</code></nobr>       | `-f`  | Skip confirmation prompts                   |
>
> ---
>
> ### upgrade
>
> Updates opencode to the latest version or a specific version.
>
> ```bash
> opencode upgrade [target]
> ```
>
> To upgrade to the latest version.
>
> ```bash
> opencode upgrade
> ```
>
> To upgrade to a specific version.
>
> ```bash
> opencode upgrade v0.1.48
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                                                       |
> | -------------------------------------- | ----- | ----------------------------------------------------------------- |
> | <nobr><code>{"--method"}</code></nobr> | `-m`  | The installation method that was used; curl, npm, pnpm, bun, brew |
>
> ---
>
>
> | <nobr><code>{"--username"}</code></nobr> | `-u`  | Basic auth username (defaults to `OPENCODE_SERVER_USERNAME` or `opencode`) |
>
> ---
>
> ### auth
>
> Command to manage credentials and login for providers.
>
> ```bash
> opencode auth [command]
> ```
>
> ---
>
> #### login
>
> OpenCode is powered by the provider list at [Models.dev](https://models.dev), so you can use `opencode auth login` to configure API keys for any provider you'd like to use. This is stored in `~/.local/share/opencode/auth.json`.
>
> ```bash
> opencode auth login
> ```
>
> When OpenCode starts up it loads the providers from the credentials file. And if there are any keys defined in your environments or a `.env` file in your project.
>
> ##### Flags
>
> | Flag                                     | Short | Description                                          |
> | ---------------------------------------- | ----- | ---------------------------------------------------- |
> | <nobr><code>{"--provider"}</code></nobr> | `-p`  | Provider ID or name to log in to                     |
> | <nobr><code>{"--method"}</code></nobr>   | `-m`  | Login method label to use, skipping method selection |
>
> ---
>
> #### list
>
> Lists all the authenticated providers as stored in the credentials file.
>
> ```bash
> opencode auth list
> ```
>
> Or the short version.
>
> ```bash
> opencode auth ls
> ```
>
> ---
>
> #### logout
>
> Logs you out of a provider by clearing it from the credentials file.
>
> ```bash
> opencode auth logout
> ```
>
> ---
>
> ### github
>
> Manage the GitHub agent for repository automation.
>
> ```bash
> opencode github [command]
> ```
>
> ---
>
> #### install
>
> Install the GitHub agent in your repository.
>
> ```bash
> opencode github install
> ```
>
> This sets up the necessary GitHub Actions workflow and guides you through the configuration process. [Learn more](/docs/github).
>
> ---
>
> #### run
>
> Run the GitHub agent. This is typically used in GitHub Actions.
>
> ```bash
> opencode github run
> ```
>
> ##### Flags
>
> | Flag                                  | Description                            |
> | ------------------------------------- | -------------------------------------- |
> | <nobr><code>{"--event"}</code></nobr> | GitHub mock event to run the agent for |
> | <nobr><code>{"--token"}</code></nobr> | GitHub personal access token           |
>
> ---
>
> ### mcp
>
> Manage Model Context Protocol servers.
>
> ```bash
> opencode mcp [command]
> ```
>
> ---
>
> #### add
>
> Add an MCP server to your configuration.
>
> ```bash
> opencode mcp add
> ```
>
> This command will guide you through adding either a local or remote MCP server.
>
> ---
>
> #### list
>
> List all configured MCP servers and their connection status.
>
> ```bash
> opencode mcp list
> ```
>
> Or use the short version.
>
> ```bash
> opencode mcp ls
> ```
>
> ---
>
> #### auth
>
> Authenticate with an OAuth-enabled MCP server.
>
> ```bash
> opencode mcp auth [name]
> ```
>
> If you don't provide a server name, you'll be prompted to select from available OAuth-capable servers.
>
> You can also list OAuth-capable servers and their authentication status.
>
> ```bash
> opencode mcp auth list
> ```
>
> Or use the short version.
>
> ```bash
> opencode mcp auth ls
> ```
>
> ---
>
> #### logout
>
> Remove OAuth credentials for an MCP server.
>
> ```bash
> opencode mcp logout [name]
> ```
>
> ---
>
> #### debug
>
> Debug OAuth connection issues for an MCP server.
>
> ```bash
> opencode mcp debug <name>
> ```
>
> ---
>
> ### models
>
> List all available models from configured providers.
>
> ```bash
> opencode models [provider]
> ```
>
> This command displays all models available across your configured providers in the format `provider/model`.
>
> This is useful for figuring out the exact model name to use in [your config](/docs/config/).
>
> You can optionally pass a provider ID to filter models by that provider.
>
> ```bash
> opencode models anthropic
> ```
>
> #### Flags
>
> | Flag                                    | Description                                                  |
> | --------------------------------------- | ------------------------------------------------------------ |
> | <nobr><code>{"--refresh"}</code></nobr> | Refresh the models cache from models.dev                     |
> | <nobr><code>{"--verbose"}</code></nobr> | Use more verbose model output (includes metadata like costs) |
>
> Use the `--refresh` flag to update the cached model list. This is useful when new models have been added to a provider and you want to see them in OpenCode.
>
> ```bash
> opencode models --refresh
> ```
>
> ---
>
> ### run
>
> Run opencode in non-interactive mode by passing a prompt directly.
>
> ```bash
> opencode run [message..]
> ```
>
> This is useful for scripting, automation, or when you want a quick answer without launching the full TUI. For example.
>
> ```bash "opencode run"
> opencode run Explain the use of context in Go
> ```
>
> You can also attach to a running `opencode serve` instance to avoid MCP server cold boot times on every run:
>
> ```bash
> # Start a headless server in one terminal
> opencode serve
>
> # In another terminal, run commands that attach to it
> opencode run --attach http://localhost:4096 "Explain async/await in JavaScript"
> ```
>
> #### Flags
>
> | Flag                                     | Short | Description                                                                |
> | ---------------------------------------- | ----- | -------------------------------------------------------------------------- |
> | <nobr><code>{"--command"}</code></nobr>  |       | The command to run, use message for args                                   |
> | <nobr><code>{"--continue"}</code></nobr> | `-c`  | Continue the last session                                                  |
> | <nobr><code>{"--session"}</code></nobr>  | `-s`  | Session ID to continue                                                     |
> | <nobr><code>{"--fork"}</code></nobr>     |       | Fork the session when continuing (use with `--continue` or `--session`)    |
> | <nobr><code>{"--share"}</code></nobr>    |       | Share the session                                                          |
> | <nobr><code>{"--model"}</code></nobr>    | `-m`  | Model to use in the form of provider/model                                 |
> | <nobr><code>{"--agent"}</code></nobr>    |       | Agent to use                                                               |
> | <nobr><code>{"--file"}</code></nobr>     | `-f`  | File(s) to attach to message                                               |
> | <nobr><code>{"--format"}</code></nobr>   |       | Format: default (formatted) or json (raw JSON events)                      |
> | <nobr><code>{"--title"}</code></nobr>    |       | Title for the session (uses truncated prompt if no value provided)         |
> | <nobr><code>{"--attach"}</code></nobr>   |       | Attach to a running opencode server (e.g., http://localhost:4096)          |
> | <nobr><code>{"--password"}</code></nobr> | `-p`  | Basic auth password (defaults to `OPENCODE_SERVER_PASSWORD`)               |
> | <nobr><code>{"--username"}</code></nobr> | `-u`  | Basic auth username (defaults to `OPENCODE_SERVER_USERNAME` or `opencode`) |
> | <nobr><code>{"--dir"}</code></nobr>      |       | Directory to run in, or path on the remote server when attaching           |
> | <nobr><code>{"--port"}</code></nobr>     |       | Port for the local server (defaults to random port)                        |
> | <nobr><code>{"--variant"}</code></nobr>  |       | Model variant (provider-specific reasoning effort)                         |
> | <nobr><code>{"--thinking"}</code></nobr> |       | Show thinking blocks                                                       |
> | <nobr><code>{"--auto"}</code></nobr>     |       | Auto-approve permissions that are not explicitly denied                    |
>
> ---
>
> ### serve
>
> Start a headless OpenCode server for API access. Check out the [server docs](/docs/server) for the full HTTP interface.
>
> ```bash
> opencode serve
> ```
>
> This starts an HTTP server that provides API access to opencode functionality without the TUI interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### session
>
> Manage OpenCode sessions.
>
> ```bash
> opencode session [command]
> ```
>
> ---
>
> #### list
>
> List all OpenCode sessions.
>
> ```bash
> opencode session list
> ```
>
> ##### Flags
>
> | Flag                                      | Short | Description                          |
> | ----------------------------------------- | ----- | ------------------------------------ |
> | <nobr><code>{"--max-count"}</code></nobr> | `-n`  | Limit to N most recent sessions      |
> | <nobr><code>{"--format"}</code></nobr>    |       | Output format: table or json (table) |
>
> ---
>
> #### delete
>
> Delete an OpenCode session.
>
> ```bash
> opencode session delete <sessionID>
> ```
>
> ---
>
> ### stats
>
> Show token usage and cost statistics for your OpenCode sessions.
>
> ```bash
> opencode stats
> ```
>
> #### Flags
>
> | Flag                                    | Description                                                                 |
> | --------------------------------------- | --------------------------------------------------------------------------- |
> | <nobr><code>{"--days"}</code></nobr>    | Show stats for the last N days (all time)                                   |
> | <nobr><code>{"--tools"}</code></nobr>   | Number of tools to show (all)                                               |
> | <nobr><code>{"--models"}</code></nobr>  | Show model usage breakdown (hidden by default). Pass a number to show top N |
> | <nobr><code>{"--project"}</code></nobr> | Filter by project (all projects, empty string: current project)             |
>
> ---
>
> ### export
>
> Export session data as JSON.
>
> ```bash
> opencode export [sessionID]
> ```
>
> If you don't provide a session ID, you'll be prompted to select from available sessions.
>
> #### Flags
>
> | Flag                                     | Description                           |
> | ---------------------------------------- | ------------------------------------- |
> | <nobr><code>{"--sanitize"}</code></nobr> | Redact sensitive transcript/file data |
>
> ---
>
> ### import
>
> Import session data from a JSON file or OpenCode share URL.
>
> ```bash
> opencode import <file>
> ```
>
> You can import from a local file or an OpenCode share URL.
>
> ```bash
> opencode import session.json
> opencode import https://opncd.ai/s/abc123
> ```
>
> ---
>
> ### web
>
> Start a headless OpenCode server with a web interface.
>
> ```bash
> opencode web
> ```
>
> This starts an HTTP server and opens a web browser to access OpenCode through a web interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### acp
>
> Start an ACP (Agent Client Protocol) server.
>
> ```bash
> opencode acp
> ```
>
> This command starts an ACP server that communicates via stdin/stdout using nd-JSON.
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--cwd"}</code></nobr>         | Working directory                          |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### plugin
>
> Install a plugin and update your config.
>
> ```bash
> opencode plugin <module>
> ```
>
> Or use the alias.
>
> ```bash
> opencode plug <module>
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                     |
> | -------------------------------------- | ----- | ------------------------------- |
> | <nobr><code>{"--global"}</code></nobr> | `-g`  | Install in global config        |
> | <nobr><code>{"--force"}</code></nobr>  | `-f`  | Replace existing plugin version |
>
> ---
>
> ### pr
>
> Fetch and checkout a GitHub PR branch, then run OpenCode.
>
> ```bash
> opencode pr <number>
> ```
>
> ---
>
> ### db
>
> Database tools.
>
> ```bash
> opencode db [query]
> ```
>
> #### Flags
>
> | Flag                                   | Description                    |
> | -------------------------------------- | ------------------------------ |
> | <nobr><code>{"--format"}</code></nobr> | Output format: `json` or `tsv` |
>
> ---
>
> #### path
>
> Print the database path.
>
> ```bash
> opencode db path
> ```
>
> ---
>
> ### debug
>
> Debugging and troubleshooting tools.
>
> ```bash
> opencode debug [command]
> ```
>
> ---
>
> ### uninstall
>
> Uninstall OpenCode and remove all related files.
>
> ```bash
> opencode uninstall
> ```
>
> #### Flags
>
> | Flag                                        | Short | Description                                 |
> | ------------------------------------------- | ----- | ------------------------------------------- |
> | <nobr><code>{"--keep-config"}</code></nobr> | `-c`  | Keep configuration files                    |
> | <nobr><code>{"--keep-data"}</code></nobr>   | `-d`  | Keep session data and snapshots             |
> | <nobr><code>{"--dry-run"}</code></nobr>     |       | Show what would be removed without removing |
> | <nobr><code>{"--force"}</code></nobr>       | `-f`  | Skip confirmation prompts                   |
>
> ---
>
> ### upgrade
>
> Updates opencode to the latest version or a specific version.
>
> ```bash
> opencode upgrade [target]
> ```
>
> To upgrade to the latest version.
>
> ```bash
> opencode upgrade
> ```
>
> To upgrade to a specific version.
>
> ```bash
> opencode upgrade v0.1.48
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                                                       |
> | -------------------------------------- | ----- | ----------------------------------------------------------------- |
> | <nobr><code>{"--method"}</code></nobr> | `-m`  | The installation method that was used; curl, npm, pnpm, bun, brew |
>
> ---
>
>
> | <nobr><code>{"--password"}</code></nobr> | `-p`  | Basic auth password (defaults to `OPENCODE_SERVER_PASSWORD`)               |
> | <nobr><code>{"--username"}</code></nobr> | `-u`  | Basic auth username (defaults to `OPENCODE_SERVER_USERNAME` or `opencode`) |
> | <nobr><code>{"--dir"}</code></nobr>      |       | Directory to run in, or path on the remote server when attaching           |
> | <nobr><code>{"--port"}</code></nobr>     |       | Port for the local server (defaults to random port)                        |
> | <nobr><code>{"--variant"}</code></nobr>  |       | Model variant (provider-specific reasoning effort)                         |
> | <nobr><code>{"--thinking"}</code></nobr> |       | Show thinking blocks                                                       |
> | <nobr><code>{"--auto"}</code></nobr>     |       | Auto-approve permissions that are not explicitly denied                    |
>
> ---
>
> ### serve
>
> Start a headless OpenCode server for API access. Check out the [server docs](/docs/server) for the full HTTP interface.
>
> ```bash
> opencode serve
> ```
>
> This starts an HTTP server that provides API access to opencode functionality without the TUI interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### session
>
> Manage OpenCode sessions.
>
> ```bash
> opencode session [command]
> ```
>
> ---
>
> #### list
>
> List all OpenCode sessions.
>
> ```bash
> opencode session list
> ```
>
> ##### Flags
>
> | Flag                                      | Short | Description                          |
> | ----------------------------------------- | ----- | ------------------------------------ |
> | <nobr><code>{"--max-count"}</code></nobr> | `-n`  | Limit to N most recent sessions      |
> | <nobr><code>{"--format"}</code></nobr>    |       | Output format: table or json (table) |
>
> ---
>
> #### delete
>
> Delete an OpenCode session.
>
> ```bash
> opencode session delete <sessionID>
> ```
>
> ---
>
> ### stats
>
> Show token usage and cost statistics for your OpenCode sessions.
>
> ```bash
> opencode stats
> ```
>
> #### Flags
>
> | Flag                                    | Description                                                                 |
> | --------------------------------------- | --------------------------------------------------------------------------- |
> | <nobr><code>{"--days"}</code></nobr>    | Show stats for the last N days (all time)                                   |
> | <nobr><code>{"--tools"}</code></nobr>   | Number of tools to show (all)                                               |
> | <nobr><code>{"--models"}</code></nobr>  | Show model usage breakdown (hidden by default). Pass a number to show top N |
> | <nobr><code>{"--project"}</code></nobr> | Filter by project (all projects, empty string: current project)             |
>
> ---
>
> ### export
>
> Export session data as JSON.
>
> ```bash
> opencode export [sessionID]
> ```
>
> If you don't provide a session ID, you'll be prompted to select from available sessions.
>
> #### Flags
>
> | Flag                                     | Description                           |
> | ---------------------------------------- | ------------------------------------- |
> | <nobr><code>{"--sanitize"}</code></nobr> | Redact sensitive transcript/file data |
>
> ---
>
> ### import
>
> Import session data from a JSON file or OpenCode share URL.
>
> ```bash
> opencode import <file>
> ```
>
> You can import from a local file or an OpenCode share URL.
>
> ```bash
> opencode import session.json
> opencode import https://opncd.ai/s/abc123
> ```
>
> ---
>
> ### web
>
> Start a headless OpenCode server with a web interface.
>
> ```bash
> opencode web
> ```
>
> This starts an HTTP server and opens a web browser to access OpenCode through a web interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### acp
>
> Start an ACP (Agent Client Protocol) server.
>
> ```bash
> opencode acp
> ```
>
> This command starts an ACP server that communicates via stdin/stdout using nd-JSON.
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--cwd"}</code></nobr>         | Working directory                          |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### plugin
>
> Install a plugin and update your config.
>
> ```bash
> opencode plugin <module>
> ```
>
> Or use the alias.
>
> ```bash
> opencode plug <module>
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                     |
> | -------------------------------------- | ----- | ------------------------------- |
> | <nobr><code>{"--global"}</code></nobr> | `-g`  | Install in global config        |
> | <nobr><code>{"--force"}</code></nobr>  | `-f`  | Replace existing plugin version |
>
> ---
>
> ### pr
>
> Fetch and checkout a GitHub PR branch, then run OpenCode.
>
> ```bash
> opencode pr <number>
> ```
>
> ---
>
> ### db
>
> Database tools.
>
> ```bash
> opencode db [query]
> ```
>
> #### Flags
>
> | Flag                                   | Description                    |
> | -------------------------------------- | ------------------------------ |
> | <nobr><code>{"--format"}</code></nobr> | Output format: `json` or `tsv` |
>
> ---
>
> #### path
>
> Print the database path.
>
> ```bash
> opencode db path
> ```
>
> ---
>
> ### debug
>
> Debugging and troubleshooting tools.
>
> ```bash
> opencode debug [command]
> ```
>
> ---
>
> ### uninstall
>
> Uninstall OpenCode and remove all related files.
>
> ```bash
> opencode uninstall
> ```
>
> #### Flags
>
> | Flag                                        | Short | Description                                 |
> | ------------------------------------------- | ----- | ------------------------------------------- |
> | <nobr><code>{"--keep-config"}</code></nobr> | `-c`  | Keep configuration files                    |
> | <nobr><code>{"--keep-data"}</code></nobr>   | `-d`  | Keep session data and snapshots             |
> | <nobr><code>{"--dry-run"}</code></nobr>     |       | Show what would be removed without removing |
> | <nobr><code>{"--force"}</code></nobr>       | `-f`  | Skip confirmation prompts                   |
>
> ---
>
> ### upgrade
>
> Updates opencode to the latest version or a specific version.
>
> ```bash
> opencode upgrade [target]
> ```
>
> To upgrade to the latest version.
>
> ```bash
> opencode upgrade
> ```
>
> To upgrade to a specific version.
>
> ```bash
> opencode upgrade v0.1.48
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                                                       |
> | -------------------------------------- | ----- | ----------------------------------------------------------------- |
> | <nobr><code>{"--method"}</code></nobr> | `-m`  | The installation method that was used; curl, npm, pnpm, bun, brew |
>
> ---
>
>
> | <nobr><code>{"--username"}</code></nobr> | `-u`  | Basic auth username (defaults to `OPENCODE_SERVER_USERNAME` or `opencode`) |
> | <nobr><code>{"--dir"}</code></nobr>      |       | Directory to run in, or path on the remote server when attaching           |
> | <nobr><code>{"--port"}</code></nobr>     |       | Port for the local server (defaults to random port)                        |
> | <nobr><code>{"--variant"}</code></nobr>  |       | Model variant (provider-specific reasoning effort)                         |
> | <nobr><code>{"--thinking"}</code></nobr> |       | Show thinking blocks                                                       |
> | <nobr><code>{"--auto"}</code></nobr>     |       | Auto-approve permissions that are not explicitly denied                    |
>
> ---
>
> ### serve
>
> Start a headless OpenCode server for API access. Check out the [server docs](/docs/server) for the full HTTP interface.
>
> ```bash
> opencode serve
> ```
>
> This starts an HTTP server that provides API access to opencode functionality without the TUI interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### session
>
> Manage OpenCode sessions.
>
> ```bash
> opencode session [command]
> ```
>
> ---
>
> #### list
>
> List all OpenCode sessions.
>
> ```bash
> opencode session list
> ```
>
> ##### Flags
>
> | Flag                                      | Short | Description                          |
> | ----------------------------------------- | ----- | ------------------------------------ |
> | <nobr><code>{"--max-count"}</code></nobr> | `-n`  | Limit to N most recent sessions      |
> | <nobr><code>{"--format"}</code></nobr>    |       | Output format: table or json (table) |
>
> ---
>
> #### delete
>
> Delete an OpenCode session.
>
> ```bash
> opencode session delete <sessionID>
> ```
>
> ---
>
> ### stats
>
> Show token usage and cost statistics for your OpenCode sessions.
>
> ```bash
> opencode stats
> ```
>
> #### Flags
>
> | Flag                                    | Description                                                                 |
> | --------------------------------------- | --------------------------------------------------------------------------- |
> | <nobr><code>{"--days"}</code></nobr>    | Show stats for the last N days (all time)                                   |
> | <nobr><code>{"--tools"}</code></nobr>   | Number of tools to show (all)                                               |
> | <nobr><code>{"--models"}</code></nobr>  | Show model usage breakdown (hidden by default). Pass a number to show top N |
> | <nobr><code>{"--project"}</code></nobr> | Filter by project (all projects, empty string: current project)             |
>
> ---
>
> ### export
>
> Export session data as JSON.
>
> ```bash
> opencode export [sessionID]
> ```
>
> If you don't provide a session ID, you'll be prompted to select from available sessions.
>
> #### Flags
>
> | Flag                                     | Description                           |
> | ---------------------------------------- | ------------------------------------- |
> | <nobr><code>{"--sanitize"}</code></nobr> | Redact sensitive transcript/file data |
>
> ---
>
> ### import
>
> Import session data from a JSON file or OpenCode share URL.
>
> ```bash
> opencode import <file>
> ```
>
> You can import from a local file or an OpenCode share URL.
>
> ```bash
> opencode import session.json
> opencode import https://opncd.ai/s/abc123
> ```
>
> ---
>
> ### web
>
> Start a headless OpenCode server with a web interface.
>
> ```bash
> opencode web
> ```
>
> This starts an HTTP server and opens a web browser to access OpenCode through a web interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### acp
>
> Start an ACP (Agent Client Protocol) server.
>
> ```bash
> opencode acp
> ```
>
> This command starts an ACP server that communicates via stdin/stdout using nd-JSON.
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--cwd"}</code></nobr>         | Working directory                          |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### plugin
>
> Install a plugin and update your config.
>
> ```bash
> opencode plugin <module>
> ```
>
> Or use the alias.
>
> ```bash
> opencode plug <module>
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                     |
> | -------------------------------------- | ----- | ------------------------------- |
> | <nobr><code>{"--global"}</code></nobr> | `-g`  | Install in global config        |
> | <nobr><code>{"--force"}</code></nobr>  | `-f`  | Replace existing plugin version |
>
> ---
>
> ### pr
>
> Fetch and checkout a GitHub PR branch, then run OpenCode.
>
> ```bash
> opencode pr <number>
> ```
>
> ---
>
> ### db
>
> Database tools.
>
> ```bash
> opencode db [query]
> ```
>
> #### Flags
>
> | Flag                                   | Description                    |
> | -------------------------------------- | ------------------------------ |
> | <nobr><code>{"--format"}</code></nobr> | Output format: `json` or `tsv` |
>
> ---
>
> #### path
>
> Print the database path.
>
> ```bash
> opencode db path
> ```
>
> ---
>
> ### debug
>
> Debugging and troubleshooting tools.
>
> ```bash
> opencode debug [command]
> ```
>
> ---
>
> ### uninstall
>
> Uninstall OpenCode and remove all related files.
>
> ```bash
> opencode uninstall
> ```
>
> #### Flags
>
> | Flag                                        | Short | Description                                 |
> | ------------------------------------------- | ----- | ------------------------------------------- |
> | <nobr><code>{"--keep-config"}</code></nobr> | `-c`  | Keep configuration files                    |
> | <nobr><code>{"--keep-data"}</code></nobr>   | `-d`  | Keep session data and snapshots             |
> | <nobr><code>{"--dry-run"}</code></nobr>     |       | Show what would be removed without removing |
> | <nobr><code>{"--force"}</code></nobr>       | `-f`  | Skip confirmation prompts                   |
>
> ---
>
> ### upgrade
>
> Updates opencode to the latest version or a specific version.
>
> ```bash
> opencode upgrade [target]
> ```
>
> To upgrade to the latest version.
>
> ```bash
> opencode upgrade
> ```
>
> To upgrade to a specific version.
>
> ```bash
> opencode upgrade v0.1.48
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                                                       |
> | -------------------------------------- | ----- | ----------------------------------------------------------------- |
> | <nobr><code>{"--method"}</code></nobr> | `-m`  | The installation method that was used; curl, npm, pnpm, bun, brew |
>
> ---
>
>
> This starts an HTTP server that provides API access to opencode functionality without the TUI interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### session
>
> Manage OpenCode sessions.
>
> ```bash
> opencode session [command]
> ```
>
> ---
>
> #### list
>
> List all OpenCode sessions.
>
> ```bash
> opencode session list
> ```
>
> ##### Flags
>
> | Flag                                      | Short | Description                          |
> | ----------------------------------------- | ----- | ------------------------------------ |
> | <nobr><code>{"--max-count"}</code></nobr> | `-n`  | Limit to N most recent sessions      |
> | <nobr><code>{"--format"}</code></nobr>    |       | Output format: table or json (table) |
>
> ---
>
> #### delete
>
> Delete an OpenCode session.
>
> ```bash
> opencode session delete <sessionID>
> ```
>
> ---
>
> ### stats
>
> Show token usage and cost statistics for your OpenCode sessions.
>
> ```bash
> opencode stats
> ```
>
> #### Flags
>
> | Flag                                    | Description                                                                 |
> | --------------------------------------- | --------------------------------------------------------------------------- |
> | <nobr><code>{"--days"}</code></nobr>    | Show stats for the last N days (all time)                                   |
> | <nobr><code>{"--tools"}</code></nobr>   | Number of tools to show (all)                                               |
> | <nobr><code>{"--models"}</code></nobr>  | Show model usage breakdown (hidden by default). Pass a number to show top N |
> | <nobr><code>{"--project"}</code></nobr> | Filter by project (all projects, empty string: current project)             |
>
> ---
>
> ### export
>
> Export session data as JSON.
>
> ```bash
> opencode export [sessionID]
> ```
>
> If you don't provide a session ID, you'll be prompted to select from available sessions.
>
> #### Flags
>
> | Flag                                     | Description                           |
> | ---------------------------------------- | ------------------------------------- |
> | <nobr><code>{"--sanitize"}</code></nobr> | Redact sensitive transcript/file data |
>
> ---
>
> ### import
>
> Import session data from a JSON file or OpenCode share URL.
>
> ```bash
> opencode import <file>
> ```
>
> You can import from a local file or an OpenCode share URL.
>
> ```bash
> opencode import session.json
> opencode import https://opncd.ai/s/abc123
> ```
>
> ---
>
> ### web
>
> Start a headless OpenCode server with a web interface.
>
> ```bash
> opencode web
> ```
>
> This starts an HTTP server and opens a web browser to access OpenCode through a web interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### acp
>
> Start an ACP (Agent Client Protocol) server.
>
> ```bash
> opencode acp
> ```
>
> This command starts an ACP server that communicates via stdin/stdout using nd-JSON.
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--cwd"}</code></nobr>         | Working directory                          |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### plugin
>
> Install a plugin and update your config.
>
> ```bash
> opencode plugin <module>
> ```
>
> Or use the alias.
>
> ```bash
> opencode plug <module>
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                     |
> | -------------------------------------- | ----- | ------------------------------- |
> | <nobr><code>{"--global"}</code></nobr> | `-g`  | Install in global config        |
> | <nobr><code>{"--force"}</code></nobr>  | `-f`  | Replace existing plugin version |
>
> ---
>
> ### pr
>
> Fetch and checkout a GitHub PR branch, then run OpenCode.
>
> ```bash
> opencode pr <number>
> ```
>
> ---
>
> ### db
>
> Database tools.
>
> ```bash
> opencode db [query]
> ```
>
> #### Flags
>
> | Flag                                   | Description                    |
> | -------------------------------------- | ------------------------------ |
> | <nobr><code>{"--format"}</code></nobr> | Output format: `json` or `tsv` |
>
> ---
>
> #### path
>
> Print the database path.
>
> ```bash
> opencode db path
> ```
>
> ---
>
> ### debug
>
> Debugging and troubleshooting tools.
>
> ```bash
> opencode debug [command]
> ```
>
> ---
>
> ### uninstall
>
> Uninstall OpenCode and remove all related files.
>
> ```bash
> opencode uninstall
> ```
>
> #### Flags
>
> | Flag                                        | Short | Description                                 |
> | ------------------------------------------- | ----- | ------------------------------------------- |
> | <nobr><code>{"--keep-config"}</code></nobr> | `-c`  | Keep configuration files                    |
> | <nobr><code>{"--keep-data"}</code></nobr>   | `-d`  | Keep session data and snapshots             |
> | <nobr><code>{"--dry-run"}</code></nobr>     |       | Show what would be removed without removing |
> | <nobr><code>{"--force"}</code></nobr>       | `-f`  | Skip confirmation prompts                   |
>
> ---
>
> ### upgrade
>
> Updates opencode to the latest version or a specific version.
>
> ```bash
> opencode upgrade [target]
> ```
>
> To upgrade to the latest version.
>
> ```bash
> opencode upgrade
> ```
>
> To upgrade to a specific version.
>
> ```bash
> opencode upgrade v0.1.48
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                                                       |
> | -------------------------------------- | ----- | ----------------------------------------------------------------- |
> | <nobr><code>{"--method"}</code></nobr> | `-m`  | The installation method that was used; curl, npm, pnpm, bun, brew |
>
> ---
>
>
> This starts an HTTP server and opens a web browser to access OpenCode through a web interface. Set `OPENCODE_SERVER_PASSWORD` to enable HTTP basic auth (username defaults to `opencode`).
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### acp
>
> Start an ACP (Agent Client Protocol) server.
>
> ```bash
> opencode acp
> ```
>
> This command starts an ACP server that communicates via stdin/stdout using nd-JSON.
>
> #### Flags
>
> | Flag                                        | Description                                |
> | ------------------------------------------- | ------------------------------------------ |
> | <nobr><code>{"--cwd"}</code></nobr>         | Working directory                          |
> | <nobr><code>{"--port"}</code></nobr>        | Port to listen on                          |
> | <nobr><code>{"--hostname"}</code></nobr>    | Hostname to listen on                      |
> | <nobr><code>{"--mdns"}</code></nobr>        | Enable mDNS discovery                      |
> | <nobr><code>{"--mdns-domain"}</code></nobr> | Custom mDNS domain name                    |
> | <nobr><code>{"--cors"}</code></nobr>        | Additional browser origin(s) to allow CORS |
>
> ---
>
> ### plugin
>
> Install a plugin and update your config.
>
> ```bash
> opencode plugin <module>
> ```
>
> Or use the alias.
>
> ```bash
> opencode plug <module>
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                     |
> | -------------------------------------- | ----- | ------------------------------- |
> | <nobr><code>{"--global"}</code></nobr> | `-g`  | Install in global config        |
> | <nobr><code>{"--force"}</code></nobr>  | `-f`  | Replace existing plugin version |
>
> ---
>
> ### pr
>
> Fetch and checkout a GitHub PR branch, then run OpenCode.
>
> ```bash
> opencode pr <number>
> ```
>
> ---
>
> ### db
>
> Database tools.
>
> ```bash
> opencode db [query]
> ```
>
> #### Flags
>
> | Flag                                   | Description                    |
> | -------------------------------------- | ------------------------------ |
> | <nobr><code>{"--format"}</code></nobr> | Output format: `json` or `tsv` |
>
> ---
>
> #### path
>
> Print the database path.
>
> ```bash
> opencode db path
> ```
>
> ---
>
> ### debug
>
> Debugging and troubleshooting tools.
>
> ```bash
> opencode debug [command]
> ```
>
> ---
>
> ### uninstall
>
> Uninstall OpenCode and remove all related files.
>
> ```bash
> opencode uninstall
> ```
>
> #### Flags
>
> | Flag                                        | Short | Description                                 |
> | ------------------------------------------- | ----- | ------------------------------------------- |
> | <nobr><code>{"--keep-config"}</code></nobr> | `-c`  | Keep configuration files                    |
> | <nobr><code>{"--keep-data"}</code></nobr>   | `-d`  | Keep session data and snapshots             |
> | <nobr><code>{"--dry-run"}</code></nobr>     |       | Show what would be removed without removing |
> | <nobr><code>{"--force"}</code></nobr>       | `-f`  | Skip confirmation prompts                   |
>
> ---
>
> ### upgrade
>
> Updates opencode to the latest version or a specific version.
>
> ```bash
> opencode upgrade [target]
> ```
>
> To upgrade to the latest version.
>
> ```bash
> opencode upgrade
> ```
>
> To upgrade to a specific version.
>
> ```bash
> opencode upgrade v0.1.48
> ```
>
> #### Flags
>
> | Flag                                   | Short | Description                                                       |
> | -------------------------------------- | ----- | ----------------------------------------------------------------- |
> | <nobr><code>{"--method"}</code></nobr> | `-m`  | The installation method that was used; curl, npm, pnpm, bun, brew |
>
> ---
>
>
> ## Environment variables
>
> OpenCode can be configured using environment variables.
>
> | Variable                              | Type    | Description                                       |
> | ------------------------------------- | ------- | ------------------------------------------------- |
> | `OPENCODE_AUTO_SHARE`                 | boolean | Automatically share sessions                      |
> | `OPENCODE_GIT_BASH_PATH`              | string  | Path to Git Bash executable on Windows            |
> | `OPENCODE_CONFIG`                     | string  | Path to config file                               |
> | `OPENCODE_TUI_CONFIG`                 | string  | Path to TUI config file                           |
> | `OPENCODE_CONFIG_DIR`                 | string  | Path to config directory                          |
> | `OPENCODE_CONFIG_CONTENT`             | string  | Inline json config content                        |
> | `OPENCODE_DISABLE_AUTOUPDATE`         | boolean | Disable automatic update checks                   |
> | `OPENCODE_DISABLE_PRUNE`              | boolean | Disable pruning of old data                       |
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_AUTO_SHARE`                 | boolean | Automatically share sessions                      |
> | `OPENCODE_GIT_BASH_PATH`              | string  | Path to Git Bash executable on Windows            |
> | `OPENCODE_CONFIG`                     | string  | Path to config file                               |
> | `OPENCODE_TUI_CONFIG`                 | string  | Path to TUI config file                           |
> | `OPENCODE_CONFIG_DIR`                 | string  | Path to config directory                          |
> | `OPENCODE_CONFIG_CONTENT`             | string  | Inline json config content                        |
> | `OPENCODE_DISABLE_AUTOUPDATE`         | boolean | Disable automatic update checks                   |
> | `OPENCODE_DISABLE_PRUNE`              | boolean | Disable pruning of old data                       |
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_GIT_BASH_PATH`              | string  | Path to Git Bash executable on Windows            |
> | `OPENCODE_CONFIG`                     | string  | Path to config file                               |
> | `OPENCODE_TUI_CONFIG`                 | string  | Path to TUI config file                           |
> | `OPENCODE_CONFIG_DIR`                 | string  | Path to config directory                          |
> | `OPENCODE_CONFIG_CONTENT`             | string  | Inline json config content                        |
> | `OPENCODE_DISABLE_AUTOUPDATE`         | boolean | Disable automatic update checks                   |
> | `OPENCODE_DISABLE_PRUNE`              | boolean | Disable pruning of old data                       |
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_CONFIG`                     | string  | Path to config file                               |
> | `OPENCODE_TUI_CONFIG`                 | string  | Path to TUI config file                           |
> | `OPENCODE_CONFIG_DIR`                 | string  | Path to config directory                          |
> | `OPENCODE_CONFIG_CONTENT`             | string  | Inline json config content                        |
> | `OPENCODE_DISABLE_AUTOUPDATE`         | boolean | Disable automatic update checks                   |
> | `OPENCODE_DISABLE_PRUNE`              | boolean | Disable pruning of old data                       |
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_TUI_CONFIG`                 | string  | Path to TUI config file                           |
> | `OPENCODE_CONFIG_DIR`                 | string  | Path to config directory                          |
> | `OPENCODE_CONFIG_CONTENT`             | string  | Inline json config content                        |
> | `OPENCODE_DISABLE_AUTOUPDATE`         | boolean | Disable automatic update checks                   |
> | `OPENCODE_DISABLE_PRUNE`              | boolean | Disable pruning of old data                       |
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_CONFIG_DIR`                 | string  | Path to config directory                          |
> | `OPENCODE_CONFIG_CONTENT`             | string  | Inline json config content                        |
> | `OPENCODE_DISABLE_AUTOUPDATE`         | boolean | Disable automatic update checks                   |
> | `OPENCODE_DISABLE_PRUNE`              | boolean | Disable pruning of old data                       |
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_CONFIG_CONTENT`             | string  | Inline json config content                        |
> | `OPENCODE_DISABLE_AUTOUPDATE`         | boolean | Disable automatic update checks                   |
> | `OPENCODE_DISABLE_PRUNE`              | boolean | Disable pruning of old data                       |
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_AUTOUPDATE`         | boolean | Disable automatic update checks                   |
> | `OPENCODE_DISABLE_PRUNE`              | boolean | Disable pruning of old data                       |
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_PRUNE`              | boolean | Disable pruning of old data                       |
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_TERMINAL_TITLE`     | boolean | Disable automatic terminal title updates          |
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_PERMISSION`                 | string  | Inlined json permissions config                   |
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_DEFAULT_PLUGINS`    | boolean | Disable default plugins                           |
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_LSP_DOWNLOAD`       | boolean | Disable automatic LSP server downloads            |
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_ENABLE_EXPERIMENTAL_MODELS` | boolean | Enable experimental models                        |
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_AUTOCOMPACT`        | boolean | Disable automatic context compaction              |
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_CLAUDE_CODE`        | boolean | Disable reading from `.claude` (prompt + skills)  |
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_CLAUDE_CODE_PROMPT` | boolean | Disable reading `~/.claude/CLAUDE.md`             |
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS` | boolean | Disable loading `.claude/skills`                  |
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_MODELS_FETCH`       | boolean | Disable fetching models from remote sources       |
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_DISABLE_MOUSE`              | boolean | Disable mouse capture in the TUI                  |
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_FAKE_VCS`                   | string  | Fake VCS provider for testing purposes            |
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_CLIENT`                     | string  | Client identifier (defaults to `cli`)             |
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_ENABLE_EXA`                 | boolean | Enable Exa web search tools                       |
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_ENABLE_PARALLEL`            | boolean | Enable Parallel web search tools                  |
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_SERVER_PASSWORD`            | string  | Enable basic auth for `serve`/`web`               |
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_SERVER_USERNAME`            | string  | Override basic auth username (default `opencode`) |
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_MODELS_URL`                 | string  | Custom URL for fetching models configuration      |
>
> ---
>
> ### Experimental
>
> These environment variables enable experimental features that may change or be removed.
>
> | Variable                                        | Type    | Description                             |
> | ----------------------------------------------- | ------- | --------------------------------------- |
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL`                         | boolean | Enable the experimental umbrella flag   |
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_ICON_DISCOVERY`          | boolean | Enable icon discovery                   |
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_DISABLE_COPY_ON_SELECT`  | boolean | Disable copy on select in TUI           |
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_BASH_DEFAULT_TIMEOUT_MS` | number  | Default timeout for bash commands in ms |
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_OUTPUT_TOKEN_MAX`        | number  | Max output tokens for LLM responses     |
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_FILEWATCHER`             | boolean | Enable file watcher for entire dir      |
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_OXFMT`                   | boolean | Enable oxfmt formatter                  |
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_LSP_TOOL`                | boolean | Enable experimental LSP tool            |
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_DISABLE_FILEWATCHER`     | boolean | Disable file watcher                    |
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_EXA`                     | boolean | Enable experimental Exa features        |
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_LSP_TY`                  | boolean | Enable TY LSP for python files          |
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_PLAN_MODE`               | boolean | Enable plan mode                        |
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_BACKGROUND_SUBAGENTS`    | boolean | Enable background subagent tasks        |
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_EVENT_SYSTEM`            | boolean | Enable experimental event system        |
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_NATIVE_LLM`              | boolean | Enable native LLM request path          |
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_PARALLEL`                | boolean | Enable parallel web search execution    |
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_SCOUT`                   | boolean | Enable Scout subagent                   |
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |
>
> | `OPENCODE_EXPERIMENTAL_WORKSPACES`              | boolean | Enable workspace support                |

### Source: OpenCode Config — Experimental env-related

> 3. **Custom config** (`OPENCODE_CONFIG` env var) - custom overrides
> 4. **Project config** (`opencode.json` in project) - project-specific settings
> 5. **`.opencode` directories** - agents, commands, plugins
> 6. **Inline config** (`OPENCODE_CONFIG_CONTENT` env var) - runtime overrides
> 7. **Managed config files** (`/Library/Application Support/opencode/` on macOS) - admin-controlled
> 8. **macOS managed preferences** (`.mobileconfig` via MDM) - highest priority, not user-overridable
>
> This means project configs can override global defaults, and global configs can override remote organizational defaults. Managed settings override everything.
>
> :::note
> The `.opencode` and `~/.config/opencode` directories use **plural names** for subdirectories: `agents/`, `commands/`, `modes/`, `plugins/`, `skills/`, `tools/`, and `themes/`. Singular names (e.g., `agent/`) are also supported for backwards compatibility.
> :::
>
> ---
>
> ### Remote
>
> Organizations can provide default configuration via the `.well-known/opencode` endpoint. This is fetched automatically when you authenticate with a provider that supports it.
>
> Remote config is loaded first, serving as the base layer. All other config sources (global, project) can override these defaults.
>
> For example, if your organization provides MCP servers that are disabled by default:
>
> ```json title="Remote config from .well-known/opencode"
> {
>   "mcp": {
>     "jira": {
>       "type": "remote",
>       "url": "https://jira.example.com/mcp",
>       "enabled": false
>     }
>   }
> }
> ```
>
> You can enable specific servers in your local config:
>
> ```json title="opencode.json"
> {
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
> ---
>
> ### Global
>
> Place your global OpenCode config in `~/.config/opencode/opencode.json`. Use global config for user-wide server/runtime preferences like providers, models, and permissions.
>
> For TUI-specific settings, use `~/.config/opencode/tui.json`.
>
> Global config overrides remote organizational defaults.
>
> ---
>
> ### Per project
>
> Add `opencode.json` in your project root. Project config has the highest precedence among standard config files - it overrides both global and remote configs.
>
> For project-specific TUI settings, add `tui.json` alongside it.
>
> :::tip
> Place project specific config in the root of your project.
> :::
>
> When OpenCode starts up, it first looks for a config file in the current directory, then traverses up to the nearest Git directory.
>
> This is also safe to be checked into Git and uses the same schema as the global one.
>
> ---
>
> ### Custom path
>
> Specify a custom config file path using the `OPENCODE_CONFIG` environment variable.
>
> ```bash
> export OPENCODE_CONFIG=/path/to/my/custom-config.json
> opencode run "Hello world"
> ```
>
> Custom config is loaded between global and project configs in the precedence order.
>
> ---
>
> ### Custom directory
>
> Specify a custom config directory using the `OPENCODE_CONFIG_DIR`
> environment variable. This directory will be searched for agents, commands,
> modes, and plugins just like the standard `.opencode` directory, and should
> follow the same structure.
>
> ```bash
> export OPENCODE_CONFIG_DIR=/path/to/my/config-directory
> opencode run "Hello world"
> ```
>
> The custom directory is loaded after the global config and `.opencode` directories, so it **can override** their settings.
>
> ---
>
> ### Managed settings
>
> Organizations can enforce configuration that users cannot override. Managed settings are loaded at the highest priority tier.
>
> #### File-based
>
> Drop an `opencode.json` or `opencode.jsonc` file in the system managed config directory:
>
> | Platform | Path                                     |
> | -------- | ---------------------------------------- |
> | macOS    | `/Library/Application Support/opencode/` |
> | Linux    | `/etc/opencode/`                         |
> | Windows  | `%ProgramData%\opencode`                 |
>
> These directories require admin/root access to write, so users cannot modify them.
>
> #### macOS managed preferences
>
> On macOS, OpenCode reads managed preferences from the `ai.opencode.managed` preference domain. Deploy a `.mobileconfig` via MDM (Jamf, Kandji, FleetDM) and the settings are enforced automatically.
>
> OpenCode checks these paths:
>
> 1. `/Library/Managed Preferences/<user>/ai.opencode.managed.plist`
> 2. `/Library/Managed Preferences/ai.opencode.managed.plist`
>
> The plist keys map directly to `opencode.json` fields. MDM metadata keys (`PayloadUUID`, `PayloadType`, etc.) are stripped automatically.
>
> **Creating a `.mobileconfig`**
>
> Use the `ai.opencode.managed` PayloadType. The OpenCode config keys go directly in the payload dict:
>
> ```xml
> <?xml version="1.0" encoding="UTF-8"?>
> <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
>   "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
> <plist version="1.0">
> <dict>
>   <key>PayloadContent</key>
>   <array>
>     <dict>
>       <key>PayloadType</key>
>       <string>ai.opencode.managed</string>
>       <key>PayloadIdentifier</key>
>       <string>com.example.opencode.config</string>
>       <key>PayloadUUID</key>
>       <string>GENERATE-YOUR-OWN-UUID</string>
>       <key>PayloadVersion</key>
>       <integer>1</integer>
>       <key>share</key>
>       <string>disabled</string>
>       <key>server</key>
>       <dict>
>         <key>hostname</key>
>         <string>127.0.0.1</string>
>       </dict>
>       <key>permission</key>
>       <dict>
>         <key>*</key>
>         <string>ask</string>
>         <key>bash</key>
>         <dict>
>           <key>*</key>
>           <string>ask</string>
>           <key>rm -rf *</key>
>           <string>deny</string>
>         </dict>
>       </dict>
>     </dict>
>   </array>
>   <key>PayloadType</key>
>   <string>Configuration</string>
>   <key>PayloadIdentifier</key>
>   <string>com.example.opencode</string>
>   <key>PayloadUUID</key>
>   <string>GENERATE-YOUR-OWN-UUID</string>
>   <key>PayloadVersion</key>
>   <integer>1</integer>
> </dict>
> </plist>
> ```
>
> Generate unique UUIDs with `uuidgen`. Customize the settings to match your organization's requirements.
>
> **Deploying via MDM**
>
> - **Jamf Pro:** Computers > Configuration Profiles > Upload > scope to target devices or smart groups
> - **FleetDM:** Add the `.mobileconfig` to your gitops repo under `mdm.macos_settings.custom_settings` and run `fleetctl apply`
>
> **Verifying on a device**
>
> Double-click the `.mobileconfig` to install locally for testing (shows in System Settings > Privacy & Security > Profiles), then run:
>
> ```bash
> opencode debug config
> ```
>
> All managed preference keys appear in the resolved config and cannot be overridden by user or project configuration.
>
> ---
>
>
> 6. **Inline config** (`OPENCODE_CONFIG_CONTENT` env var) - runtime overrides
> 7. **Managed config files** (`/Library/Application Support/opencode/` on macOS) - admin-controlled
> 8. **macOS managed preferences** (`.mobileconfig` via MDM) - highest priority, not user-overridable
>
> This means project configs can override global defaults, and global configs can override remote organizational defaults. Managed settings override everything.
>
> :::note
> The `.opencode` and `~/.config/opencode` directories use **plural names** for subdirectories: `agents/`, `commands/`, `modes/`, `plugins/`, `skills/`, `tools/`, and `themes/`. Singular names (e.g., `agent/`) are also supported for backwards compatibility.
> :::
>
> ---
>
> ### Remote
>
> Organizations can provide default configuration via the `.well-known/opencode` endpoint. This is fetched automatically when you authenticate with a provider that supports it.
>
> Remote config is loaded first, serving as the base layer. All other config sources (global, project) can override these defaults.
>
> For example, if your organization provides MCP servers that are disabled by default:
>
> ```json title="Remote config from .well-known/opencode"
> {
>   "mcp": {
>     "jira": {
>       "type": "remote",
>       "url": "https://jira.example.com/mcp",
>       "enabled": false
>     }
>   }
> }
> ```
>
> You can enable specific servers in your local config:
>
> ```json title="opencode.json"
> {
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
> ---
>
> ### Global
>
> Place your global OpenCode config in `~/.config/opencode/opencode.json`. Use global config for user-wide server/runtime preferences like providers, models, and permissions.
>
> For TUI-specific settings, use `~/.config/opencode/tui.json`.
>
> Global config overrides remote organizational defaults.
>
> ---
>
> ### Per project
>
> Add `opencode.json` in your project root. Project config has the highest precedence among standard config files - it overrides both global and remote configs.
>
> For project-specific TUI settings, add `tui.json` alongside it.
>
> :::tip
> Place project specific config in the root of your project.
> :::
>
> When OpenCode starts up, it first looks for a config file in the current directory, then traverses up to the nearest Git directory.
>
> This is also safe to be checked into Git and uses the same schema as the global one.
>
> ---
>
> ### Custom path
>
> Specify a custom config file path using the `OPENCODE_CONFIG` environment variable.
>
> ```bash
> export OPENCODE_CONFIG=/path/to/my/custom-config.json
> opencode run "Hello world"
> ```
>
> Custom config is loaded between global and project configs in the precedence order.
>
> ---
>
> ### Custom directory
>
> Specify a custom config directory using the `OPENCODE_CONFIG_DIR`
> environment variable. This directory will be searched for agents, commands,
> modes, and plugins just like the standard `.opencode` directory, and should
> follow the same structure.
>
> ```bash
> export OPENCODE_CONFIG_DIR=/path/to/my/config-directory
> opencode run "Hello world"
> ```
>
> The custom directory is loaded after the global config and `.opencode` directories, so it **can override** their settings.
>
> ---
>
> ### Managed settings
>
> Organizations can enforce configuration that users cannot override. Managed settings are loaded at the highest priority tier.
>
> #### File-based
>
> Drop an `opencode.json` or `opencode.jsonc` file in the system managed config directory:
>
> | Platform | Path                                     |
> | -------- | ---------------------------------------- |
> | macOS    | `/Library/Application Support/opencode/` |
> | Linux    | `/etc/opencode/`                         |
> | Windows  | `%ProgramData%\opencode`                 |
>
> These directories require admin/root access to write, so users cannot modify them.
>
> #### macOS managed preferences
>
> On macOS, OpenCode reads managed preferences from the `ai.opencode.managed` preference domain. Deploy a `.mobileconfig` via MDM (Jamf, Kandji, FleetDM) and the settings are enforced automatically.
>
> OpenCode checks these paths:
>
> 1. `/Library/Managed Preferences/<user>/ai.opencode.managed.plist`
> 2. `/Library/Managed Preferences/ai.opencode.managed.plist`
>
> The plist keys map directly to `opencode.json` fields. MDM metadata keys (`PayloadUUID`, `PayloadType`, etc.) are stripped automatically.
>
> **Creating a `.mobileconfig`**
>
> Use the `ai.opencode.managed` PayloadType. The OpenCode config keys go directly in the payload dict:
>
> ```xml
> <?xml version="1.0" encoding="UTF-8"?>
> <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
>   "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
> <plist version="1.0">
> <dict>
>   <key>PayloadContent</key>
>   <array>
>     <dict>
>       <key>PayloadType</key>
>       <string>ai.opencode.managed</string>
>       <key>PayloadIdentifier</key>
>       <string>com.example.opencode.config</string>
>       <key>PayloadUUID</key>
>       <string>GENERATE-YOUR-OWN-UUID</string>
>       <key>PayloadVersion</key>
>       <integer>1</integer>
>       <key>share</key>
>       <string>disabled</string>
>       <key>server</key>
>       <dict>
>         <key>hostname</key>
>         <string>127.0.0.1</string>
>       </dict>
>       <key>permission</key>
>       <dict>
>         <key>*</key>
>         <string>ask</string>
>         <key>bash</key>
>         <dict>
>           <key>*</key>
>           <string>ask</string>
>           <key>rm -rf *</key>
>           <string>deny</string>
>         </dict>
>       </dict>
>     </dict>
>   </array>
>   <key>PayloadType</key>
>   <string>Configuration</string>
>   <key>PayloadIdentifier</key>
>   <string>com.example.opencode</string>
>   <key>PayloadUUID</key>
>   <string>GENERATE-YOUR-OWN-UUID</string>
>   <key>PayloadVersion</key>
>   <integer>1</integer>
> </dict>
> </plist>
> ```
>
> Generate unique UUIDs with `uuidgen`. Customize the settings to match your organization's requirements.
>
> **Deploying via MDM**
>
> - **Jamf Pro:** Computers > Configuration Profiles > Upload > scope to target devices or smart groups
> - **FleetDM:** Add the `.mobileconfig` to your gitops repo under `mdm.macos_settings.custom_settings` and run `fleetctl apply`
>
> **Verifying on a device**
>
> Double-click the `.mobileconfig` to install locally for testing (shows in System Settings > Privacy & Security > Profiles), then run:
>
> ```bash
> opencode debug config
> ```
>
> All managed preference keys appear in the resolved config and cannot be overridden by user or project configuration.
>
> ---
>
>
> Specify a custom config file path using the `OPENCODE_CONFIG` environment variable.
>
> ```bash
> export OPENCODE_CONFIG=/path/to/my/custom-config.json
> opencode run "Hello world"
> ```
>
> Custom config is loaded between global and project configs in the precedence order.
>
> ---
>
> ### Custom directory
>
> Specify a custom config directory using the `OPENCODE_CONFIG_DIR`
> environment variable. This directory will be searched for agents, commands,
> modes, and plugins just like the standard `.opencode` directory, and should
> follow the same structure.
>
> ```bash
> export OPENCODE_CONFIG_DIR=/path/to/my/config-directory
> opencode run "Hello world"
> ```
>
> The custom directory is loaded after the global config and `.opencode` directories, so it **can override** their settings.
>
> ---
>
> ### Managed settings
>
> Organizations can enforce configuration that users cannot override. Managed settings are loaded at the highest priority tier.
>
> #### File-based
>
> Drop an `opencode.json` or `opencode.jsonc` file in the system managed config directory:
>
> | Platform | Path                                     |
> | -------- | ---------------------------------------- |
> | macOS    | `/Library/Application Support/opencode/` |
> | Linux    | `/etc/opencode/`                         |
> | Windows  | `%ProgramData%\opencode`                 |
>
> These directories require admin/root access to write, so users cannot modify them.
>
> #### macOS managed preferences
>
> On macOS, OpenCode reads managed preferences from the `ai.opencode.managed` preference domain. Deploy a `.mobileconfig` via MDM (Jamf, Kandji, FleetDM) and the settings are enforced automatically.
>
> OpenCode checks these paths:
>
> 1. `/Library/Managed Preferences/<user>/ai.opencode.managed.plist`
> 2. `/Library/Managed Preferences/ai.opencode.managed.plist`
>
> The plist keys map directly to `opencode.json` fields. MDM metadata keys (`PayloadUUID`, `PayloadType`, etc.) are stripped automatically.
>
> **Creating a `.mobileconfig`**
>
> Use the `ai.opencode.managed` PayloadType. The OpenCode config keys go directly in the payload dict:
>
> ```xml
> <?xml version="1.0" encoding="UTF-8"?>
> <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
>   "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
> <plist version="1.0">
> <dict>
>   <key>PayloadContent</key>
>   <array>
>     <dict>
>       <key>PayloadType</key>
>       <string>ai.opencode.managed</string>
>       <key>PayloadIdentifier</key>
>       <string>com.example.opencode.config</string>
>       <key>PayloadUUID</key>
>       <string>GENERATE-YOUR-OWN-UUID</string>
>       <key>PayloadVersion</key>
>       <integer>1</integer>
>       <key>share</key>
>       <string>disabled</string>
>       <key>server</key>
>       <dict>
>         <key>hostname</key>
>         <string>127.0.0.1</string>
>       </dict>
>       <key>permission</key>
>       <dict>
>         <key>*</key>
>         <string>ask</string>
>         <key>bash</key>
>         <dict>
>           <key>*</key>
>           <string>ask</string>
>           <key>rm -rf *</key>
>           <string>deny</string>
>         </dict>
>       </dict>
>     </dict>
>   </array>
>   <key>PayloadType</key>
>   <string>Configuration</string>
>   <key>PayloadIdentifier</key>
>   <string>com.example.opencode</string>
>   <key>PayloadUUID</key>
>   <string>GENERATE-YOUR-OWN-UUID</string>
>   <key>PayloadVersion</key>
>   <integer>1</integer>
> </dict>
> </plist>
> ```
>
> Generate unique UUIDs with `uuidgen`. Customize the settings to match your organization's requirements.
>
> **Deploying via MDM**
>
> - **Jamf Pro:** Computers > Configuration Profiles > Upload > scope to target devices or smart groups
> - **FleetDM:** Add the `.mobileconfig` to your gitops repo under `mdm.macos_settings.custom_settings` and run `fleetctl apply`
>
> **Verifying on a device**
>
> Double-click the `.mobileconfig` to install locally for testing (shows in System Settings > Privacy & Security > Profiles), then run:
>
> ```bash
> opencode debug config
> ```
>
> All managed preference keys appear in the resolved config and cannot be overridden by user or project configuration.
>
> ---
>
>
> export OPENCODE_CONFIG=/path/to/my/custom-config.json
> opencode run "Hello world"
> ```
>
> Custom config is loaded between global and project configs in the precedence order.
>
> ---
>
> ### Custom directory
>
> Specify a custom config directory using the `OPENCODE_CONFIG_DIR`
> environment variable. This directory will be searched for agents, commands,
> modes, and plugins just like the standard `.opencode` directory, and should
> follow the same structure.
>
> ```bash
> export OPENCODE_CONFIG_DIR=/path/to/my/config-directory
> opencode run "Hello world"
> ```
>
> The custom directory is loaded after the global config and `.opencode` directories, so it **can override** their settings.
>
> ---
>
> ### Managed settings
>
> Organizations can enforce configuration that users cannot override. Managed settings are loaded at the highest priority tier.
>
> #### File-based
>
> Drop an `opencode.json` or `opencode.jsonc` file in the system managed config directory:
>
> | Platform | Path                                     |
> | -------- | ---------------------------------------- |
> | macOS    | `/Library/Application Support/opencode/` |
> | Linux    | `/etc/opencode/`                         |
> | Windows  | `%ProgramData%\opencode`                 |
>
> These directories require admin/root access to write, so users cannot modify them.
>
> #### macOS managed preferences
>
> On macOS, OpenCode reads managed preferences from the `ai.opencode.managed` preference domain. Deploy a `.mobileconfig` via MDM (Jamf, Kandji, FleetDM) and the settings are enforced automatically.
>
> OpenCode checks these paths:
>
> 1. `/Library/Managed Preferences/<user>/ai.opencode.managed.plist`
> 2. `/Library/Managed Preferences/ai.opencode.managed.plist`
>
> The plist keys map directly to `opencode.json` fields. MDM metadata keys (`PayloadUUID`, `PayloadType`, etc.) are stripped automatically.
>
> **Creating a `.mobileconfig`**
>
> Use the `ai.opencode.managed` PayloadType. The OpenCode config keys go directly in the payload dict:
>
> ```xml
> <?xml version="1.0" encoding="UTF-8"?>
> <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
>   "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
> <plist version="1.0">
> <dict>
>   <key>PayloadContent</key>
>   <array>
>     <dict>
>       <key>PayloadType</key>
>       <string>ai.opencode.managed</string>
>       <key>PayloadIdentifier</key>
>       <string>com.example.opencode.config</string>
>       <key>PayloadUUID</key>
>       <string>GENERATE-YOUR-OWN-UUID</string>
>       <key>PayloadVersion</key>
>       <integer>1</integer>
>       <key>share</key>
>       <string>disabled</string>
>       <key>server</key>
>       <dict>
>         <key>hostname</key>
>         <string>127.0.0.1</string>
>       </dict>
>       <key>permission</key>
>       <dict>
>         <key>*</key>
>         <string>ask</string>
>         <key>bash</key>
>         <dict>
>           <key>*</key>
>           <string>ask</string>
>           <key>rm -rf *</key>
>           <string>deny</string>
>         </dict>
>       </dict>
>     </dict>
>   </array>
>   <key>PayloadType</key>
>   <string>Configuration</string>
>   <key>PayloadIdentifier</key>
>   <string>com.example.opencode</string>
>   <key>PayloadUUID</key>
>   <string>GENERATE-YOUR-OWN-UUID</string>
>   <key>PayloadVersion</key>
>   <integer>1</integer>
> </dict>
> </plist>
> ```
>
> Generate unique UUIDs with `uuidgen`. Customize the settings to match your organization's requirements.
>
> **Deploying via MDM**
>
> - **Jamf Pro:** Computers > Configuration Profiles > Upload > scope to target devices or smart groups
> - **FleetDM:** Add the `.mobileconfig` to your gitops repo under `mdm.macos_settings.custom_settings` and run `fleetctl apply`
>
> **Verifying on a device**
>
> Double-click the `.mobileconfig` to install locally for testing (shows in System Settings > Privacy & Security > Profiles), then run:
>
> ```bash
> opencode debug config
> ```
>
> All managed preference keys appear in the resolved config and cannot be overridden by user or project configuration.
>
> ---
>
>
> Specify a custom config directory using the `OPENCODE_CONFIG_DIR`
> environment variable. This directory will be searched for agents, commands,
> modes, and plugins just like the standard `.opencode` directory, and should
> follow the same structure.
>
> ```bash
> export OPENCODE_CONFIG_DIR=/path/to/my/config-directory
> opencode run "Hello world"
> ```
>
> The custom directory is loaded after the global config and `.opencode` directories, so it **can override** their settings.
>
> ---
>
> ### Managed settings
>
> Organizations can enforce configuration that users cannot override. Managed settings are loaded at the highest priority tier.
>
> #### File-based
>
> Drop an `opencode.json` or `opencode.jsonc` file in the system managed config directory:
>
> | Platform | Path                                     |
> | -------- | ---------------------------------------- |
> | macOS    | `/Library/Application Support/opencode/` |
> | Linux    | `/etc/opencode/`                         |
> | Windows  | `%ProgramData%\opencode`                 |
>
> These directories require admin/root access to write, so users cannot modify them.
>
> #### macOS managed preferences
>
> On macOS, OpenCode reads managed preferences from the `ai.opencode.managed` preference domain. Deploy a `.mobileconfig` via MDM (Jamf, Kandji, FleetDM) and the settings are enforced automatically.
>
> OpenCode checks these paths:
>
> 1. `/Library/Managed Preferences/<user>/ai.opencode.managed.plist`
> 2. `/Library/Managed Preferences/ai.opencode.managed.plist`
>
> The plist keys map directly to `opencode.json` fields. MDM metadata keys (`PayloadUUID`, `PayloadType`, etc.) are stripped automatically.
>
> **Creating a `.mobileconfig`**
>
> Use the `ai.opencode.managed` PayloadType. The OpenCode config keys go directly in the payload dict:
>
> ```xml
> <?xml version="1.0" encoding="UTF-8"?>
> <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
>   "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
> <plist version="1.0">
> <dict>
>   <key>PayloadContent</key>
>   <array>
>     <dict>
>       <key>PayloadType</key>
>       <string>ai.opencode.managed</string>
>       <key>PayloadIdentifier</key>
>       <string>com.example.opencode.config</string>
>       <key>PayloadUUID</key>
>       <string>GENERATE-YOUR-OWN-UUID</string>
>       <key>PayloadVersion</key>
>       <integer>1</integer>
>       <key>share</key>
>       <string>disabled</string>
>       <key>server</key>
>       <dict>
>         <key>hostname</key>
>         <string>127.0.0.1</string>
>       </dict>
>       <key>permission</key>
>       <dict>
>         <key>*</key>
>         <string>ask</string>
>         <key>bash</key>
>         <dict>
>           <key>*</key>
>           <string>ask</string>
>           <key>rm -rf *</key>
>           <string>deny</string>
>         </dict>
>       </dict>
>     </dict>
>   </array>
>   <key>PayloadType</key>
>   <string>Configuration</string>
>   <key>PayloadIdentifier</key>
>   <string>com.example.opencode</string>
>   <key>PayloadUUID</key>
>   <string>GENERATE-YOUR-OWN-UUID</string>
>   <key>PayloadVersion</key>
>   <integer>1</integer>
> </dict>
> </plist>
> ```
>
> Generate unique UUIDs with `uuidgen`. Customize the settings to match your organization's requirements.
>
> **Deploying via MDM**
>
> - **Jamf Pro:** Computers > Configuration Profiles > Upload > scope to target devices or smart groups
> - **FleetDM:** Add the `.mobileconfig` to your gitops repo under `mdm.macos_settings.custom_settings` and run `fleetctl apply`
>
> **Verifying on a device**
>
> Double-click the `.mobileconfig` to install locally for testing (shows in System Settings > Privacy & Security > Profiles), then run:
>
> ```bash
> opencode debug config
> ```
>
> All managed preference keys appear in the resolved config and cannot be overridden by user or project configuration.
>
> ---
>
>
> export OPENCODE_CONFIG_DIR=/path/to/my/config-directory
> opencode run "Hello world"
> ```
>
> The custom directory is loaded after the global config and `.opencode` directories, so it **can override** their settings.
>
> ---
>
> ### Managed settings
>
> Organizations can enforce configuration that users cannot override. Managed settings are loaded at the highest priority tier.
>
> #### File-based
>
> Drop an `opencode.json` or `opencode.jsonc` file in the system managed config directory:
>
> | Platform | Path                                     |
> | -------- | ---------------------------------------- |
> | macOS    | `/Library/Application Support/opencode/` |
> | Linux    | `/etc/opencode/`                         |
> | Windows  | `%ProgramData%\opencode`                 |
>
> These directories require admin/root access to write, so users cannot modify them.
>
> #### macOS managed preferences
>
> On macOS, OpenCode reads managed preferences from the `ai.opencode.managed` preference domain. Deploy a `.mobileconfig` via MDM (Jamf, Kandji, FleetDM) and the settings are enforced automatically.
>
> OpenCode checks these paths:
>
> 1. `/Library/Managed Preferences/<user>/ai.opencode.managed.plist`
> 2. `/Library/Managed Preferences/ai.opencode.managed.plist`
>
> The plist keys map directly to `opencode.json` fields. MDM metadata keys (`PayloadUUID`, `PayloadType`, etc.) are stripped automatically.
>
> **Creating a `.mobileconfig`**
>
> Use the `ai.opencode.managed` PayloadType. The OpenCode config keys go directly in the payload dict:
>
> ```xml
> <?xml version="1.0" encoding="UTF-8"?>
> <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
>   "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
> <plist version="1.0">
> <dict>
>   <key>PayloadContent</key>
>   <array>
>     <dict>
>       <key>PayloadType</key>
>       <string>ai.opencode.managed</string>
>       <key>PayloadIdentifier</key>
>       <string>com.example.opencode.config</string>
>       <key>PayloadUUID</key>
>       <string>GENERATE-YOUR-OWN-UUID</string>
>       <key>PayloadVersion</key>
>       <integer>1</integer>
>       <key>share</key>
>       <string>disabled</string>
>       <key>server</key>
>       <dict>
>         <key>hostname</key>
>         <string>127.0.0.1</string>
>       </dict>
>       <key>permission</key>
>       <dict>
>         <key>*</key>
>         <string>ask</string>
>         <key>bash</key>
>         <dict>
>           <key>*</key>
>           <string>ask</string>
>           <key>rm -rf *</key>
>           <string>deny</string>
>         </dict>
>       </dict>
>     </dict>
>   </array>
>   <key>PayloadType</key>
>   <string>Configuration</string>
>   <key>PayloadIdentifier</key>
>   <string>com.example.opencode</string>
>   <key>PayloadUUID</key>
>   <string>GENERATE-YOUR-OWN-UUID</string>
>   <key>PayloadVersion</key>
>   <integer>1</integer>
> </dict>
> </plist>
> ```
>
> Generate unique UUIDs with `uuidgen`. Customize the settings to match your organization's requirements.
>
> **Deploying via MDM**
>
> - **Jamf Pro:** Computers > Configuration Profiles > Upload > scope to target devices or smart groups
> - **FleetDM:** Add the `.mobileconfig` to your gitops repo under `mdm.macos_settings.custom_settings` and run `fleetctl apply`
>
> **Verifying on a device**
>
> Double-click the `.mobileconfig` to install locally for testing (shows in System Settings > Privacy & Security > Profiles), then run:
>
> ```bash
> opencode debug config
> ```
>
> All managed preference keys appear in the resolved config and cannot be overridden by user or project configuration.
>
> ---
>
>
> Use `OPENCODE_TUI_CONFIG` to point to a custom TUI config file.
>
> When `cursor.style` is `"default"`, the terminal default cursor is restored, so `cursor.blinking` has no effect.
>
> Set `attention.enabled` to turn on TUI desktop notifications and sounds. See [TUI attention](/docs/tui#attention).
>
> Legacy `theme`, `keybinds`, and `tui` keys in `opencode.json` are deprecated and automatically migrated when possible.
>
> ---
>
> ### Server
>
> You can configure server settings for the `opencode serve` and `opencode web` commands through the `server` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "server": {
>     "port": 4096,
>     "hostname": "0.0.0.0",
>     "mdns": true,
>     "mdnsDomain": "myproject.local",
>     "cors": ["http://localhost:5173"]
>   }
> }
> ```
>
> Available options:
>
> - `port` - Port to listen on.
> - `hostname` - Hostname to listen on. When `mdns` is enabled and no hostname is set, defaults to `0.0.0.0`.
> - `mdns` - Enable mDNS service discovery. This allows other devices on the network to discover your OpenCode server.
> - `mdnsDomain` - Custom domain name for mDNS service. Defaults to `opencode.local`. Useful for running multiple instances on the same network.
> - `cors` - Additional origins to allow for CORS when using the HTTP server from a browser-based client. Values must be full origins (scheme + host + optional port), eg `https://app.example.com`.
>
> [Learn more about the server here](/docs/server).
>
> ---
>
> ### Shell
>
> You can configure the shell used for the interactive terminal using the `shell` option. Compatible shells are also used for agent tool calls.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "shell": "pwsh"
> }
> ```
>
> If not specified, OpenCode will automatically discover and use a sensible default based on your operating system (e.g. `pwsh` or `cmd.exe` on Windows, `/bin/zsh` or `/bin/bash` on macOS/Linux). You can provide an absolute path or a short name.
>
> ---
>
> ### Tools
>
> You can manage the tools an LLM can use through the `tools` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "tools": {
>     "write": false,
>     "bash": false
>   }
> }
> ```
>
> [Learn more about tools here](/docs/tools).
>
> ---
>
> ### Models
>
> You can configure the providers and models you want to use in your OpenCode config through the `provider`, `model` and `small_model` options.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "provider": {},
>   "model": "anthropic/claude-sonnet-4-5",
>   "small_model": "anthropic/claude-haiku-4-5"
> }
> ```
>
> The `small_model` option configures a separate model for lightweight tasks like title generation. By default, OpenCode tries to use a cheaper model if one is available from your provider, otherwise it falls back to your main model.
>
> Provider options can include `timeout`, `chunkTimeout`, and `setCacheKey`:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "provider": {
>     "anthropic": {
>       "options": {
>         "timeout": 600000,
>         "chunkTimeout": 30000,
>         "setCacheKey": true
>       }
>     }
>   }
> }
> ```
>
> - `timeout` - Request timeout in milliseconds (default: 300000). Set to `false` to disable.
> - `chunkTimeout` - Timeout in milliseconds between streamed response chunks. If no chunk arrives in time, the request is aborted.
> - `setCacheKey` - Ensure a cache key is always set for designated provider.
>
> You can also configure [local models](/docs/models#local). [Learn more](/docs/models).
>
> ---
>
> ### Policies
>
> Use the `experimental.policies` option to allow or deny OpenCode actions on configured resources. Currently, policies can control which providers OpenCode may use.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "experimental": {
>     "policies": [
>       {
>         "effect": "deny",
>         "action": "provider.use",
>         "resource": "openai"
>       }
>     ]
>   }
> }
> ```
>
> [Learn more about policies here](/docs/policies).
>
> ---
>
> ### Image attachments
>
> OpenCode normalizes image attachments before sending them to the model. By default, images are resized when they exceed `2000x2000` pixels or `5242880` base64 bytes.
>
> Configure image attachment limits with the `attachment.image` option:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "attachment": {
>     "image": {
>       "auto_resize": true,
>       "max_width": 2000,
>       "max_height": 2000,
>       "max_base64_bytes": 5242880
>     }
>   }
> }
> ```
>
> - `auto_resize` - Resize images that exceed the configured limits before provider requests. Set to `false` to reject oversized images instead.
> - `max_width` - Maximum image width in pixels before resizing or rejection.
> - `max_height` - Maximum image height in pixels before resizing or rejection.
> - `max_base64_bytes` - Maximum encoded image payload size. This is the base64 payload size, not the original file size.
>
> If an image still cannot fit after resizing, OpenCode omits oversized tool-result images or fails oversized user-provided images with an image size error.
>
> ---
>
> #### Provider-Specific Options
>
> Some providers support additional configuration options beyond the generic `timeout` and `apiKey` settings.
>
> ##### Amazon Bedrock
>
> Amazon Bedrock supports AWS-specific configuration:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "provider": {
>     "amazon-bedrock": {
>       "options": {
>         "region": "us-east-1",
>         "profile": "my-aws-profile",
>         "endpoint": "https://bedrock-runtime.us-east-1.vpce-xxxxx.amazonaws.com"
>       }
>     }
>   }
> }
> ```
>
> - `region` - AWS region for Bedrock (defaults to `AWS_REGION` env var or `us-east-1`)
> - `profile` - AWS named profile from `~/.aws/credentials` (defaults to `AWS_PROFILE` env var)
> - `endpoint` - Custom endpoint URL for VPC endpoints. This is an alias for the generic `baseURL` option using AWS-specific terminology. If both are specified, `endpoint` takes precedence.
>
> :::note
> Bearer tokens (`AWS_BEARER_TOKEN_BEDROCK` or `/connect`) take precedence over profile-based authentication. See [authentication precedence](/docs/providers#authentication-precedence) for details.
> :::
>
> [Learn more about Amazon Bedrock configuration](/docs/providers#amazon-bedrock).
>
> ---
>
> ### Themes
>
> Set your UI theme in `tui.json`.
>
> ```json title="tui.json"
> {
>   "$schema": "https://opencode.ai/tui.json",
>   "theme": "tokyonight"
> }
> ```
>
> [Learn more here](/docs/themes).
>
> ---
>
> ### Agents
>
> You can configure specialized agents for specific tasks through the `agent` option.
>
> ```jsonc title="opencode.jsonc"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "agent": {
>     "code-reviewer": {
>       "description": "Reviews code for best practices and potential issues",
>       "model": "anthropic/claude-sonnet-4-5",
>       "prompt": "You are a code reviewer. Focus on security, performance, and maintainability.",
>       "tools": {
>         // Disable file modification tools for review-only agent
>         "write": false,
>         "edit": false,
>       },
>     },
>   },
> }
> ```
>
> You can also define agents using markdown files in `~/.config/opencode/agents/` or `.opencode/agents/`. [Learn more here](/docs/agents).
>
> ---
>
> ### Default agent
>
> You can set the default agent using the `default_agent` option. This determines which agent is used when none is explicitly specified.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "default_agent": "plan"
> }
> ```
>
> The default agent must be a primary agent (not a subagent). This can be a built-in agent like `"build"` or `"plan"`, or a [custom agent](/docs/agents) you've defined. If the specified agent doesn't exist or is a subagent, OpenCode will fall back to `"build"` with a warning.
>
> This setting applies across all interfaces: TUI, CLI (`opencode run`), desktop app, and GitHub Action.
>
> ---
>
> ### Subagent depth
>
> You can control how deeply subagents can invoke other subagents using the `subagent_depth` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "subagent_depth": 2
> }
> ```
>
> The default is `1`, which allows primary agents to launch subagents but prevents those subagents from launching additional subagents. Set it to `2` to allow one additional level of nested subagents, or `0` to prevent all subagent launches.
>
> ---
>
> ### Sharing
>
> You can configure the [share](/docs/share) feature through the `share` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "share": "manual"
> }
> ```
>
> This takes:
>
> - `"manual"` - Allow manual sharing via commands (default)
> - `"auto"` - Automatically share new conversations
> - `"disabled"` - Disable sharing entirely
>
> By default, sharing is set to manual mode where you need to explicitly share conversations using the `/share` command.
>
> ---
>
> ### Commands
>
> You can configure custom commands for repetitive tasks through the `command` option.
>
> ```jsonc title="opencode.jsonc"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "command": {
>     "test": {
>       "template": "Run the full test suite with coverage report and show any failures.\nFocus on the failing tests and suggest fixes.",
>       "description": "Run tests with coverage",
>       "agent": "build",
>       "model": "anthropic/claude-haiku-4-5",
>     },
>     "component": {
>       "template": "Create a new React component named $ARGUMENTS with TypeScript support.\nInclude proper typing and basic structure.",
>       "description": "Create a new component",
>     },
>   },
> }
> ```
>
> You can also define commands using markdown files in `~/.config/opencode/commands/` or `.opencode/commands/`. [Learn more here](/docs/commands).
>
> ---
>
> ### Keybinds
>
> Customize TUI keyboard shortcuts in `tui.json` with `keybinds`.
>
> ```json title="tui.json"
> {
>   "$schema": "https://opencode.ai/tui.json",
>   "keybinds": {
>     "command_list": "ctrl+p"
>   }
> }
> ```
>
> `keybinds` is merged with built-in defaults, so you only need to configure the shortcuts you want to change.
>
> [Learn more here](/docs/keybinds).
>
> ---
>
> ### Snapshot
>
> OpenCode uses snapshots to track file changes during agent operations, enabling you to undo and revert changes within a session. Snapshots are enabled by default.
>
> For large repositories or projects with many submodules, the snapshot system can cause slow indexing and significant disk usage as it tracks all changes using an internal git repository. You can disable snapshots using the `snapshot` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "snapshot": false
> }
> ```
>
> Note that disabling snapshots means changes made by the agent cannot be rolled back through the UI.
>
> ---
>
> ### Autoupdate
>
> OpenCode will automatically download any new updates when it starts up. You can disable this with the `autoupdate` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "autoupdate": false
> }
> ```
>
> If you don't want updates but want to be notified when a new version is available, set `autoupdate` to `"notify"`.
> Notice that this only works if it was not installed using a package manager such as Homebrew.
>
> ---
>
> ### Formatters
>
> You can enable and configure code formatters through the `formatter` option. Omit it to keep formatters disabled.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "formatter": true
> }
> ```
>
> Use an object to keep built-ins enabled while configuring overrides or custom formatters.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "formatter": {
>     "prettier": {
>       "disabled": true
>     },
>     "custom-prettier": {
>       "command": ["npx", "prettier", "--write", "$FILE"],
>       "environment": {
>         "NODE_ENV": "development"
>       },
>       "extensions": [".js", ".ts", ".jsx", ".tsx"]
>     }
>   }
> }
> ```
>
> [Learn more about formatters here](/docs/formatters).
>
> ---
>
> ### LSP Servers
>
> You can enable and configure LSP servers through the `lsp` option. Omit it to keep LSP disabled.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "lsp": true
> }
> ```
>
> Use an object to keep built-ins enabled while configuring overrides or custom LSP servers.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "lsp": {
>     "typescript": {
>       "disabled": true
>     }
>   }
> }
> ```
>
> [Learn more about LSP servers here](/docs/lsp).
>
> ---
>
> ### Permissions
>
> By default, opencode **allows all operations** without requiring explicit approval. You can change this using the `permission` option.
>
> For example, to ensure that the `edit` and `bash` tools require user approval:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "edit": "ask",
>     "bash": "ask"
>   }
> }
> ```
>
> [Learn more about permissions here](/docs/permissions).
>
> ---
>
> ### Compaction
>
> You can control context compaction behavior through the `compaction` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "compaction": {
>     "auto": true,
>     "prune": false,
>     "reserved": 10000
>   }
> }
> ```
>
> - `auto` - Automatically compact the session when context is full (default: `true`).
> - `prune` - Remove old tool outputs to save tokens (default: `false`). Set to `true` to enable pruning.
> - `reserved` - Token buffer for compaction. Leaves enough window to avoid overflow during compaction.
>
> ---
>
> ### Watcher
>
> You can configure file watcher ignore patterns through the `watcher` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "watcher": {
>     "ignore": ["node_modules/**", "dist/**", ".git/**"]
>   }
> }
> ```
>
> Patterns follow glob syntax. Use this to exclude noisy directories from file watching.
>
> ---
>
> ### MCP servers
>
> You can configure MCP servers you want to use through the `mcp` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "mcp": {}
> }
> ```
>
> [Learn more here](/docs/mcp-servers).
>
> ---
>
> ### Plugins
>
> [Plugins](/docs/plugins) extend OpenCode with custom tools, hooks, and integrations.
>
> Place plugin files in `.opencode/plugins/` or `~/.config/opencode/plugins/`. You can also load plugins from npm through the `plugin` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "plugin": ["opencode-helicone-session", "@my-org/custom-plugin"]
> }
> ```
>
> [Learn more here](/docs/plugins).
>
> ---
>
> ### Instructions
>
> You can configure the instructions for the model you're using through the `instructions` option.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "instructions": ["CONTRIBUTING.md", "docs/guidelines.md", ".cursor/rules/*.md"]
> }
> ```
>
> This takes an array of paths and glob patterns to instruction files. [Learn more
> about rules here](/docs/rules).
>
> ---
>
> ### Disabled providers
>
> You can disable providers that are loaded automatically through the `disabled_providers` option. This is useful when you want to prevent certain providers from being loaded even if their credentials are available.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "disabled_providers": ["openai", "gemini"]
> }
> ```
>
> :::note
> The `disabled_providers` takes priority over `enabled_providers`.
> :::
>
> The `disabled_providers` option accepts an array of provider IDs. When a provider is disabled:
>
> - It won't be loaded even if environment variables are set.
> - It won't be loaded even if API keys are configured through the `/connect` command.
> - The provider's models won't appear in the model selection list.
>
> ---
>
> ### Enabled providers
>
> You can specify an allowlist of providers through the `enabled_providers` option. When set, only the specified providers will be enabled and all others will be ignored.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "enabled_providers": ["anthropic", "openai"]
> }
> ```
>
> This is useful when you want to restrict OpenCode to only use specific providers rather than disabling them one by one.
>
> :::note
> The `disabled_providers` takes priority over `enabled_providers`.
> :::
>
> If a provider appears in both `enabled_providers` and `disabled_providers`, the `disabled_providers` takes priority for backwards compatibility.
>
> ---
>
> ### Experimental
>
> The `experimental` key contains options that are under active development.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "experimental": {}
> }
> ```
>
> :::caution
> Experimental options are not stable. They may change or be removed without notice.
> :::
>
> ---
>
>
> ### Experimental
>
> The `experimental` key contains options that are under active development.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "experimental": {}
> }
> ```
>
> :::caution
> Experimental options are not stable. They may change or be removed without notice.
> :::
>
> ---
>
>
> Experimental options are not stable. They may change or be removed without notice.
> :::
>
> ---
>
>
>   "model": "{env:OPENCODE_MODEL}",
>   "provider": {
>     "anthropic": {
>       "models": {},
>       "options": {
>         "apiKey": "{env:ANTHROPIC_API_KEY}"
>       }
>     }
>   }
> }
> ```
>
> If the environment variable is not set, it will be replaced with an empty string.
>
> ---
>
> ### Files
>
> Use `{file:path/to/file}` to substitute the contents of a file:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "instructions": ["./custom-instructions.md"],
>   "provider": {
>     "openai": {
>       "options": {
>         "apiKey": "{file:~/.secrets/openai-key}"
>       }
>     }
>   }
> }
> ```
>
> File paths can be:
>
> - Relative to the config file directory
> - Or absolute paths starting with `/` or `~`
>
> These are useful for:
>
> - Keeping sensitive data like API keys in separate files.
> - Including large instruction files without cluttering your config.
> - Sharing common configuration snippets across multiple config files.
