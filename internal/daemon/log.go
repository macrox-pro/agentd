package daemon

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/macrox-pro/agentd/internal/config"
)

// SetupLogOptions configures daemon operational logging.
type SetupLogOptions struct {
	Logging    config.LoggingConfig
	Foreground bool
	LogLevel   string
	LogFile    string
}

// SetupLog configures slog for the daemon process: append to a log file and
// optionally mirror to stderr in foreground mode. Returns cleanup to close the
// file and restore the previous default logger.
func SetupLog(opts SetupLogOptions) (*slog.Logger, func(), error) {
	level, err := opts.Logging.EffectiveLevel(opts.LogLevel)
	if err != nil {
		return nil, nil, fmt.Errorf("log level: %w", err)
	}
	path, err := resolveLogPath(opts)
	if err != nil {
		return nil, nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, nil, fmt.Errorf("log dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}

	var w io.Writer = f
	if opts.Foreground {
		w = io.MultiWriter(f, os.Stderr)
	}

	prev := slog.Default()
	h := slog.NewTextHandler(w, &slog.HandlerOptions{Level: level})
	log := slog.New(h)
	slog.SetDefault(log)

	cleanup := func() {
		_ = f.Close()
		slog.SetDefault(prev)
	}
	return log, cleanup, nil
}

func resolveLogPath(opts SetupLogOptions) (string, error) {
	path := opts.Logging.EffectiveFile(opts.LogFile)
	if path == "" {
		path = config.DefaultLogPath()
	}
	if path == "" {
		return "", fmt.Errorf("log file: state directory unavailable")
	}
	return path, nil
}
