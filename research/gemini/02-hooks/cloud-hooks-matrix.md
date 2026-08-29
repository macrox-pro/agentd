---
primary_sources:
  - id: T1-HOOKS
    title: "Gemini CLI hooks"
    url: "https://geminicli.com/docs/hooks.md"
    section: "Configuration; precedence"
  - id: T2-AGENT-HOOKS
    title: "Managed Agents API hooks"
    url: "https://ai.google.dev/gemini-api/docs/agent-hooks"
    section: "Lifecycle events; configuration discovery"
  - id: T4-TRANSITION
    title: "Gemini CLI transition timeline"
    url: "https://geminicli.com/docs/transition-timeline.md"
    section: "Antigravity CLI feature parity (hooks)"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Cloud hooks matrix

> **Applicability:** Verbatim excerpts comparing Gemini CLI local hooks vs Managed Agents API sandbox hooks (snapshot 2026-08-29). Antigravity CLI uses a third wire (`PreToolUse`, `.agents/hooks.json`) — see [09-migration-and-antigravity/gcli-migration.md](../09-migration-and-antigravity/gcli-migration.md).

### Source: Gemini CLI hooks — Configuration

> Hooks are configured in `settings.json`. Gemini CLI merges configurations from
> multiple layers in the following order of precedence (highest to lowest):
>
> 1.  **Project settings**: `.gemini/settings.json` in the current directory.
> 2.  **User settings**: `~/.gemini/settings.json`.
> 3.  **System settings**: `/etc/gemini-cli/settings.json`.
> 4.  **Extensions**: Hooks defined by installed extensions.
>
> ### Configuration schema
>
> ```json
> {
>   "hooks": {
>     "BeforeTool": [
>       {
>         "matcher": "write_file|replace",
>         "hooks": [
>           {
>             "name": "security-check",
>             "type": "command",
>             "command": "$GEMINI_PROJECT_DIR/.gemini/hooks/security.sh",
>             "timeout": 5000
>           }
>         ]
>       }
>     ]
>   }
> }
> ```
>
> #### Hook configuration fields
>
> | Field         | Type   | Required  | Description                                                          |
> | :------------ | :----- | :-------- | :------------------------------------------------------------------- |
> | `type`        | string | **Yes**   | The execution engine. Currently only `"command"` is supported.       |
> | `command`     | string | **Yes\*** | The shell command to execute. (Required when `type` is `"command"`). |
> | `name`        | string | No        | A friendly name for identifying the hook in logs and CLI commands.   |
> | `timeout`     | number | No        | Execution timeout in milliseconds (default: 60000).                  |
> | `description` | string | No        | A brief explanation of the hook's purpose.                           |

### Source: Managed Agents API hooks — Supported lifecycle events

> Hooks support 2 events inside the sandbox:
>
> | Event | When it fires | What it does |
> | --- | --- | --- |
> | `pre_tool_execution` | Right before a tool runs | Can approve (`allow`) or block (`deny`) the tool before it executes. When blocked, the model sees your rejection reason and adapts. |
> | `post_tool_execution` | Right after a tool finishes | Runs follow-up tasks like formatting code, running unit tests, or logging telemetry. Cannot block or undo completed actions. |

### Source: Managed Agents API hooks — Configuration discovery

> The runtime automatically discovers hook definitions from`.agents/hooks.json` or`/.agents/hooks.json` inside the sandbox environment. You can provide`hooks.json` alongside your custom scripts using any supported environment source:
>
> - Repository mount: A Git repository containing`.agents/hooks.json` alongside`AGENTS.md`.
> - Cloud Storage (`gcs`): A GCS bucket containing`hooks.json` copied into the environment.
> - Inline sources: Raw JSON string and script contents passed in`environment.sources` when calling`client.interactions.create`.
>
> ### hooks.json schema
>
> A`hooks.json` file groups event definitions (`pre_tool_execution` or`post_tool_execution`) under custom names. You can enable or disable each group independently:
>
> ```
> {
>   "security-gate": {
>     "enabled": true,
>     "pre_tool_execution": [
>       {
>         "matcher": "code_execution",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "python3 /.agents/hooks-scripts/gate.py",
>             "timeout": 10
>           }
>         ]
>       }
>     ]
>   },
>   "auto-format": {
>     "post_tool_execution": [
>       {
>         "matcher": "*",
>         "hooks": [
>           {
>             "type": "command",
>             "command": "python3 /.agents/hooks-scripts/auto_lint.py",
>             "timeout": 15
>           }
>         ]
>       }
>     ]
>   }
> }
> ```

### Source: Gemini CLI transition timeline — Antigravity CLI hooks parity

> While there won't be 1:1 feature parity right out of the gate, we made sure Antigravity CLI keeps the most critical features of Gemini CLI: Agent Skills, Hooks, Subagents, and Extensions (now as Antigravity plugins). Whether you use Gemini CLI to get quick, grounded answers, scaffold and build out a new coding project, or help provision your cloud infrastructure, you can still do all of that right in Antigravity CLI.
