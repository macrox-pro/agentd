package cmd

import (
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
		return config.FormatRoutes(cmd.OutOrStdout(), store.Current().Routes, dispatchRoutesJSON)
	},
}
