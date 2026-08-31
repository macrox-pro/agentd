package install

import "errors"

var (
	// ErrDirRequired means install scope or provider requires an explicit target directory.
	ErrDirRequired = errors.New("install target directory is required")
	// ErrHomeRequired means user-scope install needs a home directory.
	ErrHomeRequired = errors.New("home directory is required")
	// ErrCommandRequired means the agentd binary path is required for install.
	ErrCommandRequired = errors.New("install command is required")
)
