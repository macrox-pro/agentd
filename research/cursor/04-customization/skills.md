---
primary_sources:
  - id: T1-SKILLS
    title: "Full page"
    url: "https://cursor.com/docs/skills.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Agent Skills

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Agent Skills

> # Agent Skills
>
> Agent Skills is an open standard for extending AI agents with specialized capabilities. Skills package domain-specific knowledge and workflows that agents can use to perform specific tasks.
>
> ## What are skills?
>
> A skill is a portable, version-controlled package that teaches agents how to perform domain-specific tasks. Skills can include scripts, templates, and references that agents may act on using their tools.
>
> ### Portable
>
> Skills work across any agent that supports the Agent Skills standard.
>
> ### Version-controlled
>
> Skills are stored as files and can be tracked in your repository, or installed via GitHub repository links.
>
> ### Actionable
>
> Skills can include scripts, templates, and references that agents act on using their tools.
>
> ### Progressive
>
> Skills load resources on demand, keeping context usage efficient.
>
> ## How skills work
>
> When Cursor starts, it automatically discovers skills from skill directories and makes them available to Agent. The agent is presented with available skills and decides when they are relevant based on context.
>
> Skills can also be manually invoked by typing `/` in Agent chat and searching for the skill name. A skill invoked this way attaches to one message. To keep a skill on for the whole session, use it as a Custom Mode with Option+Enter (Mac) or Alt+Enter (Windows). See [Custom Modes](https://cursor.com/docs/agent/prompting.md#custom-modes).
>
> ## Built-in Cursor skills
>
> Cursor includes a small set of built in skills to improve your general workflows. These skills are managed by Cursor and appear alongside the skills you add yourself.
>
> | Skill                     | What it does                                                                                         |
> | ------------------------- | ---------------------------------------------------------------------------------------------------- |
> | `/automate`               | Creates Cursor Automations triggered by schedules, Slack messages, GitHub events, and other sources. |
> | `/babysit`                | Monitors a pull request and addresses feedback, conflicts, failing checks, and follow-up work.       |
> | `/canvas`                 | Creates interactive React artifacts that render alongside the conversation.                          |
> | `/create-hook`            | Creates Cursor hooks and updates `hooks.json` for agent lifecycle events.                            |
> | `/create-rule`            | Creates Cursor rules with the appropriate scope and instructions.                                    |
> | `/create-skill`           | Creates Agent Skills, including their structure and `SKILL.md` files.                                |
> | `/create-subagent`        | Creates custom subagents with focused roles and delegation instructions.                             |
> | `/cursor-blame`           | Investigates AI-authored changes and the prompts that produced them.                                 |
> | `/loop`                   | Runs a prompt or skill repeatedly at a specified interval.                                           |
> | `/migrate-to-skills`      | Converts eligible dynamic rules and slash commands into Agent Skills.                                |
> | `/review`                 | Selects and runs the appropriate code-review agent.                                                  |
> | `/review-bugbot`          | Reviews code for likely bugs and regressions with Bugbot.                                            |
> | `/review-security`        | Reviews code for security vulnerabilities with Security Review.                                      |
> | `/sdk`                    | Helps you build applications and integrations with the Cursor SDK.                                   |
> | `/shell`                  | Runs the provided text as a literal shell command.                                                   |
> | `/split-to-prs`           | Splits large changes into smaller pull requests.                                                     |
> | `/statusline`             | Configures the Cursor CLI status line.                                                               |
> | `/update-cli-config`      | Updates Cursor CLI settings in `~/.cursor/cli-config.json`.                                          |
> | `/update-cursor-settings` | Finds and updates the appropriate Cursor or VS Code setting.                                         |
>
> You can run any built-in skill by typing `/` in Agent chat and selecting its name. Agent may also use some built-in skills automatically when your request clearly matches their purpose.
>
> ## Skill directories
>
> Skills are automatically loaded from these locations:
>
> | Location            | Scope               |
> | ------------------- | ------------------- |
> | `.agents/skills/`   | Project-level       |
> | `.cursor/skills/`   | Project-level       |
> | `~/.agents/skills/` | User-level (global) |
> | `~/.cursor/skills/` | User-level (global) |
>
> For compatibility, Cursor also loads skills from Claude and Codex directories: `.claude/skills/`, `.codex/skills/`, `~/.claude/skills/`, and `~/.codex/skills/`.
>
> Each skill should be a folder containing a `SKILL.md` file:
>
> ```text
> .agents/
> └── skills/
>     └── my-skill/
>         └── SKILL.md
> ```
>
> Skills can also include optional directories for scripts, references, and assets:
>
> ```text
> .agents/
> └── skills/
>     └── deploy-app/
>         ├── SKILL.md
>         ├── scripts/
>         │   ├── deploy.sh
>         │   └── validate.py
>         ├── references/
>         │   └── REFERENCE.md
>         └── assets/
>             └── config-template.json
> ```
>
> ### Nested skill directories
>
> Skill directories can be organized into subdirectories. This is useful for grouping related skills by category, team, or domain. Cursor walks the skills root recursively and picks up any `SKILL.md` it finds:
>
> ```text
> .cursor/
> └── skills/
>     ├── shipping/
>     │   ├── land-it/
>     │   │   └── SKILL.md
>     │   └── careful-merge-conflicts/
>     │       └── SKILL.md
>     ├── debugging/
>     │   └── using-datadog-mcp/
>     │       └── SKILL.md
>     └── workflow/
>         └── tdd/
>             └── SKILL.md
> ```
>
> The category folder is purely organizational. The skill's identity comes from the folder containing `SKILL.md` (here `land-it`, `tdd`, etc.), not the parent category.
>
> Cursor also discovers skills inside nested project subdirectories. A `.cursor/skills/` (or `.agents/skills/`) folder anywhere inside your repository is picked up, so monorepos can colocate skills with the package they apply to:
>
> ```text
> my-monorepo/
> ├── .cursor/skills/         # repo-wide skills
> │   └── land-it/SKILL.md
> └── apps/
>     └── web/
>         └── .cursor/skills/  # app-specific skills
>             └── deploy-web/SKILL.md
> ```
>
> Skills in nested project directories are automatically scoped to files inside that directory. In the example above, `deploy-web` is only surfaced when the agent works with files under `apps/web/`, while skills in the repo-wide `.cursor/skills/` are available everywhere. This is similar to the [`paths` frontmatter field](https://cursor.com/docs/skills.md#scoping-a-skill-to-specific-files) — you don't need to set `paths` on a nested skill to scope it to its directory.
>
> ## SKILL.md file format
>
> Each skill is defined in a `SKILL.md` file with YAML frontmatter:
>
> ```markdown
