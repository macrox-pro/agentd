package config

import (
	"os"
	"path/filepath"
)

const projectConfigFileName = ".agentd.yaml"

// FindProjectConfig walks ancestors of startDir looking for .agentd.yaml.
// When projectRoot is non-empty and contains .agentd.yaml, that path wins.
// Returns the absolute path and true when found.
func FindProjectConfig(startDir, projectRoot string) (string, bool) {
	if projectRoot != "" {
		if p, ok := projectFileIfExists(projectRoot); ok {
			return p, true
		}
	}
	if startDir == "" {
		return "", false
	}
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", false
	}
	for {
		if p, ok := projectFileIfExists(dir); ok {
			return p, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func projectFileIfExists(dir string) (string, bool) {
	if dir == "" {
		return "", false
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}
	p := filepath.Join(abs, projectConfigFileName)
	fi, err := os.Stat(p)
	if err != nil || fi.IsDir() {
		return "", false
	}
	return p, true
}
