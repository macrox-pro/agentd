package importer_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

func TestImportDispatchesByProvider(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		provider provider.ID
		path     string
		wantErr  bool
		wantProv string
	}{
		{
			name:     "claude-code",
			provider: provider.ClaudeCode,
			path:     filepath.Join("testdata", "claude_session.jsonl"),
			wantProv: "claude-code",
		},
		{
			name:     "cursor",
			provider: provider.Cursor,
			path:     filepath.Join("testdata", "cursor_transcript.jsonl"),
			wantProv: "cursor",
		},
		{
			name:     "codex",
			provider: provider.Codex,
			path:     filepath.Join("testdata", "codex_transcript.jsonl"),
			wantProv: "codex",
		},
		{
			name:     "none",
			provider: provider.Gemini,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := importer.Import(tt.provider, importer.ImportOptions{
				SessionID:      "dispatch-" + tt.name,
				TranscriptPath: tt.path,
				Cfg:            config.TrajectoryConfig{},
			})
			if tt.wantErr {
				require.Error(t, err, "Import(%q)", tt.provider)
				return
			}
			require.NoError(t, err, "Import(%q)", tt.provider)
			require.NotEmpty(t, result.Events)
			for _, e := range result.Events {
				assert.Equal(t, tt.wantProv, e.Provider)
			}
		})
	}
}

func TestImportSetsProjectsRootFromConfig(t *testing.T) {
	t.Parallel()
	claudeRoot := t.TempDir()
	sid := "cfg-claude"
	require.NoError(t, copyFile(
		filepath.Join("testdata", "claude_session.jsonl"),
		filepath.Join(claudeRoot, sid+".jsonl"),
	))

	tests := []struct {
		name     string
		provider provider.ID
		cfgPath  string
		session  string
		explicit string
	}{
		{
			name:     "claude-code",
			provider: provider.ClaudeCode,
			cfgPath:  claudeRoot,
			session:  sid,
		},
		{
			name:     "cursor",
			provider: provider.Cursor,
			cfgPath:  t.TempDir(),
			explicit: filepath.Join("testdata", "cursor_transcript.jsonl"),
		},
		{
			name:     "codex",
			provider: provider.Codex,
			cfgPath:  t.TempDir(),
			explicit: filepath.Join("testdata", "codex_transcript.jsonl"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := config.TrajectoryConfig{
				Import: map[string]config.ImportProviderConfig{
					string(tt.provider): {Path: tt.cfgPath},
				},
			}
			result, err := importer.Import(tt.provider, importer.ImportOptions{
				SessionID:      tt.session,
				TranscriptPath: tt.explicit,
				Cfg:            cfg,
			})
			require.NoError(t, err, "Import(%q)", tt.provider)
			require.NotEmpty(t, result.Events)
			if tt.explicit == "" {
				assert.Contains(t, result.TranscriptPath, tt.cfgPath)
			}
		})
	}
}

func TestSessionIDFromTranscriptPath(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "s1", importer.SessionIDFromTranscriptPath("/tmp/s1.jsonl"))
	assert.Equal(t, "plain", importer.SessionIDFromTranscriptPath("plain"))
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
