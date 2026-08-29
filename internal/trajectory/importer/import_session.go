package importer

import (
	"context"
	"fmt"
	"os"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/provider"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

// ImportSessionOptions configures an offline transcript import into the session ledger.
type ImportSessionOptions struct {
	Provider       provider.ID
	SessionID      string
	TranscriptPath string
	DryRun         bool
	EmitOnly       bool // skip ledger, checkpoint, hub; populate Events when true
	ConfigPath     string
	SessionsRoot   string
	SnapshotCfg    *config.TrajectoryConfig
	Hub            *trajectory.Hub
}

// ImportSessionResult summarizes one import run.
type ImportSessionResult struct {
	Provider       string             `json:"provider"`
	SessionID      string             `json:"session_id"`
	ImporterStatus string             `json:"importer_status"`
	Imported       int                `json:"imported"`
	TranscriptPath string             `json:"transcript_path,omitempty"`
	LastLineIndex  int                `json:"last_line_index,omitempty"`
	DryRun         bool               `json:"dry_run,omitempty"`
	Events         []trajectory.Event `json:"-"`
}

// ImportSession imports provider transcript JSONL into the local session ledger.
func ImportSession(ctx context.Context, opts ImportSessionOptions) (ImportSessionResult, error) {
	prov := opts.Provider
	provStr := string(prov)

	sessionsRoot := opts.SessionsRoot
	if sessionsRoot == "" {
		sessionsRoot = trajectory.DefaultSessionsDir()
	}
	if sessionsRoot == "" {
		return ImportSessionResult{}, trajectory.ErrSessionsDirUnavailable
	}

	cfg := trajectoryConfigForImport(ctx, opts)
	sid := opts.SessionID
	if sid == "" && opts.TranscriptPath != "" {
		sid = SessionIDFromTranscriptPath(opts.TranscriptPath)
	}

	startIndex := 0
	if sid != "" {
		cp, err := trajectory.LoadImportCheckpoint(trajectory.ImportSidecarPath(sessionsRoot, provStr, sid))
		if err != nil {
			return ImportSessionResult{}, fmt.Errorf("load import checkpoint: %w", err)
		}
		if cp.SourcePath != "" {
			startIndex = cp.LastLineIndex + 1
		}
	}

	result, err := Import(prov, ImportOptions{
		SessionID:      sid,
		TranscriptPath: opts.TranscriptPath,
		StartIndex:     startIndex,
		Cfg:            cfg,
	})
	if err != nil {
		return ImportSessionResult{}, err
	}
	if sid == "" {
		sid = SessionIDFromTranscriptPath(result.TranscriptPath)
	}

	status := ProviderImporterStatus(provStr)
	out := ImportSessionResult{
		Provider:       provStr,
		SessionID:      sid,
		ImporterStatus: string(status),
		Imported:       len(result.Events),
		TranscriptPath: result.TranscriptPath,
		LastLineIndex:  result.LastLineIndex,
		DryRun:         opts.DryRun,
	}

	if opts.EmitOnly {
		if len(result.Events) > 0 {
			key := trajectory.ResolveSessionKeyID(prov, sid, "", "")
			if err := trajectory.AssignImportedSeq(sessionsRoot, key, result.Events); err != nil {
				return ImportSessionResult{}, fmt.Errorf("assign imported seq: %w", err)
			}
			out.Events = result.Events
		}
		return out, nil
	}

	if opts.DryRun || len(result.Events) == 0 {
		return out, nil
	}

	key := trajectory.ResolveSessionKeyID(prov, sid, "", "")
	if err := trajectory.AppendImported(sessionsRoot, key, result.Events); err != nil {
		return ImportSessionResult{}, fmt.Errorf("append imported events: %w", err)
	}
	if opts.Hub != nil {
		opts.Hub.Publish(result.Events)
	}

	st, err := os.Stat(result.TranscriptPath)
	if err != nil {
		return ImportSessionResult{}, fmt.Errorf("stat transcript: %w", err)
	}
	cp := trajectory.ImportCheckpoint{
		LastLineIndex: result.LastLineIndex,
		SourcePath:    result.TranscriptPath,
		SourceModTime: st.ModTime().UTC(),
	}
	if err := trajectory.SaveImportCheckpoint(trajectory.ImportSidecarPath(sessionsRoot, provStr, sid), cp); err != nil {
		return ImportSessionResult{}, fmt.Errorf("save import checkpoint: %w", err)
	}

	return out, nil
}

func trajectoryConfigForImport(ctx context.Context, opts ImportSessionOptions) config.TrajectoryConfig {
	if opts.SnapshotCfg != nil {
		return *opts.SnapshotCfg
	}
	return importTrajectoryConfig(ctx, opts.ConfigPath)
}

func importTrajectoryConfig(ctx context.Context, path string) config.TrajectoryConfig {
	if path != "" {
		store, err := config.Load(ctx, path)
		if err == nil {
			return store.Current().Trajectory
		}
	}
	return trajectory.DefaultImportConfig()
}
