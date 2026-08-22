package importer

import (
	"fmt"
	"os"
	"path/filepath"
)

// resolveExplicitOrSessionPath finds a transcript via --path or session id under configuredRoot.
func resolveExplicitOrSessionPath(sessionID, explicitPath, configuredRoot, label string) (string, error) {
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err != nil {
			return "", fmt.Errorf("transcript path: %w", err)
		}
		return explicitPath, nil
	}
	if sessionID == "" {
		return "", fmt.Errorf("%s import requires --path (or --session with configured import path)", label)
	}
	if configuredRoot == "" {
		return "", fmt.Errorf("%s import requires --path (no stable default transcript root)", label)
	}
	want := sessionID + ".jsonl"
	candidate := filepath.Join(configuredRoot, want)
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}
	var found string
	err := filepath.WalkDir(configuredRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == want {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil && found == "" {
		return "", fmt.Errorf("scan %s transcripts: %w", label, err)
	}
	if found == "" {
		return "", fmt.Errorf("%s transcript not found for session %q under %s", label, sessionID, configuredRoot)
	}
	return found, nil
}
