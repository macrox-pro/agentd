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
	"google.golang.org/protobuf/types/known/timestamppb"

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

	t.Run("record decision project", func(t *testing.T) {
		before := store.Current().Generation
		fp := config.ApprovalFingerprint(config.ApprovalKindSecrets, "Bash", "aws_key")
		expires := time.Now().UTC().Add(2 * time.Hour)
		resp, err := cfg.RecordDecision(ctx, &agentdv1.RecordDecisionRequest{
			ApprovalFingerprint: fp,
			Scope:               agentdv1.ConfigLayer_CONFIG_LAYER_PROJECT,
			ProjectRoot:         projDir,
			ExpiresAt:           timestamppb.New(expires),
		})
		require.NoError(t, err)
		assert.Greater(t, resp.GetConfig().GetGeneration(), before)
		snap := store.Current()
		require.Len(t, snap.Approvals.Secrets, 1)
		assert.Equal(t, fp, snap.Approvals.Secrets[0].Fingerprint)
		assert.Equal(t, config.ApprovalScopeProject, snap.Approvals.Secrets[0].Scope)
		assert.Equal(t, projDir, snap.Approvals.Secrets[0].Project)
	})

	t.Run("record decision session requires session_id", func(t *testing.T) {
		fp := config.ApprovalFingerprint(config.ApprovalKindShell, "Bash", "curl")
		_, err := cfg.RecordDecision(ctx, &agentdv1.RecordDecisionRequest{
			ApprovalFingerprint: fp,
			Scope:               agentdv1.ConfigLayer_CONFIG_LAYER_RUNTIME,
		})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("record decision session", func(t *testing.T) {
		fp := config.ApprovalFingerprint(config.ApprovalKindShell, "Bash", "curl")
		resp, err := cfg.RecordDecision(ctx, &agentdv1.RecordDecisionRequest{
			ApprovalFingerprint: fp,
			Scope:               agentdv1.ConfigLayer_CONFIG_LAYER_RUNTIME,
			SessionId:           "sess-1",
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.GetConfig().GetGeneration())
		assert.True(t, store.Current().Approvals.HasApproval(
			config.ApprovalKindShell, fp, "", "sess-1", time.Now().UTC()))
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
