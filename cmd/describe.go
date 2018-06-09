package cmd

import "github.com/spf13/cobra"

var describeCmd = &cobra.Command{
	Use:   "describe [title]",
	Short: "Describe a snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE:  describe,
}

func describe(cmd *cobra.Command, args []string) error {
	conf, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// find snippet title
	var title string
	if len(args) == 0 {
		title, err = filterSnippetTitle(conf.FilterCmd, snippetsMeta.GetSnippetTitles())
		if err != nil {
			return err
		}
	} else {
		title = args[0]
	}
	// find snippet
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
