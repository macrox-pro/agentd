package trajectory

import (
	"os"
	"path/filepath"
)

const sessionsSubdir = "sessions"

// DefaultSessionsDir returns the ledger root ($XDG_STATE_HOME/agentd/sessions,
// else ~/.local/state/agentd/sessions when XDG_STATE_HOME is unset).
func DefaultSessionsDir() string {
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return filepath.Join(dir, "agentd", sessionsSubdir)
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".local", "state", "agentd", sessionsSubdir)
}

// SessionFilePath returns the JSONL path for one session key.
func SessionFilePath(root string, key SessionKey) string {
	return filepath.Join(root, string(key.Provider), SessionFileName(key.SessionID))
}
