package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/config"
)

var dispatchRoutesJSON bool

func init() {
	dispatchRoutesCmd.Flags().BoolVar(&dispatchRoutesJSON, "json", false, "print routes as JSON")
}

var dispatchRoutesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List active hook routes",
	Long: `List the active hook routes loaded by the running service.

Shows match order and whether each route waits for a decision or runs in the
background. For operators and debugging; agents do not call this command.

Compiles defaults merged with the user config file offline (no daemon required).`,
	Example: `  agentd dispatch routes
  agentd dispatch routes --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		store, err := config.Load(cmd.Context(), resolveConfigPath())
		if err != nil {
			return err
		}
		return formatRoutes(cmd.OutOrStdout(), store.Current().Routes, dispatchRoutesJSON)
	},
}

func formatRoutes(w io.Writer, routes []config.CompiledRoute, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(routes)
	}
	for _, r := range routes {
		kinds := r.Kind
		if len(r.Match.Kinds) > 0 {
			kinds = strings.Join(r.Match.Kinds, ",")
		}
		syncKinds := routeTargetKinds(r.Sync)
		asyncKinds := routeTargetKinds(r.Async)
		if _, err := fmt.Fprintf(w, "%s\tmatch.kind=%s\tmode=%s\tsync=[%s]\tasync=[%s]\n",
			r.Name, kinds, r.Mode, syncKinds, asyncKinds); err != nil {
			return err
		}
	}
	return nil
}

func routeTargetKinds(ts []config.CompiledTarget) string {
	if len(ts) == 0 {
		return ""
	}
	parts := make([]string, len(ts))
	for i, t := range ts {
		parts[i] = string(t.Kind)
	}
	return strings.Join(parts, ",")
}
