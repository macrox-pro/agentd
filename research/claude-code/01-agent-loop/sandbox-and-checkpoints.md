---
primary_sources:
  - id: T1-SANDBOX
    title: "Sandboxing"
    url: "https://code.claude.com/docs/en/sandboxing.md"
    section: ""
  - id: T1-SANDBOX-ENV
    title: "Sandbox environments"
    url: "https://code.claude.com/docs/en/sandbox-environments.md"
    section: ""
  - id: T1-CHECKPOINT
    title: "Checkpointing"
    url: "https://code.claude.com/docs/en/checkpointing.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Sandbox and checkpoints

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Sandboxing
>
> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Configure the sandboxed Bash tool
>
> > Learn how Claude Code's sandboxed Bash tool provides filesystem and network isolation for safer, more autonomous agent execution.
>
> The Bash sandbox lets Claude run most shell commands without stopping to ask permission. Instead of approving each command, you define which files and network domains commands can touch, and the operating system enforces that boundary for every Bash command and its child processes.
>
>   To compare other isolation approaches such as dev containers, custom containers, and virtual machines, see [Sandbox environments](/docs/en/sandbox-environments). To reduce permission prompts for tools other than Bash, see [permission modes](/docs/en/permission-modes).
>
> ## Get started
>
> The sandbox is built into Claude Code and runs on macOS, Linux, and WSL2. Native Windows is not supported. On Windows, run Claude Code inside a WSL2 distribution.
>
> On macOS, there is nothing to install: sandboxing uses the built-in Seatbelt framework. On Linux and WSL2, the sandbox relies on two packages, covered in [Set up Linux and WSL2](#set-up-linux-and-wsl2). Even if you haven't installed them yet, you can start with `/sandbox`, because its panel shows whether anything is missing.
>
>   **Run /sandbox**: Start a Claude Code session and run the `/sandbox` command:
>
>     ```text
>     /sandbox
>     ```
>
>     This opens the sandbox panel with three tabs, plus a Dependencies tab on Linux when the optional seccomp filter is missing:
>
>     * **Mode**: choose how sandboxed commands are approved, covered in the next step
>     * **Overrides**: choose whether commands that fail under the sandbox can fall back to running unsandboxed. This is the [`allowUnsandboxedCommands`](/docs/en/settings-reference#sandbox-allowunsandboxedcommands) setting
>     * **Config**: view the resolved sandbox settings
>
>     If the panel shows only a Dependencies tab, a required package is missing. Install it as described in [Set up Linux and WSL2](#set-up-linux-and-wsl2), restart Claude Code, and run `/sandbox` again.
>
>   **Choose a mode**: On the Mode tab, select auto-allow or regular permissions. Auto-allow runs sandboxed commands without prompting, and regular permissions keeps the regular permission prompts even when commands are sandboxed. See [Sandbox modes](#sandbox-modes) for which commands still prompt in auto-allow mode.
>
>   **Run a Bash command**: Ask Claude to run a command, such as a build or a test suite. By default, commands inside the sandbox can write only to the working directory and the session temp directory. The first time a command needs a new network domain, Claude Code prompts for approval, or in [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) sends the request to the classifier.
>
>     Commands that cannot run sandboxed fall back to the regular permission flow. To widen or narrow these boundaries, see [Configure sandboxing](#configure-sandboxing).
>
> When you select a mode in the panel, Claude Code saves it to your project's local settings at `.claude/settings.local.json`, which apply to the current project. Claude Code adds that file to your global gitignore when it saves a setting there. To enable the sandbox across all of your projects, set [`sandbox.enabled`](/docs/en/settings-reference#sandbox-enabled) to `true` in your user settings at `~/.claude/settings.json`. To enforce sandboxing for every developer in an organization, use [managed settings](#enforce-sandboxing-with-managed-settings).
>
>   By default, if the sandbox cannot start because dependencies are missing or the platform is unsupported, Claude Code shows a warning and runs commands without sandboxing. To make this a hard failure instead, set [`sandbox.failIfUnavailable`](/docs/en/settings-reference#sandbox-failifunavailable) to `true`. This is intended for managed deployments that require sandboxing as a security gate.
>
> ### Set up Linux and WSL2
>
> On Linux and WSL2, the sandbox relies on two packages:
>
> * [`bubblewrap`](https://github.com/containers/bubblewrap): the unprivileged sandboxing tool that enforces filesystem isolation
> * [`socat`](http://www.dest-unreach.org/socat/): the relay used to route network traffic through the sandbox proxy
>
> Install them with your distribution's package manager:
>
>   #### Ubuntu/Debian
>
> ```bash
>     sudo apt-get install bubblewrap socat
>     ```
>
>   #### Fedora
>
> ```bash
>     sudo dnf install bubblewrap socat
>     ```
>
> When a dependency is missing, the Dependencies tab in `/sandbox` lists which of `ripgrep`, `bubblewrap`, `socat`, and the seccomp filter your platform lacks. If you don't see the tab after installing and restarting Claude Code, all dependencies are present.
>
> Ripgrep is bundled with the native Claude Code binary. The seccomp filter is optional and adds Unix domain socket blocking. Install it with `npm install -g @anthropic-ai/sandbox-runtime` if it is missing.
>
> When a required dependency is missing, the Dependencies tab is the only tab shown until you install it. When only the optional seccomp filter is missing, the Dependencies tab appears alongside the other tabs. The dependency check runs at startup, so restart Claude Code after installing packages for `/sandbox` to detect them.
>
>
>     On Ubuntu 24.04 and later, the default AppArmor policy prevents bubblewrap from creating the user namespaces it needs for isolation.
>
>     To check whether your environment enforces this restriction, including inside WSL2, run `sysctl kernel.apparmor_restrict_unprivileged_userns`. If the command returns `0`, skip this step. If it prints a `No such file or directory` error, the key doesn't exist and you can skip this step. If it returns `1`, add an AppArmor profile that grants `bwrap` this capability:
>
>     ```bash
>     sudo tee /etc/apparmor.d/bwrap > /dev/null <<'EOF'
>     abi <abi/4.0>,
>     include <tunables/global>
>
>     profile bwrap /usr/bin/bwrap flags=(unconfined) {
>       userns,
>       include if exists <local/bwrap>
>     }
>     EOF
>     ```
>
>     The profile applies only to `bwrap` itself, not to the commands it runs inside the sandbox. Reload AppArmor to apply it:
>
>     ```bash
>     sudo systemctl reload apparmor
>     ```
>
>
>
>     Check your WSL version with `wsl -l -v` from PowerShell. If you see `Sandboxing requires WSL2`, your distribution is running WSL1. Upgrade it to WSL2 or run Claude Code without sandboxing.
>
>     On WSL2, WSL hands a launch of a Windows binary such as `cmd.exe`, `powershell.exe`, or anything under `/mnt/c/` to the Windows host over a Unix socket, so whether a sandboxed command can launch one follows the sandbox's [Unix-socket settings](/docs/en/settings-reference#sandbox-network-allowunixsockets): the optional seccomp filter has to be installed to block the socket in the first place. To allow these launches, set `allowAllUnixSockets`; to keep them out of the sandbox entirely, add the command to [`excludedCommands`](/docs/en/settings-reference#sandbox-excludedcommands).
>
>
> ### Sandbox modes
>
> Claude Code offers two sandbox modes. In both, the sandbox enforces the same filesystem and network restrictions; the difference is only in whether sandboxed commands are auto-approved or require explicit permission.
>
> #### Auto-allow mode
>
> When a command can be sandboxed, Claude Code runs it inside the sandbox and approves it automatically, without asking your permission. Commands that cannot be sandboxed, such as those needing network access to non-allowed hosts, fall back to the regular permission flow, where Claude Code checks your [permission rules](/docs/en/permissions) and gates any command those rules do not already allow, with a prompt in Manual mode or the classifier in [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode).
>
> Even in auto-allow mode, the following still apply:
>
> * Explicit [deny rules](/docs/en/permissions) are always respected
> * `rm` or `rmdir` commands that target a [critical path](/docs/en/permission-modes#critical-paths) still go through the regular permission flow
> * Content-scoped [ask rules](/docs/en/permissions) like `Bash(git push *)` still force a prompt even for sandboxed commands
> * A bare `Bash` ask rule, or the equivalent `Bash(*)` form, is skipped for commands that run sandboxed; it still applies to commands that fall back to the regular permission flow. In [plan mode](/docs/en/permission-modes#analyze-before-you-edit-with-plan-mode), the rule isn't skipped: it prompts for sandboxed commands too, including read-only ones. Before v2.1.212, the skip applied in plan mode as well
>
>   Auto-allow mode works independently of your permission mode setting, with one exception: [plan mode](/docs/en/permission-modes#analyze-before-you-edit-with-plan-mode). Even if you're not in "accept edits" mode, sandboxed Bash commands run automatically when auto-allow is enabled. This means Bash commands that modify files within the sandbox boundaries execute without prompting, even in Manual mode, where the file edit tools would prompt.
>
>   In plan mode, auto-allow doesn't widen approvals; see [plan mode](/docs/en/permission-modes#analyze-before-you-edit-with-plan-mode) for how Claude Code gates commands while you plan. Before v2.1.212, auto-allow ran sandboxed commands without a prompt in plan mode too.
>
> #### Regular permissions mode
>
> All Bash commands go through the regular permission flow, even when sandboxed. This provides more control but requires more approvals.
>
> #### The unsandboxed retry escape hatch
>
> Some commands can't run inside the sandbox at all, such as tools that are incompatible with it or that need a host you haven't allowed. When a command fails after the sandbox denied it access, Claude Code appends the violation details to the failed command's output, so Claude sees which file path or network host the sandbox blocked. Rather than failing the task or requiring you to turn sandboxing off, Claude Code includes an escape hatch: when a command fails because of sandbox restrictions, Claude analyzes the failure and may retry the command with the `dangerouslyDisableSandbox` parameter.
>
> The retried command runs outside the sandbox, so it goes through the regular permission flow: in Manual mode you get a confirmation prompt; in [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) the classifier evaluates the underlying command instead of prompting you. To be prompted on every unsandboxed retry even in auto mode, add an [ask rule](/docs/en/permissions#match-by-input-parameter) for `Bash(dangerouslyDisableSandbox:true)`.
>
> You can disable this escape hatch by setting `"allowUnsandboxedCommands": false` in your [sandbox settings](/docs/en/settings-reference#sandbox-settings). When disabled, which the `/sandbox` Overrides tab shows as **Strict sandbox mode**, the `dangerouslyDisableSandbox` parameter is completely ignored and all commands must run sandboxed or be explicitly listed in `excludedCommands`.
>
> #### Temporary directories
>
> The session temp directory is writable inside the sandbox by default, alongside the working directory. Unless you [disable filesystem isolation](#disable-filesystem-isolation), Claude Code sets `$TMPDIR` to this directory for sandboxed commands, so tools that write temporary files work without extra configuration. Unsandboxed commands inherit your shell's `$TMPDIR` unchanged, so while filesystem isolation is on, sandboxed and unsandboxed commands resolve `$TMPDIR` to different directories. To pass temporary files between the two, write them under the working directory instead.
>
> ## Configure sandboxing
>
> Customize sandbox behavior through your `settings.json` file. See [Settings](/docs/en/settings-reference#sandbox-settings) for the complete configuration reference.
>
> By default, sandboxed commands can write only to the current working directory and the session temp directory. If subprocess commands like `kubectl`, `terraform`, or `npm` need to write outside those directories, use `sandbox.filesystem.allowWrite` to grant access to specific paths:
>
> ```json
> {
>   "sandbox": {
>     "enabled": true,
>     "filesystem": {
>       "allowWrite": ["~/.kube", "/tmp/build"]
>     }
>   }
> }
> ```
>
> These paths are enforced at the OS level, so all commands running inside the sandbox, including their child processes, respect them. This is the recommended approach when a tool needs write access to a specific location, rather than excluding the tool from the sandbox entirely with `excludedCommands`.
>
> When the same filesystem array is defined in multiple [settings scopes](/docs/en/settings#settings-precedence), the arrays are merged: paths from every scope are combined, not replaced. When you edit these lists during a session, Claude Code [applies the change to the running session](/docs/en/settings#when-edits-take-effect), so the next sandboxed command runs under the new paths.
>
> Path prefixes control how paths are resolved:
>
> | Prefix            | Meaning                                                                                | Example                                                                   |
> | :---------------- | :------------------------------------------------------------------------------------- | :------------------------------------------------------------------------ |
> | `/`               | Absolute path from filesystem root                                                     | `/tmp/build` stays `/tmp/build`                                           |
> | `~/`              | Relative to home directory                                                             | `~/.kube` becomes `$HOME/.kube`                                           |
> | `./` or no prefix | Relative to the project root for project settings, or to `~/.claude` for user settings | `./output` in `.claude/settings.json` resolves to `<project-root>/output` |
>
> This syntax differs from [Read and Edit permission rules](/docs/en/permissions#read-and-edit), which use `//path` for absolute and `/path` for project-relative. Sandbox filesystem paths use standard conventions: `/tmp/build` is absolute. For how Claude Code treats a trailing slash or a wildcard in these paths, see [Sandbox path prefixes](/docs/en/settings-reference#sandbox-path-prefixes).
>
> You can also deny write or read access using `sandbox.filesystem.denyWrite` and `sandbox.filesystem.denyRead`, and re-allow specific paths within a denied region using `sandbox.filesystem.allowRead`. When read rules overlap, the more specific path wins:
>
> | Example rules                                           | Result                                                                                                                                                                                                   |
> | :------------------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `"denyRead": ["~/"]` with `"allowRead": ["~/projects"]` | `~/projects` is readable and the rest of the home directory stays blocked. The narrower allow re-opens that part of the denied region                                                                    |
> | `"allowRead": ["~/"]` with `"denyRead": ["~/.env"]`     | `~/.env` stays blocked and the rest of the home directory is readable. The deny holds inside a wider allow, so a broad allow can't silently re-expose a secret                                           |
> | `"allowRead": ["~/"]` with `"denyRead": ["~/**/.env"]`  | Every `.env` under the home directory stays blocked and the rest is readable. A [wildcard deny](/docs/en/settings-reference#sandbox-path-prefixes) holds inside a wider allow the same way an exact path does |
>
> The example below blocks reading from the entire home directory while still allowing reads from the current project. Place it in your project's `.claude/settings.json`, because the relative path `.` resolves to the project root only when the configuration lives in project settings:
>
> ```json
> {
>   "sandbox": {
>     "enabled": true,
>     "filesystem": {
>       "denyRead": ["~/"],
>       "allowRead": ["."]
>     }
>   }
> }
> ```
>
> If you placed the same configuration in `~/.claude/settings.json`, `.` would resolve to `~/.claude` instead, and project files would remain blocked by the `denyRead` rule.
>
> ### Disable filesystem isolation
>
> Set `sandbox.filesystem.disabled` to `true` to skip filesystem isolation while keeping network isolation. The example below turns off filesystem isolation while keeping an allowlist of network domains:
>
> ```json
> {
>   "sandbox": {
>     "enabled": true,
>     "filesystem": {
>       "disabled": true
>     },
>     "network": {
>       "allowedDomains": ["github.com", "*.npmjs.org"]
>     }
>   }
> }
> ```
>
> The sandbox has two independent layers: [filesystem isolation](#filesystem-isolation) controls which paths sandboxed commands can read and write, and [network isolation](#network-isolation) controls which domains they can reach. With the filesystem layer off, sandboxed commands get unrestricted read and write access to the host filesystem, while their network egress stays confined to your allowed domains. Turn the layer off when you sandbox to control where commands connect rather than what they write.
>
> The setting is off by default and applies on the platforms where the sandbox runs: macOS, Linux, and WSL2. Requires Claude Code v2.1.216 or later.
>
>   With filesystem isolation off and commands auto-allowed, a sandboxed command can write files that later commands run or read, such as shell startup files, executables on `$PATH`, or `~/.claude/settings.json`, and use them to widen its own access on the next run. Set `filesystem.disabled` to `true` only for workloads you trust not to escalate their own access. Locking network domains with [`allowManagedDomainsOnly`](#keep-developers-from-widening-the-policy) narrows the risk but doesn't remove it, since that lock applies only to commands running inside the sandbox.
>
> #### Which settings can disable it
>
> Because turning filesystem isolation off widens what sandboxed commands can do, Claude Code honors `filesystem.disabled` from these settings sources only:
>
> * User settings, managed settings, and the `--settings` CLI flag can set it. Project settings in `.claude/settings.json` and `.claude/settings.local.json` can't, so a checked-out project can't switch filesystem isolation off.
> * When managed settings configure `sandbox.filesystem` at all, or list any `sandbox.credentials.files` entry with `"mode": "deny"`, only managed settings can set the key. This keeps administrator-deployed filesystem restrictions in force; to relax such a deployment, set `"disabled": true` in managed settings.
> * When [`CLAUDE_CODE_SUBPROCESS_ENV_SCRUB`](/docs/en/env-vars) is set, Claude Code ignores `filesystem.disabled` from every source, including managed settings, and keeps filesystem isolation on.
>
> Whether a managed `credentials.files` entry pins `filesystem.disabled`, locking the key to managed settings so developers can't turn filesystem isolation off, depends on the entry's `mode` and what happens to the entry when the sandbox starts:
>
> | Managed entry                                                                                                  | Pins `filesystem.disabled`   | What protects the file when isolation is off                                                                                    |
> | -------------------------------------------------------------------------------------------------------------- | ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
> | `"mode": "deny"`                                                                                               | Yes                          | Nothing: the read block is part of the filesystem layer                                                                         |
> | `"mode": "mask"`, applied as a mask                                                                            | No                           | Masking itself: the [sentinel copy and proxy](#mask-credential-files) on Linux and WSL2, the sandbox's own read rules on macOS  |
> | `"mode": "mask"`, [fallen back to `deny`](#mask-credential-files) at setup                                     | No                           | Nothing, same as `deny`. List a path that can't be masked, such as a directory, as an explicit `deny` entry, which pins the key |
> | `"mode": "mask"`, [degraded to `deny` by validation](/docs/en/managed-settings#invalid-entries-in-managed-settings) | Yes, like an explicit `deny` | Nothing, same as `deny`                                                                                                         |
>
> A fallback happens when the sandbox starts, after Claude Code has already read the settings the pin check runs on, so a fallen-back entry never pins. Validation rewrites an invalid entry to `deny` while settings load, so a degraded entry pins like one you wrote as `deny`.
>
> #### What changes when filesystem isolation is off
>
> Setting `filesystem.disabled` lifts the protections the filesystem layer itself enforces. Protections that other layers enforce keep applying:
>
> | Protection                                                                               | With filesystem isolation off                                                                                                                                |
> | ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `filesystem.denyRead` and [`credentials.files`](#protect-credentials) `deny` read blocks | Not enforced. The filesystem layer applies both                                                                                                              |
> | `credentials.envVars` `deny` and `mask` entries                                          | Enforced. Environment variable scrubbing is independent of the filesystem layer                                                                              |
> | [`credentials.files` `mask` entries](#mask-credential-files) applied as masks            | Enforced: masking is independent of the filesystem layer. An entry that [fell back to `deny`](#mask-credential-files) is not enforced, like any `deny` entry |
>
> Two other things change:
>
> * Sandboxed commands inherit your shell's `$TMPDIR` instead of the session temp directory, because every temp directory is writable and Claude Code no longer redirects commands to the session one. On Linux the variable is often unset in the parent shell, so it can expand empty inside sandboxed commands; Claude Code tells Claude through its Bash tool guidance to create scratch directories with `mktemp -d` instead of relying on `$TMPDIR`.
> * [`autoAllowBashIfSandboxed`](/docs/en/settings-reference#sandbox-autoallowbashifsandboxed) still defaults to `true`, so sandboxed commands keep running without prompts. Set it to `false` to prompt for sandboxed commands.
>
> ### Protect credentials
>
> The `sandbox.credentials` setting declares credential files and environment variables to protect from sandboxed commands. Each entry names a file path or an environment variable and a `mode`. The dedicated `credentials` block keeps credential rules grouped together and separate from general filesystem rules. Requires Claude Code v2.1.187 or later.
>
> For entries with `"mode": "deny"`, file paths are denied for reads inside the sandbox, the same restriction that `filesystem.denyRead` applies, and environment variables are unset before each sandboxed command runs. The file protection is part of the filesystem layer, so it doesn't apply if you [disable filesystem isolation](#disable-filesystem-isolation); the environment variable protection still does.
>
> The example below blocks reads of the AWS credentials file and the SSH directory and removes `GITHUB_TOKEN` and `NPM_TOKEN` from the environment of sandboxed commands:
>
> ```json
> {
>   "sandbox": {
>     "enabled": true,
>     "credentials": {
>       "files": [
>         { "path": "~/.aws/credentials", "mode": "deny" },
>         { "path": "~/.ssh", "mode": "deny" }
>       ],
>       "envVars": [
>         { "name": "GITHUB_TOKEN", "mode": "deny" },
>         { "name": "NPM_TOKEN", "mode": "deny" }
>       ]
>     }
>   }
> }
> ```
>
> Environment variable entries and file entries also accept `"mode": "mask"`, described under [Mask credentials](#mask-credentials).
>
> File paths follow the same [prefix rules](/docs/en/settings-reference#sandbox-path-prefixes) as `sandbox.filesystem.*` settings, and `deny` entries from every [settings scope](/docs/en/settings#settings-precedence) are merged. A `deny` entry only ever narrows access, so any scope can add one, but no scope can remove one that another scope added.
>
> There is no built-in credential deny list, so only the files and variables you list are restricted. The setting affects sandboxed Bash commands only. To strip Anthropic and cloud provider credentials from all subprocesses regardless of sandboxing, set [`CLAUDE_CODE_SUBPROCESS_ENV_SCRUB`](/docs/en/env-vars).
>
> ### Mask credentials
>
> Masking goes further than a `deny` entry under [Protect credentials](#protect-credentials). Instead of blocking a credential, Claude Code shows sandboxed commands a placeholder, the sentinel, and the [sandbox proxy](#network-isolation) swaps in the real value on outbound requests to hosts you allow. For files, the substitution is Linux and WSL2 behavior; [macOS blocks the file instead](#mask-credential-files).
>
> #### Mask environment variables
>
> `"mode": "mask"` protects a credential while keeping the tools that authenticate with it working. `deny` removes the variable entirely, which also breaks tools that need it, such as `gh` or `npm`. Requires Claude Code v2.1.199 or later.
>
> With `mask`, the sandboxed command sees a per-session sentinel value instead of the real one. Each `mask` entry can list `injectHosts`, the hosts the real value is allowed to reach. When a request leaves the sandbox for one of them, the [sandbox proxy](#network-isolation) replaces the sentinel with the real value. The command and anything it logs never hold the real credential, but its requests still authenticate.
>
> The proxy substitutes the credential inside request contents, so it has to see them. Set [`network.tlsTerminate`](/docs/en/settings-reference#sandbox-network-tlsterminate) so the proxy terminates TLS itself. Without it, masking fails without exposing anything: the command still sees only the sentinel, but the sentinel reaches the server unchanged and authentication fails. Claude Code reports this misconfiguration at startup.
>
> Substitution covers headers and request bodies. Requests that authenticate with a signature derived from the credential, rather than the credential itself, need re-signing at the proxy; [Re-sign AWS requests](#re-sign-aws-requests) covers how that works for AWS.
>
> The example below masks two tokens. `GH_TOKEN` is substituted only on requests to `api.github.com`, while `NPM_TOKEN` has no `injectHosts` and is substituted on requests to every host in `network.allowedDomains`. The proxy injects only on connections the [domain allowlist](#network-isolation) admits, so each `injectHosts` destination must also be reachable through `network.allowedDomains`.
>
> ```json
> {
>   "sandbox": {
>     "enabled": true,
>     "network": {
>       "tlsTerminate": {},
>       "allowedDomains": ["*.github.com", "registry.npmjs.org"]
>     },
>     "credentials": {
>       "envVars": [
>         { "name": "GH_TOKEN", "mode": "mask", "injectHosts": ["api.github.com"] },
>         { "name": "NPM_TOKEN", "mode": "mask" }
>       ]
>     }
>   }
> }
> ```
>
> * **`network.allowedDomains`**: the [bracketed form domain lists use](#ipv6-addresses-in-domain-lists), such as `"[::1]"`. The proxy checks this list to admit the connection.
> * **`injectHosts`**: the bare address in its canonical compressed form, such as `"::1"` or `"2001:db8::1"`. The proxy matches each entry against the connection's bare destination address, ignoring ports, so a bracketed, zone-ID, or differently compressed spelling never matches and the proxy never injects the credential there.
>
> `/doctor` flags `injectHosts` entries that can never match with the warning `Sandbox credential injectHosts entries can never match their destination`.
>
> Unlike `deny`, masking authorizes the proxy to send your real credential to the listed hosts, so it is honored only from settings you or your administrator control: user settings, managed settings, and the `--settings` CLI flag. `mask` entries, `network.tlsTerminate`, and [`credentials.allowPlaintextInject`](/docs/en/settings-reference#sandbox-credentials-allowplaintextinject), which lets the proxy inject credentials into unencrypted requests, are all ignored in a repository's `.claude/settings.json` or `.claude/settings.local.json`.
>
> When the same variable is listed with `deny` in any scope, `deny` takes precedence.
>
> Masking replaces the variable's entire value by default, which suits a bare token. Optional entry fields, which require Claude Code v2.1.224 or later, handle values with structure:
>
> * `extract`: a regular expression Claude Code applies across the value, replacing only the text captured by group 1 of each match, so a tool that parses the value, such as a `DATABASE_URL` connection string, still works inside the sandbox. The pattern must contain at least one capturing group.
> * `onExtractNoMatch` controls what happens when the pattern matches nothing:
>   * `warn`, the default, warns and passes the variable through unmasked
>   * `deny` unsets the variable inside the sandbox
>   * `error` stops sandbox setup until you fix the configuration
> * `decode: "jwt"`: for a variable holding a JSON Web Token (JWT). Claude Code verifies the value is a JWT and replaces it with a structurally valid fake token, so code inside the sandbox that decodes the token keeps working. Add `maskClaims` to list top-level payload claims to mask individually instead of replacing the whole token; the other claims stay readable. When the value doesn't verify as a JWT, or no listed claim matches, Claude Code passes the variable through unmasked with a warning. `decode` can't be combined with `extract`.
>
> See the [`credentials.envVars[]` rows in the settings reference](/docs/en/settings-reference#sandbox-settings) for the full field list.
>
> #### Re-sign AWS requests
>
> AWS requests carry SigV4 signatures over the request contents, so mask `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` together. The proxy detects a SigV4 request by the access key's sentinel and re-signs it after substituting the real values. Masking the secret alone leaves requests signed with the placeholder, which the proxy can't detect, so they fail at AWS; Claude Code warns about this case at startup, but not when only the access key ID is masked. A detected request the proxy can't re-sign, such as one missing its `x-amz-date` header, fails with a proxy error instead of reaching the server with a broken signature.
>
> Claude Code links the conventional `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN` variables into one credential automatically when you mask their whole values. If your AWS credential lives in variables with other names, group them yourself with [`credentials.awsPairs`](/docs/en/settings-reference#sandbox-credentials-awspairs), which requires Claude Code v2.1.224 or later. This example adds the pairing to a configuration that already masks `MY_KEY_ID`, `MY_SECRET_KEY`, and `MY_SESSION_TOKEN` whole-value, as in the [masking configuration above](#mask-environment-variables):
>
> ```json
> {
>   "sandbox": {
>     "credentials": {
>       "awsPairs": [
>         {
>           "accessKeyIdVar": "MY_KEY_ID",
>           "secretAccessKeyVar": "MY_SECRET_KEY",
>           "sessionTokenVar": "MY_SESSION_TOKEN"
>         }
>       ]
>     }
>   }
> }
> ```
>
> Each entry follows these rules:
>
> * `accessKeyIdVar` and `secretAccessKeyVar` name the masked `envVars` entries holding the access key ID and the secret key. The optional `sessionTokenVar` names the entry holding the session token for temporary credentials; when set, the proxy sends the real token as `x-amz-security-token` on re-signed requests.
> * Each named variable must be a `mask` entry that masks its entire value, without `extract` or `decode`.
> * The proxy re-signs requests on the hosts listed in the access key ID entry's `injectHosts`.
> * Naming any of the conventional variables in a pair replaces the automatic pairing.
>
> Like `mask` entries, `awsPairs` is honored only from user settings, managed settings, and the `--settings` CLI flag.
>
> Three AWS request forms carry signatures the proxy can't recompute. When such a request is signed with a masked pair's placeholder, the proxy fails it rather than forward a broken signature; requests signed with unmasked credentials are never affected. The [`credentials.sigv4`](/docs/en/settings-reference#sandbox-credentials-sigv4) setting, which requires Claude Code v2.1.224 or later, relaxes this per form: setting a form's key to `passthrough` forwards the request with its placeholder-derived signature, so the calling tool receives AWS's own rejection response instead of a proxy error. Like `awsPairs`, `sigv4` is honored only from user settings, managed settings, and the `--settings` CLI flag.
>
> | Request form                  | `sigv4` key | Why the proxy can't re-sign it                                                                    |
> | :---------------------------- | :---------- | :------------------------------------------------------------------------------------------------ |
> | aws-chunked streaming uploads | `streaming` | Per-chunk signatures chain off the seed signature, so re-signing would require rewriting the body |
> | Presigned URLs                | `presigned` | The signature lives in the URL itself, with no `Authorization` header                             |
> | SigV4A asymmetric signatures  | `sigv4a`    | There is no shared-key HMAC to recompute                                                          |
>
> #### Mask credential files
>
> File entries also accept `"mode": "mask"`, which requires Claude Code v2.1.221 or later. What a sandboxed command sees depends on the platform:
>
> * **Linux and WSL2**: sandboxed commands read a sentinel copy of the file, a stand-in whose secret is replaced with a placeholder value, and the [sandbox proxy](#network-isolation) substitutes the real value on egress.
> * **macOS**: sandboxed commands can't read the listed file at all. Claude Code builds no sentinel copy and substitutes nothing on egress, so tools that authenticate with the file don't work inside the sandbox, the same effect as `deny`. Unlike a `deny` entry, the read block holds even when you [disable filesystem isolation](#disable-filesystem-isolation).
>
> On every platform, the [`network.tlsTerminate`](/docs/en/settings-reference#sandbox-network-tlsterminate) requirement, `injectHosts`, and the settings-source restriction work the same way as for [masked environment variables](#mask-environment-variables).
>
> The example below masks a GitHub token stored in `~/.config/gh/hosts.yml`; the `extract` pattern, covered below, tells Claude Code which part of the file is the secret. On Linux and WSL2, sandboxed commands that read the file get a sentinel in place of the token, and the proxy substitutes the real token on requests to `api.github.com`:
>
> ```json
> {
>   "sandbox": {
>     "enabled": true,
>     "network": {
>       "tlsTerminate": {},
>       "allowedDomains": ["*.github.com"]
>     },
>     "credentials": {
>       "files": [
>         {
>           "path": "~/.config/gh/hosts.yml",
>           "mode": "mask",
>           "extract": "oauth_token:\\s*(\\S+)",
>           "injectHosts": ["api.github.com"]
>         }
>       ]
>     }
>   }
> }
> ```
>
> To confirm the mask is active, ask Claude to run `cat ~/.config/gh/hosts.yml` in a sandboxed command: on Linux and WSL2 the output shows a sentinel value in place of the token, and on macOS the read fails instead.
>
> On Linux and WSL2, the `extract` pattern is what keeps the rest of `hosts.yml` readable. Claude Code applies the regular expression across the whole file and replaces only the text captured by group 1 of each match, so `gh` still parses its config and only the token is a placeholder. Use `extract` for any structured file that tools parse, such as `.netrc`, JSON, or YAML; the pattern must contain at least one capturing group. Without `extract`, Claude Code replaces the entire file content with one sentinel value, which suits a file that holds a single bare secret and nothing else.
>
> For a file that holds a JSON Web Token (JWT), set `decode: "jwt"` instead of, or together with, `extract`. `decode` requires Claude Code v2.1.224 or later. Claude Code finds JWT candidates with a built-in pattern, or with your `extract` pattern when set, verifies each candidate is a JWT, and replaces it with a structurally valid fake token, so code that decodes the token inside the sandbox keeps working. Add `maskClaims` to mask only the named top-level payload claims inside each verified token and leave the other claims readable. When no candidate verifies, or no named claim matches, the `onExtractNoMatch` field below governs the outcome, as it does for a pattern that matches nothing.
>
> Two optional fields refine how matching behaves. Both apply only when `mode` is `mask` and `extract` or `decode` is set. On macOS, Claude Code applies `mask` entries as `deny` before the pattern runs whenever filesystem isolation is on, so these fields, and the no-match outcomes below, take effect there only when [filesystem isolation is off](#disable-filesystem-isolation):
>
> * `onExtractNoMatch` controls what happens when matching finds nothing to mask in the file:
>
>   * `warn`, the default, warns and skips the entry, so sandboxed commands can read the real file unmasked. The default suits credentials that may be legitimately absent; if the secret might be present but the pattern might miss it, use `deny`
>   * `deny` makes the file unreadable instead
>   * `error` stops sandbox setup until you fix the configuration
>
>   Claude Code treats `deny` as `error` whenever the read block wouldn't be enforced: when you [disable filesystem isolation](#disable-filesystem-isolation), and when a `filesystem.allowRead` entry from any settings source re-opens the file's path.
> * `maskDuplicates` also replaces verbatim copies of each masked credential value, an `extract` capture or a `decode`-verified token, found outside the matched spans, for a secret repeated where matching doesn't reach. It matches raw substrings, so a short or common value would be replaced everywhere it appears; reserve it for long, high-entropy secrets. Default: false.
>
> `mask` applies to a single file, so list each credential file individually. Claude Code falls back to `deny` for a `mask` entry it can't mask safely: a directory path, a glob pattern, a file larger than 8 MiB, or a file that isn't UTF-8 text. Write directories as explicit `deny` entries instead; the table under [Which settings can disable it](#which-settings-can-disable-it) covers whether each form pins `filesystem.disabled` and how it behaves with filesystem isolation off.
>
> ## How sandboxing works
>
> ### Filesystem isolation
>
> The sandboxed Bash tool restricts### Source: Sandbox environments
>
> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Choose a sandbox environment
>
> > Compare Claude Code sandbox options: the built-in sandboxed Bash tool, sandbox runtime, dev containers, Docker, and VMs. Choose the right isolation for your threat model.
>
> Isolating Claude Code limits what a session can read, write, and reach on the network. This matters most when you let Claude work with fewer permission prompts, run it unattended, or point it at code you do not fully trust.
>
> Claude Code can run in several kinds of isolated environments, ranging from a lightweight per-command sandbox to a fully separate virtual machine. This page compares them by what they isolate and what they require, helps you choose one for your threat model, and shows how to enforce that choice across an organization.
>
>   For the broader security model, see [Security](/docs/en/security). For Agent SDK deployments, see [Secure deployment](/docs/en/agent-sdk/secure-deployment).
>
> ## Compare sandboxing approaches
>
> The first two approaches in the table below run on the host operating system without containers. The rest place Claude Code inside a container or virtual machine.
>
> | Approach                                          | What is isolated                                                            | Requires Docker | Setup effort                                                                            |
> | :------------------------------------------------ | :-------------------------------------------------------------------------- | :-------------- | :-------------------------------------------------------------------------------------- |
> | [Sandboxed Bash tool](#sandboxed-bash-tool)       | Bash commands and their child processes                                     | No              | Minimal on macOS; low on Linux and WSL2                                                 |
> | [Sandbox runtime](#sandbox-runtime)               | The whole Claude Code process, including file tools, MCP servers, and hooks | No              | Low                                                                                     |
> | [Dev container](#dev-containers)                  | Full development environment                                                | Yes             | Medium                                                                                  |
> | [Custom container](#custom-container)             | Full development environment                                                | Yes             | Medium to high                                                                          |
> | [Virtual machine](#virtual-machine)               | Full operating system                                                       | No              | High                                                                                    |
> | [Claude Code on the web](#claude-code-on-the-web) | Full operating system, hosted by Anthropic                                  | No              | None; requires a Claude subscription, and GitHub when you launch from the web interface |
>
> The [sandboxed Bash tool](/docs/en/sandboxing) is built into Claude Code and restricts only Bash commands. Built-in file tools, MCP servers, and hooks still run directly on your host. Every other approach in the table puts the whole Claude Code process inside the isolation boundary, so file tools, MCP servers, and hooks are restricted too.
>
>   Sandbox isolation reduces the impact of a breach, but it does not eliminate risk. Any approach that allows network egress can still leak data the agent can read, and any approach that mounts your project directory writable can still modify that code. Review the [security limitations](/docs/en/sandboxing#security-limitations) before relying on a sandbox as a hard control.
>
>   Isolation also does not change what is sent to the model. Your prompts and the files Claude reads are transmitted to the Anthropic API or your configured provider with or without a sandbox. See [Data usage](/docs/en/data-usage) for what Claude Code sends and how to reduce it.
>
> ## Choose an approach
>
> Match your goal to a row below, then read the detail section that follows.
>
> | You want to                                                                   | Start with                                                                                                                                                                             |
> | :---------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | Reduce permission prompts during everyday work on your own machine            | The [sandboxed Bash tool](/docs/en/sandboxing), configured with `/sandbox`                                                                                                                  |
> | Let Claude work unattended with `--dangerously-skip-permissions` or auto mode | The preconfigured [dev container](/docs/en/devcontainer), any container or VM, or the [sandbox runtime](#sandbox-runtime)                                                                   |
> | Isolate MCP servers and hooks as well as Bash, without Docker                 | The sandbox runtime                                                                                                                                                                    |
> | Work on an untrusted repository                                               | A dedicated virtual machine, or [Claude Code on the web](/docs/en/claude-code-on-the-web) if you have a Claude subscription; GitHub is only required when you launch from the web interface |
> | Standardize a sandboxed environment across a team                             | The preconfigured [dev container](/docs/en/devcontainer), copied into your repository                                                                                                       |
> | Use Claude Code from a device with no local setup                             | [Claude Code on the web](/docs/en/claude-code-on-the-web), which requires a Claude subscription and a connected GitHub account                                                              |
> | Require isolation for every developer in your organization                    | [Enforce isolation across an organization](#enforce-isolation-across-an-organization)                                                                                                  |
> | Work on a native Windows host                                                 | A container or VM, or run the Bash sandbox inside WSL2                                                                                                                                 |
>
> ### How isolation relates to permission modes
>
> [Permission modes](/docs/en/permission-modes) decide whether a tool call runs and whether you are prompted first. Isolation restricts what a command can access once it runs. The two work together: when a permission mode lets actions run without asking you, an isolation boundary limits what those actions can reach.
>
> When you pass `--dangerously-skip-permissions`, Claude acts without asking you first. The [actions no mode auto-approves](/docs/en/permission-modes#actions-no-mode-auto-approves) still apply.
>
> With no prompts to catch mistakes, the isolation boundary you choose is what protects your system. Always run `--dangerously-skip-permissions` sessions inside a container, a VM, or the [sandbox runtime](#sandbox-runtime), so that file tools, MCP servers, and hooks are also inside the boundary. On Linux and macOS, Claude Code refuses to start with this flag when running as root, so run the container, VM, or sandbox runtime as a non-root user.
>
> [Auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) replaces the prompt with a classifier that reviews actions. The classifier is a per-action control, not an isolation boundary, so an isolation boundary still adds defense in depth for unattended runs, and is not required the way it is for `--dangerously-skip-permissions`.
>
> The [sandboxed Bash tool](#sandboxed-bash-tool) on its own constrains only Bash, so it is not sufficient for fully unattended runs in either mode. You can layer approaches: running the sandboxed Bash tool inside a container or VM gives you OS-level command restrictions on top of the outer environment boundary. For how the Bash sandbox itself interacts with permission rules and modes, see [How sandboxing relates to permissions and permission modes](/docs/en/sandboxing#how-sandboxing-relates-to-permissions-and-permission-modes).
>
> ## Sandboxed Bash tool
>
>   This option does not support native Windows. On Windows hosts, use WSL2 or one of the container or VM approaches below.
>
> The sandboxed Bash tool is built into Claude Code. It uses operating system primitives to restrict the filesystem and network access of every Bash command Claude runs.
>
> Run the `/sandbox` command to open the sandbox panel and choose a mode. The [Sandboxing](/docs/en/sandboxing) guide covers the approval modes, the default boundary, and how to widen or narrow it.
>
> The per-command sandbox does not cover everything that runs in a session:
>
> * Other [built-in tools](/docs/en/tools-reference) such as Read, Edit, and WebFetch run inside the Claude Code process and do not spawn arbitrary code. [Permission rules](/docs/en/permissions) for path or domain gate them instead.
> * [MCP](/docs/en/mcp) servers and hooks are separate processes that run unconstrained on the host.
>
> To put built-in tools, MCP servers, and hooks all behind one OS boundary, run the whole Claude Code process inside the [sandbox runtime](#sandbox-runtime), the [dev container](#dev-containers), or a [custom container](#custom-container).
>
> ## Sandbox runtime
>
> The [`@anthropic-ai/sandbox-runtime`](https://github.com/anthropic-experimental/sandbox-runtime) package wraps an entire process in the same Seatbelt or bubblewrap isolation that the built-in Bash sandbox uses. Running Claude Code through the runtime constrains every tool, hook, and MCP server in the session, not only Bash. The runtime is a beta research preview, and its configuration format may change as the package evolves.
>
> This section covers what you configure and what the runtime enforces on its own. For deploying the runtime in Agent SDK applications, see the [secure deployment guide](/docs/en/agent-sdk/secure-deployment#sandbox-runtime).
>
> ### Set up and launch the runtime
>
> On Linux and WSL2, the runtime relies on the same `bubblewrap` and `socat` packages as the built-in sandbox, plus `ripgrep`, which Claude Code bundles but the standalone runtime resolves from your PATH. Install `bubblewrap` and `socat` as described in [Set up Linux and WSL2](/docs/en/sandboxing#set-up-linux-and-wsl2), and `ripgrep` from your distribution's package manager. On macOS you need no additional packages. The runtime uses the built-in Seatbelt sandbox there.
>
> By default the runtime denies network access and confines writes to a small set of built-in runtime paths, so configure it before launching Claude Code through it. Put your configuration in `~/.srt-settings.json`, or in a file you pass with `--settings`. The package [README](https://github.com/anthropic-experimental/sandbox-runtime) documents the full configuration schema.
>
> Allow write access to at least:
>
> * Your project directory.
> * Claude Code's configuration paths `~/.claude` and `~/.claude.json`.
> * `/tmp`, where Claude Code writes runtime files.
>
> Allow the network domains your session needs:
>
> * `api.anthropic.com`, or your configured provider's endpoint. On a third-party provider, keep `api.anthropic.com` as well: the WebFetch domain safety check still calls it by default unless you set `skipWebFetchPreflight: true`.
> * `claude.ai` and `platform.claude.com`, which [OAuth sign-in and token refresh](/docs/en/network-config#network-access-requirements) require. Runs authenticated with an API key can drop these two.
>
> On Linux and WSL2, the runtime applies write grants only to paths that already exist. In a fresh environment, create Claude Code's configuration paths before the first launch:
>
> ```bash
> mkdir -p ~/.claude && echo '{}' > ~/.claude.json
> ```
>
> Once the settings file is in place, launch Claude Code with `npx` and pass `claude` as the command to wrap:
>
> ```bash
> npx @anthropic-ai/sandbox-runtime claude
> ```
>
> Claude Code starts inside the sandbox with the filesystem and network boundaries you configured. The same command works for sandboxing standalone MCP servers or other helper processes.
>
> ### What the runtime blocks on its own
>
> The runtime blocks the highest-risk writes without any configuration from you:
>
> * `denyWrite` takes precedence over `allowWrite`.
> * At the project root, the runtime denies `.git/hooks`, denies `.git/config` unless you set `filesystem.allowGitConfig: true`, and denies `.mcp.json`, `.claude/commands`, `.claude/agents`, and shell startup files.
> * On macOS, these denies are checked when a write happens, so they also cover nested files and repositories created during the session.
> * On Linux and WSL2, the runtime builds the deny list once at launch. It reliably covers the project root, makes a best-effort shallow scan for nested copies that exist at that point, and does not cover anything the session creates later, such as `git init`, `git clone`, or scaffolding. The README's `mandatoryDenySearchDepth` section describes the scan's exact semantics.
> * Without a valid `~/.srt-settings.json`, the runtime starts anyway, blocks network access, and confines writes to built-in runtime paths such as `/tmp/claude`, `~/.npm/_logs`, and `~/.claude/debug`. Don't take a clean start as proof your settings loaded.
> * When you pass `--settings`, the runtime refuses to start if the file fails to load.
>
> Your write grants still include other paths Claude Code loads configuration from, so deny those with `denyWrite`. A sandboxed session that can write them can persist hooks, permission rules, or MCP servers that run unsandboxed the next time you launch Claude Code.
>
> ### After unattended runs
>
> Review the paths you kept writable. On Linux and WSL2, also review anything the session created.
>
> ## Dev containers
>
> A dev container runs Claude Code inside a Docker container that VS Code or a compatible editor manages, with your project mounted in. You can define your own with a `.devcontainer/` directory in your repository.
>
> The claude-code repository publishes an [example dev container](/docs/en/devcontainer) with a default-deny iptables firewall as a starting point. Copy it into your repository and adjust the firewall allowlist, base image, and pinned Claude Code version to fit your environment. Be### Source: Checkpointing
>
> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Checkpointing
>
> > Track, rewind, and summarize Claude's edits and conversation to manage session state.
>
> Claude Code automatically tracks Claude's file edits as you work, allowing you to quickly undo changes and rewind to previous states if anything gets off track.
>
> ## How checkpoints work
>
> As you work with Claude, checkpointing automatically captures the state of your code before each user prompt.
>
> ### Automatic tracking
>
> Claude Code tracks all changes made by its file editing tools:
>
> * Every user prompt creates a new checkpoint
> * Claude Code keeps file snapshots for the 100 most recent checkpoints in a session. Discarding an older checkpoint deletes the snapshot files that no remaining checkpoint references, except each file's first snapshot, which the VS Code extension uses as the baseline for its session diffs.
> * Claude Code saves checkpoints with the conversation, so you can still run `/rewind` after you resume a session
> * Claude Code deletes checkpoints along with sessions after 30 days, following the [retention sweep rules](/docs/en/claude-directory#cleaned-up-automatically); change the period with [`cleanupPeriodDays`](/docs/en/settings-reference#cleanupperioddays)
>
> ### Rewind and summarize
>
> Run `/rewind`, or press `Esc` twice when the prompt input is empty, to open the rewind menu.
>
>   If the prompt input contains text, double `Esc` clears it instead of opening the menu. The cleared text is saved to your input history, so press `Up` to recall it after you finish in the rewind menu.
>
> The rewind menu lists each prompt you sent during the session. Select the point you want to act on, then choose an action:
>
> * **Restore code and conversation**: revert both code and conversation to that point
> * **Restore conversation**: rewind to that message while keeping current code
> * **Restore code**: revert file changes while keeping the conversation
> * **Summarize from here**: compress the conversation from this point forward into a summary, freeing context window space
> * **Summarize up to here**: compress the conversation before this point into a summary, keeping later messages intact
> * **Never mind**: return to the message list without making changes
>
> The two code restore options appear only when the selected checkpoint has tracked file changes to revert. If no file edits were captured after that point, the menu offers only **Restore conversation**, the summarize options, and **Never mind**.
>
> After restoring the conversation or choosing Summarize from here, the original prompt from the selected message is restored into the input field so you can re-send or edit it.
>
> Choosing Summarize up to here leaves you at the end of the conversation with the input empty. With either summarize option, a **Summarized conversation** marker appears in the conversation where the compressed messages were.
>
> #### Rewind past a cleared conversation
>
> If you ran `/clear` earlier in the same Claude Code process, the rewind menu shows an additional entry at the top of the list labeled `/resume <session-id> (previous session)`. Select it to resume the conversation that was active before `/clear` ran. The entry is available until you exit Claude Code or resume a different session, and requires Claude Code v2.1.191 or later. On earlier versions, run `/resume` and pick the previous session from the list instead.
>
> #### Guide a summary
>
> Summarizing doesn't change files on disk, and the original messages stay in the session transcript, so Claude can still reference the details. To guide what the summary focuses on, highlight a **Summarize** option with the arrow keys and type instructions where the row reads **add context (optional)**, then press `Enter`. Selecting the option with its number key summarizes immediately without instructions.
>
>   Summarize keeps you in the same session and compresses context, like a targeted `/compact`. To branch off and try a different approach while preserving the original session intact, use [`/branch`](/docs/en/sessions#branch-a-session) or `claude --continue --fork-session` instead.
>
> ## Common use cases
>
> Checkpoints are particularly useful when:
>
> * **Exploring alternatives**: try different implementation approaches without losing your starting point
> * **Recovering from mistakes**: quickly undo changes that introduced bugs or broke functionality
> * **Iterating on features**: experiment with variations knowing you can revert to working states
> * **Freeing context space**: summarize a verbose debugging session from the midpoint forward, keeping your initial instructions intact
>
> ## Limitations
>
> ### Bash command changes not tracked
>
> Checkpointing does not track files modified by bash commands. For example, if Claude Code runs:
>
> ```bash
> rm file.txt
> mv old.txt new.txt
> cp source.txt dest.txt
> ```
>
> These file modifications cannot be undone through rewind. Only direct file edits made through Claude's file editing tools are tracked.
>
> ### Subagent edits not restored
>
> A [subagent](/docs/en/sub-agents) makes edits with Claude's file editing tools, but Claude Code usually doesn't capture those edits in your session's checkpoints. Whether rewinding restores them depends on how the subagent runs:
>
> * **Foreground forked skill**: a [skill with `context: fork`](/docs/en/skills#run-skills-in-a-subagent) that runs in the foreground edits your working tree during your own turn, so rewinding restores its edits as usual. Set `background: false` to run a fork in the foreground; a few situations, [listed on the skills page](/docs/en/skills#run-skills-in-a-subagent), run it there regardless of the setting.
> * **Any other subagent**: rewinding doesn't restore the edits. Use git to revert them. This includes a forked skill that runs in the background, the default, and a background [`/code-review --fix`](/docs/en/code-review) run.
>
> ### External changes not tracked
>
> Checkpointing only tracks files that have been edited within the current session. Manual changes you make to files outside of Claude Code and edits from other concurrent sessions are normally not captured, unless they happen to modify the same files as the current session.
>
> ### Symlinked and hard-linked paths not restored
>
> Checkpointing doesn't rewind symlinked or hard-linked files. When you pick **Restore code** or **Restore code and conversation** from the `/rewind` menu, Claude Code skips any tracked path that is a symlink or hard link and shows a `Restored the code, but skipped N files` warning. The skipped files keep their current contents. To undo the session's changes to one of them, ask Claude to reverse the edit or edit the file yourself. Config files a dotfile manager symlinks into your project and files pnpm hard-links into place both fall into this category.
>
> To see which paths a restore skips, turn on debug logging with `/debug` before you restore: the debug log at `~/.claude/debug/<session-id>.txt` names each skipped path. For every skip reason and the recovery steps, see [the skipped-files entry in the error reference](/docs/en/errors#restored-the-code-but-skipped-files).
>
> ### Not a replacement for version control
>
> Checkpoints are designed for quick, session-level recovery. For permanent version history and collaboration, continue using version control, such as Git, for commits, branches, and long-term history.
>
> ## See also
>
> * [Interactive mode](/docs/en/interactive-mode) - Keyboard shortcuts and session controls
> * [Commands](/docs/en/commands) - Accessing checkpoints using `/rewind`
> * [CLI reference](/docs/en/cli-reference) - Command-line options