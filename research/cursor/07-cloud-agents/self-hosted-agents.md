---
primary_sources:
  - id: T1-CLOUD-SELF
    title: "Full page"
    url: "https://cursor.com/docs/cloud-agent/self-hosted.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-09-03"
cursor_docs_snapshot: "2026-09-03"
applicability: "current"
---
# Self-hosted machines

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-09-03).

### Source: Self-Hosted Machines

> # Self-Hosted Machines
>
> Self-Hosted Machines moves Cloud Agent tool execution to a machine you manage. Your team still uses Cloud Agents through the desktop app, [cursor.com](https://cursor.com/agents), and [mobile](https://cursor.com/docs/cloud-agent/mobile.md). Cursor runs the agent loop, inference, and planning. Your worker performs file edits and terminal commands. It also runs computer-use tools and local MCP servers.
>
> Cursor-managed Cloud Agents are the recommended path for most teams and the
> fastest way to get started. See [Choose where Cloud Agents
> run](https://cursor.com/docs/cloud-agent/self-hosted/choose-runtime.md) before bringing
> Cloud Agents to your own machines.
>
> ## Who should use Self-Hosted Machines
>
> Use Self-Hosted Machines when Cursor's managed cloud can't meet your constraints:
>
> - You have strict network requirements, and code or services can't be reached from outside your network. For private source control or package registries, start with managed Cloud Agents and [private connectivity](https://cursor.com/docs/cloud-agent/private-connectivity.md) ([AWS PrivateLink](https://cursor.com/docs/cloud-agent/private-connectivity.md#aws-privatelink) or [Cloudflare Tunnel](https://cursor.com/docs/cloud-agent/private-connectivity.md#cloudflare-tunnel)).
>
> - You have custom hardware, such as GPU machines or Macs for iOS development. Use a machine you already run, or a VM from a [partner host](https://cursor.com/docs/cloud-agent/self-hosted/integrations.md) such as AWS Lambda, Cloudflare, Namespace, Modal, Daytona, E2B, or Vercel.
>
> - You have custom images, such as a different operating system or an existing build pipeline, that are difficult to save as a [Cloud Agent build](https://cursor.com/docs/cloud-agent/builds.md).
>
> If you want to try this out, check out the [My Machines
> quickstart](https://cursor.com/docs/cloud-agent/self-hosted/my-machines.md#quickstart).
>
> ### Prefer managed Cloud Agents when
>
> You don't need to own the compute to keep agents inside your security perimeter. Stay on Cursor-managed Cloud Agents when these controls cover your requirements:
>
> - Isolated VMs per agent, provisioned and torn down by Cursor, with no worker fleet to size, patch, scale, or keep on call.
>
> - [Network allowlists](https://cursor.com/docs/cloud-agent/security-network.md#network-access) that restrict outbound domains by user, team, or environment.
>
> - [Private connectivity](https://cursor.com/docs/cloud-agent/private-connectivity.md) over AWS PrivateLink or Cloudflare Tunnel to self-hosted GitHub Enterprise Server, GitLab Enterprise, and private source control APIs, plus [Tailscale or a similar client](https://cursor.com/docs/cloud-agent/security-network.md#private-network-access) for services in your VPC.
>
> - Privacy Mode and customer-controlled secrets scoped to the environments you choose. See [Cloud Agent security and network](https://cursor.com/docs/cloud-agent/security-network.md) for the full model.
>
> ## What leaves your network
>
> The full checkout, build cache, and machine-local credentials stay on your machine. During a run, the worker sends Cursor the content the agent needs, such as file contents, terminal output, diffs, screenshots, local MCP results, and routing metadata. If you enable [desktop sharing](https://cursor.com/docs/cloud-agent/self-hosted/computer-use.md#share-the-agent-desktop), it also streams the agent desktop.
>
> The worker uploads Cloud Agent artifacts, such as screenshots, videos, and log references, to Cursor-managed storage so they can appear in pull requests and the dashboard. Keep secrets out of tool output and artifacts.
>
> [Privacy Mode](/data-use) also applies to Self-Hosted Machines. When enabled,
> code sent from the worker is not used for training by Cursor or model
> providers.
>
> Workers need outbound HTTPS access to:
>
> - `api2.cursor.sh` and `api2direct.cursor.sh` for the agent session
>
> - `cloud-agent-artifacts.s3.us-east-1.amazonaws.com` for artifact uploads
>
> No inbound ports, public IPs, or VPN tunnels are required. If you use a proxy, set `HTTPS_PROXY` or `https_proxy` in the worker environment.
>
> To disable artifact uploads, block outbound traffic to `cloud-agent-artifacts.s3.us-east-1.amazonaws.com` on the worker. This only prevents artifacts from uploading. The agent continues to work, including tool calls and results, but its artifacts won't appear in pull requests or the dashboard.
>
> You can connect up to 200 workers per user and 1000 per team. For
> larger company-wide deployments, [contact
> us](https://cursor.com/contact-sales?source=self-hosted-agents) to discuss
> scaling.
>
> ## How it works
>
> **Worker****Team Pool**`gpu``ios`**Controller**
>
> | Term | Definition | Example |
> | --- | --- | --- |
> |  | A machine you own, registered to Cursor with the Cursor CLI. The place where the agent gets work done: editing files, running commands, and accessing code. | A Linux VM in your AWS account, or a Mac mini on your desk. |
> |  | A routing target you can select in the Cursor client UI. Chats wait in the team pool until a worker claims them. Once a worker claims a chat, all activity in that chat is forwarded to the worker. | A  team pool routes requests that need GPUs, served only by machines with GPUs. An  team pool is served only by Macs for chats related to iOS development. |
> |  | Code you run that adjusts worker capacity based on demand. | A request arrives on a team pool with no idle workers. Your controller notices and starts a new machine. |
>
> Use the Cursor CLI to start a worker on your machine. `agent worker start` opens a long-lived outbound HTTPS connection to Cursor's backend, and Cursor sends agent tool calls over that connection. Cursor never connects into your network: the outbound connection from your machine to Cursor is all that's required.
>  Forward Return Inference (LLM) loop
>
> Workers come in two configurations:
>
> 1. **My Machines.** Best for personal workflows and one-offs: your devbox, a spare VM, or a machine with state you don't want to recreate. Connect a machine to your account. Multiple agents can run on the same machine. See [My Machines](https://cursor.com/docs/cloud-agent/self-hosted/my-machines.md).
>
> 1. **Team Pools.** Best for teams and enterprises: shared capacity, service account authentication, and centrally managed images. Register machines under a pool name, and Cursor routes each new chat to an available machine in the pool, one agent per machine. Run a controller to scale the pool up and down. See [Team Pools](https://cursor.com/docs/cloud-agent/self-hosted/pool.md).
>
> To run Team Pool workers on a third-party VM or sandbox, see [Integrations](https://cursor.com/docs/cloud-agent/self-hosted/integrations.md).
>
> ## Supported deployment patterns
>
> Run a worker anywhere you can install the Cursor CLI and its dependencies:
>
> - **Personal machines.** Connect a laptop, devbox, Mac, or remote VM through [My Machines](https://cursor.com/docs/cloud-agent/self-hosted/my-machines.md).
>
> - **Persistent hosts or containers.** Run one or more pool workers under `systemd`, `launchd`, Docker, or another process manager.
>
> - **Dynamic infrastructure.** Use the built-in [worker controller](https://cursor.com/docs/cloud-agent/self-hosted/pool.md#worker-controller) or the Cloud Agents API to start machines when requests arrive.
>
> - **Kubernetes.** Follow the [Kubernetes](https://cursor.com/docs/cloud-agent/self-hosted/kubernetes.md) guide to run pool workers with the Cursor operator and Helm chart.
>
> - **Partner hosts and templates.** Run pool workers on AWS Lambda, Cloudflare, Namespace, Modal, Daytona, E2B, Vercel, or Coder with partner guides, or clone a Cursor reference template for AWS Lambda MicroVMs, Cloudflare Containers, or Kubernetes. See [Integrations](https://cursor.com/docs/cloud-agent/self-hosted/integrations.md).
>
> Deployment guides, partner guides, and templates are reference architectures. You own the worker image, infrastructure, secrets, scaling policy, and production validation.
>
> ## Cost
>
> Every runtime option uses the selected model and follows its [pricing](https://cursor.com/docs/models-and-pricing.md#model-pricing). Cursor-managed Cloud Agents include the execution infrastructure. With Self-Hosted Machines, you also pay for and operate your machines, containers, or cluster.
>
> ## Requirements
>
> Every worker needs the [Cursor CLI](https://cursor.com/docs/cli/overview.md) and outbound HTTPS access. Install the CLI on each machine:
>
> ```
> curl https://cursor.com/install -fsS | bash
> ```
>
> **My Machines**
>
> - 
>
> A personal credential: browser login, or a personal user API key from [Cursor Dashboard → API Keys](https://cursor.com/dashboard/api).
>
> ```
> agent login
> ```
>
> **Team Pools**
>
> - 
>
> A **Cursor Enterprise plan**.
>
> - 
>
> A [service account API key](https://cursor.com/docs/account/enterprise/service-accounts.md) for worker authentication. Personal logins and other API key types can't start pool workers.
>
> ```
> export CURSOR_API_KEY="<team service-account API key>"
> ```
>
> - 
>
> Self-hosted settings configured by a team admin in the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents#self-hosted-agents): **Allow Self-Hosted Machines** lets users opt in, and **Require Self-Hosted Machines** routes every Cloud Agent run to your workers.
>
> **Computer use** (optional)
>
> - 
>
> A macOS worker with a signed-in desktop session. The CLI installs the **Cursor Computer Use** helper app on first start; grant it **Accessibility** and **Screen Recording** in System Settings.
>
> - 
>
> Or a Linux worker with the desktop packages. The worker creates its own desktop or reuses an existing X11 display:
>
> ```
> sudo apt-get install -y --no-install-recommends \  dbus-x11 ffmpeg tigervnc-standalone-server \  x11-utils x11-xserver-utils xdotool xfce4
> ```
>
> See [Computer use and desktop sharing](https://cursor.com/docs/cloud-agent/self-hosted/computer-use.md) for the macOS permission steps and the full Linux setup.
>
> ## Next steps
>
> - [Choose where Cloud Agents run](https://cursor.com/docs/cloud-agent/self-hosted/choose-runtime.md): compare managed Cloud Agents, My Machines, and Team Pools.
>
> - [My Machines](https://cursor.com/docs/cloud-agent/self-hosted/my-machines.md): connect your first worker in a few minutes, then configure personal workers, workspace roots, and local MCP servers.
>
> - [Team Pools](https://cursor.com/docs/cloud-agent/self-hosted/pool.md): organize workers into team pools, and [scale worker capacity](https://cursor.com/docs/cloud-agent/self-hosted/pool.md#worker-controller) with a controller.
>
> - [Kubernetes](https://cursor.com/docs/cloud-agent/self-hosted/kubernetes.md): deploy pool workers with the operator and Helm chart.
>
> - [Integrations](https://cursor.com/docs/cloud-agent/self-hosted/integrations.md): partner guides for AWS Lambda, Cloudflare, Namespace, Modal, Daytona, E2B, Vercel, and Coder, and reference templates for AWS Lambda MicroVMs, Cloudflare Containers, and Kubernetes.
>
> - [Computer use](https://cursor.com/docs/cloud-agent/self-hosted/computer-use.md): let agents drive a desktop and browser on your workers.
>
> - [API reference](https://cursor.com/docs/cloud-agent/api/endpoints.md#workers-and-pools): endpoints for workers, pools, the pending-request queue (list, SSE watch, claim, and release), and worker tokens.
>
> - [Self-Hosted Machines](https://cursor.com/help/ai-features/self-hosted-machines): short answers for setup, Team Pools, integrations, and troubleshooting.
>
> ### Bring Self-Hosted Machines to your enterprise
>
> Team Pools require an Enterprise plan. Talk to sales about worker fleets, private connectivity, and rollout. [Contact sales](https://cursor.com/contact-sales?source=docs-cloud-agent-self-hosted)
>
> - 
> - 
> - 
> - 
> - 
> - 
> [Skip to main content](#main-content)[](/docs)Sign inDownload
