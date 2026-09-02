// Package dispatch routes hook Invoke through sync and async pipelines.
//
// Owns: route match, Engine, queue, session lock; list-level sync merge (first_conclusive);
// calls targets factories for Kind→impl.
// Must not: decode provider wire (hookedge), compile YAML (config), guard Decide (targets/builtin + agenthooks),
// Kind switch on target type (targets factories only).
//
// Invariants:
//   - Sync response never waits on async queue drain.
//   - Invoke uses config.Snapshot only; no disk I/O on hot path.
//   - policy.fail maps sync pipeline errors in Engine.Invoke (not server neutral).
//   - Optional Observer on Engine records invoke/async histograms (nil = no-op).
//
// Entry: Engine.Invoke, Queue.Enqueue.
// See DESIGN.md §1.5 (invoke_sync, async_side), §2.
package dispatch
