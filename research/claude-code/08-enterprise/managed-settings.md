---
primary_sources:
  - id: T1-MANAGED
    title: "Managed settings"
    url: "https://code.claude.com/docs/en/managed-settings.md"
    section: ""
  - id: T1-SERVER-MANAGED
    title: "Server-managed"
    url: "https://code.claude.com/docs/en/server-managed-settings.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Managed settings

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Managed settings
>
> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Deploy managed settings
>
> > Deploy managed settings to every developer's machine: delivery mechanisms per OS, which managed source Claude Code uses, and how to verify enforcement.
>
> Managed settings are the settings your organization deploys to every developer's machine. Claude Code applies them above every other level, so no user, project, local, or `--settings` value overrides them, apart from a few [security-sensitive exceptions](/docs/en/settings#exceptions-to-managed-settings-precedence) where a stricter value from a lower level still counts.
>
> This page is for the administrator who deploys managed settings or debugs why one isn't applying. To decide what to enforce, start with the [Decide what to enforce](/docs/en/admin-setup#decide-what-to-enforce) table. For the claude.ai console path, see [Server-managed settings](/docs/en/server-managed-settings). For which file a developer's own values go in, see [Settings](/docs/en/settings).
>
> ## Deploy a managed settings file
>
> This is the quickest way to put a policy on each machine: a `managed-settings.json` file. If you haven't picked how to deliver managed settings yet, or your devices are under MDM or developers run cloud sessions, read [Choose a delivery mechanism](#choose-a-delivery-mechanism) first.
>
>   **Write managed-settings.json**: Write a `managed-settings.json` that holds the keys you've decided to enforce, in the same JSON shape as `settings.json`. The [Decide what to enforce](/docs/en/admin-setup#decide-what-to-enforce) table lists the keys behind each control, and each entry in the [settings reference](/docs/en/settings-reference) says whether a managed source can set it. This file blocks two file reads, turns off bypass mode, and makes Claude Code ignore permission rules from user, project, and local files and from `--allowedTools`:
>
>     ```json managed-settings.json
>     {
>       "permissions": {
>         "deny": [
>           "Read(./.env)",
>           "Read(./secrets/**)"
>         ],
>         "disableBypassPermissionsMode": "disable"
>       },
>       "allowManagedPermissionRulesOnly": true
>     }
>     ```
>
>     For a fuller example that shows the shape of more managed keys, including the login method, models, MCP servers, and marketplaces, see [An organization's managed settings](/docs/en/settings-example#an-organizations-managed-settings).
>
>   **Place the file on each machine**: Save the file as `managed-settings.json` in the system directory for the operating system, using whatever tooling already places files on your fleet:
>
>     * **macOS**: `/Library/Application Support/ClaudeCode/managed-settings.json`
>     * **Linux and WSL**: `/etc/claude-code/managed-settings.json`
>     * **Windows**: `C:\Program Files\ClaudeCode\managed-settings.json`
>
>   **Confirm the policy applied**: On one machine, run `/status` inside Claude Code. The `Setting sources` line shows `Enterprise managed settings (file)`. Roll out to the rest of the fleet after that; [Check that a policy is in force](#check-that-a-policy-is-in-force) covers what to look at when the line is missing.
>
> ## Choose a delivery mechanism
>
> The file in the steps above is one of four ways to get managed settings onto a machine. Every mechanism carries the same policy keys as a `settings.json` file, so the [settings reference](/docs/en/settings-reference) applies to all of them. A few keys are tied to particular sources, and each entry's Scope line says which:
>
> * **Delivery controls**: [`policyHelper`](/docs/en/settings-reference#policyhelper) and [`wslInheritsWindowsSettings`](/docs/en/settings-reference#wslinheritswindowssettings)
> * **Gateway login keys**: [`forceLoginGatewayUrl`](/docs/en/settings-reference#forcelogingatewayurl) and the `"gateway"` value of [`forceLoginMethod`](/docs/en/settings-reference#forceloginmethod)
>
> A managed settings file, an MDM profile, or the claude.ai console applies one policy to everyone it reaches. To give one group of developers a different policy, deploy a different file or profile to that group; the claude.ai console [can't target a group yet](/docs/en/server-managed-settings#current-limitations), while a self-hosted [Claude apps gateway](/docs/en/claude-apps-gateway) delivers managed settings per IdP group.
>
> When more than one mechanism delivers a policy to the same machine, Claude Code uses one and ignores the others; [Which managed source Claude Code uses](#which-managed-source-claude-code-uses) gives the order.
>
> The MDM and file rows are together called endpoint-managed settings, because the policy is stored on the developer's device, as opposed to the server-managed row, where Claude Code fetches it.
>
> Pick a mechanism by how you already manage devices, using the table below.
>
> | Mechanism                                              | How you deliver it                                                                                                                                                                                                | When Claude Code reads it                                                                                                                                                                                                                | Use it when                                                                                    |
> | :----------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------- |
> | [Server-managed settings](/docs/en/server-managed-settings) | In the claude.ai admin console, or on a self-hosted [Claude apps gateway](/docs/en/claude-apps-gateway)                                                                                                                | Fetched at startup and polled hourly; see [changes that need approval](#where-and-when-a-policy-applies)                                                                                                                                 | You want one place to change policy for a claude.ai organization without touching each machine |
> | MDM or OS-level policy                                 | As a macOS configuration profile or a Windows `HKLM` registry value, through Jamf, Intune, Group Policy, or a similar tool; see [where each mechanism stores the policy](#where-each-mechanism-stores-the-policy) | Read at startup and checked for changes every 30 minutes                                                                                                                                                                                 | You already manage devices with MDM or Group Policy                                            |
> | File-based                                             | As `managed-settings.json` in a system directory on each machine; see [where each mechanism stores the policy](#where-each-mechanism-stores-the-policy)                                                           | Read at startup and reloaded when a file changes                                                                                                                                                                                         | Machines without MDM, Linux hosts, or images you build yourself                                |
> | HKCU registry, Windows and WSL                         | As a Windows `HKCU` registry value; see [where each mechanism stores the policy](#where-each-mechanism-stores-the-policy)                                                                                         | Read at startup and checked for changes every 30 minutes; Claude Code uses it only when no other managed source delivers a policy key and no [host-supplied parent settings](#let-an-embedding-host-add-policy) supply a restrictive key | You can't write the machine-level `HKLM` key                                                   |
>
> Starter templates for Jamf, Iru, Intune, and Group Policy are in the [MDM examples repository](https://github.com/anthropics/claude-code/tree/main/examples/mdm).
>
> For managed MCP servers, which you deploy alongside any of these through `managed-mcp.json`, see [Managed MCP configuration](/docs/en/managed-mcp).
>
> ### Where and when a policy applies
>
> A deployed policy reaches the developer's sessions as follows:
>
> * **Surfaces**: every surface that runs Claude Code on the machine reads these sources: the terminal, the VS Code and JetBrains extensions, the desktop app, and [Agent SDK](/docs/en/agent-sdk/typescript) sessions, which load managed settings even when `settingSources` excludes the user, project, and local files.
> * **Cloud sessions**: a session in an Anthropic-hosted environment doesn't read a device's MDM profile or file, so policy for it has to come from server-managed settings. A session in a [self-hosted environment](/docs/en/self-hosted-environments) reads the managed settings file in its runner image only when server-managed settings deliver no keys, apart from the [keys Claude Code reads from every admin source](#keys-read-from-every-admin-source).
> * **Running sessions**: a session picks up most changes on the schedule in the table without a restart. Claude Code reads [`forceRemoteSettingsRefresh`](/docs/en/settings-reference#forceremotesettingsrefresh) and [`requiredMinimumVersion`](/docs/en/settings-reference#requiredminimumversion) only at session start, arms a new or changed [`policyHelper`](/docs/en/settings-reference#policyhelper) entry at the next launch, and reads [some user-editable keys once at session start](/docs/en/settings#when-edits-take-effect).
> * **Changes that need approval**: apart from the [updates that wait for the next launch](/docs/en/server-managed-settings#fetch-and-caching-behavior), a server-managed change to a setting that [needs approval](/docs/en/server-managed-settings#security-approval-dialogs), such as a hook or an `env` variable, waits for the developer to accept the dialog in an interactive session, and applies for the current run in a session an IDE extension or the Agent SDK hosts. Other server-managed changes apply on the next poll.
> * **Long-lived sessions**: a session left open for weeks can still lag a rollout. [`requiredMinimumVersion`](/docs/en/settings-reference#requiredminimumversion) blocks an outdated binary from starting and doesn't end a session that's already running.
>
> ### Where each mechanism stores the policy
>
> The keys are the same everywhere, but each mechanism stores them in a different place and shape:
>
> * **Server-managed**: Anthropic's servers, or your gateway, hold the policy. Claude Code keeps a local cache that it applies at startup and [replaces on each successful fetch](/docs/en/server-managed-settings#security-considerations).
> * **macOS configuration profile**: the `com.anthropic.claudecode` managed preferences domain. Use the same top-level keys as `managed-settings.json`, with nested settings as dictionaries and lists as plist arrays.
> * **Windows HKLM registry**: the JSON as a `REG_SZ` or `REG_EXPAND_SZ` value named `Settings` under `HKLM\SOFTWARE\Policies\ClaudeCode`.
> * **File-based**: `managed-settings.json`, an optional `managed-settings.d/` directory, and `managed-mcp.json` in the system directory: `/Library/Application Support/ClaudeCode/` on macOS, `/etc/claude-code/` on Linux and WSL, and `C:\Program Files\ClaudeCode\` on Windows. Claude Code doesn't read the legacy Windows path `C:\ProgramData\ClaudeCode\managed-settings.json`.
> * **Windows HKCU registry**: the same `Settings` value under `HKCU\SOFTWARE\Policies\ClaudeCode`.
>
> ### Split a file-based policy across teams
>
> If several teams own parts of one policy, put each part in its own file in `managed-settings.d/`, next to `managed-settings.json` in the same system directory, instead of editing one shared file.
>
> Claude Code merges `managed-settings.json` first, then every `*.json` file in the directory in alphabetical order. Name the files with numeric prefixes to control the order, such as `10-telemetry.json` and `20-security.json`. Claude Code ignores hidden files and files that don't end in `.json`.
>
> When two files set the same key, Claude Code combines them by these rules:
>
> * **Single values**, such as `"model": "opus"` or `"cleanupPeriodDays": 7`: the later file's value replaces the earlier one
> * **Lists**, such as `permissions.deny` or `sandbox.network.allowedDomains`: the two lists combine, with duplicates removed
> * **Nested blocks**, such as `env` or `sandbox`: the two blocks merge key by key, and each key inside follows these same rules
> * **`fallbackModel`**: the later chain replaces the earlier one whole
> * **[`extraKnownMarketplaces`](/docs/en/settings-reference#extraknownmarketplaces)**: a later entry with the same name replaces the earlier one whole
> * **[`modelPicker`](/docs/en/settings-reference#modelpicker)**: the later lineup replaces the earlier one whole
>
> ## Which managed source Claude Code uses
>
> When your organization delivers more than one managed source, Claude Code uses the first of these sources that delivers at least one policy key and ignores the rest rather than merging them, apart from the few cross-source keys covered in [Keys read from every admin source](#keys-read-from-every-admin-source).
>
> Claude Code shows no warning for the sources it skips. To see which source it used, run `/status`.
>
> Two terms recur in this section:
>
> * **Policy key**: any settings key other than the control key [`wslInheritsWindowsSettings`](/docs/en/settings-reference#wslinheritswindowssettings). A managed settings file or MDM policy that contains only that key doesn't count, and Claude Code moves on to the next source.
> * **Admin source**: one of the first three sources below. The HKCU registry is user-writable and isn't one.
>
> Claude Code checks the sources in this order, highest priority first:
>
> 1. Remote settings, delivered from claude.ai as [server-managed settings](/docs/en/server-managed-settings) or by a [Claude apps gateway](/docs/en/claude-apps-gateway). Claude Code fetches this source only when the session authenticates to Anthropic's API directly with an [eligible login or key](/docs/en/server-managed-settings#platform-availability), or signs in to a gateway with `/login`. On other providers, or when `ANTHROPIC_BASE_URL` points somewhere other than Anthropic's API, it starts at the next source
> 2. MDM or OS-level policies: the macOS plist or the HKLM registry key
> 3. Managed settings files, `managed-settings.d/*.json` and `managed-settings.json` merged together
> 4. The HKCU registry, on Windows, and on WSL once the HKLM registry or the Windows managed settings file turns [`wslInheritsWindowsSettings`](/docs/en/settings-reference#wslinheritswindowssettings) on and the HKCU value also sets it. Claude Code reads it only when no source above it delivers a policy key and no [host-supplied parent settings](#let-an-embedding-host-add-policy) supply a restrictive key
>
> This diagram shows the ranking, with examples of the cross-source keys Claude Code reads from the first three sources:
>
> ### Keys read from every admin source
>
> For most keys, Claude Code reads only the [source it selected](#which-managed-source-claude-code-uses). A value in a lower-ranked source is ignored even when the selected source leaves that key unset.
>
> A few keys work differently. Claude Code reads them from every admin source, so a lower-ranked MDM policy or managed settings file can still set them when the selected source doesn't. Claude Code leaves the user-writable HKCU registry out of that scan; when HKCU is the only source and no host supplies parent settings, HKCU applies like any selected source.
>
> The cross-source keys include:
>
> * `sandbox.network.allowManagedDomainsOnly` and `sandbox.filesystem.allowManagedReadPathsOnly`: a `true` in any admin source turns the lock on. While a lock is on, Claude Code unions the allowlist it locks, `sandbox.network.allowedDomains` together with `WebFetch(domain:...)` allow rules, or `sandbox.filesystem.allowRead`, across every admin source. Without the lock, an unselected admin source's allowlist is ignored like any other key from an unselected source
> * `allowAllClaudeAiMcps`
> * The sandbox binary paths `sandbox.bwrapPath` and `sandbox.socatPath`
> * The sandbox `ripgrep` binary, [`sandbox.ripgrep`](/docs/en/settings-reference#sandbox-ripgrep)
> * `sandbox.filesystem.disabled` and `sandbox.network.strictAllowlist`
> * [`useAutoModeDuringPlan`](/docs/en/settings-reference#useautomodeduringplan) and [`syncClaudeAiSkills`](/docs/en/settings-reference#syncclaudeaiskills), where a `false` from any admin source turns the behavior off. A `false` in the developer's user or local settings turns it off too; each key can only deny
> * A commit-trailer opt-out in `attribution`, or in the deprecated `includeCoAuthoredBy`, from any tier
> * [`forceRemoteSettingsRefresh`](/docs/en/server-managed-settings)
> * `env`, merged per variable across the admin sources: each variable comes from the highest-priority source that defines it, so lower sources fill in variables the higher ones leave unset. A few variables follow their own rules; [Per-key exceptions across managed sources](/docs/en/server-managed-settings#per-key-exceptions-across-managed-sources) names each one. Requires Claude Code v2.1.223 or later. Before v2.1.223, Claude Code applied the selected source's whole `env` block only
>
> ### Compute the policy with a helper program
>
> A [`policyHelper`](/docs/en/settings-reference#policyhelper) is an executable your MDM policy or managed settings file names, and Claude Code runs it to compute managed settings at startup. When the selected source configures one and the helper emits a `managedSettings` object, that output changes what Claude Code reads:
>
> * **The emitted `managedSettings` object is the only managed settings for the session**, including for the [keys it otherwise reads from every admin source](#keys-read-from-every-admin-source), apart from `forceRemoteSettingsRefresh`, which Claude Code checks in every admin source at startup before the helper runs. A helper that exits 0 without emitting one contributes nothing, and the sources apply as usual; a helper that fails stops Claude Code from starting, as the [`policyHelper`](/docs/en/settings-reference#policyhelper) entry describes
>
> Claude Code selects the source at startup, and that selection decides whether a helper runs. The [`policyHelper`](/docs/en/settings-reference#policyhelper) entry says which sources can configure a helper.
>
> ### Let an embedding host add policy
>
> When another application launches Claude Code, such as Claude Desktop, an IDE extension, or an Agent SDK app, that host can pass its own managed settings through the SDK `managedSettings` option. Claude Code calls these parent settings.
>
> By default, Claude Code ignores parent settings whenever an admin source is present: server-managed settings, an MDM or OS-level policy, or a managed settings file.
>
> To have Claude Code merge parent settings alongside an admin source, set [`parentSettingsBehavior`](/docs/en/settings-reference#parentsettingsbehavior) to `"merge"` in the highest-priority managed source; Claude Code reads the key from that source only.
>
> Claude Code then keeps only the host's values that restrict what Claude can do, with one gap to know about: unless you also set the `allowManaged*Only` locks, the host's permission allow rules and sandbox allowlists still apply. See [Restrict parent settings](/docs/en/claude-apps-gateway#restrict-parent-settings) for the locks.
>
> A [`policyHelper`](/docs/en/settings-reference#policyhelper) can turn parent merging off regardless of this key; its entry says when.
>
> Claude Code also applies these checks to parent-supplied values on their own:
>
> * When any admin source sets `allowManagedPermissionRulesOnly`, Claude Code drops [parent-supplied](/docs/en/claude-apps-gateway#restrict-parent-settings) permission allow rules and `additionalDirectories` as it reads them, even when a higher-priority source leaves the key unset. The key's effect on your own permission rules still comes only from the selected source, or from parent settings you've chosen to merge
> * A `forceLoginOrgUUID` or `allowedMcpServers` value in the highest-priority admin source blocks a parent-supplied one, and Claude Code enforces the admin value. A value in an admin source that isn't the highest-priority one neither applies nor blocks the parent's. Before v2.1.223, a value in any admin source blocked the parent's
> * An `availableModels` value follows the same rule as `forceLoginOrgUUID` and `allowedMcpServers`
>
> ### What a developer can change
>
> A developer's own settings files, `--settings` values, and project files never override a managed value; the [exceptions](/docs/en/settings#exceptions-to-managed-settings-precedence) only let a stricter lower-level value count. Four things sit outside that rule:
>
> * **The model for a session**: a managed `model` is a default, not a lock. `--model` and `ANTHROPIC_MODEL` still pick the model for that session, so deploy [`availableModels`](/docs/en/settings-reference#availablemodels) to restrict the choice.
> * **Local admin rights**: a developer who is an administrator on the machine can edit the managed source itself, which is why MDM tooling can redeploy the profile or file on a schedule and why the HKLM registry and the macOS managed preferences domain exist.
> * **The server-managed cache**: server-managed settings come from Anthropic's servers, and an edit to the local cache [lasts only until the next successful fetch](/docs/en/server-managed-settings#security-considerations).
> * **Other tools**: managed settings bind Claude Code only. A developer who calls the API from another tool isn't under them.
>
> ## Check that a policy is in force
>
> A developer reports that a policy isn't applying, or you want to confirm a rollout landed before pushing it to the fleet. Two commands on that machine answer it: `/status` shows which managed source Claude Code selected, and `claude doctor` lists what it dropped.
>
> ### Read the source in /status
>
> On the developer's machine, run `/status` inside Claude Code and read the `Setting sources` line. When a managed source is in effect, the line lists `Enterprise managed settings` with the source Claude Code selected in parentheses:
>
> * `(remote)`: server-managed settings from claude.ai or a gateway
> * `(plist)` or `(HKLM)`: an MDM or OS policy
> * `(file)`, `(drop-ins)`, or `(file + drop-ins)`: `managed-settings.json`, the drop-in directory, or both
> * `(HKCU)`: the user-writable registry fallback
> * `(parent process)`: an [embedding host](#let-an-embedding-host-add-policy) supplied restrictive settings
> * `(helper)`: a [`policyHelper`](/docs/en/settings-reference#policyhelper) configured by the selected MDM or file source
>
> When the policy isn't applying, the line tells you which of two problems you have:
>
> * **The line is missing**: Claude Code found no managed source that delivers a policy key. Check that the file sits at the path for the OS, that it's valid JSON, and that it contains a key other than `wslInheritsWindowsSettings`.
> * **The line names a source other than the one you deployed**: a higher-priority source is present and Claude Code ignored yours. [Which managed source Claude Code uses](#which-managed-source-claude-code-uses) gives the order.
>
> ### Find entries Claude Code dropped
>
> When a managed settings file, MDM profile, registry value, or server-managed payload fails schema validation, Claude Code first skips the individual entries it can repair, such as one invalid permission rule, with a warning for each, then drops any top-level key whose value still fails and keeps enforcing every remaining valid key. Claude Code is stricter with the `managedSettings` a [`policyHelper`](/docs/en/settings-reference#policyhelper) emits: it makes the same entry repairs, but any schema violation that survives fails the whole helper run, and at startup Claude Code refuses to start, the same as for a helper that exits non-zero. A managed settings file or drop-in file that isn't valid JSON contributes no settings at all; Claude Code reports it with the other validation errors and reads the remaining sources as usual.
>
> To find a dropped entry, look in one of three places:
>
> * Interactive sessions show a dialog at startup listing the invalid entries.
> * Non-interactive runs with `-p` print a summary to stderr.
> * [`claude doctor`](/docs/en/debug-your-config) lists each invalid entry with its source and field.
>
> #### Keys that fail closed
>
> A few enforcement keys aren't dropped when invalid. Claude Code enforces a stricter fallback until the value is fixed; the table shows what it enforces for each key:
>
> | Field                         | Behavior when present but invalid                                                                                                                                                                                                                |
> | :---------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `allowedMcpServers`           | Enforced as an empty allowlist, so no MCP servers are admitted until the value is fixed. An individual invalid entry is stripped and the valid subset is enforced.                                                                               |
> | `allowManagedHooksOnly`       | Treated as `true` until fixed: the [hook restrictions](/docs/en/settings-reference#allowmanagedhooksonly) apply and, unless `disableCommandPluginSources` is explicitly `false`, command-sourced plugins are disabled.                                |
> | `allowManagedMcpServersOnly`  | Treated as `true`.                                                                                                                                                                                                                               |
> | `disableCommandPluginSources` | Treated as `true`, so command-sourced plugins stay disabled until the value is fixed.                                                                                                                                                            |
> | `availableModels`             | Enforced as an empty allowlist until fixed, so only the Default model is available; a non-string entry is stripped and the valid subset enforced.                                                                                                |
> | `enforceAvailableModels`      | Treated as `true`.                                                                                                                                                                                                                               |
> | `forceLoginOrgUUID`           | No organization is permitted to log in until the value is fixed.                                                                                                                                                                                 |
> | `deniedMcpServers`            | An individual invalid entry is stripped and the valid subset is enforced. A wholly invalid value is dropped with a warning, since denying every server would block servers the policy never named.                                               |
> | `sandbox.credentials`         | A recoverable invalid entry is degraded to `mode: "deny"` with a warning; an unrecoverable one is stripped; valid entries stay enforced. See [invalid credential entries](/docs/en/settings-reference#invalid-credential-entries-in-managed-settings) |
>
> `requiredMinimumVersion` and `requiredMaximumVersion` fail open by design: an invalid value is dropped rather than enforced, so a bad policy push can't prevent Claude Code from starting.
>
> This tolerance applies only to managed settings. User, project, and local settings files remain strict: a file whose JSON or top-level shape fails validation is rejected as a whole and reported, and an individual entry that fails, such as a malformed permission rule, is skipped with a warning while the rest of the file applies.
>
> ## Keys only a managed source can set
>
> Claude Code reads the following keys only from a managed source; placing them in user or project settings files has no effect.
>
> Most of them are locks: the value a lock governs, such as permission rules or `sandbox.network.allowedDomains`, is an ordinary key that any level can set, and the lock tells Claude Code to honor only the managed value.
>
> The table covers the permission, plugin, and delivery controls. For any key not listed here, the Scope column of the [settings reference](/docs/en/settings-reference#all-settings) index says whether it's managed-only; the remaining managed-only keys there include the gateway login URL, version, browser, mobile-simulator, SSH host, Desktop local-session, sandbox binary path, and CLAUDE.md controls.
>
> | Setting                                                                                                               | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
> | :-------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | [`allowAllClaudeAiMcps`](/docs/en/settings-reference#allowallclaudeaimcps)                                                 | Load the claude.ai connectors alongside a deployed `managed-mcp.json` instead of suppressing them                                                                                                                                                                                                                                                                                                                                                                              |
> | [`allowedChannelPlugins`](/docs/en/settings-reference#allowedchannelplugins)                                               | Allowlist of channel plugins that may push messages. Replaces the default Anthropic allowlist when set. Requires `channelsEnabled: true`. See [Restrict which channel plugins can run](/docs/en/channels#restrict-which-channel-plugins-can-run)                                                                                                                                                                                                                                    |
> | [`allowManagedHooksOnly`](/docs/en/settings-reference#allowmanagedhooksonly)                                               | When `true`, restricts which hooks run; see [what runs under `allowManagedHooksOnly`](/docs/en/settings-reference#what-runs-under-allowmanagedhooksonly) for the full effect list                                                                                                                                                                                                                                                                                                   |
> | [`allowManagedMcpServersOnly`](/docs/en/settings-reference#allowmanagedmcpserversonly)                                     | When `true`, only `allowedMcpServers` from managed settings are respected. `deniedMcpServers` still merges from all sources. See [Managed MCP configuration](/docs/en/managed-mcp)                                                                                                                                                                                                                                                                                                  |
> | [`allowManagedPermissionRulesOnly`](/docs/en/settings-reference#allowmanagedpermissionrulesonly)                           | Only managed permission rules apply; the entry lists every source it ignores                                                                                                                                                                                                                                                                                                                                                                                                   |
> | [`blockedMarketplaces`](/docs/en/settings-reference#blockedmarketplaces)                                                   | Blocklist of marketplace sources. Blocked sources are checked before downloading, so they never touch the filesystem. See [managed marketplace restrictions](/docs/en/plugin-marketplaces#managed-marketplace-restrictions)                                                                                                                                                                                                                                                         |
> | [`channelsEnabled`](/docs/en/settings-reference#channelsenabled)                                                           | Allow [channels](/docs/en/channels) for the organization. See [enterprise controls](/docs/en/channels#enterprise-controls) for the default on each plan                                                                                                                                                                                                                                                                                                                                  |
> | [`disableCommandPluginSources`](/docs/en/settings-reference#disablecommandpluginsources)                                   | When `true`, blocks [`command` plugin sources](/docs/en/plugin-marketplaces#command-sources) entirely, so the marketplace-declared command never runs. Also blocks marketplace [`headersHelper` commands](/docs/en/plugin-marketplaces#authenticate-archive-downloads), except for a marketplace that managed settings themselves declare. When unset, follows `allowManagedHooksOnly`. Requires Claude Code v2.1.229 or later, and the `headersHelper` block requires v2.1.238 or later |
> | [`disableSideloadFlags`](/docs/en/settings-reference#disablesideloadflags)                                                 | Reject the `--plugin-dir`, `--plugin-url`, `--agents`, and `--mcp-config` flags at startup. In cloud sessions, Claude Code drops the MCP servers the server delivered through `--mcp-config`, other than in-process `type: "sdk"` entries, and starts the session. Requires Claude Code v2.1.193 or later                                                                                                                                                                      |
> | [`forceRemoteSettingsRefresh`](/docs/en/settings-reference#forceremotesettingsrefresh)                                     | When `true`, blocks CLI startup until remote managed settings are freshly fetched and exits if the fetch fails. See [fail-closed enforcement](/docs/en/server-managed-settings#enforce-fail-closed-startup)                                                                                                                                                                                                                                                                         |
> | [`parentSettingsBehavior`](/docs/en/settings-reference#parentsettingsbehavior)                                             | Whether host-supplied parent settings merge under the managed policy                                                                                                                                                                                                                                                                                                                                                                                                           |
> | [`pluginSuggestionMarketplaces`](/docs/en/settings-reference#pluginsuggestionmarketplaces)                                 | Marketplaces whose plugins Claude Code may suggest to users                                                                                                                                                                                                                                                                                                                                                                                                                    |
> | [`pluginTrustMessage`](/docs/en/settings-reference#plugintrustmessage)                                                     | Custom message appended to the plugin trust warning shown before installation                                                                                                                                                                                                                                                                                                                                                                                                  |
> | [`policyHelper`](/docs/en/settings-reference#policyhelper)                                                                 | Executable that computes managed settings at startup; see [Compute managed settings with a policy helper](/docs/en/settings-reference#policyhelper)                                                                                                                                                                                                                                                                                                                                 |
> | [`sandbox.filesystem.allowManagedReadPathsOnly`](/docs/en/settings-reference#sandbox-filesystem-allowmanagedreadpathsonly) | When `true`, only `filesystem.allowRead` paths from managed settings are respected. `denyRead` still merges from all sources                                                                                                                                                                                                                                                                                                                                                   |
> | [`sandbox.network.allowManagedDomainsOnly`](/docs/en/settings-reference#sandbox-network-allowmanageddomainsonly)           |### Source: Server-managed settings
>
> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Configure server-managed settings
>
> > Centrally configure Claude Code for your organization through server-delivered settings, without requiring device management infrastructure.
>
> Server-managed settings let organization Owners centrally configure Claude Code from [**Admin Settings > Claude Code > Managed settings**](https://claude.ai/admin-settings/claude-code) in the claude.ai console. Claude Code clients fetch these settings automatically when users authenticate with a Team or Enterprise OAuth login, an OAuth token supplied through `CLAUDE_CODE_OAUTH_TOKEN`, or a directly configured API key, on platforms where server-managed delivery is supported. See [Platform availability](#platform-availability).
>
>   Server-managed settings are available for [Claude for Teams](https://claude.com/pricing?utm_source=claude_code\&utm_medium=docs\&utm_content=server_settings_teams#team-&-enterprise) and [Claude for Enterprise](https://anthropic.com/contact-sales?utm_source=claude_code\&utm_medium=docs\&utm_content=server_settings_enterprise) customers.
>
> ## Requirements
>
> To use server-managed settings, you need:
>
> * Claude for Teams or Claude for Enterprise plan
> * The Owner or Primary Owner role in your Claude organization, to view and edit the configuration
> * Network access to `api.anthropic.com`
>
> ## Choose between server-managed and endpoint-managed settings
>
> Claude Code supports two approaches for centralized configuration. Server-managed settings deliver configuration from Anthropic's servers. [Endpoint-managed settings](/docs/en/managed-settings#delivery-mechanisms) are deployed directly to devices through native OS policies (macOS managed preferences, Windows registry) or managed settings files.
>
> | Approach                                                                  | Best for                                                 | Security model                                                                                                |
> | :------------------------------------------------------------------------ | :------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------ |
> | **Server-managed settings**                                               | Organizations without MDM, or users on unmanaged devices | Settings that Claude Code fetches from Anthropic's servers at startup and refreshes hourly during the session |
> | **[Endpoint-managed settings](/docs/en/managed-settings#delivery-mechanisms)** | Organizations with MDM or endpoint management            | Settings deployed to devices via MDM configuration profiles, registry policies, or managed settings files     |
>
> If your devices are enrolled in an MDM or endpoint management solution, endpoint-managed settings provide stronger security guarantees because the settings file can be protected from user modification at the OS level. Endpoint-managed settings don't reach [cloud sessions](/docs/en/model-config#surface-coverage) in Anthropic-hosted environments, so organizations using Claude Code on the web should configure server-managed settings as well. Sessions in a [self-hosted environment](/docs/en/self-hosted-environments) read the managed settings file in the runner image, but only when server-managed settings deliver no keys, per the [settings precedence](#settings-precedence) below and its [per-key exceptions](#per-key-exceptions-across-managed-sources).
>
> ## Configure server-managed settings
>
>   **Open the admin console**: In the claude.ai console, go to [**Admin Settings > Claude Code > Managed settings**](https://claude.ai/admin-settings/claude-code).
>
>     If the link redirects you to a different Admin Settings page instead of the Claude Code page, your account doesn't have the required role. Admin and other non-Owner roles can't view or edit managed settings, so ask an Owner or Primary Owner in your organization to make the change. See [Access control](#access-control).
>
>   **Define your settings**: Add your configuration as JSON. All [settings available in `settings.json`](/docs/en/settings-reference#all-settings) are supported except those restricted to OS-level policy delivery; see [Current limitations](#current-limitations) for that short list. This includes [hooks](/docs/en/hooks), [environment variables](/docs/en/env-vars), and [managed-only settings](/docs/en/managed-settings#managed-only-settings) like `allowManagedPermissionRulesOnly`.
>
>     This example enforces a permission deny list, prevents users from bypassing permissions, and restricts permission rules to those defined in managed settings:
>
>     ```json
>     {
>       "permissions": {
>         "deny": [
>           "Bash(curl *)",
>           "Read(./.env)",
>           "Read(./.env.*)",
>           "Read(./secrets/**)"
>         ],
>         "disableBypassPermissionsMode": "disable"
>       },
>       "allowManagedPermissionRulesOnly": true
>     }
>     ```
>
>     Hooks use the same format as in `settings.json`.
>
>     This example runs an audit script after every file edit across the organization:
>
>     ```json
>     {
>       "hooks": {
>         "PostToolUse": [
>           {
>             "matcher": "Edit|Write",
>             "hooks": [
>               { "type": "command", "command": "/usr/local/bin/audit-edit.sh" }
>             ]
>           }
>         ]
>       }
>     }
>     ```
>
>     Because hooks execute shell commands, users in interactive sessions see a [security approval dialog](#security-approval-dialogs) before Claude Code applies them.
>
>     To configure the [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) classifier so it knows which repos, buckets, and domains your organization trusts, deliver an `autoMode` block the same way; see [Configure auto mode](/docs/en/auto-mode-config) for how the `autoMode` entries affect what the classifier blocks and important warnings about the `environment`, `allow`, `soft_deny`, and `hard_deny` fields.
>
>   **Save and deploy**: Save your changes. Claude Code clients receive the updated settings on their next startup or hourly polling cycle.
>
> ### Verify settings delivery
>
> To confirm that settings are being applied, ask a user to restart Claude Code. If the configuration includes settings that trigger the [security approval dialog](#security-approval-dialogs), the user sees a prompt describing the managed settings the next time Claude Code fetches them: at the next start, or within an hour in a running interactive session. You can also verify that managed permission rules are active by having a user run `/permissions` to view their effective permission rules.
>
> ### Access control
>
> The following roles can manage server-managed settings:
>
> * **Primary Owner**
> * **Owner**
>
> Restrict access to trusted personnel, as settings changes apply to all users in the organization.
>
> ### Managed-only settings
>
> Most [settings keys](/docs/en/settings-reference#all-settings) work in any scope. A handful of keys are only read from managed settings and have no effect when placed in user or project settings files. See [managed-only settings](/docs/en/managed-settings#managed-only-settings) for the permission and plugin controls, or read the Scope column of the [All settings](/docs/en/settings-reference#all-settings) index for the full set.
>
> ### Current limitations
>
> Server-managed settings have the following limitations:
>
> * Settings apply uniformly to all users in the organization. Per-group configurations are not yet supported.
> * A [`managed-mcp.json`](/docs/en/managed-mcp) file can't be distributed through server-managed settings. Deliver the `allowedMcpServers` and `deniedMcpServers` policy keys there instead. Claude Code reads a `managed-mcp.json` deployed at its [system path](/docs/en/managed-mcp#exclusive-control-with-managed-mcp-json) separately from the managed settings tier, so the file still applies when server-managed settings are in effect.
> * Settings restricted to OS-level policy sources, such as `policyHelper` and `wslInheritsWindowsSettings`, aren't honored. Deploy them through MDM or a system `managed-settings.json` file instead. A `policyHelper` deployed that way runs only when its source is the one selected under [precedence within the managed tier](/docs/en/managed-settings#precedence-within-the-managed-tier).
>
> ## Settings delivery
>
> ### Settings precedence
>
> Server-managed settings and [endpoint-managed settings](/docs/en/managed-settings#delivery-mechanisms) both occupy the highest tier in the Claude Code [settings hierarchy](/docs/en/settings#settings-precedence). No other settings level can override them, including command line arguments, apart from the [exceptions to managed settings precedence](/docs/en/settings#exceptions-to-managed-settings-precedence).
>
> Within the managed tier, Claude Code uses the first source that delivers at least one policy key, checking server-managed settings first and then endpoint-managed settings; [Precedence within the managed tier](/docs/en/managed-settings#precedence-within-the-managed-tier) has the full ranking and the carve-out for `wslInheritsWindowsSettings`. Apart from the [exception keys covered next](#per-key-exceptions-across-managed-sources), sources don't merge: if server-managed settings deliver a policy key, other endpoint-managed settings are ignored. If server-managed settings deliver nothing, endpoint-managed settings apply.
>
> If the selected source is an MDM policy or managed settings file whose [`policyHelper`](/docs/en/settings-reference#policyhelper) supplies managed settings, the helper's output replaces that source as the only managed configuration for the run. Claude Code doesn't consult a `policyHelper` configured in MDM or file-based settings while server-managed settings deliver a policy key.
>
> If you clear your server-managed configuration in the admin console with the intent of falling back to an endpoint-managed plist or registry policy, be aware that [cached settings](#fetch-and-caching-behavior) persist on client machines until the next successful fetch, and the keys that [apply only at the next launch](#fetch-and-caching-behavior), such as `model`, stay in effect until each client relaunches. Run `/status` to see which managed source is active.
>
> ### Per-key exceptions across managed sources
>
> Two kinds of keys are exceptions to the no-merge rule:
>
> * **Cross-source lock keys**: a small set of keys, such as the sandbox allowlist locks, [listed on the managed settings page](/docs/en/managed-settings#precedence-within-the-managed-tier). They are honored when any admin-controlled managed source sets them; the user-writable HKCU registry tier is excluded, and when a [`policyHelper`](/docs/en/settings-reference#policyhelper) supplies managed settings, its output is the only source these checks read. The startup [`forceRemoteSettingsRefresh`](#enforce-fail-closed-startup) check runs before the helper and reads any admin source.
> * **The `env` block**: apart from the telemetry unit and routing variables paired with a credential key, both covered below, it merges per key across the admin-controlled sources. For each environment variable, the highest-priority source defining it wins, and lower admin sources fill in variables the higher sources leave unset. An endpoint-managed `env` entry therefore applies whenever the server-managed configuration leaves that variable unset, or while a cached server value for it is [withheld pending server confirmation](#fetch-and-caching-behavior). Requires Claude Code v2.1.223 or later. Before v2.1.223, Claude Code applies the winning source's whole `env` block only.
>   * **Telemetry unit**: the `OTEL_EXPORTER_OTLP_*` exporter keys, the `OTEL_LOG_*` content-capture toggles, `OTEL_LOGS_EXPORTER`, and the beta tracing variables `ENABLE_BETA_TRACING_DETAILED` and `BETA_TRACING_ENDPOINT` follow the highest source that sets any of them as a unit. A source that delivers the `otelHeadersHelper` credential key claims the unit too, but lands these variables only when it is the winning source: a non-winning source that delivers the key contributes none of them and still blocks lower sources from filling them in. Either way, an exporter endpoint from one source can never pair with credentials from another.
>   * **Credential-paired routing**: a source that pairs routing variables with a winner-only credential key, such as `apiKeyHelper` or `otelHeadersHelper`, contributes those routing variables only when it wins the slot.
>
> ### Fetch and caching behavior
>
> Claude Code fetches settings from Anthropic's servers at startup and polls for updates hourly during active sessions.
>
> A client signed in through a [Claude apps gateway](#platform-availability) fetches its settings from the gateway and waits for that fetch before the session starts, so the asynchronous fetch in the lists below doesn't apply to it. [Enforce fail-closed startup](#enforce-fail-closed-startup) covers what happens when that fetch fails.
>
> **First launch without cached settings:**
>
> * Claude Code fetches settings asynchronously
> * If the fetch fails, Claude Code continues without server-managed settings; endpoint-managed settings still apply. If a managed source sets [`forceRemoteSettingsRefresh`](#enforce-fail-closed-startup), Claude Code exits instead
> * Until the fetch completes, Claude Code doesn't yet enforce the server-managed settings
>
> **Subsequent launches with cached settings:**
>
> * Cached settings apply immediately at startup, except for the withheld environment variables described below
> * Claude Code fetches fresh settings in the background
> * Cached settings persist through network failures. The withheld environment variables remain withheld until a fetch succeeds
>
> Claude Code withholds several categories of variables in the cached `env` block until the server confirms the payload for the session. This keeps a cached proxy, certificate authority, endpoint, or credential value from redirecting, intercepting, or re-authenticating the settings fetch that confirms the payload. The hardening applies only to the server-fetched settings cache: [endpoint-managed settings](/docs/en/managed-settings#delivery-mechanisms) deployed through MDM or `managed-settings.json` are unaffected. The withholding requires Claude Code v2.1.198 or later; before v2.1.198, the whole cached `env` block applies at startup. The withheld categories include:
>
> * Proxy and TLS configuration, such as `HTTPS_PROXY`, `NODE_EXTRA_CA_CERTS`, and the mTLS client certificate variables `CLAUDE_CODE_CLIENT_CERT` and `CLAUDE_CODE_CLIENT_KEY`
> * API routing and provider selection, including `ANTHROPIC_BASE_URL`, the provider selection variables such as `CLAUDE_CODE_USE_BEDROCK` and `CLAUDE_CODE_USE_VERTEX`, and the provider endpoint URLs such as `ANTHROPIC_BEDROCK_BASE_URL`
> * Authentication credentials, such as `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, and `CLAUDE_CODE_OAUTH_TOKEN`
> * The configuration-directory selector `CLAUDE_CONFIG_DIR`
> * Credential-source and configuration-directory selectors, in Claude Code v2.1.223 or later: the Workload Identity Federation variables such as `ANTHROPIC_FEDERATION_RULE_ID` and `ANTHROPIC_IDENTITY_TOKEN`, the profile and configuration-directory selectors `ANTHROPIC_PROFILE` and `ANTHROPIC_CONFIG_DIR`, and the operating-system directory variables `HOME`, `XDG_CONFIG_HOME`, `APPDATA`, and `USERPROFILE`
>
> Claude Code reads the Workload Identity Federation variables and the `ANTHROPIC_PROFILE` and `ANTHROPIC_CONFIG_DIR` selectors only at startup, so a server-delivered value for them doesn't switch the session's credential source even after the fetch succeeds. To deliver those selectors on Claude Code v2.1.223 or later, use [endpoint-managed settings](/docs/en/managed-settings#delivery-mechanisms) such as MDM or `managed-settings.json`. For `CLAUDE_CONFIG_DIR` and the operating-system directory variables, the withholding itself is the protection: the cached value stays out of the environment until the server confirms the payload.
>
> Every other key in the cached `env` block, such as telemetry and OpenTelemetry configuration, applies at startup as before. Once the server confirms the payload, and you approve it if it needs [security approval](#security-approval-dialogs), the withheld variables apply for the rest of the session; the startup-only selectors covered above reach the environment but don't switch the running session's credential source.
>
> If your organization needs a proxy to reach `api.anthropic.com`, the withholding only affects the server-delivered `env` block itself: a proxy set in an [endpoint-managed](/docs/en/managed-settings#delivery-mechanisms) `env` block through MDM or `managed-settings.json`, in the shell environment, or in [user settings](/docs/en/settings#where-settings-live) reaches the settings fetch. The endpoint-managed source requires Claude Code v2.1.223 or later: the cached server-managed proxy value is withheld until the fetch confirms it, so the endpoint-managed value fills in per key and reaches the fetch itself. Before v2.1.223, use the shell environment or user settings so the proxy applies alongside a cached server payload. The first launch has no cache, so an endpoint-managed source, the shell environment, or user settings is still required for the initial fetch.
>
> Claude Code applies most settings updates to running sessions without a restart. Some updates apply only at the next launch, including OpenTelemetry exporter configuration, the `model` key, and the removal of a variable from the `env` block.
>
> ### Invalid entries in delivered settings
>
> Delivered payloads parse tolerantly with the same rules as the other managed sources. When part of a payload fails schema validation, Claude Code surfaces a validation error and applies every remaining valid setting; [Invalid entries in managed settings](/docs/en/managed-settings#invalid-entries-in-managed-settings) says what it drops and which keys fall back to a stricter value. Requires Claude Code v2.1.169 or later.
>
> Server-managed delivery adds these behaviors:
>
> * The cache at `~/.claude/remote-settings.json` stores the salvaged payload with invalid entries removed. The raw invalid payload is never persisted.
> * When no field in the payload can be salvaged, Claude Code rejects the payload, keeps the last-accepted cached settings, and writes `Remote settings: Settings validation failed - no fields could be salvaged` to the debug log. With `forceRemoteSettingsRefresh` set, the CLI exits instead.
> * The [security approval dialog](#security-approval-dialogs) evaluates the salvaged payload, so a stripped invalid entry is never presented for approval and never executes.
>
> To debug delivery issues, run `claude --debug-file <path>` and search the log for `Remote settings`. Validate a payload change with `claude doctor` on a test machine before rolling it out to the organization.
>
> ### Enforce fail-closed startup
>
> By default, if the remote settings fetch fails at startup, the CLI continues with the settings cached from the last successful fetch, except for the [withheld environment variables](#fetch-and-caching-behavior), or without server-managed settings on a machine that has never fetched them. With no cache, Claude Code still applies any [endpoint-managed settings](/docs/en/managed-settings#delivery-mechanisms) on the device. For environments where either window is unacceptable, set `forceRemoteSettingsRefresh: true` in your managed settings. Clients signed in through a [Claude apps gateway](#platform-availability) exit when the startup fetch fails, without this setting.
>
> When this setting is active in a session that fetches server-managed settings, the CLI blocks at startup until remote settings are freshly fetched. If the fetch fails, the CLI exits rather than proceeding without the policy. This setting self-perpetuates: once delivered from the server, it is also cached locally so that subsequent startups enforce the same behavior even before the first successful fetch of a new session. A session that [doesn't fetch server-managed settings](#platform-availability) starts without waiting.
>
> To enable this, add the key to your managed settings configuration:
>
> ```json
> {
>   "forceRemoteSettingsRefresh": true
> }
> ```
>
> You can also set this key in an [endpoint-managed](/docs/en/managed-settings#delivery-mechanisms) MDM profile or system `managed-settings.json` file to enforce fail-closed behavior on first launch, before any server payload has been delivered. In Claude Code v2.1.191 or later, this flag is an exception to the [precedence rule](#settings-precedence) above: it is honored when set in any admin-controlled managed source even if a cached server-managed payload is also present, so an MDM-delivered value is not ignored when server-managed settings exist. When a [`policyHelper`](/docs/en/settings-reference#policyhelper) supplies managed settings, its output replaces every other managed source for the keys Claude Code reads after startup; the startup fail-closed check itself reads this key from any admin-controlled source before the helper runs. The entry says which sources Claude Code reads the helper from and when it runs.
>
> The settings fetch also sends a `Cache-Control: no-cache` header so intermediate HTTP proxies don't serve a stale response.
>
> Before enabling this setting, ensure your network policies allow connectivity to `api.anthropic.com`. If that endpoint is unreachable, the CLI exits at startup and users cannot start Claude Code.
>
> The `claude auth` subcommands such as `claude auth login` are exempt from this check and from the gateway startup exit, so users can re-authenticate when expired credentials are the reason the settings fetch fails.
>
> ### Security approval dialogs
>
> Certain settings that could pose security risks require explicit user approval before Claude Code applies them in an interactive session:
>
> * **Shell command settings**: settings that execute shell commands, such as `apiKeyHelper`, `statusLine`, and `otelHeadersHelper`
> * **Sandbox binary settings**: `sandbox.bwrapPath`, `sandbox.socatPath`, and `sandbox.ripgrep`. Each of these settings points at an executable, and Claude Code runs that executable
> * **Custom environment variables**: delivered `env` variables that require the user's approval, such as proxy and base-URL variables; see [Environment variables and the approval dialog](#environment-variables-and-the-approval-dialog)
> * **Hook configurations**: any hook definition
> * **Managed CLAUDE.md content**: a `claudeMd` value delivered through managed settings
>
> When these settings are present, users see a security dialog explaining what is being configured. Users must approve to proceed. If a user rejects the settings, Claude Code exits.
>
> #### Approval memory
>
> Claude Code records your approval in your configuration directory, `~/.claude` unless you set [`CLAUDE_CONFIG_DIR`](/docs/en/env-vars). What it records depends on the credential the settings fetch uses:
>
> * **A claude.ai login saved by `/login` or `claude auth login`**: one approval per organization, held by the account that approved most recently.
> * **Any other credential**, such as a [Claude apps gateway](/docs/en/claude-apps-gateway), an API key, or `CLAUDE_CODE_OAUTH_TOKEN`: one approval for the delivered settings, kept with the cached copy of the settings in that configuration directory. Claude Code shows the dialog again when the settings that require approval change, and after you run `/logout` or `claude auth logout`, which delete the cached copy.
>
> With a saved claude.ai login:
>
> * If you sign out and back in, or switch to another organization and later return, Claude Code doesn't show the dialog again while those settings are unchanged, unless another account approved them for that organization in the same configuration directory in between.
> * If you sign in to the same organization with a different account, Claude Code shows the dialog again even when the settings are unchanged. That account's approval replaces the previous one, so when you switch back, Claude Code shows the dialog once more.
>
> If an interactive session can't show the dialog, Claude Code doesn't apply the delivered settings and keeps the last-approved settings; the dialog appears in the next session that can show it. Requires Claude Code v2.1.211 or later.
>
> If an error closes the dialog before you answer, Claude Code doesn't apply the delivered settings and keeps the last-approved settings; Claude Code shows the dialog again in the next session that can show it.
>
>   A non-interactive run, such as `claude -p` or an Agent SDK session, can't show the dialog. When the delivered settings would require approval, Claude Code applies them for that run only: it doesn't record them as approved or write them to the [local cache](#fetch-and-caching-behavior), and the next interactive session shows the dialog. Until a user approves in an interactive session, each non-interactive run fetches the settings again at startup. Before v2.1.207, a non-interactive run saved the settings as approved, so later interactive sessions never showed the dialog for them.
>
> #### Environment variables and the approval dialog
>
> Claude Code applies some delivered `env` variables without showing the user the approval dialog, including:
>
> * Feature and command toggles
> * Model selection and behavior settings, such as `ANTHROPIC_MODEL`, `DISABLE_PROMPT_CACHING`, and `CLAUDE_CODE_EFFORT_LEVEL`
> * Context window and compaction settings, such as `DISABLE_AUTO_COMPACT`
> * Terminal UI and accessibility options
> * Numeric limits, budgets, and timeouts
>
> Other delivered variables can require the user's approval before they take effect; a non-empty proxy, base-URL, or `OTEL_EXPORTER_OTLP_ENDPOINT` value always does. When a delivered variable needs approval, the dialog names it, so the user sees exactly what the policy is asking to set. Before v2.1.218, Claude Code applied fewer variables without asking the user, so settings such as `DISABLE_AUTO_COMPACT` triggered the dialog at any non-empty value.
>
> Claude Code decides whether four privacy toggles need approval by the delivered value rather than by the variable name: `CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC`, `DISABLE_ERROR_REPORTING`, `DISABLE_TELEMETRY`, and `DO_NOT_TRACK`. A truthy value such as `1` or `true` only turns tracking, reporting, or other nonessential traffic off, so Claude Code applies it without asking the user. For any other non-empty value, Claude Code shows the dialog. Before v2.1.218, all of them except `DO_NOT_TRACK` applied without approval at any value, and `DO_NOT_TRACK` triggered the dialog at any non-empty value.
>
> ## Platform availability
>
> Server-managed settings require a direct connection to `api.anthropic.com`, and delivery requires the session to authenticate with a Team or Enterprise OAuth login, an OAuth token supplied through `CLAUDE_CODE_OAUTH_TOKEN`, or a directly configured API key. Keys returned by an [`apiKeyHelper`](/docs/en/settings-reference#apikeyhelper) script don't trigger the settings fetch.
>
> Server-managed settings are not available when using third-party model providers:
>
> * Amazon Bedrock
> * Google Cloud's Agent Platform
> * Microsoft Foundry
> * [Claude Platform on AWS](/docs/en/claude-platform-on-aws)
> * Custom API endpoints via `ANTHROPIC_BASE_URL` or third-party [LLM gateways](/docs/en/llm-gateway)
>
> If you export a `CLAUDE_CODE_USE_*` provider variable or a non-default `ANTHROPIC_BASE_URL` in your shell, Claude Code skips the settings fetch for your sessions. You can't clear the export with a server-managed `env` block, because the block arrives through the fetch that the export prevents. An [endpoint-managed settings](/docs/en/managed-settings#delivery-mechanisms) `env` block doesn't restore the fetch either: Claude Code checks eligibility before it applies managed `env` blocks, so the override changes the session's provider selection but the fetch stays skipped.
>
> To restore server-managed delivery, remove the export from your shell, or set the variable to `""` in your user settings `env` block, which applies before the eligibility check. To enforce policy without relying on users to change their shells, deliver the settings through the endpoint-managed channel instead.
>
> For Amazon Bedrock, Google Cloud's Agent Platform, Microsoft Foundry, and [Claude Platform on AWS](/docs/en/claude-platform-on-aws) deployments, a self-hosted [Claude apps gateway](/docs/en/claude-apps-gateway) provides the equivalent remote managed-settings delivery: gateway-signed-in clients fetch managed settings from the gateway instead of `api.anthropic.com`. The failure semantics differ at startup: a gateway client that can't reach the gateway exits with an error instead of falling back to cached settings, while the hourly background refresh is fail-open on both channels.
>
> ## Audit logging
>
> Audit log events for settings changes are available through the compliance API or audit log export. Contact your Anthropic account team for access.
>
> Audit events include the type of action performed, the account and device that performed the action, and references to the previous and new values.
>
> ## Security considerations
>
> Server-managed settings provide centralized policy enforcement, but they operate as a client-side control, not a security boundary. On unmanaged devices, a user doesn't need admin or sudo access to bypass them.
>
> | Scenario                                                               | Behavior