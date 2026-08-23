// Package targets implements dispatch target adapters: builtin, exec, http, grpc, log, file.
//
// Owns: sync/async Invoke per target type; Kind→impl factories; builtin wraps agenthooks Runner.Decide.
// Must not: route matching (dispatch), YAML compile (config), wire I/O (hookedge).
//
// Invariants:
//   - Async targets never block sync pipeline return.
//   - Kind→impl mapping lives only in NewSyncInvoker / NewAsyncInvoker.
//
// Entry: NewSyncInvoker, NewAsyncInvoker, Builtin.Decide, per-target Invoke.
// See DESIGN.md §1.5 (invoke_sync, async_side), §2.
package targets
