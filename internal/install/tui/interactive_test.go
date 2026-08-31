package tui_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/install/tui"
)

func TestInteractive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		getenv func(string) string
		stdout *os.File
		want   bool
	}{
		{
			name: "no_tui_env",
			getenv: func(key string) string {
				if key == "AGENTD_NO_TUI" {
					return "1"
				}
				return ""
			},
			stdout: os.Stdout,
			want:   false,
		},
		{
			name: "ci_true",
			getenv: func(key string) string {
				if key == "CI" {
					return "true"
				}
				return ""
			},
			stdout: os.Stdout,
			want:   false,
		},
		{
			name: "nil_stdout",
			getenv: func(string) string {
				return ""
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tui.Interactive(tt.getenv, tt.stdout)
			assert.Equal(t, tt.want, got, "Interactive(%q)", tt.name)
		})
	}
}
