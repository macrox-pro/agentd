---
primary_sources:
  - id: T1-CONFIG
    title: "Config"
    url: "https://opencode.ai/docs/config.md"
    section: "Locations"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Config locations and precedence

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Config — Format and Locations

> ## Format
>
> OpenCode supports both **JSON** and **JSONC** (JSON with Comments) formats.
>
> ```jsonc title="opencode.jsonc"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "model": "anthropic/claude-sonnet-4-5",
>   "autoupdate": true,
>   "server": {
>     "port": 4096,
>   },
> }
> ```
>
> ---
>
> ---
>
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
