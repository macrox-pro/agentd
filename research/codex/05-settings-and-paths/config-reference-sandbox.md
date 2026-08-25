---
primary_sources:
  - id: T1-CFG-REF
    title: "Configuration Reference"
    url: "https://learn.chatgpt.com/docs/config-file/config-reference.md"
    section: "approval_policy; sandbox_mode; sandbox_workspace_write.*; permissions.*; default_permissions; allowed_permission_profiles"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Config reference — sandbox and permissions

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Configuration Reference — sandbox and approval keys (cross-links)

> For sandbox and approval keys (`approval_policy`, `sandbox_mode`, and `sandbox_workspace_write.*`), pair this reference with [Sandbox and approvals](https://learn.chatgpt.com/docs/agent-approvals-security#sandbox-and-approvals), [Protected paths in writable roots](https://learn.chatgpt.com/docs/agent-approvals-security#protected-paths-in-writable-roots), and [Network access](https://learn.chatgpt.com/docs/agent-approvals-security#network-access). For beta permission profiles, see [Permissions](https://learn.chatgpt.com/docs/permissions).

### Source: Configuration Reference — `config.toml` approval, sandbox, permissions

> <ConfigTable
>   options={[
>     {
>       key: "approval_policy",
>       type: "untrusted | on-request | never | { granular = { sandbox_approval = bool, rules = bool, mcp_elicitations = bool, request_permissions = bool, skill_approval = bool } }",
>       description:
>         "Controls when Codex pauses for approval before executing commands. You can also use `approval_policy = { granular = { ... } }` to allow or auto-reject specific prompt categories while keeping other prompts interactive. `on-failure` is deprecated; use `on-request` for interactive runs or `never` for non-interactive runs.",
>     },
>     {
>       key: "approval_policy.granular.sandbox_approval",
>       type: "boolean",
>       description:
>         "When `true`, sandbox escalation approval prompts are allowed to surface.",
>     },
>     {
>       key: "approval_policy.granular.rules",
>       type: "boolean",
>       description:
>         "When `true`, approvals triggered by execpolicy `prompt` rules are allowed to surface.",
>     },
>     {
>       key: "approval_policy.granular.mcp_elicitations",
>       type: "boolean",
>       description:
>         "When `true`, MCP elicitation prompts are allowed to surface instead of being auto-rejected.",
>     },
>     {
>       key: "approval_policy.granular.request_permissions",
>       type: "boolean",
>       description:
>         "When `true`, prompts from the `request_permissions` tool are allowed to surface.",
>     },
>     {
>       key: "approval_policy.granular.skill_approval",
>       type: "boolean",
>       description:
>         "When `true`, skill-script approval prompts are allowed to surface.",
>     },
>     {
>       key: "sandbox_mode",
>       type: "read-only | workspace-write | danger-full-access",
>       description:
>         "Sandbox policy for filesystem and network access during command execution.",
>     },
>     {
>       key: "sandbox_workspace_write.writable_roots",
>       type: "array<string>",
>       description:
>         'Additional writable roots when `sandbox_mode = "workspace-write"`.',
>     },
>     {
>       key: "sandbox_workspace_write.network_access",
>       type: "boolean",
>       description:
>         "Allow outbound network access inside the workspace-write sandbox.",
>     },
>     {
>       key: "sandbox_workspace_write.exclude_tmpdir_env_var",
>       type: "boolean",
>       description:
>         "Exclude `$TMPDIR` from writable roots in workspace-write mode.",
>     },
>     {
>       key: "sandbox_workspace_write.exclude_slash_tmp",
>       type: "boolean",
>       description:
>         "Exclude `/tmp` from writable roots in workspace-write mode.",
>     },
>     {
>       key: "default_permissions",
>       type: "string",
>       description:
>         "Name of the default permissions profile to apply to sandboxed tool calls. Built-ins are `:read-only`, `:workspace`, and `:danger-full-access`; custom profile names require matching `[permissions.<name>]` tables. Don't combine with `sandbox_mode` or `[sandbox_workspace_write]`.",
>     },
>     {
>       key: "permissions.<name>.description",
>       type: "string",
>       description:
>         "Human-readable description for this named profile. A profile does not inherit its parent's description through `extends`.",
>     },
>     {
>       key: "permissions.<name>.extends",
>       type: "string",
>       description:
>         "Optional parent profile applied before this named profile. Set it to another named profile, `:read-only`, or `:workspace`; `:danger-full-access`, undefined parents, and cycles are rejected.",
>     },
>     {
>       key: "permissions.<name>.workspace_roots",
>       type: "table",
>       description:
>         "Profile-defined workspace roots that receive `:workspace_roots` filesystem rules alongside the session's runtime workspace roots.",
>     },
>     {
>       key: "permissions.<name>.workspace_roots.<path>",
>       type: "boolean",
>       description:
>         "Opt a path into the profile's workspace root set when `true`. Disabled entries remain inactive.",
>     },
>     {
>       key: "permissions.<name>.filesystem",
>       type: "table",
>       description:
>         "Named filesystem permission profile. Each key is an absolute path or special token such as `:minimal` or `:workspace_roots`.",
>     },
>     {
>       key: "permissions.<name>.filesystem.glob_scan_max_depth",
>       type: "number",
>       description:
>         "Maximum depth for expanding deny-read glob patterns on platforms that snapshot matches before sandbox startup. Must be at least `1` when set.",
>     },
>     {
>       key: "permissions.<name>.filesystem.<path-or-glob>",
>       type: '"read" | "write" | "deny" | table',
>       description:
>         'Grant direct access for a path, glob pattern, or special token, or scope nested entries under that root. Use `"deny"` to deny reads for matching paths.',
>     },
>     {
>       key: 'permissions.<name>.filesystem.":workspace_roots".<subpath-or-glob>',
>       type: '"read" | "write" | "deny"',
>       description:
>         'Scoped filesystem access relative to each effective workspace root. Use `"."` for the root itself; glob subpaths such as `"**/*.env"` can deny reads with `"deny"`.',
>     },
>     {
>       key: "permissions.<name>.network.enabled",
>       type: "boolean",
>       description:
>         "Enable network access for commands in this permission profile. This does not start the network proxy. Without `features.network_proxy` or enabled administrator-managed networking requirements, command network access is direct and profile domain rules are not enforced.",
>     },
>     {
>       key: "permissions.<name>.network.proxy_url",
>       type: "string",
>       description:
>         "HTTP listener URL used when this permissions profile enables sandboxed networking.",
>     },
>     {
>       key: "permissions.<name>.network.enable_socks5",
>       type: "boolean",
>       description:
>         "Expose SOCKS5 support when this permissions profile enables sandboxed networking.",
>     },
>     {
>       key: "permissions.<name>.network.socks_url",
>       type: "string",
>       description: "SOCKS5 proxy endpoint used by this permissions profile.",
>     },
>     {
>       key: "permissions.<name>.network.enable_socks5_udp",
>       type: "boolean",
>       description: "Allow UDP over the SOCKS5 listener when enabled.",
>     },
>     {
>       key: "permissions.<name>.network.allow_upstream_proxy",
>       type: "boolean",
>       description:
>         "Allow sandboxed networking to chain through another upstream proxy.",
>     },
>     {
>       key: "permissions.<name>.network.dangerously_allow_non_loopback_proxy",
>       type: "boolean",
>       description:
>         "Permit non-loopback bind addresses for sandboxed networking listeners. Enabling it can expose listeners beyond localhost.",
>     },
>     {
>       key: "permissions.<name>.network.dangerously_allow_all_unix_sockets",
>       type: "boolean",
>       description:
>         "Allow arbitrary Unix socket destinations instead of the default restricted set. Use only in tightly controlled environments.",
>     },
>     {
>       key: "permissions.<name>.network.mode",
>       type: "limited | full",
>       description: "Network proxy mode used for subprocess traffic.",
>     },
>     {
>       key: "permissions.<name>.network.domains",
>       type: "table",
>       description:
>         "Domain rules for sandboxed commands. Enforced only when `features.network_proxy` or enabled administrator-managed networking requirements activate the proxy. Supports exact hosts, `*.example.com`, `**.example.com`, and global `*` allow rules; `deny` wins. Does not restrict web search, apps, or MCP servers.",
>     },
>     {
>       key: "permissions.<name>.network.domains.<pattern>",
>       type: "allow | deny",
>       description:
>         "Allow or deny an exact host or scoped wildcard pattern such as `*.example.com` or `**.example.com`.",
>     },
>     {
>       key: "permissions.<name>.network.unix_sockets",
>       type: "table",
>       description:
>         "Unix socket allowlist overrides for sandboxed networking. Use socket paths as keys; `allow` adds a path, and `deny` rejects it.",
>     },
>     {
>       key: "permissions.<name>.network.unix_sockets.<path>",
>       type: "allow | deny",
>       description:
>         "Add an absolute Unix socket path to the effective allowlist with `allow`, or reject it with `deny`. Denied entries are omitted from the effective allowlist.",
>     },
>     {
>       key: "permissions.<name>.network.allow_local_binding",
>       type: "boolean",
>       description:
>         "Permit broader local/private-network access through sandboxed networking. Exact local IP literal or `localhost` allow rules can still permit specific local targets when this stays `false`.",
>     }
>   ]}
>   client:load
> />

### Source: Configuration Reference — `requirements.toml` allowed_permission_profiles / default_permissions

> Use `allowed_sandbox_modes` with `sandbox_mode`. For permission-profile
> deployments, use `allowed_permission_profiles` with managed
> `default_permissions`.
>
> The `[models.new_thread]` table supplies managed defaults, not enforcement.
> Explicit launch choices from dedicated CLI flags or `--config` overrides take
> precedence. An explicit model or reasoning-effort override skips both managed
> model fields; `service_tier` is independent.

> <ConfigTable
>   options={[
>     {
>       key: "allowed_permission_profiles",
>       type: "table<boolean>",
>       description:
>         "Complete list of allowed permission profiles. Profiles set to `true` are allowed. Profiles that are omitted or set to `false` are denied, including profiles added in future versions. When requirements sources are combined, entries are matched by profile name.",
>     },
>     {
>       key: "allowed_permission_profiles.<name>",
>       type: "boolean",
>       description:
>         "Allow or deny a built-in or custom permission profile defined in a loaded config or requirements source. A later, higher-precedence requirements source can use `false` to turn off a profile allowed by an earlier, lower-precedence source.",
>     },
>     {
>       key: "default_permissions",
>       type: "string",
>       description:
>         "Managed default permission profile. The profile must be allowed by `allowed_permission_profiles`. Set this explicitly for predictable behavior; if omitted, Codex defaults to `:workspace` only when both `:workspace` and `:read-only` are explicitly allowed.",
>     },
>     {
>       key: "permissions.<name>",
>       type: "table",
>       description:
>         "Admin-defined permission profile. The name can't start with `:`, use the reserved name `filesystem`, or duplicate a profile from a loaded config. Uses the same profile fields as `config.toml`; see the Permissions guide for the complete profile schema.",
>     },
>     {
>       key: "permissions.filesystem.deny_read",
>       type: "array<string>",
>       description:
>         "Admin-enforced filesystem read denials. Entries can be paths or glob patterns, and users cannot weaken them with local config.",
>     }
>   ]}
>   client:load
> />
