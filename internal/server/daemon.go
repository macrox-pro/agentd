package server

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
)

const healthStatusOK = "ok"

func (d *daemonService) Health(context.Context, *agentdv1.HealthRequest) (*agentdv1.HealthResponse, error) {
	return &agentdv1.HealthResponse{Status: healthStatusOK}, nil
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
		if s := d.opts.Engine.Sessions(); s != nil {
			resp.ActiveSessions = s.Active()
		}
	}
	resp.ConfigLayers = layerInfos(d.opts.Store, snap)
	return resp, nil
}

func layerInfos(store *config.Store, snap *config.Snapshot) []*agentdv1.LayerInfo {
	if store == nil || snap == nil {
		return nil
	}
	var out []*agentdv1.LayerInfo
	if p := store.UserPath(); p != "" {
		out = append(out, &agentdv1.LayerInfo{
			Layer:       agentdv1.ConfigLayer_CONFIG_LAYER_USER,
			Path:        p,
			Fingerprint: snap.Fingerprint,
		})
	}
	if p := store.RuntimePath(); p != "" {
		out = append(out, &agentdv1.LayerInfo{
			Layer: agentdv1.ConfigLayer_CONFIG_LAYER_RUNTIME,
			Path:  p,
		})
	}
	for _, p := range store.ProjectPaths() {
		out = append(out, &agentdv1.LayerInfo{
			Layer: agentdv1.ConfigLayer_CONFIG_LAYER_PROJECT,
			Path:  p,
		})
	}
	return out
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
