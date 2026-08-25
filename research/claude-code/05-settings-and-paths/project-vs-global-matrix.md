---
primary_sources:
  - id: T1-SETTINGS
    title: "Claude Code settings"
    url: "https://code.claude.com/docs/en/settings.md"
    section: "Settings precedence; Settings files"
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "Hook locations"
  - id: T1-MCP
    title: "MCP"
    url: "https://code.claude.com/docs/en/mcp.md"
    section: "Configuration"
  - id: T1-MEMORY
    title: "Memory"
    url: "https://code.claude.com/docs/en/memory.md"
    section: "CLAUDE.md; rules"
  - id: T1-SKILLS
    title: "Skills"
    url: "https://code.claude.com/docs/en/skills.md"
    section: "Where skills live"
  - id: T1-SUBAGENTS
    title: "Sub-agents"
    url: "https://code.claude.com/docs/en/sub-agents.md"
    section: "Where subagents live"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Project vs global path matrix

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

## Summary table

| Concern | Paths (from docs) | Notes from docs |
| ------- | ----------------- | --------------- |
| settings.json layers | `managed-settings.json` / MDM; CLI `--settings`; `.claude/settings.local.json`; `.claude/settings.json`; `~/.claude/settings.json` | Higher scope overrides lower; array settings merge |
| hooks | Under `"hooks"` key in settings files; plugin `hooks/hooks.json` | No standalone project hooks.json except plugins |
| MCP | `.mcp.json` (project root); `~/.claude.json` (user/local MCP); not in settings.json | Project servers need one-time approval |
| CLAUDE.md | managed policy path; `~/.claude/CLAUDE.md`; `./CLAUDE.md`; `./.claude/CLAUDE.md`; nested subdirs on demand | Loaded at session start or on Read |
| rules | `~/.claude/rules/`; `.claude/rules/*.md` | Path-scoped via frontmatter |
| skills | `.claude/skills/<name>/SKILL.md`; `~/.claude/skills/`; plugin-bundled | Folder + SKILL.md required |
| subagents | `.claude/agents/`; `~/.claude/agents/` | Markdown agent definitions |
| plugins | marketplace install; `.claude-plugin/` for plugin scope | Bundles hooks, skills, MCP |
| env vars | `env` in settings.json; process environment | `CLAUDE_CONFIG_DIR` bypasses ~/.claude |
| transcripts | `~/.claude/projects/**/transcript.jsonl` | Session resume, export, trajectory import |


### Source: Hooks reference — Hook locations
>
> ### Hook locations
>
> Where you define a hook determines its scope:
>
> | Location                                 | Scope                                                                                                            | Shareable                                             |
> | :--------------------------------------- | :--------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------- |
> | `~/.claude/settings.json`                | All your projects                                                                                                | No, local to your machine                             |
> | `.claude/settings.json`                  | Single project                                                                                                   | Yes, can be committed to the repo                     |
> | `.claude/settings.local.json`            | Single project                                                                                                   | No, gitignored when Claude Code saves a setting to it |
> | Managed policy settings                  | Organization-wide                                                                                                | Yes, admin-controlled                                 |
> | [Plugin](/docs/en/plugins) `hooks/hooks.json` | When plugin is enabled                                                                                           | Yes, bundled with the plugin                          |
> | [Skill](/docs/en/skills) frontmatter          | The rest of the session once the skill is invoked. See [Hooks in skills and agents](#hooks-in-skills-and-agents) | Yes, defined in the skill file                        |
> | [Subagent](/docs/en/sub-agents) frontmatter   | While that subagent is running                                                                                   | Yes, defined in the subagent file                     |
>
> Cloud sessions on [Claude Code on the web](/docs/en/claude-code-on-the-web) don't read your local `~/.claude/settings.json`; hooks there come from the repo and from your organization's server-managed settings. In a [self-hosted environment](/docs/en/self-hosted-environments-configuration#permissions-and-tool-approval), Claude Code also runs the hooks the operator seeded from the runner host's `~/.claude/`, and it runs the hooks in the runner image's managed settings file when neither [server-managed settings nor an MDM-delivered Claude Code policy](/docs/en/settings#precedence-within-the-managed-tier) supplies the managed tier. See [what carries over from your setup](/docs/en/cloud-environments#what-carries-over-from-your-setup) for which files reach a cloud session.
>
> For details on settings file resolution, see [settings](/docs/en/settings).
>
> Hooks from settings files, managed policy settings, and plugins also run inside [subagents](/docs/en/sub-agents). When a subagent calls a tool, tool events such as `PreToolUse` and `PostToolUse` fire the same configured hooks as in the main conversation, and the input carries the `agent_id` and `agent_type` [common input fields](#common-input-fields) that identify the subagent.
>
> Enterprise administrators can use `allowManagedHooksOnly` to restrict which hooks run:
>
> * Your user, project, local, and plugin hooks are blocked. Hooks from plugins force-enabled in managed settings `enabledPlugins` are exempt
> * Claud### Source: Claude Code settings — Settings files
>
> ## Settings files and who they affect
>
> Claude Code reads settings from four files, and an organization can also deliver managed settings from the claude.ai console. Each source has a scope: the set of people and projects a setting saved in it applies to, whether that's just you, everyone in a project, or everyone in your organization.
>
> | Scope          | File                                                                                          | Who it affects                                                                                                                                                       | Use it for                                                                         |
> | :------------- | :-------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------- |
> | User           | `~/.claude/settings.json`                                                                     | You, in every project on this machine                                                                                                                                | Personal preferences: theme, editor mode, default model, your own permission rules |
> | Shared project | `.claude/settings.json`                                                                       | Everyone who starts Claude Code in the folder that contains it. In a git repository, commit it so teammates get it                                                   | Team permissions, hooks, plugins, and the environment variables the project needs  |
> | Project local  | `.claude/settings.local.json`                                                                 | You, in this one project only. Claude Code keeps it out of git when it creates the file; if you create it by hand, add it to `.gitignore` yourself                   | Personal overrides for one project, and testing before you share                   |
> | Managed        | `managed-settings.json` and other [managed sources](/docs/en/managed-settings#delivery-mechanisms) | Everyone your organization deploys it to; nothing you set overrides it, apart from a few [security-sensitive exceptions](#exceptions-to-managed-settings-precedence) | Security policy and compliance requirements                                        |
>
> In the File column, `~/.claude` is the `.claude` folder in your home directory, and a bare `.claude` is the `.claude` folder inside the project you start Claude Code in.
>
> ### Compare the scope of each settings file
>
> Suppose you have three projects on your machine, `website/`, `api/`, and `acme-app/`, a teammate has their own clone of `acme-app/`, and you start a [cloud session](#settings-in-cloud-sessions) on `acme-app/`.
>
> The graphic below shows which of those folders a setting applies in when you start Claude Code from them. Click a settings file to see the folders it reaches.
>
> <SettingsScope />
>
> * **`~/.claude/settings.json`**: every project on your machine, and nothing on your teammate's or in the cloud session
> * **`acme-app/.claude/settings.json`**: your `acme-app/`. It reaches your teammate's clone and the cloud session only if you commit the file to### Source: Claude Code settings — Settings precedence
>
> ## Settings precedence
>
> When the same key appears in more than one place, Claude Code uses the value from the highest level that sets it. The stack below shows the levels, highest on top; a key at a higher level overrides the same key anywhere below it.
>
> <SettingsPrecedence />
>
> In order, highest precedence first:
>
> 1. **Managed settings**: settings your organization deploys, by a `managed-settings.json` file, an MDM policy, or [server-managed settings](/docs/en/server-managed-settings) from the claude.ai console. Nothing you set overrides them: a key you pass with `--settings` doesn't override the same managed key, and a flag such as `--model` picks only from the models your organization allows. A managed `model` sets the model each session starts with, and you can still switch with `/model`; the lock is [`availableModels`](/docs/en/settings-reference#availablemodels), which constrains `/model`, `--model`, and the `model` key in your own files. When your organization delivers more than one managed source, the rules for [precedence within the managed tier](/docs/en/managed-settings#precedence-within-the-managed-tier) say what Claude Code reads from each.
> 2. **Command line arguments**: flags you pass when you start `claude` from a terminal, for one session; see [Change a setting for one session](#change-a-setting-for-one-session). Claude Code merges JSON you pass with `--settings <file-or-json>` with your settings files by the same rules as the other levels: it takes a key you set here over the same key in local, project, or user settings, and keeps the lower-level value for a key you omit.
> 3. **Project local settings** (`.claude/settings.local.json`): your personal settings for this project.
> 4. **Shared project settings** (`.claude/settings.json`): settings your team checks into source control.
> 5. **User settings** (`~/.claude/settings.json`): your personal settings for every project.
>
> Environment variables aren't a level in this stack. When a behavior has both a shell variable and a settings key, which one applies is decided per pair, not by level: `ANTHROPIC_MODEL` exported in your shell applies over the `model` key from any file, while `ANTHROPIC_DEFAULT_MODEL` applies only when no file sets `model`. The [environment variables reference](/docs/en/env-vars#precedence) says which keys have a pair and which one Claude Code reads first. An `env` block inside a settings file is an ordinary key and follows the levels above.
>
> For a few security-sensitive keys, Claude Code honors a stricter value from a lower level over a managed value; [Exceptions to managed settings precedence](#exceptions-to-managed-settings-precedence) lists them.
>
> ### Lists merge instead of overriding
>
> When you set the same list key, such as `permissions.allow`, in more than one file, Claude Code combines the lists instead of picking one, so each file can add entries without removing another file's. Three keys that hold model lists follow their own rules:
>
> * [`fallbackModel`](/docs/en/settings-reference#fallbackmodel) is an ordered chain where position carries meaning, so Claude Code takes the whole value from the highest-precedence file that defines it.
> * [`modelPicker`](/docs/en/settings-reference#modelpicker) holds one ordered list of rows plus a replace flag, so Claude Code never merges rows from two sources. It takes the whole value from the highest of managed settings, `--settings`, and user settings that defines it, and ignores the key in project and local settings### Source: MCP — Configuration locations
>
> mcpServers` JSON block**: configuration written for another client's settings file.
>
> Each is one of the inputs the four options in [Installing MCP servers](#installing-mcp-servers) take. Find the shape you have below to turn it into the command Claude Code accepts. Each command writes to [local scope](#local-scope) unless you add `--scope project` or `--scope user`.
>
> #### From a URL
>
> A URL means the server is remote. For an `https://` endpoint, add it with `--transport http`, or with `--transport sse` when the instructions say the endpoint uses SSE. For a `wss://` endpoint, use [Option 4](#option-4-add-a-remote-websocket-server) instead, since `--transport` doesn't accept `ws`:
>
> ```bash
> claude mcp add --transport http example https://mcp.example.com/mcp
> ```
>
> If the instructions also give an API key or token header, pass it with `--header` as shown in [Option 1](#option-1-add-a-remote-http-server).
>
> #### From an `npx`, `uvx`, or binary command
>
> A launch command means the server runs as a local stdio process. Put the whole command after `--`, so Claude Code passes flags such as `-y` to the command that starts the server instead of reading them as its own options. Pass any environment variables the instructions ask for with `--env`, after the server name and before `--`:
>
> ```bash
> claude mcp add example --env API_KEY=your-key -- npx -y @example/mcp-server
> ```
>
> [Option 3](#option-3-add-a-local-stdio-server) covers the `--` separator in full.
>
> #### From an `mcpServers` JSON block
>
> An `mcpServers` block written for another MCP client, such as Claude Desktop, uses the wrapper key and entry shape Claude Code reads. Pass `claude mcp add-json` the object inside `mcpServers`, not the wrapper. Two entries need a repair first:
>
> * **A `url` with no `type`**: add `"type": "http"`, `"type": "sse"`, or `"type": "ws"` to match the endpoint. Claude Code reads an entry with no `type` as a stdio server, so a `url` entry without a `type` fails.
> * **A key with cha### Source: Memory — CLAUDE.md locations
>
> ## Choose where to put CLAUDE.md files
>
> CLAUDE.md files can live in several locations, each with a different scope. The table below lists them in load order, from broadest scope to most specific, so a project instruction appears in context after a user instruction.
>
> | Scope                    | Location                                                                                                                                                                | Purpose                                                    | Use case examples                                                    | Shared with                     |
> | ------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------- | -------------------------------------------------------------------- | ------------------------------- |
> | **Managed policy**       | • macOS: `/Library/Application Support/ClaudeCode/CLAUDE.md`<br />• Linux and WSL: `/etc/claude-code/CLAUDE.md`<br />• Windows: `C:\Program Files\ClaudeCode\CLAUDE.md` | Organization-wide instructions managed by IT/DevOps        | Company coding standards, security policies, compliance requirements | All users in organization       |
> | **User instructions**    | `~/.claude/CLAUDE.md`                                                                                                                                                   | Personal preferences for all projects                      | Code styling preferences, personal tooling shortcuts                 | Just you (all projects)         |
> | **Project instructions** | `./CLAUDE.md` or `./.claude/CLAUDE.md`                                                                                                                                  | Team-shared instructions for the project                   | Project architecture, coding standards, common workflows             | Team members via source control |
> | **Local instructions**   | `./CLAUDE.local.md`                                                                                                                                                     | Personal project-specific preferences; add to `.gitignore` | Your sandbox URLs, preferred test data                               | Just you (current project)      |
>
> CLAUDE.md and CLAUDE.local.md files in the directory hierarchy above the working directory are loaded at launch. Files in subdirectories load on demand when Claude reads files in those directories. See [How CLAUDE.md files load](#how-claude-md-files-load) for the full resolution order.
>
> For large projects, you can break instructions into topic-specific files using [project rules](#organize-rules-with-claude/rules/). Rules let you scope instructions to specific file types or subdirectories.
>
> ### Set up a project CLAUDE.md
>
> A project CLAUDE.md can be stored in either `./CLAUDE.md` or `./.claude/CLAUDE.md`. Create this file and add instructions that apply to anyone working on the project: build and test commands, coding standards, architectural decisions, naming conventions, and common workflows. These instructions are shared with your team through version control, so focus on project-level standards rather than personal preferences. To confirm the file loaded, run `/context` in a session and check the list under **Memory files**### Source: Memory — Rules
>
> ## Organize rules with `.claude/rules/`
>
> For larger projects, you can organize instructions into multiple files using the `.claude/rules/` directory. This keeps instructions modular and easier for teams to maintain. Rules can also be [scoped to specific file paths](#path-specific-rules), so they only load into context when Claude works with matching files, reducing noise and saving context space.
>
> > [Note] Rules load into context every session or when matching files are opened. For task-specific instructions that don't need to be in context all the time, use [skills](/docs/en/skills) instead, which only load when you invoke them or when Claude determines they're relevant to your prompt.
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
> Clau### Source: Skills — Where skills live
>
> ## Where skills live
>
> Where you store a skill determines who can use it:
>
> | Location   | Path                                         | Applies to                     |
> | :--------- | :------------------------------------------- | :----------------------------- |
> | Enterprise | See [managed settings](/docs/en/managed-settings) | All users in your organization |
> | Personal   | `~/.claude/skills/<skill-name>/SKILL.md`     | All your projects              |
> | Project    | `.claude/skills/<skill-name>/SKILL.md`       | This project only              |
> | Plugin     | `<plugin>/skills/<skill-name>/SKILL.md`      | Where plugin is enabled        |
>
> When skills share the same name, Claude Code resolves the conflict by source:
>
> * Across levels, enterprise overrides personal, and personal overrides project.
>   * For example, with a `deploy` skill in both `~/.claude/skills/` and your project's `.claude/skills/`, `/deploy` runs the personal one.
> * A skill at any of these levels also overrides a bundled skill with the same name, but not the bundled skill's aliases.
>   * For example, a `code-review` skill in your project's `.claude/skills/` replaces the bundled `/code-review`, and typing the bundled alias `/review` never runs your skill.
> * Plugin skills use a `plugin-name:skill-name` namespace, so they can't conflict with other levels.
>   * For example, `my-plugin/skills/deploy/SKILL.md` becomes `/my-plugin:deploy` and loads alongside a `deploy` skill in your project's `.claude/skills/`.
> * If you have files in `.claude/commands/`, those work the same way, but if a skill and a command share the same name, the skill takes precedence.
>   * For example, with both `.claude/commands/deploy.md` and `.claude/skills/deploy/SKILL.md`, `/deploy` runs the skill.
> * A skill or command from any of these sources overrides a skill [synced from your claude.ai account](#when-a-synced-skill-name-matches-another-command) with the same name.
>   * For example, with a `deploy` skill enabled on claude.ai and another in your project's `.claude/skills/`, `/deploy` runs the project one.
>
> Skills also load from nested `.claude/skills/` directories below your working directory. When Claude reads or edits a file in a subdirectory, skills from that subdirectory's `.claude/skills/` become available. This lets a monorepo package provide its own skills that apply when working on that package, even if the session started at the repo root.
>
> If a nested skill shares a name with another skill, both stay available. For example, with a `deploy` skill at the project root and another in `apps/web/.claude/skills/`:
>
> * The nested one appears under a directory-qualified name, `apps/web:deploy`.
> * Its description says which directory it applies to.
> * Claude picks the variant that matches the files it is working on.
>
> Typing `/deploy` runs the project-root skill. Type the qualified name `/apps/web:deploy` to run the nested variant explicitly.
>
> When you or Claude invoke the unqualified name, the project-root skill loads, and Claude Code appends a list of the directory-qualified variants to its content with an instruction to also invoke any variant whose directory holds the files Claude is working on. A nested skill therefore still applies to work in its directory when only the unqualified name is invoked.
>
> The folder name `synced` is reserved in the enterprise, personal, and project skills locations, in any capitalization. Claude Code [downloads the skills you enable on claude.ai](/docs/en/env-vars#variab### Source: Sub-agents — Where subagents live
>
> .claude/agents/` directly. Subagent files, frontmatter fields, and the `.claude/agents/` and `~/.claude/agents/` locations are unchanged; only the terminal wizard is removed.
>
> This walkthrough creates a user-level subagent that reviews code and suggests improvements.
>
> 1. **Ask Claude to create the subagent**: In Claude Code, describe the subagent you want and where to save it:
>
>     ```text wrap
>     Create a personal code-improver subagent in ~/.claude/agents/ that scans
>     files and suggests improvements for readability, performance, and best
>     practices. It should explain each issue, show the current code, and
>     provide an improved version. Make it read-only and have it use Sonnet.
>     ```
>
>     Claude writes the file with a `name`, a `description`, a `tools` list, a `model`, and a system prompt.
>
> 2. **Review the file**: Open `~/.claude/agents/code-improver.md` and confirm the frontmatter matches what you asked for. The result looks like this:
>
>     ```markdown
>     ---
>     name: code-improver
>     description: Scans files and suggests improvements for readability, performance, and best practices. Use after writing or modifying code.
>     tools: Read, Grep, Glob
>     model: sonnet
>     ---
>
>     You are a code improvement specialist. For each issue you find, explain
>     the problem, show the current code, and provide an improved version.
>     ```
>
>     Because the file lives in `~/.claude/agents/`, the subagent is available in every project on your machine. To scope it to one project instead, move it to that project's `.claude/agents/` directory. [Choose the subagent scope](#choose-the-subagent-scope) compares the two.
>
> 3. **Try it out**: Ask Claude to delegate to the new subagent:
>
>     ```text wrap
>     Use the code-improver agent to suggest improvements in this project
>     ```
>
>     Claude delegates to your new subagent, which scans the codebase and returns improvement suggestions. In the transcript, the delegation appears as a tool call row showing the subagent's name followed by a short task description, such as `code-improver (Suggest code improvements)`.
>
>     If Claude can't find the new subagent, restart Claude Code and try again. This happens only when `~/.claude/agents/` didn't exist before the session started, because a running session doesn't detect a newly created `agents` directory.
>
> You now have a subagent you can use in any project on your machine to analyze codebases and suggest improvements.
>
> You can also write subagent files by hand, define them via CLI flags, or distribute them through plugins. The following sections cover all configuration options.
>
> > [Note] On Claude Code v2.1.197 and earlier, `/agents` opens an interactive wizard with a **Running** tab that lists live subagents and a **Library** tab for creating, editing, and deleting them.&#x20;
>
> ## Configure subagents
>
> A subagent's file location determines who it's available to, and its frontmatter determines what it can do. This section covers where subagent files live and every field they support.
>
> ### Choose the subagent scope
>
> Store subagent files in different locations depending on scope. When multiple subagents share the same name, Claude Code uses the one from the higher-priority location.
>
> | Location                     | Scope                   | Priority    | How to create                                 |
> | :-------------------