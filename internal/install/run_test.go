package install_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/install"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name          string
		provider      string
		scope         string
		workSubdir    string
		installSubdir string
		wantErr       bool
		wantPath      string
		wantSubstr    []string
		wantAbsent    []string
		checkShim     bool
		useDefaultDir bool
	}{
		{
			name:          "claude project",
			provider:      "claude-code",
			scope:         "project",
			wantPath:      filepath.Join(".claude", "settings.json"),
			wantSubstr:    []string{"agenthooks", "run", "--provider=claude-code", "SubagentStart", "PreCompact", "PostToolUseFailure"},
			useDefaultDir: true,
		},
		{
			name:          "cursor user defaults to home/.cursor",
			provider:      "cursor",
			scope:         "user",
			installSubdir: filepath.Join("home", ".cursor"),
			wantPath:      "hooks.json",
			wantSubstr:    []string{"agenthooks", "--provider=cursor", "subagentStart", "preCompact"},
			wantAbsent:    []string{"afterAgentThought", "afterFileEdit"},
			useDefaultDir: true,
		},
		{
			name:          "project cursor writes .cursor/hooks.json under cwd",
			provider:      "cursor",
			scope:         "project",
			wantPath:      filepath.Join(".cursor", "hooks.json"),
			wantSubstr:    []string{"agenthooks", "--provider=cursor", "subagentStart"},
			wantAbsent:    []string{"afterAgentThought", "afterFileEdit"},
			useDefaultDir: true,
		},
		{
			name:          "project codex defaults to cwd/.codex",
			provider:      "codex",
			scope:         "project",
			workSubdir:    "repo",
			installSubdir: filepath.Join("repo", ".codex"),
			wantPath:      "hooks.json",
			wantSubstr:    []string{"agenthooks", "--provider=codex"},
			useDefaultDir: true,
		},
		{
			name:          "gemini project skips unmapped natives",
			provider:      "gemini",
			scope:         "project",
			wantPath:      filepath.Join(".gemini", "settings.json"),
			wantSubstr:    []string{"agenthooks", "--provider=gemini"},
			wantAbsent:    []string{"subagentStart", "SubagentStart"},
			useDefaultDir: true,
		},
		{
			name:          "opencode shim",
			provider:      "opencode",
			scope:         "project",
			wantPath:      filepath.Join(".opencode", "plugin", "agenthooks.ts"),
			checkShim:     true,
			wantSubstr:    []string{"opencode"},
			useDefaultDir: true,
		},
		{
			name:          "kimi project still errors",
			provider:      "kimi-code",
			scope:         "project",
			wantErr:       true,
			useDefaultDir: true,
		},
		{
			name:     "unknown provider",
			provider: "nope",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			workDir := root
			if tt.workSubdir != "" {
				workDir = filepath.Join(root, tt.workSubdir)
				require.NoError(t, os.MkdirAll(workDir, 0o755))
			}
			homeDir := filepath.Join(root, "home")
			require.NoError(t, os.MkdirAll(homeDir, 0o755))
			t.Setenv("HOME", homeDir)

			cwd, err := os.Getwd()
			require.NoError(t, err)
			require.NoError(t, os.Chdir(workDir))
			t.Cleanup(func() { _ = os.Chdir(cwd) })

			bin := filepath.Join(root, "agentd")
			require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

			opts := install.Options{
				Provider: tt.provider,
				Scope:    tt.scope,
				Command:  []string{bin},
			}
			if !tt.useDefaultDir {
				opts.DirFlagSet = true
			}

			result, err := install.Run(context.Background(), opts)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			installRoot := workDir
			if tt.installSubdir != "" {
				installRoot = filepath.Join(root, tt.installSubdir)
			}

			path := filepath.Join(installRoot, tt.wantPath)
			b, err := os.ReadFile(path)
			require.NoError(t, err, "ReadFile(%q)", path)
			body := string(b)
			assert.Contains(t, body, bin)
			for _, s := range tt.wantSubstr {
				assert.Contains(t, body, s)
			}
			for _, s := range tt.wantAbsent {
				assert.NotContains(t, body, s)
			}
			if tt.checkShim {
				assert.True(t,
					strings.Contains(body, `"agenthooks", "serve"`) || strings.Contains(body, "agenthooks\", \"serve\""),
				)
			}

			var buf bytes.Buffer
			require.NoError(t, install.WriteReport(&buf, result))
			report := buf.String()
			assert.Contains(t, report, "provider="+tt.provider)
			assert.Contains(t, report, "scope="+tt.scope)
			assert.Contains(t, report, path)
		})
	}
}

func TestWriteReport_allUnchangedStillPrintsReport(t *testing.T) {
	t.Parallel()

	result := install.Result{
		Provider: "cursor",
		Scope:    "user",
		Dir:      "/home/me/.cursor",
		Changes: []install.FileChange{
			{Path: "hooks.json", State: install.StateUnchanged},
		},
	}
	var buf bytes.Buffer
	require.NoError(t, install.WriteReport(&buf, result))
	report := buf.String()
	assert.Contains(t, report, "provider=cursor scope=user dir=/home/me/.cursor")
	assert.Contains(t, report, "unchanged")
	assert.Contains(t, report, "/home/me/.cursor/hooks.json")
}

func TestRun_dryRunNoWrite(t *testing.T) {
	root := t.TempDir()
	workDir := filepath.Join(root, "repo")
	cursorDir := filepath.Join(workDir, ".cursor")
	require.NoError(t, os.MkdirAll(cursorDir, 0o755))
	homeDir := filepath.Join(root, "home")
	require.NoError(t, os.MkdirAll(homeDir, 0o755))
	t.Setenv("HOME", homeDir)

	bin := filepath.Join(root, "agentd")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	_, err := install.Run(context.Background(), install.Options{
		Provider:   "cursor",
		Scope:      "project",
		Dir:        workDir,
		DirFlagSet: true,
		Command:    []string{bin},
		DryRun:     true,
	})
	require.NoError(t, err, "Run(dry-run)")
	_, err = os.Stat(filepath.Join(cursorDir, "hooks.json"))
	assert.True(t, os.IsNotExist(err), "dry-run must not write hooks.json")
}
