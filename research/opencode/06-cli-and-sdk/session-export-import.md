---
primary_sources:
  - id: T1-CLI
    title: "CLI"
    url: "https://opencode.ai/docs/cli.md"
    section: "export and import"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Session export and import

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode CLI — export and import

> import { Tabs, TabItem } from "@astrojs/starlight/components"
>
> The OpenCode CLI by default starts the [TUI](/docs/tui) when run without any arguments.
>
> ```bash
> opencode
> ```
>
> But it also accepts commands as documented on this page. This allows you to interact with OpenCode programmatically.
>
> ```bash
> opencode run "Explain how closures work in JavaScript"
> ```
>
> ---
>
> ### tui
>
> Start the OpenCode terminal user interface.
>
> ```bash
> opencode [project]
> ```
>
> #### Flags
>
> | Flag                                        | Short | Description                                                             |
> | ------------------------------------------- | ----- | ----------------------------------------------------------------------- |
> | <nobr><code>{"--continue"}</code></nobr>    | `-c`  | Continue the last session                                               |
> | <nobr><code>{"--session"}</code></nobr>     | `-s`  | Session ID to continue                                                  |
> | <nobr><code>{"--fork"}</code></nobr>        |       | Fork the session when continuing (use with `--continue` or `--session`) |
> | <nobr><code>{"--prompt"}</code></nobr>      |       | Prompt to use                                                           |
> | <nobr><code>{"--model"}</code></nobr>       | `-m`  | Model to use in the form of provider/model                              |
> | <nobr><code>{"--agent"}</code></nobr>       |       | Agent to use                                                            |
> | <nobr><code>{"--auto"}</code></nobr>        |       | Auto-approve permissions that are not explicitly denied                 |
> | <nobr><code>{"--port"}</code></nobr>        |       | Port to listen on                                                       |
> | <nobr><code>{"--hostname"}</code></nobr>    |       | Hostname to listen on                                                   |
> | <nobr><code>{"--mdns"}</code></nobr>        |       | Enable mDNS discovery                                                   |
> | <nobr><code>{"--mdns-domain"}</code></nobr> |       | Custom mDNS domain name                                                 |
> | <nobr><code>{"--cors"}</code></nobr>        |       | Additional browser origin(s) to allow CORS                              |
>
> ---
>
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
