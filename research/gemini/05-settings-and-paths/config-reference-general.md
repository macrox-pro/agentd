---
primary_sources:
  - id: T1-CONFIG
    title: "Gemini CLI configuration"
    url: "https://geminicli.com/docs/reference/configuration.md"
    section: "General; hooksConfig"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Configuration reference — general

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Gemini CLI configuration — Part 1

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
>   - **Description:** Interval in seconds to poll the terminal background color.
>   - **Default:** `60`
>
> - **`ui.customThemes`** (object):
>
>   - **Description:** Custom theme definitions.
>   - **Default:** `{}`
>
> - **`ui.hideWindowTitle`** (boolean):
>
>   - **Description:** Hide the window title bar
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`ui.inlineThinkingMode`** (enum):
>
>   - **Description:** Display model thinking inline: off or full.
>   - **Default:** `"off"`
>   - **Values:** `"off"`, `"full"`
>
> - **`ui.showStatusInTitle`** (boolean):
>
>   - **Description:** Show Gemini CLI model thoughts in the terminal window title
>     during the working phase
>   - **Default:** `false`
>
> - **`ui.dynamicWindowTitle`** (boolean):
>
>   - **Description:** Update the terminal window title with current status icons
>     (Ready: ◇, Action Required: ✋, Working: ✦)
>   - **Default:** `true`
>
> - **`ui.showHomeDirectoryWarning`** (boolean):
>
>   - **Description:** Show a warning when running Gemini CLI in the home
>     directory.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`ui.showCompatibilityWarnings`** (boolean):
>
>   - **Description:** Show warnings about terminal or OS compatibility issues.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`ui.hideTips`** (boolean):
>
>   - **Description:** Hide helpful tips in the UI
>   - **Default:** `false`
>
> - **`ui.escapePastedAtSymbols`** (boolean):
>
>   - **Description:** When enabled, @ symbols in pasted text are escaped to
>     prevent unintended @path expansion.
>   - **Default:** `false`
>
> - **`ui.showShortcutsHint`** (boolean):
>
>   - **Description:** Show the "? for shortcuts" hint above the input.
>   - **Default:** `true`
>
> - **`ui.compactToolOutput`** (boolean):
>
>   - **Description:** Display tool outputs (like directory listings and file
>     reads) in a compact, structured format.
>   - **Default:** `true`
>
> - **`ui.hideBanner`** (boolean):
>
>   - **Description:** Hide the application banner
>   - **Default:** `false`
>
> - **`ui.hideContextSummary`** (boolean):
>
>   - **Description:** Hide the context summary (GEMINI.md, MCP servers) above the
>     input.
>   - **Default:** `false`
>
> - **`ui.footer.items`** (array):
>
>   - **Description:** List of item IDs to display in the footer. Rendered in
>     order
>   - **Default:** `undefined`
>
> - **`ui.footer.showLabels`** (boolean):
>
>   - **Description:** Display a second line above the footer items with
>     descriptive headers (e.g., /model).
>   - **Default:** `true`
>
> - **`ui.footer.hideCWD`** (boolean):
>
>   - **Description:** Hide the current working directory in the footer.
>   - **Default:** `false`
>
> - **`ui.footer.hideSandboxStatus`** (boolean):
>
>   - **Description:** Hide the sandbox status indicator in the footer.
>   - **Default:** `false`
>
> - **`ui.footer.hideModelInfo`** (boolean):
>
>   - **Description:** Hide the model name and context usage in the footer.
>   - **Default:** `false`
>
> - **`ui.footer.hideContextPercentage`** (boolean):
>
>   - **Description:** Hides the context window usage percentage.
>   - **Default:** `true`
>
> - **`ui.hideFooter`** (boolean):
>
>   - **Description:** Hide the footer from the UI
>   - **Default:** `false`
>
> - **`ui.collapseDrawerDuringApproval`** (boolean):
>
>   - **Description:** Whether to collapse the UI drawer when a tool is awaiting
>     confirmation.
>   - **Default:** `true`
>
> - **`ui.showMemoryUsage`** (boolean):
>
>   - **Description:** Display memory usage information in the UI
>   - **Default:** `false`
>
> - **`ui.showLineNumbers`** (boolean):
>
>   - **Description:** Show line numbers in the chat.
>   - **Default:** `true`
>
> - **`ui.showCitations`** (boolean):
>
>   - **Description:** Show citations for generated text in the chat.
>   - **Default:** `false`
>
> - **`ui.showModelInfoInChat`** (boolean):
>
>   - **Description:** Show the model name in the chat for each model turn.
>   - **Default:** `false`
>
> - **`ui.showUserIdentity`** (boolean):
>
>   - **Description:** Show the signed-in user's identity (e.g. email) in the UI.
>   - **Default:** `true`
>
> - **`ui.useAlternateBuffer`** (boolean):
>
>   - **Description:** Use an alternate screen buffer for the UI, preserving shell
>     history.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`ui.renderProcess`** (boolean):
>
>   - **Description:** Enable Ink render process for the UI.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`ui.terminalBuffer`** (boolean):
>
>   - **Description:** Use the new terminal buffer architecture for rendering.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`ui.useBackgroundColor`** (boolean):
>
>   - **Description:** Whether to use background colors in the UI.
>   - **Default:** `true`
>
> - **`ui.incrementalRendering`** (boolean):
>
>   - **Description:** Enable incremental rendering for the UI. This option will
>     reduce flickering but may cause rendering artifacts. Only supported when
>     useAlternateBuffer is enabled.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`ui.showSpinner`** (boolean):
>
>   - **Description:** Show the spinner during operations.
>   - **Default:** `true`
>
> - **`ui.loadingPhrases`** (enum):
>
>   - **Description:** What to show while the model is working: tips, witty
>     comments, all, or off.
>   - **Default:** `"off"`
>   - **Values:** `"tips"`, `"witty"`, `"all"`, `"off"`
>
> - **`ui.errorVerbosity`** (enum):
>
>   - **Description:** Controls whether recoverable errors are hidden (low) or
>     fully shown (full).
>   - **Default:** `"low"`
>   - **Values:** `"low"`, `"full"`
>
> - **`ui.customWittyPhrases`** (array):
>
>   - **Description:** Custom witty phrases to display during loading. When
>     provided, the CLI cycles through these instead of the defaults.
>   - **Default:** `[]`
>
> - **`ui.accessibility.enableLoadingPhrases`** (boolean):
>
>   - **Description:** @deprecated Use ui.loadingPhrases instead. Enable loading
>     phrases during operations.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`ui.accessibility.screenReader`** (boolean):
>   - **Description:** Render output in plain-text to be more screen reader
>     accessible
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> #### `ide`
>
> - **`ide.enabled`** (boolean):
>
>   - **Description:** Enable IDE integration mode.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`ide.hasSeenNudge`** (boolean):
>   - **Description:** Whether the user has seen the IDE integration nudge.
>   - **Default:** `false`
>
> #### `privacy`
>
> - **`privacy.usageStatisticsEnabled`** (boolean):
>   - **Description:** Enable collection of usage statistics
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> #### `billing`
>
> - **`billing.overageStrategy`** (enum):
>
>   - **Description:** How to handle quota exhaustion when AI credits are
>     available. 'ask' prompts each time, 'always' automatically uses credits,
>     'never' disables credit usage.
>   - **Default:** `"ask"`
>   - **Values:** `"ask"`, `"always"`, `"never"`
>
> - **`billing.vertexAi.requestType`** (enum):
>
>   - **Description:** Sets the X-Vertex-AI-LLM-Request-Type header for Vertex AI
>     requests.
>   - **Default:** `undefined`
>   - **Values:** `"dedicated"`, `"shared"`
>   - **Requires restart:** Yes
>
> - **`billing.vertexAi.sharedRequestType`** (enum):
>   - **Description:** Sets the X-Vertex-AI-LLM-Shared-Request-Type header for
>     Vertex AI requests.
>   - **Default:** `undefined`
>   - **Values:** `"priority"`, `"flex"`
>   - **Requires restart:** Yes
>
> #### `model`
>
> - **`model.name`** (string):
>
>   - **Description:** The Gemini model to use for conversations.
>   - **Default:** `undefined`
>
> - **`model.maxSessionTurns`** (number):
>
>   - **Description:** Maximum number of user/model/tool turns to keep in a
>     session. -1 means unlimited.
>   - **Default:** `-1`
>
> - **`model.summarizeToolOutput`** (object):
>
>   - **Description:** Enables or disables summarization of tool output. Configure
>     per-tool token budgets (for example {"run_shell_command": {"tokenBudget":
>     2000}}). Currently only the run_shell_command tool supports summarization.
>   - **Default:** `undefined`
>
> - **`model.compressionThreshold`** (number):
>
>   - **Description:** The fraction of context usage at which to trigger context
>     compression (e.g. 0.2, 0.3).
>   - **Default:** `0.5`
>   - **Requires restart:** Yes
>
> - **`model.disableLoopDetection`** (boolean):
>
>   - **Description:** Disable automatic detection and prevention of infinite
>     loops.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`model.skipNextSpeakerCheck`** (boolean):
>   - **Description:** Skip the next speaker check.
>   - **Default:** `true`
>
> #### `modelConfigs`
>
> - **`modelConfigs.aliases`** (object):
>
>   - **Description:** Named presets for model configs. Can be used in place of a
>     model name and can inherit from other aliases using an `extends` property.
>   - **Default:**
>
>     ```json
>     {
>       "base": {
>         "modelConfig": {
>           "generateContentConfig": {
>             "temperature": 0,
>             "topP": 1
>           }
>         }
>       },
>       "chat-base": {
>         "extends": "base",
>         "modelConfig": {
>           "generateContentConfig": {
>             "thinkingConfig": {
>               "includeThoughts": true
>             },
>             "temperature": 1,
>             "topP": 0.95,
>             "topK": 64
>           }
>         }
>       },
>       "chat-base-2.5": {
>         "extends": "chat-base",
>         "modelConfig": {
>           "generateContentConfig": {
>             "thinkingConfig": {
>               "thinkingBudget": 8192
>             }
>           }
>         }
>       },
>       "chat-base-3": {
>         "extends": "chat-base",
>         "modelConfig": {
>           "generateContentConfig": {
>             "thinkingConfig": {
>               "thinkingLevel": "HIGH"
>             }
>           }
>         }
>       },
>       "gemini-3-pro-preview": {
>         "extends": "chat-base-3",
>         "modelConfig": {
>           "model": "gemini-3-pro-preview"
>         }
>       },
>       "gemini-3-flash-preview": {
>         "extends": "chat-base-3",
>         "modelConfig": {
>           "model": "gemini-3-flash-preview"
>         }
>       },
>       "gemini-3.1-pro-preview": {
>         "extends": "chat-base-3",
>         "modelConfig": {
>           "model": "gemini-3.1-pro-preview"
>         }
>       },
>       "gemini-3.1-pro-preview-customtools": {
>         "extends": "chat-base-3",
>         "modelConfig": {
>           "model": "gemini-3.1-pro-preview-customtools"
>         }
>       },
>       "gemini-3.1-flash-lite-preview": {
>         "extends": "chat-base-3",
>         "modelConfig": {
>           "model": "gemini-3.1-flash-lite-preview"
>         }
>       },
>       "gemini-2.5-pro": {
>         "extends": "chat-base-2.5",
>         "modelConfig": {
>           "model": "gemini-2.5-pro"
>         }
>       },
>       "gemini-2.5-flash": {
>         "extends": "chat-base-2.5",
>         "modelConfig": {
>           "model": "gemini-2.5-flash"
>         }
>       },
>       "gemini-2.5-flash-lite": {
>         "extends": "chat-base-2.5",
>         "modelConfig": {
>           "model": "gemini-2.5-flash-lite"
>         }
>       },
>       "gemini-3.1-flash-lite": {
>         "extends": "chat-base-3",
>         "modelConfig": {
>           "model": "gemini-3.1-flash-lite"
>         }
>       },
>       "gemini-3.5-flash": {
>         "extends": "chat-base-3",
>         "modelConfig": {
>           "model": "gemini-3.5-flash"
>         }
>       },
>       "gemma-4-31b-it": {
>         "extends": "chat-base-3",
>         "modelConfig": {
>           "model": "gemma-4-31b-it"
>         }
>       },
>       "gemma-4-26b-a4b-it": {
>         "extends": "chat-base-3",
>         "modelConfig": {
>           "model": "gemma-4-26b-a4b-it"
>         }
>       },
>       "gemini-2.5-flash-base": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "gemini-2.5-flash"
>         }
>       },
>       "gemini-3-flash-base": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "gemini-3-flash-preview"
>         }
>       },
>       "gemini-3.5-flash-base": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "gemini-3.5-flash"
>         }
>       },
>       "classifier": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "flash-lite",
>           "generateContentConfig": {
>             "maxOutputTokens": 1024,
>             "thinkingConfig": {
>               "thinkingBudget": 512
>             }
>           }
>         }
>       },
>       "prompt-completion": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "flash-lite",
>           "generateContentConfig": {
>             "temperature": 0.3,
>             "maxOutputTokens": 16000,
>             "thinkingConfig": {
>               "thinkingBudget": 0
>             }
>           }
>         }
>       },
>       "fast-ack-helper": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "flash-lite",
>           "generateContentConfig": {
>             "temperature": 0.2,
>             "maxOutputTokens": 120,
>             "thinkingConfig": {
>               "thinkingBudget": 0
>             }
>           }
>         }
>       },
>       "edit-corrector": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "flash-lite",
>           "generateContentConfig": {
>             "thinkingConfig": {
>               "thinkingBudget": 0
>             }
>           }
>         }
>       },
>       "summarizer-default": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "flash-lite",
>           "generateContentConfig": {
>             "maxOutputTokens": 2000
>           }
>         }
>       },
>       "summarizer-shell": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "flash-lite",
>           "generateContentConfig": {
>             "maxOutputTokens": 2000
>           }
>         }
>       },
>       "web-search": {
>         "extends": "gemini-3-flash-base",
>         "modelConfig": {
>           "generateContentConfig": {
>             "tools": [
>               {
>                 "googleSearch": {}
>               }
>             ]
>           }
>         }
>       },
>       "web-fetch": {
>         "extends": "gemini-3-flash-base",
>         "modelConfig": {
>           "generateContentConfig": {
>             "tools": [
>               {
>                 "urlContext": {}
>               }
>             ]
>           }
>         }
>       },
>       "web-fetch-fallback": {
>         "extends": "gemini-3-flash-base",
>         "modelConfig": {}
>       },
>       "loop-detection": {
>         "extends": "gemini-3-flash-base",
>         "modelConfig": {}
>       },
>       "loop-detection-double-check": {
>         "extends": "base",
>         "modelConfig": {
>           "model": "gemini-3-pro-preview"
>         }
>       },
>       "llm-edit-fixer": {
>         "extends": "gemini-3-flash-base",
>         "modelConfig": {}
>       },
>       "next-speaker-checker": {
>         "extends": "gemini-3-flash-base",
>         "modelConfig": {}
>       },
>       "context-snapshotter": {
>         "extends": "gemini-3-flash-base",
>         "modelConfig": {
>           "generateContentConfig": {
>             "thinkingConfig": {
>               "thinkingLevel": "HIGH"
>             },
>             "temperature": 1,
>             "topP": 0.95,
>             "topK": 64
>           }
>         }
>       },
>       "chat-compression-3-pro": {
>         "modelConfig": {
>           "model": "gemini-3-pro-preview"
>         }
>       },
>       "chat-compression-3-flash": {
>         "modelConfig": {
>           "model": "gemini-3-flash-preview"
>         }
>       },
>       "chat-compression-3.1-flash-lite": {
>         "modelConfig": {
>           "model": "gemini-3.1-flash-lite"
>         }
>       },
>       "chat-compression-2.5-pro": {
>         "modelConfig": {
>           "model": "gemini-2.5-pro"
>         }
>       },
>       "chat-compression-2.5-flash": {
>         "modelConfig": {
>           "model": "gemini-2.5-flash"
>         }
>       },
>       "chat-compression-2.5-flash-lite": {
>         "modelConfig": {
>           "model": "gemini-2.5-flash-lite"
>         }
>       },
>       "chat-compression-default": {
>         "modelConfig": {
>           "model": "gemini-3-pro-preview"
>         }
>       },
>       "agent-history-provider-summarizer": {
>         "modelConfig": {
>           "model": "gemini-3-flash-preview"
>         }
>       }
>     }
>     ```
>
> - **`modelConfigs.customAliases`** (object):
>
>   - **Description:** Custom named presets for model configs. These are merged
>     with (and override) the built-in aliases.
>   - **Default:** `{}`
>
> - **`modelConfigs.customOverrides`** (array):
>
>   - **Description:** Custom model config overrides. These are merged with (and
>     added to) the built-in overrides.
>   - **Default:** `[]`
>
> - **`modelConfigs.overrides`** (array):
>
>   - **Description:** Apply specific configuration overrides based on matches,
>     with a primary key of model (or alias). The most specific match will be
>     used.
>   - **Default:** `[]`
>
> - **`modelConfigs.modelDefinitions`** (object):
>
>   - **Description:** Registry of model metadata, including tier, family, and
>     features.
>   - **Default:**
>
>     ```json
>     {
>       "gemini-3.1-flash-lite": {
>         "tier": "flash-lite",
>         "family": "gemini-3",
>         "isPreview": false,
>         "isVisible": true,
>         "features": {
>           "thinking": false,
>           "multimodalToolUse": true
>         }
>       },
>       "gemini-3.1-pro-preview": {
>         "tier": "pro",
>         "family": "gemini-3",
>         "isPreview": true,
>         "isVisible": true,
>         "features": {
>           "thinking": true,
>           "multimodalToolUse": true
>         }
>       },
>       "gemini-3.1-pro-preview-customtools": {
>         "tier": "pro",
>         "family": "gemini-3",
>         "isPreview": true,
>         "isVisible": false,
>         "features": {
>           "thinking": true,
>           "multimodalToolUse": true
>         }
>       },
>       "gemini-3-pro-preview": {
>         "tier": "pro",
>         "family": "gemini-3",
>         "isPreview": true,
>         "isVisible": true,
>         "features": {
>           "thinking": true,
>           "multimodalToolUse": true
>         }
>       },
>       "gemini-3-flash-preview": {
>         "tier": "flash",
>         "family": "gemini-3",
>         "isPreview": true,
>         "isVisible": true,
>         "features": {
>           "thinking": false,
>           "multimodalToolUse": true
>         }
>       },
>       "gemini-3.5-flash": {
>         "tier": "flash",
>         "family": "gemini-3",
>         "isPreview": false,
>         "isVisible": true,
>         "features": {
>           "thinking": false,
>           "multimodalToolUse": true
>         }
>       },
>       "gemini-2.5-pro": {
>         "tier": "pro",
>         "family": "gemini-2.5",
>         "isPreview": false,
>         "isVisible": true,
>         "features": {
>           "thinking": false,
>           "multimodalToolUse": false
>         }
>       },
>       "gemini-2.5-flash": {
>         "tier": "flash",
>         "family": "gemini-2.5",
>         "isPreview": false,
>         "isVisible": true,
>         "features": {
>           "thinking": false,
>           "multimodalToolUse": false
>         }
>       },
>       "gemini-2.5-flash-lite": {
>         "tier": "flash-lite",
>         "family": "gemini-2.5",
>         "isPreview": false,
>         "isVisible": true,
>         "features": {
>           "thinking": false,
>           "multimodalToolUse": false
>         }
>       },
>       "gemma-4-31b-it": {
>         "displayName": "gemma-4-31b-it",
>         "tier": "custom",
>         "family": "gemma-4",
>         "isPreview": false,
>         "isVisible": true,
>         "features": {
>           "thinking": true,
>           "multimodalToolUse": false
>         }
>       },
>       "gemma-4-26b-a4b-it": {
>         "displayName": "gemma-4-26b-a4b-it",
>         "tier": "custom",
>         "family": "gemma-4",
>         "isPreview": false,
>         "isVisible": true,
>         "features": {
>           "thinking": true,
>           "multimodalToolUse": false
>         }
>       },
>       "auto": {
>         "displayName": "Auto",
>         "tier": "auto",
>         "isPreview": true,
>         "isVisible": true,
>         "features": {
>           "thinking": true,
>           "multimodalToolUse": false
>         }
>       },
>       "pro": {
>         "tier": "pro",
>         "isPreview": false,
>         "isVisible": false,
>         "features": {
>           "thinking": true,
>           "multimodalToolUse": false
>         }
>       },
>       "flash": {
>         "tier": "flash",
>         "isPreview": false,
>         "isVisible": false,
>         "features": {
>           "thinking": false,
>           "multimodalToolUse": false
>         }
>       },
>       "flash-lite": {
>         "tier": "flash-lite",
>         "isPreview": false,
>         "isVisible": false,
>         "features": {
>           "thinking": false,
>           "multimodalToolUse": false
>         }
>       },
>       "auto-gemini-3": {
>         "tier": "auto",
>         "family": "gemini-3",
>         "isPreview": true,
>         "isVisible": false
>       },
>       "auto-gemini-2.5": {
>         "tier": "auto",
>         "family": "gemini-2.5",
>         "isPreview": false,
>         "isVisible": false
>       }
>     }
>     ```
>
>   - **Requires restart:** Yes
>
> - **`modelConfigs.modelIdResolutions`** (object):
>
>   - **Description:** Rules for resolving requested model names to concrete model
>     IDs based on context.
>   - **Default:**
>
>     ```json
>     {
>       "gemma-4-31b-it": {
>         "default": "gemma-4-31b-it"
>       },
>       "gemma-4-26b-a4b-it": {
>         "default": "gemma-4-26b-a4b-it"
>       },
>       "gemini-3.1-pro-preview": {
>         "default": "gemini-3.1-pro-preview",
>         "contexts": [
>           {
>             "condition": {
>               "hasAccessToPreview": false
>             },
>             "target": "gemini-2.5-pro"
>           },
>           {
>             "condition": {
>               "useCustomTools": true
>             },
>             "target": "gemini-3.1-pro-preview-customtools"
>           }
>         ]
>       },
>       "gemini-3.1-pro-preview-customtools": {
>         "default": "gemini-3.1-pro-preview-customtools",
>         "contexts": [
>           {
>             "condition": {
>               "hasAccessToPreview": false
>             },
>             "target": "gemini-2.5-pro"
>           }
>         ]
>       },
>       "gemini-3-flash-preview": {
>         "default": "gemini-3-flash-preview",
>         "contexts": [
>           {
>             "condition": {
>               "hasAccessToPreview": false,
>               "useGemini3_5Flash": true
>             },
>             "target": "gemini-3.5-flash"
>           },
>           {
>             "condition": {
>               "hasAccessToPreview": false,
>               "useGemini3_5Flash": false
>             },
>             "target": "gemini-2.5-flash"
>           }
>         ]
>       },
>       "gemini-3.5-flash": {
>         "default": "gemini-3.5-flash",
>         "contexts": [
>           {
>             "condition": {
>               "useGemini3_5Flash": false,
>               "hasAccessToPreview": false
>             },
>             "target": "gemini-2.5-flash"
>           },
>           {
>             "condition": {
>               "useGemini3_5Flash": false
>             },
>             "target": "gemini-3-flash-preview"
>           }
>         ]
>       },
>       "gemini-2.5-flash": {
>         "default": "gemini-2.5-flash",
>         "contexts": [
>           {
>             "condition": {
>               "useGemini3_5Flash": true
>             },
>             "target": "gemini-3.5-flash"
>           }
>         ]
>       },
>       "gemini-3-pro-preview": {
>         "default": "gemini-3-pro-preview",
>         "contexts": [
>           {
>             "condition": {
>               "hasAccessToPreview": false
>             },
>             "target": "gemini-2.5-pro"
>           },
>           {
>             "condition": {
>               "useGemini3_1": true,
>               "useCustomTools": true
>             },
>             "target": "gemini-3.1-pro-preview-customto
