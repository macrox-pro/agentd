// Package install writes provider hook configs via agenthooks/install.
package install

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/speakeasy-api/agenthooks"
	ahinstall "github.com/speakeasy-api/agenthooks/install"
)

const (
	identityName        = "agentd"
	identityVersion     = "0.1.0"
	identityDescription = "agentd hook proxy"

	toolPreTimeout   = 30 * time.Second
	shortHookTimeout = 5 * time.Second
)

// Options configures an install into a provider's hook settings.
type Options struct {
	Provider string
	Scope    string
	Dir      string
	Command  []string // required: abs path to agentd binary
}

// Run writes provider hook configs via agenthooks/install.
func Run(ctx context.Context, opts Options) error {
	if len(opts.Command) == 0 {
		return fmt.Errorf("command is required")
	}
	provider, err := parseProvider(opts.Provider)
	if err != nil {
		return err
	}
	scope, err := parseScope(opts.Scope)
	if err != nil {
		return err
	}
	dir := opts.Dir
	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getwd: %w", err)
		}
		dir = cwd
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("abs dir: %w", err)
	}

	m := ahinstall.Manifest{
		Command: append([]string(nil), opts.Command...),
		Hooks:   defaultHooks(),
		Identity: ahinstall.Identity{
			Name:        identityName,
			Version:     identityVersion,
			Description: identityDescription,
		},
		Fail: agenthooks.FailClosed,
	}
	return ahinstall.Install(ctx, m, ahinstall.Target{
		Provider: provider,
		Scope:    scope,
		Dir:      absDir,
	})
}

func parseProvider(s string) (agenthooks.Provider, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "claude-code":
		return agenthooks.ProviderClaudeCode, nil
	case "cursor":
		return agenthooks.ProviderCursor, nil
	case "codex":
		return agenthooks.ProviderCodex, nil
	case "gemini":
		return agenthooks.ProviderGemini, nil
	case "opencode":
		return agenthooks.ProviderOpenCode, nil
	case "kimicode", "kimi-code":
		return agenthooks.ProviderKimi, nil
	case "":
		return "", fmt.Errorf("provider is required")
	default:
		return "", fmt.Errorf("unknown provider %q", s)
	}
}

func parseScope(s string) (ahinstall.Scope, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "project":
		return ahinstall.ScopeProject, nil
	case "user":
		return ahinstall.ScopeUser, nil
	case "plugin":
		return ahinstall.ScopePlugin, nil
	default:
		return "", fmt.Errorf("unknown scope %q", s)
	}
}

func defaultHooks() []ahinstall.HookSpec {
	return []ahinstall.HookSpec{
		{Kind: agenthooks.KindToolPre, Blocking: true, Timeout: toolPreTimeout},
		{Kind: agenthooks.KindToolPost, Blocking: false, Timeout: shortHookTimeout},
		{Kind: agenthooks.KindPromptSubmitted, Blocking: true, Timeout: toolPreTimeout},
		{Kind: agenthooks.KindStop, Blocking: true, Timeout: shortHookTimeout},
		{Kind: agenthooks.KindSessionStart, Blocking: false, Timeout: shortHookTimeout},
		{Kind: agenthooks.KindSessionEnd, Blocking: false, Timeout: shortHookTimeout},
		{Kind: agenthooks.KindNotification, Blocking: false, Timeout: shortHookTimeout},
	}
}
