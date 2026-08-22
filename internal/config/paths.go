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
	defaultLogFileName    = "agentd.log"
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

// DefaultStateDir returns the agentd state directory (parent of runtime.yaml).
func DefaultStateDir() string {
	p := DefaultRuntimePath()
	if p == "" {
		return ""
	}
	return filepath.Dir(p)
}

// DefaultLogPath returns the default daemon operational log file path.
func DefaultLogPath() string {
	dir := DefaultStateDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, defaultLogFileName)
}
