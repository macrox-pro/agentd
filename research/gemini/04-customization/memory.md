---
primary_sources:
  - id: T1-AUTO-MEM
    title: "Auto Memory"
    url: "https://geminicli.com/docs/cli/auto-memory.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Memory and context

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Auto Memory — Full page

> Auto Memory is an experimental feature that mines your past Gemini CLI sessions
> in the background and proposes durable memory updates and reusable
> [Agent Skills](/docs/cli/skills). You review each candidate before it becomes
> available to future sessions: apply memory updates, promote skills, or discard
> anything you do not want.
>
> <!-- prettier-ignore -->
> > [!NOTE]
> > This is an experimental feature currently under active development.
>
> ## Overview
>
> Every session you run with Gemini CLI is recorded locally as a transcript. Auto
> Memory scans those transcripts for durable facts, preferences, workflow
> constraints, and procedural patterns that recur across sessions. It can draft
> memory updates as unified diff `.patch` files and draft reusable procedures as
> `SKILL.md` files. All candidates are held in a project-local inbox until you
> approve or discard them.
>
> You'll use Auto Memory when you want to:
>
> - **Capture team workflows** that you find yourself walking the agent through
>   more than once.
> - **Preserve durable project context** such as repeated verification commands,
>   local constraints, or personal project notes.
> - **Codify hard-won fixes** for project-specific landmines so future sessions
>   avoid them.
> - **Bootstrap a skills library** without writing every `SKILL.md` by hand.
>
> Auto Memory complements direct memory-file editing. The agent can still persist
> explicit user instructions by editing the appropriate Markdown memory file; Auto
> Memory infers candidates from past sessions, writes reviewable patches or skill
> drafts, and never applies them without your approval.
>
> ## Prerequisites
>
> - Gemini CLI installed and authenticated.
> - At least one idle project session with 10 or more user messages. Auto Memory
>   ignores active, trivial, and sub-agent sessions.
>
> ## How to enable Auto Memory
>
> Auto Memory is off by default. Enable it in your settings file:
>
> 1.  Open your global settings file at `~/.gemini/settings.json`. If you only
>     want Auto Memory in one project, edit `.gemini/settings.json` in that
>     project instead.
>
> 2.  Add the experimental flag:
>
>     ```json
>     {
>       "experimental": {
>         "autoMemory": true
>       }
>     }
>     ```
>
> 3.  Restart Gemini CLI. The flag requires a restart because the extraction
>     service starts during session boot.
>
> ## How Auto Memory works
>
> Auto Memory runs as a background task on session startup. It does not block the
> UI, consume your interactive turns, or surface tool prompts.
>
> 1.  **Eligibility scan.** The service indexes recent sessions from
>     `~/.gemini/tmp/<project>/chats/`. Sessions are eligible only if they have
>     been idle for at least three hours and contain at least 10 user messages.
> 2.  **Lock acquisition.** A lock file in the project's memory directory
>     coordinates across multiple CLI instances so extraction runs at most once at
>     a time. A state file records processed session versions, and extraction is
>     throttled so short back-to-back CLI launches do not repeatedly scan history.
> 3.  **Candidate extraction.** A background extraction agent reviews the session
>     index, reads any sessions that look like they contain durable memory or
>     repeated procedural workflows, and drafts candidates. It defaults to
>     creating no artifacts unless the evidence is strong, so many runs produce no
>     inbox items.
> 4.  **Safety boundaries.** Auto Memory writes candidates to a review inbox. It
>     cannot directly edit active memory files, settings, credentials, or project
>     `GEMINI.md` files.
> 5.  **Patch validation.** Skill update patches are parsed and dry-run before
>     they are surfaced. Memory patches are parsed, target-allowlisted, and
>     applied atomically only when you approve them from the inbox.
> 6.  **Notification.** When a run produces new candidates, Gemini CLI surfaces an
>     inline message telling you how many items are waiting.
>
> ## How to review extracted items
>
> Use the `/memory inbox` slash command to open the inbox dialog at any time:
>
> **Command:** `/memory inbox`
>
> The dialog groups pending items into new skills, skill updates, and memory
> updates. From there you can:
>
> - **Read** the full `SKILL.md` body before deciding.
> - **Promote** a skill to your user (`~/.gemini/skills/`) or workspace
>   (`.gemini/skills/`) directory.
> - **Discard** a skill you do not want.
> - **Apply** or reject a `.patch` proposal against an existing skill.
> - **Review** memory diffs before they touch active files.
> - **Apply** or dismiss private and global memory patches. Private patches target
>   the project memory directory; global patches target only your personal
>   `~/.gemini/GEMINI.md` file.
>
> Promoted skills become discoverable in the next session and follow the standard
> [skill discovery precedence](/docs/cli/skills#skill-discovery-tiers). Applied memory
> patches update the underlying memory files and reload memory for the current
> session.
>
> ## How to disable Auto Memory
>
> To turn off background extraction, set the flag back to `false` in your settings
> file and restart Gemini CLI:
>
> ```json
> {
>   "experimental": {
>     "autoMemory": false
>   }
> }
> ```
>
> Disabling the flag stops the background service immediately on the next session
> start. Existing inbox items remain on disk; you can either drain them with
> `/memory inbox` first or remove the project memory directory manually.
>
> ## Data and privacy
>
> - Auto Memory only reads session files that already exist locally on your
>   machine.
> - Auto Memory uses model calls to analyze selected local transcript content
>   during extraction. No candidates are applied automatically, but transcript
>   excerpts may be sent to the configured model as part of those calls.
> - The extraction agent is instructed to redact secrets, tokens, and credentials
>   it encounters and to never copy large tool outputs verbatim.
> - Drafted skills and memory patches live in your project's memory directory
>   until you promote, apply, dismiss, or discard them. They are not automatically
>   loaded into any session.
>
> ## Limitations
>
> - The extraction agent runs on a preview Gemini Flash model. Extraction quality
>   depends on the model's ability to recognize durable patterns versus one-off
>   incidents.
> - Auto Memory does not extract memory or skills from the current session. It
>   only considers sessions that have been idle for three hours or more.
> - Project or workspace shared instructions in project `GEMINI.md` files are not
>   auto-extractable. Auto Memory can propose private project memory, global
>   personal memory, and skills.
> - Inbox items are stored per project. Skills extracted in one workspace are not
>   visible from another until you promote them to the user-scope skills
>   directory.
>
> ## Next steps
>
> - Learn how skills are discovered and activated in [Agent Skills](/docs/cli/skills).
> - Explore the [memory management tutorial](/docs/cli/tutorials/memory-management) for
>   the complementary explicit-memory and `GEMINI.md` workflows.
> - Review the experimental settings catalog in
>   [Settings](/docs/cli/settings#experimental).

### Source: Memory files — Full page

> Gemini CLI persists durable facts, user preferences, and project details by
> editing Markdown memory files directly.
>
> ## Technical reference
>
> The agent routes memories to the appropriate Markdown file: shared project
> instructions go in repository `GEMINI.md` files, private project notes go in the
> per-project private memory folder, and cross-project personal preferences go in
> the global `~/.gemini/GEMINI.md` file.
>
> ## Technical behavior
>
> - **Storage:** Edits Markdown files with `write_file` or `replace`.
> - **Loading:** The stored facts are automatically included in the hierarchical
>   context system for all future sessions.
> - **Format:** Keeps durable instructions concise and avoids duplicating the same
>   fact across multiple memory tiers.
>
> ## Use cases
>
> - Persisting user preferences (for example, "I prefer functional programming").
> - Saving project-wide architectural decisions.
> - Storing frequently used aliases or system configurations.
>
> ## Next steps
>
> - Follow the [Memory management guide](/docs/cli/tutorials/memory-management)
>   for practical examples.
> - Learn how the [Project context (GEMINI.md)](/docs/cli/gemini-md) system loads
>   this information.

### Source: Memory Import Processor — Full page

> The Memory Import Processor is a feature that lets you modularize your GEMINI.md
> files by importing content from other files using the `@file.md` syntax.
>
> ## Overview
>
> This feature enables you to break down large GEMINI.md files into smaller, more
> manageable components that can be reused across different contexts. The import
> processor supports both relative and absolute paths, with built-in safety
> features to prevent circular imports and ensure file access security.
>
> ## Syntax
>
> Use the `@` symbol followed by the path to the file you want to import:
>
> ```markdown
> # Main GEMINI.md file
>
> This is the main content.
>
> @./components/instructions.md
>
> More content here.
>
> @./shared/configuration.md
> ```
>
> ## Supported path formats
>
> ### Relative paths
>
> - `@./file.md` - Import from the same directory
> - `@../file.md` - Import from parent directory
> - `@./components/file.md` - Import from subdirectory
>
> ### Absolute paths
>
> - `@/absolute/path/to/file.md` - Import using absolute path
>
> ## Examples
>
> ### Basic import
>
> ```markdown
> # My GEMINI.md
>
> Welcome to my project!
>
> @./get-started.md
>
> ## Features
>
> @./features/overview.md
> ```
>
> ### Nested imports
>
> The imported files can themselves contain imports, creating a nested structure:
>
> ```markdown
> # main.md
>
> @./header.md @./content.md @./footer.md
> ```
>
> ```markdown
> # header.md
>
> # Project Header
>
> @./shared/title.md
> ```
>
> ## Safety features
>
> ### Circular import detection
>
> The processor automatically detects and prevents circular imports:
>
> ```markdown
> # file-a.md
>
> @./file-b.md
> ```
>
> ```markdown
> # file-b.md
>
> @./file-a.md <!-- This will be detected and prevented -->
> ```
>
> ### File access security
>
> The `validateImportPath` function ensures that imports are only allowed from
> specified directories, preventing access to sensitive files outside the allowed
> scope.
>
> ### Maximum import depth
>
> To prevent infinite recursion, there's a configurable maximum import depth
> (default: 5 levels).
>
> ## Error handling
>
> ### Missing files
>
> If a referenced file doesn't exist, the import will fail gracefully with an
> error comment in the output.
>
> ### File access errors
>
> Permission issues or other file system errors are handled gracefully with
> appropriate error messages.
>
> ## Code region detection
>
> The import processor uses the `marked` library to detect code blocks and inline
> code spans, ensuring that `@` imports inside these regions are properly ignored.
> This provides robust handling of nested code blocks and complex Markdown
> structures.
>
> ## Import tree structure
>
> The processor returns an import tree that shows the hierarchy of imported files,
> similar to Claude's `/memory` feature. This helps users debug problems with
> their GEMINI.md files by showing which files were read and their import
> relationships.
>
> Example tree structure:
>
> ```
> Memory Files
>  L project: GEMINI.md
>             L a.md
>               L b.md
>                 L c.md
>               L d.md
>                 L e.md
>                   L f.md
>             L included.md
> ```
>
> The tree preserves the order that files were imported and shows the complete
> import chain for debugging purposes.
>
> ## Comparison to Claude Code's `/memory` (`claude.md`) approach
>
> Claude Code's `/memory` feature (as seen in `claude.md`) produces a flat, linear
> document by concatenating all included files, always marking file boundaries
> with clear comments and path names. It does not explicitly present the import
> hierarchy, but the LLM receives all file contents and paths, which is sufficient
> for reconstructing the hierarchy if needed.
>
> > [!NOTE] The import tree is mainly for clarity during development and has
> > limited relevance to LLM consumption.
>
> ## API reference
>
> ### `processImports(content, basePath, debugMode?, importState?)`
>
> Processes import statements in GEMINI.md content.
>
> **Parameters:**
>
> - `content` (string): The content to process for imports
> - `basePath` (string): The directory path where the current file is located
> - `debugMode` (boolean, optional): Whether to enable debug logging (default:
>   false)
> - `importState` (ImportState, optional): State tracking for circular import
>   prevention
>
> **Returns:** Promise&lt;ProcessImportsResult&gt; - Object containing processed
> content and import tree
>
> ### `ProcessImportsResult`
>
> ```typescript
> interface ProcessImportsResult {
>   content: string; // The processed content with imports resolved
>   importTree: MemoryFile; // Tree structure showing the import hierarchy
> }
> ```
>
> ### `MemoryFile`
>
> ```typescript
> interface MemoryFile {
>   path: string; // The file path
>   imports?: MemoryFile[]; // Direct imports, in the order they were imported
> }
> ```
>
> ### `validateImportPath(importPath, basePath, allowedDirectories)`
>
> Validates import paths to ensure they are safe and within allowed directories.
>
> **Parameters:**
>
> - `importPath` (string): The import path to validate
> - `basePath` (string): The base directory for resolving relative paths
> - `allowedDirectories` (string[]): Array of allowed directory paths
>
> **Returns:** boolean - Whether the import path is valid
>
> ### `findProjectRoot(startDir)`
>
> Finds the project root by searching for a `.git` directory upwards from the
> given start directory. Implemented as an **async** function using non-blocking
> file system APIs to avoid blocking the Node.js event loop.
>
> **Parameters:**
>
> - `startDir` (string): The directory to start searching from
>
> **Returns:** Promise&lt;string&gt; - The project root directory (or the start
> directory if no `.git` is found)
>
> ## Best Practices
>
> 1. **Use descriptive file names** for imported components
> 2. **Keep imports shallow** - avoid deeply nested import chains
> 3. **Document your structure** - maintain a clear hierarchy of imported files
> 4. **Test your imports** - ensure all referenced files exist and are accessible
> 5. **Use relative paths** when possible for better portability
>
> ## Troubleshooting
>
> ### Common issues
>
> 1. **Import not working**: Check that the file exists and the path is correct
> 2. **Circular import warnings**: Review your import structure for circular
>    references
> 3. **Permission errors**: Ensure the files are readable and within allowed
>    directories
> 4. **Path resolution issues**: Use absolute paths if relative paths aren't
>    resolving correctly
>
> ### Debug mode
>
> Enable debug mode to see detailed logging of the import process:
>
> ```typescript
> const result = await processImports(content, basePath, true);
> ```
