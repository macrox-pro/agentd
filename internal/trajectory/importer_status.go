package trajectory

import "github.com/macrox-pro/agentd/internal/provider"

// ImporterStatus is L2 import support for a provider (DESIGN §14.6).
type ImporterStatus string

const (
	ImporterSupported ImporterStatus = "supported"
	ImporterPartial   ImporterStatus = "partial"
	ImporterNone      ImporterStatus = "none"
)

// ProviderImporterStatus returns the importer tier for a canonical provider id.
func ProviderImporterStatus(name string) ImporterStatus {
	id, ok := provider.Lookup(name)
	if !ok {
		return ImporterNone
	}
	switch id {
	case provider.ClaudeCode, provider.Codex:
		return ImporterSupported
	case provider.Cursor:
		return ImporterPartial
	default:
		return ImporterNone
	}
}
