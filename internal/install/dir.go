package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ahinstall "github.com/speakeasy-api/agenthooks/install"

	"github.com/macrox-pro/agentd/internal/provider"
)

const (
	codexHomeEnv = "CODEX_HOME"
	kimiHomeEnv  = "KIMI_CODE_HOME"
)

// ResolveDir picks the filesystem root for hook config install.
func ResolveDir(id provider.ID, scope ahinstall.Scope, dirFlag string, dirFlagSet bool, cwd, home string, getenv func(string) string) (string, error) {
	if dirFlagSet {
		if strings.TrimSpace(dirFlag) == "" {
			return "", ErrDirRequired
		}
		return absDir(dirFlag)
	}

	switch scope {
	case ahinstall.ScopePlugin:
		return "", ErrDirRequired
	case ahinstall.ScopeUser:
		root, err := userRoot(id, home, getenv)
		if err != nil {
			return "", err
		}
		return absDir(root)
	default:
		root, err := projectRoot(id, cwd)
		if err != nil {
			return "", err
		}
		return absDir(root)
	}
}

func userRoot(id provider.ID, home string, getenv func(string) string) (string, error) {
	if home == "" {
		return "", ErrHomeRequired
	}
	switch id {
	case provider.OpenCode:
		return "", ErrDirRequired
	case provider.Codex:
		if v := strings.TrimSpace(getenv(codexHomeEnv)); v != "" {
			return v, nil
		}
		return filepath.Join(home, ".codex"), nil
	case provider.KimiCode:
		if v := strings.TrimSpace(getenv(kimiHomeEnv)); v != "" {
			return v, nil
		}
		return filepath.Join(home, ".kimi-code"), nil
	case provider.Cursor:
		return filepath.Join(home, ".cursor"), nil
	case provider.ClaudeCode:
		return filepath.Join(home, ".claude"), nil
	case provider.Gemini:
		return filepath.Join(home, ".gemini"), nil
	default:
		return "", fmt.Errorf("unknown provider %q", id)
	}
}

func projectRoot(id provider.ID, cwd string) (string, error) {
	if cwd == "" {
		return "", fmt.Errorf("working directory is required")
	}
	if id == provider.Codex {
		return filepath.Join(cwd, ".codex"), nil
	}
	return cwd, nil
}

func absDir(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("abs dir: %w", err)
	}
	return abs, nil
}

func resolveWorkingDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}
	return cwd, nil
}

func resolveHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return home, nil
}
