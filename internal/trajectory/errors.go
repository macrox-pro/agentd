package trajectory

import "errors"

var (
	// ErrSessionsDirUnavailable means the default sessions state directory is unavailable.
	ErrSessionsDirUnavailable = errors.New("sessions dir unavailable")
)
