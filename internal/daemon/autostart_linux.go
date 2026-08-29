//go:build linux

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func systemdUnitPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", fmt.Errorf("home directory unavailable")
	}
	return filepath.Join(home, ".config", "systemd", "user", systemdUnitName), nil
}

func registerAutostart(spec AutostartSpec) error {
	if err := linuxAutostartPreflight(); err != nil {
		return err
	}
	path, err := systemdUnitPath()
	if err != nil {
		return err
	}
	body := renderSystemdUnit(spec)
	if err := writeFileAtomic(path, []byte(body), manifestFilePerm); err != nil {
		return err
	}
	if out, err := exec.Command("systemctl", "--user", "daemon-reload").CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl daemon-reload: %w: %s", err, strings.TrimSpace(string(out)))
	}
	if out, err := exec.Command("systemctl", "--user", "enable", "--now", systemdUnitName).CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl enable: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func unregisterAutostart() error {
	if _, err := systemdUnitPath(); err != nil {
		return nil
	}
	path, err := systemdUnitPath()
	if err != nil {
		return nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	if out, err := exec.Command("systemctl", "--user", "disable", systemdUnitName).CombinedOutput(); err != nil {
		if !strings.Contains(string(out), "not loaded") && !strings.Contains(string(out), "No such file") {
			return fmt.Errorf("systemctl disable: %w: %s", err, strings.TrimSpace(string(out)))
		}
	}
	_ = os.Remove(path)
	_, _ = exec.Command("systemctl", "--user", "daemon-reload").CombinedOutput()
	return nil
}

func readAutostartState() (AutostartReport, error) {
	path, err := systemdUnitPath()
	if err != nil {
		return AutostartReport{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return AutostartReport{Backend: BackendSystemd, ManifestPath: path}, nil
		}
		return AutostartReport{}, err
	}
	exe := parseSystemdExecStart(string(data))
	return AutostartReport{
		Enabled:       true,
		Backend:       BackendSystemd,
		ManifestPath:  path,
		RegisteredExe: exe,
	}, nil
}

func linuxAutostartPreflight() error {
	if os.Getenv("XDG_RUNTIME_DIR") == "" {
		return ErrAutostartNotAvailable
	}
	if out, err := exec.Command("systemctl", "--user", "show-environment").CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", ErrAutostartNotAvailable, strings.TrimSpace(string(out)))
	}
	return nil
}
