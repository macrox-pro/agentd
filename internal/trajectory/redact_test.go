package trajectory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

func TestPrepareRaw(t *testing.T) {
	t.Parallel()

	secretJSON := []byte(`{"token":"AKIAIOSFODNN7EXAMPLE"}`)

	tests := []struct {
		name string
		raw  []byte
		cfg  config.TrajectoryConfig
		want string
		nil  bool
	}{
		{name: "disabled", raw: []byte(`{}`), cfg: config.TrajectoryConfig{}, nil: true},
		{name: "no include raw", raw: []byte(`{}`), cfg: config.TrajectoryConfig{Enabled: true}, nil: true},
		{name: "empty raw", raw: nil, cfg: config.TrajectoryConfig{Enabled: true, IncludeRaw: true}, nil: true},
		{
			name: "pass through",
			raw:  []byte(`{"ok":true}`),
			cfg:  config.TrajectoryConfig{Enabled: true, IncludeRaw: true},
			want: `{"ok":true}`,
		},
		{
			name: "redact secret",
			raw:  secretJSON,
			cfg:  config.TrajectoryConfig{Enabled: true, IncludeRaw: true, RedactSecretRules: true},
			want: `{"token":"[REDACTED]"}`,
		},
		{
			name: "truncate",
			raw:  []byte(`{"hello":"world"}`),
			cfg:  config.TrajectoryConfig{Enabled: true, IncludeRaw: true, MaxEventBytes: 5},
			want: `{"hel`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := trajectory.PrepareRaw(tt.raw, tt.cfg)
			if tt.nil {
				assert.Nil(t, got, "PrepareRaw(%q)", tt.name)
				return
			}
			assert.Equal(t, tt.want, string(got), "PrepareRaw(%q)", tt.name)
		})
	}
}

func TestPrepareTranscriptText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		cfg  config.TrajectoryConfig
		want string
	}{
		{name: "empty", text: "", want: ""},
		{name: "pass through", text: "hello", cfg: config.TrajectoryConfig{}, want: "hello"},
		{
			name: "redact secret",
			text: "key AKIAIOSFODNN7EXAMPLE here",
			cfg:  config.TrajectoryConfig{RedactSecretRules: true},
			want: "[REDACTED]",
		},
		{
			name: "truncate",
			text: "abcdefghij",
			cfg:  config.TrajectoryConfig{MaxEventBytes: 4},
			want: "abcd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, trajectory.PrepareTranscriptText(tt.text, tt.cfg), "PrepareTranscriptText(%q)", tt.name)
		})
	}
}
