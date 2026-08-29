//go:build windows

package daemon_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macrox-pro/agentd/internal/daemon"
)

func TestAutostartWindows_schtasks_quoting(t *testing.T) {
	spec := daemon.AutostartSpecForTest(`C:\Program Files\agentd\agentd.exe`, []string{"daemon", "start", "--foreground"})
	got := daemon.SchtasksTRForTest(spec)
	assert.Contains(t, got, `"C:\Program Files\agentd\agentd.exe"`)
}

func TestAutostartWindows_parse_query_exe(t *testing.T) {
	const xml = `<Task xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task"><Actions><Exec><Command>C:\agentd.exe</Command></Exec></Actions></Task>`
	assert.Equal(t, `C:\agentd.exe`, daemon.ParseSchtasksQueryForTest(xml))
}
