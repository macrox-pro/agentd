// Package version holds the build-time agentd version string.
//
// Owns: link-time Version var (this binary; agentd version).
// Must not: runtime feature logic; daemon Status (running process version).
package version

// Version is overridden at link time via:
//
//	-ldflags "-X github.com/macrox-pro/agentd/internal/version.Version=v1.0.0"
var Version = "dev"
