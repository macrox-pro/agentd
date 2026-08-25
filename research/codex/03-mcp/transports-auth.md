---
primary_sources:
  - id: T1-MCP
    title: "Model Context Protocol"
    url: "https://learn.chatgpt.com/docs/extend/mcp.md"
    section: "Supported MCP features; STDIO; Streamable HTTP; OAuth client registration; timeouts; tools; approval_mode"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# MCP transports and auth

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Model Context Protocol — Features / Transports / Auth

> The ChatGPT desktop app, Codex CLI, and IDE extension support MCP servers and
> share MCP configuration for the same Codex host.
>
> The supported server features below apply to MCP servers configured on a Codex
> host. Hosted plugin tools can have different capabilities.
>
> ## Supported MCP features
>
> - **STDIO servers**: Servers that run as a local process (started by a command).
>   - Environment variables
> - **Streamable HTTP servers**: Servers that you access at an address.
>   - Bearer token authentication
>   - OAuth authentication, including Client ID Metadata Documents (CIMD) and
>     Dynamic Client Registration (DCR)
>   - ChatGPT session authentication for trusted first-party servers
> - **Server instructions**: Codex reads the MCP `instructions` field returned during initialization and uses it as server-wide guidance alongside the server's tools.
>
> If you build or maintain an MCP server for Codex, use `instructions` for cross-tool workflows, constraints, and rate limits that apply across the server. Keep the first 512 characters self-contained so the most important guidance is available when Codex is deciding how to use the server.
>
> #### STDIO servers
>
> - `command` (required): The command that starts the server.
> - `args` (optional): Arguments to pass to the server.
> - `env` (optional): Environment variables to set for the server.
> - `env_vars` (optional): Environment variables to allow and forward.
> - `cwd` (optional): Working directory to start the server from.
> - `experimental_environment` (optional): Set to `remote` to start the stdio
>   server through a remote executor environment when one is available.
>
> `env_vars` can contain plain variable names or objects with a source:
>
> ```toml
> env_vars = ["LOCAL_TOKEN", { name = "REMOTE_TOKEN", source = "remote" }]
> ```
>
> String entries and `source = "local"` read from Codex's local environment.
> `source = "remote"` reads from the remote executor environment and requires
> remote MCP stdio.
>
> #### Streamable HTTP servers
>
> - `url` (required): The server address.
> - `auth` (optional): Authentication to try after configured bearer tokens and
>   authorization headers. Use `oauth` (the default) for stored MCP OAuth
>   credentials. Use `chatgpt` to use the current ChatGPT session for the trusted
>   first-party ChatGPT origin, with stored OAuth as a fallback.
> - `bearer_token_env_var` (optional): Environment variable name for a bearer token to send in `Authorization`.
> - `http_headers` (optional): Map of header names to static values.
> - `env_http_headers` (optional): Map of header names to environment variable names (values pulled from the environment).
>
> If no credential source resolves, Codex can connect to the server without
> authentication. Run `codex mcp login <server-name>` separately to start an MCP
> OAuth login.
>
> #### Other configuration options
>
> - `startup_timeout_sec` (optional): Timeout (seconds) for the server to start. Default: `10`.
> - `tool_timeout_sec` (optional): Timeout (seconds) for the server to run a tool. Default: `60`.
> - `enabled` (optional): Set `false` to disable a server without deleting it.
> - `required` (optional): Set `true` to make startup fail if this enabled server can't initialize.
> - `enabled_tools` (optional): Tool allow list.
> - `disabled_tools` (optional): Tool deny list (applied after `enabled_tools`).
> - `default_tools_approval_mode` (optional): Default approval behavior for
>   tools from this server. Supported values are `auto`, `prompt`, `writes`, and
>   `approve`. The `writes` mode prompts for tools that aren't marked read-only.
> - `tools.<tool>.approval_mode` (optional): Per-tool approval behavior override.
>
> If your OAuth provider requires a fixed callback port, set the top-level `mcp_oauth_callback_port` in `config.toml`. If unset, Codex binds to an ephemeral port.
>
> If your MCP OAuth flow must use a specific callback URL (for example, a remote Devbox ingress URL or a custom callback path), set `mcp_oauth_callback_url`. Codex uses this value as the base callback URL, then appends a server-specific callback ID to produce the OAuth `redirect_uri` it sends during login. Register the full derived `redirect_uri` with your OAuth provider, including the appended callback ID and any configured path, query, or port, rather than registering only the base host or path without that suffix. Local callback URLs (for example `localhost`) bind on the local interface; non-local callback URLs bind on `0.0.0.0` so the callback can reach the host.
>
> If the MCP server advertises `scopes_supported`, Codex prefers those
> server-advertised scopes during OAuth login. Otherwise, Codex falls back to the
> scopes configured in `config.toml`.
>
> #### OAuth client registration
>
> Codex supports [OAuth Client ID Metadata Documents (CIMD)](https://datatracker.ietf.org/doc/draft-ietf-oauth-client-id-metadata-document/)
> and Dynamic Client Registration (DCR). By default, Codex automatically chooses
> CIMD when the authorization server advertises
> `client_id_metadata_document_supported: true`, includes `none` in
> `token_endpoint_auth_methods_supported`, and the callback uses a supported
> loopback URL. Otherwise, Codex uses DCR when available. A configured OAuth client
> ID always takes precedence and skips client registration.
>
> For CIMD, Codex uses a ChatGPT-hosted metadata document specific to the MCP
> server:
>
> ```text
> https://chatgpt.com/oauth/codex/<callback_id>/client.json
> ```
>
> Codex derives `<callback_id>` from the MCP server URL and includes it in the
> loopback redirect URI, such as
> `http://127.0.0.1:<port>/callback/<callback_id>`. The metadata document registers
> the matching loopback URI without a port. Authorization servers must accept the
> port selected at login while matching the host and path exactly, as required by
> [RFC 8252](https://www.rfc-editor.org/rfc/rfc8252.html#section-7.3). Custom
> callback hosts, paths, or query parameters require DCR or a configured OAuth
> client ID.
>
> Support for a stable, shared CIMD document is in development and coming soon:
>
> ```text
> https://chatgpt.com/oauth/codex/client.json
> ```
>
> Codex will use the stable document with the shared `/callback` path when the
> authorization server advertises
> `authorization_response_iss_parameter_supported: true`, provides a valid
> `issuer` in its metadata, and includes a matching `iss` in authorization
> responses. Servers without issuer-bound responses will continue using the
> callback-specific document.
>
> To choose a registration method for one CLI login, use
> `--oauth-client-registration`:
>
> ```bash
> codex mcp login <server-name> --oauth-client-registration cimd
> codex mcp login <server-name> --oauth-client-registration dcr
> ```
>
> The default is `auto`. Registration choices apply only to the current login and
> aren't stored in `config.toml`.
