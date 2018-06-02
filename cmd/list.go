package cmd

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all snippets",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Do stuff
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
