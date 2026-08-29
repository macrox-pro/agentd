package trajectory

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// WriteEvents encodes events as JSONL to w (one JSON object per line).
func WriteEvents(w io.Writer, events []Event) error {
	enc := json.NewEncoder(w)
	for i := range events {
		stampSchemaVersion(&events[i])
		if err := enc.Encode(events[i]); err != nil {
			return fmt.Errorf("encode event: %w", err)
		}
	}
	return nil
}

// WriteEventsToFile writes events as JSONL to path (truncate/create).
func WriteEventsToFile(path string, events []Event) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("mkdir output dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open output file: %w", err)
	}
	defer f.Close()
	if err := WriteEvents(f, events); err != nil {
		return err
	}
	return nil
}
