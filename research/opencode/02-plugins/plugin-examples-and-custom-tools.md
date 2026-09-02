---
primary_sources:
  - id: T1-PLUGINS
    title: "Plugins"
    url: "https://opencode.ai/docs/plugins.md"
    section: "Examples"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Plugin examples and custom tools

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Plugins — Examples

> ## Examples
>
> Here are some examples of plugins you can use to extend opencode.
>
> ---
>
> ### Send notifications
>
> Send notifications when certain events occur:
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
