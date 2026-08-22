package server

import (
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
			if err := stream.Send(&agentdv1.SubscribeResponse{Event: trajectory.EventToSessionEvent(ev)}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}
