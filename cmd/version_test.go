package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/version"
)

func TestVersionCLI(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "ldflags semver",
			version: "0.0.2",
			want:    "0.0.2\n",
		},
		{
			name:    "ldflags v-prefixed",
			version: "v0.0.4",
			want:    "v0.0.4\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prev := version.Version
			version.Version = tt.version
			t.Cleanup(func() { version.Version = prev })

			got := executeRoot(t, execOpts{args: []string{"version"}})
			require.NoError(t, got.err)
			assert.Equal(t, tt.want, got.out)
		})
	}
}
