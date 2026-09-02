---
primary_sources:
  - id: T1-CONFIG
    title: "Config"
    url: "https://opencode.ai/docs/config.md"
    section: "Locations"
  - id: T1-RULES
    title: "Rules"
    url: "https://opencode.ai/docs/rules.md"
    section: "Types"
  - id: T1-PLUGINS
    title: "Plugins"
    url: "https://opencode.ai/docs/plugins.md"
    section: "From local files"
  - id: T1-SKILLS
    title: "Skills"
    url: "https://opencode.ai/docs/skills.md"
    section: "Place files"
  - id: T1-TROUBLESHOOTING
    title: "Troubleshooting"
    url: "https://opencode.ai/docs/troubleshooting.md"
    section: "Storage"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Project vs global path matrix

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: Synthesized matrix (index table)

> ## Summary table
>
> | Concern | Project path | Global (`~/`) path | Notes from docs |
> | ------- | ------------ | ------------------- | ---------------- |
> | Main config | `opencode.json` in project root | `~/.config/opencode/opencode.json` | Config files merge; later overrides earlier for conflicting keys |
> | TUI config | `tui.json` alongside project config | `~/.config/opencode/tui.json` | Separate schema at opencode.ai/tui.json |
> | Rules / AGENTS.md | `AGENTS.md` in project (walk up to git root) | `~/.config/opencode/AGENTS.md` | Claude fallback: `CLAUDE.md`, `~/.claude/CLAUDE.md` |
> | Custom instructions | `instructions` in project `opencode.json` | same in global config | Remote URLs supported |
> | Agents | `.opencode/agents/` | `~/.config/opencode/agents/` | Singular `agent/` also supported |
> | Commands | `.opencode/commands/` | `~/.config/opencode/commands/` | |
> | Modes | `.opencode/modes/` | via config dirs | Listed in config custom directory |
> | Plugins (official) | `.opencode/plugins/` | `~/.config/opencode/plugins/` | agentd uses `.opencode/plugin/` (non-official) |
> | Skills | `.opencode/skills/*/SKILL.md` | `~/.config/opencode/skills/` | Also `.claude/skills`, `.agents/skills` |
> | Custom tools | `.opencode/tools/` | `~/.config/opencode/tools/` | |
> | Themes | `.opencode/themes/` | `~/.config/opencode/themes/` | |
> | Plugin npm deps | `.opencode/package.json` | config dir `package.json` | Bun install at startup |
> | Auth credentials | — | `~/.local/share/opencode/auth.json` | Via `opencode auth login` |
> | MCP OAuth tokens | — | `~/.local/share/opencode/mcp-auth.json` | |
> | Session storage | `./<hash>/storage/` in git repo | `./global/storage/` if not git | See troubleshooting |
> | Cache | — | `~/.cache/opencode/` | Provider packages, node_modules |
> | Managed config | — | `/Library/Application Support/opencode/` (macOS), `/etc/opencode/` (Linux) | Highest priority |
> | Remote org defaults | — | `.well-known/opencode` | Loaded first as base layer |

### Source: OpenCode Config — Locations

> ## Locations
>
> You can place your config in a couple of different locations and they have a
> different order of precedence.
>
> :::note
> Configuration files are **merged together**, not replaced.
> :::
>
> Configuration files are merged together, not replaced. Settings from the following config locations are combined. Later configs override earlier ones only for conflicting keys. Non-conflicting settings from all configs are preserved.
>
> For example, if your global config sets `autoupdate: true` and your project config sets `model: "anthropic/claude-sonnet-4-5"`, the final configuration will include both settings.
>
> ---
>
> ### Precedence order
>
> Config sources are loaded in this order (later sources override earlier ones):
>
> 1. **Remote config** (from `.well-known/opencode`) - organizational defaults
> 2. **Global config** (`~/.config/opencode/opencode.json`) - user preferences
> 3. **Custom config** (`OPENCODE_CONFIG` env var) - custom overrides
> 4. **Project config** (`opencode.json` in project) - project-specific settings
> 5. **`.opencode` directories** - agents, commands, plugins
> 6. **Inline config** (`OPENCODE_CONFIG_CONTENT` env var) - runtime overrides
> 7. **Managed config files** (`/Library/Application Support/opencode/` on macOS) - admin-controlled
> 8. **macOS managed preferences** (`.mobileconfig` via MDM) - highest priority, not user-overridable
>
> This means project configs can override global defaults, and global configs can override remote organizational defaults. Managed settings override everything.
>
> :::note
> The `.opencode` and `~/.config/opencode` directories use **plural names** for subdirectories: `agents/`, `commands/`, `modes/`, `plugins/`, `skills/`, `tools/`, and `themes/`. Singular names (e.g., `agent/`) are also supported for backwards compatibility.
> :::
>
> ---
>
> ### Remote
>
> Organizations can provide default configuration via the `.well-known/opencode` endpoint. This is fetched automatically when you authenticate with a provider that supports it.
>
> Remote config is loaded first, serving as the base layer. All other config sources (global, project) can override these defaults.
>
> For example, if your organization provides MCP servers that are disabled by default:
>
> ```json title="Remote config from .well-known/opencode"
> {
>   "mcp": {
>     "jira": {
>       "type": "remote",
>       "url": "https://jira.example.com/mcp",
>       "enabled": false
>     }
>   }
> }
> ```
>
> You can enable specific servers in your local config:
>
> ```json title="opencode.json"
> {
>   "mcp": {
>     "jira": {
>       "type": "remote",
>       "url": "https://jira.example.com/mcp",
>       "enabled": true
>     }
>   }
> }
> ```
>
> ---
>
> ### Global
>
> Place your global OpenCode config in `~/.config/opencode/opencode.json`. Use global config for user-wide server/runtime preferences like providers, models, and permissions.
>
> For TUI-specific settings, use `~/.config/opencode/tui.json`.
>
> Global config overrides remote organizational defaults.
>
> ---
>
> ### Per project
>
> Add `opencode.json` in your project root. Project config has the highest precedence among standard config files - it overrides both global and remote configs.
>
> For project-specific TUI settings, add `tui.json` alongside it.
>
> :::tip
> Place project specific config in the root of your project.
> :::
>
> When OpenCode starts up, it first looks for a config file in the current directory, then traverses up to the nearest Git directory.
>
> This is also safe to be checked into Git and uses the same schema as the global one.
>
> ---
>
> ### Custom path
>
> Specify a custom config file path using the `OPENCODE_CONFIG` environment variable.
>
> ```bash
> export OPENCODE_CONFIG=/path/to/my/custom-config.json
> opencode run "Hello world"
> ```
>
> Custom config is loaded between global and project configs in the precedence order.
>
> ---
>
> ### Custom directory
>
> Specify a custom config directory using the `OPENCODE_CONFIG_DIR`
> environment variable. This directory will be searched for agents, commands,
> modes, and plugins just like the standard `.opencode` directory, and should
> follow the same structure.
>
> ```bash
> export OPENCODE_CONFIG_DIR=/path/to/my/config-directory
> opencode run "Hello world"
> ```
>
> The custom directory is loaded after the global config and `.opencode` directories, so it **can override** their settings.
>
> ---
>
> ### Managed settings
>
> Organizations can enforce configuration that users cannot override. Managed settings are loaded at the highest priority tier.
>
> #### File-based
>
> Drop an `opencode.json` or `opencode.jsonc` file in the system managed config directory:
>
> | Platform | Path                                     |
> | -------- | ---------------------------------------- |
> | macOS    | `/Library/Application Support/opencode/` |
> | Linux    | `/etc/opencode/`                         |
> | Windows  | `%ProgramData%\opencode`                 |
>
> These directories require admin/root access to write, so users cannot modify them.
>
> #### macOS managed preferences
>
> On macOS, OpenCode reads managed preferences from the `ai.opencode.managed` preference domain. Deploy a `.mobileconfig` via MDM (Jamf, Kandji, FleetDM) and the settings are enforced automatically.
>
> OpenCode checks these paths:
>
> 1. `/Library/Managed Preferences/<user>/ai.opencode.managed.plist`
> 2. `/Library/Managed Preferences/ai.opencode.managed.plist`
>
> The plist keys map directly to `opencode.json` fields. MDM metadata keys (`PayloadUUID`, `PayloadType`, etc.) are stripped automatically.
>
> **Creating a `.mobileconfig`**
>
> Use the `ai.opencode.managed` PayloadType. The OpenCode config keys go directly in the payload dict:
>
> ```xml
> <?xml version="1.0" encoding="UTF-8"?>
> <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
>   "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
> <plist version="1.0">
> <dict>
>   <key>PayloadContent</key>
>   <array>
>     <dict>
>       <key>PayloadType</key>
>       <string>ai.opencode.managed</string>
>       <key>PayloadIdentifier</key>
>       <string>com.example.opencode.config</string>
>       <key>PayloadUUID</key>
>       <string>GENERATE-YOUR-OWN-UUID</string>
>       <key>PayloadVersion</key>
>       <integer>1</integer>
>       <key>share</key>
>       <string>disabled</string>
>       <key>server</key>
>       <dict>
>         <key>hostname</key>
>         <string>127.0.0.1</string>
>       </dict>
>       <key>permission</key>
>       <dict>
>         <key>*</key>
>         <string>ask</string>
>         <key>bash</key>
>         <dict>
>           <key>*</key>
>           <string>ask</string>
>           <key>rm -rf *</key>
>           <string>deny</string>
>         </dict>
>       </dict>
>     </dict>
>   </array>
>   <key>PayloadType</key>
>   <string>Configuration</string>
>   <key>PayloadIdentifier</key>
>   <string>com.example.opencode</string>
>   <key>PayloadUUID</key>
>   <string>GENERATE-YOUR-OWN-UUID</string>
>   <key>PayloadVersion</key>
>   <integer>1</integer>
> </dict>
> </plist>
> ```
>
> Generate unique UUIDs with `uuidgen`. Customize the settings to match your organization's requirements.
>
> **Deploying via MDM**
>
> - **Jamf Pro:** Computers > Configuration Profiles > Upload > scope to target devices or smart groups
> - **FleetDM:** Add the `.mobileconfig` to your gitops repo under `mdm.macos_settings.custom_settings` and run `fleetctl apply`
>
> **Verifying on a device**
>
> Double-click the `.mobileconfig` to install locally for testing (shows in System Settings > Privacy & Security > Profiles), then run:
>
> ```bash
> opencode debug config
> ```
>
> All managed preference keys appear in the resolved config and cannot be overridden by user or project configuration.
>
> ---

### Source: OpenCode Rules — Project and Global

> You can provide custom instructions to opencode by creating an `AGENTS.md` file. This is similar to Cursor's rules. It contains instructions that will be included in the LLM's context to customize its behavior for your specific project.
>
> ---
>
>
> To create a new `AGENTS.md` file, you can run the `/init` command in opencode.
>
> :::tip
> You should commit your project's `AGENTS.md` file to Git.
> :::
>
> `/init` scans the important files in your repo, may ask a couple of targeted questions when the codebase cannot answer them, and then creates or updates `AGENTS.md` with concise project-specific guidance.
>
> It focuses on the things future agent sessions are most likely to need:
>
> - build, lint, and test commands
> - command order and focused verification steps when they matter
> - architecture and repo structure that are not obvious from filenames alone
> - project-specific conventions, setup quirks, and operational gotchas
> - references to existing instruction sources like Cursor or Copilot rules
>
> If you already have an `AGENTS.md`, `/init` will improve it in place instead of blindly replacing it.
>
> ---
>
>
> You should commit your project's `AGENTS.md` file to Git.
> :::
>
> `/init` scans the important files in your repo, may ask a couple of targeted questions when the codebase cannot answer them, and then creates or updates `AGENTS.md` with concise project-specific guidance.
>
> It focuses on the things future agent sessions are most likely to need:
>
> - build, lint, and test commands
> - command order and focused verification steps when they matter
> - architecture and repo structure that are not obvious from filenames alone
> - project-specific conventions, setup quirks, and operational gotchas
> - references to existing instruction sources like Cursor or Copilot rules
>
> If you already have an `AGENTS.md`, `/init` will improve it in place instead of blindly replacing it.
>
> ---
>
>
> `/init` scans the important files in your repo, may ask a couple of targeted questions when the codebase cannot answer them, and then creates or updates `AGENTS.md` with concise project-specific guidance.
>
> It focuses on the things future agent sessions are most likely to need:
>
> - build, lint, and test commands
> - command order and focused verification steps when they matter
> - architecture and repo structure that are not obvious from filenames alone
> - project-specific conventions, setup quirks, and operational gotchas
> - references to existing instruction sources like Cursor or Copilot rules
>
> If you already have an `AGENTS.md`, `/init` will improve it in place instead of blindly replacing it.
>
> ---
>
>
> If you already have an `AGENTS.md`, `/init` will improve it in place instead of blindly replacing it.
>
> ---
>
>
> You can also just create this file manually. Here's an example of some things you can put into an `AGENTS.md` file.
>
> ```markdown title="AGENTS.md"
> # SST v3 Monorepo Project
>
> This is an SST v3 monorepo with TypeScript. The project uses bun workspaces for package management.
>
>
> ```markdown title="AGENTS.md"
> # SST v3 Monorepo Project
>
> This is an SST v3 monorepo with TypeScript. The project uses bun workspaces for package management.
>
>
> opencode also supports reading the `AGENTS.md` file from multiple locations. And this serves different purposes.
>
> ### Project
>
> Place an `AGENTS.md` in your project root for project-specific rules. These only apply when you are working in this directory or its sub-directories.
>
> ### Global
>
> You can also have global rules in a `~/.config/opencode/AGENTS.md` file. This gets applied across all opencode sessions.
>
> Since this isn't committed to Git or shared with your team, we recommend using this to specify any personal rules that the LLM should follow.
>
> ### Claude Code Compatibility
>
> For users migrating from Claude Code, OpenCode supports Claude Code's file conventions as fallbacks:
>
> - **Project rules**: `CLAUDE.md` in your project directory (used if no `AGENTS.md` exists)
> - **Global rules**: `~/.claude/CLAUDE.md` (used if no `~/.config/opencode/AGENTS.md` exists)
> - **Skills**: `~/.claude/skills/` — see [Agent Skills](/docs/skills/) for details
>
> To disable Claude Code compatibility, set one of these environment variables:
>
> ```bash
> export OPENCODE_DISABLE_CLAUDE_CODE=1        # Disable all .claude support
> export OPENCODE_DISABLE_CLAUDE_CODE_PROMPT=1 # Disable only ~/.claude/CLAUDE.md
> export OPENCODE_DISABLE_CLAUDE_CODE_SKILLS=1 # Disable only .claude/skills
> ```
>
> ---
>
>
> ### Project
>
> Place an `AGENTS.md` in your project root for project-specific rules. These only apply when you are working in this directory or its sub-directories.
>
> ### Global
>
> You can also have global rules in a `~/.config/opencode/AGENTS.md` file. This gets applied across all opencode sessions.
>
> Since this isn't committed to Git or shared with your team, we recommend using this to specify any personal rules that the LLM should follow.
>
> ### Claude Code Compatibility
>
> For users migrating from Claude Code, OpenCode supports Claude Code's file conventions as fallbacks:
>
> - **Project rules**: `CLAUDE.md` in your project directory (used if no `AGENTS.md` exists)
> - **Global rules**: `~/.claude/CLAUDE.md` (used if no `~/.config/opencode/AGENTS.md` exists)
> - **Skills**: `~/.claude/skills/` — see [Agent Skills](/docs/skills/) for details
>
> To disable Claude Code compatibility, set one of these environment variables:
>
> ```bash
> export OPENCODE_DISABLE_CLAUDE_CODE=1        # Disable all .claude support
> export OPENCODE_DISABLE_CLAUDE_CODE_PROMPT=1 # Disable only ~/.claude/CLAUDE.md
> export OPENCODE_DISABLE_CLAUDE_CODE_SKILLS=1 # Disable only .claude/skills
> ```
>
> ---
>
>
> Place an `AGENTS.md` in your project root for project-specific rules. These only apply when you are working in this directory or its sub-directories.
>
> ### Global
>
> You can also have global rules in a `~/.config/opencode/AGENTS.md` file. This gets applied across all opencode sessions.
>
> Since this isn't committed to Git or shared with your team, we recommend using this to specify any personal rules that the LLM should follow.
>
> ### Claude Code Compatibility
>
> For users migrating from Claude Code, OpenCode supports Claude Code's file conventions as fallbacks:
>
> - **Project rules**: `CLAUDE.md` in your project directory (used if no `AGENTS.md` exists)
> - **Global rules**: `~/.claude/CLAUDE.md` (used if no `~/.config/opencode/AGENTS.md` exists)
> - **Skills**: `~/.claude/skills/` — see [Agent Skills](/docs/skills/) for details
>
> To disable Claude Code compatibility, set one of these environment variables:
>
> ```bash
> export OPENCODE_DISABLE_CLAUDE_CODE=1        # Disable all .claude support
> export OPENCODE_DISABLE_CLAUDE_CODE_PROMPT=1 # Disable only ~/.claude/CLAUDE.md
> export OPENCODE_DISABLE_CLAUDE_CODE_SKILLS=1 # Disable only .claude/skills
> ```
>
> ---
>
>
> ### Global
>
> You can also have global rules in a `~/.config/opencode/AGENTS.md` file. This gets applied across all opencode sessions.
>
> Since this isn't committed to Git or shared with your team, we recommend using this to specify any personal rules that the LLM should follow.
>
> ### Claude Code Compatibility
>
> For users migrating from Claude Code, OpenCode supports Claude Code's file conventions as fallbacks:
>
> - **Project rules**: `CLAUDE.md` in your project directory (used if no `AGENTS.md` exists)
> - **Global rules**: `~/.claude/CLAUDE.md` (used if no `~/.config/opencode/AGENTS.md` exists)
> - **Skills**: `~/.claude/skills/` — see [Agent Skills](/docs/skills/) for details
>
> To disable Claude Code compatibility, set one of these environment variables:
>
> ```bash
> export OPENCODE_DISABLE_CLAUDE_CODE=1        # Disable all .claude support
> export OPENCODE_DISABLE_CLAUDE_CODE_PROMPT=1 # Disable only ~/.claude/CLAUDE.md
> export OPENCODE_DISABLE_CLAUDE_CODE_SKILLS=1 # Disable only .claude/skills
> ```
>
> ---
>
>
> You can also have global rules in a `~/.config/opencode/AGENTS.md` file. This gets applied across all opencode sessions.
>
> Since this isn't committed to Git or shared with your team, we recommend using this to specify any personal rules that the LLM should follow.
>
> ### Claude Code Compatibility
>
> For users migrating from Claude Code, OpenCode supports Claude Code's file conventions as fallbacks:
>
> - **Project rules**: `CLAUDE.md` in your project directory (used if no `AGENTS.md` exists)
> - **Global rules**: `~/.claude/CLAUDE.md` (used if no `~/.config/opencode/AGENTS.md` exists)
> - **Skills**: `~/.claude/skills/` — see [Agent Skills](/docs/skills/) for details
>
> To disable Claude Code compatibility, set one of these environment variables:
>
> ```bash
> export OPENCODE_DISABLE_CLAUDE_CODE=1        # Disable all .claude support
> export OPENCODE_DISABLE_CLAUDE_CODE_PROMPT=1 # Disable only ~/.claude/CLAUDE.md
> export OPENCODE_DISABLE_CLAUDE_CODE_SKILLS=1 # Disable only .claude/skills
> ```
>
> ---
>
>
> - **Project rules**: `CLAUDE.md` in your project directory (used if no `AGENTS.md` exists)
> - **Global rules**: `~/.claude/CLAUDE.md` (used if no `~/.config/opencode/AGENTS.md` exists)
> - **Skills**: `~/.claude/skills/` — see [Agent Skills](/docs/skills/) for details
>
> To disable Claude Code compatibility, set one of these environment variables:
>
> ```bash
> export OPENCODE_DISABLE_CLAUDE_CODE=1        # Disable all .claude support
> export OPENCODE_DISABLE_CLAUDE_CODE_PROMPT=1 # Disable only ~/.claude/CLAUDE.md
> export OPENCODE_DISABLE_CLAUDE_CODE_SKILLS=1 # Disable only .claude/skills
> ```
>
> ---
>
>
> - **Global rules**: `~/.claude/CLAUDE.md` (used if no `~/.config/opencode/AGENTS.md` exists)
> - **Skills**: `~/.claude/skills/` — see [Agent Skills](/docs/skills/) for details
>
> To disable Claude Code compatibility, set one of these environment variables:
>
> ```bash
> export OPENCODE_DISABLE_CLAUDE_CODE=1        # Disable all .claude support
> export OPENCODE_DISABLE_CLAUDE_CODE_PROMPT=1 # Disable only ~/.claude/CLAUDE.md
> export OPENCODE_DISABLE_CLAUDE_CODE_SKILLS=1 # Disable only .claude/skills
> ```
>
> ---
>
>
> 1. **Local files** by traversing up from the current directory (`AGENTS.md`, `CLAUDE.md`)
> 2. **Global file** at `~/.config/opencode/AGENTS.md`
> 3. **Claude Code file** at `~/.claude/CLAUDE.md` (unless disabled)
>
> The first matching file wins in each category. For example, if you have both `AGENTS.md` and `CLAUDE.md`, only `AGENTS.md` is used. Similarly, `~/.config/opencode/AGENTS.md` takes precedence over `~/.claude/CLAUDE.md`.
>
> ---
>
>
> 2. **Global file** at `~/.config/opencode/AGENTS.md`
> 3. **Claude Code file** at `~/.claude/CLAUDE.md` (unless disabled)
>
> The first matching file wins in each category. For example, if you have both `AGENTS.md` and `CLAUDE.md`, only `AGENTS.md` is used. Similarly, `~/.config/opencode/AGENTS.md` takes precedence over `~/.claude/CLAUDE.md`.
>
> ---
>
>
> The first matching file wins in each category. For example, if you have both `AGENTS.md` and `CLAUDE.md`, only `AGENTS.md` is used. Similarly, `~/.config/opencode/AGENTS.md` takes precedence over `~/.claude/CLAUDE.md`.
>
> ---
>
>
> You can specify custom instruction files in your `opencode.json` or the global `~/.config/opencode/opencode.json`. This allows you and your team to reuse existing rules rather than having to duplicate them to AGENTS.md.
>
> Example:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "instructions": ["CONTRIBUTING.md", "docs/guidelines.md", ".cursor/rules/*.md"]
> }
> ```
>
> You can also use remote URLs to load instructions from the web.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "instructions": ["https://raw.githubusercontent.com/my-org/shared-rules/main/style.md"]
> }
> ```
>
> Remote instructions are fetched with a 5 second timeout.
>
> All instruction files are combined with your `AGENTS.md` files.
>
> ---
>
>
> All instruction files are combined with your `AGENTS.md` files.
>
> ---
>
>
> While opencode doesn't automatically parse file references in `AGENTS.md`, you can achieve similar functionality in two ways:
>
> ### Using opencode.json
>
> The recommended approach is to use the `instructions` field in `opencode.json`:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "instructions": ["docs/development-standards.md", "test/testing-guidelines.md", "packages/*/AGENTS.md"]
> }
> ```
>
> ### Manual Instructions in AGENTS.md
>
> You can teach opencode to read external files by providing explicit instructions in your `AGENTS.md`. Here's a practical example:
>
> ```markdown title="AGENTS.md"
> # TypeScript Project Rules
>
>
>   "instructions": ["docs/development-standards.md", "test/testing-guidelines.md", "packages/*/AGENTS.md"]
> }
> ```
>
> ### Manual Instructions in AGENTS.md
>
> You can teach opencode to read external files by providing explicit instructions in your `AGENTS.md`. Here's a practical example:
>
> ```markdown title="AGENTS.md"
> # TypeScript Project Rules
>
>
> ### Manual Instructions in AGENTS.md
>
> You can teach opencode to read external files by providing explicit instructions in your `AGENTS.md`. Here's a practical example:
>
> ```markdown title="AGENTS.md"
> # TypeScript Project Rules
>
>
> You can teach opencode to read external files by providing explicit instructions in your `AGENTS.md`. Here's a practical example:
>
> ```markdown title="AGENTS.md"
> # TypeScript Project Rules
>
>
> ```markdown title="AGENTS.md"
> # TypeScript Project Rules
>
>
> - Keep AGENTS.md concise while referencing detailed guidelines
> - Ensure opencode loads files only when needed for the specific task
>
> :::tip
> For monorepos or projects with shared standards, using `opencode.json` with glob patterns (like `packages/*/AGENTS.md`) is more maintainable than manual instructions.
> :::
>
> For monorepos or projects with shared standards, using `opencode.json` with glob patterns (like `packages/*/AGENTS.md`) is more maintainable than manual instructions.
> :::

### Source: OpenCode Plugins — Plugin directories

> ### From local files
>
> Place JavaScript or TypeScript files in the plugin directory.
>
> - `.opencode/plugins/` - Project-level plugins
> - `~/.config/opencode/plugins/` - Global plugins
>
> Files in these directories are automatically loaded at startup.
>
> ---
>
> ### From npm
>
> Specify npm packages in your config file.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "plugin": ["opencode-helicone-session", "opencode-wakatime", "@my-org/custom-plugin"]
> }
> ```
>
> Both regular and scoped npm packages are supported.
>
> Browse available plugins in the [ecosystem](/docs/ecosystem#plugins).
>
> ---
>
> ### How plugins are installed
>
> **npm plugins** are installed automatically using Bun at startup. Packages and their dependencies are cached in `~/.cache/opencode/node_modules/`.
>
> **Local plugins** are loaded directly from the plugin directory. To use external packages, you must create a `package.json` within your config directory (see [Dependencies](#dependencies)), or publish the plugin to npm and [add it to your config](/docs/config#plugins).
>
> ---
>
> ### Load order
>
> Plugins are loaded from all sources and all hooks run in sequence. The load order is:
>
> 1. Global config (`~/.config/opencode/opencode.json`)
> 2. Project config (`opencode.json`)
> 3. Global plugin directory (`~/.config/opencode/plugins/`)
> 4. Project plugin directory (`.opencode/plugins/`)
>
> Duplicate npm packages with the same name and version are loaded once. However, a local plugin and an npm plugin with similar names are both loaded separately.
>
> ---
>
>
> - `.opencode/plugins/` - Project-level plugins
> - `~/.config/opencode/plugins/` - Global plugins
>
> Files in these directories are automatically loaded at startup.
>
> ---
>
> ### From npm
>
> Specify npm packages in your config file.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "plugin": ["opencode-helicone-session", "opencode-wakatime", "@my-org/custom-plugin"]
> }
> ```
>
> Both regular and scoped npm packages are supported.
>
> Browse available plugins in the [ecosystem](/docs/ecosystem#plugins).
>
> ---
>
> ### How plugins are installed
>
> **npm plugins** are installed automatically using Bun at startup. Packages and their dependencies are cached in `~/.cache/opencode/node_modules/`.
>
> **Local plugins** are loaded directly from the plugin directory. To use external packages, you must create a `package.json` within your config directory (see [Dependencies](#dependencies)), or publish the plugin to npm and [add it to your config](/docs/config#plugins).
>
> ---
>
> ### Load order
>
> Plugins are loaded from all sources and all hooks run in sequence. The load order is:
>
> 1. Global config (`~/.config/opencode/opencode.json`)
> 2. Project config (`opencode.json`)
> 3. Global plugin directory (`~/.config/opencode/plugins/`)
> 4. Project plugin directory (`.opencode/plugins/`)
>
> Duplicate npm packages with the same name and version are loaded once. However, a local plugin and an npm plugin with similar names are both loaded separately.
>
> ---
>
>
> 4. Project plugin directory (`.opencode/plugins/`)
>
> Duplicate npm packages with the same name and version are loaded once. However, a local plugin and an npm plugin with similar names are both loaded separately.
>
> ---
>
>
> ```ts title=".opencode/plugins/my-plugin.ts"
> import { escape } from "shescape"
>
> export const MyPlugin = async (ctx) => {
>   return {
>     "tool.execute.before": async (input, output) => {
>       if (input.tool === "bash") {
>         output.args.command = escape(output.args.command)
>       }
>     },
>   }
> }
> ```
>
> ---
>
> ### Basic structure
>
> ```js title=".opencode/plugins/example.js"
> export const MyPlugin = async ({ project, client, $, directory, worktree }) => {
>   console.log("Plugin initialized!")
>
>   return {
>     // Hook implementations go here
>   }
> }
> ```
>
> The plugin function receives:
>
> - `project`: The current project information.
> - `directory`: The current working directory.
> - `worktree`: The git worktree path.
> - `client`: An opencode SDK client for interacting with the AI.
> - `$`: Bun's [shell API](https://bun.com/docs/runtime/shell) for executing commands.
>
> ---
>
> ### TypeScript support
>
> For TypeScript plugins, you can import types from the plugin package:
>
> ```ts title="my-plugin.ts" {1}
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const MyPlugin: Plugin = async ({ project, client, $, directory, worktree }) => {
>   return {
>     // Type-safe hook implementations
>   }
> }
> ```
>
> ---
>
> ### Events
>
> Plugins can subscribe to events as seen below in the Examples section. Here is a list of the different events available.
>
> #### Command Events
>
> - `command.executed`
>
> #### File Events
>
> - `file.edited`
> - `file.watcher.updated`
>
> #### Installation Events
>
> - `installation.updated`
>
> #### LSP Events
>
> - `lsp.client.diagnostics`
> - `lsp.updated`
>
> #### Message Events
>
> - `message.part.removed`
> - `message.part.updated`
> - `message.removed`
> - `message.updated`
>
> #### Permission Events
>
> - `permission.asked`
> - `permission.replied`
>
> #### Server Events
>
> - `server.connected`
>
> #### Session Events
>
> - `session.created`
> - `session.compacted`
> - `session.deleted`
> - `session.diff`
> - `session.error`
> - `session.idle`
> - `session.status`
> - `session.updated`
>
> #### Todo Events
>
> - `todo.updated`
>
> #### Shell Events
>
> - `shell.env`
>
> #### Tool Events
>
> - `tool.execute.after`
> - `tool.execute.before`
>
> #### TUI Events
>
> - `tui.prompt.append`
> - `tui.command.execute`
> - `tui.toast.show`
>
> ---
>
>
> ```js title=".opencode/plugins/example.js"
> export const MyPlugin = async ({ project, client, $, directory, worktree }) => {
>   console.log("Plugin initialized!")
>
>   return {
>     // Hook implementations go here
>   }
> }
> ```
>
> The plugin function receives:
>
> - `project`: The current project information.
> - `directory`: The current working directory.
> - `worktree`: The git worktree path.
> - `client`: An opencode SDK client for interacting with the AI.
> - `$`: Bun's [shell API](https://bun.com/docs/runtime/shell) for executing commands.
>
> ---
>
> ### TypeScript support
>
> For TypeScript plugins, you can import types from the plugin package:
>
> ```ts title="my-plugin.ts" {1}
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const MyPlugin: Plugin = async ({ project, client, $, directory, worktree }) => {
>   return {
>     // Type-safe hook implementations
>   }
> }
> ```
>
> ---
>
> ### Events
>
> Plugins can subscribe to events as seen below in the Examples section. Here is a list of the different events available.
>
> #### Command Events
>
> - `command.executed`
>
> #### File Events
>
> - `file.edited`
> - `file.watcher.updated`
>
> #### Installation Events
>
> - `installation.updated`
>
> #### LSP Events
>
> - `lsp.client.diagnostics`
> - `lsp.updated`
>
> #### Message Events
>
> - `message.part.removed`
> - `message.part.updated`
> - `message.removed`
> - `message.updated`
>
> #### Permission Events
>
> - `permission.asked`
> - `permission.replied`
>
> #### Server Events
>
> - `server.connected`
>
> #### Session Events
>
> - `session.created`
> - `session.compacted`
> - `session.deleted`
> - `session.diff`
> - `session.error`
> - `session.idle`
> - `session.status`
> - `session.updated`
>
> #### Todo Events
>
> - `todo.updated`
>
> #### Shell Events
>
> - `shell.env`
>
> #### Tool Events
>
> - `tool.execute.after`
> - `tool.execute.before`
>
> #### TUI Events
>
> - `tui.prompt.append`
> - `tui.command.execute`
> - `tui.toast.show`
>
> ---
>
>
> ```js title=".opencode/plugins/notification.js"
> export const NotificationPlugin = async ({ project, client, $, directory, worktree }) => {
>   return {
>     event: async ({ event }) => {
>       // Send notification on session completion
>       if (event.type === "session.idle") {
>         await $`osascript -e 'display notification "Session completed!" with title "opencode"'`
>       }
>     },
>   }
> }
> ```
>
> We are using `osascript` to run AppleScript on macOS. Here we are using it to send notifications.
>
> :::note
> If you’re using the OpenCode desktop app, it can send system notifications automatically when a response is ready or when a session errors.
> :::
>
> ---
>
> ### .env protection
>
> Prevent opencode from reading `.env` files:
>
> ```javascript title=".opencode/plugins/env-protection.js"
> export const EnvProtection = async ({ project, client, $, directory, worktree }) => {
>   return {
>     "tool.execute.before": async (input, output) => {
>       if (input.tool === "read" && output.args.filePath.includes(".env")) {
>         throw new Error("Do not read .env files")
>       }
>     },
>   }
> }
> ```
>
> ---
>
> ### Inject environment variables
>
> Inject environment variables into all shell execution (AI tools and user terminals):
>
> ```javascript title=".opencode/plugins/inject-env.js"
> export const InjectEnvPlugin = async () => {
>   return {
>     "shell.env": async (input, output) => {
>       output.env.MY_API_KEY = "secret"
>       output.env.PROJECT_ROOT = input.cwd
>     },
>   }
> }
> ```
>
> ---
>
> ### Custom tools
>
> Plugins can also add custom tools to opencode:
>
> ```ts title=".opencode/plugins/custom-tools.ts"
> import { type Plugin, tool } from "@opencode-ai/plugin"
>
> export const CustomToolsPlugin: Plugin = async (ctx) => {
>   return {
>     tool: {
>       mytool: tool({
>         description: "This is a custom tool",
>         args: {
>           foo: tool.schema.string(),
>         },
>         async execute(args, context) {
>           const { directory, worktree } = context
>           return `Hello ${args.foo} from ${directory} (worktree: ${worktree})`
>         },
>       }),
>     },
>   }
> }
> ```
>
> The `tool` helper creates a custom tool that opencode can call. It takes a Zod schema function and returns a tool definition with:
>
> - `description`: What the tool does
> - `args`: Zod schema for the tool's arguments
> - `execute`: Function that runs when the tool is called
>
> Your custom tools will be available to opencode alongside built-in tools.
>
> :::note
> If a plugin tool uses the same name as a built-in tool, the plugin tool takes precedence.
> :::
>
> ---
>
> ### Logging
>
> Use `client.app.log()` instead of `console.log` for structured logging:
>
> ```ts title=".opencode/plugins/my-plugin.ts"
> export const MyPlugin = async ({ client }) => {
>   await client.app.log({
>     body: {
>       service: "my-plugin",
>       level: "info",
>       message: "Plugin initialized",
>       extra: { foo: "bar" },
>     },
>   })
> }
> ```
>
> Levels: `debug`, `info`, `warn`, `error`. See [SDK documentation](https://opencode.ai/docs/sdk) for details.
>
> ---
>
> ### Compaction hooks
>
> Customize the context included when a session is compacted:
>
> ```ts title=".opencode/plugins/compaction.ts"
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const CompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Inject additional context into the compaction prompt
>       output.context.push(`
>
> ```javascript title=".opencode/plugins/env-protection.js"
> export const EnvProtection = async ({ project, client, $, directory, worktree }) => {
>   return {
>     "tool.execute.before": async (input, output) => {
>       if (input.tool === "read" && output.args.filePath.includes(".env")) {
>         throw new Error("Do not read .env files")
>       }
>     },
>   }
> }
> ```
>
> ---
>
> ### Inject environment variables
>
> Inject environment variables into all shell execution (AI tools and user terminals):
>
> ```javascript title=".opencode/plugins/inject-env.js"
> export const InjectEnvPlugin = async () => {
>   return {
>     "shell.env": async (input, output) => {
>       output.env.MY_API_KEY = "secret"
>       output.env.PROJECT_ROOT = input.cwd
>     },
>   }
> }
> ```
>
> ---
>
> ### Custom tools
>
> Plugins can also add custom tools to opencode:
>
> ```ts title=".opencode/plugins/custom-tools.ts"
> import { type Plugin, tool } from "@opencode-ai/plugin"
>
> export const CustomToolsPlugin: Plugin = async (ctx) => {
>   return {
>     tool: {
>       mytool: tool({
>         description: "This is a custom tool",
>         args: {
>           foo: tool.schema.string(),
>         },
>         async execute(args, context) {
>           const { directory, worktree } = context
>           return `Hello ${args.foo} from ${directory} (worktree: ${worktree})`
>         },
>       }),
>     },
>   }
> }
> ```
>
> The `tool` helper creates a custom tool that opencode can call. It takes a Zod schema function and returns a tool definition with:
>
> - `description`: What the tool does
> - `args`: Zod schema for the tool's arguments
> - `execute`: Function that runs when the tool is called
>
> Your custom tools will be available to opencode alongside built-in tools.
>
> :::note
> If a plugin tool uses the same name as a built-in tool, the plugin tool takes precedence.
> :::
>
> ---
>
> ### Logging
>
> Use `client.app.log()` instead of `console.log` for structured logging:
>
> ```ts title=".opencode/plugins/my-plugin.ts"
> export const MyPlugin = async ({ client }) => {
>   await client.app.log({
>     body: {
>       service: "my-plugin",
>       level: "info",
>       message: "Plugin initialized",
>       extra: { foo: "bar" },
>     },
>   })
> }
> ```
>
> Levels: `debug`, `info`, `warn`, `error`. See [SDK documentation](https://opencode.ai/docs/sdk) for details.
>
> ---
>
> ### Compaction hooks
>
> Customize the context included when a session is compacted:
>
> ```ts title=".opencode/plugins/compaction.ts"
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const CompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Inject additional context into the compaction prompt
>       output.context.push(`
>
> ```javascript title=".opencode/plugins/inject-env.js"
> export const InjectEnvPlugin = async () => {
>   return {
>     "shell.env": async (input, output) => {
>       output.env.MY_API_KEY = "secret"
>       output.env.PROJECT_ROOT = input.cwd
>     },
>   }
> }
> ```
>
> ---
>
> ### Custom tools
>
> Plugins can also add custom tools to opencode:
>
> ```ts title=".opencode/plugins/custom-tools.ts"
> import { type Plugin, tool } from "@opencode-ai/plugin"
>
> export const CustomToolsPlugin: Plugin = async (ctx) => {
>   return {
>     tool: {
>       mytool: tool({
>         description: "This is a custom tool",
>         args: {
>           foo: tool.schema.string(),
>         },
>         async execute(args, context) {
>           const { directory, worktree } = context
>           return `Hello ${args.foo} from ${directory} (worktree: ${worktree})`
>         },
>       }),
>     },
>   }
> }
> ```
>
> The `tool` helper creates a custom tool that opencode can call. It takes a Zod schema function and returns a tool definition with:
>
> - `description`: What the tool does
> - `args`: Zod schema for the tool's arguments
> - `execute`: Function that runs when the tool is called
>
> Your custom tools will be available to opencode alongside built-in tools.
>
> :::note
> If a plugin tool uses the same name as a built-in tool, the plugin tool takes precedence.
> :::
>
> ---
>
> ### Logging
>
> Use `client.app.log()` instead of `console.log` for structured logging:
>
> ```ts title=".opencode/plugins/my-plugin.ts"
> export const MyPlugin = async ({ client }) => {
>   await client.app.log({
>     body: {
>       service: "my-plugin",
>       level: "info",
>       message: "Plugin initialized",
>       extra: { foo: "bar" },
>     },
>   })
> }
> ```
>
> Levels: `debug`, `info`, `warn`, `error`. See [SDK documentation](https://opencode.ai/docs/sdk) for details.
>
> ---
>
> ### Compaction hooks
>
> Customize the context included when a session is compacted:
>
> ```ts title=".opencode/plugins/compaction.ts"
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const CompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Inject additional context into the compaction prompt
>       output.context.push(`
>
> ```ts title=".opencode/plugins/custom-tools.ts"
> import { type Plugin, tool } from "@opencode-ai/plugin"
>
> export const CustomToolsPlugin: Plugin = async (ctx) => {
>   return {
>     tool: {
>       mytool: tool({
>         description: "This is a custom tool",
>         args: {
>           foo: tool.schema.string(),
>         },
>         async execute(args, context) {
>           const { directory, worktree } = context
>           return `Hello ${args.foo} from ${directory} (worktree: ${worktree})`
>         },
>       }),
>     },
>   }
> }
> ```
>
> The `tool` helper creates a custom tool that opencode can call. It takes a Zod schema function and returns a tool definition with:
>
> - `description`: What the tool does
> - `args`: Zod schema for the tool's arguments
> - `execute`: Function that runs when the tool is called
>
> Your custom tools will be available to opencode alongside built-in tools.
>
> :::note
> If a plugin tool uses the same name as a built-in tool, the plugin tool takes precedence.
> :::
>
> ---
>
> ### Logging
>
> Use `client.app.log()` instead of `console.log` for structured logging:
>
> ```ts title=".opencode/plugins/my-plugin.ts"
> export const MyPlugin = async ({ client }) => {
>   await client.app.log({
>     body: {
>       service: "my-plugin",
>       level: "info",
>       message: "Plugin initialized",
>       extra: { foo: "bar" },
>     },
>   })
> }
> ```
>
> Levels: `debug`, `info`, `warn`, `error`. See [SDK documentation](https://opencode.ai/docs/sdk) for details.
>
> ---
>
> ### Compaction hooks
>
> Customize the context included when a session is compacted:
>
> ```ts title=".opencode/plugins/compaction.ts"
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const CompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Inject additional context into the compaction prompt
>       output.context.push(`
>
> ```ts title=".opencode/plugins/my-plugin.ts"
> export const MyPlugin = async ({ client }) => {
>   await client.app.log({
>     body: {
>       service: "my-plugin",
>       level: "info",
>       message: "Plugin initialized",
>       extra: { foo: "bar" },
>     },
>   })
> }
> ```
>
> Levels: `debug`, `info`, `warn`, `error`. See [SDK documentation](https://opencode.ai/docs/sdk) for details.
>
> ---
>
> ### Compaction hooks
>
> Customize the context included when a session is compacted:
>
> ```ts title=".opencode/plugins/compaction.ts"
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const CompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Inject additional context into the compaction prompt
>       output.context.push(`
>
> ```ts title=".opencode/plugins/compaction.ts"
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const CompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Inject additional context into the compaction prompt
>       output.context.push(`
>
> ```ts title=".opencode/plugins/custom-compaction.ts"
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const CustomCompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Replace the entire compaction prompt
>       output.prompt = `
> You are generating a continuation prompt for a multi-agent swarm session.
>
> Summarize:
> 1. The current task and its status
> 2. Which files are being modified and by whom
> 3. Any blockers or dependencies between agents
> 4. The next steps to complete the work
>
> Format as a structured prompt that a new agent can use to resume work.
> `
>     },
>   }
> }
> ```
>
> When `output.prompt` is set, it completely replaces the default compaction prompt. The `output.context` array is ignored in this case.

### Source: OpenCode Skills — Skill paths

> ## Place files
>
> Create one folder per skill name and put a `SKILL.md` inside it.
> OpenCode searches these locations:
>
> - Project config: `.opencode/skills/<name>/SKILL.md`
> - Global config: `~/.config/opencode/skills/<name>/SKILL.md`
> - Project Claude-compatible: `.claude/skills/<name>/SKILL.md`
> - Global Claude-compatible: `~/.claude/skills/<name>/SKILL.md`
> - Project agent-compatible: `.agents/skills/<name>/SKILL.md`
> - Global agent-compatible: `~/.agents/skills/<name>/SKILL.md`
>
> ---
>
>
> - Project config: `.opencode/skills/<name>/SKILL.md`
> - Global config: `~/.config/opencode/skills/<name>/SKILL.md`
> - Project Claude-compatible: `.claude/skills/<name>/SKILL.md`
> - Global Claude-compatible: `~/.claude/skills/<name>/SKILL.md`
> - Project agent-compatible: `.agents/skills/<name>/SKILL.md`
> - Global agent-compatible: `~/.agents/skills/<name>/SKILL.md`
>
> ---
>
>
> Create `.opencode/skills/git-release/SKILL.md` like this:
>
> ```markdown
> ---
> name: git-release
> description: Create consistent releases and changelogs
> license: MIT
> compatibility: opencode
> metadata:
>   audience: maintainers
>   workflow: github
> ---

### Source: OpenCode Troubleshooting — Storage paths

> ## Storage
>
> opencode stores session data and other application data on disk at:
>
> - **macOS/Linux**: `~/.local/share/opencode/`
> - **Windows**: Press `WIN+R` and paste `%USERPROFILE%\.local\share\opencode`
>
> This directory contains:
>
> - `auth.json` - Authentication data like API keys, OAuth tokens
> - `log/` - Application logs
> - `project/` - Project-specific data like session and message data
>   - If the project is within a Git repo, it is stored in `./<project-slug>/storage/`
>   - If it is not a Git repo, it is stored in `./global/storage/`
>
> ---
