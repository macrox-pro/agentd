---
primary_sources:
  - id: T1-EP-SEC
    title: "Endpoint"
    url: "https://cursor.com/docs/enterprise/endpoint-security.md"
    section: "Endpoint"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Endpoint security and network

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Endpoint security

> # Endpoint Security Configuration
>
> Cursor loads JavaScript modules and performs file I/O during startup. Endpoint security software that intercepts file operations or injects into processes can slow startup past internal timeouts, causing features like Agent to fail. This page covers how to configure exclusions so Cursor works alongside your security stack.
>
> ## What to exclude
>
> Add the following processes and paths to your security product's exclusion list.
>
> ### Windows
>
> **Processes:** each process has a user install path and a system install path. Add the path that matches your install type.
>
> | Process            | Install type | Path                                                                                             |
> | ------------------ | ------------ | ------------------------------------------------------------------------------------------------ |
> | `Cursor.exe`       | User         | `%LOCALAPPDATA%\Programs\cursor\Cursor.exe`                                                      |
> | `Cursor.exe`       | System       | `%ProgramFiles%\cursor\Cursor.exe`                                                               |
> | `rg.exe`           | User         | `%LOCALAPPDATA%\Programs\cursor\resources\app\node_modules\@vscode\ripgrep\bin\rg.exe`           |
> | `rg.exe`           | System       | `%ProgramFiles%\cursor\resources\app\node_modules\@vscode\ripgrep\bin\rg.exe`                    |
> | `inno_updater.exe` | User         | `%LOCALAPPDATA%\Programs\cursor\resources\app\node_modules\cursor-inno-updater\inno_updater.exe` |
> | `inno_updater.exe` | System       | `%ProgramFiles%\cursor\resources\app\node_modules\cursor-inno-updater\inno_updater.exe`          |
>
> **Paths:**
>
> | Path                              | Description                                               |
> | --------------------------------- | --------------------------------------------------------- |
> | `%LOCALAPPDATA%\Programs\cursor\` | Application binaries and bundled modules (user install)   |
> | `%ProgramFiles%\cursor\`          | Application binaries and bundled modules (system install) |
> | `%APPDATA%\Cursor\`               | User data, settings, and workspace storage                |
>
> ### macOS
>
> **Processes:** `Cursor.app`
>
> **Paths:**
>
> | Path                        | Description        |
> | --------------------------- | ------------------ |
> | `/Applications/Cursor.app/` | Application bundle |
>
> ## Why exclusions may be needed
>
> Cursor's extension host reads JavaScript files from its own install directory at startup. When security software adds per-file scanning latency, the cumulative delay can exceed Cursor's startup timeout.
>
> This primarily affects startup. Once modules are loaded into memory, ongoing file operations are infrequent and unlikely to cause issues.
>
> Cursor's own files are code-signed binaries and bundled JavaScript, not user-generated content. Excluding them from real-time scanning is low-risk and does not reduce protection for user files or network traffic.
>
> Both **process exclusions** and **path exclusions** may be needed. Some products use kernel-level minifilter drivers that scan all file I/O regardless of which process is reading. A process-only exclusion may not be sufficient — add path exclusions for the Cursor install directory as well.
>
> ## Identifying active security software
>
> These commands can help identify which products are running so you know where you may need to configure exclusions. On Windows, run in an **Administrator PowerShell** window:
>
> ```powershell
> # Registered AV products
> Get-CimInstance -Namespace root/SecurityCenter2 -ClassName AntiVirusProduct |
>   Select-Object displayName, pathToSignedProductExe
>
> # Kernel-level filesystem filter drivers
> fltmc
>
> # Check for EDR process injection via environment variables
> [System.Environment]::GetEnvironmentVariables() |
>   Where-Object { $_.Keys -match "BPP|COR_PROFILER|COMPLUS|__COMPAT" }
>
> # Windows Defender status
> Get-MpComputerStatus |
>   Select-Object IsTamperProtected, RealTimeProtectionEnabled, AMRunningMode
> ```
>
> **How to read `fltmc` output:** Standard Windows drivers you can ignore include `WdFilter`, `storqosflt`, `wcifs`, `CldFlt`, `bfs`, `FileCrypt`, `luafv`, `Wof`, `FileInfo`, `npsvctrig`, `bindflt`, and `UnionFS`. Other drivers are likely from third-party security software.
>
> **How to read the environment variable output:** If it returns any results, an EDR product is injecting code into every new process on the machine, and an exclusion may be necessary.
>
> ## Verifying exclusions are working
>
> After applying exclusions, restart Cursor and verify that Agent features work without timing out. If you previously saw empty Extension Host logs (Cmd/Ctrl+Shift+P → "Output" → "Extension Host"), they should now show normal startup output.
>
> ## Troubleshooting checklist
>
> 1. Run the [identification commands above](https://cursor.com/docs/enterprise/endpoint-security.md#identifying-active-security-software) to determine which security products are running
> 2. Add both process and path exclusions for the identified products in their management consoles
> 3. Restart Cursor and test Agent — this is the definitive test of whether exclusions are working
> 4. If exclusions don't resolve the issue, [export logs](https://cursor.com/help/troubleshooting/agent-issues.md#what-if-i-see-agent-execution-timed-out) and contact Cursor support with the diagnostic output\`
>
>

### Source: Network configuration

> # Network Configuration
>
> Cursor needs to communicate with backend services and AI providers. This documentation covers how to configure Cursor to work within your network infrastructure, including proxies, firewalls, and encryption requirements.
>
> ## Proxy configuration
>
> Many enterprises route traffic through proxy servers for monitoring and security. Cursor works with most proxy configurations, but some proxy settings can cause issues with streaming responses.
>
> ### HTTP/2 vs HTTP/1.1
>
> Cursor uses HTTP/2 bidirectional streaming by default for real-time chat and agent experiences. Some enterprise proxies don't support HTTP/2 streaming correctly. Zscaler is the most widely used proxy with this limitation.
>
> If you experience issues with streaming, Cursor automatically falls back to HTTP/1.1 Server-Sent Events (SSE) mode. This fallback was specifically designed to work with Zscaler and similar proxies that buffer or break HTTP/2 streams. The fallback happens transparently when HTTP/2 bidirectional streaming doesn't work.
>
> ### SSL inspection and DLP
>
> Many corporate proxies perform SSL man-in-the-middle inspection to scan traffic for security threats or data loss prevention (DLP). This replaces Cursor's SSL certificates with your proxy's certificates.
>
> When Cursor traffic goes through Secure Web Gateways (SWG), SSL inspection, or DLP, it often causes timeouts, slowness, or errors when using Cursor's Agent capabilities. This is one of the most common deployment blockers for enterprise customers. For endpoint security software (AV, EDR, DLP) that runs on the machine itself rather than at the network level, see [Endpoint Security Configuration](https://cursor.com/docs/enterprise/endpoint-security.md).
>
> Cursor's services are already encrypted end-to-end. We recommend disabling SSL inspection for these domains:
>
> - `.cursor.sh`
> - `cursor-cdn.com`
> - `marketplace.cursorapi.com`
> - `authenticate.cursor.sh`
> - `authenticator.cursor.sh`
> - `*.cursorvm.com`
> - `*.*.cursorvm.com`
>
> If your security policy requires SSL inspection on all traffic, your proxy must support:
>
> - HTTP/2 bidirectional streaming (or that Cursor's HTTP/1.1 fallback works)
> - Server-Sent Events (SSE) passthrough without buffering
> - Long-running connections without forced timeouts
> - Disabling response buffering for streaming content types
>
> ### Testing proxy connectivity
>
> If you experience connection issues, you can test connectivity manually using curl commands. These commands simulate the requests Cursor makes to backend services.
>
> **Test basic connectivity:**
>
> ```bash
> curl -v https://api2.cursor.sh |& grep -C1 issuer:
> ```
>
> This shows which SSL certificate is in use. You should see Amazon RSA. If you see your proxy provider (like Zscaler), SSL inspection is active.
>
> **Test HTTP/1.1 streaming:**
>
> ```bash
> echo -ne "\x0\x0\x0\x0\x11{\"payload\":\"foo\"}" | \
>   curl --http1.1 -No - -XPOST \
>   -H "Content-Type: application/connect+json" \
>   --data-binary @- \
>   https://api2.cursor.sh/aiserver.v1.HealthService/StreamSSE
> ```
>
> You should see output appear line by line over 5 seconds. If it appears all at once after 5 seconds, your proxy is buffering streaming responses.
>
> **Test HTTP/2 bidirectional streaming:**
>
> ```bash
> (for i in 1 2 3 4 5; do \
>   echo -ne "\x0\x0\x0\x0\x12{\"payload\":\"foo$i\"}"; \
>   sleep 1; \
> done) | curl -No - -XPOST \
>   -H "Content-Type: application/connect+json" \
>   -T - \
>   https://api2.cursor.sh/aiserver.v1.HealthService/StreamBidi
> ```
>
> Output should appear once per second. If buffered for 5 seconds, your proxy doesn't support HTTP/2 bidirectional streaming.
>
> ## IP allowlisting
>
> If your network uses IP-based access controls, you need to allow traffic to Cursor's backend services.
>
> Rather than maintaining IP address lists (which can change), configure your firewall to allow traffic to these domain patterns:
>
> - `*.cursor.sh`
> - `*.cursor-cdn.com`
> - `*.cursorapi.com`
> - `*.cursorvm.com`
> - `*.*.cursorvm.com`
>
> We generally recommend allowlisting with the domain patterns above. However, if your firewall mandates granular subdomain allowlists without wildcards, use the following list:
>
> - `api2.cursor.sh`: Used for most API requests.
> - `api5.cursor.sh`: Used for Cursor's agent requests.
> - `api3.cursor.sh`: Used for Cursor Tab requests (HTTP/2 only).
> - `repo42.cursor.sh`: Used for codebase indexing (HTTP/2 only).
> - `api4.cursor.sh`, `us-asia.gcpp.cursor.sh`, `us-eu.gcpp.cursor.sh`, `us-only.gcpp.cursor.sh`: Used for Cursor Tab requests depending on your location (HTTP/2 only).
> - `adminportal42.cursor.sh`: Used to configure SSO and domain verification.
> - `marketplace.cursorapi.com`, `cursor-cdn.com`, `downloads.cursor.com`, `anysphere-binaries.s3.us-east-1.amazonaws.com`: Used for client updates and downloading extensions from the extension marketplace.
> - `api5.cursor.sh`: Used for network access layer (NAL) requests. These subdomains are also used:
>   - `agent.api5.cursor.sh`
>   - `agentn.api5.cursor.sh`
>   - `agent.us.api5.cursor.sh`
>   - `agentn.us.api5.cursor.sh`
>   - `agent.global.api5.cursor.sh`
>   - `agentn.global.api5.cursor.sh`
> - `authenticate.cursor.sh`: Authorization endpoint.
> - `authenticator.cursor.sh`: Auth UI and login webview.
> - `prod.authentication.cursor.sh`: Production token issuer.
> - `authentication.cursor.sh`: JWT issuer (backend).
>
> ## Private connectivity
>
> Cursor supports [private connectivity](https://cursor.com/docs/cloud-agent/private-connectivity.md) for Enterprise teams that need Cloud Agents, Bugbot, or Cursor backend services to access private source control systems and package registries. Supported options include AWS PrivateLink and Cloudflare Tunnel.
>
> Cursor does not currently offer VPC peering or customer-facing Google Private Service Connect.
>
> When you run Cursor agents in the editor or via the CLI, they inherit your existing network configuration. If you run Cursor on a machine within your VPC, agent operations inherit:
>
> - Your network security groups
> - Your firewall rules
> - Your DNS configuration
> - Your VPN or private network access
>
> This means Cursor agents can access internal resources that the machine can reach, while following your existing network security controls.
>
> ## Encryption
>
> Cursor encrypts data both in transit and at rest.
>
> ### In transit
>
> - TLS 1.2 or higher for all connections to Cursor services
> - TLS for connections to third-party AI providers
> - Certificate pinning for critical services
>
> ### At rest
>
> - AES-256 encryption for stored data
> - Encrypted vector database storage
> - Encrypted code storage for Cloud Agents (when enabled)
>
> ### Key management
>
> Cursor manages encryption keys. Keys are rotated according to security best practices and stored in secure key management systems.
>
> For enhanced security control, enterprise customers can use Customer Managed Encryption Keys (CMEK) for encrypting data stored in Cursor's infrastructure. See [Data Encryption](https://cursor.com/docs/enterprise/privacy-and-data-governance.md#data-encryption) for details.
>
> ## LLM gateways
>
> Some enterprises want to route LLM traffic through their own gateways for additional monitoring and control.
>
> Custom gateways can introduce additional latency, rate limiting, and compatibility issues. We instead recommend using Cursor's built-in hooks feature to implement your own security controls.
>
> Cursor's [Zero Data Retention policy](https://cursor.com/docs/account/teams/dashboard.md#settings) does not apply when using your own API keys. Your data handling will be subject to the privacy policies of your chosen AI provider (OpenAI, Anthropic, Google, Azure, or AWS).
>
> See [Hooks](https://cursor.com/docs/hooks.md) and [Security Guardrails](https://cursor.com/docs/enterprise/llm-safety-and-controls.md) for details.
>
> ## Cloud Agents networking
>
> Cloud Agents run on Cursor's infrastructure, not your local network. They can access:
>
> - Public GitHub repositories
> - GitHub Enterprise Cloud repositories you've granted access to
> - GitHub Enterprise Server (self-hosted GitHub Enterprise)
> - On-prem and cloud-based GitLab
> - Bitbucket Cloud repositories
> - Public package registries (npm, PyPI, etc.)
>
> Cloud Agents cannot access:
>
> - Resources behind your corporate firewall
> - On-premises GitHub Enterprise Server
> - Private package registries without internet access
>
> If your development workflow requires access to internal resources, use the Cursor editor on machines within your network instead of Cloud Agents.
>
> ## Troubleshooting checklist
>
> If you experience connection issues:
>
> 1. **Test basic connectivity** to `api2.cursor.sh`
> 2. **Check if SSL inspection is active** and consider excluding Cursor domains
> 3. **Verify streaming works** using the curl tests above
> 4. **Check firewall rules** allow traffic to `*.cursor.sh` and related domains
> 5. **Review proxy logs** for connection errors or timeouts
> 6. **Test from a machine outside your network** to isolate network-specific issues
>
> Most connectivity issues stem from proxies buffering streaming responses. Work with your network team to disable buffering for Cursor domains or implement proper streaming support.
>
> ### Need help with enterprise network setup?
>
> Contact our team for deployment assistance and priority support.
>
>
