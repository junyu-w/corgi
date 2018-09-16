package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Update corgi configuration",
	RunE:  configure,
}

var editor string
var filterCmd string
var snippetsDir string

func configure(cmd *cobra.Command, args []string) error {
	conf, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	if editor != "" {
		conf.Editor = editor
		if err := conf.Save(); err != nil {
			return err
		}
	}
	if filterCmd != "" {
		conf.FilterCmd = filterCmd
		if err := conf.Save(); err != nil {
			return err
		}
	}
	if snippetsDir != "" {
		conf.SnippetsDir = snippetsDir
		if err := conf.Save(); err != nil {
			return err
		}
		snippetsMeta.IsMetaDirty = true
		if err := snippetsMeta.Save(); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	configCmd.Flags().StringVar(&filterCmd, "filter-cmd", "", "Set the filter command to use for fuzzy searching snippet (default to fzf)")
	configCmd.Flags().StringVar(&editor, "editor", "", "Set the text editor you would like to use to edit snippet")
	configCmd.Flags().StringVar(&snippetsDir, "snippets-dir", "", "Set the path where all snippets are located")
	rootCmd.AddCommand(configCmd)
}
