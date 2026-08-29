//go:build !linux && !darwin && !windows

package daemon

func registerAutostart(AutostartSpec) error {
	return ErrAutostartUnsupported
}

func unregisterAutostart() error {
	return ErrAutostartUnsupported
}

func readAutostartState() (AutostartReport, error) {
	return AutostartReport{}, ErrAutostartUnsupported
}
