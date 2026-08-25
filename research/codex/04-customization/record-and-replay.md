---
primary_sources:
  - id: T2-RECORD
    title: "Record & Replay"
    url: "https://learn.chatgpt.com/docs/extend/record-and-replay.md"
    section: "full page"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Record & Replay

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Record & Replay — overview; Start a recording

> Record & Replay is available on macOS. Computer Use must also be available and
>   enabled.
>
> Record & Replay lets you demonstrate a workflow on your
> Mac and turn it into a reusable skill. Use it when the workflow is repetitive,
> depends on your preferences, or is easier to show than to describe in a prompt.

> ## Start a recording
>
> 1. In the ChatGPT desktop app, select ChatGPT and turn on Work in the switcher, or select Codex. Then open **Plugins**.
> 2. Open the **+** menu.
> 3. Select **Record a skill**.
> 4. Review the suggested prompt, add any helpful context, and submit it.
> 5. When the chat asks for permission to record your actions, approve the
>    request once you are ready to demonstrate the workflow.
> 6. Perform the workflow on your Mac.
> 7. When you are done, stop recording from the menu bar or overlay, or tell the
>    chat that you are done.
>
> After you stop recording, ChatGPT or Codex inspects the captured workflow and
> drafts a skill. The skill explains when to use the workflow, what inputs it
> needs, what steps to follow, and how to verify the result. You can also ask for
> further refinements.

### Source: Record & Replay — Troubleshooting

> ### I don't see Record & Replay
>
> If your organization manages Codex with `requirements.toml`, the
> `[features].computer_use` requirement controls Record & Replay too. Setting
> `computer_use = false` makes both features unavailable.
