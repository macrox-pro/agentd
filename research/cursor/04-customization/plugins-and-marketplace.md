---
primary_sources:
  - id: T1-PLUGINS
    title: "Full page"
    url: "https://cursor.com/docs/plugins.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Plugins and marketplace

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Plugins and marketplace

> # Plugins
>
> Plugins package rules, skills, agents, commands, MCP servers, and hooks into distributable bundles.
>
> Cursor supports the [Agent Plugins](https://agent-plugins.org) open standard alongside its own plugin format. Install and manage them from the [Customize](https://cursor.com/docs/customize-cursor.md) page or browse official plugins in the [Cursor Marketplace](/marketplace). For community plugins and MCP servers, browse [cursor.directory](https://cursor.directory). You can also [build your own](https://cursor.com/docs/plugins.md#creating-plugins) to share with other developers.
>
> ## What plugins contain
>
> A plugin can bundle any combination of these components:
>
> | Component       | Available in   | Description                                                |
> | :-------------- | :------------- | :--------------------------------------------------------- |
> | **Rules**       | Cursor Plugins | Persistent AI guidance and coding standards (`.mdc` files) |
> | **Skills**      | Both formats   | Specialized agent capabilities for complex tasks           |
> | **Agents**      | Cursor Plugins | Custom agent configurations and prompts                    |
> | **Commands**    | Cursor Plugins | Agent-executable command files                             |
> | **MCP Servers** | Both formats   | Model Context Protocol integrations                        |
> | **Hooks**       | Cursor Plugins | Automation scripts triggered by events                     |
>
> ## The Agent Plugins standard
>
> Plugins bundle reusable components an agent can use. [Agent Plugins](https://agent-plugins.org)
> is the open standard for packaging portable skills and MCP servers, much like
> [Agent Skills](https://cursor.com/docs/skills.md) defines a standard for individual skills. Cursor
> supports Agent Plugins alongside Cursor Plugins.
>
> - **Agent Plugins**: spec-conformant plugins with a `plugin.json` manifest at the plugin root, packaging skills and MCP servers
> - **Cursor Plugins**: plugins with a `.cursor-plugin/plugin.json` manifest, which add rules, agents, commands, hooks, and [variables](https://cursor.com/docs/reference/plugins.md#variables)
>
> A plugin that follows the Agent Plugins specification loads in Cursor without changes. Cursor Plugins continue to develop in parallel with the standard, so Cursor-specific components and marketplace features keep working as they do today.
>
> Learn more at [agent-plugins.org](https://agent-plugins.org) or read the [specification on GitHub](https://github.com/agentplugins/agent-plugins-spec).
>
> ## Cursor Plugin canvases
>
> Plugins now ship with prebuilt **canvases**: shared setup templates your team can open and reuse.
>
> - **Hex Canvas** — Build data visualizations. At Cursor, we use the Hex Canvas to explore and share analytics.
> - **Atlassian Canvas** — See a realtime view of your issues, projects, and documents from Jira and Confluence.
>
> Open a canvas from an installed plugin in Customize to get a guided starting point instead of configuring everything from scratch.
>
> ## The marketplace
>
> The [Cursor Marketplace](/marketplace) is where you discover and install official Cursor Plugins. Plugins are distributed as Git repositories and submitted through the Cursor team.
>
> Every plugin is [manually reviewed](https://cursor.com/help/security-and-privacy/marketplace-security.md) before it's listed. Browse official plugins at [cursor.com/marketplace](https://cursor.com/marketplace) or search by keyword in **Customize**. For community plugins and MCP servers, browse [cursor.directory](https://cursor.directory).
>
> ## Team marketplaces
>
> Team marketplaces are available on Teams and Enterprise plans.
> They can distribute Agent Plugins and Cursor Plugins through the same
> marketplace.
>
> - Teams plan: up to 1 team marketplace
> - Enterprise plan: unlimited team marketplaces
>
> [Contact sales](https://cursor.com/contact-sales?source=docs-plugins) for unlimited team marketplaces and Enterprise admin controls.
>
> Open **Dashboard -> Plugins** to manage Team Marketplaces.
>
> On Enterprise plans, only admins can add team marketplaces from **Dashboard
> -> Plugins**.
>
> ### Default team marketplace
>
> The **Default** team marketplace connects shared plugins and MCP servers across Cursor. Admins can add Team MCP servers that are already available to Cloud Agents, then make the same servers available for teammates to install and configure in the Agent Window, IDE, and CLI.
>
> Adding a Team MCP server to the Default marketplace does not install or enable it for every developer. Admins still control marketplace access and plugin installation modes. Each developer may also need to authenticate with the MCP provider.
>
> ### Migrate existing Team MCPs
>
> Admins can link standalone Team MCP servers to the Default marketplace:
>
> 1. Open **Dashboard -> Integrations & MCP**.
> 2. Find **Team MCP Servers**.
> 3. Select **Add to Team Marketplace** in the migration prompt.
> 4. Open **Dashboard -> Plugins** to review the Default marketplace, its access, and plugin installation modes.
>
> Cursor creates the Default marketplace if needed and links the existing MCP servers to it. The servers remain available to Cloud Agents while teammates gain the option to install and configure them locally.
>
> Removing a linked MCP plugin from the marketplace or deleting the marketplace
> can delete the Team MCP server. This removes it for local users and Cloud
> Agents. Review the confirmation message before continuing.
>
> ### Marketplace access
>
> Team marketplaces are available to everyone in their team by default. Under **Marketplace Settings -> Marketplace Access**, admins can restrict a marketplace to selected [Organization Groups](https://cursor.com/docs/enterprise/organization-groups.md). Only members of the marketplace's team who belong to a selected group receive access. Team admins retain access.
>
> ### How does SCIM work?
>
> Organization Groups can sync membership from your identity provider through [SCIM](https://cursor.com/docs/account/teams/scim.md). Manage membership in your identity provider, and Cursor syncs those updates to the Organization Group.
>
> Existing marketplaces that use team-level SCIM directory groups keep that configuration. Cursor does not migrate those assignments automatically. Organizations without Organization Groups continue to use SCIM directory groups.
>
> ### Plugin installation modes
>
> After setting marketplace access, choose how each plugin is distributed to that audience:
>
> - **Default Off**: Developers can find the plugin and choose whether to install it.
> - **Default On**: The plugin is installed by default, but developers can opt out.
> - **Required**: The plugin is always installed and cannot be uninstalled.
>
> ## Add a team marketplace
>
> Use this flow to import a GitHub repository as a team marketplace:
>
> 1. Go to **Dashboard -> Plugins**.
> 2. In **Team Marketplaces**, click **Add Marketplace**.
> 3. Follow the instructions to create a marketplace from scratch, or use "Import from Repo" if importing from GitHub.
> 4. Add and review plugins using "Add to Marketplace".
> 5. Under **Marketplace Settings**, set **Marketplace Access**, optionally enable Auto Refresh, then save.
>
> Example repository to try:
>
> - [fieldsphere/cursor-team-marketplace-template](https://github.com/fieldsphere/cursor-team-marketplace-template)
>
> ## Keep plugins up to date
>
> When importing from GitHub, plugins are indexed when you first import the repository. You can refresh plugins in two ways:
>
> - **Automatically**: Turn on **Enable Auto Refresh** to update plugins automatically whenever changes are pushed to the branch the marketplace tracks. This requires the [Cursor GitHub App](https://cursor.com/docs/integrations/github.md) installed on the repository. Cursor re-indexes a marketplace at most once every 10 minutes, batching rapid pushes to the latest commit.
> - **Manually**: Click "Refresh" to manually update.
>
> For marketplaces created with "Import from Repo", Auto Refresh re-reads the full manifest on each push, so new plugins added to the repository are picked up automatically.
>
> For marketplaces where plugins were added individually, Auto Refresh only updates existing plugins. Re-import the repository URL to pick up newly added plugins.
>
> ## Where developers find team marketplaces
>
> Developers can find team marketplaces in Customize.
>
> - Open **Customize** in the sidebar
> - Look for plugins from your team marketplace.
> - Install Default Off plugins directly from that panel.
> - Default On plugins are installed automatically, but developers can opt out.
> - Required plugins are installed automatically and cannot be uninstalled.
> - Install and configure marketplace MCP servers for use in the Agent Window, IDE, and CLI.
>
> ## Installing plugins
>
> Install Agent Plugins and Cursor Plugins from a marketplace. Cursor detects the
> format from the plugin manifest, so the installation flow is the same for both:
>
> 1. Open **Customize** in the sidebar.
> 2. Find the plugin you want to use.
> 3. Select **Install** and choose a project or user scope.
>
> An Agent Plugin has a `plugin.json` manifest at its root. A Cursor Plugin has a
> `.cursor-plugin/plugin.json` manifest. Team marketplaces can distribute either
> format using the same access and installation modes described above.
>
> You can also load either format directly from `~/.cursor/plugins/local` while
> developing a plugin. See [Test plugins locally](https://cursor.com/docs/plugins.md#test-plugins-locally).
>
> ### MCP Apps deeplinks
>
> Share MCP server configurations using install links:
>
> ```text
> cursor://anysphere.cursor-deeplink/mcp/install?name=$NAME&config=$BASE64_ENCODED_CONFIG
> ```
>
> See [MCP install links](https://cursor.com/docs/mcp/install-links.md) for details on generating these links.
>
> ## Managing installed plugins
>
> Open **Customize** in the sidebar to manage installed Agent Plugins, Cursor
> Plugins, MCP servers, rules, and skills from one page. Filter by user, workspace,
> or team scope to see what is installed.
>
> ### MCP servers
>
> Toggle personal and team-distributed MCP servers on or off from Customize:
>
> 1. Open **Customize** in the sidebar
> 2. Find the MCP server you want to change
> 3. Use the toggle to enable or disable it
>
> Disabled servers won't load or appear in chat.
>
> ### Rules and skills
>
> Manage rules and skills from Customize. Toggle individual rules between **Always**, **Agent Decides**, and **Manual** modes. Skills appear in the **Agent Decides** section and can be invoked manually with `/skill-name` in chat.
>
> ## Using the workspaceOpen hook
>
> A `workspaceOpen` hook can return plugin paths to load on workspace open, which is useful when the set of plugins depends on the workspace itself.
>
> ### Hooks reference
>
> Register plugin paths from a `workspaceOpen` hook script
>
> ## Creating plugins
>
> A plugin is a directory with a manifest and its components. Choose Agent Plugins
> when you want to package portable skills and MCP servers. Choose Cursor Plugins
> when you also need Cursor-specific components such as rules, agents, commands,
> hooks, or variables.
>
> ### Agent Plugin
>
> ```text
> my-plugin/
> ├── plugin.json
> ├── skills/
> │   └── code-reviewer/
> │       └── SKILL.md
> └── mcp.json
> ```
>
> Agent Plugins require a root `plugin.json` with the standard's schema identifier:
>
> ```json
> {
>   "$schema": "https://agent-plugins.org/schemas/1.0.0/plugin.schema.json",
>   "name": "my-plugin",
>   "description": "Portable code review tools",
>   "version": "1.0.0",
>   "author": { "name": "Your Name" }
> }
> ```
>
> ### Cursor Plugin
>
> ```text
> my-plugin/
> ├── .cursor-plugin/
> │   └── plugin.json
> ├── rules/
> │   └── coding-standards.mdc
> ├── skills/
> │   └── code-reviewer/
> │       └── SKILL.md
> └── mcp.json
> ```
>
> Cursor Plugin manifests only require a `name`. Components are discovered from
> their default directories, or you can specify custom paths in the manifest.
>
> ```json
> {
>   "name": "my-plugin",
>   "description": "Custom development tools",
>   "version": "1.0.0",
>   "author": { "name": "Your Name" }
> }
> ```
>
> Start from the [Cursor Plugin template repository](https://github.com/cursor/plugin-template),
> or read the [Agent Plugins authoring guide](https://agent-plugins.org/plugin-authors)
> to create an Agent Plugin.
>
> ### Test plugins locally
>
> Before you publish, load either plugin format from `~/.cursor/plugins/local`:
>
> 1. Create a folder for your plugin:
>    `~/.cursor/plugins/local/my-plugin`
> 2. Copy your plugin files into that folder. Include either a root `plugin.json`
>    for an Agent Plugin or `.cursor-plugin/plugin.json` for a Cursor Plugin.
> 3. Restart Cursor, or run **Developer: Reload Window**.
> 4. Verify your plugin components load in Cursor, such as rules, skills, or MCP servers.
>
> For faster iteration, symlink your plugin repository:
>
> ```bash
> ln -s /path/to/my-plugin ~/.cursor/plugins/local/my-plugin
> ```
>
> When your plugin is ready, submit it for review at [cursor.com/marketplace/publish](https://cursor.com/marketplace/publish).
> Cursor Plugins can use `.cursor-plugin/marketplace.json` for multi-plugin
> repositories.
>
> See the [Plugins reference](https://cursor.com/docs/reference/plugins.md) for the full manifest schema, component formats, and submission checklist.
>
> ### Team and Enterprise marketplaces
>
> Upgrade for private team marketplaces and organization-wide plugin distribution.
>
> ## FAQ
>
> ### Are marketplace plugins reviewed for security?
>
> Yes. Every plugin is manually reviewed before it's listed. All plugins must be open source, and we review each update before publishing. See [Marketplace security](https://cursor.com/help/security-and-privacy/marketplace-security.md) for details on vetting, update reviews, and how to report issues.
>
> ### How do I create a plugin?
>
> For a portable Agent Plugin, add a root `plugin.json` manifest and package
> skills or MCP servers. For a Cursor Plugin, add a
> `.cursor-plugin/plugin.json` manifest and any Cursor components you need.
> See the [Plugins reference](https://cursor.com/docs/reference/plugins.md) for examples of both
> formats.
>
> ### How do Cursor Plugins relate to the Agent Plugins standard?
>
> [Agent Plugins](https://agent-plugins.org) is an open, vendor-neutral specification for packaging skills and MCP servers into portable plugins. Cursor supports the standard, so spec-conformant plugins load in Cursor without changes. Cursor Plugins are developed in parallel and add Cursor-specific components like rules, agents, commands, hooks, and variables.
>
> ## Related
>
> - [Plugins help](https://cursor.com/help/customization/plugins.md)
>
>

---
primary_sources:
  - id: T1-PLUGINS-REF
    title: "Full page"
    url: "https://cursor.com/docs/reference/plugins.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Plugin manifest reference

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Plugin manifest reference

> # Plugins reference
>
> Reference documentation for building, structuring, and submitting Cursor plugins. Plugins package rules, skills, agents, commands, MCP servers, and hooks into distributable bundles that work in the Cursor IDE.
>
> If you're starting from scratch, use the [plugin template repository](https://github.com/cursor/plugin-template).
>
> ## Supported plugin formats
>
> Cursor loads plugins in two formats, identified by their manifest location:
>
> | Format                                                     | Manifest location                | Components                                                     |
> | :--------------------------------------------------------- | :------------------------------- | :------------------------------------------------------------- |
> | [Agent Plugins](https://agent-plugins.org) (open standard) | `plugin.json` at the plugin root | Skills, MCP servers                                            |
> | Cursor Plugins                                             | `.cursor-plugin/plugin.json`     | Skills, MCP servers, rules, agents, commands, hooks, variables |
>
> A plugin that conforms to the [Agent Plugins specification](https://github.com/agentplugins/agent-plugins-spec) loads in Cursor without changes. The rest of this reference documents the Cursor plugin format, which is developed in parallel with the standard and supports the full set of Cursor components.
>
> ## Plugin structure
>
> A plugin is a directory with a manifest file and your plugin assets:
>
> ### Agent Plugin
>
> ```text
> my-plugin/
> ├── plugin.json            # Required: Agent Plugins manifest
> ├── skills/                # Agent Skills
> │   └── code-reviewer/
> │       └── SKILL.md
> └── mcp.json               # MCP server definitions
> ```
>
> The Agent Plugins standard defines portable skills and MCP servers. See the
> [Agent Plugins authoring guide](https://agent-plugins.org/plugin-authors) for
> the full package and schema reference.
>
> ### Cursor Plugin
>
> ```text
> my-plugin/
> ├── .cursor-plugin/
> │   └── plugin.json        # Required: Cursor Plugin manifest
> ├── rules/                 # Cursor rules (.mdc files)
> │   ├── coding-standards.mdc
> │   └── review-checklist.mdc
> ├── skills/                # Agent skills
> │   └── code-reviewer/
> │       └── SKILL.md
> ├── agents/                # Custom agent configurations
> │   └── security-reviewer.md
> ├── commands/              # Agent-executable commands
> │   └── deploy.md
> ├── hooks/                 # Hook definitions
> │   └── hooks.json
> ├── mcp.json               # MCP server definitions
> ├── assets/                # Logos and static assets
> │   └── logo.svg
> ├── scripts/               # Hook and utility scripts
> │   └── format-code.py
> └── README.md
> ```
>
> ## Cursor Plugin manifest
>
> Every Cursor Plugin requires a `.cursor-plugin/plugin.json` manifest file. The
> sections below document Cursor Plugin fields, components, and marketplace
> features. For a root Agent Plugins manifest, use the
> [standard's manifest reference](https://agent-plugins.org/plugin-authors/manifest).
>
> ### Required fields
>
> | Field  | Type   | Description                                                                                                                                                              |
> | :----- | :----- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `name` | string | Plugin identifier. Lowercase, kebab-case (alphanumerics, hyphens, and periods). Must start and end with an alphanumeric character. Examples: `my-plugin`, `prompts.chat` |
>
> ### Optional fields
>
> | Field         | Type                     | Description                                                                                                                                                                                                                                                                                         |
> | :------------ | :----------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `description` | string                   | Brief plugin description                                                                                                                                                                                                                                                                            |
> | `version`     | string                   | Semantic version (e.g., `1.0.0`)                                                                                                                                                                                                                                                                    |
> | `author`      | object                   | Author info: `name` (required), `email` (optional)                                                                                                                                                                                                                                                  |
> | `homepage`    | string                   | URL to plugin homepage                                                                                                                                                                                                                                                                              |
> | `repository`  | string                   | URL to plugin repository                                                                                                                                                                                                                                                                            |
> | `license`     | string                   | License identifier (e.g., `MIT`)                                                                                                                                                                                                                                                                    |
> | `keywords`    | array                    | Tags for discovery and categorization                                                                                                                                                                                                                                                               |
> | `logo`        | string                   | Relative path to a logo file in the repo (e.g., `assets/logo.svg`), or an absolute URL. Relative paths resolve to `raw.githubusercontent.com` URLs. Preferred: commit the logo to your repo and use a relative path.                                                                                |
> | `rules`       | string or array          | Path(s) to rule files or directories                                                                                                                                                                                                                                                                |
> | `agents`      | string or array          | Path(s) to agent files or directories                                                                                                                                                                                                                                                               |
> | `skills`      | string or array          | Path(s) to skill directories                                                                                                                                                                                                                                                                        |
> | `commands`    | string or array          | Path(s) to command files or directories                                                                                                                                                                                                                                                             |
> | `hooks`       | string or object         | Path to hooks config file, or inline hook config                                                                                                                                                                                                                                                    |
> | `mcpServers`  | string, object, or array | Path to MCP config file, inline MCP server config, or an array of either. Overrides default `mcp.json` discovery.                                                                                                                                                                                   |
> | `variables`   | object                   | JSON Schema that declares variable **names** (tokens, connection strings). The plugin does not store secret values; users set them in the dashboard (**Plugins** → **Configure**). Substituted into `${VAR}` placeholders. See [Variables](https://cursor.com/docs/reference/plugins.md#variables). |
>
> ### Example manifest
>
> ```json
> {
>   "name": "enterprise-plugin",
>   "version": "1.2.0",
>   "description": "Enterprise development tools with security scanning and compliance checks",
>   "author": {
>     "name": "ACME DevTools",
>     "email": "devtools@acme.com"
>   },
>   "keywords": ["enterprise", "security", "compliance"],
>   "logo": "assets/logo.svg"
> }
> ```
>
> ## Variables
>
> Use `variables` to declare the **names** (and types/descriptions) of user-specified configuration — for example an API token for an HTTP MCP server. The plugin only defines the schema; it does not include the secret values themselves.
>
> Team admins set the actual values in the dashboard under **Plugins** (at install time, or later via **Configure** on the plugin).
>
> Do not put secret values in the plugin repo. In `mcp.json` and other plugin config, include only `${VAR}` placeholders that match property names in the schema.
>
> ```json title=".cursor-plugin/plugin.json"
> {
>   "name": "example-plugin",
>   "variables": {
>     "type": "object",
>     "properties": {
>       "API_TOKEN": {
>         "type": "string",
>         "title": "API token",
>         "description": "Bearer token for the example HTTP MCP"
>       }
>     },
>     "required": ["API_TOKEN"]
>   }
> }
> ```
>
> ```json title="mcp.json"
> {
>   "mcpServers": {
>     "example-api": {
>       "url": "https://mcp.example.com/mcp",
>       "headers": {
>         "Authorization": "Bearer ${API_TOKEN}"
>       }
>     }
>   }
> }
> ```
>
> The top level must be `{ "type": "object", "properties": { ... } }`. Only a fixed set of JSON Schema keywords is accepted (`type`, `title`, `description`, `default`, `enum`, `const`, `properties`, `required`, `items`, and common length/numeric constraints).
>
> ## Cursor Plugin component discovery
>
> When the manifest does not specify explicit paths for a component type, the parser uses **automatic folder-based discovery**:
>
> | Component   | Default location          | How it's discovered                                                                        |
> | :---------- | :------------------------ | :----------------------------------------------------------------------------------------- |
> | Skills      | `skills/`                 | Each subdirectory containing a `SKILL.md` file                                             |
> | Rules       | `rules/`                  | All `.md`, `.mdc`, or `.markdown` files                                                    |
> | Agents      | `agents/`                 | All `.md`, `.mdc`, or `.markdown` files                                                    |
> | Commands    | `commands/`               | All `.md`, `.mdc`, `.markdown`, or `.txt` files                                            |
> | Hooks       | `hooks/hooks.json`        | Parsed for hook event names                                                                |
> | MCP Servers | `mcp.json`                | Parsed for server entries                                                                  |
> | Root Skill  | `SKILL.md` at plugin root | Treated as a single-skill plugin (only if no `skills/` dir and no manifest `skills` field) |
>
> If a manifest field **is** specified (e.g., `"skills": "./my-skills/"`), it **replaces** folder discovery for that component. The default folder is not also scanned.
>
> ## Rules format
>
> Rules are `.mdc` files providing persistent guidance to the AI. Place them in the `rules/` directory.
>
> Rules require YAML frontmatter with metadata:
>
> ```markdown title="rules/prefer-const.mdc"
