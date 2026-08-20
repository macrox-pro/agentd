package transport

import "errors"

// ErrUnsupported indicates Listen/Dial is not available on this platform.
var ErrUnsupported = errors.New("transport: unsupported platform")
