---
primary_sources:
  - id: T1-TROUBLESHOOTING
    title: "Troubleshooting"
    url: "https://opencode.ai/docs/troubleshooting.md"
    section: "Storage"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Storage and runtime paths

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Troubleshooting — Logs and Storage

> ## Logs
>
> Log files are written to:
>
> - **macOS/Linux**: `~/.local/share/opencode/log/`
> - **Windows**: Press `WIN+R` and paste `%USERPROFILE%\.local\share\opencode\log`
>
> Log files are named with timestamps (e.g., `2025-01-09T123456.log`) and the most recent 10 log files are kept.
>
> You can set the log level with the `--log-level` command-line option to get more detailed debug information. For example, `opencode --log-level DEBUG`.
>
> ---
>
> ---
>
> ## Storage
>
> opencode stores session data and other application data on disk at:
>
> - **macOS/Linux**: `~/.local/share/opencode/`
> - **Windows**: Press `WIN+R` and paste `%USERPROFILE%\.local\share\opencode`
>
> This directory contains:
>
> - `auth.json` - Authentication data like API keys, OAuth tokens
> - `log/` - Application logs
> - `project/` - Project-specific data like session and message data
>   - If the project is within a Git repo, it is stored in `./<project-slug>/storage/`
>   - If it is not a Git repo, it is stored in `./global/storage/`
>
> ---
