// Package guard implements declarative checks: secrets, shell, MCP, paths.
//
// Owns: pattern matching and Ask/Deny/Allow signals for guard rules.
// Must not: route selection (dispatch), wire encode (hookedge), config merge (config).
//
// Invariants:
//   - Guards honor policy.ask_fallback when Ask is unsupported on the event.
//   - cfg.Action deny for secrets is unconditional Deny regardless of ask_fallback.
//
// Entry: AttachCheckers registry; secrets, shell, mcp, paths attach functions (called from targets/builtin).
// See DESIGN.md §1.5 (invoke_sync), §2.
package guard
