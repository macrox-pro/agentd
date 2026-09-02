---
primary_sources:
  - id: T1-AGENTS
    title: "Agents"
    url: "https://opencode.ai/docs/agents.md"
    section: "Built-in"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Built-in agents and modes

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Agents — Types and Built-in

> ## Types
>
> There are two types of agents in OpenCode; primary agents and subagents.
>
> ---
>
> ### Primary agents
>
> Primary agents are the main assistants you interact with directly. You can cycle through them using the **Tab** key, or your configured `switch_agent` keybind. These agents handle your main conversation. Tool access is configured via permissions — for example, Build has all tools enabled while Plan is restricted.
>
> :::tip
> You can use the **Tab** key to switch between primary agents during a session.
> :::
>
> OpenCode comes with two built-in primary agents, **Build** and **Plan**. We'll
> look at these below.
>
> ---
>
> ### Subagents
>
> Subagents are specialized assistants that primary agents can invoke for specific tasks. You can also manually invoke them by **@ mentioning** them in your messages.
>
> OpenCode comes with three built-in subagents, **General**, **Explore**, and **Scout**. We'll look at this below.
>
> ---
>
> ---
>
> ## Built-in
>
> OpenCode comes with two built-in primary agents and three built-in subagents.
>
> ---
>
> ### Use build
>
> _Mode_: `primary`
>
> Build is the **default** primary agent with all tools enabled. This is the standard agent for development work where you need full access to file operations and system commands.
>
> ---
>
> ### Use plan
>
> _Mode_: `primary`
>
> A restricted agent designed for planning and analysis. We use a permission system to give you more control and prevent unintended changes.
> By default, all of the following are set to `ask`:
>
> - `file edits`: All writes, patches, and edits
> - `bash`: All bash commands
>
> This agent is useful when you want the LLM to analyze code, suggest changes, or create plans without making any actual modifications to your codebase.
>
> ---
>
> ### Use general
>
> _Mode_: `subagent`
>
> A general-purpose agent for researching complex questions and executing multi-step tasks. Has full tool access (except todo), so it can make file changes when needed. Use this to run multiple units of work in parallel.
>
> ---
>
> ### Use explore
>
> _Mode_: `subagent`
>
> A fast, read-only agent for exploring codebases. Cannot modify files. Use this when you need to quickly find files by patterns, search code for keywords, or answer questions about the codebase.
>
> ---
>
> ### Use scout
>
> _Mode_: `subagent`
>
> A read-only agent for external docs and dependency research. Use this when you need to clone a dependency repository into OpenCode's managed cache, inspect library source, or cross-reference local code against upstream implementations without modifying your workspace.
>
> ---
>
> ### Use compaction
>
> _Mode_: `primary`
>
> Hidden system agent that compacts long context into a smaller summary. It runs automatically when needed and is not selectable in the UI.
>
> ---
>
> ### Use title
>
> _Mode_: `primary`
>
> Hidden system agent that generates short session titles. It runs automatically and is not selectable in the UI.
>
> ---
>
> ### Use summary
>
> _Mode_: `primary`
>
> Hidden system agent that creates session summaries. It runs automatically and is not selectable in the UI.
>
> ---
>
> ---
>
> ## Usage
>
> 1. For primary agents, use the **Tab** key to cycle through them during a session. You can also use your configured `switch_agent` keybind.
>
> 2. Subagents can be invoked:
>    - **Automatically** by primary agents for specialized tasks based on their descriptions.
>    - Manually by **@ mentioning** a subagent in your message. For example.
>
>      ```txt frame="none"
>      @general help me search for this function
>      ```
>
> 3. **Navigation between sessions**: When subagents create child sessions, use `session_child_first` (default: **\<Leader>+Down**) to enter the first child session from the parent.
>
> 4. Once you are in a child session, use:
>    - `session_child_cycle` (default: **Right**) to cycle to the next child session
>    - `session_child_cycle_reverse` (default: **Left**) to cycle to the previous child session
>    - `session_parent` (default: **Up**) to return to the parent session
>
>    This lets you switch between the main conversation and specialized subagent work.
>
> ---
