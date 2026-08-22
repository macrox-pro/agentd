package importer

import "errors"

var (
	// ErrImportNotSupported means the provider has no transcript importer.
	ErrImportNotSupported = errors.New("transcript import not supported")
	// ErrTranscriptRootRequired means a configured transcript root is required to resolve by session id.
	ErrTranscriptRootRequired = errors.New("configured transcript root required")
)
