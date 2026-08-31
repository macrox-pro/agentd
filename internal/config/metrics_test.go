package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestCompileMetrics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		yaml    string
		want    config.MetricsConfig
		wantErr bool
	}{
		{
			name: "defaults",
			yaml: `version: 1`,
			want: config.MetricsConfig{Enabled: false, Listen: "127.0.0.1:2112"},
		},
		{
			name: "user override",
			yaml: `version: 1
metrics:
  enabled: true
  listen: "127.0.0.1:9090"
`,
			want: config.MetricsConfig{Enabled: true, Listen: "127.0.0.1:9090"},
		},
		{
			name: "enabled false explicit vs omit",
			yaml: `version: 1
metrics:
  enabled: false
`,
			want: config.MetricsConfig{Enabled: false, Listen: "127.0.0.1:2112"},
		},
		{
			name: "invalid listen 2112",
			yaml: `version: 1
metrics:
  listen: "2112"
`,
			wantErr: true,
		},
		{
			name: "invalid listen host",
			yaml: `version: 1
metrics:
  listen: "host"
`,
			wantErr: true,
		},
		{
			name: "invalid listen host colon",
			yaml: `version: 1
metrics:
  listen: "host:"
`,
			wantErr: true,
		},
		{
			name: "enabled without listen uses default",
			yaml: `version: 1
metrics:
  enabled: true
`,
			want: config.MetricsConfig{Enabled: true, Listen: "127.0.0.1:2112"},
		},
		{
			name: "IPv6",
			yaml: `version: 1
metrics:
  enabled: true
  listen: "[::1]:2112"
`,
			want: config.MetricsConfig{Enabled: true, Listen: "[::1]:2112"},
		},
		{
			name: "colon-port",
			yaml: `version: 1
metrics:
  enabled: true
  listen: ":2112"
`,
			want: config.MetricsConfig{Enabled: true, Listen: ":2112"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "agentd.yaml")
			require.NoError(t, os.WriteFile(path, []byte(tt.yaml), 0o600), "TestCompileMetrics(%q)", tt.name)
			store, err := config.Load(context.Background(), path)
			if tt.wantErr {
				require.Error(t, err, "TestCompileMetrics(%q)", tt.name)
				return
			}
			require.NoError(t, err, "TestCompileMetrics(%q)", tt.name)
			got := store.Current().Metrics
			assert.Equal(t, tt.want, got, "TestCompileMetrics(%q)", tt.name)
		})
	}
}

func TestMetricsConfigEffective(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		cfg      config.MetricsConfig
		override string
		wantEn   bool
		wantAddr string
		wantErr  bool
	}{
		{
			name:     "CLI override enables",
			cfg:      config.MetricsConfig{Enabled: false},
			override: "127.0.0.1:9999",
			wantEn:   true,
			wantAddr: "127.0.0.1:9999",
		},
		{
			name:     "empty CLI uses YAML",
			cfg:      config.MetricsConfig{Enabled: true, Listen: "127.0.0.1:2112"},
			override: "",
			wantEn:   true,
			wantAddr: "127.0.0.1:2112",
		},
		{
			name:     "disabled when YAML off and no override",
			cfg:      config.MetricsConfig{Enabled: false},
			override: "",
			wantEn:   false,
		},
		{
			name:     "invalid CLI listen",
			cfg:      config.MetricsConfig{Enabled: false},
			override: "not-a-host-port",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			en, addr, err := tt.cfg.EffectiveListen(tt.override)
			if tt.wantErr {
				require.Error(t, err, "EffectiveListen(%q)", tt.name)
				return
			}
			require.NoError(t, err, "EffectiveListen(%q)", tt.name)
			assert.Equal(t, tt.wantEn, en, "EffectiveListen(%q) enabled", tt.name)
			assert.Equal(t, tt.wantAddr, addr, "EffectiveListen(%q) addr", tt.name)
		})
	}
}
