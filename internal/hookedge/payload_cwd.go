package hookedge

import (
	"encoding/json"
	"os"
)

// ResolveCWD picks cwd from hook JSON: cwd, then workspace_roots[0], then os.Getwd().
func ResolveCWD(payload []byte) string {
	var meta struct {
		CWD            string   `json:"cwd"`
		WorkspaceRoots []string `json:"workspace_roots"`
	}
	if err := json.Unmarshal(payload, &meta); err == nil {
		if meta.CWD != "" {
			return meta.CWD
		}
		if len(meta.WorkspaceRoots) > 0 && meta.WorkspaceRoots[0] != "" {
			return meta.WorkspaceRoots[0]
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return cwd
}
