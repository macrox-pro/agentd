// Package version holds the build-time agentd version string.
//
// Owns: link-time Version; String/Resolve (BuildInfo for go install / local VCS).
// Must not: runtime feature logic; daemon Status (running process version).
//
// Invariants:
//   - Non-dev linked Version always wins over BuildInfo.
//   - go install @tag/@latest uses info.Main.Version (semver or pseudo-version).
//   - Local (devel) builds fall back to short vcs.revision when stamped.
//
// Entry: String, Version.
package version

import "runtime/debug"

const (
	devVersion  = "dev"
	develMain   = "(devel)"
	shortRevLen = 7
)

// Version is overridden at link time via:
//
//	-ldflags "-X github.com/macrox-pro/agentd/internal/version.Version=v1.0.0"
//
// Leave as "dev" so String falls back to BuildInfo (go install / local VCS).
var Version = devVersion

// String returns the effective version for this binary.
func String() string {
	info, _ := debug.ReadBuildInfo()
	return Resolve(Version, info)
}

// Resolve selects the version string. Used by String; exported for table tests.
func Resolve(linked string, info *debug.BuildInfo) string {
	if linked != "" && linked != devVersion {
		return linked
	}
	if info == nil {
		return devVersion
	}
	if v := info.Main.Version; v != "" && v != develMain {
		return v
	}

	var (
		rev   string
		dirty bool
	)
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}
	if rev == "" {
		return devVersion
	}
	if len(rev) > shortRevLen {
		rev = rev[:shortRevLen]
	}
	if dirty {
		return devVersion + "+" + rev + "-dirty"
	}
	return devVersion + "+" + rev
}
