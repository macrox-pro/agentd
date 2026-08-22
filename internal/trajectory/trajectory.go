// Package trajectory implements the append-only session ledger (M9–M11).
//
// Owns: event catalog, in-memory store, async JSONL persist, export/list/search,
// import append, policy replay, log fork.
// Must not: hook wire decode (hookedge), route match (dispatch), config compile (config).
//
// Invariants:
//   - No disk I/O on the sync Invoke path; enqueue only.
//   - Contiguous seq per session; events immutable after append.
//   - Opt-in via config.Trajectory.Enabled (default off).
//
// Entry: Recorder.Record, ListSessions, ExportSession, Search, AppendImported, ReplayPolicy, ForkSession.
// See DESIGN.md §1.5 (async_side), §14.
package trajectory
