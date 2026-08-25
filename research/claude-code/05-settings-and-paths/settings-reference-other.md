---
primary_sources:
  - id: T1-SETTINGS-REF
    title: "Settings reference"
    url: "https://code.claude.com/docs/en/settings-reference.md"
    section: "Remaining keys by category"
  - id: T1-SETTINGS-EXAMPLE
    title: "Example settings files"
    url: "https://code.claude.com/docs/en/settings-example.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Settings reference — other keys

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Settings reference — remaining keys

> ## Model and responses
>
> Choose which models Claude Code uses and how it responds. For how these settings interact with the `/model` command and environment variables, see [Model configuration](/docs/en/model-config).
>
> ### `advisorModel`
>
> Pick which model answers when Claude calls the server-side [advisor tool](/docs/en/advisor). Unset it to turn the advisor off. The advisor must be at least as capable as your main model; when it isn't, Claude Code sends requests without the advisor. See [Choose an advisor model](/docs/en/advisor#choose-an-advisor-model).
>
> You don't usually edit this key by hand. Run `/advisor` to open a picker that shows the current choice, the models that can advise, and **No advisor**. Claude Code saves your pick to this key in `~/.claude/settings.json`. In a session attached to a remote worker, the pick applies to that session only.
>
> To pick Fable, first accept the [usage-credits consent](/docs/en/advisor#fable-advisor-and-usage-credits) by running `/model fable`. Until you do, picking Fable in `/advisor` saves nothing and Claude Code tells you to run `/model fable` first.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of the aliases `"fable"`, `"opus"`, or `"sonnet"`, which resolve to Claude Code's current default version of that model family, or a full model ID such as `"claude-opus-5"`
> * **Default**: unset, so the advisor is off
> * **Per-session overrides**: `--advisor` takes precedence over this key for one session. [`CLAUDE_CODE_DISABLE_ADVISOR_TOOL`](/docs/en/env-vars) turns the advisor off, and this key can't turn it back on
>
> ```json settings.json
> {
>   "advisorModel": "opus"
> }
> ```
>
> The key has no effect on Amazon Bedrock, Google Cloud's Agent Platform, or Microsoft Foundry. `"fable"` requires [Fable 5 access](/docs/en/advisor#choose-an-advisor-model).
>
> ### `alwaysThinkingEnabled`
>
> Turn [extended thinking](/docs/en/model-config#extended-thinking) off for every session by setting this to `false`. Thinking is on by default, so `true` changes nothing. Most people set this through `/config` rather than by editing the file.
>
> On models that always think, such as Fable 5, `false` has no effect. On [third-party providers](/docs/en/third-party-integrations) Claude Code omits the `thinking` parameter instead of turning thinking off, so adaptive-reasoning models may still think.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: no effect; thinking is already on
>   * `false`: Claude Code turns extended thinking off for every session
> * **Default**: unset, so thinking is on for models that support it
> * **Per-session overrides**: [`MAX_THINKING_TOKENS`](/docs/en/env-vars) takes precedence over this key for one session: `0` turns thinking off, under the same model and provider limits as `false`, and a positive value turns thinking on even when this key is `false`. On adaptive-reasoning models the number itself is ignored
>
> ```json settings.json
> {
>   "alwaysThinkingEnabled": false
> }
> ```
>
> ### `availableModels`
>
> Restrict which models people can select for the main session, [subagents](/docs/en/sub-agents), [skills](/docs/en/skills), and the [advisor](/docs/en/advisor). A managed list constrains `/model`, `--model`, and the `model` key in a developer's own files; a model outside it can't be selected. On its own this doesn't touch the Default option; pair it with [`enforceAvailableModels`](#enforceavailablemodels) for that.
>
> * **Scope**: [`Any file`](#scopes). Deploy it in managed settings to enforce it for an organization.
> * **Type**: array of model aliases or IDs
> * **Default**: unset, so every model is available
>
> This example lets people select only Sonnet and Haiku models:
>
> ```json settings.json
> {
>   "availableModels": ["sonnet", "haiku"]
> }
> ```
>
> See [Restrict model selection](/docs/en/model-config#restrict-model-selection).
>
> ### `effortLevel`
>
> Keep an [effort level](/docs/en/model-config#adjust-effort-level) across sessions. Lower levels are faster and cheaper on straightforward tasks, and higher levels reason more deeply on complex problems. Claude Code writes this key to your user settings when you run `/effort low`, `medium`, `high`, or `xhigh` in an interactive session on your machine. In a `-p` run, the Agent SDK, or a session attached to a remote worker, `/effort` applies to that session only. The message `/effort` prints says which happened.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"low"`: the least reasoning, for short, scoped, latency-sensitive tasks that aren't intelligence-sensitive
>   * `"medium"`: reduces token usage for cost-sensitive work that can trade off some intelligence
>   * `"high"`: balances token usage and intelligence
>   * `"xhigh"`: deeper reasoning at higher token spend
> * **Default**: unset
> * **Per-session overrides**: `--effort` takes precedence over this key for one session, and [`CLAUDE_CODE_EFFORT_LEVEL`](/docs/en/env-vars) takes precedence over both
>
> ```json settings.json
> {
>   "effortLevel": "xhigh"
> }
> ```
>
> On Opus 4.7, Opus 4.8, and Fable 5, Claude Code holds that model's default effort until you change effort once with `/effort`, `--effort`, or the `/model` picker. After that, it reads this key. See [Adjust effort level](/docs/en/model-config#adjust-effort-level).
>
> ### `enforceAvailableModels`
>
> The `/model` picker has a **Default** option that resolves to your [organization default model](/docs/en/model-config#organization-default-model) when one applies, and otherwise to your account type's default. An [`availableModels`](#availablemodels) allowlist limits the models you can name, but on its own it leaves **Default** alone, so **Default** can still resolve to a model outside the list. This key closes that gap. Requires Claude Code v2.1.175 or later.
>
> When your organization deploys any managed settings, Claude Code reads this key from the managed source alone and ignores it in your other files.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: when **Default** would resolve to a model outside `availableModels`, Claude Code resolves it to the first available model in the list
>   * `false`: **Default** resolves as usual, even to a model outside `availableModels`
> * **Default**: `false`
>
> This example restricts named selections to Sonnet and Haiku models and makes **Default** resolve to the first of them that is available:
>
> ```json settings.json
> {
>   "availableModels": ["sonnet", "haiku"],
>   "enforceAvailableModels": true
> }
> ```
>
> This key has no effect when `availableModels` is unset or empty. See [Enforce the allowlist for the Default model](/docs/en/model-config#enforce-the-allowlist-for-the-default-model). Requires Claude Code v2.1.175 or later.
>
> ### `fallbackModel`
>
> Name backup models for Claude Code to try, in order, when your primary model is overloaded or unavailable. Claude Code switches to the next available model in the chain for the rest of the turn and shows a notice. Without a chain, Claude Code retries the same model and then surfaces the server's error, and you retry or switch models yourself.
>
> A switch means one turn with a cold [prompt cache](/docs/en/prompt-caching#switching-models) on the fallback model; your next message tries the primary model first again.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of model aliases or IDs; `"default"` expands to the default model
> * **Default**: unset, so a failed request isn't retried on another model
> * **Per-session overrides**: `--fallback-model` takes precedence over this key for one session
>
> This example tries Sonnet 5 first, then Haiku 4.5, when your primary model fails:
>
> ```json settings.json
> {
>   "fallbackModel": ["claude-sonnet-5", "claude-haiku-4-5"]
> }
> ```
>
> Unlike most array settings, this key doesn't merge across settings files: the highest-precedence file that defines it supplies the whole chain. If your project file sets `["claude-sonnet-5"]` and your user file sets `["claude-haiku-4-5"]`, the chain is `["claude-sonnet-5"]` only. Claude Code keeps at most three distinct allowed models from the list and ignores the rest. See [Fallback model chains](/docs/en/model-config#fallback-model-chains).
>
> ### `fastMode`
>
> Turn [fast mode](/docs/en/fast-mode) on for sessions where it's available, for interactive work like rapid iteration or live debugging where you want speed at a higher cost per token. You don't usually edit this key by hand: running `/fast` writes `fastMode: true` to `~/.claude/settings.json`, and running it again to turn fast mode off removes the key. Fast mode runs only on Opus 5 and Opus 4.8: turning it on from another model switches you to Opus, and switching to an unsupported model turns it off. See [Switch models while fast mode is on](/docs/en/fast-mode#switch-models-while-fast-mode-is-on).
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code turns fast mode on for sessions where it's available
>   * `false`: fast mode stays off
> * **Default**: unset, so fast mode is off
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_FAST_MODE`](/docs/en/env-vars) turns fast mode off for one session, and this key can't turn it back on
>
> ```json settings.json
> {
>   "fastMode": true
> }
> ```
>
> ### `fastModePerSessionOptIn`
>
> Normally, running `/fast` saves [`fastMode`](#fastmode) to a person's user settings, so fast mode is on at the start of every later session. Set this key to `true` to stop that: a saved `fastMode: true` no longer turns fast mode on at session start, and each person has to run `/fast` in each session they want it. Claude Code leaves the `fastMode` key in their file, so turning this key off restores the old behavior. Owners on Team or Enterprise plans can deploy it organization-wide through [server-managed settings](/docs/en/server-managed-settings).
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: a saved `fastMode: true` no longer turns fast mode on at session start, so each person runs `/fast` in each session they want it; a `fastMode: true` passed with `--settings` still counts for that session unless managed settings set this key
>   * `false`: a saved `fastMode: true` turns fast mode on at the start of every later session
> * **Default**: `false`
>
> ```json settings.json
> {
>   "fastModePerSessionOptIn": true
> }
> ```
>
> See [Require per-session opt-in](/docs/en/fast-mode#require-per-session-opt-in).
>
> ### `language`
>
> Have Claude respond in a language other than English by default. There is no fixed list for responses: Claude Code adds the value verbatim to the system prompt as an instruction to always respond in that language, so any language name Claude can read works. Claude Code doesn't check the value, so a misspelled name reaches Claude as written rather than producing an error. The same value sets the language for [voice dictation](/docs/en/voice-dictation#change-the-dictation-language), which does have a fixed list of [supported dictation languages](/docs/en/voice-dictation#change-the-dictation-language), and for auto-generated session titles.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, any language name, such as `"japanese"`, `"spanish"`, or `"french"`; Claude Code doesn't validate it
> * **Default**: unset; session titles then match the language of your conversation
>
> ```json settings.json
> {
>   "language": "japanese"
> }
> ```
>
> ### `model`
>
> Set the model every new session uses, so you don't have to pick one with `/model` each time. Setting it here doesn't stop you from switching mid-session. If your admin set an [organization default model](/docs/en/model-config#organization-default-model) to override user selection, you get that model even when you set this key in user, project, or local settings.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, a model alias or full model ID
> * **Default**: unset, so Claude Code uses your account's default model
> * **Per-session overrides**: `--model` takes precedence over [`ANTHROPIC_MODEL`](/docs/en/env-vars), and both take precedence over this key for one session, including over a managed `model`; an [`availableModels`](#availablemodels) list still applies to the pick
>
> ```json settings.json
> {
>   "model": "claude-sonnet-5"
> }
> ```
>
> A value here outranks [`ANTHROPIC_DEFAULT_MODEL`](/docs/en/model-config#set-a-default-model-for-new-sessions), which Claude Code uses only when nothing else selects a model.
>
> ### `modelOverrides`
>
> Map Anthropic model IDs to provider-specific model IDs, such as Amazon Bedrock inference profile ARNs. Each model picker entry then uses its mapped value when calling the provider API. Administrators use this on [Amazon Bedrock, Google Cloud's Agent Platform, and Microsoft Foundry](/docs/en/model-config#override-model-ids-per-version) to route each model version to a specific inference profile, version name, or deployment for governance, cost allocation, or regional routing.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object mapping model ID to provider model ID
> * **Default**: unset
>
> This example routes every call for Opus 4.6 to the named Bedrock inference profile:
>
> ```json settings.json
> {
>   "modelOverrides": {
>     "claude-opus-4-6": "arn:aws:bedrock:us-east-1:123456789012:inference-profile/example"
>   }
> }
> ```
>
> See [Override model IDs per version](/docs/en/model-config#override-model-ids-per-version).
>
> ### `modelPicker`
>
> List the models the `/model` picker offers, in the order you write them and under labels you choose, so the picker lists the models your organization runs, after the built-in lineup or instead of it. Each row's `model` is taken verbatim, so it accepts anything `--model` accepts: an alias such as `opus`, an Anthropic model ID, or a provider-format ID for Amazon Bedrock, Google Cloud's Agent Platform, Microsoft Foundry, or an LLM gateway. Requires Claude Code v2.1.242 or later.
>
> * **Scope**: [`User or managed`](#scopes). Claude Code reads the key from managed settings, `--settings`, and user settings, and ignores it in project and local settings so a repository you clone can't relabel the picker. The highest of those three that sets the key supplies the whole lineup, and Claude Code never combines lineups from two sources.
> * **Type**: object with an `options` array of rows and an optional `replaceBuiltInOptions` Boolean
> * **Default**: unset, so the picker shows the built-in lineup
>
> This example adds two Bedrock deployments after the built-in lineup, under names your team recognizes:
>
> ```json managed-settings.json
> {
>   "modelPicker": {
>     "options": [
>       { "model": "us.anthropic.claude-opus-4-8", "label": "Opus (production)" },
>       {
>         "model": "us.anthropic.claude-sonnet-4-6",
>         "label": "Sonnet (production)",
>         "description": "Day-to-day work"
>       }
>     ]
>   }
> }
> ```
>
>
>
> #### Fields for `modelPicker`
>
> The key takes two fields, one for the rows themselves and one for whether they replace the built-in lineup or add to it.
>
> | Field                   | Type                                                                                  | What it does                                                                                                                                                                                                                                                                  |
> | :---------------------- | :------------------------------------------------------------------------------------ | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `options`               | array of rows, each with a required `model` and an optional `label` and `description` | The rows the picker shows, in this order, except that a grayed-out row moves to the bottom. Without a `label`, Claude Code titles the row with the built-in name for a model it knows, or the model ID otherwise, and without a `description` it writes a generic second line |
> | `replaceBuiltInOptions` | Boolean, default `false`                                                              | Set it to `true` to show only these rows, **Default**, and a row for the model the session is already using. Leave it unset to add these rows after the built-in lineup                                                                                                       |
>
> With `replaceBuiltInOptions` on, Claude Code hides every other row: the built-in lineup, the rows it adds for [`availableModels`](#availablemodels) entries, the models [gateway discovery](/docs/en/llm-gateway-protocol#model-discovery) found, and [`ANTHROPIC_CUSTOM_MODEL_OPTION`](/docs/en/model-config#add-a-custom-model-option). With it off, Claude Code skips a listed model that the built-in lineup already covers. A label changes what the picker shows, not which model Claude Code runs.
>
> An [`availableModels`](#availablemodels) allowlist still applies to these rows. Before you add a listed model to the allowlist, read [Merge behavior](/docs/en/model-config#merge-behavior): a specific model ID narrows its family's wildcard entry. Claude Code also checks each row against the session before it shows the picker:
>
> * **Dropped**: a row Claude Code can't serve, such as a retired model or a model your organization has no access to
> * **Grayed out**: a row you can't select yet, shown with the reason
> * **No row survives**: Claude Code keeps the built-in lineup, filtered by the allowlist as usual
>
> Claude Code drops a row it can't parse and keeps the rest. See [Fix a broken settings file](/docs/en/settings#fix-a-broken-settings-file).
>
> ### `outputStyle`
>
> Select an [output style](/docs/en/output-styles) by name. An output style is a saved set of instructions that Claude Code adds to the system prompt to change Claude's role, tone, and output format, such as the built-in Explanatory and Learning styles or one you wrote yourself.
>
> Claude Code builds the style into the system prompt once per conversation. An edit to this key takes effect after you run `/clear` or start a new session.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, the name of a [built-in](/docs/en/output-styles#built-in-output-styles) or [custom](/docs/en/output-styles#create-a-custom-output-style) output style
> * **Default**: unset, so Claude Code uses the default style
>
> This example selects the built-in Explanatory style, which adds educational insights between tasks:
>
> ```json settings.json
> {
>   "outputStyle": "Explanatory"
> }
> ```
>
> ### `promptCacheTtl`
>
> Choose how long the [prompt cache](/docs/en/prompt-caching) holds the main conversation. This key applies to your interactive, `-p`, and Agent SDK turns, together with the helpers Claude Code runs inline with them. The one-hour lifetime keeps the cache warm across longer breaks, and the API [bills each cache write at a higher rate](https://platform.claude.com/docs/en/build-with-claude/prompt-caching#pricing) than at the five-minute lifetime. Requires Claude Code v2.1.242 or later.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"5m"`: the cache holds for five minutes
>   * `"1h"`: the cache holds for an hour
> * **Default**: unset, so each main-conversation request gets [its default lifetime](/docs/en/prompt-caching#which-ttl-each-request-gets)
> * **Per-session overrides**: [`FORCE_PROMPT_CACHING_5M`](/docs/en/env-vars) takes precedence over everything else, then [`CLAUDE_CODE_PROMPT_CACHE_TTL`](/docs/en/env-vars), then this key, and last [`ENABLE_PROMPT_CACHING_1H`](/docs/en/env-vars)
>
> This example keeps the main conversation on the one-hour lifetime and leaves subagents on five minutes:
>
> ```json settings.json
> {
>   "promptCacheTtl": "1h",
>   "subagentPromptCacheTtl": "5m"
> }
> ```
>
> For what each lifetime costs, see [Cache lifetime](/docs/en/prompt-caching#cache-lifetime).
>
> ### `showThinkingSummaries`
>
> See summaries of Claude's [extended thinking](/docs/en/model-config#extended-thinking) in interactive sessions. Set it if you want the full summaries when you expand thinking with `Ctrl+O`. When unset or `false`, the Anthropic API redacts thinking blocks and Claude Code shows a collapsed stub; third-party providers don't redact.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: you see full thinking summaries when you expand thinking with `Ctrl+O`
>   * `false`: the Anthropic API redacts thinking blocks and Claude Code shows a collapsed stub
> * **Default**: `false`
>
> ```json settings.json
> {
>   "showThinkingSummaries": true
> }
> ```
>
> Redaction only changes what you see, not what the model generates: to reduce thinking spend, [lower the budget or disable thinking](/docs/en/model-config#extended-thinking) instead. This setting has no effect in non-interactive mode (`-p`), the Agent SDK, or IDE extensions such as VS Code.
>
> ### `subagentPromptCacheTtl`
>
> Choose how long the [prompt cache](/docs/en/prompt-caching) holds the requests Claude Code makes outside the main conversation. This key applies to [subagents](/docs/en/sub-agents), [workflows](/docs/en/workflows), and Claude Code's own background and helper requests, such as compaction and session titles. The one-hour lifetime keeps the cache warm across longer breaks, and the API [bills each cache write at a higher rate](https://platform.claude.com/docs/en/build-with-claude/prompt-caching#pricing) than at the five-minute lifetime. Requires Claude Code v2.1.242 or later.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"5m"`: the cache holds for five minutes
>   * `"1h"`: the cache holds for an hour
> * **Default**: unset, so each of these requests gets [its default lifetime](/docs/en/prompt-caching#which-ttl-each-request-gets)
> * **Per-session overrides**: [`FORCE_PROMPT_CACHING_5M`](/docs/en/env-vars) takes precedence over everything else, then [`CLAUDE_CODE_SUBAGENT_PROMPT_CACHE_TTL`](/docs/en/env-vars), then this key, and last [`ENABLE_PROMPT_CACHING_1H`](/docs/en/env-vars), which asks for the one-hour lifetime on every request
>
> This example gives subagents and the other requests outside the main conversation the one-hour lifetime:
>
> ```json settings.json
> {
>   "subagentPromptCacheTtl": "1h"
> }
> ```
>
> This key covers the requests [`promptCacheTtl`](#promptcachettl) doesn't, so set both to choose a lifetime for every request Claude Code makes. For how a subagent's cache differs from the main conversation's, see [Subagents and the cache](/docs/en/prompt-caching#subagents-and-the-cache).
>
> ### `switchModelsOnFlag`
>
> Choose what happens when a [safety classifier flags a request](/docs/en/model-config#automatic-model-fallback): switch to the fallback model and continue, or pause so you can choose between switching and editing the prompt.
>
> * **Scope**: [`Any file`](#scopes). Appears in `/config` as **Switch models when a message is flagged**.
> * **Type**: Boolean
>   * `true`: Claude Code switches to the fallback model and continues
>   * `false`: in an interactive session Claude Code pauses so you can choose between switching and editing the prompt; where no dialog can show, such as a `-p` run, the flagged request ends as an error
> * **Default**: `true`, switch automatically
>
> ```json settings.json
> {
>   "switchModelsOnFlag": false
> }
> ```
>
> See [Ask before switching](/docs/en/model-config#ask-before-switching). Requires Claude Code v2.1.170 or later.
>
> ### `ultracode`
>
> Start sessions with [ultracode](/docs/en/workflows#let-claude-decide-with-ultracode) on. With it on, Claude plans a workflow for each substantive task instead of waiting for you to ask. Claude plans workflows only when [dynamic workflows](/docs/en/workflows) are enabled for you and your model supports `xhigh` effort. Either way, `ultracode: true` runs the session at `xhigh` effort. Claude Code reads this key but never writes it: `/effort ultracode` turns ultracode on for the current session only.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: sessions start at `xhigh` effort, with ultracode on when dynamic workflows are enabled for you and your model supports `xhigh`
>   * `false`: sessions start with ultracode off
> * **Default**: unset, so ultracode is off
> * **Per-session overrides**: `/effort ultracode` turns ultracode on for one session without this key. So does `--effort ultracode`, which requires Claude Code v2.1.203 or later
>
> ```json settings.json
> {
>   "ultracode": true
> }
> ```
>
> Ultracode runs the session at `xhigh` effort and takes precedence over `effortLevel`. An Agent SDK `apply_flag_settings` control request also accepts the key.
>
> ## Sandbox settings
>
> Isolate the commands Claude runs from your filesystem, your network, and your credentials. For how sandboxing works and platform requirements, see [Sandboxing](/docs/en/sandboxing).
>
> ### `sandbox`
>
> Isolate the Bash commands Claude runs from your filesystem and network with [sandboxing](/docs/en/sandboxing). Turn the sandbox on with `enabled`, then narrow or widen what sandboxed commands can touch with the `filesystem`, `network`, and `credentials` sub-objects. The sandbox runs on macOS, Linux, and WSL2.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object with `enabled`, `failIfUnavailable`, `autoAllowBashIfSandboxed`, `excludedCommands`, `allowUnsandboxedCommands`, `enableWeakerNestedSandbox`, `enableWeakerNetworkIsolation`, `allowAppleEvents`, `bwrapPath`, `socatPath`, `ignoreViolations`, and `ripgrep`, plus the `filesystem`, `network`, and `credentials` objects
> * **Default**: unset, so Claude Code runs commands without a sandbox
>
> This turns the sandbox on, skips permission prompts for sandboxed commands, runs `docker` outside the sandbox, opens two extra write paths, hides your AWS credentials file, and pre-allows GitHub and npm:
>
> ```json settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "autoAllowBashIfSandboxed": true,
>     "excludedCommands": ["docker *"],
>     "filesystem": {
>       "allowWrite": ["/tmp/build", "~/.kube"],
>       "denyRead": ["~/.aws/credentials"]
>     },
>     "network": {
>       "allowedDomains": ["github.com", "*.npmjs.org"]
>     }
>   }
> }
> ```
>
> Boolean keys take the value from the highest-precedence settings file that sets them, so a managed `enabled` or `failIfUnavailable` overrides anything a developer sets. Array keys merge across every settings file, so a developer can append entries; see [Keep developers from widening the policy](/docs/en/sandboxing#keep-developers-from-widening-the-policy) for the managed-only locks. To require the sandbox for an organization, see [Enforce sandboxing with managed settings](/docs/en/sandboxing#enforce-sandboxing-with-managed-settings).
>
> ### `sandbox.enabled`
>
> Turn on [sandboxing](/docs/en/sandboxing) for Bash commands. When you pick a mode in the `/sandbox` panel, Claude Code writes this key to `.claude/settings.local.json` for the current project; set it in `~/.claude/settings.json` to sandbox every project.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code sandboxes Bash commands
>   * `false`: Bash commands run unsandboxed
> * **Default**: `false`
>
> ```json settings.json
> {
>   "sandbox": {
>     "enabled": true
>   }
> }
> ```
>
> On Linux and WSL2 the sandbox needs `bubblewrap` and `socat`; see [Set up Linux and WSL2](/docs/en/sandboxing#set-up-linux-and-wsl2). When the sandbox can't start, Claude Code shows a warning and runs commands unsandboxed unless you also set [`failIfUnavailable`](#sandbox-failifunavailable).
>
> ### `sandbox.failIfUnavailable`
>
> Make Claude Code exit with an error at startup when `sandbox.enabled` is `true` but the sandbox can't start, because a dependency is missing or the platform is unsupported. Without it, Claude Code shows a warning and runs commands unsandboxed. Use it in managed settings when your organization requires sandboxing as a hard gate.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code exits with an error at startup when `sandbox.enabled` is `true` but the sandbox can't start
>   * `false`: Claude Code shows a warning and runs commands unsandboxed
> * **Default**: `false`
>
> This makes every managed machine sandbox commands or refuse to start:
>
> ```json managed-settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "failIfUnavailable": true
>   }
> }
> ```
>
> See [Enforce sandboxing with managed settings](/docs/en/sandboxing#enforce-sandboxing-with-managed-settings).
>
> ### `sandbox.autoAllowBashIfSandboxed`
>
> Let Claude Code run sandboxed Bash commands without a permission prompt. Commands that can't run in the sandbox still go through the regular permission flow, and `deny` rules and content-scoped `ask` rules such as `Bash(git push *)` still apply; a bare `Bash` ask rule is skipped for sandboxed commands. Set it to `false` to send sandboxed commands through the regular permission flow too, which the `/sandbox` **Mode** tab calls regular permissions mode.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code runs sandboxed Bash commands without a permission prompt, subject to `deny` rules and content-scoped `ask` rules; `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB` turns auto-allow off
>   * `false`: sandboxed commands go through the regular permission flow, so your allow rules and permission mode decide. The `/sandbox` **Mode** tab calls this regular permissions mode
> * **Default**: `true`
>
> This keeps the sandbox on and sends sandboxed commands through the regular permission flow:
>
> ```json settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "autoAllowBashIfSandboxed": false
>   }
> }
> ```
>
> See [Sandbox modes](/docs/en/sandboxing#sandbox-modes) for what auto-allow mode still prompts on and how it behaves in plan mode.
>
> ### `sandbox.excludedCommands`
>
> Name commands that Claude Code always runs outside the sandbox, such as tools that don't work under it. Each entry uses the same syntax as the content of a `Bash(...)` [permission rule](/docs/en/permissions#permission-rule-syntax): an exact command, a prefix such as `docker *`, or a wildcard pattern. When any part of a compound command matches an entry, Claude Code runs the whole command unsandboxed.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of command patterns
> * **Default**: unset, so every command Claude Code can sandbox runs sandboxed
>
> ```json settings.json
> {
>   "sandbox": {
>     "excludedCommands": ["docker *"]
>   }
> }
> ```
>
> Excluded commands still go through the regular permission flow. Exclusion is a convenience, not a security boundary: prefer [`filesystem.allowWrite`](#sandbox-filesystem-allowwrite) when a tool only needs to write somewhere specific. Entries merge across every settings file, and there is no managed-only lock for this list, so keep a managed list narrow.
>
> ### `sandbox.allowUnsandboxedCommands`
>
> Let Claude retry a command outside the sandbox with the `dangerouslyDisableSandbox` parameter after the sandbox blocks it. Set it to `false` so Claude Code ignores that parameter completely and every command must run sandboxed or appear in [`excludedCommands`](#sandbox-excludedcommands), which the `/sandbox` **Overrides** tab shows as **Strict sandbox mode**. Use `false` in managed settings for policies that require strict sandboxing.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude can retry a command outside the sandbox with the `dangerouslyDisableSandbox` parameter after the sandbox blocks it
>   * `false`: Claude Code ignores that parameter, so every command runs sandboxed or appears in `excludedCommands`
> * **Default**: `true`
>
> This enforces strict sandbox mode for everyone the managed settings cover:
>
> ```json managed-settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "allowUnsandboxedCommands": false
>   }
> }
> ```
>
> An unsandboxed retry goes through the regular permission flow: a prompt in Manual mode, the classifier in auto mode. See [The unsandboxed retry escape hatch](/docs/en/sandboxing#the-unsandboxed-retry-escape-hatch).
>
> ### `sandbox.filesystem`
>
> Control which paths sandboxed commands can read and write. By default they can write to the working directory, any directories you add with `--add-dir`, and the session temp directory, and can read the rest of the filesystem, including credential files. Widen or narrow that with the four path lists, or switch the filesystem layer off with `disabled`. See [Filesystem isolation](/docs/en/sandboxing#filesystem-isolation) for the default boundaries.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object with `allowWrite`, `denyWrite`, `denyRead`, and `allowRead` arrays, plus the `allowManagedReadPathsOnly` and `disabled` Booleans
> * **Default**: unset, so the default read and write boundaries apply
>
> This lets sandboxed commands write to a build directory and your kubeconfig, and hides your AWS credentials file:
>
> ```json settings.json
> {
>   "sandbox": {
>     "filesystem": {
>       "allowWrite": ["/tmp/build", "~/.kube"],
>       "denyRead": ["~/.aws/credentials"]
>     }
>   }
> }
> ```
>
> Claude Code enforces these lists at the OS sandbox boundary, so they apply to every subprocess a sandboxed command starts, such as `kubectl`, `terraform`, or `npm`, not only to Claude's file tools. Your [permission rules](/docs/en/sandboxing#permission-rules) feed the same lists: `Edit` allow and deny rules join `allowWrite` and `denyWrite`, `Read` deny rules join `denyRead`, and `WebFetch(domain:...)` allow and deny rules join the [`network`](#sandbox-network) domain lists. Every list merges across settings files. When you edit a list during a session, Claude Code [applies the change to the running session](/docs/en/settings#when-edits-take-effect).
>
> #### Sandbox path prefixes
>
> Paths in `allowWrite`, `denyWrite`, `denyRead`, `allowRead`, and [`credentials.files`](#sandbox-credentials-files) resolve by their prefix:
>
> | Prefix            | Meaning                                                                                | Example                                                                   |
> | :---------------- | :------------------------------------------------------------------------------------- | :------------------------------------------------------------------------ |
> | `/`               | Absolute path from filesystem root                                                     | `/tmp/build` stays `/tmp/build`                                           |
> | `~/`              | Relative to home directory                                                             | `~/.kube` becomes `$HOME/.kube`                                           |
> | `./` or no prefix | Relative to the project root for project settings, or to `~/.claude` for user settings | `./output` in `.claude/settings.json` resolves to `<project-root>/output` |
>
> The `//path` prefix for absolute paths also works. If you use single-slash `/path` expecting project-relative resolution, switch to `./path`. This syntax differs from [Read and Edit permission rules](/docs/en/permissions#read-and-edit), which use `//path` for absolute and `/path` for project-relative: sandbox filesystem paths use standard conventions, so `/tmp/build` is an absolute path.
>
> Claude Code strips a trailing slash from a directory path, so `~/.aws` and `~/.aws/` match the same directory. Before v2.1.224, Claude Code passed the trailing slash through to the sandbox, and Claude could still read or write paths under a `denyRead` or `denyWrite` entry written with one.
>
> Claude Code also removes a trailing `/**`, so `~/build/**` and `~/build` cover the same directory. Whether a wildcard such as `*` works depends on which list the entry is in and on the platform:
>
> * **`allowWrite` and `denyWrite`**: on macOS, wildcards work. On Linux and WSL2, the sandbox mounts concrete paths, so Claude Code skips an entry that contains `*`, `?`, or `[` once the trailing `/**` is removed, and that entry has no effect. Claude Code adds the paths from your `Edit` permission rules to these lists, so the same limit applies to them, and the **Config** tab of `/sandbox` warns about `Edit` and `Read` permission rules that contain wildcards.
> * **`denyRead` and `allowRead`**: wildcards work on every platform. On Linux and WSL2, Claude Code expands a read entry to the concrete paths it matches, which it doesn't do for the write lists.
>
> ### `sandbox.filesystem.allowWrite`
>
> Add paths where sandboxed commands can write, beyond the working directory and the session temp directory. Use it when a subprocess such as `kubectl` or a build tool needs to write outside the project.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of path strings, using the [sandbox path prefixes](#sandbox-path-prefixes)
> * **Default**: unset, so sandboxed commands can write only to the working directory and the session temp directory
>
> This lets a build write under `/tmp/build` and lets `kubectl` update your kubeconfig:
>
> ```json settings.json
> {
>   "sandbox": {
>     "filesystem": {
>       "allowWrite": ["/tmp/build", "~/.kube"]
>     }
>   }
> }
> ```
>
> Entries merge across every settings file: user, project, local, and managed paths combine rather than replace each other, and Claude Code adds the paths from your `Edit(...)` allow permission rules. An `allowWrite` entry can't lift a [protected path](/docs/en/sandboxing#protected-paths).
>
> ### `sandbox.filesystem.denyWrite`
>
> Block sandboxed commands from writing to specific paths, including paths inside a directory that is otherwise writable.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of path strings, using the [sandbox path prefixes](#sandbox-path-prefixes)
> * **Default**: unset
>
> This keeps sandboxed commands from changing system configuration or installing binaries:
>
> ```json settings.json
> {
>   "sandbox": {
>     "filesystem": {
>       "denyWrite": ["/etc", "/usr/local/bin"]
>     }
>   }
> }
> ```
>
> Entries merge across every settings file, and Claude Code adds the paths from your `Edit(...)` deny permission rules.
>
> ### `sandbox.filesystem.denyRead`
>
> Block sandboxed commands from reading specific paths, such as credential files that the default read policy would otherwise expose. To protect a credential file and keep it usable through the sandbox proxy, see [`sandbox.credentials`](#sandbox-credentials) instead.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of path strings, using the [sandbox path prefixes](#sandbox-path-prefixes)
> * **Default**: unset, so sandboxed commands keep the [default read access](/docs/en/sandboxing#filesystem-isolation), which includes credential files such as `~/.aws/credentials`
>
> ```json settings.json
> {
>   "sandbox": {
>     "filesystem": {
>       "denyRead": ["~/.aws/credentials"]
>     }
>   }
> }
> ```
>
> Entries merge across every settings file, and Claude Code adds the paths from your `Read(...)` deny permission rules. When [`filesystem.disabled`](#sandbox-filesystem-disabled) is `true`, Claude Code doesn't enforce these entries.
>
> ### `sandbox.filesystem.allowRead`
>
> Re-open reading for specific paths inside a region that [`denyRead`](#sandbox-filesystem-denyread) blocks, to build workspace-only read access. An exact or wildcard `denyRead` entry stays blocked inside a broader `allowRead`, as the [overlap table](/docs/en/sandboxing#configure-sandboxing) shows. When a wildcard `denyRead` entry such as `~/**/.env` matches a directory, Claude Code blocks reads of its contents as well. Before v2.1.236 on macOS, Claude Code re-opened the paths a wildcard `denyRead` entry matched wherever a broader `allowRead` entry covered them, and left a matched directory's contents readable.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of path strings, using the [sandbox path prefixes](#sandbox-path-prefixes)
> * **Default**: unset
>
> This blocks reads of your home directory except the project itself:
>
> ```json settings.json
> {
>   "sandbox": {
>     "filesystem": {
>       "denyRead": ["~/"],
>       "allowRead": ["."]
>     }
>   }
> }
> ```
>
> Place a `.` entry in project settings: it resolves to the project root there and to `~/.claude` in user settings. Entries merge across every settings file unless [`allowManagedReadPathsOnly`](#sandbox-filesystem-allowmanagedreadpathsonly) is set.
>
> ### `sandbox.filesystem.allowManagedReadPathsOnly`
>
> Honor only the [`allowRead`](#sandbox-filesystem-allowread) entries that come from managed settings, so developers can't re-open read access to paths your organization blocked. `denyRead` entries still merge from every settings file.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code honors only the `allowRead` entries from managed settings
>   * `false`: `allowRead` entries merge from every settings file
> * **Default**: `false`
>
> This blocks reads of the home directory, re-opens `~/work`, and stops developers from re-opening anything else:
>
> ```json managed-settings.json
> {
>   "sandbox": {
>     "filesystem": {
>       "denyRead": ["~/"],
>       "allowRead": ["~/work"],
>       "allowManagedReadPathsOnly": true
>     }
>   }
> }
> ```
>
> See [Keep developers from widening the policy](/docs/en/sandboxing#keep-developers-from-widening-the-policy).
>
> ### `sandbox.filesystem.disabled`
>
> Skip filesystem isolation while keeping network isolation. Sandboxed commands get unrestricted read and write access to the host filesystem, and their network egress stays confined to [`network.allowedDomains`](#sandbox-network-alloweddomains). Use it when you sandbox to control where commands connect rather than what they write. Requires Claude Code v2.1.216 or later.
>
> * **Scope**: [`User or managed`](#scopes). When managed settings configure `sandbox.filesystem` at all, or list a `sandbox.credentials.files` entry with `"mode": "deny"`, only managed settings can set it.
> * **Type**: Boolean
>   * `true`: Claude Code skips filesystem isolation and keeps network isolation
>   * `false`: filesystem isolation stays on
> * **Default**: `false`, so filesystem isolation stays on
>
> This leaves the filesystem open and confines network egress to GitHub and npm:
>
> ```json settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "filesystem": {
>       "disabled": true
>     },
>     "network": {
>       "allowedDomains": ["github.com", "*.npmjs.org"]
>     }
>   }
> }
> ```
>
> With the layer off, Claude Code doesn't enforce `denyRead` or `credentials.files` `deny` entries, while `credentials.envVars` entries and applied `mask` entries keep working. [`autoAllowBashIfSandboxed`](#sandbox-autoallowbashifsandboxed) still defaults to `true`, so set it to `false` to keep prompting. See [Disable filesystem isolation](/docs/en/sandboxing#disable-filesystem-isolation) for the full list of sources that can set it and what changes when isolation is off. Requires Claude Code v2.1.216 or later.
>
> ### `sandbox.ignoreViolations`
>
> Silence sandbox violation reports for paths you expect a command to probe and be refused, such as a tool that checks `/etc/hosts` on startup, so those denials don't show up as violations or in what Claude sees. The sandbox still blocks the access; only the report is suppressed. Keys are substrings to match against the command, with `*` matching every command, and values are substrings of the violation to ignore for that command, such as a filesystem path.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object mapping a command substring to an array of violation substrings, usually paths
> * **Default**: unset, so every violation is reported
>
> ```json settings.json
> {
>   "sandbox": {
>     "ignoreViolations": {
>       "*": ["/etc/hosts"]
>     }
>   }
> }
> ```
>
> ### `sandbox.enableWeakerNestedSandbox`
>
> Run the Linux sandbox inside an unprivileged Docker container, where bubblewrap can't mount a fresh `/proc`. Instead the inner sandbox bind-mounts the container's existing `/proc`, which exposes process information that a fresh mount would hide. This reduces security; use it only when the outer container already provides the isolation you need.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: the inner sandbox bind-mounts the container's existing `/proc` instead of mounting a fresh one
>   * `false`: the sandbox mounts a fresh `/proc`, which doesn't work in an unprivileged Docker container
> * **Default**: `false`
>
> ```json settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "enableWeakerNestedSandbox": true
>   }
> }
> ```
>
> Linux and WSL2 only. See [Bubblewrap fails to start inside a container](/docs/en/sandboxing#troubleshooting).
>
> ### `sandbox.enableWeakerNetworkIsolation`
>
> Let sandboxed commands on macOS reach the system TLS trust service, `com.apple.trustd.agent`. Go-based tools such as `gh`, `gcloud`, and `terraform` need it to verify TLS certificates when you use [`network.httpProxyPort`](#sandbox-network-httpproxyport) with a MITM proxy and a custom CA. This reduces security by opening a potential data exfiltration path through the trust service.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: sandboxed commands on macOS can reach `com.apple.trustd.agent`
>   * `false`: sandboxed commands on macOS can't reach the system TLS trust service
> * **Default**: `false`
>
> ```json settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "enableWeakerNetworkIsolation": true
>   }
> }
> ```
>
> If you don't use a MITM proxy, list the failing tools in [`excludedCommands`](#sandbox-excludedcommands) instead; see [Go-based CLIs fail TLS verification on macOS](/docs/en/sandboxing#troubleshooting).
>
> ### `sandbox.allowAppleEvents`
>
> Let sandboxed commands on macOS send Apple Events, which `open`, `osascript`, and tools that open URLs in a browser need; without it they fail with error `-600`. This removes code-execution isolation: sandboxed commands can launch other applications unsandboxed with no user prompt, and can send AppleScript commands to running applications such as Terminal, subject to the per-app macOS automation-consent prompt (TCC).
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: Boolean
>   * `true`: sandboxed commands on macOS can send Apple Events
>   * `false`: sandboxed commands on macOS can't send Apple Events, so `open` and `osascript` fail with error `-600`
> * **Default**: `false`
>
> ```json settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "allowAppleEvents": true
>   }
> }
> ```
>
> To keep isolation and still run one such tool, add it to [`excludedCommands`](#sandbox-excludedcommands) instead. See [Apple Events on macOS](/docs/en/sandboxing#security-limitations).
>
> ### `sandbox.ripgrep`
>
> Point the sandbox at a ripgrep binary of your own instead of the one Claude Code uses, for example when your platform needs a differently built `rg`.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: object with `command`, the path to the ripgrep binary, and optional `args`, an array of arguments to prepend
> * **Default**: unset, so the sandbox uses the same ripgrep binary as Claude Code. That is the bundled binary unless you set [`USE_BUILTIN_RIPGREP`](/docs/en/env-vars) to `0`
>
> ```json settings.json
> {
>   "sandbox": {
>     "ripgrep": {
>       "command": "/usr/local/bin/rg"
>     }
>   }
> }
> ```
>
> ### `sandbox.bwrapPath`
>
> Point the sandbox at a bubblewrap binary installed outside `PATH`, such as a vendored copy on an air-gapped host. Claude Code uses the path both for the startup dependency check and when it wraps each sandboxed command.
>
> * **Scope**: [`Managed`](#scopes). Claude Code reads it only from managed settings so that a user, project, or local file can't point the sandbox at a different binary.
> * **Type**: string, an absolute path; Claude Code drops a relative path and falls back to `PATH` lookup
> * **Default**: unset, so Claude Code finds `bwrap` on `PATH`
>
> ```json managed-settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "bwrapPath": "/opt/admin/bwrap"
>   }
> }
> ```
>
> Linux and WSL2 only.
>
> ### `sandbox.socatPath`
>
> Point the sandbox network proxy at a `socat` binary installed outside `PATH`.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: string, an absolute path; Claude Code drops a relative path and falls back to `PATH` lookup
> * **Default**: unset, so Claude Code finds `socat` on `PATH`
>
> ```json managed-settings.json
> {
>   "sandbox": {
>     "enabled": true,
>     "socatPath": "/opt/admin/socat"
>   }
> }
> ```
>
> Linux and WSL2 only.
>
> ### `sandbox.credentials`
>
> Declare the credential files and environment variables to [protect from sandboxed commands](/docs/en/sandboxing#protect-credentials). Each entry names a file `path` or a variable `name` and a `mode`: `deny` hides the credential inside the sandbox, and `mask` shows sandboxed commands a placeholder while the [sandbox proxy](/docs/en/sandboxing#mask-credentials) substitutes the real value on outbound requests. Claude Code protects only the entries you list; there is no built-in credential deny list. Requires Claude Code v2.1.187 or later.
>
> * **Scope**: [`Any file`](#scopes). `deny` entries merge from every scope, and Claude Code honors `mask` entries, `allowPlaintextInject`, `awsPairs`, and `sigv4` only from user settings, managed settings, and the `--settings` flag.
> * **Type**: object with `files`, `envVars`, `allowPlaintextInject`, `awsPairs`, and `sigv4`
> * **Default**: unset, so no credentials are protected
>
> This hides your AWS credentials file and removes `GITHUB_TOKEN` from sandboxed commands:
>
> ```json settings.json
> {
>   "sandbox": {
>     "credentials": {
>       "files": [{ "path": "~/.aws/credentials", "mode": "deny" }],
>       "envVars": [{ "name": "GITHUB_TOKEN", "mode": "deny" }]
>     }
>   }
> }
> ```
>
> The `deny` file protection is part of the filesystem layer, so it doesn't apply when you [disable filesystem isolation](/docs/en/sandboxing#disable-filesystem-isolation); the environment variable protection still does. Requires Claude Code v2.1.187 or later.
>
> #### Invalid credential entries in managed settings
>
> When a managed `sandbox.credentials` entry fails validation, Claude Code keeps protecting the credential where it can:
>
> * An entry in `files` or `envVars` that still has a valid `path` or `name` and a `mode` of `mask` or `deny`, such as one whose `extract` pattern has no capturing group, is degraded to `mode: "deny"` with a warning, so the credential stays blocked, not masked, until you fix the entry. A degraded `files` entry pins [`filesystem.disabled`](/docs/en/sandboxing#disable-filesystem-isolation) like an explicit `deny` entry, and the warning notes that its read block isn't enforced if managed settings turn filesystem isolation off.
> * An entry with an unknown `mode` or an invalid `path` or `name` is stripped.
> * Each case warns; whether an entry is degraded or stripped, the remaining valid entries are still enforced, and a wholly invalid `credentials` value is dropped while the rest of `sandbox` still applies.
>
> Applies in v2.1.191 and later; before v2.1.221, every invalid entry was stripped. For the other managed keys with per-field handling, see [Invalid entries in managed settings](/docs/en/managed-settings#invalid-entries-in-managed-settings).
>
> ### `sandbox.credentials.files`
>
> Protect credential files or directories from sandboxed commands. With `"mode": "deny"`, Claude Code blocks reads of the path inside the sandbox, the same read block as [`sandbox.filesystem.denyRead`](#sandbox-filesystem-denyread). With `"mode": "mask"`, sandboxed commands on Linux and WSL2 read a sentinel copy of the file, and the sandbox proxy substitutes the real value on outbound requests to that entry's `injectHosts`; on macOS the file is unreadable inside the sandbox instead. Requires Claude Code v2.1.187 or later, and `"mode": "mask"` requires v2.1.221 or later.
>
> * **Scope**: [`Any file`](#scopes). Claude Code drops `mask` entries from project `.claude/settings.json` and local `.claude/settings.local.json`.
> * **Type**: array of objects, each with `path` and a `mode` of `"deny"` or `"mask"`, plus the optional [mask fields for files](#mask-fields-for-files)
> * **Default**: unset, so no credential files are protected
>
> This hides your AWS credentials file and masks the `gh` hosts file, substituting the real value only on requests to `api.github.com`:
>
> ```json settings.json
> {
>   "sandbox": {
>     "credentials": {
>       "files": [
>         { "path": "~/.aws/credentials", "mode": "deny" },
>         { "path": "~/.config/gh/hosts.yml", "mode": "mask", "injectHosts": ["api.github.com"] }
>       ]
>     }
>   }
> }
> ```
>
> Paths use the same [prefixes](#sandbox-path-prefixes) as the `sandbox.filesystem.*` settings, and Claude Code merges the arrays from every settings scope. `mask` substitution runs only through the sandbox proxy, so set [`sandbox.network.tlsTerminate`](#sandbox-network-tlsterminate), or [`allowPlaintextInject`](#sandbox-credentials-allowplaintextinject) for plain-HTTP test networks. `mask` applies to a single file, so list each credential file individually. Claude Code accepts but ignores the `mask` fields on a `deny` entry. [Mask credential files](/docs/en/sandboxing#mask-credential-files) covers which settings sources are honored and when an entry falls back to `deny`. Requires Claude Code v2.1.187 or later; `mask` entries require v2.1.221 or later.
>
>
>
>
>
>
>
> #### Mask fields for files
>
> A `mask` entry accepts these optional fields. Without `extract` or `decode`, Claude Code replaces the entire file content with one sentinel. On macOS with filesystem isolation on, Claude Code applies a `mask` entry as `deny` before `extract` or `decode` runs; see [Mask credential files](/docs/en/sandboxing#mask-credential-files).
>
> | Field              | Type                                                                                                               | What it does                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
> | :----------------- | :----------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `extract`          | string, a regular expression with at least one capturing group                                                     | Mask only the text captured by group 1 of each match, so the rest of the file stays parseable. With `decode` also set, Claude Code checks each capture as a possible JWT instead of replacing it outright. Requires v2.1.221 or later                                                                                                                                                                                                                                                                                                                               |
> | `onExtractNoMatch` | `"warn"`, `"deny"`, or `"error"`; default `"warn"`                                                                 | What happens when `extract` or `decode` finds nothing to mask. `warn` leaves the file readable as-is inside the sandbox, `deny` makes it unreadable, and `error` stops sandbox setup until you fix the configuration. Claude Code treats `deny` as `error` when the read block wouldn't be enforced, because you [disable filesystem isolation](/docs/en/sandboxing#disable-filesystem-isolation) or a [`sandbox.filesystem.allowRead`](#sandbox-filesystem-allowread) entry re-opens the path. Requires v2.1.221 or later; the `decode` case requires v2.1.224 or later |
> | `decode`           | the string `"jwt"`                                                                                                 | Find JSON Web Tokens (JWTs) in the file, with a built-in pattern or with `extract` when set, verify each candidate, and replace it with a structurally valid fake token, so code inside the sandbox that decodes the token keeps working. When no candidate verifies, `onExtractNoMatch` governs the outcome. Requires v2.1.224 or later                                                                                                                                                                                                                            |
> | `maskClaims`       | array of strings, at least one claim name; requires `decode`                                                       | Mask only the named top-level payload claims inside each verified JWT and rebuild the token around the modified payload, so the other claims stay readable. When no named claim matches, `onExtractNoMatch` governs the outcome. Requires v2.1.224 or later                                                                                                                                                                                                                                                                                                         |
> | `maskDuplicates`   | Boolean, default `false`                                                                                           | Also replace verbatim copies of each masked value elsewhere in the file, such as a secret pasted into a comment. Claude Code matches raw substrings, so reserve it for long, high-entropy secrets. Consulted only when `extract` or `decode` is set. Requires v2.1.221 or later                                                                                                                                                                                                                                                                                     |
> | `injectHosts`      | array of strings, each a host that [`sandbox.network.allowedDomains`](#sandbox-network-alloweddomains) also admits | Narrow the hosts where the sandbox proxy substitutes the real value. When unset, the proxy substitutes it on requests to every host in `sandbox.network.allowedDomains`. Requires v2.1.221 or later                                                                                                                                                                                                                                                                                                                                                                 |
>
> This masks only the `oauth_token` value in the `gh` hosts file, replaces every other copy of it in the file, makes the file unreadable if the pattern matches nothing, and substitutes the real token only on requests to `api.github.com`:
>
> ```json settings.json
> {
>   "sandbox": {
>     "credentials": {
>       "files": [
>         {
>           "path": "~/.config/gh/hosts.yml",
>           "mode": "mask",
>           "extract": "oauth_token:\\s*(\\S+)",
>           "maskDuplicates": true,
>           "onExtractNoMatch": "deny",
>           "injectHosts": ["api.github.com"]
>         }
>       ]
>     }
>   }
> }
> ```
>
> ### `sandbox.credentials.envVars`
>
> Protect environment variables from sandboxed commands. With `"mode": "deny"`, Claude Code removes the variable from the environment of sandboxed commands. With `"mode": "mask"`, sandboxed commands see a per-session sentinel value, and the sandbox proxy substitutes the real value on outbound requests to that entry's `injectHosts`, so tools such as `gh` and `npm` keep authenticating without ever holding the real credential. Requires Claude Code v2.1.187 or later, and `"mode": "mask"` requires v2.1.199 or later.
>
> * **Scope**: [`Any file`](#scopes). Claude Code drops `mask` entries from project `.claude/settings.json` and local `.claude/settings.local.json`.
> * **Type**: array of objects, each with `name` and a `mode` of `"deny"` or `"mask"`, plus the optional [mask fields for environment variables](#mask-fields-for-environment-variables)
> * **Default**: unset, so no environment variables are protected
>
> This removes `NPM_TOKEN` from sandboxed commands and masks `GITHUB_TOKEN`, substituting the real value only on requests to `api.github.com`:
>
> ```json settings.json
> {
>   "sandbox": {
>     "credentials": {
>       "envVars": [
>         { "name": "NPM_TOKEN", "mode": "deny" },
>         { "name": "GITHUB_TOKEN", "mode": "mask", "injectHosts": ["api.github.com"] }
>       ]
>     }
>   }
> }
> ```
>
> The `name` must start with a letter or underscore and contain only letters, digits, and underscores. Claude Code merges the arrays from every settings scope, and `deny` takes precedence when the same variable appears with both modes. `mask` substitution runs only through the sandbox proxy, so set [`sandbox.network.tlsTerminate`](#sandbox-network-tlsterminate), or [`allowPlaintextInject`](#sandbox-credentials-allowplaintextinject) for plain-HTTP test networks; see [Mask environment variables](/docs/en/sandboxing#mask-environment-variables). Claude Code accepts but ignores the `mask` fields on a `deny` entry. Requires Claude Code v2.1.187 or later; `mask` entries require v2.1.199 or later.
>
>
>
>
>
>
> #### Mask fields for environment variables
>
> A `mask` entry accepts these optional fields. Without `extract` or `decode`, Claude Code replaces the entire value with one sentinel. `extract` and `decode` can't be combined on the same entry.
>
> | Field              | Type                                                                                                               | What it does                                                                                                                                                                                                                                                                                                                                                                                      |
> | :----------------- | :----------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
> | `extract`          | string, a regular expression with at least one capturing group                                                     | Mask only the text captured by group 1 of each match, such as the password inside a `DATABASE_URL` connection string, so the rest of the value stays parseable. Requires v2.1.224 or later                                                                                                                                                                                                        |
> | `onExtractNoMatch` | `"warn"`, `"deny"`, or `"error"`; default `"warn"`. On an entry with `decode`, only `"warn"` is accepted           | What happens when `extract` matches nothing. `warn` passes the variable through unmasked, `deny` unsets it inside the sandbox, and `error` stops sandbox setup until you fix the configuration. Requires v2.1.224 or later                                                                                                                                                                        |
> | `decode`           | the string `"jwt"`                                                                                                 | Verify the whole value is a JWT and replace it with a structurally valid fake token, so code inside the sandbox that decodes the token keeps working; the proxy substitutes the whole real token on egress. A value that doesn't verify passes through unmasked with a warning. Requires v2.1.224 or later                                                                                        |
> | `maskClaims`       | array of strings, at least one claim name; requires `decode`                                                       | Mask only the named top-level payload claims inside the decoded JWT and rebuild the token around the modified payload, so the other claims stay readable. When no named claim matches, the variable passes through unmasked with a warning. Requires v2.1.224 or later                                                                                                                            |
> | `injectHosts`      | array of strings, each a host that [`sandbox.network.allowedDomains`](#sandbox-network-alloweddomains) also admits | Narrow the hosts where the sandbox proxy substitutes the real value. When unset, the proxy substitutes it on requests to every host in `sandbox.network.allowedDomains`. Write an IPv6 destination as the bare compressed address, such as `"::1"`, not the bracketed form; see [IPv6 destinations in `injectHosts`](/docs/en/sandboxing#ipv6-destinations-in-injecthosts). Requires v2.1.199 or later |
>
> This masks only the password inside `DATABASE_URL`, unsets the variable if the pattern matches nothing, and masks a JWT in `SERVICE_JWT` while leaving every claim except `api_key` readable:
>
> ```json settings.json
> {
>   "sandbox": {
>     "credentials": {
>       "envVars": [
>         {
>           "name": "DATABASE_URL",
>           "mode": "mask",
>           "extract": "://[^:]+:([^@]+)@",
>           "onExtractNoMatch": "deny"
>         },
>         {
>           "name": "SERVICE_JWT",
>           "mode": "mask",
>           "decode": "jwt",
>           "maskClaims": ["api_key"]
>         }
>       ]
>     }
>   }
> }
> ```
>
> ### `sandbox.credentials.allowPlaintextInject`
>
> Allow `mask` substitution on plain HTTP requests as well as TLS-terminated HTTPS. On plain HTTP the upstream identity is unverified and the credential travels in cleartext, so leave this off outside trusted test networks. Requires Claude Code v2.1.199 or later.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code allows `mask` substitution on plain HTTP requests as well as TLS-terminated HTTPS
>   * `false`: Claude Code allows `mask` substitution only on TLS-terminated HTTPS
> * **Default**: `false`
>
> ```json settings.json
> {
>   "sandbox": {
>     "credentials": {
>       "allowPlaintextInject": true
>     }
>   }
> }
> ```
>
> Requires Claude Code v2.1.199 or later.
>
> ### `sandbox.credentials.awsPairs`
>
> Group masked environment variables that form one AWS credential for [SigV4 re-signing](/docs/en/sandboxing#re-sign-aws-requests) when your credential lives in variables with non-standard names. Claude Code links the conventional `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN` trio automatically when you mask their whole values, so you need this key only for other names. Requires Claude Code v2.1.224 or later.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: array of objects, each with `accessKeyIdVar`, `secretAccessKeyVar`, and optionally `sessionTokenVar`, naming `sandbox.credentials.envVars` entries
> * **Default**: unset, so only the conventional trio is paired
>
> This links three custom-named variables into one AWS credential for re-signing:
>
> ```json settings.json
> {
>   "sandbox": {
>     "credentials": {
>       "awsPairs": [
>         {
>           "accessKeyIdVar": "MY_KEY_ID",
>           "secretAccessKeyVar": "MY_SECRET_KEY",
>           "sessionTokenVar": "MY_SESSION_TOKEN"
>         }
>       ]
>     }
>   }
> }
> ```
>
> Each named variable must be a whole-value `mask` entry in [`sandbox.credentials.envVars`](#sandbox-credentials-envvars), without `extract` or `decode`, and can fill only one slot across all pairs.
>
> ### `sandbox.credentials.sigv4`
>
> Choose what the sandbox proxy does with AWS request forms it [can't re-sign](/docs/en/sandboxing#re-sign-aws-requests): `streaming` for aws-chunked streaming uploads, `presigned` for presigned URLs, and `sigv4a` for SigV4A asymmetric signatures. This applies only to requests signed with a masked pair's placeholder access key ID. Requires Claude Code v2.1.224 or later.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: object with `streaming`, `presigned`, and `sigv4a`, each one of:
>   * `"deny"`: the proxy fails the request
>   * `"passthrough"`: the proxy forwards the request signed with the masked placeholder, so the tool receives AWS's own rejection
> * **Default**: unset, so every form is `"deny"`
>
> This forwards streaming uploads instead of failing them at the proxy:
>
> ```json settings.json
> {
>   "sandbox": {
>     "credentials": {
>       "sigv4": {
>         "streaming": "passthrough"
>       }
>     }
>   }
> }
> ```
>
> With `deny`, the proxy fails the request. With `passthrough`, the proxy forwards the request with its signature computed from the masked placeholder, so AWS rejects it and the calling tool receives AWS's own response instead of a proxy error.
>
> ### `sandbox.network`
>
> Control which hosts, ports, and sockets sandboxed commands can reach. The sandbox routes outbound traffic through a proxy that enforces these lists; see [Network isolation](/docs/en/sandboxing#network-isolation) for how the proxy decides and when it prompts.
>
> * **Scope**: [`Any file`](#scopes). `strictAllowlist`, `allowManagedDomainsOnly`, and `tlsTerminate` are read from fewer sources, as their entries say.
> * **Type**: object with the sub-keys below
> * **Default**: unset, so no domains are pre-allowed and the sandbox prompts for each new host
>
> This pre-allows GitHub and npm, blocks `uploads.github.com`, and lets commands bind to localhost:
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "allowedDomains": ["github.com", "*.npmjs.org"],
>       "deniedDomains": ["uploads.github.com"],
>       "allowLocalBinding": true
>     }
>   }
> }
> ```
>
> Claude Code merges the array sub-keys across settings scopes and deduplicates them, so a project can add domains to your user list. `WebFetch(domain:...)` allow and deny [permission rules](/docs/en/sandboxing#permission-rules) feed the same allow and deny lists.
>
> ### `sandbox.network.allowUnixSockets`
>
> List the Unix socket paths sandboxed commands can connect to on macOS. Claude Code ignores this list on Linux and WSL2, where the seccomp filter can't inspect socket paths; use [`allowAllUnixSockets`](#sandbox-network-allowallunixsockets) there instead.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of strings, each a socket path
> * **Default**: unset, so the macOS sandbox blocks every Unix socket
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "allowUnixSockets": ["~/.ssh/agent-socket"]
>     }
>   }
> }
> ```
>
> A socket path can grant broad access: allowing `/var/run/docker.sock`, for example, lets a sandboxed command control the Docker daemon. See [Security limitations](/docs/en/sandboxing#security-limitations).
>
> ### `sandbox.network.allowAllUnixSockets`
>
> Let sandboxed commands connect to every Unix socket. On Linux and WSL2, the sandbox's [seccomp filter](/docs/en/sandboxing#set-up-linux-and-wsl2) blocks `socket(AF_UNIX, ...)` calls, so this is the only way to permit Unix sockets there. When the filter is missing, which `/sandbox` reports on its Dependencies tab, the sandbox doesn't block Unix-socket calls. See [Set up Linux and WSL2](/docs/en/sandboxing#set-up-linux-and-wsl2) for where the filter comes from.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: sandboxed commands can connect to every Unix socket
>   * `false`: the sandbox blocks Unix-socket connections: on macOS except the paths in `allowUnixSockets`, and on Linux and WSL2 through the seccomp filter when it's present
> * **Default**: `false`
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "allowAllUnixSockets": true
>     }
>   }
> }
> ```
>
> On WSL2, `true` also reopens the interop socket that launches Windows binaries such as `cmd.exe` and `powershell.exe`.
>
> ### `sandbox.network.allowLocalBinding`
>
> Let sandboxed commands bind to localhost ports on macOS, for example to start a dev server.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: sandboxed commands can bind to localhost ports on macOS
>   * `false`: sandboxed commands on macOS can't bind to localhost ports
> * **Default**: `false`
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "allowLocalBinding": true
>     }
>   }
> }
> ```
>
> ### `sandbox.network.allowMachLookup`
>
> List additional XPC and Mach service names the macOS sandbox may look up. Tools that communicate over XPC, such as the iOS Simulator or Playwright, need their services listed here.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of strings, each a service name; a single trailing `*` matches a prefix, and `"*"` alone matches every service
> * **Default**: unset
>
> This allows every service under the `com.apple.coresimulator.` prefix:
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "allowMachLookup": ["com.apple.coresimulator.*"]
>     }
>   }
> }
> ```
>
> ### `sandbox.network.allowedDomains`
>
> Pre-allow domains for outbound traffic from sandboxed commands, so the sandbox doesn't prompt for them. Wildcards such as `*.example.com` match subdomains, and an optional `:port` suffix limits an entry to one port; an entry without a port matches every port.
>
> * **Scope**: [`Any file`](#scopes). Only managed settings when [`allowManagedDomainsOnly`](#sandbox-network-allowmanageddomainsonly) is set.
> * **Type**: array of strings, each a domain, wildcard pattern, or IP literal, with an optional `:port` suffix
> * **Default**: unset, so the sandbox prompts the first time a command reaches a new host
>
> This pre-allows GitHub on every port, every npm subdomain, and one API host on port 443 only:
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "allowedDomains": ["github.com", "*.npmjs.org", "api.example.com:443"]
>     }
>   }
> }
> ```
>
> Write IPv6 literals bracketed, with an optional port: `"[::1]"` allows every port and `"[::1]:443"` one port. The bracketed form requires Claude Code v2.1.229 or later. See [IPv6 addresses in domain lists](/docs/en/sandboxing#ipv6-addresses-in-domain-lists).
>
> ### `sandbox.network.deniedDomains`
>
> Block domains for outbound traffic from sandboxed commands, using the same wildcard, port, and IPv6 syntax as [`allowedDomains`](#sandbox-network-alloweddomains). A denied domain stays blocked even when an `allowedDomains` entry matches it too.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of strings, each a domain, wildcard pattern, or IP literal, with an optional `:port` suffix
> * **Default**: unset
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "deniedDomains": ["sensitive.cloud.example.com"]
>     }
>   }
> }
> ```
>
> Claude Code merges this list from every settings source even when `allowManagedDomainsOnly` is set, so a developer can always tighten the deny list. For IPv6 literals, see [IPv6 addresses in domain lists](/docs/en/sandboxing#ipv6-addresses-in-domain-lists).
>
> ### `sandbox.network.strictAllowlist`
>
> Deny sandboxed commands access to hosts outside the allowlist instead of prompting for approval. The allowlist is [`allowedDomains`](#sandbox-network-alloweddomains) plus domains from `WebFetch(domain:...)` allow rules, or only the managed settings entries when [`allowManagedDomainsOnly`](#sandbox-network-allowmanageddomainsonly) is set. Requires Claude Code v2.1.219 or later.
>
> * **Scope**: [`User or managed`](#scopes). A repository can't turn it on or off.
> * **Type**: Boolean
>   * `true`: Claude Code denies sandboxed commands access to hosts outside the allowlist
>   * `false`: unless another trusted settings file sets `true`, Claude Code decides a host outside the allowlist by permission mode instead of denying it outright: it runs the classifier in auto mode, denies in `dontAsk` mode, allows in `bypassPermissions` mode and in plan mode when bypass is available, and otherwise asks you
> * **Default**: `false`
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "strictAllowlist": true
>     }
>   }
> }
> ```
>
> Claude Code enforces this for sandboxed commands only; in-process tools such as `WebFetch` still follow their [permission rules](/docs/en/sandboxing#permission-rules). When any of the honored sources sets it to `true`, it stays on. See [Network isolation](/docs/en/sandboxing#network-isolation). Requires Claude Code v2.1.219 or later.
>
> ### `sandbox.network.allowManagedDomainsOnly`
>
> Lock the network allowlist to what managed settings define. Claude Code then honors only `allowedDomains` and `WebFetch(domain:...)` allow rules from managed settings, ignores domains from user, project, local, and `--settings` settings, and blocks a non-allowed domain automatically instead of prompting.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code honors only `allowedDomains` and `WebFetch(domain:...)` allow rules from managed settings and blocks a non-allowed domain instead of prompting
>   * `false`: domains from user, project, local, and `--settings` settings merge into the allowlist
> * **Default**: `false`
>
> This locks the allowlist to GitHub and npm and ignores any domains developers add:
>
> ```json managed-settings.json
> {
>   "sandbox": {
>     "network": {
>       "allowManagedDomainsOnly": true,
>       "allowedDomains": ["github.com", "*.npmjs.org"]
>     }
>   }
> }
> ```
>
> Denied domains still merge from every source. See [Keep developers from widening the policy](/docs/en/sandboxing#keep-developers-from-widening-the-policy).
>
> ### `sandbox.network.httpProxyPort`
>
> Point the sandbox at your own HTTP proxy instead of the one Claude Code runs. Organizations do this to inspect HTTPS traffic, apply their own filtering rules, or log every request. When unset, Claude Code starts its own proxy for HTTP traffic.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: number, a local TCP port
> * **Default**: unset, so Claude Code runs its own proxy
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "httpProxyPort": 8080
>     }
>   }
> }
> ```
>
> Set [`socksProxyPort`](#sandbox-network-socksproxyport) too if your proxy should carry SOCKS traffic as well; with only one of the two set, Claude Code still runs its own proxy for the other protocol. See [Custom proxy configuration](/docs/en/sandboxing#custom-proxy-configuration).
>
> ### `sandbox.network.socksProxyPort`
>
> Point the sandbox at your own SOCKS5 proxy instead of the one Claude Code runs. When unset, Claude Code starts its own proxy for SOCKS traffic.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: number, a local TCP port
> * **Default**: unset, so Claude Code runs its own proxy
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "socksProxyPort": 8081
>     }
>   }
> }
> ```
>
> See [Custom proxy configuration](/docs/en/sandboxing#custom-proxy-configuration).
>
> ### `sandbox.network.tlsTerminate`
>
> Make the sandbox proxy terminate TLS so it can read the contents of HTTPS requests. This is experimental, and `mask` [credential substitution](/docs/en/sandboxing#mask-credentials) requires it. Set `{}` to generate an ephemeral certificate authority for the session, or set `caCertPath` and `caKeyPath` to use your own.
>
> * **Scope**: [`User or managed`](#scopes). A repository can't switch it on or supply a certificate authority.
> * **Type**: object with optional `caCertPath` and `caKeyPath` strings, each a file path
> * **Default**: unset, so the proxy doesn't terminate or inspect TLS
>
> ```json settings.json
> {
>   "sandbox": {
>     "network": {
>       "tlsTerminate": {}
>     }
>   }
> }
> ```
>
> When more than one honored source sets it, Claude Code uses the value from the highest-precedence source: managed settings, then the `--settings` flag, then user settings. Requires Claude Code v2.1.199 or later.
>
>
> ## Memory and context
>
> Control what Claude Code loads into context, how it compacts, and where it keeps memory and plans. See [Manage context](/docs/en/context-window) and [Memory](/docs/en/memory).
>
> ### `autoCompactEnabled`
>
> Have Claude Code [compact the conversation automatically](/docs/en/context-window#when-your-context-fills-up) when context approaches the limit. Appears in `/config` as **Auto-compact**, and toggling it there writes this key to your user settings.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code compacts the conversation automatically when context approaches the limit
>   * `false`: Claude Code doesn't compact automatically
> * **Default**: `true`
> * **Per-session overrides**: [`DISABLE_AUTO_COMPACT`](/docs/en/env-vars) turns auto-compact off for one session; whichever of the two turns it off, the other can't turn it back on
>
> ```json settings.json
> {
>   "autoCompactEnabled": false
> }
> ```
>
> The manual `/compact` command keeps working while auto-compact is off.
>
> ### `autoCompactWindow`
>
> Set how full the context window gets before Claude Code [compacts automatically](/docs/en/context-window#when-your-context-fills-up).
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: number of tokens, from `100000` to `1000000`. Claude Code caps the value at your model's context window; the [models overview](https://platform.claude.com/docs/en/about-claude/models/overview) lists each model's window
> * **Default**: unset, so Claude Code picks a window tuned for your model
> * **Per-session overrides**: [`--autocompact`](/docs/en/cli-reference#cli-flags) takes precedence over this key for one session, and [`CLAUDE_CODE_AUTO_COMPACT_WINDOW`](/docs/en/env-vars) takes precedence over both
>
> ```json settings.json
> {
>   "autoCompactWindow": 500000
> }
> ```
>
> Set it with the [`/autocompact`](/docs/en/commands#all-commands) command, which writes this key to your user settings. [Set the auto-compact window](/docs/en/model-config#set-the-auto-compact-window) covers how the command, flag, variable, and setting interact.
>
> ### `autoMemoryDirectory`
>
> Store [auto memory](/docs/en/memory#storage-location) in a directory of your choice instead of the per-project default.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, an absolute or `~/`-prefixed directory path
> * **Default**: unset, so Claude Code uses `~/.claude/projects/<project>/memory/`
>
> ```json settings.json
> {
>   "autoMemoryDirectory": "~/my-memory-dir"
> }
> ```
>
> From project or local settings, Claude Code honors this key under the same [workspace trust rule as hooks](/docs/en/permissions#what-runs-before-you-trust-a-folder), since a cloned repository can supply those files.
>
> ### `autoMemoryEnabled`
>
> Turn [auto memory](/docs/en/memory#enable-or-disable-auto-memory) on or off. When `false`, Claude doesn't read from or write to the auto memory directory. You can also toggle it with `/memory` during a session, which writes this key to your user settings.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: the same as unset; auto memory stays on unless something that outranks this key turns it off for the session, such as `--bare`, safe mode, or `CLAUDE_CODE_DISABLE_AUTO_MEMORY`
>   * `false`: Claude doesn't read from or write to the auto memory directory
> * **Default**: `true`
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_AUTO_MEMORY`](/docs/en/env-vars) takes precedence over this key for one session, in either direction
>
> ```json settings.json
> {
>   "autoMemoryEnabled": false
> }
> ```
>
> ### `claudeMd`
>
> Inject CLAUDE.md-style instructions as organization-managed memory without deploying a separate file. Claude Code loads the text as a managed memory entry ahead of user and project CLAUDE.md files.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: string, the text of a CLAUDE.md file; write it as you would the file, Markdown included, with line breaks as `\n`
> * **Default**: unset
>
> This example deploys two rules as a short Markdown list:
>
> ```json managed-settings.json
> {
>   "claudeMd": "# Engineering rules\n\n- Always run make lint before committing.\n- Never push directly to main."
> }
> ```
>
> See [Deploy organization-wide CLAUDE.md](/docs/en/memory#deploy-organization-wide-claude-md).
>
> ### `claudeMdExcludes`
>
> Skip specific `CLAUDE.md` files when Claude Code loads [memory](/docs/en/memory#exclude-specific-claude-md-files). In a large monorepo, use it to skip CLAUDE.md files from other teams that aren't relevant to your work; [Exclude irrelevant CLAUDE.md files](/docs/en/large-codebases#exclude-irrelevant-claude-md-files) in the large-codebases guide walks through that case. Patterns match against absolute file paths.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of strings, each a glob pattern or absolute path
> * **Default**: unset, so Claude Code loads every CLAUDE.md it finds
>
> ```json settings.json
> {
>   "claudeMdExcludes": ["**/vendor/**/CLAUDE.md"]
> }
> ```
>
> Exclusions apply only to user, project, and local memory files; managed policy CLAUDE.md files can't be excluded.
>
>
> ### `fileCheckpointingEnabled`
>
> Have Claude Code snapshot files before each edit so [`/rewind`](/docs/en/checkpointing) can restore them. Appears in `/config` as **Rewind code (checkpoints)**, and toggling it there writes this key to your user settings.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code snapshots files before each edit so `/rewind` can restore them
>   * `false`: Claude Code doesn't snapshot files, so `/rewind` can't restore them
> * **Default**: `true`
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING`](/docs/en/env-vars) turns checkpointing off for one session; whichever of the two turns it off, the other can't turn it back on
>
> ```json settings.json
> {
>   "fileCheckpointingEnabled": false
> }
> ```
>
> In a `-p` run or an Agent SDK session, Claude Code ignores this key. The SDK turns checkpointing on with its `enableFileCheckpointing` option, and a bare `-p` run needs `CLAUDE_CODE_ENABLE_SDK_FILE_CHECKPOINTING=true`. See [File checkpointing in the Agent SDK](/docs/en/agent-sdk/file-checkpointing).
>
> ### `plansDirectory`
>
> Choose where Claude Code stores the plan files it writes in [plan mode](/docs/en/permission-modes#analyze-before-you-edit-with-plan-mode). Claude Code resolves the path relative to the project root and keeps the default when the path resolves outside it.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, a path relative to the project root
> * **Default**: unset, so Claude Code uses `~/.claude/plans`
>
> ```json settings.json
> {
>   "plansDirectory": "./plans"
> }
> ```
>
> ### `skillListingBudgetFraction`
>
> Each turn, Claude sees a [listing of your skills](/docs/en/skills#skill-descriptions-are-cut-short) with their descriptions, and Claude Code caps that listing at a share of the context window. When the listing is over the cap, Claude Code keeps every skill's name but drops the descriptions of the least-used skills, so Claude can still invoke those skills but is less likely to choose one on its own. Raise this key to keep more descriptions visible at the cost of more context per turn.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: number, a fraction greater than `0` and at most `1`
> * **Default**: `0.01`, which reserves 1% of the context window
>
> ```json settings.json
> {
>   "skillListingBudgetFraction": 0.02
> }
> ```
>
> To see how much context the listing uses and which skills contribute most, run `/doctor`.
>
> ### `skillListingMaxDescChars`
>
> Each turn, Claude sees a [listing of your skills](/docs/en/skills#skill-descriptions-are-cut-short) that shows each skill's `description` and `when_to_use` text. This key caps how many characters of that text Claude Code shows per skill; longer text is cut at the cap.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: number of characters, a positive integer
> * **Default**: `1536`
>
> ```json settings.json
> {
>   "skillListingMaxDescChars": 2048
> }
> ```
>
> Raise it to keep long descriptions intact at the cost of more context per turn; lower it to fit more skills under [`skillListingBudgetFraction`](#skilllistingbudgetfraction).
>
> ## Interface and terminal
>
> Change how Claude Code looks and behaves in your terminal: theme, editor mode, status line, spinner, notifications inside the session, and accessibility. See [Terminal configuration](/docs/en/terminal-config).
>
> ### `askUserQuestionTimeout`
>
> Let an unanswered [`AskUserQuestion`](/docs/en/tools-reference) dialog auto-continue after a period of idle time, submitting whatever options you had already selected. Set it when you step away and want Claude to continue without you. With the default, questions wait until you answer them. Requires Claude Code v2.1.200 or later.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: string, one of `"60s"`, `"5m"`, `"10m"`, or `"never"`
> * **Default**: `"never"`
> * **Per-session overrides**: [`CLAUDE_AFK_TIMEOUT_MS`](/docs/en/env-vars) takes precedence over this key for one session
>
> ```json settings.json
> {
>   "askUserQuestionTimeout": "5m"
> }
> ```
>
> Appears in `/config` as **Question auto-continue timeout**, which writes this key to user settings; Claude Code hides the row while managed settings or the `--settings` flag set the key. Requires Claude Code v2.1.200 or later.
>
> ### `autoContinueAtUsageLimit`
>
> After a claude.ai usage limit stops your session, wait in the open session and continue the task automatically after the reset. See [Turn automatic continue off](/docs/en/interactive-mode#turn-automatic-continue-off). Requires Claude Code v2.1.234 or later.
>
> * **Scope**: [`User or managed`](#scopes). Read from user settings, `--settings`, and managed settings only. When none of those sets the key, a project or local settings file that sets it turns the feature off rather than being ignored.
> * **Type**: Boolean
>   * `true`: after a claude.ai usage limit stops your session, Claude Code waits in the open session and continues the task automatically after the reset
>   * `false`: Claude Code doesn't start the wait on its own. You can still [start a wait yourself](/docs/en/interactive-mode#start-a-wait-yourself) from the usage-limit options menu
> * **Default**: `true`
>
> ```json settings.json
> {
>   "autoContinueAtUsageLimit": false
> }
> ```
>
> Appears in `/config` as **Continue automatically at usage limit**, which writes this key to user settings; Claude Code hides the row while managed settings or the `--settings` flag set the key.
>
> ### `autoScrollEnabled`
>
> Follow new output to the bottom of the conversation in [fullscreen rendering](/docs/en/fullscreen). Turn it off to stay where you scrolled while Claude keeps working; permission prompts still scroll into view.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: the conversation follows new output to the bottom
>   * `false`: you stay where you scrolled while Claude keeps working; permission prompts still appear below the transcript
> * **Default**: `true`
>
> ```json settings.json
> {
>   "autoScrollEnabled": false
> }
> ```
>
> Appears in `/config` as **Auto-scroll** when fullscreen rendering is on, which writes this key to user settings.
>
> ### `axScreenReader`
>
> Render screen-reader friendly output: flat text without decorative borders or animations. Screen-reader mode uses the classic renderer, so the `tui` setting has no effect while it is active; attached [background sessions](/docs/en/agent-view) still render fullscreen. Requires Claude Code v2.1.181 or later.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code renders flat text without decorative borders or animations, using the classic renderer
>   * `false`: Claude Code renders normally
> * **Default**: unset, so screen-reader mode is off
> * **Per-session overrides**: [`--ax-screen-reader`](/docs/en/cli-reference#cli-flags) takes precedence over [`CLAUDE_AX_SCREEN_READER`](/docs/en/env-vars), and both take precedence over this key for one session
>
> ```json settings.json
> {
>   "axScreenReader": true
> }
> ```
>
> Requires Claude Code v2.1.181 or later.
>
> ### `companyAnnouncements`
>
> Show your organization's announcements to users at startup. When you list more than one, Claude Code picks one at random for each session; on a person's very first launch it shows the first entry.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of strings
> * **Default**: unset, so no announcement shows
>
> ```json settings.json
> {
>   "companyAnnouncements": [
>     "Welcome to Acme Corp! Review our code guidelines at docs.example.com"
>   ]
> }
> ```
>
> ### `defaultShell`
>
> Choose whether Bash or PowerShell runs the shell commands you type with the [`!` prefix](/docs/en/interactive-mode#shell-mode-with-prefix) in the input box, the ones Claude Code runs directly and adds to the session.
>
> `"powershell"` works only while the [PowerShell tool](/docs/en/tools-reference#powershell-tool) is on. The tool is on by default on Windows without Git Bash, and on Windows with Git Bash for claude.ai and Console accounts. In Amazon Bedrock, Google Cloud's Agent Platform, and Microsoft Foundry sessions, and on macOS, Linux, and WSL, set `CLAUDE_CODE_USE_POWERSHELL_TOOL=1` to turn the tool on. Set that variable to `0` to turn the tool off.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"bash"`: Claude Code runs your `!` commands in Bash
>   * `"powershell"`: Claude Code runs your `!` commands in PowerShell
> * **Default**: `"bash"`, or `"powershell"` on Windows when Bash isn't available
>
> ```json settings.json
> {
>   "defaultShell": "powershell"
> }
> ```
>
> If the shell you name isn't available, Claude Code uses the other one: `"powershell"` falls back to Bash when the PowerShell tool is off, and `"bash"` falls back to PowerShell when Bash isn't installed.
>
> ### `dialogExpiry`
>
> Set the deadline for dialogs Claude Code [forwards to a remote client](/docs/en/remote-control#limitations), such as a Remote Control or SDK host, for the approval dialog for a [held cross-session message](/docs/en/cross-session-messaging#control-inbound-messages), and for the mid-session [Fable 5 usage-credits consent prompt](/docs/en/model-config#fable-5-and-usage-credits) in a session that may have nobody at the terminal. When no answer arrives before the deadline, Claude Code cancels the dialog and continues with its no-action default. Requires Claude Code v2.1.224 or later.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: string, one of `"60s"`, `"5m"`, `"10m"`, or `"never"`, which disables the deadline
> * **Default**: `"5m"`
> * **Per-session overrides**: [`CLAUDE_CODE_USER_DIALOG_TIMEOUT_MS`](/docs/en/env-vars) takes precedence over this key for one session
>
> ```json settings.json
> {
>   "dialogExpiry": "10m"
> }
> ```
>
> Permission prompts and [`AskUserQuestion`](/docs/en/tools-reference#askuserquestion-tool-behavior) questions use their own flows and aren't governed by this deadline. Appears in `/config` as **Dialog expiry**, which writes this key to user settings; the row requires Claude Code v2.1.232 or later, and Claude Code hides it while managed settings or the `--settings` flag set the key.
>
> ### `editorMode`
>
> Choose the key binding mode for the input prompt.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"normal"`: standard key bindings in the prompt input
>   * `"vim"`: vim-style editing with NORMAL, INSERT, and VISUAL modes
> * **Default**: `"normal"`
>
> ```json settings.json
> {
>   "editorMode": "vim"
> }
> ```
>
> Appears in `/config` as **Editor mode**, which writes this key to user settings.
>
> ### `emojiCompletionEnabled`
>
> Show emoji suggestions when you type `:` plus a shortcode in the prompt input, and replace a completed shortcode such as `:heart:` with its emoji. Set it to `false` to turn off both.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code shows emoji suggestions after `:` and replaces a completed shortcode with its emoji
>   * `false`: Claude Code neither suggests emoji nor replaces shortcodes
> * **Default**: `true`
>
> ```json settings.json
> {
>   "emojiCompletionEnabled": false
> }
> ```
>
> See [Emoji shortcodes](/docs/en/interactive-mode#emoji-shortcodes). Requires Claude Code v2.1.217 or later.
>
>
> ### `fileSuggestion`
>
> Run your own command to supply `@` file path autocomplete instead of the built-in file suggestion. The built-in suggestion uses fast filesystem traversal; a large monorepo may do better with project-specific indexing such as a pre-built file index.
>
> * **Scope**: [`Any file`](#scopes). Under the [status line and file suggestion gates](#status-line-and-file-suggestion-gates), Claude Code turns the command off or runs only a managed value, and skips yours without warning.
> * **Type**: object with `type`, always `"command"`, and `command`, the shell command to run
> * **Default**: unset, so Claude Code uses the built-in file suggestion
>
> ```json settings.json
> {
>   "fileSuggestion": {
>     "type": "command",
>     "command": "~/.claude/file-suggestion.sh"
>   }
> }
> ```
>
> After you save this, type `@` followed by part of a path in the prompt: the suggestions come from your command's output.
>
> #### Command input and output
>
> Claude Code runs the command with the same environment variables as [hooks](/docs/en/hooks), including `CLAUDE_PROJECT_DIR`, and stops waiting after five seconds. The command receives JSON on stdin with a `query` field holding what you've typed so far:
>
> ```json
> {"query": "src/comp"}
> ```
>
> Print newline-separated file paths to stdout. Claude Code shows at most 15:
>
> ```text
> src/components/Button.tsx
> src/components/Modal.tsx
> src/components/Form.tsx
> ```
>
> The following script reads the query and hands it to a repository file index:
>
> ```bash
> #!/bin/bash
> query=$(cat | jq -r '.query')
> # Replace your-repo-file-index with your own file search command
> your-repo-file-index --query "$query" | head -20
> ```
>
>
> ### `footerLinksRegexes`
>
> Render extra clickable badges in the footer below the input box when a regex matches turn output: tool results, including file contents and fetched pages, and Claude's own responses. Use it to turn IDs printed by project CLIs, such as review tools and issue trackers, into session links. Requires Claude Code v2.1.176 or later.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: array of objects, each with `type` set to `"regex"`, a `pattern` regex, a `url` template, and an optional `label`; `{name}` placeholders in `url` and `label` are filled from named capture groups in `pattern`
> * **Default**: unset, so no badges render
>
> This example matches issue keys such as `PROJ-1234` and builds each link from the captured key:
>
> ```json settings.json
> {
>   "footerLinksRegexes": [
>     {
>       "type": "regex",
>       "pattern": "\\b(?<key>PROJ-\\d+)\\b",
>       "url": "https://issues.example.com/browse/{key}",
>       "label": "{key}"
>     }
>   ]
> }
> ```
>
> With this configured, when `PROJ-1234` appears in a tool result or in Claude's reply, a `PROJ-1234` badge appears in the footer linking to `https://issues.example.com/browse/PROJ-1234`. Requires Claude Code v2.1.176 or later.
>
> #### Badge constraints
>
> Each entry's URL, label, and badge count are bounded as follows:
>
> | Constraint  | Behavior                                                                                                                                                                                           |
> | :---------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | URL origin  | Captured values are URL-encoded and the constructed URL must share the template's literal origin. A capture can fill a path segment or query value but can't change where the link points          |
> | URL length  | Constructed URLs longer than 2048 characters are dropped                                                                                                                                           |
> | URL scheme  | Must be `https`, `http`, or a recognized editor or workspace deep-link scheme: `vscode`, `vscode-insiders`, `cursor`, `windsurf`, `zed`, `jetbrains`, `idea`, `slack`, `linear`, `notion`, `figma` |
> | Label       | Defaults to the matched text and is truncated to 28 display columns                                                                                                                                |
> | Badge count | At most 5 badges render. The oldest is displaced by newer matches and `/clear` removes them                                                                                                        |
>
> When a turn completes, Claude Code matches each entry's `pattern` regex against the turn output on the main thread, so a slow regex blocks the UI until it finishes. Nested quantifiers such as `(a+)+$` can take exponentially long against certain inputs and freeze the session, so keep each `pattern` linear and avoid nesting `+` or `*`.
>
> Footer badges render alongside a [custom status line](/docs/en/statusline) when one is configured; neither replaces the other. Use a status line for a script-driven row that computes its own content from session data, and footer badges to turn IDs from the conversation into links without a script.
>
> ### `keybindingFlavor`
>
> Choose which convention `Ctrl+W` follows in the prompt input. Set it to `"readline"` to make `Ctrl+W` delete back to the previous whitespace, as Bash does, so a path or a `--flag=value` goes in one press. Requires Claude Code v2.1.238 or later.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"classic"`: `Ctrl+W` deletes the previous word
>   * `"readline"`: `Ctrl+W` deletes back to the previous whitespace
> * **Default**: `"classic"`
>
> ```json settings.json
> {
>   "keybindingFlavor": "readline"
> }
> ```
>
> See [Make editing keys follow readline conventions](/docs/en/interactive-mode#make-ctrl-w-delete-back-to-whitespace) for the per-key behavior.
>
> ### `prefersReducedMotion`
>
> Reduce or turn off interface animations such as the spinner, shimmer, and flash effects. Appears in `/config` as **Reduce motion**.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code reduces or turns off interface animations such as the spinner, shimmer, and flash effects
>   * `false`: the same as unset; Claude Code shows its animations
> * **Default**: `false`
>
> ```json settings.json
> {
>   "prefersReducedMotion": true
> }
> ```
>
> ### `promptSuggestionEnabled`
>
> Show or hide [prompt suggestions](/docs/en/interactive-mode#prompt-suggestions), the grayed-out predictions that appear in your prompt input. Set it to `false`, or turn off **Prompt suggestions** in `/config`, to hide them.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: you see prompt suggestions in your prompt input
>   * `false`: Claude Code hides prompt suggestions
> * **Default**: `true`
> * **Per-session overrides**: [`CLAUDE_CODE_ENABLE_PROMPT_SUGGESTION`](/docs/en/env-vars) takes precedence over this key for one session
>
> ```json settings.json
> {
>   "promptSuggestionEnabled": false
> }
> ```
>
> Prompt suggestions need a claude.ai or Console account with telemetry on. On Amazon Bedrock, Google Cloud's Agent Platform, and Microsoft Foundry, or with telemetry turned off, such as by [`DISABLE_TELEMETRY`](/docs/en/env-vars), this key has no effect and only `CLAUDE_CODE_ENABLE_PROMPT_SUGGESTION=1` turns them on.
>
> ### `respectGitignore`
>
> Control whether the `@` file picker leaves out files that match `.gitignore` patterns. Appears in `/config` as **Respect .gitignore in file picker**.
>
> * **Scope**: [`Any file`](#scopes). When no settings file sets it, Claude Code falls back to `respectGitignore` in `~/.claude.json`, which the `/config` toggle writes.
> * **Type**: Boolean
>   * `true`: the `@` file picker leaves out files that match `.gitignore` patterns
>   * `false`: the `@` file picker includes files that match `.gitignore` patterns
> * **Default**: `true`
>
> ```json settings.json
> {
>   "respectGitignore": false
> }
> ```
>
> ### `respondToBashCommands`
>
> Choose whether Claude responds after you run a shell command with the [`!` prefix](/docs/en/interactive-mode#shell-mode-with-prefix) in the input box. By default, Claude Code adds the command's output to the conversation and Claude replies to it. Set this key to `false` to add the output to context without a reply, so you can run several commands and ask about them together. Requires Claude Code v2.1.186 or later.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code adds the command's output to the conversation and Claude replies to it
>   * `false`: Claude Code adds the output to context without a reply
> * **Default**: `true`
>
> ```json settings.json
> {
>   "respondToBashCommands": false
> }
> ```
>
> See [Shell mode with `!` prefix](/docs/en/interactive-mode#shell-mode-with-prefix). Requires Claude Code v2.1.186 or later.
>
> ### `showClearContextOnPlanAccept`
>
> When Claude finishes a plan in [plan mode](/docs/en/permission-modes#review-and-approve-a-plan), it shows an approval menu. Planning can use a lot of context, so this key adds a first option to that menu, **Yes, clear context and …**, that approves the plan, clears the conversation context, and starts implementing from the plan alone. The rest of the label names the permission mode the session continues in, and shows how much of your context the planning used.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: the plan approval menu gets a first option, **Yes, clear context and …**, that approves the plan and clears the conversation context
>   * `false`: the plan approval menu shows no clear-context option
> * **Default**: `false`
>
> ```json settings.json
> {
>   "showClearContextOnPlanAccept": true
> }
> ```
>
> ### `showTurnDuration`
>
> Show or hide the turn duration message after each response, such as "Cooked for 1m 6s". Appears in `/config` as **Show turn duration**.
>
> * **Scope**: [`Any file`](#scopes). A value in `~/.claude.json` from an older version applies when no settings file sets it.
> * **Type**: Boolean
>   * `true`: you see the turn duration message after each response
>   * `false`: Claude Code hides the turn duration message
> * **Default**: `true`
>
> ```json settings.json
> {
>   "showTurnDuration": false
> }
> ```
>
> ### `spellcheck`
>
> Underline misspelled words in the prompt input as you type, using a spell checker you install. Claude Code checks only the text in the input box. [Check spelling as you type](/docs/en/interactive-mode#check-spelling-as-you-type) covers installing aspell, hunspell, or ispell and what the checker covers. Requires Claude Code v2.1.235 or later.
>
> * **Scope**: [`User or managed`](#scopes). The block from the highest tier that sets it applies as a whole.
> * **Type**: object with `enabled` (Boolean), `checker` (`"aspell"`, `"hunspell"`, `"ispell"`, or `"auto"`), `language` (string, passed to the checker as its dictionary name), and `color` (string, a terminal color name, `#rrggbb`, `rgb(r,g,b)`, `ansi256(n)`, or `ansi:<name>`)
> * **Default**: unset, so spell checking is off; `checker` defaults to `"auto"`, the first of the three found on `PATH`; `language` defaults to the checker's own dictionary; `color` defaults to the theme's error color
>
> ```json settings.json
> {
>   "spellcheck": { "enabled": true, "language": "en_GB" }
> }
> ```
>
> ### `spinnerTipsEnabled`
>
> While Claude works, the spinner line rotates through short tips about Claude Code features, such as "Use Plan Mode to prepare for a complex request before making changes. Press Shift+Tab twice to enable." Set this key to `false` to hide them. Appears in `/config` as **Show tips**.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: you see tips in the spinner while Claude is working
>   * `false`: Claude Code hides spinner tips
> * **Default**: `true`
>
> ```json settings.json
> {
>   "spinnerTipsEnabled": false
> }
> ```
>
> ### `spinnerTipsOverride`
>
> Replace or extend the [spinner tips](#spinnertipsenabled), the short hints Claude Code rotates through while Claude works, with your own strings, such as a team reminder to run a review skill. Set `excludeDefault` to `true` and list at least one tip to show only your tips; when it's `false` or absent, or `tips` is empty, Claude Code keeps the built-in tips and adds yours.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object with a `tips` array of strings and an optional `excludeDefault` Boolean
> * **Default**: unset, so Claude Code shows only the built-in tips
>
> This example replaces the built-in tips with a single tip of your own:
>
> ```json settings.json
> {
>   "spinnerTipsOverride": {
>     "excludeDefault": true,
>     "tips": ["Run /review before opening a PR"]
>   }
> }
> ```
>
> ### `spinnerVerbs`
>
> While a turn is in progress, the spinner shows a rotating verb such as "Accomplishing", "Architecting", or "Baking". Use this key to add your own verbs to that rotation or replace the built-in list with yours.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object with a `verbs` array of strings and `mode`, one of:
>   * `"append"`: Claude Code adds your verbs to the built-in set
>   * `"replace"`: Claude Code shows only your verbs
> * **Default**: unset, so Claude Code uses the built-in verbs
>
> This example adds two verbs to the built-in set:
>
> ```json settings.json
> {
>   "spinnerVerbs": {
>     "mode": "append",
>     "verbs": ["Pondering", "Crafting"]
>   }
> }
> ```
>
> In `"replace"` mode with an empty `verbs` array, Claude Code keeps the built-in verbs.
>
> ### `statusLine`
>
> Run your own command to render a [status line](/docs/en/statusline) below the prompt with context such as the model, cost, or git branch. Optional fields adjust spacing, add periodic re-runs, and hide the built-in vim mode indicator when your script renders `vim.mode` itself.
>
> * **Scope**: [`Any file`](#scopes). When [`allowManagedHooksOnly`](#allowmanagedhooksonly) is on, or [`disableAllHooks`](#disableallhooks) is set outside managed settings, only the managed settings value runs.
> * **Type**: object with `type` set to `"command"` and a `command` string, plus optional `padding` as a number of characters, `refreshInterval` as a number of seconds, minimum `1`, and `hideVimModeIndicator` as a Boolean
> * **Default**: unset, so no status line
>
> This example prints the model name and context usage, and adds two characters of horizontal spacing:
>
> ```json settings.json
> {
>   "statusLine": {
>     "type": "command",
>     "command": "jq -r '\"[\\(.model.display_name)] \\(.context_window.used_percentage // 0)% context\"'",
>     "padding": 2
>   }
> }
> ```
>
> The example needs [`jq`](https://jqlang.org/) installed and runs in a shell. For PowerShell and Git Bash equivalents, see [Windows configuration](/docs/en/statusline#windows-configuration); for the full setup, see [Manually configure a status line](/docs/en/statusline#manually-configure-a-status-line).
>
> ### `subagentStatusLine`
>
> When Claude runs [subagents](/docs/en/sub-agents), Claude Code lists them in a task display below the prompt, one row per subagent showing `name · description · token count`. This key lets you run your own command to rewrite those rows, for example to show each subagent's context usage as a percentage. On each refresh, Claude Code sends the visible rows as one JSON object on stdin, with a `tasks` array carrying each subagent's `id`, `name`, `status`, `model`, `tokenCount`, and more, and replaces the row for each `id` you write back as a `{"id", "content"}` line. Rows you don't write back keep the default rendering.
>
> * **Scope**: [`Any file`](#scopes). When [`allowManagedHooksOnly`](#allowmanagedhooksonly) is on, or [`disableAllHooks`](#disableallhooks) is set outside managed settings, only the managed settings value runs.
> * **Type**: object with `type` set to `"command"` and a `command` string
> * **Default**: unset, so Claude Code renders the default rows
>
> ```json settings.json
> {
>   "subagentStatusLine": {
>     "type": "command",
>     "command": "jq -c '.tasks[] | {id, content: \"\\(.name): \\(.tokenCount) tokens\"}'"
>   }
> }
> ```
>
> See [Subagent status lines](/docs/en/statusline#subagent-status-lines).
>
> ### `syntaxHighlightingDisabled`
>
> Claude Code colors code by language in the diffs, code blocks, and file previews it shows in the terminal, with its built-in highlighter; no plugin or language server is involved. Set this key to `true` to show them as plain text instead, for example if the colors clash with your terminal theme or slow a screen reader.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code turns off syntax highlighting in diffs, code blocks, and file previews
>   * `false`: Claude Code highlights syntax
> * **Default**: `false`
>
> ```json settings.json
> {
>   "syntaxHighlightingDisabled": true
> }
> ```
>
> ### `terminalProgressBarEnabled`
>
> Some terminals can show a progress indicator on the tab or in the taskbar for the program running in them. While Claude is working, Claude Code reports an in-progress state to the terminal and clears it when the turn ends, so you can see from another tab or window whether Claude is still busy. It does so only in terminals that support the indicator: ConEmu, Ghostty 1.2.0 or later, and iTerm2 3.6.6 or later. Set this key to `false` to stop reporting it. Appears in `/config` as **Terminal progress bar**.
>
> * **Scope**: [`Any file`](#scopes). A value in `~/.claude.json` from an older version applies when no settings file sets it.
> * **Type**: Boolean
>   * `true`: you see the terminal progress bar in terminals that support it
>   * `false`: Claude Code hides the terminal progress bar
> * **Default**: `true`
>
> ```json settings.json
> {
>   "terminalProgressBarEnabled": false
> }
> ```
>
> ### `terminalTitleFromRename`
>
> Claude Code sets your terminal tab's title. By default it uses a title it generates from the conversation, and once you give the session a [name](/docs/en/sessions#name-your-sessions) with `/rename` or `--name`, the tab shows that name instead. Set this key to `false` to keep the generated title on the tab even after you name the session. The name itself still applies, so `/resume <name>` and the session picker find it.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: the terminal tab title shows the session name you set
>   * `false`: the tab keeps the title Claude Code generates from your conversation
> * **Default**: `true`
>
> ```json settings.json
> {
>   "terminalTitleFromRename": false
> }
> ```
>
> To stop Claude Code from updating the terminal title at all, set [`CLAUDE_CODE_DISABLE_TERMINAL_TITLE`](/docs/en/env-vars) to `1` instead.
>
> ### `theme`
>
> Pick the color theme for the interface. Appears in `/config` as **Theme**.
>
> * **Scope**: [`Any file`](#scopes). A value in `~/.claude.json` from an older version applies when no settings file sets it.
> * **Type**: string, one of:
>   * `"auto"`: matches your terminal's light or dark background
>   * `"dark"`: the dark theme
>   * `"light"`: the light theme
>   * `"dark-daltonized"`: the dark theme with colorblind-friendly colors
>   * `"light-daltonized"`: the light theme with colorblind-friendly colors
>   * `"dark-ansi"`: the dark theme using only your terminal's ANSI color palette
>   * `"light-ansi"`: the light theme using only your terminal's ANSI color palette
>   * `"custom:<slug>"` or `"custom:<plugin-name>:<slug>"`: a custom theme from `~/.claude/themes/` or a plugin
> * **Default**: `"dark"`
>
> ```json settings.json
> {
>   "theme": "light-daltonized"
> }
> ```
>
> See [Create a custom theme](/docs/en/terminal-config#create-a-custom-theme).
>
> ### `tui`
>
> Choose the terminal UI renderer. Use `"fullscreen"` for the flicker-free [alt-screen renderer](/docs/en/fullscreen) with virtualized scrollback, or `"default"` for the classic main-screen renderer. Running `/tui fullscreen` or `/tui default` writes this key for you.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"default"`: the classic main-screen renderer
>   * `"fullscreen"`: the flicker-free alt-screen renderer with virtualized scrollback
> * **Default**: unset, so Claude Code [picks the renderer for you](/docs/en/fullscreen#fullscreen-by-default)
> * **Per-session overrides**: [`CLAUDE_CODE_NO_FLICKER`](/docs/en/env-vars) and [`CLAUDE_CODE_DISABLE_ALTERNATE_SCREEN`](/docs/en/env-vars) take precedence over this key for one session: `CLAUDE_CODE_NO_FLICKER=1` turns fullscreen on, and `CLAUDE_CODE_NO_FLICKER=0` or `CLAUDE_CODE_DISABLE_ALTERNATE_SCREEN=1` turns it off; when both are set, Claude Code turns it off
>
> ```json settings.json
> {
>   "tui": "fullscreen"
> }
> ```
>
> Under tmux `-CC` or over SSH to Windows, Claude Code keeps the classic renderer unless you set `CLAUDE_CODE_NO_FLICKER=1`. Background sessions opened from [agent view](/docs/en/agent-view) always use the fullscreen renderer regardless of this setting.
>
> ### `verbose`
>
> By default, the transcript collapses each tool call to a short summary, such as the command Claude ran and a line count of its output, and you press `Ctrl+O` to switch the whole transcript to the expanded view when you want the details. Set this key to `true` to show every tool call's full input and output inline as it happens, which is useful when you're debugging a hook, an MCP server, or a long shell command. Appears in `/config` as **Verbose output**.
>
> * **Scope**: [`Any file`](#scopes). A value in `~/.claude.json` from an older version applies when no settings file sets it.
> * **Type**: Boolean
>   * `true`: you see full tool output
>   * `false`: you see truncated summaries of tool output
> * **Default**: `false`
> * **Per-session overrides**: [`--verbose`](/docs/en/cli-reference#cli-flags) takes precedence over this key for one session
>
> ```json settings.json
> {
>   "verbose": true
> }
> ```
>
> A [`viewMode`](#viewmode) value or a sticky `/focus` selection overrides this key every session.
>
> ### `viewMode`
>
> Set the transcript view Claude Code starts in: `"default"`, `"verbose"`, or `"focus"`. When set, it overrides both the sticky `/focus` selection and the [`verbose`](#verbose) setting.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"default"`: the normal transcript with truncated tool output
>   * `"verbose"`: the transcript with full tool output
>   * `"focus"`: only your last prompt, a one-line summary of tool calls with edit diffstats, and the final response. Focus view needs the [fullscreen renderer](#tui)
> * **Default**: unset, so the `verbose` setting and your last `/focus` choice apply
> * **Per-session overrides**: [`--verbose`](/docs/en/cli-reference#cli-flags) takes precedence over this key for one session
>
> ```json settings.json
> {
>   "viewMode": "focus"
> }
> ```
>
> ### `vimInsertModeRemaps`
>
> Map two-key INSERT-mode sequences to Escape in [vim editor mode](/docs/en/interactive-mode#vim-editor-mode). Each key is exactly two printable characters typed in sequence, and `"<Esc>"` is the only supported target; Claude Code ignores other entries. Requires Claude Code v2.1.208 or later.
>
> * **Scope**: [`User or managed`](#scopes). A repository can't remap your keystrokes.
> * **Type**: object mapping a two-character sequence to `"<Esc>"`
> * **Default**: unset
>
> ```json settings.json
> {
>   "vimInsertModeRemaps": {
>     "jj": "<Esc>"
>   }
> }
> ```
>
> Has no effect unless `editorMode` is `"vim"`. See [Remap INSERT-mode key sequences](/docs/en/interactive-mode#remap-insert-mode-key-sequences). Requires Claude Code v2.1.208 or later.
>
> ### `voice`
>
> Turn on [voice dictation](/docs/en/voice-dictation) and choose how the dictation key behaves. Claude Code writes this object for you when you run `/voice`.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object with `enabled` as a Boolean, `autoSubmit` as a Boolean that applies in hold mode only, and `mode`, one of:
>   * `"hold"`: you hold the dictation key while speaking and release it to stop
>   * `"tap"`: you tap the key once to start recording and again to send
> * **Default**: unset, so dictation is off; when `enabled` is `true` and `mode` is unset, Claude Code uses `"hold"`
>
> This example turns dictation on and makes the key tap once to start recording and again to send:
>
> ```json settings.json
> {
>   "voice": {
>     "enabled": true,
>     "mode": "tap"
>   }
> }
> ```
>
> `autoSubmit` sends the prompt when you release the key in hold mode. Voice dictation requires a claude.ai account.
>
> ### `voiceEnabled`
>
> Turn voice dictation on with the single Boolean form that predates the `voice` object. When both are set, `voice.enabled` applies.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: voice dictation is on when you're logged in with a claude.ai account and your organization's policy allows voice, unless `voice.enabled` is set
>   * `false`: voice dictation is off, unless `voice.enabled` is set
> * **Default**: unset
>
> ```json settings.json
> {
>   "voiceEnabled": true
> }
> ```
>
> ### `wheelScrollAccelerationEnabled`
>
> Accelerate mouse-wheel scroll speed during fast scrolls in [fullscreen rendering](/docs/en/fullscreen#mouse-wheel-scrolling). Set it to `false` for a constant scroll rate per wheel notch. Requires Claude Code v2.1.174 or later.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code accelerates mouse-wheel scroll speed during fast scrolls
>   * `false`: Claude Code scrolls at a constant rate per wheel notch
> * **Default**: `true`
>
> ```json settings.json
> {
>   "wheelScrollAccelerationEnabled": false
> }
> ```
>
> Requires Claude Code v2.1.174 or later.
>
> ## Git and attribution
>
> Control the attribution Claude Code adds to commits and pull requests and how it works with git.
>
>
> ### `attribution`
>
> Customize the attribution Claude Code adds to git commits and pull requests. Commits get a [git trailer](https://git-scm.com/docs/git-interpret-trailers) such as `Co-Authored-By` by default; pull request descriptions get plain text. Set each part separately with the sub-keys below.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object with `commit` and `pr` strings and a `sessionUrl` Boolean
> * **Default**: unset, so Claude Code uses the standard attribution shown under each sub-key
>
> This example replaces the commit attribution, removes pull request attribution, and drops the session link:
>
> ```json settings.json
> {
>   "attribution": {
>     "commit": "Generated with AI\n\nCo-Authored-By: AI <ai@example.com>",
>     "pr": "",
>     "sessionUrl": false
>   }
> }
> ```
>
> To hide all attribution, set [`commit`](#attribution-commit) and [`pr`](#attribution-pr) to empty strings and [`sessionUrl`](#attribution-sessionurl) to `false`. Once you set `commit` or `pr`, Claude Code ignores the deprecated `includeCoAuthoredBy` setting and uses its default text for whichever of the two you left unset.
>
> ### `includeCoAuthoredBy`
>
> Use [`attribution`](#attribution) instead, which replaces this key and lets you change or hide the commit trailer, the pull request text, and the session link separately. Claude Code still honors `includeCoAuthoredBy: false` from settings files that predate `attribution`, but ignores it once you set `attribution.commit` or `attribution.pr`.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: the same as unset; Claude Code adds the commit trailer and the pull request attribution text
>   * `false`: Claude Code omits both the commit trailer and the pull request attribution text, unless `attribution` sets `commit` or `pr`, in which case the [`attribution`](#attribution) rules apply
> * **Default**: `true`
>
> ```json settings.json
> {
>   "includeCoAuthoredBy": false
> }
> ```
>
> To hide all attribution today, set [`attribution.commit`](#attribution-commit) and [`attribution.pr`](#attribution-pr) to empty strings and [`attribution.sessionUrl`](#attribution-sessionurl) to `false`.
>
> ### `includeGitInstructions`
>
> At session start, Claude Code adds two git-related pieces to Claude's prompt: its built-in instructions for how to write commits and pull requests, in the Bash tool's description, and a git status snapshot of your repository in the system prompt, meaning the current branch, the main branch, `git status` output, and recent commits. Set this key to `false` to leave both out, for example when you use your own git workflow skills.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code includes its built-in commit and pull request workflow instructions and the git status snapshot. Cloud sessions never include the snapshot
>   * `false`: Claude Code leaves both out
> * **Default**: `true`
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_GIT_INSTRUCTIONS`](/docs/en/env-vars) takes precedence over this key for one session
>
> ```json settings.json
> {
>   "includeGitInstructions": false
> }
> ```
>
> ### `prUrlTemplate`
>
> Point the PR links Claude Code renders, in the footer badge and in tool-result summaries, at an internal code-review tool instead of `github.com`. Claude Code substitutes `{host}`, `{owner}`, `{repo}`, `{number}`, and `{url}` from the `gh`-reported PR URL. The [GitLab merge request badge](/docs/en/interactive-mode#gitlab-merge-requests) keeps its GitLab URL.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, a URL template using any of the five placeholders
> * **Default**: unset, so links use the `gh`-reported URL
>
> ```json settings.json
> {
>   "prUrlTemplate": "https://reviews.example.com/{owner}/{repo}/pull/{number}"
> }
> ```
>
> Claude Code applies the template only to the links it renders itself; a PR number Claude writes in a message, such as `#123`, stays as Claude wrote it. A URL that doesn't have the `/pull/<number>` shape is left unchanged.
>
> ### `attribution.commit`
>
> Set the attribution text Claude Code adds to git commits, including any trailers. Set it to an empty string to hide commit attribution.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string
> * **Default**: unset, so Claude Code adds `Co-Authored-By: <model name> <noreply@anthropic.com>`, where the model name reflects the active model for the session, such as `Claude Sonnet 5`, or `Claude` alone when the session's model isn't a public model
>
> This example replaces the default trailer with a custom line and a custom `Co-Authored-By` trailer:
>
> ```json settings.json
> {
>   "attribution": {
>     "commit": "Generated with AI\n\nCo-Authored-By: AI <ai@example.com>"
>   }
> }
> ```
>
> ### `attribution.pr`
>
> Set the attribution text Claude Code adds to pull request descriptions. Set it to an empty string to hide pull request attribution.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string
> * **Default**: unset, so Claude Code adds `🤖 Generated with [Claude Code](https://claude.com/claude-code)`
>
> ```json settings.json
> {
>   "attribution": {
>     "pr": ""
>   }
> }
> ```
>
> ### `attribution.sessionUrl`
>
> Choose whether Claude Code appends the claude.ai session link when it commits or opens a pull request from a [cloud](/docs/en/claude-code-on-the-web) or [Remote Control](/docs/en/remote-control) session. Claude Code adds the link as a `Claude-Session` trailer on commits and as a link in pull request descriptions. Set it to `false` to omit the link.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code appends the claude.ai session link when it commits or opens a pull request from a cloud or Remote Control session
>   * `false`: Claude Code omits the link
> * **Default**: `true`
>
> ```json settings.json
> {
>   "attribution": {
>     "sessionUrl": false
>   }
> }
> ```
>
>
>
> ## Plugins and skills
>
> Enable plugins, register marketplaces, restrict which plugin sources an organization allows, and control which skills load. For installing and building plugins, see [Plugins](/docs/en/plugins).
>
> ### `disableBundledSkills`
>
> Turn off the [skills](/docs/en/skills) and workflows included with Claude Code. Claude Code removes bundled skills and workflows entirely, while built-in commands such as `/init` stay typable but are hidden from the model.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code removes bundled skills and workflows and hides built-in commands such as `/init` from the model
>   * `false`: bundled skills load
> * **Default**: unset, so bundled skills load
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_BUNDLED_SKILLS`](/docs/en/env-vars) set to `1` turns bundled skills off for one session; whichever of the two turns them off, the other can't turn them back on
>
> ```json settings.json
> {
>   "disableBundledSkills": true
> }
> ```
>
> Skills from plugins, `.claude/skills/`, and `.claude/commands/` are unaffected. `/doctor` stays typable like the built-in commands; to hide it, set [`DISABLE_DOCTOR_COMMAND`](/docs/en/env-vars) instead.
>
> ### `disableSkillShellExecution`
>
> Turn off inline shell execution for `` !`...` `` and ` ```! ` blocks in [skills](/en/skills) and custom commands from user, project, plugin, or additional-directory sources. Claude Code replaces each command with `[shell command execution disabled by policy]` instead of running it.
>
> * **Scope**: [`Any file`](#scopes). A `true` in managed settings can't be overridden by `false` elsewhere.
> * **Type**: Boolean
>   * `true`: Claude Code replaces each inline shell command with `[shell command execution disabled by policy]` instead of running it
>   * `false`: inline shell runs
> * **Default**: unset, so inline shell runs
>
> ```json settings.json
> {
>   "disableSkillShellExecution": true
> }
> ```
>
> Bundled skills and skills deployed through managed settings are unaffected.
>
> ### `skillOverrides`
>
> Hide or collapse a [skill](/docs/en/skills#override-skill-visibility-from-settings) without editing its `SKILL.md`. Claude Code applies the value under each skill's name to the skill list Claude sees and to your `/` autocomplete.
>
> * **Scope**: [`Any file`](#scopes). The `/skills` menu writes to `.claude/settings.local.json`.
> * **Type**: object mapping skill name to one of:
>   * `"on"`: Claude sees the skill and you can type `/name`
>   * `"name-only"`: Claude sees the skill by name without its description
>   * `"user-invocable-only"`: Claude doesn't see the skill, but you can still type `/name`
>   * `"off"`: Claude doesn't see the skill and `/name` is hidden from autocomplete
> * **Default**: unset, so every skill is `"on"`
>
> This example lists `legacy-context` to Claude by name only and hides `deploy` from Claude and from `/` autocomplete:
>
> ```json settings.json
> {
>   "skillOverrides": {
>     "legacy-context": "name-only",
>     "deploy": "off"
>   }
> }
> ```
>
> `"name-only"` lists the skill to the model without its description, `"user-invocable-only"` hides it from the model but keeps `/name` typable, and `"off"` hides it from both. Overrides don't apply to plugin skills, which you manage through `/plugin`.
>
> ### `syncClaudeAiSkills`
>
> Turn off the download of the [skills you enable on claude.ai](/docs/en/skills#how-synced-skills-behave). Claude Code downloads them into `~/.claude/skills/synced/` when you run it in [non-interactive mode](/docs/en/headless) with the `-p` flag and [`CLAUDE_CODE_SYNC_SKILLS`](/docs/en/env-vars#variables) set. Set `false` to stop that download and hide the skills it already synced. Claude Code honors only `false`: `true` is the same as unset and doesn't turn syncing on.
>
> * **Scope**: [`User, local, or managed`](#scopes). A repository can't turn it off for you.
> * **Type**: Boolean
>   * `false`: Claude Code stops downloading synced skills and hides the ones already in `~/.claude/skills/synced/`. In user or managed settings, it also moves them to `~/.claude/skills/.trash/`
>   * `true`: the same as unset
> * **Default**: unset, so a non-interactive run with `CLAUDE_CODE_SYNC_SKILLS` set downloads the skills
>
> This example keeps a machine from downloading the account's skills, whatever a session sets in its environment:
>
> ```json settings.json
> {
>   "syncClaudeAiSkills": false
> }
> ```
>
> ### `allowedChannelPlugins`
>
> Choose which [channel](/docs/en/channels) plugins can push messages into sessions in your organization. When you set it, Claude Code uses your list in place of the default Anthropic allowlist; each entry names a plugin and the marketplace it comes from.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: array of objects, each with `marketplace` and `plugin` strings
> * **Default**: unset, so Claude Code uses the default Anthropic allowlist
>
> This example turns channels on and allows only the Telegram plugin from the official Anthropic marketplace:
>
> ```json managed-settings.json
> {
>   "channelsEnabled": true,
>   "allowedChannelPlugins": [
>     { "marketplace": "claude-plugins-official", "plugin": "telegram" }
>   ]
> }
> ```
>
> An empty array blocks every channel plugin. This key takes effect once channels pass the [`channelsEnabled`](#channelsenabled) gate for the account: on Team and Enterprise plans, and on Console accounts with managed settings, that means `channelsEnabled: true`. See [Restrict which channel plugins can run](/docs/en/channels#restrict-which-channel-plugins-can-run).
>
> ### `blockedMarketplaces`
>
> Block plugin marketplace sources for your organization. Claude Code checks the blocklist on marketplace add and on plugin install, update, refresh, and auto-update, so a marketplace someone added before you set the policy can't be used to fetch plugins either. Blocked sources are checked before download, so they never touch the filesystem.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: array of marketplace source objects, in the same forms as [`strictKnownMarketplaces`](#allowed-source-types)
> * **Default**: unset, so no marketplace is blocked
>
> This example blocks one GitHub repository as a marketplace source:
>
> ```json managed-settings.json
> {
>   "blockedMarketplaces": [
>     { "source": "github", "repo": "untrusted/plugins" }
>   ]
> }
> ```
>
> A `github` entry may use the [owner-wildcard form](#owner-wildcards) `"owner/*"` to block every repository under that GitHub owner, which requires Claude Code v2.1.223 or later. Add `{ "source": "skills-dir" }` to stop Claude Code loading [`@skills-dir` plugins](/docs/en/plugins-reference#skills-directory-plugins) from `~/.claude/skills/` without restricting any marketplace. See [Managed marketplace restrictions](/docs/en/plugin-marketplaces#managed-marketplace-restrictions).
>
> ### `channelsEnabled`
>
> Allow [channels](/docs/en/channels) for your organization. On claude.ai Team and Enterprise plans, Claude Code blocks channels until you set this to `true`. For [Anthropic Console](/docs/en/authentication#claude-console-authentication) accounts that authenticate with an API key, channels are allowed by default. If your organization deploys managed settings, Claude Code blocks channels on those accounts too until you set this key to `true`.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code allows channels for your organization
>   * `false`: the same as unset; whether channels are blocked depends on your plan, as the Default says
> * **Default**: unset; channels are blocked on Team and Enterprise plans and on Console accounts with managed settings, and allowed on Pro and Max plans and on Console accounts without managed settings
>
> ```json managed-settings.json
> {
>   "channelsEnabled": true
> }
> ```
>
> To restrict which plugins can register as channels once they're enabled, set [`allowedChannelPlugins`](#allowedchannelplugins). See [Enterprise controls](/docs/en/channels#enterprise-controls).
>
> ### `disableCommandPluginSources`
>
> Block the [`command` plugin source](/docs/en/plugin-marketplaces#command-sources), which installs a plugin by running a marketplace-declared command on the user's machine. When you set it to `true`, Claude Code never runs the command, doesn't install or update command-sourced plugins, and stops loading the ones already installed. Set it to `false` to allow them explicitly. Whenever it blocks command sources, whether you set it to `true` or leave it unset under [`allowManagedHooksOnly`](#allowmanagedhooksonly), it also blocks marketplace [`headersHelper` commands](/docs/en/plugin-marketplaces#authenticate-archive-downloads), except for a marketplace that managed settings themselves declare. Requires Claude Code v2.1.229 or later, and the `headersHelper` block requires v2.1.238 or later.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code never runs the marketplace-declared command, doesn't install or update command-sourced plugins, and stops loading the ones already installed
>   * `false`: Claude Code allows command-sourced plugins explicitly
> * **Default**: unset, so Claude Code follows [`allowManagedHooksOnly`](#allowmanagedhooksonly): an organization that restricts hook execution to managed settings gets command sources disabled too
>
> ```json managed-settings.json
> {
>   "disableCommandPluginSources": true
> }
> ```
>
> Requires Claude Code v2.1.229 or later.
>
> ### `pluginSuggestionMarketplaces`
>
> Name the marketplaces whose plugins can appear as contextual install suggestions, in spinner tips and pinned at the top of the `/plugin` **Discover** tab. The built-in first-party frontend-design tip is unaffected. Suggestions come from each plugin's `relevance` declaration in its marketplace entry.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: array of marketplace names
> * **Default**: unset, so no marketplace-declared suggestions surface
>
> ```json managed-settings.json
> {
>   "pluginSuggestionMarketplaces": ["acme-corp-plugins"]
> }
> ```
>
> A name takes effect only when the marketplace is registered on the machine and its registered source is also declared in the same managed settings, either as the [`extraKnownMarketplaces`](#extraknownmarketplaces) entry for that name or as an entry of [`strictKnownMarketplaces`](#strictknownmarketplaces). Claude Code ignores a marketplace registered from a different source under an allowlisted name. The official marketplace is exempt from the source requirement: allowlisting its name alone suffices, since that name can only register from the official Anthropic source. See [Suggest plugins by context](/docs/en/plugin-relevance).
>
> ### `pluginTrustMessage`
>
> Add your organization's own text to the plugin trust warning Claude Code shows before installation, for example to confirm that plugins from your internal marketplace are vetted.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: string
> * **Default**: unset, so Claude Code shows the standard warning alone
>
> ```json managed-settings.json
> {
>   "pluginTrustMessage": "All plugins from our marketplace are approved by IT"
> }
> ```
>
> ### `strictKnownMarketplaces`
>
> Restrict which plugin marketplace sources people in your organization can add and install plugins from. Claude Code enforces the allowlist on marketplace add and on plugin install, update, refresh, and auto-update, before any network or filesystem operation, so a marketplace someone added before you set the policy can't be used to fetch plugins once its source no longer matches. Blocked users see an error naming the managed policy.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: array of marketplace source objects; see [Allowed source types](#allowed-source-types)
> * **Default**: unset, so users can add any marketplace. An empty array is a complete lockdown that blocks every marketplace source, including the official Anthropic marketplace
>
> This example allows two GitHub repositories, one pinned to the `v2.0` ref, and one hosted `marketplace.json` URL:
>
> ```json managed-settings.json
> {
>   "strictKnownMarketplaces": [
>     { "source": "github", "repo": "acme-corp/approved-plugins" },
>     { "source": "github", "repo": "acme-corp/security-tools", "ref": "v2.0" },
>     { "source": "url", "url": "https://plugins.example.com/marketplace.json" }
>   ]
> }
> ```
>
> You can also write this key as `allowedMarketplaces`; [Marketplace key aliases](#marketplace-key-aliases) describes how Claude Code treats the alias and which version accepts it. This key is a policy gate: it controls what users may add but registers nothing. To restrict and pre-register in one file, see [Combine with `extraKnownMarketplaces`](#combine-with-extraknownmarketplaces). For the user-facing view, see [Managed marketplace restrictions](/docs/en/plugin-marketplaces#managed-marketplace-restrictions).
>
> #### Allowed source types
>
> Each entry below shows one allowlist entry per source type and the fields it accepts. Most types match exactly; `hostPattern` and `pathPattern` match by regex, and `github` entries can use an [owner wildcard](#owner-wildcards).
>
> | Source        | Example entry                                                                                                                   | Fields                                                                                         |
> | :------------ | :------------------------------------------------------------------------------------------------------------------------------ | :--------------------------------------------------------------------------------------------- |
> | `github`      | `{ "source": "github", "repo": "acme-corp/plugins", "ref": "main", "path": "marketplace" }`                                     | `repo` required; `ref` is a branch or tag; `path` is a subdirectory                            |
> | `git`         | `{ "source": "git", "url": "https://gitlab.example.com/tools/plugins.git", "ref": "production" }`                               | `url` required; `ref` and `path` as for `github`                                               |
> | `url`         | `{ "source": "url", "url": "https://plugins.example.com/marketplace.json", "headers": { "Authorization": "Bearer ${TOKEN}" } }` | `url` required; `headers` adds HTTP headers for authenticated access                           |
> | `npm`         | `{ "source": "npm", "package": "@acme-corp/claude-plugins" }`                                                                   | `package` required, the npm package that contains `marketplace.json`                           |
> | `file`        | `{ "source": "file", "path": "/opt/acme-corp/plugins/marketplace.json" }`                                                       | `path` required, the absolute path to a `marketplace.json` file                                |
> | `directory`   | `{ "source": "directory", "path": "/opt/acme-corp/approved-marketplaces" }`                                                     | `path` required, the absolute path to a directory containing `.claude-plugin/marketplace.json` |
> | `hostPattern` | `{ "source": "hostPattern", "hostPattern": "^github\\.example\\.com$" }`                                                        | `hostPattern` required, a regex matched against the marketplace host                           |
> | `pathPattern` | `{ "source": "pathPattern", "pathPattern": "^/opt/approved/" }`                                                                 | `pathPattern` required, a regex matched against the `path` of `file` and `directory` sources   |
> | `skills-dir`  | `{ "source": "skills-dir" }`                                                                                                    | No fields. Opts the `~/.claude/skills/` plugin scan back in                                    |
>
> Three source types carry rules beyond the table:
>
> * **`url`**: a URL marketplace downloads only the `marketplace.json` file, and Claude Code doesn't fetch plugin files by relative path from that server, so its plugins must use a [plugin source](/docs/en/plugin-marketplaces#plugin-sources) other than a relative path, such as an archive URL, which can be on the same host. For plugins with relative paths, use a Git-based marketplace instead. See [Plugins with relative paths fail in URL-based marketplaces](/docs/en/plugin-marketplaces#plugins-with-relative-paths-fail-in-url-based-marketplaces).
> * **`hostPattern`**: use it to allow every marketplace on an internal GitHub Enterprise or GitLab server without listing each repository. Claude Code matches `github` sources against `github.com`, takes the hostname from `url` sources, and takes it from `git` sources depending on the [git URL](https://git-scm.com/docs/git-clone#_git_urls)'s form:
>
>   * A URL with a scheme, such as `https://` or `ssh://`: the hostname in the URL.
>   * An SSH address without a scheme, in git's `user@host:path` form, such as `git@git.example.com:tools/plugins.git`: the host between `@` and `:`, which is the host git connects to.
>   * Any other form without a scheme: no host, so no `strictKnownMarketplaces` `hostPattern` entry matches it. For a `blockedMarketplaces` `hostPattern`, Claude Code takes a host from a wider set of forms, so a blocklist entry can still match such a form. Before v2.1.234, a `strictKnownMarketplaces` `hostPattern` also matched some forms that git doesn't treat as SSH addresses.
>
>   `file` and `directory` sources have no host and never match a `hostPattern` entry.
> * **`pathPattern`**: use it to allow filesystem marketplaces alongside `hostPattern` entries for network sources. `".*"` allows every local path; a narrower pattern such as `"^/opt/approved/"` restricts to a directory.
>
> Any allowlist, even an empty one, also stops Claude Code loading [`@skills-dir` plugins](/docs/en/plugins-reference#skills-directory-plugins) from `~/.claude/skills/`. Add the `{ "source": "skills-dir" }` entry to keep loading them; the entry has no meaning outside this key and `blockedMarketplaces`.
>
> #### Owner wildcards
>
> A `github` entry whose `repo` value is `"<owner>/*"` matches every repository under that GitHub owner. Owner wildcards require Claude Code v2.1.223 or later and work only in `strictKnownMarketplaces` and `blockedMarketplaces`. Everywhere else a `github` source appears, such as `extraKnownMarketplaces` or `/plugin marketplace add`, the `repo` value must name a single repository. Before v2.1.223, Claude Code compared the entry literally, so an allowlist entry matched no repository and a blocklist entry blocked nothing; single-repository entries are enforced on every version.
>
> This entry allows any marketplace repository in the `acme-corp` organization:
>
> ```json managed-settings.json
> {
>   "strictKnownMarketplaces": [
>     { "source": "github", "repo": "acme-corp/*" }
>   ]
> }
> ```
>
> Only the whole repository-name position can be a wildcard. Claude Code compares entries such as `*`, `*/plugins`, or `acme-corp/tools-*` literally, so they match no repository.
>
> The matching rules differ between the two settings:
>
> | Rule                      | `strictKnownMarketplaces`                                                                                                                                             | `blockedMarketplaces`                                                           |
> | ------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
> | Matching source spellings | `owner/repo` form only. A git URL that clones the same repository doesn't match                                                                                       | Any spelling, including git URLs that resolve to the same github.com repository |
> | Owner case                | Case-sensitive, like exact-entry matching                                                                                                                             | Case-insensitive                                                                |
> | `ref`                     | Follows the exact-entry rules: an entry with a `ref` matches only sources with that exact ref, and an entry without one matches only sources that don't specify a ref | An entry without a `ref` blocks all refs of the repositories it matches         |
> | `path`                    | Looser than the exact-entry rules: an entry with a `path` requires that exact value, while an entry without one matches any path inside the repository                | An entry without a `path` blocks all paths of the repositories it matches       |
>
> #### Exact matching
>
> For every source type except owner-wildcard `github` entries and the regex-matched `hostPattern` and `pathPattern` entries, Claude Code allows a user's addition only when the marketplace source matches an entry exactly. For the git-based sources `github` and `git`, exact matching includes the optional fields:
>
> * The `repo` or `url` must match exactly
> * The `ref` field must match exactly, or both must be undefined
> * The `path` field must match exactly, or both must be undefined
>
> For example, Claude Code treats each pair below as two different sources:
>
> * `{ "source": "github", "repo": "acme-corp/plugins" }` and `{ "source": "github", "repo": "acme-corp/plugins", "ref": "main" }`
> * `{ "source": "github", "repo": "acme-corp/plugins", "path": "marketplace" }` and `{ "source": "github", "repo": "acme-corp/plugins" }`
>
> #### Allow only the official marketplace
>
> To allow the official Anthropic marketplace and nothing else, list its repository:
>
> ```json managed-settings.json
> {
>   "strictKnownMarketplaces": [
>     { "source": "github", "repo": "anthropics/claude-plugins-official" }
>   ]
> }
> ```
>
> With this entry, Claude Code keeps an already-registered official marketplace available and, on a fresh machine, registers the marketplace automatically the first time you start Claude Code interactively. Automatic registration most commonly misses:
>
> * Non-interactive environments that run before the machine's first interactive launch.
> * Machines where Claude Code already ran interactively under a policy that blocked the marketplace, such as the empty-array lockdown. Claude Code records the blocked attempt and doesn't retry after the policy changes.
>
> On these machines, add the marketplace to [`extraKnownMarketplaces`](#extraknownmarketplaces) in the same `managed-settings.json` so Claude Code registers it automatically, or run `claude plugin marketplace add anthropics/claude-plugins-official`.
>
> #### Combine with `extraKnownMarketplaces`
>
> The two keys do different jobs. This table compares them:
>
> | Aspect            | `strictKnownMarketplaces`                | `extraKnownMarketplaces`                                                                             |
> | ----------------- | ---------------------------------------- | ---------------------------------------------------------------------------------------------------- |
> | Purpose           | Organizational policy enforcement        | Team convenience                                                                                     |
> | Settings file     | Managed settings only                    | Any settings file                                                                                    |
> | Behavior          | Blocks non-allowlisted additions         | Registers missing marketplaces                                                                       |
> | When enforced     | Before network and filesystem operations | Immediately from user or managed settings; after the workspace trust dialog for a repository's files |
> | Can be overridden | No, highest precedence                   | Yes, by higher-precedence settings                                                                   |
> | Source format     | Direct source object                     | Named marketplace with a nested `source` object                                                      |
>
> To both restrict and pre-register a marketplace for all users, set both in `managed-settings.json`:
>
> ```json managed-settings.json
> {
>   "strictKnownMarketplaces": [
>     { "source": "github", "repo": "acme-corp/plugins" }
>   ],
>   "extraKnownMarketplaces": {
>     "acme-tools": {
>       "source": { "source": "github", "repo": "acme-corp/plugins" }
>     }
>   }
> }
> ```
>
> With only `strictKnownMarketplaces` set, users can still add an allowed marketplace themselves with `/plugin marketplace add`. The official Anthropic marketplace is the only one Claude Code registers automatically, and only when the allowlist allows it. [Allow only the official marketplace](#allow-only-the-official-marketplace) lists the machines it misses.
>
> ### `strictPluginOnlyCustomization`
>
> Block skills, agents, hooks, and MCP servers from user and project sources, so they can come only from plugins or managed settings. Combine it with [`strictKnownMarketplaces`](#strictknownmarketplaces) to control the full customization supply chain: the marketplace allowlist controls which plugins users can install.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: `true` to lock all four kinds of customization, or an array naming the kinds to lock, from `"skills"`, `"agents"`, `"hooks"`, and `"mcp"`
> * **Default**: unset, so nothing is locked
>
> This example locks skills and hooks and leaves agents and MCP servers unlocked:
>
> ```json managed-settings.json
> {
>   "strictPluginOnlyCustomization": ["skills", "hooks"]
> }
> ```
>
> The four sub-key entries below list what each surface blocks and what still loads. Claude Code ignores surface names it doesn't recognize rather than failing the settings file, so you can add new surface names before every client has updated.
>
> ### `strictPluginOnlyCustomization.skills`
>
> Lock the `skills` surface. Claude Code stops loading skills from `~/.claude/skills/` and `.claude/skills/`, custom commands from `~/.claude/commands/` and `.claude/commands/`, skills under `--add-dir` directories, and skills synced from your claude.ai account, and keeps loading plugin skills, bundled skills, and skills in the managed policy directory.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: the string `"skills"` in the [`strictPluginOnlyCustomization`](#strictpluginonlycustomization) array
> * **Default**: not locked
>
> ```json managed-settings.json
> {
>   "strictPluginOnlyCustomization": ["skills"]
> }
> ```
>
> ### `strictPluginOnlyCustomization.agents`
>
> Lock the `agents` surface. Claude Code stops loading agents from `~/.claude/agents/` and `.claude/agents/`, and keeps loading plugin agents, built-in agents, and agents in the managed policy directory.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: the string `"agents"` in the [`strictPluginOnlyCustomization`](#strictpluginonlycustomization) array
> * **Default**: not locked
>
> ```json managed-settings.json
> {
>   "strictPluginOnlyCustomization": ["agents"]
> }
> ```
>
> ### `strictPluginOnlyCustomization.hooks`
>
> Lock the `hooks` surface. Claude Code stops running hooks from user, project, and local `settings.json`, and keeps running plugin hooks and hooks in managed settings.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: the string `"hooks"` in the [`strictPluginOnlyCustomization`](#strictpluginonlycustomization) array
> * **Default**: not locked
>
> ```json managed-settings.json
> {
>   "strictPluginOnlyCustomization": ["hooks"]
> }
> ```
>
> ### `strictPluginOnlyCustomization.mcp`
>
> Lock the `mcp` surface. Claude Code stops loading MCP servers from `~/.claude.json` and `.mcp.json`, and keeps loading plugin MCP servers and [`managed-mcp.json`](/docs/en/managed-mcp) servers.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: the string `"mcp"` in the [`strictPluginOnlyCustomization`](#strictpluginonlycustomization) array
> * **Default**: not locked
>
> ```json managed-settings.json
> {
>   "strictPluginOnlyCustomization": ["mcp"]
> }
> ```
>
> ### `enabledPlugins`
>
> Turn individual [plugins](/docs/en/plugins) on or off, keyed by `plugin-name@marketplace-name`. A plugin with no entry at any scope falls back to its [`defaultEnabled`](/docs/en/plugins-reference#default-enablement) value. When you enable or disable a plugin with `/plugin` or `claude plugin enable`, Claude Code writes this key for you.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object mapping `plugin-name@marketplace-name` to a Boolean
> * **Default**: unset, so each plugin follows its `defaultEnabled` value
>
> This example enables two plugins from the `team-tools` marketplace and disables one from `personal`:
>
> ```json settings.json
> {
>   "enabledPlugins": {
>     "code-formatter@team-tools": true,
>     "deployment-tools@team-tools": true,
>     "experimental-features@personal": false
>   }
> }
> ```
>
> Each scope serves a different purpose:
>
> * **User settings**: your personal plugin preferences
> * **Project settings**: plugins shared with everyone in the repository
> * **Local settings**: per-machine overrides, gitignored when Claude Code saves a setting there
> * **Managed settings**: organization-wide policy. A plugin set to `false` here is blocked from installation at every scope and hidden from the marketplace
>
> Project settings take precedence over user settings, so setting a plugin to `false` in `~/.claude/settings.json` doesn't disable a plugin that the project's `.claude/settings.json` enables. To opt out of a project-enabled plugin on your machine, set it to `false` in `.claude/settings.local.json` instead. Plugins force-enabled by managed settings can't be disabled this way, since managed settings override local settings.
>
> Enabling a plugin from an external source such as a GitHub repository or npm package in a project's `.claude/settings.json` doesn't install it for other people. On every path that loads plugins, Claude Code reports the plugin as not installed until each user [installs it themselves](/docs/en/discover-plugins#configure-team-marketplaces).
>
> ### `extraKnownMarketplaces`
>
> Register additional plugin marketplaces by name, so that people who open the repository, or everyone your managed settings reach, get the marketplace without adding it themselves. Claude Code registers each marketplace it doesn't already know. Whether a plugin that [`enabledPlugins`](#enabledplugins) names from it installs depends on the plugin's source and which file enables it; that entry has the rules.
>
> * **Scope**: [`Any file`](#scopes). Claude Code honors entries in a repository's `.claude/settings.json` or `.claude/settings.local.json` only after you accept the workspace trust dialog for that folder; in a folder you haven't trusted, including a `-p` run there, it ignores them without a message.
> * **Type**: object mapping a marketplace name to an object with a `source` object and an optional `autoUpdate` Boolean
> * **Default**: unset
>
> This example registers a GitHub marketplace and a marketplace from a self-hosted git URL:
>
> ```json settings.json
> {
>   "extraKnownMarketplaces": {
>     "acme-tools": {
>       "source": {
>         "source": "github",
>         "repo": "acme-corp/claude-plugins"
>       }
>     },
>     "security-plugins": {
>       "source": {
>         "source": "git",
>         "url": "https://git.example.com/security/plugins.git"
>       }
>     }
>   }
> }
> ```
>
> [What runs before you trust a folder](/docs/en/permissions#what-runs-before-you-trust-a-folder) compares the trust gate with the other content a repository can supply. You can also write this key as `additionalMarketplaces`; see [Marketplace key aliases](#marketplace-key-aliases).
>
> Set `"autoUpdate": true` alongside `source` to make Claude Code refresh that marketplace and update its installed plugins in the background after startup. When omitted, `claude-plugins-official` and most other official Anthropic marketplaces default to `true`, and third-party marketplaces default to `false`. See [Configure auto-updates](/docs/en/discover-plugins#configure-auto-updates).
>
> When more than one settings file defines a marketplace entry under the same name, Claude Code uses the entry from the [highest-precedence file](/docs/en/settings#settings-precedence) whole. That entry replaces the lower-precedence entry and inherits none of its fields, so a redefinition can't combine one file's `source.headers` credential with a URL another file controls. Before v2.1.228, Claude Code merged same-name entries field by field, so an entry in a higher-precedence file could inherit fields it didn't set, including another file's `headers`.
>
> #### Marketplace source types
>
> The `source` object takes one of these forms:
>
> * **`github`**: a GitHub repository, with `repo`
> * **`git`**: any git URL, with `url`
> * **`url`**: a direct URL to a `marketplace.json` file, with `url` and optional `headers` and `headersHelper` for authenticated access. `headersHelper` names a command that prints headers whose values are too short-lived to list in `headers`, and requires Claude Code v2.1.238 or later
> * **`file`**: a local path to a `marketplace.json` file, with `path`
> * **`directory`**: a local filesystem path, with `path`, for development only
> * **`settings`**: an inline marketplace declared directly in the settings file without a hosted repository, with `name` and `plugins`
>
> The `git` source type works with any git hosting service, including self-hosted GitLab and Bitbucket. Claude Code clones the repository with the same authentication that `git clone` would use on that machine: configured credential helpers or SSH keys. A provider token such as `GITHUB_TOKEN` takes effect only through a credential helper that reads it. See [Private repositories](/docs/en/plugin-marketplaces#private-repositories) for setup details.
>
> For `github` and `git` sources, set `"skipLfs": true` inside the `source` object, alongside `repo` or `url`, to skip Git LFS downloads when Claude Code clones or updates the marketplace repository. LFS pointer files remain as pointers instead of downloading their content. Use this when the repository contains large LFS objects unrelated to plugin content. Requires Claude Code v2.1.153 or later.
>
> For a `url` source, set `headersHelper` inside the `source` object when the credential in `headers` expires and a command has to produce a fresh one. Requires Claude Code v2.1.238 or later. For what the command must print and where Claude Code runs it, see [Write the headersHelper command](/docs/en/plugin-marketplaces#write-the-headershelper-command), and for the cases where Claude Code doesn't run it, see [When Claude Code skips a headersHelper command](/docs/en/plugin-marketplaces#when-claude-code-skips-a-headershelper-command-or-drops-its-output). Once you set `headersHelper` on an `https://` marketplace URL, Claude Code runs the command at two points, reusing one run's output for up to 60 seconds:
>
> * Before each fetch of that marketplace's `marketplace.json`, including a later refresh. Claude Code sends the printed headers with that fetch.
> * Before each plugin archive download on the marketplace URL's origin, meaning the same scheme, host, and port. Claude Code sends the output with that download, and no other download gets the headers.
>
> Claude Code ignores any `headersHelper` set in the `.claude/settings.json` or `.claude/settings.local.json` of a directory you add with [`--add-dir`](/docs/en/permissions#what-runs-before-you-trust-a-folder), on a `url` source and on an inline plugin entry alike, and sends only the fixed `headers` set in that file. [How users accept a headersHelper command](/docs/en/plugin-marketplaces#how-users-accept-a-headershelper-command) covers the other settings files.
>
> Plugins listed in a `settings` source must reference external sources such as GitHub or npm, and the `name` must match the marketplace key. You still enable each plugin separately in `enabledPlugins`. This example declares one plugin inline:
>
> ```json settings.json
> {
>   "extraKnownMarketplaces": {
>     "team-tools": {
>       "source": {
>         "source": "settings",
>         "name": "team-tools",
>         "plugins": [
>           {
>             "name": "code-formatter",
>             "source": {
>               "source": "github",
>               "repo": "acme-corp/code-formatter"
>             }
>           }
>         ]
>       }
>     }
>   }
> }
> ```
>
> A plugin entry under `source: 'settings'` whose own `source` is an [`archive`](/docs/en/plugin-marketplaces#zip-archives) can set `headers` for the archive download. If the value you would put in `headers` is short-lived, such as a token your registry mints on request, set a `headersHelper` command instead. An entry may set both. Both fields require Claude Code v2.1.238 or later.
>
> Claude Code sends the entry's `headers`, and whatever the command prints, with that plugin's archive download and with no other download. Claude Code runs the command only when a user [installs or updates that one plugin by itself](/docs/en/plugin-marketplaces#how-users-accept-a-headershelper-command). Three further rules depend on which file holds the entry:
>
> * **`strict`**: unlike an entry in a marketplace's `marketplace.json`, an entry in settings doesn't need `"strict": false`, because a settings file carries no manifest fields to inline. See [Strict mode](/docs/en/plugin-marketplaces#strict-mode).
> * **Folder trust**: for an entry in a project's `.claude/settings.json` or `.claude/settings.local.json`, Claude Code runs the command only after the user has also [trusted that folder](/docs/en/permissions#what-runs-before-you-trust-a-folder).
> * **Header filter**: Claude Code drops [request-routing and client-identity header names](/docs/en/plugin-marketplaces#when-claude-code-skips-a-headershelper-command-or-drops-its-output) from an entry in a project's `.claude/settings.json` or `.claude/settings.local.json`, because a repository can supply those files. Claude Code applies the same filter to a catalog entry and to an entry in an `--add-dir` directory's settings, and no filter to an entry in your user settings, a `--settings` file, or managed settings.
>
> #### Marketplace key aliases
>
> On Claude Code v2.1.232 or later, you can write `extraKnownMarketplaces` as `additionalMarketplaces` and `strictKnownMarketplaces` as `allowedMarketplaces`. Claude Code treats each alias as follows:
>
> * Earlier versions ignore the alias, so keep the canonical spelling in a file that older versions also read, such as a managed settings file for a fleet with mixed Claude Code versions.
> * In any settings file that accepts the canonical key, Claude Code reads the alias exactly as it reads the canonical key.
> * Claude Code may rewrite `additionalMarketplaces` to `extraKnownMarketplaces` when it updates the file.
> * If you set both spellings in one file, Claude Code uses the canonical value and ignores the alias.
>
> ### `pluginConfigs`
>
> Store the non-sensitive answers you give a plugin's [`userConfig`](/docs/en/plugins-reference#user-configuration) configuration dialog, keyed by plugin ID. Claude Code writes this key to your user settings when you fill in the dialog, so you don't need to edit it by hand. Sensitive options go to the macOS Keychain instead, or to `~/.claude/.credentials.json` on platforms without a supported keychain.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: object mapping a plugin ID to an object with an `options` field, mapping each option name to a string, number, Boolean, or array of strings, and an optional `mcpServers` field holding per-server user configuration values in the same shape
> * **Default**: unset
>
> This example stores the `api_endpoint` option for the `deployer` plugin from `acme-tools`:
>
> ```json settings.json
> {
>   "pluginConfigs": {
>     "deployer@acme-tools": {
>       "options": {
>         "api_endpoint": "https://api.example.com"
>       }
>     }
>   }
> }
> ```
>
> Claude Code ignores project and local entries because it substitutes these values into plugin hook, MCP, and LSP configurations, and a cloned repository must not be able to supply them. Before v2.1.207, project and local settings were also read.
>
> ## MCP
>
> Control which MCP servers Claude Code connects to and which an organization allows. See [Connect to external tools with MCP](/docs/en/mcp) and [Managed MCP configuration](/docs/en/managed-mcp).
>
> ### `allowAllClaudeAiMcps`
>
> Load the [claude.ai connectors](/docs/en/mcp#use-mcp-servers-from-claude-ai) Claude Code fetches itself alongside a deployed `managed-mcp.json`. Without this key, `managed-mcp.json` takes exclusive control of MCP servers and suppresses those connectors.
>
> * **Scope**: [`Managed`](#scopes). Users can't re-enable connectors that exclusive control suppressed.
> * **Type**: Boolean
>   * `true`: Claude Code loads the claude.ai connectors alongside a deployed `managed-mcp.json`
>   * `false`: a deployed `managed-mcp.json` takes exclusive control of MCP servers and suppresses claude.ai connectors
> * **Default**: `false`, so a deployed `managed-mcp.json` suppresses claude.ai connectors
>
> ```json managed-settings.json
> {
>   "allowAllClaudeAiMcps": true
> }
> ```
>
> [`allowedMcpServers`](#allowedmcpservers) and [`deniedMcpServers`](#deniedmcpservers) still apply to the connectors this key loads. Connectors delivered to [cloud sessions](/docs/en/claude-code-on-the-web) stay suppressed. See [Allow claude.ai connectors alongside the managed set](/docs/en/managed-mcp#allow-claude-ai-connectors-alongside-the-managed-set).
>
> ### `allowedMcpServers`
>
> Allowlist the MCP servers people can use. Claude Code blocks any server that doesn't match an entry, whichever settings file or `.mcp.json` defined it, including servers from `managed-mcp.json`. Built-in servers such as Claude in Chrome, IDE-provided servers, and servers the CLI itself configures are exempt from the allowlist; the denylist still applies to them.
>
> * **Scope**: [`Any file`](#scopes). Entries from every file merge into one allowlist unless [`allowManagedMcpServersOnly`](#allowmanagedmcpserversonly) is set. Deploy it in managed settings to enforce it.
> * **Type**: array of objects, each with exactly one key: `serverName`, a string limited to letters, numbers, hyphens, and underscores; `serverCommand`, an array of the command and its arguments matched exactly; or `serverUrl`, a URL pattern with `*` wildcards
> * **Default**: unset, so every server is allowed; an empty array blocks every server
>
> This example allows only the stdio server that the listed `npx` command starts:
>
> ```json settings.json
> {
>   "allowedMcpServers": [
>     { "serverCommand": ["npx", "-y", "@modelcontextprotocol/server-filesystem"] }
>   ]
> }
> ```
>
> A [`deniedMcpServers`](#deniedmcpservers) entry takes precedence, so a server on both lists is blocked. Once the list contains any `serverCommand` entry, a stdio server must match a `serverCommand` entry, and once it contains any `serverUrl` entry, a remote server must match a `serverUrl` entry: a `serverName` match no longer admits that kind of server. See [Policy-based control with allowlists and denylists](/docs/en/managed-mcp#policy-based-control-with-allowlists-and-denylists).
>
> ### `allowManagedMcpServersOnly`
>
> Make the managed allowlist the only one that applies. Claude Code then reads [`allowedMcpServers`](#allowedmcpservers) from managed settings alone and ignores allowlists in user, project, and local settings; [`deniedMcpServers`](#deniedmcpservers) still merges from every file, so users can still block servers for themselves. Administrators set it so a user's own settings can't broaden what the managed allowlist permits.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code reads `allowedMcpServers` from managed settings alone and ignores allowlists in user, project, and local settings
>   * `false`: allowlists from every settings file merge
> * **Default**: `false`, so allowlists from every settings file merge
>
> This example locks the allowlist to managed settings and allows only the server named `github`:
>
> ```json managed-settings.json
> {
>   "allowManagedMcpServersOnly": true,
>   "allowedMcpServers": [
>     { "serverName": "github" }
>   ]
> }
> ```
>
> Users can still add MCP servers of their own; only servers that match the managed allowlist load. See [Restrict the allowlist to managed settings only](/docs/en/managed-mcp#restrict-the-allowlist-to-managed-settings-only).
>
> ### `deniedMcpServers`
>
> Block specific MCP servers. Claude Code refuses to load a matching server in every scope, including servers from `managed-mcp.json` and [claude.ai connectors](/docs/en/mcp#use-mcp-servers-from-claude-ai).
>
> * **Scope**: [`Any file`](#scopes). Entries from every file merge into one denylist, and [`allowManagedMcpServersOnly`](#allowmanagedmcpserversonly) doesn't change that. Deploy it in managed settings to enforce it.
> * **Type**: array of objects, each with exactly one key: `serverName`, any non-empty string, so a claude.ai connector's display name such as `"claude.ai Slack"` works; `serverCommand`, an array of the command and its arguments matched exactly; or `serverUrl`, a URL pattern with `*` wildcards
> * **Default**: unset, so no server is blocked; an empty array also blocks nothing
>
> ```json settings.json
> {
>   "deniedMcpServers": [
>     { "serverName": "filesystem" }
>   ]
> }
> ```
>
> The denylist takes precedence over [`allowedMcpServers`](#allowedmcpservers), so a server on both lists is blocked. See [Policy-based control with allowlists and denylists](/docs/en/managed-mcp#policy-based-control-with-allowlists-and-denylists).
>
> ### `disableClaudeAiConnectors`
>
> Turn off [claude.ai MCP connectors](/docs/en/mcp#use-mcp-servers-from-claude-ai) so Claude Code neither fetches nor connects them. A `true` in any settings file applies: a checked-in project `.claude/settings.json` can opt a repository out of cloud connectors, but a project-level `false` can't override a user- or managed-level `true`. Requires Claude Code v2.1.182 or later.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code neither fetches nor connects claude.ai MCP connectors
>   * `false`: the same as unset; Claude Code fetches your connectors unless another settings file or `ENABLE_CLAUDEAI_MCP_SERVERS` turns them off
> * **Default**: `false`, so Claude Code fetches your connectors
> * **Per-session overrides**: [`ENABLE_CLAUDEAI_MCP_SERVERS`](/docs/en/env-vars) set to `false` turns connectors off for one session; whichever of the two turns them off, the other can't turn them back on
>
> ```json settings.json
> {
>   "disableClaudeAiConnectors": true
> }
> ```
>
> Servers you pass explicitly with `--mcp-config` are unaffected. To block individual connectors instead of all of them, use [`deniedMcpServers`](#deniedmcpservers). See [Disable claude.ai connectors](/docs/en/mcp#disable-claude-ai-connectors). Requires Claude Code v2.1.182 or later.
>
> ### `disabledMcpjsonServers`
>
> Reject specific servers defined in a project's `.mcp.json` file so Claude Code never connects them or asks you to approve them. A rejection in any settings file applies, including a project `.claude/settings.json` checked into the repository.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of strings, the server names as they appear in `.mcp.json`
> * **Default**: unset
>
> ```json settings.json
> {
>   "disabledMcpjsonServers": ["filesystem"]
> }
> ```
>
> Claude Code writes this key to `.claude/settings.local.json` when you reject a server in the approval dialog. `claude mcp get <name>` shows a rejected server as `✘ Rejected (see disabledMcpjsonServers in settings)`. Rejection takes precedence over [`enabledMcpjsonServers`](#enabledmcpjsonservers) and [`enableAllProjectMcpServers`](#enableallprojectmcpservers).
>
> ### `enableAllProjectMcpServers`
>
> Approve every MCP server defined in project `.mcp.json` files without a prompt. Claude Code writes this key to `.claude/settings.local.json` when you choose to approve all servers in the approval dialog.
>
> * **Scope**: [`Any file`](#scopes). In a folder whose trust dialog you haven't accepted, Claude Code honors it from user settings, managed settings, and `--settings` and ignores it in the shared project file, both in the session and for `claude mcp list` and `claude mcp get`; [Project server approvals and workspace trust](/docs/en/mcp#project-server-approvals-and-workspace-trust) says when an untracked `.claude/settings.local.json` counts too.
> * **Type**: Boolean
>   * `true`: Claude Code approves every MCP server defined in project `.mcp.json` files without a prompt
>   * `false`: Claude Code asks you to approve each server. In a trusted folder, a `false` in a higher-precedence file overrides a `true` in a lower one; in a folder you haven't trusted, a `true` in any honored file is enough
> * **Default**: unset, so Claude Code asks you to approve each server
>
> ```json settings.json
> {
>   "enableAllProjectMcpServers": true
> }
> ```
>
> A [`disabledMcpjsonServers`](#disabledmcpjsonservers) entry still rejects a server.
>
> ### `enabledMcpjsonServers`
>
> Approve specific servers defined in project `.mcp.json` files so Claude Code connects them without asking. Claude Code writes this key to `.claude/settings.local.json` when you approve a server in the approval dialog.
>
> * **Scope**: [`Any file`](#scopes). In a folder whose trust dialog you haven't accepted, Claude Code honors it from user settings, managed settings, and `--settings` and ignores it in the shared project file, both in the session and for `claude mcp list` and `claude mcp get`; [Project server approvals and workspace trust](/docs/en/mcp#project-server-approvals-and-workspace-trust) says when an untracked `.claude/settings.local.json` counts too.
> * **Type**: array of strings, the server names as they appear in `.mcp.json`
> * **Default**: unset
>
> This example approves the `memory` and `github` servers from the project's `.mcp.json`:
>
> ```json settings.json
> {
>   "enabledMcpjsonServers": ["memory", "github"]
> }
> ```
>
> A [`disabledMcpjsonServers`](#disabledmcpjsonservers) entry still rejects a server.
>
> ## Agents, sessions, and worktrees
>
> Set the default agent, control teammates and cross-session messaging, and configure worktrees. See [Subagents](/docs/en/sub-agents) and [Worktrees](/docs/en/worktrees).
>
> ### `agent`
>
> Run the main thread as a named [subagent](/docs/en/sub-agents#invoke-subagents-explicitly), so Claude Code applies that subagent's system prompt, tool restrictions, and model to your session. The same key sets the default agent for sessions you dispatch from `claude agents`.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, the name of a built-in or custom agent
> * **Default**: unset, so the main thread runs as Claude Code's default agent
> * **Per-session overrides**: `--agent` takes precedence over this key for one session
>
> ```json settings.json
> {
>   "agent": "code-reviewer"
> }
> ```
>
> A plugin's own `settings.json` can also supply this key; see [Ship default settings with your plugin](/docs/en/plugins#ship-default-settings-with-your-plugin).
>
> ### `crossSessionInbound`
>
> Choose what this session does with [messages arriving from your other Claude Code sessions](/docs/en/cross-session-messaging#control-inbound-messages). When no value applies, Claude Code decides per message from the two sessions' permission-mode classes. Requires Claude Code v2.1.224 or later.
>
> * **Scope**: [`Any file`](#scopes). A project or local value applies only when it's stricter than the value managed settings, the `--settings` flag, or user settings give.
> * **Type**: string, one of:
>   * `"accept"`: Claude Code delivers the message to Claude
>   * `"hold"`: Claude Code shows a notice for the message without delivering it
>   * `"refuse"`: Claude Code drops the message
> * **Default**: unset, so Claude Code decides per message
>
> ```json settings.json
> {
>   "crossSessionInbound": "hold"
> }
> ```
>
> Claude Code reads managed settings first, then the `--settings` flag, then user settings, and applies the first value found. `refuse` is stricter than `hold`, and `hold` is stricter than `accept`. When none of the trusted sources sets a value, a project or local `hold` or `refuse` still applies, replacing the per-message default. In sessions with cross-session messaging, this key appears in `/config` as **Messages from your other sessions**, which writes it to user settings; the row requires Claude Code v2.1.232 or later, and Claude Code hides it while the `--settings` flag or managed settings set the key.
>
> ### `disableAgentView`
>
> Turn off [background agents and agent view](/docs/en/agent-view): `claude agents`, `--bg`, `/background`, and the on-demand supervisor. Set it in [managed settings](/docs/en/managed-settings) to enforce it for an organization.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code turns off `claude agents`, `--bg`, `/background`, and the on-demand supervisor
>   * `false`: agent view is available
> * **Default**: unset, so agent view is available
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_AGENT_VIEW`](/docs/en/env-vars) turns agent view off for one session; whichever of the two turns it off, the other can't turn it back on
>
> ```json settings.json
> {
>   "disableAgentView": true
> }
> ```
>
> ### `isolatePeerMachines`
>
> Require your explicit approval before Claude's `SendMessage` reaches one of your sessions beyond this machine; see [Require approval for cross-machine messages](/docs/en/cross-session-messaging#require-approval-for-cross-machine-messages). The approval prompt appears even in [`bypassPermissions` mode](/docs/en/permission-modes#skip-all-checks-with-bypasspermissions-mode).
>
> * **Scope**: [`Any file`](#scopes). A `true` from any scope applies, so a checked-in project file can turn the requirement on but not off.
> * **Type**: Boolean
>   * `true`: Claude Code asks for your approval before Claude's `SendMessage` reaches one of your sessions beyond this machine
>   * `false`: cross-machine messages don't prompt
> * **Default**: unset, so cross-machine messages don't prompt
>
> ```json settings.json
> {
>   "isolatePeerMachines": true
> }
> ```
>
> The cross-machine `SendMessage` approval requires Claude Code v2.1.224 or later.
>
> ### `processWrapper`
>
> On macOS and Linux, place a corporate launcher command in front of the [background processes Claude Code starts](/docs/en/corporate-launcher#what-the-launcher-covers). Claude Code runs the launcher with its own command line appended, so the launcher must exec into Claude Code; see [Run Claude Code behind a corporate launcher](/docs/en/corporate-launcher) for the launcher contract. Requires Claude Code v2.1.210 or later.
>
> * **Scope**: [`User or managed`](#scopes)
> * **Type**: string, the launcher command as an argv prefix, such as an absolute path with optional arguments
> * **Default**: unset, so background processes start unwrapped
> * **Per-session overrides**: [`CLAUDE_CODE_PROCESS_WRAPPER`](/docs/en/env-vars) takes precedence over this key for one session
>
> ```json settings.json
> {
>   "processWrapper": "/opt/corp/launcher --profile claude"
> }
> ```
>
> Claude Code ignores the launcher on Windows and starts every process unwrapped. Requires Claude Code v2.1.210 or later.
>
> ### `teammateMode`
>
> Choose where Claude Code shows [agent team](/docs/en/agent-teams) teammates: inside your main terminal pane, or in split panes when your terminal supports them. See [Choose a display mode](/docs/en/agent-teams#choose-a-display-mode).
>
> * **Scope**: [`Any file`](#scopes). Claude Code also reads a value left in `~/.claude.json` by older versions.
> * **Type**: string, one of:
>   * `"in-process"`: teammates run inside your main terminal pane
>   * `"auto"`: split panes when you're running inside tmux, or inside iTerm2 with `it2` on your `PATH` or tmux installed; in-process otherwise
>   * `"tmux"`: split panes using tmux or iTerm2, detected from your terminal
>   * `"iterm2"`: iTerm2 native split panes through the `it2` CLI, in Claude Code v2.1.186 or later
> * **Default**: `"in-process"`
> * **Per-session overrides**: `--teammate-mode` takes precedence over this key for one session
>
> ```json settings.json
> {
>   "teammateMode": "auto"
> }
> ```
>
> Before v2.1.179, the default was `auto`. The `iterm2` value requires Claude Code v2.1.186 or later.
>
>
> ### `worktree`
>
> Configure how Claude Code creates and manages [git worktrees](/docs/en/worktrees) for `--worktree`, the `EnterWorktree` tool, and isolated subagents and background sessions.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: object with `baseRef`, `symlinkDirectories`, `sparsePaths`, and `bgIsolation`
> * **Default**: unset
>
> This example branches new worktrees from your current `HEAD` and symlinks `node_modules` into each one:
>
> ```json settings.json
> {
>   "worktree": {
>     "baseRef": "head",
>     "symlinkDirectories": ["node_modules"]
>   }
> }
> ```
>
> To copy gitignored files like `.env` into new worktrees, add a [`.worktreeinclude` file](/docs/en/worktrees#copy-gitignored-files-into-worktrees) to your project root instead of a setting.
>
> ### `worktree.baseRef`
>
> Choose which ref new worktrees branch from. `"fresh"` branches from `origin/<default-branch>` for a clean tree matching the remote; `"head"` branches from your current local `HEAD`, so unpushed commits and feature-branch state are present in the worktree.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"fresh"`: new worktrees branch from `origin/<default-branch>`
>   * `"head"`: new worktrees branch from your current local `HEAD`, including unpushed commits
> * **Default**: `"fresh"`
>
> ```json settings.json
> {
>   "worktree": {
>     "baseRef": "head"
>   }
> }
> ```
>
> Inside a linked worktree, `"head"` resolves to that worktree's `HEAD`, not the main checkout's.
>
> ### `worktree.symlinkDirectories`
>
> Symlink directories from the main repository into each worktree so you don't duplicate large directories on disk.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of strings, directory paths relative to the repository root
> * **Default**: unset, so Claude Code symlinks no directories
>
> This example symlinks `node_modules` and `.cache` from the main repository into every new worktree:
>
> ```json settings.json
> {
>   "worktree": {
>     "symlinkDirectories": ["node_modules", ".cache"]
>   }
> }
> ```
>
> ### `worktree.sparsePaths`
>
> Check out only the listed directories in each worktree through git sparse-checkout. Claude Code writes only those directories plus root-level files to disk, which is faster in large monorepos; see [Check out only the directories you need](/docs/en/large-codebases#check-out-only-the-directories-you-need).
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: array of strings, directory paths relative to the repository root
> * **Default**: unset, so each worktree checks out the whole tree
>
> This example checks out only `packages/my-app` and `shared/utils`, plus root-level files, in each worktree:
>
> ```json settings.json
> {
>   "worktree": {
>     "sparsePaths": ["packages/my-app", "shared/utils"]
>   }
> }
> ```
>
> While a sparse worktree exists, git enables `extensions.worktreeConfig` in the repository's shared `.git/config`.
>
> ### `worktree.bgIsolation`
>
> Choose how [background sessions](/docs/en/agent-view#how-file-edits-are-isolated) isolate their file edits. With `"worktree"`, Claude Code blocks `Edit` and `Write` in the main checkout until the session calls `EnterWorktree`; with `"none"`, background jobs edit the working copy directly. Set `"none"` for a repository where git worktrees are impractical.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, one of:
>   * `"worktree"`: Claude Code blocks `Edit` and `Write` in the main checkout until the session calls `EnterWorktree`
>   * `"none"`: background jobs edit the working copy directly
> * **Default**: `"worktree"`
>
> ```json settings.json
> {
>   "worktree": {
>     "bgIsolation": "none"
>   }
> }
> ```
>
> Outside a git repository, a [`WorktreeCreate` hook](/docs/en/worktrees#non-git-version-control) that fails releases the block so the session can edit the working directory in place; that release requires Claude Code v2.1.203 or later.
>
> ## Remote, desktop, and notifications
>
> Configure Remote Control, cloud environments, the desktop app, and the notifications Claude Code sends when it needs you. See [Remote Control](/docs/en/remote-control).
>
> ### `agentPushNotifEnabled`
>
> Allow Claude to send a push notification to your phone when it decides one is worth sending, for example when a long task finishes. Claude Code syncs this choice to your account, and pushes arrive while [Remote Control](/docs/en/remote-control) is connected. Appears in `/config` as **Push when Claude decides**.
>
> * **Scope**: [`Any file`](#scopes). Claude Code also reads a value left in `~/.claude.json` by older versions.
> * **Type**: Boolean
>   * `true`: Claude can send a push notification to your phone when it decides one is worth sending
>   * `false`: Claude doesn't send those notifications
> * **Default**: `false`
>
> ```json settings.json
> {
>   "agentPushNotifEnabled": true
> }
> ```
>
> See [Mobile push notifications](/docs/en/remote-control#mobile-push-notifications).
>
> ### `awaySummaryEnabled`
>
> Show a one-line session recap when you return to the terminal after a few minutes away. Set it to `false`, or turn off **Session recap** in `/config`, to stop the recap.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: you see a one-line session recap when you return after a few minutes away
>   * `false`: Claude Code shows no recap
> * **Default**: unset, so the recap is on
> * **Per-session overrides**: [`CLAUDE_CODE_ENABLE_AWAY_SUMMARY`](/docs/en/env-vars) takes precedence over this key for one session, in either direction
>
> ```json settings.json
> {
>   "awaySummaryEnabled": false
> }
> ```
>
> Claude Code never shows the recap in non-interactive mode.
>
> ### `disableArtifact`
>
> Turn off the [Artifact](/docs/en/artifacts) tool, which publishes session output as a private web page on claude.ai, for everyone your settings reach, such as an organization through managed settings. To turn the tool off just for yourself, use [`enableArtifact`](#enableartifact) instead, which the **Artifacts** toggle in `/config` writes to your user settings.
>
> Don't put either key in a project's `.claude/settings.json`: Claude Code ignores `enableArtifact` there, and a `disableArtifact` there is overridden by any higher-precedence file rather than acting as a lock.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code turns off the Artifact tool for everyone your settings reach
>   * `false`: the Artifact tool follows `enableArtifact` and your account's availability
> * **Default**: `false`
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_ARTIFACT`](/docs/en/env-vars) set to `1` turns the tool off for one session; whichever of the two turns it off, the other can't turn it back on
>
> ```json settings.json
> {
>   "disableArtifact": true
> }
> ```
>
> A managed `disableArtifact` takes precedence over a user's [`enableArtifact`](#enableartifact) choice.
>
> ### `disableDeepLinkRegistration`
>
> Stop Claude Code from registering the `claude-cli://` protocol handler with the operating system, which it otherwise does after you send the first prompt of an interactive session. [Deep links](/docs/en/deep-links) let external tools open a Claude Code session with a pre-filled prompt. Set this in environments where protocol handler registration is restricted or managed separately.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: the string `"disable"`
> * **Default**: unset, so Claude Code registers the handler
>
> ```json settings.json
> {
>   "disableDeepLinkRegistration": "disable"
> }
> ```
>
> ### `disableDesktopLocalSessions`
>
> Turn off Code sessions that run on the device in the [desktop app](/docs/en/desktop#local-sessions-on-managed-devices), for deployments where developers should work on remote machines over SSH. In the Code tab, the **Local** environment stays in the environment dropdown but is grayed out and can't be selected, with a tooltip saying your organization turned it off; on Windows the WSL entry is grayed out the same way, though whether WSL sessions run on a managed device at all is [governed separately](/docs/en/admin-setup#wsl-sessions-in-claude-code-desktop). New sessions default to the first [SSH connection](/docs/en/desktop#ssh-sessions) if one is configured, and the app refuses to start or resume a session on the device, including an SSH connection back to the same machine. SSH sessions to other hosts and cloud sessions are unaffected. The desktop app reads this key; the terminal CLI ignores it.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean; only the JSON Boolean `true` takes effect
>   * `true`: the desktop app offers no on-device Code sessions; existing local sessions stay listed but can't continue
>   * `false`: local sessions stay available
> * **Default**: unset, so local sessions are available
>
> ```json managed-settings.json
> {
>   "disableDesktopLocalSessions": true
> }
> ```
>
> The desktop app ignores any other value, and a value that isn't a Boolean, such as the string `"true"` or `1`, also logs a warning. Pair it with [`sshConfigs`](#sshconfigs) so users land on a working connection, and with [`sshHostAllowlist`](#sshhostallowlist) to limit which hosts they can reach. See [Local sessions on managed devices](/docs/en/desktop#local-sessions-on-managed-devices).
>
> Claude Desktop supplies Code sessions with policy derived from your desktop configuration, for example the egress allowlist, filesystem sandbox, and MCP restrictions in third-party deployments. Claude Code ignores those parent settings whenever an [admin source](/docs/en/managed-settings#which-managed-source-claude-code-uses) is present: server-managed settings, an MDM or OS-level policy, or a managed settings file. Deploying this key through one of those on a device that had none before, as in third-party deployments, therefore stops the desktop-derived policies from applying. [Let an embedding host add policy](/docs/en/managed-settings#let-an-embedding-host-add-policy) covers when parent settings can still merge; this holds for any key you deploy that way, not only this one.
>
> ### `disableRemoteControl`
>
> Turn off [Remote Control](/docs/en/remote-control): Claude Code then refuses `claude remote-control`, the `--remote-control` flag, auto-start, and the in-session toggle, and reports that your organization's policy disabled it. Place it in [managed settings](/docs/en/managed-settings) for per-device MDM enforcement.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code refuses `claude remote-control`, the `--remote-control` flag, auto-start, and the in-session toggle
>   * `false`: Remote Control stays available
> * **Default**: `false`
>
> ```json settings.json
> {
>   "disableRemoteControl": true
> }
> ```
>
> ### `enableArtifact`
>
> Turn the [Artifact](/docs/en/artifacts) tool on or off for yourself. When unset, the default follows the feature's [availability](/docs/en/artifacts#availability) for your account. The **Artifacts** row in `/config` writes this key, so you don't usually edit it by hand. To turn artifacts off for everyone from managed settings, use [`disableArtifact`](#disableartifact) instead. Requires Claude Code v2.1.196 or later.
>
> * **Scope**: [`User or managed`](#scopes). Claude Code ignores the key in project and local settings so that a repository you clone can't turn the tool on for you.
> * **Type**: Boolean
>   * `true`: Claude Code turns the Artifact tool on for you where it's available for your account
>   * `false`: Claude Code turns the Artifact tool off for you
> * **Default**: unset, so the tool follows your account's availability
>
> ```json settings.json
> {
>   "enableArtifact": false
> }
> ```
>
> A managed [`disableArtifact`](#disableartifact) and your organization's [admin setting](/docs/en/artifacts#manage-artifacts-for-your-organization) take precedence over this key. Requires Claude Code v2.1.196 or later.
>
> ### `inputNeededNotifEnabled`
>
> Get a push notification on your phone when a permission prompt or question is waiting for your input. Claude Code sends these only while [Remote Control](/docs/en/remote-control) is connected. Appears in `/config` as **Push when actions required**.
>
> * **Scope**: [`Any file`](#scopes). Claude Code also reads a value left in `~/.claude.json` by older versions.
> * **Type**: Boolean
>   * `true`: you get a push notification on your phone when a permission prompt or question is waiting, while Remote Control is connected
>   * `false`: Claude Code sends no such notifications
> * **Default**: `false`
>
> ```json settings.json
> {
>   "inputNeededNotifEnabled": true
> }
> ```
>
> See [Mobile push notifications](/docs/en/remote-control#mobile-push-notifications).
>
> ### `preferredNotifChannel`
>
> Choose how Claude Code notifies you when a task completes or a permission prompt is waiting. Appears in `/config` as **Local notifications**.
>
> * **Scope**: [`Any file`](#scopes). Claude Code also reads a value left in `~/.claude.json` by older versions.
> * **Type**: string, one of:
>   * `"auto"`: Claude Code sends a desktop notification in iTerm2, Ghostty, and Kitty, rings the bell in Terminal.app only when its audible bell is off, and does nothing elsewhere
>   * `"terminal_bell"`: Claude Code rings the bell character in any terminal
>   * `"iterm2"`: Claude Code sends an iTerm2 desktop notification
>   * `"iterm2_with_bell"`: Claude Code sends an iTerm2 desktop notification and rings the bell
>   * `"kitty"`: Claude Code sends a Kitty desktop notification
>   * `"ghostty"`: Claude Code sends a Ghostty desktop notification
>   * `"notifications_disabled"`: Claude Code sends no notification
> * **Default**: `"auto"`
>
> ```json settings.json
> {
>   "preferredNotifChannel": "terminal_bell"
> }
> ```
>
> With `"auto"`, Claude Code sends a desktop notification in iTerm2, Ghostty, and Kitty. In Terminal.app it rings the bell character only when you have turned Terminal's audible bell off, and in other terminals it does nothing. Set `"terminal_bell"` to ring the bell character in any terminal. See [Get a terminal bell or notification](/docs/en/terminal-config#get-a-terminal-bell-or-notification).
>
> ### `remote.defaultEnvironmentId`
>
> Pick the default [cloud environment](/docs/en/cloud-environments) for cloud sessions you create from the CLI, such as with `claude --cloud`. Claude Code writes this key to your user settings when you pick an environment with [`/remote-env`](/docs/en/cloud-environments#select-an-environment-from-the-cli).
>
> * **Scope**: [`Any file`](#scopes). For a self-hosted environment ID, user or managed settings, or the `--settings` flag only.
> * **Type**: string, an environment ID such as `env_...` or `ccpool_...`
> * **Default**: unset, so Claude Code uses the Anthropic-hosted environment when one is in your list, and otherwise the first environment it finds
> * **Per-session overrides**: `--environment` takes precedence over this key for the one cloud session it creates
>
> ```json settings.json
> {
>   "remote": {
>     "defaultEnvironmentId": "env_0123abcd"
>   }
> }
> ```
>
> An Anthropic-hosted environment ID, which starts with `env_`, follows the standard settings precedence, so a value in a repository's project settings overrides your user-level pick. A [self-hosted environment](/docs/en/self-hosted-environments) ID, which starts with `ccpool_`, is honored only from user settings, managed settings, and the `--settings` flag; Claude Code ignores one in a repository's project or local settings, and `/remote-env` shows which value it ignored, so a checked-in file can't steer sessions onto a self-hosted environment you didn't choose.
>
> ### `remoteControlAtStartup`
>
> Connect [Remote Control](/docs/en/remote-control) automatically when each interactive session starts, instead of waiting for `/remote-control`. Set it to `true` to turn auto-connect on, `false` to turn it off. Appears in `/config` as **Enable Remote Control for all sessions**.
>
> * **Scope**: [`Any file`](#scopes). Claude Code also reads a value left in `~/.claude.json` by older versions.
> * **Type**: Boolean
>   * `true`: Claude Code connects Remote Control automatically when each interactive session starts
>   * `false`: Claude Code waits for `/remote-control`
> * **Default**: unset, so auto-connect follows your organization's admin default when one is set, and otherwise Claude Code's current default
> * **Per-session overrides**: `--remote-control` turns Remote Control on for one session even when this key is `false`, and no flag turns it off for one session
>
> ```json settings.json
> {
>   "remoteControlAtStartup": true
> }
> ```
>
> Claude Code ignores a `true` from project or local settings, so a repository can turn auto-connect off for its checkout but can't turn it on. For the full per-scope behavior, see [Enable Remote Control for all sessions](/docs/en/remote-control#enable-remote-control-for-all-sessions) and the [security keys where the stricter value applies](/docs/en/settings#security-keys-where-the-stricter-value-applies).
>
> ### `sshConfigs`
>
> Add SSH connections to the [Desktop](/docs/en/desktop#pre-configure-ssh-connections-for-your-team) environment dropdown. Administrators use it to distribute shared connections to a team. Connections you define in managed settings show as managed, so users can select them but can't edit or delete them in the app.
>
> * **Scope**: [`User or managed`](#scopes). The desktop app reads this key.
> * **Type**: array of objects, each with required `id`, `name`, and `sshHost` and optional `sshPort` and `sshIdentityFile`
> * **Default**: unset
>
> This example adds one connection named `Dev VM` that connects to `user@dev.example.com`:
>
> ```json settings.json
> {
>   "sshConfigs": [
>     {
>       "id": "dev-vm",
>       "name": "Dev VM",
>       "sshHost": "user@dev.example.com"
>     }
>   ]
> }
> ```
>
> ### `sshHostAllowlist`
>
> Limit the hosts a [Desktop SSH session](/docs/en/desktop#restrict-which-ssh-hosts-users-can-connect-to) can connect to. Only the Desktop app reads this key; the CLI doesn't. Patterns are case-insensitive: `*` matches any host, `*.example.com` matches `example.com` and every subdomain, and anything else is an exact match against the hostname after `~/.ssh/config` resolution. An empty array turns SSH sessions off.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: array of hostname patterns
> * **Default**: unset, so any host is allowed
>
> This example allows `devboxes.example.com` and its subdomains, plus the exact host `bastion.example.com`:
>
> ```json managed-settings.json
> {
>   "sshHostAllowlist": ["*.devboxes.example.com", "bastion.example.com"]
> }
> ```
>
>
> ## Authentication and providers
>
> Supply credentials through helper scripts and, for organizations, force a login method or organization. See [Authentication](/docs/en/authentication).
>
> ### `apiKeyHelper`
>
> Run your own command to produce the credential Claude Code sends with model requests. Claude Code runs the command through the system shell, `/bin/sh` on macOS and Linux and `cmd` on Windows, and sends its output as both the `X-Api-Key` and `Authorization: Bearer` headers. Use it for dynamic or rotating credentials, such as short-lived tokens fetched from a vault.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, a shell command line
> * **Default**: unset, so Claude Code doesn't run a helper
>
> ```json settings.json
> {
>   "apiKeyHelper": "/bin/generate_temp_api_key.sh"
> }
> ```
>
> Claude Code caches the value and reruns the command after the interval you set with [`CLAUDE_CODE_API_KEY_HELPER_TTL_MS`](/docs/en/env-vars). In interactive sessions, when the command comes from project or local settings, Claude Code doesn't run it until you accept the workspace trust prompt. See [Credential management](/docs/en/authentication#credential-management).
>
> ### `awsAuthRefresh`
>
> Run your own command, such as `aws sso login`, to refresh the credentials in your `.aws` directory when the ones Claude Code has for [Amazon Bedrock](/docs/en/amazon-bedrock) stop working. Claude Code checks the current credentials against STS first and runs the command only when that check fails, then reads the refreshed `.aws` directory.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, a shell command line
> * **Default**: unset, so Claude Code doesn't refresh AWS credentials for you
>
> ```json settings.json
> {
>   "awsAuthRefresh": "aws sso login --profile myprofile"
> }
> ```
>
> Use this key when your refresh flow writes to `.aws`; use [`awsCredentialExport`](#awscredentialexport) when it prints credentials instead. See [advanced credential configuration](/docs/en/amazon-bedrock#advanced-credential-configuration).
>
> ### `awsCredentialExport`
>
> Run your own command that prints AWS credentials as JSON, so Claude Code can call [Amazon Bedrock](/docs/en/amazon-bedrock) with credentials that don't live in your `.aws` directory. Claude Code accepts the `aws sts` output shape and the flat `aws configure export-credentials` shape, and scopes the credentials to its own Bedrock client, so the shell commands Claude runs still see your ambient credentials.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, a shell command line
> * **Default**: unset, so Claude Code uses the ambient AWS credential chain
>
> ```json settings.json
> {
>   "awsCredentialExport": "/bin/generate_aws_grant.sh"
> }
> ```
>
> Unlike [`awsAuthRefresh`](#awsauthrefresh), Claude Code always runs this command when it's set, without checking the ambient credentials first. See [advanced credential configuration](/docs/en/amazon-bedrock#advanced-credential-configuration).
>
> ### `forceLoginMethod`
>
> Restrict which kind of account people can log in with. Set `"claudeai"` to allow only claude.ai accounts, `"console"` to allow only Claude Console accounts, or `"gateway"` to send people to a [cloud gateway](/docs/en/claude-apps-gateway) instead of a first-party login. Administrators set it in managed settings and pair it with [`forceLoginOrgUUID`](#forceloginorguuid) to keep developers' claude.ai logins inside one organization.
>
> * **Scope**: [`Any file`](#scopes). Claude Code honors `"gateway"` only from a managed source on the machine: `managed-settings.json`, the macOS plist or Windows HKLM registry, or a policy helper. It treats `"gateway"` as unset in user, project, local, HKCU, and server-managed settings, the same rule as [`forceLoginGatewayUrl`](#forcelogingatewayurl).
> * **Type**: string, one of:
>   * `"claudeai"`: only claude.ai accounts can log in
>   * `"console"`: only Claude Console accounts can log in
>   * `"gateway"`: Claude Code sends people to a cloud gateway instead of a first-party login
> * **Default**: unset, so people pick a login method
>
> ```json settings.json
> {
>   "forceLoginMethod": "claudeai"
> }
> ```
>
> Every first-party login path applies the restriction, including the [VS Code extension](/docs/en/vs-code), the Agent SDK, `claude setup-token`, and `/install-github-app`, except the terminal's interactive login screen, reached by `/login` or first-run onboarding, which pre-selects the method without enforcing it. Before v2.1.212, only terminal logins applied it. See [Restrict login to your organization](/docs/en/authentication#restrict-login-to-your-organization) for how each login path, environment credentials, and third-party providers are handled.
>
> ### `forceLoginGatewayUrl`
>
> Set the gateway URL the `/login` Cloud gateway screen connects to, so people reach your [cloud gateway](/docs/en/claude-apps-gateway) without typing its address. The screen has no URL field: with this key set, it shows your gateway URL and connects when the person presses Enter; without it, it tells them to contact their IT administrator. When `forceLoginMethod` is unset, this key alone opens the Cloud gateway screen. `forceLoginMethod: "gateway"` also opens it and removes the login-method picker, and a `claudeai` or `console` value there takes precedence over this key. Set both keys so the screen connects instead of showing an error.
>
> * **Scope**: [`Managed`](#scopes). Read only from a source on the machine: `managed-settings.json`, the macOS plist or Windows HKLM registry, or a policy helper. Claude Code ignores it in HKCU and server-managed settings.
> * **Type**: string, a full URL including the scheme
> * **Default**: unset, so the Cloud gateway screen shows an error telling people to contact their IT administrator
>
> ```json managed-settings.json
> {
>   "forceLoginGatewayUrl": "https://claude-gateway.example.com"
> }
> ```
>
> A value that isn't a valid URL is dropped on its own; the rest of the managed settings file still applies. See [Set the gateway URL](/docs/en/claude-apps-gateway#set-the-gateway-url).
>
> ### `forceLoginOrgUUID`
>
> From a managed source, require claude.ai account logins to belong to one Anthropic organization, a single UUID, or to any of several, an array. From any settings file, a single UUID also pre-selects that organization during a claude.ai or Claude Console login; an array pre-selects nothing.
>
> * **Scope**: [`Any file`](#scopes). Only a managed source enforces the restriction; a single UUID in any other settings file pre-selects the organization during login without restricting it.
> * **Type**: string, one UUID, or array of strings, several UUIDs
> * **Default**: unset, so any organization can log in
>
> This example accepts logins from either of two organizations without pre-selecting one:
>
> ```json managed-settings.json
> {
>   "forceLoginOrgUUID": ["xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"]
> }
> ```
>
> An empty array in a managed source blocks every login with a misconfiguration message, and so does a value Claude Code can't parse. See [Restrict login to your organization](/docs/en/authentication#restrict-login-to-your-organization) for how Claude Code treats Claude Console logins, the other login paths, and environment credentials.
>
> ### `gcpAuthRefresh`
>
> Run your own command to refresh Google Cloud Application Default Credentials when Claude Code finds they've expired or can't be loaded, so [Google Cloud's Agent Platform](/docs/en/google-vertex-ai) requests keep working without you re-authenticating by hand.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, a shell command line
> * **Default**: unset, so Claude Code's credential error tells you to run `gcloud auth application-default login` yourself
>
> ```json settings.json
> {
>   "gcpAuthRefresh": "gcloud auth application-default login"
> }
> ```
>
> See [advanced credential configuration](/docs/en/google-vertex-ai#advanced-credential-configuration).
>
> ### `otelHeadersHelper`
>
> Run your own command to generate the headers Claude Code sends with OpenTelemetry exports, for backends whose tokens rotate. Claude Code runs it at startup and periodically after that, and expects a JSON object of string header values on stdout.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: string, an executable path or a shell command line
> * **Default**: unset, so Claude Code adds no helper-generated headers
>
> ```json settings.json
> {
>   "otelHeadersHelper": "/bin/generate_otel_headers.sh"
> }
> ```
>
> Set the refresh interval with [`CLAUDE_CODE_OTEL_HEADERS_HELPER_DEBOUNCE_MS`](/docs/en/env-vars). See [Dynamic headers](/docs/en/monitoring-usage#dynamic-headers) for the script requirements and where Claude Code reports a failing helper.
>
> ## Updates and versioning
>
> Choose an update channel and, for organizations, pin the versions people can run. See [Update Claude Code](/docs/en/setup#update-claude-code).
>
> ### `autoUpdatesChannel`
>
> Choose which [release channel](/docs/en/setup#configure-release-channel) background auto-updates and `claude update` follow. Set `"stable"` for a version that is typically about one week old and skips releases with major regressions, or `"latest"` for the most recent release.
>
> * **Scope**: [`Any file`](#scopes). Set it in managed settings to enforce one channel across your organization.
> * **Type**: string, one of:
>   * `"latest"`: updates follow the most recent release
>   * `"stable"`: updates follow a version that is typically about one week old and skips releases with major regressions
> * **Default**: unset, so Claude Code follows `"latest"`
>
> ```json settings.json
> {
>   "autoUpdatesChannel": "stable"
> }
> ```
>
> Claude Code writes `"stable"` to your user settings when you pick it under **Auto-update channel** in `/config`, and removes the key when you switch back to latest there. `claude install stable` and `claude install latest` also save the channel you name. Switching from `"latest"` to `"stable"` in `/config` asks whether to allow a downgrade or stay on your current version; staying sets [`minimumVersion`](#minimumversion). Homebrew installs ignore this key: the `claude-code` cask tracks stable and `claude-code@latest` tracks latest, and `claude update` defers to `brew upgrade`. To turn auto-updates off entirely, set [`DISABLE_AUTOUPDATER`](/docs/en/setup#disable-auto-updates) in `env`.
>
> ### `minimumVersion`
>
> Keep background auto-updates and `claude update` from installing any version below this one, so moving to the `"stable"` channel doesn't downgrade you from a newer `"latest"` build. Claude Code writes this key for you when you choose to stay on your current version while switching channels in `/config`, and clears it when you switch back to `"latest"`.
>
> * **Scope**: [`Any file`](#scopes). Set it in managed settings to pin an organization-wide minimum that user and project settings can't lower.
> * **Type**: string, a version number such as `"2.1.100"`
> * **Default**: unset, so updates can install any version the channel offers
>
> This example follows the stable channel and refuses to install any version below 2.1.100:
>
> ```json settings.json
> {
>   "autoUpdatesChannel": "stable",
>   "minimumVersion": "2.1.100"
> }
> ```
>
> This key only constrains updates. To make Claude Code refuse to start below a version, use [`requiredMinimumVersion`](#requiredminimumversion) instead. See [Pin a minimum version](/docs/en/setup#pin-a-minimum-version).
>
> ### `requiredMaximumVersion`
>
> Set the newest Claude Code version your organization allows to start. When the running version is newer, Claude Code exits at startup and tells the user to install an approved version through your organization's approved method; `claude install <version>` may also work. Requires Claude Code v2.1.163 or later.
>
> * **Scope**: [`Managed`](#scopes). Claude Code gives no warning when it ignores the key elsewhere.
> * **Type**: string, a version number such as `"2.1.150"`; a value that isn't a valid version is ignored
> * **Default**: unset, so no ceiling applies
>
> ```json managed-settings.json
> {
>   "requiredMaximumVersion": "2.1.150"
> }
> ```
>
> Background auto-updates and `claude update` skip versions above the ceiling, so an installation inside the range stays inside it. `claude update`, `claude install`, and `claude doctor` keep working above the ceiling so users can recover. Pair it with [`requiredMinimumVersion`](#requiredminimumversion) to enforce a range. Requires Claude Code v2.1.163 or later.
>
> ### `requiredMinimumVersion`
>
> Set the oldest Claude Code version your organization allows to start. When the running version is older, Claude Code exits at startup and tells the user to update through your organization's approved method. The check runs at startup only, so a session that's already running continues. Requires Claude Code v2.1.163 or later.
>
> * **Scope**: [`Managed`](#scopes). Claude Code gives no warning when it ignores the key elsewhere.
> * **Type**: string, a version number such as `"2.1.150"`; a value that isn't a valid version is ignored
> * **Default**: unset, so no floor applies
>
> ```json managed-settings.json
> {
>   "requiredMinimumVersion": "2.1.150"
> }
> ```
>
> `claude update`, `claude install`, and `claude doctor` keep working below the floor so users can recover. Unlike [`minimumVersion`](#minimumversion), which only prevents downgrades, this key blocks startup. Pair it with [`requiredMaximumVersion`](#requiredmaximumversion) to enforce a range. Requires Claude Code v2.1.163 or later.
>
> ## Tools
>
> Turn off specific tools in the [Claude Code desktop app](/docs/en/desktop). The terminal CLI ignores these keys. For the tools themselves, see [Tools available to Claude](/docs/en/tools-reference).
>
> ### `browserExternalPageTools`
>
> Stop Claude from using its tools to read or act on external pages in the desktop app's [Browser pane](/docs/en/desktop#browse-external-sites). People in your organization can still open external sites themselves, and local dev server previews keep working with Claude's tools. The desktop app reads this key; the terminal CLI ignores it.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: string, `"disabled"`; the desktop app also accepts `"disable"`, in either case
> * **Default**: unset, so Claude's tools work on external pages
>
> ```json managed-settings.json
> {
>   "browserExternalPageTools": "disabled"
> }
> ```
>
> Any other value leaves Claude's tools on, and a non-empty string that isn't one of the two accepted values logs a warning. To block external sites for people and Claude alike, set [`disableBrowserExternalNavigation`](#disablebrowserexternalnavigation) instead. See [Restrict external browsing for your organization](/docs/en/desktop#restrict-external-browsing-for-your-organization).
>
> ### `disableBrowserExternalNavigation`
>
> Turn off external browsing in the desktop app's [Browser pane](/docs/en/desktop#browse-external-sites) for people and Claude alike. Localhost dev server previews keep working. The desktop app reads this key; the terminal CLI ignores it.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean; only the JSON Boolean `true` takes effect
>   * `true`: the desktop app turns off external browsing in the Browser pane for people and Claude alike; localhost previews keep working
>   * `false`: external browsing stays on
> * **Default**: unset, so external browsing is on
>
> ```json managed-settings.json
> {
>   "disableBrowserExternalNavigation": true
> }
> ```
>
> The desktop app ignores any other value, and a value that isn't a Boolean, such as the string `"true"` or `1`, also logs a warning. To leave external browsing on but keep Claude's tools off external pages, set [`browserExternalPageTools`](#browserexternalpagetools) instead. See [Restrict external browsing for your organization](/docs/en/desktop#restrict-external-browsing-for-your-organization).
>
> ### `disableMobileSimulatorTools`
>
> Block Claude's tools for the desktop app's [iOS Simulator pane](/docs/en/desktop-ios-simulator#turn-off-simulator-access). People keep manual use of the pane; only Claude's access is removed, and nobody can turn it back on from inside the app. The desktop app reads this key; the terminal CLI ignores it.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean; only the JSON Boolean `true` takes effect
>   * `true`: the desktop app blocks Claude's tools for the iOS Simulator pane
>   * `false`: Claude's simulator tools follow each person's settings toggle in the desktop app
> * **Default**: unset, so Claude's simulator tools follow each person's settings toggle in the desktop app
>
> ```json managed-settings.json
> {
>   "disableMobileSimulatorTools": true
> }
> ```
>
> The desktop app ignores any other value, and a value that isn't a Boolean, such as the string `"true"` or `1`, also logs a warning.
>
>
> ## Privacy and telemetry
>
> Control how long Claude Code keeps session data and what it sends. The switches that turn off usage metrics and error reports are environment variables, not settings keys: set `DISABLE_TELEMETRY`, `DISABLE_ERROR_REPORTING`, or `CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC` in the [`env`](#env) key or in the shell. [Telemetry services](/docs/en/data-usage#telemetry-services) says what each one stops. The session survey is the exception: [`feedbackSurveyRate`](#feedbacksurveyrate) below turns it off from a settings file.
>
> ### `cleanupPeriodDays`
>
> Set how many days Claude Code keeps [session transcripts and other application data](/docs/en/claude-directory#cleaned-up-automatically) before deleting them. Claude Code runs the deletion as a background sweep after a session starts, as long as it can safely determine the retention period.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: number of days, a whole number, minimum `1`
> * **Default**: `30`
>
> ```json settings.json
> {
>   "cleanupPeriodDays": 20
> }
> ```
>
> Setting `0` fails validation, so pick a large value such as `3650` for long retention. To stop Claude Code from writing transcripts at all, see [Plaintext storage](/docs/en/claude-directory#plaintext-storage).
>
> ### `feedbackSurveyRate`
>
> Set the probability that the [session quality survey](/docs/en/data-usage#session-quality-surveys) appears when a session is eligible for it. Set `0` to keep the survey from appearing.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: number between `0` and `1`
> * **Default**: unset, so Claude Code uses the rate Anthropic sets remotely, or its built-in rate of `0.005` on Amazon Bedrock, Google Cloud's Agent Platform, and Microsoft Foundry, which don't receive remote configuration
> * **Per-session overrides**: [`CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY`](/docs/en/env-vars) set to `1` turns the survey off for one session whatever rate this key sets
>
> ```json settings.json
> {
>   "feedbackSurveyRate": 0.05
> }
> ```
>
> The same rate applies to the survey in the VS Code extension.
>
> ### `skipWebFetchPreflight`
>
> Skip the [WebFetch domain safety check](/docs/en/data-usage#webfetch-domain-safety-check), which sends each requested hostname to `api.anthropic.com` before fetching. Set `true` in environments that block traffic to Anthropic, such as Amazon Bedrock, Google Cloud's Agent Platform, or Microsoft Foundry deployments with restrictive egress.
>
> * **Scope**: [`Any file`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code skips the WebFetch domain safety check
>   * `false`: the check runs before the first fetch to each hostname in a session, and again for a hostname whose earlier check was blocked or failed
> * **Default**: unset, so the check runs before the first fetch to each hostname in a session
>
> ```json settings.json
> {
>   "skipWebFetchPreflight": true
> }
> ```
>
> With the check skipped, WebFetch attempts any URL without consulting the blocklist, so pair it with [`WebFetch` permission rules](/docs/en/permissions#webfetch) if you need to restrict which domains Claude can reach.
>
>
> ## Enterprise and managed settings
>
> Keys an organization uses to compute, refresh, and combine managed settings. See [Set up managed settings](/docs/en/admin-setup).
>
> ### `disableSideloadFlags`
>
> Reject the `--plugin-dir`, `--plugin-url`, `--agents`, and `--mcp-config` CLI flags at startup, which users could otherwise pass to bypass [`strictKnownMarketplaces`](#strictknownmarketplaces) for a single run. Claude Code exits with an error naming the rejected flags, and applies the same check to surfaces that start the CLI with these flags internally, currently [Cowork](/docs/en/desktop) local sessions in the desktop app. In [cloud sessions](/docs/en/claude-code-on-the-web), Claude Code drops the MCP servers the server delivered through `--mcp-config`, other than in-process `type: "sdk"` entries, and starts the session. Requires Claude Code v2.1.193 or later.
>
> * **Scope**: [`Managed`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code rejects `--plugin-dir`, `--plugin-url`, `--agents`, and `--mcp-config` at startup and exits with an error naming them, except that in cloud sessions it drops the MCP servers the server delivered through `--mcp-config`, other than in-process `type: "sdk"` entries, and starts the session
>   * `false`: Claude Code accepts those flags
> * **Default**: `false`
>
> ```json managed-settings.json
> {
>   "disableSideloadFlags": true
> }
> ```
>
> Claude Code still accepts a `--mcp-config` whose servers are all in-process `type: "sdk"` entries, so the Agent SDK and VS Code extension keep working. Users can still add servers with `claude mcp add` or a `.mcp.json` file; for per-server control, set [`allowedMcpServers`](/docs/en/managed-mcp) as well. Requires Claude Code v2.1.193 or later.
>
> In cloud sessions, Claude Code also ignores server-delivered mid-session MCP updates, the path behind cloud session configuration and SDK `setMcpServers()` on remote workers. In-process `type: "sdk"` entries stay exempt there too. Before v2.1.239, a server-delivered `--mcp-config` blocked a cloud session from starting.
>
> ### `forceRemoteSettingsRefresh`
>
> Block CLI startup until Claude Code has freshly fetched [server-managed settings](/docs/en/server-managed-settings). If the fetch fails, Claude Code exits instead of continuing with cached or no settings. When the key is unset, startup continues without waiting for remote settings. A Cloud gateway session always waits, and exits if the gateway can't be reached. Set it when your environment can't accept even a brief window in which a session runs without its managed policy.
>
> * **Scope**: [`Managed`](#scopes). Claude Code honors a `true` from any admin-controlled managed source, even one that isn't the highest-priority source.
> * **Type**: Boolean
>   * `true`: Claude Code blocks startup until it has freshly fetched server-managed settings, and exits if the fetch fails
>   * `false`: startup continues without waiting for remote settings
> * **Default**: `false`
>
> ```json managed-settings.json
> {
>   "forceRemoteSettingsRefresh": true
> }
> ```
>
> Set it in an MDM profile or the managed settings file to enforce fail-closed startup before the first server payload arrives. Claude Code applies the check only in sessions that fetch server-managed settings, so a session that [doesn't fetch them](/docs/en/server-managed-settings#platform-availability) starts without waiting. The `claude auth` subcommands are exempt, so users can re-authenticate when expired credentials are why the fetch fails. See [Enforce fail-closed startup](/docs/en/server-managed-settings#enforce-fail-closed-startup).
>
> ### `parentSettingsBehavior`
>
> Choose whether Claude Code applies managed settings supplied by an embedding host process, such as the Agent SDK or an IDE extension, when an admin-deployed managed tier is also present. With `"first-wins"`, Claude Code drops the host-supplied settings; with `"merge"`, it applies them under the admin tier through a restrictive-only filter. Set `"merge"` when a host needs to pass its own restrictions to the sessions it launches, for example Claude Desktop delivering a gateway's egress allowlist.
>
> * **Scope**: [`Managed`](#scopes). Claude Code reads it from the highest-priority admin-controlled managed source.
> * **Type**: string, one of:
>   * `"first-wins"`: Claude Code drops the host-supplied settings when an admin-deployed managed tier is present
>   * `"merge"`: Claude Code applies the host-supplied settings under the admin tier through a restrictive-only filter
> * **Default**: `"first-wins"`
>
> ```json managed-settings.json
> {
>   "parentSettingsBehavior": "merge"
> }
> ```
>
> This key has no effect when no admin-deployed managed tier exists: the host's settings then apply as the only managed tier, still filtered to restrictive values. For the filter's limits and how the managed sources interact, see [Parent settings from embedding hosts](/docs/en/managed-settings#parent-settings-from-embedding-hosts) and [Restrict parent settings](/docs/en/claude-apps-gateway#restrict-parent-settings).
>
>
> ### `policyHelper`
>
> Run an executable you deploy that computes managed settings at startup, so you can derive policy from device posture, identity, or a remote service instead of a static file. Claude Code runs the helper before it accepts the first prompt and treats the settings it emits as the managed settings for the session.
>
> * **Scope**: [`Managed`](#scopes). Read from the macOS plist, the Windows HKLM registry, or the managed settings file. Claude Code reads the key from the highest-priority managed source that delivers settings and runs the helper only when that source is one of those three; it ignores the key in server-managed settings, the HKCU registry, and host-supplied parent settings.
> * **Type**: object with `path`, `timeoutMs`, and `refreshIntervalMs`
> * **Default**: unset, so no helper runs
>
> This example runs the helper with a 5-second timeout and re-runs it every five minutes:
>
> ```json managed-settings.json
> {
>   "policyHelper": {
>     "path": "/usr/local/bin/claude-policy",
>     "timeoutMs": 5000,
>     "refreshIntervalMs": 300000
>   }
> }
> ```
>
> #### Write the helper output
>
> Claude Code runs the helper with no arguments, sets `CLAUDE_CODE_VERSION` in its environment, and reads a JSON envelope from stdout, capped at 1 MB. Put the settings under a `managedSettings` key. A bare settings object with no `managedSettings` key parses with `managedSettings` undefined and applies nothing, and Claude Code reports no error:
>
> ```json
> {
>   "managedSettings": {
>     "permissions": { "deny": ["Read(//etc/secrets/**)"] }
>   }
> }
> ```
>
> When the helper emits `managedSettings`, that object becomes the only managed settings source for the run: Claude Code ignores the MDM, file, and HKCU sources, reads the [cross-source keys](/docs/en/managed-settings#keys-read-from-every-admin-source) from the helper's output alone, and never merges [parent settings](/docs/en/managed-settings#parent-settings-from-embedding-hosts). The startup `forceRemoteSettingsRefresh` check runs before the helper and reads any admin source. A helper that exits 0 without emitting `managedSettings` contributes no managed settings, and the other sources apply as usual. When the helper exits non-zero at startup, Claude Code prints the error and refuses to start, so a helper that needs outage resilience should serve from its own cache and exit `0`.
>
> ### `policyHelper.path`
>
> Name the helper executable Claude Code runs. Claude Code refuses to start when the path isn't absolute, or on Windows when it doesn't end in `.exe`.
>
> * **Scope**: [`Managed`](#scopes). Read from the macOS plist, the Windows HKLM registry, or the managed settings file, wherever [`policyHelper`](#policyhelper) is read.
> * **Type**: string, an absolute path in normalized form, without `.` or `..` segments
> * **Default**: none; required when `policyHelper` is set
>
> ```json managed-settings.json
> {
>   "policyHelper": {
>     "path": "/usr/local/bin/claude-policy"
>   }
> }
> ```
>
> ### `policyHelper.timeoutMs`
>
> Set how long Claude Code waits for the helper before treating the run as failed. A timed-out run fails the same way as a non-zero exit, so at startup Claude Code refuses to start.
>
> * **Scope**: [`Managed`](#scopes). Read from the macOS plist, the Windows HKLM registry, or the managed settings file, wherever [`policyHelper`](#policyhelper) is read.
> * **Type**: integer, milliseconds, minimum `1000`
> * **Default**: `10000`
>
> ```json managed-settings.json
> {
>   "policyHelper": {
>     "path": "/usr/local/bin/claude-policy",
>     "timeoutMs": 5000
>   }
> }
> ```
>
> ### `policyHelper.refreshIntervalMs`
>
> Have Claude Code re-run the helper in the background on an interval so policy changes reach a running session. When a refresh succeeds, its output replaces the previous managed settings without a restart; when a refresh fails, Claude Code keeps the policy it already has.
>
> * **Scope**: [`Managed`](#scopes). Read from the macOS plist, the Windows HKLM registry, or the managed settings file, wherever [`policyHelper`](#policyhelper) is read.
> * **Type**: integer, milliseconds: `0` to disable refresh, otherwise at least `60000`
> * **Default**: unset, so Claude Code runs the helper once at startup
>
> This example re-runs the helper every five minutes:
>
> ```json managed-settings.json
> {
>   "policyHelper": {
>     "path": "/usr/local/bin/claude-policy",
>     "refreshIntervalMs": 300000
>   }
> }
> ```
>
> ### `wslInheritsWindowsSettings`
>
> Have Claude Code on WSL read managed settings from the Windows policy chain, with HKLM and the Windows managed settings file taking priority over `/etc/claude-code` and HKCU below it. While the chain is on, Claude Code reads `/etc/claude-code` only when no managed settings file or drop-in under `C:\Program Files\ClaudeCode\` delivers a policy key other than `wslInheritsWindowsSettings` itself. Set it to extend the policy you already deploy on Windows to WSL sessions on the same machine, so they follow the same rules as host sessions. Claude Code honors it only when set in the HKLM registry key or in a managed settings file or drop-in under `C:\Program Files\ClaudeCode\`, both of which require Windows admin to write.
>
> * **Scope**: [`Managed`](#scopes). In an admin-controlled Windows source.
> * **Type**: Boolean
>   * `true`: Claude Code on WSL reads managed settings from the Windows policy chain, and reads `/etc/claude-code` only when no managed settings file or drop-in under `C:\Program Files\ClaudeCode\` delivers a policy key other than `wslInheritsWindowsSettings`
>   * `false`: WSL reads only `/etc/claude-code`
> * **Default**: `false`, so WSL reads only `/etc/claude-code`
>
> ```json managed-settings.json
> {
>   "wslInheritsWindowsSettings": true
> }
> ```
>
> Once an admin source turns the chain on, HKCU policy joins it on WSL only when HKCU also sets the key to `true`. That copy doesn't turn the chain on by itself. A Windows source that contains only this key doesn't count as a policy source, so a lower-priority source still supplies the policy. This key has no effect on native Windows.
>
> ## Global config settings
>
> Save these keys in `~/.claude.json`, not in a settings file. Claude Code ignores them anywhere else. Claude Code and `/config` write most of them for you, and you can also edit them by hand.
>
> ### `autoConnectIde`
>
> Connect to a running IDE automatically when you start Claude Code from an external terminal. Appears in `/config` as **Auto-connect to IDE (external terminal)** when you run Claude Code outside a VS Code or JetBrains terminal.
>
> * **Scope**: [`Global config`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code connects to a running IDE automatically when you start it from an external terminal
>   * `false`: Claude Code doesn't connect automatically from an external terminal; inside a VS Code or JetBrains terminal, or with `--ide`, it still connects
> * **Default**: `false`
> * **Per-session overrides**: [`CLAUDE_CODE_AUTO_CONNECT_IDE`](/docs/en/env-vars) takes precedence over this key for one session, in either direction
>
> ```json ~/.claude.json
> {
>   "autoConnectIde": true
> }
> ```
>
> Claude Code ignores this key in `settings.json`.
>
> ### `autoInstallIdeExtension`
>
> Install the Claude Code IDE extension automatically when you run Claude Code from a VS Code terminal. Appears in `/config` as **Auto-install IDE extension** when you run Claude Code inside a VS Code or JetBrains terminal.
>
> * **Scope**: [`Global config`](#scopes)
> * **Type**: Boolean
>   * `true`: Claude Code installs the IDE extension automatically when you run it from a VS Code terminal
>   * `false`: Claude Code doesn't install the extension automatically
> * **Default**: `true`
> * **Per-session overrides**: [`CLAUDE_CODE_IDE_SKIP_AUTO_INSTALL`](/docs/en/env-vars) set to `1` skips the install for one session even when this key is `true`
>
> ```json ~/.claude.json
> {
>   "autoInstallIdeExtension": false
> }
> ```
>
> Claude Code ignores this key in `settings.json`.
>
> ### `diffTool`
>
> Choose where Claude Code shows the diff of an `Edit` or `Write` change it proposes when a [VS Code](/docs/en/vs-code) or [JetBrains](/docs/en/jetbrains#features) IDE is connected: `"auto"` opens it in the IDE's diff viewer, `"terminal"` keeps it in the terminal. Appears in `/config` as **Diff tool** only while Claude Code is connected to a VS Code or JetBrains IDE.
>
> * **Scope**: [`Global config`](#scopes)
> * **Type**: string, one of:
>   * `"auto"`: Claude Code opens the diff in the IDE's diff viewer when a VS Code or JetBrains IDE is connected
>   * `"terminal"`: Claude Code keeps the diff in the terminal
> * **Default**: `"auto"`
>
> ```json ~/.claude.json
> {
>   "diffTool": "terminal"
> }
> ```
>
> Claude Code ignores this key in `settings.json`.
>
> ### `externalEditorContext`
>
> When you press `Ctrl+G`, Claude Code opens the prompt you're typing in your [external editor](/docs/en/interactive-mode#general-controls). With this key on, the editor buffer starts with Claude's previous response as `#` comment lines, so you can read it while you write, and Claude Code strips those lines when you save. Appears in `/config` as **Show last response in external editor**.
>
> * **Scope**: [`Global config`](#scopes)
> * **Type**: Boolean
>   * `true`: the editor buffer starts with Claude's previous response as `#` comment lines, which Claude Code strips when you save
>   * `false`: the editor buffer opens with only your prompt
> * **Default**: `false`
>
> ```json ~/.claude.json
> {
>   "externalEditorContext": true
> }
> ```
>
> With it on, the buffer Claude Code opens looks like this, and only the text below the marker line is sent as your prompt:
>
> ```text
> # ─── Claude's last response (for reference; removed on save) ───
> # I added the retry loop to fetchUser in src/api.ts and a test
> # for the timeout case. Want me to wire the same retry into
> # fetchOrders?
> # ─── Write your reply below this line ──────────────────────────
>
> Yes, and cap it at three attempts.
> ```
>
> Claude Code keeps the last 50 lines of the response and marks the cut with `# … (earlier output truncated)`.
>
> Claude Code ignores this key in `settings.json`.
>
> ### `permissionExplainerEnabled`
>
> When Claude asks permission to run a Bash or PowerShell command, you can press `Ctrl+E` on the prompt to get a model-generated [explanation of the command](/docs/en/permissions#permission-system): what it does, why Claude is running it, and what could go wrong, labeled **Low risk**, **Med risk**, or **High risk**. Claude Code asks the model for the explanation only when you press the shortcut, and showing it doesn't run the command. Set this key to `false` to turn the shortcut off.
>
> * **Scope**: [`Global config`](#scopes)
> * **Type**: Boolean
>   * `true`: you can press `Ctrl+E` on a Bash or PowerShell permission prompt to get a model-generated explanation of the command
>   * `false`: Claude Code turns the `Ctrl+E` shortcut off
> * **Default**: `true`
>
> ```json ~/.claude.json
> {
>   "permissionExplainerEnabled": false
> }
> ```
>
> Claude Code ignores this key in `settings.json`.
>
> ### `teammateDefaultModel`
>
> Through v2.1.233, you set this key to the model for [agent team](/docs/en/agent-teams#specify-teammates-and-models) teammates your prompt didn't name a model for: an alias such as `"sonnet"`, or `null` to follow the lead's model. For the model Claude Code picks for such teammates now, see [specify teammates and models](/docs/en/agent-teams#specify-teammates-and-models).
>
> * **Scope**: [`Global config`](#scopes). On v2.1.233 and earlier.
> * **Type**: string, a model alias or full model ID, or `null`
> * **Default**: unset

### Source: Example settings files

> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Example settings files
>
> > Realistic settings.json files for a developer, a team, and an organization: copy one, keep the keys you want, and change the values.
>
> This page holds three example `settings.json` files, one for each place you save a setting:
>
> * A developer's `~/.claude/settings.json`
> * A team's `.claude/settings.json`, committed to the repository
> * An organization's `managed-settings.json`
>
> Each one is a plausible file for that reader, so you can see the shape and copy the parts you want. None of them is a recommended baseline. Every value comes from the key's entry on the [settings reference](/docs/en/settings-reference), which has its type, default, and where it can be set.
>
> Each example has two tabs. **Copyable settings file** is the file as you'd save it. **What each key does** is the same file with a comment above each key; Claude Code doesn't accept comments in a settings file, so copy from the first tab.
>
> ## Your own settings
>
> One developer's personal settings. It picks a model and effort, adjusts the terminal, and pre-approves a read-only command and one file read. Everything not listed keeps its default. A file like this goes in `~/.claude/settings.json`, where it applies to every project you open.
>
>
> </Tabs>
>
> <h2 id="a-teams-shared-settings">
>   A team's shared settings
> </h2>
>
> One team's shared settings, committed to the repository so everyone who clones it gets the same permissions, hooks, telemetry, and plugin marketplace. Save a file like this at `.claude/settings.json` at the top of the repository. Three things to know before you commit one:
>
> * **Cloud sessions read it too.** A [cloud session](/docs/en/settings#settings-in-cloud-sessions) on Claude Code on the web starts from a clone of the repository, so the committed file applies there as well.
> * **Allow rules wait for trust.** Allow rules and `extraKnownMarketplaces` entries take effect after each person [trusts this folder itself](/docs/en/permissions#project-allow-rules-and-workspace-trust), not only a parent folder; deny and ask rules apply in every session, trusted or not.
> * **The hook is a script in the repo.** This file's hook runs `.claude/hooks/block-rm.sh`; [How a hook resolves](/docs/en/hooks#how-a-hook-resolves) walks through writing it.
>
>
> </Tabs>
>
> <h2 id="an-organizations-managed-settings">
>   An organization's managed settings
> </h2>
>
> A `managed-settings.json` file that shows the shape of the managed keys, with one plausible value for each. It isn't a recommended policy: pick the keys that match your own requirements and set your own values. The example sets these keys:
>
> * `forceLoginMethod` and `forceLoginOrgUUID` pin the login method and organization
> * `availableModels` and `enforceAvailableModels` restrict which models sessions can use
> * `permissions.deny` blocks two file reads and `curl`, and `disableBypassPermissionsMode` removes the bypass permission mode
> * `allowManagedPermissionRulesOnly` and `allowManagedMcpServersOnly` make the managed permission and MCP allowlists the only ones that apply
> * `allowedMcpServers` pins the MCP server by URL
> * `strictKnownMarketplaces` allows one plugin marketplace
> * `sandbox` sandboxes commands with a fixed network allowlist and no unsandboxed retry
> * `requiredMinimumVersion` sets a minimum Claude Code version
> * `cleanupPeriodDays` shortens retention of session transcripts and other local data to seven days
> * `companyAnnouncements` shows a message at startup
>
> Administrators deploy a file like this as `managed-settings.json`, or the same JSON through MDM or [server-managed settings](/docs/en/server-managed-settings). One deployed file applies to every machine or account it reaches. To give a group different values, deploy a different file or profile to that group, since [server-managed settings don't support per-group policy yet](/docs/en/server-managed-settings#current-limitations).
>
>
> </Tabs>
