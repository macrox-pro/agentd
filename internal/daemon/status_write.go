package daemon

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// WriteStatus writes a StatusReport as JSON or a short human line.
func WriteStatus(w io.Writer, rep StatusReport, asJSON bool) error {
	if asJSON {
		payload := map[string]any{
			"running": rep.Running,
			"socket":  rep.Socket,
		}
		if rep.Running {
			payload["version"] = rep.Version
			payload["started_at"] = rep.StartedAt.UTC().Format(time.RFC3339)
			payload["generation"] = rep.Generation
			payload["fingerprint"] = rep.Fingerprint
			payload["async_queue_depth"] = rep.AsyncQueueDepth
			payload["async_dropped_count"] = rep.AsyncDroppedCount
			payload["trajectory_dropped_count"] = rep.TrajectoryDroppedCount
			payload["compiled_route_count"] = rep.CompiledRouteCount
		}
		return json.NewEncoder(w).Encode(payload)
	}
	if !rep.Running {
		_, err := fmt.Fprintln(w, "agentd: not running")
		return err
	}
	_, err := fmt.Fprintf(w, "agentd: running (version %s, generation %d)\n",
		rep.Version, rep.Generation)
	return err
}
