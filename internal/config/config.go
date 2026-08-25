// Package config merges and compiles agentd configuration into immutable snapshots.
//
// Owns: four-layer merge, Compile, Store hot-path snapshot, debounced reload, persist,
// OfflineFor (edge unreachable path).
// Must not: dispatch routing (dispatch), hook wire (hookedge).
//
// Invariants:
//   - Hot path: Store.Current() only — no disk I/O per Invoke.
//   - Reload debounced; atomic snapshot swap.
//   - OfflineFor may read disk; used only when the daemon is unreachable.
//
// Entry: Store.Current, Compile, NewStore, OfflineFor.
// See DESIGN.md §1.5 (config_reload), §7.
package config
