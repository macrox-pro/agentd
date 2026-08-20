package targets

import (
	"context"
	"log/slog"
	"strings"
)

// Log writes structured async audit lines.
type Log struct {
	Logger *slog.Logger
}

// InvokeAsync logs event metadata.
func (t *Log) InvokeAsync(ctx context.Context, req AsyncRequest) error {
	log := t.Logger
	if log == nil {
		log = slog.Default()
	}
	level := parseLogLevel(req.Target.Level)
	attrs := []any{
		"target", "log",
		"provider", req.Provider,
		"kind", req.EventKind,
	}
	if req.SyncOutcome != nil {
		attrs = append(attrs, "decision", req.SyncOutcome.Kind.String(), "reason", req.SyncOutcome.Reason)
	}
	log.Log(ctx, level, "dispatch async", attrs...)
	return nil
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
