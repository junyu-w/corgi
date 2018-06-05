package cmd

import "github.com/spf13/cobra"

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a snippet",
	Args:  cobra.ExactArgs(1),
	RunE:  describe,
}

func describe(cmd *cobra.Command, args []string) error {
	title := args[0]
	_, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	s, err := snippetsMeta.FindSnippet(title)
	if err != nil {
		return err
	}
	s.Describe()
	return nil
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
