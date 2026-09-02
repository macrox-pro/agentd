---
primary_sources:
  - id: T1-THEMES
    title: "Themes"
    url: "https://opencode.ai/docs/themes.md"
    section: "Full page"
  - id: T1-KEYBINDS
    title: "Keybinds"
    url: "https://opencode.ai/docs/keybinds.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Themes and keybinds

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Themes

> With OpenCode you can select from one of several built-in themes, use a theme that adapts to your terminal theme, or define your own custom theme.
>
> By default, OpenCode uses our own `opencode` theme.
>
> ---
>
> ## Terminal requirements
>
> For themes to display correctly with their full color palette, your terminal must support **truecolor** (24-bit color). Most modern terminals support this by default, but you may need to enable it:
>
> - **Check support**: Run `echo $COLORTERM` - it should output `truecolor` or `24bit`
> - **Enable truecolor**: Set the environment variable `COLORTERM=truecolor` in your shell profile
> - **Terminal compatibility**: Ensure your terminal emulator supports 24-bit color (most modern terminals like iTerm2, Alacritty, Kitty, Windows Terminal, and recent versions of GNOME Terminal do)
>
> Without truecolor support, themes may appear with reduced color accuracy or fall back to the nearest 256-color approximation.
>
> ---
>
> ## Built-in themes
>
> OpenCode comes with several built-in themes.
>
> | Name                   | Description                                                                  |
> | ---------------------- | ---------------------------------------------------------------------------- |
> | `system`               | Adapts to your terminal’s background color                                   |
> | `tokyonight`           | Based on the [Tokyonight](https://github.com/folke/tokyonight.nvim) theme    |
> | `everforest`           | Based on the [Everforest](https://github.com/sainnhe/everforest) theme       |
> | `ayu`                  | Based on the [Ayu](https://github.com/ayu-theme) dark theme                  |
> | `catppuccin`           | Based on the [Catppuccin](https://github.com/catppuccin) theme               |
> | `catppuccin-macchiato` | Based on the [Catppuccin](https://github.com/catppuccin) theme               |
> | `gruvbox`              | Based on the [Gruvbox](https://github.com/morhetz/gruvbox) theme             |
> | `kanagawa`             | Based on the [Kanagawa](https://github.com/rebelot/kanagawa.nvim) theme      |
> | `nord`                 | Based on the [Nord](https://github.com/nordtheme/nord) theme                 |
> | `matrix`               | Hacker-style green on black theme                                            |
> | `one-dark`             | Based on the [Atom One](https://github.com/Th3Whit3Wolf/one-nvim) Dark theme |
>
> And more, we are constantly adding new themes.
>
> ---
>
> ## System theme
>
> The `system` theme is designed to automatically adapt to your terminal's color scheme. Unlike traditional themes that use fixed colors, the _system_ theme:
>
> - **Generates gray scale**: Creates a custom gray scale based on your terminal's background color, ensuring optimal contrast.
> - **Uses ANSI colors**: Leverages standard ANSI colors (0-15) for syntax highlighting and UI elements, which respect your terminal's color palette.
> - **Preserves terminal defaults**: Uses `none` for text and background colors to maintain your terminal's native appearance.
>
> The system theme is for users who:
>
> - Want OpenCode to match their terminal's appearance
> - Use custom terminal color schemes
> - Prefer a consistent look across all terminal applications
>
> ---
>
> ## Using a theme
>
> You can select a theme by bringing up the theme select with the `/theme` command. Or you can specify it in `tui.json`.
>
> ```json title="tui.json" {3}
> {
>   "$schema": "https://opencode.ai/tui.json",
>   "theme": "tokyonight"
> }
> ```
>
> ---
>
> ## Custom themes
>
> OpenCode supports a flexible JSON-based theme system that allows users to create and customize themes easily.
>
> ---
>
> ### Hierarchy
>
> Themes are loaded from multiple directories in the following order where later directories override earlier ones:
>
> 1. **Built-in themes** - These are embedded in the binary
> 2. **User config directory** - Defined in `~/.config/opencode/themes/*.json` or `$XDG_CONFIG_HOME/opencode/themes/*.json`
> 3. **Project root directory** - Defined in the `<project-root>/.opencode/themes/*.json`
> 4. **Current working directory** - Defined in `./.opencode/themes/*.json`
>
> If multiple directories contain a theme with the same name, the theme from the directory with higher priority will be used.
>
> ---
>
> ### Creating a theme
>
> To create a custom theme, create a JSON file in one of the theme directories.
>
> For user-wide themes:
>
> ```bash no-frame
> mkdir -p ~/.config/opencode/themes
> vim ~/.config/opencode/themes/my-theme.json
> ```
>
> And for project-specific themes.
>
> ```bash no-frame
> mkdir -p .opencode/themes
> vim .opencode/themes/my-theme.json
> ```
>
> ---
>
> ### JSON format
>
> Themes use a flexible JSON format with support for:
>
> - **Hex colors**: `"#ffffff"`
> - **ANSI colors**: `3` (0-255)
> - **Color references**: `"primary"` or custom definitions
> - **Dark/light variants**: `{"dark": "#000", "light": "#fff"}`
> - **No color**: `"none"` - Uses the terminal's default color or transparent
>
> ---
>
> ### Color definitions
>
> The `defs` section is optional and it allows you to define reusable colors that can be referenced in the theme.
>
> ---
>
> ### Terminal defaults
>
> The special value `"none"` can be used for any color to inherit the terminal's default color. This is particularly useful for creating themes that blend seamlessly with your terminal's color scheme:
>
> - `"text": "none"` - Uses terminal's default foreground color
> - `"background": "none"` - Uses terminal's default background color
>
> ---
>
> ### Example
>
> Here's an example of a custom theme:
>
> ```json title="my-theme.json"
> {
>   "$schema": "https://opencode.ai/theme.json",
>   "defs": {
>     "nord0": "#2E3440",
>     "nord1": "#3B4252",
>     "nord2": "#434C5E",
>     "nord3": "#4C566A",
>     "nord4": "#D8DEE9",
>     "nord5": "#E5E9F0",
>     "nord6": "#ECEFF4",
>     "nord7": "#8FBCBB",
>     "nord8": "#88C0D0",
>     "nord9": "#81A1C1",
>     "nord10": "#5E81AC",
>     "nord11": "#BF616A",
>     "nord12": "#D08770",
>     "nord13": "#EBCB8B",
>     "nord14": "#A3BE8C",
>     "nord15": "#B48EAD"
>   },
>   "theme": {
>     "primary": {
>       "dark": "nord8",
>       "light": "nord10"
>     },
>     "secondary": {
>       "dark": "nord9",
>       "light": "nord9"
>     },
>     "accent": {
>       "dark": "nord7",
>       "light": "nord7"
>     },
>     "error": {
>       "dark": "nord11",
>       "light": "nord11"
>     },
>     "warning": {
>       "dark": "nord12",
>       "light": "nord12"
>     },
>     "success": {
>       "dark": "nord14",
>       "light": "nord14"
>     },
>     "info": {
>       "dark": "nord8",
>       "light": "nord10"
>     },
>     "text": {
>       "dark": "nord4",
>       "light": "nord0"
>     },
>     "textMuted": {
>       "dark": "nord3",
>       "light": "nord1"
>     },
>     "background": {
>       "dark": "nord0",
>       "light": "nord6"
>     },
>     "backgroundPanel": {
>       "dark": "nord1",
>       "light": "nord5"
>     },
>     "backgroundElement": {
>       "dark": "nord1",
>       "light": "nord4"
>     },
>     "border": {
>       "dark": "nord2",
>       "light": "nord3"
>     },
>     "borderActive": {
>       "dark": "nord3",
>       "light": "nord2"
>     },
>     "borderSubtle": {
>       "dark": "nord2",
>       "light": "nord3"
>     },
>     "diffAdded": {
>       "dark": "nord14",
>       "light": "nord14"
>     },
>     "diffRemoved": {
>       "dark": "nord11",
>       "light": "nord11"
>     },
>     "diffContext": {
>       "dark": "nord3",
>       "light": "nord3"
>     },
>     "diffHunkHeader": {
>       "dark": "nord3",
>       "light": "nord3"
>     },
>     "diffHighlightAdded": {
>       "dark": "nord14",
>       "light": "nord14"
>     },
>     "diffHighlightRemoved": {
>       "dark": "nord11",
>       "light": "nord11"
>     },
>     "diffAddedBg": {
>       "dark": "#3B4252",
>       "light": "#E5E9F0"
>     },
>     "diffRemovedBg": {
>       "dark": "#3B4252",
>       "light": "#E5E9F0"
>     },
>     "diffContextBg": {
>       "dark": "nord1",
>       "light": "nord5"
>     },
>     "diffLineNumber": {
>       "dark": "nord2",
>       "light": "nord4"
>     },
>     "diffAddedLineNumberBg": {
>       "dark": "#3B4252",
>       "light": "#E5E9F0"
>     },
>     "diffRemovedLineNumberBg": {
>       "dark": "#3B4252",
>       "light": "#E5E9F0"
>     },
>     "markdownText": {
>       "dark": "nord4",
>       "light": "nord0"
>     },
>     "markdownHeading": {
>       "dark": "nord8",
>       "light": "nord10"
>     },
>     "markdownLink": {
>       "dark": "nord9",
>       "light": "nord9"
>     },
>     "markdownLinkText": {
>       "dark": "nord7",
>       "light": "nord7"
>     },
>     "markdownCode": {
>       "dark": "nord14",
>       "light": "nord14"
>     },
>     "markdownBlockQuote": {
>       "dark": "nord3",
>       "light": "nord3"
>     },
>     "markdownEmph": {
>       "dark": "nord12",
>       "light": "nord12"
>     },
>     "markdownStrong": {
>       "dark": "nord13",
>       "light": "nord13"
>     },
>     "markdownHorizontalRule": {
>       "dark": "nord3",
>       "light": "nord3"
>     },
>     "markdownListItem": {
>       "dark": "nord8",
>       "light": "nord10"
>     },
>     "markdownListEnumeration": {
>       "dark": "nord7",
>       "light": "nord7"
>     },
>     "markdownImage": {
>       "dark": "nord9",
>       "light": "nord9"
>     },
>     "markdownImageText": {
>       "dark": "nord7",
>       "light": "nord7"
>     },
>     "markdownCodeBlock": {
>       "dark": "nord4",
>       "light": "nord0"
>     },
>     "syntaxComment": {
>       "dark": "nord3",
>       "light": "nord3"
>     },
>     "syntaxKeyword": {
>       "dark": "nord9",
>       "light": "nord9"
>     },
>     "syntaxFunction": {
>       "dark": "nord8",
>       "light": "nord8"
>     },
>     "syntaxVariable": {
>       "dark": "nord7",
>       "light": "nord7"
>     },
>     "syntaxString": {
>       "dark": "nord14",
>       "light": "nord14"
>     },
>     "syntaxNumber": {
>       "dark": "nord15",
>       "light": "nord15"
>     },
>     "syntaxType": {
>       "dark": "nord7",
>       "light": "nord7"
>     },
>     "syntaxOperator": {
>       "dark": "nord9",
>       "light": "nord9"
>     },
>     "syntaxPunctuation": {
>       "dark": "nord4",
>       "light": "nord0"
>     }
>   }
> }
> ```

### Source: OpenCode Keybinds

> OpenCode has a list of keybinds that you can customize through `tui.json`.
>
> ```json title="tui.json"
> {
>   "$schema": "https://opencode.ai/tui.json",
>   "leader_timeout": 2000,
>   "keybinds": {
>     "leader": "ctrl+x",
>     "app_exit": "ctrl+c,ctrl+d,<leader>q",
>     "app_debug": "none",
>     "app_console": "none",
>     "app_heap_snapshot": "none",
>     "app_toggle_animations": "none",
>     "app_toggle_file_context": "none",
>     "app_toggle_diffwrap": "none",
>     "app_toggle_paste_summary": "none",
>     "app_toggle_session_directory_filter": "none",
>     "command_list": "ctrl+p",
>     "help_show": "none",
>     "docs_open": "none",
>
>     "editor_open": "<leader>e",
>     "theme_list": "<leader>t",
>     "theme_switch_mode": "none",
>     "theme_mode_lock": "none",
>     "sidebar_toggle": "<leader>b",
>     "scrollbar_toggle": "none",
>     "status_view": "<leader>s",
>
>     "session_export": "<leader>x",
>     "session_copy": "none",
>     "session_move": "none",
>     "session_new": "<leader>n",
>     "session_list": "<leader>l",
>     "session_timeline": "<leader>g",
>     "session_fork": "none",
>     "session_rename": "ctrl+r",
>     "session_delete": "ctrl+d",
>     "session_share": "none",
>     "session_unshare": "none",
>     "session_interrupt": "escape",
>     "session_compact": "<leader>c",
>     "session_toggle_timestamps": "none",
>     "session_toggle_generic_tool_output": "none",
>     "session_child_first": "<leader>down",
>     "session_child_cycle": "right",
>     "session_child_cycle_reverse": "left",
>     "session_parent": "up",
>
>     "stash_delete": "ctrl+d",
>     "model_provider_list": "ctrl+a",
>     "model_favorite_toggle": "ctrl+f",
>     "model_list": "<leader>m",
>     "model_cycle_recent": "f2",
>     "model_cycle_recent_reverse": "shift+f2",
>     "model_cycle_favorite": "none",
>     "model_cycle_favorite_reverse": "none",
>     "mcp_list": "none",
>     "provider_connect": "none",
>     "console_org_switch": "none",
>     "agent_list": "<leader>a",
>     "agent_cycle": "tab",
>     "agent_cycle_reverse": "shift+tab",
>     "variant_cycle": "ctrl+t",
>     "variant_list": "none",
>
>     "messages_page_up": "pageup,ctrl+alt+b",
>     "messages_page_down": "pagedown,ctrl+alt+f",
>     "messages_line_up": "ctrl+alt+y",
>     "messages_line_down": "ctrl+alt+e",
>     "messages_half_page_up": "ctrl+alt+u",
>     "messages_half_page_down": "ctrl+alt+d",
>     "messages_first": "ctrl+g,home",
>     "messages_last": "ctrl+alt+g,end",
>     "messages_next": "none",
>     "messages_previous": "none",
>     "messages_last_user": "none",
>     "messages_copy": "<leader>y",
>     "messages_undo": "<leader>u",
>     "messages_redo": "<leader>r",
>     "messages_toggle_conceal": "<leader>h",
>     "tool_details": "none",
>     "display_thinking": "none",
>
>     "prompt_submit": "none",
>     "prompt_editor_context_clear": "none",
>     "prompt_skills": "none",
>     "prompt_stash": "none",
>     "prompt_stash_pop": "none",
>     "prompt_stash_list": "none",
>     "workspace_set": "none",
>
>     "input_clear": "ctrl+c",
>     "input_paste": {
>       "key": "ctrl+v",
>       "preventDefault": false
>     },
>     "input_submit": "return",
>     "input_newline": "shift+return,ctrl+return,alt+return,ctrl+j",
>     "input_move_left": "left,ctrl+b",
>     "input_move_right": "right,ctrl+f",
>     "input_move_up": "up",
>     "input_move_down": "down",
>     "input_select_left": "shift+left",
>     "input_select_right": "shift+right",
>     "input_select_up": "shift+up",
>     "input_select_down": "shift+down",
>     "input_line_home": "ctrl+a",
>     "input_line_end": "ctrl+e",
>     "input_select_line_home": "ctrl+shift+a",
>     "input_select_line_end": "ctrl+shift+e",
>     "input_visual_line_home": "alt+a",
>     "input_visual_line_end": "alt+e",
>     "input_select_visual_line_home": "alt+shift+a",
>     "input_select_visual_line_end": "alt+shift+e",
>     "input_buffer_home": "home",
>     "input_buffer_end": "end",
>     "input_select_buffer_home": "shift+home",
>     "input_select_buffer_end": "shift+end",
>     "input_delete_line": "ctrl+shift+d",
>     "input_delete_to_line_end": "ctrl+k",
>     "input_delete_to_line_start": "ctrl+u",
>     "input_backspace": "backspace,shift+backspace",
>     "input_delete": "ctrl+d,delete,shift+delete",
>     "input_undo": "ctrl+-,super+z",
>     "input_redo": "ctrl+.,super+shift+z",
>     "input_word_forward": "alt+f,alt+right,ctrl+right",
>     "input_word_backward": "alt+b,alt+left,ctrl+left",
>     "input_select_word_forward": "alt+shift+f,alt+shift+right",
>     "input_select_word_backward": "alt+shift+b,alt+shift+left",
>     "input_delete_word_forward": "alt+d,alt+delete,ctrl+delete",
>     "input_delete_word_backward": "ctrl+w,ctrl+backspace,alt+backspace",
>     "input_select_all": "super+a",
>     "history_previous": "up",
>     "history_next": "down",
>
>     "dialog.select.prev": "up,ctrl+p",
>     "dialog.select.next": "down,ctrl+n",
>     "dialog.select.page_up": "pageup",
>     "dialog.select.page_down": "pagedown",
>     "dialog.select.home": "home",
>     "dialog.select.end": "end",
>     "dialog.select.submit": "return",
>     "dialog.prompt.submit": "return",
>     "dialog.mcp.toggle": "space",
>     "prompt.autocomplete.prev": "up,ctrl+p",
>     "prompt.autocomplete.next": "down,ctrl+n",
>     "prompt.autocomplete.hide": "escape",
>     "prompt.autocomplete.select": "return",
>     "prompt.autocomplete.complete": "tab",
>     "permission.prompt.fullscreen": "ctrl+f",
>     "plugins.toggle": "space",
>     "dialog.plugins.install": "shift+i",
>
>     "terminal_suspend": "ctrl+z",
>     "terminal_title_toggle": "none",
>     "tips_toggle": "<leader>h",
>     "plugin_manager": "none",
>     "plugin_install": "none",
>
>     "which_key_toggle": "ctrl+alt+k",
>     "which_key_layout_toggle": "ctrl+alt+shift+k",
>     "which_key_pending_toggle": "ctrl+alt+shift+p",
>     "which_key_group_previous": "ctrl+alt+left,ctrl+alt+[",
>     "which_key_group_next": "ctrl+alt+right,ctrl+alt+]",
>     "which_key_scroll_up": "ctrl+alt+up,ctrl+alt+p",
>     "which_key_scroll_down": "ctrl+alt+down,ctrl+alt+n",
>     "which_key_page_up": "ctrl+alt+pageup",
>     "which_key_page_down": "ctrl+alt+pagedown",
>     "which_key_home": "ctrl+alt+home",
>     "which_key_end": "ctrl+alt+end"
>   }
> }
> ```
>
> :::note
> On Windows, the defaults for `input_undo` and `terminal_suspend` are different:
>
> - `input_undo` defaults to `ctrl+z,ctrl+-,super+z` when it is not explicitly configured. The `ctrl+z` binding is added because Windows terminals do not support POSIX suspend.
> - `terminal_suspend` is forced to `none` because native Windows terminals do not support POSIX suspend.
>   :::
>
> ---
>
> ## Leader Key
>
> OpenCode uses a `leader` key for many keybinds. This avoids conflicts in your terminal.
>
> By default, `ctrl+x` is the leader key and many actions require you to first press the leader key and then the shortcut. For example, to start a new session you first press `ctrl+x` and then press `n`.
>
> You don't need to use a leader key for your keybinds but we recommend doing so.
>
> Some navigation keybinds intentionally do not use the leader key by default. For subagent sessions, the defaults are `session_child_first` = `<leader>down`, `session_child_cycle` = `right`, `session_child_cycle_reverse` = `left`, and `session_parent` = `up`.
>
> `leader_timeout` controls how long OpenCode waits for the next key after the leader key. It defaults to `2000` milliseconds.
>
> ---
>
> ## Binding Values
>
> A string can contain one shortcut or multiple comma-separated shortcuts. You can also use an array for multiple shortcuts.
>
> For advanced cases, use an object with `key`, `event`, `preventDefault`, or `fallthrough`.
>
> ```json title="tui.json"
> {
>   "$schema": "https://opencode.ai/tui.json",
>   "keybinds": {
>     "messages_copy": ["<leader>y", "ctrl+shift+c"],
>     "input_paste": {
>       "key": "ctrl+v",
>       "preventDefault": false
>     }
>   }
> }
> ```
>
> ---
>
> ## Disable Keybind
>
> You can disable a keybind by adding the key to `tui.json` with a value of `"none"` or `false`.
>
> ```json title="tui.json"
> {
>   "$schema": "https://opencode.ai/tui.json",
>   "keybinds": {
>     "session_compact": "none"
>   }
> }
> ```
>
> ---
>
> ## Desktop Prompt Shortcuts
>
> The OpenCode desktop app prompt input supports common Readline/Emacs-style shortcuts for editing text. These are built-in and currently not configurable via `opencode.json`.
>
> | Shortcut | Action                                   |
> | -------- | ---------------------------------------- |
> | `ctrl+a` | Move to start of current line            |
> | `ctrl+e` | Move to end of current line              |
> | `ctrl+b` | Move cursor back one character           |
> | `ctrl+f` | Move cursor forward one character        |
> | `alt+b`  | Move cursor back one word                |
> | `alt+f`  | Move cursor forward one word             |
> | `ctrl+d` | Delete character under cursor            |
> | `ctrl+k` | Kill to end of line                      |
> | `ctrl+u` | Kill to start of line                    |
> | `ctrl+w` | Kill previous word                       |
> | `alt+d`  | Kill next word                           |
> | `ctrl+t` | Transpose characters                     |
> | `ctrl+g` | Cancel popovers / abort running response |
>
> ---
>
> ## Shift+Enter
>
> Some terminals don't send modifier keys with Enter by default. You may need to configure your terminal to send `Shift+Enter` as an escape sequence.
>
> ### Windows Terminal
>
> Open your `settings.json` at:
>
> ```
> %LOCALAPPDATA%\Packages\Microsoft.WindowsTerminal_8wekyb3d8bbwe\LocalState\settings.json
> ```
>
> Add this to the root-level `actions` array:
>
> ```json
> "actions": [
>   {
>     "command": {
>       "action": "sendInput",
>       "input": "\u001b[13;2u"
>     },
>     "id": "User.sendInput.ShiftEnterCustom"
>   }
> ]
> ```
>
> Add this to the root-level `keybindings` array:
>
> ```json
> "keybindings": [
>   {
>     "keys": "shift+enter",
>     "id": "User.sendInput.ShiftEnterCustom"
>   }
> ]
> ```
>
> Save the file and restart Windows Terminal or open a new tab.
