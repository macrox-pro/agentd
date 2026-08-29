---
primary_sources:
  - id: T1-EXT-WRITE
    title: "Build Gemini CLI extensions"
    url: "https://geminicli.com/docs/extensions/writing-extensions.md"
    section: "Full page"
  - id: T1-EXT-REF
    title: "Extension reference"
    url: "https://geminicli.com/docs/extensions/reference.md"
    section: "Full page"
  - id: T1-EXT-BP
    title: "Gemini CLI extension best practices"
    url: "https://geminicli.com/docs/extensions/best-practices.md"
    section: "Full page"
  - id: T1-EXT-REL
    title: "Release extensions"
    url: "https://geminicli.com/docs/extensions/releasing.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Extensions — developer guide

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Build Gemini CLI extensions — Full page

> Gemini CLI extensions let you expand the capabilities of Gemini CLI by adding
> custom tools, commands, and context. This guide walks you through creating your
> first extension, from setting up a template to adding custom functionality and
> linking it for local development.
>
> ## Prerequisites
>
> Before you start, ensure you have Gemini CLI installed and a basic understanding
> of Node.js.
>
> ## Extension features
>
> Extensions offer several ways to customize Gemini CLI. Use this table to decide
> which features your extension needs.
>
> | Feature                                                        | What it is                                                                                                                | When to use it                                                                                                                                                                                                                                                                                 | Invoked by            |
> | :------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :-------------------- |
> | **[MCP server](/docs/extensions/reference#mcp-servers)**                     | A standard way to expose new tools and data sources to the model.                                                         | Use this when you want the model to be able to _do_ new things, like fetching data from an internal API, querying a database, or controlling a local application. We also support MCP resources (which can replace custom commands) and system instructions (which can replace custom context) | Model                 |
> | **[Custom commands](/docs/cli/custom-commands)**               | A shortcut (like `/my-cmd`) that executes a pre-defined prompt or shell command.                                          | Use this for repetitive tasks or to save long, complex prompts that you use frequently. Great for automation.                                                                                                                                                                                  | User                  |
> | **[Context file (`GEMINI.md`)](/docs/extensions/reference#contextfilename)** | A markdown file containing instructions that are loaded into the model's context at the start of every session.           | Use this to define the "personality" of your extension, set coding standards, or provide essential knowledge that the model should always have.                                                                                                                                                | CLI provides to model |
> | **[Agent skills](/docs/cli/skills)**                           | A specialized set of instructions and workflows that the model activates only when needed.                                | Use this for complex, occasional tasks (like "create a PR" or "audit security") to avoid cluttering the main context window when the skill isn't being used.                                                                                                                                   | Model                 |
> | **[Hooks](/docs/hooks)**                                 | A way to intercept and customize the CLI's behavior at specific lifecycle events (for example, before/after a tool call). | Use this when you want to automate actions based on what the model is doing, like validating tool arguments, logging activity, or modifying the model's input/output.                                                                                                                          | CLI                   |
> | **[Custom themes](/docs/extensions/reference#themes)**                       | A set of color definitions to personalize the CLI UI.                                                                     | Use this to provide a unique visual identity for your extension or to offer specialized high-contrast or thematic color schemes.                                                                                                                                                               | User (via /theme)     |
>
> ## Step 1: Create a new extension
>
> The easiest way to start is by using a built-in template. We'll use the
> `mcp-server` example as our foundation.
>
> Run the following command to create a new directory called `my-first-extension`
> with the template files:
>
> ```bash
> gemini extensions new my-first-extension mcp-server
> ```
>
> This creates a directory with the following structure:
>
> ```
> my-first-extension/
> ├── example.js
> ├── gemini-extension.json
> └── package.json
> ```
>
> ## Step 2: Understand the extension files
>
> Your new extension contains several key files that define its behavior.
>
> ### `gemini-extension.json`
>
> The manifest file tells Gemini CLI how to load and use your extension.
>
> ```json
> {
>   "name": "mcp-server-example",
>   "version": "1.0.0",
>   "mcpServers": {
>     "nodeServer": {
>       "command": "node",
>       "args": ["${extensionPath}${/}example.js"],
>       "cwd": "${extensionPath}"
>     }
>   }
> }
> ```
>
> - `name`: The unique name for your extension.
> - `version`: The version of your extension.
> - `mcpServers`: Defines Model Context Protocol (MCP) servers to add new tools.
>   - `command`, `args`, `cwd`: Specify how to start your server. The
>     `${extensionPath}` variable is replaced with the absolute path to your
>     extension's directory.
>
> ### `example.js`
>
> This file contains the source code for your MCP server. It uses the
> `@modelcontextprotocol/sdk` to define tools.
>
> ```javascript
> /**
>  * @license
>  * Copyright 2025 Google LLC
>  * SPDX-License-Identifier: Apache-2.0
>  */
>
> import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
> import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
> import { z } from 'zod';
>
> const server = new McpServer({
>   name: 'prompt-server',
>   version: '1.0.0',
> });
>
> // Registers a new tool named 'fetch_posts'
> server.registerTool(
>   'fetch_posts',
>   {
>     description: 'Fetches a list of posts from a public API.',
>     inputSchema: z.object({}).shape,
>   },
>   async () => {
>     const apiResponse = await fetch(
>       'https://jsonplaceholder.typicode.com/posts',
>     );
>     const posts = await apiResponse.json();
>     const response = { posts: posts.slice(0, 5) };
>     return {
>       content: [
>         {
>           type: 'text',
>           text: JSON.stringify(response),
>         },
>       ],
>     };
>   },
> );
>
> const transport = new StdioServerTransport();
> await server.connect(transport);
> ```
>
> ### `package.json`
>
> The standard configuration file for a Node.js project. It defines dependencies
> and scripts for your extension.
>
> ## Step 3: Add extension settings
>
> Some extensions need configuration, such as API keys or user preferences. Let's
> add a setting for an API key.
>
> 1.  Open `gemini-extension.json`.
> 2.  Add a `settings` array to the configuration:
>
>     ```json
>     {
>       "name": "mcp-server-example",
>       "version": "1.0.0",
>       "settings": [
>         {
>           "name": "API Key",
>           "description": "The API key for the service.",
>           "envVar": "MY_SERVICE_API_KEY",
>           "sensitive": true
>         }
>       ],
>       "mcpServers": {
>         // ...
>       }
>     }
>     ```
>
> When a user installs this extension, Gemini CLI will prompt them to enter the
> "API Key". The value will be stored securely in the system keychain (because
> `sensitive` is true) and injected into the MCP server's process as the
> `MY_SERVICE_API_KEY` environment variable.
>
> > **Important (Environment Variable Sanitization):** For security reasons,
> > sensitive environment variables are filtered out and not passed to extensions
> > or MCP servers by default. Extensions will _only_ have access to environment
> > variables that are explicitly declared in the `settings` array using the
> > `envVar` property, plus a few standard safe variables. Do not expect host
> > environment variables to be available otherwise.
>
> ## Step 4: Link your extension
>
> Link your extension to your Gemini CLI installation for local development.
>
> 1.  **Install dependencies:**
>
>     ```bash
>     cd my-first-extension
>     npm install
>     ```
>
> 2.  **Link the extension:**
>
>     The `link` command creates a symbolic link from Gemini CLI extensions
>     directory to your development directory. Changes you make are reflected
>     immediately.
>
>     ```bash
>     gemini extensions link .
>     ```
>
> Restart your Gemini CLI session to use the new `fetch_posts` tool. Test it by
> asking: "fetch posts".
>
> ## Step 5: Add a custom command
>
> Custom commands create shortcuts for complex prompts.
>
> 1.  Create a `commands` directory and a subdirectory for your command group:
>
>     **macOS/Linux**
>
>     ```bash
>     mkdir -p commands/fs
>     ```
>
>     **Windows (PowerShell)**
>
>     ```powershell
>     New-Item -ItemType Directory -Force -Path "commands\fs"
>     ```
>
> 2.  Create a file named `commands/fs/grep-code.toml`:
>
>     ```toml
>     prompt = """
>     Please summarize the findings for the pattern `{{args}}`.
>
>     Search Results:
>     !{grep -r {{args}} .}
>     """
>     ```
>
>     This command, `/fs:grep-code`, takes an argument, runs the `grep` shell
>     command, and pipes the results into a prompt for summarization.
>
> After saving the file, restart Gemini CLI. Run `/fs:grep-code "some pattern"` to
> use your new command.
>
> ## Step 6: Add a custom `GEMINI.md`
>
> Provide persistent context to the model by adding a `GEMINI.md` file to your
> extension. This is useful for setting behavior or providing essential tool
> information.
>
> 1.  Create a file named `GEMINI.md` in the root of your extension directory:
>
>     ```markdown
>     # My First Extension Instructions
>
>     You are an expert developer assistant. When the user asks you to fetch
>     posts, use the `fetch_posts` tool. Be concise in your responses.
>     ```
>
> 2.  Update your `gemini-extension.json` to load this file:
>
>     ```json
>     {
>       "name": "my-first-extension",
>       "version": "1.0.0",
>       "contextFileName": "GEMINI.md",
>       "mcpServers": {
>         "nodeServer": {
>           "command": "node",
>           "args": ["${extensionPath}${/}example.js"],
>           "cwd": "${extensionPath}"
>         }
>       }
>     }
>     ```
>
> Restart Gemini CLI. The model now has the context from your `GEMINI.md` file in
> every session where the extension is active.
>
> ## (Optional) Step 7: Add an Agent Skill
>
> [Agent Skills](/docs/cli/skills) bundle specialized expertise and workflows.
> Skills are activated only when needed, which saves context tokens.
>
> 1.  Create a `skills` directory and a subdirectory for your skill:
>
>     **macOS/Linux**
>
>     ```bash
>     mkdir -p skills/security-audit
>     ```
>
>     **Windows (PowerShell)**
>
>     ```powershell
>     New-Item -ItemType Directory -Force -Path "skills\security-audit"
>     ```
>
> 2.  Create a `skills/security-audit/SKILL.md` file:
>
>     ```markdown
>     ---
>     name: security-audit
>     description:
>       Expertise in auditing code for security vulnerabilities. Use when the user
>       asks to "check for security issues" or "audit" their changes.
>     ---
>
>     # Security Auditor
>
>     You are an expert security researcher. When auditing code:
>
>     1. Look for common vulnerabilities (OWASP Top 10).
>     2. Check for hardcoded secrets or API keys.
>     3. Suggest remediation steps for any findings.
>     ```
>
> Gemini CLI automatically discovers skills bundled with your extension. The model
> activates them when it identifies a relevant task.
>
> ## Step 8: Release your extension
>
> When your extension is ready, share it with others via a Git repository or
> GitHub Releases. Refer to the [Extension Releasing Guide](/docs/extensions/releasing) for
> detailed instructions and learn how to list your extension in the gallery.
>
> ## Next steps
>
> - [Extension reference](/docs/extensions/reference): Deeply understand the extension format,
>   commands, and configuration.
> - [Best practices](/docs/extensions/best-practices): Learn strategies for building great
>   extensions.

### Source: Extension reference — Full page

> This guide covers the `gemini extensions` commands and the structure of the
> `gemini-extension.json` configuration file.
>
> ## Manage extensions
>
> Use the `gemini extensions` command group to manage your extensions from the
> terminal.
>
> Note that commands like `gemini extensions install` are not supported within the
> CLI's interactive mode. However, you can use the `/extensions list` command to
> view installed extensions. All management operations, including updates to slash
> commands, take effect only after you restart the CLI session.
>
> ### Install an extension
>
> Install an extension by providing its GitHub repository URL or a local file
> path.
>
> Gemini CLI creates a copy of the extension during installation. You must run
> `gemini extensions update` to pull changes from the source. To install from
> GitHub, you must have `git` installed on your machine.
>
> ```bash
> gemini extensions install <source> [--ref <ref>] [--auto-update] [--pre-release] [--consent] [--skip-settings]
> ```
>
> - `<source>`: The GitHub URL or local path of the extension.
> - `--ref`: The git ref (branch, tag, or commit) to install.
> - `--auto-update`: Enable automatic updates for this extension.
> - `--pre-release`: Enable installation of pre-release versions.
> - `--consent`: Acknowledge security risks and skip the confirmation prompt.
> - `--skip-settings`: Skip the configuration on install process.
>
> ### Uninstall an extension
>
> To uninstall one or more extensions, use the `uninstall` command:
>
> ```bash
> gemini extensions uninstall <name...>
> ```
>
> ### Disable an extension
>
> Extensions are enabled globally by default. You can disable an extension
> entirely or for a specific workspace.
>
> ```bash
> gemini extensions disable <name> [--scope <scope>]
> ```
>
> - `<name>`: The name of the extension to disable.
> - `--scope`: The scope to disable the extension in (`user` or `workspace`).
>
> ### Enable an extension
>
> Re-enable a disabled extension using the `enable` command:
>
> ```bash
> gemini extensions enable <name> [--scope <scope>]
> ```
>
> - `<name>`: The name of the extension to enable.
> - `--scope`: The scope to enable the extension in (`user` or `workspace`).
>
> ### Update an extension
>
> Update an extension to the version specified in its `gemini-extension.json`
> file.
>
> ```bash
> gemini extensions update <name>
> ```
>
> To update all installed extensions at once:
>
> ```bash
> gemini extensions update --all
> ```
>
> ### Create an extension from a template
>
> Create a new extension directory using a built-in template.
>
> ```bash
> gemini extensions new <path> [template]
> ```
>
> - `<path>`: The directory to create.
> - `[template]`: The template to use (for example, `mcp-server`, `context`,
>   `custom-commands`).
>
> ### Link a local extension
>
> Create a symbolic link between your development directory and Gemini CLI
> extensions directory. This lets you test changes immediately without
> reinstalling.
>
> ```bash
> gemini extensions link <path>
> ```
>
> ## Extension format
>
> Gemini CLI loads extensions from `<home>/.gemini/extensions`. Each extension
> must have a `gemini-extension.json` file in its root directory.
>
> ### `gemini-extension.json`
>
> The manifest file defines the extension's behavior and configuration.
>
> ```json
> {
>   "name": "my-extension",
>   "version": "1.0.0",
>   "description": "My awesome extension",
>   "mcpServers": {
>     "my-server": {
>       "command": "node",
>       "args": ["${extensionPath}/my-server.js"],
>       "cwd": "${extensionPath}"
>     }
>   },
>   "contextFileName": "GEMINI.md",
>   "excludeTools": ["run_shell_command"],
>   "migratedTo": "https://github.com/new-owner/new-extension-repo",
>   "plan": {
>     "directory": ".gemini/plans"
>   }
> }
> ```
>
> - `name`: The name of the extension. This is used to uniquely identify the
>   extension and for conflict resolution when extension commands have the same
>   name as user or project commands. The name should be lowercase or numbers and
>   use dashes instead of underscores or spaces. This is how users will refer to
>   your extension in the CLI. Note that we expect this name to match the
>   extension directory name.
> - `version`: The version of the extension.
> - `description`: A short description of the extension. This will be displayed on
>   [geminicli.com/extensions](https://geminicli.com/extensions).
> - `migratedTo`: The URL of the new repository source for the extension. If this
>   is set, the CLI will automatically check this new source for updates and
>   migrate the extension's installation to the new source if an update is found.
> - `mcpServers`: A map of MCP servers to settings. The key is the name of the
>   server, and the value is the server configuration. These servers will be
>   loaded on startup just like MCP servers defined in a
>   [`settings.json` file](/docs/reference/configuration). If both an extension
>   and a `settings.json` file define an MCP server with the same name, the server
>   defined in the `settings.json` file takes precedence.
>   - Note that all MCP server configuration options are supported except for
>     `trust`.
>   - For portability, you should use `${extensionPath}` to refer to files within
>     your extension directory.
>   - Separate your executable and its arguments using `command` and `args`
>     instead of putting them both in `command`.
> - `contextFileName`: The name of the file that contains the context for the
>   extension. This will be used to load the context from the extension directory.
>   If this property is not used but a `GEMINI.md` file is present in your
>   extension directory, then that file will be loaded.
> - `excludeTools`: An array of tool names to exclude from the model. You can also
>   specify command-specific restrictions for tools that support it, like the
>   `run_shell_command` tool. For example,
>   `"excludeTools": ["run_shell_command(rm -rf)"]` will block the `rm -rf`
>   command. Note that this differs from the MCP server `excludeTools`
>   functionality, which can be listed in the MCP server config.
> - `plan`: Planning features configuration.
>   - `directory`: The directory where planning artifacts are stored. This serves
>     as a fallback if the user hasn't specified a plan directory in their
>     settings. If not specified by either the extension or the user, the default
>     is `~/.gemini/tmp/<project>/<session-id>/plans/`.
>
> When Gemini CLI starts, it loads all the extensions and merges their
> configurations. If there are any conflicts, the workspace configuration takes
> precedence.
>
> ### Extension settings
>
> Extensions can define settings that users provide during installation, such as
> API keys or URLs. These values are stored in a `.env` file within the extension
> directory.
>
> To define settings, add a `settings` array to your manifest:
>
> ```json
> {
>   "name": "my-api-extension",
>   "version": "1.0.0",
>   "settings": [
>     {
>       "name": "API Key",
>       "description": "Your API key for the service.",
>       "envVar": "MY_API_KEY",
>       "sensitive": true
>     }
>   ]
> }
> ```
>
> - `name`: The setting's display name.
> - `description`: A clear explanation of the setting.
> - `envVar`: The environment variable name where the value is stored.
> - `sensitive`: If `true`, the value is stored in the system keychain and
>   obfuscated in the UI.
>
> To update an extension's settings:
>
> ```bash
> gemini extensions config <name> [setting] [--scope <scope>]
> ```
>
> #### Environment variable sanitization
>
> For security reasons, sensitive environment variables are filtered out and not
> passed to extensions or MCP servers by default.
>
> Extensions **will not** inherit the user's full shell environment variables.
> They will only have access to:
>
> 1. Standard safe variables (e.g., `HOME`, `PATH`, `TMPDIR`).
> 2. Variables explicitly declared and requested in the `gemini-extension.json`
>    manifest via the `settings` array (using the `envVar` property).
>
> If your extension requires specific environment variables (like an API key,
> custom host, or config path), you **must** declare them in the `settings` array
> so the CLI can allowlist them for use within the extension.
>
> ### Custom commands
>
> Provide [custom commands](/docs/cli/custom-commands) by placing TOML files in a
> `commands/` subdirectory. Gemini CLI uses the directory structure to determine
> the command name.
>
> For an extension named `gcp`:
>
> - `commands/deploy.toml` becomes `/deploy`
> - `commands/gcs/sync.toml` becomes `/gcs:sync` (namespaced with a colon)
>
> ### Hooks
>
> Intercept and customize CLI behavior using [hooks](/docs/hooks). Define
> hooks in a `hooks/hooks.json` file within your extension directory. Note that
> hooks are not defined in the `gemini-extension.json` manifest.
>
> ### Agent skills
>
> Bundle [agent skills](/docs/cli/skills) to provide specialized workflows. Place
> skill definitions in a `skills/` directory. For example,
> `skills/security-audit/SKILL.md` exposes a `security-audit` skill.
>
> ### Sub-agents
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > Sub-agents are a preview feature currently under active development.
>
> Provide [sub-agents](/docs/core/subagents) that users can delegate tasks to. Add
> agent definition files (`.md`) to an `agents/` directory in your extension root.
>
> ### <a id="policy-engine"></a>Policy Engine
>
> Extensions can contribute policy rules and safety checkers to Gemini CLI
> [Policy Engine](/docs/reference/policy-engine). These rules are defined in
> `.toml` files and take effect when the extension is activated.
>
> To add policies, create a `policies/` directory in your extension's root and
> place your `.toml` policy files inside it. Gemini CLI automatically loads all
> `.toml` files from this directory.
>
> Rules contributed by extensions run in their own tier (tier 2), alongside
> workspace-defined policies. This tier has higher priority than the default rules
> but lower priority than user or admin policies.
>
> <!-- prettier-ignore -->
> > [!WARNING]
> > For security, Gemini CLI ignores any `allow` decisions or `yolo`
> > mode configurations in extension policies. This ensures that an extension
> > cannot automatically approve tool calls or bypass security measures without
> > your confirmation.
>
> **Example `policies.toml`**
>
> ```toml
> [[rule]]
> mcpName = "my_server"
> toolName = "dangerous_tool"
> decision = "ask_user"
> priority = 100
>
> [[safety_checker]]
> mcpName = "my_server"
> toolName = "write_data"
> priority = 200
> [safety_checker.checker]
> type = "in-process"
> name = "allowed-path"
> required_context = ["environment"]
> ```
>
> ### Themes
>
> Extensions can provide custom themes to personalize the CLI UI. Themes are
> defined in the `themes` array in `gemini-extension.json`.
>
> **Example**
>
> ```json
> {
>   "name": "my-green-extension",
>   "version": "1.0.0",
>   "themes": [
>     {
>       "name": "shades-of-green",
>       "type": "custom",
>       "background": {
>         "primary": "#1a362a"
>       },
>       "text": {
>         "primary": "#a6e3a1",
>         "secondary": "#6e8e7a",
>         "link": "#89e689"
>       },
>       "status": {
>         "success": "#76c076",
>         "warning": "#d9e689",
>         "error": "#b34e4e"
>       },
>       "border": {
>         "default": "#4a6c5a"
>       },
>       "ui": {
>         "comment": "#6e8e7a"
>       }
>     }
>   ]
> }
> ```
>
> Custom themes provided by extensions can be selected using the `/theme` command
> or by setting the `ui.theme` property in your `settings.json` file. Note that
> when referring to a theme from an extension, the extension name is appended to
> the theme name in parentheses, for example,
> `shades-of-green (my-green-extension)`.
>
> ### Conflict resolution
>
> Extension commands have the lowest precedence. If an extension command name
> conflicts with a user or project command, the extension command is prefixed with
> the extension name (for example, `/gcp.deploy`) using a dot separator.
>
> ## Variables
>
> Gemini CLI supports variable substitution in `gemini-extension.json` and
> `hooks/hooks.json`.
>
> | Variable           | Description                                     |
> | :----------------- | :---------------------------------------------- |
> | `${extensionPath}` | The absolute path to the extension's directory. |
> | `${workspacePath}` | The absolute path to the current workspace.     |
> | `${/}`             | The platform-specific path separator.           |

### Source: Gemini CLI extension best practices — Full page

> This guide covers best practices for developing, securing, and maintaining
> Gemini CLI extensions.
>
> ## Development
>
> Developing extensions for Gemini CLI is a lightweight, iterative process. Use
> these strategies to build robust and efficient extensions.
>
> ### Structure your extension
>
> While simple extensions may contain only a few files, we recommend a organized
> structure for complex projects.
>
> ```text
> my-extension/
> ├── package.json
> ├── tsconfig.json
> ├── gemini-extension.json
> ├── src/
> │   ├── index.ts
> │   └── tools/
> └── dist/
> ```
>
> - **Use TypeScript:** We strongly recommend using TypeScript for type safety and
>   improved developer experience.
> - **Separate source and build:** Keep your source code in `src/` and output
>   build artifacts to `dist/`.
> - **Bundle dependencies:** If your extension has many dependencies, bundle them
>   using a tool like `esbuild` to reduce installation time and avoid conflicts.
>
> ### Iterate with `link`
>
> Use the `gemini extensions link` command to develop locally without reinstalling
> your extension after every change.
>
> ```bash
> cd my-extension
> gemini extensions link .
> ```
>
> Changes to your code are immediately available in the CLI after you rebuild the
> project and restart the session.
>
> ### Use `GEMINI.md` effectively
>
> Your `GEMINI.md` file provides essential context to the model.
>
> - **Focus on goals:** Explain the high-level purpose of the extension and how to
>   interact with its tools.
> - **Be concise:** Avoid dumping exhaustive documentation into the file. Use
>   clear, direct language.
> - **Provide examples:** Include brief examples of how the model should use
>   specific tools or commands.
>
> ## Security
>
> Follow the principle of least privilege and rigorous input validation when
> building extensions.
>
> ### Minimal permissions
>
> Only request the permissions your MCP server needs to function. Avoid giving the
> model broad access (such as full shell access) if restricted tools are
> sufficient.
>
> If your extension uses powerful tools like `run_shell_command`, restrict them in
> your `gemini-extension.json` file:
>
> ```json
> {
>   "name": "my-safe-extension",
>   "excludeTools": ["run_shell_command(rm -rf *)"]
> }
> ```
>
> This ensures the CLI blocks dangerous commands even if the model attempts to
> execute them.
>
> ### Validate inputs
>
> Your MCP server runs on the user's machine. Always validate tool inputs to
> prevent arbitrary code execution or unauthorized filesystem access.
>
> ```typescript
> // Example: Validating paths
> if (!path.resolve(inputPath).startsWith(path.resolve(allowedDir) + path.sep)) {
>   throw new Error('Access denied');
> }
> ```
>
> ### Secure sensitive settings
>
> If your extension requires API keys or other secrets, use the `sensitive: true`
> option in your manifest. This ensures keys are stored in the system keychain and
> obfuscated in the CLI output.
>
> ```json
> "settings": [
>   {
>     "name": "API Key",
>     "envVar": "MY_API_KEY",
>     "sensitive": true
>   }
> ]
> ```
>
> ## Release
>
> Follow standard versioning and release practices to ensure a smooth experience
> for your users.
>
> ### Semantic versioning
>
> Follow [Semantic Versioning (SemVer)](https://semver.org/) to communicate
> changes clearly.
>
> - **Major:** Breaking changes (for example, renaming tools or changing
>   arguments).
> - **Minor:** New features (for example, adding new tools or commands).
> - **Patch:** Bug fixes and performance improvements.
>
> ### Release channels
>
> Use Git branches to manage release channels. This lets users choose between
> stability and the latest features.
>
> ```bash
> # Install the stable version (default branch)
> gemini extensions install github.com/user/repo
>
> # Install the development version
> gemini extensions install github.com/user/repo --ref dev
> ```
>
> ### Clean artifacts
>
> When using GitHub Releases, ensure your archives only contain necessary files
> (such as `dist/`, `gemini-extension.json`, and `package.json`). Exclude
> `node_modules/` and `src/` to minimize download size.
>
> ## Test and verify
>
> Test your extension thoroughly before releasing it to users.
>
> - **Manual verification:** Use `gemini extensions link` to test your extension
>   in a live CLI session. Verify that tools appear in the debug console (F12) and
>   that custom commands resolve correctly.
> - **Automated testing:** If your extension includes an MCP server, write unit
>   tests for your tool logic using a framework like Vitest or Jest. You can test
>   MCP tools in isolation by mocking the transport layer.
>
> ## Troubleshooting
>
> Use these tips to diagnose and fix common extension issues.
>
> ### Extension not loading
>
> If your extension doesn't appear in `/extensions list`:
>
> - **Check the manifest:** Ensure `gemini-extension.json` is in the root
>   directory and contains valid JSON.
> - **Verify the name:** The `name` field in the manifest must match the extension
>   directory name exactly.
> - **Restart the CLI:** Extensions are loaded at the start of a session. Restart
>   Gemini CLI after making changes to the manifest or linking a new extension.
>
> ### MCP server failures
>
> If your tools aren't working as expected:
>
> - **Check the logs:** View the CLI logs to see if the MCP server failed to
>   start.
> - **Test the command:** Run the server's `command` and `args` directly in your
>   terminal to ensure it starts correctly outside of Gemini CLI.
> - **Debug console:** In interactive mode, press **F12** to open the debug
>   console and inspect tool calls and responses.
>
> ### Command conflicts
>
> If a custom command isn't responding:
>
> - **Check precedence:** Remember that user and project commands take precedence
>   over extension commands. Use the prefixed name (for example,
>   `/extension.command`) to verify the extension's version.
> - **Help command:** Run `/help` to see a list of all available commands and
>   their sources.

### Source: Release extensions — Full page

> Release Gemini CLI extensions to your users through a Git repository or GitHub
> Releases. This guide explains how to share your work, list it in the gallery,
> and manage updates.
>
> Git repository releases are the simplest approach and offer the most flexibility
> for managing development branches. GitHub Releases are more efficient for
> initial installations because they ship as single archives rather than requiring
> a full `git clone`. Use GitHub Releases if you need to include platform-specific
> binary files.
>
> ## List your extension in the gallery
>
> The [Gemini CLI extension gallery](https://geminicli.com/extensions/browse/)
> automatically indexes public extensions to help users discover your work. You
> don't need to submit an issue or email us to list your extension.
>
> To have your extension automatically discovered and listed:
>
> 1.  **Use a public repository:** Ensure your extension is hosted in a public
>     GitHub repository.
> 2.  **Add the GitHub topic:** Add the `gemini-cli-extension` topic to your
>     repository's **About** section. Our crawler uses this topic to find new
>     extensions.
> 3.  **Place the manifest at the root:** Ensure your `gemini-extension.json` file
>     is in the absolute root of the repository or the release archive.
>
> Our system crawls tagged repositories daily. Once you tag your repository, your
> extension will appear in the gallery if it passes validation.
>
> ## Release through a Git repository
>
> Releasing through Git is the most flexible option. Create a public Git
> repository and provide the URL to your users. They can then install your
> extension using `gemini extensions install <your-repo-uri>`.
>
> Users can optionally depend on a specific branch, tag, or commit using the
> `--ref` argument. For example:
>
> ```bash
> gemini extensions install <your-repo-uri> --ref=stable
> ```
>
> Whenever you push commits to the referenced branch, the CLI prompts users to
> update their installation. The `HEAD` commit is always treated as the latest
> version.
>
> ### Manage release channels
>
> You can use branches or tags to manage different release channels, such as
> `stable`, `preview`, or `dev`.
>
> We recommend using your default branch as the stable release channel. This
> ensures that the default installation command always provides the most reliable
> version of your extension. You can then use a `dev` branch for active
> development and merge it into the default branch when you are ready for a
> release.
>
> ## Release through GitHub Releases
>
> Distributing extensions through
> [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases)
> provides a faster installation experience by avoiding a repository clone.
>
> Gemini CLI checks for updates by looking for the **Latest** release on GitHub.
> Users can also install specific versions using the `--ref` argument with a
> release tag. Use the `--pre-release` flag to install the latest version even if
> it isn't marked as **Latest**.
>
> ### Custom pre-built archives
>
> You can attach custom archives directly to your GitHub Release as assets. This
> is useful if your extension requires a build step or includes platform-specific
> binaries.
>
> Custom archives must be fully self-contained and follow the required
> [archive structure](#archive-structure). If your extension is
> platform-independent, provide a single generic asset.
>
> #### Platform-specific archives
>
> To let Gemini CLI find the correct asset for a user's platform, use the
> following naming convention:
>
> 1.  **Platform and architecture-specific:**
>     `{platform}.{arch}.{name}.{extension}`
> 2.  **Platform-specific:** `{platform}.{name}.{extension}`
> 3.  **Generic:** A single asset will be used as a fallback if no specific match
>     is found.
>
> Use these values for the placeholders:
>
> - `{name}`: Your extension name.
> - `{platform}`: Use `darwin` (macOS), `linux`, or `win32` (Windows).
> - `{arch}`: Use `x64` or `arm64`.
> - `{extension}`: Use `.tar.gz` or `.zip`.
>
> **Examples:**
>
> - `darwin.arm64.my-tool.tar.gz` (specific to Apple Silicon Macs)
> - `darwin.my-tool.tar.gz` (fallback for all Macs, for example Intel)
> - `linux.x64.my-tool.tar.gz`
> - `win32.my-tool.zip`
>
> #### Archive structure
>
> Archives must be fully contained extensions. The `gemini-extension.json` file
> must be at the root of the archive. The rest of the layout should match a
> standard extension structure.
>
> #### Example GitHub Actions workflow
>
> Use this example workflow to build and release your extension for multiple
> platforms:
>
> ```yaml
> name: Release Extension
>
> on:
>   push:
>     tags:
>       - 'v*'
>
> jobs:
>   release:
>     runs-on: ubuntu-latest
>     steps:
>       - uses: actions/checkout@v3
>
>       - name: Set up Node.js
>         uses: actions/setup-node@v3
>         with:
>           node-version: '20'
>
>       - name: Install dependencies
>         run: npm ci
>
>       - name: Build extension
>         run: npm run build
>
>       - name: Create release assets
>         run: |
>           npm run package -- --platform=darwin --arch=arm64
>           npm run package -- --platform=linux --arch=x64
>           npm run package -- --platform=win32 --arch=x64
>
>       - name: Create GitHub Release
>         uses: softprops/action-gh-release@v1
>         with:
>           files: |
>             release/darwin.arm64.my-tool.tar.gz
>             release/linux.arm64.my-tool.tar.gz
>             release/win32.arm64.my-tool.zip
> ```
>
> ## Migrate an extension repository
>
> If you move your extension to a new repository or rename it, use the
> `migratedTo` property in `gemini-extension.json` to seamlessly transition your
> users.
>
> 1.  **Create the new repository:** Set up your extension in its new location.
> 2.  **Update the old repository:** In your original repository, update the
>     `gemini-extension.json` file to include the `migratedTo` property pointing
>     to the new repository URL, and increment the version number.
>     ```json
>     {
>       "name": "my-extension",
>       "version": "1.1.0",
>       "migratedTo": "https://github.com/new-owner/new-extension-repo"
>     }
>     ```
> 3.  **Release the update:** Publish this new version in your old repository.
>
> When users check for updates, Gemini CLI detects the `migratedTo` field,
> verifies the new repository, and automatically updates their local installation
> to track the new source. All settings migrate automatically.
>
> ## How updates work
>
> Gemini CLI automatically checks for extension updates based on the installation
> method. Understanding these mechanisms helps you ensure your users always have
> the latest version.
>
> ### Sync manifest and tags
>
> For GitHub releases, always ensure the `version` in `gemini-extension.json`
> matches your GitHub release tag. While the CLI uses tags for update detection,
> it displays the manifest version in the UI. Keeping them in sync prevents
> confusion.
>
> ### Update mechanisms
>
> <details>
> <summary>Technical update details</summary>
>
> The CLI uses different strategies depending on the installation type:
>
> - **GitHub releases:** The CLI queries the GitHub API for the latest release
>   tag. It ignores the `version` field in the manifest for detection.
> - **Git clones:** The CLI runs `git ls-remote` to compare the latest remote
>   commit hash with your local `HEAD`.
> - **Local extensions:** The CLI compares the `version` field in the source
>   directory's manifest with the installed version.
>
> To verify an extension's installation type, inspect the `type` field in the
> metadata file at `~/.gemini/extensions/<name>/.gemini-extension-install.json`.
>
> </details>
>
> <!-- prettier-ignore -->
> > [!IMPORTANT]
> > The `migratedTo` flow requires at least one release on the new repository for
> > the CLI to recognize it as a valid update source.
