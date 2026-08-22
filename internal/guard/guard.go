// Package guard implements declarative checks: secrets, shell, MCP, paths.
//
// Owns: pattern matching and Ask/Deny/Allow signals for guard rules.
// Must not: route selection (dispatch), wire encode (hookedge), config merge (config).
//
// Invariants:
//   - Guards honor policy fail mode and provider caps from compiled config.
//
// Entry: secrets, shell, mcp, paths check functions (called from targets/builtin).
// See DESIGN.md §1.5 (invoke_sync), §5.
package guard
