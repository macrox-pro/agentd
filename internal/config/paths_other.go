//go:build !unix && !windows

package config

import (
	"os"
	"path/filepath"
)

// DefaultRuntimePath returns a best-effort state path or "".
func DefaultRuntimePath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, stateSubdir, runtimeConfigFileName)
}
