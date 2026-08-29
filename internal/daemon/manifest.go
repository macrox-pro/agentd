package daemon

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	launchdLabel      = "io.github.macrox-pro.agentd"
	systemdUnitName   = "agentd.service"
	windowsTaskName   = "agentd"
	manifestFilePerm  = 0o600
	manifestDirPerm   = 0o700
	manifestTmpSuffix = ".tmp"
)

// Backend names login autostart integrations for status JSON.
type Backend string

const (
	BackendSystemd     Backend = "systemd"
	BackendLaunchd     Backend = "launchd"
	BackendSchtasks    Backend = "schtasks"
	BackendUnsupported Backend = ""
)

// AutostartSpec holds the registered binary and service argv.
type AutostartSpec struct {
	Exe  string
	Args []string
}

// AutostartOptions configures Enable.
type AutostartOptions struct {
	Socket     string
	ConfigPath string
	Version    string
}

// AutostartReport is the login autostart snapshot for status JSON.
type AutostartReport struct {
	Enabled       bool
	Backend       Backend
	ManifestPath  string
	RegisteredExe string
	Stale         bool
}

func renderSystemdUnit(spec AutostartSpec) string {
	execParts := append([]string{spec.Exe}, spec.Args...)
	line := strings.Join(shellQuoteJoin(execParts), " ")
	return fmt.Sprintf(`[Unit]
Description=agentd user daemon

[Service]
Type=simple
ExecStart=%s
Restart=no

[Install]
WantedBy=default.target
`, line)
}

func renderLaunchdPlist(spec AutostartSpec) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	b.WriteString("\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">")
	b.WriteString("\n<plist version=\"1.0\">\n<dict>\n")
	b.WriteString("  <key>Label</key>\n  <string>")
	xml.EscapeText(&b, []byte(launchdLabel))
	b.WriteString("</string>\n")
	b.WriteString("  <key>RunAtLoad</key>\n  <true/>\n")
	b.WriteString("  <key>ProgramArguments</key>\n  <array>\n")
	for _, arg := range append([]string{spec.Exe}, spec.Args...) {
		b.WriteString("    <string>")
		xml.EscapeText(&b, []byte(arg))
		b.WriteString("</string>\n")
	}
	b.WriteString("  </array>\n</dict>\n</plist>\n")
	return b.String()
}

func schtasksTR(spec AutostartSpec) string {
	parts := append([]string{spec.Exe}, spec.Args...)
	if len(parts) == 1 {
		return quoteWindowsArg(parts[0])
	}
	exe := quoteWindowsArg(parts[0])
	rest := strings.Join(parts[1:], " ")
	return exe + " " + rest
}

func quoteWindowsArg(s string) string {
	if strings.ContainsAny(s, " \t\"") {
		return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return s
}

func shellQuoteJoin(parts []string) []string {
	out := make([]string, len(parts))
	for i, p := range parts {
		if strings.ContainsAny(p, " \t\"'\\") {
			out[i] = `"` + strings.ReplaceAll(p, `"`, `\"`) + `"`
		} else {
			out[i] = p
		}
	}
	return out
}

func parseSystemdExecStart(unit string) string {
	for line := range strings.SplitSeq(unit, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "ExecStart=") {
			continue
		}
		val := strings.TrimPrefix(line, "ExecStart=")
		return firstToken(val)
	}
	return ""
}

func parseLaunchdProgram(plist string) string {
	dec := xml.NewDecoder(strings.NewReader(plist))
	inArray := false
	for {
		tok, err := dec.Token()
		if err != nil {
			return ""
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "array" {
				inArray = true
			}
			if inArray && t.Name.Local == "string" {
				var s string
				if err := dec.DecodeElement(&s, &t); err != nil {
					return ""
				}
				return s
			}
		case xml.EndElement:
			if t.Name.Local == "array" {
				return ""
			}
		}
	}
}

func parseSchtasksQuery(xmlRaw string) string {
	// Task Scheduler XML: <Command>...</Command> or nested in Actions
	_, after, ok := strings.Cut(xmlRaw, "<Command>")
	if !ok {
		return ""
	}
	rest := after
	before, _, ok := strings.Cut(rest, "</Command>")
	if !ok {
		return ""
	}
	return strings.TrimSpace(before)
}

func firstToken(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if s[0] == '"' {
		end := strings.Index(s[1:], `"`)
		if end >= 0 {
			return s[1 : 1+end]
		}
	}
	if i := strings.IndexAny(s, " \t"); i >= 0 {
		return s[:i]
	}
	return s
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, manifestDirPerm); err != nil {
		return fmt.Errorf("mkdir autostart dir: %w", err)
	}
	tmp := path + manifestTmpSuffix
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return fmt.Errorf("write autostart tmp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename autostart file: %w", err)
	}
	return nil
}
