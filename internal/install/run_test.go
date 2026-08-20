package install_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		provider   string
		wantErr    bool
		wantPath   string
		wantSubstr []string
		checkShim  bool
	}{
		{
			name:     "claude project",
			provider: "claude-code",
			wantPath: filepath.Join(".claude", "settings.json"),
			wantSubstr: []string{
				"agenthooks",
				"run",
				"--provider=claude-code",
			},
		},
		{
			name:      "opencode shim",
			provider:  "opencode",
			wantPath:  filepath.Join(".opencode", "plugin", "agenthooks.ts"),
			checkShim: true,
			wantSubstr: []string{
				"opencode",
			},
		},
		{
			name:     "unknown provider",
			provider: "nope",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			bin := filepath.Join(dir, "agentd")
			require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755), "WriteFile(%q)", bin)

			err := install.Run(context.Background(), install.Options{
				Provider: tt.provider,
				Scope:    "project",
				Dir:      dir,
				Command:  []string{bin},
			})
			if tt.wantErr {
				require.Error(t, err, "Run(%q)", tt.provider)
				return
			}
			require.NoError(t, err, "Run(%q)", tt.provider)

			path := filepath.Join(dir, tt.wantPath)
			b, err := os.ReadFile(path)
			require.NoError(t, err, "ReadFile(%q)", path)
			body := string(b)
			assert.Contains(t, body, bin, "Run(%q) binary path", tt.provider)
			for _, s := range tt.wantSubstr {
				assert.Contains(t, body, s, "Run(%q) body", tt.provider)
			}
			if tt.checkShim {
				assert.True(t,
					strings.Contains(body, `"agenthooks", "serve"`) || strings.Contains(body, "agenthooks\", \"serve\""),
					"Run(%q) serve shim", tt.provider,
				)
			}
		})
	}
}
