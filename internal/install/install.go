// Package install writes provider hook configs via agenthooks/install.
//
// Owns: manifest generation and hook shim install/uninstall.
// Must not: daemon lifecycle (daemon), hook wire at runtime (hookedge).
//
// Entry: Run.
// See DESIGN.md §6 (agentd install).
package install
