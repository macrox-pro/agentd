---
primary_sources:
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "Managed hooks from requirements.toml; Plugin-bundled hooks"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Managed and plugin-bundled hooks

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex Hooks — Managed hooks from requirements.toml

> ## Managed hooks from `requirements.toml`
>
> Enterprise-managed requirements can also define hooks inline under `[hooks]`.
> This is useful when admins want to enforce the hook configuration while
> delivering the actual scripts through MDM or another device-management system.
> To enforce managed hooks even for users who disabled hooks locally, pin
> `[features].hooks = true` in `requirements.toml` alongside `[hooks]`. To ignore
> user, project, session, and plugin hooks while still allowing administrator
> managed hooks, set `allow_managed_hooks_only = true`.
>
> ```toml
> allow_managed_hooks_only = true
>
> [features]
> hooks = true
>
> [hooks]
> managed_dir = "/enterprise/hooks"
> windows_managed_dir = 'C:\enterprise\hooks'
>
> [[hooks.PreToolUse]]
> matcher = "^Bash$"
>
> [[hooks.PreToolUse.hooks]]
> type = "command"
> command = "python3 /enterprise/hooks/pre_tool_use_policy.py"
> command_windows = 'py -3 C:\enterprise\hooks\pre_tool_use_policy.py'
> timeout = 30
> statusMessage = "Checking managed Bash command"
> ```
>
> Notes for managed hooks:
>
> - `managed_dir` is used on macOS and Linux.
> - `windows_managed_dir` is used on Windows.
> - Codex doesn't distribute the scripts in `managed_dir`; your enterprise
>   tooling must install and update them separately.
> - Managed hook commands should use absolute script paths under the configured
>   managed directory.
> - `allow_managed_hooks_only = true` skips hooks from user, project, session, and
>   plugin sources, but still loads managed hooks from `requirements.toml` and
>   other managed config layers.

### Source: Codex Hooks — Plugin-bundled hooks

> ## Plugin-bundled hooks
>
> When a plugin is enabled, Codex can load lifecycle hooks from that plugin
> alongside user, project, and managed hooks.
>
> By default, Codex looks for `hooks/hooks.json` inside the plugin root. A plugin
> manifest can override that default with a `hooks` entry in
> `.codex-plugin/plugin.json`. The manifest entry can be a `./`-prefixed path, an
> array of `./`-prefixed paths, an inline hooks object, or an array of inline
> hooks objects.
>
> ```json
> {
>   "name": "repo-policy",
>   "hooks": "./hooks/hooks.json"
> }
> ```
>
> Manifest hook paths are resolved relative to the plugin root and must stay
> inside that root. If a manifest defines `hooks`, Codex uses those manifest
> entries instead of the default `hooks/hooks.json`.
>
> Plugin hook commands receive these environment variables:
>
> - `PLUGIN_ROOT` is a Codex-specific extension that points to the installed
>   plugin root.
> - `PLUGIN_DATA` is a Codex-specific extension that points to the plugin's
>   writable data directory.
> - Codex also sets `CLAUDE_PLUGIN_ROOT` and `CLAUDE_PLUGIN_DATA` for
>   compatibility with existing plugin hooks.
>
> Plugin hooks use the same event schema as other hooks. Installing or enabling a
> plugin doesn't automatically trust its hooks; Codex skips plugin-bundled hooks
> until you review and trust the current hook definition.
