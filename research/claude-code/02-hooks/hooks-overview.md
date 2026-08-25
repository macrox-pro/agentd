---
primary_sources:
  - id: T1-HOOKS-GUIDE
    title: "Automate actions with hooks"
    url: "https://code.claude.com/docs/en/hooks-guide.md"
    section: "Overview"
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hooks overview

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Hooks guide — Introduction

> # Automate actions with hooks
>
> > Run shell commands automatically when Claude Code edits files, finishes tasks, or needs input. Format code, send notifications, validate commands, and enforce project rules.
>
> Hooks are user-defined shell commands. Claude Code runs them at specific points in its lifecycle, which gives you deterministic control: certain actions always happen rather than relying on the LLM to choose to run them. Use hooks to enforce project rules, automate repetitive tasks, and integrate Claude Code with your existing tools.
>
> For decisions that require judgment rather than deterministic rules, you can also use [prompt-based hooks](#prompt-based-hooks) or [agent-based hooks](#agent-based-hooks) that use a Claude model to evaluate conditions.
>
> For other ways to extend Claude Code, see [skills](/docs/en/skills) for giving Claude additional instructions and executable commands, [subagents](/docs/en/sub-agents) for running tasks in isolated contexts, and [plugins](/docs/en/plugins) for packaging extensions to share across projects.
>
>
>   This guide covers common use cases and how to get started. For full event schemas, JSON input/output formats, and advanced features like async hooks and MCP tool hooks, see the [Hooks reference](/docs/en/hooks).
