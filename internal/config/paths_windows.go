//go:build windows

package config

import (
	"os"
	"path/filepath"
)

// DefaultRuntimePath returns %LOCALAPPDATA%\agentd\runtime.yaml.
// Returns "" when LOCALAPPDATA cannot be resolved.
func DefaultRuntimePath() string {
	dir := os.Getenv("LOCALAPPDATA")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil || home == "" {
			return ""
		}
		dir = filepath.Join(home, "AppData", "Local")
	}
	return filepath.Join(dir, stateSubdir, runtimeConfigFileName)
}
