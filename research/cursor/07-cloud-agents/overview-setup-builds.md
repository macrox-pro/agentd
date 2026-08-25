---
primary_sources:
  - id: T1-CLOUD
    title: "Cloud"
    url: "https://cursor.com/docs/cloud-agent.md"
    section: "Cloud"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Cloud agents overview

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Cloud agents

> # Cloud Agents
>
> Cloud agents use the same [agent fundamentals](https://cursor.com/learn/agents.md) but run in isolated VMs in the cloud with full development environments instead of on your local machine. The development environment is similar to the setup on your laptop: cloned repos, installed dependencies, secrets, startup commands, and network access.
>
> Effective development environments give agents full context on your codebase and organization, so they can test and verify their work.
>
> ## Why use Cloud Agents?
>
> You can run as many agents as you want in parallel, and they do not require your local machine to be connected to the internet.
>
> Because they have access to their own virtual machine, cloud agents can build, test, and interact with the changed software. They can also use computers to control the desktop and browser. Cloud agents support [MCP servers](https://cursor.com/docs/mcp.md), giving them access to external tools and data sources like databases, APIs, and third-party services.
>
> Cloud agents can also run in multi-repo environments. Use one when a task spans separate frontend, backend, infrastructure, or shared-library repositories. The agent can inspect the full workspace, make coordinated changes, and open pull requests in the repos it changes. Long-running is not available for multi-repo environments yet.
>
> ## How to access
>
> Before anyone can start a cloud agent from a repository, a Cursor account admin needs to connect source control for the account. Set up [GitHub (Cloud and Enterprise Server)](https://cursor.com/docs/integrations/github.md), [GitLab (Cloud and Self-Hosted)](https://cursor.com/docs/integrations/gitlab.md), [Bitbucket Cloud](https://cursor.com/docs/integrations/bitbucket.md), or [Azure DevOps](https://cursor.com/docs/integrations/azure-devops.md).
>
> You can kick off cloud agents from wherever you work:
>
> 1. **Cursor for iOS**: Start and manage agents from the [Cursor iOS app](https://cursor.com/docs/cloud-agent/mobile.md)
> 2. **Cursor Web**: Start and manage agents from [cursor.com/agents](https://cursor.com/agents) on any device
> 3. **Cursor Desktop**: Select **Cloud** in the dropdown under the agent input
> 4. **Slack**: Use the @cursor command to kick off an agent
> 5. **GitHub or Bitbucket**: Comment `@cursor` on a GitHub PR or issue, or on a Bitbucket PR, to kick off an agent
> 6. **Linear**: Use the @cursor command to kick off an agent
> 7. **API**: Use the API to kick off an agent
>
> On **Android**, use [cursor.com/agents](https://cursor.com/agents) in Chrome
> and tap **Install App** for a Progressive Web App (PWA). See [Cursor for
> iOS](https://cursor.com/docs/cloud-agent/mobile.md) for the native iPhone and iPad app and more mobile
> options.
>
> ### Use Cursor in Slack
>
> Learn more about setting up and using the Slack integration, including
> triggering agents and receiving notifications.
>
> ## How it works
>
> ### Repository provider connection
>
> Cloud agents clone your repo from GitHub, GitLab, Azure DevOps Services, or Bitbucket Cloud and work on a separate branch, then push changes to your repo for handoff.
>
> You need read-write privileges to your repo and any dependent repos or submodules.
>
> ### Environments
>
> Agents are only as capable as the environments they run in. An agent that can write code but can't run tests, query services, or reach APIs cannot close the loop on its work.
>
> Not setting up a development environment for your cloud agents is like not giving your engineers a computer. This is why environment setup is the most important step to improve the effectiveness of cloud agents. It lets cloud agents work like engineers do: write code, test and verify work, and ship software.
>
> You can configure environments with agent-led setup, a saved snapshot, or a Dockerfile in `.cursor/environment.json`. See [Cloud agent setup](https://cursor.com/docs/cloud-agent/setup.md) to get started. [Builds](https://cursor.com/docs/cloud-agent/builds.md) prepare each environment in the background so agents start with repositories and dependencies ready.
>
> The Cloud Agents dashboard shows which environment and Build an agent used, along with environment details and version history. On the agent page, hover over the repository name at the top of the page to inspect the environment used for that run. See [Cloud agent setup](https://cursor.com/docs/cloud-agent/setup.md) for configuration details.
>
> ### Runtime and environment controls
>
> Cursor manages VM provisioning, isolation, snapshots, startup, artifacts, and capacity for every Cloud Agent. You can add secrets, restrict outbound domains, connect to private networks with Tailscale or a similar client, and use private connectivity for supported source control paths.
>
> See [Cloud Agent security and network](https://cursor.com/docs/cloud-agent/security-network.md) for the full set of environment and network controls. If you're weighing whether to self-host, see [why most teams start with Cursor Cloud](https://cursor.com/docs/cloud-agent/self-hosted.md).
>
> ## Models
>
> Cloud Agents use a curated selection of models. You can select the context window size for supported models.
>
> ## MCP support
>
> Cloud agents can use [MCP (Model Context Protocol)](https://cursor.com/docs/mcp.md) servers configured for your team. Add and manage MCP servers through the MCP dropdown in [cursor.com/agents](https://cursor.com/agents).
>
> Both HTTP and stdio transports are supported. OAuth is supported for MCP servers that need it. See [Cloud Agent capabilities](https://cursor.com/docs/cloud-agent/capabilities.md) for setup details.
>
> Cloud Agents also include a built-in [Cursor Cloud MCP](https://cursor.com/docs/cloud-agent/capabilities.md#cursor-cloud-mcp) for run diagnostics, including transcripts, run events, environment details, and setup logs.
>
> ## Hooks support
>
> Cloud agents run command-based hooks from `.cursor/hooks.json` in your repository. On Enterprise plans, they also run team hooks and enterprise-managed hooks.
>
> This keeps formatters, audit scripts, and policy checks active when work runs in the cloud. Supported hooks include tool and file hooks (`preToolUse`, `beforeShellExecution`, `afterFileEdit`), plus lifecycle hooks (`beforeSubmitPrompt`, `subagentStart` / `subagentStop`, `preCompact`, `afterAgentResponse` / `afterAgentThought`, and `stop`).
>
> Hooks do not run during early exploratory turns in a read-only environment; they start once the agent has a writable environment. Some hooks are IDE-specific (Tab hooks, `workspaceOpen`). User-level hooks from `~/.cursor/hooks.json` are also not available since cloud VMs don't have access to your local home directory.
>
> See [Hooks: Cloud agent support](https://cursor.com/docs/hooks.md#cloud-agent-support) for the full support matrix and details.
>
> ## Artifacts and remote desktop control
>
> Cloud agents produce merge-ready PRs with artifacts to demo their changes. You can also control the agent's remote desktop to use the modified software.
>
> - **Artifacts**: Agents produce screenshots, videos, and logs so you can see exactly what changed and how the agent verified its work.
> - **Remote desktop control**: Take control of the agent's desktop to test the software yourself in a full development environment without checking out the branch locally. Release control back to the agent for it to keep working.
>
> See [Cloud agent capabilities](https://cursor.com/docs/cloud-agent/capabilities.md) for details on artifacts, computer use, and remote desktop control.
>
> ## Share agents with your team
>
> Send an agent's URL to a teammate and they can open the run to see the conversation, the code changes, and the artifacts the agent produced.
>
> Agents are visible to members of the Cursor team they were started under. Cursor also verifies each viewer's repository access: to open a teammate's agent, connect your own source control account under [Integrations](https://cursor.com/dashboard/integrations) and make sure you have access to the repository the agent worked in. Team membership alone doesn't grant access.
>
> Viewing is read-only. To let teammates send follow-up messages and continue the work, a team admin can enable [team follow-ups](https://cursor.com/docs/cloud-agent/settings.md#team-follow-ups).
>
> ## Related pages
>
> - Learn more about [Cloud agent capabilities](https://cursor.com/docs/cloud-agent/capabilities.md).
> - Learn more about [Cloud agent setup](https://cursor.com/docs/cloud-agent/setup.md).
> - Learn more about [Cloud Agent Builds](https://cursor.com/docs/cloud-agent/builds.md).
> - Learn more about [Cloud agent security](https://cursor.com/docs/cloud-agent/security-network.md).
> - Learn more about [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md).
> - Learn more about [Agent metadata](https://cursor.com/docs/cloud-agent/metadata.md).
> - Learn more about [Cloud agent settings](https://cursor.com/docs/cloud-agent/settings.md).
>
> ## Billing
>
> Cloud Agents are charged at API pricing for the selected [model](https://cursor.com/docs/models-and-pricing.md#model-pricing). You can select the context window size, and a larger context window can increase token usage and costs. You'll be asked to set a spend limit when you first start using them.
>
> ## Troubleshooting
>
> ### Agent runs are not starting
>
> - Ensure you're logged in and have connected your GitHub, GitLab, Azure DevOps, or Bitbucket account.
> - Check that you have the necessary repository permissions.
> - You need to be on a paid Cursor plan.
>
> ### My secrets aren't available to the cloud agent
>
> - Ensure you've added secrets in [cursor.com/dashboard/cloud-agents](https://cursor.com/dashboard/cloud-agents)
> - Secrets are workspace/team-scoped; make sure you're using the correct account
> - Try restarting the cloud agent after adding new secrets
>
> ### Can't find the Secrets tab
>
> - If you don't see it, ensure you have the necessary permissions
>
> ### Do snapshots copy .env.local files?
>
> Snapshots save your base environment configuration (installed packages, system dependencies, etc.).
> If you include `.env.local` files during snapshot creation, they will be saved. However, using the Secrets tab
> in Cursor Settings is the recommended approach for managing environment variables.
>
> ### A teammate can't open my agent
>
> - The teammate must belong to the same Cursor team as you.
> - They need to connect their own source control account under [Integrations](https://cursor.com/dashboard/integrations) and have access to the repository the agent worked in.
> - See [Share agents with your team](https://cursor.com/docs/cloud-agent.md#share-agents-with-your-team) for details.
>
> ### Slack integration not working
>
> Verify that your workspace admin has installed the Cursor Slack app and that
> you have the proper permissions.
>
> ## Naming History
>
> Cloud Agents were formerly called Background Agents.
>
>

### Source: Setup

> # Cloud Environment Setup
>
> Cloud agents run on isolated Ubuntu machines. Configure the environment so the agent has the same repos, tools, dependencies, secrets, and network access a developer would use.
>
> Create a new environment in your [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents#environments).
>
> ## What is a cloud agent environment?
>
> The development environment for a cloud agent is similar to the setup on your laptop: cloned repos, installed dependencies, secrets, startup commands, and network access.
>
> Effective development environments give agents full context on your codebase and organization, so they can test and verify their work.
>
> ![Cloud agent development environment architecture](https://ptht05hbb1ssoooe.public.blob.vercel-storage.com/assets/blog/cloud-agents-architecture-light.png)
>
> ## Why does environment configuration matter?
>
> Agents are only as capable as the environments they run in. An agent that can write code but can't run tests, query services, or reach APIs cannot close the loop on its work.
>
> To take engineering tasks from start to finish, cloud agents need a configured development environment with all the repositories, tools, dependencies, and context to stay autonomous and productive.
>
> Development environments also make agent sessions faster. [Builds](https://cursor.com/docs/cloud-agent/builds.md) prepare repositories, tools, and dependencies in the background so agents start from a ready-to-use machine.
>
> Environment setup is the most important step to improve the effectiveness of your cloud agents.
>
> ## Environment setup options
>
> There are two main ways to configure the environment for your cloud agent:
>
> 1. Let Cursor's agent set up its own environment from the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents#environments). The agent installs dependencies, verifies the environment, and creates its first Build.
> 2. Manually configure the environment with a Dockerfile. If you choose this option, you can specify the Dockerfile in a `.cursor/environment.json` file.
>
> Both options let you specify an install script. Cursor runs it while creating a Build so dependencies are ready before an agent starts.
>
> ### Multi-repo environments
>
> Use a multi-repo environment when an agent needs to work across more than one repository. Select multiple repositories when you create the environment. Cursor clones each selected repo into the agent machine and reuses the environment for future agent runs and automations that use the same repo group.
>
> Multi-repo environments are useful when your frontend, backend, infrastructure, or shared libraries live in separate repos. The agent can inspect the full workspace, make coordinated changes, run tests across repos, and open pull requests in the repos it changes.
>
> You can see which environment is active, along with all past active versions, by visiting the environment's configuration page on the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents#environments).
>
> ### Environment resolution order
>
> Cursor resolves environment configuration by repository or repo group, using the first match:
>
> 1. `.cursor/environment.json` in the repository
> 2. A personal saved environment
> 3. A team saved environment
>
> This gives you predictable defaults at the team level while still letting individual users override with a personal environment when a repo-level `.cursor/environment.json` is not present. User overrides are also useful to allow testing out a new environment configuration before rolling it out to the entire team.
>
> ### Agent-driven setup (recommended)
>
> Cursor can set up your dev environment in the cloud in less than 10 minutes. Start guided setup from the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents#environments) or from the [Agents Window](https://cursor.com/docs/agent/agents-window.md) in the Cursor desktop app.
>
> You will be asked to connect your GitHub, GitLab, Azure DevOps, or Bitbucket account and select one or more repositories.
>
> Then, you provide Cursor with the environment variables and secrets it will need to install dependencies and run the code.
>
> As the agent works, you can watch its progress in a shared terminal session while it handles setup tasks like installing dependencies. Cursor saves the environment after it verifies the code and completes a successful Build.
>
> ![Cloud environment setup in a shared terminal session](https://ptht05hbb1ssoooe.public.blob.vercel-storage.com/assets/changelog/cloud-environment-setup.png)
>
> Future Cloud Agents start from the active Build and can test changes by running your software. Commit the configuration to `.cursor/environment.json` so your whole team benefits.
>
> ### Manual setup with Dockerfile (advanced)
>
> For advanced cases, configure the environment with a Dockerfile:
>
> - Create a Dockerfile to install system-level dependencies, use specific compiler versions, install debuggers, or switch the base OS image
> - Do not `COPY` the full project; Cursor manages the workspace and checks out the correct commit
> - Edit `.cursor/environment.json` directly to configure runtime settings
> - Use build secrets for private package registries or build-time credentials
>
> Here's an example `.cursor/environment.json` referencing a `.cursor/Dockerfile` (relative path) and a `custom_script.sh` install script:
>
> ```json
> {
>   "build": {
>     "dockerfile": "Dockerfile",
>     "context": ".."
>   },
>   "install": "pnpm install && ./custom_script.sh"
> }
> ```
>
> If your repo needs Docker, Tailscale, or Cloudflare Tunnel, see [Running Docker](https://cursor.com/docs/cloud-agent/setup.md#running-docker), [Running Tailscale](https://cursor.com/docs/cloud-agent/setup.md#running-tailscale), and [Running Cloudflare Tunnel](https://cursor.com/docs/cloud-agent/setup.md#running-cloudflare-tunnel) below.
>
> You configure the environment with a Dockerfile; you do not get direct access to the remote machine.
>
> Dockerfile builds use layer caching. When you change a Dockerfile, Cursor rebuilds the changed layers instead of rebuilding every layer from scratch.
>
> ### Cursor-configured Dockerfiles (private beta)
>
> For teams that do not want to write a Dockerfile from scratch, Cursor can configure one for you. During setup, Cursor inspects your repos, identifies tools and dependencies, and produces a Dockerfile-based environment configuration you can edit and version.
>
> This flow is in private beta for Enterprise teams. To request access, contact your Cursor account representative or email [hi@cursor.com](mailto:hi@cursor.com) from your team admin account.
>
> ### Computer Use Support for Dockerfile Repos
>
> Computer use is supported for repos with Dockerfiles based on Debian/Ubuntu-based Linux distributions. If you require support for a different Linux distribution, please contact support.
>
> ### Resource limits
>
> Each cloud agent runs on a default VM profile with limited memory and CPU. If you are on an Enterprise plan and your repo needs more resources, contact support and we can increase limits for your workspace.
>
> Self-serve custom resource configuration is coming soon.
>
> ## Install script
>
> The install script was previously called the update script in the dashboard and docs.
>
> Cursor runs the install script (`install` in `environment.json`) when it creates a [Build](https://cursor.com/docs/cloud-agent/builds.md). The script completes in the background instead of delaying each agent start.
>
> Use `install` for work Cursor can prepare ahead of time. Examples include installing dependencies, generating code, compiling artifacts, and warming disk caches.
>
> Install scripts can read [agent metadata](https://cursor.com/docs/cloud-agent/metadata.md) and mint [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md) from the same local socket the agent uses.
>
> ### Install script idempotency
>
> The install script must be idempotent. It runs for every Build and may run on
> previously prepared disk state.
>
> ### How Builds use the install script
>
> Cursor starts from the environment's base image, clones its repositories, and runs `install` to completion. A successful Build captures the resulting disk state and becomes active. New agents start from the active Build.
>
> Make the script complete. Expensive setup belongs in `install` because it runs before an agent request instead of during startup. Commands such as `pnpm install` can still reuse prepared state and only update changed dependencies.
>
> Builds preserve disk state only. Running processes, exported shell variables, and in-memory caches don't continue into an agent run. Start services with `start` or `terminals`.
>
> ### Environment configuration recovery
>
> A failed Build doesn't replace the active Build. Agents continue to start from the most recent successful environment while you inspect the failure and create a replacement.
>
> Open the environment's **Builds** tab to inspect logs, start an agent from the failed Build, or select another successful Build. See [Cloud Agent Builds](https://cursor.com/docs/cloud-agent/builds.md) for Build controls and debugging.
>
> ### How to decide what to put in your install script
>
> Put every repeatable preparation step in `install`. Include full dependency installation, code generation, artifact compilation, and other work that writes reusable results to disk.
>
> Keep long-running processes out of `install`. Put Docker, databases, tunnels, and dev servers in startup commands. You can also [add instructions in AGENTS.md](https://cursor.com/docs/cloud-agent/setup.md#add-cloud-specific-instructions-to-agentsmd) for services an agent only needs for specific tasks.
>
> ## Startup commands
>
> After an agent boots from a Build, Cursor runs the `start` command and then any configured `terminals`. Use these for processes that should stay alive while the agent runs.
>
> You can skip `start` in many repos. If your environment depends on Docker, add `sudo service docker start` in `start`.
>
> `terminals` are for app code processes. These terminals run in a `tmux` session shared by you and the agent.
>
> ## Add cloud-specific instructions to `AGENTS.md`
>
> Cloud agents read `AGENTS.md` files. We recommend adding a dedicated section for Cloud-only setup and testing instructions, with a title such as `Cursor Cloud specific instructions`.
>
> If this section gets large, we recommend including references to other files that can contain detailed instructions for specific tasks.
>
> See our [AGENTS.md docs](https://cursor.com/docs/rules.md#agentsmd) for more information.
>
> ## Environment variables and secrets
>
> In order to fully run and test code like a human developer, Cloud agents often need environment variables and secrets such as API keys and database credentials.
>
> ### Recommended: use the Secrets tab in Cursor settings
>
> The easiest way to manage secrets is through [cursor.com](https://cursor.com/dashboard/cloud-agents). These are exposed to the cloud agent as environment variables.
>
> For more about the different types of secrets, see our [Secrets documentation](https://cursor.com/docs/cloud-agent/security-network.md#secret-protection). To grant cloud-role access without long-lived keys, see [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md).
>
> ### Environment-scoped secrets
>
> Use environment-scoped secrets when a credential should only be available to agents that use one environment. This is useful for multi-repo environments, staging credentials, or repository groups with different access needs.
>
> Environment-scoped secrets apply to every repo in that environment. They are not available to other environments.
>
> ### Sign-in credentials and 2FA
>
> If your app requires login, add the same credentials you use locally as secrets, such as a username, email, and password.
>
> If your login flow uses TOTP-based 2FA, add the TOTP secret, sometimes called the shared or root secret, as a secret too. The agent can generate the current 6-digit code with `oathtool --totp -b "$TOTP_SECRET"`.
>
> ### Monorepos with multiple `.env` files
>
> If your monorepo has multiple `.env.local` files:
>
> - Add values from all `.env.local` files to the same Secrets tab
> - Use unique variable names when keys overlap, such as `NEXTJS_*` and `CONVEX_*`
> - Reference those variables from each app as needed
>
> If you include `.env.local` files while taking a snapshot, they can be saved and available to cloud agents. The Secrets tab remains the recommended approach for security and management.
>
> ### Using AWS IAM Roles
>
> Cursor supports assuming customer-provided IAM roles for deeper integration with AWS. This allows you to grant specific AWS permissions to cloud agents without sharing long-lived credentials.
>
> 1. **Create the IAM role**: In your AWS account, create the IAM role that you'd like the cloud agent to assume, and note its ARN (e.g. `arn:aws:iam::123456789012:role/acmeRole`).
>
> 2. **Configure the IAM role secret**: Navigate to [Cursor Dashboard → Cloud Agents](https://cursor.com/dashboard?tab=cloud-agents), and add a user or team secret named `CURSOR_AWS_ASSUME_IAM_ROLE_ARN` set to the ARN of the IAM role you created.
>
> 3. **Generate an external ID**: A team admin must do this from the **Advanced** section of team settings. Navigate to [Cursor Dashboard → Settings → Advanced](https://cursor.com/dashboard?tab=settings) and find the External ID settings. If you don't see an external ID displayed, enter a placeholder value in the "AWS IAM Role ARN" field, click "Validate & Save", and reload the page. This will generate an external ID for your team (e.g. `cursor-xxx-yyy-zzz`).
>
> 4. **Configure IAM role trust policy**: In your AWS account, update the IAM role's trust policy to trust Cursor's role assumer. The trust policy should look like this:
>
> ```json
> {
>   "Version": "2012-10-17",
>   "Statement": [
>     {
>       "Sid": "AllowCursorAssume",
>       "Effect": "Allow",
>       "Principal": {
>         "AWS": "arn:aws:iam::289469326074:role/roleAssumer"
>       },
>       "Action": "sts:AssumeRole",
>       "Condition": {
>         "StringEquals": {
>           "sts:ExternalId": "cursor-xxx-yyy-zzz"
>         }
>       }
>     }
>   ]
> }
> ```
>
> Replace `cursor-xxx-yyy-zzz` with the external ID generated for your team.
>
> **Environment variables:**
>
> When configured, Cursor sets these environment variables so AWS tooling uses the `cursor-cloud-agent` profile:
>
> - `AWS_CONFIG_FILE` points to a Cursor-managed AWS config file
> - `AWS_PROFILE` is set to `cursor-cloud-agent`
> - `AWS_SDK_LOAD_CONFIG` is set to `1`
>
> The AWS CLI and AWS SDKs that use the default credential chain pick up this profile automatically during setup commands and while the agent is running. You don't need to export `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, or `AWS_SESSION_TOKEN` yourself.
>
> Cursor assumes the role with STS credentials that expire after 1 hour.
> When the agent wakes, Cursor refreshes credentials that are missing, invalid, or within 15 minutes of expiration.
>
> To federate with AWS STS `AssumeRoleWithWebIdentity`, or with GCP, Azure, or another OIDC verifier, mint [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md) from the agent VM instead of storing long-lived cloud keys.
>
> ## Configuration in code with environment.json
>
> If you prefer to keep your environment configuration defined in code, you can commit a `.cursor/environment.json` to your repository.
>
> Builds use the configuration from the environment's default branch. For feature branch changes, commit and push the configuration, then start an agent on the branch. Cursor checks out the requested branch on top of the active Build, and the agent can rerun the install command when the branch changes dependencies.
>
> Sample `environment.json` using a snapshot-based config (the snapshot ID is accessible from the environments page of the dashboard):
>
> ```json
> {
>   "snapshot": "snapshot-20260212-00000000-0000-0000-0000-000000000000",
>   "install": "npm install"
> }
> ```
>
> Here is a sample `.cursor/environment.json` referencing a `.cursor/Dockerfile` (relative path) and a `custom_script.sh` install script:
>
> ```json
> {
>   "build": {
>     "dockerfile": "Dockerfile",
>     "context": ".."
>   },
>   "install": "pnpm install && ./custom_script.sh"
> }
> ```
>
> ### Important path behavior
>
> The `dockerfile` and `context` paths in `build` are relative to `.cursor`. When
> you omit `context`, it defaults to `.cursor`. The values `.`, `./`, and `..` are
> special-cased to mean the repository root rather than `.cursor`, so to `COPY`
> files that live in `.cursor` with bare filenames, omit `context`. The `install`
> command runs from your project root.
>
> The full schema is [defined here](https://www.cursor.com/schemas/environment.schema.json).
>
> ## Running Docker
>
> Cloud agents support Docker workflows. We use this internally for full-stack repos that run many services.
>
> For simple setups, installing Docker is often enough. Commands like `docker run hello-world` usually work once Docker is installed and the daemon is running.
>
> Docker has edge cases in Cloud Agents because it runs inside another container layer. Simple workflows usually work. More complex setups should start from the `fuse-overlayfs` and `iptables-legacy` configuration below.
>
> For more complex Docker setups, use `fuse-overlayfs`, `iptables-legacy`, and make sure your cloud agent user can run Docker.
>
> ### Recommended Dockerfile for complex Docker setups
>
> ```docker
> ########################################################
> # DOCKER INSTALLATION
> ########################################################
>
> # Install Docker
> RUN install -m 0755 -d /etc/apt/keyrings && \
>     curl --retry 3 --retry-delay 5 -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
>     chmod a+r /etc/apt/keyrings/docker.gpg && \
>     echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
>     $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
>     apt-get update && \
>     apt-get install -y \
>     docker-ce=5:28.5.2-1~ubuntu.24.04~noble \
>     docker-ce-cli=5:28.5.2-1~ubuntu.24.04~noble \
>     containerd.io \
>     docker-buildx-plugin \
>     docker-compose-plugin \
>     && rm -rf /var/lib/apt/lists/*
>
> RUN apt-get update && apt-get install -y fuse-overlayfs && rm -rf /var/lib/apt/lists/*
> RUN mkdir -p /etc/docker && \
>     printf '%s\n' '{' \
>     '  "storage-driver": "fuse-overlayfs"' \
>     '}' > /etc/docker/daemon.json
> RUN apt-get update && apt-get install -y iptables && rm -rf /var/lib/apt/lists/*
> RUN update-alternatives --set iptables /usr/sbin/iptables-legacy && \
>     update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy
>
> ########################################################
> # CONFIG UBUNTU USER
> ########################################################
>
> # ensure no password authentication
> RUN echo 'PasswordAuthentication no\nChallengeResponseAuthentication no\nUsePAM no' > /etc/ssh/sshd_config.d/disable_password_auth.conf
>
> # Create non-root user (only if it doesn't exist)
> RUN id -u ubuntu &>/dev/null || useradd -m -s /bin/bash ubuntu
> # Create docker group if it doesn't exist and add ubuntu user to it
> RUN groupadd -f docker && usermod -aG docker ubuntu
> RUN usermod -aG sudo ubuntu
> # Configure passwordless sudo for ubuntu user
> RUN echo "ubuntu ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/ubuntu
> # Set a password for ubuntu user
> RUN echo "ubuntu:ubuntu" | chpasswd
> ```
>
> ## Running Tailscale
>
> Tailscale does not work in its default networking mode in Cloud agent VMs. Use userspace networking mode instead.
>
> This lets the agent reach private services and data stores through your tailnet without exposing those services to the public internet.
>
> Start `tailscaled` with:
>
> ```bash
> tailscaled --tun=userspace-networking \
>   --outbound-http-proxy-listen=localhost:1054 \
>   --socks5-server=localhost:1055
> ```
>
> Then export these proxy variables in the shell where you want traffic to flow through Tailscale:
>
> ```bash
> export ALL_PROXY=socks5h://localhost:1055/
> export HTTP_PROXY=http://localhost:1054/
> export HTTPS_PROXY=http://localhost:1054/
> ```
>
> After that, run your usual `tailscale up ...` flow.
>
> If you want a working reference, some customers have used [`tailscale-orb`](https://circleci.com/developer/orbs/orb/orbiously/tailscale#commands-connect) successfully because its Docker mode follows this pattern.
>
> Userspace networking does not let the VM appear as a tailnet exit node.
>
> ## Running Cloudflare Tunnel
>
> Cloudflare Tunnel works in Cloud Agent VMs because `cloudflared` runs in userspace.
>
> Use this pattern when a Cloud Agent needs to reach a private HTTP service in a VPC or intranet:
>
> - Install `cloudflared` in your environment Dockerfile or install script.
> - Run a `cloudflared` connector inside your private network.
> - Route an authenticated hostname, such as `vpc.example.com`, through the tunnel to the private origin.
> - Add that hostname to the Cloud Agent network allowlist if your environment uses restricted egress.
> - Store the Cloudflare Access service token values as Cursor Secrets. For example, use `CF_ACCESS_CLIENT_ID` and `CF_ACCESS_CLIENT_SECRET`.
>
> The Cloud Agent can then call the private service over normal HTTPS with the `CF-Access-Client-Id` and `CF-Access-Client-Secret` headers. The connector makes the outbound connection to Cloudflare and forwards the request to your private origin. Your services and data stores stay on your private network, and the connector does not need inbound ports open.
>
> For private TCP services, such as databases, configure a Cloudflare TCP Access app and run `cloudflared access tcp` in your startup command. Point your app or test command at the local listener that `cloudflared` creates.
>
> Keep tunnel tokens and Access service token secrets in Cursor Secrets, not in
> your repository. Rotate them after testing if they were created for a proof of
> concept.
>
>

### Source: Builds

> # Cloud Agent Builds
>
> Builds prepare your Cloud Agent environment in the background. Each agent starts from a pre-built machine with your repositories, tools, and dependencies ready.
>
> With Builds, you get:
>
> - **Faster starts**: Clone, install, and dependency work happen ahead of time, so agents boot from a ready environment instead of waiting on every start.
> - **Reliable starts**: Agents always start from the latest successful Build. A failed install or bad config doesn't replace the working one.
> - **Observable environments**: You can see every Build, inspect its logs and commits, and trace which Build each agent used.
>
> ## How Builds work
>
> A Build is a bootable snapshot of a prepared Cloud Agent environment. Cursor creates Builds ahead of agent runs and keeps the latest successful one ready to start.
>
> Each Build follows this lifecycle:
>
> 1. **Trigger**: A Build starts on a schedule, after you save an environment version, from a manual request, or at an agent's request. See [When Builds occur](https://cursor.com/docs/cloud-agent/builds.md#when-builds-occur).
> 2. **Prepare**: Cursor starts from your base image, clones every repository in the environment at its default branch, and runs the `install` command to completion.
> 3. **Snapshot**: Cursor saves the machine's disk state with the environment version and exact commit SHA for each repository.
> 4. **Activate**: A successful Build becomes active.
> 5. **Start agents**: New agents, automations, and code reviews start from the active Build.
>
> Cursor keeps pre-warmed copies of active Builds ready. This removes repository cloning and dependency installation from the agent startup path.
>
> If a new Build fails, agents continue to use the last successful Build. A broken dependency update, install command, or Dockerfile doesn't replace the active environment.
>
> ## When Builds occur
>
> Cursor starts a Build for four reasons. The Builds tab labels each one with its trigger type.
>
> | Trigger              | When it runs                                                          |
> | :------------------- | :-------------------------------------------------------------------- |
> | Recurring            | On a regular schedule for every environment                           |
> | Configuration change | When you save the environment configuration or change its secrets     |
> | Manual               | When you select **Trigger build** in the Builds tab                   |
> | Agent-requested      | When an agent runs a test Build, for example during environment setup |
>
> ### Recurring Builds
>
> Cursor regularly checks each environment and rebuilds it when something changed. This keeps the active Build close to the head of each repository's default branch, so agents start with fresh code and warm dependency caches instead of pulling and reinstalling at startup.
>
> ### Skipped Builds
>
> A recurring check skips the Build when nothing changed since the last completed one: no new commits on the default branch of any repository in the environment, and no configuration or secret changes. The Builds tab records these checks with a **Skipped** status. They complete in seconds, run no install commands, and leave the active Build in place.
>
> A steady stream of Recurring entries mixing Skipped and Success statuses is the expected state for a healthy environment. Quiet repositories produce mostly Skipped entries. Active repositories rebuild more often.
>
> Cursor only skips recurring Builds. Manual, agent-requested, and configuration-change Builds always run.
>
> ## What runs during a Build and an agent start
>
> Use each environment command for a distinct phase:
>
> | Command     | When it runs                   | Use it for                                                                             |
> | :---------- | :----------------------------- | :------------------------------------------------------------------------------------- |
> | `install`   | During each Build              | Installing dependencies, generating code, compiling artifacts, and warming disk caches |
> | `start`     | At the start of each agent run | Starting Docker, databases, tunnels, and other services                                |
> | `terminals` | At the start of each agent run | Starting app processes in `tmux` terminals shared with the agent                       |
>
> Make `install` complete and idempotent. It can run repeatedly and may run on top of previously prepared disk state. Commands such as `npm install`, `pnpm install`, and `pip install` already support this pattern.
>
> Builds preserve disk state only. Running processes, shell exports, and
> in-memory caches stop when Cursor snapshots the machine. Put services and
> other session-specific work in `start` or `terminals`.
>
> Your existing environment inputs still apply. Builds use saved snapshots, `.cursor/environment.json`, Dockerfiles, install and startup commands, secrets, and network settings.
>
> ## How Builds handle Git state
>
> A Build records the commit checked out for each repository when it runs.
>
> - **Default branch runs** start from the commit recorded in the active Build. Scheduled Builds refresh that commit in the background. When **Update stale builds** is on and a Build is older than your **Staleness threshold**, agents pull the latest default-branch code at start. When that setting is off, agents use the Build's recorded commit as-is. The default threshold is 24 hours. Set it to `0` to always pull.
> - **Feature branch runs** start from the active Build's prepared disk, then Cursor checks out the requested branch. The source code matches the branch you selected while reusing dependencies from the Build.
> - **Multi-repo environments** record one commit per repository and prepare the complete workspace together.
>
> If a feature branch changes dependencies, the agent receives your environment context and install command so it can refresh the environment before testing.
>
> ## How secrets work with Builds
>
> Builds can access team and environment secrets. Use these for private package registries, artifact stores, and other credentials required by `install`.
>
> User secrets are added only when an agent starts. They aren't available during Builds and don't become part of a shared snapshot.
>
> Saving environment configuration or changing its secrets triggers a new Build.
>
> ## Manage Builds
>
> Open an environment's **Builds** tab to:
>
> - See every Build's type, status, and start time
> - Open a Build to inspect its details and logs
> - Select **Trigger build** to run a Build on demand
> - Activate a draft Build or deactivate a Build
> - Cancel an in-progress Build
> - Start an agent from a specific Build
> - Configure **Update stale builds** and the **Staleness threshold**
>
> Every agent run records the Build it started from. Use this provenance to compare environment behavior with the exact configuration and repository commits in the Build.
>
> ## Debug a Build
>
> Open a failed Build to inspect its events and logs. Agents continue to start from the active successful Build while you diagnose the failure.
>
> For an exact reproduction, start an agent from the failed Build. The agent opens the machine in its failed state, where it can inspect logs, update the environment, run a test Build, and verify the result.
>
> You can also ask a Cloud Agent to inspect and manage Builds through the built-in [Cursor Cloud MCP](https://cursor.com/docs/cloud-agent/capabilities.md#cursor-cloud-mcp). For example:
>
> ```text
> Inspect the latest failed Build for this environment. Fix the environment
> configuration, run a test Build, and verify it before proposing the final
> install and start commands.
> ```
>
> ## Build behavior reference
>
> ### Which Build does an agent use?
>
> By default, an agent uses the latest successful active Build for its environment.
>
> You can also start an agent from a specific Build when testing or debugging.
>
> ### What happens before the first successful Build?
>
> Agents use the standard environment startup flow until the first Build completes successfully. A failed Build doesn't interrupt existing agent workflows.
>
> ### How fresh is the source code?
>
> Feature branch runs check out the requested branch after the Build starts. Default branch runs begin at the commit recorded by the active Build. If **Update stale builds** is on and the Build is older than your **Staleness threshold**, agents pull the latest default-branch code at start.
>
> ### Do Builds replace snapshots or Dockerfiles?
>
> No. A saved snapshot or Dockerfile defines the base machine used to create a Build. Cursor then clones the repositories, runs `install`, and creates a fresh bootable snapshot.
>
> ### Do Builds support multiple repositories?
>
> Yes. One Build prepares all repositories in the environment and records the commit used for each one.
>
> ### Do Builds cost extra?
>
> No. Builds are included with Cloud Agents.
>
> ## Related
>
> - [Cloud Environment Setup](https://cursor.com/docs/cloud-agent/setup.md)
> - [Cloud Agent capabilities](https://cursor.com/docs/cloud-agent/capabilities.md)
> - [Cloud Agents settings](https://cursor.com/docs/cloud-agent/settings.md)
>
>

### Source: Best practices

> # Best Practices
>
> Use these recommendations to get more reliable Cloud Agent runs.
>
> ## Set up the environment first
>
> Use [Cloud agent setup](https://cursor.com/docs/cloud-agent/setup.md) so that Cursor has its environment configured. Like a human developer, Cursor does better work if its environment is set up correctly.
>
> ## Ensure the agent can access what it needs
>
> Before running a Cloud Agent, verify these prerequisites:
>
> - **Secrets**: Make sure the agent has access to required secrets (API keys, database credentials, etc.) through the [Secrets tab](https://cursor.com/dashboard/cloud-agents) in your dashboard. For cloud roles, prefer [OIDC tokens](https://cursor.com/docs/cloud-agent/identity.md) over long-lived access keys.
> - **Egress controls**: If you have [network access](https://cursor.com/docs/cloud-agent/security-network.md) restrictions enabled, ensure all URLs your local development requires are whitelisted.
> - **Local testability**: Your repo should be set up to run well locally without requiring external services that cannot be reached from a VM. If it is hard for a human developer to test locally, it will also be hard for an agent.
>
> ## Use skills and agents.md to configure your agent
>
> If the cloud agent is having difficulty testing its changes, we recommend using [skills](https://cursor.com/docs/skills.md) and agents.md to configure your agent.
>
> Think of the agent as a smart, but low-context human developer. The best way to make sure it does the right thing is to give it the context it needs to understand what to do.
>
> For example, at Cursor our agents.md lists tips for running and debugging the most commonly used microservices in our mono-repo. We also have lots of skills about how to test and debug key services, each with clear instructions on when to use the skill.
>
> The skills contain in-depth details, such as how to debug a specific microservice or how to set up a third-party dependency when needed for testing.
>
> ## Use rules to enforce conventions
>
> Cloud Agents can read and follow [Rules](https://cursor.com/docs/rules.md) at three levels:
>
> - **User rules**: Set in Cursor Settings, these apply to your sessions across all repositories. Best for rules you only want to apply to you personally.
> - **Team rules**: Set in the [Rules, Commands, Hooks dashboard](https://cursor.com/dashboard/team-content), these apply to all team members across every repository. Best for org-wide conventions.
> - **Repo rules**: `.cursor/rules/*.mdc` files committed to the repository, these apply to all agents using that repository. Best for repo/project-specific conventions.
>
> ## Give the agent the tools it needs
>
> We have often found that agents are limited by the tools they have access to. We recommend using MCP and creating custom tools so that the agent has access to the same systems a human developer would.
>
> ## Mold the tools to the agent
>
> It is important to create tools that the agent is good at using. We recommend creating tools, and iterating based on observations of how the agent uses them.
>
> For example, at Cursor we have created a custom CLI for the model to run micro-services in our codebase. We found that when running custom dev commands, e.g. from a package.json file, some models would forget arguments, or agents would get distracted by noisy build logs which human developers knew to ignore.
>
>

### Source: Settings

> # Cloud Agents settings
>
> Workspace admins can configure Cloud Agents from the [Cloud Agents dashboard](https://cursor.com/dashboard/cloud-agents).
>
> ## Environment management
>
> The **Environments** view lists the saved environments available to your team. Environments can be scoped to one repo or to a group of repos.
>
> Open an environment to review:
>
> - The repositories it applies to
> - Whether it uses a snapshot or `.cursor/environment.json`
> - The install script that runs during [Builds](https://cursor.com/docs/cloud-agent/builds.md)
> - Runtime secrets and build secrets
> - Network access settings
> - Version history, Builds, and setup runs
>
> Use **Update with Agent** when you want Cursor to inspect the current environment and propose a new setup. Use **New Setup Run** when you want Cursor to start setting up the environment fresh. Use **Restore** from version history to make a prior environment version active again.
>
> The **Builds** tab shows the prepared environment versions available to Cloud Agents. You can inspect logs, trigger a Build, choose or pin the active Build, and start an agent from a specific Build. See [Cloud Agent Builds](https://cursor.com/docs/cloud-agent/builds.md) for details.
>
> ## Default settings
>
> - **Default model** – the model used when a run does not specify one. Pick any model available for cloud agents.
> - **Default repository** – when empty, agents ask the user to choose a repo. Supplying a repo here lets users skip that step.
> - **Base branch** – the branch agents fork from when creating pull requests. Leave blank to use the repository’s default branch.
>
> ## Network access settings
>
> Control which network resources Cloud Agents can reach. User and team settings support three modes:
>
> - **Allow all network access** – no domain restrictions.
> - **Default + allowlist** – the [default domains](https://cursor.com/docs/agent/security/run-modes.md#network-access) plus any domains you add.
> - **Allowlist only** – only domains you explicitly add.
>
> Users, team admins, and environment owners can configure network access. Environment-level settings can inherit user or team policy, add an environment allowlist, or define their own access mode. See [Network Access](https://cursor.com/docs/cloud-agent/security-network.md) for full details.
>
> ## Security settings
>
> All security options require admin privileges.
>
> - **Display agent summary** – controls whether Cursor shows the agent's file-diff images and code snippets. Disable this if you prefer not to expose file paths or code in the sidebar.
> - **Display agent summary in external channels** – extends the previous toggle to Slack or any external channel you've connected.
> - **Team follow-ups** – controls whether team members can send follow-up messages to cloud agents created by other users on the team. See [team follow-ups](https://cursor.com/docs/cloud-agent/settings.md#team-follow-ups) below.
>
> ## Team feature settings
>
> Team admins can enable or disable these features for their team:
>
> - **Long running agents** – controls whether team members can run agents for extended durations. Admins can enable or restrict this capability at the team level. Long-running is not available for multi-repo environments yet. Selecting a multi-repo environment disables the toggle.
> - **Computer use** – controls whether agents can use computer interaction capabilities (available to enterprise teams only).
>
> Changes save instantly and affect new agents immediately.
>
> ### Team follow-ups
>
> Team members can send follow-up messages to cloud agents created by other users on the same team. This is useful when a teammate starts an agent and you need to course-correct, add context, or continue the work while they're unavailable.
>
> Follow-ups build on agent visibility. A teammate must be able to view the agent before they can send follow-ups: they must belong to the same Cursor team, and they need their own access to the agent's repository. See [Share agents with your team](https://cursor.com/docs/cloud-agent.md#share-agents-with-your-team).
>
> Team admins control this behavior from the [Cloud Agents security settings](https://cursor.com/dashboard/cloud-agents) with three options:
>
> | Setting                   | Behavior                                                                                                                                                                                   |
> | ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | **Disabled**              | Only the original creator can send follow-ups to their agent. No team follow-ups are allowed.                                                                                              |
> | **Service accounts only** | Team members can send follow-ups to agents created by a [service account](https://cursor.com/docs/account/enterprise/service-accounts.md), but not to agents created by other human users. |
> | **All**                   | Any team member can send follow-ups to any agent on the team, regardless of who created it.                                                                                                |
>
> ### Lateral movement and secret exposure
>
> Enabling team follow-ups means a user can influence the execution of a cloud agent that runs with *another user's* secrets and credentials. A follow-up message can instruct the agent to read environment variables, print secrets to logs, push credentials to an external endpoint, or perform actions using the original creator's access tokens.
>
> A team member with limited permissions could escalate their access by directing an agent that holds a more privileged user's secrets. Treat this setting with the same care you would give shared SSH keys or service credentials.
>
>
