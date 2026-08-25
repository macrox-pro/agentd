---
primary_sources:
  - id: T1-CFG-ADV
    title: "Advanced Configuration"
    url: "https://learn.chatgpt.com/docs/config-file/config-advanced.md"
    section: "Config and state locations; History persistence"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Sessions and history

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Advanced Configuration — Config and state locations

> ## Config and state locations
>
> Codex stores its local state under `CODEX_HOME` (defaults to `~/.codex`).
>
> Common files you may see there:
>
> - `config.toml` (your local configuration)
> - `auth.json` (if you use file-based credential storage) or your OS keychain/keyring
> - `history.jsonl` (if history persistence is enabled)
> - Other per-user state such as logs and caches
>
> For authentication details (including credential storage modes), see [Authentication](https://learn.chatgpt.com/docs/auth). For the full list of configuration keys, see [Configuration Reference](https://learn.chatgpt.com/docs/config-file/config-reference).
>
> For shared defaults, rules, and skills checked into repos or system paths, see [Team Config](https://learn.chatgpt.com/docs/enterprise/admin-setup#step-4-standardize-local-configuration-with-team-config).
>
> If you just need to point the built-in OpenAI provider at an LLM proxy, router, or data-residency enabled project, set `openai_base_url` in `config.toml` instead of defining a new provider. This changes the base URL for the built-in `openai` provider without requiring a separate `model_providers.<id>` entry.
>
> ```toml
> openai_base_url = "https://us.api.openai.com/v1"
> ```

### Source: Advanced Configuration — History persistence

> ## History persistence
>
> By default, Codex saves local session transcripts under `CODEX_HOME` (for example, `~/.codex/history.jsonl`). To disable local history persistence:
>
> ```toml
> [history]
> persistence = "none"
> ```
>
> To cap the history file size, set `history.max_bytes`. When the file exceeds the cap, Codex drops the oldest entries and compacts the file while keeping the newest records.
>
> ```toml
> [history]
> max_bytes = 104857600 # 100 MiB
> ```

## Research note

Official docs describe history.jsonl under CODEX_HOME; nested sessions/YYYY/MM/DD/rollout-*-{id}.jsonl layout not found in this source.
