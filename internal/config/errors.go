package config

import "errors"

// ErrParseConfig indicates YAML in a config file could not be unmarshaled.
var ErrParseConfig = errors.New("parse config")

// ErrUnknownToggle indicates the feature name is not in the curated toggle catalog.
var ErrUnknownToggle = errors.New("unknown toggle")

// ErrToggleAlreadySet indicates the target layer already matches the requested value.
var ErrToggleAlreadySet = errors.New("toggle already set")
