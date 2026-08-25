---
primary_sources:
  - id: T2-CUSTOM-PROMPTS
    title: "Custom Prompts"
    url: "https://learn.chatgpt.com/docs/custom-prompts.md"
    section: "deprecated; ~/.codex/prompts"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "deprecated"
---
# Custom prompts (deprecated)

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Custom Prompts — deprecation and path

> Custom prompts are deprecated. Use [skills](https://learn.chatgpt.com/docs/build-skills) for reusable
>   instructions that Codex can invoke explicitly or implicitly.
>
> Custom prompts (deprecated) let you turn Markdown files into reusable prompts that you can invoke as slash commands in both the Codex CLI and the Codex IDE extension.
>
> Custom prompts require explicit invocation and live in your local Codex home directory (for example, `~/.codex`), so they're not shared through your repository. If you want to share a prompt (or want Codex to implicitly invoke it), [use skills](https://learn.chatgpt.com/docs/build-skills).
>
> 1. Create the prompts directory:
>
> ```bash
>    mkdir -p ~/.codex/prompts
> ```
>
> Manage prompts by editing or deleting files under `~/.codex/prompts/`. Codex scans only the top-level Markdown files in that folder, so place each custom prompt directly under `~/.codex/prompts/` rather than in subdirectories.
