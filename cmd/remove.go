package cmd

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:   "remove [title]",
	Short: "Remove a snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE:  remove,
}

func remove(cmd *cobra.Command, args []string) error {
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
	if err = snippetsMeta.DeleteSnippet(title); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
