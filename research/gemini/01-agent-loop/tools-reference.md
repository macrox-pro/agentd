---
primary_sources:
  - id: T1-TOOLS
    title: "Tools reference"
    url: "https://geminicli.com/docs/reference/tools.md"
    section: "Full page"
  - id: T1-ASK-USER
    title: "Ask User Tool"
    url: "https://geminicli.com/docs/tools/ask-user.md"
    section: "Full page"
  - id: T1-FS
    title: "File system tools reference"
    url: "https://geminicli.com/docs/tools/file-system.md"
    section: "Full page"
  - id: T1-SHELL
    title: "Shell tool (`run_shell_command`)"
    url: "https://geminicli.com/docs/tools/shell.md"
    section: "Full page"
  - id: T1-PLAN-TOOLS
    title: "Gemini CLI planning tools"
    url: "https://geminicli.com/docs/tools/planning.md"
    section: "Full page"
  - id: T1-TODOS
    title: "Todo tool (`write_todos`)"
    url: "https://geminicli.com/docs/tools/todos.md"
    section: "Full page"
  - id: T1-TRACKER
    title: "Tracker tools (`tracker_*`)"
    url: "https://geminicli.com/docs/tools/tracker.md"
    section: "Full page"
  - id: T1-WEB-SEARCH
    title: "Web search tool (`google_web_search`)"
    url: "https://geminicli.com/docs/tools/web-search.md"
    section: "Full page"
  - id: T1-WEB-FETCH
    title: "Web fetch tool (`web_fetch`)"
    url: "https://geminicli.com/docs/tools/web-fetch.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Tools reference

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Tools reference — Full page

> Gemini CLI uses tools to interact with your local environment, access
> information, and perform actions on your behalf. These tools extend the model's
> capabilities beyond text generation, letting it read files, execute commands,
> and search the web.
>
> ## How to use Gemini CLI's tools
>
> Tools are generally invoked automatically by Gemini CLI when it needs to perform
> an action. However, you can also trigger specific tools manually using shorthand
> syntax.
>
> ### Automatic execution and security
>
> When the model wants to use a tool, Gemini CLI evaluates the request against its
> security policies.
>
> - **User confirmation:** You must manually approve tools that modify files or
>   execute shell commands (mutators). The CLI shows you a diff or the exact
>   command before you confirm.
> - **Sandboxing:** You can run tool executions in secure, containerized
>   environments to isolate changes from your host system. For more details, see
>   the [Sandboxing](/docs/cli/sandbox) guide.
> - **Trusted folders:** You can configure which directories allow the model to
>   use system tools. For more details, see the
>   [Trusted folders](/docs/cli/trusted-folders) guide.
>
> Review confirmation prompts carefully before allowing a tool to execute.
>
> ### How to use manually-triggered tools
>
> You can directly trigger key tools using special syntax in your prompt:
>
> - **[File access](/docs/tools/file-system#read_many_files) (`@`):** Use the `@`
>   symbol followed by a file or directory path to include its content in your
>   prompt. This triggers the `read_many_files` tool.
> - **[Shell commands](/docs/tools/shell) (`!`):** Use the `!` symbol followed by
>   a system command to execute it directly. This triggers the `run_shell_command`
>   tool.
>
> ## How to manage tools
>
> Using built-in commands, you can inspect available tools and configure how they
> behave.
>
> ### Tool discovery
>
> Use the `/tools` command to see what tools are currently active in your session.
>
> - **`/tools`**: Lists all registered tools with their display names.
> - **`/tools desc`**: Lists all tools with their full descriptions.
>
> This is especially useful for verifying that
> [MCP servers](/docs/tools/mcp-server) or custom tools are loaded correctly.
>
> ### Tool configuration
>
> You can enable, disable, or configure specific tools in your settings. For
> example, you can set a specific pager for shell commands or configure the
> browser used for web searches. See the [Settings](/docs/cli/settings) guide for
> details.
>
> ## Available tools
>
> The following sections list all available tools, categorized by their primary
> function. For detailed parameter information, see the linked documentation for
> each tool.
>
> ### Execution
>
> | Tool                                     | Kind      | Description                                                                                                              |
> | :--------------------------------------- | :-------- | :----------------------------------------------------------------------------------------------------------------------- |
> | [`run_shell_command`](/docs/tools/shell) | `Execute` | Executes arbitrary shell commands. Supports interactive sessions and background processes. Requires manual confirmation. |
>
> ### File System
>
> | Tool                                         | Kind     | Description                                                                                           |
> | :------------------------------------------- | :------- | :---------------------------------------------------------------------------------------------------- |
> | [`glob`](/docs/tools/file-system)            | `Search` | Finds files matching specific glob patterns across the workspace.                                     |
> | [`grep_search`](/docs/tools/file-system)     | `Search` | Searches for a regular expression pattern within file contents. Legacy alias: `search_file_content`.  |
> | [`list_directory`](/docs/tools/file-system)  | `Read`   | Lists the names of files and subdirectories within a specified path.                                  |
> | [`read_file`](/docs/tools/file-system)       | `Read`   | Reads the content of a specific file. Supports text, images, audio, and PDF.                          |
> | [`read_many_files`](/docs/tools/file-system) | `Read`   | Reads and concatenates content from multiple files. Often triggered by the `@` symbol in your prompt. |
> | [`replace`](/docs/tools/file-system)         | `Edit`   | Performs precise text replacement within a file. Requires manual confirmation.                        |
> | [`write_file`](/docs/tools/file-system)      | `Edit`   | Creates or overwrites a file with new content. Requires manual confirmation.                          |
>
> ### Interaction
>
> | Tool                               | Kind          | Description                                                                            |
> | :--------------------------------- | :------------ | :------------------------------------------------------------------------------------- |
> | [`ask_user`](/docs/tools/ask-user) | `Communicate` | Requests clarification or missing information via an interactive dialog.               |
> | [`write_todos`](/docs/tools/todos) | `Other`       | Maintains an internal list of subtasks. The model uses this to track its own progress. |
>
> ### Task Tracker (Experimental)
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > This is an experimental feature currently under active development. Enable via `experimental.taskTracker`.
>
> | Tool                                            | Kind    | Description                                                                 |
> | :---------------------------------------------- | :------ | :-------------------------------------------------------------------------- |
> | [`tracker_create_task`](/docs/tools/tracker)    | `Other` | Creates a new task in the experimental tracker.                             |
> | [`tracker_update_task`](/docs/tools/tracker)    | `Other` | Updates an existing task's status, description, or dependencies.            |
> | [`tracker_get_task`](/docs/tools/tracker)       | `Other` | Retrieves the full details of a specific task.                              |
> | [`tracker_list_tasks`](/docs/tools/tracker)     | `Other` | Lists tasks in the tracker, optionally filtered by status, type, or parent. |
> | [`tracker_add_dependency`](/docs/tools/tracker) | `Other` | Adds a dependency between two tasks, ensuring topological execution.        |
> | [`tracker_visualize`](/docs/tools/tracker)      | `Other` | Renders an ASCII tree visualization of the current task graph.              |
>
> ### MCP
>
> | Tool                                              | Kind     | Description                                                            |
> | :------------------------------------------------ | :------- | :--------------------------------------------------------------------- |
> | [`list_mcp_resources`](/docs/tools/mcp-resources) | `Search` | Lists all available resources exposed by connected MCP servers.        |
> | [`read_mcp_resource`](/docs/tools/mcp-resources)  | `Read`   | Reads the content of a specific Model Context Protocol (MCP) resource. |
>
> ### Memory
>
> | Tool                                             | Kind    | Description                                                                          |
> | :----------------------------------------------- | :------ | :----------------------------------------------------------------------------------- |
> | [`activate_skill`](/docs/tools/activate-skill)   | `Other` | Loads specialized procedural expertise from the `.gemini/skills` directory.          |
> | [`get_internal_docs`](/docs/tools/internal-docs) | `Think` | Accesses Gemini CLI's own documentation for accurate answers about its capabilities. |
>
> ### Planning
>
> | Tool                                      | Kind   | Description                                                                              |
> | :---------------------------------------- | :----- | :--------------------------------------------------------------------------------------- |
> | [`enter_plan_mode`](/docs/tools/planning) | `Plan` | Switches the CLI to a safe, read-only "Plan Mode" for researching complex changes.       |
> | [`exit_plan_mode`](/docs/tools/planning)  | `Plan` | Finalizes a plan, presents it for review, and requests approval to start implementation. |
>
> ### System
>
> | Tool            | Kind    | Description                                                                                                        |
> | :-------------- | :------ | :----------------------------------------------------------------------------------------------------------------- |
> | `complete_task` | `Other` | Finalizes a subagent's mission and returns the result to the parent agent. This tool is not available to the user. |
>
> ### Task Tracking
>
> | Tool                     | Kind    | Description                                                                 |
> | :----------------------- | :------ | :-------------------------------------------------------------------------- |
> | `tracker_add_dependency` | `Think` | Adds a dependency between two existing tasks in the tracker.                |
> | `tracker_create_task`    | `Think` | Creates a new task in the internal tracker to monitor progress.             |
> | `tracker_get_task`       | `Think` | Retrieves the details and current status of a specific tracked task.        |
> | `tracker_list_tasks`     | `Think` | Lists all tasks currently being tracked.                                    |
> | `tracker_update_task`    | `Think` | Updates the status or details of an existing task.                          |
> | `tracker_visualize`      | `Think` | Generates a visual representation of the current task dependency graph.     |
> | `update_topic`           | `Think` | Updates the current topic and status to keep the user informed of progress. |
>
> ### Web
>
> | Tool                                          | Kind     | Description                                                                                                                                                                                                                                                                     |
> | :-------------------------------------------- | :------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | [`google_web_search`](/docs/tools/web-search) | `Search` | Performs a Google Search to find up-to-date information.                                                                                                                                                                                                                        |
> | [`web_fetch`](/docs/tools/web-fetch)          | `Fetch`  | Retrieves and processes content from specific URLs. **Warning:** This tool can access local and private network addresses (for example, localhost), which may pose a security risk if used with untrusted prompts. In Plan Mode, this tool requires explicit user confirmation. |
>
> ### Tool argument keys
>
> When writing [`argsPattern`](/docs/reference/policy-engine#arguments-pattern) rules for the
> [policy engine](/docs/reference/policy-engine), you need to know the JSON argument keys for
> each tool. The following table lists the keys that appear in the JSON
> representation of each tool's arguments.
>
> | Tool                     | JSON argument keys                                                                                                                                                                                   |
> | :----------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `run_shell_command`      | `command`, `description`, `dir_path`, `is_background`                                                                                                                                                |
> | `glob`                   | `pattern`, `dir_path`, `case_sensitive`, `respect_git_ignore`, `respect_gemini_ignore`                                                                                                               |
> | `grep_search`            | `pattern`, `dir_path`, `include_pattern`, `exclude_pattern`, `names_only`, `case_sensitive`, `fixed_strings`, `context`, `after`, `before`, `no_ignore`, `max_matches_per_file`, `total_max_matches` |
> | `list_directory`         | `dir_path`, `ignore`, `file_filtering_options`                                                                                                                                                       |
> | `read_file`              | `file_path`, `start_line`, `end_line`                                                                                                                                                                |
> | `read_many_files`        | `include`, `exclude`, `recursive`, `useDefaultExcludes`                                                                                                                                              |
> | `write_file`             | `file_path`, `content`                                                                                                                                                                               |
> | `replace`                | `file_path`, `old_string`, `new_string`, `instruction`, `allow_multiple`                                                                                                                             |
> | `ask_user`               | `questions` (array of `question`, `header`, `type`, `options`)                                                                                                                                       |
> | `write_todos`            | `todos` (array of `description`, `status`)                                                                                                                                                           |
> | `activate_skill`         | `name`                                                                                                                                                                                               |
> | `get_internal_docs`      | `path`                                                                                                                                                                                               |
> | `enter_plan_mode`        | `reason`                                                                                                                                                                                             |
> | `exit_plan_mode`         | `plan_path`                                                                                                                                                                                          |
> | `tracker_create_task`    | `title`, `description`, `type`                                                                                                                                                                       |
> | `tracker_update_task`    | `id`, `title`, `description`, `status`, `dependencies`                                                                                                                                               |
> | `tracker_get_task`       | `id`                                                                                                                                                                                                 |
> | `tracker_list_tasks`     | `status`, `type`, `parentId`                                                                                                                                                                         |
> | `tracker_add_dependency` | `taskId`, `dependencyId`                                                                                                                                                                             |
> | `tracker_visualize`      | _(none)_                                                                                                                                                                                             |
> | `update_topic`           | `title`, `summary`, `strategic_intent`                                                                                                                                                               |
> | `google_web_search`      | `query`                                                                                                                                                                                              |
> | `web_fetch`              | `prompt`                                                                                                                                                                                             |
>
> For example, to write a policy rule that blocks any `write_file` call targeting
> a `.env` file, you would match against the `file_path` key:
>
> ```toml
> [[rule]]
> toolName = "write_file"
> argsPattern = '"file_path":".*\.env"'
> decision = "deny"
> priority = 100
> denyMessage = "Writing to .env files is not allowed."
> ```
>
> For full argument descriptions and types, see the individual tool pages linked
> in the [tables above](#available-tools).
>
> ## Under the hood
>
> For developers, the tool system is designed to be extensible and robust. The
> `ToolRegistry` class manages all available tools.
>
> You can extend Gemini CLI with custom tools by configuring
> `tools.discoveryCommand` in your settings or by connecting to MCP servers.
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > For a deep dive into the internal Tool API and how to implement your
> > own tools in the codebase, see the `packages/core/src/tools/` directory in
> > GitHub.
>
> ## Next steps
>
> - Learn how to [Set up an MCP server](/docs/tools/mcp-server).
> - Explore [Agent Skills](/docs/cli/skills) for specialized expertise.
> - See the [Command reference](/docs/reference/commands) for slash commands.

### Source: Ask User Tool — Full page

> The `ask_user` tool lets Gemini CLI ask you one or more questions to gather
> preferences, clarify requirements, or make decisions. It supports multiple
> question types including multiple-choice, free-form text, and Yes/No
> confirmation.
>
> ## `ask_user` (Ask User)
>
> - **Tool name:** `ask_user`
> - **Display name:** Ask User
> - **File:** `ask-user.ts`
> - **Parameters:**
>
>   - `questions` (array of objects, required): A list of 1 to 4 questions to ask.
>     Each question object has the following properties:
>     - `question` (string, required): The complete question text.
>     - `header` (string, required): A short label (max 16 chars) displayed as a
>       chip/tag (for example, "Auth", "Database").
>     - `type` (string, optional): The type of question. Defaults to `'choice'`.
>       - `'choice'`: Multiple-choice with options (supports multi-select).
>       - `'text'`: Free-form text input.
>       - `'yesno'`: Yes/No confirmation.
>     - `options` (array of objects, optional): Required for `'choice'` type. 2-4
>       selectable options.
>       - `label` (string, required): Display text (1-5 words).
>       - `description` (string, required): Brief explanation.
>     - `multiSelect` (boolean, optional): For `'choice'` type, allows selecting
>       multiple options. Automatically adds an "All the above" option if there
>       are multiple standard options.
>     - `placeholder` (string, optional): Hint text for input fields.
>
> - **Behavior:**
>
>   - Presents an interactive dialog to the user with the specified questions.
>   - Pauses execution until the user provides answers or dismisses the dialog.
>   - Returns the user's answers to the model.
>
> - **Output (`llmContent`):** A JSON string containing the user's answers,
>   indexed by question position (for example,
>   `{"answers":{"0": "Option A", "1": "Some text"}}`).
>
> - **Confirmation:** Yes. The tool inherently involves user interaction.
>
> ## Usage Examples
>
> ### Multiple Choice Question
>
> ```json
> {
>   "questions": [
>     {
>       "header": "Database",
>       "question": "Which database would you like to use?",
>       "type": "choice",
>       "options": [
>         {
>           "label": "PostgreSQL",
>           "description": "Powerful, open source object-relational database system."
>         },
>         {
>           "label": "SQLite",
>           "description": "C-library that implements a SQL database engine."
>         }
>       ]
>     }
>   ]
> }
> ```
>
> ### Text Input Question
>
> ```json
> {
>   "questions": [
>     {
>       "header": "Project Name",
>       "question": "What is the name of your new project?",
>       "type": "text",
>       "placeholder": "for example, my-awesome-app"
>     }
>   ]
> }
> ```
>
> ### Yes/No Question
>
> ```json
> {
>   "questions": [
>     {
>       "header": "Deploy",
>       "question": "Do you want to deploy the application now?",
>       "type": "yesno"
>     }
>   ]
> }
> ```

### Source: File system tools reference — Full page

> Gemini CLI core provides a suite of tools for interacting with the local file
> system. These tools allow the model to explore and modify your codebase.
>
> ## Technical reference
>
> All file system tools operate within a `rootDirectory` (the current working
> directory or workspace root) for security.
>
> ### `list_directory` (ReadFolder)
>
> Lists the names of files and subdirectories directly within a specified path.
>
> - **Tool name:** `list_directory`
> - **Arguments:**
>   - `dir_path` (string, required): Absolute or relative path to the directory.
>   - `ignore` (array, optional): Glob patterns to exclude.
>   - `file_filtering_options` (object, optional): Configuration for `.gitignore`
>     and `.geminiignore` compliance.
>
> ### `read_file` (ReadFile)
>
> Reads and returns the content of a specific file. Supports text, images, audio,
> and PDF.
>
> - **Tool name:** `read_file`
> - **Arguments:**
>   - `file_path` (string, required): Path to the file.
>   - `offset` (number, optional): Start line for text files (0-based).
>   - `limit` (number, optional): Maximum lines to read.
>
> ### `write_file` (WriteFile)
>
> Writes content to a specified file, overwriting it if it exists or creating it
> if not.
>
> - **Tool name:** `write_file`
> - **Arguments:**
>   - `file_path` (string, required): Path to the file.
>   - `content` (string, required): Data to write.
> - **Confirmation:** Requires manual user approval.
>
> ### `glob` (FindFiles)
>
> Finds files matching specific glob patterns across the workspace.
>
> - **Tool name:** `glob`
> - **Display name:** FindFiles
> - **File:** `glob.ts`
> - **Parameters:**
>   - `pattern` (string, required): The glob pattern to match against (for
>     example, `"*.py"`, `"src/**/*.js"`).
>   - `path` (string, optional): The absolute path to the directory to search
>     within. If omitted, searches the tool's root directory.
>   - `case_sensitive` (boolean, optional): Whether the search should be
>     case-sensitive. Defaults to `false`.
>   - `respect_git_ignore` (boolean, optional): Whether to respect .gitignore
>     patterns when finding files. Defaults to `true`.
> - **Behavior:**
>   - Searches for files matching the glob pattern within the specified directory.
>   - Returns a list of absolute paths, sorted with the most recently modified
>     files first.
>   - Ignores common nuisance directories like `node_modules` and `.git` by
>     default.
> - **Output (`llmContent`):** A message like:
>   `Found 5 file(s) matching "*.ts" within src, sorted by modification time (newest first):\nsrc/file1.ts\nsrc/subdir/file2.ts...`
> - **Confirmation:** No.
>
> ### `grep_search` (SearchText)
>
> `grep_search` searches for a regular expression pattern within the content of
> files in a specified directory. Can filter files by a glob pattern. Returns the
> lines containing matches, along with their file paths and line numbers.
>
> - **Tool name:** `grep_search`
> - **Display name:** SearchText
> - **File:** `grep.ts`
> - **Parameters:**
>   - `pattern` (string, required): The regular expression (regex) to search for
>     (for example, `"function\s+myFunction"`).
>   - `path` (string, optional): The absolute path to the directory to search
>     within. Defaults to the current working directory.
>   - `include` (string, optional): A glob pattern to filter which files are
>     searched (for example, `"*.js"`, `"src/**/*.{ts,tsx}"`). If omitted,
>     searches most files (respecting common ignores).
> - **Behavior:**
>   - Uses `git grep` if available in a Git repository for speed; otherwise, falls
>     back to system `grep` or a JavaScript-based search.
>   - Returns a list of matching lines, each prefixed with its file path (relative
>     to the search directory) and line number.
> - **Output (`llmContent`):** A formatted string of matches, for example:
>   ```
>   Found 3 matches for pattern "myFunction" in path "." (filter: "*.ts"):
>   ---
>   File: src/utils.ts
>   L15: export function myFunction() {
>   L22:   myFunction.call();
>   ---
>   File: src/index.ts
>   L5: import { myFunction } from './utils';
>   ---
>   ```
> - **Confirmation:** No.
>
> ### `replace` (Edit)
>
> `replace` replaces text within a file. By default, the tool expects to find and
> replace exactly ONE occurrence of `old_string`. If you want to replace multiple
> occurrences of the exact same string, set `allow_multiple` to `true`. This tool
> is designed for precise, targeted changes and requires significant context
> around the `old_string` to ensure it modifies the correct location.
>
> - **Tool name:** `replace`
> - **Arguments:**
>   - `file_path` (string, required): Path to the file.
>   - `instruction` (string, required): Semantic description of the change.
>   - `old_string` (string, required): Exact literal text to find.
>   - `new_string` (string, required): Exact literal text to replace with.
>   - `allow_multiple` (boolean, optional): If `true`, replaces all occurrences.
>     If `false` (default), only succeeds if exactly one occurrence is found.
> - **Confirmation:** Requires manual user approval.
>
> ## Next steps
>
> - Follow the [File management tutorial](/docs/cli/tutorials/file-management) for
>   practical examples.
> - Learn about [Trusted folders](/docs/cli/trusted-folders) to manage access
>   permissions.

### Source: Shell tool (`run_shell_command`) — Full page

> The `run_shell_command` tool allows the Gemini model to execute commands
> directly on your system's shell. It is the primary mechanism for the agent to
> interact with your environment beyond simple file edits.
>
> ## Technical reference
>
> On Windows, commands execute with `powershell.exe -NoProfile -Command`. On other
> platforms, they execute with `bash -c`.
>
> ### Arguments
>
> - `command` (string, required): The exact shell command to execute.
> - `description` (string, optional): A brief description shown to the user for
>   confirmation.
> - `dir_path` (string, optional): The absolute path or relative path from
>   workspace root where the command runs.
> - `is_background` (boolean, optional): Whether to move the process to the
>   background immediately after starting.
>
> ### Policy engine shorthands
>
> The [policy engine](/docs/reference/policy-engine) provides two convenience
> fields for writing rules that target shell commands:
>
> - `commandPrefix`: Matches if the `command` argument starts with a given string.
> - `commandRegex`: Matches if the `command` argument matches a given regular
>   expression.
>
> These are syntactic sugar for combining `toolName = "run_shell_command"` with an
> `argsPattern` in a policy TOML file. They are **not** arguments of
> `run_shell_command` itself.
>
> For details on writing shell-specific policy rules, see
> [Special syntax for `run_shell_command`](/docs/reference/policy-engine#special-syntax-for-run_shell_command)
> in the policy engine reference.
>
> ### Return values
>
> The tool returns a JSON object containing:
>
> - `Command`: The executed string.
> - `Directory`: The execution path.
> - `Stdout` / `Stderr`: The output streams.
> - `Exit Code`: The process return code.
> - `Background PIDs`: PIDs of any started background processes.
>
> ## Configuration
>
> You can configure the behavior of the `run_shell_command` tool by modifying your
> `settings.json` file or by using the `/settings` command in Gemini CLI.
>
> ### Enabling interactive commands
>
> To enable interactive commands, you need to set the
> `tools.shell.enableInteractiveShell` setting to `true`. This will use `node-pty`
> for shell command execution, which allows for interactive sessions. If
> `node-pty` is not available, it will fall back to the `child_process`
> implementation, which does not support interactive commands.
>
> **Example `settings.json`:**
>
> ```json
> {
>   "tools": {
>     "shell": {
>       "enableInteractiveShell": true
>     }
>   }
> }
> ```
>
> ### Showing color in output
>
> To show color in the shell output, you need to set the `tools.shell.showColor`
> setting to `true`. This setting only applies when
> `tools.shell.enableInteractiveShell` is enabled.
>
> **Example `settings.json`:**
>
> ```json
> {
>   "tools": {
>     "shell": {
>       "showColor": true
>     }
>   }
> }
> ```
>
> ### Setting the pager
>
> You can set a custom pager for the shell output by setting the
> `tools.shell.pager` setting. The default pager is `cat`. This setting only
> applies when `tools.shell.enableInteractiveShell` is enabled.
>
> **Example `settings.json`:**
>
> ```json
> {
>   "tools": {
>     "shell": {
>       "pager": "less"
>     }
>   }
> }
> ```
>
> ## Interactive commands
>
> The `run_shell_command` tool now supports interactive commands by integrating a
> pseudo-terminal (pty). This lets you run commands that require real-time user
> input, such as text editors (`vim`, `nano`), terminal-based UIs (`htop`), and
> interactive version control operations (`git rebase -i`).
>
> When an interactive command is running, you can send input to it from the Gemini
> CLI. To focus on the interactive shell, press `Tab`. The terminal output,
> including complex TUIs, will be rendered correctly.
>
> ## Important notes
>
> - **Security:** Be cautious when executing commands, especially those
>   constructed from user input, to prevent security vulnerabilities.
> - **Error handling:** Check the `Stderr`, `Error`, and `Exit Code` fields to
>   determine if a command executed successfully.
> - **Background processes:** When a command is run in the background with `&`,
>   the tool will return immediately and the process will continue to run in the
>   background. The `Background PIDs` field will contain the process ID of the
>   background process.
>
> ## Environment variables
>
> When `run_shell_command` executes a command, it sets the `GEMINI_CLI=1`
> environment variable in the subprocess's environment. This allows scripts or
> tools to detect if they are being run from within Gemini CLI.
>
> ## Command restrictions
>
> <!-- prettier-ignore -->
> > [!WARNING]
> > The `tools.core` setting is an **allowlist for _all_ built-in
> > tools**, not just shell commands. When you set `tools.core` to any value,
> > _only_ the tools explicitly listed will be enabled. This includes all built-in
> > tools like `read_file`, `write_file`, `glob`, `grep_search`, `list_directory`,
> > `replace`, etc.
>
> You can restrict the commands that can be executed by the `run_shell_command`
> tool by using the `tools.core` and `tools.exclude` settings in your
> configuration file.
>
> - `tools.core`: To restrict `run_shell_command` to a specific set of commands,
>   add entries to the `core` list under the `tools` category in the format
>   `run_shell_command(<command>)`. For example,
>   `"tools": {"core": ["run_shell_command(git)"]}` will only allow `git`
>   commands. Including the generic `run_shell_command` acts as a wildcard,
>   allowing any command not explicitly blocked.
> - `tools.exclude` [DEPRECATED]: To block specific commands, use the
>   [Policy Engine](/docs/reference/policy-engine). Historically, this setting
>   allowed adding entries to the `exclude` list under the `tools` category in the
>   format `run_shell_command(<command>)`. For example,
>   `"tools": {"exclude": ["run_shell_command(rm)"]}` will block `rm` commands.
>
> The validation logic is designed to be secure and flexible:
>
> 1.  **Command chaining disabled**: The tool automatically splits commands
>     chained with `&&`, `||`, or `;` and validates each part separately. If any
>     part of the chain is disallowed, the entire command is blocked.
> 2.  **Prefix matching**: The tool uses prefix matching. For example, if you
>     allow `git`, you can run `git status` or `git log`.
> 3.  **Blocklist precedence**: The `tools.exclude` list is always checked first.
>     If a command matches a blocked prefix, it will be denied, even if it also
>     matches an allowed prefix in `tools.core`.
>
> ### Command restriction examples
>
> **Allow only specific command prefixes**
>
> To allow only `git` and `npm` commands, and block all others:
>
> ```json
> {
>   "tools": {
>     "core": ["run_shell_command(git)", "run_shell_command(npm)"]
>   }
> }
> ```
>
> - `git status`: Allowed
> - `npm install`: Allowed
> - `ls -l`: Blocked
>
> **Block specific command prefixes**
>
> To block `rm` and allow all other commands:
>
> ```json
> {
>   "tools": {
>     "core": ["run_shell_command"],
>     "exclude": ["run_shell_command(rm)"]
>   }
> }
> ```
>
> - `rm -rf /`: Blocked
> - `git status`: Allowed
> - `npm install`: Allowed
>
> **Blocklist takes precedence**
>
> If a command prefix is in both `tools.core` and `tools.exclude`, it will be
> blocked.
>
> - **`tools.shell.enableInteractiveShell`**: (boolean) Uses `node-pty` for
>   real-time interaction.
> - **`tools.shell.showColor`**: (boolean) Preserves ANSI colors in output.
> - **`tools.shell.inactivityTimeout`**: (number) Seconds to wait for output
>   before killing the process.
>
> ### Command restrictions
>
> You can limit which commands the agent is allowed to request using these
> settings:
>
> - **`tools.core`**: An allowlist of command prefixes (for example,
>   `["git", "npm test"]`).
> - **`tools.exclude`**: A blocklist of command prefixes.
>
> ## Use cases
>
> - Running build scripts and test suites.
> - Initializing or managing version control systems.
> - Installing project dependencies.
> - Starting development servers or background watchers.
>
> ## Next steps
>
> - Follow the [Shell commands tutorial](/docs/cli/tutorials/shell-commands) for
>   practical examples.
> - Learn about [Sandboxing](/docs/cli/sandbox) to isolate command execution.

### Source: Gemini CLI planning tools — Full page

> Planning tools let Gemini CLI switch into a safe, read-only "Plan Mode" for
> researching and planning complex changes, and to signal the finalization of a
> plan to the user.
>
> ## 1. `enter_plan_mode` (EnterPlanMode)
>
> `enter_plan_mode` switches the CLI to Plan Mode. This tool is typically called
> by the agent when you ask it to "start a plan" using natural language. In this
> mode, the agent is restricted to read-only tools to allow for safe exploration
> and planning.
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > This tool is not available when the CLI is in YOLO mode.
>
> - **Tool name:** `enter_plan_mode`
> - **Display name:** Enter Plan Mode
> - **File:** `enter-plan-mode.ts`
> - **Parameters:**
>   - `reason` (string, optional): A short reason explaining why the agent is
>     entering plan mode (for example, "Starting a complex feature
>     implementation").
> - **Behavior:**
>   - Switches the CLI's approval mode to `PLAN`.
>   - Notifies the user that the agent has entered Plan Mode.
> - **Output (`llmContent`):** A message indicating the switch, for example,
>   `Switching to Plan mode.`
> - **Confirmation:** Yes. The user is prompted to confirm entering Plan Mode.
>
> ## 2. `exit_plan_mode` (ExitPlanMode)
>
> `exit_plan_mode` signals that the planning phase is complete. It presents the
> finalized plan to the user and requests formal approval to start the
> implementation. The agent MUST reach an informal agreement with the user in the
> chat regarding the proposed strategy BEFORE calling this tool.
>
> - **Tool name:** `exit_plan_mode`
> - **Display name:** Exit Plan Mode
> - **File:** `exit-plan-mode.ts`
> - **Parameters:**
>   - `plan_path` (string, required): The path to the finalized Markdown plan
>     file. This file MUST be located within the project's temporary plans
>     directory (for example, `~/.gemini/tmp/<project>/plans/`).
> - **Behavior:**
>   - Validates that the `plan_path` is within the allowed directory and that the
>     file exists and has content.
>   - Presents the plan to the user for formal review.
>   - If the user approves the plan:
>     - Switches the CLI's approval mode to the user's chosen approval mode (
>       `DEFAULT` or `AUTO_EDIT`).
>     - Marks the plan as approved for implementation.
>   - If the user rejects the plan:
>     - Stays in Plan Mode.
>     - Returns user feedback to the model to refine the plan.
> - **Output (`llmContent`):**
>   - On approval: A message indicating the plan was approved and the new approval
>     mode.
>   - On rejection: A message containing the user's feedback.
> - **Confirmation:** Yes. Shows the finalized plan and asks for user formal
>   approval to proceed with implementation.

### Source: Todo tool (`write_todos`) — Full page

> The `write_todos` tool allows the Gemini agent to maintain an internal list of
> subtasks for multi-step requests.
>
> ## Technical reference
>
> The agent uses this tool to manage its execution plan and provide progress
> updates to the CLI interface.
>
> ### Arguments
>
> - `todos` (array of objects, required): The complete list of tasks. Each object
>   includes:
>   - `description` (string): Technical description of the task.
>   - `status` (enum): `pending`, `in_progress`, `completed`, `cancelled`, or
>     `blocked`.
>
> ## Technical behavior
>
> - **Interface:** Updates the progress indicator above the CLI input prompt.
> - **Exclusivity:** Only one task can be marked `in_progress` at any time.
> - **Persistence:** Todo state is scoped to the current session.
> - **Interaction:** Users can toggle the full list view using **Ctrl+T**.
>
> ## Use cases
>
> - Breaking down a complex feature implementation into manageable steps.
> - Coordinating multi-file refactoring tasks.
> - Providing visibility into the agent's current focus during long-running tasks.
>
> ## Next steps
>
> - Follow the [Task planning tutorial](/docs/cli/tutorials/task-planning) for
>   usage details.
> - Learn about [Session management](/docs/cli/session-management) for context.

### Source: Tracker tools (`tracker_*`) — Full page

> <!-- prettier-ignore -->
> > [!NOTE]
> > This is an experimental feature currently under active development.
>
> The `tracker_*` tools allow the Gemini agent to maintain an internal, persistent
> graph of tasks and dependencies for multi-step requests. This suite of tools
> provides a more robust and granular way to manage execution plans than the
> legacy `write_todos` tool.
>
> ## Technical reference
>
> The agent uses these tools to manage its execution plan, decompose complex goals
> into actionable sub-tasks, and provide real-time progress updates to the CLI
> interface. The task state is stored in the `.gemini/tmp/tracker/<session-id>`
> directory, allowing the agent to manage its plan for the current session.
>
> ### Available Tools
>
> - `tracker_create_task`: Creates a new task in the tracker. You can specify a
>   title, description, and task type (`epic`, `task`, `bug`).
> - `tracker_update_task`: Updates an existing task's status (`open`,
>   `in_progress`, `blocked`, `closed`), description, or dependencies.
> - `tracker_get_task`: Retrieves the full details of a specific task by its
>   6-character hex ID.
> - `tracker_list_tasks`: Lists tasks in the tracker, optionally filtered by
>   status, type, or parent ID.
> - `tracker_add_dependency`: Adds a dependency between two tasks, ensuring
>   topological execution.
> - `tracker_visualize`: Renders an ASCII tree visualization of the current task
>   graph.
>
> ## Technical behavior
>
> - **Interface:** Updates the progress indicator and task tree above the CLI
>   input prompt.
> - **Persistence:** Task state is saved automatically to the
>   `.gemini/tmp/tracker/<session-id>` directory. Task states are session-specific
>   and do not persist across different sessions.
> - **Dependencies:** Tasks can depend on other tasks, forming a directed acyclic
>   graph (DAG). The agent must resolve dependencies before starting blocked
>   tasks.
> - **Interaction:** Users can view the current state of the tracker by asking the
>   agent to visualize it, or by running `gemini-cli` commands if implemented.
>
> ## Use cases
>
> - Coordinating multi-file refactoring projects.
> - Breaking down a mission into a hierarchy of epics and tasks for better
>   visibility.
> - Tracking bugs and feature requests directly within the context of an active
>   codebase.
> - Providing visibility into the agent's current focus and remaining work.
>
> ## Next steps
>
> - Follow the [Task planning tutorial](/docs/cli/tutorials/task-planning) for
>   usage details and migration from the legacy todo list.
> - Learn about [Session management](/docs/cli/session-management) for context on
>   persistent state.

### Source: Web search tool (`google_web_search`) — Full page

> The `google_web_search` tool allows the Gemini agent to retrieve up-to-date
> information, news, and facts from the internet via Google Search.
>
> ## Technical reference
>
> The agent uses this tool when your request requires knowledge of current events
> or specific online documentation not available in its internal training data.
>
> ### Arguments
>
> - `query` (string, required): The search query to be executed.
>
> ## Technical behavior
>
> - **Grounding:** Returns a generated summary based on search results.
> - **Citations:** Includes source URIs and titles for factual grounding.
> - **Processing:** The Gemini API processes the search results before returning a
>   synthesized response to the agent.
>
> ## Use cases
>
> - Researching the latest version of a software library or API.
> - Finding solutions to recent software bugs or security vulnerabilities.
> - Retrieving news or documentation updated after the model's knowledge cutoff.
>
> ## Next steps
>
> - Follow the [Web tools guide](/docs/cli/tutorials/web-tools) for practical
>   usage examples.
> - Explore the [Web fetch tool reference](/docs/tools/web-fetch) for direct URL access.

### Source: Web fetch tool (`web_fetch`) — Full page

> The `web_fetch` tool allows the Gemini agent to retrieve and process content
> from specific URLs provided in your prompt.
>
> ## Technical reference
>
> The agent uses this tool when you include URLs in your prompt and request
> specific operations like summarization or extraction.
>
> ### Arguments
>
> - `prompt` (string, required): A request containing up to 20 valid URLs
>   (starting with `http://` or `https://`) and instructions on how to process
>   them.
>
> ## Technical behavior
>
> - **Confirmation:** Triggers a confirmation dialog showing the converted URLs.
> - **Plan Mode:** In [Plan Mode](/docs/cli/plan-mode), `web_fetch` is available
>   but always requires explicit user confirmation (`ask_user`) due to security
>   implications of accessing external or private network addresses.
> - **Processing:** Uses the Gemini API's `urlContext` for retrieval.
> - **Fallback:** If API access fails, the tool attempts to fetch raw content
>   directly from your local machine.
> - **Formatting:** Returns a synthesized response with source attribution.
>
> ## Use cases
>
> - Summarizing technical articles or blog posts.
> - Comparing data between two or more web pages.
> - Extracting specific information from a documentation site.
>
> ## Next steps
>
> - Follow the [Web tools guide](/docs/cli/tutorials/web-tools) for practical
>   usage examples.
> - See the [Web search tool reference](/docs/tools/web-search) for general queries.
