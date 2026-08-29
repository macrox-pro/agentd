package daemon

import "errors"

var (
	// ErrAlreadyRunning means another daemon holds the lock or a live PID.
	ErrAlreadyRunning = errors.New("daemon already running")
	// ErrNotRunning means no daemon PID/socket was found.
	ErrNotRunning = errors.New("daemon not running")
	// ErrAutostartUnsupported means login autostart is not implemented on this GOOS.
	ErrAutostartUnsupported = errors.New("login autostart unsupported on this platform")
	// ErrAutostartNotAvailable means the OS user session backend is unavailable.
	ErrAutostartNotAvailable = errors.New("login autostart backend unavailable")
)
