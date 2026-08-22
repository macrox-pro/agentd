package importer_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

func TestImportDispatchesByProvider(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		provider string
		path     string
		wantErr  bool
		wantProv string
	}{
		{
			name:     "claude-code",
			provider: "claude-code",
			path:     filepath.Join("testdata", "claude_session.jsonl"),
			wantProv: "claude-code",
		},
		{
			name:     "cursor",
			provider: "cursor",
			path:     filepath.Join("testdata", "cursor_transcript.jsonl"),
			wantProv: "cursor",
		},
		{
			name:     "codex",
			provider: "codex",
			path:     filepath.Join("testdata", "codex_transcript.jsonl"),
			wantProv: "codex",
		},
		{
			name:     "none",
			provider: "gemini",
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

func TestSessionIDFromTranscriptPath(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "s1", importer.SessionIDFromTranscriptPath("/tmp/s1.jsonl"))
	assert.Equal(t, "plain", importer.SessionIDFromTranscriptPath("plain"))
}
