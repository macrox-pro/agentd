# Troubleshooting

> **Language:** [English](./troubleshooting.md) · [Русский](../ru/troubleshooting.md)

Common failures and what agentd does in each case.

## Daemon not running

`hook run|notify|serve` prints `daemon not running` on stderr, then applies `policy.offline` from local config (defaults merged with user merged with project merged with runtime).

| `policy.offline` | Behavior |
|------------------|----------|
| `fail_open` (default) | Exit **0**; sync hooks encode a neutral decision so the agent continues |
| `fail_closed` | Exit **1** (blocks when provider hooks are fail-closed) |

```bash
agentd daemon start
agentd daemon status --json
agentd daemon enable   # optional: start automatically when you log in
```

Check `--socket` matches the edge process. Daemon operational logs: `agentd.log` in the [state directory](./configuration.md#state-directory).

**Start lock:** `daemon start` acquires `agentd.lock` next to the socket (in the state directory on Windows) before removing a stale socket or PID from a crashed prior run. If another start is in progress or a live daemon holds the lock, start fails with “already running.” Only one process cleans stale files, and only under that lock.

To start agentd automatically on login: `agentd daemon enable`. If enable failed but `daemon status --json` shows `"autostart":{"enabled":true}`, fix your config and run `daemon start` or log in again — you do not need to run enable twice. See [Operations → Autostart at login](./operations.md#autostart-at-login).

## Daemon start fails (invalid user config)

`agentd daemon start` validates the user config before loading. Invalid YAML or compile errors print:

```text
agentd: invalid user config /path/to/.agentd.yaml: …
```

Start aborts; the file is **not** renamed or quarantined. Diagnose offline:

```bash
agentd config validate --config ~/.agentd.yaml
```

Fix the file, then `agentd daemon start` again. `config validate` and `config show` do not create a missing user file.

If start fails with `unknown kind`, fix the typo in `dispatch_defaults` or `match.kind` ([Dispatch](./dispatch.md#event-kinds-kind)).

## Timeouts

Sync budget ≈ 90% of provider timeout, optionally capped by route `sync_timeout` ([Dispatch](./dispatch.md)). CLI `--timeout 0` leaves deadline unset; daemon then uses kind defaults (30s / 5s).

If sync exceeds budget, context cancels; **`policy.fail`** then maps the error to allow or deny/block on the daemon side.

`policy.fail` applies when the sync pipeline returns an **error** to the engine (typical case: budget deadline). It does **not** remap a normal guard deny. For `grpc` sync targets, peer failures are usually handled by the target's `on_error` before `policy.fail` is consulted ([Dispatch](./dispatch.md)).

## OpenCode serve: daemon lost mid-stream

`hook serve` dials the daemon once at start. If a later `Invoke` fails (daemon stopped, socket error), the edge caches `policy.offline` for the rest of that NDJSON session: stderr prints `daemon not running` once, then each frame follows the cached offline mode. Restart the daemon and restart `hook serve` (or reload the OpenCode plugin) to recover daemon policy.

## Undecodable hook payload (daemon up)

If the daemon cannot decode a payload, it returns a **neutral** wire decision and skips trajectory recording for that call. This is not a hook-edge exit failure. Fix the wire JSON or provider mismatch; empty Codex/Kimi stdout is often normal ([Codex empty no-op](#codex-empty-no-op)).

## Ask loops

Same tool keeps Asking → no matching approval. Extract `approval_fingerprint=` from the Ask message and run `config record-decision` ([Approvals](./approvals.md)). Wrong `--session-id` does not match session-scoped approvals.

## Async drops

`async_dropped_count` rising → queue full (`queue_capacity`). Raise capacity/workers or reduce async fan-out; `on_overflow: log` adds warnings but still drops.

## Windows named pipe

Default pipe is scoped to your Windows user **SID** (security identifier): `\\.\pipe\agentd-<SID>`. Mismatched user or session → dial failures that look like [Daemon not running](#daemon-not-running).

## Codex empty no-op

Codex/Kimi no-op is **empty stdout**, exit 0 — not `{}`. Do not treat empty as failure.

## Offline policy field

`policy.offline` controls soft vs hard failure when the daemon is unreachable (see [Daemon not running](#daemon-not-running)). Invalid local YAML is treated as `fail_closed`.

## Out of scope for the current release

No agent auth, transcript pipelines, Go plugins, hooks DSL, async retry storms, or sync `exec` decisions ([DESIGN.md §11](../../DESIGN.md#11-non-goals-v1)).

See also: [CLI](./cli.md), [Operations](./operations.md).
