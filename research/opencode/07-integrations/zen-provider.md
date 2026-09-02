---
primary_sources:
  - id: T1-ZEN
    title: "Zen"
    url: "https://opencode.ai/docs/zen.md"
    section: "Gateway and connect"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# OpenCode Zen provider

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Zen — overview and /connect (trimmed)

> import config from "../../../config.mjs"
> export const console = config.console
> export const email = `mailto:${config.email}`
>
> OpenCode Zen is a list of tested and verified models provided by the OpenCode team.
>
> Zen works like any other provider in OpenCode. You login to OpenCode Zen and get
> your API key. It's **completely optional** and you don't need to use it to use
> OpenCode.
>
> ---
>
> ## Background
>
> There are a large number of models out there but only a few of
> these models work well as coding agents. Additionally, most providers are
> configured very differently; so you get very different performance and quality.
>
> :::tip
> We tested a select group of models and providers that work well with OpenCode.
> :::
>
> So if you are using a model through something like OpenRouter, you can never be
> sure if you are getting the best version of the model you want.
>
> To fix this, we did a couple of things:
>
> 1. We tested a select group of models and talked to their teams about how to
>    best run them.
> 2. We then worked with a few providers to make sure these were being served
>    correctly.
> 3. Finally, we benchmarked the combination of the model/provider and came up
>    with a list that we feel good recommending.
>
> OpenCode Zen is an AI gateway that gives you access to these models.
>
> ---
>
> ## How it works
>
> OpenCode Zen works like any other provider in OpenCode.
>
> 1. You sign in to **<a href={console}>OpenCode Zen</a>**, add your billing
>    details, and copy your API key.
> 2. You run the `/connect` command in the TUI, select OpenCode Zen, and paste your API key.
> 3. Run `/models` in the TUI to see the list of models we recommend.
>
> You are charged per request and you can add credits to your account.
>
> ---
>
> ## Endpoints
>
> You can also access our models through the following API endpoints.
>
> | Model                           | Model ID                        | Endpoint                                                  | AI SDK Package              |
> | ------------------------------- | ------------------------------- | --------------------------------------------------------- | --------------------------- |
> | GPT 5.6 Sol                     | gpt-5.6-sol                     | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.6 Terra                   | gpt-5.6-terra                   | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.6 Luna                    | gpt-5.6-luna                    | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.5                         | gpt-5.5                         | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.5 Pro                     | gpt-5.5-pro                     | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.4                         | gpt-5.4                         | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.4 Pro                     | gpt-5.4-pro                     | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.4 Mini                    | gpt-5.4-mini                    | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.4 Nano                    | gpt-5.4-nano                    | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.3 Codex                   | gpt-5.3-codex                   | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.3 Codex Spark             | gpt-5.3-codex-spark             | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.2                         | gpt-5.2                         | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.2 Codex                   | gpt-5.2-codex                   | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.1                         | gpt-5.1                         | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.1 Codex                   | gpt-5.1-codex                   | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.1 Codex Max               | gpt-5.1-codex-max               | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5.1 Codex Mini              | gpt-5.1-codex-mini              | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5                           | gpt-5                           | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5 Codex                     | gpt-5-codex                     | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | GPT 5 Nano                      | gpt-5-nano                      | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | Claude Fable 5.1                | claude-fable-5-1                | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Fable 5                  | claude-fable-5                  | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Opus 5                   | claude-opus-5                   | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Opus 4.8                 | claude-opus-4-8                 | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Opus 4.7                 | claude-opus-4-7                 | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Opus 4.6                 | claude-opus-4-6                 | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Opus 4.5                 | claude-opus-4-5                 | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Sonnet 5                 | claude-sonnet-5                 | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Sonnet 4.6               | claude-sonnet-4-6               | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Sonnet 4.5               | claude-sonnet-4-5               | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Claude Haiku 4.5                | claude-haiku-4-5                | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Gemini 3.8 Flash                | gemini-3.8-flash                | `https://opencode.ai/zen/v1/models/gemini-3.8-flash`      | `@ai-sdk/google`            |
> | Gemini 3.7 Flash                | gemini-3.7-flash                | `https://opencode.ai/zen/v1/models/gemini-3.7-flash`      | `@ai-sdk/google`            |
> | Gemini 3.6 Flash                | gemini-3.6-flash                | `https://opencode.ai/zen/v1/models/gemini-3.6-flash`      | `@ai-sdk/google`            |
> | Gemini 3.5 Flash                | gemini-3.5-flash                | `https://opencode.ai/zen/v1/models/gemini-3.5-flash`      | `@ai-sdk/google`            |
> | Gemini 3.5 Flash Lite           | gemini-3.5-flash-lite           | `https://opencode.ai/zen/v1/models/gemini-3.5-flash-lite` | `@ai-sdk/google`            |
> | Gemini 3.1 Pro                  | gemini-3.1-pro                  | `https://opencode.ai/zen/v1/models/gemini-3.1-pro`        | `@ai-sdk/google`            |
> | Gemini 3 Flash                  | gemini-3-flash                  | `https://opencode.ai/zen/v1/models/gemini-3-flash`        | `@ai-sdk/google`            |
> | Grok 4.6                        | grok-4.6                        | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | Grok 4.5                        | grok-4.5                        | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | Grok Build 0.1                  | grok-build-0.1                  | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | Muse Spark 1.2                  | muse-spark-1.2                  | `https://opencode.ai/zen/v1/responses`                    | `@ai-sdk/openai`            |
> | Qwen3.7 Max                     | qwen3.7-max                     | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Qwen3.7 Plus                    | qwen3.7-plus                    | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Qwen3.6 Plus                    | qwen3.6-plus                    | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | Qwen3.5 Plus                    | qwen3.5-plus                    | `https://opencode.ai/zen/v1/messages`                     | `@ai-sdk/anthropic`         |
> | DeepSeek V4 Pro                 | deepseek-v4-pro                 | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | DeepSeek V4 Flash               | deepseek-v4-flash               | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | MiniMax M3                      | minimax-m3                      | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | MiniMax M2.7                    | minimax-m2.7                    | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | MiniMax M2.5                    | minimax-m2.5                    | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | GLM 5.2                         | glm-5.2                         | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | GLM 5.1                         | glm-5.1                         | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | GLM 5                           | glm-5                           | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | Kimi K2.5                       | kimi-k2.5                       | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | Kimi K2.6                       | kimi-k2.6                       | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | Kimi K2.7 Code                  | kimi-k2.7-code                  | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | Kimi K3                         | kimi-k3                         | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | Big Pickle                      | big-pickle                      | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | MiMo-V2.5 Free                  | mimo-v2.5-free                  | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | Ling 3.0 Flash Fin Free         | ling-3.0-flash-fin-free         | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | Nemotron 3 Ultra Free           | nemotron-3-ultra-free           | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
> | Nemotron 3.5 Lightning Free     | nemotron-3.5-lightning-free     | `https://opencode.ai/zen/v1/chat/completions`             | `@ai-sdk/openai-compatible` |
