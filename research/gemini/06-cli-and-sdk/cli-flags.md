---
primary_sources:
  - id: T1-CLI-REF
    title: "Gemini CLI cheatsheet"
    url: "https://geminicli.com/docs/cli/cli-reference.md"
    section: "Full page"
  - id: T1-INSTALL
    title: "Gemini CLI installation, execution, and releases"
    url: "https://geminicli.com/docs/get-started/installation.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# CLI flags and cheatsheet

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Gemini CLI cheatsheet — Full page

> This page provides a reference for commonly used Gemini CLI commands, options,
> and parameters.
>
> ## CLI commands
>
> | Command                            | Description                        | Example                                                      |
> | ---------------------------------- | ---------------------------------- | ------------------------------------------------------------ |
> | `gemini`                           | Start interactive REPL             | `gemini`                                                     |
> | `gemini -p "query"`                | Query non-interactively            | `gemini -p "summarize README.md"`                            |
> | gemini "query"                     | Query and continue interactively   | gemini "explain this project"                                |
> | `cat file \| gemini`               | Process piped content              | `cat logs.txt \| gemini`<br>`Get-Content logs.txt \| gemini` |
> | `gemini -i "query"`                | Execute and continue interactively | `gemini -i "What is the purpose of this project?"`           |
> | `gemini -r "latest"`               | Continue most recent session       | `gemini -r "latest"`                                         |
> | `gemini -r "latest" "query"`       | Continue session with a new prompt | `gemini -r "latest" "Check for type errors"`                 |
> | `gemini -r "<session-id>" "query"` | Resume session by ID               | `gemini -r "abc123" "Finish this PR"`                        |
> | `gemini update`                    | Update to latest version           | `gemini update`                                              |
> | `gemini extensions`                | Manage extensions                  | See [Extensions Management](#extensions-management)          |
> | `gemini mcp`                       | Configure MCP servers              | See [MCP Server Management](#mcp-server-management)          |
>
> ### Positional arguments
>
> | Argument | Type              | Description                                                                                                |
> | -------- | ----------------- | ---------------------------------------------------------------------------------------------------------- |
> | `query`  | string (variadic) | Positional prompt. Defaults to interactive mode in a TTY. Use `-p/--prompt` for non-interactive execution. |
>
> ## Interactive commands
>
> These commands are available within the interactive REPL.
>
> | Command              | Description                                     |
> | -------------------- | ----------------------------------------------- |
> | `/skills reload`     | Reload discovered skills from disk              |
> | `/agents reload`     | Reload the agent registry                       |
> | `/commands list`     | List available custom slash commands            |
> | `/commands reload`   | Reload custom slash commands                    |
> | `/memory reload`     | Reload context files (for example, `GEMINI.md`) |
> | `/mcp reload`        | Restart and reload MCP servers                  |
> | `/extensions reload` | Reload all active extensions                    |
> | `/help`              | Show help for all commands                      |
> | `/quit`              | Exit the interactive session                    |
>
> ## CLI Options
>
> | Option                           | Alias | Type    | Default   | Description                                                                                                                                                            |
> | -------------------------------- | ----- | ------- | --------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `--debug`                        | `-d`  | boolean | `false`   | Run in debug mode with verbose logging                                                                                                                                 |
> | `--version`                      | `-v`  | -       | -         | Show CLI version number and exit                                                                                                                                       |
> | `--help`                         | `-h`  | -       | -         | Show help information                                                                                                                                                  |
> | `--model`                        | `-m`  | string  | `auto`    | Model to use. See [Model Selection](#model-selection) for available values.                                                                                            |
> | `--prompt`                       | `-p`  | string  | -         | Prompt text. Appended to stdin input if provided. Forces non-interactive mode.                                                                                         |
> | `--prompt-interactive`           | `-i`  | string  | -         | Execute prompt and continue in interactive mode                                                                                                                        |
> | `--worktree`                     | `-w`  | string  | -         | Start Gemini in a new git worktree. If no name is provided, one is generated automatically. Requires `experimental.worktrees: true` in settings.                       |
> | `--sandbox`                      | `-s`  | boolean | `false`   | Run in a sandboxed environment for safer execution                                                                                                                     |
> | `--skip-trust`                   | -     | boolean | `false`   | Trust the current workspace for this session, skipping the folder trust check.                                                                                         |
> | `--approval-mode`                | -     | string  | `default` | Approval mode for tool execution. Choices: `default`, `auto_edit`, `yolo`, `plan`                                                                                      |
> | `--yolo`                         | `-y`  | boolean | `false`   | **Deprecated.** Auto-approve all actions. Use `--approval-mode=yolo` instead.                                                                                          |
> | `--experimental-acp`             | -     | boolean | -         | Start in ACP (Agent Code Pilot) mode. **Experimental feature.**                                                                                                        |
> | `--experimental-zed-integration` | -     | boolean | -         | Run in Zed editor integration mode. **Experimental feature.**                                                                                                          |
> | `--allowed-mcp-server-names`     | -     | array   | -         | Allowed MCP server names (comma-separated or multiple flags)                                                                                                           |
> | `--allowed-tools`                | -     | array   | -         | **Deprecated.** Use the [Policy Engine](/docs/reference/policy-engine) instead. Tools that are allowed to run without confirmation (comma-separated or multiple flags) |
> | `--extensions`                   | `-e`  | array   | -         | List of extensions to use. If not provided, all extensions are enabled (comma-separated or multiple flags)                                                             |
> | `--list-extensions`              | `-l`  | boolean | -         | List all available extensions and exit                                                                                                                                 |
> | `--resume`                       | `-r`  | string  | -         | Resume a previous session. Use `"latest"` for most recent or index number (for example `--resume 5`)                                                                   |
> | `--list-sessions`                | -     | boolean | -         | List available sessions for the current project and exit                                                                                                               |
> | `--delete-session`               | -     | string  | -         | Delete a session by index number (use `--list-sessions` to see available sessions)                                                                                     |
> | `--include-directories`          | -     | array   | -         | Additional directories to include in the workspace (comma-separated or multiple flags)                                                                                 |
> | `--screen-reader`                | -     | boolean | -         | Enable screen reader mode for accessibility                                                                                                                            |
> | `--output-format`                | `-o`  | string  | `text`    | The format of the CLI output. Choices: `text`, `json`, `stream-json`                                                                                                   |
>
> ## Model selection
>
> The `--model` (or `-m`) flag lets you specify which Gemini model to use. You can
> use either model aliases (user-friendly names) or concrete model names.
>
> ### Model aliases
>
> These are convenient shortcuts that map to specific models:
>
> | Alias        | Resolves To                                | Description                                                                                                               |
> | ------------ | ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------- |
> | `auto`       | `gemini-2.5-pro` or `gemini-3-pro-preview` | **Default.** Resolves to the preview model if preview features are enabled, otherwise resolves to the standard pro model. |
> | `pro`        | `gemini-2.5-pro` or `gemini-3-pro-preview` | For complex reasoning tasks. Uses preview model if enabled.                                                               |
> | `flash`      | `gemini-2.5-flash`                         | Fast, balanced model for most tasks.                                                                                      |
> | `flash-lite` | `gemini-2.5-flash-lite`                    | Fastest model for simple tasks.                                                                                           |
>
> ## Extensions management
>
> | Command                                            | Description                                  | Example                                                                        |
> | -------------------------------------------------- | -------------------------------------------- | ------------------------------------------------------------------------------ |
> | `gemini extensions install <source>`               | Install extension from Git URL or local path | `gemini extensions install https://github.com/user/my-extension`               |
> | `gemini extensions install <source> --ref <ref>`   | Install from specific branch/tag/commit      | `gemini extensions install https://github.com/user/my-extension --ref develop` |
> | `gemini extensions install <source> --auto-update` | Install with auto-update enabled             | `gemini extensions install https://github.com/user/my-extension --auto-update` |
> | `gemini extensions uninstall <name>`               | Uninstall one or more extensions             | `gemini extensions uninstall my-extension`                                     |
> | `gemini extensions list`                           | List all installed extensions                | `gemini extensions list`                                                       |
> | `gemini extensions update <name>`                  | Update a specific extension                  | `gemini extensions update my-extension`                                        |
> | `gemini extensions update --all`                   | Update all extensions                        | `gemini extensions update --all`                                               |
> | `gemini extensions enable <name>`                  | Enable an extension                          | `gemini extensions enable my-extension`                                        |
> | `gemini extensions disable <name>`                 | Disable an extension                         | `gemini extensions disable my-extension`                                       |
> | `gemini extensions link <path>`                    | Link local extension for development         | `gemini extensions link /path/to/extension`                                    |
> | `gemini extensions new <path>`                     | Create new extension from template           | `gemini extensions new ./my-extension`                                         |
> | `gemini extensions validate <path>`                | Validate extension structure                 | `gemini extensions validate ./my-extension`                                    |
>
> See [Extensions Documentation](/docs/extensions) for more details.
>
> ## MCP server management
>
> | Command                                                       | Description                     | Example                                                                                              |
> | ------------------------------------------------------------- | ------------------------------- | ---------------------------------------------------------------------------------------------------- |
> | `gemini mcp add <name> <command>`                             | Add stdio-based MCP server      | `gemini mcp add github npx -y @modelcontextprotocol/server-github`                                   |
> | `gemini mcp add <name> <url> --transport http`                | Add HTTP-based MCP server       | `gemini mcp add api-server http://localhost:3000 --transport http`                                   |
> | `gemini mcp add <name> <command> --env KEY=value`             | Add with environment variables  | `gemini mcp add slack node server.js --env SLACK_TOKEN=xoxb-xxx`                                     |
> | `gemini mcp add <name> <command> --scope user`                | Add with user scope             | `gemini mcp add db node db-server.js --scope user`                                                   |
> | `gemini mcp add <name> <command> --include-tools tool1,tool2` | Add with specific tools         | `gemini mcp add github npx -y @modelcontextprotocol/server-github --include-tools list_repos,get_pr` |
> | `gemini mcp remove <name>`                                    | Remove an MCP server            | `gemini mcp remove github`                                                                           |
> | `gemini mcp list`                                             | List all configured MCP servers | `gemini mcp list`                                                                                    |
>
> See [MCP Server Integration](/docs/tools/mcp-server) for more details.
>
> ## Skills management
>
> | Command                          | Description                           | Example                                           |
> | -------------------------------- | ------------------------------------- | ------------------------------------------------- |
> | `gemini skills list`             | List all discovered agent skills      | `gemini skills list`                              |
> | `gemini skills install <source>` | Install skill from Git, path, or file | `gemini skills install https://github.com/u/repo` |
> | `gemini skills link <path>`      | Link local agent skills via symlink   | `gemini skills link /path/to/my-skills`           |
> | `gemini skills uninstall <name>` | Uninstall an agent skill              | `gemini skills uninstall my-skill`                |
> | `gemini skills enable <name>`    | Enable an agent skill                 | `gemini skills enable my-skill`                   |
> | `gemini skills disable <name>`   | Disable an agent skill                | `gemini skills disable my-skill`                  |
> | `gemini skills enable --all`     | Enable all skills                     | `gemini skills enable --all`                      |
> | `gemini skills disable --all`    | Disable all skills                    | `gemini skills disable --all`                     |
>
> See [Agent Skills Documentation](/docs/cli/skills) for more details.

### Source: Gemini CLI installation, execution, and releases — Full page

> import { Tabs, TabItem } from '@astrojs/starlight/components';
>
>
>
> This document provides an overview of Gemini CLI's system requirements,
> installation methods, and release types.
>
> ## Recommended system specifications
>
> - **Operating System:**
>   - macOS 15+
>   - Windows 11 24H2+
>   - Ubuntu 20.04+
> - **Hardware:**
>   - "Casual" usage: 4GB+ RAM (short sessions, common tasks and edits)
>   - "Power" usage: 16GB+ RAM (long sessions, large codebases, deep context)
> - **Runtime:** Node.js 20.0.0+
> - **Shell:** Bash, Zsh, or PowerShell
> - **Location:**
>   [Gemini Code Assist supported locations](https://developers.google.com/gemini-code-assist/resources/available-locations#americas)
> - **Internet connection required**
>
> ## Install Gemini CLI
>
> We recommend most users install Gemini CLI using one of the following
> installation methods. Note that Gemini CLI comes pre-installed on
> [**Cloud Shell**](https://docs.cloud.google.com/shell/docs) and
> [**Cloud Workstations**](https://cloud.google.com/workstations).
>
> <Tabs>
>   <TabItem label="npm">
>
>   Install globally with npm:
>
>   ```bash
>   npm install -g @google/gemini-cli
>   ```
>
>   </TabItem>
>   <TabItem label="Homebrew (macOS/Linux)">
>
>   Install globally with Homebrew:
>
>   ```bash
>   brew install gemini-cli
>   ```
>
>   </TabItem>
>   <TabItem label="MacPorts (macOS)">
>
>   Install globally with MacPorts:
>
>   ```bash
>   sudo port install gemini-cli
>   ```
>
>   </TabItem>
>   <TabItem label="Anaconda">
>
>   Install with Anaconda (for restricted environments):
>
>   ```bash
>   # Create and activate a new environment
>   conda create -y -n gemini_env -c conda-forge nodejs
>   conda activate gemini_env
>
>   # Install Gemini CLI globally via npm (inside the environment)
>   npm install -g @google/gemini-cli
>   ```
>
>   </TabItem>
> </Tabs>
>
> ## Run Gemini CLI
>
> For most users, we recommend running Gemini CLI with the `gemini` command:
>
> ```bash
> gemini
> ```
>
> For a list of options and additional commands, see the
> [CLI cheatsheet](/docs/cli/cli-reference).
>
> You can also run Gemini CLI using one of the following advanced methods:
>
> <Tabs>
>   <TabItem label="npx">
>
>   Run instantly with npx. You can run Gemini CLI without permanent installation.
>
>   ```bash
>   # Using npx (no installation required)
>   npx @google/gemini-cli
>   ```
>
>   You can also execute the CLI directly from the main branch on GitHub, which is
>   helpful for testing features still in development:
>
>   ```bash
>   npx https://github.com/google-gemini/gemini-cli
>   ```
>
>   </TabItem>
>   <TabItem label="Docker/Podman Sandbox">
>
>   For security and isolation, Gemini CLI can be run inside a container. This is
>   the default way that the CLI executes tools that might have side effects.
>
>   - **Directly from the registry:** You can run the published sandbox image
>     directly. This is useful for environments where you only have Docker and want
>     to run the CLI.
>     ```bash
>     # Run the published sandbox image for a specified CLI version
>     docker run --rm -it us-docker.pkg.dev/gemini-code-dev/gemini-cli/sandbox:0.42.0-nightly.20260428.g59b2dea0e
>     ```
>   - **Using the `--sandbox` flag:** If you have Gemini CLI installed locally
>     (using the standard installation described above), you can instruct it to run
>     inside the sandbox container.
>     ```bash
>     gemini --sandbox -y -p "your prompt here"
>     ```
>
>   </TabItem>
>   <TabItem label="From source">
>
>   Contributors to the project will want to run the CLI directly from the source
>   code.
>
>   - **Development mode:** This method provides hot-reloading and is useful for
>     active development.
>     ```bash
>     # From the root of the repository
>     npm run start
>     ```
>   - **Production mode (React optimizations):** This method runs the CLI with React
>     production mode enabled, which is useful for testing performance without
>     development overhead.
>     ```bash
>     # From the root of the repository
>     npm run start:prod
>     ```
>   - **Production-like mode (linked package):** This method simulates a global
>     installation by linking your local package. It's useful for testing a local
>     build in a production workflow.
>
>     ```bash
>     # Link the local cli package to your global node_modules
>     npm link packages/cli
>
>     # Now you can run your local version using the `gemini` command
>     gemini
>     ```
>
>   </TabItem>
> </Tabs>
>
> ## Releases
>
> Gemini CLI has three release channels: stable, preview, and nightly. For most
> users, we recommend the stable release, which is the default installation.
>
> <Tabs>
>   <TabItem label="Stable">
>
>   Stable releases are published each week. A stable release is created from the
>   previous week's preview release along with any bug fixes. The stable release
>   uses the `latest` tag. Omitting the tag also installs the latest stable
>   release by default.
>
>   ```bash
>   # Both commands install the latest stable release.
>   npm install -g @google/gemini-cli
>   npm install -g @google/gemini-cli@latest
>   ```
>
>   </TabItem>
>   <TabItem label="Preview">
>
>   New preview releases will be published each week. These releases are not fully
>   vetted and may contain regressions or other outstanding issues. Try out the
>   preview release by using the `preview` tag:
>
>   ```bash
>   npm install -g @google/gemini-cli@preview
>   ```
>
>   </TabItem>
>   <TabItem label="Nightly">
>
>   Nightly releases are published every day. The nightly release includes all
>   changes from the main branch at time of release. It should be assumed there are
>   pending validations and issues. You can help test the latest changes by
>   installing with the `nightly` tag:
>
>   ```bash
>   npm install -g @google/gemini-cli@nightly
>   ```
>
>   </TabItem>
> </Tabs>
