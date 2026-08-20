// Package version holds the build-time agentd version string.
package version

// Version is overridden at link time via:
//
//	-ldflags "-X github.com/macrox-pro/agentd/internal/version.Version=v1.0.0"
var Version = "dev"
