package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/config"
)

type configService struct {
	agentdv1.UnimplementedConfigServiceServer
	store *config.Store
}

func (c *configService) GetConfig(_ context.Context, req *agentdv1.GetConfigRequest) (*agentdv1.GetConfigResponse, error) {
	if c.store == nil {
		return nil, status.Error(codes.FailedPrecondition, "config store unavailable")
	}
	layer, err := protoToLayer(req.GetLayer())
	if err != nil {
		return nil, err
	}
	cwd := req.GetCwd()
	if layer == config.LayerProject || layer == config.LayerMerged {
		_, _ = c.store.EnsureProject(cwd, "")
	}
	raw, err := c.store.LayerYAML(layer, cwd, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "layer yaml: %v", err)
	}
	snap := c.store.SnapshotFor(cwd, "")
	return &agentdv1.GetConfigResponse{
		Layer:       req.GetLayer(),
		YamlContent: raw,
		Config: &agentdv1.ConfigGeneration{
			Generation:  snap.Generation,
			Fingerprint: snap.Fingerprint,
		},
	}, nil
}

func (c *configService) PatchConfig(_ context.Context, req *agentdv1.PatchConfigRequest) (*agentdv1.PatchConfigResponse, error) {
	if c.store == nil {
		return nil, status.Error(codes.FailedPrecondition, "config store unavailable")
	}
	if err := c.store.PatchRuntime(req.GetYamlPatch()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "patch runtime: %v", err)
	}
	snap := c.store.Current()
	return &agentdv1.PatchConfigResponse{
		Config: &agentdv1.ConfigGeneration{
			Generation:  snap.Generation,
			Fingerprint: snap.Fingerprint,
		},
	}, nil
}

func (c *configService) RecordDecision(_ context.Context, req *agentdv1.RecordDecisionRequest) (*agentdv1.RecordDecisionResponse, error) {
	if c.store == nil {
		return nil, status.Error(codes.FailedPrecondition, "config store unavailable")
	}
	scope, err := protoToApprovalScope(req.GetScope())
	if err != nil {
		return nil, err
	}
	opts := config.RecordDecisionOptions{
		Fingerprint: req.GetApprovalFingerprint(),
		Scope:       scope,
		Project:     req.GetProjectRoot(),
		SessionID:   req.GetSessionId(),
	}
	if ts := req.GetExpiresAt(); ts != nil {
		opts.ExpiresAt = ts.AsTime().UTC()
	}
	if err := c.store.RecordDecision(opts); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "record decision: %v", err)
	}
	snap := c.store.Current()
	return &agentdv1.RecordDecisionResponse{
		Config: &agentdv1.ConfigGeneration{
			Generation:  snap.Generation,
			Fingerprint: snap.Fingerprint,
		},
	}, nil
}

func protoToApprovalScope(l agentdv1.ConfigLayer) (config.ApprovalScope, error) {
	switch l {
	case agentdv1.ConfigLayer_CONFIG_LAYER_PROJECT:
		return config.ApprovalScopeProject, nil
	case agentdv1.ConfigLayer_CONFIG_LAYER_RUNTIME:
		return config.ApprovalScopeSession, nil
	default:
		return "", status.Errorf(codes.InvalidArgument, "scope must be PROJECT or RUNTIME (session), got %v", l)
	}
}

func protoToLayer(l agentdv1.ConfigLayer) (config.Layer, error) {
	switch l {
	case agentdv1.ConfigLayer_CONFIG_LAYER_USER:
		return config.LayerUser, nil
	case agentdv1.ConfigLayer_CONFIG_LAYER_PROJECT:
		return config.LayerProject, nil
	case agentdv1.ConfigLayer_CONFIG_LAYER_RUNTIME:
		return config.LayerRuntime, nil
	case agentdv1.ConfigLayer_CONFIG_LAYER_MERGED, agentdv1.ConfigLayer_CONFIG_LAYER_UNSPECIFIED:
		return config.LayerMerged, nil
	default:
		return "", status.Errorf(codes.InvalidArgument, "unknown layer %v", l)
	}
}
