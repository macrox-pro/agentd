// Package config merges and compiles agentd configuration into immutable snapshots.
//
// Owns: four-layer merge, CompileMerged, Store hot-path snapshot, debounced reload, persist,
// PrepareUserConfig (daemon-start user bootstrap), SetToggle/GetToggle (CLI feature toggles),
// OfflineFor (edge unreachable path).
// Must not: dispatch routing (dispatch), hook wire (hookedge).
//
// Invariants:
//   - Hot path: Store.Current() only — no disk I/O per Invoke.
//   - Reload debounced; atomic snapshot swap.
//   - OfflineFor may read disk; used only when the daemon is unreachable.
//   - PrepareUserConfig runs only from daemon start; Load/LoadWith never bootstrap.
//   - SetToggle may bootstrap the user config file when missing (same shape as PrepareUserConfig).
//
// Entry: Store.Current, CompileMerged, PrepareUserConfig, SetToggle, GetToggle, LookupToggle, OfflineFor.
// See DESIGN.md §1.5 (config_reload), §7.
package config
