package config

import "fmt"

const (
	defaultMaxEventBytes     = 262144
	defaultTrajectoryQueueCap = 1024
)

// TrajectoryConfig is compiled trajectory ledger settings.
type TrajectoryConfig struct {
	Enabled           bool
	IncludeRaw        bool
	RedactSecretRules bool
	MaxEventBytes     int
	QueueCapacity     int
}

func defaultTrajectory() TrajectoryConfig {
	return TrajectoryConfig{
		Enabled:           false,
		IncludeRaw:        false,
		RedactSecretRules: true,
		MaxEventBytes:     defaultMaxEventBytes,
		QueueCapacity:     defaultTrajectoryQueueCap,
	}
}

type fileTrajectory struct {
	Enabled           *bool  `yaml:"enabled"`
	IncludeRaw        *bool  `yaml:"include_raw"`
	RedactSecretRules *bool  `yaml:"redact_secret_rules"`
	MaxEventBytes     *int   `yaml:"max_event_bytes"`
	QueueCapacity     *int   `yaml:"queue_capacity"`
}

func parseTrajectory(in *fileTrajectory, base TrajectoryConfig) (TrajectoryConfig, error) {
	out := base
	if in == nil {
		return out, nil
	}
	if in.Enabled != nil {
		out.Enabled = *in.Enabled
	}
	if in.IncludeRaw != nil {
		out.IncludeRaw = *in.IncludeRaw
	}
	if in.RedactSecretRules != nil {
		out.RedactSecretRules = *in.RedactSecretRules
	}
	if in.MaxEventBytes != nil {
		if *in.MaxEventBytes < 0 {
			return TrajectoryConfig{}, fmt.Errorf("trajectory.max_event_bytes must be >= 0")
		}
		out.MaxEventBytes = *in.MaxEventBytes
	}
	if in.QueueCapacity != nil {
		if *in.QueueCapacity < 1 {
			return TrajectoryConfig{}, fmt.Errorf("trajectory.queue_capacity must be >= 1")
		}
		out.QueueCapacity = *in.QueueCapacity
	}
	return out, nil
}

func mergeTrajectoryPtr(base, user *fileTrajectory) *fileTrajectory {
	if base == nil && user == nil {
		return nil
	}
	out := fileTrajectory{}
	if base != nil {
		out = *base
	}
	if user == nil {
		return &out
	}
	if user.Enabled != nil {
		out.Enabled = user.Enabled
	}
	if user.IncludeRaw != nil {
		out.IncludeRaw = user.IncludeRaw
	}
	if user.RedactSecretRules != nil {
		out.RedactSecretRules = user.RedactSecretRules
	}
	if user.MaxEventBytes != nil {
		out.MaxEventBytes = user.MaxEventBytes
	}
	if user.QueueCapacity != nil {
		out.QueueCapacity = user.QueueCapacity
	}
	return &out
}
