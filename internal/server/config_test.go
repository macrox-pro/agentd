package server_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/server"
)

func TestConfigService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dir := t.TempDir()
	userPath := filepath.Join(dir, "user.yaml")
	require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\n"), 0o600))

	projDir := filepath.Join(dir, "proj")
	require.NoError(t, os.MkdirAll(projDir, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(projDir, ".agentd.yaml"), []byte("version: 1\npolicy:\n  fail: fail_open\n"), 0o600))

	store, err := config.Load(ctx, userPath)
	require.NoError(t, err)

	q := dispatch.NewQueue(store.Current().Async, nil)
	t.Cleanup(func() { q.Close(2 * time.Second) })
	eng := dispatch.NewEngine(q, nil)

	srv := server.New(server.Options{
		Store:     store,
		Engine:    eng,
		StartedAt: time.Now().UTC(),
		Version:   "test",
	})
	conn := dialBuf(t, srv)
	cfg := agentdv1.NewConfigServiceClient(conn)
	hook := agentdv1.NewHookServiceClient(conn)

	t.Run("get merged", func(t *testing.T) {
		resp, err := cfg.GetConfig(ctx, &agentdv1.GetConfigRequest{
			Layer: agentdv1.ConfigLayer_CONFIG_LAYER_MERGED,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.GetYamlContent())
		assert.GreaterOrEqual(t, resp.GetConfig().GetGeneration(), uint64(1))
	})

	t.Run("patch bumps generation", func(t *testing.T) {
		before := store.Current().Generation
		resp, err := cfg.PatchConfig(ctx, &agentdv1.PatchConfigRequest{
			YamlPatch: []byte("version: 1\nasync:\n  queue_capacity: 7\n"),
		})
		require.NoError(t, err)
		assert.Greater(t, resp.GetConfig().GetGeneration(), before)
		assert.Equal(t, 7, store.Current().Async.QueueCapacity)
	})

	t.Run("record decision unimplemented", func(t *testing.T) {
		_, err := cfg.RecordDecision(ctx, &agentdv1.RecordDecisionRequest{
			ApprovalFingerprint: "sha256:x",
		})
		require.Error(t, err)
		assert.Equal(t, codes.Unimplemented, status.Code(err))
	})

	t.Run("invoke with project cwd", func(t *testing.T) {
		resp, err := hook.Invoke(ctx, &agentdv1.InvokeRequest{
			Provider:       agentdv1.Provider_PROVIDER_CLAUDE_CODE,
			RawPayload:     claudeToolPre(t, "go test"),
			InvocationMode: agentdv1.InvocationMode_INVOCATION_MODE_STDIN,
			Cwd:            projDir,
		})
		require.NoError(t, err)
		assert.Equal(t, agentdv1.DecisionKind_DECISION_KIND_NO_DECISION, resp.GetDecision().GetKind())
		// project fail_open is in snap; fingerprint should differ from base if project loaded
		baseFP := store.Current().Fingerprint
		assert.NotEqual(t, baseFP, resp.GetConfig().GetFingerprint())
	})
}
