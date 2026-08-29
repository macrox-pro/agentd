---
primary_sources:
  - id: T1-HOOKS
    title: "Gemini CLI hooks"
    url: "https://geminicli.com/docs/hooks.md"
    section: "Precedence"
  - id: T1-CONFIG
    title: "Gemini CLI configuration"
    url: "https://geminicli.com/docs/reference/configuration.md"
    section: "Locations"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Project vs global path matrix

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Gemini CLI hooks — Hook precedence

> Hooks are scripts or programs that Gemini CLI executes at specific points in the
> agentic loop, allowing you to intercept and customize behavior without modifying
> the CLI's source code.
>
> ## What are hooks?
>
> Hooks run synchronously as part of the agent loop—when a hook event fires,
> Gemini CLI waits for all matching hooks to complete before continuing.
>
> With hooks, you can:
>
> - **Add context:** Inject relevant information (like git history) before the
>   model processes a request.
> - **Validate actions:** Review tool arguments and block potentially dangerous
>   operations.
> - **Enforce policies:** Implement security scanners and compliance checks.
> - **Log interactions:** Track tool usage and model responses for auditing.
> - **Optimize behavior:** Dynamically filter available tools or adjust model
>   parameters.
>
> ### Getting started
>
> - **[Writing hooks guide](/docs/hooks/writing-hooks)**: A tutorial on creating
>   your first hook with comprehensive examples.
> - **[Best practices](/docs/hooks/best-practices)**: Guidelines on security,
>   performance, and debugging.
> - **[Hooks reference](/docs/hooks/reference)**: The definitive technical
>   specification of I/O schemas and exit codes.
>
> ## Core concepts
>
> ### Hook events
>
> Hooks are triggered by specific events in Gemini CLI's lifecycle.
>
> | Event                 | When It Fires                                  | Impact                 | Common Use Cases                             |
> | --------------------- | ---------------------------------------------- | ---------------------- | -------------------------------------------- |
> | `SessionStart`        | When a session begins (startup, resume, clear) | Inject Context         | Initialize resources, load context           |
> | `SessionEnd`          | When a session ends (exit, clear)              | Advisory               | Clean up, save state                         |
> | `BeforeAgent`         | After user submits prompt, before planning     | Block Turn / Context   | Add context, validate prompts, block turns   |
> | `AfterAgent`          | When agent loop ends                           | Retry / Halt           | Review output, force retry or halt execution |
> | `BeforeModel`         | Before sending request to LLM                  | Block Turn / Mock      | Modify prompts, swap models, mock responses  |
> | `AfterModel`          | After receiving LLM response                   | Block Turn / Redact    | Filter/redact responses, log interactions    |
> | `BeforeToolSelection` | Before LLM selects tools                       | Filter Tools           | Filter available tools, optimize selection   |
> | `BeforeTool`          | Before a tool executes                         | Block Tool / Rewrite   | Validate arguments, block dangerous ops      |
> | `AfterTool`           | After a tool executes                          | Block Result / Context | Process results, run tests, hide results     |
> | `PreCompress`         | Before context compression                     | Advisory               | Save state, notify user                      |
> | `Notification`        | When a system notification occurs              | Advisory               | Forward to desktop alerts, logging           |
>
> ### Global mechanics
>
> Understanding these core principles is essential for building robust hooks.
>
> #### Strict JSON requirements (The "Golden Rule")
>
> Hooks communicate via `stdin` (Input) and `stdout` (Output).
>
> 1. **Silence is Mandatory**: Your script **must not** print any plain text to
>    `stdout` other than the final JSON object. **Even a single `echo` or `print`
>    call before the JSON will break parsing.**
> 2. **Pollution = Failure**: If `stdout` contains non-JSON text, parsing will
>    fail. The CLI will default to "Allow" and treat the entire output as a
>    `systemMessage`.
> 3. **Debug via Stderr**: Use `stderr` for **all** logging and debugging (for
>    example, `echo "debug" >&2`). Gemini CLI captures `stderr` but never attempts
>    to parse it as JSON.
>
> #### Exit codes
>
> Gemini CLI uses exit codes to determine the high-level outcome of a hook
> execution:
>
> | Exit Code | Label            | Behavioral Impact                                                                                                                                                            |
> | --------- | ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | **0**     | **Success**      | The `stdout` is parsed as JSON. **Preferred code** for all logic, including intentional blocks (for example, `{"decision": "deny"}`).                                        |
> | **2**     | **System Block** | **Critical Block**. The target action (tool, turn, or stop) is aborted. `stderr` is used as the rejection reason. High severity; used for security stops or script failures. |
> | **Other** | **Warning**      | Non-fatal failure. A warning is shown, but the interaction proceeds using original parameters.                                                                               |
>
> #### Matchers
>
> You can filter which specific tools or triggers fire your hook using the
> `matcher` field.
>
> - **Tool events** (`BeforeTool`, `AfterTool`): Matchers are **Regular
>   Expressions**. (for example, `"write_.*"`).
> - **Lifecycle events**: Matchers are **Exact Strings**. (for example,
>   `"startup"`).
> - **Wildcards**: `"*"` or `""` (empty string) matches all occurrences.
>
> ## Configuration
>
> Hooks are configured in `settings.json`. Gemini CLI merges configurations from
> multiple layers in the following order of precedence (highest to lowest):
>
> 1.  **Project settings**: `.gemini/settings.json` in the current directory.
> 2.  **User settings**: `~/.gemini/settings.json`.
> 3.  **System settings**: `/etc/gemini-cli/settings.json`.
> 4.  **Extensions**: Hooks defined by installed extensions.
>
> ### Configuration schema
>
> ```json
> {
>   "hooks": {
>     "BeforeTool": [
>       {
>         "matcher": "write_file|replace",
>         "hooks": [
>           {
>             "name": "security-check",
>             "type": "command",
>             "command": "$GEMINI_PROJECT_DIR/.gemini/hooks/security.sh",
>             "timeout": 5000
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> #### Hook configuration fields
>
> | Field         | Type   | Required  | Description                                                          |
> | :------------ | :----- | :-------- | :------------------------------------------------------------------- |
> | `type`        | string | **Yes**   | The execution engine. Currently only `"command"` is supported.       |
> | `command`     | string | **Yes\*** | The shell command to execute. (Required when `type` is `"command"`). |
> | `name`        | string | No        | A friendly name for identifying the hook in logs and CLI commands.   |
> | `timeout`     | number | No        | Execution timeout in milliseconds (default: 60000).                  |
> | `description` | string | No        | A brief explanation of the hook's purpose.                           |
>
> ---
>
> ### Environment variables
>
> Hooks are executed with a sanitized environment.
>
> - `GEMINI_PROJECT_DIR`: The absolute path to the project root.
> - `GEMINI_PLANS_DIR`: The absolute path to the plans directory.
> - `GEMINI_SESSION_ID`: The unique ID for the current session.
> - `GEMINI_CWD`: The current working directory.
> - `CLAUDE_PROJECT_DIR`: (Alias) Provided for compatibility.
>
> ## Security and risks
>
> <!-- prettier-ignore -->
> > [!WARNING]
> > Hooks execute arbitrary code with your user privileges. By
> > configuring hooks, you are allowing scripts to run shell commands on your
> > machine.
>
> **Project-level hooks** are particularly risky when opening untrusted projects.
> Gemini CLI **fingerprints** project hooks. If a hook's name or command changes
> (for example, via `git pull`), it is treated as a **new, untrusted hook** and
> you will be warned before it executes.
>
> See [Security Considerations](/docs/hooks/best-practices#using-hooks-securely)
> for a detailed threat model.
>
> ## Managing hooks
>
> Use the CLI commands to manage hooks without editing JSON manually:
>
> - **View hooks:** `/hooks panel`
> - **Enable/Disable all:** `/hooks enable-all` or `/hooks disable-all`
> - **Toggle individual:** `/hooks enable <name>` or `/hooks disable <name>`

### Source: Gemini CLI configuration — Configuration locations

> Gemini CLI offers several ways to configure its behavior, including environment
> variables, command-line arguments, and settings files. This document outlines
> the different configuration methods and available settings.
>
> ## Configuration layers
>
> Configuration is applied in the following order of precedence (lower numbers are
> overridden by higher numbers):
>
> 1.  **Default values:** Hardcoded defaults within the application.
> 2.  **System defaults file:** System-wide default settings that can be
>     overridden by other settings files.
> 3.  **User settings file:** Global settings for the current user.
> 4.  **Project settings file:** Project-specific settings.
> 5.  **System settings file:** System-wide settings that override all other
>     settings files.
> 6.  **Environment variables:** System-wide or session-specific variables,
>     potentially loaded from `.env` files.
> 7.  **Command-line arguments:** Values passed when launching the CLI.
>
> ## Settings files
>
> Gemini CLI uses JSON settings files for persistent configuration. There are four
> locations for these files:
>
> <!-- prettier-ignore -->
> > [!TIP]
> > JSON-aware editors can use autocomplete and validation by pointing to
> > the generated schema at `schemas/settings.schema.json` in this repository.
> > When working outside the repo, reference the hosted schema at
> > `https://raw.githubusercontent.com/google-gemini/gemini-cli/main/schemas/settings.schema.json`.
>
> - **System defaults file:**
>   - **Location:** `/etc/gemini-cli/system-defaults.json` (Linux),
>     `C:\ProgramData\gemini-cli\system-defaults.json` (Windows) or
>     `/Library/Application Support/GeminiCli/system-defaults.json` (macOS). The
>     path can be overridden using the `GEMINI_CLI_SYSTEM_DEFAULTS_PATH`
>     environment variable.
>   - **Scope:** Provides a base layer of system-wide default settings. These
>     settings have the lowest precedence and are intended to be overridden by
>     user, project, or system override settings.
> - **User settings file:**
>   - **Location:** `~/.gemini/settings.json` (where `~` is your home directory).
>   - **Scope:** Applies to all Gemini CLI sessions for the current user. User
>     settings override system defaults.
> - **Project settings file:**
>   - **Location:** `.gemini/settings.json` within your project's root directory.
>   - **Scope:** Applies only when running Gemini CLI from that specific project.
>     Project settings override user settings and system defaults.
> - **System settings file:**
>   - **Location:** `/etc/gemini-cli/settings.json` (Linux),
>     `C:\ProgramData\gemini-cli\settings.json` (Windows) or
>     `/Library/Application Support/GeminiCli/settings.json` (macOS). The path can
>     be overridden using the `GEMINI_CLI_SYSTEM_SETTINGS_PATH` environment
>     variable.
>   - **Scope:** Applies to all Gemini CLI sessions on the system, for all users.
>     System settings act as overrides, taking precedence over all other settings
>     files. May be useful for system administrators at enterprises to have
>     controls over users' Gemini CLI setups.
>
> **Note on environment variables in settings:** String values within your
> `settings.json` and `gemini-extension.json` files can reference environment
> variables using `$VAR_NAME`, `${VAR_NAME}`, or `${VAR_NAME:-DEFAULT_VALUE}`
> syntax. These variables will be automatically resolved when the settings are
> loaded. For example, if you have an environment variable `MY_API_TOKEN`, you
> could use it in `settings.json` like this: `"apiKey": "$MY_API_TOKEN"`. If you
> want to provide a fallback value, use `${MY_API_TOKEN:-default-token}`.
> Additionally, each extension can have its own `.env` file in its directory,
> which will be loaded automatically.
>
> **Note for Enterprise Users:** For guidance on deploying and managing Gemini CLI
> in a corporate environment, see the
> [Enterprise Configuration](/docs/cli/enterprise) documentation.
>
> ### The `.gemini` directory in your project
>
> In addition to a project settings file, a project's `.gemini` directory can
> contain other project-specific files related to Gemini CLI's operation, such as:
>
> - [Custom sandbox profiles](#sandboxing) (for example,
>   `.gemini/sandbox-macos-custom.sb`, `.gemini/sandbox.Dockerfile`).
>
> ### Available settings in `settings.json`
>
> Settings are organized into categories. All settings should be placed within
> their corresponding top-level category object in your `settings.json` file.
>
> <!-- SETTINGS-AUTOGEN:START -->
>
> #### `policyPaths`
>
> - **`policyPaths`** (array):
>   - **Description:** Additional policy files or directories to load.
>   - **Default:** `[]`
>   - **Requires restart:** Yes
>
> #### `adminPolicyPaths`
>
> - **`adminPolicyPaths`** (array):
>   - **Description:** Additional admin policy files or directories to load.
>   - **Default:** `[]`
>   - **Requires restart:** Yes
>
> #### `general`
>
> - **`general.preferredEditor`** (enum):
>
>   - **Description:** The preferred editor to open files in. Must be one of the
>     built-in supported identifiers. Use /editor in the CLI to pick
>     interactively, or leave unset to use $VISUAL/$EDITOR.
>   - **Default:** `undefined`
>   - **Values:** `"vscode"`, `"vscodium"`, `"windsurf"`, `"cursor"`, `"zed"`,
>     `"antigravity"`, `"sublimetext"`, `"lapce"`, `"nova"`, `"bbedit"`, `"vim"`,
>     `"neovim"`, `"emacs"`, `"hx"`, `"emacsclient"`, `"micro"`
>
> - **`general.openEditorInNewWindow`** (boolean):
>
>   - **Description:** Open VS Code-family editors in a new window when editing
>     files.
>   - **Default:** `false`
>
> - **`general.vimMode`** (boolean):
>
>   - **Description:** Enable Vim keybindings
>   - **Default:** `false`
>
> - **`general.defaultApprovalMode`** (enum):
>
>   - **Description:** The default approval mode for tool execution. 'default'
>     prompts for approval, 'auto_edit' auto-approves edit tools, and 'plan' is
>     read-only mode. YOLO mode (auto-approve all actions) can only be enabled via
>     command line (--yolo or --approval-mode=yolo).
>   - **Default:** `"default"`
>   - **Values:** `"default"`, `"auto_edit"`, `"plan"`
>
> - **`general.devtools`** (boolean):
>
>   - *
