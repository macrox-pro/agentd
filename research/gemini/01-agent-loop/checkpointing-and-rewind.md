---
primary_sources:
  - id: T1-CP
    title: "Checkpointing"
    url: "https://geminicli.com/docs/cli/checkpointing.md"
    section: "Full page"
  - id: T1-RW
    title: "Rewind"
    url: "https://geminicli.com/docs/cli/rewind.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Checkpointing and rewind

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Checkpointing — Full page

> Gemini CLI includes a Checkpointing feature that automatically saves a snapshot
> of your project's state before any file modifications are made by AI-powered
> tools. This lets you safely experiment with and apply code changes, knowing you
> can instantly revert back to the state before the tool was run.
>
> ## How it works
>
> When you approve a tool that modifies the file system (like `write_file` or
> `replace`), the CLI automatically creates a "checkpoint." This checkpoint
> includes:
>
> 1.  **A Git snapshot:** A commit is made in a special, shadow Git repository
>     located in your home directory (`~/.gemini/history/<project_hash>`). This
>     snapshot captures the complete state of your project files at that moment.
>     It does **not** interfere with your own project's Git repository.
> 2.  **Conversation history:** The entire conversation you've had with the agent
>     up to that point is saved.
> 3.  **The tool call:** The specific tool call that was about to be executed is
>     also stored.
>
> If you want to undo the change or simply go back, you can use the `/restore`
> command. Restoring a checkpoint will:
>
> - Revert all files in your project to the state captured in the snapshot.
> - Restore the conversation history in the CLI.
> - Re-propose the original tool call, allowing you to run it again, modify it, or
>   simply ignore it.
>
> All checkpoint data, including the Git snapshot and conversation history, is
> stored locally on your machine. The Git snapshot is stored in the shadow
> repository while the conversation history and tool calls are saved in a JSON
> file in your project's temporary directory, typically located at
> `~/.gemini/tmp/<project_hash>/checkpoints`.
>
> ## Enabling the feature
>
> The Checkpointing feature is disabled by default. To enable it, you need to edit
> your `settings.json` file.
>
> <!-- prettier-ignore -->
> > [!CAUTION]
> > The `--checkpointing` command-line flag was removed in version
> > 0.11.0. Checkpointing can now only be enabled through the `settings.json`
> > configuration file.
>
> Add the following key to your `settings.json`:
>
> ```json
> {
>   "general": {
>     "checkpointing": {
>       "enabled": true
>     }
>   }
> }
> ```
>
> ## Using the `/restore` command
>
> Once enabled, checkpoints are created automatically. To manage them, you use the
> `/restore` command.
>
> ### List available checkpoints
>
> To see a list of all saved checkpoints for the current project, simply run:
>
> ```
> /restore
> ```
>
> The CLI will display a list of available checkpoint files. These file names are
> typically composed of a timestamp, the name of the file being modified, and the
> name of the tool that was about to be run (for example,
> `2025-06-22T10-00-00_000Z-my-file.txt-write_file`).
>
> ### Restore a specific checkpoint
>
> To restore your project to a specific checkpoint, use the checkpoint file from
> the list:
>
> ```
> /restore <checkpoint_file>
> ```
>
> For example:
>
> ```
> /restore 2025-06-22T10-00-00_000Z-my-file.txt-write_file
> ```
>
> After running the command, your files and conversation will be immediately
> restored to the state they were in when the checkpoint was created, and the
> original tool prompt will reappear.

### Source: Rewind — Full page

> The `/rewind` command lets you go back to a previous state in your conversation
> and, optionally, revert any file changes made by the AI during those
> interactions. This is a powerful tool for undoing mistakes, exploring different
> approaches, or simply cleaning up your session history.
>
> ## Usage
>
> To use the rewind feature, simply type `/rewind` into the input prompt and press
> **Enter**.
>
> Alternatively, you can use the keyboard shortcut: **Press `Esc` twice**.
>
> ## Interface
>
> When you trigger a rewind, an interactive list of your previous interactions
> appears.
>
> 1.  **Select interaction:** Use the **Up/Down arrow keys** to navigate through
>     the list. The most recent interactions are at the bottom.
> 2.  **Preview:** As you select an interaction, you'll see a preview of the user
>     prompt and, if applicable, the number of files changed during that step.
> 3.  **Confirm selection:** Press **Enter** on the interaction you want to rewind
>     back to.
> 4.  **Action selection:** After selecting an interaction, you'll be presented
>     with a confirmation dialog with up to three options:
>     - **Rewind conversation and revert code changes:** Reverts both the chat
>       history and the file modifications to the state before the selected
>       interaction.
>     - **Rewind conversation:** Only reverts the chat history. File changes are
>       kept.
>     - **Revert code changes:** Only reverts the file modifications. The chat
>       history is kept.
>     - **Do nothing (esc):** Cancels the rewind operation.
>
> If no code changes were made since the selected point, the options related to
> reverting code changes will be hidden.
>
> ## Key considerations
>
> - **Destructive action:** Rewinding is a destructive action for your current
>   session history and potentially your files. Use it with care.
> - **Agent awareness:** When you rewind the conversation, the AI model loses all
>   memory of the interactions that were removed. If you only revert code changes,
>   you may need to inform the model that the files have changed.
> - **Manual edits:** Rewinding only affects file changes made by the AI's edit
>   tools. It does **not** undo manual edits you've made or changes triggered by
>   the shell tool (`!`).
> - **Compression:** Rewind works across chat compression points by reconstructing
>   the history from stored session data.
