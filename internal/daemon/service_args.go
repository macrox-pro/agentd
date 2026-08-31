package daemon

// serviceStartArgs builds argv for "agentd daemon start --foreground" used by detach and OS autostart manifests.
func serviceStartArgs(opts StartOptions) []string {
	args := []string{"daemon", "start", "--foreground"}
	if opts.ConfigPath != "" {
		args = append(args, "--config", opts.ConfigPath)
	}
	if opts.Socket != "" {
		args = append(args, "--socket", opts.Socket)
	}
	if opts.LogLevel != "" {
		args = append(args, "--log-level", opts.LogLevel)
	}
	if opts.LogFile != "" {
		args = append(args, "--log-file", opts.LogFile)
	}
	if opts.MetricsListen != "" {
		args = append(args, "--metrics-listen", opts.MetricsListen)
	}
	return args
}
