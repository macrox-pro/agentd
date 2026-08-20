package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/hookclient"
)

var configPatch string

func init() {
	configPatchCmd.Flags().StringVar(&configPatch, "file", "", "YAML file with runtime overrides to apply (required)")
	_ = configPatchCmd.MarkFlagRequired("file")
}

var configPatchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Apply temporary runtime setting overrides",
	Long: `Apply temporary overrides to the running service without editing your
user or project config files.

Requires a running agentd service. Changes apply in-memory until the next
full reload from disk (runtime persist lands in a later milestone).`,
	Example: `  agentd config patch --file runtime-delta.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		raw, err := os.ReadFile(configPatch)
		if err != nil {
			return fmt.Errorf("read patch file: %w", err)
		}
		cli, err := hookclient.Dial(cmd.Context(), resolveSocket())
		if err != nil {
			return fmt.Errorf("daemon not running: %w", err)
		}
		defer cli.Close()

		resp, err := cli.PatchConfig(cmd.Context(), &agentdv1.PatchConfigRequest{YamlPatch: raw})
		if err != nil {
			return err
		}
		cfg := resp.GetConfig()
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "generation=%d fingerprint=%s\n",
			cfg.GetGeneration(), cfg.GetFingerprint())
		return err
	},
}
