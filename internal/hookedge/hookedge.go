// Package hookedge is the CLI wire edge: decode/encode provider hooks via
// agenthooks and forward Invoke to the local daemon.
//
// Owns: provider wire decode/encode, single gRPC Invoke forward.
// Must not: policy Decide (daemon dispatch + targets/builtin), config compile.
//
// Invariants:
//   - Never log to stdout on the hook path.
//   - Preserve Event.Raw verbatim.
//
// Entry: Run, Notify, Serve.
// See DESIGN.md §1.5 (invoke_sync).
package hookedge
