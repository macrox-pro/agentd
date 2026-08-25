---
primary_sources:
  - id: T1-AGENT-TERMINAL
    title: "Tools"
    url: "https://cursor.com/docs/agent/tools/terminal.md"
    section: "Tools"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Agent tools

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Terminal tool

> # Terminal
>
> Cursor runs shell commands directly in your terminal. Your [Run Mode](https://cursor.com/docs/agent/security/run-modes.md) controls when commands run, when Cursor asks, and when terminal commands enter the sandbox.
>
> ## Sandbox
>
> The sandbox runs terminal commands in a restricted environment that blocks unauthorized file access and network activity. For platform requirements, network modes, environment variables, and `sandbox.json` configuration, read [Run Modes > Sandboxing](https://cursor.com/docs/agent/security/run-modes.md#sandboxing).
>
> ## Troubleshooting
>
> Some shell themes (for example, Powerlevel9k/Powerlevel10k) can interfere with
> the inline terminal output. If your command output looks truncated or
> misformatted, disable the theme or switch to a simpler prompt when Cursor runs.
>
> ### Disable heavy prompts for Cursor sessions
>
> Use the `CURSOR_AGENT` environment variable in your shell config to detect when
> Cursor is running and skip initializing fancy prompts/themes.
>
> ```zsh
> # ~/.zshrc - disable Powerlevel10k when Cursor runs
> if [[ -n "$CURSOR_AGENT" ]]; then
>   # Skip theme initialization for better compatibility
> else
>   [[ -r ~/.p10k.zsh ]] && source ~/.p10k.zsh
> fi
> ```
>
> ```bash
> # ~/.bashrc - fall back to a simple prompt in Cursor sessions
> if [[ -n "$CURSOR_AGENT" ]]; then
>   PS1='\u@\h \W \$ '
> fi
> ```
>
> ## Related
>
> - [Terminal help](https://cursor.com/help/ai-features/terminal.md)
>
>

### Source: Browser tool

> # Browser
>
> Agent can control a web browser to test applications, audit accessibility, convert designs into code, and more. With full access to console logs and network traffic, Agent can debug issues and automate comprehensive testing workflows.
>
> For enterprise customers, browser controls are governed by MCP allowlist or denylist.
>
> ## Native integration
>
> Agent displays browser actions like screenshots and actions in the chat, as well as the browser window itself either in a separate window or an inline pane.
>
> We've optimized the browser tools to be more efficient and reduce token usage, as well as:
>
> - **Efficient log handling**: Browser logs are written to files that Agent can grep and selectively read. Instead of summarizing verbose output after every action, Agent reads only the relevant lines it needs. This preserves full context while minimizing token usage.
> - **Visual feedback with images**: Screenshots are integrated directly with the file reading tool, so Agent actually sees the browser state as images rather than relying on text descriptions. This enables better understanding of visual layouts and UI elements.
> - **Smart prompting**: Agent receives additional context about browser logs, including total line counts and preview snippets, helping it make informed decisions about what to inspect.
> - **Development server awareness**: Agent is prompted to detect running development servers and use the correct ports instead of starting duplicate servers or guessing port numbers.
>
> You can use Browser without installing or configuring any external tools.
>
> ## Browser capabilities
>
> Agent has access to the following browser tools:
>
> ### Navigate
>
> Visit URLs and browse web pages. Agent can navigate anywhere on the web by visiting URLs, following links, going back and forward in history, and refreshing pages.
>
> ### Click
>
> Interact with buttons, links, and form elements. Agent can identify and interact with page elements, performing click, double-click, right-click, and hover actions on any visible element.
>
> ### Type
>
> Enter text into input fields and forms. Agent can fill out forms, submit data, and interact with form fields, search boxes, and text areas.
>
> ### Scroll
>
> Navigate through long pages and content. Agent can scroll to reveal additional content, find specific elements, and explore lengthy documents.
>
> ### Screenshot
>
> Capture visual representations of web pages. Screenshots help Agent understand page layout, verify visual elements, and provide you with confirmation of browser actions.
>
> ### Console Output
>
> Read browser console messages, errors, and logs. Agent can monitor JavaScript errors, debugging output, and network warnings to troubleshoot issues and verify page behavior.
>
> ### Network Traffic
>
> Monitor HTTP requests and responses made by the page. Agent can track API calls, analyze request payloads, check response status codes, and diagnose network-related issues. This is currently only available in the Agent panel, coming soon to the layout.
>
> ## Session persistence
>
> Browser state persists between Agent sessions based on your workspace. This means:
>
> - **Cookies**: Authentication cookies and session data remain available across browser sessions
> - **Local Storage**: Data stored in `localStorage` and `sessionStorage` persists
> - **IndexedDB**: Database content is retained between sessions
>
> The browser context is isolated per workspace, ensuring that different projects maintain separate storage and cookie states.
>
> ## Use cases
>
> ### Web development workflow
>
> Browser integrates into web development workflows alongside tools like Figma and Linear. See the [Web Development cookbook](https://cursor.com/for/web-development.md) for a complete guide on using Browser with design systems, project management tools, and component libraries.
>
> ### Accessibility improvements
>
> Agent can audit and improve web accessibility to meet WCAG compliance standards.
>
> @browser Check color contrast ratios, verify semantic HTML and ARIA labels, test keyboard navigation, and identify missing alt text
>
> ### Automated testing
>
> Agent can execute comprehensive test suites and capture screenshots for visual regression testing.
>
> @browser Fill out forms with test data, click through workflows, test responsive designs, validate error messages, and monitor console for JavaScript errors
>
> ### Design to code
>
> Agent can convert designs into working code with responsive layouts.
>
> @browser Analyze this design mockup, extract colors and typography, and generate pixel-perfect HTML and CSS code
>
> ### Adjusting UI design from screenshots
>
> Agent can refine existing interfaces by identifying visual discrepancies and updating component styles.
>
> @browser Compare current UI against this design screenshot and adjust spacing, colors, and typography to match
>
> ## Security
>
> Browser runs as a secure web view and is controlled using an MCP server running as an extension. Multiple layers protect you from unauthorized access and malicious actions.
> Cursor's Browser integrations have also been reviewed by multiple external security auditors.
>
> ### Authentication and isolation
>
> The browser implements several security measures:
>
> - **Token authentication**: Agent layout generates a random authentication token before each browser session starts
> - **Tab isolation**: Each browser tab receives a unique random ID to prevent cross-tab interference
> - **Session-based security**: Tokens regenerate for each new browser session
>
> ### Tool approval
>
> Browser tools require your approval by default. Review each action before Agent executes it. This prevents unexpected navigation, data submission, or script execution.
>
> You can configure approval settings in Agent Settings. Available modes:
>
> | Mode                     | Description                                                                 |
> | :----------------------- | :-------------------------------------------------------------------------- |
> | **Manual approval**      | Review and approve each browser action individually (recommended)           |
> | **Allow-listed actions** | Actions matching your allow list run automatically; others require approval |
> | **Auto-run**             | All browser actions execute immediately without approval (use with caution) |
>
> ### Allow and block lists
>
> Browser tools integrate with Cursor's [security guardrails](https://cursor.com/docs/agent/security.md). Configure which browser actions run automatically:
>
> - **Allow list**: Specify trusted actions that skip approval prompts
> - **Block list**: Define actions that should always be blocked
> - Access settings through: **Cursor Settings > Agents > Auto-Run**
>
> The allow/block list system provides best-effort protection. AI behavior can be unpredictable due to prompt injection and other issues. Review auto-approved actions regularly.
>
> Never use auto-run mode with untrusted code or unfamiliar websites. Agent could execute malicious scripts or submit sensitive data without your knowledge.
>
> ### Browser context
>
> The browser opens as a pane within Cursor, giving Agent full control through MCP tools.
>
> ## Recommended models
>
> We recommend using Sonnet 4.5, GPT-5, and Auto for the best performance.
>
> ## Enterprise usage
>
> For enterprise customers, browser functionality is managed through toggling availability under MCP controls. Admins have granular controls over each MCP server, as well as over browser access.
>
> ### Enabling browser for enterprise
>
> To enable browser capabilities for your enterprise team:
>
> 1. Navigate to your [Settings Dashboard](https://cursor.com/dashboard/settings)
> 2. Go to **MCP Configuration**
> 3. Toggle "browser features"
>
> Once configured, users in your organization will have access to browser tools based on your MCP allowlist or denylist settings.
>
> ### Origin allowlist
>
> Enterprise administrators can configure an origin allowlist that restricts which sites the agent can automatically navigate to and where MCP tools can run. This provides granular control over browser access for security and compliance.
>
> The Browser Origin Allowlist feature must be enabled for your organization before it appears in your dashboard. Contact your Cursor account team to request access.
>
> #### Configuration
>
> To configure the origin allowlist:
>
> 1. Navigate to your [Admin Dashboard](https://cursor.com/dashboard/settings)
> 2. Go to **MCP Configuration**
> 3. Ensure **Enable Browser Automation Features (v2.0+)** is enabled
> 4. Under **Browser Origin Allowlist (v2.1+)**, click **Add Origin**
> 5. Enter the origins you want to allow (e.g., `*`, `http://localhost:3000`, `https://internal.example.com`)
>
> Leave the allowlist empty to allow all origins. Each origin should be added separately using the Add Origin button.
>
> ![MCP Configuration showing Browser Origin Allowlist settings with Add Origin button](/docs-static/images/agent/browser-origin-allowlist.png)
>
> #### Behavior
>
> When an origin allowlist is configured:
>
> - **Automatic navigation**: The agent can only use the `browser_navigate` tool to visit URLs matching origins in the allowlist
> - **MCP tool execution**: MCP tools can only run on origins that are in the allowlist
> - **Manual navigation**: Users can still manually navigate the browser to any URL, including origins outside the allowlist (useful for viewing documentation or inspecting external sites)
> - **Tool restrictions**: Once the browser is on an origin not in the allowlist, browser tools (click, type, navigate) are blocked, even if the user navigated there manually
>
> #### Edge cases
>
> The origin allowlist provides best-effort protection. Be aware of these behaviors:
>
> - **Link navigation**: If the agent clicks a link on an allowed domain that navigates to a non-allowed origin, the navigation will succeed
> - **Redirects**: If the agent navigates to an allowed origin that subsequently redirects to a non-allowed origin, the redirect will be permitted
> - **JavaScript navigation**: Client-side navigation (via `window.location` or similar) from an allowed origin to a non-allowed origin will succeed
>
> The origin allowlist restricts automatic agent navigation but cannot prevent all navigation paths. Review your allowlist regularly and consider the security implications of allowing access to domains that may redirect or link to external sites.
>
> ## Related
>
> - [Browser tool help](https://cursor.com/help/ai-features/browser.md)
>
>

### Source: Search tool

> # Search
>
> ## Instant Grep
>
> The fastest way to find code is an exact match: a function name, variable, error string, or regex pattern. Agent uses grep automatically when you reference specific symbols.
>
> Cursor ships with [Instant Grep](/changelog/2-1#instant-grep-beta), a custom search engine that outperforms `ripgrep` on large codebases. It runs automatically; no configuration needed.
>
> Instant Grep supports full regex and word-boundary matching, so Agent can construct patterns like `import.*PaymentService` or `PaymentFailedError` to trace references across files.
>
> ## Privacy and security
>
> File paths are encrypted before being sent to Cursor's servers. Code content is never stored in plaintext; it is held in memory during indexing, then discarded.
>
> ## Explore subagent
>
> Agent can spawn an [Explore subagent](https://cursor.com/docs/subagents.md) that runs in its own context window with a faster model. It executes many parallel searches without bloating the main conversation, returning only the relevant findings.
>
> Agent uses the Explore subagent automatically when it decides a task benefits from broad search. You can also request it directly: "use a subagent to find all the places we validate user input."
>
> This is useful for context management. Searching through many files generates a lot of context. The subagent keeps the main conversation focused by summarizing results instead of dumping raw file contents.
>
> ## FAQ
>
> ### Is my source code stored on Cursor servers?
>
> No. Cursor creates embeddings without storing filenames or source code. Filenames are obfuscated and code chunks are encrypted. When Agent searches, Cursor retrieves the embeddings and decrypts the chunks on the client side.
>
> ### Can I customize path encryption?
>
> Create a `.cursor/keys` file in your workspace root:
>
> ```json
> {
>   "path_decryption_key": "your-custom-key-here"
> }
> ```
>
> ### How does team sharing work?
>
> Indexes can be shared across team members for faster indexing of similar codebases. Cursor respects file access permissions and only shares accessible content.
>
> ### Does Cursor support multi-root workspaces?
>
> Yes. Cursor supports [multi-root workspaces](https://code.visualstudio.com/docs/editor/workspaces#_multiroot-workspaces). All codebases get indexed automatically, and each codebase's context is available to Agent. Some features that rely on a single git root, like worktrees, are disabled for multi-root workspaces. Cloud Agents do not support multi-root workspaces.
>
>

### Source: Canvas tool

> # Canvases
>
> Canvases let Cursor create interactive artifacts that render next to the chat. Instead of scrolling through a long markdown table or code block, you get a standalone view, laid out with sections, stats, and tables, that you can reopen, edit, and iterate on.
>
> Ask agents for a dashboard, analysis, audit, or report, and Cursor opens the result in a canvas when that is a better fit.
>
> ## How it works
>
> 1. Cursor decides that your task benefits from a visual or interactive view, or you ask for one directly.
> 2. Cursor builds the canvas and inserts a reference to it in your chat.
> 3. You review the rendered view, switch to the source to tweak it, or ask Cursor to change it.
> 4. Cursor saves the canvas so you can reopen and rerun it later with fresh data.
>
> Each canvas appears in your workspace's canvas list, so you can jump back to past ones without rerunning them.
>
> ## Opening a canvas
>
> - **From Cursor**: when Cursor creates a canvas, a card appears at the end of the response. Click it to open.
> - **Command Palette**: run **Open Canvas** from the palette, listed under View.
> - **Agents Window**: open a canvas tab directly from the new tab menu in the [Agents Window](https://cursor.com/docs/agent/agents-window.md).
>
> ## Sharing canvases
>
> Shared canvases turn an interactive artifact into something your whole team can open, not just you. When you share a canvas, Cursor uploads a live snapshot of the view and gives you a link teammates can open in the browser — same layout, charts, and tables, without rerunning the agent or digging through chat history. Use **Publish** from the canvas toolbar to publish or refresh a share; browse everything your team has published from **Shared Canvases** on the [dashboard](https://cursor.com/dashboard).
>
> Shared canvases are available on paid plans (Pro, Teams, and Enterprise). Free accounts cannot create shares. Because each share is team-visible, you need to be on a team — Pro users on a team can share too. Sharing also requires a privacy mode that allows data storage (Legacy Privacy Mode blocks it).
>
> Team admins can turn shared canvases off for the organization from [team settings](https://cursor.com/dashboard/settings#shared-canvases) under **Shared Canvases**.
>
> ## Iterating on a canvas
>
> Canvases are designed to be easy to refine.
>
> - If the layout isn't right, tell Cursor what to change instead of editing by hand.
> - If the numbers look stale or off, ask Cursor to rerun the underlying query or show its work.
> - For larger reworks, revert and prompt Cursor again with more details. This is usually faster than nudging through small follow-ups.
> - For small tweaks, you can also manually edit the source code.
>
> ## Packaging in skills
>
> Common canvas workflows can be packaged as [skills](https://cursor.com/docs/skills.md) so Cursor produces a consistent layout every time you ask.
>
> A canvas skill typically includes:
>
> - **A trigger description** so Cursor knows when to reach for it, like "quarterly revenue report" or "dependency audit".
> - **Layout instructions** that define the sections, stats, and tables the canvas should contain.
> - **Data sources and queries** Cursor should run to populate the view, such as a SQL query, API call, or shell command.
> - **Formatting rules** like units, date ranges, or sort order.
>
> Once the skill is in place, a short prompt is enough to regenerate the canvas with fresh data, and every teammate using the skill gets the same output shape.
>
> ## Related
>
> - [Agents Window](https://cursor.com/docs/agent/agents-window.md)
> - [Skills](https://cursor.com/docs/skills.md)
> - [Prompting](https://cursor.com/docs/agent/prompting.md)
>
>

### Source: Agents window

> # Agents Window
>
> The Agents Window is Cursor's agent-first interface. It provides a unified workspace to build with agents across repos and environments, including local, cloud, remote SSH, and more. It combines the power of parallel agents with the depth and control of a development environment.
>
> You can switch back to the editor anytime, or have both open simultaneously.
>
> ## Open the Agents Window
>
> If you're in the editor, type Cmd+Shift+P → Open Agents Window to open the Agents Window.
>
> ![Command Palette showing the Open Agents Window command](/docs-static/images/agent/open-agents-window-final.png)
>
> ## Switch Back to the IDE
>
> To return to the classic Cursor IDE, type Cmd+Shift+P → Open IDE. This opens the current workspace in the editor.
>
> ![Actions menu showing the Open IDE command](/docs-static/images/agent/open-editor-window-final.png)
>
> If you want to view or edit files without leaving the Agents Window, you can type Cmd+P to search files, or Cmd+Shift+F to search all files.
>
> ![Agents Window showing file search and file viewing](/docs-static/images/agent/file-agents-window-final.png)
>
> ## Features Available Only in the Agents Window
>
> The following features are available in the Agents Window:
>
> - **Multi-workspace:** work with agents across all your projects from one place.
> - **New diffs view:** review and commit changes, and manage PRs without leaving Cursor.
> - **Parallel agents:** run many parallel agents in the cloud (and work with them from your phone, web, Slack, GitHub, and Linear).
> - **Easier handoff between local and cloud:** quickly move an agent from cloud to local to iterate quickly, and move it back to the cloud so it keeps working on its own.
> - **Cloud subagents:** hand off a task to a [cloud subagent](https://cursor.com/docs/subagents.md#cloud-subagents) with `/in-cloud`, or `/babysit` a PR, so long-running work runs on its own VM and branch while you keep working locally.
> - **Worktrees:** [run agents in isolated Git checkouts](https://cursor.com/docs/configuration/worktrees.md) so each task has its own files and changes.
>
> ## Choosing Between Agents Window and Editor
>
> The Agents Window works well when you want to run and manage many agents in parallel. If you are using agents to write most of your code, the Agents Window helps pull you up to a higher level of abstraction.
>
> The editor works well when you want the classic IDE with VS Code extensions and flexible screen splitting to see many files at once.
>
> You can move between the two interfaces, and we will continue to support and improve both experiences.
>
> ## Enterprise access
>
> Agents Window is generally available with Cursor 3, released on April 2, 2026. For the two weeks following launch, Enterprise Admins can control rollout within their organizations by giving access to their entire team or to specific users via Team settings. After the rollout period, all users will have access by default.
>
>

### Source: Agent review

> # Agent Review
>
> Agent Review runs a dedicated code review on your local changes from inside Cursor.
>
> ## Setup
>
> To configure Agent Review:
>
> 1. Open **Cursor Settings**
> 2. Go to **Agents**
> 3. Find **Agent Review** and configure your preferences
>
> Starting in Cursor 3.11, this setting moves to **Git & PRs** > **Pull Requests**.
>
> Agent Review also reads repository rules from `BUGBOT.md` files. To set up these rule files, see [BugBot docs](https://cursor.com/docs/bugbot.md).
>
> You can set it to run automatically after every agent task, or leave it manual and trigger it yourself.
>
> ## Running a review
>
> There are three ways to start a review:
>
> - **Automatic**: When enabled in settings, Agent Review runs after every commit is made.
> - **Slash command**: Type `/agent-review` in the agent window input to trigger a review on demand.
> - **Source Control tab**: Open the Source Control tab and run Agent Review to compare all local changes against your main branch. This catches issues across your full set of changes, not only the latest edit.
>
> [Media](https://ptht05hbb1ssoooe.public.blob.vercel-storage.com/assets/changelog/changelog-2-1-1.mp4)
>
> ## Review depth
>
> Agent Review supports two depth levels. Choose based on the thoroughness of review you need.
>
> | Depth     | Speed | Cost | Best for                                                   |
> | :-------- | :---- | :--- | :--------------------------------------------------------- |
> | **Quick** | Fast  | Low  | Small diffs, formatting changes, or a fast sanity check    |
> | **Deep**  | Slow  | High | Complex logic, security-sensitive code, or large refactors |
>
>
