package daemon_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestRenderSystemdUnit(t *testing.T) {
	t.Parallel()
	spec := daemon.AutostartSpecForTest("/usr/bin/agentd", []string{"daemon", "start", "--foreground", "--config", "/cfg.yaml"})
	got := daemon.RenderSystemdUnitForTest(spec)
	require.Contains(t, got, "ExecStart=/usr/bin/agentd")
	require.Contains(t, got, "--config")
	exe := daemon.ParseSystemdExecStartForTest(got)
	assert.Equal(t, "/usr/bin/agentd", exe)
}

func TestRenderLaunchdPlist(t *testing.T) {
	t.Parallel()
	spec := daemon.AutostartSpecForTest("/usr/bin/agentd", []string{"daemon", "start", "--foreground"})
	got := daemon.RenderLaunchdPlistForTest(spec)
	exe := daemon.ParseLaunchdProgramForTest(got)
	assert.Equal(t, "/usr/bin/agentd", exe)
}

func TestSchtasksTR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		exe  string
		want string
	}{
		{name: "windows_exe_with_spaces", exe: `C:\Program Files\agentd\agentd.exe`, want: `"C:\Program Files\agentd\agentd.exe"`},
		{name: "plain", exe: `C:\agentd.exe`, want: `C:\agentd.exe`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			spec := daemon.AutostartSpecForTest(tt.exe, []string{"daemon", "start", "--foreground"})
			got := daemon.SchtasksTRForTest(spec)
			assert.Contains(t, got, tt.want)
		})
	}
}

func TestParseSchtasksQuery(t *testing.T) {
	t.Parallel()
	const xml = `<Task><Actions><Exec><Command>C:\agentd.exe</Command></Exec></Actions></Task>`
	assert.Equal(t, `C:\agentd.exe`, daemon.ParseSchtasksQueryForTest(xml))
}

func TestCorruptManifestParse(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "", daemon.ParseSystemdExecStartForTest("not a unit"))
	assert.Equal(t, "", daemon.ParseLaunchdProgramForTest("<broken"))
}
