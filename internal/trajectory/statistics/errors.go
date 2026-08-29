package statistics

import "errors"

var (
	// ErrDisabled is returned when trajectory.enabled is false.
	ErrDisabled = errors.New("trajectory disabled")
	// ErrStatsOff is returned when trajectory.statistics is false.
	ErrStatsOff = errors.New("trajectory statistics disabled")
)
