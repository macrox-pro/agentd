// Package install writes provider hook configs via agenthooks/install.
//
// Owns: manifest generation, install target resolution, discovery, plan, doctor report.
// Must not: daemon lifecycle (daemon), hook wire at runtime (hookedge).
//
// Invariants:
//   - Install merges into existing agent config files; it never removes agent home dirs.
//
// Entry: Run, RunAll, Plan, Discover, Report, ResolveDir, WriteReport.
// See DESIGN.md §6 (agentd install).
package install
