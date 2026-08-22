package trajectory

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Persister writes JSONL lines asynchronously with debounced flush.
type Persister struct {
	root string
	mu   sync.Mutex
	pending map[SessionKey][]Event
	timer *time.Timer
}

// NewPersister returns a persister rooted at dir (typically DefaultSessionsDir()).
func NewPersister(root string) *Persister {
	return &Persister{
		root:    root,
		pending: map[SessionKey][]Event{},
	}
}

// Schedule queues events for debounced disk append.
func (p *Persister) Schedule(key SessionKey, events []Event) {
	if p == nil || p.root == "" || len(events) == 0 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pending[key] = append(p.pending[key], events...)
	if p.timer == nil {
		p.timer = time.AfterFunc(persistDebounce, func() {
			_ = p.Flush(context.Background())
		})
	} else {
		p.timer.Reset(persistDebounce)
	}
}

// Flush writes all pending events to disk.
func (p *Persister) Flush(_ context.Context) error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	batch := p.pending
	p.pending = map[SessionKey][]Event{}
	if p.timer != nil {
		p.timer.Stop()
		p.timer = nil
	}
	p.mu.Unlock()

	var firstErr error
	for key, events := range batch {
		if err := p.appendFile(key, events); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (p *Persister) appendFile(key SessionKey, events []Event) error {
	path := SessionFilePath(p.root, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("mkdir sessions: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("open session log: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range events {
		if err := enc.Encode(e); err != nil {
			return fmt.Errorf("encode event: %w", err)
		}
	}
	return nil
}
