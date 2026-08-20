package server

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

// Options configures the gRPC services.
type Options struct {
	Store      *config.Store
	Engine     *dispatch.Engine
	StartedAt  time.Time
	Version    string
	OnShutdown func()
	Log        *slog.Logger
}

type daemonService struct {
	agentdv1.UnimplementedDaemonServiceServer
	opts Options
}

type hookService struct {
	agentdv1.UnimplementedHookServiceServer
	store  *config.Store
	engine *dispatch.Engine
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
	agentdv1.RegisterHookServiceServer(s, &hookService{store: opts.Store, engine: opts.Engine})
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
		CompiledRouteCount: uint32(len(snap.Routes)),
	}
	if d.opts.Engine != nil {
		if q := d.opts.Engine.Queue(); q != nil {
			resp.AsyncQueueDepth = uint32(q.Depth())
		}
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

func (h *hookService) Invoke(ctx context.Context, req *agentdv1.InvokeRequest) (*agentdv1.InvokeResponse, error) {
	snap := h.store.Current()
	resp := &agentdv1.InvokeResponse{
		Config: &agentdv1.ConfigGeneration{
			Generation:  snap.Generation,
			Fingerprint: snap.Fingerprint,
		},
		Decision: dispatch.NeutralDecision(),
	}
	if h.engine == nil {
		return resp, nil
	}
	in := dispatch.InvokeInput{
		Provider:   req.GetProvider(),
		RawPayload: req.GetRawPayload(),
		Snap:       snap,
	}
	if dl := req.GetDeadline(); dl != nil {
		in.Deadline = dl.AsTime()
	}
	result, err := h.engine.Invoke(ctx, in)
	if err != nil {
		// Match agenthooks wire: undecodable payloads become a neutral no-op.
		return resp, nil
	}
	resp.Decision = result.Decision
	resp.AsyncDispatchedCount = result.AsyncDispatchedCount
	return resp, nil
}
