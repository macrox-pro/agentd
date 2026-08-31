# Operations

> **Language:** [English](./operations.md) · [Русский](../ru/operations.md)

Day-2 controls for the user daemon.

## One daemon per user

Lock + socket prevent a second start. Override IPC with `--socket`.

## Status

```bash
agentd daemon status
agentd daemon status --json
```

JSON when running:

| Field | Meaning |
|-------|---------|
| `running` | bool |
| `socket` | path / pipe |
| `version` | daemon version (`dev` / `dev+rev` / module version depending on how that process was built) |
| `started_at` | RFC3339 UTC |
| `generation` | config generation |
| `fingerprint` | merged config fingerprint |
| `async_queue_depth` | queued async jobs |
| `async_dropped_count` | overflow drops (monotonic) |
| `trajectory_dropped_count` | trajectory queue overflow (monotonic; when ledger enabled) |
| `compiled_route_count` | routes in snapshot |
| `metrics_listen` | Prometheus scrape address when metrics enabled; empty when disabled |

`--json` always includes an `autostart` object (even when `running` is false):

| Field | Meaning |
|-------|---------|
| `autostart.enabled` | Login autostart is registered with the OS |
| `autostart.backend` | `systemd`, `launchd`, or `schtasks` |
| `autostart.manifest_path` | Path to unit/plist when applicable |
| `autostart.registered_exe` | Binary path stored for login start |
| `autostart.stale` | `true` when `registered_exe` differs from this CLI binary (re-run `daemon enable` after upgrade) |

Human line: `agentd: running (version …, generation …)` — does not print depth/drops or autostart (use `--json`).

## Prometheus metrics

Enable in user config or at start:

```yaml
metrics:
  enabled: true
  listen: "127.0.0.1:2112"
```

Or one-shot override:

```bash
agentd daemon start --metrics-listen 127.0.0.1:2112
```

When running, `agentd daemon status --json` includes `metrics_listen` (empty when disabled). Scrape URL: `http://<metrics_listen>/metrics`.

Example Prometheus job:

```yaml
scrape_configs:
  - job_name: agentd
    static_configs:
      - targets: ["127.0.0.1:2112"]
```

Metrics listen address is fixed at **daemon start**. After changing `metrics` in YAML, run `agentd daemon restart` — `daemon reload` updates config snapshots but does not rebind the HTTP listener.

Default bind is loopback (`127.0.0.1`). Exposing metrics on `0.0.0.0` is possible but not recommended on shared machines.

## Autostart at login

Run `agentd daemon enable` once if you want agentd to start automatically when you log in. Run `agentd daemon disable` to remove that registration. Disabling autostart does **not** stop a daemon that is already running.

```bash
agentd daemon enable
agentd daemon disable
```

### What `enable` does

`daemon enable` registers agentd with your OS (systemd user unit, launchd LaunchAgent, or Windows Task Scheduler), then tries to start the daemon immediately if it is not already running.

### Command exited with an error, but autostart may still be on

If `daemon enable` prints an error, read the message carefully. **Autostart may already be turned on** even when the daemon did not start right now.

This usually means your OS is set to start agentd on login, but something blocked the immediate start — often an invalid `~/.agentd.yaml`.

**What to do:**

1. Check autostart: `agentd daemon status --json` — look at `autostart.enabled`.
2. Fix the problem shown in the error (often `agentd config validate --config ~/.agentd.yaml`).
3. Either run `agentd daemon start`, or log out and back in (or reboot) — agentd should start on login without running `enable` again.
4. To turn off autostart only: `agentd daemon disable`.

### After upgrading agentd

If you install a new binary (for example `go install …@latest`), run `agentd daemon enable` again so login autostart points at the new path. Check `autostart.stale` in `daemon status --json`.

### `disable` vs `stop`

`daemon disable` removes login autostart only. It does **not** stop the daemon. Use `daemon stop` when you want to shut down the running process.

## Reload / stop

```bash
agentd daemon reload
agentd daemon stop --timeout 10s
```

Stop drains sync then async (up to timeout), then removes socket/PID.

## Trajectory statistics

```bash
agentd config enable trajectory
agentd config enable trajectory-statistics
agentd trajectory stats [--provider ID] [--json]
agentd session stats SESSION_ID --provider ID [--json]
```

`trajectory stats` reads in-memory daemon counters (`TrajectoryService.Statistics`) and needs a running daemon; counters reset on restart. `session stats` scans a local JSONL ledger and does not need the daemon. Both require `trajectory.enabled` and `trajectory.statistics`. See [Trajectory](./trajectory.md#daemon-statistics).

## Logging

Hook path: never debug on stdout. The daemon appends operational logs to `agentd.log` in the [state directory](./configuration.md#state-directory); `agentd daemon start --foreground` also mirrors to stderr. Async dispatch `target: log` uses the same slog logger.

See also: [Dispatch](./dispatch.md) async overflow, [Configuration](./configuration.md) reload.
