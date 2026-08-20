//go:build !unix && !windows

package transport

// DefaultSocketPath returns empty on unsupported platforms.
func DefaultSocketPath() string {
	return ""
}
