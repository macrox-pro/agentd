package dispatch

import (
	"sync"
	"sync/atomic"

	"github.com/speakeasy-api/agenthooks"
)

// Sessions serializes sync Invoke per session id (DESIGN §8).
type Sessions struct {
	active atomic.Int32
	gates  sync.Map // string -> *sync.Mutex
}

// Lock acquires the per-session mutex. Empty id is a no-op.
// The returned unlock must be called once.
func (s *Sessions) Lock(id string) (unlock func()) {
	if s == nil || id == "" {
		return func() {}
	}
	v, _ := s.gates.LoadOrStore(id, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	s.active.Add(1)
	return func() {
		s.active.Add(-1)
		mu.Unlock()
	}
}

// Active returns how many session locks are currently held.
func (s *Sessions) Active() uint32 {
	if s == nil {
		return 0
	}
	n := s.active.Load()
	if n < 0 {
		return 0
	}
	return uint32(n)
}

// SessionIDOf extracts the normalized session id from a typed event.
func SessionIDOf(typed any) string {
	if base := agenthooks.EventOf(typed); base != nil {
		return base.Session.ID
	}
	return ""
}
