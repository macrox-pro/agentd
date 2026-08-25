// Package hookclient provides a gRPC client for the local agentd daemon.
//
// Owns: dial transport, HookService/ConfigService/SessionService RPC wrappers.
// Must not: hook wire decode/encode (hookedge).
//
// Entry: Dial, DialReady, Client.Invoke, Client.Subscribe, Client.Health.
// See DESIGN.md §1.5 (invoke_sync, async_side).
package hookclient

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/transport"
)

const grpcDialTarget = "passthrough:///agentd"

// Client wraps daemon, hook, config, and session service clients.
type Client struct {
	conn    *grpc.ClientConn
	daemon  agentdv1.DaemonServiceClient
	hook    agentdv1.HookServiceClient
	config  agentdv1.ConfigServiceClient
	session agentdv1.SessionServiceClient
}

// Dial connects to the daemon socket.
// Note: gRPC dial is lazy — prefer DialReady on the hook path when the process
// must fail closed/open immediately if the daemon is unreachable.
func Dial(ctx context.Context, socket string) (*Client, error) {
	if socket == "" {
		socket = transport.DefaultSocketPath()
	}
	path := socket
	conn, err := grpc.NewClient(grpcDialTarget,
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return transport.Dial(ctx, path)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial daemon: %w", err)
	}
	return &Client{
		conn:    conn,
		daemon:  agentdv1.NewDaemonServiceClient(conn),
		hook:    agentdv1.NewHookServiceClient(conn),
		config:  agentdv1.NewConfigServiceClient(conn),
		session: agentdv1.NewSessionServiceClient(conn),
	}, nil
}

// DialReady dials and confirms Health. Use on the hook hot path so an unreachable
// daemon is detected before the first Invoke (grpc.NewClient connects lazily).
func DialReady(ctx context.Context, socket string) (*Client, error) {
	cli, err := Dial(ctx, socket)
	if err != nil {
		return nil, err
	}
	if _, err := cli.Health(ctx); err != nil {
		_ = cli.Close()
		return nil, fmt.Errorf("daemon health: %w", err)
	}
	return cli, nil
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
	return c.daemon.Health(ctx, &agentdv1.HealthRequest{})
}

// Status returns daemon status.
func (c *Client) Status(ctx context.Context) (*agentdv1.StatusResponse, error) {
	return c.daemon.Status(ctx, &agentdv1.StatusRequest{})
}

// Shutdown requests daemon shutdown.
func (c *Client) Shutdown(ctx context.Context) error {
	_, err := c.daemon.Shutdown(ctx, &agentdv1.ShutdownRequest{})
	return err
}

// Reload reloads config on the daemon.
func (c *Client) Reload(ctx context.Context) (*agentdv1.ReloadConfigResponse, error) {
	return c.daemon.ReloadConfig(ctx, &agentdv1.ReloadConfigRequest{})
}

// Invoke sends a hook invocation.
func (c *Client) Invoke(ctx context.Context, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	return c.hook.Invoke(ctx, req)
}

// GetConfig fetches a config layer from the daemon.
func (c *Client) GetConfig(ctx context.Context, req *agentdv1.GetConfigRequest) (*agentdv1.GetConfigResponse, error) {
	return c.config.GetConfig(ctx, req)
}

// PatchConfig applies a runtime overlay patch on the daemon.
func (c *Client) PatchConfig(ctx context.Context, req *agentdv1.PatchConfigRequest) (*agentdv1.PatchConfigResponse, error) {
	return c.config.PatchConfig(ctx, req)
}

// RecordDecision records an approval in the runtime overlay.
func (c *Client) RecordDecision(ctx context.Context, req *agentdv1.RecordDecisionRequest) (*agentdv1.RecordDecisionResponse, error) {
	return c.config.RecordDecision(ctx, req)
}

// Subscribe opens a live trajectory event stream from the daemon.
func (c *Client) Subscribe(ctx context.Context, req *agentdv1.SubscribeRequest) (agentdv1.SessionService_SubscribeClient, error) {
	return c.session.Subscribe(ctx, req)
}
