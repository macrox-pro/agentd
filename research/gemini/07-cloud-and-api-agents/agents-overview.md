---
primary_sources:
  - id: T2-AGENTS
    title: "Agents overview"
    url: "https://ai.google.dev/gemini-api/docs/agents"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Agents overview

> **Applicability:** Verbatim excerpts (snapshot 2026-08-29).

### Source: Agents overview — Full page

> Antigravity agent 
>
>  Building managed agents 
>
>  Environments 
>
>  Hooks 
>
>  Deep Research agent 
>
>  
>  Tools 
>  
>
>  Overview 
>
>  Google Search 
>
>  Google Maps 
>
>  Code execution 
>
>  URL context 
>
>  Computer use 
>
>  File search 
>
>  Combine tools and function calling 
>
>  
>  Live API 
>  
>
>  Overview 
>
>  
>  
>  Get started 
>  Get started using the GenAI SDK 
>  Get started using raw WebSockets 
>  
>
>  Capabilities 
>
>  Live transcription 
>
>  Live translation 
>
>  Tool use 
>
>  Session management 
>
>  Ephemeral tokens 
>
>  Best practices 
>
>  
>  Optimization 
>  
>
>  Overview 
>
>  Batch API 
>
>  Webhooks 
>
>  Flex inference 
>
>  Priority inference 
>
>  Context caching 
>
>  
>  Guides 
>  
>
>  Interactions API 
>
>  Streaming 
>
>  Background execution 
>
>  
>  
>  File input 
>  Input methods 
>  Files API 
>  
>
>  OpenAI compatibility 
>
>  Media resolution 
>
>  Token counting 
>
>  Prompt engineering 
>
>  
>  
>  Logs and datasets 
>  Get started with logs 
>  Data logging and sharing 
>  
>
>  
>  
>  Safety 
>  Safety settings 
>  Safety guidance 
>  
>
>  
>  
>  Frameworks 
>  LangChain & LangGraph 
>  CrewAI 
>  LlamaIndex 
>  Vercel AI SDK 
>  Temporal 
>  
>
>  
>  Resources 
>  
>
>  Release notes 
>
>  Deprecations 
>
>  Libraries 
>
>  
>  
>  Migration 
>  Migrate to Gen AI SDK 
>  Migrate to Interactions API 
>  Interactions breaking changes (May 2026) 
>  
>
>  Rate limits 
>
>  Billing info 
>
>  API troubleshooting 
>
>  API errors 
>
>  Status 
>
>  Partner and library integrations 
>
>  
>  
>  Google AI Studio 
>  Quickstart 
>  Google AI plans 
>  Vibe code in Build mode 
>  Developing full-stack apps 
>  Build Android apps 
>  Deploying your app 
>  Agents in AI Studio Playground 
>  Try out LearnLM 
>  Troubleshooting 
>  Access for Workspace users 
>  
>
>  
>  
>  Google Cloud Platform 
>  Gemini Enterprise Agent Platform Gemini API 
>  OAuth authentication 
>  
>
>  
>  Policies 
>  
>
>  Terms of service 
>
>  Available regions 
>
>  Abuse monitoring 
>
>  Feedback information 
>
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>
>  
>  
>  
>  
>  
>  
>  
>  Gemini 3.7 Flash is now available. Try it out .
>  
>  
>  
>  
>  
>  
>
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  Home
>  
>  
>  
>  
>  
>  
>
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  Gemini API
>  
>  
>  
>  
>  
>  
>
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  Docs
>  
>  
>  
>  
>  
>  
>
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>
>  
>  
>  
>  Send feedback
>  
>  
>  
>  
>  
>  Agents overview 
>  
>  
>  
>  
>
>  
>  
>
>  
>  
>  
>  
>  
>
>  
>  
>  
>  
>
>  
>
>  Managed agents on the Gemini API give you a configurable agent
> harness. A single API call provisions a Linux sandbox where the agent reasons,
> executes code, manages files, and browses the web autonomously.
>
>  
>  
>  
>  
>  rocket_launch 
>  
>  
>  
>  Quickstart
>
>  
>  Make your first agent call, stream responses, and build a custom agent.
>  
>
>  
>  
>  
>  
>  
>  smart_toy 
>  
>  
>  
>  Antigravity Agent
>
>  
>  Capabilities, tools, multimodal input, and pricing for the default agent.
>  
>
>  
>  
>  
>  
>  
>  experiment 
>  
>  
>  
>  Agents in AI Studio
>
>  
>  Visual playground for prototyping agents without writing code.
>  
>
>  
>  
>  
>
>  Available Managed agents
>
>  
>  Antigravity agent : General-purpose
> managed agent powered by Gemini 3.7 Flash. Runs code, manages files, and
> searches the web inside a secure Linux sandbox hosted by Google. You can
> configure the underlying model (such as Gemini 3.7 Flash, Gemini 3.6 Flash, or Gemini 3.5 Flash)
> using agent_config , and extend it with your own instructions, skills, and data to
>  build a custom agent .
>
>  Deep Research : Autonomous research agent
> that plans, executes, and synthesizes multi-step research tasks for use cases
> like market analysis, due diligence, and literature reviews.
>
>  
>
>  Security and best practices
>
>  Note: Managed agents are in Public Preview . Review the agent&#39;s actions and
> outputs before relying on them in sensitive workflows. 
>  Every agent runs in a sandboxed environment that is isolated at the OS level.
> The sandbox has unrestricted outbound network access by default. You can
> restrict or disable network access using an allowlist.
>
>  Network access
>
>  By default, environments have unrestricted outbound network access. Use a
>  network allowlist to restrict outbound traffic to specific domains or
> wildcard patterns. For configuration details, see
>  Network Allow List (AI
> Studio) or Network rules 
> (API).
>
>  External tools and APIs
>
>  You can connect external tools and APIs to extend the agent. Only use tools
> from trusted sources and scope permissions to the minimum required. Credentials
> can be injected securely via egress proxy header transformations and are never
> exposed inside the sandbox. The agent may use any credential it has access to,
> so only provide credentials whose full scope you are willing to grant.
>
>  
>  Use least-privilege service accounts or API keys.
>
>  Prefer short-lived tokens over long-lived keys.
>
>  Only provide credentials whose full scope you are willing to grant.
>
>  Rotate credentials on a regular schedule.
>
>  
>
>  For details on configuring header transformations, see
>  Credentials .
>
>  Human oversight
>
>  Always verify outputs (generated code, data transformations, configuration
> changes) before deploying them, especially for tasks that modify data or
> interact with external systems.
>
>  Pricing
>
>  Managed agents use a pay-as-you-go model 
> based on Gemini model tokens and tool usage. A single interaction can trigger
> multiple reasoning loops, typically consuming 100k to 3M tokens. Environment
> compute is not billed during the preview. See estimated costs 
> for per-task breakdowns. Managed agents are also available on the free tier with
> a free rate limit and usage quota.
>
>  Limits
>
>  
>  
>  
>  Limit 
>  Description 
>  
>  
>
>  
>  
>  Environment Lifetime 
>  Environments are permanently deleted after 7 days of inactivity. 
>  
>  
>  VM Spin-down 
>  VMs shut down after a brief period of inactivity to conserve resources. The next request restores the state (with a cold start). 
>  
>  
>  Pre-installed Software 
>  Ubuntu-based environment with Python 3.12 and Node.js 22. For more information on the environment&#39;s base image, see Pre-installed software . 
>  
>  
>  Max agents 
>  You can have up to 1,000 managed agents. 
>  
>  
>  
>
>  Agent frameworks
>
>  You can also build agents with Gemini using these frameworks and SDKs:
>
>  
>  LangChain / LangGraph : Build
> stateful, complex application flows and multi-agent systems using graph
> structures.
>
>  LlamaIndex : Connect Gemini agents to
> your private data for RAG-enhanced workflows.
>
>  CrewAI : Orchestrate collaborative,
> role-playing autonomous AI agents.
>
>  Vercel AI SDK : Build
> AI-powered user interfaces and agents in JavaScript/TypeScript.
>
>  Google ADK : An
> open-source framework for building and orchestrating interoperable AI
> agents.
>
>  Antigravity SDK : Build
> autonomous AI agents using the same tools, agent loop, and context
> management that power Google Antigravity, programmable in Python.
>
>  
>  
>  
>
>  
>  
>
>  
>  
>  
>  
>  
>  
>  
>  
>  
>
>  
>  
>  
>  Send feedback
>  
>  
>  
>  
>  
>  
>  
>
>  
>
>  
>  Except as otherwise noted, the content of this page is licensed under the Creative Commons Attribution 4.0 License , and code samples are licensed under the Apache 2.0 License . For details, see the Google Developers Site Policies . Java is a registered trademark of Oracle and/or its affiliates.
>
>  Last updated 2026-08-18 UTC.
>
>  
>
>  
>  
>
>  
>  
>  
>  
>  
>  
>  
>
>  
>  
>  Need to tell us more?
>  
>  
>  
>  
>  
>  
>  
>  [[["Easy to understand","easyToUnderstand","thumb-up"],["Solved my problem","solvedMyProblem","thumb-up"],["Other","otherUp","thumb-up"]],[["Missing the information I need","missingTheInformationINeed","thumb-down"],["Too complicated / too many steps","tooComplicatedTooManySteps","thumb-down"],["Out of date","outOfDate","thumb-down"],["Samples / code issue","samplesCodeIssue","thumb-down"],["Other","otherDown","thumb-down"]],["Last updated 2026-08-18 UTC."],[],[]]
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  
>
>  
>  
>
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  Terms
>  
>  
>  
>
>  
>  
>  
>  
>  
>  Privacy
>  
>  
>  
>
>  
>  
>  
>  
>  
>  Manage cookies
>  
>  
>  
>
>  
>  
>  
>  
>  
>  
>  
>  
>  
>  English 
>  
>
>  
>  
>  Deutsch 
>  
>
>  
>  
>  Español – América Latina 
>  
>
>  
>  
>  Français 
>  
>
>  
>  
>  Indonesia 
>  
>
>  
>  
>  Italiano 
>  
>
>  
>  
>  Polski 
>  
>
>  
>  
>  Português – Brasil 
>  
>
>  
>  
>  Shqip 
>  
>
>  
>  
>  Tiếng Việt 
>  
>
>  
>  
>  Türkçe 
>  
>
>  
>  
>  Русский 
>  
>
>  
>  
>  עברית 
>  
>
>  
>  
>  العربيّة 
>  
>
>  
>  
>  فارسی 
>  
>
>  
>  
>  हिंदी 
>  
>
>  
>  
>  বাংলা 
>  
>
>  
>  
>  ภาษาไทย 
>  
>
>  
>  
>  中文 – 简体 
>  
>
>  
>  
>  中文 – 繁體 
>  
>
>  
>  
>  日本語 
>  
>
>  
>  
>  한국어
