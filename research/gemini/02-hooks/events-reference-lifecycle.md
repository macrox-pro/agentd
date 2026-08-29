---
primary_sources:
  - id: T1-HOOKS-REF
    title: "Hooks reference"
    url: "https://geminicli.com/docs/hooks/reference.md"
    section: "SessionStart; SessionEnd; PreCompress; Notification"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Hook events — lifecycle

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Hooks reference — Lifecycle hooks

> ### `SessionStart`
>
> Fires on application startup, resuming a session, or after a `/clear` command.
> Used for loading initial context.
>
> - **Input fields**:
>   - `source`: (`"startup" | "resume" | "clear"`)
> - **Relevant output fields**:
>   - `hookSpecificOutput.additionalContext`: (`string`)
>     - **Interactive**: Injected as the first turn in history.
>     - **Non-interactive**: Prepended to the user's prompt.
>   - `systemMessage`: Shown at the start of the session.
> - **Advisory only**: `continue` and `decision` fields are **ignored**. Startup
>   is never blocked.
>
> ### `SessionEnd`
>
> Fires when the CLI exits or a session is cleared. Used for cleanup or final
> telemetry.
>
> - **Input Fields**:
>   - `reason`: (`"exit" | "clear" | "logout" | "prompt_input_exit" | "other"`)
> - **Relevant Output Fields**:
>   - `systemMessage`: Displayed to the user during shutdown.
> - **Best Effort**: The CLI **will not wait** for this hook to complete and
>   ignores all flow-control fields (`continue`, `decision`).
>
> ### `PreCompress`
>
> Fires before the CLI summarizes history to save tokens. Used for logging or
> state saving.
>
> - **Input Fields**:
>   - `trigger`: (`"auto" | "manual"`)
> - **Relevant Output Fields**:
>   - `systemMessage`: Displayed to the user before compression.
> - **Advisory Only**: Fired asynchronously. It **cannot** block or modify the
>   compression process. Flow-control fields are ignored.
>
> ---
>
> ## Stable Model API
>
> Gemini CLI uses these structures to ensure hooks don't break across SDK updates.
>
> **LLMRequest**:
>
> ```typescript
> {
>   "model": string,
>   "messages": Array<{
>     "role": "user" | "model" | "system",
>     "content": string // Non-text parts are filtered out for hooks
>   }>,
>   "config": { "temperature": number, ... },
>   "toolConfig": { "mode": string, "allowedFunctionNames": string[] }
> }
>
> ```
>
> **LLMResponse**:
>
> ```typescript
> {
>   "candidates": Array<{
>     "content": { "role": "model", "parts": string[] },
>     "finishReason": string
>   }>,
>   "usageMetadata": { "totalTokenCount": number }
> }
> ```
>
> ### `Notification`
>
> Fires when the CLI emits a system alert (for example, Tool Permissions). Used
> for external logging or cross-platform alerts.
>
> - **Input Fields**:
>   - `notification_type`: (`"ToolPermission"`)
>   - `message`: Summary of the alert.
>   - `details`: JSON object with alert-specific metadata (for example, tool name,
>     file path).
> - **Relevant Output Fields**:
>   - `systemMessage`: Displayed alongside the system alert.
> - **Observability Only**: This hook **cannot** block alerts or grant permissions
>   automatically. Flow-control fields are ignored.
