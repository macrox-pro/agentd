---
primary_sources:
  - id: T1-AGENT-OVERVIEW
    title: "Overview and modes"
    url: "https://cursor.com/docs/agent/overview.md"
    section: "Overview and modes"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Agent overview and modes

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Agent overview

> # Cursor Agent
>
> Agent is Cursor's assistant that can complete complex coding tasks independently, run terminal commands, and edit code. Access in sidepane with Cmd+I.
>
> Learn more about [how agents work](https://cursor.com/learn/agents.md) and help you build faster.
>
> ## How Agent works
>
> An agent is built on three components:
>
> 1. **Instructions**: The system prompt and [rules](https://cursor.com/docs/rules.md) that guide agent behavior
> 2. **Tools**: File editing, codebase search, terminal execution, and more
> 3. **Model**: The agent model you pick for the task
>
> Cursor's agent orchestrates these components for each model we support, tuning instructions and tools specifically for every frontier model. As new models are released, you can focus on building software while Cursor handles the model-specific optimizations.
>
> ## Tools
>
> Tools are the building blocks of Agent. They are used to search your codebase and the web to find relevant information, make edits to your files, run terminal commands, and more.
>
> To understand how tool calling works under the hood, see our [tool calling fundamentals](https://cursor.com/learn/tool-calling.md).
>
> There is no limit on the number of tool calls Agent can make during a task.
>
> ### Search files and folders
>
> Search for files by name, read directory structures, and find exact keywords or patterns within files.
>
> ### Web
>
> Generate search queries and perform web searches.
>
> ### Fetch Rules
>
> Retrieve specific [rules](https://cursor.com/docs/rules.md) based on type and description.
>
> ### Read files
>
> Intelligently read the content of a file. Also supports image files (.png, .jpg, .gif, .webp, .svg) and includes them in the conversation context for analysis by vision-capable models.
>
> ### Edit files
>
> Suggest edits to files and apply them automatically.
>
> ### Run shell commands
>
> Execute terminal commands and monitor output. By default, Cursor uses the first terminal profile available.
>
> To set your preferred terminal profile:
>
> 1. Open Command Palette (`Cmd/Ctrl+Shift+P`)
> 2. Search for "Terminal: Select Default Profile"
> 3. Choose your desired profile
>
> ### Browser
>
> Control a browser to take screenshots, test applications, and verify visual changes. Agent can navigate pages, interact with elements, and capture the current state for analysis. See the [Browser documentation](https://cursor.com/docs/agent/tools/browser.md) for details.
>
> ### Image generation
>
> Generate images from text descriptions or reference images. Useful for creating UI mockups, product assets, and visualizing architecture diagrams. Images are saved to your project's `assets/` folder by default and shown inline in chat.
>
> ### Ask questions
>
> Ask clarifying questions during a task. While waiting for your response, the agent continues reading files, making edits, or running commands. Your answer is incorporated as soon as it arrives.
>
> ## Checkpoints
>
> Checkpoints save snapshots of your codebase during an Agent session. Agent automatically creates them before making significant changes, capturing the state of all modified files.
>
> If Agent takes a wrong turn, click any checkpoint in the chat timeline to preview your files at that point, then restore to revert all files to that state. You can also restore from the `Restore Checkpoint` button on previous requests or the + button when hovering over a message.
>
> Checkpoints are useful for exploratory work, complex refactoring, and iterative development where you want safe rollback points.
>
> Checkpoints are stored locally and separate from Git. Only use them for undoing Agent changes; use Git for permanent version control.
>
> ## Queued messages
>
> You have two ways to talk to an agent while it works. Queue a message and it waits for the current task to finish. [Send a follow-up now](https://cursor.com/docs/agent/overview.md#steer-a-running-agent) and it steers the active turn at the agent's next tool call.
>
> [Media](/docs-static/images/agent/planning/agent-queue.mp4)
>
> ### Using the queue
>
> 1. While Agent is working, type your next instruction
> 2. Press Enter to add it to the queue
> 3. Messages appear in order below the active task
> 4. Drag to reorder queued messages as needed
> 5. Agent processes them sequentially after finishing
>
> ### Keyboard shortcuts
>
> While Agent is working:
>
> - Press Enter to queue your message (it waits until Agent finishes the current task)
> - Press Cmd+Enter to send immediately, bypassing the queue
>
> ### Immediate messaging
>
> When you use Cmd+Enter to send immediately, your message is appended to the most recent user message in the chat and processed right away without waiting in the queue.
>
> - Your message attaches to tool results and sends immediately
> - This creates a more responsive experience for urgent follow-ups
> - Use this when you need to interrupt or redirect Agent's current work
>
> ### Steer a running agent
>
> You can send a follow-up to steer the agent while it's working, without interrupting it. Type a follow-up and hit **Send now**, or press Enter twice. The message is delivered at the agent's next tool call instead of cutting off work mid-action, which preserves in-flight work and keeps the agent on task.
>
> This is available on [cursor.com/agents](https://cursor.com/agents) now and rolling out in the [Agents Window](https://cursor.com/docs/agent/agents-window.md). Press Tab to queue the message for after the turn instead.
>
> In the [CLI](https://cursor.com/docs/cli/overview.md), pressing Enter while the agent works steers the active run at a safe boundary, and pressing Enter again interrupts the turn.
>
> ## Goals with /goal
>
> Agent reads each message as a new job. Use `/goal` to give the agent a long-lived objective to work towards until it's fully complete:
>
> ```text
> /goal fix all flaky tests and make CI green
> ```
>
> In the [CLI](https://cursor.com/docs/cli/overview.md), Ctrl+C pauses the goal. Pair a goal with a [Custom Mode](https://cursor.com/docs/agent/prompting.md#custom-modes) when you want the agent to follow a playbook, or with the built-in [`/loop`](https://cursor.com/docs/skills.md#built-in-cursor-skills) skill for recurring check-ins while it pursues the objective.
>
> `/goal` is rolling out. If you don't see it, try it in a new chat.
>
>

### Source: Plan Mode

> # Plan Mode
>
> Plan Mode creates detailed implementation plans before writing any code. Agent researches your codebase, asks clarifying questions, and generates a reviewable plan you can edit before building.
>
> Press Shift+Tab from the chat input to rotate to Plan Mode. Cursor also suggests it automatically when you type keywords that indicate complex tasks.
>
> ## How it works
>
> 1. Agent asks clarifying questions to understand your requirements
> 2. Researches your codebase to gather relevant context
> 3. Creates a comprehensive implementation plan
> 4. You review and edit the plan through chat or markdown files
> 5. Click to build the plan when ready
>
> Plans are saved by default in your home directory. Click "Save to workspace" to move it to your workspace for future reference, team sharing, and documentation.
>
> ## When to use Plan Mode
>
> Plan Mode works best for:
>
> - Complex features with multiple valid approaches
> - Tasks that touch many files or systems
> - Unclear requirements where you need to explore before understanding scope
> - Architectural decisions where you want to review the approach first
>
> For quick changes or tasks you've done many times before, jumping straight to Agent mode is fine.
>
> ## Starting over from a plan
>
> Sometimes Agent builds something that doesn't match what you wanted. Instead of trying to fix it through follow-up prompts, go back to the plan.
>
> Revert the changes, refine the plan to be more specific about what you need, and run it again. This is often faster than fixing an in-progress agent, and produces cleaner results.
>
> For larger changes, spend extra time creating a precise, well-scoped plan. The hard part is often figuring out **what** change should be made. With the right instructions, delegate implementation to Agent.
>
> ## Switching modes
>
> - Use the mode picker dropdown in Agent
> - Press Shift+Tab for quick switching
>
> ## Related
>
> - [Plan mode help](https://cursor.com/help/ai-features/plan-mode.md)
>
>

### Source: Debug Mode

> # Debug Mode
>
> Debug Mode helps you find root causes and fix tricky bugs that are hard to reproduce or understand. Instead of immediately writing code, the agent generates hypotheses, adds log statements, and uses runtime information to pinpoint the exact issue before making a targeted fix.
>
> ## When to use Debug Mode
>
> Debug Mode works best for:
>
> - **Bugs you can reproduce but can't figure out**: When you know something is wrong but the cause isn't obvious from reading the code
> - **Race conditions and timing issues**: Problems that depend on execution order or async behavior
> - **Performance problems and memory leaks**: Issues that require runtime profiling to understand
> - **Regressions where something used to work**: When you need to trace what changed
>
> When standard Agent interactions struggle with a bug, Debug Mode provides a different approach using runtime evidence rather than guessing at fixes.
>
> ## How it works
>
> 1. **Explore and hypothesize**: The agent explores relevant files, builds context, and generates multiple hypotheses about potential root causes.
>
> 2. **Add instrumentation**: The agent adds log statements that send data to a local debug server running in a Cursor extension.
>
> 3. **Reproduce the bug**: Debug Mode asks you to reproduce the bug and provides specific steps. This keeps you in the loop and ensures the agent captures real runtime behavior.
>
> 4. **Analyze logs**: After reproduction, the agent reviews the collected logs to identify the actual root cause based on runtime evidence.
>
> 5. **Make targeted fix**: The agent makes a focused fix that directly addresses the root cause, often just a few lines of code.
>
> 6. **Verify and clean up**: You can re-run the reproduction steps to verify the fix. Once confirmed, the agent removes all instrumentation.
>
> ## Tips for Debug Mode
>
> - **Provide detailed context**: The more you describe the bug and how to reproduce it, the better the agent's instrumentation will be. Include error messages, stack traces, and specific steps.
> - **Follow reproduction steps exactly**: Execute the steps the agent provides to ensure logs capture the actual issue.
> - **Reproduce multiple times if needed**: Reproducing the bug multiple times may help the agent identify tricky problems like race conditions.
> - **Be specific about expected vs. actual behavior**: Help the agent understand what should happen versus what is happening.
>
> ## Switching modes
>
> - Use the mode picker dropdown in Agent
> - Press Shift+Tab for quick switching
>
> ## Related
>
> - [Debug mode help](https://cursor.com/help/ai-features/debug-mode.md)
>
>

### Source: Design Mode

> # Design Mode
>
> Design Mode lets you direct agents with visual prompts. From the browser in the [Agents Window](https://cursor.com/docs/agent/agents-window.md), you can click an element, draw on the page, or describe a change by voice. Cursor captures the context it needs and edits the code while you move on to the next change.
>
> UI work tends to be spatial. Instead of describing a change in a sentence, your instruction can include the selected element, the code behind it, the surrounding layout, and the visual relationships on the page. This tightens the loop between noticing something and fixing it.
>
> Click an element in the running app, prompt against that selected element, and let the agent edit the code.
>
> ## Open Design Mode
>
> Design Mode lives in the browser inside the Agents Window. Open the browser, then toggle Design Mode with Cmd + Shift + D. Toggle it off with the same shortcut to return to normal browsing.
>
> ## Ways to direct the agent
>
> Design Mode gives you several ways to convey intent.
>
> ### Select an element
>
> Click any element in the running product to target it. The agent gets the element and its code, so you can prompt against the exact thing you see without leaving the app.
>
> ### Select multiple elements
>
> Multi-select helps when the change depends on a relationship between elements. Reference two components and ask the agent to make one match the other, remove repeated content, or adjust a group together.
>
> Select multiple elements and describe how they should change together.
>
> ### Draw on the page
>
> Drawing tells the agent which area of the page your instruction applies to. Circle a crowded section, box in a region, or mark part of an animated page. The annotation sits over a frozen frame of the viewport, so the agent sees the exact page state you were responding to.
>
> ### Narrate by voice
>
> You can narrate instructions with your voice instead of typing. The mic stays available while agents run, so you can queue the next change without waiting.
>
> Use voice input and drawing together to describe a change.
>
> ## Keyboard shortcuts
>
> | Action               | Shortcut        |
> | :------------------- | :-------------- |
> | Toggle Design Mode   | Cmd + Shift + D |
> | Select an area       | Shift + drag    |
> | Add element to chat  | Cmd + L         |
> | Add element to input | Option + click  |
>
> ## What the agent sees
>
> Picking an element adds two complementary signals to context:
>
> - **Element identity**: the xpath, the component, attributes, computed styles, and props from the fiber tree. This helps the agent find the source and edit the right code.
> - **A screenshot**: the layout, surrounding elements, and the exact page state. This gives the agent spatial context for the change.
>
> ## Work in flow
>
> When you refine an interface, one edit usually leads to the next. You adjust a component, notice the spacing around it, then see how another component should match.
>
> Design Mode lets you send those edits away as you notice them. Point at one element, describe the change, move to another part of the page, and send another edit before the first one finishes. This makes it easy to multitask and manage several subagents at once. As agents finish, the app hot reloads and your changes appear in the running product.
>
> This flow works best with a fast model that is strong at interface work. We recommend [Composer 2.5](/blog/composer-2-5).
>
> ## Related
>
> - [Agents Window](https://cursor.com/docs/agent/agents-window.md)
> - [Browser](https://cursor.com/docs/agent/tools/browser.md)
>
>

### Source: Prompting

> # Prompting agents
>
> Direct Agent with text prompts in the chat input. You can attach context, images, and voice, and switch models at any point.
>
> ## @ mentions
>
> Type `@` in the chat input to attach specific context to your prompt. Start typing after `@` and Cursor shows matching suggestions.
>
> - **Files & Folders**: `@auth.ts` or `@src/components/` to include files or folders (type `/` after selecting a folder to navigate deeper)
> - **Terminals**: `@Terminals` to include terminal output as context
> - **Chats**: `@Chats` to reference context from a previous conversation
> - **Git diffs**: `@Commit (Diff of Working State)` for uncommitted changes, or `@Branch (Diff with Main)` for your full branch diff
> - **Browser**: `@Browser` to attach context from the built-in browser
>
> Use @ mentions when you know which files are relevant. If you're not sure which files matter, skip it — Agent finds relevant files through its own search.
>
> ## Custom Modes
>
> Type `/` in the chat input to invoke a [skill](https://cursor.com/docs/skills.md). Pressing Enter attaches the skill to one message, and it fades as the conversation moves on. Use any skill as a Custom Mode to keep the agent focused on it while it works.
>
> Pick the skill from the `/` menu and press Option+Enter (Mac) or Alt+Enter (Windows) instead. You can also select **Use as Mode** from the skill entry. Inside a mode, the skill stays in context on every turn, even as the agent works for hours, until you exit the mode.
>
> Custom Modes work well for skills that describe how to work rather than a one-shot task. Keep a code-review checklist active while you move through several files, or hold a team playbook like `/tdd` on for an entire feature.
>
> Custom Modes are available in the [Agents Window](https://cursor.com/docs/agent/agents-window.md) and the [CLI](https://cursor.com/docs/cli/overview.md). Any skill with a valid frontmatter block can back a mode, and the optional `icon` and `color` frontmatter fields style the mode's badge. See [Using a skill as a Custom Mode](https://cursor.com/docs/skills.md#using-a-skill-as-a-custom-mode).
>
> ## Image input
>
> Attach images to your prompt to provide visual context for UI work, debugging, and design implementation.
>
> - **Drag and drop** an image file into the chat input
> - **Paste from clipboard** with Cmd+V, including screenshots
>
> This is useful for implementing design mockups, debugging visual issues, and referencing error messages or stack traces without manual transcription.
>
> ## Voice input
>
> Click the microphone icon in the chat input to dictate your prompt instead of typing. Speak naturally, include technical details like file and function names, and review the transcription before sending.
>
> ## Context usage
>
> Every chat shares a fixed context window with the model. As you add files, run tools, and exchange messages, those tokens fill up. When the window gets close to full, Cursor compresses older parts of the conversation into a summary to leave more room for new conversation.
>
> The context ring next to your prompt input shows how full the window is at a glance. Click the ring to open the breakdown tray, which shows the total tokens used split by category:
>
> - **System prompt**: Cursor's built-in instructions for the model
> - **Tools**: definitions of every tool available to the agent
> - **Rules**: project and user rules included in the prompt
> - **Skills**: skill descriptions injected into the system context
> - **MCP**: instructions and catalog from connected MCP servers
> - **Subagents**: documentation for subagent types the agent can launch
> - **Summarized conversation**: compressed summaries of earlier turns
> - **Conversation**: your messages, the agent's replies, and tool results
>
> Hover a segment in the bar or a row in the list to highlight that category.
>
> ## Changing models
>
> Use the model picker dropdown at the top of the chat input to switch models, or press Cmd / to cycle through models. The change applies to the current conversation going forward. Set a default model in **Cursor Settings > Models**.
>
> - **Faster models** work well for quick edits and routine tasks
> - **More capable models** are better for complex reasoning and multi-file refactoring
>
> You can switch models mid-conversation, for example when a faster model handled exploration but you need deeper reasoning for implementation. See [Models & Pricing](https://cursor.com/docs/models-and-pricing.md) for the full list.
>
>
