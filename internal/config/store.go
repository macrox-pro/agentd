package config

import "sync/atomic"

// Snapshot is an immutable compiled configuration generation.
type Snapshot struct {
	Generation  uint64
	Fingerprint string
}

// Store holds the current config snapshot for lock-free reads on the hot path.
type Store struct {
	snap atomic.Pointer[Snapshot]
}

// NewStore returns a store with an empty initial snapshot.
func NewStore() *Store {
	s := &Store{}
	s.snap.Store(&Snapshot{})
	return s
}

// Current returns the active snapshot.
func (s *Store) Current() *Snapshot {
	return s.snap.Load()
}
