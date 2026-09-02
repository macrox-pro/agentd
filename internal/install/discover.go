package install

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/macrox-pro/agentd/internal/provider"
)

// Confidence classifies how strongly a provider install target was detected.
type Confidence string

const (
	ConfidenceHigh   Confidence = "high"
	ConfidenceMedium Confidence = "medium"
)

// DiscoverEnv supplies filesystem and environment dependencies for Discover.
type DiscoverEnv struct {
	Getenv   func(string) string
	LookPath func(string) (string, error)
	Stat     func(string) (fs.FileInfo, error)
	Cwd      string
	Home     string
}

// Finding is one detected provider signal (no hook status — see Plan).
type Finding struct {
	Provider   provider.ID
	Confidence Confidence
	ProjectDir string // marker dir under Cwd when present
	UserDir    string // marker dir under home when present
	Binary     string // PATH binary when medium confidence
	Note       string
}

type discoverRule struct {
	id         provider.ID
	projectRel string // marker under cwd; empty when unsupported
	userRel    string // relative to home when getenv empty
	userEnv    string
	binaries   []string
}

var discoverRules = []discoverRule{
	{id: provider.ClaudeCode, projectRel: ".claude", userRel: ".claude", binaries: []string{"claude"}},
	{id: provider.Cursor, projectRel: ".cursor", userRel: ".cursor", binaries: []string{"cursor-agent", "cursor"}},
	{id: provider.Codex, projectRel: ".codex", userRel: ".codex", userEnv: codexHomeEnv, binaries: []string{"codex"}},
	{id: provider.Gemini, projectRel: ".gemini", userRel: ".gemini", binaries: []string{"gemini"}},
	{id: provider.OpenCode, projectRel: ".opencode", binaries: []string{"opencode"}},
	{id: provider.KimiCode, userRel: ".kimi-code", userEnv: kimiHomeEnv, binaries: []string{"kimi", "kimi-code"}},
}

// Discover scans cwd, home, and PATH for installed coding agents.
func Discover(ctx context.Context, env DiscoverEnv) ([]Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	env = env.withDefaults()
	if env.Cwd == "" {
		return nil, fmt.Errorf("working directory is required")
	}
	if env.Home == "" {
		return nil, ErrHomeRequired
	}

	var out []Finding
	for _, rule := range discoverRules {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		f, ok := discoverOne(rule, env)
		if ok {
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Provider < out[j].Provider
	})
	return out, nil
}

func (env DiscoverEnv) withDefaults() DiscoverEnv {
	if env.Getenv == nil {
		env.Getenv = os.Getenv
	}
	if env.LookPath == nil {
		env.LookPath = exec.LookPath
	}
	if env.Stat == nil {
		env.Stat = func(path string) (fs.FileInfo, error) {
			return os.Stat(path)
		}
	}
	return env
}

func discoverOne(rule discoverRule, env DiscoverEnv) (Finding, bool) {
	var projectDir, userDir string
	if rule.projectRel != "" {
		p := filepath.Join(env.Cwd, rule.projectRel)
		if isDir(env.Stat, p) {
			projectDir = p
		}
	}
	if rule.userRel != "" || rule.userEnv != "" {
		u := userMarkerDir(rule, env)
		if u != "" && isDir(env.Stat, u) {
			userDir = u
		}
	}

	if projectDir != "" || userDir != "" {
		return Finding{
			Provider:   rule.id,
			Confidence: ConfidenceHigh,
			ProjectDir: projectDir,
			UserDir:    userDir,
		}, true
	}

	for _, bin := range rule.binaries {
		path, err := env.LookPath(bin)
		if err == nil && path != "" {
			return Finding{
				Provider:   rule.id,
				Confidence: ConfidenceMedium,
				Binary:     path,
				Note:       "binary in PATH only; config dir not found",
			}, true
		}
	}
	return Finding{}, false
}

func userMarkerDir(rule discoverRule, env DiscoverEnv) string {
	if rule.userEnv != "" {
		if v := filepath.Clean(env.Getenv(rule.userEnv)); v != "" && v != "." {
			return v
		}
	}
	if rule.userRel == "" {
		return ""
	}
	return filepath.Join(env.Home, rule.userRel)
}

func isDir(stat func(string) (fs.FileInfo, error), path string) bool {
	info, err := stat(path)
	return err == nil && info.IsDir()
}
