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
	userNotify := config.CompiledRoute{
		Name:  "user-notify",
		Match: config.RouteMatch{Kinds: []string{"notification"}},
		Mode:  config.ModeAsyncOnly,
	}
	routes := []config.CompiledRoute{user, userNotify, def, other}

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
			name: "user beats catch-all",
			ev: &agenthooks.NotificationEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderCodex, Kind: agenthooks.KindNotification},
			},
			want: "user-notify",
		},
		{
			name: "subagent.stop uses catch-all not agent.stop",
			ev: &agenthooks.StopEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindSubagentStop},
			},
			want: "default-other",
		},
		{
			name: "exact default beats catch-all",
			ev: &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
				Tool:  agenthooks.ToolCall{Name: "Read", Canonical: agenthooks.ToolFileRead},
			},
			want: "default-tool.pre",
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

func TestMatchRouteCompiledKinds(t *testing.T) {
	t.Parallel()
	snap := testSnap(t)
	tests := []struct {
		name string
		ev   any
		want string
	}{
		{
			name: "subagent.stop vs agent.stop",
			ev: &agenthooks.StopEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindSubagentStop},
			},
			want: "default-subagent.stop",
		},
		{
			name: "agent.stop exact default",
			ev: &agenthooks.StopEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindStop},
			},
			want: "default-agent.stop",
		},
		{
			name: "catch-all does not steal tool.pre",
			ev: &agenthooks.ToolPreEvent{
				Event: agenthooks.Event{Provider: agenthooks.ProviderClaudeCode, Kind: agenthooks.KindToolPre},
				Tool:  agenthooks.ToolCall{Name: "Read", Canonical: agenthooks.ToolFileRead},
			},
			want: "default-tool.pre",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := dispatch.MatchRoute(snap.Routes, tt.ev)
			require.NotNil(t, got, "MatchRoute(%q)", tt.name)
			assert.Equal(t, tt.want, got.Name, "MatchRoute(%q)", tt.name)
		})
	}
}
