---
primary_sources:
  - id: T2-MEMORIES
    title: "Memories"
    url: "https://learn.chatgpt.com/docs/customization/memories.md"
    section: "Local memory storage; Configure local memories"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Memories

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Memories — Local memory storage

> ## Local memory storage
>
> Codex stores memories under your Codex home directory. By default, that's
> `~/.codex`. See [Config and state locations](https://learn.chatgpt.com/docs/config-file/config-advanced#config-and-state-locations)
> for how Codex uses `CODEX_HOME`.
>
> The main memory files live under `~/.codex/memories/` and include summaries,
> durable entries, recent inputs, and supporting evidence from prior chats.
>
> Treat these files as generated state. You can inspect them when troubleshooting
> or before sharing your Codex home directory, but don't rely on editing them by
> hand as your primary control surface.

### Source: Memories — Configure local memories

> ## Configure local memories
>
> Local Codex memories are off by default. In the ChatGPT desktop app, open
> **Settings > Personalization** and turn on **Enable memories**.
>
> For config-based setup, add the feature flag to `config.toml`:
>
> ```toml
> [features]
> memories = true
> ```
>
> For config file locations and the full list of memory-related settings, see
> [Config basics](https://learn.chatgpt.com/docs/config-file/config-basic) and the [configuration
> reference](https://learn.chatgpt.com/docs/config-file/config-reference).
>
> Common memory-specific settings include:
>
> - `memories.generate_memories`: controls whether newly created chats can be
>   stored as memory-generation inputs.
> - `memories.use_memories`: controls whether Codex injects existing memories into
>   future sessions.
> - `memories.disable_on_external_context`: when `true`, keeps chats that used
>   external context such as MCP tool calls, web search, or tool search out of
>   memory generation. The older `memories.no_memories_if_mcp_or_web_search` key
>   is still accepted as an alias.
> - `memories.min_rate_limit_remaining_percent`: controls the minimum remaining
>   Codex rate-limit percentage required before memory generation starts.
> - `memories.extract_model`: overrides the model used for per-chat memory
>   extraction.
> - `memories.consolidation_model`: overrides the model used for global memory
>   consolidation.
