package transport

import (
	"context"
	"fmt"
	"net"
	"strings"
)

const (
	unixScheme  = "unix://"
	npipeScheme = "npipe:"
)

// ParseEndpoint normalizes a grpc target endpoint to a dial path.
// Accepted forms: unix:///abs/path, bare unix path, \\.\pipe\name, npipe:\\.\pipe\name.
func ParseEndpoint(endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", fmt.Errorf("endpoint is required")
	}
	switch {
	case strings.HasPrefix(endpoint, unixScheme):
		path := strings.TrimPrefix(endpoint, unixScheme)
		if path == "" {
			return "", fmt.Errorf("unix endpoint path is empty")
		}
		return path, nil
	case strings.HasPrefix(endpoint, npipeScheme):
		path := strings.TrimPrefix(endpoint, npipeScheme)
		if path == "" {
			return "", fmt.Errorf("npipe endpoint path is empty")
		}
		return path, nil
	default:
		return endpoint, nil
	}
}

// DialEndpoint dials a grpc forward endpoint (unix socket or Windows named pipe).
func DialEndpoint(ctx context.Context, endpoint string) (net.Conn, error) {
	path, err := ParseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}
	return Dial(ctx, path)
}
