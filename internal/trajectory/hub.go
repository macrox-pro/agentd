package trajectory

import (
	"log/slog"
	"sync"
	"sync/atomic"
)

const hubBufferSize = 64

// SubscribeFilter selects events for one subscriber.
type SubscribeFilter struct {
	Provider  string
	SessionID string
	Source    string
}

type subscriber struct {
	ch     chan Event
	filter SubscribeFilter
}

// Hub fans out post-commit ledger events to live subscribers.
type Hub struct {
	mu     sync.RWMutex
	subs   map[uint64]*subscriber
	next   atomic.Uint64
	closed atomic.Bool
	log    *slog.Logger
}

// NewHub returns an empty subscriber registry.
func NewHub(log *slog.Logger) *Hub {
	return &Hub{
		subs: make(map[uint64]*subscriber),
		log:  log,
	}
}

// Register adds a filtered subscriber. The returned unregister func removes it.
func (h *Hub) Register(filter SubscribeFilter) (<-chan Event, func()) {
	ch := make(chan Event, hubBufferSize)
	if h == nil {
		close(ch)
		return ch, func() {}
	}
	if h.closed.Load() {
		close(ch)
		return ch, func() {}
	}
	id := h.next.Add(1)
	sub := &subscriber{ch: ch, filter: filter}
	h.mu.Lock()
	if h.closed.Load() {
		h.mu.Unlock()
		close(ch)
		return ch, func() {}
	}
	h.subs[id] = sub
	h.mu.Unlock()
	unregister := func() {
		if h == nil {
			return
		}
		h.mu.Lock()
		delete(h.subs, id)
		h.mu.Unlock()
	}
	return ch, unregister
}

// Publish delivers events to matching subscribers without blocking callers.
func (h *Hub) Publish(events []Event) {
	if h == nil || h.closed.Load() || len(events) == 0 {
		return
	}
	h.mu.RLock()
	subs := make([]*subscriber, 0, len(h.subs))
	for _, s := range h.subs {
		subs = append(subs, s)
	}
	h.mu.RUnlock()
	for _, ev := range events {
		for _, s := range subs {
			if !matchSubscribeFilter(s.filter, ev) {
				continue
			}
			select {
			case s.ch <- ev:
			default:
				if h.log != nil {
					h.log.Warn("trajectory subscribe drop; slow consumer",
						"provider", ev.Provider,
						"session_id", ev.SessionID,
						"type", ev.Type,
					)
				}
			}
		}
	}
}

// Close ends all subscriptions.
func (h *Hub) Close() {
	if h == nil || !h.closed.CompareAndSwap(false, true) {
		return
	}
	h.mu.Lock()
	for _, s := range h.subs {
		close(s.ch)
	}
	h.subs = nil
	h.mu.Unlock()
}

func matchSubscribeFilter(f SubscribeFilter, ev Event) bool {
	if f.Provider != "" && CanonicalProvider(f.Provider) != ev.Provider {
		return false
	}
	if f.SessionID != "" && f.SessionID != ev.SessionID {
		return false
	}
	if f.Source != "" && f.Source != ev.Source {
		return false
	}
	return true
}
