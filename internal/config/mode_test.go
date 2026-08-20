package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestNormalizeMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   config.DispatchMode
		want config.DispatchMode
	}{
		{name: "alias sync_then_async", in: config.ModeSyncThenAsync, want: config.ModeAfterSync},
		{name: "parallel unchanged", in: config.ModeParallel, want: config.ModeParallel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, config.NormalizeMode(tt.in), "NormalizeMode(%q)", tt.in)
		})
	}
}
