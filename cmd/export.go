package cmd

import "github.com/spf13/cobra"

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export snippet definition to shell script",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Do stuff
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
