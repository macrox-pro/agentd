// Package statistics implements trajectory counters and session stats aggregation.
//
// Owns: daemon-lifetime Collector, session FromEvents/Load, token extractors, Gate.
// Must not: gRPC handlers, Cobra, ledger persist, hook wire decode.
//
// Invariants:
//   - Gate requires trajectory.enabled && trajectory.statistics.
//   - Daemon Observe is one increment per successful Invoke (not per ledger line).
//   - No per-session map in daemon Collector (global + per-provider only).
//   - Token fields increment only when extractors find values.
//
// Entry: Gate, NewCollector, Collector.Observe, Collector.Snapshot, Load, FromEvents, HookKind.
// See DESIGN.md §14.6 (async_side statistics).
package statistics
