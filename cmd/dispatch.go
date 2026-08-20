package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/macrox-pro/agentd/internal/config"
)

func init() {
	rootCmd.AddCommand(dispatchCmd)
}

var dispatchRoutesJSON bool

var dispatchCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Inspect how hook events are routed",
	Long: `Inspect how the running service routes hook events.

Use these commands when debugging why an event took a particular path.`,
}

func init() {
	dispatchCmd.AddCommand(dispatchRoutesCmd)
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
		routes := store.Current().Routes
		if dispatchRoutesJSON {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(routes)
		}
		for _, r := range routes {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\tkind=%s\tmode=%s\tsync=%d\tasync=%d\n",
				r.Name, r.Kind, r.Mode, len(r.Sync), len(r.Async))
		}
		return nil
	},
}
