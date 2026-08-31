// Package tui runs the interactive install/setup wizard.
//
// Owns: TTY / AGENTD_NO_TUI / CI gate; huh-based target selection and confirm.
// Must not: discovery, target-scope rules, or hook writes (call install).
//
// Entry: Interactive, RunWizard.
package tui
