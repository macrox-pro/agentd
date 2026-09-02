---
primary_sources:
  - id: T1-CONFIG
    title: "Config"
    url: "https://opencode.ai/docs/config.md"
    section: "Experimental"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Config reference — experimental

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Config — Deprecated providers and Experimental

> ### Disabled providers
>
> You can disable providers that are loaded automatically through the `disabled_providers` option. This is useful when you want to prevent certain providers from being loaded even if their credentials are available.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "disabled_providers": ["openai", "gemini"]
> }
> ```
>
> :::note
> The `disabled_providers` takes priority over `enabled_providers`.
> :::
>
> The `disabled_providers` option accepts an array of provider IDs. When a provider is disabled:
>
> - It won't be loaded even if environment variables are set.
> - It won't be loaded even if API keys are configured through the `/connect` command.
> - The provider's models won't appear in the model selection list.
>
> ---
>
> ### Enabled providers
>
> You can specify an allowlist of providers through the `enabled_providers` option. When set, only the specified providers will be enabled and all others will be ignored.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "enabled_providers": ["anthropic", "openai"]
> }
> ```
>
> This is useful when you want to restrict OpenCode to only use specific providers rather than disabling them one by one.
>
> :::note
> The `disabled_providers` takes priority over `enabled_providers`.
> :::
>
> If a provider appears in both `enabled_providers` and `disabled_providers`, the `disabled_providers` takes priority for backwards compatibility.
>
> ---
>
> ### Experimental
>
> The `experimental` key contains options that are under active development.
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "experimental": {}
> }
> ```
>
> :::caution
> Experimental options are not stable. They may change or be removed without notice.
> :::
>
> ---
