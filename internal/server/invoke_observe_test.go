package server_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/server"
	"github.com/macrox-pro/agentd/internal/trajectory/statistics"
)

func TestInvoke_skipsObserveOnError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "agentd.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`version: 1
trajectory:
  enabled: true
  statistics: true
`), 0o600))
	store, err := config.Load(ctx, path)
	require.NoError(t, err)
	collector := statistics.NewCollector()
	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	srv := server.New(server.Options{
		Store:     store,
		Engine:    dispatch.NewEngine(q, nil),
		Collector: collector,
	})
	conn := dialBuf(t, srv)
	hook := agentdv1.NewHookServiceClient(conn)
	_, err = hook.Invoke(ctx, &agentdv1.InvokeRequest{
		Provider:   agentdv1.Provider_PROVIDER_CLAUDE_CODE,
		RawPayload: []byte(`not-json`),
	})
	require.NoError(t, err)
	assert.Empty(t, collector.Snapshot(agentdv1.Provider_PROVIDER_UNSPECIFIED).HooksByKind)
}
