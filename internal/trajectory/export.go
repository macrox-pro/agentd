package trajectory

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Export writes session JSONL to w. When sessionID is empty, concatenates all sessions
// for providerFilter (or all providers when filter is empty).
func Export(w io.Writer, root, providerFilter, sessionID string) error {
	if root == "" {
		root = DefaultSessionsDir()
	}
	if sessionID != "" {
		path, err := FindSessionPath(root, providerFilter, sessionID)
		if err != nil {
			return err
		}
		return copyFile(w, path)
	}
	summaries, err := ListSessions(root, providerFilter)
	if err != nil {
		return err
	}
	for i, s := range summaries {
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if err := copyFile(w, s.Path); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

// ExportToFile writes export output to outPath.
func ExportToFile(outPath, root, providerFilter, sessionID string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return Export(f, root, providerFilter, sessionID)
}
