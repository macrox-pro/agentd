---
primary_sources:
  - id: T1-CLI-OVERVIEW
    title: "CLI"
    url: "https://cursor.com/docs/cli/overview.md"
    section: "CLI"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# CLI overview and usage

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: CLI installation

> # Installation
>
> ## Installation
>
> ### macOS, Linux and Windows (WSL)
>
> Install Cursor CLI with a single command:
>
> ```bash
> curl https://cursor.com/install -fsS | bash
> ```
>
> ### Windows (native)
>
> Install Cursor CLI on Windows using PowerShell:
>
> ```powershell
> irm 'https://cursor.com/install?win32=true' | iex
> ```
>
> ### Verification
>
> After installation, verify that Cursor CLI is working correctly:
>
> ```bash
> agent --version
> ```
>
> ## Post-installation setup
>
> 1. **Add \~/.local/bin to your PATH:**
>
>    For bash:
>
>    ```bash
>    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
>    source ~/.bashrc
>    ```
>
>    For zsh:
>
>    ```bash
>    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
>    source ~/.zshrc
>    ```
>
> 2. **Start using Cursor Agent:**
>    ```bash
>    agent
>    ```
>
> ## Updates
>
> Cursor CLI will try to auto-update by default to ensure you always have the latest version.
>
> To manually update Cursor CLI to the latest version:
>
> ```bash
> agent update
> ```
>
>

### Source: CLI overview

> # Cursor CLI
>
> Cursor CLI lets you interact with AI agents directly from your terminal to write, review, and modify code. Whether you prefer an interactive terminal interface or print automation for scripts and CI pipelines, the CLI provides powerful coding assistance right where you work.
>
> ## Getting started
>
> ```bash
> # Install (macOS, Linux, WSL)
> curl https://cursor.com/install -fsS | bash
>
> # Install (Windows PowerShell)
> irm 'https://cursor.com/install?win32=true' | iex
>
> # Run interactive session
> agent
> ```
>
> [Media](https://ptht05hbb1ssoooe.public.blob.vercel-storage.com/assets/uploads/plan-mode.mp4)
>
> ## Interactive mode
>
> Start a conversational session with the agent to describe your goals, review proposed changes, and approve commands:
>
> ```bash
> # Start interactive session
> agent
>
> # Start with initial prompt
> agent "refactor the auth module to use JWT tokens"
> ```
>
> ## Modes
>
> The CLI supports the same modes as the editor. Switch between modes using slash commands, keyboard shortcuts, or the `--mode` flag.
>
> | Mode      | Description                                                  | Shortcut                                    |
> | :-------- | :----------------------------------------------------------- | :------------------------------------------ |
> | **Agent** | Full access to all tools for complex coding tasks            | Default (no `--mode` value needed)          |
> | **Plan**  | Design your approach before coding with clarifying questions | Shift+Tab, `/plan`, `--plan`, `--mode=plan` |
> | **Ask**   | Read-only exploration without making changes                 | `/ask`, `--mode=ask`                        |
>
> ## Non-interactive mode
>
> Use print mode for non-interactive scenarios like scripts, CI pipelines, or automation:
>
> ```bash
> # Run with specific prompt and model
> agent -p "find and fix performance issues" --model "gpt-5"
>
> # Use with git changes included for review
> agent -p "review these changes for security issues" --output-format text
> ```
>
> ## Cloud Agent handoff
>
> Push your conversation to a [Cloud Agent](https://cursor.com/docs/cloud-agent.md) to continue running while you're away. Prepend `&` to any message:
>
> ```bash
> # Send a task to Cloud Agent mid-conversation
> & refactor the auth module and add comprehensive tests
> ```
>
> Pick up your Cloud Agent tasks on web or mobile at [cursor.com/agents](https://cursor.com/agents).
>
> ## Sessions
>
> Resume previous conversations to maintain context across multiple interactions:
>
> ```bash
> # Open previous chats and resume one
> agent ls
>
> # Resume latest conversation
> agent resume
>
> # Continue the previous session
> agent --continue
>
> # Resume specific conversation
> agent --resume="chat-id-here"
> ```
>
> ## Sandbox controls
>
> Configure command execution settings with `/sandbox` or the `--sandbox <mode>` flag (`enabled` or `disabled`). Toggle sandbox mode on or off and control network access through an interactive menu. Settings persist across sessions.
>
> [Media](https://ptht05hbb1ssoooe.public.blob.vercel-storage.com/assets/uploads/sandox.mp4)
>
> ## Sudo password prompting
>
> Run commands requiring elevated privileges without leaving the CLI. When a command needs `sudo`, Cursor displays a secure, masked password prompt. Your password flows directly to `sudo` via a secure IPC channel; the AI model never sees it.
>
>

### Source: Using the CLI

> # Using Agent in CLI
>
> ## Modes
>
> The CLI supports the same [modes](https://cursor.com/docs/agent/overview.md) as the editor. Switch modes using slash commands or the `--mode` flag.
>
> ### Plan mode
>
> Use Plan mode to design your approach before coding. The agent asks clarifying questions to refine your plan.
>
> - Press Shift+Tab to rotate to Plan mode
> - Use `/plan` to switch to Plan mode
> - Start with `--plan` or `--mode=plan` flag
>
> ### Ask mode
>
> Use Ask mode to explore code without making changes. The agent searches your codebase and provides answers without editing files.
>
> - Use `/ask` to switch to Ask mode
> - Start with `--mode=ask` flag
>
> ## Prompting
>
> Stating intent clearly is recommended for the best results. For example, you can use the prompt "do not write any code" to ensure that the agent won't edit any files. This is generally helpful when planning tasks before implementing them.
>
> Agent has tools for file operations, searching, running shell commands, and web access.
>
> ## MCP
>
> Agent supports [MCP (Model Context Protocol)](/marketplace) for extended functionality and integrations. The CLI will automatically detect and respect your `mcp.json` configuration file, enabling the same MCP servers and tools that you've configured for the editor.
>
> ## ACP
>
> Agent also supports [ACP (Agent Client Protocol)](https://cursor.com/docs/cli/acp.md) for custom client integrations. Use `agent acp` to run Cursor CLI as an ACP server over `stdio` with JSON-RPC messaging.
>
> ## Rules
>
> The CLI agent supports the same [rules system](https://cursor.com/docs/rules.md) as the editor. You can create rules in the `.cursor/rules` directory to provide context and guidance to the agent. These rules will be automatically loaded and applied based on their configuration, allowing you to customize the agent's behavior for different parts of your project or specific file types.
>
> The CLI also reads `AGENTS.md` and `CLAUDE.md` at the project root (if
> present) and applies them as rules alongside `.cursor/rules`.
>
> ## Working with Agent
>
> ### Navigation
>
> Previous messages can be accessed using arrow up (ArrowUp) where you can cycle through them.
>
> ### Input shortcuts
>
> - Shift+Tab — Rotate between modes (Agent, Plan, Ask)
> - Shift+Enter — Insert a newline instead of submitting, making it easier to write multi-line prompts.
> - Ctrl+D — Exit the CLI. Follows standard shell behavior, requiring a double-press to exit.
> - Ctrl+J or +Enter — Universal alternatives for inserting newlines that work in all terminals.
>
> Shift+Enter works in iTerm2, Ghostty, Kitty, Warp, and Zed. For tmux users, use Ctrl+J instead. See [Terminal setup](https://cursor.com/docs/cli/reference/terminal-setup.md) for configuration options and troubleshooting.
>
> ### Review
>
> Review changes with Ctrl+R. Press i to add follow-up instructions. Use ArrowUp/ArrowDown to scroll, and ArrowLeft/ArrowRight to switch files.
>
> ### Selecting context
>
> Select files and folders to include in context with @. Free up space in the context window by running `/summarize`. `/compress` remains an alias.
>
> ### Custom Modes
>
> Pick a [skill](https://cursor.com/docs/skills.md) from the `/` menu and press Enter to attach it to one message. Press Option+Enter instead to invoke it as a Custom Mode that stays active until you exit it.
>
> ## Cloud Agent handoff
>
> Push your conversation to a [Cloud Agent](https://cursor.com/docs/cloud-agent.md) and let it keep running while you're away. Prepend `&` to any message to send it to the cloud. Pick it back up on web or mobile at [cursor.com/agents](https://cursor.com/agents).
>
> ```bash
> # Send a task to Cloud Agent mid-conversation
> & refactor the auth module and add comprehensive tests
> ```
>
> ## CLI worktrees
>
> Pass `-w` or `--worktree [name]` to run the agent in a new Git worktree instead of editing your current checkout directly. Cursor creates these checkouts under `~/.cursor/worktrees/<reponame>/<name>`, alongside worktrees created from the editor. If you omit `name`, Cursor generates one.
>
> Cursor cleans up CLI worktrees with the same retention rules it uses for editor worktrees. For cleanup settings and limits, see [How are old worktrees cleaned up?](https://cursor.com/docs/configuration/worktrees.md#worktrees-cleanup).
>
> Combine `--workspace <path>` when you need an explicit repository root. Otherwise the CLI uses the current working directory. `--worktree` only changes where the agent makes file edits inside that project.
>
> ```bash
> # Create a temporary worktree from the current repository with a generated name
> agent --worktree "upgrade the test runner and fix any broken snapshots"
>
> # Create a named worktree from another repository
> agent --workspace ~/src/my-app --worktree auth-fix "fix the flaky auth test and open a PR"
> ```
>
> ## History
>
> Continue from an existing thread with `--resume [thread id]` to load prior context.
>
> To resume the most recent conversation, use `agent resume`, `--continue`, or the `/resume` slash command.
>
> You can also run `agent ls` to open previous chats and resume one.
>
> ## Command approval
>
> Before running terminal commands, CLI will ask you to approve (y) or reject (n) execution.
>
> ## Non-interactive mode
>
> Use `-p` or `--print` to run Agent in non-interactive mode. This will print the response to the console.
>
> With non-interactive mode, you can invoke Agent in a non-interactive way. This allows you to integrate it in scripts, CI pipelines, etc.
>
> You can combine this with `--output-format` to control how the output is formatted. For example, use `--output-format json` for structured output that's easier to parse in scripts, or `--output-format text` for plain text output of the agent's final response.
>
> Cursor has full write access in non-interactive mode.
>
>
