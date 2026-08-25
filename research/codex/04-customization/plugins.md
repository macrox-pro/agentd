---
primary_sources:
  - id: T1-PLUGINS
    title: "Plugins"
    url: "https://learn.chatgpt.com/docs/plugins.md"
    section: "full page"
  - id: T1-BUILD-PLUGINS
    title: "Build plugins"
    url: "https://learn.chatgpt.com/docs/build-plugins.md"
    section: "full page"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Plugins

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Plugins — full page

> # Plugins
>
> > For the complete documentation index, see [llms.txt](https://learn.chatgpt.com/llms.txt). Markdown versions of documentation pages are available by appending `.md` to the page URL.
>
> ## Overview
>
> Plugins bundle capabilities into reusable workflows in ChatGPT and Codex. They
> can include skills, connectors, or both. Both products use one universal plugin
> directory, so the same public plugins are discoverable from their supported
> surfaces.
>
> Plugins work in Chat and Work across ChatGPT on the web, desktop, and mobile,
> and in Codex in the ChatGPT desktop app. Codex CLI also has a plugin browser
> for Codex environments. The IDE extension doesn't support plugins.
>
> On mobile, you can use plugins available to your account in Chat or Work.
>
> <ContentModeSwitch group="codex-surface" id="app">
>
> Open the **Plugins** tab to browse and install plugins. After installation, you
> can use plugins in Chat or Work in ChatGPT, or in Codex. Installed plugins can
> add skills, connectors, and MCP tools to new chats.
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="web">
>
> Open the **Plugins** tab to browse and install plugins. After installation, you
> can use plugins in Chat or Work. A plugin can prompt you to connect an external
> service before its tools become available.
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="cli">
>
> In Codex CLI, enter `/plugins` to open the plugin browser. Install a plugin from
> a configured marketplace, then start a new session before using its bundled
> skills or tools.
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="ide">
>
> <a id="plugin-directory-in-the-ide-extension"></a>
>
> ### Use plugins from a supported surface
>
> Plugins aren't available in the IDE extension. To browse and install plugins
> for Codex, use the ChatGPT desktop app or Codex CLI.
>
> </ContentModeSwitch>
>
> Extend what ChatGPT and Codex can do, for example:
>
> - Install the Codex Security plugin to scan authorized code and confirm
>   plausible vulnerability findings.
> - Install the Gmail plugin to work with Gmail.
> - Install the Google Drive plugin to work across Drive, Docs, Sheets, and
>   Slides.
> - Install the Slack plugin to summarize channels or draft replies.
>
> A plugin can contain one or more of these parts:
>
> - **Skills:** reusable instructions for specific kinds of work. ChatGPT and
>   Codex can load them when needed so they follow the right steps and use the
>   right references or helper scripts for a task.
> - **Connectors:** connections to tools like GitHub, Slack, or Google Drive, so
>   ChatGPT and Codex can read information from those tools and take actions in
>   them. Connectors expose tools and can optionally include custom UI.
> - **MCP servers:** services that give ChatGPT and Codex access to more tools or
>   shared information, often from systems outside your local project. They're
>   also the services behind connectors. They define tools, enforce auth, return
>   structured data, and perform actions against external systems.
> - **Browser extensions:** browser capabilities that a plugin needs for its
>   workflow.
> - **Hooks:** commands that run at configured lifecycle points. Review and trust
>   plugin hooks before you enable them.
> - **Scheduled task templates:** reusable starting points for recurring tasks
>   where scheduled tasks are available.
>
> You can share plugins by publishing them through a marketplace source, such as a
> repo marketplace for a project or team. See [Build plugins](https://developers.openai.com/plugins/build/plugins)
> for marketplace setup, packaging, and distribution guidance.
>
> If you are building an integration, start with
> [Build an MCP server](https://developers.openai.com/plugins/build/mcp-server).
> If the plugin needs custom UI, use the
> [optional UI guide](https://developers.openai.com/plugins/build/chatgpt-ui).
>
> ## Use and install plugins
>
> <a id="plugin-directory-in-the-codex-app"></a>
>
> <ContentModeSwitch group="codex-surface" ids="app,web">
>
> ### Universal plugin directory
>
> ChatGPT and Codex use the same public plugin catalog. On the web or in the
> ChatGPT desktop app, open the **Plugins** tab to browse and install plugins.
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="app">
>
>
>   
>
> > Illustration: Plugins Directory in the ChatGPT desktop app
>
>
>
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" ids="app,web">
>
> The Plugins Directory organizes plugins into tabs:
>
> - **OpenAI:** plugins built by OpenAI.
> - **Your workspace name:** plugins provided by your workspace.
> - **Personal:** personal marketplace plugins, including **Created by me** and
>   **Shared with me** sections when those plugins are available.
>
> Use the separate **Installed** row to review plugins you already installed.
>
> ### Install and use a plugin
>
> Once you open the Plugins Directory:
>
> <WorkflowSteps>
>
> 1. Search or browse for a plugin, then open its details.
> 2. Select the plus button to install the plugin.
> 3. If the plugin needs a connector, connect it when prompted. Some plugins
>    ask you to authenticate during install. Others wait until the first time you
>    use them.
> 4. After installation, start a new chat and ask ChatGPT or Codex to use the
>    plugin.
>
> </WorkflowSteps>
>
> ### Connect supported partners with Sign in with ChatGPT
>
> **Sign in with ChatGPT** is rolling out in beta for supported plugins and
> partner sites, including Airtable, GitLab, HubSpot, Notion, Supabase, and
> Vercel. When the option is available, select **Sign in with ChatGPT** while
> connecting the plugin to create or link your account with that service.
>
> Signing in shares only your name, email address, and profile picture, when
> available, with the partner. It doesn't grant the plugin access to your data or
> approve actions automatically. Review and approve the plugin's requested
> permissions as a separate step before using the connection.
>
> After you install a plugin, you can use it directly in the prompt window:
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="app">
>
>
>   
>
> > Illustration: Installed plugin on the Plugins page
>
>
>
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" ids="app,web">
>
>
>
>   
>
>     
> Describe the task directly
>
>     
>
>       Ask for the outcome you want, such as "Summarize unread Gmail threads
>       from today" or "Pull the latest launch notes from Google Drive."
>     
>
>     
>
>       Use this when you want ChatGPT to choose the right installed tools for the
>       task.
>     
>
>   
>
>
>   
>
>     
> Choose a specific plugin
>
>     
>
>       Type `@` to invoke the plugin or one of its bundled skills
>       explicitly.
>     
>
>     
>
>       Use this when you want to be specific about which plugin or skill ChatGPT
>       should use. See [Skills & Plugins](https://learn.chatgpt.com/docs/skills-and-plugins).
>     
>
>   
>
>
>
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="app">
>
> ### Use Apple Messages from Codex
>
> The Apple Messages plugin is available on all plans in the ChatGPT desktop app
> for macOS. In Codex and ChatGPT Work, it can read and search iMessage, SMS, and
> RCS chats on your Mac and send messages on your behalf through the Messages app.
> It doesn't let you interact with ChatGPT remotely through Messages, and it
> doesn't work in regular ChatGPT chats.
>
> For this release, the Messages plugin is included only in the Apple Silicon
> (arm64) build of the ChatGPT desktop app.
>
> <WorkflowSteps>
>
> 1. Open **Plugins**, find the Apple Messages plugin, and install it.
> 2. Start a new Codex or ChatGPT Work chat and ask it to find, summarize, draft,
>    or send a message.
> 3. Grant the requested macOS permissions before ChatGPT reads Messages.
> 4. Review the message and its recipients before allowing a send.
>
> </WorkflowSteps>
>
> By default, ChatGPT sends messages only after you approve the message and its
> recipients. Choose **Allow once** to approve only that send. If you select
> **Always allow sending to this chat**, ChatGPT can send future messages to that
> Messages chat without another send approval.
>
> Keep per-send approval for chats that may contain untrusted or misleading
> instructions. Persistent approval removes your final chance to review a message
> before ChatGPT sends it as you. Use it only when you accept that risk.
>
> To restore per-send approval, open **Settings** > **Computer use** and select
> **Manage** next to **Messages**. Under **Always allowed to send**, select the
> trash icon next to the chat, then confirm **Remove**. ChatGPT will ask
> before sending to that chat again.
>
> **Known issue:** If your task is set to **Full access** or otherwise disables
> approval prompts, Apple Messages may be unable to show the confirmation needed
> to send. Switch to **Ask for approval** or **Approve for me** and try again.
>
> Apple Messages runs on your Mac. It isn't directly available in ChatGPT on the
> web or mobile, Codex CLI, or the IDE extension.
>
> In managed workspaces, administrators can disable Apple Messages through the
> existing Computer Use control.
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" id="cli">
>
> <a id="plugin-directory-in-codex-cli"></a>
>
> ### Plugin browser in Codex CLI
>
> In Codex CLI, run the following command to open the plugin browser:
>
> ```text
> codex
> /plugins
> ```
>
>
>   
>
> > Illustration: Plugins list in Codex CLI
>
>
>
>
> The CLI plugin browser groups plugins by marketplace. Use the marketplace tabs
> to switch sources, open a plugin to inspect details, install or uninstall
> marketplace entries, and press <kbd>Space</kbd> on an installed plugin to turn it
> on or off.
>
> </ContentModeSwitch>
>
> <a id="api-key-availability"></a>
>
> <ContentModeSwitch group="codex-surface" ids="app,cli">
>
> ### API key availability
>
> If you [sign in to Codex with an OpenAI API
> key](https://learn.chatgpt.com/docs/auth#sign-in-with-an-api-key), you can browse, install, and manage
> supported OpenAI-curated plugins in Codex CLI and Codex in the ChatGPT desktop
> app. Some plugins aren't available with API key authentication because their
> connection flows require unsupported OAuth capabilities. Review plugin usage
> on the [Platform Usage page](https://platform.openai.com/usage).
>
> </ContentModeSwitch>
>
> ### How permissions and data sharing work
>
> <ContentModeSwitch group="codex-surface" id="web">
>
> In ChatGPT on the web, Chat and Work use the workspace permissions and tools
> available to that chat. Connectors still require their own sign-in and access.
>
> </ContentModeSwitch>
>
> <ContentModeSwitch group="codex-surface" ids="app,cli">
>
> When a plugin capability runs through a Codex host, the host's [sandbox and
> approval policy](https://learn.chatgpt.com/docs/agent-approvals-security) applies.
> Connections to external services use that service's own authentication and
> access controls.
>
> </ContentModeSwitch>
>
> - Bundled skills become available when you start a new chat or CLI session
>   after installation.
> - If a plugin includes connectors, the active product may prompt you to install
>   or sign in to those connectors during setup or the first time you use them.
> - If a plugin includes MCP servers, they may require extra setup or
>   authentication before you can use them.
> - When ChatGPT sends data through a bundled connector, that service's terms and privacy
>   policy apply.
>
> ### Remove a plugin
>
> To remove a plugin, open it from a supported plugin browser and select
> **Uninstall plugin** when that action is available. Workspace-installed or
> default plugins may not offer that action; your workspace administrator controls
> them instead.
>
> Uninstalling a plugin removes the plugin bundle from that ChatGPT or Codex
> environment, but bundled connectors stay connected until you manage them in
> ChatGPT.
>
> ## Build your own plugin
>
> If you want to create, test, or distribute your own plugin, see
> [Build plugins](https://developers.openai.com/plugins/build/plugins). That page covers local scaffolding,
> manual marketplace setup, workspace sharing, plugin manifests, and packaging
> guidance.
>
> If your plugin includes server-backed capabilities, see
> [Build an MCP server](https://developers.openai.com/plugins/build/mcp-server).
> MCP tools can work without custom UI or return UI when a visual surface helps
> the workflow.
>
> When your plugin is ready for review, see
> [Submit plugins](https://developers.openai.com/plugins/deploy/submission) for the OpenAI Platform submission
> flow, required permissions, review materials, MCP checks, and test case
> requirements.
>
> ## Plugin guides
>
> - [Record & Replay](https://learn.chatgpt.com/docs/extend/record-and-replay): Show ChatGPT a workflow
>   once and turn it into a reusable skill.
> - [Codex Security plugin](https://learn.chatgpt.com/docs/security/plugin): Scan authorized code,
>   confirm findings, and prepare reviewed fixes.

### Source: Build plugins — full page

> # Build plugins
>
> > For the complete documentation index, see [llms.txt](https://learn.chatgpt.com/llms.txt). Markdown versions of documentation pages are available by appending `.md` to the page URL.
>
> To build or submit a plugin, use the complete
> [builder documentation on developers.openai.com](https://developers.openai.com/plugins).
>
>
>
>   <ButtonLink href="/plugins" color="primary" variant="solid" size="lg">
>     Build and submit a plugin
>   </ButtonLink>
>
>
>
> This page provides a brief introduction. A plugin is an installable package
> that can include skills, an MCP server, or both. An MCP server can also return
> optional UI.
>
> ChatGPT and Codex share one universal plugin directory. Publish a public plugin
> once to make the same listing discoverable from supported surfaces in both
> products. During development, use a local marketplace to test the package
> before submitting it to the universal directory.
>
> Start with a skill when you are still iterating on one personal workflow.
> Build a plugin when you want to share that workflow, package related skills,
> connect to an external service, or distribute a stable capability to a team.
>
> ## Create a plugin with `@plugin-creator`
>
> For the fastest setup, use the built-in `@plugin-creator` skill in ChatGPT Work
> mode or `$plugin-creator` in Codex.
>
>
>   
>
> > Illustration: Plugin creator skill in ChatGPT
>
>
>
>
> Describe the outcome, the skills or MCP server to include, and whether you want
> a local marketplace entry for testing. For example:
>
> ```text
> @plugin-creator Create a plugin named meeting-follow-up.
> Include a skill that turns meeting notes into decisions, owners, and next steps.
> Add it to a personal marketplace so I can test it locally.
> ```
>
> The skill creates the required `.codex-plugin/plugin.json` manifest, organizes
> the plugin folder, and can add the plugin to a local marketplace.
>
>
>   
>
> > Illustration: Invoking the plugin creator skill
>
>
>
>
> After it finishes:
>
> 1. Review `.codex-plugin/plugin.json`.
> 2. Check each bundled skill under `skills/`.
> 3. Refresh ChatGPT or Codex and install the plugin from its local marketplace
>    source.
> 4. Test the plugin in a new conversation with representative requests.
>
> If the plugin includes an MCP server, first build and test that server, then
> give `@plugin-creator` the registered connection details. Follow the complete
> [MCP server workflow](https://developers.openai.com/plugins/build/mcp-server)
> for tools, authentication, deployment, and testing.
>
> ## Create a skills-only plugin manually
>
> A minimal plugin contains a manifest and at least one skill:
>
> ```text
> meeting-follow-up/
> ├── .codex-plugin/
> │   └── plugin.json
> └── skills/
>     └── meeting-follow-up/
>         └── SKILL.md
> ```
>
> Create `.codex-plugin/plugin.json`:
>
> ```json
> {
>   "name": "meeting-follow-up",
>   "version": "1.0.0",
>   "description": "Turn meeting notes into decisions and next steps",
>   "skills": "./skills/"
> }
> ```
>
> Then add `skills/meeting-follow-up/SKILL.md`:
>
> ```md
> ---
> name: meeting-follow-up
> description: Extract decisions, owners, and next steps from meeting notes.
> ---
>
> Review the meeting notes. Return:
>
> 1. Decisions
> 2. Action items with owners
> 3. Open questions
> ```
>
> Use a stable plugin name in kebab case. Keep the skill description specific
> enough for ChatGPT and Codex to recognize when the workflow applies.
>
> Use `@plugin-creator` to add the folder to a local marketplace, then install and
> test it before sharing it.
>
> ## Continue with the builder documentation
>
> For complete builder documentation, use the
> [Plugins documentation](https://developers.openai.com/plugins/). It covers:
>
> - [Plugin architecture](https://developers.openai.com/plugins/concepts/plugins)
> - [Building skills](https://developers.openai.com/plugins/build/skills)
> - [Building an MCP server](https://developers.openai.com/plugins/build/mcp-server)
> - [Adding optional UI](https://developers.openai.com/plugins/build/chatgpt-ui)
> - [Packaging a plugin](https://developers.openai.com/plugins/build/plugins)
> - [Testing a plugin](https://developers.openai.com/plugins/deploy/connect-chatgpt)
> - [Submitting and publishing](https://developers.openai.com/plugins/deploy/submission)
>
> To browse, install, enable, or remove plugins, see [Use
> plugins](https://learn.chatgpt.com/docs/plugins).
