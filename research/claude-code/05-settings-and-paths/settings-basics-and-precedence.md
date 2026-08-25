---
primary_sources:
  - id: T1-SETTINGS
    title: "Claude Code settings"
    url: "https://code.claude.com/docs/en/settings.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Settings basics and precedence

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Claude Code settings

> # Claude Code settings
>
> > Change Claude Code settings, pick the scope a key belongs in, verify the change, and learn which value Claude Code uses when a key is set in several places.
>
> ## Settings files and who they affect
>
> Claude Code reads settings from four files, and an organization can also deliver managed settings from the claude.ai console. Each source has a scope: the set of people and projects a setting saved in it applies to, whether that's just you, everyone in a project, or everyone in your organization.
>
> | Scope          | File                                                                                          | Who it affects                                                                                                                                                       | Use it for                                                                         |
> | :------------- | :-------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------- |
> | User           | `~/.claude/settings.json`                                                                     | You, in every project on this machine                                                                                                                                | Personal preferences: theme, editor mode, default model, your own permission rules |
> | Shared project | `.claude/settings.json`                                                                       | Everyone who starts Claude Code in the folder that contains it. In a git repository, commit it so teammates get it                                                   | Team permissions, hooks, plugins, and the environment variables the project needs  |
> | Project local  | `.claude/settings.local.json`                                                                 | You, in this one project only. Claude Code keeps it out of git when it creates the file; if you create it by hand, add it to `.gitignore` yourself                   | Personal overrides for one project, and testing before you share                   |
> | Managed        | `managed-settings.json` and other [managed sources](/docs/en/managed-settings#delivery-mechanisms) | Everyone your organization deploys it to; nothing you set overrides it, apart from a few [security-sensitive exceptions](#exceptions-to-managed-settings-precedence) | Security policy and compliance requirements                                        |
>
> In the File column, `~/.claude` is the `.claude` folder in your home directory, and a bare `.claude` is the `.claude` folder inside the project you start Claude Code in.
>
>
> ### Compare the scope of each settings file
>
> Suppose you have three projects on your machine, `website/`, `api/`, and `acme-app/`, a teammate has their own clone of `acme-app/`, and you start a [cloud session](#settings-in-cloud-sessions) on `acme-app/`.
>
> The graphic below shows which of those folders a setting applies in when you start Claude Code from them. Click a settings file to see the folders it reaches.
>
>
> * **`~/.claude/settings.json`**: every project on your machine, and nothing on your teammate's or in the cloud session
> * **`acme-app/.claude/settings.json`**: your `acme-app/`. It reaches your teammate's clone and the cloud session only if you commit the file to version control; until you do, it's a file on your disk like any other and nobody else has it
> * **`acme-app/.claude/settings.local.json`**: your `acme-app/` only. Claude Code adds it to your global git excludes the first time it writes the file, so it stays out of your commits; if you create the file by hand, [add it to `.gitignore` yourself](#keep-personal-settings-out-of-a-repository)
> * **Managed settings**, whether a `managed-settings.json` file, an MDM policy, or [server-managed settings](/docs/en/server-managed-settings) from the claude.ai console: every project on every machine your organization deploys it to, or that you sign in to with your organization account. Only server-managed settings reach the cloud session
>
>
> ### Find or create your settings files
>
> Installing Claude Code doesn't create any settings file. If your machine or project already has one, it came from one of these sources:
>
> * **Managed**: your organization deploys it. You don't create or edit it.
> * **Shared project**: a project that already uses Claude Code may have one committed. If not, create it at `.claude/settings.json` in the project folder.
> * **User** and **Project local**: create them yourself, or let Claude Code create them. It writes `~/.claude/settings.json` the first time you change an option in the `/config` menu that it stores in user settings, such as the theme, and `.claude/settings.local.json` the first time you give a standing approval on a permission prompt, such as "Yes, and don't ask again" for a Bash command. A few `/config` options, including **Show tips**, save to `.claude/settings.local.json` instead of the user file.
>
>
>   On Windows, `~/.claude` means `%USERPROFILE%\.claude`. To keep the home-directory files somewhere else, set [`CLAUDE_CONFIG_DIR`](/docs/en/env-vars); Claude Code then stores your settings, session history, and plugins there instead.
>
>
> Claude Code also keeps a fifth file, [`~/.claude.json`](/docs/en/claude-directory#ce-claude-json), that it writes for itself; you don't need to edit it. It holds your sign-in session, [MCP server](/docs/en/mcp) configurations, per-project state such as trust decisions, and the [global config keys](/docs/en/settings-reference#global-config-settings) that `/config` writes for you.
>
> ### Share settings with your team
>
> Commit `.claude/settings.json` so everyone who clones the repository gets the same permissions, hooks, telemetry, and plugins. Each teammate can still override it for themselves in their own `.claude/settings.local.json`, so personal exceptions don't need a commit. For a complete team file, see [a team's shared settings](/docs/en/settings-example#a-teams-shared-settings).
>
> Some of what you commit waits until each teammate [trusts the folder](/docs/en/permissions#project-allow-rules-and-workspace-trust), and a few keys never take effect from a repository file; [Troubleshoot a setting that doesn't apply](#common-cases) covers both.
>
>
> ### Keep personal settings out of a repository
>
> To change a setting for yourself in one project without changing it for your teammates, save it in `.claude/settings.local.json` inside the project. Claude Code applies that file over the committed `.claude/settings.json`, so if your team's file sets `"model": "claude-sonnet-5"` and you want Opus, put `"model": "claude-opus-4-8"` in your local file and only your sessions change.
>
> Three things to know about the local file:
>
> * **Claude Code writes it too.** When Claude asks permission to run a Bash command and you choose "Yes, and don't ask again", Claude Code saves that [permission approval](/docs/en/permissions#permission-system) here as an `allow` rule.
> * **You don't need to gitignore it yourself, unless you created it by hand.** The first time Claude Code writes the file in a git repository that doesn't already ignore it, it adds `**/.claude/settings.local.json` to your global git excludes file, so the file stays out of your commits in every repository. That file is `core.excludesFile` when your global git config sets it to an absolute or `~`-prefixed path; otherwise it's `$XDG_CONFIG_HOME/git/ignore`, or `~/.config/git/ignore` when `XDG_CONFIG_HOME` is unset. If you created the file by hand and Claude Code hasn't written to it yet, add it to `.gitignore` yourself.
> * **Its allow rules don't wait for trust while the file stays untracked.** Because the file is yours and not the repository's, Claude Code applies its `allow` rules without the [workspace trust](/docs/en/permissions#project-allow-rules-and-workspace-trust) step it requires for the committed file. If the file is tracked by git, the trust step applies to it too; see [When your local settings file needs trust](/docs/en/permissions#when-your-local-settings-file-needs-trust).
>
>
> #### Where Claude Code keeps the local file in a git repository
>
> When Claude asks permission to run a Bash command and you choose "Yes, and don't ask again", Claude Code saves that approval as an allow rule in `.claude/settings.local.json`. If you started Claude Code in a subdirectory or a [worktree](/docs/en/worktrees) of a git repository, it reads and writes that file at the repository root, so the approval applies across the whole repository. The shared `.claude/settings.json` doesn't move: Claude Code reads it only from the folder you start in, so start at the repository root to pick up a committed file there. Two details follow from the root location:
>
> * **When the file stays in the starting directory instead**: outside a git repository, when the repository root is your home directory, on Windows, or when the repository root or its `.git` or `.claude` entry isn't owned by your user.
> * **Paths in the file still resolve from where you started**: a permission rule that starts with `/` or a relative sandbox path keeps covering the directory you started Claude Code in, not the repository root.
>
> Before v2.1.211, Claude Code kept the file in the starting directory. It still reads a file an earlier version left there alongside the root file; where both set the same key, the root's value applies, and permission rules from both files apply. The Agent SDK's [`resolveSettings()`](/docs/en/agent-sdk/typescript#resolvesettings) helper always reads the file from the starting directory.
>
>
> ### Check what your organization enforces
>
> If your organization manages Claude Code, some settings are decided for you and nothing you put in your own files changes them. To see which, run `/status`: the `Setting sources` line names the managed source that applies to you. Managed settings apply wherever Claude Code runs on this machine; [What a developer can change](/docs/en/managed-settings#what-a-developer-can-change) covers local admin rights and tools other than Claude Code.
>
> Managed settings reach you through the [delivery mechanisms](/docs/en/managed-settings#delivery-mechanisms) on the managed settings page, most commonly:
>
> * [Server-managed settings](/docs/en/server-managed-settings), which Claude Code fetches from the claude.ai admin console or a self-hosted [Claude apps gateway](/docs/en/claude-apps-gateway)
> * MDM or OS-level policies, and `managed-settings.json` files in a system directory
> * An embedding host such as Claude Desktop, through the SDK `managedSettings` option; see [Control policy from an embedding host](/docs/en/managed-settings#parent-settings-from-embedding-hosts)
>
> If you're the administrator, [Set up Claude Code for your organization](/docs/en/admin-setup) walks through choosing what to enforce, and [Deploy managed settings](/docs/en/managed-settings) covers delivery and how to confirm a policy is in force.
>
> ## Change a setting
>
> You can change a setting from the `/config` menu, by editing a settings file, or for one session from the command line.
>
>
> Claude Code's system prompt isn't published. To give Claude standing instructions, use [`CLAUDE.md` files](/docs/en/memory) or the `--append-system-prompt` flag.
>
> ### Use the /config menu
>
> Run `/config` inside Claude Code and open the **Config** tab. It lists a short set of personal options such as theme, editor mode, and verbose output, not every settings key. Select an option to change it; Claude Code saves it for you:
>
> * **Most options**: `~/.claude/settings.json`
> * **A few options, such as Show tips**: `.claude/settings.local.json`
> * **The [global config options](/docs/en/settings-reference#global-config-settings)**: `~/.claude.json`
>
> To set one option without the menu, pass `key=value`, such as `/config verbose=true`.
>
>
>   `/config` is part of the terminal interface. The [VS Code](/docs/en/vs-code) chat panel and the [desktop app](/docs/en/desktop) don't open it; change settings there by editing a settings file or through those apps' own settings.
>
>
> ### Edit a settings file
>
> Open the settings file for the scope you want in your editor and add or change a key. Settings files are strict JSON: a `//` comment or a trailing comma is a syntax error, and Claude Code reports the file as a [Settings Error](#fix-a-broken-settings-file) at the next start. For example, to let Claude Code run your lint and test commands without asking and stop it reading `.env` files, add this to `~/.claude/settings.json`:
>
> ```json ~/.claude/settings.json
> {
>   "$schema": "https://json.schemastore.org/claude-code-settings.json",
>   "permissions": {
>     "allow": [
>       "Bash(npm run lint)",
>       "Bash(npm run test *)"
>     ],
>     "deny": [
>       "Read(./.env)",
>       "Read(./.env.*)"
>     ]
>   }
> }
> ```
>
> Each entry under `permissions` is a rule that names a tool and what it may do; [Configure permissions](/docs/en/permissions) explains the syntax. The `$schema` line points to the [published JSON schema](https://json.schemastore.org/claude-code-settings.json) for Claude Code settings, which gives you autocomplete and inline validation in VS Code, Cursor, and any other editor that supports JSON schema. The schema can lag behind the newest CLI releases, so a validation warning on a recently documented key doesn't mean your configuration is invalid.
>
> After you save, run `/status` inside Claude Code to confirm the file loaded; [Confirm what loaded](#check-what-loaded) says what the `Setting sources` line shows and how a broken file is reported.
>
> For a complete personal file, team file, and organization file, each shown with a comment on every key it sets, see the [example settings files](/docs/en/settings-example).
>
>
> ### Change a setting for one session
>
> To try a value without saving it, set it when you start Claude Code. The value applies to that session and your settings files stay as they were. You have three ways to do it:
>
> * **`--settings`**: pass a key as JSON, inline or as a path to a file. Claude Code applies it above your user, project, and local files and below managed settings. It can set any key your user settings file can set; it can't set `Managed` or `Global config` keys.
> * **A flag for that key**: some keys have their own flag, such as `--model` for `model` and `--effort` for `effortLevel`.
> * **An environment variable**: export the key's paired variable before you run `claude`, such as `ANTHROPIC_MODEL` for `model`.
>
> Each key's entry on the [settings reference](/docs/en/settings-reference) lists its per-session overrides and which one takes precedence, so check the entry for the key you want to change.
>
> Commands you run inside a session mostly save your choice: `/config` writes to your settings files, and `/model`
