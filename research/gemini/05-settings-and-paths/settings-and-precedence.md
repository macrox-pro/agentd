---
primary_sources:
  - id: T1-SETTINGS
    title: "Gemini CLI settings (`/settings` command)"
    url: "https://geminicli.com/docs/cli/settings.md"
    section: "Full page"
  - id: T1-CONFIG
    title: "Gemini CLI configuration"
    url: "https://geminicli.com/docs/reference/configuration.md"
    section: "Settings file locations"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Settings and precedence

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Gemini CLI settings (`/settings` command) — Full page

> Control your Gemini CLI experience with the `/settings` command. The `/settings`
> command opens a dialog to view and edit all your Gemini CLI settings, including
> your UI experience, keybindings, and accessibility features.
>
> Your Gemini CLI settings are stored in a `settings.json` file. In addition to
> using the `/settings` command, you can also edit them in one of the following
> locations:
>
> - **User settings**: `~/.gemini/settings.json`
> - **Workspace settings**: `your-project/.gemini/settings.json`
>
> <!-- prettier-ignore -->
> > [!IMPORTANT]
> > Workspace settings override user settings.
>
> ## Settings reference
>
> Here is a list of all the available settings, grouped by category and ordered as
> they appear in the UI.
>
> <!-- SETTINGS-AUTOGEN:START -->
>
> ### General
>
> | UI Label                      | Setting                            | Description                                                                                                                                                                                                                                                   | Default     |
> | ----------------------------- | ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
> | Vim Mode                      | `general.vimMode`                  | Enable Vim keybindings                                                                                                                                                                                                                                        | `false`     |
> | Default Approval Mode         | `general.defaultApprovalMode`      | The default approval mode for tool execution. 'default' prompts for approval, 'auto_edit' auto-approves edit tools, and 'plan' is read-only mode. YOLO mode (auto-approve all actions) can only be enabled via command line (--yolo or --approval-mode=yolo). | `"default"` |
> | Enable Auto Update            | `general.enableAutoUpdate`         | Enable automatic updates.                                                                                                                                                                                                                                     | `true`      |
> | Enable Terminal Notifications | `general.enableNotifications`      | Enable terminal run-event notifications for action-required prompts and session completion.                                                                                                                                                                   | `false`     |
> | Terminal Notification Method  | `general.notificationMethod`       | How to send terminal notifications.                                                                                                                                                                                                                           | `"auto"`    |
> | Enable Plan Mode              | `general.plan.enabled`             | Enable Plan Mode for read-only safety during planning.                                                                                                                                                                                                        | `true`      |
> | Plan Directory                | `general.plan.directory`           | The directory where planning artifacts are stored. If not specified, defaults to the system temporary directory. A custom directory requires a policy to allow write access in Plan Mode.                                                                     | `undefined` |
> | Plan Model Routing            | `general.plan.modelRouting`        | Automatically switch between Pro and Flash models based on Plan Mode status. Uses Pro for the planning phase and Flash for the implementation phase.                                                                                                          | `true`      |
> | Retry Fetch Errors            | `general.retryFetchErrors`         | Retry on "exception TypeError: fetch failed sending request" errors.                                                                                                                                                                                          | `true`      |
> | Max Chat Model Attempts       | `general.maxAttempts`              | Maximum number of attempts for requests to the main chat model. Cannot exceed 10.                                                                                                                                                                             | `10`        |
> | Debug Keystroke Logging       | `general.debugKeystrokeLogging`    | Enable debug logging of keystrokes to the console.                                                                                                                                                                                                            | `false`     |
> | Enable Session Cleanup        | `general.sessionRetention.enabled` | Enable automatic session cleanup                                                                                                                                                                                                                              | `true`      |
> | Keep chat history             | `general.sessionRetention.maxAge`  | Automatically delete chats older than this time period (e.g., "30d", "7d", "24h", "1w")                                                                                                                                                                       | `"30d"`     |
> | Topic & Update Narration      | `general.topicUpdateNarration`     | Enable the Topic & Update communication model for reduced chattiness and structured progress reporting.                                                                                                                                                       | `true`      |
> | Log RAG Snippets              | `general.logRagSnippets`           | Log full Code Customization (RAG) retrieved snippets to a local file for debugging.                                                                                                                                                                           | `false`     |
>
> ### Output
>
> | UI Label      | Setting         | Description                                            | Default  |
> | ------------- | --------------- | ------------------------------------------------------ | -------- |
> | Output Format | `output.format` | The format of the CLI output. Can be `text` or `json`. | `"text"` |
>
> ### UI
>
> | UI Label                             | Setting                                | Description                                                                                                                                                       | Default |
> | ------------------------------------ | -------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
> | Auto Theme Switching                 | `ui.autoThemeSwitching`                | Automatically switch between default light and dark themes based on terminal background color.                                                                    | `true`  |
> | Terminal Background Polling Interval | `ui.terminalBackgroundPollingInterval` | Interval in seconds to poll the terminal background color.                                                                                                        | `60`    |
> | Hide Window Title                    | `ui.hideWindowTitle`                   | Hide the window title bar                                                                                                                                         | `false` |
> | Inline Thinking                      | `ui.inlineThinkingMode`                | Display model thinking inline: off or full.                                                                                                                       | `"off"` |
> | Show Thoughts in Title               | `ui.showStatusInTitle`                 | Show Gemini CLI model thoughts in the terminal window title during the working phase                                                                              | `false` |
> | Dynamic Window Title                 | `ui.dynamicWindowTitle`                | Update the terminal window title with current status icons (Ready: ◇, Action Required: ✋, Working: ✦)                                                            | `true`  |
> | Show Home Directory Warning          | `ui.showHomeDirectoryWarning`          | Show a warning when running Gemini CLI in the home directory.                                                                                                     | `true`  |
> | Show Compatibility Warnings          | `ui.showCompatibilityWarnings`         | Show warnings about terminal or OS compatibility issues.                                                                                                          | `true`  |
> | Hide Tips                            | `ui.hideTips`                          | Hide helpful tips in the UI                                                                                                                                       | `false` |
> | Escape Pasted @ Symbols              | `ui.escapePastedAtSymbols`             | When enabled, @ symbols in pasted text are escaped to prevent unintended @path expansion.                                                                         | `false` |
> | Show Shortcuts Hint                  | `ui.showShortcutsHint`                 | Show the "? for shortcuts" hint above the input.                                                                                                                  | `true`  |
> | Compact Tool Output                  | `ui.compactToolOutput`                 | Display tool outputs (like directory listings and file reads) in a compact, structured format.                                                                    | `true`  |
> | Hide Banner                          | `ui.hideBanner`                        | Hide the application banner                                                                                                                                       | `false` |
> | Hide Context Summary                 | `ui.hideContextSummary`                | Hide the context summary (GEMINI.md, MCP servers) above the input.                                                                                                | `false` |
> | Hide CWD                             | `ui.footer.hideCWD`                    | Hide the current working directory in the footer.                                                                                                                 | `false` |
> | Hide Sandbox Status                  | `ui.footer.hideSandboxStatus`          | Hide the sandbox status indicator in the footer.                                                                                                                  | `false` |
> | Hide Model Info                      | `ui.footer.hideModelInfo`              | Hide the model name and context usage in the footer.                                                                                                              | `false` |
> | Hide Context Window Percentage       | `ui.footer.hideContextPercentage`      | Hides the context window usage percentage.                                                                                                                        | `true`  |
> | Hide Footer                          | `ui.hideFooter`                        | Hide the footer from the UI                                                                                                                                       | `false` |
> | Show Memory Usage                    | `ui.showMemoryUsage`                   | Display memory usage information in the UI                                                                                                                        | `false` |
> | Show Line Numbers                    | `ui.showLineNumbers`                   | Show line numbers in the chat.                                                                                                                                    | `true`  |
> | Show Citations                       | `ui.showCitations`                     | Show citations for generated text in the chat.                                                                                                                    | `false` |
> | Show Model Info In Chat              | `ui.showModelInfoInChat`               | Show the model name in the chat for each model turn.                                                                                                              | `false` |
> | Show User Identity                   | `ui.showUserIdentity`                  | Show the signed-in user's identity (e.g. email) in the UI.                                                                                                        | `true`  |
> | Use Alternate Screen Buffer          | `ui.useAlternateBuffer`                | Use an alternate screen buffer for the UI, preserving shell history.                                                                                              | `false` |
> | Render Process                       | `ui.renderProcess`                     | Enable Ink render process for the UI.                                                                                                                             | `true`  |
> | Terminal Buffer                      | `ui.terminalBuffer`                    | Use the new terminal buffer architecture for rendering.                                                                                                           | `false` |
> | Use Background Color                 | `ui.useBackgroundColor`                | Whether to use background colors in the UI.                                                                                                                       | `true`  |
> | Incremental Rendering                | `ui.incrementalRendering`              | Enable incremental rendering for the UI. This option will reduce flickering but may cause rendering artifacts. Only supported when useAlternateBuffer is enabled. | `true`  |
> | Show Spinner                         | `ui.showSpinner`                       | Show the spinner during operations.                                                                                                                               | `true`  |
> | Loading Phrases                      | `ui.loadingPhrases`                    | What to show while the model is working: tips, witty comments, all, or off.                                                                                       | `"off"` |
> | Error Verbosity                      | `ui.errorVerbosity`                    | Controls whether recoverable errors are hidden (low) or fully shown (full).                                                                                       | `"low"` |
> | Screen Reader Mode                   | `ui.accessibility.screenReader`        | Render output in plain-text to be more screen reader accessible                                                                                                   | `false` |
>
> ### IDE
>
> | UI Label | Setting       | Description                  | Default |
> | -------- | ------------- | ---------------------------- | ------- |
> | IDE Mode | `ide.enabled` | Enable IDE integration mode. | `false` |
>
> ### Billing
>
> | UI Label         | Setting                   | Description                                                                                                                                                | Default |
> | ---------------- | ------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
> | Overage Strategy | `billing.overageStrategy` | How to handle quota exhaustion when AI credits are available. 'ask' prompts each time, 'always' automatically uses credits, 'never' disables credit usage. | `"ask"` |
>
> ### Model
>
> | UI Label                      | Setting                      | Description                                                                            | Default     |
> | ----------------------------- | ---------------------------- | -------------------------------------------------------------------------------------- | ----------- |
> | Model                         | `model.name`                 | The Gemini model to use for conversations.                                             | `undefined` |
> | Max Session Turns             | `model.maxSessionTurns`      | Maximum number of user/model/tool turns to keep in a session. -1 means unlimited.      | `-1`        |
> | Context Compression Threshold | `model.compressionThreshold` | The fraction of context usage at which to trigger context compression (e.g. 0.2, 0.3). | `0.5`       |
> | Disable Loop Detection        | `model.disableLoopDetection` | Disable automatic detection and prevention of infinite loops.                          | `false`     |
> | Skip Next Speaker Check       | `model.skipNextSpeakerCheck` | Skip the next speaker check.                                                           | `true`      |
>
> ### Agents
>
> | UI Label                  | Setting                                  | Description                                                                                   | Default |
> | ------------------------- | ---------------------------------------- | --------------------------------------------------------------------------------------------- | ------- |
> | Confirm Sensitive Actions | `agents.browser.confirmSensitiveActions` | Require manual confirmation for sensitive browser actions (e.g., fill_form, evaluate_script). | `false` |
> | Block File Uploads        | `agents.browser.blockFileUploads`        | Hard-block file upload requests from the browser agent.                                       | `false` |
>
> ### Context
>
> | UI Label                             | Setting                                           | Description                                                                                                                                                                                                                                 | Default |
> | ------------------------------------ | ------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
> | Memory Discovery Max Dirs            | `context.discoveryMaxDirs`                        | Maximum number of directories to search for memory.                                                                                                                                                                                         | `200`   |
> | Load Memory From Include Directories | `context.loadMemoryFromIncludeDirectories`        | Controls how /memory reload loads GEMINI.md files. When true, include directories are scanned; when false, only the current directory is used.                                                                                              | `false` |
> | Respect .gitignore                   | `context.fileFiltering.respectGitIgnore`          | Respect .gitignore files when searching.                                                                                                                                                                                                    | `true`  |
> | Respect .geminiignore                | `context.fileFiltering.respectGeminiIgnore`       | Respect .geminiignore files when searching.                                                                                                                                                                                                 | `true`  |
> | Enable Recursive File Search         | `context.fileFiltering.enableRecursiveFileSearch` | Enable recursive file search functionality when completing @ references in the prompt.                                                                                                                                                      | `true`  |
> | Enable Fuzzy Search                  | `context.fileFiltering.enableFuzzySearch`         | Enable fuzzy search when searching for files.                                                                                                                                                                                               | `true`  |
> | Custom Ignore File Paths             | `context.fileFiltering.customIgnoreFilePaths`     | Additional ignore file paths to respect. These files take precedence over .geminiignore and .gitignore. Files earlier in the array take precedence over files later in the array, e.g. the first file takes precedence over the second one. | `[]`    |
>
> ### Tools
>
> | UI Label                         | Setting                              | Description                                                                                                                                                                | Default |
> | -------------------------------- | ------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
> | Sandbox Allowed Paths            | `tools.sandboxAllowedPaths`          | List of additional paths that the sandbox is allowed to access.                                                                                                            | `[]`    |
> | Sandbox Network Access           | `tools.sandboxNetworkAccess`         | Whether the sandbox is allowed to access the network.                                                                                                                      | `false` |
> | Enable Interactive Shell         | `tools.shell.enableInteractiveShell` | Use node-pty for an interactive shell experience. Fallback to child_process still applies.                                                                                 | `true`  |
> | Show Color                       | `tools.shell.showColor`              | Show color in shell output.                                                                                                                                                | `true`  |
> | Use Ripgrep                      | `tools.useRipgrep`                   | Use ripgrep for file content search instead of the fallback implementation. Provides faster search performance.                                                            | `true`  |
> | Tool Output Truncation Threshold | `tools.truncateToolOutputThreshold`  | Maximum characters to show when truncating large tool outputs. Set to 0 or negative to disable truncation.                                                                 | `40000` |
> | Disable LLM Correction           | `tools.disableLLMCorrection`         | Disable LLM-based error correction for edit tools. When enabled, tools will fail immediately if exact string matches are not found, instead of attempting to self-correct. | `true`  |
>
> ### Security
>
> | UI Label                              | Setting                                         | Description                                                                                                                                                                                                                          | Default |
> | ------------------------------------- | ----------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------- |
> | Tool Sandboxing                       | `security.toolSandboxing`                       | Tool-level sandboxing. Isolates individual tools instead of the entire CLI process.                                                                                                                                                  | `false` |
> | Disable YOLO Mode                     | `security.disableYoloMode`                      | Disable YOLO mode, even if enabled by a flag.                                                                                                                                                                                        | `false` |
> | Disable Always Allow                  | `security.disableAlwaysAllow`                   | Disable "Always allow" options in tool confirmation dialogs.                                                                                                                                                                         | `false` |
> | Allow Permanent Tool Approval         | `security.enablePermanentToolApproval`          | Enable the "Allow for all future sessions" option in tool confirmation dialogs.                                                                                                                                                      | `false` |
> | Auto-add to Policy by Default         | `security.autoAddToPolicyByDefault`             | When enabled, the "Allow for all future sessions" option becomes the default choice for low-risk tools in trusted workspaces.                                                                                                        | `false` |
> | Blocks extensions from Git            | `security.blockGitExtensions`                   | Blocks installing and loading extensions from Git.                                                                                                                                                                                   | `false` |
> | Extension Source Regex Allowlist      | `security.allowedExtensions`                    | List of Regex patterns for allowed extensions. If nonempty, only extensions that match the patterns in this list are allowed. Overrides the blockGitExtensions setting.                                                              | `[]`    |
> | Folder Trust                          | `security.folderTrust.enabled`                  | Setting to track whether Folder trust is enabled.                                                                                                                                                                                    | `true`  |
> | Enable Environment Variable Redaction | `security.environmentVariableRedaction.enabled` | Enable redaction of environment variables that may contain secrets.                                                                                                                                                                  | `false` |
> | Enable Context-Aware Security         | `security.enableConseca`                        | Enable the context-aware security checker. This feature uses an LLM to dynamically generate and enforce security policies for tool use based on your prompt, providing an additional layer of protection against unintended actions. | `false` |
>
> ### Advanced
>
> | UI Label                          | Setting                        | Description                                                                                                                                                                                                           | Default |
> | --------------------------------- | ------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
> | Auto Configure Max Old Space Size | `advanced.autoConfigureMemory` | Automatically configure Node.js memory limits. Note: Because memory is allocated during the initial process boot, this setting is only read from the global user settings file and ignores workspace-level overrides. | `true`  |
> | Ignore Local .env                 | `advanced.ignoreLocalEnv`      | Whether to ignore generic .env files in the project directory.                                                                                                                                                        | `false` |
>
> ### Experimental
>
> | UI Label                                             | Setting                                         | Description                                                                                                                                                                                                                                                            | Default              |
> | ---------------------------------------------------- | ----------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------- |
> | Gemma Models                                         | `experimental.gemma`                            | Enable access to Gemma 4 models via Gemini API.                                                                                                                                                                                                                        | `true`               |
> | Voice Mode                                           | `experimental.voiceMode`                        | Enable experimental voice dictation and commands (/voice, /voice model).                                                                                                                                                                                               | `false`              |
> | Voice Activation Mode                                | `experimental.voice.activationMode`             | How to trigger voice recording with the Space key.                                                                                                                                                                                                                     | `"push-to-talk"`     |
> | Voice Transcription Backend                          | `experimental.voice.backend`                    | The backend to use for voice transcription. Note: When using the Gemini Live backend, voice recordings are sent to Google Cloud for transcription.                                                                                                                     | `"gemini-live"`      |
> | Whisper Model                                        | `experimental.voice.whisperModel`               | The Whisper model to use for local transcription.                                                                                                                                                                                                                      | `"ggml-base.en.bin"` |
> | Voice Stop Grace Period (ms)                         | `experimental.voice.stopGracePeriodMs`          | How long to wait for final transcription after stopping recording.                                                                                                                                                                                                     | `4000`               |
> | Enable Git Worktrees                                 | `experimental.worktrees`                        | Enable automated Git worktree management for parallel work.                                                                                                                                                                                                            | `false`              |
> | Use OSC 52 Paste                                     | `experimental.useOSC52Paste`                    | Use OSC 52 for pasting. This may be more robust than the default system when using remote terminal sessions (if your terminal is configured to allow it).                                                                                                              | `false`              |
> | Use OSC 52 Copy                                      | `experimental.useOSC52Copy`                     | Use OSC 52 for copying. This may be more robust than the default system when using remote terminal sessions (if your terminal is configured to allow it).                                                                                                              | `false`              |
> | Model Steering                                       | `experimental.modelSteering`                    | Enable model steering (user hints) to guide the model during tool execution.                                                                                                                                                                                           | `false`              |
> | Direct Web Fetch                                     | `experimental.directWebFetch`                   | Enable web fetch behavior that bypasses LLM summarization.                                                                                                                                                                                                             | `false`              |
> | Enable Gemma Model Router                            | `experimental.gemmaModelRouter.enabled`         | Enable the Gemma Model Router (experimental). Requires a local endpoint serving Gemma via the Gemini API using LiteRT-LM shim.                                                                                                                                         | `false`              |
> | Auto-start LiteRT Server                             | `experimental.gemmaModelRouter.autoStartServer` | Automatically start the LiteRT-LM server when Gemini CLI starts and the Gemma router is enabled.                                                                                                                                                                       | `false`              |
> | Auto Memory                                          | `experimental.autoMemory`                       | Automatically extract memory patches and skills from past sessions in the background. Every change is written as a unified diff `.patch` file under `<projectMemoryDir>/.inbox/<kind>/` and held for review in /memory inbox; nothing is applied until you approve it. | `false`              |
> | Use the generalist profile to manage agent contexts. | `experimental.generalistProfile`                | Suitable for general coding and software development tasks.                                                                                                                                                                                                            | `false`              |
> | Enable Context Management                            | `experimental.contextManagement`                | Enable logic for context management.                                                                                                                                                                                                                                   | `false`              |
>
> ### Skills
>
> | UI Label            | Setting          | Description          | Default |
> | ------------------- | ---------------- | -------------------- | ------- |
> | Enable Agent Skills | `skills.enabled` | Enable Agent Skills. | `true`  |
>
> ### HooksConfig
>
> | UI Label           | Setting                     | Description                                                                      | Default |
> | ------------------ | --------------------------- | -------------------------------------------------------------------------------- | ------- |
> | Enable Hooks       | `hooksConfig.enabled`       | Canonical toggle for the hooks system. When disabled, no hooks will be executed. | `true`  |
> | Hook Notifications | `hooksConfig.notifications` | Show visual indicators when hooks are executing.                                 | `true`  |
>
> <!-- SETTINGS-AUTOGEN:END -->

### Source: Gemini CLI configuration — Settings file locations (excerpt)

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
