package config

import (
	"fmt"

	"github.com/macrox-pro/agentd/internal/provider"
)

const (
	defaultMaxEventBytes      = 262144
	defaultTrajectoryQueueCap = 1024
)

// ImportProviderConfig is compiled per-provider transcript import settings.
type ImportProviderConfig struct {
	Enabled bool
	Path    string
}

// TrajectoryConfig is compiled trajectory ledger settings.
type TrajectoryConfig struct {
	Enabled           bool
	Statistics        bool
	IncludeRaw        bool
	RedactSecretRules bool
	MaxEventBytes     int
	QueueCapacity     int
	Import            map[string]ImportProviderConfig
}

func defaultTrajectory() TrajectoryConfig {
	return TrajectoryConfig{
		Enabled:           true,
		Statistics:        true,
		IncludeRaw:        true,
		RedactSecretRules: true,
		MaxEventBytes:     defaultMaxEventBytes,
		QueueCapacity:     defaultTrajectoryQueueCap,
		Import:            defaultImportProviders(),
	}
}

func defaultImportProviders() map[string]ImportProviderConfig {
	return map[string]ImportProviderConfig{
		string(provider.ClaudeCode): {Enabled: false, Path: ""},
		string(provider.Cursor):     {Enabled: false, Path: ""},
		string(provider.Codex):      {Enabled: false, Path: ""},
	}
}

type fileTrajectory struct {
	Enabled           *bool                          `yaml:"enabled,omitempty"`
	Statistics        *bool                          `yaml:"statistics,omitempty"`
	IncludeRaw        *bool                          `yaml:"include_raw,omitempty"`
	RedactSecretRules *bool                          `yaml:"redact_secret_rules,omitempty"`
	MaxEventBytes     *int                           `yaml:"max_event_bytes,omitempty"`
	QueueCapacity     *int                           `yaml:"queue_capacity,omitempty"`
	Import            map[string]*fileImportProvider `yaml:"import,omitempty"`
}

type fileImportProvider struct {
	Enabled *bool  `yaml:"enabled,omitempty"`
	Path    string `yaml:"path,omitempty"`
}

func parseTrajectory(in *fileTrajectory, base TrajectoryConfig) (TrajectoryConfig, error) {
	out := base
	if in == nil {
		return out, nil
	}
	if in.Enabled != nil {
		out.Enabled = *in.Enabled
	}
	if in.Statistics != nil {
		out.Statistics = *in.Statistics
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
	if len(in.Import) > 0 {
		merged := defaultImportProviders()
		for name, prov := range in.Import {
			if prov == nil {
				continue
			}
			key := name
			if id, ok := provider.Lookup(name); ok {
				key = string(id)
			}
			cur := merged[key]
			if prov.Enabled != nil {
				cur.Enabled = *prov.Enabled
			}
			if prov.Path != "" {
				cur.Path = prov.Path
			}
			merged[key] = cur
		}
		out.Import = merged
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
	if user.Statistics != nil {
		out.Statistics = user.Statistics
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
	if len(user.Import) > 0 {
		if out.Import == nil {
			out.Import = map[string]*fileImportProvider{}
		}
		for name, prov := range user.Import {
			if prov == nil {
				continue
			}
			existing := out.Import[name]
			if existing == nil {
				existing = &fileImportProvider{}
			}
			if prov.Enabled != nil {
				existing.Enabled = prov.Enabled
			}
			if prov.Path != "" {
				existing.Path = prov.Path
			}
			out.Import[name] = existing
		}
	}
	return &out
}

// ClaudeImport returns compiled claude-code import settings.
func (c TrajectoryConfig) ClaudeImport() ImportProviderConfig {
	return c.importProvider(string(provider.ClaudeCode))
}

// CursorImport returns compiled cursor import settings.
func (c TrajectoryConfig) CursorImport() ImportProviderConfig {
	return c.importProvider(string(provider.Cursor))
}

// CodexImport returns compiled codex import settings.
func (c TrajectoryConfig) CodexImport() ImportProviderConfig {
	return c.importProvider(string(provider.Codex))
}

func (c TrajectoryConfig) importProvider(name string) ImportProviderConfig {
	if c.Import == nil {
		return ImportProviderConfig{}
	}
	return c.Import[name]
}
