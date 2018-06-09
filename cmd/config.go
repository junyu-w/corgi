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

func configure(cmd *cobra.Command, args []string) error {
	conf, _, err := loadConfigAndSnippetsMeta()
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
	return nil
}

func init() {
	configCmd.Flags().StringVar(&filterCmd, "filter-cmd", "", "Select the text editor you would like to use to edit snippet")
	configCmd.Flags().StringVar(&editor, "editor", "", "Select the text editor you would like to use to edit snippet")
	rootCmd.AddCommand(configCmd)
}
