package trajectory

// ImporterStatus is L2 import support for a provider (DESIGN §14.6).
type ImporterStatus string

const (
	ImporterSupported ImporterStatus = "supported"
	ImporterPartial   ImporterStatus = "partial"
	ImporterNone      ImporterStatus = "none"
)

// ProviderImporterStatus returns the importer tier for a canonical provider id.
func ProviderImporterStatus(provider string) ImporterStatus {
	switch CanonicalProvider(provider) {
	case "claude-code":
		return ImporterSupported
	case "cursor", "codex":
		return ImporterPartial
	default:
		return ImporterNone
	}
}
