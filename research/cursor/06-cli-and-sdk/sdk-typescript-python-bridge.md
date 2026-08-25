---
primary_sources:
  - id: T1-SDK-TS
    title: "SDK"
    url: "https://cursor.com/docs/sdk/typescript.md"
    section: "SDK"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Cursor SDK

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: TypeScript SDK

> # Cursor TypeScript SDK
>
> The `@cursor/sdk` package lets you call Cursor's agent from your own code. The same agent that runs in the Cursor IDE, CLI, and web app is now scriptable from TypeScript. Run the `/sdk` skill inside Cursor to get started.
>
> ### Cookbook
>
> End-to-end examples live in the [Cursor
> Cookbook](https://github.com/cursor/cookbook): a [SDK
> quickstart](https://github.com/cursor/cookbook/tree/main/sdk/quickstart), an
> [app-builder prototyping
> tool](https://github.com/cursor/cookbook/tree/main/sdk/app-builder), a [kanban
> board for cloud
> agents](https://github.com/cursor/cookbook/tree/main/sdk/agent-kanban), and a
> [coding-agent
> CLI](https://github.com/cursor/cookbook/tree/main/sdk/coding-agent-cli). Good
> starting points for CI auto-fix bots, bug triage workers, code-review passes,
> embedded in-product agents, and orchestrators.
>
> ## Overview
>
> The SDK wraps local and cloud runtimes behind one interface. You write the same code regardless of where the agent runs.
>
> | Runtime                   | What it does                                                           | When to use                                                                                                                |
> | :------------------------ | :--------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------- |
> | **Local**                 | Runs the agent loop inline in your Node process. Files come from disk. | Dev scripts and CI checks against a working tree.                                                                          |
> | **Cloud (Cursor-hosted)** | Runs in an isolated VM with your repo cloned in. Cursor runs the VMs.  | When the caller doesn't have the repo, you want many agents in parallel, or runs need to survive the caller disconnecting. |
>
> ### Local means local agent loop, not local model
>
> "Local" describes where the agent loop and filesystem access run, not where
> the model runs. All inference goes through Cursor's hosted models in both
> modes. Local mode keeps your files on your machine; cloud mode runs in a
> Cursor environment. The model itself is hosted in either case.
>
> Runtime is picked by which key you pass to `Agent.create()` (`local` or `cloud`). Use the same `CURSOR_API_KEY` for either.
>
> For the REST API, see the [Cloud Agents API](https://cursor.com/docs/cloud-agent/api/endpoints.md). For other languages, see the [SDK Bridge](https://cursor.com/docs/sdk/bridge.md).
>
> ## Authentication
>
> Set `CURSOR_API_KEY` (or pass `apiKey`) before creating an agent. For interactive hosts without a pre-provisioned key, [`Cursor.auth.login()`](https://cursor.com/docs/sdk/typescript.md#cursorauth) mints and stores one through a browser login.
>
> The SDK accepts user API keys and service account API keys for both local and cloud runs. Team Admin API keys are not yet supported.
>
> - **User API key** from [Cursor Dashboard → API Keys](https://cursor.com/dashboard/api)
> - **Service account API key** from [Team settings](https://cursor.com/dashboard/team-settings). See [Service accounts](https://cursor.com/docs/account/enterprise/service-accounts.md)
>
> ```bash
> export CURSOR_API_KEY="your-key"
> ```
>
> ## Usage and billing
>
> SDK runs follow the same pricing, request pools, and Privacy Mode rules as runs from the IDE and Cloud Agents. Spend shows up in your team's [usage dashboard](https://cursor.com/dashboard/usage) under the SDK tag.
>
> Service account API keys bill to the team that owns the service account. User API keys bill to that user's plan.
>
> To read per-run token counts in code, see [Token usage](https://cursor.com/docs/sdk/typescript.md#token-usage). To fetch billed usage and dollar cost for an agent's runs, see [`Agent.getUsage()`](https://cursor.com/docs/sdk/typescript.md#agentgetusage).
>
> ## Core concepts
>
> | Concept        | Description                                                                                                        |
> | :------------- | :----------------------------------------------------------------------------------------------------------------- |
> | **Agent**      | Durable container that holds conversation state, workspace config, and settings. Survives across multiple prompts. |
> | **Run**        | One prompt submission. Owns its own stream, status, result, and cancellation.                                      |
> | **SDKMessage** | Normalized stream events emitted during a run. Same shape across all runtimes.                                     |
>
> ## Installation
>
> ```bash
> npm install @cursor/sdk
> ```
>
> The package name starts with `@`. The bare `cursor/sdk` doesn't exist on npm.
>
> ### Runtime support
>
> The SDK requires Node.js 22.13 or later. It ships per-platform `@cursor/sdk-<os>-<arch>` binaries for sandboxing and ripgrep, so it is a Node-first package.
>
> Importing `@cursor/sdk` does not eagerly load the local agent stack. The local executor loads on the first local `acquire`, so cloud-only and type-only consumers don't pay the local import cost. The first local agent in a process pays a one-time import, then the module stays cached.
>
> `@cursor/sdk` publishes self-contained `.d.ts` files, so types resolve without pulling in unpublished workspace packages. After upgrading, re-run your typecheck. Stream types such as `TurnEndedUpdate` resolve to real types instead of `any`.
>
> ### Single-file bundles and compiled executables
>
> `@cursor/sdk/bundled` is a self-contained, single-file build of the SDK with the same public API as `@cursor/sdk`. Use it when your app ships as one file: a standalone binary from `bun build --compile`, or a single-file bundle from esbuild.
>
> The default build loads parts of itself lazily at runtime. Single-file bundlers can't follow those loads, so a compiled app fails on the first `Agent.create()` with an error like `Cannot find module './986.js'`. The bundled entries put everything in one file, so your bundler embeds the whole SDK up front.
>
> | Entry                        | Contents                                                |
> | :--------------------------- | :------------------------------------------------------ |
> | `@cursor/sdk/bundled`        | Everything `@cursor/sdk` exports.                       |
> | `@cursor/sdk/bundled/sqlite` | `SqliteLocalAgentStore`, matching `@cursor/sdk/sqlite`. |
>
> ```typescript
> import { Agent } from "@cursor/sdk/bundled";
>
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   local: { cwd: process.cwd() },
> });
> ```
>
> Compile with Bun as usual:
>
> ```bash
> bun build --compile main.ts --outfile my-agent
> ```
>
> A few things to know:
>
> - **The bundled entries run on Bun**, including executables from `bun build --compile`. They load on Node too, but the SQLite store is unavailable there, so the default [local agent store](https://cursor.com/docs/sdk/typescript.md#local-agent-stores) falls back to JSONL. Keep importing `@cursor/sdk` in Node apps that don't ship as one file.
> - **`zod`, `@bufbuild/protobuf`, and the `@connectrpc/*` packages resolve from your own install.** They come with `@cursor/sdk`, and your bundler embeds one shared copy, so Zod schemas you pass to [custom tools](https://cursor.com/docs/sdk/typescript.md#custom-tools) keep working.
> - **Native binaries can't live inside a JavaScript bundle.** Sandboxing and the built-in ripgrep ship in the per-platform `@cursor/sdk-<os>-<arch>` packages. Place `node_modules/@cursor/sdk-<os>-<arch>/` next to your compiled executable and the SDK finds it there. Without it, search falls back to `rg` on `PATH`, and enabling [`sandboxOptions`](https://cursor.com/docs/sdk/typescript.md#sandbox-options) throws a `ConfigurationError`.
>
> Types resolve for the bundled entries the same way they do for `@cursor/sdk`. No TypeScript config changes needed.
>
> ## Quick start
>
> The fastest way in: a local agent against your current working tree, streaming events as they come in. Cloud setup is in [Creating agents](https://cursor.com/docs/sdk/typescript.md#creating-agents) below.
>
> ```typescript
> import { Agent } from "@cursor/sdk";
>
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   local: { cwd: process.cwd() },
> });
>
> const run = await agent.send("Summarize what this repository does");
>
> for await (const event of run.stream()) {
>   console.log(event);
> }
> ```
>
> Each event is a discriminated `SDKMessage`. [Streaming](https://cursor.com/docs/sdk/typescript.md#streaming) shows how to extract assistant text, handle tool calls, and clean up with `await using`. For a one-shot prompt (create, run, dispose), see [Agent.prompt()](https://cursor.com/docs/sdk/typescript.md#agentprompt).
>
> ### Quickstart approves tool calls automatically
>
> The default local agent runs tool calls (shell, edit, write, etc.) without
> asking for approval; there's no human-in-the-loop prompt in headless mode. To
> gate tool calls, configure [hooks](https://cursor.com/docs/sdk/typescript.md#hooks) (such as `beforeShellExecution` or
> `preToolUse`) or run with [`local.sandboxOptions.enabled:
>   true`](https://cursor.com/docs/sdk/typescript.md#sandbox-options).
>
> ## Creating agents
>
> ```typescript
> function Agent.create(options: AgentOptions): Promise<SDKAgent>;
> ```
>
> `Agent.create()` validates options and returns a handle immediately. Pass either `local` or `cloud` to pick a runtime.
>
> ```typescript
> // Local agent
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   local: { cwd: "/path/to/repo" },
> });
>
> // Cloud agent
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   cloud: {
>     repos: [{ url: "https://github.com/your-org/your-repo", startingRef: "main" }],
>     autoCreatePR: true,
>   },
> });
> ```
>
> `agent.agentId` is populated immediately. Local agents get an `agent-<uuid>` ID; cloud agents get a `bc-<uuid>` ID.
>
> Cloud agents started by the SDK are filtered out of the default agent list. To
> view them in Cursor Web or a Cursor window, click **Filter > Source > SDK**.
>
> ### No-repo cloud agents
>
> Cloud agents can run on an empty VM with no repository. Pass `cloud` with an empty `repos` list, or omit `repos` entirely. Omitting `cloud` selects the local runtime instead.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   cloud: { repos: [] },
> });
>
> const run = await agent.send(
>   "Research the top 3 TypeScript testing frameworks and summarize."
> );
> console.log((await run.wait()).result);
> ```
>
> No-repo agents must be enabled for your account or team. Repository-scoped API keys can't create them; use an unrestricted service account key or a user API key instead.
>
> ### Session environment variables
>
> For cloud agents, pass `cloud.envVars` when a run needs short-lived credentials or other values that should live only with that agent.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   cloud: {
>     repos: [{ url: "https://github.com/your-org/your-repo" }],
>     envVars: {
>       STAGING_API_TOKEN: process.env.STAGING_API_TOKEN!,
>     },
>   },
> });
> ```
>
> These values are encrypted at rest, injected into the cloud agent's shell, and deleted with the agent. `envVars` can't be used with a caller-supplied `agentId`; omit `agentId` and read the server-minted ID from `agent.agentId`. Variable names can't start with `CURSOR_`.
>
> For values that should only exist during a single run, pass them on `agent.send()` instead. See [Per-run environment variables](https://cursor.com/docs/sdk/typescript.md#per-run-environment-variables).
>
> ### Agent metadata
>
> Attach your own string tags to a cloud agent with `cloud.metadata`. The tags are
> persisted with the agent and returned on `SDKAgentInfo.metadata` from
> `Agent.get()` and `Agent.list()`. These tags are not the in-VM
> [agent metadata](https://cursor.com/docs/cloud-agent/metadata.md) API, which exposes the current
> run's id, owner, turn, and workspace from inside the VM.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   cloud: {
>     repos: [{ url: "https://github.com/your-org/your-repo" }],
>     metadata: {
>       end_user_id: "user-123",
>       ticket_id: "ENG-456",
>     },
>   },
> });
> ```
>
> ### Model parameters
>
> Use `model.params` to pass per-model options such as reasoning effort. Parameter ids and values vary by model. Use [`Cursor.models.list()`](https://cursor.com/docs/sdk/typescript.md#cursormodelslist) to discover supported parameters and preset variants for your account.
>
> On legacy request-based plans, Cursor enables [Max Mode](https://cursor.com/help/ai-features/max-mode.md) automatically when the selected model requires it.
>
> ### Composer 2 reroutes to Composer 2.5
>
> Composer 2 is retired. SDK requests that still pass `composer-2` or
> `composer-2-fast` are rerouted to Composer 2.5 at auth time, so existing
> scripts keep working. If you relied on the `composer-2-fast` variant, confirm
> the fast behavior still matches what you expect.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: {
>     id: "composer-2.5",
>     params: [{ id: "fast", value: "true" }],
>   },
>   local: { cwd: process.cwd() },
> });
> ```
>
> ### Cursor Router
>
> [Cursor Router](https://cursor.com/docs/cursor-router.md) selects a model for each Auto request. In the SDK, Router is the `auto-smart` model with an `optimize_for` parameter. It is available on Teams and Enterprise. Enterprise admins must enable Router for the team before `auto-smart` appears in the catalog.
>
> The Cursor SDK is an agent SDK, not a standalone model-inference or chat-completions API. Router picks models for Cursor agent runs that can reason over a workspace, call tools, run commands, and edit files. Cursor does not currently document a raw Router endpoint for arbitrary model calls.
>
> #### Select Cost, Balance, or Intelligence
>
> Pass `auto-smart` and set `optimize_for` explicitly:
>
> | Product label | SDK value      |
> | :------------ | :------------- |
> | Cost          | `cost`         |
> | Balance       | `balanced`     |
> | Intelligence  | `intelligence` |
>
> Use **Balance** in product copy. Use `balanced` only as the SDK wire value.
>
> ```typescript
> import { Agent } from "@cursor/sdk";
>
> await using agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: {
>     id: "auto-smart",
>     params: [{ id: "optimize_for", value: "balanced" }],
>   },
>   local: { cwd: process.cwd() },
> });
>
> const run = await agent.send("Find and fix the failing authentication test");
> const result = await run.wait();
>
> console.log(result.status);
> console.log(result.requestId);
> ```
>
> Always pass `optimize_for`. Do not omit it and do not send a legacy `default` value; discovery through the catalog is the supported contract.
>
> #### Discover Router in the model catalog
>
> `Cursor.models.list()` returns the models, parameter definitions, and preset variants available to the API key's current account and team. Cursor Router appears as `auto-smart` when Router is available. Team administrators can disable Router or restrict which optimization modes members may select.
>
> Treat the catalog as the source of truth before hard-coding a selection:
>
> ```typescript
> import { Cursor, type ModelSelection } from "@cursor/sdk";
>
> const models = await Cursor.models.list();
> const router = models.find((model) => model.id === "auto-smart");
> const optimizeFor = router?.parameters?.find(
>   (parameter) => parameter.id === "optimize_for",
> );
>
> if (!router || !optimizeFor) {
>   throw new Error(
>     "Cursor Router is not available for this API key. Verify that Router is enabled for the key's team.",
>   );
> }
>
> const requestedMode = "balanced";
> const allowedValues = new Set(
>   optimizeFor.values.map(({ value }) => value),
> );
>
> if (!allowedValues.has(requestedMode)) {
>   throw new Error(
>     `Router mode "${requestedMode}" is not enabled for this team.`,
>   );
> }
>
> const model: ModelSelection = {
>   id: router.id,
>   params: [{ id: optimizeFor.id, value: requestedMode }],
> };
> ```
>
> #### Switch modes per run
>
> Override the model on `agent.send()` to change Router mode for a run:
>
> ```typescript
> const run = await agent.send("Handle this complex migration", {
>   model: {
>     id: "auto-smart",
>     params: [{ id: "optimize_for", value: "intelligence" }],
>   },
> });
> ```
>
> Per-run model overrides are sticky. Later sends without an override keep using the new selection. See [Per-run model override](https://cursor.com/docs/sdk/typescript.md#per-run-model-override).
>
> #### Model ids: `auto-smart`, `auto`, and `default`
>
> | Selection                                     | Meaning                                                                                                                                     |
> | :-------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------ |
> | `auto-smart` with `optimize_for`              | Cursor Router. Use this when you want Cost, Balance, or Intelligence.                                                                       |
> | `{ id: "auto" }`                              | Server-selected Auto fallback when a specific model is missing from the catalog. Prefer `auto-smart` when you need an explicit Router mode. |
> | Omitting `optimize_for`, or sending `default` | Not a supported Router contract. Always discover allowed values and pass `cost`, `balanced`, or `intelligence`.                             |
>
> #### Billing and routing pool
>
> - All Auto modes bill at the list price of the model each request is routed to.
> - The underlying model can change between requests. Prefer a fixed model id when you need reproducible comparisons.
> - Enterprise model allowlists shape the routing pool. Blocking required models can disable Router.
>
> For current rates and the routing pool, see [Cursor Router](https://cursor.com/docs/cursor-router.md) and [Auto modes](https://cursor.com/docs/models-and-pricing.md#auto-modes).
>
> #### Troubleshooting missing Router
>
> If `auto-smart` is missing or an optimization mode is rejected:
>
> 1. Call `Cursor.models.list()`.
> 2. Confirm `auto-smart` is in the result.
> 3. Confirm `optimize_for` includes the value you want (`cost`, `balanced`, or `intelligence`).
> 4. Confirm Router is enabled for the team tied to the API key.
> 5. If you belong to multiple teams, confirm the key is operating in the intended team context.
> 6. Check team model-access policy if Router is unavailable or cannot choose a valid underlying model.
>
> ### SDKAgent
>
> The handle returned by `Agent.create()` and `Agent.resume()`.
>
> ```typescript
> interface SDKAgent {
>   readonly agentId: string;
>   readonly model: ModelSelection | undefined;
>
>   send(message: string | SDKUserMessage, options?: SendOptions): Promise<Run>;
>   close(): void;
>   reload(): Promise<void>;
>   [Symbol.asyncDispose](): Promise<void>;
>
>   listArtifacts(): Promise<SDKArtifact[]>;
>   downloadArtifact(path: string): Promise<Buffer>;
>   getUsage(options?: GetUsageOptions): Promise<AgentUsage>;
> }
> ```
>
> | Member                  | Description                                                                                                                                                                  |
> | :---------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `agentId`               | Stable agent identifier. `agent-<uuid>` for local, `bc-<uuid>` for cloud.                                                                                                    |
> | `model`                 | Current model selection. Updates after every successful `send({ model })`. `undefined` until something sets it (including resumed agents whose caller did not pass `model`). |
> | `send`                  | Start a new run with the given prompt. Returns a `Run` handle.                                                                                                               |
> | `close`                 | Begin disposal without awaiting. Fire-and-forget.                                                                                                                            |
> | `reload`                | Re-read filesystem config (hooks, project MCP, subagents) without disposing.                                                                                                 |
> | `[Symbol.asyncDispose]` | Async disposal. Pair with `await using` for automatic cleanup.                                                                                                               |
> | `listArtifacts`         | List files produced by the agent (cloud only; local returns empty).                                                                                                          |
> | `downloadArtifact`      | Download a file by path (cloud only; local throws).                                                                                                                          |
> | `getUsage`              | Fetch billed token usage and dollar cost for the agent.                                                                                                                      |
>
> ### Agent.prompt()
>
> ```typescript
> function Agent.prompt(message: string, options?: AgentOptions): Promise<RunResult>;
> ```
>
> One-shot convenience: creates an agent, sends a single prompt, waits for the run to finish, and disposes.
>
> ```typescript
> const result = await Agent.prompt("What does the auth middleware do?", {
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   local: { cwd: process.cwd() },
> });
> ```
>
> ## Sending messages
>
> Each `agent.send()` returns a `Run`. The agent retains conversation context across runs; the run is the unit of work for one prompt.
>
> ### Run
>
> ```typescript
> type RunStatus = "running" | "finished" | "error" | "cancelled";
> type RunOperation = "stream" | "wait" | "cancel" | "conversation";
>
> interface Run {
>   readonly id: string;
>   readonly requestId?: string;
>   readonly agentId: string;
>   readonly status: RunStatus;
>   readonly result?: string;
>   readonly error?: RunError;
>   readonly model?: ModelSelection;
>   readonly durationMs?: number;
>   readonly usage?: TokenUsage;
>   readonly git?: RunGitInfo;
>   readonly createdAt?: number;
>
>   stream(): AsyncGenerator<SDKMessage, void>;
>   wait(): Promise<RunResult>;
>   cancel(): Promise<void>;
>   conversation(): Promise<ConversationTurn[]>;
>
>   supports(operation: RunOperation): boolean;
>   unsupportedReason(operation: RunOperation): string | undefined;
>   onDidChangeStatus(listener: (status: RunStatus) => void): () => void;
> }
>
> interface RunGitInfo {
>   branches: Array<{ repoUrl: string; branch?: string; prUrl?: string }>;
> }
>
> interface RunError {
>   message: string;
>   code?: string;
> }
>
> interface TokenUsage {
>   inputTokens: number;
>   outputTokens: number;
>   cacheReadTokens: number;
>   cacheWriteTokens: number;
>   totalTokens: number;
>   reasoningTokens?: number;
> }
>
> interface RunResult {
>   id: string;
>   requestId?: string;
>   status: "finished" | "error" | "cancelled";
>   result?: string;
>   error?: RunError;
>   model?: ModelSelection;
>   durationMs?: number;
>   usage?: TokenUsage;
>   git?: RunGitInfo;
> }
> ```
>
> ### Streaming
>
> ```typescript
> const run = await agent.send("Find the bug in src/auth.ts");
>
> for await (const event of run.stream()) {
>   switch (event.type) {
>     case "assistant":
>       for (const block of event.message.content) {
>         if (block.type === "text") process.stdout.write(block.text);
>       }
>       break;
>     case "thinking":
>       process.stdout.write(event.text);
>       break;
>     case "tool_call":
>       console.log(`[tool] ${event.name}: ${event.status}`);
>       break;
>     case "status":
>       console.log(`[status] ${event.status}`);
>       break;
>   }
> }
>
> // Follow-up on the same agent. Conversation state from the previous
> // run is loaded automatically.
> const run2 = await agent.send("Fix it and add a regression test");
> await run2.wait();
> ```
>
> To send images alongside text:
>
> ```typescript
> const run = await agent.send({
>   text: "What's in this screenshot?",
>   images: [{ data: base64Png, mimeType: "image/png" }],
> });
> ```
>
> ### Waiting without streaming
>
> ```typescript
> const result = await run.wait();
>
> console.log(result.status);      // "finished" | "error" | "cancelled"
> console.log(result.result);      // final assistant text, if any
> console.log(result.error);       // { message, code? } when the run failed
> console.log(result.model);       // resolved ModelSelection used for this run
> console.log(result.durationMs);
> console.log(result.usage);       // cumulative TokenUsage, or undefined if unavailable
> console.log(result.git);         // { branches: [{ repoUrl, branch?, prUrl? }] } on cloud
> ```
>
> The final assistant text is on `result.result` as a string. There's no `text`, `message`, `messages`, or `content` field to dig through. If you need the per-step transcript instead, call `run.conversation()` for a structured `ConversationTurn[]` view:
>
> ```typescript
> const result = await run.wait();
> const finalText = result.result ?? "";
>
> const turns = await run.conversation();
> const lastAssistant = turns
>   .flatMap((t) => (t.type === "agentConversationTurn" ? t.turn.steps : []))
>   .filter((s) => s.type === "assistantMessage")
>   .at(-1);
>
> console.log(lastAssistant?.message.text);
> ```
>
> ### Cancelling a run
>
> ```typescript
> await run.cancel();
> ```
>
> Cancels the run. The status moves to `"cancelled"`, the live stream aborts, in-flight tool calls stop, and `run.wait()` resolves with `status: "cancelled"`. Partial output (assistant text written so far) stays on the Run object.
>
> Cancel is supported on running local and cloud runs and is a no-op if the run already finished.
>
> ### Reading run state
>
> ```typescript
> console.log(run.status);  // "running" | "finished" | "error" | "cancelled"
>
> const stop = run.onDidChangeStatus((status) => {
>   console.log(`status changed to ${status}`);
> });
> // Call `stop()` to remove the listener.
>
> // Structured per-turn view of the conversation accumulated in this run
> const turns = await run.conversation();
> ```
>
> `run.conversation()` returns the run's `ConversationTurn[]` (an agent turn with steps, or a shell turn with command and output). Use it to render or persist the run's structured history without subscribing to the live stream.
>
> ### Token usage
>
> Runs report token usage when the runtime provides it. Read the cumulative total from `run.usage` while the run is in flight, or from `result.usage` after `run.wait()`. Both hold a `TokenUsage` summed across every turn that reported usage, and both are `undefined` when no turn did (a cancelled run that never finished a turn, or a runtime that doesn't surface usage).
>
> ```typescript
> interface TokenUsage {
>   inputTokens: number;
>   outputTokens: number;
>   cacheReadTokens: number;
>   cacheWriteTokens: number;
>   totalTokens: number;
>   reasoningTokens?: number;
> }
> ```
>
> | Field              | Description                                                                                       |
> | :----------------- | :------------------------------------------------------------------------------------------------ |
> | `inputTokens`      | Prompt tokens sent to the model.                                                                  |
> | `outputTokens`     | Tokens generated by the model.                                                                    |
> | `cacheReadTokens`  | Tokens served from the prompt cache.                                                              |
> | `cacheWriteTokens` | Tokens written to the prompt cache.                                                               |
> | `totalTokens`      | `inputTokens + outputTokens + cacheReadTokens + cacheWriteTokens`. Excludes `reasoningTokens`.    |
> | `reasoningTokens`  | Reasoning tokens, a subset of `outputTokens`. Omitted when the model or runtime didn't report it. |
>
> ```typescript
> const result = await run.wait();
>
> if (result.usage) {
>   console.log(`total: ${result.usage.totalTokens}`);
>   console.log(`in: ${result.usage.inputTokens}, out: ${result.usage.outputTokens}`);
>   console.log(
>     `cache read/write: ${result.usage.cacheReadTokens}/${result.usage.cacheWriteTokens}`
>   );
> } else {
>   console.log("no usage reported for this run");
> }
> ```
>
> `reasoningTokens` is already counted inside `outputTokens`, so `totalTokens` leaves it out to avoid double-counting.
>
> For per-turn numbers as they stream, handle the `usage` [stream event](https://cursor.com/docs/sdk/typescript.md#stream-events) (`SDKUsageMessage`). It fires once at the end of each turn that reported usage and carries that turn's `TokenUsage`. `run.usage` and `result.usage` stay cumulative across the run.
>
> ```typescript
> for await (const event of run.stream()) {
>   if (event.type === "usage") {
>     console.log(`turn used ${event.usage.totalTokens} tokens`);
>   }
> }
> ```
>
> Token counts are what the runtime reports; they say nothing about cost. For billed usage and the dollar cost of an agent's runs, call [`agent.getUsage()`](https://cursor.com/docs/sdk/typescript.md#agentgetusage).
>
> ### Run correlation with requestId
>
> Every `agent.send()` gets a platform-generated UUID, exposed as `requestId` on both the `Run` and the `RunResult`. Use it to tie a script or CI run to backend logs, analytics, and support threads instead of guessing from `agentId` alone.
>
> ```typescript
> const run = await agent.send("Audit the auth middleware");
> console.log(run.requestId); // e.g. "6e0d261c-86a2-4383-89f0-9162c1c10662"
>
> const result = await run.wait();
> logger.info({ requestId: result.requestId }, "run finished");
> ```
>
> `requestId` persists with the run, so it round-trips through the in-memory, SQLite, and JSONL [local stores](https://cursor.com/docs/sdk/typescript.md#local-agent-stores) and is set on cloud runs when the backend returns one. Log it alongside `error.requestId` from [errors](https://cursor.com/docs/sdk/typescript.md#errors) so a single identifier spans both success and failure paths.
>
> ### Per-run model override
>
> The `model` you pass to `agent.send()` overrides the agent's selection for that run, then becomes sticky: subsequent sends without an override continue to use the new model. To switch back, pass another `model` override or read the current selection from `agent.model`.
>
> ```typescript
> const run = await agent.send("Plan the refactor", {
>   model: { id: "composer-2.5", params: [{ id: "fast", value: "true" }] },
> });
>
> console.log(agent.model);  // updated to the override after the send succeeds
> ```
>
> `run.model` and `result.model` reflect the selection that this specific run actually used and are immutable once the run starts.
>
> ### Per-run environment variables
>
> Cloud agents can also take environment variables for a single run. Pass `cloud.envVars` on `agent.send()` and the values are injected into the agent's shell for that run only — when the run finishes, they're removed from the VM and the next run doesn't see them. This is the right shape for credentials that rotate between turns, like a short-lived deploy token you mint right before asking the agent to use it.
>
> ```typescript
> const run = await agent.send("Deploy the preview environment", {
>   cloud: {
>     envVars: {
>       DEPLOY_TOKEN: await mintShortLivedToken(),
>     },
>   },
> });
> ```
>
> If a run-scoped variable has the same name as an agent-scoped one from [`cloud.envVars` on `Agent.create()`](https://cursor.com/docs/sdk/typescript.md#session-environment-variables), the run-scoped value wins for that run, then the agent-scoped value comes back on the next run.
>
> Per-run variables work on the first send too. The SDK passes them along with agent creation, scoped to the initial run, so they aren't persisted on the agent. Like agent-scoped variables, they're encrypted at rest and names can't start with `CURSOR_`.
>
> Per-run environment variables are cloud agents only, and they aren't available for agents running against public repositories. For local agents, the agent process inherits your own environment, so set variables on the process before calling `send()`.
>
> ### Conversation mode
>
> Pass `mode: "plan"` or `mode: "agent"` to control whether a run explores and plans first or implements changes directly. See [Plan mode](https://cursor.com/help/ai-features/plan-mode.md) for what plan mode does in the product.
>
> Set `mode` on `Agent.create()` to seed the first run. On follow-up `agent.send()` calls, omit `mode` to keep the conversation's current mode, or pass `mode` to switch for that run only.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   mode: "plan",
>   cloud: {
>     repos: [{ url: "https://github.com/your-org/your-repo" }],
>   },
> });
>
> await (await agent.send("Design the auth refactor")).wait();
> await (await agent.send("Looks good, start building", { mode: "agent" })).wait();
> ```
>
> ### Streaming raw deltas
>
> `run.stream()` yields normalized `SDKMessage` events. For lower-level updates (per-token text, tool-call args streaming in, thinking deltas, nested task updates, step boundaries), pass `onDelta` and `onStep` callbacks to `send()`:
>
> ```typescript
> const run = await agent.send("Refactor the utils module", {
>   onDelta: ({ update }) => {
>     if (update.type === "text-delta") process.stdout.write(update.text);
>     if (update.type === "thinking-delta") process.stdout.write(update.text);
>   },
>   onStep: ({ step }) => {
>     console.log(`[step] ${step.type}`);
>   },
> });
> ```
>
> The callbacks are awaited before the next update is processed, so you can apply backpressure. `InteractionUpdate` covers `text-delta`, `thinking-delta`, `thinking-completed`, `tool-call-started`, `tool-call-completed`, `tool-call-delta`, `partial-tool-call`, `token-delta`, `step-started`, `step-completed`, `turn-ended`, and a handful of summary and shell-output deltas.
>
> ### Per-send options
>
> | Property            | Type                                          | Description                                                                                                                                                                                                                                       |
> | :------------------ | :-------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `model`             | `ModelSelection`                              | Per-send model override. If omitted, uses `agent.model`. Sticky: a successful send updates `agent.model`.                                                                                                                                         |
> | `mode`              | `"agent" \| "plan"`                           | Per-send conversation mode override. If omitted on follow-ups, keeps the conversation's current mode.                                                                                                                                             |
> | `mcpServers`        | `Record<string, McpServerConfig>`             | Inline MCP server definitions. Fully replaces creation-time servers for this run.                                                                                                                                                                 |
> | `onStep`            | `(args: { step }) => void \| Promise<void>`   | Callback after each completed conversation step (text, thinking, or tool batch).                                                                                                                                                                  |
> | `onDelta`           | `(args: { update }) => void \| Promise<void>` | Callback per raw `InteractionUpdate`.                                                                                                                                                                                                             |
> | `idempotencyKey`    | `string`                                      | Optional client-generated idempotency key for the send.                                                                                                                                                                                           |
> | `cloud.envVars`     | `Record<string, string>`                      | Cloud agents only. [Per-run environment variables](https://cursor.com/docs/sdk/typescript.md#per-run-environment-variables) injected for this run and removed when it finishes. Overrides agent-scoped `cloud.envVars` by name for this run only. |
> | `local.force`       | `boolean`                                     | Local agents only. Defaults to `false`. Expire a stuck active run before starting this message. Cloud returns `409 agent_busy` server-side, so no equivalent is needed.                                                                           |
> | `local.customTools` | `Record<string, SDKCustomTool>`               | Local agents only. [Custom tools](https://cursor.com/docs/sdk/typescript.md#custom-tools) for this run. Replaces the agent's creation-time `local.customTools` for that run.                                                                      |
>
> ***
>
> The next three sections are detailed reference for `SDKMessage`, `InteractionUpdate`, and `ConversationTurn`. Skim or skip on a first read; [Resuming agents](https://cursor.com/docs/sdk/typescript.md#resuming-agents) picks up the narrative.
>
> ## Stream events
>
> Events from `run.stream()`. Discriminate on `type`. All events include `agent_id` and `run_id`.
>
> ```typescript
> type SDKMessage =
>   | SDKSystemMessage
>   | SDKUserMessageEvent
>   | SDKAssistantMessage
>   | SDKThinkingMessage
>   | SDKToolUseMessage
>   | SDKStatusMessage
>   | SDKTaskMessage
>   | {
>       type: "request";
>       agent_id: string;
>       run_id: string;
>       request_id: string;
>     }
>   | SDKUsageMessage;
> ```
>
> | `type`        | Description                                                                                      | Key fields                                                                      |
> | :------------ | :----------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------ |
> | `"system"`    | Init metadata. Emitted once at the start of a run.                                               | `subtype?` (`"init"`), `model?`, `tools?`                                       |
> | `"user"`      | Echo of the user prompt for this run.                                                            | `message.content: TextBlock[]`                                                  |
> | `"assistant"` | Model text output.                                                                               | `message.content: (TextBlock \| ToolUseBlock)[]`                                |
> | `"thinking"`  | Reasoning content.                                                                               | `text`, `thinking_duration_ms?`                                                 |
> | `"tool_call"` | Tool invocation lifecycle. Emitted at start with `args`, then again on completion with `result`. | `call_id`, `name`, `status`, `args?`, `result?`, `truncated?`                   |
> | `"status"`    | Cloud run lifecycle transitions.                                                                 | `status`, `message?`                                                            |
> | `"task"`      | Task-level milestones and summaries.                                                             | `status?`, `text?`                                                              |
> | `"request"`   | Awaiting user input or approval.                                                                 | `request_id`                                                                    |
> | `"usage"`     | Per-turn token usage, emitted once at turn end when the runtime reported it.                     | `usage` ([`TokenUsage`](https://cursor.com/docs/sdk/typescript.md#token-usage)) |
>
> Result data (final text, model, duration, cumulative token usage, git metadata) lives on the `Run` object after the stream completes. Use `run.wait()` to read it.
>
> > **Tool call schema is not stable.** The `args` and `result` payloads on `tool_call` events reflect each tool's internal shape and can change as tools evolve. Tool names can also be renamed or replaced. Treat `args` and `result` as `unknown` and parse defensively. The event envelope (`type`, `call_id`, `name`, `status`) is stable.
>
> ### Message types
>
> ```typescript
> interface SDKSystemMessage {
>   type: "system";
>   subtype?: "init";
>   agent_id: string;
>   run_id: string;
>   model?: ModelSelection;
>   tools?: string[];
> }
>
> interface SDKUserMessageEvent {
>   type: "user";
>   agent_id: string;
>   run_id: string;
>   message: { role: "user"; content: TextBlock[] };
> }
>
> interface SDKAssistantMessage {
>   type: "assistant";
>   agent_id: string;
>   run_id: string;
>   message: {
>     role: "assistant";
>     content: Array<TextBlock | ToolUseBlock>;
>   };
> }
>
> interface SDKThinkingMessage {
>   type: "thinking";
>   agent_id: string;
>   run_id: string;
>   text: string;
>   thinking_duration_ms?: number;
> }
>
> interface SDKToolUseMessage {
>   type: "tool_call";
>   agent_id: string;
>   run_id: string;
>   call_id: string;
>   name: string;
>   status: "running" | "completed" | "error";
>   args?: unknown;
>   result?: unknown;
>   truncated?: { args?: boolean; result?: boolean };
> }
>
> interface SDKStatusMessage {
>   type: "status";
>   agent_id: string;
>   run_id: string;
>   status: "CREATING" | "RUNNING" | "FINISHED" | "ERROR" | "CANCELLED" | "EXPIRED";
>   message?: string;
> }
>
> interface SDKTaskMessage {
>   type: "task";
>   agent_id: string;
>   run_id: string;
>   status?: string;
>   text?: string;
> }
>
> interface SDKUsageMessage {
>   type: "usage";
>   agent_id: string;
>   run_id: string;
>   usage: TokenUsage;
> }
>
> interface TextBlock {
>   type: "text";
>   text: string;
> }
>
> interface ToolUseBlock {
>   type: "tool_use";
>   id: string;
>   name: string;
>   input: unknown;
> }
> ```
>
> `SDKToolUseMessage` is emitted twice for most tool calls: first with `status: "running"` and `args` populated, then again on completion with `status: "completed"` (or `"error"`) and `result` populated. `truncated` flags whether the SDK truncated `args` or `result` because the payload was too large.
>
> `SDKStatusMessage` covers cloud-side lifecycle transitions. `CREATING` covers VM provisioning and repo cloning; `RUNNING` is the agent doing work; the rest are terminal.
>
> `SDKUsageMessage` is emitted once at the end of each turn that reported token usage, carrying that turn's [`TokenUsage`](https://cursor.com/docs/sdk/typescript.md#token-usage). The cumulative total across turns stays on `run.usage` and `result.usage`. See [Token usage](https://cursor.com/docs/sdk/typescript.md#token-usage).
>
> ## Interaction updates
>
> `InteractionUpdate` is the raw delta type passed to the `onDelta` callback on `agent.send()`. Updates are finer-grained than `SDKMessage` events: text streams in token-by-token, tool calls report partial state as args accumulate, thinking arrives as it happens.
>
> ```typescript
> type InteractionUpdate =
>   | TextDeltaUpdate
>   | ThinkingDeltaUpdate
>   | ThinkingCompletedUpdate
>   | ToolCallStartedUpdate
>   | ToolCallCompletedUpdate
>   | ToolCallDeltaUpdate
>   | PartialToolCallUpdate
>   | TokenDeltaUpdate
>   | StepStartedUpdate
>   | StepCompletedUpdate
>   | TurnEndedUpdate
>   | UserMessageAppendedUpdate
>   | SummaryUpdate
>   | SummaryStartedUpdate
>   | SummaryCompletedUpdate
>   | ShellOutputDeltaUpdate;
> ```
>
> ### Update types
>
> ```typescript
> interface TextDeltaUpdate {
>   type: "text-delta";
>   text: string;
> }
>
> interface ThinkingDeltaUpdate {
>   type: "thinking-delta";
>   text: string;
> }
>
> interface ThinkingCompletedUpdate {
>   type: "thinking-completed";
>   thinkingDurationMs: number;
> }
>
> interface ToolCallStartedUpdate {
>   type: "tool-call-started";
>   callId: string;
>   toolCall: ToolCall;
>   modelCallId: string;
> }
>
> interface PartialToolCallUpdate {
>   type: "partial-tool-call";
>   callId: string;
>   toolCall: ToolCall;
>   modelCallId: string;
> }
>
> interface ToolCallCompletedUpdate {
>   type: "tool-call-completed";
>   callId: string;
>   toolCall: ToolCall;
>   modelCallId: string;
> }
>
> interface ToolCallDeltaUpdate {
>   type: "tool-call-delta";
>   callId: string;
>   modelCallId: string;
>   taskUpdate: NestedTaskUpdate;
> }
>
> type NestedTaskUpdate =
>   | TextDeltaUpdate
>   | ToolCallStartedUpdate
>   | ToolCallCompletedUpdate
>   | ThinkingDeltaUpdate
>   | ThinkingCompletedUpdate
>   | PartialToolCallUpdate
>   | StepStartedUpdate
>   | StepCompletedUpdate;
>
> interface TokenDeltaUpdate {
>   type: "token-delta";
>   tokens: number;
> }
>
> interface StepStartedUpdate {
>   type: "step-started";
>   stepId: number;
> }
>
> interface StepCompletedUpdate {
>   type: "step-completed";
>   stepId: number;
>   stepDurationMs: number;
> }
>
> interface TurnEndedUpdate {
>   type: "turn-ended";
>   usage?: {
>     inputTokens: number;
>     outputTokens: number;
>     cacheReadTokens: number;
>     cacheWriteTokens: number;
>     reasoningTokens?: number;
>   };
> }
>
> interface UserMessageAppendedUpdate {
>   type: "user-message-appended";
>   userMessage: UserMessage;
> }
>
> interface SummaryUpdate {
>   type: "summary";
>   summary: string;
> }
>
> interface SummaryStartedUpdate {
>   type: "summary-started";
> }
>
> interface SummaryCompletedUpdate {
>   type: "summary-completed";
> }
>
> interface ShellOutputDeltaUpdate {
>   type: "shell-output-delta";
>   event: Record<string, unknown>;
> }
> ```
>
> `ToolCallDeltaUpdate` carries one level of nested interaction updates from a task or subagent tool call. `PartialToolCallUpdate` is emitted as the model streams arguments into a tool call before it commits. The same stability disclaimer that applies to `SDKToolUseMessage.args` applies here.
>
> ## Conversation types
>
> The structured per-turn view of a run, returned by `run.conversation()` and used in the `onStep` callback's argument.
>
> ```typescript
> type ConversationTurn =
>   | { type: "agentConversationTurn"; turn: AgentConversationTurn }
>   | { type: "shellConversationTurn"; turn: ShellConversationTurn };
>
> interface AgentConversationTurn {
>   userMessage?: UserMessage;
>   steps: ConversationStep[];
> }
>
> interface ShellConversationTurn {
>   shellCommand?: ShellCommand;
>   shellOutput?: ShellOutput;
> }
>
> type ConversationStep =
>   | { type: "assistantMessage"; message: AssistantMessage }
>   | { type: "toolCall"; message: ToolCall }
>   | { type: "thinkingMessage"; message: ThinkingMessage };
>
> interface AssistantMessage {
>   text: string;
> }
>
> interface ThinkingMessage {
>   text: string;
>   thinkingDurationMs?: number;
> }
>
> interface UserMessage {
>   text: string;
> }
>
> interface ShellCommand {
>   command: string;
>   workingDirectory?: string;
> }
>
> interface ShellOutput {
>   stdout: string;
>   stderr: string;
>   exitCode: number;
> }
> ```
>
> `ToolCall` is a discriminated union over every built-in tool (shell, edit, read, write, glob, grep, ls, semSearch, mcp, task, and others). Its shape is internal-facing; see the [stability note](https://cursor.com/docs/sdk/typescript.md#stream-events) under Stream events.
>
> ## Resuming agents
>
> ```typescript
> function Agent.resume(agentId: string, options?: Partial<AgentOptions>): Promise<SDKAgent>;
> ```
>
> Use `Agent.resume()` to reattach to an existing agent by ID. Common flows: reconnecting to a long-running cloud agent that was kicked off earlier, or continuing a conversation after the local process restarted. Runtime is auto-detected from the ID prefix (`bc-` is cloud, anything else is local).
>
> ```typescript
> await using agent = await Agent.resume("bc-abc123", {
>   apiKey: process.env.CURSOR_API_KEY!,
> });
>
> const run = await agent.send("Also update the changelog");
> await run.wait();
> ```
>
> `agent.model` is `undefined` on resume unless you pass `model` again. Inline `mcpServers` are not persisted across resume — they often carry secrets and live in memory only. Pass them again on resume, or use file-based MCP config (`.cursor/mcp.json` + `local.settingSources`) for servers that should survive.
>
> ## Inspecting agents and runs
>
> List, fetch, and reload past agents. List endpoints return `{ items, nextCursor? }` for cursor-based pagination.
>
> ### Agent.list()
>
> ```typescript
> function Agent.list(options?: ListAgentsOptions): Promise<ListResult<SDKAgentInfo>>;
>
> type ListAgentsOptions = {
>   limit?: number;
>   cursor?: string;
> } & (
>   | { runtime?: undefined }
>   | { runtime: "local"; cwd?: string; store?: LocalAgentStore }
>   | {
>       runtime: "cloud";
>       prUrl?: string;
>       includeArchived?: boolean;
>       apiKey?: string;
>     }
> );
> ```
>
> ```typescript
> const { items, nextCursor } = await Agent.list({
>   runtime: "local",
>   cwd: process.cwd(),
> });
> ```
>
> ### Agent.get()
>
> ```typescript
> function Agent.get(agentId: string, options?: GetAgentOptions): Promise<SDKAgentInfo>;
>
> interface GetAgentOptions {
>   cwd?: string;       // local routing
>   apiKey?: string;    // cloud routing
>   store?: LocalAgentStore;
> }
> ```
>
> Runtime is auto-detected from the agent ID prefix (`bc-` → cloud, otherwise local).
>
> ### Agent.listRuns()
>
> ```typescript
> function Agent.listRuns(agentId: string, options?: ListRunsOptions): Promise<ListResult<Run>>;
>
> type ListRunsOptions = {
>   limit?: number;
>   cursor?: string;
> } & (
>   | { runtime?: "local"; cwd?: string; store?: LocalAgentStore }
>   | { runtime: "cloud"; apiKey?: string }
> );
> ```
>
> ### Agent.getRun()
>
> ```typescript
> function Agent.getRun(runId: string, options?: GetRunOptions): Promise<Run>;
>
> type GetRunOptions =
>   | { runtime?: "local"; cwd?: string; store?: LocalAgentStore }
>   | { runtime: "cloud"; agentId: string; apiKey?: string };
> ```
>
> Cloud `getRun` requires the parent `agentId`.
>
> ### Agent.cancelRun()
>
> ```typescript
> function Agent.cancelRun(runId: string, options?: GetRunOptions): Promise<void>;
> ```
>
> Cancels a run when you have its ID but do not have a `Run` handle.
>
> ### Agent.messages.list()
>
> ```typescript
> Agent.messages.list(
>   agentId: string,
>   options?: GetAgentMessagesOptions
> ): Promise<AgentMessage[]>;
>
> interface GetAgentMessagesOptions {
>   limit?: number;
>   offset?: number;
>   runtime?: "local";
>   cwd?: string;
>   store?: LocalAgentStore;
> }
>
> interface AgentMessage {
>   type: "user" | "assistant";
>   uuid: string;
>   agent_id: string;
>   message: unknown;
> }
> ```
>
> Returns the stored user and assistant messages for a local agent.
>
> ### Agent.getUsage()
>
> Fetch billed token usage and dollar cost for an agent's runs. Call it on a handle, or statically by ID when you don't have one. Cloud agents return a per-run breakdown; local agents return a per-turn breakdown. Pass `runId` to restrict the result to one entry: for cloud agents a `run-<uuid>` run ID, for local agents an ID from a previous `getUsage().runs[].runId`.
>
> ```typescript
> agent.getUsage(options?: GetUsageOptions): Promise<AgentUsage>;
> function Agent.getUsage(
>   agentId: string,
>   options?: GetUsageOptions & { apiKey?: string }
> ): Promise<AgentUsage>;
>
> interface GetUsageOptions {
>   runId?: string;
> }
>
> interface AgentUsage {
>   usage: TokenUsage;   // summed across `runs`
>   cost?: UsageCost;    // summed across `runs`
>   runs: RunUsage[];
> }
>
> interface RunUsage {
>   runId: string;
>   usage: TokenUsage;
>   cost?: UsageCost;
> }
>
> interface UsageCost {
>   rawCostCents: number;   // undiscounted model token cost; 0 for request-priced usage
>   chargedCents: number;   // amount charged, discounts and the Cursor Token Fee included
> }
> ```
>
> ```typescript
> const { usage, cost, runs } = await agent.getUsage();
>
> console.log(`tokens: ${usage.totalTokens}`);
> if (cost) {
>   console.log(`charged: $${(cost.chargedCents / 100).toFixed(2)}`);
> }
> for (const run of runs) {
>   console.log(run.runId, run.usage.totalTokens, run.cost?.chargedCents);
> }
> ```
>
> Cost includes discounts and can take a moment to settle after a run ends; `cost` is absent until it does. `chargedCents` is `0` for plan-included, BYOK, and credit-grant usage.
>
> This is a different view from [Token usage](https://cursor.com/docs/sdk/typescript.md#token-usage): `run.usage` is the live token count for one run, while `getUsage()` is the billed record across the agent's runs.
>
> ### Cloud agent lifecycle
>
> Cloud agents stay in your team's workspace until you archive or delete them. `Agent.list({ runtime: "cloud" })` hides archived agents by default; pass `includeArchived: true` to see them. Filter by `prUrl` to find the agent that opened a specific pull request.
>
> ```typescript
> function Agent.archive(agentId: string, options?: AgentOperationOptions): Promise<void>;
> function Agent.unarchive(agentId: string, options?: AgentOperationOptions): Promise<void>;
> function Agent.delete(agentId: string, options?: AgentOperationOptions): Promise<void>;
>
> interface AgentOperationOptions {
>   cwd?: string;
>   apiKey?: string;
>   store?: LocalAgentStore;
> }
> ```
>
> ```typescript
> await Agent.archive(agentId);     // soft-delete; transcript stays readable
> await Agent.unarchive(agentId);   // restore an archived agent
> await Agent.delete(agentId);      // permanent; subsequent reads return 404
> ```
>
> ### SDKAgentInfo
>
> The metadata shape returned by `Agent.list()` and `Agent.get()`.
>
> ```typescript
> type SDKAgentInfo = {
>   agentId: string;
>   name: string;
>   summary: string;
>   lastModified: number;
>   status?: "running" | "finished" | "error";
>   createdAt?: number;
>   archived?: boolean;
> } & (
>   | { runtime?: undefined }
>   | { runtime: "local"; cwd?: string }
>   | {
>       runtime: "cloud";
>       env?: { type: "cloud" | "pool" | "machine"; name?: string };
>       repos?: string[];
>       metadata?: Record<string, string>;
>     }
> );
> ```
>
> ## The Cursor namespace
>
> Account-level reads, catalog reads, and process-wide SDK configuration. The read methods take an optional `{ apiKey }` and otherwise fall back to `CURSOR_API_KEY`, then to a stored [browser login](https://cursor.com/docs/sdk/typescript.md#cursorauth).
>
> ### Cursor.auth
>
> Interactive login for hosts without a pre-provisioned API key. `Cursor.auth.login()` opens the Cursor website's login page in a browser, waits for completion, and mints a user API key (90 days by default) stored in `~/.cursor/sdk/auth.json`. After login, `Agent.create()`, `Cursor.me()`, and the other reads work without `apiKey` or `CURSOR_API_KEY`.
>
> ```typescript
> import { Cursor } from "@cursor/sdk";
>
> await Cursor.auth.login();
>
> const status = await Cursor.auth.status();
> // { status: "logged-in", backendUrl, email?, apiKeyExpiresAtMs? }
> // | { status: "logged-out" }
>
> await Cursor.auth.logout();
> ```
>
> | Option        | Description                                                                                                                                                                              |
> | :------------ | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `backendUrl`  | API base URL. Defaults to `CURSOR_BACKEND_URL`, then production.                                                                                                                         |
> | `websiteUrl`  | Browser login base URL. Defaults to `CURSOR_WEBSITE_URL`, then production.                                                                                                               |
> | `openBrowser` | `true` (default) opens the system browser when that is likely to work; `false` never opens one; a function is a custom opener. Skipped in SSH sessions or when `NO_OPEN_BROWSER` is set. |
> | `onLoginUrl`  | Called with the login URL before waiting, so the host can display it. When omitted and no browser opened, the URL is written to stderr.                                                  |
> | `signal`      | `AbortSignal` that cancels the wait; `login` then throws an `AuthenticationError`.                                                                                                       |
> | `store`       | Where to persist the credentials. Defaults to `~/.cursor/sdk/auth.json`; pass `null` to only receive the key in the result.                                                              |
> | `apiKeyName`  | Display name of the minted key in the dashboard's API-keys list.                                                                                                                         |
> | `apiKeyTtlMs` | Lifetime of the minted key in milliseconds. Defaults to 90 days.                                                                                                                         |
>
> `Cursor.auth.login()` returns
> `{ apiKey, email?, apiKeyExpiresAtMs }`. Use `FileCredentialStore` or
> `InMemoryCredentialStore` to provide a custom `store` to `login()`, `status()`,
> or `logout()`.
>
> Credential resolution order everywhere in the SDK: explicit `apiKey`, then `CURSOR_API_KEY`, then the stored login. The stored login does not read credentials from a local Cursor app installation; it only holds keys minted by `Cursor.auth.login()`.
>
> ### Cursor.configure()
>
> ```typescript
> function Cursor.configure(options: CursorConfigureOptions): void;
>
> interface CursorConfigureOptions {
>   local?: {
>     store?: LocalAgentStore | null;
>     useHttp1ForAgent?: boolean | null;
>     workspaceScanCacheTtlMs?: number | null;
>   };
> }
> ```
>
> Set defaults for local agents that apply to later `Agent.*` calls. Fields on an individual call override these values; pass `null` to clear a previous default.
>
> | Option                          | Description                                                                                                                                                                                                                                                                                                                                                                               |
> | :------------------------------ | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `local.store`                   | Default [local agent store](https://cursor.com/docs/sdk/typescript.md#local-agent-stores) when a call omits `local.store`. The SDK uses on-disk SQLite when the SQLite backend is available and falls back to `JsonlLocalAgentStore` otherwise.                                                                                                                                           |
> | `local.useHttp1ForAgent`        | Force local agent backend streams to use HTTP/1.1 with SSE instead of HTTP/2. Useful behind proxies or on fetch stacks that don't support HTTP/2.Bun defaults to HTTP/1.1 due to upstream HTTP/2 compatibility issues.                                                                                                                                                                    |
> | `local.workspaceScanCacheTtlMs` | How long the SDK reuses a workspace scan (rules, skills, `AGENTS.md`, ignore files), in milliseconds. Defaults to 20 seconds. Raise it in a long-lived host serving a checkout that only changes on deploy; the trade is freshness, since a rule added after the process started can go unseen for this long. The `CURSOR_RIPWALK_CACHE_TTL_MS` environment variable sets the same value. |
>
> ```typescript
> import { Cursor, JsonlLocalAgentStore } from "@cursor/sdk";
>
> Cursor.configure({
>   local: {
>     store: new JsonlLocalAgentStore("/var/lib/cursor-agents"),
>     useHttp1ForAgent: true,
>   },
> });
> ```
>
> ### Cursor.me()
>
> ```typescript
> function Cursor.me(options?: CursorRequestOptions): Promise<SDKUser>;
>
> interface CursorRequestOptions {
>   apiKey?: string;
> }
>
> interface SDKUser {
>   apiKeyName: string;
>   userId?: number;
>   userEmail?: string;
>   userFirstName?: string;
>   userLastName?: string;
>   createdAt: string;
> }
> ```
>
> ### Cursor.models.list()
>
> ```typescript
> function Cursor.models.list(options?: CursorRequestOptions): Promise<SDKModel[]>;
>
> type SDKModel = ModelListItem;
>
> interface ModelListItem {
>   id: string;
>   displayName: string;
>   description?: string;
>   aliases?: string[];
>   parameters?: ModelParameterDefinition[];
>   variants?: ModelVariant[];
> }
>
> interface ModelParameterDefinition {
>   id: string;
>   displayName?: string;
>   values: Array<{ value: string; displayName?: string }>;
> }
>
> interface ModelVariant {
>   params: ModelParameterValue[];
>   displayName: string;
>   description?: string;
>   isDefault?: boolean;
> }
> ```
>
> Use `Cursor.models.list()` to discover valid `model` ids and per-model `params` before calling `Agent.create()` or `agent.send()`. Parameters are model-specific. Common examples include reasoning effort and Cursor Router's `optimize_for` on `auto-smart`.
>
> The catalog is account- and team-specific. Cursor Router only appears as `auto-smart` when Router is available for the API key's team. See [Cursor Router](https://cursor.com/docs/sdk/typescript.md#cursor-router).
>
> ```typescript
> const models = await Cursor.models.list();
> const composer = models.find((model) => model.id === "composer-2.5");
>
> console.log(composer?.parameters);
> // [
> //   {
> //     id: "fast",
> //     displayName: "Fast",
> //     values: [
> //       { value: "false" },
> //       { value: "true", displayName: "Fast" },
> //     ],
> //   },
> // ]
> ```
>
> Pass selected parameter values through `model.params`. Preset `variants` already contain valid `params`, so you can copy them into a model selection.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: {
>     id: "composer-2.5",
>     params: [{ id: "fast", value: "true" }],
>   },
>   local: { cwd: process.cwd() },
> });
> ```
>
> #### Best practices
>
> - **Discover, don't hard-code.** Call `Cursor.models.list()` at startup (or once per process) and cache the result. Model ids and parameter shapes can change as new models ship.
> - **Pass parameters explicitly when the model expects them.** A model whose `parameters` array is non-empty is a parameterized model. Send the params you want; otherwise the run uses each parameter's first allowed value, which may not match what you intend. For Cursor Router, always pass `optimize_for` explicitly.
> - **Resolve by capability, not id.** If you want "the current default in fast mode" rather than a specific model, look it up:
>
>   ```typescript
>   const models = await Cursor.models.list();
>   const composer = models.find((m) => m.id === "composer-2.5");
>   const fast = composer?.parameters?.find((p) => p.id === "fast");
>   const fastValue = fast?.values.find((v) => v.value === "true")?.value;
>
>   const model = composer
>     ? {
>         id: composer.id,
>         params: fastValue ? [{ id: "fast", value: fastValue }] : undefined,
>       }
>     : {
>         id: "auto-smart",
>         params: [{ id: "optimize_for", value: "balanced" }],
>       };
>   ```
>
>   Prefer an explicit Router selection (`auto-smart` + `optimize_for`) when the target model is missing. Fall back to `{ id: "auto" }` only when you want server-selected Auto without choosing Cost, Balance, or Intelligence.
>
> ### Cursor.repositories.list()
>
> ```typescript
> function Cursor.repositories.list(options?: CursorRequestOptions): Promise<SDKRepository[]>;
>
> interface SDKRepository {
>   url: string;
> }
> ```
>
> Returns the GitHub repositories connected for the calling user's team. Cloud only.
>
> ## Configuration sources at a glance
>
> MCP servers, subagents, and hooks all resolve from a mix of inline options and on-disk config. The precedence is the same shape across the three: per-send inline > creation-time inline > project files > user files > team / dashboard config.
>
> | Feature              | Inline option                                               | Local file (project)                                                       | Local file (user)                        | Cloud / dashboard                                                                      | Precedence                                                                                              |
> | :------------------- | :---------------------------------------------------------- | :------------------------------------------------------------------------- | :--------------------------------------- | :------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------ |
> | **MCP servers**      | `mcpServers` on `Agent.create()` and `agent.send()`         | `.cursor/mcp.json` (gated by `local.settingSources` including `"project"`) | `~/.cursor/mcp.json` (gated by `"user"`) | Servers configured at [cursor.com/agents](https://cursor.com/agents) (cloud only)      | Send > create > plugins > project > user (local); Send > create > dashboard (cloud)                     |
> | **Subagents**        | `agents` on `Agent.create()`                                | `.cursor/agents/*.md` (frontmatter: `name`, `description`, `model?`)       | n/a                                      | Cloud picks up the same project files when the agent runs against the cloned repo      | Inline overrides file-based with the same name                                                          |
> | **Hooks**            | None — file-based only                                      | `.cursor/hooks.json` (+ scripts)                                           | `~/.cursor/hooks.json`                   | Cloud runs project hooks. On Enterprise plans, also team and enterprise-managed hooks. | File-based; project layered with user / team / enterprise per [Hooks](https://cursor.com/docs/hooks.md) |
> | **Settings sources** | `local.settingSources` selects which on-disk layers to load | `.cursor/`                                                                 | `~/.cursor/`                             | n/a                                                                                    | Cloud always loads `project` / `team` / `plugins` and ignores `local.settingSources`.                   |
>
> Inline values are good for secrets that should never touch disk (per-run API keys, tenant-scoped tokens). File-based config is good for policy: hooks especially are a project boundary, not a per-run knob.
>
> ## MCP servers
>
> Agents can pick up MCP servers from several sources. Inline definitions in `Agent.create()` or `agent.send()` are the most common path. File-based and dashboard-managed configs are also supported.
>
> ### What gets loaded
>
> **Local agents** load servers from up to five sources, with first-match-wins precedence on conflicting names:
>
> 1. `mcpServers` on `agent.send()`. Fully replaces creation-time servers for that run (not merged).
> 2. `mcpServers` on `Agent.create()`. Used when no per-send override is provided.
> 3. Plugin servers, if `local.settingSources` includes `"plugins"`.
> 4. Project servers from `.cursor/mcp.json`, if `local.settingSources` includes `"project"`.
> 5. User servers from `~/.cursor/mcp.json`, if `local.settingSources` includes `"user"`.
>
> Without `local.settingSources`, only inline servers are loaded. If a local MCP server requires OAuth login, the SDK can't prompt you to sign in. It only works if you've already signed in to that server from the Cursor app, in which case the SDK reuses that saved login.
>
> **Cloud agents** load servers from:
>
> 1. `mcpServers` on `agent.send()`. Fully replaces creation-time servers for that run (not merged).
> 2. `mcpServers` on `Agent.create()`. Used when no per-send override is provided.
> 3. Your user and team MCP servers from [cursor.com/agents](https://cursor.com/agents).
>
> If an inline server doesn't include `auth` or `headers` and you've previously authorized that server URL on cursor.com/agents, runs authenticated with a personal API token reuse those OAuth tokens automatically. Service account API keys cannot fall back to user auth as they are not associated with a user.
>
> `local.settingSources` does not apply to cloud agents.
>
> ### Local
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "auto" },
>   local: { cwd: process.cwd() },
>   mcpServers: {
>     docs: {
>       type: "http",
>       url: "https://example.com/mcp",
>       auth: {
>         CLIENT_ID: "client-id",
>         scopes: ["read", "write"],
>       },
>     },
>     filesystem: {
>       type: "stdio",
>       command: "npx",
>       args: ["-y", "@modelcontextprotocol/server-filesystem", process.cwd()],
>       cwd: process.cwd(),
>     },
>   },
> });
> ```
>
> ### Cloud
>
> Cloud agents can receive authenticated MCP configs inline too. Use HTTP auth when Cursor should proxy a remote MCP through the backend. Use stdio `env` when the server runs inside the cloud VM and reads credentials from environment variables.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   cloud: {
>     repos: [{ url: "https://github.com/your-org/your-repo", startingRef: "main" }],
>   },
>   mcpServers: {
>     linear: {
>       type: "http",
>       url: "https://mcp.linear.app/sse",
>       headers: {
>         Authorization: `Bearer ${process.env.LINEAR_API_KEY!}`,
>       },
>     },
>     figma: {
>       type: "http",
>       url: "https://api.figma.com/mcp",
>       auth: {
>         CLIENT_ID: process.env.FIGMA_CLIENT_ID!,
>         CLIENT_SECRET: process.env.FIGMA_CLIENT_SECRET!,
>         scopes: ["file_content:read"],
>       },
>     },
>     github: {
>       type: "stdio",
>       command: "npx",
>       args: ["-y", "@modelcontextprotocol/server-github"],
>       env: {
>         GITHUB_TOKEN: process.env.GITHUB_TOKEN!,
>       },
>     },
>   },
> });
> ```
>
> Use `headers` for static API keys or Bearer tokens — Cursor passes them through on every request. Use `auth` for OAuth-protected servers. For cloud, Cursor runs the OAuth flow once server-side and reuses the token across runs. Locally, the SDK can't open a browser to sign you in; it only reuses tokens you've already obtained by signing in through the Cursor app.
>
> - HTTP `headers` and `auth` are handled by Cursor's backend. Sensitive fields are redacted and do not enter the VM.
> - Stdio `env` values are passed into the VM because the server runs there. Treat them like any other runtime secret.
> - OAuth for MCP servers configured on cursor.com/agents stays per-user, even for team-level servers.
>
> See [MCP](https://cursor.com/docs/mcp.md) for the full config format and [Cloud Agent capabilities](https://cursor.com/docs/cloud-agent/capabilities.md#mcp-tools) for cloud-specific behavior.
>
> ## Subagents
>
> Define named subagents that the main agent spawns via the `Agent` tool. Pass them inline:
>
> ```typescript
> const agent = await Agent.create({
>   model: { id: "composer-2.5" },
>   apiKey: process.env.CURSOR_API_KEY!,
>   local: { cwd: process.cwd() },
>   agents: {
>     "code-reviewer": {
>       description: "Expert code reviewer for quality and security.",
>       prompt: "Review code for bugs, security issues, and proven approaches.",
>       model: "inherit",
>     },
>     "test-writer": {
>       description: "Writes tests for code changes.",
>       prompt: "Write comprehensive tests for the given code.",
>     },
>   },
> });
> ```
>
> Subagents committed to the repo at `.cursor/agents/*.md` (with `name`, `description`, and optional `model` frontmatter) are also picked up. Inline definitions override file-based ones with the same name.
>
> ### Nested subagents
>
> Subagents can spawn their own subagents, within a nesting limit. When a subagent uses the `Agent` tool, the SDK hands it the same subagent executor the parent has, so a parent can delegate to a subagent that delegates further. Each level reaches the same set of named subagents and [custom tools](https://cursor.com/docs/sdk/typescript.md#custom-tools). The top-level agent and its direct subagents can launch subagents, but a subagent launched by another subagent can't launch further ones.
>
> ## Restricting the toolset
>
> `tools` allowlists the built-in tools offered to the model; `disallowedTools` removes tools and keeps the rest, including tools added to the platform after your SDK version was released. Both are local agents only for now, and neither persists on the agent: pass them again on `Agent.resume()` to keep the restriction for follow-up runs.
>
> ```typescript
> // Read-only agent: only these tools are offered.
> const reader = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   tools: ["read", "grep", "glob", "ls"],
>   local: { cwd: process.cwd() },
> });
>
> // Everything except shell access.
> const noShell = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   disallowedTools: ["shell"],
>   local: { cwd: process.cwd() },
> });
> ```
>
> - `tools: undefined` (default) offers the standard toolset for the selected model; `tools: []` offers no built-in tools, so the model can only respond with text.
> - Both fields accept the `ToolName` union: public names (`"read"`, `"edit"`, `"task"`, `"webSearch"`, ...), the capability groups `"shell"` and `"mcp"`, and raw proto tool names. Unknown names throw a `ConfigurationError` at `Agent.create()` / `Agent.resume()`.
> - Deny wins: a tool must be in `tools` (when set) and not in `disallowedTools` to be offered.
> - Disallowing `"mcp"` also removes [custom tools](https://cursor.com/docs/sdk/typescript.md#custom-tools). Disallowing `"task"` prevents [subagents](https://cursor.com/docs/sdk/typescript.md#subagents); otherwise subagents keep their own curated toolsets.
>
> ## Custom tools
>
> Custom tools let you expose your own functions to the agent without standing up a separate MCP server. Pass them on `local.customTools` and the SDK registers them as an MCP server named `custom-user-tools`. The agent discovers and calls them through the same MCP path as any other server. Deny rules and [sandbox](https://cursor.com/docs/sdk/typescript.md#sandbox-options) limits still apply, but custom tools skip interactive approval, so [sandboxed](https://cursor.com/docs/sdk/typescript.md#sandbox-options) and [auto-review](https://cursor.com/docs/sdk/typescript.md#auto-review) runs call them without prompting. Custom tools reach [subagents](https://cursor.com/docs/sdk/typescript.md#subagents) (including nested ones) too.
>
> Custom tools are local agents only. Passing `local.customTools` to a cloud agent throws a `ConfigurationError`.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   local: {
>     cwd: process.cwd(),
>     customTools: {
>       get_deployment_status: {
>         description: "Look up the current deployment status for a service.",
>         inputSchema: {
>           type: "object",
>           properties: {
>             service: { type: "string", description: "Service name" },
>           },
>           required: ["service"],
>         },
>         async execute({ service }) {
>           const res = await fetch(`https://deploys.internal/api/${service}`);
>           const body = await res.json();
>           return `Service ${service} is ${body.status} (build ${body.build}).`;
>         },
>       },
>     },
>   },
> });
>
> await agent.send("Is the checkout service deployed yet?").then((r) => r.wait());
> ```
>
> Set custom tools once on `Agent.create()` to apply them to every run, or pass `local.customTools` on a single `agent.send()` to replace them for that run.
>
> ```typescript
> await agent.send("Roll forward if the canary is healthy", {
>   local: {
>     customTools: {
>       promote_canary: {
>         description: "Promote the current canary build to production.",
>         async execute() {
>           await promoteCanary();
>           return { content: [{ type: "text", text: "Promoted." }] };
>         },
>       },
>     },
>   },
> });
> ```
>
> ### Tool definition
>
> ```typescript
> interface SDKCustomTool {
>   description?: string;
>   inputSchema?: Record<string, SDKJsonValue>;
>   execute: (
>     args: Record<string, SDKJsonValue>,
>     context: SDKCustomToolContext
>   ) => SDKCustomToolResult | Promise<SDKCustomToolResult>;
> }
>
> interface SDKCustomToolContext {
>   toolCallId?: string;
> }
> ```
>
> | Field         | Description                                                                                                                                    |
> | :------------ | :--------------------------------------------------------------------------------------------------------------------------------------------- |
> | `description` | Shown to the model so it knows when to call the tool. Defaults to an empty string.                                                             |
> | `inputSchema` | JSON Schema for the arguments. Defaults to an open object that accepts any properties.                                                         |
> | `execute`     | Your callback. Receives the parsed `args` and a `context` with the `toolCallId`. Runs in your process, so it can reach anything your code can. |
>
> ### Tool results
>
> `execute` can return a plain string, any JSON value, or a structured envelope. The map key is the tool name the model calls.
>
> ```typescript
> type SDKCustomToolResult =
>   | string
>   | SDKJsonValue
>   | {
>       content: SDKCustomToolContent[];
>       isError?: boolean;
>       structuredContent?: Record<string, SDKJsonValue>;
>     };
>
> type SDKCustomToolContent =
>   | { type: "text"; text: string }
>   | { type: "image"; data: string; mimeType?: string };
> ```
>
> - Return a string for plain text output.
> - Return any JSON value to send it back as text; objects also populate `structuredContent`.
> - Return the envelope for full control: mix text and base64 image `content`, set `isError: true` to report a failure, or attach `structuredContent` for the model to parse. Throwing from `execute` is also reported back to the agent as a tool error.
>
> ## Hooks
>
> Hooks are file-based only. There is no programmatic hook callback. Hooks are a project policy boundary, not a per-run knob.
>
> - **Local:** Add `.cursor/hooks.json` to the repo passed as `local.cwd`, or add `~/.cursor/hooks.json` for user-level hooks.
> - **Cloud:** Commit `.cursor/hooks.json` and its scripts to the repo passed in `cloud.repos`. SDK-created cloud agents load project hooks automatically. On Enterprise plans, they also run team hooks and enterprise-managed hooks.
>
> See [Hooks](https://cursor.com/docs/hooks.md) for the configuration format and [Cloud Agents hooks support](https://cursor.com/docs/cloud-agent.md#hooks-support) for cloud behavior.
>
> ## Sandbox options
>
> Local agents run with `local.sandboxOptions.enabled: false` by default. The agent can read and write the working directory, execute shell commands, and reach the network without restriction. There's no human-in-the-loop approval flow in headless SDK runs, so a sandbox-by-default would either block legitimate tool calls silently or require a callback that doesn't fit a script.
>
> When you enable the sandbox, the SDK constrains every shell tool call and shell-spawned process:
>
> - **Filesystem** — Writes are limited to the working directory (`local.cwd`) and a small set of allowed paths. Reads outside the workspace are blocked.
> - **Shell** — Commands run inside a platform sandbox (`bubblewrap` on Linux, `seatbelt` on macOS, the bundled `@cursor/sdk-<os>-<arch>` helper). Privileged operations are denied.
> - **Network** — Outbound network is denied by default. To allow specific hosts, drop a `.cursor/sandbox.json` in the workspace listing the allowed hosts. The SDK reads the same per-user policy at `~/.cursor/sandbox.json` if present.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   local: {
>     cwd: process.cwd(),
>     sandboxOptions: { enabled: true },
>   },
> });
> ```
>
> If sandboxing isn't supported on the host (older Linux without `bubblewrap`, missing helper binary), the SDK throws a `ConfigurationError` with a message that names the missing dependency. Disable `sandboxOptions.enabled` or run in cloud mode to recover.
>
> Cloud runs always execute inside an isolated VM, so `sandboxOptions` doesn't apply.
>
> ## Auto-review
>
> By default a local agent runs every tool call without restriction, since headless runs have no human to approve them. Set `local.autoReview: true` to route local tool calls through [Auto-review](https://cursor.com/docs/agent/security/run-modes.md) instead, the same classifier the IDE uses to allow or block Shell, MCP, and Fetch calls based on safety and how well each call matches the run's intent.
>
> ```typescript
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   local: {
>     cwd: process.cwd(),
>     autoReview: true,
>   },
> });
> ```
>
> Auto-review needs the classifier enabled on the connected backend; when it isn't available, runs fall back to the default behavior. Because there's no interactive approval in a headless run, a call the classifier blocks is denied rather than escalated, and the agent gets the block reason and can try another approach. Steer the classifier with a `permissions.json` `autoRun` block in the workspace, the same as in the IDE. See [permissions.json](https://cursor.com/docs/reference/permissions.md) for the format.
>
> Auto-review is local agents only. Cloud runs already execute in an isolated VM. The classifier is best-effort convenience, not a security boundary; combine it with [`sandboxOptions`](https://cursor.com/docs/sdk/typescript.md#sandbox-options) or an [allowlist](https://cursor.com/docs/agent/security/run-modes.md) for strict control.
>
> ## Artifacts
>
> List and download files from the agent's workspace.
>
> ```typescript
> interface SDKArtifact {
>   path: string;
>   sizeBytes: number;
>   updatedAt: string;
> }
> ```
>
> ```typescript
> const artifacts: SDKArtifact[] = await agent.listArtifacts();
>
> for (const artifact of artifacts) {
>   console.log(artifact.path, artifact.sizeBytes);
> }
>
> const buffer = await agent.downloadArtifact(artifacts[0].path);
> ```
>
> Artifact support is runtime-dependent. Local SDK agents currently return no artifacts and throw for `downloadArtifact`.
>
> ## Resource management
>
> Always dispose agents when done. The cleanest pattern is `await using`:
>
> ```typescript
> await using agent = await Agent.create({ /* ... */ });
> // disposed automatically when the block exits
> ```
>
> To dispose explicitly:
>
> ```typescript
> await agent[Symbol.asyncDispose]();
> ```
>
> `agent.close()` is the documented way to start disposal without awaiting. `Symbol.asyncDispose` works (`await using` is built on it) but `close()` is the path you should reach for in code that doesn't use the `await using` syntax. `agent.reload()` picks up filesystem config changes (hooks, project MCP, subagents) without disposing.
>
> ## Agent lifecycle
>
> ### Prewarm a local workspace
>
> Resolving a local workspace (rules, skills, MCP servers, ignore files) is the slowest part of a local agent's first turn, and on a large repo it can dominate it. By default that cost lands inside the first `send()`. A host that knows where its agents will run can pay it earlier with `prewarmLocalWorkspace()`:
>
> ```typescript
> import { createAgentPlatform } from "@cursor/sdk";
>
> const platform = await createAgentPlatform();
> const release = await platform.prewarmLocalWorkspace({
>   apiKey: process.env.CURSOR_API_KEY!,
>   local: { cwd: "/srv/checkout", settingSources: ["project"] },
> });
>
> // The first send() against this workspace starts immediately.
>
> await release(); // on shutdown
> ```
>
> Pass the same `AgentOptions` your agents will use; prewarming only helps sends whose workspace options match. Call the returned release function when the host shuts down.
>
> ### Reattach to an existing agent
>
> `Agent.resume(agentId)` returns a fresh handle to an agent that already exists. The runtime is auto-detected from the ID prefix (`bc-` is cloud, anything else is local), and conversation state is loaded from the cloud (cloud) or the local checkpoint store (local). This is how you continue work after a process restart, or how a different worker picks up an agent another process started.
>
> ```typescript
> const agent = await Agent.resume("bc-abc123", {
>   apiKey: process.env.CURSOR_API_KEY!,
> });
>
> const run = await agent.send("Apply the suggested fix");
> const result = await run.wait();
> ```
>
> If the run was already running when you reattached, `Agent.getRun(runId, { runtime: "cloud", agentId })` (or the local equivalent) returns a `Run` you can `stream()`, `wait()`, or `cancel()` against.
>
> ### Conversation context
>
> Local agents persist conversation state in a checkpoint store. By default this is on-disk SQLite under your home directory; swap it for JSONL or a custom backend with [`local.store`](https://cursor.com/docs/sdk/typescript.md#local-agent-stores). Each call to `agent.send()` loads the latest checkpoint for that agent and passes it to the model, so follow-ups see the same context the previous run finished with. The store survives process restarts, which means `Agent.resume(agentId)` from a brand-new process picks up where the previous one left off.
>
> Cloud agents persist state server-side. Reattaching from anywhere returns the same conversation.
>
> A few things that look like context loss but aren't:
>
> - A new `Agent.create()` always starts a fresh agent with a new `agentId`. To continue an existing conversation, capture `agent.agentId` from the first call and use `Agent.resume(agentId)` later.
> - `Agent.prompt()` creates, runs, and disposes in one shot. There's no second turn; that's the contract.
> - Inline `mcpServers` aren't persisted across `Agent.resume()` because they often carry secrets. Pass them again on resume, or use file-based MCP config.
>
> ### Dispatcher pattern
>
> A dispatcher owns a pool of agents and hands work to them as it arrives. The shape is straightforward: keep a map of `agentId` to long-lived `SDKAgent`, route incoming prompts by some key (user, repo, ticket), and `Agent.resume()` from disk if a process restart wiped the in-memory map.
>
> ```typescript
> import { Agent, type SDKAgent } from "@cursor/sdk";
>
> const agents = new Map<string, SDKAgent>();
>
> async function getAgent(key: string, savedId?: string): Promise<SDKAgent> {
>   const existing = agents.get(key);
>   if (existing) return existing;
>
>   const agent = savedId
>     ? await Agent.resume(savedId, {
>         apiKey: process.env.CURSOR_API_KEY!,
>       })
>     : await Agent.create({
>         apiKey: process.env.CURSOR_API_KEY!,
>         model: { id: "composer-2.5" },
>         local: { cwd: process.cwd() },
>       });
>
>   agents.set(key, agent);
>   return agent;
> }
>
> async function handleMessage(key: string, prompt: string, savedId?: string) {
>   const agent = await getAgent(key, savedId);
>   const run = await agent.send(prompt);
>   return run.wait();
> }
> ```
>
> Cloud SSE streams retain backlog for a window after the run starts, so a dispatcher that streams to many subscribers can call `run.stream()` from each subscriber without losing earlier events. For really long-running cloud runs, dispatchers usually fan out to `run.wait()` and let subscribers poll `run.conversation()` if they need the structured transcript.
>
> ## Local agent stores
>
> Local agents persist agent metadata, conversation checkpoints, runs, and run events to disk so that follow-ups and `Agent.resume()` survive process restarts. By default the SDK uses on-disk SQLite when the SQLite backend is available and falls back to `JsonlLocalAgentStore` otherwise. You can swap in a different backend with `local.store`.
>
> The SDK ships two backends and lets you bring your own:
>
> | Store                    | Import               | When to use                                                                                                          |
> | :----------------------- | :------------------- | :------------------------------------------------------------------------------------------------------------------- |
> | `SqliteLocalAgentStore`  | `@cursor/sdk/sqlite` | On-disk SQLite under the workspace state root.                                                                       |
> | `JsonlLocalAgentStore`   | `@cursor/sdk`        | Portable newline-delimited JSON (NDJSON) files under a directory you choose. Easy to inspect, copy, and diff.        |
> | Custom `LocalAgentStore` | Your code            | Persist to anything: in-memory, Redis, Postgres, or a hosted database. Implement the interface or compose substores. |
>
> Cloud agents persist server-side, so `local.store` applies to local agents only.
>
> ### JSONL store
>
> `JsonlLocalAgentStore` writes four NDJSON files (`agents.ndjson`, `runs.ndjson`, `run_events.ndjson`, `checkpoints.ndjson`) under the directory you pass. Construct one and pass it on `local.store`.
>
> ```typescript
> import { Agent, JsonlLocalAgentStore } from "@cursor/sdk";
>
> const store = new JsonlLocalAgentStore("/var/lib/cursor-agents");
>
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   local: { cwd: process.cwd(), store },
> });
> ```
>
> Pass the same store instance on `Agent.resume()` and on the local list and get APIs (`Agent.list`, `Agent.get`, `Agent.listRuns`, `Agent.getRun`) so they read the same data.
>
> ### Set a process-wide default
>
> To avoid threading a store through every call, set a default once with [`Cursor.configure()`](https://cursor.com/docs/sdk/typescript.md#cursorconfigure). Per-call `local.store` still wins when you pass it.
>
> ```typescript
> import { Cursor, JsonlLocalAgentStore } from "@cursor/sdk";
>
> Cursor.configure({ local: { store: new JsonlLocalAgentStore("/var/lib/cursor-agents") } });
>
> // Later calls use the configured store unless they pass their own.
> const agent = await Agent.create({
>   apiKey: process.env.CURSOR_API_KEY!,
>   model: { id: "composer-2.5" },
>   local: { cwd: process.cwd() },
> });
> ```
>
> Pass `store: null` to `Cursor.configure({ local: { store: null } })` to clear a previous default and return to the SDK's default local store selection.
>
> ### Custom stores
>
> To persist somewhere else (a shared Postgres, Redis, or an in-memory map for tests), implement `LocalAgentStore`. It's four substores, each a small CRUD surface the SDK calls:
>
> ```typescript
> interface LocalAgentStore {
>   readonly agents: LocalAgentStoreAgents;         // agent metadata rows
>   readonly checkpoints: LocalAgentStoreCheckpoints; // content-addressed conversation blobs
>   readonly runs: LocalAgentStoreRuns;             // run rows
>   readonly runEvents: LocalAgentStoreRunEvents;   // append-only run event log
> }
> ```
>
> Implement the interface directly, or build each substore separately and combine them with `composeLocalAgentStore`:
>
> ```typescript
> import { composeLocalAgentStore } from "@cursor/sdk";
>
> const store = composeLocalAgentStore({
>   agents: myAgentsTable,
>   checkpoints: myCheckpointBlobs,
>   runs: myRunsTable,
>   runEvents: myRunEventLog,
> });
> ```
>
> The substores mirror the default SQLite tables: `agents` holds one row per agent (with a slim `latestCheckpoint.rootBlobId` pointer), `checkpoints` holds the content-addressed conversation blobs those pointers reference, `runs` holds one row per run, and `runEvents` is the append-only stream log. Catalog substores paginate with an opaque `cursor` / `nextCursor`; the run event log resumes with an exclusive `afterOffset` / `nextOffset`. See the exported `LocalAgentStore`, `LocalAgentDocument`, `LocalAgentRunDocument`, and related types for the exact shapes.
>
> ## Configuration reference
>
> ### AgentOptions
>
> | Property          | Type                              | Default                                                             | Description                                                                                                                                                                                          |
> | :---------------- | :-------------------------------- | :------------------------------------------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `model`           | `ModelSelection`                  | Required for local; cloud falls back to the server-resolved default | Model to use. See [`ModelSelection`](https://cursor.com/docs/sdk/typescript.md#modelselection).                                                                                                      |
> | `apiKey`          | `string`                          | `CURSOR_API_KEY` env                                                | User API key or service account key. Team Admin keys are not yet supported.                                                                                                                          |
> | `name`            | `string`                          | Auto-generated                                                      | Human-readable agent name returned as `name` in `Agent.list()` / `Agent.get()`.                                                                                                                      |
> | `local`           | `LocalAgentOptions`               |                                                                     | Local agent config. See [`LocalAgentOptions`](https://cursor.com/docs/sdk/typescript.md#localagentoptions).                                                                                          |
> | `cloud`           | `CloudAgentOptions`               |                                                                     | Cloud agent config.                                                                                                                                                                                  |
> | `mcpServers`      | `Record<string, McpServerConfig>` |                                                                     | Inline MCP server definitions.                                                                                                                                                                       |
> | `agents`          | `Record<string, AgentDefinition>` |                                                                     | Subagent definitions.                                                                                                                                                                                |
> | `tools`           | `ToolName[]`                      | Default toolset                                                     | [Restrict the toolset](https://cursor.com/docs/sdk/typescript.md#restricting-the-toolset): only the listed built-in tools are offered to the model. `[]` means no built-in tools. Local agents only. |
> | `disallowedTools` | `ToolName[]`                      |                                                                     | [Remove tools](https://cursor.com/docs/sdk/typescript.md#restricting-the-toolset) from the toolset; everything else stays available. Deny wins when combined with `tools`. Local agents only.        |
> | `agentId`         | `string`                          | Auto-generated                                                      | Durable agent ID. Pass to keep a stable ID across invocations.                                                                                                                                       |
> | `idempotencyKey`  | `string`                          | Auto-generated for cloud                                            | Optional client-generated idempotency key.                                                                                                                                                           |
> | `mode`            | `"agent" \| "plan"`               | `"agent"`                                                           | Initial conversation mode for the agent's first run. See [Conversation mode](https://cursor.com/docs/sdk/typescript.md#conversation-mode).                                                           |
>
> ### LocalAgentOptions
>
> Config for local agents, passed as `local` on `Agent.create()`. Also exported as a standalone type for `Partial<LocalAgentOptions>`.
>
> | Property             | Type                            | Default              | Description                                                                                                                                              |
> | :------------------- | :------------------------------ | :------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `cwd`                | `string`                        |                      | Primary working directory for the default shell and agent-store scoping.                                                                                 |
> | `dirs`               | `string[]`                      |                      | Additional workspace folders for multi-root setups. Merged with `cwd` (duplicates dropped) so rules, skills, and workspace context load from every path. |
> | `settingSources`     | `SettingSource[]`               |                      | Ambient settings layers to load: `"project"`, `"user"`, `"team"`, `"mdm"`, `"plugins"`, or `"all"`.                                                      |
> | `sandboxOptions`     | `{ enabled: boolean }`          | `{ enabled: false }` | [Sandbox](https://cursor.com/docs/sdk/typescript.md#sandbox-options) config.                                                                             |
> | `autoReview`         | `boolean`                       | `false`              | Route local tool calls through [Auto-review](https://cursor.com/docs/sdk/typescript.md#auto-review).                                                     |
> | `customTools`        | `Record<string, SDKCustomTool>` |                      | [Custom tools](https://cursor.com/docs/sdk/typescript.md#custom-tools) exposed as the `custom-user-tools` MCP server.                                    |
> | `store`              | `LocalAgentStore`               | SDK default store    | [Local agent store](https://cursor.com/docs/sdk/typescript.md#local-agent-stores) backing persistence.                                                   |
> | `enableAgentRetries` | `boolean`                       | `true`               | Enable transport and stall auto-retry for local agent runs. Set `false` to surface transport errors on the first failure.                                |
>
> ### CloudAgentOptions
>
> | Property                | Type                                                                                                        | Default                                                | Description                                                                                                                                                                                                                                                                                                                                            |
> | :---------------------- | :---------------------------------------------------------------------------------------------------------- | :----------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `env`                   | `{ type: "cloud"; name?: string } \| { type: "pool"; name?: string } \| { type: "machine"; name?: string }` | `{ type: "cloud" }`                                    | Execution environment target. `cloud` uses Cursor-hosted VMs; set `name` to use a saved Cursor-hosted environment. `pool` and `machine` route to self-hosted workers you run. Omit `repos` and leave `env` at the default for a no-repo agent with an empty workspace. Named Cursor-hosted environments and explicit `repos` are mutually exclusive.   |
> | `repos`                 | `Array<{ url: string; startingRef?: string; prUrl?: string }>`                                              |                                                        | Repositories to clone into the VM. Pass one entry for a single-repo agent, or up to 20 for a multi-repo agent. Omit or pass `[]` for a [no-repo agent](https://cursor.com/docs/sdk/typescript.md#no-repo-cloud-agents). Mutually exclusive with a named `env.name` for Cursor-hosted environments. Pass `prUrl` to attach the agent to an existing PR. |
> | `workOnCurrentBranch`   | `boolean`                                                                                                   | `false`                                                | Push commits to the existing branch instead of a new one.                                                                                                                                                                                                                                                                                              |
> | `autoCreatePR`          | `boolean`                                                                                                   | `false`                                                | Open a PR when the run finishes.                                                                                                                                                                                                                                                                                                                       |
> | `openAsCursorGithubApp` | `boolean`                                                                                                   | `true` for service-account keys, `false` for user keys | Open PRs as the Cursor GitHub App instead of the API key's owner. The resolved value is echoed on create, get, and list.                                                                                                                                                                                                                               |
> | `skipReviewerRequest`   | `boolean`                                                                                                   | `false`                                                | Skip requesting the calling user as a reviewer on the PR.                                                                                                                                                                                                                                                                                              |
> | `envVars`               | `Record<string, string>`                                                                                    |                                                        | Session-scoped environment variables for cloud agents.                                                                                                                                                                                                                                                                                                 |
> | `metadata`              | `Record<string, string>`                                                                                    |                                                        | Caller-owned string tags persisted on the cloud agent. See [Agent metadata](https://cursor.com/docs/sdk/typescript.md#agent-metadata).                                                                                                                                                                                                                 |
>
> ### AgentDefinition
>
> | Property      | Type                                               | Default     | Description                                                                                     |
> | :------------ | :------------------------------------------------- | :---------- | :---------------------------------------------------------------------------------------------- |
> | `description` | `string`                                           | *required*  | When to use this subagent. Shown to the parent agent so it knows when to spawn.                 |
> | `prompt`      | `string`                                           | *required*  | System prompt for the subagent.                                                                 |
> | `model`       | `ModelSelection \| "inherit"`                      | `"inherit"` | Model override. Pass `"inherit"` to use the parent's selection.                                 |
> | `mcpServers`  | `Array<string \| Record<string, McpServerConfig>>` |             | MCP servers available to this subagent. Names reference servers from the parent's `mcpServers`. |
>
> ### ModelSelection
>
> ```typescript
> interface ModelSelection {
>   id: string;
>   params?: ModelParameterValue[];
> }
>
> interface ModelParameterValue {
>   id: string;
>   value: string;
> }
> ```
>
> `id` is the model identifier (for example, `"composer-2.5"` or `"auto-smart"`). `params` carries per-model parameters such as reasoning effort or Router's `optimize_for`. Use [`Cursor.models.list()`](https://cursor.com/docs/sdk/typescript.md#cursormodelslist) to discover valid ids, parameter definitions, and preset variants for your account. See [Cursor Router](https://cursor.com/docs/sdk/typescript.md#cursor-router) for the Router selection contract.
>
> ### McpServerConfig
>
> ```typescript
> type McpServerConfig =
>   // stdio
>   | {
>       type?: "stdio";
>       command: string;
>       args?: string[];
>       env?: Record<string, string>;
>       cwd?: string;       // local only; cloud rejects this field
>     }
>   // HTTP / SSE
>   | {
>       type?: "http" | "sse";
>       url: string;
>       headers?: Record<string, string>;   // passed through; Authorization here works
>       auth?: {
>         CLIENT_ID: string;
>         CLIENT_SECRET?: string;
>         scopes?: string[];
>       };
>     };
> ```
>
> For HTTP servers running in the cloud, `headers` and `auth` are handled by Cursor's backend. Sensitive fields are redacted before the VM sees them. For stdio servers in the cloud, `env` values are passed into the VM (treat them like any runtime secret).
>
> ### SDKUserMessage
>
> ```typescript
> interface SDKUserMessage {
>   text: string;
>   images?: SDKImage[];
> }
> ```
>
> The structured form of `agent.send()`'s message argument. Use it to send images alongside text.
>
> ### SDKImage
>
> ```typescript
> type SDKImage =
>   | { url: string; dimension?: SDKImageDimension }
>   | { data: string; mimeType: string; dimension?: SDKImageDimension };
>
> interface SDKImageDimension {
>   width: number;
>   height: number;
> }
> ```
>
> Pass either a remote `url` or base64 `data` with a `mimeType`.
>
> ### SettingSource
>
> ```typescript
> type SettingSource =
>   | "project"
>   | "user"
>   | "team"
>   | "mdm"
>   | "plugins"
>   | "all";
> ```
>
> Controls which on-disk settings layers a local agent loads. Cloud agents always load `project` / `team` / `plugins` and ignore this field.
>
> | Value       | Source                                  |
> | :---------- | :-------------------------------------- |
> | `"project"` | `.cursor/` in the workspace             |
> | `"user"`    | `~/.cursor/`                            |
> | `"team"`    | Team settings synced from the dashboard |
> | `"mdm"`     | MDM-managed enterprise settings         |
> | `"plugins"` | Plugin-provided settings                |
> | `"all"`     | Shorthand for all of the above          |
>
> ### ListResult
>
> ```typescript
> interface ListResult<T> {
>   items: T[];
>   nextCursor?: string;
> }
> ```
>
> Returned by `Agent.list()` and `Agent.listRuns()`. `nextCursor` is absent when there are no more pages.
>
> ## Errors
>
> All SDK errors extend `CursorSdkError` (re-exported as `CursorAgentError` for backwards compatibility). Use `isRetryable` to drive retry logic, and `code` / `status` / `requestId` for diagnostics.
>
> ```typescript
> class CursorSdkError extends Error {
>   readonly isRetryable: boolean;
>   readonly code?: string;       // stable SDK / backend code
>   readonly status?: number;     // HTTP status if available
>   readonly cause?: unknown;     // wrapped underlying error
>   readonly endpoint?: string;
>   readonly requestId?: string;
>   readonly operation?: string;  // SDK operation that produced the error
> }
> ```
>
> | Error class                    | Typical message                                               | Likely cause                                                                                                               | Recommended fix                                                                                                                                                                                                 |
> | :----------------------------- | :------------------------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `AuthenticationError`          | "Invalid API key"                                             | Missing or wrong `CURSOR_API_KEY`, expired token, or admin disabled the key.                                               | Generate a new key from [API Keys](https://cursor.com/dashboard/api) (user) or [Team settings](https://cursor.com/dashboard/team-settings) (service account). Confirm the key has permission for the operation. |
> | `RateLimitError`               | "Rate limit exceeded" or "Usage limit exceeded"               | Burst limit or monthly usage cap.                                                                                          | Back off using exponential delay (the SDK reports `isRetryable: true` for transient cases). For monthly cap, raise the plan's [usage limit](https://cursor.com/help/models-and-usage/usage-limits.md).          |
> | `ConfigurationError`           | "Bad model name", "API key not supported", "File unsupported" | Invalid `model.id`, missing required `params`, unsupported file in a tool call, or an admin policy blocking the request.   | Call `Cursor.models.list()` to confirm the id and params. Check repo / file paths exist.                                                                                                                        |
> | `AgentBusyError`               | "Agent is busy"                                               | Sending a follow-up while the same cloud agent already has a run in `CREATING` or `RUNNING` state.                         | Wait for the active run to finish, cancel it, or poll `Agent.listRuns()` before sending again.                                                                                                                  |
> | `IntegrationNotConnectedError` | "\[provider] integration is not connected"                    | Creating a cloud agent for a repo whose SCM provider isn't connected to your Cursor team.                                  | Open `error.helpUrl` to reconnect the provider, then retry.                                                                                                                                                     |
> | `NetworkError`                 | "Service unavailable", "Timeout"                              | Transient backend issue, network partition, or deadline exceeded.                                                          | Retry with backoff. Inspect `error.requestId` if you need to file a support ticket.                                                                                                                             |
> | `UnsupportedRunOperationError` | "Operation "stream" is not supported on this runtime"         | Calling a `Run` method the current runtime can't satisfy (e.g. streaming on a re-fetched local run that already finished). | Guard with `run.supports(operation)` / `run.unsupportedReason(operation)` first.                                                                                                                                |
> | `AgentNotFoundError`           | "Agent not found"                                             | The requested agent does not exist or is not visible under the resolved local workspace.                                   | Check the agent ID, `cwd`, and `local.store`.                                                                                                                                                                   |
> | `UnknownAgentError`            | Server-defined message                                        | Unclassified backend or runtime error.                                                                                     | Inspect `error.code` and `error.cause` for the underlying detail.                                                                                                                                               |
>
> ### Check error.helpUrl
>
> Some errors carry a one-click resolution link. The most common is
> `IntegrationNotConnectedError`, but more error types may add `helpUrl` over
> time. When you catch an error, log `error.helpUrl` if present and surface it
> to the user.
>
> ### IntegrationNotConnectedError
>
> ```typescript
> class IntegrationNotConnectedError extends ConfigurationError {
>   readonly provider: string;   // e.g. "github", "gitlab", "azuredevops"
>   readonly helpUrl: string;    // dashboard link to reconnect
> }
> ```
>
> The default error message doesn't include `helpUrl`, so log it explicitly:
>
> ```typescript
> import { Agent, IntegrationNotConnectedError } from "@cursor/sdk";
>
> try {
>   await Agent.create({
>     apiKey: process.env.CURSOR_API_KEY!,
>     cloud: {
>       repos: [{ url: "https://github.com/your-org/private-repo" }],
>     },
>   });
> } catch (err) {
>   if (err instanceof IntegrationNotConnectedError) {
>     console.error(err.provider, err.helpUrl);
>   }
> }
> ```
>
> ### AgentBusyError
>
> ```typescript
> class AgentBusyError extends CursorAgentError {}
> ```
>
> `isRetryable` is `false` for `agent_busy`. Retrying immediately will keep failing until the active run reaches a terminal status or you cancel it. Other `409` responses, such as `agent_archived`, throw `ConfigurationError` instead.
>
> Wait for the active run to finish, cancel it with `run.cancel()`, or poll `Agent.listRuns()` before sending again:
>
> ```typescript
> import { Agent, AgentBusyError } from "@cursor/sdk";
>
> const agent = await Agent.resume("bc-00000000-0000-0000-0000-000000000001");
>
> try {
>   await agent.send({ text: "Also add tests for the auth middleware." });
> } catch (err) {
>   if (err instanceof AgentBusyError) {
>     const runs = await Agent.listRuns(agent.agentId, { runtime: "cloud", limit: 1 });
>     const active = runs.items[0];
>     if (active?.status === "running") {
>       await active.cancel();
>     }
>     await agent.send({ text: "Also add tests for the auth middleware." });
>     return;
>   }
>   throw err;
> }
> ```
>
> Local agents do not return `agent_busy`. Use `send({ local: { force: true } })` to expire a stuck local run before starting a new one.
>
> ### UnsupportedRunOperationError
>
> ```typescript
> class UnsupportedRunOperationError extends ConfigurationError {
>   readonly operation: RunOperation;
> }
> ```
>
> Thrown when a `Run` operation isn't available on the current runtime. Use `run.supports(operation)` and `run.unsupportedReason(operation)` to check before calling.
>
> ## Known limitations
>
> - Inline `mcpServers` are not persisted across `Agent.resume()`. Pass them again on resume if needed.
> - Custom tools (`local.customTools`), Auto-review (`local.autoReview`), custom stores (`local.store`), and toolset restrictions (`tools`, `disallowedTools`) are local agents only. Cloud agents reject `local.customTools` and persist server-side.
> - `tools` and `disallowedTools` are not persisted on the agent. Pass them again on `Agent.resume()` to keep the restriction.
> - Artifact download is not implemented for local agents (`agent.listArtifacts()` returns an empty list and `agent.downloadArtifact()` throws).
> - `local.settingSources` (and the file-based MCP / subagent paths it gates) does not apply to cloud agents. Cloud always loads `project` / `team` / `plugins`.
> - Hooks are file-based only (`.cursor/hooks.json`). No programmatic callbacks.
> - The SDK doesn't auto-discover credentials from a local Cursor app installation. Set `CURSOR_API_KEY` (or pass `apiKey`) explicitly, or mint a key with [`Cursor.auth.login()`](https://cursor.com/docs/sdk/typescript.md#cursorauth).
> - Local mode requires Node.js 22.13 or later and platform sandbox-helper support. The default store falls back to `JsonlLocalAgentStore` when the SQLite backend isn't available.
>
>

### Source: Python SDK

> # Cursor Python SDK
>
> The `cursor-sdk` package lets you call Cursor's agent from your own Python code. The same agent that runs in the Cursor IDE, CLI, and web app is scriptable from Python with sync and async clients, typed dataclasses, and ordinary iteration for streams and pages. Run the `/sdk` skill inside Cursor to get started.
>
> For the REST API, see the [Cloud Agents API](https://cursor.com/docs/cloud-agent/api/endpoints.md). For other languages, see the [SDK Bridge](https://cursor.com/docs/sdk/bridge.md).
>
> ## Overview
>
> The SDK wraps local and cloud runtimes behind one interface. You write the same code regardless of where the agent runs.
>
> | Runtime                   | What it does                                                          | When to use                                                                                                                |
> | :------------------------ | :-------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------- |
> | **Local**                 | Runs the agent against local files on disk.                           | Dev scripts and CI checks against a working tree.                                                                          |
> | **Cloud (Cursor-hosted)** | Runs in an isolated VM with your repo cloned in. Cursor runs the VMs. | When the caller doesn't have the repo, you want many agents in parallel, or runs need to survive the caller disconnecting. |
>
> Set the runtime by passing `local` or `cloud` to `Agent.create()`.
>
> ## Authentication
>
> Set `CURSOR_API_KEY` or pass `api_key` before creating an agent.
>
> The SDK accepts user API keys and service account API keys for both local and cloud runs. Team Admin API keys are not yet supported.
>
> - **User API key** from [Cursor Dashboard -> API Keys](https://cursor.com/dashboard/api)
> - **Service account API key** from [Team settings](https://cursor.com/dashboard/team-settings). See [Service accounts](https://cursor.com/docs/account/enterprise/service-accounts.md)
>
> ```bash
> export CURSOR_API_KEY="your-key"
> ```
>
> ## Usage and billing
>
> SDK runs follow the same pricing, request pools, and Privacy Mode rules as runs from the IDE and Cloud Agents. Spend shows up in your team's [usage dashboard](https://cursor.com/dashboard/usage) under the SDK tag.
>
> To read per-run token counts in code, see [Token usage](https://cursor.com/docs/sdk/python.md#token-usage). To fetch billed usage and dollar cost for an agent's runs, see [`agent.get_usage()`](https://cursor.com/docs/sdk/python.md#agentget_usage).
>
> ## Core concepts
>
> | Concept          | Description                                                                                                                      |
> | :--------------- | :------------------------------------------------------------------------------------------------------------------------------- |
> | **Agent**        | Durable handle that holds conversation state, workspace config, model selection, and settings. Survives across multiple prompts. |
> | **Run**          | One prompt submission. Owns its own stream, status, result, conversation, and cancellation.                                      |
> | **SDKMessage**   | Typed stream message yielded during a run. Same shape across local and cloud runtimes.                                           |
> | **CursorClient** | Explicit client for lifecycle control, custom HTTP options, or multiple workspaces in one process. `Client` is an alias.         |
> | **AsyncClient**  | Async-mirror client. Required for all async operations.                                                                          |
>
> ## Installation
>
> ```bash
> pip install cursor-sdk
> ```
>
> Requires Python 3.10 or later.
>
> ## Quick start
>
> ```python
> import os
>
> from cursor_sdk import Agent, LocalAgentOptions
>
> with Agent.create(
>     model="composer-2.5",
>     api_key="crsr_key",
>     local=LocalAgentOptions(cwd=os.getcwd()),
> ) as agent:
>     print(agent.send("Summarize what this repository does").text())
> ```
>
> [Stream events](https://cursor.com/docs/sdk/python.md#stream-events) shows how to extract assistant text, handle tool calls, and read run state. For a one-shot prompt (create, run, finish), see [`Agent.prompt()`](https://cursor.com/docs/sdk/python.md#agentprompt).
>
> ### Cloud quick start
>
> The Python SDK has native support for Cursor's cloud agents. You can list connected repositories, start an agent against one of them, wait for the run, and review the final result.
>
> ```python
> from cursor_sdk import Agent, CloudAgentOptions, CloudRepository
>
> with Agent.create(
>     model="composer-2.5",
>     api_key="crsr_key",
>     cloud=CloudAgentOptions(
>         repos=[CloudRepository(url="https://github.com/your-org/your-repo", starting_ref="main")],
>         auto_create_pr=True,
>     ),
> ) as agent:
>     print(agent.send("Add structured logging to the auth middleware").text())
> ```
>
> Cloud agents started by the SDK are filtered out of the default agent list. To view them in Cursor Web or the Cursor agents window, click **Filter > Source > SDK**.
>
> ## Async usage
>
> The async client mirrors the sync surface and is recommended for servers, bots, and concurrent agent orchestration. `AsyncAgent`, `AsyncClient`, `AsyncRun`, and `AsyncCursor` are exported from both `cursor_sdk` and `cursor_sdk.asyncio`.
>
> ```python
> import asyncio
> import os
>
> from cursor_sdk import AsyncClient, LocalAgentOptions
>
> async def main():
>     async with await AsyncClient.launch_bridge(workspace=os.getcwd()) as client:
>         async with await client.agents.create(
>             model="composer-2.5",
>             api_key="crsr_key",
>             local=LocalAgentOptions(cwd=os.getcwd()),
>         ) as agent:
>             run = await agent.send("Summarize what this repository does")
>             print(await run.text())
>
> asyncio.run(main())
> ```
>
> There is no global async default client. Instantiate `AsyncClient` explicitly, or use `AsyncClient.launch_bridge(...)` as an async context manager, so each event loop owns its own client. Do not mix sync and async clients in the same code path.
>
> Direct `AsyncAgent` class methods require `client=`. Use
> `await client.agents.create(...)` or
> `await AsyncAgent.create(..., client=client)`.
>
> | Sync                      | Async                               |
> | :------------------------ | :---------------------------------- |
> | `CursorClient` / `Client` | `AsyncClient` / `AsyncCursorClient` |
> | `Agent`                   | `AsyncAgent`                        |
> | `Run`                     | `AsyncRun`                          |
> | `Cursor`                  | `AsyncCursor`                       |
> | `ListResult`              | `AsyncListResult`                   |
> | `DefaultHttpxClient`      | `DefaultAsyncHttpxClient`           |
>
> ## Creating agents
>
> `Agent.create()` validates options and returns a handle immediately. Pass either `local` or `cloud` to pick a runtime.
>
> ```python
> from cursor_sdk import Agent, CloudAgentOptions, CloudRepository, LocalAgentOptions
>
> agent = Agent.create(
>     model="composer-2.5",
>     local=LocalAgentOptions(cwd="."),
> )
>
> cloud_agent = Agent.create(
>     model="composer-2.5",
>     cloud=CloudAgentOptions(
>         repos=[CloudRepository(url="https://github.com/your-org/your-repo", starting_ref="main")],
>         auto_create_pr=True,
>     ),
> )
> ```
>
> `agent.agent_id` is populated immediately. Local agents get an `agent-<uuid>` ID; cloud agents get a `bc-<uuid>` ID. `agent.model` is a typed `ModelSelection`, so `agent.model.id` and `agent.model.params` work directly.
>
> Cloud agents started by the SDK are filtered out of the default agent list. To
> view them in Cursor Web or a Cursor agent window, click **Filter > Source > SDK**.
>
> ### No-repo cloud agents
>
> Cloud agents can run on an empty VM with no repository. Pass `cloud` with an empty `repos` list, or omit `repos` entirely. Omitting `cloud` selects the local runtime instead.
>
> ```python
> from cursor_sdk import Agent, CloudAgentOptions
>
> with Agent.create(cloud=CloudAgentOptions(repos=[])) as agent:
>     run = agent.send("Research the top 3 Python testing frameworks and summarize.")
>     print(run.wait().result)
> ```
>
> No-repo agents must be enabled for your account or team. Repository-scoped API keys can't create them; use an unrestricted service account key or a user API key instead.
>
> ### Session environment variables
>
> For cloud agents, pass `env_vars` when a run needs short-lived credentials or other values that should live only with that agent.
>
> ```python
> import os
>
> agent = Agent.create(
>     model="composer-2.5",
>     cloud=CloudAgentOptions(
>         repos=[CloudRepository(url="https://github.com/your-org/your-repo")],
>         env_vars={
>             "STAGING_API_TOKEN": os.environ["STAGING_API_TOKEN"],
>         },
>     ),
> )
> ```
>
> These values are encrypted at rest, injected into the cloud agent's shell, and deleted with the agent. `env_vars` can't be used with a caller-supplied `agent_id`; omit `agent_id` and read the server-minted ID from `agent.agent_id`. Variable names can't start with `CURSOR_`.
>
> For values that should only exist during a single run, pass them on `agent.send()` instead. See [Per-run environment variables](https://cursor.com/docs/sdk/python.md#per-run-environment-variables).
>
> ### Agent metadata
>
> Attach your own identifiers to a cloud agent when you create it. Metadata can link an agent to a user, tenant, workflow, or ticket in your system, and is read back on `SDKAgentInfo.metadata` from `client.agents.get()` and `client.agents.list()`. These tags are not the in-VM [agent metadata](https://cursor.com/docs/cloud-agent/metadata.md) API, which exposes the current run's id, owner, turn, and workspace from inside the VM.
>
> ```python
> from cursor_sdk import Agent, CloudAgentOptions, CloudRepository
>
> with Agent.create(
>     model="composer-2.5",
>     cloud=CloudAgentOptions(
>         repos=[CloudRepository(url="https://github.com/your-org/your-repo")],
>         metadata={
>             "end_user_id": "user-123",
>             "ticket_id": "ENG-456",
>         },
>     ),
> ) as agent:
>     print(agent.agent_id)
> ```
>
> Metadata is available for cloud agents at creation time. You can attach up to 50 key-value pairs. Keys must be non-empty and no more than 255 characters. Values must be strings no larger than 4096 bytes. Empty string values are allowed, and an empty mapping is treated as no metadata.
>
> If metadata isn't enabled for the API key's account, creating an agent with a
> non-empty map returns `403 feature_unavailable`.
>
> ### Model parameters
>
> Use `ModelSelection.params` to pass per-model options such as reasoning effort or Cursor Router's `optimize_for`. Parameter IDs and values vary by model. Use [`Cursor.models.list()`](https://cursor.com/docs/sdk/python.md#the-cursor-namespace) to discover supported parameters and preset variants for your account.
>
> ```python
> from cursor_sdk import Agent, LocalAgentOptions, ModelParameterValue, ModelSelection
>
> agent = Agent.create(
>     model=ModelSelection(
>         id="composer-2.5",
>         params=[ModelParameterValue(id="fast", value="true")],
>     ),
>     local=LocalAgentOptions(cwd="."),
> )
> ```
>
> Use [`Cursor.models.list()`](https://cursor.com/docs/sdk/python.md#the-cursor-namespace) to discover the parameter IDs and preset variants for a given model. See [Cursor Router](https://cursor.com/docs/sdk/python.md#cursor-router) for the `auto-smart` selection contract.
>
> ### Cursor Router
>
> [Cursor Router](https://cursor.com/docs/cursor-router.md) selects a model for each Auto request. In the SDK, Router is the `auto-smart` model with an `optimize_for` parameter. It is available on Teams and Enterprise. Enterprise admins must enable Router for the team before `auto-smart` appears in the catalog.
>
> The Cursor SDK is an agent SDK, not a standalone model-inference or chat-completions API. Router picks models for Cursor agent runs that can reason over a workspace, call tools, run commands, and edit files. Cursor does not currently document a raw Router endpoint for arbitrary model calls.
>
> #### Select Cost, Balance, or Intelligence
>
> Pass `auto-smart` and set `optimize_for` explicitly:
>
> | Product label | SDK value      |
> | :------------ | :------------- |
> | Cost          | `cost`         |
> | Balance       | `balanced`     |
> | Intelligence  | `intelligence` |
>
> Use **Balance** in product copy. Use `balanced` only as the SDK wire value.
>
> ```python
> import os
>
> from cursor_sdk import Agent, LocalAgentOptions, ModelParameterValue, ModelSelection
>
> with Agent.create(
>     model=ModelSelection(
>         id="auto-smart",
>         params=[ModelParameterValue(id="optimize_for", value="balanced")],
>     ),
>     local=LocalAgentOptions(cwd=os.getcwd()),
> ) as agent:
>     run = agent.send("Find and fix the failing authentication test")
>     result = run.wait()
>
>     print(result.status)
> ```
>
> Always pass `optimize_for`. Do not omit it and do not send a legacy `default` value; discovery through the catalog is the supported contract.
>
> #### Discover Router in the model catalog
>
> `Cursor.models.list()` returns the models, parameter definitions, and preset variants available to the API key's current account and team. Cursor Router appears as `auto-smart` when Router is available. Team administrators can disable Router or restrict which optimization modes members may select.
>
> Treat the catalog as the source of truth before hard-coding a selection:
>
> ```python
> from cursor_sdk import Cursor, ModelParameterValue, ModelSelection
>
> models = Cursor.models.list()
> router = next((model for model in models if model.id == "auto-smart"), None)
> optimize_for = next(
>     (
>         parameter
>         for parameter in (router.parameters if router else [])
>         if parameter.id == "optimize_for"
>     ),
>     None,
> )
>
> if router is None or optimize_for is None:
>     raise RuntimeError(
>         "Cursor Router is not available for this API key. "
>         "Verify that Router is enabled for the key's team."
>     )
>
> requested_mode = "balanced"
> allowed_values = {entry.value for entry in optimize_for.values}
>
> if requested_mode not in allowed_values:
>     raise RuntimeError(
>         f'Router mode "{requested_mode}" is not enabled for this team.'
>     )
>
> model = ModelSelection(
>     id=router.id,
>     params=[ModelParameterValue(id=optimize_for.id, value=requested_mode)],
> )
> ```
>
> #### Switch modes per run
>
> Override the model on `agent.send()` to change Router mode for a run:
>
> ```python
> from cursor_sdk import ModelParameterValue, ModelSelection, SendOptions
>
> run = agent.send(
>     "Handle this complex migration",
>     SendOptions(
>         model=ModelSelection(
>             id="auto-smart",
>             params=[ModelParameterValue(id="optimize_for", value="intelligence")],
>         ),
>     ),
> )
> ```
>
> Per-run model overrides are sticky. Later sends without an override keep using the new selection. See [Per-run model override](https://cursor.com/docs/sdk/python.md#per-run-model-override).
>
> #### Model ids: `auto-smart`, `auto`, and `default`
>
> | Selection                                     | Meaning                                                                                                                                     |
> | :-------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------ |
> | `auto-smart` with `optimize_for`              | Cursor Router. Use this when you want Cost, Balance, or Intelligence.                                                                       |
> | `ModelSelection(id="auto")`                   | Server-selected Auto fallback when a specific model is missing from the catalog. Prefer `auto-smart` when you need an explicit Router mode. |
> | Omitting `optimize_for`, or sending `default` | Not a supported Router contract. Always discover allowed values and pass `cost`, `balanced`, or `intelligence`.                             |
>
> #### Billing and routing pool
>
> - All Auto modes bill at the list price of the model each request is routed to.
> - The underlying model can change between requests. Prefer a fixed model id when you need reproducible comparisons.
> - Enterprise model allowlists shape the routing pool. Blocking required models can disable Router.
>
> For current rates and the routing pool, see [Cursor Router](https://cursor.com/docs/cursor-router.md) and [Auto modes](https://cursor.com/docs/models-and-pricing.md#auto-modes).
>
> #### Troubleshooting missing Router
>
> If `auto-smart` is missing or an optimization mode is rejected:
>
> 1. Call `Cursor.models.list()`.
> 2. Confirm `auto-smart` is in the result.
> 3. Confirm `optimize_for` includes the value you want (`cost`, `balanced`, or `intelligence`).
> 4. Confirm Router is enabled for the team tied to the API key.
> 5. If you belong to multiple teams, confirm the key is operating in the intended team context.
> 6. Check team model-access policy if Router is unavailable or cannot choose a valid underlying model.
>
> ### Raw dictionaries
>
> Typed dataclasses are preferred for application code because IDE autocomplete and type checking work better. The SDK also accepts plain dictionaries for short scripts or externally supplied JSON. Snake-case keys are normalized.
>
> ```python
> from cursor_sdk import Agent
>
> with Agent.create(
>     {
>         "api_key": "crsr_key",
>         "model": {"id": "composer-2.5"},
>         "local": {"cwd": "."},
>     }
> ) as agent:
>     ...
> ```
>
> ## Agent
>
> The handle returned by `Agent.create()`, `Agent.resume()`, `client.agents.create()`, and `client.agents.resume()`.
>
> ```python
> class Agent:
>     agent_id: str
>     model: ModelSelection | None
>     client: CursorClient
>
>     def send(
>         self,
>         message: str | Mapping[str, Any] | UserMessage,
>         options: SendOptions | Mapping[str, Any] | None = None,
>         *,
>         idempotency_key: str | None = None,
>     ) -> Run: ...
>
>     def reload(self) -> None: ...
>     def close(self) -> None: ...
>
>     def list_messages(
>         self, options: Mapping[str, Any] | None = None
>     ) -> list[AgentMessage]: ...
>     def list_artifacts(self) -> list[SDKArtifact]: ...
>     def download_artifact(self, path: str) -> bytes: ...
>     def get_usage(self, *, run_id: str | None = None) -> AgentUsage: ...
>
>     def archive(self, options: Mapping[str, Any] | None = None) -> None: ...
>     def unarchive(self, options: Mapping[str, Any] | None = None) -> None: ...
>     def delete(self, options: Mapping[str, Any] | None = None) -> None: ...
> ```
>
> | Member                             | Description                                                                           |
> | :--------------------------------- | :------------------------------------------------------------------------------------ |
> | `agent_id`                         | Stable agent identifier. `agent-<uuid>` for local, `bc-<uuid>` for cloud.             |
> | `model`                            | Current typed model selection. Updates after a successful send with a model override. |
> | `send`                             | Start a new run with the given prompt. Returns a `Run` handle.                        |
> | `reload`                           | Re-read filesystem config (hooks, project MCP, subagents) without disposing.          |
> | `close`                            | Close the agent and release resources.                                                |
> | `list_messages`                    | List message history for the agent.                                                   |
> | `list_artifacts`                   | List files produced by the agent (cloud only; local returns empty).                   |
> | `download_artifact`                | Download a file by path (cloud only; local raises).                                   |
> | `get_usage`                        | Fetch billed token usage and dollar cost for the agent.                               |
> | `archive` / `unarchive` / `delete` | Manage cloud agent lifecycle.                                                         |
>
> Use a context manager for automatic cleanup:
>
> ```python
> with Agent.create(model="composer-2.5", local=LocalAgentOptions(cwd=".")) as agent:
>     print(agent.send("Explain this repository").text())
> ```
>
> When you use the sync `Agent.*` or `Cursor.*` helpers without passing `client=`, the SDK starts or reuses a module-level default client. It is closed automatically at process exit, and you can close it explicitly:
>
> ```python
> from cursor_sdk import close_default_client
>
> close_default_client()
> ```
>
> ### Agent.prompt()
>
> ```python
> Agent.prompt(
>     message: str | Mapping[str, Any] | UserMessage,
>     options: AgentOptions | Mapping[str, Any] | None = None,
>     *,
>     client: CursorClient | None = None,
> ) -> RunResult
> ```
>
> One-shot convenience: creates an agent, sends a single prompt, waits for the run to finish, and disposes.
>
> ```python
> from cursor_sdk import Agent, AgentOptions, LocalAgentOptions
>
> result = Agent.prompt(
>     "What does the auth middleware do?",
>     AgentOptions(model="composer-2.5", local=LocalAgentOptions(cwd=".")),
> )
> print(result.result)
> ```
>
> Async equivalent (assumes you already have an `AsyncClient` open):
>
> ```python
> from cursor_sdk import AgentOptions, AsyncAgent, LocalAgentOptions
>
> result = await AsyncAgent.prompt(
>     "What does the auth middleware do?",
>     AgentOptions(model="composer-2.5", local=LocalAgentOptions(cwd=".")),
>     client=client,
> )
> ```
>
> ## CursorClient
>
> Use `CursorClient` when you want explicit lifecycle control, a custom bridge endpoint, custom HTTP options, or multiple workspaces in one process. `Client` remains available as an alias.
>
> ```python
> from cursor_sdk import CursorClient, LocalAgentOptions
>
> with CursorClient.launch_bridge(workspace=".") as client:
>     with client.agents.create(
>         model="composer-2.5",
>         api_key="crsr_key",
>         local=LocalAgentOptions(cwd="."),
>     ) as agent:
>         print(agent.send("Summarize what this repository does").text())
> ```
>
> ### Resources
>
> Explicit clients expose resource namespaces:
>
> | Resource       | Sync method examples                                                             | Async method examples                                              |
> | :------------- | :------------------------------------------------------------------------------- | :----------------------------------------------------------------- |
> | `agents`       | `client.agents.create(...)`, `client.agents.list(...)`, `client.agents.get(...)` | `await client.agents.create(...)`, `await client.agents.list(...)` |
> | `models`       | `client.models.list()`                                                           | `await client.models.list()`                                       |
> | `repositories` | `client.repositories.list()`                                                     | `await client.repositories.list()`                                 |
>
> Top-level methods such as `client.create_agent(...)` and `client.list_agents(...)` remain available, but resource namespaces are the preferred shape for application code.
>
> ### Custom HTTP clients
>
> Both sync and async clients accept a custom `httpx` client for proxies, transports, and other advanced HTTP configuration:
>
> ```python
> from cursor_sdk import CursorClient, DefaultHttpxClient
>
> with CursorClient.launch_bridge(
>     workspace=".",
>     http_client=DefaultHttpxClient(proxy="http://proxy.example.com"),
> ) as client:
>     ...
> ```
>
> ```python
> from cursor_sdk import AsyncClient, DefaultAsyncHttpxClient
>
> async with await AsyncClient.launch_bridge(
>     workspace=".",
>     http_client=DefaultAsyncHttpxClient(proxy="http://proxy.example.com"),
> ) as client:
>     ...
> ```
>
> `DefaultHttpxClient` and `DefaultAsyncHttpxClient` keep the SDK's default timeout and redirect behavior. Plain `httpx.Client` and `httpx.AsyncClient` use httpx defaults instead.
>
> ### Configuring timeouts and retries
>
> Both clients expose `with_options(...)`, which returns a shallow copy that shares connection settings and overrides defaults. Use `timeout` for all requests, or set `unary_timeout` and `stream_timeout` separately. `max_retries` controls client retries:
>
> ```python
> short = client.with_options(timeout=5.0, max_retries=2)
> agent = short.agents.create(model="composer-2.5", local=LocalAgentOptions(cwd="."))
> ```
>
> Async equivalent:
>
> ```python
> short_async = async_client.with_options(timeout=5.0, max_retries=2)
> agent = await short_async.agents.create(model="composer-2.5", local=LocalAgentOptions(cwd="."))
> ```
>
> ## Sending messages
>
> Each `agent.send()` returns a `Run`. Each `await async_agent.send()` returns an `AsyncRun`. The agent retains conversation context across runs; the run is the unit of work for one prompt.
>
> ```python
> print(agent.send("Find the bug in src/auth.py").text())
>
> # Same agent, full conversation context is preserved.
> print(agent.send("Fix it and add a regression test").text())
> ```
>
> Async equivalent:
>
> ```python
> run = await agent.send("Find the bug in src/auth.py")
> print(await run.text())
>
> run = await agent.send("Fix it and add a regression test")
> print(await run.text())
> ```
>
> To send images alongside text:
>
> ```python
> run = agent.send(
>     {
>         "text": "What's in this screenshot?",
>         "images": [{"data": base64_png, "mime_type": "image/png"}],
>     }
> )
> ```
>
> You can also use helper dataclasses. `SDKImage.from_file(path)` reads from disk and handles base64 encoding for you:
>
> ```python
> from cursor_sdk import SDKImage, UserMessage
>
> run = agent.send(
>     UserMessage(
>         text="What's in this screenshot?",
>         images=[SDKImage.from_file("screenshot.png")],
>     )
> )
> ```
>
> `SDKImage.data_image(base64_data, mime_type)` and `SDKImage.url_image(url)` are also available for callers that already have encoded bytes or a remote URL.
>
> ### Run
>
> ```python
> class Run:
>     id: str
>     agent_id: str
>     status: str  # "running" | "finished" | "error" | "cancelled" | "expired"
>     result: str
>     model: ModelSelection | None
>     duration_ms: int
>     git: RunGitInfo | None
>     created_at: str | None
>     usage: TokenUsage | None  # cumulative; property on the live handle
>
>     def stream(self) -> Iterator[SDKMessage]: ...
>     def messages(self) -> Iterator[SDKMessage]: ...
>     def events(self) -> Iterator[RunStreamEvent]: ...
>     def iter_text(self) -> Iterator[str]: ...
>     def text(self) -> str: ...
>     def wait(self) -> RunResult: ...
>     def cancel(self) -> None: ...
>     def conversation(self) -> list[ConversationTurn]: ...
>     def conversation_json(self) -> str: ...
>     def observe(self, *, after_offset: str | None = None) -> Iterator[RunStreamEvent]: ...
>
>     def supports(self, operation: str) -> bool: ...
>     def unsupported_reason(self, operation: str) -> str | None: ...
>     def on_did_change_status(
>         self, listener: Callable[[str], None]
>     ) -> Callable[[], None]: ...
> ```
>
> `run.stream()` is an alias for `run.messages()`. Iterating `run` directly yields `RunStreamEvent` envelopes, the same as `run.events()`.
>
> `AsyncRun` exposes the same state fields, including `usage`. Methods that do I/O are async: `async for message in run.stream()`, `async for message in run.messages()`, `async for event in run.events()`, `async for text in run.iter_text()`, `await run.text()`, `await run.wait()`, `await run.cancel()`, `await run.conversation()`, `await run.conversation_json()`, and `async for event in run.observe()`.
>
> ### Streaming
>
> ```python
> run = agent.send("Find the bug in src/auth.py")
>
> for message in run.messages():
>     if message.type == "assistant":
>         for block in message.message.content:
>             if block.type == "text":
>                 print(block.text, end="")
>     elif message.type == "thinking":
>         print(message.text, end="")
>     elif message.type == "tool_call":
>         print(f"[tool] {message.name}: {message.status}")
>     elif message.type == "status":
>         print(f"[status] {message.status}")
>     elif message.type == "usage":
>         print(f"[usage] turn total={message.usage.total_tokens}")
> ```
>
> A run stream is consumable once. `run.messages()`, `run.events()`, and `run.iter_text()` all draw from the same underlying stream and advance it. Once the stream completes, the run holds the terminal result (`run.result`, `run.status`, `run.usage`, `run.git`, ...). Call `run.wait()` to drain any remaining events and return the typed `RunResult`.
>
> ### Waiting without streaming
>
> ```python
> result = run.wait()
>
> print(result.status)       # "finished" | "error" | "cancelled" | "expired"
> print(result.result)       # final assistant text, if any
> print(result.model)        # resolved ModelSelection used for this run
> print(result.duration_ms)
> print(result.usage)        # cumulative TokenUsage, or None if unavailable
> print(result.git)          # RunGitInfo on cloud
> ```
>
> Async equivalent:
>
> ```python
> result = await run.wait()
> ```
>
> ### Token usage
>
> Runs report token usage when the runtime provides it. Read the cumulative total from `run.usage` on the live handle (while streaming or after `wait()`), or from `result.usage` on the `RunResult` returned by `run.wait()`. Both hold a `TokenUsage` summed across every turn that reported usage, and both are `None` when no turn did—for example a cancelled run that never finished a turn, a runtime that doesn't surface usage, or a detached cloud snapshot that hasn't reconciled usage yet.
>
> ```python
> @dataclass(frozen=True)
> class TokenUsage:
>     input_tokens: int
>     output_tokens: int
>     cache_read_tokens: int
>     cache_write_tokens: int
>     total_tokens: int
>     reasoning_tokens: int | None = None
> ```
>
> | Field                | Description                                                                                           |
> | :------------------- | :---------------------------------------------------------------------------------------------------- |
> | `input_tokens`       | Prompt tokens sent to the model.                                                                      |
> | `output_tokens`      | Tokens generated by the model.                                                                        |
> | `cache_read_tokens`  | Tokens served from the prompt cache.                                                                  |
> | `cache_write_tokens` | Tokens written to the prompt cache.                                                                   |
> | `total_tokens`       | `input_tokens + output_tokens + cache_read_tokens + cache_write_tokens`. Excludes `reasoning_tokens`. |
> | `reasoning_tokens`   | Reasoning tokens, a subset of `output_tokens`. `None` when the model or runtime didn't report it.     |
>
> ```python
> result = run.wait()
>
> if result.usage is not None:
>     print(f"total: {result.usage.total_tokens}")
>     print(f"in: {result.usage.input_tokens}, out: {result.usage.output_tokens}")
>     print(
>         f"cache read/write: {result.usage.cache_read_tokens}/{result.usage.cache_write_tokens}"
>     )
> else:
>     print("no usage reported for this run")
> ```
>
> `reasoning_tokens` is already counted inside `output_tokens`, so `total_tokens` leaves it out to avoid double-counting.
>
> For per-turn numbers as they stream, handle the `usage` [stream event](https://cursor.com/docs/sdk/python.md#stream-events) (`SDKUsageMessage`). It fires once at the end of each turn that reported usage and carries that turn's `TokenUsage`. `run.usage` and `result.usage` stay cumulative across the run. After stream turns, the handle prefers those summed totals; otherwise it uses usage from `wait()` or from a `get_run` / `list_runs` snapshot when the bridge supplies it.
>
> ```python
> for message in run.messages():
>     if message.type == "usage":
>         print(f"turn used {message.usage.total_tokens} tokens")
>
> # Or after wait / without consuming messages yourself:
> result = run.wait()
> print(run.usage, result.usage)
> ```
>
> Async equivalent: `async for message in run.messages()` and `await run.wait()`. `run.usage` is still a sync property on `AsyncRun`.
>
> `TokenUsage` is exported from `cursor_sdk` (plus `to_token_usage` / `sum_token_usage` for advanced callers). Wire JSON is camelCase (`inputTokens`, …); the Python dataclasses use snake\_case.
>
> Token counts are what the runtime reports; they say nothing about cost. For billed usage and the dollar cost of an agent's runs, call [`agent.get_usage()`](https://cursor.com/docs/sdk/python.md#agentget_usage).
>
> ### Reading text output
>
> `iter_text()` yields assistant text as it streams. `text()` returns the final terminal text, blocking on `wait()` if the run is still running.
>
> ```python
> for chunk in run.iter_text():
>     print(chunk, end="")
>
> final_text = run.text()
> ```
>
> Async equivalent:
>
> ```python
> async for chunk in run.iter_text():
>     print(chunk, end="")
>
> final_text = await run.text()
> ```
>
> ### Cancelling a run
>
> ```python
> run.cancel()
> ```
>
> Async equivalent:
>
> ```python
> await run.cancel()
> ```
>
> `run.cancel()` requests cancellation of an active run. The status moves to `"cancelled"`, the live stream stops, in-flight tool calls stop, and `run.wait()` resolves with `status: "cancelled"`. Partial output (assistant text written so far) stays on the `Run` object.
>
> Cancelling a run that is already terminal (`"finished"`, `"error"`, `"cancelled"`, `"expired"`) raises `UnsupportedRunOperationError`. Guard with `run.status` when in doubt:
>
> ```python
> if run.status == "running":
>     run.cancel()
> ```
>
> ### Reading run state
>
> ```python
> print(run.id)
> print(run.status)  # "running" | "finished" | "error" | "cancelled" | "expired"
>
> stop = run.on_did_change_status(lambda status: print(f"status changed to {status}"))
> stop()  # remove the listener
>
> turns = run.conversation()
> ```
>
> `run.conversation()` returns a typed `list[ConversationTurn]`. Use it to render or persist structured history without subscribing to the live stream. `run.conversation_json()` returns the raw JSON string.
>
> For async runs, use `await run.conversation()` and `await run.conversation_json()`.
>
> ### Per-run model override
>
> The `model` you pass to `agent.send()` overrides the agent's selection for that run, then becomes sticky: subsequent sends without an override continue to use the new model. To switch back, pass another `model` override or read the current selection from `agent.model`.
>
> ```python
> from cursor_sdk import ModelParameterValue, ModelSelection, SendOptions
>
> run = agent.send(
>     "Plan the refactor",
>     SendOptions(
>         model=ModelSelection(
>             id="composer-2.5",
>             params=[ModelParameterValue(id="fast", value="true")],
>         ),
>     ),
> )
> ```
>
> `run.model` and `result.model` reflect the selection this run used and are immutable once the run starts.
>
> ### Per-run environment variables
>
> Cloud agents can also take environment variables for a single run. Pass `cloud.env_vars` in `SendOptions` and the values are injected into the agent's shell for that run only — when the run finishes, they're removed from the VM and the next run doesn't see them. This is the right shape for credentials that rotate between turns, like a short-lived deploy token you mint right before asking the agent to use it.
>
> ```python
> from cursor_sdk import CloudSendOptions, SendOptions
>
> run = agent.send(
>     "Deploy the preview environment",
>     SendOptions(
>         cloud=CloudSendOptions(env_vars={"DEPLOY_TOKEN": mint_short_lived_token()}),
>     ),
> )
> ```
>
> If a run-scoped variable has the same name as an agent-scoped one from [`env_vars` on `CloudAgentOptions`](https://cursor.com/docs/sdk/python.md#session-environment-variables), the run-scoped value wins for that run, then the agent-scoped value comes back on the next run.
>
> Per-run variables work on the first send too. The SDK passes them along with agent creation, scoped to the initial run, so they aren't persisted on the agent. Like agent-scoped variables, they're encrypted at rest and names can't start with `CURSOR_`.
>
> Per-run environment variables are cloud agents only, and they aren't available for agents running against public repositories. For local agents, the agent process inherits your own environment, so set variables on the process before calling `send()`.
>
> ### Conversation mode
>
> Pass `mode="plan"` or `mode="agent"` to control whether a run explores and plans first or implements changes directly. See [Plan mode](https://cursor.com/help/ai-features/plan-mode.md) for what plan mode does in the product.
>
> Set `mode` in `AgentOptions` passed to `Agent.create()` to seed the first run. On follow-up `agent.send()` calls, omit `mode` to keep the conversation's current mode, or pass `mode` to switch for that run only.
>
> ```python
> from cursor_sdk import Agent, AgentOptions, CloudAgentOptions, CloudRepository, SendOptions
>
> with Agent.create(
>     AgentOptions(
>         model="composer-2.5",
>         mode="plan",
>         cloud=CloudAgentOptions(
>             repos=[CloudRepository(url="https://github.com/your-org/your-repo")],
>         ),
>     )
> ) as agent:
>     agent.send("Design the auth refactor").wait()
>     agent.send(
>         "Looks good, start building",
>         SendOptions(mode="agent"),
>     ).wait()
> ```
>
> ### Streaming raw deltas
>
> Pass `on_delta` and `on_step` callbacks in `SendOptions` for lower-level updates. Sync callbacks are called inline. Async callbacks may be sync or async; awaitable return values are awaited before the next event is processed.
>
> ```python
> from cursor_sdk import SendOptions
>
> def on_delta(update):
>     if update.type in ("text-delta", "thinking-delta"):
>         print(update.text, end="")
>
> run = agent.send(
>     "Refactor the utils module",
>     SendOptions(on_delta=on_delta, on_step=lambda step: print(f"[step] {step.type}")),
> )
> run.wait()
> ```
>
> The concrete update and step subclasses live in `cursor_sdk.events`:
>
> ```python
> from cursor_sdk.events import TextDeltaUpdate, ToolCallStartedUpdate
>
> if isinstance(update, TextDeltaUpdate):
>     print(update.text)
> ```
>
> They remain importable from `cursor_sdk` for backward compatibility, but new code should import from `cursor_sdk.events`.
>
> ### SendOptions
>
> | Property          | Type                                         | Description                                                                                                                                                                                                                              |
> | :---------------- | :------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `model`           | `str \| ModelSelection \| Mapping[str, Any]` | Per-send model override. If omitted, uses `agent.model`. Sticky after a successful send.                                                                                                                                                 |
> | `mode`            | `"agent" \| "plan"`                          | Per-send conversation mode override. If omitted on follow-ups, keeps the conversation's current mode.                                                                                                                                    |
> | `mcp_servers`     | `Mapping[str, McpServerConfig]`              | Inline MCP server definitions. Fully replaces creation-time servers for this run.                                                                                                                                                        |
> | `cloud.env_vars`  | `Mapping[str, str]`                          | Cloud agents only. [Per-run environment variables](https://cursor.com/docs/sdk/python.md#per-run-environment-variables) injected for this run and removed when it finishes. Overrides agent-scoped `env_vars` by name for this run only. |
> | `local.force`     | `bool`                                       | Local agents only. Defaults to `None` (unset). Set `True` to expire a stuck active run before starting this message. Cloud returns `409 agent_busy` server-side, so no equivalent is needed.                                             |
> | `idempotency_key` | `str`                                        | Optional client-generated idempotency key for the send.                                                                                                                                                                                  |
> | `on_step`         | `Callable[[ConversationStep], Any]`          | Callback after each completed conversation step (text, thinking, or tool batch).                                                                                                                                                         |
> | `on_delta`        | `Callable[[InteractionUpdate], Any]`         | Callback per raw `InteractionUpdate`.                                                                                                                                                                                                    |
>
> ***
>
> The next three sections are detailed reference for `SDKMessage`, `InteractionUpdate`, and `ConversationTurn`. Skim or skip on a first read; [Resuming agents](https://cursor.com/docs/sdk/python.md#resuming-agents) picks up the narrative.
>
> ## Stream events
>
> `run.messages()` yields typed SDK message dataclasses. Discriminate on `message.type`. All messages include `agent_id` and `run_id` when the runtime provides them.
>
> ```python
> SDKMessage = (
>     SDKSystemMessage
>     | SDKUserMessageEvent
>     | SDKAssistantMessage
>     | SDKThinkingMessage
>     | SDKToolUseMessage
>     | SDKStatusMessage
>     | SDKTaskMessage
>     | SDKRequestMessage
>     | SDKUsageMessage
>     | Mapping[str, Any]
> )
> ```
>
> | `type`        | Dataclass             | Key fields                                                                  |
> | :------------ | :-------------------- | :-------------------------------------------------------------------------- |
> | `"system"`    | `SDKSystemMessage`    | `subtype`, `model`, `tools`                                                 |
> | `"user"`      | `SDKUserMessageEvent` | `message.content`                                                           |
> | `"assistant"` | `SDKAssistantMessage` | `message.content` with `TextBlock` and `ToolUseBlock` values                |
> | `"thinking"`  | `SDKThinkingMessage`  | `text`, `thinking_duration_ms`                                              |
> | `"tool_call"` | `SDKToolUseMessage`   | `call_id`, `name`, `status`, `args`, `result`, `truncated`                  |
> | `"status"`    | `SDKStatusMessage`    | `status`, `message`                                                         |
> | `"task"`      | `SDKTaskMessage`      | `status`, `text`                                                            |
> | `"request"`   | `SDKRequestMessage`   | `request_id`                                                                |
> | `"usage"`     | `SDKUsageMessage`     | `usage` ([`TokenUsage`](https://cursor.com/docs/sdk/python.md#token-usage)) |
>
> `SDKToolUseMessage` is emitted twice for most tool calls: first with `status="running"` and `args` populated, then again on completion with `status="completed"` (or `"error"`) and `result` populated. `truncated` flags whether the SDK truncated `args` or `result` because the payload was too large.
>
> `SDKUsageMessage` is emitted once at the end of each turn that reported token usage, carrying that turn's [`TokenUsage`](https://cursor.com/docs/sdk/python.md#token-usage). The cumulative total across turns stays on `run.usage` and `result.usage`. See [Token usage](https://cursor.com/docs/sdk/python.md#token-usage).
>
> ```python
> @dataclass(frozen=True)
> class SDKUsageMessage:
>     type: Literal["usage"]
>     agent_id: str
>     run_id: str
>     usage: TokenUsage
> ```
>
> Result data (final text, model, duration, cumulative token usage, git metadata) lives on the `Run` object after the stream completes. Use `run.wait()` to read it, including `result.usage` when the runtime reported it.
>
> > **Tool call schema is not stable.** The `args` and `result` payloads on `tool_call` events reflect each tool's internal shape and can change as tools evolve. Tool names can also be renamed or replaced. Treat `args` and `result` as untyped data and parse defensively. The event envelope (`type`, `call_id`, `name`, `status`) is stable.
>
> `run.events()` yields lower-level `RunStreamEvent` envelopes. Use it when you need offsets, terminal result envelopes, or raw interaction updates:
>
> ```python
> for event in run.events():
>     print(event.kind, event.offset)
> ```
>
> ## Interaction updates
>
> `InteractionUpdate` is the raw delta type passed to the `on_delta` callback on `agent.send()`. Updates are finer-grained than `SDKMessage` events: text streams in token-by-token and tool calls report partial state as args accumulate.
>
> ```python
> InteractionUpdate = (
>     TextDeltaUpdate
>     | ThinkingDeltaUpdate
>     | ThinkingCompletedUpdate
>     | ToolCallStartedUpdate
>     | ToolCallCompletedUpdate
>     | PartialToolCallUpdate
>     | TokenDeltaUpdate
>     | StepStartedUpdate
>     | StepCompletedUpdate
>     | TurnEndedUpdate
>     | UserMessageAppendedUpdate
>     | SummaryUpdate
>     | SummaryStartedUpdate
>     | SummaryCompletedUpdate
>     | ShellOutputDeltaUpdate
>     | UnknownInteractionUpdate
>     | Mapping[str, Any]
> )
> ```
>
> `PartialToolCallUpdate` is emitted as the model streams arguments into a tool call before it commits. The same stability disclaimer that applies to `SDKToolUseMessage.args` applies here.
>
> ## Conversation types
>
> The structured per-turn view of a run, returned by `run.conversation()`. Each item is a wrapper that carries the turn `type` discriminator alongside the typed payload in `turn`.
>
> ```python
> @dataclass(frozen=True)
> class ConversationTurn:
>     type: str  # "agentConversationTurn" | "shellConversationTurn"
>     turn: AgentConversationTurn | ShellConversationTurn | Mapping[str, Any]
>
> @dataclass(frozen=True)
> class AgentConversationTurn:
>     user_message: Mapping[str, Any] | None = None
>     steps: Sequence[ConversationStep] = ()
>
> @dataclass(frozen=True)
> class ShellConversationTurn:
>     shell_command: ShellCommand | None = None
>     shell_output: ShellOutput | None = None
>
> ConversationStep = (
>     AssistantConversationStep
>     | ToolCallConversationStep
>     | ThinkingConversationStep
>     | Mapping[str, Any]
> )
> ```
>
> Discriminate on `turn.type` and read the payload through `turn.turn`:
>
> ```python
> for turn in run.conversation():
>     if turn.type == "agentConversationTurn":
>         for step in turn.turn.steps:
>             print(step.type)
>     elif turn.type == "shellConversationTurn":
>         print(turn.turn.shell_command, turn.turn.shell_output)
> ```
>
> `run.conversation()` from `on_step` callbacks fires per `ConversationStep`, not per turn. Tool-call conversation steps carry a `Mapping[str, Any]` payload. Treat tool-call payload details as untyped data; see the [stability note](https://cursor.com/docs/sdk/python.md#stream-events) under Stream events.
>
> ## Resuming agents
>
> ```python
> Agent.resume(
>     agent_id: str,
>     options: AgentOptions | Mapping[str, Any] | None = None,
>     *,
>     client: CursorClient | None = None,
> ) -> Agent
> ```
>
> Use `Agent.resume()` or `client.agents.resume()` to reattach to an existing agent by ID. Common flows: reconnecting to a long-running cloud agent that was kicked off earlier, or continuing a conversation after the local process restarted. Runtime is auto-detected from the ID prefix (`bc-` is cloud, anything else is local).
>
> ```python
> agent = Agent.resume("bc-abc123")
> run = agent.send("Also update the changelog")
> run.wait()
> ```
>
> Async equivalent:
>
> ```python
> agent = await client.agents.resume("bc-abc123")
> run = await agent.send("Also update the changelog")
> await run.wait()
> ```
>
> `agent.model` is `None` on resume unless you pass `model` again. Inline MCP servers are not persisted across resume; they often carry secrets and live in memory only. Pass them again on resume, or use file-based MCP config (`.cursor/mcp.json` plus `local.setting_sources`) for servers that should survive.
>
> ### Local persistence
>
> Local agents persist conversation state and run metadata through the bridge, so follow-ups and `Agent.resume()` survive a process restart. The bridge keeps this under a per-workspace state root on disk by default. Cloud agents persist server-side, so resuming a cloud agent from anywhere returns the same conversation.
>
> Local persistence is workspace-scoped. When the bridge runs as a long-lived sidecar or subprocess, give it the same workspace as the agent so local list, get, and resume calls resolve the right agents. Set it once on the client and pass `cwd` to the local list and get calls:
>
> ```python
> from cursor_sdk import CursorClient
>
> with CursorClient.launch_bridge(workspace="/path/to/repo") as client:
>     agents = client.agents.list(runtime="local", cwd="/path/to/repo")
>     info = client.agents.get(agents.items[0].agent_id, cwd="/path/to/repo")
> ```
>
> ## Inspecting agents and runs
>
> Use `CursorClient` for list, get, and pagination APIs.
>
> ```python
> from cursor_sdk import CursorClient
>
> with CursorClient.launch_bridge(workspace=".") as client:
>     agents = client.agents.list(runtime="local", cwd=".")
>
>     for agent_info in agents.auto_paging_iter():
>         print(agent_info.agent_id)
>
>     info = client.agents.get(agents.items[0].agent_id)
>     runs = client.agents.list_runs(info.agent_id)
>     run = client.agents.get_run(runs.items[0].id)
> ```
>
> Async equivalent:
>
> ```python
> agents = await client.agents.list(runtime="local", cwd=".")
>
> async for agent_info in agents.auto_paging_iter():
>     print(agent_info.agent_id)
>
> info = await client.agents.get(agents.items[0].agent_id)
> runs = await client.agents.list_runs(info.agent_id)
> run = await client.agents.get_run(runs.items[0].id)
> ```
>
> Use `agent.list_messages()` on an agent handle to read message history. `Agent.messages.list(agent_id)` is a typed-attribute convenience for the same call when you only have an ID.
>
> Use `Agent.get_run(run_id)` or `client.agents.get_run(run_id)` to fetch a run
> without an agent handle. Cancel it with
> `Agent.cancel_run(run_id, agent_id=...)` or
> `client.agents.cancel_run(run_id, agent_id=...)`. The async client methods are
> awaitable and use the same arguments.
>
> `AgentMessage` is distinct from a streamed `SDKMessage`:
>
> ```python
> @dataclass(frozen=True)
> class AgentMessage:
>     type: str
>     uuid: str
>     agent_id: str
>     message: Any = None
> ```
>
> List endpoints return `ListResult[T]`. Use `.items` and `.next_cursor` directly, iterate the current page with `for item in page`, or iterate all pages with `.auto_paging_iter()`. Async list endpoints return `AsyncListResult[T]`; `async for item in page` walks the current page, and `async for item in page.auto_paging_iter()` walks every page in the result set.
>
> ### SDKAgentInfo
>
> The metadata shape returned by `Agent.list()`, `Agent.get()`, `client.agents.list()`, and `client.agents.get()`.
>
> ```python
> @dataclass(frozen=True)
> class SDKAgentInfo:
>     agent_id: str
>     name: str
>     summary: str
>     last_modified: str | None = None
>     status: str | None = None  # "running" | "finished" | "error"
>     created_at: str | None = None
>     archived: bool = False
>     runtime: Literal["local", "cloud"] | None = None
>     cwd: str = ""
>     env: CloudEnvironment | None = None
>     repos: Sequence[str] = ()
>     metadata: Mapping[str, str] = {}  # from CloudAgentOptions.metadata; empty for local agents
> ```
>
> ### Cloud agent lifecycle
>
> Cloud agents stay in your team's workspace until you archive or delete them. `client.agents.list(runtime="cloud")` hides archived agents by default; pass `include_archived=True` to see them. Filter by `pr_url` to find the agent that opened a specific pull request.
>
> ```python
> # By ID, no agent handle required:
> Agent.archive(agent_id)
> Agent.unarchive(agent_id)
> Agent.delete(agent_id)
>
> # Through an explicit client:
> client.agents.archive(agent_id)
> client.agents.unarchive(agent_id)
> client.agents.delete(agent_id)
>
> # On an existing agent handle:
> agent.archive()
> agent.unarchive()
> agent.delete()
> ```
>
> `archive` soft-deletes the agent so the transcript stays readable. `unarchive` restores it. `delete` is permanent; subsequent reads return `NotFoundError`.
>
> Async lifecycle methods use the same names and are awaitable.
>
> ### agent.get\_usage()
>
> Fetch billed token usage and dollar cost for an agent's runs. Cloud agents return a per-run breakdown; local agents return a per-turn breakdown. Pass `run_id` to restrict the result to one entry: for cloud agents a `run-<uuid>` run ID, for local agents an ID from a previous `get_usage().runs[].run_id`.
>
> ```python
> usage = agent.get_usage()
>
> print(f"tokens: {usage.usage.total_tokens}")
> if usage.cost is not None:
>     print(f"charged: ${usage.cost.charged_cents / 100:.2f}")
> for run in usage.runs:
>     print(run.run_id, run.usage.total_tokens)
> ```
>
> ```python
> @dataclass(frozen=True)
> class AgentUsage:
>     usage: TokenUsage              # summed across `runs`
>     runs: Sequence[RunUsage] = ()
>     cost: UsageCost | None = None  # summed across `runs`
>
> @dataclass(frozen=True)
> class RunUsage:
>     run_id: str
>     usage: TokenUsage
>     cost: UsageCost | None = None
>
> @dataclass(frozen=True)
> class UsageCost:
>     raw_cost_cents: float  # undiscounted model token cost; 0 for request-priced usage
>     charged_cents: float   # amount charged, discounts and the Cursor Token Rate included
> ```
>
> Cost includes discounts and can take a moment to settle after a run ends; `cost` is `None` until it does. `charged_cents` is `0.0` for plan-included, BYOK, and credit-grant usage.
>
> This is a different view from [Token usage](https://cursor.com/docs/sdk/python.md#token-usage): `run.usage` is the live token count for one run, while `get_usage()` is the billed record across the agent's runs. On async agents, `await agent.get_usage()` matches. `AgentUsage`, `RunUsage`, and `UsageCost` are exported from `cursor_sdk`.
>
> ## The Cursor namespace
>
> Account-level and catalog reads. Sync methods take optional `api_key` and otherwise fall back to `CURSOR_API_KEY`.
>
> ```python
> from cursor_sdk import Cursor
>
> me = Cursor.me()
> models = Cursor.models.list()
> repositories = Cursor.repositories.list()
> ```
>
> Explicit-client equivalent:
>
> ```python
> me = client.me()
> models = client.models.list()
> repositories = client.repositories.list()
> ```
>
> Async equivalent:
>
> ```python
> from cursor_sdk import AsyncCursor
>
> me = await AsyncCursor.me(client=client)
> models = await AsyncCursor.models.list(client=client)
> repositories = await AsyncCursor.repositories.list(client=client)
> ```
>
> `Cursor.me()` returns an `SDKUser` with `api_key_name`, `created_at`, and
> optional `user_id`, `user_email`, `user_first_name`, and `user_last_name`
> fields.
>
> Use `Cursor.models.list()` to discover valid model IDs and per-model parameters before calling `Agent.create()` or `agent.send()`. Parameters are model-specific. Common examples are reasoning effort and Cursor Router's `optimize_for` on `auto-smart`.
>
> The catalog is account- and team-specific. Cursor Router only appears as `auto-smart` when Router is available for the API key's team. See [Cursor Router](https://cursor.com/docs/sdk/python.md#cursor-router).
>
> ```python
> models = Cursor.models.list()
> composer = next((model for model in models if model.id == "composer-2.5"), None)
>
> print(composer.parameters if composer else [])
> # [
> #   ModelParameterDefinition(
> #       id="fast",
> #       display_name="Fast",
> #       values=(
> #           ModelParameterDefinitionValue(value="false"),
> #           ModelParameterDefinitionValue(value="true", display_name="Fast"),
> #       ),
> #   ),
> # ]
> ```
>
> Preset `variants` on each `SDKModel` already contain valid `params`, so you can copy them into a `ModelSelection`.
>
> Prefer an explicit Router selection (`auto-smart` + `optimize_for`) when a target model is missing and you want Cost, Balance, or Intelligence. Fall back to `ModelSelection(id="auto")` only when you want server-selected Auto without choosing a Router mode. For Cursor Router, always pass `optimize_for` explicitly.
>
> `Cursor.repositories.list()` returns the SCM repositories (GitHub, GitLab, Bitbucket, Azure DevOps, depending on what's connected) available for cloud agents on the calling account or team. Each item exposes a `url`. Use these to populate `CloudAgentOptions.repos`.
>
> ## MCP servers
>
> Agents can pick up MCP servers from inline definitions, project/user settings, plugins, and dashboard-managed configuration depending on the runtime.
>
> ```python
> from cursor_sdk import (
>     Agent,
>     AgentOptions,
>     HttpMcpServerConfig,
>     LocalAgentOptions,
>     McpAuth,
>     StdioMcpServerConfig,
> )
>
> agent = Agent.create(
>     AgentOptions(
>         model="composer-2.5",
>         local=LocalAgentOptions(cwd="."),
>         mcp_servers={
>             "docs": HttpMcpServerConfig(
>                 url="https://example.com/mcp",
>                 auth=McpAuth(client_id="client-id", scopes=["read", "write"]),
>             ),
>             "filesystem": StdioMcpServerConfig(
>                 command="npx",
>                 args=["-y", "@modelcontextprotocol/server-filesystem", "."],
>             ),
>         },
>     )
> )
> ```
>
> Flat dictionaries (`{"type": "http", "url": ...}` and `{"type": "stdio", "command": ...}`) are also accepted as a quick-script convenience.
>
> ### What gets loaded
>
> **Local agents** load servers from up to five sources, with first-match-wins precedence on conflicting names:
>
> 1. `mcp_servers` on `agent.send()`. Fully replaces creation-time servers for that run (not merged).
> 2. `mcp_servers` on `Agent.create()`. Used when no per-send override is provided.
> 3. Plugin servers, if `local.setting_sources` includes `"plugins"`.
> 4. Project servers from `.cursor/mcp.json`, if `local.setting_sources` includes `"project"`.
> 5. User servers from `~/.cursor/mcp.json`, if `local.setting_sources` includes `"user"`.
>
> Without `local.setting_sources`, only inline servers are loaded. If a local MCP server requires OAuth login, the SDK can reuse a saved login from the Cursor app, but it cannot open a browser to sign you in.
>
> **Cloud agents** load servers from:
>
> 1. `mcp_servers` on `agent.send()`. Fully replaces creation-time servers for that run (not merged).
> 2. `mcp_servers` on `Agent.create()`. Used when no per-send override is provided.
> 3. Your user and team MCP servers from [cursor.com/agents](https://cursor.com/agents).
>
> If an inline server doesn't include `auth` or `headers` and you've previously authorized that server URL on cursor.com/agents, runs authenticated with a personal API token reuse those OAuth tokens automatically. Service account API keys cannot fall back to user auth as they are not associated with a user.
>
> `local.setting_sources` does not apply to cloud agents.
>
> ### Cloud
>
> Cloud agents accept authenticated MCP configs inline too. Cloud MCP supports HTTP and stdio transports. Use HTTP `headers` for static API keys or Bearer tokens. Use HTTP `auth` for OAuth-protected servers. Use stdio `env` when the server runs inside the cloud VM and reads credentials from environment variables.
>
> ```python
> from cursor_sdk import (
>     Agent,
>     AgentOptions,
>     CloudAgentOptions,
>     CloudRepository,
>     HttpMcpServerConfig,
>     StdioMcpServerConfig,
> )
>
> agent = Agent.create(
>     AgentOptions(
>         model="composer-2.5",
>         cloud=CloudAgentOptions(
>             repos=[CloudRepository(url="https://github.com/your-org/your-repo")],
>         ),
>         mcp_servers={
>             "linear": HttpMcpServerConfig(
>                 url="https://mcp.linear.app/mcp",
>                 headers={"Authorization": "Bearer linear_pat_xxx"},
>             ),
>             "github": StdioMcpServerConfig(
>                 command="npx",
>                 args=["-y", "@modelcontextprotocol/server-github"],
>                 env={"GITHUB_TOKEN": "ghp_xxx"},
>             ),
>         },
>     )
> )
> ```
>
> - HTTP `headers` and `auth` are handled by Cursor's backend. Sensitive fields are redacted and do not enter the VM.
> - Stdio `env` values are passed into the VM because the server runs there. Treat them like any other runtime secret.
> - OAuth for MCP servers configured on cursor.com/agents stays per-user, even for team-level servers.
>
> See [MCP](https://cursor.com/docs/mcp.md) for the full config format and [Cloud Agent capabilities](https://cursor.com/docs/cloud-agent/capabilities.md#mcp-tools) for cloud-specific behavior.
>
> ## Subagents
>
> Define named subagents that the main agent can spawn via the `Agent` tool. Pass them inline:
>
> ```python
> from cursor_sdk import Agent, AgentDefinition, AgentOptions, LocalAgentOptions
>
> agent = Agent.create(
>     AgentOptions(
>         model="composer-2.5",
>         local=LocalAgentOptions(cwd="."),
>         agents={
>             "code-reviewer": AgentDefinition(
>                 description="Expert code reviewer for quality and security.",
>                 prompt="Review code for bugs, security issues, and proven approaches.",
>                 model="inherit",
>             ),
>             "test-writer": AgentDefinition(
>                 description="Writes tests for code changes.",
>                 prompt="Write comprehensive tests for the given code.",
>             ),
>         },
>     )
> )
> ```
>
> Subagents committed to the repo at `.cursor/agents/*.md` (with `name`, `description`, and optional `model` frontmatter) are also picked up. Inline definitions override file-based ones with the same name.
>
> ### Nested subagents
>
> Subagents can spawn their own subagents, within a nesting limit. When a subagent uses the `Agent` tool, it reaches the same subagent executor the parent has, so a parent can delegate to a subagent that delegates further. Each level sees the same set of named subagents. The top-level agent and its direct subagents can launch subagents, but a subagent launched by another subagent can't launch further ones.
>
> ## Restricting the toolset
>
> `tools` allowlists the built-in tools offered to the model; `disallowed_tools` removes tools and keeps the rest, including tools added to the platform after your SDK version was released. Both are local agents only for now, and neither persists on the agent: pass them again on resume to keep the restriction.
>
> ```python
> from cursor_sdk import Agent, AgentOptions, LocalAgentOptions
>
> # Read-only agent: only these tools are offered.
> reader = Agent.create(
>     AgentOptions(
>         model="composer-2.5",
>         tools=["read", "grep", "glob", "ls"],
>         local=LocalAgentOptions(cwd="."),
>     )
> )
>
> # Everything except shell access.
> no_shell = Agent.create(
>     AgentOptions(
>         model="composer-2.5",
>         disallowed_tools=["shell"],
>         local=LocalAgentOptions(cwd="."),
>     )
> )
> ```
>
> - Omitting `tools` offers the standard toolset for the selected model; `tools=[]` offers no built-in tools, so the model can only respond with text.
> - Both fields accept public names (`"read"`, `"edit"`, `"task"`, `"webSearch"`, ...) and the capability groups `"shell"` and `"mcp"`. Unknown names raise `BadRequestError` at creation.
> - Deny wins: a tool must be in `tools` (when set) and not in `disallowed_tools` to be offered.
> - Disallowing `"mcp"` also removes [custom tools](https://cursor.com/docs/sdk/python.md#custom-tools). Disallowing `"task"` prevents [subagents](https://cursor.com/docs/sdk/python.md#subagents); otherwise subagents keep their own curated toolsets.
>
> ## Custom tools
>
> Custom tools let you expose Python functions to local agents without standing up a separate MCP server. Pass them on `LocalAgentOptions.custom_tools`.
>
> ```python
> from cursor_sdk import Agent, CustomTool, CustomToolContext, LocalAgentOptions
>
> def get_deployment_status(args, context: CustomToolContext):
>     service = args["service"]
>     return f"Service {service} is healthy."
>
> with Agent.create(
>     model="composer-2.5",
>     local=LocalAgentOptions(
>         cwd=".",
>         custom_tools={
>             "get_deployment_status": CustomTool(
>                 description="Look up the current deployment status for a service.",
>                 input_schema={
>                     "type": "object",
>                     "properties": {
>                         "service": {"type": "string", "description": "Service name"},
>                     },
>                     "required": ["service"],
>                 },
>                 execute=get_deployment_status,
>             ),
>         },
>     ),
> ) as agent:
>     agent.send("Is the checkout service healthy?").wait()
> ```
>
> `execute` receives the parsed arguments and a `CustomToolContext` with `tool_call_id` when available. It can return a string, a JSON-compatible value, or a mapping with a `content` list. Custom tools are local agents only.
>
> ## Hooks
>
> Hooks are file-based only. There is no programmatic hook callback. Hooks are a project policy boundary, not a per-run knob.
>
> - **Local:** Add `.cursor/hooks.json` to the repo passed as `local.cwd`, or add `~/.cursor/hooks.json` for user-level hooks.
> - **Cloud:** Commit `.cursor/hooks.json` and its scripts to the repo passed in `cloud.repos`. SDK-created cloud agents load project hooks automatically. On Enterprise plans, they also run team hooks and enterprise-managed hooks.
>
> See [Hooks](https://cursor.com/docs/hooks.md) for the configuration format and [Cloud Agents hooks support](https://cursor.com/docs/cloud-agent.md#hooks-support) for cloud behavior.
>
> ## Artifacts
>
> List and download files from the agent's workspace.
>
> ```python
> @dataclass(frozen=True)
> class SDKArtifact:
>     path: str
>     size_bytes: int = 0
>     updated_at: str = ""
> ```
>
> ```python
> from pathlib import Path
>
> artifacts = agent.list_artifacts()
>
> for artifact in artifacts:
>     print(artifact.path, artifact.size_bytes)
>
> # Download a single artifact to disk.
> content = agent.download_artifact(artifacts[0].path)
> Path("review.md").write_bytes(content)
> ```
>
> Async agents expose `await agent.list_artifacts()` and `await agent.download_artifact(path)`.
>
> Artifact support is runtime-dependent. Local SDK agents return an empty list from `list_artifacts()` and raise from `download_artifact()`.
>
> ## Resource management
>
> Always close agents when done. The cleanest sync pattern is a context manager:
>
> ```python
> from cursor_sdk import Agent, LocalAgentOptions
>
> with Agent.create(model="composer-2.5", local=LocalAgentOptions(cwd=".")) as agent:
>     agent.send("Summarize the repository").wait()
> ```
>
> To dispose explicitly:
>
> ```python
> agent.close()
> ```
>
> Async agents and clients support async context managers and `await` cleanup:
>
> ```python
> from cursor_sdk import AsyncClient, LocalAgentOptions
>
> async with await AsyncClient.launch_bridge(workspace=".") as client:
>     async with await client.agents.create(
>         model="composer-2.5",
>         local=LocalAgentOptions(cwd="."),
>     ) as agent:
>         run = await agent.send("Summarize the repository")
>         await run.wait()
> ```
>
> To dispose explicitly:
>
> ```python
> await agent.close()
> await client.aclose()
> ```
>
> The module-level sync default client is closed automatically at process exit. Long-running processes can close and reset it explicitly:
>
> ```python
> from cursor_sdk import close_default_client
>
> close_default_client()
> ```
>
> ## Configuration reference
>
> The Python SDK accepts helper dataclasses and raw dictionaries. Dataclasses use Python `snake_case` fields and are preferred for application code.
>
> ### AgentOptions
>
> | Property           | Type                                                 | Default                                                             | Description                                                                                                                                                                           |
> | :----------------- | :--------------------------------------------------- | :------------------------------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `model`            | `str \| ModelSelection \| Mapping[str, Any]`         | Required for local; cloud falls back to the server-resolved default | Model to use. See [`ModelSelection`](https://cursor.com/docs/sdk/python.md#modelselection).                                                                                           |
> | `api_key`          | `str`                                                | `CURSOR_API_KEY` env                                                | User API key or service account key. Team Admin keys are not yet supported.                                                                                                           |
> | `name`             | `str`                                                | Auto-generated                                                      | Human-readable agent name surfaced in `client.agents.list()` / `client.agents.get()`.                                                                                                 |
> | `local`            | `LocalAgentOptions \| Mapping[str, Any]`             | `None`                                                              | Local agent config. Pass to create a local agent.                                                                                                                                     |
> | `cloud`            | `CloudAgentOptions \| Mapping[str, Any]`             | `None`                                                              | Cloud agent config. Pass to create a cloud agent.                                                                                                                                     |
> | `mcp_servers`      | `Mapping[str, McpServerConfig]`                      | `None`                                                              | Inline MCP server definitions.                                                                                                                                                        |
> | `agents`           | `Mapping[str, AgentDefinition \| Mapping[str, Any]]` | `None`                                                              | Subagent definitions.                                                                                                                                                                 |
> | `tools`            | `Sequence[str]`                                      | Default toolset                                                     | Only the listed built-in tools are offered to the model. `[]` means no built-in tools; the model can only respond with text. Local agents only.                                       |
> | `disallowed_tools` | `Sequence[str]`                                      | `None`                                                              | Removes the listed built-in tools; everything else stays available. Deny wins when combined with `tools`. Local agents only.                                                          |
> | `agent_id`         | `str`                                                | Auto-generated                                                      | Durable agent ID. Pass to keep a stable ID across invocations.                                                                                                                        |
> | `idempotency_key`  | `str`                                                | Auto-generated for cloud                                            | Optional client-generated idempotency key. Cloud only.                                                                                                                                |
> | `mode`             | `"agent" \| "plan"`                                  | `None`                                                              | Initial conversation mode for the agent's first run. When omitted, the server starts in agent mode. See [Conversation mode](https://cursor.com/docs/sdk/python.md#conversation-mode). |
>
> ### LocalAgentOptions
>
> | Property          | Type                                            | Default | Description                                                                                                                         |
> | :---------------- | :---------------------------------------------- | :------ | :---------------------------------------------------------------------------------------------------------------------------------- |
> | `cwd`             | `str \| os.PathLike`                            | `None`  | Primary working directory. Multi-entry lists are rejected; use `dirs` for multi-root.                                               |
> | `dirs`            | `Sequence[str \| os.PathLike]`                  | `None`  | Additional workspace folders for multi-root setups. Merged with `cwd` so rules, skills, and workspace context load from every path. |
> | `setting_sources` | `Sequence[SettingSource]`                       | `None`  | Ambient settings layers: `"project"`, `"user"`, `"team"`, `"mdm"`, `"plugins"`, or `"all"`.                                         |
> | `sandbox_options` | `SandboxOptions \| Mapping[str, Any]`           | `None`  | Local sandbox options.                                                                                                              |
> | `store`           | `LocalAgentStoreConfig \| Mapping[str, Any]`    | `None`  | Local store config passed to the bridge.                                                                                            |
> | `auto_review`     | `bool`                                          | `None`  | Route local tool calls through Auto-review when the connected backend supports it.                                                  |
> | `custom_tools`    | `Mapping[str, CustomTool \| Mapping[str, Any]]` | `None`  | [Custom tools](https://cursor.com/docs/sdk/python.md#custom-tools) exposed to local agents.                                         |
>
> ### CloudAgentOptions
>
> | Property                    | Type                                             | Default                                                | Description                                                                                                                                                                                                                    |
> | :-------------------------- | :----------------------------------------------- | :----------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `env`                       | `CloudEnvironment \| Mapping[str, Any]`          | `None`                                                 | Execution environment. When omitted, the server uses Cursor-hosted cloud VMs. `pool` and `machine` target self-hosted workers you run.                                                                                         |
> | `repos`                     | `Sequence[CloudRepository \| Mapping[str, Any]]` | `None`                                                 | Repositories to clone into the VM. Omit or pass `[]` for a [no-repo agent](https://cursor.com/docs/sdk/python.md#no-repo-cloud-agents) with an empty workspace. Pass `pr_url` on a repo to attach the agent to an existing PR. |
> | `work_on_current_branch`    | `bool`                                           | `None`                                                 | Push commits to the existing branch instead of a new one. The server treats an omitted value as `False`.                                                                                                                       |
> | `auto_create_pr`            | `bool`                                           | `None`                                                 | Open a PR when the run finishes. The server treats an omitted value as `False`.                                                                                                                                                |
> | `open_as_cursor_github_app` | `bool`                                           | `True` for service-account keys, `False` for user keys | Open PRs as the Cursor GitHub App instead of the API key's owner. The resolved value is echoed on create, get, and list.                                                                                                       |
> | `skip_reviewer_request`     | `bool`                                           | `None`                                                 | Skip requesting the calling user as a reviewer on the PR. The server treats an omitted value as `False`.                                                                                                                       |
> | `env_vars`                  | `Mapping[str, str]`                              | `None`                                                 | Session-scoped environment variables for cloud agents.                                                                                                                                                                         |
> | `metadata`                  | `Mapping[str, str]`                              | `None`                                                 | Caller-owned string tags persisted on the cloud agent. See [Agent metadata](https://cursor.com/docs/sdk/python.md#agent-metadata).                                                                                             |
>
> ### AgentDefinition
>
> | Property      | Type                                                             | Default    | Description                                                                                      |
> | :------------ | :--------------------------------------------------------------- | :--------- | :----------------------------------------------------------------------------------------------- |
> | `description` | `str`                                                            | *required* | When to use this subagent. Shown to the parent agent so it knows when to spawn.                  |
> | `prompt`      | `str`                                                            | *required* | System prompt for the subagent.                                                                  |
> | `model`       | `str \| ModelSelection \| Mapping[str, Any] \| "inherit"`        | `None`     | Model override. `None` and `"inherit"` both use the parent's selection.                          |
> | `mcp_servers` | `Sequence[str \| AgentDefinitionMcpServer \| Mapping[str, Any]]` | `None`     | MCP servers available to this subagent. Names reference servers from the parent's `mcp_servers`. |
>
> ### CustomTool
>
> ```python
> @dataclass
> class CustomTool:
>     execute: Callable[[Mapping[str, Any], CustomToolContext], Any]
>     description: str | None = None
>     input_schema: Mapping[str, Any] | None = None
>
> class CustomToolContext:
>     tool_call_id: str | None = None
> ```
>
> ### ModelSelection
>
> ```python
> @dataclass(frozen=True)
> class ModelSelection:
>     id: str
>     params: Sequence[ModelParameterValue] = ()
>
> @dataclass(frozen=True)
> class ModelParameterValue:
>     id: str
>     value: str
> ```
>
> `id` is the model identifier (for example, `"composer-2.5"` or `"auto-smart"`). `params` carries per-model parameters such as reasoning effort or Router's `optimize_for`. Use `Cursor.models.list()` to discover valid IDs, parameter definitions, and preset variants for your account. See [Cursor Router](https://cursor.com/docs/sdk/python.md#cursor-router) for the Router selection contract.
>
> ### McpServerConfig
>
> ```python
> from cursor_sdk.types import McpServerConfig
>
> @dataclass(frozen=True)
> class HttpMcpServerConfig:
>     url: str
>     type: Literal["http", "sse"] | str = "http"
>     headers: Mapping[str, str] | None = None
>     auth: McpAuth | Mapping[str, Any] | None = None
>
> @dataclass(frozen=True)
> class SseMcpServerConfig(HttpMcpServerConfig):
>     type: Literal["sse"] = "sse"
>
> @dataclass(frozen=True)
> class StdioMcpServerConfig:
>     command: str
>     args: Sequence[str] | None = None
>     env: Mapping[str, str] | None = None
>     cwd: str | os.PathLike | None = None  # local only; cloud rejects this field
>
> @dataclass(frozen=True)
> class McpAuth:
>     client_id: str
>     client_secret: str | None = None
>     scopes: Sequence[str] = ()
> ```
>
> For HTTP servers running in the cloud, `headers` and `auth` are handled by Cursor's backend. Sensitive fields are redacted before the VM sees them. For stdio servers in the cloud, `env` values are passed into the VM (treat them like any runtime secret).
>
> ### UserMessage
>
> ```python
> @dataclass(frozen=True)
> class UserMessage:
>     text: str
>     images: Sequence[SDKImage | Mapping[str, Any]] | None = None
> ```
>
> The structured form of `agent.send()`'s message argument. Use it to send images alongside text.
>
> ### SDKImage
>
> ```python
> @dataclass(frozen=True)
> class SDKImage:
>     url: str | None = None
>     data: str | None = None
>     mime_type: str | None = None
>     dimension: SDKImageDimension | Mapping[str, Any] | None = None
>
>     @classmethod
>     def from_url(cls, url: str, dimension=None) -> SDKImage: ...
>
>     @classmethod
>     def from_data(cls, data: bytes | str, mime_type: str, dimension=None) -> SDKImage: ...
>
>     @classmethod
>     def url_image(cls, url: str, dimension=None) -> SDKImage: ...
>
>     @classmethod
>     def data_image(cls, data: str, mime_type: str, dimension=None) -> SDKImage: ...
>
>     @classmethod
>     def from_file(cls, path, *, mime_type=None, dimension=None) -> SDKImage: ...
> ```
>
> Pass either a remote `url` or base64 `data` with a `mime_type`. `from_data()` accepts bytes or a base64 string. `from_file()` reads a file from disk and base64-encodes it.
>
> ### SettingSource
>
> `SettingSource` is available from `cursor_sdk.types`.
>
> ```python
> from cursor_sdk.types import SettingSource
> ```
>
> Controls which on-disk settings layers a local agent loads. Cloud agents always load `project`, `team`, and `plugins` and ignore this field.
>
> | Value       | Source                                  |
> | :---------- | :-------------------------------------- |
> | `"project"` | `.cursor/` in the workspace             |
> | `"user"`    | `~/.cursor/`                            |
> | `"team"`    | Team settings synced from the dashboard |
> | `"mdm"`     | MDM-managed enterprise settings         |
> | `"plugins"` | Plugin-provided settings                |
> | `"all"`     | Shorthand for all of the above          |
>
> ### ListResult
>
> ```python
> @dataclass(frozen=True)
> class ListResult(Generic[T]):
>     items: list[T]
>     next_cursor: str = ""
>
>     def to_dict(self) -> dict[str, Any]: ...
>     def has_next_page(self) -> bool: ...
>     def next_page_info(self) -> dict[str, str]: ...
>     def get_next_page(self) -> ListResult[T]: ...
>     def auto_paging_iter(self) -> Iterator[T]: ...
> ```
>
> Returned by `client.agents.list()`, `client.agents.list_runs()`, and `Agent.list()`. `next_cursor` is empty when there are no more pages. Async list endpoints return `AsyncListResult[T]` with awaitable equivalents.
>
> ## Errors
>
> All SDK errors extend `CursorAgentError`. `CursorSDKError` is the backward-compatible alias root for older callers. Use `is_retryable` and `retry_after` to drive retry logic.
>
> ```python
> class CursorAgentError(Exception):
>     message: str
>     code: str | None
>     status: int | None
>     status_code: int | None
>     details: list[Mapping[str, Any]]
>     is_retryable: bool
>     cause: BaseException | None
>     proto_error_code: str | None
>     request_id: str | None
>     headers: Mapping[str, str]
>     retry_after: str | None
> ```
>
> | Error                          | When                                                                                                                    |
> | :----------------------------- | :---------------------------------------------------------------------------------------------------------------------- |
> | `AuthenticationError`          | Invalid API key or not logged in.                                                                                       |
> | `PermissionDeniedError`        | Authenticated caller does not have permission for the requested operation.                                              |
> | `RateLimitError`               | Too many requests or usage limits exceeded.                                                                             |
> | `ConfigurationError`           | Invalid model, missing required configuration, or bad request parameters.                                               |
> | `AgentBusyError`               | Sending a follow-up while the agent already has a run in `CREATING` or `RUNNING` state (HTTP `409`, code `agent_busy`). |
> | `BadRequestError`              | Request is malformed.                                                                                                   |
> | `IntegrationNotConnectedError` | Creating a cloud agent for a repo whose SCM provider is not connected.                                                  |
> | `NetworkError`                 | Service unavailable or network failure.                                                                                 |
> | `APITimeoutError`              | Request timed out.                                                                                                      |
> | `InternalServerError`          | Cursor service returned a server error.                                                                                 |
> | `NotFoundError`                | Requested resource was not found.                                                                                       |
> | `AgentNotFoundError`           | Agent does not exist or isn't visible under the current working directory.                                              |
> | `UnsupportedRunOperationError` | Run operation is not supported for the current run state.                                                               |
>
> ### Retrying with backoff
>
> `is_retryable` and `retry_after` drive caller-side retry logic. `retry_after` is an HTTP-style string (seconds, or an HTTP date) supplied by the server when it's set.
>
> ```python
> import time
>
> from cursor_sdk import Agent, AgentOptions, CursorAgentError, LocalAgentOptions, RateLimitError
>
> for attempt in range(3):
>     try:
>         result = Agent.prompt(
>             "Audit the auth middleware for missing input validation",
>             AgentOptions(model="composer-2.5", local=LocalAgentOptions(cwd=".")),
>         )
>         break
>     except RateLimitError as err:
>         time.sleep(float(err.retry_after) if err.retry_after else 2**attempt)
>     except CursorAgentError as err:
>         if not err.is_retryable:
>             raise
>         time.sleep(2**attempt)
> ```
>
> Every `CursorAgentError` includes `request_id` when the server returned one. Log it whenever you surface an error so support has a handle on the failure.
>
> ### IntegrationNotConnectedError
>
> ```python
> class IntegrationNotConnectedError(ConfigurationError):
>     provider: str   # e.g. "github", "gitlab", "azuredevops"
>     help_url: str   # dashboard link to reconnect
> ```
>
> Use `help_url` to point the user at the right reconnect flow. New providers may be added without an SDK release.
>
> ### AgentBusyError
>
> Cloud agents allow only one active run at a time. `AgentBusyError` is raised when you call `agent.send()` (or otherwise create a run) while another run on the same agent is still `CREATING` or `RUNNING`.
>
> `is_retryable` is `False`. Retrying immediately will keep failing until the active run reaches a terminal status or you cancel it. Other `409` responses, such as `agent_archived`, raise `ConfigurationError` instead.
>
> Wait for the active run to finish, cancel it with `run.cancel()`, or poll `Agent.list_runs()` before sending again:
>
> ```python
> from cursor_sdk import Agent, AgentBusyError
>
> agent = Agent.resume("bc-00000000-0000-0000-0000-000000000001")
>
> try:
>     agent.send("Also add tests for the auth middleware.")
> except AgentBusyError:
>     runs = Agent.list_runs(agent.agent_id, {"runtime": "cloud", "limit": 1})
>     active = runs.items[0] if runs.items else None
>     if active is not None and active.status == "running":
>         active.cancel()
>     agent.send("Also add tests for the auth middleware.")
> ```
>
> Local agents do not raise `AgentBusyError`. Pass `local={"force": True}` on `send()` to expire a stuck local run before starting a new one.
>
> ### UnsupportedRunOperationError
>
> ```python
> class UnsupportedRunOperationError(ConfigurationError):
>     operation: str
> ```
>
> Raised when a `Run` operation is not allowed on the current run. The most common case is `run.cancel()` on a run that's already terminal.
>
> `run.supports(operation)` and `run.unsupported_reason(operation)` report SDK-level capability for an operation name (`"stream"`, `"wait"`, `"cancel"`, `"conversation"`) and do not check run state. Read `run.status` to guard state-sensitive calls.
>
> ## Troubleshooting
>
> Set `CURSOR_SDK_LOG=debug` (or `info`) to attach a stderr handler to the SDK's own logger. The SDK only configures its own `cursor_sdk` logger, so this won't interfere with the host application's logging setup.
>
> ```bash
> CURSOR_SDK_LOG=debug python my_script.py
> ```
>
> The bundled bridge binary is installed as `cursor-sdk-bridge` on PATH alongside the package. Run it directly to confirm the build shipped with your wheel:
>
> ```bash
> cursor-sdk-bridge --help
> ```
>
> ## Known limitations
>
> - Tool-call payload schemas are intentionally not strongly typed.
> - Inline MCP servers are not persisted across `Agent.resume()`. Pass them again on resume if needed.
> - Custom tools (`local.custom_tools`) and toolset restrictions (`tools`, `disallowed_tools`) are local agents only. The restrictions don't persist on the agent; pass them again on resume.
> - Artifact download is not implemented for local agents.
> - `local.setting_sources` (and the file-based MCP and subagent paths it gates) does not apply to cloud agents. Cloud always loads `project`, `team`, and `plugins`.
> - Hooks are file-based only (`.cursor/hooks.json`). No programmatic callbacks.
>
>

### Source: SDK Bridge

> # Cursor SDK Bridge
>
> The SDK Bridge is a small local server that embeds the TypeScript SDK and exposes the same agent surface over a stable Connect/protobuf protocol. Use it to script Cursor agents from languages without a first-party SDK.
>
> If you write TypeScript or Python, install the first-party [TypeScript](https://cursor.com/docs/sdk/typescript.md) or [Python](https://cursor.com/docs/sdk/python.md) SDK instead. Python already talks to a bundled copy of the bridge.
>
> The protocol, standalone binaries, and adapter guide live in [`cursor/sdk-bridge`](https://github.com/cursor/sdk-bridge). Pin a release, then point a Cursor agent at that repo to build a thin adapter.
>
> Cursor publishes and supports the `sdk.v1` contract and bridge binaries.
> Adapters in other languages are not first-party SDKs. Lead with TypeScript or
> Python unless you need a language those packages do not cover.
>
> ## When to use it
>
> | Path                                                                     | Use when                                                           |
> | :----------------------------------------------------------------------- | :----------------------------------------------------------------- |
> | [TypeScript SDK](https://cursor.com/docs/sdk/typescript.md)              | You're writing TypeScript or JavaScript.                           |
> | [Python SDK](https://cursor.com/docs/sdk/python.md)                      | You're writing Python.                                             |
> | [SDK Bridge](https://github.com/cursor/sdk-bridge)                       | You need Go, Rust, Java, C#, or another language.                  |
> | [Cloud Agents API](https://cursor.com/docs/cloud-agent/api/endpoints.md) | You only need cloud agents over HTTP, with no local agent runtime. |
>
> The bridge is for SDK authors and platform teams. Application code should depend on `@cursor/sdk` or `cursor-sdk`.
>
> ## How it works
>
> ```mermaid
> flowchart LR
>   adapter["Your adapter"]
>   bridge["cursor-sdk-bridge"]
>   api["Cursor API"]
>   adapter -->|"sdk.v1 Connect RPCs"| bridge
>   bridge -->|"HTTPS"| api
>   bridge -->|"tool and store callbacks"| adapter
> ```
>
> Your adapter spawns `cursor-sdk-bridge`, or attaches to one your platform already runs. The bridge binds a loopback HTTP/1.1 port and serves the `sdk.v1` services. Because it embeds `@cursor/sdk`, new agent features land once in the bridge. Adapters pick them up by bumping the binary.
>
> Classic gRPC over HTTP/2 will not connect. Use a [Connect](https://connectrpc.com/docs/protocol) client, or plain `POST`s with protobuf or JSON bodies.
>
> ## Get started
>
> ### Get an API key
>
> SDK runs accept user API keys and service account API keys. Team Admin API keys are not supported yet.
>
> - [Cursor Dashboard → API Keys](https://cursor.com/dashboard/api)
> - [Service accounts](https://cursor.com/docs/account/enterprise/service-accounts.md)
>
> ```bash
> export CURSOR_API_KEY="your-key"
> ```
>
> ### Pin a bridge release
>
> Each GitHub release tag matches the TypeScript and Python SDK version. Download the standalone archive for your platform from [GitHub releases](https://github.com/cursor/sdk-bridge/releases). Each archive unpacks to:
>
> - `bin/cursor-sdk-bridge` (`.exe` on Windows)
> - `proto/sdk/v1/` (the contract for that binary)
> - `manifest.json`
>
> Use `darwin`, `linux`, or `win32` with `x64` or `arm64`. Windows is `x64` only.
>
> The same binary ships inside `cursor-sdk` wheels. After `pip install cursor-sdk`, `cursor-sdk-bridge` is on your `PATH`.
>
> ### Point an agent at the repo
>
> Open Agent and run this prompt. It points Cursor at [`cursor/sdk-bridge`](https://github.com/cursor/sdk-bridge) and the adapter build guide.
>
> Read https\://github.com/cursor/sdk-bridge and follow the Agent: start here guide in the README. Build a thin Cursor SDK adapter in this repository's primary language. Cover codegen from proto/sdk/v1, bridge process lifecycle, streaming, errors, and callback servers.
>
> Confirm a fresh binary before you debug adapter code:
>
> ```bash
> cursor-sdk-bridge --help
> ```
>
> When an RPC fails and your adapter can't see why, run the bridge with `--verbose` (or set `CURSOR_SDK_BRIDGE_LOG=1`) to log each RPC's name, outcome, duration, and full error to stderr. Request and response payloads are never logged.
>
> The repo also has a [curl-only smoke test](https://github.com/cursor/sdk-bridge/blob/main/docs/smoke-test.md) that exercises spawn, `Ping`, `Me`, `CreateAgent`, and `Send` with no adapter code.
>
> ## Adapter shape
>
> An adapter is a library another developer can install without knowing the bridge exists. First-party SDKs converge on this shape:
>
> | Piece                     | Role                                                                                                                    |
> | :------------------------ | :---------------------------------------------------------------------------------------------------------------------- |
> | **Bridge manager**        | Find or spawn the binary, complete the ready-line handshake, and shut it down. Allow attaching to an existing endpoint. |
> | **Transport**             | Connect over HTTP/1.1: unary POSTs and streamed responses, with bearer auth on every call.                              |
> | **Client**                | Low-level typed RPCs for agents, runs, models, and repositories.                                                        |
> | **Agent and Run handles** | The public API: create, send, stream events, wait, and cancel.                                                          |
> | **Errors**                | Map Connect codes and `sdk.v1` error details onto exceptions or result types in your language.                          |
> | **Callback servers**      | Optional loopback servers so users can define custom tools and stores in your language.                                 |
>
> Ship a one-prompt helper (create, send, wait, close) and a context-manager or RAII form so the bridge process cannot leak.
>
> ## Protocol
>
> The wire contract is protobuf package `sdk.v1`:
>
> | Proto                                    | Role                                                                       |
> | :--------------------------------------- | :------------------------------------------------------------------------- |
> | `sdk_agent_service.proto`                | Create and resume agents, send prompts, stream runs, artifacts, and usage. |
> | `sdk_cursor_service.proto`               | Identity, models, and repositories.                                        |
> | `sdk_bridge_control_service.proto`       | Ping, version, shutdown, and tool-callback registration.                   |
> | `sdk_custom_tool_callback_service.proto` | Hosted by your adapter. The bridge calls it to run user-defined tools.     |
> | `sdk_store_callback_service.proto`       | Hosted by your adapter for custom agent stores.                            |
> | `sdk_messages.proto`                     | Shared messages and the run-stream envelope.                               |
> | `sdk_errors.proto`                       | Structured error details.                                                  |
>
> Leave `proto/` untouched when you vendor it. Cursor regenerates those files on every SDK release.
>
> Details stay in the repo:
>
> - [Lifecycle and handshake](https://github.com/cursor/sdk-bridge/blob/main/docs/protocol.md)
> - [Services](https://github.com/cursor/sdk-bridge/blob/main/docs/services.md)
> - [Streaming](https://github.com/cursor/sdk-bridge/blob/main/docs/streaming.md)
> - [Errors](https://github.com/cursor/sdk-bridge/blob/main/docs/errors.md)
> - [Curl smoke test](https://github.com/cursor/sdk-bridge/blob/main/docs/smoke-test.md)
> - [Versioning](https://github.com/cursor/sdk-bridge/blob/main/docs/versioning.md)
>
> ## Authentication
>
> Two separate secrets:
>
> 1. **Cursor API key.** Set `options.api_key` on create, resume, and catalog calls such as `ListModels`. Also export `CURSOR_API_KEY` in the bridge process environment. Catalog calls require the per-call key.
> 2. **Bridge bearer token.** Generated per process during the ready-line handshake. Send `Authorization: Bearer <token>` on every RPC, including streams. The bridge listens on `127.0.0.1` by default.
>
> See [protocol.md](https://github.com/cursor/sdk-bridge/blob/main/docs/protocol.md) for spawn flags, the ready line, and shutdown order.
>
> ## Versioning
>
> `sdk.v1` changes additively. Existing fields are not renumbered or reused. A breaking change would ship as `sdk.v2` alongside `v1`.
>
> Pin codegen to a release tag, and prefer a bridge whose `manifest.json` `sdkVersion` matches. Older adapters keep working against newer bridges. New RPCs stay invisible until you regenerate.
>
> Call `SdkBridgeControlService.GetVersion` when you need to gate on `protocol_version` or `capabilities` at runtime.
>
> ## Support
>
> - **Supported:** the published `sdk.v1` protos, standalone `cursor-sdk-bridge` binaries, and the first-party TypeScript and Python SDKs.
> - **Your responsibility:** community or in-house adapters built on the bridge. You own versioning, support, and security review for those libraries.
>
> SDK runs follow the same pricing, request pools, and Privacy Mode rules as the IDE and Cloud Agents. Spend appears on the [usage dashboard](https://cursor.com/dashboard/usage) under the SDK tag.
>
> ## Related
>
> - [TypeScript SDK](https://cursor.com/docs/sdk/typescript.md)
> - [Python SDK](https://cursor.com/docs/sdk/python.md)
> - [Cloud Agents API](https://cursor.com/docs/cloud-agent/api/endpoints.md)
> - [SDK changelog](https://cursor.com/docs/sdk/changelog.md)
> - [`cursor/sdk-bridge` on GitHub](https://github.com/cursor/sdk-bridge)
>
>

### Source: SDK Changelog

> # SDK Changelog
>
> The latest features, improvements, and fixes shipping to the Cursor SDK, covering `@cursor/sdk` on npm and `cursor-sdk` on PyPI.
>
> ## 1.0.27
>
> - **Restrict the agent's toolset.** `tools` allowlists the built-in tools offered to the model (`[]` means text-only), and `disallowedTools` removes tools while keeping the rest. Both take public names like `"read"` or capability groups like `"shell"` and `"mcp"`, in TypeScript and Python (`tools`, `disallowed_tools`). Local agents only for now, and not persisted across `resume`.
> - **Log in from the browser in TypeScript.** `Cursor.auth.login()` opens a browser login, mints an API key, and stores it in `~/.cursor/sdk/auth.json`; `Cursor.auth.status()` and `Cursor.auth.logout()` round it out. After login, `Agent.create()` and the `Cursor.*` reads work without `apiKey` or `CURSOR_API_KEY`.
> - **Usage and cost for local agents.** `agent.getUsage()` in TypeScript and `agent.get_usage()` in Python now work for local agents too, returning a per-turn breakdown. Pass a `runId` from a previous result to narrow to one turn.
> - **Open PRs as the Cursor GitHub App.** `cloud.openAsCursorGithubApp` in TypeScript and `open_as_cursor_github_app` in Python control PR authorship. Service-account keys default to the app; user keys default to the key's owner.
> - **Multi-root local workspaces.** Pass `local.dirs` to load rules, skills, and project context from several folders; `cwd` stays the single primary working directory. Replaces the `cwd` array form, which only ever used the first entry.
> - **Clearer Python errors.** Failures that previously surfaced as a bare "internal error" now carry the underlying message and code.
> - **Admin command denylists apply to local runs.** Shell commands matching your team's admin denylist are rejected with a policy message before they execute, including on paths that skip approval prompts.
>
> ## 1.0.26
>
> - **Warm up a local workspace before the first send.** `platform.prewarmLocalWorkspace(options)` resolves rules, skills, MCP servers, and ignore files ahead of time, so the first `send()` against that workspace starts immediately. It returns a release function to call on shutdown.
> - **Control how long workspace scans stay cached.** `configureCursorSdk({ local: { workspaceScanCacheTtlMs } })` sets the cache lifetime for workspace scans, and the `CURSOR_RIPWALK_CACHE_TTL_MS` environment variable sets the same value for hosted deployments. Long-lived servers on stable checkouts can now skip repeated re-scans.
> - **Custom tools run without approval prompts.** Host-defined tools passed via `customTools` no longer fail with an interactive-approval error on sandboxed or auto-review local runs. Deny rules and sandbox limits still apply.
> - **Signed macOS binaries.** The `@cursor/sdk` macOS platform packages now ship code-signed binaries, so Gatekeeper and endpoint security tools no longer block them.
> - **Cleaner Python exception hierarchy.** `PermissionDeniedError`, `BadRequestError`, and `InternalServerError` now inherit directly from `CursorSDKError` instead of `AuthenticationError`, `ConfigurationError`, and `NetworkError`, so `except` blocks catch what their names say.
> - **Fixed intermittent startup failures in Python.** Roughly 1 in 64 agent launches failed before reaching the first send. Launches are now reliable.
>
> ## 1.0.25
>
> - **Billed usage and cost on demand.** `agent.getUsage()` in TypeScript and `agent.get_usage()` in Python return token usage, billed cost, and a per-run breakdown for cloud agents, and `Agent.getUsage(agentId)` works without a handle. Cost is server-derived, includes discounts, and settles shortly after a run ends. Cloud-only for now; local runs throw a typed configuration error.
>
> ## 1.0.24
>
> - **TypeScript and Python now release together.** Starting with 1.0.24, `@cursor/sdk` on npm and `cursor-sdk` on PyPI ship from the same release and share a version number. Python releases no longer trail TypeScript.
> - **More reliable long-running streams.** Streaming responses on heavy runs no longer drop mid-stream, which previously surfaced as network errors in Python clients on long turns.
>
> ## 1.0.23
>
> - **Per-send environment variables for cloud runs.** Pass `send(prompt, { cloud: { envVars } })` to scope env vars to a single run, including the first send that creates the agent. `Agent.create({ cloud: { envVars } })` still sets agent-scoped defaults.
> - **Error details on failed runs.** Local and cloud runs that fail now expose a structured error with `message` and `code` fields, so you can tell what went wrong without parsing logs. `run.wait()` behaves the same as before.
> - **Token usage in Python.** Run streams emit typed `usage` messages with per-turn token counts, and cumulative totals are available on `run.usage` and `RunResult.usage`, matching TypeScript from 1.0.22.
> - **Sturdier local run history.** Run history on disk now survives interrupted writes, fixing a class of failures where a crashed process left runs that could not be resumed.
> - **Fixed streaming stalls on Bun.** Run streams under Bun no longer stall on long responses.
>
> ## 1.0.22
>
> - **Token usage on every run.** Local runs emit per-turn `usage` events on `run.stream()` and cumulative totals on `run.wait()`. Cloud runs surface the same usage on their stream and `wait()` results, and totals persist for detached local handles so a process that reattaches still gets them.
>
> ## 1.0.21
>
> - **Run agents under Bun.** `agent.send()` now works under Bun with the same behavior as Node. This also fixes fresh Node installs that could miss a required dependency.
> - **Friendlier runtime names in Python.** List APIs and `get_run` accept `runtime="cloud"`, `"local"`, and `"auto"`, matching the documented values.
>
> ## 1.0.20
>
> - **The SDK imports cleanly under Bun.** Importing `@cursor/sdk` no longer crashes under Bun. Running agents under Bun follows in 1.0.21.
>
>
