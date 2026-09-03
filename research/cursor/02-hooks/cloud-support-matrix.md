---
primary_sources:
  - id: T1-HOOKS
    title: "Cloud agent support"
    url: "https://cursor.com/docs/hooks.md"
    section: "Cloud agent support"
  - id: T1-CLOUD-SELF-POOL
    title: "Team Pools — Hooks"
    url: "https://cursor.com/docs/cloud-agent/self-hosted/pool.md"
    section: "Hooks"
also_cited_in: []
studied_at: "2026-09-03"
cursor_docs_snapshot: "2026-09-03"
applicability: "current"
---
# Cloud agent hook support

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-09-03).

### Source: Cursor Hooks — Cloud agent support

> ## Cloud agent support
>
> Cloud agents run command-based hooks from your repository. If you have hooks defined in `.cursor/hooks.json` at the root of your project, cloud agents pick them up and run them during their work.
>
> On Enterprise plans, cloud agents also run team hooks and enterprise-managed hooks configured through the [web dashboard](https://cursor.com/dashboard/team-content?section=hooks).
>
> Cloud agents sometimes begin in a read-only environment for early exploratory turns. Hooks do not run during those turns. They start once the agent has a writable environment.
>
> ### Supported hooks
>
> The following hooks run in cloud agents:
>
> | Hook                   | Supported |
> | ---------------------- | --------- |
> | `beforeShellExecution` | Yes       |
> | `afterShellExecution`  | Yes       |
> | `beforeReadFile`       | Yes       |
> | `afterFileEdit`        | Yes       |
> | `preToolUse`           | Yes       |
> | `postToolUse`          | Yes       |
> | `postToolUseFailure`   | Yes       |
> | `subagentStart`        | Yes       |
> | `subagentStop`         | Yes       |
> | `beforeSubmitPrompt`   | Yes       |
> | `preCompact`           | Yes       |
> | `afterAgentResponse`   | Yes       |
> | `afterAgentThought`    | Yes       |
> | `stop`                 | Yes       |
>
> ### Hooks not available in cloud agents
>
> Some hooks don't apply to cloud agents due to differences in the execution environment:
>
> | Hook                                       | Reason                                                                                                                                                                                                   |
> | ------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `sessionStart`                             | Deferred while cloud agents can still start in a read-only environment. Hooks don't load there, so a cloud `sessionStart` would fire too late (after the first write) rather than at true session start. |
> | `sessionEnd`                               | Cloud agents have no editor-lifetime session boundary. `sessionEnd` is tied to the IDE session, not a cloud agent chat.                                                                                  |
> | `beforeMCPExecution` / `afterMCPExecution` | Deferred while cloud agents can still start in a read-only environment, where hooks don't load and MCP hook timing is unclear.                                                                           |
> | `beforeTabFileRead` / `afterTabFileEdit`   | Tab completions are an IDE feature and don't run in cloud agents.                                                                                                                                        |
> | `workspaceOpen`                            | This is an IDE lifecycle hook and doesn't apply to cloud agents.                                                                                                                                         |
>
> ### Configuration sources
>
> Cloud agents load hooks from these sources:
>
> - **Project hooks** (`.cursor/hooks.json` in your repo): Loaded and run during cloud agent work.
> - **Team hooks** (Enterprise): Distributed from the dashboard and run in cloud agents.
> - **Enterprise hooks** (Enterprise): System-wide managed hooks run in cloud agents.
>
> User-level hooks (`~/.cursor/hooks.json`) are not available in cloud agents. Cloud agent VMs don't have access to your local home directory configuration.
>
> Self-Hosted Machines workers (Pools and My Machines) run the same command-based project hooks, and on Enterprise also run team and enterprise-managed hooks. On those workers, `sessionStart` and `sessionEnd` fire when a session claims the worker and when that claim is released. See [Hooks on Pools](https://cursor.com/docs/cloud-agent/self-hosted/pool.md#hooks).
>
> ### Execution type limits
>
> Cloud agents run **command-based hooks** only. Prompt-based hooks require authentication wiring between the hook and the agent loop, which isn't available in the cloud execution environment.

### Source: Team Pools — Hooks

> ## Hooks
>
> Self-Hosted Machines workers run command-based hooks during Cloud Agent sessions. They load configuration from the workspace the worker serves.
>
> - Project hooks. Commit `.cursor/hooks.json` and the scripts it references in the repository or workspace directory the worker uses. That is the git checkout, a `--worker-dir` root, or the directory you pass for an any-repo pool.
> - Team and enterprise-managed hooks. On Enterprise, workers also run hooks configured in the web dashboard.
>
> The Hooks reference covers the schema, events, and examples. Cloud agent support lists which events the Cloud Agent loop runs.
>
> The Cloud Agent limits that also apply on these tool-execution workers:
>
> - Command-based hooks only. Prompt-based hooks do not run.
> - No IDE-only hooks. Tab hooks (`beforeTabFileRead`, `afterTabFileEdit`) and `workspaceOpen` do not run on workers.
>
> `sessionStart` and `sessionEnd` run on Self-Hosted Machines workers. They fire when a Cloud Agent session claims the worker and when that claim is released. Cursor-managed Cloud Agents skip those hooks.
>
> Kubernetes workers use this same hooks model.
>
> ### Do hooks run on self-hosted workers?
>
> Yes. Self-Hosted Machines workers run project hooks from
> `.cursor/hooks.json`. On Enterprise, they also run team and
> enterprise-managed hooks. See Hooks.
