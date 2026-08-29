---
primary_sources:
  - id: T1-IDE
    title: "IDE Integration"
    url: "https://geminicli.com/docs/ide-integration.md"
    section: "Full page"
  - id: T1-ACP
    title: "ACP Mode"
    url: "https://geminicli.com/docs/cli/acp-mode.md"
    section: "Full page"
  - id: T1-IDE-SPEC
    title: "Gemini CLI companion plugin: Interface specification"
    url: "https://geminicli.com/docs/ide-integration/ide-companion-spec.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# IDE integration and ACP

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: IDE Integration — Full page

> Gemini CLI can integrate with your IDE to provide a more seamless and
> context-aware experience. This integration allows the CLI to understand your
> workspace better and enables powerful features like native in-editor diffing.
>
> There are two primary ways to integrate Gemini CLI with an IDE:
>
> 1.  **VS Code companion extension**: Install the "Gemini CLI Companion"
>     extension on [Antigravity](https://antigravity.google),
>     [Visual Studio Code](https://code.visualstudio.com/), or other VS Code
>     compatible editors.
> 2.  **Agent Client Protocol (ACP)**: An open protocol for interoperability
>     between AI coding agents and IDEs. This method is used for integrations with
>     tools like JetBrains and Zed, which leverage the ACP Agent Registry for easy
>     discovery and installation of compatible agents like Gemini CLI.
>
> ## VS Code companion extension
>
> The **Gemini CLI Companion extension** grants Gemini CLI direct access to your
> VS Code compatible IDEs and improves your experience by providing real-time
> context such as open files, cursor positions, and text selection. The extension
> also enables a native diffing interface so you can seamlessly review and apply
> AI-generated code changes directly within your editor.
>
> ### Features
>
> - **Workspace context:** The CLI automatically gains awareness of your workspace
>   to provide more relevant and accurate responses. This context includes:
>
>   - The **10 most recently accessed files** in your workspace.
>   - Your active cursor position.
>   - Any text you have selected (up to a 16KB limit; longer selections will be
>     truncated).
>
> - **Native diffing:** When Gemini suggests code modifications, you can view the
>   changes directly within your IDE's native diff viewer. This lets you review,
>   edit, and accept or reject the suggested changes seamlessly.
>
> - **VS Code commands:** You can access Gemini CLI features directly from the VS
>   Code Command Palette (`Cmd+Shift+P` or `Ctrl+Shift+P`):
>   - `Gemini CLI: Run`: Starts a new Gemini CLI session in the integrated
>     terminal.
>   - `Gemini CLI: Accept Diff`: Accepts the changes in the active diff editor.
>   - `Gemini CLI: Close Diff Editor`: Rejects the changes and closes the active
>     diff editor.
>   - `Gemini CLI: View Third-Party Notices`: Displays the third-party notices for
>     the extension.
>
> ### Installation and setup
>
> There are three ways to set up the IDE integration:
>
> #### 1. Automatic nudge (recommended)
>
> When you run Gemini CLI inside a supported editor, it will automatically detect
> your environment and prompt you to connect. Answering "Yes" will automatically
> run the necessary setup, which includes installing the companion extension and
> enabling the connection.
>
> #### 2. Manual installation from CLI
>
> If you previously dismissed the prompt or want to install the extension
> manually, you can run the following command inside Gemini CLI:
>
> ```
> /ide install
> ```
>
> This will find the correct extension for your IDE and install it.
>
> #### 3. Manual installation from a marketplace
>
> You can also install the extension directly from a marketplace.
>
> - **For Visual Studio Code:** Install from the
>   [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=google.gemini-cli-vscode-ide-companion).
> - **For VS Code forks:** To support forks of VS Code, the extension is also
>   published on the
>   [Open VSX Registry](https://open-vsx.org/extension/google/gemini-cli-vscode-ide-companion).
>   Follow your editor's instructions for installing extensions from this
>   registry.
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > The "Gemini CLI Companion" extension may appear towards the bottom of
> > search results. If you don't see it immediately, try scrolling down or
> > sorting by "Newly Published".
> >
> > After manually installing the extension, you must run `/ide enable` in the CLI
> > to activate the integration.
>
> ### Usage
>
> #### Enabling and disabling
>
> You can control the IDE integration from within the CLI:
>
> - To enable the connection to the IDE, run:
>   ```
>   /ide enable
>   ```
> - To disable the connection, run:
>   ```
>   /ide disable
>   ```
>
> When enabled, Gemini CLI will automatically attempt to connect to the IDE
> companion extension.
>
> #### Checking the status
>
> To check the connection status and see the context the CLI has received from the
> IDE, run:
>
> ```
> /ide status
> ```
>
> If connected, this command will show the IDE it's connected to and a list of
> recently opened files it is aware of.
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > The file list is limited to 10 recently accessed files within your
> > workspace and only includes local files on disk.
>
> #### Working with diffs
>
> When you ask Gemini to modify a file, it can open a diff view directly in your
> editor.
>
> **To accept a diff**, you can perform any of the following actions:
>
> - Click the **checkmark icon** in the diff editor's title bar.
> - Save the file (for example, with `Cmd+S` or `Ctrl+S`).
> - Open the Command Palette and run **Gemini CLI: Accept Diff**.
> - Respond with `yes` in the CLI when prompted.
>
> **To reject a diff**, you can:
>
> - Click the **'x' icon** in the diff editor's title bar.
> - Close the diff editor tab.
> - Open the Command Palette and run **Gemini CLI: Close Diff Editor**.
> - Respond with `no` in the CLI when prompted.
>
> You can also **modify the suggested changes** directly in the diff view before
> accepting them.
>
> If you select ‘Allow for this session’ in the CLI, changes will no longer show
> up in the IDE as they will be auto-accepted.
>
> ## Agent Client Protocol (ACP)
>
> ACP is an open protocol that standardizes how AI coding agents communicate with
> code editors and IDEs. It addresses the challenge of fragmented distribution,
> where agents traditionally needed custom integrations for each client. With ACP,
> developers can implement their agent once, and it becomes compatible with any
> ACP-compliant editor.
>
> For a comprehensive introduction to ACP, including its architecture and
> benefits, refer to the official
> [ACP Introduction](https://agentclientprotocol.com/get-started/introduction)
> documentation.
>
> ### The ACP Agent Registry
>
> Gemini CLI is officially available in the **ACP Agent Registry**. This allows
> you to install and update Gemini CLI directly within supporting IDEs and
> eliminates the need for manual downloads or IDE-specific extensions.
>
> Using the registry ensures:
>
> - **Ease of use**: Discover and install agents directly within your IDE
>   settings.
> - **Latest versions**: Ensures users always have access to the most up-to-date
>   agent implementations.
>
> For more details on how the registry works, visit the official
> [ACP Agent Registry](https://agentclientprotocol.com/get-started/registry) page.
> You can learn about how specific IDEs leverage this integration in the following
> section.
>
> ### IDE-specific integration
>
> Gemini CLI is an ACP-compatible agent available in the ACP Agent Registry.
> Here’s how different IDEs leverage the ACP and the registry:
>
> #### JetBrains IDEs
>
> JetBrains IDEs (like IntelliJ IDEA, PyCharm, or GoLand) offer built-in registry
> support, allowing users to find and install ACP-compatible agents directly.
>
> For more details, refer to the official
> [JetBrains AI Blog announcement](https://blog.jetbrains.com/ai/2026/01/acp-agent-registry/).
>
> #### Zed
>
> Zed, a modern code editor, also integrates with the ACP Agent Registry. This
> allows Zed users to easily browse, install, and manage ACP agents.
>
> Learn more about Zed's integration with the ACP Registry in their
> [blog post](https://zed.dev/blog/acp-registry).
>
> #### Other ACP-compatible IDEs
>
> Any other IDE that supports the ACP Agent Registry can install Gemini CLI
> directly through their in-built registry features.
>
> ## Using with sandboxing
>
> If you are using Gemini CLI within a sandbox, be aware of the following:
>
> - **On macOS:** The IDE integration requires network access to communicate with
>   the IDE companion extension. You must use a Seatbelt profile that allows
>   network access.
> - **In a Docker container:** If you run Gemini CLI inside a Docker (or Podman)
>   container, the IDE integration can still connect to the VS Code extension
>   running on your host machine. The CLI is configured to automatically find the
>   IDE server on `host.docker.internal`. No special configuration is usually
>   required, but you may need to ensure your Docker networking setup allows
>   connections from the container to the host.
>
> ## Troubleshooting
>
> ### VS Code companion extension errors
>
> #### Connection errors
>
> - **Message:**
>   `🔴 Disconnected: Failed to connect to IDE companion extension in [IDE Name]. Please ensure the extension is running. To install the extension, run /ide install.`
>
>   - **Cause:** Gemini CLI could not find the necessary environment variables
>     (`GEMINI_CLI_IDE_WORKSPACE_PATH` or `GEMINI_CLI_IDE_SERVER_PORT`) to connect
>     to the IDE. This usually means the IDE companion extension is not running or
>     did not initialize correctly.
>   - **Solution:**
>     1.  Make sure you have installed the **Gemini CLI Companion** extension in
>         your IDE and that it is enabled.
>     2.  Open a new terminal window in your IDE to ensure it picks up the correct
>         environment.
>
> - **Message:**
>   `🔴 Disconnected: IDE connection error. The connection was lost unexpectedly. Please try reconnecting by running /ide enable`
>   - **Cause:** The connection to the IDE companion was lost.
>   - **Solution:** Run `/ide enable` to try and reconnect. If the issue
>     continues, open a new terminal window or restart your IDE.
>
> #### Manual PID override
>
> If automatic IDE detection fails, or if you are running Gemini CLI in a
> standalone terminal and want to manually associate it with a specific IDE
> instance, you can set the `GEMINI_CLI_IDE_PID` environment variable to the
> process ID (PID) of your IDE.
>
> **macOS/Linux**
>
> ```bash
> export GEMINI_CLI_IDE_PID=12345
> ```
>
> **Windows (PowerShell)**
>
> ```powershell
> $env:GEMINI_CLI_IDE_PID=12345
> ```
>
> When this variable is set, Gemini CLI will skip automatic detection and attempt
> to connect using the provided PID.
>
> #### Configuration errors
>
> - **Message:**
>   `🔴 Disconnected: Directory mismatch. Gemini CLI is running in a different location than the open workspace in [IDE Name]. Please run the CLI from one of the following directories: [List of directories]`
>
>   - **Cause:** The CLI's current working directory is outside the workspace you
>     have open in your IDE.
>   - **Solution:** `cd` into the same directory that is open in your IDE and
>     restart the CLI.
>
> - **Message:**
>   `🔴 Disconnected: To use this feature, please open a workspace folder in [IDE Name] and try again.`
>   - **Cause:** You have no workspace open in your IDE.
>   - **Solution:** Open a workspace in your IDE and restart the CLI.
>
> #### General errors
>
> - **Message:**
>   `IDE integration is not supported in your current environment. To use this feature, run Gemini CLI in one of these supported IDEs: [List of IDEs]`
>
>   - **Cause:** You are running Gemini CLI in a terminal or environment that is
>     not a supported IDE.
>   - **Solution:** Run Gemini CLI from the integrated terminal of a supported
>     IDE, like Antigravity or VS Code.
>
> - **Message:**
>   `No installer is available for IDE. Please install Gemini CLI Companion extension manually from the marketplace.`
>   - **Cause:** You ran `/ide install`, but the CLI does not have an automated
>     installer for your specific IDE.
>   - **Solution:** Open your IDE's extension marketplace, search for "Gemini CLI
>     Companion", and
>     [install it manually](#3-manual-installation-from-a-marketplace).
>
> ### ACP integration errors
>
> For issues related to ACP integration, refer to the debugging and telemetry
> section in the [ACP Mode](/docs/cli/acp-mode) documentation.

### Source: ACP Mode — Full page

> ACP (Agent Client Protocol) mode is a special operational mode of Gemini CLI
> designed for programmatic control, primarily for IDE and other developer tool
> integrations. It uses a JSON-RPC protocol over stdio to communicate between
> Gemini CLI agent and a client.
>
> To start Gemini CLI in ACP mode, use the `--acp` flag:
>
> ```bash
> gemini --acp
> ```
>
> ## Agent Client Protocol (ACP)
>
> ACP is an open protocol that standardizes how AI coding agents communicate with
> code editors and IDEs. It addresses the challenge of fragmented distribution,
> where agents traditionally needed custom integrations for each client. With ACP,
> developers can implement their agent once, and it becomes compatible with any
> ACP-compliant editor.
>
> For a comprehensive introduction to ACP, including its architecture and
> benefits, refer to the official
> [ACP Introduction](https://agentclientprotocol.com/get-started/introduction)
> documentation.
>
> ### Existing integrations using ACP
>
> The ACP Agent Registry simplifies the distribution and management of
> ACP-compatible agents across various IDEs. Gemini CLI is an ACP-compatible agent
> and can be found in this registry.
>
> For more general information about the registry, and how to use it with specific
> IDEs like JetBrains and Zed, refer to the
> [IDE Integration](/docs/ide-integration) documentation.
>
> You can also find more information on the official
> [ACP Agent Registry](https://agentclientprotocol.com/get-started/registry) page.
>
> ## Architecture and protocol basics
>
> ACP mode establishes a client-server relationship between your tool (the client)
> and Gemini CLI (the server).
>
> - **Communication:** The entire communication happens over standard input/output
>   (stdio) using the JSON-RPC 2.0 protocol.
> - **Client's role:** The client is responsible for sending requests (for
>   example, prompts) and handling responses and notifications from Gemini CLI.
> - **Gemini CLI's role:** In ACP mode, Gemini CLI listens for incoming JSON-RPC
>   requests, processes them, and sends back responses.
>
> The core of the ACP implementation can be found in
> `packages/cli/src/acp/acpClient.ts`.
>
> ### Extending with MCP
>
> ACP can be used with the Model Context Protocol (MCP). This lets an ACP client
> (like an IDE) expose its own functionality as "tools" that the Gemini model can
> use.
>
> 1.  The client implements an **MCP server** that advertises its tools.
> 2.  During the ACP `initialize` handshake, the client provides the connection
>     details for its MCP server.
> 3.  Gemini CLI connects to the MCP server, discovers the available tools, and
>     makes them available to the AI model.
> 4.  When the model decides to use one of these tools, Gemini CLI sends a tool
>     call request to the MCP server.
>
> This mechanism lets for a powerful, two-way integration where the agent can
> leverage the IDE's capabilities to perform tasks. The MCP client logic is in
> `packages/core/src/tools/mcp-client.ts`.
>
> ## Capabilities and supported methods
>
> The ACP protocol exposes a number of methods for ACP clients (for example IDEs)
> to control Gemini CLI.
>
> ### Core methods
>
> - `initialize`: Establishes the initial connection and lets the client to
>   register its MCP server.
> - `authenticate`: Authenticates the user.
> - `newSession`: Starts a new chat session.
> - `loadSession`: Loads a previous session.
> - `prompt`: Sends a prompt to the agent.
> - `cancel`: Cancels an ongoing prompt.
>
> ### Session control
>
> - `setSessionMode`: Allows changing the approval level for tool calls (for
>   example, to `auto-approve`).
> - `unstable_setSessionModel`: Changes the model for the current session.
>
> ### File system proxy
>
> ACP includes a proxied file system service. This means that when the agent needs
> to read or write files, it does so through the ACP client. This is a security
> feature that ensures the agent only has access to the files that the client (and
> by extension, the user) has explicitly allowed.
>
> ## Debugging and telemetry
>
> You can get insights into the ACP communication and the agent's behavior through
> debugging logs and telemetry.
>
> ### Debugging logs
>
> To enable general debugging logs, start Gemini CLI with the `--debug` flag:
>
> ```bash
> gemini --acp --debug
> ```
>
> ### Telemetry
>
> For more detailed telemetry, you can use the following environment variables to
> capture telemetry data to a file:
>
> - `GEMINI_TELEMETRY_ENABLED=true`
> - `GEMINI_TELEMETRY_TARGET=local`
> - `GEMINI_TELEMETRY_OUTFILE=/path/to/your/log.json`
>
> This will write a JSON log file containing detailed information about all the
> events happening within the agent, including ACP requests and responses. The
> integration test `integration-tests/acp-telemetry.test.ts` provides a working
> example of how to set this up.

### Source: Gemini CLI companion plugin: Interface specification — Full page

> > Last Updated: September 15, 2025
>
> This document defines the contract for building a companion plugin to enable
> Gemini CLI's IDE mode. For VS Code, these features (native diffing, context
> awareness) are provided by the official extension
> ([marketplace](https://marketplace.visualstudio.com/items?itemName=Google.gemini-cli-vscode-ide-companion)).
> This specification is for contributors who wish to bring similar functionality
> to other editors like JetBrains IDEs, Sublime Text, etc.
>
> ## I. The communication interface
>
> Gemini CLI and the IDE plugin communicate through a local communication channel.
>
> ### 1. Transport layer: MCP over HTTP
>
> The plugin **MUST** run a local HTTP server that implements the **Model Context
> Protocol (MCP)**.
>
> - **Protocol:** The server must be a valid MCP server. We recommend using an
>   existing MCP SDK for your language of choice if available.
> - **Endpoint:** The server should expose a single endpoint (for example, `/mcp`)
>   for all MCP communication.
> - **Port:** The server **MUST** listen on a dynamically assigned port (that is,
>   listen on port `0`).
>
> ### 2. Discovery mechanism: The port file
>
> For Gemini CLI to connect, it needs to discover which IDE instance it's running
> in and what port your server is using. The plugin **MUST** facilitate this by
> creating a "discovery file."
>
> - **How the CLI finds the file:** The CLI determines the Process ID (PID) of the
>   IDE it's running in by traversing the process tree. It then looks for a
>   discovery file that contains this PID in its name.
> - **File location:** The file must be created in a specific directory:
>   `os.tmpdir()/gemini/ide/`. Your plugin must create this directory if it
>   doesn't exist.
> - **File naming convention:** The filename is critical and **MUST** follow the
>   pattern: `gemini-ide-server-${PID}-${PORT}.json`
>   - `${PID}`: The process ID of the parent IDE process. Your plugin must
>     determine this PID and include it in the filename.
>   - `${PORT}`: The port your MCP server is listening on.
> - **File content and workspace validation:** The file **MUST** contain a JSON
>   object with the following structure:
>
>   ```json
>   {
>     "port": 12345,
>     "workspacePath": "/path/to/project1:/path/to/project2",
>     "authToken": "a-very-secret-token",
>     "ideInfo": {
>       "name": "vscode",
>       "displayName": "VS Code"
>     }
>   }
>   ```
>
>   - `port` (number, required): The port of the MCP server.
>   - `workspacePath` (string, required): A list of all open workspace root paths,
>     delimited by the OS-specific path separator (`:` for Linux/macOS, `;` for
>     Windows). The CLI uses this path to ensure it's running in the same project
>     folder that's open in the IDE. If the CLI's current working directory is not
>     a sub-directory of `workspacePath`, the connection will be rejected. Your
>     plugin **MUST** provide the correct, absolute path(s) to the root of the
>     open workspace(s).
>   - `authToken` (string, required): A secret token for securing the connection.
>     The CLI will include this token in an `Authorization: Bearer <token>` header
>     on all requests.
>   - `ideInfo` (object, required): Information about the IDE.
>     - `name` (string, required): A short, lowercase identifier for the IDE (for
>       example, `vscode`, `jetbrains`).
>     - `displayName` (string, required): A user-friendly name for the IDE (for
>       example, `VS Code`, `JetBrains IDE`).
>
> - **Authentication:** To secure the connection, the plugin **MUST** generate a
>   unique, secret token and include it in the discovery file. The CLI will then
>   include this token in the `Authorization` header for all requests to the MCP
>   server (for example, `Authorization: Bearer a-very-secret-token`). Your server
>   **MUST** validate this token on every request and reject any that are
>   unauthorized.
> - **Tie-breaking with environment variables (recommended):** For the most
>   reliable experience, your plugin **SHOULD** both create the discovery file and
>   set the `GEMINI_CLI_IDE_SERVER_PORT` environment variable in the integrated
>   terminal. The file serves as the primary discovery mechanism, but the
>   environment variable is crucial for tie-breaking. If a user has multiple IDE
>   windows open for the same workspace, the CLI uses the
>   `GEMINI_CLI_IDE_SERVER_PORT` variable to identify and connect to the correct
>   window's server.
>
> ## II. The context interface
>
> To enable context awareness, the plugin **MAY** provide the CLI with real-time
> information about the user's activity in the IDE.
>
> ### `ide/contextUpdate` notification
>
> The plugin **MAY** send an `ide/contextUpdate`
> [notification](https://modelcontextprotocol.io/specification/2025-06-18/basic/index#notifications)
> to the CLI whenever the user's context changes.
>
> - **Triggering events:** This notification should be sent (with a recommended
>   debounce of 50ms) when:
>   - A file is opened, closed, or focused.
>   - The user's cursor position or text selection changes in the active file.
> - **Payload (`IdeContext`):** The notification parameters **MUST** be an
>   `IdeContext` object:
>
>   ```typescript
>   interface IdeContext {
>     workspaceState?: {
>       openFiles?: File[];
>       isTrusted?: boolean;
>     };
>   }
>
>   interface File {
>     // Absolute path to the file
>     path: string;
>     // Last focused Unix timestamp (for ordering)
>     timestamp: number;
>     // True if this is the currently focused file
>     isActive?: boolean;
>     cursor?: {
>       // 1-based line number
>       line: number;
>       // 1-based character number
>       character: number;
>     };
>     // The text currently selected by the user
>     selectedText?: string;
>   }
>   ```
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > The `openFiles` list should only include files that exist on disk.
> > Virtual files (for example, unsaved files without a path, editor settings pages)
> > **MUST** be excluded.
>
> ### How the CLI uses this context
>
> After receiving the `IdeContext` object, the CLI performs several normalization
> and truncation steps before sending the information to the model.
>
> - **File ordering:** The CLI uses the `timestamp` field to determine the most
>   recently used files. It sorts the `openFiles` list based on this value.
>   Therefore, your plugin **MUST** provide an accurate Unix timestamp for when a
>   file was last focused.
> - **Active file:** The CLI considers only the most recent file (after sorting)
>   to be the "active" file. It will ignore the `isActive` flag on all other files
>   and clear their `cursor` and `selectedText` fields. Your plugin should focus
>   on setting `isActive: true` and providing cursor/selection details only for
>   the currently focused file.
> - **Truncation:** To manage token limits, the CLI truncates both the file list
>   (to 10 files) and the `selectedText` (to 16KB).
>
> While the CLI handles the final truncation, it is highly recommended that your
> plugin also limits the amount of context it sends.
>
> ## III. The diffing interface
>
> To enable interactive code modifications, the plugin **MAY** expose a diffing
> interface. This allows the CLI to request that the IDE open a diff view, showing
> proposed changes to a file. The user can then review, edit, and ultimately
> accept or reject these changes directly within the IDE.
>
> ### `openDiff` tool
>
> The plugin **MUST** register an `openDiff` tool on its MCP server.
>
> - **Description:** This tool instructs the IDE to open a modifiable diff view
>   for a specific file.
> - **Request (`OpenDiffRequest`):** The tool is invoked via a `tools/call`
>   request. The `arguments` field within the request's `params` **MUST** be an
>   `OpenDiffRequest` object.
>
>   ```typescript
>   interface OpenDiffRequest {
>     // The absolute path to the file to be diffed.
>     filePath: string;
>     // The proposed new content for the file.
>     newContent: string;
>   }
>   ```
>
> - **Response (`CallToolResult`):** The tool **MUST** immediately return a
>   `CallToolResult` to acknowledge the request and report whether the diff view
>   was successfully opened.
>
>   - On Success: If the diff view was opened successfully, the response **MUST**
>     contain empty content (that is, `content: []`).
>   - On Failure: If an error prevented the diff view from opening, the response
>     **MUST** have `isError: true` and include a `TextContent` block in the
>     `content` array describing the error.
>
>   The actual outcome of the diff (acceptance or rejection) is communicated
>   asynchronously via notifications.
>
> ### `closeDiff` tool
>
> The plugin **MUST** register a `closeDiff` tool on its MCP server.
>
> - **Description:** This tool instructs the IDE to close an open diff view for a
>   specific file.
> - **Request (`CloseDiffRequest`):** The tool is invoked via a `tools/call`
>   request. The `arguments` field within the request's `params` **MUST** be an
>   `CloseDiffRequest` object.
>
>   ```typescript
>   interface CloseDiffRequest {
>     // The absolute path to the file whose diff view should be closed.
>     filePath: string;
>   }
>   ```
>
> - **Response (`CallToolResult`):** The tool **MUST** return a `CallToolResult`.
>   - On Success: If the diff view was closed successfully, the response **MUST**
>     include a single **TextContent** block in the content array containing the
>     file's final content before closing.
>   - On Failure: If an error prevented the diff view from closing, the response
>     **MUST** have `isError: true` and include a `TextContent` block in the
>     `content` array describing the error.
>
> ### `ide/diffAccepted` notification
>
> When the user accepts the changes in a diff view (for example, by clicking an
> "Apply" or "Save" button), the plugin **MUST** send an `ide/diffAccepted`
> notification to the CLI.
>
> - **Payload:** The notification parameters **MUST** include the file path and
>   the final content of the file. The content may differ from the original
>   `newContent` if the user made manual edits in the diff view.
>
>   ```typescript
>   {
>     // The absolute path to the file that was diffed.
>     filePath: string;
>     // The full content of the file after acceptance.
>     content: string;
>   }
>   ```
>
> ### `ide/diffRejected` notification
>
> When the user rejects the changes (for example, by closing the diff view without
> accepting), the plugin **MUST** send an `ide/diffRejected` notification to the
> CLI.
>
> - **Payload:** The notification parameters **MUST** include the file path of the
>   rejected diff.
>
>   ```typescript
>   {
>     // The absolute path to the file that was diffed.
>     filePath: string;
>   }
>   ```
>
> ## IV. The lifecycle interface
>
> The plugin **MUST** manage its resources and the discovery file correctly based
> on the IDE's lifecycle.
>
> - **On activation (IDE startup/plugin enabled):**
>   1.  Start the MCP server.
>   2.  Create the discovery file.
> - **On deactivation (IDE shutdown/plugin disabled):**
>   1.  Stop the MCP server.
>   2.  Delete the discovery file.
