# Troubleshooting

> **Language:** [English](./troubleshooting.md) · [Русский](../ru/troubleshooting.md)

Failures seen in the field and how agentd behaves.

## Daemon not running

`hook run|notify|serve` prints `daemon not running` on stderr, then applies `policy.offline` from local config (defaults ⊕ user ⊕ project ⊕ runtime).

| `policy.offline` | Behavior |
|------------------|----------|
| `fail_open` (default) | Exit **0**; sync hooks encode a neutral decision so the agent continues |
| `fail_closed` | Exit **1** (blocks when provider hooks are fail-closed) |

```bash
agentd daemon start
agentd daemon status --json
```

Check `--socket` matches the edge process. Stale sockets are cleaned only under the start lock. Daemon operational logs: `agentd.log` in the [state directory](./configuration.md#state-directory).

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

## Timeouts

Sync budget ≈ 90% of provider timeout, optionally capped by route `sync_timeout` ([Dispatch](./dispatch.md)). CLI `--timeout 0` leaves deadline unset; daemon then uses kind defaults (30s / 5s).

If sync exceeds budget, context cancels; policy `fail` governs broader failure modes on the daemon side.

## Ask loops

Same tool keeps Asking → no matching approval. Extract `approval_fingerprint=` from the Ask message and run `config record-decision` ([Approvals](./approvals.md)). Wrong `--session-id` does not match session-scoped approvals.

## Async drops

`async_dropped_count` rising → queue full (`queue_capacity`). Raise capacity/workers or reduce async fan-out; `on_overflow: log` adds warnings but still drops.

## Windows named pipe

Default pipe is SID-scoped (`\\.\pipe\agentd-<SID>`). Mismatched user/session → dial failures that look like “daemon not running”.

## Codex empty no-op

Codex/Kimi no-op is **empty stdout**, exit 0 — not `{}`. Do not treat empty as failure.

## Offline policy field

`policy.offline` controls soft vs hard failure when the daemon is unreachable (see [Daemon not running](#daemon-not-running)). Invalid local YAML is treated as `fail_closed`.

## Not in v1

No agent auth, transcript pipelines, Go plugins, hooks DSL, async retry storms, or sync `exec` decisions ([DESIGN.md §11](../../DESIGN.md#11-non-goals-v1)).

See also: [CLI](./cli.md), [Operations](./operations.md).
