# Getting started

> **Language:** [English](./getting-started.md) · [Русский](../ru/getting-started.md)

Five steps: install the binary, start the background service, optional user config, connect a coding agent, confirm it works.

Why this exists: [Why agentd](./why.md). Terms: [Glossary](./glossary.md).

## Terms

| Term | Meaning |
|------|---------|
| **Coding agent** | A tool that edits code for you (Claude Code, Cursor, Codex, Gemini CLI, OpenCode, Kimi Code). In commands it is `--provider`. |
| **Hook** | A small program the agent runs at an event — for example *before* it executes a shell command. |
| **Daemon** | One long-lived `agentd` process per user. Hooks send events here; policy is decided here. |
| **Install** | Write the agent’s settings so it calls `agentd` on those events. Does not start the daemon. |
| **Scope `project`** | Settings in the current repository. |
| **Scope `user`** | Settings in the agent’s home folder (for example `~/.cursor`). |

## 1. Install the binary

See [Installation](./installation.md). Shortest path (requires [Go 1.26+](https://go.dev/dl/)):

```bash
go install github.com/macrox-pro/agentd@latest
```

Ensure `$(go env GOPATH)/bin` is on your `PATH`.

## 2. Start the daemon

One process per user. The default socket depends on the OS (`--socket` overrides it).

```bash
agentd daemon start
agentd daemon status
```

`--foreground` stays attached to the terminal (useful while debugging). Logs go to `agentd.log` in the [state directory](./configuration.md#state-directory). With `--foreground`, the same lines also go to stderr.

If `~/.agentd.yaml` (or `--config`) is missing, **`daemon start` writes a minimal user config** and continues ([details](./configuration.md#user-config-bootstrap)). Edit the file after the daemon is up.

To start agentd when you log in:

```bash
agentd daemon enable
```

See [Operations → Autostart at login](./operations.md#autostart-at-login).

## 3. User config (optional at first start)

Default path: `~/.agentd.yaml` (or `--config`). Skip creating it by hand if you already ran `daemon start` — the file matches:

```yaml
version: 1
policy:
  fail: fail_closed
  # offline defaults to fail_open — agents keep working if the daemon is down
  # offline: fail_closed  # block agents when the daemon is unreachable
guards:
  secrets:
    enabled: true
    action: ask
```

| Key | Meaning |
|-----|---------|
| `fail: fail_closed` | If a check cannot run, treat it as deny. |
| `offline: fail_open` | If the daemon is down, let the agent continue (default). |
| `guards.secrets.action: ask` | On a secret match, ask the user instead of silently denying. |

Saving the file is enough: the daemon reloads it. `agentd daemon reload` forces a re-merge of all layers.

The **session ledger** (a log of hook calls) is on by default:

```bash
agentd config get trajectory    # trajectory: on (default)
```

To turn it off: `agentd config disable trajectory`.

## 4. Connect an agent (install hooks)

See what agentd finds on this machine (read-only — writes nothing):

```bash
agentd doctor
```

Preview or write hooks for every agent that already has a config folder ([Providers → Auto-detection](./providers.md#auto-detection)):

```bash
agentd install --all-detected          # print the plan; do not write
agentd install --all-detected --yes    # write the agent’s hook files
```

In an interactive terminal, use the install wizard: `agentd setup` (full flow) or `agentd install` with no flags (short flow). In CI or scripts set `AGENTD_NO_TUI=1` or `CI=true`.

One agent, this repository:

```bash
agentd install --provider=claude-code --scope=project
```

Generated settings call `agentd agenthooks …` (hidden alias of `agentd hook …` — [why](./cli.md#agenthooks-hidden)). In docs and manual settings, prefer `hook run` / `hook serve` / `hook notify`.

Per-agent paths and limits: [Providers](./providers.md).

## 5. Confirm

```bash
agentd daemon status --json
```

Expect `"running": true`. The JSON report includes config generation, fingerprint, async queue depth, and drop counters — see [Operations → Status](./operations.md#status).

Use a tool in the agent, or send a test payload:

```bash
echo '{"session_id":"s","cwd":"/tmp","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}' \
  | agentd hook run --provider=claude-code
```

A clean “before tool” event with defaults usually returns a no-op the agent understands (Claude: `{}`).

## Next

- [Configuration](./configuration.md) — layers and YAML keys
- [Guards](./guards.md) / [Dispatch](./dispatch.md) — checks and routing
- [Approvals](./approvals.md) — ask once, then allow
- [CLI](./cli.md) — commands and flags
