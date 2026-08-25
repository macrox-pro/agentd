---
primary_sources:
  - id: T2-IMPORT
    title: "Import from another agent"
    url: "https://learn.chatgpt.com/docs/import.md"
    section: "What ChatGPT can import; What to review after importing"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Import from another agent

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Import from another agent — overview; CLI import

> Use the import flow to bring instructions, settings, skills, plugins, projects,
> and recent work from another agent into the ChatGPT desktop app or Codex CLI.
> The desktop app can import from **Claude Code**, **Claude Cowork**,
> or **Cursor**. Codex CLI can import from **Claude Code** or **Cursor**.
>
> Importing doesn't change or delete your existing agent setup.

> ### Import in Codex CLI
>
> 1. Start a local Codex CLI session and type `/import`.
> 2. Choose **Claude Code** or **Cursor**.
> 3. Select the supported setup, project files, and recent chats you want to
>    import.
> 4. Review the imported configuration and continue working in Codex.
>
> Codex CLI imports up to 50 chats from the last 30 days. The `/import` command
> isn't available during a running task, in a remote session, or while connected
> to a local app-server daemon. See [CLI slash
> commands](https://learn.chatgpt.com/docs/developer-commands?surface=cli#cli-import-claude-code-or-cursor-setup-with-import).

### Source: Import from another agent — What ChatGPT can import

> ## What ChatGPT can import
>
> | Imported item                     | Destination                                             |
> | --------------------------------- | ------------------------------------------------------- |
> | Instruction files                 | [`AGENTS.md`](https://learn.chatgpt.com/docs/agent-configuration/agents-md)     |
> | `settings.json`                   | [`config.toml`](https://learn.chatgpt.com/docs/config-file/config-basic)        |
> | Skills                            | [Skills](https://learn.chatgpt.com/docs/build-skills)                           |
> | Plugins                           | Plugins                                                 |
> | Existing project folders          | Projects using the same folders                         |
> | Project memories from Claude Code | [Memories](https://learn.chatgpt.com/docs/customization/memories)               |
> | Chats from the last 30 days       | ChatGPT chats                                           |
> | MCP server configuration          | [Codex MCP configuration](https://learn.chatgpt.com/docs/extend/mcp)            |
> | Hooks                             | [Codex hooks](https://learn.chatgpt.com/docs/hooks)                             |
> | Slash commands                    | [Skills](https://learn.chatgpt.com/docs/build-skills)                           |
> | Subagents                         | [Codex subagents](https://learn.chatgpt.com/docs/agent-configuration/subagents) |

### Source: Import from another agent — What to review after importing

> ## What to review after importing
>
> Review imported setup before you rely on it, especially:
>
> - Tool restrictions or permissions in imported skills and agents.
> - MCP server settings that use custom authentication, headers, environment
>   variables, or transports. You may need to sign in again.
> - Hooks whose behavior may differ after import.
> - Plugins, marketplaces, or other setup that needs manual follow-up.
> - Prompt templates or command-style prompts that depend on arguments, shell
>   interpolation, or file-path placeholders.
