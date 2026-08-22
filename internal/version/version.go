// Package version holds the build-time agentd version string.
//
// Owns: link-time Version var.
// Must not: runtime feature logic.
package version

// Version is overridden at link time via:
//
//	-ldflags "-X github.com/macrox-pro/agentd/internal/version.Version=v1.0.0"
var Version = "dev"
