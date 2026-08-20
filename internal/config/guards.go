package config

import "fmt"

const (
	guardNameSecrets = "secrets"
	guardNameShell   = "shell"
	guardNameMCP     = "mcp"
	guardNamePaths   = "paths"
)

type GuardAction string

const (
	GuardAsk  GuardAction = "ask"
	GuardDeny GuardAction = "deny"
)

// SecretsGuard is compiled secrets guard settings.
type SecretsGuard struct {
	Enabled bool
	Action  GuardAction
	Rules   []string
}

// ShellGuard is compiled shell command guard settings.
type ShellGuard struct {
	Enabled      bool
	DenyPatterns []string
	AskOn        []string
}

// MCPGuard is compiled MCP server deny settings.
type MCPGuard struct {
	Enabled     bool
	DenyServers []string
}

// PathsGuard is compiled filesystem path deny settings.
type PathsGuard struct {
	Enabled   bool
	DenyRead  []string
	DenyWrite []string
}

// Guards holds compiled guard settings.
type Guards struct {
	Secrets SecretsGuard
	Shell   ShellGuard
	MCP     MCPGuard
	Paths   PathsGuard
}

func knownGuardName(name string) bool {
	switch name {
	case guardNameSecrets, guardNameShell, guardNameMCP, guardNamePaths:
		return true
	default:
		return false
	}
}

// enabledGuardNames returns enabled guard names in attach order.
func enabledGuardNames(g Guards) []string {
	var out []string
	if g.Secrets.Enabled {
		out = append(out, guardNameSecrets)
	}
	if g.Shell.Enabled {
		out = append(out, guardNameShell)
	}
	if g.MCP.Enabled {
		out = append(out, guardNameMCP)
	}
	if g.Paths.Enabled {
		out = append(out, guardNamePaths)
	}
	return out
}

func parseGuards(fg *fileGuards, def Guards) (Guards, error) {
	out := def
	if fg == nil {
		return out, nil
	}
	if fg.Secrets != nil {
		s := fg.Secrets
		if s.Enabled != nil {
			out.Secrets.Enabled = *s.Enabled
		}
		if s.Action != "" {
			switch GuardAction(s.Action) {
			case GuardAsk, GuardDeny:
				out.Secrets.Action = GuardAction(s.Action)
			default:
				return Guards{}, fmt.Errorf("guards.secrets.action: unknown %q", s.Action)
			}
		}
		if s.Rules != nil {
			out.Secrets.Rules = append([]string(nil), s.Rules...)
		}
	}
	if fg.Shell != nil {
		s := fg.Shell
		if s.Enabled != nil {
			out.Shell.Enabled = *s.Enabled
		}
		if s.DenyPatterns != nil {
			out.Shell.DenyPatterns = append([]string(nil), s.DenyPatterns...)
		}
		if s.AskOn != nil {
			out.Shell.AskOn = append([]string(nil), s.AskOn...)
		}
	}
	if fg.MCP != nil {
		m := fg.MCP
		if m.Enabled != nil {
			out.MCP.Enabled = *m.Enabled
		}
		if m.DenyServers != nil {
			out.MCP.DenyServers = append([]string(nil), m.DenyServers...)
		}
	}
	if fg.Paths != nil {
		p := fg.Paths
		if p.Enabled != nil {
			out.Paths.Enabled = *p.Enabled
		}
		if p.DenyRead != nil {
			out.Paths.DenyRead = append([]string(nil), p.DenyRead...)
		}
		if p.DenyWrite != nil {
			out.Paths.DenyWrite = append([]string(nil), p.DenyWrite...)
		}
	}
	return out, nil
}
