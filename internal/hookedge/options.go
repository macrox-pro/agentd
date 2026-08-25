package hookedge

import (
	"io"
	"time"
)

// Options configures hookedge Run, Notify, and Serve.
type Options struct {
	Socket      string
	ConfigPath  string
	Provider    string
	ArgvPayload bool
	Timeout     time.Duration
	PayloadArg  string
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
}
