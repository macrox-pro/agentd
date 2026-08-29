---
primary_sources:
  - id: T4-MIGRATION
    title: "Migrating from Gemini CLI"
    url: "https://antigravity.google/docs/cli/gcli-migration"
    section: "Full page (HTML → text)"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Gemini CLI to Antigravity migration

> **Applicability:** Verbatim excerpts (snapshot 2026-08-29). HTML sources converted to plain text where noted.

### Source: Migrating from Gemini CLI — Full page

> Migration | Google Antigravity Docs
>
>  
>
>  
>
> Skip to content
>
> Google Antigravity Docs
>
> Search
>
> Ctrl
>
> K
>
> Cancel
>
>  
>
> Select theme
>
> Dark
>
> Light
>
> Auto
>
>  
>
>  
>
>  
>
> Home
>
> Antigravity 2.0
>
> v2.11.0
>
> Overview
>
> Getting Started
>
> Build with Google
>
> Feature Overview
>
> Models
>
> Projects
>
> Remote Control
>
> Settings
>
> Overview
>
> Agent Settings
>
> Artifact Review
>
> Customizations
>
> MCP
>
> Skills
>
> Rules
>
> Plugins
>
> Hooks
>
> Sidecars
>
> Agent Capabilities
>
> Slash commands overview
>
> Teamwork agent teams (/teamwork-preview)
>
> Permissions
>
> Subagents
>
> Artifacts
>
> Overview
>
> Plan
>
> Walkthrough
>
> Screenshots
>
> Antigravity CLI
>
> v1.1.22
>
> Overview
>
> Getting Started
>
> Installation &amp; Auth
>
> Tutorial
>
> Using AGY CLI
>
> Features
>
> Migration
>
> Prompting &amp; Interaction
>
> Artifacts
>
> Reviewing Artifacts
>
> Managing Conversations
>
> Agent Capabilities
>
> Choose an execution mode
>
> Headless Mode
>
> Background Tasks &amp; Subagents
>
> Sandbox
>
> Permissions
>
> Projects
>
> Settings
>
> Settings, Rendering &amp; Keybindings
>
> Vim Editor Mode
>
> AI Credits
>
> Customizations
>
> MCP
>
> Plugins &amp; Skills
>
> Status Line Customization
>
> Terminal Title Customization
>
> Commands
>
> Agents Command (/agents)
>
> Code Search Command (/codesearch)
>
> AI Credits Command (/credits)
>
> Diff Command (/diff)
>
> Permissions Command (/permissions)
>
> Resume Command (/resume)
>
> Status Line Command (/statusline)
>
> Teamwork Command (/teamwork-preview)
>
> Window Title Command (/title)
>
> Model Quotas (/usage)
>
> Voice Dictation (/voice)
>
> Best Practices
>
> Troubleshooting
>
> CLI Reference
>
> Antigravity SDK
>
> v0.1.15
>
> Overview + Quick Start
>
> Personas
>
> Tools &amp; Skills
>
> Customizations
>
> MCP
>
> Policies
>
> Subagents
>
> Structured Output
>
> Lifecycle &amp; Hooks
>
> Antigravity for IDEs
>
> v2.5.5
>
> Overview
>
> Getting Started
>
> IDE Extensions
>
> Overview
>
> Visual Studio Code
>
> Visual Studio
>
> JetBrains
>
> Zed
>
> Xcode
>
> Features
>
> Tab
>
> Side Panel
>
> Review Changes
>
> Artifacts
>
> Plan
>
> Walkthrough
>
> Screenshots
>
> Browser Recordings
>
> Browser
>
> Overview
>
> Allowlist / Denylist
>
> Separate Chrome Profile
>
> Customizations
>
> MCP
>
> Skills
>
> Rules
>
> Workflows
>
> Plugins
>
> Hooks
>
> Settings
>
> Migration
>
> Firebase Studio Migration
>
> Enterprise
>
> Plans
>
> FAQ
>
>  
>
> Select theme
>
> Dark
>
> Light
>
> Auto
>
>  
>
>  
>
> On this page
>
> Overview
>
> Overview
>
> First-launch onboarding
>
> Converting extensions to plugins
>
> Expected import output
>
> Context files and workspace rules
>
> Updated skills paths
>
> MCP config formatting changes
>
> Directory mapping
>
> Required schema updates
>
> Next steps
>
>  
>
> On this page
>
> Overview
>
> Overview
>
> First-launch onboarding
>
> Converting extensions to plugins
>
> Expected import output
>
> Context files and workspace rules
>
> Updated skills paths
>
> MCP config formatting changes
>
> Directory mapping
>
> Required schema updates
>
> Next steps
>
>  
>
> Antigravity CLI
>
> Migration
>
> Markdown
>
> keyboard_arrow_down
>
> content_copy
>
> Copy Markdown
>
> open_in_new
>
> View Markdown
>
>  
>  
>
> Migrating from Gemini CLI
>
> Section titled “Migrating from Gemini CLI”
>
> Convert your legacy configurations, import Gemini CLI extensions as native plugins, adapt custom skills paths, and reformat Model Context Protocol configurations.
>
> Overview
>
> Section titled “Overview”
>
> Antigravity CLI preserves backward compatibility with the core developer-experience constructs popularized by Gemini CLI. To ensure a seamless upgrade, the CLI offers automatic onboarding conversion alongside explicit CLI migration command sequences.
>
> First-launch onboarding
>
> Section titled “First-launch onboarding”
>
> When you execute 
> agy
>  for the first time in an environment containing legacy configurations, the CLI automatically detects your existing profiles. An interactive checklist prompts you to choose which assets to migrate:
>
> Auto-conversion
> : Select the extensions and global configurations you wish to convert.
>
> Keyring storage
> : The CLI migrates your active session tokens securely into your operating system’s native keyring storage.
>
> Settings alignment
> : Default visual parameters and rendering buffers map automatically to your new settings profile.
>
> Note
>
> Partial Parity
> : While we preserve support for workspace skills, rules, and MCP servers, certain customized terminal themes or experimental visual overlays from Gemini CLI may not be supported.
>
> Converting extensions to plugins
>
> Section titled “Converting extensions to plugins”
>
> Since Gemini CLI launched, the industry has standardized on the term 
> plugins
> . You can manually convert your legacy Gemini extensions to native Antigravity plugins by executing:
>
> agy
>
>  plugin
>
>  import
>
>  gemini
>
> This utility searches your legacy local directories, parses your extension manifests, and converts files into native layout blocks.
>
> Expected import output
>
> Section titled “Expected import output”
>
> [ok] conductor-tools
>
>  - skills : skipped (none detected)
>
>  - agents : skipped (none detected)
>
>  ✔ commands : 4 legacy commands converted to skills
>
>  - mcpServers : skipped (none detected)
>
> [ok] google-workspace
>
>  ✔ skills : 5 skills processed
>
>  - agents : skipped (none detected)
>
>  ✔ commands : 2 legacy commands converted to skills
>
>  ✔ mcpServers : 1 server definition migrated to mcp_config.json
>
> Context files and workspace rules
>
> Section titled “Context files and workspace rules”
>
> Both CLI platforms utilize identical workspace context rules. No modifications are needed to your existing rule documents:
>
> Workspace local context
> : The agent continues to parse and enforce rule constraints defined inside your active directory’s 
> GEMINI.md
>  and 
> AGENTS.md
>  files.
>
> Global developer context
> : The agent automatically consults and enforces your global constraints located at 
> ~/.gemini/GEMINI.md
> .
>
> Updated skills paths
>
> Section titled “Updated skills paths”
>
> While global shared skills remain in your user home directory, the target folder path for local workspace-specific skills has been updated.
>
> Configuration
>
> Gemini CLI
>
> Antigravity CLI
>
> Global shared path
>
> ~/.gemini/skills/
>
> ~/.gemini/antigravity-cli/skills/
>
> Workspace project path
>
> .gemini/skills/
>
> .agents/skills/
>
> Note
>
> Action Required
> : If your project contains custom workspace skills defined in 
> .gemini/skills/
> , you must manually rename or relocate the folder to 
> .agents/skills/
>  for the Antigravity agent to recognize them as active slash commands.
>
> MCP config formatting changes
>
> Section titled “MCP config formatting changes”
>
> Antigravity CLI separates Model Context Protocol servers into dedicated, lightweight JSON profiles instead of nesting them inside your primary preferences configuration.
>
> Directory mapping
>
> Section titled “Directory mapping”
>
> Legacy Gemini Config
> : Servers were declared inline within 
> ~/.gemini/settings.json
> .
>
> Antigravity CLI Config
> : Servers are defined inside a standalone 
> mcp_config.json
>  profile:
>
> Global servers: 
> ~/.gemini/config/mcp_config.json
>
> Workspace servers: 
> .agents/mcp_config.json
>
> Required schema updates
>
> Section titled “Required schema updates”
>
> When manually migrating remote websocket or SSE server definitions, update the URI key parameter to match the current standard:
>
> Legacy schema keys
> : 
> url
>  or 
> httpUrl
>
> Modern schema key
> : 
> serverUrl
>
> {
>
>  "mcpServers"
>
> : {
>
>  "remote-indexer"
>
> : {
>
>  "serverUrl"
>
> : 
>
> "https://mcp.internal.enterprise.com/sse"
>
> ,
>
>  "env"
>
> : {
>
>  "AUTH_TOKEN"
>
> : 
>
> "secure_alpha_token"
>
>  }
>
>  }
>
>  }
>
> }
>
> Next steps
>
> Section titled “Next steps”
>
> Begin configuring your new visual parameters and troubleshooting any setup anomalies:
>
> Settings, Rendering &#x26; Keybindings
>
> : Customize keyboard hotkeys, themes, and screen buffers.
>
> Troubleshooting
>
> : Learn how to resolve authentication lockouts or path issues.
>
> CLI Reference
>
> : Access standard parameters lists and slash command mappings.
>
> Previous
>
> Features
>
> Next
>
> Prompting &amp; Interaction
