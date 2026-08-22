package provider

import "errors"

var (
	// ErrProviderRequired means the provider id argument was empty.
	ErrProviderRequired = errors.New("provider is required")
	// ErrUnknownProvider means the provider id is not a known canonical id or alias.
	ErrUnknownProvider = errors.New("unknown provider")
)
