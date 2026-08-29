---
primary_sources:
  - id: T1-CONFIG
    title: "Gemini CLI configuration"
    url: "https://geminicli.com/docs/reference/configuration.md"
    section: "Directory layout"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Gemini directory layout

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Gemini CLI configuration — Settings paths and scopes

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
>   - **Description:** Enable DevTools inspector on launch.
>   - **Default:** `false`
>
> - **`general.enableAutoUpdate`** (boolean):
>
>   - **Description:** Enable automatic updates.
>   - **Default:** `true`
>
> - **`general.enableAutoUpdateNotification`** (boolean):
>
>   - **Description:** Enable update notification prompts.
>   - **Default:** `true`
>
> - **`general.enableNotifications`** (boolean):
>
>   - **Description:** Enable terminal run-event notifications for action-required
>     prompts and session completion.
>   - **Default:** `false`
>
> - **`general.notificationMethod`** (enum):
>
>   - **Description:** How to send terminal notifications.
>   - **Default:** `"auto"`
>   - **Values:** `"auto"`, `"osc9"`, `"osc777"`, `"bell"`
>
> - **`general.checkpointing.enabled`** (boolean):
>
>   - **Description:** Enable session checkpointing for recovery
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`general.plan.enabled`** (boolean):
>
>   - **Description:** Enable Plan Mode for read-only safety during planning.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`general.plan.directory`** (string):
>
>   - **Description:** The directory where planning artifacts are stored. If not
>     specified, defaults to the system temporary directory. A custom directory
>     requires a policy to allow write access in Plan Mode.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`general.plan.modelRouting`** (boolean):
>
>   - **Description:** Automatically switch between Pro and Flash models based on
>     Plan Mode status. Uses Pro for the planning phase and Flash for the
>     implementation phase.
>   - **Default:** `true`
>
> - **`general.retryFetchErrors`** (boolean):
>
>   - **Description:** Retry on "exception TypeError: fetch failed sending
>     request" errors.
>   - **Default:** `true`
>
> - **`general.maxAttempts`** (number):
>
>   - **Description:** Maximum number of attempts for requests to the main chat
>     model. Cannot exceed 10.
>   - **Default:** `10`
>
> - **`general.debugKeystrokeLogging`** (boolean):
>
>   - **Description:** Enable debug logging of keystrokes to the console.
>   - **Default:** `false`
>
> - **`general.sessionRetention.enabled`** (boolean):
>
>   - **Description:** Enable automatic session cleanup
>   - **Default:** `true`
>
> - **`general.sessionRetention.maxAge`** (string):
>
>   - **Description:** Automatically delete chats older than this time period
>     (e.g., "30d", "7d", "24h", "1w")
>   - **Default:** `"30d"`
>
> - **`general.sessionRetention.maxCount`** (number):
>
>   - **Description:** Alternative: Maximum number of sessions to keep (most
>     recent)
>   - **Default:** `undefined`
>
> - **`general.sessionRetention.minRetention`** (string):
>
>   - **Description:** Minimum retention period (safety limit, defaults to "1d")
>   - **Default:** `"1d"`
>
> - **`general.topicUpdateNarration`** (boolean):
>
>   - **Description:** Enable the Topic & Update communication model for reduced
>     chattiness and structured progress reporting.
>   - **Default:** `true`
>
> - **`general.logRagSnippets`** (boolean):
>   - **Description:** Log full Code Customization (RAG) retrieved snippets to a
>     local file for debugging.
>   - **Default:** `false`
>
> #### `output`
>
> - **`output.format`** (enum):
>   - **Description:** The format of the CLI output. Can be `text` or `json`.
>   - **Default:** `"text"`
>   - **Values:** `"text"`, `"json"`
>
> #### `ui`
>
> - **`ui.debugRainbow`** (boolean):
>
>   - **Description:** Enable debug rainbow rendering. Only useful for debugging
>     rendering bugs and performance issues.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`ui.theme`** (string):
>
>   - **Description:** The color theme for the UI. See the CLI themes guide for
>     available options.
>   - **Default:** `undefined`
>
> - **`ui.autoThemeSwitching`** (boolean):
>
>   - **Description:** Automatically switch between default light and dark themes
>     based on terminal background color.
>   - **Default:** `true`
>
> - **`ui.terminalBackgroundPollingInterval`** (number):
>
>   - **Description:** Interval in seconds to poll the termi
