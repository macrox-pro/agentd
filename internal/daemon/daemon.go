// Package daemon runs the user-level agentd process: start, stop, lock, status, reload,
// and login autostart registration.
//
// Owns: process lifecycle, single-instance lock, status JSON, SIGHUP reload signal,
// operational slog setup (SetupLog), login autostart (Enable, Disable, AutostartStatus).
// Must not: hook dispatch (dispatch), config compile (config), gRPC handlers (server).
//
// Invariants:
//   - One daemon per user (lock file).
//   - Reload is debounced in config.Store, not here.
//   - Logger configured once at startup; default log file under state dir.
//   - Disable removes OS autostart only; never stops a running daemon.
//
// Entry: Start, Stop, Enable, Disable, AutostartStatus, WriteStatus, SetupLog.
// See DESIGN.md §1.5 (config_reload).
package daemon
