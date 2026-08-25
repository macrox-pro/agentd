---
primary_sources:
  - id: T1-CLOUD
    title: "Codex cloud"
    url: "https://learn.chatgpt.com/docs/cloud.md"
    section: "hooks mentions"
  - id: T1-ENV-MODES
    title: "Environment modes"
    url: "https://learn.chatgpt.com/docs/environments/modes.md"
    section: "environment placement (no hooks)"
  - id: T1-MANAGED
    title: "Managed configuration"
    url: "https://learn.chatgpt.com/docs/enterprise/managed-configuration.md"
    section: "managed hooks"
also_cited_in: [T2-REMOTE]
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Cloud hooks matrix

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex cloud — hooks mentions

> ### 2. Connect GitHub or GitLab
>
> Connect GitHub or GitLab (Beta) when prompted. For GitHub, choose the repositories Codex can access; for GitLab, select a project when you create the environment. For GitLab setup, webhook permissions, and merge request reviews, see [Use Codex with GitLab (Beta)](https://learn.chatgpt.com/docs/third-party/gitlab).

### Source: Environment modes — environment placement (no hooks mentions)

> In the ChatGPT desktop app, open the ChatGPT dropdown and select **Codex**.
> When starting a Codex chat, choose where it runs:
>
> - **Local**: work directly in your current project directory.
> - **Worktree**: isolate changes in a Git worktree. [Learn more](https://learn.chatgpt.com/docs/environments/git-worktrees).
> - **Cloud**: run remotely in a configured cloud environment.
>
> Both **Local** and **Worktree** chats run on your computer.

### Source: Managed configuration — managed hooks

> Requirements constrain security-sensitive settings (approval policy, approvals reviewer, automatic review policy, sandbox mode, permission profiles, web search mode, managed hooks, which MCP servers users can enable, and which user-configured plugin marketplace sources they can add, install from, or refresh). When resolving configuration (for example from `config.toml`, [profile files](https://learn.chatgpt.com/docs/config-file/config-advanced#profiles), or CLI config overrides), if a value conflicts with an enforced rule, the local client falls back to a compatible value and notifies the user. If you configure an `mcp_servers` allowlist, the client enables an MCP server only when both its name and identity match an approved entry; otherwise, the client disables it.

> Higher-precedence layers override ordinary scalar and list values from lower
> layers. Tables merge by key, while requirements such as rules, hooks, and
> filesystem restrictions have field-specific composition behavior. Use the
> [`requirements.toml` reference](https://learn.chatgpt.com/docs/config-file/config-reference#requirementstoml)
> for the current schema instead of assuming that every field merges the same
> way.

> ### Enforce managed hooks from requirements
>
> Admins can also define managed lifecycle hooks directly in `requirements.toml`.
> Use `[hooks]` for the hook configuration itself, and point `managed_dir` at the
> directory where your MDM or endpoint-management tooling installs the referenced
> scripts.
>
> To enforce managed hooks even for users who turned hooks off locally, pin
> `[features].hooks = true` alongside `[hooks]`. To skip user, project, session,
> and plugin hooks while still allowing managed hooks, set
> `allow_managed_hooks_only = true`.
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
> Notes:
>
> - The local runtime enforces the hook configuration from `requirements.toml`,
>   but it doesn't distribute the scripts in `managed_dir`.
> - Deliver those scripts with your MDM or device-management solution.
> - Managed hook commands should reference absolute script paths under the
>   configured managed directory.
> - `allow_managed_hooks_only = true` skips hooks from user, project, session, and
>   plugin sources, but still loads hooks from `requirements.toml` and other
>   managed config layers.

Research note: Official docs reviewed for this dump do not state whether Codex cloud loads user-level ~/.codex/hooks.json. Local and Worktree modes run on the user's computer per environments/modes.

### Source: Codex Remote — connected host approvals

> ## Start, guide, and review coding tasks from your phone
>
> Follow progress, approve actions, and send instructions from your phone. Codex runs each task on your connected computer.
>
> Use the ChatGPT mobile app with a connected Mac or Windows PC. Availability depends on rollout and your workspace settings.

> ## Codex Remote advantages
>
> - **Start tasks from your phone:** Choose a connected computer and project, describe the task, and let Codex get to work.
> - **Guide work as it happens:** Open a task, follow its progress, and send new instructions without returning to your desk.
> - **Approve requested actions:** Review requested commands and actions before Codex continues on your connected computer.
> - **Review the result:** Inspect responses, changed files, diffs, and test results, then decide what happens next.

> Start, approve, and review tasks from your phone. Your connected computer runs the work under your organization’s security policies.

> ### 2. Approve requests
>
> Review commands and requested actions before Codex continues working on your connected computer.

