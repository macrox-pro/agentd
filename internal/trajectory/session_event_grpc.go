package trajectory

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/timestamppb"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

// EventFromSessionEvent maps agentd.v1.SessionEvent to a ledger Event.
func EventFromSessionEvent(ev *agentdv1.SessionEvent) Event {
	if ev == nil {
		return Event{}
	}
	out := Event{
		SchemaVersion:  ev.GetSchemaVersion(),
		Seq:            ev.GetSeq(),
		Type:           ev.GetType(),
		Source:         ev.GetSource(),
		Provider:       ev.GetProvider(),
		InvocationMode: ev.GetInvocationMode(),
		SessionID:      ev.GetSessionId(),
		ProjectRoot:    ev.GetProjectRoot(),
		CWD:            ev.GetCwd(),
		Ignorable:      ev.GetIgnorable(),
	}
	if ts := ev.GetTs(); ts != nil {
		out.TS = ts.AsTime().UTC()
	}
	if len(ev.GetData()) > 0 {
		out.Data = append(json.RawMessage(nil), ev.GetData()...)
	}
	if len(ev.GetRaw()) > 0 {
		out.Raw = append(json.RawMessage(nil), ev.GetRaw()...)
	}
	return out
}

// EventToSessionEvent maps a ledger Event to agentd.v1.SessionEvent.
func EventToSessionEvent(ev Event) *agentdv1.SessionEvent {
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
