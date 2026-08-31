package daemon

import (
	"context"
	"net/http"
	"os"
)

func ServiceStartArgsForTest(opts StartOptions) []string { return serviceStartArgs(opts) }

func StartMetricsServerForTest(ctx context.Context, listen string, handler http.Handler) (addr string, shutdown func(context.Context), err error) {
	ms, err := startMetricsServer(ctx, listen, handler)
	if err != nil {
		return "", nil, err
	}
	return ms.addr(), ms.shutdown, nil
}

func AutostartSpecForTest(exe string, args []string) AutostartSpec {
	return AutostartSpec{Exe: exe, Args: args}
}

func RenderSystemdUnitForTest(spec AutostartSpec) string { return renderSystemdUnit(spec) }

func RenderLaunchdPlistForTest(spec AutostartSpec) string { return renderLaunchdPlist(spec) }

func ParseSystemdExecStartForTest(unit string) string { return parseSystemdExecStart(unit) }

func ParseLaunchdProgramForTest(plist string) string { return parseLaunchdProgram(plist) }

func SchtasksTRForTest(spec AutostartSpec) string { return schtasksTR(spec) }

func ParseSchtasksQueryForTest(xmlRaw string) string { return parseSchtasksQuery(xmlRaw) }

func WriteFileAtomicForTest(path string, data []byte, perm int) error {
	return writeFileAtomic(path, data, os.FileMode(perm))
}
