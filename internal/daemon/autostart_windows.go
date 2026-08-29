//go:build windows

package daemon

import (
	"fmt"
	"os/exec"
	"strings"
)

func schtasksCreateArgs(spec AutostartSpec) []string {
	tr := schtasksTR(spec)
	return []string{
		"/Create",
		"/TN", windowsTaskName,
		"/SC", "ONLOGON",
		"/TR", tr,
		"/F",
	}
}

func registerAutostart(spec AutostartSpec) error {
	if out, err := exec.Command("schtasks", schtasksCreateArgs(spec)...).CombinedOutput(); err != nil {
		return fmt.Errorf("schtasks create: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func unregisterAutostart() error {
	out, err := exec.Command("schtasks", "/Delete", "/TN", windowsTaskName, "/F").CombinedOutput()
	if err != nil {
		msg := strings.ToLower(string(out))
		if strings.Contains(msg, "cannot find") || strings.Contains(msg, "does not exist") {
			return nil
		}
		return fmt.Errorf("schtasks delete: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func readAutostartState() (AutostartReport, error) {
	out, err := exec.Command("schtasks", "/Query", "/TN", windowsTaskName, "/XML").CombinedOutput()
	if err != nil {
		msg := strings.ToLower(string(out))
		if strings.Contains(msg, "cannot find") || strings.Contains(msg, "does not exist") {
			return AutostartReport{Backend: BackendSchtasks}, nil
		}
		return AutostartReport{}, fmt.Errorf("schtasks query: %w: %s", err, strings.TrimSpace(string(out)))
	}
	exe := parseSchtasksQuery(string(out))
	return AutostartReport{
		Enabled:       true,
		Backend:       BackendSchtasks,
		RegisteredExe: exe,
	}, nil
}
