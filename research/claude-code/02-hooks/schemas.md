---
primary_sources:
  - id: T1-HOOKS
    title: "Hooks reference"
    url: "https://code.claude.com/docs/en/hooks.md"
    section: "JSON output; Hook events"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook schemas

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks reference — JSON schemas pointer

> Per-event JSON input/output examples are extracted under `schemas/hooks/` (31 `*.input.example.json` files).
>
> No official generated JSON Schema repository was found in Anthropic open-source repos (unlike Codex `hooks/schema/generated`). See [anthropics/claude-code hook-development SKILL.md](https://github.com/anthropics/claude-code/blob/main/plugins/plugin-dev/skills/hook-development/SKILL.md) for plugin hook format notes.
>
> ### JSON output
>
> Exit codes only let you block or stay silent, but JSON output gives you finer-grained control. Instead of exiting with code 2 to block, exit 0 and print a JSON object to stdout. Claude Code reads specific fields from that JSON to control behavior, including [decision control](#decision-control) for blocking, allowing, or escalating to the user.
>
>
>   Choose one approach per hook: either use exit codes alone for signaling, or exit 0 and print JSON for structured control. If you mix them, exit 2 keeps its [blocking effect](#exit-code-2-behavior-per-event), and Claude Code still reads the JSON fields, with the one elicitation exception noted under [Exit code 2](#exit-code-2).
>
>
> Your hook's stdout must contain only the JSON object. If your shell profile prints text on startup, it can interfere with JSON parsing. See [Hook JSON has no effect](/docs/en/hooks-guide#hook-json-has-no-effect) in the troubleshooting guide.
>
> Hook output strings, including `additionalContext`, `systemMessage`, and plain stdout, are capped at 10,000 characters. Output that exceeds this limit is saved to a file and replaced with a preview and file path, the same way a large valid Bash result is handled under [Output limits](/docs/en/tools-reference#output-limits).
>
> The JSON object supports three kinds of fields:
>
> * **Universal fields** like `continue` are listed in the table below. Every event accepts them, but some events discard them or deliver `systemMessage` somewhere other than the transcript. Each event's section says so. `terminalSequence` works on those events too, with the exceptions listed under [Emit terminal notifications](#emit-terminal-notifications).
> * **Top-level `decision` and `reason`** are used by some events to block or provide feedback.
> * **`hookSpecificOutput`** is a nested object for events that need richer control. It requires a `hookEventName` field set to the event name.
>
> | Field              | Default | Description                                                                                                                                                                                                                                                                                                                          |
> | :----------------- | :------ | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `continue`         | `true`  | If `false`, Claude stops processing entirely after the hook runs. Takes precedence over any event-specific decision fields                                                                                                                                                                                                           |
> | `stopReason`       | none    | Message shown to the user when `continue` is `false`. Not shown to Claude                                                                                                                                                                                                                                                            |
> | `suppressOutput`   | `false` | Has no effect: Claude Code accepts the field but doesn't act on it. A successful hook's stdout is never shown in the transcript and is recorded in the debug log                                                                                                                                                                     |
> | `systemMessage`    | none    | Warning message shown to the user. In [Agent SDK](/docs/en/agent-sdk/overview) and [`--output-format stream-json`](/docs/en/headless) output, it can arrive as an [`SDKInformationalMessage`](/docs/en/agent-sdk/typescript#sdkinformationalmessage)                                                                                                |
> | `terminalSequence` | none    | A terminal escape sequence for Claude Code to emit on your behalf, such as a desktop notification, window title, or bell. Restricted to OSC `0`/`1`/`2`/`9`/`99`/`777` and BEL. If the value contains anything outside the allowlist, the field is ignored. Use this instead of writing to `/dev/tty`, which is unavailable to hooks |
>
> To stop Claude entirely:
>
> ```json
> { "continue": false, "stopReason": "Build failed, fix errors before continuing" }
> ```
>
> For `PreToolUse` and `PostToolUse` hooks, the stop applies even when the tool call fails or completes while Claude is still streaming a response.
>
> #### Emit terminal notifications
>
> Hooks run without a controlling terminal, so writing escape sequences directly to `/dev/tty` fails. Instead, return the escape sequence in the `terminalSequence` field and Claude Code emits it for you through its own terminal write path. This is race-free, works inside tmux and GNU screen, and works on Windows where there is no `/dev/tty`.
>
> The field accepts a string of one or more allowlisted escape sequences:
>
> * OSC `0`, `1`, `2`: window and icon titles
> * OSC `9`: iTerm2, ConEmu, Windows Terminal, and WezTerm notifications, including `9;4` taskbar progress
> * OSC `99`: Kitty notifications
> * OSC `777`: urxvt, Ghostty, and Warp notifications
> * Bare BEL
>
> Sequences may be terminated with BEL or with ST. Anything outside the allowlist, including CSI cursor and color sequences, OSC palette sequences, OSC 8 hyperlinks, OSC 52 clipboard writes, and OSC 1337, is rejected and the field is ignored.
>
> Claude Code writes the sequence itself when it processes your hook's output, so the field works on events that discard `systemMessage` and `continue`, such as `Notification` and `StopFailure`. It has two limits:
>
> * Claude Code writes the sequence only in an interactive session, and only while its interface is on screen. In non-interactive mode
