---
primary_sources:
  - id: T1-CLOUD
    title: "Codex cloud"
    url: "https://learn.chatgpt.com/docs/cloud.md"
    section: "full page"
  - id: T1-ENV-CLOUD
    title: "Cloud environments"
    url: "https://learn.chatgpt.com/docs/environments/cloud-environment.md"
    section: "full page"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Codex cloud

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex cloud — full page

> # Codex cloud
>
> > For the complete documentation index, see [llms.txt](https://learn.chatgpt.com/llms.txt). Markdown versions of documentation pages are available by appending `.md` to the page URL.
>
> ## Run coding tasks in parallel cloud environments
>
> Run tasks in isolated cloud environments, work in parallel, and start work from the web, GitHub, GitLab, Linear, or Slack.
>
> > Illustration: Codex cloud chat composer and chat list with interactive archiving
>
> ### Start here
>
> - [Open Codex cloud](https://chatgpt.com/codex)
> - [Set up Codex cloud](#getting-started)
>
> ### Why use Codex cloud
>
> - **Run work in parallel:** Give longer tasks dedicated environments and let them continue while you work on something else.
> - **Reproduce the environment:** Configure the dependencies, tools, variables, and setup steps each repository needs.
> - **Review before you merge:** Inspect the summary and diff, request a follow-up, or open a pull request when the result is ready.
>
> ## Getting started
>
> **Set up Codex cloud.**
>
> Connect GitHub or GitLab, create an environment, and start your first cloud chat.
>
> ### 1. Open Codex and sign in
>
> Go to [Codex](https://chatgpt.com/codex) and sign in with your ChatGPT account.
>
> ### 2. Connect GitHub or GitLab
>
> Connect GitHub or GitLab (Beta) when prompted. For GitHub, choose the repositories Codex can access; for GitLab, select a project when you create the environment. For GitLab setup, webhook permissions, and merge request reviews, see [Use Codex with GitLab (Beta)](https://learn.chatgpt.com/docs/third-party/gitlab).
>
> ### 3. Create an environment
>
> Open [environment settings](https://chatgpt.com/codex/settings/environments) and create an environment for the repository you selected. Configure any dependencies, tools, environment variables, or secrets the task needs.
>
> For configuration details, see [Cloud environments](https://learn.chatgpt.com/docs/environments/cloud-environment).
>
> ### 4. Start your first task
>
> Return to [Codex](https://chatgpt.com/codex), choose your environment, and describe the result you want. You can watch the task logs or let the task run in the background.
>
> ### 5. Review the result
>
> Review the summary and diff. Ask Codex to make follow-up changes, or open a pull request when the work is ready.
>
> ### Next steps
>
> - [Customize the cloud environment](https://learn.chatgpt.com/docs/environments/cloud-environment)
> - [Configure agent internet access](https://learn.chatgpt.com/docs/cloud/internet-access)
> - [Use Codex with GitHub](https://learn.chatgpt.com/docs/third-party/github)
> - [Use Codex with GitLab (Beta)](https://learn.chatgpt.com/docs/third-party/gitlab)
> - [Use Codex in Linear](https://learn.chatgpt.com/docs/third-party/linear)
> - [Use Codex in Slack](https://learn.chatgpt.com/docs/third-party/slack)
>
> ## See what Codex cloud can do
>
> Give each task the environment it needs, then review the result on your schedule.
>
> - [Delegate several tasks](https://learn.chatgpt.com/docs/environments/cloud-environment): Start work in parallel and return as each task reaches a reviewable result.
> - [Build a reproducible environment](https://learn.chatgpt.com/docs/environments/cloud-environment): Configure the dependencies, tools, variables, and setup steps a repository needs.
> - [Delegate from your integrations](https://learn.chatgpt.com/docs/developers): Start work in Codex cloud from GitHub pull requests, GitLab merge requests and issues, Linear issues, or Slack channels and threads.
>
> ## Use Codex cloud when…
>
> - [Work needs to run in the background](https://learn.chatgpt.com/docs/environments/cloud-environment): Delegate a longer task and return when it is ready.
> - [You want to compare several attempts](https://learn.chatgpt.com/docs/environments/cloud-environment): Run tasks in parallel without tying up your local machine.
> - [Work starts in GitHub, GitLab, Linear, or Slack](https://learn.chatgpt.com/docs/developers): Use integrations to hand off work without leaving the pull request, merge request, issue, channel, or thread.
> - [You are away from your development machine](https://learn.chatgpt.com/docs/environments/cloud-environment): Start and review work from the web or Codex CLI.

### Source: Cloud environments — full page

> # Cloud environments
>
> > For the complete documentation index, see [llms.txt](https://learn.chatgpt.com/llms.txt). Markdown versions of documentation pages are available by appending `.md` to the page URL.
>
> Use environments to control what Codex installs and runs during cloud chats. For example, you can add dependencies, install tools like linters and formatters, and set environment variables.
>
> Configure environments in [Codex settings](https://chatgpt.com/codex/settings/environments).
>
> <a id="how-codex-cloud-tasks-run"></a>
>
> ## How Codex cloud chats run
>
> Here's what happens when you submit a prompt:
>
> 1. Codex creates a container and checks out your repo at the selected branch or commit SHA.
> 2. Codex runs your setup script, plus an optional maintenance script when a cached container is resumed.
> 3. Codex applies your internet access settings. Setup scripts run with internet access. Agent internet access is off by default, but you can enable limited or unrestricted access if needed. See [agent internet access](https://learn.chatgpt.com/docs/cloud/internet-access).
> 4. The agent runs terminal commands in a loop. It edits code, runs checks, and tries to validate its work. If your repo includes `AGENTS.md`, the agent uses it to find project-specific lint and test commands.
> 5. When the agent finishes, it shows its answer and a diff of any files it changed. You can open a PR or ask follow-up questions.
>
> ## Default universal image
>
> The Codex agent runs in a default container image called `universal`, which comes pre-installed with common languages, packages, and tools.
>
> In environment settings, select **Set package versions** to pin versions of Python, Node.js, and other runtimes.
>
> For details on what's installed, see
>   [openai/codex-universal](https://github.com/openai/codex-universal) for a
>   reference Dockerfile and an image that can be pulled and tested locally.
>
> While `codex-universal` comes with languages pre-installed for speed and convenience, you can also install additional packages to the container using [setup scripts](#manual-setup).
>
> ## Environment variables and secrets
>
> **Environment variables** are set for the full duration of the chat (including setup scripts and the agent phase).
>
> **Secrets** are similar to environment variables, except:
>
> - They are stored with an additional layer of encryption and are only decrypted for task execution.
> - They are only available to setup scripts. For security reasons, secrets are removed before the agent phase starts.
>
> ## Automatic setup
>
> For projects using common package managers (`npm`, `yarn`, `pnpm`, `pip`, `pipenv`, and `poetry`), Codex can automatically install dependencies and tools.
>
> ## Manual setup
>
> If your development setup is more complex, you can also provide a custom setup script. For example:
>
> ```bash
> # Install type checker
> pip install pyright
>
> # Install dependencies
> poetry install --with test
> pnpm install
> ```
>
> Setup scripts run in a separate Bash session from the agent, so commands like
>   `export` do not persist into the agent phase. To persist environment
>   variables, add them to `~/.bashrc` or configure them in environment settings.
>
> ## Container caching
>
> Codex caches container state for up to 12 hours to speed up new chats and follow-ups.
>
> When an environment is cached:
>
> - Codex clones the repository and checks out the default branch.
> - Codex runs the setup script and caches the resulting container state.
>
> When a cached container is resumed:
>
> - Codex checks out the branch specified for the chat.
> - Codex runs the maintenance script (optional). This is useful when the setup script ran on an older commit and dependencies need to be updated.
>
> Codex automatically invalidates the cache if you change the setup script, maintenance script, environment variables, or secrets. If your repo changes in a way that makes the cached state incompatible, select **Reset cache** on the environment page.
>
> For Business and Enterprise users, caches are shared across all users who have
>   access to the environment. Invalidating the cache will affect all users of the
>   environment in your workspace.
>
> ## Internet access and network proxy
>
> Internet access is available during the setup script phase to install dependencies. During the agent phase, internet access is off by default, but you can configure limited or unrestricted access. See [agent internet access](https://learn.chatgpt.com/docs/cloud/internet-access).
>
> Environments run behind an HTTP/HTTPS network proxy for security and abuse prevention purposes. All outbound internet traffic passes through this proxy.
