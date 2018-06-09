package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var appVersion = "v0.2.2"

var rootCmd = &cobra.Command{
	Use:          "corgi",
	Short:        "Corgi is a smart dog that helps you manage your CLI workflow",
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
