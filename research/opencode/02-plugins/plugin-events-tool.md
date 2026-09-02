---
primary_sources:
  - id: T1-PLUGINS
    title: "Plugins"
    url: "https://opencode.ai/docs/plugins.md"
    section: "Tool Events"
  - id: T1-TOOLS
    title: "Tools"
    url: "https://opencode.ai/docs/tools.md"
    section: "apply_patch"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Plugin events — tool

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Plugins — Tool Events

> #### Tool Events
>
> - `tool.execute.after`
> - `tool.execute.before`

### Source: OpenCode Tools — apply_patch hook notes

> The `write` tool is controlled by the `edit` permission, which covers all file modifications (`edit`, `write`, `apply_patch`).
> :::
>
> ---
>
> ### read
>
> Read file contents from your codebase.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "read": "allow"
>   }
> }
> ```
>
> This tool reads files and returns their contents. It supports reading specific line ranges for large files.
>
> ---
>
> ### grep
>
> Search file contents using regular expressions.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "grep": "allow"
>   }
> }
> ```
>
> Fast content search across your codebase. Supports full regex syntax and file pattern filtering.
>
> ---
>
> ### glob
>
> Find files by pattern matching.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "glob": "allow"
>   }
> }
> ```
>
> Search for files using glob patterns like `**/*.js` or `src/**/*.ts`. Returns matching file paths sorted by modification time.
>
> ---
>
> ### lsp (experimental)
>
> Interact with your configured LSP servers to get code intelligence features like definitions, references, hover info, and call hierarchy.
>
> :::note
> This tool is only available when `OPENCODE_EXPERIMENTAL_LSP_TOOL=true` (or `OPENCODE_EXPERIMENTAL=true`).
> :::
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "lsp": "allow"
>   }
> }
> ```
>
> Supported operations include `goToDefinition`, `findReferences`, `hover`, `documentSymbol`, `workspaceSymbol`, `goToImplementation`, `prepareCallHierarchy`, `incomingCalls`, and `outgoingCalls`.
>
> To configure which LSP servers are available for your project, see [LSP Servers](/docs/lsp).
>
> ---
>
> ### apply_patch
>
> Apply patches to files.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "edit": "allow"
>   }
> }
> ```
>
> This tool applies patch files to your codebase. Useful for applying diffs and patches from various sources.
>
> When handling `tool.execute.before` or `tool.execute.after` hooks, check `input.tool === "apply_patch"` (not `"patch"`).
>
> `apply_patch` uses `output.args.patchText` instead of `output.args.filePath`. Paths are embedded in marker lines within `patchText` and are relative to the project root (for example: `*** Add File: src/new-file.ts`, `*** Update File: src/existing.ts`, `*** Move to: src/renamed.ts`, `*** Delete File: src/obsolete.ts`).
>
> :::note
> The `apply_patch` tool is controlled by the `edit` permission, which covers all file modifications (`edit`, `write`, `apply_patch`).
> :::
>
> ---
>
> ### skill
>
> Load a [skill](/docs/skills) (a `SKILL.md` file) and return its content in the conversation.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "skill": "allow"
>   }
> }
> ```
>
> ---
>
> ### todowrite
>
> Manage todo lists during coding sessions.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "todowrite": "allow"
>   }
> }
> ```
>
> Creates and updates task lists to track progress during complex operations. The LLM uses this to organize multi-step tasks.
>
> :::note
> This tool is disabled for subagents by default, but you can enable it manually. [Learn more](/docs/agents/#permissions)
> :::
>
> ---
>
> ### webfetch
>
> Fetch web content.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "webfetch": "allow"
>   }
> }
> ```
>
> Allows the LLM to fetch and read web pages. Useful for looking up documentation or researching online resources.
>
> ---
>
> ### websearch
>
> Search the web for information.
>
> :::note
> This tool is only available when using the OpenCode or OpenCode Go provider, or when either the `OPENCODE_ENABLE_EXA` or `OPENCODE_ENABLE_PARALLEL` environment variable is set to any truthy value (e.g., `true` or `1`).
>
> To enable when launching OpenCode:
>
> ```bash
> OPENCODE_ENABLE_EXA=1 opencode
> # or
> OPENCODE_ENABLE_PARALLEL=1 opencode
> ```
>
> :::
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "websearch": "allow"
>   }
> }
> ```
>
> Performs web searches using Exa or Parallel to find relevant information online. Useful for researching topics, finding current events, or gathering information beyond the training data cutoff.
>
> No API key is required — the tool connects directly to the backend's hosted MCP service without authentication.
>
> :::tip
> Use `websearch` when you need to find information (discovery), and `webfetch` when you need to retrieve content from a specific URL (retrieval).
> :::
>
> ---
>
> ### question
>
> Ask the user questions during execution.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "question": "allow"
>   }
> }
> ```
>
> This tool allows the LLM to ask the user questions during a task. It's useful for:
>
> - Gathering user preferences or requirements
> - Clarifying ambiguous instructions
> - Getting decisions on implementation choices
> - Offering choices about what direction to take
>
> Each question includes a header, the question text, and a list of options. Users can select from the provided options or type a custom answer. When there are multiple questions, users can navigate between them before submitting all answers.
>
> ---
>
>
> ### apply_patch
>
> Apply patches to files.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "edit": "allow"
>   }
> }
> ```
>
> This tool applies patch files to your codebase. Useful for applying diffs and patches from various sources.
>
> When handling `tool.execute.before` or `tool.execute.after` hooks, check `input.tool === "apply_patch"` (not `"patch"`).
>
> `apply_patch` uses `output.args.patchText` instead of `output.args.filePath`. Paths are embedded in marker lines within `patchText` and are relative to the project root (for example: `*** Add File: src/new-file.ts`, `*** Update File: src/existing.ts`, `*** Move to: src/renamed.ts`, `*** Delete File: src/obsolete.ts`).
>
> :::note
> The `apply_patch` tool is controlled by the `edit` permission, which covers all file modifications (`edit`, `write`, `apply_patch`).
> :::
>
> ---
>
> ### skill
>
> Load a [skill](/docs/skills) (a `SKILL.md` file) and return its content in the conversation.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "skill": "allow"
>   }
> }
> ```
>
> ---
>
> ### todowrite
>
> Manage todo lists during coding sessions.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "todowrite": "allow"
>   }
> }
> ```
>
> Creates and updates task lists to track progress during complex operations. The LLM uses this to organize multi-step tasks.
>
> :::note
> This tool is disabled for subagents by default, but you can enable it manually. [Learn more](/docs/agents/#permissions)
> :::
>
> ---
>
> ### webfetch
>
> Fetch web content.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "webfetch": "allow"
>   }
> }
> ```
>
> Allows the LLM to fetch and read web pages. Useful for looking up documentation or researching online resources.
>
> ---
>
> ### websearch
>
> Search the web for information.
>
> :::note
> This tool is only available when using the OpenCode or OpenCode Go provider, or when either the `OPENCODE_ENABLE_EXA` or `OPENCODE_ENABLE_PARALLEL` environment variable is set to any truthy value (e.g., `true` or `1`).
>
> To enable when launching OpenCode:
>
> ```bash
> OPENCODE_ENABLE_EXA=1 opencode
> # or
> OPENCODE_ENABLE_PARALLEL=1 opencode
> ```
>
> :::
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "websearch": "allow"
>   }
> }
> ```
>
> Performs web searches using Exa or Parallel to find relevant information online. Useful for researching topics, finding current events, or gathering information beyond the training data cutoff.
>
> No API key is required — the tool connects directly to the backend's hosted MCP service without authentication.
>
> :::tip
> Use `websearch` when you need to find information (discovery), and `webfetch` when you need to retrieve content from a specific URL (retrieval).
> :::
>
> ---
>
> ### question
>
> Ask the user questions during execution.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "question": "allow"
>   }
> }
> ```
>
> This tool allows the LLM to ask the user questions during a task. It's useful for:
>
> - Gathering user preferences or requirements
> - Clarifying ambiguous instructions
> - Getting decisions on implementation choices
> - Offering choices about what direction to take
>
> Each question includes a header, the question text, and a list of options. Users can select from the provided options or type a custom answer. When there are multiple questions, users can navigate between them before submitting all answers.
>
> ---
>
>
> When handling `tool.execute.before` or `tool.execute.after` hooks, check `input.tool === "apply_patch"` (not `"patch"`).
>
> `apply_patch` uses `output.args.patchText` instead of `output.args.filePath`. Paths are embedded in marker lines within `patchText` and are relative to the project root (for example: `*** Add File: src/new-file.ts`, `*** Update File: src/existing.ts`, `*** Move to: src/renamed.ts`, `*** Delete File: src/obsolete.ts`).
>
> :::note
> The `apply_patch` tool is controlled by the `edit` permission, which covers all file modifications (`edit`, `write`, `apply_patch`).
> :::
>
> ---
>
> ### skill
>
> Load a [skill](/docs/skills) (a `SKILL.md` file) and return its content in the conversation.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "skill": "allow"
>   }
> }
> ```
>
> ---
>
> ### todowrite
>
> Manage todo lists during coding sessions.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "todowrite": "allow"
>   }
> }
> ```
>
> Creates and updates task lists to track progress during complex operations. The LLM uses this to organize multi-step tasks.
>
> :::note
> This tool is disabled for subagents by default, but you can enable it manually. [Learn more](/docs/agents/#permissions)
> :::
>
> ---
>
> ### webfetch
>
> Fetch web content.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "webfetch": "allow"
>   }
> }
> ```
>
> Allows the LLM to fetch and read web pages. Useful for looking up documentation or researching online resources.
>
> ---
>
> ### websearch
>
> Search the web for information.
>
> :::note
> This tool is only available when using the OpenCode or OpenCode Go provider, or when either the `OPENCODE_ENABLE_EXA` or `OPENCODE_ENABLE_PARALLEL` environment variable is set to any truthy value (e.g., `true` or `1`).
>
> To enable when launching OpenCode:
>
> ```bash
> OPENCODE_ENABLE_EXA=1 opencode
> # or
> OPENCODE_ENABLE_PARALLEL=1 opencode
> ```
>
> :::
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "websearch": "allow"
>   }
> }
> ```
>
> Performs web searches using Exa or Parallel to find relevant information online. Useful for researching topics, finding current events, or gathering information beyond the training data cutoff.
>
> No API key is required — the tool connects directly to the backend's hosted MCP service without authentication.
>
> :::tip
> Use `websearch` when you need to find information (discovery), and `webfetch` when you need to retrieve content from a specific URL (retrieval).
> :::
>
> ---
>
> ### question
>
> Ask the user questions during execution.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "question": "allow"
>   }
> }
> ```
>
> This tool allows the LLM to ask the user questions during a task. It's useful for:
>
> - Gathering user preferences or requirements
> - Clarifying ambiguous instructions
> - Getting decisions on implementation choices
> - Offering choices about what direction to take
>
> Each question includes a header, the question text, and a list of options. Users can select from the provided options or type a custom answer. When there are multiple questions, users can navigate between them before submitting all answers.
>
> ---
>
>
> `apply_patch` uses `output.args.patchText` instead of `output.args.filePath`. Paths are embedded in marker lines within `patchText` and are relative to the project root (for example: `*** Add File: src/new-file.ts`, `*** Update File: src/existing.ts`, `*** Move to: src/renamed.ts`, `*** Delete File: src/obsolete.ts`).
>
> :::note
> The `apply_patch` tool is controlled by the `edit` permission, which covers all file modifications (`edit`, `write`, `apply_patch`).
> :::
>
> ---
>
> ### skill
>
> Load a [skill](/docs/skills) (a `SKILL.md` file) and return its content in the conversation.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "skill": "allow"
>   }
> }
> ```
>
> ---
>
> ### todowrite
>
> Manage todo lists during coding sessions.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "todowrite": "allow"
>   }
> }
> ```
>
> Creates and updates task lists to track progress during complex operations. The LLM uses this to organize multi-step tasks.
>
> :::note
> This tool is disabled for subagents by default, but you can enable it manually. [Learn more](/docs/agents/#permissions)
> :::
>
> ---
>
> ### webfetch
>
> Fetch web content.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "webfetch": "allow"
>   }
> }
> ```
>
> Allows the LLM to fetch and read web pages. Useful for looking up documentation or researching online resources.
>
> ---
>
> ### websearch
>
> Search the web for information.
>
> :::note
> This tool is only available when using the OpenCode or OpenCode Go provider, or when either the `OPENCODE_ENABLE_EXA` or `OPENCODE_ENABLE_PARALLEL` environment variable is set to any truthy value (e.g., `true` or `1`).
>
> To enable when launching OpenCode:
>
> ```bash
> OPENCODE_ENABLE_EXA=1 opencode
> # or
> OPENCODE_ENABLE_PARALLEL=1 opencode
> ```
>
> :::
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "websearch": "allow"
>   }
> }
> ```
>
> Performs web searches using Exa or Parallel to find relevant information online. Useful for researching topics, finding current events, or gathering information beyond the training data cutoff.
>
> No API key is required — the tool connects directly to the backend's hosted MCP service without authentication.
>
> :::tip
> Use `websearch` when you need to find information (discovery), and `webfetch` when you need to retrieve content from a specific URL (retrieval).
> :::
>
> ---
>
> ### question
>
> Ask the user questions during execution.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "question": "allow"
>   }
> }
> ```
>
> This tool allows the LLM to ask the user questions during a task. It's useful for:
>
> - Gathering user preferences or requirements
> - Clarifying ambiguous instructions
> - Getting decisions on implementation choices
> - Offering choices about what direction to take
>
> Each question includes a header, the question text, and a list of options. Users can select from the provided options or type a custom answer. When there are multiple questions, users can navigate between them before submitting all answers.
>
> ---
>
>
> The `apply_patch` tool is controlled by the `edit` permission, which covers all file modifications (`edit`, `write`, `apply_patch`).
> :::
>
> ---
>
> ### skill
>
> Load a [skill](/docs/skills) (a `SKILL.md` file) and return its content in the conversation.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "skill": "allow"
>   }
> }
> ```
>
> ---
>
> ### todowrite
>
> Manage todo lists during coding sessions.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "todowrite": "allow"
>   }
> }
> ```
>
> Creates and updates task lists to track progress during complex operations. The LLM uses this to organize multi-step tasks.
>
> :::note
> This tool is disabled for subagents by default, but you can enable it manually. [Learn more](/docs/agents/#permissions)
> :::
>
> ---
>
> ### webfetch
>
> Fetch web content.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "webfetch": "allow"
>   }
> }
> ```
>
> Allows the LLM to fetch and read web pages. Useful for looking up documentation or researching online resources.
>
> ---
>
> ### websearch
>
> Search the web for information.
>
> :::note
> This tool is only available when using the OpenCode or OpenCode Go provider, or when either the `OPENCODE_ENABLE_EXA` or `OPENCODE_ENABLE_PARALLEL` environment variable is set to any truthy value (e.g., `true` or `1`).
>
> To enable when launching OpenCode:
>
> ```bash
> OPENCODE_ENABLE_EXA=1 opencode
> # or
> OPENCODE_ENABLE_PARALLEL=1 opencode
> ```
>
> :::
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "websearch": "allow"
>   }
> }
> ```
>
> Performs web searches using Exa or Parallel to find relevant information online. Useful for researching topics, finding current events, or gathering information beyond the training data cutoff.
>
> No API key is required — the tool connects directly to the backend's hosted MCP service without authentication.
>
> :::tip
> Use `websearch` when you need to find information (discovery), and `webfetch` when you need to retrieve content from a specific URL (retrieval).
> :::
>
> ---
>
> ### question
>
> Ask the user questions during execution.
>
> ```json title="opencode.json" {4}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "permission": {
>     "question": "allow"
>   }
> }
> ```
>
> This tool allows the LLM to ask the user questions during a task. It's useful for:
>
> - Gathering user preferences or requirements
> - Clarifying ambiguous instructions
> - Getting decisions on implementation choices
> - Offering choices about what direction to take
>
> Each question includes a header, the question text, and a list of options. Users can select from the provided options or type a custom answer. When there are multiple questions, users can navigate between them before submitting all answers.
>
> ---
