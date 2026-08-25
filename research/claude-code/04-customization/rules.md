---
primary_sources:
  - id: T1-MEMORY
    title: "Memory"
    url: "https://code.claude.com/docs/en/memory.md"
    section: "Rules"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Rules

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Memory — Rules
>
> Organize rules with `.claude/rules/`
>
> For larger projects, you can organize instructions into multiple files using the `.claude/rules/` directory. This keeps instructions modular and easier for teams to maintain. Rules can also be [scoped to specific file paths](#path-specific-rules), so they only load into context when Claude works with matching files, reducing noise and saving context space.
>
>   Rules load into context every session or when matching files are opened. For task-specific instructions that don't need to be in context all the time, use [skills](/docs/en/skills) instead, which only load when you invoke them or when Claude determines they're relevant to your prompt.
>
> #### Set up rules
>
> Place markdown files in your project's `.claude/rules/` directory. Each file should cover one topic, with a descriptive filename like `testing.md` or `api-design.md`. All `.md` files are discovered recursively, so you can organize rules into subdirectories like `frontend/` or `backend/`:
>
> ```text
> your-project/
> ├── .claude/
> │   ├── CLAUDE.md           # Main project instructions
> │   └── rules/
> │       ├── code-style.md   # Code style guidelines
> │       ├── testing.md      # Testing conventions
> │       └── security.md     # Security requirements
> ```
>
> Rules without [`paths` frontmatter](#path-specific-rules) are loaded at launch with the same priority as `.claude/CLAUDE.md`.
>
> Project rules are skipped if you exclude `project` from [`--setting-sources`](/docs/en/cli-reference). Before v2.1.211, rules that load on demand, including path-scoped rules and rules in nested `.claude/rules/` directories, loaded even when `project` was excluded.
>
> #### Path-specific rules
>
> Rules can be scoped to specific files using YAML frontmatter with the `paths` field. These conditional rules only apply when Claude is working with files matching the specified patterns.
>
> ```markdown
> ---
> paths:
>   - "src/api/**/*.ts"
> ---
>
> # API Development Rules
>
> - All API endpoints must include input validation
> - Use the standard error response format
> - Include OpenAPI documentation comments
> ```
>
> Rules without a `paths` field are loaded unconditionally and apply to all files. Path-scoped rules trigger when Claude reads files matching the pattern, not on every tool use. As of v2.1.198, matching also works when Claude reaches a file through a symlinked path to the project directory, for example in a symlinked checkout.
>
> Use glob patterns in the `paths` field to match files by extension, directory, or any combination:
>
> | Pattern                | Matches                                  |
> | ---------------------- | ---------------------------------------- |
> | `**/*.ts`              | All TypeScript files in any directory    |
> | `src/**/*`             | All files under `src/` directory         |
> | `*.md`                 | Markdown files in the project root       |
> | `src/components/*.tsx` | React components in a specific directory |
>
> You can specify multiple patterns and use brace expansion to match multiple extensions in one pattern:
>
> ```markdown
> ---
> paths:
>   - "src/**/*.{ts,tsx}"
>   - "lib/**/*.ts"
>   - "tests/**/*.test.ts"
> ---
> ```
>
> Each brace group multiplies the number of expanded patterns: `src/*.{ts,tsx}` expands to two patterns, and `{a,b}/{c,d}/*.{ts,tsx}` to eight. To keep expansion bounded, a rule's whole `paths` list shares one budget of 1,000 expanded patterns and 4 MiB, and patterns without braces don't count against it.
>
> Claude Code uses any pattern that would exceed the budget unexpanded, and its literal braces match no files. Before v2.1.217, a `paths` value with many brace groups stalled or crashed the CLI at startup.
>
> Glob syntax treats `[` as the start of a bracket expression such as `[abc]`. A pattern with a `[` that can't be read as a bracket expression, such as `photos [2024/**`, is invalid: it matches nothing, and the rule's other patterns keep working. To match a literal `[` in a file name, escape it as `photos \[2024/**`. Before v2.1.207, one invalid pattern made the Read tool fail for every file the rule was evaluated against, instead of matching nothing.
>
> #### Share rules across projects with symlinks
>
> The `.claude/rules/` directory supports symlinks, so you can maintain a shared set of rules and link them into multiple projects. Symlinks are resolved and loaded normally, and circular symlinks are detected and handled gracefully.
>
> This example links both a shared directory and an individual file:
>
> ```bash
> ln -s ~/shared-claude-rules .claude/rules/shared
> ln -s ~/company-standards/security.md .claude/rules/security.md
> ```
>
> #### User-level rules
>
> Personal rules in `~/.claude/rules/` apply to every project on your machine. Use them for preferences that aren't project-specific:
>
> ```text
> ~/.claude/rules/
> ├── preferences.md    # Your personal coding preferences
> └── workflows.md      # Your preferred workflows
> ```
>
> User-level rules are loaded before project rules, giving project rules higher priority.
>
> ### Manage CLAUDE.md for large teams
>
> For organizations deploying Claude Code across teams, you can centralize instructions and control which CLAUDE.md files are loaded.
>
> #### Deploy organization-wide CLAUDE.md
>
> Organizations can deploy a centrally managed CLAUDE.md that applies to all users on a machine. This file cannot be excluded by individual settings.
>
>   **Create the file at the managed policy location**: * macOS: `/Library/Application Support/ClaudeCode/CLAUDE.md`
>     * Linux and WSL: `/etc/claude-code/CLAUDE.md`
>     * Windows: `C:\Program Files\ClaudeCode\CLAUDE.md`
>
>   **Deploy with your configuration management system**: Use MDM, Group Policy, Ansible, or similar tools to distribute the file across developer machines. See [managed settings](/docs/en/managed-settings) for other organization-wide configuration options.
>
> The `claudeMd` key lets you put managed CLAUDE.md content directly inside `managed-settings.json` instead of deploying a separate file.
>
> **Scope**: every Claude Code session on the machine, in every repository. For repository-specific guidance, commit a project CLAUDE.md instead.
>
> **Precedence**: same as a managed CLAUDE.md file. Loads before user and project CLAUDE.md.
>
> **Where it's honored**: managed and policy settings only. Setting `claudeMd` in user, project, or local settings has no effect.
>
> The example below adds behavioral instructions directly in a managed settings file:
>
> ```json
> {
>   "claudeMd": "Always run `make lint` before committing.\nNever push directly to main."
> }
> ```
>
> A managed CLAUDE.md and [managed settings](/docs/en/managed-settings) serve different purposes. Use settings for technical enforcement and CLAUDE.md for behavioral guidance:
>
> | Concern                                        | Configure in                                              |
> | :--------------------------------------------- | :-------------------------------------------------------- |
> | Block specific tools, commands, or file paths  | Managed settings: `permissions.deny`                      |
> | Enforce sandbox isolation                      | Managed settings: `sandbox.enabled`                       |
> | Environment variables and API provider routing | Managed settings: `env`                                   |
> | Login method and organization restrictions     | Managed settings: `forceLoginMethod`, `forceLoginOrgUUID` |
> | Code style and quality guidelines              | Managed CLAUDE.md                                         |
> | Data handling and compliance reminders         | Managed CLAUDE.md                                         |
> | Behavioral instructions for Claude             | Managed CLAUDE.md                                         |
>
> Settings rules are enforced by the client regardless of what Claude decides to do. CLAUDE.md instructions shape Claude's behavior but are not a hard enforcement layer.
>
> #### Exclude specific CLAUDE.md files
>
> In large monorepos, ancestor CLAUDE.md files may contain instructions that aren't relevant to your work. The `claudeMdExcludes` setting lets you skip specific files by path or glob pattern.
>
> This example excludes a top-level CLAUDE.md and a rules directory from a parent folder. Add it to `.claude/settings.local.json` so the exclusion stays local to your machine:
>
> ```json
> {
>   "claudeMdExcludes": [
>     "**/monorepo/CLAUDE.md",
>     "/home/user/monorepo/other-team/.claude/rules/**"
>   ]
> }
> ```
>
> Patterns are matched against absolute file paths using glob syntax. You can configure `claudeMdExcludes` at any [settings layer](/docs/en/settings#where-settings-live): user, project, local, or managed policy. Arrays merge across layers.
>
> Managed policy CLAUDE.md files cannot be excluded. This ensures organization-wide instructions always apply regardless of individual settings.