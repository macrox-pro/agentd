package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install"
)

func TestDoctorCLI(t *testing.T) {
	tests := []struct {
		name    string
		args    func(t *testing.T) []string
		wantErr bool
		contain []string
		check   func(t *testing.T, out string)
	}{
		{
			name: "doctor_empty_findings",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				return []string{"doctor"}
			},
			contain: []string{"no coding agents detected"},
		},
		{
			name: "doctor_human_output",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				require.NoError(t, os.MkdirAll(filepath.Join(home, ".cursor"), 0o755))
				return []string{"doctor"}
			},
			contain: []string{"cursor", "confidence=high", "daemon: unreachable"},
		},
		{
			name: "doctor_json",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				require.NoError(t, os.MkdirAll(filepath.Join(home, ".cursor"), 0o755))
				return []string{"doctor", "--json"}
			},
			check: func(t *testing.T, out string) {
				t.Helper()
				var report install.DoctorReport
				require.NoError(t, json.Unmarshal([]byte(out), &report))
				require.NotEmpty(t, report.Findings)
				assert.Equal(t, "cursor", report.Findings[0].Provider.String())
				assert.False(t, report.DaemonReachable)
			},
		},
		{
			name: "doctor_daemon_unreachable_ok",
			args: func(t *testing.T) []string {
				t.Helper()
				home := t.TempDir()
				t.Setenv("HOME", home)
				sock := filepath.Join(t.TempDir(), "missing.sock")
				return []string{"doctor", "--socket", sock}
			},
			contain: []string{"daemon: unreachable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executeRoot(t, execOpts{args: tt.args(t)})
			if tt.wantErr {
				require.Error(t, got.err)
			} else {
				require.NoError(t, got.err)
			}
			for _, s := range tt.contain {
				assert.Contains(t, got.out, s)
			}
			if tt.check != nil {
				tt.check(t, got.out)
			}
		})
	}
}
