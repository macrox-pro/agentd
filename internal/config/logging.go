package config

import (
	"fmt"
	"log/slog"
)

// LogLevel is the daemon operational log verbosity.
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LoggingConfig is compiled daemon operational logging settings.
type LoggingConfig struct {
	Level LogLevel
	File  string // empty = default state-dir agentd.log
}

func defaultLogging() LoggingConfig {
	return LoggingConfig{
		Level: LogLevelInfo,
		File:  "",
	}
}

type fileLogging struct {
	Level string `yaml:"level,omitempty"`
	File  string `yaml:"file,omitempty"`
}

func parseLogging(in *fileLogging, base LoggingConfig) (LoggingConfig, error) {
	out := base
	if in == nil {
		return out, nil
	}
	if in.Level != "" {
		lvl, err := parseLogLevel(in.Level)
		if err != nil {
			return LoggingConfig{}, err
		}
		out.Level = lvl
	}
	if in.File != "" {
		out.File = in.File
	}
	return out, nil
}

func parseLogLevel(s string) (LogLevel, error) {
	switch LogLevel(s) {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return LogLevel(s), nil
	default:
		return "", fmt.Errorf("logging.level unknown %q", s)
	}
}

// EffectiveLevel returns the slog level, applying a non-empty CLI override first.
func (c LoggingConfig) EffectiveLevel(override string) (slog.Level, error) {
	lvl := c.Level
	if override != "" {
		parsed, err := parseLogLevel(override)
		if err != nil {
			return slog.LevelInfo, err
		}
		lvl = parsed
	}
	return logLevelToSlog(lvl), nil
}

// EffectiveFile returns the configured log file path, or "" when the default
// state-dir path should be used. CLI override wins when non-empty.
func (c LoggingConfig) EffectiveFile(override string) string {
	if override != "" {
		return override
	}
	return c.File
}

func logLevelToSlog(l LogLevel) slog.Level {
	switch l {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func mergeLoggingPtr(base, user *fileLogging) *fileLogging {
	if base == nil && user == nil {
		return nil
	}
	out := fileLogging{}
	if base != nil {
		out = *base
	}
	if user == nil {
		return &out
	}
	if user.Level != "" {
		out.Level = user.Level
	}
	if user.File != "" {
		out.File = user.File
	}
	return &out
}
