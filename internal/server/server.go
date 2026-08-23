// Package server implements thin gRPC mapping for DaemonService, HookService, ConfigService, and SessionService.
//
// Owns: proto ↔ internal type mapping; HookService maps via Invoker and SnapshotSource.
// Must not: policy logic, route match, guard checks, config compile.
//
// Entry: HookService.Invoke (Invoker, SnapshotSource), ConfigService handlers, SessionService.Subscribe.
// See DESIGN.md §1.5 (invoke_sync, config_reload, async_side).
package server

import (
	"log/slog"
	"time"

	"google.golang.org/grpc"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/dispatch"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

const defaultVersion = "dev"

// Options configures the gRPC services.
type Options struct {
	Store      *config.Store
	Engine     *dispatch.Engine
	Recorder   *trajectory.Recorder
	Logger     *slog.Logger
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
	snap     SnapshotSource
	engine   Invoker
	recorder *trajectory.Recorder
	log      *slog.Logger
}

// New registers DaemonService, HookService, and ConfigService on a new gRPC server.
func New(opts Options) *grpc.Server {
	if opts.StartedAt.IsZero() {
		opts.StartedAt = time.Now().UTC()
	}
	if opts.Version == "" {
		opts.Version = defaultVersion
	}
	s := grpc.NewServer()
	agentdv1.RegisterDaemonServiceServer(s, &daemonService{opts: opts})
	var snap SnapshotSource
	if opts.Store != nil {
		snap = opts.Store
	}
	var inv Invoker
	if opts.Engine != nil {
		inv = opts.Engine
	}
	agentdv1.RegisterHookServiceServer(s, &hookService{
		snap:     snap,
		engine:   inv,
		recorder: opts.Recorder,
		log:      opts.Logger,
	})
	agentdv1.RegisterConfigServiceServer(s, &configService{store: opts.Store})
	if opts.Recorder != nil {
		agentdv1.RegisterSessionServiceServer(s, &sessionService{hub: opts.Recorder.Hub()})
	} else {
		agentdv1.RegisterSessionServiceServer(s, &sessionService{})
	}
	return s
}
