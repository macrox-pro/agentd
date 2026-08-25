---
primary_sources:
  - id: T1-CFG-ADV
    title: "Advanced Configuration"
    url: "https://learn.chatgpt.com/docs/config-file/config-advanced.md"
    section: "Notifications"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Notifications (notify argv)

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Advanced Configuration — Notifications

> ## Notifications
>
> Use `notify` to trigger an external program whenever Codex emits supported events (currently only `agent-turn-complete`). This is handy for desktop toasts, chat webhooks, CI updates, or any side-channel alerting that the built-in TUI notifications don't cover.
>
> ```toml
> notify = ["python3", "/path/to/notify.py"]
> ```
>
> Example `notify.py` (truncated) that reacts to `agent-turn-complete`:
>
> ```python
> #!/usr/bin/env python3
> import json, subprocess, sys
>
> def main() -> int:
>     notification = json.loads(sys.argv[1])
>     if notification.get("type") != "agent-turn-complete":
>         return 0
>     title = f"Codex: {notification.get('last-assistant-message', 'Turn Complete!')}"
>     message = " ".join(notification.get("input-messages", []))
>     subprocess.check_output([
>         "terminal-notifier",
>         "-title", title,
>         "-message", message,
>         "-group", "codex-" + notification.get("thread-id", ""),
>         "-activate", "com.googlecode.iterm2",
>     ])
>     return 0
>
> if __name__ == "__main__":
>     sys.exit(main())
> ```
>
> The script receives a single JSON argument. Common fields include:
>
> - `type` (currently `agent-turn-complete`)
> - `thread-id` (session identifier)
> - `turn-id` (turn identifier)
> - `cwd` (working directory)
> - `input-messages` (user messages that led to the turn)
> - `last-assistant-message` (last assistant message text)
>
> Place the script somewhere on disk and point `notify` to it.
>
> #### `notify` vs `tui.notifications`
>
> - `notify` runs an external program (good for webhooks, desktop notifiers, CI hooks).
> - `tui.notifications` is built in to the TUI and can optionally filter by event type (for example, `agent-turn-complete` and `approval-requested`).
> - `tui.notification_method` controls how the TUI emits terminal notifications (`auto`, `osc9`, or `bel`).
> - `tui.notification_condition` controls whether TUI notifications fire only when
>   the terminal is `unfocused` or `always`.
>
> In `auto` mode, Codex prefers OSC 9 notifications (a terminal escape sequence some terminals interpret as a desktop notification) and falls back to BEL (`\x07`) otherwise.
>
> See [Configuration Reference](https://learn.chatgpt.com/docs/config-file/config-reference) for the exact keys.
