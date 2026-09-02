// Package hookedge is the CLI wire edge: decode/encode provider hooks via
// agenthooks and forward Invoke to the local daemon.
//
// Owns: provider wire decode/encode, single gRPC Invoke forward (including Cwd on run/notify/serve),
// policy.offline when the daemon is unreachable.
// Must not: full Decide / guards / route compile (daemon dispatch + targets/builtin).
//
// Invariants:
//   - Never log to stdout on the hook path.
//   - Preserve Event.Raw verbatim.
//   - Unreachable daemon → config.OfflineFor; fail_open encodes Neutral / exit 0.
//
// Entry: Run, Notify, Serve.
// See DESIGN.md §1.5 (invoke_sync).
package hookedge
