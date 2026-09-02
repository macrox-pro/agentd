---
primary_sources:
  - id: T1-INTRO
    title: "Intro"
    url: "https://opencode.ai/docs/"
    section: "Usage"
  - id: T1-COMMANDS
    title: "Commands"
    url: "https://opencode.ai/docs/commands.md"
    section: "Built-in"
  - id: T1-CLI
    title: "CLI"
    url: "https://opencode.ai/docs/cli.md"
    section: "Commands"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Built-in slash commands

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Intro — slash-related usage

> /init
> ```
>
> This will get OpenCode to analyze your project and create an `AGENTS.md` file in
> the project root.
>
> :::tip
> You should commit your project's `AGENTS.md` file to Git.
> :::
>
> This helps OpenCode understand the project structure and the coding patterns
> used.
>
> ---
>
>
>    OpenCode has a _Plan mode_ that disables its ability to make changes and
>    instead suggest _how_ it'll implement the feature.
>
>    Switch to it using the **Tab** key. You'll see an indicator for this in the lower right corner.
>
>    ```bash frame="none" title="Switch to Plan mode"
>    <TAB>
>    ```
>
>    Now let's describe what we want it to do.
>
>    ```txt frame="none"
>    When a user deletes a note, we'd like to flag it as deleted in the database.
>    Then create a screen that shows all the recently deleted notes.
>    From this screen, the user can undelete a note or permanently delete it.
>    ```
>
>    You want to give OpenCode enough details to understand what you want. It helps
>    to talk to it like you are talking to a junior developer on your team.
>
>    :::tip
>    Give OpenCode plenty of context and examples to help it understand what you
>    want.
>    :::
>
> 2. **Iterate on the plan**
>
>    Once it gives you a plan, you can give it feedback or add more details.
>
>    ```txt frame="none"
>    We'd like to design this new screen using a design I've used before.
>    [Image #1] Take a look at this image and use it as a reference.
>    ```
>
>    :::tip
>    Drag and drop images into the terminal to add them to the prompt.
>    :::
>
>    OpenCode can scan any images you give it and add them to the prompt. You can
>    do this by dragging and dropping an image into the terminal.
>
> 3. **Build the feature**
>
>    Once you feel comfortable with the plan, switch back to _Build mode_ by
>    hitting the **Tab** key again.
>
>    ```bash frame="none"
>    <TAB>
>    ```
>
>    And asking it to make the changes.
>
>    ```bash frame="none"
>    Sounds good! Go ahead and make the changes.
>    ```
>
> ---
>
> ### Make changes
>
> For more straightforward changes, you can ask OpenCode to directly build it
> without having to review the plan first.
>
> ```txt frame="none" "@packages/functions/src/settings.ts" "@packages/functions/src/notes.ts"
> We need to add authentication to the /settings route. Take a look at how this is
> handled in the /notes route in @packages/functions/src/notes.ts and implement
> the same logic in @packages/functions/src/settings.ts
> ```
>
> You want to make sure you provide a good amount of detail so OpenCode makes the right
> changes.
>
> ---
>
> ### Undo changes
>
> Let's say you ask OpenCode to make some changes.
>
> ```txt frame="none" "@packages/functions/src/api/index.ts"
> Can you refactor the function in @packages/functions/src/api/index.ts?
> ```
>
> But you realize that it is not what you wanted. You **can undo** the changes
> using the `/undo` command.
>
> ```bash frame="none"
> /undo
> ```
>
> OpenCode will now revert the changes you made and show your original message
> again.
>
> ```txt frame="none" "@packages/functions/src/api/index.ts"
> Can you refactor the function in @packages/functions/src/api/index.ts?
> ```
>
> From here you can tweak the prompt and ask OpenCode to try again.
>
> :::tip
> You can run `/undo` multiple times to undo multiple changes.
> :::
>
> Or you **can redo** the changes using the `/redo` command.
>
> ```bash frame="none"
> /redo
> ```
>
> ---
>
>
>    ```bash frame="none" title="Switch to Plan mode"
>    <TAB>
>    ```
>
>    Now let's describe what we want it to do.
>
>    ```txt frame="none"
>    When a user deletes a note, we'd like to flag it as deleted in the database.
>    Then create a screen that shows all the recently deleted notes.
>    From this screen, the user can undelete a note or permanently delete it.
>    ```
>
>    You want to give OpenCode enough details to understand what you want. It helps
>    to talk to it like you are talking to a junior developer on your team.
>
>    :::tip
>    Give OpenCode plenty of context and examples to help it understand what you
>    want.
>    :::
>
> 2. **Iterate on the plan**
>
>    Once it gives you a plan, you can give it feedback or add more details.
>
>    ```txt frame="none"
>    We'd like to design this new screen using a design I've used before.
>    [Image #1] Take a look at this image and use it as a reference.
>    ```
>
>    :::tip
>    Drag and drop images into the terminal to add them to the prompt.
>    :::
>
>    OpenCode can scan any images you give it and add them to the prompt. You can
>    do this by dragging and dropping an image into the terminal.
>
> 3. **Build the feature**
>
>    Once you feel comfortable with the plan, switch back to _Build mode_ by
>    hitting the **Tab** key again.
>
>    ```bash frame="none"
>    <TAB>
>    ```
>
>    And asking it to make the changes.
>
>    ```bash frame="none"
>    Sounds good! Go ahead and make the changes.
>    ```
>
> ---
>
> ### Make changes
>
> For more straightforward changes, you can ask OpenCode to directly build it
> without having to review the plan first.
>
> ```txt frame="none" "@packages/functions/src/settings.ts" "@packages/functions/src/notes.ts"
> We need to add authentication to the /settings route. Take a look at how this is
> handled in the /notes route in @packages/functions/src/notes.ts and implement
> the same logic in @packages/functions/src/settings.ts
> ```
>
> You want to make sure you provide a good amount of detail so OpenCode makes the right
> changes.
>
> ---
>
> ### Undo changes
>
> Let's say you ask OpenCode to make some changes.
>
> ```txt frame="none" "@packages/functions/src/api/index.ts"
> Can you refactor the function in @packages/functions/src/api/index.ts?
> ```
>
> But you realize that it is not what you wanted. You **can undo** the changes
> using the `/undo` command.
>
> ```bash frame="none"
> /undo
> ```
>
> OpenCode will now revert the changes you made and show your original message
> again.
>
> ```txt frame="none" "@packages/functions/src/api/index.ts"
> Can you refactor the function in @packages/functions/src/api/index.ts?
> ```
>
> From here you can tweak the prompt and ask OpenCode to try again.
>
> :::tip
> You can run `/undo` multiple times to undo multiple changes.
> :::
>
> Or you **can redo** the changes using the `/redo` command.
>
> ```bash frame="none"
> /redo
> ```
>
> ---
>
>
> using the `/undo` command.
>
> ```bash frame="none"
> /undo
> ```
>
> OpenCode will now revert the changes you made and show your original message
> again.
>
> ```txt frame="none" "@packages/functions/src/api/index.ts"
> Can you refactor the function in @packages/functions/src/api/index.ts?
> ```
>
> From here you can tweak the prompt and ask OpenCode to try again.
>
> :::tip
> You can run `/undo` multiple times to undo multiple changes.
> :::
>
> Or you **can redo** the changes using the `/redo` command.
>
> ```bash frame="none"
> /redo
> ```
>
> ---
>
>
> /undo
> ```
>
> OpenCode will now revert the changes you made and show your original message
> again.
>
> ```txt frame="none" "@packages/functions/src/api/index.ts"
> Can you refactor the function in @packages/functions/src/api/index.ts?
> ```
>
> From here you can tweak the prompt and ask OpenCode to try again.
>
> :::tip
> You can run `/undo` multiple times to undo multiple changes.
> :::
>
> Or you **can redo** the changes using the `/redo` command.
>
> ```bash frame="none"
> /redo
> ```
>
> ---
>
>
> You can run `/undo` multiple times to undo multiple changes.
> :::
>
> Or you **can redo** the changes using the `/redo` command.
>
> ```bash frame="none"
> /redo
> ```
>
> ---
>
>
> Or you **can redo** the changes using the `/redo` command.
>
> ```bash frame="none"
> /redo
> ```
>
> ---
>
>
> /redo
> ```
>
> ---
>
>
> team](/docs/share).
>
> ```bash frame="none"
> /share
> ```
>
> This will create a link to the current conversation and copy it to your clipboard.
>
> :::note
> Conversations are not shared by default.
> :::
>
> Here's an [example conversation](https://opencode.ai/s/4XP1fce5) with OpenCode.
>
> ---
>
>
> /share
> ```
>
> This will create a link to the current conversation and copy it to your clipboard.
>
> :::note
> Conversations are not shared by default.
> :::
>
> Here's an [example conversation](https://opencode.ai/s/4XP1fce5) with OpenCode.
>
> ---

### Source: OpenCode Commands — Built-in

> Custom commands are in addition to the built-in commands like `/init`, `/undo`, `/redo`, `/share`, `/help`. [Learn more](/docs/tui#commands).
>
> ---
>
>
> ## Built-in
>
> opencode includes several built-in commands like `/init`, `/undo`, `/redo`, `/share`, `/help`; [learn more](/docs/tui#commands).
>
> :::note
> Custom commands can override built-in commands.
> :::
>
> If you define a custom command with the same name, it will override the built-in command.
>
> opencode includes several built-in commands like `/init`, `/undo`, `/redo`, `/share`, `/help`; [learn more](/docs/tui#commands).
>
> :::note
> Custom commands can override built-in commands.
> :::
>
> If you define a custom command with the same name, it will override the built-in command.

### Source: OpenCode CLI — commands reference

> ## Commands
>
> The OpenCode CLI also has the following commands.
>
> ---
>
> ### agent
>
> Manage agents for OpenCode.
>
> ```bash
> opencode agent [command]
> ```
>
> ---
>
> #### create
>
> Create a new agent with custom configuration.
>
> ```bash
> opencode agent create
> ```
>
> This command will guide you through creating a new agent with a custom system prompt and permission configuration. Anything you don't allow is denied in the generated agent's frontmatter.
>
> #### Flags
>
> | Flag                                        | Short | Description                                                                                                                                                                                                                |
> | ------------------------------------------- | ----- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | <nobr><code>{"--path"}</code></nobr>        |       | Directory to write the agent file to (defaults to global or `.opencode/agent` based on the prompt)                                                                                                                         |
> | <nobr><code>{"--description"}</code></nobr> |       | What the agent should do                                                                                                                                                                                                   |
> | <nobr><code>{"--mode"}</code></nobr>        |       | Agent mode: `all`, `primary`, or `subagent`                                                                                                                                                                                |
> | <nobr><code>{"--permissions"}</code></nobr> |       | Comma-separated list of permissions to allow (default: all). Available: `bash`, `read`, `edit`, `glob`, `grep`, `webfetch`, `task`, `todowrite`, `websearch`, `lsp`, `skill`. Anything omitted is denied. Alias: `--tools` |
> | <nobr><code>{"--model"}</code></nobr>       | `-m`  | Model to use, in `provider/model` format                                                                                                                                                                                   |
>
> Passing all of `--path`, `--description`, `--mode`, and `--permissions` runs the command non-interactively.
>
> ---
>
> #### list
>
> List all available agents.
>
> ```bash
> opencode agent list
> ```
>
> ---
>
> ### attach
>
> Attach a terminal to an already running OpenCode backend server started via `serve` or `web` commands.
>
> ```bash
> opencode attach [url]
> ```
>
> This allows using the TUI with a remote OpenCode backend. For example:
>
> ```bash
> # Start the backend server for web/mobile access
> opencode web --port 4096 --hostname 0.0.0.0
>
> # In another terminal, attach the TUI to the running backend
> opencode attach http://10.20.30.40:4096
> ```
>
> #### Flags
>
> | Flag                                     | Short | Description                                                                |
> | ---------------------------------------- | ----- | -------------------------------------------------------------------------- |
> | <nobr><code>{"--dir"}</code></nobr>      |       | Working directory to start TUI in                                          |
> | <nobr><code>{"--continue"}</code></nobr> | `-c`  | Continue the last session                                                  |
> | <nobr><code>{"--session"}</code></nobr>  | `-s`  | Session ID to continue                                                     |
> | <nobr><code>{"--fork"}</code></nobr>     |       | Fork the session when continuing (use with `--continue` or `--session`)    |
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
