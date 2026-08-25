---
primary_sources:
  - id: T2-PLUGIN-MKT
    title: "Plugin marketplaces"
    url: "https://code.claude.com/docs/en/plugin-marketplaces.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Plugin marketplaces

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Plugin marketplaces
>
> # Create and distribute a plugin marketplace
>
> > Build and host plugin marketplaces to distribute Claude Code extensions across teams and communities.
>
> A **plugin marketplace** is a catalog that lets you distribute plugins to others. Marketplaces provide centralized discovery, version tracking, automatic updates, and support for multiple source types, including git repositories and local paths. This guide shows you how to create your own marketplace to share plugins with your team or community.
>
> Looking to install plugins from an existing marketplace? See [Discover and install prebuilt plugins](/docs/en/discover-plugins).
>
> ## Overview
>
> Creating and distributing a marketplace involves:
>
> 1. **Create plugins**: build one or more plugins with skills, agents, hooks, MCP servers, or LSP servers. This guide assumes you already have plugins to distribute; see [Create plugins](/docs/en/plugins) for details on how to create them.
> 2. **Create the marketplace file**: define a `marketplace.json` that lists your plugins and where to find them. See [Create the marketplace file](#create-the-marketplace-file).
> 3. **Host the marketplace**: push to GitHub, GitLab, or another git host. See [Host and distribute marketplaces](#host-and-distribute-marketplaces).
> 4. **Share with users**: users add your marketplace with `/plugin marketplace add` and install individual plugins. See [Discover and install plugins](/docs/en/discover-plugins).
>
> Once your marketplace is live, you can update it by pushing changes to your repository. Users refresh their local copy with `/plugin marketplace update`.
>
> ## Walkthrough: create a local marketplace
>
> This example creates a marketplace with one plugin: a `quality-review` skill for code reviews. You'll create the directory structure, add a skill, create the plugin manifest and marketplace catalog, then install and test it.
>
>   **Create the directory structure**: ```bash
>     mkdir -p my-marketplace/.claude-plugin
>     mkdir -p my-marketplace/plugins/quality-review-plugin/.claude-plugin
>     mkdir -p my-marketplace/plugins/quality-review-plugin/skills/quality-review
>     ```
>
>   **Create the skill**: Create a `SKILL.md` file that defines what the `quality-review` skill does.
>
>     ```markdown my-marketplace/plugins/quality-review-plugin/skills/quality-review/SKILL.md
>     ---
>     description: Review code for bugs, security, and performance
>     ---
>
>     Review the code I've selected or the recent changes for:
>     - Potential bugs or edge cases
>     - Security concerns
>     - Performance issues
>     - Readability improvements
>
>     Be concise and actionable.
>     ```
>
>   **Create the plugin manifest**: Create a `plugin.json` file that describes the plugin. The manifest goes in the `.claude-plugin/` directory.
>
>     ```json my-marketplace/plugins/quality-review-plugin/.claude-plugin/plugin.json
>     {
>       "name": "quality-review-plugin",
>       "description": "Adds a quality-review skill for quick code reviews",
>       "version": "1.0.0",
>       "author": {
>         "name": "Your Name"
>       }
>     }
>     ```
>
>
>
> > [Note] Setting `version` means users only receive updates when you change this field, so bump it on every release. A plugin with a [`command` source](#command-sources) isn't pinned by this field. If you omit `version`, the version comes from the next source in [version management](/docs/en/plugins-reference#version-management).
>
>   **Create the marketplace file**: Create the marketplace catalog that lists your plugin.
>
>     ```json my-marketplace/.claude-plugin/marketplace.json
>     {
>       "name": "my-plugins",
>       "owner": {
>         "name": "Your Name"
>       },
>       "plugins": [
>         {
>           "name": "quality-review-plugin",
>           "source": "./plugins/quality-review-plugin",
>           "description": "Adds a quality-review skill for quick code reviews"
>         }
>       ]
>     }
>     ```
>
>   **Add and install**: From the directory that contains `my-marketplace`, start Claude Code and run the following commands. The install command opens a plugin details view where you select an installation scope to confirm the install. Check the install summary: if it reports `Run /reload-plugins to activate.`, run that command.
>
>     ```shell
>     /plugin marketplace add ./my-marketplace
>     /plugin install quality-review-plugin@my-plugins
>     ```
>
>   **Try it out**: Select some code in your editor and run your new skill. Plugin skills are namespaced with the plugin name.
>
>     ```shell
>     /quality-review-plugin:quality-review
>     ```
>
> To learn more about what plugins can do, including hooks, agents, MCP servers, and LSP servers, see [Plugins](/docs/en/plugins).
>
>   **How plugins are installed**: when users install a plugin, Claude Code copies the plugin directory to a cache location, except for a [`command` source in link mode](#copy-mode-and-link-mode), which is used in place. Copied plugins can't reference files outside their directory using paths like `../shared-utils`, because those files won't be copied.
>
>   If you need to share files across plugins, use symlinks. See [Plugin caching and file resolution](/docs/en/plugins-reference#plugin-caching-and-file-resolution) for details.
>
> ## Create the marketplace file
>
> Create `.claude-plugin/marketplace.json` in your repository root. This file defines your marketplace's name, owner information, and a list of plugins with their sources.
>
> Each plugin entry needs at minimum a `name` and a `source` that tells Claude Code where to fetch it from. See the [full schema](#marketplace-schema) below for all available fields.
>
> ```json
> {
>   "name": "company-tools",
>   "owner": {
>     "name": "DevTools Team",
>     "email": "devtools@example.com"
>   },
>   "plugins": [
>     {
>       "name": "code-formatter",
>       "source": "./plugins/formatter",
>       "description": "Automatic code formatting on save",
>       "version": "2.1.0",
>       "author": {
>         "name": "DevTools Team"
>       }
>     },
>     {
>       "name": "deployment-tools",
>       "source": {
>         "source": "github",
>         "repo": "company/deploy-plugin"
>       },
>       "description": "Deployment automation tools"
>     }
>   ]
> }
> ```
>
> ## Marketplace schema
>
> ### Required fields
>
> | Field     | Type   | Description                                                                                                                                                                                                                                                                                                                                                                                                                           | Example        |
> | :-------- | :----- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | :------------- |
> | `name`    | string | Marketplace identifier (kebab-case, no spaces). This is public-facing: users see it when installing plugins (for example, `/plugin install my-tool@your-marketplace`). Each user can register only one marketplace per name: adding a second marketplace with the same name replaces the first. To publish multiple plugins under one marketplace name, list them all in a [single `marketplace.json`](#create-the-marketplace-file). | `"acme-tools"` |
> | `owner`   | object | Marketplace maintainer information ([see fields below](#owner-fields))                                                                                                                                                                                                                                                                                                                                                                |                |
> | `plugins` | array  | List of available plugins                                                                                                                                                                                                                                                                                                                                                                                                             | See below      |
>
>   **Reserved names**: the following marketplace names are reserved for official Anthropic use and can't be used by third-party marketplaces: `claude-code-marketplace`, `claude-code-plugins`, `claude-plugins-official`, `claude-plugins-community`, `claude-community`, `anthropic-marketplace`, `anthropic-plugins`, `agent-skills`, `anthropic-agent-skills`, `knowledge-work-plugins`, `life-sciences`, `claude-for-legal`, `claude-for-financial-services`, `financial-services-plugins`, `first-party-plugins`, `healthcare`. Names that impersonate official marketplaces, such as `official-claude-plugins` or `anthropic-plugins-v2`, are also blocked. Reserving these names prevents a third-party marketplace from presenting itself as an Anthropic-published source.
>
>   Claude Code re-checks reserved names every time it loads a marketplace, not only when you add one. A marketplace that was registered under one of these names before the name became reserved stops loading and reports that it is [registered from an untrusted source](/docs/en/errors#marketplace-is-registered-from-an-untrusted-source). Remove that marketplace and re-add it from the official Anthropic source. A third-party marketplace affected by a newly reserved name loads again as soon as you re-add it under a different name. Before v2.1.205, `first-party-plugins` and `healthcare` weren't reserved, and a marketplace already registered under a reserved name kept loading.
>
> ### Owner fields
>
> | Field   | Type   | Required | Description                                  |
> | :------ | :----- | :------- | :------------------------------------------- |
> | `name`  | string | Yes      | Name of the maintainer or team               |
> | `email` | string | No       | Contact email for the maintainer             |
> | `url`   | string | No       | Website, GitHub profile, or organization URL |
>
> ### Optional fields
>
> | Field                                 | Type   | Description                                                                                                                                                                                                                                                                                  |
> | :------------------------------------ | :----- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `$schema`                             | string | JSON Schema URL for editor autocomplete and validation. Claude Code ignores this field at load time.                                                                                                                                                                                         |
> | `description`                         | string | Brief marketplace description                                                                                                                                                                                                                                                                |
> | `version`                             | string | Marketplace manifest version                                                                                                                                                                                                                                                                 |
> | `metadata.pluginRoot`                 | string | Directory that Claude Code resolves bare plugin source names under. See [Relative paths](#relative-paths). Requires Claude Code v2.1.239 or later.                                                                                                                                           |
> | `allowCrossMarketplaceDependenciesOn` | array  | Other marketplaces that plugins in this marketplace may depend on. Dependencies from a marketplace not listed here are blocked at install. See [Depend on a plugin from another marketplace](/docs/en/plugin-dependencies#depend-on-a-plugin-from-another-marketplace).                           |
> | `renames`                             | object | Map from a former plugin `name` to its current name, or to `null` if the plugin was removed. Lets existing users migrate automatically when you rename or remove an entry in `plugins`. See [Rename or remove a plugin](#rename-or-remove-a-plugin). Requires Claude Code v2.1.193 or later. |
>
> `description` and `version` are also accepted under `metadata` for backward compatibility.
>
> ## Plugin entries
>
> Each plugin entry in the `plugins` array describes a plugin and where to find it. You can include any field from the [plugin manifest schema](/docs/en/plugins-reference#plugin-manifest-schema), such as `description`, `version`, `author`, `commands`, and `hooks`, plus these marketplace-specific fields: `source`, `category`, `tags`, `strict`, `relevance`, `headers`, and `headersHelper`.
>
> ### Required fields
>
> | Field    | Type           | Description                                                                                                                                            |
> | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `name`   | string         | Plugin identifier (kebab-case, no spaces). This is public-facing: users see it when installing (for example, `/plugin install my-plugin@marketplace`). |
> | `source` | string\|object | Where to fetch the plugin from (see [Plugin sources](#plugin-sources) below)                                                                           |
>
> ### Optional plugin fields
>
> **Standard metadata fields:**
>
> | Field            | Type    | Description                                                                                                                                                                                                                                                                                                                                                  |
> | :--------------- | :------ | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `displayName`    | string  | Human-readable name shown in UI surfaces. Falls back to `name` when omitted. May contain spaces and any casing. Not used for namespacing or lookup.                                                                                                                                                                                                          |
> | `description`    | string  | Brief plugin description                                                                                                                                                                                                                                                                                                                                     |
> | `version`        | string  | Plugin version. If set (here or in `plugin.json`), the plugin is pinned to this string and users only receive updates when it changes. A plugin with a [`command` source](#command-sources) isn't pinned by either field. If set in neither place, the version comes from the next source in [version management](/docs/en/plugins-reference#version-management). |
> | `author`         | object  | Plugin author information (`name` required; `email` and `url` optional)                                                                                                                                                                                                                                                                                      |
> | `homepage`       | string  | Plugin homepage or documentation URL                                                                                                                                                                                                                                                                                                                         |
> | `repository`     | string  | Source code repository URL                                                                                                                                                                                                                                                                                                                                   |
> | `license`        | string  | SPDX license identifier (for example, MIT, Apache-2.0)                                                                                                                                                                                                                                                                                                       |
> | `keywords`       | array   | Tags for plugin discovery and categorization                                                                                                                                                                                                                                                                                                                 |
> | `metadata`       | object  | Free-form object for your own fields, such as entitlement or catalog data. Claude Code doesn't read it. Before v2.1.222, `claude plugin validate` reported the key as an unrecognized field.                                                                                                                                                                 |
> | `category`       | string  | Plugin category for organization                                                                                                                                                                                                                                                                                                                             |
> | `tags`           | array   | Tags for searchability                                                                                                                                                                                                                                                                                                                                       |
> | `strict`         | boolean | Controls whether `plugin.json` is the authority for component definitions (default: true). See [Strict mode](#strict-mode) below.                                                                                                                                                                                                                            |
> | `relevance`      | object  | Signals that tell Claude Code when to suggest this plugin to users. Takes effect only for marketplaces an administrator allowlists in managed settings. See [Recommend plugins for your org](/docs/en/plugin-relevance). Requires Claude Code v2.1.152 or later.                                                                                                  |
> | `defaultEnabled` | boolean | Whether the plugin is enabled after install (default: true). Set to `false` to install the plugin disabled until the user opts in. Takes precedence over the same field in the plugin's `plugin.json`. See [Default enablement](/docs/en/plugins-reference#default-enablement). Requires Claude Code v2.1.154 or later.                                           |
>
> **Component configuration fields:**
>
> | Field        | Type           | Description                                                    |
> | :----------- | :------------- | :------------------------------------------------------------- |
> | `skills`     | string\|array  | Custom paths to skill directories containing `<name>/SKILL.md` |
> | `commands`   | string\|array  | Custom paths to flat `.md` skill files or directories          |
> | `agents`     | string\|array  | Custom paths to agent files                                    |
> | `hooks`      | string\|object | Custom hooks configuration or path to hooks file               |
> | `mcpServers` | string\|object | MCP server configurations or path to MCP config                |
> | `lspServers` | string\|object | LSP server configurations or path to LSP config                |
>
> **Archive authentication fields:**
>
> Set these when the entry has an [`archive` source](#zip-archives) on a server that requires credentials.
>
> | Field           | Type   | Description                                                                                                                                                                                                                                                                                         |
> | :-------------- | :----- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `headers`       | object | HTTP headers Claude Code sends when it downloads this entry's archive. Overrides the marketplace's headers of the same name. Requires Claude Code v2.1.238 or later.                                                                                                                                |
> | `headersHelper` | string | Command that prints the HTTP headers for this entry's archive download as one JSON object, for a credential that expires. See [Authenticate archive downloads](#authenticate-archive-downloads). The entry must also set [`"strict": false`](#strict-mode). Requires Claude Code v2.1.238 or later. |
>
> ## Plugin sources
>
> Plugin sources tell Claude Code where to get each individual plugin listed in your marketplace. These are set in the `source` field of each plugin entry in `marketplace.json`.
>
> Claude Code copies each installed plugin into the local versioned plugin cache at `~/.claude/plugins/cache`, except for a [`command` source in link mode](#copy-mode-and-link-mode), which Claude Code uses in place. Claude Code also [installs the plugin's eligible Node.js package dependencies](/docs/en/plugins-reference#node-js-package-dependencies) into the cached copy.
>
> | Source        | Type                            | Fields                             | Notes                                                                                                                                                                                                                                               |
> | ------------- | ------------------------------- | ---------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | Relative path | `string` (e.g. `"./my-plugin"`) | none                               | Local directory within the marketplace repo. Must start with `./`, unless you write a [bare name under `metadata.pluginRoot`](#relative-paths). Claude Code resolves the path relative to the marketplace root, not the `.claude-plugin/` directory |
> | `github`      | object                          | `repo`, `ref?`, `sha?`             |                                                                                                                                                                                                                                                     |
> | `url`         | object                          | `url`, `ref?`, `sha?`              | Git URL source                                                                                                                                                                                                                                      |
> | `git-subdir`  | object                          | `url`, `path`, `ref?`, `sha?`      | Subdirectory within a git repo. Clones sparsely to minimize bandwidth for monorepos                                                                                                                                                                 |
> | `npm`         | object                          | `package`, `version?`, `registry?` | Installed via `npm install`                                                                                                                                                                                                                         |
> | `archive`     | object                          | `url`, `sha256?`                   | Zip archive downloaded over HTTPS. Works without git or npm on the user's machine. Requires Claude Code v2.1.224 or later                                                                                                                           |
> | `command`     | object                          | `command`, `timeout?`, `mode?`     | Plugin directory produced by running a local command, re-run once per session to pick up changes. Requires Claude Code v2.1.229 or later                                                                                                            |
>
>   **Marketplace sources vs plugin sources**: These are different concepts that control different things.
>
>   * **Marketplace source**: where to fetch the `marketplace.json` catalog itself. Set when users run `/plugin marketplace add` or in `extraKnownMarketplaces` settings. Git-based marketplace sources support `ref` (branch/tag) but not `sha`.
>   * **Plugin source**: where to fetch an individual plugin listed in the marketplace. Set in the `source` field of each plugin entry inside `marketplace.json`. Git-based plugin sources support both `ref` (branch/tag) and `sha` (exact commit).
>
>   For example, a marketplace hosted at `acme-corp/plugin-catalog` (marketplace source) can list a plugin fetched from `acme-corp/code-formatter` (plugin source). The marketplace source and plugin source point to different repositories and are pinned independently.
>
> The git-based source types below are `github`, `url`, and `git-subdir`. When both `ref` and `sha` are set on any of them, the `sha` is the effective pin. Claude Code fetches and checks out the pinned commit directly.
>
> On most git hosts, including GitHub, GitLab, and Bitbucket, this means installation succeeds even if the branch or tag named by `ref` has since been deleted upstream, as long as the commit is still reachable from the repository. Some servers, such as AWS CodeCommit, don't support fetching commits by SHA. On those servers the `ref` must still exist and the pinned commit must be reachable from it.
>
> If you distribute plugins through **Organization settings > Plugins**, only some source types are allowed. See [Distribute through organization settings](#distribute-through-organization-settings).
>
> ### Relative paths
>
> For plugins in the same repository, use a path starting with `./`:
>
> ```json
> {
>   "name": "my-plugin",
>   "source": "./plugins/my-plugin"
> }
> ```
>
> Paths resolve relative to the marketplace root, which is the directory containing `.claude-plugin/`. In the example above, `./plugins/my-plugin` points to `<repo>/plugins/my-plugin`, even though `marketplace.json` lives at `<repo>/.claude-plugin/marketplace.json`. Don't use `../` to reference paths outside the marketplace root.
>
> A bare name is a single directory name with no `/`, such as `"formatter"`. To write bare names instead of `./` paths, set [`metadata.pluginRoot`](#optional-fields) to the directory they resolve under. With `"pluginRoot": "./plugins"`, Claude Code resolves `"source": "formatter"` to `./plugins/formatter`. Requires Claude Code v2.1.239 or later.
>
> `metadata.pluginRoot` must itself be a relative path inside the marketplace. Claude Code ignores it for a source that already starts with `./`. A source that contains a `/`, such as `team-a/formatter`, isn't a bare name and still needs the `./` prefix, even when `metadata.pluginRoot` is set.
>
>   Claude Code resolves relative paths against a local copy of the marketplace, so they work when users add your marketplace from a git source or a local directory. If users add your marketplace via a direct URL to the `marketplace.json` file, relative paths won't resolve, because Claude Code downloads only that file. For URL-based distribution, use any other [plugin source](#plugin-sources) instead. See [Troubleshooting](#plugins-with-relative-paths-fail-in-url-based-marketplaces) for details.
>
> ### GitHub repositories
>
> ```json
> {
>   "name": "github-plugin",
>   "source": {
>     "source": "github",
>     "repo": "owner/plugin-repo"
>   }
> }
> ```
>
> You can pin to a specific branch, tag, or commit:
>
> ```json
> {
>   "name": "github-plugin",
>   "source": {
>     "source": "github",
>     "repo": "owner/plugin-repo",
>     "ref": "v2.0.0",
>     "sha": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0"
>   }
> }
> ```
>
> | Field  | Type   | Description                                                           |
> | :----- | :----- | :-------------------------------------------------------------------- |
> | `repo` | string | Required. GitHub repository in `owner/repo` format                    |
> | `ref`  | string | Optional. Git branch or tag (defaults to repository default branch)   |
> | `sha`  | string | Optional. Full 40-character git commit SHA to pin to an exact version |
>
> ### Git repositories
>
> ```json
> {
>   "name": "git-plugin",
>   "source": {
>     "source": "url",
>     "url": "https://gitlab.com/team/plugin.git"
>   }
> }
> ```
>
> You can pin to a specific branch, tag, or commit:
>
> ```json
> {
>   "name": "git-plugin",
>   "source": {
>     "source": "url",
>     "url": "https://gitlab.com/team/plugin.git",
>     "ref": "main",
>     "sha": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0"
>   }
> }
> ```
>
> | Field | Type   | Description                                                                                                                                              |
> | :---- | :----- | :------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `url` | string | Required. Full git repository URL (`https://` or `git@`). The `.git` suffix is optional, so Azure DevOps and AWS CodeCommit URLs without the suffix work |
> | `ref` | string | Optional. Git branch or tag (defaults to repository default branch)                                                                                      |
> | `sha` | string | Optional. Full 40-character git commit SHA to pin to an exact version                                                                                    |
>
> ### Git subdirectories
>
> Use `git-subdir` to point to a plugin that lives inside a subdirectory of a git repository. Claude Code uses a sparse, partial clone to fetch only the subdirectory, minimizing bandwidth for large monorepos.
>
> ```json
> {
>   "name": "my-plugin",
>   "source": {
>     "source": "git-subdir",
>     "url": "https://github.com/acme-corp/monorepo.git",
>     "path": "tools/claude-plugin"
>   }
> }
> ```
>
> You can pin to a specific branch, tag, or commit:
>
> ```json
> {
>   "name": "my-plugin",
>   "source": {
>     "source": "git-subdir",
>     "url": "https://github.com/acme-corp/monorepo.git",
>     "path": "tools/claude-plugin",
>     "ref": "v2.0.0",
>     "sha": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0"
>   }
> }
> ```
>
> The `url` field also accepts a GitHub shorthand (`owner/repo`) or SSH URLs (`git@github.com:owner/repo.git`).
>
> | Field  | Type   | Description                                                                                              |
> | :----- | :----- | :------------------------------------------------------------------------------------------------------- |
> | `url`  | string | Required. Git repository URL, GitHub `owner/repo` shorthand, or SSH URL                                  |
> | `path` | string | Required. Subdirectory path within the repo containing the plugin (for example, `"tools/claude-plugin"`) |
> | `ref`  | string | Optional. Git branch or tag (defaults to repository default branch)                                      |
> | `sha`  | string | Optional. Full 40-character git commit SHA to pin to an exact version                                    |
>
> ### npm packages
>
> Plugins distributed as npm packages are installed using `npm install`. This works with any package on the public npm registry or a private registry your team hosts.
>
> ```json
> {
>   "name": "my-npm-plugin",
>   "source": {
>     "source": "npm",
>     "package": "@acme/claude-plugin"
>   }
> }
> ```
>
> To pin to a specific version, add the `version` field:
>
> ```json
> {
>   "name": "my-npm-plugin",
>   "source": {
>     "source": "npm",
>     "package": "@acme/claude-plugin",
>     "version": "2.1.0"
>   }
> }
> ```
>
> To install from a private or internal registry, add the `registry` field:
>
> ```json
> {
>   "name": "my-npm-plugin",
>   "source": {
>     "source": "npm",
>     "package": "@acme/claude-plugin",
>     "version": "^2.0.0",
>     "registry": "https://npm.example.com"
>   }
> }
> ```
>
> | Field      | Type   | Description                                                                                  |
> | :--------- | :----- | :------------------------------------------------------------------------------------------- |
> | `package`  | string | Required. Package name or scoped package (for example, `@org/plugin`)                        |
> | `version`  | string | Optional. Version or version range (for example, `2.1.0`, `^2.0.0`, `~1.5.0`)                |
> | `registry` | string | Optional. Custom npm registry URL. Defaults to the system npm registry (typically npmjs.org) |
>
> ### Zip archives
>
> Use `archive` to distribute a plugin as a zip file that Claude Code downloads over HTTPS, so installs work without git or npm on the user's machine. Host the file on any static file server or artifact repository, such as an S3 bucket, an Artifactory generic repository, or nginx. Requires Claude Code v2.1.224 or later. On versions v2.1.120 through v2.1.223, installing the plugin fails with `This plugin uses a source type your Claude Code version does not support. Update Claude Code and try again.`; on older versions, a marketplace containing an `archive` entry fails to load entirely.
>
> This entry installs the plugin from a zip file on an artifact server:
>
> ```json
> {
>   "name": "my-plugin",
>   "source": {
>     "source": "archive",
>     "url": "https://artifacts.example.com/claude-plugins/my-plugin-2.1.0.zip"
>   }
> }
> ```
>
> When you build the zip, you can zip the plugin's contents directly or zip the plugin folder itself. Claude Code looks for `.claude-plugin/` at the top of the archive, then inside a single top-level folder, so both layouts install:
>
> ```text
> my-plugin.zip          my-plugin.zip
> ├── .claude-plugin/    └── my-plugin/
> │   └── plugin.json        ├── .claude-plugin/
> └── commands/              │   └── plugin.json
>                            └── commands/
> ```
>
> Claude Code doesn't look deeper than one folder, so a plugin nested further down fails to install. Claude Code refuses archives larger than 256 MiB.
>
> To pin the exact file, add a `sha256` field with the archive's digest:
>
> ```json
> {
>   "name": "my-plugin",
>   "source": {
>     "source": "archive",
>     "url": "https://artifacts.example.com/claude-plugins/my-plugin-2.1.0.zip",
>     "sha256": "6bfa50e3d2e00c052b46abe51fff89346ac803e45771f76dcf6df1ab74cca5e1"
>   }
> }
> ```
>
> If the downloaded file doesn't match the pin, Claude Code refuses the install and reports [`Plugin archive integrity check failed`](/docs/en/errors#plugin-archive-integrity-check-failed).
>
> Archive sources accept these fields:
>
> | Field    | Type   | Description                                                                                                                                                                                                                |
> | :------- | :----- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `url`    | string | Required. HTTPS URL of the zip archive. Claude Code rejects `http://` URLs, along with loopback, link-local, and cloud-metadata hosts. Every redirect hop must satisfy the same rules, or Claude Code refuses the download |
> | `sha256` | string | Optional. SHA-256 digest of the archive as 64 hex characters, uppercase or lowercase. Claude Code verifies every download against it and refuses the install on a mismatch                                                 |
>
> The `sha256` digest also serves as the plugin's version when neither `plugin.json` nor the marketplace entry declares one. See [Version management](/docs/en/plugins-reference#version-management). If you declare a `version`, that version string is the update signal, so after changing the zip and its digest, bump the version too, or users keep the cached copy.
>
> #### Authenticate archive downloads
>
> To authenticate an archive download, such as a download from a private registry, set the HTTP headers Claude Code sends with it. Set `headers` on the `url` source you registered the marketplace from, such as an [`extraKnownMarketplaces`](/docs/en/settings-reference#extraknownmarketplaces) entry. On Claude Code v2.1.238 or later, you can set it on the plugin's entry instead, beside `source`.
>
> If the value you would put in `headers` is short-lived, such as a token your registry mints on request, set a `headersHelper` command in the same place instead. Claude Code runs the command and sends the JSON object it prints as that place's headers. Requires Claude Code v2.1.238 or later.
>
> The place you choose decides which downloads get the headers and when Claude Code runs the command:
>
> | Place                    | Downloads that get the headers                                                             | When Claude Code runs a `headersHelper` set there                                                                                                                   |
> | :----------------------- | :----------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | Marketplace `url` source | Archive downloads on the marketplace URL's origin, meaning the same scheme, host, and port | Before each fetch of the marketplace's `marketplace.json` and before each archive download on that origin. Claude Code reuses one run's output for up to 60 seconds |
> | Plugin entry             | That entry's download only                                                                 | Only when a user installs or updates that one plugin by itself and [accepts the command](#how-users-accept-a-headershelper-command)                                 |
>
> Where both places set a header of the same name, Claude Code sends the entry's value. Within one place, a header the command prints overrides a header of the same name listed in `headers`.
>
> ##### Add a headersHelper to a plugin entry
>
> This entry sets `headersHelper` beside `source`. It also sets `"strict": false`, which Claude Code requires of a `marketplace.json` entry that sets `headersHelper`. With [`"strict": false`](#strict-mode), the marketplace entry is the plugin's entire definition, so a user can review what the plugin contains before accepting the command:
>
> ```json
> {
>   "name": "my-plugin",
>   "description": "Formatting commands for internal services",
>   "strict": false,
>   "commands": "./commands",
>   "source": {
>     "source": "archive",
>     "url": "https://registry.example.com/plugins/my-plugin-2.1.0.zip"
>   },
>   "headersHelper": "/opt/bin/mint-registry-token.sh"
> }
> ```
>
> To check the entry, run `claude plugin install my-plugin@your-marketplace`. Claude Code shows you the command and the archive URL, and downloads the zip after you accept.
>
> Before v2.1.238, Claude Code downloaded an entry's archive without its `headers` or `headersHelper`, so an install that relied on them failed with `HTTP 401 while downloading plugin archive from`, followed by the URL, with the registry's status code in place of 401.
>
> #### Write the headersHelper command
>
> Whether you set `headersHelper` on a marketplace's `url` source or on a plugin entry, write the command to meet these requirements:
>
> * **Command text**: at most 500 characters of printable ASCII, with no run of four or more spaces.
> * **Output**: print one JSON object of header names and string values on stdout, then exit 0 within 10 seconds.
> * **Shell and working directory**: Claude Code runs the command through `sh`, or `cmd.exe` on Windows, from the configuration directory, `~/.claude` or [`CLAUDE_CONFIG_DIR`](/docs/en/env-vars#variables). Give an absolute path or a command on `PATH`, because a relative path resolves against that directory, not the user's project.
> * **Variables Claude Code removes**: from the environment of a command set in a `marketplace.json` entry or in a project's `.claude/settings.json` or `.claude/settings.local.json`, Claude Code removes every variable whose name contains a word such as `TOKEN`, `SECRET`, `KEY`, or `AUTH`, including `ANTHROPIC_API_KEY`. Claude Code doesn't apply this removal to a command set in user settings, a `--settings` file, or managed settings.
> * **Variables Claude Code sets**: `CLAUDE_CODE_MARKETPLACE_URL` and `CLAUDE_CODE_MARKETPLACE_NAME` for a `url` source's command, and `CLAUDE_CODE_PLUGIN_NAME` and `CLAUDE_CODE_PLUGIN_ARCHIVE_URL` for an entry's command. `CLAUDE_CODE_MARKETPLACE_NAME` is unset on the first fetch after a user adds a marketplace by URL, because that fetch is what supplies the name.
>
> A command that mints a bearer token prints an object like this one:
>
> ```json
> {"Authorization": "Bearer eyJhbGciOiJSUzI1NiJ9"}
> ```
>
> #### When Claude Code skips a headersHelper command or drops its output
>
> Claude Code doesn't run a `headersHelper` command, or drops headers that came from `headers` or from the command's output, in these situations:
>
> * **Command fails**: if the command exits non-zero, runs past 10 seconds, or prints anything other than a JSON object of string values, Claude Code doesn't make the fetch or download it ran the command for.
> * **Marketplace URL doesn't start with `https://`**: Claude Code doesn't run that `url` source's command and sends only the headers listed in its `headers` field.
> * **Redirect leaves the origin**: when a download is redirected off the archive URL's origin, Claude Code drops the `headers` values and command output of both the marketplace `url` source and the plugin entry.
> * **Entry sets a routing or identity header**: Claude Code drops request-routing and client-identity names such as `Host`, `Cookie`, and `X-Forwarded-*` from an entry's `headers` and command output, and keeps authentication names such as `Authorization`. Claude Code filters every `marketplace.json` entry this way, and an [inline settings entry](/docs/en/settings-reference#extraknownmarketplaces) depending on which file declares it.
> * **Command set in an `--add-dir` directory's settings**: Claude Code ignores it, on a `url` source and on an [inline plugin entry](/docs/en/settings-reference#extraknownmarketplaces) alike, and sends only that file's `headers`.
> * **Managed settings block the command**: setting [`disableCommandPluginSources`](/docs/en/settings-reference#disablecommandpluginsources) to `true` blocks `headersHelper` commands, and [`allowManagedHooksOnly`](/docs/en/settings-reference#allowmanagedhooksonly) blocks them too unless `disableCommandPluginSources` is explicitly `false`. Under either block, Claude Code still runs the command for a marketplace that managed settings themselves declare.
>
> #### How users accept a headersHelper command
>
> A user accepts a plugin entry's command each time they install or update that one plugin by itself, from the plugin's own view in `/plugin` or with `claude plugin install` or `claude plugin update`. Claude Code shows the command and the archive URL, and runs the command only after the user accepts. In a non-interactive shell, pass [`--yes`](/docs/en/plugins-reference#plugin-install) to accept it.
>
> Claude Code runs only the command it showed, for the archive URL it showed. If the entry's command or archive URL changed in between, Claude Code refuses the install or update. A change in the query string alone doesn't count.
>
> ##### Installs and updates that refuse the command instead of asking
>
> On any operation other than a single-plugin install or update, Claude Code neither runs an entry's command nor downloads its archive, so the plugin stays at its installed version or stays uninstalled. What the user sees depends on the operation:
>
> * **Installing several plugins at once, from a plugin suggestion, or as another plugin's dependency**: Claude Code refuses the plugin that has the command and points the user at that plugin's own view in `/plugin`. The other plugins in a bulk install still install. A plugin that depends on the refused plugin fails to install until the user installs the refused plugin by itself.
> * **Background auto-update, or session start for a plugin whose archive was never downloaded**: Claude Code lists the plugin in the `/plugin` Errors tab so the user knows to install or update it by hand. An auto-update that finds the entry still advertises the installed version lists nothing.
>
> ##### When a marketplace `url` source's command runs
>
> A marketplace `url` source's `headersHelper` is declared in a settings file, such as an [`extraKnownMarketplaces`](/docs/en/settings-reference#extraknownmarketplaces) entry, rather than in the catalog the marketplace publishes, so Claude Code doesn't ask the user to accept it on each install or update. The settings file that declares it decides when Claude Code runs it:
>
> | Settings file                                                                 | When Claude Code runs the command                                                                                                                                                                                                            |
> | :---------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | User settings, a `--settings` file, or a managed settings file on the machine | Without asking, including during a background marketplace refresh                                                                                                                                                                            |
> | A project's `.claude/settings.json` or `.claude/settings.local.json`          | Only after the user accepts the [workspace trust dialog](/docs/en/permissions#what-runs-before-you-trust-a-folder) for that folder itself. A `-p` or SDK session doesn't count as accepting it, and neither does trust granted to a parent folder |
> | Server-managed settings                                                       | Only after the user approves the delivered settings in the [security approval dialog](/docs/en/server-managed-settings#security-approval-dialogs)                                                                                                 |
>
> In a `-p` or SDK session, Claude Code can't show the security approval dialog. It applies the other delivered settings, but the marketplace fetch, and any archive download that needs the command, fails until a user has approved in an interactive session.
>
> For an [inline plugin entry](/docs/en/settings-reference#extraknownmarketplaces) in one of these files, Claude Code requires the same folder trust or settings approval as for a marketplace-level command in that file, and the user also accepts the entry's command on each install or update.
>
> ### Command sources
>
> Use `command` when a locally installed tool produces the plugin directory, such as an IDE that renders its plugin for the currently selected toolchain. Claude Code runs the command when the user installs the plugin and re-runs it in the background once per session, so your users pick up the tool's changed output without reinstalling. Requires Claude Code v2.1.229 or later. On v2.1.120 through v2.1.228, installing the plugin fails with `This plugin uses a source type your Claude Code version does not support. Update Claude Code and try again.`, and on older versions the whole marketplace fails to load.
>
> This entry insta