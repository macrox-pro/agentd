// Package server implements thin gRPC mapping for DaemonService, HookService, and ConfigService.
//
// Owns: proto ↔ internal type mapping; delegates Invoke to dispatch.Engine.
// Must not: policy logic, route match, guard checks, config compile.
//
// Entry: HookService.Invoke, ConfigService handlers.
// See DESIGN.md §1.5 (invoke_sync, config_reload).
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
	store    *config.Store
	engine   *dispatch.Engine
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
	agentdv1.RegisterHookServiceServer(s, &hookService{
		store:    opts.Store,
		engine:   opts.Engine,
		recorder: opts.Recorder,
		log:      opts.Logger,
	})
	agentdv1.RegisterConfigServiceServer(s, &configService{store: opts.Store})
	return s
}
