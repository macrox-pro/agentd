package importer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/trajectory/importer"
)

func TestProviderImporterStatusTable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		provider string
		want     importer.ImporterStatus
	}{
		{name: "claude-code", provider: "claude-code", want: importer.ImporterSupported},
		{name: "cursor", provider: "cursor", want: importer.ImporterPartial},
		{name: "codex", provider: "codex", want: importer.ImporterSupported},
		{name: "gemini", provider: "gemini", want: importer.ImporterNone},
		{name: "opencode", provider: "opencode", want: importer.ImporterNone},
		{name: "kimi-code", provider: "kimi-code", want: importer.ImporterNone},
		{name: "kimi alias", provider: "kimicode", want: importer.ImporterNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, importer.ProviderImporterStatus(tt.provider))
		})
	}
}
