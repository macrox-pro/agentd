package server

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

type sessionService struct {
	agentdv1.UnimplementedSessionServiceServer
	hub *trajectory.Hub
}

func (s *sessionService) Subscribe(req *agentdv1.SubscribeRequest, stream agentdv1.SessionService_SubscribeServer) error {
	if s.hub == nil {
		<-stream.Context().Done()
		return stream.Context().Err()
	}
	filter := trajectory.SubscribeFilter{
		Provider:  req.GetProvider(),
		SessionID: req.GetSessionId(),
		Source:    req.GetSource(),
	}
	ch, unregister := s.hub.Register(filter)
	defer unregister()

	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(&agentdv1.SubscribeResponse{Event: sessionEventToProto(ev)}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func sessionEventToProto(ev trajectory.Event) *agentdv1.SessionEvent {
	out := &agentdv1.SessionEvent{
		SchemaVersion:  ev.SchemaVersion,
		Seq:            ev.Seq,
		Type:           ev.Type,
		Source:         ev.Source,
		Provider:       ev.Provider,
		InvocationMode: ev.InvocationMode,
		SessionId:      ev.SessionID,
		ProjectRoot:    ev.ProjectRoot,
		Cwd:            ev.CWD,
		Ignorable:      ev.Ignorable,
	}
	if !ev.TS.IsZero() {
		out.Ts = timestamppb.New(ev.TS)
	}
	if len(ev.Data) > 0 {
		out.Data = append([]byte(nil), ev.Data...)
	}
	if len(ev.Raw) > 0 {
		out.Raw = append([]byte(nil), ev.Raw...)
	}
	return out
}
