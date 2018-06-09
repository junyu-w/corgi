package cmd

import (
	"corgi/snippet"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [file1] [file2...]",
	Short: "Import a snippet from one or multiple json files",
	Args:  cobra.MinimumNArgs(1),
	RunE:  importSnippet,
}

func importSnippet(cmd *cobra.Command, args []string) error {
	// load config and snippets
	conf, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// import snippet
	snippetJsonFiles := args
	for _, f := range snippetJsonFiles {
		s, err := snippet.LoadSnippet(f)
		if err != nil {
			return err
		}
		if err = snippetsMeta.SaveNewSnippet(s, conf.SnippetsDir); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(importCmd)
}
