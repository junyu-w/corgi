package cmd

import (
	"corgi/config"
	"corgi/snippet"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all snippets",
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	// load config
	conf, err := config.Load()
	if err != nil {
		return err
	}
	// Load snippets
	snippets, err := snippet.LoadSnippets(conf.SnippetsFile)
	// display
	fmt.Println("Here is the list of corgi snippets saved on your system:")
	for _, s := range snippets.Snippets {
		color.Yellow("[ %s ]", s.Title)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
