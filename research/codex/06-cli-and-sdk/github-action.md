---
primary_sources:
  - id: T2-GHA
    title: "Codex GitHub Action"
    url: "https://learn.chatgpt.com/docs/github-action.md"
    section: "full page"
also_cited_in: []
studied_at: "2026-08-25"
codex_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Codex GitHub Action

> **Applicability:** Verbatim excerpts from Codex documentation (snapshot 2026-08-25).

### Source: Codex GitHub Action — overview; Configure codex exec; Manage privileges

> Use the Codex GitHub Action (`openai/codex-action@v1`) to run Codex in CI/CD jobs, apply patches, or post reviews from a GitHub Actions workflow.
> The action installs the Codex CLI, starts the Responses API proxy when you provide an API key, and runs `codex exec` under the permissions you specify.

> ## Prerequisites
>
> - Store your OpenAI key as a GitHub secret (for example `OPENAI_API_KEY`) and reference it in the workflow.
> - Run the job on a Linux or macOS runner. For Windows, set `safety-strategy: unsafe`.
> - Check out your code before invoking the action so Codex can read the repository contents.
> - Decide which prompts you want to run. You can provide inline text via `prompt` or point to a file committed in the repo with `prompt-file`.

> ## Configure `codex exec`
>
> Fine-tune how Codex runs by setting the action inputs that map to `codex exec` options:
>
> - `prompt` or `prompt-file` (choose one): Inline instructions or a repository path to Markdown or text with your task. Consider storing prompts in `.github/codex/prompts/`.
> - `codex-args`: Extra CLI flags. Provide a JSON array (for example `["--ephemeral"]`) or a shell string (`--profile ci`) to configure sessions, profiles, or MCP settings.
> - `model` and `effort`: Pick the Codex agent configuration you want; leave empty for defaults.
> - `sandbox`: Match the sandbox mode (`workspace-write`, `read-only`, `danger-full-access`) to the permissions Codex needs during the run.
> - `output-file`: Save the final Codex message to disk so later steps can upload or diff it.
> - `codex-version`: Pin a specific CLI release. Leave blank to use the latest published version.
> - `codex-home`: Point to a shared Codex home directory if you want to reuse configuration files or MCP setups across steps.

> ## Manage privileges
>
> Codex has broad access on GitHub-hosted runners unless you restrict it. Use these inputs to control exposure:
>
> - `safety-strategy` (default `drop-sudo`) removes `sudo` before running Codex. This is irreversible for the job and protects secrets in memory. On Windows you must set `safety-strategy: unsafe`.
> - `unprivileged-user` pairs `safety-strategy: unprivileged-user` with `codex-user` to run Codex as a specific account. Ensure the user can read and write the repository checkout (see the [`unprivileged-user` example](https://github.com/openai/codex-action/blob/main/examples/unprivileged-user.yml) for an ownership fix).
> - `read-only` keeps Codex from changing files or using the network, but it still runs with elevated privileges. Don't rely on `read-only` alone to protect secrets.
> - `sandbox` limits filesystem and network access within Codex itself. Choose the narrowest option that still lets the task complete.
> - `allow-users` and `allow-bots` restrict who can trigger the workflow. By default only users with write access can run the action; list extra trusted accounts explicitly or leave the field empty for the default behavior.
