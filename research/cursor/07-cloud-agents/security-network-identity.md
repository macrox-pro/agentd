---
primary_sources:
  - id: T1-CLOUD-SEC
    title: "Security"
    url: "https://cursor.com/docs/cloud-agent/security.md"
    section: "Security"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Cloud security and identity

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Security

> # Security overview
>
> This page explains how Cloud Agents are built and secured. It walks through what happens when an agent runs, how access is granted, where the code and data live, how they're isolated and encrypted, and what controls you have over each stage. It answers the questions that come up when a team evaluates Cloud Agents against its security requirements.
>
> For the configuration reference, including secret types, network access modes, and egress IP ranges, see [Secrets & Network](https://cursor.com/docs/cloud-agent/security-network.md). For federating a VM into AWS, GCP, Azure, or a custom verifier without long-lived keys, see [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md). This page explains the model behind those controls; those pages tell you how to set them.
>
> This is a companion to Cursor's broader security material. See the [Trust Center](https://trust.cursor.com/) for certifications, subprocessors, and architecture, and the [Security and Privacy Hardening](https://cursor.com/docs/enterprise/security-hardening.md) reference for the controls you own across all of Cursor. Cursor is SOC 2 Type 2 compliant and commits to at-least-annual penetration testing by reputable third parties.
>
> ## How Cloud Agents work
>
> A Cloud Agent is a coding agent that runs in a virtual machine in Cursor's cloud instead of on a developer's laptop. The VM holds a full development environment: the cloned repository, installed dependencies, configured secrets, and network access.
>
> A single run moves through these stages:
>
> 1. **Start.** A user or integration starts a task from the web app, IDE, CLI, API, Slack, or a linked issue or pull request.
> 2. **Provision.** Cursor provisions an isolated VM for that agent and clones the authorized repository into it.
> 3. **Run.** The agent runs code and tools inside the VM and streams its progress, output, and artifacts back to the user.
> 4. **Persist.** Conversation state, metadata, and artifacts are saved to Cursor-managed storage so you can review and resume the run.
> 5. **Hand off.** The agent pushes its branch and opens a draft pull request for a human to review before anything merges.
> 6. **Recycle.** VM runtime resources are hibernated and then deleted on lifecycle timers once the run is idle.
>
> ## Access and authorization
>
> Cloud Agents reach your code through the Cursor GitHub or GitLab App, not through any single person's credentials.
>
> - **Admins install the app.** Enabling Cloud Agents takes admin privileges on both Cursor and your Git provider. An admin installs the Cursor app on your Git organization and grants access only to the repositories you choose.
> - **Users connect their own account.** Once the app is installed, each user who wants to start an agent connects their own Git account. This is a second, per-user layer on top of the org-level install.
> - **Access is inherited, never widened.** A Cloud Agent can only reach repositories the triggering user could already reach. Starting an agent never grants access to a repository the user didn't already have.
>
> Team admins can go further and lock a Git organization to your Cursor organization with [Protected Git Scopes](https://cursor.com/docs/enterprise/model-and-integration-management.md#protected-git-scopes), so only your teams can start Cloud Agents on its repositories. You can also keep sensitive repositories out entirely with a [repository blocklist](https://cursor.com/docs/enterprise/model-and-integration-management.md#git-repository-blocklist).
>
> Cursor employees do not have access to the code inside Cloud Agent VMs. Access attempts are monitored by Cursor's security team.
>
> ## Isolation and infrastructure
>
> Each agent runs in its own VM boundary, not a shared process sandbox. One agent cannot see another agent's code, environment, or state.
>
> - **Per-agent VMs.** Every agent gets a dedicated environment, isolated from other agents and other users.
> - **MicroVM isolation.** Runtime workspaces run on Firecracker-based microVM infrastructure.
> - **Account-level separation.** Cloud Agent VMs run in a separate AWS account from the rest of Cursor's production infrastructure, so the code-execution environment is walled off from Cursor's other services.
>
> ## Encryption
>
> Cursor encrypts Cloud Agent data in transit and at rest.
>
> - **In transit.** TLS 1.2 or higher for service-to-service and client-to-service traffic.
> - **At rest.** AES-256, with per-agent keys so each agent's session data is encrypted under its own key.
> - **Customer-managed keys.** Enterprise teams can map a customer-managed KMS key (CMEK/BYOK) for Cloud Agent server-side encryption, so you control key rotation and access. See [Data encryption](https://cursor.com/docs/enterprise/privacy-and-data-governance.md#data-encryption).
>
> ## What data is stored, where, and for how long
>
> A Cloud Agent touches four kinds of data. Each is stored in a different place and follows its own retention rule.
>
> | Data                   | What it holds                                                                                          | Where it lives                                            | Retention                                                                                           |
> | :--------------------- | :----------------------------------------------------------------------------------------------------- | :-------------------------------------------------------- | :-------------------------------------------------------------------------------------------------- |
> | **Runtime workspace**  | The checked-out repository, build artifacts, and tool-execution context for a live run                 | The isolated Cloud Agent VM                               | Recycled automatically after the run goes idle; the timer refreshes when you send follow-up prompts |
> | **VM snapshots**       | Point-in-time copies of the VM disk (including cloned code) used to start and resume without recloning | Snapshot and cache layer outside the active VM, encrypted | Rolling 90 days of inactivity; each start or resume extends it, then automatic deletion             |
> | **Conversation state** | Prompts, model responses, tool calls, diff context, and demo artifacts that make up the transcript     | Cursor backend, encrypted with per-agent keys             | Kept indefinitely by default so you can revisit and resume runs; deletable on demand                |
> | **Secrets and tokens** | Cloud Agent secrets, OAuth tokens, and API credentials you configure                                   | Encrypted credential stores in Cursor's backend           | Kept until you delete or remove them                                                                |
>
> The [Delete Agent API](https://cursor.com/docs/cloud-agent/api/endpoints.md#delete-an-agent-permanently) removes an agent's conversation transcript and artifacts on demand. Snapshots can't be deleted on demand; they follow the 90-day inactivity window above. Enterprise teams can also cap conversation retention with [retention policies](https://cursor.com/docs/cloud-agent/security-network.md#cloud-agent-retention-policies). For the full detail on retention and deletion, see [Data retention](https://cursor.com/docs/cloud-agent/security-network.md#data-retention).
>
> ## Privacy and model data
>
> Cloud Agents run in [Privacy Mode](https://cursor.com/privacy-overview). With Privacy Mode on, Cursor never trains on code accessed by Cloud Agents or on the prompts and responses their runs generate. Most models also run under Cursor's zero-data-retention agreements, so providers don't store or train on requests and responses. See [Privacy and Data Governance](https://cursor.com/docs/enterprise/privacy-and-data-governance.md) for the model-by-model detail and the exceptions.
>
> Legacy Privacy Mode is not supported for Cloud Agents, because agents need to store code and environment data in the cloud while they run. Enforce standard Privacy Mode org-wide so every run inherits its zero-data-retention guarantees.
>
> ## Autonomy and prompt injection
>
> Cloud Agents auto-run terminal commands so they can iterate on tests without stopping for approval on every step. This is more autonomous than the foreground agent, and it changes the risk model: an attacker who plants instructions in content the agent reads (a prompt-injection attack) could try to make the agent exfiltrate code to an external host. See [OpenAI's explanation of prompt-injection risk for cloud agents](https://platform.openai.com/docs/codex/agent-network#risks-of-agent-internet-access).
>
> The layers that contain this risk:
>
> - **Network egress controls.** Restrict outbound traffic to a default set plus your allowlist, or to your allowlist only, so a compromised agent has nowhere to send data. Enterprise admins can lock the policy org-wide. See [Network access](https://cursor.com/docs/cloud-agent/security-network.md#network-access).
> - **Redacted runtime secrets.** Mark secrets as [Runtime Secrets](https://cursor.com/docs/cloud-agent/security-network.md#runtime-secrets) so their values are stripped from the transcript, tool output, and commits and never reach the model.
> - **File exclusion.** Add sensitive paths to [`.cursorignore`](https://cursor.com/docs/reference/ignore-file.md) to keep them out of the agent's context.
> - **Human-in-the-loop handoff.** Agents open draft pull requests. Nothing merges until a person reviews the change.
> - **Signed commits.** Every agent commit is signed with an HSM-backed Ed25519 key and shows a "Verified" badge, so agent-authored changes are attributable and can satisfy signed-commit branch protection. See [Signed commits](https://cursor.com/docs/cloud-agent/security-network.md#signed-commits).
>
> For deeper defense, pair these with [hooks](https://cursor.com/docs/hooks.md) to enforce policy and log activity at agent lifecycle points, and have [Bugbot](https://cursor.com/docs/bugbot.md) or [Security Agents](https://cursor.com/docs/security-agents.md) review agent output before it ships.
>
> ## Risk considerations
>
> | Risk                                | Mitigation                                                                                                                                                                                                                    |
> | :---------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | **Full codebase in the cloud**      | Isolated per-agent VMs, AES-256 encryption, and automatic VM and snapshot deletion on lifecycle timers.                                                                                                                       |
> | **Third-party or insider access**   | No Cursor employee access to code in agent VMs, with monitored access attempts. VMs run in a separate AWS account from other Cursor services.                                                                                 |
> | **Agent autonomy**                  | Scope bounded by the repository and the triggering user's access. External reach is limited to configured tools and terminal commands, gated by network egress controls and reviewed through draft PRs.                       |
> | **Network access and exfiltration** | Internet access is on by default but can be restricted to allowlisted domains, down to allowlist-only, and locked org-wide.                                                                                                   |
> | **Secret exposure**                 | Encrypted secret storage, redacted runtime secrets kept out of the model, build-only secrets scoped to the Docker build, and [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md) for short-lived cloud federation. |
>
> ## Auditability
>
> Cloud Agent activity is logged and attributable.
>
> - **Session logging.** Runs are logged, and team admins can review activity from the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents).
> - **Attributed changes.** Every commit and pull request an agent creates is attributed and visible in your Git history, with signed, verified commits.
> - **Audit logs.** Authentication and admin events flow to your [audit logs](https://cursor.com/docs/enterprise/compliance-and-monitoring.md#audit-logs), which Enterprise teams can stream to a SIEM, webhook, or S3.
> - **Run diagnostics.** The built-in [Cursor Cloud MCP](https://cursor.com/docs/cloud-agent/capabilities.md#cursor-cloud-mcp) exposes transcripts, run events, environment details, and setup logs for a run.
>
> ## Data deletion
>
> | Mechanism                         | What it removes                                  | How                                                                                                          |
> | :-------------------------------- | :----------------------------------------------- | :----------------------------------------------------------------------------------------------------------- |
> | **Archive**                       | Hides an agent from the dashboard                | Archive from the [dashboard](https://cursor.com/dashboard/cloud-agents)                                      |
> | **Delete Agent API**              | An agent's conversation transcript and artifacts | [Delete Agent API](https://cursor.com/docs/cloud-agent/api/endpoints.md#delete-an-agent-permanently)         |
> | **Snapshot expiry**               | VM snapshots and cached code                     | Automatic after 90 days of inactivity                                                                        |
> | **Retention policy (Enterprise)** | Conversations older than your chosen window      | [Retention policies](https://cursor.com/docs/cloud-agent/security-network.md#cloud-agent-retention-policies) |
> | **Account deletion**              | The account and its associated data              | [Delete account](https://cursor.com/help/account-and-billing/delete-account.md)                              |
>
> ## FAQ
>
> ### Are Cloud Agents less secure than local agents?
>
> They carry a different risk profile, not a worse one. Running an agent unattended in an isolated sandbox with egress restrictions and minimal permissions can be tighter than a developer's laptop, which usually has full internet access and elevated privileges.
>
> ### Does Cursor store an entire repository permanently?
>
> No. Cursor clones the repository to run the agent, and that clone can live in VM snapshots to speed up future starts, but it isn't kept indefinitely. Snapshots are deleted after 90 days of inactivity.
>
> ### Can a Cloud Agent access repositories the developer can't?
>
> No. Access is gated by the triggering developer's Git access. A Cloud Agent can't reach a repository the developer didn't already have access to.
>
> ### Can we restrict which repositories Cloud Agents reach?
>
> Yes. Cloud Agents can only reach repositories you authorize through your Git provider connection. You control which repositories are available, and admins can lock scopes with [Protected Git Scopes](https://cursor.com/docs/enterprise/model-and-integration-management.md#protected-git-scopes) or exclude repositories with a [blocklist](https://cursor.com/docs/enterprise/model-and-integration-management.md#git-repository-blocklist).
>
> ### Can a Cloud Agent read secrets and credentials?
>
> Configure secrets through the Secrets tab in your dashboard. They're encrypted at rest with KMS, encrypted in transit, and injected as environment variables at runtime. Mark sensitive values as [Runtime Secrets](https://cursor.com/docs/cloud-agent/security-network.md#runtime-secrets) to keep them out of the transcript, tool output, and commits. As a matter of practice, keep secrets out of the repository; if sensitive files must live there, add them to [`.cursorignore`](https://cursor.com/docs/reference/ignore-file.md). For cloud roles, mint [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md) from the VM instead of storing long-lived access keys.
>
> ### Can we audit Cloud Agent activity?
>
> Yes. Sessions are logged, admins can review activity from the dashboard, and every commit and pull request an agent creates is attributed in your Git history. Enterprise teams can stream audit logs to a SIEM.
>
> ### How is data deleted?
>
> Archive an agent from the dashboard, or use the [Delete Agent API](https://cursor.com/docs/cloud-agent/api/endpoints.md#delete-an-agent-permanently) to remove its transcript and artifacts. Full account deletion and Enterprise retention policies remove data on a broader schedule.
>
> ## Related pages
>
> - [Secrets & Network](https://cursor.com/docs/cloud-agent/security-network.md) for secret types, network access modes, egress IP ranges, and signed commits.
> - [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md) for short-lived JWTs and cloud federation.
> - [Privacy and Data Governance](https://cursor.com/docs/enterprise/privacy-and-data-governance.md) for data flows, Privacy Mode, and encryption.
> - [Security and Privacy Hardening](https://cursor.com/docs/enterprise/security-hardening.md) for the controls you configure across Cursor.
> - [Trust Center](https://trust.cursor.com/) for certifications, subprocessors, and architecture.
>
> This summary aids understanding and does not override the [MSA](https://cursor.com/terms/msa), [DPA](https://cursor.com/terms/dpa), or other binding contractual terms.
>
>

### Source: Security and network

> # Secrets & Network
>
> Cloud Agents are available in Privacy Mode. We never train on your code and only retain code for running the agent. [Learn more about Privacy mode](https://www.cursor.com/privacy-overview).
>
> For a walkthrough of how Cloud Agents are architected and secured, including the run lifecycle, access model, isolation, encryption, and data handling, see the [Security overview](https://cursor.com/docs/cloud-agent/security.md). This page is the configuration reference for the controls it describes.
>
> **Privacy Mode (Legacy)** is not supported. Legacy privacy mode blocks cloud
> data storage, and Cloud Agents need to store code and environment data in the
> cloud while they run. Switch to Privacy Mode from [Dashboard → Cloud
> Agents](https://cursor.com/dashboard/cloud-agents) before using Cloud Agents.
>
> ## Secret protection
>
> Secrets provided to Cloud Agents are encrypted at rest and in transit. They are not visible to anyone other than the Cloud Agent user.
>
> Secrets can be set as Environment Variables, Runtime Secrets, or Build Secrets.
>
> ### Environment Variables
>
> Secrets set with type `Environment Variable` are visible to the cloud agent. These are best used for non-sensitive configuration that is helpful for the agent to view, such as flags or public URLs. They are still encrypted at rest and in transit as with other secret types.
>
> ### Runtime secrets
>
> Previously, Runtime Secrets were called Redacted Secrets.
>
> Secrets set with type `Runtime Secret` are still loaded as environment variables, but their contents are redacted from the agent's tool call results, chat transcript, commits, and commit messages, and replaced with the placeholder string `[REDACTED]`. These are best used for sensitive credentials that should not be exposed to the agent and should never be committed to the repository.
>
> Because Runtime Secrets still function internally as environment variables, while they are not shown to the agent, they are still visible to users interacting with the agent's environment via the Terminal.
>
> ### Build secrets
>
> Secrets set with type `Build Secret` are only available to the [Docker build process](https://cursor.com/docs/cloud-agent/security-network.md#manual-setup-with-dockerfile-advanced) (if you have configured one) and are not exposed to the running agent's environment. These are best used for private package registries or build-time credentials that should not be exposed to the agent.
>
> In order to securely use a Build Secret within your Dockerfile, reference them from a `RUN` step using a [Docker secret mount](https://docs.docker.com/build/building/secrets/#secret-mounts), for example:
>
> ```docker
> RUN --mount=type=secret,id=MY_TOKEN,env=MY_TOKEN,required=true \
>     ./scripts/install-private-deps.sh
> ```
>
> ## OIDC identity tokens
>
> For cloud roles and internal APIs, prefer short-lived [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md) over long-lived keys in Secrets. A Cloud Agent can mint a Cursor-signed JWT from a local socket and present it to AWS, GCP, Azure, Vault, or any OIDC verifier.
>
> ## Signed commits
>
> Cloud Agents sign every commit with a HSM-backed Ed25519 key. On GitHub and GitLab, these commits display a "Verified" badge so your team can confirm the commit came from Cursor.
>
> This works automatically for all Cloud Agents. No setup is required.
>
> If your repository enforces branch protection rules that require signed commits, Cloud Agent PRs satisfy those rules without extra configuration.
>
> ## Protected Git Scopes
>
> Team admins can lock a Git organization to your Cursor organization so only your teams can start Cloud Agents on its repositories. See [Protected Git Scopes](https://cursor.com/docs/enterprise/model-and-integration-management.md#protected-git-scopes).
>
> ## What you should know
>
> 1. Grant read-write privileges to our GitHub app for repos you want to edit. We use this to clone the repo and make changes.
> 2. Your code runs inside our AWS infrastructure in isolated VMs and is stored on VM disks while the agent is accessible.
> 3. The agent has internet access by default. You can configure [network egress controls](https://cursor.com/docs/cloud-agent/security-network.md#network-access) for users, teams, and saved environments to restrict the domains the agent can access.
> 4. The agent auto-runs all terminal commands, letting it iterate on tests. This differs from the foreground agent, which requires user approval for every command. Auto-running introduces data exfiltration risk: attackers could execute prompt injection attacks, tricking the agent to upload code to malicious websites. See [OpenAI's explanation about risks of prompt injection for cloud agents](https://platform.openai.com/docs/codex/agent-network#risks-of-agent-internet-access).
> 5. If privacy mode is disabled, we collect prompts and dev environments to improve the product.
> 6. If you disable privacy mode when starting a cloud agent, then enable it during the agent's run, the agent continues with privacy mode disabled until it completes.
>
> ## Data retention
>
> Cloud Agents store two types of data for every run:
>
> - **Conversation history.** The prompts, model responses, tool calls, and demo artifacts that make up the agent's transcript. This is the data you see when you open an agent on the web or from a desktop client.
> - **Environment snapshots.** Encrypted point-in-time copies of the virtual machine disk. Snapshots let you customize VM environments and allow agents to start or resume without recloning the repository or running the setup again.
>
> Conversation history is kept indefinitely by default so you can revisit and resume past runs. Environment snapshots are stored for a maximum of **90 days** of inactivity. Each time an agent starts or resumes from a snapshot, its expiry extends for another 90 days. Once a snapshot goes unused for 90 days, it's deleted automatically, regardless of plan or policy.
>
> You can use the [Delete Agent API](https://cursor.com/docs/cloud-agent/api/endpoints.md#delete-an-agent-permanently) to explicitly delete a cloud agent's conversation history. This endpoint removes the conversation transcript and its artifacts. It doesn't delete environment snapshots, which can't be deleted on demand and instead follow the retention window above.
>
> ### Cloud agent retention policies
>
> Custom retention windows are in early access for select Enterprise teams. [Contact sales](https://cursor.com/contact-sales?source=docs-cloud-agent-retention) to request access.
>
> Enterprise team admins can cap how long the team's Cloud Agent data is kept from **Team Settings** on the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents). The available windows are **Indefinite** and **90 days**.
>
> When you set the policy to **90 days**:
>
> - A background job deletes conversations older than the retention policy window.
> - Environment snapshots continue to follow the rolling 90-day inactivity window described above.
> - The policy applies to every agent run the team owns, including runs from saved environments and the [API](https://cursor.com/docs/cloud-agent/api/v0.md).
>
> Switching back to **Indefinite** stops further conversation deletions but doesn't restore data that's already been removed.
>
> ## Network access
>
> Control which network resources your Cloud Agents can reach. These settings are available on the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents) for individual users, saved environments, and team admins.
>
> ### Private network access
>
> Cloud Agents do not need to run on your hardware to reach private resources. For services in a VPC or intranet, use Tailscale userspace networking, Cloudflare Tunnel, or a similar private-network client in the Cloud Agent environment. See [Running Tailscale](https://cursor.com/docs/cloud-agent/setup.md#running-tailscale) and [Running Cloudflare Tunnel](https://cursor.com/docs/cloud-agent/setup.md#running-cloudflare-tunnel) for setup notes.
>
> With either Tailscale or Cloudflare Tunnel, your private services do not need to accept inbound traffic from the public internet. The agent connects through an authenticated network path, while the service stays on your private network.
>
> Cloudflare Tunnel is a good fit when the agent can reach the private service through an authenticated HTTPS hostname. A connector in your network dials out to Cloudflare, and the Cloud Agent calls that hostname like any other external URL. You can protect the hostname with Cloudflare Access service tokens, store the token values as Cursor Secrets, and add the hostname to your Cloud Agent allowlist.
>
> For TCP targets such as private databases, use a tunnel client that exposes a local TCP listener in the agent environment. The agent then connects to `localhost`, while the tunnel forwards traffic to the private origin.
>
> For private GitHub Enterprise Server, GitLab Enterprise, source control APIs, package registries such as Artifactory or Nexus, and related webhook traffic, Enterprise teams can use [private connectivity](https://cursor.com/docs/cloud-agent/private-connectivity.md) with AWS PrivateLink or Cloudflare Tunnel.
>
> ### Access modes
>
> Three modes control outbound network access for Cloud Agents:
>
> | Mode                         | Behavior                                                                                                                                                     |
> | :--------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | **Allow all network access** | Cloud Agents can reach any external host. No domain restrictions apply.                                                                                      |
> | **Default + allowlist**      | Cloud Agents can reach the [default domains](https://cursor.com/docs/agent/security/run-modes.md#network-access) plus any domains you add to your allowlist. |
> | **Allowlist only**           | Cloud Agents can only reach the domains you explicitly add to your allowlist.                                                                                |
>
> Even in **Allowlist only** mode, a small set of domains remain accessible so Cloud Agents can function. These include Cursor's own services and source control management (SCM) providers.
>
> ### Artifact uploads
>
> Cloud Agents upload [artifacts](https://cursor.com/docs/cloud-agent/capabilities.md#demos-and-artifacts) (screenshots, videos, and log references shown on PRs) to `cloud-agent-artifacts.s3.us-east-1.amazonaws.com`.
>
> If you use **Default + allowlist** or **Allowlist only**, add the exact host to your allowlist so artifact uploads succeed. Don't broaden the entry to `*.s3.us-east-1.amazonaws.com`: the wildcard opens egress to every bucket in the region and creates an exfiltration path for a prompt-injected agent. Blocking the host disables uploads; agent sessions and other tool calls keep working.
>
> ### User-level settings
>
> Individual users can configure their network access mode from the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents) under the **Security** header. Your user-level setting applies to all Cloud Agents you create.
>
> When you select a mode that includes an allowlist (**Default + allowlist** or **Allowlist only**), an allowlist configuration section appears below the setting where you can add your custom domains.
>
> ### Environment-level settings
>
> Saved environments can have their own network access mode and allowlist. Use environment-level settings when one repo or repo group needs stricter egress than the rest of your team.
>
> For example, you can keep a production-adjacent environment on **Allowlist only** while leaving a less sensitive environment on **Default + allowlist**. Agents that use the stricter environment inherit those restrictions.
>
> Environment-level settings include two inheritance options:
>
> | Mode                                         | Behavior                                                                                  |
> | :------------------------------------------- | :---------------------------------------------------------------------------------------- |
> | **Inherit settings**                         | Uses the applicable user or team network access setting.                                  |
> | **Inherit settings + environment allowlist** | Uses the applicable user or team setting and adds domains from the environment allowlist. |
>
> You can also set an environment directly to **Allow all network access**, **Default + allowlist**, or **Allowlist only**.
>
> ### Team-level settings
>
> Team admins can set a default network access mode for the entire team from the same dashboard. The team-level allowlist is the same allowlist that admins configure for the [sandbox default network allowlist](https://cursor.com/docs/agent/security/run-modes.md#network-access). There is no separate allowlist to manage; one allowlist controls both Cloud Agent network access and the sandbox defaults.
>
> When a team-level setting exists:
>
> - If an environment defines its own mode, the **environment setting applies** to agents that use that environment.
> - If an environment inherits settings and a user has configured their own setting, the **user setting takes precedence**.
> - If neither the environment nor the user has configured a setting, the **team default applies**.
>
> ### Locking the setting (Enterprise)
>
> Locking is available for Enterprise teams only.
>
> Enterprise team admins can lock the network access setting using the **Lock Network Access Policy** option. When locked:
>
> - The team-level setting applies to every member, regardless of their individual preference.
> - Users cannot override the locked setting from their own dashboard.
>
> This gives admins full control over Cloud Agent network access across the organization.
>
> ### Relationship to sandbox network policy
>
> The "Default" domains in the **Default + allowlist** mode are the same [default network allowlist](https://cursor.com/docs/agent/security/run-modes.md#network-access) used by the desktop Agent's sandbox. The team-level allowlist is also shared: when an admin configures an allowlist on the dashboard, it applies to both Cloud Agent network access and the [sandbox network policy](https://cursor.com/docs/reference/sandbox.md).
>
> ## Egress IP ranges
>
> Cloud Agents make network connections from specific IP address ranges when accessing external services, APIs, or repositories.
>
> ### API endpoint
>
> The IP ranges are available via a [JSON API endpoint](https://cursor.com/docs/ips.json):
>
> ```bash
> curl https://cursor.com/docs/ips.json
> ```
>
> #### Response format
>
> ```json
> {
>   "version": 1,
>   "modified": "2025-09-24T16:00:00.000Z",
>   "cloudAgents": {
>     "us3p": ["100.26.13.169/32", "34.195.201.10/32", "..."],
>     "us4p": ["54.184.235.255/32", "35.167.37.158/32", "..."],
>     "us5p": ["3.12.82.200/32", "52.14.104.140/32", "..."]
>   },
>   "gitEgressProxy": ["184.73.225.134/32", "3.209.66.12/32", "52.44.113.131/32"]
> }
> ```
>
> - **version**: Schema version number for the API response
> - **modified**: ISO 8601 timestamp of when the IP ranges were last updated
> - **cloudAgents**: Object containing IP ranges, keyed by cluster
> - **gitEgressProxy**: IP addresses used by the [git egress proxy](https://cursor.com/docs/cloud-agent/security-network.md#git-egress-proxy-and-ip-allow-list)
>
> IP ranges published in [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing). You can use an online conversion tool to convert from CIDR notation to IP address ranges if needed.
>
> ### Using the IP ranges
>
> These published IP ranges may be used by Cloud Agents to:
>
> - Clone and push to remote repositories (unless using the [git egress proxy](https://cursor.com/docs/integrations/github.md#ip-allow-list-configuration))
> - Download packages and dependencies
> - Make API calls to external services
> - Access web resources during agent execution
>
> If your organization uses firewall rules or IP allowlists to control network access, you may need to allowlist these IP ranges to ensure Cloud Agents can properly access your services.
>
> **Important considerations:**
>
> - We make changes to our IP addresses from time to time for scaling and operational needs.
> - We do not recommend allowlisting by IP address as your primary security mechanism.
> - If you must use these IP ranges, we strongly encourage regular monitoring of the JSON API endpoint.
>
> ### Git egress proxy and IP allow list
>
> Cursor supports a similar but distinct feature to [use a git egress proxy for IP allow lists](https://cursor.com/docs/integrations/github.md#ip-allow-list-configuration). This proxy routes all git traffic through a narrower set of IPs and works across all git hosts, including GitHub, GitLab, Azure DevOps, and Bitbucket.
>
> For git hosts specifically, we recommend the IP allow list configuration described in the link above, as it integrates directly with the Cursor GitHub app.
>
> If you need to add the proxy IPs directly to an allowlist, use these addresses:
>
> ```text
> 184.73.225.134
> 3.209.66.12
> 52.44.113.131
> ```
>
> ### Cursor Review IPs
>
> If your team uses Cloud Agents alongside [Cursor Review](https://cursor.com/docs/cursor-review/overview.md), allowlist these additional IPs on top of the git egress proxy IPs above:
>
> ```text
> 34.192.39.182
> 50.16.106.255
> 44.217.29.124
> 3.223.245.201
> 54.164.185.10
> 34.194.133.23
> 35.170.116.221
> ```
>
> These IP addresses are stable. If the list ever changes, teams using IP allow
> lists will get advance notice before any address is added or removed.
>
> Enterprise customers with private GitHub Enterprise Server or GitLab Enterprise deployments can use [private connectivity options](https://cursor.com/docs/cloud-agent/private-connectivity.md), so Cloud Agents and Bugbot can access private source control systems.
>
>

### Source: Identity

> # OIDC tokens
>
> Cloud Agents can mint short-lived [OIDC](https://openid.net/specs/openid-connect-core-1_0.html) JWTs from inside the VM and use them to assume cloud roles or call internal services without storing long-lived credentials in [Secrets](https://cursor.com/docs/cloud-agent/security-network.md#secret-protection).
>
> Agents call this API with their terminal tools. You don't need to run these requests yourself.
>
> To have an agent mint tokens, include this in your prompt:
>
> ```text
> To mint OIDC tokens, follow the instructions at
> https://cursor.com/docs/cloud-agent/identity
> ```
>
> This API is local to the agent VM. It is unrelated to the [Cloud Agents API](https://cursor.com/docs/cloud-agent/api/endpoints.md), which uses Cursor API keys and manages agents from outside the VM. The same socket also serves [agent metadata](https://cursor.com/docs/cloud-agent/metadata.md) for values that don't belong in a credential.
>
> Cursor-managed Cloud Agent VMs serve the token socket. Every token they mint carries `agent_runtime: managed`.
>
> ## How it works
>
> 1. The agent calls the local socket and asks for a token with an audience the verifier expects.
> 2. Cursor signs an RS256 JWT bound to that agent and owner.
> 3. The agent sends the JWT to your cloud or verifier (AWS STS, GCP, Azure, Vault, or a service you run).
> 4. The verifier checks the signature against Cursor's published JWKS and authorizes on claims such as `sub`, `team_id`, or `cloud_agent_id`.
>
> ## Mint a token
>
> The agent mints a token over the Unix socket at `CURSOR_AGENT_SOCKET`. On Cursor-managed VMs the default is `/run/cursor/api.sock`.
>
> ```bash
> curl --unix-socket "${CURSOR_AGENT_SOCKET:-/run/cursor/api.sock}" \
>   -H 'Content-Type: application/json' \
>   -d '{"aud":"sts.amazonaws.com"}' \
>   http://cursor-agent/v1/tokens/oidc
> ```
>
> Requests are HTTP over a Unix socket. The hostname in the URL is ignored.
>
> Include an optional `nonce` when the verifier expects replay binding:
>
> ```bash
> curl --unix-socket "${CURSOR_AGENT_SOCKET:-/run/cursor/api.sock}" \
>   -H 'Content-Type: application/json' \
>   -d '{"aud":"https://oidc.example.com","nonce":"unpredictable-value"}' \
>   http://cursor-agent/v1/tokens/oidc
> ```
>
> ### Request
>
> `POST /v1/tokens/oidc` over the Unix socket. `Content-Type: application/json` is required. Maximum body size is 4 KB.
>
> | Field       | Required | Description                                                                                                                                                                                                                                                                                                                                                                                        |
> | :---------- | :------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `aud`       | Yes      | Audience string your verifier checks. Printable ASCII, no whitespace, up to 512 characters. Examples: `sts.amazonaws.com`, `https://oidc.example.com`.                                                                                                                                                                                                                                             |
> | `nonce`     | No       | Opaque string echoed into the JWT `nonce` claim. Up to 512 characters.                                                                                                                                                                                                                                                                                                                             |
> | `sub_claim` | No       | Claim name to put in `sub` as `<name>:<value>`, for verifiers that only match `sub` and `aud`. Up to 64 characters. Discovery lists the supported names in `x_cursor_sub_claims_supported`; currently `team_id`. Unsupported names are rejected. If the claim has no value for this agent, such as `team_id` on a personal account, the mint fails instead of falling back to the default subject. |
>
> Cursor doesn't allowlist audiences. Your verifier must reject unexpected `aud` values.
>
> ### Response
>
> ```json
> {
>   "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ii4uLiJ9...",
>   "expires_at": 1785500000
> }
> ```
>
> | Field        | Description                                              |
> | :----------- | :------------------------------------------------------- |
> | `token`      | Signed JWT.                                              |
> | `expires_at` | Expiration as Unix seconds. Matches the JWT `exp` claim. |
>
> Tokens are valid for **5 minutes**. There is no refresh endpoint. Mint again when a new token is needed.
>
> ### When claims appear
>
> [Install scripts](https://cursor.com/docs/cloud-agent/setup.md) can mint on the same socket. A token only includes claims that have a value when it is minted: `turn_id` and `turn_start` are absent until a coding turn starts, and `branch_name` is absent until the run records a branch. Owner, team, and repository claims are set from agent creation onward.
>
> If the socket is missing right after boot, retry the connection.
>
> ## Verify a token
>
> Publish these URLs to your identity provider or resource server:
>
> | Endpoint  | URL                                                       |
> | :-------- | :-------------------------------------------------------- |
> | Issuer    | `https://api.cursor.com`                                  |
> | Discovery | `https://api.cursor.com/.well-known/openid-configuration` |
> | JWKS      | `https://api.cursor.com/keys`                             |
>
> ```bash
> curl -sS https://api.cursor.com/.well-known/openid-configuration
> curl -sS https://api.cursor.com/keys
> ```
>
> Discovery follows [OpenID Connect Discovery 1.0](https://openid.net/specs/openid-connect-discovery-1_0.html). Tokens are minted on the agent VM, so the discovery document has no `authorization_endpoint` or `token_endpoint`.
>
> ### Older issuer URL
>
> Cursor still serves a second discovery document at
> `https://api2.cursor.sh/cloud-agent/identity`. Minted tokens no longer carry
> that issuer. Point verifiers at `https://api.cursor.com`.
>
> Check at least:
>
> - Signature with RS256 and the JWKS `kid`
> - `iss` is `https://api.cursor.com`
> - `aud` is the audience your service expects
> - `nbf` / `exp` with a small clock-skew allowance (`nbf` is 5 seconds before `iat`)
> - `sub` or other claims your policy uses
>
> Discovery includes `x_cursor_audience_bound: true`. Every token is minted for the caller-supplied `aud`. Don't accept a token issued for a different audience. Discovery also publishes `x_cursor_sub_claims_supported`, the claim names a mint request can project into `sub` with `sub_claim`.
>
> ## JWT claims
>
> Header: `alg=RS256`, `typ=JWT`, plus `kid`.
>
> | Claim                      | Always present        | Description                                                                                                                                                                                                                  |
> | :------------------------- | :-------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `iss`                      | Yes                   | `https://api.cursor.com`                                                                                                                                                                                                     |
> | `sub`                      | Yes                   | Stable owner subject: `user:<id>` or `service_account:<id>` by default, or `<claim>:<value>` (for example `team_id:123`) when the mint request set `sub_claim`. Not an email.                                                |
> | `aud`                      | Yes                   | Audience from the mint request.                                                                                                                                                                                              |
> | `iat`                      | Yes                   | Issued-at, Unix seconds.                                                                                                                                                                                                     |
> | `nbf`                      | Yes                   | Not-before (`iat - 5`).                                                                                                                                                                                                      |
> | `exp`                      | Yes                   | Expiration (`iat + 300`).                                                                                                                                                                                                    |
> | `jti`                      | Yes                   | Unique id per mint.                                                                                                                                                                                                          |
> | `cloud_agent_id`           | Yes                   | Cloud Agent id (`bcId`).                                                                                                                                                                                                     |
> | `nonce`                    | No                    | Present only when the mint request included one.                                                                                                                                                                             |
> | `agent_runtime`            | Yes                   | `managed` on Cursor-managed Cloud Agent VMs.                                                                                                                                                                                 |
> | `owner_email`              | When known            | Lowercased user email. Prefer `sub` or `owner_user_id` for allowlists; email can change.                                                                                                                                     |
> | `owner_user_id`            | When known            | Cursor user id, as a decimal string.                                                                                                                                                                                         |
> | `owner_service_account_id` | When known            | Service account id when a service account owns the agent.                                                                                                                                                                    |
> | `team_id`                  | When known            | Owning team id, as a decimal string.                                                                                                                                                                                         |
> | `turn_id`                  | When a turn is active | Id of this coding turn. Different from `cloud_agent_id`, which is the Cloud Agent id (`bcId`).                                                                                                                               |
> | `turn_start`               | When a turn is active | Run start, Unix seconds.                                                                                                                                                                                                     |
> | `repo_url`                 | When known            | Primary repository in `host/path` form, such as `github.com/acme/widgets`. Hostname is lowercased, with no scheme, credentials, port, query, or `.git` suffix. On a multi-repo agent, this is only the primary repository.   |
> | `repo_urls`                | When known            | Every repository in the workspace, same form as `repo_url`. Primary repository first, then the rest sorted. Present only when the set is known complete. Missing means the set isn't known, not that there is only one repo. |
> | `repo_count`               | When known            | Number of entries in `repo_urls`. Present exactly when `repo_urls` is. Use this with `repo_url` when your verifier can only match a single value (`repo_count == 1`).                                                        |
> | `branch_name`              | When known            | Current branch.                                                                                                                                                                                                              |
> | `environment_id`           | When known            | Id of the Cursor environment this run used.                                                                                                                                                                                  |
> | `source`                   | When known            | How the agent was started, such as `WEBSITE`, `API`, `SLACK`, or `AUTOMATIONS`.                                                                                                                                              |
> | `automation_id`            | For automations       | Automation id when `source` is automations.                                                                                                                                                                                  |
>
> `repo_url` is the primary repository. To confine an agent to specific repositories, pin the complete set with `repo_urls`.
>
> ## Trust model
>
> The token identifies the Cloud Agent run, not a specific process inside the VM. Any process that can reach the socket can mint a token: the agent, code it runs, and hooks. Scope permissions to what you would grant that run as a whole.
>
> You don't choose which agent the token is for. Cursor fills claims from this run, so a process in the VM can't mint a token for a different agent.
>
> ## Rate limits and errors
>
> Each agent VM can mint **30 tokens per minute**, in bursts of up to 10. The socket also accepts at most 8 connections at once. That cap is shared with [agent metadata](https://cursor.com/docs/cloud-agent/metadata.md). Cache a token until it expires instead of minting per call.
>
> Retry `429`, `503`, `500`, `502`, and `504` with backoff. Treat `403` as fatal: this agent isn't allowed to mint.
>
> Error bodies carry a machine-readable code. Invalid-request errors (400, 404, 405, 413, and 415) also include a `usage` string that restates the full request contract. Rate-limit and saturation errors stay code-only:
>
> ```json
> { "error": "invalid_aud", "usage": "POST /v1/tokens/oidc ..." }
> ```
>
> ```json
> { "error": "rate_limited" }
> ```
>
> | HTTP      | `error`                                                                | When                                                                                                                                                                     |
> | :-------- | :--------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | 400       | `invalid_json`, `invalid_aud`, `invalid_nonce`, or `invalid_sub_claim` | Bad request body                                                                                                                                                         |
> | 404       | `not_found`                                                            | Wrong path                                                                                                                                                               |
> | 405       | `method_not_allowed`                                                   | Not `POST`                                                                                                                                                               |
> | 413       | `body_too_large`                                                       | Body over 4 KB                                                                                                                                                           |
> | 415       | `invalid_content_type`                                                 | Missing or non-JSON `Content-Type`                                                                                                                                       |
> | 429       | `rate_limited`                                                         | Over the per-agent mint budget; honor `Retry-After`                                                                                                                      |
> | 503       | `saturated`                                                            | Too many connections; honor `Retry-After`                                                                                                                                |
> | 500       | `host_error`                                                           | Internal error; retry                                                                                                                                                    |
> | 502 / 504 | `backend_unreachable`                                                  | Cursor couldn't mint the token; retry                                                                                                                                    |
> | Other     | `backend_error`                                                        | Cursor rejected the mint. `400` means fix the request (for example an unsupported `sub_claim`, or one with no value for this agent). `403` is fatal. `503` is retryable. |
>
> ## AWS IAM example
>
> Use OIDC when you want AWS to trust Cursor-signed JWTs with `AssumeRoleWithWebIdentity`. For the simpler Cursor-managed assume-role flow (External ID + `CURSOR_AWS_ASSUME_IAM_ROLE_ARN`), see [Using AWS IAM Roles](https://cursor.com/docs/cloud-agent/setup.md#using-aws-iam-roles).
>
> 1. Create an IAM OIDC identity provider whose URL is `https://api.cursor.com`.
> 2. Set the audience to `sts.amazonaws.com` (or another audience your role expects).
> 3. Trust the role only for subjects and teams you intend to allow.
>
> Example trust policy:
>
> ```json
> {
>   "Version": "2012-10-17",
>   "Statement": [
>     {
>       "Effect": "Allow",
>       "Principal": {
>         "Federated": "arn:aws:iam::123456789012:oidc-provider/api.cursor.com"
>       },
>       "Action": "sts:AssumeRoleWithWebIdentity",
>       "Condition": {
>         "StringEquals": {
>           "api.cursor.com:aud": "sts.amazonaws.com"
>         },
>         "StringLike": {
>           "api.cursor.com:sub": "user:*"
>         }
>       }
>     }
>   ]
> }
> ```
>
> Tighten this with an exact `sub`, such as `user:42` for one user or `service_account:<id>` for an agent that runs as a service account. AWS trust policies only match `aud` and `sub`, so scope trust to a team by minting with `"sub_claim":"team_id"` and matching the projected subject:
>
> ```json
> "StringEquals": {
>   "api.cursor.com:aud": "sts.amazonaws.com",
>   "api.cursor.com:sub": "team_id:123"
> }
> ```
>
> Follow current [AWS IAM OIDC](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_oidc.html) instructions for provider creation and thumbprints.
>
> The agent mints with `"aud":"sts.amazonaws.com"` (plus `"sub_claim":"team_id"` when the trust policy matches the team subject) and passes the JWT to STS. If you use [network allowlists](https://cursor.com/docs/cloud-agent/security-network.md#network-access), allow `sts.amazonaws.com` (and any regional STS host you call).
>
> ## Other verifiers
>
> The same tokens work with any OIDC-compliant verifier:
>
> - **GCP** Workload Identity Federation
> - **Azure** federated credentials / Entra ID
> - **Vault** JWT/OIDC auth
> - Internal APIs that validate RS256 JWTs
>
> Point the provider at the discovery URL, require your audience, and authorize on claims such as `sub`, `team_id`, or `cloud_agent_id`. To confine an agent to specific repositories, pin the complete set with `repo_urls`; `repo_url` names only the primary repository.
>
> Minting uses the local socket only. Exchanging the JWT with AWS, GCP, Azure, or your service still needs outbound network access to those hosts.
>
> ## Related pages
>
> - [Agent metadata](https://cursor.com/docs/cloud-agent/metadata.md) for key-value run metadata on the same socket
> - [Secrets & Network](https://cursor.com/docs/cloud-agent/security-network.md) for dashboard secrets and egress controls
> - [Cloud agent setup](https://cursor.com/docs/cloud-agent/setup.md#using-aws-iam-roles) for Cursor-managed AWS role assumption
> - [Security overview](https://cursor.com/docs/cloud-agent/security.md) for isolation and access model
> - [Service accounts](https://cursor.com/docs/account/enterprise/service-accounts.md) when agents run as a team service account
>
>

### Source: Private connectivity

> # Private Connectivity
>
> Cursor supports private network connectivity for Enterprise teams that need Cursor to work with systems that are not reachable from the public internet. This includes self-hosted GitHub Enterprise Server, GitLab Enterprise, Bitbucket Data Center, Artifactory, Nexus, private source control APIs, and webhook traffic from those systems back to Cursor.
>
> The same private connectivity setup is used across Cursor services that need access to your source control system, including [Cloud Agents](https://cursor.com/docs/cloud-agent.md), [Bugbot](https://cursor.com/docs/bugbot.md), and Cursor backend services.
>
> To set up private connectivity, contact [hi@cursor.com](mailto:hi@cursor.com) or your Cursor sales representative.
>
> ## Supported options
>
> | Option            | Best for                                                                                                                | Cloud provider                             | Status    |
> | :---------------- | :---------------------------------------------------------------------------------------------------------------------- | :----------------------------------------- | :-------- |
> | AWS PrivateLink   | Private connectivity between Cursor and your Git provider or package registry, including webhook traffic back to Cursor | AWS                                        | Supported |
> | Cloudflare Tunnel | Cursor accessing a private origin when AWS PrivateLink is not practical                                                 | Any environment that can run `cloudflared` | Supported |
>
> ## How to choose
>
> Use AWS PrivateLink when your private Git provider or package registry is in AWS or can sit behind an AWS Network Load Balancer. This is the preferred path for self-hosted GitHub Enterprise Server and GitLab Enterprise.
>
> AWS PrivateLink can cover two traffic directions:
>
> - Cursor accessing your private Git provider to clone repositories and call Git APIs.
> - Your Git provider sending webhooks or callbacks to Cursor over `api2.cursor.sh` without public internet egress.
>
> Use Cloudflare Tunnel when you cannot publish an AWS endpoint service or when you need a deployment model that only requires an outbound tunnel from your network.
>
> If your team requires Google Private Service Connect (PSC), contact Cursor. Cursor does not currently offer a customer-facing PSC service.
>
> ## Prerequisites
>
> Before starting, make sure you have:
>
> - A Cursor Enterprise workspace
> - A self-hosted GitHub Enterprise Server, GitLab Enterprise, Bitbucket Data Center, or private package registry (such as Artifactory or Nexus) reachable over HTTPS on port 443
> - A publicly trusted TLS certificate for the Git or registry hostname
> - DNS ownership for that hostname
> - AWS permissions to create endpoint services or interface VPC endpoints, if using AWS PrivateLink
> - Permission to run `cloudflared`, if using Cloudflare Tunnel
>
> Cursor does not support self-signed certificates, unencrypted connections, SSH, custom ports, or IPv6-only endpoint services for these private connectivity paths.
>
> If you run a proxy in front of GitHub Enterprise Server, make sure it allows Cursor's GitHub App integration to use authenticated GitHub REST and GraphQL APIs.
>
> ## AWS PrivateLink
>
> AWS PrivateLink supports private traffic in either direction between Cursor and your Git provider or package registry. You may need one direction or both, depending on your network policy.
>
> ### Direction 1: Cursor to your Git provider or package registry
>
> Use this option when Cursor needs to clone repositories, call Git APIs, or reach a private package registry such as Artifactory or Nexus.
>
> #### 1. Create an AWS endpoint service
>
> Create a Network Load Balancer in front of your Git provider or package registry HTTPS endpoint. Publish that load balancer as an AWS VPC endpoint service.
>
> Send Cursor:
>
> - Endpoint service name, for example `com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0`
> - AWS region
> - Git or registry hostname, for example `github.example.com` or `artifactory.example.com`
> - Whether your endpoint service has AWS-managed private DNS enabled
> - Whether your Network Load Balancer preserves client IPs or your backend filters source IPs
>
> If your endpoint service is outside `us-east-1`, enable cross-region access on the endpoint service.
>
> #### 2. Allow Cursor's AWS principal
>
> Cursor will provide the AWS principal to add to your endpoint service allowed principals. Add the exact principal Cursor provides:
>
> ```text
> arn:aws:iam::<cursor-aws-account-id>:role/<cursor-provided-role>
> ```
>
> Cursor cannot create its interface endpoint until this principal is allowed. If the principal is missing or does not match exactly, AWS returns `InvalidServiceName`.
>
> If your load balancer preserves client IPs, or if your backend filters source IPs, allow these Cursor PrivateLink subnet CIDRs:
>
> ```text
> 10.2.8.0/21
> 10.2.24.0/21
> 10.2.40.0/21
> ```
>
> #### 3. Accept the endpoint connection
>
> After Cursor creates the interface endpoint, accept the endpoint connection in your AWS account if your endpoint service requires manual acceptance.
>
> #### 4. Configure DNS
>
> If your endpoint service exposes AWS-managed private DNS for your Git or registry hostname, Cursor enables private DNS on its interface endpoint.
>
> If your endpoint service does not expose private DNS, Cursor creates private DNS on its side and maps that hostname to the endpoint DNS name.
>
> Use the same hostname in Cursor that appears on the TLS certificate and in DNS.
>
> ### Direction 2: Your Git provider to `api2.cursor.sh`
>
> Use this option when your GitHub Enterprise Server or GitLab Enterprise host cannot reach the public internet but still needs to send webhooks or callbacks to Cursor.
>
> Cursor publishes an AWS PrivateLink endpoint service for `api2.cursor.sh`. You create an interface VPC endpoint in your AWS account and enable private DNS so `api2.cursor.sh` resolves to private endpoint IPs from your Git provider network.
>
> #### Endpoint service details
>
> Cursor will confirm your AWS principal is allowlisted before you create the endpoint.
>
> | Field                      | Value                                                                                |
> | :------------------------- | :----------------------------------------------------------------------------------- |
> | Service name               | `com.amazonaws.vpce.us-east-1.vpce-svc-054b15427d4bea2b7`                            |
> | Service ID                 | `vpce-svc-054b15427d4bea2b7`                                                         |
> | Home region                | `us-east-1`                                                                          |
> | Supported consumer regions | `us-east-1`, `us-east-2`, `us-west-2`, `eu-central-1`, `eu-west-1`, `ap-southeast-2` |
> | IP address types           | IPv4 only                                                                            |
> | Private DNS name           | `api2.cursor.sh`                                                                     |
>
> #### Mode 1: AWS-managed private DNS
>
> This is the recommended mode. Set `private_dns_enabled = true`.
>
> ```hcl
> resource "aws_vpc_endpoint" "cursor_api2" {
>   vpc_id              = aws_vpc.app.id
>   service_name        = "com.amazonaws.vpce.us-east-1.vpce-svc-054b15427d4bea2b7"
>   service_region      = "us-east-1"
>   vpc_endpoint_type   = "Interface"
>   subnet_ids          = [for subnet in aws_subnet.app_private : subnet.id]
>   private_dns_enabled = true
>   security_group_ids  = [aws_security_group.cursor_api2_endpoint.id]
> }
> ```
>
> AWS associates your VPC with the managed private hosted zone for `api2.cursor.sh`. Inside the VPC, `api2.cursor.sh` resolves to the endpoint ENI IPs. No Route 53 record is required.
>
> #### Mode 2: Customer-managed private hosted zone
>
> Use this mode if you want to own the DNS record. Set `private_dns_enabled = false`, then create a private hosted zone for `api2.cursor.sh` scoped to the consumer VPC.
>
> ```hcl
> resource "aws_vpc_endpoint" "cursor_api2" {
>   vpc_id              = aws_vpc.app.id
>   service_name        = "com.amazonaws.vpce.us-east-1.vpce-svc-054b15427d4bea2b7"
>   service_region      = "us-east-1"
>   vpc_endpoint_type   = "Interface"
>   subnet_ids          = [for subnet in aws_subnet.app_private : subnet.id]
>   private_dns_enabled = false
>   security_group_ids  = [aws_security_group.cursor_api2_endpoint.id]
> }
>
> resource "aws_route53_zone" "cursor_api2" {
>   name    = "api2.cursor.sh"
>   comment = "Customer-managed PHZ for api2.cursor.sh scoped to the app VPC."
>
>   vpc {
>     vpc_id = aws_vpc.app.id
>   }
> }
>
> resource "aws_route53_record" "cursor_api2_a" {
>   zone_id = aws_route53_zone.cursor_api2.zone_id
>   name    = "api2.cursor.sh"
>   type    = "A"
>
>   alias {
>     name                   = aws_vpc_endpoint.cursor_api2.dns_entry[0].dns_name
>     zone_id                = aws_vpc_endpoint.cursor_api2.dns_entry[0].hosted_zone_id
>     evaluate_target_health = false
>   }
> }
> ```
>
> If GitHub Enterprise Server or GitLab Enterprise uses DNS outside the endpoint VPC, forward `api2.cursor.sh` queries to the VPC resolver or create an equivalent private DNS override. Do not create a public DNS override.
>
> ## Cloudflare Tunnel
>
> Use Cloudflare Tunnel when AWS PrivateLink is not a fit.
>
> Cursor creates the tunnel and shares:
>
> - A public hostname under Cursor-controlled DNS
> - A tunnel token through a secure 1Password share
> - A sample `cloudflared` configuration
>
> Your network runs `cloudflared` and opens outbound connections to Cloudflare. No inbound firewall rule is required.
>
> Example `cloudflared` configuration:
>
> ```yaml
> ingress:
>   - hostname: <cursor-provided-hostname>
>     service: https://<your-internal-service>:443
>   - service: http_status:404
> ```
>
> Example run command:
>
> ```bash
> docker run -d --restart=always --name cloudflared \
>   -v /path/to/config.yml:/etc/cloudflared/config.yml \
>   cloudflare/cloudflared:latest \
>   tunnel --config /etc/cloudflared/config.yml \
>   run --token <TUNNEL_TOKEN>
> ```
>
> Keep the tunnel token secret. Do not send it through email or chat.
>
> ## Complete the source control connection
>
> After private networking is configured, complete the source control setup in Cursor:
>
> - For GitHub Enterprise Server, follow the [GitHub integration setup](https://cursor.com/docs/integrations/github.md#setup).
> - For GitLab Enterprise, follow the [GitLab integration setup](https://cursor.com/docs/integrations/gitlab.md#setup).
> - For Bitbucket Data Center, follow the [Bitbucket integration setup](https://cursor.com/docs/integrations/bitbucket.md#setup).
> - Use the same hostname that is covered by your TLS certificate and private DNS configuration.
> - If a proxy sits in front of your Git provider, make sure it allows the authenticated API traffic described in [Prerequisites](https://cursor.com/docs/cloud-agent/private-connectivity.md#prerequisites).
>
> Cursor uses the connected source control integration for Cloud Agents, Bugbot, and other Cursor services that need repository access.
>
> ### Check the private webhook path
>
> If your Git provider sends webhooks to Cursor through the `api2.cursor.sh` PrivateLink path, run these checks from the same network path used by GitHub Enterprise Server or GitLab Enterprise:
>
> ```bash
> getent hosts api2.cursor.sh
> # or, if dig is available
> dig +short api2.cursor.sh
> curl -sS https://api2.cursor.sh/
> ```
>
> Every resolved IP should be inside your consumer VPC CIDR. If you see public IPs such as `3.x.x.x` or `44.x.x.x`, private DNS is not in effect.
>
> The `curl` request should return HTTP `200` with a body that starts with `Welcome to Cursor.` That response means the request reached a live Cursor `api2` backend.
>
> ## Troubleshooting
>
> | Symptom                                                                                            | Likely cause                                                                                        | Fix                                                                                                                                       |
> | :------------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------- |
> | Cursor cannot complete the private connection to your Git provider                                 | Cursor cannot reach or attach to the endpoint service                                               | Confirm the endpoint service name, region, and allowed principal match the values Cursor provided, then contact Cursor with the timestamp |
> | Cursor reports that the endpoint connection is waiting for customer action                         | The endpoint service requires approval in your AWS account                                          | Review pending endpoint connection requests for the service and approve the Cursor request                                                |
> | Bugbot or Cloud Agents connect to GHES but fail during app setup, repo sync, or webhook processing | A proxy in front of GHES is blocking or rewriting authenticated GitHub REST or GraphQL API requests | Allow Cursor's GitHub App integration to use authenticated GitHub REST and GraphQL APIs                                                   |
> | `api2.cursor.sh` resolves to public IPs                                                            | Private DNS is not in the resolver path used by GitHub Enterprise Server or GitLab Enterprise       | Enable AWS-managed private DNS, or forward DNS to the endpoint VPC resolver                                                               |
> | TCP to `api2.cursor.sh:443` times out                                                              | Security group, NACL, route table, or firewall blocks traffic to endpoint ENIs                      | Allow TCP 443 from your Git provider network to the endpoint ENIs                                                                         |
> | TLS fails for `api2.cursor.sh`                                                                     | DNS points to the wrong target or the client is not using SNI                                       | Check endpoint DNS and retry with SNI enabled                                                                                             |
> | `curl https://api2.cursor.sh/` does not return `Welcome to Cursor.`                                | Traffic is not reaching a healthy Cursor backend                                                    | Send Cursor the timestamp, source VPC, and resolved endpoint IPs                                                                          |
> | Cloudflare Tunnel does not connect                                                                 | `cloudflared` cannot reach Cloudflare or the token/config is wrong                                  | Check outbound firewall rules, token, and `cloudflared` logs                                                                              |
>
> ## Google Private Service Connect
>
> Cursor does not currently offer customer-facing Google Private Service Connect.
>
> If you need private connectivity from a GCP VPC to Cursor services, or from Cursor to a private service in your GCP project, contact Cursor so we can scope the requirement. Today, use AWS PrivateLink or Cloudflare Tunnel when those deployment models fit.
>
> ## What to send Cursor
>
> For AWS PrivateLink to your Git provider or package registry:
>
> - Endpoint service name
> - AWS region
> - Git or registry hostname
> - Whether private DNS is enabled
> - Whether your load balancer preserves client IPs or filters source IPs
>
> For `api2.cursor.sh` over AWS PrivateLink:
>
> - AWS principal Cursor should allowlist
> - VPC and region where you will create the interface endpoint
> - Whether you plan to use AWS-managed private DNS or customer-managed DNS
>
> For Cloudflare Tunnel:
>
> - Internal origin URL
> - Customer contacts for the secure 1Password share
> - Any hostname or naming restrictions
>
> ## Further reading
>
> - [AWS: Create an endpoint service](https://docs.aws.amazon.com/vpc/latest/privatelink/create-endpoint-service.html)
> - [AWS: Manage DNS names for VPC endpoint services](https://docs.aws.amazon.com/vpc/latest/privatelink/manage-dns-names.html)
> - [AWS: Access an AWS service using an interface VPC endpoint](https://docs.aws.amazon.com/vpc/latest/privatelink/create-interface-endpoint.html)
> - [Cloudflare Tunnel documentation](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/)
>
>
