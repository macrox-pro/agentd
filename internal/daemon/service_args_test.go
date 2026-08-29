package daemon_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestServiceStartArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts daemon.StartOptions
		want []string
	}{
		{
			name: "default_args",
			opts: daemon.StartOptions{},
			want: []string{"daemon", "start", "--foreground"},
		},
		{
			name: "includes_config_and_socket",
			opts: daemon.StartOptions{
				ConfigPath: "/home/me/.agentd.yaml",
				Socket:     "/tmp/agentd.sock",
			},
			want: []string{
				"daemon", "start", "--foreground",
				"--config", "/home/me/.agentd.yaml",
				"--socket", "/tmp/agentd.sock",
			},
		},
		{
			name: "omits_empty_log_flags",
			opts: daemon.StartOptions{LogLevel: "", LogFile: ""},
			want: []string{"daemon", "start", "--foreground"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := daemon.ServiceStartArgsForTest(tt.opts)
			assert.Equal(t, tt.want, got)
		})
	}
}
