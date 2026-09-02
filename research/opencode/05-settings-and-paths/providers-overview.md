---
primary_sources:
  - id: T1-PROVIDERS
    title: "Providers"
    url: "https://opencode.ai/docs/providers.md"
    section: "Overview and auth"
also_cited_in: []
studied_at: "2026-09-02"
opencode_docs_snapshot: "2026-09-02"
applicability: "current"
---
# Providers overview

> **Applicability:** Verbatim excerpts from OpenCode documentation (snapshot 2026-09-02).

### Source: OpenCode Providers — essentials (auth and setup)

> import config from "../../../config.mjs"
> export const console = config.console
>
> OpenCode uses the [AI SDK](https://ai-sdk.dev/) and [Models.dev](https://models.dev) to support **75+ LLM providers** and it supports running local models.
>
> To add a provider you need to:
>
> 1. Add the API keys for the provider using the `/connect` command.
> 2. Configure the provider in your OpenCode config.
>
> ---
>
> ### Credentials
>
> When you add a provider's API keys with the `/connect` command, they are stored
> in `~/.local/share/opencode/auth.json`.
>
> ---
>
> ### Config
>
> You can customize the providers through the `provider` section in your OpenCode
> config.
>
> ---
>
> #### Base URL
>
> You can customize the base URL for any provider by setting the `baseURL` option. This is useful when using proxy services or custom endpoints.
>
> ```json title="opencode.json" {6}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "provider": {
>     "anthropic": {
>       "options": {
>         "baseURL": "https://api.anthropic.com/v1"
>       }
>     }
>   }
> }
> ```
>
> ---
>
> #### Hiding models
>
> You can hide specific models from the `/models` picker for a provider using the `blacklist` option. This is useful when a provider exposes models you don't want to use or select.
>
> ```json title="opencode.json" {6}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "provider": {
>     "anthropic": {
>       "blacklist": ["claude-opus-4-20250514"]
>     }
>   }
> }
> ```
>
> The inverse `whitelist` option hides every model except the ones listed.
>
> ```json title="opencode.json" {6}
> {
>   "$schema": "https://opencode.ai/config.json",
>   "provider": {
>     "anthropic": {
>       "whitelist": ["claude-sonnet-4-20250514"]
>     }
>   }
> }
> ```
>
> Both options take an array of model IDs — the same IDs shown in the `/models` picker.
>
> - `blacklist` removes the listed models from the picker.
> - `whitelist` keeps only the listed models and hides the rest.
> - You can combine them: `whitelist` narrows the set, then `blacklist` removes entries from it.
>
> ---
>
> ## OpenCode Zen
>
> OpenCode Zen is a list of models provided by the OpenCode team that have been
> tested and verified to work well with OpenCode. [Learn more](/docs/zen).
>
> :::tip
> If you are new, we recommend starting with OpenCode Zen.
> :::
>
> 1. Run the `/connect` command in the TUI, select `OpenCode Zen`, and head to [opencode.ai/auth](https://opencode.ai/zen).
>
>    ```txt
>    /connect
>    ```
>
> 2. Sign in, add your billing details, and copy your API key.
>
> 3. Paste your API key.
>
>    ```txt
>    ┌ API key
>    │
>    │
>    └ enter
>    ```
>
> 4. Run `/models` in the TUI to see the list of models we recommend.
>
>    ```txt
>    /models
>    ```
>
> It works like any other provider in OpenCode and is completely optional to use.
>
> ---
>
> ## OpenCode Go
>
> OpenCode Go is a low cost subscription plan that provides reliable access to popular open coding models provided by the OpenCode team that have been
> tested and verified to work well with OpenCode.
>
> 1. Run the `/connect` command in the TUI, select `OpenCode Go`, and head to [opencode.ai/auth](https://opencode.ai/zen).
>
>    ```txt
>    /connect
>    ```
>
> 2. Sign in, add your billing details, and copy your API key.
>
> 3. Paste your API key.
>
>    ```txt
>    ┌ API key
>    │
>    │
>    └ enter
>    ```
>
> 4. Run `/models` in the TUI to see the list of models we recommend.
>
>    ```txt
>    /models
>    ```
>
> It works like any other provider in OpenCode and is completely optional to use.
>
> ---
>
> ## Directory
>
> Let's look at some of the providers in detail. If you'd like to add a provider to the
> list, feel free to open a PR.
>
> :::note
> Don't see a provider here? Submit a PR.
> :::
>
> ---
>
> ### 302.AI
>
> 1. Head over to the [302.AI console](https://302.ai/), create an account, and generate an API key.
>
> 2. Run the `/connect` command and search for **302.AI**.
>
>    ```txt
>    /connect
>    ```
>
> 3. Enter your 302.AI API key.
>
>    ```txt
>    ┌ API key
>    │
>    │
>    └ enter
>    ```
>
> 4. Run the `/models` command to select a model.
>
>    ```txt
>    /models
>    ```
>
> ---
>
> ### Amazon Bedrock
>
> To use Amazon Bedrock with OpenCode:
>
> 1. Head over to the **Model catalog** in the Amazon Bedrock console and request
>    access to the models you want.
>
>    :::tip
>    You need to have access to the model you want in Amazon Bedrock.
>    :::
>
> 2. **Configure authentication** using one of the following methods:
