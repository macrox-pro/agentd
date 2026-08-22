package trajectory

import (
	"fmt"
	"os"
	"time"

	"github.com/macrox-pro/agentd/internal/config"
)

// AppendImported assigns contiguous seq after existing ledger events and persists.
func AppendImported(root string, key SessionKey, events []Event) error {
	if len(events) == 0 {
		return nil
	}
	if root == "" {
		root = DefaultSessionsDir()
	}
	path := SessionFilePath(root, key)
	nextSeq := uint64(1)
	if existing, err := ReadEvents(path); err == nil && len(existing) > 0 {
		for _, e := range existing {
			if e.Seq >= nextSeq {
				nextSeq = e.Seq + 1
			}
		}
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read existing ledger: %w", err)
	}
	now := time.Now().UTC()
	for i := range events {
		events[i].Seq = nextSeq
		nextSeq++
		if events[i].TS.IsZero() {
			events[i].TS = now
		}
		events[i].Provider = key.Provider
		events[i].SessionID = key.SessionID
		if events[i].ProjectRoot == "" {
			events[i].ProjectRoot = key.ProjectRoot
		}
	}
	return AppendEvents(root, key, events)
}

// DefaultImportConfig returns conservative defaults for offline import.
func DefaultImportConfig() config.TrajectoryConfig {
	cfg := config.TrajectoryConfig{
		RedactSecretRules: true,
		MaxEventBytes:     262144,
	}
	return cfg
}
