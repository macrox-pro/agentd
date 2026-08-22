// Package importer maps provider on-disk transcripts into trajectory events.
//
// Owns: provider transcript import (claude-code, cursor, codex), mapping helpers.
// Must not: hook wire decode (hookedge), config compile (config), Invoke hot path.
//
// Invariants:
//   - Never invent thinking or tool output absent from transcript files.
//   - Imported events use source=transcript and append-only seq assignment (caller).
//
// Entry: ImportClaude, ImportCursor, ImportCodex.
// See DESIGN.md §14.3, §14.6.
package importer
