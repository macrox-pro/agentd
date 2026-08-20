package config

import "fmt"

// DispatchMode is a route dispatch mode.
type DispatchMode string

const (
	ModeSyncOnly      DispatchMode = "sync_only"
	ModeAsyncOnly     DispatchMode = "async_only"
	ModeParallel      DispatchMode = "parallel"
	ModeAfterSync     DispatchMode = "after_sync"
	ModeSyncThenAsync DispatchMode = "sync_then_async" // alias for after_sync
)

func parseDispatchMode(s string) (DispatchMode, error) {
	switch DispatchMode(s) {
	case ModeSyncOnly, ModeAsyncOnly, ModeParallel, ModeAfterSync, ModeSyncThenAsync:
		return DispatchMode(s), nil
	default:
		return "", fmt.Errorf("unknown %q", s)
	}
}

// NormalizeMode maps aliases to canonical modes.
func NormalizeMode(m DispatchMode) DispatchMode {
	if m == ModeSyncThenAsync {
		return ModeAfterSync
	}
	return m
}
