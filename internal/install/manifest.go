package install

import (
	"github.com/speakeasy-api/agenthooks"
	ahinstall "github.com/speakeasy-api/agenthooks/install"
)

// Manifest builds the agenthooks install manifest for the agentd binary.
func Manifest(command []string) (ahinstall.Manifest, error) {
	if len(command) == 0 {
		return ahinstall.Manifest{}, ErrCommandRequired
	}
	return ahinstall.Manifest{
		Command: append([]string(nil), command...),
		Hooks:   defaultHooks(),
		Identity: ahinstall.Identity{
			Name:        identityName,
			Version:     identityVersion,
			Description: identityDescription,
		},
		Fail: agenthooks.FailClosed,
	}, nil
}
