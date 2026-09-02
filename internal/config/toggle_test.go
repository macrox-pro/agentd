package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/config"
)

func TestLookupToggle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		wantErr error
		wantDef bool
	}{
		{name: "unknown_toggle", in: "not-a-feature", wantErr: config.ErrUnknownToggle},
		{name: "list_toggle_names_sorted", in: "", wantDef: false},
		{name: "trajectory", in: "trajectory", wantDef: true},
		{name: "trajectory-raw", in: "trajectory-raw", wantDef: true},
		{name: "trajectory-statistics", in: "trajectory-statistics", wantDef: true},
		{name: "guard-shell", in: "guard-shell", wantDef: true},
		{name: "guard-mcp", in: "guard-mcp", wantDef: true},
		{name: "guard-paths", in: "guard-paths", wantDef: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.name == "list_toggle_names_sorted" {
				names := config.ListToggleNames()
				require.True(t, sort.StringsAreSorted(names), "ListToggleNames() sorted")
				assert.Len(t, names, 6, "ListToggleNames()")
				return
			}
			_, err := config.LookupToggle(tt.in)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr, "LookupToggle(%q)", tt.in)
				return
			}
			require.NoError(t, err, "LookupToggle(%q)", tt.in)
		})
	}
}

func TestSetToggle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T) (userPath, projectDir string)
		toggle     string
		scope      config.ToggleScope
		enabled    bool
		wantErr    error
		after      func(t *testing.T, userPath, projectDir string)
		runTwice   bool
		secondWant error
	}{
		{
			name: "enable_trajectory_creates_user_bootstrap",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return filepath.Join(dir, "user.yaml"), dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: false,
			after: func(t *testing.T, userPath, _ string) {
				raw, err := os.ReadFile(userPath)
				require.NoError(t, err, "ReadFile(%q)", userPath)
				body := string(raw)
				assert.Contains(t, body, "trajectory:", "user config")
				assert.Contains(t, body, "enabled: false", "trajectory disabled")
				assert.Contains(t, body, "fail_closed", "bootstrap policy")
			},
		},
		{
			name: "enable_trajectory_idempotent",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\ntrajectory:\n  enabled: true\n"), 0o600))
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: true,
			wantErr: config.ErrToggleAlreadySet,
		},
		{
			name: "trajectory-statistics",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\ntrajectory:\n  enabled: true\n"), 0o600))
				return userPath, dir
			},
			toggle:  "trajectory-statistics",
			scope:   config.ToggleScopeUser,
			enabled: false,
			after: func(t *testing.T, userPath, _ string) {
				raw, err := os.ReadFile(userPath)
				require.NoError(t, err, "ReadFile(%q)", userPath)
				assert.Contains(t, string(raw), "statistics: false", "trajectory statistics")
			},
		},
		{
			name: "disable_trajectory_explicit_false",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\ntrajectory:\n  enabled: true\n"), 0o600))
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: false,
			after: func(t *testing.T, userPath, _ string) {
				raw, err := os.ReadFile(userPath)
				require.NoError(t, err, "ReadFile(%q)", userPath)
				assert.Contains(t, string(raw), "enabled: false", "explicit false")
			},
		},
		{
			name: "enable_guard_shell_project_new_file",
			setup: func(t *testing.T) (string, string) {
				return "", t.TempDir()
			},
			toggle:  "guard-shell",
			scope:   config.ToggleScopeProject,
			enabled: true,
			after: func(t *testing.T, _, projectDir string) {
				raw, err := os.ReadFile(filepath.Join(projectDir, ".agentd.yaml"))
				require.NoError(t, err, "ReadFile project config")
				assert.Contains(t, string(raw), "shell:", "shell guard")
				assert.Contains(t, string(raw), "enabled: true", "shell enabled")
			},
		},
		{
			name: "enable_guard_shell_user_scope",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return filepath.Join(dir, "user.yaml"), dir
			},
			toggle:  "guard-shell",
			scope:   config.ToggleScopeUser,
			enabled: true,
			after: func(t *testing.T, userPath, _ string) {
				raw, err := os.ReadFile(userPath)
				require.NoError(t, err, "ReadFile(%q)", userPath)
				assert.Contains(t, string(raw), "shell:", "shell guard")
			},
		},
		{
			name: "invalid_yaml_blocks_write",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(userPath, []byte(":\n  bad:\n"), 0o600))
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: true,
			wantErr: config.ErrParseConfig,
			after: func(t *testing.T, userPath, _ string) {
				raw, err := os.ReadFile(userPath)
				require.NoError(t, err, "ReadFile(%q)", userPath)
				assert.Equal(t, ":\n  bad:\n", string(raw), "unchanged corrupt file")
			},
		},
		{
			name: "enable_each_catalog_entry",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return filepath.Join(dir, "user.yaml"), dir
			},
			toggle:  "guard-paths",
			scope:   config.ToggleScopeProject,
			enabled: true,
			after: func(t *testing.T, userPath, projectDir string) {
				for _, name := range config.ListToggleNames() {
					scope := config.ToggleScopeUser
					dir := projectDir
					enabled := true
					if strings.HasPrefix(name, "guard-") {
						scope = config.ToggleScopeProject
						dir = filepath.Join(projectDir, name)
						require.NoError(t, os.MkdirAll(dir, 0o700))
					}
					if name == "trajectory" || name == "trajectory-raw" || name == "trajectory-statistics" {
						enabled = false
					}
					_, err := config.SetToggle(config.SetToggleOptions{
						Name: name, Scope: scope, Enabled: enabled,
						UserPath: userPath, ProjectDir: dir,
					})
					require.NoError(t, err, "SetToggle(%q)", name)
				}
			},
		},
		{
			name: "disable_effective_off_idempotent",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				_, err := config.SetToggle(config.SetToggleOptions{
					Name: "trajectory", Scope: config.ToggleScopeUser, Enabled: false,
					UserPath: userPath, ProjectDir: dir,
				})
				require.NoError(t, err, "SetToggle(trajectory off)")
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: false,
			wantErr: config.ErrToggleAlreadySet,
		},
		{
			name: "user_path_is_directory",
			setup: func(t *testing.T) (string, string) {
				return t.TempDir(), t.TempDir()
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: true,
		},
		{
			name: "empty_user_path",
			setup: func(t *testing.T) (string, string) {
				return "", t.TempDir()
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: true,
		},
		{
			name: "invalid_compile_blocks_write",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: nope\n"), 0o600))
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: true,
			after: func(t *testing.T, userPath, _ string) {
				raw, err := os.ReadFile(userPath)
				require.NoError(t, err, "ReadFile(%q)", userPath)
				assert.Contains(t, string(raw), "fail: nope", "unchanged invalid compile")
			},
		},
		{
			name: "invalid_parse_blocks_write",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy: ["), 0o600))
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: true,
			wantErr: config.ErrParseConfig,
		},
		{
			name: "empty_user_file_enable",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(userPath, nil, 0o600))
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: false,
			after: func(t *testing.T, userPath, _ string) {
				_, err := config.Load(t.Context(), userPath)
				require.NoError(t, err, "Load(%q)", userPath)
			},
		},
		{
			name: "sequential_enables_preserve_keys",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return filepath.Join(dir, "user.yaml"), dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: false,
			after: func(t *testing.T, userPath, projectDir string) {
				_, err := config.SetToggle(config.SetToggleOptions{
					Name: "trajectory-raw", Scope: config.ToggleScopeUser, Enabled: false,
					UserPath: userPath, ProjectDir: projectDir,
				})
				require.NoError(t, err, "SetToggle(trajectory-raw)")
				raw, err := os.ReadFile(userPath)
				require.NoError(t, err, "ReadFile(%q)", userPath)
				body := string(raw)
				assert.Contains(t, body, "enabled: false", "trajectory disabled")
				assert.Contains(t, body, "include_raw: false", "trajectory-raw disabled")
			},
		},
		{
			name: "enable_preserves_sibling_keys",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				content := "version: 1\npolicy:\n  fail: fail_open\n"
				require.NoError(t, os.WriteFile(userPath, []byte(content), 0o600))
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: false,
			after: func(t *testing.T, userPath, _ string) {
				raw, err := os.ReadFile(userPath)
				require.NoError(t, err, "ReadFile(%q)", userPath)
				assert.Contains(t, string(raw), "fail_open", "policy preserved")
			},
		},
		{
			name: "trajectory_raw_without_trajectory_ok",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return filepath.Join(dir, "user.yaml"), dir
			},
			toggle:  "trajectory-raw",
			scope:   config.ToggleScopeUser,
			enabled: false,
			after: func(t *testing.T, userPath, _ string) {
				_, err := config.Load(t.Context(), userPath)
				require.NoError(t, err, "Load(%q)", userPath)
			},
		},
		{
			name: "trajectory_project_scope_write",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return filepath.Join(dir, "user.yaml"), dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeProject,
			enabled: false,
			after: func(t *testing.T, _, projectDir string) {
				raw, err := os.ReadFile(filepath.Join(projectDir, ".agentd.yaml"))
				require.NoError(t, err, "ReadFile project")
				assert.Contains(t, string(raw), "trajectory:", "project trajectory")
			},
		},
		{
			name: "set_toggle_no_runtime_io",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				runtimePath := filepath.Join(dir, "runtime.yaml")
				runtimeBefore := "version: 1\ntrajectory:\n  enabled: true\n"
				require.NoError(t, os.WriteFile(runtimePath, []byte(runtimeBefore), 0o600))
				t.Cleanup(func() {
					raw, err := os.ReadFile(runtimePath)
					if err == nil {
						assert.Equal(t, runtimeBefore, string(raw), "runtime unchanged")
					}
				})
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: false,
			after: func(t *testing.T, userPath, _ string) {
				_, err := os.Stat(userPath)
				require.NoError(t, err, "Stat(%q)", userPath)
			},
		},
		{
			name: "project_cwd_created_on_write",
			setup: func(t *testing.T) (string, string) {
				base := t.TempDir()
				projectDir := filepath.Join(base, "nested", "repo")
				return filepath.Join(base, "user.yaml"), projectDir
			},
			toggle:  "guard-shell",
			scope:   config.ToggleScopeProject,
			enabled: true,
			after: func(t *testing.T, _, projectDir string) {
				_, err := os.Stat(filepath.Join(projectDir, ".agentd.yaml"))
				require.NoError(t, err, "Stat project config")
			},
		},
		{
			name: "write_blocked_parent_fails",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				block := filepath.Join(dir, "block")
				require.NoError(t, os.WriteFile(block, []byte("x"), 0o600), "WriteFile(%q)", block)
				return filepath.Join(block, "user.yaml"), dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: false,
		},
		{
			name: "enable_effective_on_idempotent",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(filepath.Join(dir, ".agentd.yaml"), []byte("version: 1\ntrajectory:\n  enabled: true\n"), 0o600))
				return userPath, dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: true,
			wantErr: config.ErrToggleAlreadySet,
			after: func(t *testing.T, userPath, _ string) {
				_, err := os.Stat(userPath)
				assert.True(t, os.IsNotExist(err), "Stat(%q)", userPath)
			},
		},
		{
			name: "written_yaml_no_null_keys",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return filepath.Join(dir, "user.yaml"), dir
			},
			toggle:  "trajectory",
			scope:   config.ToggleScopeUser,
			enabled: false,
			after: func(t *testing.T, userPath, _ string) {
				raw, err := os.ReadFile(userPath)
				require.NoError(t, err, "ReadFile(%q)", userPath)
				assert.NotContains(t, string(raw), "null:", "yaml output")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userPath, projectDir := tt.setup(t)
			_, err := config.SetToggle(config.SetToggleOptions{
				Name:       tt.toggle,
				Scope:      tt.scope,
				Enabled:    tt.enabled,
				UserPath:   userPath,
				ProjectDir: projectDir,
			})
			if tt.wantErr != nil {
				require.Error(t, err, "SetToggle(%q)", tt.name)
				if errors.Is(tt.wantErr, config.ErrParseConfig) {
					require.ErrorIs(t, err, config.ErrParseConfig, "SetToggle(%q)", tt.name)
				} else {
					require.ErrorIs(t, err, tt.wantErr, "SetToggle(%q)", tt.name)
				}
			} else if tt.name == "user_path_is_directory" || tt.name == "empty_user_path" || tt.name == "write_blocked_parent_fails" || tt.name == "invalid_compile_blocks_write" {
				require.Error(t, err, "SetToggle(%q)", tt.name)
			} else {
				require.NoError(t, err, "SetToggle(%q)", tt.name)
			}
			if tt.after != nil {
				tt.after(t, userPath, projectDir)
			}
		})
	}
}

func TestGetToggle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T) (userPath, projectDir string)
		toggle     string
		wantOn     bool
		wantSource config.ToggleSource
	}{
		{
			name: "trajectory-statistics",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				_, err := config.SetToggle(config.SetToggleOptions{
					Name: "trajectory-statistics", Scope: config.ToggleScopeUser, Enabled: false,
					UserPath: userPath, ProjectDir: dir,
				})
				require.NoError(t, err, "SetToggle(trajectory-statistics off)")
				return userPath, dir
			},
			toggle:     "trajectory-statistics",
			wantOn:     false,
			wantSource: config.ToggleSourceUser,
		},
		{
			name: "get_shows_default_on",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return filepath.Join(dir, "missing.yaml"), dir
			},
			toggle:     "trajectory",
			wantOn:     true,
			wantSource: config.ToggleSourceDefault,
		},
		{
			name: "get_shows_user_layer",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				_, err := config.SetToggle(config.SetToggleOptions{
					Name: "trajectory", Scope: config.ToggleScopeUser, Enabled: false,
					UserPath: userPath, ProjectDir: dir,
				})
				require.NoError(t, err, "SetToggle(trajectory off)")
				return userPath, dir
			},
			toggle:     "trajectory",
			wantOn:     false,
			wantSource: config.ToggleSourceUser,
		},
		{
			name: "get_project_overrides_user",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(userPath, []byte("version: 1\npolicy:\n  fail: fail_closed\ntrajectory:\n  enabled: true\n"), 0o600))
				require.NoError(t, os.WriteFile(filepath.Join(dir, ".agentd.yaml"), []byte("version: 1\ntrajectory:\n  enabled: false\n"), 0o600))
				return userPath, dir
			},
			toggle:     "trajectory",
			wantOn:     false,
			wantSource: config.ToggleSourceProject,
		},
		{
			name: "get_shows_project_layer_only",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				userPath := filepath.Join(dir, "user.yaml")
				require.NoError(t, os.WriteFile(filepath.Join(dir, ".agentd.yaml"), []byte("version: 1\ntrajectory:\n  enabled: false\n"), 0o600))
				return userPath, dir
			},
			toggle:     "trajectory",
			wantOn:     false,
			wantSource: config.ToggleSourceProject,
		},
		{
			name: "get_ignores_runtime_overlay",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				stateDir := filepath.Join(dir, "state", "agentd")
				require.NoError(t, os.MkdirAll(stateDir, 0o700))
				require.NoError(t, os.WriteFile(filepath.Join(stateDir, "runtime.yaml"), []byte("version: 1\ntrajectory:\n  enabled: true\n"), 0o600))
				userPath := filepath.Join(dir, "user.yaml")
				return userPath, dir
			},
			toggle:     "trajectory",
			wantOn:     true,
			wantSource: config.ToggleSourceDefault,
		},
		{
			name: "get_finds_ancestor_project",
			setup: func(t *testing.T) (string, string) {
				base := t.TempDir()
				root := filepath.Join(base, "repo")
				nested := filepath.Join(root, "pkg", "sub")
				require.NoError(t, os.MkdirAll(nested, 0o700))
				require.NoError(t, os.WriteFile(filepath.Join(root, ".agentd.yaml"), []byte("version: 1\ntrajectory:\n  enabled: false\n"), 0o600))
				return filepath.Join(base, "user.yaml"), nested
			},
			toggle:     "trajectory",
			wantOn:     false,
			wantSource: config.ToggleSourceProject,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userPath, projectDir := tt.setup(t)
			got, err := config.GetToggle(config.GetToggleOptions{
				Name:       tt.toggle,
				UserPath:   userPath,
				ProjectDir: projectDir,
			})
			require.NoError(t, err, "GetToggle(%q)", tt.name)
			assert.Equal(t, tt.wantOn, got.Enabled, "GetToggle(%q) enabled", tt.name)
			assert.Equal(t, tt.wantSource, got.Source, "GetToggle(%q) source", tt.name)
		})
	}
}
