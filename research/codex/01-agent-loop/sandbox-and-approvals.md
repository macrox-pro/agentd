---
primary_sources:
  - id: T1-SANDBOX
    title: "Sandbox"
    url: "https://learn.chatgpt.com/docs/sandboxing.md"
    section: "full page"
  - id: T1-APPROVALS
    title: "Agent approvals & security"
    url: "https://learn.chatgpt.com/docs/agent-approvals-security.md"
    section: "full page"
also_cited_in: [T1-WIN-SANDBOX, T1-WSL]
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Sandbox and approvals

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Sandbox — full page

> # Sandbox
>
> > For the complete documentation index, see [llms.txt](https://learn.chatgpt.com/llms.txt). Markdown versions of documentation pages are available by appending `.md` to the page URL.
>
> <ContentModeSwitch group="codex-surface" ids="app,cli,ide">
>
> The sandbox is the boundary that lets the agent act autonomously without giving it
> unrestricted access to your machine. When a local chat runs commands in the
> **ChatGPT desktop app**, **Codex CLI**, or **IDE extension**, those commands run inside a
> constrained environment instead of running with full access by default.
>
> That environment defines what the agent can do on its own, such as which files it
> can modify and whether commands can use the network. When a task stays inside
> those boundaries, the agent can keep moving without stopping for confirmation. When
> it needs to go beyond them, the approval flow takes over.
>
> Sandboxing and approvals are different controls that work together. The
>   sandbox defines technical boundaries. The approval policy decides when the
>   agent must stop and ask before crossing them.
>
> ## What the sandbox does
>
> The sandbox applies to spawned commands, not just to built-in file
> operations. If the agent runs tools like `git`, package managers, or test runners,
> those commands inherit the same sandbox boundaries.
>
> Codex uses platform-native enforcement on each OS. The implementation differs
> between macOS, Linux, WSL2, and native Windows, but the idea is the same across
> surfaces: give the agent a bounded place to work so routine tasks can run
> autonomously inside clear limits.
>
> ## Why it matters
>
> The sandbox reduces approval fatigue. Instead of asking you to confirm every
> low-risk command, the agent can read files, make edits, and run routine project
> commands within the boundary you already approved.
>
> It also gives you a clearer trust model for agentic work. You aren't just
> trusting the agent's intentions; you are trusting that the agent is operating
> inside enforced limits. That makes it easier to let the agent work independently
> while still knowing when it will stop and ask for help.
>
> ## Getting started
>
> The default permissions mode applies sandboxing automatically.
>
> ### Prerequisites
>
> On **macOS**, sandboxing works out of the box using the built-in Seatbelt
> framework.
>
> On **Windows**, Codex uses the native [Windows
> sandbox](https://learn.chatgpt.com/docs/windows/windows-sandbox#windows-sandbox) when you run in PowerShell and the
> Linux sandbox implementation when you run in WSL2.
>
> On **Linux and WSL2**, install `bubblewrap` with your package manager first:
>
> <Tabs
>   id="codex-sandboxing-prerequisites"
>   param="sandbox-os"
>   tabs={[
>     { id: "ubuntu-debian", label: "Ubuntu/Debian" },
>     { id: "fedora", label: "Fedora" },
>   ]}
> >
>   
>
>
> ```bash
> sudo apt install bubblewrap
> ```
>
>   
>
>
>   
>
>
> ```bash
> sudo dnf install bubblewrap
> ```
>
>   
>
> </Tabs>
>
> Codex uses the first `bwrap` executable it finds on `PATH`. If no `bwrap`
> executable is available, Codex falls back to a bundled helper, but that helper
> requires support for unprivileged user namespace creation. Installing the
> distribution package that provides `bwrap` keeps this setup reliable.
>
> Codex surfaces a startup warning when `bwrap` is missing or when the helper
> can't create the needed user namespace. On distributions that restrict this
> AppArmor setting, prefer loading the `bwrap` AppArmor profile so `bwrap` can
> keep working without disabling the restriction globally.
>
> **Ubuntu AppArmor note:** On Ubuntu 25.04, installing `bubblewrap` from
>   Ubuntu's package repository should work without extra AppArmor setup. The
>   `bwrap-userns-restrict` profile ships in the `apparmor` package at
>   `/etc/apparmor.d/bwrap-userns-restrict`.
>
> On Ubuntu 24.04, Codex may still warn that it can't create the needed user
> namespace after `bubblewrap` is installed. Copy and load the extra profile:
>
> ```bash
> sudo apt update
> sudo apt install apparmor-profiles apparmor-utils
> sudo install -m 0644 \
>   /usr/share/apparmor/extra-profiles/bwrap-userns-restrict \
>   /etc/apparmor.d/bwrap-userns-restrict
> sudo apparmor_parser -r /etc/apparmor.d/bwrap-userns-restrict
> ```
>
> `apparmor_parser -r` loads the profile into the kernel without a reboot. You
> can also reload all AppArmor profiles:
>
> ```bash
> sudo systemctl reload apparmor.service
> ```
>
> If that profile is unavailable or does not resolve the issue, you can disable
> the AppArmor unprivileged user namespace restriction with:
>
> ```bash
> sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0
> ```
>
> </ContentModeSwitch>
>
> ## How permissions work
>
> <ContentModeSwitch group="codex-surface" ids="app,cli,ide">
>
> Use the permissions control for your surface to change how Codex handles local
> actions.
>
> Approvals determine when Codex pauses before an action, while the sandbox
> determines which files and network resources commands can access. When an
> approval offers different scopes, such as approving once or for the session,
> choose the narrowest scope that lets the task continue. Keep the project
> boundary as the default; use separate projects or worktrees instead of
> broadening access across unrelated repositories.
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="web">
>
> ChatGPT Work runs code and shell commands in a managed, isolated environment.
> Workspace policy and tool-specific controls determine which capabilities are
> available. When the setting is available, use **Settings > Data controls > Work
> network access** to manage network access for code and shell commands. Turn on
> **Allow public internet access** to let those commands reach the public
> internet. When it's off, commands can reach only required hostnames from a
> managed allowlist.
>
> Web search, plugins, and the remote browser have separate controls.
> Changes take effect after the current code or shell run finishes and Work
> refreshes its execution environment. ChatGPT web doesn't expose the local
> Codex sandbox or approval-mode selector.
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="app">
>
> In the ChatGPT desktop app, use the permissions control beneath the composer.
> Depending on your configuration, the menu can include **Ask for approval**,
> **Approve for me** for eligible approval requests, **Full access**, and named or
> custom permissions profiles.
>
> <PermissionModeSelectorDemo client:load />
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="cli">
>
> In the CLI, enter
> [`/permissions`](https://learn.chatgpt.com/docs/developer-commands?surface=cli#cli-update-permissions-with-permissions)
> to open the permissions picker and change the active permissions profile.
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="ide">
>
> In the IDE extension, use the permissions control beneath the composer.
> Depending on your configuration, the menu can include **Ask for approval**,
> **Approve for me** for eligible approval requests, **Full access**, and named or
> custom permissions profiles.
>
>
>
>   
>     
>
> > Illustration: Codex approval mode selector in the IDE extension
>
>
>   
>
>
>
> </ContentModeSwitch>
>
> <a id="configure-defaults"></a>
>
> <ContentModeSwitch group="codex-surface" ids="app,cli,ide">
>
> ## Configure defaults
>
> To start with the same behavior every time, set defaults in `config.toml`.
> [Config basics](https://learn.chatgpt.com/docs/config-file/config-basic) explains how it works, and the
> [Configuration reference](https://learn.chatgpt.com/docs/config-file/config-reference) documents the exact keys for
> `sandbox_mode`, `approval_policy`, `approvals_reviewer`, and
> `sandbox_workspace_write.writable_roots`. Use those settings to decide how much
> autonomy the agent gets by default, which directories it can write to, when it
> should pause for approval, and who reviews eligible approval requests.
>
> At a high level, the common sandbox modes are:
>
> - `read-only`: The agent can inspect files, but it can't edit files or run
>   commands without approval.
> - `workspace-write`: The agent can read files, edit within the workspace, and run
>   routine local commands inside that boundary. This is the default low-friction
>   mode for local work.
> - `danger-full-access`: The agent runs without sandbox restrictions. This removes
>   the filesystem and network boundaries and should be used only when you want
>   the agent to act with full access.
>
> The common approval policies are:
>
> - `untrusted`: The agent asks before running commands that aren't in its trusted
>   set.
> - `on-request`: The agent works inside the sandbox by default and asks when it
>   needs to go beyond that boundary.
> - `never`: The agent doesn't stop for approval prompts.
>
> When approvals are interactive, you can also choose who reviews them with
> `approvals_reviewer`:
>
> - `user`: approval prompts surface to the user. This is the default.
> - `auto_review`: eligible approval prompts go to a reviewer agent (see
>   [automatic review](https://learn.chatgpt.com/docs/sandboxing/auto-review)).
>
> Full access means using `sandbox_mode = "danger-full-access"` together with
> `approval_policy = "never"`. By contrast, the lower-risk local automation
> preset is `sandbox_mode = "workspace-write"` together with
> `approval_policy = "on-request"`, or the matching CLI flags
> `--sandbox workspace-write --ask-for-approval on-request`. You can then keep
> `approvals_reviewer = "user"` for manual approvals or set
> `approvals_reviewer = "auto_review"` for automatic approval review.
>
> If you need the agent to work across more than one directory, writable roots let
> you extend the places it can modify without removing the sandbox entirely. If
> you need a broader or narrower trust boundary, adjust the default sandbox mode
> and approval policy instead of relying on one-off exceptions.
>
> When a workflow needs a specific exception, use [rules](https://learn.chatgpt.com/docs/agent-configuration/rules). Rules
> let you allow, prompt, or forbid command prefixes outside the sandbox, which is
> often a better fit than broadly expanding access. For IDE-specific settings
> entry points, see [Codex IDE extension settings](https://learn.chatgpt.com/docs/developer-settings?surface=ide).
>
> Automatic review, when available, doesn't change the sandbox boundary. It's
> one possible `approvals_reviewer` for approval requests at that boundary, such
> as sandbox escalations, blocked network access, or side-effecting tool calls
> that still need approval. Actions already allowed inside the sandbox run
> without extra review. For the reviewer lifecycle, trigger types, denial
> semantics, and configuration details, see
> [automatic review](https://learn.chatgpt.com/docs/sandboxing/auto-review).
>
> Platform details live in the platform-specific docs. For native Windows setup,
> behavior, and troubleshooting, see [Windows](https://learn.chatgpt.com/docs/windows/windows-sandbox). For admin
> requirements and organization-level constraints on sandboxing and approvals, see
> [Agent approvals & security](https://learn.chatgpt.com/docs/agent-approvals-security).
>
> </ContentModeSwitch>

### Source: Agent approvals & security — full page

> # Agent approvals & security
>
> > For the complete documentation index, see [llms.txt](https://learn.chatgpt.com/llms.txt). Markdown versions of documentation pages are available by appending `.md` to the page URL.
>
> Codex helps protect your code and data and reduces the risk of misuse.
>
> This page covers how to operate Codex safely, including sandboxing, approvals,
>   and network access. If you are looking for Codex Security, the product for
>   scanning connected GitHub repositories, see [Codex Security](https://learn.chatgpt.com/docs/security).
>
> By default, the agent runs with network access turned off. Locally, Codex uses an OS-enforced sandbox that limits what it can touch (typically to the current workspace), plus an approval policy that controls when it must stop and ask you before acting.
>
> For a high-level explanation of how sandboxing works across the ChatGPT desktop app,
> Codex CLI, and IDE extension, see [sandboxing](https://learn.chatgpt.com/docs/sandboxing).
> For a broader enterprise security overview, see the [Codex security white paper](https://trust.openai.com/?itemUid=382f924d-54f3-43a8-a9df-c39e6c959958&source=click).
>
> ## Sandbox and approvals
>
> Codex security controls come from two layers that work together:
>
> - **Sandbox mode**: What Codex can do technically (for example, where it can write and whether it can reach the network) when it executes model-generated commands.
> - **Approval policy**: When Codex must ask you before it executes an action (for example, leaving the sandbox, using the network, or running commands outside a trusted set).
>
> Codex uses different sandbox modes depending on where you run it:
>
> - **Codex cloud**: Runs in isolated OpenAI-managed containers, preventing access to your host system or unrelated data. Uses a two-phase runtime model: setup runs before the agent phase and can access the network to install specified dependencies, then the agent phase runs offline by default unless you enable internet access for that environment. Secrets configured for cloud environments are available only during setup and are removed before the agent phase starts.
> - **Codex CLI / IDE extension**: OS-level mechanisms enforce sandbox policies. Defaults include no network access and write permissions limited to the active workspace. You can configure the sandbox, approval policy, and network settings based on your risk tolerance.
>
> In the `Auto` preset (for example, `--sandbox workspace-write --ask-for-approval on-request`), Codex can read files, make edits, and run commands in the working directory automatically.
>
> Codex asks for approval to edit files outside the workspace or to run commands that require network access. If you want to chat or plan without making changes, switch to `read-only` mode with the `/permissions` command.
>
> Codex can also elicit approval for app (connector) tool calls that advertise side effects, even when the action isn't a shell command or file change. Destructive app/MCP tool calls always require approval when the tool advertises a destructive annotation (unless the tool advertises a read annotation, which takes priority).
>
> ## Network access <ElevatedRiskBadge class="ml-2" />
>
> For Codex cloud, see [agent internet access](https://learn.chatgpt.com/docs/cloud/internet-access) to enable full internet access or a domain allow list.
>
> For the ChatGPT desktop app, Codex CLI, or IDE extension, the default `workspace-write` sandbox mode keeps network access turned off unless you enable it in your configuration:
>
> ```toml
> [sandbox_workspace_write]
> network_access = true
> ```
>
> ### Network isolation
>
> Network access is controlled through destination rules that apply to scripts,
> programs, and subprocesses spawned by commands. When command network access is
> already enabled, turn on the `network_proxy` feature to constrain that traffic
> to the network policy you configure. Adding domain rules does not enable the
> proxy by itself.
>
> ```toml
> [features.network_proxy]
> enabled = true
> domains = { "api.openai.com" = "allow", "example.com" = "deny" }
> ```
>
> For a one-off CLI session, use the boolean shorthand when you only need the
> toggle, and the table form when you also set policy options:
>
> ```bash
> codex \
>   -c 'features.network_proxy=true' \
>   -c 'sandbox_workspace_write.network_access=true'
>
> codex \
>   -c 'features.network_proxy.enabled=true' \
>   -c 'features.network_proxy.domains={ "api.openai.com" = "allow", "example.com" = "deny" }' \
>   -c 'sandbox_workspace_write.network_access=true'
> ```
>
> The feature changes how enabled network access is enforced; it does not grant
> network access by itself. Use `sandbox_workspace_write.network_access` with
> `workspace-write` config to decide whether commands have network access at all:
>
> - Network off + `network_proxy` on: network stays off, and the feature does nothing.
> - Network on + `network_proxy` off: network stays on with unrestricted direct
>   outbound access.
> - Network on + `network_proxy` on: network stays on, and outbound traffic is
>   constrained by the configured network policy.
>
> The proxy feature also applies to [permission profiles](https://learn.chatgpt.com/docs/permissions#network-permissions).
> A profile's `network.enabled = true` grants command network access, while
> `features.network_proxy = true` activates enforcement of that profile's domain
> rules:
>
> ```toml
> default_permissions = "project-edit"
>
> [features]
> network_proxy = true
>
> [permissions.project-edit]
> extends = ":workspace"
>
> [permissions.project-edit.network]
> enabled = true
>
> [permissions.project-edit.network.domains]
> "api.openai.com" = "allow"
> ```
>
> If you omit the proxy feature in this example, commands have direct network
> access and the `api.openai.com` allow rule does not restrict their destinations.
>
> Admin-managed `experimental_network` requirements are separate from the user
> feature toggle. They can configure and start sandboxed networking without
> `features.network_proxy`, but they do not turn on network access when the active
> sandbox keeps it off. See [Managed configuration](https://learn.chatgpt.com/docs/enterprise/managed-configuration#configure-network-access-requirements)
> for the administrator-side `requirements.toml` shape.
>
> #### Network policy
>
> Domain rules are allowlist-first:
>
> - Exact hosts match only themselves.
> - `*.example.com` matches subdomains such as `api.example.com`, but not
>   `example.com`.
> - `**.example.com` matches both the apex and subdomains.
> - A global `*` allow rule matches any public host that is not denied. Treat `*`
>   as broad network access and prefer scoped rules when you can.
> - `deny` always wins over `allow`, and global `*` is only valid for allow rules.
>
> #### Local and private destinations
>
> By default, `allow_local_binding = false` blocks loopback, link-local, and
> private destinations:
>
> - Specific exceptions: add an exact local IP literal or `localhost` allow rule
>   when a command needs one local target.
> - Broader access: set `allow_local_binding = true` only when you intentionally
>   want wider local/private reach.
> - Wildcards: wildcard rules do not count as explicit local exceptions.
> - Resolved addresses: hostnames that resolve to local/private IPs stay blocked
>   even if they match the allowlist.
>
> #### DNS rebinding protections
>
> Before allowing a hostname, Codex performs a best-effort DNS and IP
> classification check:
>
> - Lookups that fail or time out are blocked.
> - Hostnames that resolve to non-public addresses are blocked.
> - The check reduces DNS rebinding risk, but it does not eliminate it. Preventing
>   rebinding completely would require pinning resolved IPs through the transport
>   layer.
>
> If hostile DNS is in scope, enforce egress controls at a lower layer too.
>
> #### Dangerous settings
>
> Two settings deliberately widen the trust boundary:
>
> - `dangerously_allow_non_loopback_proxy = true` can expose proxy listeners beyond
>   loopback.
> - `dangerously_allow_all_unix_sockets = true` bypasses the Unix socket allowlist.
>
> Use them only in tightly controlled environments. When Unix socket proxying is
> enabled, listeners stay loopback-only even if non-loopback binding was requested,
> so sandboxed networking does not become a remote bridge into local daemons.
>
> `network_proxy` is off by default. When you enable it:
>
> | Setting                                | Default | Behavior                                                                                                                                                                              |
> | -------------------------------------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `enabled`                              | `false` | Starts sandboxed networking only when command network access is already on.                                                                                                           |
> | `domains`                              | unset   | Uses allowlist behavior, so no external destinations are allowed until you add `allow` rules. Supports exact hosts, scoped wildcards, and global `*` allow rules; `deny` always wins. |
> | `unix_sockets`                         | unset   | No Unix socket destinations are allowed until you add explicit `allow` rules.                                                                                                         |
> | `allow_local_binding`                  | `false` | Blocks local and private-network destinations unless you add an exact local IP literal or `localhost` allow rule, or explicitly opt into broader local/private access.                |
> | `enable_socks5`                        | `true`  | Exposes SOCKS5 support when policy allows it.                                                                                                                                         |
> | `enable_socks5_udp`                    | `true`  | Allows UDP over SOCKS5 when SOCKS5 is available.                                                                                                                                      |
> | `allow_upstream_proxy`                 | `true`  | Lets sandboxed networking honor an upstream proxy from the environment.                                                                                                               |
> | `dangerously_allow_non_loopback_proxy` | `false` | Keeps listener endpoints on loopback unless you deliberately expose them beyond localhost.                                                                                            |
> | `dangerously_allow_all_unix_sockets`   | `false` | Keeps Unix socket access allowlist-based unless you deliberately bypass that protection.                                                                                              |
>
> ### Traffic outside the command network proxy
>
> The network proxy filters scripts, programs, and child processes that run
> inside the local command sandbox. It does not filter web search, app or
> connector tool calls, MCP server connections, browser or Computer Use activity,
> Codex cloud tasks, or the client's model and authentication requests. These
> surfaces use separate service connections, feature settings, workspace
> policies, or environment controls.
>
> For managed users, combine command network policy with controls such as
> `allowed_web_search_modes`, approved `mcp_servers`, and feature requirements
> for apps, plugins, browsers, or Computer Use. See
> [Managed configuration](https://learn.chatgpt.com/docs/enterprise/managed-configuration).
>
> You can also control the [web search tool](https://platform.openai.com/docs/guides/tools-web-search) without granting full network access to spawned commands. Codex defaults to using a web search cache to access results. The cache is an OpenAI-maintained index of web results, so cached mode returns pre-indexed results instead of fetching live pages. This reduces exposure to prompt injection from arbitrary live content, but you should still treat web results as untrusted. If you are using `--yolo` or another [full access sandbox setting](#common-sandbox-and-approval-combinations), web search defaults to live results. Use `--search` or set `web_search = "live"` to allow live browsing, or set it to `"disabled"` to turn the tool off:
>
> ```toml
> web_search = "cached"  # default
> # web_search = "disabled"
> # web_search = "live"  # same as --search
> ```
>
> Set `web_search = "indexed"` when external web access should be gated by the
> search index. Use caution when enabling network access or web search in Codex.
> Prompt injection can cause the agent to fetch and follow untrusted instructions.
>
> ## Defaults and recommendations
>
> - On launch, Codex detects whether the folder is version-controlled and recommends:
>   - Version-controlled folders: `Auto` (workspace write + on-request approvals)
>   - Non-version-controlled folders: `read-only`
> - Depending on your setup, Codex may also start in `read-only` until you explicitly trust the working directory (for example, via an onboarding prompt or `/permissions`).
> - The workspace includes the current directory and temporary directories like `/tmp`. Use the `/status` command to see which directories are in the workspace.
> - To accept the defaults, run `codex`.
> - You can set these explicitly:
>   - `codex --sandbox workspace-write --ask-for-approval on-request`
>   - `codex --sandbox read-only --ask-for-approval on-request`
>
> ### Protected paths in writable roots
>
> In the default `workspace-write` sandbox policy, writable roots still include protected paths:
>
> - `<writable_root>/.git` is protected as read-only whether it appears as a directory or file.
> - If `<writable_root>/.git` is a pointer file (`gitdir: ...`), the resolved Git directory path is also protected as read-only.
> - `<writable_root>/.agents` is protected as read-only when it exists as a directory.
> - `<writable_root>/.codex` is protected as read-only when it exists as a directory.
> - Protection is recursive, so everything under those paths is read-only.
>
> ### Run without approval prompts
>
> You can disable approval prompts with `--ask-for-approval never` or `-a never` (shorthand).
>
> This option works with all `--sandbox` modes, so you still control Codex's level of autonomy. Codex makes a best effort within the constraints you set.
>
> If you need Codex to read files, make edits, and run commands with network access without approval prompts, use `--sandbox danger-full-access` (or the `--dangerously-bypass-approvals-and-sandbox` flag). Use caution before doing so.
>
> For a middle ground, `approval_policy = { granular = { ... } }` lets you keep specific approval prompt categories interactive while automatically rejecting others. The granular policy covers sandbox approvals, execpolicy-rule prompts, MCP prompts, `request_permissions` prompts, and skill-script approvals.
>
> ### Automatic approval reviews
>
> By default, approval requests route to you:
>
> ```toml
> approvals_reviewer = "user"
> ```
>
> Automatic approval reviews apply when approvals are interactive, such as
> `approval_policy = "on-request"` or a granular approval policy. Set
> `approvals_reviewer = "auto_review"` to route eligible approval requests
> through a reviewer agent before Codex runs the request:
>
> ```toml
> approval_policy = "on-request"
> approvals_reviewer = "auto_review"
> ```
>
> For the full reviewer lifecycle, trigger conditions, configuration precedence,
> and failure behavior, see
> [Auto-review](https://learn.chatgpt.com/docs/sandboxing/auto-review).
>
> The reviewer evaluates only actions that already need approval, such as sandbox
> escalations, blocked network requests, `request_permissions` prompts, or
> side-effecting app and MCP tool calls. Actions that stay inside the sandbox
> continue without an extra review step.
>
> The reviewer policy checks for data exfiltration, credential probing, persistent
> security weakening, and destructive actions. Low-risk and medium-risk actions
> can proceed when policy allows them. The policy denies critical-risk actions.
> High-risk actions require enough user authorization and no matching deny rule.
> Prompt-build, review-session, and parse failures fail closed. Timeouts are
> surfaced separately, but the action still does not run.
>
> The [default reviewer policy](https://github.com/openai/codex/blob/main/codex-rs/core/src/guardian/policy.md)
> is in the open-source Codex repository. Enterprises can replace its
> tenant-specific section with `guardian_policy_config` in managed requirements.
> Local `[auto_review].policy` text is also supported, but managed requirements
> take precedence. For setup details, see
> [Managed configuration](https://learn.chatgpt.com/docs/enterprise/managed-configuration#configure-automatic-review-policy).
>
> In the ChatGPT desktop app, these reviews appear as automatic review items with a status
> such as Reviewing, Approved, Denied, Aborted, or Timed out. They can also
> include a risk level and user-authorization assessment for the reviewed
> request.
>
> Automatic review uses extra model calls, so it can add to Codex usage. Admins
> can constrain it with `allowed_approvals_reviewers`.
>
> ### Common sandbox and approval combinations
>
> | Intent                                                            | Flags / config                                                                                                                      | Effect                                                                                                                                           |
> | ----------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
> | Auto (preset)                                                     | _no flags needed_ or `--sandbox workspace-write --ask-for-approval on-request`                                                      | Codex can read files, make edits, and run commands in the workspace. Codex requires approval to edit outside the workspace or to access network. |
> | Safe read-only browsing                                           | `--sandbox read-only --ask-for-approval on-request`                                                                                 | Codex can read files and answer questions. Codex requires approval to make edits, run commands, or access network.                               |
> | Read-only non-interactive (CI)                                    | `--sandbox read-only --ask-for-approval never`                                                                                      | Codex can only read files; never asks for approval.                                                                                              |
> | Automatically edit but ask for approval to run untrusted commands | `--sandbox workspace-write --ask-for-approval untrusted`                                                                            | Codex can read and edit files but asks for approval before running untrusted commands.                                                           |
> | Auto-review mode                                                  | `--sandbox workspace-write --ask-for-approval on-request -c approvals_reviewer=auto_review` or `approvals_reviewer = "auto_review"` | Same sandbox boundary as standard on-request mode, but eligible approval requests are reviewed by Auto-review instead of surfacing to the user.  |
> | Dangerous full access                                             | `--dangerously-bypass-approvals-and-sandbox` (alias: `--yolo`)                                                                      | <ElevatedRiskBadge /> No sandbox; no approvals _(not recommended)_                                                                               |
>
> For non-interactive runs, use `codex exec --sandbox workspace-write`; Codex keeps older `codex exec --full-auto` invocations as a deprecated compatibility path and prints a warning.
>
> With `--ask-for-approval untrusted`, Codex runs only known-safe read operations automatically. Commands that can mutate state or trigger external execution paths (for example, destructive Git operations or Git output/config-override flags) require approval.
>
> #### Configuration in `config.toml`
>
> For the broader configuration workflow, see [Config basics](https://learn.chatgpt.com/docs/config-file/config-basic), [Advanced Config](https://learn.chatgpt.com/docs/config-file/config-advanced#approval-policies-and-sandbox-modes), and the [Configuration Reference](https://learn.chatgpt.com/docs/config-file/config-reference).
>
> ```toml
> # Always ask for approval mode
> approval_policy = "untrusted"
> sandbox_mode    = "read-only"
> allow_login_shell = false # optional hardening: disallow login shells for shell-based tools
>
> # Optional: Allow network in workspace-write mode
> [sandbox_workspace_write]
> network_access = true
>
> # Optional: granular approval policy
> # approval_policy = { granular = {
> #   sandbox_approval = true,
> #   rules = true,
> #   mcp_elicitations = true,
> #   request_permissions = false,
> #   skill_approval = false
> # } }
> ```
>
> You can also save presets as [profile files](https://learn.chatgpt.com/docs/config-file/config-advanced#profiles), then select them with `codex --profile profile-name`:
>
> ```toml
> # ~/.codex/full_auto.config.toml
> approval_policy = "on-request"
> sandbox_mode    = "workspace-write"
> ```
>
> ```toml
> # ~/.codex/readonly_quiet.config.toml
> approval_policy = "never"
> sandbox_mode    = "read-only"
> ```
>
> ### Test the sandbox locally
>
> To see what happens when a command runs under the Codex sandbox, use these Codex CLI commands:
>
> ```bash
> # macOS
> codex sandbox macos [--permissions-profile <name>] [--log-denials] [COMMAND]...
> # Linux
> codex sandbox linux [--permissions-profile <name>] [COMMAND]...
> # Windows
> codex sandbox windows [--permissions-profile <name>] [COMMAND]...
> ```
>
> The `sandbox` command is also available as `codex debug`, and the platform helpers have aliases (for example `codex sandbox seatbelt` and `codex sandbox landlock`).
>
> ## OS-level sandbox
>
> Codex enforces the sandbox differently depending on your OS:
>
> - **macOS** uses Seatbelt policies and runs commands using `sandbox-exec` with a profile (`-p`) that corresponds to the `--sandbox` mode you selected. When restricted read access enables platform defaults, Codex appends a curated macOS platform policy (instead of broadly allowing `/System`) to preserve common tool compatibility.
> - **Linux** uses `bwrap` plus `seccomp` by default.
> - **Windows** uses the Linux sandbox implementation when running in [Windows Subsystem for Linux 2 (WSL2)](https://learn.chatgpt.com/docs/windows/wsl). WSL1 was supported through Codex `0.114`; starting in `0.115`, the Linux sandbox moved to `bwrap`, so WSL1 is no longer supported. When running natively on Windows, Codex uses a [Windows sandbox](https://learn.chatgpt.com/docs/windows/windows-sandbox#windows-sandbox) implementation.
>
> If you use the Codex IDE extension on Windows, it supports WSL2 directly. Set the following in your VS Code settings to keep the agent inside WSL2 whenever it's available:
>
> ```json
> {
>   "chatgpt.runCodexInWindowsSubsystemForLinux": true
> }
> ```
>
> This ensures the IDE extension inherits Linux sandbox semantics for commands, approvals, and filesystem access even when the host OS is Windows. Learn more in the [WSL guide](https://learn.chatgpt.com/docs/windows/wsl).
>
> When running natively on Windows, configure the native sandbox mode in `config.toml`:
>
> ```toml
> [windows]
> sandbox = "unelevated" # or "elevated"
> # sandbox_private_desktop = true  # default; set false only for compatibility
> ```
>
> See the [Windows setup guide](https://learn.chatgpt.com/docs/windows/windows-sandbox#windows-sandbox) for details.
>
> When you run Linux in a containerized environment such as Docker, the sandbox may not work if the host or container configuration blocks the namespace, setuid `bwrap`, or `seccomp` operations that Codex needs.
>
> In that case, configure your Docker container to provide the isolation you need, then run `codex` with `--sandbox danger-full-access` (or the `--dangerously-bypass-approvals-and-sandbox` flag) inside the container.
>
> ### Run Codex in Dev Containers
>
> If your host cannot run the Linux sandbox directly, or if your organization already standardizes on containerized development, run Codex with Dev Containers and let Docker provide the outer isolation boundary. This works with Visual Studio Code Dev Containers and compatible tools.
>
> Use the [Codex secure devcontainer example](https://github.com/openai/codex/tree/main/.devcontainer) as a reference implementation. The example installs Codex, common development tools, `bubblewrap`, and firewall-based outbound controls.
>
> Devcontainers provide substantial protection, but they do not prevent every
>   attack. If you run Codex with `--sandbox danger-full-access` or
>   `--dangerously-bypass-approvals-and-sandbox` inside the container, a malicious
>   project can exfiltrate anything available inside the devcontainer, including
>   Codex credentials. Use this pattern only with trusted repositories, and
>   monitor Codex activity as you would in any other elevated environment.
>
> The reference implementation includes:
>
> - an Ubuntu 24.04 base image with Codex and common development tools installed;
> - an allowlist-driven firewall profile for outbound access;
> - VS Code settings and extension recommendations for reopening the workspace in a container;
> - persistent mounts for command history and Codex configuration;
> - `bubblewrap`, so Codex can still use its Linux sandbox when the container grants the needed capabilities.
>
> To try it:
>
> 1. Install Visual Studio Code and the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).
> 2. Copy the Codex example `.devcontainer` setup into your repository, or start from the Codex repository directly.
> 3. In VS Code, run **Dev Containers: Open Folder in Container...** and select `.devcontainer/devcontainer.secure.json`.
> 4. After the container starts, open a terminal and run `codex`.
>
> You can also start the container from the CLI:
>
> ```bash
> devcontainer up --workspace-folder . --config .devcontainer/devcontainer.secure.json
> ```
>
> The example has three main pieces:
>
> - `.devcontainer/devcontainer.secure.json` controls container settings, capabilities, mounts, environment variables, and VS Code extensions.
> - `.devcontainer/Dockerfile.secure` defines the Ubuntu-based image and installed tools.
> - `.devcontainer/init-firewall.sh` applies the outbound network policy.
>
> The reference firewall is intentionally a starting point. If you depend on domain allowlisting for isolation, implement DNS rebinding and DNS refresh protections that fit your environment, such as TTL-aware refreshes or a DNS-aware firewall.
>
> Inside the container, choose one of these modes:
>
> - Keep Codex's Linux sandbox enabled if the Dev Container profile grants the capabilities needed for `bwrap` to create the inner sandbox.
> - If the container is your intended security boundary, run Codex with `--sandbox danger-full-access` inside the container so Codex does not try to create a second sandbox layer.
>
> ## Version control
>
> Codex works best with a version control workflow:
>
> - Work on a feature branch and keep `git status` clean before delegating. This keeps Codex patches easier to isolate and revert.
> - Prefer patch-based workflows (for example, `git diff`/`git apply`) over editing tracked files directly. Commit frequently so you can roll back in small increments.
> - Treat Codex suggestions like any other PR: run targeted verification, review diffs, and document decisions in commit messages for auditing.
>
> ## Monitoring and telemetry
>
> Codex supports opt-in monitoring via OpenTelemetry (OTel) to help teams audit usage, investigate issues, and meet compliance requirements without weakening local security defaults. Telemetry is off by default; enable it explicitly in your configuration.
>
> ### Overview
>
> - Codex turns off OTel export by default to keep local runs self-contained.
> - When enabled, Codex emits structured log events covering chats, API requests, SSE/WebSocket stream activity, user prompts (redacted by default), tool approval decisions, and tool results.
> - Codex tags exported events with `service.name` (originator), CLI version, and an environment label to separate dev/staging/prod traffic.
>
> ### Enable OTel (opt-in)
>
> Add an `[otel]` block to your Codex configuration (typically `~/.codex/config.toml`), choosing an exporter and whether to log prompt text.
>
> ```toml
> [otel]
> environment = "staging"   # dev | staging | prod
> exporter = "none"          # none | otlp-http | otlp-grpc
> log_user_prompt = false     # redact prompt text unless policy allows
> ```
>
> - `exporter = "none"` leaves instrumentation active but doesn't send data anywhere.
> - To send events to your own collector, pick one of:
>
> ```toml
> [otel]
> exporter = { otlp-http = {
>   endpoint = "https://otel.example.com/v1/logs",
>   protocol = "binary",
>   headers = { "x-otlp-api-key" = "${OTLP_TOKEN}" }
> }}
> ```
>
> ```toml
> [otel]
> exporter = { otlp-grpc = {
>   endpoint = "https://otel.example.com:4317",
>   headers = { "x-otlp-meta" = "abc123" }
> }}
> ```
>
> Codex batches events and flushes them on shutdown. Codex exports only telemetry produced by its OTel module.
>
> ### Event categories
>
> Representative event types include:
>
> - `codex.conversation_starts` (model, reasoning settings, sandbox/approval policy)
> - `codex.api_request` (attempt, status/success, duration, and error details)
> - `codex.sse_event` (stream event kind, success/failure, duration, plus token counts on `response.completed`)
> - `codex.websocket_request` and `codex.websocket_event` (request duration plus per-message kind/success/error)
> - `codex.user_prompt` (length; content redacted unless explicitly enabled)
> - `codex.tool_decision` (approved/denied, source: configuration vs. user)
> - `codex.tool_result` (duration, success, output snippet)
>
> Associated OTel metrics (counter plus duration histogram pairs) include `codex.api_request`, `codex.sse_event`, `codex.websocket.request`, `codex.websocket.event`, and `codex.tool.call` (with corresponding `.duration_ms` instruments).
>
> For the full event catalog and configuration reference, see the [Codex configuration documentation on GitHub](https://github.com/openai/codex/blob/main/docs/config.md#otel).
>
> ### Security and privacy guidance
>
> - Keep `log_user_prompt = false` unless policy explicitly permits storing prompt contents. Prompts can include source code and sensitive data.
> - Route telemetry only to collectors you control; apply retention limits and access controls aligned with your compliance requirements.
> - Treat tool arguments and outputs as sensitive. Favor redaction at the collector or SIEM when possible.
> - Review local data retention settings (for example, `history.persistence` / `history.max_bytes`) if you don't want Codex to save session transcripts under `CODEX_HOME`. See [Advanced Config](https://learn.chatgpt.com/docs/config-file/config-advanced#history-persistence) and [Configuration Reference](https://learn.chatgpt.com/docs/config-file/config-reference).
> - If you run the CLI with network access turned off, OTel export can't reach your collector. To export, allow network access in `workspace-write` mode for the OTel endpoint, or export from Codex cloud with the collector domain on your approved list.
> - Review events periodically for approval/sandbox changes and unexpected tool executions.
>
> OTel is optional and designed to complement, not replace, the sandbox and approval protections described above.
>
> ## Managed configuration
>
> Enterprise admins can configure Codex security settings for their workspace in [Managed configuration](https://learn.chatgpt.com/docs/enterprise/managed-configuration). See that page for setup and policy details.

### Source: Windows sandbox — Configure the Windows sandbox

> ## Configure the Windows sandbox
>
> When you run Codex natively on Windows, agent mode uses a Windows sandbox to
> block filesystem writes outside the working folder and prevent network access
> without your explicit approval.
>
> Native Windows sandbox support includes two modes that you can configure in
> `config.toml`:
>
> ```toml
> [windows]
> sandbox = "elevated" # or "unelevated"
> ```
>
> `elevated` is the preferred native Windows sandbox. It uses dedicated
> lower-privilege sandbox users, filesystem permission boundaries, firewall
> rules, and local policy changes needed for commands that run in the sandbox.
>
> `unelevated` is the fallback native Windows sandbox. It runs commands with a
> restricted Windows token derived from your current user, applies ACL-based
> filesystem boundaries, and uses environment-level offline controls instead of
> the dedicated offline-user firewall rule. It's weaker than `elevated`, but it
> is still useful when administrator-approved setup is blocked by local or
> enterprise policy.
>
> If both modes are available, use `elevated`. If the default native sandbox
> doesn't work in your environment, use `unelevated` as a fallback while you
> troubleshoot the setup.
>
> Enterprise administrators can constrain which native sandbox implementations
> Codex can use through [`requirements.toml`](https://learn.chatgpt.com/docs/enterprise/managed-configuration#admin-enforced-requirements-requirementstoml):
>
> ```toml
> [windows]
> allowed_sandbox_implementations = ["elevated"]
> ```
>
> This example requires the `elevated` sandbox and prevents users from falling
> back to `unelevated`. To permit either implementation, include both values;
> Codex prefers `elevated` when no mode is selected. See the
> [`requirements.toml` reference](https://learn.chatgpt.com/docs/config-file/config-reference#requirementstoml) for
> the supported values.
>
> By default, both sandbox modes also use a private desktop for stronger UI
> isolation. Set `windows.sandbox_private_desktop = false` only if you need the
> older `Winsta0\\Default` behavior for compatibility.

### Source: Windows sandbox — Grant sandbox read access; diagnostics paths

> ### Grant sandbox read access
>
> When a command fails because the Windows sandbox can't read a directory, use:
>
> ```text
> /sandbox-add-read-dir C:\absolute\directory\path
> ```
>
> The path must be an existing absolute directory. After the command succeeds, later commands that run in the sandbox can read that directory during the current session.

> If you still have problems, send:
>
> - `CODEX_HOME/.sandbox/sandbox.log`
>
> Do not send:
>
> - the contents of `CODEX_HOME/.sandbox-secrets/`

### Source: WSL — platform behavior

> When you use WSL2, Codex runs inside the Linux environment instead of using the
> native [Windows sandbox](https://learn.chatgpt.com/docs/windows/windows-sandbox). Choose WSL2 when you need Linux-native
> tooling, your repositories and developer workflow already live in WSL2, or
> neither native Windows sandbox mode works for your environment.
>
> WSL1 was supported through Codex `0.114`. Starting in Codex `0.115`, the Linux
> sandbox moved to `bubblewrap`, so WSL1 is no longer supported.
