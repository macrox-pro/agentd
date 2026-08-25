---
primary_sources:
  - id: T1-CLOUD-CAP
    title: "Capabilities"
    url: "https://cursor.com/docs/cloud-agent/capabilities.md"
    section: "Capabilities"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Cloud capabilities and metadata

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Capabilities

> # Capabilities
>
> ## Computer use
>
> Each cloud agent runs in its own isolated VM with a full desktop environment. Agents can use a mouse and keyboard to control the desktop and browser, allowing them to interact with the software they build like a human developer.
>
> This means agents can start dev servers, open the app in a browser, click through UI flows, and verify their changes work before pushing a PR. Read more in the [announcement blog post](/blog/agent-computer-use).
>
> ## Demos and Artifacts
>
> Agents create artifacts such as screenshots, videos, and log references to demonstrate their work. These artifacts are attached to the PR so you can quickly validate changes without checking out the branch locally.
>
> ### Artifacts in GitHub
>
> You can opt-in to have Cloud Agents embed artifacts directly into GitHub pull request descriptions by enabling the **Allow posting artifacts to GitHub** setting in the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents#my-pull-requests).
>
> GitHub's [image proxy](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-anonymized-urls) requires public URLs, so artifacts in PR descriptions use long, unguessable URLs that are viewable without authentication. For context, GitHub used public URLs for all issue and PR attachments until [May 2023](https://github.blog/changelog/2023-05-08-more-secure-private-attachments).
>
> ## Remote desktop control
>
> You can take control of the agent's remote desktop to interact with the software the agent is building. Hand control back to the agent at any time to let it keep working.
>
> Cloud agents run in a remote VM that can be fully onboarded with your repo, dependencies, tooling, and setup scripts. This allows you to test changes directly in the agent's VM without checking out the branch on your local machine.
>
> ## MCP tools
>
> Cloud agents can use [MCP (Model Context Protocol)](https://cursor.com/docs/mcp.md) servers configured for your team. This gives agents access to external tools and data sources like databases, APIs, and third-party services during their runs.
>
> Add and enable personal MCP servers through the MCP dropdown in [cursor.com/agents](https://cursor.com/agents). Team admins configure shared servers under **Dashboard -> Integrations & MCP**.
>
> Admins can link shared Team MCP servers to the [Default team marketplace](https://cursor.com/docs/plugins.md#migrate-existing-team-mcps). Linking keeps the servers available to Cloud Agents and also makes them available for teammates to install and configure in the Agent Window, IDE, and CLI.
>
> Cloud agents support OAuth for MCP servers that need it. OAuth is per-user, including for MCP servers shared at the team level.
>
> ### Custom MCP servers
>
> You can add custom MCP servers using either **HTTP** or **stdio** transport. SSE and `mcp-remote` are not supported.
>
> MCP configurations are encrypted at rest. Sensitive fields are redacted and cannot be read back by any user after saving:
>
> - **`env`** — environment variables for stdio servers
> - **`headers`** — request headers for HTTP servers
> - **`CLIENT_SECRET`** — OAuth client secret for HTTP servers
>
> ### HTTP vs stdio
>
> - **HTTP (recommended)** — server configurations are never present in the cloud agent's VM environment. The agent does not have access to refresh tokens, headers, or other credentials. Tool calls are proxied through the backend.
> - **Stdio** — servers run inside the cloud agent's VM, so the agent has access to the server's configuration and environment variables. This is similar to how stdio MCPs work in the Cursor IDE.
>
> Stdio servers depend on the VM environment to execute. We cannot verify that a stdio server will run successfully until a cloud agent is launched. We recommend using HTTP MCPs when possible, and configuring your [environment setup](https://cursor.com/docs/cloud-agent/setup.md) correctly if you use stdio servers.
>
> ### Cursor Cloud MCP
>
> The Cursor Cloud MCP is a built-in diagnostics server available during Cloud Agent runs. It lets an agent inspect the current run, browse related runs in the same environment, and fetch transcripts, diff metadata, environment details, run events, and setup logs without manually collecting links and files.
>
> Team admins can disable Cursor Cloud MCP for their team from **MCP Configuration** in [team settings](https://cursor.com/dashboard/settings). See [Team dashboard](https://cursor.com/docs/account/teams/dashboard.md#integrations) for more on MCP admin controls.
>
> #### Access and permissions
>
> Cloud Agent conversations can include prompts, code, tool output, and secrets. All tools enforce access checks on every request.
>
> | Role       | What you can access                                                                                                                                |
> | :--------- | :------------------------------------------------------------------------------------------------------------------------------------------------- |
> | Team admin | List and fetch details (including transcripts) for Cloud Agent runs across the team, for repositories and environments they already have access to |
> | Non-admin  | Only your own runs and transcripts. You cannot view other team members' chats through this MCP                                                     |
>
> Even when listing runs in a shared environment, non-admins only see agents they started or own. Service accounts follow the same rules as the user or team context they run under.
>
> #### What you can inspect
>
> | Category      | Examples                                                                                                                                                                                                                         |
> | :------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | Current run   | Run ID, URL, repo, branch, model, owner, lifecycle status, and where the run was started (Cursor, Slack, GitHub, API, and others)                                                                                                |
> | Events        | Setup, pull request, artifact, and MCP authentication outcomes shown on the run dashboard. Use `get-events` for the current run, or `batch-fetch-details` with `include_events` for other runs. See the event kinds under Tools. |
> | Related runs  | Other Cloud Agents in the same environment, or on the same repository when no saved environment is attached                                                                                                                      |
> | Environment   | Environment version, full environment config, dashboard URL, and effective egress network policy                                                                                                                                 |
> | Transcript    | Full user-agent conversation, including tool calls when available                                                                                                                                                                |
> | Diff metadata | Whether the agent changed code, how much changed, and whether it opened a PR                                                                                                                                                     |
> | Setup logs    | Raw logs from environment setup and image-build steps                                                                                                                                                                            |
>
> #### Tools
>
> Depending on your MCP client, tool names may include a server prefix (for example, `cursor-cloud-run-info`). The underlying tools are:
>
> | Tool                                | Purpose                                                                                                                                                                                                                 |
> | :---------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `run-info`                          | Get the current run's identity, metadata, and URL. Start here.                                                                                                                                                          |
> | `environment-info`                  | Get the current run's environment version, config, dashboard URL, and effective egress policy.                                                                                                                          |
> | `get-events`                        | List the current run's dashboard events, oldest first.                                                                                                                                                                  |
> | `list-cloud-agents`                 | Browse Cloud Agent runs visible to you in this environment. Filter by source, status, date, code changes, PR creation, and archived state.                                                                              |
> | `batch-fetch-details`               | Fetch details for specific run IDs (`bcId`s). Optionally include transcripts, diff metadata, setup logs, environment info, and run events via `include_events` (writes `events.json` per run; up to 50 runs per batch). |
> | `get-automation`                    | Get an automation's details like name and owner from its ID.                                                                                                                                                            |
> | `list-environment-builds`           | List recent [Builds](https://cursor.com/docs/cloud-agent/builds.md) for the current environment and inspect their status.                                                                                               |
> | `environment-build-logs`            | Download the install and setup logs for a Build.                                                                                                                                                                        |
> | `trigger-environment-build`         | Run a test Build with the current configuration or proposed install and start commands.                                                                                                                                 |
> | `propose-environment-json`          | Present install and start commands for you to review before saving the environment.                                                                                                                                     |
> | `take-environment-snapshot`         | Snapshot a machine after the agent verifies its environment setup.                                                                                                                                                      |
> | `check-environment-snapshot`        | Check whether an environment snapshot is ready.                                                                                                                                                                         |
> | `request-environment-setup-actions` | Request user actions that block environment setup, such as adding a secret.                                                                                                                                             |
>
> Dashboard event `kind` values from `get-events` and from `batch-fetch-details` with `include_events` are:
>
> | `kind`               | Meaning                                                                          |
> | :------------------- | :------------------------------------------------------------------------------- |
> | `setup_started`      | Environment setup began.                                                         |
> | `setup_completed`    | Environment setup finished.                                                      |
> | `setup_failed`       | Environment setup failed.                                                        |
> | `pr_created`         | Pull request opened.                                                             |
> | `pr_creation_failed` | Pull request creation failed.                                                    |
> | `artifact_created`   | Walkthrough artifact uploaded.                                                   |
> | `mcp_auth_error`     | MCP server authentication failed; its tools were skipped, and the run continued. |
>
> A typical diagnostics flow is `run-info` → `get-events` → `environment-info` → `list-cloud-agents` → `batch-fetch-details` (set `include_events` when you need other runs' dashboard events).
>
> ## Subscriptions
>
> Agent tasks rarely end with the last commit. CI has to pass. Reviewers leave comments. A teammate needs to answer a question in Slack. Subscriptions let a cloud agent wait for those events and keep working when they happen, without you re-prompting it.
>
> The agent subscribes to an event source, ends its turn, and wakes when a matching event arrives. Events land as follow-ups in the same conversation, so the agent continues with full context:
>
> - Open a PR, then respond to review comments and CI failures until it merges
> - Ask a question in Slack and continue once someone replies
> - Check back on a long-running job with a timer
>
> To subscribe, describe the wait in your prompt. For example, "open a PR and keep it green until merge" or "ask in #releases and wait for approval". You can also invoke the built-in `/subscribe` skill, which works the same way: tell it what to watch and the agent picks the right subscription.
>
> Agents can subscribe to events from these integrations:
>
> | Integration | Events                                                                                                                                                                                                                     |
> | :---------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | GitHub      | Pull request activity (comments, reviews, and lifecycle changes) for one PR, a whole repo, or one author's PRs, and CI results on a branch. Uses the [GitHub integration](https://cursor.com/docs/integrations/github.md). |
> | Slack       | Replies in a thread, messages in a channel, and newly created public channels. Uses the [Slack integration](https://cursor.com/docs/integrations/slack.md).                                                                |
> | Linear      | Issues created or changing state, and new comments on issues. Uses the [Linear integration](https://cursor.com/docs/integrations/linear.md).                                                                               |
> | Timers      | A point in time: a one-off reminder after a delay, or a recurring cron schedule. Recurring loops are also available as the built-in [`/loop`](https://cursor.com/docs/skills.md#built-in-cursor-skills) skill.             |
>
> ### How subscriptions work
>
> - Subscriptions belong to a single agent conversation. Events wake that agent as follow-up messages.
> - Bursts coalesce. Several events arriving close together can wake the agent once, and the agent re-reads the source (the PR, thread, or issue) before acting.
> - A subscription lasts at most 180 days. Agents also unsubscribe on their own when the wait is over.
>
> ## Fixing CI Failures
>
> Cloud Agents automatically try to fix CI failures in PRs they create. This currently supports GitHub Actions only.
>
> Cloud Agents skip automatic CI follow-ups if:
>
> - You've pushed a new commit to the branch; cloud agents do not auto-fix CI failures on human commits.
> - You've sent a follow-up message to the agent.
> - The same check is already failing on the base commit of the PR.
> - The PR has already had 10 CI-failure follow-ups.
>
> To disable this feature on all your personal Cloud Agents, go to [Cursor Dashboard → Cloud Agents → My Settings](https://cursor.com/dashboard/cloud-agents) and disable the "Automatically fix CI Failures" option.
>
> To disable this feature on a specific Cloud Agent PR, you can comment `@cursor autofix off` on the PR. To re-enable it, comment `@cursor autofix on`.
>
> If you want cloud agents to fix CI failures in your own PRs, you can simply ask them by tagging Cursor in a comment as normal. For example, `@cursor please fix the CI failures`, or `@cursor fix the CI lint check failure`.
>
> Automatically fixing CI failures is currently only available on Teams; support for non-Teams accounts is coming soon. In the meantime, if you want similar behavior, you can ask the cloud agent explicitly to monitor and fix CI failures on the PR.
>
> ## OIDC identity tokens
>
> Cursor-managed Cloud Agent VMs can mint short-lived OIDC JWTs from a local socket. Agents use them to assume cloud roles or call internal APIs without storing long-lived keys. See [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md).
>
> ## Agent metadata
>
> The same socket also serves [agent metadata](https://cursor.com/docs/cloud-agent/metadata.md). Agents, hooks, and scripts can read the agent id, owner, current turn, and workspace as plain text.
>
>

### Source: Metadata

> # Agent metadata
>
> ### Preview
>
> Agent metadata is in preview and subject to change, including breaking
> changes.
>
> Cloud Agents can read key-value metadata about the current run from inside the VM: the agent id, who owns it, who submitted this turn, which model is serving, and which repos are checked out. Hooks and install scripts can read those values too.
>
> Agents call this API with their terminal tools. You don't need to run these requests yourself.
>
> To have an agent read metadata, include this in your prompt:
>
> ```text
> To read agent metadata, follow the instructions at
> https://cursor.com/docs/cloud-agent/metadata
> ```
>
> This API is local to the agent VM. It is not the caller-owned `metadata` tags you set when creating an agent with the [SDK](https://cursor.com/docs/sdk/typescript.md#agent-metadata) or [Cloud Agents API](https://cursor.com/docs/cloud-agent/api/endpoints.md). Those APIs use Cursor API keys and manage agents from outside the VM.
>
> When something outside the VM needs to verify the agent's identity, have the agent mint an [OIDC token](https://cursor.com/docs/cloud-agent/identity.md) instead. Those JWTs are signed and audience-bound. Metadata is not a credential. It can include the current turn's submitter and serving model, which a token shouldn't carry.
>
> Cursor-managed Cloud Agent VMs serve metadata on the same socket as OIDC tokens. Self-hosted workers do not serve this API yet.
>
> ## Read a value
>
> The agent reads keys over the Unix socket at `CURSOR_AGENT_SOCKET`. On Cursor-managed VMs the default is `/run/cursor/api.sock`.
>
> ```bash
> curl --unix-socket "${CURSOR_AGENT_SOCKET:-/run/cursor/api.sock}" \
>   http://cursor-agent/v1/meta-data/agent/id
> ```
>
> Requests are HTTP over a Unix socket. The hostname in the URL is ignored.
>
> List a prefix to see which keys exist:
>
> ```bash
> curl --unix-socket "${CURSOR_AGENT_SOCKET:-/run/cursor/api.sock}" \
>   http://cursor-agent/v1/meta-data/
> ```
>
> ```text
> agent/
> owner/
> turn/
> workspace/
> ```
>
> Then request a key:
>
> ```bash
> curl --unix-socket "${CURSOR_AGENT_SOCKET:-/run/cursor/api.sock}" \
>   http://cursor-agent/v1/meta-data/owner/user-id
> ```
>
> ### Request
>
> `GET /v1/meta-data[/<path>]` over the Unix socket. No body or extra headers. Trailing slashes are allowed, so a listed `agent/` can be requested as `/v1/meta-data/agent/`.
>
> A missing key returns `404`.
>
> ### Response
>
> Successful reads are `text/plain; charset=utf-8`. A key response is the value as text, with nothing else.
>
> | Kind   | Body                                                                                                |
> | :----- | :-------------------------------------------------------------------------------------------------- |
> | Key    | The value as a string. Keys with several values put one entry per line.                             |
> | Prefix | One child per line, sorted. Nested prefixes end with `/`. The listing ends with a trailing newline. |
>
> Error responses are JSON. See [Rate limits and errors](https://cursor.com/docs/cloud-agent/metadata.md#rate-limits-and-errors).
>
> ### When keys appear
>
> [Install scripts](https://cursor.com/docs/cloud-agent/setup.md) can read the same socket. A key is present only when it has a value: `turn/` is missing until a coding turn starts, and `workspace/branch-name` is missing until the run records a branch. Owner, team, and repository keys are available from agent creation onward.
>
> If the socket is missing right after boot, retry the connection.
>
> ## Keys
>
> Missing keys are omitted from listings and return `404` if requested directly. A listing only includes keys that exist right now.
>
> ### `agent/`
>
> | Key             | When present | Description                                                                     |
> | :-------------- | :----------- | :------------------------------------------------------------------------------ |
> | `agent/id`      | Always       | Cloud Agent id (`bcId`).                                                        |
> | `agent/name`    | When known   | Name shown in the dashboard.                                                    |
> | `agent/source`  | When known   | How the agent was started, such as `WEBSITE`, `API`, `SLACK`, or `AUTOMATIONS`. |
> | `agent/runtime` | Always       | `managed` on Cursor-managed Cloud Agent VMs.                                    |
>
> ### `owner/`
>
> | Key                        | When present | Description                                                                                    |
> | :------------------------- | :----------- | :--------------------------------------------------------------------------------------------- |
> | `owner/user-id`            | When known   | Cursor user id of the agent owner, as a decimal string. Prefer this over email for allowlists. |
> | `owner/user-email`         | When known   | Lowercased owner email. Email can change.                                                      |
> | `owner/service-account-id` | When known   | Service account id when a service account owns the agent.                                      |
> | `owner/team-id`            | When known   | Owning team id, as a decimal string.                                                           |
>
> ### `turn/`
>
> `turn/` exists only while a coding turn is active. Between turns those keys are gone. If `turn/` is missing, there is no active turn.
>
> Values under `turn/` always reflect the current turn. Don't cache them across turns.
>
> | Key               | When present  | Description                                                                                                                                                                                             |
> | :---------------- | :------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `turn/id`         | During a turn | Id of this coding turn. Different from `agent/id`, which is the Cloud Agent id (`bcId`).                                                                                                                |
> | `turn/user-id`    | When known    | Cursor user id of the person who submitted this turn, as a decimal string. On a [team follow-up](https://cursor.com/docs/cloud-agent/settings.md#team-follow-ups) this can differ from `owner/user-id`. |
> | `turn/user-email` | When known    | Lowercased email of that person.                                                                                                                                                                        |
> | `turn/started-at` | During a turn | Turn start as Unix seconds.                                                                                                                                                                             |
> | `turn/model`      | When known    | Model serving this turn. If you selected Auto, this is the model that served, not `Auto`.                                                                                                               |
>
> OIDC tokens don't include who submitted the turn or which model is serving, because a token can outlive the turn. Read those keys from metadata instead.
>
> ### `workspace/`
>
> | Key                        | When present          | Description                                                                                                                                                                                                                |
> | :------------------------- | :-------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `workspace/repo-url`       | When known            | Primary repository in `host/path` form, such as `github.com/acme/widgets`. Hostname is lowercased, with no scheme, credentials, port, query, or `.git` suffix. On a multi-repo agent, this is only the primary repository. |
> | `workspace/repo-urls`      | When the set is known | Every repository in the workspace, same form as `repo-url`. Primary repository first, then the rest sorted, one URL per line. Missing means the set isn't known, not that there is only one repo.                          |
> | `workspace/branch-name`    | When known            | Branch on the primary repository.                                                                                                                                                                                          |
> | `workspace/environment-id` | When known            | Id of the Cursor environment this run used.                                                                                                                                                                                |
> | `workspace/automation-id`  | For automations       | Automation id when `agent/source` is automations.                                                                                                                                                                          |
>
> `workspace/repo-url` is the primary repository. For the full set, read `workspace/repo-urls`.
>
> ## Who can read metadata
>
> Any process that can reach the socket can read every key: the agent, code it runs, and hooks. Treat these values as visible to the whole run.
>
> Metadata is not signed. To prove identity to AWS, GCP, Vault, or your own service, have the agent mint an [OIDC token](https://cursor.com/docs/cloud-agent/identity.md) and verify the JWT. Don't forward metadata values as a credential.
>
> ## Rate limits and errors
>
> Each agent VM can make **120** metadata requests per minute, in bursts of up to 20. The socket also accepts at most 8 connections at once. That cap is shared with OIDC minting.
>
> Retry `429`, `503`, `500`, `502`, and `504` with backoff. Treat `403` as fatal: this agent isn't allowed to read metadata.
>
> `404` and `405` responses include a `usage` string that restates how to call the API. Rate-limit and saturation errors stay code-only:
>
> ```json
> { "error": "not_found", "usage": "GET /v1/meta-data[/<path>] ..." }
> ```
>
> ```json
> { "error": "rate_limited" }
> ```
>
> | HTTP      | `error`               | When                                                            |
> | :-------- | :-------------------- | :-------------------------------------------------------------- |
> | 404       | `not_found`           | Unknown or missing key                                          |
> | 405       | `method_not_allowed`  | Not `GET`                                                       |
> | 429       | `rate_limited`        | Over the per-agent request budget; honor `Retry-After`          |
> | 503       | `saturated`           | Too many connections; honor `Retry-After`                       |
> | 500       | `host_error`          | Internal error; retry                                           |
> | 502 / 504 | `backend_unreachable` | Cursor couldn't return metadata; retry                          |
> | Other     | `backend_error`       | Cursor rejected the request. `403` is fatal; `503` is retryable |
>
> ## Examples
>
> An agent or hook can compare the turn submitter to the owner. A teammate's follow-up can take a stricter path:
>
> ```bash
> SOCKET="${CURSOR_AGENT_SOCKET:-/run/cursor/api.sock}"
>
> owner="$(curl -fsS --unix-socket "$SOCKET" \
>   http://cursor-agent/v1/meta-data/owner/user-id)"
> turn_user="$(curl -fsS --unix-socket "$SOCKET" \
>   http://cursor-agent/v1/meta-data/turn/user-id || true)"
>
> if [ -n "$turn_user" ] && [ "$turn_user" != "$owner" ]; then
>   echo "follow-up from user $turn_user; owner is $owner"
> fi
> ```
>
> An agent or hook can tag logs with the agent id and the model that served the turn:
>
> ```bash
> SOCKET="${CURSOR_AGENT_SOCKET:-/run/cursor/api.sock}"
>
> agent_id="$(curl -fsS --unix-socket "$SOCKET" \
>   http://cursor-agent/v1/meta-data/agent/id)"
> model="$(curl -fsS --unix-socket "$SOCKET" \
>   http://cursor-agent/v1/meta-data/turn/model || true)"
>
> echo "cloud_agent_id=$agent_id model=${model:-unknown}"
> ```
>
> List every repository in the workspace. `repo-urls` is one URL per line:
>
> ```bash
> curl -fsS --unix-socket "${CURSOR_AGENT_SOCKET:-/run/cursor/api.sock}" \
>   http://cursor-agent/v1/meta-data/workspace/repo-urls
> ```
>
> ```text
> github.com/acme/widgets
> github.com/acme/docs
> ```
>
> ## Related pages
>
> - [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md) for signed JWTs and cloud federation
> - [Secrets & Network](https://cursor.com/docs/cloud-agent/security-network.md) for dashboard secrets and egress controls
> - [Cloud agent setup](https://cursor.com/docs/cloud-agent/setup.md) for install scripts that can read this socket
> - [Hooks](https://cursor.com/docs/hooks.md) for running this API at tool and conversation boundaries
> - [Service accounts](https://cursor.com/docs/account/enterprise/service-accounts.md) when agents run as a team service account
>
>
