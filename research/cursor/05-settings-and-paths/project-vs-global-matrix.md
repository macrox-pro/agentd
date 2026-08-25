---
primary_sources:
  - id: T1-HOOKS
    title: "Configuration"
    url: "https://cursor.com/docs/hooks.md"
    section: "Configuration"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Project vs global matrix

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).
------|---------|----------------------------|-------------------|
| Hooks | `.cursor/hooks.json`, `.cursor/hooks/*` | `~/.cursor/hooks.json`, `~/.cursor/hooks/*` | Dashboard team hooks; Enterprise: `/Library/Application Support/Cursor/hooks.json` (macOS), `/etc/cursor/hooks.json` (Linux), `C:\ProgramData\Cursor\hooks.json` (Windows) |
| MCP | `.cursor/mcp.json` | `~/.cursor/mcp.json` | Team MCP via dashboard / marketplace |
| Rules | `.cursor/rules/*.mdc` | User Rules: Customize → Rules (UI only; no on-disk path in official docs) | Team Rules: dashboard |
| AGENTS.md | repo root + nested subdirs | — | — |
| Skills | `.cursor/skills/`, `.agents/skills/` | `~/.cursor/skills/` | via plugins |
| Subagents | `.cursor/agents/` | `~/.cursor/agents/` | — |
| CLI permissions | `.cursor/cli.json` | `~/.cursor/cli-config.json` | — |
| Cloud environment | `.cursor/environment.json` | personal saved env (dashboard) | team saved env |
| Bugbot rules | `.cursor/BUGBOT.md` | personal Bugbot settings (dashboard) | team rules |
| Ignore | `.cursorignore`, `.cursorindexingignore` | global ignore patterns in IDE user settings (UI) | — |
| Plugins | workspace scope in Customize | user scope | team marketplaces |

### Source: Cursor Hooks — Configuration sources and precedence

> Priority: Enterprise → Team → Project → User (Claude third-party hooks lower).

### Source: Cursor Rules — Precedence

> Team Rules → Project Rules → User Rules.

### Source: Ignore file — Global ignore files

> Set ignore patterns for all projects in user settings to exclude sensitive files without per-project configuration.
