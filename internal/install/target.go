package install

import (
	"fmt"

	ahinstall "github.com/speakeasy-api/agenthooks/install"

	"github.com/macrox-pro/agentd/internal/provider"
)

// Target is a resolved install destination (domain type; not ahinstall.Target).
type Target struct {
	Provider provider.ID
	Scope    string // project | user
	Dir      string // absolute install root
}

// TargetsFromHighConfidence builds install targets from high-confidence findings.
func TargetsFromHighConfidence(findings []Finding, env DiscoverEnv) ([]Target, error) {
	env = env.withDefaults()
	var out []Target
	for _, f := range findings {
		if f.Confidence != ConfidenceHigh {
			continue
		}
		ts, err := targetsFromFinding(f, env)
		if err != nil {
			return nil, err
		}
		out = append(out, ts...)
	}
	return out, nil
}

func targetsFromFinding(f Finding, env DiscoverEnv) ([]Target, error) {
	switch f.Provider {
	case provider.KimiCode:
		if f.UserDir == "" {
			return nil, nil
		}
		dir, err := ResolveDir(f.Provider, ahinstall.ScopeUser, "", false, env.Cwd, env.Home, env.Getenv)
		if err != nil {
			return nil, err
		}
		return []Target{{Provider: f.Provider, Scope: "user", Dir: dir}}, nil
	case provider.OpenCode:
		if f.ProjectDir == "" {
			return nil, nil
		}
		dir, err := ResolveDir(f.Provider, ahinstall.ScopeProject, "", false, env.Cwd, env.Home, env.Getenv)
		if err != nil {
			return nil, err
		}
		return []Target{{Provider: f.Provider, Scope: "project", Dir: dir}}, nil
	default:
		var out []Target
		if f.ProjectDir != "" {
			dir, err := ResolveDir(f.Provider, ahinstall.ScopeProject, "", false, env.Cwd, env.Home, env.Getenv)
			if err != nil {
				return nil, err
			}
			out = append(out, Target{Provider: f.Provider, Scope: "project", Dir: dir})
		}
		if f.UserDir != "" {
			dir, err := ResolveDir(f.Provider, ahinstall.ScopeUser, "", false, env.Cwd, env.Home, env.Getenv)
			if err != nil {
				return nil, err
			}
			out = append(out, Target{Provider: f.Provider, Scope: "user", Dir: dir})
		}
		return out, nil
	}
}

func resolveAHScope(scope string) (ahinstall.Scope, string, error) {
	switch scope {
	case "project":
		return ahinstall.ScopeProject, "project", nil
	case "user":
		return ahinstall.ScopeUser, "user", nil
	case "plugin":
		return ahinstall.ScopePlugin, "plugin", nil
	default:
		return "", "", fmt.Errorf("unknown scope %q", scope)
	}
}
