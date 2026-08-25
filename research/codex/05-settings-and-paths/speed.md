---
primary_sources:
  - id: T2-SPEED
    title: "Speed"
    url: "https://learn.chatgpt.com/docs/agent-configuration/speed.md"
    section: "Fast mode"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Speed

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Speed — Fast mode

> ## Fast mode
>
> Codex offers the ability to increase the speed of the model for increased
> credit consumption.
>
> Fast mode increases supported model speed by 1.5x and consumes credits at a
> higher rate than Standard mode. It currently supports GPT-5.6, GPT-5.5, and
> GPT-5.4. GPT-5.6 and GPT-5.5 consume credits at 2.5x the Standard rate;
> GPT-5.4 consumes credits at 2x the Standard rate.
>
> Use `/fast on`, `/fast off`, or `/fast status` in the CLI to change or inspect
> the current setting. You can also persist the default with `service_tier =
> "fast"` plus `[features].fast_mode = true` in `config.toml`. Fast mode is
> available in the ChatGPT desktop app, Codex CLI, and IDE extension when you
> sign in with ChatGPT. Fast mode is a ChatGPT credit feature. With an API key,
> Codex uses API token pricing instead, and ChatGPT credit multipliers don't
> apply. API Priority processing has its own billing rate; for GPT-5.6, it costs
> 2x the Standard API token rate.
