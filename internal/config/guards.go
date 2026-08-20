package config

import "fmt"

const guardNameSecrets = "secrets"

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

// Guards holds compiled guard settings.
type Guards struct {
	Secrets SecretsGuard
}

func parseGuards(fg *fileGuards, def Guards) (Guards, error) {
	out := def
	if fg == nil || fg.Secrets == nil {
		return out, nil
	}
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
	return out, nil
}
