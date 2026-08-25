---
primary_sources:
  - id: T1-CFG-BASIC
    title: "Config basics"
    url: "https://learn.chatgpt.com/docs/config-file/config-basic.md"
    section: "Configuration precedence; project trust"
  - id: T1-CFG-ADV
    title: "Advanced Configuration"
    url: "https://learn.chatgpt.com/docs/config-file/config-advanced.md"
    section: "Profiles; Config and state locations; Project config files; Hooks; Notifications; History persistence"
  - id: T1-HOOKS
    title: "Codex Hooks"
    url: "https://learn.chatgpt.com/docs/hooks.md"
    section: "Where Codex looks for hooks"
  - id: T1-SKILLS
    title: "Build skills"
    url: "https://learn.chatgpt.com/docs/build-skills.md"
    section: "Where Codex loads local skills"
  - id: T1-RULES
    title: "Rules"
    url: "https://learn.chatgpt.com/docs/agent-configuration/rules.md"
    section: "Create a rules file"
  - id: T1-AGENTS-MD
    title: "Custom instructions with AGENTS.md"
    url: "https://learn.chatgpt.com/docs/agent-configuration/agents-md.md"
    section: "How Codex discovers guidance"
  - id: T1-SUBAGENTS
    title: "Subagents"
    url: "https://learn.chatgpt.com/docs/agent-configuration/subagents.md"
    section: "Custom agents"
  - id: T1-CFG-REF
    title: "Configuration Reference"
    url: "https://learn.chatgpt.com/docs/config-file/config-reference.md"
    section: "config.toml opening; projects.<path>.trust_level"
  - id: T1-CFG-ENV
    title: "Environment variables"
    url: "https://learn.chatgpt.com/docs/config-file/environment-variables.md"
    section: "CODEX_HOME"
  - id: T1-MCP
    title: "Model Context Protocol"
    url: "https://learn.chatgpt.com/docs/extend/mcp.md"
    section: "Connect Codex to an MCP server"
  - id: T1-DEV-SETTINGS
    title: "Developer settings"
    url: "https://learn.chatgpt.com/docs/developer-settings.md"
    section: "IDE extension chatgpt.* keys"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Project vs global path matrix

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

## Summary table

| Concern | Paths (verbatim from docs) | Notes from docs |
| ------- | -------------------------- | --------------- |
| `config.toml` project vs `~/.codex` vs `/etc/codex` | `.codex/config.toml`; `~/.codex/config.toml`; `/etc/codex/config.toml` | Precedence: project (trusted) → profile → user → system → built-ins |
| profiles `$CODEX_HOME/<name>.config.toml` | `$CODEX_HOME/profile-name.config.toml`; `~/.codex/profile-name.config.toml` | Selected with `--profile profile-name` |
| hooks.json / inline hooks (accumulate) | `hooks.json`; inline `[hooks]`; `~/.codex/hooks.json`; `~/.codex/config.toml`; `<repo>/.codex/hooks.json`; `<repo>/.codex/config.toml` | Matching hooks from multiple files all run; layers accumulate |
| rules | `~/.codex/rules/default.rules`; `~/.codex/rules/`; `<repo>/.codex/rules/` | `.rules` under `rules/` next to an active config layer |
| projects trust_level | `projects.<path>.trust_level` (`"trusted"` / `"untrusted"`) | Untrusted skips project-scoped `.codex/` layers |
| AGENTS.md / AGENTS.override.md | `AGENTS.override.md`; `AGENTS.md`; under `~/.codex` / `CODEX_HOME` and project dirs | Global then project walk |
| skills | `.agents/skills`; `$CWD/.agents/skills`; `$HOME/.agents/skills`; `/etc/codex/skills` | Also `$CWD/../.agents/skills`, `$REPO_ROOT/.agents/skills` |
| subagents | `~/.codex/agents/`; `.codex/agents/` | Personal vs project-scoped custom agent TOML |
| MCP | `~/.codex/config.toml`; `.codex/config.toml` | MCP config in `config.toml`; project only when trusted |
| notify (user only) | `notify` in `~/.codex/config.toml` | Ignored in project-local `.codex/config.toml` |
| sessions/history under CODEX_HOME | `CODEX_HOME` (defaults `~/.codex`); `history.jsonl`; sessions under `CODEX_HOME` | Env var lists config, auth, logs, sessions, skills |
| auth.json | `auth.json` under `CODEX_HOME` (defaults `~/.codex`); `~/.codex/auth.json` | File-based credential storage |
| IDE `chatgpt.*` keys | `chatgpt.*` editor settings (e.g. `chatgpt.commentCodeLensEnabled`) | IDE-only; do not go in `config.toml` |

### Source: Config basics — Configuration precedence

> Codex reads configuration details from more than one location. Your personal defaults live in `~/.codex/config.toml`, and you can add project overrides with `.codex/config.toml` files. For security, Codex loads project `.codex/` layers only when you trust the project.
>
> ## Configuration precedence
>
> Codex resolves values in this order (highest precedence first):
>
> 1. CLI flags and `--config` overrides
> 2. Project config files: `.codex/config.toml`, ordered from the project root down to your current working directory (closest wins; trusted projects only)
> 3. [Profile](https://learn.chatgpt.com/docs/config-file/config-advanced#profiles) files selected with `--profile profile-name` (`~/.codex/profile-name.config.toml`)
> 4. User config: `~/.codex/config.toml`
> 5. System config (if present): `/etc/codex/config.toml` on Unix
> 6. Built-in defaults
>
> If you mark a project as untrusted, Codex skips project-scoped `.codex/` layers, including project-local config, hooks, and rules. User and system config still load, including user/global hooks and rules.

### Source: Configuration Reference — `config.toml` opening / profiles / trust_level

> User-level configuration lives in `~/.codex/config.toml`. You can also add project-scoped overrides in `.codex/config.toml` files. Codex loads project-scoped config files only when you trust the project.
>
> Project-scoped config can't override machine-local provider, auth,
> host-owned app request metadata, notification, configuration profile selection,
> or telemetry routing keys. Codex ignores `openai_base_url`,
> `chatgpt_base_url`, `apps_mcp_product_sku`, `model_provider`,
> `model_providers`, `notify`, `profile`, `profiles`,
> `experimental_realtime_ws_base_url`, and `otel` when they appear in a
> project-local `.codex/config.toml`; put provider, notification, and telemetry
> keys in user-level config instead. Config [profile files](https://learn.chatgpt.com/docs/config-file/config-advanced#profiles) live next to
> `config.toml` as `$CODEX_HOME/profile-name.config.toml`; select one with
> `--profile profile-name`.
>
> <ConfigTable
>   options={[
>     {
>       key: "projects.<path>.trust_level",
>       type: "string",
>       description:
>         'Mark a project or worktree as trusted or untrusted (`"trusted"` | `"untrusted"`). Untrusted projects skip project-scoped `.codex/` layers, including project-local config, hooks, and rules.',
>     }
>   ]}
>   client:load
> />

### Source: Advanced Configuration — Profiles

> ## Profiles
>
> Profiles let you save named configuration layers and switch between them from
> the CLI. When you pass `--profile profile-name`, Codex loads
> `~/.codex/config.toml`, then overlays `~/.codex/profile-name.config.toml`.
> Profile names can contain letters, numbers, hyphens, and underscores.
>
> Create a separate TOML file for each profile. Use top-level config keys in the
> profile file; don't nest them under `[profiles.profile-name]`.
>
> ```toml
> # ~/.codex/deep-review.config.toml
> model = "gpt-5.5"
> model_reasoning_effort = "xhigh"
> approval_policy = "on-request"
> model_catalog_json = "/Users/me/.codex/model-catalogs/deep-review.json"
> ```

### Source: Advanced Configuration — Project config files (notify user-only)

> Project config files can't override settings that redirect credentials, alter
> host-owned app request metadata, change provider auth, select config profiles,
> or run machine-local notification/telemetry commands. Codex ignores the
> following keys in project-local `.codex/config.toml` and prints a startup
> warning when it sees them: `openai_base_url`, `chatgpt_base_url`,
> `apps_mcp_product_sku`, `model_provider`, `model_providers`, `notify`,
> `profile`, `profiles`, `experimental_realtime_ws_base_url`, and `otel`. Set
> provider, notification, and telemetry keys in your user-level
> `~/.codex/config.toml`; select config profiles with `--profile profile-name`
> and `~/.codex/profile-name.config.toml`.

### Source: Codex Hooks — Where Codex looks for hooks

> ## Where Codex looks for hooks
>
> Codex discovers hooks next to active config layers in either of these forms:
>
> - `hooks.json`
> - inline `[hooks]` tables inside `config.toml`
>
> In practice, the four most useful locations are:
>
> - `~/.codex/hooks.json`
> - `~/.codex/config.toml`
> - `<repo>/.codex/hooks.json`
> - `<repo>/.codex/config.toml`
>
> If more than one hook source exists, Codex loads all matching hooks.
> Higher-precedence config layers don't replace lower-precedence hooks.
> If a single layer contains both `hooks.json` and inline `[hooks]`, Codex
> merges them and warns at startup. Prefer one representation per layer.
>
> Project-local hooks load only when the project `.codex/` layer is trusted. In
> untrusted projects, Codex still loads user and system hooks from their own
> active config layers.

### Source: Rules — Create a rules file

> 1. Create a `.rules` file under a `rules/` folder next to an active config layer (for example, `~/.codex/rules/default.rules`).
>
> Codex scans `rules/` under every active config layer at startup, including [Team Config](https://learn.chatgpt.com/docs/enterprise/admin-setup#step-4-standardize-local-configuration-with-team-config) locations and the user layer at `~/.codex/rules/`. Project-local rules under `<repo>/.codex/rules/` load only when the project `.codex/` layer is trusted.
>
> When you add a command to the allow list in the TUI, Codex writes to the user layer at `~/.codex/rules/default.rules` so future runs can skip the prompt.

### Source: Custom instructions with AGENTS.md — How Codex discovers guidance

> 1. **Global scope:** In your Codex home directory (defaults to `~/.codex`, unless you set `CODEX_HOME`), Codex reads `AGENTS.override.md` if it exists. Otherwise, Codex reads `AGENTS.md`. Codex uses only the first non-empty file at this level.
> 2. **Project scope:** Starting at the project root (typically the Git root), Codex walks down to your current working directory. If Codex cannot find a project root, it only checks the current directory. In each directory along the path, it checks for `AGENTS.override.md`, then `AGENTS.md`, then any fallback names in `project_doc_fallback_filenames`. Codex includes at most one file per directory.
>
> Use `~/.codex/AGENTS.override.md` when you need a temporary global override without deleting the base file. Remove the override to restore the shared guidance.

### Source: Build skills — Where Codex loads local skills

> Codex reads skills from repository, user, admin, and system locations. For repositories, Codex scans `.agents/skills` in every directory from your current working directory up to the repository root. If two skills share the same `name`, Codex doesn't merge them; both can appear in skill selectors.
>
> | Skill Scope | Location                                                                                                  | Suggested use                                                                                                                                                                                        |
> | :---------- | :-------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `REPO`      | `$CWD/.agents/skills` <br /> Current working directory: where you launch Codex.                           | If you're in a repository or code environment, teams can check in skills relevant to a working folder. For example, skills only relevant to a microservice or a module.                              |
> | `REPO`      | `$CWD/../.agents/skills` <br /> A folder above CWD when you launch Codex inside a Git repository.         | If you're in a repository with nested folders, organizations can check in skills relevant to a shared area in a parent folder.                                                                       |
> | `REPO`      | `$REPO_ROOT/.agents/skills` <br /> The topmost root folder when you launch Codex inside a Git repository. | If you're in a repository with nested folders, organizations can check in skills relevant to everyone using the repository. These serve as root skills available to any subfolder in the repository. |
> | `USER`      | `$HOME/.agents/skills` <br /> Any skills checked into the user's personal folder.                         | Use to curate skills relevant to a user that apply to any repository the user may work in.                                                                                                           |
> | `ADMIN`     | `/etc/codex/skills` <br /> Any skills checked into the machine or container in a shared, system location. | Use for SDK scripts, automation, and for checking in default admin skills available to each user on the machine.                                                                                     |
> | `SYSTEM`    | Bundled with Codex by OpenAI.                                                                             | Useful skills relevant to a broad audience such as the skill-creator and plan skills. Available to everyone when they start Codex.                                                                   |

### Source: Subagents — Custom agents

> To define your own custom agents, add standalone TOML files under
> `~/.codex/agents/` for personal agents or `.codex/agents/` for project-scoped
> agents.

### Source: Model Context Protocol — Connect Codex to an MCP server

> Codex stores MCP configuration in `config.toml` alongside other Codex configuration settings. By default this is `~/.codex/config.toml`, but you can also scope MCP servers to a project with `.codex/config.toml` (trusted projects only).

### Source: Advanced Configuration — Notifications (notify)

> ## Notifications
>
> Use `notify` to trigger an external program whenever Codex emits supported events (currently only `agent-turn-complete`). This is handy for desktop toasts, chat webhooks, CI updates, or any side-channel alerting that the built-in TUI notifications don't cover.
>
> ```toml
> notify = ["python3", "/path/to/notify.py"]
> ```

### Source: Advanced Configuration — Config and state locations / History persistence

> ## Config and state locations
>
> Codex stores its local state under `CODEX_HOME` (defaults to `~/.codex`).
>
> Common files you may see there:
>
> - `config.toml` (your local configuration)
> - `auth.json` (if you use file-based credential storage) or your OS keychain/keyring
> - `history.jsonl` (if history persistence is enabled)
> - Other per-user state such as logs and caches
>
> ## History persistence
>
> By default, Codex saves local session transcripts under `CODEX_HOME` (for example, `~/.codex/history.jsonl`). To disable local history persistence:
>
> ```toml
> [history]
> persistence = "none"
> ```

### Source: Environment variables — `CODEX_HOME`

> | `CODEX_HOME`        | CLI, IDE extension, app-server, installers | `~/.codex`   | Sets the root for Codex state, including config, auth, logs, sessions, skills, and standalone package metadata. If you set it, the directory must already exist. |

### Source: Developer settings — IDE `chatgpt.*` keys

> The Codex IDE extension has two settings layers:
>
> - **Codex settings** control agent behavior shared with Codex CLI, including the
>   model, reasoning effort, permissions, sandbox, MCP servers, and
>   personalization. Codex reads these settings from `config.toml`.
> - **Editor settings** control how the extension behaves inside VS Code and
>   compatible editors. These settings use `chatgpt.*` keys in the editor's
>   settings system.
>
> | Setting                                      | Default        | Description                                                                                                                                                                                                                                                                                |
> | -------------------------------------------- | -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `chatgpt.commentCodeLensEnabled`             | `true`         | Show CodeLens above `TODO` comments so Codex can address them.                                                                                                                                                                                                                             |
> | `chatgpt.openOnStartup`                      | `false`        | Focus the Codex sidebar when the extension finishes starting.                                                                                                                                                                                                                              |
> | `chatgpt.followUpQueueMode`                  | `queue`        | Choose whether messages sent during a run wait for the next run (`queue`) or steer the current run (`steer`). The extension treats the legacy `interrupt` value as `steer`. Press <kbd>Cmd</kbd>/<kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>Enter</kbd> to invert the behavior for one message. |
> | `chatgpt.composerEnterBehavior`              | `enter`        | Choose whether <kbd>Enter</kbd> always sends (`enter`), <kbd>Cmd</kbd>/<kbd>Ctrl</kbd>+<kbd>Enter</kbd> sends multiline prompts (`cmdIfMultiline`), or the modifier is always required (`cmdAlways`).                                                                                      |
> | `chatgpt.reviewDelivery`                     | `inline`       | Run `/review` in the current chat when possible (`inline`) or start a separate review chat (`detached`).                                                                                                                                                                                   |
> | `chatgpt.localeOverride`                     | Auto           | Set the preferred language for the Codex UI. Leave empty to detect it automatically.                                                                                                                                                                                                       |
> | `chatgpt.runCodexInWindowsSubsystemForLinux` | `false`        | Windows only: Run Codex in WSL when WSL is available. Use this when your repositories and tooling live in WSL2 or when you need Linux-native tooling. Changing this setting reloads VS Code.                                                                                               |
> | `chatgpt.cliExecutable`                      | Unset          | Development only: Set the path to the Codex CLI executable. You don't need this setting unless you're developing the Codex CLI; manually overriding the bundled executable can prevent parts of the extension from working.                                                                |
>
> The `chatgpt.*` keys above belong to the IDE extension and don't go in
> `config.toml`.
