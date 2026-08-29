---
primary_sources:
  - id: T1-GET-STARTED
    title: "Get started with Gemini CLI"
    url: "https://geminicli.com/docs/get-started.md"
    section: "Full page"
  - id: T1-PLAN
    title: "Plan Mode"
    url: "https://geminicli.com/docs/cli/plan-mode.md"
    section: "Full page"
  - id: T1-MODEL-ROUTE
    title: "Model routing"
    url: "https://geminicli.com/docs/cli/model-routing.md"
    section: "Full page"
  - id: T1-MODEL
    title: "Gemini CLI model selection (`/model` command)"
    url: "https://geminicli.com/docs/cli/model.md"
    section: "Full page"
  - id: T1-MODEL-STEER
    title: "Model steering (experimental)"
    url: "https://geminicli.com/docs/cli/model-steering.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Overview and modes

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Get started with Gemini CLI — Full page

> Welcome to Gemini CLI! This guide will help you install, configure, and start
> using Gemini CLI to enhance your workflow right from your terminal.
>
> ## Quickstart: Install, authenticate, configure, and use Gemini CLI
>
> Gemini CLI brings the power of advanced language models directly to your command
> line interface. As an AI-powered assistant, Gemini CLI can help you with a
> variety of tasks, from understanding and generating code to reviewing and
> editing documents.
>
> ## Install
>
> The standard method to install and run Gemini CLI uses `npm`:
>
> ```bash
> npm install -g @google/gemini-cli
> ```
>
> Once Gemini CLI is installed, run Gemini CLI from your command line:
>
> ```bash
> gemini
> ```
>
> For more installation options, see
> [Gemini CLI Installation](/docs/get-started/installation).
>
> ## Authenticate
>
> To begin using Gemini CLI, you must authenticate with a Google service. In most
> cases, you can log in with your existing Google account:
>
> 1. Run Gemini CLI after installation:
>
>    ```bash
>    gemini
>    ```
>
> 2. When asked "How would you like to authenticate for this project?" select **1.
>    Sign in with Google**.
>
> 3. Select your Google account.
>
> 4. Click on **Sign in**.
>
> Certain account types may require you to configure a Google Cloud project. For
> more information, including other authentication methods, see
> [Gemini CLI Authentication Setup](/docs/get-started/authentication).
>
> ## Configure
>
> Gemini CLI offers several ways to configure its behavior, including environment
> variables, command-line arguments, and settings files.
>
> To explore your configuration options, see
> [Gemini CLI Configuration](/docs/reference/configuration).
>
> ## Use
>
> Once installed and authenticated, you can start using Gemini CLI by issuing
> commands and prompts in your terminal. Ask it to generate code, explain files,
> and more.
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > These examples demonstrate potential capabilities. Your actual
> > results can vary based on the model used and your project environment.
>
> ### Rename your photographs based on content
>
> You can use Gemini CLI to automate file management tasks that require visual
> analysis. In this example, Gemini CLI renames images based on their actual
> subject matter.
>
> Scenario: You have a folder containing the following files:
>
> ```bash
> photos/photo1.png
> photos/photo2.png
> photos/photo3.png
> ```
>
> Give Gemini the following prompt:
>
> ```cli
> Rename the photos in my "photos" directory based on their contents.
> ```
>
> Result: Gemini asks for permission to rename your files.
>
> Select **Allow once** and your files are renamed:
>
> ```bash
> photos/yellow_flowers.png
> photos/antique_dresser.png
> photos/green_android_robot.png
> ```
>
> ### Explain a repository by reading its code
>
> Gemini CLI is effective for rapid codebase exploration. The following example
> shows how to ask Gemini CLI to fetch, analyze, and summarize a remote project.
>
> Scenario: You want to understand how a popular open-source utility works by
> inspecting its code, not just its README.
>
> Give Gemini CLI the following prompt:
>
> ```cli
> Clone the 'chalk' repository from https://github.com/chalk/chalk, read its key source files, and explain how it works.
> ```
>
> Result: Gemini performs a sequence of actions to answer your request.
>
> 1.  First, it asks for permission to run `git clone` to download the repository.
> 2.  Next, it finds the important source files and asks for permission to read
>     them.
> 3.  Finally, after analyzing the code, it provides a summary.
>
> Gemini CLI returns an explanation based on the actual source code:
>
> ```markdown
> The `chalk` library is a popular npm package for styling terminal output with
> colors. After analyzing the source code, here's how it works:
>
> - **Core functionality:** The main file sets up a chainable API. Each color or
>   modifier (like `bold` or `italic`) is a getter that appends the corresponding
>   ANSI escape code to an internal stack.
>
> - **Chainable API:** When you call `chalk.red.bold('Hello')`, you are accessing
>   getters. The `red` getter adds the red color code, and the `bold` getter adds
>   the bold code.
>
> - **Output generation:** When the chain is treated as a string (for example, in
>   `console.log`), a final `toString()` method is called. This method joins all
>   the stored ANSI codes, wraps them around the input string ('Hello'), and adds
>   a reset code at the end. This produces the final, styled string that the
>   terminal can render.
> ```
>
> ### Combine two spreadsheets into one spreadsheet
>
> Gemini CLI can process and transform data across multiple files. Use this
> capability to merge reports or reformat data sets without manual copying.
>
> Scenario: You have two .csv files: `Revenue - 2023.csv` and
> `Revenue - 2024.csv`. Each file contains monthly revenue figures.
>
> Give Gemini CLI the following prompt:
>
> ```cli
> Combine the two .csv files into a single .csv file, with each year a different column.
> ```
>
> Result: Gemini CLI reads each file and then asks for permission to write a new
> file. Provide your permission and Gemini CLI provides the combined data:
>
> ```csv
> Month,2023,2024
> January,0,1000
> February,0,1200
> March,0,2400
> April,900,500
> May,1000,800
> June,1000,900
> July,1200,1000
> August,1800,400
> September,2000,2000
> October,2400,3400
> November,3400,1800
> December,2100,9000
> ```
>
> ### Run unit tests
>
> Gemini CLI can generate boilerplate code and tests based on your existing
> implementation. This example demonstrates how to request code coverage for a
> JavaScript component.
>
> Scenario: You've written a simple login page. You wish to write unit tests to
> ensure that your login page has code coverage.
>
> Give Gemini CLI the following prompt:
>
> ```cli
> Write unit tests for Login.js.
> ```
>
> Result: Gemini CLI asks for permission to write a new file and creates a test
> for your login page.
>
> ## Check usage and quota
>
> You can check your current token usage and quota information using the
> `/stats model` command. This command provides a snapshot of your current
> session's token usage, as well as your overall quota and usage for the supported
> models.
>
> For more information on the `/stats` command and its subcommands, see the
> [Command Reference](/docs/reference/commands#stats).
>
> ## Next steps
>
> - Follow the [File management](/docs/cli/tutorials/file-management) guide to
>   start working with your codebase.
> - See [Shell commands](/docs/cli/tutorials/shell-commands) to learn about
>   terminal integration.

### Source: Plan Mode — Full page

> Plan Mode is a read-only environment for architecting robust solutions before
> implementation. With Plan Mode, you can:
>
> - **Research:** Explore the project in a read-only state to prevent accidental
>   changes.
> - **Design:** Understand problems, evaluate trade-offs, and choose a solution.
> - **Plan:** Align on an execution strategy before any code is modified.
>
> Plan Mode is enabled by default. You can manage this setting using the
> `/settings` command.
>
> ## How to enter Plan Mode
>
> Plan Mode integrates seamlessly into your workflow, letting you switch between
> planning and execution as needed.
>
> You can either configure Gemini CLI to start in Plan Mode by default or enter
> Plan Mode manually during a session.
>
> ### Launch in Plan Mode
>
> To start Gemini CLI directly in Plan Mode by default:
>
> 1.  Use the `/settings` command.
> 2.  Set **Default Approval Mode** to `Plan`.
>
> To launch Gemini CLI in Plan Mode once:
>
> 1. Use `gemini --approval-mode=plan` when launching Gemini CLI.
>
> ### Enter Plan Mode manually
>
> To start Plan Mode while using Gemini CLI:
>
> - **Keyboard shortcut:** Press `Shift+Tab` to cycle through approval modes
>   (`Default` -> `Auto-Edit` -> `Plan`). Plan Mode is automatically removed from
>   the rotation when Gemini CLI is actively processing or showing confirmation
>   dialogs.
>
> - **Command:** Type `/plan [goal]` in the input box. The `[goal]` is optional;
>   for example, `/plan implement authentication` will switch to Plan Mode and
>   immediately submit the prompt to the model.
>
> - **Natural Language:** Ask Gemini CLI to "start a plan for...". Gemini CLI
>   calls the
>   [`enter_plan_mode`](/docs/tools/planning#1-enter_plan_mode-enterplanmode) tool
>   to switch modes. This tool is not available when Gemini CLI is in
>   [YOLO mode](/docs/reference/configuration#command-line-arguments).
>
> ## How to use Plan Mode
>
> Plan Mode lets you collaborate with Gemini CLI to design a solution before
> Gemini CLI takes action.
>
> 1.  **Provide a goal:** Start by describing what you want to achieve. Gemini CLI
>     will then enter Plan Mode (if it's not already) to research the task.
> 2.  **Discuss and agree on strategy:** As Gemini CLI analyzes your codebase, it
>     will discuss its findings and proposed strategy with you to ensure
>     alignment. It may ask you questions or present different implementation
>     options using [`ask_user`](/docs/tools/ask-user). **Gemini CLI will stop and
>     wait for your confirmation** before drafting the formal plan. You should
>     reach an informal agreement on the approach before proceeding.
> 3.  **Review the plan:** Once you've agreed on the strategy, Gemini CLI creates
>     a detailed implementation plan as a Markdown file in your plans directory.
>
>     - **View:** You can open and read this file to understand the proposed
>       changes.
>     - **Edit:** Press `Ctrl+X` to open the plan directly in your configured
>       external editor.
>
> 4.  **Approve or iterate:** Gemini CLI will present the finalized plan for your
>     formal approval.
>     - **Approve:** If you're satisfied with the plan, approve it to start the
>       implementation immediately: **Yes, automatically accept edits** or **Yes,
>       manually accept edits**.
>     - **Iterate:** If the plan needs adjustments, provide feedback in the input
>       box or [edit the plan file directly](#collaborative-plan-editing). Gemini
>       CLI will refine the strategy and update the plan.
>     - **Cancel:** You can cancel your plan with `Esc`.
>
> For more complex or specialized planning tasks, you can
> [customize the planning workflow with skills](#custom-planning-with-skills).
>
> ### Collaborative plan editing
>
> You can collaborate with Gemini CLI by making direct changes or leaving comments
> in the implementation plan. This is often faster and more precise than
> describing complex changes in natural language.
>
> 1.  **Open the plan:** Press `Ctrl+X` when Gemini CLI presents a plan for
>     review.
> 2.  **Edit or comment:** The plan opens in your configured external editor (for
>     example, VS Code or Vim). You can:
>     - **Modify steps:** Directly reorder, delete, or rewrite implementation
>       steps.
>     - **Leave comments:** Add inline questions or feedback (for example, "Wait,
>       shouldn't we use the existing `Logger` class here?").
> 3.  **Save and close:** Save your changes and close the editor.
> 4.  **Review and refine:** Gemini CLI automatically detects the changes, reviews
>     your comments, and adjusts the implementation strategy. It then presents the
>     refined plan for your final approval.
>
> ## How to exit Plan Mode
>
> You can exit Plan Mode at any time, whether you have finalized a plan or want to
> switch back to another mode.
>
> - **Approve a plan:** When Gemini CLI presents a finalized plan, approving it
>   automatically exits Plan Mode and starts the implementation.
> - **Keyboard shortcut:** Press `Shift+Tab` to cycle to the desired mode.
> - **Natural language:** Ask Gemini CLI to "exit plan mode" or "stop planning."
>
> ## Tool Restrictions
>
> Plan Mode enforces strict safety policies to prevent accidental changes.
>
> These are the only allowed tools:
>
> - **FileSystem (Read):**
>   [`read_file`](/docs/tools/file-system#2-read_file-readfile),
>   [`list_directory`](/docs/tools/file-system#1-list_directory-readfolder),
>   [`glob`](/docs/tools/file-system#4-glob-findfiles)
> - **Search:** [`grep_search`](/docs/tools/file-system#5-grep_search-searchtext),
>   [`google_web_search`](/docs/tools/web-search),
>   [`web_fetch`](/docs/tools/web-fetch) (requires explicit confirmation),
>   [`get_internal_docs`](/docs/tools/internal-docs)
> - **Research Subagents:**
>   [`codebase_investigator`](/docs/core/subagents#codebase-investigator),
>   [`cli_help`](/docs/core/subagents#cli-help-agent)
> - **Interaction:** [`ask_user`](/docs/tools/ask-user)
> - **MCP tools (Read):** Read-only [MCP tools](/docs/tools/mcp-server) (for
>   example, `github_read_issue`, `postgres_read_schema`) and core
>   [MCP resource tools](/docs/tools/mcp-resources) (`list_mcp_resources`,
>   `read_mcp_resource`) are allowed.
> - **Planning (Write):**
>   [`write_file`](/docs/tools/file-system#3-write_file-writefile) and
>   [`replace`](/docs/tools/file-system#6-replace-edit) only allowed for `.md`
>   files in the `~/.gemini/tmp/<project>/<session-id>/plans/` directory or your
>   [custom plans directory](#custom-plan-directory-and-policies).
> - **Skills:** [`activate_skill`](/docs/cli/skills) (allows loading specialized
>   instructions and resources in a read-only manner)
>
> ## Customization and best practices
>
> Plan Mode is secure by default, but you can adapt it to fit your specific
> workflows. You can customize how Gemini CLI plans by using skills, adjusting
> safety policies, changing where plans are stored, or adding hooks.
>
> ### Custom planning with skills
>
> You can use [Agent Skills](/docs/cli/skills) to customize how Gemini CLI
> approaches planning for specific types of tasks. When a skill is activated
> during Plan Mode, its specialized instructions and procedural workflows will
> guide the research, design, and planning phases.
>
> For example:
>
> - A **"Database Migration"** skill could ensure the plan includes data safety
>   checks and rollback strategies.
> - A **"Security Audit"** skill could prompt Gemini CLI to look for specific
>   vulnerabilities during codebase exploration.
> - A **"Frontend Design"** skill could guide Gemini CLI to use specific UI
>   components and accessibility standards in its proposal.
>
> To use a skill in Plan Mode, you can explicitly ask Gemini CLI to "use the
> `<skill-name>` skill to plan..." or Gemini CLI may autonomously activate it
> based on the task description.
>
> ### Custom policies
>
> Plan Mode's default tool restrictions are managed by the
> [policy engine](/docs/reference/policy-engine) and defined in the built-in
> [`plan.toml`] file. The built-in policy (Tier 1) enforces the read-only state,
> but you can customize these rules by creating your own policies in your
> `~/.gemini/policies/` directory (Tier 2).
>
> #### Global vs. mode-specific rules
>
> As described in the
> [policy engine documentation](/docs/reference/policy-engine#approval-modes), any
> rule that does not explicitly specify `modes` is considered "always active" and
> will apply to Plan Mode as well.
>
> To maintain the integrity of Plan Mode as a safe research environment,
> persistent tool approvals are context-aware. Approvals granted in modes like
> Default or Auto-Edit do not apply to Plan Mode, ensuring that tools trusted for
> implementation don't automatically execute while you're researching. However,
> approvals granted while in Plan Mode are treated as intentional choices for
> global trust and apply to all modes.
>
> If you want to manually restrict a rule to other modes but _not_ to Plan Mode,
> you must explicitly specify the target modes. For example, to allow `npm test`
> in default and Auto-Edit modes but not in Plan Mode:
>
> ```toml
> [[rule]]
> toolName = "run_shell_command"
> commandPrefix = "npm test"
> decision = "allow"
> priority = 100
> # By omitting "plan", this rule will not be active in Plan Mode.
> modes = ["default", "autoEdit"]
> ```
>
> #### Example: Automatically approve read-only MCP tools
>
> By default, read-only MCP tools require user confirmation in Plan Mode. You can
> use `toolAnnotations` and the `mcpName` wildcard to customize this behavior for
> your specific environment.
>
> `~/.gemini/policies/mcp-read-only.toml`
>
> ```toml
> [[rule]]
> toolName = "*"
> mcpName = "*"
> toolAnnotations = { readOnlyHint = true }
> decision = "allow"
> priority = 100
> modes = ["plan"]
> ```
>
> For more information on how the policy engine works, see the
> [policy engine](/docs/reference/policy-engine) docs.
>
> #### Example: Allow git commands in Plan Mode
>
> This rule lets you check the repository status and see changes while in Plan
> Mode.
>
> `~/.gemini/policies/git-research.toml`
>
> ```toml
> [[rule]]
> toolName = "run_shell_command"
> commandPrefix = ["git status", "git diff"]
> decision = "allow"
> priority = 100
> modes = ["plan"]
> ```
>
> #### Example: Enable custom subagents in Plan Mode
>
> Built-in research [subagents](/docs/core/subagents) like
> [`codebase_investigator`](/docs/core/subagents#codebase-investigator) and
> [`cli_help`](/docs/core/subagents#cli-help-agent) are enabled by default in Plan
> Mode. You can enable additional
> [custom subagents](/docs/core/subagents#creating-custom-subagents) by adding a
> rule to your policy.
>
> `~/.gemini/policies/research-subagents.toml`
>
> ```toml
> [[rule]]
> toolName = "my_custom_subagent"
> decision = "allow"
> priority = 100
> modes = ["plan"]
> ```
>
> Tell Gemini CLI it can use these tools in your prompt, for example: _"You can
> check ongoing changes in git."_
>
> ### Custom plan directory and policies
>
> By default, planning artifacts are stored in a managed temporary directory
> outside your project: `~/.gemini/tmp/<project>/<session-id>/plans/`.
>
> You can configure a custom directory for plans in your `settings.json`. For
> example, to store plans in a `.gemini/plans` directory within your project:
>
> ```json
> {
>   "general": {
>     "plan": {
>       "directory": ".gemini/plans"
>     }
>   }
> }
> ```
>
> To maintain the safety of Plan Mode, user-configured paths for the plans
> directory are restricted to the project root. This ensures that custom planning
> locations defined within a project's workspace cannot be used to escape and
> overwrite sensitive files elsewhere. Any user-configured directory must reside
> within the project boundary.
>
> Using a custom directory requires updating your
> [policy engine](/docs/reference/policy-engine) configurations to allow
> `write_file` and `replace` in that specific location. For example, to allow
> writing to the `.gemini/plans` directory within your project, create a policy
> file at `~/.gemini/policies/plan-custom-directory.toml`:
>
> ```toml
> [[rule]]
> toolName = ["write_file", "replace"]
> decision = "allow"
> priority = 100
> modes = ["plan"]
> # Adjust the pattern to match your custom directory.
> # This example matches any .md file in a .gemini/plans directory within the project.
> argsPattern = "\"file_path\":\"[^\"]+[\\\\/]+\\.gemini[\\\\/]+plans[\\\\/]+[\\w-]+\\.md\""
> ```
>
> ### Using hooks with Plan Mode
>
> You can use the [hook system](/docs/hooks/writing-hooks) to automate parts of
> the planning workflow or enforce additional checks when Gemini CLI transitions
> into or out of Plan Mode.
>
> Hooks such as `BeforeTool` or `AfterTool` can be configured to intercept the
> `enter_plan_mode` and `exit_plan_mode` tool calls.
>
> > [!WARNING] When hooks are triggered by **tool executions**, they do **not**
> > run when you manually toggle Plan Mode using the `/plan` command or the
> > `Shift+Tab` keyboard shortcut. If you need hooks to execute on mode changes,
> > ensure the transition is initiated by the agent (for example, by asking "start
> > a plan for...").
>
> #### Example: Archive approved plans to GCS (`AfterTool`)
>
> If your organizational policy requires a record of all execution plans, you can
> use an `AfterTool` hook to securely copy the plan artifact to Google Cloud
> Storage whenever Gemini CLI exits Plan Mode to start the implementation.
>
> **`.gemini/hooks/archive-plan.sh`:**
>
> ```bash
> #!/usr/bin/env bash
> # Extract the plan filename from the tool input JSON
> plan_filename=$(jq -r '.tool_input.plan_filename // empty')
>
> # Construct the absolute path using the GEMINI_PLANS_DIR environment variable
> plan_path="$GEMINI_PLANS_DIR/$plan_filename"
>
> if [ -f "$plan_path" ]; then
>   # Generate a unique filename using a timestamp
>   filename="$(date +%s)_$(basename "$plan_path")"
>
>   # Upload the plan to GCS in the background so it doesn't block the CLI
>   gsutil cp "$plan_path" "gs://my-audit-bucket/gemini-plans/$filename" > /dev/null 2>&1 &
> fi
>
> # AfterTool hooks should generally allow the flow to continue
> echo '{"decision": "allow"}'
> ```
>
> To register this `AfterTool` hook, add it to your `settings.json`:
>
> ```json
> {
>   "hooks": {
>     "AfterTool": [
>       {
>         "matcher": "exit_plan_mode",
>         "hooks": [
>           {
>             "name": "archive-plan",
>             "type": "command",
>             "command": "~/.gemini/hooks/archive-plan.sh"
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> ## Commands
>
> - **`/plan copy`**: Copy the currently approved plan to your clipboard.
>
> ## Planning workflows
>
> Plan Mode provides building blocks for structured research and design. These are
> implemented as [extensions](/docs/extensions) using core planning tools
> like [`enter_plan_mode`](/docs/tools/planning#1-enter_plan_mode-enterplanmode),
> [`exit_plan_mode`](/docs/tools/planning#2-exit_plan_mode-exitplanmode), and
> [`ask_user`](/docs/tools/ask-user).
>
> ### Built-in planning workflow
>
> The built-in planner uses an adaptive workflow to analyze your project, consult
> you on trade-offs via [`ask_user`](/docs/tools/ask-user), and draft a plan for
> your approval.
>
> ### Custom planning workflows
>
> You can install or create specialized planners to suit your workflow.
>
> #### Conductor
>
> [Conductor] is designed for spec-driven development. It organizes work into
> "tracks" and stores persistent artifacts in your project's `conductor/`
> directory:
>
> - **Automate transitions:** Switches to read-only mode via
>   [`enter_plan_mode`](/docs/tools/planning#1-enter_plan_mode-enterplanmode).
> - **Streamline decisions:** Uses [`ask_user`](/docs/tools/ask-user) for
>   architectural choices.
> - **Maintain project context:** Stores artifacts in the project directory using
>   [custom plan directory and policies](#custom-plan-directory-and-policies).
> - **Handoff execution:** Transitions to implementation via
>   [`exit_plan_mode`](/docs/tools/planning#2-exit_plan_mode-exitplanmode).
>
> #### Build your own
>
> Since Plan Mode is built on modular building blocks, you can develop your own
> custom planning workflow as an [extensions](/docs/extensions). By
> leveraging core tools and [custom policies](#custom-policies), you can define
> how Gemini CLI researches and stores plans for your specific domain.
>
> To build a custom planning workflow, you can use:
>
> - **Tool usage:** Use core tools like
>   [`enter_plan_mode`](/docs/tools/planning#1-enter_plan_mode-enterplanmode),
>   [`ask_user`](/docs/tools/ask-user), and
>   [`exit_plan_mode`](/docs/tools/planning#2-exit_plan_mode-exitplanmode) to
>   manage the research and design process.
> - **Customization:** Set your own storage locations and policy rules using
>   [custom plan directories](#custom-plan-directory-and-policies) and
>   [custom policies](#custom-policies).
>
> <!-- prettier-ignore -->
> > [!TIP]
> > Use [Conductor] as a reference when building your own custom
> > planning workflow.
>
> By using Plan Mode as its execution environment, your custom methodology can
> enforce read-only safety during the design phase while benefiting from
> high-reasoning model routing.
>
> ## Automatic Model Routing
>
> When using an [auto model](/docs/reference/configuration#model), Gemini CLI
> automatically optimizes [model routing](/docs/cli/telemetry#model-routing) based
> on the current phase of your task:
>
> 1.  **Planning Phase:** While in Plan Mode, the CLI routes requests to a
>     high-reasoning **Pro** model to ensure robust architectural decisions and
>     high-quality plans.
> 2.  **Implementation Phase:** Once a plan is approved and you exit Plan Mode,
>     the CLI detects the existence of the approved plan and automatically
>     switches to a high-speed **Flash** model. This provides a faster, more
>     responsive experience during the implementation of the plan.
>
> If the high-reasoning model is unavailable or you don't have access to it,
> Gemini CLI automatically and silently falls back to a faster model to ensure
> your workflow isn't interrupted.
>
> This behavior is enabled by default to provide the best balance of quality and
> performance. You can disable this automatic switching in your settings:
>
> ```json
> {
>   "general": {
>     "plan": {
>       "modelRouting": false
>     }
>   }
> }
> ```
>
> ## Cleanup
>
> By default, Gemini CLI automatically cleans up old session data, including all
> associated plan files and task trackers.
>
> - **Default behavior:** Sessions (and their plans) are retained for **30 days**.
> - **Configuration:** You can customize this behavior via the `/settings` command
>   (search for **Enable Session Cleanup** or **Keep chat history**) or in your
>   `settings.json` file. See
>   [session retention](/docs/cli/session-management#session-retention) for more
>   details.
>
> Manual deletion also removes all associated artifacts:
>
> - **Command Line:** Use `gemini --delete-session <index|id>`.
> - **Session Browser:** Press `/resume`, navigate to a session, and press `x`.
>
> If you use a [custom plans directory](#custom-plan-directory-and-policies),
> those files are not automatically deleted and must be managed manually.
>
> ## Non-interactive execution
>
> When running Gemini CLI in non-interactive environments (such as headless
> scripts or CI/CD pipelines), Plan Mode optimizes for automated workflows:
>
> - **Automatic transitions:** The policy engine automatically approves the
>   `enter_plan_mode` and `exit_plan_mode` tools without prompting for user
>   confirmation.
> - **Automated implementation:** When exiting Plan Mode to execute the plan,
>   Gemini CLI automatically switches to
>   [YOLO mode](/docs/reference/policy-engine#approval-modes) instead of the
>   standard Default mode. This allows the CLI to execute the implementation steps
>   automatically without hanging on interactive tool approvals.
>
> **Example:**
>
> ```bash
> gemini --approval-mode plan -p "Analyze telemetry and suggest improvements"
> ```
>
> [`plan.toml`]:
>   https://github.com/google-gemini/gemini-cli/blob/main/packages/core/src/policy/policies/plan.toml
> [Conductor]: https://github.com/gemini-cli-extensions/conductor
> [open an issue]: https://github.com/google-gemini/gemini-cli/issues

### Source: Model routing — Full page

> Gemini CLI includes a model routing feature that automatically switches to a
> fallback model in case of a model failure. This feature is enabled by default
> and provides resilience when the primary model is unavailable.
>
> ## How it works
>
> Model routing is managed by the `ModelAvailabilityService`, which monitors model
> health and automatically routes requests to available models based on defined
> policies.
>
> 1.  **Model failure:** If the currently selected model fails (for example, due
>     to quota or server errors), the CLI will initiate the fallback process.
>
> 2.  **User consent:** Depending on the failure and the model's policy, the CLI
>     may prompt you to switch to a fallback model (by default always prompts
>     you).
>
>     Some internal utility calls (such as prompt completion and classification)
>     use a silent fallback chain for `gemini-2.5-flash-lite` and will fall back
>     to `gemini-2.5-flash` and `gemini-2.5-pro` without prompting or changing the
>     configured model.
>
> 3.  **Model switch:** If approved, or if the policy allows for silent fallback,
>     the CLI will use an available fallback model for the current turn or the
>     remainder of the session.
>
> ### Local Model Routing (Experimental)
>
> Gemini CLI supports using a local model for routing decisions. When configured,
> Gemini CLI will use a locally-running **Gemma** model to make routing decisions
> (instead of sending routing decisions to a hosted model). This feature can help
> reduce costs associated with hosted model usage while offering similar routing
> decision latency and quality.
>
> The easiest way to set this up is using the automated `gemini gemma setup`
> command.
>
> For more details on how to configure local model routing, see
> [`gemini gemma` — Local Model Routing Setup](/docs/core/gemma-setup).
>
> ### Model selection precedence
>
> The model used by Gemini CLI is determined by the following order of precedence:
>
> 1.  **`--model` command-line flag:** A model specified with the `--model` flag
>     when launching the CLI will always be used.
> 2.  **`GEMINI_MODEL` environment variable:** If the `--model` flag is not used,
>     the CLI will use the model specified in the `GEMINI_MODEL` environment
>     variable.
> 3.  **`model.name` in `settings.json`:** If neither of the above are set, the
>     model specified in the `model.name` property of your `settings.json` file
>     will be used.
> 4.  **Local model (experimental):** If the Gemma local model router is enabled
>     in your `settings.json` file, the CLI will use the local Gemma model
>     (instead of Gemini models) to route the request to an appropriate model.
> 5.  **Default model:** If none of the above are set, the default model will be
>     used. The default model is `auto`

### Source: Gemini CLI model selection (`/model` command) — Full page

> Select your Gemini CLI model. The `/model` command lets you configure the model
> used by Gemini CLI, giving you more control over your results. Use **Pro**
> models for complex tasks and reasoning, **Flash** models for high speed results,
> or the (recommended) **Auto** setting to choose the best model for your tasks.
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > The `/model` command (and the `--model` flag) does not override the
> > model used by sub-agents. Consequently, even when using the `/model` flag you
> > may see other models used in your model usage reports.
>
> ## How to use the `/model` command
>
> Use the following command in Gemini CLI:
>
> ```
> /model
> ```
>
> Running this command will open a dialog with your options:
>
> | Option            | Description                                                    | Models                                       |
> | ----------------- | -------------------------------------------------------------- | -------------------------------------------- |
> | Auto (Gemini 3)   | Let the system choose the best Gemini 3 model for your task.   | gemini-3-pro-preview, gemini-3-flash-preview |
> | Auto (Gemini 2.5) | Let the system choose the best Gemini 2.5 model for your task. | gemini-2.5-pro, gemini-2.5-flash             |
> | Manual            | Select a specific model.                                       | Any available model.                         |
>
> We recommend selecting one of the above **Auto** options. However, you can
> select **Manual** to select a specific model from those available.
>
> You can also use the `--model` flag to specify a particular Gemini model on
> startup. For more details, refer to the
> [configuration documentation](/docs/reference/configuration).
>
> Changes to these settings will be applied to all subsequent interactions with
> Gemini CLI.
>
> ## Best practices for model selection
>
> - **Default to Auto.** For most users, the _Auto_ option model provides a
>   balance between speed and performance, automatically selecting the correct
>   model based on the complexity of the task. Example: Developing a web
>   application could include a mix of complex tasks (building architecture and
>   scaffolding the project) and simple tasks (generating CSS).
>
> - **Switch to Pro if you aren't getting the results you want.** If you think you
>   need your model to be a little "smarter," you can manually select Pro. Pro
>   will provide you with the highest levels of reasoning and creativity. Example:
>   A complex or multi-stage debugging task.
>
> - **Switch to Flash or Flash-Lite if you need faster results.** If you need a
>   simple response quickly, Flash or Flash-Lite is the best option. Example:
>   Converting a JSON object to a YAML string.

### Source: Model steering (experimental) — Full page

> Model steering lets you provide real-time guidance and feedback to Gemini CLI
> while it is actively executing a task. This lets you correct course, add missing
> context, or skip unnecessary steps without having to stop and restart the agent.
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > This is an experimental feature currently under active development and
> > may need to be enabled under `/settings`.
>
> Model steering is particularly useful during complex [Plan Mode](/docs/cli/plan-mode)
> workflows or long-running subagent executions where you want to ensure the agent
> stays on the right track.
>
> ## Enabling model steering
>
> Model steering is an experimental feature and is disabled by default. You can
> enable it using the `/settings` command or by updating your `settings.json`
> file.
>
> 1.  Type `/settings` in Gemini CLI.
> 2.  Search for **Model Steering**.
> 3.  Set the value to **true**.
>
> Alternatively, add the following to your `settings.json`:
>
> ```json
> {
>   "experimental": {
>     "modelSteering": true
>   }
> }
> ```
>
> ## Using model steering
>
> When model steering is enabled, Gemini CLI treats any text you type while the
> agent is working as a steering hint.
>
> 1.  Start a task (for example, "Refactor the database service").
> 2.  While the agent is working (the spinner is visible), type your feedback in
>     the input box.
> 3.  Press **Enter**.
>
> Gemini CLI acknowledges your hint with a brief message and injects it directly
> into the model's context for the very next turn. The model then re-evaluates its
> current plan and adjusts its actions accordingly.
>
> ### Common use cases
>
> You can use steering hints to guide the model in several ways:
>
> - **Correcting a path:** "Actually, the utilities are in `src/common/utils`."
> - **Skipping a step:** "Skip the unit tests for now and just focus on the
>   implementation."
> - **Adding context:** "The `User` type is defined in `packages/core/types.ts`."
> - **Redirecting the effort:** "Stop searching the codebase and start drafting
>   the plan now."
> - **Handling ambiguity:** "Use the existing `Logger` class instead of creating a
>   new one."
>
> ## How it works
>
> When you submit a steering hint, Gemini CLI performs the following actions:
>
> 1.  **Immediate acknowledgment:** It uses a small, fast model to generate a
>     one-sentence acknowledgment so you know your hint was received.
> 2.  **Context injection:** It prepends an internal instruction to your hint that
>     tells the main agent to:
>     - Re-evaluate the active plan.
>     - Classify the update (for example, as a new task or extra context).
>     - Apply minimal-diff changes to affected tasks.
> 3.  **Real-time update:** The hint is delivered to the agent at the beginning of
>     its next turn, ensuring the most immediate course correction possible.
>
> ## Next steps
>
> - Tackle complex tasks with [Plan Mode](/docs/cli/plan-mode).
> - Build custom [Agent Skills](/docs/cli/skills).
