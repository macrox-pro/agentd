package config_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestFormatRoutes(t *testing.T) {
	t.Parallel()

	routes := []config.CompiledRoute{
		{
			Name:  "pre_tool",
			Kind:  "PreToolUse",
			Mode:  config.ModeSyncOnly,
			Match: config.RouteMatch{Kinds: []string{"PreToolUse"}},
			Sync:  []config.CompiledTarget{{Kind: config.TargetBuiltin}},
			Async: []config.CompiledTarget{{Kind: config.TargetLog}, {Kind: config.TargetFile}},
		},
	}

	tests := []struct {
		name   string
		routes []config.CompiledRoute
		asJSON bool
		want   string
		check  func(t *testing.T, raw []byte)
	}{
		{
			name:   "empty human",
			routes: nil,
			want:   "",
		},
		{
			name:   "human sync async",
			routes: routes,
			want:   "pre_tool\tmatch.kind=PreToolUse\tmode=sync_only\tsync=[builtin]\tasync=[log,file]\n",
		},
		{
			name:   "json",
			routes: routes,
			asJSON: true,
			check: func(t *testing.T, raw []byte) {
				t.Helper()
				var got []config.CompiledRoute
				require.NoError(t, json.Unmarshal(raw, &got), "json")
				require.Len(t, got, 1)
				assert.Equal(t, "pre_tool", got[0].Name)
				assert.Equal(t, config.TargetBuiltin, got[0].Sync[0].Kind)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			err := config.FormatRoutes(&buf, tt.routes, tt.asJSON)
			require.NoError(t, err, "FormatRoutes(%s)", tt.name)
			if tt.check != nil {
				tt.check(t, buf.Bytes())
				return
			}
			assert.Equal(t, tt.want, buf.String(), "FormatRoutes(%s)", tt.name)
		})
	}
}
