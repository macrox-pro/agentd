// Package statistics implements trajectory counters and session stats aggregation.
//
// Owns: daemon-lifetime Collector, session FromEvents/Load, token extractors, Gate.
// Must not: gRPC handlers, Cobra, ledger persist, hook wire decode.
//
// Invariants:
//   - Gate requires trajectory.enabled && trajectory.statistics.
//   - Daemon Observe is one increment per successful Invoke (not per ledger line).
//   - Snapshot exposes global + per-provider rollups only (no per-session state in API).
//   - Internal cursorStopLast tracks Cursor cumulative stop billing deltas until daemon restart.
//   - Token fields increment only when extractors find values.
//   - Token extraction falls back to provider transcript tail-scan when hook raw carries no usage (Codex Stop).
//
// Entry: Gate, NewCollector, Collector.Observe, Collector.Snapshot, Load, FromEvents, HookKind.
// See DESIGN.md §14.6 (async_side statistics).
package statistics
