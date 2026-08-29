package statistics_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

func TestGate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		cfg     config.TrajectoryConfig
		wantErr error
	}{
		{
			name: "ok",
			cfg:  config.TrajectoryConfig{Enabled: true, Statistics: true},
		},
		{
			name:    "trajectory_disabled",
			cfg:     config.TrajectoryConfig{Enabled: false, Statistics: true},
			wantErr: statistics.ErrDisabled,
		},
		{
			name:    "statistics_disabled",
			cfg:     config.TrajectoryConfig{Enabled: true, Statistics: false},
			wantErr: statistics.ErrStatsOff,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := statistics.Gate(tt.cfg)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
