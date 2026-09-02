//go:build windows

package transport_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"

	"github.com/macrox-pro/agentd/internal/transport"
)

func TestDefaultSocketPathUsesSID(t *testing.T) {
	t.Parallel()
	path := transport.DefaultSocketPath()
	assert.True(t, strings.HasPrefix(path, `\\.\pipe\agentd-`), "path=%q", path)

	token, err := windows.OpenCurrentProcessToken()
	require.NoError(t, err)
	defer token.Close()
	tokUser, err := token.GetTokenUser()
	require.NoError(t, err)
	assert.Contains(t, path, tokUser.User.Sid.String())
}

func TestIsPipePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "default socket", in: transport.DefaultSocketPath(), want: true},
		{name: "pipe path", in: `\\.\pipe\agentd-test`, want: true},
		{name: "upper case namespace", in: `\\.\PIPE\agentd-test`, want: true},
		{name: "namespace without name", in: `\\.\pipe\`, want: false},
		{name: "file path", in: `C:\Temp\agentd\s.sock`, want: false},
		{name: "empty", in: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, transport.IsPipePath(tt.in), "IsPipePath(%q)", tt.in)
		})
	}
}

func TestListenDialRoundTripWindows(t *testing.T) {
	t.Parallel()
	path := `\\.\pipe\agentd-test-` + strings.ReplaceAll(t.Name(), "/", "-")
	ln, err := transport.Listen(path)
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
	conn, err := transport.Dial(ctx, path)
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func TestParseEndpointNpipe(t *testing.T) {
	t.Parallel()
	got, err := transport.ParseEndpoint(`npipe:\\.\pipe\agentd-sid`)
	require.NoError(t, err)
	assert.Equal(t, `\\.\pipe\agentd-sid`, got)
}
