---
primary_sources:
  - id: T1-CONFIG
    title: "Config"
    url: "https://opencode.ai/docs/config.md"
    section: "Variables"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Config variables substitution

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Config — Variables

> ## Variables
>
> You can use variable substitution in your config files to reference environment variables and file contents.
>
> ---
>
> ### Env vars
>
> Use `{env:VARIABLE_NAME}` to substitute environment variables:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "model": "{env:OPENCODE_MODEL}",
>   "provider": {
>     "anthropic": {
>       "models": {},
>       "options": {
>         "apiKey": "{env:ANTHROPIC_API_KEY}"
>       }
>     }
>   }
> }
> ```
>
> If the environment variable is not set, it will be replaced with an empty string.
>
> ---
>
> ### Files
>
> Use `{file:path/to/file}` to substitute the contents of a file:
>
> ```json title="opencode.json"
> {
>   "$schema": "https://opencode.ai/config.json",
>   "instructions": ["./custom-instructions.md"],
>   "provider": {
>     "openai": {
>       "options": {
>         "apiKey": "{file:~/.secrets/openai-key}"
>       }
>     }
>   }
> }
> ```
>
> File paths can be:
>
> - Relative to the config file directory
> - Or absolute paths starting with `/` or `~`
>
> These are useful for:
>
> - Keeping sensitive data like API keys in separate files.
> - Including large instruction files without cluttering your config.
> - Sharing common configuration snippets across multiple config files.
