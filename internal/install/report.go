package install

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
)

// ChangeState classifies one installed file.
type ChangeState string

const (
	StateCreate    ChangeState = "create"
	StateUpdate    ChangeState = "update"
	StateUnchanged ChangeState = "unchanged"
)

// FileChange is one file-level install outcome.
type FileChange struct {
	Path  string
	State ChangeState
}

// Result summarizes one install invocation.
type Result struct {
	Provider string
	Scope    string
	Dir      string
	Changes  []FileChange
}

// WriteReport prints a human-readable install summary.
func WriteReport(w io.Writer, r Result) error {
	if _, err := fmt.Fprintf(w, "provider=%s scope=%s dir=%s\n", r.Provider, r.Scope, r.Dir); err != nil {
		return err
	}
	changes := append([]FileChange(nil), r.Changes...)
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Path < changes[j].Path
	})
	for _, c := range changes {
		absPath := c.Path
		if !filepath.IsAbs(absPath) {
			absPath = filepath.Join(r.Dir, filepath.FromSlash(c.Path))
		}
		if _, err := fmt.Fprintf(w, "  %-9s %s\n", c.State, absPath); err != nil {
			return err
		}
	}
	return nil
}
