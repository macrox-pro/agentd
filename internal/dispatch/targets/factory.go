package targets

import (
	"fmt"
	"log/slog"

	"github.com/macrox-pro/agentd/internal/config"
)

// NewAsyncInvoker builds an AsyncInvoker for a compiled target.
func NewAsyncInvoker(t config.CompiledTarget, builtin *Builtin, log *slog.Logger) (AsyncInvoker, error) {
	switch t.Kind {
	case config.TargetBuiltin:
		return &BuiltinObserve{Inner: builtin}, nil
	case config.TargetLog:
		return &Log{Logger: log}, nil
	case config.TargetFile:
		return &File{Logger: log}, nil
	case config.TargetHTTP:
		return &HTTP{Logger: log}, nil
	case config.TargetExec:
		return &Exec{Logger: log}, nil
	case config.TargetGRPC:
		return &GRPC{Logger: log}, nil
	default:
		return nil, fmt.Errorf("unknown async target %q", t.Kind)
	}
}
