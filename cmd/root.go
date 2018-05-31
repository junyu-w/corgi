package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "corgi",
	Short: "Corgi is a smart dog that helps you organize your command flow for future usage",
	Run: func(cmd *cobra.Command, args []string) {
		// Do stuff here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
