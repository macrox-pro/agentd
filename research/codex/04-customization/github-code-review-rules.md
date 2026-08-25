---
primary_sources:
  - id: T2-GH-REVIEW
    title: "Review GitHub pull requests with Codex"
    url: "https://learn.chatgpt.com/docs/third-party/github.md"
    section: "Customize what Codex reviews"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# GitHub code review rules (AGENTS.md)

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Review GitHub pull requests with Codex — Customize what Codex reviews

> ## Customize what Codex reviews
>
> Codex searches your repository for `AGENTS.md` files and follows the applicable
> code review rules. Add a `## Code Review Rules` section to the file closest to
> the code the rules govern. Use `###` headings to group related checks when
> helpful.
>
> For example, an experiment-reporting service can keep post-exposure behavior
> from changing a comparison cohort:
>
> ```md
> ## Code Review Rules
>
> ### Experiment cohorts
>
> - Do not filter treatment comparisons on post-exposure behavior, including conversion or retention.
>   Safe path: build cohorts from assignment or exposure; report conversion as an outcome.
> ```
>
> Put repository-wide rules in the root `AGENTS.md` and service-specific rules
> in a nested file, such as `services/experiment_reporting/AGENTS.md`. Codex
> applies the root and more-specific guidance that covers each changed file, so
> unrelated changes don't have to carry service-specific context.
>
> Start with two or three concise rules that encode checks reviewers often explain. Useful rules:
>
> - **Focus on consequential, repository-specific behavior.** Describe the
>   compatibility constraint, data boundary, or unsafe side effect to flag and
>   why it matters.
> - **State the safe path or exception.** Give Codex enough context to distinguish
>   a real issue from expected behavior.
> - **Keep rules scoped and durable.** Prefer outcomes over function names that
>   can change, and place guidance near the code it governs.
> - **Leave mechanical checks in CI.** Keep formatting, lint, and other
>   deterministic checks out of review rules.
>
> Open a representative pull request and request a review with `@codex review`.
> Refine the rules based on the findings and feedback you see, and narrow or
> remove guidance that produces noise.
>
> Code review rules guide Codex; they don't replace tests, branch protections, or
> required approvals.
>
> For a one-off focus, add it to your pull request comment:
>
> `@codex review for issues in the database migration`

### Source: Review GitHub pull requests with Codex — Request a Codex review

> ## Request a Codex review
>
> 1. In a pull request comment, mention `@codex review`.
> 2. Wait for Codex to react (👀) and post a review.
>
> Codex posts a review on the pull request, just like a teammate would. In
> GitHub, Codex flags only P0 and P1 issues so review comments stay focused on
> high-priority risks.
