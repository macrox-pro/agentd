//go:build darwin

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func launchdPlistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", fmt.Errorf("home directory unavailable")
	}
	return filepath.Join(home, "Library", "LaunchAgents", launchdLabel+".plist"), nil
}

func registerAutostart(spec AutostartSpec) error {
	path, err := launchdPlistPath()
	if err != nil {
		return err
	}
	body := renderLaunchdPlist(spec)
	if err := writeFileAtomic(path, []byte(body), manifestFilePerm); err != nil {
		return err
	}
	domain, err := launchdDomain()
	if err != nil {
		return err
	}
	_ = exec.Command("launchctl", "bootout", domain, path).Run()
	if out, err := exec.Command("launchctl", "bootstrap", domain, path).CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl bootstrap: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func unregisterAutostart() error {
	path, err := launchdPlistPath()
	if err != nil {
		return nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	if domain, err := launchdDomain(); err == nil {
		_ = exec.Command("launchctl", "bootout", domain, path).Run()
	}
	_ = os.Remove(path)
	return nil
}

func readAutostartState() (AutostartReport, error) {
	path, err := launchdPlistPath()
	if err != nil {
		return AutostartReport{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return AutostartReport{Backend: BackendLaunchd, ManifestPath: path}, nil
		}
		return AutostartReport{}, err
	}
	exe := parseLaunchdProgram(string(data))
	return AutostartReport{
		Enabled:       true,
		Backend:       BackendLaunchd,
		ManifestPath:  path,
		RegisteredExe: exe,
	}, nil
}

func launchdDomain() (string, error) {
	uid := strconv.Itoa(os.Getuid())
	return "gui/" + uid, nil
}
