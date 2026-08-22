package config

import "fmt"

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
	IncludeRaw        bool
	RedactSecretRules bool
	MaxEventBytes     int
	QueueCapacity     int
	Import            map[string]ImportProviderConfig
}

func defaultTrajectory() TrajectoryConfig {
	return TrajectoryConfig{
		Enabled:           false,
		IncludeRaw:        false,
		RedactSecretRules: true,
		MaxEventBytes:     defaultMaxEventBytes,
		QueueCapacity:     defaultTrajectoryQueueCap,
		Import:            defaultImportProviders(),
	}
}

func defaultImportProviders() map[string]ImportProviderConfig {
	return map[string]ImportProviderConfig{
		"claude-code": {Enabled: false, Path: ""},
		"cursor":      {Enabled: false, Path: ""},
		"codex":       {Enabled: false, Path: ""},
	}
}

type fileTrajectory struct {
	Enabled           *bool                         `yaml:"enabled"`
	IncludeRaw        *bool                         `yaml:"include_raw"`
	RedactSecretRules *bool                         `yaml:"redact_secret_rules"`
	MaxEventBytes     *int                          `yaml:"max_event_bytes"`
	QueueCapacity     *int                          `yaml:"queue_capacity"`
	Import            map[string]*fileImportProvider `yaml:"import"`
}

type fileImportProvider struct {
	Enabled *bool  `yaml:"enabled"`
	Path    string `yaml:"path"`
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
	if len(in.Import) > 0 {
		merged := defaultImportProviders()
		for name, prov := range in.Import {
			if prov == nil {
				continue
			}
			key := canonicalImportProvider(name)
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

func canonicalImportProvider(name string) string {
	switch name {
	case "claude-code":
		return "claude-code"
	case "cursor":
		return "cursor"
	case "codex":
		return "codex"
	default:
		return name
	}
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
	return c.importProvider("claude-code")
}

// CursorImport returns compiled cursor import settings.
func (c TrajectoryConfig) CursorImport() ImportProviderConfig {
	return c.importProvider("cursor")
}

// CodexImport returns compiled codex import settings.
func (c TrajectoryConfig) CodexImport() ImportProviderConfig {
	return c.importProvider("codex")
}

func (c TrajectoryConfig) importProvider(name string) ImportProviderConfig {
	if c.Import == nil {
		return ImportProviderConfig{}
	}
	return c.Import[name]
}
