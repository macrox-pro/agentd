// Package trajectory implements the append-only session ledger (M9 L0 live path).
//
// Owns: event catalog, in-memory store, async JSONL persist, export/list read paths.
// Must not: hook wire decode (hookedge), route match (dispatch), config compile (config).
//
// Invariants:
//   - No disk I/O on the sync Invoke path; enqueue only.
//   - Contiguous seq per session; events immutable after append.
//   - Opt-in via config.Trajectory.Enabled (default off).
//
// Entry: Recorder.Record, ListSessions, ExportSession.
// See DESIGN.md §1.5 (async_side), §14.
package trajectory
