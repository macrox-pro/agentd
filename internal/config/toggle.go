package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"go.yaml.in/yaml/v3"
)

// ToggleScope is the config layer written by SetToggle.
type ToggleScope string

const (
	ToggleScopeUser    ToggleScope = "user"
	ToggleScopeProject ToggleScope = "project"
)

// ToggleSource names which merge layer wins for GetToggle.
type ToggleSource string

const (
	ToggleSourceDefault ToggleSource = "default"
	ToggleSourceUser    ToggleSource = "user"
	ToggleSourceProject ToggleSource = "project"
)

type toggleKind int

const (
	toggleTrajectory toggleKind = iota
	toggleTrajectoryRaw
	toggleGuardShell
	toggleGuardMCP
	toggleGuardPaths
)

type toggleDef struct {
	name         string
	defaultScope ToggleScope
	kind         toggleKind
}

var toggleCatalog = []toggleDef{
	{name: "trajectory", defaultScope: ToggleScopeUser, kind: toggleTrajectory},
	{name: "trajectory-raw", defaultScope: ToggleScopeUser, kind: toggleTrajectoryRaw},
	{name: "guard-shell", defaultScope: ToggleScopeProject, kind: toggleGuardShell},
	{name: "guard-mcp", defaultScope: ToggleScopeProject, kind: toggleGuardMCP},
	{name: "guard-paths", defaultScope: ToggleScopeProject, kind: toggleGuardPaths},
}

// SetToggleOptions configures a persistent user or project layer write.
type SetToggleOptions struct {
	Name       string
	Scope      ToggleScope
	Enabled    bool
	UserPath   string
	ProjectDir string
}

// SetToggleResult describes the outcome of SetToggle.
type SetToggleResult struct {
	Name       string
	Enabled    bool
	Scope      ToggleScope
	ConfigPath string
	AlreadySet bool
}

// GetToggleOptions configures effective toggle inspection.
type GetToggleOptions struct {
	Name       string
	UserPath   string
	ProjectDir string
}

// GetToggleResult is the effective toggle state after merge (runtime excluded).
type GetToggleResult struct {
	Name    string
	Enabled bool
	Source  ToggleSource
}

// LookupToggle returns a catalog entry by CLI feature name.
func LookupToggle(name string) (toggleDef, error) {
	for _, t := range toggleCatalog {
		if t.name == name {
			return t, nil
		}
	}
	return toggleDef{}, ErrUnknownToggle
}

// ListToggleNames returns sorted curated feature names.
func ListToggleNames() []string {
	out := make([]string, len(toggleCatalog))
	for i, t := range toggleCatalog {
		out[i] = t.name
	}
	sort.Strings(out)
	return out
}

// SetToggle writes a bool toggle to the user or project config file.
func SetToggle(opts SetToggleOptions) (SetToggleResult, error) {
	def, err := LookupToggle(opts.Name)
	if err != nil {
		return SetToggleResult{}, err
	}
	scope := opts.Scope
	if scope == "" {
		scope = def.defaultScope
	}

	targetPath, err := toggleTargetPath(scope, opts.UserPath, opts.ProjectDir)
	if err != nil {
		return SetToggleResult{}, err
	}

	layerFC, fileExists, err := readLayerFC(targetPath, scope)
	if err != nil {
		return SetToggleResult{}, err
	}

	userFC, projectFC, err := loadPeerLayers(scope, opts.UserPath, opts.ProjectDir, targetPath, layerFC, fileExists)
	if err != nil {
		return SetToggleResult{}, err
	}

	if alreadySet(def, layerFC, userFC, projectFC, scope, targetPath, opts.Enabled) {
		return SetToggleResult{
			Name:       def.name,
			Enabled:    opts.Enabled,
			Scope:      scope,
			ConfigPath: targetPath,
			AlreadySet: true,
		}, ErrToggleAlreadySet
	}

	setLayerBool(layerFC, def.kind, opts.Enabled)

	if scope == ToggleScopeUser {
		userFC = layerFC
	} else {
		projectFC = layerFC
	}

	if _, err := CompileMerged(userFC, projectFC, nil); err != nil {
		return SetToggleResult{}, fmt.Errorf("compile config: %w", err)
	}

	raw, err := yaml.Marshal(layerFC)
	if err != nil {
		return SetToggleResult{}, fmt.Errorf("marshal config: %w", err)
	}
	if err := writeYAMLAtomic(targetPath, raw); err != nil {
		return SetToggleResult{}, fmt.Errorf("write config: %w", err)
	}

	return SetToggleResult{
		Name:       def.name,
		Enabled:    opts.Enabled,
		Scope:      scope,
		ConfigPath: targetPath,
	}, nil
}

// GetToggle returns the effective toggle state (defaults ⊕ user ⊕ project; no runtime).
func GetToggle(opts GetToggleOptions) (GetToggleResult, error) {
	def, err := LookupToggle(opts.Name)
	if err != nil {
		return GetToggleResult{}, err
	}

	userFC, _ := loadOptionalLayerFC(opts.UserPath)
	projectPath, _ := FindProjectConfig(opts.ProjectDir, "")
	var projectFC *fileConfig
	if projectPath != "" {
		projectFC, _ = loadOptionalLayerFC(projectPath)
	}

	defRes, _ := CompileMerged(nil, nil, nil)
	defVal := effectiveBool(defRes, def.kind)

	userVal := defVal
	if userFC != nil {
		res, err := CompileMerged(userFC, nil, nil)
		if err != nil {
			return GetToggleResult{}, fmt.Errorf("compile user config: %w", err)
		}
		userVal = effectiveBool(res, def.kind)
	}

	projectVal := userVal
	if projectFC != nil {
		res, err := CompileMerged(userFC, projectFC, nil)
		if err != nil {
			return GetToggleResult{}, fmt.Errorf("compile project config: %w", err)
		}
		projectVal = effectiveBool(res, def.kind)
	}

	source := ToggleSourceDefault
	enabled := defVal
	if userFC != nil && userVal != defVal {
		source = ToggleSourceUser
		enabled = userVal
	}
	if projectFC != nil && projectVal != userVal {
		source = ToggleSourceProject
		enabled = projectVal
	}

	return GetToggleResult{
		Name:    def.name,
		Enabled: enabled,
		Source:  source,
	}, nil
}

func toggleTargetPath(scope ToggleScope, userPath, projectDir string) (string, error) {
	switch scope {
	case ToggleScopeUser:
		if userPath == "" {
			return "", fmt.Errorf("user config path unavailable")
		}
		abs, err := filepath.Abs(userPath)
		if err != nil {
			return "", fmt.Errorf("user config path: %w", err)
		}
		return abs, nil
	case ToggleScopeProject:
		if projectDir == "" {
			return "", fmt.Errorf("project directory unavailable")
		}
		abs, err := filepath.Abs(projectDir)
		if err != nil {
			return "", fmt.Errorf("project directory: %w", err)
		}
		return filepath.Join(abs, projectConfigFileName), nil
	default:
		return "", fmt.Errorf("unknown scope %q", scope)
	}
}

func readLayerFC(path string, scope ToggleScope) (*fileConfig, bool, error) {
	st, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return newLayerTemplate(scope), false, nil
		}
		return nil, false, fmt.Errorf("stat config: %w", err)
	}
	if st.IsDir() {
		return nil, false, fmt.Errorf("read config: %q is a directory", path)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("read config: %w", err)
	}
	if len(raw) == 0 {
		return &fileConfig{Version: 1}, true, nil
	}
	var fc fileConfig
	if err := yaml.Unmarshal(raw, &fc); err != nil {
		return nil, false, fmt.Errorf("%w %q: %w", ErrParseConfig, path, err)
	}
	return &fc, true, nil
}

func newLayerTemplate(scope ToggleScope) *fileConfig {
	if scope == ToggleScopeUser {
		return bootstrapUserFileConfig()
	}
	return &fileConfig{Version: 1}
}

func loadOptionalLayerFC(path string) (*fileConfig, error) {
	if path == "" {
		return nil, nil
	}
	st, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("stat config: %w", err)
	}
	if st.IsDir() {
		return nil, fmt.Errorf("read config: %q is a directory", path)
	}
	_, fc, err := readFileConfig(path)
	return fc, err
}

func loadPeerLayers(scope ToggleScope, userPath, projectDir, targetPath string, layerFC *fileConfig, fileExists bool) (*fileConfig, *fileConfig, error) {
	var userFC, projectFC *fileConfig

	switch scope {
	case ToggleScopeUser:
		if fileExists {
			userFC = layerFC
		} else {
			userFC = cloneFileConfig(layerFC)
		}
		projectPath, ok := FindProjectConfig(projectDir, "")
		if ok {
			projectFC, err := loadOptionalLayerFC(projectPath)
			if err != nil {
				return nil, nil, err
			}
			return userFC, projectFC, nil
		}
	case ToggleScopeProject:
		projectFC = layerFC
		if userPath != "" {
			var err error
			userFC, err = loadOptionalLayerFC(userPath)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return userFC, projectFC, nil
}

func alreadySet(def toggleDef, layerFC, userFC, projectFC *fileConfig, scope ToggleScope, _ string, desired bool) bool {
	if v := layerBool(layerFC, def.kind); v != nil && *v == desired {
		return true
	}
	if layerBool(layerFC, def.kind) != nil {
		return false
	}

	var mergeUser, mergeProject *fileConfig
	switch scope {
	case ToggleScopeUser:
		mergeUser = userFC
		mergeProject = projectFC
	case ToggleScopeProject:
		mergeUser = userFC
		mergeProject = projectFC
	}
	res, err := CompileMerged(mergeUser, mergeProject, nil)
	if err != nil {
		return false
	}
	return effectiveBool(res, def.kind) == desired
}

func layerBool(fc *fileConfig, kind toggleKind) *bool {
	if fc == nil {
		return nil
	}
	switch kind {
	case toggleTrajectory:
		if fc.Trajectory != nil {
			return fc.Trajectory.Enabled
		}
	case toggleTrajectoryRaw:
		if fc.Trajectory != nil {
			return fc.Trajectory.IncludeRaw
		}
	case toggleGuardShell:
		if fc.Guards != nil && fc.Guards.Shell != nil {
			return fc.Guards.Shell.Enabled
		}
	case toggleGuardMCP:
		if fc.Guards != nil && fc.Guards.MCP != nil {
			return fc.Guards.MCP.Enabled
		}
	case toggleGuardPaths:
		if fc.Guards != nil && fc.Guards.Paths != nil {
			return fc.Guards.Paths.Enabled
		}
	}
	return nil
}

func setLayerBool(fc *fileConfig, kind toggleKind, v bool) {
	if fc == nil {
		return
	}
	switch kind {
	case toggleTrajectory:
		if fc.Trajectory == nil {
			fc.Trajectory = &fileTrajectory{}
		}
		fc.Trajectory.Enabled = &v
	case toggleTrajectoryRaw:
		if fc.Trajectory == nil {
			fc.Trajectory = &fileTrajectory{}
		}
		fc.Trajectory.IncludeRaw = &v
	case toggleGuardShell:
		if fc.Guards == nil {
			fc.Guards = &fileGuards{}
		}
		if fc.Guards.Shell == nil {
			fc.Guards.Shell = &fileShellGuard{}
		}
		fc.Guards.Shell.Enabled = &v
	case toggleGuardMCP:
		if fc.Guards == nil {
			fc.Guards = &fileGuards{}
		}
		if fc.Guards.MCP == nil {
			fc.Guards.MCP = &fileMCPGuard{}
		}
		fc.Guards.MCP.Enabled = &v
	case toggleGuardPaths:
		if fc.Guards == nil {
			fc.Guards = &fileGuards{}
		}
		if fc.Guards.Paths == nil {
			fc.Guards.Paths = &filePathsGuard{}
		}
		fc.Guards.Paths.Enabled = &v
	}
}

func effectiveBool(res CompileResult, kind toggleKind) bool {
	switch kind {
	case toggleTrajectory:
		return res.Trajectory.Enabled
	case toggleTrajectoryRaw:
		return res.Trajectory.IncludeRaw
	case toggleGuardShell:
		return res.Guards.Shell.Enabled
	case toggleGuardMCP:
		return res.Guards.MCP.Enabled
	case toggleGuardPaths:
		return res.Guards.Paths.Enabled
	default:
		return false
	}
}

func cloneFileConfig(in *fileConfig) *fileConfig {
	if in == nil {
		return nil
	}
	raw, err := yaml.Marshal(in)
	if err != nil {
		return in
	}
	var out fileConfig
	if err := yaml.Unmarshal(raw, &out); err != nil {
		return in
	}
	return &out
}
