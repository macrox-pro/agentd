package trajectory

import "errors"

var (
	// ErrSessionsDirUnavailable means the default sessions state directory is unavailable.
	ErrSessionsDirUnavailable = errors.New("sessions dir unavailable")
	// ErrSessionNotFound means no session JSONL exists for the provider and session id.
	ErrSessionNotFound = errors.New("session not found")
	// ErrNewSessionIDRequired means fork requires a non-empty destination session id.
	ErrNewSessionIDRequired = errors.New("new session id required")
	// ErrSourceSessionEmpty means the fork source ledger has no events.
	ErrSourceSessionEmpty = errors.New("source session is empty")
	// ErrSessionAlreadyExists means the fork destination session already exists.
	ErrSessionAlreadyExists = errors.New("session already exists")
	// ErrReplayNoRaw means no hook/invoked events have stored Raw payloads.
	ErrReplayNoRaw = errors.New("policy replay requires stored raw payloads at record time")
	// ErrReplayNoEvents means the session has no hook/invoked events to replay.
	ErrReplayNoEvents = errors.New("no hook/invoked events to replay")
	// ErrReplaySeqNotFound means no hook/invoked event with Raw exists for the requested seq.
	ErrReplaySeqNotFound = errors.New("no hook/invoked event with raw payload for seq")
	// ErrNilConfigSnap means ReplayPolicy was called without a config snapshot.
	ErrNilConfigSnap = errors.New("nil config snapshot")
	// ErrNilEngine means ReplayPolicy was called without a dispatch engine.
	ErrNilEngine = errors.New("nil dispatch engine")
)
