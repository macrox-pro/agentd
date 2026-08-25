---
primary_sources:
  - id: T1-CLOUD-SELF
    title: "Full page"
    url: "https://cursor.com/docs/cloud-agent/self-hosted-agents.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Self-hosted agents

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Self-hosted agents

> # Cursor Self-Hosted Agents
>
> Run Cursor agents on your own infrastructure. Your team still works in the Cursor app, on [cursor.com](https://cursor.com/agents), and in the Cursor mobile app. Cursor still runs the agent loop, including inference and planning. The agent gets work done by executing on a machine you own instead of in Cursor's managed cloud.
>
> - Your team keeps every Cursor surface: native app, web, and mobile.
> - Cursor owns the agent loop: inference, planning, and orchestration.
> - Tool calls run on your machines, inside your network.
>
> Cursor-managed Cloud Agents are the recommended path for most teams and the
> fastest way to get started. See [Choose where Cloud Agents
> run](https://cursor.com/docs/cloud-agent/self-hosted-guides/choose-runtime.md) before taking on
> your own worker fleet.
>
> ## Who should use Self-Hosted Agents
>
> Use Self-Hosted Agents when Cursor's managed cloud can't meet your constraints:
>
> - You have strict network requirements, and code or services can't be reached from outside your network.
> - You have custom hardware, such as GPU machines or Macs for iOS development.
>
> ## How it works
>
> | Term                   | Definition                                                                                                                                                                                     | Example                                                                                                                                                    |
> | :--------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | **Self-hosted worker** | Your machine or VM, a place where the Cursor agent can get work done: editing files, running commands, and accessing code.                                                                     | A Linux VM in your AWS account, or a Mac mini on your desk.                                                                                                |
> | **Pool**               | A routing target you can select in the Cursor client UI. Chats wait in the pool until a worker claims them. Once a worker claims a chat, all activity in that chat is forwarded to the worker. | A `gpu` pool routes requests that need GPUs, served only by machines with GPUs. An `ios` pool is served only by Macs for chats related to iOS development. |
> | **Orchestrator**       | Code you run that adjusts worker capacity based on demand.                                                                                                                                     | A request arrives on a pool with no idle workers. Your orchestrator notices and starts a new machine.                                                      |
>
> A self-hosted worker is a machine with the Cursor CLI that has run `agent worker start`. The CLI opens a long-lived outbound HTTPS connection to Cursor's backend, and Cursor sends agent tool calls over that connection. Cursor never connects into your network: the outbound connection from your machine to Cursor is all that's required.
>
> Self-hosted workers come in two configurations:
>
> 1. **My Machines.** Connect a single machine you own. Multiple agents can run on the same machine. My Machines requires a personal credential: browser login or a personal user API key. See the [quickstart](https://cursor.com/docs/cloud-agent/self-hosted-agents/quickstart.md).
> 2. **Pools.** Register machines under a pool name, and Cursor routes each new chat to an available machine in the pool, one agent per machine. Run an orchestrator to scale the pool up and down. Pools require a Cursor Enterprise plan and a [service account API key](https://cursor.com/docs/account/enterprise/service-accounts.md). See [Pools](https://cursor.com/docs/cloud-agent/self-hosted-agents/pools.md).
>
> ## What leaves your network
>
> Two things leave your network: file chunks the model reads during inference, and Cloud Agent artifacts (screenshots, videos, and log references) the worker uploads to Cursor-managed storage so they can appear in PRs and the dashboard. Your repos, build caches, and secrets stay on your machines.
>
> Self-hosted workers need outbound HTTPS access to:
>
> - `api2.cursor.sh` and `api2direct.cursor.sh` for the agent session
> - `cloud-agent-artifacts.s3.us-east-1.amazonaws.com` for artifact uploads
>
> No inbound ports, public IPs, or VPN tunnels are required. If you use a proxy, set `HTTPS_PROXY` or `https_proxy` in the worker environment.
>
> Self-Hosted Agents supports up to 10 workers per user and 50 per team. For
> larger company-wide deployments, [contact
> us](https://cursor.com/contact-sales?source=self-hosted-agents) to discuss
> scaling.
>
> ## Get started
>
> - [Quickstart](https://cursor.com/docs/cloud-agent/self-hosted-agents/quickstart.md): connect your first self-hosted worker in a few minutes.
> - [Pools](https://cursor.com/docs/cloud-agent/self-hosted-agents/pools.md): organize workers into pools for your team.
> - [Pool orchestration](https://cursor.com/docs/cloud-agent/self-hosted-agents/orchestration.md): scale worker capacity with the fleet management API.
> - [API reference](https://cursor.com/docs/cloud-agent/api/endpoints.md#fleet-management): endpoints for workers, pools, the pending-request queue (list, SSE watch, and claim), and worker tokens.
>
>
