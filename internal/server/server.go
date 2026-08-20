package server

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
)

// Options configures the gRPC services.
type Options struct {
	Store      *config.Store
	StartedAt  time.Time
	Version    string
	OnShutdown func()
}

type daemonService struct {
	agentdv1.UnimplementedDaemonServiceServer
	opts Options
}

type hookService struct {
	agentdv1.UnimplementedHookServiceServer
	store *config.Store
}

// New registers DaemonService and HookService on a new gRPC server.
func New(opts Options) *grpc.Server {
	if opts.StartedAt.IsZero() {
		opts.StartedAt = time.Now().UTC()
	}
	if opts.Version == "" {
		opts.Version = "dev"
	}
	s := grpc.NewServer()
	agentdv1.RegisterDaemonServiceServer(s, &daemonService{opts: opts})
	agentdv1.RegisterHookServiceServer(s, &hookService{store: opts.Store})
	return s
}

func (d *daemonService) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: "ok"}, nil
}

func (d *daemonService) Status(context.Context, *agentdv1.StatusRequest) (*agentdv1.StatusResponse, error) {
	snap := d.opts.Store.Current()
	resp := &agentdv1.StatusResponse{
		Version:   d.opts.Version,
		StartedAt: timestamppb.New(d.opts.StartedAt),
		Config: &agentdv1.ConfigGeneration{
			Generation:  snap.Generation,
			Fingerprint: snap.Fingerprint,
		},
	}
	if snap.UserPath != "" {
		resp.ConfigLayers = []*agentdv1.LayerInfo{{
			Layer:       agentdv1.ConfigLayer_CONFIG_LAYER_USER,
			Path:        snap.UserPath,
			Fingerprint: snap.Fingerprint,
		}}
	}
	return resp, nil
}

func (d *daemonService) ReloadConfig(ctx context.Context, _ *agentdv1.ReloadConfigRequest) (*agentdv1.ReloadConfigResponse, error) {
	if err := d.opts.Store.Reload(ctx); err != nil {
		return nil, err
	}
	snap := d.opts.Store.Current()
	return &agentdv1.ReloadConfigResponse{
		Config: &agentdv1.ConfigGeneration{
			Generation:  snap.Generation,
			Fingerprint: snap.Fingerprint,
		},
	}, nil
}

func (d *daemonService) Shutdown(context.Context, *agentdv1.ShutdownRequest) (*agentdv1.ShutdownResponse, error) {
	if d.opts.OnShutdown != nil {
		go d.opts.OnShutdown()
	}
	return &agentdv1.ShutdownResponse{}, nil
}

func (h *hookService) Invoke(_ context.Context, _ *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	snap := h.store.Current()
	return &agentdv1.InvokeResponse{
		Decision: &agentdv1.Decision{
			Kind: agentdv1.DecisionKind_DECISION_KIND_NO_DECISION,
		},
		Config: &agentdv1.ConfigGeneration{
			Generation:  snap.Generation,
			Fingerprint: snap.Fingerprint,
		},
		AsyncDispatchedCount: 0,
	}, nil
}
