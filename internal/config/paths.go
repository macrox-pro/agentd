package config

import (
	"os"
	"path/filepath"
)

const userConfigFileName = ".agentd.yaml"

// DefaultUserPath returns the default user config path ($HOME/.agentd.yaml).
// It returns "" when the home directory cannot be resolved.
func DefaultUserPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, userConfigFileName)
}
