package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/macrox-pro/agentd/internal/config"
	"github.com/macrox-pro/agentd/internal/transport"
)

var (
	cfgFile    string
	socketPath string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "agentd",
	Short: "Control coding-agent hooks from one background service",
	Long: `agentd runs a background service that applies your hook policies when
coding agents (Claude Code, Cursor, Codex, Gemini CLI, OpenCode, Kimi Code)
fire hooks.

Use "agentd daemon" to start and manage the service, and "agentd hook" as the
command agents call from their hook settings.`,
	Example: `  agentd daemon start
  agentd hook run --provider=claude-code`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// RootCommand returns the root cobra command (for external tests).
func RootCommand() *cobra.Command {
	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file (default $HOME/.agentd.yaml)")
	rootCmd.PersistentFlags().StringVar(&socketPath, "socket", "", "path to the daemon socket")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "print extra messages to stderr")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("socket", rootCmd.PersistentFlags().Lookup("socket"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".agentd")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "using config file:", viper.ConfigFileUsed())
	}
}

func resolveSocket() string {
	if socketPath != "" {
		return socketPath
	}
	return transport.DefaultSocketPath()
}

func resolveConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	return config.DefaultUserPath()
}
