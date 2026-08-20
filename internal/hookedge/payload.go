package hookedge

import (
	"bytes"
	"fmt"
	"io"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
)

const maxPayloadBytes = 8 << 20

func readPayload(opts Options) ([]byte, agentdv1.InvocationMode, error) {
	if opts.ArgvPayload {
		if opts.PayloadArg == "" {
			return nil, 0, fmt.Errorf("empty argv payload")
		}
		return []byte(opts.PayloadArg), agentdv1.InvocationMode_INVOCATION_MODE_ARGV, nil
	}
	stdin := opts.Stdin
	if stdin == nil {
		stdin = bytes.NewReader(nil)
	}
	payload, err := io.ReadAll(io.LimitReader(stdin, maxPayloadBytes))
	if err != nil {
		return nil, 0, fmt.Errorf("read stdin: %w", err)
	}
	if len(payload) == 0 {
		return nil, 0, fmt.Errorf("empty stdin")
	}
	return payload, agentdv1.InvocationMode_INVOCATION_MODE_STDIN, nil
}
