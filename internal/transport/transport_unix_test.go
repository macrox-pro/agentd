//go:build unix

package transport_test

import (
	"context"
	"io"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/macrox-pro/agentd/internal/transport"
)

func TestListenDialRoundTrip(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "s.sock")

	ln, err := transport.Listen(path)
	require.NoError(t, err, "Listen(%q)", path)
	t.Cleanup(func() { _ = ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c io.ReadWriteCloser) {
				defer c.Close()
				_, _ = io.Copy(c, c)
			}(conn)
		}
	}()

	tests := []struct {
		name string
		msg  []byte
	}{
		{name: "ascii ping", msg: []byte("ping")},
		{name: "empty payload", msg: []byte("")},
		{name: "binary payload", msg: []byte{0x00, 0xff, 0x42}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			conn, err := transport.Dial(ctx, path)
			require.NoError(t, err, "Dial(%q)", tt.name)
			defer conn.Close()

			_, err = conn.Write(tt.msg)
			require.NoError(t, err, "Write(%q)", tt.name)

			buf := make([]byte, len(tt.msg))
			if len(tt.msg) > 0 {
				_, err = io.ReadFull(conn, buf)
				require.NoError(t, err, "Read(%q)", tt.name)
			}
			assert.Equal(t, tt.msg, buf, "round-trip(%q)", tt.name)
		})
	}
}

func TestDefaultSocketPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		wantNonEmp bool
		wantSubstr string
	}{
		{name: "non empty", wantNonEmp: true},
		{name: "contains agentd", wantSubstr: "agentd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := transport.DefaultSocketPath()
			if tt.wantNonEmp {
				assert.NotEmpty(t, got, "DefaultSocketPath()")
			}
			if tt.wantSubstr != "" {
				assert.Contains(t, got, tt.wantSubstr, "DefaultSocketPath()")
			}
		})
	}
}
