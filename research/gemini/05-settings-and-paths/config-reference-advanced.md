---
primary_sources:
  - id: T1-CONFIG
    title: "Gemini CLI configuration"
    url: "https://geminicli.com/docs/reference/configuration.md"
    section: "Remaining categories"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Configuration reference — advanced

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Gemini CLI configuration — Part 3

> axTokens`** (number):
>
>   - **Description:** The number of tokens to allow before triggering
>     compression.
>   - **Default:** `150000`
>   - **Requires restart:** Yes
>
> - **`contextManagement.historyWindow.retainedTokens`** (number):
>
>   - **Description:** The number of tokens to always retain.
>   - **Default:** `40000`
>   - **Requires restart:** Yes
>
> - **`contextManagement.messageLimits.normalMaxTokens`** (number):
>
>   - **Description:** The target number of tokens to budget for a normal
>     conversation turn.
>   - **Default:** `2500`
>   - **Requires restart:** Yes
>
> - **`contextManagement.messageLimits.retainedMaxTokens`** (number):
>
>   - **Description:** The maximum number of tokens a single conversation turn can
>     consume before truncation.
>   - **Default:** `12000`
>   - **Requires restart:** Yes
>
> - **`contextManagement.messageLimits.normalizationHeadRatio`** (number):
>
>   - **Description:** The ratio of tokens to retain from the beginning of a
>     truncated message (0.0 to 1.0).
>   - **Default:** `0.25`
>   - **Requires restart:** Yes
>
> - **`contextManagement.tools.distillation.maxOutputTokens`** (number):
>
>   - **Description:** Maximum tokens to show to the model when truncating large
>     tool outputs.
>   - **Default:** `10000`
>   - **Requires restart:** Yes
>
> - **`contextManagement.tools.distillation.summarizationThresholdTokens`**
>   (number):
>
>   - **Description:** Threshold above which truncated tool outputs will be
>     summarized by an LLM.
>   - **Default:** `20000`
>   - **Requires restart:** Yes
>
> - **`contextManagement.tools.outputMasking.protectionThresholdTokens`**
>   (number):
>
>   - **Description:** Minimum number of tokens to protect from masking (most
>     recent tool outputs).
>   - **Default:** `50000`
>   - **Requires restart:** Yes
>
> - **`contextManagement.tools.outputMasking.minPrunableThresholdTokens`**
>   (number):
>
>   - **Description:** Minimum prunable tokens required to trigger a masking pass.
>   - **Default:** `30000`
>   - **Requires restart:** Yes
>
> - **`contextManagement.tools.outputMasking.protectLatestTurn`** (boolean):
>   - **Description:** Ensures the absolute latest turn is never masked,
>     regardless of token count.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> #### `admin`
>
> - **`admin.secureModeEnabled`** (boolean):
>
>   - **Description:** If true, disallows YOLO mode and "Always allow" options
>     from being used.
>   - **Default:** `false`
>
> - **`admin.extensions.enabled`** (boolean):
>
>   - **Description:** If false, disallows extensions from being installed or
>     used.
>   - **Default:** `true`
>
> - **`admin.mcp.enabled`** (boolean):
>
>   - **Description:** If false, disallows MCP servers from being used.
>   - **Default:** `true`
>
> - **`admin.mcp.config`** (object):
>
>   - **Description:** Admin-configured MCP servers (allowlist).
>   - **Default:** `{}`
>
> - **`admin.mcp.requiredConfig`** (object):
>
>   - **Description:** Admin-required MCP servers that are always injected.
>   - **Default:** `{}`
>
> - **`admin.skills.enabled`** (boolean):
>   - **Description:** If false, disallows agent skills from being used.
>   - **Default:** `true`
>   <!-- SETTINGS-AUTOGEN:END -->
>
> #### `mcpServers`
>
> Configures connections to one or more Model-Context Protocol (MCP) servers for
> discovering and using custom tools. Gemini CLI attempts to connect to each
> configured MCP server to discover available tools. Every discovered tool is
> prepended with the `mcp_` prefix and its server alias to form a fully qualified
> name (FQN) (for example, `mcp_serverAlias_actualToolName`) to avoid conflicts.
> Note that the system might strip certain schema properties from MCP tool
> definitions for compatibility. At least one of `command`, `url`, or `httpUrl`
> must be provided. If multiple are specified, the order of precedence is
> `httpUrl`, then `url`, then `command`.
>
> <!-- prettier-ignore -->
> > [!WARNING]
> > Avoid using underscores (`_`) in your server aliases (for example, use
> > `my-server` instead of `my_server`). The underlying policy engine parses Fully
> > Qualified Names (`mcp_server_tool`) using the first underscore after the
> > `mcp_` prefix. An underscore in your server alias will cause the parser to
> > misidentify the server name, which can cause security policies to fail
> > silently.
>
> - **`mcpServers.<SERVER_NAME>`** (object): The server parameters for the named
>   server.
>   - `command` (string, optional): The command to execute to start the MCP server
>     via standard I/O.
>   - `args` (array of strings, optional): Arguments to pass to the command.
>   - `env` (object, optional): Environment variables to set for the server
>     process.
>   - `cwd` (string, optional): The working directory in which to start the
>     server.
>   - `url` (string, optional): The URL of an MCP server that uses Server-Sent
>     Events (SSE) for communication.
>   - `httpUrl` (string, optional): The URL of an MCP server that uses streamable
>     HTTP for communication.
>   - `headers` (object, optional): A map of HTTP headers to send with requests to
>     `url` or `httpUrl`.
>   - `timeout` (number, optional): Timeout in milliseconds for requests to this
>     MCP server.
>   - `trust` (boolean, optional): Trust this server and bypass all tool call
>     confirmations.
>   - `description` (string, optional): A brief description of the server, which
>     may be used for display purposes.
>   - `includeTools` (array of strings, optional): List of tool names to include
>     from this MCP server. When specified, only the tools listed here will be
>     available from this server (allowlist behavior). If not specified, all tools
>     from the server are enabled by default.
>   - `excludeTools` (array of strings, optional): List of tool names to exclude
>     from this MCP server. Tools listed here will not be available to the model,
>     even if they are exposed by the server. **Note:** `excludeTools` takes
>     precedence over `includeTools` - if a tool is in both lists, it will be
>     excluded.
>
> #### `telemetry`
>
> Configures logging and metrics collection for Gemini CLI. For more information,
> see [Telemetry](/docs/cli/telemetry).
>
> - **Properties:**
>   - **`enabled`** (boolean): Whether or not telemetry is enabled.
>   - **`traces`** (boolean): Whether detailed traces with large attributes (like
>     tool outputs and file reads) are captured. Defaults to `false`.
>   - **`target`** (string): The destination for collected telemetry. Supported
>     values are `local` and `gcp`.
>   - **`otlpEndpoint`** (string): The endpoint for the OTLP Exporter.
>   - **`otlpProtocol`** (string): The protocol for the OTLP Exporter (`grpc` or
>     `http`).
>   - **`logPrompts`** (boolean): Whether or not to include the content of user
>     prompts in the logs.
>   - **`outfile`** (string): The file to write telemetry to when `target` is
>     `local`.
>   - **`useCollector`** (boolean): Whether to use an external OTLP collector.
>
> ### Example `settings.json`
>
> Here is an example of a `settings.json` file with the nested structure, new as
> of v0.3.0:
>
> ```json
> {
>   "general": {
>     "vimMode": true,
>     "preferredEditor": "code",
>     "sessionRetention": {
>       "enabled": true,
>       "maxAge": "30d",
>       "maxCount": 100
>     }
>   },
>   "ui": {
>     "theme": "GitHub",
>     "hideBanner": true,
>     "hideTips": false,
>     "customWittyPhrases": [
>       "You forget a thousand things every day. Make sure this is one of ’em",
>       "Connecting to AGI"
>     ]
>   },
>   "tools": {
>     "sandbox": "docker",
>     "discoveryCommand": "bin/get_tools",
>     "callCommand": "bin/call_tool",
>     "exclude": ["write_file"]
>   },
>   "mcpServers": {
>     "mainServer": {
>       "command": "bin/mcp_server.py"
>     },
>     "anotherServer": {
>       "command": "node",
>       "args": ["mcp_server.js", "--verbose"]
>     }
>   },
>   "telemetry": {
>     "enabled": true,
>     "target": "local",
>     "otlpEndpoint": "http://localhost:4317",
>     "logPrompts": true
>   },
>   "privacy": {
>     "usageStatisticsEnabled": true
>   },
>   "model": {
>     "name": "gemini-1.5-pro-latest",
>     "maxSessionTurns": 10,
>     "summarizeToolOutput": {
>       "run_shell_command": {
>         "tokenBudget": 100
>       }
>     }
>   },
>   "context": {
>     "fileName": ["CONTEXT.md", "GEMINI.md"],
>     "includeDirectories": ["path/to/dir1", "~/path/to/dir2", "../path/to/dir3"],
>     "loadFromIncludeDirectories": true,
>     "fileFiltering": {
>       "respectGitIgnore": false
>     }
>   },
>   "advanced": {
>     "excludedEnvVars": ["DEBUG", "DEBUG_MODE", "NODE_ENV"]
>   }
> }
> ```
>
> ## Shell history
>
> The CLI keeps a history of shell commands you run. To avoid conflicts between
> different projects, this history is stored in a project-specific directory
> within your user's home folder.
>
> - **Location:** `~/.gemini/tmp/<project_hash>/shell_history`
>   - `<project_hash>` is a unique identifier generated from your project's root
>     path.
>   - The history is stored in a file named `shell_history`.
>
> ## Environment variables and `.env` files
>
> Environment variables are a common way to configure applications, especially for
> sensitive information like API keys or for settings that might change between
> environments. For authentication setup, see the
> [Authentication documentation](/docs/get-started/authentication) which covers
> all available authentication methods.
>
> The CLI automatically loads environment variables from an `.env` file. The
> loading order is:
>
> 1.  `.env` file in the current working directory.
> 2.  If not found, it searches upwards in parent directories until it finds an
>     `.env` file or reaches the project root (identified by a `.git` folder) or
>     the home directory.
> 3.  If still not found, it looks for `~/.env` (in the user's home directory).
>
> **Environment variable exclusion:** Some environment variables (like `DEBUG` and
> `DEBUG_MODE`) are automatically excluded from being loaded from project `.env`
> files to prevent interference with gemini-cli behavior. Variables from
> `.gemini/.env` files are never excluded. You can customize this behavior using
> the `advanced.excludedEnvVars` setting in your `settings.json` file.
>
> - **`GEMINI_API_KEY`**:
>   - Your API key for the Gemini API.
>   - One of several available
>     [authentication methods](/docs/get-started/authentication).
>   - Set this in your shell profile (for example, `~/.bashrc`, `~/.zshrc`) or an
>     `.env` file.
> - **`GEMINI_MODEL`**:
>   - Specifies the default Gemini model to use.
>   - Overrides the hardcoded default
>   - Example: `export GEMINI_MODEL="gemini-3-flash-preview"` (Windows PowerShell:
>     `$env:GEMINI_MODEL="gemini-3-flash-preview"`)
> - **`GEMINI_CLI_TRUST_WORKSPACE`**:
>   - If set to `"true"`, trusts the current workspace for the duration of the
>     session, bypassing the folder trust check.
>   - Useful for headless environments (for example, CI/CD pipelines).
> - **`GEMINI_CLI_TRUSTED_FOLDERS_PATH`**:
>   - Overrides the default location for the `trustedFolders.json` file.
>   - Useful if you want to store this configuration in a custom location instead
>     of the default `~/.gemini/`.
> - **`GEMINI_CLI_IDE_PID`**:
>   - Manually specifies the PID of the IDE process to use for integration. This
>     is useful when running Gemini CLI in a standalone terminal while still
>     wanting to associate it with a specific IDE instance.
>   - Overrides the automatic IDE detection logic.
> - **`GEMINI_CLI_HOME`**:
>   - Specifies the root directory for Gemini CLI's user-level configuration and
>     storage.
>   - By default, this is the user's system home directory. The CLI will create a
>     `.gemini` folder inside this directory.
>   - Useful for shared compute environments or keeping CLI state isolated.
>   - Example: `export GEMINI_CLI_HOME="/path/to/user/config"` (Windows
>     PowerShell: `$env:GEMINI_CLI_HOME="C:\path\to\user\config"`)
> - **`GEMINI_CLI_SURFACE`**:
>   - Specifies a custom label to include in the `User-Agent` header for API
>     traffic reporting.
>   - This is useful for tracking specific internal tools or distribution
>     channels.
>   - Example: `export GEMINI_CLI_SURFACE="my-custom-tool"` (Windows PowerShell:
>     `$env:GEMINI_CLI_SURFACE="my-custom-tool"`)
> - **`GOOGLE_API_KEY`**:
>   - Your Google Cloud API key.
>   - Required for using Vertex AI in express mode.
>   - Ensure you have the necessary permissions.
>   - Example: `export GOOGLE_API_KEY="YOUR_GOOGLE_API_KEY"` (Windows PowerShell:
>     `$env:GOOGLE_API_KEY="YOUR_GOOGLE_API_KEY"`).
> - **`GOOGLE_CLOUD_PROJECT`**:
>   - Your Google Cloud Project ID.
>   - Required for using Code Assist or Vertex AI.
>   - If using Vertex AI, ensure you have the necessary permissions in this
>     project.
>   - **Cloud Shell note:** When running in a Cloud Shell environment, this
>     variable defaults to a special project allocated for Cloud Shell users. If
>     you have `GOOGLE_CLOUD_PROJECT` set in your global environment in Cloud
>     Shell, it will be overridden by this default. To use a different project in
>     Cloud Shell, you must define `GOOGLE_CLOUD_PROJECT` in a `.env` file.
>   - Example: `export GOOGLE_CLOUD_PROJECT="YOUR_PROJECT_ID"` (Windows
>     PowerShell: `$env:GOOGLE_CLOUD_PROJECT="YOUR_PROJECT_ID"`).
> - **`GOOGLE_APPLICATION_CREDENTIALS`** (string):
>   - **Description:** The path to your Google Application Credentials JSON file.
>   - **Example:**
>     `export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/credentials.json"`
>     (Windows PowerShell:
>     `$env:GOOGLE_APPLICATION_CREDENTIALS="C:\path\to\your\credentials.json"`)
> - **`GOOGLE_GENAI_API_VERSION`**:
>   - Specifies the API version to use for Gemini API requests.
>   - When set, overrides the default API version used by the SDK.
>   - Example: `export GOOGLE_GENAI_API_VERSION="v1"` (Windows PowerShell:
>     `$env:GOOGLE_GENAI_API_VERSION="v1"`)
> - **`GOOGLE_GEMINI_BASE_URL`**:
>   - Overrides the default base URL for Gemini API requests (when using
>     `gemini-api-key` authentication).
>   - Must be a valid URL. For security, it must use HTTPS unless pointing to
>     `localhost` (or `127.0.0.1` / `[::1]`).
>   - Example: `export GOOGLE_GEMINI_BASE_URL="https://my-proxy.com"` (Windows
>     PowerShell: `$env:GOOGLE_GEMINI_BASE_URL="https://my-proxy.com"`)
> - **`GOOGLE_VERTEX_BASE_URL`**:
>   - Overrides the default base URL for Vertex AI API requests (when using
>     `vertex-ai` authentication).
>   - Must be a valid URL. For security, it must use HTTPS unless pointing to
>     `localhost` (or `127.0.0.1` / `[::1]`).
>   - Example: `export GOOGLE_VERTEX_BASE_URL="https://my-vertex-proxy.com"`
>     (Windows PowerShell:
>     `$env:GOOGLE_VERTEX_BASE_URL="https://my-vertex-proxy.com"`)
> - **`OTLP_GOOGLE_CLOUD_PROJECT`**:
>   - Your Google Cloud Project ID for Telemetry in Google Cloud
>   - Example: `export OTLP_GOOGLE_CLOUD_PROJECT="YOUR_PROJECT_ID"` (Windows
>     PowerShell: `$env:OTLP_GOOGLE_CLOUD_PROJECT="YOUR_PROJECT_ID"`).
> - **`GEMINI_TELEMETRY_ENABLED`**:
>   - Set to `true` or `1` to enable telemetry. Any other value is treated as
>     disabling it.
>   - Overrides the `telemetry.enabled` setting.
> - **`GEMINI_TELEMETRY_TRACES_ENABLED`**:
>   - Set to `true` or `1` to enable detailed tracing with large attributes. Any
>     other value is treated as disabling it.
>   - Overrides the `telemetry.traces` setting.
> - **`GEMINI_TELEMETRY_TARGET`**:
>   - Sets the telemetry target (`local` or `gcp`).
>   - Overrides the `telemetry.target` setting.
> - **`GEMINI_TELEMETRY_OTLP_ENDPOINT`**:
>   - Sets the OTLP endpoint for telemetry.
>   - Overrides the `telemetry.otlpEndpoint` setting.
> - **`GEMINI_TELEMETRY_OTLP_PROTOCOL`**:
>   - Sets the OTLP protocol (`grpc` or `http`).
>   - Overrides the `telemetry.otlpProtocol` setting.
> - **`GEMINI_TELEMETRY_LOG_PROMPTS`**:
>   - Set to `true` or `1` to enable or disable logging of user prompts. Any other
>     value is treated as disabling it.
>   - Overrides the `telemetry.logPrompts` setting.
> - **`GEMINI_TELEMETRY_OUTFILE`**:
>   - Sets the file path to write telemetry to when the target is `local`.
>   - Overrides the `telemetry.outfile` setting.
> - **`GEMINI_TELEMETRY_USE_COLLECTOR`**:
>   - Set to `true` or `1` to enable or disable using an external OTLP collector.
>     Any other value is treated as disabling it.
>   - Overrides the `telemetry.useCollector` setting.
> - **`GOOGLE_CLOUD_LOCATION`**:
>   - Your Google Cloud Project Location (for example, us-central1).
>   - Required for using Vertex AI in non-express mode.
>   - Example: `export GOOGLE_CLOUD_LOCATION="YOUR_PROJECT_LOCATION"` (Windows
>     PowerShell: `$env:GOOGLE_CLOUD_LOCATION="YOUR_PROJECT_LOCATION"`).
> - **`GEMINI_SANDBOX`**:
>   - Alternative to the `sandbox` setting in `settings.json`.
>   - Accepts `true`, `false`, `docker`, `podman`, or a custom command string.
> - **`GEMINI_SYSTEM_MD`**:
>   - Replaces the built‑in system prompt with content from a Markdown file.
>   - `true`/`1`: Use project default path `./.gemini/system.md`.
>   - Any other string: Treat as a path (relative/absolute supported, `~`
>     expands).
>   - `false`/`0` or unset: Use the built‑in prompt. See
>     [System Prompt Override](/docs/cli/system-prompt).
> - **`GEMINI_WRITE_SYSTEM_MD`**:
>   - Writes the current built‑in system prompt to a file for review.
>   - `true`/`1`: Write to `./.gemini/system.md`. Otherwise treat the value as a
>     path.
>   - Run the CLI once with this set to generate the file.
> - **`SEATBELT_PROFILE`** (macOS specific):
>   - Switches the Seatbelt (`sandbox-exec`) profile on macOS.
>   - `permissive-open`: (Default) Denies operations by default, confining writes
>     to the project folder (and a few other folders, see
>     `packages/cli/src/utils/sandbox-macos-permissive-open.sb`) while allowing
>     broad file reads and network access.
>   - `restrictive-open`: Declines operations by default, allows network.
>   - `strict-open`: Restricts both reads and writes to the working directory,
>     allows network.
>   - `strict-proxied`: Same as `strict-open` but routes network through proxy.
>   - `<profile_name>`: Uses a custom profile. To define a custom profile, create
>     a file named `sandbox-macos-<profile_name>.sb` in your project's `.gemini/`
>     directory (for example, `my-project/.gemini/sandbox-macos-custom.sb`).
> - **`DEBUG` or `DEBUG_MODE`** (often used by underlying libraries or the CLI
>   itself):
>   - Set to `true` or `1` to enable verbose debug logging, which can be helpful
>     for troubleshooting.
>   - **Note:** These variables are automatically excluded from project `.env`
>     files by default to prevent interference with gemini-cli behavior. Use
>     `.gemini/.env` files if you need to set these for gemini-cli specifically.
> - **`NO_COLOR`**:
>   - Set to any value to disable all color output in the CLI.
> - **`CLI_TITLE`**:
>   - Set to a string to customize the title of the CLI.
> - **`CODE_ASSIST_ENDPOINT`**:
>   - Specifies the endpoint for the code assist server.
>   - This is useful for development and testing.
>
> ### Environment variable redaction
>
> To prevent accidental leakage of sensitive information, Gemini CLI automatically
> redacts potential secrets from environment variables when executing tools (such
> as shell commands). This "best effort" redaction applies to variables inherited
> from the system or loaded from `.env` files.
>
> **Default Redaction Rules:**
>
> - **By Name:** Variables are redacted if their names contain sensitive terms
>   like `TOKEN`, `SECRET`, `PASSWORD`, `KEY`, `AUTH`, `CREDENTIAL`, `PRIVATE`, or
>   `CERT`.
> - **By Value:** Variables are redacted if their values match known secret
>   patterns, such as:
>   - Private keys (RSA, OpenSSH, PGP, etc.)
>   - Certificates
>   - URLs containing credentials
>   - API keys and tokens (GitHub, Google, AWS, Stripe, Slack, etc.)
> - **Specific Blocklist:** Certain variables like `CLIENT_ID`, `DB_URI`,
>   `DATABASE_URL`, and `CONNECTION_STRING` are always redacted by default.
>
> **Allowlist (Never Redacted):**
>
> - Common system variables (for example, `PATH`, `HOME`, `USER`, `SHELL`, `TERM`,
>   `LANG`).
> - Variables starting with `GEMINI_CLI_`.
> - GitHub Action specific variables.
>
> **Configuration:**
>
> You can customize this behavior in your `settings.json` file:
>
> - **`security.allowedEnvironmentVariables`**: A list of variable names to
>   _never_ redact, even if they match sensitive patterns.
> - **`security.blockedEnvironmentVariables`**: A list of variable names to
>   _always_ redact, even if they don't match sensitive patterns.
>
> ```json
> {
>   "security": {
>     "allowedEnvironmentVariables": ["MY_PUBLIC_KEY", "NOT_A_SECRET_TOKEN"],
>     "blockedEnvironmentVariables": ["INTERNAL_IP_ADDRESS"]
>   }
> }
> ```
>
> ## Command-line arguments
>
> Arguments passed directly when running the CLI can override other configurations
> for that specific session.
>
> - **`--acp`**:
>   - Starts the agent in Agent Communication Protocol (ACP) mode.
> - **`--allowed-mcp-server-names`**:
>   - A comma-separated list of MCP server names to allow for the session.
> - **`--allowed-tools <tool1,tool2,...>`**:
>   - A comma-separated list of tool names that will bypass the confirmation
>     dialog.
>   - Example: `gemini --allowed-tools "ShellTool(git status)"`
> - **`--approval-mode <mode>`**:
>   - Sets the approval mode for tool calls. Available modes:
>     - `default`: Prompt for approval on each tool call (default behavior)
>     - `auto_edit`: Automatically approve edit tools (replace, write_file) while
>       prompting for others
>     - `yolo`: Automatically approve all tool calls (equivalent to `--yolo`)
>     - `plan`: Read-only mode for tool calls (requires experimental planning to
>       be enabled).
>       > **Note:** This mode is currently under development and not yet fully
>       > functional.
>   - Cannot be used together with `--yolo`. Use `--approval-mode=yolo` instead of
>     `--yolo` for the new unified approach.
>   - Example: `gemini --approval-mode auto_edit`
> - **`--debug`** (**`-d`**):
>   - Enables debug mode for this session, providing more verbose output. Open the
>     debug console with F12 to see the additional logging.
> - **`--delete-session <identifier>`**:
>   - Delete a specific chat session by its index number or full session UUID.
>   - Use `--list-sessions` first to see available sessions, their indices, and
>     UUIDs.
>   - Example: `gemini --delete-session 3` or
>     `gemini --delete-session a1b2c3d4-e5f6-7890-abcd-ef1234567890`
> - **`--extensions <extension_name ...>`** (**`-e <extension_name ...>`**):
>   - Specifies a list of extensions to use for the session. If not provided, all
>     available extensions are used.
>   - Use the special term `gemini -e none` to disable all extensions.
>   - Example: `gemini -e my-extension -e my-other-extension`
> - **`--fake-responses`**:
>   - Path to a file with fake model responses for testing.
> - **`--help`** (or **`-h`**):
>   - Displays help information about command-line arguments.
> - **`--include-directories <dir1,dir2,...>`**:
>   - Includes additional directories in the workspace for multi-directory
>     support.
>   - Can be specified multiple times or as comma-separated values.
>   - 5 directories can be added at maximum.
>   - Example: `--include-directories /path/to/project1,/path/to/project2` or
>     `--include-directories /path/to/project1 --include-directories /path/to/project2`
> - **`--list-extensions`** (**`-l`**):
>   - Lists all available extensions and exits.
> - **`--list-sessions`**:
>   - List all available chat sessions for the current project and exit.
>   - Shows session indices, dates, message counts, and preview of first user
>     message.
>   - Example: `gemini --list-sessions`
> - **`--model <model_name>`** (**`-m <model_name>`**):
>   - Specifies the Gemini model to use for this session.
>   - Example: `npm start -- --model gemini-3-pro-preview`
> - **`--output-format <format>`**:
>   - **Description:** Specifies the format of the CLI output for non-interactive
>     mode.
>   - **Values:**
>     - `text`: (Default) The standard human-readable output.
>     - `json`: A machine-readable JSON output.
>     - `stream-json`: A streaming JSON output that emits real-time events.
>   - **Note:** For structured output and scripting, use the
>     `--output-format json` or `--output-format stream-json` flag.
> - **`--prompt <your_prompt>`** (**`-p <your_prompt>`**):
>   - Used to pass a prompt directly to the command. This invokes Gemini CLI in a
>     non-interactive mode.
> - **`--prompt-interactive <your_prompt>`** (**`-i <your_prompt>`**):
>   - Starts an interactive session with the provided prompt as the initial input.
>   - The prompt is processed within the interactive session, not before it.
>   - Cannot be used when piping input from stdin.
>   - Example: `gemini -i "explain this code"`
> - **`--record-responses`**:
>   - Path to a file to record model responses for testing.
> - **`--resume [session_id]`** (**`-r [session_id]`**):
>   - Resume a previous chat session. Use "latest" for the most recent session,
>     provide a session index number, or provide a full session UUID.
>   - If no session_id is provided, defaults to "latest".
>   - Example: `gemini --resume 5` or `gemini --resume latest` or
>     `gemini --resume a1b2c3d4-e5f6-7890-abcd-ef1234567890` or `gemini --resume`
>   - See [Session Management](/docs/cli/session-management) for more details.
> - **`--sandbox`** (**`-s`**):
>   - Enables sandbox mode for this session.
> - **`--screen-reader`**:
>   - Enables screen reader mode, which adjusts the TUI for better compatibility
>     with screen readers.
> - **`--version`**:
>   - Displays the version of the CLI.
> - **`--yolo`**:
>   - Enables YOLO mode, which automatically approves all tool calls.
>
> ## Context files (hierarchical instructional context)
>
> While not strictly configuration for the CLI's _behavior_, context files
> (defaulting to `GEMINI.md` but configurable via the `context.fileName` setting)
> are crucial for configuring the _instructional context_ (also referred to as
> "memory") provided to the Gemini model. This powerful feature lets you give
> project-specific instructions, coding style guides, or any relevant background
> information to the AI, making its responses more tailored and accurate to your
> needs. The CLI includes UI elements, such as an indicator in the footer showing
> the number of loaded context files, to keep you informed about the active
> context.
>
> - **Purpose:** These Markdown files contain instructions, guidelines, or context
>   that you want the Gemini model to be aware of during your interactions. The
>   system is designed to manage this instructional context hierarchically.
>
> ### Example context file content (for example, `GEMINI.md`)
>
> Here's a conceptual example of what a context file at the root of a TypeScript
> project might contain:
>
> ```markdown
> # Project: My Awesome TypeScript Library
>
> ## General Instructions:
>
> - When generating new TypeScript code, follow the existing coding style.
> - Ensure all new functions and classes have JSDoc comments.
> - Prefer functional programming paradigms where appropriate.
> - All code should be compatible with TypeScript 5.0 and Node.js 20+.
>
> ## Coding Style:
>
> - Use 2 spaces for indentation.
> - Interface names should be prefixed with `I` (for example, `IUserService`).
> - Private class members should be prefixed with an underscore (`_`).
> - Always use strict equality (`===` and `!==`).
>
> ## Specific Component: `src/api/client.ts`
>
> - This file handles all outbound API requests.
> - When adding new API call functions, ensure they include robust error handling
>   and logging.
> - Use the existing `fetchWithRetry` utility for all GET requests.
>
> ## Regarding Dependencies:
>
> - Avoid introducing new external dependencies unless absolutely necessary.
> - If a new dependency is required, state the reason.
> ```
>
> This example demonstrates how you can provide general project context, specific
> coding conventions, and even notes about particular files or components. The
> more relevant and precise your context files are, the better the AI can assist
> you. Project-specific context files are highly encouraged to establish
> conventions and context.
>
> - **Hierarchical loading and precedence:** The CLI implements a sophisticated
>   hierarchical memory system by loading context files (for example, `GEMINI.md`)
>   from several locations. Content from files lower in this list (more specific)
>   typically overrides or supplements content from files higher up (more
>   general). The exact concatenation order and final context can be inspected
>   using the `/memory show` command. The typical loading order is:
>   1.  **Global context file:**
>       - Location: `~/.gemini/<configured-context-filename>` (for example,
>         `~/.gemini/GEMINI.md` in your user home directory).
>       - Scope: Provides default instructions for all your projects.
>   2.  **Project root and ancestors context files:**
>       - Location: The CLI searches for the configured context file in the
>         current working directory and then in each parent directory up to either
>         the project root (identified by a `.git` folder) or your home directory.
>       - Scope: Provides context relevant to the entire project or a significant
>         portion of it.
>   3.  **Sub-directory context files (contextual/local):**
>       - Location: The CLI also scans for the configured context file in
>         subdirectories _below_ the current working directory (respecting common
>         ignore patterns like `node_modules`, `.git`, etc.). The breadth of this
>         search is limited to 200 directories by default, but can be configured
>         with the `context.discoveryMaxDirs` setting in your `settings.json`
>         file.
>       - Scope: Allows for highly specific instructions relevant to a particular
>         component, module, or subsection of your project.
> - **Concatenation and UI indication:** The contents of all found context files
>   are concatenated (with separators indicating their origin and path) and
>   provided as part of the system prompt to the Gemini model. The CLI footer
>   displays the count of loaded context files, giving you a quick visual cue
>   about the active instructional context.
> - **Importing content:** You can modularize your context files by importing
>   other Markdown files using the `@path/to/file.md` syntax. For more details,
>   see the [Memory Import Processor documentation](/docs/reference/memport).
> - **Commands for memory management:**
>   - Use `/memory refresh` to force a re-scan and reload of all context files
>     from all configured locations. This updates the AI's instructional context.
>   - Use `/memory show` to display the combined instructional context currently
>     loaded, allowing you to verify the hierarchy and content being used by the
>     AI.
>   - See the [Commands documentation](/docs/reference/commands#memory) for full details on
>     the `/memory` command and its sub-commands (`show` and `reload`).
>
> By understanding and utilizing these configuration layers and the hierarchical
> nature of context files, you can effectively manage the AI's memory and tailor
> Gemini CLI's responses to your specific needs and projects.
>
> ## Sandboxing
>
> Gemini CLI can execute potentially unsafe operations (like shell commands and
> file modifications) within a sandboxed environment to protect your system.
>
> Sandboxing is disabled by default, but you can enable it in a few ways:
>
> - Using `--sandbox` or `-s` flag.
> - Setting `GEMINI_SANDBOX` environment variable.
> - Sandbox is enabled when using `--yolo` or `--approval-mode=yolo` by default.
>
> By default, it uses a pre-built `gemini-cli-sandbox` Docker image.
>
> For project-specific sandboxing needs, you can create a custom Dockerfile at
> `.gemini/sandbox.Dockerfile` in your project's root directory. This Dockerfile
> can be based on the base sandbox image:
>
> ```dockerfile
> FROM gemini-cli-sandbox
>
> # Add your custom dependencies or configurations here.
> # Note: The base image runs as the non-root 'node' user.
> # You must switch to 'root' to install system packages.
> # For example:
> # USER root
> # RUN apt-get update && apt-get install -y some-package
> # USER node
> # COPY ./my-config /app/my-config
> ```
>
> When `.gemini/sandbox.Dockerfile` exists, you can use `BUILD_SANDBOX`
> environment variable when running Gemini CLI to automatically build the custom
> sandbox image:
>
> ```bash
> BUILD_SANDBOX=1 gemini -s
> ```
>
> Building a custom sandbox with `BUILD_SANDBOX` is only supported when running
> Gemini CLI from source. If you installed the CLI with npm, build the Docker
> image separately and reference that image in your sandbox configuration.
>
> ## Usage statistics
>
> To help us improve Gemini CLI, we collect anonymized usage statistics. This data
> helps us understand how the CLI is used, identify common issues, and prioritize
> new features.
>
> **What we collect:**
>
> - **Tool calls:** We log the names of the tools that are called, whether they
>   succeed or fail, and how long they take to execute. We do not collect the
>   arguments passed to the tools or any data returned by them.
> - **API requests:** We log the Gemini model used for each request, the duration
>   of the request, and whether it was successful. We do not collect the content
>   of the prompts or responses.
> - **Session information:** We collect information about the configuration of the
>   CLI, such as the enabled tools and the approval mode.
>
> **What we DON'T collect:**
>
> - **Personally identifiable information (PII):** We do not collect any personal
>   information, such as your name, email address, or API keys.
> - **Prompt and response content:** We do not log the content of your prompts or
>   the responses from the Gemini model.
> - **File content:** We do not log the content of any files that are read or
>   written by the CLI.
>
> **How to opt out:**
>
> You can opt out of usage statistics collection at any time by setting the
> `usageStatisticsEnabled` property to `false` under the `privacy` category in
> your `settings.json` file:
>
> ```json
> {
>   "privacy": {
>     "usageStatisticsEnabled": false
>   }
> }
> ```
