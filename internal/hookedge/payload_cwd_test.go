package hookedge_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/hookedge"
)

func TestResolveCWD(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   map[string]any
		want string
	}{
		{
			name: "cwd field",
			in:   map[string]any{"cwd": "/proj"},
			want: "/proj",
		},
		{
			name: "workspace roots",
			in:   map[string]any{"workspace_roots": []string{"/ws", "/other"}},
			want: "/ws",
		},
		{
			name: "cwd wins over workspace roots",
			in: map[string]any{
				"cwd":             "/proj",
				"workspace_roots": []string{"/ws"},
			},
			want: "/proj",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			payload, err := json.Marshal(tt.in)
			require.NoError(t, err)
			assert.Equal(t, tt.want, hookedge.ResolveCWD(payload))
		})
	}
}
