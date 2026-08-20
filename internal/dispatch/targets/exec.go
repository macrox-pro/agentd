package targets

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
)

// Exec runs an external command with optional raw stdin.
type Exec struct {
	Logger *slog.Logger
}

// InvokeAsync runs command with context timeout from the worker.
func (t *Exec) InvokeAsync(ctx context.Context, req AsyncRequest) error {
	if len(req.Target.Command) == 0 {
		return fmt.Errorf("exec target: empty command")
	}
	cmd := exec.CommandContext(ctx, req.Target.Command[0], req.Target.Command[1:]...)
	if req.Target.Stdin == "raw" {
		cmd.Stdin = bytes.NewReader(req.Raw)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		if t.Logger != nil {
			t.Logger.Warn("exec target failed", "error", err, "output", string(out))
		}
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
