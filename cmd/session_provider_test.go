package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/cmd"
)

func resetCommandFlags(c *cobra.Command) {
	if c == nil {
		return
	}
	c.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	c.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, sub := range c.Commands() {
		resetCommandFlags(sub)
	}
}

func TestSessionProviderCLI(t *testing.T) {
	state := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(state, "agentd", "sessions"), 0o700))
	t.Setenv("XDG_STATE_HOME", state)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name: "list all providers",
			args: []string{"session", "list"},
		},
		{
			name:     "list unknown provider",
			args:     []string{"session", "list", "--provider", "nope"},
			wantErr:  true,
			contains: "unknown provider",
		},
		{
			name:     "import cursor without path",
			args:     []string{"session", "import", "--provider", "cursor", "--session", "x"},
			wantErr:  true,
			contains: "cursor import requires --path",
		},
		{
			name:     "import claude without session or path",
			args:     []string{"session", "import", "--provider", "claude-code"},
			wantErr:  true,
			contains: "import requires --session or --path",
		},
		{
			name:     "import unsupported provider",
			args:     []string{"session", "import", "--provider", "gemini", "--session", "x"},
			wantErr:  true,
			contains: "not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := cmd.RootCommand()
			resetCommandFlags(root)
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&out)
			root.SetArgs(tt.args)
			err := root.Execute()
			if tt.wantErr {
				require.Error(t, err, "args=%v", tt.args)
				if tt.contains != "" {
					assert.Contains(t, err.Error(), tt.contains)
				}
				return
			}
			require.NoError(t, err, "args=%v", tt.args)
			if tt.name == "list all providers" {
				assert.Contains(t, out.String(), "no sessions")
			}
		})
	}
}
