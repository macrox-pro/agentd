---
primary_sources:
  - id: T2-SDK
    title: "Codex SDK"
    url: "https://learn.chatgpt.com/docs/codex-sdk.md"
    section: "full page"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Codex SDK

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex SDK — overview; TypeScript; Python; Sandbox presets

> If you use Codex through Codex CLI, the IDE extension, or Codex cloud, you can also control it programmatically.
>
> Use the SDK when you need to:
>
> - Control Codex as part of your CI/CD pipeline
> - Create your own agent that can engage with Codex to perform complex engineering tasks
> - Build Codex into your own internal tools and workflows
> - Integrate Codex within your own application
>
> Use the Codex SDK for coding-focused Codex threads. If Codex is one specialist inside a broader orchestrated workflow, [run Codex CLI as an MCP server and orchestrate it with the Agents SDK](https://learn.chatgpt.com/docs/mcp-server).

> ## TypeScript library
>
> The TypeScript library lets your application start, continue, and resume local Codex threads.
>
> Use the library server-side; it requires Node.js 18 or later.
>
> ### Installation
>
> To get started, install the Codex SDK using `npm`:
>
> ```bash
> npm install @openai/codex-sdk
> ```
>
> For more details, check out the [TypeScript repo](https://github.com/openai/codex/tree/main/sdk/typescript).

> ## Python library
>
> The Python SDK controls the local Codex app-server over JSON-RPC. It requires Python 3.10 or later. Published SDK builds include a pinned Codex CLI runtime dependency.
>
> ### Installation
>
> To install the SDK run:
>
> ```bash
> pip install openai-codex
> ```
>
> Published SDK builds automatically use their pinned runtime. Pass `CodexConfig(codex_bin=...)` only when you intentionally want to run against a specific local Codex executable.

> ### Sandbox presets
>
> Use the same `Sandbox` presets when creating a thread or changing its filesystem
> access for a later turn:
>
> ```python
> from openai_codex import Codex, Sandbox
>
> with Codex() as codex:
>     thread = codex.thread_start(sandbox=Sandbox.workspace_write)
>     thread.run("Make the requested change.")
>     review = thread.run("Review the diff only.", sandbox=Sandbox.read_only)
> ```
>
> Available presets:
>
> - `Sandbox.read_only`: Read files without allowing writes.
> - `Sandbox.workspace_write`: Read files and write inside the workspace and configured writable roots.
> - `Sandbox.full_access`: Run without filesystem access restrictions.
>
> When you omit `sandbox=`, app-server uses its configured default. A sandbox
> passed to `run(...)` or `turn(...)` applies to that turn and later turns
> on the thread.
>
> For more details, check out the [Python repo](https://github.com/openai/codex/tree/main/sdk/python).
