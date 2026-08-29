package config

import "errors"

// ErrParseConfig indicates YAML in a config file could not be unmarshaled.
var ErrParseConfig = errors.New("parse config")
