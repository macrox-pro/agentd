package trajectory

import (
	"sync"
)

type sessionLog struct {
	mu      sync.Mutex
	events  []Event
	nextSeq uint64
	opened  bool
}

// Store holds in-memory append-only session logs.
type Store struct {
	mu       sync.RWMutex
	sessions map[SessionKey]*sessionLog
}

// NewStore returns an empty trajectory store.
func NewStore() *Store {
	return &Store{sessions: map[SessionKey]*sessionLog{}}
}

// Append adds events to the session log and returns copies of appended records.
func (s *Store) Append(key SessionKey, events []Event) []Event {
	if s == nil || len(events) == 0 {
		return nil
	}
	log := s.sessionLog(key)
	log.mu.Lock()
	defer log.mu.Unlock()

	out := make([]Event, 0, len(events))
	for _, e := range events {
		e.SessionID = key.SessionID
		e.Provider = key.Provider
		e.ProjectRoot = key.ProjectRoot
		if e.Type == TypeSessionOpen {
			if log.opened {
				continue
			}
			log.opened = true
		}
		stampSchemaVersion(&e)
		log.nextSeq++
		e.Seq = log.nextSeq
		log.events = append(log.events, e)
		out = append(out, e)
	}
	return out
}

// Events returns a copy of all events for a session key.
func (s *Store) Events(key SessionKey) []Event {
	if s == nil {
		return nil
	}
	log := s.sessionLog(key)
	log.mu.Lock()
	defer log.mu.Unlock()
	out := make([]Event, len(log.events))
	copy(out, log.events)
	return out
}

func (s *Store) sessionLog(key SessionKey) *sessionLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	if log, ok := s.sessions[key]; ok {
		return log
	}
	log := &sessionLog{}
	s.sessions[key] = log
	return log
}
