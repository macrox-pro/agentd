---
primary_sources:
  - id: T1-SUBAGENTS
    title: "Full page"
    url: "https://cursor.com/docs/subagents.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Subagents

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Subagents

> # Subagents
>
> Subagents are specialized AI assistants that Cursor's agent can delegate tasks to. Each subagent operates in its own context window, handles specific types of work, and returns its result to the parent agent. Use subagents to break down complex tasks, do work in parallel, and preserve context in the main conversation.
>
> You can use subagents in the editor, CLI, and [Cloud Agents](https://cursor.com/docs/cloud-agent.md).
>
> ### Context isolation
>
> Each subagent has its own context window. Long research or exploration tasks don't consume space in your main conversation.
>
> ### Parallel execution
>
> Launch multiple subagents simultaneously. Work on different parts of your codebase without waiting for sequential completion.
>
> ### Specialized expertise
>
> Configure subagents with custom prompts, tool access, and models for domain-specific tasks.
>
> ### Reusability
>
> Define custom subagents and use them across projects.
>
> ## How subagents work
>
> When Agent encounters a complex task, it can launch a subagent automatically. The subagent receives a prompt with all necessary context, works autonomously, and returns a final message with its results.
>
> Subagents start with a clean context. The parent agent includes relevant information in the prompt since subagents don't have access to prior conversation history.
>
> ### Foreground vs background
>
> Subagents run in one of two modes:
>
> | Mode           | Behavior                                                             | Best for                                    |
> | :------------- | :------------------------------------------------------------------- | :------------------------------------------ |
> | **Foreground** | Blocks until the subagent completes. Returns the result immediately. | Sequential tasks where you need the output. |
> | **Background** | Returns immediately. The subagent works independently.               | Long-running tasks or parallel workstreams. |
>
> ## Built-in subagents
>
> Cursor includes three built-in subagents that handle context-heavy operations automatically. These subagents were designed based on analysis of agent conversations where context window limits were hit.
>
> | Subagent    | Purpose                         | Why it's a subagent                                                                                                                            |
> | :---------- | :------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------------------- |
> | **Explore** | Searches and analyzes codebases | Codebase exploration generates large intermediate output that would bloat the main context. Uses a faster model to run many parallel searches. |
> | **Bash**    | Runs series of shell commands   | Command output is often verbose. Isolating it keeps the parent focused on decisions, not logs.                                                 |
> | **Browser** | Controls browser via MCP tools  | Browser interactions produce noisy DOM snapshots and screenshots. The subagent filters this down to relevant results.                          |
>
> ### Why these subagents exist
>
> These three operations share common traits: they generate noisy intermediate output, benefit from specialized prompts and tools, and can consume significant context. Running them as subagents solves several problems:
>
> - **Context isolation** — Intermediate output stays in the subagent. The parent only sees the final summary.
> - **Model flexibility** — The explore subagent uses a faster model by default. This enables running 10 parallel searches in the time a single main-agent search would take.
> - **Specialized configuration** — Each subagent has prompts and tool access tuned for its specific task.
> - **Cost efficiency** — Faster models cost less. Isolating token-heavy work in subagents with appropriate model choices reduces overall cost.
>
> You don't need to configure these subagents. Agent uses them automatically when appropriate.
>
> ## When to use subagents
>
> | Use subagents when...                                     | Use skills when...                                      |
> | :-------------------------------------------------------- | :------------------------------------------------------ |
> | You need context isolation for long research tasks        | The task is single-purpose (generate changelog, format) |
> | Running multiple workstreams in parallel                  | You want a quick, repeatable action                     |
> | The task requires specialized expertise across many steps | The task completes in one shot                          |
> | You want an independent verification of work              | You don't need a separate context window                |
>
> If you find yourself creating a subagent for a simple, single-purpose task like "generate a changelog" or "format imports," consider using a [skill](https://cursor.com/docs/skills.md) instead.
>
> ## Quick start
>
> Agent automatically uses subagents when appropriate. You can also create a custom subagent by asking Agent:
>
> Create a subagent file at .cursor/agents/verifier.md with YAML frontmatter (name, description) followed by the prompt. The verifier subagent should validate completed work, check that implementations are functional, run tests, and report what passed vs what's incomplete.
>
> For more control, create custom subagents manually in your project or user directory.
>
> ## Custom subagents
>
> Define custom subagents to encode specialized knowledge, enforce team standards, or automate repetitive workflows.
>
> ### File locations
>
> | Type                  | Location            | Scope                                                |
> | :-------------------- | :------------------ | :--------------------------------------------------- |
> | **Project subagents** | `.cursor/agents/`   | Current project only                                 |
> |                       | `.claude/agents/`   | Current project only (Claude compatibility)          |
> |                       | `.codex/agents/`    | Current project only (Codex compatibility)           |
> | **User subagents**    | `~/.cursor/agents/` | All projects for current user                        |
> |                       | `~/.claude/agents/` | All projects for current user (Claude compatibility) |
> |                       | `~/.codex/agents/`  | All projects for current user (Codex compatibility)  |
>
> Project subagents take precedence when names conflict. When multiple locations contain subagents with the same name, `.cursor/` takes precedence over `.claude/` or `.codex/`.
>
> ### File format
>
> Each subagent is a markdown file with YAML frontmatter:
>
> ```markdown
