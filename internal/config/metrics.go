package config

import (
	"fmt"
	"net"
)

const defaultMetricsListen = "127.0.0.1:2112"

// MetricsConfig is compiled daemon Prometheus scrape settings.
type MetricsConfig struct {
	Enabled bool
	Listen  string
}

func defaultMetrics() MetricsConfig {
	return MetricsConfig{
		Enabled: false,
		Listen:  defaultMetricsListen,
	}
}

type fileMetrics struct {
	Enabled *bool  `yaml:"enabled,omitempty"`
	Listen  string `yaml:"listen,omitempty"`
}

func parseMetrics(in *fileMetrics, base MetricsConfig) (MetricsConfig, error) {
	out := base
	if in == nil {
		return out, nil
	}
	if in.Enabled != nil {
		out.Enabled = *in.Enabled
	}
	if in.Listen != "" {
		if err := validateListenAddr(in.Listen); err != nil {
			return MetricsConfig{}, err
		}
		out.Listen = in.Listen
	}
	return out, nil
}

func validateListenAddr(addr string) error {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("metrics.listen invalid %q: %w", addr, err)
	}
	if port == "" {
		return fmt.Errorf("metrics.listen invalid %q: missing port", addr)
	}
	return nil
}

// EffectiveListen returns whether metrics HTTP is enabled and the listen address.
// A non-empty override enables metrics and overrides the configured listen address.
func (c MetricsConfig) EffectiveListen(override string) (enabled bool, listen string, err error) {
	if override != "" {
		if err := validateListenAddr(override); err != nil {
			return false, "", err
		}
		return true, override, nil
	}
	if !c.Enabled {
		return false, "", nil
	}
	listen = c.Listen
	if listen == "" {
		listen = defaultMetricsListen
	}
	if err := validateListenAddr(listen); err != nil {
		return false, "", err
	}
	return true, listen, nil
}

func mergeMetricsPtr(base, user *fileMetrics) *fileMetrics {
	if base == nil && user == nil {
		return nil
	}
	out := fileMetrics{}
	if base != nil {
		out = *base
	}
	if user == nil {
		return &out
	}
	if user.Enabled != nil {
		out.Enabled = user.Enabled
	}
	if user.Listen != "" {
		out.Listen = user.Listen
	}
	return &out
}
