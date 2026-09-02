---
primary_sources:
  - id: T1-PLUGINS
    title: "Plugins"
    url: "https://opencode.ai/docs/plugins.md"
    section: "Session Events"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Plugin events — session

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Plugins — Session Events

> #### Session Events
>
> - `session.created`
> - `session.compacted`
> - `session.deleted`
> - `session.diff`
> - `session.error`
> - `session.idle`
> - `session.status`
> - `session.updated`

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
