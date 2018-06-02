package cmd

import "github.com/spf13/cobra"

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a snippet",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Do stuff
		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
