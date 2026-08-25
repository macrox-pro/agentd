---
primary_sources:
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "Hooks overview; Where Codex looks for hooks; Turn hooks off"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hooks overview

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex Hooks — Hooks / Where Codex looks for hooks

> # Hooks
>
> > For the complete documentation index, see [llms.txt](https://learn.chatgpt.com/llms.txt). Markdown versions of documentation pages are available by appending `.md` to the page URL.
>
> Hooks are an extensibility framework for Codex. They let you run scripts or MCP
> tools during the agentic loop, enabling features such as:
>
> - Send the chat to a custom logging/analytics engine
> - Scan your team's prompts to block accidentally pasting API keys
> - Summarize chats to create persistent memories automatically
> - Run a custom validation check when a chat turn stops, enforcing standards
> - Customize prompting when in a certain directory
>
> Runtime behavior to keep in mind:
>
> - Matching hooks from multiple files all run.
> - Multiple matching command hooks for the same event are launched concurrently,
>   so one hook can't prevent another matching hook from starting.
> - Non-managed hooks must be reviewed and trusted before they run.
>
> Hooks run at different points in a conversation:
>
> | When                              | Hooks                                                                                                                     |
> | --------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
> | During a turn                     | `PreToolUse`, `PermissionRequest`, `PostToolUse`, `PreCompact`, `PostCompact`, `UserPromptSubmit`, `SubagentStop`, `Stop` |
> | When a session or subagent starts | `SessionStart`, `SubagentStart`                                                                                           |
> | When the main thread ends         | `SessionEnd` (doesn't run for subagents)                                                                                  |
>
> ## Where Codex looks for hooks
>
> Codex discovers hooks next to active config layers in either of these forms:
>
> - `hooks.json`
> - inline `[hooks]` tables inside `config.toml`
>
> Installed plugins can also bundle lifecycle config through their plugin
> manifest or a default `hooks/hooks.json` file. See [Build
> plugins](https://developers.openai.com/plugins/build/plugins#bundled-mcp-servers-and-lifecycle-hooks) for the
> plugin packaging rules.
>
> In practice, the four most useful locations are:
>
> - `~/.codex/hooks.json`
> - `~/.codex/config.toml`
> - `<repo>/.codex/hooks.json`
> - `<repo>/.codex/config.toml`
>
> If more than one hook source exists, Codex loads all matching hooks.
> Higher-precedence config layers don't replace lower-precedence hooks.
> If a single layer contains both `hooks.json` and inline `[hooks]`, Codex
> merges them and warns at startup. Prefer one representation per layer.
>
> Codex can also discover hooks bundled with enabled plugins. Plugin-bundled
> hooks load alongside other hook sources and use the same trust-review flow as
> other non-managed hooks.
>
> Project-local hooks load only when the project `.codex/` layer is trusted. In
> untrusted projects, Codex still loads user and system hooks from their own
> active config layers.

### Source: Codex Hooks — Turn hooks off

> ## Turn hooks off
>
> Hooks are enabled by default. To turn them off in `config.toml`, set:
>
> ```toml
> [features]
> hooks = false
> ```
>
> Use `hooks` as the canonical feature key. `codex_hooks` still works as a
> deprecated alias. Admins can force hooks off the same way in
> `requirements.toml` with `[features].hooks = false`.
