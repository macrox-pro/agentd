// Package daemon runs the user-level agentd process: start, stop, lock, status, reload.
//
// Owns: process lifecycle, single-instance lock, status JSON, SIGHUP reload signal.
// Must not: hook dispatch (dispatch), config compile (config), gRPC handlers (server).
//
// Invariants:
//   - One daemon per user (lock file).
//   - Reload is debounced in config.Store, not here.
//
// Entry: Start, Stop, WriteStatus.
// See DESIGN.md §1.5 (config_reload).
package daemon
