---
primary_sources:
  - id: T2-MONITORING
    title: "Monitoring"
    url: "https://code.claude.com/docs/en/monitoring-usage.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Monitoring and OTEL

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Monitoring

> # Monitoring
>
> > Learn how to enable and configure OpenTelemetry for Claude Code.
>
> Track Claude Code usage, costs, and tool activity across your organization by exporting telemetry data through OpenTelemetry (OTel). Claude Code exports metrics as time series data via the standard metrics protocol, events via the logs/events protocol, and optionally distributed traces via the [traces protocol](#traces-beta).
>
> ## Quick start
>
> Configure OpenTelemetry using environment variables:
>
> ```bash
> # 1. Enable telemetry
> export CLAUDE_CODE_ENABLE_TELEMETRY=1
>
> # 2. Choose exporters (both are optional - configure only what you need)
> export OTEL_METRICS_EXPORTER=otlp       # Options: otlp, prometheus, console, none
> export OTEL_LOGS_EXPORTER=otlp          # Options: otlp, console, none
>
> # 3. Configure OTLP endpoint (for OTLP exporter)
> export OTEL_EXPORTER_OTLP_PROTOCOL=grpc
> export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
>
> # 4. Set authentication (if required)
> export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Bearer your-token"
>
> # 5. For debugging: reduce export intervals, and reset them for production use
> export OTEL_METRIC_EXPORT_INTERVAL=10000  # 10 seconds (default: 60000ms)
> export OTEL_LOGS_EXPORT_INTERVAL=5000     # 5 seconds (default: 5000ms)
>
> # 6. Run Claude Code
> claude
> ```
>
> To verify a setup that exports metrics, check your backend for the `claude_code.session.count` metric, which Claude Code emits when a session starts. To verify a logs-only setup, submit a prompt and check for the `claude_code.user_prompt` event. If nothing arrives, run `claude --debug` and check the debug log for OTel export errors.
>
> For full configuration options, see the [OpenTelemetry specification](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/exporter.md#configuration-options).
>
> ## Administrator configuration
>
> Administrators can configure OpenTelemetry settings for all users through the [managed settings file](/docs/en/managed-settings#delivery-mechanisms). See the [settings precedence](/docs/en/settings#settings-precedence) for more information about how settings are applied.
>
> Example managed settings configuration:
>
> ```json
> {
>   "env": {
>     "CLAUDE_CODE_ENABLE_TELEMETRY": "1",
>     "OTEL_METRICS_EXPORTER": "otlp",
>     "OTEL_LOGS_EXPORTER": "otlp",
>     "OTEL_EXPORTER_OTLP_PROTOCOL": "grpc",
>     "OTEL_EXPORTER_OTLP_ENDPOINT": "http://collector.example.com:4317",
>     "OTEL_EXPORTER_OTLP_HEADERS": "Authorization=Bearer example-token"
>   }
> }
> ```
>
> Claude Code doesn't pass `OTEL_*` environment variables to the subprocesses it spawns, including the Bash tool, hooks, MCP servers, and language servers. An OpenTelemetry-instrumented application that you run through the Bash tool doesn't inherit Claude Code's exporter endpoint or headers, so set those variables directly in the command if that application needs to export its own telemetry.
>
> ### How managed settings lock the OTLP destination
>
> When you set an `OTEL_EXPORTER_OTLP_*` variable in managed settings, Claude Code removes conflicting developer-set variables at startup and logs a warning you can see with `claude --debug`. What it removes depends on which variable you set:
>
> * **Endpoints**: when you set `OTEL_EXPORTER_OTLP_ENDPOINT`, Claude Code removes every developer-set per-signal endpoint. Developers can't point one signal at a different collector, so you don't need to also set the per-signal endpoint variables in managed settings.
> * **Protocols**: when you set `OTEL_EXPORTER_OTLP_PROTOCOL`, Claude Code removes every developer-set per-signal protocol.
> * **Credentials**: when you set `OTEL_EXPORTER_OTLP_HEADERS`, `OTEL_EXPORTER_OTLP_CLIENT_KEY`, or `OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE`, Claude Code removes the developer-set per-signal versions of that variable, plus every developer-set endpoint variable, generic or per-signal, since those credentials would otherwise reach a collector the managed settings didn't choose.
> * **Exporter selectors**: `OTEL_METRICS_EXPORTER`, `OTEL_LOGS_EXPORTER`, and the beta `OTEL_TRACES_EXPORTER` follow normal per-key precedence. A developer's setting can still disable a signal or switch it to the console exporter, so set the selectors in managed settings too if you need them locked. Across [admin sources](/docs/en/managed-settings#precedence-within-the-managed-tier), `OTEL_LOGS_EXPORTER` follows the [telemetry unit](/docs/en/server-managed-settings#per-key-exceptions-across-managed-sources) while the other two selectors merge per key. Requires Claude Code v2.1.223 or later.
>
> Claude Code doesn't remove per-signal variables that you set in managed settings itself, so you can route one signal to a different collector by setting its variable there, as the [SIEM example](#send-events-to-a-siem) does. If you set a per-signal credential there, Claude Code removes the developer-set endpoint for that signal.
>
> This removal behavior changes where telemetry is delivered, not what Claude Code collects.
>
> Before v2.1.217, every variable followed per-key settings precedence independently, so a signal-specific endpoint set in user settings or the shell redirected that signal away from the managed collector.
>
> ## Configuration details
>
> ### Common configuration variables
>
> These variables configure exporters, endpoints, and export behavior for all deployments. If you set a per-signal endpoint or protocol variable, such as `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT`, Claude Code uses it instead of the generic variable for that signal. If you set a per-signal headers variable, such as `OTEL_EXPORTER_OTLP_METRICS_HEADERS`, Claude Code merges it with the generic `OTEL_EXPORTER_OTLP_HEADERS` for that signal. On machines with managed settings, see [How managed settings lock the OTLP destination](#how-managed-settings-lock-the-otlp-destination) for what Claude Code removes.
>
> | Environment Variable                                | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      | Example Values                                                                                                                                                 |
> | --------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `CLAUDE_CODE_ENABLE_TELEMETRY`                      | Enables telemetry collection (required)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          | `1`                                                                                                                                                            |
> | `OTEL_METRICS_EXPORTER`                             | Metrics exporter types, comma-separated. Use `none` to disable                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | `console`, `otlp`, `prometheus`, `none`                                                                                                                        |
> | `OTEL_LOGS_EXPORTER`                                | Logs/events exporter types, comma-separated. Use `none` to disable                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               | `console`, `otlp`, `none`                                                                                                                                      |
> | `OTEL_EXPORTER_OTLP_PROTOCOL`                       | Protocol for OTLP exporter, applies to all signals. Claude Code has no default protocol, so set this or the signal-specific protocol variable for each `otlp` exporter you enable                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                | `grpc`, `http/json`, `http/protobuf`                                                                                                                           |
> | `OTEL_EXPORTER_OTLP_ENDPOINT`                       | OTLP collector endpoint for all signals                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          | `http://localhost:4317`                                                                                                                                        |
> | `OTEL_EXPORTER_OTLP_METRICS_PROTOCOL`               | Protocol for metrics, overrides general setting                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | `grpc`, `http/json`, `http/protobuf`                                                                                                                           |
> | `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT`               | OTLP metrics endpoint, overrides general setting                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | `http://localhost:4318/v1/metrics`                                                                                                                             |
> | `OTEL_EXPORTER_OTLP_LOGS_PROTOCOL`                  | Protocol for logs, overrides general setting                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | `grpc`, `http/json`, `http/protobuf`                                                                                                                           |
> | `OTEL_EXPORTER_OTLP_LOGS_ENDPOINT`                  | OTLP logs endpoint, overrides general setting                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | `http://localhost:4318/v1/logs`                                                                                                                                |
> | `OTEL_EXPORTER_OTLP_HEADERS`                        | Authentication headers for OTLP                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | `Authorization=Bearer token`                                                                                                                                   |
> | `OTEL_EXPORTER_OTLP_METRICS_HEADERS`                | Authentication headers for metrics, merged with the general headers                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | `Authorization=Bearer token`                                                                                                                                   |
> | `OTEL_EXPORTER_OTLP_LOGS_HEADERS`                   | Authentication headers for logs, merged with the general headers                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | `Authorization=Bearer token`                                                                                                                                   |
> | `OTEL_METRIC_EXPORT_INTERVAL`                       | Export interval in milliseconds (default: 60000)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | `5000`, `60000`                                                                                                                                                |
> | `OTEL_LOGS_EXPORT_INTERVAL`                         | Logs export interval in milliseconds (default: 5000)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | `1000`, `10000`                                                                                                                                                |
> | `OTEL_LOG_USER_PROMPTS`                             | Enable logging of user prompt content (default: disabled)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | `1` to enable                                                                                                                                                  |
> | `OTEL_LOG_ASSISTANT_RESPONSES`                      | Enable logging of assistant response text on `assistant_response` events (default: disabled). When unset, falls back to the value of `OTEL_LOG_USER_PROMPTS`. Requires Claude Code v2.1.193 or later                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | `1` to enable, `0` to keep redacted                                                                                                                            |
> | `OTEL_LOG_TOOL_DETAILS`                             | Enable logging of tool parameters and input arguments in tool events and trace span attributes: Bash commands, MCP server and tool names, skill names, user-authored workflow names, and tool input. Also enables custom, plugin, and MCP command names on `user_prompt` events (default: disabled). For Claude Desktop's built-in servers, in sessions Claude Desktop owns, `mcp_server_name`/`mcp_tool_name` emit on `tool_decision`/`tool_result` even with the flag off. The exception requires Claude Code v2.1.214 or later                                                                                                                                                | `1` to enable                                                                                                                                                  |
> | `OTEL_LOG_TOOL_CONTENT`                             | Enable logging of tool input and output content in span events (default: disabled). Requires [tracing](#traces-beta). Content is truncated at the content limit (60 KB by default)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               | `1` to enable                                                                                                                                                  |
> | `OTEL_LOG_RAW_API_BODIES`                           | Emit the full Anthropic Messages API request and response JSON as `api_request_body` / `api_response_body` log events (default: disabled). Bodies include the entire conversation history. Enabling this implies consent to everything `OTEL_LOG_USER_PROMPTS`, `OTEL_LOG_TOOL_DETAILS`, and `OTEL_LOG_TOOL_CONTENT` would reveal                                                                                                                                                                                                                                                                                                                                                | `1` for inline bodies truncated at the content limit (60 KB by default), or `file:<dir>` for untruncated bodies on disk with a `body_ref` pointer in the event |
> | `CLAUDE_CODE_OTEL_CONTENT_MAX_LENGTH`               | Content limit: the maximum length of content-bearing attributes such as model responses, tool content, system prompts, and raw API bodies, truncation marker included, in UTF-16 code units (default: 61440, i.e. 60 KB). The default is sized for backends that cap attribute values at 64 KB; raise it only if your backend accepts larger values, or lower it to cut telemetry volume. When an OpenTelemetry SDK attribute limit, `OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT` or one of its logrecord and span variants, is set lower, Claude Code truncates at that smaller value so the `[TRUNCATED ...]` marker stays within the SDK limit. Requires Claude Code v2.1.214 or later | `262144`                                                                                                                                                       |
> | `OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE` | Metrics temporality preference (default: `delta`). Set to `cumulative` if your backend expects cumulative temporality                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            | `delta`, `cumulative`                                                                                                                                          |
> | `CLAUDE_CODE_OTEL_HEADERS_HELPER_DEBOUNCE_MS`       | Interval for refreshing dynamic headers (default: 1740000ms / 29 minutes)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | `900000`                                                                                                                                                       |
>
> For the `http/protobuf` and `http/json` protocols, Claude Code sends each export request with a `Content-Length` header. Before v2.1.212, Claude Code versions from v2.1.191 onward sent these requests with chunked transfer encoding; Azure Monitor and other endpoints that require a declared length rejected them with `411 Length Required` or `400` errors.
>
> ### mTLS authentication
>
> How you configure client certificates for the OTLP exporter depends on the OTLP protocol in use for that signal, set via `OTEL_EXPORTER_OTLP_PROTOCOL` or the per-signal override. The same configuration applies to metrics, logs, and traces.
>
> | Protocol                     | Client certificate variables                                                                                                                                                                      | Trust the collector's CA with    |
> | :--------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | :------------------------------- |
> | `http/protobuf`, `http/json` | `CLAUDE_CODE_CLIENT_CERT`, `CLAUDE_CODE_CLIENT_KEY`, and optionally `CLAUDE_CODE_CLIENT_KEY_PASSPHRASE`. See [Network configuration](/docs/en/network-config#mtls-authentication)                      | `NODE_EXTRA_CA_CERTS`            |
> | `grpc`                       | `OTEL_EXPORTER_OTLP_CLIENT_KEY` and `OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE`, or the per-signal variants such as `OTEL_EXPORTER_OTLP_METRICS_CLIENT_KEY` to use a different certificate per signal | `OTEL_EXPORTER_OTLP_CERTIFICATE` |
>
> For `grpc`, the OpenTelemetry SDK reads the standard OTLP variables directly, so existing configurations that set the per-signal metrics variables continue to work. On machines with managed settings, Claude Code [may remove developer-set per-signal credentials and endpoints](#how-managed-settings-lock-the-otlp-destination) at startup.
>
> ### Metrics cardinality control
>
> The following environment variables control which attributes are included in metrics to manage cardinality:
>
> | Environment Variable                       | Description                                                                     | Default Value | Example to Disable |
> | ------------------------------------------ | ------------------------------------------------------------------------------- | ------------- | ------------------ |
> | `OTEL_METRICS_INCLUDE_SESSION_ID`          | Include session.id attribute in metrics                                         | `true`        | `false`            |
> | `OTEL_METRICS_INCLUDE_VERSION`             | Include app.version attribute in metrics                                        | `false`       | `true`             |
> | `OTEL_METRICS_INCLUDE_ACCOUNT_UUID`        | Include user.account\_uuid and user.account\_id attributes in metrics           | `true`        | `false`            |
> | `OTEL_METRICS_INCLUDE_ENTRYPOINT`          | Include app.entrypoint attribute in metrics                                     | `false`       | `true`             |
> | `OTEL_METRICS_INCLUDE_RESOURCE_ATTRIBUTES` | Include keys from `OTEL_RESOURCE_ATTRIBUTES` as attributes on metric datapoints | `true`        | `false`            |
>
> Lower cardinality generally means better performance and lower storage costs but less granular data for analysis.
>
> ### Traces (beta)
>
> Distributed tracing exports spans that link each user prompt to the API requests and tool executions it triggers, so you can view a full request as a single trace in your tracing backend.
>
> Tracing is off by default. To enable it, set both `CLAUDE_CODE_ENABLE_TELEMETRY=1` and `CLAUDE_CODE_ENHANCED_TELEMETRY_BETA=1`, then set `OTEL_TRACES_EXPORTER` to choose where spans are sent. Traces reuse the [common OTLP configuration](#common-configuration-variables) for endpoint, protocol, headers, and [mTLS](#mtls-authentication). On machines with managed settings, Claude Code [may remove developer-set per-signal credentials and endpoints](#how-managed-settings-lock-the-otlp-destination) at startup.
>
> | Environment Variable                  | Description                                                                       | Example Values                       |
> | ------------------------------------- | --------------------------------------------------------------------------------- | ------------------------------------ |
> | `CLAUDE_CODE_ENHANCED_TELEMETRY_BETA` | Enable span tracing (required). `ENABLE_ENHANCED_TELEMETRY_BETA` is also accepted | `1`                                  |
> | `OTEL_TRACES_EXPORTER`                | Traces exporter types, comma-separated. Use `none` to disable                     | `console`, `otlp`, `none`            |
> | `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL`  | Protocol for traces, overrides `OTEL_EXPORTER_OTLP_PROTOCOL`                      | `grpc`, `http/json`, `http/protobuf` |
> | `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`  | OTLP traces endpoint, overrides `OTEL_EXPORTER_OTLP_ENDPOINT`                     | `http://localhost:4318/v1/traces`    |
> | `OTEL_EXPORTER_OTLP_TRACES_HEADERS`   | Authentication headers for traces, merged with `OTEL_EXPORTER_OTLP_HEADERS`       | `Authorization=Bearer token`         |
> | `OTEL_TRACES_EXPORT_INTERVAL`         | Span batch export interval in milliseconds (default: 5000)                        | `1000`, `10000`                      |
>
> Spans redact user prompt text, tool input details, and tool content by default. Set `OTEL_LOG_USER_PROMPTS=1`, `OTEL_LOG_TOOL_DETAILS=1`, and `OTEL_LOG_TOOL_CONTENT=1` to include them.
>
> When tracing is active, Bash and PowerShell subprocesses automatically inherit a `TRACEPARENT` environment variable containing the W3C trace context of the active tool execution span. This lets any subprocess that reads `TRACEPARENT` parent its own spans under the same trace, enabling end-to-end distributed tracing through scripts and commands that Claude runs.
>
> When tracing is active and Claude Code is connected directly to the Anthropic API, each model request carries a W3C `traceparent` header set to the `claude_code.llm_request` span's context, and the API's `traceresponse` header is recorded as a span link. Together these connect Claude Code's client-side spans to the server-side trace through any compliant intermediary. Outbound HTTP MCP requests carry `traceparent` the same way. The header is not sent to third-party providers.
>
> By default, the `traceparent` header on model and HTTP MCP requests is sent only when `ANTHROPIC_BASE_URL` is unset or points at the Anthropic API, since some proxies reject unrecognized headers. The subprocess `TRACEPARENT` variable is controlled by the same switch for consistency. If you run Claude Code through a custom `ANTHROPIC_BASE_URL` proxy and want trace context propagated, set `CLAUDE_CODE_PROPAGATE_TRACEPARENT=1`.
>
> In Agent SDK and non-interactive sessions started with `-p`, Claude Code also reads `TRACEPARENT` and `TRACESTATE` from its own environment when starting each interaction span. This lets an embedding process pass its active W3C trace context into the subprocess so Claude Code's spans appear as children of the caller's distributed trace. Interactive sessions ignore inbound `TRACEPARENT` to avoid accidentally inheriting ambient values from CI or container environments.
>
> The inbound trace context also applies to [events](#events). In Agent SDK and `-p` sessions with `TRACEPARENT` set, each OTLP event log record carries `trace_id` and `span_id` values that join it to your application's trace, even when the traces exporter isn't configured, so your logging backend can correlate events with the rest of the trace.
>
> A record emitted while an interaction is active carries the interaction span's IDs, even when Claude Code emits it outside the span's async context, such as in a permission prompt callback or for a record buffered during startup and exported later. A record emitted with no active interaction span carries the inbound `TRACEPARENT` IDs directly. Before v2.1.214, records emitted outside the span's async context carried the inbound `TRACEPARENT` IDs instead of the span's IDs. Before v2.1.212, event records emitted outside an active span didn't carry `trace_id` or `span_id`.
>
> #### Span hierarchy
>
> Each user prompt starts a `claude_code.interaction` root span. API calls, tool calls, and hook executions are recorded as its children. Tool spans have two child spans of their own: one for the time spent waiting on a permission decision and one for the execution itself. When the Agent tool, or legacy Task tool, spawns a subagent, the subagent's API and tool spans nest under the parent's `claude_code.tool` span.
>
> ```text
> claude_code.interaction
> ├── claude_code.llm_request
> ├── claude_code.hook                    (requires detailed beta tracing)
> └── claude_code.tool
>     ├── claude_code.tool.blocked_on_user
>     ├── claude_code.tool.execution
>     └── (Agent tool) subagent claude_code.llm_request / claude_code.tool spans
> ```
>
> In Agent SDK and `claude -p` sessions, `claude_code.interaction` itself becomes a child of the caller's span when `TRACEPARENT` is set in the environment.
>
> #### Span attributes
>
> Every span carries the [standard attributes](#standard-attributes) plus a `span.type` attribute matching its name. The tables below list the additional attributes set on each span. The `llm_request`, `tool.execution`, and `hook` spans set OpenTelemetry status `ERROR` when they record a failure; the other spans always end with status `UNSET`.
>
> **`claude_code.interaction`**
>
> | Attribute                 | Description                                               | Gated by                |
> | ------------------------- | --------------------------------------------------------- | ----------------------- |
> | `user_prompt`             | Prompt text. Value is `<REDACTED>` unless the gate is set | `OTEL_LOG_USER_PROMPTS` |
> | `user_prompt_length`      | Prompt length in characters                               |                         |
> | `interaction.sequence`    | 1-based counter of interactions in this session           |                         |
> | `interaction.duration_ms` | Wall-clock duration of the turn                           |                         |
>
> **`claude_code.llm_request`**
>
> | Attribute                        | Description                                                                                                                                   | Gated by                |
> | -------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------- |
> | `model`                          | Model identifier                                                                                                                              |                         |
> | `gen_ai.system`                  | Always `anthropic`. OpenTelemetry GenAI semantic convention                                                                                   |                         |
> | `gen_ai.request.model`           | Same value as `model`. OpenTelemetry GenAI semantic convention                                                                                |                         |
> | `query_source`                   | Subsystem that issued the request, such as `repl_main_thread` or a subagent name                                                              |                         |
> | `agent_id`                       | Identifier of the subagent or teammate that issued the request. Absent on the main session                                                    |                         |
> | `parent_agent_id`                | Identifier of the agent that spawned this one. Absent for the main session and for agents spawned directly from it                            |                         |
> | `workflow.run_id`                | Run identifier of the [Workflow](/docs/en/workflows) tool run that spawned this agent, prefixed `wf_`. Absent for agents not spawned by a workflow |                         |
> | `workflow.name`                  | Name of the workflow that spawned this agent. User-authored names are replaced with `custom` unless the gate is set                           | `OTEL_LOG_TOOL_DETAILS` |
> | `speed`                          | `fast` or `normal`                                                                                                                            |                         |
> | `llm_request.context`            | `interaction`, `tool`, or `standalone` depending on the parent span                                                                           |                         |
> | `duration_ms`                    | Wall-clock duration including retries                                                                                                         |                         |
> | `ttft_ms`                        | Time to first token in milliseconds                                                                                                           |                         |
> | `input_tokens`                   | Input token count from the API usage block                                                                                                    |                         |
> | `output_tokens`                  | Output token count                                                                                                                            |                         |
> | `cache_read_tokens`              | Tokens read from prompt cache                                                                                                                 |                         |
> | `cache_creation_tokens`          | Tokens written to prompt cache                                                                                                                |                         |
> | `request_id`                     | Anthropic API request ID from the `request-id` response header                                                                                |                         |
> | `gen_ai.response.id`             | Same value as `request_id`. OpenTelemetry GenAI semantic convention                                                                           |                         |
> | `client_request_id`              | Client-generated `x-client-request-id` of the final attempt                                                                                   |                         |
> | `attempt`                        | Total attempts made for this request                                                                                                          |                         |
> | `success`                        | `true` or `false`                                                                                                                             |                         |
> | `status_code`                    | HTTP status code when the request failed                                                                                                      |                         |
> | `error`                          | Error message when the request failed                                                                                                         |                         |
> | `response.has_tool_call`         | `true` when the response contained tool-use blocks                                                                                            |                         |
> | `stop_reason`                    | API response `stop_reason`, such as `end_turn`, `tool_use`, `max_tokens`, `stop_sequence`, `pause_turn`, or `refusal`                         |                         |
> | `gen_ai.response.finish_reasons` | Same value as `stop_reason`, wrapped in a string array. OpenTelemetry GenAI semantic convention                                               |                         |
>
> Each retry attempt is also recorded as a `gen_ai.request.attempt` span event with `attempt` and `client_request_id` attributes.
>
> **`claude_code.tool`**
>
> | Attribute             | Description                                                                                                                                                                                                                          | Gated by                |
> | --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------- |
> | `tool_name`           | Tool name                                                                                                                                                                                                                            |                         |
> | `duration_ms`         | Wall-clock duration including permission wait and execution                                                                                                                                                                          |                         |
> | `result_tokens`       | Approximate token size of the tool result                                                                                                                                                                                            |                         |
> | `agent_id`            | Identifier of the subagent or teammate that ran the tool. Absent on the main session                                                                                                                                                 |                         |
> | `parent_agent_id`     | Identifier of the agent that spawned this one. Absent for the main session and for agents spawned directly from it                                                                                                                   |                         |
> | `workflow.run_id`     | Run identifier of the Workflow tool run that spawned this agent, prefixed `wf_`. Absent for agents not spawned by a workflow                                                                                                         |                         |
> | `workflow.name`       | Name of the workflow that spawned this agent. User-authored names are replaced with `custom` unless the gate is set                                                                                                                  | `OTEL_LOG_TOOL_DETAILS` |
> | `tool_use_id`         | The model's `tool_use` block id for this call. Matches the `tool_use_id` on the [tool\_result](#tool-result-event) and [tool\_decision](#tool-decision-event) events and in hook payloads, so you can join the span to those records |                         |
> | `gen_ai.tool.call.id` | Same value as `tool_use_id`. OpenTelemetry GenAI semantic convention                                                                                                                                                                 |                         |
> | `file_path`           | Target file path for Read, Edit, and Write tools                                                                                                                                                                                     | `OTEL_LOG_TOOL_DETAILS` |
> | `full_command`        | Command string for the Bash tool                                                                                                                                                                                                     | `OTEL_LOG_TOOL_DETAILS` |
> | `skill_name`          | Skill name for the Skill tool                                                                                                                                                                                                        | `OTEL_LOG_TOOL_DETAILS` |
> | `subagent_type`       | Subagent type for the Agent tool or legacy Task tool                                                                                                                                                                                 | `OTEL_LOG_TOOL_DETAILS` |
>
> When `OTEL_LOG_TOOL_CONTENT=1`, this span also records a `tool.output` span event whose attributes contain the tool's input and output bodies, truncated at the content limit (60 KB by default) per attribute.
>
> **`claude_code.tool.blocked_on_user`**
>
> | Attribute     | Description                                                               | Gated by |
> | ------------- | ------------------------------------------------------------------------- | -------- |
> | `duration_ms` | Time spent waiting for the permission decision                            |          |
> | `decision`    | `accept` or `reject`                                                      |          |
> | `source`      | Decision source, matching the [Tool decision event](#tool-decision-event) |          |
>
> **`claude_code.tool.execution`**
>
> | Attribute             | Description                                                                                                                                       | Gated by                |
> | --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------- |
> | `duration_ms`         | Time spent running the tool body                                                                                                                  |                         |
> | `tool_use_id`         | Same value as on the parent `claude_code.tool` span                                                                                               |                         |
> | `gen_ai.tool.call.id` | Same value as `tool_use_id`. OpenTelemetry GenAI semantic convention                                                                              |                         |
> | `success`             | `true` or `false`                                                                                                                                 |                         |
> | `error`               | Error category string when execution failed, such as `Error:ENOENT` or `ShellError`. Contains the full error message instead when the gate is set | `OTEL_LOG_TOOL_DETAILS` |
>
> **`claude_code.hook`**
>
> This span is emitted only when detailed beta tracing is active, which requires `ENABLE_BETA_TRACING_DETAILED=1` and `BETA_TRACING_ENDPOINT` in addition to the trace exporter configuration above. In interactive CLI sessions, this also requires your organization to be allowlisted for the feature. Agent SDK and non-interactive `-p` sessions are not gated. It is not emitted when only `CLAUDE_CODE_ENHANCED_TELEMETRY_BETA` is set.
>
> | Attribute                | Description                                      | Gated by                |
> | ------------------------ | ------------------------------------------------ | ----------------------- |
> | `hook_event`             | Hook event type, such as `PreToolUse`            |                         |
> | `hook_name`              | Full hook name, such as `PreToolUse:Write`       |                         |
> | `num_hooks`              | Number of matching hook commands executed        |                         |
> | `hook_definitions`       | JSON-serialized hook configuration               | `OTEL_LOG_TOOL_DETAILS` |
> | `duration_ms`            | Wall-clock duration of all matching hooks        |                         |
> | `num_success`            | Count of hooks that completed successfully       |                         |
> | `num_blocking`           | Count of hooks that returned a blocking decision |                         |
> | `num_non_blocking_error` | Count of hooks that failed without blocking      |                         |
> | `num_cancelled`          | Count of hooks cancelled before completion       |                         |
>
>
>   Additional content-bearing attributes such as `new_context`, `system_prompt_preview`, `user_system_prompt`, `tool_input`, and `response.model_output` are emitted only when detailed beta tracing is active. They are not part of the stable span schema.
>
>   `user_system_prompt` additionally requires `OTEL_LOG_USER_PROMPTS=1`. It carries only the system prompt text you provide via the `systemPrompt` SDK option or `--system-prompt` and `--append-system-prompt` flags, truncated at the content limit (60 KB by default), and is emitted once per session rather than per request.
>
>
> ### Dynamic headers
>
> For enterprise environments that require dynamic authentication, you can configure a script to generate headers dynamically. Dynamic headers apply only to the `http/protobuf` and `http/json` protocols. With the `grpc` protocol, Claude Code uses only the static headers variables, `OTEL_EXPORTER_OTLP_HEADERS` and its per-signal variants.
>
> #### Settings configuration
>
> Add to your `.claude/settings.json`, replacing the path with your own script:
>
> ```json
> {
>   "otelHeadersHelper": "/path/to/generate-otel-headers.sh"
> }
> ```
>
> The value can be the path to an executable file, including a path that contains spaces, or a shell command line with arguments. On Windows, the value always runs through the shell, so quote a path that contains spaces inside the JSON value.
>
> #### Script requirements
>
> The script must output valid JSON with string key-value pairs representing HTTP headers:
>
> ```bash
> #!/bin/bash
> # Example: Multiple headers
> echo "{\"Authorization\": \"Bearer $(get-token.sh)\", \"X-API-Key\": \"$(get-api-key.sh)\"}"
> ```
>
> If the helper fails or prints output that doesn't meet these requirements, Claude Code reports the error in:
>
> * `/status` output
> * The debug log, when running with [`--debug`](/docs/en/cli-reference#cli-flags) or after running `/debug` in the session
> * stderr, in non-interactive sessions started with `-p`
>
> #### Refresh behavior
>
> The headers helper script runs at startup and periodically thereafter to support token refresh. By default, the script runs every 29 minutes. Customize the interval with the `CLAUDE_CODE_OTEL_HEADERS_HELPER_DEBOUNCE_MS` environment variable.
>
> ### Multi-team organization support
>
> Organizations with multiple teams or departments can add custom attributes to distinguish between different groups using the `OTEL_RESOURCE_ATTRIBUTES` environment variable:
>
> ```bash
> # Add custom attributes for team identification
> export OTEL_RESOURCE_ATTRIBUTES="department=engineering,team.id=platform,cost_center=eng-123"
> ```
>
> These custom attributes are included in all metrics and events, allowing you to:
>
> * Filter metrics by team or department
> * Track costs per cost center
> * Create team-specific dashboards
> * Set up alerts for specific teams
>
> Claude Code attaches these values as attributes on every metric datapoint and event record, in addition to sending them in the OTLP resource block. Because most metrics backends expose datapoint attributes as queryable labels, you can group and filter metrics by your custom keys directly. Custom keys never override the [standard attributes](#standard-attributes) such as `user.id` or `session.id`: when a key collides, Claude Code keeps the built-in value.
>
> Each custom key becomes a label on every metric series, so high-cardinality values increase storage cost in your metrics backend. To send custom attributes in the resource block only and omit them from datapoint labels, set `OTEL_METRICS_INCLUDE_RESOURCE_ATTRIBUTES=false`. See [Metrics cardinality control](#metrics-cardinality-control).
>
>
>   The `OTEL_RESOURCE_ATTRIBUTES` environment variable uses comma-separated key=value pairs with strict formatting requirements:
>
>   * **No spaces allowed**: values can't contain spaces. For example, `user.organizationName=My Company` is invalid
>   * **Format**: must be comma-separated key=value pairs: `key1=value1,key2=value2`
>   * **Allowed characters**: only US-ASCII characters excluding control characters, whitespace, double quotes, commas, semicolons, and backslashes
>   * **Special characters**: characters outside the allowed range must be percent-encoded
>
>   For a value that would need a space, use underscores or camelCase instead. The following examples set `org.name` with each form:
>
>   ```bash
>   export OTEL_RESOURCE_ATTRIBUTES="org.name=Johns_Organization"
>   export OTEL_RESOURCE_ATTRIBUTES="org.name=JohnsOrganization"
>   ```
>
>   You can percent-encode any character, not only the excluded ones. This example encodes both the space and the apostrophe:
>
>   ```bash
>   export OTEL_RESOURCE_ATTRIBUTES="org.name=John%27s%20Organization"
>   ```
>
>   Wrapping values in quotes doesn't escape spaces. For example, `org.name="My Company"` results in the literal value `"My Company"` with the quotes included, not `My Company`.
>
>
> ### Example configurations
>
> Set these environment variables before running `claude`. Each scenario below shows a complete configuration, and each variable is described under [Common configuration variables](#common-configuration-variables). To confirm a configuration took effect, check your backend for the `claude_code.session.count` metric after starting a session; the [Quick start](#quick-start) covers logs-only verification and what to check when nothing arrives.
>
> For console debugging with a 1-second export interval:
>
> ```bash
> export CLAUDE_CODE_ENABLE_TELEMETRY=1
> export OTEL_METRICS_EXPORTER=console
> export OTEL_METRIC_EXPORT_INTERVAL=1000
> ```
>
> For OTLP over gRPC:
>
> ```bash
> export CLAUDE_CODE_ENABLE_TELEMETRY=1
> export OTEL_METRICS_EXPORTER=otlp
> export OTEL_EXPORTER_OTLP_PROTOCOL=grpc
> export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
> ```
>
> For Prometheus, scraped from `http://localhost:9464/metrics`:
>
> ```bash
> export CLAUDE_CODE_ENABLE_TELEMETRY=1
> export OTEL_METRICS_EXPORTER=prometheus
> ```
>
> On a [self-hosted environment](/docs/en/self-hosted-environments-reference#pass-through-session-child-metrics), the session binds port 9464 only at the runner's default capacity of one. At higher capacity, the runner re-exposes session counters and gauges on its own `/metrics` endpoint instead.
>
> To send metrics to multiple exporters:
>
> ```bash
> export CLAUDE_CODE_ENABLE_TELEMETRY=1
> export OTEL_METRICS_EXPORTER=console,otlp
> export OTEL_EXPORTER_OTLP_PROTOCOL=http/json
> ```
>
> To send metrics and logs to different endpoints or backends:
>
> ```bash
> export CLAUDE_CODE_ENABLE_TELEMETRY=1
> export OTEL_METRICS_EXPORTER=otlp
> export OTEL_LOGS_EXPORTER=otlp
> export OTEL_EXPORTER_OTLP_METRICS_PROTOCOL=http/protobuf
> export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://metrics.example.com:4318
> export OTEL_EXPORTER_OTLP_LOGS_PROTOCOL=grpc
> export OTEL_EXPORTER_OTLP_LOGS_ENDPOINT=http://logs.example.com:4317
> ```
>
> To export metrics only, without events or logs:
>
> ```bash
> export CLAUDE_CODE_ENABLE_TELEMETRY=1
> export OTEL_METRICS_EXPORTER=otlp
> export OTEL_EXPORTER_OTLP_PROTOCOL=grpc
> export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
> ```
>
> To export events and logs only, without metrics:
>
> ```bash
> export CLAUDE_CODE_ENABLE_TELEMETRY=1
> export OTEL_LOGS_EXPORTER=otlp
> export OTEL_EXPORTER_OTLP_PROTOCOL=grpc
> export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
> ```
>
> ## Available metrics and events
>
> ### Standard attributes
>
> All metrics and events share these standard attributes:
>
> | Attribute                            | Description                                                                                                                                                                                                                          | Controlled By                                              |
> | ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------- |
> | `session.id`                         | Unique session identifier                                                                                                                                                                                                            | `OTEL_METRICS_INCLUDE_SESSION_ID` (default: true)          |
> | `app.version`                        | Current Claude Code version                                                                                                                                                                                                          | `OTEL_METRICS_INCLUDE_VERSION` (default: false)            |
> | `app.entrypoint`                     | How the session was launched, such as `cli
