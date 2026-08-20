package config

import (
	"os"
	"path/filepath"
)

// Shared path constants. Runtime default path is platform-specific
// (paths_unix.go / paths_windows.go / paths_other.go).

const (
	userConfigFileName    = ".agentd.yaml"
	runtimeConfigFileName = "runtime.yaml"
	stateSubdir           = "agentd"
)

// DefaultUserPath returns the default user config path ($HOME/.agentd.yaml).
// It returns "" when the home directory cannot be resolved.
func DefaultUserPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, userConfigFileName)
}
