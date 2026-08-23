// Package trajectory implements the append-only session ledger (M9–M12).
//
// Owns: event catalog, in-memory store, async JSONL persist, live Hub fan-out,
// export/list/search, import append, policy replay, log fork.
// Must not: hook wire decode (hookedge), config compile (config), or own route
// matching. ReplayPolicy accepts an injected dispatch.Invoker for offline policy
// dry-run only (does not compile routes or match on the live hot path).
//
// Invariants:
//   - No disk I/O on the sync Invoke path; enqueue only.
//   - Contiguous seq per session; events immutable after append.
//   - Opt-in via config.Trajectory.Enabled (default off).
//   - schema_version frozen at SchemaVersion for v1.1 contract.
//
// Entry: Recorder.Record, Hub.Publish, ListSessions, Export, ExportToFile, Search,
// AppendImported, ReplayPolicy, ReplayPolicyFromConfig, ForkSession,
// ResolveSessionKey, ResolveSessionKeyID, EventFromSessionEvent, EventToSessionEvent.
// Import orchestration and L2 importer status: importer.ImportSession, importer.ProviderImporterStatus.
// See DESIGN.md §1.5 (async_side), §14.
package trajectory
