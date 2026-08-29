---
primary_sources:
  - id: T1-CONFIG
    title: "Gemini CLI configuration"
    url: "https://geminicli.com/docs/reference/configuration.md"
    section: "mcpServers; modelConfigs"
  - id: T1-GEN-SETTINGS
    title: "Advanced Model Configuration"
    url: "https://geminicli.com/docs/cli/generation-settings.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Configuration reference — MCP and model

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Gemini CLI configuration — Part 2

> ols"
>           },
>           {
>             "condition": {
>               "useGemini3_1": true
>             },
>             "target": "gemini-3.1-pro-preview"
>           }
>         ]
>       },
>       "auto": {
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
>             "target": "gemini-3.1-pro-preview-customtools"
>           },
>           {
>             "condition": {
>               "useGemini3_1": true
>             },
>             "target": "gemini-3.1-pro-preview"
>           }
>         ]
>       },
>       "pro": {
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
>             "target": "gemini-3.1-pro-preview-customtools"
>           },
>           {
>             "condition": {
>               "useGemini3_1": true
>             },
>             "target": "gemini-3.1-pro-preview"
>           }
>         ]
>       },
>       "gemini-3.1-flash-lite": {
>         "default": "gemini-3.1-flash-lite"
>       },
>       "flash": {
>         "default": "gemini-3-flash-preview",
>         "contexts": [
>           {
>             "condition": {
>               "useGemini3_5Flash": true
>             },
>             "target": "gemini-3.5-flash"
>           },
>           {
>             "condition": {
>               "hasAccessToPreview": false
>             },
>             "target": "gemini-2.5-flash"
>           }
>         ]
>       },
>       "flash-lite": {
>         "default": "gemini-3.1-flash-lite"
>       },
>       "auto-gemini-3": {
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
>             "target": "gemini-3.1-pro-preview-customtools"
>           },
>           {
>             "condition": {
>               "useGemini3_1": true
>             },
>             "target": "gemini-3.1-pro-preview"
>           }
>         ]
>       },
>       "auto-gemini-2.5": {
>         "default": "gemini-2.5-pro"
>       }
>     }
>     ```
>
>   - **Requires restart:** Yes
>
> - **`modelConfigs.classifierIdResolutions`** (object):
>
>   - **Description:** Rules for resolving classifier tiers (flash, pro) to
>     concrete model IDs.
>   - **Default:**
>
>     ```json
>     {
>       "flash": {
>         "default": "gemini-3-flash-preview",
>         "contexts": [
>           {
>             "condition": {
>               "useGemini3_5Flash": true
>             },
>             "target": "gemini-3.5-flash"
>           },
>           {
>             "condition": {
>               "hasAccessToPreview": false
>             },
>             "target": "gemini-2.5-flash"
>           },
>           {
>             "condition": {
>               "requestedModels": ["gemini-2.5-pro", "auto-gemini-2.5"]
>             },
>             "target": "gemini-2.5-flash"
>           }
>         ]
>       },
>       "pro": {
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
>               "requestedModels": ["gemini-2.5-pro", "auto-gemini-2.5"]
>             },
>             "target": "gemini-2.5-pro"
>           },
>           {
>             "condition": {
>               "useGemini3_1": true,
>               "useCustomTools": true
>             },
>             "target": "gemini-3.1-pro-preview-customtools"
>           },
>           {
>             "condition": {
>               "useGemini3_1": true
>             },
>             "target": "gemini-3.1-pro-preview"
>           }
>         ]
>       }
>     }
>     ```
>
>   - **Requires restart:** Yes
>
> - **`modelConfigs.modelChains`** (object):
>
>   - **Description:** Availability policy chains defining fallback behavior for
>     models.
>   - **Default:**
>
>     ```json
>     {
>       "preview": [
>         {
>           "model": "gemini-3-pro-preview",
>           "actions": {
>             "terminal": "prompt",
>             "transient": "prompt",
>             "not_found": "prompt",
>             "unknown": "prompt"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "terminal",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         },
>         {
>           "model": "gemini-3-flash-preview",
>           "isLastResort": true,
>           "maxAttempts": 10,
>           "actions": {
>             "terminal": "prompt",
>             "transient": "prompt",
>             "not_found": "prompt",
>             "unknown": "prompt"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "terminal",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         }
>       ],
>       "auto-preview": [
>         {
>           "model": "gemini-3-pro-preview",
>           "maxAttempts": 3,
>           "actions": {
>             "terminal": "prompt",
>             "transient": "silent",
>             "not_found": "prompt",
>             "unknown": "prompt"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "sticky_retry",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         },
>         {
>           "model": "gemini-3-flash-preview",
>           "isLastResort": true,
>           "maxAttempts": 10,
>           "actions": {
>             "terminal": "prompt",
>             "transient": "prompt",
>             "not_found": "prompt",
>             "unknown": "prompt"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "terminal",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         }
>       ],
>       "default": [
>         {
>           "model": "gemini-2.5-pro",
>           "actions": {
>             "terminal": "prompt",
>             "transient": "prompt",
>             "not_found": "prompt",
>             "unknown": "prompt"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "sticky_retry",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         },
>         {
>           "model": "gemini-2.5-flash",
>           "isLastResort": true,
>           "maxAttempts": 10,
>           "actions": {
>             "terminal": "prompt",
>             "transient": "prompt",
>             "not_found": "prompt",
>             "unknown": "prompt"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "terminal",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         }
>       ],
>       "auto-default": [
>         {
>           "model": "gemini-2.5-pro",
>           "maxAttempts": 3,
>           "actions": {
>             "terminal": "prompt",
>             "transient": "silent",
>             "not_found": "prompt",
>             "unknown": "prompt"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "sticky_retry",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         },
>         {
>           "model": "gemini-2.5-flash",
>           "isLastResort": true,
>           "maxAttempts": 10,
>           "actions": {
>             "terminal": "prompt",
>             "transient": "prompt",
>             "not_found": "prompt",
>             "unknown": "prompt"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "terminal",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         }
>       ],
>       "lite": [
>         {
>           "model": "flash-lite",
>           "actions": {
>             "terminal": "silent",
>             "transient": "silent",
>             "not_found": "silent",
>             "unknown": "silent"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "terminal",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         },
>         {
>           "model": "gemini-2.5-flash",
>           "actions": {
>             "terminal": "silent",
>             "transient": "silent",
>             "not_found": "silent",
>             "unknown": "silent"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "terminal",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         },
>         {
>           "model": "gemini-2.5-pro",
>           "isLastResort": true,
>           "actions": {
>             "terminal": "silent",
>             "transient": "silent",
>             "not_found": "silent",
>             "unknown": "silent"
>           },
>           "stateTransitions": {
>             "terminal": "terminal",
>             "transient": "terminal",
>             "not_found": "terminal",
>             "unknown": "terminal"
>           }
>         }
>       ]
>     }
>     ```
>
>   - **Requires restart:** Yes
>
> #### `agents`
>
> - **`agents.overrides`** (object):
>
>   - **Description:** Override settings for specific agents, e.g. to disable the
>     agent, set a custom model config, or run config.
>   - **Default:** `{}`
>   - **Requires restart:** Yes
>
> - **`agents.browser.sessionMode`** (enum):
>
>   - **Description:** Session mode: 'persistent', 'isolated', or 'existing'.
>   - **Default:** `"persistent"`
>   - **Values:** `"persistent"`, `"isolated"`, `"existing"`
>   - **Requires restart:** Yes
>
> - **`agents.browser.headless`** (boolean):
>
>   - **Description:** Run browser in headless mode.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`agents.browser.profilePath`** (string):
>
>   - **Description:** Path to browser profile directory for session persistence.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`agents.browser.visualModel`** (string):
>
>   - **Description:** Model for the visual agent's analyze_screenshot tool. When
>     set, enables the tool.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`agents.browser.allowedDomains`** (array):
>
>   - **Description:** A list of allowed domains for the browser agent (e.g.,
>     ["github.com", "*.google.com"]).
>   - **Default:**
>
>     ```json
>     ["github.com", "*.google.com", "localhost"]
>     ```
>
>   - **Requires restart:** Yes
>
> - **`agents.browser.disableUserInput`** (boolean):
>
>   - **Description:** Disable user input on browser window during automation.
>   - **Default:** `true`
>
> - **`agents.browser.maxActionsPerTask`** (number):
>
>   - **Description:** The maximum number of tool calls allowed per browser task.
>     Enforcement is hard: the agent will be terminated when the limit is reached.
>   - **Default:** `100`
>
> - **`agents.browser.confirmSensitiveActions`** (boolean):
>
>   - **Description:** Require manual confirmation for sensitive browser actions
>     (e.g., fill_form, evaluate_script).
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`agents.browser.blockFileUploads`** (boolean):
>   - **Description:** Hard-block file upload requests from the browser agent.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> #### `context`
>
> - **`context.fileName`** (string | string[]):
>
>   - **Description:** The name of the context file or files to load into memory.
>     Accepts either a single string or an array of strings.
>   - **Default:** `undefined`
>
> - **`context.importFormat`** (string):
>
>   - **Description:** The format to use when importing memory.
>   - **Default:** `undefined`
>
> - **`context.includeDirectoryTree`** (boolean):
>
>   - **Description:** Whether to include the directory tree of the current
>     working directory in the initial request to the model.
>   - **Default:** `true`
>
> - **`context.discoveryMaxDirs`** (number):
>
>   - **Description:** Maximum number of directories to search for memory.
>   - **Default:** `200`
>
> - **`context.memoryBoundaryMarkers`** (array):
>
>   - **Description:** File or directory names that mark the boundary for
>     GEMINI.md discovery. The upward traversal stops at the first directory
>     containing any of these markers. An empty array disables parent traversal.
>   - **Default:**
>
>     ```json
>     [".git"]
>     ```
>
>   - **Requires restart:** Yes
>
> - **`context.includeDirectories`** (array):
>
>   - **Description:** Additional directories to include in the workspace context.
>     Missing directories will be skipped with a warning.
>   - **Default:** `[]`
>
> - **`context.loadMemoryFromIncludeDirectories`** (boolean):
>
>   - **Description:** Controls how /memory reload loads GEMINI.md files. When
>     true, include directories are scanned; when false, only the current
>     directory is used.
>   - **Default:** `false`
>
> - **`context.fileFiltering.respectGitIgnore`** (boolean):
>
>   - **Description:** Respect .gitignore files when searching.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`context.fileFiltering.respectGeminiIgnore`** (boolean):
>
>   - **Description:** Respect .geminiignore files when searching.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`context.fileFiltering.enableFileWatcher`** (boolean):
>
>   - **Description:** Enable file watcher updates for @ file suggestions
>     (experimental).
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`context.fileFiltering.enableRecursiveFileSearch`** (boolean):
>
>   - **Description:** Enable recursive file search functionality when completing
>     @ references in the prompt.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`context.fileFiltering.enableFuzzySearch`** (boolean):
>
>   - **Description:** Enable fuzzy search when searching for files.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`context.fileFiltering.customIgnoreFilePaths`** (array):
>   - **Description:** Additional ignore file paths to respect. These files take
>     precedence over .geminiignore and .gitignore. Files earlier in the array
>     take precedence over files later in the array, e.g. the first file takes
>     precedence over the second one.
>   - **Default:** `[]`
>   - **Requires restart:** Yes
>
> #### `tools`
>
> - **`tools.sandbox`** (string):
>
>   - **Description:** Legacy full-process sandbox execution environment. Set to a
>     boolean to enable or disable the sandbox, provide a string path to a sandbox
>     profile, or specify an explicit sandbox command (e.g., "docker", "podman",
>     "lxc", "windows-native").
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`tools.sandboxAllowedPaths`** (array):
>
>   - **Description:** List of additional paths that the sandbox is allowed to
>     access.
>   - **Default:** `[]`
>   - **Requires restart:** Yes
>
> - **`tools.sandboxNetworkAccess`** (boolean):
>
>   - **Description:** Whether the sandbox is allowed to access the network.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`tools.shell.enableInteractiveShell`** (boolean):
>
>   - **Description:** Use node-pty for an interactive shell experience. Fallback
>     to child_process still applies.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`tools.shell.backgroundCompletionBehavior`** (enum):
>
>   - **Description:** Controls what happens when a background shell command
>     finishes. 'silent' (default): quietly exits in background. 'inject':
>     automatically returns output to agent. 'notify': shows brief message in
>     chat.
>   - **Default:** `"silent"`
>   - **Values:** `"silent"`, `"inject"`, `"notify"`
>
> - **`tools.shell.pager`** (string):
>
>   - **Description:** The pager command to use for shell output. Defaults to
>     `cat`.
>   - **Default:** `"cat"`
>
> - **`tools.shell.showColor`** (boolean):
>
>   - **Description:** Show color in shell output.
>   - **Default:** `true`
>
> - **`tools.shell.inactivityTimeout`** (number):
>
>   - **Description:** The maximum time in seconds allowed without output from the
>     shell command. Defaults to 5 minutes.
>   - **Default:** `300`
>
> - **`tools.shell.enableShellOutputEfficiency`** (boolean):
>
>   - **Description:** Enable shell output efficiency optimizations for better
>     performance.
>   - **Default:** `true`
>
> - **`tools.core`** (array):
>
>   - **Description:** Restrict the set of built-in tools with an allowlist. Match
>     semantics mirror tools.allowed; see the built-in tools documentation for
>     available names.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`tools.allowed`** (array):
>
>   - **Description:** Tool names that bypass the confirmation dialog. Useful for
>     trusted commands (for example ["run_shell_command(git)",
>     "run_shell_command(npm test)"]). See shell tool command restrictions for
>     matching details.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`tools.confirmationRequired`** (array):
>
>   - **Description:** Tool names that always require user confirmation. Takes
>     precedence over allowed tools and core tool allowlists.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`tools.exclude`** (array):
>
>   - **Description:** Tool names to exclude from discovery.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`tools.discoveryCommand`** (string):
>
>   - **Description:** Command to run for tool discovery.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`tools.callCommand`** (string):
>
>   - **Description:** Defines a custom shell command for invoking discovered
>     tools. The command must take the tool name as the first argument, read JSON
>     arguments from stdin, and emit JSON results on stdout.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`tools.useRipgrep`** (boolean):
>
>   - **Description:** Use ripgrep for file content search instead of the fallback
>     implementation. Provides faster search performance.
>   - **Default:** `true`
>
> - **`tools.truncateToolOutputThreshold`** (number):
>
>   - **Description:** Maximum characters to show when truncating large tool
>     outputs. Set to 0 or negative to disable truncation.
>   - **Default:** `40000`
>   - **Requires restart:** Yes
>
> - **`tools.disableLLMCorrection`** (boolean):
>   - **Description:** Disable LLM-based error correction for edit tools. When
>     enabled, tools will fail immediately if exact string matches are not found,
>     instead of attempting to self-correct.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> #### `mcp`
>
> - **`mcp.serverCommand`** (string):
>
>   - **Description:** Command to start an MCP server.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`mcp.allowed`** (array):
>
>   - **Description:** A list of MCP servers to allow.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`mcp.excluded`** (array):
>   - **Description:** A list of MCP servers to exclude.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> #### `useWriteTodos`
>
> - **`useWriteTodos`** (boolean):
>   - **Description:** Enable the write_todos tool.
>   - **Default:** `true`
>
> #### `security`
>
> - **`security.toolSandboxing`** (boolean):
>
>   - **Description:** Tool-level sandboxing. Isolates individual tools instead of
>     the entire CLI process.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`security.disableYoloMode`** (boolean):
>
>   - **Description:** Disable YOLO mode, even if enabled by a flag.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`security.disableAlwaysAllow`** (boolean):
>
>   - **Description:** Disable "Always allow" options in tool confirmation
>     dialogs.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`security.enablePermanentToolApproval`** (boolean):
>
>   - **Description:** Enable the "Allow for all future sessions" option in tool
>     confirmation dialogs.
>   - **Default:** `false`
>
> - **`security.autoAddToPolicyByDefault`** (boolean):
>
>   - **Description:** When enabled, the "Allow for all future sessions" option
>     becomes the default choice for low-risk tools in trusted workspaces.
>   - **Default:** `false`
>
> - **`security.blockGitExtensions`** (boolean):
>
>   - **Description:** Blocks installing and loading extensions from Git.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`security.allowedExtensions`** (array):
>
>   - **Description:** List of Regex patterns for allowed extensions. If nonempty,
>     only extensions that match the patterns in this list are allowed. Overrides
>     the blockGitExtensions setting.
>   - **Default:** `[]`
>   - **Requires restart:** Yes
>
> - **`security.folderTrust.enabled`** (boolean):
>
>   - **Description:** Setting to track whether Folder trust is enabled.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`security.environmentVariableRedaction.allowed`** (array):
>
>   - **Description:** Environment variables to always allow (bypass redaction).
>   - **Default:** `[]`
>   - **Requires restart:** Yes
>
> - **`security.environmentVariableRedaction.blocked`** (array):
>
>   - **Description:** Environment variables to always redact.
>   - **Default:** `[]`
>   - **Requires restart:** Yes
>
> - **`security.environmentVariableRedaction.enabled`** (boolean):
>
>   - **Description:** Enable redaction of environment variables that may contain
>     secrets.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`security.auth.selectedType`** (string):
>
>   - **Description:** The currently selected authentication type.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`security.auth.enforcedType`** (string):
>
>   - **Description:** The required auth type. If this does not match the selected
>     auth type, the user will be prompted to re-authenticate.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`security.auth.useExternal`** (boolean):
>
>   - **Description:** Whether to use an external authentication flow.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`security.enableConseca`** (boolean):
>   - **Description:** Enable the context-aware security checker. This feature
>     uses an LLM to dynamically generate and enforce security policies for tool
>     use based on your prompt, providing an additional layer of protection
>     against unintended actions.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> #### `advanced`
>
> - **`advanced.autoConfigureMemory`** (boolean):
>
>   - **Description:** Automatically configure Node.js memory limits. Note:
>     Because memory is allocated during the initial process boot, this setting is
>     only read from the global user settings file and ignores workspace-level
>     overrides.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`advanced.dnsResolutionOrder`** (string):
>
>   - **Description:** The DNS resolution order.
>   - **Default:** `undefined`
>   - **Requires restart:** Yes
>
> - **`advanced.excludedEnvVars`** (array):
>
>   - **Description:** Environment variables to exclude from project context.
>   - **Default:**
>
>     ```json
>     ["DEBUG", "DEBUG_MODE"]
>     ```
>
> - **`advanced.ignoreLocalEnv`** (boolean):
>
>   - **Description:** Whether to ignore generic .env files in the project
>     directory.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`advanced.bugCommand`** (object):
>   - **Description:** Configuration for the bug report command.
>   - **Default:** `undefined`
>
> #### `experimental`
>
> - **`experimental.gemma`** (boolean):
>
>   - **Description:** Enable access to Gemma 4 models via Gemini API.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`experimental.voiceMode`** (boolean):
>
>   - **Description:** Enable experimental voice dictation and commands (/voice,
>     /voice model).
>   - **Default:** `false`
>
> - **`experimental.voice.activationMode`** (enum):
>
>   - **Description:** How to trigger voice recording with the Space key.
>   - **Default:** `"push-to-talk"`
>   - **Values:** `"push-to-talk"`, `"toggle"`
>
> - **`experimental.voice.backend`** (enum):
>
>   - **Description:** The backend to use for voice transcription. Note: When
>     using the Gemini Live backend, voice recordings are sent to Google Cloud for
>     transcription.
>   - **Default:** `"gemini-live"`
>   - **Values:** `"gemini-live"`, `"whisper"`
>
> - **`experimental.voice.whisperModel`** (enum):
>
>   - **Description:** The Whisper model to use for local transcription.
>   - **Default:** `"ggml-base.en.bin"`
>   - **Values:** `"ggml-tiny.en.bin"`, `"ggml-base.en.bin"`,
>     `"ggml-large-v3-turbo-q5_0.bin"`, `"ggml-large-v3-turbo-q8_0.bin"`
>
> - **`experimental.voice.stopGracePeriodMs`** (number):
>
>   - **Description:** How long to wait for final transcription after stopping
>     recording.
>   - **Default:** `4000`
>
> - **`experimental.adk.agentSessionNoninteractiveEnabled`** (boolean):
>
>   - **Description:** Enable non-interactive agent sessions.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.adk.agentSessionInteractiveEnabled`** (boolean):
>
>   - **Description:** Enable the agent session implementation for the interactive
>     CLI.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.adk.agentSessionSubagentEnabled`** (boolean):
>
>   - **Description:** Route subagent invocations through the AgentSession
>     protocol instead of legacy executors.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.enableAgents`** (boolean):
>
>   - **Description:** Enable local and remote subagents.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`experimental.worktrees`** (boolean):
>
>   - **Description:** Enable automated Git worktree management for parallel work.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.extensionManagement`** (boolean):
>
>   - **Description:** Enable extension management features.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`experimental.extensionConfig`** (boolean):
>
>   - **Description:** Enable requesting and fetching of extension settings.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`experimental.extensionRegistry`** (boolean):
>
>   - **Description:** Enable extension registry explore UI.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.extensionRegistryURI`** (string):
>
>   - **Description:** The URI (web URL or local file path) of the extension
>     registry.
>   - **Default:** `"https://geminicli.com/extensions.json"`
>   - **Requires restart:** Yes
>
> - **`experimental.extensionReloading`** (boolean):
>
>   - **Description:** Enables extension loading/unloading within the CLI session.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.useOSC52Paste`** (boolean):
>
>   - **Description:** Use OSC 52 for pasting. This may be more robust than the
>     default system when using remote terminal sessions (if your terminal is
>     configured to allow it).
>   - **Default:** `false`
>
> - **`experimental.useOSC52Copy`** (boolean):
>
>   - **Description:** Use OSC 52 for copying. This may be more robust than the
>     default system when using remote terminal sessions (if your terminal is
>     configured to allow it).
>   - **Default:** `false`
>
> - **`experimental.taskTracker`** (boolean):
>
>   - **Description:** Enable task tracker tools.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.modelSteering`** (boolean):
>
>   - **Description:** Enable model steering (user hints) to guide the model
>     during tool execution.
>   - **Default:** `false`
>
> - **`experimental.directWebFetch`** (boolean):
>
>   - **Description:** Enable web fetch behavior that bypasses LLM summarization.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.dynamicModelConfiguration`** (boolean):
>
>   - **Description:** Enable dynamic model configuration (definitions,
>     resolutions, and chains) via settings.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.gemmaModelRouter.enabled`** (boolean):
>
>   - **Description:** Enable the Gemma Model Router (experimental). Requires a
>     local endpoint serving Gemma via the Gemini API using LiteRT-LM shim.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.gemmaModelRouter.autoStartServer`** (boolean):
>
>   - **Description:** Automatically start the LiteRT-LM server when Gemini CLI
>     starts and the Gemma router is enabled.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.gemmaModelRouter.binaryPath`** (string):
>
>   - **Description:** Custom path to the LiteRT-LM binary. Leave empty to use the
>     default location (~/.gemini/bin/litert/).
>   - **Default:** `""`
>   - **Requires restart:** Yes
>
> - **`experimental.gemmaModelRouter.classifier.host`** (string):
>
>   - **Description:** The host of the classifier.
>   - **Default:** `"http://localhost:9379"`
>   - **Requires restart:** Yes
>
> - **`experimental.gemmaModelRouter.classifier.model`** (string):
>
>   - **Description:** The model to use for the classifier. Only tested on
>     `gemma3-1b-gpu-custom`.
>   - **Default:** `"gemma3-1b-gpu-custom"`
>   - **Requires restart:** Yes
>
> - **`experimental.stressTestProfile`** (boolean):
>
>   - **Description:** Significantly lowers token limits to force early garbage
>     collection and distillation for testing purposes.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.autoMemory`** (boolean):
>
>   - **Description:** Automatically extract memory patches and skills from past
>     sessions in the background. Every change is written as a unified diff
>     `.patch` file under `<projectMemoryDir>/.inbox/<kind>/` and held for review
>     in /memory inbox; nothing is applied until you approve it.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.generalistProfile`** (boolean):
>
>   - **Description:** Suitable for general coding and software development tasks.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.powerUserProfile`** (boolean):
>
>   - **Description:** Less cache friendly version of the generalist profile.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.contextManagement`** (boolean):
>
>   - **Description:** Enable logic for context management.
>   - **Default:** `false`
>   - **Requires restart:** Yes
>
> - **`experimental.topicUpdateNarration`** (boolean):
>   - **Description:** Deprecated: Use general.topicUpdateNarration instead.
>   - **Default:** `false`
>
> #### `skills`
>
> - **`skills.enabled`** (boolean):
>
>   - **Description:** Enable Agent Skills.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`skills.disabled`** (array):
>   - **Description:** List of disabled skills.
>   - **Default:** `[]`
>   - **Requires restart:** Yes
>
> #### `hooksConfig`
>
> - **`hooksConfig.enabled`** (boolean):
>
>   - **Description:** Canonical toggle for the hooks system. When disabled, no
>     hooks will be executed.
>   - **Default:** `true`
>   - **Requires restart:** Yes
>
> - **`hooksConfig.disabled`** (array):
>
>   - **Description:** List of hook names (commands) that should be disabled.
>     Hooks in this list will not execute even if configured.
>   - **Default:** `[]`
>
> - **`hooksConfig.notifications`** (boolean):
>   - **Description:** Show visual indicators when hooks are executing.
>   - **Default:** `true`
>
> #### `hooks`
>
> - **`hooks.BeforeTool`** (array):
>
>   - **Description:** Hooks that execute before tool execution. Can intercept,
>     validate, or modify tool calls.
>   - **Default:** `[]`
>
> - **`hooks.AfterTool`** (array):
>
>   - **Description:** Hooks that execute after tool execution. Can process
>     results, log outputs, or trigger follow-up actions.
>   - **Default:** `[]`
>
> - **`hooks.BeforeAgent`** (array):
>
>   - **Description:** Hooks that execute before agent loop starts. Can set up
>     context or initialize resources.
>   - **Default:** `[]`
>
> - **`hooks.AfterAgent`** (array):
>
>   - **Description:** Hooks that execute after agent loop completes. Can perform
>     cleanup or summarize results.
>   - **Default:** `[]`
>
> - **`hooks.Notification`** (array):
>
>   - **Description:** Hooks that execute on notification events (errors,
>     warnings, info). Can log or alert on specific conditions.
>   - **Default:** `[]`
>
> - **`hooks.SessionStart`** (array):
>
>   - **Description:** Hooks that execute when a session starts. Can initialize
>     session-specific resources or state.
>   - **Default:** `[]`
>
> - **`hooks.SessionEnd`** (array):
>
>   - **Description:** Hooks that execute when a session ends. Can perform cleanup
>     or persist session data.
>   - **Default:** `[]`
>
> - **`hooks.PreCompress`** (array):
>
>   - **Description:** Hooks that execute before chat history compression. Can
>     back up or analyze conversation before compression.
>   - **Default:** `[]`
>
> - **`hooks.BeforeModel`** (array):
>
>   - **Description:** Hooks that execute before LLM requests. Can modify prompts,
>     inject context, or control model parameters.
>   - **Default:** `[]`
>
> - **`hooks.AfterModel`** (array):
>
>   - **Description:** Hooks that execute after LLM responses. Can process
>     outputs, extract information, or log interactions.
>   - **Default:** `[]`
>
> - **`hooks.BeforeToolSelection`** (array):
>   - **Description:** Hooks that execute before tool selection. Can filter or
>     prioritize available tools dynamically.
>   - **Default:** `[]`
>
> #### `contextManagement`
>
> - **`contextManagement.historyWindow.m

### Source: Advanced Model Configuration — Full page

> This guide details the Model Configuration system within Gemini CLI. Designed
> for researchers, AI quality engineers, and advanced users, this system provides
> a rigorous framework for managing generative model hyperparameters and
> behaviors.
>
> <!-- prettier-ignore -->
> > [!WARNING]
> > This is a power-user feature. Configuration values are passed
> > directly to the model provider with minimal validation. Incorrect settings
> > (for example, incompatible parameter combinations) may result in runtime
> > errors from the API.
>
> ## 1. System Overview
>
> The Model Configuration system (`ModelConfigService`) enables deterministic
> control over model generation. It decouples the requested model identifier (for
> example, a CLI flag or agent request) from the underlying API configuration.
> This allows for:
>
> - **Precise Hyperparameter Tuning**: Direct control over `temperature`, `topP`,
>   `thinkingBudget`, and other SDK-level parameters.
> - **Environment-Specific Behavior**: Distinct configurations for different
>   operating contexts (for example, testing vs. production).
> - **Agent-Scoped Customization**: Applying specific settings only when a
>   particular agent is active.
>
> The system operates on two core primitives: **Aliases** and **Overrides**.
>
> ## 2. Configuration Primitives
>
> These settings are located under the `modelConfigs` key in your configuration
> file.
>
> ### Aliases (`customAliases`)
>
> Aliases are named, reusable configuration presets. Users should define their own
> aliases (or override system defaults) in the `customAliases` map.
>
> - **Inheritance**: An alias can `extends` another alias (including system
>   defaults like `chat-base`), inheriting its `modelConfig`. Child aliases can
>   overwrite or augment inherited settings.
> - **Abstract Aliases**: An alias is not required to specify a concrete `model`
>   if it serves purely as a base for other aliases.
>
> **Example Hierarchy**:
>
> ```json
> "modelConfigs": {
>   "customAliases": {
>     "base": {
>       "modelConfig": {
>         "generateContentConfig": { "temperature": 0.0 }
>       }
>     },
>     "chat-base": {
>       "extends": "base",
>       "modelConfig": {
>         "generateContentConfig": { "temperature": 0.7 }
>       }
>     }
>   }
> }
> ```
>
> ### Overrides (`overrides`)
>
> Overrides are conditional rules that inject configuration based on the runtime
> context. They are evaluated dynamically for each model request.
>
> - **Match Criteria**: Overrides apply when the request context matches the
>   specified `match` properties.
>   - `model`: Matches the requested model name or alias.
>   - `overrideScope`: Matches the distinct scope of the request (typically the
>     agent name, for example, `codebaseInvestigator`).
>
> **Example Override**:
>
> ```json
> "modelConfigs": {
>   "overrides": [
>     {
>       "match": {
>         "overrideScope": "codebaseInvestigator"
>       },
>       "modelConfig": {
>         "generateContentConfig": { "temperature": 0.1 }
>       }
>     }
>   ]
> }
> ```
>
> ## 3. Resolution Strategy
>
> The `ModelConfigService` resolves the final configuration through a two-step
> process:
>
> ### Step 1: Alias Resolution
>
> The requested model string is looked up in the merged map of system `aliases`
> and user `customAliases`.
>
> 1.  If found, the system recursively resolves the `extends` chain.
> 2.  Settings are merged from parent to child (child wins).
> 3.  This results in a base `ResolvedModelConfig`.
> 4.  If not found, the requested string is treated as the raw model name.
>
> ### Step 2: Override Application
>
> The system evaluates the `overrides` list against the request context (`model`
> and `overrideScope`).
>
> 1.  **Filtering**: All matching overrides are identified.
> 2.  **Sorting**: Matches are prioritized by **specificity** (the number of
>     matched keys in the `match` object).
>     - Specific matches (for example, `model` + `overrideScope`) override broad
>       matches (for example, `model` only).
>     - Tie-breaking: If specificity is equal, the order of definition in the
>       `overrides` array is preserved (last one wins).
> 3.  **Merging**: The configurations from the sorted overrides are merged
>     sequentially onto the base configuration.
>
> ## 4. Configuration Reference
>
> The configuration follows the `ModelConfigServiceConfig` interface.
>
> ### `ModelConfig` Object
>
> Defines the actual parameters for the model.
>
> | Property                | Type     | Description                                                               |
> | :---------------------- | :------- | :------------------------------------------------------------------------ |
> | `model`                 | `string` | The identifier of the model to be called (for example, `gemini-2.5-pro`). |
> | `generateContentConfig` | `object` | The configuration object passed to the `@google/genai` SDK.               |
>
> ### `GenerateContentConfig` (Common Parameters)
>
> Directly maps to the SDK's `GenerateContentConfig`. Common parameters include:
>
> - **`temperature`**: (`number`) Controls output randomness. Lower values (0.0)
>   are deterministic; higher values (>0.7) are creative.
> - **`topP`**: (`number`) Nucleus sampling probability.
> - **`maxOutputTokens`**: (`number`) Limit on generated response length.
> - **`thinkingConfig`**: (`object`) Configuration for models with reasoning
>   capabilities (for example, `thinkingBudget`, `includeThoughts`).
>
> ## 5. Practical Examples
>
> ### Defining a Deterministic Baseline
>
> Create an alias for tasks requiring high precision, extending the standard chat
> configuration but enforcing zero temperature.
>
> ```json
> "modelConfigs": {
>   "customAliases": {
>     "precise-mode": {
>       "extends": "chat-base",
>       "modelConfig": {
>         "generateContentConfig": {
>           "temperature": 0.0,
>           "topP": 1.0
>         }
>       }
>     }
>   }
> }
> ```
>
> ### Agent-Specific Parameter Injection
>
> Enforce extended thinking budgets for a specific agent without altering the
> global default, for example for the `codebaseInvestigator`.
>
> ```json
> "modelConfigs": {
>   "overrides": [
>     {
>       "match": {
>         "overrideScope": "codebaseInvestigator"
>       },
>       "modelConfig": {
>         "generateContentConfig": {
>           "thinkingConfig": { "thinkingBudget": 4096 }
>         }
>       }
>     }
>   ]
> }
> ```
>
> ### Experimental Model Evaluation
>
> Route traffic for a specific alias to a preview model for A/B testing, without
> changing client code.
>
> ```json
> "modelConfigs": {
>   "overrides": [
>     {
>       "match": {
>         "model": "gemini-2.5-pro"
>       },
>       "modelConfig": {
>         "model": "gemini-2.5-pro-experimental-001"
>       }
>     }
>   ]
> }
> ```
