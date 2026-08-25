// Package importer maps provider on-disk transcripts into trajectory events.
//
// Owns: provider registry (dispatch + L2 status), transcript import (claude-code,
// cursor, codex rollout), mapping helpers, ImportSession facade.
// Must not: hook wire decode (hookedge), config compile (config), Invoke hot path.
//
// Invariants:
//   - Never invent thinking or tool output absent from transcript files.
//   - Codex thinking only from plaintext event_msg.agent_reasoning (not encrypted_content).
//   - Imported events use source=transcript and append-only seq assignment (caller).
//
// Entry: Import, ImportSession, ProviderImporterStatus, ImportClaude, ImportCursor, ImportCodex.
// See DESIGN.md §14.2, §14.3.
package importer
