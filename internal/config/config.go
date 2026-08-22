// Package config merges and compiles agentd configuration into immutable snapshots.
//
// Owns: four-layer merge, Compile, Store hot-path snapshot, debounced reload, persist.
// Must not: dispatch routing (dispatch), hook wire (hookedge).
//
// Invariants:
//   - Hot path: Store.Current() only — no disk I/O per Invoke.
//   - Reload debounced; atomic snapshot swap.
//
// Entry: Store.Current, Compile, NewStore.
// See DESIGN.md §1.5 (config_reload), §7.
package config
