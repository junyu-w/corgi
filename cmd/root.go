package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var appVersion = "v0.1.1-alpha"

var rootCmd = &cobra.Command{
	Use:          "corgi",
	Short:        "Corgi is a smart dog that helps you organize your command flow for future usage",
	Version:      appVersion,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
