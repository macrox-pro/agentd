package trajectory

import (
	"fmt"
	"os"
	"time"
)

// ForkResult summarizes a successful log fork.
type ForkResult struct {
	Provider      string `json:"provider"`
	ParentSession string `json:"parent_session"`
	NewSessionID  string `json:"new_session_id"`
	BoundarySeq   uint64 `json:"boundary_seq"`
	Copied        int    `json:"copied"`
	Path          string `json:"path"`
}

// ForkSession copies events with seq <= atSeq (or all if atSeq == 0) into a new
// session ledger and appends session/fork + session/end-seed metadata.
// The source JSONL is never modified.
func ForkSession(root string, src SessionKey, newSessionID string, atSeq uint64) (ForkResult, error) {
	if newSessionID == "" {
		return ForkResult{}, fmt.Errorf("new session id required")
	}
	if root == "" {
		root = DefaultSessionsDir()
	}
	src.Provider = CanonicalProvider(src.Provider)
	srcPath, err := FindSessionPath(root, src.Provider, src.SessionID)
	if err != nil {
		return ForkResult{}, err
	}
	events, err := ReadEvents(srcPath)
	if err != nil {
		return ForkResult{}, fmt.Errorf("read source session: %w", err)
	}
	if len(events) == 0 {
		return ForkResult{}, fmt.Errorf("source session is empty")
	}

	dst := ResolveSessionKey(src.Provider, newSessionID, src.ProjectRoot, "")
	dstPath := SessionFilePath(root, dst)
	if _, err := os.Stat(dstPath); err == nil {
		return ForkResult{}, fmt.Errorf("session %q already exists for provider %q", newSessionID, dst.Provider)
	} else if !os.IsNotExist(err) {
		return ForkResult{}, fmt.Errorf("stat destination: %w", err)
	}

	boundary := atSeq
	if boundary == 0 {
		for _, e := range events {
			if e.Seq > boundary {
				boundary = e.Seq
			}
		}
	}

	var copied []Event
	now := time.Now().UTC()
	var nextSeq uint64 = 1
	for _, e := range events {
		if e.Seq > boundary {
			continue
		}
		c := e
		c.Seq = nextSeq
		nextSeq++
		c.SessionID = dst.SessionID
		c.Provider = dst.Provider
		copied = append(copied, c)
	}
	if len(copied) == 0 {
		return ForkResult{}, fmt.Errorf("no events with seq <= %d", boundary)
	}

	forkData := mustJSON(SessionForkData{
		ParentProvider: src.Provider,
		ParentSession:  src.SessionID,
		BoundarySeq:    boundary,
	})
	seedData := mustJSON(SessionEndSeedData{
		ParentProvider: src.Provider,
		ParentSession:  src.SessionID,
		BoundarySeq:    boundary,
	})
	copied = append(copied,
		Event{
			Seq:       nextSeq,
			Type:      TypeSessionFork,
			Source:    SourceSystem,
			TS:        now,
			Provider:  dst.Provider,
			SessionID: dst.SessionID,
			Data:      forkData,
			Ignorable: true,
		},
		Event{
			Seq:       nextSeq + 1,
			Type:      TypeSessionEndSeed,
			Source:    SourceSystem,
			TS:        now,
			Provider:  dst.Provider,
			SessionID: dst.SessionID,
			Data:      seedData,
			Ignorable: true,
		},
	)

	if err := AppendEvents(root, dst, copied); err != nil {
		return ForkResult{}, fmt.Errorf("write forked session: %w", err)
	}
	return ForkResult{
		Provider:      dst.Provider,
		ParentSession: src.SessionID,
		NewSessionID:  dst.SessionID,
		BoundarySeq:   boundary,
		Copied:        len(copied) - 2,
		Path:          dstPath,
	}, nil
}
