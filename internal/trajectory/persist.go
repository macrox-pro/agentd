package trajectory

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Persister writes JSONL lines asynchronously with debounced flush.
type Persister struct {
	root    string
	log     *slog.Logger
	mu      sync.Mutex
	flushMu sync.Mutex
	pending map[SessionKey][]Event
	wake    chan struct{}
	done    chan struct{}
	once    sync.Once
}

// NewPersister returns a persister rooted at dir (typically DefaultSessionsDir()).
func NewPersister(root string, log *slog.Logger) *Persister {
	p := &Persister{
		root:    root,
		log:     log,
		pending: map[SessionKey][]Event{},
		wake:    make(chan struct{}, 1),
		done:    make(chan struct{}),
	}
	p.once.Do(func() { go p.loop() })
	return p
}

// Schedule queues events for debounced disk append.
func (p *Persister) Schedule(key SessionKey, events []Event) {
	if p == nil || p.root == "" || len(events) == 0 {
		return
	}
	p.mu.Lock()
	p.pending[key] = append(p.pending[key], events...)
	p.mu.Unlock()
	select {
	case p.wake <- struct{}{}:
	default:
	}
}

func (p *Persister) loop() {
	for {
		select {
		case <-p.done:
			return
		case <-p.wake:
		}
		timer := time.NewTimer(persistDebounce)
		for {
			select {
			case <-p.done:
				timer.Stop()
				return
			case <-p.wake:
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(persistDebounce)
			case <-timer.C:
				if err := p.Flush(context.Background()); err != nil && p.log != nil {
					p.log.Warn("trajectory persist flush failed", "error", err)
				}
				goto next
			}
		}
	next:
	}
}

// Flush writes all pending events to disk.
func (p *Persister) Flush(_ context.Context) error {
	if p == nil {
		return nil
	}
	p.flushMu.Lock()
	defer p.flushMu.Unlock()

	p.mu.Lock()
	if len(p.pending) == 0 {
		p.mu.Unlock()
		return nil
	}
	batch := p.pending
	p.pending = map[SessionKey][]Event{}
	p.mu.Unlock()

	var firstErr error
	failed := map[SessionKey][]Event{}
	for key, events := range batch {
		if err := p.appendFile(key, events); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			failed[key] = events
		}
	}
	if len(failed) > 0 {
		p.mu.Lock()
		for key, events := range failed {
			p.pending[key] = append(events, p.pending[key]...)
		}
		select {
		case p.wake <- struct{}{}:
		default:
		}
		p.mu.Unlock()
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

// AppendEvents writes events to the session JSONL file immediately (offline import).
func AppendEvents(root string, key SessionKey, events []Event) error {
	if len(events) == 0 {
		return nil
	}
	p := NewPersister(root, nil)
	return p.appendFile(key, events)
}
