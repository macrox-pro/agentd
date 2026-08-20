//go:build unix

package transport_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/transport"
)

func TestDialEndpointUnix(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	sock := filepath.Join(dir, "peer.sock")
	ln, err := transport.Listen(sock)
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		_ = conn.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := transport.DialEndpoint(ctx, "unix://"+sock)
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}
