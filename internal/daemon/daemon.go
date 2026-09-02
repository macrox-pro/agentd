// Package daemon runs the user-level agentd process: start, stop, lock, status, reload,
// and login autostart registration.
//
// Owns: process lifecycle, single-instance lock, status JSON, SIGHUP reload signal,
// operational slog setup (SetupLog), login autostart (Enable, Disable, AutostartStatus),
// opt-in Prometheus metrics HTTP listen lifecycle.
// Must not: hook dispatch (dispatch), config compile (config), gRPC handlers (server).
//
// Invariants:
//   - One daemon per user (lock file).
//   - PID and lock live next to a file socket; a Windows named pipe has no
//     filesystem parent, so pipe endpoints use the per-user state directory.
//   - Reload is debounced in config.Store, not here.
//   - Logger configured once at startup; default log file under state dir.
//   - Disable removes OS autostart only; never stops a running daemon.
//
// Entry: Start, Stop, Enable, Disable, AutostartStatus, WriteStatus, SetupLog.
// See DESIGN.md §1.5 (config_reload).
package daemon
