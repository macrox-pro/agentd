//go:build unix

package config

import (
	"os"
	"path/filepath"
)

const localStateSubdir = ".local/state"

// DefaultRuntimePath returns $XDG_STATE_HOME/agentd/runtime.yaml, or
// $HOME/.local/state/agentd/runtime.yaml when XDG_STATE_HOME is unset.
// Returns "" when no home/state directory can be resolved.
func DefaultRuntimePath() string {
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return filepath.Join(dir, stateSubdir, runtimeConfigFileName)
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, localStateSubdir, stateSubdir, runtimeConfigFileName)
}
