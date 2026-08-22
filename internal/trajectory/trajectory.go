// Package trajectory implements the append-only session ledger (M9–M11).
//
// Owns: event catalog, in-memory store, async JSONL persist, export/list/search,
// import append, policy replay, log fork.
// Must not: hook wire decode (hookedge), config compile (config), or own route
// matching. ReplayPolicy accepts an injected *dispatch.Engine for offline policy
// dry-run only (does not compile routes or match on the live hot path).
//
// Invariants:
//   - No disk I/O on the sync Invoke path; enqueue only.
//   - Contiguous seq per session; events immutable after append.
//   - Opt-in via config.Trajectory.Enabled (default off).
//
// Entry: Recorder.Record, ListSessions, Export, ExportToFile, Search, AppendImported, ReplayPolicy, ForkSession.
// See DESIGN.md §1.5 (async_side), §14.
package trajectory
