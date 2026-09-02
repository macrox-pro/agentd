---
primary_sources:
  - id: T1-NETWORK
    title: "Network"
    url: "https://opencode.ai/docs/network.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Network

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: Network — Full page

> OpenCode supports standard proxy environment variables and custom certificates for enterprise network environments.
>
> ---
>
> ## Proxy
>
> OpenCode respects standard proxy environment variables.
>
> ```bash
> # HTTPS proxy (recommended)
> export HTTPS_PROXY=https://proxy.example.com:8080
>
> # HTTP proxy (if HTTPS not available)
> export HTTP_PROXY=http://proxy.example.com:8080
>
> # Bypass proxy for local server (required)
> export NO_PROXY=localhost,127.0.0.1
> ```
>
> :::caution
> The TUI communicates with a local HTTP server. You must bypass the proxy for this connection to prevent routing loops.
> :::
>
> You can configure the server's port and hostname using [CLI flags](/docs/cli#run).
>
> ---
>
> ### Authenticate
>
> If your proxy requires basic authentication, include credentials in the URL.
>
> ```bash
> export HTTPS_PROXY=http://username:password@proxy.example.com:8080
> ```
>
> :::caution
> Avoid hardcoding passwords. Use environment variables or secure credential storage.
> :::
>
> For proxies requiring advanced authentication like NTLM or Kerberos, consider using an LLM Gateway that supports your authentication method.
>
> ---
>
> ## Custom certificates
>
> If your enterprise uses custom CAs for HTTPS connections, configure OpenCode to trust them.
>
> ```bash
> export NODE_EXTRA_CA_CERTS=/path/to/ca-cert.pem
> ```
>
> This works for both proxy connections and direct API access.
