package statistics

import "github.com/macrox-pro/agentd/internal/config"

// Gate returns nil when trajectory statistics are enabled for both daemon and CLI surfaces.
func Gate(cfg config.TrajectoryConfig) error {
	if !cfg.Enabled {
		return ErrDisabled
	}
	if !cfg.Statistics {
		return ErrStatsOff
	}
	return nil
}
