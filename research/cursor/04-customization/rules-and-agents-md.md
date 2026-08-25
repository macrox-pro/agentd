---
primary_sources:
  - id: T1-RULES
    title: "Full page"
    url: "https://cursor.com/docs/rules.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Rules and AGENTS.md

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Rules and AGENTS.md

> # Rules
>
> Rules provide system-level instructions to Agent. They bundle prompts, scripts, and more together, making it easy to manage and share workflows across your team.
>
> Cursor supports four types of rules:
>
> ### Project Rules
>
> Stored in `.cursor/rules`, version-controlled and scoped to your codebase.
>
> ### User Rules
>
> Global to your Cursor environment. Used by Agent (Chat).
>
> ### Team Rules
>
> Team-wide rules managed from the dashboard. Available on Team and [Enterprise](https://cursor.com/docs/enterprise.md) plans.
>
> ### AGENTS.md
>
> Agent instructions in markdown format. Simple alternative to
> `.cursor/rules`.
>
> ## How rules work
>
> Large language models don't retain memory between completions. Rules provide persistent, reusable context at the prompt level.
>
> When applied, rule contents are included at the start of the model context. This gives the AI consistent guidance for generating code, interpreting edits, or helping with workflows.
>
> ## Project rules
>
> Project rules live in `.cursor/rules` as `.mdc` files and are version-controlled. They are scoped using path patterns, invoked manually, or included based on relevance.
>
> Use project rules to:
>
> - Encode domain-specific knowledge about your codebase
> - Automate project-specific workflows or templates
> - Standardize style or architecture decisions
>
> ### Rule file structure
>
> Each rule is an `.mdc` file that you can name anything you want. Project rules must use the `.mdc` extension. A plain `.md` file in `.cursor/rules` is ignored by the rules system because it has no frontmatter to specify `description`, `globs`, and `alwaysApply`. If you prefer plain markdown, use [AGENTS.md](https://cursor.com/docs/rules.md#agentsmd) instead.
>
> ```bash
> .cursor/rules/
>   react-patterns.mdc       # Recognized as a project rule
>   api-guidelines.md        # Ignored (wrong extension)
>   frontend/                # Organize rules in folders
>     components.mdc
> ```
>
> ### Rule anatomy
>
> Each rule is a markdown file with frontmatter metadata and content. Control how rules are applied from the type dropdown which changes properties `description`, `globs`, `alwaysApply`.
>
> | Rule Type                 | Description                                           |
> | :------------------------ | :---------------------------------------------------- |
> | `Always Apply`            | Apply to every chat session                           |
> | `Apply Intelligently`     | When Agent decides it's relevant based on description |
> | `Apply to Specific Files` | When file matches a specified pattern                 |
> | `Apply Manually`          | When @-mentioned in chat (e.g., `@my-rule`)           |
>
> Under the hood, the three frontmatter fields interact to determine when a rule is included:
>
> | `alwaysApply` | `description` | `globs`  | Behavior                                                         |
> | :------------ | :------------ | :------- | :--------------------------------------------------------------- |
> | `true`        | —             | —        | Always included. Globs and description are ignored.              |
> | `false`       | —             | provided | Auto-attached when a matching file is in context.                |
> | `false`       | provided      | omitted  | Agent reads the description and pulls the rule in when relevant. |
> | `false`       | omitted       | omitted  | Included only when you `@`-mention the rule in chat.             |
>
> ```md title="Always applied"
