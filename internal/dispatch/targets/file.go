package targets

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// File appends one JSONL audit record per event.
type File struct {
	Logger *slog.Logger
}

type fileRecord struct {
	Time     string `json:"time"`
	Provider string `json:"provider"`
	Kind     string `json:"kind"`
	Decision string `json:"decision,omitempty"`
	Reason   string `json:"reason,omitempty"`
	RawLen   int    `json:"raw_len"`
}

// InvokeAsync appends a JSONL line to the configured path.
func (t *File) InvokeAsync(ctx context.Context, req AsyncRequest) error {
	_ = ctx
	path := req.Target.Path
	if path == "" {
		return fmt.Errorf("file target: empty path")
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("file target mkdir: %w", err)
		}
	}
	rec := fileRecord{
		Time:     time.Now().UTC().Format(time.RFC3339Nano),
		Provider: req.Provider,
		Kind:     req.EventKind,
		RawLen:   len(req.Raw),
	}
	if req.SyncOutcome != nil {
		rec.Decision = req.SyncOutcome.Kind.String()
		rec.Reason = req.SyncOutcome.Reason
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("file target open: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("file target write: %w", err)
	}
	return nil
}
