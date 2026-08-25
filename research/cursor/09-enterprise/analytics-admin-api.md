---
primary_sources:
  - id: T1-ANALYTICS
    title: "Analytics"
    url: "https://cursor.com/docs/account/teams/analytics-api.md"
    section: "Analytics"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Analytics and Admin API

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Analytics API

> # Analytics API
>
> The Analytics API provides comprehensive insights into your team's Cursor usage, including AI-assisted coding metrics, active users, model usage, and more.
>
> - The Analytics API uses [Basic Authentication](https://cursor.com/docs/api.md#basic-authentication). Most endpoints require an admin-scoped API key with `admin:*` scope. Bugbot review analytics require `read:*` scope. Generate a key from [Cursor Dashboard → API Keys](https://cursor.com/dashboard/api).
> - For details on authentication, rate limits, and best practices, see the [API Overview](https://cursor.com/docs/api.md).
> - **Availability**: Only for enterprise teams
>
> ### Available Endpoints
>
> ### Agent Edits
>
> /analytics/team/agent-edits
>
> Get metrics on AI-suggested code edits accepted by your team with Cursor.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/agent-edits" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "event_date": "2025-01-15",
>       "total_suggested_diffs": 145,
>       "total_accepted_diffs": 98,
>       "total_rejected_diffs": 47,
>       "total_green_lines_accepted": 820,
>       "total_red_lines_accepted": 160,
>       "total_green_lines_rejected": 210,
>       "total_red_lines_rejected": 60,
>       "total_green_lines_suggested": 1030,
>       "total_red_lines_suggested": 220,
>       "total_lines_suggested": 1250,
>       "total_lines_accepted": 980
>     },
>     {
>       "event_date": "2025-01-16",
>       "total_suggested_diffs": 132,
>       "total_accepted_diffs": 89,
>       "total_rejected_diffs": 43,
>       "total_green_lines_accepted": 740,
>       "total_red_lines_accepted": 150,
>       "total_green_lines_rejected": 185,
>       "total_red_lines_rejected": 55,
>       "total_green_lines_suggested": 925,
>       "total_red_lines_suggested": 175,
>       "total_lines_suggested": 1100,
>       "total_lines_accepted": 890
>     }
>   ],
>   "params": {
>     "metric": "agent-edits",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Tab Usage
>
> /analytics/team/tabs
>
> Get metrics on Tab autocomplete usage across your team.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/tabs" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "event_date": "2025-01-15",
>       "total_suggestions": 5420,
>       "total_accepts": 3210,
>       "total_rejects": 2210,
>       "total_green_lines_accepted": 4120,
>       "total_red_lines_accepted": 2000,
>       "total_green_lines_rejected": 1480,
>       "total_red_lines_rejected": 730,
>       "total_green_lines_suggested": 5600,
>       "total_red_lines_suggested": 2740,
>       "total_lines_suggested": 8340,
>       "total_lines_accepted": 6120
>     },
>     {
>       "event_date": "2025-01-16",
>       "total_suggestions": 4980,
>       "total_accepts": 3050,
>       "total_rejects": 1930,
>       "total_green_lines_accepted": 3890,
>       "total_red_lines_accepted": 1890,
>       "total_green_lines_rejected": 1350,
>       "total_red_lines_rejected": 580,
>       "total_green_lines_suggested": 5240,
>       "total_red_lines_suggested": 2650,
>       "total_lines_suggested": 7890,
>       "total_lines_accepted": 5780
>     }
>   ],
>   "params": {
>     "metric": "tabs",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Daily Active Users (DAU)
>
> /analytics/team/dau
>
> Get daily active user counts for your team. DAU is the number of unique users who have used Cursor in a given day.
> An active user is a user who has used at least one AI feature in Cursor.
>
> Response includes DAU breakdown metrics for the Cursor CLI, Cloud Agents, and BugBot.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/dau?startDate=14d&endDate=today" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "date": "2025-01-15",
>       "dau": 42,
>       "cli_dau": 5,
>       "cloud_agent_dau": 37,
>       "bugbot_dau": 10
>     },
>     {
>       "date": "2025-01-16",
>       "dau": 38,
>       "cli_dau": 4,
>       "cloud_agent_dau": 34,
>       "bugbot_dau": 12
>     }
>   ],
>   "params": {
>     "metric": "dau",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Client Versions
>
> /analytics/team/client-versions
>
> Get distribution of Cursor client versions used by your team (defaults to last 7 days). We report the latest version for each user per day (if a user has installed multiple versions, we report the latest).
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/client-versions" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "event_date": "2025-01-01",
>       "client_version": "0.42.3",
>       "user_count": 35,
>       "percentage": 0.833
>     },
>     {
>       "event_date": "2025-01-01",
>       "client_version": "0.42.2",
>       "user_count": 7,
>       "percentage": 0.167
>     }
>   ],
>   "params": {
>     "metric": "client-versions",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Model Usage
>
> /analytics/team/models
>
> Get metrics on AI model usage across your team.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/models" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "date": "2025-01-15",
>       "model_breakdown": {
>         "claude-sonnet-4.5": {
>           "messages": 1250,
>           "users": 28
>         },
>         "gpt-4o": {
>           "messages": 450,
>           "users": 15
>         },
>         "claude-opus-4.5": {
>           "messages": 320,
>           "users": 12
>         }
>       }
>     },
>     {
>       "date": "2025-01-16",
>       "model_breakdown": {
>         "claude-sonnet-4.5": {
>           "messages": 1180,
>           "users": 26
>         },
>         "gpt-4o": {
>           "messages": 420,
>           "users": 14
>         }
>       }
>     }
>   ],
>   "params": {
>     "metric": "models",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Top File Extensions
>
> /analytics/team/top-file-extensions
>
> Get the most frequently edited files across your team in Cursor. Returns the top 5 file extensions per day by suggestion volume.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/top-file-extensions?startDate=30d&endDate=today" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "event_date": "2025-01-15",
>       "file_extension": "tsx",
>       "total_files": 156,
>       "total_accepts": 98,
>       "total_rejects": 45,
>       "total_lines_suggested": 3230,
>       "total_lines_accepted": 2340,
>       "total_lines_rejected": 890
>     },
>     {
>       "event_date": "2025-01-15",
>       "file_extension": "ts",
>       "total_files": 142,
>       "total_accepts": 89,
>       "total_rejects": 38,
>       "total_lines_suggested": 2850,
>       "total_lines_accepted": 2100,
>       "total_lines_rejected": 750
>     }
>   ],
>   "params": {
>     "metric": "top-files",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### MCP Adoption
>
> /analytics/team/mcp
>
> Get metrics on MCP (Model Context Protocol) tool adoption across your team. Returns daily adoption counts broken down by tool name and MCP server name.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/mcp" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "event_date": "2025-01-15",
>       "tool_name": "read_file",
>       "mcp_server_name": "filesystem",
>       "usage": 245
>     },
>     {
>       "event_date": "2025-01-15",
>       "tool_name": "search_web",
>       "mcp_server_name": "brave-search",
>       "usage": 128
>     },
>     {
>       "event_date": "2025-01-16",
>       "tool_name": "read_file",
>       "mcp_server_name": "filesystem",
>       "usage": 231
>     }
>   ],
>   "params": {
>     "metric": "mcp",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Commands Adoption
>
> /analytics/team/commands
>
> Get metrics on Cursor command adoption across your team. Returns daily adoption counts broken down by command name.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/commands" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "event_date": "2025-01-15",
>       "command_name": "explain",
>       "usage": 89
>     },
>     {
>       "event_date": "2025-01-15",
>       "command_name": "refactor",
>       "usage": 45
>     },
>     {
>       "event_date": "2025-01-16",
>       "command_name": "explain",
>       "usage": 92
>     }
>   ],
>   "params": {
>     "metric": "commands",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Plans Adoption
>
> /analytics/team/plans
>
> Get metrics on Plan mode adoption across your team. Returns daily adoption counts broken down by AI model used for plan generation.
>
> The API returns `default` as the model name when a user has the Auto model selection enabled. This corresponds to what users see as "Auto" in the Cursor UI.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/plans" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "event_date": "2025-01-15",
>       "model": "claude-sonnet-4.5",
>       "usage": 156
>     },
>     {
>       "event_date": "2025-01-15",
>       "model": "default",
>       "usage": 42
>     },
>     {
>       "event_date": "2025-01-16",
>       "model": "claude-sonnet-4.5",
>       "usage": 148
>     }
>   ],
>   "params": {
>     "metric": "plans",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Skills Adoption
>
> /analytics/team/skills
>
> Get metrics on Skills adoption across your team. Returns daily adoption counts broken down by skill name.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/skills" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "event_date": "2025-01-15",
>       "skill_name": "react-best-practices",
>       "usage": 53
>     },
>     {
>       "event_date": "2025-01-15",
>       "skill_name": "usage-billing",
>       "usage": 41
>     },
>     {
>       "event_date": "2025-01-16",
>       "skill_name": "react-best-practices",
>       "usage": 48
>     }
>   ],
>   "params": {
>     "metric": "skills",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Ask Mode Adoption
>
> /analytics/team/ask-mode
>
> Get metrics on Ask mode adoption across your team. Returns daily adoption counts broken down by AI model used for Ask mode queries.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `users` string
>
> Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/ask-mode" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "event_date": "2025-01-15",
>       "model": "claude-sonnet-4.5",
>       "usage": 203
>     },
>     {
>       "event_date": "2025-01-15",
>       "model": "gpt-4o",
>       "usage": 67
>     },
>     {
>       "event_date": "2025-01-16",
>       "model": "claude-sonnet-4.5",
>       "usage": 198
>     }
>   ],
>   "params": {
>     "metric": "ask-mode",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31"
>   }
> }
> ```
>
> ### Conversation Insights
>
> /analytics/team/conversation-insights
>
> Get the same aggregate Conversation Insights data you see in the dashboard. This endpoint returns aggregate insights, not raw conversation exports or raw conversation content.
>
> Available only for enterprise teams with Conversation Insights enabled. If **Disable Conversation Insights** is turned on in team settings, this endpoint returns `401`.
>
> For user-level filtering, use the shared `users` query parameter described in [Team-Level Endpoints](https://cursor.com/docs/account/teams/analytics-api.md#team-level-endpoints). SCIM group filtering is available in the dashboard UI only and isn't supported in the Analytics API.
>
> `intents` and `complexity` describe whole conversations.
>
> `categories`, `guidanceLevels`, and `workTypes` describe work across conversation segments.
>
> #### Parameters
>
> `startDate` string
>
> Start date for the analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for the analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `include` string | string\[]
>
> Required. Select which Conversation Insights slices to return. Supported values: `intents`, `complexity`, `categories`, `guidanceLevels`, and `workTypes`. You can pass `include` as a comma-separated list like `include=intents,complexity` or repeat it like `include=intents&include=workTypes`.
>
> `users` string
>
> Optional. Filter Conversation Insights to specific users. Pass comma-separated emails or user IDs, for example `users=alice@example.com,user_abc123`.
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/team/conversation-insights?startDate=2026-03-01&endDate=2026-03-07&include=intents,complexity,categories,guidanceLevels,workTypes&users=alice@example.com,bob@example.com" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "intents": {
>       "distribution": [
>         {
>           "intent": "Write Code",
>           "count": 18
>         },
>         {
>           "intent": "Ask",
>           "count": 7
>         },
>         {
>           "intent": "Plan",
>           "count": 3
>         }
>       ],
>       "topValues": [
>         {
>           "intent": "Write Code",
>           "count": 18
>         },
>         {
>           "intent": "Ask",
>           "count": 7
>         }
>       ],
>       "timeSeries": [
>         {
>           "date": "2026-03-01",
>           "intent": "Ask",
>           "count": 2
>         },
>         {
>           "date": "2026-03-02",
>           "intent": "Write Code",
>           "count": 6
>         }
>       ],
>       "subcategories": {
>         "askMode": [
>           {
>             "subcategory": "error_fix",
>             "count": 4
>           }
>         ],
>         "planMode": [
>           {
>             "subcategory": "implementation",
>             "count": 3
>           }
>         ],
>         "writeCode": [
>           {
>             "subcategory": "feature",
>             "count": 11
>           }
>         ]
>       }
>     },
>     "complexity": {
>       "distribution": [
>         {
>           "complexity": "high",
>           "count": 12
>         },
>         {
>           "complexity": "medium",
>           "count": 10
>         }
>       ],
>       "timeSeries": [
>         {
>           "date": "2026-03-01",
>           "complexity": "medium",
>           "count": 4
>         },
>         {
>           "date": "2026-03-02",
>           "complexity": "high",
>           "count": 5
>         }
>       ]
>     },
>     "categories": {
>       "distribution": [
>         {
>           "category": "New Features",
>           "count": 9
>         },
>         {
>           "category": "Bug Fixing & Debugging",
>           "count": 6
>         }
>       ],
>       "timeSeries": [
>         {
>           "date": "2026-03-01",
>           "category": "Bug Fixing & Debugging",
>           "count": 2
>         },
>         {
>           "date": "2026-03-02",
>           "category": "New Features",
>           "count": 4
>         }
>       ]
>     },
>     "guidanceLevels": {
>       "distribution": [
>         {
>           "guidanceLevel": "high",
>           "count": 8
>         },
>         {
>           "guidanceLevel": "medium",
>           "count": 7
>         }
>       ],
>       "timeSeries": [
>         {
>           "date": "2026-03-01",
>           "guidanceLevel": "medium",
>           "count": 3
>         },
>         {
>           "date": "2026-03-02",
>           "guidanceLevel": "high",
>           "count": 4
>         }
>       ]
>     },
>     "workTypes": {
>       "distribution": [
>         {
>           "workType": "new_feature",
>           "count": 9
>         },
>         {
>           "workType": "bug",
>           "count": 6
>         }
>       ],
>       "timeSeries": [
>         {
>           "date": "2026-03-01",
>           "workType": "bug",
>           "count": 2
>         },
>         {
>           "date": "2026-03-02",
>           "workType": "new_feature",
>           "count": 4
>         }
>       ]
>     }
>   },
>   "params": {
>     "metric": "conversation-insights",
>     "teamId": 12345,
>     "startDate": "2026-03-01",
>     "endDate": "2026-03-07",
>     "include": [
>       "intents",
>       "complexity",
>       "categories",
>       "guidanceLevels",
>       "workTypes"
>     ]
>   }
> }
> ```
>
> ### Leaderboard
>
> /analytics/team/leaderboard
>
> Get a leaderboard of team members ranked by AI usage metrics.
>
> **Behavior:**
>
> - **Without user filtering**: Returns users ranked by the specified metric (default: combined lines accepted)
> - **With user filtering**: Returns users that match the filter (with their actual team-wide rankings)
> - Supports pagination for teams with many members
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number for pagination (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 10, max: 500)
>
> `users` string
>
> Filter to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> Returns separate leaderboards for Tab autocomplete and Agent edits. When filtering by users, those users appear with their **actual team-wide rank**, not a filtered rank. For example, if you request a user who ranks #45 overall, they'll appear with `rank: 45`.
>
> ```bash
> # Get first page of leaderboard (top 10 users)
> curl -X GET "https://api.cursor.com/analytics/team/leaderboard" \
>   -u YOUR_API_KEY:
> ```
>
> ```bash
> # Get second page with custom page size
> curl -X GET "https://api.cursor.com/analytics/team/leaderboard?page=2&pageSize=20" \
>   -u YOUR_API_KEY:
> ```
>
> ```bash
> # Filter by specific users
> curl -X GET "https://api.cursor.com/analytics/team/leaderboard?users=alice@example.com,bob@example.com" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "tab_leaderboard": {
>       "data": [
>         {
>           "email": "alice@example.com",
>           "user_id": "user_abc123",
>           "profile_picture_url": "https://example.com/avatars/alice.jpg",
>           "total_accepts": 1334,
>           "total_lines_accepted": 3455,
>           "total_lines_suggested": 15307,
>           "line_acceptance_ratio": 0.2256519892590384,
>           "accept_ratio": 0.2330827067669173,
>           "rank": 1
>         },
>         {
>           "email": "bob@example.com",
>           "user_id": "user_def789",
>           "profile_picture_url": "https://example.com/avatars/bob.jpg",
>           "total_accepts": 796,
>           "total_lines_accepted": 2090,
>           "total_lines_suggested": 7689,
>           "line_acceptance_ratio": 0.2718168812589414,
>           "accept_ratio": 0.2731256599787746,
>           "rank": 2
>         }
>       ],
>       "total_users": 142
>     },
>     "agent_leaderboard": {
>       "data": [
>         {
>           "email": "alice@example.com",
>           "user_id": "user_abc123",
>           "profile_picture_url": "https://example.com/avatars/alice.jpg",
>           "total_accepts": 914,
>           "total_lines_accepted": 65947,
>           "total_lines_suggested": 201467,
>           "line_acceptance_ratio": 0.3273465219182842,
>           "rank": 1
>         },
>         {
>           "email": "bob@example.com",
>           "user_id": "user_def789",
>           "profile_picture_url": "https://example.com/avatars/bob.jpg",
>           "total_accepts": 843,
>           "total_lines_accepted": 61709,
>           "total_lines_suggested": 51092,
>           "line_acceptance_ratio": 1.2077924536684573,
>           "rank": 2
>         }
>       ],
>       "total_users": 142
>     }
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 10,
>     "totalUsers": 142,
>     "totalPages": 15,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "leaderboard",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 10
>   }
> }
> ```
>
> ### Bugbot Analytics
>
> /analytics/team/bugbot
>
> Get per-PR Bugbot review analytics for your team, including issue counts by severity and how many issues were resolved.
>
> For per-review data including billed cost and individual findings, use [Bugbot review analytics](https://cursor.com/docs/account/teams/analytics-api.md#bugbot-review-analytics).
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `prState` string
>
> PR state filter. Allowed values: `merged` or `all`. Default: `merged`. Use `merged` for merged PR analytics only. Use `all` for analytics across PR states.
>
> `repo` string
>
> Optional repository filter. Accepts full URLs or host/path formats (for example, `https://github.com/org/repo.git` or `github.com/org/repo`). Normalized to `host/owner/repo`.
>
> `page` number
>
> Page number for pagination (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of PRs per page (default: `100`, max: `250`)
>
> ```bash
> # Get Bugbot PR analytics for last 7 days (default window)
> curl -X GET "https://api.cursor.com/analytics/team/bugbot" \
>   -u YOUR_API_KEY:
> ```
>
> ```bash
> # Filter by repository and date range
> curl -X GET "https://api.cursor.com/analytics/team/bugbot?repo=github.com/acme/app&startDate=2025-01-01&endDate=2025-01-31" \
>   -u YOUR_API_KEY:
> ```
>
> ```bash
> # Paginate results
> curl -X GET "https://api.cursor.com/analytics/team/bugbot?page=2&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": [
>     {
>       "repo": "github.com/acme/app",
>       "pr_number": 42,
>       "timestamp": "2025-01-21T00:00:00.000Z",
>       "reviews": 3,
>       "issues": {
>         "total": 5,
>         "by_severity": {
>           "high": 1,
>           "medium": 2,
>           "low": 2
>         }
>       },
>       "issues_resolved": {
>         "total": 2,
>         "by_severity": {
>           "high": 1,
>           "medium": 1,
>           "low": 0
>         }
>       }
>     }
>   ],
>   "pagination": {
>     "page": 1,
>     "pageSize": 100,
>     "totalItems": 1,
>     "totalPages": 1,
>     "hasNextPage": false,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "bugbot",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "repo": "github.com/acme/app",
>     "prState": "merged",
>     "page": 1,
>     "pageSize": 100
>   }
> }
> ```
>
> ### Bugbot review analytics
>
> /analytics/team/bugbot-reviews
>
> Return one item per completed Bugbot review, including the reviewed commit, findings count, billed cost, and per-finding resolution data.
>
> Includes both posted reviews and dry-run reviews. Posted findings are identified by `comment_id` and `resolution_status`. Dry-run findings return `title`, `description`, and `locations` instead because nothing is posted to the SCM.
>
> Requires an API key with `read:*` scope.
>
> #### Parameters
>
> `startDate` string
>
> Start of the analytics range. Defaults to 7 days ago. See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats).
>
> `endDate` string
>
> End of the analytics range. Defaults to now. See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats).
>
> `repo` string
>
> Optional repository filter in `host/owner/repo` form. Protocol and `.git` suffix are optional.
>
> `prNumber` number
>
> Optional pull request or merge request number.
>
> `page` number
>
> Page number for pagination (1-indexed). Default: `1`.
>
> `pageSize` number
>
> Number of reviews per page. Default: `100`, max: `250`.
>
> `dryRun` boolean
>
> Optional filter for dry-run (`true`) or posted (`false`) reviews only.
>
> ```bash
> curl --get https://api.cursor.com/analytics/team/bugbot-reviews \
>   -u YOUR_API_KEY: \
>   --data-urlencode 'startDate=2026-06-01' \
>   --data-urlencode 'endDate=2026-06-29' \
>   --data-urlencode 'repo=github.com/your-org/your-repo' \
>   --data-urlencode 'prNumber=42' \
>   --data-urlencode 'page=1' \
>   --data-urlencode 'pageSize=100'
> ```
>
> ```bash
> curl --get https://api.cursor.com/analytics/team/bugbot-reviews \
>   -u YOUR_API_KEY: \
>   --data-urlencode 'dryRun=true' \
>   --data-urlencode 'repo=github.com/your-org/your-repo' \
>   --data-urlencode 'prNumber=42'
> ```
>
> **Response (posted review):**
>
> ```json
> {
>   "data": [
>     {
>       "request_id": "6e0d261c-86a2-4383-89f0-9162c1c10662",
>       "timestamp": "2026-06-29T19:42:18.000Z",
>       "repo": "github.com/your-org/your-repo",
>       "repo_node_id": "R_kgDOABCDEF",
>       "pr_number": 42,
>       "commit_sha": "9f3c2a1b7d8e4f5061728394a5b6c7d8e9f0a1b2",
>       "bugs_found": 2,
>       "cost_cents": 42.5,
>       "dry_run": false,
>       "publication_status": "posted",
>       "bugs": [
>         {
>           "comment_id": "2147483999",
>           "resolution_status": "resolved",
>           "severity": "high"
>         },
>         {
>           "comment_id": "2147484000",
>           "resolution_status": "unresolved",
>           "severity": "medium"
>         }
>       ]
>     }
>   ],
>   "pagination": {
>     "page": 1,
>     "pageSize": 100,
>     "totalItems": 1,
>     "totalPages": 1,
>     "hasNextPage": false,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "bugbot-reviews",
>     "teamId": 12345,
>     "startDate": "2026-06-01",
>     "endDate": "2026-06-29",
>     "repo": "github.com/your-org/your-repo",
>     "prNumber": 42,
>     "page": 1,
>     "pageSize": 100
>   }
> }
> ```
>
> **Response (dry-run review):**
>
> ```json
> {
>   "data": [
>     {
>       "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
>       "timestamp": "2026-06-29T20:15:03.000Z",
>       "repo": "github.com/your-org/your-repo",
>       "repo_node_id": "R_kgDOABCDEF",
>       "pr_number": 42,
>       "commit_sha": "9f3c2a1b7d8e4f5061728394a5b6c7d8e9f0a1b2",
>       "bugs_found": 1,
>       "cost_cents": null,
>       "dry_run": true,
>       "publication_status": "dry_run",
>       "bugs": [
>         {
>           "comment_id": null,
>           "resolution_status": null,
>           "severity": "medium",
>           "title": "Unbounded retry loop",
>           "description": "retry() recurses without a ceiling.",
>           "locations": [
>             { "file": "src/net.ts", "start_line": 5, "end_line": 9 }
>           ]
>         }
>       ]
>     }
>   ],
>   "pagination": {
>     "page": 1,
>     "pageSize": 100,
>     "totalItems": 1,
>     "totalPages": 1,
>     "hasNextPage": false,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "bugbot-reviews",
>     "teamId": 12345,
>     "startDate": "2026-06-01",
>     "endDate": "2026-06-29",
>     "repo": "github.com/your-org/your-repo",
>     "prNumber": 42,
>     "dryRun": true,
>     "page": 1,
>     "pageSize": 100
>   }
> }
> ```
>
> `repo_node_id`, `pr_number`, `commit_sha`, `cost_cents`, `bugs[].comment_id`, `bugs[].resolution_status`, and `bugs[].severity` may be `null` when unavailable. `cost_cents` is `null` when the review is not billed separately. For dry-run reviews, `bugs[].title`, `bugs[].description`, and `bugs[].locations` carry the finding content. Dry-run findings have `comment_id: null` and `resolution_status: null` because nothing is posted to the SCM.
>
> To trigger a dry-run review, call `POST /bugbot/review` with `"dryRun": true`. See the [Bugbot API docs](https://cursor.com/docs/bugbot.md#trigger-a-review).
>
> ***
>
> ## By-User Endpoints
>
> By-user endpoints provide the same metrics as team-level endpoints, but organized by individual users with pagination support. These are ideal for generating per-user reports or processing large teams in batches.
>
> ### Common Query Parameters
>
> | Parameter   | Type        | Required | Description                                                                                               |
> | ----------- | ----------- | -------- | --------------------------------------------------------------------------------------------------------- |
> | `startDate` | Date string | No       | Start date for the analytics period (default: 7 days ago)                                                 |
> | `endDate`   | Date string | No       | End date for the analytics period (default: today)                                                        |
> | `page`      | number      | No       | Page number (default: 1)                                                                                  |
> | `pageSize`  | number      | No       | Number of users per page (default: 100, max: 500)                                                         |
> | `users`     | string      | No       | Limit pagination to specific users (comma-separated emails or IDs, e.g., `alice@example.com,user_abc123`) |
>
> **User Filtering:**
> When you provide the `users` parameter to by-user endpoints:
>
> - **Pagination is filtered**: Only the specified users are included in the result set and pagination counts
> - **Useful for**: Getting detailed data for specific team members without paginating through all users
> - Example: If you have 500 users but only want data for 3 specific users, filter by their emails to get all 3 in a single page
>
> **Note:** By-user endpoints support the same date formats and shortcuts as team-level endpoints. See the [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats) section above.
>
> ### Response Format
>
> All by-user endpoints return data in this format:
>
> ```json
> {
>   "data": {
>     "user1@example.com": [ /* user's data */ ],
>     "user2@example.com": [ /* user's data */ ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 100,
>     "totalUsers": 250,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "agent-edits",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 100,
>     "userMappings": [
>       { "id": "user_abc123", "email": "user1@example.com" },
>       { "id": "user_def456", "email": "user2@example.com" }
>     ]
>   }
> }
> ```
>
> **Response Structure:**
>
> - `data` - Object keyed by user email addresses, each containing an array of that user's metrics
> - `pagination` - Pagination information
> - `params` - Request parameters echoed back
>   - `userMappings` - Array mapping email addresses to public user IDs for this page. Useful for cross-referencing with other APIs or creating links to user profiles.
>
> ### Available Endpoints
>
> All by-user endpoints follow the pattern: `/analytics/by-user/{metric}`
>
> - `GET /analytics/by-user/agent-edits` - Agent edits by user
> - `GET /analytics/by-user/tabs` - Tab usage by user
> - `GET /analytics/by-user/models` - Model usage by user
> - `GET /analytics/by-user/top-file-extensions` - Top files by user
> - `GET /analytics/by-user/client-versions` - Client versions by user
> - `GET /analytics/by-user/mcp` - MCP adoption by user
> - `GET /analytics/by-user/commands` - Commands adoption by user
> - `GET /analytics/by-user/plans` - Plans adoption by user
> - `GET /analytics/by-user/skills` - Skills adoption by user
> - `GET /analytics/by-user/ask-mode` - Ask mode adoption by user
>
> ### Agent Edits By User
>
> /analytics/by-user/agent-edits
>
> Get agent edits metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/agent-edits?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/agent-edits?users=alice@example.com,bob@example.com,carol@example.com" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "total_suggested_diffs": 145,
>         "total_accepted_diffs": 98,
>         "total_rejected_diffs": 47,
>         "total_green_lines_accepted": 820,
>         "total_red_lines_accepted": 160,
>         "total_green_lines_rejected": 210,
>         "total_red_lines_rejected": 60,
>         "total_green_lines_suggested": 1030,
>         "total_red_lines_suggested": 220,
>         "total_lines_suggested": 1250,
>         "total_lines_accepted": 980
>       },
>       {
>         "event_date": "2025-01-16",
>         "total_suggested_diffs": 132,
>         "total_accepted_diffs": 89,
>         "total_rejected_diffs": 43,
>         "total_green_lines_accepted": 740,
>         "total_red_lines_accepted": 150,
>         "total_green_lines_rejected": 185,
>         "total_red_lines_rejected": 55,
>         "total_green_lines_suggested": 925,
>         "total_red_lines_suggested": 175,
>         "total_lines_suggested": 1100,
>         "total_lines_accepted": 890
>       }
>     ],
>     "bob@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "total_suggested_diffs": 95,
>         "total_accepted_diffs": 72,
>         "total_rejected_diffs": 23,
>         "total_green_lines_accepted": 450,
>         "total_red_lines_accepted": 90,
>         "total_green_lines_rejected": 120,
>         "total_red_lines_rejected": 35,
>         "total_green_lines_suggested": 570,
>         "total_red_lines_suggested": 125,
>         "total_lines_suggested": 695,
>         "total_lines_accepted": 540
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "agent-edits",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ### Tab Usage By User
>
> /analytics/by-user/tabs
>
> Get Tab autocomplete metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/tabs?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "total_suggestions": 320,
>         "total_accepts": 210,
>         "total_rejects": 110,
>         "total_green_lines_accepted": 280,
>         "total_red_lines_accepted": 120,
>         "total_green_lines_rejected": 90,
>         "total_red_lines_rejected": 45,
>         "total_green_lines_suggested": 370,
>         "total_red_lines_suggested": 165,
>         "total_lines_suggested": 535,
>         "total_lines_accepted": 400
>       }
>     ],
>     "bob@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "total_suggestions": 180,
>         "total_accepts": 120,
>         "total_rejects": 60,
>         "total_green_lines_accepted": 150,
>         "total_red_lines_accepted": 70,
>         "total_green_lines_rejected": 50,
>         "total_red_lines_rejected": 25,
>         "total_green_lines_suggested": 200,
>         "total_red_lines_suggested": 95,
>         "total_lines_suggested": 295,
>         "total_lines_accepted": 220
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "tabs",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ### Model Usage By User
>
> /analytics/by-user/models
>
> Get model usage metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/models?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "date": "2025-01-15",
>         "model_breakdown": {
>           "claude-sonnet-4.5": {
>             "messages": 85,
>             "users": 1
>           },
>           "gpt-4o": {
>             "messages": 32,
>             "users": 1
>           }
>         }
>       }
>     ],
>     "bob@example.com": [
>       {
>         "date": "2025-01-15",
>         "model_breakdown": {
>           "claude-sonnet-4.5": {
>             "messages": 64,
>             "users": 1
>           }
>         }
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "models",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ### Top File Extensions By User
>
> /analytics/by-user/top-file-extensions
>
> Get top file extension metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/top-file-extensions?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "file_extension": "tsx",
>         "total_files": 45,
>         "total_accepts": 32,
>         "total_rejects": 10,
>         "total_lines_suggested": 890,
>         "total_lines_accepted": 650,
>         "total_lines_rejected": 240
>       },
>       {
>         "event_date": "2025-01-15",
>         "file_extension": "ts",
>         "total_files": 38,
>         "total_accepts": 28,
>         "total_rejects": 8,
>         "total_lines_suggested": 720,
>         "total_lines_accepted": 540,
>         "total_lines_rejected": 180
>       }
>     ],
>     "bob@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "file_extension": "py",
>         "total_files": 22,
>         "total_accepts": 18,
>         "total_rejects": 4,
>         "total_lines_suggested": 410,
>         "total_lines_accepted": 340,
>         "total_lines_rejected": 70
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "top-files",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ### Client Versions By User
>
> /analytics/by-user/client-versions
>
> Get client version metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/client-versions?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "client_version": "0.42.3",
>         "user_count": 1,
>         "percentage": 1.0
>       }
>     ],
>     "bob@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "client_version": "0.42.2",
>         "user_count": 1,
>         "percentage": 1.0
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "client-versions",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ### MCP Adoption By User
>
> /analytics/by-user/mcp
>
> Get MCP tool adoption metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/mcp?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "tool_name": "read_file",
>         "mcp_server_name": "filesystem",
>         "usage": 45
>       },
>       {
>         "event_date": "2025-01-16",
>         "tool_name": "read_file",
>         "mcp_server_name": "filesystem",
>         "usage": 38
>       }
>     ],
>     "bob@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "tool_name": "search_web",
>         "mcp_server_name": "brave-search",
>         "usage": 23
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "mcp",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ### Commands Adoption By User
>
> /analytics/by-user/commands
>
> Get command adoption metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/commands?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "command_name": "explain",
>         "usage": 12
>       },
>       {
>         "event_date": "2025-01-16",
>         "command_name": "explain",
>         "usage": 15
>       }
>     ],
>     "bob@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "command_name": "refactor",
>         "usage": 8
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "commands",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ### Plans Adoption By User
>
> /analytics/by-user/plans
>
> Get Plan mode adoption metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/plans?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "model": "claude-sonnet-4.5",
>         "usage": 23
>       },
>       {
>         "event_date": "2025-01-16",
>         "model": "claude-sonnet-4.5",
>         "usage": 19
>       }
>     ],
>     "bob@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "model": "gpt-4o",
>         "usage": 12
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "plans",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ### Skills Adoption By User
>
> /analytics/by-user/skills
>
> Get Skills adoption metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/skills?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "skill_name": "react-best-practices",
>         "usage": 8
>       },
>       {
>         "event_date": "2025-01-15",
>         "skill_name": "create-rule",
>         "usage": 3
>       },
>       {
>         "event_date": "2025-01-16",
>         "skill_name": "react-best-practices",
>         "usage": 5
>       }
>     ],
>     "bob@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "skill_name": "commit-message-helper",
>         "usage": 5
>       },
>       {
>         "event_date": "2025-01-15",
>         "skill_name": "create-skill",
>         "usage": 2
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "skills",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ### Ask Mode Adoption By User
>
> /analytics/by-user/ask-mode
>
> Get Ask mode adoption metrics organized by individual users with pagination support.
>
> #### Parameters
>
> `startDate` string
>
> Start date for analytics period (default: 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `endDate` string
>
> End date for analytics period (default: today). See [Date Formats](https://cursor.com/docs/account/teams/analytics-api.md#date-formats)
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of users per page (default: 100, max: 500)
>
> `users` string
>
> Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)
>
> ```bash
> curl -X GET "https://api.cursor.com/analytics/by-user/ask-mode?page=1&pageSize=50" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "data": {
>     "alice@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "model": "claude-sonnet-4.5",
>         "usage": 34
>       },
>       {
>         "event_date": "2025-01-16",
>         "model": "claude-sonnet-4.5",
>         "usage": 28
>       }
>     ],
>     "bob@example.com": [
>       {
>         "event_date": "2025-01-15",
>         "model": "gpt-4o",
>         "usage": 15
>       }
>     ]
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 50,
>     "totalUsers": 120,
>     "totalPages": 3,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "metric": "ask-mode",
>     "teamId": 12345,
>     "startDate": "2025-01-01",
>     "endDate": "2025-01-31",
>     "page": 1,
>     "pageSize": 50,
>     "userMappings": [
>       { "id": "user_abc123", "email": "alice@example.com" },
>       { "id": "user_def456", "email": "bob@example.com" }
>     ]
>   }
> }
> ```
>
> ***
>
> ## Team-Level Endpoints
>
> Team-level endpoints provide aggregated metrics for your entire team or filtered subsets of users. All endpoints support date range filtering and optional user filtering.
>
> ### Common Query Parameters
>
> | Parameter   | Type        | Required | Description                                                                                                                                                                |
> | ----------- | ----------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `startDate` | Date string | No       | Start date for the analytics period (default: 7 days ago)                                                                                                                  |
> | `endDate`   | Date string | No       | End date for the analytics period (default: today)                                                                                                                         |
> | `users`     | string      | No       | Filter data to specific users (comma-separated). Each value can be an email (e.g., `alice@example.com`) or public user ID (e.g., `user_abc123`). You can mix both formats. |
>
> **User Filtering:**
> The `users` parameter accepts a comma-separated list of identifiers. Each identifier can be:
>
> - **Email address** (e.g., `alice@example.com`) - Auto-detected by the presence of `@`
> - **Public user ID** (e.g., `user_abc123`) - Auto-detected by the `user_` prefix
> - **Mixed format** - You can combine emails and IDs in the same request
>
> **Examples:**
>
> ```bash
> # Filter by emails only
> ?users=alice@example.com,bob@example.com,carol@example.com
>
> # Filter by public user IDs only
> ?users=user_abc123,user_def456,user_ghi789
>
> # Mix emails and IDs
> ?users=alice@example.com,user_def456,bob@example.com
> ```
>
> When you filter by users, the API returns data **only for those specific users**. This is useful for:
>
> - Analyzing specific team members or groups (e.g., engineering leads, specific project teams)
> - Generating reports for a subset of users
> - Comparing metrics across selected individuals
>
> ### Date Formats
>
> **Default Behavior:**
> If you omit both `startDate` and `endDate`, the API defaults to the **last 7 days** (from 7 days ago until today). This is perfect for quick queries without specifying dates.
>
> **Standard Formats:**
>
> - `YYYY-MM-DD` - Simple date format (e.g., `2025-01-15`) **← Recommended**
> - ISO 8601 timestamps (e.g., `2025-01-15T00:00:00Z`)
>
> **Shortcuts:**
>
> - `now` or `today` - Current date (at 00:00:00)
> - `yesterday` - Yesterday's date (at 00:00:00)
> - `<number>d` - Days ago (e.g., `7d` = 7 days ago, `30d` = 30 days ago)
>
> **Important Notes:**
>
> - **Time is ignored**: All dates are resolved to the day level (00:00:00 UTC). Sending `2025-01-15T14:30:00Z` is the same as `2025-01-15`.
> - **Use recommended formats**: Use `YYYY-MM-DD` or shortcuts for better HTTP caching support. Different time values (like `T14:30:00Z` vs `T08:00:00Z`) prevent cache hits even though they resolve to the same day.
> - **Date ranges**: Limited to a maximum of 30 days.
>
> **Examples:**
>
> ```bash
> # Omit dates for last 7 days (simplest and best for caching)
> curl "https://api.cursor.com/analytics/team/agent-edits"
>
> # Using YYYY-MM-DD format for specific date range (recommended)
> ?startDate=2025-01-01&endDate=2025-01-31
>
> # Using shortcuts for last 30 days
> ?startDate=30d&endDate=today
>
> # Using shortcuts for last 14 days
> ?startDate=14d&endDate=now
>
> # ❌ Don't use timestamps - prevents caching and time is ignored anyway
> ?startDate=2025-01-15T14:30:00Z&endDate=2025-01-31T23:59:59Z
> ```
>
> ## Rate Limits
>
> Rate limits are enforced per team and reset every minute:
>
> - **Team-level endpoints**: 100 requests per minute per team
> - **By-user endpoints**: 50 requests per minute per team
>
> **What happens when you exceed the rate limit?**
>
> When you exceed the rate limit, you'll receive a `429 Too Many Requests` response:
>
> ```json
> {
>   "error": "Too Many Requests",
>   "message": "Rate limit exceeded. Please try again later."
> }
> ```
>
> ## Best Practices
>
> For general API best practices including exponential backoff, caching strategies, and error handling, see the [API Overview Best Practices](https://cursor.com/docs/api.md#best-practices).
>
> 1. **Use pagination for large teams**: If your team has more than 100 users, use the by-user endpoints with pagination to avoid timeouts.
> 2. **Leverage caching**: Both Team and User level endpoints support ETags. Store the ETag and use `If-None-Match` headers to reduce unnecessary data transfer.
> 3. **Filter by users when possible**: If you only need data for specific users, use the `users` parameter to reduce query time.
> 4. **Date ranges**: Keep date ranges reasonable (e.g., 1-3 months) for optimal performance.
>
>

### Source: Admin API

> # Admin API
>
> The Admin API lets you programmatically access your team's data, including member information, usage metrics, spending details, and model access.
>
> - The Admin API uses [Basic Authentication](https://cursor.com/docs/api.md#basic-authentication) with your API key as the username.
> - For details on creating API keys, authentication methods, rate limits, and best practices, see the [API Overview](https://cursor.com/docs/api.md).
>
> For org-wide actions across your teams, see [Organizations](https://cursor.com/docs/enterprise/organizations.md) and the [Organization API](https://cursor.com/docs/account/organizations/organization-admin-api.md).
>
> ## Endpoints
>
> ### Get Team Members
>
> /teams/members
>
> Retrieve all team members and their details.
>
> #### Response Fields
>
> `teamMembers` array
>
> Array of team member objects, each containing:
>
> - `id` string - Encoded user ID for the team member (e.g., `user_PDSPmvukpYgZEDXsoNirw3CFhy`)
> - `email` string - Email address of the team member
> - `name` string - Display name of the team member
> - `role` string - Role in the team (e.g., `member`, `owner`)
> - `isRemoved` boolean - Whether the member has been removed from the team
>
> ```bash
> curl -X GET https://api.cursor.com/teams/members \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "teamMembers": [
>     {
>       "id": "user_PDSPmvukpYgZEDXsoNirw3CFhy",
>       "name": "Alex",
>       "email": "developer@company.com",
>       "role": "member",
>       "isRemoved": false
>     },
>     {
>       "id": "user_kljUvI0ASZORvSEXf9hV0ydcso",
>       "name": "Sam",
>       "email": "admin@company.com",
>       "role": "owner",
>       "isRemoved": false
>     }
>   ]
> }
> ```
>
> ### Get Audit Logs
>
> /teams/audit-logs
>
> Retrieve audit log events for your team with filtering. Track team activity, security events, and configuration changes. Rate limited to 20 requests per minute per team. See [rate limits and best practices](https://cursor.com/docs/api.md#rate-limits).
>
> #### Parameters
>
> `startTime` string | number
>
> Start time (defaults to 7 days ago). See [Date Formats](https://cursor.com/docs/account/teams/admin-api.md#date-formats)
>
> `endTime` string | number
>
> End time (defaults to now). See [Date Formats](https://cursor.com/docs/account/teams/admin-api.md#date-formats)
>
> `eventTypes` string
>
> Comma-separated event types to filter by. Possible values: `login`, `logout`, `add_user`, `remove_user`, `update_user_role`, `team_settings`, `mcp_server_config`, `team_api_key`, `user_api_key`, `privacy_mode`, `user_spend_limit`, `team_rule`, `team_repo`, `team_hook`, `team_command`, `create_directory_group`, `delete_directory_group`, `update_directory_group`, `update_directory_group_permissions`, `add_user_to_directory_group`, `remove_user_from_directory_group`, `bugbot_installation`, `bugbot_installation_settings`, `bugbot_repo_settings`, `bugbot_team_rule`, `bugbot_team_settings`, `bugbot_bulk_repo_update`
>
> `search` string
>
> Search term to filter events
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Results per page (1-500). Default: `100`
>
> `users` string
>
> Filter by users. See [User Filtering](https://cursor.com/docs/account/teams/admin-api.md#user-filtering) below
>
> Date range cannot exceed 30 days. Make multiple requests for longer periods.
>
> #### Date Formats
>
> The `startTime` and `endTime` parameters support multiple formats:
>
> - **Relative shortcuts**: `now`, `today`, `yesterday`, `7d` (7 days ago), `5h` (5 hours ago), `300s` (300 seconds ago)
> - **ISO 8601 strings**: `2024-01-15T12:00:00Z` or `2024-01-15T10:00:00-05:00`
> - **YYYY-MM-DD format**: `2024-01-15` (time defaults to 00:00:00 UTC)
> - **Unix timestamps**: `1705315200` (seconds) or `1705315200000` (milliseconds)
>
> **Examples:**
>
> - `?startTime=7d&endTime=now` - Last 7 days
> - `?startTime=5h&endTime=now` - Last 5 hours
> - `?startTime=2024-01-15&endTime=2024-01-20` - Specific date range
> - `?startTime=1705315200000&endTime=1705401600000` - Unix timestamps
>
> #### User Filtering
>
> The `users` parameter accepts multiple formats, comma-separated:
>
> - **Email addresses**: `developer@company.com,admin@company.com`
> - **Encoded user IDs**: `user_PDSPmvukpYgZEDXsoNirw3CFhy,user_kljUvI0ASZORvSEXf9hV0ydcso`
>
> You can mix formats: `developer@company.com,12345,user_PDSPmvukpYgZEDXsoNirw3CFhy`
>
> Maximum number of users per request equals `pageSize`.
>
> ```bash
> curl -X GET "https://api.cursor.com/teams/audit-logs?users=admin@company.com,developer@company.com&eventTypes=login,add_user" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "events": [
>     {
>       "event_id": "evt_abc123",
>       "timestamp": "2024-01-15T12:30:00.000Z",
>       "ip_address": "203.0.113.42",
>       "user_email": "admin@company.com",
>       "event_type": "add_user",
>       "event_data": {
>         "email": "admin@company.com",
>         "method": "manual"
>       }
>     },
>     {
>       "event_id": "evt_def456",
>       "timestamp": "2024-01-15T10:15:00.000Z",
>       "ip_address": "192.168.1.1",
>       "user_email": "developer@company.com",
>       "event_type": "login",
>       "event_data": {
>         "ip_address": "192.168.1.1",
>         "user_agent": "Cursor/0.42.0"
>       }
>     }
>   ],
>   "pagination": {
>     "page": 1,
>     "pageSize": 100,
>     "totalCount": 2,
>     "totalPages": 1,
>     "hasNextPage": false,
>     "hasPreviousPage": false
>   },
>   "params": {
>     "teamId": 12345,
>     "startDate": 1704729600000,
>     "endDate": 1705334400000
>   }
> }
> ```
>
> ### Get Daily Usage Data
>
> /teams/daily-usage-data
>
> Retrieve daily usage metrics for your team. Data is aggregated at the hourly level - we recommend polling this endpoint at most once per hour. Rate limited to 20 requests per minute per team. See [best practices](https://cursor.com/docs/api.md#best-practices).
>
> #### Parameters
>
> `startDate` number Required
>
> Start date in epoch milliseconds
>
> `endDate` number Required
>
> End date in epoch milliseconds
>
> `page` number
>
> Page number (1-indexed). When provided along with `pageSize`, enables pagination and returns data for **all team members with a membership during the requested date range**.
>
> `pageSize` number
>
> Number of users per page. When provided along with `page`, enables pagination and returns data for **all team members with a membership during the requested date range**.
>
> Without pagination parameters, this endpoint only returns **active users** (those with activity during the date range). To get **all team members**, include both `page` and `pageSize` parameters.
>
> When using pagination, the response includes an `isActive` field for each user indicating whether they had activity on that day. Members who joined after the requested period are excluded.
>
> Date range cannot exceed 30 days. Make multiple requests for longer periods.
>
> The fields `subscriptionIncludedReqs`, `usageBasedReqs`, and `apiKeyReqs` count raw usage events, not billable request units in older request-based pricing. To get accurate billable request counts, use the [`/teams/filtered-usage-events`](https://cursor.com/docs/account/teams/admin-api.md#get-usage-events-data) endpoint and sum the `requestsCosts` field.
>
> #### Response Fields
>
> Each object in the `data` array contains:
>
> - `userId` number - Unique identifier for the user
> - `day` string - The date this record covers (ISO date, e.g., `2024-03-18`)
> - `date` number - Date as epoch milliseconds
> - `email` string - User's email address
> - `isActive` boolean - Whether the user had activity on this day (only present with pagination)
> - `totalLinesAdded` number - Total lines of code added
> - `totalLinesDeleted` number - Total lines of code deleted
> - `acceptedLinesAdded` number - AI-suggested lines added that were accepted
> - `acceptedLinesDeleted` number - AI-suggested lines deleted that were accepted
> - `totalApplies` number - Total AI code apply actions
> - `totalAccepts` number - Total accepted AI suggestions
> - `totalRejects` number - Total rejected AI suggestions
> - `totalTabsShown` number - Total Tab completions shown to the user
> - `totalTabsAccepted` number - Total Tab completions accepted by the user
> - `composerRequests` number - Number of Composer requests made
> - `chatRequests` number - Number of chat requests made
> - `agentRequests` number - Number of Agent mode requests made
> - `cmdkUsages` number - Number of Cmd+K inline edit usages
> - `subscriptionIncludedReqs` number - Requests included in the subscription plan
> - `apiKeyReqs` number - Requests made via API key
> - `usageBasedReqs` number - Usage-based (overage) requests
> - `bugbotUsages` number - Number of Bugbot usages
> - `mostUsedModel` string | null - Most frequently used AI model for the day
> - `applyMostUsedExtension` string | null - Most common file extension for apply actions
> - `tabMostUsedExtension` string | null - Most common file extension for Tab completions
> - `clientVersion` string | null - Cursor client version used
>
> ```bash
> # Get data for active users only (no pagination)
> curl -X POST https://api.cursor.com/teams/daily-usage-data \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "startDate": 1710720000000,
>     "endDate": 1710892800000
>   }'
>
> # Get data for ALL team members (with pagination)
> curl -X POST https://api.cursor.com/teams/daily-usage-data \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "startDate": 1710720000000,
>     "endDate": 1710892800000,
>     "page": 1,
>     "pageSize": 1000
>   }'
> ```
>
> **Response (without pagination - active users only):**
>
> ```json
> {
>   "data": [
>     {
>       "userId": 12345,
>       "day": "2024-03-18",
>       "date": 1710720000000,
>       "isActive": true,
>       "totalLinesAdded": 1543,
>       "totalLinesDeleted": 892,
>       "acceptedLinesAdded": 1102,
>       "acceptedLinesDeleted": 645,
>       "totalApplies": 87,
>       "totalAccepts": 73,
>       "totalRejects": 14,
>       "totalTabsShown": 342,
>       "totalTabsAccepted": 289,
>       "composerRequests": 45,
>       "chatRequests": 128,
>       "agentRequests": 12,
>       "cmdkUsages": 67,
>       "subscriptionIncludedReqs": 180,
>       "apiKeyReqs": 0,
>       "usageBasedReqs": 5,
>       "bugbotUsages": 3,
>       "mostUsedModel": "gpt-5",
>       "applyMostUsedExtension": ".tsx",
>       "tabMostUsedExtension": ".ts",
>       "clientVersion": "0.25.1",
>       "email": "developer@company.com"
>     }
>   ],
>   "period": {
>     "startDate": 1710720000000,
>     "endDate": 1710892800000
>   }
> }
> ```
>
> **Response (with pagination - all team members):**
>
> ```json
> {
>   "data": [
>     {
>       "userId": 12345,
>       "day": "2024-03-18",
>       "date": 1710720000000,
>       "isActive": true,
>       "totalLinesAdded": 1543,
>       "totalLinesDeleted": 892,
>       "acceptedLinesAdded": 1102,
>       "acceptedLinesDeleted": 645,
>       "totalApplies": 87,
>       "totalAccepts": 73,
>       "totalRejects": 14,
>       "totalTabsShown": 342,
>       "totalTabsAccepted": 289,
>       "composerRequests": 45,
>       "chatRequests": 128,
>       "agentRequests": 12,
>       "cmdkUsages": 67,
>       "subscriptionIncludedReqs": 180,
>       "apiKeyReqs": 0,
>       "usageBasedReqs": 5,
>       "bugbotUsages": 3,
>       "mostUsedModel": "gpt-5",
>       "applyMostUsedExtension": ".tsx",
>       "tabMostUsedExtension": ".ts",
>       "clientVersion": "0.25.1",
>       "email": "developer@company.com"
>     },
>     {
>       "userId": 12346,
>       "day": "2024-03-18",
>       "date": 1710720000000,
>       "isActive": false,
>       "totalLinesAdded": 0,
>       "totalLinesDeleted": 0,
>       "acceptedLinesAdded": 0,
>       "acceptedLinesDeleted": 0,
>       "totalApplies": 0,
>       "totalAccepts": 0,
>       "totalRejects": 0,
>       "totalTabsShown": 0,
>       "totalTabsAccepted": 0,
>       "composerRequests": 0,
>       "chatRequests": 0,
>       "agentRequests": 0,
>       "cmdkUsages": 0,
>       "subscriptionIncludedReqs": 0,
>       "apiKeyReqs": 0,
>       "usageBasedReqs": 0,
>       "bugbotUsages": 0,
>       "mostUsedModel": null,
>       "applyMostUsedExtension": null,
>       "tabMostUsedExtension": null,
>       "clientVersion": null,
>       "email": "inactive-user@company.com"
>     }
>   ],
>   "period": {
>     "startDate": 1710720000000,
>     "endDate": 1710892800000
>   },
>   "pagination": {
>     "page": 1,
>     "pageSize": 1000,
>     "totalUsers": 150,
>     "totalPages": 1,
>     "hasNextPage": false,
>     "hasPreviousPage": false
>   }
> }
> ```
>
> ### Get Spending Data
>
> /teams/spend
>
> Retrieve spending information for the current billing cycle with search, sorting, and pagination.
>
> #### Parameters
>
> `searchTerm` string
>
> Search in user names and emails
>
> `sortBy` string
>
> Sort by: `amount`, `date`, `user`. Default: `date`
>
> `sortDirection` string
>
> Sort direction: `asc`, `desc`. Default: `desc`
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Results per page
>
> #### Response Fields
>
> Each object in `teamMemberSpend` contains:
>
> - `userId` string - Encoded user ID (e.g., `user_PDSPmvukpYgZEDXsoNirw3CFhy`). Shares the same identifier namespace as `teamMembers[].id` from [`/teams/members`](https://cursor.com/docs/account/teams/admin-api.md#get-team-members).
> - `name` string - Display name of the user
> - `email` string - Email address of the user
> - `role` string - Role in the team (e.g., `member`, `owner`)
> - `spendCents` number - On-demand spend in cents for the current billing cycle (excludes included usage)
> - `overallSpendCents` number - Total spend in cents for the current billing cycle, including both on-demand and included usage
> - `fastPremiumRequests` number - Number of usage-based premium requests made during the billing cycle
> - `hardLimitOverrideDollars` number - Custom hard spending limit override in dollars for this user (0 means no override)
> - `monthlyLimitDollars` number | null - Monthly spending limit in dollars set for this user, or `null` if no limit is set
> - `effectivePerUserLimitDollars` number - Currently enforced per-user spending limit in dollars, derived from `monthlyLimitDollars` and `hardLimitOverrideDollars`
>
> On June 4th, 2026 we added additional precision to the spendCents and overallSpendCents fields to avoid rounding errors when comparing results to invoice amounts.
>
> ```bash
> curl -X POST https://api.cursor.com/teams/spend \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "searchTerm": "alex@company.com",
>     "page": 2,
>     "pageSize": 25
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "teamMemberSpend": [
>     {
>       "userId": "user_PDSPmvukpYgZEDXsoNirw3CFhy",
>       "spendCents": 2450.125487,
>       "overallSpendCents": 2450.125487,
>       "fastPremiumRequests": 1250,
>       "name": "Alex",
>       "email": "developer@company.com",
>       "role": "member",
>       "hardLimitOverrideDollars": 100,
>       "monthlyLimitDollars": 200,
>       "effectivePerUserLimitDollars": 100
>     },
>     {
>       "userId": "user_kljUvI0ASZORvSEXf9hV0ydcso",
>       "spendCents": 1875.500123,
>       "overallSpendCents": 3200.750456,
>       "fastPremiumRequests": 980,
>       "name": "Sam",
>       "email": "admin@company.com",
>       "role": "owner",
>       "hardLimitOverrideDollars": 0,
>       "monthlyLimitDollars": null,
>       "effectivePerUserLimitDollars": 50
>     }
>   ],
>   "subscriptionCycleStart": 1708992000000,
>   "totalMembers": 15,
>   "totalPages": 1
> }
> ```
>
> ### Get Usage Events Data
>
> /teams/filtered-usage-events
>
> Retrieve detailed usage events for your team with filtering, search, and pagination options. This endpoint provides granular insights into API calls, model usage, token consumption, and costs. Data is aggregated at the hourly level. We recommend polling this endpoint at most once per hour. Rate limited to 60 requests per minute per team. See the [API guidance](https://cursor.com/docs/api.md#best-practices).
>
> **Cost Calculation**: To reconcile event-level costs with `/teams/spend` totals, sum the `chargedCents` field across events. This field includes both the model cost and the Cursor Token Rate when a request is eligible for the rate, matching the dashboard totals. It works for both token-based and request-based billing plans.
>
> The `cursorTokenFee` field represents the Cursor Token Rate and is only present when the rate applies to a third-party model request. This includes when Auto routes to a third-party model. First-party Cursor models such as Grok and Composer, and request-based enterprise accounts do not include this fee. See [Cursor Token Rate](https://cursor.com/help/models-and-usage/token-rate.md).
>
> #### Parameters
>
> `startDate` number
>
> Start date in epoch milliseconds. This bound is inclusive.
>
> `endDate` number
>
> End date in epoch milliseconds. This bound is inclusive.
>
> `startDate` and `endDate` are points in time with millisecond precision, and
> both bounds are inclusive. An event exactly at `2026-05-08T00:00:00.000Z` is
> included when `endDate` is `1778198400000`. For non-overlapping daily
> ingestion windows, set the previous window's `endDate` to the final
> millisecond of the day, such as `2026-05-07T23:59:59.999Z`.
>
> `userId` number
>
> Filter by specific user ID
>
> `page` number
>
> Page number (1-indexed). Default: `1`
>
> `pageSize` number
>
> Number of results per page. Default: `100`. Maximum: `1000`.
>
> `email` string
>
> Filter by user email address
>
> `serviceAccountId` string
>
> Filter by service account ID
>
> `cloudAgentId` string
>
> Filter by a specific cloud agent run ID. Pass `*` to return events from all cloud agent runs.
>
> `automationId` string
>
> Filter by a specific automation UUID. Pass `*` to return events from all automations.
>
> `hostingType` string
>
> Filter cloud agent (background agent) runs by where they executed. Use this to isolate inference spend for self-hosted agents from Cursor-hosted runs. Accepted values:
>
> - `CLOUD` - Cursor-hosted runs
> - `SELF_HOSTED` - any self-hosted run (a self-hosted pool worker or a personal "My Machine" worker)
> - `SELF_HOSTED_POOL` - team self-hosted pool workers only
> - `SELF_HOSTED_MACHINE` - personal "My Machine" workers only
>
> An unrecognized `hostingType` value returns a `400` error rather than an empty result, so a typo can't be mistaken for genuinely zero self-hosted spend. This filter covers inference spend only; self-hosted compute runs on your own machines and is never metered by Cursor.
>
> When you pass multiple filters, the endpoint combines them with `AND`. For example, `automationId` and `serviceAccountId` return events that match both values.
>
> #### Response Fields
>
> Each object in `usageEvents` contains:
>
> - `timestamp` string - Event timestamp in epoch milliseconds (as a string)
> - `userEmail` string - Email address of the user who made the request
> - `serviceAccountId` string | undefined - ID of the service account that made the request. Omitted for human user events.
> - `serviceAccountName` string | undefined - Display name of the service account that made the request. Omitted for human user events.
> - `cloudAgentId` string | undefined - ID of the cloud agent run attributed to this event. Omitted for events outside cloud agents.
> - `automationId` string | undefined - UUID of the automation attributed to this event. Omitted for events outside automations.
> - `conversationId` string | undefined - ID of the conversation (agent session) that generated this event. Use it to attribute spend to a session or as a join key with other sources that expose conversation IDs, such as the [AI Code Tracking API](https://cursor.com/docs/account/teams/ai-code-tracking-api.md). Omitted for events without an associated conversation.
> - `model` string - AI model used for the request
> - `kind` string - Billing category (e.g., `Usage-based`, `Included in Business`)
> - `maxMode` boolean - Whether the request used max mode
> - `requestsCosts` number - Cost in request units
> - `isTokenBasedCall` boolean - Whether the request was billed by token usage
> - `isChargeable` boolean - Whether this event incurs a charge
> - `isHeadless` boolean - Whether this request was made without a connected client (e.g., background agents)
> - `tokenUsage` object | undefined - Token usage details (present when `isTokenBasedCall` is `true`):
>   - `inputTokens` number - Input tokens consumed
>   - `outputTokens` number - Output tokens generated
>   - `cacheWriteTokens` number - Tokens written to cache
>   - `cacheReadTokens` number - Tokens read from cache
>   - `totalCents` number - Total model cost in cents
>   - `discountPercentOff` number | undefined - Discount percentage applied, if any
> - `chargedCents` number - Total amount charged in cents for this event. For third-party model requests subject to the Cursor Token Rate, this includes model cost plus the Cursor Token Rate. Use this field to reconcile event-level costs with `/teams/spend` totals. Works for both token-based and request-based billing plans.
> - `cursorTokenFee` number | undefined - Cursor Token Rate in cents. Present only when the rate applies to a third-party model request (including when Auto routes to a third-party model).
>
> ```bash
> curl -X POST https://api.cursor.com/teams/filtered-usage-events \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "startDate": 1748411762359,
>     "endDate": 1751003762359,
>     "email": "developer@company.com",
>     "page": 1,
>     "pageSize": 25
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "totalUsageEventsCount": 113,
>   "pagination": {
>     "numPages": 5,
>     "currentPage": 1,
>     "pageSize": 25,
>     "hasNextPage": true,
>     "hasPreviousPage": false
>   },
>   "usageEvents": [
>     {
>       "timestamp": "1750979225854",
>       "userEmail": "developer@company.com",
>       "conversationId": "8f2e4a1b-6c3d-4e5f-9a7b-2d1c8e6f4a3b",
>       "model": "claude-4.5-sonnet",
>       "kind": "Usage-based",
>       "maxMode": true,
>       "requestsCosts": 5,
>       "isTokenBasedCall": true,
>       "isChargeable": true,
>       "isHeadless": false,
>       "tokenUsage": {
>         "inputTokens": 126,
>         "outputTokens": 450,
>         "cacheWriteTokens": 6112,
>         "cacheReadTokens": 11964,
>         "totalCents": 20.18232
>       },
>       "chargedCents": 21.36232,
>       "cursorTokenFee": 1.18
>     },
>     {
>       "timestamp": "1750979173824",
>       "userEmail": "developer@company.com",
>       "conversationId": "8f2e4a1b-6c3d-4e5f-9a7b-2d1c8e6f4a3b",
>       "model": "claude-4.5-sonnet",
>       "kind": "Usage-based",
>       "maxMode": true,
>       "requestsCosts": 10,
>       "isTokenBasedCall": true,
>       "isChargeable": true,
>       "isHeadless": false,
>       "tokenUsage": {
>         "inputTokens": 5805,
>         "outputTokens": 311,
>         "cacheWriteTokens": 11964,
>         "cacheReadTokens": 0,
>         "totalCents": 40.167,
>         "discountPercentOff": 10
>       },
>       "chargedCents": 37.33,
>       "cursorTokenFee": 1.18
>     },
>     {
>       "timestamp": "1750978339901",
>       "userEmail": "admin@company.com",
>       "model": "claude-4-sonnet-thinking",
>       "kind": "Included in Business",
>       "maxMode": true,
>       "requestsCosts": 1.4,
>       "isTokenBasedCall": false,
>       "isChargeable": false,
>       "isHeadless": false,
>       "chargedCents": 8
>     }
>   ],
>   "period": {
>     "startDate": 1748411762359,
>     "endDate": 1751003762359
>   }
> }
> ```
>
> **Service account usage example:**
>
> ```bash
> curl -X POST https://api.cursor.com/teams/filtered-usage-events \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "startDate": 1748411762359,
>     "endDate": 1751003762359,
>     "serviceAccountId": "sa_abc123",
>     "page": 1,
>     "pageSize": 10
>   }'
> ```
>
> **Service account response:**
>
> ```json
> {
>   "totalUsageEventsCount": 1,
>   "pagination": {
>     "numPages": 1,
>     "currentPage": 1,
>     "pageSize": 10,
>     "hasNextPage": false,
>     "hasPreviousPage": false
>   },
>   "usageEvents": [
>     {
>       "timestamp": "1750979225854",
>       "userEmail": "agent-runner@company.com",
>       "serviceAccountId": "sa_abc123",
>       "serviceAccountName": "Nightly CI Agent",
>       "conversationId": "3b9d7c2e-1f4a-4b8c-a6d5-e9f0a2b4c6d8",
>       "model": "claude-4.5-sonnet",
>       "kind": "Usage-based",
>       "maxMode": true,
>       "requestsCosts": 5,
>       "isTokenBasedCall": true,
>       "isChargeable": true,
>       "isHeadless": true,
>       "tokenUsage": {
>         "inputTokens": 126,
>         "outputTokens": 450,
>         "cacheWriteTokens": 6112,
>         "cacheReadTokens": 11964,
>         "totalCents": 20.18232
>       },
>       "chargedCents": 21.36232,
>       "cursorTokenFee": 1.18
>     }
>   ],
>   "period": {
>     "startDate": 1748411762359,
>     "endDate": 1751003762359
>   }
> }
> ```
>
> **Automation usage example:**
>
> Use an automation UUID to retrieve its usage events. Automation attribution works for automations that run as a user or a service account.
>
> ```bash
> curl -X POST https://api.cursor.com/teams/filtered-usage-events \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "startDate": 1748411762359,
>     "endDate": 1751003762359,
>     "automationId": "7fc64f90-6d7a-4a5d-91b1-bd1f529a85dd",
>     "page": 1,
>     "pageSize": 100
>   }'
> ```
>
> Each matching event includes its `automationId` and `cloudAgentId`. Sum `chargedCents` across the events to calculate the automation's total cost.
>
> **Self-hosted agent spend example:**
>
> ```bash
> curl -X POST https://api.cursor.com/teams/filtered-usage-events \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "startDate": 1748411762359,
>     "endDate": 1751003762359,
>     "hostingType": "SELF_HOSTED",
>     "page": 1,
>     "pageSize": 10
>   }'
> ```
>
> ### Set User Spend Limit
>
> /teams/user-spend-limit
>
> Set spending limits for individual team members. This allows you to control how much each user can spend on AI usage within your team. Rate limited to 250 requests per minute per team. See [rate limits](https://cursor.com/docs/api.md#rate-limits).
>
> #### Parameters
>
> `userEmail` string Required
>
> Email address of the team member
>
> `spendLimitDollars` number | null Required
>
> Spending limit in dollars (integer only, no decimals). Set to `null` to remove the limit.
>
> - **Availability**: Enterprise only
> - The user must already be a member of your team
> - Only integer values are accepted (no decimal amounts)
> - Setting `spendLimitDollars` to 0 will set the limit to $0
> - Setting `spendLimitDollars` to `null` will clear/remove the limit entirely
>
> ```bash
> curl -X POST https://api.cursor.com/teams/user-spend-limit \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "userEmail": "developer@company.com",
>     "spendLimitDollars": 100
>   }'
> ```
>
> **Successful response:**
>
> ```json
> {
>   "outcome": "success",
>   "message": "Spend limit set to $100 for user developer@company.com"
> }
> ```
>
> **Error response:**
>
> ```json
> {
>   "outcome": "error",
>   "message": "Invalid email format"
> }
> ```
>
> ### Remove Team Member
>
> /teams/remove-member
>
> Remove a member from your team programmatically. This is useful for automating offboarding workflows or integrating with HR systems. Rate limited to 50 requests per minute per team. See [rate limits](https://cursor.com/docs/api.md#rate-limits).
>
> #### Parameters
>
> `userId` string
>
> Encoded user ID (e.g., `user_PDSPmvukpYgZEDXsoNirw3CFhy`). Required if `email` is not provided.
>
> `email` string
>
> Email address of the team member. Required if `userId` is not provided.
>
> - **Availability**: Enterprise only
> - Provide either `userId` or `email`, but not both
> - At least one paid member must remain on the team after removal
> - At least one admin (owner or free-owner) must remain on the team after removal
>
> ```bash
> curl -X POST https://api.cursor.com/teams/remove-member \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "email": "developer@company.com"
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "success": true,
>   "userId": "user_PDSPmvukpYgZEDXsoNirw3CFhy",
>   "hasBillingCycleUsage": true
> }
> ```
>
> **Remove by user ID:**
>
> ```bash
> curl -X POST https://api.cursor.com/teams/remove-member \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "userId": "user_PDSPmvukpYgZEDXsoNirw3CFhy"
>   }'
> ```
>
> **Error responses:**
>
> ```json
> {
>   "error": "User is not a member of this team"
> }
> ```
>
> ```json
> {
>   "error": "Either userId or email must be provided"
> }
> ```
>
> ```json
> {
>   "error": "Only one of userId or email should be provided, not both"
> }
> ```
>
> ### Get Team Repo Blocklists
>
> /settings/repo-blocklists/repos
>
> Retrieve all repository blocklists configured for your team. Add repositories and use patterns to prevent files or directories from being indexed or used as context.
>
> #### Pattern Examples
>
> Common blocklist patterns:
>
> - `*` - Block entire repository
> - `*.env` - Block all .env files
> - `config/*` - Block all files in config directory
> - `**/*.secret` - Block all .secret files in any subdirectory
> - `src/api/keys.ts` - Block specific file
>
> ```bash
> curl -X GET https://api.cursor.com/settings/repo-blocklists/repos \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "repos": [
>     {
>       "id": "repo_123",
>       "url": "https://github.com/company/sensitive-repo",
>       "patterns": ["*.env", "config/*", "secrets/**"]
>     },
>     {
>       "id": "repo_456",
>       "url": "https://github.com/company/internal-tools",
>       "patterns": ["*"]
>     }
>   ]
> }
> ```
>
> ### Upsert Repo Blocklists
>
> /settings/repo-blocklists/repos/upsert
>
> Replace existing repository blocklists for the provided repos. This endpoint will only overwrite the patterns for the repositories provided. All other repos will be unaffected.
>
> #### Parameters
>
> `repos` array Required
>
> Array of repository blocklist objects. Each repository object must contain:
>
> - `url` string - Repository URL to blocklist
> - `patterns` string\[] - Array of file patterns to block (glob patterns supported)
>
> ```bash
> curl -X POST https://api.cursor.com/settings/repo-blocklists/repos/upsert \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "repos": [
>       {
>         "url": "https://github.com/company/sensitive-repo",
>         "patterns": ["*.env", "config/*", "secrets/**"]
>       },
>       {
>         "url": "https://github.com/company/internal-tools",
>         "patterns": ["*"]
>       }
>     ]
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "repos": [
>     {
>       "id": "repo_123",
>       "url": "https://github.com/company/sensitive-repo",
>       "patterns": ["*.env", "config/*", "secrets/**"]
>     },
>     {
>       "id": "repo_456",
>       "url": "https://github.com/company/internal-tools",
>       "patterns": ["*"]
>     }
>   ]
> }
> ```
>
> ### Delete Repo Blocklist
>
> /settings/repo-blocklists/repos/:repoId
>
> Remove a specific repository from the blocklist. Returns 204 No Content on successful deletion.
>
> #### Parameters
>
> `repoId` string Required
>
> ID of the repository blocklist to delete
>
> ```bash
> curl -X DELETE https://api.cursor.com/settings/repo-blocklists/repos/repo_123 \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```text
> 204 No Content
> ```
>
> ## Billing Groups
>
> [Billing groups](https://cursor.com/docs/account/enterprise/billing-groups.md) allow Enterprise admins to understand and manage spend across groups of users. This functionality is useful for reporting, internal chargebacks, and budgeting.
>
> Members can only be in one billing group at a time. Members not assigned to any group are placed in a reserved `Unassigned` group.
>
> ### List Groups
>
> /teams/groups
>
> Retrieve all billing groups for your team with spend data for the current billing cycle.
>
> #### Parameters
>
> `billingCycle` string
>
> ISO date string (e.g., `2025-01-15`) to specify which billing cycle to query. Defaults to current cycle.
>
> ```bash
> curl -X GET "https://api.cursor.com/teams/groups?billingCycle=2025-01-15" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "groups": [
>     {
>       "id": "group_PDSPmvukpYgZEDXsoNirw3CFhy",
>       "name": "Engineering",
>       "type": "BILLING",
>       "directoryGroupId": null,
>       "memberCount": 12,
>       "createdAt": "2024-01-15T10:30:00.000Z",
>       "updatedAt": "2024-01-20T14:22:00.000Z",
>       "spendCents": 245000,
>       "currentMembers": [
>         {
>           "userId": "user_abc123",
>           "name": "Alex Developer",
>           "email": "alex@company.com",
>           "joinedAt": "2024-01-15T10:30:00.000Z",
>           "leftAt": null,
>           "spendCents": 12500
>         }
>       ],
>       "formerMembers": [],
>       "dailySpend": [
>         { "date": "2025-01-15", "spendCents": 8500 },
>         { "date": "2025-01-16", "spendCents": 9200 }
>       ]
>     },
>     {
>       "id": "group_kljUvI0ASZORvSEXf9hV0ydcso",
>       "name": "Design",
>       "type": "BILLING",
>       "directoryGroupId": "dir_group_abc123xyz",
>       "memberCount": 5,
>       "createdAt": "2024-01-16T09:00:00.000Z",
>       "updatedAt": "2024-01-16T09:00:00.000Z",
>       "spendCents": 87500,
>       "currentMembers": [],
>       "formerMembers": [],
>       "dailySpend": []
>     }
>   ],
>   "unassignedGroup": {
>     "id": "group_unassigned",
>     "name": "Unassigned",
>     "type": "BILLING",
>     "directoryGroupId": null,
>     "memberCount": 3,
>     "createdAt": "2024-01-01T00:00:00.000Z",
>     "updatedAt": "2024-01-01T00:00:00.000Z",
>     "spendCents": 15000,
>     "currentMembers": [],
>     "formerMembers": [],
>     "dailySpend": []
>   },
>   "billingCycle": {
>     "cycleStart": "2025-01-01T00:00:00.000Z",
>     "cycleEnd": "2025-02-01T00:00:00.000Z"
>   }
> }
> ```
>
> ### Get Group
>
> /teams/groups/:groupId
>
> Retrieve a single billing group with its members and spend data for the current billing cycle.
>
> #### Parameters
>
> `groupId` string Required
>
> The encoded group ID (e.g., `group_PDSPmvukpYgZEDXsoNirw3CFhy`)
>
> `billingCycle` string
>
> ISO date string (e.g., `2025-01-15`) to specify which billing cycle to query. Defaults to current cycle.
>
> ```bash
> curl -X GET "https://api.cursor.com/teams/groups/group_PDSPmvukpYgZEDXsoNirw3CFhy?billingCycle=2025-01-15" \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "group": {
>     "id": "group_PDSPmvukpYgZEDXsoNirw3CFhy",
>     "name": "Engineering",
>     "type": "BILLING",
>     "directoryGroupId": null,
>     "memberCount": 3,
>     "createdAt": "2024-01-15T10:30:00.000Z",
>     "updatedAt": "2024-01-20T14:22:00.000Z",
>     "spendCents": 125000,
>     "currentMembers": [
>       {
>         "userId": "user_abc123",
>         "name": "Alex Developer",
>         "email": "alex@company.com",
>         "joinedAt": "2024-01-15T10:30:00.000Z",
>         "leftAt": null,
>         "spendCents": 75000,
>         "dailySpend": [
>           { "date": "2025-01-15", "spendCents": 5000 },
>           { "date": "2025-01-16", "spendCents": 7500 }
>         ]
>       },
>       {
>         "userId": "user_def456",
>         "name": "Sam Engineer",
>         "email": "sam@company.com",
>         "joinedAt": "2024-01-16T09:15:00.000Z",
>         "leftAt": null,
>         "spendCents": 50000,
>         "dailySpend": [
>           { "date": "2025-01-15", "spendCents": 3500 },
>           { "date": "2025-01-16", "spendCents": 4200 }
>         ]
>       }
>     ],
>     "formerMembers": [
>       {
>         "userId": "user_xyz789",
>         "name": "Former Member",
>         "email": "former@company.com",
>         "joinedAt": "2024-01-10T08:00:00.000Z",
>         "leftAt": "2024-01-14T17:00:00.000Z",
>         "spendCents": 0
>       }
>     ],
>     "dailySpend": [
>       { "date": "2025-01-15", "spendCents": 8500 },
>       { "date": "2025-01-16", "spendCents": 11700 }
>     ]
>   },
>   "billingCycle": {
>     "cycleStart": "2025-01-01T00:00:00.000Z",
>     "cycleEnd": "2025-02-01T00:00:00.000Z"
>   }
> }
> ```
>
> ### Create Group
>
> /teams/groups
>
> Create a new billing group. Rate limited to 20 requests per minute per team.
>
> #### Parameters
>
> `name` string Required
>
> Name of the group
>
> `type` string
>
> Group type. Currently only `BILLING` is supported. Default: `BILLING`
>
> ```bash
> curl -X POST https://api.cursor.com/teams/groups \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "name": "Engineering"
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "group": {
>     "id": "group_PDSPmvukpYgZEDXsoNirw3CFhy",
>     "name": "Engineering",
>     "type": "BILLING",
>     "directoryGroupId": null,
>     "memberCount": 0,
>     "createdAt": "2024-01-15T10:30:00.000Z",
>     "updatedAt": "2024-01-15T10:30:00.000Z",
>     "members": []
>   }
> }
> ```
>
> ### Update Group
>
> /teams/groups/:groupId
>
> Update a billing group's name or directory group attachment. Rate limited to 20 requests per minute per team.
>
> Only one field can be updated per request. To update both name and directory attachment, make separate requests.
>
> #### Parameters
>
> `groupId` string Required
>
> The encoded group ID
>
> `name` string
>
> New name for the group
>
> `directoryGroupId` string | null
>
> Directory group ID to sync with, or `null` to detach from directory sync
>
> ```bash
> curl -X PATCH https://api.cursor.com/teams/groups/group_PDSPmvukpYgZEDXsoNirw3CFhy \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "name": "Platform Engineering"
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "group": {
>     "id": "group_PDSPmvukpYgZEDXsoNirw3CFhy",
>     "name": "Platform Engineering",
>     "type": "BILLING",
>     "directoryGroupId": null,
>     "memberCount": 3,
>     "createdAt": "2024-01-15T10:30:00.000Z",
>     "updatedAt": "2024-01-25T16:45:00.000Z",
>     "members": [
>       {
>         "userId": "user_abc123",
>         "name": "Alex Developer",
>         "email": "alex@company.com",
>         "joinedAt": "2024-01-15T10:30:00.000Z"
>       }
>     ]
>   }
> }
> ```
>
> ### Delete Group
>
> /teams/groups/:groupId
>
> Delete a billing group. Returns 204 No Content on success. Rate limited to 20 requests per minute per team.
>
> Deleting a billing group is a destructive operation; data cannot be recovered. All historical usage for deleted groups is reassigned retroactively to the `Unassigned` group.
>
> #### Parameters
>
> `groupId` string Required
>
> The encoded group ID to delete
>
> ```bash
> curl -X DELETE https://api.cursor.com/teams/groups/group_PDSPmvukpYgZEDXsoNirw3CFhy \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```text
> 204 No Content
> ```
>
> ### Add Members to Group
>
> /teams/groups/:groupId/members
>
> Add team members to a billing group. Users must already be members of your team and not currently assigned to another group. Rate limited to 20 requests per minute per team.
>
> Billing groups synced with SCIM cannot be modified via the API. All member assignment for SCIM-synced groups must be handled via [SCIM](https://cursor.com/docs/account/teams/scim.md).
>
> #### Parameters
>
> `groupId` string Required
>
> The encoded group ID
>
> `userIds` string\[] Required
>
> Array of encoded user IDs to add (e.g., `["user_abc123", "user_def456"]`)
>
> ```bash
> curl -X POST https://api.cursor.com/teams/groups/group_PDSPmvukpYgZEDXsoNirw3CFhy/members \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "userIds": ["user_abc123", "user_def456"]
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "group": {
>     "id": "group_PDSPmvukpYgZEDXsoNirw3CFhy",
>     "name": "Engineering",
>     "type": "BILLING",
>     "directoryGroupId": null,
>     "memberCount": 2,
>     "createdAt": "2024-01-15T10:30:00.000Z",
>     "updatedAt": "2024-01-25T16:50:00.000Z",
>     "members": [
>       {
>         "userId": "user_abc123",
>         "name": "Alex Developer",
>         "email": "alex@company.com",
>         "joinedAt": "2024-01-25T16:50:00.000Z"
>       },
>       {
>         "userId": "user_def456",
>         "name": "Sam Engineer",
>         "email": "sam@company.com",
>         "joinedAt": "2024-01-25T16:50:00.000Z"
>       }
>     ]
>   }
> }
> ```
>
> ### Remove Members from Group
>
> /teams/groups/:groupId/members
>
> Remove team members from a billing group. Removed members are moved to the `Unassigned` group. Rate limited to 20 requests per minute per team.
>
> Billing groups synced with SCIM cannot be modified via the API. All member changes for SCIM-synced groups must be handled via [SCIM](https://cursor.com/docs/account/teams/scim.md).
>
> #### Parameters
>
> `groupId` string Required
>
> The encoded group ID
>
> `userIds` string\[] Required
>
> Array of encoded user IDs to remove
>
> ```bash
> curl -X DELETE https://api.cursor.com/teams/groups/group_PDSPmvukpYgZEDXsoNirw3CFhy/members \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "userIds": ["user_def456"]
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "group": {
>     "id": "group_PDSPmvukpYgZEDXsoNirw3CFhy",
>     "name": "Engineering",
>     "type": "BILLING",
>     "directoryGroupId": null,
>     "memberCount": 1,
>     "createdAt": "2024-01-15T10:30:00.000Z",
>     "updatedAt": "2024-01-25T17:00:00.000Z",
>     "members": [
>       {
>         "userId": "user_abc123",
>         "name": "Alex Developer",
>         "email": "alex@company.com",
>         "joinedAt": "2024-01-25T16:50:00.000Z"
>       }
>     ]
>   }
> }
> ```
>
> ## Model access
>
> Model access routes are in preview and may change. Paths, response fields, and error behavior can shift before general availability.
>
> Read and update the team's [model access](https://cursor.com/docs/enterprise/model-and-integration-management.md#model-access-control) policy: whether a custom policy is on, defaults for new providers and models, per-provider / per-model toggles, and per-model settings such as Fast and reasoning effort.
>
> Enabling a model without parameter settings leaves it on the catalog defaults. Use per-model settings when those defaults, such as Fast, do not match your team's policy.
>
> These routes return the **team baseline**. Organization Groups can still widen access for some members; group allowlists are not part of this API. Personal API key (BYOK) controls stay in the dashboard.
>
> For org-wide reads and bulk toggles across linked teams, see the [Organization API model access](https://cursor.com/docs/account/organizations/organization-admin-api.md#model-access) routes.
>
> - **Availability**: Teams with model access control enabled
> - **Authentication**: Team API key (Basic auth). Reads require **`models:read`**. Writes require **`models:*`**. Keys with **`admin:*`** work for both. Generic **`read:*`** keys cannot call these routes.
> - **Provider and model IDs**: Path segments are catalog ids such as `anthropic` and `claude-opus-4-6`, not display names. GET responses include display names.
> - **Configuration first**: Provider and model reads and writes return **409** while `state` is `unrestricted` (or `legacy`). The first `PUT /teams/model-access/configuration` with defaults on an unrestricted team turns policy on and seeds the current catalog (same idea as the first save on the Models page). Later configuration PUTs with defaults update defaults only and leave existing toggles in place.
> - **Return to unrestricted**: `PUT /teams/model-access/configuration` with `{ "state": "unrestricted" }` clears the custom policy so `state` becomes `unrestricted` again.
> - **Rate limits**: 20 requests per minute. Writes appear in team audit logs as `team_settings` events. See [rate limits and best practices](https://cursor.com/docs/api.md#rate-limits).
>
> ### Get Model Access Configuration
>
> /teams/model-access/configuration
>
> Return whether the team has a custom model-access policy and the defaults for newly seen providers and models.
>
> #### Response Fields
>
> `teamId` number
>
> Integer team ID implied by the API key.
>
> `state` string
>
> One of `unrestricted`, `custom`, or `legacy`.
>
> `newProviderDefault` string | null
>
> `enabled` or `disabled` when `state` is `custom`. Otherwise `null`.
>
> `newModelDefault` string | null
>
> `enabled` or `disabled` when `state` is `custom`. Otherwise `null`.
>
> ```bash
> curl -X GET https://api.cursor.com/teams/model-access/configuration \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "teamId": 7,
>   "state": "unrestricted",
>   "newProviderDefault": null,
>   "newModelDefault": null
> }
> ```
>
> ### Update Model Access Configuration
>
> /teams/model-access/configuration
>
> Create a custom policy, update defaults, or return the team to unrestricted.
>
> Send either:
>
> - `{ "state": "unrestricted" }` to clear the custom policy (and legacy allowed/blocked lists) so `state` becomes `unrestricted`
> - `{ "newProviderDefault", "newModelDefault" }` to create or update a custom policy (backward-compatible shorthand for `state: "custom"`)
>
> The first defaults PUT on an unrestricted team creates a custom policy and seeds catalog entries. Later defaults PUTs update defaults only and leave existing toggles in place.
>
> #### Request body
>
> `state` string
>
> Optional. Use `unrestricted` to clear policy. Omit when sending defaults.
>
> `newProviderDefault` string
>
> `enabled` or `disabled`. Required when creating or updating a custom policy; omit when `state` is `unrestricted`.
>
> `newModelDefault` string
>
> `enabled` or `disabled`. Required when creating or updating a custom policy; omit when `state` is `unrestricted`.
>
> ```bash
> curl -X PUT https://api.cursor.com/teams/model-access/configuration \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "newProviderDefault": "disabled",
>     "newModelDefault": "enabled"
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "teamId": 7,
>   "state": "custom",
>   "newProviderDefault": "disabled",
>   "newModelDefault": "enabled"
> }
> ```
>
> Return the team to unrestricted:
>
> ```bash
> curl -X PUT https://api.cursor.com/teams/model-access/configuration \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{ "state": "unrestricted" }'
> ```
>
> **Response:**
>
> ```json
> {
>   "teamId": 7,
>   "state": "unrestricted",
>   "newProviderDefault": null,
>   "newModelDefault": null
> }
> ```
>
> ### List Model Access Providers
>
> /teams/model-access/providers
>
> List catalog providers and models with resolved enabled flags and per-model `parameters`. Returns **409** when the team does not have a custom policy.
>
> Each model includes a catalog-driven `parameters` array. Parameter ids and supported values come from the model catalog (for example `fast`, `reasoning`, `effort`, `context`). Use this GET to discover which parameters a model supports before writing.
>
> #### Model `parameters` fields
>
> `id` string
>
> Parameter id (for example `fast` or `reasoning`).
>
> `displayName` string
>
> Human-readable label.
>
> `supportedValues` string\[]
>
> All values the catalog allows for this parameter on this model.
>
> `allowedValues` string\[]
>
> Values currently allowed by the team policy.
>
> `configuredDefaultValue` string | null
>
> Admin-pinned default, or `null` when unset.
>
> `catalogDefaultValue` string | null
>
> Catalog default for the parameter on this model.
>
> ```bash
> curl -X GET https://api.cursor.com/teams/model-access/providers \
>   -u YOUR_API_KEY:
> ```
>
> **Response:**
>
> ```json
> {
>   "teamId": 7,
>   "state": "custom",
>   "providers": [
>     {
>       "id": "anthropic",
>       "displayName": "Anthropic",
>       "enabled": true,
>       "models": [
>         {
>           "id": "claude-opus-4-6",
>           "displayName": "Opus 4.6",
>           "enabled": true,
>           "parameters": [
>             {
>               "id": "fast",
>               "displayName": "Fast",
>               "supportedValues": ["false", "true"],
>               "allowedValues": ["false", "true"],
>               "configuredDefaultValue": null,
>               "catalogDefaultValue": "true"
>             }
>           ]
>         }
>       ]
>     },
>     {
>       "id": "openai",
>       "displayName": "OpenAI",
>       "enabled": true,
>       "models": [
>         {
>           "id": "gpt-5.4",
>           "displayName": "GPT-5.4",
>           "enabled": true,
>           "parameters": [
>             {
>               "id": "reasoning",
>               "displayName": "Reasoning",
>               "supportedValues": ["low", "medium", "high", "xhigh", "max"],
>               "allowedValues": ["low", "medium", "high"],
>               "configuredDefaultValue": "high",
>               "catalogDefaultValue": "medium"
>             }
>           ]
>         }
>       ]
>     }
>   ]
> }
> ```
>
> ### Update Model Access Provider
>
> /teams/model-access/providers/:provider
>
> Enable or disable a provider. Returns **409** when the team is still `unrestricted` or `legacy`.
>
> #### Parameters
>
> `provider` string Required
>
> Catalog provider id (for example `openai` or `anthropic`).
>
> #### Request body
>
> `enabled` boolean Required
>
> ```bash
> curl -X PUT https://api.cursor.com/teams/model-access/providers/openai \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{"enabled": false}'
> ```
>
> ### List Models for a Provider
>
> /teams/model-access/providers/:provider/models
>
> List models for one provider with resolved enabled flags and per-model `parameters`. The parameter fields match the [providers response](https://cursor.com/docs/account/teams/admin-api.md#list-model-access-providers). Returns **409** when the team does not have a custom policy.
>
> #### Parameters
>
> `provider` string Required
>
> Catalog provider id (for example `anthropic`).
>
> ```bash
> curl -X GET https://api.cursor.com/teams/model-access/providers/anthropic/models \
>   -u YOUR_API_KEY:
> ```
>
> ### Update Model Access Model
>
> /teams/model-access/providers/:provider/models/:model
>
> Enable or disable a single model, and optionally set per-model parameter restrictions and defaults. Returns **409** when the team is still `unrestricted` or `legacy`.
>
> #### Parameters
>
> `provider` string Required
>
> Catalog provider id (for example `anthropic`).
>
> `model` string Required
>
> Catalog model id (for example `claude-opus-4-6`).
>
> #### Request body
>
> `enabled` boolean Required
>
> `parameters` object
>
> Optional map from parameter id to settings. Omitted parameters and fields are left unchanged.
>
> - `allowedValues` string\[] | null: Restrict which values members may pick. Pass `null` to clear the restriction.
> - `defaultValue` string | null: Default value for the team. Must be within `allowedValues` when a restriction is set. Pass `null` to restore the catalog default.
>
> Unknown parameter ids or values, empty `allowedValues` arrays, defaults outside `allowedValues`, and settings that resolve to no valid model variant return **400**.
>
> Disable Fast on a model:
>
> ```bash
> curl -X PUT https://api.cursor.com/teams/model-access/providers/anthropic/models/claude-opus-4-6 \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "enabled": true,
>     "parameters": {
>       "fast": { "allowedValues": ["false"] }
>     }
>   }'
> ```
>
> Set allowed reasoning levels and a default:
>
> ```bash
> curl -X PUT https://api.cursor.com/teams/model-access/providers/openai/models/gpt-5.4 \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "enabled": true,
>     "parameters": {
>       "reasoning": {
>         "allowedValues": ["low", "medium", "high"],
>         "defaultValue": "high"
>       }
>     }
>   }'
> ```
>
> Clear a restriction and restore the catalog default:
>
> ```bash
> curl -X PUT https://api.cursor.com/teams/model-access/providers/openai/models/gpt-5.4 \
>   -u YOUR_API_KEY: \
>   -H "Content-Type: application/json" \
>   -d '{
>     "enabled": true,
>     "parameters": {
>       "reasoning": {
>         "allowedValues": null,
>         "defaultValue": null
>       }
>     }
>   }'
> ```
>
> **Response:**
>
> ```json
> {
>   "id": "gpt-5.4",
>   "displayName": "GPT-5.4",
>   "enabled": true,
>   "provider": "openai",
>   "parameters": [
>     {
>       "id": "reasoning",
>       "displayName": "Reasoning",
>       "supportedValues": ["low", "medium", "high", "xhigh", "max"],
>       "allowedValues": ["low", "medium", "high", "xhigh", "max"],
>       "configuredDefaultValue": null,
>       "catalogDefaultValue": "medium"
>     }
>   ]
> }
> ```
>
> ### Errors
>
> Error bodies use:
>
> ```json
> { "code": "error", "message": "…" }
> ```
>
> | Status | When                                                                                                                                                                                                                              |
> | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `401`  | Bad key, or missing `models:read` / `models:*` (or `admin:*`)                                                                                                                                                                     |
> | `403`  | Model access control is not available for that team                                                                                                                                                                               |
> | `409`  | Provider or model read or write while `state` is `unrestricted` or `legacy`                                                                                                                                                       |
> | `400`  | Unknown provider, model, parameter id, or parameter value; invalid body; empty `allowedValues`; default outside `allowedValues`; settings that resolve to no valid model variant; or a Smart Auto required model would be blocked |
>
>
