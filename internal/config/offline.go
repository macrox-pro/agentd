package config

import "context"

// OfflineFor returns compiled policy.offline for defaults ⊕ user ⊕ project(cwd) ⊕ runtime.
// Disk I/O is intentional — for the hook edge when the daemon is unreachable only.
// If opts.RuntimePath is empty, DefaultRuntimePath() is used (missing file is OK).
func OfflineFor(opts LoadOptions, cwd string) (FailMode, error) {
	if opts.RuntimePath == "" {
		opts.RuntimePath = DefaultRuntimePath()
	}
	store, err := LoadWith(context.Background(), opts)
	if err != nil {
		return "", err
	}
	snap := store.SnapshotFor(cwd, "")
	if snap == nil {
		return defaultPolicy().Offline, nil
	}
	return snap.Policy.Offline, nil
}
