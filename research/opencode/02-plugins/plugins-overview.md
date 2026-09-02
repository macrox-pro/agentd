---
primary_sources:
  - id: T1-PLUGINS
    title: "Plugins"
    url: "https://opencode.ai/docs/plugins.md"
    section: "Use a plugin"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Plugins overview

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Plugins — Use and Create (structure)

> ## Use a plugin
>
> There are two ways to load plugins.
>
> ---
>
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
> ---
>
> ## Create a plugin
>
> A plugin is a **JavaScript/TypeScript module** that exports one or more plugin
> functions. Each function receives a context object and returns a hooks object.
>
> ---
>
> ### Dependencies
>
> Local plugins and custom tools can use external npm packages. Add a `package.json` to your config directory with the dependencies you need.
>
> ```json title=".opencode/package.json"
> {
>   "dependencies": {
>     "shescape": "^2.1.0"
>   }
> }
> ```
>
> OpenCode runs `bun install` at startup to install these. Your plugins and tools can then import them.
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
