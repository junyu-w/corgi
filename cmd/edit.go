package cmd

import (
	"corgi/util"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a snippet",
	Args:  cobra.ExactArgs(1),
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	title := args[0]
	// load config & snippets
	config, snippets, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// find snippet
	s, err := snippets.FindSnippet(title)
	if err != nil {
		return err
	}
	command := fmt.Sprintf("%s %s", config.Editor, s.GetFilePath())
	if err := util.Execute(command, os.Stdin, os.Stdout); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(editCmd)
}
