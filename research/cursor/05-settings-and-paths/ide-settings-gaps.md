---
primary_sources:
  - id: T1-IGNORE
    title: "Gaps"
    url: "https://cursor.com/docs/reference/ignore-file.md"
    section: "Gaps"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# IDE settings gaps

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).
# IDE settings gaps

> **Applicability:** Documents configuration surfaces not fully specified on disk in official `/docs/`.

### Source: Ignore file — Global ignore files

> Set ignore patterns for all projects in user settings to exclude sensitive files without per-project configuration. The global ignore list is empty by default.

### Source: Cursor Rules — User Rules

> User Rules are global preferences defined in **Customize → Rules** that apply across all projects. They are used by Agent (Chat).

### Source: Cursor Rules — Team Rules

> Team administrators can create and manage rules directly from the Cursor dashboard.

**Not documented as on-disk paths in official `/docs/` (as of snapshot):**

- IDE `settings.json` (VS Code/Cursor editor settings) — commonly `~/Library/Application Support/Cursor/User/settings.json` (macOS) per ecosystem convention; no dedicated Cursor `/docs/` page found.
- User Rules storage location on filesystem — UI/dashboard only per `rules.md`.

### Source: Ignore file — Full page (excerpt)

> # Ignore file
>
> Cursor reads and indexes your project's codebase to power its features. Control which directories and files Cursor can access using a `.cursorignore` file in your root directory.
>
> Cursor blocks access to files listed in `.cursorignore` from:
>
> - Code accessible by [Agent](https://cursor.com/docs/agent/overview.md), Tab, and Inline Edit
> - Code accessible via [@ mention references](https://cursor.com/docs/agent/prompting.md)
>
> The terminal and MCP server tools used by Agent cannot block access to code
> governed by `.cursorignore`
>
> ## Why ignore files?
>
> **Security**: Restrict access to API keys, credentials, and secrets. While Cursor blocks ignored files, complete protection isn't guaranteed due to LLM unpredictability.
>
> **Performance**: In large codebases or monorepos, exclude irrelevant portions for faster indexing and more accurate file discovery.
>
> ## Configuring `.cursorignore`
>
> Create a `.cursorignore` file in your root directory using `.gitignore` syntax.
>
> ### Pattern syntax
>
> - `*` matches any characters except `/`
> - `**` matches any characters including `/`
> - `?` matches a single character
> - `!` negates a pattern (un-ignores a previously ignored path)
> - Lines starting with `#` are comments
> - Trailing spaces are ignored unless escaped with `\`
>
> ### Pattern examples
>
> ```sh
> config.json      # Specific file
> dist/           # Directory
> *.log           # File extension
> **/logs         # Nested directories
> !app/           # Exclude from ignore (negate)
> ```
>
> ### Hierarchical ignore
>
> Enable `Cursor Settings` > `Features` > `Editor` > `Hierarchical Cursor Ignore` to search parent directories for `.cursorignore` files.
>
> Starting in Cursor 3.11, this setting moves to `Cursor Settings` > `Indexing` > `Ignore Files` > `Hierarchical Cursor Ignore`.
>
> ## Global ignore files
>
> Set ignore patterns for all projects in user settings to exclude sensitive files without per-project configuration. The global ignore list is empty by default.
>
> Common patterns to add:
>
> - Environ