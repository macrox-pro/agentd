package trajectory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImportCheckpoint tracks incremental import progress for one transcript source.
type ImportCheckpoint struct {
	LastLineIndex int       `json:"last_line_index"`
	SourcePath    string    `json:"source_path"`
	SourceModTime time.Time `json:"source_mod_time"`
}

// ImportSidecarPath returns the import checkpoint path beside a session ledger.
func ImportSidecarPath(sessionsRoot, provider, sessionID string) string {
	base := strings.TrimSuffix(SessionFileName(sessionID), ".jsonl") + ".import.json"
	return filepath.Join(sessionsRoot, CanonicalProvider(provider), base)
}

// LoadImportCheckpoint reads the sidecar for a session, or zero value if missing.
func LoadImportCheckpoint(sidecarPath string) (ImportCheckpoint, error) {
	b, err := os.ReadFile(sidecarPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ImportCheckpoint{}, nil
		}
		return ImportCheckpoint{}, err
	}
	var cp ImportCheckpoint
	if err := json.Unmarshal(b, &cp); err != nil {
		return ImportCheckpoint{}, err
	}
	return cp, nil
}

// SaveImportCheckpoint writes the sidecar atomically.
func SaveImportCheckpoint(sidecarPath string, cp ImportCheckpoint) error {
	if err := os.MkdirAll(filepath.Dir(sidecarPath), 0o700); err != nil {
		return err
	}
	b, err := json.Marshal(cp)
	if err != nil {
		return err
	}
	tmp := sidecarPath + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, sidecarPath)
}
