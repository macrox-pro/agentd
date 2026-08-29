---
primary_sources:
  - id: T1-SKILLS
    title: "Agent Skills"
    url: "https://geminicli.com/docs/cli/skills.md"
    section: "Full page"
  - id: T1-SKILLS-CREATE
    title: "Creating Agent Skills"
    url: "https://geminicli.com/docs/cli/creating-skills.md"
    section: "Full page"
  - id: T1-SKILLS-USE
    title: "Managing Agent Skills"
    url: "https://geminicli.com/docs/cli/using-agent-skills.md"
    section: "Full page"
  - id: T1-SKILLS-BP
    title: "Agent Skill best practices"
    url: "https://geminicli.com/docs/cli/skills-best-practices.md"
    section: "Full page"
  - id: T1-ACTIVATE-SKILL
    title: "Activate skill tool (`activate_skill`)"
    url: "https://geminicli.com/docs/tools/activate-skill.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Agent Skills

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Agent Skills — Full page

> Agent Skills let you extend Gemini CLI with specialized expertise, procedural
> workflows, and task-specific resources. Based on the
> [Agent Skills](https://agentskills.io) open standard, a "skill" is a
> self-contained directory that packages instructions and assets into a
> discoverable capability.
>
> Unlike general context files ([GEMINI.md](/docs/cli/gemini-md)), which provide
> persistent workspace-wide background, Skills represent **on-demand expertise**.
> This lets Gemini CLI maintain a vast library of specialized capabilities—such as
> security auditing, cloud deployments, or codebase migrations—without cluttering
> the model's immediate context window.
>
> ## How it works
>
> The lifecycle of an Agent Skill involves discovery, activation, and conditional
> resource access.
>
> 1.  **Discovery**: At the start of a session, Gemini CLI scans the discovery
>     tiers and injects the name and description of all enabled skills into the
>     system prompt.
> 2.  **Activation**: When Gemini identifies a task matching a skill's
>     description, it calls the `activate_skill` tool.
> 3.  **Consent**: You will see a confirmation prompt in the UI detailing the
>     skill's name, purpose, and the directory path it will gain access to.
> 4.  **Injection**: Upon your approval:
>     - The `SKILL.md` body and folder structure is added to the conversation
>       history.
>     - The skill's directory is added to the agent's allowed file paths, granting
>       it permission to read any bundled assets.
> 5.  **Execution**: The model proceeds with the specialized expertise active. It
>     is instructed to prioritize the skill's procedural guidance within reason.
>
> ## Discovery tiers
>
> Gemini CLI discovers skills from several locations, following a specific order
> of precedence (lowest to highest):
>
> 1.  **Built-in skills**: Standard skills included with Gemini CLI that provide
>     foundational capabilities.
> 2.  **Extension skills**: Skills bundled within installed
>     [extensions](/docs/extensions).
> 3.  **User skills**: Located in `~/.gemini/skills/` or the `~/.agents/skills/`
>     alias.
> 4.  **Workspace skills**: Located in `.gemini/skills/` or the `.agents/skills/`
>     alias. Workspace skills are shared with your team via version control.
>
> ### Precedence and aliases
>
> If multiple skills share the same name, the version from the higher-precedence
> location is used. Within the same tier (user or workspace), the
> `.agents/skills/` alias takes precedence over the `.gemini/skills/` directory.
>
> The `.agents/skills/` alias provides an interoperable path for managing
> agent-specific expertise that remains compatible across different AI tools.
>
> ## Key benefits
>
> Agent Skills provide several advantages for managing specialized knowledge and
> complex workflows.
>
> - **Shared expertise**: Package complex workflows (like a specific team's PR
>   review process) into a folder that anyone can use.
> - **Repeatable workflows**: Ensure complex multi-step tasks are performed
>   consistently by providing a procedural framework.
> - **Resource bundling**: Include scripts, templates, or example data alongside
>   instructions so the agent has everything it needs.
> - **Progressive disclosure**: Only skill metadata (name and description) is
>   loaded initially. Detailed instructions and resources are only disclosed when
>   the model explicitly activates the skill, saving context tokens.
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > `/skills disable` and `/skills enable` default to the `user` scope. Use
> > `--scope workspace` to manage workspace-specific settings.
>
> To see all available skills in your current session, use the `/skills list`
> command.
>
> ## Managing skills
>
> You can manage Agent Skills through interactive session commands or directly
> from your terminal.
>
> ### In an interactive session
>
> Use the `/skills` slash command to view and manage available expertise:
>
> - `/skills list [all] [nodesc]`: Shows discovered skills. Use `all` to include
>   built-in skills and `nodesc` to hide descriptions.
> - `/skills link <path> [--scope user|workspace]`: Links skills from a local
>   directory.
> - `/skills disable <name>`: Prevents a specific skill from being used.
> - `/skills enable <name>`: Re-enables a disabled skill.
> - `/skills reload` (or `/skills refresh`): Refreshes the list of discovered
>   skills from all tiers.
>
> ### From the terminal
>
> The `gemini skills` command provides management utilities:
>
> ```bash
> # List all discovered skills. Use --all to include built-in skills.
> gemini skills list --all
>
> # Install a skill from a Git repository or local directory.
> # Use --consent to skip the security confirmation prompt.
> gemini skills install https://github.com/user/repo.git --consent
>
> # Uninstall a skill.
> gemini skills uninstall my-skill --scope workspace
> ```
>
> #### Command options
>
> The skill management commands support several global and command-specific
> options.
>
> - `--scope`: Either `user` (global, default) or `workspace` (local to the
>   project).
> - `--path`: The sub-directory within a Git repository containing the skill.
> - `--consent`: Acknowledge security risks and skip the interactive confirmation
>   during installation.
>
> For more details on CLI commands, see the
> [CLI reference](/docs/cli/cli-reference#skills-management).
>
> ## Next steps
>
> Explore these resources to refine your skills and understand the framework
> better.
>
> - [Get started with Agent Skills](/docs/cli/tutorials/skills-getting-started): A
>   quick walkthrough of triggering and using skills.
> - [Creating Agent Skills](/docs/cli/creating-skills): Create your first skill and
>   bundle custom logic.
> - [Using Agent Skills](/docs/cli/using-agent-skills): Learn how to leverage built-in
>   and custom skills.
> - [Best practices](/docs/cli/skills-best-practices): Learn strategies for building
>   effective skills.

### Source: Creating Agent Skills — Full page

> Agent Skills let you extend Gemini CLI with specialized expertise, procedural
> workflows, and task-specific resources. This guide walks you through both
> automated and manual methods for creating and organizing your skills.
>
> ## Quickstart: Create a skill with a prompt
>
> The fastest way to create a new skill is to use the built-in `skill-creator`.
> This meta-skill guides you through designing, scaffolding, and validating your
> expertise.
>
> Simply ask Gemini CLI to create a skill for you:
>
> > "Create a new skill called 'code-reviewer' that analyzes local files for
> > common errors and style violations."
>
> Gemini will then:
>
> 1.  Generate a new directory for your skill (for example, `my-new-skill/`).
> 2.  Create a `SKILL.md` file with the necessary YAML frontmatter (`name` and
>     `description`).
> 3.  Create the standard resource directories: `scripts/`, `references/`, and
>     `assets/`.
>
> Once created, you can find your new skill in `.gemini/skills/code-reviewer/`.
>
> ## Manual creation
>
> 1.  **Create a directory** for your skill (for example, `my-new-skill/`).
> 2.  **Create a `SKILL.md` file** inside the new directory.
>
> ### 1. Create the directory structure
>
> The first step is to create the necessary folders for your skill and its
> scripts.
>
> **macOS/Linux**
>
> ```bash
> mkdir -p .gemini/skills/code-reviewer/scripts
> ```
>
> **Windows (PowerShell)**
>
> ```powershell
> New-Item -ItemType Directory -Force -Path ".gemini\skills\code-reviewer\scripts"
> ```
>
> ### 2. Define the skill (`SKILL.md`)
>
> The `SKILL.md` file defines the skill's purpose and instructions for the agent.
> Create a file at `.gemini/skills/code-reviewer/SKILL.md`.
>
> ```markdown
> ---
> name: code-reviewer
> description:
>   Expertise in reviewing code changes for correctness, security, and style. Use
>   when the user asks to "review" their code or a PR.
> ---
>
> # Code Reviewer Instructions
>
> You act as a senior software engineer specialized in code quality. When this
> skill is active, you MUST:
>
> 1.  **Analyze**: Review the provided code for logical errors, security
>     vulnerabilities, and style violations.
> 2.  **Review**: Use the bundled `scripts/review.js` utility to perform an
>     automated check.
> 3.  **Feedback**: Provide constructive feedback, clearly distinguishing between
>     critical issues and minor improvements.
> ```
>
> ### 3. Add the tool logic
>
> Skills can bundle resources like scripts to perform deterministic tasks. Create
> a file at `.gemini/skills/code-reviewer/scripts/review.js`.
>
> ```javascript
> // .gemini/skills/code-reviewer/scripts/review.js
> const file = process.argv[2];
>
> if (!file) {
>   console.error('Usage: node review.js <file>');
>   process.exit(1);
> }
>
> console.log(`Reviewing ${file}...`);
> // Simple mock review logic
> setTimeout(() => {
>   console.log(`Result: Success (No major issues found in ${file})`);
> }, 500);
> ```
>
> ### 4. Test the skill
>
> Gemini CLI automatically discovers skills in the `.gemini/skills` directory.
>
> 1.  Start a new session and ask a question that triggers the skill's
>     description: "Can you review index.js"
> 2.  Gemini identifies the request matches the `code-reviewer` description and
>     asks for permission to activate it.
> 3.  Once you approve, Gemini executes the bundled script:
>     `node .gemini/skills/code-reviewer/scripts/review.js index.js`
>
> To determine whether your skill has been correctly loaded, run the command:
>
> ```bash
> /skills
> ```
>
> ### 5. Optional: Share your skill
>
> You can share your skills in several ways depending on your target audience.
>
> - **Workspace skills**: Commit your skill to a `.gemini/skills/` directory in
>   your project repository.
> - **Extensions**: Bundle your skill within a
>   [Gemini CLI extension](/docs/extensions/writing-extensions).
> - **Git repositories**: Share the skill directory as a standalone Git repo and
>   install it using `gemini skills install <url>`.
>
> ---
>
> ## Core concepts
>
> Now that you've built your first skill, let's explore the core components and
> workflows for developing more complex expertise.
>
> ### Skill structure
>
> While a `SKILL.md` file is the only required component, we recommend the
> following structure for organizing your skill's resources.
>
> ```text
> my-skill/
> ├── SKILL.md       (Required) Instructions and metadata
> ├── scripts/       (Optional) Executable scripts
> ├── references/    (Optional) Static documentation
> └── assets/        (Optional) Templates and other resources
> ```
>
> When a skill is activated, the model is granted access to this entire directory.
> You can instruct the model to use the tools and files found within these
> folders.
>
> ### Metadata and triggers
>
> The `SKILL.md` file uses YAML frontmatter for metadata.
>
> - **`name`**: A unique identifier for the skill. This should match the directory
>   name.
> - **`description`**: **CRITICAL.** This is how Gemini decides when to use the
>   skill. Be specific about the tasks it handles and the keywords that should
>   trigger it.
>
> ### Discovery tiers
>
> Gemini CLI discovers skills from several locations, following a specific order
> of precedence (lowest to highest):
>
> 1.  **Built-in Skills**: Included with Gemini CLI (pre-approved).
> 2.  **Extension Skills**: Bundled within [extensions](/docs/extensions).
> 3.  **User Skills**: `~/.gemini/skills/` or the `~/.agents/skills/` alias.
> 4.  **Workspace Skills**: `.gemini/skills/` or the `.agents/skills/` alias.
>
> ### Discovery aliases
>
> You can use `.agents/skills` as an alternative to `.gemini/skills`. This alias
> is compatible with other AI agent tools following the
> [Agent Skills](https://agentskills.io) standard.
>
> ## Advanced development
>
> Once you've built a basic skill, you can use specialized scripts and workflows
> to streamline your development process.
>
> ### Creation scripts
>
> If you are developing a skill and want to use the same scripts the built-in
> tools use, you can find them in the core package. These scripts help automate
> the initialization, validation, and packaging of skills.
>
> - **Initialize**: `node scripts/init_skill.cjs <name> --path <dir>`
> - **Validate**: `node scripts/validate_skill.cjs <path/to/skill>`
> - **Package**: `node scripts/package_skill.cjs <path/to/skill>` (Creates a
>   `.skill` zip file)
>
> ### Linking for local development
>
> If you are developing a skill in a separate directory, you can link it to your
> user skills directory for testing:
>
> ```bash
> gemini skills link .
> ```
>
> ## Next steps
>
> - [Skill best practices](/docs/cli/skills-best-practices): Learn strategies for
>   building reliable and effective skills.
> - [Agent Skills overview](/docs/cli/skills): Deep dive into discovery tiers and the
>   skill lifecycle.
> - [Get started with Agent Skills](/docs/cli/tutorials/skills-getting-started): A
>   quick walkthrough of triggering and using skills.

### Source: Managing Agent Skills — Full page

> Agent Skills provide Gemini CLI with specialized expertise on demand. This guide
> covers advanced management techniques, including using slash commands, terminal
> utilities, and understanding discovery tiers.
>
> ## Discovery tiers
>
> Gemini CLI discovers skills from several locations, following a specific order
> of precedence (lowest to highest):
>
> 1.  **Built-in Skills**: Included with Gemini CLI and always available.
> 2.  **Extension Skills**: Bundled within [extensions](/docs/extensions).
> 3.  **User Skills**: Located in `~/.gemini/skills/` or the `~/.agents/skills/`
>     alias. These are available across all your projects.
> 4.  **Workspace Skills**: Located in `.gemini/skills/` or the `.agents/skills/`
>     alias within your current directory. These are project-specific.
>
> > **Tip:** If multiple skills share the same name, the version from the
> > higher-precedence location is used.
>
> ## In-session management
>
> Use the `/skills` slash command during an interactive session to manage your
> available expertise.
>
> - **List skills**: `/skills list` (shows discovered skills).
>   - Use `/skills list all` to include internal built-in skills.
>   - Use `/skills list nodesc` to hide descriptions.
> - **Reload skills**: `/skills reload` (or `/skills refresh`) to scan for new or
>   modified skills without restarting the CLI.
> - **Toggle status**:
>   - `/skills disable <name>`: Prevents a skill from being triggered.
>   - `/skills enable <name>`: Re-enables a disabled skill.
> - **Link local skills**: `/skills link <path> [--scope user|workspace]` to
>   immediately use a skill you are developing.
>
> ## Terminal utilities
>
> The `gemini skills` command provides management utilities directly from your
> system shell.
>
> ### Install a skill
>
> To install a skill from a remote repository or a local `.skill` package:
>
> ```bash
> gemini skills install https://github.com/user/my-awesome-skill
> ```
>
> By default, this installs to your **user profile**. Use `--scope workspace` to
> install it only for the current project.
>
> ### Link for development
>
> If you are developing a skill, use the `link` command to create a reference to
> your local directory:
>
> ```bash
> gemini skills link ./path/to/my-skill
> ```
>
> ### Uninstall a skill
>
> To completely remove an installed or linked skill:
>
> ```bash
> gemini skills uninstall <name>
> ```
>
> ## Security and consent
>
> Agent Skills can execute scripts and access your files. To protect your
> environment:
>
> 1.  **Installation consent**: When installing from a remote URL, you will be
>     asked to confirm the source.
> 2.  **Activation consent**: Every time a skill is triggered during a session,
>     the agent must ask for permission to activate it and gain access to its
>     resources.
>
> ## Next steps
>
> - [Get started with Agent Skills](/docs/cli/tutorials/skills-getting-started): A
>   walkthrough for creating your first skill.
> - [Creating Agent Skills](/docs/cli/creating-skills): Detailed guide on bundling
>   scripts and assets.
> - [Skill best practices](/docs/cli/skills-best-practices): Strategies for building
>   reliable expertise.

### Source: Agent Skill best practices — Full page

> Create high-quality, reliable Agent Skills by following these established design
> principles and patterns.
>
> ## Design for discovery
>
> The most important part of a skill is its `description`. This is the only
> information the model has before activation.
>
> - **Be specific**: Use keywords that are likely to appear in user prompts (for
>   example, "audit," "security," "refactor," "migration").
> - **Define the trigger**: Clearly state _when_ the skill should be used (for
>   example, "Use this skill when the user asks to review a PR for performance
>   regressions").
> - **Avoid overlap**: Ensure your skill descriptions are distinct from one
>   another and from the general capabilities of the model.
>
> ## Progressive disclosure
>
> The "context window" is a shared resource. Use a three-level loading system to
> manage context efficiently.
>
> 1.  **Metadata (name + description)**: Always in context (~100 words).
> 2.  **`SKILL.md` body**: Loaded only after the skill triggers (<5k words).
> 3.  **Bundled resources**: Loaded only as needed by the model.
>
> **Best practice**: Keep the `SKILL.md` body focused on core procedural
> instructions. Move detailed reference material, schemas, and examples into
> separate files in a `references/` directory.
>
> ## Degrees of freedom
>
> Match the level of instruction specificity to the task's fragility.
>
> - **High freedom (text-based instructions)**: Use when multiple approaches are
>   valid or decisions depend heavily on context.
> - **Medium freedom (pseudocode or scripts with parameters)**: Use when a
>   preferred pattern exists but some variation is acceptable.
> - **Low freedom (specific scripts, few parameters)**: Use when operations are
>   fragile and error-prone, or a specific sequence MUST be followed.
>
> ## Bundle resources effectively
>
> Leverage the skill's ability to include scripts and assets to extend the agent's
> capabilities.
>
> - **Use scripts for deterministic tasks**: If a task can be automated with a
>   script (for example, running a linter, fetching data from an API), bundle it
>   in the `scripts/` folder.
> - **Agentic ergonomics**: Ensure scripts output LLM-friendly stdout. Suppress
>   verbose tracebacks and provide clear, concise success/failure messages.
> - **Provide templates**: Include common file headers or boilerplate code in the
>   `assets/` folder to ensure the agent produces consistent output.
>
> ## Anatomy of a great skill
>
> A well-structured skill directory organizes its resources into specialized
> sub-folders.
>
> ```text
> my-skill/
> ├── SKILL.md       (Required) Core instructions and metadata
> ├── scripts/       (Optional) Executable logic (Node.js, Python, etc.)
> ├── references/    (Optional) Documentation to be loaded as needed
> └── assets/        (Optional) Templates and non-executable resources
> ```
>
> ## Security and privacy
>
> Design your skills with security in mind to protect your workspace and data.
>
> - **Avoid hardcoded secrets**: Never include API keys or passwords in your
>   skill's scripts or documentation.
> - **Review third-party skills**: Inspect the `SKILL.md` and scripts of any skill
>   before installing it from an untrusted source.
> - **Limit scope**: Design skills to be as focused as possible to minimize the
>   potential impact of errors.

### Source: Activate skill tool (`activate_skill`) — Full page

> The `activate_skill` tool lets Gemini CLI load specialized procedural expertise
> and resources when they are relevant to your request.
>
> ## Description
>
> Skills are packages of instructions and tools designed for specific engineering
> tasks, such as reviewing code or creating pull requests. Gemini CLI uses this
> tool to "activate" a skill, which provides it with detailed guidelines and
> specialized tools tailored to that task.
>
> ### Arguments
>
> `activate_skill` takes one argument:
>
> - `name` (enum, required): The name of the skill to activate (for example,
>   `code-reviewer`, `pr-creator`, or `docs-writer`).
>
> ## Usage
>
> The `activate_skill` tool is used exclusively by the Gemini agent. You cannot
> invoke this tool manually.
>
> When the agent identifies that a task matches a discovered skill, it requests to
> activate that skill. Once activated, the agent's behavior is guided by the
> skill's specific instructions until the task is complete.
>
> ## Behavior
>
> The agent uses this tool to provide professional-grade assistance:
>
> - **Specialized logic:** Skills contain expert-level procedures for complex
>   workflows.
> - **Dynamic capability:** Activating a skill can grant the agent access to new,
>   task-specific tools.
> - **Contextual awareness:** Skills help the agent focus on the most relevant
>   standards and conventions for a particular task.
>
> ## Next steps
>
> - Learn how to [Use Agent Skills](/docs/cli/skills).
> - See the [Build agent skills](/docs/cli/creating-skills) guide.
