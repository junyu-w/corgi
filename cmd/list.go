package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snippets",
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	// load config & snippets
	_, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// display
	fmt.Println("Here is the list of corgi snippets saved on your system:")
	for _, s := range snippetsMeta.Snippets {
		color.Yellow("- %s", s.Title)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
