package install

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/speakeasy-api/agenthooks"
	ahinstall "github.com/speakeasy-api/agenthooks/install"

	"github.com/macrox-pro/agentd/internal/provider"
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
	Provider   string
	Scope      string
	Dir        string
	DirFlagSet bool
	Command    []string // required: abs path to agentd binary
	DryRun     bool
}

// Run writes provider hook configs via agenthooks/install.
func Run(ctx context.Context, opts Options) (Result, error) {
	m, err := Manifest(opts.Command)
	if err != nil {
		return Result{}, err
	}
	id, err := provider.Parse(opts.Provider)
	if err != nil {
		return Result{}, err
	}
	ahProv, err := id.Agenthooks()
	if err != nil {
		return Result{}, err
	}
	scope, scopeLabel, err := parseScope(opts.Scope)
	if err != nil {
		return Result{}, err
	}
	cwd, err := resolveWorkingDir()
	if err != nil {
		return Result{}, err
	}
	home, err := resolveHomeDir()
	if err != nil {
		return Result{}, err
	}
	absDir, err := ResolveDir(id, scope, opts.Dir, opts.DirFlagSet, cwd, home, os.Getenv)
	if err != nil {
		return Result{}, err
	}

	target := ahinstall.Target{
		Provider: ahProv,
		Scope:    scope,
		Dir:      absDir,
	}
	changes, err := ahinstall.Diff(m, target)
	if err != nil {
		return Result{}, err
	}
	var installOpts []ahinstall.InstallOption
	if opts.DryRun {
		installOpts = append(installOpts, ahinstall.WithDryRun())
	}
	if err := ahinstall.Install(ctx, m, target, installOpts...); err != nil {
		return Result{}, err
	}
	return Result{
		Provider: id.String(),
		Scope:    scopeLabel,
		Dir:      absDir,
		Changes:  mapChanges(changes),
	}, nil
}

func mapChanges(changes []ahinstall.Change) []FileChange {
	out := make([]FileChange, len(changes))
	for i, c := range changes {
		out[i] = FileChange{
			Path:  c.Path,
			State: ChangeState(c.State),
		}
	}
	return out
}

func parseScope(s string) (ahinstall.Scope, string, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "project":
		return ahinstall.ScopeProject, "project", nil
	case "user":
		return ahinstall.ScopeUser, "user", nil
	case "plugin":
		return ahinstall.ScopePlugin, "plugin", nil
	default:
		return "", "", fmt.Errorf("unknown scope %q", s)
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
