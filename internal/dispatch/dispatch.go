// Package dispatch routes hook Invoke through sync and async pipelines.
//
// Owns: route match, Engine, queue, session lock, target dispatch.
// Must not: decode provider wire (hookedge), compile YAML (config), guard Decide (targets/builtin + agenthooks).
//
// Invariants:
//   - Sync response never waits on async queue drain.
//   - Invoke uses config.Snapshot only; no disk I/O on hot path.
//
// Entry: Engine.Invoke, Queue.Enqueue.
// See DESIGN.md §1.5 (invoke_sync, async_side), §2.
package dispatch
