package guard

import (
	"github.com/speakeasy-api/agenthooks"

	"github.com/macrox-pro/agentd/internal/config"
)

const (
	checkerSecrets = "secrets"
	checkerShell   = "shell"
	checkerMCP     = "mcp"
	checkerPaths   = "paths"
)

// Checker attaches one named guard to an agenthooks Runner.
type Checker func(r *agenthooks.Runner, guards config.Guards, dctx DecisionContext)

var checkers = map[string]Checker{
	checkerSecrets: func(r *agenthooks.Runner, guards config.Guards, dctx DecisionContext) {
		AttachSecrets(r, guards.Secrets, dctx)
	},
	checkerShell: func(r *agenthooks.Runner, guards config.Guards, dctx DecisionContext) {
		AttachShell(r, guards.Shell, dctx)
	},
	checkerMCP: func(r *agenthooks.Runner, guards config.Guards, _ DecisionContext) {
		AttachMCP(r, guards.MCP)
	},
	checkerPaths: func(r *agenthooks.Runner, guards config.Guards, _ DecisionContext) {
		AttachPaths(r, guards.Paths)
	},
}

// AttachCheckers registers named guards on r in names order.
func AttachCheckers(r *agenthooks.Runner, guards config.Guards, dctx DecisionContext, names []string) {
	if r == nil {
		return
	}
	for _, name := range names {
		if attach, ok := checkers[name]; ok {
			attach(r, guards, dctx)
		}
	}
}
