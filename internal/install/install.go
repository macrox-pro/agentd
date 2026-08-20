package install

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/speakeasy-api/agenthooks"
	ahinstall "github.com/speakeasy-api/agenthooks/install"
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
		Command:  append([]string(nil), opts.Command...),
		Hooks:    defaultHooks(),
		Identity: ahinstall.Identity{Name: "agentd", Version: "0.1.0", Description: "agentd hook proxy"},
		Fail:     agenthooks.FailClosed,
	}
	return ahinstall.Install(ctx, m, ahinstall.Target{
		Provider: provider,
		Scope:    scope,
		Dir:      absDir,
	})
}

func parseProvider(s string) (agenthooks.Provider, error) {
	switch s {
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
	switch s {
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
		{Kind: agenthooks.KindToolPre, Blocking: true, Timeout: 30 * time.Second},
		{Kind: agenthooks.KindToolPost, Blocking: false, Timeout: 5 * time.Second},
		{Kind: agenthooks.KindPromptSubmitted, Blocking: true, Timeout: 30 * time.Second},
		{Kind: agenthooks.KindStop, Blocking: true, Timeout: 5 * time.Second},
		{Kind: agenthooks.KindSessionStart, Blocking: false, Timeout: 5 * time.Second},
		{Kind: agenthooks.KindSessionEnd, Blocking: false, Timeout: 5 * time.Second},
		{Kind: agenthooks.KindNotification, Blocking: false, Timeout: 5 * time.Second},
	}
}
