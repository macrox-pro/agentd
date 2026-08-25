---
primary_sources:
  - id: T1-ENTERPRISE
    title: "Enterprise"
    url: "https://cursor.com/docs/enterprise.md"
    section: "Enterprise"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Enterprise overview

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Enterprise

> # Enterprise
>
> Cursor provides enterprise-grade security, compliance, and administrative controls for organizations deploying AI-assisted development at scale.
>
> If you are rolling out multiple linked teams, start with [Organizations](https://cursor.com/docs/enterprise/organizations.md).
>
> ## Security and compliance resources
>
> For security reviews and compliance assessments, start with these resources:
>
> - [Trust Center](https://trust.cursor.com/) - Security practices, certifications, and compliance information
> - [Security page](https://cursor.com/security) - Detailed security architecture and controls
> - [Privacy Overview](https://cursor.com/privacy-overview) - Data handling and privacy guarantees
> - [Data Processing Agreement](https://cursor.com/terms/dpa) - GDPR-compliant DPA with data protection commitments
>
> Our certifications include SOC2 Type II, and we maintain GDPR compliance. Visit the [Trust Center](https://trust.cursor.com/) for the latest certification documents and third-party assessment reports.
>
> ## Enterprise documentation
>
> Learn how to deploy, configure, and manage Cursor for your organization. This documentation covers:
>
> - [Admin Setup Guide](https://cursor.com/docs/enterprise/admin-setup-guide.md) - Short checklist of what admins should plan when rolling out Cursor
> - [Security and Privacy Hardening](https://cursor.com/docs/enterprise/security-hardening.md) - One-page checklist of security and privacy controls with links to configure each one
> - [Organizations](https://cursor.com/docs/enterprise/organizations.md) - Org-wide team membership sync and organization groups
> - [Identity & access](https://cursor.com/docs/enterprise/identity-and-access-management.md) - SSO, SCIM, RBAC, and MDM policies
> - [Privacy & data governance](https://cursor.com/docs/enterprise/privacy-and-data-governance.md) - Data flows, Privacy Mode, and data residency
> - [Network configuration](https://cursor.com/docs/enterprise/network-configuration.md) - Proxy setup, IP allowlisting, and encryption
> - [Private connectivity](https://cursor.com/docs/cloud-agent/private-connectivity.md) - AWS PrivateLink and Cloudflare Tunnel for private source control access
> - [Endpoint security](https://cursor.com/docs/enterprise/endpoint-security.md) - Configure antivirus, EDR, and DLP software
> - [LLM safety & controls](https://cursor.com/docs/enterprise/llm-safety-and-controls.md) - Hooks, terminal sandboxing, and agent controls
> - [Models & integrations](https://cursor.com/docs/enterprise/model-and-integration-management.md) - Model controls, MCP, and third-party integrations
> - [Cyber Safeguards](https://cursor.com/docs/account/enterprise/cyber-safeguards.md) - Apply for Anthropic's Cyber Verification Program (CVP) to use eligible Claude models without cyber safeguards
> - [Spend Limits](https://cursor.com/help/account-and-billing/spend-limits.md) - Configure spending limits to control costs
> - [Enterprise help articles](https://cursor.com/help/account-and-billing/enterprise.md) - Common questions about enterprise plans
> - [Compliance & monitoring](https://cursor.com/docs/enterprise/compliance-and-monitoring.md) - Audit logs and tracking
> - [OpenTelemetry Export](https://cursor.com/docs/enterprise/opentelemetry-export.md) - Usage metrics and logs delivered to your observability stack over OTLP
> - [HIPAA Business Associate Agreements](https://cursor.com/docs/enterprise/baa.md) - Request BAA support for Enterprise customers
> - [Deployment patterns](https://cursor.com/docs/enterprise/deployment-patterns.md) - MDM-managed editor vs self-hosted CLI
>
> ## Key features
>
> ### Identity and access
>
> - [SSO and SAML](https://cursor.com/docs/account/teams/sso.md) - Single sign-on for streamlined authentication
> - [SCIM](https://cursor.com/docs/account/teams/scim.md) - Automated user provisioning and deprovisioning
> - [MDM policies](https://cursor.com/docs/enterprise/identity-and-access-management.md#mdm-policies) - Enforce allowed team IDs and extensions on user devices
>
> ### Privacy and security
>
> - [Privacy Mode](https://cursor.com/privacy-overview) - No training on your data by Cursor or other AI providers
> - [Agent Security](https://cursor.com/docs/agent/security.md) - Guardrails for agent tool execution
> - [Hooks](https://cursor.com/docs/hooks.md) - Custom security and compliance workflows
>
> ### Administrative controls
>
> - [Dashboard](https://cursor.com/docs/account/teams/dashboard.md) - Team management, settings, and monitoring
> - [Admin API](https://cursor.com/docs/account/teams/admin-api.md) - Programmatic access to admin features
> - [Analytics](https://cursor.com/docs/account/teams/analytics.md) - Usage metrics and insights
> - [Conversation Insights](https://cursor.com/docs/account/teams/analytics.md#conversation-insights) - Understand the type of work being done with Cursor (Enterprise only)
> - [AI Code Tracking API](https://cursor.com/docs/account/teams/ai-code-tracking-api.md) - Per-commit AI usage metrics (Enterprise only)
> - [Cursor Blame](https://cursor.com/docs/integrations/cursor-blame.md) - AI-aware git blame that shows AI vs human code attribution (Enterprise only)
> - [Analytics API](https://cursor.com/docs/account/teams/analytics-api.md) - Usage metrics and insights
> - [Billing Groups](https://cursor.com/docs/account/enterprise/billing-groups.md) - Manage spend across groups of users for reporting and chargebacks (Enterprise only)
> - [Service Accounts](https://cursor.com/docs/account/enterprise/service-accounts.md) - Non-human accounts for automated workflows (Enterprise only)
>
> ### Models and integrations
>
> - [Models](https://cursor.com/docs/models-and-pricing.md) - Available models and configuration
> - [Cyber Safeguards](https://cursor.com/docs/account/enterprise/cyber-safeguards.md) - Anthropic Cyber Verification Program (CVP) access for security groups (Enterprise only)
> - [MCP](https://cursor.com/docs/mcp.md) - Model Context Protocol server trust management
> - [Slack](https://cursor.com/docs/integrations/slack.md) - Cloud Agents in Slack
> - [GitHub](https://cursor.com/docs/integrations/github.md) - Repository integration
> - [Linear](https://cursor.com/docs/integrations/linear.md) - Issue tracking integration
> - [Bugbot](https://cursor.com/docs/bugbot.md) - Automated bug detection and fixing
>
> ### Monitoring and compliance
>
> - Audit logs - Track authentication, user management, and administrative actions (Enterprise only)
> - SIEM integration - Stream audit logs to your security tools
> - [OpenTelemetry Export](https://cursor.com/docs/enterprise/opentelemetry-export.md) - Stream usage metrics and logs to your observability stack over OTLP (Enterprise only)
> - [HIPAA Business Associate Agreements](https://cursor.com/docs/enterprise/baa.md) - BAA support for Enterprise customers
>
> ## Getting started
>
> Start with the [Admin Setup Guide](https://cursor.com/docs/enterprise/admin-setup-guide.md) for a short rollout checklist. In short:
>
> 1. Review the [Trust Center](https://trust.cursor.com/) and [Security page](https://cursor.com/security) for your security assessment
> 2. Read through the [enterprise documentation](https://cursor.com/docs/enterprise.md) to understand deployment options
> 3. Set up [SSO](https://cursor.com/docs/account/teams/sso.md) and [SCIM](https://cursor.com/docs/account/teams/scim.md) for user management
> 4. Deploy Cursor and configure [MDM policies](https://cursor.com/docs/enterprise/deployment-patterns.md#mdm-configuration) to enforce team IDs and extensions
> 5. Review the [Dashboard](https://cursor.com/docs/account/teams/dashboard.md) to monitor team usage
>
> ## Plan Comparison
>
> ### Team Admin & Billing
>
> | Capability                                                                                                          | Individual Plans | Teams                                                                     | Enterprise                                                                                                                                                                                                                                                        |
> | ------------------------------------------------------------------------------------------------------------------- | ---------------- | ------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | Centralized Billing                                                                                                 |                  | ✓                                                                         | ✓                                                                                                                                                                                                                                                                 |
> | Usage Spend Controls                                                                                                | Personal limits  | Team limits                                                               | [Pooled usage + admin-only controls](https://cursor.com/help/account-and-billing/spend-limits.md#what-types-of-spend-limits-are-available)                                                                                                                        |
> | [Billing Groups](https://cursor.com/docs/account/enterprise/billing-groups.md)                                      |                  |                                                                           | ✓                                                                                                                                                                                                                                                                 |
> | [Team Usage Analytics](https://cursor.com/docs/account/teams/analytics.md)                                          |                  | [Analytics Dashboard](https://cursor.com/docs/account/teams/analytics.md) | [Analytics Dashboard](https://cursor.com/docs/account/teams/analytics.md),[AI Code Tracking API](https://cursor.com/docs/account/teams/ai-code-tracking-api.md),[Conversation Insights](https://cursor.com/docs/account/teams/analytics.md#conversation-insights) |
> | [Cursor Blame](https://cursor.com/docs/integrations/cursor-blame.md)                                                |                  |                                                                           | ✓                                                                                                                                                                                                                                                                 |
> | [SSO (SAML/OIDC)](https://cursor.com/docs/enterprise/identity-and-access-management.md#single-sign-on-sso-and-saml) |                  | ✓                                                                         | ✓                                                                                                                                                                                                                                                                 |
> | [SCIM Provisioning](https://cursor.com/docs/account/teams/scim.md)                                                  |                  |                                                                           | ✓                                                                                                                                                                                                                                                                 |
> | [Audit Logs](https://cursor.com/docs/enterprise/compliance-and-monitoring.md#audit-logs)                            |                  |                                                                           | ✓                                                                                                                                                                                                                                                                 |
> | [Service Accounts](https://cursor.com/docs/account/enterprise/service-accounts.md)                                  |                  |                                                                           | ✓                                                                                                                                                                                                                                                                 |
>
> ### Marketplace
>
> | Capability                        | Individual Plans | Teams                 | Enterprise                                                                                                         |
> | --------------------------------- | ---------------- | --------------------- | ------------------------------------------------------------------------------------------------------------------ |
> | Team marketplaces                 |                  | 1 team marketplace    | Unlimited team marketplaces                                                                                        |
> | Community plugin import           |                  | On by default         | Off by default                                                                                                     |
> | Marketplace edits                 |                  | All teams can edit    | Only admin edits                                                                                                   |
> | SCIM distribution & access gating |                  | No SCIM access gating | Scope distribution and gate access via SCIM (control who sees which marketplace based on identity provider groups) |
>
> ### Centralized Agent Controls
>
> | Capability                                                                                                              | Individual Plans | Teams                  | Enterprise                                                                           |
> | ----------------------------------------------------------------------------------------------------------------------- | ---------------- | ---------------------- | ------------------------------------------------------------------------------------ |
> | [Privacy Mode](https://cursor.com/docs/enterprise/privacy-and-data-governance.md#privacy-mode-enforcement)              | User choice      | Enforce org-wide       | Enforce org-wide                                                                     |
> | [Team Rules](https://cursor.com/docs/rules.md#team-rules)                                                               |                  | Enforceable + Optional | Enforceable + Optional                                                               |
> | [Hooks for Logging,Auditing, and more](https://cursor.com/docs/hooks.md)                                                | ✓                | MDM Distribution       | [MDM & Server-side distribution](https://cursor.com/docs/hooks.md#team-distribution) |
> | [Agent Sandbox Mode](https://cursor.com/docs/agent/security/run-modes.md#sandboxing)                                    | ✓                | ✓                      | Enforce org-wide                                                                     |
> | [Repository Blocklist](https://cursor.com/docs/enterprise/model-and-integration-management.md#git-repository-blocklist) |                  |                        | ✓                                                                                    |
> | [Model Access Restrictions](https://cursor.com/docs/enterprise/model-and-integration-management.md)                     |                  |                        | ✓                                                                                    |
> | [Auto-run, Browser, and Network Controls](https://cursor.com/docs/enterprise/llm-safety-and-controls.md)                |                  |                        | ✓                                                                                    |
>
> ### User Access Controls
>
> | Capability   | Individual & Teams Plans | Enterprise                                     |
> | ------------ | ------------------------ | ---------------------------------------------- |
> | Cursor CLI   |                          | Restrict which users can access agents via CLI |
> | Cloud Agents |                          | Restrict which users can create Cloud Agents   |
> | Analytics    |                          | Restrict analytics dashboard to admins only    |
> | BYOK         |                          | Disable users from using their own API keys    |
>
> ### Support & Legal
>
> | Capability        | Individual Plans                                          | Teams                                                     | Enterprise                                                          |
> | ----------------- | --------------------------------------------------------- | --------------------------------------------------------- | ------------------------------------------------------------------- |
> | Technical Support | [Community & Standard Support](https://forum.cursor.com/) | [Community & Standard Support](https://forum.cursor.com/) | First human response times: 8 hours (critical), 24 hours (standard) |
> | Terms             | [Online Terms](https://cursor.com/terms-of-service)       | [MSA & DPA](https://cursor.com/terms/msa)                 | [MSA & DPA](https://cursor.com/terms/msa)                           |
>
> For security vulnerabilities, see our [responsible disclosure program](https://cursor.com/docs/agent/security.md#responsible-disclosure).
>
> ### Ready to deploy Cursor at scale?
>
> Contact our team to discuss your organization's needs.
>
>

### Source: LLM safety and controls

> # LLM Safety and Controls
>
> AI models can behave unexpectedly. This documentation covers how to control what agents can do, set up safety guardrails, and guide LLM behavior toward desired outcomes.
>
> ## Understanding model behavior
>
> LLMs generate text based on probability distributions, not by retrieving facts from a database or executing deterministic logic. They can produce different outputs for the same input, hallucinate facts or code that seems plausible but is wrong, and be influenced by carefully crafted prompts (prompt injection).
>
> You can't rely on LLMs to always make safe decisions. Instead, you combine two approaches: **security controls** that enforce hard boundaries on what agents can do, and **steering mechanisms** that guide LLM behavior toward better outcomes.
>
> For a deeper understanding of how LLMs work, see [How AI Models Work](https://cursor.com/learn/how-ai-models-work.md).
>
> ## Two approaches to safety
>
> Cursor provides two complementary approaches to managing AI agent behavior:
>
> **Security controls (deterministic enforcement)**: Hard boundaries that block dangerous operations regardless of what the LLM suggests. These include terminal command restrictions, enforcement hooks that reject operations, approval workflows, and sandboxing. Security controls are your primary defense against harmful agent actions.
>
> **LLM steering (non-deterministic guidance)**: Mechanisms that guide the LLM toward better behavior by shaping its context and available actions. These include Rules that add instructions to prompts, Commands that provide reusable workflows, and integrations that enrich the agent's knowledge. Steering improves agent quality but doesn't guarantee prevention of harmful actions.
>
> Use both approaches together. Security controls provide the safety net. Steering reduces how often agents attempt problematic actions in the first place.
>
> ## Security controls
>
> These deterministic controls enforce hard boundaries on what agents can do. They work regardless of what the LLM suggests.
>
> ### Terminal command restrictions
>
> By default, Cursor requires your approval before executing any terminal command. This protects against destructive commands (deleting files, dropping databases), commands that expose sensitive data, and commands with unintended side effects.
>
> When an agent wants to run a command, you see a prompt showing the full command. You can approve and run it, deny it, or modify it before running.
>
> #### Auto-approval risks
>
> You can enable auto-approval for terminal commands, but understand the risks. Agents might run destructive commands without your knowledge, commands execute before you can review them, and bugs or prompt injection could cause unintended operations.
>
> #### Run Mode configuration
>
> Enterprise teams can configure Run Mode policies in the team dashboard. In Cursor 3.6 and above, end users choose between **Auto-review** (the default), **Allowlist**, and **Run Everything** modes. **Auto-review** runs allowlisted calls, sandboxes shell commands when it can, and routes the rest through an LLM classifier that returns allow or block based on safety and how well the call matches the user's intent. You can create an allowlist of commands that don't require approval, such as `npm install`, `pip install`, `cargo build`, or `make test`.
>
> The allowlist is best-effort, not a security boundary. Determined agents or prompt injection might bypass it. Always combine allowlists with other security controls like hooks.
>
> See [Run Modes](https://cursor.com/docs/agent/security/run-modes.md#run-mode) and [Agent Security](https://cursor.com/docs/agent/security.md) for details.
>
> ### Enforcement hooks
>
> Hooks let you run custom logic at key points in the agent loop.
>
> - Before prompt submission: Scan prompts for sensitive data before they're sent to LLMs. Block submissions that contain API keys or credentials, personal identifiable information (PII), or proprietary information.
> - Before file reading: Scan files before agents read them. Redact or block access to configuration files with secrets, PII in databases or logs, or proprietary algorithms.
> - After code generation: Scan generated code before it's written to disk. Check for security vulnerabilities (SQL injection, XSS), licensed code that might cause IP issues, or API keys and credentials in code.
> - Before terminal execution: Block dangerous commands or route them through approval workflows. For example, block all `git push` commands, require approval for any `sudo` command, or block database `DROP` statements.
>
> #### Example: Blocking git commands
>
> This hook intercepts shell commands and blocks raw git usage, directing users to the GitHub CLI instead:
>
> ```bash
> #!/bin/bash
> input=$(cat)
> command=$(echo "$input" | jq -r '.command')
>
> if [[ "$command" =~ git[[:space:]] ]]; then
>     cat << EOF
> {
>   "permission": "deny",
>   "userMessage": "Git command blocked. Please use gh tool instead.",
>   "agentMessage": "Use 'gh' commands instead of raw git."
> }
> EOF
> fi
> ```
>
> #### Example: Redacting secrets
>
> This hook scans file contents for GitHub API keys and blocks access if found:
>
> ```bash
> #!/bin/bash
> input=$(cat)
> content=$(echo "$input" | jq -r '.content')
>
> if echo "$content" | grep -qE 'gh[ps]_[A-Za-z0-9]{36}'; then
>     cat << EOF
> {
>   "permission": "deny"
> }
> EOF
>     exit 3
> fi
> ```
>
> See [Hooks](https://cursor.com/docs/hooks.md) for complete documentation and more examples.
>
> ### Protecting sensitive files
>
> Not all files in your repositories should be accessible to AI. Configuration files, secrets, and sensitive data need protection.
>
> #### .cursorignore
>
> The `.cursorignore` file works like `.gitignore` but controls what Cursor can access. Files matching patterns in `.cursorignore` are excluded from:
>
> - Agent file reading
> - Context selection
>
> `.cursorignore` is not a security boundary. It's a convenience feature to exclude files from AI processing, but:
>
> - Users can manually read ignored files
> - Agents might find ways to access ignored content
> - It doesn't prevent file access, only excludes from indexing
>
> For true security, use file system permissions or encrypt sensitive data.
>
> See [Ignore Files](https://cursor.com/docs/reference/ignore-file.md) for detailed syntax.
>
> #### .cursor directory protection
>
> The `.cursor` directory in repositories contains project-specific settings, rules, and cache files. Enterprise teams can prevent agents from modifying this directory.
>
> When enabled, agents cannot:
>
> - Modify files in `.cursor/`
> - Delete the `.cursor/` directory
> - Change cursor rules or settings files
>
> Users can still manually edit these files, but agents require approval.
>
> Configure in the [team dashboard](https://cursor.com/docs/account/teams/dashboard.md) under ".cursor Directory Protection" (Enterprise only).
>
> #### Browser origin controls
>
> Enterprise teams can restrict which websites agents can navigate to when using the [browser tool](https://cursor.com/docs/agent/tools/browser.md). Define an allowlist of approved domains—agents attempting to visit other origins are blocked.
>
> Configure in the [team dashboard](https://cursor.com/docs/account/teams/dashboard.md) under "Browser Controls" (Enterprise only).
>
> ### Integration with DLP tools
>
> Many enterprises have existing Data Loss Prevention (DLP) tools that scan for sensitive data. You can integrate Cursor with your DLP tools in three ways.
>
> #### Endpoint DLP agents
>
> Most endpoint DLP software can inspect Cursor's network traffic. Configure your DLP to monitor traffic to `*.cursor.sh` domains, scan for sensitive patterns in outbound requests, and block or alert on policy violations.
>
> Network DLP may impact performance. See [Network Configuration](https://cursor.com/docs/enterprise/network-configuration.md) for proxy considerations.
>
> #### Hooks-based DLP
>
> Use Cursor's hooks feature to implement custom DLP logic:
>
> **Before prompt submission:**
> Scan prompts for sensitive patterns before sending to LLMs:
>
> ```bash
> #!/bin/bash
> input=$(cat)
> prompt=$(echo "$input" | jq -r '.prompt')
>
> # Check for API keys
> if echo "$prompt" | grep -qE 'api[_-]?key.*[A-Za-z0-9]{32}'; then
>     cat << EOF
> {
>   "continue": false,
>   "userMessage": "Prompt contains what looks like an API key. Remove it and try again."
> }
> EOF
>     exit 1
> fi
>
> # Allow if no sensitive data found
> cat << EOF
> {
>   "continue": true
> }
> EOF
> ```
>
> **After code generation:**
> Scan generated code before it's written to disk:
>
> ```bash
> #!/bin/bash
> input=$(cat)
> file_path=$(echo "$input" | jq -r '.file_path')
> edits=$(echo "$input" | jq -r '.edits[].new_string')
>
> # Check for hardcoded credentials
> if echo "$edits" | grep -qE 'password.*=.*["\047][^"\047]+["\047]'; then
>     # Send to your DLP API for analysis
>     curl -X POST "https://dlp.yourcompany.com/scan" \
>       -H "Content-Type: application/json" \
>       -d "{\"content\":\"$edits\",\"file\":\"$file_path\"}"
>     
>     # Check API response and act accordingly
> fi
> ```
>
> #### Third-party DLP integration
>
> Call your existing DLP vendor's API from hooks:
>
> ```bash
> #!/bin/bash
> input=$(cat)
> content=$(echo "$input" | jq -r '.content')
>
> # Send to DLP API
> response=$(curl -s -X POST "https://dlp-api.company.com/analyze" \
>   -H "Authorization: Bearer $DLP_API_TOKEN" \
>   -H "Content-Type: application/json" \
>   -d "{\"text\":\"$content\"}")
>
> # Parse response
> is_allowed=$(echo "$response" | jq -r '.allowed')
>
> if [ "$is_allowed" = "true" ]; then
>     cat << EOF
> {
>   "permission": "allow"
> }
> EOF
> else
>     violation=$(echo "$response" | jq -r '.violation_type')
>     cat << EOF
> {
>   "permission": "deny",
>   "userMessage": "Content blocked by DLP policy: $violation"
> }
> EOF
> fi
> ```
>
> This approach gives you centralized DLP policy management across all development tools.
>
> ### Approval workflows
>
> You can configure Cursor to ask for approval on every agent action. Users can set their agent to always ask before reading files, editing files, running terminal commands, or making network requests.
>
> However, this approach significantly slows down the development experience. Agents need multiple actions to complete tasks, and requiring approval for each action makes the workflow tedious. Most teams instead choose to use hooks to block dangerous operations automatically.
>
> ### Model provider safety
>
> All model providers (OpenAI, Anthropic, Google, SpaceXAI) implement safety systems that filter harmful content. These systems reject prompts requesting harmful information, refuse to generate dangerous code, and filter outputs for safety.
>
> Cursor works with providers to ensure models meet safety standards before deployment to users. Providers continuously evaluate models for safety issues. However, these are not security boundaries. Safety systems can be bypassed or tricked. Always implement your own controls through hooks and access policies.
>
> ### Sandboxing considerations
>
> Cursor agents run on your local machine by default. They can read files you can read, write files you can write, execute commands you can execute, and access network resources you can access.
>
> There is no security boundary between agents and your user account. If your account can delete files, agents can delete files (with approval by default).
>
> #### Sandboxing options
>
> If you need stronger isolation, run Cursor in a separate VM using Cloud Agents, use file system permissions to limit what the Cursor process can access, or run Cursor on a dedicated development machine with limited access to production systems.
>
> For most enterprises, the built-in approval requirements and hooks provide sufficient control.
>
> ### File system permissions
>
> For further defense, use file system permissions to protect sensitive files:
>
> **Restrict access to secret files:**
>
> ```bash
> # Make secrets readable only by specific users
> chmod 600 .env
> chown app-user:app-user .env
>
> # Or use separate directories with restricted access
> chmod 700 /etc/app/secrets
> ```
>
> **Separate sensitive repos:**
> Keep highly sensitive code in separate repositories with restricted access. Don't clone these repositories to machines where Cursor runs.
>
> **Encrypted filesystems:**
> For very sensitive data, use encrypted filesystems that require explicit mounting. Don't mount these filesystems in directories where Cursor has access.
>
> ## LLM steering
>
> Security controls block harmful actions after the LLM suggests them. Steering mechanisms guide the LLM to make better suggestions in the first place. These are non-deterministic. They improve outcomes but don't guarantee prevention.
>
> ### Rules
>
> Rules add instructions to the LLM's context window before every request. Use rules to establish coding standards, enforce architectural patterns, set security requirements, or define project-specific conventions.
>
> Rules work at three scopes:
>
> **User rules**: Apply to all projects for a specific user. Use these for personal preferences like code style or preferred libraries.
>
> **Project rules**: Apply to everyone working on a project. Use these for project-specific standards like naming conventions or framework usage.
>
> **Team rules**: Apply to all projects in your organization. Use these for company-wide standards like security requirements or compliance rules.
>
> The LLM sees all applicable rules when generating responses. It will attempt to follow them, but rules are suggestions, not guarantees. Combine rules with enforcement hooks for requirements that must be followed.
>
> See [Rules](https://cursor.com/docs/rules.md) for configuration and examples.
>
> ### Commands and workflows
>
> Commands package reusable prompts that agents can invoke with slash commands like `/test` or `/deploy`. Commands help standardize common workflows across your team.
>
> **Workflows**: Create multi-step processes that guide agents through complex tasks. For example, a `/security-review` command might instruct the agent to scan for SQL injection, check for exposed secrets, validate input sanitization, and generate a security report.
>
> **Prompt libraries**: Build a collection of tested prompts for common tasks. This reduces variation in agent behavior and captures institutional knowledge.
>
> Commands are scoped to teams, projects, or users. Team admins can create organization-wide commands that appear for all developers.
>
> See [Commands](https://cursor.com/help/customization/rules.md) for configuration and examples.
>
> ### Context enrichment with MCPs
>
> Model Context Protocol (MCP) servers let agents access external data sources. Use MCPs to pull in company documentation, query internal APIs, access knowledge bases, or integrate with development tools.
>
> MCPs enrich the agent's context with information it wouldn't otherwise have. For example, an MCP might provide access to your API specifications, so agents can generate code that correctly calls your internal services.
>
> MCPs are scoped to teams or users. Unlike hooks, MCPs don't enforce policies—they provide information that helps agents make better decisions.
>
> See [MCP Integration](https://cursor.com/docs/mcp.md) for configuration and examples.
>
> ### Advanced safety controls for Enterprise
>
> Contact our team to learn about org-wide enforcement and security policies.
>
>
