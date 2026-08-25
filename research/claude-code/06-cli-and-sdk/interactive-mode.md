---
primary_sources:
  - id: T2-INTERACTIVE
    title: "Interactive mode"
    url: "https://code.claude.com/docs/en/interactive-mode.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Interactive mode

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Interactive mode

> # Interactive mode
>
> > Complete reference for keyboard shortcuts, input modes, and interactive features in Claude Code sessions.
>
> ## Keyboard shortcuts
>
>
>   Keyboard shortcuts may vary by platform and terminal. In [fullscreen rendering](/docs/en/fullscreen), press `?` in the transcript viewer to see available shortcuts there.
>
>   **macOS users**: Option/Alt key shortcuts (`Alt+B`, `Alt+F`, `Alt+D`, `Alt+Y`, `Alt+P`) require configuring Option as Meta in your terminal. See [Enable Option key shortcuts on macOS](/docs/en/terminal-config#enable-option-key-shortcuts-on-macos) for the setting in each terminal.
>
>
> ### General controls
>
> | Shortcut                                                                                     | Description                                                                                                                                                                                                                                                                | Context                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
> | :------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `Ctrl+C`                                                                                     | Interrupt, or clear input                                                                                                                                                                                                                                                  | Interrupts a running operation. If nothing is running, the first press clears the prompt input and a second press exits Claude Code                                                                                                                                                                                                                                                                                                                                                                                              |
> | `Ctrl+X Ctrl+K`                                                                              | Stop all running [background subagents](/docs/en/sub-agents#run-subagents-in-foreground-or-background) in this session, and turn off [artifact auto-replies](/docs/en/artifacts#let-claude-reply-to-comments-on-its-own) for the rest of it. Press twice within 3 seconds to confirm | Subagent control                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
> | `Ctrl+D`                                                                                     | Exit Claude Code session                                                                                                                                                                                                                                                   | The first press shows a confirmation hint and a second press within 800ms exits. When the prompt has text, `Ctrl+D` deletes the character after the cursor instead                                                                                                                                                                                                                                                                                                                                                               |
> | `Ctrl+G` or `Ctrl+X Ctrl+E`                                                                  | Open in default text editor                                                                                                                                                                                                                                                | Edit your prompt or custom response in your default text editor. `Ctrl+X Ctrl+E` is the readline-native binding. Turn on **Show last response in external editor** in `/config` to prepend Claude's previous reply as `#`-commented context above your prompt; Claude Code strips the comment block when you save                                                                                                                                                                                                                |
> | `Ctrl+L`                                                                                     | Redraw screen                                                                                                                                                                                                                                                              | Forces a full terminal redraw, keeping input and conversation history. Use this to recover if the display becomes garbled or partially blank                                                                                                                                                                                                                                                                                                                                                                                     |
> | `Ctrl+O`                                                                                     | Toggle transcript viewer                                                                                                                                                                                                                                                   | Shows detailed tool usage and execution, with a timestamp and the model used on each assistant message. Also expands MCP calls, which collapse to a single line like "Called slack 3 times" by default                                                                                                                                                                                                                                                                                                                           |
> | `Ctrl+R`                                                                                     | Reverse search command history                                                                                                                                                                                                                                             | Search through previous commands interactively                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
> | `Ctrl+V` or `Cmd+V` (iTerm2) or `Alt+V` (Windows and WSL)                                    | Paste image from clipboard                                                                                                                                                                                                                                                 | Inserts an `[Image #N]` chip at the cursor so you can reference it positionally in your prompt. On WSL, both `Ctrl+V` and `Alt+V` are bound; use `Alt+V` if your terminal intercepts `Ctrl+V`                                                                                                                                                                                                                                                                                                                                    |
> | `Ctrl+B`                                                                                     | Background running tasks                                                                                                                                                                                                                                                   | Backgrounds Bash commands and agents. Tmux users press twice                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
> | `Ctrl+T`                                                                                     | Toggle Claude's task checklist                                                                                                                                                                                                                                             | Show or hide [Claude's to-do checklist](#task-list) in the status area. This is not the background-task view; use [`/tasks`](/docs/en/commands) to see running shells and subagents                                                                                                                                                                                                                                                                                                                                                   |
> | `Ctrl+S`                                                                                     | Stash or restore prompt                                                                                                                                                                                                                                                    | With text in the input, stashes it and clears the prompt. Pressed again on an empty prompt, restores the stashed text, cursor position, and pasted content                                                                                                                                                                                                                                                                                                                                                                       |
> | `Ctrl+Z`                                                                                     | Suspend Claude Code                                                                                                                                                                                                                                                        | Unix only. Suspends the process to your shell; run `fg` to resume                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
> | `Left/Right arrows`                                                                          | Cycle through dialog tabs                                                                                                                                                                                                                                                  | Navigate between tabs in permission dialogs and menus                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
> | `Tab`                                                                                        | Accept an autocomplete suggestion, or add a comment to a permission answer                                                                                                                                                                                                 | While autocomplete suggestions are showing in the prompt input, accepts the selected suggestion. On most permission prompts, with **Yes** or **No** focused, opens a comment field on that option, and pressing it again closes the field. See [add a comment when you answer a permission prompt](/docs/en/permissions#add-a-comment-when-you-answer-a-permission-prompt)                                                                                                                                                            |
> | `Up/Down arrows` or `Ctrl+P`/`Ctrl+N`                                                        | Move cursor or navigate command history                                                                                                                                                                                                                                    | When the input spans more than one visual row, whether wrapped or multiline, first moves the cursor within the prompt. Once the cursor is on the first or last visual row, pressing again navigates command history. While you have messages queued, `Up` from the first row instead [takes them back](#take-back-what-you-queued)                                                                                                                                                                                               |
> | `Esc`                                                                                        | Interrupt Claude, or close a dialog                                                                                                                                                                                                                                        | Stop the current response or tool call mid-turn so you can redirect. Claude keeps the work done so far. If you have [messages queued](#queue-messages-while-claude-works), Claude Code sends them next. When a dialog is open, `Esc` closes the dialog. On a permission prompt, `Esc` declines the action, the same as [**No** without a comment](/docs/en/permissions#add-a-comment-when-you-answer-a-permission-prompt)                                                                                                             |
> | `Esc` + `Esc`                                                                                | Clear input draft, or rewind                                                                                                                                                                                                                                               | When the prompt input contains text, double `Esc` clears it and saves the draft to history so `Up` recalls it. When the input is empty, double `Esc` opens the [rewind menu](/docs/en/checkpointing) to restore or summarize code and conversation from a previous point                                                                                                                                                                                                                                                              |
> | `Shift+Tab`, or `Alt+M` on Windows when the Node or Bun runtime doesn't enable VT input mode | Cycle permission modes                                                                                                                                                                                                                                                     | Cycle through `default` (labeled Manual in the mode indicator), `acceptEdits`, `plan`, and, when available, `bypassPermissions` and then `auto`. From `auto`, the first press switches to `default`. See [permission modes](/docs/en/permission-modes). On a file permission prompt, the same key closes an open [comment field](/docs/en/permissions#add-a-comment-when-you-answer-a-permission-prompt). With no field open, it selects the option that allows the action for the rest of the session, when the prompt offers that option |
> | `Option+P` (macOS) or `Alt+P` (Windows/Linux)                                                | Switch model                                                                                                                                                                                                                                                               | Switch models without clearing your prompt                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
> | `Option+T` (macOS) or `Alt+T` (Windows/Linux)                                                | Toggle extended thinking                                                                                                                                                                                                                                                   | Enable or disable extended thinking mode. Has no effect on Fable 5, which always uses extended thinking. Works on macOS without configuring Option as Meta                                                                                                                                                                                                                                                                                                                                                                       |
> | `Option+O` (macOS) or `Alt+O` (Windows/Linux)                                                | Toggle fast mode                                                                                                                                                                                                                                                           | Enable or disable [fast mode](/docs/en/fast-mode)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
>
> ### Text editing
>
> | Shortcut                   | Description                          | Context                                                                                                                                                                                                                                                                        |
> | :------------------------- | :----------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `Ctrl+A`                   | Move cursor to start of current line | In multiline input, moves to the start of the current logical line                                                                                                                                                                                                             |
> | `Ctrl+E`                   | Move cursor to end of current line   | In multiline input, moves to the end of the current logical line                                                                                                                                                                                                               |
> | `Ctrl+K`                   | Delete to end of line                | Stores deleted text for pasting                                                                                                                                                                                                                                                |
> | `Ctrl+U`                   | Delete from cursor to line start     | Stores deleted text for pasting. Repeat to clear across lines in multiline input. On macOS, terminal emulators including iTerm2 and Terminal.app map `Cmd+Backspace` to this shortcut                                                                                          |
> | `Ctrl+W`                   | Delete previous word                 | Stores deleted text for pasting. On macOS, `Option+Delete` deletes the previous word, and on Windows, `Ctrl+Backspace` does. To make `Ctrl+W` delete back to the previous whitespace instead, [set `keybindingFlavor` to `"readline"`](#make-ctrl-w-delete-back-to-whitespace) |
> | `Ctrl+Y`                   | Paste deleted text                   | Paste text deleted with `Ctrl+K`, `Ctrl+U`, `Ctrl+W`, or, under [`keybindingFlavor: "readline"`](#make-ctrl-w-delete-back-to-whitespace), `Alt+D`                                                                                                                              |
> | `Alt+Y` (after `Ctrl+Y`)   | Cycle paste history                  | After pasting, cycle through previously deleted text. Requires [Option as Meta](#keyboard-shortcuts) on macOS                                                                                                                                                                  |
> | `Alt+B`                    | Move cursor back one word            | Word navigation. Requires [Option as Meta](#keyboard-shortcuts) on macOS                                                                                                                                                                                                       |
> | `Alt+F`                    | Move cursor forward one word         | Moves to the start of the next word, or to the end of the current word when [`keybindingFlavor` is `"readline"`](#make-ctrl-w-delete-back-to-whitespace). Requires [Option as Meta](#keyboard-shortcuts) on macOS                                                              |
> | `Alt+D`                    | Delete next word                     | Deletes through the space after the word, or to the end of the word when [`keybindingFlavor` is `"readline"`](#make-ctrl-w-delete-back-to-whitespace). Requires [Option as Meta](#keyboard-shortcuts) on macOS                                                                 |
> | `Ctrl+_` or `Ctrl+Shift+-` | Undo last input edit                 | Restores the previous input text and cursor position                                                                                                                                                                                                                           |
>
>
>   Make editing keys follow readline conventions
>
>
> Set [`keybindingFlavor`](/docs/en/settings-reference#keybindingflavor) to `"readline"` to make the prompt's editing keys follow GNU readline conventions, as in Bash. The default value is `"classic"`. Requires Claude Code v2.1.238 or later. Under `"readline"`:
>
> * `Ctrl+W` deletes back to the previous whitespace.
> * A word is a run of letters and digits, so punctuation such as `_`, `.`, and `/` separates words. `Alt+F` and `Alt+D` stop at the end of the current word, and `Ctrl+Y` can paste back text that `Alt+D` deleted. Before v2.1.239, Claude Code applied `"readline"` only to `Ctrl+W`.
>
> Add the setting to `~/.claude/settings.json`:
>
> ```json
> {
>   "keybindingFlavor": "readline"
> }
> ```
>
> To confirm, type `fix the bug in src/utils/foo.ts` in the prompt and press `Ctrl+W`. Claude Code removes `src/utils/foo.ts`. Under `"classic"` it removes only `foo.ts`.
>
> This setting is separate from the [keybindings configuration file](/docs/en/keybindings): the word-editing commands aren't actions there, so you can't remap them in `keybindings.json`.
>
> ### Theme and display
>
> | Shortcut | Description                                | Context                                                                                                      |
> | :------- | :----------------------------------------- | :----------------------------------------------------------------------------------------------------------- |
> | `Ctrl+T` | Toggle syntax highlighting for code blocks | Only works inside the `/theme` picker menu. Controls whether code in Claude's responses uses syntax coloring |
>
> ### Multiline input
>
> | Method           | Shortcut       | Context                                                                                                                                                                            |
> | :--------------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | Quick escape     | `\` + `Enter`  | Works in all terminals                                                                                                                                                             |
> | Option key       | `Option+Enter` | After enabling [Option as Meta](/docs/en/terminal-config#enable-option-key-shortcuts-on-macos) on macOS                                                                                 |
> | Shift+Enter      | `Shift+Enter`  | Native in iTerm2, WezTerm, Ghostty, Kitty, Warp, Apple Terminal, Windows Terminal. For other terminals, see [Enter multiline prompts](/docs/en/terminal-config#enter-multiline-prompts) |
> | Control sequence | `Ctrl+J`       | Works in any terminal without configuration                                                                                                                                        |
> | Paste mode       | Paste directly | For code blocks, logs                                                                                                                                                              |
>
> ### Quick commands
>
> | Shortcut           | Description                    | Notes                                                                                                                                                                                                                                                                                                                                            |
> | :----------------- | :----------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `/` at start       | Command or skill               | See [commands](#commands) and [skills](/docs/en/skills)                                                                                                                                                                                                                                                                                               |
> | `!` at start       | Shell mode                     | Run a command directly, add its output to the session, and have Claude respond to it                                                                                                                                                                                                                                                             |
> | `@`                | File path mention              | Trigger file path autocomplete. In sessions with [cross-session messaging](/docs/en/cross-session-messaging#message-another-session), when you type at least one letter after the `@`, Claude Code also suggests your other live sessions on this machine, so you can tell Claude to message the one you pick. Requires Claude Code v2.1.232 or later |
> | `:`                | Emoji shortcode                | Type a full `:name:` to insert the emoji, or two or more characters for suggestions. See [Emoji shortcodes](#emoji-shortcodes). Requires Claude Code v2.1.217 or later                                                                                                                                                                           |
> | `?` on empty input | Toggle the shortcut help panel | Typing `?` when the input already contains text inserts the character                                                                                                                                                                                                                                                                            |
>
> ### Transcript viewer
>
> When the transcript viewer is open (toggled with `Ctrl+O`), these shortcuts are available. Run `/tui` with no argument to check which renderer is active. `Ctrl+E` can be rebound via [`transcript:toggleShowAll`](/docs/en/keybindings).
>
> | Shortcut             | Description                                                                                                                                                                                                           |
> | :------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `?`                  | Toggle the keyboard shortcut help panel. Requires [fullscreen rendering](/docs/en/fullscreen)                                                                                                                              |
> | `{` / `}`            | Jump to the previous or next user prompt, like vim paragraph motion. Requires [fullscreen rendering](/docs/en/fullscreen)                                                                                                  |
> | `Ctrl+E`             | Toggle show all content. Available in the classic renderer only, not in [fullscreen rendering](/docs/en/fullscreen)                                                                                                        |
> | `[`                  | Write the full conversation to your terminal's native scrollback so `Cmd+F`, tmux copy mode, and other native tools can search it. Requires [fullscreen rendering](/docs/en/fullscreen#search-and-review-the-conversation) |
> | `v`                  | Write the conversation to a temporary file and open it in `$VISUAL` or `$EDITOR`. Requires [fullscreen rendering](/docs/en/fullscreen)                                                                                     |
> | `q`, `Ctrl+C`, `Esc` | Exit transcript view. All three can be rebound via [`transcript:exit`](/docs/en/keybindings)                                                                                                                               |
>
> ### Voice input
>
> | Shortcut            | Description     | Notes                                                                                                                                                                            |
> | :------------------ | :-------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | Hold or tap `Space` | Voice dictation | Requires [voice dictation](/docs/en/voice-dictation) to be enabled. Hold to record, or run `/voice tap` for tap-to-toggle. [Rebindable](/docs/en/voice-dictation#rebind-the-dictation-key) |
>
> ## Commands
>
> Type `/` in Claude Code to see the commands available to you, or type `/` followed by any letters to filter. The `/` menu lists built-in commands, bundled and user-authored [skills](/docs/en/skills), and commands contributed by [plugins](/docs/en/plugins) and [MCP servers](/docs/en/mcp#use-mcp-prompts-as-commands). Not all built-in commands are visible to every user since some depend on your platform or plan, and [a few available commands are hidden from the menu by design](/docs/en/commands#how-the-command-menu-matches-what-you-type) and run when you type their full name.
>
> In [fullscreen rendering](/docs/en/fullscreen#use-the-mouse), the `/` command and `@` file suggestion lists also respond to the mouse: hovering highlights a row and clicking accepts it.
>
> See the [commands reference](/docs/en/commands) for the full list of commands included in Claude Code.
>
> ## Vim editor mode
>
> Enable vim-style editing via `/config` → Editor mode.
>
> Claude Code keeps your vim mode and cursor position when you toggle the [transcript viewer](#transcript-viewer) with `Ctrl+O` or open and close a panel such as `/config`. If you leave the prompt in NORMAL mode, it's still in NORMAL mode when you return, with the cursor where you left it.
>
> ### Mode switching
>
> | Command | Action                                | From mode      |
> | :------ | :------------------------------------ | :------------- |
> | `Esc`   | Enter NORMAL mode                     | INSERT, VISUAL |
> | `i`     | Insert before cursor                  | NORMAL         |
> | `I`     | Insert at beginning of line           | NORMAL         |
> | `a`     | Insert after cursor                   | NORMAL         |
> | `A`     | Insert at end of line                 | NORMAL         |
> | `o`     | Open line below                       | NORMAL         |
> | `O`     | Open line above                       | NORMAL         |
> | `v`     | Start character-wise visual selection | NORMAL         |
> | `V`     | Start line-wise visual selection      | NORMAL         |
>
> ### Remap INSERT-mode key sequences
>
> The [`vimInsertModeRemaps`](/docs/en/settings-reference#viminsertmoderemaps) setting maps a two-key INSERT-mode sequence to Escape, so a mapping like `jj` returns you to NORMAL mode. Requires Claude Code v2.1.208 or later.
>
> The following `~/.claude/settings.json` example turns on vim mode and maps `jj` to Escape:
>
> ```json
> {
>   "editorMode": "vim",
>   "vimInsertModeRemaps": { "jj": "<Esc>" }
> }
> ```
>
> Each key is exactly two printable characters typed in sequence, and `"<Esc>"` is the only supported target. Entries with a different length or target are ignored.
>
> Typing the first character of a sequence inserts it normally. Pressing the second character within one second removes that pending character and switches to NORMAL mode, leaving neither character in your input. After the one-second window, or if a different key follows, both characters stay as literal text, so you can still type a word containing the sequence by pausing between the two keys.
>
> Claude Code reads this setting from your user settings file, the `--settings` flag, and [managed settings](/docs/en/managed-settings) only. Entries in a project's `.claude/settings.json` or `.claude/settings.local.json` are ignored, so a checked-out repository can't remap your keystrokes.
>
> ### Navigation (NORMAL mode)
>
> | Command         | Action                                                                                                                                              |
> | :-------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `h`/`j`/`k`/`l` | Move left/down/up/right                                                                                                                             |
> | `Space`         | Move right                                                                                                                                          |
> | `w`             | Next word                                                                                                                                           |
> | `e`             | End of word                                                                                                                                         |
> | `b`             | Previous word                                                                                                                                       |
> | `0`             | Beginning of line                                                                                                                                   |
> | `$`             | End of line                                                                                                                                         |
> | `^`             | First non-blank character                                                                                                                           |
> | `gg`            | Beginning of input                                                                                                                                  |
> | `G`             | End of input                                                                                                                                        |
> | `f{char}`       | Jump to next occurrence of character                                                                                                                |
> | `F{char}`       | Jump to previous occurrence of character                                                                                                            |
> | `t{char}`       | Jump to just before next occurrence of character                                                                                                    |
> | `T{char}`       | Jump to just after previous occurrence of character                                                                                                 |
> | `;`             | Repeat last f/F/t/T motion                                                                                                                          |
> | `,`             | Repeat last f/F/t/T motion in reverse                                                                                                               |
> | `/`             | Open reverse history search, same as `Ctrl+R`. The empty search prompt shows a hint: press `Esc` then `i` then `/` to open the command menu instead |
>
>
>   In vim NORMAL mode, if the cursor is at the beginning or end of input and can't move further, `j`/`k` and `↑`/`↓` navigate command history instead. `←` on an empty prompt opens [agent view](/docs/en/agent-view) from NORMAL mode as well as INSERT; before v2.1.219, `←` on an empty prompt did nothing in NORMAL mode.
>
>
> ### Editing (NORMAL mode)
>
> | Command        | Action                                                                                                                    |
> | :------------- | :------------------------------------------------------------------------------------------------------------------------ |
> | `x`            | Delete character                                                                                                          |
> | `dd`           | Delete line                                                                                                               |
> | `D`            | Delete to end of line                                                                                                     |
> | `dw`/`de`/`db` | Delete word/to end/back                                                                                                   |
> | `cc`           | Change line                                                                                                               |
> | `C`            | Change to end of line                                                                                                     |
> | `cw`/`ce`/`cb` | Change word/to end/back                                                                                                   |
> | `s`            | Substitute character: delete the character under the cursor and enter INSERT mode. Requires Claude Code v2.1.211 or later |
> | `S`            | Substitute line: clear the line and enter INSERT mode. Requires Claude Code v2.1.211 or later                             |
> | `yy`/`Y`       | Yank (copy) line                                                                                                          |
> | `yw`/`ye`/`yb` | Yank word/to end/back                                                                                                     |
> | `p`            | Paste after cursor                                                                                                        |
> | `P`            | Paste before cursor                                                                                                       |
> | `>>`           | Indent line                                                                                                               |
> | `<<`           | Dedent line                                                                                                               |
> | `J`            | Join lines                                                                                                                |
> | `u`            | Undo                                                                                                                      |
> | `.`            | Repeat last change                                                                                                        |
>
> ### Text objects (NORMAL mode)
>
> Text objects work with operators like `d`, `c`, and `y`:
>
> | Command   | Action                                   |
> | :-------- | :--------------------------------------- |
> | `iw`/`aw` | Inner/around word                        |
> | `iW`/`aW` | Inner/around WORD (whitespace-delimited) |
> | `i"`/`a"` | Inner/around double quotes               |
> | `i'`/`a'` | Inner/around single quotes               |
> | `i(`/`a(` | Inner/around parentheses                 |
> | `i[`/`a[` | Inner/around brackets                    |
> | `i{`/`a{` | Inner/around braces                      |
>
> ### Visual mode
>
> Press `v` for character-wise selection or `V` for line-wise selection. Motions extend the selection, and operators act on it directly.
>
> | Command          | Action                                               |
> | :--------------- | :--------------------------------------------------- |
> | `d`/`x`          | Delete selection                                     |
> | `y`              | Yank selection                                       |
> | `c`/`s`          | Change selection                                     |
> | `p`              | Replace selection with register contents             |
> | `r{char}`        | Replace every selected character with `{char}`       |
> | `~`/`u`/`U`      | Toggle, lowercase, or uppercase selection            |
> | `>`/`<`          | Indent or dedent selected lines                      |
> | `J`              | Join selected lines                                  |
> | `o`              | Swap cursor and anchor                               |
> | `iw`/`aw`/`i"`/… | Select a text object                                 |
> | `v`/`V`          | Toggle between character-wise and line-wise, or exit |
>
> Block-wise visual mode with `Ctrl+V` is not supported.
>
> ## Command history
>
> Claude Code keeps a history of the prompts you type, and Up-arrow recall reaches prompts from past sessions of the same project:
>
> * Input history is stored per working directory
> * Running `/clear` starts a new session: recall then lists the new session's prompts first, with earlier sessions' prompts after them. The previous session's conversation is preserved and can be resumed.
> * Submitting the same prompt twice in a row records one history entry, so pressing Up steps to the previous distinct prompt
> * When you recall a prompt that included pasted text, Claude Code sends the full pasted content again when you resubmit. If the content has since been [cleaned up](/docs/en/claude-directory#cleaned-up-automatically), Claude Code doesn't send the literal `[Pasted text #N]` string; see [Paste large content](/docs/en/terminal-config#paste-large-content) for what happens to the prompt
> * History expansion with `!` is disabled by default
>
> ### Reverse search with Ctrl+R
>
> Press `Ctrl+R` to interactively search through your command history. In [fullscreen rendering](/docs/en/fullscreen), `Ctrl+R` opens a search dialog instead: type to filter, press `Up` and `Down` to move through matches, and press `Ctrl+S` to cycle the scope through this session, this project, and all projects. Press `Enter` or `Tab` to place a match in the prompt input, or `Esc` to cancel. The steps below describe the classic renderer's inline search:
>
> 1. **Start search**: press `Ctrl+R` to activate reverse history search
> 2. **Type query**: enter text to search for in previous commands. The search term is highlighted in matching results
> 3. **Navigate matches**: press `Ctrl+R` again to cycle through older matches
> 4. **Search scope**: the inline search always searches prompts from all projects
> 5. **Accept match**:
>    * Press `Tab` or `Esc` to accept the current match and continue editing
>    * Press `Enter` to accept and execute the command immediately
> 6. **Cancel search**:
>    * Press `Ctrl+C` to cancel and restore your original input
>    * Press `Backspace` on empty search to cancel
>
> The inline search scans your full prompt history, newest first, with duplicates collapsed to the newest occurrence. The fullscreen dialog searches your whole prompt history in the selected scope, newest first, with duplicates collapsed to the newest occurrence: the most recent prompts appear immediately, and matches from older prompts fill in as Claude Code loads the rest. Matching prompts display with the search term highlighted, so you can find and reuse previous inputs.
>
> Accepting a match or canceling the search takes effect immediately, even while Claude Code is still loading the history.
>
> ## Background Bash commands
>
> Claude Code supports running Bash commands in the background, allowing you to continue working while long-running processes execute.
>
> ### How backgrounding works
>
> When Claude Code runs a command in the background, it runs the command asynchronously and immediately returns a background task ID. Claude Code can respond to new prompts while the command continues executing in the background.
>
> To run commands in the background, you can either:
>
> * Prompt Claude Code to run a command in the background
> * Press `Ctrl+B` to move a regular Bash tool invocation to the background. Tmux users must press `Ctrl+B` twice due to tmux's prefix key.
>
> **Key features:**
>
> * Output is written to a file and Claude can retrieve it using the Read tool
> * Background tasks have unique IDs for tracking and output retrieval
> * Background tasks are automatically cleaned up when Claude Code exits. If you background the session instead of exiting it, Claude Code hands them to the background session, where they keep running. See [background a running session](/docs/en/agent-view#from-inside-a-session)
> * Background tasks are automatically terminated if output exceeds 5GB, with a note in stderr explaining why
> * On macOS and Linux, Claude Code terminates running background tasks when the operating system signals memory pressure, provided the session has been idle for at least 30 minutes a
