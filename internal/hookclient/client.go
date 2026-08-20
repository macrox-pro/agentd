package hookclient

import (
	"context"
	"fmt"
	"net"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps daemon and hook service clients.
type Client struct {
	conn   *grpc.ClientConn
	Daemon agentdv1.DaemonServiceClient
	Hook   agentdv1.HookServiceClient
}

// Dial connects to the daemon socket.
func Dial(ctx context.Context, socket string) (*Client, error) {
	if socket == "" {
		socket = transport.DefaultSocketPath()
	}
	path := socket
	conn, err := grpc.NewClient("passthrough:///agentd",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return transport.Dial(ctx, path)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial daemon: %w", err)
	}
	return &Client{
		conn:   conn,
		Daemon: agentdv1.NewDaemonServiceClient(conn),
		Hook:   agentdv1.NewHookServiceClient(conn),
	}, nil
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Health returns daemon health.
func (c *Client) Health(ctx context.Context) (*agentdv1.HealthResponse, error) {
	return c.Daemon.Health(ctx, &agentdv1.HealthRequest{})
}

// Status returns daemon status.
func (c *Client) Status(ctx context.Context) (*agentdv1.StatusResponse, error) {
	return c.Daemon.Status(ctx, &agentdv1.StatusRequest{})
}

// Shutdown requests daemon shutdown.
func (c *Client) Shutdown(ctx context.Context) error {
	_, err := c.Daemon.Shutdown(ctx, &agentdv1.ShutdownRequest{})
	return err
}

// Reload reloads config on the daemon.
func (c *Client) Reload(ctx context.Context) (*agentdv1.ReloadConfigResponse, error) {
	return c.Daemon.ReloadConfig(ctx, &agentdv1.ReloadConfigRequest{})
}

// Invoke sends a hook invocation.
func (c *Client) Invoke(ctx context.Context, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	return c.Hook.Invoke(ctx, req)
}
