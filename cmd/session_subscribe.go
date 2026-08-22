package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/hookclient"
	"github.com/macrox-pro/agentd/internal/trajectory"
)

var (
	sessionSubscribeProvider string
	sessionSubscribeSession  string
	sessionSubscribeSource   string
	sessionSubscribeJSON     bool
)

func init() {
	sessionSubscribeCmd.Flags().StringVar(&sessionSubscribeProvider, "provider", "", "filter by provider id")
	sessionSubscribeCmd.Flags().StringVar(&sessionSubscribeSession, "session", "", "filter by session id")
	sessionSubscribeCmd.Flags().StringVar(&sessionSubscribeSource, "source", "", "filter by source (hook, decision, transcript, system)")
	sessionSubscribeCmd.Flags().BoolVar(&sessionSubscribeJSON, "json", false, "print one JSON object per line (NDJSON; ledger field names)")
}

var sessionSubscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Stream live trajectory events from the daemon",
	Long: `Stream trajectory ledger events from the running daemon (live firehose from dial time).

Requires a running daemon with trajectory.enabled. History is not replayed — use
session show or session export for past events. Offline import/fork paths do not
publish to this stream unless the daemon Claude import watcher appends live.`,
	Example: `  agentd session subscribe --json
  agentd session subscribe --provider claude-code --json
  agentd session subscribe --provider cursor --session s1 --source hook`,
	RunE: runSessionSubscribe,
}

func runSessionSubscribe(cmd *cobra.Command, _ []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cli, err := hookclient.Dial(ctx, resolveSocket())
	if err != nil {
		return err
	}
	defer cli.Close()

	stream, err := cli.Subscribe(ctx, &agentdv1.SubscribeRequest{
		Provider:  sessionSubscribeProvider,
		SessionId: sessionSubscribeSession,
		Source:    sessionSubscribeSource,
	})
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	out := cmd.OutOrStdout()
	enc := json.NewEncoder(out)
	var flush func() error
	if w, ok := out.(*os.File); ok {
		bw := bufio.NewWriter(w)
		enc = json.NewEncoder(bw)
		flush = bw.Flush
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		ev := subscribeEventToLedger(msg.GetEvent())
		if sessionSubscribeJSON {
			if err := enc.Encode(ev); err != nil {
				return err
			}
			if flush != nil {
				if err := flush(); err != nil {
					return err
				}
			}
			continue
		}
		fmt.Fprintf(out, "%d\t%s\t%s\t%s\n",
			ev.Seq, ev.Type, ev.Source, ev.Provider)
	}
}

func subscribeEventToLedger(ev *agentdv1.SessionEvent) trajectory.Event {
	if ev == nil {
		return trajectory.Event{}
	}
	out := trajectory.Event{
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
