package version_test

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/version"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name   string
		linked string
		info   *debug.BuildInfo
		want   string
	}{
		{
			name:   "ldflags semver wins",
			linked: "v0.0.3",
			info:   &debug.BuildInfo{Main: debug.Module{Version: "v9.9.9"}},
			want:   "v0.0.3",
		},
		{
			name:   "nil info defaults",
			linked: "dev",
			info:   nil,
			want:   "dev",
		},
		{
			name:   "empty linked treated as dev",
			linked: "",
			info:   nil,
			want:   "dev",
		},
		{
			name:   "go install module version",
			linked: "dev",
			info:   &debug.BuildInfo{Main: debug.Module{Version: "v0.0.3"}},
			want:   "v0.0.3",
		},
		{
			name:   "go install pseudo version",
			linked: "dev",
			info: &debug.BuildInfo{
				Main: debug.Module{Version: "v0.0.0-20260826120000-abcdef123456"},
			},
			want: "v0.0.0-20260826120000-abcdef123456",
		},
		{
			name:   "devel with short revision",
			linked: "dev",
			info: &debug.BuildInfo{
				Main: debug.Module{Version: "(devel)"},
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "abcdef1234567890"},
				},
			},
			want: "dev+abcdef1",
		},
		{
			name:   "devel dirty revision",
			linked: "dev",
			info: &debug.BuildInfo{
				Main: debug.Module{Version: "(devel)"},
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "abcdef1234567890"},
					{Key: "vcs.modified", Value: "true"},
				},
			},
			want: "dev+abcdef1-dirty",
		},
		{
			name:   "devel without vcs",
			linked: "dev",
			info:   &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}},
			want:   "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, version.Resolve(tt.linked, tt.info))
		})
	}
}
