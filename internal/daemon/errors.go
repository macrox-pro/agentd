package daemon

import "errors"

var (
	// ErrAlreadyRunning means another daemon holds the lock or a live PID.
	ErrAlreadyRunning = errors.New("daemon already running")
	// ErrNotRunning means no daemon PID/socket was found.
	ErrNotRunning = errors.New("daemon not running")
)
