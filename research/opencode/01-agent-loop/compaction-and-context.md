---
primary_sources:
  - id: T1-PLUGINS
    title: "Plugins"
    url: "https://opencode.ai/docs/plugins.md"
    section: "Compaction hooks"
  - id: T1-CONFIG
    title: "Config"
    url: "https://opencode.ai/docs/config.md"
    section: "Compaction"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Compaction and context

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Plugins — Compaction hooks

> ### Compaction hooks
>
> Customize the context included when a session is compacted:
>
> ```ts title=".opencode/plugins/compaction.ts"
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const CompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Inject additional context into the compaction prompt
>       output.context.push(`
>
> export const CompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Inject additional context into the compaction prompt
>       output.context.push(`
>
>     "experimental.session.compacting": async (input, output) => {
>       // Inject additional context into the compaction prompt
>       output.context.push(`
>
> The `experimental.session.compacting` hook fires before the LLM generates a continuation summary. Use it to inject domain-specific context that the default compaction prompt would miss.
>
> You can also replace the compaction prompt entirely by setting `output.prompt`:
>
> ```ts title=".opencode/plugins/custom-compaction.ts"
> import type { Plugin } from "@opencode-ai/plugin"
>
> export const CustomCompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Replace the entire compaction prompt
>       output.prompt = `
> You are generating a continuation prompt for a multi-agent swarm session.
>
> Summarize:
> 1. The current task and its status
> 2. Which files are being modified and by whom
> 3. Any blockers or dependencies between agents
> 4. The next steps to complete the work
>
> Format as a structured prompt that a new agent can use to resume work.
> `
>     },
>   }
> }
> ```
>
> When `output.prompt` is set, it completely replaces the default compaction prompt. The `output.context` array is ignored in this case.
>
> export const CustomCompactionPlugin: Plugin = async (ctx) => {
>   return {
>     "experimental.session.compacting": async (input, output) => {
>       // Replace the entire compaction prompt
>       output.prompt = `
> You are generating a continuation prompt for a multi-agent swarm session.
>
> Summarize:
> 1. The current task and its status
> 2. Which files are being modified and by whom
> 3. Any blockers or dependencies between agents
> 4. The next steps to complete the work
>
> Format as a structured prompt that a new agent can use to resume work.
> `
>     },
>   }
> }
> ```
>
> When `output.prompt` is set, it completely replaces the default compaction prompt. The `output.context` array is ignored in this case.
>
>     "experimental.session.compacting": async (input, output) => {
>       // Replace the entire compaction prompt
>       output.prompt = `
> You are generating a continuation prompt for a multi-agent swarm session.
>
> Summarize:
> 1. The current task and its status
> 2. Which files are being modified and by whom
> 3. Any blockers or dependencies between agents
> 4. The next steps to complete the work
>
> Format as a structured prompt that a new agent can use to resume work.
> `
>     },
>   }
> }
> ```
>
> When `output.prompt` is set, it completely replaces the default compaction prompt. The `output.context` array is ignored in this case.
