package dispatch_test

import (
	"testing"

	"github.com/speakeasy-api/agenthooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

func TestMatchRouteCriteria(t *testing.T) {
	t.Parallel()

	user := config.CompiledRoute{
		Name:  "gate",
		Match: config.RouteMatch{Kinds: []string{"tool.pre"}, Providers: []string{"*"}, Tools: []string{"Shell"}},
		Mode:  config.ModeParallel,
	}
	def := config.CompiledRoute{
		Name:    "default-tool.pre",
		Kind:    "tool.pre",
		Match:   config.RouteMatch{Kinds: []string{"tool.pre"}},
		Mode:    config.ModeParallel,
		Default: true,
	}
	other := config.CompiledRoute{
		Name:    "default-other",
		Kind:    "other",
		Match:   config.RouteMatch{Kinds: []string{"other"}},
		Mode:    config.ModeAsyncOnly,
		Default: true,
	}
	routes := []config.CompiledRoute{user, def, other}

	tests := []struct {
		name string
		ev   any
		want string
	}{
		{
			name: "user shell tool.pre",
			ev: &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
				Tool:  agenthooks.ToolCall{Name: "Bash", Canonical: agenthooks.ToolShell},
			},
			want: "gate",
		},
		{
			name: "default when tools do not match",
			ev: &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
				Tool:  agenthooks.ToolCall{Name: "Read", Canonical: agenthooks.ToolFileRead},
			},
			want: "default-tool.pre",
		},
		{
			name: "no match returns nil",
			ev: &agenthooks.NotificationEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderCodex, Kind: agenthooks.KindNotification},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := dispatch.MatchRoute(routes, tt.ev)
			if tt.want == "" {
				assert.Nil(t, got, "MatchRoute(%q)", tt.name)
				return
			}
			require.NotNil(t, got, "MatchRoute(%q)", tt.name)
			assert.Equal(t, tt.want, got.Name, "MatchRoute(%q)", tt.name)
		})
	}
}
