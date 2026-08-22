package trajectory

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SessionSummary is one on-disk session ledger.
type SessionSummary struct {
	Provider       string `json:"provider"`
	SessionID      string `json:"session_id"`
	Path           string `json:"path"`
	Lines          int    `json:"lines,omitempty"`
	ImporterStatus string `json:"importer_status,omitempty"`
}

// ListSessions scans root for session JSONL files, optionally filtered by provider.
func ListSessions(root, providerFilter string) ([]SessionSummary, error) {
	if root == "" {
		root = DefaultSessionsDir()
	}
	if root == "" {
		return nil, ErrSessionsDirUnavailable
	}
	filter := CanonicalProvider(providerFilter)
	var out []SessionSummary
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".jsonl") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		parts := strings.Split(filepath.ToSlash(rel), "/")
		if len(parts) != 2 {
			return nil
		}
		prov := parts[0]
		if filter != "" && prov != filter {
			return nil
		}
		sid := strings.TrimSuffix(parts[1], ".jsonl")
		out = append(out, SessionSummary{
			Provider:       prov,
			SessionID:      sid,
			Path:           path,
			ImporterStatus: string(ProviderImporterStatus(prov)),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Provider != out[j].Provider {
			return out[i].Provider < out[j].Provider
		}
		return out[i].SessionID < out[j].SessionID
	})
	return out, nil
}

// ReadEvents loads all events from a session JSONL file.
func ReadEvents(path string) ([]Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return decodeEvents(f)
}

// FindSessionPath resolves provider + session id to a JSONL path under root.
func FindSessionPath(root, provider, sessionID string) (string, error) {
	if root == "" {
		root = DefaultSessionsDir()
	}
	key := ResolveSessionKey(provider, sessionID, "", "")
	path := SessionFilePath(root, key)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	// Fallback: scan provider dir for filename match
	provDir := filepath.Join(root, CanonicalProvider(provider))
	entries, err := os.ReadDir(provDir)
	if err != nil {
		return "", ErrSessionNotFound
	}
	want := SessionFileName(sessionID)
	for _, e := range entries {
		if e.Name() == want {
			return filepath.Join(provDir, e.Name()), nil
		}
	}
	return "", ErrSessionNotFound
}

func decodeEvents(r io.Reader) ([]Event, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var out []Event
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e Event
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("decode event: %w", err)
		}
		out = append(out, e)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
