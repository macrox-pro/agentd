package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestApprovalFingerprintStable(t *testing.T) {
	t.Parallel()
	a := config.ApprovalFingerprint(config.ApprovalKindSecrets, "Bash", "aws_key")
	b := config.ApprovalFingerprint(config.ApprovalKindSecrets, "Bash", "aws_key")
	c := config.ApprovalFingerprint(config.ApprovalKindSecrets, "Bash", "github_pat")
	assert.Equal(t, a, b)
	assert.NotEqual(t, a, c)
	assert.True(t, len(a) > len("sha256:"))
	assert.Contains(t, a, "sha256:secrets/")
}

func TestSecretsStableKey(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "aws_key,github_pat", config.SecretsStableKey([]string{"github_pat", "aws_key"}))
	assert.Equal(t, "", config.SecretsStableKey(nil))
}

func TestCompileMergedApprovalsDropsExpired(t *testing.T) {
	t.Parallel()
	past := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)
	future := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	fpLive := config.ApprovalFingerprint(config.ApprovalKindSecrets, "Bash", "aws_key")
	fpDead := config.ApprovalFingerprint(config.ApprovalKindSecrets, "Bash", "jwt")

	runtimeYAML := []byte(`
version: 1
approvals:
  secrets:
    - fingerprint: "` + fpDead + `"
      scope: project
      project: /repo
      expires_at: "` + past + `"
      granted_by: ask_user
    - fingerprint: "` + fpLive + `"
      scope: project
      project: /repo
      expires_at: "` + future + `"
      granted_by: ask_user
`)
	store, err := config.LoadWith(t.Context(), config.LoadOptions{})
	require.NoError(t, err)
	require.NoError(t, store.PatchRuntime(runtimeYAML))

	snap := store.Current()
	require.Len(t, snap.Approvals.Secrets, 1)
	assert.Equal(t, fpLive, snap.Approvals.Secrets[0].Fingerprint)
	assert.True(t, snap.Approvals.HasApproval(config.ApprovalKindSecrets, fpLive, "/repo", "", time.Now().UTC()))
	assert.False(t, snap.Approvals.HasApproval(config.ApprovalKindSecrets, fpDead, "/repo", "", time.Now().UTC()))
}

func TestMergeApprovalsUpsertByFingerprint(t *testing.T) {
	t.Parallel()
	fp := config.ApprovalFingerprint(config.ApprovalKindShell, "Bash", "curl")
	base := []byte(`
version: 1
approvals:
  shell:
    - fingerprint: "` + fp + `"
      scope: session
      session_id: s1
      granted_by: ask_user
`)
	overlay := []byte(`
version: 1
approvals:
  shell:
    - fingerprint: "` + fp + `"
      scope: session
      session_id: s2
      granted_by: ask_user
`)
	store, err := config.LoadWith(t.Context(), config.LoadOptions{})
	require.NoError(t, err)
	require.NoError(t, store.PatchRuntime(base))
	require.NoError(t, store.PatchRuntime(overlay))
	snap := store.Current()
	require.Len(t, snap.Approvals.Shell, 1)
	assert.Equal(t, "s2", snap.Approvals.Shell[0].SessionID)
}

func TestCompileMergedTemporaryBlocks(t *testing.T) {
	t.Parallel()
	past := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)
	future := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	runtimeYAML := []byte(`
version: 1
blocks:
  temporary:
    - tool: shell
      pattern: "curl * | sh"
      reason: expired
      until: "` + past + `"
    - tool: Bash
      pattern: "dangerous"
      reason: live
      until: "` + future + `"
`)
	store, err := config.LoadWith(t.Context(), config.LoadOptions{})
	require.NoError(t, err)
	require.NoError(t, store.PatchRuntime(runtimeYAML))
	snap := store.Current()
	require.Len(t, snap.TemporaryBlocks, 1)
	assert.Equal(t, "Bash", snap.TemporaryBlocks[0].Tool)
	hit := config.MatchTemporaryBlock(snap.TemporaryBlocks, "Bash", "do dangerous things", time.Now().UTC())
	require.NotNil(t, hit)
	assert.Equal(t, "live", hit.Reason)
	assert.Nil(t, config.MatchTemporaryBlock(snap.TemporaryBlocks, "Bash", "safe", time.Now().UTC()))
}
