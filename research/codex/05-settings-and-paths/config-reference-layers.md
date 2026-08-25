---
primary_sources:
  - id: T1-CFG-REF
    title: "Configuration Reference"
    url: "https://learn.chatgpt.com/docs/config-file/config-reference.md"
    section: "config.toml opening; projects.<path>.trust_level; requirements.toml intro"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Config reference — layers

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Configuration Reference — `config.toml` (user vs project)

> # Configuration Reference
>
> > For the complete documentation index, see [llms.txt](https://learn.chatgpt.com/llms.txt). Markdown versions of documentation pages are available by appending `.md` to the page URL.
>
> Use this page as a searchable reference for Codex configuration files. For conceptual guidance and examples, start with [Config basics](https://learn.chatgpt.com/docs/config-file/config-basic) and [Advanced Config](https://learn.chatgpt.com/docs/config-file/config-advanced).
>
> ## `config.toml`
>
> User-level configuration lives in `~/.codex/config.toml`. You can also add project-scoped overrides in `.codex/config.toml` files. Codex loads project-scoped config files only when you trust the project.
>
> Project-scoped config can't override machine-local provider, auth,
> host-owned app request metadata, notification, configuration profile selection,
> or telemetry routing keys. Codex ignores `openai_base_url`,
> `chatgpt_base_url`, `apps_mcp_product_sku`, `model_provider`,
> `model_providers`, `notify`, `profile`, `profiles`,
> `experimental_realtime_ws_base_url`, and `otel` when they appear in a
> project-local `.codex/config.toml`; put provider, notification, and telemetry
> keys in user-level config instead. Config [profile files](https://learn.chatgpt.com/docs/config-file/config-advanced#profiles) live next to
> `config.toml` as `$CODEX_HOME/profile-name.config.toml`; select one with
> `--profile profile-name`.
>
> For sandbox and approval keys (`approval_policy`, `sandbox_mode`, and `sandbox_workspace_write.*`), pair this reference with [Sandbox and approvals](https://learn.chatgpt.com/docs/agent-approvals-security#sandbox-and-approvals), [Protected paths in writable roots](https://learn.chatgpt.com/docs/agent-approvals-security#protected-paths-in-writable-roots), and [Network access](https://learn.chatgpt.com/docs/agent-approvals-security#network-access). For beta permission profiles, see [Permissions](https://learn.chatgpt.com/docs/permissions).

### Source: Configuration Reference — `projects.<path>.trust_level`

> <ConfigTable
>   options={[
>     {
>       key: "projects.<path>.trust_level",
>       type: "string",
>       description:
>         'Mark a project or worktree as trusted or untrusted (`"trusted"` | `"untrusted"`). Untrusted projects skip project-scoped `.codex/` layers, including project-local config, hooks, and rules.',
>     }
>   ]}
>   client:load
> />

### Source: Configuration Reference — `requirements.toml` (intro)

> ## `requirements.toml`
>
> `requirements.toml` is an admin-enforced configuration file that constrains security-sensitive settings users can't override. For details, locations, and examples, see [Admin-enforced requirements](https://learn.chatgpt.com/docs/enterprise/managed-configuration#admin-enforced-requirements-requirementstoml).
>
> For ChatGPT Business and Enterprise users, Codex can also apply cloud-fetched
> requirements. See the security page for precedence details.
>
> Use `[features]` in `requirements.toml` to pin runtime feature flags by the same
> canonical keys that `config.toml` uses. Requirements can also include documented
> app-only keys that don't belong in `config.toml`. Omitted keys remain
> unconstrained.

