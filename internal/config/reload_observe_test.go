package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestReloadObserved(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T) (*config.Store, string)
		reloadErr  bool
		wantResult string
	}{
		{
			name: "result ok",
			setup: func(t *testing.T) (*config.Store, string) {
				path := filepath.Join(t.TempDir(), "agentd.yaml")
				require.NoError(t, os.WriteFile(path, []byte("version: 1\n"), 0o600))
				store, err := config.Load(context.Background(), path)
				require.NoError(t, err)
				return store, path
			},
			wantResult: "ok",
		},
		{
			name: "result error",
			setup: func(t *testing.T) (*config.Store, string) {
				path := filepath.Join(t.TempDir(), "agentd.yaml")
				require.NoError(t, os.WriteFile(path, []byte("version: 1\n"), 0o600))
				store, err := config.Load(context.Background(), path)
				require.NoError(t, err)
				require.NoError(t, os.WriteFile(path, []byte("not: [valid\n"), 0o600))
				return store, path
			},
			reloadErr:  true,
			wantResult: "error",
		},
		{
			name: "PatchRuntime does not increment",
			setup: func(t *testing.T) (*config.Store, string) {
				path := filepath.Join(t.TempDir(), "agentd.yaml")
				require.NoError(t, os.WriteFile(path, []byte("version: 1\n"), 0o600))
				store, err := config.Load(context.Background(), path)
				require.NoError(t, err)
				return store, path
			},
			wantResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store, _ := tt.setup(t)
			var results []string
			store.SetOnReload(func(result string) {
				results = append(results, result)
			})

			if tt.name == "PatchRuntime does not increment" {
				require.NoError(t, store.PatchRuntime([]byte("version: 1\n")))
				assert.Empty(t, results, "PatchRuntime must not call onReload")
				return
			}

			err := store.Reload(context.Background())
			if tt.reloadErr {
				require.Error(t, err, "Reload")
			} else {
				require.NoError(t, err, "Reload")
			}
			require.Len(t, results, 1, "onReload calls")
			assert.Equal(t, tt.wantResult, results[0], "result")
		})
	}
}
